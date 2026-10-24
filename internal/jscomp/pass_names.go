// Port of com/google/javascript/jscomp/PassNames.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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

// PassNames: if the name of a pass is used in more than one place in the
// source, it's good to create a symbolic name here.
type PassNames struct{}

const (
	PassNames_AFTER_EARLY_OPTIMIZATION_LOOP                = "afterEarlyOptimizationLoop"
	PassNames_AFTER_MAIN_OPTIMIZATIONS                     = "afterMainOptimizations"
	PassNames_AFTER_STANDARD_CHECKS                        = "afterStandardChecks"
	PassNames_AMBIGUATE_PROPERTIES                         = "ambiguateProperties"
	PassNames_ANALYZER_CHECKS                              = "analyzerChecks"
	PassNames_ANGULAR_PASS                                 = "angularPass"
	PassNames_BEFORE_EARLY_OPTIMIZATION_LOOP               = "beforeEarlyOptimizationLoop"
	PassNames_BEFORE_EARLY_OPTIMIZATIONS_TRANSPILATION     = "beforeEarlyOptimizationsTranspilation"
	PassNames_BEFORE_STANDARD_OPTIMIZATIONS                = "beforeStandardOptimizations"
	PassNames_OBFUSCATION_PASS_MARKER                      = "obfuscationPassMarker"
	PassNames_BEFORE_MAIN_OPTIMIZATIONS                    = "beforeMainOptimizations"
	PassNames_BEFORE_TYPE_CHECKING                         = "beforeTypeChecking"
	PassNames_BEFORE_SERIALIZATION                         = "beforeSerialization"
	PassNames_BEFORE_VARIABLE_RENAMING                     = "beforeVariableRenaming"
	PassNames_BEFORE_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS = "beforeExtractPrototypeMemberDeclarations"
	PassNames_CHECK_CONFORMANCE                            = "checkConformance"
	PassNames_CHECK_REG_EXP                                = "checkRegExp"
	PassNames_CHECK_TYPES                                  = "checkTypes"
	PassNames_CHECK_VARIABLE_REFERENCES                    = "checkVariableReferences"
	PassNames_CHECK_VARS                                   = "checkVars"
	PassNames_COALESCE_VARIABLE_NAMES                      = "coalesceVariableNames"
	PassNames_COLLAPSE_ANONYMOUS_FUNCTIONS                 = "collapseAnonymousFunctions"
	PassNames_COLLAPSE_OBJECT_LITERALS                     = "collapseObjectLiterals"
	PassNames_COLLAPSE_PROPERTIES                          = "collapseProperties"
	PassNames_COLLAPSE_VARIABLE_DECLARATIONS               = "collapseVariableDeclarations"
	PassNames_CONVERT_TO_DOTTED_PROPERTIES                 = "convertToDottedProperties"
	PassNames_CREATE_MODULE_MAP                            = "createModuleMap"
	PassNames_CROSS_CHUNK_CODE_MOTION                      = "crossChunkCodeMotion"
	PassNames_CROSS_CHUNK_METHOD_MOTION                    = "crossChunkMethodMotion"
	PassNames_DEAD_ASSIGNMENT_ELIMINATION                  = "deadAssignmentsElimination"
	PassNames_DECLARED_GLOBAL_EXTERNS_ON_WINDOW            = "declaredGlobalExternsOnWindow"
	PassNames_DESERIALIZE_COMPILER_STATE                   = "deserializeCompilerState"
	PassNames_DEVIRTUALIZE_METHODS                         = "devirtualizeMethods"
	PassNames_DISAMBIGUATE_PROPERTIES                      = "disambiguateProperties"
	PassNames_ES6_NORMALIZE_CLASSES                        = "es6NormalizeClasses"
	PassNames_EXPORT_TEST_FUNCTIONS                        = "exportTestFunctions"
	PassNames_EXTERN_EXPORTS                               = "externExports"
	PassNames_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS        = "extractPrototypeMemberDeclarations"
	PassNames_FLOW_SENSITIVE_INLINE_VARIABLES              = "flowSensitiveInlineVariables"
	PassNames_GATHER_GETTERS_AND_SETTERS                   = "gatherGettersAndSetters"
	PassNames_GATHER_MODULE_METADATA                       = "gatherModuleMetadata"
	PassNames_GATHER_RAW_EXPORTS                           = "gatherRawExports"
	PassNames_GENERATE_EXPORTS                             = "generateExports"
	PassNames_INFER_CONSTS                                 = "inferConsts"
	PassNames_INFER_TYPES                                  = "inferTypes"
	PassNames_INLINE_FUNCTIONS                             = "inlineFunctions"
	PassNames_INLINE_PROPERTIES                            = "inlineProperties"
	PassNames_INLINE_TYPE_ALIASES                          = "inlineTypeAliases"
	PassNames_INLINE_VARIABLES                             = "inlineVariables"
	PassNames_LATE_PEEPHOLE_OPTIMIZATIONS                  = "latePeepholeOptimizations"
	PassNames_LINT_CHECKS                                  = "lintChecks"
	PassNames_MARK_UNNORMALIZED                            = "markUnnormalized"
	PassNames_NORMALIZE                                    = "normalize"
	PassNames_OPTIMIZATIONS_HALFWAY_POINT                  = "optimizationsHalfwayPoint"
	PassNames_OPTIMIZE_CALLS                               = "optimizeCalls"
	PassNames_PARSE_INPUTS                                 = "parseInputs"
	PassNames_PEEPHOLE_OPTIMIZATIONS                       = "peepholeOptimizations"
	PassNames_POST_NORMALIZE_PEEPHOLE                      = "postNormalizePeephole"
	PassNames_REWRITE_COMMON_JS_MODULES                    = "rewriteCommonJsModules"
	PassNames_REWRITE_SCRIPTS_TO_ES6_MODULES               = "rewriteScriptsToEs6Modules"
	PassNames_REMOVE_UNUSED_CODE                           = "removeUnusedCode"
	PassNames_REPLACE_ID_GENERATORS                        = "replaceIdGenerators"
	PassNames_REPLACE_MESSAGES                             = "replaceMessages"
	PassNames_RESOLVE_TYPES                                = "resolveTypes"
	PassNames_REWRITE_FUNCTION_EXPRESSIONS                 = "rewriteFunctionExpressions"
	PassNames_RENAME_PROPERTIES                            = "renameProperties"
	PassNames_BEFORE_RENAME_PROPERTIES                     = "beforeRenameProperties"
	PassNames_STRIP_SIDE_EFFECT_PROTECTION                 = "stripSideEffectProtection"
	PassNames_WIZ_PASS                                     = "wizPass"
)
