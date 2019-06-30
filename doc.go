// Package middleware provide simple middleware framework.
//
// It preserve "http.HandlerFunc" signature from "net/http" package, which is good thing
// because it will always compatible with other library that follow this standard library signature.
//
// Why use "http.HandlerFunc" instead of "http.Handler"?,
// because "http.Handler.ServeHTTP" is "http.HandlerFunc", so it more easy to convert between them
//
// Middleware
//
// Middleware is basically a decorator function, it take "http.HandlerFunc"
// and also return "http.HandlerFunc"
//
//	func(http.HandlerFunc) http.HandlerFunc
//
// You can think "Compile" as a way to adding decorator into single http handler,
// for example:
// 	var h = middleware.Compile(
// 		a(someparam),
// 		b,
// 		func(w http.ResponseWriter, r *http.Request) { ... }
// 	)
// is semantically equivalent with python code:
// 	@a(someparam)
// 	@b
// 	def h(w, r):
// 		...
//
// supose you have:
//
// 	var m1 func(http.HandlerFunc) http.HandlerFunc = ...
// 	var m2 func(http.HandlerFunc) http.HandlerFunc = ...
// 	var m3 func(http.HandlerFunc) http.HandlerFunc = ...
// 	var m4 func(http.HandlerFunc) http.HandlerFunc = ...
// 	var m5 func(http.HandlerFunc) http.HandlerFunc = ...
// 	var m6 func(http.HandlerFunc) http.HandlerFunc = ...
// 	var m7 func(http.HandlerFunc) http.HandlerFunc = ...
// 	handler := http.NewServeMux()
// 	var compiled http.HandlerFunc = middleware.C(
// 		m1,
// 		[]interface{}{
// 			m2,
// 			m3,
// 			[]interface{}{
// 				m4, m5,
// 			},
// 			m6,
// 		},
// 		m7,
// 		handler,
// 	)
// 	panic(http.ListenAndServe(":8080", compiled))
//
// then compiled will have same value with
//
// 	compiled := m1(m2(m3(m4(m5(m6(m7(handler.ServeHTTP)))))))
//
// as you can see, this is fast when "compiled" handle the real request
// because all middleware is only compiled once into "compiled",
// and used multiple times
//
// see https://github.com/payfazz/go-middleware
package middleware
