package pathsegment

import "strings"

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

func (s *state) len() int {
	return len(s.list)
}

func (s *state) i() int {
	return s.index
}

func (s *state) get(i int) string {
	if i < 0 || i >= len(s.list) {
		return ""
	}
	return s.list[i]
}

func newState(path string) *state {
	return &state{strings.Split(path, "/"), 0}
}
