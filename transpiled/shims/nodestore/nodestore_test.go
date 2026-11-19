package nodestore

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

func TestNodeDataMemoryLayout(t *testing.T) {
	sz := reflect.TypeFor[NodeData]().Size()
	if sz != 32 {
		t.Fatalf("expected NodeData size to be 32 bytes, got %d", sz)
	}

	typ := reflect.TypeFor[NodeData]()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		switch field.Type.Kind() {
		case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Interface, reflect.Func, reflect.Chan, reflect.String, reflect.UnsafePointer:
			t.Fatalf("NodeData contains Go pointer/reference type in field %s (%v); noscan violation", field.Name, field.Type)
		}
	}
}

func TestNodeHandleProperties(t *testing.T) {
	if !NilNode.IsNil() {
		t.Errorf("NilNode.IsNil() should be true")
	}
	if NilNode.ID() != 0 {
		t.Errorf("NilNode.ID() should be 0")
	}

	n := Node(42)
	if n.IsNil() {
		t.Errorf("Node(42).IsNil() should be false")
	}
	if n.ID() != 42 {
		t.Errorf("Node(42).ID() should be 42")
	}
	if !n.Equals(Node(42)) {
		t.Errorf("Node(42).Equals(Node(42)) should be true")
	}
	if n.Equals(Node(43)) {
		t.Errorf("Node(42).Equals(Node(43)) should be false")
	}
	var idVal uint32 = 42
	if n.HashCode() != int32(idVal*0x9E3779B1) {
		t.Errorf("unexpected HashCode for Node(42)")
	}
}

func TestArenaMultiBlockAllocation(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	// Allocate more than 2 blocks (> 32,768 nodes)
	const count = 35000
	nodes := make([]Node, count)
	for i := 0; i < count; i++ {
		nodes[i] = s.Alloc()
		nodes[i].SetToken(uint8(i % 250))
		nodes[i].SetLinenoCharno(i+1, (i*3)%100)
		nodes[i].SetLength(int32(i % 500))
	}

	for i := 0; i < count; i++ {
		n := nodes[i]
		if n.Token() != uint8(i%250) {
			t.Fatalf("node %d: token mismatch: got %d, want %d", i, n.Token(), i%250)
		}
		if n.Lineno() != i+1 {
			t.Fatalf("node %d: lineno mismatch: got %d, want %d", i, n.Lineno(), i+1)
		}
		if n.Charno() != (i*3)%100 {
			t.Fatalf("node %d: charno mismatch: got %d, want %d", i, n.Charno(), (i*3)%100)
		}
		if n.Length() != int32(i%500) {
			t.Fatalf("node %d: length mismatch: got %d, want %d", i, n.Length(), i%500)
		}
	}
}

func TestNodeHierarchyLinks(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	alloc := NewNodeAllocator(s)
	c1 := alloc.AllocWithToken(10)
	c2 := alloc.AllocWithToken(20)
	parent := alloc.Alloc2ChildrenCachedTemplate(100, c1, c2)

	if parent.Token() != 100 {
		t.Errorf("parent token: got %d, want 100", parent.Token())
	}
	if parent.FirstChild() != c1 {
		t.Errorf("parent.FirstChild(): got %d, want %d", parent.FirstChild(), c1)
	}
	if c1.Parent() != parent {
		t.Errorf("c1.Parent(): got %d, want %d", c1.Parent(), parent)
	}
	if c1.Next() != c2 {
		t.Errorf("c1.Next(): got %d, want %d", c1.Next(), c2)
	}
	if c2.Parent() != parent {
		t.Errorf("c2.Parent(): got %d, want %d", c2.Parent(), parent)
	}
	if c2.Previous() != c1 {
		t.Errorf("c2.Previous(): got %d, want %d", c2.Previous(), c1)
	}
}

func TestSideTables(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	n := s.Alloc()
	n.SetTC(12345)
	if n.TC() != 12345 {
		t.Errorf("n.TC(): got %d, want 12345", n.TC())
	}

	n.SetOrigNameID(6789)
	if n.OrigNameID() != 6789 {
		t.Errorf("n.OrigNameID(): got %d, want 6789", n.OrigNameID())
	}

	n.SetAux(NodeAux{JSDoc: "doc", Extra: 42})
	aux, ok := n.GetAux()
	if !ok || aux.JSDoc != "doc" || aux.Extra != 42 {
		t.Errorf("n.GetAux(): got %+v, ok=%v", aux, ok)
	}
}

type testPropItem struct {
	key  int
	val  string
	next *testPropItem
}

func TestPagedAnySideColumnsAndNilReset(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	if NilNode.TypeOrColor() != nil {
		t.Fatalf("expected NilNode.TypeOrColor() == nil, got %v", NilNode.TypeOrColor())
	}
	if NilNode.PropList() != nil {
		t.Fatalf("expected NilNode.PropList() == nil, got %v", NilNode.PropList())
	}

	// Setting nil/0 on a freshly allocated node in an unmapped column must be a no-op
	// without allocating NodeColumn / PagedColumn / propShards.
	n0 := s.Alloc()
	n0.SetTypeOrColor(nil)
	n0.SetPropList(nil)
	n0.SetOrigNameID(0)
	for i := range s.propShards {
		if s.propShards[i].m != nil {
			t.Fatalf("setting nil PropList allocated propShard[%d] map eagerly", i)
		}
	}
	if len(s.tcs) != 0 || len(s.origNames) != 0 {
		t.Fatalf("setting zero/nil values allocated columns eagerly: tc=%d orig=%d",
			len(s.tcs), len(s.origNames))
	}

	// Allocate across multiple 16,384-node blocks and 64-node pages.
	const total = 36000
	nodes := make([]Node, total)
	for i := 0; i < total; i++ {
		n := s.Alloc()
		nodes[i] = n
		if i%3 == 0 {
			n.SetTypeOrColor(i * 10)
		}
		if i%5 == 0 {
			n.SetPropList(&testPropItem{key: i, val: "prop"})
		}
		if i%7 == 0 {
			n.SetOrigNameID(StringID(i + 1))
		}
	}

	for i, n := range nodes {
		if i%3 == 0 {
			if got := n.TypeOrColor(); got != i*10 {
				t.Fatalf("node %d TypeOrColor: got %v, want %d", i, got, i*10)
			}
		} else if got := n.TypeOrColor(); got != nil {
			t.Fatalf("node %d TypeOrColor: got %v, want nil", i, got)
		}

		if i%5 == 0 {
			p, ok := n.PropList().(*testPropItem)
			if !ok || p == nil || p.key != i || p.val != "prop" {
				t.Fatalf("node %d PropList: got %#v", i, n.PropList())
			}
		} else if got := n.PropList(); got != nil {
			t.Fatalf("node %d PropList: got %v, want nil", i, got)
		}

		if i%7 == 0 {
			if got := n.OrigNameID(); got != StringID(i+1) {
				t.Fatalf("node %d OrigNameID: got %d, want %d", i, got, i+1)
			}
		} else if got := n.OrigNameID(); got != 0 {
			t.Fatalf("node %d OrigNameID: got %d, want 0", i, got)
		}
	}

	// Overwrite and reset back to nil / 0, verifying values are cleared cleanly.
	for i := 0; i < total; i += 15 {
		n := nodes[i]
		n.SetTypeOrColor("updated")
		if got := n.TypeOrColor(); got != "updated" {
			t.Fatalf("node %d updated TypeOrColor: got %v, want updated", i, got)
		}
		n.SetTypeOrColor(nil)
		if got := n.TypeOrColor(); got != nil {
			t.Fatalf("node %d nil-reset TypeOrColor: got %v, want nil", i, got)
		}

		n.SetPropList(nil)
		if got := n.PropList(); got != nil {
			t.Fatalf("node %d nil-reset PropList: got %v, want nil", i, got)
		}

		n.SetOrigNameID(0)
		if got := n.OrigNameID(); got != 0 {
			t.Fatalf("node %d zero-reset OrigNameID: got %d, want 0", i, got)
		}
	}
}

func TestNodeStoreSwitchingSideColumns(t *testing.T) {
	s1 := NewNodeStore()
	s2 := NewNodeStore()

	tmpl := NewTemplateNode(77)
	tmpl.SetPropList("static-prop")
	defer tmpl.SetPropList(nil)

	prev := SetCurrentNodeStore(s1)
	defer SetCurrentNodeStore(prev)

	a1 := s1.Alloc()
	a1.SetTypeOrColor("s1-type")
	a1.SetPropList(&testPropItem{key: 1, val: "s1-prop"})
	a1.SetOrigNameID(101)

	// Switch to s2: s1 must unregister, s2 starts clean, static template remains accessible.
	if old := SetCurrentNodeStore(s2); old != s1 {
		t.Fatalf("SetCurrentNodeStore(s2) returned %v, want s1", old)
	}
	if a1.TypeOrColor() != nil || a1.PropList() != nil || a1.OrigNameID() != 0 {
		t.Fatalf("unregistered s1 node still visible while s2 active")
	}
	if tmpl.PropList() != "static-prop" {
		t.Fatalf("static template PropList lost across store switch: %v", tmpl.PropList())
	}

	a2 := s2.Alloc()
	a2.SetTypeOrColor("s2-type")
	a2.SetPropList(&testPropItem{key: 2, val: "s2-prop"})
	a2.SetOrigNameID(202)

	if a2.TypeOrColor() != "s2-type" || a2.OrigNameID() != 202 {
		t.Fatalf("s2 side columns mismatch: tc=%v orig=%d", a2.TypeOrColor(), a2.OrigNameID())
	}

	// Switch back to s1: s2 must unregister and s1 columns must be restored intact.
	if old := SetCurrentNodeStore(s1); old != s2 {
		t.Fatalf("SetCurrentNodeStore(s1) returned %v, want s2", old)
	}
	if a2.TypeOrColor() != nil || a2.PropList() != nil || a2.OrigNameID() != 0 {
		t.Fatalf("unregistered s2 node still visible after switching back to s1")
	}
	if a1.TypeOrColor() != "s1-type" {
		t.Fatalf("restored s1 TypeOrColor: got %v, want s1-type", a1.TypeOrColor())
	}
	if p, ok := a1.PropList().(*testPropItem); !ok || p.val != "s1-prop" {
		t.Fatalf("restored s1 PropList: got %#v", a1.PropList())
	}
	if a1.OrigNameID() != 101 {
		t.Fatalf("restored s1 OrigNameID: got %d, want 101", a1.OrigNameID())
	}
}

func TestConcurrentSideColumnsAcrossBlocksAndPages(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	// 1. Pre-allocate consecutive nodes so adjacent handles in the SAME 64-node page
	// are written/read concurrently by distinct goroutines (stressing CAS page allocation).
	const sharedCount = 4096
	shared := make([]Node, sharedCount)
	for i := 0; i < sharedCount; i++ {
		shared[i] = s.Alloc()
	}

	const goroutines = 8
	const perGoroutine = 5000 // 8 * 5000 = 40,000 > 2 * 16,384 blocks
	var wg sync.WaitGroup
	wg.Add(goroutines)

	batchNodes := make([][]Node, goroutines)
	for g := 0; g < goroutines; g++ {
		g := g
		batchNodes[g] = make([]Node, perGoroutine)
		go func() {
			defer wg.Done()
			// Concurrent writes/reads on interleaved slots within shared 64-node pages.
			for i := g; i < sharedCount; i += goroutines {
				n := shared[i]
				n.SetTypeOrColor(i)
				n.SetPropList(&testPropItem{key: i, val: "shared"})
				n.SetOrigNameID(StringID(i + 10))
				if got := n.TypeOrColor(); got != i {
					t.Errorf("shared node %d TypeOrColor: got %v, want %d", i, got, i)
				}
			}

			// Concurrent NodeAllocator allocations + side-column reads/writes across blocks.
			alloc := NewNodeAllocator(s)
			for i := 0; i < perGoroutine; i++ {
				n := alloc.AllocWithToken(uint8(g + 1))
				n.SetTypeOrColor(g*100000 + i)
				n.SetPropList(&testPropItem{key: i, val: "batch"})
				n.SetOrigNameID(StringID(g*100000 + i + 1))
				if i%4 == 0 {
					n.SetTypeOrColor(nil)
				}
				batchNodes[g][i] = n
			}
		}()
	}
	wg.Wait()

	for i, n := range shared {
		if got := n.TypeOrColor(); got != i {
			t.Fatalf("shared[%d] TypeOrColor: got %v, want %d", i, got, i)
		}
		if p, ok := n.PropList().(*testPropItem); !ok || p.key != i {
			t.Fatalf("shared[%d] PropList: got %#v", i, n.PropList())
		}
		if got := n.OrigNameID(); got != StringID(i+10) {
			t.Fatalf("shared[%d] OrigNameID: got %d, want %d", i, got, i+10)
		}
	}

	for g := 0; g < goroutines; g++ {
		for i, n := range batchNodes[g] {
			wantTC := any(g*100000 + i)
			if i%4 == 0 {
				wantTC = nil
			}
			if got := n.TypeOrColor(); got != wantTC {
				t.Fatalf("batch[%d][%d] TypeOrColor: got %v, want %v", g, i, got, wantTC)
			}
			if p, ok := n.PropList().(*testPropItem); !ok || p.key != i || p.val != "batch" {
				t.Fatalf("batch[%d][%d] PropList: got %#v", g, i, n.PropList())
			}
			if got := n.OrigNameID(); got != StringID(g*100000+i+1) {
				t.Fatalf("batch[%d][%d] OrigNameID: got %d, want %d", g, i, got, g*100000+i+1)
			}
		}
	}
}

func TestConcurrentBatchAllocation(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	const goroutines = 8
	const perGoroutine = 5000
	var wg sync.WaitGroup
	wg.Add(goroutines)

	allNodes := make([][]Node, goroutines)
	for g := 0; g < goroutines; g++ {
		g := g
		allNodes[g] = make([]Node, perGoroutine)
		go func() {
			defer wg.Done()
			alloc := NewNodeAllocator(s)
			for i := 0; i < perGoroutine; i++ {
				n := alloc.AllocWithToken(uint8(g + 1))
				n.SetLinenoCharno(i+1, g)
				allNodes[g][i] = n
			}
		}()
	}
	wg.Wait()

	seen := make(map[Node]bool, goroutines*perGoroutine)
	for g := 0; g < goroutines; g++ {
		for i := 0; i < perGoroutine; i++ {
			n := allNodes[g][i]
			if seen[n] {
				t.Fatalf("duplicate node handle detected: %d", n)
			}
			seen[n] = true
			if n.Token() != uint8(g+1) {
				t.Fatalf("token mismatch for node %d: got %d, want %d", n, n.Token(), g+1)
			}
			if n.Lineno() != i+1 {
				t.Fatalf("lineno mismatch for node %d: got %d, want %d", n, n.Lineno(), i+1)
			}
			if n.Charno() != g {
				t.Fatalf("charno mismatch for node %d: got %d, want %d", n, n.Charno(), g)
			}
		}
	}
}

type fakeTCHolder struct {
	name     string
	tcHandle atomic.Uint32
}

func (f *fakeTCHolder) IsJavaNil_() bool            { return f == nil }
func (f *fakeTCHolder) TCHandlePtr_() *atomic.Uint32 { return &f.tcHandle }

type fakeJavaNilOnly struct {
	label string
}

func (f *fakeJavaNilOnly) IsJavaNil_() bool { return f == nil }

func hasGoPointers(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Interface, reflect.Func, reflect.Chan, reflect.String, reflect.UnsafePointer:
		return true
	case reflect.Array:
		return hasGoPointers(typ.Elem())
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			if hasGoPointers(typ.Field(i).Type) {
				return true
			}
		}
	}
	return false
}

func TestNodeColumnNoscanAndTCHandleHolder(t *testing.T) {
	if sz := reflect.TypeFor[NodeColumn]().Size(); sz != blockLen*4 {
		t.Fatalf("expected NodeColumn size %d, got %d", blockLen*4, sz)
	}
	if hasGoPointers(reflect.TypeFor[NodeColumn]()) {
		t.Fatalf("NodeColumn contains Go pointers; noscan violation")
	}

	s1 := NewNodeStore()
	s2 := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer SetCurrentNodeStore(prev)

	// Typed-nil TCHandleHolder and typed-nil javaNil must NOT allocate columns or deref nil.
	var typedNilHolder *fakeTCHolder
	var typedNilPlain *fakeJavaNilOnly
	nNil := s1.Alloc()
	nNil.SetTypeOrColor(any(typedNilHolder))
	nNil.SetTypeOrColor(any(typedNilPlain))
	if len(s1.tcs) != 0 {
		t.Fatalf("typed-nil values should not allocate TC column, got %d", len(s1.tcs))
	}
	if got := nNil.TypeOrColor(); got != nil {
		t.Fatalf("expected nil after typed-nil SetTypeOrColor, got %#v", got)
	}

	singleton := &fakeTCHolder{name: "StandardColors.UNKNOWN"}
	n1 := s1.Alloc()
	n2 := s1.Alloc()
	n1.SetTypeOrColor(singleton)
	n2.SetTypeOrColor(singleton)

	h1 := n1.TC()
	h2 := n2.TC()
	if h1 == 0 || h1 != h2 {
		t.Fatalf("expected identical non-zero cached TC handle on n1 and n2, got %d vs %d", h1, h2)
	}
	if got := singleton.tcHandle.Load(); got != h1 {
		t.Fatalf("holder cached handle = %d, want %d", got, h1)
	}
	if n1.TypeOrColor() != singleton || n2.TypeOrColor() != singleton {
		t.Fatalf("TypeOrColor mismatch for singleton")
	}

	// Resetting with typed-nil clears TC back to 0.
	n2.SetTypeOrColor(any(typedNilHolder))
	if n2.TC() != 0 || n2.TypeOrColor() != nil {
		t.Fatalf("typed-nil reset failed: TC=%d val=%#v", n2.TC(), n2.TypeOrColor())
	}

	// Global singleton handle works across NodeStore switches and static template nodes.
	tmpl := NewTemplateNode(55)
	tmpl.SetTypeOrColor(singleton)
	defer tmpl.SetTypeOrColor(nil)

	SetCurrentNodeStore(s2)
	m1 := s2.Alloc()
	m1.SetTypeOrColor(singleton)
	if m1.TC() != h1 || m1.TypeOrColor() != singleton {
		t.Fatalf("singleton handle was not reused in s2: TC=%d, want %d", m1.TC(), h1)
	}
	if tmpl.TC() != h1 || tmpl.TypeOrColor() != singleton {
		t.Fatalf("static template singleton TypeOrColor mismatch")
	}
}

func TestConcurrentTCHandleHolderRace(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	// Shared fresh holders whose tcHandle starts at 0, contended by all goroutines simultaneously.
	const numHolders = 64
	holders := make([]*fakeTCHolder, numHolders)
	for i := range holders {
		holders[i] = &fakeTCHolder{name: "color"}
	}
	var typedNil *fakeTCHolder

	const goroutines = 8
	const perGoroutine = 2000
	nodes := make([][]Node, goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		g := g
		nodes[g] = make([]Node, perGoroutine)
		go func() {
			defer wg.Done()
			alloc := NewNodeAllocator(s)
			for i := 0; i < perGoroutine; i++ {
				n := alloc.Alloc()
				h := holders[(i+g)%numHolders]
				n.SetTypeOrColor(h)
				if got := n.TypeOrColor(); got != h {
					t.Errorf("concurrent TypeOrColor mismatch: got %#v, want %#v", got, h)
				}
				if i%9 == 0 {
					n.SetTypeOrColor(any(typedNil))
				}
				nodes[g][i] = n
			}
		}()
	}
	wg.Wait()

	for idx, h := range holders {
		if h.tcHandle.Load() == 0 {
			t.Fatalf("holder %d handle was not cached", idx)
		}
	}

	for g := 0; g < goroutines; g++ {
		for i, n := range nodes[g] {
			if i%9 == 0 {
				if n.TypeOrColor() != nil || n.TC() != 0 {
					t.Fatalf("nodes[%d][%d] expected nil after typed-nil reset, got %v", g, i, n.TypeOrColor())
				}
				continue
			}
			want := holders[(i+g)%numHolders]
			if got := n.TypeOrColor(); got != want {
				t.Fatalf("nodes[%d][%d] got %#v, want %#v", g, i, got, want)
			}
			if n.TC() != want.tcHandle.Load() {
				t.Fatalf("nodes[%d][%d] TC handle %d != holder cached handle %d", g, i, n.TC(), want.tcHandle.Load())
			}
		}
	}
}

type fakeSourceFileHolder struct {
	name     string
	sfHandle atomic.Uint32
}

func (f *fakeSourceFileHolder) IsJavaNil_() bool                    { return f == nil }
func (f *fakeSourceFileHolder) SourceFileHandlePtr_() *atomic.Uint32 { return &f.sfHandle }

type fakeColdProp struct {
	propType byte
	val      string
	next     *fakeColdProp
}

func (f *fakeColdProp) IsJavaNil_() bool { return f == nil }

func TestColumnPageNumBlockNoscanAndSourceFileHandleHolder(t *testing.T) {
	if !hasGoPointers(reflect.TypeFor[*int]()) {
		t.Fatalf("hasGoPointers[*int] sanity check failed")
	}
	if sz := reflect.TypeFor[columnPage]().Size(); sz != columnPageLen*4 {
		t.Fatalf("expected columnPage size %d, got %d", columnPageLen*4, sz)
	}
	if hasGoPointers(reflect.TypeFor[columnPage]()) {
		t.Fatalf("columnPage contains Go pointers; noscan violation")
	}
	if hasGoPointers(reflect.TypeFor[NodeColumn]()) {
		t.Fatalf("NodeColumn contains Go pointers; noscan violation")
	}
	if hasGoPointers(reflect.TypeFor[numBlock]()) {
		t.Fatalf("numBlock contains Go pointers; noscan violation")
	}

	if NilNode.SourceFileObj() != nil {
		t.Fatalf("expected NilNode.SourceFileObj() == nil, got %#v", NilNode.SourceFileObj())
	}

	s1 := NewNodeStore()
	s2 := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer SetCurrentNodeStore(prev)

	var typedNilSF *fakeSourceFileHolder
	var typedNilPlain *fakeJavaNilOnly
	var typedNilProp *fakeColdProp
	nNil := AllocNodeInit(3)
	if nNil.Token() != 0xFF || nNil.Lineno() != -1 || nNil.Charno() != -1 || nNil.KindFlags() != 3 {
		t.Fatalf("AllocNodeInit fields mismatch: tok=%d line=%d char=%d kind=%d",
			nNil.Token(), nNil.Lineno(), nNil.Charno(), nNil.KindFlags())
	}
	nNil.SetSourceFileObj(any(typedNilSF))
	nNil.SetSourceFileObj(any(typedNilPlain))
	if nNil.SourceFile() != 0 || nNil.SourceFileObj() != nil {
		t.Fatalf("typed-nil SetSourceFileObj should leave SourceFile=0, got %d / %#v",
			nNil.SourceFile(), nNil.SourceFileObj())
	}
	nNil.SetPropList(any(typedNilProp))
	nNil.SetPropList(any(typedNilPlain))
	if got := nNil.PropList(); got != nil {
		t.Fatalf("expected nil after typed-nil SetPropList, got %#v", got)
	}

	sf := &fakeSourceFileHolder{name: "input.js"}
	n1 := AllocNodeInit(1)
	n2 := AllocNodeInit(2)
	n1.SetSourceFileObj(sf)
	n2.SetSourceFileObj(sf)

	h1 := n1.SourceFile()
	h2 := n2.SourceFile()
	if h1 == 0 || h1 != h2 {
		t.Fatalf("expected identical non-zero cached SourceFile handle on n1 and n2, got %d vs %d", h1, h2)
	}
	if got := sf.sfHandle.Load(); got != h1 {
		t.Fatalf("holder cached sfHandle = %d, want %d", got, h1)
	}
	if n1.SourceFileObj() != sf || n2.SourceFileObj() != sf {
		t.Fatalf("SourceFileObj mismatch for sf")
	}

	// Non-holder comparable fallback reuses handle too.
	n1.SetSourceFileObj("plain-source.js")
	n2.SetSourceFileObj("plain-source.js")
	if n1.SourceFile() == 0 || n1.SourceFile() != n2.SourceFile() || n1.SourceFileObj() != "plain-source.js" {
		t.Fatalf("non-holder fallback SourceFileObj mismatch: h1=%d h2=%d obj=%v",
			n1.SourceFile(), n2.SourceFile(), n1.SourceFileObj())
	}

	// Resetting with typed-nil clears SourceFile back to 0.
	n2.SetSourceFileObj(any(typedNilSF))
	if n2.SourceFile() != 0 || n2.SourceFileObj() != nil {
		t.Fatalf("typed-nil reset failed: handle=%d val=%#v", n2.SourceFile(), n2.SourceFileObj())
	}

	// Global SourceFileHandleHolder works across NodeStore switches and static template nodes.
	tmpl := NewTemplateNode(55)
	tmpl.SetSourceFileObj(sf)
	defer tmpl.SetSourceFileObj(nil)

	SetCurrentNodeStore(s2)
	m1 := AllocNodeInit(5)
	m1.SetSourceFileObj(sf)
	if m1.SourceFile() != h1 || m1.SourceFileObj() != sf {
		t.Fatalf("sf handle was not reused in s2: handle=%d, want %d", m1.SourceFile(), h1)
	}
	if tmpl.SourceFile() != h1 || tmpl.SourceFileObj() != sf {
		t.Fatalf("static template sf SourceFileObj mismatch")
	}

	// Verify fallback to staticStore when CurrentNodeStore is nil.
	SetCurrentNodeStore(nil)
	sn := AllocNodeInit(7)
	if sn.ID()&staticBit == 0 || sn.Token() != 0xFF || sn.Lineno() != -1 || sn.KindFlags() != 7 {
		t.Fatalf("AllocNodeInit on nil current store did not initialize static node properly: id=%x", sn.ID())
	}
}

func TestNumberValPropFlagsAndKindFlagsIsolation(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	// +0.0 fast path must NOT allocate any numTable entry (Payload == 0).
	nZero := AllocNodeInit(1)
	nZero.SetNumberVal(0.0)
	if nZero.Payload() != 0 || s.nums.count.Load() != 0 {
		t.Fatalf("+0.0 allocated numTable entry: payload=%d count=%d", nZero.Payload(), s.nums.count.Load())
	}
	if got := nZero.NumberVal(); got != 0.0 || math.Signbit(got) {
		t.Fatalf("+0.0 roundtrip failed: got %v signbit=%v", got, math.Signbit(got))
	}

	// -0.0 must preserve negative sign bit (bits != 0).
	nNegZero := AllocNodeInit(1)
	nNegZero.SetNumberVal(math.Copysign(0, -1))
	if nNegZero.Payload() == 0 {
		t.Fatalf("-0.0 should allocate non-zero handle")
	}
	if got := nNegZero.NumberVal(); got != 0.0 || !math.Signbit(got) {
		t.Fatalf("-0.0 lost sign bit: got %v signbit=%v", got, math.Signbit(got))
	}

	// Special IEEE-754 values + coexistence with PropFlags (int32) and KindFlags + PropList!
	specials := []float64{
		0.0,
		math.Copysign(0, -1),
		1.0,
		-1.0,
		3.141592653589793,
		math.MaxFloat64,
		math.SmallestNonzeroFloat64,
		math.Inf(1),
		math.Inf(-1),
		math.NaN(),
		math.Float64frombits(0x7ff8000000123456),
	}
	propFlagVals := []int32{0, 1, -1, math.MinInt32, math.MaxInt32, -0x12345678}
	for idx, want := range specials {
		n := AllocNodeInit(1)
		pf := propFlagVals[idx%len(propFlagVals)]
		n.SetPropFlags(pf)
		n.SetNumberVal(want)
		n.SetPropList(&fakeColdProp{propType: byte(idx), val: "jsdoc"})
		// Mutating KindFlags while cold prop is set must preserve both KindFlags and PropList!
		n.SetKindFlags(4)
		if n.KindFlags() != 4 {
			t.Fatalf("KindFlags = %d, want 4", n.KindFlags())
		}
		if n.PropFlags() != pf {
			t.Fatalf("NumberVal clobbered PropFlags: got %d, want %d", n.PropFlags(), pf)
		}
		got := n.NumberVal()
		if math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("special[%d] bits mismatch: got 0x%x, want 0x%x",
				idx, math.Float64bits(got), math.Float64bits(want))
		}
		if cp, ok := n.PropList().(*fakeColdProp); !ok || cp.val != "jsdoc" {
			t.Fatalf("PropList lost after SetKindFlags: got %#v", n.PropList())
		}
		n.SetPropList(nil)
		if n.KindFlags() != 4 || n.PropList() != nil {
			t.Fatalf("after nil SetPropList: KindFlags=%d PropList=%#v", n.KindFlags(), n.PropList())
		}
	}

	// Cross numBlock boundary (> 4096 non-zero numbers) + overwrite in place vs new value.
	const numNodes = 5000
	nodes := make([]Node, numNodes)
	for i := 0; i < numNodes; i++ {
		n := AllocNodeInit(1)
		n.SetNumberVal(float64(i) + 0.25)
		nodes[i] = n
	}
	for i, n := range nodes {
		if got := n.NumberVal(); got != float64(i)+0.25 {
			t.Fatalf("node %d NumberVal = %v, want %v", i, got, float64(i)+0.25)
		}
	}
	// Overwrite with new non-zero value and then reset to +0.0.
	nodes[10].SetNumberVal(999.5)
	if got := nodes[10].NumberVal(); got != 999.5 {
		t.Fatalf("overwritten NumberVal = %v, want 999.5", got)
	}
	nodes[10].SetNumberVal(0.0)
	if got := nodes[10].NumberVal(); got != 0.0 || nodes[10].Payload() != 0 {
		t.Fatalf("reset +0.0 NumberVal = %v payload=%d", got, nodes[10].Payload())
	}

	// Static template node NumberVal works across store resets.
	tmpl := NewTemplateNode(39)
	tmpl.SetNumberVal(42.5)
	defer tmpl.SetNumberVal(0)
	SetCurrentNodeStore(nil)
	if got := tmpl.NumberVal(); got != 42.5 {
		t.Fatalf("static template NumberVal = %v, want 42.5", got)
	}
}

