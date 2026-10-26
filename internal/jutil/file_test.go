package jutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Go-only: tests of java.io.File, Files.asCharSource, CharStreams.readLines
// and PrintStream(File) against results printed by JDK 21 / Guava.
func TestFileCharSourceReadLines(t *testing.T) {
	for in, want := range map[string]string{
		"":           "",
		"a":          "[a]",
		"a\n":        "[a]",
		"\n":         "[]",
		"a\nb":       "[a][b]",
		"a\r\nb\r\n": "[a][b]",
		"a\rb\r":     "[a][b]",
		"\r\n\r\n":   "[][]",
		"a\n\nb":     "[a][][b]",
		"\n\r":       "[][]",
		"x\r\r\ny":   "[x][][y]",
		"a b":        "[a b]",
	} {
		var sb strings.Builder
		for _, l := range CharStreams_ReadLines(NewStringReader(in)) {
			sb.WriteString("[" + l + "]")
		}
		if sb.String() != want {
			t.Errorf("readLines(%q) = %q, want %q", in, sb.String(), want)
		}
	}
	for in, want := range map[string]string{
		"a//b/": "a/b", "/": "/", "": "", "//x//": "/x",
		"a//b": "a/b", "a/": "a", "///": "/", "//": "/", "a///": "a", "./a/": "./a", "a/./b": "a/./b",
	} {
		if got := NewFile(in).GetPath(); got != want {
			t.Errorf("new File(%q).getPath() = %q, want %q", in, got, want)
		}
	}
	expectPanic := func(want string, f func()) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("no exception, want %s", want)
				return
			}
			if got := r.(Throwable).Error(); got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		}()
		f()
	}
	expectPanic("java.io.FileNotFoundException: /nonexistent/zz (No such file or directory)", func() {
		Files_AsCharSource(NewFile("/nonexistent/zz"), StandardCharsets_UTF_8).OpenStream()
	})
	expectPanic("java.io.FileNotFoundException: /nonexistent/zz (No such file or directory)", func() {
		NewPrintStreamFromFile(NewFile("/nonexistent/zz"))
	})
	dir := t.TempDir()
	expectPanic("java.io.FileNotFoundException: "+dir+" (Is a directory)", func() {
		Files_AsCharSource(NewFile(dir), StandardCharsets_UTF_8).OpenStream()
	})
	path := filepath.Join(dir, "out.txt")
	ps := NewPrintStreamFromFile(NewFile(path))
	ps.Println("xé")
	ps.Close()
	b, _ := os.ReadFile(path)
	if string(b) != "xé\n" {
		t.Errorf("PrintStream(File) wrote %q", b)
	}
	r := Files_AsCharSource(NewFile(path), StandardCharsets_UTF_8).OpenStream()
	if lines := CharStreams_ReadLines(r); len(lines) != 1 || lines[0] != "xé" {
		t.Errorf("read back %q", lines)
	}
	r.Close()
}

// Go-only: CharSource.read as Guava documents it.
func TestCharSourceWrapRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(path, []byte("xé\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := Files_AsCharSource(NewFile(path), StandardCharsets_UTF_8).Read(); got != "xé\n" {
		t.Errorf("Files.asCharSource(...).read() = %q", got)
	}
}
