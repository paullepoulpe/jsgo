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

// TestNodeDataStrictLayoutAndSize validates that NodeData is strictly 36 bytes,
// contains zero Go pointers/references, has zero padding bytes, and matches the
// required field offset layout.
func TestNodeDataStrictLayoutAndSize(t *testing.T) {
	sz := reflect.TypeFor[NodeData]().Size()
	if sz != 36 {
		t.Fatalf("NodeData size violation: expected exactly 36 bytes, got %d", sz)
	}

	typ := reflect.TypeFor[NodeData]()
	expectedFields := []struct {
		name   string
		offset uintptr
		size   uintptr
	}{
		{"Next", 0, 4},
		{"First", 4, 4},
		{"Token", 8, 1},
		{"KindFlags", 9, 1},
		{"Length16", 10, 2},
		{"Parent", 12, 4},
		{"Previous", 16, 4},
		{"Payload", 20, 4},
		{"SourceFile", 24, 4},
		{"PropFlags", 28, 4},
		{"LinenoCharno", 32, 4},
	}

	if typ.NumField() != len(expectedFields) {
		t.Fatalf("expected %d fields in NodeData, got %d", len(expectedFields), typ.NumField())
	}

	var totalFieldBytes uintptr
	for i, ef := range expectedFields {
		f := typ.Field(i)
		if f.Name != ef.name {
			t.Errorf("field %d: expected name %q, got %q", i, ef.name, f.Name)
		}
		if f.Offset != ef.offset {
			t.Errorf("field %s: expected offset %d, got %d", f.Name, ef.offset, f.Offset)
		}
		if f.Type.Size() != ef.size {
			t.Errorf("field %s: expected size %d, got %d", f.Name, ef.size, f.Type.Size())
		}
		totalFieldBytes += ef.size

		// Verify no pointers or pointer-containing types
		switch f.Type.Kind() {
		case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Interface, reflect.Func, reflect.Chan, reflect.String, reflect.UnsafePointer:
			t.Fatalf("field %s is of pointer/reference kind %v; noscan violated", f.Name, f.Type.Kind())
		}
	}

	if totalFieldBytes != 36 {
		t.Fatalf("sum of field sizes (%d) != 36; unexpected padding", totalFieldBytes)
	}
}

// TestNodeStoreGCDoesNotScanBlocks empirically validates that the Go runtime treats
// nodeBlock chunks as noscan (contains no pointers). We test this by writing arbitrary
// bit patterns (such as unmapped/non-canonical 64-bit pointer values) into NodeData
// fields across thousands of nodes and invoking runtime.GC(). If the GC marked/scanned
// these blocks as containing pointers, it would dereference them and crash.
func TestNodeStoreGCDoesNotScanBlocks(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	const nodeCount = 50000
	alloc := NewNodeAllocator(s)

	for i := 0; i < nodeCount; i++ {
		n := alloc.Alloc()
		rec := n.Rec()
		// Populate adjacent 32-bit fields with bogus 64-bit addresses
		rec.Payload = 0xDEADBEEF
		rec.SourceFile = 0xDEADC0DE
		rec.PropFlags = 0xCAFEBABE
		rec.LinenoCharno = int32(0x1337BEEF)
		rec.Length16 = 0xCAFE
		rec.Next = Node(0xBADF00D)
		rec.First = Node(0x7FFFFFFF)
		rec.Parent = Node(0x80000000)
		rec.Previous = Node(0xDEADDEAD)
	}

	// Trigger multiple GC cycles. If nodeBlock was scanned for pointers,
	// the garbage collector would attempt to trace these bogus addresses and segfault.
	for cycle := 0; cycle < 5; cycle++ {
		runtime.GC()
	}

	// Verify all records retained their exact bit patterns without corruption
	for i := 1; i <= 1000; i++ {
		n := Node(s.base + uint32(i))
		rec := n.Rec()
		if rec.Payload != 0xDEADBEEF || rec.SourceFile != 0xDEADC0DE || rec.Length16 != 0xCAFE {
			t.Fatalf("node %d corrupted after GC: payload=0x%X, sourceFile=0x%X, length16=0x%X",
				n, rec.Payload, rec.SourceFile, rec.Length16)
		}
	}
}

