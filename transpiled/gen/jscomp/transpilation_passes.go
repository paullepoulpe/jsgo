/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/TranspilationPasses.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_transpilation_passes_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_4 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_5 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_6 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_7 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_8 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_9 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_10 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_11 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_12 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_13 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_14 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_15 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_16 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_17 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_18 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_19 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_20 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_21 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_22 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_23 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_24 j2gort.Ext_java_util_function_Function
var lambdaSingleton_transpilation_passes_25 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_transpilation_passes_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_1 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_1)
		_ = compiler
		var peepholeTranspilations j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
		_ = peepholeTranspilations
		peepholeTranspilations.Add_1(NewReportUntranspilableFeatures(compiler, compiler.GetOptions().GetBrowserFeaturesetYearObject(), compiler.GetOptions().GetOutputFeatureSet()))
		if compiler.GetOptions().NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CATCH_BINDING()) {
			peepholeTranspilations.Add_1(NewRewriteCatchWithNoBinding(compiler))
		}
		if compiler.GetOptions().NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES()) {
			peepholeTranspilations.Add_1(NewEs6NormalizeShorthandProperties(compiler))
		}
		if compiler.GetOptions().NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_NEW_TARGET()) {
			peepholeTranspilations.Add_1(NewRewriteNewDotTarget(compiler))
		}
		return j2gort.PtrToAny(PeepholeTranspilationsPass_Create(compiler, peepholeTranspilations))
	})
	lambdaSingleton_transpilation_passes_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewEs6RewriteModulesToCommonJsModules(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_transpilation_passes_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewEs6RelativizeImportPaths(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_4 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(RewriteAsyncFunctions_Create(j2gort.AsIface[AbstractCompiler](mrp0_)))
	})
	lambdaSingleton_transpilation_passes_5 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(RewriteAsyncIteration_Create(j2gort.AsIface[AbstractCompiler](mrp0_)))
	})
	lambdaSingleton_transpilation_passes_6 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewRewriteObjectSpread(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_7 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewEs7RewriteExponentialOperator(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_8 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewEs6NormalizeClasses(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_9 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_2 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_2)
		_ = compiler
		return NewEs6RewriteClass(compiler, compiler.GetOptions().GetEs6SubclassTranspilation())
	})
	lambdaSingleton_transpilation_passes_10 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewEs6RewriteArrowFunction(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_11 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_3 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_3)
		_ = compiler
		return NewRewritePolyfills_5(compiler, compiler.GetOptions().GetRewritePolyfills(), compiler.GetOptions().GetIsolatePolyfills(), compiler.GetOptions().GetInjectPolyfillsNewerThan(), compiler.GetOptions().ShouldDisambiguateProperties())
	})
	lambdaSingleton_transpilation_passes_12 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_4 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_4)
		_ = compiler
		return NewInstrumentAsyncContext(compiler, compiler.GetOptions().GetOutputFeatureSet().Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS()))
	})
	lambdaSingleton_transpilation_passes_13 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewEs6ConvertSuper(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_14 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewInjectTranspilationRuntimeLibraries_1(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_transpilation_passes_15 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewEs6RewriteRestParameters(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_16 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewEs6RewriteSpreadExpressions(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_17 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewLateEs6ToEs3Converter(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_18 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewEs6ForOfConverter(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_19 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewEs6RewriteBlockScopedFunctionDeclaration(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_transpilation_passes_20 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewEs6RewriteBlockScopedDeclaration(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_transpilation_passes_21 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewEs6RewriteGenerators(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_22 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewRewriteLogicalAssignmentOperatorsPass(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_transpilation_passes_23 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewRewriteOptionalChainingOperator_1(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_transpilation_passes_24 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewRewriteNullishCoalesceOperator(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_transpilation_passes_25 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_9 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_9)
		_ = compiler
		return (NewLambda_CompilerPass(func(externs pkg_rhino.Node, root pkg_rhino.Node) {
			_ = externs
			_ = root
			TranspilationPasses_PostTranspileCheckUnsupportedFeaturesRemoved(compiler)
		}))
	})
}

// Java class: com.google.javascript.jscomp.TranspilationPasses
type TranspilationPasses struct {
	idHash_ atomic.Int32
}

func (this_ *TranspilationPasses) BindSelf_(s any) {
	_ = s
}
func (this_ *TranspilationPasses) IsJavaNil_() bool { return this_ == nil }
func (*TranspilationPasses) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TranspilationPasses", "TranspilationPasses"
}
func (this_ *TranspilationPasses) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TranspilationPasses) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TranspilationPasses", "com.google.javascript.jscomp.TranspilationPasses", func(v_ any) bool {
		return j2gort.InstanceOf[*TranspilationPasses](v_)
	})
}

var static_TranspilationPasses_peepholeTranspilationsPasses PassFactory

func StaticGet_TranspilationPasses_PeepholeTranspilationsPasses() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_peepholeTranspilationsPasses
}
func StaticSet_TranspilationPasses_PeepholeTranspilationsPasses(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_peepholeTranspilationsPasses = v
}

var static_TranspilationPasses_es6RewriteModuleToCjs PassFactory

func StaticGet_TranspilationPasses_Es6RewriteModuleToCjs() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_es6RewriteModuleToCjs
}
func StaticSet_TranspilationPasses_Es6RewriteModuleToCjs(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_es6RewriteModuleToCjs = v
}

var static_TranspilationPasses_es6RelativizeImportPaths PassFactory

func StaticGet_TranspilationPasses_Es6RelativizeImportPaths() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_es6RelativizeImportPaths
}
func StaticSet_TranspilationPasses_Es6RelativizeImportPaths(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_es6RelativizeImportPaths = v
}

