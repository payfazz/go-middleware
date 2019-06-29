package paniclogger_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/paniclogger"
	"github.com/payfazz/go-middleware/util/printer"
)

type checkWrite struct {
	written bool
}

func (cw *checkWrite) Write(data []byte) (int, error) {
	cw.written = true
	return len(data), nil
}

func TestNormal(t *testing.T) {
	cw := &checkWrite{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	defer func() {
		if rec := recover(); rec != nil {
			t.FailNow()
		}
	}()

	h := middleware.C(
		paniclogger.New(
			20,
			paniclogger.DefaultLogger(printer.Wrap(cw)),
		),
		func(w http.ResponseWriter, r *http.Request) {
		},
	)

	h(w, r)
}

func TestAbort(t *testing.T) {
	cw := &checkWrite{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	defer func() {
		if rec := recover(); rec != http.ErrAbortHandler {
			t.FailNow()
		}
	}()

	h := middleware.C(
		paniclogger.New(
			20,
			paniclogger.DefaultLogger(printer.Wrap(cw)),
		),
		func(w http.ResponseWriter, r *http.Request) {
			panic(http.ErrAbortHandler)
		},
	)

	h(w, r)
}

func TestPanic(t *testing.T) {
	cw := &checkWrite{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	defer func() {
		if rec := recover(); rec == nil {
			t.FailNow()
		}
		if !cw.written {
			t.FailNow()
		}
		if w.Code != 500 {
			t.FailNow()
		}
	}()

	h := middleware.C(
		paniclogger.New(
			20,
			paniclogger.DefaultLogger(printer.Wrap(cw)),
		),
		func(w http.ResponseWriter, r *http.Request) {
			panic(errors.New("test-panic"))
		},
	)

	h(w, r)
}

func TestPanic2(t *testing.T) {
	cw := &checkWrite{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	defer func() {
		if rec := recover(); rec == nil {
			t.FailNow()
		}
		if !cw.written {
			t.FailNow()
		}
		if w.Code != 0 {
			t.FailNow()
		}
		if len(w.Body.Bytes()) != 0 {
			t.FailNow()
		}
	}()

	h := middleware.C(
		paniclogger.New(
			20,
			paniclogger.DefaultLogger(printer.Wrap(cw)),
		),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "test")
			panic(errors.New("test-panic"))
		},
	)

	h(w, r)
}
