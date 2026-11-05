package jscomp_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/resources"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type testCustomPass struct {
	called *int
}

func (p testCustomPass) Process(_, _ rhino.Node) {
	*p.called++
}

type testCssMap struct{}

func (testCssMap) Get(value string) (string, bool) { return "css_" + value, true }
func (testCssMap) GetStyle() jscomp.CssRenamingMap_Style {
	return jscomp.CssRenamingMap_Style_BY_WHOLE
}

type forwardedHolder struct {
	fn rhino.Node
}

func (h *forwardedHolder) ForwardASTNodes(fwd func(rhino.Node) rhino.Node) {
	h.fn = fwd(h.fn)
}

func TestCompilerOptionsEmbedderSettersAndTypedAstStream(t *testing.T) {
	var customPassCalls int
	var reconcileHookCalled bool
	var typedAstBuf bytes.Buffer

	opts := jscomp.NewCompilerOptions()
	opts.SetChecksOnly(true)
	opts.SetCheckTypes(true)
	opts.AddCustomPass(jscomp.CustomPassExecutionTime_BEFORE_CHECKS, testCustomPass{called: &customPassCalls})
	opts.SetCssRenamingMap(testCssMap{})
	opts.SetNameGenerator(jscomp.NewDefaultNameGenerator())
	opts.SetXidHashFunction(func(s string) int32 { return int32(len(s)) })
	opts.SetChunkIdHashFunction(func(s string) int32 { return int32(len(s) + 1) })
	opts.SetTypedAstOutputStream(&typedAstBuf)
	opts.SetLicenseFilter(func(lic string) bool {
		return !strings.Contains(lic, "DROP_ME")
	})
	opts.SetReconcileOptionsHook(func(o *jscomp.CompilerOptions) {
		if o == opts {
			reconcileHookCalled = true
		}
	})

	var compressCalled, decompressCalled bool
	opts.SetStateCompressionWrapper(func(w io.Writer) io.Writer {
		compressCalled = true
		return w
	})
	opts.SetStateDecompressionWrapper(func(r io.Reader) io.Reader {
		decompressCalled = true
		return r
	})
	if opts.GetStateCompressionWrapper() == nil || opts.GetStateDecompressionWrapper() == nil {
		t.Fatal("expected state compression/decompression wrappers to be non-nil")
	}

	c := jscomp.NewCompiler()
	inFile := jscomp.SourceFile_FromCode("in.js", "/** @license KEEP_ME */\n/** @type {number} */ var x = 1;\n")
	res := c.Compile(
		[]*jscomp.SourceFile{},
		[]*jscomp.SourceFile{inFile},
		opts,
	)
	if !res.Success {
		t.Fatalf("compile failed: %v", res.Errors)
	}
	if customPassCalls != 1 {
		t.Fatalf("customPassCalls = %d, want 1", customPassCalls)
	}
	if !reconcileHookCalled {
		t.Fatal("expected ReconcileOptionsHook to be invoked")
	}
	if typedAstBuf.Len() == 0 {
		t.Fatal("expected in-memory TypedAstOutputStream to receive serialized TypedAST bytes")
	}

	var stateBuf bytes.Buffer
	c.SaveState(&stateBuf)
	if !compressCalled {
		t.Fatal("expected StateCompressionWrapper to be called during SaveState")
	}
	c2 := jscomp.NewCompiler()
	c2.Init([]*jscomp.SourceFile{}, []*jscomp.SourceFile{inFile}, opts)
	c2.RestoreState(&stateBuf)
	if !decompressCalled {
		t.Fatal("expected StateDecompressionWrapper to be called during RestoreState")
	}

	// Verify licenseFilter suppresses DROP_ME while keeping KEEP_ME:
	optsCodegen := jscomp.NewCompilerOptions()
	optsCodegen.SetLicenseFilter(func(lic string) bool {
		return !strings.Contains(lic, "DROP_ME")
	})
	c3 := jscomp.NewCompiler()
	c3.Compile(
		[]*jscomp.SourceFile{},
		[]*jscomp.SourceFile{
			jscomp.SourceFile_FromCode("a.js", "/** @license DROP_ME */\nvar a = 1;\n"),
			jscomp.SourceFile_FromCode("b.js", "/** @license KEEP_ME */\nvar b = 2;\n"),
		},
		optsCodegen,
	)
	src := c3.ToSource()
	if strings.Contains(src, "DROP_ME") || !strings.Contains(src, "KEEP_ME") {
		t.Fatalf("unexpected licenseFilter output: %q", src)
	}
}