func TestPropListLeakFreeGCReclamationAndZeroAllocNilRead(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	nEmpty := AllocNodeInit(0)
	allocs := testing.AllocsPerRun(200, func() {
		if nEmpty.PropList() != nil {
			t.Fatal("expected nil PropList")
		}
	})
	if allocs > 0 {
		t.Fatalf("PropList() on node without cold props allocated %f objects, want 0", allocs)
	}

	var finalized atomic.Int32
	n := AllocNodeInit(2)
	func() {
		first := &fakeColdProp{propType: 1, val: "ephemeral-1"}
		second := &fakeColdProp{propType: 2, val: "ephemeral-2"}
		runtime.SetFinalizer(first, func(*fakeColdProp) { finalized.Add(1) })
		runtime.SetFinalizer(second, func(*fakeColdProp) { finalized.Add(1) })
		n.SetPropList(first)
		// Overwrite first -> first must become unreachable immediately.
		n.SetPropList(second)
		// Clear second -> second must become unreachable immediately.
		n.SetPropList(nil)
	}()

	for range 6 {
		runtime.GC()
		if finalized.Load() == 2 {
			break
		}
	}
	if got := finalized.Load(); got != 2 {
		t.Fatalf("expected both overwritten/cleared PropList items to be GC-reclaimed (0 leak), finalized=%d/2", got)
	}
}

func TestConcurrentSourceFileNumberValAndShardedPropListRace(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	const numSF = 64
	sfs := make([]*fakeSourceFileHolder, numSF)
	for i := range sfs {
		sfs[i] = &fakeSourceFileHolder{name: "file.js"}
	}
	var typedNilSF *fakeSourceFileHolder
	var typedNilProp *fakeColdProp

	const goroutines = 8
	const perGoroutine = 5000 // 40,000 > 2 blocks & > 4096 numBlock entries
	nodes := make([][]Node, goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		g := g
		nodes[g] = make([]Node, perGoroutine)
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				n := AllocNodeInit(uint8(g) & kindFlagsMask)
				sf := sfs[(i+g)%numSF]
				n.SetSourceFileObj(sf)
				n.SetPropFlags(int32((g+1)*100000 + i))
				n.SetNumberVal(float64(g*10000+i) + 0.5)
				n.SetOrigNameID(StringID(i + 1))
				if i%16 == 0 {
					p := &fakeColdProp{propType: byte(g), val: "cold"}
					n.SetPropList(p)
					if got := n.PropList(); got != p {
						t.Errorf("concurrent PropList mismatch: got %#v, want %#v", got, p)
					}
					if i%32 == 0 {
						n.SetPropList(any(typedNilProp))
					}
				}
				if i%11 == 0 {
					n.SetSourceFileObj(any(typedNilSF))
				}
				nodes[g][i] = n
			}
		}()
	}
	wg.Wait()

	for idx, sf := range sfs {
		if sf.sfHandle.Load() == 0 {
			t.Fatalf("sfHolder %d handle was not cached", idx)
		}
	}

	for g := 0; g < goroutines; g++ {
		for i, n := range nodes[g] {
			if n.Token() != 0xFF || n.Lineno() != -1 || n.KindFlags() != uint8(g)&kindFlagsMask {
				t.Fatalf("nodes[%d][%d] AllocNodeInit corrupted: tok=%d line=%d kind=%d",
					g, i, n.Token(), n.Lineno(), n.KindFlags())
			}
			if n.PropFlags() != int32((g+1)*100000+i) {
				t.Fatalf("nodes[%d][%d] PropFlags=%d", g, i, n.PropFlags())
			}
			if n.NumberVal() != float64(g*10000+i)+0.5 {
				t.Fatalf("nodes[%d][%d] NumberVal=%v", g, i, n.NumberVal())
			}
			if n.OrigNameID() != StringID(i+1) {
				t.Fatalf("nodes[%d][%d] OrigNameID=%d, want %d", g, i, n.OrigNameID(), i+1)
			}
			if i%16 == 0 && i%32 != 0 {
				if p, ok := n.PropList().(*fakeColdProp); !ok || p == nil || p.propType != byte(g) {
					t.Fatalf("nodes[%d][%d] cold PropList=%#v", g, i, n.PropList())
				}
			} else if n.PropList() != nil {
				t.Fatalf("nodes[%d][%d] expected nil PropList, got %#v", g, i, n.PropList())
			}
			if i%11 == 0 {
				if n.SourceFileObj() != nil || n.SourceFile() != 0 {
					t.Fatalf("nodes[%d][%d] expected nil SourceFileObj, got %v", g, i, n.SourceFileObj())
				}
				continue
			}
			wantSF := sfs[(i+g)%numSF]
			if got := n.SourceFileObj(); got != wantSF {
				t.Fatalf("nodes[%d][%d] SourceFileObj got %#v, want %#v", g, i, got, wantSF)
			}
			if n.SourceFile() != wantSF.sfHandle.Load() {
				t.Fatalf("nodes[%d][%d] SourceFile handle %d != cached %d",
					g, i, n.SourceFile(), wantSF.sfHandle.Load())
			}
		}
	}
}

type fakeJSDocHolder struct {
	comment     string
	jsdocHandle atomic.Uint32
}

func (f *fakeJSDocHolder) IsJavaNil_() bool               { return f == nil }
func (f *fakeJSDocHolder) JSDocHandlePtr_() *atomic.Uint32 { return &f.jsdocHandle }

func TestAnyTableZeroBoxAllocAndLargeHandles(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	// Unmapped / out-of-range handles in TC must safely return nil without panicking.
	nBogus := AllocNodeInit(0)
	nBogus.SetTC(60000000)
	if got := nBogus.TypeOrColor(); got != nil {
		t.Fatalf("expected nil for unmapped large TC handle, got %#v", got)
	}
	nBogus.SetTC(0xFFFFFFFF)
	if got := nBogus.TypeOrColor(); got != nil {
		t.Fatalf("expected nil for out-of-bounds TC handle, got %#v", got)
	}

	// Populate > 5000 distinct TCHandleHolders to cross 4096-entry anyBlock boundaries,
	// and verify 0 *any heap box allocations once an anyBlock is already allocated.
	const count = 5200
	holders := make([]*fakeTCHolder, count)
	nodes := make([]Node, count)
	for i := 0; i < count; i++ {
		holders[i] = &fakeTCHolder{name: "t"}
		nodes[i] = AllocNodeInit(0)
	}
	// Warm up column and first block, then check zero-alloc insertion for a fresh holder
	// within an already allocated anyBlock.
	nodes[0].SetTypeOrColor(holders[0])
	curCount := tcHandles.count.Load()
	if (curCount+1)&tableBlockMask != 0 {
		fresh := holders[1]
		n1 := nodes[1]
		allocs := testing.AllocsPerRun(1, func() {
			n1.SetTypeOrColor(fresh)
		})
		if allocs > 0 {
			t.Fatalf("SetTypeOrColor with new holder in warm anyBlock allocated %v objects, want 0 (no *any box)", allocs)
		}
	}

	for i := 2; i < count; i++ {
		nodes[i].SetTypeOrColor(holders[i])
	}
	for i := 0; i < count; i++ {
		if got := nodes[i].TypeOrColor(); got != holders[i] {
			t.Fatalf("node %d TypeOrColor got %#v, want %#v", i, got, holders[i])
		}
		if tcHandles.Peek(nodes[i].TC()) != holders[i] {
			t.Fatalf("AnyTable.Peek(%d) mismatch", nodes[i].TC())
		}
	}
	if tcHandles.Peek(0) != nil || tcHandles.Peek(0x7FFFFFFF) != nil {
		t.Fatalf("AnyTable.Peek on 0 / out-of-range should be nil")
	}
}

func TestJSDocColumnAndKindFlagsIsolation(t *testing.T) {
	s1 := NewNodeStore()
	s2 := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer SetCurrentNodeStore(prev)

	if NilNode.JSDocObj() != nil {
		t.Fatalf("expected NilNode.JSDocObj() == nil, got %#v", NilNode.JSDocObj())
	}

	var typedNilDoc *fakeJSDocHolder
	var typedNilPlain *fakeJavaNilOnly
	n0 := AllocNodeInit(0x07) // max 3-bit kind flag (7)
	if n0.KindFlags() != 0x07 {
		t.Fatalf("KindFlags = %d, want 7", n0.KindFlags())
	}
	n0.SetJSDocObj(nil)
	n0.SetJSDocObj(any(typedNilDoc))
	n0.SetJSDocObj(any(typedNilPlain))
	if s1.auxLen() != 0 {
		t.Fatalf("nil/typed-nil SetJSDocObj allocated aux entry eagerly: %d", s1.auxLen())
	}
	allocs := testing.AllocsPerRun(100, func() {
		if n0.JSDocObj() != nil {
			t.Fatal("expected nil JSDocObj")
		}
	})
	if allocs > 0 {
		t.Fatalf("JSDocObj() on node without JSDoc allocated %f objects, want 0", allocs)
	}

	doc1 := &fakeJSDocHolder{comment: "/** @const */"}
	doc2 := &fakeJSDocHolder{comment: "/** @type {number} */"}
	prop := &fakeColdProp{propType: 7, val: "cold"}

	n1 := AllocNodeInit(4)
	n2 := AllocNodeInit(2)
	n1.SetJSDocObj(doc1)
	n2.SetJSDocObj(doc1)
	if doc1.jsdocHandle.Load() != 0 {
		t.Fatalf("expected JSDocObj stored directly in NodeStore.aux without handle indirection, got handle=%d", doc1.jsdocHandle.Load())
	}
	if n1.JSDocObj() != doc1 || n2.JSDocObj() != doc1 {
		t.Fatalf("JSDocObj mismatch for doc1")
	}

	// Verify KindFlags (0..7) and internal flags never clobber one another.
	n1.SetPropList(prop)
	n1.SetKindFlags(7)
	if n1.KindFlags() != 7 || n1.JSDocObj() != doc1 || n1.PropList() != prop {
		t.Fatalf("SetKindFlags(7) corrupted JSDoc or PropList: kind=%d doc=%#v prop=%#v",
			n1.KindFlags(), n1.JSDocObj(), n1.PropList())
	}
	n1.SetJSDocObj(doc2)
	if n1.KindFlags() != 7 || n1.JSDocObj() != doc2 || n1.PropList() != prop {
		t.Fatalf("SetJSDocObj(doc2) corrupted KindFlags or PropList")
	}
	n1.SetJSDocObj(any(typedNilDoc))
	if n1.KindFlags() != 7 || n1.JSDocObj() != nil || n1.PropList() != prop {
		t.Fatalf("clearing JSDocObj corrupted KindFlags or PropList")
	}
	n1.SetJSDocObj(doc1)
	n1.SetPropList(nil)
	if n1.KindFlags() != 7 || n1.JSDocObj() != doc1 || n1.PropList() != nil {
		t.Fatalf("clearing PropList corrupted KindFlags or JSDocObj")
	}

	// Store switching + static template node support.
	tmpl := NewTemplateNode(50)
	tmpl.SetJSDocObj(doc2)
	defer tmpl.SetJSDocObj(nil)

	SetCurrentNodeStore(s2)
	if n1.JSDocObj() != nil {
		t.Fatalf("unregistered s1 JSDocObj still visible while s2 active")
	}
	if tmpl.JSDocObj() != doc2 {
		t.Fatalf("static template JSDocObj lost on store switch")
	}
	m := AllocNodeInit(1)
	m.SetJSDocObj(doc1)
	if m.JSDocObj() != doc1 {
		t.Fatalf("s2 JSDocObj mismatch")
	}

	SetCurrentNodeStore(s1)
	if m.JSDocObj() != nil {
		t.Fatalf("unregistered s2 JSDocObj still visible after switching back to s1")
	}
	if n1.JSDocObj() != doc1 {
		t.Fatalf("restored s1 JSDocObj mismatch: got %#v", n1.JSDocObj())
	}
}

func TestClonePropsFromAllCombinations(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	sf := &fakeSourceFileHolder{name: "app.js"}
	doc := &fakeJSDocHolder{comment: "/** @export */"}
	prop := &fakeColdProp{propType: 9, val: "extra"}

	// 1. Hot IRFactory template case: SourceFile only -> zero heap allocations!
	tmpl := AllocNodeInit(0)
	tmpl.SetSourceFileObj(sf)
	dstHot := AllocNodeInit(3)
	hotAllocs := testing.AllocsPerRun(100, func() {
		d := AllocNodeInit(3)
		if !d.ClonePropsFrom(tmpl) {
			t.Fatal("ClonePropsFrom returned false on clean node")
		}
	})
	if hotAllocs > 0 {
		t.Fatalf("hot ClonePropsFrom allocated %f objects, want 0", hotAllocs)
	}
	if !dstHot.ClonePropsFrom(tmpl) {
		t.Fatal("ClonePropsFrom failed")
	}
	if dstHot.SourceFileObj() != sf || dstHot.KindFlags() != 3 || dstHot.PropFlags() != 0 || dstHot.JSDocObj() != nil || dstHot.PropList() != nil {
		t.Fatalf("unexpected cloned props on dstHot")
	}

	// 2. Full case: SourceFile + PropFlags + JSDoc + cold PropList.
	srcFull := AllocNodeInit(1)
	srcFull.SetSourceFileObj(sf)
	srcFull.SetPropFlags(-0x55AA1234)
	srcFull.SetJSDocObj(doc)
	srcFull.SetPropList(prop)

	dstFull := AllocNodeInit(4)
	if !dstFull.ClonePropsFrom(srcFull) {
		t.Fatal("ClonePropsFrom failed on full src")
	}
	if dstFull.KindFlags() != 4 ||
		dstFull.SourceFileObj() != sf ||
		dstFull.PropFlags() != -0x55AA1234 ||
		dstFull.JSDocObj() != doc ||
		dstFull.PropList() != prop {
		t.Fatalf("full ClonePropsFrom mismatch: kind=%d sf=%#v pf=%x doc=%#v prop=%#v",
			dstFull.KindFlags(), dstFull.SourceFileObj(), dstFull.PropFlags(), dstFull.JSDocObj(), dstFull.PropList())
	}

	// 3. Rejects cloning into a node that already has ANY property (even if other == NilNode!).
	cases := []func(Node){
		func(n Node) { n.SetSourceFileObj(sf) },
		func(n Node) { n.SetPropFlags(1) },
		func(n Node) { n.SetJSDocObj(doc) },
		func(n Node) { n.SetPropList(prop) },
	}
	for idx, setup := range cases {
		dirty := AllocNodeInit(2)
		setup(dirty)
		if dirty.ClonePropsFrom(tmpl) {
			t.Fatalf("case %d: ClonePropsFrom should return false when target has existing props", idx)
		}
		if dirty.ClonePropsFrom(NilNode) {
			t.Fatalf("case %d: ClonePropsFrom(NilNode) should return false before dereferencing NilNode", idx)
		}
	}

	// 4. NilNode receiver or NilNode other (when target is clean) panics cleanly.
	assertPanics := func(name string, fn func()) {
		t.Helper()
		defer func() {
			if recover() == nil {
				t.Fatalf("%s: expected panic on NilNode dereference", name)
			}
		}()
		fn()
	}
	assertPanics("NilNode.Rec", func() { _ = NilNode.Rec().Token })
	assertPanics("NilNode.ClonePropsFrom", func() { NilNode.ClonePropsFrom(tmpl) })
	clean := AllocNodeInit(0)
	assertPanics("clean.ClonePropsFrom(NilNode)", func() { clean.ClonePropsFrom(NilNode) })
}

func TestSingleOwnerBatchAllocAndNestedPoolRace(t *testing.T) {
	s1 := NewNodeStore()
	s2 := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer func() {
		ExitSingleOwnerAlloc()
		SetCurrentNodeStore(prev)
	}()

	if !CanUseSingleOwnerAlloc() || IsSingleOwnerAllocActive() {
		t.Fatalf("unexpected initial single-owner state")
	}

	// 1. Single-owner multi-block allocation + exact NodeCount + FlushFastAlloc zeroing check.
	if !TryEnterSingleOwnerAlloc() {
		t.Fatal("TryEnterSingleOwnerAlloc should succeed")
	}
	if TryEnterSingleOwnerAlloc() {
		t.Fatal("nested TryEnterSingleOwnerAlloc must return false")
	}

	const multiBlkCount = 35000 // > 2 * 16384 blocks, crossing exact block boundaries with 1024 batches
	firstBatch := make([]Node, multiBlkCount)
	for i := 0; i < multiBlkCount; i++ {
		kind := uint8(i % 5)
		n := AllocNodeInit(kind)
		if n.Token() != 0xFF || n.Lineno() != -1 || n.Charno() != -1 || n.KindFlags() != kind {
			t.Fatalf("node %d bad init: tok=%d line=%d char=%d kind=%d",
				i, n.Token(), n.Lineno(), n.Charno(), n.KindFlags())
		}
		firstBatch[i] = n
		if i%777 == 0 && s1.NodeCount() != i+1 {
			t.Fatalf("step %d: NodeCount() = %d, want %d", i, s1.NodeCount(), i+1)
		}
	}
	if s1.NodeCount() != multiBlkCount {
		t.Fatalf("NodeCount() = %d, want %d", s1.NodeCount(), multiBlkCount)
	}

	// Flush mid-batch and verify s1.Alloc() gets a completely zeroed NodeData (not 0xFF / -1).
	FlushFastAlloc()
	if s1.NodeCount() != multiBlkCount {
		t.Fatalf("after FlushFastAlloc NodeCount() = %d, want %d", s1.NodeCount(), multiBlkCount)
	}
	zeroNode := s1.Alloc()
	if zeroNode.Token() != 0 || zeroNode.Rec().LinenoCharno != 0 || zeroNode.KindFlags() != 0 {
		t.Fatalf("rewound slot was not cleared before s.Alloc(): %+v", *zeroNode.Rec())
	}

	// Switching stores while single-owner active flushes s1 and switches batching to s2.
	_ = AllocNodeInit(2) // start a new batch in s1
	wantS1Count := multiBlkCount + 2
	SetCurrentNodeStore(s2)
	if s1.NodeCount() != wantS1Count {
		t.Fatalf("s1 NodeCount after store switch = %d, want %d", s1.NodeCount(), wantS1Count)
	}
	s2Node := AllocNodeInit(4)
	if s2Node.KindFlags() != 4 || s2.NodeCount() != 1 {
		t.Fatalf("s2 single-owner alloc mismatch: kind=%d count=%d", s2Node.KindFlags(), s2.NodeCount())
	}
	ExitSingleOwnerAlloc()
	SetCurrentNodeStore(s1)

	// 2. Full lifecycle race stress:
	// main (single-owner) -> 1-worker CompilerExecutor -> nested 4-worker pool (+ nested Thread.Start)
	// -> 1-worker CompilerExecutor (CHECKS/OPT fast path resumed) -> main (single-owner).
	sf := &fakeSourceFileHolder{name: "concurrent.js"}
	doc := &fakeJSDocHolder{comment: "/** @struct */"}
	tc := &fakeTCHolder{name: "NumberType"}

	for iter := 0; iter < 5; iter++ {
		iterStore := NewNodeStore()
		SetCurrentNodeStore(iterStore)

		var collectedMu sync.Mutex
		var allIterNodes []Node
		recordNodes := func(ns []Node) {
			collectedMu.Lock()
			allIterNodes = append(allIterNodes, ns...)
			collectedMu.Unlock()
		}

		// Phase A: main goroutine single-owner allocations.
		EnterSingleOwnerAlloc()
		mainPre := make([]Node, 300)
		for i := range mainPre {
			n := AllocNodeInit(1)
			n.SetSourceFileObj(sf)
			mainPre[i] = n
		}
		recordNodes(mainPre)

		// Hand off to 1-worker CompilerExecutor goroutine.
		ExitSingleOwnerAlloc()
		compilerDone := make(chan struct{})
		go func() {
			entered := TryEnterSingleOwnerAlloc()
			if !entered {
				t.Errorf("iter %d: compiler worker failed TryEnterSingleOwnerAlloc", iter)
			}
			defer func() {
				if entered {
					ExitSingleOwnerAlloc()
				}
				close(compilerDone)
			}()

			// Phase B: single-owner setup on compiler thread.
			compPre := make([]Node, 1500)
			for i := range compPre {
				n := AllocNodeInit(2)
				n.SetSourceFileObj(sf)
				n.SetJSDocObj(doc)
				compPre[i] = n
			}
			recordNodes(compPre)
			tmplNode := compPre[0]

			// Phase C: nested 4-worker pool (PrebuildAst --threads 4) + nested Thread.Start().
			const numWorkers = 4
			const numTasks = 12
			tasksCh := make(chan func(), numTasks)
			var poolWg sync.WaitGroup
			poolWg.Add(numWorkers)
			for w := 0; w < numWorkers; w++ {
				go func() {
					defer poolWg.Done()
					for fn := range tasksCh {
						fn()
					}
				}()
			}

			futures := make([]chan struct{}, numTasks)
			for taskIdx := 0; taskIdx < numTasks; taskIdx++ {
				taskIdx := taskIdx
				done := make(chan struct{})
				futures[taskIdx] = done
				// Must be called on submitting goroutine BEFORE sending task on channel!
				SuspendSingleOwnerAlloc()
				tasksCh <- func() {
					defer func() {
						ResumeSingleOwnerAlloc()
						close(done)
					}()
					if TryEnterSingleOwnerAlloc() {
						t.Errorf("worker inside suspended pool must not enter single-owner alloc")
					}
					local := make([]Node, 400)
					for j := range local {
						n := AllocNodeInit(3)
						n.ClonePropsFrom(tmplNode)
						n.SetTypeOrColor(tc)
						local[j] = n
					}
					// Also spawn a nested Java Thread from one task to stress nested suspend/resume.
					if taskIdx%4 == 0 {
						threadDone := make(chan struct{})
						SuspendSingleOwnerAlloc()
						go func() {
							defer func() {
								ResumeSingleOwnerAlloc()
								close(threadDone)
							}()
							tn := AllocNodeInit(4)
							tn.ClonePropsFrom(tmplNode)
							recordNodes([]Node{tn})
						}()
						<-threadDone
					}
					recordNodes(local)
				}
			}
			close(tasksCh)

			// Drain futures while ALSO allocating on the compiler thread between future.Get() calls!
			drainNodes := make([]Node, numTasks)
			for i, done := range futures {
				<-done
				drainNodes[i] = AllocNodeInit(1)
			}
			poolWg.Wait()
			recordNodes(drainNodes)

			if !CanUseSingleOwnerAlloc() || !IsSingleOwnerAllocActive() {
				t.Errorf("iter %d: expected single-owner fast path to resume after pool drain", iter)
			}

			// Phase D: CHECKS & OPT on compiler thread (resumed fast 1024-batch path!).
			compPost := make([]Node, 2500)
			for i := range compPost {
				n := AllocNodeInit(0)
				n.ClonePropsFrom(tmplNode)
				n.SetTypeOrColor(tc)
				compPost[i] = n
			}
			recordNodes(compPost)
		}()
		<-compilerDone

		// Phase E: back on main goroutine.
		EnterSingleOwnerAlloc()
		mainPost := make([]Node, 300)
		for i := range mainPost {
			mainPost[i] = AllocNodeInit(2)
		}
		recordNodes(mainPost)
		ExitSingleOwnerAlloc()

		// Verify 100% unique handles and exact NodeCount.
		if iterStore.NodeCount() != len(allIterNodes) {
			t.Fatalf("iter %d: NodeCount() = %d, want %d", iter, iterStore.NodeCount(), len(allIterNodes))
		}
		seen := make(map[Node]struct{}, len(allIterNodes))
		for _, n := range allIterNodes {
			if n == 0 {
				t.Fatalf("iter %d: allocated NilNode", iter)
			}
			if _, dup := seen[n]; dup {
				t.Fatalf("iter %d: duplicate node handle %d", iter, n)
			}
			seen[n] = struct{}{}
			if n.Token() != 0xFF || n.Lineno() != -1 {
				t.Fatalf("iter %d: node %d corrupted: tok=%d line=%d", iter, n, n.Token(), n.Lineno())
			}
		}
	}
}

func TestNilNodeRecRuntimeErrorAndOverflowLengths(t *testing.T) {
	assertRuntimeError := func(name string, fn func()) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("%s: expected runtime.Error panic on Node(0), got nil", name)
			}
			if _, ok := r.(runtime.Error); !ok {
				t.Fatalf("%s: expected runtime.Error on Node(0), got %T (%v)", name, r, r)
			}
		}()
		fn()
	}

	assertRuntimeError("Node(0).Rec()", func() {
		_ = Node(0).Rec().Next
	})
	assertRuntimeError("Node(0).SourceFile()", func() {
		_ = Node(0).SourceFile()
	})
	assertRuntimeError("Node(0).SetSourceFile(1)", func() {
		Node(0).SetSourceFile(1)
	})
	assertRuntimeError("Node(0).SetSourceFile(70000)", func() {
		Node(0).SetSourceFile(70000)
	})
	assertRuntimeError("Node(0).SetLength(10)", func() {
		Node(0).SetLength(10)
	})
	assertRuntimeError("Node(0).SetLength(70000)", func() {
		Node(0).SetLength(70000)
	})
	assertRuntimeError("Node(0).SetLength(-1)", func() {
		Node(0).SetLength(-1)
	})

	s1 := NewNodeStore()
	s2 := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer SetCurrentNodeStore(prev)

	// Verify SetLength(0) does NOT allocate any aux overflow entry.
	nInline := AllocNodeToken(0, 32)
	nInline.SetLength(0)
	if s1.auxLen() != 0 {
		t.Fatalf("SetLength(0) allocated aux overflow entry eagerly: %d", s1.auxLen())
	}
	if got := nInline.Length(); got != 0 {
		t.Fatalf("Length() = %d, want 0", got)
	}

	// Verify full int32 range transitions: 100000 -> 42 -> 65535 -> -7 -> MinInt32 -> MaxInt32 -> 0.
	seq := []int32{100000, 42, 65534, 65535, -7, math.MinInt32, math.MaxInt32, 0}
	for _, want := range seq {
		nInline.SetLength(want)
		if got := nInline.Length(); got != want {
			t.Fatalf("after SetLength(%d): Length() = %d", want, got)
		}
	}

	// Store switching + static template node length preservation.
	nS1 := AllocNodeToken(3, 35)
	nS1.SetLength(123456)
	tmpl := NewTemplateNode(88)
	tmpl.SetLength(234567)
	defer tmpl.SetLength(0)

	SetCurrentNodeStore(s2)
	if got := nS1.Length(); got != 0 {
		t.Fatalf("unregistered s1 Length visible in s2: got %d", got)
	}
	if got := tmpl.Length(); got != 234567 {
		t.Fatalf("static template Length lost across store switch: got %d", got)
	}
	nS2 := AllocNodeToken(1, 33)
	nS2.SetLength(-42)
	if got := nS2.Length(); got != -42 {
		t.Fatalf("s2 negative Length = %d, want -42", got)
	}

	SetCurrentNodeStore(s1)
	if got := nS1.Length(); got != 123456 {
		t.Fatalf("restored s1 Length = %d, want 123456", got)
	}

	// SourceFileByHandle roundtrip + zero/out-of-bounds safety.
	sf := &fakeSourceFileHolder{name: "bundle.js"}
	nS1.SetSourceFileObj(sf)
	if got := SourceFileByHandle(nS1.SourceFile()); got != sf {
		t.Fatalf("SourceFileByHandle(%d) = %#v, want %#v", nS1.SourceFile(), got, sf)
	}
	if SourceFileByHandle(0) != nil || SourceFileByHandle(0xFFFFFFFF) != nil {
		t.Fatalf("SourceFileByHandle on 0 / 0xFFFFFFFF should return nil")
	}
}

