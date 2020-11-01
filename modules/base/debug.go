package base

import (
	"fmt"
	"os"
	"strings"
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

// IsTrue fun
func IsTrue(s string) bool {
	lower := strings.ToLower(s)
	return lower == "true" || lower == "yes" || lower == "on" || lower == "1"
}
