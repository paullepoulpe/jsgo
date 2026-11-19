package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIFlagsAndEndToEndDirectoryRun(t *testing.T) {
	tmpDir := t.TempDir()
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module example.com/clitest\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	j2Dir := filepath.Join(tmpDir, "j2gort")
	if err := os.MkdirAll(j2Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	j2Src := `package j2gort

func CheckNotNullIface[I interface{ IsJavaNil_() bool }](v I) I {
	if any(v) == nil || v.IsJavaNil_() {
		panic("NPE")
	}
	return v
}

func AsIface[I any](v any) I {
	return v.(I)
}
`
	if err := os.WriteFile(filepath.Join(j2Dir, "rt.go"), []byte(j2Src), 0o644); err != nil {
		t.Fatal(err)
	}

	pkgDir := filepath.Join(tmpDir, "genpkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	origSrc := `package genpkg

import "example.com/clitest/j2gort"

type Greeter interface {
	IsJavaNil_() bool
	Greet() string
}

type Impl struct{ msg string }

func NewImpl(m string) *Impl         { return &Impl{msg: m} }
func (i *Impl) IsJavaNil_() bool     { return i == nil }
func (i *Impl) Greet() string        { return i.msg }

func TransformMe(x int32) Greeter {
	_ = int32(x)
	var g Greeter = j2gort.CheckNotNullIface(j2gort.AsIface[Greeter](NewImpl("hi")))
	return g
}
`
	srcFile := filepath.Join(pkgDir, "gen.go")
	if err := os.WriteFile(srcFile, []byte(origSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1. -list
	var stdout, stderr bytes.Buffer
	if code := runMain([]string{"-list"}, &stdout, &stderr); code != 0 {
		t.Fatalf("-list failed (code=%d): %s", code, stderr.String())
	}
	for _, p := range []string{"native_slice_lowering", "dead_cast_cleanup", "j2gort_peephole"} {
		if !strings.Contains(stdout.String(), p) {
			t.Errorf("-list output missing %q: %s", p, stdout.String())
		}
	}
	for _, removed := range []string{"arraylist_slice", "stringbuilder"} {
		if strings.Contains(stdout.String(), removed) {
			t.Errorf("-list output still contains deleted pass %q: %s", removed, stdout.String())
		}
	}

	// 2. -passes=none -stats (must NOT modify file)
	stdout.Reset()
	stderr.Reset()
	if code := runMain([]string{"-passes=none", "-stats", pkgDir}, &stdout, &stderr); code != 0 {
		t.Fatalf("-passes=none failed (code=%d): %s", code, stderr.String())
	}
	afterNone, _ := os.ReadFile(srcFile)
	if string(afterNone) != origSrc {
		t.Fatalf("-passes=none modified file!")
	}

	// 3. -passes=all -dry-run -stats (computes 3 sites, does NOT modify file on disk)
	stdout.Reset()
	stderr.Reset()
	if code := runMain([]string{"-passes=all", "-dry-run", "-stats", pkgDir}, &stdout, &stderr); code != 0 {
		t.Fatalf("-dry-run failed (code=%d): %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "total:               3 sites") {
		t.Fatalf("expected 3 total sites in -dry-run stats, got:\n%s", stdout.String())
	}
	afterDry, _ := os.ReadFile(srcFile)
	if string(afterDry) != origSrc {
		t.Fatalf("-dry-run modified file on disk!")
	}

	// 4. -passes=all -stats (in-place write + unused j2gort import pruned!)
	stdout.Reset()
	stderr.Reset()
	if code := runMain([]string{"-passes=all", "-stats", pkgDir}, &stdout, &stderr); code != 0 {
		t.Fatalf("-passes=all failed (code=%d): %s", code, stderr.String())
	}
	afterWrite, _ := os.ReadFile(srcFile)
	got := string(afterWrite)
	if !strings.Contains(got, `var g Greeter = NewImpl("hi")`) || strings.Contains(got, "j2gort") || strings.Contains(got, "int32(x)") {
		t.Fatalf("unexpected transformed file:\n%s", got)
	}

	// 5. Re-running -passes=all -json-stats is strictly idempotent (0 sites).
	stdout.Reset()
	stderr.Reset()
	if code := runMain([]string{"-passes=all", "-json-stats", pkgDir}, &stdout, &stderr); code != 0 {
		t.Fatalf("second -passes=all failed (code=%d): %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"total_sites": 0`) {
		t.Fatalf("expected 0 sites on idempotent rerun, got:\n%s", stdout.String())
	}

	// 6. Invalid pass name returns exit code 2.
	stdout.Reset()
	stderr.Reset()
	if code := runMain([]string{"-passes=nonexistent", tmpDir}, &stdout, &stderr); code != 2 {
		t.Fatalf("expected exit code 2 for unknown pass, got %d", code)
	}
}
