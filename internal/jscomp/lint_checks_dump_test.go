package jscomp

// Differential oracle for the 23 jscomp/lint checks (CheckArrayWithGoogObject
// ... CheckVar).
//
// This file mirrors tools/oracle/java/LintChecksDump.java step by step: each
// input is compiled on its own by a fresh Compiler (one empty externs file,
// the parsing options of CompilerTestCase.getOptions(), module-load
// diagnostics off), Compiler.parseInputs() builds the AST, then each check
// runs over it as a fresh instance, in a fixed order, and every JSError it
// reports is captured raw (the Java tool overrides Compiler.report; here a
// wrapper around the Compiler, which the checks get as their AbstractCompiler,
// overrides Report). See the Java file for the output format.
//
// TestLintChecksDumpGolden always runs: it compares the Go output with the
// checked-in Java output for testdata/lint_checks/src (hand-written inputs, one
// or more per check) and for the repository's testdata/fixtures. Regenerate
// the golden files with (CP=.cache/closure-compiler.jar:<LintChecksDump build>):
//
//	cd internal/jscomp/testdata/lint_checks
//	java -Xss1g -cp $CP com.google.javascript.jscomp.LintChecksDump \
//	  --convention=default src-default.txt src @inputs.txt
//	(cd ../../../../testdata/fixtures && find . -name '*.js' | sed 's|^./||' | LC_ALL=C sort) > /tmp/fixtures.txt
//	java -Xss1g -cp $CP com.google.javascript.jscomp.LintChecksDump \
//	  --convention=default --gzip fixtures-default.txt.gz ../../../../testdata/fixtures @/tmp/fixtures.txt
//
// and the same with --convention=google into src-google.txt and
// fixtures-google.txt.gz.
//
// TestLintChecksDumpCorpus runs on any Java output: set LINT_ORACLE=<the
// LintChecksDump output, .gz or not>, LINT_ROOT=<the root it ran on>,
// LINT_FILES=<the list it ran on>, LINT_CONVENTION=<its --convention>
// (default: google), and optionally LINT_DIFF_DIR=<dir> to write the Go
// output there (lint.txt).

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

// lintChecksDumpCheck is one lint check: its name and how to run it over
// (externsRoot, mainRoot).
type lintChecksDumpCheck struct {
	name string
	run  func(c AbstractCompiler, externsRoot, mainRoot *rhino.Node)
}

func lintChecksDumpPass(name string, create func(c AbstractCompiler) CompilerPass) lintChecksDumpCheck {
	return lintChecksDumpCheck{name, func(c AbstractCompiler, externsRoot, mainRoot *rhino.Node) {
		create(c).Process(externsRoot, mainRoot)
	}}
}

