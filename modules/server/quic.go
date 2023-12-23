package server

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

// QuicServer quic
type QuicServer struct {
	Handler       func(net.Conn, string)
	IdleTimeout   time.Duration
	MaxTimeout    time.Duration
	MaxConnetions int //WIP
	mu            sync.RWMutex
	listenerWg    sync.WaitGroup
	listeners     map[*quic.Listener]struct{}
	conns         map[quic.Stream]struct{}
	connWg        sync.WaitGroup
	doneChan      chan struct{}
}

func (srv *QuicServer) trackConn(c quic.Stream, add bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.conns == nil {
		srv.conns = make(map[quic.Stream]struct{})
	}
	if add {
		srv.conns[c] = struct{}{}
		srv.connWg.Add(1)
	} else {
		delete(srv.conns, c)
		srv.connWg.Done()
	}
}

func (srv *QuicServer) trackListener(ln *quic.Listener, add bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.listeners == nil {
		srv.listeners = make(map[*quic.Listener]struct{})
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
func (srv *QuicServer) getDoneChanLocked() chan struct{} {
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	return srv.doneChan
}

func (srv *QuicServer) getDoneChan() <-chan struct{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	return srv.getDoneChanLocked()
}

func (srv *QuicServer) closeDoneChanLocked() {
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

func (srv *QuicServer) closeListenersLocked() error {
	var err error
	for ln := range srv.listeners {
		if cerr := ln.Close(); cerr != nil && err == nil {
			err = cerr
		}
		delete(srv.listeners, ln)
	}
	return err
}

// Serve serve
func (srv *QuicServer) Serve(ln *quic.Listener) error {
	srv.trackListener(ln, true)
	defer srv.trackListener(ln, false)
	var tempDelay time.Duration
	for {
		conn, e := ln.Accept(context.Background())
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

// ListenAndServeQUIC listen and serve
func (srv *QuicServer) ListenAndServeQUIC(listen string, certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	config := &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-git"},
	}
	ln, err := quic.ListenAddr(listen, config, nil)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

// Handle handle
func (srv *QuicServer) handle(conn quic.Connection) error {
	sm, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	srv.trackConn(sm, true)
	defer srv.trackConn(sm, false)
	qc := &quicConn{Stream: sm, conn: conn, idleTimeout: srv.IdleTimeout}
	if srv.MaxTimeout != 0 {
		qc.maxDeadline = time.Now().Add(srv.MaxTimeout)
	}
	if srv.Handler != nil {
		srv.Handler(qc, "QUIC")
	}
	return nil
}

// Shutdown todo
func (srv *QuicServer) Shutdown(ctx context.Context) error {
	if srv == nil {
		return nil
	}
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
