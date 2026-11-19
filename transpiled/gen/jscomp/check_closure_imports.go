/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CheckClosureImports.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CheckClosureImports
type CheckClosureImports struct {
	compiler       AbstractCompiler             `gson:"compiler"`
	checker        *CheckClosureImports_Checker `gson:"checker"`
	namespacesSeen j2gort.OrderedSet[string]    `gson:"namespacesSeen"`
	chunkGraph     *JSChunkGraph                `gson:"chunkGraph"`
	idHash_        atomic.Int32
}

func (this_ *CheckClosureImports) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckClosureImports) IsJavaNil_() bool { return this_ == nil }
func (*CheckClosureImports) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckClosureImports", "CheckClosureImports"
}
func (this_ *CheckClosureImports) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckClosureImports) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckClosureImports", "com.google.javascript.jscomp.CheckClosureImports", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckClosureImports](v_)
	})
}
func (*CheckClosureImports) Implements_jscomp_CompilerPass() {}

var static_CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING *DiagnosticType

func StaticGet_CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING() *DiagnosticType {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING
}
func StaticSet_CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING(v *DiagnosticType) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING = v
}

var static_CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION *DiagnosticType

func StaticGet_CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION() *DiagnosticType {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION
}
func StaticSet_CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION(v *DiagnosticType) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION = v
}

var static_CheckClosureImports_INVALID_CLOSURE_IMPORT_CALL *DiagnosticType

func StaticGet_CheckClosureImports_INVALID_CLOSURE_IMPORT_CALL() *DiagnosticType {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_INVALID_CLOSURE_IMPORT_CALL
}
func StaticSet_CheckClosureImports_INVALID_CLOSURE_IMPORT_CALL(v *DiagnosticType) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_INVALID_CLOSURE_IMPORT_CALL = v
}

var static_CheckClosureImports_LATE_PROVIDE_ERROR *DiagnosticType

func StaticGet_CheckClosureImports_LATE_PROVIDE_ERROR() *DiagnosticType {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_LATE_PROVIDE_ERROR
}
func StaticSet_CheckClosureImports_LATE_PROVIDE_ERROR(v *DiagnosticType) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_LATE_PROVIDE_ERROR = v
}

var static_CheckClosureImports_LET_CLOSURE_IMPORT *DiagnosticType

func StaticGet_CheckClosureImports_LET_CLOSURE_IMPORT() *DiagnosticType {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_LET_CLOSURE_IMPORT
}
func StaticSet_CheckClosureImports_LET_CLOSURE_IMPORT(v *DiagnosticType) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_LET_CLOSURE_IMPORT = v
}

var static_CheckClosureImports_NO_CLOSURE_IMPORT_DESTRUCTURING *DiagnosticType

func StaticGet_CheckClosureImports_NO_CLOSURE_IMPORT_DESTRUCTURING() *DiagnosticType {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_NO_CLOSURE_IMPORT_DESTRUCTURING
}
func StaticSet_CheckClosureImports_NO_CLOSURE_IMPORT_DESTRUCTURING(v *DiagnosticType) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_NO_CLOSURE_IMPORT_DESTRUCTURING = v
}

var static_CheckClosureImports_LHS_OF_CLOSURE_IMPORT_MUST_BE_CONST_IN_ES_MODULE *DiagnosticType

func StaticGet_CheckClosureImports_LHS_OF_CLOSURE_IMPORT_MUST_BE_CONST_IN_ES_MODULE() *DiagnosticType {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_LHS_OF_CLOSURE_IMPORT_MUST_BE_CONST_IN_ES_MODULE
}
func StaticSet_CheckClosureImports_LHS_OF_CLOSURE_IMPORT_MUST_BE_CONST_IN_ES_MODULE(v *DiagnosticType) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_LHS_OF_CLOSURE_IMPORT_MUST_BE_CONST_IN_ES_MODULE = v
}

var static_CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR *DiagnosticType

func StaticGet_CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR() *DiagnosticType {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR
}
func StaticSet_CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR(v *DiagnosticType) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR = v
}

var static_CheckClosureImports_GOOG_REQUIRE pkg_rhino.Node

func StaticGet_CheckClosureImports_GOOG_REQUIRE() pkg_rhino.Node {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_GOOG_REQUIRE
}
func StaticSet_CheckClosureImports_GOOG_REQUIRE(v pkg_rhino.Node) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_GOOG_REQUIRE = v
}

var static_CheckClosureImports_GOOG_MODULE_GET pkg_rhino.Node

func StaticGet_CheckClosureImports_GOOG_MODULE_GET() pkg_rhino.Node {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_GOOG_MODULE_GET
}
func StaticSet_CheckClosureImports_GOOG_MODULE_GET(v pkg_rhino.Node) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_GOOG_MODULE_GET = v
}

var static_CheckClosureImports_GOOG_FORWARD_DECLARE pkg_rhino.Node

