/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ProcessDefines.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

// Java class: com.google.javascript.jscomp.ProcessDefines
type ProcessDefines struct {
	compiler                            AbstractCompiler                                  `gson:"compiler"`
	registry                            *pkg_rhino.JSTypeRegistry                         `gson:"registry"`
	replacementValuesFromFlags          j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"replacementValuesFromFlags"`
	namespaceSupplier                   j2gort.Ext_java_util_function_Supplier            `gson:"namespaceSupplier"`
	enableZonesDefineName               string                                            `gson:"enableZonesDefineName"`
	zoneInputPattern                    *j2gort.Ext_java_util_regex_Pattern               `gson:"zoneInputPattern"`
	unknownDefinesToIgnore              j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"unknownDefinesToIgnore"`
	knownDefineJsdocs                   j2gort.PlainSet[*pkg_rhino.JSDocInfo]             `gson:"knownDefineJsdocs"`
	knownGoogDefineCalls                j2gort.PlainSet[pkg_rhino.Node]                   `gson:"knownGoogDefineCalls"`
	defineByDefineName                  j2gort.OrderedMap[string, *ProcessDefines_Define] `gson:"defineByDefineName"`
	replacementValuesFromClosureDefines j2gort.OrderedMap[string, pkg_rhino.Node]         `gson:"replacementValuesFromClosureDefines"`
	validDefineValueExpressions         j2gort.PlainSet[pkg_rhino.Node]                   `gson:"validDefineValueExpressions"`
	namespace                           *GlobalNamespace                                  `gson:"namespace"`
	idHash_                             atomic.Int32
	mode                                ProcessDefines_Mode `gson:"mode"`
	recognizeClosureDefines             bool                `gson:"recognizeClosureDefines"`
	hasZoneInput                        bool                `gson:"hasZoneInput"`
}

func (this_ *ProcessDefines) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessDefines) IsJavaNil_() bool { return this_ == nil }
func (*ProcessDefines) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessDefines", "ProcessDefines"
}
func (this_ *ProcessDefines) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessDefines) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessDefines", "com.google.javascript.jscomp.ProcessDefines", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessDefines](v_)
	})
}
func (*ProcessDefines) Implements_jscomp_CompilerPass() {}

var static_ProcessDefines_KNOWN_DEFINES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ProcessDefines_KNOWN_DEFINES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ProcessDefines()
	return static_ProcessDefines_KNOWN_DEFINES
}
func StaticSet_ProcessDefines_KNOWN_DEFINES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ProcessDefines()
	static_ProcessDefines_KNOWN_DEFINES = v
}

var static_ProcessDefines_GOOG_DEFINE pkg_rhino.Node

func StaticGet_ProcessDefines_GOOG_DEFINE() pkg_rhino.Node {
	clinit_ProcessDefines()
	return static_ProcessDefines_GOOG_DEFINE
}
func StaticSet_ProcessDefines_GOOG_DEFINE(v pkg_rhino.Node) {
	clinit_ProcessDefines()
	static_ProcessDefines_GOOG_DEFINE = v
}

var static_ProcessDefines_UNKNOWN_DEFINE_WARNING *DiagnosticType

func StaticGet_ProcessDefines_UNKNOWN_DEFINE_WARNING() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_UNKNOWN_DEFINE_WARNING
}
func StaticSet_ProcessDefines_UNKNOWN_DEFINE_WARNING(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_UNKNOWN_DEFINE_WARNING = v
}

var static_ProcessDefines_INVALID_DEFINE_NAME_ERROR *DiagnosticType

func StaticGet_ProcessDefines_INVALID_DEFINE_NAME_ERROR() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_INVALID_DEFINE_NAME_ERROR
}
func StaticSet_ProcessDefines_INVALID_DEFINE_NAME_ERROR(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_INVALID_DEFINE_NAME_ERROR = v
}

var static_ProcessDefines_MISSING_DEFINE_ANNOTATION *DiagnosticType

func StaticGet_ProcessDefines_MISSING_DEFINE_ANNOTATION() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_MISSING_DEFINE_ANNOTATION
}
func StaticSet_ProcessDefines_MISSING_DEFINE_ANNOTATION(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_MISSING_DEFINE_ANNOTATION = v
}

var static_ProcessDefines_INVALID_DEFINE_TYPE *DiagnosticType

func StaticGet_ProcessDefines_INVALID_DEFINE_TYPE() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_INVALID_DEFINE_TYPE
}
func StaticSet_ProcessDefines_INVALID_DEFINE_TYPE(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_INVALID_DEFINE_TYPE = v
}

var static_ProcessDefines_INVALID_DEFINE_VALUE *DiagnosticType

func StaticGet_ProcessDefines_INVALID_DEFINE_VALUE() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_INVALID_DEFINE_VALUE
}
func StaticSet_ProcessDefines_INVALID_DEFINE_VALUE(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_INVALID_DEFINE_VALUE = v
}

var static_ProcessDefines_INVALID_DEFINE_LOCATION *DiagnosticType

func StaticGet_ProcessDefines_INVALID_DEFINE_LOCATION() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_INVALID_DEFINE_LOCATION
}
func StaticSet_ProcessDefines_INVALID_DEFINE_LOCATION(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_INVALID_DEFINE_LOCATION = v
}

var static_ProcessDefines_NON_CONST_DEFINE *DiagnosticType

func StaticGet_ProcessDefines_NON_CONST_DEFINE() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_NON_CONST_DEFINE
}
func StaticSet_ProcessDefines_NON_CONST_DEFINE(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_NON_CONST_DEFINE = v
}

var static_ProcessDefines_CLOSURE_DEFINES_ERROR *DiagnosticType

func StaticGet_ProcessDefines_CLOSURE_DEFINES_ERROR() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_CLOSURE_DEFINES_ERROR
}
func StaticSet_ProcessDefines_CLOSURE_DEFINES_ERROR(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_CLOSURE_DEFINES_ERROR = v
}

var static_ProcessDefines_CLOSURE_DEFINES_MULTIPLE *DiagnosticType

func StaticGet_ProcessDefines_CLOSURE_DEFINES_MULTIPLE() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_CLOSURE_DEFINES_MULTIPLE
}
func StaticSet_ProcessDefines_CLOSURE_DEFINES_MULTIPLE(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_CLOSURE_DEFINES_MULTIPLE = v
}

var static_ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR *DiagnosticType

func StaticGet_ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR
}
func StaticSet_ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR = v
}

var static_ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT *DiagnosticType

func StaticGet_ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT
}
func StaticSet_ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT = v
}

var static_ProcessDefines_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT *DiagnosticType

func StaticGet_ProcessDefines_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT
}
func StaticSet_ProcessDefines_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT = v
}

var static_ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE *DiagnosticType

func StaticGet_ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE() *DiagnosticType {
	clinit_ProcessDefines()
	return static_ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE
}
func StaticSet_ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE(v *DiagnosticType) {
	clinit_ProcessDefines()
	static_ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE = v
}

var clinitDone_ProcessDefines uint32
var clinitState_ProcessDefines j2gort.ClinitState

func clinit_ProcessDefines() {
	if atomic.LoadUint32(&clinitDone_ProcessDefines) != 0 {
		return
	}
	clinitSlow_ProcessDefines()
}

