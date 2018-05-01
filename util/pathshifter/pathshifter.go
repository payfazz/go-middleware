// Package pathshifter provide simple utility for routing.
//
// kv middleware is required to be installed in middleware chain.
package pathshifter

import (
	"net/http"
	"strings"

	"github.com/payfazz/go-middleware/common/kv"
)

type keyType struct{}

var key keyType

type state struct {
	list  []string
	index int
}

func (s *state) shiftEnd() (string, bool) {
	if s.index == len(s.list) {
		return "", true
	}
	ret := s.list[s.index]
	s.index++
	return ret, s.index == len(s.list)
}

func (s *state) unshift() {
	if s.index == 0 {
		return
	}
	s.index--
}

// Shift return first segment of path, and remove it from its internal state.
func Shift(r *http.Request) string {
	ret, _ := shiftEnd(r)
	return ret
}

// ShiftEnd is same as Shift, but also indicate the end of the path
func ShiftEnd(r *http.Request) (string, bool) {
	return shiftEnd(r)
}

func shiftEnd(r *http.Request) (string, bool) {
	tryToInit(r)
	s := kv.Get(r, key).(*state)
	ret, ok := s.shiftEnd()
	kv.Set(r, key, s)
	return ret, ok
}

// Unshift do the reverse of Shift
func Unshift(r *http.Request) {
	tryToInit(r)
	s := kv.Get(r, key).(*state)
	s.unshift()
	kv.Set(r, key, s)
}

func tryToInit(r *http.Request) {
	if kv.Get(r, key) == nil {
		path := r.URL.EscapedPath()
		path = strings.TrimPrefix(path, "/")
		kv.Set(r, key, &state{strings.Split(path, "/"), 0})
	}
}