func StaticGet_CheckClosureImports_GOOG_FORWARD_DECLARE() pkg_rhino.Node {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_GOOG_FORWARD_DECLARE
}
func StaticSet_CheckClosureImports_GOOG_FORWARD_DECLARE(v pkg_rhino.Node) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_GOOG_FORWARD_DECLARE = v
}

var static_CheckClosureImports_GOOG_REQUIRE_TYPE pkg_rhino.Node

func StaticGet_CheckClosureImports_GOOG_REQUIRE_TYPE() pkg_rhino.Node {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_GOOG_REQUIRE_TYPE
}
func StaticSet_CheckClosureImports_GOOG_REQUIRE_TYPE(v pkg_rhino.Node) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_GOOG_REQUIRE_TYPE = v
}

var static_CheckClosureImports_GOOG_REQUIRE_DYNAMIC pkg_rhino.Node

func StaticGet_CheckClosureImports_GOOG_REQUIRE_DYNAMIC() pkg_rhino.Node {
	clinit_CheckClosureImports()
	return static_CheckClosureImports_GOOG_REQUIRE_DYNAMIC
}
func StaticSet_CheckClosureImports_GOOG_REQUIRE_DYNAMIC(v pkg_rhino.Node) {
	clinit_CheckClosureImports()
	static_CheckClosureImports_GOOG_REQUIRE_DYNAMIC = v
}

var clinitDone_CheckClosureImports uint32
var clinitState_CheckClosureImports j2gort.ClinitState

func clinit_CheckClosureImports() {
	if atomic.LoadUint32(&clinitDone_CheckClosureImports) != 0 {
		return
	}
	clinitSlow_CheckClosureImports()
}

func clinitSlow_CheckClosureImports() {
	if !j2gort.EnterClinit(&clinitState_CheckClosureImports, &clinitDone_CheckClosureImports) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckClosureImports, &clinitDone_CheckClosureImports)
	static_CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING = DiagnosticType_Error("JSC_INVALID_CLOSURE_IMPORT_DESTRUCTURING", "Destructuring {0} must be a simple object pattern.")
	static_CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION = DiagnosticType_Error("JSC_ONE_CLOSURE_IMPORT_PER_DECLARATION", "There may only be one {0} per var/let/const declaration.")
	static_CheckClosureImports_INVALID_CLOSURE_IMPORT_CALL = DiagnosticType_Error("JSC_INVALID_CLOSURE_IMPORT_CALL", "{0} parameter must be a string literal.")
	static_CheckClosureImports_LATE_PROVIDE_ERROR = DiagnosticType_Error("JSC_LATE_PROVIDE_ERROR", "Required namespace \"{0}\" not provided yet.")
	static_CheckClosureImports_LET_CLOSURE_IMPORT = DiagnosticType_Disabled("JSC_LET_CLOSURE_IMPORT", "Module imports must be constant. Please use ''const'' instead of ''let''.")
	static_CheckClosureImports_NO_CLOSURE_IMPORT_DESTRUCTURING = DiagnosticType_Error("JSC_NO_CLOSURE_IMPORT_DESTRUCTURING", "Cannot destructure the return value of {0}")
	static_CheckClosureImports_LHS_OF_CLOSURE_IMPORT_MUST_BE_CONST_IN_ES_MODULE = DiagnosticType_Error("JSC_LHS_OF_CLOUSRE_IMPORT_MUST_BE_CONST_IN_ES_MODULE", "The left side of a {0} must use ''const'' (not ''let'' or ''var'') in an ES module.")
	static_CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR = DiagnosticType_Warning("JSC_XMODULE_REQUIRE_ERROR", "namespace \"{0}\" is required in chunk {2} but provided in chunk {1}. Is chunk {2} missing a dependency on chunk {1}?")
	static_CheckClosureImports_GOOG_REQUIRE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "require")
	static_CheckClosureImports_GOOG_MODULE_GET = pkg_rhino.IR_Getprop_3(pkg_rhino.IR_Name("goog"), "module", []string{"get"})
	static_CheckClosureImports_GOOG_FORWARD_DECLARE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "forwardDeclare")
	static_CheckClosureImports_GOOG_REQUIRE_TYPE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "requireType")
	static_CheckClosureImports_GOOG_REQUIRE_DYNAMIC = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "requireDynamic")
}

func (this_ *CheckClosureImports) initFields_CheckClosureImports() {
	_ = this_
}

// Java constructor: CheckClosureImports(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.modules.ModuleMetadataMap moduleMetadataMap)
func NewCheckClosureImports(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *CheckClosureImports {
	clinit_CheckClosureImports()
	this_ := &CheckClosureImports{}
	InitCheckClosureImports(this_, compiler, moduleMetadataMap)
	return this_
}

func InitCheckClosureImports(this_ *CheckClosureImports, compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) {
	_ = compiler
	_ = moduleMetadataMap
	this_.initFields_CheckClosureImports()
	this_.compiler = compiler
	this_.checker = NewCheckClosureImports_Checker(this_, compiler, moduleMetadataMap)
	this_.namespacesSeen = j2gort.NewOrderedSet[string]()
	this_.chunkGraph = compiler.GetChunkGraph()
}

// Java method: CheckClosureImports.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckClosureImports) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, externs, ToIface_NodeTraversal_Callback(this_.checker))
	NodeTraversal_Traverse(this_.compiler, root, ToIface_NodeTraversal_Callback(this_.checker))
}

