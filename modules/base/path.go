package base

import (
	"errors"
	"os"
	"path/filepath"
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
	if len(out) <= len(parent) || out[len(parent)] != os.PathSeparator || !strings.HasPrefix(out,parent) {
		return "", ErrDangerousPathAccessDenied
	}
	return cleanedPath, nil
}

func JoinSanitizePathSlow(parent string, elem ...string) (string, error) {
	parent = filepath.Clean(parent)
	return JoinSanitizePath(parent, elem...)
}
