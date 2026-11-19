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
// Port of com/google/javascript/jscomp/GatherModuleMetadata.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.GatherModuleMetadata
type GatherModuleMetadata struct {
	modulesByPath          j2gort.Ext_java_util_Map                    `gson:"modulesByPath"`
	modulesByGoogNamespace j2gort.Ext_java_util_Map                    `gson:"modulesByGoogNamespace"`
	currentModule          *GatherModuleMetadata_ModuleMetadataBuilder `gson:"currentModule"`
	parentModule           *GatherModuleMetadata_ModuleMetadataBuilder `gson:"parentModule"`
	compiler               AbstractCompiler                            `gson:"compiler"`
	loadModuleCall         pkg_rhino.Node                              `gson:"loadModuleCall"`
	idHash_                atomic.Int32
	processCommonJsModules bool                        `gson:"processCommonJsModules"`
	moduleResolutionMode   ModuleLoader_ResolutionMode `gson:"moduleResolutionMode"`
}

func (this_ *GatherModuleMetadata) BindSelf_(s any) {
	_ = s
}
func (this_ *GatherModuleMetadata) IsJavaNil_() bool { return this_ == nil }
func (*GatherModuleMetadata) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GatherModuleMetadata", "GatherModuleMetadata"
}
func (this_ *GatherModuleMetadata) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GatherModuleMetadata) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GatherModuleMetadata", "com.google.javascript.jscomp.GatherModuleMetadata", func(v_ any) bool {
		return j2gort.InstanceOf[*GatherModuleMetadata](v_)
	})
}
func (*GatherModuleMetadata) Implements_jscomp_CompilerPass() {}

var static_GatherModuleMetadata_MIXED_MODULE_TYPE *DiagnosticType

func StaticGet_GatherModuleMetadata_MIXED_MODULE_TYPE() *DiagnosticType {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_MIXED_MODULE_TYPE
}
func StaticSet_GatherModuleMetadata_MIXED_MODULE_TYPE(v *DiagnosticType) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_MIXED_MODULE_TYPE = v
}

var static_GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID *DiagnosticType

func StaticGet_GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID() *DiagnosticType {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID
}
func StaticSet_GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID(v *DiagnosticType) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID = v
}

var static_GatherModuleMetadata_INVALID_DECLARE_MODULE_ID_CALL *DiagnosticType

func StaticGet_GatherModuleMetadata_INVALID_DECLARE_MODULE_ID_CALL() *DiagnosticType {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_INVALID_DECLARE_MODULE_ID_CALL
}
func StaticSet_GatherModuleMetadata_INVALID_DECLARE_MODULE_ID_CALL(v *DiagnosticType) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_INVALID_DECLARE_MODULE_ID_CALL = v
}

var static_GatherModuleMetadata_DECLARE_MODULE_ID_OUTSIDE_ES6_MODULE *DiagnosticType

func StaticGet_GatherModuleMetadata_DECLARE_MODULE_ID_OUTSIDE_ES6_MODULE() *DiagnosticType {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_DECLARE_MODULE_ID_OUTSIDE_ES6_MODULE
}
func StaticSet_GatherModuleMetadata_DECLARE_MODULE_ID_OUTSIDE_ES6_MODULE(v *DiagnosticType) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_DECLARE_MODULE_ID_OUTSIDE_ES6_MODULE = v
}

var static_GatherModuleMetadata_MULTIPLE_DECLARE_MODULE_NAMESPACE *DiagnosticType

func StaticGet_GatherModuleMetadata_MULTIPLE_DECLARE_MODULE_NAMESPACE() *DiagnosticType {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_MULTIPLE_DECLARE_MODULE_NAMESPACE
}
func StaticSet_GatherModuleMetadata_MULTIPLE_DECLARE_MODULE_NAMESPACE(v *DiagnosticType) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_MULTIPLE_DECLARE_MODULE_NAMESPACE = v
}

var static_GatherModuleMetadata_INVALID_REQUIRE_TYPE *DiagnosticType

func StaticGet_GatherModuleMetadata_INVALID_REQUIRE_TYPE() *DiagnosticType {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_INVALID_REQUIRE_TYPE
}
func StaticSet_GatherModuleMetadata_INVALID_REQUIRE_TYPE(v *DiagnosticType) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_INVALID_REQUIRE_TYPE = v
}

var static_GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC *DiagnosticType

func StaticGet_GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC() *DiagnosticType {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC
}
func StaticSet_GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC(v *DiagnosticType) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC = v
}

var static_GatherModuleMetadata_INVALID_MAYBE_REQUIRE *DiagnosticType

func StaticGet_GatherModuleMetadata_INVALID_MAYBE_REQUIRE() *DiagnosticType {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_INVALID_MAYBE_REQUIRE
}
func StaticSet_GatherModuleMetadata_INVALID_MAYBE_REQUIRE(v *DiagnosticType) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_INVALID_MAYBE_REQUIRE = v
}

var static_GatherModuleMetadata_INVALID_SET_TEST_ONLY *DiagnosticType

func StaticGet_GatherModuleMetadata_INVALID_SET_TEST_ONLY() *DiagnosticType {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_INVALID_SET_TEST_ONLY
}
func StaticSet_GatherModuleMetadata_INVALID_SET_TEST_ONLY(v *DiagnosticType) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_INVALID_SET_TEST_ONLY = v
}

var static_GatherModuleMetadata_INVALID_NESTED_LOAD_MODULE *DiagnosticType

func StaticGet_GatherModuleMetadata_INVALID_NESTED_LOAD_MODULE() *DiagnosticType {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_INVALID_NESTED_LOAD_MODULE
}
func StaticSet_GatherModuleMetadata_INVALID_NESTED_LOAD_MODULE(v *DiagnosticType) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_INVALID_NESTED_LOAD_MODULE = v
}

var static_GatherModuleMetadata_INVALID_TOGGLE_USAGE *DiagnosticType

func StaticGet_GatherModuleMetadata_INVALID_TOGGLE_USAGE() *DiagnosticType {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_INVALID_TOGGLE_USAGE
}
func StaticSet_GatherModuleMetadata_INVALID_TOGGLE_USAGE(v *DiagnosticType) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_INVALID_TOGGLE_USAGE = v
}

var static_GatherModuleMetadata_GOOG_PROVIDE pkg_rhino.Node