// Java method: CheckClosureImports.verifyRequireOrder(java.lang.String namespace, com.google.javascript.rhino.Node call, com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata requiredModule)
func (this_ *CheckClosureImports) priv_CheckClosureImports_verifyRequireOrder_1033(namespace string, call pkg_rhino.Node, t *NodeTraversal, requiredModule ModuleMetadataMap_ModuleMetadata) {
	_ = namespace
	_ = call
	_ = t
	_ = requiredModule
	if !(this_.namespacesSeen.Contains(namespace)) {
		t.Report_3(call, StaticGet_CheckClosureImports_LATE_PROVIDE_ERROR(), []string{namespace})
		return
	}
	if (requiredModule.RootNode() == 0) || requiredModule.RootNode().IsFromExterns() {
		return
	}
	var requiredInputId *pkg_rhino.InputId = NodeUtil_GetInputId(requiredModule.RootNode())
	_ = requiredInputId
	var requiredInput *CompilerInput = j2gort.CheckNotNullPtr(this_.compiler.GetInput(requiredInputId))
	_ = requiredInput
	var requiredChunk *JSChunk = requiredInput.GetChunk()
	_ = requiredChunk
	var currentChunk *JSChunk = t.GetChunk()
	_ = currentChunk
	if (currentChunk != requiredChunk) && !(this_.chunkGraph.DependsOn(currentChunk, requiredChunk)) {
		this_.compiler.Report_1(JSError_Make_3(call, StaticGet_CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR(), []string{namespace, requiredChunk.GetName(), currentChunk.GetName()}))
	}
}

// Java enum: com.google.javascript.jscomp.CheckClosureImports.ClosureImport
type CheckClosureImports_ClosureImport struct {
	j2gort.EnumBase
	fn_AllowDestructuring  func(*CheckClosureImports_ClosureImport) bool
	fn_AliasMustBeConstant func(*CheckClosureImports_ClosureImport) bool
	fn_MustBeOrdered       func(*CheckClosureImports_ClosureImport) bool
	fn_CallName            func(*CheckClosureImports_ClosureImport) string
}

func (this_ *CheckClosureImports_ClosureImport) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *CheckClosureImports_ClosureImport) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*CheckClosureImports_ClosureImport) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*CheckClosureImports_ClosureImport) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *CheckClosureImports_ClosureImport) initFields_CheckClosureImports_ClosureImport() {
	_ = this_
}

func initEnumCheckClosureImports_ClosureImport(this_ *CheckClosureImports_ClosureImport, enumName_ string, enumOrd_ int32) *CheckClosureImports_ClosureImport {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_CheckClosureImports_ClosureImport()
	return this_
}

var enum_CheckClosureImports_ClosureImport_REQUIRE *CheckClosureImports_ClosureImport

func StaticGet_CheckClosureImports_ClosureImport_REQUIRE() *CheckClosureImports_ClosureImport {
	clinit_CheckClosureImports_ClosureImport()
	return enum_CheckClosureImports_ClosureImport_REQUIRE
}

var enum_CheckClosureImports_ClosureImport_FORWARD_DECLARE *CheckClosureImports_ClosureImport

func StaticGet_CheckClosureImports_ClosureImport_FORWARD_DECLARE() *CheckClosureImports_ClosureImport {
	clinit_CheckClosureImports_ClosureImport()
	return enum_CheckClosureImports_ClosureImport_FORWARD_DECLARE
}

var enum_CheckClosureImports_ClosureImport_REQUIRE_TYPE *CheckClosureImports_ClosureImport

func StaticGet_CheckClosureImports_ClosureImport_REQUIRE_TYPE() *CheckClosureImports_ClosureImport {
	clinit_CheckClosureImports_ClosureImport()
	return enum_CheckClosureImports_ClosureImport_REQUIRE_TYPE
}

var enum_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC *CheckClosureImports_ClosureImport

func StaticGet_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC() *CheckClosureImports_ClosureImport {
	clinit_CheckClosureImports_ClosureImport()
	return enum_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC
}

var clinitDone_CheckClosureImports_ClosureImport uint32
var clinitState_CheckClosureImports_ClosureImport j2gort.ClinitState

func clinit_CheckClosureImports_ClosureImport() {
	if atomic.LoadUint32(&clinitDone_CheckClosureImports_ClosureImport) != 0 {
		return
	}
	clinitSlow_CheckClosureImports_ClosureImport()
}

