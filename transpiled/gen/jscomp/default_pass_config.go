/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/DefaultPassConfig.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_default_pass_config_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_4 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_5 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_6 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_7 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_8 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_9 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_10 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_11 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_12 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_13 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_14 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_15 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_16 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_17 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_18 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_19 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_20 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_21 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_22 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_23 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_24 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_25 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_26 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_27 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_28 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_29 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_30 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_31 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_32 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_33 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_34 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_35 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_36 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_37 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_38 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_39 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_40 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_41 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_42 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_43 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_44 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_45 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_46 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_47 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_48 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_49 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_50 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_51 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_52 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_53 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_54 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_55 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_56 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_57 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_58 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_59 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_60 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_61 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_62 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_63 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_64 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_65 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_66 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_67 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_68 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_69 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_70 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_71 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_72 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_73 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_74 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_75 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_76 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_77 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_78 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_79 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_80 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_81 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_82 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_83 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_84 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_85 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_86 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_87 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_88 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_89 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_90 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_91 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_92 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_93 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_94 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_95 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_96 j2gort.Ext_java_util_function_Function
var lambdaSingleton_default_pass_config_97 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_default_pass_config_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewForbidDynamicImportUsage(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_1 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_1)
		_ = compiler
		return NewRewriteDynamicImports(compiler, compiler.GetOptions().GetDynamicImportAlias(), compiler.GetOptions().GetChunkOutputType())
	})
	lambdaSingleton_default_pass_config_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewCheckSideEffects_StripProtection(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_4 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_5 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_5)
		_ = compiler
		return NewCheckMissingRequires(compiler, compiler.GetModuleMetadataMap())
	})
	lambdaSingleton_default_pass_config_5 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_6 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_6)
		_ = compiler
		return DefaultPassConfig_CombineChecks(compiler, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(NewCheckJSDoc(compiler), NewEs6CheckModule(compiler))))
	})
	lambdaSingleton_default_pass_config_6 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewConvertToTypedInterface(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_7 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewExtraRequireRemover(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_8 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_9 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_9)
		_ = compiler
		var convention CodingConvention = compiler.GetCodingConvention()
		_ = convention
		if convention.GetExportSymbolFunction() != j2gort.NullString {
			return NewExportTestFunctions(compiler, convention.GetExportSymbolFunction(), convention.GetExportPropertyFunction())
		} else {
			return NewErrorPass_2_AbstractCompiler_DiagnosticType(compiler, StaticGet_DefaultPassConfig_GENERATE_EXPORTS_ERROR())
		}
	})
	lambdaSingleton_default_pass_config_9 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewAngularPass(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_default_pass_config_10 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_15 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_15)
		_ = compiler
		return j2gort.PtrToAny(ScopedAliases_S_Builder(compiler).Build())
	})
	lambdaSingleton_default_pass_config_11 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_17 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_17)
		_ = compiler
		return NewInjectRuntimeLibraries(compiler, j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](compiler.GetOptions().GetForceLibraryInjectionList())))
	})
	lambdaSingleton_default_pass_config_12 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewRemoveWeakSources(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_13 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewDeclaredGlobalExternsOnWindow(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_14 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewCheckTypeImportCodeReferences(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_15 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_18 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_18)
		_ = compiler
		return NewClosureCheckModule(compiler, compiler.GetModuleMetadataMap())
	})
	lambdaSingleton_default_pass_config_16 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_20 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_20)
		_ = compiler
		return NewCheckClosureImports(compiler, compiler.GetModuleMetadataMap())
	})
	lambdaSingleton_default_pass_config_17 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_21 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_21)
		_ = compiler
		return NewRewriteGoogJsImports(compiler, StaticGet_RewriteGoogJsImports_Mode_LINT_AND_REWRITE(), compiler.GetModuleMap())
	})
	lambdaSingleton_default_pass_config_18 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_24 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_24)
		_ = compiler
		var useTypesForOptimization bool = compiler.GetOptions().ShouldUseTypesForLocalOptimization()
		_ = useTypesForOptimization
		var peepholeOptimizations j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
		_ = peepholeOptimizations
		peepholeOptimizations.Add_1(NewPeepholeRemoveDeadCode())
		if compiler.GetOptions().GetJ2clPass().ShouldAddJ2clPasses() {
			peepholeOptimizations.Add_1(NewJ2clEqualitySameRewriterPass(useTypesForOptimization))
		}
		return NewPeepholeOptimizationsPass_3_AbstractCompiler_string_Ext_java_util_List(compiler, "earlyPeepholeOptimizations", peepholeOptimizations)
	})
	lambdaSingleton_default_pass_config_19 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_26 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_26)
		_ = compiler
		return DefaultPassConfig_CreatePeepholeOptimizationsPass(compiler, "peepholeOptimizations", true)
	})
	lambdaSingleton_default_pass_config_20 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_27 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_27)
		_ = compiler
		return DefaultPassConfig_CreatePeepholeOptimizationsPass(compiler, "peepholeOptimizations", true)
	})
	lambdaSingleton_default_pass_config_21 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_28 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_28)
		_ = compiler
		return DefaultPassConfig_CreatePeepholeOptimizationsPass(compiler, "peepholeOptimizations", false)
	})
	lambdaSingleton_default_pass_config_22 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewVarCheck_1(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_default_pass_config_23 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewInferConsts(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_default_pass_config_24 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewVariableReferenceCheck(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_25 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewCheckSuper(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_default_pass_config_26 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_32 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_32)
		_ = compiler
		return NewLambda_CompilerPass(func(externs pkg_rhino.Node, root pkg_rhino.Node) {
			_ = externs
			_ = root
			compiler.ClearTypedScopeCreator()
		})
	})
	lambdaSingleton_default_pass_config_27 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_33 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_33)
		_ = compiler
		return NewLambda_CompilerPass(func(externs pkg_rhino.Node, root pkg_rhino.Node) {
			_ = externs
			_ = root
			compiler.SetTopScope(nil)
			{
				cur_34 := j2gort.NewIterCursor(j2gort.CheckNotNull(compiler.GetInputsInOrder()))
				for cur_34.HasNext() {
					rawElem_35 := cur_34.Next()
					compilerInput := j2gort.CastPtr[CompilerInput](rawElem_35)
					_ = compilerInput
					compilerInput.SetTypedScope(nil)
				}
			}
		})
	})
	lambdaSingleton_default_pass_config_28 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_36 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_36)
		_ = compiler
		return (NewLambda_CompilerPass(func(unused pkg_rhino.Node, srcRoot pkg_rhino.Node) {
			_ = unused
			_ = srcRoot
			var globalRoot pkg_rhino.Node = srcRoot.GetParent()
			_ = globalRoot
			var val_37 TypeInferencePass
			InitTypeInferencePass(&val_37, compiler, compiler.GetReverseAbstractInterpreter(), j2gort.CastPtr[TypedScopeCreator](compiler.GetTypedScopeCreator()))
			var inferencePass *TypeInferencePass = &val_37
			_ = inferencePass
			compiler.SetTypeCheckingHasRun(true)
			compiler.SetTopScope(inferencePass.InferAllScopes(globalRoot))
		}))
	})
	lambdaSingleton_default_pass_config_29 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewInferJSDocInfo(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_30 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewCheckAccessControls(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_31 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_40 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_40)
		_ = compiler
		var callbacks *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}).AddSelf_(NewCheckConstPrivateProperties(compiler)).AddSelf_(NewCheckConstantCaseNames(compiler)).AddSelf_(NewCheckEmptyStatements(compiler)).AddSelf_(NewCheckEnums(compiler)).AddSelf_(NewCheckEs6ModuleFileStructure(compiler)).AddSelf_(NewCheckEs6Modules(compiler)).AddSelf_(NewCheckNoMutatedEs6Exports(compiler)).AddSelf_(NewCheckGoogModuleTypeScriptName(compiler)).AddSelf_(NewCheckInterfaces(compiler)).AddSelf_(NewCheckJSDocStyle(compiler)).AddSelf_(NewCheckMissingSemicolon(compiler)).AddSelf_(NewCheckNullabilityModifiers(compiler)).AddSelf_(NewCheckPrimitiveAsObject(compiler)).AddSelf_(NewCheckPrototypeProperties(compiler)).AddSelf_(NewCheckUnusedPrivateProperties(compiler)).AddSelf_(NewCheckUnusedLabels(compiler)).AddSelf_(NewCheckUselessBlocks(compiler)).AddSelf_(NewCheckVar(compiler))
		_ = callbacks
		return DefaultPassConfig_CombineChecks(compiler, j2gort.AsIface[j2gort.Ext_java_util_List](callbacks.BuildImmList_()))
	})
	lambdaSingleton_default_pass_config_32 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_41 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_41)
		_ = compiler
		var callbacks j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.Ext_com_google_common_collect_ImmutableList_Of_3(NewCheckArrayWithGoogObject(compiler), NewImplicitNullabilityCheck(compiler), NewCheckNestedNames(compiler))
		_ = callbacks
		return DefaultPassConfig_CombineChecks(compiler, j2gort.AsIface[j2gort.Ext_java_util_List](callbacks))
	})
	lambdaSingleton_default_pass_config_33 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_42 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_42)
		_ = compiler
		return DefaultPassConfig_CombineChecks(compiler, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(NewCheckProvidesSorted(StaticGet_CheckProvidesSorted_Mode_COLLECT_AND_REPORT()), NewCheckRequiresSorted(StaticGet_CheckRequiresSorted_Mode_COLLECT_AND_REPORT()))))
	})
	lambdaSingleton_default_pass_config_34 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewConstCheck(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_default_pass_config_35 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewRewriteCallerCodeLocation(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_36 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewReplaceToggles(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_37 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_49 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_49)
		_ = compiler
		return NewClosureOptimizePrimitives(compiler, compiler.GetOptions().GetOutputFeatureSet().Contains_1_FeatureSet(pkg_parser.StaticGet_FeatureSet_ES2015()))
	})
	lambdaSingleton_default_pass_config_38 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewConvertChunksToESModules(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_39 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_51 any) any {
		options := j2gort.CastPtr[CompilerOptions](lp_51)
		_ = options
		return NewPassFactory_PreconditionResult(((options.GetPropertyCollapseLevel() == StaticGet_CompilerOptions_PropertyCollapseLevel_NONE()) || options.GetAssumePropertiesAreStaticallyAnalyzable()), "requires assumePropertiesAreStaticallyAnalyzable to be enabled")
	})
	lambdaSingleton_default_pass_config_40 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_53 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_53)
		_ = compiler
		return NewInlineObjectLiterals(compiler, j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](compiler.GetUniqueNameIdSupplier()))
	})
	lambdaSingleton_default_pass_config_41 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(DefaultPassConfig_RequirePropertiesAreStaticallyAnalyzable(j2gort.CastPtr[CompilerOptions](mrp0_)))
	})
	lambdaSingleton_default_pass_config_42 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(DefaultPassConfig_RequirePropertiesAreStaticallyAnalyzable(j2gort.CastPtr[CompilerOptions](mrp0_)))
	})
	lambdaSingleton_default_pass_config_43 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_55 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_55)
		_ = compiler
		return j2gort.PtrToAny(OptimizeCalls_S_Builder().SetCompiler(compiler).SetConsiderExterns(false).AddPass(NewDevirtualizeMethods(compiler)).Build())
	})
	lambdaSingleton_default_pass_config_44 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(DefaultPassConfig_RequirePropertiesAreStaticallyAnalyzable(j2gort.CastPtr[CompilerOptions](mrp0_)))
	})
	lambdaSingleton_default_pass_config_45 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_56 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_56)
		_ = compiler
		return j2gort.PtrToAny(OptimizeCalls_S_Builder().SetCompiler(compiler).SetConsiderExterns(false).AddPass(NewOptimizeReturns(compiler)).AddPass(NewOptimizeParameters(compiler)).Build())
	})
	lambdaSingleton_default_pass_config_46 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_57 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_57)
		_ = compiler
		return j2gort.PtrToAny(OptimizeCalls_S_Builder().SetCompiler(compiler).SetConsiderExterns(false).AddPass(NewOptimizeConstructors(compiler)).Build())
	})
	lambdaSingleton_default_pass_config_47 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(DefaultPassConfig_RequirePropertiesAreStaticallyAnalyzable(j2gort.CastPtr[CompilerOptions](mrp0_)))
	})
	lambdaSingleton_default_pass_config_48 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewPureFunctionIdentifier_Driver(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_49 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_59 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_59)
		_ = compiler
		return NewInlineVariables(compiler, StaticGet_InlineVariables_Mode_CONSTANTS_ONLY())
	})
	lambdaSingleton_default_pass_config_50 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(DefaultPassConfig_RequirePropertiesAreStaticallyAnalyzable(j2gort.CastPtr[CompilerOptions](mrp0_)))
	})
	lambdaSingleton_default_pass_config_51 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewInlineSimpleMethods(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_52 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewDeadAssignmentsElimination(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_53 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(DefaultPassConfig_RequirePropertiesAreStaticallyAnalyzable(j2gort.CastPtr[CompilerOptions](mrp0_)))
	})
	lambdaSingleton_default_pass_config_54 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewDeadPropertyAssignmentElimination(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_55 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewInlineProperties(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_56 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewIsolatePolyfills_1(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_57 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewFlowSensitiveInlineVariables(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_58 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewCollapseVariableDeclarations(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_59 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(DefaultPassConfig_RequirePropertiesAreStaticallyAnalyzable(j2gort.CastPtr[CompilerOptions](mrp0_)))
	})
	lambdaSingleton_default_pass_config_60 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewFunctionRewriter(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_61 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewCollapseAnonymousFunctions(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_62 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewRewriteGlobalDeclarationsForTryCatchWrapping(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_63 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(DefaultPassConfig_RequirePropertiesAreStaticallyAnalyzable(j2gort.CastPtr[CompilerOptions](mrp0_)))
	})
	lambdaSingleton_default_pass_config_64 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(Normalize_CreateNormalizeForOptimizations(j2gort.AsIface[AbstractCompiler](mrp0_)))
	})
	lambdaSingleton_default_pass_config_65 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewExternExportsPass(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_66 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return MakeDeclaredNamesUnique_GetContextualRenameInverter(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_67 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(DefaultPassConfig_RequirePropertiesAreStaticallyAnalyzable(j2gort.CastPtr[CompilerOptions](mrp0_)))
	})
	lambdaSingleton_default_pass_config_68 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewRemovePropertyRenamingCalls(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_69 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewRenameLabels_1(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_70 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewConvertToDottedProperties(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_71 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewAstValidator_1(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_72 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_73 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_73)
		_ = compiler
		return NewVarCheck_2(compiler, true)
	})
	lambdaSingleton_default_pass_config_73 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewPolymerPass(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_default_pass_config_74 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewChromePass(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_default_pass_config_75 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewJ2clConstantHoisterPass(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_76 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_75 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_75)
		_ = compiler
		var changedScopeNodes j2gort.Ext_java_util_List = compiler.GetChangeTracker().GetChangedScopeNodesForPass("j2clClinitPass")
		_ = changedScopeNodes
		return NewJ2clClinitPrunerPass(compiler, changedScopeNodes)
	})
	lambdaSingleton_default_pass_config_77 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewJ2clPropertyInlinerPass(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_78 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewJ2clPass(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_default_pass_config_79 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewJ2clUtilGetDefineRewriterPass(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_80 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewJ2clAssertRemovalPass(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_81 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewJ2clSourceFileChecker(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_82 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewJ2clChecksPass(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_83 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewRemoveCastNodes(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_84 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_77 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_77)
		_ = compiler
		return NewLambda_CompilerPass(func(externs pkg_rhino.Node, js pkg_rhino.Node) {
			_ = externs
			_ = js
			{
				var val_78 ConvertTypesToColors
				InitConvertTypesToColors(&val_78, compiler, SerializationOptions_S_Builder().SetIncludeDebugInfo(compiler.GetOptions().ShouldSerializeExtraDebugInfo()).Build())
				(&val_78).Process_2(externs, js)
			}
			compiler.SetLifeCycleStage(StaticGet_AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC())
		})
	})
	lambdaSingleton_default_pass_config_85 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewRemoveUnnecessarySyntheticExterns(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_86 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewSubstituteEs6Syntax(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_87 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewWhitespaceWrapGoogModules(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_88 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewProcessCommonJSModules(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_89 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewEs6RewriteScriptsToModules(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_90 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_81 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_81)
		_ = compiler
		return NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap())
	})
	lambdaSingleton_default_pass_config_91 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewGatherGetterAndSetterProperties(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_default_pass_config_92 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_82 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_82)
		_ = compiler
		return NewLambda_CompilerPass(func(externs pkg_rhino.Node, js pkg_rhino.Node) {
			_ = externs
			_ = js
			compiler.InitializeSyntheticCodeInput()
		})
	})
	lambdaSingleton_default_pass_config_93 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_83 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_83)
		_ = compiler
		return NewLambda_CompilerPass(func(externs pkg_rhino.Node, js pkg_rhino.Node) {
			_ = externs
			_ = js
			compiler.RemoveSyntheticCodeInput_0()
		})
	})
	lambdaSingleton_default_pass_config_94 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_84 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_84)
		_ = compiler
		return NewLambda_CompilerPass(func(externs pkg_rhino.Node, js pkg_rhino.Node) {
			_ = externs
			_ = js
			compiler.MergeSyntheticCodeInput()
		})
	})
	lambdaSingleton_default_pass_config_95 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_85 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_85)
		_ = compiler
		return NewLambda_CompilerPass(func(externs pkg_rhino.Node, root pkg_rhino.Node) {
			_ = externs
			_ = root
			NewLocaleDataPasses_ProtectGoogLocale(compiler).Process_2(externs, root)
		})
	})
	lambdaSingleton_default_pass_config_96 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_87 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_87)
		_ = compiler
		return NewTranspileAndOptimizeClosureUnaware(compiler, StaticGet_NestedCompilerRunner_Mode_TRANSPILE_ONLY())
	})
	lambdaSingleton_default_pass_config_97 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_88 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_88)
		_ = compiler
		return NewTranspileAndOptimizeClosureUnaware(compiler, StaticGet_NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE())
	})
}

