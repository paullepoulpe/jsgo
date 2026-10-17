package jscomp

// Go-only: test-helper — the differential oracle for CheckJSDoc,
// CheckJSDocStyle, VariableReferenceCheck, CheckSideEffects (+
// StripProtection), CheckAccessControls and CheckConformance (with
// ConformanceRules): the Go mirror of tools/oracle/java/W5ChecksDump.java,
// whose header documents the options, the steps (DefaultPassConfig.getChecks()
// order, on type-checked code) and the output format. The diagnostic lines,
// exceptions, escaping and section splitting are ChecksDump's
// (checks_oracle_test.go).
//
//   - TestW5ChecksDumpGolden compares against the checked-in Java output
//     (always runs):
//     testdata/w5_checks/expected-src.txt for the hand-written inputs in
//     testdata/w5_checks/src (multi-file inputs split at "//// file: " lines)
//     and testdata/w5_checks/expected-typecheck.txt.gz for the corpus sample
//     of testdata/typecheck (files.txt, sources in src/). Regenerate with
//     (TD=internal/jscomp/testdata/w5_checks)
//       java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//         com.google.javascript.jscomp.W5ChecksDump \
//         --externs=tools/oracle/typecheck-externs.js --externs=$TD/externs-extra.js \
//         --conformance=example:$TD/conformance-example.textproto \
//         --conformance=all:$TD/conformance-all.textproto \
//         $TD/expected-src.txt $TD/src @<sorted list of $TD/src>
//     and the same with $TD/expected-typecheck.txt.gz,
//     internal/jscomp/testdata/typecheck/src and
//     @internal/jscomp/testdata/typecheck/files.txt.
//   - TestW5ChecksDumpOracle compares a whole run when W5CHECKS_ORACLE (the
//     W5ChecksDump output, .gz or plain), W5CHECKS_ROOT (the root of the
//     sources) and W5CHECKS_FILES (the list given to W5ChecksDump) are set;
//     W5CHECKS_DIFF_DIR, if set, receives <rel>.go.txt and <rel>.java.txt
//     for every mismatching file; W5CHECKS_NOCATCH=1 makes an exception
//     thrown by a step fail the test with its Go stack. The externs and conformance configs are
//     always the ones above. Use absolute paths: `go test` runs in the
//     package directory.
//
// The diagnostics of each CheckConformance step are sorted (see the Java
// header: Java's order between rules depends on identity hash codes).
//
// Both report, besides the mismatching files, how many files differ in each
// step ("## <name>" section), so that a type-checking difference is told
// apart from a difference in the checks.
//
// Known, accepted difference: Java catches StackOverflowError (the tool runs
// with -Xss1g); a Go stack overflow is fatal, so such an input stops the run.

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// w5ChecksDump_noCatch (W5CHECKS_NOCATCH=1) lets an exception thrown by a
// step propagate, with its Go stack, instead of printing it: a debugging aid.
var w5ChecksDump_noCatch = os.Getenv("W5CHECKS_NOCATCH") == "1"

// w5ChecksDump_marker is W5ChecksDump.MARKER.
const w5ChecksDump_marker = "//// file: "

// w5ChecksDump_namedConfig is W5ChecksDump.NamedConfig.
type w5ChecksDump_namedConfig struct {
	name   string
	config *ConformanceConfig
}

// w5ChecksDump_setup is what W5ChecksDump.main reads from its --externs and
// --conformance arguments: the externs SourceFiles (each named by its file
// name) and the parsed conformance configs.
type w5ChecksDump_setup struct {
	externs []*SourceFile
	configs []w5ChecksDump_namedConfig
}

