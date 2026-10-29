// Port of com/google/javascript/jscomp/InjectTranspilationRuntimeLibraries.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InjectTranspilationRuntimeLibraries injects JS library code that may be
// needed by the transpiled form of the input source code.
//
// The intention here is to add anything that could be needed and rely on
// RemoveUnusedCode to remove the parts that don't end up getting used. This
// pass should run before type checking so the type checking code can add
// type information to the injected JavaScript for checking and optimization
// purposes.
//
// TODO(b/120486392): consider merging this pass with InjectRuntimeLibraries
// and RewritePolyfills.
type InjectTranspilationRuntimeLibraries struct {
	compiler AbstractCompiler

	runtimeLibs                   *RuntimeJsLibManager
	injectedClassExtendsLibraries bool
	shouldInstrumentAsyncContext  bool
}

// Java: InjectTranspilationRuntimeLibraries(AbstractCompiler)
func NewInjectTranspilationRuntimeLibraries(compiler AbstractCompiler) *InjectTranspilationRuntimeLibraries {
	return NewInjectTranspilationRuntimeLibrariesWithInstrumentAsyncContext(
		compiler, compiler.GetOptions().GetInstrumentAsyncContext())
}

// Java: InjectTranspilationRuntimeLibraries(AbstractCompiler, boolean)
func NewInjectTranspilationRuntimeLibrariesWithInstrumentAsyncContext(
	compiler AbstractCompiler, shouldInstrumentAsyncContext bool) *InjectTranspilationRuntimeLibraries {
	i := &InjectTranspilationRuntimeLibraries{}
	i.compiler = compiler
	i.runtimeLibs = compiler.GetRuntimeJsLibManager()
	i.injectedClassExtendsLibraries = false
	i.shouldInstrumentAsyncContext = shouldInstrumentAsyncContext
	return i
}