func TestPassConfigDecoratorsPassListBuilderAndCustomNodeMetadataCompaction(t *testing.T) {
	opts := jscomp.NewCompilerOptions()
	jscomp.CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)

	var decoratorOrder []string
	var holder *forwardedHolder
	var metaMap *jscomp.OrderedMap[rhino.Node, any]

	c := jscomp.NewCompiler()
	pc := jscomp.NewDefaultPassConfig(opts)
	pc.SetChecksDecorator(func(b *jscomp.PassListBuilder) {
		if b.Len() == 0 || b.At(0) == nil {
			t.Fatal("expected non-empty checks PassListBuilder")
		}
		b.InsertAt(0, jscomp.NewPassFactory("embedderRecordMeta", func(comp jscomp.AbstractCompiler) jscomp.CompilerPass {
			return jscomp.CompilerPassFunc(func(_, jsRoot rhino.Node) {
				decoratorOrder = append(decoratorOrder, "checks")
				metaMap = comp.CustomNodeMetadata()
				// Record every FUNCTION node before AST compactions run.
				var walk func(rhino.Node)
				walk = func(n rhino.Node) {
					if n.IsFunction() {
						nameNode := n.GetFirstChild()
						name := ""
						if !nameNode.IsNil() && nameNode.IsName() {
							name = nameNode.GetString()
						}
						if holder == nil {
							holder = &forwardedHolder{fn: n}
							metaMap.Put(n, holder)
						} else {
							metaMap.Put(n, []rhino.Node{n})
						}
						_ = name
					}
					for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
						walk(ch)
					}
				}
				walk(jsRoot)
			})
		}))
	})
	pc.SetOptimizationsDecorator(func(b *jscomp.PassListBuilder) {
		b.InsertAt(b.Len(), jscomp.NewPassFactory("embedderPostOpt", func(jscomp.AbstractCompiler) jscomp.CompilerPass {
			return jscomp.CompilerPassFunc(func(_, _ rhino.Node) {
				decoratorOrder = append(decoratorOrder, "optimizations")
			})
		}))
	})
	pc.SetFinalizationsDecorator(func(b *jscomp.PassListBuilder) {
		b.InsertAt(b.Len(), jscomp.NewPassFactory("embedderPostFinal", func(jscomp.AbstractCompiler) jscomp.CompilerPass {
			return jscomp.CompilerPassFunc(func(_, _ rhino.Node) {
				decoratorOrder = append(decoratorOrder, "finalizations")
			})
		}))
	})
	c.SetPassConfig(pc)

	// Include dead code so AST compaction (>10% dead + typesToColors preorder compaction) definitely runs.
	var deadBuilder strings.Builder
	for i := range 200 {
		deadBuilder.WriteString("function unusedFn_")
		deadBuilder.WriteString(strconv.Itoa(i))
		deadBuilder.WriteString("(a, b) { return a + b + 123; }\n")
	}
	deadBuilder.WriteString("window['keep'] = function keptFn(v) { return v * 2; };\n")

	res := c.Compile(
		[]*jscomp.SourceFile{jscomp.SourceFile_FromCode("ext.js", "/** @externs */ var window;")},
		[]*jscomp.SourceFile{jscomp.SourceFile_FromCode("app.js", deadBuilder.String())},
		opts,
	)
	if !res.Success {
		t.Fatalf("compile failed: %v", res.Errors)
	}
	if strings.Join(decoratorOrder, ",") != "checks,optimizations,finalizations" {
		t.Fatalf("decoratorOrder = %v", decoratorOrder)
	}
	if metaMap == nil || metaMap.IsEmpty() || holder == nil || holder.fn.IsNil() || !holder.fn.IsFunction() {
		t.Fatalf("CustomNodeMetadata was not preserved/forwarded across compaction: holder=%+v", holder)
	}
	for k := range metaMap.Keys() {
		if k.IsNil() || !k.IsFunction() {
			t.Fatalf("forwarded key is not a valid FUNCTION node: %v", k)
		}
	}
}

func TestOptimizeCallsFunctionNodeHelpers(t *testing.T) {
	c := jscomp.NewCompiler()
	opts := jscomp.NewCompilerOptions()
	opts.AddCustomPass(jscomp.CustomPassExecutionTime_BEFORE_CHECKS, jscomp.CompilerPassFunc(func(_, _ rhino.Node) {
		script := c.GetScriptNode("a.js")
		fnName := script.GetFirstChild().GetFirstChild()
		varName := script.GetSecondChild().GetFirstChild()

		fn, count := jscomp.OptimizeCalls_GetSingleFunctionNode([]rhino.Node{fnName})
		if count != 1 || fn.IsNil() || !fn.IsFunction() {
			t.Fatalf("OptimizeCalls_GetSingleFunctionNode(fnName) = (%v, %d)", fn, count)
		}
		fns := jscomp.OptimizeCalls_AppendDefinitionFunctionNodesFor(nil, varName)
		if len(fns) != 1 || !fns[0].IsFunction() {
			t.Fatalf("OptimizeCalls_AppendDefinitionFunctionNodesFor(varName) = %v", fns)
		}
		_, multiCount := jscomp.OptimizeCalls_GetSingleFunctionNode([]rhino.Node{fnName, varName})
		if multiCount < 2 {
			t.Fatalf("expected multiCount >= 2, got %d", multiCount)
		}
	}))
	res := c.Compile([]*jscomp.SourceFile{}, []*jscomp.SourceFile{
		jscomp.SourceFile_FromCode("a.js", "function foo(x) { return x; } var bar = function(y) { return y; };"),
	}, opts)
	if !res.Success {
		t.Fatalf("compile failed: %v", res.Errors)
	}
}

