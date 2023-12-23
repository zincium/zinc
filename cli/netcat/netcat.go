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

	"github.com/quic-go/quic-go"
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
		fmt.Fprintf(os.Stderr, "[TLS] unable dial %s : %v\n", address, err)
		os.Exit(1)
	}
	defer conn.Close()
	_ = base.GroupExecute(
		func() error {
			if _, err := base.Copy(conn, os.Stdin); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "[tls] copy stdin to conn %v\n", err)
			}
			return nil
		},
		func() error {
			if _, err := base.Copy(os.Stdout, conn); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "[TLS] copy conn to stdout %v\n", err)
			}
			return nil
		},
	)

}

// ExchangeQUIC exchange QUIC
func ExchangeQUIC(address string) {
	session, err := quic.DialAddr(context.Background(), address, generateTLSConfig(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[QUIC] unable dial %s : %v\n", address, err)
		os.Exit(1)
	}
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "[QUIC] unable open stream %s : %v", address, err)
		os.Exit(1)
	}
	defer stream.Close()
	_ = base.GroupExecute(
		func() error {
			if _, err := base.Copy(stream, os.Stdin); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "[QUIC] copy stdin to conn %v\n", err)
			}
			return nil
		},
		func() error {
			if _, err := base.Copy(os.Stdout, stream); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "[QUIC] copy conn to stdout %v\n", err)
			}
			return nil
		},
	)
}

// ExchangeTCP exchange tls
func ExchangeTCP(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[TCP] unable dial %s : %v\n", address, err)
		os.Exit(1)
	}
	defer conn.Close()
	_ = base.GroupExecute(
		func() error {
			if _, err := base.Copy(conn, os.Stdin); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "[TCP] copy stdin to conn %v\n", err)
			}
			return nil
		},
		func() error {
			if _, err := base.Copy(os.Stdout, conn); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "[TCP] copy conn to stdout %v\n", err)
			}
			return nil
		},
	)
}

// version info
var (
	VERSION     = "1.0"
	BUILDTIME   string
	BUILDCOMMIT string
	BUILDBRANCH string
	GOVERSION   string
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
  -P|--port        Use a specific port to connect to the server

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

// Port port
func (e ExchangeMode) Port() string {
	switch e {
	case ModeTCP:
		return "9418"
	case ModeTLS:
		return "9419"
	case ModeQUIC:
		return "9420"
	}
	return "80"
}

type options struct {
	mode    ExchangeMode
	address string
	port    int
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
	case 'P':
		p, err := strconv.Atoi(oa)
		if err != nil {
			return base.ErrorCat("invaild port number: ", oa)
		}
		o.port = p
	}
	return nil
}

func (o *options) makeAddress(host string) string {
	a, _, err := net.SplitHostPort(host)
	if err != nil {
		if !strings.Contains(err.Error(), "missing port in address") {
			return host
		}
		if o.port != 0 {
			return net.JoinHostPort(a, strconv.Itoa(o.port))
		}
		return net.JoinHostPort(host, o.mode.Port())
	}
	if o.port != 0 {
		return net.JoinHostPort(a, strconv.Itoa(o.port))
	}
	return host
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
	pa.Add("port", base.NOARG, 'P')
	if err := pa.Execute(os.Args, o); err != nil {
		return err
	}
	args := pa.Unresolved()
	if len(args) == 0 {
		return errors.New("missing address input")
	}
	base.DbgPrint("Args %v", args)
	if base.IsTrue(os.Getenv("ZINC_INSECURE_TLS")) {
		InsecureSkipVerify = true
	}
	if len(args) == 1 || o.port != 0 {
		o.address = o.makeAddress(args[0]) // ignore port args
		return nil
	}
	o.address = net.JoinHostPort(args[0], args[1])
	return nil
}

func main() {
	var o options
	if err := o.ParseArgv(); err != nil {
		fmt.Fprintf(os.Stderr, "netcat parse argv error: \x1b[31m%v\x1b[0m\n", err)
		os.Exit(1)
	}
	base.DbgPrint("Use netcat to connect: %s insecure tls: %v", o.address, InsecureSkipVerify)
	switch o.mode {
	case ModeTCP:
		ExchangeTCP(o.address)
	case ModeTLS:
		ExchangeTLS(o.address)
	case ModeQUIC:
		ExchangeQUIC(o.address)
	}
}
