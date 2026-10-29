package jutil

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
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
			r[4], strconv.Itoa(int(b.HashCode())),
			strconv.Itoa(b.NextSetBit(from)), r[7], r[8]}
		if strings.Join(got, "\t") != strings.Join(r, "\t") {
			t.Errorf("bitset:\n got %v\nwant %v", got, r)
		}
	}
}

func TestLoggerLevelsAndHandlers(t *testing.T) {
	l := Logger_GetLogger("com.example.jutil.Test")
	h := &recordingHandler{}
	l.addHandler(h)
	l.setUseParentHandlers(false)
	defer l.removeHandler(h)
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

func (h *recordingHandler) publish(r *logRecordEntry) {
	h.got = append(h.got, r.getLevel().GetName()+":"+formatter_FormatMessage(r))
}

type noopHandler struct{}

func (noopHandler) publish(*logRecordEntry) {}

func TestLoggerConcurrentGetLoggerAndParse(t *testing.T) {
	const workers = 16
	const iters = 100
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func(id int) {
			defer wg.Done()
			parent := Logger_GetLogger(fmt.Sprintf("com.example.concurrent.p%d", id%4))
			parent.setUseParentHandlers(false)
			for i := 0; i < iters; i++ {
				child := Logger_GetLogger(fmt.Sprintf("com.example.concurrent.p%d.c%d", id%4, i%8))
				h := &noopHandler{}
				child.addHandler(h)
				if i%2 == 0 {
					parent.SetLevel(Level_FINE)
				} else {
					parent.SetLevel(Level_WARNING)
				}
				_ = child.IsLoggable(Level_INFO)
				_ = child.GetParent()
				child.Fine("concurrent fine")
				if lv := Level_Parse("WARNING"); lv != Level_WARNING {
					t.Errorf("Level_Parse(WARNING) = %v", lv)
				}
				if lv := Level_Parse("850"); lv.GetName() != "850" || lv.value != 850 {
					t.Errorf("Level_Parse(850) = %v/%d", lv.GetName(), lv.value)
				}
				child.removeHandler(h)
			}
		}(w)
	}
	wg.Wait()
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