func StaticGet_GatherModuleMetadata_GOOG_PROVIDE() pkg_rhino.Node {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_PROVIDE
}
func StaticSet_GatherModuleMetadata_GOOG_PROVIDE(v pkg_rhino.Node) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_PROVIDE = v
}

var static_GatherModuleMetadata_GOOG_MODULE pkg_rhino.Node

func StaticGet_GatherModuleMetadata_GOOG_MODULE() pkg_rhino.Node {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_MODULE
}
func StaticSet_GatherModuleMetadata_GOOG_MODULE(v pkg_rhino.Node) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_MODULE = v
}

var static_GatherModuleMetadata_GOOG_MODULE_GET pkg_rhino.Node

func StaticGet_GatherModuleMetadata_GOOG_MODULE_GET() pkg_rhino.Node {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_MODULE_GET
}
func StaticSet_GatherModuleMetadata_GOOG_MODULE_GET(v pkg_rhino.Node) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_MODULE_GET = v
}

var static_GatherModuleMetadata_GOOG_REQUIRE pkg_rhino.Node

func StaticGet_GatherModuleMetadata_GOOG_REQUIRE() pkg_rhino.Node {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_REQUIRE
}
func StaticSet_GatherModuleMetadata_GOOG_REQUIRE(v pkg_rhino.Node) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_REQUIRE = v
}

var static_GatherModuleMetadata_GOOG_REQUIRE_TYPE pkg_rhino.Node

func StaticGet_GatherModuleMetadata_GOOG_REQUIRE_TYPE() pkg_rhino.Node {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_REQUIRE_TYPE
}
func StaticSet_GatherModuleMetadata_GOOG_REQUIRE_TYPE(v pkg_rhino.Node) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_REQUIRE_TYPE = v
}

var static_GatherModuleMetadata_GOOG_MAYBE_REQUIRE pkg_rhino.Node

func StaticGet_GatherModuleMetadata_GOOG_MAYBE_REQUIRE() pkg_rhino.Node {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_MAYBE_REQUIRE
}
func StaticSet_GatherModuleMetadata_GOOG_MAYBE_REQUIRE(v pkg_rhino.Node) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_MAYBE_REQUIRE = v
}

var static_GatherModuleMetadata_GOOG_REQUIRE_DYNAMIC pkg_rhino.Node

func StaticGet_GatherModuleMetadata_GOOG_REQUIRE_DYNAMIC() pkg_rhino.Node {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_REQUIRE_DYNAMIC
}
func StaticSet_GatherModuleMetadata_GOOG_REQUIRE_DYNAMIC(v pkg_rhino.Node) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_REQUIRE_DYNAMIC = v
}

var static_GatherModuleMetadata_GOOG_SET_TEST_ONLY pkg_rhino.Node

func StaticGet_GatherModuleMetadata_GOOG_SET_TEST_ONLY() pkg_rhino.Node {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_SET_TEST_ONLY
}
func StaticSet_GatherModuleMetadata_GOOG_SET_TEST_ONLY(v pkg_rhino.Node) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_SET_TEST_ONLY = v
}

var static_GatherModuleMetadata_GOOG_MODULE_DECLARELEGACYNAMESPACE pkg_rhino.Node

func StaticGet_GatherModuleMetadata_GOOG_MODULE_DECLARELEGACYNAMESPACE() pkg_rhino.Node {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_MODULE_DECLARELEGACYNAMESPACE
}
func StaticSet_GatherModuleMetadata_GOOG_MODULE_DECLARELEGACYNAMESPACE(v pkg_rhino.Node) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_MODULE_DECLARELEGACYNAMESPACE = v
}

var static_GatherModuleMetadata_GOOG_DECLARE_MODULE_ID pkg_rhino.Node

func StaticGet_GatherModuleMetadata_GOOG_DECLARE_MODULE_ID() pkg_rhino.Node {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_DECLARE_MODULE_ID
}
func StaticSet_GatherModuleMetadata_GOOG_DECLARE_MODULE_ID(v pkg_rhino.Node) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_DECLARE_MODULE_ID = v
}

var static_GatherModuleMetadata_GOOG_MODULE_DECLARNAMESPACE pkg_rhino.Node

func StaticGet_GatherModuleMetadata_GOOG_MODULE_DECLARNAMESPACE() pkg_rhino.Node {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_MODULE_DECLARNAMESPACE
}
func StaticSet_GatherModuleMetadata_GOOG_MODULE_DECLARNAMESPACE(v pkg_rhino.Node) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_MODULE_DECLARNAMESPACE = v
}
func StaticGet_GatherModuleMetadata_TOGGLE_NAME_PREFIX() string { return "TOGGLE_" }

var static_GatherModuleMetadata_GOOG_LOADMODULE pkg_rhino.QualifiedName

func StaticGet_GatherModuleMetadata_GOOG_LOADMODULE() pkg_rhino.QualifiedName {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_GOOG_LOADMODULE
}
func StaticSet_GatherModuleMetadata_GOOG_LOADMODULE(v pkg_rhino.QualifiedName) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_GOOG_LOADMODULE = v
}

var static_GatherModuleMetadata_NAMESPACE_SEGMENT_REGEX *j2gort.Ext_java_util_regex_Pattern

func StaticGet_GatherModuleMetadata_NAMESPACE_SEGMENT_REGEX() *j2gort.Ext_java_util_regex_Pattern {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_NAMESPACE_SEGMENT_REGEX
}
func StaticSet_GatherModuleMetadata_NAMESPACE_SEGMENT_REGEX(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_NAMESPACE_SEGMENT_REGEX = v
}

var static_GatherModuleMetadata_DOT_SPLITTER *j2gort.Ext_com_google_common_base_Splitter

func StaticGet_GatherModuleMetadata_DOT_SPLITTER() *j2gort.Ext_com_google_common_base_Splitter {
	clinit_GatherModuleMetadata()
	return static_GatherModuleMetadata_DOT_SPLITTER
}
func StaticSet_GatherModuleMetadata_DOT_SPLITTER(v *j2gort.Ext_com_google_common_base_Splitter) {
	clinit_GatherModuleMetadata()
	static_GatherModuleMetadata_DOT_SPLITTER = v
}

var clinitDone_GatherModuleMetadata uint32
var clinitState_GatherModuleMetadata j2gort.ClinitState

func clinit_GatherModuleMetadata() {
	if atomic.LoadUint32(&clinitDone_GatherModuleMetadata) != 0 {
		return
	}
	clinitSlow_GatherModuleMetadata()
}

