package panicreporter

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/payfazz/go-errors/v2"
	"github.com/payfazz/go-errors/v2/trace"
	"github.com/payfazz/go-middleware/util/responsewriter"
)

func New(reporter func(error)) func(http.HandlerFunc) http.HandlerFunc {
	if reporter == nil {
		matcher := regexp.MustCompile(`` +
			`^(` +
			`github.com/payfazz/go-middleware|` +
			`github.com/payfazz/httpchain|` +
			`net/http` +
			`)(\.|\/)`,
		)
		reporter = func(err error) {
			fmt.Fprint(os.Stderr, errors.FormatWithFilter(err,
				func(l trace.Location) bool {
					return !matcher.MatchString(l.Func())
				},
			))
		}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w2 := responsewriter.Wrap(w)

			err := errors.Catch(func() error {
				next(w2, r)
				return nil
			})

			if err == nil {
				return
			}

			errors.Catch(func() error {
				reporter(err)
				return nil
			})

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
	fmt.Fprint(w, "😵 Internal Server Error")
	return true
}
