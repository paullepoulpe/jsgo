package jutil

// Go-only: language-adapter — java.io.Writer. An IOException is thrown as
// a panic (jutil.IOException), like every other Java exception.

// Writer is java.io.Writer (the methods ported code uses).
type Writer interface {
	// Write is Writer.write(String).
	Write(s string)
	// Flush is Writer.flush().
	Flush()
	// Close is Writer.close().
	Close()
}
