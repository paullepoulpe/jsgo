package jscomp

// Go-only: test-helper — the differential oracle for the TypedAST
// serialization passes (TypedAstSerializer, SerializeTypedAstPass,
// ConvertTypesToColors) and, through the round trip, the deserializer: the
// Go mirror of tools/oracle/java/SerializationDump.java, whose header
// documents the variants (typed, colors, untyped) and the output format.
//
//   - TestSerializationDumpGolden compares against the checked-in Java output
//     testdata/serialization/golden.txt.gz for the inputs listed in
//     testdata/serialization/files.txt (sources in testdata/serialization/src;
//     always runs). Regenerate with
//       java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//         com.google.javascript.jscomp.SerializationDump --externs=tools/oracle/typecheck-externs.js \
//         --full-every=4 internal/jscomp/testdata/serialization/golden.txt.gz \
//         internal/jscomp/testdata/serialization/src --list=internal/jscomp/testdata/serialization/files.txt
//   - TestSerializationDumpOracle compares any other run when
//     SERIALIZATION_ORACLE (the Java <out.gz>), SERIALIZATION_ROOT and
//     SERIALIZATION_FILES (the list given with --list) are set, with
//     SERIALIZATION_SAMPLE and SERIALIZATION_FULL_EVERY matching the Java
//     run's --sample and --full-every (SERIALIZATION_FULL=1 for --full);
//     SERIALIZATION_DIFF_DIR, if set, receives <rel>.go.txt and <rel>.java.txt
//     for every mismatching file. Use absolute paths.

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var serializationDump_VARIANTS = []string{"typed", "colors", "untyped"}

// serializationDump_setup is SerializationDump.Setup.
type serializationDump_setup struct {
	errors      *typeCheckDump_errorManager
	compiler    *Compiler
	options     *CompilerOptions
	externsFile *SourceFile
	inputFile   *SourceFile
}

func newSerializationDump_setup(externsFile, inputFile *SourceFile) *serializationDump_setup {
	s := &serializationDump_setup{
		errors:      newTypeCheckDump_errorManager(),
		compiler:    NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard)),
		options:     typeCheckDump_options(),
		externsFile: externsFile,
		inputFile:   inputFile,
	}
	s.compiler.SetErrorManager(s.errors)
	s.compiler.Init([]*SourceFile{externsFile}, []*SourceFile{inputFile}, s.options)
	return s
}

// serializationDump_dumpFile is SerializationDump.dumpFile.
func serializationDump_dumpFile(sb *strings.Builder, externs, rel, src string, fullThis bool) {
	typeCheckDump_line(sb, "## file "+rel)
	externsFile := SourceFile_FromCode("externs.js", externs)
	inputFile := SourceFile_FromCode(rel, src)
	{
		s := newSerializationDump_setup(externsFile, inputFile)
		var thrown jutil.Throwable
		jutil.Try(func() { s.compiler.Parse() }, func(e jutil.IRuntimeException) { thrown = e })
		if thrown != nil {
			typeCheckDump_line(sb, "# parse exception "+typeCheckDump_className(thrown)+": "+
				jutil.ToJavaString(jutil.GetMessageOrNull(thrown)))
			return
		}
		ok := s.errors.GetErrorCount() == 0
		if ok {
			typeCheckDump_line(sb, "# parse ok")
		} else {
			typeCheckDump_line(sb, "# parse errors")
			for _, d := range s.errors.lines {
				typeCheckDump_line(sb, "diag "+d)
			}
			return
		}
	}
	for _, variant := range serializationDump_VARIANTS {
		typeCheckDump_line(sb, "## variant "+variant)
		serializationDump_dumpVariant(sb, variant, externsFile, inputFile, fullThis)
	}
}

