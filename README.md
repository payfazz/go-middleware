# go-middleware

[![GoDoc](https://godoc.org/github.com/payfazz/go-middleware?status.svg)][godoc]

Fast golang middleware.

This project is about simple http middleware that preserve `http.Handler` and `http.HandlerFunc` signature.

It heavily use clousure and tail call (so it will be faster when tail-cail-optimization implemented on golang in the future). The final middleware chain is precompute, so it should be fast.

see [godoc][godoc] and examples for more information

see also https://github.com/payfazz/go-router for router

## TODO

* More documentation and examples

[godoc]: https://godoc.org/github.com/payfazz/go-middleware
