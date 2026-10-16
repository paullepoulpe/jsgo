package resources

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// testdata/properties_parser.txt is the Java PropertiesParser.parse output
// (tools/oracle/java/PropertiesParserDump.java) for the inputs in
// testdata/properties_parser_inputs.txt: hand-written edge cases and real
// .properties files from the pinned jar. Closure has no PropertiesParserTest.

type propertiesRecord struct {
	input string
	want  string
}

func readPropertiesOracle(t *testing.T) []propertiesRecord {
	t.Helper()
	f, err := os.Open("testdata/properties_parser.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<26)
	var records []propertiesRecord
	var cur *strings.Builder
	flush := func() {
		if cur != nil {
			records[len(records)-1].want = cur.String()
		}
	}
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "=== ") {
			flush()
			records = append(records, propertiesRecord{input: unescape(t, line[len("=== "):])})
			cur = &strings.Builder{}
			continue
		}
		cur.WriteString(line)
		cur.WriteByte('\n')
	}
	flush()
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return records
}

func dumpProperties(source string) string {
	var sb strings.Builder
	var m *jutil.LinkedHashMap[string, string]
	var thrown jutil.Throwable
	jutil.Try(func() { m = PropertiesParser_Parse(source) }, func(e jutil.Throwable) { thrown = e })
	if thrown != nil {
		sb.WriteString(escape("exception " + thrown.Error()))
		sb.WriteByte('\n')
		return sb.String()
	}
	for k, v := range m.All() {
		sb.WriteString(escape(k + "=" + v))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func TestPropertiesParserMatchesJava(t *testing.T) {
	records := readPropertiesOracle(t)
	if len(records) < 70 {
		t.Fatalf("only %d records", len(records))
	}
	for _, r := range records {
		if got := dumpProperties(r.input); got != r.want {
			t.Errorf("parse(%q):\n--- java\n%s--- go\n%s", r.input, r.want, got)
		}
	}
}

// escape and unescape mirror RegExpDump.escape/unescape.
func escape(s string) string {
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

func unescape(t *testing.T, s string) string {
	var units []uint16
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' {
			units = append(units, uint16(c))
			continue
		}
		i++
		switch s[i] {
		case '\\':
			units = append(units, '\\')
		case 'u':
			v, err := strconv.ParseUint(s[i+1:i+5], 16, 16)
			if err != nil {
				t.Fatalf("bad escape in %q", s)
			}
			units = append(units, uint16(v))
			i += 4
		default:
			t.Fatalf("bad escape in %q", s)
		}
	}
	return jutil.StringFromUTF16(units)
}
