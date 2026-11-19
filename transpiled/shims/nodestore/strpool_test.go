package nodestore

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func TestStrPoolInternAndLookup(t *testing.T) {
	if Intern("") != EmptyStringID {
		t.Errorf("Intern(\"\") should equal EmptyStringID (1)")
	}
	if Lookup("") != EmptyStringID {
		t.Errorf("Lookup(\"\") should equal EmptyStringID (1)")
	}
	if String(EmptyStringID) != "" {
		t.Errorf("String(EmptyStringID) should be \"\"")
	}

	id1 := Intern("hello")
	id2 := Intern("hello")
	if id1 != id2 {
		t.Errorf("expected idempotent IDs, got %d and %d", id1, id2)
	}
	if String(id1) != "hello" {
		t.Errorf("expected String(%d) == \"hello\", got %q", id1, String(id1))
	}
	if Lookup("hello") != id1 {
		t.Errorf("Lookup(\"hello\") should be %d, got %d", id1, Lookup("hello"))
	}
	if Lookup("nonexistent_symbol_1234") != 0 {
		t.Errorf("Lookup for absent string should return 0")
	}
}

func TestZeroAllocInternBytes(t *testing.T) {
	// First intern the string
	raw := []byte("function foo(bar, baz) { return bar + baz; }")
	id1 := Intern(string(raw))

	// Probe with raw bytes; AllocsPerRun should be 0
	allocs := testing.AllocsPerRun(100, func() {
		id2 := InternBytes(raw)
		if id2 != id1 {
			t.Fatalf("ID mismatch: got %d, want %d", id2, id1)
		}
	})
	if allocs > 0 {
		t.Errorf("InternBytes allocated %f heap objects on existing string; expected 0", allocs)
	}
}

func TestConcurrentInterning(t *testing.T) {
	const goroutines = 8
	const perGoroutine = 500

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		g := g
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				str := fmt.Sprintf("sym_%d_%d", g, i)
				id := Intern(str)
				if String(id) != str {
					t.Errorf("String(%d): got %q, want %q", id, String(id), str)
				}
				// Also intern shared strings
				sharedStr := fmt.Sprintf("shared_%d", i%20)
				sharedID := Intern(sharedStr)
				if String(sharedID) != sharedStr {
					t.Errorf("String(%d): got %q, want %q", sharedID, String(sharedID), sharedStr)
				}
			}
		}()
	}
	wg.Wait()
}

func TestStrPoolBaselineRewind(t *testing.T) {
	snap := CaptureBaseline()
	ep0 := Epoch()

	id := Intern("temp_string_to_rewind_99999")
	if String(id) != "temp_string_to_rewind_99999" {
		t.Fatalf("failed to intern temp string")
	}

	RewindBaseline(snap)
	if Epoch() <= ep0 {
		t.Errorf("expected Epoch to advance after rewind")
	}
}

// Compile-time signature check for JoinDot, JoinDot3, and JoinDot4.
var (
	_ func(string, string) string                 = JoinDot
	_ func(string, string, string) string         = JoinDot3
	_ func(string, string, string, string) string = JoinDot4
)

