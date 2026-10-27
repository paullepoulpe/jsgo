package cli

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

	"github.com/paullepoulpe/jsgo/internal/jscomp"
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
		// CommandLineRunner.main, with the exit status that the
		// SystemExitCodeReceiver would pass to System.exit, as the shell
		// sees it (tools/oracle/java/CliE2E).
		runner := NewCommandLineRunnerWithOutAndErr(
			c.args, jutil.NewPrintStream(&out), jutil.NewPrintStream(&errOut))
		exitStatus := func(code int) int {
			b := int8(int32(code))
			if b == 0 && code != 0 {
				b = -1
			}
			return int(uint8(b))
		}
		status := 0
		if runner.ShouldRunCompiler() {
			runner.SetExitCodeReceiver(func(code int) { status = exitStatus(code) })
			runner.Run()
		} else if runner.HasErrors() {
			status = exitStatus(-1)
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
// to restore with a NullPointerException, as Java's does since v20261015:
// JSCompilerStateProto.parseDelimitedFrom returns null at the end of the
// stream, and deserializeCompilerState dereferences it.
func TestCompilerRestoreState_Truncated(t *testing.T) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Close()
	options := jscomp.NewCompilerOptions()
	compiler := jscomp.NewCompiler()
	compiler.Init(nil, []*jscomp.SourceFile{jscomp.SourceFile_FromCode("in.js", "var x;")}, options)
	var caught jutil.Throwable
	jutil.Try(func() { compiler.RestoreState(&buf) }, func(e jutil.IRuntimeException) { caught = e })
	if caught == nil || !strings.Contains(jutil.ToJavaString(caught), "java.lang.NullPointerException") {
		t.Fatalf("RestoreState of an empty stream: got %v, want a NullPointerException", caught)
	}
}

// An ES module that re-exports the default import of a CommonJS module
// crashes Es6RewriteModules in Java too (an INTERNAL COMPILER ERROR from
// checkState(boundVariableQualifiedName.isSimple(), "unexpected qualified
// name: %s", ...)). The message embeds a QualifiedName, which has Java's
// Object.toString(); the port printed Go's %v of the struct. The hash is
// Java's identity hash, which changes from run to run, so only the form is
// compared.
func TestCliParity_UnexpectedQualifiedNameMessage(t *testing.T) {
	t.Chdir(t.TempDir())
	os.WriteFile("lib.js", []byte("module.exports = {x: 1};\n"), 0o666)
	os.WriteFile("es.js", []byte("import Lib from './lib.js';\nexport { Lib };\n"), 0o666)
	var out, errOut bytes.Buffer
	runner := NewCommandLineRunnerWithOutAndErr(
		[]string{"--process_common_js_modules", "--module_resolution=NODE", "--dependency_mode=PRUNE",
			"--entry_point=es.js", "--js", "lib.js", "--js", "es.js"},
		jutil.NewPrintStream(&out), jutil.NewPrintStream(&errOut))
	status := 0
	runner.SetExitCodeReceiver(func(s int) { status = s })
	runner.Run()
	const want = "java.lang.RuntimeException: INTERNAL COMPILER ERROR.\nPlease report this problem.\n\n" +
		"unexpected qualified name: com.google.javascript.rhino.QualifiedName$GetpropQname@"
	if got := errOut.String(); !strings.HasPrefix(got, want) {
		t.Errorf("stderr:\n%s\nwant a prefix:\n%s", got, want)
	}
	if status != -2 || out.Len() != 0 {
		t.Errorf("exit status %d, stdout %q; want -2 and nothing", status, out.String())
	}
}

// A super in an object literal's method, transpiled to ES5: since
// v20261015, Es6ConvertSuper reports it as JSC_CANNOT_CONVERT (it used to find
// no enclosing class and dereference null).
func TestCliParity_SuperInObjectLiteralTranspiled(t *testing.T) {
	cliParityCase{
		files: map[string]string{"in.js": "({ a() { super.b(); } });\n"},
		args:  []string{"--language_out=ECMASCRIPT5", "--js", "in.js"},
		wantErr: "in.js:1:1: WARNING - [JSC_USELESS_CODE] Suspicious code. This code lacks side-effects. Is there a bug?\n" +
			"  1| ({ a() { super.b(); } });\n" +
			"      ^^^^^^^^^^^^^^^^^^^^^^\n" +
			"\n" +
			"in.js:1:9: ERROR - [JSC_CANNOT_CONVERT] This code cannot be transpiled. super is not supported in object literal methods.\n" +
			"  1| ({ a() { super.b(); } });\n" +
			"              ^^^^^\n" +
			"\n" +
			"1 error(s), 1 warning(s)\n",
		wantExit: 1,
	}.run(t)
}

// A function type as the type argument of an @extends, mentioning the
// subclass's own template: Java resolves FnBox<string>'s T to
// function(string): string and reports both mismatches. The template map's
// fast path used to return function types unchanged (they matched the
// PrototypeObjectType case), leaving T as function(U): U, and the port
// reported nothing.
func TestCliParity_FunctionTypeTemplateArgumentOfSupertype(t *testing.T) {
	cliParityCase{
		files: map[string]string{"t2.js": `/** @constructor @template T */
function Box() {}
/** @constructor @extends {Box<function(U): U>} @template U */
function FnBox() { Box.call(this); }
/** @type {!FnBox<string>} */
var b = new FnBox();
/** @type {!Box<function(number): number>} */
var x = b;
/** @type {!Box<function(string): string>} */
var y = b;
/** @param {!Box<function(number): number>} p */
function take(p) {}
take(b);
`},
		args: []string{"--checks_only", "--jscomp_warning=checkTypes", "--js", "t2.js"},
		wantErr: "t2.js:8:8: WARNING - [JSC_TYPE_MISMATCH] initializing variable\nfound   : FnBox<string>\nrequired: Box<function(number): number>\n  8| var x = b;\n             ^\n\n" +
			"t2.js:13:5: WARNING - [JSC_TYPE_MISMATCH] actual parameter 1 of take does not match formal parameter\nfound   : FnBox<string>\nrequired: Box<function(number): number>\n  13| take(b);\n           ^\n\n" +
			"0 error(s), 2 warning(s), 100.0% typed\n",
	}.run(t)
}
