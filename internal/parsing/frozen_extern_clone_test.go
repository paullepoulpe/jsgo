package parsing

import (
	"archive/zip"
	"bytes"
	"io"
	"runtime"
	"strings"
	"testing"
	"weak"

	"github.com/paullepoulpe/jsgo/internal/resources"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func frozenCloneWalk(n *rhino.Node, visit func(*rhino.Node)) {
	visit(n)
	if d := n.GetJSDocInfo(); d != nil {
		for _, root := range d.GetTypeNodes() {
			frozenCloneWalk(root, visit)
		}
		for root := range d.GetTypeTransformations().Values() {
			frozenCloneWalk(root, visit)
		}
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		frozenCloneWalk(c, visit)
	}
}

func TestFrozenExternCloneDifferentialAndIsolation(t *testing.T) {
	config := ParserRunner_CreateConfig(Config_LanguageMode_ES_NEXT, nil, Config_StrictMode_SLOPPY)
	cases := []string{
		`function f(/** number */ x) {} /** @type {string} */ var value;`,
		`/** @constructor @template T */ function Box() {} /** @param {{x:!Array<string>, y:number=}} x @return {?string} */ Box.prototype.f = function(x) {};`,
		"/** @template T := 'string' =: */ function transform() {}",
		`/** @interface */ function A() {} /** @interface */ function B() {} /** @constructor @implements {A} @implements {B} */ function C() {}`,
		`/** @typedef {{x:number,y:(string|undefined)}} */ var Record; /** @enum {number} */ var Kind = {A:1, B:2};`,
	}
	for _, source := range cases {
		file := rhino.NewSimpleSourceFile("externs.zip//frozen.js", rhino.StaticSourceFile_SourceKind_EXTERN)
		original := parserRunnerParseUncached(file, source, config, rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
		template := cloneBuiltinExternResult(original, file, true)
		frozen := rhino.Node_FreezeParseTree(template.Ast)
		if frozen == nil {
			t.Fatal("supported parse template unexpectedly fell back")
		}
		want := builtinCacheFingerprint(cloneBuiltinExternResult(template, file, false).Ast)
		for iteration := 0; iteration < 3; iteration++ {
			current := rhino.NewSimpleSourceFile(file.GetName(), rhino.StaticSourceFile_SourceKind_EXTERN)
			got := frozen.CloneWithSourceFile(current)
			if builtinCacheFingerprint(got) != want {
				t.Fatal("frozen clone differs from recursive clone")
			}
			frozenCloneWalk(got, func(n *rhino.Node) {
				if n.GetStaticSourceFile() != current {
					t.Fatal("AST/type/TTL source metadata not rebound")
				}
				var last *rhino.Node
				for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
					if c.GetParent() != n {
						t.Fatal("child parent relation corrupted")
					}
					last = c
				}
				if n.GetLastChild() != last {
					t.Fatal("last-child relation corrupted")
				}
				// Every mutable returned node/property/JSDoc type subtree is
				// disturbed before the next clone is requested.
				n.SetLength(900 + iteration)
				n.SetChangeTime(int32(iteration + 1))
				if d := n.GetJSDocInfo(); d != nil {
					if entries := d.GetImplementedInterfaces(); len(entries) > 0 {
						entries[0] = rhino.NewJSTypeExpression(rhino.IR_Name("changed"), "changed")
					}
					d.GetTypeTransformations().Put("changed", rhino.IR_Name("changed"))
				}
			})
			got.DetachChildren()
			if builtinCacheFingerprint(template.Ast) != want {
				t.Fatal("clone mutation reached private template")
			}
		}
	}
}

//go:noinline
func frozenCloneLifetimeTemplate() (*rhino.FrozenParseTree, weak.Pointer[rhino.Node], weak.Pointer[rhino.Node]) {
	rhino.Node_AllocationBoundary()
	file := rhino.NewSimpleSourceFile("externs.zip//frozen-lifetime.js", rhino.StaticSourceFile_SourceKind_EXTERN)
	config := ParserRunner_CreateConfig(Config_LanguageMode_ES_NEXT, nil, Config_StrictMode_SLOPPY)
	original := parserRunnerParseUncached(file, `/** @param {Array<string>} x */ function f(x) {}`, config, rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
	private := cloneBuiltinExternResult(original, file, true)
	return rhino.Node_FreezeParseTree(private.Ast), weak.Make(original.Ast), weak.Make(private.Ast)
}

func TestFrozenExternCloneDoesNotRetainPooledAST(t *testing.T) {
	frozen, original, private := frozenCloneLifetimeTemplate()
	if frozen == nil {
		t.Fatal("freeze failed")
	}
	for i := 0; i < 10 && (original.Value() != nil || private.Value() != nil); i++ {
		runtime.GC()
	}
	if original.Value() != nil || private.Value() != nil {
		t.Fatal("frozen metadata retained a compiler/private AST allocation group")
	}
	runtime.KeepAlive(frozen)
}

func TestFrozenExternCloneCacheOptIn(t *testing.T) {
	old := builtinExternFrozenCloneEnabled
	builtinExternFrozenCloneEnabled = true
	defer func() { builtinExternFrozenCloneEnabled = old }()
	ParserRunner_EnableBuiltinExternCache(true)
	defer ParserRunner_EnableBuiltinExternCache(false)
	file := rhino.NewSimpleSourceFile("externs.zip//frozen-cache.js", rhino.StaticSourceFile_SourceKind_EXTERN)
	config := ParserRunner_CreateConfig(Config_LanguageMode_ES_NEXT, nil, Config_StrictMode_SLOPPY)
	source := `function f(/** number */ x) {} /** @type {string} */ var value;`
	first := ParserRunner_Parse(file, source, config, rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
	want := builtinCacheFingerprint(first.Ast)
	if len(builtinExternCache.groups) != 1 {
		t.Fatal("cache group missing")
	}
	entry := builtinExternCache.groups[0].files[builtinExternCacheKey{file.GetName(), source}]
	if entry == nil || entry.frozen == nil {
		t.Fatal("opt-in did not freeze eligible template")
	}
	first.Ast.DetachChildren()
	second := ParserRunner_Parse(file, source, config, rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
	if builtinCacheFingerprint(second.Ast) != want {
		t.Fatal("cache's frozen path changed result")
	}
}

// Exercise the actual pinned browser/core extern archive, reporting acceptance
// counts so a nominal optimization cannot pass by silently falling back on all
// real inputs. Every accepted file is compared with the original deep clone.
func TestFrozenExternCloneBuiltinCorpusAcceptance(t *testing.T) {
	data := resources.ResourceLoader_GetResourceAsStream("com.google.javascript.jscomp.AbstractCommandLineRunner", "/externs.zip")
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	config := ParserRunner_CreateConfig(Config_LanguageMode_ES_NEXT, nil, Config_StrictMode_SLOPPY)
	files, nodes, fallback := 0, 0, 0
	for _, entry := range archive.File {
		if !strings.HasSuffix(entry.Name, ".js") || (strings.Contains(entry.Name, "/") && !strings.HasPrefix(entry.Name, "browser/")) {
			continue
		}
		reader, err := entry.Open()
		if err != nil {
			t.Fatal(err)
		}
		source, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			t.Fatal(err)
		}
		file := rhino.NewSimpleSourceFile("externs.zip//"+strings.TrimPrefix(entry.Name, "browser/"), rhino.StaticSourceFile_SourceKind_EXTERN)
		reporter := &frozenCloneQuietReporter{}
		parsed := parserRunnerParseUncached(file, string(source), config, reporter)
		if reporter.reported || parsed.Ast == nil || len(parsed.Comments) != 0 || len(parsed.TopLevelStatementRanges) != 0 || parsed.SourceMapURL != nil {
			continue
		}
		private := cloneBuiltinExternResult(parsed, file, true)
		frozen := rhino.Node_FreezeParseTree(private.Ast)
		if frozen == nil {
			fallback++
			continue
		}
		files++
		got := frozen.CloneWithSourceFile(file)
		frozenCloneWalk(got, func(*rhino.Node) { nodes++ })
		if builtinCacheFingerprint(got) != builtinCacheFingerprint(cloneBuiltinExternResult(private, file, false).Ast) {
			t.Fatalf("builtin clone mismatch: %s", entry.Name)
		}
	}
	t.Logf("frozen clone accepted %d builtin files / %d AST+type+TTL nodes; unsupported fallback %d files", files, nodes, fallback)
	if files < 10 || nodes < 1000 {
		t.Fatal("frozen clone failed to cover a substantial real builtin corpus")
	}
}

type frozenCloneQuietReporter struct{ reported bool }

func (r *frozenCloneQuietReporter) Warning(message, source string, line, offset int) {
	r.reported = true
}
func (r *frozenCloneQuietReporter) Error(message, source string, line, offset int) { r.reported = true }
