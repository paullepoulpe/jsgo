package nodestore

import (
	"sync"
	"testing"
)

func TestPlainUint32ColumnsConcurrentPageAndBlockPublication(t *testing.T) {
	s := NewNodeStore()
	prev := SetCurrentNodeStore(s)
	defer SetCurrentNodeStore(prev)

	// Pre-allocate nodes spanning multiple 16,384-node blocks and 64-node pages,
	// leaving half untouched (must read back as 0/nil even while neighbors allocate
	// columns/pages concurrently) and mutating interleaved slots across 16 goroutines.
	const total = 40000
	nodes := make([]Node, total)
	for i := range nodes {
		nodes[i] = AllocNodeToken(uint8(i%5), uint8(i%200))
	}

	tc := &fakeTCHolder{name: "StringType"}
	doc := &fakeJSDocHolder{comment: "/** @const */"}
	sf := &fakeSourceFileHolder{name: "mod.js"}

	const workers = 16
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		w := w
		go func() {
			defer wg.Done()
			for i := w; i < total; i += workers {
				n := nodes[i]
				if i%2 == 0 {
					// Read zero-initialized neighbors on unallocated/newly-allocated columns & pages.
					if n.TypeOrColor() != nil || n.Length() != 0 || n.OrigNameID() != 0 || n.JSDocObj() != nil || n.SourceFileObj() != nil {
						t.Errorf("untouched node %d non-zero", i)
						return
					}
					continue
				}
				n.SetTypeOrColor(tc)
				n.SetLength(int32(i * 3))
				n.SetOrigNameID(StringID(i + 100))
				n.SetJSDocObj(doc)
				n.SetSourceFileObj(sf)
				if n.TypeOrColor() != tc || n.Length() != int32(i*3) || n.OrigNameID() != StringID(i+100) || n.JSDocObj() != doc || n.SourceFileObj() != sf {
					t.Errorf("node %d readback mismatch", i)
					return
				}
			}
		}()
	}
	wg.Wait()

	// Parallel read-only verification across all workers (mirrors CodePrinter parallel reads).
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		w := w
		go func() {
			defer wg.Done()
			for i := w; i < total; i += workers {
				n := nodes[i]
				if i%2 == 0 {
					if n.TypeOrColor() != nil || n.Length() != 0 || n.OrigNameID() != 0 || n.JSDocObj() != nil || n.SourceFileObj() != nil {
						t.Errorf("verify untouched node %d non-zero", i)
						return
					}
				} else {
					if n.TypeOrColor() != tc || n.Length() != int32(i*3) || n.OrigNameID() != StringID(i+100) || n.JSDocObj() != doc || n.SourceFileObj() != sf {
						t.Errorf("verify node %d mismatch", i)
						return
					}
				}
			}
		}()
	}
	wg.Wait()
}