func TestSourceFileOverflowHandlesAndClonePropsFrom(t *testing.T) {
	s1 := NewNodeStore()
	s2 := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer SetCurrentNodeStore(prev)

	// 1. SourceFile16 inline (< 0xFFFF) + sfOverflowTop (>= 0xFFFF) across the full uint32 range!
	n := AllocNodeInit(1)
	fullSeq := []uint32{0, 1, 42, 65534, 0xFFFF, 65536, 100000, 5, 0xCAFEBABE, 0, 0xFFFFFFFF, 65534, 77777}
	for _, want := range fullSeq {
		n.SetSourceFile(want)
		if got := n.SourceFile(); got != want {
			t.Fatalf("after SetSourceFile(%d): SourceFile() = %d", want, got)
		}
		wantSF16 := uint16(0xFFFF)
		if want < 0xFFFF {
			wantSF16 = uint16(want)
		}
		if n.Rec().SourceFile16 != wantSF16 {
			t.Fatalf("expected NodeData.SourceFile16 = 0x%X, got 0x%X", wantSF16, n.Rec().SourceFile16)
		}
	}

	// 2. Real SourceFileHandleHolder objects pushed past > 65535 handles:
	// test exact boundary 65534 (0xFFFE), 65535 (0xFFFF), 65536 (0x10000), and > 65540!
	for sourceFileHandles.count.Load() < 65533 {
		_ = sourceFileHandles.Add(&sourceFileHandlesMu, &fakeSourceFileHolder{name: "pad.js"})
	}
	overflowHolders := make([]*fakeSourceFileHolder, 12)
	for i := range overflowHolders {
		overflowHolders[i] = &fakeSourceFileHolder{name: "overflow_file.js"}
	}

	doc := &fakeJSDocHolder{comment: "/** @fileoverview */"}
	prop := &fakeColdProp{propType: 3, val: "meta"}

	for i, sf := range overflowHolders {
		src := AllocNodeToken(2, 40)
		src.SetSourceFileObj(sf)
		src.SetPropFlags(int32(1000 + i))
		if i%2 == 0 {
			src.SetJSDocObj(doc)
		}
		if i%3 == 0 {
			src.SetPropList(prop)
		}

		h := src.SourceFile()
		if h == 0 || sf.sfHandle.Load() != h {
			t.Fatalf("holder[%d] invalid cached handle: %d vs %d", i, h, sf.sfHandle.Load())
		}
		if got := src.SourceFileObj(); got != sf {
			t.Fatalf("holder[%d] (handle=%d) SourceFileObj() = %#v, want %#v", i, h, got, sf)
		}
		if got := SourceFileByHandle(h); got != sf {
			t.Fatalf("holder[%d] (handle=%d) SourceFileByHandle() = %#v, want %#v", i, h, got, sf)
		}

		// ClonePropsFrom must copy full 32-bit SourceFile handle!
		dst := AllocNodeToken(4, 41)
		if !dst.ClonePropsFrom(src) {
			t.Fatalf("holder[%d] (handle=%d) ClonePropsFrom returned false on clean node", i, h)
		}
		if dst.SourceFile() != h || dst.SourceFileObj() != sf {
			t.Fatalf("holder[%d] (handle=%d) cloned SourceFile mismatch: handle=%d obj=%#v",
				i, h, dst.SourceFile(), dst.SourceFileObj())
		}
		if dst.PropFlags() != int32(1000+i) || dst.KindFlags() != 4 {
			t.Fatalf("holder[%d] cloned PropFlags/KindFlags mismatch", i)
		}
		// Rejects second ClonePropsFrom when dst already has SourceFile set!
		if dst.ClonePropsFrom(src) {
			t.Fatalf("holder[%d] ClonePropsFrom should return false when target already has SourceFile", i)
		}
	}

	// Verify at least one holder actually crossed > 65535!
	lastSF := overflowHolders[len(overflowHolders)-1]
	lastH := lastSF.sfHandle.Load()
	if lastH <= 65535 {
		t.Fatalf("expected last overflow holder handle > 65535, got %d", lastH)
	}

	// 4. Store switching & static template node with overflow SourceFile (> 65535).
	nS1 := AllocNodeInit(1)
	nS1.SetSourceFileObj(lastSF)
	tmpl := NewTemplateNode(91)
	tmpl.SetSourceFileObj(lastSF)
	defer tmpl.SetSourceFileObj(nil)

	SetCurrentNodeStore(s2)
	if tmpl.SourceFile() != lastH || tmpl.SourceFileObj() != lastSF {
		t.Fatalf("static template overflow SourceFile lost across store switch: h=%d", tmpl.SourceFile())
	}
	nS2Clone := AllocNodeInit(2)
	if !nS2Clone.ClonePropsFrom(tmpl) || nS2Clone.SourceFile() != lastH || nS2Clone.SourceFileObj() != lastSF {
		t.Fatalf("cloning overflow SourceFile from static template into s2 failed")
	}

	SetCurrentNodeStore(s1)
	if nS1.SourceFile() != lastH || nS1.SourceFileObj() != lastSF {
		t.Fatalf("restored s1 overflow SourceFile mismatch: h=%d", nS1.SourceFile())
	}
}

func TestAllocNodeTokenAndFastBatchFlushAndConcurrentPools(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer func() {
		ExitSingleOwnerAlloc()
		SetCurrentNodeStore(prev)
	}()

	if !TryEnterSingleOwnerAlloc() {
		t.Fatal("TryEnterSingleOwnerAlloc should succeed")
	}

	// 1. Interleave AllocNodeToken, AllocNodeFast, and AllocNodeInit across > 2 blocks,
	// flushing mid-batch repeatedly to verify rewound unconsumed slots remain pristine zeroed NodeData{}.
	const count = 35000
	nodes := make([]Node, 0, count)
	wantTok := make([]uint8, 0, count)
	wantKind := make([]uint8, 0, count)

	for i := 0; i < count; i++ {
		kind := uint8(i % 5)
		tok := uint8((i * 7) % 200)
		var n Node
		switch i % 3 {
		case 0:
			// Pass internal flag bits (0xC0) on every 9th call to verify masking.
			rawKind := kind
			if i%9 == 0 {
				rawKind |= kindFlagInternalMask
			}
			n = AllocNodeToken(rawKind, tok)
		case 1:
			n = AllocNodeFast()
			rec := n.Rec()
			if *rec != (NodeData{}) {
				t.Fatalf("step %d: AllocNodeFast returned non-zero NodeData: %+v", i, *rec)
			}
			rec.Token = tok
			rec.LinenoCharno = -1
			rec.KindFlags = kind
		default:
			n = AllocNodeInit(kind | kindFlagInternalMask)
			tok = 0xFF
		}
		nodes = append(nodes, n)
		wantTok = append(wantTok, tok)
		wantKind = append(wantKind, kind)

		if i%509 == 250 {
			FlushFastAlloc()
			if s.NodeCount() != len(nodes) {
				t.Fatalf("after FlushFastAlloc at %d: NodeCount()=%d, want %d", i, s.NodeCount(), len(nodes))
			}
			z := s.Alloc()
			if *z.Rec() != (NodeData{}) {
				t.Fatalf("after FlushFastAlloc at %d: s.Alloc() slot dirty: %+v", i, *z.Rec())
			}
			z.Rec().Token = 77
			z.Rec().LinenoCharno = -1
			z.Rec().KindFlags = 2
			nodes = append(nodes, z)
			wantTok = append(wantTok, 77)
			wantKind = append(wantKind, 2)
		}
	}

	FlushFastAlloc()
	if s.NodeCount() != len(nodes) {
		t.Fatalf("final NodeCount() = %d, want %d", s.NodeCount(), len(nodes))
	}
	for i, n := range nodes {
		if n.Token() != wantTok[i] || n.Lineno() != -1 || n.Charno() != -1 || n.KindFlags() != wantKind[i] {
			t.Fatalf("node[%d] (%d) mismatch: tok=%d(want %d) line=%d kind=%d(want %d)",
				i, n, n.Token(), wantTok[i], n.Lineno(), n.KindFlags(), wantKind[i])
		}
		if n.JSDocObj() != nil || n.PropList() != nil {
			t.Fatalf("node[%d] leaked internal kindFlags", i)
		}
	}

	// 2. Concurrent pool stress mixing AllocNodeToken, AllocNodeFast, SetLength overflow, and Suspend/Resume.
	const workers = 6
	const tasks = 18
	tasksCh := make(chan func(), tasks)
	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for fn := range tasksCh {
				fn()
			}
		}()
	}

	type nodeExpect struct {
		n      Node
		tok    uint8
		kind   uint8
		length int32
	}
	var mu sync.Mutex
	var conc []nodeExpect
	dones := make([]chan struct{}, tasks)
	for ti := range tasks {
		ti := ti
		done := make(chan struct{})
		dones[ti] = done
		SuspendSingleOwnerAlloc()
		tasksCh <- func() {
			defer func() {
				ResumeSingleOwnerAlloc()
				close(done)
			}()
			local := make([]nodeExpect, 300)
			for j := range local {
				kind := uint8((ti + j) % 5)
				tok := uint8((ti*13 + j) % 180)
				var n Node
				if j%2 == 0 {
					n = AllocNodeToken(kind, tok)
				} else {
					n = AllocNodeFast()
					rec := n.Rec()
					rec.Token = tok
					rec.LinenoCharno = -1
					rec.KindFlags = kind
				}
				ln := int32(j * 300) // crosses 65535 for j >= 219!
				if j%17 == 0 {
					ln = int32(-j - 1)
				}
				n.SetLength(ln)
				local[j] = nodeExpect{n: n, tok: tok, kind: kind, length: ln}
			}
			mu.Lock()
			conc = append(conc, local...)
			mu.Unlock()
		}
	}
	close(tasksCh)
	for _, d := range dones {
		<-d
	}
	wg.Wait()

	// Resume single-owner fast path after pool completes.
	post := AllocNodeToken(4, 99)
	post.SetLength(80000)
	ExitSingleOwnerAlloc()

	if post.Token() != 99 || post.KindFlags() != 4 || post.Length() != 80000 {
		t.Fatalf("post-pool AllocNodeToken mismatch: tok=%d kind=%d len=%d",
			post.Token(), post.KindFlags(), post.Length())
	}

	seen := make(map[Node]struct{}, len(nodes)+len(conc)+1)
	for _, n := range nodes {
		seen[n] = struct{}{}
	}
	seen[post] = struct{}{}
	for idx, e := range conc {
		if _, dup := seen[e.n]; dup {
			t.Fatalf("concurrent duplicate handle %d at idx %d", e.n, idx)
		}
		seen[e.n] = struct{}{}
		if e.n.Token() != e.tok || e.n.KindFlags() != e.kind || e.n.Lineno() != -1 || e.n.Length() != e.length {
			t.Fatalf("conc[%d] corrupted: tok=%d/%d kind=%d/%d line=%d len=%d/%d",
				idx, e.n.Token(), e.tok, e.n.KindFlags(), e.kind, e.n.Lineno(), e.n.Length(), e.length)
		}
	}
}

func assertPanicsNPE(t *testing.T, label string, fn func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("%s: expected NullPointerException panic, got nil", label)
		}
		if _, ok := r.(*javalang.NullPointerException); !ok {
			t.Fatalf("%s: expected *javalang.NullPointerException, got %T (%v)", label, r, r)
		}
	}()
	fn()
}

func assertPanicsIllegalState(t *testing.T, label string, fn func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("%s: expected IllegalStateException panic, got nil", label)
		}
		if _, ok := r.(*javalang.IllegalStateException); !ok {
			t.Fatalf("%s: expected *javalang.IllegalStateException, got %T (%v)", label, r, r)
		}
	}()
	fn()
}

func TestNodeAllocatorInitLocalFlushTemplateAndHelpers(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	// 1. InitLocal + Flush returns unused trailing IDs in single-owner/sequential use.
	tmplA := s.Alloc()
	tmplA.SetToken(94)
	tmplA.SetSourceFile(17)
	tmplA.SetPropFlags(0x01000000)

	tmplB := s.Alloc()
	tmplB.SetToken(94)
	tmplB.SetSourceFile(70005) // 0xFFFF overflow column!
	tmplB.SetPropFlags(0x02000004)

	tmplCold := s.Alloc()
	tmplCold.SetToken(94)
	tmplCold.SetSourceFile(29)
	tmplCold.SetPropFlags(0x00000010)
	tmplCold.SetPropList(&fakeColdProp{propType: 7, val: "cold-prop"})
	tmplCold.SetJSDocObj("jsdoc-header")

	beforeLocal := s.NodeCount()
	var alloc NodeAllocator
	alloc.InitLocal()
	alloc.SetTemplate(tmplA)

	strHello := Intern("hello_v17")
	leaf1 := alloc.AllocNode1(10, tmplA)
	strNode := alloc.AllocStringNode(32, strHello, tmplA)
	numZero := alloc.AllocNumberNode(33, 0.0, tmplA)
	numNegZero := alloc.AllocNumberNode(33, math.Copysign(0, -1), tmplA)
	numPi := alloc.AllocNumberNode(33, 3.141592653589793, tmplA)
	p2 := alloc.AllocNode2(40, leaf1, tmplA)
	p3 := alloc.AllocNode3(50, strNode, numZero, tmplA)
	p4 := alloc.AllocNode4(60, p2, p3, numNegZero, tmplA)

	// Switch template mid-batch to 0xFFFF overflow template, then cold template, then clear.
	alloc.SetTemplate(tmplB)
	overLeaf := alloc.AllocNode1(11, tmplB)
	overNum := alloc.AllocNumberNode(33, 42.5, tmplB)
	overParent := alloc.AllocNode2(41, overLeaf, tmplB)

	alloc.SetTemplate(tmplCold)
	coldStr := alloc.AllocStringNode(35, strHello, tmplCold)
	coldLeaf := alloc.AllocNode1(12, tmplCold)
	coldParent := alloc.AllocNode2(42, coldLeaf, tmplCold)

	alloc.SetTemplate(0)
	bareLeaf := alloc.AllocNode1(13, 0)

	alloc.Flush()
	// Exactly 15 nodes allocated from alloc; Flush must rewind the remaining 241 IDs!
	if got, want := s.NodeCount()-beforeLocal, 15; got != want {
		t.Fatalf("after Flush, NodeCount delta = %d, want %d", got, want)
	}

	// Verify inline template stamping, tokens, LinenoCharno, KindFlags, and child links.
	for _, n := range []Node{leaf1, strNode, numZero, numNegZero, numPi, p2, p3, p4} {
		if n.SourceFile() != 17 || n.PropFlags() != 0x01000000 || n.Lineno() != -1 || n.Charno() != -1 {
			t.Fatalf("tmplA node %d mismatch: sf=%d pf=0x%x line=%d char=%d",
				n, n.SourceFile(), n.PropFlags(), n.Lineno(), n.Charno())
		}
	}
	if strNode.KindFlags() != 3 || String(StringID(strNode.Payload())) != "hello_v17" {
		t.Fatalf("strNode mismatch: kind=%d str=%q", strNode.KindFlags(), String(StringID(strNode.Payload())))
	}
	if numZero.KindFlags() != 1 || numZero.Payload() != 0 || numZero.NumberVal() != 0 || math.Signbit(numZero.NumberVal()) {
		t.Fatalf("+0.0 NumberNode mismatch: kind=%d payload=%d val=%v sign=%v",
			numZero.KindFlags(), numZero.Payload(), numZero.NumberVal(), math.Signbit(numZero.NumberVal()))
	}
	if numNegZero.KindFlags() != 1 || numNegZero.Payload() == 0 || numNegZero.NumberVal() != 0 || !math.Signbit(numNegZero.NumberVal()) {
		t.Fatalf("-0.0 NumberNode lost sign bit: kind=%d payload=%d val=%v sign=%v",
			numNegZero.KindFlags(), numNegZero.Payload(), numNegZero.NumberVal(), math.Signbit(numNegZero.NumberVal()))
	}
	if numPi.NumberVal() != 3.141592653589793 {
		t.Fatalf("numPi = %v, want pi", numPi.NumberVal())
	}
	if p4.FirstChild() != p2 || p2.Next() != p3 || p3.Next() != numNegZero || numNegZero.Next() != 0 {
		t.Fatalf("p4 forward links corrupted")
	}
	if p2.Previous() != numNegZero || p3.Previous() != p2 || numNegZero.Previous() != p3 {
		t.Fatalf("p4 circular previous links corrupted")
	}
	if p2.Parent() != p4 || p3.Parent() != p4 || numNegZero.Parent() != p4 {
		t.Fatalf("p4 parent links corrupted")
	}
	if overParent.SourceFile() != 70005 || overParent.PropFlags() != 0x02000004 || overNum.NumberVal() != 42.5 || overNum.SourceFile() != 70005 {
		t.Fatalf("overflow template mismatch: sf=%d pf=0x%x num=%v",
			overParent.SourceFile(), overParent.PropFlags(), overNum.NumberVal())
	}
	if coldStr.KindFlags() != 3 || coldStr.SourceFile() != 29 || coldStr.JSDocObj() != "jsdoc-header" {
		t.Fatalf("coldStr mismatch: kind=%d sf=%d jsdoc=%v", coldStr.KindFlags(), coldStr.SourceFile(), coldStr.JSDocObj())
	}
	if cp, ok := coldParent.PropList().(*fakeColdProp); !ok || cp.val != "cold-prop" || coldParent.JSDocObj() != "jsdoc-header" {
		t.Fatalf("coldParent props mismatch: prop=%#v jsdoc=%v", coldParent.PropList(), coldParent.JSDocObj())
	}
	if bareLeaf.SourceFile() != 0 || bareLeaf.PropFlags() != 0 || bareLeaf.Token() != 13 {
		t.Fatalf("bareLeaf mismatch: sf=%d pf=%d tok=%d", bareLeaf.SourceFile(), bareLeaf.PropFlags(), bareLeaf.Token())
	}

	// 2. Cross-block boundary stress (> 18,000 nodes) + cross-block child linking!
	alloc.InitLocal()
	alloc.SetTemplate(tmplA)
	firstBlockChild := alloc.AllocNode1(1, tmplA)
	for range blockLen + 500 {
		_ = alloc.AllocNode1(2, tmplA)
	}
	secondBlockChild := alloc.AllocNode1(3, tmplA)
	if (uint32(firstBlockChild) >> blockShift) == (uint32(secondBlockChild) >> blockShift) {
		t.Fatalf("expected children in distinct nodeBlocks")
	}
	crossParent := alloc.AllocNode3(77, firstBlockChild, secondBlockChild, tmplA)
	alloc.Flush()
	if crossParent.FirstChild() != firstBlockChild || firstBlockChild.Next() != secondBlockChild ||
		secondBlockChild.Previous() != firstBlockChild || firstBlockChild.Previous() != secondBlockChild ||
		firstBlockChild.Parent() != crossParent || secondBlockChild.Parent() != crossParent {
		t.Fatalf("cross-block AllocNode3 links corrupted")
	}

	// 3. Validation panics (null child, already-parented child, duplicate child).
	alloc.InitLocal()
	alloc.SetTemplate(tmplA)
	assertPanicsNPE(t, "AllocNode2(nil)", func() { alloc.AllocNode2(10, 0, tmplA) })
	assertPanicsNPE(t, "AllocNode3(c1, nil)", func() {
		c := alloc.AllocNode1(1, tmplA)
		alloc.AllocNode3(10, c, 0, tmplA)
	})
	assertPanicsNPE(t, "AllocNode4(c1, c2, nil)", func() {
		c1 := alloc.AllocNode1(1, tmplA)
		c2 := alloc.AllocNode1(2, tmplA)
		alloc.AllocNode4(10, c1, c2, 0, tmplA)
	})
	assertPanicsIllegalState(t, "AllocNode2(attached)", func() {
		alloc.AllocNode2(10, firstBlockChild, tmplA)
	})
	assertPanicsIllegalState(t, "AllocNode3(dup c, c)", func() {
		c := alloc.AllocNode1(1, tmplA)
		alloc.AllocNode3(10, c, c, tmplA)
	})
	assertPanicsIllegalState(t, "AllocNode4(dup c1, c2, c1)", func() {
		c1 := alloc.AllocNode1(1, tmplA)
		c2 := alloc.AllocNode1(2, tmplA)
		alloc.AllocNode4(10, c1, c2, c1, tmplA)
	})
	alloc.Flush()

	// 4. Static store fallback when CurrentNodeStore is deactivated.
	SetCurrentNodeStore(nil)
	var staticAlloc NodeAllocator
	staticAlloc.InitLocal()
	sn := staticAlloc.AllocNumberNode(33, math.Copysign(0, -1), 0)
	staticAlloc.Flush()
	if uint32(sn)&staticBit == 0 || !math.Signbit(sn.NumberVal()) {
		t.Fatalf("static store AllocNumberNode(-0.0) failed: id=0x%x val=%v", sn, sn.NumberVal())
	}
}

func TestConcurrentNodeAllocatorWorkersAndSingleOwnerLockElision(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	EnterSingleOwnerAlloc()
	defer ExitSingleOwnerAlloc()

	const rounds = 4
	const workers = 8
	const filesPerWorker = 20

	strID := Intern("parallel_ident")

	for round := range rounds {
		// Phase 1: Single-owner fast-path (no mutexes on propShards / auxShards).
		singleNodes := make([]Node, 512)
		for i := range singleNodes {
			n := AllocNodeToken(0, uint8(i%120))
			n.SetPropList(&fakeColdProp{propType: byte(round), val: "single"})
			n.SetAux(NodeAux{JSDoc: "single-doc", Extra: i})
			if i%5 == 0 {
				n.SetPropList(nil)
				n.SetAux(NodeAux{})
			}
			singleNodes[i] = n
		}

		// Phase 2: Suspend single-owner mode and run parallel PrebuildAst-style workers
		// while the parent goroutine ALSO reads/writes PropList and Aux concurrently!
		var wg sync.WaitGroup
		workerRoots := make([][]Node, workers)
		dones := make([]chan struct{}, workers)

		for w := range workers {
			w := w
			done := make(chan struct{})
			dones[w] = done
			SuspendSingleOwnerAlloc()
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() {
					ResumeSingleOwnerAlloc()
					close(done)
				}()
				roots := make([]Node, 0, filesPerWorker)
				for f := range filesPerWorker {
					var a NodeAllocator
					a.InitLocal()
					tmpl := a.AllocNode1(94, 0)
					tmpl.SetSourceFile(uint32(w*1000 + f + 1))
					tmpl.SetPropFlags(int32(round + 1))
					a.SetTemplate(tmpl)

					var prevStmt Node
					var root Node
					for k := range 45 {
						s1 := a.AllocStringNode(32, strID, tmpl)
						n1 := a.AllocNumberNode(33, float64(k), tmpl)
						bin := a.AllocNode3(13, s1, n1, tmpl)
						stmt := a.AllocNode2(93, bin, tmpl)
						if k%7 == 0 {
							stmt.SetPropList(&fakeColdProp{propType: byte(w), val: "worker"})
							stmt.SetAux(NodeAux{JSDoc: "worker-doc", Extra: w*1000 + k})
						}
						if prevStmt != 0 {
							root = a.AllocNode3(90, prevStmt, stmt, tmpl)
							prevStmt = root
						} else {
							prevStmt = stmt
							root = stmt
						}
					}
					a.Flush()
					roots = append(roots, root)
				}
				workerRoots[w] = roots
			}()
		}

		// Parent goroutine concurrently accesses PropList & Aux on singleNodes while workers are active!
		for i, n := range singleNodes {
			if i%5 != 0 {
				if cp, ok := n.PropList().(*fakeColdProp); !ok || cp.val != "single" {
					t.Fatalf("concurrent parent read PropList mismatch: %#v", n.PropList())
				}
				if aux, ok := n.GetAux(); !ok || aux.Extra != i {
					t.Fatalf("concurrent parent read Aux mismatch: %+v", aux)
				}
			}
		}

		for _, d := range dones {
			<-d
		}
		wg.Wait()

		// Phase 3: Back in single-owner fast path (fastSuspendCount == 0)!
		// Read and mutate worker-created PropList and Aux entries without locks (must be 100% race-free!).
		for w, roots := range workerRoots {
			for f, r := range roots {
				wantSF := uint32(w*1000 + f + 1)
				if r.SourceFile() != wantSF || r.PropFlags() != int32(round+1) {
					t.Fatalf("worker[%d][%d] root template mismatch: sf=%d/%d pf=%d",
						w, f, r.SourceFile(), wantSF, r.PropFlags())
				}
				r.SetPropList(&fakeColdProp{propType: 99, val: "post-merge"})
				r.SetAux(NodeAux{JSDoc: "post-merge", Extra: f})
				if cp, ok := r.PropList().(*fakeColdProp); !ok || cp.val != "post-merge" {
					t.Fatalf("post-merge PropList mismatch")
				}
				if aux, ok := r.GetAux(); !ok || aux.Extra != f {
					t.Fatalf("post-merge Aux mismatch")
				}
				r.SetPropList(nil)
				r.SetAux(NodeAux{})
			}
		}
	}
}

