// Port of test/com/google/javascript/jscomp/deps/ModuleLoaderTest.java @ closure-compiler v20260415

package jscomp

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Java: ModuleLoaderTest.inputs(String...)
func jsaMlInputs(names ...string) []DependencyInfo {
	builder := []DependencyInfo{}
	for _, name := range names {
		builder = append(builder, jsaMlInput(name))
	}
	return builder
}

// jsaMlInput is one input of ModuleLoaderTest.inputs(String...):
// new CompilerInput(SourceFile.fromCode(name, ""), false).
func jsaMlInput(name string) DependencyInfo {
	return NewCompilerInputWithIsExtern(SourceFile_FromCode(name, ""), false)
}

func jsaMlPACKAGE_JSON_MAIN_ENTRIES() *jutil.LinkedHashMap[string, string] {
	return jsaMap(
		"/B/package.json", "/B/lib/b",
		"/node_modules/B/package.json", "/node_modules/B/lib/b.js")
}

func jsaMap(kv ...string) *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	for i := 0; i < len(kv); i += 2 {
		m.Put(kv[i], kv[i+1])
	}
	return m
}

type jsaMapBuilderT struct {
	m *jutil.LinkedHashMap[string, string]
}

func jsaMapBuilder() *jsaMapBuilderT {
	return &jsaMapBuilderT{m: jutil.NewLinkedHashMap[string, string]()}
}

func (b *jsaMapBuilderT) Put(k, v string) *jsaMapBuilderT { b.m.Put(k, v); return b }

func (b *jsaMapBuilderT) Build() *jutil.LinkedHashMap[string, string] { return b.m }

func (b *ModuleLoader_Builder) jsaBuildFunc() func() { return func() { b.Build() } }

// jsaMlCollect is the lambda (CheckLevel level, JSError error) -> errors.add(error).
type jsaMlCollector struct{ errors *[]*JSError }

func (c jsaMlCollector) Report(level *CheckLevel, error *JSError) {
	*c.errors = append(*c.errors, error)
}

func jsaMlCollect(errors *[]*JSError) ErrorHandler { return jsaMlCollector{errors: errors} }

// jsaMlCustomResolver is the anonymous ModuleResolver of testCustomResolution.
type jsaMlCustomResolver struct {
	ModuleResolverBase
}

func (r *jsaMlCustomResolver) ResolveJsModule(scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) *string {
	if strings.HasPrefix(moduleAddress, "@custom/") {
		moduleAddress = moduleAddress[8:]
	}
	return r.ModuleResolverBase.Locate(scriptAddress, moduleAddress)
}

func (r *jsaMlCustomResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) *string {
	return r.ResolveJsModule(scriptAddress, moduleAddress, "", 0, 0)
}

func jsaMlAssertUri(t *testing.T, expected string, actual *ModuleLoader_ModulePath) {
	t.Helper()
	if actual == nil {
		t.Errorf("got nil, want %q", expected)
		return
	}
	jsaEq(t, actual.String(), expected)
}

func jsaMlResolveJsModule(context *ModuleLoader_ModulePath, moduleAddress string) *ModuleLoader_ModulePath {
	return context.ResolveJsModule(moduleAddress, "", -1, -1)
}

func TestModuleLoader_WindowsAddresses(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("js\\a.js", "js\\b.js")).SetFactory(NewNodeModuleResolver_Factory()).Build()
	jsaMlAssertUri(t, "js/a.js", loader.Resolve("js\\a.js"))
	jsaMlAssertUri(t, "js/b.js", jsaMlResolveJsModule(loader.Resolve("js\\a.js"), "./b"))
}