// lintChecksDumpChecks is LintChecksDump.CHECKS, in the same order.
var lintChecksDumpChecks = []lintChecksDumpCheck{
	lintChecksDumpPass("CheckArrayWithGoogObject", func(c AbstractCompiler) CompilerPass { return NewCheckArrayWithGoogObject(c) }),
	lintChecksDumpPass("CheckConstPrivateProperties", func(c AbstractCompiler) CompilerPass { return NewCheckConstPrivateProperties(c) }),
	lintChecksDumpPass("CheckConstantCaseNames", func(c AbstractCompiler) CompilerPass { return NewCheckConstantCaseNames(c) }),
	lintChecksDumpPass("CheckDuplicateCase", func(c AbstractCompiler) CompilerPass { return NewCheckDuplicateCase(c) }),
	lintChecksDumpPass("CheckEmptyStatements", func(c AbstractCompiler) CompilerPass { return NewCheckEmptyStatements(c) }),
	lintChecksDumpPass("CheckEnums", func(c AbstractCompiler) CompilerPass { return NewCheckEnums(c) }),
	lintChecksDumpPass("CheckEs6ModuleFileStructure", func(c AbstractCompiler) CompilerPass { return NewCheckEs6ModuleFileStructure(c) }),
	lintChecksDumpPass("CheckEs6Modules", func(c AbstractCompiler) CompilerPass { return NewCheckEs6Modules(c) }),
	lintChecksDumpPass("CheckExtraRequires", func(c AbstractCompiler) CompilerPass { return NewCheckExtraRequires(c, nil) }),
	lintChecksDumpPass("CheckGoogModuleTypeScriptName", func(c AbstractCompiler) CompilerPass { return NewCheckGoogModuleTypeScriptName(c) }),
	lintChecksDumpPass("CheckInterfaces", func(c AbstractCompiler) CompilerPass { return NewCheckInterfaces(c) }),
	lintChecksDumpPass("CheckMissingSemicolon", func(c AbstractCompiler) CompilerPass { return NewCheckMissingSemicolon(c) }),
	lintChecksDumpPass("CheckNestedNames", func(c AbstractCompiler) CompilerPass { return NewCheckNestedNames(c) }),
	lintChecksDumpPass("CheckNoMutatedEs6Exports", func(c AbstractCompiler) CompilerPass { return NewCheckNoMutatedEs6Exports(c) }),
	lintChecksDumpPass("CheckNullabilityModifiers", func(c AbstractCompiler) CompilerPass { return NewCheckNullabilityModifiers(c) }),
	lintChecksDumpPass("CheckPrimitiveAsObject", func(c AbstractCompiler) CompilerPass { return NewCheckPrimitiveAsObject(c) }),
	lintChecksDumpPass("CheckPrototypeProperties", func(c AbstractCompiler) CompilerPass { return NewCheckPrototypeProperties(c) }),
	{"CheckProvidesSorted", func(c AbstractCompiler, externsRoot, mainRoot *rhino.Node) {
		NodeTraversal_Traverse(c, mainRoot, NewCheckProvidesSorted(CheckProvidesSorted_Mode_COLLECT_AND_REPORT))
	}},
	{"CheckRequiresSorted", func(c AbstractCompiler, externsRoot, mainRoot *rhino.Node) {
		NodeTraversal_Traverse(c, mainRoot, NewCheckRequiresSorted(CheckRequiresSorted_Mode_COLLECT_AND_REPORT))
	}},
	lintChecksDumpPass("CheckUnusedLabels", func(c AbstractCompiler) CompilerPass { return NewCheckUnusedLabels(c) }),
	lintChecksDumpPass("CheckUnusedPrivateProperties", func(c AbstractCompiler) CompilerPass { return NewCheckUnusedPrivateProperties(c) }),
	lintChecksDumpPass("CheckUselessBlocks", func(c AbstractCompiler) CompilerPass { return NewCheckUselessBlocks(c) }),
	lintChecksDumpPass("CheckVar", func(c AbstractCompiler) CompilerPass { return NewCheckVar(c) }),
}

// lintChecksDumpCompiler is the Java tool's anonymous Compiler subclass: it
// records every JSError reported while a check runs.
type lintChecksDumpCompiler struct {
	*Compiler
	phase string
	diags *[]string
}

func (c *lintChecksDumpCompiler) Report(error *JSError) {
	if c.phase != "" {
		*c.diags = append(*c.diags, c.phase+" "+lintChecksDumpFormat(error))
	}
	c.Compiler.Report(error)
}

// lintChecksDumpOptions is LintChecksDump.options().
func lintChecksDumpOptions(t testing.TB, convention string) *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	options.SetEmitUseStrict(false)
	options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	options.SetModuleResolutionMode(ModuleLoader_ResolutionMode_BROWSER)
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_TYPES_ONLY)
	options.SetPreserveTypeAnnotations(true)
	options.SetWarningLevel(DiagnosticGroup_ForType(ModuleLoader_LOAD_WARNING), CheckLevel_OFF)
	options.SetWarningLevel(DiagnosticGroup_ForType(ModuleLoader_INVALID_MODULE_PATH), CheckLevel_OFF)
	switch convention {
	case "google":
		options.SetCodingConvention(NewGoogleCodingConvention())
	case "closure":
	case "default":
		options.SetCodingConvention(CodingConventions_GetDefault())
	default:
		t.Fatalf("unknown convention %q", convention)
	}
	return options
}

