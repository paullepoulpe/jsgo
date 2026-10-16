package jutil

import (
	"bufio"
	"compress/gzip"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/GsonVectors.java against the
// Gson 2.9.1 shaded into the pinned closure-compiler.jar.

// gsonUnescape decodes a GsonVectors field ("\\0" is Java null).
func gsonUnescape(t *testing.T, s string) (string, bool) {
	if s == "\\0" {
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

// gsonEscape is GsonVectors.esc.
func gsonEscape(s string) string {
	var sb strings.Builder
	for _, c := range StringToUTF16(s) {
		switch {
		case c == '\\':
			sb.WriteString("\\\\")
		case c == '\t':
			sb.WriteString("\\t")
		case c == '\n':
			sb.WriteString("\\n")
		case c >= 0x20 && c <= 0x7e:
			sb.WriteByte(byte(c))
		default:
			sb.WriteString(StringFormat("\\u%04x", int32(c)))
		}
	}
	return sb.String()
}

// gsonRun is GsonVectors.run.
func gsonRun(f func() any) (s string) {
	Try(func() {
		s = "=" + ToJavaString(f())
	}, func(t Throwable) {
		s = "!" + t.Error()
	})
	return s
}

func gsonAccessors(e JsonElement) string {
	var sb strings.Builder
	sb.WriteString(gsonRun(func() any { return e.GetAsInt() }))
	sb.WriteString(" | ")
	sb.WriteString(gsonRun(func() any { return e.GetAsLong() }))
	sb.WriteString(" | ")
	sb.WriteString(gsonRun(func() any { return e.GetAsString() }))
	sb.WriteString(" | ")
	sb.WriteString(gsonRun(func() any { return e.GetAsBoolean() }))
	sb.WriteString(" | ")
	sb.WriteString(gsonRun(func() any { return e.GetAsJsonObject() }))
	sb.WriteString(" | ")
	sb.WriteString(gsonRun(func() any { return e.GetAsJsonArray() }))
	sb.WriteString(" | ")
	sb.WriteString(gsonRun(func() any { return e.IsJsonNull() }))
	sb.WriteString(ToJavaString(e.IsJsonPrimitive()))
	sb.WriteString(ToJavaString(e.IsJsonArray()))
	sb.WriteString(ToJavaString(e.IsJsonObject()))
	return sb.String()
}

func TestGsonMatchesJava(t *testing.T) {
	f, err := os.Open("testdata/gson.tsv.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<22), 1<<22)
	for sc.Scan() {
		want = append(want, sc.Text())
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}

	gson := NewGson()
	var got []string
	inputs, numbers := 0, 0
	for _, line := range want {
		fields := strings.Split(line, "\t")
		switch fields[0] {
		case "J":
			inputs++
			input, _ := gsonUnescape(t, fields[1])
			var e JsonElement
			out := "J\t" + fields[1]
			out += "\t" + gsonEscape(gsonRun(func() any { e = gson.FromJsonElement(input); return e }))
			out += "\t" + gsonEscape(gsonRun(func() any { return gson.FromJsonObject(input) }))
			out += "\t" + gsonEscape(gsonRun(func() any { return gson.FromJsonArray(input) }))
			got = append(got, out)
			if e != nil {
				got = append(got, "V\t\t"+gsonEscape(gsonAccessors(e)))
				if e.IsJsonObject() {
					for k, v := range e.GetAsJsonObject().EntrySet().All() {
						got = append(got, "V\t"+gsonEscape(k)+"\t"+gsonEscape(gsonAccessors(v)))
					}
				} else if e.IsJsonArray() {
					for m := range e.GetAsJsonArray().All() {
						got = append(got, "V\t[]\t"+gsonEscape(gsonAccessors(m)))
					}
				}
			}
		case "N":
			numbers++
			s, _ := gsonUnescape(t, fields[1])
			n := NewLazilyParsedNumber(s)
			got = append(got, "N\t"+fields[1]+"\t"+gsonEscape(gsonRun(func() any { return n.IntValue() }))+
				"\t"+gsonEscape(gsonRun(func() any { return n.LongValue() })))
		}
	}
	if inputs < 1500 || numbers < 50 {
		t.Fatalf("only %d inputs, %d numbers", inputs, numbers)
	}
	diffs := 0
	for i := 0; i < len(want) || i < len(got); i++ {
		var w, g string
		if i < len(want) {
			w = want[i]
		}
		if i < len(got) {
			g = got[i]
		}
		if w != g {
			diffs++
			if diffs <= 10 {
				t.Errorf("line %d:\n java: %s\n   go: %s", i+1, w, g)
			}
		}
	}
	if diffs > 0 {
		t.Fatalf("%d of %d lines differ", diffs, len(want))
	}
}