func TestModuleLoader_JsExtensionNode(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("js/a.js", "js/b.js")).SetFactory(NewNodeModuleResolver_FactoryWithEntries(jsaMlPACKAGE_JSON_MAIN_ENTRIES())).SetPathResolver(ModuleLoader_PathResolver_RELATIVE).Build()
	jsaMlAssertUri(t, "js/a.js", loader.Resolve("js/a.js"))
	jsaMlAssertUri(t, "js/b.js", jsaMlResolveJsModule(loader.Resolve("js/a.js"), "./b"))
	jsaMlAssertUri(t, "js/b.js", jsaMlResolveJsModule(loader.Resolve("js/a.js"), "./b.js"))
}

func TestModuleLoader_LocateJsNode(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("A/index.js", "B/index.js", "app.js")).SetFactory(NewNodeModuleResolver_FactoryWithEntries(jsaMlPACKAGE_JSON_MAIN_ENTRIES())).SetPathResolver(ModuleLoader_PathResolver_RELATIVE).Build()
	jsaMlAssertUri(t, "A/index.js", loader.Resolve("A/index.js"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "../A"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "../A/"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "../A/index"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("app.js"), "./A/index"))
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("app.js"), "A/index") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("folder/app.js"), "A/index") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("folder/app.js"), "index") == nil)
}

func TestModuleLoader_LocateNodeModuleNode(t *testing.T) {
	compilerInputs := jsaMlInputs("/A/index.js", "/A/index.json", "/node_modules/A/index.js", "/node_modules/A/foo.js", "/node_modules/A/node_modules/A/index.json", "/B/package.json", "/B/lib/b.js", "/node_modules/B/package.json", "/node_modules/B/lib/b.js")
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{}).SetInputs(compilerInputs).SetFactory(NewNodeModuleResolver_FactoryWithEntries(jsaMlPACKAGE_JSON_MAIN_ENTRIES())).SetPathResolver(ModuleLoader_PathResolver_RELATIVE).Build()
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve(" /foo.js"), "/A"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A/index.js"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "./A"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "./A/index.js"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A/index"))
	jsaMlAssertUri(t, "/A/index.json", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A/index.json"))
	jsaMlAssertUri(t, "/node_modules/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "A"))
	jsaMlAssertUri(t, "/node_modules/A/node_modules/A/index.json", jsaMlResolveJsModule(loader.Resolve("/node_modules/A/foo.js"), "A"))
	jsaMlAssertUri(t, "/node_modules/A/foo.js", jsaMlResolveJsModule(loader.Resolve("/node_modules/A/index.js"), "./foo"))
	jsaMlAssertUri(t, "/B/lib/b.js", jsaMlResolveJsModule(loader.Resolve("/app.js"), "/B"))
	jsaMlAssertUri(t, "/B/lib/b.js", jsaMlResolveJsModule(loader.Resolve("/app.js"), "/B/"))
	jsaMlAssertUri(t, "/node_modules/B/lib/b.js", jsaMlResolveJsModule(loader.Resolve("/app.js"), "B"))
}

func TestModuleLoader_JsExtensionBrowser(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("js/a.js", "js/b.js")).SetFactory(BrowserModuleResolver_FACTORY).Build()
	jsaMlAssertUri(t, "js/a.js", loader.Resolve("js/a.js"))
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("js/a.js"), "./b") == nil)
	jsaMlAssertUri(t, "js/b.js", jsaMlResolveJsModule(loader.Resolve("js/a.js"), "./b.js"))
}

func TestModuleLoader_LocateJsBrowser(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("A/index.js", "B/index.js", "app.js")).SetFactory(BrowserModuleResolver_FACTORY).Build()
	jsaMlAssertUri(t, "A/index.js", loader.Resolve("A/index.js"))
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("B/index.js"), "../A") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("B/index.js"), "../A/index") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("app.js"), "./A/index") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("app.js"), "A/index") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("folder/app.js"), "A/index") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("folder/app.js"), "index") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("B/index.js"), "../A") == nil)
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "../A/index.js"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("app.js"), "./A/index.js"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("folder/app.js"), "../A/index.js"))
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("folder/app.js"), "index") == nil)
}

