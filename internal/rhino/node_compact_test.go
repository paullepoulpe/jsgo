package rhino

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestCompactNodeLayoutAndGCRoots(t *testing.T) {
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
	for _, c := range []Node{n, n.CloneNode(), NewNode(Token_NUMBER).Srcref(n)} {
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
	if n.rec().kindFlags&kindFlagHasExtBox != 0 || n.rec().aux != 0 {
		t.Fatal("empty nodeExtBox was not released")
	}
}

// The node record holds no pointer at all: the tree links are 32-bit node
// ids (node_store.go) and types, colors and cold data are reached
// through 32-bit handles into the store's tables, so the GC never traces the
// AST (a segment of records is a pointer-free object).
func TestNodeHasNoPointers(t *testing.T) {
	typ := reflect.TypeOf(nodeData{})
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Type.Kind() == reflect.Pointer || f.Type.Kind() == reflect.UnsafePointer {
			t.Errorf("field %s is a pointer", f.Name)
		}
	}
	if typ.Size() != 48 {
		t.Errorf("nodeData is %d bytes, want 48", typ.Size())
	}
	if s := reflect.TypeOf(nodeExtBox{}).Size(); s != 32 {
		t.Errorf("nodeExtBox is %d bytes, want 32", s)
	}
	if s := reflect.TypeOf(JSTypeBase{}).Size(); s != 32 {
		t.Errorf("JSTypeBase is %d bytes, want 32", s)
	}
	if s := reflect.TypeOf(NilNode).Size(); s != 4 {
		t.Errorf("Node is %d bytes, want 4", s)
	}
}

// ReleaseNodeTypeHandles (what ClearJSTypeRegistry calls) makes every type
// handle of the registry read as absent; a type set afterwards works again.
// Colors and cold data live as long as the store.
func TestNodeTypeRelease(t *testing.T) {
	prev := SetCurrentNodeStore(NewNodeStore())
	defer SetCurrentNodeStore(prev)

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
	if colored.GetColor() != color || documented.GetJSDocInfo() != info {
		t.Fatal("releasing types dropped a color or JSDoc")
	}
	typed.SetJSType(nil)
	typed.SetJSType(num)
	if typed.GetJSType() != num {
		t.Fatal("type cannot be re-registered after release")
	}
}

// A node of a store that is no longer current cannot be used, and a
// package-level node or template (the static store) cannot be linked.
func TestNodeStoreIsolation(t *testing.T) {
	prev := SetCurrentNodeStore(NewNodeStore())
	defer SetCurrentNodeStore(prev)
	a := IR_Block()
	SetCurrentNodeStore(NewNodeStore())
	expectPanic := func(what, msg string, f func()) {
		t.Helper()
		defer func() {
			t.Helper()
			r := recover()
			if s, _ := r.(string); !strings.Contains(s, msg) {
				t.Fatalf("%s: expected a panic containing %q, got %v", what, msg, r)
			}
		}()
		f()
	}
	expectPanic("stale node", "not in the current node store", func() { a.GetFirstChild() })
	tmpl := NewTemplateNode(Token_SCRIPT)
	if !tmpl.Store().IsStatic() || CurrentNodeStore().IsStatic() {
		t.Fatal("template is not in the static store")
	}
	expectPanic("linking a template", "must never be linked", func() { IR_Block().AddChildToBack(tmpl) })
	expectPanic("nil node", "nil node", func() { NilNode.GetToken() })
	if !NilNode.IsNil() || IR_Empty().IsNil() || IR_Empty() == IR_Empty() {
		t.Fatal("nil handle or handle equality is wrong")
	}
	n := IR_Name("x")
	if n.GetFirstChild() != NilNode || n.GetParent() != NilNode {
		t.Fatal("absent links are not the nil node")
	}
}
