package xtbxml

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type saxVectorHandler struct {
	out *strings.Builder
}

func (h *saxVectorHandler) SetDocumentLocator(locator Locator) {}
func (h *saxVectorHandler) StartDocument()                     { h.out.WriteString("startDocument\n") }
func (h *saxVectorHandler) EndDocument()                       { h.out.WriteString("endDocument\n") }
func (h *saxVectorHandler) StartPrefixMapping(prefix string, uri string) {
	h.out.WriteString("startPrefixMapping " + saxVectorEscape(prefix) + "\n")
}
func (h *saxVectorHandler) EndPrefixMapping(prefix string) {
	h.out.WriteString("endPrefixMapping " + saxVectorEscape(prefix) + "\n")
}
func (h *saxVectorHandler) StartElement(uri string, localName string, qName string, atts Attributes) {
	h.out.WriteString("start " + saxVectorEscape(uri) + "|" + saxVectorEscape(localName) + "|" + saxVectorEscape(qName))
	for i := 0; i < atts.GetLength(); i++ {
		h.out.WriteString(" " + saxVectorEscape(atts.GetQName(i)) + "=" + saxVectorEscape(atts.GetValueAt(i)))
	}
	h.out.WriteString("\n")
}
func (h *saxVectorHandler) EndElement(uri string, localName string, qName string) {
	h.out.WriteString("end " + saxVectorEscape(uri) + "|" + saxVectorEscape(localName) + "|" + saxVectorEscape(qName) + "\n")
}
func (h *saxVectorHandler) Characters(ch []uint16, start int, length int) {
	h.out.WriteString("chars " + saxVectorEscape(jutil.StringFromUTF16(ch[start:start+length])) + "\n")
}
func (h *saxVectorHandler) IgnorableWhitespace(ch []uint16, start int, length int) {
	h.out.WriteString("ws " + saxVectorEscape(jutil.StringFromUTF16(ch[start:start+length])) + "\n")
}
func (h *saxVectorHandler) ProcessingInstruction(target string, data string) {
	h.out.WriteString("pi " + saxVectorEscape(target) + "|" + saxVectorEscape(data) + "\n")
}
func (h *saxVectorHandler) SkippedEntity(name string) {
	h.out.WriteString("skipped " + saxVectorEscape(name) + "\n")
}

func saxVectorEscape(s string) string {
	return matcherVectorEscape(&s)
}

func matcherVectorEscape(s *string) string {
	if s == nil {
		return "<null>"
	}
	var sb strings.Builder
	for _, c := range jutil.StringToUTF16(*s) {
		switch {
		case c == '\\':
			sb.WriteString(`\\`)
		case c == '\t':
			sb.WriteString(`\t`)
		case c == '\n':
			sb.WriteString(`\n`)
		case c == '\r':
			sb.WriteString(`\r`)
		case c < 0x20 || c > 0x7e:
			fmt.Fprintf(&sb, `\u%04x`, c)
		default:
			sb.WriteByte(byte(c))
		}
	}
	return sb.String()
}

// TestSAXParserMatchesJava parses testdata/sax/*.xml and compares the events
// with the JDK 21 ones (testdata/sax/expected.txt, from
// tools/oracle/java/JutilSaxVectors.java).
func TestSAXParserMatchesJava(t *testing.T) {
	expectedBytes, err := os.ReadFile("testdata/sax/expected.txt")
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{}
	for _, block := range strings.Split(string(expectedBytes), "== ")[1:] {
		name, body, _ := strings.Cut(block, "\n")
		expected[name] = body
	}
	files, _ := filepath.Glob("testdata/sax/*.xml")
	sort.Strings(files)
	if len(files) == 0 {
		t.Fatal("no documents")
	}
	saved := jutil.System_Err()
	jutil.System_SetErr(jutil.NewPrintStream(io.Discard))
	defer jutil.System_SetErr(saved)
	for _, f := range files {
		name := filepath.Base(f)
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		out := &strings.Builder{}
		if err := Parse(bytes.NewReader(data), &saxVectorHandler{out: out}); err != nil {
			if _, ok := err.(ISAXException); !ok {
				t.Fatalf("%s: %v", name, err)
			}
			out.WriteString("error org.xml.sax.SAXParseException\n")
		}
		want, ok := expected[name]
		if !ok {
			t.Errorf("%s: no Java vector", name)
			continue
		}
		if out.String() != want {
			t.Errorf("%s:\n got:\n%s\nJava:\n%s", name, out.String(), want)
		}
	}
}

