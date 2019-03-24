// Package middleware provide simple middleware framework.
// it preserve http.HandlerFunc signature from net/http package, which is good thing
// because it will always compatible with other library that follow this standard library signature
//
// Middleware is func(http.HandlerFunc) http.HandlerFunc.
//
// for example usage, see examples directory
//
// see also https://gist.github.com/win-t/8a243301bd227cca6135374cf94d9e98 for example usage of go-middleware and go-router
package middleware

import (
	"net/http"
	"reflect"
)

// Compile all middleware into single http.HandlerFunc.
// Compile have same argument meaning with CompileList.
// Compile in fact call CompileList for it arguments.
//
// You can think Compile as a way to adding decorator,
// for example:
//	var h = middleware.Compile(
//		a(someparam),
//		b,
//		func(w http.ResponseWriter, r *http.Request) { ... }
//	)
// is semantically equivalent with python code:
//	@a(someparam)
//	@b
//	def h(w, r):
//		...
func Compile(all ...interface{}) http.HandlerFunc {
	var f http.HandlerFunc
	list := CompileList(all...)
	for i := len(list) - 1; i >= 0; i-- {
		f = list[i](f)
	}
	return f
}

// CompileList will flatten all params into single array of middleware, basically:
// 	CompileList(m1, m2, [m3, m4, [m5, m6]], m7) -> [m1, m2, m3, m4, m5, m6, m7]
// and also will convert http.HandlerFunc and http.Handler into leaf middleware,
// that middleware will not call next, i.e. stopping the chain,
// suitable for last handler in the chain.
func CompileList(all ...interface{}) []func(http.HandlerFunc) http.HandlerFunc {
	ret := make([]func(http.HandlerFunc) http.HandlerFunc, 0, len(all))
	for _, item := range all {
		if item == nil {
			panic("middleware: invalid argument: can't be nil")
		}

		itemValue := reflect.ValueOf(item)
		itemType := itemValue.Type()

		var m func(http.HandlerFunc) http.HandlerFunc
		mValue := reflect.ValueOf(&m).Elem()
		mType := mValue.Type()
		if itemType.ConvertibleTo(mType) {
			mValue.Set(itemValue.Convert(mType))
			ret = append(ret, m)
			continue
		}

		var hf http.HandlerFunc
		hfValue := reflect.ValueOf(&hf).Elem()
		hfType := hfValue.Type()
		if itemType.ConvertibleTo(hfType) {
			hfValue.Set(itemValue.Convert(hfType))
			ret = append(ret, func(http.HandlerFunc) http.HandlerFunc {
				return hf
			})
			continue
		}

		var h http.Handler
		hValue := reflect.ValueOf(&h).Elem()
		hType := hValue.Type()
		if itemType.ConvertibleTo(hType) {
			hValue.Set(itemValue.Convert(hType))
			ret = append(ret, func(http.HandlerFunc) http.HandlerFunc {
				return h.ServeHTTP
			})
			continue
		}

		switch itemType.Kind() {
		case reflect.Slice, reflect.Array:
			args := make([]interface{}, itemValue.Len())
			for i := 0; i < itemValue.Len(); i++ {
				args[i] = itemValue.Index(i).Interface()
			}
			ret = append(ret, CompileList(args...)...)
		default:
			panic("middleware: invalid argument: " + itemType.String() + " can't be converted middleware")
		}
	}
	return ret
}

// Nop is dummy middleware, it just return next
func Nop(next http.HandlerFunc) http.HandlerFunc {
	return next
}