func TestModuleLoader_LocateNodeModuleBrowser(t *testing.T) {
	compilerInputs := jsaMlInputs("/A/index.js", "/A/index.json", "/node_modules/A/index.js", "/node_modules/A/foo.js", "/node_modules/A/node_modules/A/index.json", "/B/package.json", "/B/lib/b.js", "/node_modules/B/package.json", "/node_modules/B/lib/b.js")
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{}).SetInputs(compilerInputs).SetFactory(BrowserModuleResolver_FACTORY).Build()
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A") == nil)
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A/index.js"))
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("/foo.js"), "./A") == nil)
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "./A/index.js"))
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A/index") == nil)
	jsaMlAssertUri(t, "/A/index.json", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A/index.json"))
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("/foo.js"), "A") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("/node_modules/A/foo.js"), "A") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("/node_modules/A/index.js"), "./foo") == nil)
	jsaMlAssertUri(t, "/node_modules/A/foo.js", jsaMlResolveJsModule(loader.Resolve("/node_modules/A/index.js"), "./foo.js"))
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("/app.js"), "/B") == nil)
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("/app.js"), "B") == nil)
}

func TestModuleLoader_NormalizeUris(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"a", "b", "/c"}).SetInputs(jsaMlInputs()).SetFactory(BrowserModuleResolver_FACTORY).Build()
	jsaMlAssertUri(t, "a.js", loader.Resolve("a/a.js"))
	jsaMlAssertUri(t, "a.js", loader.Resolve("a.js"))
	jsaMlAssertUri(t, "some.js", loader.Resolve("some.js"))
	jsaMlAssertUri(t, "/x.js", loader.Resolve("/x.js"))
	jsaMlAssertUri(t, "x-y.js", loader.Resolve("x:y.js"))
	jsaMlAssertUri(t, "foo%20bar.js", loader.Resolve("foo bar.js"))
}

func TestModuleLoader_DuplicateUris(t *testing.T) {
	e := jsaThrows[jutil.Throwable](t, ModuleLoader_Builder_().SetModuleRoots([]string{"a", "b"}).SetInputs(jsaMlInputs("a/f.js", "b/f.js")).SetFactory(BrowserModuleResolver_FACTORY).jsaBuildFunc())
	if e != nil {
		jsaTrue(t, strings.Contains(e.GetMessage(), "Duplicate module path"), e.GetMessage())
	}
}

func TestModuleLoader_CanonicalizePath(t *testing.T) {
	jsaEq(t, ModuleNames_CanonicalizePath("a/b/c"), "a/b/c")
	jsaEq(t, ModuleNames_CanonicalizePath("a/b/../c"), "a/c")
	jsaEq(t, ModuleNames_CanonicalizePath("a/b/../../b/c"), "b/c")
	jsaEq(t, ModuleNames_CanonicalizePath("a/b/../../c"), "c")
	jsaEq(t, ModuleNames_CanonicalizePath("../a/b/.."), "../a")
	jsaEq(t, ModuleNames_CanonicalizePath("/a/b/../../.."), "/")
	jsaEq(t, ModuleNames_CanonicalizePath("/a/../../../b"), "/b")
	jsaEq(t, ModuleNames_CanonicalizePath("/a/.."), "/")
}

func TestModuleLoader_ToJSIdentifier(t *testing.T) {
	jsaEq(t, ModuleNames_ToJSIdentifier("com/example/test"), "com$example$test")
	jsaEq(t, ModuleNames_ToJSIdentifier("file://a/b.jar!com/example/test"), "file_$$a$b_jar$com$example$test")
}

func TestModuleLoader_EscapePath(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{}).SetInputs(jsaMlInputs("/has:special:chars.js")).SetFactory(BrowserModuleResolver_FACTORY).SetPathResolver(ModuleLoader_PathResolver_RELATIVE).SetPathEscaper(ModuleLoader_PathEscaper_ESCAPE).Build()
	// : is escaped to -
	jsaEq(t, loader.Resolve("file://my/file.js").String(), "file-//my/file.js")
	jsaEq(t, jsaMlResolveJsModule(loader.Resolve("c"), "/has:special:chars.js").String(), "/has-special-chars.js")
}

