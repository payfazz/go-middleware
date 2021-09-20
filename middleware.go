// Package middleware
//
// this package re-export Chain from github.com/payfazz/httpchain
//
// see subpackage for collections of middleware
package middleware

import (
	"net/http"

	"github.com/payfazz/httpchain"
)

// see https://pkg.go.dev/github.com/payfazz/httpchain#Chain
func Chain(all ...interface{}) http.HandlerFunc {
	return httpchain.Chain(all...)
}

// shortcut for Chain function
func C(all ...interface{}) http.HandlerFunc {
	return httpchain.Chain(all...)
}
