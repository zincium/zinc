package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/zincium/zinc/modules/base"
	"github.com/zincium/zinc/modules/env"
	"github.com/zincium/zinc/modules/process"
	"github.com/zincium/zinc/modules/server"
)

// define
var (
	null                    = []byte("\x00")
	ErrMalformedNetworkData = errors.New("Malformed network data")
)

// Server server
type Server struct {
	Root    string // repositories root
	GitPath string // git or /path/to/git
	srv     *server.Server
	tlssrv  *server.Server
	qsrv    *server.QuicServer
}

// Request request
type Request struct {
	Path    string
	Host    string
	Service string
	Version string
	Query   url.Values
}

// ListenAndServe todo
func (srv *Server) ListenAndServe(opts *Options) error {
	// srv.GitPath = opts.GitPath
	// srv.Root = opts.Root
	var wg sync.WaitGroup
	wg.Add(1)
	srv.srv = &server.Server{
		Handler:     srv.Handle,
		MaxTimeout:  opts.maxTimeout,
		IdleTimeout: opts.idleTimeout,
	}
	go func() {
		defer wg.Done()
		if err := srv.srv.ListenAndServe(opts.Listen); err != nil {

		}
	}()
	if opts.TLSListen != "" {
		wg.Add(1)
		srv.tlssrv = &server.Server{
			Handler:     srv.Handle,
			MaxTimeout:  opts.maxTimeout,
			IdleTimeout: opts.idleTimeout,
		}
		go func() {
			defer wg.Done()
			if err := srv.tlssrv.ListenAndServeTLS(opts.Listen, opts.Certificate, opts.CertificateKey); err != nil {

			}
		}()
	}
	if opts.QUICListen != "" {
		wg.Add(1)
		srv.qsrv = &server.QuicServer{
			Handler:     srv.Handle,
			MaxTimeout:  opts.maxTimeout,
			IdleTimeout: opts.idleTimeout,
		}
		go func() {
			defer wg.Done()
			if err := srv.qsrv.ListenAndServeQUIC(opts.Listen, opts.Certificate, opts.CertificateKey); err != nil {

			}
		}()
	}
	wg.Wait()
	return nil
}

func (srv *Server) readRequest(conn net.Conn) (*Request, error) {
	dec := pktline.NewScanner(conn)
	if !dec.Scan() {
		return nil, dec.Err()
	}
	parts := bytes.Split(dec.Bytes(), null)
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
	// path start with '/'
	u, err := url.Parse("git://" + req.Host + string(p0[pos+1:]))
	base.DbgPrint("URL: %s", u)
	if err != nil {
		return nil, err
	}
	req.Path = filepath.Join(srv.Root, u.Path)
	if _, err := os.Stat(req.Path); err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("repository '%s' not found", u.Path)
	}
	req.Query = u.Query()
	if len(parts) > 4 {
		req.Version = string(parts[3])
	}
	return req, nil
}

// Handle on handle
func (srv *Server) Handle(conn net.Conn) {
	enc := pktline.NewEncoder(conn)
	req, err := srv.readRequest(conn)
	if err != nil {
		enc.EncodeString(err.Error())
		return
	}
	sugar.Infof("git-%s %s", req.Service, req.Path)
	cmd := exec.Command(srv.GitPath,
		"-c", "receive.denyDeleteCurrent=false",
		req.Service,
		req.Path)
	cmd.Env = append(cmd.Env, env.Environ()...)
	if len(req.Version) != 0 {
		cmd.Env = append(cmd.Env, "GIT_PROTOCOL="+req.Version)
	}
	base.DbgPrint("cmd: %v\n%v", cmd.Args, req)
	in, err := cmd.StdinPipe()
	if err != nil {
		base.DbgPrint("unable create stdin: %v", err)
		enc.EncodeString("internal server error")
		return
	}
	defer in.Close()
	out, err := cmd.StdoutPipe()
	if err != nil {
		base.DbgPrint("unable create stdout: %v", err)
		enc.EncodeString("internal server error")
		return
	}
	defer out.Close()
	if err := cmd.Start(); err != nil {
		// recored error
		base.DbgPrint("unable create process: %v", err)
		enc.EncodeString("internal server error")
		return
	}
	defer func() {
		_ = process.Finalize(cmd)
	}()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if _, err := base.Copy(conn, out); err != nil {
			base.DbgPrint("unable create stdin: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := base.Copy(in, conn); err != nil {
			base.DbgPrint("unable create stdin: %v", err)
		}
	}()
	wg.Wait()
}