func TestModuleLoader_DoNoEscapePath(t *testing.T) {
	// : is a character that is escaped
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{}).SetInputs(jsaMlInputs("/has:special:chars.js")).SetFactory(BrowserModuleResolver_FACTORY).SetPathResolver(ModuleLoader_PathResolver_RELATIVE).SetPathEscaper(ModuleLoader_PathEscaper_CANONICALIZE_ONLY).Build()
	jsaEq(t, loader.Resolve("file://my/file.js").String(), "file://my/file.js")
	jsaEq(t, jsaMlResolveJsModule(loader.Resolve("c"), "/has:special:chars.js").String(), "/has:special:chars.js")
}

func TestModuleLoader_LocateNodeModulesNoLeadingSlash(t *testing.T) {
	compilerInputs := jsaMlInputs("/A/index.js", "/A/index.json", "node_modules/A/index.js", "node_modules/A/foo.js", "node_modules/A/node_modules/A/index.json", "/B/package.json", "/B/lib/b.js", "node_modules/B/package.json", "node_modules/B/lib/b.js")
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{}).SetInputs(compilerInputs).SetFactory(NewNodeModuleResolver_FactoryWithEntries(jsaMlPACKAGE_JSON_MAIN_ENTRIES())).SetPathResolver(ModuleLoader_PathResolver_RELATIVE).Build()
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve(" /foo.js"), "/A"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A/index.js"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "./A"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "./A/index.js"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A/index"))
	jsaMlAssertUri(t, "/A/index.json", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "/A/index.json"))
	jsaMlAssertUri(t, "/node_modules/A/index.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "A"))
	jsaMlAssertUri(t, "/node_modules/A/node_modules/A/index.json", jsaMlResolveJsModule(loader.Resolve("node_modules/A/foo.js"), "A"))
	jsaMlAssertUri(t, "node_modules/A/foo.js", jsaMlResolveJsModule(loader.Resolve("node_modules/A/index.js"), "./foo"))
	jsaMlAssertUri(t, "/B/lib/b.js", jsaMlResolveJsModule(loader.Resolve("/app.js"), "/B"))
	jsaMlAssertUri(t, "/node_modules/B/lib/b.js", jsaMlResolveJsModule(loader.Resolve("/app.js"), "B"))
}

func TestModuleLoader_LocateNodeModulesBrowserFieldAdvancedUsage(t *testing.T) {
	// case where the package.json looks like the following:
	//   {"main": "server.js",
	//    "browser": {"server.js": "client.js",
	//                "exclude/this.js": false,
	//                "replace/other.js": "with/alternative.js"}}
	packageJsonMainEntries := jsaMap("/node_modules/mymodule/package.json", "/node_modules/mymodule/server.js", "/node_modules/mymodule/server.js", "/node_modules/mymodule/client.js", "/node_modules/mymodule/override/relative.js", "/node_modules/mymodule/./with/this.js", "/node_modules/mymodule/exclude/this.js", ModuleLoader_JSC_BROWSER_SKIPLISTED_MARKER, "/node_modules/mymodule/replace/other.js", "/node_modules/mymodule/with/alternative.js")
	compilerInputs := jsaMlInputs("/node_modules/mymodule/package.json", "/node_modules/mymodule/client.js", "/node_modules/mymodule/with/alternative.js", "/node_modules/mymodule/with/this.js", "/foo.js")
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{}).SetInputs(compilerInputs).SetFactory(NewNodeModuleResolver_FactoryWithEntries(packageJsonMainEntries)).SetPathResolver(ModuleLoader_PathResolver_RELATIVE).Build()
	jsaMlAssertUri(t, "/node_modules/mymodule/client.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "mymodule"))
	jsaMlAssertUri(t, "/node_modules/mymodule/with/alternative.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "mymodule/replace/other.js"))
	jsaMlAssertUri(t, "/node_modules/mymodule/with/alternative.js", jsaMlResolveJsModule(loader.Resolve("/node_modules/mymodule/client.js"), "./replace/other.js"))
	jsaMlAssertUri(t, "/node_modules/mymodule/with/this.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "mymodule/override/relative.js"))
	jsaMlAssertUri(t, "/node_modules/mymodule/with/this.js", jsaMlResolveJsModule(loader.Resolve("/node_modules/mymodule/client.js"), "./override/relative.js"))
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("/node_modules/mymodule/client.js"), "./exclude/this.js") == nil)
}