func TestFirstNextTokHasColdPropsAndChangeTimeColumn(t *testing.T) {
	assertRuntimeError := func(name string, fn func()) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("%s: expected runtime.Error panic on Node(0), got nil", name)
			}
			if _, ok := r.(runtime.Error); !ok {
				t.Fatalf("%s: expected runtime.Error on Node(0), got %T (%v)", name, r, r)
			}
		}()
		fn()
	}

	assertRuntimeError("Node(0).FirstNextTok()", func() {
		_, _, _ = Node(0).FirstNextTok()
	})
	assertRuntimeError("Node(0).HasColdProps()", func() {
		_ = Node(0).HasColdProps()
	})
	assertRuntimeError("Node(0).ChangeTime()", func() {
		_ = Node(0).ChangeTime()
	})
	assertRuntimeError("Node(0).SetChangeTime(0)", func() {
		Node(0).SetChangeTime(0)
	})
	assertRuntimeError("Node(0).SetChangeTime(7)", func() {
		Node(0).SetChangeTime(7)
	})

	s1 := NewNodeStore()
	s2 := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer SetCurrentNodeStore(prev)

	// 1. FirstNextTok & HasColdProps fast accessors.
	var alloc NodeAllocator
	alloc.InitLocal()
	c1 := alloc.AllocNode1(11, 0)
	c2 := alloc.AllocNode1(22, 0)
	parent := alloc.AllocNode3(99, c1, c2, 0)
	alloc.Flush()

	if f, nx, tok := parent.FirstNextTok(); f != c1 || nx != 0 || tok != 99 {
		t.Fatalf("parent.FirstNextTok() = (%d, %d, %d), want (%d, 0, 99)", f, nx, tok, c1)
	}
	if f, nx, tok := c1.FirstNextTok(); f != 0 || nx != c2 || tok != 11 {
		t.Fatalf("c1.FirstNextTok() = (%d, %d, %d), want (0, %d, 11)", f, nx, tok, c2)
	}
	if f, nx, tok := c2.FirstNextTok(); f != 0 || nx != 0 || tok != 22 {
		t.Fatalf("c2.FirstNextTok() = (%d, %d, %d), want (0, 0, 22)", f, nx, tok)
	}

	if parent.HasColdProps() {
		t.Fatalf("expected HasColdProps() == false on fresh node")
	}
	var typedNilProp *fakeColdProp
	parent.SetPropList(any(typedNilProp))
	if parent.HasColdProps() {
		t.Fatalf("expected HasColdProps() == false after typed-nil SetPropList")
	}
	prop := &fakeColdProp{propType: 5, val: "cold-v19"}
	parent.SetPropList(prop)
	if !parent.HasColdProps() {
		t.Fatalf("expected HasColdProps() == true after non-nil SetPropList")
	}
	parent.SetPropList(nil)
	if parent.HasColdProps() {
		t.Fatalf("expected HasColdProps() == false after nil SetPropList")
	}

	// 2. ChangeTime lazy aux allocation + positive/negative/zero/extremes + KindFlags isolation.
	n := AllocNodeInit(0x07)
	n.SetChangeTime(0)
	if s1.auxLen() != 0 {
		t.Fatalf("SetChangeTime(0) allocated aux entry eagerly: %d", s1.auxLen())
	}
	if got := n.ChangeTime(); got != 0 {
		t.Fatalf("default ChangeTime() = %d, want 0", got)
	}
	readAllocs := testing.AllocsPerRun(100, func() {
		if n.ChangeTime() != 0 || n.HasColdProps() {
			t.Fatal("unexpected non-zero state")
		}
		_, _, _ = n.FirstNextTok()
	})
	if readAllocs > 0 {
		t.Fatalf("expected 0 allocs on ChangeTime/HasColdProps/FirstNextTok, got %v", readAllocs)
	}

	doc := &fakeJSDocHolder{comment: "/** @nosideeffects */"}
	n.SetJSDocObj(doc)
	n.SetPropList(prop)
	vals := []int32{1, 42, -1, -999999, math.MaxInt32, math.MinInt32, 0, 777}
	for _, want := range vals {
		n.SetChangeTime(want)
		if got := n.ChangeTime(); got != want {
			t.Fatalf("after SetChangeTime(%d): ChangeTime() = %d", want, got)
		}
		// Mutating KindFlags must preserve ChangeTime, JSDoc, and ColdProps!
		n.SetKindFlags(3)
		if n.KindFlags() != 3 || n.ChangeTime() != want || n.JSDocObj() != doc || !n.HasColdProps() || n.PropList() != prop {
			t.Fatalf("KindFlags/ChangeTime/JSDoc/PropList isolation broken for want=%d: kind=%d ct=%d",
				want, n.KindFlags(), n.ChangeTime())
		}
	}
	n.SetChangeTime(0)
	if n.ChangeTime() != 0 || n.Rec().KindFlags&kindFlagHasChangeTime != 0 {
		t.Fatalf("SetChangeTime(0) did not clear kindFlagHasChangeTime: ct=%d flags=0x%x",
			n.ChangeTime(), n.Rec().KindFlags)
	}

	// 3. ClonePropsFrom & NodeAllocator.SetTemplate with ChangeTime.
	src := AllocNodeToken(2, 50)
	src.SetSourceFile(19)
	src.SetPropFlags(0x1234)
	src.SetChangeTime(8888)
	dst := AllocNodeToken(4, 51)
	if !dst.ClonePropsFrom(src) {
		t.Fatalf("ClonePropsFrom returned false on clean node")
	}
	if dst.ChangeTime() != 8888 || dst.SourceFile() != 19 || dst.PropFlags() != 0x1234 || dst.KindFlags() != 4 {
		t.Fatalf("ClonePropsFrom did not copy ChangeTime/props: ct=%d sf=%d pf=0x%x kind=%d",
			dst.ChangeTime(), dst.SourceFile(), dst.PropFlags(), dst.KindFlags())
	}
	// Target already having ChangeTime must reject ClonePropsFrom!
	dirtyCT := AllocNodeToken(1, 30)
	dirtyCT.SetChangeTime(-5)
	if dirtyCT.ClonePropsFrom(src) || dirtyCT.ClonePropsFrom(NilNode) {
		t.Fatalf("ClonePropsFrom must return false when target already has ChangeTime")
	}

	alloc.InitLocal()
	alloc.SetTemplate(src)
	fromTmpl := alloc.AllocNode1(60, src)
	alloc.Flush()
	if fromTmpl.ChangeTime() != 8888 || fromTmpl.SourceFile() != 19 {
		t.Fatalf("NodeAllocator with ChangeTime template failed: ct=%d sf=%d",
			fromTmpl.ChangeTime(), fromTmpl.SourceFile())
	}

	// 4. Store switching, ResetNodeStore, and static template node ChangeTime.
	tmpl := NewTemplateNode(80)
	tmpl.SetChangeTime(314159)
	defer tmpl.SetChangeTime(0)

	SetCurrentNodeStore(s2)
	if src.ChangeTime() != 0 {
		t.Fatalf("unregistered s1 ChangeTime visible while s2 active: %d", src.ChangeTime())
	}
	if tmpl.ChangeTime() != 314159 {
		t.Fatalf("static template ChangeTime lost across store switch: %d", tmpl.ChangeTime())
	}
	s2Node := AllocNodeToken(1, 44)
	s2Node.SetChangeTime(-271828)
	if s2Node.ChangeTime() != -271828 {
		t.Fatalf("s2 ChangeTime = %d, want -271828", s2Node.ChangeTime())
	}

	ResetNodeStore()
	if src.ChangeTime() != 0 || s2Node.ChangeTime() != 0 || tmpl.ChangeTime() != 314159 {
		t.Fatalf("after ResetNodeStore: s1=%d s2=%d tmpl=%d",
			src.ChangeTime(), s2Node.ChangeTime(), tmpl.ChangeTime())
	}

	SetCurrentNodeStore(s1)
	if src.ChangeTime() != 8888 || dst.ChangeTime() != 8888 {
		t.Fatalf("restored s1 ChangeTime mismatch: src=%d dst=%d", src.ChangeTime(), dst.ChangeTime())
	}

	// 5. Concurrent ChangeTime reads/writes across shared 64-node pages and >2 blocks.
	const sharedCount = 2048
	shared := make([]Node, sharedCount)
	for i := range shared {
		shared[i] = s1.Alloc()
	}
	const workers = 6
	var wg sync.WaitGroup
	wg.Add(workers)
	for g := range workers {
		g := g
		go func() {
			defer wg.Done()
			for i := g; i < sharedCount; i += workers {
				v := int32(i*17 - 5000)
				if v == 0 {
					v = 1
				}
				shared[i].SetChangeTime(v)
				if got := shared[i].ChangeTime(); got != v {
					t.Errorf("concurrent ChangeTime[%d] = %d, want %d", i, got, v)
				}
				if i%9 == 0 {
					shared[i].SetChangeTime(0)
				}
			}
		}()
	}
	wg.Wait()
	for i, nd := range shared {
		want := int32(i*17 - 5000)
		if want == 0 {
			want = 1
		}
		if i%9 == 0 {
			want = 0
		}
		if got := nd.ChangeTime(); got != want {
			t.Fatalf("post-race shared[%d].ChangeTime() = %d, want %d", i, got, want)
		}
	}
}

type fakeColorHolder struct {
	name     string
	tcHandle atomic.Uint32
}

func (f *fakeColorHolder) IsJavaNil_() bool               { return f == nil }
func (f *fakeColorHolder) TCHandlePtr_() *atomic.Uint32   { return &f.tcHandle }
func (f *fakeColorHolder) Implements_colors_Color()       {}

func TestWave20ATypeColorFastPathsAndAllocatorInternCache(t *testing.T) {
	assertRuntimeError := func(name string, fn func()) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("%s: expected runtime.Error panic on Node(0), got nil", name)
			}
			if _, ok := r.(runtime.Error); !ok {
				t.Fatalf("%s: expected runtime.Error on Node(0), got %T (%v)", name, r, r)
			}
		}()
		fn()
	}

	assertRuntimeError("Node(0).JSTypeObj()", func() { _ = Node(0).JSTypeObj() })
	assertRuntimeError("Node(0).ColorObj()", func() { _ = Node(0).ColorObj() })
	assertRuntimeError("Node(0).SetJSTypeObj(nil)", func() { _ = Node(0).SetJSTypeObj(nil) })
	assertRuntimeError("Node(0).SetColorObj(nil)", func() { _ = Node(0).SetColorObj(nil) })

	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	jsType1 := &fakeTCHolder{name: "number"}
	jsType2 := &fakeTCHolder{name: "string"}
	color1 := &fakeColorHolder{name: "StandardColors.NUMBER"}
	color2 := &fakeColorHolder{name: "StandardColors.STRING"}
	var typedNilType *fakeTCHolder
	var typedNilColor *fakeColorHolder

	n := AllocNodeToken(4, 38)
	if n.JSTypeObj() != nil || n.ColorObj() != nil {
		t.Fatalf("fresh node should have nil JSTypeObj and ColorObj")
	}

	// SetJSTypeObj sets JSType, rejects SetColorObj, and clears cleanly via nil or typed-nil.
	if !n.SetJSTypeObj(jsType1) {
		t.Fatalf("SetJSTypeObj(jsType1) failed on unset node")
	}
	if n.JSTypeObj() != jsType1 || n.ColorObj() != nil || n.TypeOrColor() != jsType1 {
		t.Fatalf("after SetJSTypeObj: jstype=%#v color=%#v tc=%#v", n.JSTypeObj(), n.ColorObj(), n.TypeOrColor())
	}
	if n.SetColorObj(color1) || n.SetColorObj(nil) {
		t.Fatalf("SetColorObj should return false when node holds a JSType")
	}
	if !n.SetJSTypeObj(jsType2) || n.JSTypeObj() != jsType2 {
		t.Fatalf("SetJSTypeObj overwrite failed")
	}
	if !n.SetJSTypeObj(any(typedNilType)) || n.JSTypeObj() != nil || n.TC() != 0 {
		t.Fatalf("SetJSTypeObj(typedNil) failed to clear")
	}

	// Once cleared, SetColorObj succeeds and rejects SetJSTypeObj.
	if !n.SetColorObj(color1) {
		t.Fatalf("SetColorObj(color1) failed after JSType cleared")
	}
	if n.ColorObj() != color1 || n.JSTypeObj() != nil || n.TypeOrColor() != color1 || n.KindFlags() != 4 {
		t.Fatalf("after SetColorObj: color=%#v jstype=%#v kind=%d", n.ColorObj(), n.JSTypeObj(), n.KindFlags())
	}
	if n.SetJSTypeObj(jsType1) || n.SetJSTypeObj(nil) {
		t.Fatalf("SetJSTypeObj should return false when node holds a Color")
	}
	// SetKindFlags must preserve kindFlagTypeIsColor!
	n.SetKindFlags(3)
	if n.KindFlags() != 3 || n.ColorObj() != color1 || n.JSTypeObj() != nil {
		t.Fatalf("SetKindFlags clobbered kindFlagTypeIsColor: kind=%d color=%#v", n.KindFlags(), n.ColorObj())
	}
	if !n.SetColorObj(color2) || n.ColorObj() != color2 {
		t.Fatalf("SetColorObj overwrite failed")
	}
	if !n.SetColorObj(any(typedNilColor)) || n.ColorObj() != nil || n.TC() != 0 {
		t.Fatalf("SetColorObj(typedNil) failed to clear")
	}
	if !n.SetJSTypeObj(jsType1) || n.JSTypeObj() != jsType1 {
		t.Fatalf("SetJSTypeObj after clearing Color failed")
	}

	// SetTypeOrColor must also maintain kindFlagTypeIsColor for CopyTypeFrom / ConvertTypesToColors.
	n.SetTypeOrColor(color1)
	if n.ColorObj() != color1 || n.JSTypeObj() != nil {
		t.Fatalf("SetTypeOrColor(color1) did not set kindFlagTypeIsColor")
	}
	n.SetTypeOrColor(jsType2)
	if n.JSTypeObj() != jsType2 || n.ColorObj() != nil {
		t.Fatalf("SetTypeOrColor(jsType2) did not clear kindFlagTypeIsColor")
	}
	n.SetTypeOrColor(color2)
	n.SetTypeOrColor(nil)
	if n.ColorObj() != nil || n.JSTypeObj() != nil || !n.SetJSTypeObj(jsType1) {
		t.Fatalf("SetTypeOrColor(nil) did not clear kindFlagTypeIsColor")
	}

	// Zero-allocation verification on warm fast paths.
	allocs := testing.AllocsPerRun(100, func() {
		_ = n.SetJSTypeObj(jsType1)
		_ = n.JSTypeObj()
		_ = n.ColorObj()
		_ = n.SetJSTypeObj(nil)
		_ = n.SetColorObj(color1)
		_ = n.ColorObj()
		_ = n.JSTypeObj()
		_ = n.SetColorObj(nil)
	})
	if allocs > 0 {
		t.Fatalf("expected 0 allocs on JSTypeObj/ColorObj fast paths, got %v", allocs)
	}

	// NodeAllocator.Intern L1 direct-mapped cache parity, collision handling, and Flush cleanup.
	var alloc NodeAllocator
	alloc.InitLocal()
	if got := alloc.Intern(""); got != EmptyStringID {
		t.Fatalf("alloc.Intern(\"\") = %d, want %d", got, EmptyStringID)
	}
	words := make([]string, 600) // > 256 slots to force L1 direct-mapped collisions & evictions
	wantIDs := make([]StringID, len(words))
	for i := range words {
		words[i] = fmt.Sprintf("ident_w20a_%d_prop", i)
		wantIDs[i] = Intern(words[i])
	}
	for pass := range 3 {
		for i, w := range words {
			if got := alloc.Intern(w); got != wantIDs[i] {
				t.Fatalf("pass %d word[%d] (%q): alloc.Intern=%d, want %d", pass, i, w, got, wantIDs[i])
			}
		}
	}
	hitWord := words[len(words)-1]
	internAllocs := testing.AllocsPerRun(100, func() {
		if alloc.Intern(hitWord) != wantIDs[len(words)-1] {
			t.Fatal("unexpected L1 intern mismatch")
		}
	})
	if internAllocs > 0 {
		t.Fatalf("expected 0 allocs on NodeAllocator.Intern L1 hit, got %v", internAllocs)
	}
	if sz := reflect.TypeFor[NodeAllocator]().Size(); sz > 96 {
		t.Fatalf("expected compact NodeAllocator <= 96 bytes (pooled internCache), got %d", sz)
	}
	usedCache := alloc.internCache
	if usedCache == nil {
		t.Fatalf("expected non-nil pooled internCache after Intern calls")
	}
	alloc.Flush()
	if alloc.internDirty || alloc.internCache != nil {
		t.Fatalf("expected Flush() to reset internDirty and return internCache to pool")
	}
	for slot, e := range usedCache {
		if e != 0 {
			t.Fatalf("Flush() left stale entry at slot %d: %v", slot, e)
		}
	}

	// Verify lazy Intern() without InitLocal() (e.g. NewNodeAllocator or zero-value NodeAllocator).
	directAlloc := NewNodeAllocator(s)
	if got := directAlloc.Intern(hitWord); got != wantIDs[len(words)-1] {
		t.Fatalf("NewNodeAllocator.Intern mismatch: got %d, want %d", got, wantIDs[len(words)-1])
	}
	directAlloc.Flush()
}

func TestWave21APlainNodeInlineLengthAndLargeTCHandles(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	// 1. Lengths 0..65534 store directly in NodeData.Length16 without allocating s.lengthOverflows!
	var alloc NodeAllocator
	alloc.InitLocal()
	inlineLengths := []int32{0, 1, 42, 65534}
	plainNodes := make([]Node, len(inlineLengths))
	for i, ln := range inlineLengths {
		n := alloc.AllocNode1(uint8(10+i), 0)
		n.SetLength(ln)
		plainNodes[i] = n
	}
	alloc.Flush()

	if s.auxLen() != 0 {
		t.Fatalf("inline SetLength (< 65535) allocated s.aux entry eagerly: got %d entries, want 0", s.auxLen())
	}
	for i, n := range plainNodes {
		if got := n.Length(); got != inlineLengths[i] {
			t.Fatalf("plainNodes[%d].Length() = %d, want %d", i, got, inlineLengths[i])
		}
	}

	// Overflow lengths (>= 65535 or < 0) round-trip via lengthOverflows PagedColumn.
	overflowLengths := []int32{65535, 100000, -1, -999999, math.MinInt32, math.MaxInt32}
	for i, ln := range overflowLengths {
		n := s.Alloc()
		n.SetLength(ln)
		if got := n.Length(); got != ln {
			t.Fatalf("overflow[%d].Length() = %d, want %d", i, got, ln)
		}
	}

	// 2. StringNode (KindFlags == 3) and NumberNode (KindFlags == 1) isolate Payload and Length cleanly.
	alloc.InitLocal()
	sid := alloc.Intern("wave21a_identifier")
	strNode := alloc.AllocStringNode(32, sid, 0)
	numNode := alloc.AllocNumberNode(33, 2.718281828, 0)
	alloc.Flush()

	if strNode.Length() != 0 || numNode.Length() != 0 {
		t.Fatalf("fresh StringNode/NumberNode should have Length() == 0, got %d / %d", strNode.Length(), numNode.Length())
	}
	strNode.SetLength(18)
	numNode.SetLength(11)
	if strNode.Length() != 18 || String(StringID(strNode.Payload())) != "wave21a_identifier" {
		t.Fatalf("StringNode Payload/Length clobbered: len=%d str=%q", strNode.Length(), String(StringID(strNode.Payload())))
	}
	if numNode.Length() != 11 || numNode.NumberVal() != 2.718281828 {
		t.Fatalf("NumberNode Payload/Length clobbered: len=%d num=%v", numNode.Length(), numNode.NumberVal())
	}

	// 3. Raw SetPayload, SetNumberVal, and SetKindFlags migrate any existing inline plain-node Length cleanly.
	nMig1 := s.Alloc()
	nMig1.SetLength(77777)
	nMig1.SetPayload(0xCAFEBABE)
	if nMig1.Payload() != 0xCAFEBABE || nMig1.Length() != 77777 {
		t.Fatalf("SetPayload after SetLength clobbered: payload=0x%x len=%d", nMig1.Payload(), nMig1.Length())
	}
	nMig1.SetLength(88888)
	if nMig1.Payload() != 0xCAFEBABE || nMig1.Length() != 88888 {
		t.Fatalf("SetLength after SetPayload clobbered: payload=0x%x len=%d", nMig1.Payload(), nMig1.Length())
	}

	nMig2 := s.Alloc()
	nMig2.SetLength(54321)
	nMig2.SetKindFlags(3)
	if nMig2.KindFlags() != 3 || nMig2.Payload() != 0 || nMig2.Length() != 54321 {
		t.Fatalf("SetKindFlags(3) after SetLength failed migration: kind=%d payload=%d len=%d",
			nMig2.KindFlags(), nMig2.Payload(), nMig2.Length())
	}

	nMig3 := s.Alloc()
	nMig3.SetLength(13579)
	nMig3.SetNumberVal(-12.5)
	if nMig3.NumberVal() != -12.5 || nMig3.Length() != 13579 {
		t.Fatalf("SetNumberVal after SetLength failed migration: num=%v len=%d", nMig3.NumberVal(), nMig3.Length())
	}

	// 4. JSTypeObj & ColorObj across >4096 handles (multi-block AnyTable) with 0 allocs.
	for tcHandles.count.Load() < 4100 {
		_ = tcHandles.Add(&tcHandlesMu, &fakeTCHolder{name: "pad"})
	}
	highType := &fakeTCHolder{name: "high_jstype"}
	highColor := &fakeColorHolder{name: "high_color"}
	nt := AllocNodeToken(0, 32)
	nc := AllocNodeToken(0, 32)
	if !nt.SetJSTypeObj(highType) || !nc.SetColorObj(highColor) {
		t.Fatal("failed setting high JSType / Color")
	}
	if highType.tcHandle.Load() < 4096 || highColor.tcHandle.Load() < 4096 {
		t.Fatalf("expected handles >= 4096, got %d / %d", highType.tcHandle.Load(), highColor.tcHandle.Load())
	}
	if nt.JSTypeObj() != highType || nt.ColorObj() != nil || nc.ColorObj() != highColor || nc.JSTypeObj() != nil {
		t.Fatal("high-handle JSTypeObj/ColorObj mismatch")
	}
}

func TestWave22AStoreResetHooksBlockOfAndCachedAllocators(t *testing.T) {
	assertRuntimeError := func(name string, fn func()) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("%s: expected runtime.Error panic, got nil", name)
			}
			if _, ok := r.(runtime.Error); !ok {
				t.Fatalf("%s: expected runtime.Error, got %T (%v)", name, r, r)
			}
		}()
		fn()
	}

	// 1. RegisterStoreResetHook lifecycle, ordering, non-firing cases, and -race concurrency.
	prevOrig := CurrentNodeStore()
	defer SetCurrentNodeStore(prevOrig)
	ResetNodeStore()

	var hookActive atomic.Bool
	hookActive.Store(true)
	defer hookActive.Store(false)

	var hookCalls1, hookCalls2, nestedHookCalls atomic.Int64
	var orderMu sync.Mutex
	var callOrder []int
	var registeredNested atomic.Bool

	RegisterStoreResetHook(nil) // nil hook must be safely ignored
	RegisterStoreResetHook(func() {
		if !hookActive.Load() {
			return
		}
		hookCalls1.Add(1)
		// Calling CurrentNodeStore and RegisterStoreResetHook inside a hook must not deadlock!
		_ = CurrentNodeStore()
		if registeredNested.CompareAndSwap(false, true) {
			RegisterStoreResetHook(func() {
				if hookActive.Load() {
					nestedHookCalls.Add(1)
				}
			})
		}
		orderMu.Lock()
		callOrder = append(callOrder, 1)
		orderMu.Unlock()
	})
	RegisterStoreResetHook(func() {
		if !hookActive.Load() {
			return
		}
		hookCalls2.Add(1)
		orderMu.Lock()
		callOrder = append(callOrder, 2)
		orderMu.Unlock()
	})

	// Idempotent ResetNodeStore while staticStore is active must not fire hooks.
	ResetNodeStore()
	if hookCalls1.Load() != 0 || hookCalls2.Load() != 0 {
		t.Fatalf("hooks fired on idle ResetNodeStore: %d, %d", hookCalls1.Load(), hookCalls2.Load())
	}

	s1 := NewNodeStore()
	s2 := NewNodeStore()

	// Activating s1 from staticStore must not fire reset hooks.
	SetCurrentNodeStore(s1)
	if hookCalls1.Load() != 0 || hookCalls2.Load() != 0 {
		t.Fatalf("hooks fired when activating s1 from staticStore: %d, %d", hookCalls1.Load(), hookCalls2.Load())
	}
	// Idempotent SetCurrentNodeStore(s1) must not fire reset hooks.
	SetCurrentNodeStore(s1)
	if hookCalls1.Load() != 0 || hookCalls2.Load() != 0 {
		t.Fatalf("hooks fired on idempotent SetCurrentNodeStore(s1): %d, %d", hookCalls1.Load(), hookCalls2.Load())
	}

	// Switching s1 -> s2 unregisters s1 and must fire hooks 1 and 2 in order.
	SetCurrentNodeStore(s2)
	if hookCalls1.Load() != 1 || hookCalls2.Load() != 1 || nestedHookCalls.Load() != 0 {
		t.Fatalf("after s1->s2: got (%d, %d, %d), want (1, 1, 0)",
			hookCalls1.Load(), hookCalls2.Load(), nestedHookCalls.Load())
	}

	// ResetNodeStore unregisters s2 and fires hooks 1, 2, and the newly registered nested hook.
	ResetNodeStore()
	if hookCalls1.Load() != 2 || hookCalls2.Load() != 2 || nestedHookCalls.Load() != 1 {
		t.Fatalf("after s2->nil: got (%d, %d, %d), want (2, 2, 1)",
			hookCalls1.Load(), hookCalls2.Load(), nestedHookCalls.Load())
	}
	orderMu.Lock()
	if len(callOrder) != 4 || callOrder[0] != 1 || callOrder[1] != 2 || callOrder[2] != 1 || callOrder[3] != 2 {
		t.Fatalf("unexpected hook order: %v", callOrder)
	}
	orderMu.Unlock()

	// Concurrent RegisterStoreResetHook while switching stores under -race.
	var hookRaceWG sync.WaitGroup
	hookRaceWG.Add(2)
	go func() {
		defer hookRaceWG.Done()
		for range 20 {
			RegisterStoreResetHook(func() {})
		}
	}()
	go func() {
		defer hookRaceWG.Done()
		for range 20 {
			SetCurrentNodeStore(s1)
			ResetNodeStore()
		}
	}()
	hookRaceWG.Wait()

	// 2. BlockShift, BlockOf, and (*NodeBlock).FirstNextTok (including NPE panic & cross-block traversal).
	if BlockShift != 14 || (1<<BlockShift) != blockLen {
		t.Fatalf("unexpected BlockShift=%d", BlockShift)
	}
	assertRuntimeError("BlockOf(NilNode)", func() {
		_ = BlockOf(NilNode)
	})
	assertRuntimeError("BlockOf(unmapped)", func() {
		_ = BlockOf(Node(0x3FFF_FFFF))
	})
	assertRuntimeError("(*NodeBlock)(nil).FirstNextTok(0)", func() {
		var nilBlk *NodeBlock
		_, _, _ = nilBlk.FirstNextTok(0)
	})

	s := NewNodeStore()
	SetCurrentNodeStore(s)

	// Allocate across a block boundary so siblings span two distinct NodeBlocks.
	root := AllocNodeToken(0, 100)
	var prevChild Node
	const crossBlockChildren = blockLen + 300
	expectedChildren := make([]Node, crossBlockChildren)
	for i := range crossBlockChildren {
		tok := uint8((i % 200) + 1)
		c := AllocNodeToken(0, tok)
		if i%5 == 0 {
			gc := AllocNodeToken(0, uint8((i%50)+201))
			c.SetFirstChild(gc)
			gc.SetParent(c)
			gc.SetPrevious(gc)
		}
		c.SetParent(root)
		if prevChild == 0 {
			root.SetFirstChild(c)
		} else {
			prevChild.SetNext(c)
			c.SetPrevious(prevChild)
		}
		prevChild = c
		expectedChildren[i] = c
	}
	root.FirstChild().SetPrevious(prevChild)

	rootBlk := BlockOf(root)
	curBlkIdx := uint32(root) >> BlockShift
	rf, rn, rtok := rootBlk.FirstNextTok(root)
	if rf != expectedChildren[0] || rn != 0 || rtok != 100 {
		t.Fatalf("root FirstNextTok = (%d, %d, %d), want (%d, 0, 100)", rf, rn, rtok, expectedChildren[0])
	}

	curBlk := rootBlk
	blockSwitches := 0
	idx := 0
	for c := rf; c != 0; idx++ {
		if cIdx := uint32(c) >> BlockShift; cIdx != curBlkIdx {
			curBlkIdx = cIdx
			curBlk = BlockOf(c)
			blockSwitches++
		}
		cf, cn, ctok := curBlk.FirstNextTok(c)
		wf, wn, wtok := c.FirstNextTok()
		if cf != wf || cn != wn || ctok != wtok || c != expectedChildren[idx] {
			t.Fatalf("child[%d] Block.FirstNextTok=(%d,%d,%d), want (%d,%d,%d)",
				idx, cf, cn, ctok, wf, wn, wtok)
		}
		c = cn
	}
	if idx != crossBlockChildren || blockSwitches == 0 {
		t.Fatalf("traversed %d children with %d block switches, want %d and >=1", idx, blockSwitches, crossBlockChildren)
	}

	// Verify 0 allocs and -race cleanliness when node at block offset 0 (root) is concurrently mutated
	// while another goroutine calls BlockOf(otherNodeInSameBlock) and blk.FirstNextTok.
	otherInFirstBlk := expectedChildren[10]
	travAllocs := testing.AllocsPerRun(100, func() {
		b := BlockOf(otherInFirstBlk)
		_, _, _ = b.FirstNextTok(otherInFirstBlk)
	})
	if travAllocs > 0 {
		t.Fatalf("expected 0 allocs for BlockOf + FirstNextTok, got %v", travAllocs)
	}

	var blkRaceWG sync.WaitGroup
	blkRaceWG.Add(2)
	go func() {
		defer blkRaceWG.Done()
		for i := range 2000 {
			// root is at slot 0 of s's first NodeBlock (root & blockMask == 0).
			root.SetToken(uint8(i & 0xFF))
		}
	}()
	go func() {
		defer blkRaceWG.Done()
		for range 2000 {
			b := BlockOf(otherInFirstBlk)
			f, nx, tk := b.FirstNextTok(otherInFirstBlk)
			if nx != expectedChildren[11] || tk != expectedChildren[10].Token() || f != otherInFirstBlk.FirstChild() {
				t.Errorf("concurrent FirstNextTok mismatch: (%d, %d, %d)", f, nx, tk)
				return
			}
		}
	}()
	blkRaceWG.Wait()

	// 3. AllocNode1Cached & AllocStringNodeCached (fast path, slow templates, block boundaries, 0 allocs, -race).
	fastTmpl1 := AllocNodeToken(0, 10)
	fastTmpl1.SetSourceFile(17)
	fastTmpl1.SetPropFlags(0x13579BDF)

	fastTmpl2 := AllocNodeToken(0, 11)
	fastTmpl2.SetSourceFile(42)
	fastTmpl2.SetPropFlags(0x2468ACE0)

	overflowSFTmpl := AllocNodeToken(0, 12)
	overflowSFTmpl.SetSourceFile(90000) // >= 0xFFFF -> a.tmplSlow == true -> a.end == 0
	overflowSFTmpl.SetPropFlags(0x55AA55AA)

	coldTmpl := AllocNodeToken(0, 13)
	coldTmpl.SetSourceFile(9)
	coldTmpl.SetPropFlags(0x0F0F0F0F)
	coldDoc := &fakeJSDocHolder{comment: "/** @const */"}
	coldProp := &fakeColdProp{propType: 7, val: "wave22a-cold"}
	coldTmpl.SetJSDocObj(coldDoc)
	coldTmpl.SetPropList(coldProp)
	coldTmpl.SetChangeTime(4321)

	templates := []Node{0, fastTmpl1, fastTmpl2, overflowSFTmpl, coldTmpl}

	var alloc NodeAllocator
	alloc.InitLocal()
	sidHello := alloc.Intern("hello_wave22a")
	sidWorld := alloc.Intern("world_wave22a")

	const totalCached = blockLen + 1500
	for i := range totalCached {
		tmpl := templates[i%len(templates)]
		alloc.SetTemplate(tmpl)
		tok := uint8((i % 120) + 1)
		var n Node
		isStr := i%2 == 1
		wantSID := sidHello
		if i%4 == 3 {
			wantSID = sidWorld
		}
		if isStr {
			n = alloc.AllocStringNodeCached(tok, wantSID)
		} else {
			n = alloc.AllocNode1Cached(tok)
		}

		if n.Token() != tok || n.Lineno() != -1 || n.Charno() != -1 || n.FirstChild() != 0 || n.Next() != 0 || n.Parent() != 0 || n.Previous() != 0 {
			t.Fatalf("iter %d: bad node header tok=%d line=%d char=%d", i, n.Token(), n.Lineno(), n.Charno())
		}
		if isStr {
			if n.KindFlags() != 3 || StringID(n.Payload()) != wantSID || n.Length() != 0 {
				t.Fatalf("iter %d: bad StringNode kind=%d payload=%d len=%d", i, n.KindFlags(), n.Payload(), n.Length())
			}
		} else {
			wantPay := uint32(0)
			if tmpl == coldTmpl {
				wantPay = 4321 // plain leaf node stores cloned ChangeTime inline in Payload
			}
			if n.KindFlags() != 0 || n.Payload() != wantPay || n.Length() != 0 {
				t.Fatalf("iter %d: bad leaf node kind=%d payload=%d (want %d) len=%d", i, n.KindFlags(), n.Payload(), wantPay, n.Length())
			}
		}
		if tmpl == 0 {
			if n.SourceFile() != 0 || n.PropFlags() != 0 || n.HasColdProps() || n.JSDocObj() != nil || n.ChangeTime() != 0 {
				t.Fatalf("iter %d (nil tmpl): unexpected props sf=%d pf=0x%x", i, n.SourceFile(), n.PropFlags())
			}
		} else {
			if n.SourceFile() != tmpl.SourceFile() || n.PropFlags() != tmpl.PropFlags() || n.ChangeTime() != tmpl.ChangeTime() || n.JSDocObj() != tmpl.JSDocObj() {
				t.Fatalf("iter %d: template mismatch sf=%d/%d pf=0x%x/0x%x ct=%d/%d",
					i, n.SourceFile(), tmpl.SourceFile(), n.PropFlags(), tmpl.PropFlags(), n.ChangeTime(), tmpl.ChangeTime())
			}
			if tmpl == coldTmpl && !n.HasColdProps() {
				t.Fatalf("iter %d: expected cold props from coldTmpl", i)
			}
		}
	}
	alloc.Flush()

	// Warm fast-path allocation check: stack-allocated NodeAllocator + AllocNode1Cached + AllocStringNodeCached = 0 heap allocs.
	var warmAlloc NodeAllocator
	warmAlloc.InitLocal()
	warmAlloc.SetTemplate(fastTmpl1)
	_ = warmAlloc.AllocNode1Cached(1) // prime 256-node batch inside already-allocated block
	cachedAllocs := testing.AllocsPerRun(50, func() {
		n1 := warmAlloc.AllocNode1Cached(40)
		n2 := warmAlloc.AllocStringNodeCached(38, sidHello)
		if n1 == 0 || n2 == 0 {
			t.Fatal("nil node")
		}
	})
	warmAlloc.Flush()
	if cachedAllocs > 0 {
		t.Fatalf("expected 0 heap allocs on warm AllocNode1Cached/AllocStringNodeCached, got %v", cachedAllocs)
	}

	// Concurrent multi-worker NodeAllocator stress with AllocNode1Cached & AllocStringNodeCached.
	const workers = 6
	var allocWG sync.WaitGroup
	allocWG.Add(workers)
	for w := range workers {
		w := w
		go func() {
			defer allocWG.Done()
			var la NodeAllocator
			la.InitLocal()
			defer la.Flush()
			sid := la.Intern(fmt.Sprintf("worker_%d_ident", w))
			for j := range 1200 {
				if j%300 == 0 {
					la.SetTemplate(templates[(w+j)%len(templates)])
				}
				n1 := la.AllocNode1Cached(uint8((j % 90) + 1))
				n2 := la.AllocStringNodeCached(38, sid)
				if n1.Token() != uint8((j%90)+1) || StringID(n2.Payload()) != sid || n2.KindFlags() != 3 {
					t.Errorf("worker %d iter %d mismatch", w, j)
					return
				}
			}
		}()
	}
	allocWG.Wait()
}