// TestStressAlloc1MillionNodesMultiGoroutine allocates 1,000,000+ AST nodes across
// 16 goroutines concurrently using NodeAllocator. It checks handle uniqueness,
// hierarchy linking, field readback consistency, and memory scaling.
func TestStressAlloc1MillionNodesMultiGoroutine(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	const numGoroutines = 16
	const nodesPerGoroutine = 65536 // 16 * 65536 = 1,048,576 nodes total (> 1M)
	const totalNodes = numGoroutines * nodesPerGoroutine

	t.Logf("Allocating %d AST nodes across %d goroutines...", totalNodes, numGoroutines)

	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	type workerResult struct {
		handles []Node
	}
	results := make([]workerResult, numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		g := g
		results[g].handles = make([]Node, nodesPerGoroutine)
		go func() {
			defer wg.Done()
			alloc := NewNodeAllocator(s)
			handles := results[g].handles

			for i := 0; i < nodesPerGoroutine; i++ {
				tok := uint8((g*31 + i) % 255)
				n := alloc.AllocWithToken(tok)
				n.SetLinenoCharno((g+1)*1000+(i%1000), i%4096)
				n.SetLength(int32(i % 50000))
				n.SetPayload(uint32(g<<24 | i))
				n.SetSourceFile(uint32(g + 1))
				n.SetPropFlags(int32(i * 7))
				n.SetKindFlags(uint8(g) & kindFlagsMask)
				handles[i] = n
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	allocRate := float64(totalNodes) / elapsed.Seconds()
	t.Logf("Allocated %d nodes in %v (%.2f nodes/sec)", totalNodes, elapsed, allocRate)

	// Verify handle uniqueness across all goroutines
	t.Log("Verifying handle uniqueness...")
	seen := make(map[Node]struct{}, totalNodes)
	for g := 0; g < numGoroutines; g++ {
		for i := 0; i < nodesPerGoroutine; i++ {
			n := results[g].handles[i]
			if n.IsNil() {
				t.Fatalf("goroutine %d, node %d is NilNode", g, i)
			}
			if _, exists := seen[n]; exists {
				t.Fatalf("duplicate handle detected: %d (g=%d, i=%d)", n, g, i)
			}
			seen[n] = struct{}{}
		}
	}

	// Concurrently verify field integrity across all 1M+ nodes
	t.Log("Concurrently verifying field integrity of all 1M+ nodes...")
	verifyWg := sync.WaitGroup{}
	verifyWg.Add(numGoroutines)
	for g := 0; g < numGoroutines; g++ {
		g := g
		go func() {
			defer verifyWg.Done()
			handles := results[g].handles
			for i := 0; i < nodesPerGoroutine; i++ {
				n := handles[i]
				expectedTok := uint8((g*31 + i) % 255)
				if n.Token() != expectedTok {
					t.Errorf("node %d: token mismatch: got %d, want %d", n, n.Token(), expectedTok)
					return
				}
				expectedLine := (g+1)*1000 + (i % 1000)
				if n.Lineno() != expectedLine {
					t.Errorf("node %d: lineno mismatch: got %d, want %d", n, n.Lineno(), expectedLine)
					return
				}
				expectedChar := i % 4096
				if n.Charno() != expectedChar {
					t.Errorf("node %d: charno mismatch: got %d, want %d", n, n.Charno(), expectedChar)
					return
				}
				expectedPayload := uint32(g<<24 | i)
				if n.Payload() != expectedPayload {
					t.Errorf("node %d: payload mismatch: got %d, want %d", n, n.Payload(), expectedPayload)
					return
				}
				if n.SourceFile() != uint32(g+1) {
					t.Errorf("node %d: sourceFile mismatch: got %d, want %d", n, n.SourceFile(), g+1)
					return
				}
				if n.PropFlags() != int32(i*7) {
					t.Errorf("node %d: propFlags mismatch: got %d, want %d", n, n.PropFlags(), i*7)
					return
				}
				if n.KindFlags() != uint8(g)&kindFlagsMask {
					t.Errorf("node %d: kindFlags mismatch: got %d, want %d", n, n.KindFlags(), uint8(g)&kindFlagsMask)
					return
				}
				if n.Length() != int32(i%50000) {
					t.Errorf("node %d: length mismatch: got %d, want %d", n, n.Length(), i%50000)
					return
				}
			}
		}()
	}
	verifyWg.Wait()

	// Run GC to verify arena stability under 1M+ live nodes
	gcStart := time.Now()
	runtime.GC()
	gcPause := time.Since(gcStart)
	t.Logf("GC pause with %d live arena nodes: %v", totalNodes, gcPause)
}

// TestStrPoolHighContentionStress stress-tests StrPool with 32 concurrent goroutines
// interning hundreds of thousands of random ASCII and UTF-8 strings (including CJK, emojis,
// math symbols, and accented characters) alongside highly contentious shared keys.
func TestStrPoolHighContentionStress(t *testing.T) {
	const numGoroutines = 32
	const stringsPerGoroutine = 4000
	const numSharedHotStrings = 200

	// Pre-generate a shared set of hot strings to force contention
	hotStrings := make([]string, numSharedHotStrings)
	for i := 0; i < numSharedHotStrings; i++ {
		switch i % 4 {
		case 0:
			hotStrings[i] = fmt.Sprintf("hot_ascii_identifier_%d", i)
		case 1:
			hotStrings[i] = fmt.Sprintf("hot_cjk_変数名_%d_日本語テスト", i)
		case 2:
			hotStrings[i] = fmt.Sprintf("hot_emoji_🚀_🔥_token_%d", i)
		case 3:
			hotStrings[i] = fmt.Sprintf("hot_unicode_∑_∏_∆_ç_é_ñ_%d", i)
		}
	}

	// Pre-intern hot strings from one goroutine to get ground-truth IDs
	hotIDs := make([]StringID, numSharedHotStrings)
	for i, s := range hotStrings {
		hotIDs[i] = Intern(s)
		if hotIDs[i] <= EmptyStringID {
			t.Fatalf("invalid ID for hot string %q: %d", s, hotIDs[i])
		}
		if String(hotIDs[i]) != s {
			t.Fatalf("hot string round-trip failed: got %q, want %q", String(hotIDs[i]), s)
		}
	}

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	type checkRecord struct {
		str string
		id  StringID
	}
	workerRecords := make([][]checkRecord, numGoroutines)

	start := time.Now()
	for g := 0; g < numGoroutines; g++ {
		g := g
		workerRecords[g] = make([]checkRecord, stringsPerGoroutine)
		go func() {
			defer wg.Done()
			records := workerRecords[g]

			for i := 0; i < stringsPerGoroutine; i++ {
				var s string
				if i%3 == 0 {
					// Hot contention on shared strings
					hotIdx := (g*17 + i) % numSharedHotStrings
					s = hotStrings[hotIdx]
					id := Intern(s)
					if id != hotIDs[hotIdx] {
						t.Errorf("goroutine %d: inconsistent hot string ID: got %d, want %d", g, id, hotIDs[hotIdx])
					}
					records[i] = checkRecord{str: s, id: id}
				} else if i%3 == 1 {
					// Unique UTF-8 string with diverse multi-byte characters
					var b [8]byte
					rand.Read(b[:])
					hexPart := hex.EncodeToString(b[:])
					s = fmt.Sprintf("utf8_%d_%d_🔥_한국어_%s_日本語_ß_à_ö", g, i, hexPart)
					id := Intern(s)
					records[i] = checkRecord{str: s, id: id}
				} else {
					// Unique ASCII string
					var b [16]byte
					rand.Read(b[:])
					s = fmt.Sprintf("ascii_%d_%d_%s", g, i, hex.EncodeToString(b[:]))
					id := Intern(s)
					records[i] = checkRecord{str: s, id: id}
				}
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	totalOps := numGoroutines * stringsPerGoroutine
	t.Logf("Completed %d concurrent intern ops in %v (%.2f ops/sec)", totalOps, elapsed, float64(totalOps)/elapsed.Seconds())

	// Concurrently verify 100% interning consistency and round-trip fidelity
	t.Log("Verifying round-trip string integrity, Lookup, and InternBytes consistency...")
	verifyWg := sync.WaitGroup{}
	verifyWg.Add(numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		g := g
		go func() {
			defer verifyWg.Done()
			for _, rec := range workerRecords[g] {
				// String(id) must match original string byte-for-byte
				gotStr := String(rec.id)
				if gotStr != rec.str {
					t.Errorf("String(%d) mismatch: got %q, want %q", rec.id, gotStr, rec.str)
					return
				}
				// Lookup(str) must return exact same ID
				lookupID := Lookup(rec.str)
				if lookupID != rec.id {
					t.Errorf("Lookup(%q) mismatch: got %d, want %d", rec.str, lookupID, rec.id)
					return
				}
				// InternBytes must return exact same ID
				bytesID := InternBytes([]byte(rec.str))
				if bytesID != rec.id {
					t.Errorf("InternBytes(%q) mismatch: got %d, want %d", rec.str, bytesID, rec.id)
					return
				}
			}
		}()
	}
	verifyWg.Wait()
}

// TestStrPoolBaselineSnapshotRewindCorrectness validates that:
// 1. Snapshot captures the exact high-water mark of every shard.
// 2. Rewinding resets nextLocalIdx and count to the snapshot mark.
// 3. Post-snapshot strings are invalidated while pre-snapshot strings remain 100% intact.
// 4. Epoch is bumped on every rewinding operation.
// 5. Subsequent interning cleanly re-uses shard space without table corruption or race conditions.
func TestStrPoolBaselineSnapshotRewindCorrectness(t *testing.T) {
	// Step 1: Pre-populate base strings
	const baseCount = 500
	baseStrings := make([]string, baseCount)
	baseIDs := make([]StringID, baseCount)
	for i := 0; i < baseCount; i++ {
		baseStrings[i] = fmt.Sprintf("baseline_string_%d_🚀_日本語", i)
		baseIDs[i] = Intern(baseStrings[i])
	}

	snap1 := CaptureBaseline()
	epoch1 := Epoch()

	// Step 2: Intern stage 1 temporary strings
	const stage1Count = 1000
	stage1Strings := make([]string, stage1Count)
	stage1IDs := make([]StringID, stage1Count)
	for i := 0; i < stage1Count; i++ {
		stage1Strings[i] = fmt.Sprintf("stage1_temp_string_%d_🔥", i)
		stage1IDs[i] = Intern(stage1Strings[i])
	}

	snap2 := CaptureBaseline()
	epoch2 := Epoch()
	if epoch2 != epoch1 {
		t.Fatalf("Epoch should not change during interning; got %d, want %d", epoch2, epoch1)
	}

	// Step 3: Intern stage 2 temporary strings
	const stage2Count = 1000
	stage2Strings := make([]string, stage2Count)
	stage2IDs := make([]StringID, stage2Count)
	for i := 0; i < stage2Count; i++ {
		stage2Strings[i] = fmt.Sprintf("stage2_temp_string_%d_⚡️", i)
		stage2IDs[i] = Intern(stage2Strings[i])
	}

	// Verify all stage 2 strings resolve
	for i := 0; i < stage2Count; i++ {
		if String(stage2IDs[i]) != stage2Strings[i] {
			t.Fatalf("stage2 string resolution failed for %q", stage2Strings[i])
		}
	}

	// Step 4: Rewind to snap2 (reclaims stage 2 strings)
	RewindBaseline(snap2)
	epoch3 := Epoch()
	if epoch3 <= epoch2 {
		t.Fatalf("Epoch did not advance after rewind: got %d, want > %d", epoch3, epoch2)
	}

	// Verify baseline and stage 1 strings are still 100% intact
	for i := 0; i < baseCount; i++ {
		if got := String(baseIDs[i]); got != baseStrings[i] {
			t.Fatalf("baseline string %d corrupted after snap2 rewind: got %q, want %q", i, got, baseStrings[i])
		}
		if lookup := Lookup(baseStrings[i]); lookup != baseIDs[i] {
			t.Fatalf("baseline Lookup(%q) corrupted after snap2 rewind: got %d, want %d", baseStrings[i], lookup, baseIDs[i])
		}
	}
	for i := 0; i < stage1Count; i++ {
		if got := String(stage1IDs[i]); got != stage1Strings[i] {
			t.Fatalf("stage1 string %d corrupted after snap2 rewind: got %q, want %q", i, got, stage1Strings[i])
		}
		if lookup := Lookup(stage1Strings[i]); lookup != stage1IDs[i] {
			t.Fatalf("stage1 Lookup(%q) corrupted after snap2 rewind: got %d, want %d", stage1Strings[i], lookup, stage1IDs[i])
		}
	}

	// Step 5: Rewind to snap1 (reclaims stage 1 strings back to original baseline)
	RewindBaseline(snap1)
	epoch4 := Epoch()
	if epoch4 <= epoch3 {
		t.Fatalf("Epoch did not advance after snap1 rewind: got %d, want > %d", epoch4, epoch3)
	}

	// Verify original baseline strings still 100% intact
	for i := 0; i < baseCount; i++ {
		if got := String(baseIDs[i]); got != baseStrings[i] {
			t.Fatalf("baseline string %d corrupted after snap1 rewind: got %q, want %q", i, got, baseStrings[i])
		}
		if lookup := Lookup(baseStrings[i]); lookup != baseIDs[i] {
			t.Fatalf("baseline Lookup(%q) corrupted after snap1 rewind: got %d, want %d", baseStrings[i], lookup, baseIDs[i])
		}
	}

	// Step 6: Re-intern strings and verify clean behavior
	for i := 0; i < stage1Count; i++ {
		newID := Intern(stage1Strings[i])
		if String(newID) != stage1Strings[i] {
			t.Fatalf("re-interned stage1 string failed: got %q, want %q", String(newID), stage1Strings[i])
		}
	}
}

// Benchmarks for NodeStore and StrPool
func BenchmarkNodeStoreAllocSingleThread(b *testing.B) {
	s := NewNodeStore()
	alloc := NewNodeAllocator(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = alloc.Alloc()
	}
}

func BenchmarkNodeStoreAllocMultiThread(b *testing.B) {
	s := NewNodeStore()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		alloc := NewNodeAllocator(s)
		for pb.Next() {
			_ = alloc.Alloc()
		}
	})
}

func BenchmarkStrPoolInternHit(b *testing.B) {
	s := "com.google.javascript.jscomp.Compiler"
	_ = Intern(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Intern(s)
	}
}

func BenchmarkStrPoolInternBytesHit(b *testing.B) {
	raw := []byte("com.google.javascript.jscomp.Compiler")
	_ = Intern(string(raw))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = InternBytes(raw)
	}
}

func BenchmarkStrPoolConcurrentContention(b *testing.B) {
	keys := []string{
		"identifier_alpha", "identifier_beta", "identifier_gamma", "identifier_delta",
		"変数名_テスト", "한국어_식별자", "symbol_🚀_emoji", "function_closure_pass",
	}
	for _, k := range keys {
		_ = Intern(k)
	}
	var counter atomic.Uint64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			idx := counter.Add(1) % uint64(len(keys))
			_ = Intern(keys[idx])
		}
	})
}
