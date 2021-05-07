// Package reqlogger provide logger middleware for every http request.
package reqlogger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/payfazz/go-middleware/util/responsewriter"
)

// New return logger middleware
//
// Every request will be logged to dest, if dest is nil then os.Stdout is used
func New(dest io.Writer) func(http.HandlerFunc) http.HandlerFunc {
	if dest == nil {
		dest = os.Stdout
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w2 := responsewriter.Wrap(w)

			start := time.Now()

			next(w2, r)

			var status string
			if w2.Hijacked() {
				status = "Hijacked"
			} else {
				duration := time.Since(start)
				statusCode := w2.Status()
				status = fmt.Sprintf(
					"%d %s | %s",
					statusCode,
					http.StatusText(statusCode),
					duration.Truncate(1*time.Millisecond).String(),
				)
			}

			fmt.Fprintf(dest, "%s %s | %s\n", r.Method, r.URL.EscapedPath(), status)
		}
	}
}
