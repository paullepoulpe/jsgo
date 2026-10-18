package jscomp

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestParallelSourcePreservesTextLicensesAndMappings(t *testing.T) {
	for _, pretty := range []bool{false, true} {
		for _, mode := range []string{"binary", "chunk", "array", "chunkAware"} {
			t.Run(fmt.Sprintf("%s/pretty=%v", mode, pretty), func(t *testing.T) {
				render := func(workers int) (string, string) {
					options := NewCompilerOptions()
					options.SetNumParallelThreads(workers)
					options.SetPrettyPrint(pretty)
					options.SetPrintInputDelimiter(true)
					options.SetEmitUseStrict(true)
					options.SetSourceMapOutputPath(jutil.StringPtr("out.js.map"))
					options.SetPreserveTypeAnnotations(true)
					compiler := NewCompiler()
					compiler.Init(nil, []*SourceFile{
						SourceFile_FromCode("a.js", "/** @license License A */\n/** @type {number} */ var x = 1;"),
						SourceFile_FromCode("b.js", "/** @license License B */\nfunction f(a) { return a + x; }"),
						SourceFile_FromCode("c.js", "/** @license License A */\nf(2);"),
						SourceFile_FromCode("empty.js", "/** @license Empty */"),
					}, options)
					compiler.ParseInputs()
					var source string
					switch mode {
					case "binary":
						source = compiler.ToSource()
					case "chunk":
						source = compiler.ToSourceOfJSChunk(compiler.GetChunks()[0])
					case "array":
						source = strings.Join(compiler.ToSourceArray(NewCompiler_ScriptNodeLicensesOnlyTracker(compiler), compiler.GetChunks()[0]), "\n---\n")
					case "chunkAware":
						tracker := NewCompiler_ChunkGraphAwareLicenseTracker(compiler)
						tracker.SetCurrentChunkContext(compiler.GetChunks()[0])
						source = compiler.ToSourceOfLicenseTrackerAndJSChunk(tracker, compiler.GetChunks()[0])
					}
					var mapping strings.Builder
					compiler.GetSourceMap().AppendTo(&mapping, jutil.StringPtr("out.js"))
					return source, mapping.String()
				}
				wantSource, wantMap := render(1)
				for _, workers := range []int{2, 4} {
					gotSource, gotMap := render(workers)
					if gotSource != wantSource || gotMap != wantMap {
						t.Fatalf("%d workers changed output\nsource:\n%s\nwant:\n%s\nmap:\n%s\nwant:\n%s", workers, gotSource, wantSource, gotMap, wantMap)
					}
				}
			})
		}
	}
}

type customSourceLicenseTracker struct{ nodes []*rhino.Node }

func (c *customSourceLicenseTracker) TrackLicensesForNode(node *rhino.Node) {
	c.nodes = append(c.nodes, node)
}
func (c *customSourceLicenseTracker) EmitLicenses() *jutil.LinkedHashSet[string] {
	return jutil.NewLinkedHashSet[string]()
}

func TestParallelSourceCustomTrackerRemainsSequential(t *testing.T) {
	options := NewCompilerOptions()
	options.SetNumParallelThreads(4)
	compiler := newTestCompiler(options)
	roots := []*rhino.Node{rhino.IR_Script(), rhino.IR_Script()}
	tracker := &customSourceLicenseTracker{}
	var emitted []int
	compiler.toSourceRoots(roots, tracker, func(i int, root *rhino.Node, source *parallelSource) {
		if source != nil {
			t.Fatal("custom tracker incorrectly took the parallel path")
		}
		emitted = append(emitted, i)
	})
	if !reflect.DeepEqual(emitted, []int{0, 1}) {
		t.Fatalf("publication order %v", emitted)
	}
}
