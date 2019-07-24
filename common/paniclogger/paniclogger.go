// Package paniclogger provide middleware to recover panic, it is just for logging purpose,
// and send http status 500 if possible.
// At the end, it will repanic with http.ErrAbortHandler.
//
// The purpose of this package is only for logging, because with default panic handler
// (it use http.Server.ErrorLog), you cannot reformat the error message.
//
// It is not wise to panic inside http.Handler.ServeHTTP, you should write http 500 error by yourself.
package paniclogger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/payfazz/go-middleware/util/printer"
	"github.com/payfazz/go-middleware/util/responsewriter"
)

// Event struct for callback
type Event struct {
	Error interface{}
	Stack []struct {
		File string
		Line int
	}
}

// Callback func.
type Callback func(Event)

// New return middleware that log any panic (except http.ErrAbortHandler).
// If panic occurs, it will write HTTP 500 Internal server error to client if nothing written yet
// and then close the connection, by repanic with http.ErrAbortHandler.
// If callback is nil, it will use DefaultLogger(nil).
func New(stackTraceDepth int, callback Callback) func(http.HandlerFunc) http.HandlerFunc {
	if callback == nil {
		callback = DefaultLogger(nil)
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			newW := responsewriter.Wrap(w)

			defer func() {
				if rec := recover(); rec != nil {
					if rec == http.ErrAbortHandler {
						panic(rec)
					}
					event := Event{
						Error: rec,
					}
					if stackTraceDepth > 0 {
						ptrs := make([]uintptr, stackTraceDepth)
						ptrsNum := runtime.Callers(4, ptrs)
						if ptrsNum > 0 {
							frames := runtime.CallersFrames(ptrs)
							for {
								frame, more := frames.Next()
								s := struct {
									File string
									Line int
								}{frame.File, frame.Line}
								if s.File == "" {
									s.File = "???"
								}
								event.Stack = append(event.Stack, s)
								if !more {
									break
								}
							}
						}
					}

					if !newW.Written() && !newW.Hijacked() {
						respData := []byte(fmt.Sprintf("%d %s",
							http.StatusInternalServerError,
							http.StatusText(http.StatusInternalServerError),
						))
						newW.Header().Set("Content-Type", "text/plain; charset=utf-8")
						newW.Header().Set("Content-Length", strconv.Itoa(len(respData)))
						newW.Header().Set("Connection", "close")
						newW.WriteHeader(http.StatusInternalServerError)
						newW.Write(respData)
						newW.Flush()
					}

					callback(event)

					panic(http.ErrAbortHandler)
				}
			}()

			next(newW, r)
		}
	}
}

// DefaultLogger return default callback function for this middleware.
// if logger is nil, it will use os.Stderr
func DefaultLogger(logger printer.Printer) Callback {
	if logger == nil {
		logger = log.New(os.Stderr, "", 0)
	}
	return func(event Event) {
		var errMsg interface{}
		switch err := event.Error.(type) {
		case error:
			errMsg = err.Error()
		case fmt.Stringer:
			errMsg = err.String()
		default:
			errMsg = err
		}
		output := fmt.Sprintf("%v | ERR | %#v\n", time.Now().UTC(), errMsg)
		if len(event.Stack) > 0 {
			output += "STACK:\n"
			for _, s := range event.Stack {
				output += fmt.Sprintf("- %s:%d\n", s.File, s.Line)
			}
		}
		logger.Print(output)
	}
}

// NewWithDefaultLogger is same with New(20, DefaultLogger(...)).
func NewWithDefaultLogger(logger printer.Printer) func(http.HandlerFunc) http.HandlerFunc {
	return New(20, DefaultLogger(logger))
}