type customBanAlertRule struct {
	jscomp.ConformanceRules_AbstractRule
}

func TestCustomConformanceRuleInterceptorAndProtoAccessors(t *testing.T) {
	const customClass = "com.example.CustomBanMagicLitRule"
	jscomp.RegisterCustomConformanceRule(customClass, func(c jscomp.AbstractCompiler, req *jscomp.Requirement) (jscomp.CheckConformance_Rule, error) {
		r := &customBanAlertRule{}
		err := r.InitRule(c, req, false, func(_ *jscomp.NodeTraversal, n rhino.Node) *jscomp.ConformanceRules_ConformanceResult {
			if n.IsStringLit() && n.GetString() == "FORBIDDEN_TOKEN" {
				res := jscomp.NewConformanceRules_ConformanceResultWithNote(
					jscomp.ConformanceRules_ConformanceLevel_VIOLATION, "found forbidden token")
				if !res.IsViolation() || res.GetNote() != "found forbidden token" {
					t.Fatal("ConformanceResult accessor mismatch")
				}
				return res
			}
			return jscomp.ConformanceRules_ConformanceResult_CONFORMANCE
		})
		return r, err
	})
	defer jscomp.UnregisterCustomConformanceRule(customClass)

	scopeEntry := jscomp.RequirementScopeEntry_NewBuilder().
		SetReason(jscomp.RequirementScopeEntry_Reason_LEGACY).
		SetExplanation("legacy dir").
		AddPrefix("legacy/").
		Build()
	req := jscomp.Requirement_NewBuilder().
		SetType(jscomp.Requirement_Type_CUSTOM).
		SetJavaClass(customClass).
		SetErrorMessage("Forbidden literal used").
		AddAllowlistEntry(scopeEntry).
		Build()
	origCfg := jscomp.ConformanceConfig_NewBuilder().
		AddRequirement(req).
		SetLibraryLevelNonAllowlistedConformanceViolationsBehavior(
			jscomp.ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_RECORD_ONLY).
		Build()

	parsedCfg, err := jscomp.ParseConformanceConfigBytes(origCfg.ToByteArray())
	if err != nil {
		t.Fatalf("ParseConformanceConfigBytes failed: %v", err)
	}
	clonedCfg := parsedCfg.Clone()
	if clonedCfg.GetRequirementCount() != 1 || clonedCfg.GetRequirement(0).GetErrorMessage() != "Forbidden literal used" {
		t.Fatal("ConformanceConfig Clone/GetRequirement mismatch")
	}

	opts := jscomp.NewCompilerOptions()
	opts.SetConformanceConfigs([]*jscomp.ConformanceConfig{clonedCfg})

	var intercepted []*jscomp.ViolationProto
	c := jscomp.NewCompiler()
	c.SetConformanceViolationInterceptor(func(
		r *jscomp.Requirement,
		allowlist *jscomp.RequirementScopeEntry,
		diag *jscomp.JSError,
		_ jscomp.ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior,
		_ bool,
	) bool {
		srcName, _ := diag.SourceName()
		vp := jscomp.ViolationProto_NewBuilder().
			SetRequirement(r).
			SetSourceName(srcName).
			SetLineNumber(int32(diag.GetLineNumber())).
			SetCharNumber(int32(diag.GetCharno())).
			Build()
		if allowlist != nil {
			vp.SetAllowlistEntry(allowlist)
		}
		intercepted = append(intercepted, vp)
		return true
	})

	res := c.Compile(
		[]*jscomp.SourceFile{},
		[]*jscomp.SourceFile{
			jscomp.SourceFile_FromCode("src/bad.js", "var s = 'FORBIDDEN_TOKEN';\n"),
			jscomp.SourceFile_FromCode("legacy/old.js", "var s2 = 'FORBIDDEN_TOKEN';\n"),
		},
		opts,
	)
	if !res.Success {
		t.Fatalf("expected warnings not errors: %v", res.Errors)
	}
	if len(intercepted) != 2 {
		t.Fatalf("intercepted %d violations, want 2", len(intercepted))
	}

	cvProto := jscomp.ConformanceViolationsProto_NewBuilder().
		SetInputName("target_lib").
		AddViolations(intercepted[0]).
		AddViolations(intercepted[1]).
		Build()
	roundTrip, err := jscomp.ParseConformanceViolationsBytes(cvProto.Clone().ToByteArray())
	if err != nil {
		t.Fatalf("ParseConformanceViolationsBytes: %v", err)
	}
	vList := roundTrip.GetViolationsList()
	if len(vList) != 2 ||
		vList[0].GetSourceName() != "src/bad.js" ||
		vList[0].GetLineNumber() != 1 ||
		vList[0].GetCharNumber() < 0 ||
		vList[0].GetRequirement().GetJavaClass() != customClass ||
		vList[1].GetAllowlistEntry().GetExplanation() != "legacy dir" {
		t.Fatalf("unexpected round-tripped ConformanceViolationsProto: %+v", vList)
	}
}

