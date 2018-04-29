// Package recovery provide recovery middleware, it will handle panic in subsequence middleware.
package recovery

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/payfazz/go-middleware"
)

type stringer interface {
	String() string
}

// Event struct for recovery callback
type Event struct {
	Message string
	Stack   []struct {
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
			event.ResponseWriter.WriteHeader(http.StatusInternalServerError)

			if len(event.Stack) > 0 {
				fmt.Fprintf(os.Stderr, "ERR: %s\nSTACK:\n", event.Message)
				for _, s := range event.Stack {
					fmt.Fprintf(os.Stderr, "- %s:%d\n", s.File, s.Line)
				}
			} else {
				fmt.Fprintf(os.Stderr, "ERR: %s\n", event.Message)
			}
		}
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					event := Event{
						ResponseWriter: w,
						Request:        r,
					}
					switch tmp := rec.(type) {
					case error:
						event.Message = tmp.Error()
					case stringer:
						event.Message = tmp.String()
					case string:
						event.Message = tmp
					default:
						event.Message = "unknown error"
					}
					if stackTraceDepth > 0 {
						ptrs := make([]uintptr, stackTraceDepth)
						ptrsNum := runtime.Callers(4, ptrs[:])
						for i := 0; i < ptrsNum; i++ {
							s := struct {
								File string
								Line int
							}{"*unknown", 0}
							if fn := runtime.FuncForPC(ptrs[i]); fn != nil {
								s.File, s.Line = fn.FileLine(ptrs[i])
							}
							event.Stack = append(event.Stack, s)
						}
					}

					callback(&event)
				}
			}()
			next(w, r)
		}
	}
}
