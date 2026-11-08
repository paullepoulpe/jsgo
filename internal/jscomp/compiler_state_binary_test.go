/*
 * Copyright 2026 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type nonCloseTrackingStream struct {
	buf         bytes.Buffer
	writeClosed bool
	readClosed  bool
	reader      *bytes.Reader
}

func (s *nonCloseTrackingStream) Write(p []byte) (int, error) { return s.buf.Write(p) }
func (s *nonCloseTrackingStream) Close() error {
	if s.reader != nil {
		s.readClosed = true
	} else {
		s.writeClosed = true
	}
	return nil
}
func (s *nonCloseTrackingStream) Read(p []byte) (int, error) { return s.reader.Read(p) }

type failingShortReader struct {
	data []byte
	err  error
}

func (r *failingShortReader) Read(p []byte) (int, error) {
	if len(r.data) > 0 {
		n := copy(p, r.data)
		r.data = r.data[n:]
		return n, nil
	}
	return 0, r.err
}

func buildBinaryStateTestCorpus(numScripts int, extraStrings int) (externCode string, scriptCodes []string) {
	loneSurrogate := jutil.StringFromUTF16([]uint16{0xD800, 's', 'u', 'r', 'r'})
	externCode = NewTestExternsBuilder().AddClosureExterns().AddString().AddPromise().Build() + `
/** @type {!Object} */ var window;
/** @param {...*} args */ function extSink(args) {}
/**
 * @param {!ITemplateArray} strings
 * @param {...*} values
 * @return {string}
 */
function rawTag(strings, values) {}
`
	scriptCodes = make([]string, numScripts)
	for i := range numScripts {
		origPath := fmt.Sprintf("src/orig_%02d.ts", i)
		jsPath := fmt.Sprintf("script_%02d.js", i)
		inlineMapJSON := fmt.Sprintf(
			`{"version":3,"file":%q,"sources":[%q],"sourcesContent":[%q],"names":["origSym_%d"],"mappings":"AAAA;AACA,SAASA,GAAG,CAAC,IAAI;AACjB,SAAO,IAAI,GAAG,CAAC;AACjB"}`,
			jsPath, origPath,
			fmt.Sprintf("// Café résumé λ %d\nexport function fn_%d(x: number) { return x + %d; }\n", i, i, i),
			i,
		)
		inlineMapB64 := base64.StdEncoding.EncodeToString([]byte(inlineMapJSON))

		if i%4 == 3 {
			scriptCodes[i] = fmt.Sprintf(`/**
 * @fileoverview
 * @closureUnaware
 * @license Shadow License %[1]d
 */