func TestRenderChunksPostProcessCompressAndRuntimeOverrides(t *testing.T) {
	resources.RegisterOverride(resources.JsDir+"embedder_custom_lib.js", []byte("var embedderRuntimeInjected = 99;\n"))
	defer resources.ClearOverrides()

	prevDir := jscomp.GetRuntimeLibDir()
	jscomp.SetRuntimeLibDir("custom/runtime/js")
	defer jscomp.SetRuntimeLibDir(prevDir)
	if got := jscomp.GetRuntimeLibDir(); got != "custom/runtime/js/" {
		t.Fatalf("GetRuntimeLibDir = %q, want custom/runtime/js/", got)
	}

	opts := jscomp.NewCompilerOptions()
	opts.SetNumParallelThreads(2)

	c1 := jscomp.NewJSChunk("c1")
	c1.Add(jscomp.NewCompilerInput(jscomp.SourceFile_FromCode("a.js", "var a = 1;\n")))
	c2 := jscomp.NewJSChunk("c2")
	c2.AddDependency(c1)
	c2.Add(jscomp.NewCompilerInput(jscomp.SourceFile_FromCode("b.js", "var b = a + 2;\n")))
	chunks := []*jscomp.JSChunk{c1, c2}

	comp := jscomp.NewCompiler()
	res := comp.CompileChunks([]*jscomp.SourceFile{}, chunks, opts)
	if !res.Success {
		t.Fatalf("CompileChunks failed: %v", res.Errors)
	}

	rendered, err := comp.RenderChunksAndSourceMaps(chunks, jscomp.RenderChunkOptions{
		GenerateSourceMaps: true,
		PostProcessInput: func(in *jscomp.CompilerInput, code string) string {
			return "/*in:" + in.GetName() + "*/\n" + code
		},
		CompressSourceMap: func(_ *jscomp.JSChunk, sm string) ([]byte, error) {
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			if _, wErr := gz.Write([]byte(sm)); wErr != nil {
				return nil, wErr
			}
			if cErr := gz.Close(); cErr != nil {
				return nil, cErr
			}
			return buf.Bytes(), nil
		},
	})
	if err != nil {
		t.Fatalf("RenderChunksAndSourceMaps error: %v", err)
	}
	if len(rendered) != 2 {
		t.Fatalf("len(rendered) = %d, want 2", len(rendered))
	}
	if !strings.HasPrefix(rendered[0].Code, "/*in:a.js*/\n") || !strings.HasPrefix(rendered[1].Code, "/*in:b.js*/\n") {
		t.Fatalf("PostProcessInput not applied: %q / %q", rendered[0].Code, rendered[1].Code)
	}
	for i, rc := range rendered {
		if len(rc.SourceMapBytes) == 0 {
			t.Fatalf("chunk %d missing compressed SourceMapBytes", i)
		}
		zr, zErr := gzip.NewReader(bytes.NewReader(rc.SourceMapBytes))
		if zErr != nil {
			t.Fatal(zErr)
		}
		raw, rErr := io.ReadAll(zr)
		_ = zr.Close()
		if rErr != nil {
			t.Fatal(rErr)
		}
		if string(raw) != rc.SourceMap {
			t.Fatalf("chunk %d decompressed SourceMapBytes mismatch", i)
		}
	}
}

