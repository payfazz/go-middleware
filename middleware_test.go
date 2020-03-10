package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/payfazz/go-middleware"
)

func genMFF(text string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, text)
			next(w, r)
		}
	}
}
func TestFunc(t *testing.T) {
	h := middleware.C(
		genMFF("1"),
		genMFF("2"),
		[]interface{}{
			genMFF("3"),
			genMFF("4"),
			[]interface{}{
				genMFF("5"),
			},
		},
		genMFF("6"),
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
	h := middleware.C(
		genMFF("1"),
		genMFF("2"),
		[]interface{}{
			genMFF("3"),
			genMFF("4"),
			[]interface{}{
				genMFF("5"),
			},
		},
		genMFF("6"),
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
	h := middleware.C(
		genMFF("1"),
		genMFF("2"),
		[]interface{}{
			genMFF("3"),
			genMFF("4"),
			[]interface{}{
				genMFF("5"),
			},
		},
		genMFF("6"),
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
		t.Fatal("middleware.C should not be able to process nil")
	}()
	middleware.C(nil)
}

func TestArbitrary(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			return
		}
		t.Fatal("middleware.C should not be able to arbitrary type")
	}()
	middleware.C(1, 2, 3)
}

func getMHH(text string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, text)
			next.ServeHTTP(w, r)
		})
	}
}

func TestMHH(t *testing.T) {
	h := middleware.Compile(
		getMHH("1"),
		getMHH("2"),
		[]interface{}{
			getMHH("3"),
			getMHH("4"),
			[]interface{}{
				getMHH("5"),
			},
		},
		getMHH("6"),
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

func genMHF(text string) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, text)
			next.ServeHTTP(w, r)
		}
	}
}

func TestHF(t *testing.T) {
	h := middleware.Compile(
		genMHF("1"),
		genMHF("2"),
		[]interface{}{
			genMHF("3"),
			genMHF("4"),
			[]interface{}{
				genMHF("5"),
			},
		},
		genMHF("6"),
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

func genMFH(text string) func(http.HandlerFunc) http.Handler {
	return func(next http.HandlerFunc) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, text)
			next(w, r)
		})
	}
}

func TestMFH(t *testing.T) {
	h := middleware.Compile(
		genMFH("1"),
		genMFH("2"),
		[]interface{}{
			genMFH("3"),
			genMFH("4"),
			[]interface{}{
				genMFH("5"),
			},
		},
		genMFH("6"),
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
