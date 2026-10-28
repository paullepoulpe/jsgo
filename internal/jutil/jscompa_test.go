package jutil

import (
	"encoding/hex"
	"math"
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilJscompAVectors.java against
// the pinned closure-compiler.jar (JDK 21, shaded Guava).

// unescVector undoes JutilJscompAVectors.esc; "\N" is Java null.
func unescVector(t *testing.T, s string) (string, bool) {
	t.Helper()
	if s == `\N` {
		return "", false
	}
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
		case 't':
			units = append(units, '\t')
		case 'n':
			units = append(units, '\n')
		case 'r':
			units = append(units, '\r')
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
	return StringFromUTF16(units), true
}

func escVector(s string) string {
	var sb strings.Builder
	for _, c := range StringToUTF16(s) {
		switch c {
		case '\\':
			sb.WriteString(`\\`)
		case '\t':
			sb.WriteString(`\t`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		default:
			if c < 0x20 || (c >= 0xD800 && c <= 0xDFFF) || c > 0x7e {
				sb.WriteString(`\u`)
				h := strconv.FormatUint(uint64(c), 16)
				sb.WriteString(strings.Repeat("0", 4-len(h)) + h)
			} else {
				sb.WriteByte(byte(c))
			}
		}
	}
	return sb.String()
}

// guardVector renders a result like JutilJscompAVectors.guard.
func guardVector(op func() string) (res string) {
	defer func() {
		if r := recover(); r != nil {
			th := AsThrowable(r)
			res = "EXC:" + th.throwable().className + ":" + escVector(th.GetMessage())
		}
	}()
	return escVector(op())
}

func javaNullable(p *string) string {
	if p == nil {
		return "null"
	}
	return *p
}

func TestMessageFormatMatchesJava(t *testing.T) {
	rows := readTSV(t, "jscompa_messageformat.tsv")
	if len(rows) < 1500 {
		t.Fatalf("only %d vectors", len(rows))
	}
	argSets := map[string][]any{
		"0": {},
		"4": {"A", "B", "C", "D"},
		"3": {"x'y", "{0}", "ü"},
	}
	for _, r := range rows {
		pattern, _ := unescVector(t, r[0])
		args := argSets[r[1]]
		got := guardVector(func() string { return MessageFormat_Format(pattern, args) })
		if got != r[2] {
			t.Errorf("MessageFormat.format(%q, %v) = %q, want %q", pattern, args, got, r[2])
		}
	}
}

func TestPathMatchesJava(t *testing.T) {
	for _, r := range readTSV(t, "jscompa_path.tsv") {
		p, _ := unescVector(t, r[1])
		var got, want string
		switch r[0] {
		case "of":
			want = r[2]
			got = guardVector(func() string {
				x := Path_Of(p)
				parent, name := "null", "null"
				if x.GetParent() != nil {
					parent = x.GetParent().String()
				}
				if x.GetFileName() != nil {
					name = x.GetFileName().String()
				}
				return x.String() + "|" + parent + "|" + name + "|" + strconv.Itoa(x.GetNameCount()) + "|" +
					x.Normalize().String() + "|" + strconv.Itoa(int(x.HashCode())) + "|" + strconv.FormatBool(x.IsAbsolute())
			})
		case "rel", "sib", "sw":
			q, _ := unescVector(t, r[2])
			want = r[3]
			got = guardVector(func() string {
				switch r[0] {
				case "rel":
					return Path_Of(p).Relativize(Path_Of(q)).String()
				case "sib":
					return Path_Of(p).ResolveSiblingString(q).Normalize().String()
				}
				sign := Path_Of(p).CompareTo(Path_Of(q))
				if sign > 0 {
					sign = 1
				} else if sign < 0 {
					sign = -1
				}
				return strconv.FormatBool(Path_Of(p).StartsWith(Path_Of(q))) + "," +
					strconv.FormatBool(Path_Of(p).EndsWith(Path_Of(q))) + "," + strconv.Itoa(sign)
			})
		}
		if got != want {
			t.Errorf("%v: got %q", r, got)
		}
	}
}

func TestURIMatchesJava(t *testing.T) {
	for _, r := range readTSV(t, "jscompa_uri.tsv") {
		a, _ := unescVector(t, r[1])
		var got, want string
		if r[0] == "p" {
			want = r[2]
			got = guardVector(func() string {
				x := NewURI(a)
				return strconv.FormatBool(x.IsAbsolute()) + "|" + strconv.FormatBool(x.IsOpaque()) + "|" +
					strconv.Itoa(int(x.HashCode())) + "|" + x.String() + "|" + javaNullable(x.GetScheme()) + "|" +
					javaNullable(x.GetHost()) + "|" + strconv.Itoa(x.GetPort()) + "|" + javaNullable(x.GetRawPath())
			})
		} else {
			b, _ := unescVector(t, r[2])
			want = r[3]
			got = guardVector(func() string { return strconv.FormatBool(NewURI(a).Equals(NewURI(b))) })
		}
		if got != want {
			t.Errorf("%v: got %q", r, got)
		}
	}
}

func TestFormatDecimalFloatMatchesJava(t *testing.T) {
	rows := readTSV(t, "jscompa_fixed.tsv")
	if len(rows) < 1000 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		bits, err := strconv.ParseUint(r[0], 16, 64)
		if err != nil {
			t.Fatal(err)
		}
		d := math.Float64frombits(bits)
		for i, prec := range []int{1, 0, 3} {
			if got := Formatter_FormatDecimalFloat(d, prec); got != r[1+i] {
				t.Errorf("%%.%df of %v (%s) = %q, want %q", prec, d, r[0], got, r[1+i])
			}
		}
	}
}

func TestStringsMatchJava(t *testing.T) {
	for _, r := range readTSV(t, "jscompa_strings.tsv") {
		var got, want string
		switch r[0] {
		case "split/", "splitnm":
			s, _ := unescVector(t, r[1])
			sep := "/"
			if r[0] == "splitnm" {
				sep = "/node_modules/"
			}
			parts := StringSplitLiteral(s, sep)
			got = escVector(strings.Join(parts, "|")) + "\t" + strconv.Itoa(len(parts))
			want = r[2] + "\t" + r[3]
		case "b64d":
			s, _ := unescVector(t, r[1])
			want = r[2]
			got = guardVector(func() string { return hex.EncodeToString(BaseEncoding_Base64().Decode(s)) })
		case "pe":
			s, _ := unescVector(t, r[1])
			want = r[2]
			got = guardVector(func() string { return NewPercentEscaper("-_.*", false).Escape(s) })
		case "ms":
			continue
		}
		if got != want {
			t.Errorf("%v: got %q", r, got)
		}
	}
	if got := ImmutableMultiset_Of("a", "b", "a", "c", "a", "b").String(); got != "[a x 3, b x 2, c]" {
		t.Errorf("multiset: %q", got)
	}
	if got := ImmutableMultiset_Of[string]().String(); got != "[]" {
		t.Errorf("multiset: %q", got)
	}
}

func TestTreeSetMatchesJava(t *testing.T) {
	for _, r := range readTSV(t, "jscompa_treeset.tsv") {
		var prio []int
		for _, p := range strings.Split(strings.TrimSuffix(r[0], ","), ",") {
			v, _ := strconv.Atoi(p)
			prio = append(prio, v)
		}
		order := map[int]int{}
		set := NewTreeSet(func(a, b int) int {
			if d := prio[a] - prio[b]; d != 0 {
				return d
			}
			return order[b] - order[a]
		})
		adds := 0
		for _, op := range strings.Fields(r[1]) {
			g, _ := strconv.Atoi(op[1:])
			if op[0] == 'r' {
				set.Remove(g)
			} else {
				adds++
				order[g] = adds
				set.Remove(g)
				set.Add(g)
			}
		}
		var asc, desc strings.Builder
		for x := range set.All() {
			asc.WriteString(strconv.Itoa(x) + ",")
		}
		for x := range set.Descending() {
			desc.WriteString(strconv.Itoa(x) + ",")
		}
		if asc.String() != r[2] || desc.String() != r[3] || strconv.Itoa(set.Size()) != r[4] {
			t.Errorf("%v: got %s %s %d", r, asc.String(), desc.String(), set.Size())
		}
	}
}

func TestPatternTranslation(t *testing.T) {
	cases := []struct {
		re, s string
		want  bool
	}{
		{`Unnecessary escape:.*`, "Unnecessary escape: x", true},
		{`Unnecessary escape:.*`, "Unnecessary escape: x\ny", false},
		{`Unnecessary escape:.*`, "Unnecessary escape: x y", false},
		{`.*Unknown type.*\n.*`, "a Unknown type b\nc", true},
		{`^Too many template parameters\n.*`, "Too many template parameters\nfoo", true},
		{`Missing type declaration\.`, "Missing type declaration.", true},
		{`Missing type declaration\.`, "Missing type declarationX", false},
		{Pattern_Quote(`a.b\E{0}`), `a.b\E{0}`, true},
		{`\Qx.y\E.*\Qz\E`, "x.yANYz", true},
		{`\Qx.y\E.*\Qz\E`, "xzy", false},
		{`[a.]+`, "a.a", true},
		{`\s+`, " \t\x0b", true},
	}
	for _, c := range cases {
		if got := Pattern_Compile(c.re).Matches(c.s); got != c.want {
			t.Errorf("%q matches %q = %v, want %v", c.re, c.s, got, c.want)
		}
	}
}

func TestCharsetDecodeMatchesJava(t *testing.T) {
	rows := readTSV(t, "jscompa_decode.tsv")
	if len(rows) < 1000 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		cs := Charset_ForName(r[0])
		b, err := hex.DecodeString(r[1])
		if err != nil {
			t.Fatal(err)
		}
		if got := escVector(cs.Decode(b)); got != r[2] {
			t.Errorf("new String(%s, %s) = %q, want %q", r[1], r[0], got, r[2])
		}
		var strict string
		func() {
			defer func() {
				if x := recover(); x != nil {
					if m, ok := x.(*MalformedInputException); ok {
						strict = "MALFORMED:" + strconv.Itoa(m.GetInputLength())
						return
					}
					panic(x)
				}
			}()
			strict = escVector(cs.DecodeStrict(b))
		}()
		if strict != r[3] {
			t.Errorf("strict decode of %s as %s = %q, want %q", r[1], r[0], strict, r[3])
		}
	}
}

func TestCharsetDecodeUTF16ReversedMark(t *testing.T) {
	// Checked against JDK 21: U+FFFE is decoded, not reported as malformed.
	for in, want := range map[string]string{"0041fffe": "A\ufffe", "fefffffe0041": "\ufffeA", "fffefeff": "\ufffe"} {
		b, _ := hex.DecodeString(in)
		if got := StandardCharsets_UTF_16.Decode(b); got != want {
			t.Errorf("%s: got %q, want %q", in, got, want)
		}
	}
}

func TestJsonWriterValueLong(t *testing.T) {
	var sb strings.Builder
	w := NewJsonWriter(&sb)
	w.BeginArray().ValueLong(-5).ValueLong(9007199254740993).EndArray()
	w.Flush()
	w.Close()
	if sb.String() != "[-5,9007199254740993]" {
		t.Errorf("got %s", sb.String())
	}
}