// lintChecksDumpOne is LintChecksDump.dumpOne.
func lintChecksDumpOne(t testing.TB, sb *strings.Builder, convention string, name string, src string) {
	lintChecksDumpLine(sb, "# file "+name)
	var diags []string
	compiler := &lintChecksDumpCompiler{Compiler: NewCompiler(), diags: &diags}
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs", "")},
		[]*SourceFile{SourceFile_FromCode(name, src)},
		lintChecksDumpOptions(t, convention))
	var root *rhino.Node
	jutil.Try(func() {
		root = compiler.ParseInputs()
	}, func(e jutil.IRuntimeException) {
		root = nil
	})
	if root == nil || compiler.HasErrors() {
		failed := "# parse-failed"
		for _, e := range compiler.GetErrors() {
			failed += " " + e.GetType().key
		}
		lintChecksDumpLine(sb, failed)
		return
	}
	externsRoot, mainRoot := root.GetFirstChild(), root.GetSecondChild()
	for _, check := range lintChecksDumpChecks {
		compiler.phase = check.name
		start := len(diags)
		jutil.Try(func() {
			check.run(compiler, externsRoot, mainRoot)
		}, func(e jutil.IRuntimeException) {
			ex := check.name + " # exception " + lintChecksDumpClassName(e)
			for c := e.GetCause(); c != nil; c = c.GetCause() {
				ex += " <- " + lintChecksDumpClassName(c)
			}
			diags = append(diags, ex)
		})
		if check.name == "CheckNoMutatedEs6Exports" {
			slices.SortStableFunc(diags[start:], jutil.StringCompareTo)
		}
	}
	compiler.phase = ""
	for _, d := range diags {
		lintChecksDumpLine(sb, d)
	}
}

// lintChecksDumpClassName is e.getClass().getName(): the class part of
// Throwable.toString().
func lintChecksDumpClassName(e jutil.Throwable) string {
	s := e.Error()
	if i := strings.Index(s, ": "); i >= 0 {
		return s[:i]
	}
	return s
}

// lintChecksDumpFormat is LintChecksDump.format(JSError).
func lintChecksDumpFormat(e *JSError) string {
	return e.GetType().key +
		" " + e.GetDefaultLevel().String() +
		" " + strconv.Itoa(e.GetLineno()) +
		":" + strconv.Itoa(e.GetCharno()) +
		" +" + strconv.Itoa(e.GetLength()) +
		" @" + strconv.Itoa(e.GetNodeSourceOffset()) +
		" \"" + e.GetDescription() + "\""
}

// lintChecksDumpLine is LintChecksDump.line: backslash and every UTF-16 unit
// outside printable ASCII become escapes.
func lintChecksDumpLine(sb *strings.Builder, s string) {
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			sb.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			sb.WriteByte(byte(c))
		default:
			const hex = "0123456789abcdef"
			sb.WriteString("\\u")
			sb.WriteByte(hex[c>>12&0xf])
			sb.WriteByte(hex[c>>8&0xf])
			sb.WriteByte(hex[c>>4&0xf])
			sb.WriteByte(hex[c&0xf])
		}
	}
	sb.WriteByte('\n')
}

// lintChecksDumpBlocks splits a dump into per-file blocks ("# file" line
// included), in order.
func lintChecksDumpBlocks(t testing.TB, r io.Reader) (names []string, blocks map[string]string) {
	blocks = map[string]string{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1<<20), 1<<30)
	var cur strings.Builder
	name := ""
	flush := func() {
		if name != "" {
			names = append(names, name)
			blocks[name] = cur.String()
		}
		cur.Reset()
	}
	for sc.Scan() {
		line := sc.Text()
		if rest, ok := strings.CutPrefix(line, "# file "); ok {
			flush()
			name = rest
		}
		cur.WriteString(line)
		cur.WriteByte('\n')
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	flush()
	return names, blocks
}

