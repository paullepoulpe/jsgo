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
	funcs, _, err := extractGoFunctions(path)
	if err != nil {
		t.Fatalf("extractGoFunctions: %v", err)
	}
	if len(funcs) != 4 {
		t.Errorf("expected 4 functions, got %d: %v", len(funcs), funcNames(funcs))
	}
}

func TestStubDetection(t *testing.T) {
	path := writeTempFile(t, "stub.go", stubbyGoSrc)
	funcs, _, err := extractGoFunctions(path)
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