// w5ChecksDump_loadSetup reads the externs and configs every run of the
// oracle uses (see the Java header).
func w5ChecksDump_loadSetup(t testing.TB) *w5ChecksDump_setup {
	t.Helper()
	td := filepath.Join("testdata", "w5_checks")
	s := &w5ChecksDump_setup{}
	for _, p := range []string{
		filepath.Join("..", "..", "tools", "oracle", "typecheck-externs.js"),
		filepath.Join(td, "externs-extra.js"),
	} {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		s.externs = append(s.externs, SourceFile_FromCode(filepath.Base(p), string(b)))
	}
	for _, c := range []struct{ name, file string }{
		{"example", "conformance-example.textproto"},
		{"all", "conformance-all.textproto"},
	} {
		b, err := os.ReadFile(filepath.Join(td, c.file))
		if err != nil {
			t.Fatal(err)
		}
		builder := ConformanceConfig_NewBuilder()
		TextFormat_Merge(string(b), builder)
		s.configs = append(s.configs, w5ChecksDump_namedConfig{c.name, builder.Build()})
	}
	return s
}

// w5ChecksDump_options is W5ChecksDump.options().
func w5ChecksDump_options() *CompilerOptions {
	options := NewCompilerOptions()
	// TypeCheckDump.options()
	options.SetCodingConvention(NewGoogleCodingConvention())
	options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	options.SetWarningLevel(DiagnosticGroups_MISSING_PROPERTIES, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_MISPLACED_TYPE_ANNOTATION, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_INVALID_CASTS, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_JSDOC_MISSING_TYPE, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_BOUNDED_GENERICS, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_MISSING_OVERRIDE, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_STRICT_MISSING_PROPERTIES, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS, CheckLevel_WARNING)
	options.SetCheckTypes(true)
	// The checks of this tool.
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE)
	options.SetWarningLevel(DiagnosticGroups_MISPLACED_SUPPRESS, CheckLevel_WARNING)
	options.SetWarningLevel(CheckJSDocStyle_ALL_DIAGNOSTICS, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_ACCESS_CONTROLS, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_DEPRECATED, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_CONSTANT_PROPERTY, CheckLevel_WARNING)
	options.SetWarningLevel(
		NewDiagnosticGroup(CheckJSDoc_MISPLACED_MSG_ANNOTATION, VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT),
		CheckLevel_WARNING)
	options.SetWarningLevel(
		NewDiagnosticGroup(ModuleLoader_LOAD_WARNING, ModuleLoader_INVALID_MODULE_PATH),
		CheckLevel_OFF)
	return options
}

// w5ChecksDump_inputs is W5ChecksDump.inputs.
func w5ChecksDump_inputs(rel string, src string) []*SourceFile {
	if !strings.HasPrefix(src, w5ChecksDump_marker) {
		return []*SourceFile{SourceFile_FromCode(rel, src)}
	}
	var inputs []*SourceFile
	var name *string
	var cur strings.Builder
	for _, l := range strings.Split(src, "\n") {
		if strings.HasPrefix(l, w5ChecksDump_marker) {
			if name != nil {
				inputs = append(inputs, SourceFile_FromCode(*name, cur.String()))
			}
			n := jutil.StringTrim(strings.TrimPrefix(l, w5ChecksDump_marker))
			name = &n
			cur.Reset()
		} else {
			cur.WriteString(l)
			cur.WriteByte('\n')
		}
	}
	inputs = append(inputs, SourceFile_FromCode(*name, cur.String()))
	return inputs
}