var static_TranspilationPasses_rewriteAsyncFunctions PassFactory

func StaticGet_TranspilationPasses_RewriteAsyncFunctions() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_rewriteAsyncFunctions
}
func StaticSet_TranspilationPasses_RewriteAsyncFunctions(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_rewriteAsyncFunctions = v
}

var static_TranspilationPasses_rewriteAsyncIteration PassFactory

func StaticGet_TranspilationPasses_RewriteAsyncIteration() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_rewriteAsyncIteration
}
func StaticSet_TranspilationPasses_RewriteAsyncIteration(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_rewriteAsyncIteration = v
}

var static_TranspilationPasses_rewriteObjectSpread PassFactory

func StaticGet_TranspilationPasses_RewriteObjectSpread() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_rewriteObjectSpread
}
func StaticSet_TranspilationPasses_RewriteObjectSpread(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_rewriteObjectSpread = v
}

var static_TranspilationPasses_rewriteExponentialOperator PassFactory

func StaticGet_TranspilationPasses_RewriteExponentialOperator() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_rewriteExponentialOperator
}
func StaticSet_TranspilationPasses_RewriteExponentialOperator(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_rewriteExponentialOperator = v
}

var static_TranspilationPasses_es6NormalizeClasses PassFactory

func StaticGet_TranspilationPasses_Es6NormalizeClasses() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_es6NormalizeClasses
}
func StaticSet_TranspilationPasses_Es6NormalizeClasses(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_es6NormalizeClasses = v
}

var static_TranspilationPasses_es6RewriteClass PassFactory

func StaticGet_TranspilationPasses_Es6RewriteClass() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_es6RewriteClass
}
func StaticSet_TranspilationPasses_Es6RewriteClass(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_es6RewriteClass = v
}

var static_TranspilationPasses_es6RewriteArrowFunction PassFactory

func StaticGet_TranspilationPasses_Es6RewriteArrowFunction() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_es6RewriteArrowFunction
}
func StaticSet_TranspilationPasses_Es6RewriteArrowFunction(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_es6RewriteArrowFunction = v
}

var static_TranspilationPasses_rewritePolyfills PassFactory

func StaticGet_TranspilationPasses_RewritePolyfills() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_rewritePolyfills
}
func StaticSet_TranspilationPasses_RewritePolyfills(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_rewritePolyfills = v
}

var static_TranspilationPasses_instrumentAsyncContext PassFactory

func StaticGet_TranspilationPasses_InstrumentAsyncContext() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_instrumentAsyncContext
}
func StaticSet_TranspilationPasses_InstrumentAsyncContext(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_instrumentAsyncContext = v
}

var static_TranspilationPasses_es6ConvertSuper PassFactory

