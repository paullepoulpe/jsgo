// Package flatspike benchmarks flat slice node storage against the chunk
// store on real parsed files (SPIKE; see
// docs/memory-research/round3-flat-slices.md).
//
//	FLATSPIKE_ARGS=<ws>/args.txt FLATSPIKE_WS=<ws> go test -bench . -benchtime 10x ./internal/rhino/flatspike
package flatspike

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type nopReporter struct{}

func (nopReporter) Warning(message, sourceName string, line, lineOffset int) {}
func (nopReporter) Error(message, sourceName string, line, lineOffset int)   {}

type corpus struct {
	roots    []*rhino.Node // chunk store, one SCRIPT per file
	preRoots []*rhino.Node // a second parse, compacted (pre-order chunks)
	preNodes []*rhino.Node
	allNodes []*rhino.Node // every node, allocation order
	perFile  []*rhino.FlatStore
	fileRoot []uint32
	perPre   []*rhino.FlatStore // per-file, pre-order
	preRoot  []uint32
	unified  [2]*rhino.FlatStore // [alloc order, pre order]
	uniRoots [2][]uint32
	two      *rhino.FileStores
	nodes    int
}

var (
	once sync.Once
	c    corpus
)

func load(b *testing.B) *corpus {
	once.Do(func() {
		argsFile, ws := os.Getenv("FLATSPIKE_ARGS"), os.Getenv("FLATSPIKE_WS")
		if argsFile == "" {
			return
		}
		data, err := os.ReadFile(argsFile)
		if err != nil {
			panic(err)
		}
		lines := strings.Split(string(data), "\n")
		var files []string
		for i := 0; i < len(lines)-1; i++ {
			if lines[i] == "--js" {
				files = append(files, lines[i+1])
			}
		}
		config := parsing.ParserRunner_CreateConfigWithModes(parsing.Config_LanguageMode_ES_NEXT,
			parsing.Config_JsDocParsing_TYPES_ONLY, parsing.Config_RunMode_KEEP_GOING, nil, true,
			parsing.Config_StrictMode_SLOPPY)
		rhino.NewNodeGeneration()
		parseAll := func() (roots []*rhino.Node) {
			for _, f := range files {
				src, err := os.ReadFile(filepath.Join(ws, f))
				if err != nil {
					panic(err)
				}
				r := parsing.ParserRunner_Parse(rhino.NewSimpleSourceFile(f, rhino.StaticSourceFile_SourceKind_STRONG),
					string(src), config, nopReporter{})
				if r.Ast != nil {
					roots = append(roots, r.Ast)
				}
			}
			return roots
		}
		c.roots = parseAll()
		second := parseAll()
		k := rhino.NewCompactor()
		for _, r := range second {
			nr := k.CopyTree(r)
			c.preRoots = append(c.preRoots, nr)
			c.preNodes = append(c.preNodes, rhino.CollectNodes(nr)...)
		}
		for _, r := range c.roots {
			c.allNodes = append(c.allNodes, rhino.CollectNodes(r)...)
			s := &rhino.FlatStore{}
			c.fileRoot = append(c.fileRoot, rhino.FlattenInto(s, r, rhino.FlatAllocOrder))
			c.perFile = append(c.perFile, s)
			p := &rhino.FlatStore{}
			c.preRoot = append(c.preRoot, rhino.FlattenInto(p, r, rhino.FlatPreOrder))
			c.perPre = append(c.perPre, p)
		}
		c.nodes = len(c.allNodes)
		for o := 0; o < 2; o++ {
			u := &rhino.FlatStore{Nodes: make([]rhino.Node, 0, c.nodes+1)}
			for _, r := range c.roots {
				c.uniRoots[o] = append(c.uniRoots[o], rhino.FlattenInto(u, r, rhino.FlatOrder(o)))
			}
			c.unified[o] = u
		}
		c.two = &rhino.FileStores{}
		for _, s := range c.perFile {
			c.two.Files = append(c.two.Files, s.Nodes)
		}
		runtime.GC()
	})
	if c.roots == nil {
		b.Skip("FLATSPIKE_ARGS not set")
	}
	return &c
}

var sink uint64

func report(b *testing.B, c *corpus) {
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N)/float64(c.nodes), "ns/node")
}

