package jscomp

import (
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: language-adapter — parse files independently, then publish roots,
// diagnostics, comments and source maps in submission order. Compiler maps and
// error managers are touched only by the calling goroutine. TypedAST loading
// and unusual/shared SourceFiles keep the original sequential path.
func (p *PrebuildAst) prebuildParallel(inputs []*CompilerInput) bool {
	compiler, ok := p.compiler.(*Compiler)
	if !ok || compiler == nil || compiler.typedAstFilesystem != nil || len(inputs) < 2 {
		return false
	}
	seen := make(map[*SourceFile]bool, len(inputs))
	for _, input := range inputs {
		if input == nil || input.ast == nil || seen[input.sourceFile] {
			return false
		}
		seen[input.sourceFile] = true
	}
	normal := compiler.GetParserConfig(AbstractCompiler_ConfigContext_DEFAULT)
	extern := compiler.GetParserConfig(AbstractCompiler_ConfigContext_EXTERNS)
	type outcome struct {
		parsed  *parsing.ParserRunner_ParseResult
		errors  parseBufferedErrors
		failure jutil.Throwable
	}
	outcomes := make([]outcome, len(inputs))
	work := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < min(p.numParallelThreads, len(inputs)); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range work {
				input, result := inputs[index], &outcomes[index]
				if input.ast.isParsed() {
					continue
				}
				jutil.Try(func() {
					config := normal
					if input.sourceFile.IsExtern() {
						config = extern
					}
					result.parsed = parsing.ParserRunner_Parse(input.sourceFile, input.sourceFile.GetCode(), config, &result.errors)
				}, func(t jutil.Throwable) { result.failure = t })
			}
		}()
	}
	for index := range inputs {
		work <- index
	}
	close(work)
	wg.Wait()
	var failure jutil.Throwable
	reporter := compiler.GetDefaultErrorReporter()
	for index, input := range inputs {
		if input.ast.isParsed() {
			continue
		}
		result := &outcomes[index]
		jutil.Try(func() {
			result.errors.replay(reporter)
			if result.failure != nil {
				if ioerr, ok := result.failure.(jutil.IIOException); ok {
					compiler.Report(JSError_Make(AbstractCompiler_READ_ERROR, input.sourceFile.GetName(), ioerr.GetMessage()))
				} else {
					panic(result.failure)
				}
			} else {
				input.ast.installParseResult(compiler, result.parsed)
			}
			input.ast.ensureScriptRoot()
			jutil.CheckState(input.ast.root.GetStaticSourceFile() == rhino.StaticSourceFile(input.sourceFile))
			input.ast.root.SetInputId(input.id)
			input.sourceFile.ClearCachedSource()
		}, func(t jutil.Throwable) {
			if failure == nil {
				failure = t
			}
		})
	}
	if failure != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewExecutionException(failure)))
	}
	return true
}

type parseBufferedError struct {
	warning         bool
	message, source string
	line, offset    int
}
type parseBufferedErrors struct{ events []parseBufferedError }

func (r *parseBufferedErrors) Error(message, source string, line, offset int) {
	r.events = append(r.events, parseBufferedError{false, message, source, line, offset})
}
func (r *parseBufferedErrors) Warning(message, source string, line, offset int) {
	r.events = append(r.events, parseBufferedError{true, message, source, line, offset})
}
func (r *parseBufferedErrors) replay(target rhino.ErrorReporter) {
	for _, e := range r.events {
		if e.warning {
			target.Warning(e.message, e.source, e.line, e.offset)
		} else {
			target.Error(e.message, e.source, e.line, e.offset)
		}
	}
}
