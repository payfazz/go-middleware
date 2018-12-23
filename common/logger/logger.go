// Package logger provide logger middleware.
package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	middleware "github.com/payfazz/go-middleware"
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

// Callback func.
// Do not modif Event.Request, and do not access it after the callback return
type Callback func(Event)

// New return logger middleware, callback will be called for every request.
// If callback is nil, it will log to stdout using DefaultLogger.
func New(callback Callback) middleware.Func {
	if callback == nil {
		logger := log.New(os.Stdout, "REQ ", log.LstdFlags)
		callback = DefaultLogger(logger)
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

			callback(event)
		}
	}
}

// DefaultLogger return default callback function for this middleware.
// logger can't be nil
func DefaultLogger(logger *log.Logger) Callback {
	if logger == nil {
		panic("logger: log can't be nil")
	}
	return func(event Event) {
		go func() {
			var status string
			if event.Hijacked {
				status = "Hijacked"
			} else {
				status = fmt.Sprintf("%d %s", event.Status, http.StatusText(event.Status))
			}
			logger.Printf(
				"%s | %v | %s %s\n",
				status,
				event.Duration.Truncate(1*time.Millisecond),
				event.Method,
				event.Path,
			)
		}()
	}
}