func clinitSlow_GatherModuleMetadata() {
	if !j2gort.EnterClinit(&clinitState_GatherModuleMetadata, &clinitDone_GatherModuleMetadata) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_GatherModuleMetadata, &clinitDone_GatherModuleMetadata)
	static_GatherModuleMetadata_MIXED_MODULE_TYPE = DiagnosticType_Error("JSC_MIXED_MODULE_TYPE", "A file cannot be both {0} and {1}.")
	static_GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID = DiagnosticType_Error("JSC_INVALID_NAMESPACE_OR_MODULE_ID", "Namespace and module ID must be a dot-separated sequence of legal property identifiers and must only contain ASCII, 0-9, $, ., and _. Found ''{0}''")
	static_GatherModuleMetadata_INVALID_DECLARE_MODULE_ID_CALL = DiagnosticType_Error("JSC_INVALID_DECLARE_NAMESPACE_CALL", "goog.declareModuleId parameter must be a string literal.")
	static_GatherModuleMetadata_DECLARE_MODULE_ID_OUTSIDE_ES6_MODULE = DiagnosticType_Error("JSC_DECLARE_MODULE_NAMESPACE_OUTSIDE_ES6_MODULE", "goog.declareModuleId can only be called within ES6 modules.")
	static_GatherModuleMetadata_MULTIPLE_DECLARE_MODULE_NAMESPACE = DiagnosticType_Error("JSC_MULTIPLE_DECLARE_MODULE_NAMESPACE", "goog.declareModuleId can only be called once per ES6 module.")
	static_GatherModuleMetadata_INVALID_REQUIRE_TYPE = DiagnosticType_Error("JSC_INVALID_REQUIRE_TYPE", "Argument to goog.requireType must be a string.")
	static_GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC = DiagnosticType_Error("JSC_INVALID_REQUIRE_DYNAMIC", "Argument to goog.requireDynamic must be a string.")
	static_GatherModuleMetadata_INVALID_MAYBE_REQUIRE = DiagnosticType_Error("JSC_INVALID_MAYBE_REQUIRE", "Argument to goog.maybeRequireFrameworkInternalOnlyDoNotCallOrElse must be a string.")
	static_GatherModuleMetadata_INVALID_SET_TEST_ONLY = DiagnosticType_Error("JSC_INVALID_SET_TEST_ONLY", "Optional, single argument to goog.setTestOnly must be a string.")
	static_GatherModuleMetadata_INVALID_NESTED_LOAD_MODULE = DiagnosticType_Error("JSC_INVALID_NESTED_LOAD_MODULE", "goog.loadModule cannot be nested.")
	static_GatherModuleMetadata_INVALID_TOGGLE_USAGE = DiagnosticType_Error("JSC_INVALID_TOGGLE_USAGE", "Invalid toggle usage: {0}")
	static_GatherModuleMetadata_GOOG_PROVIDE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "provide")
	static_GatherModuleMetadata_GOOG_MODULE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "module")
	static_GatherModuleMetadata_GOOG_MODULE_GET = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "module"), "get")
	static_GatherModuleMetadata_GOOG_REQUIRE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "require")
	static_GatherModuleMetadata_GOOG_REQUIRE_TYPE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "requireType")
	static_GatherModuleMetadata_GOOG_MAYBE_REQUIRE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "maybeRequireFrameworkInternalOnlyDoNotCallOrElse")
	static_GatherModuleMetadata_GOOG_REQUIRE_DYNAMIC = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "requireDynamic")
	static_GatherModuleMetadata_GOOG_SET_TEST_ONLY = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "setTestOnly")
	static_GatherModuleMetadata_GOOG_MODULE_DECLARELEGACYNAMESPACE = pkg_rhino.IR_Getprop_2(StaticGet_GatherModuleMetadata_GOOG_MODULE().CloneTree_0(), "declareLegacyNamespace")
	static_GatherModuleMetadata_GOOG_DECLARE_MODULE_ID = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "declareModuleId")
	static_GatherModuleMetadata_GOOG_MODULE_DECLARNAMESPACE = pkg_rhino.IR_Getprop_2(StaticGet_GatherModuleMetadata_GOOG_MODULE().CloneTree_0(), "declareNamespace")
	static_GatherModuleMetadata_GOOG_LOADMODULE = pkg_rhino.QualifiedName_Of("goog.loadModule")
	static_GatherModuleMetadata_NAMESPACE_SEGMENT_REGEX = j2gort.Ext_java_util_regex_Pattern_Compile_1("^[a-zA-Z_$][a-zA-Z0-9_$.]*$")
	static_GatherModuleMetadata_DOT_SPLITTER = j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(46))
}

func (this_ *GatherModuleMetadata) initFields_GatherModuleMetadata() {
	_ = this_
	this_.modulesByPath = j2gort.NewExt_java_util_LinkedHashMap()
	this_.modulesByGoogNamespace = j2gort.NewExt_java_util_LinkedHashMap()
}

// Java constructor: GatherModuleMetadata(com.google.javascript.jscomp.AbstractCompiler compiler, boolean processCommonJsModules, com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode moduleResolutionMode)
func NewGatherModuleMetadata(compiler AbstractCompiler, processCommonJsModules bool, moduleResolutionMode ModuleLoader_ResolutionMode) *GatherModuleMetadata {
	clinit_GatherModuleMetadata()
	this_ := &GatherModuleMetadata{}
	InitGatherModuleMetadata(this_, compiler, processCommonJsModules, moduleResolutionMode)
	return this_
}

func InitGatherModuleMetadata(this_ *GatherModuleMetadata, compiler AbstractCompiler, processCommonJsModules bool, moduleResolutionMode ModuleLoader_ResolutionMode) {
	_ = compiler
	_ = processCommonJsModules
	_ = moduleResolutionMode
	this_.initFields_GatherModuleMetadata()
	this_.compiler = compiler
	this_.processCommonJsModules = processCommonJsModules
	this_.moduleResolutionMode = moduleResolutionMode
}

