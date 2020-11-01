package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/lucas-clemente/quic-go"
	"github.com/zincium/zinc/modules/base"
)

// quic cat
func main() {
	if len(os.Args) < 2 {
		name := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "usage: %s host port\n       %s host:port", name, name)
		os.Exit(1)
	}
	base.IsDebugMode = base.IsTrue("ZINC_DEBUG")
	var address string
	if len(os.Args) == 2 {
		address = os.Args[1]
	} else {
		address = net.JoinHostPort(os.Args[1], os.Args[2])
	}
	base.DbgPrint("Use quicat to connet: %s", address)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	session, err := quic.DialAddr(address, tlsconfig, nil)
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
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		io.Copy(stream, os.Stdin)
		wg.Done()
	}()
	go func() {
		io.Copy(os.Stdout, stream)
		wg.Done()
	}()
	wg.Wait()
}
