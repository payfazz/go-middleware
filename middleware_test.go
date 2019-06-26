package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/payfazz/go-middleware"
)

func genM(text string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, text)
			next(w, r)
		}
	}
}
func TestFunc(t *testing.T) {
	h := middleware.Compile(
		genM("1"),
		genM("2"),
		[]interface{}{
			genM("3"),
			genM("4"),
			[]interface{}{
				genM("5"),
			},
		},
		genM("6"),
		middleware.Nop,
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "L")
		},
	)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h(w, r)

	expected := "123456L"
	found := w.Body.String()

	if found != expected {
		t.Fatalf("expected '%s', found '%s'", expected, found)
	}
}

func TestHandlerFunc(t *testing.T) {
	h := middleware.Compile(
		genM("1"),
		genM("2"),
		[]interface{}{
			genM("3"),
			genM("4"),
			[]interface{}{
				genM("5"),
			},
		},
		genM("6"),
		middleware.Nop,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "L")
		}),
	)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h(w, r)

	expected := "123456L"
	found := w.Body.String()

	if found != expected {
		t.Fatalf("expected '%s', found '%s'", expected, found)
	}
}

func TestHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "L")
	})
	h := middleware.Compile(
		genM("1"),
		genM("2"),
		[]interface{}{
			genM("3"),
			genM("4"),
			[]interface{}{
				genM("5"),
			},
		},
		genM("6"),
		middleware.Nop,
		mux,
	)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h(w, r)

	expected := "123456L"
	found := w.Body.String()

	if found != expected {
		t.Fatalf("expected '%s', found '%s'", expected, found)
	}
}

func TestNil(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			return
		}
		t.Fatal("middleware.Compile should not be able to process nil")
	}()
	middleware.Compile(nil)
}

func TestType(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			return
		}
		t.Fatal("middleware.Compile should not be able to unknown type")
	}()
	middleware.Compile(1, 2, 3)
}
