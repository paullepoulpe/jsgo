package jscomp

import (
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"
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

type sourceLicenseTestCompiler struct {
	AbstractCompiler
	scripts map[string]*rhino.Node
}

func (c *sourceLicenseTestCompiler) GetScriptNode(name string) *rhino.Node {
	return c.scripts[name]
}

func TestParallelSourceRecordsFirstEmptySourceName(t *testing.T) {
	compiler := &sourceLicenseTestCompiler{scripts: make(map[string]*rhino.Node)}
	for _, name := range []string{"prior", ""} {
		info := rhino.JSDocInfo_Builder_()
		info.RecordLicense("license:" + name)
		script := rhino.IR_Script()
		script.SetJSDocInfo(info.Build())
		compiler.scripts[name] = script
	}
	for _, scriptOnly := range []bool{false, true} {
		t.Run(fmt.Sprintf("scriptOnly=%v", scriptOnly), func(t *testing.T) {
			makeTracker := func() CodePrinter_LicenseTracker {
				if scriptOnly {
					return NewCompiler_ScriptNodeLicensesOnlyTracker(compiler)
				}
				return newCompiler_SingleBinaryLicenseTracker(compiler)
			}
			makeNode := func(name string) *rhino.Node {
				var node *rhino.Node
				if scriptOnly {
					node = rhino.IR_Script()
				} else {
					node = rhino.IR_Name("value")
				}
				return node.SetStaticSourceFile(rhino.NewSimpleSourceFile(name, rhino.StaticSourceFile_SourceKind_STRONG))
			}
			prior, empty := makeNode("prior"), makeNode("")
			wantTracker, gotTracker := makeTracker(), makeTracker()
			for _, tracker := range []CodePrinter_LicenseTracker{wantTracker, gotTracker} {
				tracker.TrackLicensesForNode(prior)
				tracker.EmitLicenses()
			}
			wantTracker.TrackLicensesForNode(empty)
			wantTracker.TrackLicensesForNode(empty)
			recorder := parallelSourceLicenses{scriptOnly: scriptOnly}
			recorder.TrackLicensesForNode(empty)
			recorder.TrackLicensesForNode(empty)
			if len(recorder.nodes) != 1 {
				t.Fatalf("first empty source transition recorded %d times", len(recorder.nodes))
			}
			for _, node := range recorder.nodes {
				gotTracker.TrackLicensesForNode(node)
			}
			want, got := wantTracker.EmitLicenses().ToSlice(), gotTracker.EmitLicenses().ToSlice()
			if !reflect.DeepEqual(want, []string{"license:"}) || !reflect.DeepEqual(got, want) {
				t.Fatalf("license replay got %v, want %v", got, want)
			}
		})
	}
}

func TestParallelSourceTypedAstRetainsSequentialPath(t *testing.T) {
	options := NewCompilerOptions()
	options.SetNumParallelThreads(4)
	compiler := newTestCompiler(options)
	compiler.typedAstFilesystem = jutil.NewLinkedHashMap[*SourceFile, func() *rhino.Node]()
	roots := []*rhino.Node{rhino.IR_Script(), rhino.IR_Script()}
	var emitted []int
	compiler.toSourceRoots(roots, NewCompiler_ScriptNodeLicensesOnlyTracker(compiler),
		func(i int, root *rhino.Node, source *parallelSource) {
			if source != nil {
				t.Fatal("TypedAST printing incorrectly took the parallel path")
			}
			emitted = append(emitted, i)
		})
	if !reflect.DeepEqual(emitted, []int{0, 1}) {
		t.Fatalf("publication order %v", emitted)
	}
}

type sourceWindowTestFile struct {
	*rhino.SimpleSourceFile
	seen  atomic.Bool
	count *atomic.Int32
}

func (f *sourceWindowTestFile) GetName() string {
	if f.seen.CompareAndSwap(false, true) {
		f.count.Add(1)
	}
	return f.SimpleSourceFile.GetName()
}

func TestParallelSourceBoundsPrivateResultsToWorkerWindow(t *testing.T) {
	for _, workers := range []int{2, 3, 4} {
		t.Run(fmt.Sprintf("workers=%d", workers), func(t *testing.T) {
			options := NewCompilerOptions()
			options.SetNumParallelThreads(workers)
			options.SetEmitUseStrict(true)
			compiler := newTestCompiler(options)
			var count atomic.Int32
			var roots []*rhino.Node
			for i := range 7 {
				file := &sourceWindowTestFile{
					SimpleSourceFile: rhino.NewSimpleSourceFile(fmt.Sprintf("window%d.js", i), rhino.StaticSourceFile_SourceKind_STRONG),
					count:            &count,
				}
				number := rhino.IR_Number(float64(i)).SetStaticSourceFile(file)
				roots = append(roots, rhino.IR_ScriptOfNodes(rhino.IR_ExprResult(number)))
			}
			var emitted []int
			compiler.toSourceRoots(roots, newCompiler_SingleBinaryLicenseTracker(compiler),
				func(i int, root *rhino.Node, source *parallelSource) {
					wantComputed := min(len(roots), (i/workers+1)*workers)
					if computed := int(count.Load()); computed != wantComputed {
						t.Fatalf("at result %d, computed %d roots, want one bounded window (%d)", i, computed, wantComputed)
					}
					if source == nil || strings.Contains(source.code.source, "use strict") != (i == 0) {
						t.Fatalf("strict prefix not limited to the first global input: result %d, %v", i, source)
					}
					emitted = append(emitted, i)
				})
			if !reflect.DeepEqual(emitted, []int{0, 1, 2, 3, 4, 5, 6}) {
				t.Fatalf("publication order %v", emitted)
			}
		})
	}
}
