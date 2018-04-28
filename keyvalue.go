package middleware

import (
	"context"
	"net/http"
)

type kvCtxType int

const kvKey kvCtxType = 0

// WithKV will inject map[interface{}]interface{} into request context,
// this will simplify passing data between middleware.
// set with SetKV and get with GetKV
//
// Internally WithKV use r.WithContext() and context.WithValue().
func WithKV() Func {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r.WithContext(context.WithValue(
				r.Context(), kvKey, make(map[interface{}]interface{})),
			))
		}
	}
}

// InjectKV will return middleware for injecting arbitary data.
// InjectKV require WithKV middleware installed
func InjectKV(key interface{}, value interface{}) Func {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			SetKV(r, key, value)
			next(rw, r)
		}
	}
}

// GetKV require WithKV middleware installed
func GetKV(r *http.Request, key interface{}) interface{} {
	m := r.Context().Value(kvKey).(map[interface{}]interface{})
	return m[key]
}

// SetKV require WithKV middleware installed
func SetKV(r *http.Request, key interface{}, value interface{}) {
	m := r.Context().Value(kvKey).(map[interface{}]interface{})
	m[key] = value
}
