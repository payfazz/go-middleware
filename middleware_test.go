package middleware_test

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/panicreporter"
	"github.com/payfazz/go-middleware/reqlogger"
)

func Example() {
	if err := http.ListenAndServe(":8080", middleware.C(
		panicreporter.New(nil),
		reqlogger.New(nil),
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.EscapedPath() {
			case "/hello":
				fmt.Fprintln(w, "hello world")
				return
			case "/random-panic":
				num := 10 / (rand.Int() % 2)
				fmt.Fprintf(w, "num = %d\n", num)
				return
			default:
				http.Error(w, "not found", 404)
				return
			}
		},
	)); err != nil {
		panic(err)
	}
}