func TestModuleLoader_LocateNodeModuleWithMultipleRootsSimple(t *testing.T) {
	// Ensure that root modules work with Node resolver.
	compilerInputs := jsaMlInputs(
		// node_modules in generated_files directory. Note that generated_files
		// is used as root directory so this node_modules should be equivalent to just
		// "/node_modules"
		"/generated_files/node_modules/second.js",

		// file that will be basis for resolving modules in this test.
		"/foo.js")
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"generated_files/"}).SetInputs(compilerInputs).SetFactory(NewNodeModuleResolver_FactoryWithEntries(jsaMap())).SetPathResolver(ModuleLoader_PathResolver_RELATIVE).Build()
	jsaMlAssertUri(t, "/node_modules/second.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "second"))
}

func TestModuleLoader_LocateNodeModuleWithMultipleRoots(t *testing.T) {
	// Ensure that root modules work with Node resolver.
	compilerInputs := jsaMlInputs(
		// node_modules in root directory.
		"/node_modules/first.js",

		// node_modules in generated_files directory. Note that generated_files
		// is used as root directory so this node_modules should be equivalent to just
		// "/node_modules"
		"/generated_files/node_modules/second.js",

		// Here node_modules is not root or insite a root path. So it should be not accessible
		// from foo.js.
		"/some_other/node_modules/third.js",

		// file that will be basis for resolving modules in this test.
		"/foo.js")
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"generated_files/"}).SetInputs(compilerInputs).SetFactory(NewNodeModuleResolver_FactoryWithEntries(jsaMap())).SetPathResolver(ModuleLoader_PathResolver_RELATIVE).Build()
	// 'first' and 'second' should resolve from foo.js
	jsaMlAssertUri(t, "/node_modules/first.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "first"))
	jsaMlAssertUri(t, "/node_modules/second.js", jsaMlResolveJsModule(loader.Resolve("/foo.js"), "second"))
	// 'third' doesn't resolve
	jsaTrue(t, jsaMlResolveJsModule(loader.Resolve("/foo.js"), "third") == nil)
}

func TestModuleLoader_Webpack(t *testing.T) {
	webpackModulesById := jsaMap("1", "A/index.js", "B/index.js", "B/index.js", "3", "app.js")
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("A/index.js", "B/index.js", "app.js")).SetFactory(NewWebpackModuleResolver_Factory(webpackModulesById)).SetPathResolver(ModuleLoader_PathResolver_RELATIVE).Build()
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("A/index.js"), "1"))
	jsaMlAssertUri(t, "/B/index.js", jsaMlResolveJsModule(loader.Resolve("A/index.js"), "B/index.js"))
	jsaMlAssertUri(t, "/app.js", jsaMlResolveJsModule(loader.Resolve("A/index.js"), "3"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "1"))
	jsaMlAssertUri(t, "/B/index.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "B/index.js"))
	jsaMlAssertUri(t, "/app.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "3"))
	jsaMlAssertUri(t, "/A/index.js", jsaMlResolveJsModule(loader.Resolve("app.js"), "1"))
	jsaMlAssertUri(t, "/B/index.js", jsaMlResolveJsModule(loader.Resolve("app.js"), "B/index.js"))
	jsaMlAssertUri(t, "/app.js", jsaMlResolveJsModule(loader.Resolve("app.js"), "3"))
	// Path loading falls back to the node resolver
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "../A"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "../A/"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "../A/index"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("app.js"), "./A/index"))
}

