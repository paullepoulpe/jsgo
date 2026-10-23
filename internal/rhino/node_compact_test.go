package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"runtime"
	"testing"
	"unsafe"
)

func TestCompactNodeLayoutAndGCRoots(t *testing.T) {
	if unsafe.Sizeof(Node{}) != 64 {
		t.Fatalf("node occupies %d bytes, expected 64", unsafe.Sizeof(Node{}))
	}
	registry := NewJSTypeRegistry(nil)
	typ := registry.GetNativeType(JSTypeNative_NUMBER_TYPE)
	n := Node_NewStringWithToken(Token_NAME, "distinct_compact_node_name")
	n.SetJSType(typ)
	typ = nil
	for i := 0; i < 3; i++ {
		runtime.GC()
	}
	if n.GetJSType() != registry.GetNativeType(JSTypeNative_NUMBER_TYPE) || n.GetString() != "distinct_compact_node_name" {
		t.Fatal("lost type or string payload")
	}
	clone := n.CloneTree()
	clone.SetString("changed")
	if n.GetString() == clone.GetString() || clone.GetJSType() != n.GetJSType() {
		t.Fatal("clone shared mutable field or lost type")
	}
	n.SetJSType(nil)
	if n.GetJSType() != nil || n.GetColor() != nil {
		t.Fatal("cleared type is non-null")
	}
	// String pointer cells must remain valid even after the bounded pool clears.
	rhinoStringPoolClearForTest()
	runtime.GC()
	if n.GetString() != "distinct_compact_node_name" {
		t.Fatal("intern table owned node payload lifetime")
	}
	n.SetOriginalName(jutil.StringPtr("original"))
}

// The length (when it does not fit in 16 bits) and the original name of a
// NUMBER or TEMPLATELIT_STRING node live in the node's nodeExtBox.
func TestNodeExtBoxLengthAndOriginalName(t *testing.T) {
	n := Node_NewNumber(42)
	n.SetLength(100000)
	n.SetOriginalName(jutil.StringPtr("answer"))
	info := NewJSDocInfo_Builder().BuildWithAlways(true)
	n.SetJSDocInfo(info)
	for _, c := range []*Node{n, n.CloneNode(), NewNode(Token_NUMBER).Srcref(n)} {
		if c.GetLength() != 100000 || c.GetOriginalName() == nil || *c.GetOriginalName() != "answer" {
			t.Fatalf("lost length %d or original name %v", c.GetLength(), c.GetOriginalName())
		}
	}
	if n.GetJSDocInfo() != info {
		t.Fatal("lost JSDocInfo")
	}
	n.SetLength(-5)
	if n.GetLength() != -5 {
		t.Fatalf("length = %d, want -5", n.GetLength())
	}
	n.SetLength(7)
	n.SetOriginalName(nil)
	n.SetJSDocInfo(nil)
	if n.GetLength() != 7 || n.GetOriginalName() != nil || n.GetJSDocInfo() != nil {
		t.Fatal("cleared fields are still set")
	}
	if n.kindFlags&kindFlagHasExtBox != 0 || n.typeOrColor != nil {
		t.Fatal("empty nodeExtBox was not released")
	}
}
