package javazip_test

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javazip"
)

func TestZipInputStreamAndZipFile(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w1, err := zw.Create("a.js")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w1.Write([]byte("var a = 1;"))
	w2, err := zw.Create("dir/b.js")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w2.Write([]byte("var b = 2;"))
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	zis := javazip.NewZipInputStream(bytes.NewReader(buf.Bytes()))
	defer zis.Close()

	e1, err := zis.GetNextEntry()
	if err != nil || e1 == nil || e1.GetName() != "a.js" {
		t.Fatalf("entry 1 mismatch: %v, err=%v", e1, err)
	}
	b1, err := javaio.ReadAllBytes(zis)
	if err != nil || string(b1) != "var a = 1;" {
		t.Fatalf("read entry 1 got %q, err=%v", string(b1), err)
	}

	e2, err := zis.GetNextEntry()
	if err != nil || e2 == nil || e2.GetName() != "dir/b.js" {
		t.Fatalf("entry 2 mismatch: %v, err=%v", e2, err)
	}
	b2, err := javaio.ReadAllBytes(zis)
	if err != nil || string(b2) != "var b = 2;" {
		t.Fatalf("read entry 2 got %q, err=%v", string(b2), err)
	}

	e3, err := zis.GetNextEntry()
	if err != nil || e3 != nil {
		t.Fatalf("expected nil EOF entry, got %v (err=%v)", e3, err)
	}

	// Test ZipFile on disk
	zipPath := filepath.Join(t.TempDir(), "test.zip")
	if err := os.WriteFile(zipPath, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	zf, err := javazip.NewZipFile(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer zf.Close()
	if zf.Size() != 2 {
		t.Fatalf("expected 2 entries, got %d", zf.Size())
	}
	ent := zf.GetEntry("dir/b.js")
	if ent == nil {
		t.Fatal("missing dir/b.js in ZipFile")
	}
	in, err := zf.GetInputStream(ent)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	raw, err := io.ReadAll(in)
	if err != nil || string(raw) != "var b = 2;" {
		t.Fatalf("ZipFile entry read got %q (err=%v)", string(raw), err)
	}
}

func TestGZIPRoundTrip(t *testing.T) {
	out := javaio.NewByteArrayOutputStream()
	gzOut, err := javazip.NewFastGzipOutputStream(out)
	if err != nil {
		t.Fatal(err)
	}
	payload := "closure typedast payload\nwith multiple lines\n"
	if _, err := gzOut.Write([]byte(payload)); err != nil {
		t.Fatal(err)
	}
	if err := gzOut.Close(); err != nil {
		t.Fatal(err)
	}

	gzIn, err := javazip.NewGZIPInputStreamWithSize(javaio.NewByteArrayInputStream(out.ToByteArray()), 4096)
	if err != nil {
		t.Fatal(err)
	}
	defer gzIn.Close()
	got, err := javaio.ReadAllBytes(gzIn)
	if err != nil || string(got) != payload {
		t.Fatalf("GZIP roundtrip got %q (err=%v)", string(got), err)
	}
}