// Java: InjectTranspilationRuntimeLibraries.process(Node, Node)
func (i *InjectTranspilationRuntimeLibraries) Process(externs rhino.Node, root rhino.Node) {
	used := parser.FeatureSet_ES3
	for script := root.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		used = used.WithFeatureSet(InjectTranspilationRuntimeLibraries_GetScriptFeatures(script))
	}

	outputFeatures := i.compiler.GetOptions().GetOutputFeatureSet()

	// Check for references to class `extends` clauses
	if !outputFeatures.Contains(used) {
		NodeUtil_VisitPostOrderWithPred(
			root, i.CheckForClassExtends, func(unused rhino.Node) bool { return !i.injectedClassExtendsLibraries })
	}

	mustBeCompiledAway := used.Without(outputFeatures)

	// We will need these runtime methods when we transpile, but we want the runtime
	// functions to be have JSType applied to it by the type inferrence.

	if mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_TEMPLATE_LITERALS) {
		i.runtimeLibs.InjectLibForField("$jscomp.createTemplateTagFirstArg")
		i.runtimeLibs.InjectLibForField("$jscomp.createTemplateTagFirstArgWithRaw")
	}

	if mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_FOR_OF) ||
		mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING) ||
		mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_OBJECT_PATTERN_REST) {
		// `makeIterator` isn't needed directly for `OBJECT_PATTERN_REST`, but when we transpile
		// a destructuring case that contains it, we transpile the entire destructured assignment,
		// which may also include `ARRAY_DESTRUCTURING`.
		i.runtimeLibs.InjectLibForField("$jscomp.makeIterator")
	}

	if mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_FOR_OF) {
		i.runtimeLibs.InjectLibForField("$jscomp.iteratorClose")
	}

	if mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_ARRAY_PATTERN_REST) {
		i.runtimeLibs.InjectLibForField("$jscomp.arrayFromIterator")
	}

	if mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_SPREAD_EXPRESSIONS) {
		// We must automatically generate the default constructor for descendent classes,
		// and those must call super(...arguments), so we end up injecting our own spread
		// expressions for such cases.
		i.runtimeLibs.InjectLibForField("$jscomp.arrayFromIterable")
	}

	if (mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD) ||
		mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_OBJECT_PATTERN_REST)) &&
		!outputFeatures.Contains(parser.FeatureSet_ES2015) {
		// We need `Object.assign` to transpile `obj = {a, ...rest};` or `const {a, ...rest} = obj;`,
		// but the output language level doesn't indicate that it is guaranteed to be present, so
		// we'll include our polyfill.
		// Use "ensureLibraryInjected" instead of "injectLibForField" because we're injecting the
		// polyfill
		// for Object.assign, which is thus not an actual $jscomp.* method we can lookup.
		i.runtimeLibs.EnsureLibraryInjected("es6/object/assign", false /* force= */)
	}

	if mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER) {
		i.runtimeLibs.InjectLibForField("$jscomp.global")
	}

	if mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_GENERATORS) {
		i.runtimeLibs.InjectLibForField("$jscomp.generator.createGenerator")
		i.runtimeLibs.InjectLibForField("$jscomp.generator.Context")
	}

	if mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS) {
		i.runtimeLibs.InjectLibForField("$jscomp.asyncExecutePromiseGeneratorFunction")
		if !outputFeatures.ContainsFeature(parser.FeatureSet_Feature_GENERATORS) {
			i.runtimeLibs.InjectLibForField("$jscomp.asyncExecutePromiseGeneratorProgram")
			i.runtimeLibs.InjectLibForField("$jscomp.generator.Context")
		}
	}

	if mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_ASYNC_GENERATORS) {
		i.runtimeLibs.InjectLibForField("$jscomp.asyncExecutePromiseGeneratorFunction")
		i.runtimeLibs.InjectLibForField("$jscomp.AsyncGeneratorWrapper")
		i.runtimeLibs.InjectLibForField("$jscomp.AsyncGeneratorWrapper$ActionRecord")
		i.runtimeLibs.InjectLibForField("$jscomp.AsyncGeneratorWrapper$ActionEnum.AWAIT_VALUE")
		i.runtimeLibs.InjectLibForField("$jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE")
		i.runtimeLibs.InjectLibForField("$jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_STAR")
		if !outputFeatures.ContainsFeature(parser.FeatureSet_Feature_GENERATORS) {
			i.runtimeLibs.InjectLibForField("$jscomp.asyncExecutePromiseGeneratorProgram")
			i.runtimeLibs.InjectLibForField("$jscomp.generator.Context")
		}
	}

	if mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_FOR_AWAIT_OF) {
		i.runtimeLibs.InjectLibForField("$jscomp.makeAsyncIterator")
	}

	if mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_REST_PARAMETERS) {
		i.runtimeLibs.InjectLibForField("$jscomp.getRestArguments")
	}

	if used.ContainsFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS) &&
		mustBeCompiledAway.ContainsFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS) {
		i.runtimeLibs.InjectLibForField("$jscomp.PrivateMap")
	}

	if i.shouldInstrumentAsyncContext &&
		// NOTE: async functions only matter for output features, since we don't bother
		// instrumenting them if they're being transpiled away.  Generators are relevant
		// regardless of whether they're transpiled or not.
		(outputFeatures.ContainsFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS) ||
			used.ContainsFeature(parser.FeatureSet_Feature_GENERATORS) ||
			used.ContainsFeature(parser.FeatureSet_Feature_ASYNC_GENERATORS)) {
		i.runtimeLibs.InjectLibForField("$jscomp.asyncContextStart")
	}
}

// Java: InjectTranspilationRuntimeLibraries.getScriptFeatures(Node)
func InjectTranspilationRuntimeLibraries_GetScriptFeatures(script rhino.Node) *parser.FeatureSet {
	features := NodeUtil_GetFeatureSetOfScript(script)
	if features != nil {
		return features
	}
	return parser.FeatureSet_ES3
}

// Java: InjectTranspilationRuntimeLibraries.checkForClassExtends(Node)
func (i *InjectTranspilationRuntimeLibraries) CheckForClassExtends(n rhino.Node) {
	if !n.IsClass() {
		return
	}
	// This is technically an optimization - we could just always inject these when we see
	// Feature.CLASSES. That's fine for real code, but just makes some unit testing
	// harder because more runtime libraries are injected.
	superclass := n.GetSecondChild()
	if !i.injectedClassExtendsLibraries && !superclass.IsEmpty() {
		i.runtimeLibs.InjectLibForField("$jscomp.construct")
		i.runtimeLibs.InjectLibForField("$jscomp.inherits")
		// We must automatically generate the default constructor for descendent classes,
		// and those must call super(...arguments), so we end up injecting our own spread
		// expressions for such cases.
		i.runtimeLibs.InjectLibForField("$jscomp.arrayFromIterable")
		i.injectedClassExtendsLibraries = true
	}
}
