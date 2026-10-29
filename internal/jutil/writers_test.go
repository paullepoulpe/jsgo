package jutil

import (
	"bytes"
	"fmt"
	"testing"
)

// Expected bytes come from the JDK 21 runtime (OutputStreamWriter over a
// ByteArrayOutputStream, UTF_8).

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
