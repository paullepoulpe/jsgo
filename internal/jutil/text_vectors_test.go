package jutil

import (
	"bufio"
	"encoding/hex"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// Vectors from tools/oracle/java/JutilTextVectors.java (random inputs, the
// seed is the generator's second argument). JUTIL_TEXT_VECTORS=<dir> reads
// vectors written with another seed instead of testdata.

func readTextVectors(t *testing.T, name string, minRows int) [][]string {
	t.Helper()
	dir := os.Getenv("JUTIL_TEXT_VECTORS")
	if dir == "" {
		dir = "testdata"
	}
	f, err := os.Open(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var rows [][]string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		rows = append(rows, strings.Split(sc.Text(), "\t"))
	}
	if len(rows) < minRows {
		t.Fatalf("%s: only %d vectors", name, len(rows))
	}
	return rows
}

// textGuard renders a result like JutilTextVectors.guard (a null exception
// message is \N).
func textGuard(op func() string) (res string) {
	defer func() {
		if r := recover(); r != nil {
			th := AsThrowable(r)
			msg := `\N`
			if th.throwable().hasMsg {
				msg = escVector(th.GetMessage())
			}
			res = "EXC:" + th.throwable().className + ":" + msg
		}
	}()
	return escVector(op())
}

func javaNull(s *string) string {
	if s == nil {
		return "null"
	}
	return *s
}

func pathOrNull(p *Path) string {
	if p == nil {
		return "null"
	}
	return p.String()
}

func describeURI(x *URI) string {
	ssp := x.ssp
	if !x.IsOpaque() {
		// The raw scheme-specific part of a hierarchical URI.
		v := x.str
		if x.scheme != nil {
			v = StringSubstringFrom(v, StringLength(*x.scheme)+1)
		}
		if x.fragment != nil {
			v = StringSubstring(v, 0, StringLength(v)-StringLength(*x.fragment)-1)
		}
		ssp = &v
	}
	return strconv.FormatBool(x.IsAbsolute()) + "|" + strconv.FormatBool(x.IsOpaque()) + "|" +
		strconv.Itoa(int(x.HashCode())) + "|" + x.String() + "|" + javaNull(x.GetScheme()) + "|" +
		javaNull(ssp) + "|" + javaNull(x.GetRawAuthority()) + "|" + javaNull(x.userInfo) + "|" +
		javaNull(x.GetHost()) + "|" + strconv.Itoa(x.GetPort()) + "|" + javaNull(x.GetRawPath()) + "|" +
		javaNull(x.GetRawQuery()) + "|" + javaNull(x.GetRawFragment())
}

func TestURITextVectors(t *testing.T) {
	for _, r := range readTextVectors(t, "text_uri.tsv", 5000) {
		a, _ := unescVector(t, r[1])
		var got, want string
		switch r[0] {
		case "p":
			want = r[2]
			got = textGuard(func() string { return describeURI(NewURI(a)) })
		case "eq":
			b, _ := unescVector(t, r[2])
			want = r[3]
			got = textGuard(func() string {
				x, y := NewURI(a), NewURI(b)
				return strconv.FormatBool(x.Equals(y)) + "|" + strconv.FormatBool(x.HashCode() == y.HashCode())
			})
		}
		if got != want {
			t.Errorf("%v:\n got %s\nwant %s", r, got, want)
		}
	}
}

func describePath(x *Path) string {
	var names []string
	for i := 0; i < x.GetNameCount(); i++ {
		names = append(names, x.GetName(i).String())
	}
	return x.String() + "|" + pathOrNull(x.GetParent()) + "|" + pathOrNull(x.GetFileName()) + "|" +
		pathOrNull(x.GetRoot()) + "|" + strconv.Itoa(x.GetNameCount()) + "|" + strings.Join(names, ",") + "|" +
		x.Normalize().String() + "|" + strconv.Itoa(int(x.HashCode())) + "|" + strconv.FormatBool(x.IsAbsolute())
}

func TestPathTextVectors(t *testing.T) {
	for _, r := range readTextVectors(t, "text_path.tsv", 5000) {
		p, _ := unescVector(t, r[1])
		var got, want string
		switch r[0] {
		case "of":
			want = r[2]
			got = textGuard(func() string { return describePath(Path_Of(p)) })
		case "of2":
			q, _ := unescVector(t, r[2])
			more, _ := unescVector(t, r[3])
			want = r[4]
			got = textGuard(func() string { return describePath(Path_Of(p, q, more)) })
		case "rel", "res", "cmp":
			q, _ := unescVector(t, r[2])
			want = r[3]
			got = textGuard(func() string {
				switch r[0] {
				case "rel":
					return Path_Of(p).Relativize(Path_Of(q)).String()
				case "res":
					return Path_Of(p).ResolveString(q).String() + "|" + Path_Of(p).ResolveSiblingString(q).String() + "|" +
						Path_Of(p).Resolve(Path_Of(q)).String()
				}
				return strconv.FormatBool(Path_Of(p).StartsWith(Path_Of(q))) + "|" +
					strconv.FormatBool(Path_Of(p).EndsWith(Path_Of(q))) + "|" +
					strconv.FormatBool(Path_Of(p).StartsWithString(q)) + "|" +
					strconv.FormatBool(Path_Of(p).EndsWithString(q)) + "|" +
					strconv.Itoa(Path_Of(p).CompareTo(Path_Of(q))) + "|" +
					strconv.FormatBool(Path_Of(p).Equals(Path_Of(q)))
			})
		case "sub":
			b, _ := strconv.Atoi(r[2])
			e, _ := strconv.Atoi(r[3])
			want = r[4] + "\t" + r[5]
			got = textGuard(func() string { return Path_Of(p).Subpath(b, e).String() }) + "\t" +
				textGuard(func() string { return Path_Of(p).GetName(b).String() })
		}
		if got != want {
			t.Errorf("%v:\n got %s\nwant %s", r, got, want)
		}
	}
}

// textArgSets are JutilTextVectors.argSets() with the Go types of the
// corresponding Java boxes.
var textArgSets = [][]any{
	{}, {"A", "B", "C"}, {7, int64(-1234567), 2.5}, {nil, "x'y", 0.0005}, {1.5, 0, math.Copysign(0, -1)},
	{math.NaN(), math.Inf(1), 1234567.891},
}

func TestMessageFormatTextVectors(t *testing.T) {
	for _, r := range readTextVectors(t, "text_messageformat.tsv", 5000) {
		pattern, _ := unescVector(t, r[1])
		var got string
		switch r[0] {
		case "mf":
			set, _ := strconv.Atoi(r[2])
			got = textGuard(func() string { return MessageFormat_Format(pattern, textArgSets[set]) })
		case "ml":
			v, _ := strconv.ParseInt(r[2], 10, 64)
			got = textGuard(func() string { return MessageFormat_Format(pattern, []any{v}) })
		case "md":
			bits, _ := strconv.ParseUint(r[2], 16, 64)
			got = textGuard(func() string { return MessageFormat_Format(pattern, []any{math.Float64frombits(bits)}) })
		case "mF":
			bits, _ := strconv.ParseUint(r[2], 16, 32)
			got = textGuard(func() string { return MessageFormat_Format(pattern, []any{math.Float32frombits(uint32(bits))}) })
		}
		if got != r[3] {
			t.Errorf("%v:\n got %s\nwant %s", r, got, r[3])
		}
	}
}

func TestPropertiesTextVectors(t *testing.T) {
	for _, r := range readTextVectors(t, "text_properties.tsv", 2000) {
		text, _ := unescVector(t, r[1])
		got := textGuard(func() string {
			p := PropertiesLoad(text)
			keys := p.KeySlice()
			sortJavaStrings(keys)
			var sb strings.Builder
			for i, k := range keys {
				if i > 0 {
					sb.WriteString("\x01")
				}
				sb.WriteString(k + "\x02" + p.Get(k))
			}
			return sb.String()
		})
		if got != r[2] {
			t.Errorf("%v:\n got %s\nwant %s", r, got, r[2])
		}
	}
}

// textSplitRegexes is JutilTextVectors.SPLIT_REGEX.
var textSplitRegexes = []string{"\r?\n", ",", ",+", ";|,", `\s+`, "[,;]", "(,)", "\u00e9", "\U0001F600|,", "a|b",
	"a*", "", "b?", ",|", "x*", ",*"}

func TestSplitTextVectors(t *testing.T) {
	for _, r := range readTextVectors(t, "text_split.tsv", 2000) {
		s, _ := unescVector(t, r[1])
		ri, _ := strconv.Atoi(r[2])
		limit, _ := strconv.Atoi(r[3])
		got := textGuard(func() string {
			re := regexp.MustCompile(textSplitRegexes[ri])
			var parts []string
			if limit == 0 {
				parts = StringSplitRegexp(s, re)
			} else {
				parts = stringSplitRegexpLimit(s, re, limit)
			}
			return strconv.Itoa(len(parts)) + "\x01" + strings.Join(parts, "\x01")
		})
		want := r[4]
		if !strings.Contains(want, `\u0001`) {
			want += `\u0001` // no elements
		}
		if got != want {
			t.Errorf("%v:\n got %s\nwant %s", r, got, want)
		}
	}
}

func TestCaseTextVectors(t *testing.T) {
	for _, r := range readTextVectors(t, "text_case.tsv", 3500) {
		switch r[0] {
		case "r":
			a, _ := unescVector(t, r[1])
			b, _ := unescVector(t, r[2])
			to, _ := strconv.Atoi(r[4])
			oo, _ := strconv.Atoi(r[5])
			n, _ := strconv.Atoi(r[6])
			got := UTF16RegionMatches(r[3] == "true", StringToUTF16(a), to, StringToUTF16(b), oo, n)
			if strconv.FormatBool(got) != r[7] {
				t.Errorf("%v: got %v", r, got)
			}
		case "c":
			cp, _ := strconv.ParseUint(r[1], 16, 32)
			u, _ := strconv.ParseUint(r[2], 16, 32)
			l, _ := strconv.ParseUint(r[3], 16, 32)
			if unicode15ToUpper(rune(cp)) != rune(u) || unicode15ToLower(rune(cp)) != rune(l) {
				t.Errorf("case mapping of U+%X: Go %X %X, Java %X %X", cp, unicode15ToUpper(rune(cp)), unicode15ToLower(rune(cp)), u, l)
			}
		case "q":
			s, _ := unescVector(t, r[1])
			q := Pattern_Quote(s)
			if escVector(q) != r[2] || strconv.FormatBool(Pattern_Compile(q).Matches(s)) != r[3] {
				t.Errorf("%v: got %s", r, escVector(q))
			}
		}
	}
}

func TestUTF16TextVectors(t *testing.T) {
	for _, r := range readTextVectors(t, "text_utf16.tsv", 2000) {
		cs := Charset_ForName(r[0])
		b, err := hex.DecodeString(r[1])
		if err != nil {
			t.Fatal(err)
		}
		if got := escVector(cs.Decode(b)); got != r[2] {
			t.Errorf("new String(%s, %s) = %q, want %q", r[1], r[0], got, r[2])
		}
		s, bad := cs.impl().decode(b, true)
		strict := escVector(s)
		if bad.length > 0 {
			strict = "MALFORMED:" + strconv.Itoa(bad.length)
		}
		if strict != r[3] {
			t.Errorf("strict decode of %s as %s = %q, want %q", r[1], r[0], strict, r[3])
		}
	}
}

// logWriter records what reaches the writer under a BufferedWriter.
type logWriter struct{ log strings.Builder }

func (w *logWriter) Write(s string) { w.log.WriteString(" w" + strconv.Itoa(StringLength(s))) }
func (w *logWriter) Flush()         { w.log.WriteString(" f") }
func (w *logWriter) Close()         { w.log.WriteString(" c") }

func TestBufferedWriterTextVectors(t *testing.T) {
	for _, r := range readTextVectors(t, "text_bufferedwriter.tsv", 400) {
		out := &logWriter{}
		w := NewBufferedWriter(out)
		func() {
			defer func() {
				if x := recover(); x != nil {
					out.log.WriteString(" !" + AsThrowable(x).GetMessage())
				}
			}()
			for _, op := range strings.Fields(r[1]) {
				switch op[0] {
				case 'w':
					n, _ := strconv.Atoi(op[1:])
					w.Write(strings.Repeat("x", n))
				case 'a':
					w.AppendChar('y')
				case 'f':
					w.Flush()
				case 'c':
					w.Close()
				}
			}
			out.log.WriteString(" .")
		}()
		if got := strings.TrimSpace(out.log.String()); got != r[2] {
			t.Errorf("%s: got %s, want %s", r[1], got, r[2])
		}
	}
}
