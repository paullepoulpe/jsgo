package jscomp

// Differential oracle for the level-11 checks, Go side of
// tools/oracle/java/ChecksDump.java (see its header for the format): each
// input is compiled the way CompilerTestCase does by default, parsed with
// Compiler.ParseInputs, and then CheckDebuggerStatement, CheckGlobalThis,
// CheckSuspiciousCode, CheckRegExp, CheckSuper,
// CheckTypeImportCodeReferences, StrictModeCheck, SourceInfoCheck,
// CheckUnreachableCode, CheckMissingReturn, J2clChecksPass, VarCheck (plus
// the synthesized externs it creates), ConstCheck, VarCheck(validityCheck)
// and ForbidDynamicImportUsage run in order on the same compiler; every
// diagnostic that reaches the ErrorManager is printed in report order.
//
//   - TestChecksDumpGolden compares against the checked-in Java output for
//     the sources in testdata/checks/src (always runs):
//     testdata/checks/expected-code.txt (each file as the code input) and
//     testdata/checks/expected-externs.txt (each file as the externs input,
//     --as-externs). Regenerate with
//       java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.ChecksDump \
//         [--as-externs] testdata/checks/expected-<mode>.txt testdata/checks/src <rel>...
//     testdata/checks/expected-program.txt is the --program run over
//     testdata/checks/program (one compilation; CheckMissingRequires first).
//   - TestChecksDumpOracle compares a whole run when CHECKS_ORACLE (the
//     ChecksDump output, .gz or plain), CHECKS_ROOT (the root of the sources)
//     and CHECKS_FILES (the list of relative paths given to ChecksDump) are
//     set; CHECKS_AS_EXTERNS=1 selects --as-externs, CHECKS_PROGRAM=1
//     --program, and CHECKS_DIFF_DIR, if set, receives the Go and Java
//     sections of every mismatching file. Use absolute paths: `go test` runs
//     in the package directory.

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

// checksDumpExterns is ChecksDump.EXTERNS (VarCheckTest's externs).
const checksDumpExterns = "var window; function alert() {}"

// checksDumpOptions is ChecksDump.options(): CompilerTestCase.getOptions()
// with its setUp() defaults.
func checksDumpOptions() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	options.SetEmitUseStrict(false)
	options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	options.SetModuleResolutionMode(ModuleLoader_ResolutionMode_BROWSER)
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_TYPES_ONLY)
	options.SetAssumeStaticInheritanceIsNotUsed(true)
	options.SetPreserveTypeAnnotations(true)
	options.SetAssumeGettersArePure(false)
	options.SetCheckSymbols(true)
	options.SetCodingConvention(NewGoogleCodingConvention())
	options.SetCheckTypes(false)
	options.SetWarningLevel(
		NewDiagnosticGroup(
			CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT,
			VarCheck_STRICT_CHUNK_DEP_ERROR,
			CheckMissingRequires_MISSING_REQUIRE_TYPE,
			CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE,
			CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE,
			CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE,
			CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE,
			CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE),
		CheckLevel_WARNING)
	options.SetWarningLevel(
		NewDiagnosticGroup(ModuleLoader_LOAD_WARNING, ModuleLoader_INVALID_MODULE_PATH),
		CheckLevel_OFF)
	return options
}

// checksDumpErrorManager is ChecksDump.RecordingErrorManager: it records
// every diagnostic that reaches the error manager, in report order.
type checksDumpErrorManager struct {
	BasicErrorManagerBase
	sb *strings.Builder
}

func newChecksDumpErrorManager(sb *strings.Builder) *checksDumpErrorManager {
	m := &checksDumpErrorManager{sb: sb}
	m.InitBasicErrorManagerBase(m)
	return m
}

