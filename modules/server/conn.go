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

// rstAvoidanceDelay is the amount of time we sleep after closing the
// write side of a TCP connection before closing the entire socket.
// By sleeping, we increase the chances that the client sees our FIN
// and processes its final data before they process the subsequent RST
// from closing a connection with known unread data.
// This RST seems to occur mostly on BSD systems. (And Windows?)
// This timeout is somewhat arbitrary (~latency around the planet).
const rstAvoidanceDelay = 500 * time.Millisecond

type closeWriter interface {
	CloseWrite() error
}

var _ closeWriter = (*net.TCPConn)(nil)

// closeWrite flushes any outstanding data and sends a FIN packet (if
// client is connected via TCP), signalling that we're done. We then
// pause for a bit, hoping the client processes it before any
// subsequent RST.
//
// See https://golang.org/issue/3595
func (c *serverConn) closeWriteAndWait() {
	if tcp, ok := c.Conn.(closeWriter); ok {
		tcp.CloseWrite()
	}
	time.Sleep(rstAvoidanceDelay)
}

func (c *serverConn) Close() error {
	c.closeWriteAndWait()
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
