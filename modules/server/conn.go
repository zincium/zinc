package server

import (
	"context"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
)

type serverConn struct {
	net.Conn

	idleTimeout   time.Duration
	maxDeadline   time.Time
	closeCanceler context.CancelFunc
}

func (c *serverConn) Write(p []byte) (n int, err error) {
	c.updateDeadline()
	n, err = c.Conn.Write(p)
	if _, isNetErr := err.(net.Error); isNetErr && c.closeCanceler != nil {
		c.closeCanceler()
	}
	return
}

func (c *serverConn) Read(b []byte) (n int, err error) {
	c.updateDeadline()
	n, err = c.Conn.Read(b)
	if _, isNetErr := err.(net.Error); isNetErr && c.closeCanceler != nil {
		c.closeCanceler()
	}
	return
}

func (c *serverConn) Close() (err error) {
	err = c.Conn.Close()
	if c.closeCanceler != nil {
		c.closeCanceler()
	}
	return
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

	idleTimeout   time.Duration
	maxDeadline   time.Time
	closeCanceler context.CancelFunc
}

func (c *quicConn) Write(p []byte) (n int, err error) {
	c.updateDeadline()
	n, err = c.Stream.Write(p)
	if _, isNetErr := err.(net.Error); isNetErr && c.closeCanceler != nil {
		c.closeCanceler()
	}
	return
}

func (c *quicConn) Read(b []byte) (n int, err error) {
	c.updateDeadline()
	n, err = c.Stream.Read(b)
	if _, isNetErr := err.(net.Error); isNetErr && c.closeCanceler != nil {
		c.closeCanceler()
	}
	return
}

func (c *quicConn) Close() (err error) {
	err = c.Stream.Close()
	if c.closeCanceler != nil {
		c.closeCanceler()
	}
	return
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
