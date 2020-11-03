package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/lucas-clemente/quic-go"
	"github.com/zincium/zinc/modules/base"
)

// ExchangeQUIC exchange QUIC
func ExchangeQUIC(address string) {
	session, err := quic.DialAddr(address, generateTLSConfig(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable dial(quic) %s : %v", address, err)
		os.Exit(1)
	}
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable open stream %s : %v", address, err)
		os.Exit(1)
	}
	defer stream.Close()
	_ = base.GroupExecute(
		func() error {
			if _, err := base.Copy(stream, os.Stdin); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "copy stdin to conn %v\n", err)
			}
			return nil
		},
		func() error {
			if _, err := base.Copy(os.Stdout, stream); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "copy conn to stdout %v\n", err)
			}
			return nil
		},
	)
}