func TestWave23AOrigNameIsStringAndSetStringPayload(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	if NilNode.OrigNameID() != 0 {
		t.Fatalf("expected NilNode.OrigNameID() == 0, got %d", NilNode.OrigNameID())
	}

	idFoo := Intern("wave23a_foo")
	idBar := Intern("wave23a_bar")
	idRenamed := Intern("wave23a_a")

	var alloc NodeAllocator
	alloc.InitLocal()
	defer alloc.Flush()

	// 1. StringNode SetOrigNameID(sameID) -> 0 origNameField / PagedColumn allocations, 0 heap allocs.
	const normCount = 2000 // spans > 30 64-node pages
	normNodes := make([]Node, normCount)
	for i := range normCount {
		sid := idFoo
		if i%3 == 1 {
			sid = idBar
		} else if i%3 == 2 {
			sid = EmptyStringID
		}
		normNodes[i] = alloc.AllocStringNodeCached(38, sid)
	}
	alloc.Flush()

	if len(s.origNames) != 0 {
		t.Fatalf("expected 0 origNames columns before Normalize, got %d", len(s.origNames))
	}

	normAllocs := testing.AllocsPerRun(1, func() {
		for i, n := range normNodes {
			sid := idFoo
			if i%3 == 1 {
				sid = idBar
			} else if i%3 == 2 {
				sid = EmptyStringID
			}
			n.SetOrigNameID(sid)
			if n.OrigNameID() != sid || n.KindFlags() != 3 {
				t.Fatalf("node %d OrigNameID=%d (want %d), KindFlags=%d", i, n.OrigNameID(), sid, n.KindFlags())
			}
		}
	})
	if normAllocs > 0 {
		t.Fatalf("expected 0 heap allocs when SetOrigNameID == StringNode Payload, got %v", normAllocs)
	}
	if len(s.origNames) != 0 {
		t.Fatalf("SetOrigNameID(sameID) allocated %d origNames PagedColumns, want 0", len(s.origNames))
	}

	// NullString StringNode (Payload == 0) + SetOrigNameID(0) must NOT set kindFlagInlineOrigName or allocate pages.
	nullStrNode := AllocNodeToken(3, 38)
	nullStrNode.SetOrigNameID(0)
	if nullStrNode.Rec().KindFlags&kindFlagInlineOrigName != 0 || nullStrNode.OrigNameID() != 0 || len(s.origNames) != 0 {
		t.Fatalf("nullStrNode SetOrigNameID(0) corrupted state: rawFlags=0x%x orig=%d cols=%d",
			nullStrNode.Rec().KindFlags, nullStrNode.OrigNameID(), len(s.origNames))
	}

	// Idempotent SetStringPayload(sameID) must NOT materialize origNameField!
	n0 := normNodes[0] // Payload == idFoo, kindFlagInlineOrigName set
	n0.SetStringPayload(idFoo)
	if len(s.origNames) != 0 || n0.Rec().KindFlags&kindFlagInlineOrigName == 0 || n0.OrigNameID() != idFoo {
		t.Fatalf("idempotent SetStringPayload(idFoo) unexpectedly materialized origNames: cols=%d", len(s.origNames))
	}

	// 2. Subsequent SetStringPayload(newID) (renaming "foo" -> "a") materializes "foo" in origNameField!
	n0.SetStringPayload(idRenamed)
	if StringID(n0.Payload()) != idRenamed {
		t.Fatalf("after SetStringPayload(idRenamed): Payload=%d, want %d", n0.Payload(), idRenamed)
	}
	if n0.OrigNameID() != idFoo {
		t.Fatalf("after SetStringPayload(idRenamed): OrigNameID=%d (%q), want %d (%q)",
			n0.OrigNameID(), String(n0.OrigNameID()), idFoo, String(idFoo))
	}
	if n0.KindFlags() != 3 || n0.Rec().KindFlags&kindFlagInlineOrigName != 0 {
		t.Fatalf("expected kindFlagInlineOrigName cleared after SetStringPayload(newID), rawFlags=0x%x", n0.Rec().KindFlags)
	}

	// SetStringPayload(0) on EmptyStringID StringNode with kindFlagInlineOrigName preserves EmptyStringID!
	nEmpty := normNodes[2] // Payload == EmptyStringID
	nEmpty.SetStringPayload(0)
	if nEmpty.Payload() != 0 || nEmpty.OrigNameID() != EmptyStringID {
		t.Fatalf("SetStringPayload(0) lost EmptyStringID originalName: payload=%d orig=%d", nEmpty.Payload(), nEmpty.OrigNameID())
	}

	// Guard checks on SetPayload, SetKindFlags, and SetNumberVal:
	nPay := normNodes[3] // Payload == idFoo
	nPay.SetPayload(uint32(idRenamed))
	if StringID(nPay.Payload()) != idRenamed || nPay.OrigNameID() != idFoo {
		t.Fatalf("SetPayload failed to materialize origName: payload=%d orig=%d", nPay.Payload(), nPay.OrigNameID())
	}

	nKind := normNodes[6] // Payload == idFoo
	nKind.SetKindFlags(4) // change away from StringNode (3 -> 4)
	nKind.SetPayload(999999)
	if nKind.KindFlags() != 4 || nKind.Payload() != 999999 || nKind.OrigNameID() != idFoo {
		t.Fatalf("SetKindFlags(4) failed to materialize origName: kind=%d payload=%d orig=%d",
			nKind.KindFlags(), nKind.Payload(), nKind.OrigNameID())
	}

	nKindSame := normNodes[9] // Payload == idFoo
	nKindSame.SetKindFlags(3) // staying StringNode preserves kindFlagInlineOrigName
	if nKindSame.Rec().KindFlags&kindFlagInlineOrigName == 0 || nKindSame.OrigNameID() != idFoo {
		t.Fatalf("SetKindFlags(3) unexpectedly cleared kindFlagInlineOrigName")
	}

	// Non-String/non-plain node (e.g. kind == 2) with Payload numerically equal to v must NOT set kindFlagInlineOrigName.
	nonStrSamePayload := AllocNodeToken(2, 30)
	nonStrSamePayload.SetPayload(uint32(idFoo))
	nonStrSamePayload.SetOrigNameID(idFoo)
	if nonStrSamePayload.Rec().KindFlags&kindFlagInlineOrigName != 0 {
		t.Fatalf("kind==2 node with Payload == idFoo must not set kindFlagInlineOrigName")
	}
	nonStrSamePayload.SetPayload(77)
	if nonStrSamePayload.Payload() != 77 || nonStrSamePayload.OrigNameID() != idFoo {
		t.Fatalf("kind==2 node SetPayload clobbered OrigNameID: payload=%d orig=%d", nonStrSamePayload.Payload(), nonStrSamePayload.OrigNameID())
	}

	// Plain node (kind == 0) stores Length in Length16 and OrigNameID inline in Payload independently!
	plainInlineLen := AllocNodeToken(0, 30)
	plainInlineLen.SetLength(1234)
	plainInlineLen.SetOrigNameID(idFoo)
	if plainInlineLen.Rec().KindFlags&kindFlagInlineOrigName == 0 {
		t.Fatalf("expected plain node without ChangeTime to store OrigNameID inline")
	}
	plainInlineLen.SetLength(77)
	if plainInlineLen.Length() != 77 || plainInlineLen.OrigNameID() != idFoo {
		t.Fatalf("plain node SetLength clobbered OrigNameID: len=%d orig=%d", plainInlineLen.Length(), plainInlineLen.OrigNameID())
	}

	// 3. SetOrigNameID(differentID) -> SetOrigNameID(sameAsPayload) -> SetOrigNameID(0) transitions.
	nTrans := normNodes[12] // Payload == idFoo, currently kindFlagInlineOrigName
	nTrans.SetOrigNameID(idBar)
	if nTrans.OrigNameID() != idBar || nTrans.Rec().KindFlags&kindFlagInlineOrigName != 0 || origNameField.get(nTrans) != uint32(idBar) {
		t.Fatalf("SetOrigNameID(differentID) mismatch: orig=%d rawFlags=0x%x col=%d",
			nTrans.OrigNameID(), nTrans.Rec().KindFlags, origNameField.get(nTrans))
	}
	// Transition back to sameAsPayload: must clear origNameField slot to 0 and set kindFlagInlineOrigName!
	nTrans.SetOrigNameID(idFoo)
	if nTrans.OrigNameID() != idFoo || nTrans.Rec().KindFlags&kindFlagInlineOrigName == 0 || origNameField.get(nTrans) != 0 {
		t.Fatalf("SetOrigNameID(sameAsPayload) did not clear origNameField & set flag: orig=%d rawFlags=0x%x col=%d",
			nTrans.OrigNameID(), nTrans.Rec().KindFlags, origNameField.get(nTrans))
	}
	// Transition to 0: must clear flag and return 0!
	nTrans.SetOrigNameID(0)
	if nTrans.OrigNameID() != 0 || nTrans.Rec().KindFlags&kindFlagInlineOrigName != 0 || origNameField.get(nTrans) != 0 {
		t.Fatalf("SetOrigNameID(0) did not clear originalName: orig=%d rawFlags=0x%x col=%d",
			nTrans.OrigNameID(), nTrans.Rec().KindFlags, origNameField.get(nTrans))
	}

	// 4. ClonePropsFrom + SetTemplate + cross-kind OrigNameID copying (StringNode same/diff & plain CALL Node).
	srcName := AllocNodeToken(3, 38)
	srcName.SetStringPayload(idFoo)
	srcName.SetOrigNameID(idFoo) // kindFlagInlineOrigName set on srcName!
	srcName.SetSourceFile(23)
	srcName.SetPropFlags(0x123456)
	doc := &fakeJSDocHolder{comment: "/** @type {string} */"}
	prop := &fakeColdProp{propType: 5, val: "wave23a-prop"}
	srcName.SetJSDocObj(doc)
	srcName.SetChangeTime(99)
	srcName.SetPropList(prop)

	// 4a. ClonePropsFrom onto StringNode with DIFFERENT string ("bar"): must NOT leak kindFlagInlineOrigName!
	dstDiffStr := AllocNodeToken(3, 38)
	dstDiffStr.SetStringPayload(idBar)
	if !dstDiffStr.ClonePropsFrom(srcName) {
		t.Fatal("ClonePropsFrom onto dstDiffStr failed")
	}
	if dstDiffStr.KindFlags() != 3 || dstDiffStr.OrigNameID() != 0 || StringID(dstDiffStr.Payload()) != idBar ||
		dstDiffStr.SourceFile() != 23 || dstDiffStr.PropFlags() != 0x123456 || dstDiffStr.JSDocObj() != doc ||
		dstDiffStr.ChangeTime() != 99 || dstDiffStr.PropList() != prop {
		t.Fatalf("ClonePropsFrom onto dstDiffStr mismatch: kind=%d orig=%d payload=%d",
			dstDiffStr.KindFlags(), dstDiffStr.OrigNameID(), dstDiffStr.Payload())
	}
	// Copying originalName (as in Srcref / CopyBaseNodeFields) onto dstDiffStr stores idFoo in origNameField:
	dstDiffStr.SetOrigNameID(srcName.OrigNameID())
	if dstDiffStr.OrigNameID() != idFoo || StringID(dstDiffStr.Payload()) != idBar || dstDiffStr.Rec().KindFlags&kindFlagInlineOrigName != 0 {
		t.Fatalf("copying OrigNameID onto dstDiffStr failed: orig=%d payload=%d", dstDiffStr.OrigNameID(), dstDiffStr.Payload())
	}

	// 4b. ClonePropsFrom onto StringNode that ALREADY has kindFlagInlineOrigName set (same string "foo"):
	dstSameStr := AllocNodeToken(3, 38)
	dstSameStr.SetStringPayload(idFoo)
	dstSameStr.SetOrigNameID(idFoo)
	if !dstSameStr.ClonePropsFrom(srcName) {
		t.Fatal("ClonePropsFrom onto dstSameStr with existing OrigNameID failed")
	}
	if dstSameStr.OrigNameID() != idFoo || dstSameStr.Rec().KindFlags&kindFlagInlineOrigName == 0 || dstSameStr.SourceFile() != 23 {
		t.Fatalf("ClonePropsFrom clobbered dstSameStr inline OrigNameID: orig=%d sf=%d", dstSameStr.OrigNameID(), dstSameStr.SourceFile())
	}

	// 4c. ClonePropsFrom + OrigNameID copy onto plain CALL Node (KindFlags == 0):
	callNode := AllocNodeToken(0, 45)
	callNode.SetLength(14)
	if !callNode.ClonePropsFrom(srcName) {
		t.Fatal("ClonePropsFrom onto plain callNode failed")
	}
	if callNode.KindFlags() != 0 || callNode.OrigNameID() != 0 || callNode.Length() != 14 || callNode.SourceFile() != 23 || callNode.ChangeTime() != 99 {
		t.Fatalf("ClonePropsFrom onto callNode corrupted plain node: kind=%d orig=%d len=%d sf=%d ct=%d",
			callNode.KindFlags(), callNode.OrigNameID(), callNode.Length(), callNode.SourceFile(), callNode.ChangeTime())
	}
	callNode.SetOrigNameID(srcName.OrigNameID())
	if callNode.KindFlags() != 0 || callNode.OrigNameID() != idFoo || callNode.Length() != 14 || callNode.ChangeTime() != 99 || callNode.Rec().KindFlags&kindFlagInlineOrigName != 0 {
		t.Fatalf("callNode.SetOrigNameID(srcName.OrigNameID()) corrupted callNode: kind=%d orig=%d len=%d ct=%d",
			callNode.KindFlags(), callNode.OrigNameID(), callNode.Length(), callNode.ChangeTime())
	}

	// 4d. SetTemplate on hot StringNode with kindFlagInlineOrigName stays on fast inline path (tmplCold == false).
	hotTmpl := AllocNodeToken(3, 38)
	hotTmpl.SetStringPayload(idFoo)
	hotTmpl.SetOrigNameID(idFoo)
	hotTmpl.SetSourceFile(31)
	hotTmpl.SetPropFlags(0x77)
	alloc.InitLocal()
	alloc.SetTemplate(hotTmpl)
	if alloc.tmplCold {
		t.Fatalf("hotTmpl with kindFlagInlineOrigName unexpectedly marked cold: cold=%v", alloc.tmplCold)
	}
	tLeaf := alloc.AllocNode1Cached(45)
	tStr := alloc.AllocStringNodeCached(38, idBar)
	alloc.Flush()
	if tLeaf.KindFlags() != 0 || tLeaf.OrigNameID() != 0 || tLeaf.SourceFile() != 31 || tLeaf.PropFlags() != 0x77 {
		t.Fatalf("AllocNode1Cached from hotTmpl leaked OrigNameID: kind=%d orig=%d sf=%d", tLeaf.KindFlags(), tLeaf.OrigNameID(), tLeaf.SourceFile())
	}
	if tStr.KindFlags() != 3 || tStr.OrigNameID() != 0 || StringID(tStr.Payload()) != idBar || tStr.SourceFile() != 31 {
		t.Fatalf("AllocStringNodeCached from hotTmpl leaked OrigNameID: kind=%d orig=%d payload=%d", tStr.KindFlags(), tStr.OrigNameID(), tStr.Payload())
	}

	// 5. Concurrent -race stress across goroutines: Normalize -> Rename -> Reset -> ClonePropsFrom across pages & blocks.
	const workers = 8
	const perWorker = 3000
	workerNodes := make([][]Node, workers)
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := range workers {
		w := w
		workerNodes[w] = make([]Node, perWorker)
		go func() {
			defer wg.Done()
			var la NodeAllocator
			la.InitLocal()
			defer la.Flush()
			la.SetTemplate(hotTmpl)
			origSID := la.Intern(fmt.Sprintf("wave23a_orig_%d", w))
			renSID := la.Intern(fmt.Sprintf("wave23a_ren_%d", w))
			altSID := la.Intern(fmt.Sprintf("wave23a_alt_%d", w))
			for i := range perWorker {
				var n Node
				switch i % 4 {
				case 0:
					// Normalized StringNode kept unrenamed (inline kindFlagOrigNameIsString).
					n = la.AllocStringNodeCached(38, origSID)
					n.SetOrigNameID(origSID)
				case 1:
					// Normalized StringNode then renamed via SetStringPayload!
					n = la.AllocStringNodeCached(38, origSID)
					n.SetOrigNameID(origSID)
					n.SetStringPayload(renSID)
				case 2:
					// Distinct origName -> reverted to Payload -> cleared to 0 -> set back to Payload!
					n = la.AllocStringNodeCached(38, origSID)
					n.SetOrigNameID(altSID)
					n.SetOrigNameID(origSID)
					if i%8 == 2 {
						n.SetOrigNameID(0)
					}
				default:
					// Plain CALL node with OrigNameID copied from normalized StringNode.
					n = la.AllocNode1Cached(45)
					n.SetLength(int32(i + 1))
					n.SetOrigNameID(origSID)
				}
				workerNodes[w][i] = n
			}
		}()
	}
	wg.Wait()

	// Concurrent read verification across all workers (mirrors parallel CodePrinter).
	wg.Add(workers)
	for w := range workers {
		w := w
		go func() {
			defer wg.Done()
			origSID := Intern(fmt.Sprintf("wave23a_orig_%d", w))
			renSID := Intern(fmt.Sprintf("wave23a_ren_%d", w))
			for i, n := range workerNodes[w] {
				switch i % 4 {
				case 0:
					if n.KindFlags() != 3 || StringID(n.Payload()) != origSID || n.OrigNameID() != origSID {
						t.Errorf("w=%d i=%d case 0 mismatch: kind=%d pay=%d orig=%d", w, i, n.KindFlags(), n.Payload(), n.OrigNameID())
						return
					}
				case 1:
					if n.KindFlags() != 3 || StringID(n.Payload()) != renSID || n.OrigNameID() != origSID {
						t.Errorf("w=%d i=%d case 1 mismatch: kind=%d pay=%d orig=%d", w, i, n.KindFlags(), n.Payload(), n.OrigNameID())
						return
					}
				case 2:
					wantOrig := origSID
					if i%8 == 2 {
						wantOrig = 0
					}
					if n.KindFlags() != 3 || StringID(n.Payload()) != origSID || n.OrigNameID() != wantOrig {
						t.Errorf("w=%d i=%d case 2 mismatch: kind=%d pay=%d orig=%d want=%d", w, i, n.KindFlags(), n.Payload(), n.OrigNameID(), wantOrig)
						return
					}
				default:
					if n.KindFlags() != 0 || n.Length() != int32(i+1) || n.OrigNameID() != origSID {
						t.Errorf("w=%d i=%d case 3 mismatch: kind=%d len=%d orig=%d", w, i, n.KindFlags(), n.Length(), n.OrigNameID())
						return
					}
				}
			}
		}()
	}
	wg.Wait()
}