func StaticGet_TranspilationPasses_Es6ConvertSuper() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_es6ConvertSuper
}
func StaticSet_TranspilationPasses_Es6ConvertSuper(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_es6ConvertSuper = v
}

var static_TranspilationPasses_injectTranspilationRuntimeLibraries PassFactory

func StaticGet_TranspilationPasses_InjectTranspilationRuntimeLibraries() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_injectTranspilationRuntimeLibraries
}
func StaticSet_TranspilationPasses_InjectTranspilationRuntimeLibraries(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_injectTranspilationRuntimeLibraries = v
}

var static_TranspilationPasses_es6RewriteRestParameters PassFactory

func StaticGet_TranspilationPasses_Es6RewriteRestParameters() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_es6RewriteRestParameters
}
func StaticSet_TranspilationPasses_Es6RewriteRestParameters(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_es6RewriteRestParameters = v
}

var static_TranspilationPasses_es6RewriteSpreadExpressions PassFactory

func StaticGet_TranspilationPasses_Es6RewriteSpreadExpressions() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_es6RewriteSpreadExpressions
}
func StaticSet_TranspilationPasses_Es6RewriteSpreadExpressions(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_es6RewriteSpreadExpressions = v
}

var static_TranspilationPasses_lateConvertEs6ToEs3 PassFactory

func StaticGet_TranspilationPasses_LateConvertEs6ToEs3() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_lateConvertEs6ToEs3
}
func StaticSet_TranspilationPasses_LateConvertEs6ToEs3(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_lateConvertEs6ToEs3 = v
}

var static_TranspilationPasses_es6ForOf PassFactory

func StaticGet_TranspilationPasses_Es6ForOf() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_es6ForOf
}
func StaticSet_TranspilationPasses_Es6ForOf(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_es6ForOf = v
}

var static_TranspilationPasses_rewriteBlockScopedFunctionDeclaration PassFactory

func StaticGet_TranspilationPasses_RewriteBlockScopedFunctionDeclaration() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_rewriteBlockScopedFunctionDeclaration
}
func StaticSet_TranspilationPasses_RewriteBlockScopedFunctionDeclaration(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_rewriteBlockScopedFunctionDeclaration = v
}

var static_TranspilationPasses_rewriteBlockScopedDeclaration PassFactory

func StaticGet_TranspilationPasses_RewriteBlockScopedDeclaration() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_rewriteBlockScopedDeclaration
}
func StaticSet_TranspilationPasses_RewriteBlockScopedDeclaration(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_rewriteBlockScopedDeclaration = v
}

var static_TranspilationPasses_rewriteGenerators PassFactory

func StaticGet_TranspilationPasses_RewriteGenerators() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_rewriteGenerators
}
func StaticSet_TranspilationPasses_RewriteGenerators(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_rewriteGenerators = v
}

var static_TranspilationPasses_rewriteLogicalAssignmentOperatorsPass PassFactory

func StaticGet_TranspilationPasses_RewriteLogicalAssignmentOperatorsPass() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_rewriteLogicalAssignmentOperatorsPass
}
func StaticSet_TranspilationPasses_RewriteLogicalAssignmentOperatorsPass(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_rewriteLogicalAssignmentOperatorsPass = v
}

var static_TranspilationPasses_rewriteOptionalChainingOperator PassFactory

func StaticGet_TranspilationPasses_RewriteOptionalChainingOperator() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_rewriteOptionalChainingOperator
}
func StaticSet_TranspilationPasses_RewriteOptionalChainingOperator(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_rewriteOptionalChainingOperator = v
}

var static_TranspilationPasses_rewriteNullishCoalesceOperator PassFactory

func StaticGet_TranspilationPasses_RewriteNullishCoalesceOperator() PassFactory {
	clinit_TranspilationPasses()
	return static_TranspilationPasses_rewriteNullishCoalesceOperator
}
func StaticSet_TranspilationPasses_RewriteNullishCoalesceOperator(v PassFactory) {
	clinit_TranspilationPasses()
	static_TranspilationPasses_rewriteNullishCoalesceOperator = v
}

var clinitDone_TranspilationPasses uint32
var clinitState_TranspilationPasses j2gort.ClinitState

