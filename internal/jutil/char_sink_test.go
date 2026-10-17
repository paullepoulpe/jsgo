package jutil

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilCharSinkVectors.java.

func charSinkUnits(t *testing.T, h string) string {
	t.Helper()
	if h == "-" {
		return ""
	}
	var u []uint16
	for i := 0; i < len(h); i += 4 {
		v, err := strconv.ParseUint(h[i:i+4], 16, 16)
		if err != nil {
			t.Fatal(err)
		}
		u = append(u, uint16(v))
	}
	return StringFromUTF16(u)
}

func TestCharSinkWriteMatchesGuava(t *testing.T) {
	rows := readTSV(t, "char_sink.tsv")
	dir := t.TempDir()
	n := 0
	for i, r := range rows {
		switch r[0] {
		case "dir", "missing":
			target := dir
			if r[0] == "missing" {
				target = filepath.Join(dir, "missing", "f")
			}
			got := "none"
			Try(func() {
				Files_AsCharSink(NewFile(target), StandardCharsets_UTF_8).Write("x")
			}, func(e IIOException) {
				got = strings.Replace(e.(Throwable).Error(), dir, "<dir>", 1)
			})
			if got != r[1] {
				t.Errorf("%s: got %q, want %q", r[0], got, r[1])
			}
		default:
			in := charSinkUnits(t, r[0])
			f := filepath.Join(dir, "f"+strconv.Itoa(i))
			Files_AsCharSink(NewFile(f), StandardCharsets_UTF_8).Write("previous longer content\n")
			Files_AsCharSink(NewFile(f), StandardCharsets_UTF_8).Write(in)
			b, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			got := hex.EncodeToString(b)
			if got == "" {
				got = "-"
			}
			if got != r[1] {
				t.Errorf("write(%s): file %s, want %s", r[0], got, r[1])
			}
			n++
		}
	}
	if n < 12 {
		t.Fatalf("only %d vectors", n)
	}
}