// saxPositionHandler records the events as JutilSaxVectors.runFull prints
// them: every line ends with the Locator position at the time of the call.
type saxPositionHandler struct {
	out     *strings.Builder
	locator Locator
}

func (h *saxPositionHandler) pos() string {
	if h.locator == nil {
		return " @null"
	}
	ids := ""
	if h.locator.GetPublicId() != nil || h.locator.GetSystemId() != nil {
		ids = " [" + matcherVectorEscape(h.locator.GetPublicId()) + "|" + matcherVectorEscape(h.locator.GetSystemId()) + "]"
	}
	return " @" + strconv.Itoa(h.locator.GetLineNumber()) + ":" + strconv.Itoa(h.locator.GetColumnNumber()) + ids
}

func (h *saxPositionHandler) SetDocumentLocator(locator Locator) {
	h.locator = locator
	h.out.WriteString("locator" + h.pos() + "\n")
}
func (h *saxPositionHandler) StartDocument() { h.out.WriteString("startDocument" + h.pos() + "\n") }
func (h *saxPositionHandler) EndDocument()   { h.out.WriteString("endDocument" + h.pos() + "\n") }
func (h *saxPositionHandler) StartPrefixMapping(prefix string, uri string) {
	h.out.WriteString("startPrefixMapping " + saxVectorEscape(prefix) + "|" + saxVectorEscape(uri) + h.pos() + "\n")
}
func (h *saxPositionHandler) EndPrefixMapping(prefix string) {
	h.out.WriteString("endPrefixMapping " + saxVectorEscape(prefix) + h.pos() + "\n")
}
func (h *saxPositionHandler) StartElement(uri string, localName string, qName string, atts Attributes) {
	h.out.WriteString("start " + saxVectorEscape(uri) + "|" + saxVectorEscape(localName) + "|" + saxVectorEscape(qName))
	for i := 0; i < atts.GetLength(); i++ {
		h.out.WriteString(" " + saxVectorEscape(atts.GetQName(i)) + "=" + saxVectorEscape(atts.GetValueAt(i)))
		if v := atts.GetValue(atts.GetQName(i)); v == nil || *v != atts.GetValueAt(i) {
			h.out.WriteString("<GetValue mismatch>")
		}
	}
	if atts.GetValue("no such attribute") != nil {
		h.out.WriteString("<GetValue of a missing attribute>")
	}
	h.out.WriteString(h.pos() + "\n")
}
func (h *saxPositionHandler) EndElement(uri string, localName string, qName string) {
	h.out.WriteString("end " + saxVectorEscape(uri) + "|" + saxVectorEscape(localName) + "|" + saxVectorEscape(qName) + h.pos() + "\n")
}
func (h *saxPositionHandler) Characters(ch []uint16, start int, length int) {
	h.out.WriteString("chars " + saxVectorEscape(jutil.StringFromUTF16(ch[start:start+length])) + h.pos() + "\n")
}
func (h *saxPositionHandler) IgnorableWhitespace(ch []uint16, start int, length int) {
	h.out.WriteString("ws " + saxVectorEscape(jutil.StringFromUTF16(ch[start:start+length])) + h.pos() + "\n")
}
func (h *saxPositionHandler) ProcessingInstruction(target string, data string) {
	h.out.WriteString("pi " + saxVectorEscape(target) + "|" + saxVectorEscape(data) + h.pos() + "\n")
}
func (h *saxPositionHandler) SkippedEntity(name string) {
	h.out.WriteString("skipped " + saxVectorEscape(name) + h.pos() + "\n")
}

// saxRunFull is the Go side of JutilSaxVectors.runFull.
func saxRunFull(doc []byte) string {
	out := &strings.Builder{}
	errBuf := &bytes.Buffer{}
	saved := jutil.System_Err()
	jutil.System_SetErr(jutil.NewPrintStream(errBuf))
	defer jutil.System_SetErr(saved)
	h := &saxPositionHandler{out: out}
	var thrown jutil.Throwable
	jutil.Try(func() {
		thrown = Parse(bytes.NewReader(doc), h)
	}, func(e jutil.Throwable) {
		thrown = e
	})
	if thrown != nil {
		switch e := thrown.(type) {
		case *SAXParseException:
			fmt.Fprintf(out, "error org.xml.sax.SAXParseException %d:%d <null>|<null> %s\n",
				e.GetLineNumber(), e.GetColumnNumber(), saxVectorEscape(e.GetMessage()))
		default:
			name, _, _ := strings.Cut(e.Error(), ": ")
			out.WriteString("error " + name + " " + saxVectorEscape(e.GetMessage()) + "\n")
		}
	}
	if errBuf.Len() > 0 {
		out.WriteString("stderr " + saxVectorEscape(errBuf.String()) + "\n")
	}
	return out.String()
}

