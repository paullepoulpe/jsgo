// Port of test/com/google/javascript/jscomp/JSChunkGraphTest.java @ closure-compiler v20260415

package jscomp

import (
	"math/rand"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type jsChunkGraphTest struct {
	t      *testing.T
	chunkA *JSChunk
	chunkB *JSChunk
	chunkC *JSChunk
	chunkD *JSChunk
	chunkE *JSChunk
	chunkF *JSChunk
	graph  *JSChunkGraph

	// For resolving dependencies only.
	compiler AbstractCompiler
}

func newJSChunkGraphTest(t *testing.T) *jsChunkGraphTest {
	return &jsChunkGraphTest{t: t, compiler: newTestCompilerUninitialized()}
}

func (x *jsChunkGraphTest) makeDeps() {
	x.chunkA = NewJSChunk("chunkA")
	x.chunkB = NewJSChunk("chunkB")
	x.chunkC = NewJSChunk("chunkC")
	x.chunkD = NewJSChunk("chunkD")
	x.chunkE = NewJSChunk("chunkE")
	x.chunkF = NewJSChunk("chunkF")
	x.chunkB.AddDependency(x.chunkA) //     __A__
	x.chunkC.AddDependency(x.chunkA) //    /  |  \
	x.chunkD.AddDependency(x.chunkB) //   B   C  |
	x.chunkE.AddDependency(x.chunkB) //  / \ /|  |
	x.chunkE.AddDependency(x.chunkC) // D   E | /
	x.chunkF.AddDependency(x.chunkA) //      \|/
	x.chunkF.AddDependency(x.chunkC) //       F
	x.chunkF.AddDependency(x.chunkE)
}

func (x *jsChunkGraphTest) makeGraph() {
	x.graph = NewJSChunkGraphFromArray([]*JSChunk{x.chunkA, x.chunkB, x.chunkC, x.chunkD, x.chunkE, x.chunkF})
}

func (x *jsChunkGraphTest) getWeakChunk() *JSChunk {
	return x.graph.GetChunkByName(JSChunk_WEAK_CHUNK_NAME)
}

func (x *jsChunkGraphTest) allChunks() []*JSChunk {
	return []*JSChunk{x.chunkA, x.chunkB, x.chunkC, x.chunkD, x.chunkE, x.chunkF}
}

// jsChunkGraphTest_containsExactly is containsExactly(...) (any order).
func jsChunkGraphTest_containsExactly[T comparable](t *testing.T, what string, got []T, want ...T) {
	t.Helper()
	ok := len(got) == len(want)
	if ok {
		rest := slices.Clone(want)
		for _, g := range got {
			i := slices.Index(rest, g)
			if i < 0 {
				ok = false
				break
			}
			rest = slices.Delete(rest, i, i+1)
		}
	}
	if !ok {
		t.Errorf("%s: got %v, want (any order) %v", what, got, want)
	}
}

func TestJSChunkGraph_MakesWeakChunkIfNotPassed(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	assertEqual(t, "chunk count", x.graph.GetChunkCount(), 7)
	if !x.graph.GetChunksByName().ContainsKey(JSChunk_WEAK_CHUNK_NAME) {
		t.Errorf("no weak chunk")
	}
	jsChunkGraphTest_containsExactly(t, "weak deps", x.getWeakChunk().GetAllDependencies().ToSlice(), x.allChunks()...)
}

func (x *jsChunkGraphTest) newWeakChunk(withF bool) *JSChunk {
	weakChunk := NewJSChunk(JSChunk_WEAK_CHUNK_NAME)
	weakChunk.AddDependency(x.chunkA)
	weakChunk.AddDependency(x.chunkB)
	weakChunk.AddDependency(x.chunkC)
	weakChunk.AddDependency(x.chunkD)
	weakChunk.AddDependency(x.chunkE)
	if withF {
		weakChunk.AddDependency(x.chunkF)
	}
	return weakChunk
}

func TestJSChunkGraph_AcceptExistingWeakChunk(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	weakChunk := x.newWeakChunk(true)

	weakChunk.AddSourceFile(SourceFile_FromCodeWithKind("weak", "", rhino.StaticSourceFile_SourceKind_WEAK))

	graph := NewJSChunkGraphFromArray(append(x.allChunks(), weakChunk))

	assertEqual(t, "chunk count", graph.GetChunkCount(), 7)
	assertEqual(t, "weak chunk", graph.GetChunkByName(JSChunk_WEAK_CHUNK_NAME), weakChunk)
}

func TestJSChunkGraph_ExistingWeakChunkMustHaveDependenciesOnAllOtherChunks(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	weakChunk := x.newWeakChunk(false)
	// Missing F

	e := assertThrows[jutil.IIllegalStateException](t, func() {
		NewJSChunkGraphFromArray(append(x.allChunks(), weakChunk))
	})
	assertEqual(t, "message", e.GetMessage(),
		"A weak chunk already exists but it does not depend on every other chunk.")
}

func TestJSChunkGraph_WeakFileCannotExistOutsideWeakChunk(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	weakChunk := x.newWeakChunk(true)

	x.chunkA.AddSourceFile(SourceFile_FromCodeWithKind("a", "", rhino.StaticSourceFile_SourceKind_WEAK))

	e := assertThrows[jutil.IIllegalStateException](t, func() {
		NewJSChunkGraphFromArray(append(x.allChunks(), weakChunk))
	})

	if !strings.Contains(e.GetMessage(), "Found these weak sources in other chunks:\n  a (in chunk chunkA)") {
		t.Errorf("message %q", e.GetMessage())
	}
}

func TestJSChunkGraph_StrongFileCannotExistInWeakChunk(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	weakChunk := x.newWeakChunk(true)

	weakChunk.AddSourceFile(SourceFile_FromCodeWithKind("a", "", rhino.StaticSourceFile_SourceKind_STRONG))

	e := assertThrows[jutil.IIllegalStateException](t, func() {
		NewJSChunkGraphFromArray(append(x.allChunks(), weakChunk))
	})
	if !strings.Contains(e.GetMessage(), "Found these strong sources in the weak chunk:\n  a") {
		t.Errorf("message %q", e.GetMessage())
	}
}

func TestJSChunkGraph_SmallerTreeBeatsDeeperTree(t *testing.T) {
	x := newJSChunkGraphTest(t)
	a := NewJSChunk("a")
	b := NewJSChunk("b")
	c := NewJSChunk("c")
	d := NewJSChunk("d")
	e := NewJSChunk("e")
	f := NewJSChunk("f")
	g := NewJSChunk("g")
	h := NewJSChunk("h")
	//   a
	//  / \
	// b   c
	b.AddDependency(a)
	c.AddDependency(a)
	//   b
	//  /|\
	// e f g
	e.AddDependency(b)
	f.AddDependency(b)
	g.AddDependency(b)
	//     c
	//     |
	//     d
	//   // \\
	//  / | | \
	// e  f g  h
	d.AddDependency(c)
	e.AddDependency(d)
	f.AddDependency(d)
	g.AddDependency(d)
	h.AddDependency(d)
	graph := NewJSChunkGraphFromArray([]*JSChunk{a, b, c, d, e, f, g, h})
	// d is deeper, but it also has an extra dependent node, so b is the better choice.
	x.assertSmallestCoveringSubtreeInGraph(b, graph, a, e, f, g)
	// However, if the parent tree we're looking at is c, then b isn't an option
	x.assertSmallestCoveringSubtreeInGraph(d, graph, c, e, f, g)
}

func TestJSChunkGraph_ChunkDepth(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	assertEqual(t, "chunkA should have depth 0", x.chunkA.GetDepth(), 0)
	assertEqual(t, "chunkB should have depth 1", x.chunkB.GetDepth(), 1)
	assertEqual(t, "chunkC should have depth 1", x.chunkC.GetDepth(), 1)
	assertEqual(t, "chunkD should have depth 2", x.chunkD.GetDepth(), 2)
	assertEqual(t, "chunkE should have depth 2", x.chunkE.GetDepth(), 2)
	assertEqual(t, "chunkF should have depth 3", x.chunkF.GetDepth(), 3)
}

func TestJSChunkGraph_DeepestCommonDep(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	A, B, C, D, E, F := x.chunkA, x.chunkB, x.chunkC, x.chunkD, x.chunkE, x.chunkF
	x.assertDeepestCommonDep(nil, A, A)
	x.assertDeepestCommonDep(nil, A, B)
	x.assertDeepestCommonDep(nil, A, C)
	x.assertDeepestCommonDep(nil, A, D)
	x.assertDeepestCommonDep(nil, A, E)
	x.assertDeepestCommonDep(nil, A, F)
	x.assertDeepestCommonDep(A, B, B)
	x.assertDeepestCommonDep(A, B, C)
	x.assertDeepestCommonDep(A, B, D)
	x.assertDeepestCommonDep(A, B, E)
	x.assertDeepestCommonDep(A, B, F)
	x.assertDeepestCommonDep(A, C, C)
	x.assertDeepestCommonDep(A, C, D)
	x.assertDeepestCommonDep(A, C, E)
	x.assertDeepestCommonDep(A, C, F)
	x.assertDeepestCommonDep(B, D, D)
	x.assertDeepestCommonDep(B, D, E)
	x.assertDeepestCommonDep(B, D, F)
	x.assertDeepestCommonDep(C, E, E)
	x.assertDeepestCommonDep(C, E, F)
	x.assertDeepestCommonDep(E, F, F)
}

func TestJSChunkGraph_DeepestCommonDepInclusive(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	A, B, C, D, E, F := x.chunkA, x.chunkB, x.chunkC, x.chunkD, x.chunkE, x.chunkF
	x.assertDeepestCommonDepInclusive(A, A, A)
	x.assertDeepestCommonDepInclusive(A, A, B)
	x.assertDeepestCommonDepInclusive(A, A, C)
	x.assertDeepestCommonDepInclusive(A, A, D)
	x.assertDeepestCommonDepInclusive(A, A, E)
	x.assertDeepestCommonDepInclusive(A, A, F)
	x.assertDeepestCommonDepInclusive(B, B, B)
	x.assertDeepestCommonDepInclusive(A, B, C)
	x.assertDeepestCommonDepInclusive(B, B, D)
	x.assertDeepestCommonDepInclusive(B, B, E)
	x.assertDeepestCommonDepInclusive(B, B, F)
	x.assertDeepestCommonDepInclusive(C, C, C)
	x.assertDeepestCommonDepInclusive(A, C, D)
	x.assertDeepestCommonDepInclusive(C, C, E)
	x.assertDeepestCommonDepInclusive(C, C, F)
	x.assertDeepestCommonDepInclusive(D, D, D)
	x.assertDeepestCommonDepInclusive(B, D, E)
	x.assertDeepestCommonDepInclusive(B, D, F)
	x.assertDeepestCommonDepInclusive(E, E, E)
	x.assertDeepestCommonDepInclusive(E, E, F)
	x.assertDeepestCommonDepInclusive(F, F, F)
}

func TestJSChunkGraph_SmallestCoveringSubtree(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	A, B, C, D, E, F := x.chunkA, x.chunkB, x.chunkC, x.chunkD, x.chunkE, x.chunkF
	x.assertSmallestCoveringSubtree(A, A, A, A)
	x.assertSmallestCoveringSubtree(A, A, A, B)
	x.assertSmallestCoveringSubtree(A, A, A, C)
	x.assertSmallestCoveringSubtree(A, A, A, D)
	x.assertSmallestCoveringSubtree(A, A, A, E)
	x.assertSmallestCoveringSubtree(A, A, A, F)
	x.assertSmallestCoveringSubtree(B, A, B, B)
	x.assertSmallestCoveringSubtree(A, A, B, C)
	x.assertSmallestCoveringSubtree(B, A, B, D)
	x.assertSmallestCoveringSubtree(B, A, B, E)
	x.assertSmallestCoveringSubtree(B, A, B, F)
	x.assertSmallestCoveringSubtree(C, A, C, C)
	x.assertSmallestCoveringSubtree(A, A, C, D)
	x.assertSmallestCoveringSubtree(C, A, C, E)
	x.assertSmallestCoveringSubtree(C, A, C, F)
	x.assertSmallestCoveringSubtree(D, A, D, D)
	x.assertSmallestCoveringSubtree(B, A, D, E)
	x.assertSmallestCoveringSubtree(B, A, D, F)
	x.assertSmallestCoveringSubtree(E, A, E, E)
	x.assertSmallestCoveringSubtree(E, A, E, F)
	x.assertSmallestCoveringSubtree(F, A, F, F)
}

func TestJSChunkGraph_GetTransitiveDepsDeepestFirst(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	A, B, C, D, E, F := x.chunkA, x.chunkB, x.chunkC, x.chunkD, x.chunkE, x.chunkF
	x.assertTransitiveDepsDeepestFirst(A)
	x.assertTransitiveDepsDeepestFirst(B, A)
	x.assertTransitiveDepsDeepestFirst(C, A)
	x.assertTransitiveDepsDeepestFirst(D, B, A)
	x.assertTransitiveDepsDeepestFirst(E, C, B, A)
	x.assertTransitiveDepsDeepestFirst(F, E, C, B, A)
}

func TestJSChunkGraph_ManageDependenciesLooseWithoutEntryPoint(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	x.setUpManageDependenciesTest()
	depOptions := DependencyOptions_PruneLegacyForEntryPoints([]*ModuleIdentifier{})
	results := x.graph.ManageDependencies(x.compiler, depOptions)

	x.assertInputs(x.chunkA, "a1", "a3")
	x.assertInputs(x.chunkB, "a2", "b2")
	x.assertInputs(x.chunkC) // no inputs
	x.assertInputs(x.chunkE, "c1", "e1", "e2")

	assertSliceEqual(t, "results", jsChunkGraphTest_sourceNames(results),
		[]string{"a1", "a3", "a2", "b2", "c1", "e1", "e2"})
}

func TestJSChunkGraph_ManageDependenciesLooseWithEntryPoint(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	x.setUpManageDependenciesTest()
	depOptions := DependencyOptions_PruneLegacyForEntryPoints(
		[]*ModuleIdentifier{ModuleIdentifier_ForClosure("c2")})
	results := x.graph.ManageDependencies(x.compiler, depOptions)

	x.assertInputs(x.chunkA, "a1", "a3")
	x.assertInputs(x.chunkB, "a2", "b2")
	x.assertInputs(x.chunkC, "c1", "c2")
	x.assertInputs(x.chunkE, "e1", "e2")

	assertSliceEqual(t, "results", jsChunkGraphTest_sourceNames(results),
		[]string{"a1", "a3", "a2", "b2", "c1", "c2", "e1", "e2"})
}

func TestJSChunkGraph_ManageDependenciesStrictWithEntryPoint(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	x.setUpManageDependenciesTest()
	depOptions := DependencyOptions_PruneForEntryPoints([]*ModuleIdentifier{ModuleIdentifier_ForClosure("c2")})
	results := x.graph.ManageDependencies(x.compiler, depOptions)

	// Everything gets pushed up into chunk c, because that's
	// the only one that has entry points.
	x.assertInputs(x.chunkA)
	x.assertInputs(x.chunkB)
	x.assertInputs(x.chunkC, "a1", "c1", "c2")
	x.assertInputs(x.chunkE)

	assertSliceEqual(t, "results", jsChunkGraphTest_sourceNames(results), []string{"a1", "c1", "c2"})
}

func TestJSChunkGraph_ManageDependenciesStrictForGoogRequireDynamic(t *testing.T) {
	x := newJSChunkGraphTest(t)
	chunkA := NewJSChunk("chunk")
	x.graph = NewJSChunkGraphFromArray([]*JSChunk{chunkA})
	inputs := []*CompilerInput{}
	compilerInputA1 := NewCompilerInput(jsChunkGraphTest_code("a1", jsChunkGraphTest_list("a1"), jsChunkGraphTest_list()))
	compilerInputA1.AddRequireDynamicImports("a2")
	chunkA.Add(compilerInputA1)
	chunkA.AddSourceFile(jsChunkGraphTest_code("a2", jsChunkGraphTest_list("a2"), jsChunkGraphTest_list()))
	inputs = append(inputs, chunkA.GetInputs()...)
	for _, input := range inputs {
		input.SetCompiler(x.compiler)
	}
	depOptions := DependencyOptions_PruneForEntryPoints([]*ModuleIdentifier{ModuleIdentifier_ForClosure("a1")})
	results := x.graph.ManageDependencies(x.compiler, depOptions)

	x.assertInputs(chunkA, "a1", "a2")
	jsChunkGraphTest_containsExactly(t, "results", jsChunkGraphTest_sourceNames(results), "a1", "a2")
}

func TestJSChunkGraph_ManageDependenciesStrictWithEntryPointWithDuplicates(t *testing.T) {
	x := newJSChunkGraphTest(t)
	a := NewJSChunk("a")
	graph := NewJSChunkGraphFromArray([]*JSChunk{a})

	// Create all the input files.
	inputs := []*CompilerInput{}
	a.AddSourceFile(jsChunkGraphTest_code("a1", jsChunkGraphTest_list("a1"), jsChunkGraphTest_list("a2")))
	a.AddSourceFile(jsChunkGraphTest_code("a2", jsChunkGraphTest_list("a2"), jsChunkGraphTest_list()))
	a.AddSourceFile(jsChunkGraphTest_code("a3", jsChunkGraphTest_list("a2"), jsChunkGraphTest_list()))
	inputs = append(inputs, a.GetInputs()...)
	for _, input := range inputs {
		input.SetCompiler(x.compiler)
	}

	depOptions := DependencyOptions_PruneForEntryPoints([]*ModuleIdentifier{ModuleIdentifier_ForClosure("a1")})
	results := graph.ManageDependencies(x.compiler, depOptions)

	// Everything gets pushed up into chunk c, because that's
	// the only one that has entry points.
	x.assertInputs(a, "a2", "a3", "a1")

	assertSliceEqual(t, "results", jsChunkGraphTest_sourceNames(results), []string{"a2", "a3", "a1"})
}

func TestJSChunkGraph_ManageDependenciesSortOnly(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	x.setUpManageDependenciesTest()
	results := x.graph.ManageDependencies(x.compiler, DependencyOptions_SortOnly())

	x.assertInputs(x.chunkA, "a1", "a2", "a3")
	x.assertInputs(x.chunkB, "b1", "b2")
	x.assertInputs(x.chunkC, "c1", "c2")
	x.assertInputs(x.chunkE, "e1", "e2")

	assertSliceEqual(t, "results", jsChunkGraphTest_sourceNames(results),
		[]string{"a1", "a2", "a3", "b1", "b2", "c1", "c2", "e1", "e2"})
}

// NOTE: The newline between the @provideGoog comment and the var statement is required.
const jsChunkGraphTest_BASEJS = "/** @fileoverview\n" +
	" * @provideGoog */\n" +
	"var COMPILED = false;\n" +
	"var goog = goog || {}\n"

func TestJSChunkGraph_ManageDependenciesSortOnlyImpl(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	x.chunkA.AddSourceFile(jsChunkGraphTest_code("a2", jsChunkGraphTest_list("a2"), jsChunkGraphTest_list("a1")))
	x.chunkA.AddSourceFile(jsChunkGraphTest_code("a1", jsChunkGraphTest_list("a1"), jsChunkGraphTest_list()))
	x.chunkA.AddSourceFile(jsChunkGraphTest_codeWithSource("base.js", jsChunkGraphTest_BASEJS, jsChunkGraphTest_list(), jsChunkGraphTest_list()))

	for _, input := range x.chunkA.GetInputs() {
		input.SetCompiler(x.compiler)
	}

	results := x.graph.ManageDependencies(x.compiler, DependencyOptions_SortOnly())

	x.assertInputs(x.chunkA, "base.js", "a1", "a2")

	assertSliceEqual(t, "results", jsChunkGraphTest_sourceNames(results), []string{"base.js", "a1", "a2"})
}

func TestJSChunkGraph_NoFiles(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	results := x.graph.ManageDependencies(x.compiler, DependencyOptions_SortOnly())
	assertEqual(t, "results", len(results), 0)
}

func TestJSChunkGraph_ToJson(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	chunks := x.graph.ToJson()
	assertEqual(t, "size", chunks.Size(), 7)
	for i := 0; i < chunks.Size(); i++ {
		m := chunks.Get(i).GetAsJsonObject()
		for _, key := range []string{"name", "dependencies", "transitive-dependencies", "inputs"} {
			if m.Get(key) == nil {
				t.Errorf("chunk %d: no %s", i, key)
			}
		}
	}
	m := chunks.Get(3).GetAsJsonObject()
	assertEqual(t, "name", m.Get("name").GetAsString(), "chunkD")
	assertEqual(t, "dependencies", m.Get("dependencies").GetAsJsonArray().String(), "[\"chunkB\"]")
	assertEqual(t, "transitive-dependencies", m.Get("transitive-dependencies").GetAsJsonArray().Size(), 2)
	assertEqual(t, "inputs", m.Get("inputs").GetAsJsonArray().String(), "[]")
}

func (x *jsChunkGraphTest) setUpManageDependenciesTest() []*CompilerInput {
	inputs := []*CompilerInput{}

	x.chunkA.AddSourceFile(jsChunkGraphTest_code("a1", jsChunkGraphTest_list("a1"), jsChunkGraphTest_list()))
	x.chunkA.AddSourceFile(jsChunkGraphTest_code("a2", jsChunkGraphTest_list("a2"), jsChunkGraphTest_list("a1")))
	x.chunkA.AddSourceFile(jsChunkGraphTest_code("a3", jsChunkGraphTest_list(), jsChunkGraphTest_list("a1")))

	x.chunkB.AddSourceFile(jsChunkGraphTest_code("b1", jsChunkGraphTest_list("b1"), jsChunkGraphTest_list("a2")))
	x.chunkB.AddSourceFile(jsChunkGraphTest_code("b2", jsChunkGraphTest_list(), jsChunkGraphTest_list("a1", "a2")))

	x.chunkC.AddSourceFile(jsChunkGraphTest_code("c1", jsChunkGraphTest_list("c1"), jsChunkGraphTest_list("a1")))
	x.chunkC.AddSourceFile(jsChunkGraphTest_code("c2", jsChunkGraphTest_list("c2"), jsChunkGraphTest_list("c1")))

	x.chunkE.AddSourceFile(jsChunkGraphTest_code("e1", jsChunkGraphTest_list(), jsChunkGraphTest_list("c1")))
	x.chunkE.AddSourceFile(jsChunkGraphTest_code("e2", jsChunkGraphTest_list(), jsChunkGraphTest_list("c1")))

	inputs = append(inputs, x.chunkA.GetInputs()...)
	inputs = append(inputs, x.chunkB.GetInputs()...)
	inputs = append(inputs, x.chunkC.GetInputs()...)
	inputs = append(inputs, x.chunkE.GetInputs()...)

	for _, input := range inputs {
		input.SetCompiler(x.compiler)
	}
	return inputs
}

func jsChunkGraphTest_shuffle(rnd *rand.Rand, files []*SourceFile) {
	rnd.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
}

func TestJSChunkGraph_GoogBaseOrderedCorrectly(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	l := jsChunkGraphTest_list
	sourceFiles := []*SourceFile{
		jsChunkGraphTest_code("a9", l("a9"), l()),
		jsChunkGraphTest_code("a8", l("a8"), l()),
		jsChunkGraphTest_code("a7", l("a7"), l()),
		jsChunkGraphTest_code("a6", l("a6"), l()),
		jsChunkGraphTest_code("a5", l("a5"), l()),
		jsChunkGraphTest_code("a4", l("a4"), l()),
		jsChunkGraphTest_code("a3", l("a3"), l()),
		jsChunkGraphTest_code("a2", l("a2"), l()),
		jsChunkGraphTest_code("a1", l("a1"), l("a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9")),
		jsChunkGraphTest_codeWithSource("base.js", jsChunkGraphTest_BASEJS, l(), l()),
	}

	depOptions := DependencyOptions_PruneForEntryPoints([]*ModuleIdentifier{ModuleIdentifier_ForClosure("a1")})
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 10; i++ {
		jsChunkGraphTest_shuffle(rnd, sourceFiles)
		x.chunkA.RemoveAll()
		for _, sourceFile := range sourceFiles {
			x.chunkA.AddSourceFile(sourceFile)
		}

		for _, input := range x.chunkA.GetInputs() {
			input.SetCompiler(x.compiler)
		}

		results := x.graph.ManageDependencies(x.compiler, depOptions)

		want := []string{"base.js", "a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9", "a1"}
		x.assertInputs(x.chunkA, want...)
		assertSliceEqual(t, "results", jsChunkGraphTest_sourceNames(results), want)
	}
}

func TestJSChunkGraph_ProperEs6ModuleOrdering(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	x.makeGraph()
	l := jsChunkGraphTest_list
	sourceFiles := []*SourceFile{
		jsChunkGraphTest_code("/entry.js", l(), l()),
		jsChunkGraphTest_code("/a/a.js", l(), l()),
		jsChunkGraphTest_code("/a/b.js", l(), l()),
		jsChunkGraphTest_code("/b/a.js", l(), l()),
		jsChunkGraphTest_code("/b/b.js", l(), l()),
		jsChunkGraphTest_code("/b/c.js", l(), l()),
		jsChunkGraphTest_code("/important.js", l(), l()),
	}

	forFile := func(s string) string { return ModuleIdentifier_ForFile(s).String() }
	orderedRequires := map[string][]string{
		"/entry.js": {
			forFile("/b/b.js"),
			forFile("/b/a.js"),
			forFile("/important.js"),
			forFile("/a/b.js"),
			forFile("/a/a.js"),
		},
		"/a/a.js":       {},
		"/a/b.js":       {},
		"/b/a.js":       {},
		"/b/b.js":       {forFile("/b/c.js")},
		"/b/c.js":       {},
		"/important.js": {},
	}

	depOptions := DependencyOptions_PruneForEntryPoints([]*ModuleIdentifier{ModuleIdentifier_ForFile("/entry.js")})
	rnd := rand.New(rand.NewSource(2))
	for iterationCount := 0; iterationCount < 10; iterationCount++ {
		jsChunkGraphTest_shuffle(rnd, sourceFiles)
		x.chunkA.RemoveAll()
		for _, sourceFile := range sourceFiles {
			x.chunkA.AddSourceFile(sourceFile)
		}

		for _, input := range x.chunkA.GetInputs() {
			input.SetCompiler(x.compiler)
			for _, require := range orderedRequires[input.GetSourceFile().GetName()] {
				input.AddOrderedRequire(DependencyInfo_Require_CompilerChunk(require))
			}
			input.SetHasFullParseDependencyInfo(true)
		}

		results := x.graph.ManageDependencies(x.compiler, depOptions)

		want := []string{"/b/c.js", "/b/b.js", "/b/a.js", "/important.js", "/a/b.js", "/a/a.js", "/entry.js"}
		x.assertInputs(x.chunkA, want...)
		assertSliceEqual(t, "results", jsChunkGraphTest_sourceNames(results), want)
	}
}

func jsChunkGraphTest_sourceFiles(inputs []*CompilerInput) []*SourceFile {
	out := []*SourceFile{}
	for _, input := range inputs {
		out = append(out, input.GetSourceFile())
	}
	return out
}

func jsChunkGraphTest_weakStrong() (weak1, weak2, strong1, strong2 *SourceFile) {
	weak1 = SourceFile_FromCodeWithKind("weak1", "", rhino.StaticSourceFile_SourceKind_WEAK)
	weak2 = SourceFile_FromCodeWithKind("weak2", "", rhino.StaticSourceFile_SourceKind_WEAK)
	strong1 = SourceFile_FromCodeWithKind("strong1", "", rhino.StaticSourceFile_SourceKind_STRONG)
	strong2 = SourceFile_FromCodeWithKind("strong2", "", rhino.StaticSourceFile_SourceKind_STRONG)
	return
}

func TestJSChunkGraph_MoveMarkedWeakSources(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()

	weak1, weak2, strong1, strong2 := jsChunkGraphTest_weakStrong()

	x.chunkA.AddSourceFile(weak1)
	x.chunkA.AddSourceFile(strong1)
	x.chunkA.AddSourceFile(weak2)
	x.chunkA.AddSourceFile(strong2)

	for _, input := range x.chunkA.GetInputs() {
		input.SetCompiler(x.compiler)
	}

	x.makeGraph()

	jsChunkGraphTest_containsExactly(t, "weak", jsChunkGraphTest_sourceFiles(x.getWeakChunk().GetInputs()), weak1, weak2)
	jsChunkGraphTest_containsExactly(t, "chunkA", jsChunkGraphTest_sourceFiles(x.chunkA.GetInputs()), strong1, strong2)
}

func TestJSChunkGraph_MoveMarkedWeakSourcesDuringManageDepsSortOnly(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()

	weak1, weak2, strong1, strong2 := jsChunkGraphTest_weakStrong()

	x.chunkA.AddSourceFile(weak1)
	x.chunkA.AddSourceFile(strong1)
	x.chunkA.AddSourceFile(weak2)
	x.chunkA.AddSourceFile(strong2)

	for _, input := range x.chunkA.GetInputs() {
		input.SetCompiler(x.compiler)
	}

	x.makeGraph()
	x.graph.ManageDependencies(x.compiler, DependencyOptions_SortOnly())

	jsChunkGraphTest_containsExactly(t, "weak", jsChunkGraphTest_sourceFiles(x.getWeakChunk().GetInputs()), weak1, weak2)
	jsChunkGraphTest_containsExactly(t, "chunkA", jsChunkGraphTest_sourceFiles(x.chunkA.GetInputs()), strong1, strong2)
}

func TestJSChunkGraph_IgnoreMarkedWeakSourcesDuringManageDepsPrune(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()

	weak1, weak2, strong1, strong2 := jsChunkGraphTest_weakStrong()

	x.chunkA.AddSourceFile(weak1)
	x.chunkA.AddSourceFile(strong1)
	x.chunkA.AddSourceFile(weak2)
	x.chunkA.AddSourceFile(strong2)

	for _, input := range x.chunkA.GetInputs() {
		input.SetCompiler(x.compiler)
	}

	x.makeGraph()
	x.graph.ManageDependencies(
		x.compiler,
		DependencyOptions_PruneForEntryPoints([]*ModuleIdentifier{
			ModuleIdentifier_ForFile("strong1"), ModuleIdentifier_ForFile("strong2")}))

	assertEqual(t, "weak", len(x.getWeakChunk().GetInputs()), 0)
	jsChunkGraphTest_containsExactly(t, "chunkA", jsChunkGraphTest_sourceFiles(x.chunkA.GetInputs()), strong1, strong2)
}

func TestJSChunkGraph_IgnoreDepsOfMarkedWeakSourcesDuringManageDepsPrune(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()

	weak := rhino.StaticSourceFile_SourceKind_WEAK
	weak1 := SourceFile_FromCodeWithKind("weak1", "goog.requireType('weak1weak');", weak)
	weak1weak := SourceFile_FromCodeWithKind("weak1weak", "goog.provide('weak1weak');", weak)
	weak2 := SourceFile_FromCodeWithKind("weak2", "goog.require('weak2strong');", weak)
	weak2strong := SourceFile_FromCodeWithKind("weak2strong", "goog.provide('weak2strong');", weak)
	strong1 := SourceFile_FromCodeWithKind("strong1", "", rhino.StaticSourceFile_SourceKind_STRONG)
	strong2 := SourceFile_FromCodeWithKind("strong2", "", rhino.StaticSourceFile_SourceKind_STRONG)

	x.chunkA.AddSourceFile(weak1)
	x.chunkA.AddSourceFile(strong1)
	x.chunkA.AddSourceFile(weak2)
	x.chunkA.AddSourceFile(strong2)
	x.chunkA.AddSourceFile(weak1weak)
	x.chunkA.AddSourceFile(weak2strong)

	for _, input := range x.chunkA.GetInputs() {
		input.SetCompiler(x.compiler)
	}

	x.makeGraph()
	x.graph.ManageDependencies(
		x.compiler,
		DependencyOptions_PruneForEntryPoints([]*ModuleIdentifier{
			ModuleIdentifier_ForFile("strong1"), ModuleIdentifier_ForFile("strong2")}))

	assertEqual(t, "weak", len(x.getWeakChunk().GetInputs()), 0)
	jsChunkGraphTest_containsExactly(t, "chunkA", jsChunkGraphTest_sourceFiles(x.chunkA.GetInputs()), strong1, strong2)
}

func TestJSChunkGraph_MoveImplicitWeakSourcesFromMoocherDuringManageDepsLegacyPrune(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()

	weak := SourceFile_FromCode("weak", "goog.provide('weak');")
	strong := SourceFile_FromCode("strong", "")
	moocher := SourceFile_FromCode("moocher", "goog.requireType('weak');")

	x.chunkA.AddSourceFile(weak)
	x.chunkA.AddSourceFile(strong)
	x.chunkA.AddSourceFile(moocher)

	for _, input := range x.chunkA.GetInputs() {
		input.SetCompiler(x.compiler)
	}

	x.makeGraph()
	x.graph.ManageDependencies(
		x.compiler,
		DependencyOptions_PruneLegacyForEntryPoints([]*ModuleIdentifier{ModuleIdentifier_ForFile("strong")}))

	jsChunkGraphTest_containsExactly(t, "weak", jsChunkGraphTest_sourceFiles(x.getWeakChunk().GetInputs()), weak)
	jsChunkGraphTest_containsExactly(t, "chunkA", jsChunkGraphTest_sourceFiles(x.chunkA.GetInputs()), strong, moocher)
}

func (x *jsChunkGraphTest) weakStrongRequireType() (weak1, weak2, strong1, strong2 *SourceFile) {
	weak1 = SourceFile_FromCode("weak1", "goog.provide('weak1');")
	weak2 = SourceFile_FromCode("weak2", "goog.provide('weak2');")
	strong1 = SourceFile_FromCode("strong1", "goog.requireType('weak1');")
	strong2 = SourceFile_FromCode("strong2", "goog.requireType('weak2');")

	x.chunkA.AddSourceFile(weak1)
	x.chunkA.AddSourceFile(strong1)
	x.chunkA.AddSourceFile(weak2)
	x.chunkA.AddSourceFile(strong2)

	for _, input := range x.chunkA.GetInputs() {
		input.SetCompiler(x.compiler)
	}
	return
}

func TestJSChunkGraph_ImplicitWeakSourcesNotMovedDuringManageDepsSortOnly(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	weak1, weak2, strong1, strong2 := x.weakStrongRequireType()

	x.makeGraph()
	x.graph.ManageDependencies(x.compiler, DependencyOptions_SortOnly())

	assertEqual(t, "weak", len(x.getWeakChunk().GetInputs()), 0)
	jsChunkGraphTest_containsExactly(t, "chunkA", jsChunkGraphTest_sourceFiles(x.chunkA.GetInputs()), weak1, strong1, weak2, strong2)
}

func TestJSChunkGraph_ImplicitWeakSourcesMovedDuringManageDepsPrune(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()
	weak1, weak2, strong1, strong2 := x.weakStrongRequireType()

	x.makeGraph()
	x.graph.ManageDependencies(
		x.compiler,
		DependencyOptions_PruneForEntryPoints([]*ModuleIdentifier{
			ModuleIdentifier_ForFile("strong1"), ModuleIdentifier_ForFile("strong2")}))

	jsChunkGraphTest_containsExactly(t, "weak", jsChunkGraphTest_sourceFiles(x.getWeakChunk().GetInputs()), weak1, weak2)
	jsChunkGraphTest_containsExactly(t, "chunkA", jsChunkGraphTest_sourceFiles(x.chunkA.GetInputs()), strong1, strong2)
}

func TestJSChunkGraph_TransitiveWeakSources(t *testing.T) {
	x := newJSChunkGraphTest(t)
	x.makeDeps()

	weak1 := SourceFile_FromCode(
		"weak1",
		"goog.provide('weak1'); goog.requireType('weak2'); goog.require('strongFromWeak');")
	strongFromWeak := SourceFile_FromCode("strongFromWeak", "goog.provide('strongFromWeak');")
	weak2 := SourceFile_FromCode("weak2", "goog.provide('weak2'); goog.requireType('weak3');")
	weak3 := SourceFile_FromCode("weak3", "goog.provide('weak3');")
	strong1 := SourceFile_FromCode("strong1", "goog.requireType('weak1');")

	x.chunkA.AddSourceFile(weak1)
	x.chunkA.AddSourceFile(strong1)
	x.chunkA.AddSourceFile(weak2)
	x.chunkA.AddSourceFile(weak3)
	x.chunkA.AddSourceFile(strongFromWeak)

	for _, input := range x.chunkA.GetInputs() {
		input.SetCompiler(x.compiler)
	}

	x.makeGraph()
	x.graph.ManageDependencies(
		x.compiler,
		DependencyOptions_PruneForEntryPoints([]*ModuleIdentifier{ModuleIdentifier_ForFile("strong1")}))

	jsChunkGraphTest_containsExactly(t, "weak", jsChunkGraphTest_sourceFiles(x.getWeakChunk().GetInputs()),
		weak1, weak2, weak3, strongFromWeak)
	jsChunkGraphTest_containsExactly(t, "chunkA", jsChunkGraphTest_sourceFiles(x.chunkA.GetInputs()), strong1)
}

func (x *jsChunkGraphTest) assertInputs(chunk *JSChunk, sourceNames ...string) {
	x.t.Helper()
	assertSliceEqual(x.t, "inputs of "+chunk.GetName(), jsChunkGraphTest_sourceNames(chunk.GetInputs()), sourceNames)
}

func jsChunkGraphTest_sourceNames(inputs []*CompilerInput) []string {
	inputNames := []string{}
	for _, input := range inputs {
		inputNames = append(inputNames, input.GetName())
	}
	return inputNames
}

func jsChunkGraphTest_code(sourceName string, provides []string, requires []string) *SourceFile {
	return jsChunkGraphTest_codeWithSource(sourceName, "", provides, requires)
}

func jsChunkGraphTest_codeWithSource(sourceName string, source string, provides []string, requires []string) *SourceFile {
	text := ""
	for _, p := range provides {
		text += "goog.provide('" + p + "');\n"
	}
	for _, r := range requires {
		text += "goog.require('" + r + "');\n"
	}
	return SourceFile_FromCode(sourceName, text+source)
}

func jsChunkGraphTest_list(strings ...string) []string {
	return append([]string{}, strings...)
}

func (x *jsChunkGraphTest) assertSmallestCoveringSubtree(expected *JSChunk, parentTree *JSChunk, chunks ...*JSChunk) {
	x.t.Helper()
	x.assertSmallestCoveringSubtreeInGraph(expected, x.graph, parentTree, chunks...)
}

func (x *jsChunkGraphTest) assertSmallestCoveringSubtreeInGraph(
	expected *JSChunk, graph *JSChunkGraph, parentTree *JSChunk, chunks ...*JSChunk) {
	x.t.Helper()
	chunksBitSet := jutil.NewBitSet()
	for _, m := range chunks {
		chunksBitSet.Set(m.GetIndex())
	}
	actual := graph.GetSmallestCoveringSubtree(parentTree, chunksBitSet)
	if actual != expected {
		x.t.Errorf("Smallest covering subtree of %v in %v should be %v but was %v",
			parentTree, chunksBitSet, expected, actual)
	}
}

func (x *jsChunkGraphTest) assertDeepestCommonDepInclusive(expected *JSChunk, m1 *JSChunk, m2 *JSChunk) {
	x.t.Helper()
	x.assertDeepestCommonDepOneWay(expected, m1, m2, true)
	x.assertDeepestCommonDepOneWay(expected, m2, m1, true)
}

func (x *jsChunkGraphTest) assertDeepestCommonDep(expected *JSChunk, m1 *JSChunk, m2 *JSChunk) {
	x.t.Helper()
	x.assertDeepestCommonDepOneWay(expected, m1, m2, false)
	x.assertDeepestCommonDepOneWay(expected, m2, m1, false)
}

func (x *jsChunkGraphTest) assertDeepestCommonDepOneWay(expected *JSChunk, m1 *JSChunk, m2 *JSChunk, inclusive bool) {
	x.t.Helper()
	var actual *JSChunk
	if inclusive {
		actual = x.graph.GetDeepestCommonDependencyInclusive(m1, m2)
	} else {
		actual = x.graph.GetDeepestCommonDependency(m1, m2)
	}
	if actual != expected {
		x.t.Errorf("Deepest common dep of %s and %s should be %v but was %v",
			m1.GetName(), m2.GetName(), expected, actual)
	}
}

func (x *jsChunkGraphTest) assertTransitiveDepsDeepestFirst(m *JSChunk, deps ...*JSChunk) {
	x.t.Helper()
	actual := x.graph.GetTransitiveDepsDeepestFirst(m)
	assertEqual(x.t, "transitive deps of "+m.GetName(), jutil.CollectionToString(actual), jutil.CollectionToString(deps))
}
