// Port of test/com/google/javascript/jscomp/LazyParsedDependencyInfoTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func lazyParsedDependencyInfoTest_flags(kv ...string) *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	for i := 0; i < len(kv); i += 2 {
		m.Put(kv[i], kv[i+1])
	}
	return m
}

// lazyParsedDependencyInfoTest_assertMapContainsExactly is
// assertThat(map).containsExactly(k1, v1, ...) (in any order).
func lazyParsedDependencyInfoTest_assertMapContainsExactly(t *testing.T, got *jutil.LinkedHashMap[string, string], kv ...string) {
	t.Helper()
	if !got.Equals(lazyParsedDependencyInfoTest_flags(kv...)) {
		t.Errorf("load flags: got %v, want %v", got, kv)
	}
}

func TestLazyParsedDependencyInfo_Delegation(t *testing.T) {
	baz := DependencyInfo_Require_GoogRequireSymbol("baz")
	qux := DependencyInfo_Require_GoogRequireSymbol("qux")
	compiler := NewCompiler() // new Compiler()
	ast := NewCompilerInput(SourceFile_FromCode("file.js", "// nothing here"))
	delegate := SimpleDependencyInfo_Builder_("path/to/1.js", "path/2.js").
		SetProvidesArray("foo", "bar").
		SetRequiresArray(baz, qux).
		Build()
	var info DependencyInfo = NewLazyParsedDependencyInfo(delegate, ast, compiler)

	assertEqual(t, "name", info.GetName(), "path/2.js")
	assertEqual(t, "path", info.GetPathRelativeToClosureBase(), "path/to/1.js")
	assertSliceEqual(t, "provides", info.GetProvides(), []string{"foo", "bar"})
	if got := info.GetRequires(); len(got) != 2 || !got[0].Equals(baz) || !got[1].Equals(qux) {
		t.Errorf("requires: got %v", got)
	}
}

func TestLazyParsedDependencyInfo_LoadFlagsParsesEs3(t *testing.T) {
	compiler := newTestCompiler(NewCompilerOptions())
	ast := NewCompilerInput(SourceFile_FromCode("file.js", "// nothing here"))
	delegate := SimpleDependencyInfo_Builder_("", "").SetLoadFlags(lazyParsedDependencyInfoTest_flags("foo", "bar")).Build()
	var info DependencyInfo = NewLazyParsedDependencyInfo(delegate, ast, compiler)

	lazyParsedDependencyInfoTest_assertMapContainsExactly(t, info.GetLoadFlags(), "foo", "bar")
	assertEqual(t, "isGoogModule", info.IsGoogModule(), false)
}

func TestLazyParsedDependencyInfo_LoadFlagsParsesEs5(t *testing.T) {
	compiler := newTestCompiler(NewCompilerOptions())
	ast := NewCompilerInput(SourceFile_FromCode("file.js", "var x = [1, 2,];"))
	delegate := SimpleDependencyInfo_Builder_("", "").
		SetLoadFlags(lazyParsedDependencyInfoTest_flags("module", "goog")).
		Build()
	var info DependencyInfo = NewLazyParsedDependencyInfo(delegate, ast, compiler)

	lazyParsedDependencyInfoTest_assertMapContainsExactly(t, info.GetLoadFlags(), "module", "goog", "lang", "es5")
	assertEqual(t, "isGoogModule", info.IsGoogModule(), true)
}

func TestLazyParsedDependencyInfo_LoadFlagsParsesEs6Impl(t *testing.T) {
	compiler := newTestCompiler(NewCompilerOptions())
	ast := NewCompilerInput(SourceFile_FromCode("file.js", "class X {}"))
	delegate := SimpleDependencyInfo_Builder_("", "").SetLoadFlags(lazyParsedDependencyInfoTest_flags("foo", "bar")).Build()
	var info DependencyInfo = NewLazyParsedDependencyInfo(delegate, ast, compiler)

	lazyParsedDependencyInfoTest_assertMapContainsExactly(t, info.GetLoadFlags(), "foo", "bar", "lang", "es6")
	assertEqual(t, "isGoogModule", info.IsGoogModule(), false)
}

func TestLazyParsedDependencyInfo_LoadFlagsParsesEs6(t *testing.T) {
	compiler := newTestCompiler(NewCompilerOptions())
	ast := NewCompilerInput(SourceFile_FromCode("file.js", "let [a, b] = [1, 2];"))
	delegate := SimpleDependencyInfo_Builder_("", "").SetLoadFlags(lazyParsedDependencyInfoTest_flags("foo", "bar")).Build()
	var info DependencyInfo = NewLazyParsedDependencyInfo(delegate, ast, compiler)

	lazyParsedDependencyInfoTest_assertMapContainsExactly(t, info.GetLoadFlags(), "foo", "bar", "lang", "es6")
	assertEqual(t, "isGoogModule", info.IsGoogModule(), false)
}

func TestLazyParsedDependencyInfo_ParseIsLazy(t *testing.T) {
	compiler := newTestCompiler(NewCompilerOptions())
	ast := NewCompilerInput(SourceFile_FromCode("file.js", "parse error"))
	delegate := SimpleDependencyInfo_Builder_("", "").Build()
	var info DependencyInfo = NewLazyParsedDependencyInfo(delegate, ast, compiler)

	info.GetName()
	info.GetPathRelativeToClosureBase()
	info.GetProvides()
	info.GetRequires()

	assertEqual(t, "errors before", compiler.GetErrorManager().GetErrorCount(), 0)
	info.GetLoadFlags()
	if compiler.GetErrorManager().GetErrorCount() < 1 {
		t.Errorf("errors after: got %d, want at least 1", compiler.GetErrorManager().GetErrorCount())
	}
}

func TestLazyParsedDependencyInfo_ModuleConflict(t *testing.T) {
	options := NewCompilerOptions()
	compiler := newTestCompiler(options)
	ast := NewCompilerInput(SourceFile_FromCode("file.js", "export let foo = 42;"))
	delegate := SimpleDependencyInfo_Builder_("", "my/js.js").
		SetLoadFlags(lazyParsedDependencyInfoTest_flags("module", "goog")).
		Build()
	var info DependencyInfo = NewLazyParsedDependencyInfo(delegate, ast, compiler)

	assertEqual(t, "warnings before", len(compiler.GetErrorManager().GetWarnings()), 0)
	lazyParsedDependencyInfoTest_assertMapContainsExactly(t, info.GetLoadFlags(), "module", "es6", "lang", "es6")
	assertEqual(t, "isEs6Module", info.IsEs6Module(), false)
	name := "my/js.js"
	want := JSError_MakeWithSourceLocation(&name, -1, -1, ModuleLoader_MODULE_CONFLICT, "my/js.js")
	got := compiler.GetErrorManager().GetWarnings()
	if len(got) != 1 || !jutil.ObjectsEquals(got[0], want) {
		t.Errorf("warnings: got %v, want [%v]", got, want)
	}
}
