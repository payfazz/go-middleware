# go-middleware

Fast golang middleware.

This project is about simple http middleware that preserve `http.Handler` and `http.HandlerFunc` signature.

I like `negroni` but they introduce new signature `func(http.ResponseWriter, *http.Request, http.HandlerFunc)` which is in my opinion that is not good, becase we already have standard signature `func(http.ResponseWriter, *http.Request)` that defined by golang itself.

It heavily use clousure and tail call, so it will be faster when tail-cail-optimization implemented on golang. The final middleware chain is precompute, so it should be faster.

for usage see examples directory

see also https://github.com/payfazz/go-router for router


## TODO

* More documentation and examples
* create testing. when i wrote this, this project was part of bigger project, all of the testing was done there.