// w5ChecksDump_dumpOne is W5ChecksDump.dumpOne.
func w5ChecksDump_dumpOne(sb *strings.Builder, setup *w5ChecksDump_setup, rel string, src string) {
	checksDumpLine(sb, "# file "+rel)
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	em := newChecksDumpErrorManager(sb)
	compiler.SetErrorManager(em)
	options := w5ChecksDump_options()
	compiler.Init(setup.externs, w5ChecksDump_inputs(rel, src), options)
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

	run := func(name string, step checksDumpPass) {
		checksDumpLine(sb, "## "+name)
		if w5ChecksDump_noCatch {
			step(compiler, externsRoot, mainRoot, sb)
			return
		}
		jutil.Try(func() {
			step(compiler, externsRoot, mainRoot, sb)
		}, func(e jutil.IRuntimeException) {
			checksDumpException(sb, e)
		})
	}

	run("modules", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewGatherModuleMetadata(c /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER).
			Process(e, m)
		NewModuleMapCreator(c, c.GetModuleMetadataMap()).Process(e, m)
	})
	run("CheckJSDoc", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NodeTraversal_Traverse(c, m, NewCheckJSDoc(c))
	})
	run("CheckJSDocStyle", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NodeTraversal_Traverse(c, m, NewCheckJSDocStyle(c))
	})
	run("VariableReferenceCheck", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewVariableReferenceCheck(c).Process(e, m)
	})
	run("CheckSideEffects", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewCheckSideEffects(c, true, true).Process(e, m)
		w5ChecksDump_protectors(out, m)
		synthesized := c.GetSynthesizedExternsInput().GetAstRoot(c)
		for _, l := range strings.Split(synthesized.ToStringTree(), "\n") {
			if l != "" {
				checksDumpLine(out, "= extern "+l)
			}
		}
	})
	run("TypeCheck", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewInferConsts(c).Process(e, m)
		// DefaultPassConfig.inferTypes
		inferencePass :=
			NewTypeInferencePass(
				c, c.GetReverseAbstractInterpreter(), c.GetTypedScopeCreator().(*TypedScopeCreator))
		c.SetTypeCheckingHasRun(true)
		c.SetTopScope(inferencePass.InferAllScopes(m.GetParent()))
		// DefaultPassConfig.checkTypes
		check :=
			NewTypeCheckWithScope(
				c,
				c.GetReverseAbstractInterpreter(),
				c.GetTypeRegistry(),
				c.GetTopScope(),
				c.GetTypedScopeCreator().(*TypedScopeCreator)).
				ReportUnknownTypes(
					options.Enables(DiagnosticGroup_ForType(TypeCheck_UNKNOWN_EXPR_TYPE))).
				ReportMissingProperties(
					!options.Disables(DiagnosticGroup_ForType(TypeCheck_INEXISTENT_PROPERTY)))
		check.Process(e, m)
		c.GetErrorManager().SetTypedPercent(check.GetTypedPercent())
	})
	run("CheckAccessControls", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewCheckAccessControls(c).Process(e, m)
	})
	for _, config := range setup.configs {
		run("CheckConformance "+config.name, func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
			var diagnostics strings.Builder
			em.sb = &diagnostics
			defer func() {
				em.sb = out
				var lines []string
				for _, l := range strings.Split(diagnostics.String(), "\n") {
					if l != "" {
						lines = append(lines, l)
					}
				}
				// The lines are escaped to printable ASCII: byte order is
				// String order.
				sort.Strings(lines)
				for _, l := range lines {
					out.WriteString(l)
					out.WriteByte('\n')
				}
			}()
			reportingMode := options.GetConformanceReportingMode()
			NewCheckConformanceWithConfigs(c, []*ConformanceConfig{config.config}, &reportingMode).Process(e, m)
		})
	}
	run("StripProtection", func(c *Compiler, e, m *rhino.Node, out *strings.Builder) {
		NewCheckSideEffects_StripProtection(c).Process(e, m)
		checksDumpLine(out, "= protectors "+strconv.Itoa(w5ChecksDump_countProtectors(m)))
	})
}

// w5ChecksDump_isProtector is W5ChecksDump.isProtector.
func w5ChecksDump_isProtector(n *rhino.Node) bool {
	return n.IsCall() &&
		n.GetFirstChild().IsName() &&
		n.GetFirstChild().GetString() == "JSCOMPILER_PRESERVE"
}

// w5ChecksDump_protectors is W5ChecksDump.protectors.
func w5ChecksDump_protectors(out *strings.Builder, n *rhino.Node) {
	if w5ChecksDump_isProtector(n) {
		checksDumpLine(out, "= protector "+jutil.ToJavaString(n.GetSourceFileName())+":"+
			strconv.Itoa(n.GetLineno())+":"+strconv.Itoa(n.GetCharno()))
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		w5ChecksDump_protectors(out, c)
	}
}

