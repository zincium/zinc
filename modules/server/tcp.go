package server

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

// ErrServerClosed define
var ErrServerClosed = errors.New("git: Server closed")

// Server server
type Server struct {
	Handler     func(conn net.Conn)
	MaxTimeout  time.Duration
	IdleTimeout time.Duration
	mu          sync.RWMutex
	listenerWg  sync.WaitGroup
	listeners   map[net.Listener]struct{}
	conns       map[net.Conn]struct{}
	connWg      sync.WaitGroup
	doneChan    chan struct{}
}

func (srv *Server) trackConn(c net.Conn, add bool) {
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

func (srv *Server) trackListener(ln net.Listener, add bool) {
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
func (srv *Server) getDoneChanLocked() chan struct{} {
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	return srv.doneChan
}

func (srv *Server) getDoneChan() <-chan struct{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	return srv.getDoneChanLocked()
}

func (srv *Server) closeListenersLocked() error {
	var err error
	for ln := range srv.listeners {
		if cerr := ln.Close(); cerr != nil && err == nil {
			err = cerr
		}
		delete(srv.listeners, ln)
	}
	return err
}

func (srv *Server) closeDoneChanLocked() {
	ch := srv.getDoneChanLocked()
	select {
	case <-ch:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		// by srv.mu.
		close(ch)
	}
}

// Serve serve
func (srv *Server) Serve(ln net.Listener) error {
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
		go srv.handle(conn)
	}
}

func (srv *Server) handle(conn net.Conn) {
	srv.trackConn(conn, true)
	defer srv.trackConn(conn, false)
	sconn := &serverConn{Conn: conn, idleTimeout: srv.IdleTimeout}
	if srv.MaxTimeout != 0 {
		sconn.maxDeadline = time.Now().Add(srv.MaxTimeout)
	}
	if srv.Handler != nil {
		srv.Handler(sconn)
	}
}

// ListenAndServe listen and serve
func (srv *Server) ListenAndServe(listen string) error {
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

//Shutdown todo
func (srv *Server) Shutdown(ctx context.Context) error {
	srv.mu.Lock()
	lnerr := srv.closeListenersLocked()
	srv.closeDoneChanLocked()
	srv.mu.Unlock()

	finished := make(chan struct{}, 1)
	go func() {
		srv.listenerWg.Wait()
		srv.connWg.Wait()
		finished <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-finished:
		return lnerr
	}
}
