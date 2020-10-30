package main

import (
	"net"
)

// Server server
type Server struct {
	listeners map[net.Listener]struct{}
}
