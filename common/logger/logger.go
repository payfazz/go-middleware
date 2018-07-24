// Package logger provide logger middleware.
package logger

import (
	"log"
	"net/http"
	"os"
	"strconv"
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
	Hijacked  bool
	Request   *http.Request
}

// New return logger middleware, callback will be called for every request.
// If callback is nil, it will log to stdout.
// Do not modif Event.Request, and do not access it after the callback return
func New(callback func(*Event)) middleware.Func {
	if callback == nil {
		logger := log.New(os.Stdout, "REQ ", 0)
		callback = func(event *Event) {
			go func() {
				var status string
				if event.Hijacked {
					status = "Hijacked"
				} else {
					status = strconv.Itoa(event.Status)
				}
				logger.Printf(
					"%s | %s | %v | %s | %s %s\n",
					event.StartTime.Format(time.RFC3339),
					status,
					event.Duration.Truncate(1*time.Millisecond),
					event.Hostname,
					event.Method,
					event.Path,
				)
			}()
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
			event.Hijacked = newW.Hijacked()

			callback(&event)
		}
	}
}
