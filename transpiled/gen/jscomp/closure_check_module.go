/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ClosureCheckModule.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ClosureCheckModule
type ClosureCheckModule struct {
	NodeTraversal_AbstractModuleCallbackBase
	currentModuleInfo *ClosureCheckModule_ModuleInfo `gson:"currentModuleInfo"`
}

func (this_ *ClosureCheckModule) BindSelf_(s any) {
	_ = s
	this_.NodeTraversal_AbstractModuleCallbackBase.BindSelf_(s)
}
func (this_ *ClosureCheckModule) IsJavaNil_() bool { return this_ == nil }
func (*ClosureCheckModule) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureCheckModule", "ClosureCheckModule"
}
func (this_ *ClosureCheckModule) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractModuleCallbackBase.IdentityHash_()
}
func (*ClosureCheckModule) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureCheckModule", "com.google.javascript.jscomp.ClosureCheckModule", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureCheckModule](v_)
	})
}
func (*ClosureCheckModule) Implements_jscomp_NodeTraversal_AbstractModuleCallback() {}
func (*ClosureCheckModule) Implements_jscomp_CompilerPass()                         {}
func (*ClosureCheckModule) Implements_jscomp_NodeTraversal_Callback()               {}

var static_ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE *DiagnosticType

func StaticGet_ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE
}
func StaticSet_ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE = v
}

var static_ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE *DiagnosticType

func StaticGet_ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE
}
func StaticSet_ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE = v
}

var static_ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE *DiagnosticType

func StaticGet_ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE
}
func StaticSet_ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE = v
}

var static_ClosureCheckModule_GOOG_MODULE_MISPLACED *DiagnosticType

func StaticGet_ClosureCheckModule_GOOG_MODULE_MISPLACED() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_GOOG_MODULE_MISPLACED
}
func StaticSet_ClosureCheckModule_GOOG_MODULE_MISPLACED(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_GOOG_MODULE_MISPLACED = v
}

var static_ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE *DiagnosticType

func StaticGet_ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE
}
func StaticSet_ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE = v
}

var static_ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS *DiagnosticType

func StaticGet_ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS
}
func StaticSet_ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS = v
}

var static_ClosureCheckModule_GOOG_MODULE_USES_THROW *DiagnosticType

func StaticGet_ClosureCheckModule_GOOG_MODULE_USES_THROW() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_GOOG_MODULE_USES_THROW
}
func StaticSet_ClosureCheckModule_GOOG_MODULE_USES_THROW(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_GOOG_MODULE_USES_THROW = v
}

var static_ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE *DiagnosticType

func StaticGet_ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE
}
func StaticSet_ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE = v
}

var static_ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE *DiagnosticType

func StaticGet_ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE
}
func StaticSet_ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE = v
}

var static_ClosureCheckModule_LET_GOOG_REQUIRE *DiagnosticType

func StaticGet_ClosureCheckModule_LET_GOOG_REQUIRE() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_LET_GOOG_REQUIRE
}
func StaticSet_ClosureCheckModule_LET_GOOG_REQUIRE(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_LET_GOOG_REQUIRE = v
}

var static_ClosureCheckModule_MULTIPLE_MODULES_IN_FILE *DiagnosticType

func StaticGet_ClosureCheckModule_MULTIPLE_MODULES_IN_FILE() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_MULTIPLE_MODULES_IN_FILE
}
func StaticSet_ClosureCheckModule_MULTIPLE_MODULES_IN_FILE(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_MULTIPLE_MODULES_IN_FILE = v
}

var static_ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION *DiagnosticType

func StaticGet_ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION
}
func StaticSet_ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION = v
}

var static_ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION *DiagnosticType

func StaticGet_ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION
}
func StaticSet_ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION = v
}

var static_ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE *DiagnosticType

func StaticGet_ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE
}
func StaticSet_ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE = v
}

var static_ClosureCheckModule_EXPORT_NOT_A_STATEMENT *DiagnosticType

func StaticGet_ClosureCheckModule_EXPORT_NOT_A_STATEMENT() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_EXPORT_NOT_A_STATEMENT
}
func StaticSet_ClosureCheckModule_EXPORT_NOT_A_STATEMENT(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_EXPORT_NOT_A_STATEMENT = v
}

var static_ClosureCheckModule_EXPORT_REPEATED_ERROR *DiagnosticType

func StaticGet_ClosureCheckModule_EXPORT_REPEATED_ERROR() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_EXPORT_REPEATED_ERROR
}
func StaticSet_ClosureCheckModule_EXPORT_REPEATED_ERROR(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_EXPORT_REPEATED_ERROR = v
}

var static_ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME *DiagnosticType

func StaticGet_ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME
}
func StaticSet_ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME = v
}

var static_ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME *DiagnosticType

func StaticGet_ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME
}
func StaticSet_ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME = v
}

var static_ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME *DiagnosticType

func StaticGet_ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME
}
func StaticSet_ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME = v
}

var static_ClosureCheckModule_USE_OF_GOOG_PROVIDE *DiagnosticType

