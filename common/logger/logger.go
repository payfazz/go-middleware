// Package logger provide logger middleware for every http request.
package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/payfazz/go-middleware/util/printer"
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
//
// Do not modif Event.Request, Event.Request.Body maybe already processed and closed
type Callback func(Event)

// New return logger middleware, callback will be called for every request.
// If callback is nil, it will use DefaultLogger(nil).
func New(callback Callback) func(http.HandlerFunc) http.HandlerFunc {
	if callback == nil {
		callback = DefaultLogger(nil)
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
// if logger is nil, it will use os.Stdout
func DefaultLogger(logger printer.Printer) Callback {
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}

	return func(event Event) {
		var status string
		if event.Hijacked {
			status = "Hijacked"
		} else {
			status = fmt.Sprintf("%d %s", event.Status, http.StatusText(event.Status))
		}
		logger.Print(fmt.Sprintf(
			"%v | %s | %v | %s %s\n",
			time.Now().UTC(),
			status,
			event.Duration.Truncate(1*time.Millisecond),
			event.Method,
			event.Request.URL.String(),
		))
	}
}

// NewWithDefaultLogger is same with New(DefaultLogger(...)).
func NewWithDefaultLogger(logger printer.Printer) func(http.HandlerFunc) http.HandlerFunc {
	return New(DefaultLogger(logger))
}