func TestRenderChunksLargeBaseStreamingAndOmitSourceMapString(t *testing.T) {
	opts := jscomp.NewCompilerOptions()
	opts.SetNumParallelThreads(4)
	opts.SetSourceMapOutputPath("out.map")

	const baseInputCount = 310 // Exceeds chunkRenderMaxBatchScripts (256) to exercise sub-batch streaming.
	baseChunk := jscomp.NewJSChunk("base")
	for i := range baseInputCount {
		name := fmt.Sprintf("base_%03d.js", i)
		code := fmt.Sprintf("/** @license SharedBaseLic */\nvar baseVar_%d = %d;\n", i, i)
		if i%50 == 0 {
			code = fmt.Sprintf("/** @license UniqueLic_%d */\nvar baseVar_%d = %d;\n", i, i, i)
		}
		baseChunk.Add(jscomp.NewCompilerInput(jscomp.SourceFile_FromCode(name, code)))
	}

	chunks := []*jscomp.JSChunk{baseChunk}
	for cIdx := range 5 {
		ch := jscomp.NewJSChunk(fmt.Sprintf("leaf_%d", cIdx))
		ch.AddDependency(baseChunk)
		ch.Add(jscomp.NewCompilerInput(jscomp.SourceFile_FromCode(
			fmt.Sprintf("leaf_%d.js", cIdx),
			fmt.Sprintf("/** @license SharedBaseLic */\nvar leafVar_%d = baseVar_0 + %d;\n", cIdx, cIdx),
		)))
		chunks = append(chunks, ch)
	}

	comp := jscomp.NewCompiler()
	res := comp.CompileChunks([]*jscomp.SourceFile{}, chunks, opts)
	if !res.Success {
		t.Fatalf("CompileChunks failed: %v", res.Errors)
	}

	// Sequential reference rendering via ToSourceOfLicenseTrackerAndJSChunk + SourceMap.AppendTo.
	seqLT := jscomp.NewCompiler_ChunkGraphAwareLicenseTracker(comp)
	wantCode := make([]string, len(chunks))
	wantMaps := make([]string, len(chunks))
	for i, ch := range chunks {
		comp.ResetAndInitializeSourceMap()
		seqLT.SetCurrentChunkContext(ch)
		wantCode[i] = comp.ToSourceOfLicenseTrackerAndJSChunk(seqLT, ch)
		var sb strings.Builder
		comp.GetSourceMap().AppendTo(&sb, ch.GetName()+".js")
		wantMaps[i] = sb.String()
	}

	rendered, err := comp.RenderChunksAndSourceMaps(chunks, jscomp.RenderChunkOptions{
		GenerateSourceMaps:  true,
		OmitSourceMapString: true,
		CompressSourceMapBytes: func(_ *jscomp.JSChunk, sm []byte) ([]byte, error) {
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			if _, wErr := gz.Write(sm); wErr != nil {
				return nil, wErr
			}
			if cErr := gz.Close(); cErr != nil {
				return nil, cErr
			}
			return buf.Bytes(), nil
		},
	})
	if err != nil {
		t.Fatalf("RenderChunksAndSourceMaps error: %v", err)
	}
	if len(rendered) != len(chunks) {
		t.Fatalf("len(rendered) = %d, want %d", len(rendered), len(chunks))
	}
	for i, rc := range rendered {
		if rc.Code != wantCode[i] {
			t.Fatalf("chunk %d (%s) JS mismatch:\ngot len=%d\nwant len=%d", i, chunks[i].GetName(), len(rc.Code), len(wantCode[i]))
		}
		if rc.SourceMap != "" {
			t.Fatalf("chunk %d (%s) expected empty SourceMap when OmitSourceMapString=true, got %d bytes", i, chunks[i].GetName(), len(rc.SourceMap))
		}
		zr, zErr := gzip.NewReader(bytes.NewReader(rc.SourceMapBytes))
		if zErr != nil {
			t.Fatal(zErr)
		}
		raw, rErr := io.ReadAll(zr)
		_ = zr.Close()
		if rErr != nil {
			t.Fatal(rErr)
		}
		if string(raw) != wantMaps[i] {
			t.Fatalf("chunk %d (%s) source map mismatch against sequential reference", i, chunks[i].GetName())
		}
	}
}

type simpleRenamingMap struct{}

func (simpleRenamingMap) Get(value string) (string, bool) { return "id_" + value, true }

