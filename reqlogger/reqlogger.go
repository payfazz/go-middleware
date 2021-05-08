// Package reqlogger provide logger middleware for every http request.
package reqlogger

// The log format is inspired by gin-gonic frameworks
// https://github.com/gin-gonic/gin/blob/4fe5f3e4b4fe62c057ec22caee4908beeef5f59c/logger.go#L143

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"

	"github.com/payfazz/go-middleware/util/responsewriter"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

func statusColorFor(enable bool, code int) string {
	if !enable {
		return ""
	}

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

func statusColorForHijacked(enabled bool) string {
	if !enabled {
		return ""
	}

	return green
}

func methodColorFor(enabled bool, method string) string {
	if !enabled {
		return ""
	}

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

func resetColor(enabled bool) string {
	if !enabled {
		return ""
	}

	return reset
}

// New return logger middleware
//
// Every request will be logged to dest, if dest is nil then os.Stdout is used
func New(dest io.Writer) func(http.HandlerFunc) http.HandlerFunc {
	if dest == nil {
		dest = os.Stdout
	}

	colorEnabled := false
	if v, ok := dest.(*os.File); ok {
		if isatty.IsTerminal(v.Fd()) {
			dest = colorable.NewColorable(v)
			colorEnabled = true
		}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w2 := responsewriter.Wrap(w)

			start := time.Now()

			next(w2, r)

			end := time.Now()
			status := w2.Status()

			latency := end.Sub(start)
			if latency > time.Minute {
				latency = latency.Truncate(time.Second)
			}
			latencyStr := latency.String()
			if w2.Hijacked() {
				latencyStr = "Hijacked"
			}

			var statusStr string
			var statusColor string
			if w2.Hijacked() {
				statusColor = statusColorForHijacked(colorEnabled)
				statusStr = "???"
			} else {
				statusColor = statusColorFor(colorEnabled, status)
				statusStr = fmt.Sprintf("%3d", status)
			}

			methodColor := methodColorFor(colorEnabled, r.Method)
			resetColor := resetColor(colorEnabled)

			fmt.Fprintf(dest, "[REQ] %v |%s %s %s| %13s |%s %-7s %s %#v\n",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusStr, resetColor,
				latencyStr,
				methodColor, r.Method, resetColor,
				r.URL.EscapedPath(),
			)
		}
	}
}
