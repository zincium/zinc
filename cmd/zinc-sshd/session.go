package main

import (
	"time"

	"github.com/gliderlabs/ssh"
)

// type Session interface {
// 	gossh.Channel

// 	// User returns the username used when establishing the SSH connection.
// 	User() string

// 	// RemoteAddr returns the net.Addr of the client side of the connection.
// 	RemoteAddr() net.Addr

// 	// LocalAddr returns the net.Addr of the server side of the connection.
// 	LocalAddr() net.Addr

// 	// Environ returns a copy of strings representing the environment set by the
// 	// user for this session, in the form "key=value".
// 	Environ() []string

// 	// Exit sends an exit status and then closes the session.
// 	Exit(code int) error

// 	// Command returns a shell parsed slice of arguments that were provided by the
// 	// user. Shell parsing splits the command string according to POSIX shell rules,
// 	// which considers quoting not just whitespace.
// 	Command() []string

// 	// RawCommand returns the exact command that was provided by the user.
// 	RawCommand() string

// 	// Subsystem returns the subsystem requested by the user.
// 	Subsystem() string

// 	// PublicKey returns the PublicKey used to authenticate. If a public key was not
// 	// used it will return nil.
// 	PublicKey() PublicKey

// 	// Context returns the connection's context. The returned context is always
// 	// non-nil and holds the same data as the Context passed into auth
// 	// handlers and callbacks.
// 	//
// 	// The context is canceled when the client's connection closes or I/O
// 	// operation fails.
// 	Context() context.Context

// 	// Permissions returns a copy of the Permissions object that was available for
// 	// setup in the auth handlers via the Context.
// 	Permissions() Permissions

// 	// Pty returns PTY information, a channel of window size changes, and a boolean
// 	// of whether or not a PTY was accepted for this session.
// 	Pty() (Pty, <-chan Window, bool)

// 	// Signals registers a channel to receive signals sent from the client. The
// 	// channel must handle signal sends or it will block the SSH request loop.
// 	// Registering nil will unregister the channel from signal sends. During the
// 	// time no channel is registered signals are buffered up to a reasonable amount.
// 	// If there are buffered signals when a channel is registered, they will be
// 	// sent in order on the channel immediately after registering.
// 	Signals(c chan<- Signal)

// 	// Break regisers a channel to receive notifications of break requests sent
// 	// from the client. The channel must handle break requests, or it will block
// 	// the request handling loop. Registering nil will unregister the channel.
// 	// During the time that no channel is registered, breaks are ignored.
// 	Break(c chan<- bool)
// }

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