// Java class: com.google.javascript.jscomp.DefaultPassConfig
type DefaultPassConfig struct {
	PassConfigBase
	preprocessorSymbolTableFactory               *PreprocessorSymbolTable_CachedInstanceFactory `gson:"-"`
	checkSideEffects                             PassFactory                                    `gson:"checkSideEffects"`
	stripSideEffectProtection                    PassFactory                                    `gson:"stripSideEffectProtection"`
	suspiciousCode                               PassFactory                                    `gson:"suspiciousCode"`
	extraRequires                                PassFactory                                    `gson:"extraRequires"`
	checkMissingRequires                         PassFactory                                    `gson:"checkMissingRequires"`
	checkJsDocAndEs6Modules                      PassFactory                                    `gson:"checkJsDocAndEs6Modules"`
	generateExports                              PassFactory                                    `gson:"generateExports"`
	generateIjs                                  PassFactory                                    `gson:"generateIjs"`
	removeExtraRequires                          PassFactory                                    `gson:"removeExtraRequires"`
	exportTestFunctions                          PassFactory                                    `gson:"exportTestFunctions"`
	gatherRawExports                             PassFactory                                    `gson:"gatherRawExports"`
	closurePrimitives                            PassFactory                                    `gson:"closurePrimitives"`
	closureProvidesRequires                      PassFactory                                    `gson:"closureProvidesRequires"`
	angularPass                                  PassFactory                                    `gson:"angularPass"`
	replaceMessagesForChrome                     PassFactory                                    `gson:"replaceMessagesForChrome"`
	closureGoogScopeAliasesForIjs                PassFactory                                    `gson:"closureGoogScopeAliasesForIjs"`
	closureGoogScopeAliases                      PassFactory                                    `gson:"closureGoogScopeAliases"`
	injectRuntimeLibraries                       PassFactory                                    `gson:"injectRuntimeLibraries"`
	removeWeakSources                            PassFactory                                    `gson:"removeWeakSources"`
	declaredGlobalExternsOnWindow                PassFactory                                    `gson:"declaredGlobalExternsOnWindow"`
	checkTypeImportCodeReferences                PassFactory                                    `gson:"checkTypeImportCodeReferences"`
	closureCheckModule                           PassFactory                                    `gson:"closureCheckModule"`
	closureRewriteModule                         PassFactory                                    `gson:"closureRewriteModule"`
	checkClosureImports                          PassFactory                                    `gson:"checkClosureImports"`
	rewriteGoogJsImports                         PassFactory                                    `gson:"rewriteGoogJsImports"`
	closureReplaceGetCssName                     PassFactory                                    `gson:"closureReplaceGetCssName"`
	createSyntheticBlocks                        PassFactory                                    `gson:"createSyntheticBlocks"`
	earlyPeepholeOptimizations                   PassFactory                                    `gson:"earlyPeepholeOptimizations"`
	earlyInlineVariables                         PassFactory                                    `gson:"earlyInlineVariables"`
	peepholeOptimizations                        PassFactory                                    `gson:"peepholeOptimizations"`
	peepholeOptimizationsOnceNormalized          PassFactory                                    `gson:"peepholeOptimizationsOnceNormalized"`
	peepholeOptimizationsOnceNonNormalized       PassFactory                                    `gson:"peepholeOptimizationsOnceNonNormalized"`
	latePeepholeOptimizations                    PassFactory                                    `gson:"latePeepholeOptimizations"`
	checkVars                                    PassFactory                                    `gson:"checkVars"`
	inferConsts                                  PassFactory                                    `gson:"inferConsts"`
	checkRegExp                                  PassFactory                                    `gson:"checkRegExp"`
	checkRegExpForOptimizations                  PassFactory                                    `gson:"checkRegExpForOptimizations"`
	checkVariableReferences                      PassFactory                                    `gson:"checkVariableReferences"`
	checkSuper                                   PassFactory                                    `gson:"checkSuper"`
	clearTypedScopeCreatorPass                   PassFactory                                    `gson:"clearTypedScopeCreatorPass"`
	clearTopTypedScopePass                       PassFactory                                    `gson:"clearTopTypedScopePass"`
	inferTypes                                   PassFactory                                    `gson:"inferTypes"`
	inferJsDocInfo                               PassFactory                                    `gson:"inferJsDocInfo"`
	checkTypes                                   PassFactory                                    `gson:"checkTypes"`
	checkControlFlow                             PassFactory                                    `gson:"checkControlFlow"`
	checkAccessControls                          PassFactory                                    `gson:"checkAccessControls"`
	lintChecks                                   PassFactory                                    `gson:"lintChecks"`
	analyzerChecks                               PassFactory                                    `gson:"analyzerChecks"`
	checkRequiresAndProvidesSorted               PassFactory                                    `gson:"checkRequiresAndProvidesSorted"`
	checkStrictMode                              PassFactory                                    `gson:"checkStrictMode"`
	processTweaks                                PassFactory                                    `gson:"processTweaks"`
	processDefinesCheck                          PassFactory                                    `gson:"processDefinesCheck"`
	processDefinesOptimize                       PassFactory                                    `gson:"processDefinesOptimize"`
	stripCode                                    PassFactory                                    `gson:"stripCode"`
	checkConsts                                  PassFactory                                    `gson:"checkConsts"`
	rewriteCallerCodeLocation                    PassFactory                                    `gson:"rewriteCallerCodeLocation"`
	replaceToggles                               PassFactory                                    `gson:"replaceToggles"`
	replaceIdGenerators                          PassFactory                                    `gson:"replaceIdGenerators"`
	replaceStrings                               PassFactory                                    `gson:"replaceStrings"`
	closureCodeRemoval                           PassFactory                                    `gson:"closureCodeRemoval"`
	closureOptimizePrimitives                    PassFactory                                    `gson:"closureOptimizePrimitives"`
	rescopeGlobalSymbols                         PassFactory                                    `gson:"rescopeGlobalSymbols"`
	convertChunksToESModules                     PassFactory                                    `gson:"convertChunksToESModules"`
	inlineAndCollapseProperties                  PassFactory                                    `gson:"inlineAndCollapseProperties"`
	collapseObjectLiterals                       PassFactory                                    `gson:"collapseObjectLiterals"`
	disambiguateProperties                       PassFactory                                    `gson:"disambiguateProperties"`
	devirtualizeMethods                          PassFactory                                    `gson:"devirtualizeMethods"`
	optimizeCalls                                PassFactory                                    `gson:"optimizeCalls"`
	optimizeConstructors                         PassFactory                                    `gson:"optimizeConstructors"`
	markPureFunctions                            PassFactory                                    `gson:"markPureFunctions"`
	inlineVariables                              PassFactory                                    `gson:"inlineVariables"`
	inlineConstants                              PassFactory                                    `gson:"inlineConstants"`
	inlineSimpleMethods                          PassFactory                                    `gson:"inlineSimpleMethods"`
	deadAssignmentsElimination                   PassFactory                                    `gson:"deadAssignmentsElimination"`
	deadPropertyAssignmentElimination            PassFactory                                    `gson:"deadPropertyAssignmentElimination"`
	inlineFunctions                              PassFactory                                    `gson:"inlineFunctions"`
	inlineProperties                             PassFactory                                    `gson:"inlineProperties"`
	isolatePolyfills                             PassFactory                                    `gson:"isolatePolyfills"`
	removeUnusedCode                             PassFactory                                    `gson:"removeUnusedCode"`
	removeUnusedCodeOnce                         PassFactory                                    `gson:"removeUnusedCodeOnce"`
	crossChunkCodeMotion                         PassFactory                                    `gson:"crossChunkCodeMotion"`
	crossChunkMethodMotion                       PassFactory                                    `gson:"crossChunkMethodMotion"`
	flowSensitiveInlineVariables                 PassFactory                                    `gson:"flowSensitiveInlineVariables"`
	coalesceVariableNames                        PassFactory                                    `gson:"coalesceVariableNames"`
	postNormalizePeephole                        PassFactory                                    `gson:"postNormalizePeephole"`
	collapseVariableDeclarations                 PassFactory                                    `gson:"collapseVariableDeclarations"`
	extractPrototypeMemberDeclarations           PassFactory                                    `gson:"extractPrototypeMemberDeclarations"`
	rewriteFunctionExpressions                   PassFactory                                    `gson:"rewriteFunctionExpressions"`
	collapseAnonymousFunctions                   PassFactory                                    `gson:"collapseAnonymousFunctions"`
	rewriteGlobalDeclarationsForTryCatchWrapping PassFactory                                    `gson:"rewriteGlobalDeclarationsForTryCatchWrapping"`
	aliasStrings                                 PassFactory                                    `gson:"aliasStrings"`
	ambiguateProperties                          PassFactory                                    `gson:"ambiguateProperties"`
	markUnnormalized                             PassFactory                                    `gson:"markUnnormalized"`
	normalize                                    PassFactory                                    `gson:"normalize"`
	externExports                                PassFactory                                    `gson:"externExports"`
	denormalize                                  PassFactory                                    `gson:"denormalize"`
	invertContextualRenaming                     PassFactory                                    `gson:"invertContextualRenaming"`
	renameProperties                             PassFactory                                    `gson:"renameProperties"`
	removePropertyRenamingCalls                  PassFactory                                    `gson:"removePropertyRenamingCalls"`
	renameVars                                   PassFactory                                    `gson:"renameVars"`
	renameLabels                                 PassFactory                                    `gson:"renameLabels"`
	convertToDottedProperties                    PassFactory                                    `gson:"convertToDottedProperties"`
	checkAstValidity                             PassFactory                                    `gson:"checkAstValidity"`
	varCheckValidity                             PassFactory                                    `gson:"varCheckValidity"`
	instrumentForCodeCoverage                    PassFactory                                    `gson:"instrumentForCodeCoverage"`
	gatherExternPropertiesCheck                  PassFactory                                    `gson:"gatherExternPropertiesCheck"`
	gatherExternPropertiesOptimize               PassFactory                                    `gson:"gatherExternPropertiesOptimize"`
	polymerPass                                  PassFactory                                    `gson:"polymerPass"`
	chromePass                                   PassFactory                                    `gson:"chromePass"`
	j2clConstantHoisterPass                      PassFactory                                    `gson:"j2clConstantHoisterPass"`
	j2clClinitPass                               PassFactory                                    `gson:"j2clClinitPass"`
	j2clPropertyInlinerPass                      PassFactory                                    `gson:"j2clPropertyInlinerPass"`
	j2clPass                                     PassFactory                                    `gson:"j2clPass"`
	j2clUtilGetDefineRewriterPass                PassFactory                                    `gson:"j2clUtilGetDefineRewriterPass"`
	j2clAssertRemovalPass                        PassFactory                                    `gson:"j2clAssertRemovalPass"`
	j2clSourceFileChecker                        PassFactory                                    `gson:"j2clSourceFileChecker"`
	j2clChecksPass                               PassFactory                                    `gson:"j2clChecksPass"`
	checkConformance                             PassFactory                                    `gson:"checkConformance"`
	removeCastNodes                              PassFactory                                    `gson:"removeCastNodes"`
	typesToColors                                PassFactory                                    `gson:"typesToColors"`
	serializeTypedAst                            PassFactory                                    `gson:"serializeTypedAst"`
	removeUnnecessarySyntheticExterns            PassFactory                                    `gson:"removeUnnecessarySyntheticExterns"`
	optimizeToEs6                                PassFactory                                    `gson:"optimizeToEs6"`
	whitespaceWrapGoogModules                    PassFactory                                    `gson:"whitespaceWrapGoogModules"`
	rewriteCommonJsModules                       PassFactory                                    `gson:"rewriteCommonJsModules"`
	rewriteScriptsToEs6Modules                   PassFactory                                    `gson:"rewriteScriptsToEs6Modules"`
	gatherModuleMetadataPass                     PassFactory                                    `gson:"gatherModuleMetadataPass"`
	createModuleMapPass                          PassFactory                                    `gson:"createModuleMapPass"`
	gatherGettersAndSetters                      PassFactory                                    `gson:"gatherGettersAndSetters"`
	addSyntheticScript                           PassFactory                                    `gson:"addSyntheticScript"`
	removeSyntheticScript                        PassFactory                                    `gson:"removeSyntheticScript"`
	mergeSyntheticScript                         PassFactory                                    `gson:"mergeSyntheticScript"`
	protectLocaleData                            PassFactory                                    `gson:"protectLocaleData"`
	substituteLocaleData                         PassFactory                                    `gson:"substituteLocaleData"`
	transpileOnlyClosureUnaware                  PassFactory                                    `gson:"transpileOnlyClosureUnaware"`
	transpileAndOptimizeClosureUnaware           PassFactory                                    `gson:"transpileAndOptimizeClosureUnaware"`
}

func (this_ *DefaultPassConfig) BindSelf_(s any) {
	_ = s
	this_.PassConfigBase.BindSelf_(s)
}
func (this_ *DefaultPassConfig) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig", "DefaultPassConfig"
}
func (this_ *DefaultPassConfig) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PassConfigBase.IdentityHash_()
}
func (*DefaultPassConfig) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DefaultPassConfig", "com.google.javascript.jscomp.DefaultPassConfig", func(v_ any) bool {
		return j2gort.InstanceOf[*DefaultPassConfig](v_)
	})
}
func (*DefaultPassConfig) Implements_jscomp_PassConfig() {}

func StaticGet_DefaultPassConfig_COMPILED_CONSTANT_NAME() string       { return "COMPILED" }
func StaticGet_DefaultPassConfig_CLOSURE_LOCALE_CONSTANT_NAME() string { return "goog.LOCALE" }

var static_DefaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_DefaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_DefaultPassConfig()
	return static_DefaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN
}
func StaticSet_DefaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_DefaultPassConfig()
	static_DefaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN = v
}

var static_DefaultPassConfig_GENERATE_EXPORTS_ERROR *DiagnosticType

func StaticGet_DefaultPassConfig_GENERATE_EXPORTS_ERROR() *DiagnosticType {
	clinit_DefaultPassConfig()
	return static_DefaultPassConfig_GENERATE_EXPORTS_ERROR
}
func StaticSet_DefaultPassConfig_GENERATE_EXPORTS_ERROR(v *DiagnosticType) {
	clinit_DefaultPassConfig()
	static_DefaultPassConfig_GENERATE_EXPORTS_ERROR = v
}

var static_DefaultPassConfig_forbidDynamicImportUsage PassFactory

func StaticGet_DefaultPassConfig_ForbidDynamicImportUsage() PassFactory {
	clinit_DefaultPassConfig()
	return static_DefaultPassConfig_forbidDynamicImportUsage
}
func StaticSet_DefaultPassConfig_ForbidDynamicImportUsage(v PassFactory) {
	clinit_DefaultPassConfig()
	static_DefaultPassConfig_forbidDynamicImportUsage = v
}

var static_DefaultPassConfig_rewriteDynamicImports PassFactory

func StaticGet_DefaultPassConfig_RewriteDynamicImports() PassFactory {
	clinit_DefaultPassConfig()
	return static_DefaultPassConfig_rewriteDynamicImports
}
func StaticSet_DefaultPassConfig_RewriteDynamicImports(v PassFactory) {
	clinit_DefaultPassConfig()
	static_DefaultPassConfig_rewriteDynamicImports = v
}

var clinitDone_DefaultPassConfig uint32
var clinitState_DefaultPassConfig j2gort.ClinitState

func clinit_DefaultPassConfig() {
	if atomic.LoadUint32(&clinitDone_DefaultPassConfig) != 0 {
		return
	}
	clinitSlow_DefaultPassConfig()
}

func clinitSlow_DefaultPassConfig() {
	if !j2gort.EnterClinit(&clinitState_DefaultPassConfig, &clinitDone_DefaultPassConfig) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DefaultPassConfig, &clinitDone_DefaultPassConfig)
	static_DefaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("^[a-zA-Z0-9$_]+$")
	static_DefaultPassConfig_GENERATE_EXPORTS_ERROR = DiagnosticType_Error("JSC_GENERATE_EXPORTS_ERROR", "Exports can only be generated if export symbol/property functions are set.")
	static_DefaultPassConfig_forbidDynamicImportUsage = PassFactory_S_Builder().SetName("FORBID_DYNAMIC_IMPORT").SetInternalFactory(lambdaSingleton_default_pass_config_1).GetPassFactory_BuilderBase().Build()
	static_DefaultPassConfig_rewriteDynamicImports = PassFactory_S_Builder().SetName("REWRITE_DYNAMIC_IMPORT").SetInternalFactory(lambdaSingleton_default_pass_config_2).GetPassFactory_BuilderBase().Build()
}

