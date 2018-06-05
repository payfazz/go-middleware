// Package recovery provide recovery middleware, it will handle panic in subsequence middleware.
package recovery

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/util/responsewriter"
)

// Event struct for recovery callback
type Event struct {
	Error interface{}
	Stack []struct {
		File string
		Line int
	}
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

// New create recovery middleware, recovery any panic in subsequence middleware.
// If callback is nil, it will write HTTP 500 internal error to client and log to stderr
func New(stackTraceDepth int, callback func(*Event)) middleware.Func {
	if callback == nil {
		callback = func(event *Event) {
			newW := responsewriter.Wrap(event.ResponseWriter)
			if !newW.Written() {
				newW.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(newW,
					fmt.Sprintf("%d %s",
						http.StatusInternalServerError,
						http.StatusText(http.StatusInternalServerError),
					),
				)
			}
			go func() {
				var errMsg interface{}
				switch err := event.Error.(type) {
				case error:
					errMsg = err.Error()
				case fmt.Stringer:
					errMsg = err.String()
				default:
					errMsg = err
				}
				now := time.Now().Format(time.RFC3339)
				if len(event.Stack) > 0 {
					fmt.Fprintf(os.Stderr, "%s | ERR | %v\nSTACK:\n", now, errMsg)
					for _, s := range event.Stack {
						fmt.Fprintf(os.Stderr, "- %s:%d\n", s.File, s.Line)
					}
				} else {
					fmt.Fprintf(os.Stderr, "%s | ERR | %#v\n", now, errMsg)
				}
			}()
		}
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					if rec == http.ErrAbortHandler {
						panic(rec)
					}
					event := Event{
						Error:          rec,
						ResponseWriter: w,
						Request:        r,
					}
					if stackTraceDepth > 0 {
						ptrs := make([]uintptr, stackTraceDepth)
						ptrsNum := runtime.Callers(4, ptrs[:])
						if ptrsNum > 0 {
							frames := runtime.CallersFrames(ptrs)
							for {
								frame, more := frames.Next()
								s := struct {
									File string
									Line int
								}{frame.File, frame.Line}
								if s.File == "" {
									s.File = "*unknown"
								}
								event.Stack = append(event.Stack, s)
								if !more {
									break
								}
							}
						}
					}

					callback(&event)

					panic(http.ErrAbortHandler)
				}
			}()

			newW := responsewriter.Wrap(w)
			next(newW, r)
		}
	}
}
