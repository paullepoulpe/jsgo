package jutil

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

// matcherVectorUnescape undoes the escaping of
// tools/oracle/java/JutilMatcherVectors.java.
func matcherVectorUnescape(s string) string {
	var units []uint16
	for i := 0; i < len(s); {
		c := s[i]
		if c != '\\' {
			units = append(units, uint16(c))
			i++
			continue
		}
		switch s[i+1] {
		case '\\':
			units = append(units, '\\')
			i += 2
		case 't':
			units = append(units, '\t')
			i += 2
		case 'n':
			units = append(units, '\n')
			i += 2
		case 'r':
			units = append(units, '\r')
			i += 2
		case 'u':
			v, err := strconv.ParseUint(s[i+2:i+6], 16, 16)
			if err != nil {
				panic(err)
			}
			units = append(units, uint16(v))
			i += 6
		default:
			panic("bad escape in " + s)
		}
	}
	return StringFromUTF16(units)
}

func matcherVectorEscape(s *string) string {
	if s == nil {
		return "<null>"
	}
	var sb strings.Builder
	for _, c := range StringToUTF16(*s) {
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

// TestMatcherMatchesJava checks Matcher against JDK 21 vectors
// (testdata/matcher.tsv, from tools/oracle/java/JutilMatcherVectors.java).
func TestMatcherMatchesJava(t *testing.T) {
	f, err := os.Open("testdata/matcher.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		cols := strings.Split(sc.Text(), "\t")
		pattern := Pattern_Compile(matcherVectorUnescape(cols[1]))
		input := matcherVectorUnescape(cols[2])
		n++
		switch cols[0] {
		case "find":
			var sb strings.Builder
			got := "ok"
			Try(func() {
				m := pattern.Matcher(input)
				for m.Find() {
					fmt.Fprintf(&sb, "%d,%d", m.Start(), m.End())
					for g := 1; g <= m.GroupCount(); g++ {
						fmt.Fprintf(&sb, ",%s@%d-%d", matcherVectorEscape(m.Group(g)), m.StartGroup(g), m.EndGroup(g))
					}
					sb.WriteByte(';')
				}
			}, func(e IUnsupportedOperationException) { got = "unsupported" })
			want := ""
			if len(cols) > 3 {
				want = cols[3]
			}
			if got == "unsupported" {
				// Resuming inside a surrogate pair after an empty match is
				// not expressible on WTF-8 offsets (matcher.go); the vector
				// must indeed resume there in Java.
				if !strings.Contains(cols[2], `\ud8`) && !strings.Contains(cols[2], `\ud9`) &&
					!strings.Contains(cols[2], `\uda`) && !strings.Contains(cols[2], `\udb`) {
					t.Errorf("find %s on %s: unexpected UnsupportedOperationException", cols[1], cols[2])
				}
				continue
			}
			if sb.String() != want {
				t.Errorf("find %s on %s: got %q, Java %q", cols[1], cols[2], sb.String(), want)
			}
		case "match":
			m := pattern.Matcher(input)
			matches := "false"
			if m.Matches() {
				matches = "true:" + strconv.Itoa(m.End())
			}
			m.Reset()
			lookingAt := "false"
			if m.LookingAt() {
				lookingAt = "true:" + strconv.Itoa(m.End())
			}
			if matches != cols[3] || lookingAt != cols[4] {
				t.Errorf("match %s on %s: got %s %s, Java %s %s", cols[1], cols[2], matches, lookingAt, cols[3], cols[4])
			}
		case "replace":
			repl := matcherVectorUnescape(cols[3])
			result := func(f func() string) (s string) {
				Try(func() { r := f(); s = matcherVectorEscape(&r) }, func(e IRuntimeException) {
					msg := e.Error()
					s = "!" + matcherVectorEscape(&msg)
				})
				return s
			}
			first := result(func() string { return pattern.Matcher(input).ReplaceFirst(repl) })
			all := result(func() string { return pattern.Matcher(input).ReplaceAll(repl) })
			if first != cols[4] || all != cols[5] {
				t.Errorf("replace %s on %s with %s: got %s %s, Java %s %s", cols[1], cols[2], cols[3], first, all, cols[4], cols[5])
			}
		default:
			t.Fatalf("bad line %q", sc.Text())
		}
	}
	if n == 0 {
		t.Fatal("no vectors")
	}
}
