// Package kv provide key-value storage middleware.
//
// It more efficient to use this middleware instead of using
// context dirrectly for key-value storage,
// because "net/http.Request.WithContext" always create shallow copy.
package kv

import (
	"context"
	"net/http"
)

type keyType struct{}

var ctxKey keyType

func ensureKV(r *http.Request) (*http.Request, map[interface{}]interface{}) {
	if iface := r.Context().Value(ctxKey); iface != nil {
		return r, iface.(map[interface{}]interface{})
	}

	kv := make(map[interface{}]interface{})
	return r.WithContext(context.WithValue(r.Context(), ctxKey, kv)), kv
}

// Get stored value inside kv
func Get(r *http.Request, key interface{}) (interface{}, bool) {
	if iface := r.Context().Value(ctxKey); iface != nil {
		m := iface.(map[interface{}]interface{})
		v, ok := m[key]
		return v, ok
	}
	return nil, false
}

// MustGet do the same thing as Get, but will panic if it never set before
func MustGet(r *http.Request, key interface{}) interface{} {
	v, ok := Get(r, key)
	if !ok {
		panic("kv: invalid key")
	}
	return v
}

// EnsureKVAndSet will return request that have kv instance in its context,
// it will also set kv entry with provided key and value
func EnsureKVAndSet(r *http.Request, key interface{}, value interface{}) *http.Request {
	r2, m := ensureKV(r)
	m[key] = value
	return r2
}

// Injector return middleware to set data, so next handler will have that value in kv
func Injector(kvKey interface{}, kvData interface{}) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, EnsureKVAndSet(r, kvKey, kvData))
		}
	}
}
