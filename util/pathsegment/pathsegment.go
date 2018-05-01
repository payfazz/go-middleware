// Package pathsegment provide simple utility for routing.
//
// kv middleware is required to be installed in middleware chain.
package pathsegment

import (
	"net/http"
	"strings"

	"github.com/payfazz/go-middleware/common/kv"
)

type keyType struct{}

var key keyType

// Shift first segment of path
func Shift(r *http.Request) string {
	ret, _ := ShiftEnd(r)
	return ret
}

// ShiftEnd is same as Shift, but also indicate the end of the path
func ShiftEnd(r *http.Request) (string, bool) {
	return getState(r).shiftEnd()
}

// Unshift do the reverse of Shift
func Unshift(r *http.Request) {
	getState(r).unshift()
}

// Len return number of segment in path
func Len(r *http.Request) int {
	return getState(r).len()
}

// I return current index of segment in path
func I(r *http.Request) int {
	return getState(r).i()
}

// Get return i-th segment in path
func Get(r *http.Request, i int) string {
	return getState(r).get(i)
}

func getState(r *http.Request) *state {
	if tmp := kv.Get(r, key); tmp != nil {
		return tmp.(*state)
	}

	s := newState(
		strings.TrimPrefix(r.URL.EscapedPath(), "/"),
	)
	kv.Set(r, key, s)

	return s
}
