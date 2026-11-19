/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/DiagnosticGroups.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.DiagnosticGroups
type DiagnosticGroups struct {
	idHash_ atomic.Int32
}

func (this_ *DiagnosticGroups) BindSelf_(s any) {
	_ = s
}
func (this_ *DiagnosticGroups) IsJavaNil_() bool { return this_ == nil }
func (*DiagnosticGroups) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DiagnosticGroups", "DiagnosticGroups"
}
func (this_ *DiagnosticGroups) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DiagnosticGroups) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DiagnosticGroups", "com.google.javascript.jscomp.DiagnosticGroups", func(v_ any) bool {
		return j2gort.InstanceOf[*DiagnosticGroups](v_)
	})
}

var static_DiagnosticGroups_UNUSED *DiagnosticType

func StaticGet_DiagnosticGroups_UNUSED() *DiagnosticType {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_UNUSED
}
func StaticSet_DiagnosticGroups_UNUSED(v *DiagnosticType) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_UNUSED = v
}

var static_DiagnosticGroups_wildcardExcludedGroups j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_DiagnosticGroups_WildcardExcludedGroups() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_wildcardExcludedGroups
}
func StaticSet_DiagnosticGroups_WildcardExcludedGroups(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_wildcardExcludedGroups = v
}

var static_DiagnosticGroups_groupsByName j2gort.Ext_java_util_Map

func StaticGet_DiagnosticGroups_GroupsByName() j2gort.Ext_java_util_Map {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_groupsByName
}
func StaticSet_DiagnosticGroups_GroupsByName(v j2gort.Ext_java_util_Map) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_groupsByName = v
}
func StaticGet_DiagnosticGroups_DIAGNOSTIC_GROUP_NAMES() string {
	return "accessControls, checkPrototypalTypes, checkRegExp, checkTypes, checkVars, conformanceViolations, const, constantProperty, deprecated, deprecatedAnnotations, defineWithoutGoogDefine, duplicateMessage, es5Strict, externsValidation, functionParams, globalThis, invalidCasts, lintVarDeclarations, misplacedTypeAnnotation, missingCallToSuper, missingOverride, missingPolyfill, missingProperties, missingProvide, missingRequire, missingReturn, missingSourcesWarnings, moduleLoad, moduleImport, msgDescriptions, nonStandardJsDocs, partialAlias, polymer, reportUnknownTypes, strictCheckTypes, strictMissingProperties, strictModuleDepCheck, strictPrimitiveOperators, suspiciousCode, typeInvalidation, undefinedVars, underscore, unknownDefines, unusedLocalVariables, uselessCode, untranspilableFeatures,visibility, weakModuleGet"
}

var static_DiagnosticGroups_UNTRANSPILABLE_FEATURES *DiagnosticGroup

func StaticGet_DiagnosticGroups_UNTRANSPILABLE_FEATURES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_UNTRANSPILABLE_FEATURES
}
func StaticSet_DiagnosticGroups_UNTRANSPILABLE_FEATURES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_UNTRANSPILABLE_FEATURES = v
}

var static_DiagnosticGroups_MODULE_LOAD *DiagnosticGroup

func StaticGet_DiagnosticGroups_MODULE_LOAD() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MODULE_LOAD
}
func StaticSet_DiagnosticGroups_MODULE_LOAD(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MODULE_LOAD = v
}

var static_DiagnosticGroups_MODULE_IMPORT *DiagnosticGroup

func StaticGet_DiagnosticGroups_MODULE_IMPORT() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MODULE_IMPORT
}
func StaticSet_DiagnosticGroups_MODULE_IMPORT(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MODULE_IMPORT = v
}

var static_DiagnosticGroups_GLOBAL_THIS *DiagnosticGroup

func StaticGet_DiagnosticGroups_GLOBAL_THIS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_GLOBAL_THIS
}
func StaticSet_DiagnosticGroups_GLOBAL_THIS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_GLOBAL_THIS = v
}

var static_DiagnosticGroups_DEPRECATED *DiagnosticGroup

func StaticGet_DiagnosticGroups_DEPRECATED() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_DEPRECATED
}
func StaticSet_DiagnosticGroups_DEPRECATED(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_DEPRECATED = v
}

var static_DiagnosticGroups_UNDERSCORE *DiagnosticGroup

func StaticGet_DiagnosticGroups_UNDERSCORE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_UNDERSCORE
}
func StaticSet_DiagnosticGroups_UNDERSCORE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_UNDERSCORE = v
}

var static_DiagnosticGroups_VISIBILITY *DiagnosticGroup

func StaticGet_DiagnosticGroups_VISIBILITY() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_VISIBILITY
}
func StaticSet_DiagnosticGroups_VISIBILITY(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_VISIBILITY = v
}

var static_DiagnosticGroups_ACCESS_CONTROLS *DiagnosticGroup

func StaticGet_DiagnosticGroups_ACCESS_CONTROLS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_ACCESS_CONTROLS
}
func StaticSet_DiagnosticGroups_ACCESS_CONTROLS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_ACCESS_CONTROLS = v
}

var static_DiagnosticGroups_NON_STANDARD_JSDOC *DiagnosticGroup

func StaticGet_DiagnosticGroups_NON_STANDARD_JSDOC() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_NON_STANDARD_JSDOC
}
func StaticSet_DiagnosticGroups_NON_STANDARD_JSDOC(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_NON_STANDARD_JSDOC = v
}

var static_DiagnosticGroups_INVALID_CASTS *DiagnosticGroup

func StaticGet_DiagnosticGroups_INVALID_CASTS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_INVALID_CASTS
}
func StaticSet_DiagnosticGroups_INVALID_CASTS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_INVALID_CASTS = v
}

var static_DiagnosticGroups_STRICT_MODULE_DEP_CHECK *DiagnosticGroup

func StaticGet_DiagnosticGroups_STRICT_MODULE_DEP_CHECK() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_STRICT_MODULE_DEP_CHECK
}
func StaticSet_DiagnosticGroups_STRICT_MODULE_DEP_CHECK(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_STRICT_MODULE_DEP_CHECK = v
}

var static_DiagnosticGroups_VIOLATED_MODULE_DEP *DiagnosticGroup

func StaticGet_DiagnosticGroups_VIOLATED_MODULE_DEP() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_VIOLATED_MODULE_DEP
}
func StaticSet_DiagnosticGroups_VIOLATED_MODULE_DEP(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_VIOLATED_MODULE_DEP = v
}

var static_DiagnosticGroups_EXTERNS_VALIDATION *DiagnosticGroup

func StaticGet_DiagnosticGroups_EXTERNS_VALIDATION() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_EXTERNS_VALIDATION
}
func StaticSet_DiagnosticGroups_EXTERNS_VALIDATION(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_EXTERNS_VALIDATION = v
}

var static_DiagnosticGroups_UNKNOWN_DEFINES *DiagnosticGroup

func StaticGet_DiagnosticGroups_UNKNOWN_DEFINES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_UNKNOWN_DEFINES
}
func StaticSet_DiagnosticGroups_UNKNOWN_DEFINES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_UNKNOWN_DEFINES = v
}

var static_DiagnosticGroups_DEFINE_WITHOUT_GOOG_DEFINE *DiagnosticGroup

func StaticGet_DiagnosticGroups_DEFINE_WITHOUT_GOOG_DEFINE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_DEFINE_WITHOUT_GOOG_DEFINE
}
func StaticSet_DiagnosticGroups_DEFINE_WITHOUT_GOOG_DEFINE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_DEFINE_WITHOUT_GOOG_DEFINE = v
}

