// Package printer for Printer interface
package printer

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

// Printer interface
type Printer interface {
	Print(v ...interface{})
}

type printer struct {
	sync.Mutex
	io.Writer
}

func (p *printer) Print(v ...interface{}) {
	var s string
	if len(v) == 1 {
		if s2, ok := v[0].(string); ok {
			s = s2
			if s[len(s)-1] != '\n' {
				s += "\n"
			}
		}
	}
	if s == "" {
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprint(v...))
		if sb.String()[sb.Len()-1] != '\n' {
			sb.WriteByte('\n')
		}
		s = sb.String()
	}

	// peform unsafe zero-copy conversion from string to byte slice
	sHeader := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bytes := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: sHeader.Data,
		Len:  sHeader.Len,
		Cap:  sHeader.Len,
	}))

	p.Lock()
	p.Write(bytes) // safe because Write only *read* the content
	p.Unlock()

	// make sure s live until here
	runtime.KeepAlive(s)
}

// Wrap io.Writer into Printer.
//
// Print method is safe to call concurently
func Wrap(writer io.Writer) Printer {
	return &printer{Writer: writer}
}
