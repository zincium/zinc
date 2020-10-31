package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/zincium/zinc/modules/server"
)

// define
var (
	null                    = []byte("\x00")
	ErrMalformedNetworkData = errors.New("Malformed network data")
)

// Server server
type Server struct {
	Root        string // repositories root
	GitPath     string // git or /path/to/git
	SecureToken string // secure token ?token=xxx&expired=yyyy
	srv         *server.Server
	qsrv        *server.QuicServer
}

// Request request
type Request struct {
	Path    string
	Host    string
	Service string
	Version string
	Query   url.Values
}

//ResolveRequest resolve request
func ResolveRequest(payload []byte) (*Request, error) {
	parts := bytes.Split(payload, null)
	if len(parts) < 2 {
		return nil, ErrMalformedNetworkData
	}
	p0 := parts[0]
	pos := bytes.IndexByte(p0, ' ')
	if pos == -1 {
		return nil, ErrMalformedNetworkData
	}
	service := string(p0[0:pos])
	if service != "git-upload-pack" && service != "git-receive-pack" && service != "git-upload-archive" {
		return nil, fmt.Errorf("unsupport git service %s", service)
	}
	req := &Request{
		Service: strings.TrimPrefix(service, "git-"),
		Host:    string(parts[1]),
	}
	u, err := url.Parse("git://" + req.Host + string(p0[pos+1]))
	if err != nil {
		return nil, err
	}
	req.Path = u.Path
	req.Query = u.Query()
	if len(parts) > 4 {
		req.Version = string(parts[3])
	}
	return req, nil
}

// Handle on handle
func (srv *Server) Handle(conn net.Conn) {
	enc := pktline.NewEncoder(conn)
	dec := pktline.NewScanner(conn)
	if !dec.Scan() {
		enc.Encodef("Protocol error: %v", dec.Err())
		return
	}
	pkl := dec.Bytes()
	fmt.Fprintf(os.Stderr, "%v", pkl)
}
