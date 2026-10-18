package parsing

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func builtinCacheFingerprint(n *rhino.Node) string {
	var out strings.Builder
	var visit func(*rhino.Node)
	visit = func(n *rhino.Node) {
		out.WriteString(n.String())
		if doc := n.GetJSDocInfo(); doc != nil {
			out.WriteString(doc.ToStringVerbose())
		}
		out.WriteByte('\n')
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			visit(child)
		}
	}
	visit(n)
	return out.String()
}

func TestBuiltinExternCacheIsolation(t *testing.T) {
	ParserRunner_EnableBuiltinExternCache(true)
	defer ParserRunner_EnableBuiltinExternCache(false)
	config := ParserRunner_CreateConfig(Config_LanguageMode_ES_NEXT, nil, Config_StrictMode_SLOPPY)
	source := "function f(/** number */ x) {}\n/** @type {string} */ var value;"
	file := rhino.NewSimpleSourceFile("externs.zip//cache-test.js", rhino.StaticSourceFile_SourceKind_EXTERN)
	original := parserRunnerParseUncached(file, source, config, rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
	want := builtinCacheFingerprint(original.Ast)
	first := ParserRunner_Parse(file, source, config, rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
	if builtinCacheFingerprint(first.Ast) != want {
		t.Fatal("first parse differs")
	}
	first.Ast.DetachChildren()
	for i := 0; i < 3; i++ {
		current := rhino.NewSimpleSourceFile(file.GetName(), rhino.StaticSourceFile_SourceKind_EXTERN)
		got := ParserRunner_Parse(current, source, config, rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
		if builtinCacheFingerprint(got.Ast) != want {
			t.Fatal("cached AST or inline JSDoc differs")
		}
		var check func(*rhino.Node)
		check = func(n *rhino.Node) {
			if n.GetStaticSourceFile() != current {
				t.Fatal("source metadata was shared across compilations")
			}
			if d := n.GetJSDocInfo(); d != nil {
				for _, typ := range d.GetTypeNodes() {
					check(typ)
				}
			}
			for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
				check(child)
			}
		}
		check(got.Ast)
		got.Ast.DetachChildren()
	}
	changed := ParserRunner_Parse(file, "/** @type {number} */ var other;", config, rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
	if builtinCacheFingerprint(changed.Ast) == want {
		t.Fatal("source-content change reused cached AST")
	}
}