func (m *checksDumpErrorManager) Report(level *CheckLevel, error *JSError) {
	checksDumpLine(m.sb,
		jutil.ToJavaString(level)+
			" "+
			error.GetType().key+
			" "+
			jutil.ToJavaString(error.GetSourceName())+
			":"+
			strconv.Itoa(error.GetLineno())+
			":"+
			strconv.Itoa(error.GetCharno())+
			" @"+
			strconv.Itoa(error.GetNodeSourceOffset())+
			"+"+
			strconv.Itoa(error.GetLength())+
			" "+
			error.GetDescription())
	m.BasicErrorManagerBase.Report(level, error)
}

func (m *checksDumpErrorManager) Println(level *CheckLevel, error *JSError) {}

func (m *checksDumpErrorManager) PrintSummary() {}

type checksDumpPass func(c *Compiler, externsRoot *rhino.Node, mainRoot *rhino.Node, out *strings.Builder)

// checksDumpOne is ChecksDump.dumpOne.
func checksDumpOne(sb *strings.Builder, rel string, src string, asExterns bool) {
	checksDumpLine(sb, "# file "+rel)
	if asExterns {
		checksDump(sb,
			[]*SourceFile{SourceFile_FromCode(rel, src)},
			[]*SourceFile{SourceFile_FromCode("code.js", "")},
			false)
	} else {
		checksDump(sb,
			[]*SourceFile{SourceFile_FromCode("externs", checksDumpExterns)},
			[]*SourceFile{SourceFile_FromCode(rel, src)},
			false)
	}
}

// checksDumpProgram is ChecksDump.main with --program: all the files are the
// code inputs of one compilation.
func checksDumpProgram(sb *strings.Builder, root string, files []string) error {
	var inputs []*SourceFile
	for _, rel := range files {
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			return err
		}
		inputs = append(inputs, SourceFile_FromCode(rel, string(src)))
	}
	checksDumpLine(sb, "# program "+strconv.Itoa(len(files))+" files")
	checksDump(sb, []*SourceFile{SourceFile_FromCode("externs", checksDumpExterns)}, inputs, true)
	return nil
}

// checksDump is ChecksDump.dump.
func checksDump(sb *strings.Builder, externs []*SourceFile, inputs []*SourceFile, program bool) {
	compiler := NewCompiler()
	compiler.SetErrorManager(newChecksDumpErrorManager(sb))
	compiler.Init(externs, inputs, checksDumpOptions())
	checksDumpLine(sb, "## parse")
	var root *rhino.Node
	threw := false
	jutil.Try(func() {
		root = compiler.ParseInputs()
	}, func(e jutil.IRuntimeException) {
		checksDumpException(sb, e)
		threw = true
	})
	if threw {
		return
	}
	if root == nil || compiler.GetErrorCount() > 0 {
		checksDumpLine(sb, "# parse-errors")
		return
	}
	externsRoot := root.GetFirstChild()
	mainRoot := root.GetLastChild()

	run := func(name string, pass checksDumpPass) {
		checksDumpLine(sb, "## "+name)
		jutil.Try(func() {
			pass(compiler, externsRoot, mainRoot, sb)
		}, func(e jutil.IRuntimeException) {
			checksDumpException(sb, e)
		})
	}

	if program {
		// As CheckMissingRequiresTest does.
		run("CheckMissingRequires", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
			NewGatherModuleMetadata(c /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER).
				Process(e, m)
			NewCheckMissingRequires(c, c.GetModuleMetadataMap()).Process(e, m)
		})
	}
	run("CheckDebuggerStatement", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewCheckDebuggerStatement(c).Process(e, m)
	})
	run("CheckGlobalThis", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NodeTraversal_Traverse(c, m, NewCheckGlobalThis(c))
	})
	run("CheckSuspiciousCode", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NodeTraversal_Traverse(c, m, NewCheckSuspiciousCode())
	})
	run("CheckRegExp", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		pass := NewCheckRegExp(c, true, true)
		pass.Process(e, m)
		checksDumpLine(out, "= globalRegExpPropertiesUsed="+strconv.FormatBool(pass.IsGlobalRegExpPropertiesUsed()))
	})
	run("CheckSuper", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewCheckSuper(c).Process(e, m)
	})
	run("CheckTypeImportCodeReferences", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewCheckTypeImportCodeReferences(c).Process(e, m)
	})
	run("StrictModeCheck", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewStrictModeCheck(c, CheckLevel_WARNING).Process(e, m)
	})
	run("SourceInfoCheck", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewSourceInfoCheck(c).Process(e, m)
	})
	run("CheckUnreachableCode", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NodeTraversal_Traverse(c, m, NewCheckUnreachableCode(c))
	})
	run("CheckMissingReturn", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NodeTraversal_Traverse(c, m, NewCheckMissingReturn(c))
	})
	run("J2clChecksPass", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewJ2clChecksPass(c).Process(e, m)
	})
	run("VarCheck", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewVarCheck(c).Process(e, m)
		synthesized := c.GetSynthesizedExternsInput().GetAstRoot(c)
		for _, l := range strings.Split(synthesized.ToStringTree(), "\n") {
			if l != "" {
				checksDumpLine(out, "= "+l)
			}
		}
	})
	run("ConstCheck", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewConstCheck(c).Process(e, m)
	})
	run("VarCheckValidity", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewVarCheckWithValidityCheck(c, true).Process(e, m)
	})
	run("ForbidDynamicImportUsage", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewForbidDynamicImportUsage(c).Process(e, m)
	})
}

