package jutil

import "testing"

type failingWriter struct{ writes int }

func (w *failingWriter) Write(string) { w.writes++; panic(NewIOException("boom")) }
func (w *failingWriter) Flush()       {}
func (w *failingWriter) Close()       {}

// Go-only: PrintWriter, StringWriter and File.exists as their Javadoc
// specifies them.
func TestPrintWriterStringWriter(t *testing.T) {
	sw := NewStringWriter()
	pw := NewPrintWriter(sw)
	pw.Print("a")
	pw.PrintChar(' ')
	pw.Println("b")
	pw.PrintNullable(nil)
	pw.PrintlnEmpty()
	pw.Flush()
	if got, want := sw.String(), "a b\nnull\n"; got != want {
		t.Errorf("StringWriter = %q, want %q", got, want)
	}
	if pw.CheckError() {
		t.Errorf("CheckError() = true without an error")
	}
	pw.Close()
	pw.Close()
	pw.Print("x")
	if !pw.CheckError() || sw.String() != "a b\nnull\n" {
		t.Errorf("write after close: CheckError() = %v, buffer %q", pw.CheckError(), sw.String())
	}
	sw.Write("c")
	if sw.String() != "a b\nnull\nc" {
		t.Errorf("StringWriter after close = %q", sw.String())
	}

	fw := &failingWriter{}
	pw = NewPrintWriter(fw)
	pw.Println("x")
	if !pw.CheckError() || fw.writes != 2 {
		t.Errorf("failing writer: CheckError() = %v, writes = %d", pw.CheckError(), fw.writes)
	}

	if NewFile("").Exists() || NewFile("/nonexistent-dir-for-test/x").Exists() || !NewFile("/").Exists() || !NewFile(".").Exists() {
		t.Errorf("File.exists() is wrong")
	}
}