func clinitSlow_ProcessDefines() {
	if !j2gort.EnterClinit(&clinitState_ProcessDefines, &clinitDone_ProcessDefines) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ProcessDefines, &clinitDone_ProcessDefines)
	static_ProcessDefines_KNOWN_DEFINES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("COMPILED", "goog.DEBUG", "$jscomp.ASSUME_ES5", "$jscomp.ASSUME_ES6", "$jscomp.ASSUME_ES2020", "$jscomp.ISOLATE_POLYFILLS", []any{"$jscomp.INSTRUMENT_ASYNC_CONTEXT"})
	static_ProcessDefines_GOOG_DEFINE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "define")
	static_ProcessDefines_UNKNOWN_DEFINE_WARNING = DiagnosticType_Warning("JSC_UNKNOWN_DEFINE_WARNING", "unknown @define variable {0}")
	static_ProcessDefines_INVALID_DEFINE_NAME_ERROR = DiagnosticType_Error("JSC_INVALID_DEFINE_NAME_ERROR", "\"{0}\" is not a valid JS identifier name")
	static_ProcessDefines_MISSING_DEFINE_ANNOTATION = DiagnosticType_Error("JSC_INVALID_MISSING_DEFINE_ANNOTATION", "Missing @define annotation")
	static_ProcessDefines_INVALID_DEFINE_TYPE = DiagnosticType_Error("JSC_INVALID_DEFINE_TYPE", "@define tag only permits primitive types")
	static_ProcessDefines_INVALID_DEFINE_VALUE = DiagnosticType_Error("JSC_INVALID_DEFINE_VALUE", "invalid initialization value for @define {0}")
	static_ProcessDefines_INVALID_DEFINE_LOCATION = DiagnosticType_Error("JSC_INVALID_DEFINE_LOCATION", "@define must be initialized on a static qualified name in global or module scope")
	static_ProcessDefines_NON_CONST_DEFINE = DiagnosticType_Error("JSC_NON_CONST_DEFINE", "@define {0} has already been set at {1}.")
	static_ProcessDefines_CLOSURE_DEFINES_ERROR = DiagnosticType_Error("JSC_CLOSURE_DEFINES_ERROR", "Invalid CLOSURE_DEFINES definition")
	static_ProcessDefines_CLOSURE_DEFINES_MULTIPLE = DiagnosticType_Error("JSC_CLOSURE_DEFINES_MULTIPLE", "Multiple CLOSURE_DEFINES definitions for {0}. First occurrence: {1}")
	static_ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR = DiagnosticType_Error("JSC_NON_GLOBAL_CLOSURE_DEFINES_ERROR", "CLOSURE_DEFINES definition must be in top-level global scope")
	static_ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT = DiagnosticType_Error("JSC_DEFINE_CALL_WITHOUT_ASSIGNMENT", "The result of a goog.define call must be assigned as an isolated statement.")
	static_ProcessDefines_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT = DiagnosticType_Error("JSC_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT", "ZoneJS is incompatible with language level ES2017 or higher (See go/ngissue/31730)\nPlease set `--language_out=ECMASCRIPT_2016` (or older) in your flags.")
	static_ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE = DiagnosticType_Disabled("JSC_DEFINE_WITHOUT_GOOG_DEFINE", "The const declaration for @define ''{0}'' must be initialized with a call to goog.define().")
}

func (this_ *ProcessDefines) initFields_ProcessDefines() {
	_ = this_
	this_.enableZonesDefineName = j2gort.NullString
	this_.knownDefineJsdocs = j2gort.NewPlainSetVal[*pkg_rhino.JSDocInfo]()
	this_.knownGoogDefineCalls = j2gort.NewPlainSetVal[pkg_rhino.Node]()
	this_.defineByDefineName = j2gort.NewOrderedMap[string, *ProcessDefines_Define]()
	this_.replacementValuesFromClosureDefines = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	this_.validDefineValueExpressions = j2gort.NewPlainSetVal[pkg_rhino.Node]()

}

// Java constructor: ProcessDefines(com.google.javascript.jscomp.ProcessDefines.Builder builder)
func NewProcessDefines(builder *ProcessDefines_Builder) *ProcessDefines {
	clinit_ProcessDefines()
	this_ := &ProcessDefines{}
	InitProcessDefines(this_, builder)
	return this_
}

func InitProcessDefines(this_ *ProcessDefines, builder *ProcessDefines_Builder) {
	_ = builder
	this_.initFields_ProcessDefines()
	this_.mode = builder.mode
	this_.compiler = builder.compiler
	if (this_.mode).FieldGet_check() {
		this_.registry = this_.compiler.GetTypeRegistry()
	} else {
		this_.registry = nil
	}
	this_.replacementValuesFromFlags = j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(builder.replacementValues)
	this_.namespaceSupplier = builder.namespaceSupplier
	this_.recognizeClosureDefines = builder.recognizeClosureDefines
	this_.enableZonesDefineName = builder.enableZonesDefineName
	this_.zoneInputPattern = builder.zoneInputPattern
	this_.unknownDefinesToIgnore = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](builder.unknownDefinesToIgnore))
}

// Java method: ProcessDefines.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ProcessDefines) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	ownsNamespace_ := this_.namespaceSupplier == nil && this_.namespace == nil
	defer func() {
		if ownsNamespace_ && this_.namespace != nil {
			this_.defineByDefineName.Clear()
			this_.namespace.ReleaseSlabs_()
			this_.namespace = nil
		}
	}()

	this_.priv_ProcessDefines_initNamespace_4538(externs, root)
	this_.priv_ProcessDefines_collectDefines_4543(root)
	this_.priv_ProcessDefines_reportInvalidDefineLocations_4551(root)
	this_.priv_ProcessDefines_collectValidDefineValueExpressions_4549()
	this_.priv_ProcessDefines_validateDefineDeclarations_4550()
	this_.priv_ProcessDefines_overrideDefines_4539()
	this_.priv_ProcessDefines_validateDefines_4555()
}

// Java method: ProcessDefines.initNamespace(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ProcessDefines) priv_ProcessDefines_initNamespace_4538(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if this_.namespaceSupplier != nil {
		this_.namespace = j2gort.CastPtr[GlobalNamespace](this_.namespaceSupplier.Get())
	}
	if this_.namespace == nil {
		this_.namespace = NewGlobalNamespace_3_AbstractCompiler_Node_Node(this_.compiler, externs, root)
	}
}

// Java method: ProcessDefines.overrideDefines()
func (this_ *ProcessDefines) priv_ProcessDefines_overrideDefines_4539() {
	if (this_.mode).FieldGet_optimize() {
		{
			cur_1 := this_.defineByDefineName.IterValues()
			for cur_1.HasNext() {

				define := cur_1.
					Next()
				_ = define
				if define.valueParent == 0 {
					continue
				}
				var inputValue pkg_rhino.Node = this_.priv_ProcessDefines_getReplacementForDefine_4540(define)
				_ = inputValue
				if (inputValue == 0) || (inputValue == define.value) {
					continue
				}
				var changed bool = (((define.value == 0) || (inputValue.RawToken_() != define.value.RawToken_())) || !(inputValue.IsEquivalentTo_1(define.value)))
				_ = changed
				if changed {
					if define.value == 0 {
						define.valueParent.AddChildToBack(inputValue.CloneTree_0())
					} else {
						define.value.ReplaceWith(inputValue.CloneTree_0())
					}
					this_.compiler.ReportChangeToEnclosingScope(define.valueParent)
				}
				if J2clSourceFileChecker_ShouldRunJ2clPasses(this_.compiler) {
					var alias string = ProcessDefines_GetGlobalDefineAlias(define.defineName)
					_ = alias
					{
						var iifeCond_1 bool
						{

							sr_3 := alias
							se_4 := define.defineName
							iifeCond_1 = j2gort.CheckNotNullString(sr_3) == se_4
						}
						if !(iifeCond_1) {
							var globalDefineName pkg_rhino.Node = NodeUtil_NewName_3_AbstractCompiler_string_Node(this_.compiler, ProcessDefines_GetGlobalDefineAlias(define.defineName), define.valueParent)
							_ = globalDefineName
							var defineLhs pkg_rhino.Node
							_ = defineLhs
							if define.valueParent.IsAssign() {
								defineLhs = define.valueParent.GetFirstChild().CloneTree_0()
							} else {
								defineLhs = define.valueParent.CloneNode_0()
							}
							if !(defineLhs.IsName() || defineLhs.IsQualifiedName()) {
								j2gort.CheckStateMsg(false, j2gort.BoxNode(defineLhs))
							}
							var globalDefineVar pkg_rhino.Node = pkg_rhino.IR_Var_2(globalDefineName, defineLhs).SrcrefTreeIfMissing(define.valueParent)
							_ = globalDefineVar
							globalDefineVar.InsertAfter(define.valueParent.GetParent())
							this_.compiler.ReportChangeToEnclosingScope(define.valueParent.GetParent())
						}
					}
				}
			}
		}
	}
	if (this_.mode).FieldGet_optimize() {
		var unusedReplacements j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_Sets_Difference(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_Sets_Union(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.replacementValuesFromFlags.KeySet_0()), j2gort.AsIface[j2gort.Ext_java_util_Set](this_.replacementValuesFromClosureDefines.LiveKeySet()))), j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_Sets_Union(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_Sets_Union(j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_ProcessDefines_KNOWN_DEFINES()), j2gort.AsIface[j2gort.Ext_java_util_Set](this_.defineByDefineName.LiveKeySet()))), j2gort.AsIface[j2gort.Ext_java_util_Set](this_.unknownDefinesToIgnore)))))
		_ = unusedReplacements
		{
			cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(unusedReplacements))
			for cur_5.HasNext() {
				unknownDefine := cur_5.NextStr()
				_ = unknownDefine
				this_.compiler.Report_1(JSError_Make_2(StaticGet_ProcessDefines_UNKNOWN_DEFINE_WARNING(), []string{unknownDefine}))
			}
		}
	}
}

