// Port of com/google/javascript/jscomp/parsing/ParserRunner.java
// @ closure-compiler v20260415
//
// ParserRunner wires the new 1:1 parser (internal/parser/ParserRewrite) through
// IRFactory to produce a rhino.Node SCRIPT root, mirroring the Java call sequence:
//
//	Parser p = new Parser(es6config, es6ErrorReporter, file);
//	ProgramTree tree = p.parseProgram();
//	IRFactory factory = IRFactory.transformTree(tree, sourceFile, config, errorReporter, file);
//	root = factory.getResultNode();
//
// The returned node is a SCRIPT-rooted rhino.Node, ready for downstream passes.
//
// Note: this file lives in internal/jscomp rather than internal/jscomp/parsing
// to avoid an import cycle: internal/parser/parser.go already imports
// internal/jscomp/parsing, so internal/jscomp/parsing cannot import internal/parser.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/legacy/internal/jscomp/parsing"
	"github.com/paullepoulpe/jsgo/legacy/internal/parser"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// parseResult mirrors ParserRunner.ParseResult. It holds the SCRIPT-rooted
// rhino.Node produced by the new pipeline.
type parseResult struct {
	// ast is the SCRIPT-rooted rhino.Node produced by Parser → IRFactory.
	ast *rhino.Node
	// hadError reports whether the parser or IRFactory encountered a hard error.
	hadError bool
	// firstError is the first error message (empty if hadError==false).
	firstError string
}

// pipelineErrorReporter satisfies both the parser.ErrorReporter interface
// (used by Scanner/ParserRewrite) and the parsing-internal errorReporter
// interface (used by IRFactory). We embed the method sets needed by each.
type pipelineErrorReporter struct {
	hadError   bool
	firstError string
}

// --- parser.ErrorReporter interface (used by Scanner / ParserRewrite) ---

func (r *pipelineErrorReporter) ReportError(position parser.SourcePosition, format string, args ...interface{}) {
	if !r.hadError {
		r.firstError = fmt.Sprintf(format, args...)
		r.hadError = true
	}
}

func (r *pipelineErrorReporter) ReportWarning(position parser.SourcePosition, format string, args ...interface{}) {
	// Warnings are not promoted to errors in the Go pipeline. Mirrors Java's
	// Es6ErrorReporter.reportWarning, which delegates to ErrorReporter.warning.
	_ = position
}

// --- parsing.ErrorReporter interface (used by IRFactory via its internal errorReporter type) ---
// IRFactory uses an unexported errorReporter interface with Error/Warning methods.
// We cannot implement it directly since it is unexported; IRFactory.NewIRFactory
// accepts a parsing.ErrorReporter adapter. See irfactoryErrorReporter below.

// irfactoryErrorReporter adapts *pipelineErrorReporter to the unexported
// errorReporter interface expected by parsing.NewIRFactory.
type irfactoryErrorReporter struct {
	inner *pipelineErrorReporter
}

func (r *irfactoryErrorReporter) Error(msg, sourceName string, lineno, charno int) {
	r.inner.ReportError(parser.SourcePosition{}, "%s", msg)
}

func (r *irfactoryErrorReporter) Warning(msg, sourceName string, lineno, charno int) {
	// Warnings are not promoted to errors; mirrors Es6ErrorReporter.reportWarning.
	_ = msg
}

// Go-only: package-bridge
// runNewPipeline runs the new 1:1 parse pipeline:
//
//	ParserRewrite → ProgramTree → IRFactory → rhino.Node
//
// It mirrors ParserRunner.parse(sourceFile, sourceString, config, errorReporter).
// The parser mode is ES8OrGreater (matching Closure's ECMASCRIPT_2021+ default).
func runNewPipeline(sourceName, src string) parseResult {
	er := &pipelineErrorReporter{}
	irEr := &irfactoryErrorReporter{inner: er}

	// Step 1: Configure and run the new parser.
	// Mirrors: new Parser(es6config, es6ErrorReporter, file); p.parseProgram();
	config := parser.NewParserConfig(parser.ParserModeES8OrGreater, false)
	p := parser.NewParserRewrite(config, er, src)
	programTree := p.Parse()

	if programTree == nil {
		return parseResult{hadError: true, firstError: "parser returned nil ProgramTree"}
	}

	// Step 2: Run IRFactory.Transform on the ProgramTree.
	// Mirrors: IRFactory factory = IRFactory.transformTree(tree, sourceFile, ...);
	//          root = factory.getResultNode();
	sourceFile := rhino.NewSimpleSourceFile(sourceName, rhino.SourceKindStrong)
	factory := parsing.NewIRFactory(sourceFile, irEr)
	// Wire source comments so IRFactory can attach JSDoc to each node.
	// Mirrors Java IRFactory.transformTree receiving the comment list from ParserRunner.
	factory.SetSourceComments(programTree.SourceComments, src)
	root := factory.Transform(&programTree.ParseTree)

	return parseResult{
		ast:        root,
		hadError:   er.hadError,
		firstError: er.firstError,
	}
}