func TestModuleLoader_BrowserWithPrefixReplacement(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("/path/to/project0/index.js", "/path/to/project1/index.js", "app.js")).SetFactory(NewBrowserWithTransformedPrefixesModuleResolver_Factory(jsaMapBuilder().Put("@project0/", "/path/to/project0/").Put("+project1/", "/path/to/project1/").Put("@root/", "/").Build())).Build()
	jsaMlAssertUri(t, "/path/to/project0/index.js", jsaMlResolveJsModule(loader.Resolve("fake.js"), "@project0/index.js"))
	jsaMlAssertUri(t, "/path/to/project1/index.js", jsaMlResolveJsModule(loader.Resolve("fake.js"), "+project1/index.js"))
	jsaMlAssertUri(t, "/app.js", jsaMlResolveJsModule(loader.Resolve("fake.js"), "@root/app.js"))
}

func TestModuleLoader_BrowserWithPrefixReplacementResolveModuleAsPath(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{".", "/path/to/project0/", "/path/to/project1/"}).SetInputs(jsaMlInputs()).SetFactory(NewBrowserWithTransformedPrefixesModuleResolver_Factory(jsaMapBuilder().Put("@project0/", "/path/to/project0/").Put("+project1/", "/path/to/project1/").Put("@root/", "/").Build())).Build()
	jsaMlAssertUri(t, "index.js", loader.Resolve("fake.js").ResolveModuleAsPath("@project0/index.js"))
	jsaMlAssertUri(t, "foo/bar/index.js", loader.Resolve("fake.js").ResolveModuleAsPath("+project1/foo/bar/index.js"))
	jsaMlAssertUri(t, "@not/a/root/index.js", loader.Resolve("fake.js").ResolveModuleAsPath("@not/a/root/index.js"))
}

func TestModuleLoader_BrowserWithPrefixReplacementAppliedMostSpecificToLeast(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("/p0/p1/p2/file.js")).SetFactory(NewBrowserWithTransformedPrefixesModuleResolver_Factory(jsaMapBuilder().Put("0/1/2/", "/p0/p1/p2/").Put("0/", "/p0/").Put("0/1/", "/p0/p1/").Build())).Build()
	jsaMlAssertUri(t, "/p0/p1/p2/file.js", jsaMlResolveJsModule(loader.Resolve("fake.js"), "0/p1/p2/file.js"))
	jsaMlAssertUri(t, "/p0/p1/p2/file.js", jsaMlResolveJsModule(loader.Resolve("fake.js"), "0/1/p2/file.js"))
	jsaMlAssertUri(t, "/p0/p1/p2/file.js", jsaMlResolveJsModule(loader.Resolve("fake.js"), "0/1/2/file.js"))
}

func TestModuleLoader_BrowserWithPrefixReplacementInvalidPrefix(t *testing.T) {
	var errors []*JSError

	loader := ModuleLoader_Builder_().SetErrorHandler(jsaMlCollect(&errors)).SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("/path/to/file.js")).SetFactory(NewBrowserWithTransformedPrefixesModuleResolver_Factory(jsaMap("prefix/", "/path/to/"))).Build()
	jsaMlAssertUri(t, "/path/to/file.js", jsaMlResolveJsModule(loader.Resolve("fake.js"), "prefix/file.js"))
	jsaEq(t, len(errors), 0)
	jsaMlResolveJsModule(loader.Resolve("fake.js"), "invalid/file.js")
	jsaEq(t, len(errors), 1)
	jsaTrue(t, errors[0].Type() == BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH)
}