func StaticGet_ClosureCheckModule_USE_OF_GOOG_PROVIDE() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_USE_OF_GOOG_PROVIDE
}
func StaticSet_ClosureCheckModule_USE_OF_GOOG_PROVIDE(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_USE_OF_GOOG_PROVIDE = v
}

var static_ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL *DiagnosticType

func StaticGet_ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL
}
func StaticSet_ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL = v
}

var static_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL *DiagnosticType

func StaticGet_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL
}
func StaticSet_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL = v
}

var static_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE *DiagnosticType

func StaticGet_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE
}
func StaticSet_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE = v
}

var static_ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT *DiagnosticType

func StaticGet_ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT() *DiagnosticType {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT
}
func StaticSet_ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT(v *DiagnosticType) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT = v
}

var static_ClosureCheckModule_GOOG_MODULE pkg_rhino.QualifiedName

func StaticGet_ClosureCheckModule_GOOG_MODULE() pkg_rhino.QualifiedName {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_GOOG_MODULE
}
func StaticSet_ClosureCheckModule_GOOG_MODULE(v pkg_rhino.QualifiedName) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_GOOG_MODULE = v
}

var static_ClosureCheckModule_GOOG_PROVIDE pkg_rhino.QualifiedName

func StaticGet_ClosureCheckModule_GOOG_PROVIDE() pkg_rhino.QualifiedName {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_GOOG_PROVIDE
}
func StaticSet_ClosureCheckModule_GOOG_PROVIDE(v pkg_rhino.QualifiedName) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_GOOG_PROVIDE = v
}

var static_ClosureCheckModule_GOOG_REQUIRE pkg_rhino.QualifiedName

func StaticGet_ClosureCheckModule_GOOG_REQUIRE() pkg_rhino.QualifiedName {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_GOOG_REQUIRE
}
func StaticSet_ClosureCheckModule_GOOG_REQUIRE(v pkg_rhino.QualifiedName) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_GOOG_REQUIRE = v
}

var static_ClosureCheckModule_GOOG_REQUIRE_TYPE pkg_rhino.QualifiedName

func StaticGet_ClosureCheckModule_GOOG_REQUIRE_TYPE() pkg_rhino.QualifiedName {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_GOOG_REQUIRE_TYPE
}
func StaticSet_ClosureCheckModule_GOOG_REQUIRE_TYPE(v pkg_rhino.QualifiedName) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_GOOG_REQUIRE_TYPE = v
}

var static_ClosureCheckModule_GOOG_MODULE_GET pkg_rhino.QualifiedName

func StaticGet_ClosureCheckModule_GOOG_MODULE_GET() pkg_rhino.QualifiedName {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_GOOG_MODULE_GET
}
func StaticSet_ClosureCheckModule_GOOG_MODULE_GET(v pkg_rhino.QualifiedName) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_GOOG_MODULE_GET = v
}

var static_ClosureCheckModule_GOOG_FORWARD_DECLARE pkg_rhino.QualifiedName

func StaticGet_ClosureCheckModule_GOOG_FORWARD_DECLARE() pkg_rhino.QualifiedName {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_GOOG_FORWARD_DECLARE
}
func StaticSet_ClosureCheckModule_GOOG_FORWARD_DECLARE(v pkg_rhino.QualifiedName) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_GOOG_FORWARD_DECLARE = v
}

var static_ClosureCheckModule_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE pkg_rhino.QualifiedName

func StaticGet_ClosureCheckModule_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE() pkg_rhino.QualifiedName {
	clinit_ClosureCheckModule()
	return static_ClosureCheckModule_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE
}
func StaticSet_ClosureCheckModule_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE(v pkg_rhino.QualifiedName) {
	clinit_ClosureCheckModule()
	static_ClosureCheckModule_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE = v
}

var clinitDone_ClosureCheckModule uint32
var clinitState_ClosureCheckModule j2gort.ClinitState

func clinit_ClosureCheckModule() {
	if atomic.LoadUint32(&clinitDone_ClosureCheckModule) != 0 {
		return
	}
	clinitSlow_ClosureCheckModule()
}

