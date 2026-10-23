package sourcemap

import (
	"math"
	"math/rand"
	"testing"
)

// Go-only: the delta-encoded mapping store must give back exactly the
// mappings that were added, in order, across block boundaries.
func TestMappingStoreRoundTrip(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	field := func() int32 {
		switch rnd.Intn(6) {
		case 0:
			return math.MaxInt32
		case 1:
			return math.MinInt32
		case 2:
			return -1
		default:
			return int32(rnd.Intn(100000))
		}
	}
	var s sourceMapGeneratorV3_mappingStore
	var want []SourceMapGeneratorV3_Mapping
	for i := 0; i < 50000; i++ {
		m := SourceMapGeneratorV3_Mapping{
			originalNameIdx: field(), sourceFileIdx: field(),
			originalLine: field(), originalCol: field(),
			startLine: field(), startCol: field(),
			endLine: field(), endCol: field(),
		}
		s.add(&m)
		want = append(want, m)
	}
	if len(s.blocks) == 0 {
		t.Fatal("expected several blocks")
	}
	r := sourceMapGeneratorV3_mappingReader{store: &s}
	var got SourceMapGeneratorV3_Mapping
	for i, w := range want {
		if !r.next(&got) {
			t.Fatalf("only %d mappings read back, want %d", i, len(want))
		}
		if got != w {
			t.Fatalf("mapping %d = %+v, want %+v", i, got, w)
		}
	}
	if r.next(&got) {
		t.Fatal("read back more mappings than were added")
	}
	s.reset()
	if r := (sourceMapGeneratorV3_mappingReader{store: &s}); r.next(&got) {
		t.Fatal("reset store is not empty")
	}
}
