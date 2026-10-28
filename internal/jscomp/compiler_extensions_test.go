package jscomp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestCompiler_DefaultCodingConventionAndIdGeneratorMap(t *testing.T) {
	c := NewCompiler()
	if c.GetDefaultCodingConvention() == nil {
		t.Fatal("expected non-nil initial defaultCodingConvention")
	}
	googleCC := NewGoogleCodingConvention()
	c.SetDefaultCodingConvention(googleCC)
	if c.GetDefaultCodingConvention() != googleCC {
		t.Fatalf("GetDefaultCodingConvention() = %T, want *GoogleCodingConvention", c.GetDefaultCodingConvention())
	}
	if got := c.GetIdGeneratorMap(); got != "" {
		t.Fatalf("GetIdGeneratorMap() before compilation = %q, want empty", got)
	}

	opts := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
	opts.SetCodingConvention(nil)
	externs := []*SourceFile{SourceFile_FromCode("externs.js", "/** @externs */\n/** @param {*} x */\nfunction sink(x) {}\n")}
	srcs := []*SourceFile{SourceFile_FromCode("in.js", `
/** @idGenerator {unique} */
function getId(name) { return name; }
var CONSTANT_NAME = getId('alpha');
sink(CONSTANT_NAME);
`)}
	res := c.Compile(externs, srcs, opts)
	if !res.Success {
		t.Fatalf("Compile failed: %v", res.Errors)
	}
	if c.GetCodingConvention() != googleCC {
		t.Fatalf("Compiler.GetCodingConvention() = %T, want custom default *GoogleCodingConvention", c.GetCodingConvention())
	}
	gotMap := c.GetIdGeneratorMap()
	if gotMap == "" {
		t.Fatal("expected non-empty GetIdGeneratorMap() after @idGenerator compilation")
	}
	if res.idGeneratorMap == nil || gotMap != *res.idGeneratorMap {
		t.Fatalf("GetIdGeneratorMap() = %q, Result.idGeneratorMap = %v", gotMap, res.idGeneratorMap)
	}
}

func TestCompiler_RenderChunksAndSourceMapsParity(t *testing.T) {
	buildCompiler := func() (*Compiler, []*JSChunk) {
		c := NewCompiler()
		opts := NewCompilerOptions()
		CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
		smPath := "%outname%.map"
		opts.SetSourceMapOutputPath(&smPath)
		opts.SetSourceMapIncludeSourcesContent(true)

		externs := []*SourceFile{
			SourceFile_FromCode("externs.js", "/** @externs */\n/** @type {!Object} */ var window;\n/** @param {*} x */ function sink(x) {}\n"),
		}
		c0 := NewJSChunk("base")
		c0.Add(NewCompilerInput(SourceFile_FromCode("base0.js", "/** @license Shared License */\nfunction shared(x) { return x * 2 + 1; }\nwindow['b'] = shared(3);\n")))
		c0.Add(NewCompilerInput(SourceFile_FromCode("base1.js", "function baseExtra(y) { return shared(y) + 5; }\nsink(baseExtra(4));\n")))

		c1 := NewJSChunk("feat1")
		c1.AddDependency(c0)
		c1.Add(NewCompilerInput(SourceFile_FromCode("feat1.js", "/** @license Feat1 License */\nsink(shared(10) + 7);\n")))

		c2 := NewJSChunk("feat2")
		c2.AddDependency(c0)
		c2.Add(NewCompilerInput(SourceFile_FromCode("feat2.js", "/** @license Shared License */\nsink(shared(20) + 9);\n")))

		chunks := []*JSChunk{c0, c1, c2}
		res := c.CompileChunks(externs, chunks, opts)
		if !res.Success {
			t.Fatalf("CompileChunks failed: %v", res.Errors)
		}
		return c, c.GetChunks()
	}

	wrapperPrefix := func(m *JSChunk) string {
		return "/* header " + m.GetName() + " */\n(function(){"
	}
	chunkFileName := func(m *JSChunk) string {
		return "out/" + m.GetName() + ".bin.js"
	}

	seqCompiler, seqChunks := buildCompiler()
	seqLT := NewCompiler_ChunkGraphAwareLicenseTracker(seqCompiler)
	var want []RenderedChunk
	for _, m := range seqChunks {
		if m.GetName() == JSChunk_WEAK_CHUNK_NAME {
			continue
		}
		seqCompiler.ResetAndIntitializeSourceMap()
		seqLT.SetCurrentChunkContext(m)
		code := seqCompiler.ToSourceOfLicenseTrackerAndJSChunk(seqLT, m)
		sm := seqCompiler.GetSourceMap()
		sm.SetWrapperPrefix(wrapperPrefix(m))
		outName := chunkFileName(m)
		var sb strings.Builder
		sm.AppendTo(&sb, &outName)
		want = append(want, RenderedChunk{Chunk: m, Code: code, SourceMap: sb.String()})
	}

	for _, workers := range []int{0, 1, 4} {
		parCompiler, parChunks := buildCompiler()
		got, err := parCompiler.RenderChunksAndSourceMaps(parChunks, RenderChunkOptions{
			GenerateSourceMaps: true,
			ChunkFileName:      chunkFileName,
			WrapperPrefix:      wrapperPrefix,
			NumWorkers:         workers,
		})
		if err != nil {
			t.Fatalf("workers=%d RenderChunksAndSourceMaps error: %v", workers, err)
		}
		if len(got) != len(want) {
			t.Fatalf("workers=%d len(got)=%d, want %d", workers, len(got), len(want))
		}
		for i := range want {
			if got[i].Chunk.GetName() != want[i].Chunk.GetName() {
				t.Fatalf("workers=%d chunk[%d] name = %q, want %q", workers, i, got[i].Chunk.GetName(), want[i].Chunk.GetName())
			}
			if got[i].Code != want[i].Code {
				t.Fatalf("workers=%d chunk %s code mismatch:\ngot:  %q\nwant: %q", workers, want[i].Chunk.GetName(), got[i].Code, want[i].Code)
			}
			if got[i].SourceMap != want[i].SourceMap {
				t.Fatalf("workers=%d chunk %s sourcemap mismatch:\ngot:  %s\nwant: %s", workers, want[i].Chunk.GetName(), got[i].SourceMap, want[i].SourceMap)
			}
		}
	}
}

