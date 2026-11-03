package jscomp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func buildTypedAstTestCorpus() (externCode string, scriptCodes []string) {
	externCode = NewTestExternsBuilder().AddClosureExterns().AddString().Build() + `
/** @type {!Object} */ var window;
/** @param {...*} args */ function extSink(args) {}
/**
 * @param {!ITemplateArray} strings
 * @param {...*} values
 * @return {string}
 */
function rawTag(strings, values) {}
`
	const numScripts = 16
	scriptCodes = make([]string, numScripts)
	for i := range numScripts {
		if i%4 == 3 {
			scriptCodes[i] = fmt.Sprintf(`/**
 * @fileoverview
 * @closureUnaware
 * @license Shadow Script License %[1]d
 */
var shadowHolder_%[1]d = 0;
/** @closureUnaware */
(function() {
  let acc = %[1]d;
  acc++;
  const obj = {'quoted_%[1]d': acc, plain_%[1]d: acc + 1};
  shadowHolder_%[1]d = obj['quoted_%[1]d'] + obj.plain_%[1]d;
})();
extSink(shadowHolder_%[1]d);
`, i)
			continue
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "/** @license Script %d License */\n", i)
		fmt.Fprintf(&sb, "/** @constructor */\nfunction Widget_%d(seed) { this.seed = seed; }\n", i)
		fmt.Fprintf(&sb, "/** @noinline @param {number} n @return {string} */\n")
		fmt.Fprintf(&sb, "Widget_%d.prototype.compute = function(n) {\n", i)
		fmt.Fprintf(&sb, "  let count = n + this.seed;\n")
		fmt.Fprintf(&sb, "  count++;\n")
		fmt.Fprintf(&sb, "  count--;\n")
		fmt.Fprintf(&sb, "  const big = %dn + 9007199254740993n;\n", i+1)
		fmt.Fprintf(&sb, "  const dict = {'quotedKey_%d': count, normalKey_%d: -0.0 + %d.5};\n", i, i, i)
		for k := range 12 {
			fmt.Fprintf(&sb, "  count += (%d * %d) ^ (count & %d);\n", k+1, i+3, k+7)
		}
		fmt.Fprintf(&sb, "  const msg = `item_${count}_${dict['quotedKey_%d']}_\\n_${String(big)}`;\n", i)
		fmt.Fprintf(&sb, "  const tagged = rawTag`raw\\unicode_%d_${msg}`;\n", i)
		fmt.Fprintf(&sb, "  return tagged;\n")
		fmt.Fprintf(&sb, "};\n")
		fmt.Fprintf(&sb, "window['w_%d'] = new Widget_%d(%d).compute(%d);\n", i, i, i*2, i+5)
		scriptCodes[i] = sb.String()
	}
	return externCode, scriptCodes
}

func makeTypedAstSourceFiles(externCode string, scriptCodes []string) ([]*SourceFile, []*SourceFile) {
	externs := []*SourceFile{
		SourceFile_FromCode("externs.js", externCode),
	}
	sources := make([]*SourceFile, len(scriptCodes))
	for i, code := range scriptCodes {
		sources[i] = SourceFile_FromCode(fmt.Sprintf("script_%02d.js", i), code)
	}
	return externs, sources
}

func serializeTestTypedAstList(t *testing.T, externCode string, scriptCodes []string) []byte {
	t.Helper()
	c := NewCompiler()
	opts := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)

	externs, sources := makeTypedAstSourceFiles(externCode, scriptCodes)
	c.Init(externs, sources, opts)
	c.ParseForCompilation()
	if c.HasErrors() {
		t.Fatalf("parse errors: %v", c.GetErrors())
	}
	c.Stage1Passes()
	if c.HasErrors() {
		t.Fatalf("stage1 errors: %v", c.GetErrors())
	}

	serializer := NewTypedAstSerializer(
		c,
		NewSerializationOptions().SetIncludeDebugInfo(true).Build(),
	)
	typedAst := serializer.SerializeRoots(c.GetExternsRoot(), c.GetJsRoot())
	return (&TypedAst_List{typedAsts_: []*TypedAst{typedAst}}).ToByteArray()
}

