package linter

// Go-only: test-helper — the differential oracle for the refactoring fixes
// (SuggestedFix, ErrorToFixMapper, ApplySuggestedFixes, FixingErrorManager):
// the Go mirror of tools/oracle/java/RefactoringFixesDump.java, whose header
// documents the steps (Linter.fixRepeatedly in memory: each round lints the
// file through LintPassConfig with Linter.Builder()'s options, dumps
// FixingErrorManager.getAllFixes() and applies getSureFixes()) and the
// output format. Escaping, sections, list reading and the exception class
// names are ChecksDump's (checks_oracle_test.go).
//
//   - TestRefactoringFixesDumpGolden compares against the checked-in Java
//     output testdata/refactoring/expected-typecheck.txt for the corpus
//     sample of testdata/typecheck (files.txt, sources in src/) (always
//     runs). Regenerate with
//     java -Xss1g -cp closure-compiler.jar:<build> \
//       com.google.javascript.jscomp.RefactoringFixesDump \
//       internal/linter/testdata/expected-typecheck.txt \
//       internal/jscomp/testdata/typecheck/src @internal/jscomp/testdata/typecheck/files.txt
//   - TestRefactoringFixesDumpOracle compares a whole run when
//     REFACTORING_FIXES_ORACLE (the RefactoringFixesDump output, .gz or
//     plain), REFACTORING_FIXES_ROOT (the root of the sources) and
//     REFACTORING_FIXES_FILES (the list given to RefactoringFixesDump) are
//     set; REFACTORING_FIXES_DIFF_DIR, if set, receives <rel>.go.txt and
//     <rel>.java.txt for every mismatching file. Use absolute paths: `go
//     test` runs in the package directory.
//
// Known, accepted difference: Java catches StackOverflowError (the tool runs
// with -Xss1g); a Go stack overflow is fatal, so such an input stops the run.

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/refactoring"
)

// refactoringFixesDump_one is RefactoringFixesDump.dumpOne.
func refactoringFixesDump_one(sb *strings.Builder, name string, src string) {
	checksDumpLine(sb, "# file "+name)
	code := src
	for round := 0; round < Linter_MAX_FIXES; round++ {
		checksDumpLine(sb, "## round "+strconv.Itoa(round))
		if refactoringFixesDump_round(sb, name, &code) {
			break
		}
	}
	checksDumpLine(sb, "## final "+refactoringFixesDump_digest(code))
}

// refactoringFixesDump_round is one iteration of RefactoringFixesDump.dumpOne's
// loop (its try block and catch clause); it reports whether the file is
// done.
func refactoringFixesDump_round(sb *strings.Builder, name string, code *string) (done bool) {
	defer func() {
		if r := recover(); r != nil {
			checksDumpLine(sb, "## exception "+refactoringFixesDump_exception(r))
			done = true
		}
	}()
	options := Linter_Builder_().Options()
	compiler := jscomp.NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	errorManager := refactoring.NewFixingErrorManager()
	compiler.SetErrorManager(errorManager)
	errorManager.SetCompiler(compiler)

	// Linter.lint(Path, Compiler), with the code in memory.
	compiler.SetPassConfig(NewLintPassConfig(options))
	compiler.DisableThreads()
	externs := jscomp.SourceFile_FromCode("<Linter externs>", "")
	compiler.Compile([]*jscomp.SourceFile{externs}, []*jscomp.SourceFile{jscomp.SourceFile_FromCode(name, *code)}, options)

	for i, fix := range errorManager.GetAllFixes() {
		refactoringFixesDump_fix(sb, i, fix)
	}
	fixes := errorManager.GetSureFixes()
	checksDumpLine(sb, "## sure "+strconv.Itoa(len(fixes)))
	if len(fixes) == 0 {
		return true
	}
	// ApplySuggestedFixes.applySuggestedFixesToFiles writes only the files
	// that have replacements.
	codeMap := map[string]string{name: *code}
	if newCode, ok := refactoring.ApplySuggestedFixes_ApplySuggestedFixesToCode(slices.Values(fixes), codeMap)[name]; ok {
		*code = newCode
	}
	checksDumpLine(sb, "## code "+refactoringFixesDump_digest(*code))
	return false
}

