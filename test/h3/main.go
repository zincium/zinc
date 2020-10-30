package main

import (
	"net/http"

	_ "net/http/pprof"

	"github.com/lucas-clemente/quic-go/http3"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("/path/dir")))
	http3.ListenAndServeQUIC("localhost:4242", "/path/to/cert/chain.pem", "/path/to/privkey.pem", nil)
}
