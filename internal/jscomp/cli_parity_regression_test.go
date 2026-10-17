package jscomp

// Go-only: end-to-end regression tests for crashes and hangs found by
// comparing the jsgo CLI with the pinned Java compiler (tools/clie2e). Each
// case runs CommandLineRunner in a fresh temporary directory, as
// `java -jar closure-compiler.jar <args>` would run there, and expects the
// Java compiler's stdout, stderr and exit status byte for byte.

import (
	"bytes"
	"compress/gzip"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type cliParityCase struct {
	files    map[string]string // inputs, relative to the working directory
	args     []string
	wantOut  string
	wantErr  string
	wantExit int
}

// run runs the case's command line in a fresh temporary directory, failing
// the test if it does not finish within a minute (a hang or a deadlock in
// the port).
func (c cliParityCase) run(t *testing.T) {
	t.Helper()
	t.Chdir(t.TempDir())
	c.runHere(t)
}

// runHere runs the case in the current directory, after writing its files.
func (c cliParityCase) runHere(t *testing.T) {
	t.Helper()
	for name, content := range c.files {
		if err := os.WriteFile(name, []byte(content), 0o666); err != nil {
			t.Fatal(err)
		}
	}
	var out, errOut bytes.Buffer
	exit := make(chan int, 1)
	go func() {
		runner := NewCommandLineRunnerWithOutAndErr(
			c.args, jutil.NewPrintStream(&out), jutil.NewPrintStream(&errOut))
		status := 0
		if runner.ShouldRunCompiler() {
			status = runner.DoRun()
		}
		exit <- status
	}()
	var status int
	select {
	case status = <-exit:
	case <-time.After(time.Minute):
		t.Fatalf("jsgo %s did not finish within a minute", strings.Join(c.args, " "))
	}
	if got := out.String(); got != c.wantOut {
		t.Errorf("stdout:\n%s\nwant:\n%s", got, c.wantOut)
	}
	if got := errOut.String(); got != c.wantErr {
		t.Errorf("stderr:\n%s\nwant:\n%s", got, c.wantErr)
	}
	if status != c.wantExit {
		t.Errorf("exit status %d, want %d", status, c.wantExit)
	}
}

// A diagnostic in a file whose sourceMappingURL names a missing map: while
// the synchronized ThreadSafeDelegatingErrorManager.generateReport formats
// it, the compiler looks the position up in that input source map, fails to
// load it and reports SOURCEMAP_RESOLVE_FAILED through the same error
// manager. Java's monitors are re-entrant (jutil.Monitor); a sync.Mutex
// deadlocked here.
func TestCliParity_ReportWhileGeneratingReport(t *testing.T) {
	cliParityCase{
		files: map[string]string{"in.js": "var x = y;\n//# sourceMappingURL=in.js.map\n"},
		args:  []string{"--jscomp_warning=undefinedVars", "--js", "in.js"},
		wantErr: "in.js:1:8: WARNING - [JSC_UNDEFINED_VARIABLE] variable y is undeclared\n" +
			"  1| var x = y;\n" +
			"             ^\n" +
			"\n" +
			"0 error(s), 2 warning(s)\n",
		wantOut:  "var x=y;\n",
		wantExit: 0,
	}.run(t)
}

// A multistage build whose stages run in separate processes: the compiler
// state that --filename_to_save_to writes must hold the state itself. The
// port kept the CompilerState in a process-wide table and wrote only its
// key, so a later process (or a second restore of the same file) failed
// with "java.io.IOException: no saved CompilerState 1". Here the first
// stage's state is restored twice, and each run must match Java's output.
func TestCliParity_MultistageStateInFile(t *testing.T) {
	t.Chdir(t.TempDir())
	base := []string{"--js", "in.js", "--compilation_level", "ADVANCED"}
	stage := func(args ...string) []string { return append(append([]string{}, base...), args...) }
	cliParityCase{
		files: map[string]string{"in.js": "function f(a){return a+1}\nwindow.x=f(2);\n"},
		args:  stage("--segment_of_compilation_to_run=CHECKS", "--filename_to_save_to=checks.state"),
	}.runHere(t)
	for i := 0; i < 2; i++ {
		cliParityCase{
			args: stage("--segment_of_compilation_to_run=OPTIMIZATIONS",
				"--filename_to_restore_from=checks.state", "--filename_to_save_to=optimizations.state"),
		}.runHere(t)
		cliParityCase{
			args:    stage("--segment_of_compilation_to_run=FINALIZATIONS", "--filename_to_restore_from=optimizations.state"),
			wantOut: "window.x=3;\n",
		}.runHere(t)
	}
}

// A compiler state file that is not one (here, an empty gzip stream) fails
// to restore with an EOFException, as Java's ObjectInputStream does on a
// truncated stream.
func TestCompilerRestoreState_Truncated(t *testing.T) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Close()
	options := NewCompilerOptions()
	compiler := NewCompiler()
	compiler.Init(nil, []*SourceFile{SourceFile_FromCode("in.js", "var x;")}, options)
	var caught jutil.Throwable
	jutil.Try(func() { compiler.RestoreState(&buf) }, func(e jutil.IRuntimeException) { caught = e })
	if caught == nil || !strings.Contains(jutil.ToJavaString(caught), "java.io.EOFException") {
		t.Fatalf("RestoreState of an empty stream: got %v, want an EOFException", caught)
	}
}