func (this_ *DefaultPassConfig) initFields_DefaultPassConfig() {
	_ = this_
	this_.preprocessorSymbolTableFactory = NewPreprocessorSymbolTable_CachedInstanceFactory()
	this_.checkSideEffects = PassFactory_S_Builder().SetName("checkSideEffects").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_2 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_2)
		_ = compiler
		return NewCheckSideEffects(compiler, this_.GetPassConfigBase().options.GetCheckSuspiciousCode(), this_.GetPassConfigBase().options.ShouldProtectHiddenSideEffects())
	})).GetPassFactory_BuilderBase().Build()
	this_.stripSideEffectProtection = PassFactory_S_Builder().SetName("stripSideEffectProtection").SetInternalFactory(lambdaSingleton_default_pass_config_3).GetPassFactory_BuilderBase().Build()
	this_.suspiciousCode = PassFactory_S_Builder().SetName("suspiciousCode").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_3 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_3)
		_ = compiler
		var sharedCallbacks j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
		_ = sharedCallbacks
		if this_.GetPassConfigBase().options.GetCheckSuspiciousCode() {
			sharedCallbacks.Add_1(NewCheckSuspiciousCode())
			sharedCallbacks.Add_1(NewCheckDuplicateCase(compiler))
		}
		if this_.GetPassConfigBase().options.Enables(StaticGet_DiagnosticGroups_GLOBAL_THIS()) {
			sharedCallbacks.Add_1(NewCheckGlobalThis(compiler))
		}
		if this_.GetPassConfigBase().options.Enables(StaticGet_DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT()) {
			sharedCallbacks.Add_1(NewCheckDebuggerStatement(compiler))
		}
		return DefaultPassConfig_CombineChecks(compiler, sharedCallbacks)
	})).GetPassFactory_BuilderBase().Build()
	this_.extraRequires = PassFactory_S_Builder().SetName("checkExtraRequires").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_4 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_4)
		_ = compiler
		return NewCheckExtraRequires(compiler, this_.GetPassConfigBase().options.GetUnusedImportsToRemove())
	})).GetPassFactory_BuilderBase().Build()
	this_.checkMissingRequires = PassFactory_S_Builder().SetName("checkMissingRequires").SetInternalFactory(lambdaSingleton_default_pass_config_4).GetPassFactory_BuilderBase().Build()
	this_.checkJsDocAndEs6Modules = PassFactory_S_Builder().SetName("checkJsDocAndEs6Modules").SetInternalFactory(lambdaSingleton_default_pass_config_5).GetPassFactory_BuilderBase().Build()
	this_.generateExports = PassFactory_S_Builder().SetName("generateExports").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_7 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_7)
		_ = compiler
		var convention CodingConvention = compiler.GetCodingConvention()
		_ = convention
		clinit_GenerateExports()
		var val_8 GenerateExports
		InitGenerateExports(&val_8, compiler, this_.GetPassConfigBase().options.ShouldExportLocalPropertyDefinitions(), convention.GetExportSymbolFunction(), convention.GetExportPropertyFunction())
		var pass *GenerateExports = &val_8
		_ = pass
		return NewDefaultPassConfig__Anon_22(this_, pass, compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.generateIjs = PassFactory_S_Builder().SetName("generateIjs").SetInternalFactory(lambdaSingleton_default_pass_config_6).GetPassFactory_BuilderBase().Build()
	this_.removeExtraRequires = PassFactory_S_Builder().SetName("removeExtraRequires").SetInternalFactory(lambdaSingleton_default_pass_config_7).GetPassFactory_BuilderBase().Build()
	this_.exportTestFunctions = PassFactory_S_Builder().SetName("exportTestFunctions").SetInternalFactory(lambdaSingleton_default_pass_config_8).GetPassFactory_BuilderBase().Build()
	this_.gatherRawExports = PassFactory_S_Builder().SetName("gatherRawExports").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_10 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_10)
		_ = compiler
		var pass *GatherRawExports = NewGatherRawExports(compiler)
		_ = pass
		return NewDefaultPassConfig__Anon_23(this_, pass, compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.closurePrimitives = PassFactory_S_Builder().SetName("closurePrimitives").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_11 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_11)
		_ = compiler
		this_.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		return NewProcessClosurePrimitives(compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.closureProvidesRequires = PassFactory_S_Builder().SetName("closureProvidesRequires").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_12 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_12)
		_ = compiler
		this_.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		clinit_ProcessClosureProvidesAndRequires()
		var val_13 ProcessClosureProvidesAndRequires
		InitProcessClosureProvidesAndRequires(&val_13, compiler, this_.GetPassConfigBase().options.ShouldPreservesGoogProvidesAndRequires())
		var pass *ProcessClosureProvidesAndRequires = &val_13
		_ = pass
		return NewLambda_CompilerPass(func(externs pkg_rhino.Node, root pkg_rhino.Node) {
			_ = externs
			_ = root
			pass.Process_2(externs, root)
			compiler.AddExportedNames(pass.GetExportedVariableNames())
		})
	})).GetPassFactory_BuilderBase().Build()
	this_.angularPass = PassFactory_S_Builder().SetName("angularPass").SetInternalFactory(lambdaSingleton_default_pass_config_9).GetPassFactory_BuilderBase().Build()
	this_.replaceMessagesForChrome = PassFactory_S_Builder().SetName("replaceMessages").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_14 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_14)
		_ = compiler
		return NewReplaceMessagesForChrome(compiler, NewGoogleJsMessageIdGenerator(this_.GetPassConfigBase().options.GetTcProjectId()))
	})).GetPassFactory_BuilderBase().Build()
	this_.closureGoogScopeAliasesForIjs = PassFactory_S_Builder().SetName("closureGoogScopeAliasesForIjs").SetInternalFactory(lambdaSingleton_default_pass_config_10).GetPassFactory_BuilderBase().Build()
	this_.closureGoogScopeAliases = PassFactory_S_Builder().SetName("closureGoogScopeAliases").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_16 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_16)
		_ = compiler
		this_.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		return j2gort.PtrToAny(ScopedAliases_S_Builder(compiler).SetPreprocessorSymbolTable(this_.preprocessorSymbolTableFactory.GetInstanceOrNull()).SetModuleMetadataMap(compiler.GetModuleMetadataMap()).SetInvalidModuleGetHandling(StaticGet_ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME()).Build())
	})).GetPassFactory_BuilderBase().Build()
	this_.injectRuntimeLibraries = PassFactory_S_Builder().SetName("InjectRuntimeLibraries").SetInternalFactory(lambdaSingleton_default_pass_config_11).GetPassFactory_BuilderBase().Build()
	this_.removeWeakSources = PassFactory_S_Builder().SetName("removeWeakSources").SetInternalFactory(lambdaSingleton_default_pass_config_12).GetPassFactory_BuilderBase().Build()
	this_.declaredGlobalExternsOnWindow = PassFactory_S_Builder().SetName("declaredGlobalExternsOnWindow").SetInternalFactory(lambdaSingleton_default_pass_config_13).GetPassFactory_BuilderBase().Build()
	this_.checkTypeImportCodeReferences = PassFactory_S_Builder().SetName("checkTypeImportCodeReferences").SetInternalFactory(lambdaSingleton_default_pass_config_14).GetPassFactory_BuilderBase().Build()
	this_.closureCheckModule = PassFactory_S_Builder().SetName("closureCheckModule").SetInternalFactory(lambdaSingleton_default_pass_config_15).GetPassFactory_BuilderBase().Build()
	this_.closureRewriteModule = PassFactory_S_Builder().SetName("closureRewriteModule").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_19 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_19)
		_ = compiler
		this_.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		return NewClosureRewriteModule(compiler, this_.preprocessorSymbolTableFactory.GetInstanceOrNull(), compiler.GetTopScope())
	})).GetPassFactory_BuilderBase().Build()
	this_.checkClosureImports = PassFactory_S_Builder().SetName("checkGoogRequires").SetInternalFactory(lambdaSingleton_default_pass_config_16).GetPassFactory_BuilderBase().Build()
	this_.rewriteGoogJsImports = PassFactory_S_Builder().SetName("rewriteGoogJsImports").SetInternalFactory(lambdaSingleton_default_pass_config_17).GetPassFactory_BuilderBase().Build()
	this_.closureReplaceGetCssName = PassFactory_S_Builder().SetName("closureReplaceGetCssName").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_22 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_22)
		_ = compiler
		return NewDefaultPassConfig__Anon_24(this_, compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.createSyntheticBlocks = PassFactory_S_Builder().SetName("createSyntheticBlocks").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_23 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_23)
		_ = compiler
		return NewCreateSyntheticBlocks(compiler, this_.GetPassConfigBase().options.GetSyntheticBlockStartMarker(), this_.GetPassConfigBase().options.GetSyntheticBlockEndMarker())
	})).GetPassFactory_BuilderBase().Build()
	this_.earlyPeepholeOptimizations = PassFactory_S_Builder().SetName("earlyPeepholeOptimizations").SetInternalFactory(lambdaSingleton_default_pass_config_18).GetPassFactory_BuilderBase().Build()
	this_.earlyInlineVariables = PassFactory_S_Builder().SetName("earlyInlineVariables").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_25 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_25)
		_ = compiler
		var mode *InlineVariables_Mode = nil
		_ = mode
		if this_.GetPassConfigBase().options.ShouldInlineVariables() {
			mode = StaticGet_InlineVariables_Mode_ALL()
		} else if this_.GetPassConfigBase().options.ShouldInlineLocalVariables() {
			mode = StaticGet_InlineVariables_Mode_LOCALS_ONLY()
		} else {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "No variable inlining option set.")))
		}
		return NewInlineVariables(compiler, mode)
	})).GetPassFactory_BuilderBase().Build()
	this_.peepholeOptimizations = PassFactory_S_Builder().SetName("peepholeOptimizations").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_default_pass_config_19).GetPassFactory_BuilderBase().Build()
	this_.peepholeOptimizationsOnceNormalized = PassFactory_S_Builder().SetName("peepholeOptimizations").SetInternalFactory(lambdaSingleton_default_pass_config_20).GetPassFactory_BuilderBase().Build()
	this_.peepholeOptimizationsOnceNonNormalized = PassFactory_S_Builder().SetName("peepholeOptimizations").SetInternalFactory(lambdaSingleton_default_pass_config_21).GetPassFactory_BuilderBase().Build()
	this_.latePeepholeOptimizations = PassFactory_S_Builder().SetName("latePeepholeOptimizations").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_29 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_29)
		_ = compiler
		var late bool = true
		_ = late
		var useTypesForOptimization bool = this_.GetPassConfigBase().options.ShouldUseTypesForLocalOptimization()
		_ = useTypesForOptimization
		return NewPeepholeOptimizationsPass_3_AbstractCompiler_string_AbstractPeepholeOptimization(compiler, "latePeepholeOptimizations", []AbstractPeepholeOptimization{NewStatementFusion(), NewPeepholeRemoveDeadCode(), NewPeepholeMinimizeConditions(late), NewPeepholeSubstituteAlternateSyntax(late), NewPeepholeReplaceKnownMethods(late, useTypesForOptimization), NewPeepholeFoldConstants(late, useTypesForOptimization)})
	})).GetPassFactory_BuilderBase().Build()
	this_.checkVars = PassFactory_S_Builder().SetName("checkVars").SetInternalFactory(lambdaSingleton_default_pass_config_22).GetPassFactory_BuilderBase().Build()
	this_.inferConsts = PassFactory_S_Builder().SetName("inferConsts").SetInternalFactory(lambdaSingleton_default_pass_config_23).GetPassFactory_BuilderBase().Build()
	this_.checkRegExp = PassFactory_S_Builder().SetName("checkRegExp").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_30 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_30)
		_ = compiler
		var pass *CheckRegExp = NewCheckRegExp(compiler, this_.GetPassConfigBase().options.GetAssumePropertiesAreStaticallyAnalyzable(), true)
		_ = pass
		return NewDefaultPassConfig__Anon_25(this_, pass, compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.checkRegExpForOptimizations = PassFactory_S_Builder().SetName("checkRegExpForOptimizations").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_31 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_31)
		_ = compiler
		var regExpCheck *CheckRegExp = NewCheckRegExp(compiler, this_.GetPassConfigBase().options.GetAssumePropertiesAreStaticallyAnalyzable(), false)
		_ = regExpCheck
		return NewDefaultPassConfig__Anon_26(this_, regExpCheck, compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.checkVariableReferences = PassFactory_S_Builder().SetName("checkVariableReferences").SetInternalFactory(lambdaSingleton_default_pass_config_24).GetPassFactory_BuilderBase().Build()
	this_.checkSuper = PassFactory_S_Builder().SetName("checkSuper").SetInternalFactory(lambdaSingleton_default_pass_config_25).GetPassFactory_BuilderBase().Build()
	this_.clearTypedScopeCreatorPass = PassFactory_S_Builder().SetName("clearTypedScopeCreatorPass").SetInternalFactory(lambdaSingleton_default_pass_config_26).GetPassFactory_BuilderBase().Build()
	this_.clearTopTypedScopePass = PassFactory_S_Builder().SetName("clearTopTypedScopePass").SetInternalFactory(lambdaSingleton_default_pass_config_27).GetPassFactory_BuilderBase().Build()
	this_.inferTypes = PassFactory_S_Builder().SetName("inferTypes").SetInternalFactory(lambdaSingleton_default_pass_config_28).GetPassFactory_BuilderBase().Build()
	this_.inferJsDocInfo = PassFactory_S_Builder().SetName("inferJsDocInfo").SetInternalFactory(lambdaSingleton_default_pass_config_29).GetPassFactory_BuilderBase().Build()
	this_.checkTypes = PassFactory_S_Builder().SetName("checkTypes").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_38 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_38)
		_ = compiler
		return NewLambda_CompilerPass(func(externs pkg_rhino.Node, root pkg_rhino.Node) {
			_ = externs
			_ = root
			var check *TypeCheck = NewTypeCheck_5(compiler, compiler.GetReverseAbstractInterpreter(), compiler.GetTypeRegistry(), compiler.GetTopScope(), j2gort.CastPtr[TypedScopeCreator](compiler.GetTypedScopeCreator())).ReportUnknownTypes(this_.GetPassConfigBase().options.Enables(DiagnosticGroup_ForType(StaticGet_TypeCheck_UNKNOWN_EXPR_TYPE()))).ReportMissingProperties(!(this_.GetPassConfigBase().options.Disables(DiagnosticGroup_ForType(StaticGet_TypeCheck_INEXISTENT_PROPERTY()))))
			_ = check
			check.Process_2(externs, root)
			compiler.GetErrorManager().SetTypedPercent(check.GetTypedPercent())
		})
	})).GetPassFactory_BuilderBase().Build()
	this_.checkControlFlow = PassFactory_S_Builder().SetName("checkControlFlow").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_39 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_39)
		_ = compiler
		var callbacks j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
		_ = callbacks
		if !(this_.GetPassConfigBase().options.Disables(StaticGet_DiagnosticGroups_CHECK_USELESS_CODE())) {
			callbacks.Add_1(NewCheckUnreachableCode(compiler))
		}
		if !(this_.GetPassConfigBase().options.Disables(StaticGet_DiagnosticGroups_MISSING_RETURN())) {
			callbacks.Add_1(NewCheckMissingReturn(compiler))
		}
		return DefaultPassConfig_CombineChecks(compiler, callbacks)
	})).GetPassFactory_BuilderBase().Build()
	this_.checkAccessControls = PassFactory_S_Builder().SetName("checkAccessControls").SetInternalFactory(lambdaSingleton_default_pass_config_30).GetPassFactory_BuilderBase().Build()
	this_.lintChecks = PassFactory_S_Builder().SetName("lintChecks").SetInternalFactory(lambdaSingleton_default_pass_config_31).GetPassFactory_BuilderBase().Build()
	this_.analyzerChecks = PassFactory_S_Builder().SetName("analyzerChecks").SetInternalFactory(lambdaSingleton_default_pass_config_32).GetPassFactory_BuilderBase().Build()
	this_.checkRequiresAndProvidesSorted = PassFactory_S_Builder().SetName("checkRequiresAndProvidesSorted").SetInternalFactory(lambdaSingleton_default_pass_config_33).GetPassFactory_BuilderBase().Build()
	this_.checkStrictMode = PassFactory_S_Builder().SetName("checkStrictMode").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_43 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_43)
		_ = compiler
		var defaultLevel *CheckLevel
		_ = defaultLevel
		if this_.GetPassConfigBase().options.ExpectStrictModeInput() {
			defaultLevel = StaticGet_CheckLevel_ERROR()
		} else {
			defaultLevel = StaticGet_CheckLevel_OFF()
		}
		return NewStrictModeCheck(compiler, defaultLevel)
	})).GetPassFactory_BuilderBase().Build()
	this_.processTweaks = PassFactory_S_Builder().SetName("processTweaks").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_44 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_44)
		_ = compiler
		return NewDefaultPassConfig__Anon_27(this_, compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.processDefinesCheck = this_.priv_DefaultPassConfig_createProcessDefines_2400(StaticGet_ProcessDefines_Mode_CHECK())
	this_.processDefinesOptimize = this_.priv_DefaultPassConfig_createProcessDefines_2400(StaticGet_ProcessDefines_Mode_OPTIMIZE())
	this_.stripCode = PassFactory_S_Builder().SetName("stripCode").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_45 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_45)
		_ = compiler
		return NewDefaultPassConfig__Anon_28(this_, compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.checkConsts = PassFactory_S_Builder().SetName("checkConsts").SetInternalFactory(lambdaSingleton_default_pass_config_34).GetPassFactory_BuilderBase().Build()
	this_.rewriteCallerCodeLocation = PassFactory_S_Builder().SetName("rewriteCallerCodeLocation").SetInternalFactory(lambdaSingleton_default_pass_config_35).GetPassFactory_BuilderBase().Build()
	this_.replaceToggles = PassFactory_S_Builder().SetName("replaceToggles").SetInternalFactory(lambdaSingleton_default_pass_config_36).GetPassFactory_BuilderBase().Build()
	this_.replaceIdGenerators = PassFactory_S_Builder().SetName("replaceIdGenerators").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_46 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_46)
		_ = compiler
		return NewDefaultPassConfig__Anon_29(this_, compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.replaceStrings = PassFactory_S_Builder().SetName("replaceStrings").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_47 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_47)
		_ = compiler
		return NewDefaultPassConfig__Anon_30(this_, compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.closureCodeRemoval = PassFactory_S_Builder().SetName("closureCodeRemoval").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_48 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_48)
		_ = compiler
		return NewClosureCodeRemoval(compiler, this_.GetPassConfigBase().options.ShouldRemoveAbstractMethods(), this_.GetPassConfigBase().options.ShouldRemoveClosureAsserts())
	})).GetPassFactory_BuilderBase().Build()
	this_.closureOptimizePrimitives = PassFactory_S_Builder().SetName("closureOptimizePrimitives").SetInternalFactory(lambdaSingleton_default_pass_config_37).GetPassFactory_BuilderBase().Build()
	this_.rescopeGlobalSymbols = PassFactory_S_Builder().SetName("rescopeGlobalSymbols").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_50 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_50)
		_ = compiler
		return NewRescopeGlobalSymbols_4(compiler, this_.GetPassConfigBase().options.GetRenamePrefixNamespace(), this_.GetPassConfigBase().options.AssumeCrossChunkNamesForRenamePrefixNamespace(), this_.GetPassConfigBase().options.GetOptimizeLocalAccessForGlobalSymbolNamespace())
	})).GetPassFactory_BuilderBase().Build()
	this_.convertChunksToESModules = PassFactory_S_Builder().SetName("convertChunksToESModules").SetInternalFactory(lambdaSingleton_default_pass_config_38).GetPassFactory_BuilderBase().Build()
	this_.inlineAndCollapseProperties = PassFactory_S_Builder().SetName("inlineAndCollapseProperties").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_52 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_52)
		_ = compiler
		return j2gort.PtrToAny(InlineAndCollapseProperties_S_Builder(compiler).SetPropertyCollapseLevel(this_.GetPassConfigBase().options.GetPropertyCollapseLevel()).SetChunkOutputType(this_.GetPassConfigBase().options.GetChunkOutputType()).SetHaveModulesBeenRewritten(this_.GetPassConfigBase().options.GetProcessCommonJSModules()).SetModuleResolutionMode(this_.GetPassConfigBase().options.GetModuleResolutionMode()).Build())
	})).SetPreconditionCheck(lambdaSingleton_default_pass_config_39).GetPassFactory_BuilderBase().Build()
	this_.collapseObjectLiterals = PassFactory_S_Builder().SetName("collapseObjectLiterals").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_default_pass_config_40).GetPassFactory_BuilderBase().Build()
	this_.disambiguateProperties = PassFactory_S_Builder().SetName("disambiguateProperties").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_54 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_54)
		_ = compiler
		return NewDisambiguateProperties(compiler, this_.GetPassConfigBase().options.GetPropertiesThatMustDisambiguate())
	})).SetPreconditionCheck(lambdaSingleton_default_pass_config_41).GetPassFactory_BuilderBase().Build()
	this_.devirtualizeMethods = PassFactory_S_Builder().SetName("devirtualizeMethods").SetInternalFactory(lambdaSingleton_default_pass_config_43).SetPreconditionCheck(lambdaSingleton_default_pass_config_42).GetPassFactory_BuilderBase().Build()
	this_.optimizeCalls = PassFactory_S_Builder().SetName("optimizeCalls").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_default_pass_config_45).SetPreconditionCheck(lambdaSingleton_default_pass_config_44).GetPassFactory_BuilderBase().Build()
	this_.optimizeConstructors = PassFactory_S_Builder().SetName("optimizeConstructors").SetRunInFixedPointLoop(false).SetInternalFactory(lambdaSingleton_default_pass_config_46).GetPassFactory_BuilderBase().Build()
	this_.markPureFunctions = PassFactory_S_Builder().SetName("markPureFunctions").SetInternalFactory(lambdaSingleton_default_pass_config_48).SetPreconditionCheck(lambdaSingleton_default_pass_config_47).GetPassFactory_BuilderBase().Build()
	this_.inlineVariables = PassFactory_S_Builder().SetName("inlineVariables").SetRunInFixedPointLoop(true).SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_58 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_58)
		_ = compiler
		var mode *InlineVariables_Mode = nil
		_ = mode
		if this_.GetPassConfigBase().options.ShouldInlineVariables() {
			mode = StaticGet_InlineVariables_Mode_ALL()
		} else if this_.GetPassConfigBase().options.ShouldInlineLocalVariables() {
			mode = StaticGet_InlineVariables_Mode_LOCALS_ONLY()
		} else {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "No variable inlining option set.")))
		}
		return NewInlineVariables(compiler, mode)
	})).GetPassFactory_BuilderBase().Build()
	this_.inlineConstants = PassFactory_S_Builder().SetName("inlineConstants").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_default_pass_config_49).GetPassFactory_BuilderBase().Build()
	this_.inlineSimpleMethods = PassFactory_S_Builder().SetName("inlineSimpleMethods").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_default_pass_config_51).SetPreconditionCheck(lambdaSingleton_default_pass_config_50).GetPassFactory_BuilderBase().Build()
	this_.deadAssignmentsElimination = PassFactory_S_Builder().SetName("deadAssignmentsElimination").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_default_pass_config_52).GetPassFactory_BuilderBase().Build()
	this_.deadPropertyAssignmentElimination = PassFactory_S_Builder().SetName("deadPropertyAssignmentElimination").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_default_pass_config_54).SetPreconditionCheck(lambdaSingleton_default_pass_config_53).GetPassFactory_BuilderBase().Build()
	this_.inlineFunctions = PassFactory_S_Builder().SetName("inlineFunctions").SetRunInFixedPointLoop(true).SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_60 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_60)
		_ = compiler
		return NewInlineFunctions(compiler, j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](compiler.GetUniqueNameIdSupplier()), this_.GetPassConfigBase().options.GetInlineFunctionsLevel(), (this_.GetPassConfigBase().options.AssumeStrictThis() || this_.GetPassConfigBase().options.ExpectStrictModeInput()), this_.GetPassConfigBase().options.GetMaxFunctionSizeAfterInlining())
	})).GetPassFactory_BuilderBase().Build()
	this_.inlineProperties = PassFactory_S_Builder().SetName("inlineProperties").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_default_pass_config_55).GetPassFactory_BuilderBase().Build()
	this_.isolatePolyfills = PassFactory_S_Builder().SetName("IsolatePolyfills").SetInternalFactory(lambdaSingleton_default_pass_config_56).GetPassFactory_BuilderBase().Build()
	this_.removeUnusedCode = PassFactory_S_Builder().SetName("removeUnusedCode").SetRunInFixedPointLoop(true).SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_61 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_61)
		_ = compiler
		return j2gort.PtrToAny(NewRemoveUnusedCode_Builder(compiler).RemoveLocalVars(this_.GetPassConfigBase().options.ShouldRemoveUnusedLocalVariables()).RemoveGlobals(this_.GetPassConfigBase().options.ShouldRemoveUnusedVariables()).PreserveFunctionExpressionNames(false).RemoveUnusedPrototypeProperties(this_.GetPassConfigBase().options.ShouldRemoveUnusedPrototypeProperties()).RemoveUnusedThisProperties(this_.GetPassConfigBase().options.IsRemoveUnusedClassProperties()).RemoveUnusedObjectDefinePropertiesDefinitions(this_.GetPassConfigBase().options.IsRemoveUnusedClassProperties()).RemoveUnusedPolyfills((this_.GetPassConfigBase().options.GetForceLibraryInjectionList().IsEmpty() && (this_.GetPassConfigBase().options.GetInjectPolyfillsNewerThan() == 0))).AssumeGettersArePure(this_.GetPassConfigBase().options.GetAssumeGettersArePure()).DisambiguateProperties(this_.GetPassConfigBase().options.ShouldDisambiguateProperties()).Build())
	})).GetPassFactory_BuilderBase().Build()
	this_.removeUnusedCodeOnce = this_.removeUnusedCode.ToBuilder().SetRunInFixedPointLoop(false).GetPassFactory_BuilderBase().Build()
	this_.crossChunkCodeMotion = PassFactory_S_Builder().SetName("crossChunkCodeMotion").SetRunInFixedPointLoop(true).SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_62 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_62)
		_ = compiler
		return NewCrossChunkCodeMotion(compiler, compiler.GetChunkGraph(), this_.GetPassConfigBase().options.GetParentChunkCanSeeSymbolsDeclaredInChildren())
	})).GetPassFactory_BuilderBase().Build()
	this_.crossChunkMethodMotion = PassFactory_S_Builder().SetName("crossChunkMethodMotion").SetRunInFixedPointLoop(true).SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_63 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_63)
		_ = compiler
		return NewCrossChunkMethodMotion(compiler, compiler.GetCrossChunkIdGenerator(), false, this_.GetPassConfigBase().options.GetCrossChunkCodeMotionNoStubMethods())
	})).GetPassFactory_BuilderBase().Build()
	this_.flowSensitiveInlineVariables = PassFactory_S_Builder().SetName("flowSensitiveInlineVariables").SetInternalFactory(lambdaSingleton_default_pass_config_57).GetPassFactory_BuilderBase().Build()
	this_.coalesceVariableNames = PassFactory_S_Builder().SetName("coalesceVariableNames").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_64 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_64)
		_ = compiler
		return NewCoalesceVariableNames(compiler, this_.GetPassConfigBase().options.ShouldGeneratePseudoNames())
	})).GetPassFactory_BuilderBase().Build()
	this_.postNormalizePeephole = PassFactory_S_Builder().SetName("postNormalizePeephole").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_65 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_65)
		_ = compiler
		var optimizations *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
		_ = optimizations
		if this_.GetPassConfigBase().options.ShouldCollapseVariableDeclarations() {
			optimizations.AddVoid_(NewExploitAssigns())
		}
		if this_.GetPassConfigBase().options.ShouldOptimizeLetAndConst() {
			optimizations.AddVoid_(NewOptimizeLetAndConstPeephole(this_.GetPassConfigBase().options.ShouldTreatGlobalScopeAsIsolated()))
		}
		return NewPeepholeOptimizationsPass_3_AbstractCompiler_string_Ext_java_util_List(compiler, "postNormalizePeephole", j2gort.AsIface[j2gort.Ext_java_util_List](optimizations.BuildImmList_()))
	})).GetPassFactory_BuilderBase().Build()
	this_.collapseVariableDeclarations = PassFactory_S_Builder().SetName("collapseVariableDeclarations").SetInternalFactory(lambdaSingleton_default_pass_config_58).GetPassFactory_BuilderBase().Build()
	this_.extractPrototypeMemberDeclarations = PassFactory_S_Builder().SetName("extractPrototypeMemberDeclarations").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_66 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_66)
		_ = compiler
		var pattern *ExtractPrototypeMemberDeclarations_Pattern = nil
		_ = pattern
		switch (this_.GetPassConfigBase().options.GetExtractPrototypeMemberDeclarationsMode()).Ordinal() {
		case int32(1) /* ExtractPrototypeMemberDeclarationsMode.USE_GLOBAL_TEMP */ :
			pattern = StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP()
		case int32(2) /* ExtractPrototypeMemberDeclarationsMode.USE_CHUNK_TEMP */ :
			pattern = StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP()
		case int32(3) /* ExtractPrototypeMemberDeclarationsMode.USE_IIFE */ :
			pattern = StaticGet_ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE()
		default:
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected")))
		}
		return NewExtractPrototypeMemberDeclarations(compiler, pattern)
	})).SetPreconditionCheck(lambdaSingleton_default_pass_config_59).GetPassFactory_BuilderBase().Build()
	this_.rewriteFunctionExpressions = PassFactory_S_Builder().SetName("rewriteFunctionExpressions").SetInternalFactory(lambdaSingleton_default_pass_config_60).GetPassFactory_BuilderBase().Build()
	this_.collapseAnonymousFunctions = PassFactory_S_Builder().SetName("collapseAnonymousFunctions").SetInternalFactory(lambdaSingleton_default_pass_config_61).GetPassFactory_BuilderBase().Build()
	this_.rewriteGlobalDeclarationsForTryCatchWrapping = PassFactory_S_Builder().SetName("rewriteGlobalDeclarationsForTryCatchWrapping").SetInternalFactory(lambdaSingleton_default_pass_config_62).GetPassFactory_BuilderBase().Build()
	this_.aliasStrings = PassFactory_S_Builder().SetName("aliasStrings").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_67 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_67)
		_ = compiler
		return NewAliasStrings(compiler, compiler.GetChunkGraph(), this_.GetPassConfigBase().options.ShouldOutputJsStringUsage(), this_.GetPassConfigBase().options.GetAliasStringsMode())
	})).GetPassFactory_BuilderBase().Build()
	this_.ambiguateProperties = PassFactory_S_Builder().SetName("ambiguateProperties").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_68 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_68)
		_ = compiler
		return NewAmbiguateProperties(compiler, j2gort.AsIface[j2gort.Ext_java_util_Set](this_.GetPassConfigBase().options.GetPropertyReservedNamingFirstChars()), j2gort.AsIface[j2gort.Ext_java_util_Set](this_.GetPassConfigBase().options.GetPropertyReservedNamingNonFirstChars()), j2gort.AsIface[j2gort.Ext_java_util_Set](compiler.GetExternProperties()))
	})).SetPreconditionCheck(lambdaSingleton_default_pass_config_63).GetPassFactory_BuilderBase().Build()
	this_.markUnnormalized = PassFactory_S_Builder().SetName("markUnnormalized").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_69 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_69)
		_ = compiler
		return NewDefaultPassConfig__Anon_31(this_, compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.normalize = PassFactory_S_Builder().SetName("normalize").SetInternalFactory(lambdaSingleton_default_pass_config_64).GetPassFactory_BuilderBase().Build()
	this_.externExports = PassFactory_S_Builder().SetName("externExports").SetInternalFactory(lambdaSingleton_default_pass_config_65).GetPassFactory_BuilderBase().Build()
	this_.denormalize = PassFactory_S_Builder().SetName("denormalize").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_70 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_70)
		_ = compiler
		return NewDenormalize(compiler, this_.GetPassConfigBase().options.GetOutputFeatureSet())
	})).GetPassFactory_BuilderBase().Build()
	this_.invertContextualRenaming = PassFactory_S_Builder().SetName("invertContextualRenaming").SetInternalFactory(lambdaSingleton_default_pass_config_66).GetPassFactory_BuilderBase().Build()
	this_.renameProperties = PassFactory_S_Builder().SetName("renameProperties").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_71 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_71)
		_ = compiler
		if !(this_.GetPassConfigBase().options.GetPropertyRenaming() == StaticGet_PropertyRenamingPolicy_ALL_UNQUOTED()) {
			j2gort.CheckState(false)
		}
		var prevPropertyMap *VariableMap = this_.GetPassConfigBase().options.GetInputPropertyMap()
		_ = prevPropertyMap
		return NewDefaultPassConfig__Anon_32(this_, compiler, prevPropertyMap)
	})).SetPreconditionCheck(lambdaSingleton_default_pass_config_67).GetPassFactory_BuilderBase().Build()
	this_.removePropertyRenamingCalls = PassFactory_S_Builder().SetName("removePropertyRenamingCalls").SetInternalFactory(lambdaSingleton_default_pass_config_68).GetPassFactory_BuilderBase().Build()
	this_.renameVars = PassFactory_S_Builder().SetName("renameVars").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_72 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_72)
		_ = compiler
		var prevVariableMap *VariableMap = this_.GetPassConfigBase().options.GetInputVariableMap()
		_ = prevVariableMap
		return NewDefaultPassConfig__Anon_33(this_, compiler, prevVariableMap)
	})).GetPassFactory_BuilderBase().Build()
	this_.renameLabels = PassFactory_S_Builder().SetName("renameLabels").SetInternalFactory(lambdaSingleton_default_pass_config_69).GetPassFactory_BuilderBase().Build()
	this_.convertToDottedProperties = PassFactory_S_Builder().SetName("convertToDottedProperties").SetInternalFactory(lambdaSingleton_default_pass_config_70).GetPassFactory_BuilderBase().Build()
	this_.checkAstValidity = PassFactory_S_Builder().SetName("checkAstValidity").SetInternalFactory(lambdaSingleton_default_pass_config_71).GetPassFactory_BuilderBase().Build()
	this_.varCheckValidity = PassFactory_S_Builder().SetName("varCheckValidity").SetInternalFactory(lambdaSingleton_default_pass_config_72).GetPassFactory_BuilderBase().Build()
	this_.instrumentForCodeCoverage = PassFactory_S_Builder().SetName("instrumentForCodeCoverage").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_74 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_74)
		_ = compiler
		return NewCoverageInstrumentationPass_4(compiler, StaticGet_CoverageInstrumentationPass_CoverageReach_CONDITIONAL(), this_.GetPassConfigBase().options.GetInstrumentForCoverageOption(), this_.GetPassConfigBase().options.GetProductionInstrumentationArrayName())
	})).GetPassFactory_BuilderBase().Build()
	this_.gatherExternPropertiesCheck = this_.priv_DefaultPassConfig_createGatherExternProperties_2402(StaticGet_GatherExternProperties_Mode_CHECK())
	this_.gatherExternPropertiesOptimize = this_.priv_DefaultPassConfig_createGatherExternProperties_2402(StaticGet_GatherExternProperties_Mode_OPTIMIZE())
	this_.polymerPass = PassFactory_S_Builder().SetName("polymerPass").SetInternalFactory(lambdaSingleton_default_pass_config_73).GetPassFactory_BuilderBase().Build()
	this_.chromePass = PassFactory_S_Builder().SetName("chromePass").SetInternalFactory(lambdaSingleton_default_pass_config_74).GetPassFactory_BuilderBase().Build()
	this_.j2clConstantHoisterPass = PassFactory_S_Builder().SetName("j2clConstantHoisterPass").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_default_pass_config_75).GetPassFactory_BuilderBase().Build()
	this_.j2clClinitPass = PassFactory_S_Builder().SetName("j2clClinitPass").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_default_pass_config_76).GetPassFactory_BuilderBase().Build()
	this_.j2clPropertyInlinerPass = PassFactory_S_Builder().SetName("j2clPropertyInlinerPass").SetInternalFactory(lambdaSingleton_default_pass_config_77).GetPassFactory_BuilderBase().Build()
	this_.j2clPass = PassFactory_S_Builder().SetName("j2clPass").SetInternalFactory(lambdaSingleton_default_pass_config_78).GetPassFactory_BuilderBase().Build()
	this_.j2clUtilGetDefineRewriterPass = PassFactory_S_Builder().SetName("j2clUtilGetDefineRewriterPass").SetInternalFactory(lambdaSingleton_default_pass_config_79).GetPassFactory_BuilderBase().Build()
	this_.j2clAssertRemovalPass = PassFactory_S_Builder().SetName("j2clAssertRemovalPass").SetInternalFactory(lambdaSingleton_default_pass_config_80).GetPassFactory_BuilderBase().Build()
	this_.j2clSourceFileChecker = PassFactory_S_Builder().SetName("j2clSourceFileChecker").SetInternalFactory(lambdaSingleton_default_pass_config_81).GetPassFactory_BuilderBase().Build()
	this_.j2clChecksPass = PassFactory_S_Builder().SetName("j2clChecksPass").SetInternalFactory(lambdaSingleton_default_pass_config_82).GetPassFactory_BuilderBase().Build()
	this_.checkConformance = PassFactory_S_Builder().SetName("checkConformance").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_76 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_76)
		_ = compiler
		return NewCheckConformance_3(compiler, this_.GetPassConfigBase().options.GetConformanceConfigs(), this_.GetPassConfigBase().options.GetConformanceReportingMode())
	})).GetPassFactory_BuilderBase().Build()
	this_.removeCastNodes = PassFactory_S_Builder().SetName("removeCastNodes").SetInternalFactory(lambdaSingleton_default_pass_config_83).GetPassFactory_BuilderBase().Build()
	this_.typesToColors = PassFactory_S_Builder().SetName("typesToColors").SetInternalFactory(lambdaSingleton_default_pass_config_84).GetPassFactory_BuilderBase().Build()
	this_.serializeTypedAst = PassFactory_S_Builder().SetName("serializeTypedAst").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_79 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_79)
		_ = compiler
		return j2gort.PtrToAny(SerializeTypedAstPass_CreateFromPath(compiler, this_.GetPassConfigBase().options.GetTypedAstOutputFile(), SerializationOptions_S_Builder().SetIncludeDebugInfo(compiler.GetOptions().ShouldSerializeExtraDebugInfo()).SetRuntimeLibraries(compiler.GetRuntimeJsLibManager().GetInjectedLibraries()).Build()))
	})).GetPassFactory_BuilderBase().Build()
	this_.removeUnnecessarySyntheticExterns = PassFactory_S_Builder().SetName("removeUnnecessarySyntheticExterns").SetInternalFactory(lambdaSingleton_default_pass_config_85).GetPassFactory_BuilderBase().Build()
	this_.optimizeToEs6 = PassFactory_S_Builder().SetName("optimizeToEs6").SetInternalFactory(lambdaSingleton_default_pass_config_86).GetPassFactory_BuilderBase().Build()
	this_.whitespaceWrapGoogModules = PassFactory_S_Builder().SetName("whitespaceWrapGoogModules").SetInternalFactory(lambdaSingleton_default_pass_config_87).GetPassFactory_BuilderBase().Build()
	this_.rewriteCommonJsModules = PassFactory_S_Builder().SetName("rewriteCommonJsModules").SetInternalFactory(lambdaSingleton_default_pass_config_88).GetPassFactory_BuilderBase().Build()
	this_.rewriteScriptsToEs6Modules = PassFactory_S_Builder().SetName("rewriteScriptsToEs6Modules").SetInternalFactory(lambdaSingleton_default_pass_config_89).GetPassFactory_BuilderBase().Build()
	this_.gatherModuleMetadataPass = PassFactory_S_Builder().SetName("gatherModuleMetadata").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_80 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_80)
		_ = compiler
		compiler.GetSynthesizedExternsInput()
		return NewGatherModuleMetadata(compiler, this_.GetPassConfigBase().options.GetProcessCommonJSModules(), this_.GetPassConfigBase().options.GetModuleResolutionMode())
	})).GetPassFactory_BuilderBase().Build()
	this_.createModuleMapPass = PassFactory_S_Builder().SetName("createModuleMap").SetInternalFactory(lambdaSingleton_default_pass_config_90).GetPassFactory_BuilderBase().Build()
	this_.gatherGettersAndSetters = PassFactory_S_Builder().SetName("gatherGettersAndSetters").SetInternalFactory(lambdaSingleton_default_pass_config_91).GetPassFactory_BuilderBase().Build()
	this_.addSyntheticScript = PassFactory_S_Builder().SetName("ADD_SYNTHETIC_SCRIPT").SetInternalFactory(lambdaSingleton_default_pass_config_92).GetPassFactory_BuilderBase().Build()
	this_.removeSyntheticScript = PassFactory_S_Builder().SetName("REMOVE_SYNTHETIC_SCRIPT").SetInternalFactory(lambdaSingleton_default_pass_config_93).GetPassFactory_BuilderBase().Build()
	this_.mergeSyntheticScript = PassFactory_S_Builder().SetName("MERGE_SYNTHETIC_SCRIPT").SetInternalFactory(lambdaSingleton_default_pass_config_94).GetPassFactory_BuilderBase().Build()
	this_.protectLocaleData = PassFactory_S_Builder().SetName("protectLocaleData").SetInternalFactory(lambdaSingleton_default_pass_config_95).GetPassFactory_BuilderBase().Build()
	this_.substituteLocaleData = PassFactory_S_Builder().SetName("SubstituteLocaleData").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_86 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_86)
		_ = compiler
		return NewDefaultPassConfig__Anon_35(this_, compiler)
	})).GetPassFactory_BuilderBase().Build()
	this_.transpileOnlyClosureUnaware = PassFactory_S_Builder().SetName("TranspileOnlyClosureUnaware").SetInternalFactory(lambdaSingleton_default_pass_config_96).GetPassFactory_BuilderBase().Build()
	this_.transpileAndOptimizeClosureUnaware = PassFactory_S_Builder().SetName("TranspileAndOptimizeClosureUnaware").SetInternalFactory(lambdaSingleton_default_pass_config_97).GetPassFactory_BuilderBase(

	// Java constructor: DefaultPassConfig(com.google.javascript.jscomp.CompilerOptions options)
	).Build()
}