// Java method: GatherModuleMetadata.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *GatherModuleMetadata) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer sc_.ReleaseSlabs_()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(NewGatherModuleMetadata_Finder(this_))
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).Traverse(externs)
	}
	{
		clinit_NodeTraversal()
		var ntb_1_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_1_)
		(&ntb_1_).SetCompiler(this_.compiler)
		(&ntb_1_).SetCallback_1_NodeTraversal_Callback(NewGatherModuleMetadata_Finder(this_))
		(&ntb_1_).SetScopeCreator(sc_)
		(&ntb_1_).Traverse(root)
	}

	this_.compiler.SetModuleMetadataMap(NewModuleMetadataMap(this_.modulesByPath, this_.modulesByGoogNamespace))
}

// Java static method: GatherModuleMetadata.isValidModuleId(java.lang.String id)
func GatherModuleMetadata_IsValidModuleId(id string) bool {
	clinit_GatherModuleMetadata()
	_ = id
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_GatherModuleMetadata_DOT_SPLITTER().Split(j2gort.BoxString(id))))
		for cur_1.HasNext() {
			segment := cur_1.NextStr()
			_ = segment
			if j2gort.CheckNotNullString(segment) == "" {
				return false
			}
			{
				var i int32 = int32(0)
				_ = i
				for ; i < j2gort.JStrLen(segment); i = i + 1 {
					if !(pkg_parser.Identifiers_IsIdentifierPart(j2gort.JStrCharAt(segment, i))) {
						return false
					}
				}
			}
		}
	}
	return StaticGet_GatherModuleMetadata_NAMESPACE_SEGMENT_REGEX().Matcher(j2gort.BoxString(id)).Matches()
}

// Java class: com.google.javascript.jscomp.GatherModuleMetadata.ModuleMetadataBuilder
type GatherModuleMetadata_ModuleMetadataBuilder struct {
	outer_                  *GatherModuleMetadata
	metadataBuilder         ModuleMetadataMap_ModuleMetadata_Builder                 `gson:"metadataBuilder"`
	googNamespaces          *j2gort.Ext_com_google_common_collect_LinkedHashMultiset `gson:"googNamespaces"`
	declaredModuleId        pkg_rhino.Node                                           `gson:"declaredModuleId"`
	declaresLegacyNamespace pkg_rhino.Node                                           `gson:"declaresLegacyNamespace"`
	idHash_                 atomic.Int32
	ambiguous               bool `gson:"ambiguous"`
	hasModuleBody           bool `gson:"hasModuleBody"`
}

func (this_ *GatherModuleMetadata_ModuleMetadataBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *GatherModuleMetadata_ModuleMetadataBuilder) IsJavaNil_() bool { return this_ == nil }
func (*GatherModuleMetadata_ModuleMetadataBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GatherModuleMetadata$ModuleMetadataBuilder", "ModuleMetadataBuilder"
}
func (this_ *GatherModuleMetadata_ModuleMetadataBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GatherModuleMetadata_ModuleMetadataBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GatherModuleMetadata.ModuleMetadataBuilder", "com.google.javascript.jscomp.GatherModuleMetadata$ModuleMetadataBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*GatherModuleMetadata_ModuleMetadataBuilder](v_)
	})
}
func (this_ *GatherModuleMetadata_ModuleMetadataBuilder) SetOuter_(o *GatherModuleMetadata) {
	this_.outer_ = o
}
func (this_ *GatherModuleMetadata_ModuleMetadataBuilder) GetOuter_() *GatherModuleMetadata {
	return this_.outer_
}

func (this_ *GatherModuleMetadata_ModuleMetadataBuilder) initFields_GatherModuleMetadata_ModuleMetadataBuilder() {
	_ = this_
	this_.googNamespaces = j2gort.Ext_com_google_common_collect_LinkedHashMultiset_Create_0()
}

// Java constructor: ModuleMetadataBuilder(com.google.javascript.rhino.Node rootNode, com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath path)
func NewGatherModuleMetadata_ModuleMetadataBuilder(outer_ *GatherModuleMetadata, rootNode pkg_rhino.Node, path *ModuleLoader_ModulePath) *GatherModuleMetadata_ModuleMetadataBuilder {
	this_ := &GatherModuleMetadata_ModuleMetadataBuilder{}
	this_.outer_ = outer_
	InitGatherModuleMetadata_ModuleMetadataBuilder(this_, rootNode, path)
	return this_
}

func InitGatherModuleMetadata_ModuleMetadataBuilder(this_ *GatherModuleMetadata_ModuleMetadataBuilder, rootNode pkg_rhino.Node, path *ModuleLoader_ModulePath) {
	_ = rootNode
	_ = path
	this_.initFields_GatherModuleMetadata_ModuleMetadataBuilder()
	this_.metadataBuilder = ModuleMetadataMap_ModuleMetadata_S_Builder().Path(path).RootNode(rootNode).ModuleType_1(StaticGet_ModuleMetadataMap_ModuleType_SCRIPT()).UsesClosure(false).IsTestOnly(false)
}

// Java method: ModuleMetadataBuilder.moduleType(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleType type, com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *GatherModuleMetadata_ModuleMetadataBuilder) ModuleType(type_ ModuleMetadataMap_ModuleType, t *NodeTraversal, n pkg_rhino.Node) {
	_ = type_
	_ = t
	_ = n
	j2gort.CheckNotNullEnum(type_)
	if this_.metadataBuilder.ModuleType_0() == type_ {
		return
	}
	if this_.metadataBuilder.ModuleType_0() == StaticGet_ModuleMetadataMap_ModuleType_SCRIPT() {
		this_.metadataBuilder.ModuleType_1(type_)
		return
	}
	this_.ambiguous = true
	t.Report_3(n, StaticGet_GatherModuleMetadata_MIXED_MODULE_TYPE(), []string{(this_.metadataBuilder.ModuleType_0()).FieldGet_description(), type_.FieldGet_description()})
}

// Java method: ModuleMetadataBuilder.recordDeclareModuleId(com.google.javascript.rhino.Node declaredModuleId)
func (this_ *GatherModuleMetadata_ModuleMetadataBuilder) RecordDeclareModuleId(declaredModuleId pkg_rhino.Node) {
	_ = declaredModuleId
	this_.declaredModuleId = declaredModuleId
}

// Java method: ModuleMetadataBuilder.recordDeclareLegacyNamespace(com.google.javascript.rhino.Node declaresLegacyNamespace)
func (this_ *GatherModuleMetadata_ModuleMetadataBuilder) RecordDeclareLegacyNamespace(declaresLegacyNamespace pkg_rhino.Node) {
	_ = declaresLegacyNamespace
	this_.declaresLegacyNamespace = declaresLegacyNamespace
}

