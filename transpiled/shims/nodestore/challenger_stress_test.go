package nodestore

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestChallengerNodeDataNoscanAndPtrdata empirically validates that NodeData and
// nodeBlock are 36-byte pointer-free structs, confirming true noscan GC semantics.
func TestChallengerNodeDataNoscanAndPtrdata(t *testing.T) {
	if hasGoPointers(reflect.TypeFor[NodeData]()) {
		t.Fatalf("noscan violation: NodeData contains Go pointers")
	}

	if hasGoPointers(reflect.TypeFor[nodeBlock]()) {
		t.Fatalf("noscan violation: nodeBlock contains Go pointers")
	}

	if sz := reflect.TypeFor[NodeData]().Size(); sz != 36 {
		t.Fatalf("NodeData size violation: expected 36, got %d", sz)
	}

	if sz := reflect.TypeFor[nodeBlock]().Size(); sz != 16384*36 {
		t.Fatalf("nodeBlock size violation: expected %d, got %d", 16384*36, sz)
	}

	// Verify that filling NodeData with invalid addresses and forcing GC does not crash
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	const count = 20000
	alloc := NewNodeAllocator(s)
	for i := 0; i < count; i++ {
		n := alloc.Alloc()
		rec := n.Rec()
		rec.Payload = 0xDEADBEEF
		rec.SourceFile = 0xBADC0FFE
		rec.PropFlags = 0xCAFEBABE
		rec.LinenoCharno = int32(0x1337BEEF)
		rec.Length16 = 0xCAFE
		rec.Next = Node(0xBADF00D)
		rec.First = Node(0x7FFFFFFF)
		rec.Parent = Node(0x80000000)
		rec.Previous = Node(0xDEADDEAD)
	}

	for i := 0; i < 5; i++ {
		runtime.GC()
	}
}

// TestChallengerPointerStability validates non-moving indexing:
// Pointers to NodeData elements in memory must remain at stable virtual addresses
// even while 1,000,000+ additional nodes are allocated.
func TestChallengerPointerStability(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	const initialCount = 5000
	type pinnedNode struct {
		node    Node
		ptrAddr string
		rawPtr  *NodeData
		token   uint8
		payload uint32
	}

	pinned := make([]pinnedNode, initialCount)
	for i := 0; i < initialCount; i++ {
		n := s.Alloc()
		tok := uint8(i%250 + 1)
		payload := uint32(0xA0000000 | i)
		n.SetToken(tok)
		n.SetPayload(payload)
		rec := n.Rec()
		pinned[i] = pinnedNode{
			node:    n,
			ptrAddr: fmt.Sprintf("%p", rec),
			rawPtr:  rec,
			token:   tok,
			payload: payload,
		}
	}

	// Now allocate 1,000,000+ more nodes across 16 goroutines
	const goroutines = 16
	const perGoroutine = 65536 // 16 * 65536 = 1,048,576 nodes
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			alloc := NewNodeAllocator(s)
			for i := 0; i < perGoroutine; i++ {
				n := alloc.Alloc()
				n.SetToken(uint8(i % 255))
			}
		}()
	}
	wg.Wait()

	// Verify that EVERY initial node has exact same pointer address and data
	for i := 0; i < initialCount; i++ {
		pn := pinned[i]
		rec := pn.node.Rec()
		currentAddr := fmt.Sprintf("%p", rec)
		if currentAddr != pn.ptrAddr {
			t.Fatalf("node %d moved in memory! original=%s, current=%s", pn.node, pn.ptrAddr, currentAddr)
		}
		if pn.rawPtr != rec {
			t.Fatalf("node %d raw pointer mismatch", pn.node)
		}
		if rec.Token != pn.token || rec.Payload != pn.payload {
			t.Fatalf("node %d data corrupted! token=%d (want %d), payload=%x (want %x)",
				pn.node, rec.Token, pn.token, rec.Payload, pn.payload)
		}
	}
}