func NewDefaultPassConfig(options *CompilerOptions) *DefaultPassConfig {
	clinit_DefaultPassConfig()
	this_ := &DefaultPassConfig{}
	this_.BindSelf_(this_)
	InitDefaultPassConfig(this_, options)
	return this_
}

func InitDefaultPassConfig(this_ *DefaultPassConfig, options *CompilerOptions) {
	_ = options
	InitPassConfig(&this_.PassConfigBase, options)
	this_.initFields_DefaultPassConfig()
}

// Java method: DefaultPassConfig.getPreprocessorSymbolTable()
func (this_ *DefaultPassConfig) GetPreprocessorSymbolTable() *PreprocessorSymbolTable {
	return this_.preprocessorSymbolTableFactory.GetInstanceOrNull()
}

// Java method: DefaultPassConfig.getTranspileOnlyPasses()
func (this_ *DefaultPassConfig) GetTranspileOnlyPasses() *PassListBuilder {
	var val_89 PassListBuilder
	InitPassListBuilder(&val_89, this_.GetPassConfigBase().options)
	var passes *PassListBuilder = &val_89
	_ = passes
	passes.MaybeAdd(this_.checkVariableReferences)
	passes.MaybeAdd(this_.checkVars)
	passes.MaybeAdd(this_.gatherModuleMetadataPass)
	passes.MaybeAdd(this_.createModuleMapPass)
	if this_.GetPassConfigBase().options.GetLanguageIn().ToFeatureSet().Has(pkg_parser.StaticGet_FeatureSet_Feature_MODULES()) {
		passes.MaybeAdd(this_.rewriteGoogJsImports)
		switch (this_.GetPassConfigBase().options.GetEs6ModuleTranspilation()).Ordinal() {
		case int32(3) /* Es6ModuleTranspilation.COMPILE */ :
			TranspilationPasses_AddEs6ModulePass(passes, this_.preprocessorSymbolTableFactory)
		case int32(2) /* Es6ModuleTranspilation.TO_COMMON_JS_LIKE_MODULES */ :
			TranspilationPasses_AddEs6ModuleToCjsPass(passes)
		case int32(1) /* Es6ModuleTranspilation.RELATIVIZE_IMPORT_PATHS */ :
			TranspilationPasses_AddEs6RewriteImportPathPass(passes)
		case int32(0) /* Es6ModuleTranspilation.NONE */ :
		}
	}
	passes.MaybeAdd(this_.checkSuper)
	TranspilationPasses_AddTranspilationRuntimeLibraries(passes)
	passes.MaybeAdd(this_.transpileOnlyClosureUnaware)
	if this_.GetPassConfigBase().options.NeedsTranspilationFrom(pkg_parser.StaticGet_FeatureSet_ES2015()) && this_.GetPassConfigBase().options.GetRewritePolyfills() {
		if this_.GetPassConfigBase().options.GetIsolatePolyfills() {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Polyfill isolation cannot be used in transpileOnly mode")))
		}
		TranspilationPasses_AddRewritePolyfillPass(passes)
	} else if this_.GetPassConfigBase().options.GetInjectPolyfillsNewerThan() != 0 {
		TranspilationPasses_AddRewritePolyfillPass(passes)
	}
	passes.MaybeAdd(this_.injectRuntimeLibraries)
	passes.MaybeAdd(this_.normalize)
	passes.MaybeAdd(this_.gatherGettersAndSetters)
	TranspilationPasses_AddTranspilationPasses(passes, this_.GetPassConfigBase().options)
	passes.MaybeAdd(this_.markUnnormalized)
	passes.MaybeAdd(this_.invertContextualRenaming)
	passes.MaybeAdd(this_.removePropertyRenamingCalls)
	passes.AssertAllOneTimePasses()
	this_.priv_DefaultPassConfig_assertValidOrderForChecks_2395(passes)
	return passes
}

