package jutil

import (
	"bufio"
	"encoding/hex"
	"os"
	"path/filepath"
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

func TestURITextVectors(t *testing.T) {
	for _, r := range readTextVectors(t, "text_uri.tsv", 5000) {
		if r[0] != "p" {
			continue
		}
		a, _ := unescVector(t, r[1])
		want := r[2]
		u, err := ParseURI(a)
		if err != nil {
			msg := strings.TrimPrefix(err.Error(), "java.net.URISyntaxException: ")
			got := "EXC:java.net.URISyntaxException:" + escVector(msg)
			if got != want {
				t.Errorf("%v:\n got %s\nwant %s", r, got, want)
			}
			continue
		}
		cols := strings.Split(want, "|")
		if got := strconv.FormatBool(u.IsAbsolute()); got != cols[0] {
			t.Errorf("%v: IsAbsolute = %s, want %s", r, got, cols[0])
		}
		if got := strconv.FormatBool(u.path == nil); got != cols[1] {
			t.Errorf("%v: isOpaque = %s, want %s", r, got, cols[1])
		}
		if got := escVector(u.String()); got != cols[3] {
			t.Errorf("%v: String = %s, want %s", r, got, cols[3])
		}
		if got := escVector(javaNull(u.path)); got != cols[10] {
			t.Errorf("%v: path = %s, want %s", r, got, cols[10])
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
					strconv.FormatBool(Path_Of(p).startsWithString(q)) + "|" +
					strconv.FormatBool(Path_Of(p).endsWithString(q)) + "|" +
					strconv.Itoa(Path_Of(p).CompareTo(Path_Of(q))) + "|" +
					strconv.FormatBool(Path_Of(p).Equals(Path_Of(q)))
			})
		case "sub":
			b, _ := strconv.Atoi(r[2])
			e, _ := strconv.Atoi(r[3])
			want = r[4] + "\t" + r[5]
			got = textGuard(func() string { return Path_Of(p).subpath(b, e).String() }) + "\t" +
				textGuard(func() string { return Path_Of(p).GetName(b).String() })
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
