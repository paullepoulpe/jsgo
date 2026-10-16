package jscomp

// Go-only: test-helper — the parsing half of Compiler.parseTestCode(String)
// for the tests of core group g2 (AstAnalyzerTest, ...), which only need
// the AST of a snippet: Compiler (above the core) is not ported yet.
//
// g2ParseTestCode parses js as Compiler.parseTestCode does for a Compiler
// whose CompilerOptions set the given languageIn and strictModeInput: the
// parsing.Config is the one Compiler.getParserConfig builds
// (ParserRunner.createConfig with the CompilerOptions defaults for
// JsDocParsing, RunMode, extra annotation names and inline source maps),
// and the source is SourceFile.fromCode("testcode", js). It returns the
// SCRIPT node and fails the test on any parse error, as the Java tests'
// `assertThat(compiler.getErrors()).isEmpty()` does. The properties that
// CompilerInput/JsAst add to the SCRIPT node afterwards (input id, feature
// set) are not added; no g2 test reads them.

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type g2CollectingReporter struct {
	errors []string
}

func (r *g2CollectingReporter) Warning(message string, sourceName string, line int, lineOffset int) {}

func (r *g2CollectingReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.errors = append(r.errors, message)
}

func g2ParseTestCode(t testing.TB, js string, languageIn parsing.Config_LanguageMode, strictModeInput bool) *rhino.Node {
	t.Helper()
	strictMode := parsing.Config_StrictMode_SLOPPY
	if strictModeInput {
		strictMode = parsing.Config_StrictMode_STRICT
	}
	config := parsing.ParserRunner_CreateConfigWithModes(
		languageIn,
		parsing.Config_JsDocParsing_TYPES_ONLY,
		parsing.Config_RunMode_STOP_AFTER_ERROR,
		nil,
		true,
		strictMode)
	reporter := &g2CollectingReporter{}
	result := parsing.ParserRunner_Parse(SourceFile_FromCode("testcode", js), js, config, reporter)
	if len(reporter.errors) != 0 {
		t.Fatalf("parse errors for %q: %v", js, reporter.errors)
	}
	if result.Ast == nil {
		t.Fatalf("no AST for %q", js)
	}
	return result.Ast
}
