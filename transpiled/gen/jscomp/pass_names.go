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
// Port of com/google/javascript/jscomp/PassNames.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.PassNames
type PassNames struct {
	idHash_ atomic.Int32
}

func (this_ *PassNames) BindSelf_(s any) {
	_ = s
}
func (this_ *PassNames) IsJavaNil_() bool { return this_ == nil }
func (*PassNames) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.PassNames", "PassNames" }
func (this_ *PassNames) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PassNames) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PassNames", "com.google.javascript.jscomp.PassNames", func(v_ any) bool {
		return j2gort.InstanceOf[*PassNames](v_)
	})
}

func StaticGet_PassNames_AFTER_EARLY_OPTIMIZATION_LOOP() string { return "afterEarlyOptimizationLoop" }
func StaticGet_PassNames_AFTER_MAIN_OPTIMIZATIONS() string { return "afterMainOptimizations" }
func StaticGet_PassNames_AFTER_STANDARD_CHECKS() string { return "afterStandardChecks" }
func StaticGet_PassNames_AMBIGUATE_PROPERTIES() string { return "ambiguateProperties" }
func StaticGet_PassNames_ANALYZER_CHECKS() string { return "analyzerChecks" }
func StaticGet_PassNames_ANGULAR_PASS() string { return "angularPass" }
func StaticGet_PassNames_BEFORE_EARLY_OPTIMIZATION_LOOP() string { return "beforeEarlyOptimizationLoop" }
func StaticGet_PassNames_BEFORE_EARLY_OPTIMIZATIONS_TRANSPILATION() string { return "beforeEarlyOptimizationsTranspilation" }
func StaticGet_PassNames_BEFORE_STANDARD_OPTIMIZATIONS() string { return "beforeStandardOptimizations" }
func StaticGet_PassNames_OBFUSCATION_PASS_MARKER() string { return "obfuscationPassMarker" }
func StaticGet_PassNames_BEFORE_MAIN_OPTIMIZATIONS() string { return "beforeMainOptimizations" }
func StaticGet_PassNames_BEFORE_TYPE_CHECKING() string { return "beforeTypeChecking" }
func StaticGet_PassNames_BEFORE_SERIALIZATION() string { return "beforeSerialization" }
func StaticGet_PassNames_BEFORE_VARIABLE_RENAMING() string { return "beforeVariableRenaming" }
func StaticGet_PassNames_BEFORE_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS() string { return "beforeExtractPrototypeMemberDeclarations" }
func StaticGet_PassNames_CHECK_CONFORMANCE() string { return "checkConformance" }
func StaticGet_PassNames_CHECK_REG_EXP() string { return "checkRegExp" }
func StaticGet_PassNames_CHECK_TYPES() string { return "checkTypes" }
func StaticGet_PassNames_CHECK_VARIABLE_REFERENCES() string { return "checkVariableReferences" }
func StaticGet_PassNames_CHECK_VARS() string { return "checkVars" }
func StaticGet_PassNames_COALESCE_VARIABLE_NAMES() string { return "coalesceVariableNames" }
func StaticGet_PassNames_COLLAPSE_ANONYMOUS_FUNCTIONS() string { return "collapseAnonymousFunctions" }
func StaticGet_PassNames_COLLAPSE_OBJECT_LITERALS() string { return "collapseObjectLiterals" }
func StaticGet_PassNames_COLLAPSE_PROPERTIES() string { return "collapseProperties" }
func StaticGet_PassNames_COLLAPSE_VARIABLE_DECLARATIONS() string { return "collapseVariableDeclarations" }
func StaticGet_PassNames_CONVERT_TO_DOTTED_PROPERTIES() string { return "convertToDottedProperties" }
func StaticGet_PassNames_CREATE_MODULE_MAP() string { return "createModuleMap" }
func StaticGet_PassNames_CROSS_CHUNK_CODE_MOTION() string { return "crossChunkCodeMotion" }
func StaticGet_PassNames_CROSS_CHUNK_METHOD_MOTION() string { return "crossChunkMethodMotion" }
func StaticGet_PassNames_DEAD_ASSIGNMENT_ELIMINATION() string { return "deadAssignmentsElimination" }
func StaticGet_PassNames_DECLARED_GLOBAL_EXTERNS_ON_WINDOW() string { return "declaredGlobalExternsOnWindow" }
func StaticGet_PassNames_DESERIALIZE_COMPILER_STATE() string { return "deserializeCompilerState" }
func StaticGet_PassNames_DEVIRTUALIZE_METHODS() string { return "devirtualizeMethods" }
func StaticGet_PassNames_DISAMBIGUATE_PROPERTIES() string { return "disambiguateProperties" }
func StaticGet_PassNames_ES6_NORMALIZE_CLASSES() string { return "es6NormalizeClasses" }
func StaticGet_PassNames_EXPORT_TEST_FUNCTIONS() string { return "exportTestFunctions" }
func StaticGet_PassNames_EXTERN_EXPORTS() string { return "externExports" }
func StaticGet_PassNames_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS() string { return "extractPrototypeMemberDeclarations" }
func StaticGet_PassNames_FLOW_SENSITIVE_INLINE_VARIABLES() string { return "flowSensitiveInlineVariables" }
func StaticGet_PassNames_GATHER_GETTERS_AND_SETTERS() string { return "gatherGettersAndSetters" }
func StaticGet_PassNames_GATHER_MODULE_METADATA() string { return "gatherModuleMetadata" }
func StaticGet_PassNames_GATHER_RAW_EXPORTS() string { return "gatherRawExports" }
func StaticGet_PassNames_GENERATE_EXPORTS() string { return "generateExports" }
func StaticGet_PassNames_INFER_CONSTS() string { return "inferConsts" }
func StaticGet_PassNames_INFER_TYPES() string { return "inferTypes" }
func StaticGet_PassNames_INLINE_FUNCTIONS() string { return "inlineFunctions" }
func StaticGet_PassNames_INLINE_PROPERTIES() string { return "inlineProperties" }
func StaticGet_PassNames_INLINE_TYPE_ALIASES() string { return "inlineTypeAliases" }
func StaticGet_PassNames_INLINE_VARIABLES() string { return "inlineVariables" }
func StaticGet_PassNames_LATE_PEEPHOLE_OPTIMIZATIONS() string { return "latePeepholeOptimizations" }
func StaticGet_PassNames_LINT_CHECKS() string { return "lintChecks" }
func StaticGet_PassNames_MARK_UNNORMALIZED() string { return "markUnnormalized" }
func StaticGet_PassNames_NORMALIZE() string { return "normalize" }
func StaticGet_PassNames_OPTIMIZATIONS_HALFWAY_POINT() string { return "optimizationsHalfwayPoint" }
func StaticGet_PassNames_OPTIMIZE_CALLS() string { return "optimizeCalls" }
func StaticGet_PassNames_PARSE_INPUTS() string { return "parseInputs" }
func StaticGet_PassNames_PEEPHOLE_OPTIMIZATIONS() string { return "peepholeOptimizations" }
func StaticGet_PassNames_POST_NORMALIZE_PEEPHOLE() string { return "postNormalizePeephole" }
func StaticGet_PassNames_REWRITE_COMMON_JS_MODULES() string { return "rewriteCommonJsModules" }
func StaticGet_PassNames_REWRITE_SCRIPTS_TO_ES6_MODULES() string { return "rewriteScriptsToEs6Modules" }
func StaticGet_PassNames_REMOVE_UNUSED_CODE() string { return "removeUnusedCode" }
func StaticGet_PassNames_REPLACE_ID_GENERATORS() string { return "replaceIdGenerators" }
func StaticGet_PassNames_REPLACE_MESSAGES() string { return "replaceMessages" }
func StaticGet_PassNames_RESOLVE_TYPES() string { return "resolveTypes" }
func StaticGet_PassNames_REWRITE_FUNCTION_EXPRESSIONS() string { return "rewriteFunctionExpressions" }
func StaticGet_PassNames_RENAME_PROPERTIES() string { return "renameProperties" }
func StaticGet_PassNames_BEFORE_RENAME_PROPERTIES() string { return "beforeRenameProperties" }
func StaticGet_PassNames_STRIP_SIDE_EFFECT_PROTECTION() string { return "stripSideEffectProtection" }
func StaticGet_PassNames_WIZ_PASS() string { return "wizPass" }
func (this_ *PassNames) initFields_PassNames() {
	_ = this_
}

// Java constructor: PassNames()
func NewPassNames() *PassNames {
	this_ := &PassNames{}
	InitPassNames(this_)
	return this_
}

func InitPassNames(this_ *PassNames) {
	this_.initFields_PassNames()
}

