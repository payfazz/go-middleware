// Package kv provide key-value storage middleware.
// It more efficient to use this middleware instead of using
// context dirrectly for key-value storage,
// because http.Request.WithContext always create shallow copy.
// kv implemented using map, so it is not safe to access it concurrently.
package kv

import (
	"context"
	"net/http"

	middleware "github.com/payfazz/go-middleware"
)

type ctxType struct{}

var ctxKey ctxType

// New return middleware for storing key-value data in request context
func New() middleware.Func {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r.WithContext(context.WithValue(
				r.Context(), ctxKey, make(map[interface{}]interface{})),
			))
		}
	}
}

// NewSetter return middleware for injecting arbitary data.
//
// kv middleware must be installed in the middleware chain. see Set.
func NewSetter(key interface{}, value interface{}) middleware.Func {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			Set(r, key, value)
			next(w, r)
		}
	}
}

// Get stored value
//
// Get will panic if kv middleware not installed in the middleware chain.
func Get(r *http.Request, key interface{}) interface{} {
	m := r.Context().Value(ctxKey).(map[interface{}]interface{})
	return m[key]
}

// Set set stored value
//
// Set will panic if kv middleware not installed in the middleware chain.
func Set(r *http.Request, key interface{}, value interface{}) {
	m := r.Context().Value(ctxKey).(map[interface{}]interface{})
	m[key] = value
}
