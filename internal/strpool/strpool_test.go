package strpool

import (
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestInternAndPointers(t *testing.T) {
	if String(0) != "" || StringPtr(0) != nil {
		t.Fatalf("ID 0 should represent nil/empty")
	}
	emptyID := Intern("")
	if emptyID != EmptyID || String(emptyID) != "" || StringPtr(emptyID) == nil || *StringPtr(emptyID) != "" {
		t.Fatalf("empty string handling broken: id=%d", emptyID)
	}
	if InternBytes(nil) != EmptyID || Lookup("") != EmptyID {
		t.Fatalf("empty InternBytes/Lookup broken")
	}
	if HashString("hello") != HashBytes([]byte("hello")) || HashString("") != HashBytes(nil) {
		t.Fatalf("HashString/HashBytes mismatch")
	}

	large := strings.Repeat("x", 40_000)
	id1 := Intern(large)
	id2 := InternBytes([]byte(large))
	if id1 != id2 || Lookup(large) != id1 || String(id1) != large {
		t.Fatalf("large string roundtrip failed")
	}
	if StringPtr(id1) != StringPtr(id2) || *StringPtr(id1) != large {
		t.Fatalf("StringPtr identity failed")
	}

	raw := []byte("zero_alloc_intern_bytes_probe")
	wantID := Intern(string(raw))
	if allocs := testing.AllocsPerRun(100, func() {
		if InternBytes(raw) != wantID {
			t.Fatal("InternBytes mismatch")
		}
	}); allocs != 0 {
		t.Fatalf("InternBytes hit allocated %v times, want 0", allocs)
	}
}

func TestConcurrentInternAndRead(t *testing.T) {
	var wg sync.WaitGroup
	for g := range 8 {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := range 2000 {
				s := "ident_" + strconv.Itoa((worker*500+i)%3000)
				id := Intern(s)
				if got := String(id); got != s {
					t.Errorf("String(%d) = %q, want %q", id, got, s)
					return
				}
				p := StringPtr(id)
				if p == nil || *p != s {
					t.Errorf("StringPtr(%d) mismatch", id)
					return
				}
				if InternBytes([]byte(s)) != id {
					t.Errorf("InternBytes(%q) != %d", s, id)
					return
				}
			}
		}(g)
	}
	wg.Wait()
}
