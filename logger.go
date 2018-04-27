package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// LoggerEvent struct for logger callback
type LoggerEvent struct {
	StartTime time.Time
	Status    int
	Duration  time.Duration
	Hostname  string
	Method    string
	Path      string
	Request   *http.Request
}

// NewLogger create logger middleware, callback will be called for every request.
// If callback is nil, it will log to stdout
func NewLogger(callback func(*LoggerEvent)) Func {
	if callback == nil {
		callback = func(le *LoggerEvent) {
			fmt.Printf(
				"%s | %d | %v | %s | %s %s\n",
				le.StartTime.Format(time.RFC3339),
				le.Status,
				le.Duration,
				le.Hostname,
				le.Method,
				le.Path,
			)
		}
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			le := LoggerEvent{
				StartTime: time.Now(),
				Hostname:  r.Host,
				Method:    r.Method,
				Path:      r.URL.Path,
				Request:   r,
			}
			newW := NewResponseWriter(w)
			next(newW, r)
			le.Duration = time.Since(le.StartTime)
			le.Status = newW.Status()
			callback(&le)
		}
	}
}
