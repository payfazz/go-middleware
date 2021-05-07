//+build go1.8

package responsewriter

import (
	"fmt"
	"net/http"
)

var (
	_ http.Pusher = (*responseWritter)(nil)
)

func (rw *responseWritter) Push(target string, opts *http.PushOptions) error {
	pusher, ok := rw.rw.(http.Pusher)
	if !ok {
		return fmt.Errorf("Push is not supported")
	}
	return pusher.Push(target, opts)
}