var static_DiagnosticGroups_TWEAKS *DiagnosticGroup

func StaticGet_DiagnosticGroups_TWEAKS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_TWEAKS
}
func StaticSet_DiagnosticGroups_TWEAKS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_TWEAKS = v
}

var static_DiagnosticGroups_MISSING_CALL_TO_SUPER *DiagnosticGroup

func StaticGet_DiagnosticGroups_MISSING_CALL_TO_SUPER() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MISSING_CALL_TO_SUPER
}
func StaticSet_DiagnosticGroups_MISSING_CALL_TO_SUPER(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MISSING_CALL_TO_SUPER = v
}

var static_DiagnosticGroups_MISSING_OVERRIDE *DiagnosticGroup

func StaticGet_DiagnosticGroups_MISSING_OVERRIDE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MISSING_OVERRIDE
}
func StaticSet_DiagnosticGroups_MISSING_OVERRIDE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MISSING_OVERRIDE = v
}

var static_DiagnosticGroups_MISSING_PROPERTIES *DiagnosticGroup

func StaticGet_DiagnosticGroups_MISSING_PROPERTIES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MISSING_PROPERTIES
}
func StaticSet_DiagnosticGroups_MISSING_PROPERTIES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MISSING_PROPERTIES = v
}

var static_DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES *DiagnosticGroup

func StaticGet_DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES
}
func StaticSet_DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES = v
}

var static_DiagnosticGroups_J2CL_CHECKS *DiagnosticGroup

func StaticGet_DiagnosticGroups_J2CL_CHECKS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_J2CL_CHECKS
}
func StaticSet_DiagnosticGroups_J2CL_CHECKS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_J2CL_CHECKS = v
}

var static_DiagnosticGroups_MISSING_RETURN *DiagnosticGroup

func StaticGet_DiagnosticGroups_MISSING_RETURN() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MISSING_RETURN
}
func StaticSet_DiagnosticGroups_MISSING_RETURN(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MISSING_RETURN = v
}

var static_DiagnosticGroups_UNDEFINED_VARIABLES *DiagnosticGroup

func StaticGet_DiagnosticGroups_UNDEFINED_VARIABLES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_UNDEFINED_VARIABLES
}
func StaticSet_DiagnosticGroups_UNDEFINED_VARIABLES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_UNDEFINED_VARIABLES = v
}

var static_DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT *DiagnosticGroup

func StaticGet_DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT
}
func StaticSet_DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT = v
}

var static_DiagnosticGroups_CHECK_REGEXP *DiagnosticGroup

func StaticGet_DiagnosticGroups_CHECK_REGEXP() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CHECK_REGEXP
}
func StaticSet_DiagnosticGroups_CHECK_REGEXP(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CHECK_REGEXP = v
}

var static_DiagnosticGroups_CHECK_TYPES *DiagnosticGroup

func StaticGet_DiagnosticGroups_CHECK_TYPES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CHECK_TYPES
}
func StaticSet_DiagnosticGroups_CHECK_TYPES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CHECK_TYPES = v
}

var static_DiagnosticGroups_ES5_INHERITANCE_DIAGNOSTIC_GROUP *DiagnosticGroup

func StaticGet_DiagnosticGroups_ES5_INHERITANCE_DIAGNOSTIC_GROUP() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_ES5_INHERITANCE_DIAGNOSTIC_GROUP
}
func StaticSet_DiagnosticGroups_ES5_INHERITANCE_DIAGNOSTIC_GROUP(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_ES5_INHERITANCE_DIAGNOSTIC_GROUP = v
}

var static_DiagnosticGroups_CHECK_PROTOTYPAL_TYPES *DiagnosticGroup

func StaticGet_DiagnosticGroups_CHECK_PROTOTYPAL_TYPES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CHECK_PROTOTYPAL_TYPES
}
func StaticSet_DiagnosticGroups_CHECK_PROTOTYPAL_TYPES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CHECK_PROTOTYPAL_TYPES = v
}

var static_DiagnosticGroups_CHECK_STATIC_OVERRIDES *DiagnosticGroup

func StaticGet_DiagnosticGroups_CHECK_STATIC_OVERRIDES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CHECK_STATIC_OVERRIDES
}
func StaticSet_DiagnosticGroups_CHECK_STATIC_OVERRIDES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CHECK_STATIC_OVERRIDES = v
}

var static_DiagnosticGroups_TOO_MANY_TYPE_PARAMS *DiagnosticGroup

func StaticGet_DiagnosticGroups_TOO_MANY_TYPE_PARAMS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_TOO_MANY_TYPE_PARAMS
}
func StaticSet_DiagnosticGroups_TOO_MANY_TYPE_PARAMS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_TOO_MANY_TYPE_PARAMS = v
}

var static_DiagnosticGroups_STRICT_MISSING_PROPERTIES *DiagnosticGroup

func StaticGet_DiagnosticGroups_STRICT_MISSING_PROPERTIES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_STRICT_MISSING_PROPERTIES
}
func StaticSet_DiagnosticGroups_STRICT_MISSING_PROPERTIES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_STRICT_MISSING_PROPERTIES = v
}

var static_DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS *DiagnosticGroup

func StaticGet_DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS
}
func StaticSet_DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS = v
}

var static_DiagnosticGroups_STRICT_CHECK_TYPES *DiagnosticGroup

func StaticGet_DiagnosticGroups_STRICT_CHECK_TYPES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_STRICT_CHECK_TYPES
}
func StaticSet_DiagnosticGroups_STRICT_CHECK_TYPES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_STRICT_CHECK_TYPES = v
}

var static_DiagnosticGroups_REPORT_UNKNOWN_TYPES *DiagnosticGroup

func StaticGet_DiagnosticGroups_REPORT_UNKNOWN_TYPES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_REPORT_UNKNOWN_TYPES
}
func StaticSet_DiagnosticGroups_REPORT_UNKNOWN_TYPES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_REPORT_UNKNOWN_TYPES = v
}

var static_DiagnosticGroups_CHECK_VARIABLES *DiagnosticGroup

func StaticGet_DiagnosticGroups_CHECK_VARIABLES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CHECK_VARIABLES
}
func StaticSet_DiagnosticGroups_CHECK_VARIABLES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CHECK_VARIABLES = v
}

var static_DiagnosticGroups_CHECK_USELESS_CODE *DiagnosticGroup

func StaticGet_DiagnosticGroups_CHECK_USELESS_CODE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CHECK_USELESS_CODE
}
func StaticSet_DiagnosticGroups_CHECK_USELESS_CODE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CHECK_USELESS_CODE = v
}

var static_DiagnosticGroups_CONST *DiagnosticGroup

func StaticGet_DiagnosticGroups_CONST() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CONST
}
func StaticSet_DiagnosticGroups_CONST(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CONST = v
}

var static_DiagnosticGroups_CONSTANT_PROPERTY *DiagnosticGroup

func StaticGet_DiagnosticGroups_CONSTANT_PROPERTY() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CONSTANT_PROPERTY
}
func StaticSet_DiagnosticGroups_CONSTANT_PROPERTY(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CONSTANT_PROPERTY = v
}

var static_DiagnosticGroups_ACCESS_CONTROLS_CONST *DiagnosticGroup