func clinitSlow_CheckClosureImports_ClosureImport() {
	if !j2gort.EnterClinit(&clinitState_CheckClosureImports_ClosureImport, &clinitDone_CheckClosureImports_ClosureImport) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckClosureImports_ClosureImport, &clinitDone_CheckClosureImports_ClosureImport)
	enum_CheckClosureImports_ClosureImport_REQUIRE = &CheckClosureImports_ClosureImport{}
	enum_CheckClosureImports_ClosureImport_REQUIRE.fn_AllowDestructuring = func(this_ *CheckClosureImports_ClosureImport) bool {
		return true
	}
	enum_CheckClosureImports_ClosureImport_REQUIRE.fn_AliasMustBeConstant = func(this_ *CheckClosureImports_ClosureImport) bool {
		return true
	}
	enum_CheckClosureImports_ClosureImport_REQUIRE.fn_MustBeOrdered = func(this_ *CheckClosureImports_ClosureImport) bool {
		return true
	}
	enum_CheckClosureImports_ClosureImport_REQUIRE.fn_CallName = func(this_ *CheckClosureImports_ClosureImport) string {
		return "goog.require"
	}
	initEnumCheckClosureImports_ClosureImport(enum_CheckClosureImports_ClosureImport_REQUIRE, "REQUIRE", int32(0))
	enum_CheckClosureImports_ClosureImport_FORWARD_DECLARE = &CheckClosureImports_ClosureImport{}
	enum_CheckClosureImports_ClosureImport_FORWARD_DECLARE.fn_AllowDestructuring = func(this_ *CheckClosureImports_ClosureImport) bool {
		return false
	}
	enum_CheckClosureImports_ClosureImport_FORWARD_DECLARE.fn_AliasMustBeConstant = func(this_ *CheckClosureImports_ClosureImport) bool {
		return false
	}
	enum_CheckClosureImports_ClosureImport_FORWARD_DECLARE.fn_MustBeOrdered = func(this_ *CheckClosureImports_ClosureImport) bool {
		return false
	}
	enum_CheckClosureImports_ClosureImport_FORWARD_DECLARE.fn_CallName = func(this_ *CheckClosureImports_ClosureImport) string {
		return "goog.forwardDeclare"
	}
	initEnumCheckClosureImports_ClosureImport(enum_CheckClosureImports_ClosureImport_FORWARD_DECLARE, "FORWARD_DECLARE", int32(1))
	enum_CheckClosureImports_ClosureImport_REQUIRE_TYPE = &CheckClosureImports_ClosureImport{}
	enum_CheckClosureImports_ClosureImport_REQUIRE_TYPE.fn_AllowDestructuring = func(this_ *CheckClosureImports_ClosureImport) bool {
		return true
	}
	enum_CheckClosureImports_ClosureImport_REQUIRE_TYPE.fn_AliasMustBeConstant = func(this_ *CheckClosureImports_ClosureImport) bool {
		return true
	}
	enum_CheckClosureImports_ClosureImport_REQUIRE_TYPE.fn_MustBeOrdered = func(this_ *CheckClosureImports_ClosureImport) bool {
		return false
	}
	enum_CheckClosureImports_ClosureImport_REQUIRE_TYPE.fn_CallName = func(this_ *CheckClosureImports_ClosureImport) string {
		return "goog.requireType"
	}
	initEnumCheckClosureImports_ClosureImport(enum_CheckClosureImports_ClosureImport_REQUIRE_TYPE, "REQUIRE_TYPE", int32(2))
	enum_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC = &CheckClosureImports_ClosureImport{}
	enum_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC.fn_AllowDestructuring = func(this_ *CheckClosureImports_ClosureImport) bool {
		return true
	}
	enum_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC.fn_AliasMustBeConstant = func(this_ *CheckClosureImports_ClosureImport) bool {
		return true
	}
	enum_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC.fn_MustBeOrdered = func(this_ *CheckClosureImports_ClosureImport) bool {
		return false
	}
	enum_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC.fn_CallName = func(this_ *CheckClosureImports_ClosureImport) string {
		return "goog.requireDynamic"
	}
	initEnumCheckClosureImports_ClosureImport(enum_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC, "REQUIRE_DYNAMIC", int32(3))
	enumValues_CheckClosureImports_ClosureImport = []*CheckClosureImports_ClosureImport{enum_CheckClosureImports_ClosureImport_REQUIRE, enum_CheckClosureImports_ClosureImport_FORWARD_DECLARE, enum_CheckClosureImports_ClosureImport_REQUIRE_TYPE, enum_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CheckClosureImports.ClosureImport", []any{enum_CheckClosureImports_ClosureImport_REQUIRE, enum_CheckClosureImports_ClosureImport_FORWARD_DECLARE, enum_CheckClosureImports_ClosureImport_REQUIRE_TYPE, enum_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CheckClosureImports$ClosureImport", []any{enum_CheckClosureImports_ClosureImport_REQUIRE, enum_CheckClosureImports_ClosureImport_FORWARD_DECLARE, enum_CheckClosureImports_ClosureImport_REQUIRE_TYPE, enum_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.CheckClosureImports.ClosureImport", "com.google.javascript.jscomp.CheckClosureImports$ClosureImport", func() {
		clinit_CheckClosureImports_ClosureImport()
	})
}

