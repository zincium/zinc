package pool

import (
	"time"

	"google.golang.org/grpc"
)

type ClientConn struct {
	*grpc.ClientConn
	target   string
	deadline time.Time
	counter  int
	closed   bool // close conn
}

// grpc.ClientConnInterface