func StaticGet_DiagnosticGroups_ACCESS_CONTROLS_CONST() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_ACCESS_CONTROLS_CONST
}
func StaticSet_DiagnosticGroups_ACCESS_CONTROLS_CONST(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_ACCESS_CONTROLS_CONST = v
}

var static_DiagnosticGroups_TYPE_INVALIDATION *DiagnosticGroup

func StaticGet_DiagnosticGroups_TYPE_INVALIDATION() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_TYPE_INVALIDATION
}
func StaticSet_DiagnosticGroups_TYPE_INVALIDATION(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_TYPE_INVALIDATION = v
}

var static_DiagnosticGroups_DUPLICATE_VARS *DiagnosticGroup

func StaticGet_DiagnosticGroups_DUPLICATE_VARS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_DUPLICATE_VARS
}
func StaticSet_DiagnosticGroups_DUPLICATE_VARS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_DUPLICATE_VARS = v
}

var static_DiagnosticGroups_ES5_STRICT *DiagnosticGroup

func StaticGet_DiagnosticGroups_ES5_STRICT() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_ES5_STRICT
}
func StaticSet_DiagnosticGroups_ES5_STRICT(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_ES5_STRICT = v
}

var static_DiagnosticGroups_WEAK_MODULE_GET *DiagnosticGroup

func StaticGet_DiagnosticGroups_WEAK_MODULE_GET() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_WEAK_MODULE_GET
}
func StaticSet_DiagnosticGroups_WEAK_MODULE_GET(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_WEAK_MODULE_GET = v
}

var static_DiagnosticGroups_MISSING_PROVIDE *DiagnosticGroup

func StaticGet_DiagnosticGroups_MISSING_PROVIDE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MISSING_PROVIDE
}
func StaticSet_DiagnosticGroups_MISSING_PROVIDE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MISSING_PROVIDE = v
}

var static_DiagnosticGroups_UNRECOGNIZED_TYPE_ERROR *DiagnosticGroup

func StaticGet_DiagnosticGroups_UNRECOGNIZED_TYPE_ERROR() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_UNRECOGNIZED_TYPE_ERROR
}
func StaticSet_DiagnosticGroups_UNRECOGNIZED_TYPE_ERROR(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_UNRECOGNIZED_TYPE_ERROR = v
}

var static_DiagnosticGroups_DANGEROUS_UNRECOGNIZED_TYPE_ERROR *DiagnosticGroup

func StaticGet_DiagnosticGroups_DANGEROUS_UNRECOGNIZED_TYPE_ERROR() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_DANGEROUS_UNRECOGNIZED_TYPE_ERROR
}
func StaticSet_DiagnosticGroups_DANGEROUS_UNRECOGNIZED_TYPE_ERROR(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_DANGEROUS_UNRECOGNIZED_TYPE_ERROR = v
}

var static_DiagnosticGroups_MISSING_REQUIRE *DiagnosticGroup

func StaticGet_DiagnosticGroups_MISSING_REQUIRE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MISSING_REQUIRE
}
func StaticSet_DiagnosticGroups_MISSING_REQUIRE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MISSING_REQUIRE = v
}

var static_DiagnosticGroups_MISSING_SOURCES_WARNINGS *DiagnosticGroup

func StaticGet_DiagnosticGroups_MISSING_SOURCES_WARNINGS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MISSING_SOURCES_WARNINGS
}
func StaticSet_DiagnosticGroups_MISSING_SOURCES_WARNINGS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MISSING_SOURCES_WARNINGS = v
}

var static_DiagnosticGroups_EXTRA_REQUIRE *DiagnosticGroup

func StaticGet_DiagnosticGroups_EXTRA_REQUIRE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_EXTRA_REQUIRE
}
func StaticSet_DiagnosticGroups_EXTRA_REQUIRE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_EXTRA_REQUIRE = v
}

var static_DiagnosticGroups_DUPLICATE_MESSAGE *DiagnosticGroup

func StaticGet_DiagnosticGroups_DUPLICATE_MESSAGE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_DUPLICATE_MESSAGE
}
func StaticSet_DiagnosticGroups_DUPLICATE_MESSAGE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_DUPLICATE_MESSAGE = v
}

var static_DiagnosticGroups_MESSAGE_DESCRIPTIONS *DiagnosticGroup

func StaticGet_DiagnosticGroups_MESSAGE_DESCRIPTIONS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MESSAGE_DESCRIPTIONS
}
func StaticSet_DiagnosticGroups_MESSAGE_DESCRIPTIONS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MESSAGE_DESCRIPTIONS = v
}

var static_DiagnosticGroups_MSG_CONVENTIONS *DiagnosticGroup

func StaticGet_DiagnosticGroups_MSG_CONVENTIONS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MSG_CONVENTIONS
}
func StaticSet_DiagnosticGroups_MSG_CONVENTIONS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MSG_CONVENTIONS = v
}

var static_DiagnosticGroups_MISPLACED_TYPE_ANNOTATION *DiagnosticGroup

func StaticGet_DiagnosticGroups_MISPLACED_TYPE_ANNOTATION() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MISPLACED_TYPE_ANNOTATION
}
func StaticSet_DiagnosticGroups_MISPLACED_TYPE_ANNOTATION(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MISPLACED_TYPE_ANNOTATION = v
}

var static_DiagnosticGroups_MISPLACED_MSG_ANNOTATION *DiagnosticGroup

func StaticGet_DiagnosticGroups_MISPLACED_MSG_ANNOTATION() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MISPLACED_MSG_ANNOTATION
}
func StaticSet_DiagnosticGroups_MISPLACED_MSG_ANNOTATION(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MISPLACED_MSG_ANNOTATION = v
}

var static_DiagnosticGroups_MISPLACED_SUPPRESS *DiagnosticGroup

func StaticGet_DiagnosticGroups_MISPLACED_SUPPRESS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MISPLACED_SUPPRESS
}
func StaticSet_DiagnosticGroups_MISPLACED_SUPPRESS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MISPLACED_SUPPRESS = v
}

var static_DiagnosticGroups_SUSPICIOUS_CODE *DiagnosticGroup

func StaticGet_DiagnosticGroups_SUSPICIOUS_CODE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_SUSPICIOUS_CODE
}
func StaticSet_DiagnosticGroups_SUSPICIOUS_CODE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_SUSPICIOUS_CODE = v
}

var static_DiagnosticGroups_FUNCTION_PARAMS *DiagnosticGroup

func StaticGet_DiagnosticGroups_FUNCTION_PARAMS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_FUNCTION_PARAMS
}
func StaticSet_DiagnosticGroups_FUNCTION_PARAMS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_FUNCTION_PARAMS = v
}

var static_DiagnosticGroups_DEPRECATED_ANNOTATIONS *DiagnosticGroup

func StaticGet_DiagnosticGroups_DEPRECATED_ANNOTATIONS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_DEPRECATED_ANNOTATIONS
}
func StaticSet_DiagnosticGroups_DEPRECATED_ANNOTATIONS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_DEPRECATED_ANNOTATIONS = v
}

var static_DiagnosticGroups_UNUSED_LOCAL_VARIABLE *DiagnosticGroup

func StaticGet_DiagnosticGroups_UNUSED_LOCAL_VARIABLE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_UNUSED_LOCAL_VARIABLE
}
func StaticSet_DiagnosticGroups_UNUSED_LOCAL_VARIABLE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_UNUSED_LOCAL_VARIABLE = v
}