// serializationDump_dumpVariant is SerializationDump.dumpVariant.
func serializationDump_dumpVariant(sb *strings.Builder, variant string, externsFile, inputFile *SourceFile, fullThis bool) {
	s := newSerializationDump_setup(externsFile, inputFile)
	stage := "checks"
	jutil.Try(func() {
		s.compiler.Parse()
		externsRoot := s.compiler.GetExternsRoot()
		jsRoot := s.compiler.GetJsRoot()
		if variant != "untyped" {
			serializationDump_typeCheck(s.compiler, s.options)
		}
		if variant == "colors" {
			NewRemoveCastNodes(s.compiler).Process(externsRoot, jsRoot)
			NewConvertTypesToColors(
				s.compiler, SerializationOptions_Builder_().SetIncludeDebugInfo(false).Build()).
				Process(externsRoot, jsRoot)
			s.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC)
		}
		typeCheckDump_line(sb, "checks "+strconv.Itoa(len(s.errors.lines))+" "+
			serializationDump_sha16(serializationDump_prefixed(s.errors.lines)))

		stage = "serialize"
		var out bytes.Buffer
		SerializeTypedAstPass_CreateFromOutputStream(
			s.compiler,
			&out,
			SerializationOptions_Builder_().
				SetIncludeDebugInfo(false).
				SetRuntimeLibraries([]string{}).
				Build()).
			Process(externsRoot, jsRoot)
		data := out.Bytes()
		sum := sha256.Sum256(data)
		typeCheckDump_line(sb, "typedast "+strconv.Itoa(len(data))+" "+hex.EncodeToString(sum[:])[:16])
		var text []string
		ast := TypedAst_List_ParseFrom(data).GetTypedAsts(0)
		text = serializationDump_addLines(text, ast.ToBuilder().ClearExternAst().ClearCodeAst().Build().String())
		for _, lazy := range ast.GetExternAstList() {
			text = append(text, "extern_ast")
			text = serializationDump_addLazyAst(text, lazy)
		}
		for _, lazy := range ast.GetCodeAstList() {
			text = append(text, "code_ast")
			text = serializationDump_addLazyAst(text, lazy)
		}
		serializationDump_section(sb, "text", text, fullThis)

		stage = "deserialize"
		d := newSerializationDump_setup(externsFile, inputFile)
		requiredInputFiles := []*SourceFile{externsFile, inputFile}
		deserialized := TypedAstDeserializer_DeserializeFullAst(
			d.compiler,
			SourceFile_FromCodeWithKind("syntheticExterns", "", rhino.StaticSourceFile_SourceKind_EXTERN),
			requiredInputFiles,
			bytes.NewReader(data),
			true,
			d.options.GetResolveSourceMapAnnotations(),
			d.options.GetParseInlineSourceMaps())
		var tree []string
		colorsSeen := jutil.NewLinkedHashMap[colors.ColorId, *colors.Color]()
		for _, f := range requiredInputFiles {
			entry, ok := deserialized.filesystem[f]
			if !ok {
				tree = append(tree, "script "+f.GetName()+" absent")
				continue
			}
			tree = append(tree, "script "+f.GetName())
			tree = serializationDump_dumpTree(tree, colorsSeen, entry.DeserializeNew(), 0)
		}
		for c := range colorsSeen.Values() {
			tree = append(tree, serializationDump_colorLine(c))
		}
		serializationDump_section(sb, "deserialized", tree, fullThis)
	}, func(e jutil.IRuntimeException) {
		typeCheckDump_line(sb, "exception "+stage+" "+typeCheckDump_className(e)+": "+
			jutil.ToJavaString(jutil.GetMessageOrNull(e)))
	})
}

// serializationDump_typeCheck is SerializationDump.typeCheck (TypeCheckDump's
// type checking steps).
func serializationDump_typeCheck(compiler *Compiler, options *CompilerOptions) {
	externsRoot := compiler.GetExternsRoot()
	jsRoot := compiler.GetJsRoot()
	NewGatherModuleMetadata(compiler, false, ModuleLoader_ResolutionMode_BROWSER).Process(externsRoot, jsRoot)
	NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).Process(externsRoot, jsRoot)
	NewInferConsts(compiler).Process(externsRoot, jsRoot)

	// DefaultPassConfig.inferTypes
	globalRoot := jsRoot.GetParent()
	inferencePass := NewTypeInferencePass(
		compiler,
		compiler.GetReverseAbstractInterpreter(),
		compiler.GetTypedScopeCreator().(*TypedScopeCreator))
	compiler.SetTypeCheckingHasRun(true)
	compiler.SetTopScope(inferencePass.InferAllScopes(globalRoot))

	// DefaultPassConfig.checkTypes
	check := NewTypeCheckWithScope(
		compiler,
		compiler.GetReverseAbstractInterpreter(),
		compiler.GetTypeRegistry(),
		compiler.GetTopScope(),
		compiler.GetTypedScopeCreator().(*TypedScopeCreator)).
		ReportUnknownTypes(
			options.Enables(DiagnosticGroup_ForType(TypeCheck_UNKNOWN_EXPR_TYPE))).
		ReportMissingProperties(
			!options.Disables(DiagnosticGroup_ForType(TypeCheck_INEXISTENT_PROPERTY)))
	check.Process(externsRoot, jsRoot)
	compiler.GetErrorManager().SetTypedPercent(check.GetTypedPercent())
}

