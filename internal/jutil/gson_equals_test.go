package jutil

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/GsonEqualsVectors.java against
// the Gson 2.9.1 shaded into the pinned closure-compiler.jar: equals for
// every ordered pair of a list of JsonElements, and their hashCodes.

// gsonEqualsElement rebuilds a GsonEqualsVectors element.
func gsonEqualsElement(t *testing.T, kind string, payload string) JsonElement {
	switch kind {
	case "json":
		return NewGson().FromJsonElement(payload)
	case "int":
		v, err := strconv.ParseInt(payload, 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		return NewJsonPrimitiveNumber(Integer(v))
	case "long":
		v, err := strconv.ParseInt(payload, 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		return NewJsonPrimitiveNumber(gson_javaInteger(v))
	case "double":
		return NewJsonPrimitiveNumber(gson_javaDouble(DoubleParseDouble(payload)))
	case "string":
		return NewJsonPrimitiveString(payload)
	case "bool":
		return NewJsonPrimitiveBoolean(payload == "true")
	case "lazy":
		return NewJsonPrimitiveNumber(NewLazilyParsedNumber(payload))
	case "null":
		return JsonNull_INSTANCE
	case "objectOfInt":
		o := NewJsonObject()
		o.Add("a", NewJsonPrimitiveNumber(Integer(1)))
		return o
	case "arrayOfLongAndDouble":
		a := NewJsonArray()
		a.Add(NewJsonPrimitiveNumber(gson_javaInteger(1)))
		a.Add(NewJsonPrimitiveNumber(gson_javaDouble(2)))
		return a
	}
	t.Fatalf("unknown kind %q", kind)
	return nil
}

func TestGsonEqualsMatchesJava(t *testing.T) {
	f, err := os.Open("testdata/gson_equals.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var elements []JsonElement
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	hashes, pairs := 0, 0
	for sc.Scan() {
		fields := strings.Split(sc.Text(), "\t")
		i, err := strconv.Atoi(fields[1])
		if err != nil {
			t.Fatal(err)
		}
		switch fields[0] {
		case "D":
			payload, _ := gsonUnescape(t, fields[3])
			if i != len(elements) {
				t.Fatalf("element %d out of order", i)
			}
			elements = append(elements, gsonEqualsElement(t, fields[2], payload))
		case "H":
			if fields[2] == "-" {
				continue
			}
			hashes++
			if got := strconv.Itoa(int(elements[i].HashCode())); got != fields[2] {
				t.Errorf("hashCode(%s) = %s, want %s", elements[i], got, fields[2])
			}
		case "E":
			for j, want := range fields[2] {
				pairs++
				if got := elements[i].Equals(elements[j]); got != (want == '1') {
					t.Errorf("%s.equals(%s) = %t, want %t", elements[i], elements[j], got, want == '1')
				}
			}
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	if hashes < 80 || pairs < 8000 {
		t.Fatalf("only %d hashes, %d pairs", hashes, pairs)
	}
}
