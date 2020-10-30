package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
)

// ServerQUIC quic
type ServerQUIC struct {
	mu         sync.RWMutex
	listenerWg sync.WaitGroup
	listeners  map[quic.Listener]struct{}
	conns      map[quic.Stream]struct{}
	connWg     sync.WaitGroup
	doneChan   chan struct{}
	so         *SecureOptions
}

func (srv *ServerQUIC) trackConn(c quic.Stream, add bool) {
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

func (srv *ServerQUIC) trackListener(ln quic.Listener, add bool) {
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
func (srv *ServerQUIC) getDoneChanLocked() chan struct{} {
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	return srv.doneChan
}

func (srv *ServerQUIC) getDoneChan() <-chan struct{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	return srv.getDoneChanLocked()
}

func (srv *ServerQUIC) generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}

// Serve serve
func (srv *ServerQUIC) Serve(ln quic.Listener) error {
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
		go srv.Handle(conn)
	}
}

// ListenAndServe listen and serve
func (srv *ServerQUIC) ListenAndServe(listen string) error {
	cert, err := tls.LoadX509KeyPair(srv.so.CertFile, srv.so.KeyFile)
	if err != nil {
		return err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	ln, err := quic.ListenAddr(listen, config, nil)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

// Handle handle
func (srv *ServerQUIC) Handle(conn quic.Session) error {
	sm, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	return srv.onStream(sm)
}

func (srv *ServerQUIC) onStream(sm quic.Stream) error {

	return nil
}