var static_DiagnosticGroups_JSDOC_MISSING_TYPE *DiagnosticGroup

func StaticGet_DiagnosticGroups_JSDOC_MISSING_TYPE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_JSDOC_MISSING_TYPE
}
func StaticSet_DiagnosticGroups_JSDOC_MISSING_TYPE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_JSDOC_MISSING_TYPE = v
}

var static_DiagnosticGroups_TYPE_IMPORT_CODE_REFERENCES *DiagnosticGroup

func StaticGet_DiagnosticGroups_TYPE_IMPORT_CODE_REFERENCES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_TYPE_IMPORT_CODE_REFERENCES
}
func StaticSet_DiagnosticGroups_TYPE_IMPORT_CODE_REFERENCES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_TYPE_IMPORT_CODE_REFERENCES = v
}

var static_DiagnosticGroups_PARTIAL_ALIAS *DiagnosticGroup

func StaticGet_DiagnosticGroups_PARTIAL_ALIAS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_PARTIAL_ALIAS
}
func StaticSet_DiagnosticGroups_PARTIAL_ALIAS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_PARTIAL_ALIAS = v
}

var static_DiagnosticGroups_USE_OF_GOOG_PROVIDE *DiagnosticGroup

func StaticGet_DiagnosticGroups_USE_OF_GOOG_PROVIDE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_USE_OF_GOOG_PROVIDE
}
func StaticSet_DiagnosticGroups_USE_OF_GOOG_PROVIDE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_USE_OF_GOOG_PROVIDE = v
}

var static_DiagnosticGroups_LINT_VAR_DECLARATIONS *DiagnosticGroup

func StaticGet_DiagnosticGroups_LINT_VAR_DECLARATIONS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_LINT_VAR_DECLARATIONS
}
func StaticSet_DiagnosticGroups_LINT_VAR_DECLARATIONS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_LINT_VAR_DECLARATIONS = v
}

var static_DiagnosticGroups_LINT_CHECKS *DiagnosticGroup

func StaticGet_DiagnosticGroups_LINT_CHECKS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_LINT_CHECKS
}
func StaticSet_DiagnosticGroups_LINT_CHECKS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_LINT_CHECKS = v
}

var static_DiagnosticGroups_STRICT_MODULE_CHECKS *DiagnosticGroup

func StaticGet_DiagnosticGroups_STRICT_MODULE_CHECKS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_STRICT_MODULE_CHECKS
}
func StaticSet_DiagnosticGroups_STRICT_MODULE_CHECKS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_STRICT_MODULE_CHECKS = v
}

var static_DiagnosticGroups_ANALYZER_CHECKS *DiagnosticGroup

func StaticGet_DiagnosticGroups_ANALYZER_CHECKS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_ANALYZER_CHECKS
}
func StaticSet_DiagnosticGroups_ANALYZER_CHECKS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_ANALYZER_CHECKS = v
}

var static_DiagnosticGroups_CLOSURE_DEP_METHOD_USAGE_CHECKS *DiagnosticGroup

func StaticGet_DiagnosticGroups_CLOSURE_DEP_METHOD_USAGE_CHECKS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CLOSURE_DEP_METHOD_USAGE_CHECKS
}
func StaticSet_DiagnosticGroups_CLOSURE_DEP_METHOD_USAGE_CHECKS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CLOSURE_DEP_METHOD_USAGE_CHECKS = v
}

var static_DiagnosticGroups_CLOSURE_CLASS_CHECKS *DiagnosticGroup

func StaticGet_DiagnosticGroups_CLOSURE_CLASS_CHECKS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CLOSURE_CLASS_CHECKS
}
func StaticSet_DiagnosticGroups_CLOSURE_CLASS_CHECKS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CLOSURE_CLASS_CHECKS = v
}

var static_DiagnosticGroups_MALFORMED_GOOG_MODULE *DiagnosticGroup

func StaticGet_DiagnosticGroups_MALFORMED_GOOG_MODULE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MALFORMED_GOOG_MODULE
}
func StaticSet_DiagnosticGroups_MALFORMED_GOOG_MODULE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MALFORMED_GOOG_MODULE = v
}

var static_DiagnosticGroups_CONFORMANCE_VIOLATIONS *DiagnosticGroup

func StaticGet_DiagnosticGroups_CONFORMANCE_VIOLATIONS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CONFORMANCE_VIOLATIONS
}
func StaticSet_DiagnosticGroups_CONFORMANCE_VIOLATIONS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CONFORMANCE_VIOLATIONS = v
}

var static_DiagnosticGroups_LATE_PROVIDE *DiagnosticGroup

func StaticGet_DiagnosticGroups_LATE_PROVIDE() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_LATE_PROVIDE
}
func StaticSet_DiagnosticGroups_LATE_PROVIDE(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_LATE_PROVIDE = v
}

var static_DiagnosticGroups_DUPLICATE_NAMESPACES *DiagnosticGroup

func StaticGet_DiagnosticGroups_DUPLICATE_NAMESPACES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_DUPLICATE_NAMESPACES
}
func StaticSet_DiagnosticGroups_DUPLICATE_NAMESPACES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_DUPLICATE_NAMESPACES = v
}

var static_DiagnosticGroups_INVALID_DEFINES *DiagnosticGroup

func StaticGet_DiagnosticGroups_INVALID_DEFINES() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_INVALID_DEFINES
}
func StaticSet_DiagnosticGroups_INVALID_DEFINES(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_INVALID_DEFINES = v
}

var static_DiagnosticGroups_MISSING_POLYFILL *DiagnosticGroup

func StaticGet_DiagnosticGroups_MISSING_POLYFILL() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_MISSING_POLYFILL
}
func StaticSet_DiagnosticGroups_MISSING_POLYFILL(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_MISSING_POLYFILL = v
}

var static_DiagnosticGroups_POLYMER *DiagnosticGroup

func StaticGet_DiagnosticGroups_POLYMER() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_POLYMER
}
func StaticSet_DiagnosticGroups_POLYMER(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_POLYMER = v
}

var static_DiagnosticGroups_BOUNDED_GENERICS *DiagnosticGroup

func StaticGet_DiagnosticGroups_BOUNDED_GENERICS() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_BOUNDED_GENERICS
}
func StaticSet_DiagnosticGroups_BOUNDED_GENERICS(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_BOUNDED_GENERICS = v
}

var static_DiagnosticGroups_PARSING *DiagnosticGroup

func StaticGet_DiagnosticGroups_PARSING() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_PARSING
}
func StaticSet_DiagnosticGroups_PARSING(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_PARSING = v
}

var static_DiagnosticGroups_CLOSURE_UNAWARE_CODE_ANNOTATION_PRESENT *DiagnosticGroup

func StaticGet_DiagnosticGroups_CLOSURE_UNAWARE_CODE_ANNOTATION_PRESENT() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_CLOSURE_UNAWARE_CODE_ANNOTATION_PRESENT
}
func StaticSet_DiagnosticGroups_CLOSURE_UNAWARE_CODE_ANNOTATION_PRESENT(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_CLOSURE_UNAWARE_CODE_ANNOTATION_PRESENT = v
}

var static_DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION *DiagnosticGroup

func StaticGet_DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION() *DiagnosticGroup {
	clinit_DiagnosticGroups()
	return static_DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION
}
func StaticSet_DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION(v *DiagnosticGroup) {
	clinit_DiagnosticGroups()
	static_DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION = v
}

