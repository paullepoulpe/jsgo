package gson

import (
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
)

func TestJsonStreamRoundTrip(t *testing.T) {
	sw := javaio.NewStringWriter()
	jw := NewJsonWriter(sw)
	jw.BeginArray()
	jw.BeginObject()
	jw.Name("src").ValueString("const a = 1 < 2;\n")
	jw.Name("path").ValueString("a.js")
	jw.EndObject()
	jw.EndArray()
	jw.Close()

	got := sw.String()
	want := `[{"src":"const a = 1 < 2;\n","path":"a.js"}]`
	if got != want {
		t.Fatalf("JsonWriter got %s, want %s", got, want)
	}

	jr := NewJsonReader(javaio.NewStringReader(got))
	jr.BeginArray()
	if !jr.HasNext() {
		t.Fatal("expected object in array")
	}
	jr.BeginObject()
	k1 := jr.NextName()
	v1 := jr.NextString()
	k2 := jr.NextName()
	v2 := jr.NextString()
	jr.EndObject()
	jr.EndArray()
	jr.Close()

	if k1 != "src" || v1 != "const a = 1 < 2;\n" || k2 != "path" || v2 != "a.js" {
		t.Fatalf("JsonReader mismatch: %q=%q, %q=%q", k1, v1, k2, v2)
	}
}

func TestGsonDOMAndCollections(t *testing.T) {
	g := NewGsonBuilder().DisableHtmlEscaping().SetPrettyPrinting().Create()
	m := javautil.NewLinkedHashMap[string, int32]()
	m.Put("alpha", 1)
	m.Put("beta", 2)

	jsonText := g.ToJson(m)
	obj := g.FromJsonObject(jsonText)
	if obj.Get("alpha").GetAsInt() != 1 || obj.Get("beta").GetAsInt() != 2 {
		t.Fatalf("unexpected DOM values from %s", jsonText)
	}

	list := g.FromJsonStringList(`["f1.js", "f2.js"]`)
	if list.Size() != 2 || list.Get(0) != "f1.js" || list.Get(1) != "f2.js" {
		t.Fatalf("FromJsonStringList failed: %v", list)
	}
}
