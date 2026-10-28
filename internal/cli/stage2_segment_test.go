package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestCommandLineRunner_Stage2SegmentSplitParity(t *testing.T) {
	cfg := NewAbstractCommandLineRunner_CommandLineConfig()
	if got := cfg.GetStage2SegmentOfCompilationToRun(); got != jscomp.CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS {
		t.Fatalf("default Stage2SegmentOfCompilationToRun = %v, want OPTIMIZATIONS", got)
	}

	dir := t.TempDir()
	extPath := filepath.Join(dir, "externs.js")
	srcPath := filepath.Join(dir, "app.js")
	if err := os.WriteFile(extPath, []byte("/** @externs */\n/** @type {!Object} */ var window;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(srcPath, []byte(`
function computeValue(a, b) {
  const sum = a + b;
  return sum * 2;
}
window['ans'] = computeValue(20, 1);
`), 0o644); err != nil {
		t.Fatal(err)
	}

	baseArgs := []string{
		"--env", "CUSTOM",
		"--compilation_level", "ADVANCED_OPTIMIZATIONS",
		"--externs", extPath,
		"--js", srcPath,
	}

	runStep := func(extraArgs []string, configure func(*AbstractCommandLineRunner_CommandLineConfig)) string {
		var outBuf, errBuf bytes.Buffer
		args := append(append([]string{}, baseArgs...), extraArgs...)
		runner := NewCommandLineRunnerWithOutAndErr(
			args, jutil.NewPrintStream(&outBuf), jutil.NewPrintStream(&errBuf))
		if configure != nil {
			configure(runner.GetCommandLineConfig())
		}
		if !runner.ShouldRunCompiler() {
			t.Fatalf("ShouldRunCompiler() = false: %s", errBuf.String())
		}
		exitCode := -1
		runner.SetExitCodeReceiver(func(code int) { exitCode = code })
		runner.Run()
		if exitCode != 0 {
			t.Fatalf("DoRun exit=%d stderr=%s", exitCode, errBuf.String())
		}
		return outBuf.String()
	}

	wantOut := runStep(nil, nil)

	stage1File := filepath.Join(dir, "stage1.state")
	stage2Half1File := filepath.Join(dir, "stage2_first_half.state")
	stage2Half2File := filepath.Join(dir, "stage2_second_half.state")

	// Stage 1: checks -> save state 1.
	runStep(nil, func(c *AbstractCommandLineRunner_CommandLineConfig) {
		c.SetSaveCompilationStateToFilename(&stage1File, 1)
	})

	// Stage 2 (first half): restore state 1 -> OPTIMIZATIONS_FIRST_HALF -> save state 2.
	runStep(nil, func(c *AbstractCommandLineRunner_CommandLineConfig) {
		c.SetContinueSavedCompilationFileName(&stage1File, 1)
		c.SetSaveCompilationStateToFilename(&stage2Half1File, 2)
		c.SetStage2SegmentOfCompilationToRun(jscomp.CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF)
	})

	// Stage 2 (second half): restore first-half state -> OPTIMIZATIONS_SECOND_HALF -> save state 2.
	runStep(nil, func(c *AbstractCommandLineRunner_CommandLineConfig) {
		c.SetContinueSavedCompilationFileName(&stage2Half1File, 1)
		c.SetSaveCompilationStateToFilename(&stage2Half2File, 2)
		c.SetStage2SegmentOfCompilationToRun(jscomp.CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF)
	})

	// Stage 3: restore second-half state -> finalizations + codegen.
	gotOut := runStep(nil, func(c *AbstractCommandLineRunner_CommandLineConfig) {
		c.SetContinueSavedCompilationFileName(&stage2Half2File, 2)
	})

	if gotOut != wantOut {
		t.Fatalf("split Stage 2 output mismatch:\ngot:  %q\nwant: %q", gotOut, wantOut)
	}
}