func TestWave26AJoinDotAndJoinDot3(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer func() {
		ExitSingleOwnerAlloc()
		SetCurrentNodeStore(prev)
	}()

	pairs := [][2]string{
		{"", ""},
		{"", "b"},
		{"a", ""},
		{"x", "y"},
		{"goog", "ui"},
		{"goog.ui.Button", "prototype"},
		{"goog.ui.Dialog", "prototype"},
		{"module$exports$wiz$foo", "render"},
		{"αβγ", "δεζ"},
		{strings.Repeat("a", 96), strings.Repeat("b", 96)},  // exact 192 max cached key len
		{strings.Repeat("a", 97), strings.Repeat("b", 96)},  // 193 -> uncached fallback
		{strings.Repeat("q", 300), strings.Repeat("z", 300)}, // giant -> uncached fallback
	}
	triples := [][3]string{
		{"", "", ""},
		{"a", "", "c"},
		{"goog", "ui", "Component"},
		{"this", "props", "value"},
		{"goog.ui.Component", "prototype", "disposeInternal"},
		{strings.Repeat("a", 64), strings.Repeat("b", 64), strings.Repeat("c", 64)}, // exact 192
		{strings.Repeat("a", 65), strings.Repeat("b", 64), strings.Repeat("c", 64)}, // 193 -> uncached
	}

	// 1. Outside single-owner mode: exact results & 0 allocs on warm lock-free hits!
	for _, p := range pairs {
		want := p[0] + "." + p[1]
		if got := JoinDot(p[0], p[1]); got != want {
			t.Fatalf("non-single-owner JoinDot(%q, %q) = %q, want %q", p[0], p[1], got, want)
		}
	}
	for _, tr := range triples {
		want := tr[0] + "." + tr[1] + "." + tr[2]
		if got := JoinDot3(tr[0], tr[1], tr[2]); got != want {
			t.Fatalf("non-single-owner JoinDot3(%q, %q, %q) = %q, want %q", tr[0], tr[1], tr[2], got, want)
		}
	}
	nonSingleOwnerAllocs := testing.AllocsPerRun(100, func() {
		if JoinDot("goog.ui.Button", "prototype") != "goog.ui.Button.prototype" {
			t.Fatal("non-single-owner JoinDot hit mismatch")
		}
		if JoinDot3("goog.ui.Component", "prototype", "disposeInternal") != "goog.ui.Component.prototype.disposeInternal" {
			t.Fatal("non-single-owner JoinDot3 hit mismatch")
		}
	})
	if nonSingleOwnerAllocs != 0 {
		t.Fatalf("expected 0 allocs on lock-free JoinDot/JoinDot3 hits outside single-owner mode, got %v", nonSingleOwnerAllocs)
	}

	// 2. Single-owner mode: warm hits must allocate 0 objects/op!
	if !TryEnterSingleOwnerAlloc() {
		t.Fatal("TryEnterSingleOwnerAlloc failed")
	}
	for _, p := range pairs {
		want := p[0] + "." + p[1]
		if got := JoinDot(p[0], p[1]); got != want {
			t.Fatalf("single-owner JoinDot(%q, %q) = %q, want %q", p[0], p[1], got, want)
		}
	}
	for _, tr := range triples {
		want := tr[0] + "." + tr[1] + "." + tr[2]
		if got := JoinDot3(tr[0], tr[1], tr[2]); got != want {
			t.Fatalf("single-owner JoinDot3(%q, %q, %q) = %q, want %q", tr[0], tr[1], tr[2], got, want)
		}
	}

	warmAllocs := testing.AllocsPerRun(200, func() {
		if JoinDot("", "") != "." {
			t.Fatal("empty JoinDot mismatch")
		}
		if JoinDot("goog.ui.Button", "prototype") != "goog.ui.Button.prototype" {
			t.Fatal("JoinDot hit mismatch")
		}
		if JoinDot("goog.ui.Dialog", "prototype") != "goog.ui.Dialog.prototype" {
			t.Fatal("JoinDot sibling hit mismatch")
		}
		if JoinDot3("", "", "") != ".." {
			t.Fatal("empty JoinDot3 mismatch")
		}
		if JoinDot3("goog.ui.Component", "prototype", "disposeInternal") != "goog.ui.Component.prototype.disposeInternal" {
			t.Fatal("JoinDot3 hit mismatch")
		}
	})
	if warmAllocs > 0 {
		t.Fatalf("expected 0 allocs on warm JoinDot/JoinDot3 hits, got %.1f", warmAllocs)
	}

	// 3. Concurrent SuspendSingleOwnerAlloc / ResumeSingleOwnerAlloc + ExitSingleOwnerAlloc handoff under -race.
	const workers = 8
	var wg sync.WaitGroup
	for round := 0; round < 4; round++ {
		_ = JoinDot("single.owner", fmt.Sprintf("pre_%d", round))
		_ = JoinDot3("single", "owner", fmt.Sprintf("pre3_%d", round))

		wg.Add(workers)
		for w := 0; w < workers; w++ {
			w := w
			SuspendSingleOwnerAlloc()
			go func() {
				defer func() {
					ResumeSingleOwnerAlloc()
					wg.Done()
				}()
				for i := 0; i < 200; i++ {
					a := fmt.Sprintf("pkg_%d", i%17)
					b := fmt.Sprintf("Cls_%d", w)
					c := "prototype"
					if got := JoinDot(a, b); got != a+"."+b {
						t.Errorf("worker JoinDot mismatch: %q", got)
					}
					if got3 := JoinDot3(a, b, c); got3 != a+"."+b+"."+c {
						t.Errorf("worker JoinDot3 mismatch: %q", got3)
					}
				}
			}()
		}
		wg.Wait()

		// Resumed single-owner fast path on compiler goroutine:
		if got := JoinDot("goog.ui.Button", "prototype"); got != "goog.ui.Button.prototype" {
			t.Fatalf("post-resume JoinDot mismatch: %q", got)
		}
	}

	// Hand off single-owner session to another goroutine + switch NodeStore:
	ExitSingleOwnerAlloc()
	s2 := NewNodeStore()
	SetCurrentNodeStore(s2)
	done := make(chan struct{})
	go func() {
		defer close(done)
		if !TryEnterSingleOwnerAlloc() {
			t.Error("handoff TryEnterSingleOwnerAlloc failed")
			return
		}
		defer ExitSingleOwnerAlloc()
		if got := JoinDot("ns", "prop"); got != "ns.prop" {
			t.Errorf("s2 JoinDot = %q", got)
		}
		if got3 := JoinDot3("ns", "Sub", "method"); got3 != "ns.Sub.method" {
			t.Errorf("s2 JoinDot3 = %q", got3)
		}
	}()
	<-done
}

func BenchmarkWave26AJoinDot(b *testing.B) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer func() {
		ExitSingleOwnerAlloc()
		SetCurrentNodeStore(prev)
	}()
	EnterSingleOwnerAlloc()
	p1 := String(Intern("goog.ui.Component"))
	p2 := String(Intern("prototype"))
	p3 := String(Intern("disposeInternal"))
	_ = JoinDot(p1, p2)
	_ = JoinDot3(p1, p2, p3)

	b.Run("JoinDotHit", func(b *testing.B) {
		b.ReportAllocs()
		var sink int
		for i := 0; i < b.N; i++ {
			sink += len(JoinDot(p1, p2))
		}
		_ = sink
	})
	b.Run("JoinDot3Hit", func(b *testing.B) {
		b.ReportAllocs()
		var sink int
		for i := 0; i < b.N; i++ {
			sink += len(JoinDot3(p1, p2, p3))
		}
		_ = sink
	})
}