var shadowOut_%[1]d = 0;
/** @closureUnaware */
(function(seed) {
  let acc = seed + 1;
  acc++;
  const d = {'quoted_%[1]d': acc, plain_%[1]d: acc * 2};
  shadowOut_%[1]d = d['quoted_%[1]d'] + d.plain_%[1]d;
})(%[1]d);
extSink(shadowOut_%[1]d);
//# sourceMappingURL=data:application/json;charset=utf-8;base64,%[2]s
`, i, inlineMapB64)
			continue
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "/** @license Module %d License */\n", i)
		fmt.Fprintf(&sb, "/** @desc Greeting %d @meaning box_%d */\nvar MSG_ITEM_%d = 'hello_%d_%s';\n", i, i, i, i, loneSurrogate)
		fmt.Fprintf(&sb, "class Box_%d {\n", i)
		fmt.Fprintf(&sb, "  constructor(seed) { /** @type {number} */ this.seed = seed; }\n")
		fmt.Fprintf(&sb, "  get value() { return this.seed; }\n")
		fmt.Fprintf(&sb, "  set value(v) { this.seed = v; }\n")
		fmt.Fprintf(&sb, "  /** @noinline @param {number} n @return {string} */\n")
		fmt.Fprintf(&sb, "  compute(n) {\n")
		fmt.Fprintf(&sb, "    let total = n + this.value;\n")
		fmt.Fprintf(&sb, "    ++total;\n")
		fmt.Fprintf(&sb, "    total++;\n")
		fmt.Fprintf(&sb, "    total--;\n")
		fmt.Fprintf(&sb, "    total ||= 1;\n")
		fmt.Fprintf(&sb, "    const big = %dn + 9007199254740993n;\n", i+1)
		fmt.Fprintf(&sb, "    const negZero = -0.0;\n")
		fmt.Fprintf(&sb, "    const nanVal = 0 / 0;\n")
		fmt.Fprintf(&sb, "    const bag = {'qKey_%d': total, plainKey_%d: negZero + %d.25, nanKey_%d: nanVal};\n", i, i, i, i)
		for k := range 8 {
			fmt.Fprintf(&sb, "    total += (%d * %d) ^ (total & %d);\n", k+1, i+2, k+5)
		}
		if i == 1 && extraStrings > 0 {
			fmt.Fprintf(&sb, "    const extraTable = [")
			for s := range extraStrings {
				if s > 0 {
					sb.WriteByte(',')
				}
				fmt.Fprintf(&sb, "'uniq_str_token_%05d'", s)
			}
			fmt.Fprintf(&sb, "];\n    extSink(extraTable.length);\n")
		}
		fmt.Fprintf(&sb, "    const cooked = `val_${total}_${bag['qKey_%d']}_${MSG_ITEM_%d}_\\n_${String(big)}`;\n", i, i)
		fmt.Fprintf(&sb, "    const rawOnly = rawTag`bad\\unicode_%d_${cooked}`;\n", i)
		fmt.Fprintf(&sb, "    return rawOnly;\n")
		fmt.Fprintf(&sb, "  }\n}\n")
		fmt.Fprintf(&sb, "async function asyncEntry_%d(x) { return await Promise.resolve(new Box_%d(x).compute(%d)); }\n", i, i, i+3)
		fmt.Fprintf(&sb, "hoistedExternHook(asyncEntry_%d(%d));\n", i, i*2)
		fmt.Fprintf(&sb, "window['exp_%d'] = new Box_%d(%d).compute(%d);\n", i, i, i, i+1)
		fmt.Fprintf(&sb, "//# sourceMappingURL=data:application/json;charset=utf-8;base64,%s\n", inlineMapB64)
		scriptCodes[i] = sb.String()
	}
	return externCode, scriptCodes
}

func makeBinaryStateChunks(externCode string, scriptCodes []string) ([]*SourceFile, []*JSChunk) {
	externs := []*SourceFile{
		SourceFile_FromCode("externs.js", externCode),
		SourceFile_FromCode("extra_externs.js", "/** @externs */\n/** @param {*} v */ function hoistedExternHook(v) {}\n"),
	}
	chunks := []*JSChunk{
		NewJSChunk("base"),
		NewJSChunk("feat1"),
		NewJSChunk("feat2"),
		NewJSChunk("feat3"),
	}
	for i := 1; i < len(chunks); i++ {
		chunks[i].AddDependency(chunks[0])
	}
	for i, code := range scriptCodes {
		sf := SourceFile_FromCode(fmt.Sprintf("script_%02d.js", i), code)
		if i == 0 {
			chunks[0].Add(NewCompilerInput(sf))
		} else {
			chunks[i%len(chunks)].Add(NewCompilerInput(sf))
		}
	}
	return externs, chunks
}

func newBinaryStateTestOptions(binaryMode bool, threads int) *CompilerOptions {
	opts := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	opts.SetRewritePolyfills(true)
	opts.SetCheckTypes(true)
	opts.SetParseInlineSourceMaps(true)
	opts.SetApplyInputSourceMaps(true)
	opts.SetSourceMapOutputPath("%outname%.map")
	opts.SetSourceMapIncludeSourcesContent(true)
	opts.SetNumParallelThreads(threads)
	opts.SetBinarySaveState(binaryMode)
	opts.serializeExtraDebugInfo = true
	return opts
}

func fmtOptValStr(s string, ok bool) string {
	if !ok {
		return "<nil>"
	}
	return fmt.Sprintf("%q", s)
}

func snapshotDeepAst(root rhino.Node) string {
	var sb strings.Builder
	sb.WriteString(root.ToStringTree())
	sb.WriteString("\n---DEEP---\n")
	var walk func(depth int, n rhino.Node)
	walk = func(depth int, n rhino.Node) {
		if n.IsNil() {
			return
		}
		colorID := "<nil>"
		if col := n.GetColor(); col != nil {
			colorID = col.GetId().String()
		}
		fmt.Fprintf(&sb, "%d:%v@%s:%d:%d props=0x%x orig=%s col=%s",
			depth, n.GetToken(), fmtOptValStr(n.SourceFileName()),
			n.GetLineno(), n.GetCharno(), n.SerializeProperties(),
			fmtOptValStr(n.OriginalName()), colorID)
		if n.IsTemplateLitString() {
			fmt.Fprintf(&sb, " tmplRaw=%q tmplCooked=%s", n.GetRawString(), fmtOptValStr(n.CookedString()))
		}
		if n.IsScript() {
			fs, _ := n.GetProp(rhino.Node_FEATURE_SET).(parser.FeatureSet)
			fmt.Fprintf(&sb, " features=%v", fs)
		}
		if j := n.GetJSDocInfo(); j != nil {
			fmt.Fprintf(&sb, " jsdoc(lic=%s,desc=%s,noinl=%v,nose=%v,cu=%v)",
				fmtOptValStr(j.License()), fmtOptValStr(j.Description()),
				j.IsNoInline(), j.IsNoSideEffects(), j.IsClosureUnawareCode())
		}
		sb.WriteByte('\n')
		if shadow := n.GetClosureUnawareShadow(); !shadow.IsNil() {
			sb.WriteString("#shadow:\n")
			walk(depth+1, shadow)
		}
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			walk(depth+1, ch)
		}
	}
	walk(0, root)
	return sb.String()
}

func gunzipPayload(t *testing.T, data []byte) []byte {
	t.Helper()
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("gzip.NewReader failed: %v", err)
	}
	defer zr.Close()
	raw, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("gunzip ReadAll failed: %v", err)
	}
	return raw
}

func TestBinarySaveState_DefaultModeUnchangedAndOptInTriggers(t *testing.T) {
	externCode, scriptCodes := buildBinaryStateTestCorpus(4, 0)

	saveWithConfig := func(optFlag bool, envVal string) ([]byte, *nonCloseTrackingStream) {
		if envVal != "" {
			t.Setenv("JSGO_BINARY_SAVE_STATE", envVal)
		} else {
			t.Setenv("JSGO_BINARY_SAVE_STATE", "")
		}
		c := NewCompiler()
		opts := newBinaryStateTestOptions(optFlag, 1)
		if CompilerOptions_IsBinarySaveState(opts) != optFlag {
			t.Fatalf("CompilerOptions_IsBinarySaveState() = %v, want %v", CompilerOptions_IsBinarySaveState(opts), optFlag)
		}
		externs, chunks := makeBinaryStateChunks(externCode, scriptCodes)
		c.InitChunks(externs, chunks, opts)
		c.ParseForCompilation()
		c.Stage1Passes()
		if c.HasErrors() {
			t.Fatalf("Stage1Passes errors: %v", c.GetErrors())
		}
		stream := &nonCloseTrackingStream{}
		c.SaveState(stream)
		if stream.writeClosed {
			t.Fatalf("SaveState must not close caller outputStream")
		}
		return stream.buf.Bytes(), stream
	}

	// 1. Default mode: gzipped protobuf (starts with gzip magic 0x1f, 0x8b).
	protoBytes, _ := saveWithConfig(false, "")
	if len(protoBytes) < 16 {
		t.Fatalf("unexpectedly short default SaveState: %d bytes", len(protoBytes))
	}
	if string(protoBytes[:8]) == binarySaveStateMagic {
		t.Fatalf("default SaveState unexpectedly wrote binary header")
	}
	if protoBytes[0] != 0x1f || protoBytes[1] != 0x8b {
		t.Fatalf("default SaveState expected gzip magic 1f 8b, got %02x %02x", protoBytes[0], protoBytes[1])
	}

	// 2. Opt-in via CompilerOptions.SetBinarySaveState(true).
	binOptBytes, _ := saveWithConfig(true, "")
	if len(binOptBytes) <= binarySaveStateHeaderSize+2 {
		t.Fatalf("unexpectedly short binary SaveState: %d bytes", len(binOptBytes))
	}
	if gotMagic := string(binOptBytes[:8]); gotMagic != binarySaveStateMagic {
		t.Fatalf("binary SaveState magic = %q, want %q", gotMagic, binarySaveStateMagic)
	}
	if gotVer := binary.LittleEndian.Uint32(binOptBytes[8:12]); gotVer != binarySaveStateVersion {
		t.Fatalf("binary SaveState version = %d, want %d", gotVer, binarySaveStateVersion)
	}
	if gotTag := binary.LittleEndian.Uint32(binOptBytes[12:16]); gotTag != binarySaveStateFormatTag {
		t.Fatalf("binary SaveState formatTag = 0x%08x, want 0x%08x", gotTag, binarySaveStateFormatTag)
	}
	if binOptBytes[16] != 0x1f || binOptBytes[17] != 0x8b {
		t.Fatalf("binary SaveState payload expected gzip magic at offset 16, got %02x %02x", binOptBytes[16], binOptBytes[17])
	}

	// 3. Opt-in via JSGO_BINARY_SAVE_STATE=1 environment variable.
	binEnvBytes, _ := saveWithConfig(false, "1")
	if string(binEnvBytes[:8]) != binarySaveStateMagic {
		t.Fatalf("JSGO_BINARY_SAVE_STATE=1 did not write binary SaveState header")
	}
	if !bytes.Equal(gunzipPayload(t, binOptBytes[16:]), gunzipPayload(t, binEnvBytes[16:])) {
		t.Fatalf("binary payload differs between SetBinarySaveState(true) and JSGO_BINARY_SAVE_STATE=1")
	}

	// 4. RestoreState auto-detects both formats even when SetBinarySaveState(false) and env unset,
	//    and preserves non-closing reader semantics.
	t.Setenv("JSGO_BINARY_SAVE_STATE", "")
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"protobuf", protoBytes},
		{"binary", binOptBytes},
	} {
		rc := &nonCloseTrackingStream{reader: bytes.NewReader(tc.data)}
		restored := NewCompiler()
		externs, chunks := makeBinaryStateChunks(externCode, scriptCodes)
		restored.InitChunks(externs, chunks, newBinaryStateTestOptions(false, 1))
		restored.RestoreState(rc)
		if rc.readClosed {
			t.Fatalf("%s RestoreState must not close caller inputStream", tc.name)
		}
		if restored.HasErrors() {
			t.Fatalf("%s RestoreState errors: %v", tc.name, restored.GetErrors())
		}
	}
}

type threeStageResult struct {
	stage1StateBytes      []byte
	stage2StateBytes      []byte
	stage2RestoredExtDump string
	stage2RestoredJsDump  string
	stage2ExternProps     []string
	stage2AttrViolations  []*ConformanceViolationsProto
	stage2UnattrViol      []*ViolationProto
	stage2LastInjected    int
	renderedChunks        []RenderedChunk
	finalTree             string
}

func runThreeStagePipeline(
	t *testing.T,
	externCode string,
	scriptCodes []string,
	binaryStage1 bool,
	binaryStage2 bool,
	threads int,
	usePassthroughCompression bool,
) threeStageResult {
	t.Helper()

	applyCompression := func(opts *CompilerOptions) {
		if usePassthroughCompression {
			opts.stateCompressionWrapper = func(w io.Writer) io.Writer { return w }
			opts.stateDecompressionWrapper = func(r io.Reader) io.Reader { return r }
		}
	}

	// Stage 1: Parse + Stage1Passes + Conformance violations + SaveState
	stage1 := NewCompiler()
	opts1 := newBinaryStateTestOptions(binaryStage1, threads)
	applyCompression(opts1)
	ext1, chunks1 := makeBinaryStateChunks(externCode, scriptCodes)
	stage1.InitChunks(ext1, chunks1, opts1)
	stage1.ParseForCompilation()
	if stage1.HasErrors() {
		t.Fatalf("Stage1 parse errors: %v", stage1.GetErrors())
	}
	stage1.Stage1Passes()
	if stage1.HasErrors() {
		t.Fatalf("Stage1Passes errors: %v", stage1.GetErrors())
	}
	v1 := ViolationProto_NewBuilder().SetSourceName("script_01.js").SetLineNumber(3).SetCharNumber(4).Build()
	stage1.AddConformanceViolations(ConformanceViolationsProto_NewBuilder().AddViolations(v1).Build())
	stage1.AddUnattributedConformanceViolation(ViolationProto_NewBuilder().SetSourceName("global.js").SetLineNumber(17).Build())

	var stage1Buf bytes.Buffer
	stage1.SaveState(&stage1Buf)

	// Inspect restored Stage 1 state on a dedicated compiler while its NodeStore is active.
	inspectStage2 := NewCompiler()
	optsInspect := newBinaryStateTestOptions(false, threads)
	applyCompression(optsInspect)
	extIns, chunksIns := makeBinaryStateChunks(externCode, scriptCodes)
	inspectStage2.InitChunks(extIns, chunksIns, optsInspect)
	inspectStage2.RestoreState(bytes.NewReader(stage1Buf.Bytes()))
	extDump := snapshotDeepAst(inspectStage2.GetExternsRoot())
	jsDump := snapshotDeepAst(inspectStage2.GetJsRoot())
	var extProps []string
	if ep := inspectStage2.GetExternPropSet(); ep != nil {
		for s := range ep.AllStrings() {
			extProps = append(extProps, s)
		}
	}
	attrViol := inspectStage2.GetConformanceViolations()
	unattrViol := inspectStage2.GetUnattributedConformanceViolations()

	// Stage 2: RestoreState + Stage2Passes(OPTIMIZATIONS) + SaveState
	stage2 := NewCompiler()
	opts2 := newBinaryStateTestOptions(binaryStage2, threads)
	applyCompression(opts2)
	ext2, chunks2 := makeBinaryStateChunks(externCode, scriptCodes)
	stage2.InitChunks(ext2, chunks2, opts2)
	stage2.RestoreState(bytes.NewReader(stage1Buf.Bytes()))
	stage2.Stage2Passes(CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS)
	if stage2.HasErrors() {
		t.Fatalf("Stage2Passes errors: %v", stage2.GetErrors())
	}
	var stage2Buf bytes.Buffer
	stage2.SaveState(&stage2Buf)

	// Stage 3: RestoreState + Stage3Passes + RenderChunksAndSourceMaps
	stage3 := NewCompiler()
	opts3 := newBinaryStateTestOptions(false, threads)
	applyCompression(opts3)
	ext3, chunks3 := makeBinaryStateChunks(externCode, scriptCodes)
	stage3.InitChunks(ext3, chunks3, opts3)
	stage3.RestoreState(bytes.NewReader(stage2Buf.Bytes()))
	stage3.Stage3Passes()
	if stage3.HasErrors() {
		t.Fatalf("Stage3Passes errors: %v", stage3.GetErrors())
	}
	finalTree := stage3.GetJsRoot().ToStringTree()
	lastInjIdx := -1
	if lastLib := stage3.runtimeJsLibManager.GetLastInjectedLibrary(); !lastLib.IsNil() && stage3.GetJsRoot().HasChildren() {
		lastInjIdx = stage3.GetJsRoot().GetFirstChild().GetIndexOfChild(lastLib)
	}
	rendered, err := stage3.RenderChunksAndSourceMaps(stage3.GetChunks(), RenderChunkOptions{
		GenerateSourceMaps: true,
		ChunkFileName:      func(m *JSChunk) string { return "out/" + m.GetName() + ".js" },
		WrapperPrefix:      func(m *JSChunk) string { return "/* " + m.GetName() + " */\n" },
		NumWorkers:         threads,
	})
	if err != nil {
		t.Fatalf("RenderChunksAndSourceMaps failed: %v", err)
	}
	stage3.PerformPostCompilationTasks()

	return threeStageResult{
		stage1StateBytes:      stage1Buf.Bytes(),
		stage2StateBytes:      stage2Buf.Bytes(),
		stage2RestoredExtDump: extDump,
		stage2RestoredJsDump:  jsDump,
		stage2ExternProps:     extProps,
		stage2AttrViolations:  attrViol,
		stage2UnattrViol:      unattrViol,
		stage2LastInjected:    lastInjIdx,
		renderedChunks:        rendered,
		finalTree:             finalTree,
	}
}

func assertThreeStageParity(t *testing.T, label string, got, want threeStageResult) {
	t.Helper()
	if got.stage2RestoredExtDump != want.stage2RestoredExtDump {
		t.Fatalf("%s: Stage 2 restored externsRoot deep AST mismatch", label)
	}
	if got.stage2RestoredJsDump != want.stage2RestoredJsDump {
		t.Fatalf("%s: Stage 2 restored jsRoot deep AST mismatch", label)
	}

	if strings.Join(got.stage2ExternProps, "\x00") != strings.Join(want.stage2ExternProps, "\x00") {
		t.Fatalf("%s: externProperties mismatch: got=%v want=%v", label, got.stage2ExternProps, want.stage2ExternProps)
	}
	if len(got.stage2AttrViolations) != len(want.stage2AttrViolations) ||
		got.stage2AttrViolations[0].GetViolationsCount() != want.stage2AttrViolations[0].GetViolationsCount() ||
		got.stage2AttrViolations[0].GetViolations(0).GetSourceName() != want.stage2AttrViolations[0].GetViolations(0).GetSourceName() ||
		got.stage2AttrViolations[0].GetViolations(0).GetLineNumber() != want.stage2AttrViolations[0].GetViolations(0).GetLineNumber() {
		t.Fatalf("%s: ConformanceViolations mismatch", label)
	}
	if len(got.stage2UnattrViol) != len(want.stage2UnattrViol) ||
		got.stage2UnattrViol[0].GetSourceName() != want.stage2UnattrViol[0].GetSourceName() ||
		got.stage2UnattrViol[0].GetLineNumber() != want.stage2UnattrViol[0].GetLineNumber() {
		t.Fatalf("%s: UnattributedConformanceViolations mismatch", label)
	}
	if got.stage2LastInjected != want.stage2LastInjected {
		t.Fatalf("%s: lastInjectedLibraryIndexInFirstScript = %d, want %d",
			label, got.stage2LastInjected, want.stage2LastInjected)
	}
	if got.finalTree != want.finalTree {
		t.Fatalf("%s: Stage 3 final AST mismatch", label)
	}
	if len(got.renderedChunks) != len(want.renderedChunks) {
		t.Fatalf("%s: rendered chunk count = %d, want %d", label, len(got.renderedChunks), len(want.renderedChunks))
	}
	for i := range want.renderedChunks {
		if got.renderedChunks[i].Chunk.GetName() != want.renderedChunks[i].Chunk.GetName() {
			t.Fatalf("%s: chunk[%d] name = %q, want %q",
				label, i, got.renderedChunks[i].Chunk.GetName(), want.renderedChunks[i].Chunk.GetName())
		}
		if got.renderedChunks[i].Code != want.renderedChunks[i].Code {
			t.Fatalf("%s: chunk %s JS mismatch:\ngot:\n%s\nwant:\n%s",
				label, want.renderedChunks[i].Chunk.GetName(), got.renderedChunks[i].Code, want.renderedChunks[i].Code)
		}
		if got.renderedChunks[i].SourceMap != want.renderedChunks[i].SourceMap {
			t.Fatalf("%s: chunk %s SourceMap mismatch:\ngot:\n%s\nwant:\n%s",
				label, want.renderedChunks[i].Chunk.GetName(), got.renderedChunks[i].SourceMap, want.renderedChunks[i].SourceMap)
		}
	}
}

func TestBinarySaveState_ThreeStageParityWithProtobuf(t *testing.T) {
	externCode, scriptCodes := buildBinaryStateTestCorpus(12, 0)

	protoBase := runThreeStagePipeline(t, externCode, scriptCodes, false, false, 1, false)
	binSeq := runThreeStagePipeline(t, externCode, scriptCodes, true, true, 1, false)
	assertThreeStageParity(t, "binary(1-thread) vs protobuf", binSeq, protoBase)

	// Mixed Stage 1 Protobuf -> Stage 2 Binary -> Stage 3 auto-detects both formats seamlessly.
	mixed := runThreeStagePipeline(t, externCode, scriptCodes, false, true, 1, false)
	assertThreeStageParity(t, "mixed(proto->bin) vs protobuf", mixed, protoBase)

	// Custom StateCompressionWrapper / StateDecompressionWrapper parity.
	binCustomWrapper := runThreeStagePipeline(t, externCode, scriptCodes, true, true, 1, true)
	assertThreeStageParity(t, "binary(custom-wrapper) vs protobuf", binCustomWrapper, protoBase)
}

func TestBinarySaveState_PreTypecheckRawStageParity(t *testing.T) {
	externCode, scriptCodes := buildBinaryStateTestCorpus(8, 0)

	configureUntyped := func(opts *CompilerOptions) {
		opts.SetCheckTypes(false)
		opts.inferTypes = false
		opts.SetDisambiguateProperties(false)
		opts.SetAmbiguateProperties(false)
	}

	runRaw := func(binaryMode bool) string {
		c1 := NewCompiler()
		opts1 := newBinaryStateTestOptions(binaryMode, 2)
		configureUntyped(opts1)
		ext1, chunks1 := makeBinaryStateChunks(externCode, scriptCodes)
		c1.InitChunks(ext1, chunks1, opts1)
		c1.ParseForCompilation()
		c1.Stage1Passes()
		if c1.HasErrors() {
			t.Fatalf("untyped Stage1 errors: %v", c1.GetErrors())
		}
		if c1.colorRegistry != nil {
			t.Fatal("expected nil colorRegistry in untyped Stage 1 compilation")
		}
		var stateBuf bytes.Buffer
		c1.SaveState(&stateBuf)

		c2 := NewCompiler()
		opts2 := newBinaryStateTestOptions(false, 2)
		configureUntyped(opts2)
		ext2, chunks2 := makeBinaryStateChunks(externCode, scriptCodes)
		c2.InitChunks(ext2, chunks2, opts2)
		c2.RestoreState(bytes.NewReader(stateBuf.Bytes()))
		c2.Stage2Passes(CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS)
		c2.Stage3Passes()
		if c2.HasErrors() {
			t.Fatalf("restored untyped compile errors: %v", c2.GetErrors())
		}
		return c2.ToSource()
	}

	want := runRaw(false)
	got := runRaw(true)
	if got != want {
		t.Fatalf("untyped SaveState/RestoreState JS mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestBinarySaveState_ParallelDeterminismAndLargeStringTable(t *testing.T) {
	// Include >4,096 distinct strings to exercise parallel readBinaryStringPool interning
	// alongside parallel script encoding and parallel PrebuildAst deserialization.
	externCode, scriptCodes := buildBinaryStateTestCorpus(16, 4500)

	protoRef := runThreeStagePipeline(t, externCode, scriptCodes, false, false, 1, false)
	bin1 := runThreeStagePipeline(t, externCode, scriptCodes, true, true, 1, false)
	assertThreeStageParity(t, "large-corpus binary(threads=1) vs protobuf", bin1, protoRef)

	rawStage1Seq := gunzipPayload(t, bin1.stage1StateBytes[binarySaveStateHeaderSize:])
	rawStage2Seq := gunzipPayload(t, bin1.stage2StateBytes[binarySaveStateHeaderSize:])

	for iter := range 2 {
		bin4 := runThreeStagePipeline(t, externCode, scriptCodes, true, true, 4, false)
		assertThreeStageParity(t, fmt.Sprintf("iter=%d binary(threads=4) vs threads=1", iter), bin4, bin1)

		rawStage1Par := gunzipPayload(t, bin4.stage1StateBytes[binarySaveStateHeaderSize:])
		rawStage2Par := gunzipPayload(t, bin4.stage2StateBytes[binarySaveStateHeaderSize:])
		if !bytes.Equal(rawStage1Par, rawStage1Seq) {
			t.Fatalf("iter=%d Stage 1 decompressed binary payload differs between threads=1 (%d B) and threads=4 (%d B)",
				iter, len(rawStage1Seq), len(rawStage1Par))
		}
		if !bytes.Equal(rawStage2Par, rawStage2Seq) {
			t.Fatalf("iter=%d Stage 2 decompressed binary payload differs between threads=1 (%d B) and threads=4 (%d B)",
				iter, len(rawStage2Seq), len(rawStage2Par))
		}
	}
}

func catchRestoreThrowable(t *testing.T, externCode string, scriptCodes []string, raw bool, data io.Reader) jutil.Throwable {
	t.Helper()
	c := NewCompiler()
	opts := newBinaryStateTestOptions(false, 1)
	if raw {
		opts.stateDecompressionWrapper = func(r io.Reader) io.Reader { return r }
	}
	externs, chunks := makeBinaryStateChunks(externCode, scriptCodes)
	c.InitChunks(externs, chunks, opts)

	var caught jutil.Throwable
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				th, ok := rec.(jutil.Throwable)
				if !ok {
					t.Fatalf("RestoreState panicked with non-Throwable (%T): %v", rec, rec)
				}
				caught = th
			}
		}()
		c.RestoreState(data)
	}()
	return caught
}

func hasIOExceptionCause(th jutil.Throwable) bool {
	for cur := th; cur != nil; cur = cur.GetCause() {
		if _, ok := cur.(jutil.IIOException); ok {
			return true
		}
	}
	return false
}

func TestBinarySaveState_TruncatedCorruptAndVersionErrors(t *testing.T) {
	externCode, scriptCodes := buildBinaryStateTestCorpus(4, 0)

	// 1. Empty stream & empty gzip stream preserve legacy protobuf NullPointerException semantics.
	var emptyGzip bytes.Buffer
	gz := gzip.NewWriter(&emptyGzip)
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	if caught := catchRestoreThrowable(t, externCode, scriptCodes, false, bytes.NewReader(emptyGzip.Bytes())); caught == nil {
		t.Fatal("expected NullPointerException on empty gzip stream")
	} else if _, ok := caught.(jutil.INullPointerException); !ok {
		t.Fatalf("empty gzip stream: got %T (%v), want INullPointerException", caught, caught)
	}

	if caught := catchRestoreThrowable(t, externCode, scriptCodes, true, bytes.NewReader(nil)); caught == nil {
		t.Fatal("expected NullPointerException on empty uncompressed stream")
	} else if _, ok := caught.(jutil.INullPointerException); !ok {
		t.Fatalf("empty uncompressed stream: got %T (%v), want INullPointerException", caught, caught)
	}

	// 2. Short non-EOF error during Peek(8) is replayed cleanly as IOException.
	shortErrStream := &failingShortReader{data: []byte("JSG"), err: errors.New("synthetic disk read failure")}
	if caught := catchRestoreThrowable(t, externCode, scriptCodes, true, shortErrStream); !hasIOExceptionCause(caught) {
		t.Fatalf("short failing reader: got %T (%v), want IOException cause", caught, caught)
	}

	// 3. Truncated binary header (8-byte magic present, < 16 bytes total).
	if caught := catchRestoreThrowable(t, externCode, scriptCodes, false, bytes.NewReader([]byte("JSGOSAV1\x01\x00"))); !hasIOExceptionCause(caught) {
		t.Fatalf("truncated binary header: got %T (%v), want IOException cause", caught, caught)
	}

	// Build a valid uncompressed binary state image so we can mutate exact offsets.
	baseRun := runThreeStagePipeline(t, externCode, scriptCodes, true, true, 1, true)
	validRaw := baseRun.stage1StateBytes

	// 4. Unsupported version -> IllegalStateException.
	for _, badVer := range []uint32{0, 2, 99} {
		mut := bytes.Clone(validRaw)
		binary.LittleEndian.PutUint32(mut[8:12], badVer)
		caught := catchRestoreThrowable(t, externCode, scriptCodes, true, bytes.NewReader(mut))
		if _, ok := caught.(jutil.IIllegalStateException); !ok {
			t.Fatalf("bad version %d: got %T (%v), want IIllegalStateException", badVer, caught, caught)
		}
	}

	// 5. Incompatible formatTag -> IllegalStateException.
	{
		mut := bytes.Clone(validRaw)
		binary.LittleEndian.PutUint32(mut[12:16], binarySaveStateFormatTag^0x1)
		caught := catchRestoreThrowable(t, externCode, scriptCodes, true, bytes.NewReader(mut))
		if _, ok := caught.(jutil.IIllegalStateException); !ok {
			t.Fatalf("bad formatTag: got %T (%v), want IIllegalStateException", caught, caught)
		}
	}

	// 6. Truncated payloads at multiple cut points -> clean Throwable (no Go runtime slice panic).
	for _, cut := range []int{16, 20, len(validRaw) / 4, len(validRaw) / 2, len(validRaw) - 5} {
		caught := catchRestoreThrowable(t, externCode, scriptCodes, true, bytes.NewReader(validRaw[:cut]))
		if caught == nil {
			t.Fatalf("cut=%d: expected Throwable on truncated binary payload", cut)
		}
	}

	// 7. Oversized section size & table count sanity caps reject corrupt headers without OOM.
	{
		mut := bytes.Clone(validRaw)
		// Offset 16 is uint32(len(stateBytes)): set to 0xFFFFFFFF (4 GiB - 1) on a tiny buffer;
		// the 4 MiB probe in (*binaryStreamReader).bytes rejects on unexpected EOF without OOM.
		binary.LittleEndian.PutUint32(mut[16:20], ^uint32(0))
		caught := catchRestoreThrowable(t, externCode, scriptCodes, true, bytes.NewReader(mut))
		if !hasIOExceptionCause(caught) {
			t.Fatalf("oversized section length: got %T (%v), want IOException cause", caught, caught)
		}
	}
	{
		var srCaught any
		func() {
			defer func() { srCaught = recover() }()
			sr := &binaryStreamReader{r: bytes.NewReader(nil)}
			_ = sr.bytes(maxBinarySaveStateSectionBytes + 1)
		}()
		if srCaught == nil || !hasIOExceptionCause(jutil.AsThrowable(srCaught)) {
			t.Fatalf("section > maxBinarySaveStateSectionBytes: got %T (%v), want IOException cause", srCaught, srCaught)
		}
	}
	{
		mut := bytes.Clone(validRaw)
		stateLen := int(binary.LittleEndian.Uint32(mut[16:20]))
		numStringsOff := 20 + stateLen + 4 // after maxLength
		binary.LittleEndian.PutUint32(mut[numStringsOff:numStringsOff+4], uint32(maxBinarySaveStateCount+1))
		caught := catchRestoreThrowable(t, externCode, scriptCodes, true, bytes.NewReader(mut))
		if !hasIOExceptionCause(caught) {
			t.Fatalf("oversized string count: got %T (%v), want IOException cause", caught, caught)
		}
	}

	// 8. Corrupt per-script binary AstNode buffers wrap IOException in MalformedTypedAstException.
	for _, tc := range []struct {
		name string
		buf  []byte
	}{
		{"short-header", []byte{1, 0, 0, 0}},
		{"zero-nodes", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"unknown-flags", []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0x00, 0x80}},
		{"bad-kind", []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 0x00, 0x00}},
		{"oob-str-ptr", []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, byte(binNodeValString), 0x00, 99}},
		{"bad-root-endIdx", []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, byte(binNodeHasChildren), 0x00, 5}},
		{"varint-overflow", []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, byte(binNodeValString), 0x00, 0xff, 0xff, 0xff, 0xff, 0x1f}},
	} {
		_ = NewCompiler()
		sd := &ScriptNodeDeserializer{
			sourceFile:   SourceFile_FromCode("bad.js", ""),
			binaryScript: tc.buf,
			stringPool:   StringPool_Empty(),
		}
		var sdCaught any
		func() {
			defer func() { sdCaught = recover() }()
			sd.DeserializeNewWithAllocator(rhino.NewNodeAllocator())
		}()
		if _, ok := sdCaught.(*MalformedTypedAstException); !ok {
			t.Fatalf("corrupt binaryScript %s: got %T (%v), want *MalformedTypedAstException", tc.name, sdCaught, sdCaught)
		}
	}
}

func TestBinaryStreamReaderBytesBoundedChunksAndVarintScriptRoundTrip(t *testing.T) {
	// 1. Payload > 4 MiB allocates exact slice capacity without bytes.Buffer doublings.
	const payloadSize = 6<<20 + 12345
	payload := make([]byte, payloadSize)
	for i := range payload {
		payload[i] = byte((i*31 + 17) & 0xff)
	}
	sr := &binaryStreamReader{r: bytes.NewReader(payload)}
	got := sr.bytes(uint64(payloadSize))
	if len(got) != payloadSize || cap(got) != payloadSize {
		t.Fatalf("sr.bytes(%d): len=%d cap=%d, want exact %d (no bytes.Buffer doubling slack)",
			payloadSize, len(got), cap(got), payloadSize)
	}
	if !bytes.Equal(got, payload) {
		t.Fatal("sr.bytes(>4MiB) payload mismatch")
	}

	// 2. Truncated > 4 MiB streams (both before and after the 4 MiB probe) fail cleanly with IOException.
	for _, truncLen := range []int{1024, 5 << 20} {
		var caught any
		func() {
			defer func() { caught = recover() }()
			r := &binaryStreamReader{r: bytes.NewReader(payload[:truncLen])}
			_ = r.bytes(uint64(payloadSize))
		}()
		if caught == nil || !hasIOExceptionCause(jutil.AsThrowable(caught)) {
			t.Fatalf("truncLen=%d: got %T (%v), want IOException", truncLen, caught, caught)
		}
	}

	// 3. Multi-byte varint script fields (> 200 strings, JSDoc tags/descriptions, cooked + raw-only templates).
	var sb strings.Builder
	sb.WriteString("/** @fileoverview @license Apache-2.0 */\n")
	sb.WriteString("function tag(s, ...v) { return s.raw[0] + (v[0] ?? ''); }\n")
	for i := 0; i < 220; i++ {
		fmt.Fprintf(&sb, "/** @desc desc_%d @meaning mean_%d */ var MSG_VAR_%d = `val_%d_${%d}`;\n", i, i, i, i, i)
	}
	sb.WriteString("export const rawEsc = tag`\\unicode_${MSG_VAR_219}`;\n")

	externCode := "/** @externs */\nfunction alert(x) {}\n"
	scriptCodes := []string{sb.String()}

	cSave := NewCompiler()
	optsSave := newBinaryStateTestOptions(true, 2)
	extSave, chunksSave := makeBinaryStateChunks(externCode, scriptCodes)
	cSave.InitChunks(extSave, chunksSave, optsSave)
	cSave.ParseForCompilation()
	if cSave.HasErrors() {
		t.Fatalf("parse errors: %v", cSave.GetErrors())
	}
	cSave.Stage1Passes()
	if cSave.HasErrors() {
		t.Fatalf("stage1 errors: %v", cSave.GetErrors())
	}
	var saved bytes.Buffer
	cSave.SaveState(&saved)

	cProto := NewCompiler()
	optsProto := newBinaryStateTestOptions(false, 2)
	extProto, chunksProto := makeBinaryStateChunks(externCode, scriptCodes)
	cProto.InitChunks(extProto, chunksProto, optsProto)
	cProto.ParseForCompilation()
	cProto.Stage1Passes()
	var savedProto bytes.Buffer
	cProto.SaveState(&savedProto)

	compileRestored := func(data []byte) string {
		c := NewCompiler()
		opts := newBinaryStateTestOptions(false, 2)
		ext, ch := makeBinaryStateChunks(externCode, scriptCodes)
		c.InitChunks(ext, ch, opts)
		c.RestoreState(bytes.NewReader(data))
		c.Stage2Passes(CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS)
		c.Stage3Passes()
		c.PerformPostCompilationTasks()
		if c.HasErrors() {
			t.Fatalf("restored compile errors: %v", c.GetErrors())
		}
		return c.ToSource()
	}

	gotBin := compileRestored(saved.Bytes())
	wantProto := compileRestored(savedProto.Bytes())
	if gotBin != wantProto {
		t.Fatalf("varint binary SaveState output mismatch vs protobuf:\ngot:\n%s\nwant:\n%s", gotBin, wantProto)
	}
}

func BenchmarkBinarySaveRestore(b *testing.B) {
	externCode, scriptCodes := buildBinaryStateTestCorpus(96, 0)
	identityOut := func(w io.Writer) io.Writer { return w }
	identityIn := func(r io.Reader) io.Reader { return r }

	prepareStage1Compiler := func(binaryMode bool, threads int, rawStream bool) (*Compiler, []byte) {
		c := NewCompiler()
		opts := newBinaryStateTestOptions(binaryMode, threads)
		if rawStream {
			opts.stateCompressionWrapper = identityOut
		}
		externs, chunks := makeBinaryStateChunks(externCode, scriptCodes)
		c.InitChunks(externs, chunks, opts)
		c.ParseForCompilation()
		if c.HasErrors() {
			b.Fatalf("parse errors: %v", c.GetErrors())
		}
		c.Stage1Passes()
		if c.HasErrors() {
			b.Fatalf("stage1 errors: %v", c.GetErrors())
		}
		var warmup bytes.Buffer
		c.SaveState(&warmup)
		return c, warmup.Bytes()
	}

	for _, threads := range []int{1, 4} {
		for _, mode := range []struct {
			name   string
			binary bool
		}{
			{"Protobuf", false},
			{"Binary", true},
		} {
			b.Run(fmt.Sprintf("%s/SaveState/threads=%d", mode.name, threads), func(b *testing.B) {
				stage1Comp, wireBytes := prepareStage1Compiler(mode.binary, threads, false)
				var buf bytes.Buffer
				buf.Grow(len(wireBytes) + 4096)
				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					buf.Reset()
					stage1Comp.SaveState(&buf)
				}
				b.ReportMetric(float64(buf.Len()), "wire-bytes")
			})

			b.Run(fmt.Sprintf("%s/RestoreState/threads=%d", mode.name, threads), func(b *testing.B) {
				_, wireBytes := prepareStage1Compiler(mode.binary, threads, false)
				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					c := NewCompiler()
					opts := newBinaryStateTestOptions(false, threads)
					externs, chunks := makeBinaryStateChunks(externCode, scriptCodes)
					c.InitChunks(externs, chunks, opts)
					c.RestoreState(bytes.NewReader(wireBytes))
				}
				b.ReportMetric(float64(len(wireBytes)), "wire-bytes")
			})
		}
	}

	for _, mode := range []struct {
		name   string
		binary bool
	}{
		{"Protobuf", false},
		{"Binary", true},
	} {
		b.Run(fmt.Sprintf("%s/RawSaveState/threads=4", mode.name), func(b *testing.B) {
			stage1Comp, wireBytes := prepareStage1Compiler(mode.binary, 4, true)
			var buf bytes.Buffer
			buf.Grow(len(wireBytes) + 4096)
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				buf.Reset()
				stage1Comp.SaveState(&buf)
			}
			b.ReportMetric(float64(buf.Len()), "wire-bytes")
		})

		b.Run(fmt.Sprintf("%s/RawRestoreState/threads=4", mode.name), func(b *testing.B) {
			_, wireBytes := prepareStage1Compiler(mode.binary, 4, true)
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				c := NewCompiler()
				opts := newBinaryStateTestOptions(false, 4)
				opts.stateDecompressionWrapper = identityIn
				externs, chunks := makeBinaryStateChunks(externCode, scriptCodes)
				c.InitChunks(externs, chunks, opts)
				c.RestoreState(bytes.NewReader(wireBytes))
			}
			b.ReportMetric(float64(len(wireBytes)), "wire-bytes")
		})
	}
}
