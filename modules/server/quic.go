package server

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
)

// QuicServer quic
type QuicServer struct {
	Handler    func(io.ReadWriter)
	mu         sync.RWMutex
	listenerWg sync.WaitGroup
	listeners  map[quic.Listener]struct{}
	conns      map[quic.Stream]struct{}
	connWg     sync.WaitGroup
	doneChan   chan struct{}
}

type quicConn struct {
	quic.Stream
	localAddr  net.Addr
	remoteAddr net.Addr
}

func (qc *quicConn) Close() error {

	return nil
}

func (qc *quicConn) LocalAddr() net.Addr {
	return qc.localAddr
}

func (qc *quicConn) RemoteAddr() net.Addr {
	return qc.remoteAddr
}

func (qc *quicConn) SetDeadline(t time.Time) error {
	return nil
}

func (qc *quicConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (qc *quicConn) SetWriteDeadline(t time.Time) error {
	return nil
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

func (srv *QuicServer) trackListener(ln quic.Listener, add bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.listeners == nil {
		srv.listeners = make(map[quic.Listener]struct{})
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

// Serve serve
func (srv *QuicServer) Serve(ln quic.Listener) error {
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
func (srv *QuicServer) handle(conn quic.Session) error {
	sm, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	return srv.onStream(sm)
}

func (srv *QuicServer) onStream(sm quic.Stream) error {
	if srv.Handler != nil {
		srv.Handler(sm)
	}
	return nil
}
