package jutil

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// The charset vectors are written by tools/oracle/java/JutilCharsetVectors.java
// (see there for the rows). JUTIL_CHARSET_VECTORS=<dir> reads vectors
// written with another seed instead of testdata.

func readCharsetVectors(t *testing.T) [][]string {
	t.Helper()
	dir := os.Getenv("JUTIL_CHARSET_VECTORS")
	if dir == "" {
		dir = "testdata"
	}
	f, err := os.Open(filepath.Join(dir, "charset_vectors.tsv.gz"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	var rows [][]string
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		rows = append(rows, strings.Split(sc.Text(), "\t"))
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return rows
}

// cvUnits parses a string written as UTF-16 units, four hex digits each.
func cvUnits(t *testing.T, s string) string {
	t.Helper()
	units := make([]uint16, 0, len(s)/4)
	for i := 0; i+4 <= len(s); i += 4 {
		v, err := strconv.ParseUint(s[i:i+4], 16, 16)
		if err != nil {
			t.Fatalf("bad units %q", s)
		}
		units = append(units, uint16(v))
	}
	return StringFromUTF16(units)
}

func cvHexUnits(s string) string { return string(cvAppendUnits(nil, s)) }

func cvAppendUnits(dst []byte, s string) []byte {
	const digits = "0123456789abcdef"
	for _, u := range StringToUTF16(s) {
		dst = append(dst, digits[u>>12], digits[u>>8&0xF], digits[u>>4&0xF], digits[u&0xF])
	}
	return dst
}

// cvReport is the REPORT-mode result as the generator writes it.
func cvReport(cs *Charset, b []byte) (res string) {
	defer func() {
		if x := recover(); x != nil {
			switch e := x.(type) {
			case *MalformedInputException:
				res = "M" + strconv.Itoa(e.GetInputLength())
			case *UnmappableCharacterException:
				res = "U" + strconv.Itoa(e.GetInputLength())
			default:
				panic(x)
			}
		}
	}()
	return cvHexUnits(cs.DecodeStrict(b))
}

func cvWrite(cs *Charset, writes ...string) []byte {
	var buf bytes.Buffer
	w := NewOutputStreamWriterWithCharset(nopWriteCloser{&buf}, cs)
	for _, s := range writes {
		w.Write(s)
	}
	w.Close()
	return buf.Bytes()
}

func cvGuard(op func() string) (res string) {
	defer func() {
		if x := recover(); x != nil {
			res = AsThrowable(x).Error()
		}
	}()
	return op()
}

func TestCharsetVectors(t *testing.T) {
	rows := readCharsetVectors(t)
	counts := map[string]int{}
	for _, r := range rows {
		counts[r[0]]++
		switch r[0] {
		case "cs":
			// A charset jutil has must carry the JDK's aliases; any other
			// must be unsupported under all its names.
			name, aliases := r[1], strings.Split(r[2], ",")
			if r[2] == "" {
				aliases = nil
			}
			if c := charsetLookup(name); c != nil && c.name == name {
				got := slices.Clone(c.aliases)
				slices.Sort(got)
				want := slices.Clone(aliases)
				slices.Sort(want)
				if !slices.Equal(got, want) {
					t.Errorf("%s aliases %v, want %v", name, got, want)
				}
				continue
			}
			for _, n := range append([]string{name}, aliases...) {
				if c := charsetLookup(n); c != nil {
					t.Errorf("%s (a name of the unsupported %s) finds %s", n, name, c.name)
				}
			}
		case "name":
			n, wantName, wantSupported := r[1], r[2], r[3]
			if !strings.HasPrefix(wantName, "java.") && charsetLookup(wantName) == nil {
				// A JDK charset that jutil does not have: it must be
				// unsupported, exactly as the JDK reports an unknown name.
				wantName = "java.nio.charset.UnsupportedCharsetException: " + n
				wantSupported = "false"
			}
			gotName := cvGuard(func() string { return Charset_ForName(n).Name() })
			gotSupported := cvGuard(func() string { return strconv.FormatBool(Charset_IsSupported(n)) })
			if gotName != wantName || gotSupported != wantSupported {
				t.Errorf("forName/isSupported(%q) = %s / %s, want %s / %s", n, gotName, gotSupported, wantName, wantSupported)
			}
		case "dec":
			cs := Charset_ForName(r[1])
			b, err := hex.DecodeString(r[2])
			if err != nil {
				t.Fatal(err)
			}
			if got := cvHexUnits(cs.Decode(b)); got != r[3] {
				t.Errorf("new String(%s, %s) = %s, want %s", r[2], r[1], got, r[3])
			}
			if got := cvHexUnits(CharStreams_ToString(NewInputStreamReader(bytes.NewReader(b), cs))); got != r[4] {
				t.Errorf("InputStreamReader(%s, %s) = %s, want %s", r[2], r[1], got, r[4])
			}
			if got := cvReport(cs, b); got != r[5] {
				t.Errorf("REPORT decode of %s as %s = %s, want %s", r[2], r[1], got, r[5])
			}
		case "enc":
			cs := Charset_ForName(r[1])
			var writes []string
			for _, w := range strings.Split(r[2], "|") {
				writes = append(writes, cvUnits(t, w))
			}
			want := ""
			if len(r) > 3 {
				want = r[3]
			}
			if got := hex.EncodeToString(cvWrite(cs, writes...)); got != want {
				t.Errorf("OutputStreamWriter(%s) writes %s: got %s, want %s", r[1], r[2], got, want)
			}
		case "can":
			cs := Charset_ForName(r[1])
			var bits strings.Builder
			for _, c := range StringToUTF16(cvUnits(t, r[2])) {
				if cs.CanEncodeChar(c) {
					bits.WriteByte('1')
				} else {
					bits.WriteByte('0')
				}
			}
			if bits.String() != r[3] {
				t.Errorf("%s canEncode of %s = %s, want %s", r[1], r[2], bits.String(), r[3])
			}
		case "digest":
			if testing.Short() {
				continue
			}
			if got := cvDigest(t, Charset_ForName(r[1]), r[2]); got != r[3] {
				t.Errorf("%s: exhaustive digest (plan %q) differs from the JDK's; CHARSET_LINES=<dir> writes the lines to <dir>/<charset>.lines, to compare with the generator's --lines output", r[1], r[2])
			}
		}
	}
	for _, k := range []string{"cs", "name", "dec", "enc", "can", "digest"} {
		if counts[k] == 0 {
			t.Errorf("no %s vectors", k)
		}
	}
}

// cvDigest computes the generator's exhaustive digest for cs.
func cvDigest(t *testing.T, cs *Charset, plan string) string {
	h := sha256.New()
	var w hashWriter = h
	if dir := os.Getenv("CHARSET_LINES"); dir != "" {
		f, err := os.Create(filepath.Join(dir, cs.Name()+".lines"))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		bw := bufio.NewWriter(f)
		defer bw.Flush()
		w = teeHash{h, bw}
	}
	var line []byte
	dec := func(b []byte) {
		line = hex.AppendEncode(line[:0], b)
		line = append(line, '\t')
		line = cvAppendUnits(line, cs.Decode(b))
		line = append(line, '\t')
		// The REPORT result without the exception (DecodeStrict wraps this
		// decode; the dec rows test it).
		if s, e := cs.impl().decode(b, true); e.length == 0 {
			line = cvAppendUnits(line, s)
		} else if e.unmappable {
			line = strconv.AppendInt(append(line, 'U'), int64(e.length), 10)
		} else {
			line = strconv.AppendInt(append(line, 'M'), int64(e.length), 10)
		}
		line = append(line, '\n')
		w.Write(line)
	}
	for _, item := range strings.Fields(plan) {
		switch item {
		case "bmp":
			for c := 0; c < 0x10000; c++ {
				line = append(cvAppendUnits(append(line[:0], 'e'), StringFromChar(uint16(c))), '\t', '0', '\t')
				if cs.CanEncodeChar(uint16(c)) {
					line[len(line)-2] = '1'
				}
				line = append(hex.AppendEncode(line, cvWrite(cs, StringFromChar(uint16(c)))), '\n')
				w.Write(line)
			}
		case "sup":
			for cp := 0x10000; cp <= 0x10FFFF; cp += 257 {
				line = append(fmt.Appendf(line[:0], "s%06x\t", cp), hex.EncodeToString(cvWrite(cs, string(rune(cp))))...)
				w.Write(append(line, '\n'))
			}
		case "gb4":
			b := make([]byte, 4)
			for b1 := 0x81; b1 <= 0xFE; b1++ {
				for b2 := 0x30; b2 <= 0x39; b2++ {
					for b3 := 0x81; b3 <= 0xFE; b3++ {
						for b4 := 0x30; b4 <= 0x39; b4++ {
							b[0], b[1], b[2], b[3] = byte(b1), byte(b2), byte(b3), byte(b4)
							dec(b)
						}
					}
				}
			}
		default:
			var pre, suf []byte
			var n int
			switch plus := strings.IndexByte(item, '+'); {
			case plus < 0:
				n, _ = strconv.Atoi(item)
			case plus == 1:
				n, _ = strconv.Atoi(item[:1])
				suf, _ = hex.DecodeString(item[2:])
			default:
				pre, _ = hex.DecodeString(item[:plus])
				n, _ = strconv.Atoi(item[plus+1:])
			}
			b := make([]byte, len(pre)+n+len(suf))
			copy(b, pre)
			copy(b[len(pre)+n:], suf)
			for v := 0; v < 1<<(8*n); v++ {
				for k := 0; k < n; k++ {
					b[len(pre)+k] = byte(v >> (8 * (n - 1 - k)))
				}
				dec(b)
			}
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

type hashWriter interface {
	Write([]byte) (int, error)
}

type teeHash struct {
	h hash.Hash
	w *bufio.Writer
}

func (t teeHash) Write(b []byte) (int, error) {
	t.w.Write(b)
	return t.h.Write(b)
}