// Java method: DefaultPassConfig.getWhitespaceOnlyPasses()
func (this_ *DefaultPassConfig) GetWhitespaceOnlyPasses() *PassListBuilder {
	var val_90 PassListBuilder
	InitPassListBuilder(&val_90, this_.GetPassConfigBase().options)
	var passes *PassListBuilder = &val_90
	_ = passes
	if this_.GetPassConfigBase().options.GetProcessCommonJSModules() {
		passes.MaybeAdd(this_.rewriteCommonJsModules)
	} else if this_.GetPassConfigBase().options.GetLanguageIn().ToFeatureSet().Has(pkg_parser.StaticGet_FeatureSet_Feature_MODULES()) {
		passes.MaybeAdd(this_.rewriteScriptsToEs6Modules)
	}
	if this_.GetPassConfigBase().options.ShouldWrapGoogModulesForWhitespaceOnly() {
		passes.MaybeAdd(this_.whitespaceWrapGoogModules)
	}
	return passes
}

// Java method: DefaultPassConfig.addModuleRewritingPasses(com.google.javascript.jscomp.PassListBuilder checks, com.google.javascript.jscomp.CompilerOptions options)
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_addModuleRewritingPasses_2386(checks *PassListBuilder, options *CompilerOptions) {
	_ = checks
	_ = options
	if options.GetLanguageIn().ToFeatureSet().Has(pkg_parser.StaticGet_FeatureSet_Feature_MODULES()) {
		checks.MaybeAdd(this_.rewriteGoogJsImports)
		TranspilationPasses_AddEs6ModulePass(checks, this_.preprocessorSymbolTableFactory)
	}
	if options.GetClosurePass() {
		checks.MaybeAdd(this_.closureRewriteModule)
	}
}

// Java method: DefaultPassConfig.getChecks()
func (this_ *DefaultPassConfig) GetChecks() *PassListBuilder {
	var val_91 PassListBuilder
	InitPassListBuilder(&val_91, this_.GetPassConfigBase().options)
	var checks *PassListBuilder = &val_91
	_ = checks
	if !(!(this_.GetPassConfigBase().options.GetSkipNonTranspilationPasses())) {
		j2gort.CheckStateStr(false, "options.getSkipNonTranspilationPasses() cannot be mixed with PassConfig::getChecks. Call PassConfig::getTranspileOnlyPasses instead.")
	}
	if this_.GetPassConfigBase().options.IsPropertyRenamingOnlyCompilationMode() {
		checks.MaybeAdd(this_.addSyntheticScript)
		checks.MaybeAdd(this_.gatherGettersAndSetters)
		checks.MaybeAdd(this_.gatherModuleMetadataPass)
		checks.MaybeAdd(this_.createModuleMapPass)
		checks.MaybeAdd(this_.declaredGlobalExternsOnWindow)
		checks.MaybeAdd(this_.checkSideEffects)
		checks.MaybeAdd(this_.angularPass)
		checks.MaybeAdd(this_.closureGoogScopeAliases)
		checks.MaybeAdd(this_.closurePrimitives)
		this_.priv_DefaultPassConfig_addModuleRewritingPasses_2386(checks, this_.GetPassConfigBase().options)
		checks.MaybeAdd(this_.clearTypedScopeCreatorPass)
		checks.MaybeAdd(this_.clearTopTypedScopePass)
		checks.MaybeAdd(this_.generateExports)
		checks.MaybeAdd(PassFactory_CreateEmptyPass("afterStandardChecks"))
		checks.MaybeAdd(this_.mergeSyntheticScript)
		checks.MaybeAdd(this_.gatherExternPropertiesCheck)
		checks.MaybeAdd(PassFactory_CreateEmptyPass("beforeSerialization"))
		return checks
	}
	if this_.GetPassConfigBase().options.ShouldGenerateTypedExterns() {
		checks.MaybeAdd(this_.addSyntheticScript)
		checks.MaybeAdd(this_.closureGoogScopeAliasesForIjs)
		checks.MaybeAdd(this_.generateIjs)
		checks.MaybeAdd(this_.removeExtraRequires)
		if this_.GetPassConfigBase().options.ShouldWrapGoogModulesForWhitespaceOnly() {
			checks.MaybeAdd(this_.whitespaceWrapGoogModules)
		}
		checks.MaybeAdd(this_.removeSyntheticScript)
		return checks
	}
	checks.MaybeAdd(this_.addSyntheticScript)
	checks.MaybeAdd(this_.gatherGettersAndSetters)
	if this_.GetPassConfigBase().options.GetLanguageIn().ToFeatureSet().Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT()) && !(this_.GetPassConfigBase().options.ShouldAllowDynamicImport()) {
		checks.MaybeAdd(StaticGet_DefaultPassConfig_ForbidDynamicImportUsage())
	}
	checks.MaybeAdd(PassFactory_CreateEmptyPass("beforeStandardChecks"))
	if !(this_.GetPassConfigBase().options.GetProcessCommonJSModules()) && this_.GetPassConfigBase().options.GetLanguageIn().ToFeatureSet().Has(pkg_parser.StaticGet_FeatureSet_Feature_MODULES()) {
		checks.MaybeAdd(this_.rewriteScriptsToEs6Modules)
	}
	checks.MaybeAdd(this_.gatherModuleMetadataPass)
	checks.MaybeAdd(this_.createModuleMapPass)
	if this_.GetPassConfigBase().options.GetProcessCommonJSModules() {
		checks.MaybeAdd(this_.rewriteCommonJsModules)
	}
	if this_.GetPassConfigBase().options.IsChromePassEnabled() {
		checks.MaybeAdd(this_.chromePass)
	}
	checks.MaybeAdd(this_.checkJsDocAndEs6Modules)
	checks.MaybeAdd(this_.checkTypeImportCodeReferences)
	if this_.GetPassConfigBase().options.Enables(StaticGet_DiagnosticGroups_LINT_CHECKS()) {
		checks.MaybeAdd(this_.lintChecks)
	}
	if this_.GetPassConfigBase().options.GetClosurePass() && this_.GetPassConfigBase().options.Enables(StaticGet_DiagnosticGroups_LINT_CHECKS()) {
		checks.MaybeAdd(this_.checkRequiresAndProvidesSorted)
	}
	if this_.GetPassConfigBase().options.Enables(StaticGet_DiagnosticGroups_EXTRA_REQUIRE()) {
		checks.MaybeAdd(this_.extraRequires)
	}
	if this_.GetPassConfigBase().options.Enables(StaticGet_DiagnosticGroups_MISSING_REQUIRE()) {
		checks.MaybeAdd(this_.checkMissingRequires)
	}
	checks.MaybeAdd(this_.declaredGlobalExternsOnWindow)
	if !(this_.GetPassConfigBase().options.GetProcessCommonJSModules()) {
		checks.MaybeAdd(this_.checkVariableReferences)
		checks.MaybeAdd(this_.checkVars)
	}
	if this_.GetPassConfigBase().options.GetClosurePass() {
		checks.MaybeAdd(this_.checkClosureImports)
	}
	checks.MaybeAdd(this_.checkStrictMode)
	if this_.GetPassConfigBase().options.GetClosurePass() {
		checks.MaybeAdd(this_.closureCheckModule)
	}
	checks.MaybeAdd(this_.checkSuper)
	checks.MaybeAdd(this_.checkSideEffects)
	if this_.GetPassConfigBase().options.GetAngularPass() {
		checks.MaybeAdd(this_.angularPass)
	}
	if this_.GetPassConfigBase().options.GetClosurePass() {
		checks.MaybeAdd(this_.closureGoogScopeAliases)
	}
	if this_.GetPassConfigBase().options.GetClosurePass() {
		checks.MaybeAdd(this_.closurePrimitives)
	}
	if this_.GetPassConfigBase().options.ShouldRewriteModulesBeforeTypechecking() {
		this_.priv_DefaultPassConfig_addModuleRewritingPasses_2386(checks, this_.GetPassConfigBase().options)
	}
	if this_.GetPassConfigBase().options.IsPolymerPassEnabled() {
		checks.MaybeAdd(this_.polymerPass)
	}
	if this_.GetPassConfigBase().options.GetProcessCommonJSModules() {
		checks.MaybeAdd(this_.checkVariableReferences)
		checks.MaybeAdd(this_.checkVars)
	}
	if this_.GetPassConfigBase().options.ShouldInferConsts() {
		checks.MaybeAdd(this_.inferConsts)
	}
	if this_.GetPassConfigBase().options.ShouldComputeFunctionSideEffects() {
		checks.MaybeAdd(this_.checkRegExp)
	}
	checks.MaybeAdd(PassFactory_CreateEmptyPass("beforeTypeChecking"))
	if this_.GetPassConfigBase().options.GetCheckTypes() || this_.GetPassConfigBase().options.GetInferTypes() {
		checks.MaybeAdd(this_.inferTypes)
		if this_.GetPassConfigBase().options.GetCheckTypes() {
			checks.MaybeAdd(this_.checkTypes)
		} else {
			checks.MaybeAdd(this_.inferJsDocInfo)
		}
	}
	if this_.GetPassConfigBase().options.Enables(StaticGet_DiagnosticGroups_ANALYZER_CHECKS()) && this_.GetPassConfigBase().options.IsTypecheckingEnabled() {
		checks.MaybeAdd(this_.analyzerChecks)
	}
	if !(this_.GetPassConfigBase().options.PreservesDetailedSourceInfo()) {
		checks.MaybeAdd(this_.clearTypedScopeCreatorPass)
	}
	if this_.GetPassConfigBase().options.ShouldRewriteModulesAfterTypechecking() {
		this_.priv_DefaultPassConfig_addModuleRewritingPasses_2386(checks, this_.GetPassConfigBase().options)
	}
	if this_.GetPassConfigBase().options.ShouldAllowDynamicImport() && this_.GetPassConfigBase().options.GetLanguageIn().ToFeatureSet().Has(pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT()) {
		checks.MaybeAdd(StaticGet_DefaultPassConfig_RewriteDynamicImports())
	}
	if !(this_.GetPassConfigBase().options.PreservesDetailedSourceInfo()) {
		checks.MaybeAdd(this_.clearTopTypedScopePass)
	}
	if (this_.GetPassConfigBase().options.GetCheckSuspiciousCode() || this_.GetPassConfigBase().options.Enables(StaticGet_DiagnosticGroups_GLOBAL_THIS())) || this_.GetPassConfigBase().options.Enables(StaticGet_DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT()) {
		checks.MaybeAdd(this_.suspiciousCode)
	}
	if this_.GetPassConfigBase().options.GetJ2clPass().ShouldAddJ2clPasses() {
		checks.MaybeAdd(this_.j2clSourceFileChecker)
	}
	if !(this_.GetPassConfigBase().options.Disables(StaticGet_DiagnosticGroups_CHECK_USELESS_CODE())) || !(this_.GetPassConfigBase().options.Disables(StaticGet_DiagnosticGroups_MISSING_RETURN())) {
		checks.MaybeAdd(this_.checkControlFlow)
	}
	if this_.GetPassConfigBase().options.IsTypecheckingEnabled() && (!(this_.GetPassConfigBase().options.Disables(StaticGet_DiagnosticGroups_ACCESS_CONTROLS())) || this_.GetPassConfigBase().options.Enables(StaticGet_DiagnosticGroups_CONSTANT_PROPERTY())) {
		checks.MaybeAdd(this_.checkAccessControls)
	}
	checks.MaybeAdd(this_.checkConsts)
	checks.MaybeAdd(this_.rewriteCallerCodeLocation)
	if len(j2gort.CheckNotNullPtr(this_.GetPassConfigBase().options).Slice_GetConformanceConfigs()) != 0 {
		checks.MaybeAdd(this_.checkConformance)
	}
	if this_.GetPassConfigBase().options.GetTweakProcessing().IsOn() {
		checks.MaybeAdd(this_.processTweaks)
	}
	checks.MaybeAdd(this_.processDefinesCheck)
	if this_.GetPassConfigBase().options.GetJ2clPass().ShouldAddJ2clPasses() {
		checks.MaybeAdd(this_.j2clChecksPass)
	}
	if this_.GetPassConfigBase().options.ShouldGenerateExports() {
		checks.MaybeAdd(this_.generateExports)
	}
	checks.MaybeAdd(PassFactory_CreateEmptyPass("afterStandardChecks"))
	checks.MaybeAdd(this_.mergeSyntheticScript)
	if this_.GetPassConfigBase().options.GetExternExportsPath() != j2gort.NullString {
		checks.MaybeAdd(this_.externExports)
	}
	if !(this_.GetPassConfigBase().options.IsChecksOnly()) {
		checks.MaybeAdd(this_.removeWeakSources)
	}
	checks.MaybeAdd(this_.gatherExternPropertiesCheck)
	checks.MaybeAdd(PassFactory_CreateEmptyPass("beforeSerialization"))
	if this_.GetPassConfigBase().options.GetTypedAstOutputFile() != nil {
		checks.MaybeAdd(this_.serializeTypedAst)
	}
	checks.AssertAllOneTimePasses()
	this_.priv_DefaultPassConfig_assertValidOrderForChecks_2395(checks)
	return checks
}

// Java method: DefaultPassConfig.getOptimizations()
func (this_ *DefaultPassConfig) GetOptimizations() *PassListBuilder {
	var val_92 PassListBuilder
	InitPassListBuilder(&val_92, this_.GetPassConfigBase().options)
	var passes *PassListBuilder = &val_92
	_ = passes
	if this_.GetPassConfigBase().options.IsPropertyRenamingOnlyCompilationMode() {
		passes.MaybeAdd(this_.removeUnnecessarySyntheticExterns)
		TranspilationPasses_AddTranspilationRuntimeLibraries(passes)
		passes.MaybeAdd(this_.closureProvidesRequires)
		passes.MaybeAdd(this_.processDefinesOptimize)
		passes.MaybeAdd(this_.normalize)
		passes.MaybeAdd(this_.gatherGettersAndSetters)
		TranspilationPasses_AddTranspilationPasses(passes, this_.GetPassConfigBase().options)
		passes.MaybeAdd(this_.gatherExternPropertiesOptimize)
		passes.MaybeAdd(PassFactory_CreateEmptyPass("beforeStandardOptimizations"))
		passes.MaybeAdd(this_.inlineAndCollapseProperties)
		passes.MaybeAdd(this_.closureOptimizePrimitives)
		passes.MaybeAdd(this_.stripSideEffectProtection)
		return passes
	}
	if this_.GetPassConfigBase().options.GetSkipNonTranspilationPasses() {
		return passes
	}
	passes.AddAll(this_.priv_DefaultPassConfig_getEarlyOptimizationPasses_2390())
	passes.MaybeAdd(PassFactory_CreateEmptyPass("optimizationsHalfwayPoint"))
	passes.AddAll(this_.priv_DefaultPassConfig_getLateOptimizationPasses_2391())
	return passes
}

