package jutil

import (
	"bufio"
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
