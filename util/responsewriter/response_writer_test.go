package responsewriter_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/payfazz/go-middleware/util/responsewriter"
)

func TestNormal(t *testing.T) {
	w := httptest.NewRecorder()
	newW := responsewriter.Wrap(w)

	newW.Write([]byte("data"))

	if w.Code != 200 {
		t.FailNow()
	}
	if w.Body.String() != "data" {
		t.FailNow()
	}
}

func TestStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	newW := responsewriter.Wrap(w)

	newW.WriteHeader(404)
	newW.Write([]byte("data"))

	if w.Code != 404 {
		t.FailNow()
	}
	if w.Body.String() != "data" {
		t.FailNow()
	}
}

func TestHeader(t *testing.T) {
	w := httptest.NewRecorder()
	newW := responsewriter.Wrap(w)

	newW.Header().Set("X-Test", "x-test")
	newW.Write([]byte("data"))

	if w.Code != 200 {
		t.FailNow()
	}
	if !reflect.DeepEqual(w.HeaderMap, http.Header{"X-Test": {"x-test"}}) {
		t.FailNow()
	}
	if w.Body.String() != "data" {
		t.FailNow()
	}
}
