package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: language-adapter — parse files or deserialize TypedAST script roots
// independently, then publish roots, diagnostics, comments and source maps in
// submission order. Compiler maps and error managers are touched only by the
// calling goroutine. Unusual/shared SourceFiles keep the original sequential
// path.
func (p *PrebuildAst) prebuildParallel(inputs []*CompilerInput) bool {
	compiler, ok := p.compiler.(*Compiler)
	if !ok || compiler == nil || len(inputs) < 2 {
		return false
	}
	seen := make(map[*SourceFile]bool, len(inputs))
	for _, input := range inputs {
		if input == nil || input.ast == nil || seen[input.sourceFile] {
			return false
		}
		seen[input.sourceFile] = true
	}
	var normal, extern *parsing.Config
	needParserConfig := compiler.typedAstFilesystem == nil
	type outcome struct {
		scriptDeserializer *ScriptNodeDeserializer
		deserializer       func() *rhino.Node
		deserializedRoot   *rhino.Node
		parsed             *parsing.ParserRunner_ParseResult
		errors             parseBufferedErrors
		failure            jutil.Throwable
		wasDeserialized    bool
	}
	outcomes := make([]outcome, len(inputs))
	if compiler.typedAstFilesystem != nil {
		for index, input := range inputs {
			if input.ast.isParsed() {
				continue
			}
			jutil.Try(func() {
				sd, fn := compiler.takeTypedAstDeserializer(input.sourceFile)
				if fn != nil {
					if sd != nil {
						outcomes[index].scriptDeserializer = sd
					} else {
						outcomes[index].deserializer = fn
					}
					outcomes[index].wasDeserialized = true
				} else {
					needParserConfig = true
				}
			}, func(t jutil.Throwable) {
				outcomes[index].failure = t
			})
		}
	}
	if needParserConfig {
		normal = compiler.GetParserConfig(AbstractCompiler_ConfigContext_DEFAULT)
		extern = compiler.GetParserConfig(AbstractCompiler_ConfigContext_EXTERNS)
	}
	failures := parallelPrebuildWorkers(len(inputs), p.numParallelThreads, func(claim func() (int, bool), recordFailure func(int, jutil.Throwable)) {
		jsdocSlabs := rhino.NewJSDocBuilderSlabs()
		for {
			index, ok := claim()
			if !ok {
				return
			}
			jutil.Try(func() {
				input, result := inputs[index], &outcomes[index]
				if input.ast.isParsed() || result.failure != nil {
					result.scriptDeserializer = nil
					result.deserializer = nil
					return
				}
				if result.wasDeserialized {
					if sd := result.scriptDeserializer; sd != nil {
						result.scriptDeserializer = nil
						result.deserializedRoot = sd.DeserializeNew()
						return
					}
					fn := result.deserializer
					result.deserializer = nil
					result.deserializedRoot = fn()
					return
				}
				config := normal
				if input.sourceFile.IsExtern() {
					config = extern
				}
				result.parsed = parsing.ParserRunner_ParseWithSlabs(
					input.sourceFile, input.sourceFile.GetCode(), config, &result.errors, &jsdocSlabs)
			}, func(t jutil.Throwable) {
				outcomes[index].scriptDeserializer = nil
				outcomes[index].deserializer = nil
				recordFailure(index, t)
			})
		}
	})
	for i, f := range failures {
		if f != nil && outcomes[i].failure == nil {
			outcomes[i].failure = f
		}
	}
	var failure jutil.Throwable
	reporter := compiler.GetDefaultErrorReporter()
	for index, input := range inputs {
		if input.ast.isParsed() {
			continue
		}
		result := &outcomes[index]
		jutil.Try(func() {
			if result.wasDeserialized {
				if result.failure != nil {
					panic(result.failure)
				}
				input.ast.installDeserializedRoot(result.deserializedRoot)
			} else {
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
			}
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