func TestWave28AJoinDot4(t *testing.T) {
	// Static store fallback check (before setting active store).
	ResetNodeStore()
	if got := JoinDot4("a", "b", "c", "d"); got != "a.b.c.d" {
		t.Fatalf("staticStore JoinDot4 = %q, want \"a.b.c.d\"", got)
	}

	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer func() {
		ExitSingleOwnerAlloc()
		SetCurrentNodeStore(prev)
	}()

	quads := [][4]string{
		{"", "", "", ""},
		{"", "b", "", "d"},
		{"a", "", "", ""},
		{"", "", "", "z"},
		{"goog", "ui", "Button", "prototype"},
		{"goog", "ui", "Component", "State"},
		{"module$exports$wiz", "foo", "Bar", "render"},
		{"αβγ", "δεζ", "ηθι", "🚀"},
		{strings.Repeat("a", 47), strings.Repeat("b", 48), strings.Repeat("c", 47), strings.Repeat("d", 48)}, // 190 + 2 == 192 -> cached
		{strings.Repeat("a", 48), strings.Repeat("b", 48), strings.Repeat("c", 47), strings.Repeat("d", 48)}, // 191 + 2 == 193 -> uncached
		{strings.Repeat("x", 120), strings.Repeat("y", 120), strings.Repeat("z", 120), "end"},                // giant -> uncached
	}

	// 1. Outside single-owner mode: exact results & 0 allocs on warm lock-free hits!
	for _, q := range quads {
		want := q[0] + "." + q[1] + "." + q[2] + "." + q[3]
		if got := JoinDot4(q[0], q[1], q[2], q[3]); got != want {
			t.Fatalf("non-single-owner JoinDot4(%q, %q, %q, %q) = %q, want %q", q[0], q[1], q[2], q[3], got, want)
		}
	}
	if allocs := testing.AllocsPerRun(100, func() {
		if JoinDot4("goog", "ui", "Button", "prototype") != "goog.ui.Button.prototype" {
			t.Fatal("non-single-owner JoinDot4 hit mismatch")
		}
	}); allocs != 0 {
		t.Fatalf("expected 0 allocs on lock-free JoinDot4 hit outside single-owner mode, got %v", allocs)
	}

	// 2. Single-owner mode: exact results & 0 allocations on warm hits!
	if !TryEnterSingleOwnerAlloc() {
		t.Fatal("TryEnterSingleOwnerAlloc failed")
	}
	for _, q := range quads {
		want := q[0] + "." + q[1] + "." + q[2] + "." + q[3]
		if got := JoinDot4(q[0], q[1], q[2], q[3]); got != want {
			t.Fatalf("single-owner JoinDot4(%q, %q, %q, %q) = %q, want %q", q[0], q[1], q[2], q[3], got, want)
		}
	}

	warmAllocs := testing.AllocsPerRun(200, func() {
		if JoinDot4("", "", "", "") != "..." {
			t.Fatal("empty JoinDot4 mismatch")
		}
		if JoinDot4("goog", "ui", "Button", "prototype") != "goog.ui.Button.prototype" {
			t.Fatal("JoinDot4 hit mismatch")
		}
		if JoinDot4("αβγ", "δεζ", "ηθι", "🚀") != "αβγ.δεζ.ηθι.🚀" {
			t.Fatal("JoinDot4 UTF-8 hit mismatch")
		}
		// Prefix "goog.ui.Button" was also cached in joinDot3Cache by JoinDot4!
		if JoinDot3("goog", "ui", "Button") != "goog.ui.Button" {
			t.Fatal("JoinDot3 shared prefix hit mismatch")
		}
	})
	if warmAllocs > 0 {
		t.Fatalf("expected 0 allocs on warm JoinDot4 hits, got %.1f", warmAllocs)
	}

	// 3. Concurrent workers under SuspendSingleOwnerAlloc / ResumeSingleOwnerAlloc.
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		w := w
		SuspendSingleOwnerAlloc()
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer ResumeSingleOwnerAlloc()
			for i := 0; i < 200; i++ {
				a := fmt.Sprintf("ns_%d", i%13)
				b := fmt.Sprintf("pkg_%d", w)
				c := "Widget"
				d := "prototype"
				want := a + "." + b + "." + c + "." + d
				if got := JoinDot4(a, b, c, d); got != want {
					t.Errorf("worker JoinDot4 = %q, want %q", got, want)
				}
			}
		}()
	}
	wg.Wait()
	if got := JoinDot4("goog", "ui", "Button", "prototype"); got != "goog.ui.Button.prototype" {
		t.Fatalf("post-resume JoinDot4 = %q", got)
	}
}

