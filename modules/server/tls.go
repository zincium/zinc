package server

import (
	"crypto/tls"
	"net"
	"time"
)

// ServeTLS serve
func (srv *Server) ServeTLS(ln net.Listener, tlsconfig *tls.Config) error {
	srv.trackListener(ln, true)
	defer srv.trackListener(ln, false)
	var tempDelay time.Duration
	for {
		conn, e := ln.Accept()
		if e != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		go srv.handleTLS(conn, tlsconfig)
	}
}

// ListenAndServeTLS listen and serve
func (srv *Server) ListenAndServeTLS(listen string, certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	//https://gist.github.com/denji/12b3a568f092ab951456
	config := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	return srv.ServeTLS(ln, config)
}

// Handle handle
func (srv *Server) handleTLS(conn net.Conn, config *tls.Config) {
	tlsconn := tls.Server(conn, config)
	if err := tlsconn.Handshake(); err != nil {
		return
	}
	if srv.Handler != nil {
		srv.Handler(tlsconn)
	}
}