var enumValues_CheckClosureImports_ClosureImport []*CheckClosureImports_ClosureImport

func CheckClosureImports_ClosureImport_Values() []*CheckClosureImports_ClosureImport {
	clinit_CheckClosureImports_ClosureImport()
	return j2gort.CloneSlice(enumValues_CheckClosureImports_ClosureImport)
}

func CheckClosureImports_ClosureImport_ValueOf(name string) *CheckClosureImports_ClosureImport {
	clinit_CheckClosureImports_ClosureImport()
	for _, c := range enumValues_CheckClosureImports_ClosureImport {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CheckClosureImports.ClosureImport." + name))
}

// Java method: ClosureImport.allowDestructuring()
func (this_ *CheckClosureImports_ClosureImport) AllowDestructuring() bool {
	if this_.fn_AllowDestructuring != nil {
		return this_.fn_AllowDestructuring(this_)
	}
	return this_.baseEnum_AllowDestructuring()
}

func (this_ *CheckClosureImports_ClosureImport) baseEnum_AllowDestructuring() bool {
	panic("abstract or native method invoked")
}

// Java method: ClosureImport.aliasMustBeConstant()
func (this_ *CheckClosureImports_ClosureImport) AliasMustBeConstant() bool {
	if this_.fn_AliasMustBeConstant != nil {
		return this_.fn_AliasMustBeConstant(this_)
	}
	return this_.baseEnum_AliasMustBeConstant()
}

func (this_ *CheckClosureImports_ClosureImport) baseEnum_AliasMustBeConstant() bool {
	panic("abstract or native method invoked")
}

// Java method: ClosureImport.mustBeOrdered()
func (this_ *CheckClosureImports_ClosureImport) MustBeOrdered() bool {
	if this_.fn_MustBeOrdered != nil {
		return this_.fn_MustBeOrdered(this_)
	}
	return this_.baseEnum_MustBeOrdered()
}

func (this_ *CheckClosureImports_ClosureImport) baseEnum_MustBeOrdered() bool {
	panic("abstract or native method invoked")
}

// Java method: ClosureImport.callName()
func (this_ *CheckClosureImports_ClosureImport) CallName() string {
	if this_.fn_CallName != nil {
		return this_.fn_CallName(this_)
	}
	return this_.baseEnum_CallName()
}

func (this_ *CheckClosureImports_ClosureImport) baseEnum_CallName() string {
	panic("abstract or native method invoked")
}

func (this_ *CheckClosureImports_ClosureImport) IsJavaNil_() bool { return this_ == nil }
func (*CheckClosureImports_ClosureImport) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckClosureImports$ClosureImport", "ClosureImport"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckClosureImports.ClosureImport", "com.google.javascript.jscomp.CheckClosureImports$ClosureImport", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckClosureImports_ClosureImport](v_)
	})
}

// Java class: com.google.javascript.jscomp.CheckClosureImports.Checker
type CheckClosureImports_Checker struct {
	NodeTraversal_AbstractModuleCallbackBase
	outer_                                          *CheckClosureImports
	f_CheckClosureImports_Checker_moduleMetadataMap *ModuleMetadataMap `gson:"moduleMetadataMap"`
}

func (this_ *CheckClosureImports_Checker) BindSelf_(s any) {
	_ = s
	this_.NodeTraversal_AbstractModuleCallbackBase.BindSelf_(s)
}
func (this_ *CheckClosureImports_Checker) IsJavaNil_() bool { return this_ == nil }
func (*CheckClosureImports_Checker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckClosureImports$Checker", "Checker"
}
func (this_ *CheckClosureImports_Checker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractModuleCallbackBase.IdentityHash_()
}
func (*CheckClosureImports_Checker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckClosureImports.Checker", "com.google.javascript.jscomp.CheckClosureImports$Checker", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckClosureImports_Checker](v_)
	})
}
func (this_ *CheckClosureImports_Checker) SetOuter_(o *CheckClosureImports)                  { this_.outer_ = o }
func (this_ *CheckClosureImports_Checker) GetOuter_() *CheckClosureImports                   { return this_.outer_ }
func (*CheckClosureImports_Checker) Implements_jscomp_NodeTraversal_AbstractModuleCallback() {}
func (*CheckClosureImports_Checker) Implements_jscomp_NodeTraversal_Callback()               {}

func (this_ *CheckClosureImports_Checker) initFields_CheckClosureImports_Checker() {
	_ = this_
}

