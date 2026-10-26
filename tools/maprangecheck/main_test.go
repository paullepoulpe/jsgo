package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRepoMapRanges is the guard: every map iteration in internal/ must say
// why its order cannot reach output (see the package comment).
func TestRepoMapRanges(t *testing.T) {
	if testing.Short() {
		t.Skip("type-checks internal/ (builds its export data)")
	}
	sites, err := Check(filepath.Join("..", ".."), DefaultPackages)
	if err != nil {
		t.Fatal(err)
	}
	if len(sites) == 0 {
		t.Fatal("found no map iterations at all; is the checker broken?")
	}
	for _, s := range sites {
		if s.Reason == "" {
			t.Errorf("%s: %s without a %q comment; keep Java's order (an insertion-ordered slice, or sorted keys) or annotate why the order cannot matter", s.Pos, s.What, "// "+Marker+" <reason>")
		}
	}
}

// TestCheckFixture checks the checker on testdata/mod: the lines marked
// "// want" are the unannotated sites, and every other map iteration there
// is annotated.
func TestCheckFixture(t *testing.T) {
	dir := filepath.Join("testdata", "mod")
	t.Setenv("GOWORK", "off")
	sites, err := Check(dir, []string{"./..."})
	if err != nil {
		t.Fatal(err)
	}
	want := map[int]bool{}
	f, err := os.Open(filepath.Join(dir, "p", "p.go"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for line := 1; sc.Scan(); line++ {
		if strings.Contains(sc.Text(), "// want") {
			want[line] = true
		}
	}
	got := map[int]bool{}
	annotated := 0
	for _, s := range sites {
		if s.Reason == "" {
			got[s.Pos.Line] = true
		} else {
			annotated++
		}
	}
	for line := range want { // order-independent: test assertions only
		if !got[line] {
			t.Errorf("line %d: not reported", line)
		}
	}
	for line := range got { // order-independent: test assertions only
		if !want[line] {
			t.Errorf("line %d: reported, but not marked // want", line)
		}
	}
	if annotated != 2 {
		t.Errorf("annotated sites = %d, want 2", annotated)
	}
}
