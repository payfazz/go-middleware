package middleware

import (
	"net/http"
	"reflect"
)

// Compile will flatten all params into single list of middleware, basically:
// 	[m1, m2, [m3, m4, [m5, m6]], m7] -> [m1, m2, m3, m4, m5, m6, m7]
// and also will convert "http.HandlerFunc" and "http.Handler" into middleware that
// doesn't call next middleware, i.e. stopping the chain,
//
// 	middleware is a value with type:
// 	- func(http.HandlerFunc) http.HandlerFunc
// 	- func(http.Handler)     http.Handlerfunc
// 	- func(http.HandlerFunc) http.Handler
// 	- func(http.Handler)     http.Handler
//
// Compile will create middleware chain based on the flattened list.
func Compile(all ...interface{}) http.HandlerFunc {
	var f http.HandlerFunc
	list := compileList(all...)
	for i := len(list) - 1; i >= 0; i-- {
		f = list[i](f)
	}
	return f
}

// C is same with Compile, it is just shortcut
func C(all ...interface{}) http.HandlerFunc {
	return Compile(all...)
}

func compileList(all ...interface{}) []func(http.HandlerFunc) http.HandlerFunc {
	ret := make([]func(http.HandlerFunc) http.HandlerFunc, 0, len(all))
	for _, item := range all {
		if item == nil {
			panic("middleware: invalid argument: can't be nil")
		}

		itemValue := reflect.ValueOf(item)
		itemType := itemValue.Type()

		var mff func(http.HandlerFunc) http.HandlerFunc
		mffValue := reflect.ValueOf(&mff).Elem()
		mffType := mffValue.Type()
		if itemType.ConvertibleTo(mffType) {
			mffValue.Set(itemValue.Convert(mffType))
			ret = append(ret, mff)
			continue
		}

		var mhh func(http.Handler) http.Handler
		mhhValue := reflect.ValueOf(&mhh).Elem()
		mhhType := mhhValue.Type()
		if itemType.ConvertibleTo(mhhType) {
			mhhValue.Set(itemValue.Convert(mhhType))
			ret = append(ret, func(next http.HandlerFunc) http.HandlerFunc {
				return mhh(next).ServeHTTP
			})
			continue
		}

		var mhf func(http.Handler) http.HandlerFunc
		mhfValue := reflect.ValueOf(&mhf).Elem()
		mhfType := mhfValue.Type()
		if itemType.ConvertibleTo(mhfType) {
			mhfValue.Set(itemValue.Convert(mhfType))
			ret = append(ret, func(next http.HandlerFunc) http.HandlerFunc {
				return mhf(next)
			})
			continue
		}

		var mfh func(http.HandlerFunc) http.Handler
		mfhValue := reflect.ValueOf(&mfh).Elem()
		mfhType := mfhValue.Type()
		if itemType.ConvertibleTo(mfhType) {
			mfhValue.Set(itemValue.Convert(mfhType))
			ret = append(ret, func(next http.HandlerFunc) http.HandlerFunc {
				return mfh(next).ServeHTTP
			})
			continue
		}

		var f http.HandlerFunc
		fValue := reflect.ValueOf(&f).Elem()
		fType := fValue.Type()
		if itemType.ConvertibleTo(fType) {
			fValue.Set(itemValue.Convert(fType))
			ret = append(ret, func(http.HandlerFunc) http.HandlerFunc {
				return f
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
			ret = append(ret, compileList(args...)...)
		default:
			panic("middleware: invalid argument: " + itemType.String() + " can't be converted middleware")
		}
	}
	return ret
}

// Nop is dummy middleware.
//
// Useful for something like:
// 	func Logger() func(http.HandlerFunc) http.HandlerFunc {
// 		if LOG_DISABLED {
// 			return middleware.Nop
// 		}
// 		return logger.New(nil)
// 	}
//
// NOTE: middleware chain only compiled once, so it doesn't have effect
// if you change "LOG_DISABLED" after it compiled
func Nop(next http.HandlerFunc) http.HandlerFunc {
	return next
}