var clinitDone_DiagnosticGroups uint32
var clinitState_DiagnosticGroups j2gort.ClinitState

func clinit_DiagnosticGroups() {
	if atomic.LoadUint32(&clinitDone_DiagnosticGroups) != 0 {
		return
	}
	clinitSlow_DiagnosticGroups()
}

func clinitSlow_DiagnosticGroups() {
	if !j2gort.EnterClinit(&clinitState_DiagnosticGroups, &clinitDone_DiagnosticGroups) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DiagnosticGroups, &clinitDone_DiagnosticGroups)
	static_DiagnosticGroups_UNUSED = DiagnosticType_Warning("JSC_UNUSED", "{0}")
	static_DiagnosticGroups_wildcardExcludedGroups = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("reportUnknownTypes", "analyzerChecks", "missingSourcesWarnings", "closureUnawareCodeAnnotationPresent")
	static_DiagnosticGroups_groupsByName = j2gort.NewExt_java_util_LinkedHashMap()
	static_DiagnosticGroups_UNTRANSPILABLE_FEATURES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("untranspilableFeatures", []*DiagnosticType{StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT()})
	static_DiagnosticGroups_MODULE_LOAD = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("moduleLoad", []*DiagnosticType{StaticGet_ModuleLoader_LOAD_WARNING(), StaticGet_ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT(), StaticGet_ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE()})
	static_DiagnosticGroups_MODULE_IMPORT = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("moduleImport", []*DiagnosticType{StaticGet_ModuleMapCreator_DOES_NOT_HAVE_EXPORT(), StaticGet_ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS()})
	static_DiagnosticGroups_GLOBAL_THIS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("globalThis", []*DiagnosticType{StaticGet_CheckGlobalThis_GLOBAL_THIS()})
	static_DiagnosticGroups_DEPRECATED = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("deprecated", []*DiagnosticType{StaticGet_CheckAccessControls_DEPRECATED_NAME(), StaticGet_CheckAccessControls_DEPRECATED_NAME_REASON(), StaticGet_CheckAccessControls_DEPRECATED_PROP(), StaticGet_CheckAccessControls_DEPRECATED_PROP_REASON(), StaticGet_CheckAccessControls_DEPRECATED_CLASS(), StaticGet_CheckAccessControls_DEPRECATED_CLASS_REASON()})
	static_DiagnosticGroups_UNDERSCORE = DiagnosticGroups_RegisterDeprecatedGroup("underscore")
	static_DiagnosticGroups_VISIBILITY = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("visibility", []*DiagnosticType{StaticGet_CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS(), StaticGet_CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS(), StaticGet_CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS(), StaticGet_CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS(), StaticGet_CheckAccessControls_EXTEND_FINAL_CLASS(), StaticGet_CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN(), StaticGet_CheckAccessControls_PRIVATE_OVERRIDE(), StaticGet_CheckAccessControls_VISIBILITY_MISMATCH()})
	static_DiagnosticGroups_ACCESS_CONTROLS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticGroup("accessControls", StaticGet_DiagnosticGroups_VISIBILITY())
	static_DiagnosticGroups_NON_STANDARD_JSDOC = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("nonStandardJsDocs", []*DiagnosticType{StaticGet_RhinoErrorReporter_BAD_JSDOC_ANNOTATION(), StaticGet_RhinoErrorReporter_INVALID_PARAM(), StaticGet_RhinoErrorReporter_JSDOC_IMPORT_TYPE_WARNING(), StaticGet_CheckJSDoc_JSDOC_IN_BLOCK_COMMENT()})
	static_DiagnosticGroups_INVALID_CASTS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("invalidCasts", []*DiagnosticType{StaticGet_TypeValidator_INVALID_CAST()})
	static_DiagnosticGroups_STRICT_MODULE_DEP_CHECK = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("strictModuleDepCheck", []*DiagnosticType{StaticGet_VarCheck_STRICT_CHUNK_DEP_ERROR(), StaticGet_CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR()})
	static_DiagnosticGroups_VIOLATED_MODULE_DEP = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("violatedModuleDep", []*DiagnosticType{StaticGet_VarCheck_VIOLATED_CHUNK_DEP_ERROR()})
	static_DiagnosticGroups_EXTERNS_VALIDATION = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("externsValidation", []*DiagnosticType{StaticGet_VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR(), StaticGet_VarCheck_UNDEFINED_EXTERN_VAR_ERROR()})
	static_DiagnosticGroups_UNKNOWN_DEFINES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("unknownDefines", []*DiagnosticType{StaticGet_ProcessDefines_UNKNOWN_DEFINE_WARNING()})
	static_DiagnosticGroups_DEFINE_WITHOUT_GOOG_DEFINE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("defineWithoutGoogDefine", []*DiagnosticType{StaticGet_ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE()})
	static_DiagnosticGroups_TWEAKS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("tweakValidation", []*DiagnosticType{StaticGet_ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING(), StaticGet_ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING()})
	static_DiagnosticGroups_MISSING_CALL_TO_SUPER = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("missingCallToSuper", []*DiagnosticType{StaticGet_CheckSuper_RETURN_WITHOUT_CALL_TO_SUPER(), StaticGet_CheckSuper_MISSING_ROOT_LEVEL_CALL_TO_SUPER()})
	static_DiagnosticGroups_MISSING_OVERRIDE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("missingOverride", []*DiagnosticType{StaticGet_TypeCheck_HIDDEN_INTERFACE_PROPERTY(), StaticGet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY(), StaticGet_TypeCheck_HIDDEN_SUPERCLASS_PROPERTY()})
	static_DiagnosticGroups_MISSING_PROPERTIES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("missingProperties", []*DiagnosticType{StaticGet_TypeCheck_INEXISTENT_PROPERTY(), StaticGet_TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION(), StaticGet_TypeCheck_POSSIBLE_INEXISTENT_PROPERTY()})
	static_DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("globallyMissingProperties", []*DiagnosticType{StaticGet_TypeCheck_POSSIBLE_INEXISTENT_PROPERTY()})
	static_DiagnosticGroups_J2CL_CHECKS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("j2clChecks", []*DiagnosticType{StaticGet_J2clChecksPass_J2CL_REFERENCE_EQUALITY()})
	static_DiagnosticGroups_MISSING_RETURN = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("missingReturn", []*DiagnosticType{StaticGet_CheckMissingReturn_MISSING_RETURN_STATEMENT()})
	static_DiagnosticGroups_UNDEFINED_VARIABLES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("undefinedVars", []*DiagnosticType{StaticGet_VarCheck_UNDEFINED_VAR_ERROR()})
	static_DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("checkDebuggerStatement", []*DiagnosticType{StaticGet_CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT()})
	static_DiagnosticGroups_CHECK_REGEXP = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("checkRegExp", []*DiagnosticType{StaticGet_CheckRegExp_REGEXP_REFERENCE(), StaticGet_CheckRegExp_MALFORMED_REGEXP()})
	static_DiagnosticGroups_CHECK_TYPES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticGroup_2("checkTypes", []*DiagnosticGroup{StaticGet_TypeValidator_ALL_DIAGNOSTICS(), StaticGet_TypeCheck_ALL_DIAGNOSTICS(), StaticGet_FunctionTypeBuilder_ALL_DIAGNOSTICS(), StaticGet_DiagnosticGroups_GLOBAL_THIS(), NewDiagnosticGroup_1_DiagnosticType_2([]*DiagnosticType{StaticGet_CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL(), StaticGet_CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION()})})
	static_DiagnosticGroups_ES5_INHERITANCE_DIAGNOSTIC_GROUP = DiagnosticGroups_RegisterGroup_2_string_DiagnosticGroup("checkEs5InheritanceCorrectnessConditions", StaticGet_TypeCheck_ES5_INHERITANCE_DIAGNOSTIC_GROUP())
	static_DiagnosticGroups_CHECK_PROTOTYPAL_TYPES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("checkPrototypalTypes", []*DiagnosticType{StaticGet_TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE(), StaticGet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY(), StaticGet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH()})
	static_DiagnosticGroups_CHECK_STATIC_OVERRIDES = StaticGet_DiagnosticGroups_CHECK_PROTOTYPAL_TYPES()
	static_DiagnosticGroups_TOO_MANY_TYPE_PARAMS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("tooManyTypeParams", []*DiagnosticType{StaticGet_RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS()})
	static_DiagnosticGroups_STRICT_MISSING_PROPERTIES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("strictMissingProperties", []*DiagnosticType{StaticGet_TypeCheck_STRICT_INEXISTENT_PROPERTY(), StaticGet_TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION(), StaticGet_TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY(), StaticGet_TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE()})
	static_DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("strictPrimitiveOperators", []*DiagnosticType{StaticGet_TypeValidator_INVALID_OPERAND_TYPE()})
	static_DiagnosticGroups_STRICT_CHECK_TYPES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticGroup_2("strictCheckTypes", []*DiagnosticGroup{StaticGet_DiagnosticGroups_STRICT_MISSING_PROPERTIES(), StaticGet_DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS()})
	static_DiagnosticGroups_REPORT_UNKNOWN_TYPES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("reportUnknownTypes", []*DiagnosticType{StaticGet_TypeCheck_UNKNOWN_EXPR_TYPE()})
	static_DiagnosticGroups_CHECK_VARIABLES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("checkVars", []*DiagnosticType{StaticGet_VarCheck_UNDEFINED_VAR_ERROR(), StaticGet_VarCheck_VAR_MULTIPLY_DECLARED_ERROR(), StaticGet_VariableReferenceCheck_EARLY_REFERENCE(), StaticGet_VariableReferenceCheck_REDECLARED_VARIABLE()})
	static_DiagnosticGroups_CHECK_USELESS_CODE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("uselessCode", []*DiagnosticType{StaticGet_CheckSideEffects_USELESS_CODE_ERROR(), StaticGet_CheckUnreachableCode_UNREACHABLE_CODE()})
	static_DiagnosticGroups_CONST = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("const", []*DiagnosticType{StaticGet_CheckAccessControls_CONST_PROPERTY_DELETED(), StaticGet_CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE(), StaticGet_ConstCheck_CONST_REASSIGNED_VALUE_ERROR()})
	static_DiagnosticGroups_CONSTANT_PROPERTY = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("constantProperty", []*DiagnosticType{StaticGet_CheckAccessControls_CONST_PROPERTY_DELETED(), StaticGet_CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE()})
	static_DiagnosticGroups_ACCESS_CONTROLS_CONST = DiagnosticGroups_RegisterGroup_2_string_DiagnosticGroup("accessControlsConst", StaticGet_DiagnosticGroups_CONSTANT_PROPERTY())
	static_DiagnosticGroups_TYPE_INVALIDATION = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("typeInvalidation", []*DiagnosticType{StaticGet_DisambiguateProperties_PROPERTY_INVALIDATION()})
	static_DiagnosticGroups_DUPLICATE_VARS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("duplicate", []*DiagnosticType{StaticGet_InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING(), StaticGet_TypeValidator_DUP_VAR_DECLARATION(), StaticGet_TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH(), StaticGet_TypeCheck_FUNCTION_MASKS_VARIABLE(), StaticGet_VarCheck_VAR_MULTIPLY_DECLARED_ERROR(), StaticGet_VariableReferenceCheck_REDECLARED_VARIABLE()})
	static_DiagnosticGroups_ES5_STRICT = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("es5Strict", []*DiagnosticType{StaticGet_RhinoErrorReporter_INVALID_OCTAL_LITERAL(), StaticGet_RhinoErrorReporter_DUPLICATE_PARAM(), StaticGet_StrictModeCheck_ARGUMENTS_ASSIGNMENT(), StaticGet_StrictModeCheck_ARGUMENTS_DECLARATION(), StaticGet_StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN(), StaticGet_StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN(), StaticGet_StrictModeCheck_DELETE_VARIABLE(), StaticGet_StrictModeCheck_DUPLICATE_MEMBER(), StaticGet_StrictModeCheck_EVAL_ASSIGNMENT(), StaticGet_StrictModeCheck_EVAL_DECLARATION(), StaticGet_StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN(), StaticGet_StrictModeCheck_FUNCTION_CALLER_FORBIDDEN(), StaticGet_StrictModeCheck_USE_OF_WITH()})
	static_DiagnosticGroups_WEAK_MODULE_GET = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("weakModuleGet", []*DiagnosticType{StaticGet_ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE()})
	static_DiagnosticGroups_MISSING_PROVIDE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("missingProvide", []*DiagnosticType{StaticGet_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE()})
	static_DiagnosticGroups_UNRECOGNIZED_TYPE_ERROR = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("unrecognizedTypeError", []*DiagnosticType{StaticGet_RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR()})
	static_DiagnosticGroups_DANGEROUS_UNRECOGNIZED_TYPE_ERROR = DiagnosticGroups_RegisterGroup_2_string_DiagnosticGroup("dangerousUnrecognizedTypeError", StaticGet_DiagnosticGroups_UNRECOGNIZED_TYPE_ERROR())
	static_DiagnosticGroups_MISSING_REQUIRE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("missingRequire", []*DiagnosticType{StaticGet_CheckMissingRequires_MISSING_REQUIRE(), StaticGet_CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE(), StaticGet_CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE(), StaticGet_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE(), StaticGet_CheckMissingRequires_MISSING_REQUIRE_TYPE(), StaticGet_CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE(), StaticGet_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE(), StaticGet_CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE(), StaticGet_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE(), StaticGet_CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE()})
	static_DiagnosticGroups_MISSING_SOURCES_WARNINGS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticGroup_2("missingSourcesWarnings", []*DiagnosticGroup{StaticGet_DiagnosticGroups_REPORT_UNKNOWN_TYPES(), StaticGet_DiagnosticGroups_UNDEFINED_VARIABLES(), StaticGet_DiagnosticGroups_MISSING_PROVIDE(), DiagnosticGroup_ForType(StaticGet_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE()), StaticGet_DiagnosticGroups_MISSING_PROPERTIES(), StaticGet_DiagnosticGroups_DUPLICATE_VARS(), DiagnosticGroup_ForType(StaticGet_ProcessDefines_INVALID_DEFINE_VALUE())})
	static_DiagnosticGroups_EXTRA_REQUIRE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("extraRequire", []*DiagnosticType{StaticGet_CheckExtraRequires_EXTRA_REQUIRE_WARNING()})
	static_DiagnosticGroups_DUPLICATE_MESSAGE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("duplicateMessage", []*DiagnosticType{StaticGet_JsMessageVisitor_MESSAGE_DUPLICATE_KEY()})
	static_DiagnosticGroups_MESSAGE_DESCRIPTIONS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("msgDescriptions", []*DiagnosticType{StaticGet_JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION()})
	static_DiagnosticGroups_MSG_CONVENTIONS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("messageConventions", []*DiagnosticType{StaticGet_JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION(), StaticGet_JsMessageVisitor_MESSAGE_HAS_NO_TEXT(), StaticGet_JsMessageVisitor_MESSAGE_TREE_MALFORMED(), StaticGet_JsMessageVisitor_MESSAGE_HAS_NO_VALUE(), StaticGet_JsMessageVisitor_MESSAGE_DUPLICATE_KEY(), StaticGet_JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY()})
	static_DiagnosticGroups_MISPLACED_TYPE_ANNOTATION = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("misplacedTypeAnnotation", []*DiagnosticType{StaticGet_CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR(), StaticGet_CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION(), StaticGet_CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL(), StaticGet_CheckJSDoc_DISALLOWED_MEMBER_JSDOC(), StaticGet_CheckJSDoc_INVALID_NO_SIDE_EFFECT_ANNOTATION(), StaticGet_CheckJSDoc_INVALID_MODIFIES_ANNOTATION(), StaticGet_CheckJSDoc_JSDOC_ON_RETURN(), StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), StaticGet_CheckJSDoc_MISPLACED_MSG_ANNOTATION()})
	static_DiagnosticGroups_MISPLACED_MSG_ANNOTATION = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("misplacedMsgAnnotation", []*DiagnosticType{StaticGet_CheckJSDoc_MISPLACED_MSG_ANNOTATION()})
	static_DiagnosticGroups_MISPLACED_SUPPRESS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("misplacedSuppress", []*DiagnosticType{StaticGet_CheckJSDoc_MISPLACED_SUPPRESS()})
	static_DiagnosticGroups_SUSPICIOUS_CODE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("suspiciousCode", []*DiagnosticType{StaticGet_CheckDuplicateCase_DUPLICATE_CASE(), StaticGet_CheckSuspiciousCode_SUSPICIOUS_SEMICOLON(), StaticGet_CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN(), StaticGet_CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN(), StaticGet_CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR(), StaticGet_CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND(), StaticGet_CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR(), StaticGet_CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR(), StaticGet_ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT(), StaticGet_TypeCheck_DETERMINISTIC_TEST()})
	static_DiagnosticGroups_FUNCTION_PARAMS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("functionParams", []*DiagnosticType{StaticGet_FunctionTypeBuilder_INEXISTENT_PARAM(), StaticGet_FunctionTypeBuilder_OPTIONAL_ARG_AT_END()})
	static_DiagnosticGroups_DEPRECATED_ANNOTATIONS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("deprecatedAnnotations", []*DiagnosticType{StaticGet_CheckJSDoc_ANNOTATION_DEPRECATED()})
	static_DiagnosticGroups_UNUSED_LOCAL_VARIABLE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("unusedLocalVariables", []*DiagnosticType{StaticGet_VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT()})
	static_DiagnosticGroups_JSDOC_MISSING_TYPE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("jsdocMissingType", []*DiagnosticType{StaticGet_RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING()})
	static_DiagnosticGroups_TYPE_IMPORT_CODE_REFERENCES = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("typeImportCodeReferences", []*DiagnosticType{StaticGet_CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE()})
	static_DiagnosticGroups_PARTIAL_ALIAS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("partialAlias", []*DiagnosticType{StaticGet_InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING()})
	static_DiagnosticGroups_USE_OF_GOOG_PROVIDE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("useOfGoogProvide", []*DiagnosticType{StaticGet_ClosureCheckModule_USE_OF_GOOG_PROVIDE()})
	static_DiagnosticGroups_LINT_VAR_DECLARATIONS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("lintVarDeclarations", []*DiagnosticType{StaticGet_CheckVar_VAR()})
	static_DiagnosticGroups_LINT_CHECKS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticGroup_2("lintChecks", []*DiagnosticGroup{StaticGet_CheckJSDocStyle_LINT_DIAGNOSTICS(), StaticGet_DiagnosticGroups_USE_OF_GOOG_PROVIDE(), StaticGet_DiagnosticGroups_LINT_VAR_DECLARATIONS(), NewDiagnosticGroup_1_DiagnosticType_2([]*DiagnosticType{StaticGet_CheckClosureImports_LET_CLOSURE_IMPORT(), StaticGet_CheckConstPrivateProperties_MISSING_CONST_PROPERTY(), StaticGet_CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME(), StaticGet_CheckConstantCaseNames_MISSING_CONST_PROPERTY(), StaticGet_CheckEmptyStatements_USELESS_EMPTY_STATEMENT(), StaticGet_CheckEnums_COMPUTED_PROP_NAME_IN_ENUM(), StaticGet_CheckEnums_DUPLICATE_ENUM_VALUE(), StaticGet_CheckEnums_ENUM_PROP_NOT_CONSTANT(), StaticGet_CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER(), StaticGet_CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM(), StaticGet_CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM(), StaticGet_CheckEs6ModuleFileStructure_MUST_COME_BEFORE(), StaticGet_CheckEs6Modules_DUPLICATE_IMPORT(), StaticGet_CheckEs6Modules_NO_DEFAULT_EXPORT(), StaticGet_CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE(), StaticGet_CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY(), StaticGet_CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS(), StaticGet_CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS(), StaticGet_CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE(), StaticGet_CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT(), StaticGet_CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS(), StaticGet_CheckMissingSemicolon_MISSING_SEMICOLON(), StaticGet_CheckNoMutatedEs6Exports_MUTATED_EXPORT(), StaticGet_CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC(), StaticGet_CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC(), StaticGet_CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC(), StaticGet_CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT(), StaticGet_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION(), StaticGet_CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER(), StaticGet_CheckProvidesSorted_PROVIDES_NOT_SORTED(), StaticGet_CheckRequiresSorted_REQUIRES_NOT_SORTED(), StaticGet_CheckUnusedLabels_UNUSED_LABEL(), StaticGet_CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY(), StaticGet_CheckUselessBlocks_USELESS_BLOCK(), StaticGet_ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE(), StaticGet_ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE(), StaticGet_ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION(), StaticGet_ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE(), StaticGet_RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING(), StaticGet_RhinoErrorReporter_UNNECESSARY_ESCAPE(), StaticGet_RhinoErrorReporter_STRING_CONTINUATION()})})
	static_DiagnosticGroups_STRICT_MODULE_CHECKS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("strictModuleChecks", []*DiagnosticType{StaticGet_ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE(), StaticGet_ClosureCheckModule_LET_GOOG_REQUIRE(), StaticGet_ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME(), StaticGet_ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME()})
	static_DiagnosticGroups_ANALYZER_CHECKS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("analyzerChecks", []*DiagnosticType{StaticGet_CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT(), StaticGet_ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC(), StaticGet_ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC(), StaticGet_CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE()})
	static_DiagnosticGroups_CLOSURE_DEP_METHOD_USAGE_CHECKS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("closureDepMethodUsageChecks", []*DiagnosticType{StaticGet_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR(), StaticGet_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR(), StaticGet_ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR(), StaticGet_ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE()})
	static_DiagnosticGroups_CLOSURE_CLASS_CHECKS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("closureClassChecks", []*DiagnosticType{StaticGet_ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR()})
	static_DiagnosticGroups_MALFORMED_GOOG_MODULE = DiagnosticGroups_RegisterUnsuppressibleGroup([]*DiagnosticType{StaticGet_ClosureCheckModule_GOOG_MODULE_MISPLACED(), StaticGet_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE()})
	DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("transitionalSuspiciousCodeWarnings", []*DiagnosticType{StaticGet_PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND()})
	static_DiagnosticGroups_CONFORMANCE_VIOLATIONS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("conformanceViolations", []*DiagnosticType{StaticGet_CheckConformance_CONFORMANCE_VIOLATION(), StaticGet_CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION()})
	static_DiagnosticGroups_LATE_PROVIDE = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("lateProvide", []*DiagnosticType{StaticGet_CheckClosureImports_LATE_PROVIDE_ERROR()})
	static_DiagnosticGroups_DUPLICATE_NAMESPACES = DiagnosticGroups_RegisterUnsuppressibleGroup([]*DiagnosticType{StaticGet_ClosurePrimitiveErrors_DUPLICATE_MODULE(), StaticGet_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE(), StaticGet_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE()})
	static_DiagnosticGroups_INVALID_DEFINES = DiagnosticGroups_RegisterUnsuppressibleGroup([]*DiagnosticType{StaticGet_ProcessDefines_INVALID_DEFINE_VALUE(), StaticGet_ProcessDefines_INVALID_DEFINE_TYPE()})
	static_DiagnosticGroups_MISSING_POLYFILL = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("missingPolyfill", []*DiagnosticType{StaticGet_RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR()})
	static_DiagnosticGroups_POLYMER = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("polymer", []*DiagnosticType{StaticGet_PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID()})
	static_DiagnosticGroups_BOUNDED_GENERICS = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("boundedGenerics", []*DiagnosticType{StaticGet_RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES(), StaticGet_RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR()})
	static_DiagnosticGroups_PARSING = DiagnosticGroup_ForType(StaticGet_RhinoErrorReporter_PARSE_ERROR())
	static_DiagnosticGroups_CLOSURE_UNAWARE_CODE_ANNOTATION_PRESENT = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("closureUnawareCodeAnnotationPresent", []*DiagnosticType{StaticGet_RhinoErrorReporter_CLOSURE_UNAWARE_ANNOTATION_PRESENT()})
	DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("polymerBehavior", []*DiagnosticType{StaticGet_PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR()})
	DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("invalidProvide", []*DiagnosticType{StaticGet_ProcessClosurePrimitives_INVALID_PROVIDE_ERROR()})
	DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("conflictingIjsFile", []*DiagnosticType{StaticGet_IjsErrors_CONFLICTING_IJS_FILE()})
	DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("implicitWeakEntryPoint", []*DiagnosticType{StaticGet_JSChunkGraph_IMPLICIT_WEAK_ENTRY_POINT_ERROR()})
	static_DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION = DiagnosticGroups_RegisterGroup_2_string_DiagnosticType("artificialFunctionPurityValidation", []*DiagnosticType{StaticGet_PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION()})
}

