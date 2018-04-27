// Package middleware provide simple middleware framework
// it preserve http.Handler signature from net/http package, which is good thing
// because it will always compatible with other library that follow this standard library signature
//
// for example usage, see example/sample.go in the source directory
package middleware

import (
	"net/http"
	"reflect"
)

// Func is function alias that take http.HandlerFunc as next middleware in the chains
//
// example:
//	middleware := func(next http.HandlerFunc) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			fmt.Println("before")
//			next(w, r)
//			fmt.Println("after")
//		}
//	}
//
// next will be nil if this middleware is last middleware in chain
type Func func(next http.HandlerFunc) http.HandlerFunc

// Compile all middleware into single http.HandlerFunc.
// Compile will call BuildList to convert all into []Func
func Compile(all ...interface{}) http.HandlerFunc {
	var f http.HandlerFunc
	list := BuildList(all...)
	for i := len(list) - 1; i >= 0; i-- {
		f = list[i](f)
	}
	return f
}

// BuildList will convert all into []Func, basically:
// 	BuildList(m1, m2, [m3, m4, [m5, m6]], m7) -> [m1, m2, m3, m4, m5, m6, m7]
// and also will convert http.Handler, http.HandlerFunc and func(http.ResponseWriter, *http.Request)
// into Func, that Func will not call next, i.e. stopping the chain,
// suitable for last handler in the chain
func BuildList(all ...interface{}) []Func {
	var ret []Func
	for _, item := range all {
		switch tmp := item.(type) {
		case Func:
			ret = append(ret, tmp)
		case func(next http.HandlerFunc) http.HandlerFunc: // alias for Func
			ret = append(ret, Func(tmp))
		case http.Handler:
			ret = append(ret, func(next http.HandlerFunc) http.HandlerFunc {
				return tmp.ServeHTTP
			})
		case func(http.ResponseWriter, *http.Request): // alias for http.HandlerFunc
			ret = append(ret, func(next http.HandlerFunc) http.HandlerFunc {
				return tmp
			})
		default:
			itemValue := reflect.ValueOf(item)
			switch itemValue.Type().Kind() {
			case reflect.Slice, reflect.Array:
				args := make([]interface{}, itemValue.Len())
				for i := 0; i < itemValue.Len(); i++ {
					args[i] = itemValue.Index(i).Interface()
				}
				ret = append(ret, BuildList(args...)...)
			default:
				panic("invalid argument")
			}
		}
	}
	return ret
}