func serializationDump_addLazyAst(text []string, lazy *LazyAst) []string {
	text = serializationDump_addLines(text, lazy.ToBuilder().ClearScript().Build().String())
	// As ScriptNodeDeserializer reads the script: without protobuf's default
	// nesting limit (100), which deep ASTs exceed.
	script := &AstNode{memoizedSize: -1}
	script.mergeFromReader(&protoReader{b: lazy.GetScript(), limit: math.MaxInt32})
	script.frozen = true
	return serializationDump_addLines(text, script.String())
}

// serializationDump_addLines is SerializationDump.addLines: s split at '\n',
// without the empty string after a final '\n'.
func serializationDump_addLines(out []string, s string) []string {
	out = append(out, strings.Split(s, "\n")...)
	if len(out) > 0 && out[len(out)-1] == "" && strings.HasSuffix(s, "\n") {
		out = out[:len(out)-1]
	}
	return out
}

func serializationDump_dumpTree(out []string, colorsSeen *jutil.LinkedHashMap[colors.ColorId, *colors.Color], n rhino.Node, level int) []string {
	line := strings.Repeat("    ", level) + n.ToStringWithFlags(true, true, false)
	if c := n.GetColor(); c != nil {
		line += " : " + c.GetId().String()
		serializationDump_collectColor(colorsSeen, c)
	}
	out = append(out, line)
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		out = serializationDump_dumpTree(out, colorsSeen, child, level+1)
	}
	return out
}

func serializationDump_collectColor(colorsSeen *jutil.LinkedHashMap[colors.ColorId, *colors.Color], c *colors.Color) {
	if colorsSeen.ContainsKey(c.GetId()) {
		return
	}
	colorsSeen.Put(c.GetId(), c)
	for _, x := range c.GetPrototypes() {
		serializationDump_collectColor(colorsSeen, x)
	}
	for _, x := range c.GetInstanceColors() {
		serializationDump_collectColor(colorsSeen, x)
	}
	for _, x := range c.GetUnionElements() {
		serializationDump_collectColor(colorsSeen, x)
	}
}

func serializationDump_colorLine(c *colors.Color) string {
	box := "null"
	if boxId, ok := c.GetBoxId(); ok {
		box = boxId.String()
	}
	return "color " + c.GetId().String() +
		" inv=" + strconv.FormatBool(c.IsInvalidating()) +
		" ctor=" + strconv.FormatBool(c.IsConstructor()) +
		" keep=" + strconv.FormatBool(c.GetPropertiesKeepOriginalName()) +
		" assert=" + strconv.FormatBool(c.IsClosureAssert()) +
		" box=" + box +
		" own=" + jutil.CollectionToString(c.GetOwnProperties()) +
		" proto=" + serializationDump_ids(c.GetPrototypes()) +
		" inst=" + serializationDump_ids(c.GetInstanceColors()) +
		" union=" + serializationDump_ids(c.GetUnionElements())
}

func serializationDump_ids(cs []*colors.Color) string {
	var out []string
	for _, c := range cs {
		out = append(out, c.GetId().String())
	}
	return "[" + strings.Join(out, ", ") + "]"
}

func serializationDump_prefixed(diags []string) []string {
	out := make([]string, 0, len(diags))
	for _, d := range diags {
		out = append(out, "diag "+d)
	}
	return out
}

func serializationDump_section(sb *strings.Builder, name string, lines []string, fullThis bool) {
	typeCheckDump_line(sb, name+" "+strconv.Itoa(len(lines))+" "+serializationDump_sha16(lines))
	if fullThis {
		for _, l := range lines {
			typeCheckDump_line(sb, l)
		}
	}
}