func (this_ *DiagnosticGroups) initFields_DiagnosticGroups() {
	_ = this_
}

// Java constructor: DiagnosticGroups()
func NewDiagnosticGroups() *DiagnosticGroups {
	clinit_DiagnosticGroups()
	this_ := &DiagnosticGroups{}
	InitDiagnosticGroups(this_)
	return this_
}

func InitDiagnosticGroups(this_ *DiagnosticGroups) {
	this_.initFields_DiagnosticGroups()
}

// Java static method: DiagnosticGroups.registerDeprecatedGroup(java.lang.String name)
func DiagnosticGroups_RegisterDeprecatedGroup(name string) *DiagnosticGroup {
	clinit_DiagnosticGroups()
	_ = name
	return DiagnosticGroups_RegisterGroup_2_string_DiagnosticGroup(name, NewDiagnosticGroup_2_string_DiagnosticType(name, []*DiagnosticType{StaticGet_DiagnosticGroups_UNUSED()}))
}

// Java static method: DiagnosticGroups.registerUnsuppressibleGroup(com.google.javascript.jscomp.DiagnosticType[] types)
func DiagnosticGroups_RegisterUnsuppressibleGroup(types []*DiagnosticType) *DiagnosticGroup {
	clinit_DiagnosticGroups()
	_ = types
	return NewDiagnosticGroup_1_DiagnosticType_2(types)
}

