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

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher

	// Status returns the status code of the response or 0 if the response has
	// not been written
	Status() int

	// Written returns whether or not the ResponseWriter has been written.
	Written() bool

	// Size returns the size of the response body.
	Size() int

	// Before allows for a function to be called before the ResponseWriter has been written to. This is
	// useful for setting headers or any other operations that must happen before a response has been written.
	Before(func())

	// Hijacked return true if the underlying ResponseWritter already hijacked
	Hijacked() bool

	// Original return the original http.ResponseWriter
	Original() http.ResponseWriter

	// internal is just empty function, the purpose is to make this interface cannot be implemented outside this package
	internal()
}

// Wrap a http.ResponseWriter
func Wrap(rw http.ResponseWriter) ResponseWriter {
	// already ResponseWriter?, return it
	if tmp, ok := rw.(ResponseWriter); ok {
		return tmp
	}

	return &responseWriter{
		ResponseWriter: rw,
	}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []func()
	hijacked    bool
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.status = s
	rw.callBefore()
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) Before(before func()) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	c, bufrw, err := hijacker.Hijack()
	rw.hijacked = err == nil
	return c, bufrw, err
}

func (rw *responseWriter) Hijacked() bool {
	return rw.hijacked
}

func (rw *responseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i]()
	}
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		if !rw.Written() {
			// The status will be StatusOK if WriteHeader has not been called yet
			rw.WriteHeader(http.StatusOK)
		}
		flusher.Flush()
	}
}

func (rw *responseWriter) Original() http.ResponseWriter {
	return rw.ResponseWriter
}

func (rw *responseWriter) internal() {}
