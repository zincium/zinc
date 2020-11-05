package server

import (
	"errors"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
)

var (
	errClosed   = errors.New("tcp: use of closed connection")
	errShutdown = errors.New("tcp: protocol is shutdown")
)

type serverConn struct {
	net.Conn

	idleTimeout time.Duration
	maxDeadline time.Time
}

func (c *serverConn) Write(p []byte) (int, error) {
	c.updateDeadline()
	return c.Conn.Write(p)
}

func (c *serverConn) Read(b []byte) (int, error) {
	c.updateDeadline()
	return c.Conn.Read(b)
}

func (c *serverConn) Close() error {
	//time.Sleep(time.Millisecond * 20)
	return c.Conn.Close()
}

func (c *serverConn) updateDeadline() {
	switch {
	case c.idleTimeout > 0:
		idleDeadline := time.Now().Add(c.idleTimeout)
		if idleDeadline.Unix() < c.maxDeadline.Unix() || c.maxDeadline.IsZero() {
			c.Conn.SetDeadline(idleDeadline)
			return
		}
		fallthrough
	default:
		c.Conn.SetDeadline(c.maxDeadline)
	}
}

type quicConn struct {
	quic.Stream
	conn quic.Session

	idleTimeout time.Duration
	maxDeadline time.Time
}

func (c *quicConn) Write(p []byte) (int, error) {
	c.updateDeadline()
	return c.Stream.Write(p)
}

func (c *quicConn) Read(b []byte) (int, error) {
	c.updateDeadline()
	return c.Stream.Read(b)
}

func (c *quicConn) updateDeadline() {
	switch {
	case c.idleTimeout > 0:
		idleDeadline := time.Now().Add(c.idleTimeout)
		if idleDeadline.Unix() < c.maxDeadline.Unix() || c.maxDeadline.IsZero() {
			c.Stream.SetDeadline(idleDeadline)
			return
		}
		fallthrough
	default:
		c.Stream.SetDeadline(c.maxDeadline)
	}
}

// LocalAddr locla addr
func (c *quicConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *quicConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