// Java constructor: Checker(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.modules.ModuleMetadataMap moduleMetadataMap)
func NewCheckClosureImports_Checker(outer_ *CheckClosureImports, compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *CheckClosureImports_Checker {
	this_ := &CheckClosureImports_Checker{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitCheckClosureImports_Checker(this_, compiler, moduleMetadataMap)
	return this_
}

func InitCheckClosureImports_Checker(this_ *CheckClosureImports_Checker, compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) {
	_ = compiler
	_ = moduleMetadataMap
	InitNodeTraversal_AbstractModuleCallback(&this_.NodeTraversal_AbstractModuleCallbackBase, compiler, moduleMetadataMap)
	this_.initFields_CheckClosureImports_Checker()
	this_.f_CheckClosureImports_Checker_moduleMetadataMap = moduleMetadataMap
}

// Java method: Checker.enterModule(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentModule, com.google.javascript.rhino.Node moduleScopeRoot)
func (this_ *CheckClosureImports_Checker) EnterModule(currentModule ModuleMetadataMap_ModuleMetadata, moduleScopeRoot pkg_rhino.Node) {
	_ = currentModule
	_ = moduleScopeRoot
	j2gort.OrderedSetAddAllStr(&this_.outer_.namespacesSeen, j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](currentModule.GoogNamespaces().ElementSet()))
}

// Java method: Checker.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.ModuleMetadataMap.@org.jspecify.annotations.Nullable ModuleMetadata currentModule, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node moduleScopeRoot)
func (this_ *CheckClosureImports_Checker) Visit_4(t *NodeTraversal, n pkg_rhino.Node, currentModule ModuleMetadataMap_ModuleMetadata, moduleScopeRoot pkg_rhino.Node) {
	_ = t
	_ = n
	_ = currentModule
	_ = moduleScopeRoot
	if currentModule == nil || currentModule.IsJavaNil_() {
		return
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if n.IsCall() {
		var kindOfCall *CheckClosureImports_ClosureImport = this_.priv_Checker_kindOfCall_1040(n)
		_ = kindOfCall
		if kindOfCall != nil {
			this_.priv_Checker_checkImport_1045(t, n, parent, currentModule, kindOfCall)
		} else if n.GetFirstChild().MatchesQualifiedName_1_Node(StaticGet_CheckClosureImports_GOOG_MODULE_GET()) {
			this_.priv_Checker_checkGoogModuleGet_1042(t, n, currentModule)
		}
	} else if n.IsName() {
		this_.priv_Checker_checkValidImportCodeReference_1041(t, n)
	}
}

// Java method: Checker.kindOfCall(com.google.javascript.rhino.Node callNode)
func (this_ *CheckClosureImports_Checker) priv_Checker_kindOfCall_1040(callNode pkg_rhino.Node) *CheckClosureImports_ClosureImport {
	_ = callNode
	if !(callNode.IsCall()) {
		j2gort.CheckState(false)
	}
	if callNode.GetFirstChild().MatchesQualifiedName_1_Node(StaticGet_CheckClosureImports_GOOG_REQUIRE()) {
		return StaticGet_CheckClosureImports_ClosureImport_REQUIRE()
	} else if callNode.GetFirstChild().MatchesQualifiedName_1_Node(StaticGet_CheckClosureImports_GOOG_FORWARD_DECLARE()) {
		return StaticGet_CheckClosureImports_ClosureImport_FORWARD_DECLARE()
	} else if callNode.GetFirstChild().MatchesQualifiedName_1_Node(StaticGet_CheckClosureImports_GOOG_REQUIRE_TYPE()) {
		return StaticGet_CheckClosureImports_ClosureImport_REQUIRE_TYPE()
	} else if callNode.GetFirstChild().MatchesQualifiedName_1_Node(StaticGet_CheckClosureImports_GOOG_REQUIRE_DYNAMIC()) {
		return StaticGet_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC()
	}
	return nil
}

// Java method: Checker.checkValidImportCodeReference(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nameNode)
func (this_ *CheckClosureImports_Checker) priv_Checker_checkValidImportCodeReference_1041(t *NodeTraversal, nameNode pkg_rhino.Node) {
	_ = t
	_ = nameNode
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(nameNode.GetString()))
	_ = var_
	if (var_ == nil) || (var_.GetNameNode() == nameNode) {
		return
	}
	var declarationNameNode pkg_rhino.Node = var_.GetNameNode()
	_ = declarationNameNode
	if (declarationNameNode == 0) || !(NodeUtil_IsDeclarationLValue(declarationNameNode)) {
		return
	}
	var decl pkg_rhino.Node
	_ = decl
	if declarationNameNode.GetParent().IsStringKey() {
		decl = declarationNameNode.GetGrandparent().GetGrandparent()
	} else {
		decl = declarationNameNode.GetParent()
	}
	if ((!(NodeUtil_IsNameDeclaration(decl)) || !(decl.HasOneChild())) || !(decl.GetFirstChild().HasChildren())) || !(decl.GetFirstChild().GetLastChild().IsCall()) {
		return
	}
	var kindOfCall *CheckClosureImports_ClosureImport = this_.priv_Checker_kindOfCall_1040(decl.GetFirstChild().GetLastChild())
	_ = kindOfCall
	if kindOfCall == nil {
		return
	}
}

// Java method: Checker.checkGoogModuleGet(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentModule)
func (this_ *CheckClosureImports_Checker) priv_Checker_checkGoogModuleGet_1042(t *NodeTraversal, call pkg_rhino.Node, currentModule ModuleMetadataMap_ModuleMetadata) {
	_ = t
	_ = call
	_ = currentModule
	if t.InModuleHoistScope() {
		t.Report_3(call, StaticGet_ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET(), []string{})
		return
	}
	if (!(currentModule.IsEs6Module()) && !(currentModule.IsGoogModule())) && (this_.GetNodeTraversal_AbstractModuleCallbackBase().compiler.GetOptions().GetModuleResolutionMode() != StaticGet_ModuleLoader_ResolutionMode_WEBPACK()) {
		if t.InGlobalScope() {
			var getSubExprRoot pkg_rhino.Node = call
			_ = getSubExprRoot
			for NodeUtil_IsNormalOrOptChainGetProp(getSubExprRoot.GetParent()) {
				getSubExprRoot = getSubExprRoot.GetParent()
			}
			var invalid bool = false
			_ = invalid
			var parent pkg_rhino.Node = getSubExprRoot.GetParent()
			_ = parent
			if parent.IsName() || parent.IsDestructuringLhs() {
				invalid = true
			} else if parent.IsAssign() {
				var target pkg_rhino.Node = parent.GetFirstChild()
				_ = target
				if target.IsName() || target.IsObjectPattern() {
					invalid = true
				}
			}
			if invalid {
				t.Report_3(call, StaticGet_ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE(), []string{})
			}
		}
	}
	if !(call.HasTwoChildren()) || !(call.GetSecondChild().IsStringLit()) {
		t.Report_3(call, StaticGet_CheckClosureImports_INVALID_CLOSURE_IMPORT_CALL(), []string{"goog.module.get"})
		return
	}
	var namespace string = call.GetSecondChild().GetString()
	_ = namespace
	var requiredModule ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.f_CheckClosureImports_Checker_moduleMetadataMap.GetModulesByGoogNamespace(), namespace))
	_ = requiredModule
	if requiredModule == nil || requiredModule.IsJavaNil_() {
		t.Report_3(call, StaticGet_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE(), []string{namespace})
		return
	}
	var maybeAssign pkg_rhino.Node = call.GetParent()
	_ = maybeAssign
	var isFillingAnAlias bool = (maybeAssign.IsAssign() && maybeAssign.GetFirstChild().IsName())
	_ = isFillingAnAlias
	var isModule bool = (currentModule.IsGoogModule() || currentModule.IsEs6Module())
	_ = isModule
	if isFillingAnAlias && isModule {
		var aliasName string = call.GetParent().GetFirstChild().GetString()
		_ = aliasName
		var aliasVar *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(aliasName))
		_ = aliasVar
		if aliasVar == nil {
			t.Report_3(call, StaticGet_ClosureRewriteModule_INVALID_GET_ALIAS(), []string{})
			return
		}
		var aliasVarNodeRhs pkg_rhino.Node = NodeUtil_GetRValueOfLValue(aliasVar.GetNode())
		_ = aliasVarNodeRhs
		if (aliasVarNodeRhs == 0) || !(NodeUtil_IsCallTo_2_Node_Node(aliasVarNodeRhs, StaticGet_CheckClosureImports_GOOG_FORWARD_DECLARE())) {
			t.Report_3(call, StaticGet_ClosureRewriteModule_INVALID_GET_ALIAS(), []string{})
			return
		}
		{
			var iifeCond_1 bool
			{

				sr_1 := namespace
				se_2 := aliasVarNodeRhs.GetLastChild().GetString()
				iifeCond_1 = j2gort.CheckNotNullString(sr_1) == se_2
			}
			if !(iifeCond_1) {
				t.Report_3(call, StaticGet_ClosureRewriteModule_INVALID_GET_ALIAS(), []string{})
				return
			}
		}
	}
}

