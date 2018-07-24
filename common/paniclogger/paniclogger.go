// Package paniclogger provide middleware to recover panic, it is just for logging purpose,
// and send http status 500 if possible.
// At the end, it will repanic with http.http.ErrAbortHandler.
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

	"github.com/payfazz/go-middleware"
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

// New return middleware that recover any panic.
// If panic occurs, it will write HTTP 500 Internal server error to client if nothing written yet
// and then close the connection, by repanic with http.ErrAbortHandler.
// If callback is nil, it will log to stderr.
func New(stackTraceDepth int, callback func(*Event)) middleware.Func {
	if callback == nil {
		logger := log.New(os.Stderr, "ERR ", 0)
		callback = func(event *Event) {
			var errMsg interface{}
			switch err := event.Error.(type) {
			case error:
				errMsg = err.Error()
			case fmt.Stringer:
				errMsg = err.String()
			default:
				errMsg = err
			}
			output := fmt.Sprintf("%s | %#v\n", time.Now().Format(time.RFC3339), errMsg)
			if len(event.Stack) > 0 {
				output += "STACK:\n"
				for _, s := range event.Stack {
					output += fmt.Sprintf("- %s:%d\n", s.File, s.Line)
				}
			}
			logger.Println(output)
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
									s.File = "*unknown"
								}
								event.Stack = append(event.Stack, s)
								if !more {
									break
								}
							}
						}
					}

					newW := responsewriter.Wrap(w)
					if !newW.Written() && !newW.Hijacked() {
						respData := []byte(fmt.Sprintf("%d %s",
							http.StatusInternalServerError,
							http.StatusText(http.StatusInternalServerError),
						))
						newW.Header().Set("Content-Type", "text/plain")
						newW.Header().Set("Content-Length", strconv.Itoa(len(respData)))
						newW.WriteHeader(http.StatusInternalServerError)
						newW.Write(respData)
						newW.Flush()
					}

					go callback(&event)

					panic(http.ErrAbortHandler)
				}
			}()

			newW := responsewriter.Wrap(w)
			next(newW, r)
		}
	}
}