// Java static method: ProcessDefines.getGlobalDefineAlias(java.lang.String defineName)
func ProcessDefines_GetGlobalDefineAlias(defineName string) string {
	clinit_ProcessDefines()
	_ = defineName
	if j2gort.SetContainsStr(StaticGet_ProcessDefines_KNOWN_DEFINES(), defineName) || (j2gort.CheckNotNullString(defineName) == "goog.LOCALE") {
		return defineName
	}
	return nodestore.JoinDollar("jscomp$defines", j2gort.JStrReplaceChar(defineName, uint16(46), uint16(36)))
}

// Java method: ProcessDefines.getReplacementForDefine(com.google.javascript.jscomp.ProcessDefines.Define define)
func (this_ *ProcessDefines) priv_ProcessDefines_getReplacementForDefine_4540(define *ProcessDefines_Define) pkg_rhino.Node {
	_ = define
	var replacement pkg_rhino.Node = this_.priv_ProcessDefines_getReplacementForDefineName_4541(define.defineName)
	_ = replacement
	if replacement != 0 {
		return replacement
	}
	if this_.priv_ProcessDefines_isGoogDefineCall_4557(define.value) && (define.value.GetChildCount() == int32(3)) {
		return define.value.GetChildAtIndex(int32(2))
	}
	return pkg_rhino.Node(0)
}

// Java method: ProcessDefines.getReplacementForDefineName(java.lang.String defineName)
func (this_ *ProcessDefines) priv_ProcessDefines_getReplacementForDefineName_4541(defineName string) pkg_rhino.Node {
	_ = defineName
	var replacementFromFlags pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](j2gort.MapGetStr(this_.replacementValuesFromFlags, defineName))
	_ = replacementFromFlags
	if replacementFromFlags != 0 {
		return replacementFromFlags
	}
	var replacementFromClosureDefines pkg_rhino.Node = this_.replacementValuesFromClosureDefines.Get(defineName)
	_ = replacementFromClosureDefines
	if replacementFromClosureDefines != 0 {
		return replacementFromClosureDefines
	}
	return pkg_rhino.Node(0)
}

// Java method: ProcessDefines.isValidDefineType(com.google.javascript.rhino.JSTypeExpression expression)
func (this_ *ProcessDefines) priv_ProcessDefines_isValidDefineType_4542(expression pkg_rhino.JSTypeExpression) bool {
	_ = expression
	var type_ pkg_rhino.JSType = this_.registry.EvaluateTypeExpressionInGlobalScope(expression)
	_ = type_
	base_type__ := type_.GetJSTypeBase()
	return (!(base_type__.IsUnknownType()) && base_type__.IsSubtypeOf_1(this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NUMBER_STRING_BOOLEAN())))
}

// Java method: ProcessDefines.collectDefines(com.google.javascript.rhino.Node root)
func (this_ *ProcessDefines) priv_ProcessDefines_collectDefines_4543(root pkg_rhino.Node) {
	_ = root
	if this_.recognizeClosureDefines {
		{
			clinit_NodeTraversal()
			var ntb_ NodeTraversal_Builder
			InitNodeTraversal_Builder(&ntb_)
			(&ntb_).SetCompiler(this_.compiler)
			(&ntb_).SetCallback_1_NodeTraversal_Callback(NewProcessDefines_ClosureDefinesCollector(this_))
			(&ntb_).Traverse(root)
		}
	}
	{
		cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.namespace.GetAllSymbols())))
		for cur_6.HasNext() {
			rawElem_7 := cur_6.Next()
			name := j2gort.CastPtr[GlobalNamespace_Name](rawElem_7)
			_ = name
			var declaration *GlobalNamespace_Ref = this_.priv_ProcessDefines_selectDefineDeclaration_4546(name)
			_ = declaration
			if declaration == nil {
				continue
			}
			var totalSets int32 = name.GetTotalSets()
			_ = totalSets
			var valueParent pkg_rhino.Node = ProcessDefines_GetValueParentForDefine(declaration)
			_ = valueParent
			var value pkg_rhino.Node
			_ = value
			if valueParent != 0 {
				value = valueParent.GetLastChild()
			} else {
				value = pkg_rhino.Node(0)
			}
			var defineName string = j2gort.NullString
			_ = defineName
			if this_.priv_ProcessDefines_isGoogDefineCall_4557(value) && this_.priv_ProcessDefines_verifyGoogDefine_4558(value) {
				var nameNode pkg_rhino.Node = value.GetSecondChild()
				_ = nameNode
				defineName = nameNode.GetString()
			} else {
				defineName = name.GetFullName()
			}
			var existingDefine *ProcessDefines_Define = this_.defineByDefineName.PutIfAbsent(defineName, NewProcessDefines_Define(defineName, name, declaration, valueParent, value))
			_ = existingDefine
			if existingDefine != nil {
				declaration = existingDefine.declaration
				totalSets = (totalSets + existingDefine.name.GetTotalSets())
			}
			if totalSets > int32(1) {
				if name.HasInlineRefs() {
					this_.priv_ProcessDefines_reportIfExtraSetRef_4544(name.InlineRef0(), declaration, defineName)
					this_.priv_ProcessDefines_reportIfExtraSetRef_4544(name.InlineRef1(), declaration, defineName)
					this_.priv_ProcessDefines_reportIfExtraSetRef_4544(name.InlineRef2(), declaration, defineName)
					this_.priv_ProcessDefines_reportIfExtraSetRef_4544(name.InlineRef3(), declaration, defineName)
				} else {
					{
						cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(name.GetRefs()))
						for cur_8.HasNext() {
							rawElem_9 := cur_8.Next()
							ref := j2gort.CastPtr[GlobalNamespace_Ref](rawElem_9)
							_ = ref
							this_.priv_ProcessDefines_reportIfExtraSetRef_4544(ref, declaration, defineName)
						}
					}
				}
			}
		}
	}
	this_.compiler.SetDefineNames(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](StaticGet_ProcessDefines_KNOWN_DEFINES()))).AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.defineByDefineName.LiveKeySet()))).Build()))
}

// Java method: ProcessDefines.reportIfExtraSetRef(com.google.javascript.jscomp.GlobalNamespace.@org.jspecify.annotations.Nullable Ref ref, com.google.javascript.jscomp.GlobalNamespace.Ref declaration, java.lang.String defineName)
func (this_ *ProcessDefines) priv_ProcessDefines_reportIfExtraSetRef_4544(ref *GlobalNamespace_Ref, declaration *GlobalNamespace_Ref, defineName string) {
	_ = ref
	_ = declaration
	_ = defineName
	if ((ref != nil) && ref.IsSet()) && !(j2gort.CheckNotNullPtr(ref) == declaration) {
		this_.compiler.Report_1(JSError_Make_3(ref.GetNode(), StaticGet_ProcessDefines_NON_CONST_DEFINE(), []string{defineName, declaration.GetNode().GetLocation()}))
	}
}