func TestThreadSafeDelegatingErrorManager_ConformanceViolationHandler(t *testing.T) {
	baseEM := NewPrintStreamErrorManager(jutil.NewPrintStream(&bytes.Buffer{}))
	tsEM := NewThreadSafeDelegatingErrorManager(baseEM)

	req := Requirement_NewBuilder().
		SetType(Requirement_Type_BANNED_NAME).
		SetErrorMessage("banned eval").
		AddValue("eval").
		Build()
	diag := JSError_Make(CheckConformance_CONFORMANCE_VIOLATION, "banned eval", "", "")

	// Without handler: delegates to base ErrorManager (false for RECORD_ONLY, true for REPORT_AS_BUILD_ERROR).
	if tsEM.ShouldReportConformanceViolation(
		req, nil, diag,
		ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_RECORD_ONLY,
		false,
	) {
		t.Fatal("expected default ShouldReportConformanceViolation(RECORD_ONLY) to return false")
	}
	if !tsEM.ShouldReportConformanceViolation(
		req, nil, diag,
		ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_REPORT_AS_BUILD_ERROR,
		false,
	) {
		t.Fatal("expected default ShouldReportConformanceViolation(REPORT_AS_BUILD_ERROR) to return true")
	}

	var recorded int
	suppressAll := true
	tsEM.SetConformanceViolationHandler(func(
		r *Requirement,
		allowlistEntry *RequirementScopeEntry,
		d *JSError,
		behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior,
		isAllowlisted bool,
	) bool {
		recorded++
		if suppressAll {
			return false
		}
		return behavior != ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_RECORD_ONLY
	})

	if tsEM.ShouldReportConformanceViolation(
		req, nil, diag,
		ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_REPORT_AS_BUILD_ERROR,
		false,
	) {
		t.Fatal("expected custom handler to suppress violation when suppressAll=true")
	}
	suppressAll = false
	if !tsEM.ShouldReportConformanceViolation(
		req, nil, diag,
		ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_REPORT_AS_BUILD_ERROR,
		false,
	) {
		t.Fatal("expected REPORT_AS_BUILD_ERROR violation to return true when suppressAll=false")
	}
	if recorded != 2 {
		t.Fatalf("handler invoked %d times, want 2", recorded)
	}
}

func TestCompiler_ConformanceViolationsSaveAndRestoreState(t *testing.T) {
	stage1 := NewCompiler()
	opts1 := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts1)
	externs := []*SourceFile{SourceFile_FromCode("externs.js", "/** @externs */\nvar window;\n")}
	srcs := []*SourceFile{SourceFile_FromCode("app.js", "window['x'] = 1 + 2;\n")}
	stage1.Init(externs, srcs, opts1)
	stage1.ParseForCompilation()
	stage1.Stage1Passes()
	if stage1.HasErrors() {
		t.Fatalf("Stage1Passes errors: %v", stage1.GetErrors())
	}

	v1 := ViolationProto_NewBuilder().SetSourceName("app.js").SetLineNumber(1).SetCharNumber(0).Build()
	attributed := ConformanceViolationsProto_NewBuilder().
		AddViolations(v1).
		Build()
	unattributed := ViolationProto_NewBuilder().SetSourceName("global.js").SetLineNumber(42).Build()
	stage1.AddConformanceViolations(attributed)
	stage1.AddUnattributedConformanceViolation(unattributed)

	var stateBuf bytes.Buffer
	stage1.SaveState(&stateBuf)

	stage2 := NewCompiler()
	opts2 := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts2)
	stage2.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "/** @externs */\nvar window;\n")},
		[]*SourceFile{SourceFile_FromCode("app.js", "window['x'] = 1 + 2;\n")},
		opts2,
	)
	stage2.RestoreState(bytes.NewReader(stateBuf.Bytes()))

	gotAttr := stage2.GetConformanceViolations()
	if len(gotAttr) != 1 || gotAttr[0].GetViolationsCount() != 1 || gotAttr[0].GetViolations(0).GetSourceName() != "app.js" {
		t.Fatalf("restored ConformanceViolations = %+v", gotAttr)
	}
	gotUnattr := stage2.GetUnattributedConformanceViolations()
	if len(gotUnattr) != 1 || gotUnattr[0].GetSourceName() != "global.js" || gotUnattr[0].GetLineNumber() != 42 {
		t.Fatalf("restored UnattributedConformanceViolations = %+v", gotUnattr)
	}
}
