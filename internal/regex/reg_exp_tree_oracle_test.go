package regex

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Differential test against the Java RegExpTree (tools/oracle/java/RegExpDump.java).
//
// testdata/regexp_oracle.txt.gz is `RegExpDump dump` over hand-written edge
// cases (testdata/regexp_edge_cases.txt) and a sample of the distinct regular
// expression literals in the pinned corpus. REGEXP_ORACLE=<dump file> (plain
// or .gz) compares another dump, e.g. the one for every literal in the corpus.
// REGEXP_ORACLE_DIFF=<file> writes the Go dump of the failing records there.

func TestRegExpTreeMatchesJavaOracle(t *testing.T) {
	compareOracle(t, "testdata/regexp_oracle.txt.gz")
}

func TestRegExpTreeMatchesJavaOracleFromEnv(t *testing.T) {
	path := os.Getenv("REGEXP_ORACLE")
	if path == "" {
		t.Skip("REGEXP_ORACLE not set")
	}
	compareOracle(t, path)
}

func compareOracle(t *testing.T, path string) {
	records := readOracle(t, path)
	if len(records) == 0 {
		t.Fatalf("%s: no records", path)
	}
	var diffOut io.Writer
	if p := os.Getenv("REGEXP_ORACLE_DIFF"); p != "" {
		f, err := os.Create(p)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		diffOut = f
	}
	bad := 0
	for _, r := range records {
		var sb strings.Builder
		dumpOne(&sb, r.pattern, r.flags)
		if got := sb.String(); got != r.want {
			bad++
			if bad <= 10 {
				t.Errorf("mismatch for %s\n--- java\n%s--- go\n%s", r.header, r.want, got)
			}
			if diffOut != nil {
				fmt.Fprintf(diffOut, "### java\n%s### go\n%s", r.want, got)
			}
		}
	}
	t.Logf("%s: %d records, %d identical", path, len(records), len(records)-bad)
	if bad > 0 {
		t.Fatalf("%d/%d records differ", bad, len(records))
	}
}

type oracleRecord struct {
	header         string
	pattern, flags string
	want           string
}

func readOracle(t *testing.T, path string) []oracleRecord {
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
	sc.Buffer(make([]byte, 1<<20), 1<<30)
	var records []oracleRecord
	var cur *strings.Builder
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "=== ") {
			if cur != nil {
				records[len(records)-1].want = cur.String()
			}
			pattern, flags, ok := strings.Cut(line[len("=== "):], "\t")
			if !ok {
				t.Fatalf("bad header %q", line)
			}
			records = append(records, oracleRecord{header: line, pattern: unescape(t, pattern), flags: unescape(t, flags)})
			cur = &strings.Builder{}
		}
		if cur == nil {
			t.Fatalf("line before first record: %q", line)
		}
		cur.WriteString(line)
		cur.WriteByte('\n')
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	if cur != nil {
		records[len(records)-1].want = cur.String()
	}
	return records
}

// escape mirrors RegExpDump.escape: backslash doubled, chars outside
// printable ASCII as \uXXXX (UTF-16 code units, so lone surrogates survive).
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

func line(sb *strings.Builder, s string) {
	sb.WriteString(escape(s))
	sb.WriteByte('\n')
}

// dumpOne mirrors RegExpDump.dumpOne.
func dumpOne(sb *strings.Builder, pattern, flags string) {
	sb.WriteString("=== " + escape(pattern) + "\t" + escape(flags) + "\n")
	var t RegExpTree
	if e := catch(func() { t = RegExpTree_ParseRegExp(pattern, flags) }); e != nil {
		line(sb, "parse: "+exceptionString(e))
		return
	}
	dumpAnalyzed(sb, "parse", t, flags)
	simplified(sb, "simplify(flags)", t, flags)
	if strings.IndexByte(flags, 'i') < 0 {
		simplified(sb, "simplify(flags+i)", t, flags+"i")
	}
}

func simplified(sb *strings.Builder, label string, t RegExpTree, flags string) {
	var s RegExpTree
	if e := catch(func() { s = t.Simplify(flags) }); e != nil {
		line(sb, label+": "+exceptionString(e))
		return
	}
	dumpAnalyzed(sb, label, s, flags)
	line(sb, "  equals(parse)="+strconv.FormatBool(s.Equals(t)))
}

func dumpAnalyzed(sb *strings.Builder, label string, t RegExpTree, flags string) {
	var out strings.Builder
	e := catch(func() {
		line(&out, label+": "+t.String())
		src := jutil.NewStringBuilder()
		t.AppendSourceCode(src)
		line(&out, "  source="+src.String())
		line(&out, fmt.Sprintf("  groups=%d hasCapturingGroup=%v caseSensitive=%v anchor=%v wholeInput=%v hash=%d",
			t.NumCapturingGroups(), t.HasCapturingGroup(), t.IsCaseSensitive(), t.ContainsAnchor(),
			RegExpTree_MatchesWholeInput(t, flags), t.HashCode()))
		dumpTree(&out, t, 1)
	})
	sb.WriteString(out.String())
	if e != nil {
		line(sb, "  "+exceptionString(e))
	}
}

func dumpTree(sb *strings.Builder, t RegExpTree, level int) {
	info := jutil.NewStringBuilder()
	t.AppendDebugInfo(info)
	line(sb, strings.Repeat("  ", level)+javaSimpleName(t)+" ["+info.String()+"]")
	for _, c := range t.Children() {
		dumpTree(sb, c, level+1)
	}
}

func javaSimpleName(t RegExpTree) string {
	return strings.TrimPrefix(fmt.Sprintf("%T", t), "*regex.RegExpTree_")
}

func catch(f func()) (thrown jutil.Throwable) {
	jutil.Try(f, func(e jutil.Throwable) { thrown = e })
	return thrown
}

// exceptionString is Java's e.getClass().getName() + ": " + e.getMessage().
func exceptionString(e jutil.Throwable) string {
	s := e.Error()
	if e.GetMessage() == "" && !strings.Contains(s, ": ") {
		s += ": null"
	}
	return s
}
