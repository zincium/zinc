package shadow

import (
	"net"
	"sync"
)

// Listener close once
type Listener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

// Close close once
func (ln *Listener) Close() error {
	ln.once.Do(ln.close)
	return ln.closeErr
}

func (ln *Listener) close() { ln.closeErr = ln.Listener.Close() }
