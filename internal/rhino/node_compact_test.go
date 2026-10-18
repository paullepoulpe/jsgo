package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"runtime"
	"testing"
	"unsafe"
)

func TestCompactNodeLayoutAndGCRoots(t *testing.T) {
	if unsafe.Sizeof(Node{}) != 80 {
		t.Fatalf("node occupies %d bytes, expected 80", unsafe.Sizeof(Node{}))
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
	rhinoStringPool_mu.Lock()
	clear(rhinoStringPool_strings)
	rhinoStringPool_mu.Unlock()
	runtime.GC()
	if n.GetString() != "distinct_compact_node_name" {
		t.Fatal("intern table owned node payload lifetime")
	}
	n.SetOriginalName(jutil.StringPtr("original"))
}
