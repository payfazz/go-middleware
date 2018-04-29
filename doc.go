/*
	Package middleware provide simple middleware framework.
	it preserve http.Handler signature from net/http package, which is good thing
	because it will always compatible with other library that follow this standard library signature


	Example

	for example usage, see example/sample.go in the source directory

	the content of example/sample.go

		package main

		import (
			"fmt"
			"net/http"

			"github.com/payfazz/go-middleware"
			"github.com/payfazz/go-middleware/common"
		)

		func main() {
			ms := []middleware.Func{m1, m2}
			ms2 := []middleware.Func{m3, m4}
			http.Handle("/", middleware.Compile(
				common.BasicPack(),
				middleware.BuildList(ms, ms2),
				m5,
				handler,
			))

			http.ListenAndServe(":8080", nil)
		}

		func m1(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("m1 before")
				next(w, r)
				fmt.Println("m1 after")
			}
		}

		func m2(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("m2 before")
				next(w, r)
				fmt.Println("m2 after")
			}
		}

		func m3(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("m3 before")
				next(w, r)
				fmt.Println("m3 after")
			}
		}

		func m4(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("m4 before")
				next(w, r)
				fmt.Println("m4 after")
			}
		}

		func m5(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("m5 before")
				next(w, r)
				fmt.Println("m5 after")
			}
		}

		func handler(w http.ResponseWriter, r *http.Request) {
			if r.URL.EscapedPath() == "/panic" {
				panic("test panic")
			}
			fmt.Println("inside handler")
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "Hello World, hai\n")
		}


*/
package middleware
