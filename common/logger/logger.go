// Package logger provide logger middleware.
package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/util/responsewriter"
)

// Event struct for logger callback
type Event struct {
	StartTime time.Time
	Status    int
	Duration  time.Duration
	Hostname  string
	Method    string
	Path      string
	Request   *http.Request
}

// New create logger middleware, callback will be called for every request.
// If callback is nil, it will log to stdout
func New(callback func(*Event)) middleware.Func {
	if callback == nil {
		callback = func(event *Event) {
			fmt.Printf(
				"%s | REQ | %d | %v | %s | %s %s\n",
				event.StartTime.Format(time.RFC3339),
				event.Status,
				event.Duration.Truncate(1*time.Millisecond),
				event.Hostname,
				event.Method,
				event.Path,
			)
		}
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			event := Event{
				StartTime: time.Now(),
				Hostname:  r.Host,
				Method:    r.Method,
				Path:      r.URL.Path,
				Request:   r,
			}
			newW := responsewriter.Wrap(w)
			next(newW, r)
			event.Duration = time.Since(event.StartTime)
			event.Status = newW.Status()

			callback(&event)
		}
	}
}
