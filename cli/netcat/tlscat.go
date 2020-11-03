package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"

	"github.com/zincium/zinc/modules/base"
)

// define
var (
	InsecureSkipVerify = false
)

func generateTLSConfig() *tls.Config {
	insecure := base.IsTrue(os.Getenv("ZINC_INSECURE_TLS"))
	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: insecure,
	}
}

// ExchangeTLS exchange tls
func ExchangeTLS(address string) {
	conn, err := tls.Dial("tcp", address, generateTLSConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable dial(quic) %s : %v", address, err)
		os.Exit(1)
	}
	defer conn.Close()
	_ = base.GroupExecute(
		func() error {
			if _, err := base.Copy(conn, os.Stdin); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "copy stdin to conn %v\n", err)
			}
			return nil
		},
		func() error {
			if _, err := base.Copy(os.Stdout, conn); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "copy conn to stdout %v\n", err)
			}
			return nil
		},
	)

}