func TestOptionsStageDecoratorsExtraOptionsWarningsAndStateMetadata(t *testing.T) {
	opts := jscomp.NewCompilerOptions()
	jscomp.CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)

	jscomp.SetExtraOption(opts, "customKey", 42)
	if v, ok := jscomp.ExtraOption(opts, "customKey"); !ok || v != 42 {
		t.Fatalf("ExtraOption(customKey) = (%v, %v), want (42, true)", v, ok)
	}
	if _, ok := jscomp.ExtraOption(opts, "missing"); ok {
		t.Fatal("expected missing key to return false")
	}
	snap := jscomp.ExtraOptionsSnapshot(opts)
	snap["customKey"] = 99
	if v, _ := jscomp.ExtraOption(opts, "customKey"); v != 42 {
		t.Fatalf("ExtraOptionsSnapshot did not defensively copy map: got %v", v)
	}

	jscomp.AddOptionWarning(opts, "deprecated custom flag --foo")

	var stages []string
	jscomp.SetChecksDecorator(opts, func(pc *jscomp.DefaultPassConfig, b *jscomp.PassListBuilder) {
		if pc == nil || b == nil {
			t.Fatal("expected non-nil DefaultPassConfig and PassListBuilder")
		}
		b.InsertAt(b.Len(), jscomp.NewPassFactory("optChecksHook", func(jscomp.AbstractCompiler) jscomp.CompilerPass {
			return jscomp.CompilerPassFunc(func(_, _ rhino.Node) {
				stages = append(stages, "checks")
			})
		}))
	})
	jscomp.SetOptimizationsDecorator(opts, func(pc *jscomp.DefaultPassConfig, b *jscomp.PassListBuilder) {
		b.InsertAt(b.Len(), jscomp.NewPassFactory("optOptHook", func(jscomp.AbstractCompiler) jscomp.CompilerPass {
			return jscomp.CompilerPassFunc(func(_, _ rhino.Node) {
				stages = append(stages, "optimizations")
			})
		}))
	})
	jscomp.SetFinalizationsDecorator(opts, func(pc *jscomp.DefaultPassConfig, b *jscomp.PassListBuilder) {
		b.InsertAt(b.Len(), jscomp.NewPassFactory("optFinalHook", func(jscomp.AbstractCompiler) jscomp.CompilerPass {
			return jscomp.CompilerPassFunc(func(_, _ rhino.Node) {
				stages = append(stages, "finalizations")
			})
		}))
	})

	c := jscomp.NewCompiler()
	res := c.Compile(
		[]*jscomp.SourceFile{jscomp.SourceFile_FromCode("ext.js", "/** @externs */ var window;")},
		[]*jscomp.SourceFile{jscomp.SourceFile_FromCode("in.js", "window['x'] = 1;\n")},
		opts,
	)
	if !res.Success {
		t.Fatalf("compile failed: %v", res.Errors)
	}
	if strings.Join(stages, ",") != "checks,optimizations,finalizations" {
		t.Fatalf("stages = %v, want [checks optimizations finalizations]", stages)
	}
	foundOptWarn := false
	for _, w := range res.Warnings {
		if w.Type() == jscomp.StrictWarningsGuard_UNRAISABLE_WARNING && strings.Contains(w.Description(), "deprecated custom flag --foo") {
			foundOptWarn = true
		}
	}
	if !foundOptWarn {
		t.Fatalf("expected UNRAISABLE_WARNING for AddOptionWarning, got %v", res.Warnings)
	}

	// Verify Compiler.Parse() also emits queued option warnings.
	parseOpts := jscomp.NewCompilerOptions()
	jscomp.AddOptionWarning(parseOpts, "parse-only warning")
	cp := jscomp.NewCompiler()
	cp.Init([]*jscomp.SourceFile{}, []*jscomp.SourceFile{jscomp.SourceFile_FromCode("p.js", "var p = 1;\n")}, parseOpts)
	cp.Parse()
	if len(cp.GetWarnings()) == 0 || !strings.Contains(cp.GetWarnings()[0].Description(), "parse-only warning") {
		t.Fatalf("expected Parse() to emit option warning, got %v", cp.GetWarnings())
	}

	// Verify GetCompilerStateProtoBuilder rejects non-empty CustomNodeMetadata and RestoreFromState clears it.
	c.CustomNodeMetadata().Put(c.GetRoot(), "non-serializable")
	panicked := false
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				if _, ok := rec.(jutil.IIllegalStateException); ok {
					panicked = true
					return
				}
				panic(rec)
			}
		}()
		_ = c.GetCompilerStateProtoBuilder()
	}()
	if !panicked {
		t.Fatal("expected GetCompilerStateProtoBuilder to panic when CustomNodeMetadata is non-empty")
	}
	c.CustomNodeMetadata().Clear()
	stateProto := c.GetCompilerStateProtoBuilder().Build()
	c.CustomNodeMetadata().Put(c.GetRoot(), "stale")
	c.RestoreFromState(stateProto)
	if !c.CustomNodeMetadata().IsEmpty() {
		t.Fatal("expected RestoreFromState to reset CustomNodeMetadata")
	}
}

