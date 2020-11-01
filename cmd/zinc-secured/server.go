package main

import (
	"bytes"
	"context"
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

// Shutdown shutdown
func (srv *Server) Shutdown(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		if srv.srv == nil {
			return
		}
		if err := srv.srv.Shutdown(ctx); err != nil {
			sugar.Errorf("shutdown git over tcp service: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		if srv.tlssrv == nil {
			return
		}
		if err := srv.tlssrv.Shutdown(ctx); err != nil {
			sugar.Errorf("shutdown git over tls service: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		if srv.qsrv == nil {
			return
		}
		if err := srv.qsrv.Shutdown(ctx); err != nil {
			sugar.Errorf("shutdown git over quic service: %v", err)
		}
	}()
	return nil
}

// ListenAndServe todo
func (srv *Server) ListenAndServe(opts *Options) {
	var wg sync.WaitGroup
	srv.srv = &server.Server{Handler: srv.Handle, MaxTimeout: opts.maxTimeout, IdleTimeout: opts.idleTimeout}
	srv.tlssrv = &server.Server{Handler: srv.Handle, MaxTimeout: opts.maxTimeout, IdleTimeout: opts.idleTimeout}
	srv.qsrv = &server.QuicServer{Handler: srv.Handle, MaxTimeout: opts.maxTimeout, IdleTimeout: opts.idleTimeout}
	wg.Add(3)
	go func() {
		defer wg.Done()
		if opts.Listen == "" {
			return
		}
		sugar.Infof("listen %s (tcp)", opts.Listen)
		if err := srv.srv.ListenAndServe(opts.Listen); err != nil && err != server.ErrServerClosed {
			sugar.Fatalf("ListenAndServe tcp://%s error: %v", opts.Listen, err)
		}
	}()
	go func() {
		defer wg.Done()
		if opts.TLSListen == "" {
			return
		}
		sugar.Infof("listen %s (tls)", opts.TLSListen)
		if err := srv.tlssrv.ListenAndServeTLS(opts.TLSListen, opts.Certificate, opts.CertificateKey); err != nil && err != server.ErrServerClosed {
			sugar.Fatalf("ListenAndServe tls://%s error: %v", opts.TLSListen, err)
		}
	}()
	go func() {
		defer wg.Done()
		if opts.QUICListen == "" {
			return
		}
		sugar.Infof("listen %s (quic)", opts.QUICListen)
		if err := srv.qsrv.ListenAndServeQUIC(opts.QUICListen, opts.Certificate, opts.CertificateKey); err != nil && err != server.ErrServerClosed {
			sugar.Fatalf("ListenAndServe quic://%s error: %v", opts.QUICListen, err)
		}
	}()
	wg.Wait()
	sugar.Infof("zinc-secured git (tcp/tls/quic) server exited")
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
		return nil, fmt.Errorf("unsupported git service %s", service)
	}
	req := &Request{
		Service: strings.TrimPrefix(service, "git-"),
		Host:    string(parts[1]),
	}
	// path start with '/'
	u, err := url.Parse("git://" + req.Host + string(p0[pos+1:]))
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
		sugar.Warnf("read %v request error: %v", conn.RemoteAddr(), err)
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
		sugar.Errorf("create stdin pipe: %v", err)
		enc.EncodeString("internal server error")
		return
	}
	defer in.Close()
	out, err := cmd.StdoutPipe()
	if err != nil {
		sugar.Errorf("create stdout pipe: %v", err)
		enc.EncodeString("internal server error")
		return
	}
	defer out.Close()
	if err := cmd.Start(); err != nil {
		// recored error
		sugar.Errorf("unable create process: %v", err)
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
			sugar.Debugf("copy out to conn: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := base.Copy(in, conn); err != nil {
			sugar.Debugf("copy stdin to conn: %v", err)
		}
	}()
	wg.Wait()
}