// Java method: Checker.isValidDestructuringImport(com.google.javascript.rhino.Node destructuringLhs)
func (this_ *CheckClosureImports_Checker) priv_Checker_isValidDestructuringImport_1043(destructuringLhs pkg_rhino.Node) bool {
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
			}
			if !(stringKey.HasChildren()) || !(stringKey.GetFirstChild().IsName()) {
				return false
			}
		}
	}
	return true
}

// Java method: Checker.checkShortName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node shortNameNode, java.lang.String namespace)
func (this_ *CheckClosureImports_Checker) priv_Checker_checkShortName_1044(t *NodeTraversal, shortNameNode pkg_rhino.Node, namespace string) {
	_ = t
	_ = shortNameNode
	_ = namespace
	var shortName string = shortNameNode.GetString()
	_ = shortName
	var lastSegment string = j2gort.JStrSubFrom(namespace, (j2gort.JStrLastIndexOfCp(namespace, int32(uint16(46))) + int32(1)))
	_ = lastSegment
	if (j2gort.CheckNotNullString(shortName) == lastSegment) || (j2gort.CheckNotNullString(lastSegment) == "") {
		return
	}
	if j2gort.Ext_com_google_common_base_Ascii_IsUpperCase(j2gort.JStrCharAt(shortName, int32(0))) != j2gort.Ext_com_google_common_base_Ascii_IsUpperCase(j2gort.JStrCharAt(lastSegment, int32(0))) {
		var newStartChar uint16
		_ = newStartChar
		if j2gort.Ext_com_google_common_base_Ascii_IsUpperCase(j2gort.JStrCharAt(shortName, int32(0))) {
			newStartChar = j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_uint16(j2gort.JStrCharAt(shortName, int32(0)))
		} else {
			newStartChar = j2gort.Ext_com_google_common_base_Ascii_ToUpperCase_1_uint16(j2gort.JStrCharAt(shortName, int32(0)))
		}
		var correctedName string = (j2gort.StrValChar(newStartChar) + j2gort.StrValString(j2gort.JStrSubFrom(shortName, int32(1))))
		_ = correctedName
		t.Report_3(shortNameNode, StaticGet_ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION(), []string{shortName, correctedName})
	}
}

