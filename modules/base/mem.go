package base

import (
	"io"
	"sync"
)

// MaxPacketSize max buffer size
const MaxPacketSize = 32 * 1024

var (
	ioCopyPool = sync.Pool{New: func() interface{} {
		b := make([]byte, MaxPacketSize)
		return &b
	}}
)

// Copy copy reader to writer
func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := ioCopyPool.Get().(*[]byte)
	defer ioCopyPool.Put(buf)
	return io.CopyBuffer(dst, src, *buf)
}