// checksDumpException is ChecksDump.exception.
func checksDumpException(sb *strings.Builder, e jutil.Throwable) {
	checksDumpLine(sb, "! "+checksDumpClassName(e))
	for c := e.GetCause(); c != nil && c != c.GetCause(); c = c.GetCause() {
		if _, ok := c.(jutil.INullPointerException); ok {
			checksDumpLine(sb, "!! cause "+checksDumpClassName(c))
		} else {
			checksDumpLine(sb, "!! cause "+checksDumpClassName(c)+": "+jutil.ToJavaString(jutil.GetMessageOrNull(c)))
		}
	}
}

// checksDumpClassName is getClass().getName(): Throwable.Error() renders
// "<class>: <message>" (or "<class>" without a message).
func checksDumpClassName(e jutil.Throwable) string {
	s := e.Error()
	if m := jutil.GetMessageOrNull(e); m != nil {
		return strings.TrimSuffix(s, ": "+*m)
	}
	return s
}

// checksDumpEscape is ChecksDump.escape: backslash and every UTF-16 code
// unit outside printable ASCII become escapes.
func checksDumpEscape(s string) string {
	var b strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			b.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(byte(c))
		default:
			fmt.Fprintf(&b, "\\u%04x", c)
		}
	}
	return b.String()
}

func checksDumpLine(sb *strings.Builder, s string) {
	sb.WriteString(checksDumpEscape(s))
	sb.WriteByte('\n')
}

// checksDumpSections splits a ChecksDump output into its per-file sections,
// keyed by relative path, in file order.
func checksDumpSections(r io.Reader) ([]string, map[string]string, error) {
	var order []string
	sections := map[string]string{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1<<20), 1<<30)
	var cur string
	var b strings.Builder
	flush := func() {
		if cur != "" {
			sections[cur] = b.String()
		}
		b.Reset()
	}
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "# file ") {
			flush()
			cur = strings.TrimPrefix(line, "# file ")
			order = append(order, cur)
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	flush()
	return order, sections, sc.Err()
}

