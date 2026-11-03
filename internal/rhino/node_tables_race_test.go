package rhino

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
)

func TestNewNodeStoreEmptyHandleTablesZeroHeapAlloc(t *testing.T) {
	s := NewNodeStore()
	if got := s.types.count.Load(); got != 0 {
		t.Fatalf("empty s.types.count = %d, want 0", got)
	}
	if got := s.colors.count.Load(); got != 0 {
		t.Fatalf("empty s.colors.count = %d, want 0", got)
	}
	for i := range s.types.dir {
		if s.types.dir[i].Load() != nil {
			t.Fatalf("expected s.types.dir[%d] to be nil on empty store", i)
		}
		if s.colors.dir[i].Load() != nil {
			t.Fatalf("expected s.colors.dir[%d] to be nil on empty store", i)
		}
	}

	var emptyTbl handleTable[JSTypeBase]
	if emptyTbl.peek(0) != nil || emptyTbl.peek(1) != nil {
		t.Fatal("expected nil peek on empty handleTable")
	}
	visited := 0
	emptyTbl.each(func(*atomic.Pointer[JSTypeBase]) {
		visited++
	})
	if visited != 0 {
		t.Fatalf("expected 0 visited slots on empty handleTable, got %d", visited)
	}

	// Verify 100 empty NewNodeStore() calls do not eagerly allocate 100 * 256 KiB (25.6 MiB).
	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	stores := make([]*NodeStore, 100)
	for i := range stores {
		stores[i] = NewNodeStore()
	}
	runtime.ReadMemStats(&after)
	runtime.KeepAlive(stores)

	delta := after.TotalAlloc - before.TotalAlloc
	const maxAllowedBytes = 2 << 20 // < 2 MiB (even under -race), vs 25.6 MiB previously
	if delta > maxAllowedBytes {
		t.Fatalf("100 NewNodeStore() calls allocated %d bytes (> %d); eager handleTable allocation regressed", delta, maxAllowedBytes)
	}
}

func TestHandleTableConcurrentSetGetAndReleaseJSType(t *testing.T) {
	const (
		numWorkers     = 8
		nodesPerWorker = 600
		totalTypes     = tableBlockLen + 300
	)

	regKeep := NewJSTypeRegistry(nil)
	regRelease := NewJSTypeRegistry(nil)

	keepTypes := make([]JSType, totalTypes)
	releaseTypes := make([]JSType, totalTypes)
	for i := range totalTypes {
		keepTypes[i] = NewRecordTypeBuilder(regKeep).
			AddProperty(fmt.Sprintf("k_%d", i), regKeep.GetNativeType(JSTypeNative_NUMBER_TYPE), NilNode).
			Build()
		releaseTypes[i] = NewRecordTypeBuilder(regRelease).
			AddProperty(fmt.Sprintf("r_%d", i), regRelease.GetNativeType(JSTypeNative_STRING_TYPE), NilNode).
			Build()
	}

	// Seed stale handles in a prior store so JSTypeBase.nodeHandle holds non-zero
	// handles 1..N when workers race on s2.types.peek(h) vs s2.types.add().
	seedStore := NewNodeStore()
	prev := SetCurrentNodeStore(seedStore)
	defer SetCurrentNodeStore(prev)
	seedNode := IR_Name("seed")
	for i := range totalTypes {
		seedNode.SetJSType(keepTypes[i])
		seedNode.SetJSType(releaseTypes[i])
	}

	store := NewNodeStore()
	SetCurrentNodeStore(store)

	workerKeepNodes := make([][]Node, numWorkers)
	workerReleaseNodes := make([][]Node, numWorkers)
	for w := range numWorkers {
		workerKeepNodes[w] = make([]Node, nodesPerWorker)
		workerReleaseNodes[w] = make([]Node, nodesPerWorker)
		for i := range nodesPerWorker {
			workerKeepNodes[w][i] = IR_Name("keep")
			workerReleaseNodes[w][i] = IR_Name("rel")
		}
	}

	start := make(chan struct{})
	var wg sync.WaitGroup
	for w := range numWorkers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			<-start
			kNodes := workerKeepNodes[workerID]
			rNodes := workerReleaseNodes[workerID]
			for i := range nodesPerWorker {
				tIdx := (workerID*nodesPerWorker + i) % totalTypes
				kt := keepTypes[tIdx]
				rt := releaseTypes[tIdx]

				kNodes[i].SetJSType(kt)
				rNodes[i].SetJSType(rt)

				if got := kNodes[i].GetJSType(); got != kt {
					t.Errorf("worker %d node %d: GetJSType() = %v, want %v", workerID, i, got, kt)
					return
				}
				// Concurrent ReleaseNodeTypeHandles may clear rt to nil, or return rt; never anything else.
				if got := rNodes[i].GetJSType(); got != nil && got != rt {
					t.Errorf("worker %d rel node %d: unexpected GetJSType() = %v", workerID, i, got)
					return
				}
			}
		}(w)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		runtime.Gosched()
		regRelease.ReleaseNodeTypeHandles()
	}()

	close(start)
	wg.Wait()

	// Final release after all workers finish attaching releaseTypes.
	regRelease.ReleaseNodeTypeHandles()

	for w := range numWorkers {
		for i := range nodesPerWorker {
			tIdx := (w*nodesPerWorker + i) % totalTypes
			if got := workerKeepNodes[w][i].GetJSType(); got != keepTypes[tIdx] {
				t.Fatalf("post-release keep node [%d][%d]: got %v, want %v", w, i, got, keepTypes[tIdx])
			}
			if got := workerReleaseNodes[w][i].GetJSType(); got != nil {
				t.Fatalf("post-release released node [%d][%d]: got %v, want nil", w, i, got)
			}
		}
	}
}

