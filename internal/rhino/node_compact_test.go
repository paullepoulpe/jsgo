package rhino

import (
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
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
	if n.kindFlags&kindFlagHasExtBox != 0 || n.aux != 0 {
		t.Fatal("empty nodeExtBox was not released")
	}
}

// The node's only pointers are its four tree links: types, colors and cold
// data are reached through 32-bit handles into the node tables.
func TestNodePointerWordsAreTheTreeLinks(t *testing.T) {
	typ := reflect.TypeOf(Node{})
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		isLink := f.Name == "next" || f.Name == "first" || f.Name == "parent" || f.Name == "previous"
		if (f.Type.Kind() == reflect.Pointer || f.Type.Kind() == reflect.UnsafePointer) != isLink {
			t.Errorf("field %s: pointer-ness does not match the 4-link layout", f.Name)
		}
	}
	if unsafe.Sizeof(nodeExtBox{}) != 32 {
		t.Errorf("nodeExtBox is %d bytes, want 32", unsafe.Sizeof(nodeExtBox{}))
	}
	if unsafe.Sizeof(JSTypeBase{}) != 32 {
		t.Errorf("JSTypeBase is %d bytes, want 32", unsafe.Sizeof(JSTypeBase{}))
	}
}

// Releasing the tables (what the Compiler does when it is collected, and
// ClearJSTypeRegistry does for types) makes every handle read as absent;
// types, colors and JSDoc set afterwards work again, and a Color or JSType
// registered under a released owner is re-registered on its next use.
func TestNodeTablesRelease(t *testing.T) {
	prev := currentNodeTables.Load()
	defer currentNodeTables.Store(prev)
	nodeTablesMu.Lock()
	older := nodeTablesLive
	nodeTablesMu.Unlock()
	for _, o := range older {
		o.Release() // this test's tables must be the oldest live ones
	}
	tables := NewNodeTables()
	tables.Install()

	registry := NewJSTypeRegistry(nil)
	num := registry.GetNativeType(JSTypeNative_NUMBER_TYPE)
	color := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("rhino_t")).Build()
	info := NewJSDocInfo_Builder().BuildWithAlways(true)

	typed := NewNode(Token_NAME).SetJSType(num)
	colored := NewNode(Token_NAME).SetColor(color)
	documented := NewNode(Token_VAR).SetJSDocInfo(info)
	documented.SetLength(1 << 20)
	if typed.GetJSType() != num || colored.GetColor() != color || documented.GetJSDocInfo() != info || documented.GetLength() != 1<<20 {
		t.Fatal("handles do not resolve")
	}

	registry.ReleaseNodeTypeHandles()
	if typed.GetJSType() != nil || typed.GetJSTypeBase() != nil {
		t.Fatal("type still resolves after the registry released its handles")
	}
	typed.SetJSType(nil)
	typed.SetJSType(num)
	if typed.GetJSType() != num {
		t.Fatal("type cannot be re-registered after release")
	}

	tables.Release()
	if colored.GetColor() != nil || documented.GetJSDocInfo() != nil || documented.GetLength() != 0 {
		t.Fatal("color or cold data still resolve after the tables were released")
	}
	if currentNodeTables.Load() == tables {
		t.Fatal("released tables are still current")
	}
	colored.SetColor(nil)
	colored.SetColor(color)
	documented.SetJSDocInfo(info)
	if colored.GetColor() != color || documented.GetJSDocInfo() != info {
		t.Fatal("color or JSDoc cannot be set again after release")
	}
}
