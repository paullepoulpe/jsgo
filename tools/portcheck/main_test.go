package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// -------- fixtures --------

// highFidelityGoSrc is a Go file that closely mirrors the fake Java source below.
const highFidelityGoSrc = `// Port of com/google/javascript/jscomp/FakeChecker.java @ closure-compiler v20260415
package jscomp

// FakeChecker checks things.
type FakeChecker struct{}

func (c *FakeChecker) check(x int) bool {
	if x > 0 {
		return true
	}
	return false
}

func (c *FakeChecker) validate(x int, y int) bool {
	return x == y
}

func (c *FakeChecker) report(msg string) {
	_ = msg
}

func (c *FakeChecker) run() {
	c.check(1)
	c.validate(1, 2)
	c.report("done")
}
`

// highFidelityJavaSrc is the Java source for the above.
const highFidelityJavaSrc = `
package com.google.javascript.jscomp;

public class FakeChecker {
  private boolean check(int x) {
    if (x > 0) { return true; }
    return false;
  }
  private boolean validate(int x, int y) {
    return x == y;
  }
  private void report(String msg) {
    // do nothing
  }
  public void run() {
    check(1);
    validate(1, 2);
    report("done");
  }
}
`

// stubbyGoSrc has many stub functions.
const stubbyGoSrc = `// Port of com/google/javascript/jscomp/FakeChecker.java @ closure-compiler v20260415
package jscomp

func check(x int) bool {
	return false
}

func validate(x int, y int) bool {
	return false
}

func report(msg string) {
}

func run() {
	panic("not implemented")
}
`

// noHeaderGoSrc has no Port of header.
const noHeaderGoSrc = `package jscomp

func doSomething() bool {
	return true
}
`

// goOnlyAdditionsSrc has Go functions not present in Java.
const goOnlyAdditionsSrc = `// Port of com/google/javascript/jscomp/FakeChecker.java @ closure-compiler v20260415
package jscomp

func check(x int) bool {
	return x > 0
}

// goSpecificHelper is a Go-only addition.
func goSpecificHelperXYZQQQ(x int) int {
	return x * 2
}

// anotherGoOnlyFunctionABCDEF
func anotherGoOnlyFunctionABCDEF() string {
	return "go-only"
}
`

// -------- helpers --------

// writeTempFile writes content to a temp file and returns the path.
func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return path
}

// -------- tests --------

func TestParseHeader_Valid(t *testing.T) {
	src := "// Port of com/google/javascript/jscomp/TypeCheck.java @ closure-compiler v20260415\n"
	hdr := parseHeader(src)
	if !hdr.valid {
		t.Error("expected valid header")
	}
	if hdr.javaPath != "com/google/javascript/jscomp/TypeCheck.java" {
		t.Errorf("unexpected javaPath: %q", hdr.javaPath)
	}
	if hdr.version != "v20260415" {
		t.Errorf("unexpected version: %q", hdr.version)
	}
}

func TestParseHeader_Missing(t *testing.T) {
	src := "package jscomp\n\nfunc Foo() {}\n"
	hdr := parseHeader(src)
	if hdr.valid {
		t.Error("expected invalid header")
	}
}

func TestParseHeader_VersionMismatch(t *testing.T) {
	src := "// Port of com/google/javascript/jscomp/Foo.java @ closure-compiler v20230101\n"
	hdr := parseHeader(src)
	if !hdr.valid {
		t.Error("expected valid header")
	}
	if hdr.version == expectedVersion {
		t.Errorf("expected version mismatch, got: %q", hdr.version)
	}
}

func TestExtractJavaMethods(t *testing.T) {
	methods := extractJavaMethods(highFidelityJavaSrc)
	want := []string{"check", "validate", "report", "run"}
	for _, w := range want {
		if _, ok := methods[w]; !ok {
			t.Errorf("expected Java method %q, got keys: %v", w, mapKeys(methods))
		}
	}
}

