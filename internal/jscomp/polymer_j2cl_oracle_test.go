package jscomp

// Go-only: test-helper — differential test against
// tools/oracle/java/PolymerJ2clDump.java: PolymerPass and the J2CL passes run
// on the JavaScript inputs of their Java unit tests (the corpus has no Polymer
// or J2CL code), set up the way CompilerTestCase sets them up, printed in the
// same format and compared block by block with the Java output.
//
//   - TestPolymerJ2clGolden replays the checked-in Java dumps
//     testdata/polymer_j2cl/{polymer,j2cl}-golden.txt.gz over the inputs
//     testdata/polymer_j2cl/{polymer,j2cl}-inputs.txt (the strings of the Java
//     tests, extracted as "=== <name> <byte length>\n<source>\n"). Regenerate with
//       CP=.cache/closure-compiler.jar:<PolymerJ2clDump build dir>
//       J="java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp $CP com.google.javascript.jscomp.PolymerJ2clDump"
//       D=internal/jscomp/testdata/polymer_j2cl
//       $J polymer $D/polymer-golden.txt.gz $D/polymer-inputs.txt $D/polymer-externs.js
//       $J j2cl $D/j2cl-golden.txt.gz $D/j2cl-inputs.txt
//     POLYMER_J2CL_DIFF_DIR=<dir> writes the Go dump of every differing input
//     there.

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const polymerJ2clOracleDir = "testdata/polymer_j2cl"

type polymerJ2clOracleInput struct {
	name, code string
}

// polymerJ2clOracleReadInputs is PolymerJ2clDump.readInputs.
func polymerJ2clOracleReadInputs(t *testing.T, path string) []polymerJ2clOracleInput {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var inputs []polymerJ2clOracleInput
	for i := 0; i < len(data); {
		nl := bytes.IndexByte(data[i:], '\n') + i
		parts := strings.Split(string(data[i:nl]), " ")
		if len(parts) != 3 || parts[0] != "===" {
			t.Fatalf("bad header %q", data[i:nl])
		}
		n, err := strconv.Atoi(parts[2])
		if err != nil {
			t.Fatal(err)
		}
		inputs = append(inputs, polymerJ2clOracleInput{parts[1], string(data[nl+1 : nl+1+n])})
		i = nl + 1 + n + 1
	}
	return inputs
}

// polymerJ2clOracleOptions is PolymerJ2clDump.testCaseOptions.
func polymerJ2clOracleOptions() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetEmitUseStrict(false)
	options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	options.SetModuleResolutionMode(ModuleLoader_ResolutionMode_BROWSER)
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_TYPES_ONLY)
	options.SetAssumeStaticInheritanceIsNotUsed(true)
	options.SetPreserveTypeAnnotations(true)
	options.SetAssumeGettersArePure(false)
	options.SetCheckSymbols(true)
	options.SetWarningLevel(DiagnosticGroups_INVALID_CASTS, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_MISPLACED_MSG_ANNOTATION, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_MISSING_PROPERTIES, CheckLevel_WARNING)
	options.SetCodingConvention(NewGoogleCodingConvention())
	polymerVersion := int32(1)
	options.SetPolymerVersion(&polymerVersion)
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS)
	return options
}

type polymerJ2clOracleDumper struct {
	sb strings.Builder
}

func (d *polymerJ2clOracleDumper) line(s string) {
	d.sb.WriteString(peepholeOracleEscape(s))
	d.sb.WriteByte('\n')
}

// polymerJ2clOracleError is PolymerJ2clDump.error.
func polymerJ2clOracleError(level string, e *JSError) string {
	return level + " " + e.GetType().key + " " + jutil.ToJavaString(e.GetSourceName()) + ":" +
		strconv.Itoa(e.GetLineno()) + ":" + strconv.Itoa(e.GetCharno()) + " " + e.GetDescription()
}

func (d *polymerJ2clOracleDumper) parseErrors(compiler *Compiler) {
	d.line("# parse errors")
	for _, e := range compiler.GetErrors() {
		d.line(polymerJ2clOracleError("ERROR", e))
	}
}

func (d *polymerJ2clOracleDumper) diagnostics(title string, compiler *Compiler) {
	d.line(title)
	for _, e := range compiler.GetErrors() {
		d.line(polymerJ2clOracleError("ERROR", e))
	}
	for _, e := range compiler.GetWarnings() {
		d.line(polymerJ2clOracleError("WARNING", e))
	}
}

func (d *polymerJ2clOracleDumper) tree(n *rhino.Node) {
	d.line("# tree")
	for _, l := range strings.Split(n.ToStringTree(), "\n") {
		d.line(l)
	}
}

// polymerJ2clOracleFirstLine is PolymerJ2clDump.firstLine.
func polymerJ2clOracleFirstLine(t jutil.Throwable) string {
	msg := "null"
	if m := jutil.GetMessageOrNull(t); m != nil {
		msg = *m
	}
	if nl := strings.IndexByte(msg, '\n'); nl >= 0 {
		msg = msg[:nl]
	}
	return msg
}