// Java method: ModuleMetadataBuilder.isScript()
func (this_ *GatherModuleMetadata_ModuleMetadataBuilder) IsScript() bool {
	return (this_.metadataBuilder.ModuleType_0() == StaticGet_ModuleMetadataMap_ModuleType_SCRIPT())
}

// Java method: ModuleMetadataBuilder.build()
func (this_ *GatherModuleMetadata_ModuleMetadataBuilder) Build() ModuleMetadataMap_ModuleMetadata {
	this_.metadataBuilder.GoogNamespacesBuilder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.googNamespaces))
	if !(this_.ambiguous) {
		if this_.hasModuleBody && (this_.metadataBuilder.ModuleType_0() == StaticGet_ModuleMetadataMap_ModuleType_SCRIPT()) {
			this_.metadataBuilder.ModuleType_1(StaticGet_ModuleMetadataMap_ModuleType_ES6_MODULE())
		}
		if (this_.declaredModuleId != 0) && (this_.metadataBuilder.ModuleType_0() != StaticGet_ModuleMetadataMap_ModuleType_ES6_MODULE()) {
			this_.outer_.compiler.Report_1(JSError_Make_3(this_.declaredModuleId, StaticGet_GatherModuleMetadata_DECLARE_MODULE_ID_OUTSIDE_ES6_MODULE(), []string{}))
		}
		if this_.declaresLegacyNamespace != 0 {
			if this_.metadataBuilder.ModuleType_0() == StaticGet_ModuleMetadataMap_ModuleType_GOOG_MODULE() {
				this_.metadataBuilder.ModuleType_1(StaticGet_ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE())
			} else {
				this_.outer_.compiler.Report_1(JSError_Make_3(this_.declaresLegacyNamespace, StaticGet_ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE(), []string{}))
			}
		}
	}
	return this_.metadataBuilder.Build()
}

// Java class: com.google.javascript.jscomp.GatherModuleMetadata.Finder
type GatherModuleMetadata_Finder struct {
	outer_            *GatherModuleMetadata
	toggleModuleNames j2gort.OrderedSet[string]     `gson:"toggleModuleNames"`
	toggleModules     j2gort.CustomOrderedSet[*Var] `gson:"toggleModules"`
	idHash_           atomic.Int32
}

func (this_ *GatherModuleMetadata_Finder) BindSelf_(s any) {
	_ = s
}
func (this_ *GatherModuleMetadata_Finder) IsJavaNil_() bool { return this_ == nil }
func (*GatherModuleMetadata_Finder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GatherModuleMetadata$Finder", "Finder"
}
func (this_ *GatherModuleMetadata_Finder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GatherModuleMetadata_Finder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GatherModuleMetadata.Finder", "com.google.javascript.jscomp.GatherModuleMetadata$Finder", func(v_ any) bool {
		return j2gort.InstanceOf[*GatherModuleMetadata_Finder](v_)
	})
}
func (this_ *GatherModuleMetadata_Finder) SetOuter_(o *GatherModuleMetadata)   { this_.outer_ = o }
func (this_ *GatherModuleMetadata_Finder) GetOuter_() *GatherModuleMetadata    { return this_.outer_ }
func (*GatherModuleMetadata_Finder) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *GatherModuleMetadata_Finder) initFields_GatherModuleMetadata_Finder() {
	_ = this_
	this_.toggleModuleNames = j2gort.NewOrderedSet[string]()
	this_.toggleModules = j2gort.NewCustomOrderedSet[*Var]()
}

// Java constructor: Finder()
func NewGatherModuleMetadata_Finder(outer_ *GatherModuleMetadata) *GatherModuleMetadata_Finder {
	this_ := &GatherModuleMetadata_Finder{}
	this_.outer_ = outer_
	InitGatherModuleMetadata_Finder(this_)
	return this_
}

func InitGatherModuleMetadata_Finder(this_ *GatherModuleMetadata_Finder) {
	this_.initFields_GatherModuleMetadata_Finder()
}

// Java method: Finder.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *GatherModuleMetadata_Finder) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		this_.priv_Finder_enterModule_3018(t, n, t.GetInput().GetPath())
	case int32(116) /* Token.IMPORT */, int32(120) /* Token.EXPORT */ :
		this_.priv_Finder_visitImportOrExport_3016(t, n)
	case int32(28) /* Token.CALL */ :
		if n.IsCall() && StaticGet_GatherModuleMetadata_GOOG_LOADMODULE().Matches(n.GetFirstChild()) {
			this_.outer_.loadModuleCall = n
			this_.priv_Finder_enterModule_3018(t, n, nil)
		}
	case int32(123) /* Token.MODULE_BODY */ :
		{
			r_2 := this_.outer_.currentModule
			r_2.hasModuleBody = true
		}
	case int32(124) /* Token.DYNAMIC_IMPORT */ :
		this_.priv_Finder_visitDynamicImport_3017(n)
	default:
	}
	return true
}

// Java method: Finder.visitImportOrExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node importOrExport)
func (this_ *GatherModuleMetadata_Finder) priv_Finder_visitImportOrExport_3016(t *NodeTraversal, importOrExport pkg_rhino.Node) {
	_ = t
	_ = importOrExport
	j2gort.CheckNotNullPtr(this_.outer_.currentModule)
	this_.outer_.currentModule.ModuleType(StaticGet_ModuleMetadataMap_ModuleType_ES6_MODULE(), t, importOrExport)
	if importOrExport.IsImport() || (importOrExport.HasTwoChildren() && importOrExport.GetLastChild().IsStringLit()) {
		this_.outer_.currentModule.metadataBuilder.Es6ImportSpecifiersBuilder().Add_1_any(j2gort.BoxString(importOrExport.GetLastChild().GetString()))
	}
}

// Java method: Finder.visitDynamicImport(com.google.javascript.rhino.Node dynamicImport)
func (this_ *GatherModuleMetadata_Finder) priv_Finder_visitDynamicImport_3017(dynamicImport pkg_rhino.Node) {
	_ = dynamicImport
	if dynamicImport.GetFirstChild().IsStringLit() {
		this_.outer_.currentModule.metadataBuilder.Es6ImportSpecifiersBuilder().Add_1_any(j2gort.BoxString(dynamicImport.GetFirstChild().GetString()))
	}
}

