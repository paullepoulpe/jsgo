package jutil

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// Expected bytes come from the JDK 21 runtime (OutputStreamWriter and
// BufferedWriter over a ByteArrayOutputStream, UTF_8).

type nopCloseBuffer struct{ bytes.Buffer }

func (*nopCloseBuffer) Close() error { return nil }

func TestOutputStreamWriterSurrogatesMatchJava(t *testing.T) {
	var out nopCloseBuffer
	w := NewOutputStreamWriter(&out)
	for _, s := range [][]uint16{{'a', 0xd83d}, {0xde00, 'b'}, {0xdc00}, {'c', 0xd800}, {'d'}, {0xd83d}} {
		w.Write(StringFromUTF16(s))
	}
	w.Close()
	if got := fmt.Sprintf("%x", out.Bytes()); got != "61f09f9880623f633f643f" {
		t.Errorf("got %s", got)
	}
	if got := throwableString(func() { w.Write("x") }); got != "java.io.IOException: Stream closed" {
		t.Errorf("write after close: %s", got)
	}
	w.Close() // no effect
}

func TestBufferedWriterMatchesJava(t *testing.T) {
	var out nopCloseBuffer
	bw := NewBufferedWriter(NewOutputStreamWriter(&out))
	// The 8192-char buffer fills between the two halves of the pair.
	bw.Write(strings.Repeat("a", 8191) + "\U0001F600z")
	bw.Close()
	b := out.Bytes()
	if got := fmt.Sprintf("%d %x", len(b), b[8189:]); got != "8196 6161f09f98807a" {
		t.Errorf("got %s", got)
	}
	if got := throwableString(func() { bw.Write("x") }); got != "java.io.IOException: Stream closed" {
		t.Errorf("write after close: %s", got)
	}
	if _, err := bw.WriteString("x"); err == nil || err.Error() != "Stream closed" {
		t.Errorf("WriteString after close: %v", err)
	}
	bw.Close() // no effect
}

func TestFilesNewOutputStreamOptions(t *testing.T) {
	p := Path_Of(t.TempDir(), "f.txt")
	if got := throwableString(func() {
		Files_NewOutputStream(p, StandardOpenOption_APPEND, StandardOpenOption_TRUNCATE_EXISTING)
	}); got != "java.lang.IllegalArgumentException: APPEND + TRUNCATE_EXISTING not allowed" {
		t.Errorf("got %s", got)
	}
	if got := throwableString(func() { Files_CreateDirectories(nil) }); !strings.HasPrefix(got, "java.lang.NullPointerException: Cannot invoke") {
		t.Errorf("got %s", got)
	}
}