// polymerJ2clOracleException is PolymerJ2clDump.exception.
func (d *polymerJ2clOracleDumper) exception(prefix string, r any) {
	e := jutil.AsThrowable(r)
	text := prefix + l11msgJavaClassName(e) + ": " + polymerJ2clOracleFirstLine(e)
	if cause := e.GetCause(); cause != nil {
		name := l11msgJavaClassName(cause)
		text += " cause " + name
		// The JVM's own messages for these name its class loaders.
		if name != "java.lang.ClassCastException" && name != "java.lang.NullPointerException" {
			text += ": " + polymerJ2clOracleFirstLine(cause)
		}
	}
	d.line(text)
}

// polymerJ2clOracleRun runs body, returning what it panicked with.
func polymerJ2clOracleRun(body func()) (exception any) {
	defer func() {
		if r := recover(); r != nil {
			exception = r
		}
	}()
	body()
	return nil
}

// polymerJ2clOracleUpdateAccessorSummary is PolymerJ2clDump.updateAccessorSummary.
func polymerJ2clOracleUpdateAccessorSummary(compiler *Compiler, root *rhino.Node) {
	compiler.SetAccessorSummary(
		AccessorSummary_Create(GatherGetterAndSetterProperties_Gather(compiler, root)))
}

// polymerJ2clOracleDumpPolymer is PolymerJ2clDump.dumpPolymer.
func polymerJ2clOracleDumpPolymer(in polymerJ2clOracleInput, externs string) string {
	d := &polymerJ2clOracleDumper{}
	d.line("=== " + in.name)
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := polymerJ2clOracleOptions()
	// PolymerPassTest.getOptions()
	options.SetPreventLibraryInjection(true)
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs", externs)},
		[]*SourceFile{SourceFile_FromCode(in.name, in.code)},
		options)
	root := compiler.ParseInputs()
	if root == nil || compiler.GetErrorCount() > 0 {
		d.parseErrors(compiler)
		return d.sb.String()
	}
	externsRoot := root.GetFirstChild()
	mainRoot := root.GetLastChild()
	if e := polymerJ2clOracleRun(func() {
		// CompilerTestCase.enableCreateModuleMap()
		NewGatherModuleMetadata(compiler /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER).
			Process(externsRoot, mainRoot)
		NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).Process(externsRoot, mainRoot)
		polymerJ2clOracleUpdateAccessorSummary(compiler, root)
		changeVerifier := NewChangeVerifier(compiler).Snapshot(mainRoot)
		NewPolymerPass(compiler).Process(externsRoot, mainRoot)
		changeVerifier.CheckRecordedChanges(mainRoot)
	}); e != nil {
		d.diagnostics("# diagnostics", compiler)
		d.exception("# exception ", e)
		return d.sb.String()
	}
	d.diagnostics("# diagnostics", compiler)
	d.tree(root)
	if e := polymerJ2clOracleRun(func() {
		compilerTestCase_createTypeCheck(compiler).ProcessForTesting(externsRoot, mainRoot)
	}); e != nil {
		d.exception("# typecheck-exception ", e)
		return d.sb.String()
	}
	d.diagnostics("# typecheck-diagnostics", compiler)
	return d.sb.String()
}

// polymerJ2clOracleDumpJ2cl is PolymerJ2clDump.dumpJ2cl.
func polymerJ2clOracleDumpJ2cl(in polymerJ2clOracleInput) string {
	d := &polymerJ2clOracleDumper{}
	d.line("=== " + in.name)
	test := in.name[:strings.IndexByte(in.name, '-')]
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	if test != "J2clSourceFileCheckerTest" {
		// The tests' createCompiler()
		J2clSourceFileChecker_MarkToRunJ2clPasses(compiler)
	}
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs", "")},
		[]*SourceFile{SourceFile_FromCode(in.name, in.code)},
		polymerJ2clOracleOptions())
	root := compiler.ParseInputs()
	if root == nil || compiler.GetErrorCount() > 0 {
		d.parseErrors(compiler)
		return d.sb.String()
	}
	externsRoot := root.GetFirstChild()
	mainRoot := root.GetLastChild()
	if e := polymerJ2clOracleRun(func() {
		normalize := test == "J2clPassTest" ||
			test == "J2clClinitPrunerPassTest" ||
			test == "J2clPropertyInlinerPassTest"
		if normalize {
			Normalize_CreateNormalizeForOptimizations(compiler).Process(externsRoot, mainRoot)
		}
		polymerJ2clOracleUpdateAccessorSummary(compiler, root)
		if test == "J2clClinitPrunerPassTest" {
			NewPureFunctionIdentifier_Driver(compiler).Process(externsRoot, mainRoot)
		}
		changeVerifier := NewChangeVerifier(compiler).Snapshot(mainRoot)
		var pass CompilerPass
		switch test {
		case "J2clPassTest":
			pass = NewJ2clPass(compiler)
		case "J2clUtilGetDefineRewriterPassTest":
			pass = CompilerPassFunc(func(externs *rhino.Node, js *rhino.Node) {
				NewProcessDefines_Builder(compiler).
					SetMode(ProcessDefines_Mode_OPTIMIZE).
					Build().
					Process(externs, js)
				NewJ2clUtilGetDefineRewriterPass(compiler).Process(externs, js)
			})
		case "J2clConstantHoisterPassTest":
			pass = NewJ2clConstantHoisterPass(compiler)
		case "J2clClinitPrunerPassTest":
			pass = NewJ2clClinitPrunerPass(
				compiler,
				compiler.GetChangeTracker().GetChangedScopeNodesForPass("J2clClinitPrunerPass"))
		case "J2clPropertyInlinerPassTest":
			pass = NewJ2clPropertyInlinerPass(compiler)
		case "J2clAssertRemovalPassTest":
			pass = NewJ2clAssertRemovalPass(compiler)
		case "J2clEqualitySameRewriterPassTest":
			pass = NewPeepholeOptimizationsPass(compiler, "equality", NewJ2clEqualitySameRewriterPass(false))
		case "J2clStringValueOfRewriterPassTest":
			pass = NewPeepholeOptimizationsPass(compiler, "valueOf", NewJ2clStringValueOfRewriterPass())
		case "J2clUndefinedChecksRewriterPassTest":
			pass = NewPeepholeOptimizationsPass(compiler, "undefined", NewJ2clUndefinedChecksRewriterPass())
		case "J2clSourceFileCheckerTest":
			pass = NewJ2clSourceFileChecker(compiler)
		default:
			panic(jutil.NewIllegalArgumentExceptionMsg(test))
		}
		compiler.BeforePass(test)
		pass.Process(externsRoot, mainRoot)
		changeVerifier.CheckRecordedChanges(mainRoot)
	}); e != nil {
		d.diagnostics("# diagnostics", compiler)
		d.exception("# exception ", e)
		return d.sb.String()
	}
	d.diagnostics("# diagnostics", compiler)
	d.tree(mainRoot)
	if test == "J2clSourceFileCheckerTest" {
		d.line("# shouldRunJ2clPasses " + strconv.FormatBool(J2clSourceFileChecker_ShouldRunJ2clPasses(compiler)))
	}
	return d.sb.String()
}