// TestChallengerDirectConcurrentAlloc tests high contention on the raw s.Alloc()
// method (without thread-local NodeAllocator buffering) to challenge block boundary
// allocation locks and atomic s.next synchronization.
func TestChallengerDirectConcurrentAlloc(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	const goroutines = 16
	const perGoroutine = 20000 // 320,000 nodes directly allocated
	var wg sync.WaitGroup
	wg.Add(goroutines)

	allHandles := make([][]Node, goroutines)
	for g := 0; g < goroutines; g++ {
		g := g
		allHandles[g] = make([]Node, perGoroutine)
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				n := s.Alloc()
				n.SetToken(uint8((g*17 + i) % 255))
				n.SetPayload(uint32(g<<20 | i))
				allHandles[g][i] = n
			}
		}()
	}
	wg.Wait()

	seen := make(map[Node]struct{}, goroutines*perGoroutine)
	for g := 0; g < goroutines; g++ {
		for i := 0; i < perGoroutine; i++ {
			n := allHandles[g][i]
			if n.IsNil() {
				t.Fatalf("goroutine %d, index %d allocated NilNode", g, i)
			}
			if _, exists := seen[n]; exists {
				t.Fatalf("duplicate node handle detected: %d", n)
			}
			seen[n] = struct{}{}

			expectedTok := uint8((g*17 + i) % 255)
			if n.Token() != expectedTok {
				t.Fatalf("token mismatch for node %d: got %d, want %d", n, n.Token(), expectedTok)
			}
			expectedPayload := uint32(g<<20 | i)
			if n.Payload() != expectedPayload {
				t.Fatalf("payload mismatch for node %d: got %x, want %x", n, n.Payload(), expectedPayload)
			}
		}
	}
}

// TestChallengerStrPoolZeroAllocCacheHits empirically verifies that cache hits
// on Intern, InternBytes, Lookup, and String incur ZERO heap allocations.
func TestChallengerStrPoolZeroAllocCacheHits(t *testing.T) {
	testStrings := []string{
		"ascii_simple_identifier",
		"utf8_변수명_関数_ident",
		"emoji_🚀_🔥_token_⚡️",
		"unicode_math_∑_∏_∆_∇_∫",
		"accented_é_à_ç_ñ_ü_ß_ø",
		"very_long_string_payload_" + hex.EncodeToString(make([]byte, 2048)),
		"",
	}

	for _, str := range testStrings {
		id := Intern(str)
		rawBytes := []byte(str)

		// Test Intern(string) hit
		allocsIntern := testing.AllocsPerRun(200, func() {
			gotID := Intern(str)
			if gotID != id {
				t.Fatalf("ID mismatch: got %d, want %d", gotID, id)
			}
		})
		if allocsIntern > 0 {
			t.Errorf("Intern(%q) allocated %f objects on hit; expected 0", str, allocsIntern)
		}

		// Test InternBytes([]byte) hit
		allocsBytes := testing.AllocsPerRun(200, func() {
			gotID := InternBytes(rawBytes)
			if gotID != id {
				t.Fatalf("ID mismatch: got %d, want %d", gotID, id)
			}
		})
		if allocsBytes > 0 {
			t.Errorf("InternBytes(%q) allocated %f objects on hit; expected 0", str, allocsBytes)
		}

		// Test Lookup(string) hit
		allocsLookup := testing.AllocsPerRun(200, func() {
			gotID := Lookup(str)
			if gotID != id {
				t.Fatalf("ID mismatch: got %d, want %d", gotID, id)
			}
		})
		if allocsLookup > 0 {
			t.Errorf("Lookup(%q) allocated %f objects on hit; expected 0", str, allocsLookup)
		}

		// Test String(id)
		allocsString := testing.AllocsPerRun(200, func() {
			gotStr := String(id)
			if gotStr != str {
				t.Fatalf("String mismatch: got %q, want %q", gotStr, str)
			}
		})
		if allocsString > 0 {
			t.Errorf("String(%d) allocated %f objects; expected 0", id, allocsString)
		}
	}
}

