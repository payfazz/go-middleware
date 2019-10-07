// Package kv provide key-value storage middleware.
//
// It more efficient to use this middleware instead of using
// context dirrectly for key-value storage,
// because "http.Request.WithContext" always create shallow copy.
// kv implemented using map, so it is not safe to access it concurrently.
package kv

import (
	"context"
	"net/http"
)

type keyType struct{}

var ctxKey keyType

// WrapRequest make sure that the request have an instance of
// kv middleware inside "r.Context()"
func WrapRequest(r *http.Request) *http.Request {
	if tmp := r.Context().Value(ctxKey); tmp != nil {
		return r
	}

	return r.WithContext(context.WithValue(
		r.Context(), ctxKey, make(map[interface{}]interface{})),
	)
}

// New return middleware for storing key-value data in request context.
//
// It is common to use this as first middleware in "middleware.Compile"
func New() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, WrapRequest(r))
		}
	}
}

// Get stored value
//
// Get will panic if kv middleware not installed in the middleware chain.
func Get(r *http.Request, key interface{}) (interface{}, bool) {
	m := r.Context().Value(ctxKey).(map[interface{}]interface{})
	v, ok := m[key]
	return v, ok
}

// MustGet do the same thing as Get, but will panic if it never set before
func MustGet(r *http.Request, key interface{}) interface{} {
	v, ok := Get(r, key)
	if !ok {
		panic("kv: invalid key")
	}
	return v
}

// Set set stored value
//
// Set will panic if kv middleware not installed in the middleware chain.
func Set(r *http.Request, key interface{}, value interface{}) {
	m := r.Context().Value(ctxKey).(map[interface{}]interface{})
	m[key] = value
}

// Injector return middleware to inject data
func Injector(kvKey interface{}, kvData interface{}) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			Set(r, kvKey, kvData)
			next(w, r)
		}
	}
}
