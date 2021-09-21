package pool

import (
	"fmt"
	"os"
	"testing"
)

func TestExtractAddressFromURL(t *testing.T) {
	sv := []string{
		"127.0.0.1:3386",
		"tls://8.8.8.8:8080",
		"tcp://9.9.9.9:10240",
		"unix:/path/to/unix.sock",
	}
	for _, s := range sv {
		ct, addr := extractAddressFromURL(s)
		fmt.Fprintf(os.Stderr, "%s --> %d %s\n", s, ct, addr)
	}
}