// polymerJ2clOracleReadBlocks splits a dump into its "=== <name>" blocks.
func polymerJ2clOracleReadBlocks(t *testing.T, path string) map[string]string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	sc := bufio.NewScanner(gz)
	sc.Buffer(make([]byte, 1<<20), 1<<28)
	blocks := map[string]string{}
	var name string
	var cur strings.Builder
	flush := func() {
		if name != "" {
			blocks[name] = cur.String()
		}
	}
	for sc.Scan() {
		l := sc.Text()
		if strings.HasPrefix(l, "=== ") {
			flush()
			name = strings.TrimPrefix(l, "=== ")
			cur.Reset()
		}
		cur.WriteString(l + "\n")
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	flush()
	return blocks
}

func TestPolymerJ2clGolden(t *testing.T) {
	externs, err := os.ReadFile(filepath.Join(polymerJ2clOracleDir, "polymer-externs.js"))
	if err != nil {
		t.Fatal(err)
	}
	diffDir := os.Getenv("POLYMER_J2CL_DIFF_DIR")
	for _, group := range []string{"polymer", "j2cl"} {
		t.Run(group, func(t *testing.T) {
			inputs := polymerJ2clOracleReadInputs(t, filepath.Join(polymerJ2clOracleDir, group+"-inputs.txt"))
			want := polymerJ2clOracleReadBlocks(t, filepath.Join(polymerJ2clOracleDir, group+"-golden.txt.gz"))
			if len(want) != len(inputs) {
				t.Fatalf("%d golden blocks for %d inputs", len(want), len(inputs))
			}
			diffs := 0
			for _, in := range inputs {
				var got string
				if group == "polymer" {
					got = polymerJ2clOracleDumpPolymer(in, string(externs))
				} else {
					got = polymerJ2clOracleDumpJ2cl(in)
				}
				if got == want[in.name] {
					continue
				}
				diffs++
				if diffs <= 5 {
					t.Errorf("%s differs:\n%s", in.name, polymerJ2clOracleFirstDiff(want[in.name], got))
				}
				if diffDir != "" {
					_ = os.MkdirAll(diffDir, 0o755)
					_ = os.WriteFile(filepath.Join(diffDir, in.name+".go.txt"), []byte(got), 0o644)
					_ = os.WriteFile(filepath.Join(diffDir, in.name+".java.txt"), []byte(want[in.name]), 0o644)
				}
			}
			if diffs > 0 {
				t.Errorf("%s: %d of %d inputs differ", group, diffs, len(inputs))
			}
		})
	}
}

// polymerJ2clOracleFirstDiff shows the first differing line of two blocks.
func polymerJ2clOracleFirstDiff(want, got string) string {
	wl := strings.Split(want, "\n")
	gl := strings.Split(got, "\n")
	for i := 0; i < len(wl) || i < len(gl); i++ {
		var w, g string
		if i < len(wl) {
			w = wl[i]
		}
		if i < len(gl) {
			g = gl[i]
		}
		if w != g {
			return "line " + strconv.Itoa(i+1) + ":\n  java: " + w + "\n  go:   " + g
		}
	}
	return "(identical lines)"
}