func clinitSlow_ClosureCheckModule() {
	if !j2gort.EnterClinit(&clinitState_ClosureCheckModule, &clinitDone_ClosureCheckModule) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ClosureCheckModule, &clinitDone_ClosureCheckModule)
	static_ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE = DiagnosticType_Error("JSC_AT_EXPORT_IN_GOOG_MODULE", "@export has no effect on top-level names in a goog.module. Consider using goog.exportSymbol instead.")
	static_ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE = DiagnosticType_Error("JSC_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE", "@export is not allowed here in a non-legacy goog.module. Consider using goog.exportSymbol instead.")
	static_ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE = DiagnosticType_Error("JSC_GOOG_MODULE_IN_NON_MODULE", "goog.module() call must be the first statement in a module.")
	static_ClosureCheckModule_GOOG_MODULE_MISPLACED = DiagnosticType_Error("JSC_GOOG_MODULE_MISPLACED", "goog.module() call must be the first statement in a file.")
	static_ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE = DiagnosticType_Error("JSC_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE", "goog.module.declareLegacyNamespace may only be called in a goog.module.")
	static_ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS = DiagnosticType_Error("JSC_GOOG_MODULE_REFERENCES_THIS", "The body of a goog.module cannot reference 'this'.")
	static_ClosureCheckModule_GOOG_MODULE_USES_THROW = DiagnosticType_Error("JSC_GOOG_MODULE_USES_THROW", "The body of a goog.module cannot use 'throw'.")
	static_ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE = DiagnosticType_Error("JSC_DUPLICATE_NAME_SHORT_REQUIRE", "Found multiple goog.require statements importing identifier ''{0}''.")
	static_ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE = DiagnosticType_Error("JSC_INVALID_DESTRUCTURING_REQUIRE", "Destructuring goog.require must be a simple object pattern.")
	static_ClosureCheckModule_LET_GOOG_REQUIRE = DiagnosticType_Disabled("JSC_LET_GOOG_REQUIRE", "Module imports must be constant. Please use ''const'' instead of ''let''.")
	static_ClosureCheckModule_MULTIPLE_MODULES_IN_FILE = DiagnosticType_Error("JSC_MULTIPLE_MODULES_IN_FILE", "There should only be a single goog.module() statement per file.")
	static_ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION = DiagnosticType_Error("JSC_ONE_REQUIRE_PER_DECLARATION", "There may only be one goog.require() per var/let/const declaration.")
	static_ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION = DiagnosticType_Disabled("JSC_INCORRECT_SHORTNAME_CAPITALIZATION", "The capitalization of short name {0} is incorrect; it should be {1}.")
	static_ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE = DiagnosticType_Error("JSC_EXPORT_NOT_AT_MODULE_SCOPE", "Exports must be at the top-level of a module")
	static_ClosureCheckModule_EXPORT_NOT_A_STATEMENT = DiagnosticType_Error("JSC_EXPORT_NOT_A_STATEMENT", "Exports should be a statement.")
	static_ClosureCheckModule_EXPORT_REPEATED_ERROR = DiagnosticType_Error("JSC_EXPORT_REPEATED_ERROR", "Name cannot be exported multiple times. Previous export on line {0}.")
	static_ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME = DiagnosticType_Error("JSC_REFERENCE_TO_MODULE_GLOBAL_NAME", "References to the global name of a module are not allowed. Perhaps you meant exports?")
	static_ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME = DiagnosticType_Disabled("JSC_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME", "Reference to fully qualified import name ''{0}''. Imports in goog.module should use the return value of goog.require / goog.forwardDeclare instead.")
	static_ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME = DiagnosticType_Disabled("JSC_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME", "Reference to fully qualified import name ''{0}''. Please use the short name ''{1}'' instead.")
	static_ClosureCheckModule_USE_OF_GOOG_PROVIDE = DiagnosticType_Disabled("JSC_USE_OF_GOOG_PROVIDE", "goog.provide is deprecated in favor of goog.module.")
	static_ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL = DiagnosticType_Error("JSC_REQUIRE_NOT_AT_TOP_LEVEL", "goog.require() must be called at file scope.")
	static_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL = DiagnosticType_Error("JSC_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL", "goog.module.declareLegacyNamespace() does not return a value")
	static_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE = DiagnosticType_Error("JSC_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL", "goog.module.declareLegacyNamespace() must be immediately after the goog.module('...'); call")
	static_ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT = DiagnosticType_Error("JSC_LEGACY_NAMESPACE_ARGUMENT", "goog.module.declareLegacyNamespace() takes no arguments")
	static_ClosureCheckModule_GOOG_MODULE = pkg_rhino.QualifiedName_Of("goog.module")
	static_ClosureCheckModule_GOOG_PROVIDE = pkg_rhino.QualifiedName_Of("goog.provide")
	static_ClosureCheckModule_GOOG_REQUIRE = pkg_rhino.QualifiedName_Of("goog.require")
	static_ClosureCheckModule_GOOG_REQUIRE_TYPE = pkg_rhino.QualifiedName_Of("goog.requireType")
	static_ClosureCheckModule_GOOG_MODULE_GET = pkg_rhino.QualifiedName_Of("goog.module.get")
	static_ClosureCheckModule_GOOG_FORWARD_DECLARE = pkg_rhino.QualifiedName_Of("goog.forwardDeclare")
	static_ClosureCheckModule_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE = pkg_rhino.QualifiedName_Of("goog.module.declareLegacyNamespace")
}

func (this_ *ClosureCheckModule) initFields_ClosureCheckModule() {
	_ = this_

}

// Java constructor: ClosureCheckModule(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.modules.ModuleMetadataMap moduleMetadataMap)
func NewClosureCheckModule(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *ClosureCheckModule {
	clinit_ClosureCheckModule()
	this_ := &ClosureCheckModule{}
	this_.BindSelf_(this_)
	InitClosureCheckModule(this_, compiler, moduleMetadataMap)
	return this_
}

func InitClosureCheckModule(this_ *ClosureCheckModule, compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) {
	_ = compiler
	_ = moduleMetadataMap
	InitNodeTraversal_AbstractModuleCallback(&this_.NodeTraversal_AbstractModuleCallbackBase, compiler, moduleMetadataMap)
	this_.initFields_ClosureCheckModule()
}

// Java method: ClosureCheckModule.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ClosureCheckModule) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.GetNodeTraversal_AbstractModuleCallbackBase().compiler, root, this_)
}

