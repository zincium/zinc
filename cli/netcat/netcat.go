package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/lucas-clemente/quic-go"
	"github.com/zincium/zinc/modules/base"
)

// define
var (
	InsecureSkipVerify = false
)

func generateTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: InsecureSkipVerify,
		NextProtos:         []string{"quic-git"},
	}
}

// ExchangeTLS exchange tls
func ExchangeTLS(address string) {
	conn, err := tls.Dial("tcp", address, generateTLSConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable dial(tls) %s : %v", address, err)
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

// ExchangeTCP exchange tls
func ExchangeTCP(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable dial(tcp) %s : %v", address, err)
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

// version info
var (
	VERSION       = "1.0"
	BUILDTIME     string
	BUILDCOMMIT   string
	BUILDBRANCH   string
	GOVERSION     string
	ServerVersion = "Zinc/" + VERSION
)

func version() {
	fmt.Fprintf(os.Stdout, `netcat - A witty netcat client
version:       %s
build branch:  %s
build commit:  %s
build time:    %s
go version:    %s
`, VERSION, BUILDBRANCH, BUILDCOMMIT, BUILDTIME, GOVERSION)

}

func usage() {
	fmt.Fprintf(os.Stdout, `netcat - A witty netcat client
usage: %s <option> url
  -h|--help        Show usage text and quit
  -v|--version     Show version number and quit
  -V|--verbose     Make the operation more talkative
  -T|--tcp         Create a tcp connection according to the input address
  -S|--tls         Create a tls connection according to the input address
  -Q|--quic        Create a quic connection according to the input address
  -K|--insecure    Allow insecure server connections when using TLS
  -R|--replace     Replace port to default port

`, os.Args[0])
}

// ExchangeMode todo
type ExchangeMode int

// define
const (
	ModeTCP ExchangeMode = iota
	ModeTLS
	ModeQUIC
)

type options struct {
	mode        ExchangeMode
	address     string
	replacePort bool
}

func (o *options) Invoke(val int, oa, raw string) error {
	switch val {
	case 'h':
		version()
		os.Exit(0)
	case 'v':
		usage()
		os.Exit(0)
	case 'V':
		base.IsDebugMode = true
	case 'T':
		o.mode = ModeTCP
	case 'S':
		o.mode = ModeTLS
	case 'Q':
		o.mode = ModeQUIC
	case 'K':
		InsecureSkipVerify = true
	case 'R':
		o.replacePort = true
	}
	return nil
}

func (o *options) buildDefaultAddress(addr string) string {
	switch o.mode {
	case ModeTCP:
		return addr + ":" + strconv.Itoa(9418)
	case ModeTLS:
		return addr + ":" + strconv.Itoa(9419)
	case ModeQUIC:
		return addr + ":" + strconv.Itoa(9420)
	}
	return addr
}

func (o *options) rebuildAddress(addr string) string {
	address, _, err := net.SplitHostPort(addr)
	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			return o.buildDefaultAddress(addr)
		}
		return addr
	}
	if o.replacePort {
		return o.buildDefaultAddress(address)
	}
	return addr
}

func (o *options) ParseArgv() error {
	var pa base.ParseArgs
	pa.Add("help", base.NOARG, 'h')
	pa.Add("version", base.NOARG, 'v')
	pa.Add("verbose", base.NOARG, 'V')
	pa.Add("tcp", base.NOARG, 'T')
	pa.Add("tls", base.NOARG, 'S')
	pa.Add("quic", base.NOARG, 'Q')
	pa.Add("insecure", base.NOARG, 'K')
	pa.Add("replace", base.NOARG, 'R')
	if err := pa.Execute(os.Args, o); err != nil {
		return err
	}
	args := pa.Unresolved()
	if len(args) == 0 {
		return errors.New("missing address input")
	}
	if base.IsTrue(os.Getenv("ZINC_INSECURE_TLS")) {
		InsecureSkipVerify = true
	}
	if len(args) == 1 {
		o.address = o.rebuildAddress(args[0])
	} else {
		if o.replacePort {
			o.address = o.buildDefaultAddress(args[0])
		} else {
			o.address = net.JoinHostPort(args[0], args[1])
		}
	}
	return nil
}

func main() {
	var o options
	if err := o.ParseArgv(); err != nil {
		fmt.Fprintf(os.Stderr, "ParseArgv: \x1b[31m%v\x1b[0m\n", err)
		os.Exit(1)
	}
	base.DbgPrint("Use netcat to connect: %s", o.address)
	switch o.mode {
	case ModeTCP:
		ExchangeTCP(o.address)
	case ModeTLS:
		ExchangeTLS(o.address)
	case ModeQUIC:
		ExchangeQUIC(o.address)
	}
}