func TestScriptNodeDeserializer_DeserializeNewWithAllocatorParityAndContiguity(t *testing.T) {
	externCode, scriptCodes := buildTypedAstTestCorpus()
	rawList := serializeTestTypedAstList(t, externCode, scriptCodes)

	deserWithAlloc := func(useAlloc bool) (*Compiler, []rhino.Node) {
		c := NewCompiler()
		opts := NewCompilerOptions()
		CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
		opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
		opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
		externs, sources := makeTypedAstSourceFiles(externCode, scriptCodes)
		c.Init(externs, sources, opts)

		required := make([]*SourceFile, 0, len(externs)+len(sources))
		required = append(required, externs...)
		required = append(required, sources...)
		ast := TypedAstDeserializer_DeserializeFullAst(
			c,
			SourceFile_FromCodeWithKind("syntheticExterns", "", rhino.StaticSourceFile_SourceKind_EXTERN),
			required,
			bytes.NewReader(rawList),
			true,
			true,
			true,
		)

		var alloc *rhino.NodeAllocator
		if useAlloc {
			alloc = rhino.NewNodeAllocator()
		}
		roots := make([]rhino.Node, len(sources))
		for i, sf := range sources {
			entry, ok := ast.filesystem[sf]
			if !ok || entry.sd == nil {
				t.Fatalf("missing ScriptNodeDeserializer for %s", sf.GetName())
			}
			if useAlloc {
				roots[i] = entry.sd.DeserializeNewWithAllocator(alloc)
			} else {
				roots[i] = entry.sd.DeserializeNew()
			}
		}
		return c, roots
	}

	_, rootsDefault := deserWithAlloc(false)
	_, rootsAlloc := deserWithAlloc(true)

	for i := range rootsDefault {
		if !rootsAlloc[i].IsEquivalentTo(rootsDefault[i]) {
			t.Fatalf("script %d AST mismatch between DeserializeNew and DeserializeNewWithAllocator:\ngot:  %s\nwant: %s",
				i, rootsAlloc[i].ToStringTree(), rootsDefault[i].ToStringTree())
		}
		fsGot, _ := rootsAlloc[i].GetProp(rhino.Node_FEATURE_SET).(parser.FeatureSet)
		fsWant, _ := rootsDefault[i].GetProp(rhino.Node_FEATURE_SET).(parser.FeatureSet)
		if fsGot != fsWant {
			t.Fatalf("script %d FeatureSet mismatch: got %v, want %v", i, fsGot, fsWant)
		}

		// Verify all nodes in a single script allocated via a fresh NodeAllocator
		// come from a contiguous slab (difference between consecutive preorder IDs
		// is bounded, with no global-counter interleaving).
		var ids []uint32
		numShadowRoots := 0
		var collect func(rhino.Node)
		collect = func(n rhino.Node) {
			if n.IsNil() {
				return
			}
			ids = append(ids, uint32(n))
			if shadow := n.GetClosureUnawareShadow(); !shadow.IsNil() {
				numShadowRoots++
				collect(shadow)
			}
			for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
				collect(ch)
			}
		}
		collect(rootsAlloc[i])
		if len(ids) < 5 {
			t.Fatalf("script %d unexpectedly small (%d nodes)", i, len(ids))
		}
		minID, maxID := ids[0], ids[0]
		for _, id := range ids[1:] {
			if id < minID {
				minID = id
			}
			if id > maxID {
				maxID = id
			}
		}
		// Each shadow root allocates 1 local closure-unaware template node from alloc.
		wantSpan := len(ids) + numShadowRoots
		if span := int(maxID - minID + 1); span != wantSpan {
			t.Fatalf("script %d expected contiguous node IDs (wantSpan=%d), got min=%d max=%d span=%d",
				i, wantSpan, minID, maxID, span)
		}
	}
}

func TestPrebuildAst_ParallelTypedAstDeserializerParity(t *testing.T) {
	externCode, scriptCodes := buildTypedAstTestCorpus()
	rawList := serializeTestTypedAstList(t, externCode, scriptCodes)

	type runSnapshot struct {
		prebuiltTree string
		featureSets  []parser.FeatureSet
		finalJS      string
		finalTree    string
	}

	runWithThreads := func(threads int) runSnapshot {
		c := NewCompiler()
		opts := NewCompilerOptions()
		CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
		opts.SetDependencyOptions(DependencyOptions_None())
		opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
		opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
		opts.SetNumParallelThreads(threads)

		externs, sources := makeTypedAstSourceFiles(externCode, scriptCodes)
		chunks := []*JSChunk{
			NewJSChunk("base"),
			NewJSChunk("feat1"),
			NewJSChunk("feat2"),
			NewJSChunk("feat3"),
		}
		for i := 1; i < len(chunks); i++ {
			chunks[i].AddDependency(chunks[0])
		}
		for i, sf := range sources {
			chunks[i%len(chunks)].Add(NewCompilerInput(sf))
		}

		c.InitChunksWithTypedAstFilesystem(externs, chunks, opts, bytes.NewReader(rawList))
		if c.HasErrors() {
			t.Fatalf("threads=%d InitChunksWithTypedAstFilesystem errors: %v", threads, c.GetErrors())
		}

		prebuiltTree := c.GetJsRoot().ToStringTree()
		var fsets []parser.FeatureSet
		for s := c.GetJsRoot().GetFirstChild(); !s.IsNil(); s = s.GetNext() {
			fs, _ := s.GetProp(rhino.Node_FEATURE_SET).(parser.FeatureSet)
			fsets = append(fsets, fs)
		}

		c.Stage2Passes(CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS)
		if c.HasErrors() {
			t.Fatalf("threads=%d Stage2Passes errors: %v", threads, c.GetErrors())
		}
		c.Stage3Passes()
		if c.HasErrors() {
			t.Fatalf("threads=%d Stage3Passes errors: %v", threads, c.GetErrors())
		}

		return runSnapshot{
			prebuiltTree: prebuiltTree,
			featureSets:  fsets,
			finalJS:      c.ToSource(),
			finalTree:    c.GetJsRoot().ToStringTree(),
		}
	}

	seq := runWithThreads(1)
	for iter := range 3 {
		par := runWithThreads(4)
		if par.prebuiltTree != seq.prebuiltTree {
			t.Fatalf("iter=%d prebuilt AST tree mismatch (threads=1 vs 4)", iter)
		}
		if len(par.featureSets) != len(seq.featureSets) {
			t.Fatalf("iter=%d featureSets len=%d, want %d", iter, len(par.featureSets), len(seq.featureSets))
		}
		for i := range seq.featureSets {
			if par.featureSets[i] != seq.featureSets[i] {
				t.Fatalf("iter=%d script %d FeatureSet=%v, want %v", iter, i, par.featureSets[i], seq.featureSets[i])
			}
		}
		if par.finalTree != seq.finalTree {
			t.Fatalf("iter=%d final AST tree mismatch (threads=1 vs 4)", iter)
		}
		if par.finalJS != seq.finalJS {
			t.Fatalf("iter=%d final JS mismatch:\ngot (threads=4):  %q\nwant (threads=1): %q", iter, par.finalJS, seq.finalJS)
		}
	}
}