func TestExtractGoFunctions_Count(t *testing.T) {
	path := writeTempFile(t, "fake.go", highFidelityGoSrc)
	funcs, _, _, err := extractGoFunctions(path)
	if err != nil {
		t.Fatalf("extractGoFunctions: %v", err)
	}
	if len(funcs) != 4 {
		t.Errorf("expected 4 functions, got %d: %v", len(funcs), funcNames(funcs))
	}
}

func TestStubDetection(t *testing.T) {
	path := writeTempFile(t, "stub.go", stubbyGoSrc)
	funcs, _, _, err := extractGoFunctions(path)
	if err != nil {
		t.Fatalf("extractGoFunctions: %v", err)
	}
	stubCount := 0
	for _, f := range funcs {
		if f.IsStub {
			stubCount++
		}
	}
	// check, validate, report, run are all stubs
	if stubCount < 3 {
		t.Errorf("expected at least 3 stubs, got %d (funcs: %v)", stubCount, funcNames(funcs))
	}
}

// makeJavaFixture creates the Java fixture file at the expected path under tmpDir:
// tmpDir/com/google/javascript/jscomp/FakeChecker.java
func makeJavaFixture(t *testing.T, tmpDir string) {
	t.Helper()
	javaDir := filepath.Join(tmpDir, "com", "google", "javascript", "jscomp")
	if err := os.MkdirAll(javaDir, 0755); err != nil {
		t.Fatal(err)
	}
	javaFile := filepath.Join(javaDir, "FakeChecker.java")
	if err := os.WriteFile(javaFile, []byte(highFidelityJavaSrc), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestHighFidelityScore(t *testing.T) {
	tmpDir := t.TempDir()
	makeJavaFixture(t, tmpDir)

	goFile := writeTempFile(t, "fakechecker.go", highFidelityGoSrc)

	result := computeScore(goFile, highFidelityGoSrc, tmpDir)
	if result.Score < 50 {
		t.Errorf("high-fidelity file scored too low: %.1f (notes: %s)", result.Score, result.Notes)
	}
	t.Logf("high-fidelity score: %.1f", result.Score)
}

func TestStubbyScore(t *testing.T) {
	tmpDir := t.TempDir()
	makeJavaFixture(t, tmpDir)

	// High-fidelity score for comparison
	hfGoFile := writeTempFile(t, "hffake.go", highFidelityGoSrc)
	hfResult := computeScore(hfGoFile, highFidelityGoSrc, tmpDir)

	goFile := writeTempFile(t, "stub.go", stubbyGoSrc)
	result := computeScore(goFile, stubbyGoSrc, tmpDir)

	// Stub count should be detected
	if result.StubCount < 3 {
		t.Errorf("expected at least 3 stubs detected, got %d", result.StubCount)
	}
	// Stubby should score strictly lower than high-fidelity equivalent
	if result.Score >= hfResult.Score {
		t.Errorf("stubby file (%.1f) should score below high-fidelity (%.1f)", result.Score, hfResult.Score)
	}
	t.Logf("stubby score: %.1f (stubs: %d, hf score: %.1f)", result.Score, result.StubCount, hfResult.Score)
}

func TestMissingHeaderScore(t *testing.T) {
	goFile := writeTempFile(t, "noheader.go", noHeaderGoSrc)
	result := computeScore(goFile, noHeaderGoSrc, "/nonexistent")
	if result.Score != 0 {
		t.Errorf("expected score=0 for missing header, got %.1f", result.Score)
	}
	if !strings.Contains(result.Notes, "missing") {
		t.Errorf("expected 'missing' in notes, got: %q", result.Notes)
	}
}

func TestGoOnlyFunctionsScore(t *testing.T) {
	tmpDir := t.TempDir()
	makeJavaFixture(t, tmpDir)

	goFile := writeTempFile(t, "goonly.go", goOnlyAdditionsSrc)
	result := computeScore(goFile, goOnlyAdditionsSrc, tmpDir)

	// Should have Go-only functions
	if result.GoOnlyCount < 2 {
		t.Errorf("expected at least 2 go-only functions, got %d", result.GoOnlyCount)
	}
	t.Logf("go-only score: %.1f (go-only count: %d)", result.Score, result.GoOnlyCount)
}

func TestJavaBodyLoC(t *testing.T) {
	loc := javaBodyLoC(highFidelityJavaSrc)
	if loc < 5 {
		t.Errorf("expected >5 body lines, got %d", loc)
	}
}

func TestComputeStats(t *testing.T) {
	results := []FileResult{
		{Score: 10},
		{Score: 50},
		{Score: 90},
	}
	stats := computeStats(results)
	if stats.TotalFiles != 3 {
		t.Errorf("expected 3 files, got %d", stats.TotalFiles)
	}
	if stats.AvgScore != 50 {
		t.Errorf("expected avg 50, got %.1f", stats.AvgScore)
	}
	if stats.MinScore != 10 {
		t.Errorf("expected min 10, got %.1f", stats.MinScore)
	}
	if stats.MaxScore != 90 {
		t.Errorf("expected max 90, got %.1f", stats.MaxScore)
	}
}

// -------- tiny wrapper + @Override tests --------

// nullTypeGoSrc mimics null_type.go: a ≤50-LoC Go file that implements a tiny
// type wrapper by satisfying a Go interface.  Its methods do not have Java-identical
// names, but the parent Java class (NullType → ValueType) contains recognisable
// method names that should trigger the tiny-wrapper bonus.
const nullTypeGoSrc = `// Port of com/google/javascript/rhino/jstype/NullType.java @ closure-compiler v20260415
package jstype

// nullType is the null type.
type nullType struct{}

func (nullType) TypeName() string { return "null" }
func (nullType) IsNullType() bool { return true }
func (nullType) IsNullable() bool { return true }
func (nullType) MatchesStringContext() bool { return false }
func (nullType) TestForEquality(that JSType) int { return 0 }
`

// nullTypeJavaSrc is the Java source for NullType — a small class with @Override methods.
const nullTypeJavaSrc = `
package com.google.javascript.rhino.jstype;

public final class NullType extends ValueType {
  @Override
  public boolean isNullType() {
    return true;
  }
  @Override
  public boolean isNullable() {
    return true;
  }
  @Override
  public boolean matchesStringContext() {
    return false;
  }
  @Override
  public JSType testForEquality(JSType that) {
    return null;
  }
  @Override
  public <T> T visit(Visitor<T> visitor) {
    return visitor.caseNullType(this);
  }
}
`

// valueTypeJavaSrc is the parent class ValueType — contains additional @Override methods.
const valueTypeJavaSrc = `
package com.google.javascript.rhino.jstype;

public abstract class ValueType extends JSType {
  @Override
  public String getDisplayName() {
    return null;
  }
  @Override
  public boolean isNullable() {
    return false;
  }
}
`

// makeNullTypeFixtures writes the Java NullType + ValueType fixtures under tmpDir.
func makeNullTypeFixtures(t *testing.T, tmpDir string) {
	t.Helper()
	javaDir := filepath.Join(tmpDir, "com", "google", "javascript", "rhino", "jstype")
	if err := os.MkdirAll(javaDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(javaDir, "NullType.java"), []byte(nullTypeJavaSrc), 0644); err != nil {
		t.Fatal(err)
	}
	// Parent class — tiny-wrapper bonus looks here.
	if err := os.WriteFile(filepath.Join(javaDir, "ValueType.java"), []byte(valueTypeJavaSrc), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestTinyWrapperBonus(t *testing.T) {
	tmpDir := t.TempDir()
	makeNullTypeFixtures(t, tmpDir)

	// Redirect the header path so computeScore finds the fixture.
	goSrc := strings.Replace(nullTypeGoSrc,
		"com/google/javascript/rhino/jstype/NullType.java",
		"com/google/javascript/rhino/jstype/NullType.java",
		1)
	goFile := writeTempFile(t, "null_type.go", goSrc)

	result := computeScore(goFile, goSrc, tmpDir)
	t.Logf("tiny-wrapper score: %.1f (notes: %s, overlap: %.0f%%)", result.Score, result.Notes, result.MethodOverlapPct)

	if result.Score < 70 {
		t.Errorf("tiny-wrapper file scored too low: %.1f — expected ≥70 (notes: %s)", result.Score, result.Notes)
	}
}

// TestJavaMethodRegexAnnotations verifies that methods preceded by multi-line
// annotation stacks (@Override\n@Nullable\n) are correctly extracted.
func TestJavaMethodRegexAnnotations(t *testing.T) {
	src := `
package com.google.javascript.rhino.jstype;

public class Example {
  @Override
  @Nullable
  public boolean isNullType() {
    return true;
  }

  @Override
  public <T> T visit(Visitor<T> v) {
    return v.caseNull(this);
  }

  @SuppressWarnings("unchecked")
  public static Example of() {
    return null;
  }
}
`
	methods := extractJavaMethods(src)
	want := []string{"isnulltype", "visit", "of"}
	for _, w := range want {
		if _, ok := methods[w]; !ok {
			t.Errorf("expected Java method %q to be extracted, got: %v", w, mapKeys(methods))
		}
	}
}

// TestGoOnlyRecognition checks that functions with a valid "// Go-only: <code>" comment
// are flagged as intentional divergences and do NOT inflate the go-only penalty.
const goOnlyValidTaggedSrc = `// Port of com/google/javascript/jscomp/FakeChecker.java @ closure-compiler v20260415
package jscomp

func check(x int) bool {
	return x > 0
}

// Go-only: language-adapter — Go interface adapter; no direct Java counterpart.
func goHelperForInterfaceAdaption() int {
	return 42
}
`

func TestGoOnlyValidTaggedFunctionNotPenalised(t *testing.T) {
	tmpDir := t.TempDir()
	makeJavaFixture(t, tmpDir) // uses highFidelityJavaSrc which has check()

	goFile := writeTempFile(t, "tagged.go", goOnlyValidTaggedSrc)
	result := computeScore(goFile, goOnlyValidTaggedSrc, tmpDir)

	// GoOnlyCount should be 0 for the validly tagged function.
	if result.GoOnlyCount > 0 {
		t.Errorf("valid Go-only tagged function should not inflate GoOnlyCount: got %d", result.GoOnlyCount)
	}
	t.Logf("go-only valid tagged score: %.1f, GoOnlyCount: %d, notes: %s", result.Score, result.GoOnlyCount, result.Notes)
}

// TestGoOnlyInvalidTagPenalized checks that a "// Go-only:" with an invalid/unknown
// reason code IS penalized (treated as bare go-only).
const goOnlyInvalidTaggedSrc = `// Port of com/google/javascript/jscomp/FakeChecker.java @ closure-compiler v20260415
package jscomp

func check(x int) bool {
	return x > 0
}

// Go-only: Go interface adapter; no direct Java counterpart.
func goHelperForInterfaceAdaptionXYZ() int {
	return 42
}
`

func TestGoOnlyInvalidTagPenalized(t *testing.T) {
	tmpDir := t.TempDir()
	makeJavaFixture(t, tmpDir)

	goFile := writeTempFile(t, "invalid_tagged.go", goOnlyInvalidTaggedSrc)
	result := computeScore(goFile, goOnlyInvalidTaggedSrc, tmpDir)

	// GoOnlyCount should be >0 because the reason code is invalid.
	if result.GoOnlyCount == 0 {
		t.Errorf("invalid Go-only tag should still penalize: GoOnlyCount=0")
	}
	t.Logf("go-only invalid tagged score: %.1f, GoOnlyCount: %d", result.Score, result.GoOnlyCount)
}

// TestTodoPortNotPenalized checks that "// TODO(port):" functions are counted
// in port_debt but NOT penalized as go-only.
const todoPortSrc = `// Port of com/google/javascript/jscomp/FakeChecker.java @ closure-compiler v20260415
package jscomp

func check(x int) bool {
	return x > 0
}

// TODO(port): FakeChecker.validate(int, int)
func validate(x int, y int) bool {
	return false
}

// TODO(port): FakeChecker.reportXYZQQQABC(String)
func reportXYZQQQABC(msg string) {
}
`

func TestTodoPortNotPenalized(t *testing.T) {
	tmpDir := t.TempDir()
	makeJavaFixture(t, tmpDir)

	goFile := writeTempFile(t, "tododebt.go", todoPortSrc)
	result := computeScore(goFile, todoPortSrc, tmpDir)

	// PortDebtCount should be 2 (both TODO(port) functions).
	if result.PortDebtCount != 2 {
		t.Errorf("expected PortDebtCount=2, got %d", result.PortDebtCount)
	}
	// GoOnlyCount should NOT include the TODO(port) functions.
	// validate IS in Java (so not go-only anyway), reportXYZQQQABC is not but is marked as debt.
	if result.GoOnlyCount > 0 {
		t.Errorf("TODO(port) functions should not be penalized as go-only, GoOnlyCount=%d", result.GoOnlyCount)
	}
	t.Logf("todo-port score: %.1f, PortDebtCount: %d, GoOnlyCount: %d", result.Score, result.PortDebtCount, result.GoOnlyCount)
}

// TestParseGoOnlyCode checks the reason code parser.
func TestParseGoOnlyCode(t *testing.T) {
	tests := []struct {
		comment string
		wantCode string
		wantValid bool
	}{
		{"// Go-only: language-adapter — some explanation", "language-adapter", true},
		{"// Go-only: package-bridge", "package-bridge", true},
		{"// Go-only: stdlib-interop", "stdlib-interop", true},
		{"// Go-only: test-helper", "test-helper", true},
		{"// Go-only: Go interface adapter; no direct Java counterpart.", "Go", false},
		{"// Go-only: constructor shim; Java uses the class constructor directly.", "constructor", false},
		{"// Go-only: Java-named alias for private helper", "Java-named", false},
	}
	for _, tt := range tests {
		code, isValid := parseGoOnlyCode(tt.comment)
		if code != tt.wantCode {
			t.Errorf("parseGoOnlyCode(%q): code=%q want=%q", tt.comment, code, tt.wantCode)
		}
		if isValid != tt.wantValid {
			t.Errorf("parseGoOnlyCode(%q): isValid=%v want=%v", tt.comment, isValid, tt.wantValid)
		}
	}
}

// TestAuditEntries checks that extractGoFunctions produces audit entries.
func TestAuditEntries(t *testing.T) {
	src := `// Port of com/google/javascript/jscomp/FakeChecker.java @ closure-compiler v20260415
package jscomp

func check(x int) bool { return x > 0 }

// Go-only: language-adapter — valid reason.
func adapterHelper() int { return 42 }

// Go-only: invalid-code — not a real reason code.
func invalidHelper() int { return 99 }

// TODO(port): FakeChecker.missingMethod(String)
func missingMethod(s string) {}
`
	goFile := writeTempFile(t, "audit_test.go", src)
	_, _, audits, err := extractGoFunctions(goFile)
	if err != nil {
		t.Fatalf("extractGoFunctions: %v", err)
	}

	validCount := 0
	invalidCount := 0
	debtCount := 0
	for _, a := range audits {
		switch a.Kind {
		case "go-only-valid":
			validCount++
		case "go-only-invalid":
			invalidCount++
		case "todo-port":
			debtCount++
		}
	}

	if validCount != 1 {
		t.Errorf("expected 1 valid go-only, got %d", validCount)
	}
	if invalidCount != 1 {
		t.Errorf("expected 1 invalid go-only, got %d", invalidCount)
	}
	if debtCount != 1 {
		t.Errorf("expected 1 todo-port, got %d", debtCount)
	}
}

// -------- helpers --------

func mapKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func funcNames(funcs []GoFunc) []string {
	names := make([]string, len(funcs))
	for i, f := range funcs {
		names[i] = f.Name
	}
	return names
}