// Java method: ClosureCheckModule.enterModule(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentModule, com.google.javascript.rhino.Node moduleScopeRoot)
func (this_ *ClosureCheckModule) EnterModule(currentModule ModuleMetadataMap_ModuleMetadata, moduleScopeRoot pkg_rhino.Node) {
	_ = currentModule
	_ = moduleScopeRoot
	if !(currentModule.IsGoogModule()) {
		return
	}
	if !(this_.currentModuleInfo == nil) {
		j2gort.CheckState(false)
	}
	if !(!(currentModule.GoogNamespaces().IsEmpty())) {
		j2gort.CheckState(false)
	}
	this_.currentModuleInfo = NewClosureCheckModule_ModuleInfo(j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetFirst(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](currentModule.GoogNamespaces()), "")))
}

// Java method: ClosureCheckModule.exitModule(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentModule, com.google.javascript.rhino.Node moduleScopeRoot)
func (this_ *ClosureCheckModule) ExitModule(currentModule ModuleMetadataMap_ModuleMetadata, moduleScopeRoot pkg_rhino.Node) {
	_ = currentModule
	_ = moduleScopeRoot
	this_.currentModuleInfo = nil
}

// Java method: ClosureCheckModule.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.ModuleMetadataMap.@org.jspecify.annotations.Nullable ModuleMetadata currentModule, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node moduleScopeRoot)
func (this_ *ClosureCheckModule) Visit_4(t *NodeTraversal, n pkg_rhino.Node, currentModule ModuleMetadataMap_ModuleMetadata, moduleScopeRoot pkg_rhino.Node) {
	_ = t
	_ = n
	_ = currentModule
	_ = moduleScopeRoot
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if this_.currentModuleInfo == nil {
		if NodeUtil_IsCallTo_2_Node_QualifiedName(n, StaticGet_ClosureCheckModule_GOOG_MODULE()) {
			t.Report_3(n, StaticGet_ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE(), []string{})
		} else if NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(n) {
			t.Report_3(n, StaticGet_ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE(), []string{})
		} else if NodeUtil_IsCallTo_2_Node_QualifiedName(n, StaticGet_ClosureCheckModule_GOOG_PROVIDE()) {
			t.Report_3(n, StaticGet_ClosureCheckModule_USE_OF_GOOG_PROVIDE(), []string{})
		}
		return
	}
	var jsDoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = jsDoc
	if jsDoc != nil {
		this_.priv_ClosureCheckModule_checkJSDoc_1140(t, jsDoc)
	}
	{
		var callee pkg_rhino.Node = pkg_rhino.Node(0)
		_ = callee
		switch n.TokenSwitchOrd_() {
		case int32(28) /* Token.CALL */ :
			callee = n.GetFirstChild()
			if StaticGet_ClosureCheckModule_GOOG_MODULE().Matches(callee) {
				{
					var iifeCond_1 bool
					{

						sr_1 := this_.currentModuleInfo.name
						se_2 := ClosureCheckModule_ExtractFirstArgumentName(n)
						iifeCond_1 = j2gort.CheckNotNullString(sr_1) == se_2
					}
					if !(iifeCond_1) {
						t.Report_3(n, StaticGet_ClosureCheckModule_MULTIPLE_MODULES_IN_FILE(), []string{})
					} else if !(ClosureCheckModule_IsFirstExpressionInGoogModule(parent)) {
						t.Report_3(n, StaticGet_ClosureCheckModule_GOOG_MODULE_MISPLACED(), []string{})
					}
				}
			} else if (StaticGet_ClosureCheckModule_GOOG_REQUIRE().Matches(callee) || StaticGet_ClosureCheckModule_GOOG_REQUIRE_TYPE().Matches(callee)) || StaticGet_ClosureCheckModule_GOOG_FORWARD_DECLARE().Matches(callee) {
				this_.priv_ClosureCheckModule_checkRequireCall_1144(t, n, parent)
			} else if StaticGet_ClosureCheckModule_GOOG_MODULE_GET().Matches(callee) && t.InModuleHoistScope() {
				t.Report_3(n, StaticGet_ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET(), []string{})
			} else if StaticGet_ClosureCheckModule_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE().Matches(callee) {
				ClosureCheckModule_CheckLegacyNamespaceCall(t, n, parent)
			}
		case int32(51) /* Token.ASSIGN */ :
			{
				if ClosureCheckModule_IsExportLhs(n.GetFirstChild()) {
					this_.priv_ClosureCheckModule_checkModuleExport_1143(t, n, parent)
				}
				break
			}
		case int32(105) /* Token.CLASS */, int32(73) /* Token.FUNCTION */ :
			if !(NodeUtil_IsStatement(n)) {
				break
			}
			fallthrough
		case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
			if (t.InModuleHoistScope() && (n.IsClass() || (NodeUtil_GetEnclosingClass(n) == 0))) && (NodeUtil_GetEnclosingType(n, pkg_rhino.StaticGet_Token_OBJECTLIT()) == 0) {
				var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
				_ = jsdoc
				if (jsdoc != nil) && jsdoc.IsExport() {
					t.Report_3(n, StaticGet_ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE(), []string{})
				}
			}
		case int32(37) /* Token.THIS */ :
			if t.InModuleHoistScope() {
				t.Report_3(n, StaticGet_ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS(), []string{})
			}
		case int32(43) /* Token.THROW */ :
			if t.InModuleHoistScope() {
				t.Report_3(n, StaticGet_ClosureCheckModule_GOOG_MODULE_USES_THROW(), []string{})
			}
		case int32(32) /* Token.NAME */, int32(26) /* Token.GETPROP */ :
			if n.MatchesQualifiedName_1_string(this_.currentModuleInfo.name) {
				if n.IsGetProp() {
					t.Report_3(n, StaticGet_ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME(), []string{})
				}
			} else {
				if !(parent.IsGetProp()) && n.IsQualifiedName() {
					var node pkg_rhino.Node = n
					_ = node
					for this_.priv_ClosureCheckModule_checkImproperReferenceToImport_1142(t, node, node.GetQualifiedName(), func() string {
						if node.GetParent().IsGetProp() {
							return node.GetParent().GetString()
						}
						return j2gort.NullString
					}()) {
						if node.IsName() || !(node.HasChildren()) {
							return
						}
						node = node.GetOnlyChild()
						if node == 0 {
							return
						}
					}
				}
			}
		default:
		}
	}
}

