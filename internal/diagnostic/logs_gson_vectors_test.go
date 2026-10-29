package diagnostic

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// testdata/logs_gson.tsv is the LogsGson part of
// tools/oracle/java/GsonToJsonVectors.java: LogsGson.toJson of the values
// below (built in Java by the same case names) from the pinned jar.

type testColor int32

const (
	testColor_RED testColor = iota
	testColor_GREEN
)

func (c testColor) String() string { return [...]string{"RED", "GREEN"}[c] }

type TestPoint struct {
	X       int32   `gson:"x"`
	Name    string  `gson:"name"`
	Missing *string `gson:"missing"`
	T       int32   `gson:"-"`
	List    []int32 `gson:"list"`
	Me      any     `gson:"me"`
	D       float64 `gson:"d"`
	C       uint16  `gson:"c"`
	B       bool    `gson:"b"`
	L       int64   `gson:"l"`
}

type testPoint3 struct {
	TestPoint
	Z int32                               `gson:"z"`
	M *jutil.LinkedHashMap[string, int32] `gson:"m"`
}

type testHolder struct {
	Value any `gson:"value"`
}

// testAbleThing is AbleThing implements LogsGson.Able.
type testAbleThing struct{ v any }

func (a *testAbleThing) ToLogsGson() any { return a.v }

func newTestPoint() *TestPoint {
	p := &TestPoint{X: 1, Name: "p<q>", T: 5, List: []int32{1, 2}, D: 0.5, C: 'c', B: true, L: 1 << 40}
	p.Me = p
	return p
}

func newTestPoint3() *testPoint3 {
	p := &testPoint3{Z: 9, M: jutil.NewLinkedHashMap[string, int32]()}
	p.TestPoint = TestPoint{X: 1, Name: "p<q>", T: 5, List: []int32{1, 2}, D: 0.5, C: 'c', B: true, L: 1 << 40}
	p.Me = p
	p.M.Put("k", 1)
	return p
}

func testCases() map[string]any {
	str := func(s string) *string { return &s }
	m := jutil.NewLinkedHashMap[string, any]()
	m.Put("a", int32(1))
	m.Put("b", nil)
	m.Put("c", []int32{1, 2})
	m.Put("d", "<")
	mi := jutil.NewLinkedHashMap[int32, string]()
	mi.Put(3, "x")
	mi.Put(1, "y")
	mn := jutil.NewLinkedHashMap[*string, int32]()
	mn.Put(nil, 1)
	me := jutil.NewLinkedHashMap[testColor, int32]()
	me.Put(testColor_RED, 1)
	ableMap := jutil.NewLinkedHashMap[string, int32]()
	ableMap.Put("a", 1)
	o := jutil.NewJsonObject()
	o.Add("a", jutil.NewJsonPrimitiveNumber(jutil.NewLazilyParsedNumber("1")))
	arr := jutil.NewJsonArray()
	arr.Add(jutil.NewJsonPrimitiveBoolean(true))
	arr.Add(jutil.JsonNull_INSTANCE)
	o.Add("b", arr)
	o.Add("c", jutil.NewJsonPrimitiveString("<x>"))
	inner := jutil.NewLinkedHashMap[string, any]()
	inner.Put("x", []any{jutil.NewLinkedHashMap[string, any]()})
	tm := jutil.NewLinkedHashMap[string, int32]()
	tm.Put("a", 2)
	tm.Put("b", 1)
	multimap := newTestMultimap()
	multimap.Put("foo", int32(0))
	multimap.Put("foo", int32(1))
	multimap.Put("bar", int32(2))
	multimap2 := newTestMultimap()
	multimap2.Put("<", "&")
	return map[string]any{
		"null":            nil,
		"string":          "a<b>&'=\"\\  \u0001\t\n\b\f\r/é\U0001F600",
		"loneSurrogate":   jutil.StringFromUTF16([]uint16{0xd800, 'x', 0xdc00}),
		"int":             int32(42),
		"long":            int64(-7),
		"short":           int16(3),
		"byte":            int8(-1),
		"intMin":          int32(math.MinInt32),
		"longMax":         int64(math.MaxInt64),
		"double1":         1.0,
		"double0.1":       0.1,
		"double1e21":      1e21,
		"doubleNegZero":   math.Copysign(0, -1),
		"double1e-7":      1e-7,
		"double123456789": 123456789.0,
		"double1e7":       1e7,
		"doubleMin":       math.SmallestNonzeroFloat64,
		"doubleNaN":       math.NaN(),
		"doubleInf":       math.Inf(-1),
		"char":            uint16('x'),
		"charLS":          uint16(0x2028),
		"charQuote":       uint16('"'),
		"true":            true,
		"listMixed":       []any{int32(1), "a", nil, true},
		"listEmpty":       []any{},
		"listNested":      [][]int32{{1}, {}},
		"intArray":        []int32{1, 2, 3},
		"stringArray":     []*string{str("x"), nil},
		"map":             m,
		"mapIntKeys":      mi,
		"mapNullKey":      mn,
		"enum":            testColor_GREEN,
		"enumList":        []testColor{testColor_RED, testColor_GREEN},
		"enumKeyMap":      me,
		"point":           newTestPoint(),
		"point3":          newTestPoint3(),
		"holderNull":      &testHolder{Value: nil},
		"holderPoint":     &testHolder{Value: newTestPoint()},
		"holderList":      &testHolder{Value: []int32{1}},
		"ableList":        &testAbleThing{v: []int32{0}},
		"ableNull":        &testAbleThing{v: nil},
		"ableString":      &testAbleThing{v: "<"},
		"ableInList":      []any{&testAbleThing{v: ableMap}},
		"holderAble":      &testHolder{Value: &testAbleThing{v: int32(1)}},
		"jsonTree":        o,
		"nested":          []any{inner},
		"treeMap":         tm,
		"linkedHashSet":   jutil.NewLinkedHashSet("b", "a"),
		"arrayList":       jutil.NewArrayListFrom([]int32{3, 1}),
		"multimap":        multimap,
		"holderMultimap":  &testHolder{Value: multimap2},
	}
}

func TestLogsGsonMatchesJava(t *testing.T) {
	f, err := os.Open("testdata/logs_gson.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	cases := testCases()
	n := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), "\t", 3)
		if parts[1] == "enumSet" {
			continue
		}
		value, ok := cases[parts[1]]
		if !ok {
			t.Fatalf("unknown case %q", parts[1])
		}
		n++
		var got string
		jutil.Try(func() { got = LogsGson_ToJson(value) }, func(e jutil.Throwable) { got = exceptionString(e) })
		if got := escape(got); got != parts[2] {
			t.Errorf("%s: got %s, Java %s", parts[1], got, parts[2])
		}
	}
	if n != len(cases) {
		t.Fatalf("%d vectors for %d cases", n, len(cases))
	}
}

func exceptionString(e jutil.Throwable) string {
	s := e.Error()
	if e.GetMessage() == "" && !strings.Contains(s, ": ") {
		s += ": null"
	}
	return s
}

// escape mirrors RegExpDump.escape.
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