// Java static method: DiagnosticGroups.registerGroup(java.lang.String name, com.google.javascript.jscomp.DiagnosticGroup group)
func DiagnosticGroups_RegisterGroup_2_string_DiagnosticGroup(name string, group *DiagnosticGroup) *DiagnosticGroup {
	clinit_DiagnosticGroups()
	_ = name
	_ = group
	j2gort.MapPutStr(StaticGet_DiagnosticGroups_GroupsByName(), name, j2gort.PtrToAny(group))
	return group
}

// Java static method: DiagnosticGroups.registerGroup(java.lang.String name, com.google.javascript.jscomp.DiagnosticType[] types)
func DiagnosticGroups_RegisterGroup_2_string_DiagnosticType(name string, types []*DiagnosticType) *DiagnosticGroup {
	clinit_DiagnosticGroups()
	_ = name
	_ = types
	clinit_DiagnosticGroup()
	var val_1 DiagnosticGroup
	InitDiagnosticGroup_2_string_DiagnosticType(&val_1, name, types)
	var group *DiagnosticGroup = &val_1
	_ = group
	j2gort.MapPutStr(StaticGet_DiagnosticGroups_GroupsByName(), name, group)
	return group
}

// Java static method: DiagnosticGroups.registerGroup(java.lang.String name, com.google.javascript.jscomp.DiagnosticGroup[] groups)
func DiagnosticGroups_RegisterGroup_2_string_DiagnosticGroup_2(name string, groups []*DiagnosticGroup) *DiagnosticGroup {
	clinit_DiagnosticGroups()
	_ = name
	_ = groups
	clinit_DiagnosticGroup()
	var val_2 DiagnosticGroup
	InitDiagnosticGroup_2_string_DiagnosticGroup(&val_2, name, groups)
	var group *DiagnosticGroup = &val_2
	_ = group
	j2gort.MapPutStr(StaticGet_DiagnosticGroups_GroupsByName(), name, group)
	return group
}

// Java static method: DiagnosticGroups.getRegisteredGroups()
func DiagnosticGroups_GetRegisteredGroups() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_DiagnosticGroups()
	return j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(StaticGet_DiagnosticGroups_GroupsByName())
}

// Java static method: DiagnosticGroups.forName(java.lang.String name)
func DiagnosticGroups_ForName(name string) *DiagnosticGroup {
	clinit_DiagnosticGroups()
	_ = name
	return j2gort.CastPtr[DiagnosticGroup](j2gort.MapGetStr(StaticGet_DiagnosticGroups_GroupsByName(), name))
}

// Java method: DiagnosticGroups.setWarningLevel(com.google.javascript.jscomp.CompilerOptions options, java.lang.String name, com.google.javascript.jscomp.CheckLevel level)
func (this_ *DiagnosticGroups) SetWarningLevel(options *CompilerOptions, name string, level *CheckLevel) {
	_ = options
	_ = name
	_ = level
	var group *DiagnosticGroup = DiagnosticGroups_ForName(name)
	_ = group
	j2gort.CheckNotNullPtr(group)
	options.SetWarningLevel(group, level)
}