// Java method: ClosureCheckModule.checkJSDoc(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.JSDocInfo jsDoc)
func (this_ *ClosureCheckModule) priv_ClosureCheckModule_checkJSDoc_1140(t *NodeTraversal, jsDoc *pkg_rhino.JSDocInfo) {
	_ = t
	_ = jsDoc
	{
		var i int32 = int32(0)
		_ = i
		var count int32 = jsDoc.GetTypeNodeCount()
		_ = count
		for ; i < count; i = i + 1 {
			this_.priv_ClosureCheckModule_checkTypeExpression_1141(t, jsDoc.GetTypeNodeAt(i))
		}
	}
}

// Java method: ClosureCheckModule.checkTypeExpression(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node)
func (this_ *ClosureCheckModule) priv_ClosureCheckModule_checkTypeExpression_1141(t *NodeTraversal, node pkg_rhino.Node) {
	_ = t
	_ = node
	if node.IsStringLit() {
		var qname string = node.GetString()
		_ = qname
		var nextQnamePart string = j2gort.NullString
		_ = nextQnamePart
		for this_.priv_ClosureCheckModule_checkImproperReferenceToImport_1142(t, node, qname, nextQnamePart) {
			var lastDot int32 = j2gort.JStrLastIndexOfCp(qname, int32(uint16(46)))
			_ = lastDot
			if lastDot < int32(0) {
				break
			}
			nextQnamePart = j2gort.JStrSubFrom(qname, (lastDot + int32(1)))
			qname = j2gort.JStrSub(qname, int32(0), lastDot)
		}
	}
	{
		var child pkg_rhino.Node = node.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_ClosureCheckModule_checkTypeExpression_1141(t, child)
		}
	}
}

// Java static method: ClosureCheckModule.isLocalVar(com.google.javascript.jscomp.Var v)
func ClosureCheckModule_IsLocalVar(v *Var) bool {
	clinit_ClosureCheckModule()
	_ = v
	if v == nil {
		return false
	}
	return v.IsLocal()
}