// TestChallengerStrPoolExtremeContention challenges StrPool with 64 concurrent goroutines
// interleaving string interning, byte interning, and lookups with random multi-byte strings.
func TestChallengerStrPoolExtremeContention(t *testing.T) {
	const goroutines = 64
	const opsPerGoroutine = 3000
	const sharedKeyCount = 100

	sharedKeys := make([]string, sharedKeyCount)
	for i := 0; i < sharedKeyCount; i++ {
		sharedKeys[i] = fmt.Sprintf("shared_key_%d_🚀_日本語_한국어", i)
		_ = Intern(sharedKeys[i])
	}

	var wg sync.WaitGroup
	wg.Add(goroutines)
	var totalSuccess atomic.Uint64

	start := time.Now()
	for g := 0; g < goroutines; g++ {
		g := g
		go func() {
			defer wg.Done()
			var localSuccess uint64
			for i := 0; i < opsPerGoroutine; i++ {
				var key string
				switch i % 5 {
				case 0, 1:
					// Shared hot key contention
					key = sharedKeys[(g*13+i)%sharedKeyCount]
					id := Intern(key)
					if String(id) != key {
						t.Errorf("String(%d) != %q", id, key)
					}
					localSuccess++
				case 2:
					// Lookup on hot key
					key = sharedKeys[(g*17+i)%sharedKeyCount]
					id := Lookup(key)
					if id <= EmptyStringID {
						t.Errorf("Lookup on %q returned %d", key, id)
					}
					localSuccess++
				case 3:
					// InternBytes on hot key
					key = sharedKeys[(g*19+i)%sharedKeyCount]
					b := []byte(key)
					id := InternBytes(b)
					if String(id) != key {
						t.Errorf("InternBytes String(%d) != %q", id, key)
					}
					localSuccess++
				case 4:
					// Unique random multi-byte string
					var rnd [8]byte
					rand.Read(rnd[:])
					key = fmt.Sprintf("unique_%d_%d_%s_⚡️_日本語", g, i, hex.EncodeToString(rnd[:]))
					id := Intern(key)
					if String(id) != key {
						t.Errorf("Unique string roundtrip mismatch for %q", key)
					}
					localSuccess++
				}
			}
			totalSuccess.Add(localSuccess)
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	expectedTotal := uint64(goroutines * opsPerGoroutine)
	if totalSuccess.Load() != expectedTotal {
		t.Fatalf("expected %d successful ops, got %d", expectedTotal, totalSuccess.Load())
	}
	t.Logf("Completed %d extreme contention StrPool operations in %v (%.2f ops/sec)",
		expectedTotal, elapsed, float64(expectedTotal)/elapsed.Seconds())
}

// TestChallengerTreeHierarchyStress constructs an AST with 50,000 subtrees,
// verifying child and sibling bidirectional navigation and parent pointers.
func TestChallengerTreeHierarchyStress(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	alloc := NewNodeAllocator(s)
	const numSubtrees = 50000

	root := alloc.AllocWithToken(1) // Token 1 = ROOT
	var prevSubtree Node

	for i := 0; i < numSubtrees; i++ {
		c1 := alloc.AllocWithToken(10)
		c2 := alloc.AllocWithToken(20)
		c3 := alloc.AllocWithToken(30)
		sub := alloc.Alloc3ChildrenCachedTemplate(100, c1, c2, c3)

		if root.FirstChild() == 0 {
			root.SetFirstChild(sub)
			sub.SetParent(root)
			sub.SetPrevious(sub)
		} else {
			prevSubtree.SetNext(sub)
			sub.SetParent(root)
			sub.SetPrevious(prevSubtree)
			root.FirstChild().SetPrevious(sub)
		}
		prevSubtree = sub
	}

	// Traverse children of root
	count := 0
	curr := root.FirstChild()
	for curr != 0 {
		count++
		if curr.Parent() != root {
			t.Fatalf("child %d has wrong parent", count)
		}
		if curr.Token() != 100 {
			t.Fatalf("child %d has wrong token", count)
		}
		// Check c1, c2, c3
		c1 := curr.FirstChild()
		if c1 == 0 || c1.Token() != 10 {
			t.Fatalf("c1 invalid in subtree %d", count)
		}
		c2 := c1.Next()
		if c2 == 0 || c2.Token() != 20 {
			t.Fatalf("c2 invalid in subtree %d", count)
		}
		c3 := c2.Next()
		if c3 == 0 || c3.Token() != 30 {
			t.Fatalf("c3 invalid in subtree %d", count)
		}
		curr = curr.Next()
	}

	if count != numSubtrees {
		t.Fatalf("expected %d traversed subtrees, got %d", numSubtrees, count)
	}
}
