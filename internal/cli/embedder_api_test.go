package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/cli"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type testMetricsRecorder struct {
	events     []string
	actionName string
	hadResult  bool
}

func (m *testMetricsRecorder) RecordActionStart() {
	m.events = append(m.events, "start")
}

func (m *testMetricsRecorder) RecordActionName(actionName string) {
	m.events = append(m.events, "name:"+actionName)
	m.actionName = actionName
}

func (m *testMetricsRecorder) RecordStartState(compiler jscomp.AbstractCompiler) {
	if compiler != nil {
		m.events = append(m.events, "state")
	}
}

func (m *testMetricsRecorder) RecordResultMetrics(compiler jscomp.AbstractCompiler, result *jscomp.Result) {
	if compiler != nil && result != nil {
		m.events = append(m.events, "result")
		m.hadResult = true
	}
}

type customRunner struct {
	cli.CommandLineRunner
	doRunCalled            bool
	setRunOptionsCalled    bool
	shouldGenMapPerChunk   bool
	writeOutputCalls       int
	customPrefixOnWriteOut string
}

func newCustomRunner(args []string, out, errStream *jutil.PrintStream) *customRunner {
	r := &customRunner{customPrefixOnWriteOut: "/*custom*/"}
	r.InitCommandLineRunnerWithStreams(r, args, strings.NewReader(""), out, errStream)
	return r
}

func (r *customRunner) DoRun() int {
	r.doRunCalled = true
	return r.CommandLineRunner.DoRun()
}

func (r *customRunner) SetRunOptions(options *jscomp.CompilerOptions) {
	r.setRunOptionsCalled = true
	r.CommandLineRunner.SetRunOptions(options)
}

func (r *customRunner) ShouldGenerateMapPerChunk(options *jscomp.CompilerOptions) bool {
	r.shouldGenMapPerChunk = true
	return r.CommandLineRunner.ShouldGenerateMapPerChunk(options)
}

func (r *customRunner) WriteOutput(
	out jutil.Appendable,
	compiler *jscomp.Compiler,
	licenseTracker jscomp.CodePrinter_LicenseTracker,
	chunk *jscomp.JSChunk,
	wrapper string,
	codePlaceholder string,
	escaper func(string) string,
	filename string,
) {
	r.writeOutputCalls++
	out.AppendCharSequence(r.customPrefixOnWriteOut)
	r.CommandLineRunner.WriteOutput(out, compiler, licenseTracker, chunk, wrapper, codePlaceholder, escaper, filename)
}

func TestCustomCommandLineRunnerHooksAndMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	inA := filepath.Join(tmpDir, "a.js")
	inB := filepath.Join(tmpDir, "b.js")
	if err := os.WriteFile(inA, []byte("var a = 1 + 2;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(inB, []byte("var b = a + 3;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	chunkPrefix := filepath.Join(tmpDir, "chunk_")

	var outBuf, errBuf bytes.Buffer
	outPS := jutil.NewPrintStream(&outBuf)
	errPS := jutil.NewPrintStream(&errBuf)

	r := newCustomRunner([]string{
		"--compilation_level=SIMPLE",
		"--js=" + inA,
		"--chunk=c1:1",
		"--js=" + inB,
		"--chunk=c2:1:c1",
		"--chunk_output_path_prefix=" + chunkPrefix,
		"--create_source_map=%outname%.map",
	}, outPS, errPS)
	cli.DisablePhaseLogger()

	if r.GetErrorPrintStream() != errPS {
		t.Fatal("GetErrorPrintStream mismatch")
	}

	rec := &testMetricsRecorder{}
	r.SetCompileMetricsRecorder(rec)
	if r.GetCompileMetricsRecorder() != rec {
		t.Fatal("GetCompileMetricsRecorder mismatch")
	}

	var exitCode int
	prevExit := jutil.SetSystemExitHook(func(status int) {
		exitCode = status
	})
	defer jutil.SetSystemExitHook(prevExit)

	r.Run()

	if exitCode != 0 {
		t.Fatalf("unexpected exitCode=%d stderr=%s", exitCode, errBuf.String())
	}
	if !r.doRunCalled || !r.setRunOptionsCalled || !r.shouldGenMapPerChunk {
		t.Fatalf("hook dispatch flags: doRun=%v setRunOptions=%v shouldGenMapPerChunk=%v",
			r.doRunCalled, r.setRunOptionsCalled, r.shouldGenMapPerChunk)
	}
	if r.writeOutputCalls != 2 {
		t.Fatalf("writeOutputCalls = %d, want 2", r.writeOutputCalls)
	}
	c1Bytes, err := os.ReadFile(chunkPrefix + "c1.js")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(c1Bytes), "/*custom*/") {
		t.Fatalf("c1 output missing custom WriteOutput prefix: %q", string(c1Bytes))
	}
	wantEvents := []string{"start", "name:full compile", "state", "result"}
	if !slices.Equal(rec.events, wantEvents) || !rec.hadResult {
		t.Fatalf("metrics events = %v, want %v", rec.events, wantEvents)
	}
}

