package responsewriter

import "net/http"

var (
	_ http.Flusher = (*responseWritter)(nil)
)

func (rw *responseWritter) Flush() {
	flusher, ok := rw.rw.(http.Flusher)
	if !ok {
		return
	}
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	flusher.Flush()
}
