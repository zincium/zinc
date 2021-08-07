package main

import (
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
)

var log *zap.Logger
var sugar *zap.SugaredLogger

// zinc-secured

func main() {
	log, _ = zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	sugar = log.Sugar()
	var opts Options
	if err := opts.ParseArgv(); err != nil {
		fmt.Fprintf(os.Stderr, "ParseArgv: %v\n", err)
		os.Exit(1)
	}
	srv := NewServer(&opts)
	sugar.Info("zincs listen: ", opts.Listen)
	if err := srv.ListenAndServe(opts.Listen); err != http.ErrServerClosed {
		sugar.Fatalf("ListenAndServer %v", err)
	}
}