func TestCommandLineConfigEmbedderSetters(t *testing.T) {
	cfg := cli.NewAbstractCommandLineRunner_CommandLineConfig()
	cfg.SetJsOutputFile("out.js").
		SetChunkOutputFiles([]string{"c1:out_c1.js"}).
		SetChunkConformanceFiles([]string{"c1:out_c1.conformance"}).
		SetChunkOutputPathPrefix("dist/").
		SetCreateSourceMap("out.js.map").
		SetSourceMapDetailLevel(jscomp.SourceMap_DetailLevel_SYMBOLS).
		SetSourceMapFormat(jscomp.SourceMap_Format_V3).
		SetOutputWrapper("(function(){%output%})();").
		SetDependencyOptions(jscomp.DependencyOptions_None()).
		SetSaveCompilationStateOpt(base.SomeString("state.bin"), 1).
		SetContinueSavedCompilationOpt(base.SomeString("state.bin"), 1).
		SetTypedAstListInputFilenameOpt(base.SomeString("libs.typedast"))

	if got := cfg.GetJsOutputFile(); got != "out.js" {
		t.Fatalf("GetJsOutputFile = %q, want out.js", got)
	}
	if got, ok := cfg.GetTypedAstListInputFilename(); !ok || got != "libs.typedast" {
		t.Fatalf("GetTypedAstListInputFilename = (%q, %v), want (libs.typedast, true)", got, ok)
	}

	var baseRunner cli.AbstractCommandLineRunnerBase
	var outBuf, errBuf bytes.Buffer
	dummy := newCustomRunner(nil, jutil.NewPrintStream(&outBuf), jutil.NewPrintStream(&errBuf))
	baseRunner.InitAbstractCommandLineRunnerBase(dummy, strings.NewReader(""), jutil.NewPrintStream(&outBuf), jutil.NewPrintStream(&errBuf))
	comp := jscomp.NewCompiler()
	baseRunner.SetCompiler(comp)
	if baseRunner.GetCompiler() != comp {
		t.Fatal("SetCompiler did not update GetCompiler")
	}

	chunk := jscomp.NewJSChunk("app")
	dummy.GetCommandLineConfig().SetChunkConformanceFiles([]string{"app:build/app.conformance"})
	if got := dummy.GetChunkConformanceFileName(chunk); got != "build/app.conformance" {
		t.Fatalf("GetChunkConformanceFileName = %q, want build/app.conformance", got)
	}
}

func TestPersistentWorkerInterceptsSystemExitAndRewindsBaseline(t *testing.T) {
	tmpDir := t.TempDir()
	inPath := filepath.Join(tmpDir, "in.js")
	outPath := filepath.Join(tmpDir, "out.js")
	if err := os.WriteFile(inPath, []byte("const x = 40 + 2.5; console.log(x);\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	w := cli.NewPersistentWorker(nil)
	for i := range 3 {
		resp := w.ExecuteRequest(&cli.WorkRequest{
			RequestId: i + 1,
			Arguments: []string{
				"--compilation_level=SIMPLE",
				"--js=" + inPath,
				"--js_output_file=" + outPath,
			},
		})
		if resp.ExitCode != 0 {
			t.Fatalf("req %d failed: exit=%d output=%s", i+1, resp.ExitCode, resp.Output)
		}
	}

	// Also verify invalid flag / error exit code does not kill the worker process:
	badResp := w.ExecuteRequest(&cli.WorkRequest{
		RequestId: 99,
		Arguments: []string{"--not_a_real_closure_flag_xyz"},
	})
	if badResp.ExitCode == 0 {
		t.Fatal("expected non-zero exit code for invalid flag in worker")
	}
}
