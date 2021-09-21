package pool

import (
	"net/url"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

type connectionType int

const (
	invalidConnection connectionType = iota
	tcpConnection
	tlsConnection
	unixConnection
)

// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
// scheme:opaque[?query][#fragment]
func extractAddressFromURL(rawAddress string) (connectionType, string) {
	if !strings.Contains(rawAddress, "/") {
		rawAddress = "tcp://" + rawAddress
	}
	u, err := url.Parse(rawAddress)
	if err != nil {
		return invalidConnection, ""
	}

	switch u.Scheme {
	case "tls":
		return tlsConnection, u.Host
	case "unix":
		return unixConnection, strings.TrimPrefix(rawAddress, "unix:")
	case "tcp":
		return tcpConnection, u.Host
	default:
		return invalidConnection, ""
	}
}

type ClientConnPool struct {
	mtx   sync.Mutex
	conns map[string]*ClientConn
}

func NewClientConnPool() *ClientConnPool {
	return &ClientConnPool{conns: make(map[string]*ClientConn)}
}

func (p *ClientConnPool) Dial(target string, opts ...grpc.DialOption) (*ClientConn, error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	return nil, nil
}
