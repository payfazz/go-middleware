// Package printer for Printer interface
package printer

// Printer interface
type Printer interface {
	Print(v ...interface{})
}

// Func is type adapter for Printer interface
type Func func(v ...interface{})

// Print call f(v...)
func (f Func) Print(v ...interface{}) {
	f(v...)
}
