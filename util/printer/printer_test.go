package printer_test

import (
	"bytes"
	"testing"

	"github.com/payfazz/go-middleware/util/printer"
)

func TestPrintSingle(t *testing.T) {
	b := &bytes.Buffer{}
	p := printer.Wrap(b)
	p.Print("test")

	expected := "test"
	found := b.String()
	if found != expected {
		t.Fatalf("expected '%#v', found '%#v'", expected, found)
	}
}

func TestPrintMultiple(t *testing.T) {
	b := &bytes.Buffer{}
	p := printer.Wrap(b)
	p.Print("test", "lala")

	expected := "testlala"
	found := b.String()
	if found != expected {
		t.Fatalf("expected '%#v', found '%#v'", expected, found)
	}
}