// Java method: ClosureCheckModule.checkImproperReferenceToImport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, java.lang.String qname, java.lang.@org.jspecify.annotations.Nullable String nextQnamePart)
func (this_ *ClosureCheckModule) priv_ClosureCheckModule_checkImproperReferenceToImport_1142(t *NodeTraversal, n pkg_rhino.Node, qname string, nextQnamePart string) bool {
	_ = t
	_ = n
	_ = qname
	_ = nextQnamePart
	if qname == j2gort.NullString {
		return true
	}
	var importLhs pkg_rhino.Node = this_.currentModuleInfo.importsByLongRequiredName.Get(qname)
	_ = importLhs
	if importLhs == 0 {
		return true
	}
	if !(j2gort.JStrIndexOf(qname, ".") >= 0) && ClosureCheckModule_IsLocalVar(j2gort.CastPtr[Var](t.GetScope().GetVar(qname))) {
		return true
	}
	if importLhs.IsName() {
		this_.GetNodeTraversal_AbstractModuleCallbackBase().compiler.Report_1(JSError_Make_3(n, StaticGet_ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME(), []string{qname, importLhs.GetString()}))
		return false
	} else if importLhs.IsDestructuringLhs() && (nextQnamePart != j2gort.NullString) {
		var objPattern pkg_rhino.Node = importLhs.GetFirstChild()
		_ = objPattern
		if objPattern.IsObjectPattern() {
			{
				var strKey pkg_rhino.Node = objPattern.GetFirstChild()
				_ = strKey
				for ; strKey != 0; strKey = strKey.GetNext() {
					if ((strKey.IsStringKey() && strKey.HasOneChild()) && strKey.GetFirstChild().IsName()) && (j2gort.CheckNotNullString(strKey.GetString()) == nextQnamePart) {
						var parent pkg_rhino.Node = n.GetParent()
						_ = parent
						this_.GetNodeTraversal_AbstractModuleCallbackBase().compiler.Report_1(JSError_Make_3(func() pkg_rhino.Node {
							if (parent != 0) && parent.IsGetProp() {
								return parent
							}
							return n
						}(), StaticGet_ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME(), []string{(j2gort.StrValString(qname) + "." + j2gort.StrValString(nextQnamePart)), strKey.GetFirstChild().GetString()}))
						return false
					}
				}
			}
		}
	}
	if importLhs.IsExprResult() && NodeUtil_IsGoogForwardDeclareCall(importLhs.GetOnlyChild()) {
		return true
	}
	this_.GetNodeTraversal_AbstractModuleCallbackBase().compiler.Report_1(JSError_Make_3(n, StaticGet_ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME(), []string{qname}))
	return false
}

// Java static method: ClosureCheckModule.isExportLhs(com.google.javascript.rhino.Node lhs)
func ClosureCheckModule_IsExportLhs(lhs pkg_rhino.Node) bool {
	clinit_ClosureCheckModule()
	_ = lhs
	if !(lhs.IsQualifiedName()) {
		return false
	}
	return (lhs.MatchesName_1_string("exports") || (lhs.IsGetProp() && lhs.GetFirstChild().MatchesName_1_string("exports")))
}

// Java method: ClosureCheckModule.checkModuleExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ClosureCheckModule) priv_ClosureCheckModule_checkModuleExport_1143(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsAssign()) {
		j2gort.CheckArgument(false)
	}
	var lhs pkg_rhino.Node = n.GetFirstChild()
	_ = lhs
	if !(ClosureCheckModule_IsExportLhs(lhs)) {
		j2gort.CheckState(false)
	}
	var previousDefinition pkg_rhino.Node = this_.currentModuleInfo.exportNodesByName.Get(lhs.GetQualifiedName())
	_ = previousDefinition
	if previousDefinition != 0 {
		var previousLine int32 = previousDefinition.GetLineno()
		_ = previousLine
		t.Report_3(n, StaticGet_ClosureCheckModule_EXPORT_REPEATED_ERROR(), []string{strconv.FormatInt(int64(previousLine), 10)})
	}
	var defaultExportNode pkg_rhino.Node = this_.currentModuleInfo.exportNodesByName.Get("exports")
	_ = defaultExportNode
	if (defaultExportNode == 0) || lhs.MatchesName_1_string("exports") {
		this_.currentModuleInfo.exportNodesByName.Put(lhs.GetQualifiedName(), lhs)
		if !(t.InModuleScope()) {
			t.Report_3(n, StaticGet_ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE(), []string{})
		} else if !(parent.IsExprResult()) {
			t.Report_3(n, StaticGet_ClosureCheckModule_EXPORT_NOT_A_STATEMENT(), []string{})
		}
	}
	if !(NodeUtil_IsPrototypeProperty(lhs)) && !(NodeUtil_IsLegacyGoogModuleFile(NodeUtil_GetEnclosingScript(n))) {
		var jsDoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
		_ = jsDoc
		if (jsDoc != nil) && jsDoc.IsExport() {
			t.Report_3(n, StaticGet_ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE(), []string{})
		}
	}
}

// Java static method: ClosureCheckModule.extractFirstArgumentName(com.google.javascript.rhino.Node callNode)
func ClosureCheckModule_ExtractFirstArgumentName(callNode pkg_rhino.Node) string {
	clinit_ClosureCheckModule()
	_ = callNode
	var firstArg pkg_rhino.Node = callNode.GetSecondChild()
	_ = firstArg
	if (firstArg != 0) && firstArg.IsStringLit() {
		return firstArg.GetString()
	}
	return j2gort.NullString
}