// Java method: ProcessDefines.matchDefineDeclarationRef(com.google.javascript.jscomp.GlobalNamespace.@org.jspecify.annotations.Nullable Ref ref)
func (this_ *ProcessDefines) priv_ProcessDefines_matchDefineDeclarationRef_4545(ref *GlobalNamespace_Ref) *GlobalNamespace_Ref {
	_ = ref
	if (ref == nil) || !(ref.IsSetFromGlobal()) {
		return nil
	}
	var refNode pkg_rhino.Node = ref.GetNode()
	_ = refNode
	if !(refNode.IsQualifiedName()) {
		return nil
	}
	var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(refNode)
	_ = jsdoc
	if (jsdoc == nil) || !(jsdoc.IsDefine()) {
		return nil
	}
	this_.knownDefineJsdocs.AddVoid(jsdoc)
	return ref
}

// Java method: ProcessDefines.selectDefineDeclaration(com.google.javascript.jscomp.GlobalNamespace.Name name)
func (this_ *ProcessDefines) priv_ProcessDefines_selectDefineDeclaration_4546(name *GlobalNamespace_Name) *GlobalNamespace_Ref {
	_ = name
	if name.GetGlobalSets() == int32(0) {
		return nil
	}
	if name.HasInlineRefs() {
		var r *GlobalNamespace_Ref = this_.priv_ProcessDefines_matchDefineDeclarationRef_4545(name.InlineRef0())
		_ = r
		if r != nil {
			return r
		}
		r = this_.priv_ProcessDefines_matchDefineDeclarationRef_4545(name.InlineRef1())
		if r != nil {
			return r
		}
		r = this_.priv_ProcessDefines_matchDefineDeclarationRef_4545(name.InlineRef2())
		if r != nil {
			return r
		}
		return this_.priv_ProcessDefines_matchDefineDeclarationRef_4545(name.InlineRef3())
	}
	{
		cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(name.GetRefs()))
		for cur_10.HasNext() {
			rawElem_11 := cur_10.Next()
			ref := j2gort.CastPtr[GlobalNamespace_Ref](rawElem_11)
			_ = ref
			var r_2 *GlobalNamespace_Ref = this_.priv_ProcessDefines_matchDefineDeclarationRef_4545(ref)
			_ = r_2
			if r_2 != nil {
				return r_2
			}
		}
	}
	return nil
}

// Java static method: ProcessDefines.getValueParentForDefine(com.google.javascript.jscomp.GlobalNamespace.Ref declaration)
func ProcessDefines_GetValueParentForDefine(declaration *GlobalNamespace_Ref) pkg_rhino.Node {
	clinit_ProcessDefines()
	_ = declaration
	var declarationNode pkg_rhino.Node = declaration.GetNode()
	_ = declarationNode
	var declarationParent pkg_rhino.Node = declarationNode.GetParent()
	_ = declarationParent
	if declarationParent.IsVar() || declarationParent.IsConst() {
		if !(declarationNode.IsName()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(declarationNode))
		}
		return declarationNode
	} else if declarationParent.IsAssign() && declarationNode.IsFirstChildOf(declarationParent) {
		return declarationParent
	}
	return pkg_rhino.Node(0)
}

// Java method: ProcessDefines.addNonSetRefNodes(com.google.javascript.jscomp.GlobalNamespace.Name name)
func (this_ *ProcessDefines) priv_ProcessDefines_addNonSetRefNodes_4547(name *GlobalNamespace_Name) {
	_ = name
	if name.HasInlineRefs() {
		var r0 *GlobalNamespace_Ref = name.InlineRef0()
		_ = r0
		if r0 != nil {
			if !(r0.IsSet()) {
				this_.validDefineValueExpressions.AddVoid(r0.GetNode())
			}
			var r1 *GlobalNamespace_Ref = name.InlineRef1()
			_ = r1
			if r1 != nil {
				if !(r1.IsSet()) {
					this_.validDefineValueExpressions.AddVoid(r1.GetNode())
				}
				var r2 *GlobalNamespace_Ref = name.InlineRef2()
				_ = r2
				if r2 != nil {
					if !(r2.IsSet()) {
						this_.validDefineValueExpressions.AddVoid(r2.GetNode())
					}
					var r3 *GlobalNamespace_Ref = name.InlineRef3()
					_ = r3
					if (r3 != nil) && !(r3.IsSet()) {
						this_.validDefineValueExpressions.AddVoid(r3.GetNode())
					}
				}
			}
		}
		return
	}
	{
		cur_12 := j2gort.NewIterCursor(j2gort.CheckNotNull(name.GetRefs()))
		for cur_12.HasNext() {
			rawElem_13 := cur_12.Next()
			ref := j2gort.CastPtr[GlobalNamespace_Ref](rawElem_13)
			_ = ref
			if !(ref.IsSet()) {
				this_.validDefineValueExpressions.AddVoid(ref.GetNode())
			}
		}
	}
}

// Java method: ProcessDefines.addAllRefNodes(com.google.javascript.jscomp.GlobalNamespace.Name name)
func (this_ *ProcessDefines) priv_ProcessDefines_addAllRefNodes_4548(name *GlobalNamespace_Name) {
	_ = name
	if name.HasInlineRefs() {
		var r0 *GlobalNamespace_Ref = name.InlineRef0()
		_ = r0
		if r0 != nil {
			this_.validDefineValueExpressions.AddVoid(r0.GetNode())
			var r1 *GlobalNamespace_Ref = name.InlineRef1()
			_ = r1
			if r1 != nil {
				this_.validDefineValueExpressions.AddVoid(r1.GetNode())
				var r2 *GlobalNamespace_Ref = name.InlineRef2()
				_ = r2
				if r2 != nil {
					this_.validDefineValueExpressions.AddVoid(r2.GetNode())
					var r3 *GlobalNamespace_Ref = name.InlineRef3()
					_ = r3
					if r3 != nil {
						this_.validDefineValueExpressions.AddVoid(r3.GetNode())
					}
				}
			}
		}
		return
	}
	{
		cur_14 := j2gort.NewIterCursor(j2gort.CheckNotNull(name.GetRefs()))
		for cur_14.HasNext() {
			rawElem_15 := cur_14.Next()
			ref := j2gort.CastPtr[GlobalNamespace_Ref](rawElem_15)
			_ = ref
			this_.validDefineValueExpressions.AddVoid(ref.GetNode())
		}
	}
}

// Java method: ProcessDefines.collectValidDefineValueExpressions()
func (this_ *ProcessDefines) priv_ProcessDefines_collectValidDefineValueExpressions_4549() {
	var defineNames *j2gort.PlainSet[*GlobalNamespace_Name] = nil
	_ = defineNames
	if !(this_.defineByDefineName.IsEmpty()) {
		defineNames = j2gort.NewPlainSetPtrCap[*GlobalNamespace_Name]((this_.defineByDefineName.Size() * int32(2)))
		{
			cur_16 := this_.defineByDefineName.IterValues()
			for cur_16.HasNext() {

				define := cur_16.
					Next()
				_ = define
				defineNames.AddVoid(define.name)
				this_.priv_ProcessDefines_addNonSetRefNodes_4547(define.name)
			}
		}
	}
	var namesToCheck []*GlobalNamespace_Name = []*GlobalNamespace_Name{}
	_ = namesToCheck
	{
		cur_18 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.namespace.GetAllSymbols())))
		for cur_18.HasNext() {
			rawElem_19 := cur_18.Next()
			name := j2gort.CastPtr[GlobalNamespace_Name](rawElem_19)
			_ = name
			if ProcessDefines_IsGlobalConst(name) && ((defineNames == nil) || !(defineNames.Contains(name))) {
				namesToCheck = append(namesToCheck, name)
			}
		}
	}
	var count int32 = int32(len(namesToCheck))
	_ = count
	for count > int32(0) {
		var size int32 = count
		_ = size
		var write int32 = int32(0)
		_ = write
		{
			var i int32 = int32(0)
			_ = i
			for ; i < size; i = i + 1 {
				var name_2 *GlobalNamespace_Name = namesToCheck[i]
				_ = name_2
				var declValue pkg_rhino.Node = ProcessDefines_GetConstantDeclValue(j2gort.CastPtr[GlobalNamespace_Ref](name_2.GetDeclaration()).GetNode())
				_ = declValue
				switch (this_.priv_ProcessDefines_isValidDefineValue_4554(declValue)).Ordinal() {
				case int32(2) /* Tri.TRUE */ :
					this_.priv_ProcessDefines_addAllRefNodes_4548(name_2)
				case int32(1) /* Tri.UNKNOWN */ :
					namesToCheck[func() int32 { old_20 := write; nw_21 := old_20 + 1; write = nw_21; return old_20 }()] = name_2
				default:
				}
			}
		}
		if write == size {
			break
		}
		count = write
	}
}

