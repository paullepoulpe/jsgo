package jscomp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestCrossChunkSingleStrong_HasSingleActiveStrongChunk(t *testing.T) {
	var nilGraph *JSChunkGraph
	if nilGraph.hasSingleActiveStrongChunk() {
		t.Fatal("expected nil JSChunkGraph to report false")
	}

	// Single strong chunk: NewJSChunkGraph appends an empty $weak$ chunk (GetChunkCount() == 2).
	c0 := NewJSChunk("c0")
	c0.AddSourceFile(SourceFile_FromCode("a.js", "var a = 1;"))
	gSingle := NewJSChunkGraph([]*JSChunk{c0})
	if gSingle.GetChunkCount() != 2 {
		t.Fatalf("expected 2 chunks (c0 + $weak$), got %d", gSingle.GetChunkCount())
	}
	if !gSingle.hasSingleActiveStrongChunk() {
		t.Fatal("expected single strong chunk + empty $weak$ to report true")
	}

	// Single strong chunk + populated $weak$ chunk (weak source moved into $weak$).
	compiler := NewCompiler()
	cStrong := NewJSChunk("c0")
	cStrong.AddSourceFile(SourceFile_FromCodeWithKind("strong.js", "var a = 1;", rhino.StaticSourceFile_SourceKind_STRONG))
	cStrong.AddSourceFile(SourceFile_FromCodeWithKind("weak.js", "var b = a;", rhino.StaticSourceFile_SourceKind_WEAK))
	for _, in := range cStrong.GetInputs() {
		in.SetCompiler(compiler)
	}
	gWithWeak := NewJSChunkGraph([]*JSChunk{cStrong})
	if gWithWeak.GetChunkCount() != 2 {
		t.Fatalf("expected 2 chunks, got %d", gWithWeak.GetChunkCount())
	}
	if gWithWeak.hasSingleActiveStrongChunk() {
		t.Fatal("expected populated $weak$ chunk to report false")
	}

	// Two strong chunks + empty $weak$ chunk.
	m0 := NewJSChunk("m0")
	m0.AddSourceFile(SourceFile_FromCode("m0.js", "var a = 1;"))
	m1 := NewJSChunk("m1")
	m1.AddDependency(m0)
	m1.AddSourceFile(SourceFile_FromCode("m1.js", "var b = a;"))
	gTwoStrong := NewJSChunkGraph([]*JSChunk{m0, m1})
	if gTwoStrong.hasSingleActiveStrongChunk() {
		t.Fatal("expected two strong chunks to report false")
	}
}

