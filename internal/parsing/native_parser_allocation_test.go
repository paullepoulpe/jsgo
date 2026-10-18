package parsing

import (
	"runtime"
	"sync"
	"testing"
	"weak"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func rejectedNativeGraph(config *Config, allocator *rhino.NodeAllocator) []weak.Pointer[rhino.Node] {
	source := "var f=function(a){return {value:a};}; let unsupported=1;"
	file := parser.NewSourceFile("rejected.js", source)
	diagnostics := &nativeParserDiagnostics{}
	diagnostics.InitErrorReporterBase(diagnostics)
	p := &nativeParser{diagnostics: diagnostics}
	p.scanner = parser.NewScanner(diagnostics, diagnostics, file, 0)
	p.factory = newIRFactoryWithNodeAllocator(
		rhino.NewSimpleSourceFile("rejected.js", rhino.StaticSourceFile_SourceKind_STRONG),
		config, diagnostics, nil, file, JsDocInfoParser_JsDocSourceKind_NORMAL, allocator)
	p.advance()
	root := p.factory.NewNode(rhino.Token_SCRIPT)
	first := p.statement()
	if p.failed() {
		panic("retention fixture rejected its supported prefix")
	}
	root.AddChildToBack(first.node)
	p.statement()
	if !p.failed() {
		panic("retention fixture did not exercise speculative rejection")
	}
	return []weak.Pointer[rhino.Node]{weak.Make(root), weak.Make(p.factory.defaultTemplateNode),
		weak.Make(p.factory.closureUnawareTemplateNode)}
}

func TestNativeParserRejectedGraphAndTemplatesDoNotShareLiveFallback(t *testing.T) {
	config := ParserRunner_CreateConfig(Config_LanguageMode_ES_NEXT, nil, Config_StrictMode_STRICT)
	before := rhino.NewNode(rhino.Token_EMPTY)
	discarded := rejectedNativeGraph(config, &rhino.NodeAllocator{})
	after := rhino.NewNode(rhino.Token_EMPTY)
	file := rhino.NewSimpleSourceFile("fallback.js", rhino.StaticSourceFile_SourceKind_STRONG)
	fallback := parserRunnerParseUncached(file, "let live=1;", config, &nativeExperimentReporter{})
	for attempt := 0; attempt < 10; attempt++ {
		runtime.GC()
		allGone := true
		for _, node := range discarded {
			if node.Value() != nil {
				allGone = false
			}
		}
		if allGone {
			runtime.KeepAlive(before)
			runtime.KeepAlive(after)
			runtime.KeepAlive(fallback)
			return
		}
	}
	t.Fatal("a rejected native graph or template shares the live global allocation stream")
}

func TestNativeParserConcurrentAcceptedAndFallback(t *testing.T) {
	config := ParserRunner_CreateConfig(Config_LanguageMode_ES_NEXT, nil, Config_StrictMode_STRICT)
	sources := []string{
		"function f(a){if(a) return {value:a}; return null;} f(1);",
		"var f=function(a){return {value:a};}; let fallback=1;",
	}
	file := rhino.NewSimpleSourceFile("parallel.js", rhino.StaticSourceFile_SourceKind_STRONG)
	expected := make([]string, len(sources))
	for i, source := range sources {
		expected[i] = parserRunnerParseUncached(file, source, config, &nativeExperimentReporter{}).Ast.ToStringTree()
	}
	var wait sync.WaitGroup
	for range 4 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for range 16 {
				for i, source := range sources {
					_, accepted := tryNativeParser(file, source, config)
					if accepted != (i == 0) {
						t.Errorf("unexpected acceptance for fixture %d: %v", i, accepted)
						return
					}
					reporter := &nativeExperimentReporter{}
					got := parserRunnerParseWithNativeFallback(file, source, config, reporter)
					if got.Ast == nil || got.Ast.ToStringTree() != expected[i] || len(reporter.diagnostics) != 0 {
						t.Errorf("concurrent accepted/fallback parse changed fixture %d", i)
						return
					}
				}
			}
		}()
	}
	wait.Wait()
}

// Negative control: with the old allocation mode, this same rejected prefix
// shares a global chunk with a live anchor and retains the whole parent graph.
// This proves the weak-reference regression would detect the original issue.
func TestNativeParserRetentionOracleDetectsGlobalSharing(t *testing.T) {
	config := ParserRunner_CreateConfig(Config_LanguageMode_ES_NEXT, nil, Config_StrictMode_STRICT)
	before := rhino.NewNode(rhino.Token_EMPTY)
	discarded := rejectedNativeGraph(config, nil)
	after := rhino.NewNode(rhino.Token_EMPTY)
	runtime.GC()
	if discarded[0].Value() == nil {
		t.Fatal("retention fixture did not detect the original shared-chunk issue")
	}
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}
