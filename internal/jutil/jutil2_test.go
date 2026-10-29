package jutil

import (
	"bufio"
	"encoding/hex"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilVectors.java against the
// pinned closure-compiler.jar (JDK 21, shaded Guava).

func readTSV(t *testing.T, name string) [][]string {
	t.Helper()
	f, err := os.Open("testdata/" + name)
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
	return rows
}

func TestFarmHashFingerprint64MatchesGuava(t *testing.T) {
	rows := readTSV(t, "farmhash.tsv")
	if len(rows) < 900 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		want, _ := strconv.ParseInt(r[1], 10, 64)
		var got int64
		if strings.HasPrefix(r[0], "L") {
			h := Hashing_FarmHashFingerprint64NewHasher()
			for _, s := range strings.Split(r[0][1:], ",") {
				v, _ := strconv.ParseInt(s, 10, 64)
				h.PutLong(v)
			}
			got = h.HashAsLong()
		} else {
			b, _ := hex.DecodeString(r[0][1:])
			got = int64(FarmHashFingerprint64(b))
		}
		if got != want {
			t.Errorf("fingerprint(%s) = %d, Guava %d", r[0], got, want)
		}
	}
}

func TestCharacterCodePointMatchesJava(t *testing.T) {
	rows := readTSV(t, "character_codepoint.tsv")
	b := func(c byte) bool { return c == '1' }
	for _, r := range rows {
		lo, _ := strconv.ParseInt(r[0], 16, 32)
		hi, _ := strconv.ParseInt(r[1], 16, 32)
		parts := strings.Split(r[2], ",")
		typ, _ := strconv.Atoi(parts[0])
		fl := parts[1]
		for cp := int(lo); cp <= int(hi); cp++ {
			if got := CharacterGetTypeOfCodePoint(cp); got != typ {
				t.Fatalf("getType(%x) = %d, Java %d", cp, got, typ)
			}
			if CharacterIsWhitespaceCodePoint(cp) != b(fl[0]) ||
				CharacterIsLetterCodePoint(cp) != b(fl[4]) || CharacterIsDigitCodePoint(cp) != b(fl[5]) {
				t.Fatalf("flags(%x) mismatch, Java %s", cp, fl)
			}
			if cp <= 0xFFFF {
				ch := uint16(cp)
				if CharacterIsJavaIdentifierStart(ch) != b(fl[1]) || CharacterIsJavaIdentifierPart(ch) != b(fl[2]) ||
					CharacterIsIdentifierIgnorable(ch) != b(fl[3]) {
					t.Fatalf("identifier flags(%x) mismatch, Java %s", cp, fl)
				}
			}
		}
	}
}

func TestBigIntegerMatchesJava(t *testing.T) {
	for _, r := range readTSV(t, "bigint.tsv") {
		x := NewBigIntegerFromString(r[0])
		if got := "x" + hex.EncodeToString(BigIntegerToByteArray(x)); got != r[1] {
			t.Errorf("toByteArray(%s) = %s, Java %s", r[0], got, r[1])
		}
		bs, _ := hex.DecodeString(r[1][1:])
		if back := NewBigIntegerFromBytes(bs); back.Cmp(x) != 0 {
			t.Errorf("new BigInteger(%s) = %s, want %s", r[1], back, r[0])
		}
		want, _ := strconv.ParseInt(r[2], 10, 32)
		if got := BigIntegerIntValue(x); int64(got) != want {
			t.Errorf("intValue(%s) = %d, Java %d", r[0], got, want)
		}
		if x.String() != r[0] {
			t.Errorf("toString = %s, want %s", x.String(), r[0])
		}
	}
}

func TestObjectsEquals(t *testing.T) {
	if !ObjectsEquals(nil, nil) || ObjectsEquals(nil, "a") || ObjectsEquals("a", nil) {
		t.Fatal("null handling")
	}
	if !ObjectsEquals(math.NaN(), math.NaN()) || ObjectsEquals(0.0, math.Copysign(0, -1)) {
		t.Fatal("Double.equals semantics")
	}
	if ObjectsEquals(int32(1), int64(1)) || !ObjectsEquals("x", "x") {
		t.Fatal("boxed equality")
	}
	l1 := NewArrayListFrom([]string{"a", "b"})
	l2 := NewArrayListFrom([]string{"a", "b"})
	if !ObjectsEquals(l1, l2) || l1.HashCode() != l2.HashCode() {
		t.Fatal("List.equals")
	}
	s1 := NewLinkedHashSet("a", "b")
	s2 := NewLinkedHashSet("b", "a")
	if !s1.Equals(s2) || s1.HashCode() != StringHashCode("a")+StringHashCode("b") {
		t.Fatal("Set.equals / hashCode")
	}
	m := NewObjectLinkedHashMap[*LinkedHashSet[string], string]()
	m.Put(s1, "x")
	if m.Get(s2) != "x" || m.Size() != 1 {
		t.Fatal("ObjectLinkedHashMap keyed by value-equal sets")
	}
	m.Put(NewLinkedHashSet("c"), "y")
	if m.String() != "{[a, b]=x, [c]=y}" {
		t.Fatalf("String() = %s", m.String())
	}
	m.Remove(s2)
	if m.Size() != 1 || m.ContainsKey(s1) {
		t.Fatal("remove")
	}
}

func TestToStringHelper(t *testing.T) {
	got := MoreObjects_ToStringHelper("JSDocInfo").Add("bitset", nil).AddBool("bit:const", true).
		Add("type", "x").OmitNullValues().String()
	if got != "JSDocInfo{bit:const=true, type=x}" {
		t.Fatalf("got %q", got)
	}
}

func TestCaseFormat(t *testing.T) {
	for in, want := range map[string]string{
		"INLINE_TYPE": "inlineType", "CONST": "const", "PURE_OR_BREAK_MY_CODE": "pureOrBreakMyCode",
		"NODTS": "nodts",
	} {
		if got := CaseFormat_UpperUnderscoreToLowerCamel(in); got != want {
			t.Errorf("%s -> %s, want %s", in, got, want)
		}
	}
}

func TestBigIntegerOverflow(t *testing.T) {
	caught := ""
	Try(func() { BigIntegerPow(NewBigIntegerFromString("5"), 800000000) }, func(e IArithmeticException) { caught = e.GetMessage() })
	if caught != "BigInteger would overflow supported range" {
		t.Fatalf("caught %q", caught)
	}
	if BigIntegerPow(NewBigIntegerFromString("5"), 3).String() != "125" {
		t.Fatal("pow")
	}
}