func TestCrossChunkMethodMotion_SingleStrongChunkFastPath(t *testing.T) {
	// 1. Standard 1-chunk build (with auto-created empty $weak$ chunk): fast return, 0 stub IDs, 0 AST changes.
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	chunk := NewJSChunk("main")
	chunk.AddSourceFile(SourceFile_FromCode("a.js",
		"function Foo() {}\n"+
			"Foo.prototype.bar = function() { return 1; };\n"+
			"class Baz { qux() { return 2; } }\n"+
			"alert((new Foo()).bar() + (new Baz()).qux());\n"))
	compiler.InitChunks([]*SourceFile{SourceFile_FromCode("externs.js", "function alert(x) {}")}, []*JSChunk{chunk}, opts)
	compiler.ParseInputs()
	if compiler.GetChunkGraph().GetChunkCount() != 2 || !compiler.GetChunkGraph().hasSingleActiveStrongChunk() {
		t.Fatalf("unexpected chunk graph setup: count=%d", compiler.GetChunkGraph().GetChunkCount())
	}

	// Unnormalized lifecycle check must still panic with IllegalStateException.
	idGen := NewIdGenerator()
	pass := NewCrossChunkMethodMotion(compiler, idGen, false, false)
	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(*jutil.IllegalStateException); ok {
					panicked = true
					return
				}
				panic(r)
			}
		}()
		pass.Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	}()
	if !panicked {
		t.Fatal("expected IllegalStateException when compiler is not normalized")
	}

	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	beforeChange := compiler.GetChangeTracker().GetChangeStamp()
	beforeSrc := compiler.ToSource()
	pass.Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	if idGen.HasGeneratedAnyIds() {
		t.Fatal("expected zero generated stub IDs on single strong chunk")
	}
	if got := compiler.GetChangeTracker().GetChangeStamp(); got != beforeChange {
		t.Fatalf("expected unchanged stamp %d, got %d", beforeChange, got)
	}
	if got := compiler.ToSource(); got != beforeSrc {
		t.Fatalf("expected identical source:\nwant:\n%s\ngot:\n%s", beforeSrc, got)
	}

	// 2. Populated $weak$ chunk (1 strong + 1 weak file) must still move methods into $weak$.
	compilerWeak := NewCompiler()
	optsWeak := NewCompilerOptions()
	cStrong := NewJSChunk("main")
	cStrong.AddSourceFile(SourceFile_FromCodeWithKind("strong.js",
		"function Foo() {}\nFoo.prototype.onlyInWeak = function() { return 42; };\n",
		rhino.StaticSourceFile_SourceKind_STRONG))
	cStrong.AddSourceFile(SourceFile_FromCodeWithKind("weak.js",
		"(new Foo()).onlyInWeak();\n",
		rhino.StaticSourceFile_SourceKind_WEAK))
	compilerWeak.InitChunks([]*SourceFile{SourceFile_FromCode("externs.js", "")}, []*JSChunk{cStrong}, optsWeak)
	compilerWeak.ParseInputs()
	compilerWeak.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	idGenWeak := NewIdGenerator()
	NewCrossChunkMethodMotion(compilerWeak, idGenWeak, false, false).Process(
		compilerWeak.GetExternsRoot(), compilerWeak.GetJsRoot())
	if !idGenWeak.HasGeneratedAnyIds() {
		t.Fatal("expected method motion into populated $weak$ chunk to generate stub ID")
	}
	weakChunk := compilerWeak.GetChunkGraph().GetChunkByName(JSChunk_WEAK_CHUNK_NAME)
	weakCode := compilerWeak.ToSourceOfJSChunk(weakChunk)
	if !strings.Contains(weakCode, "JSCompiler_unstubMethod") {
		t.Fatalf("expected unstub call in $weak$ chunk, got: %s", weakCode)
	}
}

