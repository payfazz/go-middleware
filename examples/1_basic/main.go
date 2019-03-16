package main

import (
	"fmt"
	"net/http"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common"
)

func main() {
	ms := []func(http.HandlerFunc) http.HandlerFunc{m1, m2}
	ms2 := []func(http.HandlerFunc) http.HandlerFunc{m3, m4}
	http.Handle("/", middleware.Compile(
		common.BasicPack(),
		middleware.CompileList(ms, ms2),
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
	} else if r.URL.EscapedPath() == "/hijack" {
		if hj, ok := w.(http.Hijacker); !ok {
			panic("not hijacker")
		} else {
			conn, bufrw, err := hj.Hijack()
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			fmt.Fprintf(bufrw, "HTTP/1.1 200 OK\r\n")
			fmt.Fprintf(bufrw, "Connection: Close\r\n")
			fmt.Fprintf(bufrw, "\r\n")
			fmt.Fprintf(bufrw, "test hijacker")
			bufrw.Flush()
		}
		return
	}
	fmt.Println("inside handler")
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello World, hai\n")
}
