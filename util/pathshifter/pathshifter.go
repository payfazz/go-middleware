// Package pathshifter provide simple utility for routing.
//
// kv middleware is required to be installed in middleware chain.
package pathshifter

import (
	"net/http"
	"strings"

	"github.com/payfazz/go-middleware/common/kv"
)

type keyType int

const key keyType = 0

type state struct {
	list  []string
	index int
}

func (s *state) shift() string {
	if s.index == len(s.list) {
		return ""
	}
	ret := s.list[s.index]
	s.index++
	return ret
}

func (s *state) unshift() {
	if s.index == 0 {
		return
	}
	s.index--
}

// Shift return first segment of path, and remove it from its internal state.
func Shift(r *http.Request) string {
	tryToInit(r)
	s := kv.Get(r, key).(*state)
	ret := s.shift()
	kv.Set(r, key, s)
	return ret
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
