package main

import (
	"crypto/tls"
	"net"
	"sync"
	"time"
)

// SecureOptions options
type SecureOptions struct {
	CertFile string
	KeyFile  string
}

// ServerTLS tls server
type ServerTLS struct {
	mu         sync.RWMutex
	listenerWg sync.WaitGroup
	listeners  map[net.Listener]struct{}
	conns      map[net.Conn]struct{}
	connWg     sync.WaitGroup
	doneChan   chan struct{}
	so         *SecureOptions
}

func (srv *ServerTLS) trackConn(c net.Conn, add bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.conns == nil {
		srv.conns = make(map[net.Conn]struct{})
	}
	if add {
		srv.conns[c] = struct{}{}
		srv.connWg.Add(1)
	} else {
		delete(srv.conns, c)
		srv.connWg.Done()
	}
}

func (srv *ServerTLS) trackListener(ln net.Listener, add bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.listeners == nil {
		srv.listeners = make(map[net.Listener]struct{})
	}
	if add {
		// If the *Server is being reused after a previous
		// Close or Shutdown, reset its doneChan:
		if len(srv.listeners) == 0 && len(srv.conns) == 0 {
			srv.doneChan = nil
		}
		srv.listeners[ln] = struct{}{}
		srv.listenerWg.Add(1)
	} else {
		delete(srv.listeners, ln)
		srv.listenerWg.Done()
	}
}
func (srv *ServerTLS) getDoneChanLocked() chan struct{} {
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	return srv.doneChan
}

func (srv *ServerTLS) getDoneChan() <-chan struct{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	return srv.getDoneChanLocked()
}

// Serve serve
func (srv *ServerTLS) Serve(ln net.Listener) error {
	cert, err := tls.LoadX509KeyPair(srv.so.CertFile, srv.so.KeyFile)
	if err != nil {
		return err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
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
		go func() {
			c := tls.Server(conn, config)
			if err := c.Handshake(); err != nil {
				// LOG
			}
			srv.Handle(c)
		}()
	}
}

// ListenAndServe listen and serve
func (srv *ServerTLS) ListenAndServe(listen string) error {
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

// Handle handle
func (srv *ServerTLS) Handle(conn net.Conn) error {

	return nil
}
