package env

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Expander extend
type Expander struct {
	m  map[string]string
	mu sync.RWMutex
}

// APPDIR todo
func APPDIR() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	exedir := filepath.Dir(exe)
	if filepath.Base(exedir) == "bin" {
		return filepath.Dir(exedir)
	}
	return exedir
}

// NewExpander new expander
func NewExpander() *Expander {
	e := &Expander{m: make(map[string]string)}
	if APPDIR := APPDIR(); APPDIR != "." {
		e.m["APPDIR"] = APPDIR
	}
	return e
}

// AddBashCompatible add bash compatible
func (e *Expander) AddBashCompatible() {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i := 0; i < len(os.Args); i++ {
		e.m[strconv.Itoa(i)] = os.Args[i]
	}
	e.m["$"] = strconv.Itoa(os.Getpid())
}

// Setenv setenv
func (e *Expander) Setenv(k, v string, force bool) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.m[k]; ok {
		if !force {
			return false
		}
	}
	e.m[k] = v
	return true
}

// Getenv get env
func (e *Expander) Getenv(k string) string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if v, ok := e.m[k]; ok {
		return v
	}
	return os.Getenv(k)
}

// Expand todo
func (e *Expander) Expand(s string) string {
	return os.Expand(s, e.Getenv)
}
