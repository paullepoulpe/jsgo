// Clean-room specification based on java.io.PrintStream Javadoc.

package javaio

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var (
	atExitMu    sync.Mutex
	atExitHooks []func()
)

// RegisterAtExit registers a callback to run before System.exit terminates the process.
func RegisterAtExit(fn func()) {
	atExitMu.Lock()
	defer atExitMu.Unlock()
	atExitHooks = append(atExitHooks, fn)
}

// RunAtExit runs and clears registered exit hooks in LIFO order.
func RunAtExit() {
	atExitMu.Lock()
	hooks := atExitHooks
	atExitHooks = nil
	atExitMu.Unlock()
	for i := len(hooks) - 1; i >= 0; i-- {
		hooks[i]()
	}
}

// SystemExit flushes standard streams, runs exit hooks, and calls os.Exit.
func SystemExit(code int32) {
	System_Out.Flush()
	System_Err.Flush()
	RunAtExit()
	if code < 0 || code > 255 {
		os.Exit(int(uint8(code)))
	}
	os.Exit(int(code))
}

// PrintStream adds functionality to another output stream to print representations
// of various data values conveniently.
type PrintStream struct {
	w        io.Writer
	hasError bool
}

// NewPrintStream creates a new print stream.
func NewPrintStream(w io.Writer) *PrintStream {
	return &PrintStream{w: w}
}

// Writer returns the underlying io.Writer.
func (p *PrintStream) Writer() io.Writer {
	if p == nil {
		return nil
	}
	return p.w
}

// Write implements io.Writer on PrintStream (since PrintStream extends FilterOutputStream).
func (p *PrintStream) Write(b []byte) (int, error) {
	if p == nil || p.w == nil {
		return 0, nil
	}
	n, err := p.w.Write(b)
	if err != nil {
		p.hasError = true
	}
	return n, err
}

// Print prints a value to the stream.
func (p *PrintStream) Print(v any) {
	if p == nil || p.w == nil {
		return
	}
	s := javalang.ToJavaString(v)
	_, err := io.WriteString(p.w, s)
	if err != nil {
		p.hasError = true
	}
}

// Println prints a value followed by a newline.
func (p *PrintStream) Println(v any) {
	p.Print(v)
	if p != nil && p.w != nil {
		_, err := io.WriteString(p.w, "\n")
		if err != nil {
			p.hasError = true
		}
	}
}

// Printf writes a formatted string using the specified format string and arguments.
func (p *PrintStream) Printf(format string, args ...any) {
	if p == nil || p.w == nil {
		return
	}
	_, err := fmt.Fprintf(p.w, format, args...)
	if err != nil {
		p.hasError = true
	}
}

// Format writes a formatted string using the specified format string and arguments.
func (p *PrintStream) Format(format string, args ...any) *PrintStream {
	p.Printf(format, args...)
	return p
}

// Flush flushes the stream.
func (p *PrintStream) Flush() {
	if f, ok := p.w.(interface{ Flush() error }); ok {
		_ = f.Flush()
	}
}

// Close closes the stream.
func (p *PrintStream) Close() {
	if c, ok := p.w.(io.Closer); ok {
		_ = c.Close()
	}
}

// CheckError flushes the stream and checks its error state.
func (p *PrintStream) CheckError() bool {
	p.Flush()
	return p.hasError
}

// System_Out is the standard output stream.
var System_Out = NewPrintStream(os.Stdout)

// System_Err is the standard error output stream.
var System_Err = NewPrintStream(os.Stderr)