func clinit_TranspilationPasses() {
	if atomic.LoadUint32(&clinitDone_TranspilationPasses) != 0 {
		return
	}
	clinitSlow_TranspilationPasses()
}

func clinitSlow_TranspilationPasses() {
	if !j2gort.EnterClinit(&clinitState_TranspilationPasses, &clinitDone_TranspilationPasses) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TranspilationPasses, &clinitDone_TranspilationPasses)
	static_TranspilationPasses_peepholeTranspilationsPasses = PassFactory_S_Builder().SetName("peepholeTranspilationsPasses").SetInternalFactory(lambdaSingleton_transpilation_passes_1).Build()
	static_TranspilationPasses_es6RewriteModuleToCjs = PassFactory_S_Builder().SetName("es6RewriteModuleToCjs").SetInternalFactory(lambdaSingleton_transpilation_passes_2).Build()
	static_TranspilationPasses_es6RelativizeImportPaths = PassFactory_S_Builder().SetName("es6RelativizeImportPaths").SetInternalFactory(lambdaSingleton_transpilation_passes_3).Build()
	static_TranspilationPasses_rewriteAsyncFunctions = PassFactory_S_Builder().SetName("rewriteAsyncFunctions").SetInternalFactory(lambdaSingleton_transpilation_passes_4).Build()
	static_TranspilationPasses_rewriteAsyncIteration = PassFactory_S_Builder().SetName("rewriteAsyncIteration").SetInternalFactory(lambdaSingleton_transpilation_passes_5).Build()
	static_TranspilationPasses_rewriteObjectSpread = PassFactory_S_Builder().SetName("rewriteObjectSpread").SetInternalFactory(lambdaSingleton_transpilation_passes_6).Build()
	static_TranspilationPasses_rewriteExponentialOperator = PassFactory_S_Builder().SetName("rewriteExponentialOperator").SetInternalFactory(lambdaSingleton_transpilation_passes_7).Build()
	static_TranspilationPasses_es6NormalizeClasses = PassFactory_S_Builder().SetName("es6NormalizeClasses").SetInternalFactory(lambdaSingleton_transpilation_passes_8).Build()
	static_TranspilationPasses_es6RewriteClass = PassFactory_S_Builder().SetName("Es6RewriteClass").SetInternalFactory(lambdaSingleton_transpilation_passes_9).Build()
	static_TranspilationPasses_es6RewriteArrowFunction = PassFactory_S_Builder().SetName("Es6RewriteArrowFunction").SetInternalFactory(lambdaSingleton_transpilation_passes_10).Build()
	static_TranspilationPasses_rewritePolyfills = PassFactory_S_Builder().SetName("RewritePolyfills").SetInternalFactory(lambdaSingleton_transpilation_passes_11).Build()
	static_TranspilationPasses_instrumentAsyncContext = PassFactory_S_Builder().SetName("instrumentAsyncContext").SetInternalFactory(lambdaSingleton_transpilation_passes_12).Build()
	static_TranspilationPasses_es6ConvertSuper = PassFactory_S_Builder().SetName("es6ConvertSuper").SetInternalFactory(lambdaSingleton_transpilation_passes_13).Build()
	static_TranspilationPasses_injectTranspilationRuntimeLibraries = PassFactory_S_Builder().SetName("es6InjectRuntimeLibraries").SetInternalFactory(lambdaSingleton_transpilation_passes_14).Build()
	static_TranspilationPasses_es6RewriteRestParameters = PassFactory_S_Builder().SetName("es6RewriteRestParameters").SetInternalFactory(lambdaSingleton_transpilation_passes_15).Build()
	static_TranspilationPasses_es6RewriteSpreadExpressions = PassFactory_S_Builder().SetName("es6RewriteSpreadExpressions").SetInternalFactory(lambdaSingleton_transpilation_passes_16).Build()
	static_TranspilationPasses_lateConvertEs6ToEs3 = PassFactory_S_Builder().SetName("lateConvertEs6").SetInternalFactory(lambdaSingleton_transpilation_passes_17).Build()
	static_TranspilationPasses_es6ForOf = PassFactory_S_Builder().SetName("es6ForOf").SetInternalFactory(lambdaSingleton_transpilation_passes_18).Build()
	static_TranspilationPasses_rewriteBlockScopedFunctionDeclaration = PassFactory_S_Builder().SetName("Es6RewriteBlockScopedFunctionDeclaration").SetInternalFactory(lambdaSingleton_transpilation_passes_19).Build()
	static_TranspilationPasses_rewriteBlockScopedDeclaration = PassFactory_S_Builder().SetName("Es6RewriteBlockScopedDeclaration").SetInternalFactory(lambdaSingleton_transpilation_passes_20).Build()
	static_TranspilationPasses_rewriteGenerators = PassFactory_S_Builder().SetName("rewriteGenerators").SetInternalFactory(lambdaSingleton_transpilation_passes_21).Build()
	static_TranspilationPasses_rewriteLogicalAssignmentOperatorsPass = PassFactory_S_Builder().SetName("rewriteLogicalAssignmentOperatorsPass").SetInternalFactory(lambdaSingleton_transpilation_passes_22).Build()
	static_TranspilationPasses_rewriteOptionalChainingOperator = PassFactory_S_Builder().SetName("rewriteOptionalChainingOperator").SetInternalFactory(lambdaSingleton_transpilation_passes_23).Build()
	static_TranspilationPasses_rewriteNullishCoalesceOperator = PassFactory_S_Builder().SetName("rewriteNullishCoalesceOperator").SetInternalFactory(lambdaSingleton_transpilation_passes_24).Build()
}