// saxDecodeDocBytes undoes JutilSaxVectors.escBytes.
func saxDecodeDocBytes(s string) []byte {
	var b []byte
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			if s[i+1] == 'x' {
				v, _ := strconv.ParseUint(s[i+2:i+4], 16, 8)
				b = append(b, byte(v))
				i += 3
				continue
			}
			b = append(b, s[i+1])
			i++
			continue
		}
		b = append(b, s[i])
	}
	return b
}

type saxVectorCase struct {
	name string
	doc  []byte
	want string
}

func readSaxVectors(t *testing.T, path string) []saxVectorCase {
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var r io.Reader = f
	if strings.HasSuffix(path, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			t.Fatal(err)
		}
		r = gz
	}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1<<20), 1<<26)
	var cases []saxVectorCase
	var cur *saxVectorCase
	var want strings.Builder
	finish := func() {
		if cur != nil {
			cur.want = want.String()
			cases = append(cases, *cur)
		}
		want.Reset()
	}
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "# "):
		case strings.HasPrefix(line, "== "):
			finish()
			cur = &saxVectorCase{name: line[3:]}
		case strings.HasPrefix(line, "doc "):
			cur.doc = saxDecodeDocBytes(line[4:])
		default:
			want.WriteString(line + "\n")
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	finish()
	return cases
}

// TestSAXParserVectors replays testdata/sax_vectors.txt.gz (from
// tools/oracle/java/JutilSaxVectors.java): events with Locator positions,
// exceptions with their messages and positions, and System.err, for
// hand-written and seeded random documents. SAX_VECTORS=<file> (plain or
// .gz) replays another vector file instead; SAX_ONLY=<name> runs one
// document and logs both event streams.
func TestSAXParserVectors(t *testing.T) {
	path := os.Getenv("SAX_VECTORS")
	if path == "" {
		path = "testdata/sax_vectors.txt.gz"
	}
	cases := readSaxVectors(t, path)
	if len(cases) == 0 {
		t.Fatal("no vectors")
	}
	failed := 0
	only := os.Getenv("SAX_ONLY")
	for _, c := range cases {
		if only != "" && c.name != only {
			continue
		}
		got := saxRunFull(c.doc)
		c.want = saxFatalErrorLinesOnly(c.want)
		if only != "" {
			t.Logf("doc %q\ngot:\n%s\nJava:\n%s", c.doc, got, c.want)
		}
		if got == c.want {
			continue
		}
		failed++
		if failed <= saxMaxReported() {
			gl, wl := strings.Split(got, "\n"), strings.Split(c.want, "\n")
			i := 0
			for i < len(gl) && i < len(wl) && gl[i] == wl[i] {
				i++
			}
			ctx := func(ls []string) string {
				lo := max(0, i-3)
				hi := min(len(ls), i+3)
				var b strings.Builder
				for _, l := range ls[lo:hi] {
					if len(l) > 300 {
						l = l[:300] + "..."
					}
					b.WriteString("    " + l + "\n")
				}
				return b.String()
			}
			doc := string(c.doc)
			if len(doc) > 400 {
				doc = doc[:400] + "..."
			}
			t.Errorf("%s: first difference at event line %d\n doc: %q\n got:\n%s Java:\n%s", c.name, i+1, doc, ctx(gl), ctx(wl))
		}
	}
	if failed > 0 {
		t.Errorf("%d of %d documents differ", failed, len(cases))
	}
}

// saxMaxReported is the number of differing documents to show
// (SAX_MAX_REPORTED, default 5).
func saxMaxReported() int {
	if n, err := strconv.Atoi(os.Getenv("SAX_MAX_REPORTED")); err == nil {
		return n
	}
	return 5
}

// saxFatalErrorLinesOnly keeps, of the System.err output a vector records,
// what the parser prints itself (from "[Fatal Error]" on). For some
// malformed internal subsets the JDK also prints traces of its internal
// exceptions there first; those are not reproduced.
func saxFatalErrorLinesOnly(want string) string {
	lines := strings.Split(want, "\n")
	for i, l := range lines {
		rest, ok := strings.CutPrefix(l, "stderr ")
		if !ok {
			continue
		}
		if k := strings.Index(rest, "[Fatal Error]"); k >= 0 {
			lines[i] = "stderr " + rest[k:]
		}
	}
	return strings.Join(lines, "\n")
}