// Java method: ClosureCheckModule.checkRequireCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node callNode, com.google.javascript.rhino.Node parent)
func (this_ *ClosureCheckModule) priv_ClosureCheckModule_checkRequireCall_1144(t *NodeTraversal, callNode pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = callNode
	_ = parent
	if !(callNode.IsCall()) {
		j2gort.CheckState(false)
	}
	if !(callNode.GetLastChild().IsStringLit()) {
		j2gort.CheckState(false)
	}
	switch parent.TokenSwitchOrd_() {
	case int32(93) /* Token.EXPR_RESULT */ :
		var key string = ClosureCheckModule_ExtractFirstArgumentName(callNode)
		_ = key
		this_.currentModuleInfo.importsByLongRequiredName.PutIfAbsent(key, parent)
	case int32(32) /* Token.NAME */, int32(104) /* Token.DESTRUCTURING_LHS */ :
		this_.priv_ClosureCheckModule_checkShortGoogRequireCall_1145(t, callNode, parent.GetParent())
	case int32(115) /* Token.AWAIT */ :
		var grandParentToken pkg_rhino.Token = parent.GetParent().GetToken()
		_ = grandParentToken
		{
			var iifeCond_1 bool
			{

				re_3 := grandParentToken
				ae_4 := pkg_rhino.StaticGet_Token_DESTRUCTURING_LHS()
				if re_3 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_3 == ae_4
			}
			if iifeCond_1 || func() bool {
				re_5 := grandParentToken
				ae_6 := pkg_rhino.StaticGet_Token_NAME()
				if re_5 == 0 {
					j2gort.ThrowNPE()
				}
				return re_5 == ae_6
			}() {
				this_.priv_ClosureCheckModule_checkShortGoogRequireCall_1145(t, callNode, parent.GetGrandparent())
			}
		}
	default:
		t.Report_3(callNode, StaticGet_ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL(), []string{})
	}
}

// Java method: ClosureCheckModule.checkShortGoogRequireCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node callNode, com.google.javascript.rhino.Node declaration)
func (this_ *ClosureCheckModule) priv_ClosureCheckModule_checkShortGoogRequireCall_1145(t *NodeTraversal, callNode pkg_rhino.Node, declaration pkg_rhino.Node) {
	_ = t
	_ = callNode
	_ = declaration
	if declaration.IsLet() && !(StaticGet_ClosureCheckModule_GOOG_FORWARD_DECLARE().Matches(callNode.GetFirstChild())) {
		t.Report_3(declaration, StaticGet_ClosureCheckModule_LET_GOOG_REQUIRE(), []string{})
	}
	if !(declaration.HasOneChild()) {
		t.Report_3(declaration, StaticGet_ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION(), []string{})
		return
	}
	var lhs pkg_rhino.Node = declaration.GetFirstChild()
	_ = lhs
	if lhs.IsDestructuringLhs() {
		if !(ClosureCheckModule_IsValidDestructuringImport(lhs)) {
			t.Report_3(declaration, StaticGet_ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE(), []string{})
		}
		if StaticGet_ClosureCheckModule_GOOG_FORWARD_DECLARE().Matches(callNode.GetFirstChild()) {
			t.Report_3(lhs, StaticGet_ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE(), []string{})
		}
	} else if lhs.IsName() {
		ClosureCheckModule_CheckShortName(t, lhs, callNode.GetLastChild().GetString())
	}
	this_.currentModuleInfo.importsByLongRequiredName.Put(ClosureCheckModule_ExtractFirstArgumentName(callNode), lhs)
	NodeUtil_VisitLhsNodesInNode(declaration, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_7 any) {
		nameNode := j2gort.CastNode[pkg_rhino.Node](lp_7)
		_ = nameNode
		var name string = nameNode.GetString()
		_ = name
		if !(this_.currentModuleInfo.shortImportNames.Add(name)) {
			t.Report_3(nameNode, StaticGet_ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE(), []string{name})
		}
	}))
}

// Java static method: ClosureCheckModule.checkShortName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node shortNameNode, java.lang.String namespace)
func ClosureCheckModule_CheckShortName(t *NodeTraversal, shortNameNode pkg_rhino.Node, namespace string) {
	clinit_ClosureCheckModule()
	_ = t
	_ = shortNameNode
	_ = namespace
	var nextQnamePart string = shortNameNode.GetString()
	_ = nextQnamePart
	if javalang.StringStartsWith(j2gort.CheckNotNullString(namespace), "google3.") {
		return
	}
	var lastSegment string = j2gort.JStrSubFrom(namespace, (j2gort.JStrLastIndexOfCp(namespace, int32(uint16(46))) + int32(1)))
	_ = lastSegment
	if (j2gort.CheckNotNullString(nextQnamePart) == lastSegment) || (j2gort.CheckNotNullString(lastSegment) == "") {
		return
	}
	if j2gort.Ext_com_google_common_base_Ascii_IsUpperCase(j2gort.JStrCharAt(nextQnamePart, int32(0))) != j2gort.Ext_com_google_common_base_Ascii_IsUpperCase(j2gort.JStrCharAt(lastSegment, int32(0))) {
		var newStartChar uint16
		_ = newStartChar
		if j2gort.Ext_com_google_common_base_Ascii_IsUpperCase(j2gort.JStrCharAt(nextQnamePart, int32(0))) {
			newStartChar = j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_uint16(j2gort.JStrCharAt(nextQnamePart, int32(0)))
		} else {
			newStartChar = j2gort.Ext_com_google_common_base_Ascii_ToUpperCase_1_uint16(j2gort.JStrCharAt(nextQnamePart, int32(0)))
		}
		var correctedName string = (j2gort.StrValChar(newStartChar) + j2gort.StrValString(j2gort.JStrSubFrom(nextQnamePart, int32(1))))
		_ = correctedName
		t.Report_3(shortNameNode, StaticGet_ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION(), []string{nextQnamePart, correctedName})
	}
}