// Java method: Finder.enterModule(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath path)
func (this_ *GatherModuleMetadata_Finder) priv_Finder_enterModule_3018(t *NodeTraversal, n pkg_rhino.Node, path *ModuleLoader_ModulePath) {
	_ = t
	_ = n
	_ = path
	var val_3 GatherModuleMetadata_ModuleMetadataBuilder
	val_3.outer_ = this_.outer_
	InitGatherModuleMetadata_ModuleMetadataBuilder(&val_3, n, path)
	var newModule *GatherModuleMetadata_ModuleMetadataBuilder = &val_3
	_ = newModule
	if this_.outer_.currentModule != nil {
		if this_.outer_.parentModule != nil {
			t.Report_3(n, StaticGet_GatherModuleMetadata_INVALID_NESTED_LOAD_MODULE(), []string{})
		}
		this_.outer_.parentModule = this_.outer_.currentModule
	}
	this_.outer_.currentModule = newModule
}

// Java method: Finder.leaveModule()
func (this_ *GatherModuleMetadata_Finder) priv_Finder_leaveModule_3019() {
	j2gort.CheckNotNullPtr(this_.outer_.currentModule)
	var module ModuleMetadataMap_ModuleMetadata = this_.outer_.currentModule.Build()
	_ = module
	if module.Path() != nil {
		j2gort.MapPutStr(this_.outer_.modulesByPath, module.Path().String(), module)
	}
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(module.GoogNamespaces()))
		for cur_4.HasNext() {
			namespace := cur_4.NextStr()
			_ = namespace
			j2gort.MapPutStr(this_.outer_.modulesByGoogNamespace, namespace, module)
		}
	}
	if this_.outer_.parentModule != nil {
		this_.outer_.parentModule.metadataBuilder.NestedModulesBuilder().AddVoid_(module)
	}
	this_.outer_.currentModule = this_.outer_.parentModule
	this_.outer_.parentModule = nil
}

// Java method: Finder.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *GatherModuleMetadata_Finder) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if (this_.outer_.processCommonJsModules && (this_.outer_.currentModule != nil)) && this_.outer_.currentModule.IsScript() {
		if ProcessCommonJSModules_IsCommonJsExport(t, n, this_.outer_.moduleResolutionMode) {
			this_.outer_.currentModule.ModuleType(StaticGet_ModuleMetadataMap_ModuleType_COMMON_JS(), t, n)
			return
		}
	}
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		this_.priv_Finder_leaveModule_3019()
	case int32(32) /* Token.NAME */ :
		this_.priv_Finder_visitName_3021(t, n)
	case int32(28) /* Token.CALL */ :
		if this_.outer_.loadModuleCall == n {
			this_.priv_Finder_leaveModule_3019()
			this_.outer_.loadModuleCall = pkg_rhino.Node(0)
		} else {
			this_.priv_Finder_visitGoogCall_3022(t, n)
		}
	default:
	}
}

// Java method: Finder.isFromGoogImport(com.google.javascript.jscomp.Var goog)
func (this_ *GatherModuleMetadata_Finder) priv_Finder_isFromGoogImport_3020(goog *Var) bool {
	_ = goog
	var nameNode pkg_rhino.Node = goog.GetNameNode()
	_ = nameNode
	return (((((nameNode != 0) && nameNode.IsImportStar()) && (j2gort.CheckNotNullString(nameNode.GetString()) == "goog")) && nameNode.GetParent().GetFirstChild().IsEmpty()) && javalang.StringEndsWith(j2gort.CheckNotNullString(nameNode.GetParent().GetLastChild().GetString()), "/goog.js"))
}

// Java method: Finder.visitName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *GatherModuleMetadata_Finder) priv_Finder_visitName_3021(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var name string = n.GetString()
	_ = name
	if this_.toggleModuleNames.Contains(name) {
		var nameVar *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(name))
		_ = nameVar
		if this_.toggleModules.Contains(nameVar) {
			var parent pkg_rhino.Node = n.GetParent()
			_ = parent
			if parent.IsGetProp() {
				this_.priv_Finder_addToggle_3023(t, n, parent.GetString())
			} else if !(NodeUtil_IsNameDeclaration(parent)) {
				t.Report_3(n, StaticGet_GatherModuleMetadata_INVALID_TOGGLE_USAGE(), []string{"toggle modules may not be used other than looking up properties"})
			}
		}
	}
	if !("goog" == name) {
		return
	}
	var root *Var = j2gort.CastPtr[Var](t.GetScope().GetVar("goog"))
	_ = root
	if (root != nil) && !(this_.priv_Finder_isFromGoogImport_3020(root)) {
		return
	}
	this_.outer_.currentModule.metadataBuilder.UsesClosure(true)
}