// w5ChecksDump_countProtectors is W5ChecksDump.countProtectors.
func w5ChecksDump_countProtectors(n *rhino.Node) int {
	count := 0
	if w5ChecksDump_isProtector(n) {
		count = 1
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		count += w5ChecksDump_countProtectors(c)
	}
	return count
}

// w5ChecksDump_steps splits one file's section of the output into its steps
// ("## <name>" sub-sections), keyed by name.
func w5ChecksDump_steps(section string) map[string]string {
	steps := map[string]string{}
	cur := ""
	for _, l := range strings.SplitAfter(section, "\n") {
		if strings.HasPrefix(l, "## ") {
			cur = strings.TrimSuffix(strings.TrimPrefix(l, "## "), "\n")
		}
		steps[cur] += l
	}
	return steps
}

// w5ChecksDump_compare runs the Go side over files and compares it with the
// Java output in oraclePath, file by file and step by step.
func w5ChecksDump_compare(t *testing.T, setup *w5ChecksDump_setup, oraclePath string, root string, files []string, diffDir string) {
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
	stepMismatches := map[string]int{}
	for _, rel := range files {
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		var sb strings.Builder
		w5ChecksDump_dumpOne(&sb, setup, rel, string(src))
		got := sb.String()
		if got == want[rel] {
			continue
		}
		mismatches++
		gotSteps, wantSteps := w5ChecksDump_steps(got), w5ChecksDump_steps(want[rel])
		for name := range wantSteps {
			if gotSteps[name] != wantSteps[name] {
				stepMismatches[name]++
			}
		}
		for name := range gotSteps {
			if _, ok := wantSteps[name]; !ok {
				stepMismatches[name]++
			}
		}
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
	if mismatches != 0 {
		var names []string
		for name := range stepMismatches {
			names = append(names, name)
		}
		sort.Strings(names)
		var b strings.Builder
		for _, name := range names {
			b.WriteString("\n  " + name + ": " + strconv.Itoa(stepMismatches[name]))
		}
		t.Errorf("%d of %d files differ; files differing per step:%s", mismatches, len(files), b.String())
	}
}

func TestW5ChecksDumpGolden(t *testing.T) {
	setup := w5ChecksDump_loadSetup(t)
	td := filepath.Join("testdata", "w5_checks")
	t.Run("src", func(t *testing.T) {
		var files []string
		srcRoot := filepath.Join(td, "src")
		err := filepath.WalkDir(srcRoot, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			rel, err := filepath.Rel(srcRoot, p)
			files = append(files, filepath.ToSlash(rel))
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
		sort.Strings(files)
		w5ChecksDump_compare(t, setup, filepath.Join(td, "expected-src.txt"), srcRoot, files, "")
	})
	t.Run("typecheck", func(t *testing.T) {
		dir := filepath.Join("testdata", "typecheck")
		w5ChecksDump_compare(t, setup, filepath.Join(td, "expected-typecheck.txt.gz"),
			filepath.Join(dir, "src"), checksDumpReadList(t, filepath.Join(dir, "files.txt")), "")
	})
}

func TestW5ChecksDumpOracle(t *testing.T) {
	oracle, root, list := os.Getenv("W5CHECKS_ORACLE"), os.Getenv("W5CHECKS_ROOT"), os.Getenv("W5CHECKS_FILES")
	if oracle == "" || root == "" || list == "" {
		t.Skip("set W5CHECKS_ORACLE, W5CHECKS_ROOT and W5CHECKS_FILES to compare a W5ChecksDump run")
	}
	w5ChecksDump_compare(t, w5ChecksDump_loadSetup(t), oracle, root, checksDumpReadList(t, list), os.Getenv("W5CHECKS_DIFF_DIR"))
}
