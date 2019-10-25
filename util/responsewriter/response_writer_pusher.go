//+build go1.8

package responsewriter

import (
	"fmt"
	"net/http"
)

// static type check
var (
	_ http.Pusher = (*ResponseWriter)(nil)
)

// Push from net/http.Pusher
func (rw *ResponseWriter) Push(target string, opts *http.PushOptions) error {
	pusher, ok := rw.ResponseWriter.(http.Pusher)
	if ok {
		return pusher.Push(target, opts)
	}
	return fmt.Errorf("the ResponseWriter doesn't support the Pusher interface")
}
