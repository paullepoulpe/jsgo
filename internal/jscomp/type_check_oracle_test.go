package jscomp

// Go-only: test-helper — the differential oracle for TypeCheck, TypeInference
// and TypeInferencePass (and, through them, TypedScopeCreator,
// FunctionTypeBuilder and TypeValidator): the Go mirror of
// tools/oracle/java/TypeCheckDump.java, whose header documents the pipeline
// and the output format. Each input is type-checked the way
// DefaultPassConfig's inferTypes and checkTypes passes do (after
// GatherModuleMetadata, ModuleMapCreator and InferConsts, as
// TypeCheckTestCase prepares them), with CompilerTypeTestCase's options and
// externs (tools/oracle/typecheck-externs.js), and every diagnostic, the
// typed percentage and the main AST with the JSType of every node are
// printed.
//
//   - TestTypeCheckDumpGolden compares against the checked-in Java output
//     testdata/typecheck/expected/<rel>.typecheck.gz for every <rel> of
//     testdata/typecheck/files.txt (inputs in testdata/typecheck/src, a
//     deterministic sample of the corpus; always runs). Regenerate with
//       java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//         com.google.javascript.jscomp.TypeCheckDump \
//         --gzip --externs=tools/oracle/typecheck-externs.js internal/jscomp/testdata/typecheck/expected \
//         internal/jscomp/testdata/typecheck/src @internal/jscomp/testdata/typecheck/files.txt
//   - TestTypeCheckDumpOracle compares a whole run when TYPECHECK_ROOT (the
//     root of the sources), TYPECHECK_ORACLE_DIR (the Java <outdir>) and
//     TYPECHECK_FILES (the list of relative paths given to TypeCheckDump)
//     are set; TYPECHECK_EVERY=<n> checks every n-th file only, and
//     TYPECHECK_DIFF_DIR, if set, receives <rel>.go.txt and <rel>.java.txt
//     for every mismatching file. Use absolute paths: `go test` runs in the
//     package directory.
//
// Known, accepted difference: Java catches StackOverflowError (the tool runs
// with -Xss1g); a Go stack overflow is fatal, so such an input stops the run.

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

// typeCheckDump_options is TypeCheckDump.options():
// CompilerTypeTestCase.defaultOptions() + TypeCheckTestCase.setUp's levels +
// checkTypes.
func typeCheckDump_options() *CompilerOptions {
	options := NewCompilerOptions()
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
	return options
}

// typeCheckDump_errorManager is TypeCheckDump.RecordingErrorManager: it
// records every report(level, error) in order; SortingErrorManager keeps the
// counts.
type typeCheckDump_errorManager struct {
	SortingErrorManager
	lines []string
}

func newTypeCheckDump_errorManager() *typeCheckDump_errorManager {
	m := &typeCheckDump_errorManager{}
	m.InitSortingErrorManager(nil)
	return m
}

func (m *typeCheckDump_errorManager) Report(level *CheckLevel, error_ *JSError) {
	m.lines = append(m.lines,
		jutil.ToJavaString(level)+
			" "+
			error_.GetType().key+
			" "+
			jutil.ToJavaString(error_.GetSourceName())+
			":"+
			strconv.Itoa(error_.GetLineno())+
			":"+
			strconv.Itoa(error_.GetCharno())+
			" \""+
			error_.GetDescription()+
			"\"")
	m.SortingErrorManager.Report(level, error_)
}

// typeCheckDump_dumpOne is TypeCheckDump.dumpOne.
func typeCheckDump_dumpOne(sb *strings.Builder, externs string, name string, src string) {
	errors := newTypeCheckDump_errorManager()
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	compiler.SetErrorManager(errors)
	options := typeCheckDump_options()
	var typedPercent *float64
	var thrown jutil.Throwable
	jutil.Try(func() {
		compiler.Init(
			[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
			[]*SourceFile{SourceFile_FromCode(name, src)},
			options)
		compiler.Parse()
		if errors.GetErrorCount() == 0 {
			externsRoot := compiler.GetExternsRoot()
			jsRoot := compiler.GetJsRoot()
			NewGatherModuleMetadata(compiler, false, ModuleLoader_ResolutionMode_BROWSER).
				Process(externsRoot, jsRoot)
			NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).
				Process(externsRoot, jsRoot)
			NewInferConsts(compiler).Process(externsRoot, jsRoot)

			// DefaultPassConfig.inferTypes
			globalRoot := jsRoot.GetParent()
			inferencePass :=
				NewTypeInferencePass(
					compiler,
					compiler.GetReverseAbstractInterpreter(),
					compiler.GetTypedScopeCreator().(*TypedScopeCreator))
			compiler.SetTypeCheckingHasRun(true)
			compiler.SetTopScope(inferencePass.InferAllScopes(globalRoot))

			// DefaultPassConfig.checkTypes
			check :=
				NewTypeCheckWithScope(
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
			percent := check.GetTypedPercent()
			typedPercent = &percent
			compiler.GetErrorManager().SetTypedPercent(percent)
		}
	}, func(e jutil.IRuntimeException) {
		thrown = e
	})
	typeCheckDump_line(sb, "# diagnostics")
	for _, d := range errors.lines {
		typeCheckDump_line(sb, d)
	}
	if thrown != nil {
		typeCheckDump_line(sb, "# exception "+typeCheckDump_className(thrown)+": "+
			jutil.ToJavaString(jutil.GetMessageOrNull(thrown)))
		return
	}
	if typedPercent != nil {
		typeCheckDump_line(sb, "# typedPercent "+jutil.DoubleToString(*typedPercent))
	}
	typeCheckDump_line(sb, "# tree")
	jsRoot := compiler.GetJsRoot()
	if jsRoot.IsNil() {
		typeCheckDump_line(sb, "null")
	} else {
		typeCheckDump_dumpTree(sb, jsRoot, 0)
	}
}

// typeCheckDump_dumpTree is TypeCheckDump.dumpTree.
func typeCheckDump_dumpTree(sb *strings.Builder, n rhino.Node, level int) {
	typeCheckDump_line(sb, strings.Repeat("    ", level)+n.ToStringWithFlags(true, true, true))
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		typeCheckDump_dumpTree(sb, c, level+1)
	}
}