func BenchmarkWalk(b *testing.B) {
	c := load(b)
	b.Run("chunk-api", func(b *testing.B) {
		for range b.N {
			for _, r := range c.roots {
				sink += rhino.ChunkWalk(r)
			}
		}
		report(b, c)
	})
	b.Run("chunk-inlined", func(b *testing.B) {
		for range b.N {
			for _, r := range c.roots {
				sink += rhino.ChunkWalkInl(r)
			}
		}
		report(b, c)
	})
	b.Run("chunk-api-compacted", func(b *testing.B) {
		for range b.N {
			for _, r := range c.preRoots {
				sink += rhino.ChunkWalk(r)
			}
		}
		report(b, c)
	})
	b.Run("chunk-inlined-compacted", func(b *testing.B) {
		for range b.N {
			for _, r := range c.preRoots {
				sink += rhino.ChunkWalkInl(r)
			}
		}
		report(b, c)
	})
	b.Run("per-file-slice", func(b *testing.B) {
		for range b.N {
			for i, s := range c.perFile {
				sink += s.FlatWalk(c.fileRoot[i])
			}
		}
		report(b, c)
	})
	b.Run("per-file-slice-preorder", func(b *testing.B) {
		for range b.N {
			for i, s := range c.perPre {
				sink += s.FlatWalk(c.preRoot[i])
			}
		}
		report(b, c)
	})
	b.Run("per-file-two-level-id", func(b *testing.B) {
		for range b.N {
			for i := range c.perFile {
				sink += c.two.TwoLevelWalk(uint32(i)<<rhino.FileShift | c.fileRoot[i])
			}
		}
		report(b, c)
	})
	b.Run("unified-slice-alloc-order", func(b *testing.B) {
		u := c.unified[0]
		for range b.N {
			for _, r := range c.uniRoots[0] {
				sink += u.FlatWalk(r)
			}
		}
		report(b, c)
	})
	b.Run("unified-slice-preorder", func(b *testing.B) {
		u := c.unified[1]
		for range b.N {
			for _, r := range c.uniRoots[1] {
				sink += u.FlatWalk(r)
			}
		}
		report(b, c)
	})
}

func BenchmarkParentChains(b *testing.B) {
	c := load(b)
	b.Run("chunk-api", func(b *testing.B) {
		for range b.N {
			sink += rhino.ChunkDepthSum(c.allNodes)
		}
		report(b, c)
	})
	b.Run("chunk-inlined", func(b *testing.B) {
		for range b.N {
			sink += rhino.ChunkDepthSumInl(c.allNodes)
		}
		report(b, c)
	})
	b.Run("per-file-slice", func(b *testing.B) {
		for range b.N {
			for _, s := range c.perFile {
				sink += s.FlatDepthSum(1, uint32(len(s.Nodes)))
			}
		}
		report(b, c)
	})
	b.Run("unified-slice-alloc-order", func(b *testing.B) {
		u := c.unified[0]
		for range b.N {
			sink += u.FlatDepthSum(1, uint32(len(u.Nodes)))
		}
		report(b, c)
	})
	b.Run("chunk-api-compacted", func(b *testing.B) {
		for range b.N {
			sink += rhino.ChunkDepthSum(c.preNodes)
		}
		report(b, c)
	})
	b.Run("chunk-inlined-compacted", func(b *testing.B) {
		for range b.N {
			sink += rhino.ChunkDepthSumInl(c.preNodes)
		}
		report(b, c)
	})
	b.Run("unified-slice-preorder", func(b *testing.B) {
		u := c.unified[1]
		for range b.N {
			sink += u.FlatDepthSum(1, uint32(len(u.Nodes)))
		}
		report(b, c)
	})
}

func BenchmarkMergeToUnified(b *testing.B) {
	c := load(b)
	for range b.N {
		u, _ := rhino.MergeStores(c.perPre, c.preRoot)
		sink += uint64(len(u.Nodes))
	}
	report(b, c)
}

func TestSameResults(t *testing.T) {
	if os.Getenv("FLATSPIKE_ARGS") == "" {
		t.Skip("FLATSPIKE_ARGS not set")
	}
	c := load(nil)
	var a, f, u, two, inl uint64
	for i, r := range c.roots {
		a += rhino.ChunkWalk(r)
		inl += rhino.ChunkWalkInl(r)
		f += c.perFile[i].FlatWalk(c.fileRoot[i])
		u += c.unified[1].FlatWalk(c.uniRoots[1][i])
		two += c.two.TwoLevelWalk(uint32(i)<<rhino.FileShift | c.fileRoot[i])
	}
	m, mroots := rhino.MergeStores(c.perPre, c.preRoot)
	var merged uint64
	for _, r := range mroots {
		merged += m.FlatWalk(r)
	}
	var pre uint64
	for _, r := range c.preRoots {
		pre += rhino.ChunkWalk(r)
	}
	if a != f || a != u || a != two || a != inl || a != pre || a != merged {
		t.Fatalf("walk sums differ: chunk %d inl %d perfile %d unified %d two %d", a, inl, f, u, two)
	}
	var pd uint64
	for _, s := range c.perFile {
		pd += s.FlatDepthSum(1, uint32(len(s.Nodes)))
	}
	if cd := rhino.ChunkDepthSum(c.allNodes); cd != pd {
		t.Fatalf("depth sums differ: chunk %d flat %d", cd, pd)
	}
	t.Logf("%d files, %d nodes", len(c.roots), c.nodes)
}
