package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

func gz(t *testing.T, s string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write([]byte(s)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// zipOf returns a zip archive holding the given files, in order.
func zipOf(t *testing.T, files ...[2]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, f := range files {
		w, err := zw.Create(f[0])
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(f[1])); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// gplWord is the word the check looks for, assembled so that this file does
// not contain it.
const gplWord = "G" + "PL"

func TestFileText(t *testing.T) {
	multi := append(gz(t, "first member, "), gz(t, "second member")...)
	// Compressed bytes in which the word appears by chance: a gzip header
	// whose stored file name is the word, over text that never says it.
	var chance bytes.Buffer
	zw := gzip.NewWriter(&chance)
	zw.Name = "x " + strings.ToLower(gplWord) + " y"
	zw.Write([]byte("plain text"))
	zw.Close()
	for _, tc := range []struct {
		name     string
		in       []byte
		want     string // the text; "" with skip or err
		skip     bool
		err      string // substring of the error
		mentions bool   // mentionsGPL of the text
	}{
		{name: "a.js", in: []byte("var x = 1;\n"), want: "var x = 1;\n"},
		{name: "a.txt", in: []byte("under the " + gplWord + " v2"), want: "under the " + gplWord + " v2", mentions: true},
		{name: "a.txt.gz", in: gz(t, "text, licensed under the "+gplWord), want: "text, licensed under the " + gplWord, mentions: true},
		{name: "b.txt.gz", in: gz(t, "nothing to see"), want: "nothing to see"},
		{name: "multi.gz", in: multi, want: "first member, second member"},
		{name: "chance.gz", in: chance.Bytes(), want: "plain text"},
		{name: "no-suffix", in: gz(t, "gzip without .gz"), want: "gzip without .gz"},
		{name: "c.gz", in: []byte("not gzip"), err: "named like a gzip file"},
		{name: "d.gz", in: gz(t, "truncated text")[:12], err: "d.gz"},
		{name: "e.zip", in: []byte("PK\x03\x04rest"), err: "zip file"},
		{name: "externs.zip", in: zipOf(t, [2]string{"a.js", "var a;"}, [2]string{"dir/b.js", "// under the " + gplWord}), want: "var a;\n// under the " + gplWord + "\n", mentions: true},
		{name: "nested.jar", in: zipOf(t, [2]string{"x.txt.gz", string(gz(t, "inner text"))}, [2]string{"y.bin", "\x00\xff\xfe"}), want: "inner text\n"},
		{name: "f.bin", in: []byte("PK\x03\x04rest"), err: "zip file"},
		{name: "g.jar", in: []byte("anything"), err: "jar file"},
		{name: "h.dat", in: []byte("BZh91AY&SY"), err: "bzip2 file"},
		{name: "i.dat", in: []byte("\xfd7zXZ\x00"), err: "xz file"},
		{name: "j.dat", in: []byte("\x28\xb5\x2f\xfd"), err: "zstd file"},
		{name: "k.dat", in: []byte("7z\xbc\xaf\x27\x1c"), err: "7z file"},
		{name: "l.txt", in: []byte("BZh is not bzip2 without a block size"), want: "BZh is not bzip2 without a block size"},
		{name: "image.png", in: []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\xff " + gplWord + " \xfe"), skip: true},
		{name: "nul.js", in: []byte("'\x00'; // " + gplWord), want: "'\x00'; // " + gplWord, mentions: true},
		{name: "latin1.xml", in: []byte("<a>\xe9</a>"), want: "<a>\xe9</a>"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			text, skip, err := fileText(tc.name, tc.in)
			switch {
			case tc.err != "":
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("err = %v, want one containing %q", err, tc.err)
				}
				return
			case err != nil:
				t.Fatal(err)
			case skip != tc.skip:
				t.Fatalf("skip = %v, want %v", skip, tc.skip)
			case string(text) != tc.want:
				t.Fatalf("text = %q, want %q", text, tc.want)
			}
			if got := mentionsGPL(text); got != tc.mentions {
				t.Errorf("mentionsGPL(%q) = %v, want %v", text, got, tc.mentions)
			}
		})
	}
}

func TestMentionsGPL(t *testing.T) {
	for _, tc := range []struct {
		text string
		want bool
	}{
		{"Licensed under the Apache License, Version 2.0", false},
		{"GNU General Public " + "License", true},
		{"gnu general public " + "licence", true},
		{"Lesser " + strings.ToLower(gplWord) + " here", true},
		{"L" + gplWord + "-2.1", true},
		{"an " + gplWord + "Excluded identifier", false},
		{"x" + strings.ToLower(gplWord) + "y", false},
		{"public licenses in general", false},
		{"see tools/oracle/" + corpus.ListName + " for the list", false},
		{corpus.ListName + " " + gplWord, true},
	} {
		if got := mentionsGPL([]byte(tc.text)); got != tc.want {
			t.Errorf("mentionsGPL(%q) = %v, want %v", tc.text, got, tc.want)
		}
	}
}

func TestExcludedIndex(t *testing.T) {
	x := newExcludedIndex([]string{
		"m@1.0.0/package/lib/esm/sparse/csA.js",
		"m@1.0.0/package/lib/esm/sparse/csB.js",
		"m@1.0.0/package/lib/cjs/sparse/csA.js",
		"loose.js",
	})
	for _, tc := range []struct {
		name string
		in   []byte
		want []string
	}{
		{"none", []byte("lib/esm/sparse/csC.js lib/esm/sparse/"), nil},
		{"list entry", []byte("# file m@1.0.0/package/lib/esm/sparse/csB.js\n"), []string{"m@1.0.0/package/lib/esm/sparse/csB.js"}},
		{"other root", []byte("/c/package/lib/cjs/sparse/csA.js and lib/esm/sparse/csA.js"), []string{"m@1.0.0/package/lib/esm/sparse/csA.js", "m@1.0.0/package/lib/cjs/sparse/csA.js"}},
		{"no package dir", []byte("x/loose.js"), []string{"loose.js"}},
		{"gzip golden", gz(t, "# file m@1.0.0/package/lib/esm/sparse/csA.js\nSCRIPT\n"), []string{"m@1.0.0/package/lib/esm/sparse/csA.js"}},
	} {
		text, _, err := fileText(tc.name, tc.in)
		if err != nil {
			t.Fatal(err)
		}
		if got := x.mentioned(text); !slices.Equal(got, tc.want) {
			t.Errorf("%s: mentioned = %q, want %q", tc.name, got, tc.want)
		}
	}
}
