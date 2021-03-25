package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/zincium/zinc/modules/base"
	"github.com/zincium/zinc/modules/env"
	"github.com/zincium/zinc/modules/process"
	"github.com/zincium/zinc/modules/server"
)

// define
var (
	null                    = []byte("\x00")
	ErrMalformedNetworkData = errors.New("malformed network data")
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
	return base.GroupExecute(
		func() error {
			// shutdown checked srv == nil
			if err := srv.srv.Shutdown(ctx); err != nil {
				sugar.Errorf("shutdown git over tcp service: %v", err)
				return err
			}
			return nil
		},
		func() error {
			if err := srv.tlssrv.Shutdown(ctx); err != nil {
				sugar.Errorf("shutdown git over tls service: %v", err)
				return err
			}
			return nil
		},
		func() error {
			if err := srv.qsrv.Shutdown(ctx); err != nil {
				sugar.Errorf("shutdown git over tcp service: %v", err)
				return err
			}
			return nil
		},
	)
}

func (srv *Server) listenAndServeTCP(opts *Options) error {
	if opts.Listen == "" {
		return nil
	}
	srv.srv = &server.Server{Handler: srv.Handle, MaxTimeout: opts.maxTimeout, IdleTimeout: opts.idleTimeout}
	sugar.Infof("listen %s (tcp)", opts.Listen)
	if err := srv.tlssrv.ListenAndServe(opts.Listen); err != nil && err != server.ErrServerClosed {
		sugar.Fatalf("ListenAndServe tls://%s error: %v", opts.TLSListen, err)
	}
	return nil
}

func (srv *Server) listenAndServeTLS(opts *Options) error {
	if opts.TLSListen == "" {
		return nil
	}
	srv.tlssrv = &server.Server{Handler: srv.Handle, MaxTimeout: opts.maxTimeout, IdleTimeout: opts.idleTimeout}
	sugar.Infof("listen %s (tls)", opts.TLSListen)
	if err := srv.tlssrv.ListenAndServeTLS(opts.TLSListen, opts.Certificate, opts.CertificateKey); err != nil && err != server.ErrServerClosed {
		sugar.Fatalf("ListenAndServe tls://%s error: %v", opts.TLSListen, err)
	}
	return nil
}

func (srv *Server) listenAndServeQUIC(opts *Options) error {
	if opts.QUICListen == "" {
		return nil
	}
	srv.qsrv = &server.QuicServer{Handler: srv.Handle, MaxTimeout: opts.maxTimeout, IdleTimeout: opts.idleTimeout}
	sugar.Infof("listen %s (quic)", opts.QUICListen)
	if err := srv.qsrv.ListenAndServeQUIC(opts.QUICListen, opts.Certificate, opts.CertificateKey); err != nil && err != server.ErrServerClosed {
		sugar.Fatalf("ListenAndServe quic://%s error: %v", opts.QUICListen, err)
	}
	return nil
}

// ListenAndServe todo
func (srv *Server) ListenAndServe(opts *Options) {
	_ = base.GroupExecute(
		func() error {
			return srv.listenAndServeTCP(opts)
		},
		func() error {
			return srv.listenAndServeTLS(opts)
		},
		func() error {
			return srv.listenAndServeQUIC(opts)
		},
	)
	sugar.Infof("zinc-secured git (tcp/tls/quic) server exited")
}

func (srv *Server) readRequest(conn net.Conn) (*Request, error) {
	dec := pktline.NewScanner(conn)
	_ = dec.Scan()
	if err := dec.Err(); err != nil {
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

// WriteError write error
func WriteError(conn net.Conn, msg string) {
	enc := pktline.NewEncoder(conn)
	_ = enc.EncodeString(msg + "\n")
	_ = enc.Flush()
}

// Handle on handle
func (srv *Server) Handle(conn net.Conn, mode string) {
	defer conn.Close()
	// FIXME Windows needs Fix: server closes net.Conn correctly
	req, err := srv.readRequest(conn)
	if err != nil {
		sugar.Warnf("[%s] read %v request error: %v", mode, conn.RemoteAddr(), err)
		WriteError(conn, err.Error())
		return
	}
	fmt.Fprintf(os.Stderr, "%v- %v", mode, req)
	sugar.Infof("[%s] git-%s %s", mode, req.Service, req.Path)
	cmd := exec.Command(srv.GitPath,
		"-c", "receive.denyDeleteCurrent=false",
		req.Service,
		req.Path)
	cmd.Env = append(cmd.Env, env.Environ()...)
	if len(req.Version) != 0 {
		cmd.Env = append(cmd.Env, "GIT_PROTOCOL="+req.Version)
	}
	base.DbgPrint("cmd: %v\n%v", cmd.Args, req)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		sugar.Errorf("create stdin pipe: %v", err)
		WriteError(conn, "internal server error")
		return
	}
	defer stdin.Close()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		sugar.Errorf("create stdout pipe: %v", err)
		WriteError(conn, "internal server error")
		return
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		// recored error
		sugar.Errorf("unable create process: %v", err)
		WriteError(conn, "internal server error")
		return
	}
	defer func() {
		_ = process.Finalize(cmd)
	}()
	err = base.GroupExecute(
		func() error {
			_, err := base.Copy(conn, stdout)
			stdout.Close()
			return err
		},
		func() error {
			_, err := base.Copy(stdin, conn)
			stdin.Close()
			return err
		},
	)
	if err != nil && err != io.EOF {
		sugar.Debugf("IO Exchange: %v", err)
	}
}
