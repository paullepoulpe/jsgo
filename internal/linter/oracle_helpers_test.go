package linter

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

func checksDumpClassName(e jutil.Throwable) string {
	s := e.Error()
	if m := jutil.GetMessageOrNull(e); m != nil {
		return strings.TrimSuffix(s, ": "+*m)
	}
	return s
}

func checksDumpEscape(s string) string {
	var b strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			b.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(byte(c))
		default:
			fmt.Fprintf(&b, "\\u%04x", c)
		}
	}
	return b.String()
}

func checksDumpLine(sb *strings.Builder, s string) {
	sb.WriteString(checksDumpEscape(s))
	sb.WriteByte('\n')
}

func checksDumpSections(r io.Reader) ([]string, map[string]string, error) {
	var order []string
	sections := map[string]string{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1<<20), 1<<30)
	var cur string
	var b strings.Builder
	flush := func() {
		if cur != "" {
			sections[cur] = b.String()
		}
		b.Reset()
	}
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "# file ") {
			flush()
			cur = strings.TrimPrefix(line, "# file ")
			order = append(order, cur)
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	flush()
	return order, sections, sc.Err()
}

func checksDumpOpen(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(path, ".gz") {
		return f, nil
	}
	zr, err := gzip.NewReader(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	return struct {
		io.Reader
		io.Closer
	}{zr, f}, nil
}

func checksDumpFirstDiff(want, got string) string {
	w := strings.Split(want, "\n")
	g := strings.Split(got, "\n")
	for i := 0; i < len(w) || i < len(g); i++ {
		var wl, gl string
		if i < len(w) {
			wl = w[i]
		}
		if i < len(g) {
			gl = g[i]
		}
		if wl != gl {
			return fmt.Sprintf("line %d:\n  java: %q\n  go:   %q", i+1, wl, gl)
		}
	}
	return "(no line difference)"
}

func checksDumpReadList(t *testing.T, path string) []string {
	t.Helper()
	files, err := corpus.ReadList(path)
	if err != nil {
		t.Fatal(err)
	}
	return files
}