func TestExportedOptionGettersSettersHashesAndAssertionBuilder(t *testing.T) {
	opts := jscomp.NewCompilerOptions()
	if got := jscomp.LanguageOutIsDefaultStrict(opts); got != base.Tri_UNKNOWN {
		t.Fatalf("initial LanguageOutIsDefaultStrict = %v, want UNKNOWN", got)
	}
	opts.SetLanguageOut(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	if got := jscomp.LanguageOutIsDefaultStrict(opts); got != base.Tri_TRUE {
		t.Fatalf("ES2015 LanguageOutIsDefaultStrict = %v, want TRUE", got)
	}
	opts.SetLanguageOut(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT3)
	if got := jscomp.LanguageOutIsDefaultStrict(opts); got != base.Tri_FALSE {
		t.Fatalf("ES3 LanguageOutIsDefaultStrict = %v, want FALSE", got)
	}
	opts.SetMaxOptimizationLoopIterations(7)
	if got := jscomp.OptimizationLoopMaxIterations(opts); got != 7 {
		t.Fatalf("OptimizationLoopMaxIterations = %d, want 7", got)
	}
	opts.SetPreserveClosurePrimitives(true)
	if !jscomp.PreserveClosurePrimitives(opts) {
		t.Fatal("PreserveClosurePrimitives = false, want true")
	}
	opts.SetEnableModuleRewriting(false)
	if jscomp.EnableModuleRewriting(opts) {
		t.Fatal("EnableModuleRewriting = true, want false")
	}
	opts.SetReplaceMessagesWithChromeI18n(true, base.SomeString("proj_123"))
	enabled, proj := jscomp.ReplaceMessagesWithChromeI18n(opts)
	if projStr, ok := proj.Get(); !enabled || !ok || projStr != "proj_123" {
		t.Fatalf("ReplaceMessagesWithChromeI18n = (%v, %v)", enabled, proj)
	}

	cfg := jscomp.ConformanceConfig_NewBuilder().Build()
	opts.SetConformanceConfigs([]*jscomp.ConformanceConfig{cfg})
	opts.ClearConformanceConfigs()
	if len(opts.GetConformanceConfigs()) != 0 {
		t.Fatal("expected ClearConformanceConfigs to clear configs")
	}

	opts.SetExportTestFunctions(true)
	if !opts.ShouldExportTestFunctions() {
		t.Fatal("ShouldExportTestFunctions = false")
	}
	idMap := jscomp.NewOrderedStringMap[jscomp.RenamingMap]()
	idMap.Put("myIdGen", simpleRenamingMap{})
	opts.SetIdGeneratorsOfMap(idMap)
	if opts.GetIdGenerators().Size() != 1 {
		t.Fatalf("GetIdGenerators().Size() = %d, want 1", opts.GetIdGenerators().Size())
	}
	opts.SetInferConst(false)
	if opts.ShouldInferConsts() {
		t.Fatal("ShouldInferConsts = true, want false")
	}
	opts.SetJ2clMinifierEnabled(false)
	if opts.IsJ2clMinifierEnabled() {
		t.Fatal("IsJ2clMinifierEnabled = true, want false")
	}
	opts.SetOutputJsStringUsage(true)
	if !opts.ShouldOutputJsStringUsage() {
		t.Fatal("ShouldOutputJsStringUsage = false")
	}
	opts.SetParentChunkCanSeeSymbolsDeclaredInChildren(true)
	if !opts.GetParentChunkCanSeeSymbolsDeclaredInChildren() {
		t.Fatal("GetParentChunkCanSeeSymbolsDeclaredInChildren = false")
	}
	opts.SetPrintExterns(true)
	if !opts.ShouldPrintExterns() {
		t.Fatal("ShouldPrintExterns = false")
	}
	opts.SetPropertyRenamingOnlyCompilationMode(true)
	if !opts.IsPropertyRenamingOnlyCompilationMode() {
		t.Fatal("IsPropertyRenamingOnlyCompilationMode = false")
	}
	opts.SetReplaceStringsFunctionDescriptions([]string{"Error(?)"})
	opts.SetReplaceStringsPlaceholderToken("`")
	if len(opts.GetReplaceStringsFunctionDescriptions()) != 1 || opts.GetReplaceStringsPlaceholderToken() != "`" {
		t.Fatal("ReplaceStrings setters mismatch")
	}
	opts.SetSerializeExtraDebugInfo(true)
	if !opts.ShouldSerializeExtraDebugInfo() {
		t.Fatal("ShouldSerializeExtraDebugInfo = false")
	}
	opts.SetStrictMessageReplacement(true)
	if !opts.GetStrictMessageReplacement() {
		t.Fatal("GetStrictMessageReplacement = false")
	}

	// Deterministic Jenkins lookup2 32-bit hash sanity check:
	h1 := jscomp.JenkinsHash32WithSeed([]byte("hello.world"), 0)
	h2 := jscomp.JenkinsHash32WithSeed([]byte("hello.world"), 3141592)
	if h1 == 0 || h2 == 0 || h1 == h2 {
		t.Fatalf("unexpected JenkinsHash32WithSeed values: h1=%d h2=%d", h1, h2)
	}

	// AssertionFunctionSpec builder & accessors:
	spec1 := jscomp.CodingConvention_AssertionFunctionSpec_ForTruthy().
		SetFunctionName("my.asserts.check").
		SetParamIndex(1).
		Build()
	spec2 := jscomp.CodingConvention_AssertionFunctionSpec_ForTruthy().
		SetFunctionName("my.asserts.check").
		SetParamIndex(1).
		Build()
	fnName, hasFnName := spec1.GetFunctionName().Get()
	if !hasFnName || fnName != "my.asserts.check" || spec1.GetParamIndex() != 1 || !spec1.Equals(spec2) {
		t.Fatalf("unexpected AssertionFunctionSpec: %v", spec1)
	}
}

func TestConformanceBannedScriptImportCssNamesAndLibraryBehavior(t *testing.T) {
	if jscomp.Requirement_Type_BANNED_SCRIPT_IMPORT.GetNumber() != 19 {
		t.Fatalf("BANNED_SCRIPT_IMPORT number = %d, want 19", jscomp.Requirement_Type_BANNED_SCRIPT_IMPORT.GetNumber())
	}
	reqBuilder := jscomp.Requirement_NewBuilder()
	if err := jscomp.TextFormat_MergeErr(`
type: BANNED_SCRIPT_IMPORT
value: 'goog.net.jsloader.safeLoad'
rule_id: 'ban_script_import'
error_message: 'Script import forbidden'
`, reqBuilder); err != nil {
		t.Fatalf("TextFormat_MergeErr failed: %v", err)
	}
	cfg := jscomp.ConformanceConfig_NewBuilder().
		AddRequirement(reqBuilder.Build()).
		SetLibraryLevelNonAllowlistedConformanceViolationsBehavior(
			jscomp.ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_RECORD_ONLY).
		Build()
	roundTripCfg, err := jscomp.ParseConformanceConfigBytes(cfg.ToByteArray())
	if err != nil {
		t.Fatalf("ParseConformanceConfigBytes: %v", err)
	}
	if roundTripCfg.GetRequirement(0).GetType() != jscomp.Requirement_Type_BANNED_SCRIPT_IMPORT {
		t.Fatalf("round-tripped type = %v", roundTripCfg.GetRequirement(0).GetType())
	}

	var customEvalCalls int
	jscomp.CustomBannedNameEvaluator = func(reqType jscomp.Requirement_Type, n rhino.Node) (*jscomp.ConformanceRules_ConformanceResult, bool) {
		if reqType == jscomp.Requirement_Type_BANNED_SCRIPT_IMPORT {
			customEvalCalls++
			// Suppress safe call when first arg is literal 'allowlisted.js':
			parent := n.GetParent()
			if !parent.IsNil() && parent.IsCall() && parent.GetFirstChild() == n {
				arg := n.GetNext()
				if !arg.IsNil() && arg.IsStringLit() && arg.GetString() == "allowlisted.js" {
					return jscomp.ConformanceRules_ConformanceResult_CONFORMANCE, true
				}
			}
		}
		return nil, false
	}
	defer func() { jscomp.CustomBannedNameEvaluator = nil }()

	opts := jscomp.NewCompilerOptions()
	opts.SetClosurePass(true)
	opts.SetGatherCssNames(true)
	opts.SetConformanceReportingMode(
		jscomp.CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG)
	opts.SetConformanceConfigs([]*jscomp.ConformanceConfig{roundTripCfg})

	var recordedViolations int
	c := jscomp.NewCompiler()
	c.SetConformanceViolationInterceptor(func(
		r *jscomp.Requirement,
		_ *jscomp.RequirementScopeEntry,
		_ *jscomp.JSError,
		behavior jscomp.ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior,
		_ bool,
	) bool {
		recordedViolations++
		if behavior != jscomp.ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_RECORD_ONLY {
			t.Fatalf("unexpected interceptor behavior: %v", behavior)
		}
		if got, ok := c.LookupLibraryLevelConformanceBehavior(r); !ok || got != behavior {
			t.Fatalf("LookupLibraryLevelConformanceBehavior = (%v, %v), want (%v, true)", got, ok, behavior)
		}
		return false
	})

	res := c.Compile(
		[]*jscomp.SourceFile{jscomp.SourceFile_FromCode("ext.js", `/** @externs */
var goog = {};
goog.getCssName = function(s) {};
goog.net = {jsloader: {safeLoad: function(u) {}}};
`)},
		[]*jscomp.SourceFile{jscomp.SourceFile_FromCode("app.js", `
var cls = goog.getCssName('btn-primary');
goog.net.jsloader.safeLoad('allowlisted.js');
goog.net.jsloader.safeLoad('untrusted.js');
`)},
		opts,
	)
	if !res.Success {
		t.Fatalf("compile failed: %v", res.Errors)
	}
	if customEvalCalls != 2 || recordedViolations != 1 {
		t.Fatalf("customEvalCalls=%d recordedViolations=%d, want 2 and 1", customEvalCalls, recordedViolations)
	}
	if c.CssNames() == nil || !c.CssNames().Has("btn-primary") {
		t.Fatalf("expected CssNames() to contain 'btn-primary', got %v", c.CssNames())
	}

	// Explicit Record/Lookup test:
	manualReq := jscomp.Requirement_NewBuilder().SetRuleId("manual_rule").Build()
	c.RecordLibraryLevelConformanceBehavior(
		manualReq,
		jscomp.ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_REPORT_AS_BUILD_ERROR)
	if b, ok := c.LookupLibraryLevelConformanceBehavior(manualReq); !ok ||
		b != jscomp.ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_REPORT_AS_BUILD_ERROR {
		t.Fatalf("manual LookupLibraryLevelConformanceBehavior = (%v, %v)", b, ok)
	}
}
