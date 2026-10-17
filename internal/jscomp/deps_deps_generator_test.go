// Port of test/com/google/javascript/jscomp/deps/DepsGeneratorTest.java @ closure-compiler v20260415

package jscomp

import (
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// depsGeneratorTest is the JUnit fixture of DepsGeneratorTest.
type depsGeneratorTest struct {
	t            *testing.T
	errorManager ErrorManager
}

// Java: DepsGeneratorTest.setUp()
func newDepsGeneratorTest(t *testing.T) *depsGeneratorTest {
	return &depsGeneratorTest{t: t, errorManager: NewPrintStreamErrorManager(jutil.System_Err())}
}

// depsGeneratorTest_loader is the ModuleLoader every test builds:
// ModuleLoader.builder().setErrorHandler(null).setModuleRoots(roots)
// .setInputs(of()).setFactory(BrowserModuleResolver.FACTORY)
// .setPathResolver(PathResolver.ABSOLUTE).build().
func depsGeneratorTest_loader(roots ...string) *ModuleLoader {
	return ModuleLoader_Builder_().
		SetErrorHandler(nil).
		SetModuleRoots(roots).
		SetInputs([]DependencyInfo{}).
		SetFactory(BrowserModuleResolver_FACTORY).
		SetPathResolver(ModuleLoader_PathResolver_ABSOLUTE).
		Build()
}

func (x *depsGeneratorTest) assertOutput(output *string, expected string) {
	x.t.Helper()
	if output == nil || *output != expected {
		x.t.Errorf("output = %s\nwant %q", depsGeneratorTest_str(output), expected)
	}
}

func (x *depsGeneratorTest) assertNotEmpty(output *string, message string) {
	x.t.Helper()
	if output == nil || *output == "" {
		x.t.Errorf("%s: output = %s", message, depsGeneratorTest_str(output))
	}
}

func depsGeneratorTest_str(s *string) string {
	if s == nil {
		return "null"
	}
	return "\"" + *s + "\""
}

// TODO(johnplaisted): This should eventually be an error. For now people are
// relying on this behavior for interop / ordering. Until we have official
// channels for these allow this behavior, but don't encourage it.
func TestDepsGenerator_Es6ModuleWithGoogProvide(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	srcs := []*SourceFile{
		SourceFile_FromCode(
			"/base/javascript/foo/foo.js",
			"goog.provide('my.namespace');\nimport '../closure/goog/es6.js';"),
		SourceFile_FromCode("/base/javascript/closure/goog/es6.js", "export var es6;"),
	}
	depsGenerator := NewDepsGenerator(
		[]*SourceFile{},
		srcs,
		DepsGenerator_InclusionStrategy_ALWAYS,
		"/base/javascript/closure",
		x.errorManager,
		depsGeneratorTest_loader("/base/"))
	output := depsGenerator.ComputeDependencyCalls()

	x.assertWarnings(
		"File cannot be a combination of goog.provide, goog.module, and/or ES6 " +
			"module: javascript/foo/foo.js")

	// Write the output.
	x.assertNotEmpty(output, "There should be output")

	// Write the expected output.
	expected := "goog.addDependency('../foo/foo.js', ['my.namespace'], ['goog/es6.js'], {'lang': 'es6', 'module': 'es6'});\n" +
		"goog.addDependency('goog/es6.js', [], [], {'lang': 'es6', 'module': 'es6'});\n"

	x.assertOutput(output, expected)
}

func TestDepsGenerator_Es6Modules(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	srcs := []*SourceFile{
		SourceFile_FromCode("/base/javascript/foo/foo.js", "import '../closure/goog/es6';"),
		SourceFile_FromCode("/base/javascript/closure/goog/es6.js", "export var es6;"),
	}
	depsGenerator := NewDepsGenerator(
		[]*SourceFile{},
		srcs,
		DepsGenerator_InclusionStrategy_ALWAYS,
		"/base/javascript/closure",
		x.errorManager,
		depsGeneratorTest_loader("/base/"))
	output := depsGenerator.ComputeDependencyCalls()

	x.assertNoWarnings()

	// Write the output.
	x.assertNotEmpty(output, "There should be output")

	// Write the expected output.
	expected := "goog.addDependency('../foo/foo.js', [], ['goog/es6.js'], {'lang': 'es6', 'module': 'es6'});\n" +
		"goog.addDependency('goog/es6.js', [], [], {'lang': 'es6', 'module': 'es6'});\n"

	x.assertOutput(output, expected)
}

func TestDepsGenerator_Es6ModuleDeclareModuleId(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	srcs := []*SourceFile{
		SourceFile_FromCode(
			"/base/javascript/foo/foo.js", "goog.declareModuleId('my.namespace');\nexport {};"),
		SourceFile_FromCode(
			"/base/javascript/closure/goog/googmodule.js",
			"goog.module('my.goog.module');\n"+
				"const namespace = goog.require('my.namespace');\n"),
	}
	depsGenerator := NewDepsGenerator(
		[]*SourceFile{},
		srcs,
		DepsGenerator_InclusionStrategy_ALWAYS,
		"/base/javascript/closure",
		x.errorManager,
		depsGeneratorTest_loader("/base/"))
	output := depsGenerator.ComputeDependencyCalls()

	x.assertNoWarnings()

	// Write the output.
	x.assertNotEmpty(output, "There should be output")

	// Write the expected output.
	expected := "goog.addDependency('../foo/foo.js', ['my.namespace'], [], {'lang': 'es6', 'module': 'es6'});\n" +
		"goog.addDependency('goog/googmodule.js', ['my.goog.module'], ['my.namespace'], {'lang': 'es6', 'module': 'goog'});\n"

	x.assertOutput(output, expected)
}

// Unit test for an issue run into by
// https://github.com/google/closure-compiler/pull/3026
func TestDepsGenerator_Es6ModuleScanDeps(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	// Simple ES6 modules
	srcs := []*SourceFile{
		SourceFile_FromCode("/src/css-parse.js", "export class StyleNode {}"),
		SourceFile_FromCode(
			"/src/apply-shim-utils.js", "import {StyleNode} from './css-parse.js';"),
	}

	// Run them through a DepsGenerator that is set up the same way as our
	// internal MakeJsDeps tool.
	depsGenerator := NewDepsGenerator(
		[]*SourceFile{},
		srcs,
		DepsGenerator_InclusionStrategy_ALWAYS,
		PathUtil_MakeAbsolute("/base/javascript/closure"),
		x.errorManager,
		depsGeneratorTest_loader("."))

	output := depsGenerator.ComputeDependencyCalls()

	// Make sure that there are no spurious errors.
	x.assertNoWarnings()

	expected := "goog.addDependency('../../../src/css-parse.js', [], [], {'lang': 'es6', 'module': 'es6'});\n" +
		"goog.addDependency('../../../src/apply-shim-utils.js', [], ['../../../src/css-parse.js'], {'lang': 'es6', 'module': 'es6'});\n"
	x.assertOutput(output, expected)
}

// Ensures that deps files are handled correctly both when listed as deps and
// when listed as sources.
func TestDepsGenerator_WithDepsAndSources(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	depsFile1 := SourceFile_FromCode(
		"/base/my-project/deps1.js",
		"goog.addDependency('../prescanned1/file1.js', ['dep.string'], []);\n"+
			"goog.addDependency('../prescanned1/file2.js', [], []);\n"+
			"// Test that this appears only once in the output.\n"+
			"goog.addDependency('../this/is/defined/thrice.js', [], []);\n"+
			"goog.addDependency('../this/is/defined/thrice.js', [], []);\n")
	depsFile2 := SourceFile_FromCode(
		"/base/my-project/deps2.js",
		"goog.addDependency('../prescanned2/file1.js', ['dep.bool', 'dep.number'], ['dep.string']);\n"+
			"goog.addDependency('../prescanned2/file2.js', [], []);\n"+
			"goog.addDependency('../this/is/defined/thrice.js', [], []);\n")
	srcFile1 := SourceFile_FromCode(
		"/base/my-project/src1.js",
		"goog.provide('makejsdeps.file1');\n"+
			"goog.provide('makejsdeps.file1.Test');\n"+
			"\n"+
			"/*\n"+
			"goog.require('failure1)\n"+
			"*/\n"+
			"// goog.require('failure2)\n"+
			"\n"+
			"goog.require('makejsdeps.file2');\n"+
			"// 'goog' should be silently dropped.\n"+
			"goog.require(\"goog\");\n"+
			"goog.require(\"dep.string\");\n"+
			"goog.require(\"dep.number\");\n")
	srcFile2 := SourceFile_FromCode("/base/my-project/src2.js", "goog.provide('makejsdeps.file2');")

	expected := "goog.addDependency('../../my-project/src1.js', ['makejsdeps.file1', 'makejsdeps.file1.Test'], ['makejsdeps.file2', 'dep.string', 'dep.number']);\n" +
		"goog.addDependency('../../my-project/src2.js', ['makejsdeps.file2'], []);\n" +
		"\n" +
		"// Included from: /base/my-project/deps1.js\n" +
		"goog.addDependency('../prescanned1/file1.js', ['dep.string'], []);\n" +
		"goog.addDependency('../prescanned1/file2.js', [], []);\n" +
		"\n" +
		"// Included from: /base/my-project/deps2.js\n" +
		"goog.addDependency('../this/is/defined/thrice.js', [], []);\n" +
		"goog.addDependency('../prescanned2/file1.js', ['dep.bool', 'dep.number'], ['dep.string']);\n" +
		"goog.addDependency('../prescanned2/file2.js', [], []);\n"

	depsGenerator := NewDepsGenerator(
		[]*SourceFile{depsFile1, depsFile2},
		[]*SourceFile{srcFile1, srcFile2},
		DepsGenerator_InclusionStrategy_ALWAYS,
		"/base/javascript/closure",
		x.errorManager,
		depsGeneratorTest_loader("/base/"))

	output := depsGenerator.ComputeDependencyCalls()

	x.assertNoWarnings()
	x.assertOutput(output, expected)

	// Repeat the test with the deps files listed as sources.
	// The only difference should be the addition of addDependency() calls for
	// the deps files.
	depsGenerator = NewDepsGenerator(
		[]*SourceFile{},
		[]*SourceFile{depsFile1, depsFile2, srcFile1, srcFile2},
		DepsGenerator_InclusionStrategy_ALWAYS,
		"/base/javascript/closure",
		x.errorManager,
		depsGeneratorTest_loader("/base/"))

	expectedWithDepsAsSources :=
		"goog.addDependency('../../my-project/deps1.js', [], []);\n" +
			"goog.addDependency('../../my-project/deps2.js', [], []);\n" +
			expected

	output = depsGenerator.ComputeDependencyCalls()

	x.assertNoWarnings()
	x.assertOutput(output, expectedWithDepsAsSources)
}

// Ensures that everything still works when both a deps.js and a
// deps-runfiles.js file are included. Also uses real files.
func TestDepsGenerator_DepsAsSrcs(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	depsFile1 := SourceFile_FromCode(
		"/base/deps1.js",
		"// Test deps file 1.\n"+
			"\n"+
			"goog.addDependency('../prescanned1/file1.js', ['dep.string'], []);\n"+
			"goog.addDependency('../prescanned1/file2.js', [], []);\n"+
			"// Test that this appears only once in the output. It's also defined in deps2.js\n"+
			"goog.addDependency('../this/is/defined/thrice.js', [], []);\n"+
			"goog.addDependency('../this/is/defined/thrice.js', [], []);\n")
	depsFile2 := SourceFile_FromCode(
		"/base/deps2.js",
		"// Test deps file 2.\n"+
			"\n"+
			"goog.addDependency('../prescanned2/file1.js', ['dep.bool', 'dep.number'], ['dep.string']);\n"+
			"goog.addDependency('../prescanned2/file2.js', [], []);\n"+
			"goog.addDependency('../prescanned2/generated.js', ['dep.generated'], []);\n"+
			"goog.addDependency('../this/is/defined/thrice.js', [], []);\n")
	depsGenerator := NewDepsGenerator(
		[]*SourceFile{depsFile1},
		[]*SourceFile{depsFile2},
		DepsGenerator_InclusionStrategy_ALWAYS,
		PathUtil_MakeAbsolute("/base/javascript/closure"),
		x.errorManager,
		depsGeneratorTest_loader("/base//"))

	output := depsGenerator.ComputeDependencyCalls()

	x.assertNotEmpty(output, "There should be output")
	x.assertNoWarnings()
}

func TestDepsGenerator_Externs(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	srcs := []*SourceFile{
		SourceFile_FromCode("/base/javascript/foo/externs.js", "/** @externs */"),
		SourceFile_FromCode("/base/javascript/foo/nonexterns.js", ""),
	}
	depsGenerator := NewDepsGenerator(
		[]*SourceFile{},
		srcs,
		DepsGenerator_InclusionStrategy_ALWAYS,
		"/base/javascript/closure",
		x.errorManager,
		depsGeneratorTest_loader("/base/"))
	output := depsGenerator.ComputeDependencyCalls()

	x.assertNoWarnings()

	x.assertNotEmpty(output, "There should be output")

	expected := "goog.addDependency('../foo/nonexterns.js', [], []);\n"

	x.assertOutput(output, expected)
}

func depsGeneratorTest_assertContains(t *testing.T, result string, want string, contains bool) {
	t.Helper()
	if strings.Contains(result, want) != contains {
		t.Errorf("contains(%q) = %v in %q", want, !contains, result)
	}
}

func TestDepsGenerator_MergeStrategyAlways(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	result := x.testMergeStrategyHelper(DepsGenerator_InclusionStrategy_ALWAYS)
	depsGeneratorTest_assertContains(t, result, "['a']", true)
	depsGeneratorTest_assertContains(t, result, "['b']", true)
	depsGeneratorTest_assertContains(t, result, "['c']", true)
	depsGeneratorTest_assertContains(t, result, "d.js", true)
}

func TestDepsGenerator_MergeStrategyWhenInSrcs(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	result := x.testMergeStrategyHelper(DepsGenerator_InclusionStrategy_WHEN_IN_SRCS)
	depsGeneratorTest_assertContains(t, result, "['a']", false)
	depsGeneratorTest_assertContains(t, result, "['b']", true)
	depsGeneratorTest_assertContains(t, result, "['c']", true)
	depsGeneratorTest_assertContains(t, result, "d.js", false)
}

func TestDepsGenerator_MergeStrategyDoNotDuplicate(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	result := x.testMergeStrategyHelper(DepsGenerator_InclusionStrategy_DO_NOT_DUPLICATE)
	depsGeneratorTest_assertContains(t, result, "['a']", false)
	depsGeneratorTest_assertContains(t, result, "['b']", false)
	depsGeneratorTest_assertContains(t, result, "['c']", true)
	depsGeneratorTest_assertContains(t, result, "d.js", false)
}

// Java: DepsGeneratorTest.testMergeStrategyHelper(InclusionStrategy)
func (x *depsGeneratorTest) testMergeStrategyHelper(mergeStrategy DepsGenerator_InclusionStrategy) string {
	x.t.Helper()
	dep1 := SourceFile_FromCode(
		"dep1.js",
		"goog.addDependency('../../a.js', ['a'], []);\n"+
			"goog.addDependency('../../src1.js', ['b'], []);\n"+
			"goog.addDependency('../../d.js', ['d'], []);\n")
	src1 := SourceFile_FromCode("/base//src1.js", "goog.provide('b');\n")
	src2 := SourceFile_FromCode(
		"/base//src2.js",
		"goog.provide('c');\n"+
			"goog.require('d');\n")
	depsGenerator := NewDepsGenerator(
		[]*SourceFile{dep1},
		[]*SourceFile{src1, src2},
		mergeStrategy,
		PathUtil_MakeAbsolute("/base/javascript/closure"),
		x.errorManager,
		depsGeneratorTest_loader("/base/"))

	output := depsGenerator.ComputeDependencyCalls()

	x.assertNoWarnings()
	x.assertNotEmpty(output, "There should be output files")

	if output == nil {
		return ""
	}
	return *output
}

// Java: DepsGeneratorTest.doErrorMessagesRun(List<SourceFile>, List<SourceFile>, boolean, String...)
func (x *depsGeneratorTest) doErrorMessagesRun(
	deps []*SourceFile, srcs []*SourceFile, fatal bool, errorMessages ...string) {
	x.t.Helper()
	depsGenerator := NewDepsGenerator(
		deps,
		srcs,
		DepsGenerator_InclusionStrategy_ALWAYS,
		"/javascript/closure",
		x.errorManager,
		depsGeneratorTest_loader(""))
	output := depsGenerator.ComputeDependencyCalls()

	if fatal {
		if output != nil {
			x.t.Errorf("No output should have been created: %q", *output)
		}
		x.assertErrors(errorMessages...)
	} else {
		if output == nil {
			x.t.Errorf("Output should have been created.")
		}
		x.assertWarnings(errorMessages...)
	}
}

func TestDepsGenerator_DuplicateProvides(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	dep1 := SourceFile_FromCode("dep1.js", "goog.addDependency('a.js', ['a'], []);\n")
	src1 := SourceFile_FromCode("src1.js", "goog.provide('a');\n")

	x.doErrorMessagesRun(
		[]*SourceFile{dep1},
		[]*SourceFile{src1},
		/* fatal= */ true,
		"Namespace \"a\" is already provided in other file dep1.js")
}

// Ensures that an error is thrown when the closure_path flag is set
// incorrectly.
func TestDepsGenerator_DuplicateProvidesErrorThrownIfBadClosurePathSpecified(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	// Create a stub Closure Library.
	fauxClosureDeps := SourceFile_FromCode("dep1.js", "goog.addDependency('foo/a.js', ['a'], []);\n")
	fauxClosureSrc := SourceFile_FromCode("path/to/closure/foo/a.js", "goog.provide('a');\n")
	// Create a source file that depends on the stub Closure Library.
	userSrc := SourceFile_FromCode(
		"my/package/script.js",
		"goog.require('a');\n"+
			"goog.provide('my.package.script');\n")

	// doErrorMessagesRun uses closure_path //javascript/closure and therefore
	// fails to recognize and de-dupe the stub Closure Library at
	// //path/to/closure.
	x.doErrorMessagesRun(
		[]*SourceFile{fauxClosureDeps},
		[]*SourceFile{fauxClosureSrc, userSrc},
		/* fatal= */ true,
		"Namespace \"a\" is already provided in other file dep1.js")
}

// Ensures that DepsGenerator deduplicates dependencies from custom Closure
// Library branches.
func TestDepsGenerator_DuplicateProvidesIgnoredIfInClosureDirectory(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	// Create a stub Closure Library.
	fauxClosureDeps := SourceFile_FromCode("dep1.js", "goog.addDependency('foo/a.js', ['a'], []);\n")
	fauxClosureSrc := SourceFile_FromCode("path/to/closure/foo/a.js", "goog.provide('a');\n")
	// Create a source file that depends on the stub Closure Library.
	userSrc := SourceFile_FromCode(
		"my/package/script.js",
		"goog.require('a');\n"+
			"goog.provide('my.package.script');\n")
	worker := NewDepsGenerator(
		[]*SourceFile{fauxClosureDeps},
		[]*SourceFile{fauxClosureSrc, userSrc},
		DepsGenerator_InclusionStrategy_ALWAYS,
		PathUtil_MakeAbsolute("./path/to/closure"),
		x.errorManager,
		depsGeneratorTest_loader("."))

	output := worker.ComputeDependencyCalls()

	x.assertNotEmpty(output, "Output should have been created.")
	x.assertNoWarnings()
}

func TestDepsGenerator_DuplicateProvidesSameFile(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	dep1 := SourceFile_FromCode("dep1.js", "goog.addDependency('a.js', ['a'], []);\n")
	src1 := SourceFile_FromCode(
		"src1.js",
		"goog.provide('b');\n"+
			"goog.provide('b');\n")

	x.doErrorMessagesRun(
		[]*SourceFile{dep1},
		[]*SourceFile{src1},
		/* fatal= */ false,
		"Multiple calls to goog.provide(\"b\")")
}

func TestDepsGenerator_SameFileProvideRequire(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	dep1 := SourceFile_FromCode("dep1.js", "goog.addDependency('a.js', ['a'], []);\n")
	src1 := SourceFile_FromCode(
		"src1.js",
		"goog.provide('b');\n"+
			"goog.require('b');\n")

	x.doErrorMessagesRun(
		[]*SourceFile{dep1},
		[]*SourceFile{src1},
		/* fatal= */ false,
		"Namespace \"b\" is both required and provided in the same file.")
}

func TestDepsGenerator_UnknownNamespace(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	dep1 := SourceFile_FromCode("dep1.js", "goog.addDependency('a.js', ['a'], []);\n")
	src1 := SourceFile_FromCode("src1.js", "goog.require('b');\n")

	x.doErrorMessagesRun(
		[]*SourceFile{dep1},
		[]*SourceFile{src1},
		/* fatal= */ true,
		"Namespace \"b\" is required but never provided.\n"+
			"You need to pass a library that has it in srcs or exports to your target's deps.")
}

func TestDepsGenerator_NoDepsInDepsFile(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	dep1 := SourceFile_FromCode("dep1.js", "")

	x.doErrorMessagesRun(
		[]*SourceFile{dep1},
		[]*SourceFile{},
		/* fatal= */ false,
		"No dependencies found in file")
}

func TestDepsGenerator_UnknownEs6Module(t *testing.T) {
	defer SkipIfPendingPort(t)
	x := newDepsGeneratorTest(t)
	src1 := SourceFile_FromCode("src1.js", "import './missing.js';\n")

	x.doErrorMessagesRun(
		[]*SourceFile{},
		[]*SourceFile{src1},
		/* fatal= */ true,
		"Could not find file \"./missing.js\".")
}

func depsGeneratorTest_descriptions(errors []*JSError) []string {
	out := []string{}
	for _, e := range errors {
		out = append(out, e.Description())
	}
	return out
}

// Java: DepsGeneratorTest.assertNoWarnings()
func (x *depsGeneratorTest) assertNoWarnings() {
	x.t.Helper()
	if w, e := x.errorManager.GetWarnings(), x.errorManager.GetErrors(); len(w) != 0 || len(e) != 0 {
		x.t.Errorf("warnings %q, errors %q; want none", depsGeneratorTest_descriptions(w), depsGeneratorTest_descriptions(e))
	}
}

// Java: DepsGeneratorTest.assertWarnings(String...)
func (x *depsGeneratorTest) assertWarnings(messages ...string) {
	x.t.Helper()
	if w := depsGeneratorTest_descriptions(x.errorManager.GetWarnings()); !slices.Equal(w, messages) {
		x.t.Errorf("warnings %q, want %q", w, messages)
	}
	if e := x.errorManager.GetErrors(); len(e) != 0 {
		x.t.Errorf("errors %q, want none", depsGeneratorTest_descriptions(e))
	}
}

// Java: DepsGeneratorTest.assertErrors(String...)
func (x *depsGeneratorTest) assertErrors(messages ...string) {
	x.t.Helper()
	if w := x.errorManager.GetWarnings(); len(w) != 0 {
		x.t.Errorf("warnings %q, want none", depsGeneratorTest_descriptions(w))
	}
	if e := depsGeneratorTest_descriptions(x.errorManager.GetErrors()); !slices.Equal(e, messages) {
		x.t.Errorf("errors %q, want %q", e, messages)
	}
}