func checksDumpOpen(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(path, ".gz") {
		return f, nil
	}
	zr, err := gzip.NewReader(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	return struct {
		io.Reader
		io.Closer
	}{zr, f}, nil
}

// checksDumpCompare runs the Go side over files and compares it with the
// Java output in oraclePath; it returns the number of mismatching files.
func checksDumpCompare(t *testing.T, oraclePath string, root string, files []string, asExterns bool, diffDir string) int {
	t.Helper()
	rc, err := checksDumpOpen(oraclePath)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	_, want, err := checksDumpSections(rc)
	if err != nil {
		t.Fatal(err)
	}
	mismatches := 0
	for _, rel := range files {
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		var sb strings.Builder
		checksDumpOne(&sb, rel, string(src), asExterns)
		got := sb.String()
		if got == want[rel] {
			continue
		}
		mismatches++
		if mismatches <= 10 {
			t.Errorf("%s: Go output differs from Java:\n%s", rel, checksDumpFirstDiff(want[rel], got))
		}
		if diffDir != "" {
			out := filepath.Join(diffDir, rel)
			if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(out+".go.txt", []byte(got), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(out+".java.txt", []byte(want[rel]), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return mismatches
}

func checksDumpFirstDiff(want, got string) string {
	w := strings.Split(want, "\n")
	g := strings.Split(got, "\n")
	for i := 0; i < len(w) || i < len(g); i++ {
		var wl, gl string
		if i < len(w) {
			wl = w[i]
		}
		if i < len(g) {
			gl = g[i]
		}
		if wl != gl {
			return fmt.Sprintf("line %d:\n  java: %q\n  go:   %q", i+1, wl, gl)
		}
	}
	return "(no line difference)"
}

// checksDumpReadList reads a file list, without the files of
// tools/oracle/gpl-excluded-corpus.txt.
func checksDumpReadList(t *testing.T, path string) []string {
	t.Helper()
	files, err := corpus.ReadList(path)
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// checksDumpListDir lists the files of dir in name order (the order of the
// Java side's sorted list).
func checksDumpListDir(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var files []string
	for _, e := range entries {
		files = append(files, e.Name())
	}
	return files
}

// checksDumpCompareProgram compares a --program run.
func checksDumpCompareProgram(t *testing.T, oraclePath string, root string, files []string, diffDir string) {
	t.Helper()
	rc, err := checksDumpOpen(oraclePath)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	want, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	var sb strings.Builder
	if err := checksDumpProgram(&sb, root, files); err != nil {
		t.Fatal(err)
	}
	if got := sb.String(); got != string(want) {
		t.Errorf("Go output differs from Java:\n%s", checksDumpFirstDiff(string(want), got))
		if diffDir != "" {
			if err := os.WriteFile(filepath.Join(diffDir, "program.go.txt"), []byte(got), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestChecksDumpGolden(t *testing.T) {
	dir := filepath.Join("testdata", "checks")
	files := checksDumpListDir(t, filepath.Join(dir, "src"))
	t.Run("program", func(t *testing.T) {
		checksDumpCompareProgram(t, filepath.Join(dir, "expected-program.txt"), filepath.Join(dir, "program"),
			checksDumpListDir(t, filepath.Join(dir, "program")), "")
	})
	for _, mode := range []struct {
		name      string
		asExterns bool
	}{{"code", false}, {"externs", true}} {
		t.Run(mode.name, func(t *testing.T) {
			if n := checksDumpCompare(t, filepath.Join(dir, "expected-"+mode.name+".txt"),
				filepath.Join(dir, "src"), files, mode.asExterns, ""); n != 0 {
				t.Errorf("%d of %d files differ", n, len(files))
			}
		})
	}
}

func TestChecksDumpOracle(t *testing.T) {
	oracle, root, list := os.Getenv("CHECKS_ORACLE"), os.Getenv("CHECKS_ROOT"), os.Getenv("CHECKS_FILES")
	if oracle == "" || root == "" || list == "" {
		t.Skip("set CHECKS_ORACLE, CHECKS_ROOT and CHECKS_FILES to compare a ChecksDump run")
	}
	files := checksDumpReadList(t, list)
	if os.Getenv("CHECKS_PROGRAM") == "1" {
		checksDumpCompareProgram(t, oracle, root, files, os.Getenv("CHECKS_DIFF_DIR"))
		return
	}
	n := checksDumpCompare(t, oracle, root, files, os.Getenv("CHECKS_AS_EXTERNS") == "1", os.Getenv("CHECKS_DIFF_DIR"))
	if n != 0 {
		t.Errorf("%d of %d files differ", n, len(files))
	}
}
