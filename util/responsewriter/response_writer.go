// Package responsewriter provide a wrapper around http.ResponseWriter.
//
// forked from: https://github.com/urfave/negroni
package responsewriter

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// Wrap a http.ResponseWriter
func Wrap(rw http.ResponseWriter) *ResponseWriter {
	// already ResponseWriter?, return it
	if tmp, ok := rw.(*ResponseWriter); ok {
		return tmp
	}

	return &ResponseWriter{
		ResponseWriter: rw,
	}
}

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type ResponseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []func()
	hijacked    bool
}

// static type check
var (
	_ http.ResponseWriter = (*ResponseWriter)(nil)
	_ http.Flusher        = (*ResponseWriter)(nil)
	_ http.Hijacker       = (*ResponseWriter)(nil)
)

// WriteHeader from http.ResponseWriter
func (rw *ResponseWriter) WriteHeader(s int) {
	rw.status = s
	rw.callBefore()
	rw.ResponseWriter.WriteHeader(s)
}

// Write from http.ResponseWriter
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// Status returns the status code of the response or 0 if the response has
// not been written
func (rw *ResponseWriter) Status() int {
	return rw.status
}

// Size returns the size of the response body.
func (rw *ResponseWriter) Size() int {
	return rw.size
}

// Written returns whether or not the ResponseWriter has been written.
func (rw *ResponseWriter) Written() bool {
	return rw.status != 0
}

// Before allows for a function to be called before the ResponseWriter has been written to. This is
// useful for setting headers or any other operations that must happen before a response has been written.
func (rw *ResponseWriter) Before(before func()) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

// Hijack from http.Hijacker
func (rw *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	c, bufrw, err := hijacker.Hijack()
	rw.hijacked = err == nil
	return c, bufrw, err
}

// Hijacked return true if the underlying ResponseWritter already hijacked
func (rw *ResponseWriter) Hijacked() bool {
	return rw.hijacked
}

func (rw *ResponseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i]()
	}
}

// Flush from http.Flusher
func (rw *ResponseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		if !rw.Written() {
			// The status will be StatusOK if WriteHeader has not been called yet
			rw.WriteHeader(http.StatusOK)
		}
		flusher.Flush()
	}
}
