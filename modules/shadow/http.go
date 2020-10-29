package shadow

import (
	"net/http"
	"time"
)

// ResponseWriter shadow ResponseWriter
type ResponseWriter struct {
	writer     http.ResponseWriter
	rid        string
	written    int64
	startTime  time.Time
	statusCode int
}

// NewResponseWriter bind ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{writer: w, rid: NewRID(), statusCode: http.StatusOK, startTime: time.Now()}
}

// Header get low layer header
func (w *ResponseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write data
func (w *ResponseWriter) Write(data []byte) (int, error) {
	written, err := w.writer.Write(data)
	w.written += int64(written)
	return written, err
}

// WriteHeader write header statucode
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.writer.WriteHeader(statusCode)
}

// StatusCode return statusCode
func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}

// Written return body size
func (w *ResponseWriter) Written() int64 {
	return w.written
}

// RequestID return X-Request-Id
func (w *ResponseWriter) RequestID() string {
	return w.rid
}

// Since since process request
func (w *ResponseWriter) Since() time.Duration {
	return time.Since(w.startTime)
}