func (this_ *TranspilationPasses) initFields_TranspilationPasses() {
	_ = this_
}

// Java constructor: TranspilationPasses()
func NewTranspilationPasses() *TranspilationPasses {
	clinit_TranspilationPasses()
	this_ := &TranspilationPasses{}
	InitTranspilationPasses(this_)
	return this_
}

func InitTranspilationPasses(this_ *TranspilationPasses) {
	this_.initFields_TranspilationPasses()
}

// Java static method: TranspilationPasses.addEs6ModulePass(com.google.javascript.jscomp.PassListBuilder passes, com.google.javascript.jscomp.PreprocessorSymbolTable.CachedInstanceFactory preprocessorTableFactory)
func TranspilationPasses_AddEs6ModulePass(passes *PassListBuilder, preprocessorTableFactory *PreprocessorSymbolTable_CachedInstanceFactory) {
	clinit_TranspilationPasses()
	_ = passes
	_ = preprocessorTableFactory
	passes.MaybeAdd(PassFactory_S_Builder().SetName("es6RewriteModule").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_5 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_5)
		_ = compiler
		preprocessorTableFactory.MaybeInitialize(compiler)
		return NewEs6RewriteModules_6(compiler, compiler.GetModuleMetadataMap(), compiler.GetModuleMap(), preprocessorTableFactory.GetInstanceOrNull(), compiler.GetTopScope(), compiler.GetOptions().GetChunkOutputType())
	})).Build())
}

// Java static method: TranspilationPasses.addTranspilationRuntimeLibraries(com.google.javascript.jscomp.PassListBuilder passes)
func TranspilationPasses_AddTranspilationRuntimeLibraries(passes *PassListBuilder) {
	clinit_TranspilationPasses()
	_ = passes
	passes.MaybeAdd(StaticGet_TranspilationPasses_InjectTranspilationRuntimeLibraries())
}

// Java static method: TranspilationPasses.addEs6ModuleToCjsPass(com.google.javascript.jscomp.PassListBuilder passes)
func TranspilationPasses_AddEs6ModuleToCjsPass(passes *PassListBuilder) {
	clinit_TranspilationPasses()
	_ = passes
	passes.MaybeAdd(StaticGet_TranspilationPasses_Es6RewriteModuleToCjs())
}

// Java static method: TranspilationPasses.addEs6RewriteImportPathPass(com.google.javascript.jscomp.PassListBuilder passes)
func TranspilationPasses_AddEs6RewriteImportPathPass(passes *PassListBuilder) {
	clinit_TranspilationPasses()
	_ = passes
	passes.MaybeAdd(StaticGet_TranspilationPasses_Es6RelativizeImportPaths())
}

