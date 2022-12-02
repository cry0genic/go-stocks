package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/cry0genic/go-stocks/finance"
	"github.com/cry0genic/go-stocks/history"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func stock(p history.Provider, log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err  error
			last int
		)
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()

		vars := mux.Vars(r)
		symbol, ok := vars["symbol"]
		if !ok || symbol == "" {
			log.Errorw("symbol not found in request URI!", "uri", r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal server error"))
			return
		}

		if l := r.URL.Query().Get("last"); l != "" {
			last, err = strconv.Atoi(l)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`Invalid "last" parameter`))
				return
			}
		}

		quotes, err := p.GetQuotes(r.Context(), strings.ToLower(symbol), last)
		if err != nil {
			if err == history.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("Not found"))
			} else {
				log.Error(err, zap.String("url", r.URL.String()))
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal server error"))
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(quotes)
		if err != nil {
			log.Warn(err)
		}
	}
}

func stocks(p history.Provider, log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err  error
			last int
		)
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()

		if l := r.URL.Query().Get("last"); l != "" {
			last, err = strconv.Atoi(l)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`Invalid "last" parameter`))
				return
			}
		}

		batch, err := p.GetQuotesBatch(r.Context(), finance.DefaultSymbols, last)
		if err != nil {
			if err == history.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("Not found"))
			} else {
				log.Error(err, zap.String("url", r.URL.String()))
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal server error"))
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(batch)
		if err != nil {
			log.Warn(err)
		}
	}
}
