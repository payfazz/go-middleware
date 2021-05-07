package responsewriter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/payfazz/go-middleware/util/responsewriter"
)

func TestA(t *testing.T) {
	res := httptest.NewRecorder()
	h := func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "test") }

	wrappedRes := responsewriter.Wrap(res)
	if responsewriter.Wrap(wrappedRes) != wrappedRes {
		t.Errorf("Wrap should be indempotent")
	}

	h(wrappedRes, nil)

	if !wrappedRes.Written() {
		t.Errorf("wrappedRes.Written should be true")
	}

	if wrappedRes.Status() != 200 {
		t.Errorf("wrappedRes.Status should be 200")
	}

	if wrappedRes.Size() != 4 {
		t.Errorf("wrappedRes.Size should be 4")
	}

	if res.Body.String() != "test" {
		t.Errorf(`res.Body should be "test"`)
	}
}