// Java static method: TranspilationPasses.addTranspilationPasses(com.google.javascript.jscomp.PassListBuilder passes, com.google.javascript.jscomp.CompilerOptions options)
func TranspilationPasses_AddTranspilationPasses(passes *PassListBuilder, options *CompilerOptions) {
	clinit_TranspilationPasses()
	_ = passes
	_ = options
	passes.MaybeAdd(StaticGet_TranspilationPasses_PeepholeTranspilationsPasses())
	passes.MaybeAdd(TranspilationPasses_CreateUnifiedFeatureRemovalPass("featureRemovalPasses", options))
	passes.MaybeAdd(StaticGet_TranspilationPasses_Es6NormalizeClasses())
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_RewriteLogicalAssignmentOperatorsPass())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_RewriteOptionalChainingOperator())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_NULL_COALESCE_OP()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_RewriteNullishCoalesceOperator())
	}
	if options.GetInstrumentAsyncContext() {
		passes.MaybeAdd(StaticGet_TranspilationPasses_InstrumentAsyncContext())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_FOR_AWAIT_OF()) || options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_GENERATORS()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_RewriteAsyncIteration())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD()) || options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_RewriteObjectSpread())
		if !(options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_DESTRUCTURING())) && options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST()) {
			passes.MaybeAdd(TranspilationPasses_GetEs6RewriteDestructuring(StaticGet_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST()))
		}
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_RewriteAsyncFunctions())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_EXPONENT_OP()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_RewriteExponentialOperator())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_CLASSES()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_Es6ConvertSuper())
	}
	if options.NeedsTranspilationFrom(pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7186([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_DESTRUCTURING(), pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_DESTRUCTURING()})) {
		passes.MaybeAdd(TranspilationPasses_GetEs6RewriteDestructuring(StaticGet_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS()))
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_Es6RewriteArrowFunction())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_CLASSES()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_Es6RewriteClass())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_REST_PARAMETERS()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_Es6RewriteRestParameters())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_Es6RewriteSpreadExpressions())
	}
	if options.NeedsTranspilationFrom(pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7186([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), pkg_parser.StaticGet_FeatureSet_Feature_MEMBER_DECLARATIONS(), pkg_parser.StaticGet_FeatureSet_Feature_TEMPLATE_LITERALS()})) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_LateConvertEs6ToEs3())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_FOR_OF()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_Es6ForOf())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_RewriteBlockScopedFunctionDeclaration())
	}
	if options.NeedsTranspilationFrom(pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7186([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS()})) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_RewriteBlockScopedDeclaration())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS()) {
		passes.MaybeAdd(StaticGet_TranspilationPasses_RewriteGenerators())
	}
	passes.MaybeAdd(TranspilationPasses_CreatePostTranspileUnsupportedFeaturesRemovedCheck("postTranspileUnsupportedFeaturesRemovedCheck"))
}

// Java static method: TranspilationPasses.addRewritePolyfillPass(com.google.javascript.jscomp.PassListBuilder passes)
func TranspilationPasses_AddRewritePolyfillPass(passes *PassListBuilder) {
	clinit_TranspilationPasses()
	_ = passes
	passes.MaybeAdd(StaticGet_TranspilationPasses_RewritePolyfills())
}

// Java static method: TranspilationPasses.getEs6RewriteDestructuring(com.google.javascript.jscomp.Es6RewriteDestructuring.ObjectDestructuringRewriteMode rewriteMode)
func TranspilationPasses_GetEs6RewriteDestructuring(rewriteMode Es6RewriteDestructuring_ObjectDestructuringRewriteMode) PassFactory {
	clinit_TranspilationPasses()
	_ = rewriteMode
	return PassFactory_S_Builder().SetName("Es6RewriteDestructuring").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_6 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_6)
		_ = compiler
		return j2gort.PtrToAny(NewEs6RewriteDestructuring_Builder(compiler).SetDestructuringRewriteMode(rewriteMode).Build())
	})).Build()
}

// Java static method: TranspilationPasses.doesScriptHaveAnyOfTheseFeatures(com.google.javascript.rhino.Node script, com.google.javascript.jscomp.parsing.parser.FeatureSet featureSet)
func TranspilationPasses_DoesScriptHaveAnyOfTheseFeatures(script pkg_rhino.Node, featureSet *pkg_parser.FeatureSet) bool {
	clinit_TranspilationPasses()
	_ = script
	_ = featureSet
	var features *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(script)
	_ = features
	return ((features != nil) && features.ContainsAtLeastOneOf(featureSet))
}

