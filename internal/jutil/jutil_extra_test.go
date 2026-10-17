package jutil

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilExtraVectors.java against
// the pinned closure-compiler.jar (JDK 21, shaded Guava 33.4.6).

func decodeUnits(t *testing.T, s string) []uint16 {
	t.Helper()
	if !strings.HasPrefix(s, "u") {
		t.Fatalf("bad units %q", s)
	}
	s = s[1:]
	out := make([]uint16, len(s)/4)
	for i := range out {
		v, err := strconv.ParseUint(s[4*i:4*i+4], 16, 16)
		if err != nil {
			t.Fatal(err)
		}
		out[i] = uint16(v)
	}
	return out
}

func encodeUnits(u []uint16) string {
	var sb strings.Builder
	sb.WriteByte('u')
	for _, c := range u {
		fmt.Fprintf(&sb, "%04x", c)
	}
	return sb.String()
}

func decodeHexBytes(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(strings.TrimPrefix(s, "x"))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestRandomMatchesJava(t *testing.T) {
	rows := readTSV(t, "random.tsv")
	bounds := []int32{1, 2, 3, 7, 16, 52, 64, 100, 1000, 1 << 30, 2147483647, 0x40000001}
	for _, r := range rows {
		seed, _ := strconv.ParseInt(r[0], 10, 64)
		x := NewRandomWithSeed(seed)
		var parts []string
		for i := 0; i < 5; i++ {
			parts = append(parts, strconv.Itoa(int(x.NextInt())))
		}
		for _, b := range bounds {
			parts = append(parts, strconv.Itoa(int(x.NextIntBound(b))))
		}
		for i := 0; i < 3; i++ {
			parts = append(parts, strconv.FormatInt(x.NextLong(), 10))
		}
		parts = append(parts, strconv.FormatBool(x.NextBoolean()))
		parts = append(parts, strconv.FormatInt(DoubleToRawLongBits(x.NextDouble()), 10))
		parts = append(parts, strconv.Itoa(int(int32(floatBits(x.NextFloat())))))
		if got := strings.Join(parts, ","); got != r[1] {
			t.Errorf("seed %d:\n got %s\nwant %s", seed, got, r[1])
		}
		list := make([]int, 64)
		for i := range list {
			list[i] = i
		}
		Collections_Shuffle(list, x)
		if got := strings.ReplaceAll(CollectionToString(list), " ", ""); got != r[2] {
			t.Errorf("seed %d shuffle:\n got %s\nwant %s", seed, got, r[2])
		}
	}
}

func TestMurmur3MatchesGuava(t *testing.T) {
	rows := readTSV(t, "murmur3.tsv")
	for _, r := range rows {
		switch r[0] {
		case "B":
			got := Hashing_Murmur3_128NewHasher().PutBytes(decodeHexBytes(t, r[1])).Hash().String()
			if got != r[2] {
				t.Errorf("murmur3(%s) = %s, Guava %s", r[1], got, r[2])
			}
		case "C":
			comma := strings.IndexByte(r[1], ',')
			n, _ := strconv.Atoi(r[1][:comma])
			s := StringFromUTF16(decodeUnits(t, r[1][comma+1:]))
			got := Hashing_Murmur3_128NewHasher().PutInt(int32(n)).PutUnencodedChars(s).Hash().AsInt()
			if strconv.Itoa(int(got)) != r[2] {
				t.Errorf("murmur3(%s) = %d, Guava %s", r[1], got, r[2])
			}
		}
	}
}

func TestUTF8MatchesJava(t *testing.T) {
	rows := readTSV(t, "utf8.tsv")
	for _, r := range rows {
		switch r[0] {
		case "E":
			got := StringGetBytesUTF8(StringFromUTF16(decodeUnits(t, r[1])))
			if !bytes.Equal(got, decodeHexBytes(t, r[2])) {
				t.Errorf("getBytes(%s) = %x, Java %s", r[1], got, r[2])
			}
		case "D":
			got := encodeUnits(StringToUTF16(NewStringFromUTF8Bytes(decodeHexBytes(t, r[1]))))
			if got != r[2] {
				t.Errorf("new String(%s) = %s, Java %s", r[1], got, r[2])
			}
		}
	}
}

func TestCaseFormatLowerCamelToUpperUnderscoreMatchesGuava(t *testing.T) {
	for _, r := range readTSV(t, "caseformat.tsv") {
		in := StringFromUTF16(decodeUnits(t, r[0]))
		if got := encodeUnits(StringToUTF16(CaseFormat_LowerCamelToUpperUnderscore(in))); got != r[1] {
			t.Errorf("LOWER_CAMEL.to(UPPER_UNDERSCORE, %s) = %s, Guava %s", r[0], got, r[1])
		}
	}
}

func TestCharsetMatchesJava(t *testing.T) {
	for _, r := range readTSV(t, "charset.tsv") {
		c := Charset_ForName(r[0])
		if c.Name() != r[0] {
			t.Errorf("forName(%s).name() = %s", r[0], c.Name())
		}
		for _, a := range strings.Split(r[1], ",") {
			if got := Charset_ForName(strings.ToUpper(a)); got != c {
				t.Errorf("forName(%s) = %s, want %s", a, got, c)
			}
		}
		var ranges strings.Builder
		start := -1
		for ch := 0; ch <= 0x10000; ch++ {
			ok := ch < 0x10000 && c.CanEncodeChar(uint16(ch))
			if ok && start < 0 {
				start = ch
			}
			if !ok && start >= 0 {
				fmt.Fprintf(&ranges, "%x-%x,", start, ch-1)
				start = -1
			}
		}
		if ranges.String() != r[2] {
			t.Errorf("%s canEncode ranges %s, Java %s", r[0], ranges.String(), r[2])
		}
	}
	Try(func() {
		Charset_ForName("x-no-such-charset")
		t.Error("expected UnsupportedCharsetException")
	}, func(e IIllegalArgumentException) {})
}

func TestBitSetMatchesJava(t *testing.T) {
	for _, r := range readTSV(t, "bitset.tsv") {
		b := NewBitSet()
		for _, op := range strings.Split(strings.TrimSuffix(r[0], ","), ",") {
			if op == "" {
				continue
			}
			n, _ := strconv.Atoi(op[1:])
			if op[0] == 's' {
				b.Set(n)
			} else {
				b.Clear(n)
			}
		}
		from, _ := strconv.Atoi(r[8])
		got := []string{r[0], b.String(), strconv.Itoa(b.Length()), strconv.Itoa(b.Cardinality()),
			"x" + hex.EncodeToString(b.ToByteArray()), strconv.Itoa(int(b.HashCode())),
			strconv.Itoa(b.NextSetBit(from)), strconv.Itoa(b.NextClearBit(from)), r[8]}
		if strings.Join(got, "\t") != strings.Join(r, "\t") {
			t.Errorf("bitset:\n got %v\nwant %v", got, r)
		}
	}
}

func TestBiMapMatchesGuava(t *testing.T) {
	for _, r := range readTSV(t, "bimap.tsv") {
		switch r[0] {
		case "H":
			m := HashBiMap_Create[string, string]()
			var results strings.Builder
			for _, op := range strings.Split(strings.TrimSuffix(r[1], ","), ",") {
				if op == "" {
					continue
				}
				k, v := op[1:2], op[2:3]
				res := ""
				Try(func() {
					var old string
					var had bool
					switch op[0] {
					case 'F':
						old, had = m.Lookup(k)
						m.ForcePut(k, v)
					case 'R':
						old, had = m.Lookup(k)
						m.Remove(k)
					default:
						old, had = m.Lookup(k)
						m.Put(k, v)
					}
					if had {
						res = old
					} else {
						res = "null"
					}
				}, func(e IIllegalArgumentException) { res = "!" + e.GetMessage() })
				results.WriteString(res + ";")
			}
			if got := results.String(); got != r[2] {
				t.Errorf("ops %s: results %s, Guava %s", r[1], got, r[2])
			}
			if m.String() != r[3] || m.Inverse().String() != r[4] {
				t.Errorf("ops %s: %s / %s, Guava %s / %s", r[1], m, m.Inverse(), r[3], r[4])
			}
		case "I":
			src := NewLinkedHashMap[string, string]()
			for _, op := range strings.Split(strings.TrimSuffix(r[1], ","), ",") {
				if op != "" {
					src.Put(op[0:1], op[1:2])
				}
			}
			got := ""
			Try(func() {
				m := ImmutableBiMap_CopyOf(src)
				got = m.String() + "\t" + m.Inverse().String()
			}, func(e IIllegalArgumentException) { got = "!" + e.GetMessage() + "\t-" })
			if got != r[2]+"\t"+r[3] {
				t.Errorf("copyOf(%s) = %s, Guava %s\t%s", r[1], got, r[2], r[3])
			}
		}
	}
}

func TestLoggerLevelsAndHandlers(t *testing.T) {
	l := Logger_GetLogger("com.example.jutil.Test")
	h := &recordingHandler{}
	l.AddHandler(h)
	l.SetUseParentHandlers(false)
	defer l.RemoveHandler(h)
	l.Fine("hidden")
	l.Info("shown")
	l.LogWithParam(Level_WARNING, "timers:\n{0}", "x")
	l.LogWithParam(Level_WARNING, "it''s '{0}' {0}", 7)
	want := []string{"INFO:shown", "WARNING:timers:\nx", "WARNING:it's {0} 7"}
	if strings.Join(h.got, "|") != strings.Join(want, "|") {
		t.Errorf("got %q, want %q", h.got, want)
	}
	if l.GetParent() != loggerRoot {
		t.Errorf("parent of %s is %s", l.GetName(), l.GetParent().GetName())
	}
}

type recordingHandler struct{ got []string }

func (h *recordingHandler) Publish(r *LogRecord) {
	h.got = append(h.got, r.GetLevel().GetName()+":"+Formatter_FormatMessage(r))
}

func TestBufferedReaderReadLine(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"a", []string{"a"}},
		{"a\n", []string{"a"}},
		{"a\n\nb", []string{"a", "", "b"}},
		{"a\r\nb\rc\n\r\n", []string{"a", "b", "c", ""}},
		{"\r", []string{""}},
	} {
		for _, r := range []*BufferedReader{
			NewBufferedReaderFromString(tc.in),
			NewBufferedReaderUTF8(strings.NewReader(tc.in)),
		} {
			var got []string
			for {
				line, ok := r.ReadLine()
				if !ok {
					break
				}
				got = append(got, line)
			}
			if strings.Join(got, "|") != strings.Join(tc.want, "|") || len(got) != len(tc.want) {
				t.Errorf("readLine(%q) = %q, want %q", tc.in, got, tc.want)
			}
		}
	}
}

func TestParseException(t *testing.T) {
	caught := false
	Try(func() { panic(NewParseException("Bad line: x", 0)) }, func(e IException) {
		caught = e.Error() == "java.text.ParseException: Bad line: x"
	})
	if !caught {
		t.Error("ParseException not caught as a checked exception")
	}
}

func TestSplitterOnCharOmitEmptyStrings(t *testing.T) {
	got := Splitter_OnCharOmitEmptyStrings("\na b\n\n c\n", '\n')
	want := []string{"a b", " c"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got := Splitter_OnCharOmitEmptyStrings("", ' '); len(got) != 0 {
		t.Fatalf("got %q", got)
	}
	if got := Splitter_OnCharOmitEmptyStrings("  x  y ", ' '); len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Fatalf("got %q", got)
	}
}

func TestSetMultimapContainsEntry(t *testing.T) {
	mm := NewSetMultimap[string, string]()
	mm.Put("a", "x")
	if !mm.ContainsEntry("a", "x") || mm.ContainsEntry("a", "y") || mm.ContainsEntry("b", "x") {
		t.Fatal("ContainsEntry")
	}
}
