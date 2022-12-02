package cmd

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/cry0genic/go-stocks/api"
	"github.com/cry0genic/go-stocks/finance"
	"github.com/cry0genic/go-stocks/finance/iexcloud"
	"github.com/cry0genic/go-stocks/history/sqlite"
	"github.com/cry0genic/go-stocks/poll"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	rootCmd = &cobra.Command{
		Use:   "stonks",
		Short: "Stonks helps you track your financial positions, questionable or otherwise.",
		Long: `Stonks helps you track your financial positions, questionable or otherwise.

https://www.urbandictionary.com/define.php?term=Stonks`,
		PreRun: rootPreRun,
		Run:    rootRun,
	}

	logOutput zapcore.WriteSyncer = os.Stdout
	logLevel                      = zapcore.DebugLevel
)

func init() {
	cobra.OnInitialize(func() {
		replacer := strings.NewReplacer("-", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.SetEnvPrefix("stonks")
		viper.AutomaticEnv()
	})

	rootCmd.Flags().Duration("api-idle-timeout", api.DefaultIdleTimeout, "duration clients are allowed to idle")
	rootCmd.Flags().StringP("api-listen-addr", "a", api.DefaultListenAddress, "API server host:port")
	rootCmd.Flags().Bool("api-metrics", true, "enable metrics for the API server")
	rootCmd.Flags().Duration("api-read-headers-timeout", api.DefaultReadHeaderTimeout, "duration clients have to send request headers")

	rootCmd.Flags().String("iex-batch-endpoint", iexcloud.DefaultBatchEndpoint, "IEX Cloud API batch endpoint URL")
	rootCmd.Flags().Duration("iex-call-timeout", iexcloud.DefaultTimeout, "API call timeout")
	rootCmd.Flags().Bool("iex-metrics", false, "collect metrics for IEX Cloud API calls")
	rootCmd.Flags().StringP("iex-token", "t", "", "IEX Cloud API token")

	rootCmd.Flags().StringP("log", "l", "stdout", "log file path")
	rootCmd.Flags().Bool("log-compress", false, "compress rotated log files")
	rootCmd.Flags().Bool("log-localtime", false, "log file names use local time, UTC otherwise")
	rootCmd.Flags().Int("log-max-age", 7, "max days to retain old log files")
	rootCmd.Flags().Int("log-max-backups", 5, "max number of old log files to retain")
	rootCmd.Flags().Int("log-max-size", 100, "max log file size in MB before rotation")

	rootCmd.Flags().Duration("sqlite-conn-max-lifetime", sqlite.DefaultConnsMaxLifetime, "max client connection lifetime")
	rootCmd.Flags().StringP("sqlite-database", "d", sqlite.DefaultDatabaseFile, "database file path")
	rootCmd.Flags().Int("sqlite-max-idle-conn", sqlite.DefaultMaxIdleConns, "max idle client connections")

	rootCmd.Flags().DurationP("poll", "p", poll.DefaultPollDuration, "duration between stock quote updates")
	rootCmd.Flags().String("pprof-addr", ":6060", "pprof host:port")
	rootCmd.Flags().StringSliceP("symbols", "s", finance.DefaultSymbols, "stock symbols")
	rootCmd.Flags().BoolP("verbose", "v", true, "verbose logging")

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		log.Fatalf("binding flags to viper: %s", err)
	}
}

func rootPreRun(_ *cobra.Command, _ []string) {
	if viper.GetString("iex-token") == "" {
		log.Fatal("IEX Cloud API token not set")
	}

	switch strings.ToLower(viper.GetString("log")) {
	case "stdout", "":
	default:
		logOutput = zapcore.AddSync(
			&lumberjack.Logger{
				Filename:   viper.GetString("log"),
				Compress:   viper.GetBool("log-compress"),
				LocalTime:  viper.GetBool("log-localtime"),
				MaxAge:     viper.GetInt("log-max-age"),
				MaxBackups: viper.GetInt("log-max-backups"),
				MaxSize:    viper.GetInt("log-max-size"),
			},
		)
	}

	if !viper.GetBool("verbose") {
		logLevel = zapcore.InfoLevel
	}
}

func rootRun(_ *cobra.Command, _ []string) {
	ret := 0
	defer os.Exit(ret)

	ctx, cancel := context.WithCancel(context.Background())
	zl := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			logOutput,
			logLevel,
		),
	).Sugar()
	defer func() { _ = zl.Sync() }()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		zl.Info("shutting down ...")
		cancel()
	}()

	go func() {
		_ = http.ListenAndServe(viper.GetString("pprof-addr"), nil)
	}()

	storage, err := sqlite.New(
		sqlite.ConnMaxLifetime(viper.GetDuration("sqlite-conn-max-lifetime")),
		sqlite.DatabaseFile(viper.GetString("sqlite-database")),
		sqlite.MaxIdleConnections(viper.GetInt("sqlite-max-idle-conn")),
		sqlite.Symbols(viper.GetStringSlice("symbols")),
	)
	if err != nil {
		zl.Error(err)
		gracefulExit(cancel, &ret)
	}

	defer func() {
		if err := storage.Close(); err != nil {
			zl.Errorf("closing archiver: %v", err)
		}
		zl.Debug("archiver closed")
	}()

	var iexMetrics iexcloud.Option
	if viper.GetBool("iex-metrics") {
		iexMetrics = iexcloud.InstrumentHTTPClient()
	}

	quotes, err := iexcloud.New(
		viper.GetString("iex-token"),
		iexcloud.BatchEndpoint(viper.GetString("iex-batch-endpoint")),
		iexcloud.CallTimeout(viper.GetDuration("iex-call-timeout")),
		iexMetrics,
	)
	if err != nil {
		zl.Error(err)
		gracefulExit(cancel, &ret)
	}

	poller, err := poll.New(quotes, storage, zl)
	if err != nil {
		zl.Error(err)
		gracefulExit(cancel, &ret)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		poller.Poll(
			ctx,
			viper.GetDuration("poll"),
			viper.GetStringSlice("symbols")...,
		)
		wg.Done()
	}()

	var apiMetrics api.Option
	if !viper.GetBool("api-metrics") {
		apiMetrics = api.DisableInstrumentation()
	}
	server, err := api.New(
		ctx, storage, zl,
		apiMetrics,
		api.IdleTimeout(viper.GetDuration("api-idle-timeout")),
		api.ListenAddress(viper.GetString("api-listen-addr")),
		api.ReadHeaderTimeout(viper.GetDuration("api-read-headers-timeout")),
	)
	if err != nil {
		zl.Error(err)
		gracefulExit(cancel, &ret)
	}

	wg.Add(1)
	go func() {
		sErr := server.ListenAndServe()
		if sErr != nil && sErr != http.ErrServerClosed {
			zl.Errorf("API server: %v", sErr)
		}
		wg.Done()
	}()

	wg.Wait()
}

func gracefulExit(cancel context.CancelFunc, ret *int) {
	cancel()
	*ret = 1
	runtime.Goexit()
}

func Execute() error {
	return rootCmd.Execute()
}