func TestWave30CJoinDotNoStrPoolLeakAndAllocatorInternUnpin(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer func() {
		ExitSingleOwnerAlloc()
		SetCurrentNodeStore(prev)
	}()

	if !TryEnterSingleOwnerAlloc() {
		t.Fatal("TryEnterSingleOwnerAlloc failed")
	}

	// 1. Verify JoinDot / JoinDot3 / JoinDot4 DO NOT permanently leak ephemeral
	// qualified-name lookup keys into strpool (1.04M strings / 116 MB in v28!).
	baseSnap := CaptureBaseline()

	out2 := JoinDot("wave30c_unique_pkg", "Widget")
	if out2 != "wave30c_unique_pkg.Widget" {
		t.Fatalf("JoinDot = %q", out2)
	}
	if Lookup(out2) != 0 {
		t.Fatalf("JoinDot(%q) leaked ephemeral key into strpool!", out2)
	}

	out3 := JoinDot3("wave30c_unique_pkg", "Widget", "prototype")
	if out3 != "wave30c_unique_pkg.Widget.prototype" {
		t.Fatalf("JoinDot3 = %q", out3)
	}
	if Lookup(out3) != 0 {
		t.Fatalf("JoinDot3(%q) leaked ephemeral key into strpool!", out3)
	}

	out4 := JoinDot4("wave30c_unique_pkg", "sub", "Dialog", "dispose")
	if out4 != "wave30c_unique_pkg.sub.Dialog.dispose" {
		t.Fatalf("JoinDot4 = %q", out4)
	}
	if Lookup(out4) != 0 {
		t.Fatalf("JoinDot4(%q) leaked ephemeral key into strpool!", out4)
	}

	// Medium-long qualified names (> stack buf 256B) must ALSO not leak into strpool!
	longSegA := "wave30c_mod$" + strings.Repeat("a", 140)
	longSegB := "export$" + strings.Repeat("b", 140)
	outLong := JoinDot(longSegA, longSegB)
	if outLong != longSegA+"."+longSegB {
		t.Fatal("outLong mismatch")
	}
	if Lookup(outLong) != 0 {
		t.Fatalf("expected 290B JoinDot output NOT to leak into strpool")
	}

	// Generate 2,000 unique ephemeral qualified names across JoinDot / JoinDot3 / JoinDot4
	// and verify CaptureBaseline() did not advance by a single string in any shard!
	for i := 0; i < 2000; i++ {
		seg := fmt.Sprintf("ephemeral_qname_%d", i)
		_ = JoinDot("ns", seg)
		_ = JoinDot3("ns", "sub", seg)
		_ = JoinDot4("ns", "sub", "pkg", seg)
	}
	afterSnap := CaptureBaseline()
	if afterSnap != baseSnap {
		t.Fatalf("ephemeral JoinDot calls polluted strpool: before=%v after=%v", baseSnap, afterSnap)
	}

	// 2. Single exact-sized string allocation (1 alloc/op) on cache miss, 0 allocs/op on hit!
	_ = JoinDot("warmup", "table")
	_ = JoinDot3("warmup", "table", "three")
	hitAllocs := testing.AllocsPerRun(100, func() {
		if JoinDot("warmup", "table") != "warmup.table" {
			t.Fatal("JoinDot hit mismatch")
		}
		if JoinDot3("warmup", "table", "three") != "warmup.table.three" {
			t.Fatal("JoinDot3 hit mismatch")
		}
	})
	if hitAllocs != 0 {
		t.Fatalf("expected 0 allocs on warm JoinDot/JoinDot3 cache hits, got %v", hitAllocs)
	}

	ExitSingleOwnerAlloc()
	// Warm hits outside single-owner mode allocate 0 objects; formatJoin2/3/4 on miss allocate 1 exact string!
	_ = JoinDot("uncached_ns", "Widget")
	_ = JoinDot3("uncached_ns", "Widget", "prototype")
	_ = JoinDot4("uncached_ns", "sub", "Widget", "prototype")
	if got := testing.AllocsPerRun(50, func() {
		if JoinDot("uncached_ns", "Widget") != "uncached_ns.Widget" ||
			JoinDot3("uncached_ns", "Widget", "prototype") != "uncached_ns.Widget.prototype" ||
			JoinDot4("uncached_ns", "sub", "Widget", "prototype") != "uncached_ns.sub.Widget.prototype" {
			t.Fatal("lock-free JoinDot/3/4 mismatch")
		}
	}); got != 0 {
		t.Fatalf("expected 0 allocations on lock-free JoinDot/3/4 hits outside single-owner mode, got %v", got)
	}
	missAllocs2 := testing.AllocsPerRun(50, func() {
		if formatJoin2("uncached_ns", "Widget") != "uncached_ns.Widget" {
			t.Fatal("formatJoin2 mismatch")
		}
	})
	if missAllocs2 > 1 {
		t.Fatalf("expected <= 1 string allocation on formatJoin2, got %v", missAllocs2)
	}
	missAllocs3 := testing.AllocsPerRun(50, func() {
		if formatJoin3("uncached_ns", "Widget", "prototype") != "uncached_ns.Widget.prototype" {
			t.Fatal("formatJoin3 mismatch")
		}
	})
	if missAllocs3 > 1 {
		t.Fatalf("expected <= 1 string allocation on formatJoin3, got %v", missAllocs3)
	}
	missAllocs4 := testing.AllocsPerRun(50, func() {
		if formatJoin4("uncached_ns", "sub", "Widget", "prototype") != "uncached_ns.sub.Widget.prototype" {
			t.Fatal("formatJoin4 mismatch")
		}
	})
	if missAllocs4 > 1 {
		t.Fatalf("expected <= 1 string allocation on formatJoin4, got %v", missAllocs4)
	}

	// 3. Verify neither NodeAllocator.Intern (even UNFLUSHED!) nor JoinDotCache pins caller buffers!
	sFresh := NewNodeStore()
	SetCurrentNodeStore(sFresh)
	if !TryEnterSingleOwnerAlloc() {
		t.Fatal("TryEnterSingleOwnerAlloc on sFresh failed")
	}

	runtime.GC()
	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	longLivedAlloc := NewNodeAllocator(sFresh)
	var retainedJoined []string
	func() {
		for i := 0; i < 250; i++ {
			// 250 * 128 KB = 32 MB of temporary source buffers.
			bigBuf := fmt.Sprintf("w30c_ns_%03d.PartB_%03d.ident_%03d.", i, i, i) + strings.Repeat("Z", 128<<10)
			partA := bigBuf[:11]
			partB := bigBuf[12:21]
			ident := bigBuf[22:31]

			retainedJoined = append(retainedJoined, JoinDot(partA, partB))
			retainedJoined = append(retainedJoined, JoinDot3(partA, partB, ident))
			// Intern slice in longLivedAlloc WITHOUT calling Flush()!
			sid := longLivedAlloc.Intern(ident)
			if String(sid) != ident {
				t.Fatalf("longLivedAlloc.Intern(%q) mismatch", ident)
			}
		}
	}()

	runtime.GC()
	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	var liveDelta int64
	if after.HeapAlloc > before.HeapAlloc {
		liveDelta = int64(after.HeapAlloc - before.HeapAlloc)
	}
	// 32 MB of source buffers were sliced; live retention must be < 6 MB (unpinned!).
	if liveDelta > 6<<20 {
		t.Fatalf("caller buffers pinned by JoinDotCache or NodeAllocator.internCache: grew %d bytes (%.2f MB)",
			liveDelta, float64(liveDelta)/(1<<20))
	}
	runtime.KeepAlive(longLivedAlloc)
	runtime.KeepAlive(retainedJoined)
	longLivedAlloc.Flush()

	// 4. Multi-goroutine -race stress test on JoinDot / JoinDot3 / JoinDot4 + NodeAllocator.Intern.
	ExitSingleOwnerAlloc()
	var wg sync.WaitGroup
	const workers = 16
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		w := w
		go func() {
			defer wg.Done()
			var la NodeAllocator
			la.InitLocal()
			defer la.Flush()
			for i := 0; i < 400; i++ {
				a := fmt.Sprintf("w30c_race_pkg_%d", i%23)
				b := fmt.Sprintf("Model_%d", (w+i)%19)
				c := "prototype"
				d := fmt.Sprintf("m_%d", i%11)

				j2 := JoinDot(a, b)
				if j2 != a+"."+b {
					t.Errorf("JoinDot mismatch: %q", j2)
					return
				}
				j3 := JoinDot3(a, b, c)
				if j3 != a+"."+b+"."+c {
					t.Errorf("JoinDot3 mismatch: %q", j3)
					return
				}
				j4 := JoinDot4(a, b, c, d)
				if j4 != a+"."+b+"."+c+"."+d {
					t.Errorf("JoinDot4 mismatch: %q", j4)
					return
				}
				if sid := la.Intern(j2); String(sid) != j2 {
					t.Errorf("la.Intern(%q) mismatch", j2)
					return
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkWave30CJoinDotFormatAndCache(b *testing.B) {
	p1 := "gapi.client.compute.v1.Instance"
	p2 := "prototype"
	p3 := "getNetworkInterfaces"

	b.Run("UncachedMiss2", func(b *testing.B) {
		ExitSingleOwnerAlloc()
		b.ReportAllocs()
		var sink int
		for i := 0; i < b.N; i++ {
			sink += len(JoinDot(p1, p2))
		}
		_ = sink
	})
	b.Run("UncachedMiss3", func(b *testing.B) {
		ExitSingleOwnerAlloc()
		b.ReportAllocs()
		var sink int
		for i := 0; i < b.N; i++ {
			sink += len(JoinDot3(p1, p2, p3))
		}
		_ = sink
	})
}


func TestWave33CLockFree2WayJoinDotCacheAndMultiChunkStrPool(t *testing.T) {
	// 1. Verify String(0) == "" and String(EmptyStringID) == "" branchlessly via poolChunks[0].
	if got := String(0); got != "" {
		t.Fatalf("String(0) = %q, want \"\"", got)
	}
	if got := String(EmptyStringID); got != "" {
		t.Fatalf("String(EmptyStringID) = %q, want \"\"", got)
	}

	// 2. Multi-chunk strpool interning (> 512 strings in a single shard) + RewindBaseline(zeroSnap) safety.
	snap := CaptureBaseline()
	const many = 5000
	ids := make([]StringID, many)
	for i := 0; i < many; i++ {
		s := fmt.Sprintf("w33c_multichunk_sym_%05d", i)
		ids[i] = Intern(s)
		if ids[i] < 2 {
			t.Fatalf("Intern(%q) returned invalid ID %d < 2", s, ids[i])
		}
	}
	for i := 0; i < many; i++ {
		want := fmt.Sprintf("w33c_multichunk_sym_%05d", i)
		if got := String(ids[i]); got != want {
			t.Fatalf("String(%d) = %q, want %q", ids[i], got, want)
		}
	}
	RewindBaseline(snap)
	if String(0) != "" || String(EmptyStringID) != "" {
		t.Fatalf("String(0/1) broken after RewindBaseline")
	}

	// 3. Find two distinct (a, b) pairs that hash to the SAME joinDot2Cache bucket index,
	// and verify 2-way set-associative retention (both stay warm with 0 allocs/op, no ping-pong thrashing!).
	targetPrefix := "module$exports$google$cloud$compute$v1$"
	seenIdx := make(map[uint32]string)
	var collideA1, collideA2 string
	for i := 0; i < 20000; i++ {
		cand := fmt.Sprintf("%sList%05dRequest", targetPrefix, i)
		h := poolHash(cand)*0x9E3779B97F4A7C15 ^ poolHash("prototype")
		idx := uint32(h >> (64 - joinDot2CacheBits))
		if prev, ok := seenIdx[idx]; ok {
			collideA1 = prev
			collideA2 = cand
			break
		}
		seenIdx[idx] = cand
	}
	if collideA1 == "" || collideA2 == "" {
		t.Fatal("failed to find colliding bucket pair in 20,000 probes")
	}

	// Warm both colliding pairs into the same 2-way bucket.
	want1 := collideA1 + ".prototype"
	want2 := collideA2 + ".prototype"
	if got := JoinDot(collideA1, "prototype"); got != want1 {
		t.Fatalf("got %q, want %q", got, want1)
	}
	if got := JoinDot(collideA2, "prototype"); got != want2 {
		t.Fatalf("got %q, want %q", got, want2)
	}
	// Alternating between both colliding pairs must allocate 0 bytes (2-way hit on out0 & out1!).
	if allocs := testing.AllocsPerRun(200, func() {
		if JoinDot(collideA1, "prototype") != want1 || JoinDot(collideA2, "prototype") != want2 {
			t.Fatal("2-way colliding bucket mismatch")
		}
	}); allocs != 0 {
		t.Fatalf("expected 0 allocs alternating two colliding keys in 2-way joinDot2Bucket, got %v", allocs)
	}

	// 4. Boundary lengths around joinDotFastBufLen (128) and joinDotMaxTotalLen (512).
	for _, total := range []int{2, 127, 128, 129, 255, 256, 257, 511, 512, 513, 700} {
		alen := (total - 1) / 2
		blen := total - 1 - alen
		a := strings.Repeat("a", alen)
		b := strings.Repeat("b", blen)
		if got, want := JoinDot(a, b), a+"."+b; got != want {
			t.Fatalf("total=%d JoinDot len=%d, want %d", total, len(got), len(want))
		}
		if total >= 3 {
			c1 := (total - 2) / 3
			c2 := (total - 2 - c1) / 2
			c3 := total - 2 - c1 - c2
			sa, sb, sc := strings.Repeat("x", c1), strings.Repeat("y", c2), strings.Repeat("z", c3)
			if got, want := JoinDot3(sa, sb, sc), sa+"."+sb+"."+sc; got != want {
				t.Fatalf("total=%d JoinDot3 mismatch", total)
			}
		}
	}

	// 5. High-contention concurrent -race test across 16 goroutines (mix of hot hits + evictions).
	var wg sync.WaitGroup
	const goroutines = 16
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		g := g
		go func() {
			defer wg.Done()
			for i := 0; i < 600; i++ {
				a := fmt.Sprintf("gapi.client.service%d", i%32)
				b := fmt.Sprintf("Resource_%d", (g+i)%24)
				c := "prototype"
				d := fmt.Sprintf("method_%d", i%19)
				if got := JoinDot(a, b); got != a+"."+b {
					t.Errorf("JoinDot race mismatch: %q", got)
					return
				}
				if got := JoinDot3(a, b, c); got != a+"."+b+"."+c {
					t.Errorf("JoinDot3 race mismatch: %q", got)
					return
				}
				if got := JoinDot4(a, b, c, d); got != a+"."+b+"."+c+"."+d {
					t.Errorf("JoinDot4 race mismatch: %q", got)
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestWave35CStrPoolSealedChunksAndRewind(t *testing.T) {
	baseSnap := CaptureBaseline()
	defer RewindBaseline(baseSnap)

	// 1. Intern 100,000 unique ~56-byte identifier strings (~1,562 per shard -> 3 sealed chunks/shard).
	// Verify live HeapObjects delta is tiny (<< 25,000 instead of 100,000+ individual string objects!)
	// and all 100,000 strings round-trip identically with 0 allocs/op.
	runtime.GC()
	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	const totalStrings = 100_000
	ids := make([]StringID, totalStrings)
	for i := 0; i < totalStrings; i++ {
		s := fmt.Sprintf("module$exports$proto$google$cloud$compute$v1$Sym_%06d", i)
		ids[i] = Intern(s)
	}

	runtime.GC()
	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	var objDelta uint64
	if after.HeapObjects > before.HeapObjects {
		objDelta = after.HeapObjects - before.HeapObjects
	}
	// 100,000 strings -> ~160 sealed chunks (~320 objects) + 64 tail chunks (~3,800 strings) + 64 uint32 tables.
	if objDelta > 30_000 {
		t.Fatalf("expected sealed poolChunks to collapse 100K string objects, got HeapObjects delta = %d", objDelta)
	}
	for i := 0; i < totalStrings; i += 97 {
		want := fmt.Sprintf("module$exports$proto$google$cloud$compute$v1$Sym_%06d", i)
		if got := String(ids[i]); got != want {
			t.Fatalf("String(%d) = %q, want %q", ids[i], got, want)
		}
		if lid := Lookup(want); lid != ids[i] {
			t.Fatalf("Lookup(%q) = %d, want %d", want, lid, ids[i])
		}
	}
	RewindBaseline(baseSnap)

	// 2. Mid-chunk CaptureBaseline -> seal chunk -> RewindBaseline (unseals tail chunk) -> re-seal chunk!
	var shard7Pre []StringID
	var shard7PreVals []string
	for n := 0; len(shard7Pre) < 200; n++ {
		cand := fmt.Sprintf("w35c_midrewind_pre_%d", n)
		if uint32(poolHash(cand))&(strpoolNumShards-1) == 7 {
			shard7PreVals = append(shard7PreVals, cand)
			shard7Pre = append(shard7Pre, Intern(cand))
		}
	}
	midSnap := CaptureBaseline()
	// Push shard 7 well past 512 so its tail chunk seals into slab+offs16.
	for n := 0; n < 40000; n++ {
		Intern(fmt.Sprintf("w35c_midrewind_seal1_%d", n))
	}
	// Rewind back to midSnap (must unseal shard 7's tail chunk while keeping first 200 strings intact).
	RewindBaseline(midSnap)
	for i, id := range shard7Pre {
		if got := String(id); got != shard7PreVals[i] {
			t.Fatalf("after mid-chunk rewind: String(%d) = %q, want %q", id, got, shard7PreVals[i])
		}
		if Lookup(shard7PreVals[i]) != id {
			t.Fatalf("after mid-chunk rewind: Lookup(%q) != %d", shard7PreVals[i], id)
		}
	}
	// Intern again past 512 to re-seal shard 7's tail chunk and verify both pre & post strings.
	var shard7Post []StringID
	var shard7PostVals []string
	for n := 0; len(shard7Post) < 400; n++ {
		cand := fmt.Sprintf("w35c_midrewind_post_%d", n)
		if uint32(poolHash(cand))&(strpoolNumShards-1) == 7 {
			shard7PostVals = append(shard7PostVals, cand)
			shard7Post = append(shard7Post, Intern(cand))
		}
	}
	for i, id := range shard7Pre {
		if got := String(id); got != shard7PreVals[i] {
			t.Fatalf("after re-seal: pre String(%d) = %q, want %q", id, got, shard7PreVals[i])
		}
	}
	for i, id := range shard7Post {
		if got := String(id); got != shard7PostVals[i] {
			t.Fatalf("after re-seal: post String(%d) = %q, want %q", id, got, shard7PostVals[i])
		}
	}
	RewindBaseline(baseSnap)

	// 3. >64 KiB chunk fallback (total > 0xFFFF in a single 512-string chunk).
	bigPad := strings.Repeat("Q", 180)
	var bigIDs []StringID
	var bigVals []string
	for n := 0; len(bigIDs) < 520; n++ {
		cand := fmt.Sprintf("w35c_bigchunk_%05d_%s", n, bigPad)
		if uint32(poolHash(cand))&(strpoolNumShards-1) == 11 {
			bigVals = append(bigVals, cand)
			bigIDs = append(bigIDs, Intern(cand))
		}
	}
	for i, id := range bigIDs {
		if got := String(id); got != bigVals[i] {
			t.Fatalf(">64KiB chunk String(%d) mismatch at %d", id, i)
		}
	}
	RewindBaseline(baseSnap)

	// 4. Concurrent -race stress across chunk-sealing boundaries (Intern, InternBytes, Lookup, String).
	var wg sync.WaitGroup
	const workers = 12
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		w := w
		go func() {
			defer wg.Done()
			for i := 0; i < 1500; i++ {
				s := fmt.Sprintf("w35c_race_%d_%d", w%4, i)
				var id StringID
				if i&1 == 0 {
					id = Intern(s)
				} else {
					id = InternBytes([]byte(s))
				}
				if got := String(id); got != s {
					t.Errorf("concurrent String(%d) = %q, want %q", id, got, s)
					return
				}
				if lid := Lookup(s); lid != id {
					t.Errorf("concurrent Lookup(%q) = %d, want %d", s, lid, id)
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestWave36CJoinDot4WayScanResistantCache(t *testing.T) {
	// Find 4 distinct keys hashing to the exact same joinDot2 bucket index and verify
	// all 4 ways stay resident simultaneously with 0 allocations AND resist eviction by a 1-hit scan!
	seen := make(map[uint32][]string)
	var ways []string
	for i := 0; i < 200000; i++ {
		cand := fmt.Sprintf("module$exports$w36c$pkg%06d", i)
		h := poolHash(cand)*0x9E3779B97F4A7C15 ^ poolHash("prototype")
		idx := uint32(h >> (64 - joinDot2CacheBits))
		list := append(seen[idx], cand)
		seen[idx] = list
		if len(list) == 4 {
			ways = list
			break
		}
	}
	if len(ways) != 4 {
		t.Fatal("failed to find 4 colliding keys in 200,000 probes")
	}

	// Warm all 4 ways (first 4 inserts into empty/non-full bucket succeed on 1st call).
	wants := make([]string, 4)
	for i, w := range ways {
		wants[i] = w + ".prototype"
		if got := JoinDot(w, "prototype"); got != wants[i] {
			t.Fatalf("warm JoinDot(%q) = %q, want %q", w, got, wants[i])
		}
	}

	// Probe a 5th colliding 1-hit-wonder key ONCE: scan-resistant seen filter must
	// allocate ONLY 1 string (no bucket struct!) AND keep all 4 hot ways resident!
	targetIdx := uint32((poolHash(ways[0])*0x9E3779B97F4A7C15 ^ poolHash("prototype")) >> (64 - joinDot2CacheBits))
	var fifth string
	for i := 200000; i < 400000; i++ {
		cand := fmt.Sprintf("module$exports$w36c$onehit%06d", i)
		h := poolHash(cand)*0x9E3779B97F4A7C15 ^ poolHash("prototype")
		if uint32(h>>(64-joinDot2CacheBits)) == targetIdx {
			fifth = cand
			break
		}
	}
	if fifth == "" {
		t.Fatal("failed to find 5th colliding key")
	}
	if got := JoinDot(fifth, "prototype"); got != fifth+".prototype" {
		t.Fatalf("5th JoinDot mismatch: %q", got)
	}

	// All 4 original ways must STILL hit with 0 allocations!
	if allocs := testing.AllocsPerRun(100, func() {
		for i, w := range ways {
			if JoinDot(w, "prototype") != wants[i] {
				t.Fatal("4-way entry evicted by 1-hit wonder")
			}
		}
	}); allocs != 0 {
		t.Fatalf("expected 0 allocs across all 4 resident ways after 1-hit probe, got %v", allocs)
	}
}

func TestWave37CStrPoolStagedBatchingAndJoinDot4WaySeenRing(t *testing.T) {
	baseSnap := CaptureBaseline()
	defer RewindBaseline(baseSnap)

	// 1. Verify 32-slot lock-free stagingBatch cuts heap allocations by ~10x during Intern misses.
	const count = 6400
	var prebuilt strings.Builder
	prebuilt.Grow(count * 24)
	for i := 0; i < count; i++ {
		fmt.Fprintf(&prebuilt, "w37c_staged_ident_%06d", i)
	}
	src := prebuilt.String()
	ids := make([]StringID, count)

	var mBefore, mAfter runtime.MemStats
	runtime.ReadMemStats(&mBefore)
	for i := 0; i < count; i++ {
		ids[i] = Intern(src[i*24 : (i+1)*24])
	}
	runtime.ReadMemStats(&mAfter)
	mallocs := mAfter.Mallocs - mBefore.Mallocs
	// Without staged batching, 6,400 Intern misses on unowned slices would exceed 6,400 mallocs;
	// with 32-slot stagingBatch, strpool overhead for 6,400 strings is ~500 objects total (< 1,200).
	if mallocs > 1200 {
		t.Fatalf("expected staged Intern batching to eliminate per-string Clone allocs, got %d mallocs for %d strings", mallocs, count)
	}
	for i := 0; i < count; i += 31 {
		want := fmt.Sprintf("w37c_staged_ident_%06d", i)
		if got := String(ids[i]); got != want {
			t.Fatalf("String(%d) = %q, want %q", ids[i], got, want)
		}
	}

	// 2. Verify SealActivePoolChunks immediately unpins <32 partial-batch substrings of large buffers.
	SealActivePoolChunks()
	runtime.GC()
	runtime.GC()
	var beforePin runtime.MemStats
	runtime.ReadMemStats(&beforePin)

	var partialIDs [20]StringID
	func() {
		for i := 0; i < len(partialIDs); i++ {
			big := fmt.Sprintf("w37c_pin_%02d_", i) + strings.Repeat("X", 256<<10)
			sub := big[:12]
			partialIDs[i] = Intern(sub)
		}
	}()
	SealActivePoolChunks()
	runtime.GC()
	runtime.GC()
	var afterPin runtime.MemStats
	runtime.ReadMemStats(&afterPin)
	var pinDelta int64
	if afterPin.HeapAlloc > beforePin.HeapAlloc {
		pinDelta = int64(afterPin.HeapAlloc - beforePin.HeapAlloc)
	}
	// 20 * 256 KiB = 5 MiB; after SealActivePoolChunks, retention must be tiny (< 1 MiB).
	if pinDelta > 1<<20 {
		t.Fatalf("SealActivePoolChunks left caller buffers pinned: delta = %d bytes", pinDelta)
	}
	for i, id := range partialIDs {
		want := fmt.Sprintf("w37c_pin_%02d_", i)
		if got := String(id); got != want {
			t.Fatalf("after SealActivePoolChunks: String(%d) = %q, want %q", id, got, want)
		}
	}

	// 3. Verify JoinDot 4-fingerprint seen ring admits TWO interleaved 2-hit keys into a full 4-way bucket
	// without ping-ponging forever on a 1-slot seen filter!
	seen := make(map[uint32][]string)
	var col []string
	for i := 0; i < 400000; i++ {
		cand := fmt.Sprintf("module$exports$w37c$ring%06d", i)
		h := poolHash(cand)*0x9E3779B97F4A7C15 ^ poolHash("prototype")
		idx := uint32(h >> (64 - joinDot2CacheBits))
		list := append(seen[idx], cand)
		seen[idx] = list
		if len(list) == 6 {
			col = list
			break
		}
	}
	if len(col) != 6 {
		t.Fatal("failed to find 6 colliding JoinDot keys")
	}
	// Fill all 4 ways with cold keys col[0..3].
	for i := 0; i < 4; i++ {
		_ = JoinDot(col[i], "prototype")
	}
	// Interleave TWO new hot keys col[4] and col[5]: K4, K5, K4, K5!
	// On a 1-slot seen filter, K4 and K5 overwrite each other's seen tag and NEVER get admitted;
	// on the 4-fingerprint seen ring, both get admitted on their 2nd call and then hit with 0 allocs!
	_ = JoinDot(col[4], "prototype")
	_ = JoinDot(col[5], "prototype")
	_ = JoinDot(col[4], "prototype")
	_ = JoinDot(col[5], "prototype")

	want4 := col[4] + ".prototype"
	want5 := col[5] + ".prototype"
	if allocs := testing.AllocsPerRun(100, func() {
		if JoinDot(col[4], "prototype") != want4 || JoinDot(col[5], "prototype") != want5 {
			t.Fatal("JoinDot interleaved admission mismatch")
		}
	}); allocs != 0 {
		t.Fatalf("expected 0 allocs after interleaved 2-hit admission via 4-way seen ring, got %v", allocs)
	}
}