func TestWave24A36ByteLayoutInlinePayloadAndJSTypeBeforeCast(t *testing.T) {
	assertRuntimeError := func(name string, fn func()) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("%s: expected runtime.Error panic on Node(0), got nil", name)
			}
			if _, ok := r.(runtime.Error); !ok {
				t.Fatalf("%s: expected runtime.Error on Node(0), got %T (%v)", name, r, r)
			}
		}()
		fn()
	}

	// 1. 32-byte noscan NodeData layout + widened 256-slot (1 KiB) PagedColumn verification.
	if sz := reflect.TypeFor[NodeData]().Size(); sz != 32 {
		t.Fatalf("sizeof(NodeData) = %d, want 32", sz)
	}
	if hasGoPointers(reflect.TypeFor[NodeData]()) {
		t.Fatalf("NodeData contains Go pointers")
	}
	if columnPageShift != 8 || columnPageLen != 256 || pagesPerBlock != 64 {
		t.Fatalf("unexpected PagedColumn geometry: shift=%d len=%d pages=%d",
			columnPageShift, columnPageLen, pagesPerBlock)
	}
	if sz := reflect.TypeFor[columnPage]().Size(); sz != 1024 {
		t.Fatalf("sizeof(columnPage) = %d, want 1024", sz)
	}
	if hasGoPointers(reflect.TypeFor[columnPage]()) {
		t.Fatalf("columnPage contains Go pointers")
	}
	if sz := reflect.TypeFor[PagedColumn]().Size(); sz != 64*8 {
		t.Fatalf("sizeof(PagedColumn) = %d, want 512", sz)
	}

	// NilNode NPE / nil parity.
	if NilNode.JSTypeBeforeCastObj() != nil {
		t.Fatalf("NilNode.JSTypeBeforeCastObj() = %#v, want nil", NilNode.JSTypeBeforeCastObj())
	}
	if NilNode.Length() != 0 {
		t.Fatalf("NilNode.Length() = %d, want 0", NilNode.Length())
	}
	assertRuntimeError("Node(0).SetJSTypeBeforeCastObj(nil)", func() {
		Node(0).SetJSTypeBeforeCastObj(nil)
	})
	assertRuntimeError("Node(0).SetJSTypeBeforeCastObj(non-nil)", func() {
		Node(0).SetJSTypeBeforeCastObj(&fakeTCHolder{name: "cast"})
	})

	s1 := NewNodeStore()
	s2 := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer SetCurrentNodeStore(prev)

	// 2. Length16 inline (0..65534) vs 0xFFFF overflow boundary & reset back to inline/0.
	var alloc NodeAllocator
	alloc.InitLocal()
	sid := alloc.Intern("wave24a_str")
	lenTargets := []Node{
		alloc.AllocNode1Cached(40),              // plain (kind=0)
		alloc.AllocStringNodeCached(38, sid),    // string (kind=3)
		alloc.AllocNumberNode(39, 123.456, 0),   // number (kind=1)
	}
	alloc.Flush()

	for idx, n := range lenTargets {
		for _, inl := range []int32{0, 1, 1024, 65534} {
			n.SetLength(inl)
			if got := n.Length(); got != inl {
				t.Fatalf("target[%d] inline SetLength(%d): got %d", idx, inl, got)
			}
			pg := loadLen16Page(uint32(n) >> len16PageShift)
			if pg[n&len16PageMask] != uint16(inl) {
				t.Fatalf("target[%d] expected len16=%d, got %d", idx, inl, pg[n&len16PageMask])
			}
		}
	}
	if s1.auxLen() != 0 {
		t.Fatalf("inline lengths (0..65534) allocated aux entries eagerly: %d", s1.auxLen())
	}

	for idx, n := range lenTargets {
		boundarySeq := []int32{65534, 65535, 65536, 1_000_000, -1, math.MinInt32, math.MaxInt32, 65534, 65535, 42, 0}
		for _, want := range boundarySeq {
			n.SetLength(want)
			if got := n.Length(); got != want {
				t.Fatalf("target[%d] SetLength(%d): got %d", idx, want, got)
			}
			pg := loadLen16Page(uint32(n) >> len16PageShift)
			if want >= 0 && want < 0xFFFF {
				if pg[n&len16PageMask] != uint16(want) {
					t.Fatalf("target[%d] SetLength(%d) did not restore len16: got 0x%X", idx, want, pg[n&len16PageMask])
				}
				if gotLO := s1.getAuxLengthOverflow(uint32(n)); gotLO != 0 {
					t.Fatalf("target[%d] SetLength(%d) left stale aux lengthOverflow=%d", idx, want, gotLO)
				}
			} else if pg[n&len16PageMask] != 0xFFFF {
				t.Fatalf("target[%d] SetLength(%d) expected len16=0xFFFF, got 0x%X", idx, want, pg[n&len16PageMask])
			}
		}
	}
	if StringID(lenTargets[1].Payload()) != sid || lenTargets[2].NumberVal() != 123.456 {
		t.Fatalf("Length16/overflow transitions clobbered String/Number Payload")
	}

	// 3. SourceFile overflow (>= 65535) via SetTemplate (tmplCold == true -> ClonePropsFrom).
	highSFTmpl := AllocNodeToken(0, 94)
	highSFTmpl.SetSourceFile(75_000)
	highSFTmpl.SetPropFlags(0x22446688)
	alloc.InitLocal()
	alloc.SetTemplate(highSFTmpl)
	if !alloc.tmplCold {
		t.Fatalf("expected tmplCold == true for overflow SourceFile=75000")
	}
	sfChild1 := alloc.AllocNode1Cached(10)
	sfChild2 := alloc.AllocStringNodeCached(38, sid)
	alloc.Flush()
	if sfChild1.SourceFile() != 75_000 || sfChild2.SourceFile() != 75_000 || sfChild2.PropFlags() != 0x22446688 {
		t.Fatalf("overflow SourceFile template stamping failed: %d / %d", sfChild1.SourceFile(), sfChild2.SourceFile())
	}

	// 4. Plain-node (kind == 0) inline ChangeTime + inline OrigNameID + both coexistence orders!
	sPlain := NewNodeStore()
	SetCurrentNodeStore(sPlain)

	idFn1 := Intern("wave24a_fn_orig_1")
	idFn2 := Intern("wave24a_fn_orig_2")
	idFn3 := Intern("wave24a_fn_orig_3")

	// 4a. Pure inline ChangeTime and pure inline OrigNameID allocate ZERO PagedColumns!
	fnCTOnly := AllocNodeToken(0, 52) // FUNCTION
	fnCTOnly.SetLength(240)
	fnCTOnly.SetChangeTime(19)
	fnOrigOnly := AllocNodeToken(0, 52) // FUNCTION
	fnOrigOnly.SetLength(180)
	fnOrigOnly.SetOrigNameID(idFn1)

	if sPlain.auxLen() != 0 || len(sPlain.origNames) != 0 {
		t.Fatalf("plain-node inline ChangeTime / OrigNameID allocated side storage: aux=%d orig=%d",
			sPlain.auxLen(), len(sPlain.origNames))
	}
	if fnCTOnly.ChangeTime() != 19 || fnCTOnly.OrigNameID() != 0 || fnCTOnly.Length() != 240 {
		t.Fatalf("fnCTOnly mismatch: ct=%d orig=%d len=%d", fnCTOnly.ChangeTime(), fnCTOnly.OrigNameID(), fnCTOnly.Length())
	}
	if fnOrigOnly.OrigNameID() != idFn1 || fnOrigOnly.ChangeTime() != 0 || fnOrigOnly.Length() != 180 {
		t.Fatalf("fnOrigOnly mismatch: orig=%d ct=%d len=%d", fnOrigOnly.OrigNameID(), fnOrigOnly.ChangeTime(), fnOrigOnly.Length())
	}
	// Overwrite and clear inline OrigNameID & inline ChangeTime with 0 PagedColumns/aux!
	fnOrigOnly.SetOrigNameID(idFn2)
	fnCTOnly.SetChangeTime(-42)
	if fnOrigOnly.OrigNameID() != idFn2 || fnCTOnly.ChangeTime() != -42 || len(sPlain.origNames) != 0 || sPlain.auxLen() != 0 {
		t.Fatalf("overwrite inline plain Payload allocated PagedColumns or aux")
	}
	fnOrigOnly.SetOrigNameID(0)
	fnCTOnly.SetChangeTime(0)
	if fnOrigOnly.OrigNameID() != 0 || fnOrigOnly.Payload() != 0 || fnCTOnly.ChangeTime() != 0 || fnCTOnly.Payload() != 0 {
		t.Fatalf("clearing inline plain Payload left dirty state: origPay=%d ctPay=%d", fnOrigOnly.Payload(), fnCTOnly.Payload())
	}

	// 4b. Coexistence Order A: SetOrigNameID -> SetChangeTime -> overwrite both -> clear ChangeTime -> set OrigName -> clear OrigName.
	fnOrderA := AllocNodeToken(0, 52)
	fnOrderA.SetOrigNameID(idFn1)
	fnOrderA.SetChangeTime(101)
	if fnOrderA.OrigNameID() != idFn1 || fnOrderA.ChangeTime() != 101 {
		t.Fatalf("Order A (Orig->CT) mismatch: orig=%d (want %d), ct=%d (want 101)", fnOrderA.OrigNameID(), idFn1, fnOrderA.ChangeTime())
	}
	if sPlain.auxLen() != 0 {
		t.Fatalf("Order A should keep ChangeTime inline in Payload (0 aux entries), got %d", sPlain.auxLen())
	}
	fnOrderA.SetOrigNameID(idFn2)
	fnOrderA.SetChangeTime(202)
	if fnOrderA.OrigNameID() != idFn2 || fnOrderA.ChangeTime() != 202 {
		t.Fatalf("Order A overwrite mismatch: orig=%d ct=%d", fnOrderA.OrigNameID(), fnOrderA.ChangeTime())
	}
	// Clear ChangeTime to 0 while OrigNameID stays idFn2, then set OrigNameID(idFn3) (must move back inline & clear origNameField!), then clear OrigNameID(0).
	fnOrderA.SetChangeTime(0)
	if fnOrderA.ChangeTime() != 0 || fnOrderA.OrigNameID() != idFn2 {
		t.Fatalf("Order A after SetChangeTime(0): ct=%d orig=%d", fnOrderA.ChangeTime(), fnOrderA.OrigNameID())
	}
	fnOrderA.SetOrigNameID(idFn3)
	if fnOrderA.OrigNameID() != idFn3 || fnOrderA.Rec().KindFlags&kindFlagInlineOrigName == 0 || origNameField.get(fnOrderA) != 0 {
		t.Fatalf("Order A re-inlining OrigNameID failed: orig=%d col=%d", fnOrderA.OrigNameID(), origNameField.get(fnOrderA))
	}
	fnOrderA.SetOrigNameID(0)
	if fnOrderA.OrigNameID() != 0 || fnOrderA.Payload() != 0 {
		t.Fatalf("Order A final SetOrigNameID(0) failed: orig=%d pay=%d", fnOrderA.OrigNameID(), fnOrderA.Payload())
	}

	// Also test Order A clearing OrigNameID(0) while ChangeTime remains non-zero!
	fnOrderA2 := AllocNodeToken(0, 52)
	fnOrderA2.SetOrigNameID(idFn1)
	fnOrderA2.SetChangeTime(303)
	fnOrderA2.SetOrigNameID(0)
	if fnOrderA2.OrigNameID() != 0 || fnOrderA2.ChangeTime() != 303 {
		t.Fatalf("Order A2 SetOrigNameID(0) clobbered ChangeTime: orig=%d ct=%d", fnOrderA2.OrigNameID(), fnOrderA2.ChangeTime())
	}

	// 4c. Coexistence Order B: SetChangeTime -> SetOrigNameID -> overwrite both -> clear OrigNameID -> clear ChangeTime.
	fnOrderB := AllocNodeToken(0, 52)
	fnOrderB.SetChangeTime(404)
	fnOrderB.SetOrigNameID(idFn1)
	if fnOrderB.ChangeTime() != 404 || fnOrderB.OrigNameID() != idFn1 {
		t.Fatalf("Order B (CT->Orig) mismatch: ct=%d orig=%d", fnOrderB.ChangeTime(), fnOrderB.OrigNameID())
	}
	fnOrderB.SetChangeTime(-505)
	fnOrderB.SetOrigNameID(idFn2)
	if fnOrderB.ChangeTime() != -505 || fnOrderB.OrigNameID() != idFn2 {
		t.Fatalf("Order B overwrite mismatch: ct=%d orig=%d", fnOrderB.ChangeTime(), fnOrderB.OrigNameID())
	}
	fnOrderB.SetOrigNameID(0)
	if fnOrderB.OrigNameID() != 0 || fnOrderB.ChangeTime() != -505 {
		t.Fatalf("Order B SetOrigNameID(0) mismatch: orig=%d ct=%d", fnOrderB.OrigNameID(), fnOrderB.ChangeTime())
	}
	fnOrderB.SetChangeTime(0)
	if fnOrderB.ChangeTime() != 0 || fnOrderB.Payload() != 0 {
		t.Fatalf("Order B SetChangeTime(0) mismatch: ct=%d pay=%d", fnOrderB.ChangeTime(), fnOrderB.Payload())
	}

	// 4d. Migration guards when SetPayload / SetKindFlags / SetNumberVal are called on kind==0 nodes.
	migCTPay := AllocNodeToken(0, 52)
	migCTPay.SetChangeTime(777)
	migCTPay.SetPayload(0xCAFEBABE)
	if migCTPay.Payload() != 0xCAFEBABE || migCTPay.ChangeTime() != 777 {
		t.Fatalf("SetPayload clobbered inline ChangeTime: pay=0x%x ct=%d", migCTPay.Payload(), migCTPay.ChangeTime())
	}
	migCTKind := AllocNodeToken(0, 52)
	migCTKind.SetChangeTime(888)
	migCTKind.SetKindFlags(3)
	if migCTKind.KindFlags() != 3 || migCTKind.Payload() != 0 || migCTKind.ChangeTime() != 888 {
		t.Fatalf("SetKindFlags(3) clobbered inline ChangeTime: kind=%d pay=%d ct=%d",
			migCTKind.KindFlags(), migCTKind.Payload(), migCTKind.ChangeTime())
	}
	migOrigNum := AllocNodeToken(0, 52)
	migOrigNum.SetOrigNameID(idFn1)
	migOrigNum.SetNumberVal(6.25)
	if migOrigNum.NumberVal() != 6.25 || migOrigNum.OrigNameID() != idFn1 {
		t.Fatalf("SetNumberVal clobbered inline OrigNameID: num=%v orig=%d", migOrigNum.NumberVal(), migOrigNum.OrigNameID())
	}

	SetCurrentNodeStore(s1)

	// 5. JSTypeBeforeCastObj & SetJSTypeBeforeCastObj (TCHandleHolder, typed-nil, overwrite, clear, store switch).
	castType1 := &fakeTCHolder{name: "UnknownType"}
	castType2 := &fakeTCHolder{name: "StringType"}
	mainType := &fakeTCHolder{name: "NumberType"}
	var typedNilTC *fakeTCHolder
	var typedNilJava *fakeJavaNilOnly

	nCast := AllocNodeToken(0, 42) // CAST
	nCast.SetJSTypeBeforeCastObj(nil)
	nCast.SetJSTypeBeforeCastObj(any(typedNilTC))
	nCast.SetJSTypeBeforeCastObj(any(typedNilJava))
	if s1.auxLen() != 0 || nCast.JSTypeBeforeCastObj() != nil {
		t.Fatalf("nil/typed-nil SetJSTypeBeforeCastObj allocated aux eagerly: entries=%d", s1.auxLen())
	}
	if allocs := testing.AllocsPerRun(100, func() {
		if nCast.JSTypeBeforeCastObj() != nil {
			t.Fatal("expected nil")
		}
	}); allocs > 0 {
		t.Fatalf("unset JSTypeBeforeCastObj() allocated %v objects, want 0", allocs)
	}

	if !nCast.SetJSTypeObj(mainType) {
		t.Fatal("SetJSTypeObj failed")
	}
	nCast.SetJSTypeBeforeCastObj(castType1)
	if castType1.tcHandle.Load() != 0 {
		t.Fatalf("expected SetJSTypeBeforeCastObj stored directly in NodeStore.aux without tcHandles entry, got %d", castType1.tcHandle.Load())
	}
	if nCast.JSTypeBeforeCastObj() != castType1 || nCast.JSTypeObj() != mainType || nCast.HasColdProps() {
		t.Fatalf("SetJSTypeBeforeCastObj mismatch or leaked into coldProps: before=%#v main=%#v cold=%v",
			nCast.JSTypeBeforeCastObj(), nCast.JSTypeObj(), nCast.HasColdProps())
	}
	// Warm fast path: 0 heap allocations!
	if allocs := testing.AllocsPerRun(100, func() {
		nCast.SetJSTypeBeforeCastObj(castType2)
		if nCast.JSTypeBeforeCastObj() != castType2 {
			t.Fatal("mismatch")
		}
		nCast.SetJSTypeBeforeCastObj(nil)
		if nCast.JSTypeBeforeCastObj() != nil {
			t.Fatal("expected nil")
		}
		nCast.SetJSTypeBeforeCastObj(castType1)
	}); allocs > 0 {
		t.Fatalf("warm JSTypeBeforeCastObj/SetJSTypeBeforeCastObj allocated %v objects, want 0", allocs)
	}
	nCast.SetJSTypeBeforeCastObj(any(typedNilTC))
	if nCast.JSTypeBeforeCastObj() != nil || nCast.Rec().KindFlags&kindFlagHasTypeBeforeCast != 0 {
		t.Fatalf("typed-nil SetJSTypeBeforeCastObj failed to clear flag")
	}
	nCast.SetJSTypeBeforeCastObj(castType2)

	// Store switching + static template node for JSTypeBeforeCastObj.
	staticCastTmpl := NewTemplateNode(42)
	staticCastTmpl.SetJSTypeBeforeCastObj(castType1)
	defer staticCastTmpl.SetJSTypeBeforeCastObj(nil)

	SetCurrentNodeStore(s2)
	if nCast.JSTypeBeforeCastObj() != nil {
		t.Fatalf("unregistered s1 JSTypeBeforeCastObj visible while s2 active")
	}
	if staticCastTmpl.JSTypeBeforeCastObj() != castType1 {
		t.Fatalf("static template JSTypeBeforeCastObj lost across store switch")
	}
	SetCurrentNodeStore(s1)
	if nCast.JSTypeBeforeCastObj() != castType2 {
		t.Fatalf("restored s1 JSTypeBeforeCastObj = %#v, want castType2", nCast.JSTypeBeforeCastObj())
	}

	// 6. ClonePropsFrom from plain SCRIPT template (with inline ChangeTime, SourceFile > 65535, JSDoc, TypeBeforeCast, ColdProps)
	// onto PlainNode, PlainNode-with-OrigName, StringNode, and NumberNode!
	scriptTmpl := AllocNodeToken(0, 94) // SCRIPT (kind == 0)
	scriptTmpl.SetSourceFile(88_888)
	scriptTmpl.SetPropFlags(0x11223344)
	scriptTmpl.SetChangeTime(4242) // stored inline in scriptTmpl.Rec().Payload!
	doc := &fakeJSDocHolder{comment: "/** @fileoverview */"}
	coldProp := &fakeColdProp{propType: 9, val: "script-cold"}
	scriptTmpl.SetJSDocObj(doc)
	scriptTmpl.SetJSTypeBeforeCastObj(castType1)
	scriptTmpl.SetPropList(coldProp)

	if scriptTmpl.Rec().Payload != 4242 {
		t.Fatalf("expected scriptTmpl ChangeTime inline in Payload, got %d", scriptTmpl.Rec().Payload)
	}

	dstPlain := AllocNodeToken(0, 40)
	dstPlainOrig := AllocNodeToken(0, 52)
	dstPlainOrig.SetOrigNameID(idFn1) // inline OrigNameID before ClonePropsFrom!
	dstStr := AllocNodeToken(3, 38)
	dstStr.SetStringPayload(sid)
	dstStr.SetOrigNameID(sid)
	dstNum := AllocNodeToken(1, 39)
	dstNum.SetNumberVal(-98.75)

	for idx, dst := range []Node{dstPlain, dstPlainOrig, dstStr, dstNum} {
		if !dst.ClonePropsFrom(scriptTmpl) {
			t.Fatalf("dst[%d] ClonePropsFrom(scriptTmpl) returned false", idx)
		}
		if dst.SourceFile() != 88_888 ||
			dst.PropFlags() != 0x11223344 ||
			dst.ChangeTime() != 4242 ||
			dst.JSDocObj() != doc ||
			dst.JSTypeBeforeCastObj() != castType1 ||
			dst.PropList() != coldProp {
			t.Fatalf("dst[%d] cloned props mismatch: sf=%d pf=0x%x ct=%d doc=%#v cast=%#v prop=%#v",
				idx, dst.SourceFile(), dst.PropFlags(), dst.ChangeTime(), dst.JSDocObj(), dst.JSTypeBeforeCastObj(), dst.PropList())
		}
	}
	if dstPlainOrig.OrigNameID() != idFn1 {
		t.Fatalf("ClonePropsFrom clobbered dstPlainOrig OrigNameID: got %d, want %d", dstPlainOrig.OrigNameID(), idFn1)
	}
	if dstStr.KindFlags() != 3 || StringID(dstStr.Payload()) != sid || dstStr.OrigNameID() != sid {
		t.Fatalf("ClonePropsFrom from plain SCRIPT corrupted dstStr Payload/OrigName: kind=%d pay=%d orig=%d",
			dstStr.KindFlags(), dstStr.Payload(), dstStr.OrigNameID())
	}
	if dstNum.KindFlags() != 1 || dstNum.NumberVal() != -98.75 {
		t.Fatalf("ClonePropsFrom from plain SCRIPT corrupted dstNum NumberVal: kind=%d num=%v",
			dstNum.KindFlags(), dstNum.NumberVal())
	}

	// Reject ClonePropsFrom if target already has JSTypeBeforeCastObj set!
	dirtyCast := AllocNodeToken(0, 42)
	dirtyCast.SetJSTypeBeforeCastObj(castType1)
	if dirtyCast.ClonePropsFrom(scriptTmpl) || dirtyCast.ClonePropsFrom(NilNode) {
		t.Fatalf("ClonePropsFrom must return false when target already has JSTypeBeforeCastObj")
	}

	// Also verify NodeAllocator.SetTemplate with scriptTmpl (cold template) onto Plain, String, and Number nodes!
	alloc.InitLocal()
	alloc.SetTemplate(scriptTmpl)
	if !alloc.tmplCold {
		t.Fatalf("expected tmplCold == true when scriptTmpl has cold/JSDoc/ChangeTime/TypeBeforeCast")
	}
	aPlain := alloc.AllocNode1Cached(40)
	aStr := alloc.AllocStringNodeCached(38, sid)
	aNum := alloc.AllocNumberNode(39, 3.5, scriptTmpl)
	alloc.Flush()
	if aPlain.ChangeTime() != 4242 || aPlain.JSTypeBeforeCastObj() != castType1 ||
		StringID(aStr.Payload()) != sid || aStr.ChangeTime() != 4242 || aStr.JSTypeBeforeCastObj() != castType1 ||
		aNum.NumberVal() != 3.5 || aNum.ChangeTime() != 4242 || aNum.JSTypeBeforeCastObj() != castType1 {
		t.Fatalf("NodeAllocator with scriptTmpl corrupted payloads or cold props")
	}

	// 7. Multi-goroutine -race stress across shared 256-node pages and multiple blocks.
	const sharedCount = 4096 // 16 256-node pages
	shared := make([]Node, sharedCount)
	for i := range shared {
		shared[i] = AllocNodeToken(uint8(i%4), uint8((i%100)+1))
		if i%4 == 3 {
			shared[i].SetStringPayload(sid)
		}
	}

	const workers = 8
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := range workers {
		w := w
		go func() {
			defer wg.Done()
			for i := w; i < sharedCount; i += workers {
				n := shared[i]
				if i%2 == 0 {
					n.SetLength(int32(i * 20)) // crosses 65535 for i >= 3277
				} else {
					n.SetLength(-int32(i) - 1)
				}
				n.SetSourceFile(uint32(70_000 + i))
				n.SetOrigNameID(idFn1)
				n.SetChangeTime(int32(i + 1))
				if i%3 == 0 {
					n.SetJSTypeBeforeCastObj(castType1)
				} else if i%3 == 1 {
					n.SetJSTypeBeforeCastObj(castType2)
					n.SetJSTypeBeforeCastObj(nil)
				}
			}
		}()
	}
	wg.Wait()

	for i, n := range shared {
		wantLen := int32(i * 20)
		if i%2 == 1 {
			wantLen = -int32(i) - 1
		}
		if n.Length() != wantLen || n.SourceFile() != uint32(70_000+i) || n.OrigNameID() != idFn1 || n.ChangeTime() != int32(i+1) {
			t.Fatalf("shared[%d] mismatch: len=%d/%d sf=%d orig=%d ct=%d",
				i, n.Length(), wantLen, n.SourceFile(), n.OrigNameID(), n.ChangeTime())
		}
		if i%4 == 3 && StringID(n.Payload()) != sid {
			t.Fatalf("shared[%d] StringNode payload corrupted", i)
		}
		var wantCast any
		if i%3 == 0 {
			wantCast = castType1
		}
		if n.JSTypeBeforeCastObj() != wantCast {
			t.Fatalf("shared[%d] JSTypeBeforeCastObj = %#v, want %#v", i, n.JSTypeBeforeCastObj(), wantCast)
		}
	}
}

func TestWave25ANodeDataCacheLineGroupingAndPartialWriteFastPaths(t *testing.T) {
	typ := reflect.TypeFor[NodeData]()
	if typ.Size() != 32 {
		t.Fatalf("reflect.TypeFor[NodeData]().Size() = %d, want 32", typ.Size())
	}
	// Verify Next, First, Token, KindFlags, SourceFile16 all sit in bytes 0..11,
	// and at 32 bytes (exact half cache line), 0% of NodeData records ever straddle a 64-byte cache line.
	for _, name := range []string{"Next", "First", "Token", "KindFlags", "SourceFile16"} {
		f, ok := typ.FieldByName(name)
		if !ok {
			t.Fatalf("missing field %s", name)
		}
		if f.Offset+f.Type.Size() > 12 {
			t.Fatalf("field %s ends at %d, want <= 12", name, f.Offset+f.Type.Size())
		}
	}

	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	tmpl := AllocNodeToken(0, 94)
	tmpl.SetSourceFile(12345)
	tmpl.SetPropFlags(0x00A5A500)

	var a NodeAllocator
	a.InitLocal()
	a.SetTemplate(tmpl)
	sid := a.Intern("wave25a_prop")

	s1 := a.AllocStringNodeCached(38, sid)
	s2 := a.AllocStringNode(38, sid, tmpl)
	nZero := a.AllocNumberNode(39, 0.0, tmpl)
	nNegZero := a.AllocNumberNode(39, math.Copysign(0, -1), tmpl)
	nVal := a.AllocNumberNode(39, 42.125, tmpl)
	p2 := a.AllocNode2(10, s1, tmpl)
	p3 := a.AllocNode3(11, s2, nZero, tmpl)
	p4 := a.AllocNode4(12, nNegZero, nVal, p2, tmpl)
	a.Flush()

	if StringID(s1.Payload()) != sid || s1.KindFlags() != 3 || s1.SourceFile() != 12345 || s1.Lineno() != -1 {
		t.Fatalf("AllocStringNodeCached mismatch: %+v", *s1.Rec())
	}
	if nZero.NumberVal() != 0 || math.Signbit(nZero.NumberVal()) || nZero.Payload() != 0 {
		t.Fatalf("+0.0 NumberNode mismatch: val=%v pay=%d", nZero.NumberVal(), nZero.Payload())
	}
	if nNegZero.NumberVal() != 0 || !math.Signbit(nNegZero.NumberVal()) || nNegZero.Payload() == 0 {
		t.Fatalf("-0.0 NumberNode mismatch: val=%v pay=%d", nNegZero.NumberVal(), nNegZero.Payload())
	}
	if nVal.NumberVal() != 42.125 || nVal.KindFlags() != 1 {
		t.Fatalf("42.125 NumberNode mismatch: val=%v kind=%d", nVal.NumberVal(), nVal.KindFlags())
	}
	blk := BlockOf(p4)
	if f, nx, tk := blk.FirstNextTok(p4); f != nNegZero || nx != 0 || tk != 12 {
		t.Fatalf("FirstNextTok(p4) = (%d, %d, %d)", f, nx, tk)
	}
	if p3.FirstChild() != s2 || s2.Next() != nZero || nZero.Previous() != s2 || s2.Previous() != nZero {
		t.Fatalf("AllocNode3 child links broken")
	}
}

func TestWave26AExtOrigNamePresenceGateAndJSDocCastCoexistence(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	idFoo := Intern("wave26a_foo")
	idBar := Intern("wave26a_bar")
	idShort := Intern("a")

	// Populate external origNameField on one node in page 0 of block 0 so both
	// origNameField[block] and PagedColumn[page 0] are non-nil.
	renamed := AllocNodeToken(3, 38)
	renamed.SetStringPayload(idFoo)
	renamed.SetOrigNameID(idFoo)
	renamed.SetStringPayload(idShort) // migrates idFoo to origNameField + sets kindFlagHasExtOrigName
	if renamed.Rec().KindFlags&kindFlagHasExtOrigName == 0 || renamed.Rec().KindFlags&kindFlagInlineOrigName != 0 {
		t.Fatalf("expected kindFlagHasExtOrigName set and kindFlagInlineOrigName cleared, got 0x%x", renamed.Rec().KindFlags)
	}
	if renamed.OrigNameID() != idFoo {
		t.Fatalf("renamed.OrigNameID() = %d, want %d", renamed.OrigNameID(), idFoo)
	}

	// Neighboring unset nodes in the same 256-node page must short-circuit inline via kindFlagAnyOrigName == 0.
	unsetPlain := AllocNodeToken(0, 40)
	unsetStr := AllocNodeToken(3, 38)
	unsetStr.SetStringPayload(idBar)
	if unsetPlain.OrigNameID() != 0 || unsetStr.OrigNameID() != 0 {
		t.Fatalf("unset nodes returned non-zero OrigNameID: %d / %d", unsetPlain.OrigNameID(), unsetStr.OrigNameID())
	}

	// Coexistence of JSDocObj and JSTypeBeforeCastObj on the same node under shared presence bit!
	doc1 := &fakeJSDocHolder{comment: "/** @type {string} */"}
	doc2 := &fakeJSDocHolder{comment: "/** @type {number} */"}
	cast1 := &fakeTCHolder{name: "UnknownType"}
	cast2 := &fakeTCHolder{name: "AnyType"}

	// Case A: Set both -> clear JSTypeBeforeCastObj(nil) (exact ConvertTypesToColors sequence!) -> JSDocObj must survive!
	nBothA := AllocNodeToken(0, 42)
	nBothA.SetJSDocObj(doc1)
	nBothA.SetJSTypeBeforeCastObj(cast1)
	if nBothA.JSDocObj() != doc1 || nBothA.JSTypeBeforeCastObj() != cast1 {
		t.Fatalf("Case A initial mismatch: doc=%#v cast=%#v", nBothA.JSDocObj(), nBothA.JSTypeBeforeCastObj())
	}
	nBothA.SetJSTypeBeforeCastObj(nil)
	if nBothA.JSTypeBeforeCastObj() != nil || nBothA.JSDocObj() != doc1 {
		t.Fatalf("Case A SetJSTypeBeforeCastObj(nil) lost JSDocObj: doc=%#v cast=%#v", nBothA.JSDocObj(), nBothA.JSTypeBeforeCastObj())
	}
	nBothA.SetJSDocObj(nil)
	if nBothA.JSDocObj() != nil || nBothA.Rec().KindFlags&kindFlagHasJSDocOrCast != 0 {
		t.Fatalf("Case A clearing both failed to clear kindFlagHasJSDocOrCast: 0x%x", nBothA.Rec().KindFlags)
	}

	// Case B: Set both -> clear JSDocObj(nil) -> JSTypeBeforeCastObj must survive!
	nBothB := AllocNodeToken(0, 42)
	nBothB.SetJSTypeBeforeCastObj(cast2)
	nBothB.SetJSDocObj(doc2)
	nBothB.SetJSDocObj(nil)
	if nBothB.JSDocObj() != nil || nBothB.JSTypeBeforeCastObj() != cast2 {
		t.Fatalf("Case B SetJSDocObj(nil) lost JSTypeBeforeCastObj: doc=%#v cast=%#v", nBothB.JSDocObj(), nBothB.JSTypeBeforeCastObj())
	}
	nBothB.SetJSTypeBeforeCastObj(nil)
	if nBothB.JSTypeBeforeCastObj() != nil || nBothB.Rec().KindFlags&kindFlagHasJSDocOrCast != 0 {
		t.Fatalf("Case B clearing both failed to clear kindFlagHasJSDocOrCast: 0x%x", nBothB.Rec().KindFlags)
	}

	// Case C: JSDoc-only and Cast-only ClonePropsFrom isolation.
	tmplDocOnly := AllocNodeToken(0, 40)
	tmplDocOnly.SetJSDocObj(doc1)
	dstDoc := AllocNodeToken(0, 40)
	if !dstDoc.ClonePropsFrom(tmplDocOnly) || dstDoc.JSDocObj() != doc1 || dstDoc.JSTypeBeforeCastObj() != nil {
		t.Fatalf("Case C tmplDocOnly ClonePropsFrom mismatch: doc=%#v cast=%#v", dstDoc.JSDocObj(), dstDoc.JSTypeBeforeCastObj())
	}
	tmplCastOnly := AllocNodeToken(0, 40)
	tmplCastOnly.SetJSTypeBeforeCastObj(cast1)
	dstCast := AllocNodeToken(0, 40)
	if !dstCast.ClonePropsFrom(tmplCastOnly) || dstCast.JSTypeBeforeCastObj() != cast1 || dstCast.JSDocObj() != nil {
		t.Fatalf("Case C tmplCastOnly ClonePropsFrom mismatch: doc=%#v cast=%#v", dstCast.JSDocObj(), dstCast.JSTypeBeforeCastObj())
	}
}

func BenchmarkWave26ANodeHotPaths(b *testing.B) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	idFoo := Intern("goog_ui_Component")
	idShort := Intern("a")
	seeded := AllocNodeToken(3, 38)
	seeded.SetStringPayload(idFoo)
	seeded.SetOrigNameID(idFoo)
	seeded.SetStringPayload(idShort)

	nUnset := AllocNodeToken(0, 40)
	nUnset.SetLength(42)
	nInline := AllocNodeToken(3, 38)
	nInline.SetStringPayload(idFoo)
	nInline.SetOrigNameID(idFoo)

	b.Run("NodeLengthInline", func(b *testing.B) {
		var sink int32
		for i := 0; i < b.N; i++ {
			sink += nUnset.Length()
		}
		_ = sink
	})
	b.Run("OrigNameIDUnsetSamePage", func(b *testing.B) {
		var sink StringID
		for i := 0; i < b.N; i++ {
			sink += nUnset.OrigNameID()
		}
		_ = sink
	})
	b.Run("OrigNameIDInline", func(b *testing.B) {
		var sink StringID
		for i := 0; i < b.N; i++ {
			sink += nInline.OrigNameID()
		}
		_ = sink
	})
	b.Run("SetJSTypeBeforeCastNilUnset", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nUnset.SetJSTypeBeforeCastObj(nil)
		}
	})
	b.Run("SetJSTypeNilUnset", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nUnset.SetJSTypeObj(nil)
		}
	})
}