// Java method: ProcessDefines.validateDefineDeclarations()
func (this_ *ProcessDefines) priv_ProcessDefines_validateDefineDeclarations_4550() {
	if !((this_.mode).FieldGet_check()) {
		return
	}
	{
		cur_22 := this_.defineByDefineName.IterValues()
		for cur_22.HasNext() {

			define := cur_22.
				Next()
			_ = define
			var declarationNode pkg_rhino.Node = define.declaration.GetNode()
			_ = declarationNode
			if (!(this_.priv_ProcessDefines_isGoogDefineCall_4557(define.value)) && !(j2gort.SetContainsStr(StaticGet_ProcessDefines_KNOWN_DEFINES(), define.defineName))) && !(javalang.StringStartsWith(j2gort.CheckNotNullString(define.defineName), "$jscomp.")) {
				this_.compiler.Report_1(JSError_Make_3(j2gort.CastNode[pkg_rhino.Node](j2gort.Ext_com_google_common_base_MoreObjects_FirstNonNull(j2gort.BoxNode(define.value), j2gort.BoxNode(j2gort.CastNode[pkg_rhino.Node](j2gort.Ext_com_google_common_base_MoreObjects_FirstNonNull(j2gort.BoxNode(define.valueParent), j2gort.BoxNode(declarationNode)))))), StaticGet_ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE(), []string{define.defineName, "go/js-practices/locals#defines"}))
			}
			if !(this_.priv_ProcessDefines_hasValidValue_4553(define)) {
				this_.compiler.Report_1(JSError_Make_3(j2gort.CastNode[pkg_rhino.Node](j2gort.Ext_com_google_common_base_MoreObjects_FirstNonNull(j2gort.BoxNode(define.value), j2gort.BoxNode(j2gort.CastNode[pkg_rhino.Node](j2gort.Ext_com_google_common_base_MoreObjects_FirstNonNull(j2gort.BoxNode(define.valueParent), j2gort.BoxNode(declarationNode)))))), StaticGet_ProcessDefines_INVALID_DEFINE_VALUE(), []string{define.defineName}))
			}
			var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(declarationNode)
			_ = jsdoc
			if (jsdoc == nil) || !(this_.priv_ProcessDefines_isValidDefineType_4542(jsdoc.GetType())) {
				this_.compiler.Report_1(JSError_Make_3(declarationNode, StaticGet_ProcessDefines_INVALID_DEFINE_TYPE(), []string{}))
			}
		}
	}
}

// Java method: ProcessDefines.reportInvalidDefineLocations(com.google.javascript.rhino.Node root)
func (this_ *ProcessDefines) priv_ProcessDefines_reportInvalidDefineLocations_4551(root pkg_rhino.Node) {
	_ = root
	if !((this_.mode).FieldGet_check()) {
		return
	}
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(newReleasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface(func(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
		_ = t
		_ = n
		_ = parent
		var jsdoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
		_ = jsdoc
		if ((jsdoc != nil) && jsdoc.IsDefine()) && this_.knownDefineJsdocs.Add(jsdoc) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ProcessDefines_INVALID_DEFINE_LOCATION(), []string{}))
		}
		if this_.priv_ProcessDefines_isGoogDefineCall_4557(n) && this_.knownGoogDefineCalls.Add(n) {
			this_.priv_ProcessDefines_verifyGoogDefine_4558(n)
		}
		if (n.MatchesName_1_string("CLOSURE_DEFINES") && (NodeUtil_IsNameDeclaration(parent) || (parent.IsGetElem() && parent.GetParent().IsAssign()))) && !(NodeUtil_GetEnclosingScopeRoot(n).IsRoot()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR(), []string{}))
		}
	})).Traverse(root)
}

// Java method: ProcessDefines.handleClosureDefinesValue(com.google.javascript.rhino.Node stringNode, com.google.javascript.rhino.Node valueNode, com.google.javascript.rhino.Node errorNode)
func (this_ *ProcessDefines) priv_ProcessDefines_handleClosureDefinesValue_4552(stringNode pkg_rhino.Node, valueNode pkg_rhino.Node, errorNode pkg_rhino.Node) {
	_ = stringNode
	_ = valueNode
	_ = errorNode
	if (stringNode.IsStringKey() || stringNode.IsStringLit()) && ProcessDefines_IsValidClosureDefinesValue(valueNode) {
		var key string = stringNode.GetString()
		_ = key
		if this_.replacementValuesFromClosureDefines.ContainsKey(key) {
			this_.compiler.Report_1(JSError_Make_3(errorNode, StaticGet_ProcessDefines_CLOSURE_DEFINES_MULTIPLE(), []string{key, this_.replacementValuesFromClosureDefines.Get(key).GetLocation()}))
		}
		this_.replacementValuesFromClosureDefines.Put(key, valueNode)
	} else if (this_.mode).FieldGet_check() {
		this_.compiler.Report_1(JSError_Make_3(errorNode, StaticGet_ProcessDefines_CLOSURE_DEFINES_ERROR(), []string{}))
	}
}

