package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/zincium/zinc/modules/base"
)

// ExchangeTCP exchange tls
func ExchangeTCP(address string) {
	conn, err := net.Dial("tcp", address)
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

func main() {
	// --tls --quic --tcp --insecure
}