func TestWave27AShardedAuxAndReleaseCurrentStoreJSTypes(t *testing.T) {
	if auxShards != 64 || jsdocShards != 64 {
		t.Fatalf("auxShards=%d jsdocShards=%d, want 64", auxShards, jsdocShards)
	}
	if sz := reflect.TypeFor[nodeAuxEntry]().Size(); sz != 24 {
		t.Fatalf("sizeof(nodeAuxEntry) = %d, want 24 (typeBeforeCast 16B + changeTime 4B + length 4B)", sz)
	}

	s1 := NewNodeStore()
	sShadow := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer SetCurrentNodeStore(prev)

	// 1. Verify (uint32(n) * 0x9e3779b9) >> 26 spreads consecutive nodes evenly across all 64 jsdocShards
	// while keeping s1.aux completely empty.
	const consecutive = 4096 // 64 per shard on average
	consecNodes := make([]Node, consecutive)
	sharedDoc := &fakeJSDocHolder{comment: "/** @const */"}
	for i := range consecNodes {
		n := AllocNodeToken(0, 40)
		n.SetJSDocObj(sharedDoc)
		consecNodes[i] = n
	}
	if s1.auxLen() != 0 {
		t.Fatalf("JSDoc-only nodes allocated %d s1.aux entries, want 0", s1.auxLen())
	}
	for shardIdx := range s1.jsdocShards {
		cnt := int(s1.jsdocShards[shardIdx].count)
		if cnt < 32 || cnt > 96 {
			t.Fatalf("jsdocShard %d received %d/4096 consecutive nodes (expected ~64 uniform spread)", shardIdx, cnt)
		}
	}
	for _, n := range consecNodes {
		n.SetJSDocObj(nil)
	}
	if s1.jsdocLen() != 0 || s1.auxLen() != 0 {
		t.Fatalf("clearing JSDocObj on all consecutive nodes left jsdoc=%d aux=%d, want 0", s1.jsdocLen(), s1.auxLen())
	}

	// 2. Coexistence of JSDocObj + JSTypeBeforeCastObj + Extra (GetAux/SetAux) + external ChangeTime + Length overflow on the SAME node!
	nMulti := AllocNodeToken(2, 42) // BigInt / non-plain kind=2 so ChangeTime uses aux
	docM := &fakeJSDocHolder{comment: "/** @type {bigint} */"}
	castM := &fakeTCHolder{name: "BigIntBeforeCast"}
	nMulti.SetJSDocObj(docM)
	nMulti.SetJSTypeBeforeCastObj(castM)
	nMulti.SetChangeTime(98765)
	nMulti.SetLength(120_000)
	// Simulate rhino/node.go ensureNodeExtra: GetAux -> set Extra -> SetAux.
	auxPub, ok := nMulti.GetAux()
	if !ok || auxPub.JSDoc != docM {
		t.Fatalf("GetAux before Extra mismatch: ok=%v aux=%+v", ok, auxPub)
	}
	auxPub.Extra = "bigint-999999999999999999999"
	nMulti.SetAux(auxPub)

	if nMulti.JSDocObj() != docM ||
		nMulti.JSTypeBeforeCastObj() != castM ||
		nMulti.ChangeTime() != 98765 ||
		nMulti.Length() != 120_000 {
		t.Fatalf("5-way aux coexistence readback mismatch")
	}
	if gotAux, ok := nMulti.GetAux(); !ok || gotAux.Extra != "bigint-999999999999999999999" || gotAux.JSDoc != docM {
		t.Fatalf("GetAux 5-way mismatch: ok=%v aux=%+v", ok, gotAux)
	}

	// Clear one attribute at a time and verify remaining attributes stay intact until the last one deletes the map entry!
	nMulti.SetJSDocObj(nil)
	if nMulti.JSDocObj() != nil || nMulti.JSTypeBeforeCastObj() != castM || nMulti.ChangeTime() != 98765 || nMulti.Length() != 120_000 {
		t.Fatalf("after SetJSDocObj(nil): unexpected state")
	}
	if gotAux, ok := nMulti.GetAux(); !ok || gotAux.Extra != "bigint-999999999999999999999" || gotAux.JSDoc != nil {
		t.Fatalf("after SetJSDocObj(nil): Extra lost: %+v", gotAux)
	}
	nMulti.SetJSTypeBeforeCastObj(nil)
	if nMulti.JSTypeBeforeCastObj() != nil || nMulti.Rec().KindFlags&kindFlagHasJSDocOrCast != 0 || nMulti.ChangeTime() != 98765 || nMulti.Length() != 120_000 {
		t.Fatalf("after SetJSTypeBeforeCastObj(nil): unexpected state")
	}
	nMulti.SetChangeTime(0)
	if nMulti.ChangeTime() != 0 || nMulti.Length() != 120_000 {
		t.Fatalf("after SetChangeTime(0): unexpected state")
	}
	nMulti.SetLength(12)
	if nMulti.Length() != 12 {
		t.Fatalf("after SetLength(12): got %d", nMulti.Length())
	}
	if gotAux, ok := nMulti.GetAux(); !ok || gotAux.Extra != "bigint-999999999999999999999" {
		t.Fatalf("Extra lost after clearing other cold fields: ok=%v aux=%+v", ok, gotAux)
	}
	nMulti.SetAux(NodeAux{})
	if s1.auxLen() != 0 {
		t.Fatalf("expected 0 aux entries after clearing all 5 attributes, got %d", s1.auxLen())
	}

	// 3. Immediate GC reachability release on SetJSDocObj(nil) and SetJSTypeBeforeCastObj(nil).
	var docCollected, castCollected atomic.Bool
	nGC := AllocNodeToken(0, 42)
	func() {
		ephemeralDoc := &fakeJSDocHolder{comment: "/** ephemeral */"}
		ephemeralCast := &fakeTCHolder{name: "EphemeralCastType"}
		runtime.SetFinalizer(ephemeralDoc, func(*fakeJSDocHolder) { docCollected.Store(true) })
		runtime.SetFinalizer(ephemeralCast, func(*fakeTCHolder) { castCollected.Store(true) })
		nGC.SetJSDocObj(ephemeralDoc)
		nGC.SetJSTypeBeforeCastObj(ephemeralCast)
		nGC.SetJSDocObj(nil)
		nGC.SetJSTypeBeforeCastObj(nil)
	}()
	for range 5 {
		runtime.GC()
		if docCollected.Load() && castCollected.Load() {
			break
		}
	}
	if !docCollected.Load() || !castCollected.Load() {
		t.Fatalf("expected GC collection after clearing JSDoc/Cast: doc=%v cast=%v", docCollected.Load(), castCollected.Load())
	}

	// 4. Shadow compiler store switch (s1 -> sShadow -> s1) preserves outer JSTypes,
	// and ReleaseCurrentStoreJSTypes() releases all non-Color JSTypes while preserving
	// StandardColors singletons and staticStore template handles!
	stdNumberColor := &fakeColorHolder{name: "StandardColors_NUMBER"}
	stdStringColor := &fakeColorHolder{name: "StandardColors_STRING"}
	staticSingleton := &fakeTCHolder{name: "StaticTemplateType"}
	staticTmpl := NewTemplateNode(77)
	staticTmpl.SetTypeOrColor(staticSingleton)
	defer staticTmpl.SetTypeOrColor(nil)

	nColor := AllocNodeToken(0, 39)
	if !nColor.SetColorObj(stdNumberColor) {
		t.Fatal("SetColorObj failed")
	}
	origColorHandle := stdNumberColor.tcHandle.Load()

	const numTypes = 5000 // spans multiple 4096-slot AnyTable blocks
	var releasedTypesCollected atomic.Int32
	typedNodes := make([]Node, numTypes)
	reusableType := &fakeTCHolder{name: "ReusableAfterRelease"}
	nReusable := AllocNodeToken(0, 38)
	if !nReusable.SetJSTypeObj(reusableType) {
		t.Fatal("SetJSTypeObj(reusableType) failed")
	}
	nUnvisitedCast := AllocNodeToken(0, 42)
	nUnvisitedCast.SetJSDocObj(sharedDoc)

	func() {
		for i := range numTypes {
			jt := &fakeTCHolder{name: "ObjectType"}
			runtime.SetFinalizer(jt, func(*fakeTCHolder) { releasedTypesCollected.Add(1) })
			tn := AllocNodeToken(0, 38)
			if !tn.SetJSTypeObj(jt) {
				t.Fatal("SetJSTypeObj failed")
			}
			typedNodes[i] = tn
			if i == 0 {
				nUnvisitedCast.SetJSTypeBeforeCastObj(jt)
			}
		}

		// Nested shadow compilation (RewritePolyfills-style) MUST NOT wipe outer store's JSTypes!
		saved := SetCurrentNodeStore(sShadow)
		shadowNode := AllocNodeToken(0, 39)
		if !shadowNode.SetColorObj(stdStringColor) {
			t.Fatal("shadow SetColorObj failed")
		}
		SetCurrentNodeStore(saved)

		if typedNodes[0].JSTypeObj() == nil || nUnvisitedCast.JSTypeBeforeCastObj() == nil {
			t.Fatal("shadow store switch prematurely cleared outer compilation JSTypes")
		}
	}()

	// Now simulate post-ConvertTypesToColors ReleaseCurrentStoreJSTypes()!
	ReleaseCurrentStoreJSTypes()

	if reusableType.tcHandle.Load() != 0 {
		t.Fatalf("expected ReleaseCurrentStoreJSTypes to reset reusableType.tcHandle to 0, got %d", reusableType.tcHandle.Load())
	}
	for i, tn := range typedNodes {
		if tn.JSTypeObj() != nil || tn.TC() != 0 {
			t.Fatalf("typedNodes[%d] still holds JSType after ReleaseCurrentStoreJSTypes", i)
		}
	}
	// Nodes whose JSTypes were released can now receive Colors cleanly!
	if !typedNodes[0].SetColorObj(stdStringColor) || typedNodes[0].ColorObj() != stdStringColor {
		t.Fatalf("SetColorObj failed on node after ReleaseCurrentStoreJSTypes")
	}
	// Unvisited typeBeforeCast was cleared while JSDoc on the same node survived!
	if nUnvisitedCast.JSTypeBeforeCastObj() != nil || nUnvisitedCast.JSDocObj() != sharedDoc {
		t.Fatalf("after ReleaseCurrentStoreJSTypes: cast=%#v doc=%#v", nUnvisitedCast.JSTypeBeforeCastObj(), nUnvisitedCast.JSDocObj())
	}
	// StandardColors singletons and staticStore template handles remain 100% valid across stores!
	if stdNumberColor.tcHandle.Load() != origColorHandle || nColor.ColorObj() != stdNumberColor {
		t.Fatalf("StandardColors singleton handle invalidated by ReleaseCurrentStoreJSTypes")
	}
	if staticTmpl.TypeOrColor() != staticSingleton {
		t.Fatalf("staticStore template handle invalidated by ReleaseCurrentStoreJSTypes")
	}

	// Verify GC actually reclaims the 5,000 released JSType objects!
	for range 5 {
		runtime.GC()
		if releasedTypesCollected.Load() > 0 {
			break
		}
	}
	if releasedTypesCollected.Load() == 0 {
		t.Fatalf("expected released JSType objects to be garbage collected after ReleaseCurrentStoreJSTypes")
	}

	// Verify a reset JSType used in a subsequent NodeStore gets a fresh valid handle, and StandardColors still work!
	sNext := NewNodeStore()
	SetCurrentNodeStore(sNext)
	nNextColor := AllocNodeToken(0, 39)
	nNextType := AllocNodeToken(0, 38)
	if !nNextColor.SetColorObj(stdNumberColor) || nNextColor.ColorObj() != stdNumberColor {
		t.Fatalf("StandardColors singleton failed in subsequent NodeStore")
	}
	if !nNextType.SetJSTypeObj(reusableType) || nNextType.JSTypeObj() != reusableType || reusableType.tcHandle.Load() == 0 {
		t.Fatalf("reusableType failed to acquire fresh handle in subsequent NodeStore")
	}
}

func TestWave27ASideTableMemoryFootprint(t *testing.T) {
	// Allocate 250,000 AST nodes with realistic Closure Compiler sparse ratios:
	// - 3% JSDoc (7,500 nodes)
	// - 0.3% JSTypeBeforeCast (750 nodes)
	// - 0.4% SCRIPT/FUNCTION ChangeTime (1,000 plain nodes -> inline in Payload!)
	// - 0.01% Length overflow >= 65535 (25 nodes)
	// Previously: 4 sparse pagedColumnFields + jsdocHandles allocated ~4.2 MB for 250K nodes (~400 MB at full scale).
	// Now: <= 600 KB of compact sharded aux maps!
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	const totalNodes = 250_000
	nodes := make([]Node, totalNodes)
	var alloc NodeAllocator
	alloc.InitLocal()
	for i := range nodes {
		nodes[i] = alloc.AllocNode1Cached(uint8(i%100 + 1))
	}
	alloc.Flush()

	doc := &fakeJSDocHolder{comment: "/** @param {string} x */"}
	cast := &fakeTCHolder{name: "UnknownType"}

	runtime.GC()
	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	for i, n := range nodes {
		if i%33 == 0 { // ~3% JSDoc
			n.SetJSDocObj(doc)
		}
		if i%333 == 0 { // ~0.3% @type cast
			n.SetJSTypeBeforeCastObj(cast)
		}
		if i%250 == 0 { // ~0.4% FUNCTION/SCRIPT ChangeTime
			n.SetChangeTime(int32(i + 1))
		}
		if i%10000 == 0 { // ~0.01% >64KB Length overflow
			n.SetLength(100_000 + int32(i))
		} else {
			n.SetLength(int32(i & 0x7FFF))
		}
	}

	runtime.GC()
	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	var sideTableBytes uint64
	if after.HeapAlloc > before.HeapAlloc {
		sideTableBytes = after.HeapAlloc - before.HeapAlloc
	}
	t.Logf("250,000-node sparse side-table live HeapAlloc delta: %d bytes (%.2f KiB), aux entries: %d",
		sideTableBytes, float64(sideTableBytes)/1024.0, s.auxLen())
	if sideTableBytes > 1024*1024 {
		t.Fatalf("side-table live memory too high: %d bytes (want <= 1 MiB)", sideTableBytes)
	}
	runtime.KeepAlive(s)
	runtime.KeepAlive(nodes)
}

func BenchmarkWave27ASparseFastPaths(b *testing.B) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)
	EnterSingleOwnerAlloc()
	defer ExitSingleOwnerAlloc()

	nUnset := AllocNodeToken(0, 40)
	doc := &fakeJSDocHolder{comment: "/** @const */"}
	nWithDoc := AllocNodeToken(0, 40)
	nWithDoc.SetJSDocObj(doc)

	b.Run("JSDocObjUnset", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if nUnset.JSDocObj() != nil {
				b.Fatal("unexpected")
			}
		}
	})
	b.Run("SetJSDocObjNilUnset", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nUnset.SetJSDocObj(nil)
		}
	})
	b.Run("JSDocObjSetSingleOwner", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if nWithDoc.JSDocObj() != doc {
				b.Fatal("unexpected")
			}
		}
	})
	b.Run("ChangeTimeUnset", func(b *testing.B) {
		var sink int32
		for i := 0; i < b.N; i++ {
			sink += nUnset.ChangeTime()
		}
		_ = sink
	})
	b.Run("SetChangeTimeZeroUnset", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nUnset.SetChangeTime(0)
		}
	})
}

func TestWave28ANoscanJSDocTableAndClearFastPaths(t *testing.T) {
	// 1. Verify layout & noscan invariants.
	if sz := reflect.TypeFor[jsdocSlot]().Size(); sz != 8 {
		t.Fatalf("sizeof(jsdocSlot) = %d, want 8", sz)
	}
	if hasGoPointers(reflect.TypeFor[jsdocSlot]()) {
		t.Fatalf("noscan violation: jsdocSlot contains Go pointers")
	}
	if sz := reflect.TypeFor[nodeAuxEntry]().Size(); sz != 24 {
		t.Fatalf("sizeof(nodeAuxEntry) = %d, want 24", sz)
	}

	// NilNode ClearJSDoc / ClearJSTypeBeforeCast must panic with runtime.Error (NPE).
	for _, tc := range []struct {
		name string
		fn   func()
	}{
		{"NilNode.ClearJSDoc", func() { NilNode.ClearJSDoc() }},
		{"NilNode.ClearJSTypeBeforeCast", func() { NilNode.ClearJSTypeBeforeCast() }},
	} {
		func() {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("%s: expected runtime.Error panic, got nil", tc.name)
				}
				if _, ok := r.(runtime.Error); !ok {
					t.Fatalf("%s: expected runtime.Error, got %T (%v)", tc.name, r, r)
				}
			}()
			tc.fn()
		}()
	}

	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer func() {
		ExitSingleOwnerAlloc()
		SetCurrentNodeStore(prev)
	}()
	EnterSingleOwnerAlloc()

	// 2. Zero-alloc ClearJSDoc / ClearJSTypeBeforeCast on unset & set nodes + coexistence.
	nUnset := AllocNodeToken(0, 40)
	docA := &fakeJSDocHolder{comment: "/** @type {string} */"}
	docB := &fakeJSDocHolder{comment: "/** @type {number} */"}
	castA := &fakeTCHolder{name: "StringType"}
	castB := &fakeTCHolder{name: "NumberType"}

	allocsUnset := testing.AllocsPerRun(200, func() {
		nUnset.ClearJSDoc()
		nUnset.ClearJSTypeBeforeCast()
	})
	if allocsUnset > 0 {
		t.Fatalf("ClearJSDoc/ClearJSTypeBeforeCast on unset node allocated %.1f objects", allocsUnset)
	}

	// Overwrite on same node reuses handle cleanly.
	nOverwrite := AllocNodeToken(0, 40)
	nOverwrite.SetJSDocObj(docA)
	countAfterFirst := s.jsdocs.count.Load()
	nOverwrite.SetJSDocObj(docB)
	if nOverwrite.JSDocObj() != docB {
		t.Fatalf("overwrite JSDocObj = %#v, want docB", nOverwrite.JSDocObj())
	}
	if s.jsdocs.count.Load() != countAfterFirst {
		t.Fatalf("overwrite allocated extra handle: before=%d after=%d", countAfterFirst, s.jsdocs.count.Load())
	}
	nOverwrite.ClearJSDoc()
	if nOverwrite.JSDocObj() != nil || nOverwrite.Rec().KindFlags&kindFlagHasJSDoc != 0 {
		t.Fatalf("after ClearJSDoc: doc=%#v flags=0x%x", nOverwrite.JSDocObj(), nOverwrite.Rec().KindFlags)
	}

	// Shared kindFlagHasJSDocOrCast bit with ClearJSDoc & ClearJSTypeBeforeCast in both orders.
	nCoexist := AllocNodeToken(0, 42)
	nCoexist.SetJSDocObj(docA)
	nCoexist.SetJSTypeBeforeCastObj(castA)
	nCoexist.ClearJSTypeBeforeCast()
	if nCoexist.JSTypeBeforeCastObj() != nil || nCoexist.JSDocObj() != docA || nCoexist.Rec().KindFlags&kindFlagHasJSDocOrCast == 0 {
		t.Fatalf("ClearJSTypeBeforeCast disturbed JSDocObj: doc=%#v cast=%#v", nCoexist.JSDocObj(), nCoexist.JSTypeBeforeCastObj())
	}
	nCoexist.SetJSTypeBeforeCastObj(castB)
	nCoexist.ClearJSDoc()
	if nCoexist.JSDocObj() != nil || nCoexist.JSTypeBeforeCastObj() != castB || nCoexist.Rec().KindFlags&kindFlagHasJSDocOrCast == 0 {
		t.Fatalf("ClearJSDoc disturbed JSTypeBeforeCastObj: doc=%#v cast=%#v", nCoexist.JSDocObj(), nCoexist.JSTypeBeforeCastObj())
	}
	nCoexist.ClearJSTypeBeforeCast()
	if nCoexist.JSDocObj() != nil || nCoexist.JSTypeBeforeCastObj() != nil || nCoexist.Rec().KindFlags&kindFlagHasJSDocOrCast != 0 {
		t.Fatalf("clearing both left residual state: flags=0x%x", nCoexist.Rec().KindFlags)
	}

	// 3. Tombstone creation, tombstone revival, tombstone reuse, and bounded rehash under heavy churn
	// within a single shard (spanning >4096 s.jsdocs handles across AnyTable blocks as well).
	var shard0Nodes []Node
	for len(shard0Nodes) < 300 {
		cand := AllocNodeToken(0, 40)
		if jsdocShardIndex(uint32(cand)) == 0 {
			shard0Nodes = append(shard0Nodes, cand)
		}
	}
	// Churn insert+delete on first 200 nodes in shard 0 to exercise tombstone compaction on rehash
	// without unbounded slot capacity blowup.
	for round := 0; round < 20; round++ {
		for i := 0; i < 12; i++ {
			idx := (round*12 + i) % 200
			shard0Nodes[idx].SetJSDocObj(docA)
		}
		for i := 0; i < 12; i++ {
			idx := (round*12 + i) % 200
			shard0Nodes[idx].ClearJSDoc()
		}
	}
	if cap0 := len(s.jsdocShards[0].slots); cap0 > 64 {
		t.Fatalf("expected tombstone compaction to keep shard 0 slots bounded (<= 64), got %d", cap0)
	}
	if s.jsdocShards[0].count != 0 {
		t.Fatalf("expected shard 0 count == 0 after churn clear, got %d", s.jsdocShards[0].count)
	}

	// Revive previously tombstoned keys alongside fresh keys across >5000 total JSDoc entries (multi-block s.jsdocs).
	const multiCount = 5200
	multiNodes := make([]Node, multiCount)
	multiDocs := make([]*fakeJSDocHolder, 32)
	for i := range multiDocs {
		multiDocs[i] = &fakeJSDocHolder{comment: fmt.Sprintf("/** doc %d */", i)}
	}
	copy(multiNodes, shard0Nodes)
	for i := len(shard0Nodes); i < multiCount; i++ {
		multiNodes[i] = AllocNodeToken(0, 40)
	}
	for i, n := range multiNodes {
		n.SetJSDocObj(multiDocs[i&31])
	}
	// Delete every 3rd node, overwrite every 5th surviving node, then insert more to trigger tombstone-dropping rehash.
	for i, n := range multiNodes {
		if i%3 == 0 {
			n.ClearJSDoc()
		} else if i%5 == 0 {
			n.SetJSDocObj(docB)
		}
	}
	for i, n := range multiNodes {
		got := n.JSDocObj()
		if i%3 == 0 {
			if got != nil {
				t.Fatalf("deleted multiNodes[%d] returned %#v", i, got)
			}
		} else if i%5 == 0 {
			if got != docB {
				t.Fatalf("overwritten multiNodes[%d] = %#v, want docB", i, got)
			}
		} else if got != multiDocs[i&31] {
			t.Fatalf("multiNodes[%d] = %#v, want %#v", i, got, multiDocs[i&31])
		}
	}

	// 4. ClonePropsFrom & GetAux/SetAux interoperability.
	srcClone := AllocNodeToken(0, 40)
	srcClone.SetSourceFile(7)
	srcClone.SetPropFlags(0x1234)
	srcClone.SetJSDocObj(docA)
	srcClone.SetJSTypeBeforeCastObj(castA)
	dstClone := AllocNodeToken(0, 40)
	if !dstClone.ClonePropsFrom(srcClone) {
		t.Fatal("ClonePropsFrom failed")
	}
	if dstClone.JSDocObj() != docA || dstClone.JSTypeBeforeCastObj() != castA || dstClone.SourceFile() != 7 || dstClone.PropFlags() != 0x1234 {
		t.Fatalf("ClonePropsFrom mismatch on dstClone")
	}
	// Clearing JSDoc on srcClone must NOT affect dstClone's independent handle!
	srcClone.ClearJSDoc()
	if dstClone.JSDocObj() != docA || srcClone.JSDocObj() != nil {
		t.Fatalf("ClearJSDoc on srcClone affected dstClone: dst=%#v src=%#v", dstClone.JSDocObj(), srcClone.JSDocObj())
	}

	// 5. Concurrent multi-goroutine -race stress with SuspendSingleOwnerAlloc / ResumeSingleOwnerAlloc.
	const workers = 12
	const perWorker = 400
	workerNodes := make([][]Node, workers)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		w := w
		SuspendSingleOwnerAlloc()
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer ResumeSingleOwnerAlloc()
			var a NodeAllocator
			a.InitLocal()
			ns := make([]Node, perWorker)
			myDoc := multiDocs[w&31]
			for i := 0; i < perWorker; i++ {
				n := a.AllocNode1Cached(40)
				n.SetJSDocObj(myDoc)
				if i%4 == 0 {
					n.SetJSTypeBeforeCastObj(castA)
				}
				if i%6 == 0 {
					n.SetJSDocObj(docB)
				}
				if i%9 == 0 {
					n.ClearJSDoc()
				}
				if i%8 == 0 {
					n.ClearJSTypeBeforeCast()
				}
				if i%11 == 0 {
					aux, _ := n.GetAux()
					aux.Extra = w*1000 + i
					n.SetAux(aux)
				}
				ns[i] = n
			}
			a.Flush()
			workerNodes[w] = ns
		}()
	}
	// Concurrently read pre-existing multiNodes on parent goroutine while workers mutate shards!
	for i := 0; i < 500; i++ {
		n := multiNodes[i]
		_ = n.JSDocObj()
	}
	wg.Wait()

	// Back in single-owner fast path: verify worker-written state lock-free!
	if !isSingleOwnerFastPath() {
		t.Fatal("expected single-owner fast path resumed after workers completed")
	}
	for w := 0; w < workers; w++ {
		myDoc := multiDocs[w&31]
		for i, n := range workerNodes[w] {
			var wantDoc any = myDoc
			if i%6 == 0 {
				wantDoc = docB
			}
			if i%9 == 0 {
				wantDoc = nil
			}
			if got := n.JSDocObj(); got != wantDoc {
				t.Fatalf("worker[%d][%d] JSDocObj = %#v, want %#v", w, i, got, wantDoc)
			}
			var wantCast any
			if i%4 == 0 && i%8 != 0 {
				wantCast = castA
			}
			if got := n.JSTypeBeforeCastObj(); got != wantCast {
				t.Fatalf("worker[%d][%d] JSTypeBeforeCastObj = %#v, want %#v", w, i, got, wantCast)
			}
			if i%11 == 0 {
				aux, ok := n.GetAux()
				if !ok || aux.Extra != w*1000+i || aux.JSDoc != wantDoc {
					t.Fatalf("worker[%d][%d] GetAux = (%+v, %v), want Extra=%d JSDoc=%#v", w, i, aux, ok, w*1000+i, wantDoc)
				}
			}
		}
	}
}

func TestWave30CNoscanAuxTableAndCastBloomFastPath(t *testing.T) {
	if sz := reflect.TypeFor[auxSlot]().Size(); sz != 16 {
		t.Fatalf("sizeof(auxSlot) = %d, want 16", sz)
	}
	if hasGoPointers(reflect.TypeFor[auxSlot]()) {
		t.Fatalf("noscan violation: auxSlot contains Go pointers")
	}

	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer func() {
		ExitSingleOwnerAlloc()
		SetCurrentNodeStore(prev)
	}()
	EnterSingleOwnerAlloc()

	doc := &fakeJSDocHolder{comment: "/** @param {string} x */"}
	cast1 := &fakeTCHolder{name: "StringType"}
	cast2 := &fakeTCHolder{name: "NumberType"}

	// Seed 10,000 JSDoc-only nodes and 32 @type cast nodes (realistic CHECKS ratio),
	// then run ConvertTypesToColors-style ClearJSTypeBeforeCast() across all 10,000 nodes!
	const nNodes = 10_000
	nodes := make([]Node, nNodes)
	for i := range nodes {
		n := AllocNodeToken(0, 40)
		n.SetJSDocObj(doc)
		if i%313 == 0 {
			n.SetJSTypeBeforeCastObj(cast1)
		}
		nodes[i] = n
	}

	// Verify cast handle reuse across 5,000 set/clear cycles does NOT grow s.casts past peak live count!
	peakCasts := s.casts.count.Load()
	nChurn := nodes[0]
	for i := 0; i < 5000; i++ {
		nChurn.ClearJSTypeBeforeCast()
		if i&1 == 0 {
			nChurn.SetJSTypeBeforeCastObj(cast2)
		} else {
			nChurn.SetJSTypeBeforeCastObj(cast1)
		}
	}
	if got := s.casts.count.Load(); got != peakCasts {
		t.Fatalf("expected cast handle free-list reuse to keep s.casts.count == %d, got %d", peakCasts, got)
	}

	// ConvertTypesToColors pass: ClearJSTypeBeforeCast across all 10,000 JSDoc-bearing nodes.
	allocs := testing.AllocsPerRun(5, func() {
		for _, n := range nodes {
			n.ClearJSTypeBeforeCast()
		}
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs on ClearJSTypeBeforeCast sweep, got %v", allocs)
	}
	for i, n := range nodes {
		if n.JSTypeBeforeCastObj() != nil {
			t.Fatalf("node[%d] still has JSTypeBeforeCastObj", i)
		}
		if n.JSDocObj() != doc {
			t.Fatalf("node[%d] lost JSDocObj after ClearJSTypeBeforeCast", i)
		}
	}
	for shardIdx := range s.aux {
		if c := s.aux[shardIdx].castCount.Load(); c != 0 {
			t.Fatalf("shard %d castCount = %d after clearing all casts, want 0", shardIdx, c)
		}
		if bl := s.aux[shardIdx].castBloom.Load(); bl != 0 {
			t.Fatalf("shard %d castBloom = 0x%x after clearing all casts, want 0", shardIdx, bl)
		}
	}
}

func BenchmarkWave30CTypeBeforeCastHotPaths(b *testing.B) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer func() {
		ExitSingleOwnerAlloc()
		SetCurrentNodeStore(prev)
	}()
	EnterSingleOwnerAlloc()

	doc := &fakeJSDocHolder{comment: "/** @type {string} */"}
	cast1 := &fakeTCHolder{name: "StringType"}
	cast2 := &fakeTCHolder{name: "NumberType"}

	// Seed realistic sparse casts across shards so s.hasExtCast == true and some Bloom bits are active.
	for i := 0; i < 64; i++ {
		cNode := AllocNodeToken(0, 101)
		cNode.SetJSTypeBeforeCastObj(cast1)
	}

	jsdocNodes := make([]Node, 1024)
	for i := range jsdocNodes {
		n := AllocNodeToken(0, 40)
		n.SetJSDocObj(doc)
		jsdocNodes[i] = n
	}
	castNode := AllocNodeToken(0, 101)
	castNode.SetJSTypeBeforeCastObj(cast1)

	b.Run("ClearJSTypeBeforeCastOnJSDocNode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			jsdocNodes[i&1023].ClearJSTypeBeforeCast()
		}
	})
	b.Run("GetJSTypeBeforeCastOnJSDocNode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if jsdocNodes[i&1023].JSTypeBeforeCastObj() != nil {
				b.Fatal("unexpected")
			}
		}
	})
	b.Run("SetJSTypeBeforeCastOverwrite", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if i&1 == 0 {
				castNode.SetJSTypeBeforeCastObj(cast2)
			} else {
				castNode.SetJSTypeBeforeCastObj(cast1)
			}
		}
	})
	b.Run("SetJSTypeObjRegistered", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if i&1 == 0 {
				jsdocNodes[i&1023].SetJSTypeObj(cast1)
			} else {
				jsdocNodes[i&1023].SetJSTypeObj(cast2)
			}
		}
	})
}

