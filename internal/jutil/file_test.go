package jutil

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
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
		for _, l := range splitLinesUTF16(StringToUTF16(in)) {
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
		Files_AsCharSource(NewFile("/nonexistent/zz"), StandardCharsets_UTF_8).Read()
	})
	expectPanic("java.io.FileNotFoundException: /nonexistent/zz (No such file or directory)", func() {
		NewPrintStreamFromFile(NewFile("/nonexistent/zz"))
	})
	dir := t.TempDir()
	expectPanic("java.io.FileNotFoundException: "+dir+" (Is a directory)", func() {
		Files_AsCharSource(NewFile(dir), StandardCharsets_UTF_8).Read()
	})
	path := filepath.Join(dir, "out.txt")
	ps := NewPrintStreamFromFile(NewFile(path))
	ps.Println("xé")
	ps.Close()
	b, _ := os.ReadFile(path)
	if string(b) != "xé\n" {
		t.Errorf("PrintStream(File) wrote %q", b)
	}
	lines, ex := Files_AsCharSource(NewFile(path), StandardCharsets_UTF_8).TryReadLines()
	if ex != nil || len(lines) != 1 || lines[0] != "xé" {
		t.Errorf("read back %q (err=%v)", lines, ex)
	}
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

func TestPathHashCodeConcurrent(t *testing.T) {
	p := Path_Of("some/concurrent/path/foo.js")
	var wg sync.WaitGroup
	results := make([]int32, 16)
	for i := range results {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = p.HashCode()
		}(i)
	}
	wg.Wait()
	for i, got := range results {
		if got == 0 || got != results[0] {
			t.Fatalf("results[%d] = %d, want %d", i, got, results[0])
		}
	}
	if root := Path_Of("/"); root.HashCode() != '/' {
		t.Fatalf("root HashCode() = %d, want %d", root.HashCode(), '/')
	}
	if empty := Path_Of(""); empty.HashCode() != 0 {
		t.Fatalf("empty HashCode() = %d, want 0", empty.HashCode())
	}
}