// Java static method: ClosureCheckModule.isValidDestructuringImport(com.google.javascript.rhino.Node destructuringLhs)
func ClosureCheckModule_IsValidDestructuringImport(destructuringLhs pkg_rhino.Node) bool {
	clinit_ClosureCheckModule()
	_ = destructuringLhs
	if !(destructuringLhs.IsDestructuringLhs()) {
		j2gort.CheckArgument(false)
	}
	var objectPattern pkg_rhino.Node = destructuringLhs.GetFirstChild()
	_ = objectPattern
	if !(objectPattern.IsObjectPattern()) {
		return false
	}
	{
		var stringKey pkg_rhino.Node = objectPattern.GetFirstChild()
		_ = stringKey
		for ; stringKey != 0; stringKey = stringKey.GetNext() {
			if !(stringKey.IsStringKey()) {
				return false
			} else if !(stringKey.HasChildren()) || !(stringKey.GetFirstChild().IsName()) {
				return false
			}
		}
	}
	return true
}

// Java static method: ClosureCheckModule.checkLegacyNamespaceCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node callNode, com.google.javascript.rhino.Node parent)
func ClosureCheckModule_CheckLegacyNamespaceCall(t *NodeTraversal, callNode pkg_rhino.Node, parent pkg_rhino.Node) {
	clinit_ClosureCheckModule()
	_ = t
	_ = callNode
	_ = parent
	if !(callNode.IsCall()) {
		j2gort.CheckArgument(false)
	}
	if callNode.GetChildCount() > int32(1) {
		t.Report_3(callNode, StaticGet_ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT(), []string{})
	}
	if !(parent.IsExprResult()) {
		t.Report_3(callNode, StaticGet_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL(), []string{})
		return
	}
	var prev pkg_rhino.Node = parent.GetPrevious()
	_ = prev
	if (prev == 0) || !(NodeUtil_IsGoogModuleCall(prev)) {
		t.Report_3(callNode, StaticGet_ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE(), []string{})
	}
}

// Java static method: ClosureCheckModule.isFirstExpressionInGoogModule(com.google.javascript.rhino.Node node)
func ClosureCheckModule_IsFirstExpressionInGoogModule(node pkg_rhino.Node) bool {
	clinit_ClosureCheckModule()
	_ = node
	if !(node.IsExprResult()) || (node.GetPrevious() != 0) {
		return false
	}
	var statementBlock pkg_rhino.Node = node.GetParent()
	_ = statementBlock
	return (statementBlock.IsModuleBody() || NodeUtil_IsBundledGoogModuleScopeRoot(statementBlock))
}

// Java class: com.google.javascript.jscomp.ClosureCheckModule.ModuleInfo
type ClosureCheckModule_ModuleInfo struct {
	name                      string                                    `gson:"name"`
	importsByLongRequiredName j2gort.OrderedMap[string, pkg_rhino.Node] `gson:"importsByLongRequiredName"`
	shortImportNames          j2gort.OrderedSet[string]                 `gson:"shortImportNames"`
	exportNodesByName         j2gort.OrderedMap[string, pkg_rhino.Node] `gson:"exportNodesByName"`
	idHash_                   atomic.Int32
}

func (this_ *ClosureCheckModule_ModuleInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureCheckModule_ModuleInfo) IsJavaNil_() bool { return this_ == nil }
func (*ClosureCheckModule_ModuleInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureCheckModule$ModuleInfo", "ModuleInfo"
}
func (this_ *ClosureCheckModule_ModuleInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureCheckModule_ModuleInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureCheckModule.ModuleInfo", "com.google.javascript.jscomp.ClosureCheckModule$ModuleInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureCheckModule_ModuleInfo](v_)
	})
}

func (this_ *ClosureCheckModule_ModuleInfo) initFields_ClosureCheckModule_ModuleInfo() {
	_ = this_
	this_.name = j2gort.NullString
	this_.importsByLongRequiredName = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	this_.shortImportNames = j2gort.NewOrderedSet[string]()
	this_.exportNodesByName = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
}

// Java constructor: ModuleInfo(java.lang.String moduleName)
func NewClosureCheckModule_ModuleInfo(moduleName string) *ClosureCheckModule_ModuleInfo {
	this_ := &ClosureCheckModule_ModuleInfo{}
	InitClosureCheckModule_ModuleInfo(this_, moduleName)
	return this_
}

func InitClosureCheckModule_ModuleInfo(this_ *ClosureCheckModule_ModuleInfo, moduleName string) {
	_ = moduleName
	this_.initFields_ClosureCheckModule_ModuleInfo()
	this_.name = moduleName
}
func (*ClosureCheckModule) CanReleaseSyntacticScopes_() {}
