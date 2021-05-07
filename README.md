# go-middleware

[![Go Reference](https://pkg.go.dev/badge/github.com/payfazz/go-middleware.svg)](https://pkg.go.dev/github.com/payfazz/go-middleware)

Golang middleware collection

Middleware is value that have following type
```go
func(http.HandlerFunc) http.HandlerFunc
```

In following example, `testMiddleware` is a middleware that adding some header to the response
```go
var handler http.HandlerFunc = ...
var someMiddleware func(next http.HandlerFunc) http.HanlderFunc = ...

testMiddleware := func(next http.HandlerFunc) http.HanlderFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("TestHeader", "test value")
    next(w, r)
  }
}

combinedHandler := someMiddleware(testMiddleware(handler))
```

Use [httpchain](https://pkg.go.dev/github.com/payfazz/httpchain) package to chaining multiple middleware, so you can write
```go
combinedHandler := httpchain.Chain(
  someMiddleware,
  testMiddleware,
  handler,
)
```
