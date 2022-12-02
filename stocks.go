package main

import (
	"log"

	"github.com/cry0genic/go-stocks/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
