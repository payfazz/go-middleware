// Package kv provide key-value storage middleware.
//
// It more efficient to use this middleware instead of using
// context dirrectly for key-value storage,
// because "net/http.Request.WithContext" always create shallow copy.
//
// Panic
//
// function in this package do nothing to check if kv is already present in request context.
// therefore if you forget to add it into request context, function like Get and Set will panic.
package kv

import (
	"context"
	"net/http"
)

type keyType struct{}

var ctxKey keyType

// WrapRequest make sure that the request have an kv in request context
func WrapRequest(r *http.Request) *http.Request {
	if tmp := r.Context().Value(ctxKey); tmp != nil {
		return r
	}

	return r.WithContext(context.WithValue(
		r.Context(), ctxKey, make(map[interface{}]interface{})),
	)
}

// New return middleware to make sure that next handler will have kv in request context
func New() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, WrapRequest(r))
		}
	}
}

// Get stored value inside kv
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

// Set value inside kv
func Set(r *http.Request, key interface{}, value interface{}) {
	m := r.Context().Value(ctxKey).(map[interface{}]interface{})
	m[key] = value
}

// Injector return middleware to set data, so next handler will have that value in kv
func Injector(kvKey interface{}, kvData interface{}) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			Set(r, kvKey, kvData)
			next(w, r)
		}
	}
}
