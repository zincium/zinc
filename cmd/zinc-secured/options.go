package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/zincium/zinc/modules/base"
	"github.com/zincium/zinc/modules/cast"
	"github.com/zincium/zinc/modules/env"
)

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
	fmt.Fprintf(os.Stdout, `zinc-secured - Redesigned Git Transfer Service
version:       %s
build branch:  %s
build commit:  %s
build time:    %s
go version:    %s
`, VERSION, BUILDBRANCH, BUILDCOMMIT, BUILDTIME, GOVERSION)

}

func usage() {
	fmt.Fprintf(os.Stdout, `zinc-secured - Redesigned Git Transfer Service
usage: %s <option> url
  -h|--help        Show usage text and quit
  -v|--version     Show version number and quit
  -V|--verbose     Make the operation more talkative
  -p|--profile     Set profile path. default: %s\config\zinc-secured.toml
  -D|--daemon      Run zinc-secured as daemon   

`, os.Args[0], env.AppDir())
}

// Options base
type Options struct {
	Listen         string `toml:"Listen"`
	Root           string `toml:"Root"`
	IdleTimeout    string `toml:"IdleTimeout,omitempty"`
	idleTimeout    time.Duration
	MaxTimeout     string `toml:"MaxTimeout,omitempty"`
	maxTimeout     time.Duration
	TLSListen      string `toml:"TlsListen,omitempty"`
	QUICListen     string `toml:"QuicListen,omitempty"`
	Certificate    string `toml:"Certificate,omitempty"`
	CertificateKey string `toml:"CertificateKey,omitempty"`
	profile        string
	background     bool
}

// Initialize initialize opts
func (opts *Options) Initialize(expander *env.Expander) error {
	fd, err := os.Open(opts.profile)
	if err != nil {
		return base.ErrorCat("unable open configure file: ", err.Error())
	}
	defer fd.Close()
	if toml.NewDecoder(fd).Decode(opts); err != nil {
		return base.ErrorCat("unable decode configure: ", err.Error())
	}
	opts.Certificate = expander.PathExpand(opts.Certificate)
	opts.CertificateKey = expander.PathExpand(opts.CertificateKey)
	opts.idleTimeout = cast.ToDuration(opts.IdleTimeout)
	opts.maxTimeout = cast.ToDuration(opts.maxTimeout)
	return nil
}

// Invoke invoke
func (opts *Options) Invoke(val int, oa, raw string) error {
	switch val {
	case 'h':
		version()
		os.Exit(0)
	case 'v':
		usage()
		os.Exit(0)
	case 'V':
		base.IsDebugMode = true
	case 'p':
		opts.profile = oa
	case 'D':
		opts.background = true
	}
	return nil
}

// ParseArgv parse argv
func (opts *Options) ParseArgv() error {
	var pa base.ParseArgs
	pa.Add("help", base.NOARG, 'h')
	pa.Add("version", base.NOARG, 'v')
	pa.Add("verbose", base.NOARG, 'V')
	pa.Add("profile", base.REQUIRED, 'p')
	pa.Add("daemon", base.NOARG, 'D')
	if err := pa.Execute(os.Args, opts); err != nil {
		return err
	}
	expander := env.NewExpander()
	if len(opts.profile) == 0 {
		opts.profile = expander.PathExpand("${APPDIR}/config/zinc-secured.toml")
	}
	if err := opts.Initialize(expander); err != nil {
		return err
	}
	if opts.background {

	}
	return nil
}