// Java method: Checker.checkImport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call, com.google.javascript.rhino.Node parent, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentModule, com.google.javascript.jscomp.CheckClosureImports.ClosureImport importType)
func (this_ *CheckClosureImports_Checker) priv_Checker_checkImport_1045(t *NodeTraversal, call pkg_rhino.Node, parent pkg_rhino.Node, currentModule ModuleMetadataMap_ModuleMetadata, importType *CheckClosureImports_ClosureImport) {
	_ = t
	_ = call
	_ = parent
	_ = currentModule
	_ = importType
	var atTopLevelScope bool = (t.InGlobalHoistScope() || t.InModuleScope())
	_ = atTopLevelScope
	var isModule bool = currentModule.IsModule()
	_ = isModule
	var validAssignment bool = (isModule || parent.IsExprResult())
	_ = validAssignment
	var isAliased bool = NodeUtil_IsNameDeclaration(parent.GetParent())
	_ = isAliased
	if (!(atTopLevelScope) || !(validAssignment)) && (importType != StaticGet_CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC()) {
		t.Report_3(call, StaticGet_ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR(), []string{})
		return
	}
	if !(call.HasTwoChildren()) || !(call.GetSecondChild().IsStringLit()) {
		t.Report_3(call, StaticGet_CheckClosureImports_INVALID_CLOSURE_IMPORT_CALL(), []string{importType.CallName()})
		return
	}
	if isModule && isAliased {
		var declaration pkg_rhino.Node = parent.GetParent()
		_ = declaration
		if !(declaration.HasOneChild()) {
			t.Report_3(call, StaticGet_CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION(), []string{importType.CallName()})
			return
		}
		if importType.AliasMustBeConstant() {
			if declaration.IsLet() {
				t.Report_3(call, StaticGet_CheckClosureImports_LET_CLOSURE_IMPORT(), []string{importType.CallName()})
			} else if !(declaration.IsConst()) && currentModule.IsEs6Module() {
				t.Report_3(call, StaticGet_CheckClosureImports_LHS_OF_CLOSURE_IMPORT_MUST_BE_CONST_IN_ES_MODULE(), []string{importType.CallName()})
			}
		}
		var lhs pkg_rhino.Node = declaration.GetFirstChild()
		_ = lhs
		if lhs.IsDestructuringLhs() {
			if importType.AllowDestructuring() {
				if !(this_.priv_Checker_isValidDestructuringImport_1043(lhs)) {
					t.Report_3(declaration, StaticGet_CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING(), []string{importType.CallName()})
				}
			} else {
				t.Report_3(declaration, StaticGet_CheckClosureImports_NO_CLOSURE_IMPORT_DESTRUCTURING(), []string{importType.CallName()})
			}
		} else if lhs.IsName() {
			this_.priv_Checker_checkShortName_1044(t, lhs, call.GetLastChild().GetString())
		}
	}
	var namespace string = call.GetSecondChild().GetString()
	_ = namespace
	var requiredModule ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.f_CheckClosureImports_Checker_moduleMetadataMap.GetModulesByGoogNamespace(), namespace))
	_ = requiredModule
	if requiredModule == nil || requiredModule.IsJavaNil_() {
		if importType == StaticGet_CheckClosureImports_ClosureImport_FORWARD_DECLARE() {
			if isAliased && isModule {
				t.Report_3(call, StaticGet_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE(), []string{namespace})
			}
			return
		}
		t.Report_3(call, StaticGet_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE(), []string{namespace})
		return
	}
	if importType.MustBeOrdered() {
		this_.outer_.priv_CheckClosureImports_verifyRequireOrder_1033(namespace, call, t, requiredModule)
	}
	if ((importType == StaticGet_CheckClosureImports_ClosureImport_REQUIRE()) && currentModule.IsEs6Module()) && requiredModule.IsEs6Module() {
		t.Report_3(call, StaticGet_Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE(), []string{})
	}
}
func (*CheckClosureImports_Checker) CanReleaseSyntacticScopes_() {}
