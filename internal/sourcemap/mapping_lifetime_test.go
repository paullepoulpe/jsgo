package sourcemap_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

type retainedMappingVisitor struct {
	mappings []*sourcemap.SourceMapGeneratorV3_Mapping
}

func (v *retainedMappingVisitor) Visit(m *sourcemap.SourceMapGeneratorV3_Mapping, _, _, _, _ int) {
	if m != nil {
		v.mappings = append(v.mappings, m)
	}
}

func TestSourceMapFreshMappingPublicWriteEntry(t *testing.T) {
	g := sourcemap.NewSourceMapGeneratorV3()
	var out strings.Builder
	mapper := sourcemap.NewSourceMapGeneratorV3_LineMapper(g, &out, 1)
	mapper.WriteEntry(sourcemap.NewSourceMapGeneratorV3_Mapping(), 0)
	if out.String() != "AAAA" || g.GetSourceId("") != 0 {
		t.Fatalf("default Mapping changed: %q", out.String())
	}
}

func TestSourceMapRetainedMappingResetCopyAndOtherGenerator(t *testing.T) {
	g := sourcemap.NewSourceMapGeneratorV3()
	for i := 0; i < 100; i++ {
		g.AddMappingPositions(fmt.Sprintf("source%d.js", i), nil,
			sourcemap.MakeFilePosition(i, i+1), sourcemap.MakeFilePosition(0, 2*i), sourcemap.MakeFilePosition(0, 2*i+1))
	}
	visitor := &retainedMappingVisitor{}
	sourcemap.NewSourceMapGeneratorV3_MappingTraversal(g).Traverse(visitor)
	if len(visitor.mappings) != 100 {
		t.Fatalf("expected100 nonempty segments, got%d", len(visitor.mappings))
	}
	// Dynamic mappings remain owned by their pointers; copies carry the same
	// source value rather than an index into whichever generator replays them.
	indices := []int{17, 81}
	copies := make([]sourcemap.SourceMapGeneratorV3_Mapping, len(indices))
	for i, index := range indices {
		copies[i] = *visitor.mappings[index]
	}
	g.Reset()
	// Exercise reuse of both inline and dynamic storage before replaying the
	// old dynamic records. The old initial16 cells were already Reset-owned.
	for i := 0; i < 100; i++ {
		g.AddMappingPositions("replacement.js", nil, sourcemap.MakeFilePosition(0, i),
			sourcemap.MakeFilePosition(0, 2*i), sourcemap.MakeFilePosition(0, 2*i+1))
	}
	for i, index := range indices {
		for _, m := range []*sourcemap.SourceMapGeneratorV3_Mapping{visitor.mappings[index], &copies[i]} {
			other := sourcemap.NewSourceMapGeneratorV3()
			var replay, want strings.Builder
			sourcemap.NewSourceMapGeneratorV3_LineMapper(other, &replay, 1).WriteEntry(m, 0)
			for _, value := range []int32{0, 0, int32(index), int32(index + 1)} {
				sourcemap.Base64VLQ_Encode(&want, value)
			}
			if replay.String() != want.String() || other.GetSourceId(fmt.Sprintf("source%d.js", index)) != 0 {
				t.Fatalf("retained/copy mapping%d changed source or coordinates", index)
			}
		}
		var replay strings.Builder
		sourcemap.NewSourceMapGeneratorV3_LineMapper(g, &replay, 1).WriteEntry(visitor.mappings[index], 0)
		if g.GetSourceId(fmt.Sprintf("source%d.js", index)) != i {
			t.Fatal("Reset generator did not import retained source names in replay order")
		}
	}
}

func TestSourceMapJavaInt32PositionAndPrefixArithmetic(t *testing.T) {
	g := sourcemap.NewSourceMapGeneratorV3()
	// The original source coordinates may reach Java's int bounds, including
	// a negative column; writing an entry must preserve their VLQ32 arithmetic.
	g.AddMappingPositions("bounds.js", nil, sourcemap.MakeFilePosition(math.MaxInt32, math.MinInt32),
		sourcemap.MakeFilePosition(0, 0), sourcemap.MakeFilePosition(0, 1))
	visitor := &retainedMappingVisitor{}
	sourcemap.NewSourceMapGeneratorV3_MappingTraversal(g).Traverse(visitor)
	if len(visitor.mappings) != 1 {
		t.Fatal("boundary mapping was dropped")
	}
	var got strings.Builder
	sourcemap.NewSourceMapGeneratorV3_LineMapper(g, &got, 1).WriteEntry(visitor.mappings[0], 0)
	var want strings.Builder
	for _, value := range []int32{0, 0, math.MaxInt32, math.MinInt32} {
		sourcemap.Base64VLQ_Encode(&want, value)
	}
	if got.String() != want.String() {
		t.Fatalf("Java32 source position changed: %q != %q", got.String(), want.String())
	}
	traversal := sourcemap.NewSourceMapGeneratorV3_MappingTraversal(g)
	g.SetWrapperPrefix("x")
	if traversal.GetAdjustedCol(sourcemap.MakeFilePosition(0, math.MaxInt32)) != math.MinInt32 {
		t.Fatal("prefix column addition did not wrap as a Java int")
	}
}