// Java method: DefaultPassConfig.getFinalizations()
func (this_ *DefaultPassConfig) GetFinalizations() *PassListBuilder {
	var val_93 PassListBuilder
	InitPassListBuilder(&val_93, this_.GetPassConfigBase().options)
	var passes *PassListBuilder = &val_93
	_ = passes
	if this_.GetPassConfigBase().options.IsPropertyRenamingOnlyCompilationMode() {
		if this_.GetPassConfigBase().options.ShouldRewriteGlobalDeclarationsForTryCatchWrapping() || (this_.GetPassConfigBase().options.GetRenamePrefixNamespace() != j2gort.NullString) {
			passes.MaybeAdd(this_.rewriteGlobalDeclarationsForTryCatchWrapping)
		}
		passes.MaybeAdd(this_.renameProperties)
		if (this_.GetPassConfigBase().options.GetRenamePrefixNamespace() != j2gort.NullString) && (this_.GetPassConfigBase().options.GetChunkOutputType() == StaticGet_CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE()) {
			if !(StaticGet_DefaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN().Matcher(j2gort.BoxString(this_.GetPassConfigBase().options.GetRenamePrefixNamespace())).Matches()) {
				panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Illegal character in renamePrefixNamespace name: "+j2gort.StrValString(this_.GetPassConfigBase().options.GetRenamePrefixNamespace()))))
			}
			passes.MaybeAdd(this_.rescopeGlobalSymbols)
		}
		return passes
	}
	if this_.GetPassConfigBase().options.DoLateLocalization() {
		if this_.GetPassConfigBase().options.ShouldRunReplaceMessagesPass() {
			passes.MaybeAdd(this_.priv_DefaultPassConfig_getReplaceProtectedMessagesPass_2399())
		}
		passes.MaybeAdd(this_.substituteLocaleData)
		passes.MaybeAdd(this_.peepholeOptimizationsOnceNormalized)
	}
	if this_.GetPassConfigBase().options.ShouldInlineVariables() || this_.GetPassConfigBase().options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(this_.flowSensitiveInlineVariables)
		if !(this_.GetPassConfigBase().options.DoLateLocalization()) && this_.priv_DefaultPassConfig_shouldRunRemoveUnusedCode_2394() {
			passes.MaybeAdd(this_.removeUnusedCodeOnce)
		}
	}
	if this_.GetPassConfigBase().options.DoLateLocalization() {
		passes.AddAll(this_.priv_DefaultPassConfig_getPostL10nOptimizations_2389())
	}
	passes.MaybeAdd(PassFactory_CreateEmptyPass("beforeChunkMotion"))
	if this_.GetPassConfigBase().options.ShouldRunCrossChunkCodeMotion() {
		passes.MaybeAdd(this_.crossChunkCodeMotion)
	}
	if this_.GetPassConfigBase().options.ShouldRunCrossChunkMethodMotion() {
		passes.MaybeAdd(this_.crossChunkMethodMotion)
	}
	passes.MaybeAdd(PassFactory_CreateEmptyPass("afterChunkMotion"))
	if this_.GetPassConfigBase().options.GetOptimizeESClassConstructors() && this_.GetPassConfigBase().options.GetOutputFeatureSet().Contains_1_FeatureSet(pkg_parser.StaticGet_FeatureSet_ES2015()) {
		passes.MaybeAdd(this_.optimizeConstructors)
	}
	if this_.GetPassConfigBase().options.GetIsolatePolyfills() {
		passes.MaybeAdd(this_.isolatePolyfills)
	}
	if this_.GetPassConfigBase().options.ShouldCollapseAnonymousFunctions() {
		passes.MaybeAdd(this_.collapseAnonymousFunctions)
	}
	if this_.GetPassConfigBase().options.ShouldRewriteGlobalDeclarationsForTryCatchWrapping() || (this_.GetPassConfigBase().options.GetRenamePrefixNamespace() != j2gort.NullString) {
		passes.MaybeAdd(this_.rewriteGlobalDeclarationsForTryCatchWrapping)
	}
	passes.MaybeAdd(PassFactory_CreateEmptyPass("beforeExtractPrototypeMemberDeclarations"))
	if this_.GetPassConfigBase().options.GetExtractPrototypeMemberDeclarationsMode() != StaticGet_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF() {
		passes.MaybeAdd(this_.extractPrototypeMemberDeclarations)
	}
	if (this_.GetPassConfigBase().options.ShouldAmbiguateProperties() && (this_.GetPassConfigBase().options.GetPropertyRenaming() == StaticGet_PropertyRenamingPolicy_ALL_UNQUOTED())) && this_.GetPassConfigBase().options.IsTypecheckingEnabled() {
		passes.MaybeAdd(this_.ambiguateProperties)
	}
	passes.MaybeAdd(PassFactory_CreateEmptyPass("beforeRenameProperties"))
	if this_.GetPassConfigBase().options.GetPropertyRenaming() == StaticGet_PropertyRenamingPolicy_ALL_UNQUOTED() {
		passes.MaybeAdd(this_.renameProperties)
	} else {
		passes.MaybeAdd(this_.removePropertyRenamingCalls)
	}
	if this_.GetPassConfigBase().options.ShouldReserveRawExports() {
		passes.MaybeAdd(this_.gatherRawExports)
	}
	if this_.GetPassConfigBase().options.ShouldConvertToDottedProperties() {
		passes.MaybeAdd(this_.convertToDottedProperties)
	}
	if this_.GetPassConfigBase().options.ShouldRewriteFunctionExpressions() {
		passes.MaybeAdd(this_.rewriteFunctionExpressions)
	}
	if this_.GetPassConfigBase().options.GetAliasStringsMode() != StaticGet_CompilerOptions_AliasStringsMode_NONE() {
		passes.MaybeAdd(this_.aliasStrings)
	}
	if this_.GetPassConfigBase().options.ShouldCoalesceVariableNames() {
		passes.MaybeAdd(this_.coalesceVariableNames)
		if this_.GetPassConfigBase().options.ShouldFoldConstants() {
			passes.MaybeAdd(this_.peepholeOptimizationsOnceNonNormalized)
		}
	}
	passes.MaybeAdd(this_.markUnnormalized)
	if this_.GetPassConfigBase().options.ShouldCollapseVariableDeclarations() || this_.GetPassConfigBase().options.ShouldOptimizeLetAndConst() {
		passes.MaybeAdd(this_.postNormalizePeephole)
	}
	if this_.GetPassConfigBase().options.ShouldCollapseVariableDeclarations() {
		passes.MaybeAdd(this_.collapseVariableDeclarations)
	}
	passes.MaybeAdd(this_.denormalize)
	passes.MaybeAdd(PassFactory_CreateEmptyPass("beforeVariableRenaming"))
	if this_.GetPassConfigBase().options.GetVariableRenaming() != StaticGet_VariableRenamingPolicy_ALL() {
		passes.MaybeAdd(this_.invertContextualRenaming)
	}
	if this_.GetPassConfigBase().options.GetVariableRenaming() != StaticGet_VariableRenamingPolicy_OFF() {
		passes.MaybeAdd(this_.renameVars)
	}
	if this_.GetPassConfigBase().options.ShouldRenameLabels() {
		passes.MaybeAdd(this_.renameLabels)
	}
	if this_.GetPassConfigBase().options.ShouldFoldConstants() {
		passes.MaybeAdd(this_.latePeepholeOptimizations)
	}
	if this_.GetPassConfigBase().options.ShouldProtectHiddenSideEffects() || this_.GetPassConfigBase().options.GetMergedPrecompiledLibraries() {
		passes.MaybeAdd(this_.stripSideEffectProtection)
	}
	if (this_.GetPassConfigBase().options.GetRenamePrefixNamespace() != j2gort.NullString) && (this_.GetPassConfigBase().options.GetChunkOutputType() == StaticGet_CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE()) {
		if !(StaticGet_DefaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN().Matcher(j2gort.BoxString(this_.GetPassConfigBase().options.GetRenamePrefixNamespace())).Matches()) {
			panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Illegal character in renamePrefixNamespace name: "+j2gort.StrValString(this_.GetPassConfigBase().options.GetRenamePrefixNamespace()))))
		}
		passes.MaybeAdd(this_.rescopeGlobalSymbols)
	}
	if this_.GetPassConfigBase().options.GetOutputFeatureSet().Contains_1_FeatureSet(pkg_parser.StaticGet_FeatureSet_ES2015()) {
		passes.MaybeAdd(this_.optimizeToEs6)
	}
	if this_.GetPassConfigBase().options.GetChunkOutputType() == StaticGet_CompilerOptions_ChunkOutputType_ES_MODULES() {
		passes.MaybeAdd(this_.convertChunksToESModules)
	}
	passes.MaybeAdd(this_.checkAstValidity)
	passes.MaybeAdd(this_.varCheckValidity)
	return passes
}

// Java method: DefaultPassConfig.getEarlyOptimizationLoopPasses()
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_getEarlyOptimizationLoopPasses_2388() *PassListBuilder {
	var val_94 PassListBuilder
	InitPassListBuilder(&val_94, this_.GetPassConfigBase().options)
	var earlyLoopPasses *PassListBuilder = &val_94
	_ = earlyLoopPasses
	if this_.GetPassConfigBase().options.ShouldInlineVariables() || this_.GetPassConfigBase().options.ShouldInlineLocalVariables() {
		earlyLoopPasses.MaybeAdd(this_.inlineVariables)
	} else if this_.GetPassConfigBase().options.ShouldInlineConstantVars() {
		earlyLoopPasses.MaybeAdd(this_.inlineConstants)
	}
	if this_.GetPassConfigBase().options.GetCollapseObjectLiterals() {
		earlyLoopPasses.MaybeAdd(this_.collapseObjectLiterals)
	}
	if this_.priv_DefaultPassConfig_shouldRunRemoveUnusedCode_2394() {
		earlyLoopPasses.MaybeAdd(this_.removeUnusedCode)
	}
	if this_.GetPassConfigBase().options.ShouldFoldConstants() {
		earlyLoopPasses.MaybeAdd(this_.peepholeOptimizations)
	}
	earlyLoopPasses.AssertAllLoopablePasses()
	return earlyLoopPasses
}

// Java method: DefaultPassConfig.getPostL10nOptimizations()
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_getPostL10nOptimizations_2389() *PassListBuilder {
	var val_95 PassListBuilder
	InitPassListBuilder(&val_95, this_.GetPassConfigBase().options)
	var loopPasses *PassListBuilder = &val_95
	_ = loopPasses
	if this_.GetPassConfigBase().options.ShouldOptimizeCalls() {
		loopPasses.MaybeAdd(this_.optimizeCalls)
	}
	if this_.GetPassConfigBase().options.GetJ2clPass().ShouldAddJ2clPasses() {
		loopPasses.MaybeAdd(this_.j2clConstantHoisterPass)
		loopPasses.MaybeAdd(this_.j2clClinitPass)
	}
	if this_.GetPassConfigBase().options.GetInlineFunctionsLevel() != StaticGet_CompilerOptions_Reach_NONE() {
		loopPasses.MaybeAdd(this_.inlineFunctions)
	}
	if this_.GetPassConfigBase().options.ShouldInlineVariables() || this_.GetPassConfigBase().options.ShouldInlineLocalVariables() {
		loopPasses.MaybeAdd(this_.inlineVariables)
	} else if this_.GetPassConfigBase().options.ShouldInlineConstantVars() {
		loopPasses.MaybeAdd(this_.inlineConstants)
	}
	if this_.priv_DefaultPassConfig_shouldRunRemoveUnusedCode_2394() {
		loopPasses.MaybeAdd(this_.removeUnusedCode)
	}
	if this_.GetPassConfigBase().options.ShouldFoldConstants() {
		loopPasses.MaybeAdd(this_.peepholeOptimizations)
	}
	loopPasses.AssertAllLoopablePasses()
	return loopPasses
}

// Java method: DefaultPassConfig.getEarlyOptimizationPasses()
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_getEarlyOptimizationPasses_2390() *PassListBuilder {
	var val_96 PassListBuilder
	InitPassListBuilder(&val_96, this_.GetPassConfigBase().options)
	var passes *PassListBuilder = &val_96
	_ = passes
	if this_.GetPassConfigBase().options.ShouldExportTestFunctions() {
		passes.MaybeAdd(this_.exportTestFunctions)
	}
	if this_.GetPassConfigBase().options.GetMergedPrecompiledLibraries() {
		passes.MaybeAdd(this_.removeWeakSources)
		passes.MaybeAdd(this_.checkRegExpForOptimizations)
		if this_.GetPassConfigBase().options.GetJ2clPass().ShouldAddJ2clPasses() {
			passes.MaybeAdd(this_.j2clSourceFileChecker)
		}
	} else {
		this_.priv_DefaultPassConfig_addNonTypedAstNormalizationPasses_2393(passes)
	}
	passes.MaybeAdd(this_.removeUnnecessarySyntheticExterns)
	if this_.GetPassConfigBase().options.GetSyntheticBlockStartMarker() != j2gort.NullString {
		passes.MaybeAdd(this_.createSyntheticBlocks)
	}
	if this_.GetPassConfigBase().options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(this_.j2clPass)
	}
	passes.MaybeAdd(this_.transpileAndOptimizeClosureUnaware)
	TranspilationPasses_AddTranspilationRuntimeLibraries(passes)
	if (this_.GetPassConfigBase().options.GetRewritePolyfills() || this_.GetPassConfigBase().options.GetIsolatePolyfills()) || (this_.GetPassConfigBase().options.GetInjectPolyfillsNewerThan() != 0) {
		TranspilationPasses_AddRewritePolyfillPass(passes)
	}
	passes.MaybeAdd(this_.injectRuntimeLibraries)
	if this_.GetPassConfigBase().options.GetClosurePass() {
		passes.MaybeAdd(this_.closureProvidesRequires)
	}
	if this_.GetPassConfigBase().options.ShouldRunReplaceMessagesForChrome() {
		passes.MaybeAdd(this_.replaceMessagesForChrome)
	} else if this_.GetPassConfigBase().options.ShouldRunReplaceMessagesPass() {
		if this_.GetPassConfigBase().options.DoLateLocalization() {
			passes.MaybeAdd(this_.priv_DefaultPassConfig_getProtectMessagesPass_2398())
		} else {
			passes.MaybeAdd(this_.priv_DefaultPassConfig_getFullReplaceMessagesPass_2397())
		}
	}
	if this_.GetPassConfigBase().options.DoLateLocalization() {
		passes.MaybeAdd(this_.protectLocaleData)
	}
	if this_.GetPassConfigBase().options.GetClosurePass() && !(this_.GetPassConfigBase().options.ShouldPreserveGoogLibraryPrimitives()) {
		passes.MaybeAdd(this_.closureReplaceGetCssName)
	}
	passes.MaybeAdd(this_.processDefinesOptimize)
	passes.MaybeAdd(PassFactory_CreateEmptyPass("beforeEarlyOptimizationsTranspilation"))
	passes.MaybeAdd(this_.normalize)
	passes.MaybeAdd(this_.gatherGettersAndSetters)
	TranspilationPasses_AddTranspilationPasses(passes, this_.GetPassConfigBase().options)
	if this_.GetPassConfigBase().options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(this_.j2clUtilGetDefineRewriterPass)
	}
	if this_.GetPassConfigBase().options.GetInstrumentForCoverageOption() != StaticGet_CompilerOptions_InstrumentOption_NONE() {
		passes.MaybeAdd(this_.instrumentForCodeCoverage)
	}
	passes.MaybeAdd(this_.gatherExternPropertiesOptimize)
	passes.MaybeAdd(PassFactory_CreateEmptyPass("beforeStandardOptimizations"))
	if this_.GetPassConfigBase().options.GetClosurePass() && (this_.GetPassConfigBase().options.ShouldRemoveAbstractMethods() || this_.GetPassConfigBase().options.ShouldRemoveClosureAsserts()) {
		passes.MaybeAdd(this_.closureCodeRemoval)
	}
	if this_.GetPassConfigBase().options.ShouldRemoveJ2clAsserts() {
		passes.MaybeAdd(this_.j2clAssertRemovalPass)
	}
	passes.MaybeAdd(this_.replaceToggles)
	passes.MaybeAdd(this_.inlineAndCollapseProperties)
	if ((this_.GetPassConfigBase().options.GetTweakProcessing().ShouldStrip() || !(this_.GetPassConfigBase().options.GetStripTypes().IsEmpty())) || !(this_.GetPassConfigBase().options.GetStripNameSuffixes().IsEmpty())) || !(this_.GetPassConfigBase().options.GetStripNamePrefixes().IsEmpty()) {
		passes.MaybeAdd(this_.stripCode)
	}
	if this_.GetPassConfigBase().options.ShouldReplaceIdGenerators() {
		passes.MaybeAdd(this_.replaceIdGenerators)
	}
	if this_.GetPassConfigBase().options.GetJ2clPass().ShouldAddJ2clPasses() && (this_.GetPassConfigBase().options.GetPropertyCollapseLevel() == StaticGet_CompilerOptions_PropertyCollapseLevel_ALL()) {
		passes.MaybeAdd(this_.j2clPropertyInlinerPass)
	}
	if this_.GetPassConfigBase().options.ShouldInferConsts() {
		passes.MaybeAdd(this_.inferConsts)
	}
	passes.MaybeAdd(PassFactory_CreateEmptyPass("obfuscationPassMarker"))
	if this_.GetPassConfigBase().options.GetSmartNameRemoval() {
		if this_.GetPassConfigBase().options.ShouldFoldConstants() && (this_.GetPassConfigBase().options.ShouldInlineVariables() || this_.GetPassConfigBase().options.ShouldInlineLocalVariables()) {
			passes.MaybeAdd(this_.earlyInlineVariables)
			passes.MaybeAdd(this_.earlyPeepholeOptimizations)
		}
		passes.MaybeAdd(this_.removeUnusedCodeOnce)
	}
	if this_.GetPassConfigBase().options.ShouldDisambiguateProperties() && this_.GetPassConfigBase().options.IsTypecheckingEnabled() {
		passes.MaybeAdd(this_.disambiguateProperties)
	}
	if this_.GetPassConfigBase().options.ShouldComputeFunctionSideEffects() {
		passes.MaybeAdd(this_.markPureFunctions)
	}
	passes.AssertAllOneTimePasses()
	return passes
}

// Java method: DefaultPassConfig.getLateOptimizationPasses()
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_getLateOptimizationPasses_2391() *PassListBuilder {
	var val_97 PassListBuilder
	InitPassListBuilder(&val_97, this_.GetPassConfigBase().options)
	var passes *PassListBuilder = &val_97
	_ = passes
	if this_.GetPassConfigBase().options.GetSmartNameRemoval() {
		passes.MaybeAdd(PassFactory_CreateEmptyPass("beforeEarlyOptimizationLoop"))
		passes.AddAll(this_.priv_DefaultPassConfig_getEarlyOptimizationLoopPasses_2388())
		passes.MaybeAdd(PassFactory_CreateEmptyPass("afterEarlyOptimizationLoop"))
	}
	if this_.GetPassConfigBase().options.GetClosurePass() {
		passes.MaybeAdd(this_.closureOptimizePrimitives)
	}
	if !(this_.GetPassConfigBase().options.GetReplaceStringsFunctionDescriptions().IsEmpty()) {
		passes.MaybeAdd(this_.replaceStrings)
	}
	if this_.GetPassConfigBase().options.ShouldRunCrossChunkCodeMotion() {
		passes.MaybeAdd(this_.crossChunkCodeMotion)
	}
	if this_.GetPassConfigBase().options.ShouldDevirtualizeMethods() {
		passes.MaybeAdd(this_.devirtualizeMethods)
	}
	passes.MaybeAdd(this_.priv_DefaultPassConfig_getCustomPasses_2403(StaticGet_CustomPassExecutionTime_BEFORE_OPTIMIZATION_LOOP()))
	passes.MaybeAdd(PassFactory_CreateEmptyPass("beforeMainOptimizations"))
	if this_.GetPassConfigBase().options.ShouldInlineVariables() || this_.GetPassConfigBase().options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(this_.flowSensitiveInlineVariables)
	}
	passes.AddAll(this_.priv_DefaultPassConfig_getMainOptimizationLoop_2392())
	passes.MaybeAdd(PassFactory_CreateEmptyPass("afterMainOptimizations"))
	passes.MaybeAdd(this_.priv_DefaultPassConfig_getCustomPasses_2403(StaticGet_CustomPassExecutionTime_AFTER_OPTIMIZATION_LOOP()))
	this_.priv_DefaultPassConfig_assertValidOrderForOptimizations_2396(passes)
	return passes
}

