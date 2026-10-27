package sourcemap

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

// buildRandomMappings adds a tree of nested mappings in pre-order (parents
// before children, siblings left to right, lines growing), as the code
// printer does.
func buildRandomMappings(g *SourceMapGeneratorV3, rng *rand.Rand, n int) {
	type mapping struct {
		src        string
		sym        *string
		orig       FilePosition
		start, end FilePosition
	}
	var out []*mapping
	line, col := 0, 0
	count := 0
	var rec func(depth int)
	rec = func(depth int) {
		for count < n && rng.Intn(4) != 0 {
			count++
			if rng.Intn(5) == 0 {
				line += 1 + rng.Intn(2)
				col = 0
			}
			col += rng.Intn(3)
			m := &mapping{
				src:   "src" + strconv.Itoa(rng.Intn(7)) + ".js",
				orig:  MakeFilePosition(rng.Intn(500), rng.Intn(80)),
				start: MakeFilePosition(line, col),
			}
			if rng.Intn(3) == 0 {
				s := "name" + strconv.Itoa(rng.Intn(20))
				m.sym = &s
			}
			out = append(out, m)
			if depth < 6 {
				rec(depth + 1)
			}
			col += rng.Intn(4)
			if rng.Intn(6) == 0 {
				line++
				col = rng.Intn(3)
			}
			m.end = MakeFilePosition(line, col)
		}
	}
	for count < n {
		rec(0)
		col++
	}
	for _, m := range out {
		g.AddMappingPositions(m.src, m.sym, m.orig, m.start, m.end)
	}
}

func TestParallelMappingsMatchSequential(t *testing.T) {
	saved := parallelChunkMappings
	defer func() { parallelChunkMappings = saved }()
	for seed := int64(1); seed <= 40; seed++ {
		for _, chunk := range []int{1, 2, 3, 7, 64} {
			parallelChunkMappings = chunk
			mk := func(workers int) string {
				g := NewSourceMapGeneratorV3()
				if seed%3 == 0 {
					g.SetWrapperPrefix("(function(){\n  ")
				}
				buildRandomMappings(g, rand.New(rand.NewSource(seed)), 400)
				g.SetParallelWorkers(workers)
				var sb strings.Builder
				name := "out.js"
				g.AppendTo(&sb, &name)
				return sb.String()
			}
			seq, par := mk(1), mk(4)
			if seq != par {
				t.Fatalf("seed %d chunk %d: parallel output differs\nseq: %.300s\npar: %.300s", seed, chunk, seq, par)
			}
		}
	}
}
