package env

import (
	"os"
	"sync"
)

var (
	cleanupEnv []string
	delayOnce  sync.Once
)

func delayInitializeEnv() {
	cleanupEnv = make([]string, 0, len(allowedEnv))
	for _, e := range allowedEnv {
		cleanupEnv = append(cleanupEnv, e+"="+os.Getenv(e))
	}
}

// Environ todo
func Environ() []string {
	delayOnce.Do(delayInitializeEnv)
	return cleanupEnv
}
