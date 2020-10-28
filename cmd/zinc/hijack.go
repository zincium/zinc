package main

import "net/http"

// HijackResponseWriter hijack ResponseWriter
type HijackResponseWriter struct {
	writer     http.ResponseWriter
	statusCode int
	written    int64
}

// NewHijackResponseWriter bind ResponseWriter
func NewHijackResponseWriter(w http.ResponseWriter) *HijackResponseWriter {
	return &HijackResponseWriter{writer: w, statusCode: http.StatusOK}
}

// Header get low layer header
func (hw *HijackResponseWriter) Header() http.Header {
	return hw.writer.Header()
}

// Write data hijack
func (hw *HijackResponseWriter) Write(data []byte) (int, error) {
	hw.written += int64(len(data))
	return hw.writer.Write(data)
}

// WriteHeader write header statucode hijack
func (hw *HijackResponseWriter) WriteHeader(statusCode int) {
	hw.statusCode = statusCode
	hw.writer.WriteHeader(statusCode)
}

// StatusCode return statusCode
func (hw *HijackResponseWriter) StatusCode() int {
	return hw.statusCode
}

// Written return body size
func (hw *HijackResponseWriter) Written() int64 {
	return hw.written
}