// refactoringFixesDump_exception is the "## exception" text of the catch
// clause (RuntimeException | StackOverflowError): the class names of the
// exception and its causes. A nil dereference is Java's
// NullPointerException; anything else that is not a RuntimeException goes
// on (Java does not catch it either).
func refactoringFixesDump_exception(r any) string {
	if re, ok := r.(runtime.Error); ok && strings.Contains(re.Error(), "nil pointer dereference") {
		return "java.lang.NullPointerException"
	}
	t := jutil.AsThrowable(r)
	if _, ok := t.(jutil.IRuntimeException); !ok {
		panic(r)
	}
	ex := checksDumpClassName(t)
	for c := t.GetCause(); c != nil; c = c.GetCause() {
		ex += " <- " + checksDumpClassName(c)
	}
	return ex
}

// refactoringFixesDump_fix is RefactoringFixesDump.dumpFix.
func refactoringFixesDump_fix(sb *strings.Builder, i int, fix *refactoring.SuggestedFix) {
	info := fix.GetMatchedNodeInfo()
	at := "null"
	if info != nil {
		at = info.SourceFilename() + ":" + strconv.Itoa(info.Lineno()) + ":" + strconv.Itoa(info.Charno()) +
			" closurized=" + strconv.FormatBool(info.InClosurizedFile())
	}
	checksDumpLine(sb, "fix "+strconv.Itoa(i)+" alts="+strconv.Itoa(len(fix.GetAlternatives()))+
		" desc="+jutil.ToJavaString(fix.GetDescription())+" at="+at)
	for k, r := range fix.GetReplacements().Entries() {
		checksDumpLine(sb, "  "+k+" "+r.String())
	}
}

// refactoringFixesDump_digest is RefactoringFixesDump.digest: the length in
// UTF-16 units and the SHA-256 of the UTF-8 bytes (String.getBytes(UTF_8)).
func refactoringFixesDump_digest(code string) string {
	h := sha256.Sum256(jutil.StringGetBytesUTF8(code))
	return "len=" + strconv.Itoa(jutil.StringLength(code)) + " sha256=" + hex.EncodeToString(h[:])
}

// refactoringFixesDump_compare runs the Go side over files and compares it
// with the Java output in oraclePath; it returns the number of mismatching
// files.
func refactoringFixesDump_compare(t *testing.T, oraclePath string, root string, files []string, diffDir string) int {
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
		refactoringFixesDump_one(&sb, rel, string(src))
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

func TestRefactoringFixesDumpGolden(t *testing.T) {
	td := filepath.Join("..", "jscomp", "testdata", "typecheck")
	files := checksDumpReadList(t, filepath.Join(td, "files.txt"))
	n := refactoringFixesDump_compare(t,
		filepath.Join("testdata", "expected-typecheck.txt"),
		filepath.Join(td, "src"), files, "")
	if n != 0 {
		t.Errorf("%d of %d files differ", n, len(files))
	}
}

func TestRefactoringFixesDumpOracle(t *testing.T) {
	oracle, root, list := os.Getenv("REFACTORING_FIXES_ORACLE"), os.Getenv("REFACTORING_FIXES_ROOT"), os.Getenv("REFACTORING_FIXES_FILES")
	if oracle == "" || root == "" || list == "" {
		t.Skip("set REFACTORING_FIXES_ORACLE, REFACTORING_FIXES_ROOT and REFACTORING_FIXES_FILES")
	}
	files := checksDumpReadList(t, list)
	n := refactoringFixesDump_compare(t, oracle, root, files, os.Getenv("REFACTORING_FIXES_DIFF_DIR"))
	t.Logf("%d of %d files differ", n, len(files))
	if n != 0 {
		t.Fail()
	}
}
