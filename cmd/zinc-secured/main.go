package main

import (
	"fmt"
	"os"
)

// zinc-secured

func main() {
	var opts Options
	if err := opts.ParseArgv(); err != nil {
		fmt.Fprintf(os.Stderr, "ParseArgv: %v\n", err)
		os.Exit(1)
	}
	srv := &Server{GitPath: opts.GitPath, Root: opts.Root}
	_ = srv.ListenAndServe(&opts)
}