func serializationDump_sha16(lines []string) string {
	h := sha256.New()
	for _, l := range lines {
		h.Write([]byte(typeCheckDump_escape(l) + "\n"))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// serializationDump_split splits a Java dump into its "## file" sections.
func serializationDump_split(dump string) (rels []string, sections map[string]string) {
	sections = map[string]string{}
	var cur string
	var b strings.Builder
	flush := func() {
		if cur != "" {
			sections[cur] = b.String()
		}
		b.Reset()
	}
	for _, l := range strings.SplitAfter(dump, "\n") {
		if strings.HasPrefix(l, "## file ") {
			flush()
			cur = strings.TrimSuffix(strings.TrimPrefix(l, "## file "), "\n")
			rels = append(rels, cur)
		}
		b.WriteString(l)
	}
	flush()
	return rels, sections
}

func serializationDump_readGz(t *testing.T, path string) string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// serializationDump_compare runs the Go side over files (every sample-th
// one, with the Java run's full settings) and compares each file's section
// with the Java dump.
func serializationDump_compare(t *testing.T, externs, root, oracle string, files []string, sample, fullEvery int, full bool, diffDir string) {
	t.Helper()
	oracleDump_compare(t, serializationDump_dumpFile, externs, root, oracle, files, sample, fullEvery, full, diffDir)
}

// oracleDump_compare runs dumpFile (the Go side of a per-file dump in
// SerializationDump's layout) over files (every sample-th one, with the Java
// run's full settings) and compares each file's section with the Java dump.
func oracleDump_compare(t *testing.T, dumpFile func(sb *strings.Builder, externs, rel, src string, fullThis bool),
	externs, root, oracle string, files []string, sample, fullEvery int, full bool, diffDir string) {
	t.Helper()
	_, want := serializationDump_split(serializationDump_readGz(t, oracle))
	compared, mismatches, dumped := 0, 0, 0
	for i, rel := range files {
		if i%sample != 0 {
			continue
		}
		fullThis := full || (fullEvery > 0 && dumped%fullEvery == 0)
		dumped++
		w, ok := want[rel]
		if !ok {
			t.Errorf("%s: not in the Java dump", rel)
			continue
		}
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		var sb strings.Builder
		dumpFile(&sb, externs, rel, string(src), fullThis)
		got := sb.String()
		compared++
		if got == w {
			continue
		}
		mismatches++
		if mismatches <= 10 {
			t.Errorf("%s: Go output differs from Java:\n%s", rel, typeCheckDump_firstDiff(w, got))
		}
		if diffDir != "" {
			out := filepath.Join(diffDir, rel)
			if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(out+".go.txt", []byte(got), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(out+".java.txt", []byte(w), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	t.Logf("%d files compared, %d differ", compared, mismatches)
	if mismatches > 0 {
		t.Errorf("%d of %d files differ", mismatches, compared)
	}
}

func TestSerializationDumpGolden(t *testing.T) {
	dir := filepath.Join("testdata", "serialization")
	files := typeCheckDump_readList(t, filepath.Join(dir, "files.txt"))
	serializationDump_compare(t, typeCheckDump_externs(t), filepath.Join(dir, "src"),
		filepath.Join(dir, "golden.txt.gz"), files, 1, 4, false, os.Getenv("SERIALIZATION_DIFF_DIR"))
}

func TestSerializationDumpOracle(t *testing.T) {
	oracle := os.Getenv("SERIALIZATION_ORACLE")
	root := os.Getenv("SERIALIZATION_ROOT")
	list := os.Getenv("SERIALIZATION_FILES")
	if oracle == "" || root == "" || list == "" {
		t.Skip("set SERIALIZATION_ORACLE, SERIALIZATION_ROOT and SERIALIZATION_FILES to run the corpus oracle")
	}
	atoi := func(name string, def int) int {
		s := os.Getenv(name)
		if s == "" {
			return def
		}
		n, err := strconv.Atoi(s)
		if err != nil || n < 0 {
			t.Fatalf("%s=%q: want a non-negative integer", name, s)
		}
		return n
	}
	sample := atoi("SERIALIZATION_SAMPLE", 1)
	if sample < 1 {
		sample = 1
	}
	files := typeCheckDump_readList(t, list)
	serializationDump_compare(t, typeCheckDump_externs(t), root, oracle, files, sample,
		atoi("SERIALIZATION_FULL_EVERY", 0), os.Getenv("SERIALIZATION_FULL") == "1", os.Getenv("SERIALIZATION_DIFF_DIR"))
}