// Java method: DefaultPassConfig.getMainOptimizationLoop()
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_getMainOptimizationLoop_2392() *PassListBuilder {
	var val_98 PassListBuilder
	InitPassListBuilder(&val_98, this_.GetPassConfigBase().options)
	var passes *PassListBuilder = &val_98
	_ = passes
	if this_.GetPassConfigBase().options.ShouldInlineGetters() {
		passes.MaybeAdd(this_.inlineSimpleMethods)
	}
	if this_.GetPassConfigBase().options.ShouldInlineProperties() && this_.GetPassConfigBase().options.IsTypecheckingEnabled() {
		passes.MaybeAdd(this_.inlineProperties)
	}
	if this_.GetPassConfigBase().options.ShouldRunDeadPropertyAssignmentElimination() {
		passes.MaybeAdd(this_.deadPropertyAssignmentElimination)
	}
	if this_.GetPassConfigBase().options.ShouldOptimizeCalls() {
		passes.MaybeAdd(this_.optimizeCalls)
	}
	if this_.GetPassConfigBase().options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(this_.j2clConstantHoisterPass)
		passes.MaybeAdd(this_.j2clClinitPass)
	}
	if this_.GetPassConfigBase().options.GetInlineFunctionsLevel() != StaticGet_CompilerOptions_Reach_NONE() {
		passes.MaybeAdd(this_.inlineFunctions)
	}
	if this_.GetPassConfigBase().options.ShouldInlineVariables() || this_.GetPassConfigBase().options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(this_.inlineVariables)
	} else if this_.GetPassConfigBase().options.ShouldInlineConstantVars() {
		passes.MaybeAdd(this_.inlineConstants)
	}
	if this_.GetPassConfigBase().options.ShouldRunDeadAssignmentElimination() {
		passes.MaybeAdd(this_.deadAssignmentsElimination)
	}
	if this_.GetPassConfigBase().options.GetCollapseObjectLiterals() {
		passes.MaybeAdd(this_.collapseObjectLiterals)
	}
	if this_.priv_DefaultPassConfig_shouldRunRemoveUnusedCode_2394() {
		passes.MaybeAdd(this_.removeUnusedCode)
	}
	if this_.GetPassConfigBase().options.ShouldFoldConstants() {
		passes.MaybeAdd(this_.peepholeOptimizations)
	}
	passes.AssertAllLoopablePasses()
	return passes
}

// Java method: DefaultPassConfig.addNonTypedAstNormalizationPasses(com.google.javascript.jscomp.PassListBuilder passes)
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_addNonTypedAstNormalizationPasses_2393(passes *PassListBuilder) {
	_ = passes
	passes.MaybeAdd(this_.removeCastNodes)
	passes.MaybeAdd(this_.typesToColors)
}

// Java method: DefaultPassConfig.shouldRunRemoveUnusedCode()
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_shouldRunRemoveUnusedCode_2394() bool {
	return ((((this_.GetPassConfigBase().options.ShouldRemoveUnusedVariables() || this_.GetPassConfigBase().options.ShouldRemoveUnusedLocalVariables()) || this_.GetPassConfigBase().options.ShouldRemoveUnusedPrototypeProperties()) || this_.GetPassConfigBase().options.IsRemoveUnusedClassProperties()) || this_.GetPassConfigBase().options.GetRewritePolyfills())
}

// Java method: DefaultPassConfig.assertValidOrderForChecks(com.google.javascript.jscomp.PassListBuilder checks)
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_assertValidOrderForChecks_2395(checks *PassListBuilder) {
	_ = checks
	checks.AssertPassOrder(this_.declaredGlobalExternsOnWindow, this_.checkVars, "declaredGlobalExternsOnWindow must happen before VarCheck, which adds synthetic externs")
	checks.AssertPassOrder(this_.chromePass, this_.checkJsDocAndEs6Modules, "The ChromePass must run before after JsDoc and Es6 module checking.")
	checks.AssertPassOrder(this_.closureRewriteModule, this_.processDefinesCheck, "Must rewrite goog.module before processing @define's, so that @defines in modules work.")
	checks.AssertPassOrder(this_.closurePrimitives, this_.polymerPass, "The Polymer pass must run after goog.provide processing.")
	checks.AssertPassOrder(this_.chromePass, this_.polymerPass, "The Polymer pass must run after ChromePass processing.")
	checks.AssertPassOrder(this_.polymerPass, this_.suspiciousCode, "The Polymer pass must run before suspiciousCode processing.")
	checks.AssertPassOrder(this_.addSyntheticScript, this_.gatherModuleMetadataPass, "Cannot add a synthetic script node after module metadata creation.")
	checks.AssertPassOrder(this_.closureRewriteModule, this_.removeSyntheticScript, "Synthetic script node should be removed only after module rewriting.")
	checks.AssertPassOrder(this_.closureRewriteModule, this_.rewriteCallerCodeLocation, "ClosureRewriteModule must happen before RewriteCallerCodeLocation, so that exported functions and call sites are rewritten correctly.")
	checks.AssertPassOrder(this_.closureRewriteModule, TranspilationPasses_GetEs6RewriteDestructuring(StaticGet_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS()), "RewriteCallerCodeLocation must happen before Es6RewriteDestructuring, because we need ReWriteCallerCodeLocation to run before default parameters get rewritten.")
	checks.AssertPassOrder(this_.closureRewriteModule, TranspilationPasses_GetEs6RewriteDestructuring(StaticGet_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST()), "RewriteCallerCodeLocation must happen before Es6RewriteDestructuring, because we need ReWriteCallerCodeLocation to run before default parameters get rewritten.")
	if checks.Contains(this_.closureGoogScopeAliases) {
		if !(checks.Contains(this_.checkVariableReferences)) {
			j2gort.CheckStateStr(false, "goog.scope processing requires variable checking")
		}
	}
	checks.AssertPassOrder(this_.gatherModuleMetadataPass, this_.closureCheckModule, "Need to gather module metadata before checking closure modules.")
	checks.AssertPassOrder(this_.gatherModuleMetadataPass, this_.createModuleMapPass, "Need to gather module metadata before scanning modules.")
	checks.AssertPassOrder(this_.createModuleMapPass, this_.rewriteCommonJsModules, "Need to gather module information before rewriting CommonJS modules.")
	checks.AssertPassOrder(this_.rewriteScriptsToEs6Modules, this_.gatherModuleMetadataPass, "Need to gather module information after rewriting scripts to modules.")
	checks.AssertPassOrder(this_.gatherModuleMetadataPass, this_.checkMissingRequires, "Need to gather module information before checking for missing requires.")
}

// Java method: DefaultPassConfig.assertValidOrderForOptimizations(com.google.javascript.jscomp.PassListBuilder optimizations)
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_assertValidOrderForOptimizations_2396(optimizations *PassListBuilder) {
	_ = optimizations
	optimizations.AssertPassOrder(this_.j2clPass, StaticGet_TranspilationPasses_RewriteGenerators(), "J2CL normalization should be done before generator re-writing.")
	optimizations.AssertPassOrder(StaticGet_TranspilationPasses_RewritePolyfills(), this_.processDefinesOptimize, "Polyfill injection must be done before processDefines as some polyfills reference goog.defines.")
	optimizations.AssertPassOrder(StaticGet_TranspilationPasses_InjectTranspilationRuntimeLibraries(), this_.processDefinesOptimize, "Runtime library injection must be done before processDefines some runtime libraries reference goog.defines.")
	optimizations.AssertPassOrder(this_.processDefinesOptimize, this_.j2clUtilGetDefineRewriterPass, "J2CL define re-writing should be done after processDefines since it relies on Compiler#getDefineNames to have been populated by it.")
	optimizations.AssertPassOrder(this_.removeUnusedCode, this_.isolatePolyfills, "Polyfill isolation should be done after RemovedUnusedCode. Otherwise unused polyfill removal will not find any polyfill usages and will delete all polyfills.")
	optimizations.AssertPassOrder(StaticGet_TranspilationPasses_InstrumentAsyncContext(), StaticGet_TranspilationPasses_RewriteAsyncIteration(), "AsyncContext should be instrumentated before await and/or yield is transpiled away")
	optimizations.AssertPassOrder(StaticGet_TranspilationPasses_InstrumentAsyncContext(), StaticGet_TranspilationPasses_RewriteAsyncFunctions(), "AsyncContext should be instrumentated before await and/or yield is transpiled away")
	optimizations.AssertPassOrder(StaticGet_TranspilationPasses_InstrumentAsyncContext(), StaticGet_TranspilationPasses_RewriteGenerators(), "AsyncContext should be instrumentated before await and/or yield is transpiled away")
}

// Java method: DefaultPassConfig.getFullReplaceMessagesPass()
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_getFullReplaceMessagesPass_2397() PassFactory {
	return PassFactory_S_Builder().SetName("replaceMessages").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_99 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_99)
		_ = compiler
		return NewReplaceMessages(compiler, this_.GetPassConfigBase().options.GetMessageBundle(), this_.GetPassConfigBase().options.GetStrictMessageReplacement()).GetFullReplacementPass()
	})).GetPassFactory_BuilderBase(

	// Java method: DefaultPassConfig.getProtectMessagesPass()
	).Build()
}

func (this_ *DefaultPassConfig) priv_DefaultPassConfig_getProtectMessagesPass_2398() PassFactory {
	return PassFactory_S_Builder().SetName("protectMessages").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_100 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_100)
		_ = compiler
		return NewReplaceMessages(compiler, this_.GetPassConfigBase().options.GetMessageBundle(), this_.GetPassConfigBase().options.GetStrictMessageReplacement()).GetMsgProtectionPass()
	})).GetPassFactory_BuilderBase(

	// Java method: DefaultPassConfig.getReplaceProtectedMessagesPass()
	).Build()
}

func (this_ *DefaultPassConfig) priv_DefaultPassConfig_getReplaceProtectedMessagesPass_2399() PassFactory {
	return PassFactory_S_Builder().SetName("replaceProtectedMessages").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_101 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_101)
		_ = compiler
		return NewReplaceMessages(compiler, this_.GetPassConfigBase().options.GetMessageBundle(), this_.GetPassConfigBase().options.GetStrictMessageReplacement()).GetReplacementCompletionPass()
	})).GetPassFactory_BuilderBase(

	// Java static method: DefaultPassConfig.createPeepholeOptimizationsPass(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String passName, boolean expectAstIsNormalized)
	).Build()
}

func DefaultPassConfig_CreatePeepholeOptimizationsPass(compiler AbstractCompiler, passName string, expectAstIsNormalized bool) CompilerPass {
	clinit_DefaultPassConfig()
	_ = compiler
	_ = passName
	_ = expectAstIsNormalized
	if !(expectAstIsNormalized == compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckArgumentMsg(false, j2gort.PtrToAny(compiler.GetLifeCycleStage()))
	}
	var late bool = false
	_ = late
	var useTypesForOptimization bool = compiler.GetOptions().ShouldUseTypesForLocalOptimization()
	_ = useTypesForOptimization
	var optimizations j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = optimizations
	if expectAstIsNormalized {
		optimizations.Add_1(NewMinimizeExitPoints())
	}
	optimizations.Add_1(NewPeepholeMinimizeConditions(late))
	optimizations.Add_1(NewPeepholeSubstituteAlternateSyntax(late))
	optimizations.Add_1(NewPeepholeReplaceKnownMethods(late, useTypesForOptimization))
	optimizations.Add_1(NewPeepholeRemoveDeadCode())
	if compiler.GetOptions().GetJ2clPass().ShouldAddJ2clPasses() {
		optimizations.Add_1(NewJ2clEqualitySameRewriterPass(useTypesForOptimization))
		optimizations.Add_1(NewJ2clStringValueOfRewriterPass())
		optimizations.Add_1(NewJ2clUndefinedChecksRewriterPass())
	}
	optimizations.Add_1(NewPeepholeFoldConstants(late, useTypesForOptimization))
	optimizations.Add_1(NewPeepholeCollectPropertyAssignments())
	return NewPeepholeOptimizationsPass_3_AbstractCompiler_string_Ext_java_util_List(compiler, passName, optimizations)
}

// Java static method: DefaultPassConfig.combineChecks(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.List<com.google.javascript.jscomp.NodeTraversal.Callback> callbacks)
func DefaultPassConfig_CombineChecks(compiler AbstractCompiler, callbacks j2gort.Ext_java_util_List) CompilerPass {
	clinit_DefaultPassConfig()
	_ = compiler
	_ = callbacks
	if !(!(callbacks.IsEmpty())) {
		j2gort.CheckArgument(false)
	}
	return NewCombinedCompilerPass_2_AbstractCompiler_Ext_java_util_List(compiler, callbacks)
}

// Java method: DefaultPassConfig.createProcessDefines(com.google.javascript.jscomp.ProcessDefines.Mode mode)
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_createProcessDefines_2400(mode ProcessDefines_Mode) PassFactory {
	_ = mode
	return PassFactory_S_Builder().SetName(("processDefines_" + j2gort.StrValString((mode).Name()))).SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_102 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_102)
		_ = compiler
		return j2gort.PtrToAny(NewProcessDefines_Builder(compiler).PutReplacements(DefaultPassConfig_GetAdditionalReplacements(this_.GetPassConfigBase().options)).PutReplacements(j2gort.AsIface[j2gort.Ext_java_util_Map](this_.GetPassConfigBase().options.GetDefineReplacements())).SetMode(mode).SetRecognizeClosureDefines(compiler.GetOptions().GetClosurePass()).SetEnableZonesDefineName(this_.GetPassConfigBase().options.GetEnableZonesDefineName()).SetZoneInputPattern(this_.GetPassConfigBase().options.GetZoneInputPattern()).SetUnknownDefinesToIgnore(this_.GetPassConfigBase().options.GetUnknownDefinesToIgnore()).Build())
	})).GetPassFactory_BuilderBase(

	// Java method: DefaultPassConfig.runVariableRenaming(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.VariableMap prevVariableMap, com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
	).Build()
}

func (this_ *DefaultPassConfig) priv_DefaultPassConfig_runVariableRenaming_2401(compiler AbstractCompiler, prevVariableMap *VariableMap, externs pkg_rhino.Node, root pkg_rhino.Node) *VariableMap {
	_ = compiler
	_ = prevVariableMap
	_ = externs
	_ = root
	var reservedChars j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	_ = reservedChars
	var reservedNames j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = reservedNames
	if this_.GetPassConfigBase().options.GetRenamePrefixNamespace() != j2gort.NullString {
		j2gort.SetAddStr(reservedNames, this_.GetPassConfigBase().options.GetRenamePrefixNamespace())
	}
	reservedNames.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](compiler.GetExportedNames()))
	reservedNames.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](pkg_parsing.ParserRunner_GetReservedVars()))
	clinit_RenameVars()
	var val_103 RenameVars
	InitRenameVars(&val_103, compiler, this_.GetPassConfigBase().options.GetRenamePrefix(), (this_.GetPassConfigBase().options.GetVariableRenaming() == StaticGet_VariableRenamingPolicy_LOCAL()), this_.GetPassConfigBase().options.ShouldGeneratePseudoNames(), this_.GetPassConfigBase().options.ShouldPreferStableNames(), prevVariableMap, reservedChars, reservedNames, this_.GetPassConfigBase().options.GetNameGenerator())
	var rn *RenameVars = &val_103
	_ = rn
	rn.Process_2(externs, root)
	return rn.GetVariableMap()
}

// Java method: DefaultPassConfig.createGatherExternProperties(com.google.javascript.jscomp.GatherExternProperties.Mode mode)
func (this_ *DefaultPassConfig) priv_DefaultPassConfig_createGatherExternProperties_2402(mode GatherExternProperties_Mode) PassFactory {
	_ = mode
	return PassFactory_S_Builder().SetName("gatherExternProperties").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_104 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_104)
		_ = compiler
		return NewGatherExternProperties(compiler, mode)
	})).GetPassFactory_BuilderBase(

	// Java method: DefaultPassConfig.getCustomPasses(com.google.javascript.jscomp.CustomPassExecutionTime executionTime)
	).Build()
}

func (this_ *DefaultPassConfig) priv_DefaultPassConfig_getCustomPasses_2403(executionTime CustomPassExecutionTime) PassFactory {
	_ = executionTime
	return PassFactory_S_Builder().SetName("runCustomPasses").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_105 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_105)
		_ = compiler
		return DefaultPassConfig_RunInSerial(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.GetPassConfigBase().options.GetCustomPassesAt(executionTime)))
	})).GetPassFactory_BuilderBase(

	// Java static method: DefaultPassConfig.runInSerial(java.util.Collection<com.google.javascript.jscomp.CompilerPass> passes)
	).Build()
}

func DefaultPassConfig_RunInSerial(passes j2gort.Ext_java_util_Collection) CompilerPass {
	clinit_DefaultPassConfig()
	_ = passes
	return NewDefaultPassConfig__Anon_34(passes)
}

// Java static method: DefaultPassConfig.getAdditionalReplacements(com.google.javascript.jscomp.CompilerOptions options)
func DefaultPassConfig_GetAdditionalReplacements(options *CompilerOptions) j2gort.Ext_java_util_Map {
	clinit_DefaultPassConfig()
	_ = options
	var additionalReplacements j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = additionalReplacements
	if options.ShouldMarkAsCompiled() || options.GetClosurePass() {
		j2gort.MapPutStr(additionalReplacements, "COMPILED", j2gort.BoxNode(pkg_rhino.IR_TrueNode()))
	}
	if (options.GetClosurePass() && (options.GetLocale() != j2gort.NullString)) && !(options.DoLateLocalization()) {
		j2gort.MapPutStr(additionalReplacements, "goog.LOCALE", j2gort.BoxNode(pkg_rhino.IR_String(options.GetLocale())))
	}
	return additionalReplacements
}

// Java static method: DefaultPassConfig.requirePropertiesAreStaticallyAnalyzable(com.google.javascript.jscomp.CompilerOptions options)
func DefaultPassConfig_RequirePropertiesAreStaticallyAnalyzable(options *CompilerOptions) *PassFactory_PreconditionResult {
	clinit_DefaultPassConfig()
	_ = options
	return NewPassFactory_PreconditionResult(options.GetAssumePropertiesAreStaticallyAnalyzable(), "requires assumePropertiesAreStaticallyAnalyzable to be enabled")
}