// Java static method: ProcessDefines.isValidClosureDefinesValue(com.google.javascript.rhino.Node val)
func ProcessDefines_IsValidClosureDefinesValue(val pkg_rhino.Node) bool {
	clinit_ProcessDefines()
	_ = val
	switch val.TokenSwitchOrd_() {
	case int32(35) /* Token.STRINGLIT */, int32(33) /* Token.NUMBER */, int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */ :
		return true
	case int32(22) /* Token.NEG */ :
		return val.GetFirstChild().IsNumber()
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: ProcessDefines.hasValidValue(com.google.javascript.jscomp.ProcessDefines.Define define)
func (this_ *ProcessDefines) priv_ProcessDefines_hasValidValue_4553(define *ProcessDefines_Define) bool {
	_ = define
	if define.valueParent == 0 {
		return false
	} else if define.valueParent.IsFromExterns() {
		return true
	} else {
		return this_.priv_ProcessDefines_isValidDefineValue_4554(define.value).ToBoolean(false)
	}
}

// Java static method: ProcessDefines.isGlobalConst(com.google.javascript.jscomp.GlobalNamespace.Name name)
func ProcessDefines_IsGlobalConst(name *GlobalNamespace_Name) bool {
	clinit_ProcessDefines()
	_ = name
	return (((name.GetTotalSets() == int32(1)) && (j2gort.CastPtr[GlobalNamespace_Ref](name.GetDeclaration()) != nil)) && j2gort.CastPtr[GlobalNamespace_Ref](name.GetDeclaration()).IsSetFromGlobal())
}

// Java method: ProcessDefines.isValidDefineValue(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node val)
func (this_ *ProcessDefines) priv_ProcessDefines_isValidDefineValue_4554(val pkg_rhino.Node) pkg_base.Tri {
	_ = val
	if val == 0 {
		return pkg_base.StaticGet_Tri_FALSE()
	}
	switch val.TokenSwitchOrd_() {
	case int32(35) /* Token.STRINGLIT */, int32(33) /* Token.NUMBER */, int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */ :
		return pkg_base.StaticGet_Tri_TRUE()
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */, int32(13) /* Token.ADD */, int32(3) /* Token.BITAND */, int32(20) /* Token.BITNOT */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(16) /* Token.DIV */, int32(4) /* Token.EQ */, int32(18) /* Token.EXPONENT */, int32(9) /* Token.GE */, int32(8) /* Token.GT */, int32(7) /* Token.LE */, int32(10) /* Token.LSH */, int32(6) /* Token.LT */, int32(17) /* Token.MOD */, int32(15) /* Token.MUL */, int32(5) /* Token.NE */, int32(11) /* Token.RSH */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */, int32(14) /* Token.SUB */, int32(12) /* Token.URSH */ :
		return this_.priv_ProcessDefines_isValidDefineValue_4554(val.GetFirstChild()).And(this_.priv_ProcessDefines_isValidDefineValue_4554(val.GetLastChild()))
	case int32(67) /* Token.HOOK */ :
		return this_.priv_ProcessDefines_isValidDefineValue_4554(val.GetFirstChild()).And(this_.priv_ProcessDefines_isValidDefineValue_4554(val.GetSecondChild())).And(this_.priv_ProcessDefines_isValidDefineValue_4554(val.GetLastChild()))
	case int32(19) /* Token.NOT */, int32(22) /* Token.NEG */, int32(21) /* Token.POS */ :
		return this_.priv_ProcessDefines_isValidDefineValue_4554(val.GetFirstChild())
	case int32(32) /* Token.NAME */, int32(26) /* Token.GETPROP */ :
		if val.IsQualifiedName() {
			if this_.validDefineValueExpressions.Contains(val) {
				return pkg_base.StaticGet_Tri_TRUE()
			} else {
				return pkg_base.StaticGet_Tri_UNKNOWN()
			}
		}
	case int32(28) /* Token.CALL */ :
		if !(this_.priv_ProcessDefines_isGoogDefineCall_4557(val)) {
			return pkg_base.StaticGet_Tri_FALSE()
		}
		if !(val.HasXChildren(int32(3))) {
			return pkg_base.StaticGet_Tri_TRUE()
		}
		return this_.priv_ProcessDefines_isValidDefineValue_4554(val.GetChildAtIndex(int32(2)))
	default:
	}
	return pkg_base.StaticGet_Tri_FALSE()
}

// Java static method: ProcessDefines.getConstantDeclValue(com.google.javascript.rhino.Node name)
func ProcessDefines_GetConstantDeclValue(name pkg_rhino.Node) pkg_rhino.Node {
	clinit_ProcessDefines()
	_ = name
	var parent pkg_rhino.Node = name.GetParent()
	_ = parent
	if parent == 0 {
		return pkg_rhino.Node(0)
	}
	if name.IsName() {
		if parent.IsConst() {
			return name.GetFirstChild()
		} else if !(parent.IsVar()) {
			return pkg_rhino.Node(0)
		}
		var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(name)
		_ = jsdoc
		if (jsdoc != nil) && jsdoc.IsConstant() {
			return name.GetFirstChild()
		} else {
			return pkg_rhino.Node(0)
		}
	} else if name.IsGetProp() && parent.IsAssign() {
		var jsdoc_2 *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(name)
		_ = jsdoc_2
		if (jsdoc_2 != nil) && jsdoc_2.IsConstant() {
			return name.GetNext()
		} else {
			return pkg_rhino.Node(0)
		}
	}
	return pkg_rhino.Node(0)
}

// Java method: ProcessDefines.validateDefines()
func (this_ *ProcessDefines) priv_ProcessDefines_validateDefines_4555() {
	var willOutputAsyncFunctions bool = this_.compiler.GetOptions().GetOutputFeatureSet().Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS())
	_ = willOutputAsyncFunctions
	if (this_.priv_ProcessDefines_isZoneEnabled_4556() && this_.hasZoneInput) && willOutputAsyncFunctions {
		this_.compiler.Report_1(JSError_Make_2(StaticGet_ProcessDefines_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT(), []string{}))
	}
}

// Java method: ProcessDefines.isZoneEnabled()
func (this_ *ProcessDefines) priv_ProcessDefines_isZoneEnabled_4556() bool {
	var enableZonesDefault bool = true
	_ = enableZonesDefault
	var zoneEnabled pkg_rhino.Node = this_.priv_ProcessDefines_getReplacementForDefineName_4541(this_.enableZonesDefineName)
	_ = zoneEnabled
	if zoneEnabled == 0 {
		return enableZonesDefault
	}
	return NodeUtil_GetBooleanValue(zoneEnabled).ToBoolean(enableZonesDefault)
}

// Java method: ProcessDefines.isGoogDefineCall(com.google.javascript.rhino.Node node)
func (this_ *ProcessDefines) priv_ProcessDefines_isGoogDefineCall_4557(node pkg_rhino.Node) bool {
	_ = node
	if !(this_.recognizeClosureDefines) {
		return false
	}
	if (node == 0) || !(node.IsCall()) {
		return false
	}
	return node.GetFirstChild().MatchesQualifiedName_1_Node(StaticGet_ProcessDefines_GOOG_DEFINE())
}

// Java method: ProcessDefines.verifyGoogDefine(com.google.javascript.rhino.Node callNode)
func (this_ *ProcessDefines) priv_ProcessDefines_verifyGoogDefine_4558(callNode pkg_rhino.Node) bool {
	_ = callNode
	this_.knownGoogDefineCalls.AddVoid(callNode)
	var parent pkg_rhino.Node = callNode.GetParent()
	_ = parent
	var methodName pkg_rhino.Node = callNode.GetFirstChild()
	_ = methodName
	var args pkg_rhino.Node = callNode.GetSecondChild()
	_ = args
	if NodeUtil_GetEnclosingFunction(callNode) != 0 {
		this_.compiler.Report_1(JSError_Make_3(methodName.GetParent(), StaticGet_ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR(), []string{}))
		return false
	}
	if parent.IsAssign() && parent.GetParent().IsExprResult() {
		parent = parent.GetParent()
	} else if parent.IsName() && NodeUtil_IsNameDeclaration(parent.GetParent()) {
		parent = parent.GetParent()
	} else {
		this_.compiler.Report_1(JSError_Make_3(methodName.GetParent(), StaticGet_ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT(), []string{}))
		return false
	}
	var arg pkg_rhino.Node = args
	_ = arg
	if !(this_.priv_ProcessDefines_verifyNotNull_4559(methodName, arg)) || !(this_.priv_ProcessDefines_verifyOfType_4561(methodName, arg, pkg_rhino.StaticGet_Token_STRINGLIT())) {
		return false
	}
	arg = arg.GetNext()
	if !(args.IsFromExterns()) && (!(this_.priv_ProcessDefines_verifyNotNull_4559(methodName, arg)) || !(this_.priv_ProcessDefines_verifyIsLast_4560(methodName, arg))) {
		return false
	}
	var name string = args.GetString()
	_ = name
	if !(NodeUtil_IsValidQualifiedName_2_FeatureSet_string(this_.compiler.GetOptions().GetLanguageIn().ToFeatureSet(), name)) {
		this_.compiler.Report_1(JSError_Make_3(args, StaticGet_ProcessDefines_INVALID_DEFINE_NAME_ERROR(), []string{name}))
		return false
	}
	var info *pkg_rhino.JSDocInfo = (func() pkg_rhino.Node {
		if parent.IsExprResult() {
			return parent.GetFirstChild()
		}
		return parent
	}()).GetJSDocInfo()
	_ = info
	if (info == nil) || !(info.IsDefine()) {
		this_.compiler.Report_1(JSError_Make_3(parent, StaticGet_ProcessDefines_MISSING_DEFINE_ANNOTATION(), []string{}))
		return false
	}
	return true
}

// Java method: ProcessDefines.verifyNotNull(com.google.javascript.rhino.Node methodName, com.google.javascript.rhino.Node arg)
func (this_ *ProcessDefines) priv_ProcessDefines_verifyNotNull_4559(methodName pkg_rhino.Node, arg pkg_rhino.Node) bool {
	_ = methodName
	_ = arg
	if arg == 0 {
		this_.compiler.Report_1(JSError_Make_3(methodName, StaticGet_ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR(), []string{methodName.GetQualifiedName()}))
		return false
	}
	return true
}

// Java method: ProcessDefines.verifyIsLast(com.google.javascript.rhino.Node methodName, com.google.javascript.rhino.Node arg)
func (this_ *ProcessDefines) priv_ProcessDefines_verifyIsLast_4560(methodName pkg_rhino.Node, arg pkg_rhino.Node) bool {
	_ = methodName
	_ = arg
	if arg.GetNext() != 0 {
		this_.compiler.Report_1(JSError_Make_3(methodName, StaticGet_ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR(), []string{methodName.GetQualifiedName()}))
		return false
	}
	return true
}

// Java method: ProcessDefines.verifyOfType(com.google.javascript.rhino.Node methodName, com.google.javascript.rhino.Node arg, com.google.javascript.rhino.Token desiredType)
func (this_ *ProcessDefines) priv_ProcessDefines_verifyOfType_4561(methodName pkg_rhino.Node, arg pkg_rhino.Node, desiredType pkg_rhino.Token) bool {
	_ = methodName
	_ = arg
	_ = desiredType
	if arg.GetToken() != desiredType {
		this_.compiler.Report_1(JSError_Make_3(methodName, StaticGet_ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR(), []string{methodName.GetQualifiedName()}))
		return false
	}
	return true
}

// Java enum: com.google.javascript.jscomp.ProcessDefines.Mode
type ProcessDefines_Mode uint8

const (
	enum_ProcessDefines_Mode_CHECK              ProcessDefines_Mode = 1
	enum_ProcessDefines_Mode_OPTIMIZE           ProcessDefines_Mode = 2
	enum_ProcessDefines_Mode_CHECK_AND_OPTIMIZE ProcessDefines_Mode = 3
)

var enumBaseTable_ProcessDefines_Mode = [...]j2gort.EnumBase{
	{Name_: "CHECK", Ordinal_: 0},
	{Name_: "OPTIMIZE", Ordinal_: 1},
	{Name_: "CHECK_AND_OPTIMIZE", Ordinal_: 2},
}
var enumValues_ProcessDefines_Mode = []ProcessDefines_Mode{1, 2, 3}

var enumField_ProcessDefines_Mode_check = [...]bool{
	true,
	false,
	true,
}

func (this_ ProcessDefines_Mode) FieldGet_check() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_ProcessDefines_Mode_check[this_-1]
}

