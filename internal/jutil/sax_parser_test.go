package jutil

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
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
	h.out.WriteString("chars " + saxVectorEscape(StringFromUTF16(ch[start:start+length])) + "\n")
}
func (h *saxVectorHandler) IgnorableWhitespace(ch []uint16, start int, length int) {
	h.out.WriteString("ws " + saxVectorEscape(StringFromUTF16(ch[start:start+length])) + "\n")
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

// TestSAXParserMatchesJava parses testdata/sax/*.xml and compares the events
// with the JDK 21 ones (testdata/sax/expected.txt, from
// tools/oracle/java/JutilSaxVectors.java).
func TestSAXParserMatchesJava(t *testing.T) {
	t.Skip("XMLReader.parse is not implemented (see sax_parser.go); these JDK vectors are the acceptance test for a clean-room parser")
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
	for _, f := range files {
		name := filepath.Base(f)
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		out := &strings.Builder{}
		Try(func() {
			factory := SAXParserFactory_NewInstance()
			factory.SetValidating(false)
			factory.SetXIncludeAware(false)
			factory.SetFeature("http://xml.org/sax/features/external-general-entities", false)
			factory.SetFeature("http://xml.org/sax/features/external-parameter-entities", false)
			factory.SetFeature("http://apache.org/xml/features/nonvalidating/load-external-dtd", false)
			factory.SetFeature(XMLConstants_FEATURE_SECURE_PROCESSING, true)
			reader := factory.NewSAXParser().GetXMLReader()
			reader.SetContentHandler(&saxVectorHandler{out: out})
			reader.Parse(NewInputSource(bytes.NewReader(data)))
		}, func(e ISAXException) {
			out.WriteString("error org.xml.sax.SAXParseException\n")
		})
		want, ok := expected[name]
		if !ok {
			t.Errorf("%s: no Java vector", name)
			continue
		}
		if name == "bigtext.xml" || name == "nested.xml" {
			// The JDK splits some runs of text in these two files into
			// adjacent characters() calls at points that depend on its
			// input buffering. Compare the concatenated text.
			join := func(s string) string { return strings.ReplaceAll(s, "\nchars ", "") }
			if join(out.String()) != join(want) {
				t.Errorf("%s: text differs from Java", name)
			}
			continue
		}
		if out.String() != want {
			t.Errorf("%s:\n got:\n%s\nJava:\n%s", name, out.String(), want)
		}
	}
}