func lintChecksDumpOpen(t testing.TB, path string) io.Reader {
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })
	if strings.HasSuffix(path, ".gz") {
		z, err := gzip.NewReader(f)
		if err != nil {
			t.Fatal(err)
		}
		return z
	}
	return f
}

// lintChecksDumpCompare runs the Go dump over files (relative to root) and
// compares it, file by file, with the Java dump at oracle.
func lintChecksDumpCompare(t *testing.T, oracle string, root string, files []string, convention string, diffDir string) {
	javaNames, java := lintChecksDumpBlocks(t, lintChecksDumpOpen(t, oracle))
	if !slices.Equal(javaNames, files) {
		t.Fatalf("%s covers %d files, the list has %d (or a different order)", oracle, len(javaNames), len(files))
	}
	var all strings.Builder
	diffs := 0
	for _, name := range files {
		src, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		var sb strings.Builder
		lintChecksDumpOne(t, &sb, convention, name, string(src))
		got := sb.String()
		if diffDir != "" {
			all.WriteString(got)
		}
		if want := java[name]; got != want {
			diffs++
			if diffs <= 10 {
				t.Errorf("%s: Go output differs from Java\n--- Java\n%s--- Go\n%s", name, want, got)
			}
		}
	}
	if diffDir != "" {
		if err := os.WriteFile(filepath.Join(diffDir, "lint.txt"), []byte(all.String()), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if diffs != 0 {
		t.Errorf("%d of %d files differ", diffs, len(files))
	}
}

// lintChecksDumpReadList reads a file list, without the files of
// tools/oracle/gpl-excluded-corpus.txt.
func lintChecksDumpReadList(t testing.TB, path string) []string {
	t.Helper()
	files, err := corpus.ReadList(path)
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func TestLintChecksDumpGolden(t *testing.T) {
	defer SkipIfPendingPort(t) // the Compiler reaches passes not ported yet
	srcFiles := lintChecksDumpReadList(t, "testdata/lint_checks/inputs.txt")
	const fixturesRoot = "../../testdata/fixtures"
	var fixtures []string
	err := filepath.WalkDir(fixturesRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".js") {
			rel, err := filepath.Rel(fixturesRoot, path)
			if err != nil {
				return err
			}
			fixtures = append(fixtures, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(fixtures) // LC_ALL=C sort: byte order
	for _, convention := range []string{"default", "google"} {
		t.Run(convention, func(t *testing.T) {
			defer SkipIfPendingPort(t) // the Compiler reaches passes not ported yet
			lintChecksDumpCompare(t, "testdata/lint_checks/src-"+convention+".txt", "testdata/lint_checks/src", srcFiles, convention, "")
			lintChecksDumpCompare(t, "testdata/lint_checks/fixtures-"+convention+".txt.gz", fixturesRoot, fixtures, convention, "")
		})
	}
}

func TestLintChecksDumpCorpus(t *testing.T) {
	defer SkipIfPendingPort(t) // the Compiler reaches passes not ported yet
	oracle := os.Getenv("LINT_ORACLE")
	if oracle == "" {
		t.Skip("set LINT_ORACLE, LINT_ROOT and LINT_FILES to compare with a LintChecksDump output")
	}
	convention := os.Getenv("LINT_CONVENTION")
	if convention == "" {
		convention = "google"
	}
	lintChecksDumpCompare(t, oracle, os.Getenv("LINT_ROOT"), lintChecksDumpReadList(t, os.Getenv("LINT_FILES")), convention, os.Getenv("LINT_DIFF_DIR"))
}