func TestHandleTableConcurrentColorsAndSnapshotRewind(t *testing.T) {
	const (
		numWorkers     = 8
		nodesPerWorker = 600
		totalColors    = tableBlockLen + 400
	)

	palette := make([]*colors.Color, totalColors)
	for i := range totalColors {
		palette[i] = colors.Color_SingleBuilder().
			SetId(colors.ColorId_FromUnsigned(int64(10_000 + i))).
			Build()
	}

	// Pre-populate stale handles in an old store.
	oldStore := NewNodeStore()
	prev := SetCurrentNodeStore(oldStore)
	defer SetCurrentNodeStore(prev)
	seed := IR_Name("seed")
	for _, c := range palette {
		seed.SetColor(c)
	}

	store := NewNodeStore()
	SetCurrentNodeStore(store)

	workerNodes := make([][]Node, numWorkers)
	for w := range numWorkers {
		workerNodes[w] = make([]Node, nodesPerWorker)
		for i := range nodesPerWorker {
			workerNodes[w][i] = IR_Name("colored")
		}
	}

	start := make(chan struct{})
	var wg sync.WaitGroup
	for w := range numWorkers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			<-start
			nodes := workerNodes[workerID]
			for i := range nodesPerWorker {
				c := palette[(workerID*nodesPerWorker+i)%totalColors]
				nodes[i].SetColor(c)
				if got := nodes[i].GetColor(); got != c {
					t.Errorf("worker %d node %d: GetColor() = %v, want %v", workerID, i, got, c)
					return
				}
			}
		}(w)
	}
	close(start)
	wg.Wait()

	snap := store.Snapshot()
	extraCol := colors.Color_SingleBuilder().SetId(colors.ColorId_FromUnsigned(999_999)).Build()
	workerNodes[0][0].SetColor(extraCol)
	store.RestoreSnapshot(snap)

	want0 := palette[0]
	if got := workerNodes[0][0].GetColor(); got != want0 {
		t.Fatalf("after RestoreSnapshot: GetColor() = %v, want %v", got, want0)
	}
}

func TestStandaloneStoreLockMutualExclusion(t *testing.T) {
	prev := SetCurrentNodeStore(nil)
	defer SetCurrentNodeStore(prev)

	const goroutines = 8
	const iters = 25
	active := 0
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range iters {
				LockStandaloneStore()
				active++
				if active != 1 {
					t.Errorf("concurrent critical section entry: active = %d", active)
				}
				_ = EnsureCompilationStore()
				_ = IR_Name("standalone")
				ResetNodeStore()
				active--
				UnlockStandaloneStore()
			}
		}()
	}
	wg.Wait()
}