// Java static method: TranspilationPasses.processTranspile(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node combinedRoot, com.google.javascript.jscomp.parsing.parser.FeatureSet featuresToRunFor, com.google.javascript.jscomp.NodeTraversal.Callback[] callbacks)
func TranspilationPasses_ProcessTranspile(compiler AbstractCompiler, combinedRoot pkg_rhino.Node, featuresToRunFor *pkg_parser.FeatureSet, callbacks []NodeTraversal_Callback) {
	clinit_TranspilationPasses()
	_ = compiler
	_ = combinedRoot
	_ = featuresToRunFor
	_ = callbacks
	var languageOutFeatures *pkg_parser.FeatureSet = compiler.GetOptions().GetOutputFeatureSet()
	_ = languageOutFeatures
	if languageOutFeatures.Contains_1_FeatureSet(featuresToRunFor) {
		return
	}
	{
		var singleRoot pkg_rhino.Node = combinedRoot.GetFirstChild()
		_ = singleRoot
		for ; singleRoot != 0; singleRoot = singleRoot.GetNext() {
			if !(singleRoot.IsScript()) {
				j2gort.CheckState(false)
			}
			if TranspilationPasses_DoesScriptHaveAnyOfTheseFeatures(singleRoot, featuresToRunFor) {
				for _, el_7 := range j2gort.CheckNotNullSlice(callbacks) {
					callback := el_7
					_ = callback
					NodeTraversal_Traverse(compiler, singleRoot, callback)
				}
			}
		}
	}
}

// Java static method: TranspilationPasses.maybeMarkFeatureAsTranspiledAway(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node root, com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(compiler AbstractCompiler, root pkg_rhino.Node, feature *pkg_parser.FeatureSet_Feature) {
	clinit_TranspilationPasses()
	_ = compiler
	_ = root
	_ = feature
	if !(compiler.HasHaltingErrors()) {
		compiler.MarkFeatureNotAllowed(feature)
		NodeUtil_RemoveFeatureFromAllScripts(root, feature, compiler)
	}
}

// Java static method: TranspilationPasses.maybeMarkFeaturesAsTranspiledAway(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node root, com.google.javascript.jscomp.parsing.parser.FeatureSet transpiledFeatures)
func TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(compiler AbstractCompiler, root pkg_rhino.Node, transpiledFeatures *pkg_parser.FeatureSet) {
	clinit_TranspilationPasses()
	_ = compiler
	_ = root
	_ = transpiledFeatures
	if !(compiler.HasHaltingErrors()) {
		NodeUtil_RemoveFeaturesFromAllScripts(root, transpiledFeatures, compiler)
	}
}

// Java static method: TranspilationPasses.postTranspileCheckUnsupportedFeaturesRemoved(com.google.javascript.jscomp.AbstractCompiler compiler)
func TranspilationPasses_PostTranspileCheckUnsupportedFeaturesRemoved(compiler AbstractCompiler) {
	clinit_TranspilationPasses()
	_ = compiler
	var outputFeatures *pkg_parser.FeatureSet = compiler.GetOptions().GetOutputFeatureSet()
	_ = outputFeatures
	var currentFeatures *pkg_parser.FeatureSet = compiler.GetAllowableFeatures()
	_ = currentFeatures
	if compiler.GetOptions().GetChunkOutputType() == StaticGet_CompilerOptions_ChunkOutputType_ES_MODULES() {
		currentFeatures = currentFeatures.Without_2(pkg_parser.StaticGet_FeatureSet_Feature_MODULES(), []*pkg_parser.FeatureSet_Feature{}).Without_2(pkg_parser.StaticGet_FeatureSet_Feature_IMPORT_META(), []*pkg_parser.FeatureSet_Feature{}).Without_2(pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT(), []*pkg_parser.FeatureSet_Feature{})
	}
	if outputFeatures.GetFeatures().IsEmpty() {
		return
	}
	if !(outputFeatures.Contains_1_FeatureSet(currentFeatures)) {
		var diff *pkg_parser.FeatureSet = currentFeatures.Without_1(outputFeatures)
		_ = diff
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unsupported feature(s) leaked to output code:"+j2gort.StrVal(diff.GetFeatures()))))
	}
}

