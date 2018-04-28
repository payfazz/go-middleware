package middleware

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
)

type stringer interface {
	String() string
}

// RecoveryEvent struct for recovery callback
type RecoveryEvent struct {
	Message string
	Stack   []struct {
		File string
		Line int
	}
}

// fmt.Sprintf("%s:%d", file, line)

// NewRecovery create recovery middleware, recovery any panic in subsequence middleware.
// If callback is nil, it will log to stderr
func NewRecovery(stackTraceDepth int, callback func(*RecoveryEvent)) Func {
	if callback == nil {
		callback = func(re *RecoveryEvent) {
			if len(re.Stack) > 0 {
				fmt.Fprintf(os.Stderr, "ERR: %s\nSTACK:\n", re.Message)
				for _, s := range re.Stack {
					fmt.Fprintf(os.Stderr, "- %s:%d\n", s.File, s.Line)
				}
			} else {
				fmt.Fprintf(os.Stderr, "ERR: %s\n", re.Message)
			}
		}
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					w.WriteHeader(http.StatusInternalServerError)

					re := RecoveryEvent{}

					switch tmp := rec.(type) {
					case error:
						re.Message = tmp.Error()
					case stringer:
						re.Message = tmp.String()
					case string:
						re.Message = tmp
					default:
						re.Message = "unknown error"
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
							re.Stack = append(re.Stack, s)
						}
					}

					callback(&re)
				}
			}()
			next(w, r)
		}
	}
}
