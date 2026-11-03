package jscomp

import (
	"slices"
	"strconv"
	"testing"
)

func TestOrderedMapAndSetHelpers(t *testing.T) {
	var nilMap *OrderedMap[string, int]
	nilMap.Clear()
	if idx, ok := nilMap.Index("a"); ok || idx != -1 {
		t.Fatalf("nilMap.Index = (%d, %v), want (-1, false)", idx, ok)
	}

	var nilSet *OrderedSet[string]
	nilSet.Clear()
	if idx, ok := nilSet.Index("a"); ok || idx != -1 {
		t.Fatalf("nilSet.Index = (%d, %v), want (-1, false)", idx, ok)
	}

	var s OrderedSet[string]
	// Exercise both linear mode (<= orderedMapIndexThreshold) and indexed mode (> threshold).
	for i := range 16 {
		key := "k" + strconv.Itoa(i)
		if got := s.Intern(key); got != i {
			t.Fatalf("first Intern(%q) = %d, want %d", key, got, i)
		}
		if got := s.Intern(key); got != i {
			t.Fatalf("repeat Intern(%q) = %d, want %d", key, got, i)
		}
		if idx, ok := s.Index(key); !ok || idx != i {
			t.Fatalf("Index(%q) = (%d, %v), want (%d, true)", key, idx, ok, i)
		}
	}
	if idx, ok := s.Index("missing"); ok || idx != -1 {
		t.Fatalf("Index(missing) = (%d, %v), want (-1, false)", idx, ok)
	}

	// Verify Index/Intern stay consistent across removals and hole compaction.
	if !s.Remove("k0") || !s.Remove("k5") {
		t.Fatalf("expected k0 and k5 removal to succeed")
	}
	wantSlice := s.ToSlice()
	for wantIdx, key := range wantSlice {
		if gotIdx, ok := s.Index(key); !ok || gotIdx != wantIdx {
			t.Fatalf("after Remove, Index(%q) = (%d, %v), want (%d, true)", key, gotIdx, ok, wantIdx)
		}
		if gotIdx := s.Intern(key); gotIdx != wantIdx {
			t.Fatalf("after Remove, Intern(%q) = %d, want %d", key, gotIdx, wantIdx)
		}
	}
	newIdx := s.Intern("tail")
	if newIdx != len(wantSlice) {
		t.Fatalf("Intern(tail) = %d, want %d", newIdx, len(wantSlice))
	}

	s.Clear()
	if !s.IsEmpty() || s.Len() != 0 || s.ToSlice() != nil {
		t.Fatalf("after Clear: empty=%v len=%d slice=%v", s.IsEmpty(), s.Len(), s.ToSlice())
	}
	if got := s.Intern("fresh"); got != 0 {
		t.Fatalf("after Clear, Intern(fresh) = %d, want 0", got)
	}
	if !slices.Equal(s.ToSlice(), []string{"fresh"}) {
		t.Fatalf("after Clear + Intern, slice = %v, want [fresh]", s.ToSlice())
	}

	var m OrderedMap[string, string]
	if got := m.Intern("alpha"); got != 0 {
		t.Fatalf("m.Intern(alpha) = %d, want 0", got)
	}
	m.Put("alpha", "A")
	if got := m.Intern("beta"); got != 1 {
		t.Fatalf("m.Intern(beta) = %d, want 1", got)
	}
	if got := m.Intern("alpha"); got != 0 || m.Get("alpha") != "A" {
		t.Fatalf("repeat m.Intern(alpha) = %d (val=%q), want 0 (val=A)", got, m.Get("alpha"))
	}
	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("m.Len() after Clear = %d, want 0", m.Len())
	}
}
