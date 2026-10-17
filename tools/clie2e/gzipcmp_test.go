package main

import (
	"bytes"
	"compress/gzip"
	"testing"
)

func gzipped(t *testing.T, level int, content string) []byte {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, level)
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte(content))
	w.Close()
	out := b.Bytes()
	out[8], out[9] = 0, 0xff // the header GZIPOutputStream writes (XFL 0, OS unknown)
	return out
}

func TestSameOutput(t *testing.T) {
	text := "var x = 1;\n" + string(bytes.Repeat([]byte("abcdefgh"), 1000))
	fast, best := gzipped(t, gzip.BestSpeed, text), gzipped(t, gzip.BestCompression, text)
	if bytes.Equal(fast, best) {
		t.Fatal("want two different compressions of the same content")
	}
	if !sameOutput(fast, best) {
		t.Error("the same content compressed differently should compare equal")
	}
	if sameOutput(fast, gzipped(t, gzip.BestSpeed, text+"!")) {
		t.Error("different contents should differ")
	}
	if sameOutput([]byte("a\n"), []byte("b\n")) || !sameOutput([]byte("a\n"), []byte("a\n")) {
		t.Error("plain files are compared byte for byte")
	}
	broken := append([]byte{}, fast[:len(fast)-4]...)
	if sameOutput(broken, fast) {
		t.Error("a truncated gzip file should differ")
	}
}