// Java class:
type DefaultPassConfig__Anon_22 struct {
	outer_         *DefaultPassConfig
	cap_0_pass     *GenerateExports
	cap_1_compiler AbstractCompiler
	idHash_        atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_22) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_22) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_22) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$1", ""
}
func (this_ *DefaultPassConfig__Anon_22) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_22) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_22) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_22) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_22) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_22) initFields_DefaultPassConfig__Anon_22() {
	_ = this_
}

func NewDefaultPassConfig__Anon_22(outer_ *DefaultPassConfig, cap_p_pass *GenerateExports, cap_p_compiler AbstractCompiler) *DefaultPassConfig__Anon_22 {
	this_ := &DefaultPassConfig__Anon_22{}
	this_.outer_ = outer_
	this_.cap_0_pass = cap_p_pass
	this_.cap_1_compiler = cap_p_compiler
	this_.initFields_DefaultPassConfig__Anon_22()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DefaultPassConfig__Anon_22) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.cap_0_pass.Process_2(externs, root)
	this_.cap_1_compiler.AddExportedNames(this_.cap_0_pass.GetExportedVariableNames())
}

// Java class:
type DefaultPassConfig__Anon_23 struct {
	outer_         *DefaultPassConfig
	cap_0_pass     *GatherRawExports
	cap_1_compiler AbstractCompiler
	idHash_        atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_23) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_23) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_23) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$2", ""
}
func (this_ *DefaultPassConfig__Anon_23) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_23) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_23) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_23) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_23) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_23) initFields_DefaultPassConfig__Anon_23() {
	_ = this_
}

func NewDefaultPassConfig__Anon_23(outer_ *DefaultPassConfig, cap_p_pass *GatherRawExports, cap_p_compiler AbstractCompiler) *DefaultPassConfig__Anon_23 {
	this_ := &DefaultPassConfig__Anon_23{}
	this_.outer_ = outer_
	this_.cap_0_pass = cap_p_pass
	this_.cap_1_compiler = cap_p_compiler
	this_.initFields_DefaultPassConfig__Anon_23()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DefaultPassConfig__Anon_23) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.cap_0_pass.Process_2(externs, root)
	this_.cap_1_compiler.AddExportedNames(this_.cap_0_pass.GetExportedVariableNames())
}

// Java class:
type DefaultPassConfig__Anon_24 struct {
	outer_         *DefaultPassConfig
	cap_0_compiler AbstractCompiler
	idHash_        atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_24) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_24) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_24) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$3", ""
}
func (this_ *DefaultPassConfig__Anon_24) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_24) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_24) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_24) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_24) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_24) initFields_DefaultPassConfig__Anon_24() {
	_ = this_
}

func NewDefaultPassConfig__Anon_24(outer_ *DefaultPassConfig, cap_p_compiler AbstractCompiler) *DefaultPassConfig__Anon_24 {
	this_ := &DefaultPassConfig__Anon_24{}
	this_.outer_ = outer_
	this_.cap_0_compiler = cap_p_compiler
	this_.initFields_DefaultPassConfig__Anon_24()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node jsRoot)
func (this_ *DefaultPassConfig__Anon_24) Process_2(externs pkg_rhino.Node, jsRoot pkg_rhino.Node) {
	_ = externs
	_ = jsRoot
	var cssNames *j2gort.Ext_java_util_Optional
	_ = cssNames
	if this_.outer_.GetPassConfigBase().options.ShouldGatherCssNames() {
		cssNames = j2gort.Ext_java_util_Optional_Of(j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder())
	} else {
		cssNames = j2gort.Ext_java_util_Optional_Empty()
	}
	clinit_ReplaceCssNames()
	var val_106 ReplaceCssNames
	InitReplaceCssNames(&val_106, this_.cap_0_compiler, this_.outer_.GetPassConfigBase().options.GetCssRenamingMap(), NewLambda_ReplaceCssNames_CssNameCollector(func(cssName string) {
		_ = cssName
		cssNames.IfPresent(j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_107 any) {
			builder := j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](lp_107)
			_ = builder
			j2gort.ImmSetBuilderAddStr(builder, cssName)
		}))
	}), j2gort.AsIface[j2gort.Ext_java_util_Set](this_.outer_.GetPassConfigBase().options.GetCssRenamingSkiplist()))
	var pass *ReplaceCssNames = &val_106
	_ = pass
	pass.Process_2(externs, jsRoot)
	this_.cap_0_compiler.SetCssNames(func() j2gort.Ext_java_util_Set {
		if cssNames.IsPresent() {
			return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](cssNames.Get()).Build())
		}
		return nil
	}())
}

// Java class:
type DefaultPassConfig__Anon_25 struct {
	outer_         *DefaultPassConfig
	cap_0_pass     *CheckRegExp
	cap_1_compiler AbstractCompiler
	idHash_        atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_25) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_25) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_25) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$4", ""
}
func (this_ *DefaultPassConfig__Anon_25) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_25) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_25) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_25) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_25) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_25) initFields_DefaultPassConfig__Anon_25() {
	_ = this_
}

func NewDefaultPassConfig__Anon_25(outer_ *DefaultPassConfig, cap_p_pass *CheckRegExp, cap_p_compiler AbstractCompiler) *DefaultPassConfig__Anon_25 {
	this_ := &DefaultPassConfig__Anon_25{}
	this_.outer_ = outer_
	this_.cap_0_pass = cap_p_pass
	this_.cap_1_compiler = cap_p_compiler
	this_.initFields_DefaultPassConfig__Anon_25()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DefaultPassConfig__Anon_25) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.cap_0_pass.Process_2(externs, root)
	this_.cap_1_compiler.SetHasRegExpGlobalReferences(this_.cap_0_pass.IsGlobalRegExpPropertiesUsed())
}

// Java class:
type DefaultPassConfig__Anon_26 struct {
	outer_            *DefaultPassConfig
	cap_0_regExpCheck *CheckRegExp
	cap_1_compiler    AbstractCompiler
	idHash_           atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_26) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_26) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_26) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$5", ""
}
func (this_ *DefaultPassConfig__Anon_26) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_26) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_26) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_26) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_26) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_26) initFields_DefaultPassConfig__Anon_26() {
	_ = this_
}

func NewDefaultPassConfig__Anon_26(outer_ *DefaultPassConfig, cap_p_regExpCheck *CheckRegExp, cap_p_compiler AbstractCompiler) *DefaultPassConfig__Anon_26 {
	this_ := &DefaultPassConfig__Anon_26{}
	this_.outer_ = outer_
	this_.cap_0_regExpCheck = cap_p_regExpCheck
	this_.cap_1_compiler = cap_p_compiler
	this_.initFields_DefaultPassConfig__Anon_26()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DefaultPassConfig__Anon_26) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.cap_0_regExpCheck.Process_2(externs, root)
	this_.cap_1_compiler.SetHasRegExpGlobalReferences(this_.cap_0_regExpCheck.IsGlobalRegExpPropertiesUsed())
}

// Java class:
type DefaultPassConfig__Anon_27 struct {
	outer_         *DefaultPassConfig
	cap_0_compiler AbstractCompiler
	idHash_        atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_27) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_27) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_27) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$6", ""
}
func (this_ *DefaultPassConfig__Anon_27) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_27) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_27) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_27) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_27) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_27) initFields_DefaultPassConfig__Anon_27() {
	_ = this_
}

func NewDefaultPassConfig__Anon_27(outer_ *DefaultPassConfig, cap_p_compiler AbstractCompiler) *DefaultPassConfig__Anon_27 {
	this_ := &DefaultPassConfig__Anon_27{}
	this_.outer_ = outer_
	this_.cap_0_compiler = cap_p_compiler
	this_.initFields_DefaultPassConfig__Anon_27()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node jsRoot)
func (this_ *DefaultPassConfig__Anon_27) Process_2(externs pkg_rhino.Node, jsRoot pkg_rhino.Node) {
	_ = externs
	_ = jsRoot
	{
		clinit_ProcessTweaks()
		var val_108 ProcessTweaks
		InitProcessTweaks(&val_108, this_.cap_0_compiler, this_.outer_.GetPassConfigBase().options.GetTweakProcessing().ShouldStrip())
		(&val_108).Process_2(externs, jsRoot)
	}
}

// Java class:
type DefaultPassConfig__Anon_28 struct {
	outer_         *DefaultPassConfig
	cap_0_compiler AbstractCompiler
	idHash_        atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_28) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_28) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_28) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$7", ""
}
func (this_ *DefaultPassConfig__Anon_28) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_28) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_28) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_28) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_28) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_28) initFields_DefaultPassConfig__Anon_28() {
	_ = this_
}

func NewDefaultPassConfig__Anon_28(outer_ *DefaultPassConfig, cap_p_compiler AbstractCompiler) *DefaultPassConfig__Anon_28 {
	this_ := &DefaultPassConfig__Anon_28{}
	this_.outer_ = outer_
	this_.cap_0_compiler = cap_p_compiler
	this_.initFields_DefaultPassConfig__Anon_28()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node jsRoot)
func (this_ *DefaultPassConfig__Anon_28) Process_2(externs pkg_rhino.Node, jsRoot pkg_rhino.Node) {
	_ = externs
	_ = jsRoot
	var options *CompilerOptions = this_.cap_0_compiler.GetOptions()
	_ = options
	clinit_StripCode()
	var val_109 StripCode
	InitStripCode(&val_109, this_.cap_0_compiler, options.GetStripTypes(), options.GetStripNameSuffixes(), options.GetStripNamePrefixes(), options.GetTweakProcessing().ShouldStrip())
	var pass *StripCode = &val_109
	_ = pass
	pass.Process_2(externs, jsRoot)
}

// Java class:
type DefaultPassConfig__Anon_29 struct {
	outer_         *DefaultPassConfig
	cap_0_compiler AbstractCompiler
	idHash_        atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_29) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_29) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_29) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$8", ""
}
func (this_ *DefaultPassConfig__Anon_29) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_29) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_29) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_29) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_29) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_29) initFields_DefaultPassConfig__Anon_29() {
	_ = this_
}

func NewDefaultPassConfig__Anon_29(outer_ *DefaultPassConfig, cap_p_compiler AbstractCompiler) *DefaultPassConfig__Anon_29 {
	this_ := &DefaultPassConfig__Anon_29{}
	this_.outer_ = outer_
	this_.cap_0_compiler = cap_p_compiler
	this_.initFields_DefaultPassConfig__Anon_29()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DefaultPassConfig__Anon_29) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	clinit_ReplaceIdGenerators()
	var val_110 ReplaceIdGenerators
	InitReplaceIdGenerators(&val_110, this_.cap_0_compiler, this_.outer_.GetPassConfigBase().options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_TEMPLATE_LITERALS()), this_.outer_.GetPassConfigBase().options.GetIdGenerators(), this_.outer_.GetPassConfigBase().options.ShouldGeneratePseudoNames(), this_.outer_.GetPassConfigBase().options.GetIdGeneratorsMapSerialized(), this_.outer_.GetPassConfigBase().options.GetXidHashFunction())
	var pass *ReplaceIdGenerators = &val_110
	_ = pass
	pass.Process_2(externs, root)
	this_.cap_0_compiler.SetIdGeneratorMap(pass.GetSerializedIdMappings())
}

// Java class:
type DefaultPassConfig__Anon_30 struct {
	outer_         *DefaultPassConfig
	cap_0_compiler AbstractCompiler
	idHash_        atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_30) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_30) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_30) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$9", ""
}
func (this_ *DefaultPassConfig__Anon_30) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_30) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_30) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_30) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_30) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_30) initFields_DefaultPassConfig__Anon_30() {
	_ = this_
}

func NewDefaultPassConfig__Anon_30(outer_ *DefaultPassConfig, cap_p_compiler AbstractCompiler) *DefaultPassConfig__Anon_30 {
	this_ := &DefaultPassConfig__Anon_30{}
	this_.outer_ = outer_
	this_.cap_0_compiler = cap_p_compiler
	this_.initFields_DefaultPassConfig__Anon_30()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DefaultPassConfig__Anon_30) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var pass *ReplaceStrings = NewReplaceStrings(this_.cap_0_compiler, this_.outer_.GetPassConfigBase().options.GetReplaceStringsPlaceholderToken(), j2gort.AsIface[j2gort.Ext_java_util_List](this_.outer_.GetPassConfigBase().options.GetReplaceStringsFunctionDescriptions()))
	_ = pass
	pass.Process_2(externs, root)
	this_.cap_0_compiler.SetStringMap(pass.GetStringMap())
}

// Java class:
type DefaultPassConfig__Anon_31 struct {
	outer_         *DefaultPassConfig
	cap_0_compiler AbstractCompiler
	idHash_        atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_31) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_31) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_31) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$10", ""
}
func (this_ *DefaultPassConfig__Anon_31) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_31) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_31) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_31) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_31) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_31) initFields_DefaultPassConfig__Anon_31() {
	_ = this_
}

func NewDefaultPassConfig__Anon_31(outer_ *DefaultPassConfig, cap_p_compiler AbstractCompiler) *DefaultPassConfig__Anon_31 {
	this_ := &DefaultPassConfig__Anon_31{}
	this_.outer_ = outer_
	this_.cap_0_compiler = cap_p_compiler
	this_.initFields_DefaultPassConfig__Anon_31()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DefaultPassConfig__Anon_31) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.cap_0_compiler.SetLifeCycleStage(StaticGet_AbstractCompiler_LifeCycleStage_RAW())
}

// Java class:
type DefaultPassConfig__Anon_32 struct {
	outer_                *DefaultPassConfig
	cap_0_compiler        AbstractCompiler
	cap_1_prevPropertyMap *VariableMap
	idHash_               atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_32) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_32) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_32) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$11", ""
}
func (this_ *DefaultPassConfig__Anon_32) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_32) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_32) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_32) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_32) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_32) initFields_DefaultPassConfig__Anon_32() {
	_ = this_
}

func NewDefaultPassConfig__Anon_32(outer_ *DefaultPassConfig, cap_p_compiler AbstractCompiler, cap_p_prevPropertyMap *VariableMap) *DefaultPassConfig__Anon_32 {
	this_ := &DefaultPassConfig__Anon_32{}
	this_.outer_ = outer_
	this_.cap_0_compiler = cap_p_compiler
	this_.cap_1_prevPropertyMap = cap_p_prevPropertyMap
	this_.initFields_DefaultPassConfig__Anon_32()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DefaultPassConfig__Anon_32) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	clinit_RenameProperties()
	var val_111 RenameProperties
	InitRenameProperties_6(&val_111, this_.cap_0_compiler, this_.outer_.GetPassConfigBase().options.ShouldGeneratePseudoNames(), this_.cap_1_prevPropertyMap, j2gort.AsIface[j2gort.Ext_java_util_Set](this_.outer_.GetPassConfigBase().options.GetPropertyReservedNamingFirstChars()), j2gort.AsIface[j2gort.Ext_java_util_Set](this_.outer_.GetPassConfigBase().options.GetPropertyReservedNamingNonFirstChars()), this_.outer_.GetPassConfigBase().options.GetNameGenerator())
	var rprop *RenameProperties = &val_111
	_ = rprop
	rprop.Process_2(externs, root)
	this_.cap_0_compiler.SetPropertyMap(rprop.GetPropertyMap())
}

// Java class:
type DefaultPassConfig__Anon_33 struct {
	outer_                *DefaultPassConfig
	cap_0_compiler        AbstractCompiler
	cap_1_prevVariableMap *VariableMap
	idHash_               atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_33) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_33) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_33) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$12", ""
}
func (this_ *DefaultPassConfig__Anon_33) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_33) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_33) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_33) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_33) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_33) initFields_DefaultPassConfig__Anon_33() {
	_ = this_
}

func NewDefaultPassConfig__Anon_33(outer_ *DefaultPassConfig, cap_p_compiler AbstractCompiler, cap_p_prevVariableMap *VariableMap) *DefaultPassConfig__Anon_33 {
	this_ := &DefaultPassConfig__Anon_33{}
	this_.outer_ = outer_
	this_.cap_0_compiler = cap_p_compiler
	this_.cap_1_prevVariableMap = cap_p_prevVariableMap
	this_.initFields_DefaultPassConfig__Anon_33()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DefaultPassConfig__Anon_33) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.cap_0_compiler.SetVariableMap(this_.outer_.priv_DefaultPassConfig_runVariableRenaming_2401(this_.cap_0_compiler, this_.cap_1_prevVariableMap, externs, root))
}

// Java class:
type DefaultPassConfig__Anon_35 struct {
	outer_         *DefaultPassConfig
	cap_0_compiler AbstractCompiler
	idHash_        atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_35) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_35) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_35) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$14", ""
}
func (this_ *DefaultPassConfig__Anon_35) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_35) HasIdentityEquals_()                  {}
func (this_ *DefaultPassConfig__Anon_35) SetOuter_(o *DefaultPassConfig) { this_.outer_ = o }
func (this_ *DefaultPassConfig__Anon_35) GetOuter_() *DefaultPassConfig  { return this_.outer_ }
func (*DefaultPassConfig__Anon_35) Implements_jscomp_CompilerPass()      {}

func (this_ *DefaultPassConfig__Anon_35) initFields_DefaultPassConfig__Anon_35() {
	_ = this_
}

func NewDefaultPassConfig__Anon_35(outer_ *DefaultPassConfig, cap_p_compiler AbstractCompiler) *DefaultPassConfig__Anon_35 {
	this_ := &DefaultPassConfig__Anon_35{}
	this_.outer_ = outer_
	this_.cap_0_compiler = cap_p_compiler
	this_.initFields_DefaultPassConfig__Anon_35()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DefaultPassConfig__Anon_35) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NewLocaleDataPasses_LocaleSubstitutions(this_.cap_0_compiler, this_.cap_0_compiler.GetOptions().GetLocale()).Process_2(externs, root)
}

// Java class:
type DefaultPassConfig__Anon_34 struct {
	cap_0_passes j2gort.Ext_java_util_Collection
	idHash_      atomic.Int32
}

func (this_ *DefaultPassConfig__Anon_34) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultPassConfig__Anon_34) IsJavaNil_() bool { return this_ == nil }
func (*DefaultPassConfig__Anon_34) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultPassConfig$13", ""
}
func (this_ *DefaultPassConfig__Anon_34) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultPassConfig__Anon_34) HasIdentityEquals_()             {}
func (*DefaultPassConfig__Anon_34) Implements_jscomp_CompilerPass() {}

func (this_ *DefaultPassConfig__Anon_34) initFields_DefaultPassConfig__Anon_34() {
	_ = this_
}

func NewDefaultPassConfig__Anon_34(cap_p_passes j2gort.Ext_java_util_Collection) *DefaultPassConfig__Anon_34 {
	this_ := &DefaultPassConfig__Anon_34{}
	this_.cap_0_passes = cap_p_passes
	this_.initFields_DefaultPassConfig__Anon_34()
	return this_
}

// Java method: .process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DefaultPassConfig__Anon_34) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		cur_112 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.cap_0_passes))
		for cur_112.HasNext() {
			rawElem_113 := cur_112.Next()
			pass := j2gort.AsIface[CompilerPass](rawElem_113)
			_ = pass
			pass.Process_2(externs, root)
		}
	}
}
