package kv_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/util/kv"
)

func TestNormal(t *testing.T) {
	type somekeytype struct{}
	var somekey somekeytype

	type someotherkeytype struct{}
	var someotherkey someotherkeytype

	all := middleware.Chain(
		func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				next(w, kv.WithValue(r, somekey, "test-value"))
			}
		},
		kv.Injector(20, "someint"),
		func(w http.ResponseWriter, r *http.Request) {
			val1, ok1 := kv.Get(r, somekey)
			val2, ok2 := kv.Get(r, someotherkey)
			val3 := kv.MustGet(r, 20)
			fmt.Fprintf(w, "%v:%v|%v:%v|true:%v", ok1, val1, ok2, val2, val3)
		},
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	all(w, r)

	expected := "true:test-value|false:<nil>|true:someint"
	found := w.Body.String()

	if found != expected {
		t.Fatalf("found '%s', need '%s'", found, expected)
	}
}

func TestMustGet(t *testing.T) {
	gotPanic := false
	type somekeytype struct{}
	var somekey somekeytype

	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, kv.MustGet(r, somekey).(string))
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	func() {
		defer func() { gotPanic = recover() != nil }()
		h(w, r)
	}()

	if !gotPanic {
		t.Errorf("invalid kv.MustGet")
	}
}