func TestCrossChunkCodeMotion_SingleStrongChunkFastPathAndInstanceofParity(t *testing.T) {
	// 1. Single strong chunk + debug logging enabled: preserves indexed cccm.log ("doing cccm").
	logDir := t.TempDir()
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetDebugLogDirectory(jutil.Path_Of(logDir))
	opts.SetDebugLogFilter("CrossChunkCodeMotion")
	c0 := NewJSChunk("main")
	c0.AddSourceFile(SourceFile_FromCode("a.js", "function unused() { return 1; }\nvar x = 2;\nalert(x);\n"))
	compiler.InitChunks([]*SourceFile{SourceFile_FromCode("externs.js", "function alert(x) {}")}, []*JSChunk{c0}, opts)
	compiler.ParseInputs()

	beforeSrc := compiler.ToSource()
	beforeStamp := compiler.GetChangeTracker().GetChangeStamp()
	cccm := NewCrossChunkCodeMotion(compiler, compiler.GetChunkGraph(), false)
	cccm.Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	cccm.Process(compiler.GetExternsRoot(), compiler.GetJsRoot())

	if got := compiler.GetChangeTracker().GetChangeStamp(); got != beforeStamp {
		t.Fatalf("expected no AST change on 1-chunk build, stamp %d -> %d", beforeStamp, got)
	}
	if got := compiler.ToSource(); got != beforeSrc {
		t.Fatalf("expected identical source:\nwant:\n%s\ngot:\n%s", beforeSrc, got)
	}

	// Verify both indexed log files were created with exact "doing cccm" content.
	var cccmLogs []string
	_ = filepath.WalkDir(logDir, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.Contains(d.Name(), "cccm") {
			cccmLogs = append(cccmLogs, path)
		}
		return nil
	})
	if len(cccmLogs) != 1 {
		t.Fatalf("expected 1 reopened cccm.log file, got %d (%v)", len(cccmLogs), cccmLogs)
	}
	data, err := os.ReadFile(cccmLogs[0])
	if err != nil {
		t.Fatal(err)
	}
	if content := strings.TrimSpace(string(data)); content != "doing cccm\ndoing cccm" {
		t.Fatalf("expected log %s to contain two 'doing cccm' lines, got %q", cccmLogs[0], content)
	}

	// 2. parentChunkCanSeeSymbolsDeclaredInChildren = true WITH instanceof on non-top-level symbol
	// must still insert typeof guard even on a 1-chunk build!
	compilerInst := NewCompiler()
	optsInst := NewCompilerOptions()
	cInst := NewJSChunk("main")
	cInst.AddSourceFile(SourceFile_FromCode("a.js",
		"if (cond) { var F = function() {}; }\nalert(x instanceof F);\n"))
	compilerInst.InitChunks(
		[]*SourceFile{SourceFile_FromCode("externs.js", "var cond, x; function alert(a) {}")},
		[]*JSChunk{cInst},
		optsInst,
	)
	compilerInst.ParseInputs()
	NewCrossChunkCodeMotion(compilerInst, compilerInst.GetChunkGraph(), true).Process(
		compilerInst.GetExternsRoot(), compilerInst.GetJsRoot())
	gotInst := compilerInst.ToSource()
	if !strings.Contains(gotInst, `"function"==typeof F&&x instanceof F`) &&
		!strings.Contains(gotInst, `'function' == typeof F && x instanceof F`) {
		t.Fatalf("expected instanceof typeof guard to be added, got: %s", gotInst)
	}

	// 3. parentChunkCanSeeSymbolsDeclaredInChildren = true WITHOUT instanceof takes fast path cleanly.
	compilerNoInst := NewCompiler()
	cNoInst := NewJSChunk("main")
	cNoInst.AddSourceFile(SourceFile_FromCode("a.js", "function F() {}\nalert(new F());\n"))
	compilerNoInst.InitChunks(
		[]*SourceFile{SourceFile_FromCode("externs.js", "function alert(a) {}")},
		[]*JSChunk{cNoInst},
		NewCompilerOptions(),
	)
	compilerNoInst.ParseInputs()
	stampBefore := compilerNoInst.GetChangeTracker().GetChangeStamp()
	NewCrossChunkCodeMotion(compilerNoInst, compilerNoInst.GetChunkGraph(), true).Process(
		compilerNoInst.GetExternsRoot(), compilerNoInst.GetJsRoot())
	if got := compilerNoInst.GetChangeTracker().GetChangeStamp(); got != stampBefore {
		t.Fatalf("expected no change without instanceof, stamp %d -> %d", stampBefore, got)
	}

	// 4. 1 strong chunk + populated $weak$ chunk still moves global declarations into $weak$.
	compilerWeak := NewCompiler()
	cWeak := NewJSChunk("main")
	cWeak.AddSourceFile(SourceFile_FromCodeWithKind("strong.js",
		"function movableHelper() { return 7; }\nvar keepStrong = 1;\nalert(keepStrong);\n",
		rhino.StaticSourceFile_SourceKind_STRONG))
	cWeak.AddSourceFile(SourceFile_FromCodeWithKind("weak.js",
		"alert(movableHelper());\n",
		rhino.StaticSourceFile_SourceKind_WEAK))
	compilerWeak.InitChunks(
		[]*SourceFile{SourceFile_FromCode("externs.js", "function alert(a) {}")},
		[]*JSChunk{cWeak},
		NewCompilerOptions(),
	)
	compilerWeak.ParseInputs()
	NewCrossChunkCodeMotion(compilerWeak, compilerWeak.GetChunkGraph(), false).Process(
		compilerWeak.GetExternsRoot(), compilerWeak.GetJsRoot())
	weakOut := compilerWeak.ToSourceOfJSChunk(compilerWeak.GetChunkGraph().GetChunkByName(JSChunk_WEAK_CHUNK_NAME))
	if !strings.Contains(weakOut, "function movableHelper") {
		t.Fatalf("expected movableHelper to move into populated $weak$ chunk, got: %s", weakOut)
	}
}
