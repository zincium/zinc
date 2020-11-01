package base

import (
	"fmt"
	"os"
)

// defined
var (
	IsDebugMode bool
)

// DbgPrint todo
func DbgPrint(format string, a ...interface{}) {
	if IsDebugMode {
		ss := fmt.Sprintf(format, a...)
		_, _ = os.Stderr.Write(BufferCat("\x1b[33m* ", ss, "\x1b[0m\n"))
	}
}