// Java static method: TranspilationPasses.createFeatureRemovalPass(java.lang.String passName, com.google.javascript.jscomp.parsing.parser.FeatureSet featuresToRemove)
func TranspilationPasses_CreateFeatureRemovalPass(passName string, featuresToRemove *pkg_parser.FeatureSet) PassFactory {
	clinit_TranspilationPasses()
	_ = passName
	_ = featuresToRemove
	return PassFactory_S_Builder().SetName(passName).SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_8 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_8)
		_ = compiler
		return (NewLambda_CompilerPass(func(externs pkg_rhino.Node, root pkg_rhino.Node) {
			_ = externs
			_ = root
			TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(compiler, root, featuresToRemove)
		}))
	})).Build()
}

// Java static method: TranspilationPasses.createPostTranspileUnsupportedFeaturesRemovedCheck(java.lang.String passName)
func TranspilationPasses_CreatePostTranspileUnsupportedFeaturesRemovedCheck(passName string) PassFactory {
	clinit_TranspilationPasses()
	_ = passName
	return PassFactory_S_Builder().SetName(passName).SetInternalFactory(lambdaSingleton_transpilation_passes_25).Build()
}

// Java static method: TranspilationPasses.createUnifiedFeatureRemovalPass(java.lang.String passName, com.google.javascript.jscomp.CompilerOptions options)
func TranspilationPasses_CreateUnifiedFeatureRemovalPass(passName string, options *CompilerOptions) PassFactory {
	clinit_TranspilationPasses()
	_ = passName
	_ = options
	var featuresToMarkRemoved *pkg_parser.FeatureSet = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM()
	_ = featuresToMarkRemoved
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_D()) {
		featuresToMarkRemoved = featuresToMarkRemoved.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_D())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_BIGINT()) {
		featuresToMarkRemoved = featuresToMarkRemoved.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_BIGINT())
	}
	if options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_NUMERIC_SEPARATOR()) {
		featuresToMarkRemoved = featuresToMarkRemoved.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_NUMERIC_SEPARATOR())
	}
	if options.GetChunkOutputType() != StaticGet_CompilerOptions_ChunkOutputType_ES_MODULES() {
		featuresToMarkRemoved = featuresToMarkRemoved.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_MODULES())
		featuresToMarkRemoved = featuresToMarkRemoved.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_IMPORT_META())
		featuresToMarkRemoved = featuresToMarkRemoved.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT())
	}
	if options.NeedsTranspilationFrom(pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7186([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_BINARY_LITERALS(), pkg_parser.StaticGet_FeatureSet_Feature_OCTAL_LITERALS(), pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_U(), pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_Y()})) {
		featuresToMarkRemoved = featuresToMarkRemoved.With_1_FeatureSet_Feature_7186([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_BINARY_LITERALS(), pkg_parser.StaticGet_FeatureSet_Feature_OCTAL_LITERALS(), pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_U(), pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_Y()})
	}
	if options.NeedsTranspilationFrom(pkg_parser.StaticGet_FeatureSet_ES5()) {
		featuresToMarkRemoved = featuresToMarkRemoved.With_1_FeatureSet_Feature_7186([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS(), pkg_parser.StaticGet_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES(), pkg_parser.StaticGet_FeatureSet_Feature_GETTER(), pkg_parser.StaticGet_FeatureSet_Feature_SETTER()})
	}
	featuresToMarkRemoved = featuresToMarkRemoved.With_1_FeatureSet_Feature_7186([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_STRING_CONTINUATION(), pkg_parser.StaticGet_FeatureSet_Feature_TRAILING_COMMA()})
	featuresToMarkRemoved = featuresToMarkRemoved.With_1_FeatureSet_Feature_7186([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_ES_NEXT_RUNTIME(), pkg_parser.StaticGet_FeatureSet_Feature_ES_UNSTABLE_RUNTIME()})
	return TranspilationPasses_CreateFeatureRemovalPass(passName, featuresToMarkRemoved)
}
