package handler

import (
	"bytes"
	"io"
	"net/http"
)

type ResponseWriter struct {
	statusCode  int
	httpWriter  http.ResponseWriter
	multiWriter io.Writer

	bodyBuffer *bytes.Buffer
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	b := &ResponseWriter{httpWriter: w}
	b.bodyBuffer = &bytes.Buffer{}
	b.multiWriter = io.MultiWriter(b.httpWriter, b.bodyBuffer)
	return b
}

// http.ResponseWriter interface implementations

func (r *ResponseWriter) Header() http.Header {
	return r.httpWriter.Header()
}
func (r *ResponseWriter) Write(i []byte) (int, error) {
	return r.multiWriter.Write(i)
}
func (m *ResponseWriter) WriteHeader(statusCode int) {
	m.httpWriter.WriteHeader(statusCode)
	m.statusCode = statusCode
}

// singular functions

func (r *ResponseWriter) StatusCode() int {
	return r.statusCode
}
