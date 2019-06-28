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
		kv.New(),
		func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				kv.Set(r, somekey, "test-value")
				next(w, r)
			}
		},
		func(w http.ResponseWriter, r *http.Request) {
			val1, ok1 := kv.Get(r, somekey)
			val2, ok2 := kv.Get(r, someotherkey)
			fmt.Fprintf(w, "%v:%v|%v:%v", ok1, val1, ok2, val2)
		},
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h(w, r)

	expected := "true:test-value|false:<nil>"
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
		kv.New(),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, kv.MustGet(r, somekey).(string))
		},
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h(w, r)
}
func TestWrapper(t *testing.T) {
	type somekeytype struct{}
	var somekey somekeytype

	expected := "test-value"

	h := middleware.C(
		func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				r = kv.WrapRequest(r)
				kv.Set(r, somekey, expected)
				next(w, r)
			}
		},
		func(w http.ResponseWriter, r *http.Request) {
			r = kv.WrapRequest(r)
			fmt.Fprint(w, kv.MustGet(r, somekey).(string))
		},
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h(w, r)

	found := w.Body.String()

	if found != expected {
		t.Fatalf("found '%s', need '%s'", found, expected)
	}
}
