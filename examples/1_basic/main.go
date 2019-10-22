package main

import (
	"fmt"
	"net/http"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
	"github.com/payfazz/go-middleware/common/logger"
	"github.com/payfazz/go-middleware/common/paniclogger"
)

func main() {
	basicPack := []interface{}{
		paniclogger.NewWithDefaultLogger(nil),
		logger.NewWithDefaultLogger(nil),
		kv.New(),
	}

	group1 := []interface{}{m1, m2}
	group2 := []interface{}{m3, m4}

	http.Handle("/", middleware.C(
		basicPack,
		middleware.CompileList(group1, group2),
		m5,
		handler,
	))

	http.ListenAndServe(":8080", nil)
}

const testKey = "test-key"

func m1(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("m1 before")
		kv.Set(r, testKey, "test-data")
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
	fmt.Println(kv.MustGet(r, testKey))
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello World, hai\n")
}
