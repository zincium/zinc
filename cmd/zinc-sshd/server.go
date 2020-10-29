package main

import (
	"context"
	"time"

	"github.com/gliderlabs/ssh"
)

// ServerOptions server options
type ServerOptions struct {
	PrivateKeyFile  []string
	PrivateKeyBytes []byte
	IdleTimeout     time.Duration
	MaxTimeout      time.Duration
}

// Server server
type Server struct {
	s       *ssh.Server
	options *ServerOptions
}

// ListenAndServe todo
func (srv *Server) ListenAndServe(listen string) error {
	srv.s = &ssh.Server{
		Addr:             listen,
		Handler:          srv.Handle,
		Version:          ServerVersion,
		PublicKeyHandler: srv.PublicKeyAuth,
		IdleTimeout:      srv.options.IdleTimeout,
		MaxTimeout:       srv.options.MaxTimeout,
	}
	if len(srv.options.PrivateKeyBytes) == 0 {
		for _, file := range srv.options.PrivateKeyFile {
			srv.s.SetOption(ssh.HostKeyFile(file))
		}
	} else {
		srv.s.SetOption(ssh.HostKeyPEM(srv.options.PrivateKeyBytes))
	}
	return srv.s.ListenAndServe()
}

// Shutdown shutdown
func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.s.Shutdown(ctx)
}

// PublicKeyAuth publicKeyAuth
func (srv *Server) PublicKeyAuth(ctx ssh.Context, key ssh.PublicKey) bool {

	return true
}

// Handle todo
func (srv *Server) Handle(s ssh.Session) {

}
