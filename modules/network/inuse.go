package network

import (
	"errors"
	"runtime"
	"syscall"
)

func ErrorIsAddressInUse(err error) bool {
	var netErrno syscall.Errno
	if ok := errors.As(err, &netErrno); !ok {
		return false
	}
	if netErrno == syscall.EADDRINUSE {
		return true
	}
	const WSAEADDRINUSE = 10048
	return runtime.GOOS == "windows" && netErrno == WSAEADDRINUSE
}
