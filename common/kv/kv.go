package kv

import (
	"context"
	"net/http"

	"github.com/payfazz/go-middleware"
)

type ctxType int

const ctxKey ctxType = 0

// New create middleware for storing key-value data in request context
func New() middleware.Func {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r.WithContext(context.WithValue(
				r.Context(), ctxKey, make(map[interface{}]interface{})),
			))
		}
	}
}

// NewSetter create middleware for injecting arbitary data.
//
// NewSetter will panic if key-value middleware not installed
func NewSetter(key interface{}, value interface{}) middleware.Func {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			Set(r, key, value)
			next(rw, r)
		}
	}
}

// Get get stored value
//
// Get will panic if key-value middleware not installed
func Get(r *http.Request, key interface{}) interface{} {
	m := r.Context().Value(ctxKey).(map[interface{}]interface{})
	return m[key]
}

// Set set stored value
//
// Set will panic if key-value middleware not installed
func Set(r *http.Request, key interface{}, value interface{}) {
	m := r.Context().Value(ctxKey).(map[interface{}]interface{})
	m[key] = value
}