// Java method: Finder.visitGoogCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *GatherModuleMetadata_Finder) priv_Finder_visitGoogCall_3022(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if (!(n.HasChildren()) || !(n.GetFirstChild().IsGetProp())) || !(n.GetFirstChild().IsQualifiedName()) {
		return
	}
	var getprop pkg_rhino.Node = n.GetFirstChild()
	_ = getprop
	var firstProp pkg_rhino.Node = n.GetFirstChild()
	_ = firstProp
	for firstProp.IsGetProp() {
		firstProp = firstProp.GetFirstChild()
	}
	if !(firstProp.IsName()) || !(j2gort.CheckNotNullString(firstProp.GetString()) == "goog") {
		return
	}
	var root *Var = j2gort.CastPtr[Var](t.GetScope().GetVar("goog"))
	_ = root
	if ((root != nil) && root.IsLocal()) && !(j2gort.CastPtr[Scope](root.GetScope()).IsModuleScope()) {
		return
	}
	if ((root != nil) && j2gort.CastPtr[Scope](root.GetScope()).IsModuleScope()) && !(this_.priv_Finder_isFromGoogImport_3020(root)) {
		return
	}
	if (root == nil) || (NodeUtil_GetEnclosingScript(root.GetNameNode()) != NodeUtil_GetEnclosingScript(n)) {
		this_.outer_.currentModule.metadataBuilder.UsesClosure(true)
	}
	if getprop.MatchesQualifiedName_1_Node(StaticGet_GatherModuleMetadata_GOOG_PROVIDE()) {
		this_.outer_.currentModule.ModuleType(StaticGet_ModuleMetadataMap_ModuleType_GOOG_PROVIDE(), t, n)
		if n.HasTwoChildren() && n.GetLastChild().IsStringLit() {
			var namespace string = n.GetLastChild().GetString()
			_ = namespace
			this_.priv_Finder_addNamespace_3024(this_.outer_.currentModule, StaticGet_ModuleMetadataMap_ModuleType_GOOG_PROVIDE(), namespace, t, n)
		} else {
			t.Report_3(n, StaticGet_ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE(), []string{})
			this_.outer_.currentModule.metadataBuilder.UsesClosure(false)
		}
	} else if getprop.MatchesQualifiedName_1_Node(StaticGet_GatherModuleMetadata_GOOG_MODULE()) {
		this_.outer_.currentModule.ModuleType(StaticGet_ModuleMetadataMap_ModuleType_GOOG_MODULE(), t, n)
		if n.HasTwoChildren() && n.GetLastChild().IsStringLit() {
			var namespace_2 string = n.GetLastChild().GetString()
			_ = namespace_2
			this_.priv_Finder_addNamespace_3024(this_.outer_.currentModule, StaticGet_ModuleMetadataMap_ModuleType_GOOG_MODULE(), namespace_2, t, n)
		} else {
			t.Report_3(n, StaticGet_ClosureRewriteModule_INVALID_MODULE_ID_ARG(), []string{})
			this_.outer_.currentModule.metadataBuilder.UsesClosure(false)
		}
	} else if getprop.MatchesQualifiedName_1_Node(StaticGet_GatherModuleMetadata_GOOG_MODULE_GET()) {
		if (!(n.HasTwoChildren()) || !(n.GetLastChild().IsStringLit())) || !(javalang.StringEndsWith(j2gort.CheckNotNullString(n.GetLastChild().GetString()), "$2etoggles")) {
			return
		}
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		if !(parent.IsGetProp()) {
			t.Report_3(n, StaticGet_GatherModuleMetadata_INVALID_TOGGLE_USAGE(), []string{"goog.module.get of toggles module must immediately look up a single toggle"})
			return
		}
		this_.priv_Finder_addToggle_3023(t, parent, parent.GetString())
	} else if getprop.MatchesQualifiedName_1_Node(StaticGet_GatherModuleMetadata_GOOG_MODULE_DECLARELEGACYNAMESPACE()) {
		this_.outer_.currentModule.RecordDeclareLegacyNamespace(n)
	} else if getprop.MatchesQualifiedName_1_Node(StaticGet_GatherModuleMetadata_GOOG_DECLARE_MODULE_ID()) || getprop.MatchesQualifiedName_1_Node(StaticGet_GatherModuleMetadata_GOOG_MODULE_DECLARNAMESPACE()) {
		if this_.outer_.currentModule.declaredModuleId != 0 {
			t.Report_3(n, StaticGet_GatherModuleMetadata_MULTIPLE_DECLARE_MODULE_NAMESPACE(), []string{})
		}
		if n.HasTwoChildren() && n.GetLastChild().IsStringLit() {
			this_.outer_.currentModule.RecordDeclareModuleId(n)
			var namespace_3 string = n.GetLastChild().GetString()
			_ = namespace_3
			this_.priv_Finder_addNamespace_3024(this_.outer_.currentModule, StaticGet_ModuleMetadataMap_ModuleType_GOOG_MODULE(), namespace_3, t, n)
		} else {
			t.Report_3(n, StaticGet_GatherModuleMetadata_INVALID_DECLARE_MODULE_ID_CALL(), []string{})
		}
	} else if getprop.MatchesQualifiedName_1_Node(StaticGet_GatherModuleMetadata_GOOG_REQUIRE()) {
		if n.HasTwoChildren() && n.GetLastChild().IsStringLit() {
			var namespace_4 string = n.GetLastChild().GetString()
			_ = namespace_4
			this_.outer_.currentModule.metadataBuilder.StronglyRequiredGoogNamespacesBuilder().Add_1_any(j2gort.BoxString(namespace_4))
			if javalang.StringEndsWith(j2gort.CheckNotNullString(namespace_4), "$2etoggles") {
				var callParent pkg_rhino.Node = n.GetParent()
				_ = callParent
				var lhs pkg_rhino.Node = callParent.GetFirstChild()
				_ = lhs
				if callParent.IsDestructuringLhs() {
					{
						cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(lhs.Children()))
						for cur_5.HasNext() {
							key := pkg_rhino.Node(cur_5.NextNode())
							_ = key
							if key.IsStringKey() {
								this_.priv_Finder_addToggle_3023(t, n, key.GetString())
							} else {
								t.Report_3(n, StaticGet_GatherModuleMetadata_INVALID_TOGGLE_USAGE(), []string{"must be destructured with string keys"})
							}
						}
					}
				} else if callParent.IsName() {
					var name string = callParent.GetString()
					_ = name
					var nameVar *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(name))
					_ = nameVar
					this_.toggleModules.Add(nameVar)
					this_.toggleModuleNames.Add(name)
				} else if this_.outer_.currentModule.metadataBuilder.ModuleType_0() != StaticGet_ModuleMetadataMap_ModuleType_GOOG_PROVIDE() {
					t.Report_3(n, StaticGet_GatherModuleMetadata_INVALID_TOGGLE_USAGE(), []string{"import must be assigned"})
				}
			}
		} else {
			t.Report_3(n, StaticGet_ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE(), []string{})
		}
	} else if getprop.MatchesQualifiedName_1_Node(StaticGet_GatherModuleMetadata_GOOG_REQUIRE_TYPE()) {
		if n.HasTwoChildren() && n.GetLastChild().IsStringLit() {
			this_.outer_.currentModule.metadataBuilder.WeaklyRequiredGoogNamespacesBuilder().Add_1_any(j2gort.BoxString(n.GetLastChild().GetString()))
		} else {
			t.Report_3(n, StaticGet_GatherModuleMetadata_INVALID_REQUIRE_TYPE(), []string{})
		}
	} else if getprop.MatchesQualifiedName_1_Node(StaticGet_GatherModuleMetadata_GOOG_MAYBE_REQUIRE()) {
		if n.HasTwoChildren() && n.GetLastChild().IsStringLit() {
			this_.outer_.currentModule.metadataBuilder.MaybeRequiredGoogNamespacesBuilder().Add_1_any(j2gort.BoxString(n.GetLastChild().GetString()))
		} else {
			t.Report_3(n, StaticGet_GatherModuleMetadata_INVALID_MAYBE_REQUIRE(), []string{})
		}
	} else if getprop.MatchesQualifiedName_1_Node(StaticGet_GatherModuleMetadata_GOOG_SET_TEST_ONLY()) {
		if n.HasOneChild() || (n.HasTwoChildren() && n.GetLastChild().IsStringLit()) {
			this_.outer_.currentModule.metadataBuilder.IsTestOnly(true)
		} else {
			t.Report_3(n, StaticGet_GatherModuleMetadata_INVALID_SET_TEST_ONLY(), []string{})
		}
	} else if getprop.MatchesQualifiedName_1_Node(StaticGet_GatherModuleMetadata_GOOG_REQUIRE_DYNAMIC()) {
		if n.HasTwoChildren() && n.GetLastChild().IsStringLit() {
			this_.outer_.currentModule.metadataBuilder.DynamicallyRequiredGoogNamespacesBuilder().Add_1_any(j2gort.BoxString(n.GetLastChild().GetString()))
		} else {
			t.Report_3(n, StaticGet_GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC(), []string{})
		}
	}
}