func TestModuleLoader_CustomResolution(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("A/index.js", "B/index.js", "app.js")).SetFactory(ModuleLoader_ModuleResolverFactoryFunc(func(modulePaths *jutil.LinkedHashSet[string], moduleRootPaths []string, errorHandler ErrorHandler, pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
		r := &jsaMlCustomResolver{}
		r.InitModuleResolverBase(modulePaths, moduleRootPaths, errorHandler, pathEscaper)
		return r
	})).Build()
	jsaMlAssertUri(t, "A/index.js", loader.Resolve("A/index.js"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "../A/index.js"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("app.js"), "./A/index.js"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("folder/app.js"), "../A/index.js"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("B/index.js"), "@custom/A/index.js"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("app.js"), "@custom/A/index.js"))
	jsaMlAssertUri(t, "A/index.js", jsaMlResolveJsModule(loader.Resolve("folder/app.js"), "@custom/A/index.js"))
}

func TestModuleLoader_RootsAppliedMostSpecificFirst(t *testing.T) {
	loader := ModuleLoader_Builder_().SetModuleRoots([]string{"/path/", "/path/to/project/", "/path/to/"}).SetInputs(jsaMlInputs()).SetFactory(BrowserModuleResolver_FACTORY).Build()
	jsaMlAssertUri(t, "file.js", loader.Resolve("/path/to/project/file.js"))
}

func TestModuleLoader_BrowserResolverSilentHandlesExternalUrl(t *testing.T) {
	var errors []*JSError

	loader := ModuleLoader_Builder_().SetErrorHandler(jsaMlCollect(&errors)).SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("app.js", "local/module.js")).SetFactory(BrowserModuleResolver_FACTORY).Build()
	appPath := loader.Resolve("app.js")
	result := appPath.ResolveJsModuleSilently("https://cdn.example.com/library.js")
	jsaTrue(t, result == nil)
	jsaEq(t, len(errors), 0)
}

func TestModuleLoader_BrowserResolverSilentHandlesBareModuleName(t *testing.T) {
	var errors []*JSError

	loader := ModuleLoader_Builder_().SetErrorHandler(jsaMlCollect(&errors)).SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("app.js")).SetFactory(BrowserModuleResolver_FACTORY).Build()
	appPath := loader.Resolve("app.js")
	jsaTrue(t, appPath.ResolveJsModuleSilently("react") == nil)
	jsaTrue(t, appPath.ResolveJsModuleSilently("lodash") == nil)
	jsaTrue(t, appPath.ResolveJsModuleSilently("@babel/core") == nil)
	jsaEq(t, len(errors), 0)
}

func TestModuleLoader_BrowserResolverSilentResolvesRelativePath(t *testing.T) {
	var errors []*JSError

	loader := ModuleLoader_Builder_().SetErrorHandler(jsaMlCollect(&errors)).SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("app.js", "lib/utils.js")).SetFactory(BrowserModuleResolver_FACTORY).Build()
	appPath := loader.Resolve("app.js")
	result := appPath.ResolveJsModuleSilently("./lib/utils.js")
	jsaTrue(t, result != nil)
	jsaEq(t, result.String(), "lib/utils.js")
	jsaEq(t, len(errors), 0)
}

func TestModuleLoader_BrowserResolverSilentHandlesUnresolvableRelativePath(t *testing.T) {
	var errors []*JSError

	loader := ModuleLoader_Builder_().SetErrorHandler(jsaMlCollect(&errors)).SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("app.js")).SetFactory(BrowserModuleResolver_FACTORY).Build()
	appPath := loader.Resolve("app.js")
	result := appPath.ResolveJsModuleSilently("./missing/module.js")
	jsaTrue(t, result == nil)
	jsaEq(t, len(errors), 0)
}

func TestModuleLoader_BrowserResolverDoesNotHandleUnresolvableRelativePath(t *testing.T) {
	var errorsFromRegular []*JSError

	regularLoader := ModuleLoader_Builder_().SetErrorHandler(jsaMlCollect(&errorsFromRegular)).SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("app.js")).SetFactory(BrowserModuleResolver_FACTORY).Build()
	appPath1 := regularLoader.Resolve("app.js")
	appPath1.ResolveJsModule("https://cdn.example.com/lib.js", "test.js", 1, 1)
	jsaEq(t, len(errorsFromRegular), 1)
}