// typeCheckDump_className is thrown.getClass().getName(): Throwable.Error()
// renders "<class>: <message>" (or "<class>" without a message).
func typeCheckDump_className(e jutil.Throwable) string {
	s := e.Error()
	if m := jutil.GetMessageOrNull(e); m != nil {
		return strings.TrimSuffix(s, ": "+*m)
	}
	return s
}

// typeCheckDump_escape is TypeCheckDump.escape: backslash and every UTF-16
// code unit outside printable ASCII become escapes (as ParsingDump).
func typeCheckDump_escape(s string) string {
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

// typeCheckDump_line is TypeCheckDump.line.
func typeCheckDump_line(sb *strings.Builder, s string) {
	sb.WriteString(typeCheckDump_escape(s))
	sb.WriteByte('\n')
}

// typeCheckDump_readGz reads a .typecheck.gz file of the Java side.
func typeCheckDump_readGz(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(zr)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// typeCheckDump_readList reads a list of relative paths (one per line, as
// TypeCheckDump's @<listfile>), without the files of
// tools/oracle/gpl-excluded-corpus.txt.
func typeCheckDump_readList(t *testing.T, path string) []string {
	t.Helper()
	files, err := corpus.ReadList(path)
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// typeCheckDump_firstDiff describes the first differing line.
func typeCheckDump_firstDiff(want, got string) string {
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

// typeCheckDump_compare runs the Go side over files (every every-th one) and
// compares it with the Java outputs in oracleDir; it returns the number of
// files compared and of mismatching files.
func typeCheckDump_compare(t *testing.T, externs, root, oracleDir string, files []string, every int, diffDir string) (compared, mismatches int) {
	t.Helper()
	for i, rel := range files {
		if every > 1 && i%every != 0 {
			continue
		}
		want, err := typeCheckDump_readGz(filepath.Join(oracleDir, rel+".typecheck.gz"))
		if err != nil {
			t.Fatal(err)
		}
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		var sb strings.Builder
		typeCheckDump_dumpOne(&sb, externs, rel, string(src))
		got := sb.String()
		compared++
		if got == want {
			continue
		}
		mismatches++
		if mismatches <= 10 {
			t.Errorf("%s: Go output differs from Java:\n%s", rel, typeCheckDump_firstDiff(want, got))
		}
		if diffDir != "" {
			out := filepath.Join(diffDir, rel)
			if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(out+".go.txt", []byte(got), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(out+".java.txt", []byte(want), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return compared, mismatches
}

// typeCheckDump_externs reads tools/oracle/typecheck-externs.js
// (CompilerTypeTestCase.DEFAULT_EXTERNS).
func typeCheckDump_externs(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "tools", "oracle", "typecheck-externs.js"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestTypeCheckDumpGolden(t *testing.T) {
	dir := filepath.Join("testdata", "typecheck")
	files := typeCheckDump_readList(t, filepath.Join(dir, "files.txt"))
	compared, mismatches := typeCheckDump_compare(t, typeCheckDump_externs(t),
		filepath.Join(dir, "src"), filepath.Join(dir, "expected"), files, 1, os.Getenv("TYPECHECK_DIFF_DIR"))
	if mismatches > 0 {
		t.Errorf("%d of %d files differ", mismatches, compared)
	}
}

func TestTypeCheckDumpOracle(t *testing.T) {
	root := os.Getenv("TYPECHECK_ROOT")
	oracleDir := os.Getenv("TYPECHECK_ORACLE_DIR")
	list := os.Getenv("TYPECHECK_FILES")
	if root == "" || oracleDir == "" || list == "" {
		t.Skip("set TYPECHECK_ROOT, TYPECHECK_ORACLE_DIR and TYPECHECK_FILES to run the corpus oracle")
	}
	every := 1
	if s := os.Getenv("TYPECHECK_EVERY"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 {
			t.Fatalf("TYPECHECK_EVERY=%q: want a positive integer", s)
		}
		every = n
	}
	files := typeCheckDump_readList(t, list)
	compared, mismatches := typeCheckDump_compare(t, typeCheckDump_externs(t), root, oracleDir, files, every, os.Getenv("TYPECHECK_DIFF_DIR"))
	t.Logf("%d files compared, %d differ", compared, mismatches)
	if mismatches > 0 {
		t.Errorf("%d of %d files differ", mismatches, compared)
	}
}
