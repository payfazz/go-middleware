// Package kv provide key-value storage.
//
// It more efficient to use this package instead of using context dirrectly for key-value storage,
// because "net/http.Request.WithContext" always create shallow copy.
package kv

import (
	"context"
	"net/http"
)

type kvType struct{}

var kvKey kvType

func ensureReqHaveKV(r *http.Request) (*http.Request, map[interface{}]interface{}) {
	if iface := r.Context().Value(kvKey); iface != nil {
		return r, iface.(map[interface{}]interface{})
	}

	kv := make(map[interface{}]interface{})
	return r.WithContext(context.WithValue(r.Context(), kvKey, kv)), kv
}

// Get stored value inside kv
func Get(r *http.Request, key interface{}) (interface{}, bool) {
	kvVal := r.Context().Value(kvKey)
	if kvVal == nil {
		return nil, false
	}
	m := kvVal.(map[interface{}]interface{})
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

// WithValue will return request that have kv instance in its context,
// and also set kv entry with provided key and value
func WithValue(r *http.Request, key interface{}, value interface{}) *http.Request {
	r2, m := ensureReqHaveKV(r)
	m[key] = value
	return r2
}

// Injector return middleware to set value with the key
func Injector(key interface{}, value interface{}) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, WithValue(r, key, value))
		}
	}
}
