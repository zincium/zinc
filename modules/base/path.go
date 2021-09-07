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

func isChildPathFast(path, parent string) bool {
	if len(parent) >= len(path) {
		return false
	}
	return strings.HasPrefix(path, parent) && path[len(parent)] == os.PathSeparator
}

func JoinSanitizePath(parent string, elem ...string) (string, error) {
	var buf strings.Builder
	_, _ = buf.WriteString(parent)
	for _, e := range elem {
		_ = buf.WriteByte(os.PathSeparator)
		_, _ = buf.WriteString(e)
	}
	cleanedPath := filepath.Clean(buf.String())
	if !isChildPathFast(cleanedPath, parent) {
		return "", ErrDangerousPathAccessDenied
	}
	return cleanedPath, nil
}

func JoinSanitizePathSlow(parent string, elem ...string) (string, error) {
	parent = filepath.Clean(parent)
	return JoinSanitizePath(parent, elem...)
}
