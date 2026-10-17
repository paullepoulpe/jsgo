package jutil

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"testing"
)

type nopWriteCloser struct{ *bytes.Buffer }

func (nopWriteCloser) Close() error { return nil }

// TestCliVectors checks OutputStreamWriter(out, charset), Charset.isSupported
// and Level.parse against testdata/cli_charset_vectors.tsv, written by
// tools/oracle/java/JutilCliVectors.java.
func TestCliVectors(t *testing.T) {
	f, err := os.Open("testdata/cli_charset_vectors.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	// The inputs, in the generator's order (written twice per stream).
	inputs := []string{"aéb", "x\U0001d4b3y", StringFromUTF16([]uint16{'p', 0xD835, 'q'}), StringFromUTF16([]uint16{0xDCB3, 'z'}), "€", ""}
	n := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), "\t")
		switch parts[0] {
		case "isSupported":
			var got string
			Try(func() { got = strconv.FormatBool(Charset_IsSupported(parts[1])) }, func(e Throwable) { got = e.Error() })
			want := parts[2]
			if got != want && !(strings.HasPrefix(want, "java.nio.charset.IllegalCharsetNameException") && strings.HasPrefix(got, "java.nio.charset.IllegalCharsetNameException")) {
				t.Errorf("isSupported(%q) = %s, want %s", parts[1], got, want)
			}
		case "level":
			var got string
			Try(func() { l := Level_Parse(parts[1]); got = l.GetName() + "/" + strconv.Itoa(int(l.IntValue())) }, func(e Throwable) { got = e.Error() })
			if got != parts[2] {
				t.Errorf("Level.parse(%q) = %s, want %s", parts[1], got, parts[2])
			}
		default:
			s := inputs[n%len(inputs)]
			n++
			var buf bytes.Buffer
			w := NewBufferedWriter(NewOutputStreamWriterWithCharset(nopWriteCloser{&buf}, Charset_ForName(parts[0])))
			w.Write(s)
			w.Write(s)
			w.Close()
			want := ""
			if len(parts) > 2 {
				want = parts[2]
			}
			if got := hex.EncodeToString(buf.Bytes()); got != want {
				t.Errorf("%s %q: got %s, want %s", parts[0], s, got, want)
			}
		}
	}
	if n != 36 {
		t.Errorf("read %d encoding vectors, want 36", n)
	}
}