// Java method: Finder.addToggle(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, java.lang.String name)
func (this_ *GatherModuleMetadata_Finder) priv_Finder_addToggle_3023(t *NodeTraversal, n pkg_rhino.Node, name string) {
	_ = t
	_ = n
	_ = name
	if javalang.StringStartsWith(j2gort.CheckNotNullString(name), "TOGGLE_") {
		var toggleName string = j2gort.JStrSubFrom(name, j2gort.JStrLen("TOGGLE_"))
		_ = toggleName
		this_.outer_.currentModule.metadataBuilder.ReadTogglesBuilder().Add_1_any(j2gort.BoxString(toggleName))
	} else {
		t.Report_3(n, StaticGet_GatherModuleMetadata_INVALID_TOGGLE_USAGE(), []string{"all toggle names must start with `TOGGLE_`"})
	}
}

// Java method: Finder.addNamespace(com.google.javascript.jscomp.GatherModuleMetadata.ModuleMetadataBuilder module, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleType moduleType, java.lang.String namespace, com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *GatherModuleMetadata_Finder) priv_Finder_addNamespace_3024(module *GatherModuleMetadata_ModuleMetadataBuilder, moduleType ModuleMetadataMap_ModuleType, namespace string, t *NodeTraversal, n pkg_rhino.Node) {
	_ = module
	_ = moduleType
	_ = namespace
	_ = t
	_ = n
	{
		var iifeCond_1 bool
		{

			re_6 := moduleType
			ae_7 := StaticGet_ModuleMetadataMap_ModuleType_GOOG_PROVIDE()
			if re_6 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_6 == ae_7
		}
		if iifeCond_1 || func() bool {
			re_8 := moduleType
			ae_9 := StaticGet_ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE()
			if re_8 == 0 {
				j2gort.ThrowNPE()
			}
			return re_8 == ae_9
		}() {
			if !(NodeUtil_IsValidQualifiedName_2_FeatureSet_string(this_.outer_.compiler.GetOptions().GetLanguageIn().ToFeatureSet(), namespace)) {
				this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID(), []string{namespace}))
			}
		}
	}
	{
		var iifeCond_1 bool
		{

			re_10 := moduleType
			ae_11 := StaticGet_ModuleMetadataMap_ModuleType_GOOG_MODULE()
			if re_10 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_10 == ae_11
		}
		if iifeCond_1 || func() bool {
			re_12 := moduleType
			ae_13 := StaticGet_ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE()
			if re_12 == 0 {
				j2gort.ThrowNPE()
			}
			return re_12 == ae_13
		}() {
			if !(GatherModuleMetadata_IsValidModuleId(namespace)) {
				this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID(), []string{namespace}))
			}
		}
	}
	var existingType ModuleMetadataMap_ModuleType = ModuleMetadataMap_ModuleType(0)
	_ = existingType
	var existingFileSource string = j2gort.NullString
	_ = existingFileSource
	if j2gort.MultisetContainsStr(module.googNamespaces, namespace) {
		existingType = module.metadataBuilder.ModuleType_0()
		existingFileSource = t.GetSourceName()
	} else {
		var existingModule ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.outer_.modulesByGoogNamespace, namespace))
		_ = existingModule
		if existingModule != nil && !existingModule.IsJavaNil_() {
			existingType = existingModule.ModuleType()
			existingFileSource = existingModule.RootNode().GetSourceFileName()
		}
	}
	j2gort.MultisetAddStr(this_.outer_.currentModule.googNamespaces, namespace)
	if existingType != 0 {
		switch (existingType).Ordinal() {
		case int32(0) /* ModuleType.ES6_MODULE */, int32(1) /* ModuleType.GOOG_MODULE */, int32(2) /* ModuleType.LEGACY_GOOG_MODULE */ :
			var diagnostic *DiagnosticType
			_ = diagnostic
			{
				var iifeCond_1 bool
				{

					re_14 := moduleType
					ae_15 := StaticGet_ModuleMetadataMap_ModuleType_GOOG_PROVIDE()
					if re_14 == 0 {
						j2gort.ThrowNPE()
					}
					iifeCond_1 = re_14 == ae_15
				}
				if iifeCond_1 {
					diagnostic = StaticGet_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE()
				} else {
					diagnostic = StaticGet_ClosurePrimitiveErrors_DUPLICATE_MODULE()
				}
			}
			t.Report_3(n, diagnostic, []string{namespace, existingFileSource})
			return
		case int32(4) /* ModuleType.GOOG_PROVIDE */ :
			var diagnostic_2 *DiagnosticType
			_ = diagnostic_2
			{
				var iifeCond_1 bool
				{

					re_16 := moduleType
					ae_17 := StaticGet_ModuleMetadataMap_ModuleType_GOOG_PROVIDE()
					if re_16 == 0 {
						j2gort.ThrowNPE()
					}
					iifeCond_1 = re_16 == ae_17
				}
				if iifeCond_1 {
					diagnostic_2 = StaticGet_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE()
				} else {
					diagnostic_2 = StaticGet_ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE()
				}
			}
			t.Report_3(n, diagnostic_2, []string{namespace, existingFileSource})
			return
		case int32(3) /* ModuleType.COMMON_JS */, int32(5) /* ModuleType.SCRIPT */ :
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected module type: "+j2gort.StrVal(existingType))))
	}
}