var enumField_ProcessDefines_Mode_optimize = [...]bool{
	false,
	true,
	true,
}

func (this_ ProcessDefines_Mode) FieldGet_optimize() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_ProcessDefines_Mode_optimize[this_-1]
}

func StaticGet_ProcessDefines_Mode_CHECK() ProcessDefines_Mode {
	return enum_ProcessDefines_Mode_CHECK
}
func StaticGet_ProcessDefines_Mode_OPTIMIZE() ProcessDefines_Mode {
	return enum_ProcessDefines_Mode_OPTIMIZE
}
func StaticGet_ProcessDefines_Mode_CHECK_AND_OPTIMIZE() ProcessDefines_Mode {
	return enum_ProcessDefines_Mode_CHECK_AND_OPTIMIZE
}

func ProcessDefines_Mode_Values() []ProcessDefines_Mode {
	return j2gort.CloneSlice(enumValues_ProcessDefines_Mode)
}

func ProcessDefines_Mode_ValueOf(name string) ProcessDefines_Mode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ProcessDefines_Mode {
		if enumBaseTable_ProcessDefines_Mode[i_].Name_ == name {
			return ProcessDefines_Mode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ProcessDefines.Mode." + name))
}

func (this_ ProcessDefines_Mode) IsJavaNil_() bool                              { return this_ == 0 }
func (ProcessDefines_Mode) HasIdentityEquals_()                                 {}
func (ProcessDefines_Mode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ProcessDefines_Mode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ProcessDefines_Mode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ProcessDefines_Mode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessDefines$Mode", "Mode"
}
func (this_ ProcessDefines_Mode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ProcessDefines_Mode[this_-1]
}
func (this_ ProcessDefines_Mode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ProcessDefines_Mode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ProcessDefines_Mode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ProcessDefines_Mode[this_-1].Name_
}
func (this_ ProcessDefines_Mode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ProcessDefines_Mode[this_-1].Name_
}
func (this_ ProcessDefines_Mode) Equals(other_ any) bool {
	o_, ok_ := other_.(ProcessDefines_Mode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ProcessDefines_Mode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ProcessDefines_Mode[this_-1].HashCode()
}
func (this_ ProcessDefines_Mode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ProcessDefines_Mode[this_-1].CompareTo(other_)
}
func (this_ ProcessDefines_Mode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ProcessDefines_Mode[this_-1].Name_)
}
func init() {
	vals_ := []any{ProcessDefines_Mode(1), ProcessDefines_Mode(2), ProcessDefines_Mode(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ProcessDefines.Mode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ProcessDefines$Mode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessDefines.Mode", "com.google.javascript.jscomp.ProcessDefines$Mode", func(v_ any) bool {
		return j2gort.InstanceOf[ProcessDefines_Mode](v_)
	})
}

// Java class: com.google.javascript.jscomp.ProcessDefines.Builder
type ProcessDefines_Builder struct {
	compiler                AbstractCompiler                                   `gson:"compiler"`
	replacementValues       j2gort.Ext_java_util_Map                           `gson:"replacementValues"`
	namespaceSupplier       j2gort.Ext_java_util_function_Supplier             `gson:"namespaceSupplier"`
	enableZonesDefineName   string                                             `gson:"enableZonesDefineName"`
	zoneInputPattern        *j2gort.Ext_java_util_regex_Pattern                `gson:"zoneInputPattern"`
	unknownDefinesToIgnore  j2gort.Ext_com_google_common_collect_ImmutableList `gson:"unknownDefinesToIgnore"`
	idHash_                 atomic.Int32
	mode                    ProcessDefines_Mode `gson:"mode"`
	recognizeClosureDefines bool                `gson:"recognizeClosureDefines"`
}

func (this_ *ProcessDefines_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessDefines_Builder) IsJavaNil_() bool { return this_ == nil }
func (*ProcessDefines_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessDefines$Builder", "Builder"
}
func (this_ *ProcessDefines_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessDefines_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessDefines.Builder", "com.google.javascript.jscomp.ProcessDefines$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessDefines_Builder](v_)
	})
}

func (this_ *ProcessDefines_Builder) initFields_ProcessDefines_Builder() {
	_ = this_
	this_.replacementValues = j2gort.NewExt_java_util_LinkedHashMap()
	this_.recognizeClosureDefines = true
	this_.enableZonesDefineName = j2gort.NullString

	this_.unknownDefinesToIgnore = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
}

// Java constructor: Builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewProcessDefines_Builder(compiler AbstractCompiler) *ProcessDefines_Builder {
	this_ := &ProcessDefines_Builder{}
	InitProcessDefines_Builder(this_, compiler)
	return this_
}

func InitProcessDefines_Builder(this_ *ProcessDefines_Builder, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ProcessDefines_Builder()
	this_.compiler = compiler
}

// Java method: Builder.putReplacements(java.util.Map<java.lang.String,com.google.javascript.rhino.Node> replacementValues)
func (this_ *ProcessDefines_Builder) PutReplacements(replacementValues j2gort.Ext_java_util_Map) *ProcessDefines_Builder {
	_ = replacementValues
	this_.replacementValues.PutAll(replacementValues)
	return this_
}

// Java method: Builder.setMode(com.google.javascript.jscomp.ProcessDefines.Mode x)
func (this_ *ProcessDefines_Builder) SetMode(x ProcessDefines_Mode) *ProcessDefines_Builder {
	_ = x
	this_.mode = x
	return this_
}

