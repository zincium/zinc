package base

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	ErrDangerousPathAccessDenied = errors.New("dangerous path access denied")
)

func JoinSanitizePath(parent string, elem ...string) (string, error) {
	var buf strings.Builder
	_, _ = buf.WriteString(parent)
	for _, e := range elem {
		_ = buf.WriteByte(os.PathSeparator)
		_, _ = buf.WriteString(e)
	}
	out := filepath.Clean(buf.String())
	if len(out) <= len(parent) {
		return "", ErrDangerousPathAccessDenied
	}
	if strings.HasPrefix(out, parent) && os.IsPathSeparator(out[len(parent)]) {
		return out, nil
	}
	if runtime.GOOS != "windows" && parent == "/" {
		return out, nil
	}
	return "", ErrDangerousPathAccessDenied
}

func JoinSanitizePathSlow(parent string, elem ...string) (string, error) {
	parent = filepath.Clean(parent)
	return JoinSanitizePath(parent, elem...)
}