func TestModuleLoader_NodeResolverDoesNotHandleUnresolvableRelativePath(t *testing.T) {
	var errorsFromRegular []*JSError

	regularLoader := ModuleLoader_Builder_().SetErrorHandler(jsaMlCollect(&errorsFromRegular)).SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("app.js")).SetFactory(NewNodeModuleResolver_FactoryWithEntries(jsaMlPACKAGE_JSON_MAIN_ENTRIES())).Build()
	appPath := regularLoader.Resolve("app.js")
	appPath.ResolveJsModule("https://cdn.example.com/lib.js", "test.js", 1, 1)
	jsaEq(t, len(errorsFromRegular), 1)
}

func TestModuleLoader_NodeResolverSilentHandlesAmbiguousIdentifier(t *testing.T) {
	var errors []*JSError

	loader := ModuleLoader_Builder_().SetErrorHandler(jsaMlCollect(&errors)).SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("app.js")).SetFactory(NewNodeModuleResolver_FactoryWithEntries(jsaMlPACKAGE_JSON_MAIN_ENTRIES())).Build()
	appPath := loader.Resolve("app.js")
	result := appPath.ResolveJsModuleSilently("https://cdn.example.com/lib.js")
	jsaTrue(t, result == nil)
	jsaEq(t, len(errors), 0)
}

func TestModuleLoader_NodeResolverSilentResolvesRelativePath(t *testing.T) {
	var errors []*JSError

	loader := ModuleLoader_Builder_().SetErrorHandler(jsaMlCollect(&errors)).SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("app.js", "lib/utils.js")).SetFactory(NewNodeModuleResolver_FactoryWithEntries(jsaMlPACKAGE_JSON_MAIN_ENTRIES())).Build()
	appPath := loader.Resolve("app.js")
	result := appPath.ResolveJsModuleSilently("./lib/utils.js")
	jsaTrue(t, result != nil)
	jsaEq(t, result.String(), "lib/utils.js")
	jsaEq(t, len(errors), 0)
}

func TestModuleLoader_BrowserWithPrefixesSilentHandlesAmbiguousIdentifier(t *testing.T) {
	var errors []*JSError

	loader := ModuleLoader_Builder_().SetErrorHandler(jsaMlCollect(&errors)).SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("/p0/p1/p2/file.js")).SetFactory(NewBrowserWithTransformedPrefixesModuleResolver_Factory(jsaMapBuilder().Put("0/1/2/", "/p0/p1/p2/").Put("0/", "/p0/").Put("0/1/", "/p0/p1/").Build())).Build()
	appPath := loader.Resolve("/p0/file.js")
	result := appPath.ResolveJsModuleSilently("lodash")
	jsaTrue(t, result == nil)
	jsaEq(t, len(errors), 0)
}

func TestModuleLoader_BrowserWithPrefixesSilentResolvesTransformedPrefixes(t *testing.T) {
	var errors []*JSError

	loader := ModuleLoader_Builder_().SetErrorHandler(jsaMlCollect(&errors)).SetModuleRoots([]string{"."}).SetInputs(jsaMlInputs("/p0/p1/p2/file.js")).SetFactory(NewBrowserWithTransformedPrefixesModuleResolver_Factory(jsaMapBuilder().Put("0/1/2/", "/p0/p1/p2/").Put("0/", "/p0/").Put("0/1/", "/p0/p1/").Build())).Build()
	appPath := loader.Resolve("fake.js")
	result := appPath.ResolveJsModuleSilently("0/1/2/file.js")
	jsaTrue(t, result != nil)
	jsaEq(t, result.String(), "/p0/p1/p2/file.js")
	jsaEq(t, len(errors), 0)
}
