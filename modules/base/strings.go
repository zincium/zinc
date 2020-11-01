package base

import (
	"bytes"
	"errors"
	"strings"
)

// StrSplitSkipEmpty skip empty string suggestcap is suggest cap
func StrSplitSkipEmpty(s string, sep byte, suggestcap int) []string {
	sv := make([]string, 0, suggestcap)
	var first, i int
	for ; i < len(s); i++ {
		if s[i] != sep {
			continue
		}
		if first != i {
			sv = append(sv, s[first:i])
		}
		first = i + 1
	}
	if first < len(s) {
		sv = append(sv, s[first:])
	}
	return sv
}

func splitPathInternal(p string) []string {
	sv := make([]string, 0, 8)
	var first, i int
	for ; i < len(p); i++ {
		if p[i] != '/' && p[i] != '\\' {
			continue
		}
		if first != i {
			sv = append(sv, p[first:i])
		}
		first = i + 1
	}
	if first < len(p) {
		sv = append(sv, p[first:])
	}
	return sv
}

// SplitPath skip empty string suggestcap is suggest cap
func SplitPath(p string) []string {
	svv := splitPathInternal(p)
	sv := make([]string, 0, len(svv))
	for _, s := range svv {
		if s == "." {
			continue
		}
		if s == ".." {
			if len(sv) == 0 {
				return sv
			}
			sv = sv[0 : len(sv)-1]
			continue
		}
		sv = append(sv, s)
	}
	return sv
}

// PathIsSlipVulnerability todo
func PathIsSlipVulnerability(p string) bool {
	svv := splitPathInternal(p)
	var size int
	for _, s := range svv {
		if s == "." {
			continue
		}
		if s == ".." {
			if size == 0 {
				return true
			}
			size--
			continue
		}
		size++
	}
	return size == 0
}

// StrCat cat strings:
// You should know that StrCat gradually builds advantages
// only when the number of parameters is> 2.
func StrCat(sv ...string) string {
	var sb strings.Builder
	var size int
	for _, s := range sv {
		size += len(s)
	}
	sb.Grow(size)
	for _, s := range sv {
		_, _ = sb.WriteString(s)
	}
	return sb.String()
}

// ByteCat cat strings:
// You should know that StrCat gradually builds advantages
// only when the number of parameters is> 2.
func ByteCat(sv ...[]byte) string {
	var sb bytes.Buffer
	var size int
	for _, s := range sv {
		size += len(s)
	}
	sb.Grow(size)
	for _, s := range sv {
		_, _ = sb.Write(s)
	}
	return sb.String()
}

// BufferCat todo
func BufferCat(sv ...string) []byte {
	var buf bytes.Buffer
	var size int
	for _, s := range sv {
		size += len(s)
	}
	buf.Grow(size)
	for _, s := range sv {
		_, _ = buf.WriteString(s)
	}
	return buf.Bytes()
}

// ErrorCat todo
func ErrorCat(sv ...string) error {
	return errors.New(StrCat(sv...))
}
