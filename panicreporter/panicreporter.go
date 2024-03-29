// Package panicreporter provide middleware to report any golang panic.
package panicreporter

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/payfazz/go-errors/v2"
	"github.com/payfazz/go-errors/v2/trace"

	"github.com/payfazz/go-middleware/util/responsewriter"
)

var filterOutArr = []string{
	"github.com/payfazz/go-middleware",
	"github.com/payfazz/httpchain",
	"net/http",
}

func filter(l trace.Location) bool {
	return !l.InPkg(filterOutArr...)
}

func defaultReporter() func(err error) {
	var w io.Writer = os.Stderr
	prefix := "[PANIC]"
	if isatty.IsTerminal(os.Stderr.Fd()) {
		w = colorable.NewColorable(os.Stderr)
		prefix = "\033[97;41m" + prefix + "\033[0m"
	}

	return func(err error) {
		fmt.Fprintf(w, "%s %v\n%s\n", prefix, time.Now().Format(time.RFC3339Nano), errors.FormatWithFilter(err, filter))
	}
}

// New return panic reporter middleware
//
// When panic occurs, that panic will be reported sing reporter function,
// the error will be generated using https://pkg.go.dev/github.com/payfazz/go-errors/v2 so you can use StackTrace function
// to get where the panic occurs
//
// if reporter is nil, then every panic will be printed to stderr
func New(reporter func(error)) func(http.HandlerFunc) http.HandlerFunc {
	if reporter == nil {
		reporter = defaultReporter()
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w2 := responsewriter.Wrap(w)

			err := errors.Catch(func() error { next(w2, r); return nil })
			if err == nil {
				return
			}

			errors.Catch(func() error { reporter(err); return nil })

			if !cleanWrite500(w2) {
				panic(http.ErrAbortHandler)
			}
		}
	}
}

func cleanWrite500(w responsewriter.ResponseWriter) bool {
	if w.Written() || w.Hijacked() {
		return false
	}

	w.WriteHeader(500)
	fmt.Fprintln(w, "😵 Internal Server Error")
	return true
}
