// Package responsewriter.
//
// This package provide wrapper type for http.ResponseWriter and provide additional functionality to it
package responsewriter

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter

	// Written status
	Status() int

	// Written body size
	Size() int

	// Written
	Written() bool

	// Tells if hijacked
	Hijacked() bool
}

// Wrap a net/http.ResponseWriter, return rw if rw is already wrapped
func Wrap(rw http.ResponseWriter) ResponseWriter {
	if rw, ok := rw.(*responseWritter); ok {
		return rw
	}
	return &responseWritter{rw: rw}
}

type responseWritter struct {
	rw       http.ResponseWriter
	status   int
	size     int
	hijacked bool
}

var (
	_ http.ResponseWriter = (*responseWritter)(nil)
	_ http.Flusher        = (*responseWritter)(nil)
	_ http.Hijacker       = (*responseWritter)(nil)
)

func (rw *responseWritter) Header() http.Header {
	return rw.rw.Header()
}

func (rw *responseWritter) WriteHeader(s int) {
	if !rw.Written() {
		rw.status = s
		rw.rw.WriteHeader(s)
	}
}

func (rw *responseWritter) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.rw.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWritter) Status() int {
	return rw.status
}

func (rw *responseWritter) Size() int {
	return rw.size
}

func (rw *responseWritter) Written() bool {
	return rw.status != 0
}

func (rw *responseWritter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.rw.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("Hijack is not supported")
	}
	c, bufrw, err := hijacker.Hijack()
	rw.hijacked = err == nil
	return c, bufrw, err
}

func (rw *responseWritter) Hijacked() bool {
	return rw.hijacked
}

func (rw *responseWritter) Flush() {
	flusher, ok := rw.rw.(http.Flusher)
	if !ok {
		return
	}
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	flusher.Flush()
}
