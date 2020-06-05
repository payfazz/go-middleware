package kv_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
)

func TestNormal(t *testing.T) {
	type somekeytype struct{}
	var somekey somekeytype

	type someotherkeytype struct{}
	var someotherkey someotherkeytype

	h := middleware.C(
		func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				next(w, kv.EnsureKVAndSet(r, somekey, "test-value"))
			}
		},
		func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				next(w, kv.EnsureKVAndSet(r, "stringkey", "stringvalue"))
			}
		},
		kv.Injector(20, "someint"),
		func(w http.ResponseWriter, r *http.Request) {
			val1, ok1 := kv.Get(r, somekey)
			val2, ok2 := kv.Get(r, someotherkey)
			val3, ok3 := kv.Get(r, "stringkey")
			val4 := kv.MustGet(r, 20)
			fmt.Fprintf(w, "%v:%v|%v:%v|%v:%v|true:%v", ok1, val1, ok2, val2, ok3, val3, val4)
		},
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h(w, r)

	expected := "true:test-value|false:<nil>|true:stringvalue|true:someint"
	found := w.Body.String()

	if found != expected {
		t.Fatalf("found '%s', need '%s'", found, expected)
	}
}

func TestMustGet(t *testing.T) {
	defer func() {
		if rec := recover(); rec == nil {
			t.Fatalf("MustGet must panic")
		}
	}()
	type somekeytype struct{}
	var somekey somekeytype

	h := middleware.C(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, kv.MustGet(r, somekey).(string))
		},
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h(w, r)
}
