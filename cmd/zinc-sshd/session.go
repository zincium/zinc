package main

import (
	"time"

	"github.com/gliderlabs/ssh"
)

// Session session
type Session struct {
	ssh.Session
	startTime time.Time
	written   int64
	exitCode  int
}

// NewSession new session
func NewSession(sess ssh.Session) *Session {
	return &Session{Session: sess, startTime: time.Now()}
}

// Exit override exit
func (sess *Session) Exit(code int) error {
	sess.exitCode = code
	return sess.Session.Exit(code)
}

// Write override write
func (sess *Session) Write(data []byte) (int, error) {
	w, err := sess.Session.Write(data)
	sess.written += int64(w)
	return w, err
}