func TestWave32CCastBitBlockAndInlineSetJSTypeFastPaths(t *testing.T) {
	if sz := reflect.TypeFor[castBitBlock]().Size(); sz != 2048 {
		t.Fatalf("sizeof(castBitBlock) = %d, want 2048", sz)
	}
	if hasGoPointers(reflect.TypeFor[castBitBlock]()) {
		t.Fatalf("noscan violation: castBitBlock contains Go pointers")
	}

	s1 := NewNodeStore()
	s2 := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer func() {
		ExitSingleOwnerAlloc()
		SetCurrentNodeStore(prev)
	}()
	EnterSingleOwnerAlloc()

	doc := &fakeJSDocHolder{comment: "/** @type {string} */"}
	castType := &fakeTCHolder{name: "StringType"}
	castColor := &fakeColorHolder{name: "StringColor"}
	var typedNilType *fakeTCHolder

	nDocOnly := AllocNodeToken(0, 40)
	nDocOnly.SetJSDocObj(doc)
	nCast := AllocNodeToken(0, 101)
	nCast.SetJSDocObj(doc)
	nCast.SetJSTypeBeforeCastObj(castType)
	nColorCast := AllocNodeToken(0, 101)
	nColorCast.SetJSTypeBeforeCastObj(castColor)

	// Zero-alloc SetJSTypeBeforeCastObj(nil) on JSDoc-bearing node in same block as active casts.
	allocs := testing.AllocsPerRun(100, func() {
		nDocOnly.SetJSTypeBeforeCastObj(nil)
		nDocOnly.ClearJSTypeBeforeCast()
		if nDocOnly.JSTypeBeforeCastObj() != nil {
			t.Fatal("unexpected cast on JSDoc-only node")
		}
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs on JSDoc-only SetJSTypeBeforeCastObj(nil), got %v", allocs)
	}

	// Store switch unregister/register preserves castBitField isolation.
	ExitSingleOwnerAlloc()
	SetCurrentNodeStore(s2)
	nS2 := AllocNodeToken(0, 40)
	nS2.SetJSDocObj(doc)
	if nS2.JSTypeBeforeCastObj() != nil {
		t.Fatal("unexpected cast in s2")
	}
	SetCurrentNodeStore(s1)
	EnterSingleOwnerAlloc()
	if nCast.JSTypeBeforeCastObj() != castType || nColorCast.JSTypeBeforeCastObj() != castColor {
		t.Fatal("castBitField lost across store switch")
	}

	// ReleaseCurrentStoreJSTypes clears non-Color cast bits while preserving ColorMarker cast bits.
	ReleaseCurrentStoreJSTypes()
	if nCast.JSTypeBeforeCastObj() != nil {
		t.Fatalf("expected non-Color cast cleared by ReleaseCurrentStoreJSTypes, got %#v", nCast.JSTypeBeforeCastObj())
	}
	if nColorCast.JSTypeBeforeCastObj() != castColor {
		t.Fatalf("expected Color cast preserved by ReleaseCurrentStoreJSTypes, got %#v", nColorCast.JSTypeBeforeCastObj())
	}

	// Typed-nil clears cast cleanly and resets cast bit.
	nCast.SetJSTypeBeforeCastObj(castType)
	nCast.SetJSTypeBeforeCastObj(any(typedNilType))
	if nCast.JSTypeBeforeCastObj() != nil || nCast.JSDocObj() != doc {
		t.Fatalf("typed-nil SetJSTypeBeforeCastObj mismatch")
	}

	// Concurrent -race stress on SetJSTypeObj / SetColorObj / SetTypeOrColor / SetJSTypeBeforeCastObj.
	ExitSingleOwnerAlloc()
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		g := g
		wg.Add(1)
		go func() {
			defer wg.Done()
			alloc := NewNodeAllocator(s1)
			t1 := &fakeTCHolder{name: fmt.Sprintf("T1_%d", g)}
			t2 := &fakeTCHolder{name: fmt.Sprintf("T2_%d", g)}
			c1 := &fakeColorHolder{name: fmt.Sprintf("C1_%d", g)}
			for i := 0; i < 300; i++ {
				n := alloc.AllocWithToken(40)
				n.SetJSDocObj(doc)
				if !n.SetJSTypeObj(t1) || !n.SetJSTypeObj(t2) || n.JSTypeObj() != t2 {
					t.Errorf("SetJSTypeObj failed")
					return
				}
				if !n.SetJSTypeObj(nil) || n.JSTypeObj() != nil {
					t.Errorf("SetJSTypeObj(nil) failed")
					return
				}
				if !n.SetColorObj(c1) || n.SetJSTypeObj(t1) || n.ColorObj() != c1 {
					t.Errorf("SetColorObj guard failed")
					return
				}
				n.SetTypeOrColor(t1)
				if n.JSTypeObj() != t1 {
					t.Errorf("SetTypeOrColor(t1) failed")
					return
				}
				if i%5 == 0 {
					n.SetJSTypeBeforeCastObj(t2)
					if n.JSTypeBeforeCastObj() != t2 {
						t.Errorf("JSTypeBeforeCastObj mismatch")
						return
					}
					n.SetJSTypeBeforeCastObj(nil)
				} else {
					n.SetJSTypeBeforeCastObj(nil)
					n.ClearJSTypeBeforeCast()
				}
			}
			alloc.Flush()
		}()
	}
	wg.Wait()
}


func TestWave33CSetJSTypeBeforeCastFastNoopAndLengthBoundaries(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer func() {
		ExitSingleOwnerAlloc()
		SetCurrentNodeStore(prev)
	}()
	EnterSingleOwnerAlloc()

	var typedNilTC *fakeTCHolder
	var typedNilJN *fakeJavaNilOnly
	realCast1 := &fakeTCHolder{name: "NumberType"}
	realCast2 := &fakeTCHolder{name: "StringType"}
	doc := &fakeJSDocHolder{comment: "/** @type {string} */"}

	// 1. Fresh block with nil castBitField: nil, typedNilTC, typedNilJN, and ClearJSTypeBeforeCast
	// on BOTH plain nodes AND JSDoc-bearing nodes must be 0-alloc no-ops that never allocate castBitField
	// and never clobber kindFlagHasJSDoc!
	nPlain := AllocNodeInit(0)
	nWithDoc := AllocNodeToken(3, 32)
	nWithDoc.SetJSDocObj(doc)

	for _, n := range []Node{nPlain, nWithDoc} {
		n.SetJSTypeBeforeCastObj(nil)
		n.SetJSTypeBeforeCastObj(any(typedNilTC))
		n.SetJSTypeBeforeCastObj(any(typedNilJN))
		n.ClearJSTypeBeforeCast()
	}
	if castBitField[nPlain>>blockShift].Load() != nil {
		t.Fatalf("expected castBitField to remain nil after nil/typed-nil/ClearJSTypeBeforeCast on unset nodes")
	}
	if nWithDoc.JSDocObj() != doc {
		t.Fatalf("typed-nil SetJSTypeBeforeCastObj clobbered JSDocObj: got %#v", nWithDoc.JSDocObj())
	}
	if allocs := testing.AllocsPerRun(200, func() {
		nWithDoc.SetJSTypeBeforeCastObj(nil)
		nWithDoc.SetJSTypeBeforeCastObj(any(typedNilTC))
		nWithDoc.SetJSTypeBeforeCastObj(any(typedNilJN))
		nWithDoc.ClearJSTypeBeforeCast()
	}); allocs != 0 {
		t.Fatalf("expected 0 allocs on nil/typed-nil/clear SetJSTypeBeforeCastObj, got %v", allocs)
	}

	// 2. Non-nil set -> overwrite -> typed-nil clear (preserving JSDoc) -> non-nil set -> ClearJSTypeBeforeCast.
	nWithDoc.SetJSTypeBeforeCastObj(realCast1)
	if nWithDoc.JSTypeBeforeCastObj() != realCast1 || nWithDoc.JSDocObj() != doc {
		t.Fatalf("expected realCast1 + doc, got cast=%#v doc=%#v", nWithDoc.JSTypeBeforeCastObj(), nWithDoc.JSDocObj())
	}
	nWithDoc.SetJSTypeBeforeCastObj(realCast2)
	if nWithDoc.JSTypeBeforeCastObj() != realCast2 {
		t.Fatalf("overwrite failed: got %#v", nWithDoc.JSTypeBeforeCastObj())
	}
	nWithDoc.SetJSTypeBeforeCastObj(any(typedNilTC))
	if nWithDoc.JSTypeBeforeCastObj() != nil || nWithDoc.JSDocObj() != doc {
		t.Fatalf("typedNilTC clear failed or lost JSDoc: cast=%#v doc=%#v", nWithDoc.JSTypeBeforeCastObj(), nWithDoc.JSDocObj())
	}
	nWithDoc.SetJSTypeBeforeCastObj(realCast1)
	nWithDoc.SetJSTypeBeforeCastObj(any(typedNilJN))
	if nWithDoc.JSTypeBeforeCastObj() != nil || nWithDoc.JSDocObj() != doc {
		t.Fatalf("typedNilJN clear failed or lost JSDoc: cast=%#v doc=%#v", nWithDoc.JSTypeBeforeCastObj(), nWithDoc.JSDocObj())
	}
	nPlain.SetJSTypeBeforeCastObj(realCast1)
	nPlain.ClearJSTypeBeforeCast()
	if nPlain.JSTypeBeforeCastObj() != nil || nPlain.Rec().KindFlags&kindFlagHasTypeBeforeCast != 0 {
		t.Fatalf("ClearJSTypeBeforeCast failed to clear flag on doc-free node")
	}

	// 3. Length / SetLength boundaries, including complementary bit patterns (0xFF00 -> 0x00FF -> 0xFFFF -> 0 -> -5).
	lenSeq := []int32{0, 1, 0xFF00, 0x00FF, 0xAAAA, 0x5555, 65534, 65535, 65536, 200000, 0x00FF, 0xFF00, -1, -99999, 42, 0}
	for _, want := range lenSeq {
		nPlain.SetLength(want)
		if got := nPlain.Length(); got != want {
			t.Fatalf("SetLength(%d) -> Length() = %d", want, got)
		}
	}

	// 4. Concurrent -race stress on SetJSTypeBeforeCastObj (nil, typed-nil, non-nil, overwrite, clear) & Length.
	ExitSingleOwnerAlloc()
	const workers = 12
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		w := w
		go func() {
			defer wg.Done()
			var la NodeAllocator
			la.InitLocal()
			defer la.Flush()
			cA := &fakeTCHolder{name: fmt.Sprintf("CastA_%d", w)}
			cB := &fakeTCHolder{name: fmt.Sprintf("CastB_%d", w)}
			for i := 0; i < 400; i++ {
				n := la.AllocNode1Cached(32)
				if i%2 == 0 {
					n.SetJSDocObj(doc)
				}
				// Unset no-ops:
				n.SetJSTypeBeforeCastObj(nil)
				n.SetJSTypeBeforeCastObj(any(typedNilTC))
				n.SetJSTypeBeforeCastObj(any(typedNilJN))
				n.ClearJSTypeBeforeCast()
				if i%4 == 0 {
					n.SetJSTypeBeforeCastObj(cA)
					if n.JSTypeBeforeCastObj() != cA {
						t.Errorf("expected cA")
						return
					}
					n.SetJSTypeBeforeCastObj(cB)
					if n.JSTypeBeforeCastObj() != cB {
						t.Errorf("expected cB")
						return
					}
					switch i % 12 {
					case 0:
						n.SetJSTypeBeforeCastObj(nil)
					case 4:
						n.SetJSTypeBeforeCastObj(any(typedNilTC))
					case 8:
						n.ClearJSTypeBeforeCast()
					}
					if n.JSTypeBeforeCastObj() != nil {
						t.Errorf("expected nil after clear")
						return
					}
				}
				ln := lenSeq[i%len(lenSeq)]
				n.SetLength(ln)
				if n.Length() != ln {
					t.Errorf("Length mismatch: got %d, want %d", n.Length(), ln)
					return
				}
			}
		}()
	}
	wg.Wait()
}


func TestWave38CNodeData32ByteLayoutSourceFileAndLazyLength(t *testing.T) {
	if sz := reflect.TypeFor[NodeData]().Size(); sz != 32 {
		t.Fatalf("reflect.TypeFor[NodeData]().Size() = %d, want 32", sz)
	}
	if sz := reflect.TypeFor[NodeBlock]().Size(); sz != 512*1024 {
		t.Fatalf("reflect.TypeFor[NodeBlock]().Size() = %d, want %d", sz, 512*1024)
	}
	if hasGoPointers(reflect.TypeFor[NodeData]()) {
		t.Fatal("NodeData contains Go pointers")
	}

	s1 := NewNodeStore()
	s2 := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer SetCurrentNodeStore(prev)

	// 1. Lazy Length: allocating nodes and calling SetLength(0) on unallocated pages must allocate 0 len16Pages!
	nLazy := AllocNodeToken(0, 40)
	if nLazy.Length() != 0 {
		t.Fatalf("fresh Length() = %d, want 0", nLazy.Length())
	}
	nLazy.SetLength(0)
	if len(s1.len16Pages) != 0 || len(s1.lengthOverflows) != 0 {
		t.Fatalf("SetLength(0) allocated pages eagerly: len16Pages=%d lengthOverflows=%d",
			len(s1.len16Pages), len(s1.lengthOverflows))
	}
	// Node(0).SetLength(0) must still fault with runtime.Error!
	func() {
		defer func() {
			r := recover()
			if _, ok := r.(runtime.Error); !ok {
				t.Fatalf("Node(0).SetLength(0): expected runtime.Error, got %T (%v)", r, r)
			}
		}()
		Node(0).SetLength(0)
	}()

	// 2. Length across 0, 1..65534, >= 65535, and < 0.
	lengths := []int32{0, 1, 255, 0xFF00, 0x00FF, 65534, 65535, 65536, 500_000, math.MaxInt32, -1, -42, math.MinInt32, 1234, 0}
	for _, want := range lengths {
		nLazy.SetLength(want)
		if got := nLazy.Length(); got != want {
			t.Fatalf("SetLength(%d) -> Length() = %d", want, got)
		}
	}

	// 3. SourceFile < 65535 inline (0 sfOverflows) vs >= 65535 overflow, ClonePropsFrom, and NodeAllocator!
	nSF := AllocNodeToken(0, 41)
	for _, sf := range []uint32{0, 1, 19430, 65534} {
		nSF.SetSourceFile(sf)
		if got := nSF.SourceFile(); got != sf {
			t.Fatalf("inline SetSourceFile(%d) -> %d", sf, got)
		}
	}
	if len(s1.sfOverflows) != 0 {
		t.Fatalf("inline SourceFile (< 65535) allocated sfOverflows: %d", len(s1.sfOverflows))
	}

	for _, sf := range []uint32{65535, 65536, 98765, 0xCAFEBABE, 0xFFFFFFFF, 42, 0, 77777} {
		nSF.SetSourceFile(sf)
		if got := nSF.SourceFile(); got != sf {
			t.Fatalf("overflow SetSourceFile(%d) -> %d", sf, got)
		}
	}
	nSFClone := AllocNodeToken(0, 42)
	if !nSFClone.ClonePropsFrom(nSF) || nSFClone.SourceFile() != 77777 {
		t.Fatalf("ClonePropsFrom overflow SourceFile = %d, want 77777", nSFClone.SourceFile())
	}

	// NodeAllocator with tmplSourceFile < 65535 (inline) AND >= 65535 (overflow) across all AllocNode* methods!
	for _, tmplSF := range []uint32{19430, 65535, 123456} {
		tmpl := AllocNodeToken(0, 94)
		tmpl.SetSourceFile(tmplSF)
		tmpl.SetPropFlags(0x11223344)

		var a NodeAllocator
		a.InitLocal()
		a.SetTemplate(tmpl)
		sid := a.Intern("wave38c_ident")

		c1 := a.AllocNode1Cached(10)
		c2 := a.AllocNode1(11, tmpl)
		c3 := a.AllocStringNodeCached(38, sid)
		c4 := a.AllocStringNode(38, sid, tmpl)
		c5 := a.AllocNumberNode(39, 3.5, tmpl)
		p2 := a.AllocNode2(20, c1, tmpl)
		p3 := a.AllocNode3(21, c2, c3, tmpl)
		p4 := a.AllocNode4(22, c4, c5, p2, tmpl)
		a.Flush()

		for idx, n := range []Node{c1, c2, c3, c4, c5, p2, p3, p4} {
			if n.SourceFile() != tmplSF || n.PropFlags() != 0x11223344 {
				t.Fatalf("tmplSF=%d node[%d]: SourceFile()=%d PropFlags()=0x%X",
					tmplSF, idx, n.SourceFile(), n.PropFlags())
			}
		}
	}

	// 4. Store switch unregister/re-register clears and restores both len16Top and sfOverflowTop!
	nLazy.SetLength(99999)
	nSF.SetSourceFile(88888)
	SetCurrentNodeStore(s2)
	if nLazy.Length() != 0 {
		t.Fatalf("unregistered s1 Length visible in s2: %d", nLazy.Length())
	}
	if getSourceFileOverflowSlow(nSF) != 0 {
		t.Fatalf("unregistered s1 sfOverflowTop visible in s2: %d", getSourceFileOverflowSlow(nSF))
	}
	SetCurrentNodeStore(s1)
	if nLazy.Length() != 99999 || nSF.SourceFile() != 88888 {
		t.Fatalf("re-registered s1 lost Length (%d) or SourceFile (%d)", nLazy.Length(), nSF.SourceFile())
	}
}

func TestWave39CHasNonNegativeLengthAnd16BlockSlab(t *testing.T) {
	s1 := NewNodeStore()
	s2 := NewNodeStore()
	prev := SetCurrentNodeStore(s1)
	defer SetCurrentNodeStore(prev)

	baseNeg := negLenCount.Load()

	// 1. NilNode.HasNonNegativeLength() must panic with runtime.Error both when negLenCount == 0 and > 0.
	assertNilPanics := func(label string) {
		t.Helper()
		defer func() {
			r := recover()
			if _, ok := r.(runtime.Error); !ok {
				t.Fatalf("%s: NilNode.HasNonNegativeLength() expected runtime.Error, got %T (%v)", label, r, r)
			}
		}()
		_ = NilNode.HasNonNegativeLength()
	}
	assertNilPanics("fast-path")

	// 2. Non-negative lengths (unset, 0, 1..65534, >= 65535) keep negLenCount unchanged and return true with 0 allocs.
	nA := AllocNodeToken(0, 40)
	nB := AllocNodeToken(0, 41)
	nC := AllocNodeToken(0, 42)
	for _, pos := range []int32{0, 1, 65534, 65535, 100_000, math.MaxInt32} {
		nA.SetLength(pos)
		if !nA.HasNonNegativeLength() {
			t.Fatalf("HasNonNegativeLength() = false for length %d", pos)
		}
		if got := negLenCount.Load(); got != baseNeg {
			t.Fatalf("positive SetLength(%d) changed negLenCount: got %d, want %d", pos, got, baseNeg)
		}
	}
	if allocs := testing.AllocsPerRun(100, func() {
		if !nA.HasNonNegativeLength() || !nB.HasNonNegativeLength() {
			t.Fatal("unexpected false")
		}
	}); allocs != 0 {
		t.Fatalf("expected 0 allocs on HasNonNegativeLength, got %v", allocs)
	}

	// 3. Negative length transitions: inline -> neg, neg -> neg, overflow+ -> neg, neg -> inline+, neg -> overflow+.
	nA.SetLength(50)
	nA.SetLength(-1)
	nB.SetLength(-42)
	nC.SetLength(100_000)
	nC.SetLength(math.MinInt32)
	if got := negLenCount.Load(); got != baseNeg+3 {
		t.Fatalf("negLenCount = %d, want %d", got, baseNeg+3)
	}
	assertNilPanics("slow-path")

	// Overwriting negative with another negative must not double-increment.
	nB.SetLength(-99)
	if got := negLenCount.Load(); got != baseNeg+3 {
		t.Fatalf("neg->neg overwrite changed negLenCount: got %d, want %d", got, baseNeg+3)
	}
	nPosInline := AllocNodeToken(0, 43)
	nPosInline.SetLength(123)
	nPosOverflow := AllocNodeToken(0, 44)
	nPosOverflow.SetLength(200_000)
	if nA.HasNonNegativeLength() || nB.HasNonNegativeLength() || nC.HasNonNegativeLength() {
		t.Fatal("expected HasNonNegativeLength() == false for negative-length nodes")
	}
	if !nPosInline.HasNonNegativeLength() || !nPosOverflow.HasNonNegativeLength() {
		t.Fatal("expected HasNonNegativeLength() == true for non-negative nodes while negLenCount > 0")
	}

	// Store switch s1 -> s2 must subtract s1.negLenNodes so s2 runs with negLenCount == baseNeg!
	SetCurrentNodeStore(s2)
	if got := negLenCount.Load(); got != baseNeg {
		t.Fatalf("switching to clean s2 left negLenCount = %d, want %d", got, baseNeg)
	}
	nS2 := AllocNodeToken(0, 45)
	if !nS2.HasNonNegativeLength() {
		t.Fatal("expected s2 node HasNonNegativeLength() == true")
	}

	// Switching back to s1 restores s1.negLenNodes in negLenCount!
	SetCurrentNodeStore(s1)
	if got := negLenCount.Load(); got != baseNeg+3 {
		t.Fatalf("re-registering s1 gave negLenCount = %d, want %d", got, baseNeg+3)
	}
	if nA.HasNonNegativeLength() || nB.HasNonNegativeLength() || nC.HasNonNegativeLength() {
		t.Fatal("re-registered s1 lost negative length state")
	}

	// Transition all 3 negative nodes back to non-negative (0, inline+, overflow+).
	nA.SetLength(0)
	nB.SetLength(500)
	nC.SetLength(100_000)
	if got := negLenCount.Load(); got != baseNeg {
		t.Fatalf("clearing negative lengths gave negLenCount = %d, want %d", got, baseNeg)
	}
	if !nA.HasNonNegativeLength() || !nB.HasNonNegativeLength() || !nC.HasNonNegativeLength() {
		t.Fatal("expected HasNonNegativeLength() == true after clearing negative lengths")
	}

	// Static template node negative length persists across ResetNodeStore and clears cleanly.
	tmpl := NewTemplateNode(55)
	tmpl.SetLength(-7)
	if tmpl.HasNonNegativeLength() || negLenCount.Load() != baseNeg+1 {
		t.Fatalf("template negative length tracking failed: count=%d", negLenCount.Load())
	}
	ResetNodeStore()
	if tmpl.HasNonNegativeLength() || negLenCount.Load() != baseNeg+1 {
		t.Fatalf("ResetNodeStore disturbed static template negLenCount: count=%d", negLenCount.Load())
	}
	tmpl.SetLength(0)
	if !tmpl.HasNonNegativeLength() || negLenCount.Load() != baseNeg {
		t.Fatalf("clearing static template negative length failed: count=%d", negLenCount.Load())
	}

	// 4. Verify [4]nodeBlock (2 MiB) slab batching progression on a fresh store:
	// blkIdx < 2 -> 1 block (512 KiB), blkIdx >= 2 -> 4 blocks (2 MiB).
	sSlab := NewNodeStore()
	SetCurrentNodeStore(sSlab)
	_ = sSlab.Alloc()
	if len(sSlab.blocks) != 1 {
		t.Fatalf("blkIdx=0: len(blocks) = %d, want 1", len(sSlab.blocks))
	}
	for sSlab.NodeCount() <= blockLen {
		_ = sSlab.Alloc()
	}
	if len(sSlab.blocks) != 2 {
		t.Fatalf("blkIdx=1: len(blocks) = %d, want 2", len(sSlab.blocks))
	}
	for sSlab.NodeCount() <= 2*blockLen {
		_ = sSlab.Alloc()
	}
	if len(sSlab.blocks) != 6 {
		t.Fatalf("blkIdx=2: len(blocks) = %d, want 6", len(sSlab.blocks))
	}
	for sSlab.NodeCount() <= 6*blockLen {
		_ = sSlab.Alloc()
	}
	if len(sSlab.blocks) != 10 {
		t.Fatalf("blkIdx=6: len(blocks) = %d, want 10", len(sSlab.blocks))
	}
	var aSlab NodeAllocator
	aSlab.InitLocal()
	for sSlab.NodeCount() <= 10*blockLen {
		_ = aSlab.AllocNode1Cached(40)
	}
	aSlab.Flush()
	if len(sSlab.blocks) != 14 {
		t.Fatalf("blkIdx=10: len(blocks) = %d, want 14 (4-block 2 MiB slab)", len(sSlab.blocks))
	}
	for i := 10; i < 14; i++ {
		if sSlab.blocks[i] == nil {
			t.Fatalf("sSlab.blocks[%d] == nil inside 4-block slab", i)
		}
	}

	// 5. Concurrent -race stress test on SetLength (<0, 0..65534, >=65535) and HasNonNegativeLength.
	const workers = 12
	const perWorker = 300
	workerNodes := make([][]Node, workers)
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		w := w
		go func() {
			defer wg.Done()
			var la NodeAllocator
			la.InitLocal()
			ns := make([]Node, perWorker)
			for i := 0; i < perWorker; i++ {
				n := la.AllocNode1Cached(38)
				ns[i] = n
				switch (w + i) % 5 {
				case 0:
					n.SetLength(-int32(i + 1))
					if n.HasNonNegativeLength() {
						t.Errorf("expected false for -%d", i+1)
						return
					}
				case 1:
					n.SetLength(int32(i))
					if !n.HasNonNegativeLength() {
						t.Errorf("expected true for %d", i)
						return
					}
				case 2:
					n.SetLength(70_000 + int32(i))
					if !n.HasNonNegativeLength() {
						t.Errorf("expected true for overflow %d", 70_000+i)
						return
					}
				case 3:
					n.SetLength(-5)
					n.SetLength(int32(i))
					if !n.HasNonNegativeLength() {
						t.Errorf("expected true after -5 -> %d", i)
						return
					}
				case 4:
					n.SetLength(-9)
					n.SetLength(90_000 + int32(i))
					if !n.HasNonNegativeLength() {
						t.Errorf("expected true after -9 -> overflow")
						return
					}
				}
			}
			la.Flush()
			workerNodes[w] = ns
		}()
	}
	wg.Wait()

	// Clear remaining negative lengths (case 0) and verify negLenCount returns to baseNeg.
	for w := 0; w < workers; w++ {
		for i, n := range workerNodes[w] {
			if (w+i)%5 == 0 {
				if n.HasNonNegativeLength() {
					t.Fatalf("worker[%d][%d] expected negative length before clear", w, i)
				}
				n.SetLength(0)
			}
			if !n.HasNonNegativeLength() {
				t.Fatalf("worker[%d][%d] expected non-negative length after clear", w, i)
			}
		}
	}
	if got := negLenCount.Load(); got != baseNeg {
		t.Fatalf("after concurrent clear, negLenCount = %d, want %d", got, baseNeg)
	}
}
