package main

import (
	"os"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/zincium/zinc/modules/base"
	"github.com/zincium/zinc/modules/cast"
	"github.com/zincium/zinc/modules/env"
)

// Options base
type Options struct {
	Listen         string `toml:"Listen"`
	IdleTimeout    string `toml:"IdleTimeout,omitempty"`
	idleTimeout    time.Duration
	MaxTimeout     string `toml:"MaxTimeout,omitempty"`
	maxTimeout     time.Duration
	TLSListen      string `toml:"TlsListen,omitempty"`
	QUICListen     string `toml:"QuicListen,omitempty"`
	Certificate    string `toml:"Certificate,omitempty"`
	CertificateKey string `toml:"CertificateKey,omitempty"`
}

// NewOptions new options
func NewOptions(file string) (*Options, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, base.ErrorCat("unable open configure file: ", err.Error())
	}
	defer fd.Close()
	var options Options
	if toml.NewDecoder(fd).Decode(&options); err != nil {
		return nil, base.ErrorCat("unable decode configure: ", err.Error())
	}
	expander := env.NewExpander()
	options.Certificate = expander.Expand(options.Certificate)
	options.CertificateKey = expander.Expand(options.CertificateKey)
	options.idleTimeout = cast.ToDuration(options.IdleTimeout)
	options.maxTimeout = cast.ToDuration(options.maxTimeout)
	return &options, nil
}
