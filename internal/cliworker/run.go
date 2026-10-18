// Package cliworker runs Closure's command-line front end repeatedly in one
// process. Each request creates a fresh compiler; only immutable builtin extern
// parse templates may be reused. User JavaScript is compiled on every request.
package cliworker

import (
	"bytes"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Result contains the exact bytes and shell exit status of one CLI invocation.
type Result struct {
	Exit   int    `json:"exit"`
	Stdout []byte `json:"stdout_base64"`
	Stderr []byte `json:"stderr_base64"`
}

// System streams and the ported compiler's thread-local adapters are process
// state. Requests are deliberately serialized, including callers of this API.
var runMu sync.Mutex

// Run invokes the same CommandLineRunner as jsgo. File outputs are written by
// the compiler normally. stdin supplies --js - and JSON-stream input.
func Run(args []string, stdin []byte) (result Result) {
	runMu.Lock()
	defer runMu.Unlock()
	var stdout, stderr bytes.Buffer
	outStream := jutil.NewPrintStream(&stdout)
	errStream := jutil.NewPrintStreamWithCharset(&stderr, jutil.System_NativeCharset())
	oldOut, oldErr := jutil.System_Out(), jutil.System_Err()
	jutil.System_SetOut(outStream)
	jutil.System_SetErr(errStream)
	defer func() {
		if value := recover(); value != nil {
			throwable, ok := value.(jutil.Throwable)
			if !ok {
				throwable = jutil.AsThrowable(value)
			}
			errStream.Print("Exception in thread \"main\" ")
			jutil.Throwable_PrintStackTrace(throwable, errStream)
			result.Exit = 1
		}
		outStream.Flush()
		errStream.Flush()
		result.Stdout, result.Stderr = stdout.Bytes(), stderr.Bytes()
		if result.Stdout == nil {
			result.Stdout = []byte{}
		}
		if result.Stderr == nil {
			result.Stderr = []byte{}
		}
		jutil.System_SetOut(oldOut)
		jutil.System_SetErr(oldErr)
	}()
	runner := jscomp.NewCommandLineRunnerWithStreams(args, bytes.NewReader(stdin), outStream, errStream)
	if runner.ShouldRunCompiler() {
		runner.SetExitCodeReceiver(func(code int) { result.Exit = int(uint8(code)) })
		runner.Run()
	} else if runner.HasErrors() {
		result.Exit = 255
	}
	return
}
