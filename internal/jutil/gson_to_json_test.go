package jutil

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
)

// testdata/gson_tojson.tsv is tools/oracle/java/GsonToJsonVectors.java:
// Gson.toJson of the values below (built in Java by the same case names)
// under four configurations, from the Gson 2.9.1 shaded into the jar.

type gsonTestColor int32

const (
	gsonTestColor_RED gsonTestColor = iota
	gsonTestColor_GREEN
)

func (c gsonTestColor) String() string { return [...]string{"RED", "GREEN"}[c] }

// Java: class Point { int x; String name; String missing; transient int t;
// static int s; List<Integer> list; Object me = this; double d; char c;
// boolean b; long l; }
type GsonTestPoint struct {
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

// Java: class Point3 extends Point { int z; Map<String, Integer> m; }
type gsonTestPoint3 struct {
	GsonTestPoint
	Z int32                          `gson:"z"`
	M *LinkedHashMap[string, int32] `gson:"m"`
}

type gsonTestHolder struct {
	Value any `gson:"value"`
}

type gsonTestAble struct {
	V any `gson:"v"`
}

func newGsonTestPoint() *GsonTestPoint {
	p := &GsonTestPoint{X: 1, Name: "p<q>", T: 5, List: []int32{1, 2}, D: 0.5, C: 'c', B: true, L: 1 << 40}
	p.Me = p
	return p
}

func newGsonTestPoint3() *gsonTestPoint3 {
	p := &gsonTestPoint3{Z: 9, M: NewLinkedHashMap[string, int32]()}
	p.GsonTestPoint = GsonTestPoint{X: 1, Name: "p<q>", T: 5, List: []int32{1, 2}, D: 0.5, C: 'c', B: true, L: 1 << 40}
	p.Me = p
	p.M.Put("k", 1)
	return p
}

func gsonTestCases() map[string]any {
	str := func(s string) *string { return &s }
	m := NewLinkedHashMap[string, any]()
	m.Put("a", int32(1))
	m.Put("b", nil)
	m.Put("c", []int32{1, 2})
	m.Put("d", "<")
	mi := NewLinkedHashMap[int32, string]()
	mi.Put(3, "x")
	mi.Put(1, "y")
	mn := NewLinkedHashMap[*string, int32]()
	mn.Put(nil, 1)
	me := NewLinkedHashMap[gsonTestColor, int32]()
	me.Put(gsonTestColor_RED, 1)
	ableMap := NewLinkedHashMap[string, int32]()
	ableMap.Put("a", 1)
	o := NewJsonObject()
	o.Add("a", NewJsonPrimitiveNumber(gson_javaInteger(1)))
	arr := NewJsonArray()
	arr.Add(NewJsonPrimitiveBoolean(true))
	arr.Add(JsonNull_INSTANCE)
	o.Add("b", arr)
	o.Add("c", NewJsonPrimitiveString("<x>"))
	inner := NewLinkedHashMap[string, any]()
	inner.Put("x", []any{NewLinkedHashMap[string, any]()})
	tm := NewTreeMap[string, int32](StringCompareTo)
	tm.Put("b", 1)
	tm.Put("a", 2)
	return map[string]any{
		"null":            nil,
		"string":          "a<b>&'=\"\\  \u0001\t\n\b\f\r/é\U0001F600",
		"loneSurrogate":   StringFromUTF16([]uint16{0xd800, 'x', 0xdc00}),
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
		"enum":            gsonTestColor_GREEN,
		"enumList":        []gsonTestColor{gsonTestColor_RED, gsonTestColor_GREEN},
		"enumKeyMap":      me,
		"point":           newGsonTestPoint(),
		"point3":          newGsonTestPoint3(),
		"holderNull":      &gsonTestHolder{Value: nil},
		"holderPoint":     &gsonTestHolder{Value: newGsonTestPoint()},
		"holderList":      &gsonTestHolder{Value: []int32{1}},
		"ableList":        &gsonTestAble{V: []int32{0}},
		"ableNull":        &gsonTestAble{V: nil},
		"ableString":      &gsonTestAble{V: "<"},
		"ableInList":      []any{&gsonTestAble{V: ableMap}},
		"holderAble":      &gsonTestHolder{Value: &gsonTestAble{V: int32(1)}},
		"jsonTree":        o,
		"nested":          []any{inner},
		"treeMap":         tm,
		"linkedHashSet":   NewLinkedHashSet("b", "a"),
		"arrayList":       NewArrayListFrom([]int32{3, 1}),
	}
}

func gsonTestConfigs() map[string]*Gson {
	return map[string]*Gson{
		"default":        NewGson(),
		"builder":        NewGsonBuilder().Create(),
		"noHtmlPretty":   NewGsonBuilder().DisableHtmlEscaping().SetPrettyPrinting().Create(),
		"serializeNulls": NewGsonBuilder().SerializeNulls().Create(),
	}
}

func TestGsonToJsonMatchesJava(t *testing.T) {
	f, err := os.Open("testdata/gson_tojson.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	cases := gsonTestCases()
	configs := gsonTestConfigs()
	n := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), "\t", 3)
		if parts[1] == "enumSet" {
			continue
		}
		gson, ok := configs[parts[0]]
		if !ok {
			t.Fatalf("unknown configuration %q", parts[0])
		}
		value, ok := cases[parts[1]]
		if !ok {
			t.Fatalf("unknown case %q", parts[1])
		}
		n++
		var got string
		Try(func() { got = gson.ToJson(value) }, func(e Throwable) { got = gsonTestExceptionString(e) })
		if got := gsonTestEscape(got); got != parts[2] {
			t.Errorf("%s %s: got %s, Java %s", parts[0], parts[1], got, parts[2])
		}
	}
	if n < 196 {
		t.Fatalf("only %d vectors", n)
	}
}

func gsonTestExceptionString(e Throwable) string {
	s := e.Error()
	if e.GetMessage() == "" && !strings.Contains(s, ": ") {
		s += ": null"
	}
	return s
}

// gsonTestEscape mirrors RegExpDump.escape.
func gsonTestEscape(s string) string {
	var b strings.Builder
	for _, c := range StringToUTF16(s) {
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

func TestGsonTypeHierarchyAdapterOrder(t *testing.T) {
	// GsonBuilder.create() reverses the hierarchy factories: the adapter
	// registered last wins for a value matching both.
	gson := NewGsonBuilder().
		RegisterTypeHierarchyAdapter(func(v any) bool { _, ok := v.(string); return ok }, gsonTestRawAdapter("first")).
		RegisterTypeHierarchyAdapter(func(v any) bool { _, ok := v.(string); return ok }, gsonTestRawAdapter("second")).
		Create()
	if got := gson.ToJson([]any{"x", int32(1)}); got != "[second,1]" {
		t.Errorf("got %s", got)
	}
}

type gsonTestRawAdapter string

func (a gsonTestRawAdapter) Write(out *JsonWriter, value any) {
	s := string(a)
	out.JsonValue(&s)
}

func (a gsonTestRawAdapter) Read(in *JsonReader) any { panic(NewAssertionErrorNoMessage()) }

func TestGsonToJsonUnsupportedGoKinds(t *testing.T) {
	for _, v := range []any{func() {}, map[string]int{}, uint32(1), float32(1)} {
		caught := false
		Try(func() { NewGson().ToJson(v) }, func(e IUnsupportedOperationException) { caught = true })
		if !caught {
			t.Errorf("%T: no UnsupportedOperationException", v)
		}
	}
}