// Java method: Builder.injectNamespace(java.util.function.Supplier<com.google.javascript.jscomp.GlobalNamespace> namespaceSupplier)
func (this_ *ProcessDefines_Builder) InjectNamespace(namespaceSupplier j2gort.Ext_java_util_function_Supplier) *ProcessDefines_Builder {
	_ = namespaceSupplier
	this_.namespaceSupplier = namespaceSupplier
	return this_
}

// Java method: Builder.setUnknownDefinesToIgnore(com.google.common.collect.ImmutableList<java.lang.String> unknownDefinesToIgnore)
func (this_ *ProcessDefines_Builder) SetUnknownDefinesToIgnore(unknownDefinesToIgnore j2gort.Ext_com_google_common_collect_ImmutableList) *ProcessDefines_Builder {
	_ = unknownDefinesToIgnore
	this_.unknownDefinesToIgnore = unknownDefinesToIgnore
	return this_
}

// Java method: Builder.setRecognizeClosureDefines(boolean recognizeClosureDefines)
func (this_ *ProcessDefines_Builder) SetRecognizeClosureDefines(recognizeClosureDefines bool) *ProcessDefines_Builder {
	_ = recognizeClosureDefines
	this_.recognizeClosureDefines = recognizeClosureDefines
	return this_
}

// Java method: Builder.setEnableZonesDefineName(java.lang.@org.jspecify.annotations.Nullable String enableZonesDefineName)
func (this_ *ProcessDefines_Builder) SetEnableZonesDefineName(enableZonesDefineName string) *ProcessDefines_Builder {
	_ = enableZonesDefineName
	this_.enableZonesDefineName = enableZonesDefineName
	return this_
}

// Java method: Builder.setZoneInputPattern(java.util.regex.@org.jspecify.annotations.Nullable Pattern zoneInputPattern)
func (this_ *ProcessDefines_Builder) SetZoneInputPattern(zoneInputPattern *j2gort.Ext_java_util_regex_Pattern) *ProcessDefines_Builder {
	_ = zoneInputPattern
	this_.zoneInputPattern = zoneInputPattern
	return this_
}

// Java method: Builder.build()
func (this_ *ProcessDefines_Builder) Build() *ProcessDefines {
	return NewProcessDefines(this_)
}

// Java class: com.google.javascript.jscomp.ProcessDefines.ClosureDefinesCollector
type ProcessDefines_ClosureDefinesCollector struct {
	outer_  *ProcessDefines
	idHash_ atomic.Int32
}

func (this_ *ProcessDefines_ClosureDefinesCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessDefines_ClosureDefinesCollector) IsJavaNil_() bool { return this_ == nil }
func (*ProcessDefines_ClosureDefinesCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessDefines$ClosureDefinesCollector", "ClosureDefinesCollector"
}
func (this_ *ProcessDefines_ClosureDefinesCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessDefines_ClosureDefinesCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessDefines.ClosureDefinesCollector", "com.google.javascript.jscomp.ProcessDefines$ClosureDefinesCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessDefines_ClosureDefinesCollector](v_)
	})
}
func (this_ *ProcessDefines_ClosureDefinesCollector) SetOuter_(o *ProcessDefines)         { this_.outer_ = o }
func (this_ *ProcessDefines_ClosureDefinesCollector) GetOuter_() *ProcessDefines          { return this_.outer_ }
func (*ProcessDefines_ClosureDefinesCollector) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ProcessDefines_ClosureDefinesCollector) initFields_ProcessDefines_ClosureDefinesCollector() {
	_ = this_
}

// Java constructor: ClosureDefinesCollector()
func NewProcessDefines_ClosureDefinesCollector(outer_ *ProcessDefines) *ProcessDefines_ClosureDefinesCollector {
	this_ := &ProcessDefines_ClosureDefinesCollector{}
	this_.outer_ = outer_
	InitProcessDefines_ClosureDefinesCollector(this_)
	return this_
}

func InitProcessDefines_ClosureDefinesCollector(this_ *ProcessDefines_ClosureDefinesCollector) {
	this_.initFields_ProcessDefines_ClosureDefinesCollector()
}

// Java method: ClosureDefinesCollector.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessDefines_ClosureDefinesCollector) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return ((((n.IsRoot() || n.IsScript()) || n.IsExprResult()) || n.IsAssign()) || NodeUtil_IsNameDeclaration(n))
}

// Java method: ClosureDefinesCollector.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessDefines_ClosureDefinesCollector) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() && (this_.outer_.zoneInputPattern != nil) {
		var source string = n.GetSourceFileName()
		_ = source
		if (source != j2gort.NullString) && this_.outer_.zoneInputPattern.Matcher(j2gort.BoxString(source)).Matches() {
			this_.outer_.hasZoneInput = true
		}
	}
	if NodeUtil_IsNameDeclaration(n) && n.GetFirstChild().MatchesName_1_string("CLOSURE_DEFINES") {
		var valueNode pkg_rhino.Node = n.GetFirstFirstChild()
		_ = valueNode
		if (valueNode != 0) && valueNode.IsObjectLit() {
			{
				var c pkg_rhino.Node = valueNode.GetFirstChild()
				_ = c
				for ; c != 0; c = c.GetNext() {
					this_.outer_.priv_ProcessDefines_handleClosureDefinesValue_4552(c, c.GetFirstChild(), n)
				}
			}
		}
	} else if n.IsAssign() {
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		if lhs.IsGetElem() && lhs.GetFirstChild().MatchesName_1_string("CLOSURE_DEFINES") {
			this_.outer_.priv_ProcessDefines_handleClosureDefinesValue_4552(lhs.GetSecondChild(), n.GetSecondChild(), n)
		}
	}
}

// Java class: com.google.javascript.jscomp.ProcessDefines.Define
type ProcessDefines_Define struct {
	defineName  string                `gson:"defineName"`
	name        *GlobalNamespace_Name `gson:"name"`
	declaration *GlobalNamespace_Ref  `gson:"declaration"`
	valueParent pkg_rhino.Node        `gson:"valueParent"`
	value       pkg_rhino.Node        `gson:"value"`
	idHash_     atomic.Int32
}

func (this_ *ProcessDefines_Define) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessDefines_Define) IsJavaNil_() bool { return this_ == nil }
func (*ProcessDefines_Define) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessDefines$Define", "Define"
}
func (this_ *ProcessDefines_Define) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessDefines_Define) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessDefines.Define", "com.google.javascript.jscomp.ProcessDefines$Define", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessDefines_Define](v_)
	})
}

func (this_ *ProcessDefines_Define) initFields_ProcessDefines_Define() {
	_ = this_
	this_.defineName = j2gort.NullString
}

// Java constructor: Define(java.lang.String defineName, com.google.javascript.jscomp.GlobalNamespace.Name name, com.google.javascript.jscomp.GlobalNamespace.Ref declaration, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node valueParent, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node value)
func NewProcessDefines_Define(defineName string, name *GlobalNamespace_Name, declaration *GlobalNamespace_Ref, valueParent pkg_rhino.Node, value pkg_rhino.Node) *ProcessDefines_Define {
	this_ := &ProcessDefines_Define{}
	InitProcessDefines_Define(this_, defineName, name, declaration, valueParent, value)
	return this_
}

func InitProcessDefines_Define(this_ *ProcessDefines_Define, defineName string, name *GlobalNamespace_Name, declaration *GlobalNamespace_Ref, valueParent pkg_rhino.Node, value pkg_rhino.Node) {
	_ = defineName
	_ = name
	_ = declaration
	_ = valueParent
	_ = value
	this_.initFields_ProcessDefines_Define()
	if !(((valueParent == 0) || (value == 0)) || (value.GetParent() == valueParent)) {
		j2gort.CheckState(false)
	}
	if !(declaration.IsSet()) {
		j2gort.CheckState(false)
	}
	this_.defineName = defineName
	this_.name = name
	this_.declaration = declaration
	this_.valueParent = valueParent
	this_.value = value
}
func (*ProcessDefines_ClosureDefinesCollector) CanReleaseSyntacticScopes_() {}
