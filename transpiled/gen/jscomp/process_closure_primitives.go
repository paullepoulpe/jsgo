/*
 * Copyright 2006 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ProcessClosurePrimitives.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ProcessClosurePrimitives
type ProcessClosurePrimitives struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler               AbstractCompiler                                  `gson:"compiler"`
	knownClosureSubclasses j2gort.OrderedSet[string]                         `gson:"knownClosureSubclasses"`
	closureModules         j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"closureModules"`
	googProvideHashes      []int32                                           `gson:"googProvideHashes"`
}

func (this_ *ProcessClosurePrimitives) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessClosurePrimitives) IsJavaNil_() bool { return this_ == nil }
func (*ProcessClosurePrimitives) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessClosurePrimitives", "ProcessClosurePrimitives"
}
func (this_ *ProcessClosurePrimitives) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ProcessClosurePrimitives) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessClosurePrimitives", "com.google.javascript.jscomp.ProcessClosurePrimitives", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessClosurePrimitives](v_)
	})
}
func (*ProcessClosurePrimitives) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ProcessClosurePrimitives) Implements_jscomp_CompilerPass()                            {}
func (*ProcessClosurePrimitives) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR *DiagnosticType

func StaticGet_ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR
}
func StaticSet_ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR = v
}

var static_ProcessClosurePrimitives_EXPECTED_STRING_ERROR *DiagnosticType

func StaticGet_ProcessClosurePrimitives_EXPECTED_STRING_ERROR() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_EXPECTED_STRING_ERROR
}
func StaticSet_ProcessClosurePrimitives_EXPECTED_STRING_ERROR(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_EXPECTED_STRING_ERROR = v
}

var static_ProcessClosurePrimitives_INVALID_STYLE_ERROR *DiagnosticType

func StaticGet_ProcessClosurePrimitives_INVALID_STYLE_ERROR() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_INVALID_STYLE_ERROR
}
func StaticSet_ProcessClosurePrimitives_INVALID_STYLE_ERROR(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_INVALID_STYLE_ERROR = v
}

var static_ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE *DiagnosticType

func StaticGet_ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE
}
func StaticSet_ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE = v
}

var static_ProcessClosurePrimitives_CLASS_NAMESPACE_ERROR *DiagnosticType

func StaticGet_ProcessClosurePrimitives_CLASS_NAMESPACE_ERROR() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_CLASS_NAMESPACE_ERROR
}
func StaticSet_ProcessClosurePrimitives_CLASS_NAMESPACE_ERROR(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_CLASS_NAMESPACE_ERROR = v
}

var static_ProcessClosurePrimitives_FUNCTION_NAMESPACE_ERROR *DiagnosticType

func StaticGet_ProcessClosurePrimitives_FUNCTION_NAMESPACE_ERROR() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_FUNCTION_NAMESPACE_ERROR
}
func StaticSet_ProcessClosurePrimitives_FUNCTION_NAMESPACE_ERROR(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_FUNCTION_NAMESPACE_ERROR = v
}

var static_ProcessClosurePrimitives_INVALID_PROVIDE_ERROR *DiagnosticType

func StaticGet_ProcessClosurePrimitives_INVALID_PROVIDE_ERROR() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_INVALID_PROVIDE_ERROR
}
func StaticSet_ProcessClosurePrimitives_INVALID_PROVIDE_ERROR(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_INVALID_PROVIDE_ERROR = v
}

var static_ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR *DiagnosticType

func StaticGet_ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR
}
func StaticSet_ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR = v
}

var static_ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP *DiagnosticType

func StaticGet_ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP
}
func StaticSet_ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP = v
}

var static_ProcessClosurePrimitives_BASE_CLASS_ERROR *DiagnosticType

func StaticGet_ProcessClosurePrimitives_BASE_CLASS_ERROR() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_BASE_CLASS_ERROR
}
func StaticSet_ProcessClosurePrimitives_BASE_CLASS_ERROR(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_BASE_CLASS_ERROR = v
}

var static_ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR *DiagnosticType

func StaticGet_ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR
}
func StaticSet_ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR = v
}

var static_ProcessClosurePrimitives_INVALID_FORWARD_DECLARE *DiagnosticType

func StaticGet_ProcessClosurePrimitives_INVALID_FORWARD_DECLARE() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_INVALID_FORWARD_DECLARE
}
func StaticSet_ProcessClosurePrimitives_INVALID_FORWARD_DECLARE(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_INVALID_FORWARD_DECLARE = v
}

var static_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR *DiagnosticType

func StaticGet_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR
}
func StaticSet_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR = v
}

var static_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR *DiagnosticType

func StaticGet_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR
}
func StaticSet_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR = v
}

var static_ProcessClosurePrimitives_INVALID_RENAME_FUNCTION *DiagnosticType

func StaticGet_ProcessClosurePrimitives_INVALID_RENAME_FUNCTION() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_INVALID_RENAME_FUNCTION
}
func StaticSet_ProcessClosurePrimitives_INVALID_RENAME_FUNCTION(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_INVALID_RENAME_FUNCTION = v
}

var static_ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL *DiagnosticType

func StaticGet_ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL() *DiagnosticType {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL
}
func StaticSet_ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL(v *DiagnosticType) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL = v
}
func StaticGet_ProcessClosurePrimitives_GOOG() string { return "goog" }

var static_ProcessClosurePrimitives_GOOG_INHERITS pkg_rhino.QualifiedName

func StaticGet_ProcessClosurePrimitives_GOOG_INHERITS() pkg_rhino.QualifiedName {
	clinit_ProcessClosurePrimitives()
	return static_ProcessClosurePrimitives_GOOG_INHERITS
}
func StaticSet_ProcessClosurePrimitives_GOOG_INHERITS(v pkg_rhino.QualifiedName) {
	clinit_ProcessClosurePrimitives()
	static_ProcessClosurePrimitives_GOOG_INHERITS = v
}

var clinitDone_ProcessClosurePrimitives uint32
var clinitState_ProcessClosurePrimitives j2gort.ClinitState

func clinit_ProcessClosurePrimitives() {
	if atomic.LoadUint32(&clinitDone_ProcessClosurePrimitives) != 0 {
		return
	}
	clinitSlow_ProcessClosurePrimitives()
}

func clinitSlow_ProcessClosurePrimitives() {
	if !j2gort.EnterClinit(&clinitState_ProcessClosurePrimitives, &clinitDone_ProcessClosurePrimitives) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ProcessClosurePrimitives, &clinitDone_ProcessClosurePrimitives)
	static_ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR = DiagnosticType_Error("JSC_EXPECTED_OBJECTLIT_ERROR", "method \"{0}\" expected an object literal argument")
	static_ProcessClosurePrimitives_EXPECTED_STRING_ERROR = DiagnosticType_Error("JSC_EXPECTED_STRING_ERROR", "method \"{0}\" expected a string argument")
	static_ProcessClosurePrimitives_INVALID_STYLE_ERROR = DiagnosticType_Error("JSC_INVALID_CSS_NAME_MAP_STYLE_ERROR", "Invalid CSS name map style {0}")
	static_ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE = DiagnosticType_Warning("JSC_WEAK_NAMESPACE_TYPE", "Provided symbol declared with type Object. This is rarely useful. For more information see https://github.com/google/closure-compiler/wiki/A-word-about-the-type-Object")
	static_ProcessClosurePrimitives_CLASS_NAMESPACE_ERROR = DiagnosticType_Error("JSC_CLASS_NAMESPACE_ERROR", "\"{0}\" cannot be both provided and declared as a class. Try var {0} = class '{'...'}' (metadata {1})")
	static_ProcessClosurePrimitives_FUNCTION_NAMESPACE_ERROR = DiagnosticType_Error("JSC_FUNCTION_NAMESPACE_ERROR", "\"{0}\" cannot be both provided and declared as a function. (metadata {1})")
	static_ProcessClosurePrimitives_INVALID_PROVIDE_ERROR = DiagnosticType_Error("JSC_INVALID_PROVIDE_ERROR", "\"{0}\" is not a valid {1} qualified name")
	static_ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR = DiagnosticType_Error("JSC_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR", "goog.setCssNameMapping only takes an object literal with string values")
	static_ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP = DiagnosticType_Warning("INVALID_CSS_RENAMING_MAP", "Invalid entries in css renaming map: {0}")
	static_ProcessClosurePrimitives_BASE_CLASS_ERROR = DiagnosticType_Error("JSC_BASE_CLASS_ERROR", "incorrect use of {0}.base: {1}")
	static_ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR = DiagnosticType_Error("JSC_POSSIBLE_BASE_CLASS_ERROR", "potentially incorrect use of {0}.base: {1}\nNote: if this .base method is not on a Closure subclass, this error is a false positive. Suppress with /** @suppress '{closureClassChecks}' */")
	static_ProcessClosurePrimitives_INVALID_FORWARD_DECLARE = DiagnosticType_Error("JSC_INVALID_FORWARD_DECLARE", "Malformed goog.forwardDeclare")
	static_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR = DiagnosticType_Error("JSC_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR", "Closure primitive method {0} may not be aliased")
	static_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR = DiagnosticType_Error("JSC_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR", "Closure primitive method {0} may not be aliased  outside a module (ES module, CommonJS module, or goog.module)")
	static_ProcessClosurePrimitives_INVALID_RENAME_FUNCTION = DiagnosticType_Error("JSC_INVALID_RENAME_FUNCTION", "{0} call is invalid: {1}")
	static_ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL = DiagnosticType_Error("JSC_INVALID_GOOG_WEAK_USAGE", "{0} call is invalid: {1}")
	static_ProcessClosurePrimitives_GOOG_INHERITS = pkg_rhino.QualifiedName_Of("goog.inherits")
}

func (this_ *ProcessClosurePrimitives) initFields_ProcessClosurePrimitives() {
	_ = this_
	this_.knownClosureSubclasses = j2gort.NewOrderedSet[string]()
}

// Java constructor: ProcessClosurePrimitives(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewProcessClosurePrimitives(compiler AbstractCompiler) *ProcessClosurePrimitives {
	clinit_ProcessClosurePrimitives()
	this_ := &ProcessClosurePrimitives{}
	InitProcessClosurePrimitives(this_, compiler)
	return this_
}

func InitProcessClosurePrimitives(this_ *ProcessClosurePrimitives, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ProcessClosurePrimitives()
	this_.compiler = compiler
	this_.closureModules = j2gort.CheckNotNullPtr(compiler.GetModuleMetadataMap()).GetModulesByGoogNamespace()
	var provideCount int32 = int32(0)
	_ = provideCount
	{
		cur_1 := j2gort.NewMapValCursor(this_.closureModules)
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			metadata := j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](rawElem_2)
			_ = metadata
			if metadata.GetModuleMetadataMap_ModuleMetadataBase().IsGoogProvide() {
				provideCount = provideCount + 1
			}
		}
	}
	if provideCount == int32(0) {
		this_.googProvideHashes = nil
	} else {
		var cap_ int32 = int32(8)
		_ = cap_
		for cap_ < (provideCount * int32(2)) {
			cap_ = javalang.Shl32(cap_, int(int32(1)))
		}
		var table []int32 = make([]int32, cap_)
		_ = table
		var mask int32 = (cap_ - int32(1))
		_ = mask
		{
			cur_3 := j2gort.NewMapEntryCursor(this_.closureModules)
			for cur_3.NextEntry() {
				if j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](cur_3.GetValue()).GetModuleMetadataMap_ModuleMetadataBase().IsGoogProvide() {
					var h int32 = pkg_rhino.Node_HashQualifiedNameString(cur_3.KeyStr())
					_ = h
					if h != int32(0) {
						var slot int32 = (h & mask)
						_ = slot
						for (table[slot] != int32(0)) && (table[slot] != h) {
							slot = ((slot + int32(1)) & mask)
						}
						{
							a_4 := table
							i_5 := slot
							a_4[i_5] = h
						}
					}
				}
			}
		}
		this_.googProvideHashes = table
	}
}

// Java static method: ProcessClosurePrimitives.containsHash(int[] table, int h)
func ProcessClosurePrimitives_ContainsHash(table []int32, h int32) bool {
	clinit_ProcessClosurePrimitives()
	_ = table
	_ = h
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var slot int32 = (h & mask)
	_ = slot
	for {
		var existing int32 = table[slot]
		_ = existing
		if existing == int32(0) {
			return false
		}
		if existing == h {
			return true
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: ProcessClosurePrimitives.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ProcessClosurePrimitives) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_TraverseRoots(this_.compiler, this_, externs, root)
}

// Java method: ProcessClosurePrimitives.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessClosurePrimitives) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */ :
		{
			this_.priv_ProcessClosurePrimitives_checkGoogFunctions_4425(t, n)
			this_.priv_ProcessClosurePrimitives_maybeProcessClassBaseCall_4429(n)
			this_.priv_ProcessClosurePrimitives_checkPropertyRenameCall_4437(n)
		}
	case int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */ :
		if !(t.InGlobalHoistScope()) || n.IsFromExterns() {
			return
		}
		if !(NodeUtil_IsFunctionDeclaration(n)) && !(NodeUtil_IsClassDeclaration(n)) {
			return
		}
		var name string = n.GetFirstChild().GetString()
		_ = name
		var pn ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.closureModules, name))
		_ = pn
		if (pn != nil && !pn.IsJavaNil_()) && pn.GetModuleMetadataMap_ModuleMetadataBase().IsGoogProvide() {
			this_.compiler.Report_1(JSError_Make_3(n, func() *DiagnosticType {
				if n.IsClass() {
					return StaticGet_ProcessClosurePrimitives_CLASS_NAMESPACE_ERROR()
				}
				return StaticGet_ProcessClosurePrimitives_FUNCTION_NAMESPACE_ERROR()
			}(), []string{name, javalang.ToJavaString(j2gort.CheckNotNull(pn))}))
		}
	case int32(93) /* Token.EXPR_RESULT */ :
		var assign pkg_rhino.Node = n.GetFirstChild()
		_ = assign
		if !(assign.IsAssign()) {
			break
		}
		var info *pkg_rhino.JSDocInfo = assign.GetJSDocInfo()
		_ = info
		if ((info == nil) || (this_.googProvideHashes == nil)) || n.IsFromExterns() {
			break
		}
		var target pkg_rhino.Node = assign.GetFirstChild()
		_ = target
		var h int32 = target.HashQualifiedName()
		_ = h
		if (h == int32(0)) || !(ProcessClosurePrimitives_ContainsHash(this_.googProvideHashes, h)) {
			break
		}
		var lhs string = target.GetQualifiedName()
		_ = lhs
		if lhs != j2gort.NullString {
			this_.priv_ProcessClosurePrimitives_checkPossibleGoogProvideInit_4424(lhs, info, n)
		}
	case int32(84) /* Token.VAR */, int32(97) /* Token.CONST */, int32(111) /* Token.LET */ :
		if !(n.GetFirstChild().IsName()) {
			break
		}
		this_.priv_ProcessClosurePrimitives_checkPossibleGoogProvideInit_4424(n.GetFirstChild().GetString(), n.GetJSDocInfo(), n)
	default:
	}
}

// Java method: ProcessClosurePrimitives.checkPossibleGoogProvideInit(java.lang.String namespace, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, com.google.javascript.rhino.Node definition)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_checkPossibleGoogProvideInit_4424(namespace string, info *pkg_rhino.JSDocInfo, definition pkg_rhino.Node) {
	_ = namespace
	_ = info
	_ = definition
	if (info == nil) || definition.IsFromExterns() {
		return
	}
	var metadata ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.closureModules, namespace))
	_ = metadata
	if (metadata == nil || metadata.IsJavaNil_()) || !(metadata.GetModuleMetadataMap_ModuleMetadataBase().IsGoogProvide()) {
		return
	}
	var expr pkg_rhino.JSTypeExpression = info.GetType()
	_ = expr
	if expr == (pkg_rhino.JSTypeExpression{}) {
		return
	}
	var n pkg_rhino.Node = expr.GetRoot()
	_ = n
	if n.RawToken_() == uint8(159) {
		n = n.GetFirstChild()
	}
	if (n.IsStringLit() && !(n.HasChildren())) && (j2gort.CheckNotNullString(n.GetString()) == "Object") {
		this_.compiler.Report_1(JSError_Make_3(definition, StaticGet_ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE(), []string{}))
	}
}

// Java method: ProcessClosurePrimitives.checkGoogFunctions(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_checkGoogFunctions_4425(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var callee pkg_rhino.Node = call.GetFirstChild()
	_ = callee
	if !(callee.IsGetProp()) {
		return
	}
	var receiver pkg_rhino.Node = callee.GetFirstChild()
	_ = receiver
	if !(receiver.IsName()) || !(func() bool {
		sr_6 := receiver.GetString()
		se_7 := "goog"
		return j2gort.CheckNotNullString(sr_6) == se_7
	}()) {
		return
	}
	var methodName string = callee.GetString()
	_ = methodName
	switch j2gort.CheckNotNullString(methodName) {
	case "inherits":
		this_.priv_ProcessClosurePrimitives_processInheritsCall_4432(call)
	case "addDependency":
		if this_.priv_ProcessClosurePrimitives_validateUnaliasablePrimitiveCall_4427(t, call, methodName) {
			this_.priv_ProcessClosurePrimitives_processAddDependency_4435(call)
		}
	case "setCssNameMapping":
		var unused CssRenamingMap = ProcessClosurePrimitives_ProcessSetCssNameMapping(this_.compiler, call, call.GetParent())
		_ = unused
	case "forwardDeclare":
		if this_.priv_ProcessClosurePrimitives_validatePrimitiveCallWithMessage_4428(t, call, methodName, StaticGet_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR()) {
			this_.priv_ProcessClosurePrimitives_processForwardDeclare_4436(call)
		}
	case "weakUsage":
		this_.priv_ProcessClosurePrimitives_validateWeakUsageCall_4426(call)
	default:
	}
}

// Java method: ProcessClosurePrimitives.validateWeakUsageCall(com.google.javascript.rhino.Node call)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_validateWeakUsageCall_4426(call pkg_rhino.Node) {
	_ = call
	var childCount int32 = call.GetChildCount()
	_ = childCount
	var arg pkg_rhino.Node = call.GetSecondChild()
	_ = arg
	var calleeName string = call.GetFirstChild().GetQualifiedName()
	_ = calleeName
	if childCount != int32(2) {
		this_.compiler.Report_1(JSError_Make_3(call, StaticGet_ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL(), []string{calleeName, ("should have exactly one argument, not " + j2gort.StrValInt32((childCount - int32(1))))}))
	}
	if (childCount >= int32(2)) && !(call.GetSecondChild().IsQualifiedName()) {
		this_.compiler.Report_1(JSError_Make_3(call, StaticGet_ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL(), []string{calleeName, ("argument should be a name or qualified name, not " + j2gort.StrVal(arg))}))
	}
}

// Java method: ProcessClosurePrimitives.validateUnaliasablePrimitiveCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, java.lang.String methodName)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_validateUnaliasablePrimitiveCall_4427(t *NodeTraversal, n pkg_rhino.Node, methodName string) bool {
	_ = t
	_ = n
	_ = methodName
	return this_.priv_ProcessClosurePrimitives_validatePrimitiveCallWithMessage_4428(t, n, methodName, StaticGet_ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR())
}

// Java method: ProcessClosurePrimitives.validatePrimitiveCallWithMessage(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, java.lang.String methodName, com.google.javascript.jscomp.DiagnosticType invalidAliasingError)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_validatePrimitiveCallWithMessage_4428(t *NodeTraversal, n pkg_rhino.Node, methodName string, invalidAliasingError *DiagnosticType) bool {
	_ = t
	_ = n
	_ = methodName
	_ = invalidAliasingError
	if this_.compiler.GetOptions().ShouldPreserveGoogModule() {
		return true
	}
	if !(t.InGlobalHoistScope()) && !(t.InModuleScope()) {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR(), []string{}))
		return false
	} else if !(n.GetParent().IsExprResult()) && !(t.InModuleScope()) {
		this_.compiler.Report_1(JSError_Make_3(n, invalidAliasingError, []string{("goog" + "." + j2gort.StrValString(methodName))}))
		return false
	}
	return true
}

// Java method: ProcessClosurePrimitives.maybeProcessClassBaseCall(com.google.javascript.rhino.Node n)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_maybeProcessClassBaseCall_4429(n pkg_rhino.Node) {
	_ = n
	var callTarget pkg_rhino.Node = n.GetFirstChild()
	_ = callTarget
	if !(callTarget.IsGetProp()) || !(j2gort.CheckNotNullString(callTarget.GetString()) == "base") {
		return
	}
	var baseContainerNode pkg_rhino.Node = callTarget.GetFirstChild()
	_ = baseContainerNode
	if !(baseContainerNode.IsUnscopedQualifiedName()) {
		return
	}
	var baseContainer string = callTarget.GetFirstChild().GetQualifiedName()
	_ = baseContainer
	var enclosingFnNameNode pkg_rhino.Node = ProcessClosurePrimitives_GetEnclosingDeclNameNode(n)
	_ = enclosingFnNameNode
	if (enclosingFnNameNode == 0) || !(enclosingFnNameNode.IsUnscopedQualifiedName()) {
		if this_.knownClosureSubclasses.Contains(baseContainer) {
			this_.priv_ProcessClosurePrimitives_reportBadBaseMethodUse_4433(n, baseContainer, "Could not find enclosing method.")
		} else if ProcessClosurePrimitives_BaseUsedInClass(n) {
			var clazz pkg_rhino.Node = NodeUtil_GetEnclosingClass(n)
			_ = clazz
			if (clazz.GetFirstChild().IsName() && (j2gort.CheckNotNullString(clazz.GetFirstChild().GetString()) == baseContainer)) || (clazz.GetSecondChild().IsName() && (j2gort.CheckNotNullString(clazz.GetSecondChild().GetString()) == baseContainer)) {
				this_.priv_ProcessClosurePrimitives_reportBadBaseMethodUse_4433(n, clazz.GetFirstChild().GetString(), "base method is not allowed in ES6 class. Use super instead.")
			}
		}
		return
	}
	if ProcessClosurePrimitives_BaseUsedInClass(n) {
		var clazz_2 pkg_rhino.Node = NodeUtil_GetEnclosingClass(n)
		_ = clazz_2
		var name string = NodeUtil_GetBestLValueName(clazz_2)
		_ = name
		this_.priv_ProcessClosurePrimitives_reportBadBaseMethodUse_4433(n, name, "base method is not allowed in ES6 class. Use super instead.")
		return
	}
	var enclosingQname string = enclosingFnNameNode.GetQualifiedName()
	_ = enclosingQname
	if !(j2gort.JStrIndexOf(enclosingQname, ".prototype.") >= 0) {
		this_.priv_ProcessClosurePrimitives_rewriteBaseCallInConstructor_4430(enclosingQname, baseContainer, n, enclosingFnNameNode)
	} else {
		this_.priv_ProcessClosurePrimitives_rewriteBaseCallInMethod_4431(enclosingQname, baseContainer, n, enclosingFnNameNode)
	}
}

// Java method: ProcessClosurePrimitives.rewriteBaseCallInConstructor(java.lang.String enclosingQname, java.lang.String baseContainer, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node enclosingFnNameNode)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_rewriteBaseCallInConstructor_4430(enclosingQname string, baseContainer string, n pkg_rhino.Node, enclosingFnNameNode pkg_rhino.Node) {
	_ = enclosingQname
	_ = baseContainer
	_ = n
	_ = enclosingFnNameNode
	if !(j2gort.CheckNotNullString(enclosingQname) == baseContainer) {
		if this_.knownClosureSubclasses.Contains(baseContainer) {
			this_.priv_ProcessClosurePrimitives_reportBadBaseMethodUse_4433(n, baseContainer, ("Must be used within " + j2gort.StrValString(baseContainer) + " methods"))
		}
		return
	}
	var enclosingParent pkg_rhino.Node = enclosingFnNameNode.GetParent()
	_ = enclosingParent
	var baseClassNode pkg_rhino.Node = ProcessClosurePrimitives_FindGoogInheritsCall(func() pkg_rhino.Node {
		if enclosingParent.IsAssign() {
			return enclosingParent.GetParent()
		}
		return enclosingParent
	}())
	_ = baseClassNode
	if baseClassNode == 0 {
		return
	}
	var callee pkg_rhino.Node = n.GetFirstChild()
	_ = callee
	var thisArg pkg_rhino.Node = callee.GetNext()
	_ = thisArg
	if (thisArg == 0) || !(thisArg.IsThis()) {
		this_.priv_ProcessClosurePrimitives_reportBadBaseMethodUse_4433(n, baseContainer, "First argument must be 'this'.")
		return
	}
	var methodNameNode pkg_rhino.Node = thisArg.GetNext()
	_ = methodNameNode
	if ((methodNameNode == 0) || !(methodNameNode.IsStringLit())) || !(j2gort.CheckNotNullString(methodNameNode.GetString()) == "constructor") {
		this_.priv_ProcessClosurePrimitives_reportBadBaseMethodUse_4433(n, baseContainer, "Second argument must be 'constructor'.")
		return
	}
	callee.ReplaceWith(NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.compiler, (j2gort.StrValString(baseClassNode.GetQualifiedName()) + ".call"), callee, (j2gort.StrValString(enclosingQname) + ".base")))
	methodNameNode.Detach()
	this_.compiler.ReportChangeToEnclosingScope(n)
}

// Java method: ProcessClosurePrimitives.rewriteBaseCallInMethod(java.lang.String enclosingQname, java.lang.String baseContainer, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node enclosingFnNameNode)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_rewriteBaseCallInMethod_4431(enclosingQname string, baseContainer string, n pkg_rhino.Node, enclosingFnNameNode pkg_rhino.Node) {
	_ = enclosingQname
	_ = baseContainer
	_ = n
	_ = enclosingFnNameNode
	if !(this_.knownClosureSubclasses.Contains(baseContainer)) {
		this_.priv_ProcessClosurePrimitives_reportPossibleBadBaseMethodUse_4434(n, baseContainer, ("base used outside a known Closure subclass, in: " + j2gort.StrValString(baseContainer) + "\nIf " + j2gort.StrValString(baseContainer) + " is actually a Closure subclass, the compiler is unable to statically determine this. This can happen when calling .base off of some import or other alias of the original Closure subclass, such as from a goog.require in another file. Instead of using .base, rewrite this call as ParentClass.prototype.method.call(this)"))
		return
	}
	var misuseOfBase bool = !(enclosingFnNameNode.GetFirstFirstChild().MatchesQualifiedName_1_string(baseContainer))
	_ = misuseOfBase
	if misuseOfBase {
		this_.priv_ProcessClosurePrimitives_reportBadBaseMethodUse_4433(n, baseContainer, ("Must be used within " + j2gort.StrValString(baseContainer) + " methods"))
		return
	}
	var callee pkg_rhino.Node = n.GetFirstChild()
	_ = callee
	var thisArg pkg_rhino.Node = callee.GetNext()
	_ = thisArg
	if (thisArg == 0) || !(thisArg.IsThis()) {
		this_.priv_ProcessClosurePrimitives_reportBadBaseMethodUse_4433(n, baseContainer, "First argument must be 'this'.")
		return
	}
	var methodNameNode pkg_rhino.Node = thisArg.GetNext()
	_ = methodNameNode
	if (methodNameNode == 0) || !(methodNameNode.IsStringLit()) {
		this_.priv_ProcessClosurePrimitives_reportBadBaseMethodUse_4433(n, baseContainer, "Second argument must name a method.")
		return
	}
	var methodName string = methodNameNode.GetString()
	_ = methodName
	var ending string = (".prototype." + j2gort.StrValString(methodName))
	_ = ending
	if (enclosingQname == j2gort.NullString) || !(javalang.StringEndsWith(j2gort.CheckNotNullString(enclosingQname), j2gort.CheckNotNullString(ending))) {
		this_.priv_ProcessClosurePrimitives_reportBadBaseMethodUse_4433(n, baseContainer, ("Enclosing method does not match " + j2gort.StrValString(methodName)))
		return
	}
	var className pkg_rhino.Node = enclosingFnNameNode.GetFirstFirstChild()
	_ = className
	callee.ReplaceWith(NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.compiler, (j2gort.StrValString(className.GetQualifiedName()) + ".superClass_." + j2gort.StrValString(methodName) + ".call"), callee, (j2gort.StrValString(enclosingQname) + ".base")))
	methodNameNode.Detach()
	this_.compiler.ReportChangeToEnclosingScope(n)
}

// Java static method: ProcessClosurePrimitives.findGoogInheritsCall(com.google.javascript.rhino.Node ctorDeclarationNode)
func ProcessClosurePrimitives_FindGoogInheritsCall(ctorDeclarationNode pkg_rhino.Node) pkg_rhino.Node {
	clinit_ProcessClosurePrimitives()
	_ = ctorDeclarationNode
	var maybeInheritsExpr pkg_rhino.Node = ctorDeclarationNode.GetNext()
	_ = maybeInheritsExpr
	for (maybeInheritsExpr != 0) && (maybeInheritsExpr.IsEmpty() || NodeUtil_IsExprAssign(maybeInheritsExpr)) {
		maybeInheritsExpr = maybeInheritsExpr.GetNext()
	}
	var baseClassNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = baseClassNode
	if (maybeInheritsExpr != 0) && NodeUtil_IsExprCall(maybeInheritsExpr) {
		var callNode pkg_rhino.Node = maybeInheritsExpr.GetFirstChild()
		_ = callNode
		if StaticGet_ProcessClosurePrimitives_GOOG_INHERITS().Matches(callNode.GetFirstChild()) && callNode.GetLastChild().IsQualifiedName() {
			baseClassNode = callNode.GetLastChild()
		}
	}
	return baseClassNode
}

// Java method: ProcessClosurePrimitives.processInheritsCall(com.google.javascript.rhino.Node n)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_processInheritsCall_4432(n pkg_rhino.Node) {
	_ = n
	if n.HasXChildren(int32(3)) {
		var subClass pkg_rhino.Node = n.GetSecondChild()
		_ = subClass
		var superClass pkg_rhino.Node = subClass.GetNext()
		_ = superClass
		if subClass.IsUnscopedQualifiedName() && superClass.IsUnscopedQualifiedName() {
			this_.knownClosureSubclasses.Add(subClass.GetQualifiedName())
		}
	}
}

// Java static method: ProcessClosurePrimitives.getEnclosingDeclNameNode(com.google.javascript.rhino.Node n)
func ProcessClosurePrimitives_GetEnclosingDeclNameNode(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_ProcessClosurePrimitives()
	_ = n
	var fn pkg_rhino.Node = NodeUtil_GetEnclosingFunction(n)
	_ = fn
	if fn == 0 {
		return pkg_rhino.Node(0)
	} else {
		return NodeUtil_GetNameNode(fn)
	}
}

// Java static method: ProcessClosurePrimitives.baseUsedInClass(com.google.javascript.rhino.Node n)
func ProcessClosurePrimitives_BaseUsedInClass(n pkg_rhino.Node) bool {
	clinit_ProcessClosurePrimitives()
	_ = n
	{
		var curr pkg_rhino.Node = n
		_ = curr
		for ; curr != 0; curr = curr.GetParent() {
			if curr.IsClassMembers() {
				return true
			}
		}
	}
	return false
}

// Java method: ProcessClosurePrimitives.reportBadBaseMethodUse(com.google.javascript.rhino.Node n, java.lang.String className, java.lang.String extraMessage)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_reportBadBaseMethodUse_4433(n pkg_rhino.Node, className string, extraMessage string) {
	_ = n
	_ = className
	_ = extraMessage
	this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ProcessClosurePrimitives_BASE_CLASS_ERROR(), []string{className, extraMessage}))
}

// Java method: ProcessClosurePrimitives.reportPossibleBadBaseMethodUse(com.google.javascript.rhino.Node n, java.lang.String className, java.lang.String extraMessage)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_reportPossibleBadBaseMethodUse_4434(n pkg_rhino.Node, className string, extraMessage string) {
	_ = n
	_ = className
	_ = extraMessage
	this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR(), []string{className, extraMessage}))
}

// Java static method: ProcessClosurePrimitives.processSetCssNameMapping(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func ProcessClosurePrimitives_ProcessSetCssNameMapping(compiler AbstractCompiler, n pkg_rhino.Node, parent pkg_rhino.Node) CssRenamingMap {
	clinit_ProcessClosurePrimitives()
	_ = compiler
	_ = n
	_ = parent
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	var arg pkg_rhino.Node = left.GetNext()
	_ = arg
	if !(ProcessClosurePrimitives_VerifySetCssNameMapping(compiler, left, arg)) {
		return nil
	}
	var cssNames j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = cssNames
	{
		var key pkg_rhino.Node = arg.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			var value pkg_rhino.Node = key.GetFirstChild()
			_ = value
			if (!(key.IsStringKey()) || (value == 0)) || !(value.IsStringLit()) {
				compiler.Report_1(JSError_Make_3(n, StaticGet_ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR(), []string{}))
				return nil
			}
			j2gort.MapPutStr(cssNames, key.GetString(), j2gort.BoxString(value.GetString()))
		}
	}
	var styleStr string = "BY_PART"
	_ = styleStr
	if arg.GetNext() != 0 {
		styleStr = arg.GetNext().GetString()
	}
	var style CssRenamingMap_Style = CssRenamingMap_Style(0)
	_ = style
	var didRet_8 bool
	var retVal_9 CssRenamingMap
	func() {
		defer func() {
			if rec_10 := recover(); rec_10 != nil {
				jex_11 := javalang.RecoverException(rec_10)
				_ = jex_11
				if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_11) {
					e := j2gort.AsIface[javalang.IIllegalArgumentException](jex_11)
					_ = e
					compiler.Report_1(JSError_Make_3(n, StaticGet_ProcessClosurePrimitives_INVALID_STYLE_ERROR(), []string{styleStr}))
					retVal_9 = nil
					didRet_8 = true
					return
				} else {
					panic(rec_10)
				}
			}
		}()
		style = CssRenamingMap_Style_ValueOf(styleStr)
	}()
	if didRet_8 {
		return retVal_9
	}
	if style == StaticGet_CssRenamingMap_Style_BY_PART() {
		var errors []string = []string{}
		_ = errors
		{
			cur_12 := j2gort.NewMapKeyStrCursor(cssNames)
			for cur_12.HasNext() {
				key_2 := cur_12.NextStr()
				_ = key_2
				if (j2gort.JStrIndexOf(key_2, "-") >= 0) && !(javalang.StringStartsWith(j2gort.CheckNotNullString(key_2), "--")) {
					errors = append(errors, key_2)
				}
			}
		}
		if len(errors) != 0 {
			compiler.Report_1(JSError_Make_3(n, StaticGet_ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP(), []string{j2gort.SliceFormatJavaString(errors)}))
		}
	} else if style == StaticGet_CssRenamingMap_Style_BY_WHOLE() {
		var errors_2 []string = []string{}
		_ = errors_2
		{
			cur_13 := j2gort.NewMapEntryCursor(cssNames)
			for cur_13.NextEntry() {
				if j2gort.JStrLen(cur_13.KeyStr()) > int32(10) {
					continue
				}
				{
					cur_14 := j2gort.NewMapEntryCursor(cssNames)
					for cur_14.NextEntry() {
						var combined string = j2gort.CastString(j2gort.MapGetStr(cssNames, (j2gort.StrValString(cur_14.KeyStr()) + "-" + j2gort.StrValString(cur_13.KeyStr()))))
						_ = combined
						if (combined != j2gort.NullString) && !(func() bool {
							sr_15 := combined
							se_16 := (j2gort.StrValString(j2gort.CastString(cur_14.GetValue())) + "-" + j2gort.StrValString(j2gort.CastString(cur_13.GetValue())))
							return j2gort.CheckNotNullString(sr_15) == se_16
						}()) {
							errors_2 = append(errors_2, "map("+j2gort.StrValString(cur_14.KeyStr())+"-"+j2gort.StrValString(cur_13.KeyStr())+") != map("+j2gort.StrValString(cur_14.KeyStr())+")-map("+j2gort.StrValString(cur_13.KeyStr())+")")
						}
					}
				}
			}
		}
		if len(errors_2) != 0 {
			compiler.Report_1(JSError_Make_3(n, StaticGet_ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP(), []string{j2gort.SliceFormatJavaString(errors_2)}))
		}
	}
	return NewProcessClosurePrimitives__Anon_72(cssNames, style)
}

// Java method: ProcessClosurePrimitives.processAddDependency(com.google.javascript.rhino.Node n)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_processAddDependency_4435(n pkg_rhino.Node) {
	_ = n
	var emptyNode pkg_rhino.Node = pkg_rhino.IR_Number(float64(int32(0)))
	_ = emptyNode
	n.ReplaceWith(emptyNode)
	this_.compiler.ReportChangeToEnclosingScope(emptyNode)
}

// Java method: ProcessClosurePrimitives.processForwardDeclare(com.google.javascript.rhino.Node n)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_processForwardDeclare_4436(n pkg_rhino.Node) {
	_ = n
	if !(n.GetParent().IsExprResult()) {
		return
	}
	var convention CodingConvention = this_.compiler.GetCodingConvention()
	_ = convention
	var typeDeclaration string = j2gort.NullString
	_ = typeDeclaration
	func() {
		defer func() {
			if rec_17 := recover(); rec_17 != nil {
				jex_18 := javalang.RecoverException(rec_17)
				_ = jex_18
				if j2gort.InstanceOf[javalang.INullPointerException](jex_18) || j2gort.InstanceOf[javalang.INoSuchElementException](jex_18) || j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_18) {
					e := jex_18
					_ = e
					this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ProcessClosurePrimitives_INVALID_FORWARD_DECLARE(), []string{"A single type could not identified for the goog.forwardDeclare statement"}))
				} else {
					panic(rec_17)
				}
			}
		}()
		typeDeclaration = j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](convention.IdentifyTypeDeclarationCall(n))))
	}()
	if typeDeclaration != j2gort.NullString {
		this_.compiler.ForwardDeclareType(typeDeclaration)
	}
}

// Java static method: ProcessClosurePrimitives.verifySetCssNameMapping(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node methodName, com.google.javascript.rhino.Node firstArg)
func ProcessClosurePrimitives_VerifySetCssNameMapping(compiler AbstractCompiler, methodName pkg_rhino.Node, firstArg pkg_rhino.Node) bool {
	clinit_ProcessClosurePrimitives()
	_ = compiler
	_ = methodName
	_ = firstArg
	var diagnostic *DiagnosticType = nil
	_ = diagnostic
	if firstArg == 0 {
		diagnostic = StaticGet_ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR()
	} else if !(firstArg.IsObjectLit()) {
		diagnostic = StaticGet_ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR()
	} else if firstArg.GetNext() != 0 {
		var secondArg pkg_rhino.Node = firstArg.GetNext()
		_ = secondArg
		if !(secondArg.IsStringLit()) {
			diagnostic = StaticGet_ProcessClosurePrimitives_EXPECTED_STRING_ERROR()
		} else if secondArg.GetNext() != 0 {
			diagnostic = StaticGet_ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR()
		}
	}
	if diagnostic != nil {
		compiler.Report_1(JSError_Make_3(methodName, diagnostic, []string{methodName.GetQualifiedName()}))
		return false
	}
	return true
}

// Java method: ProcessClosurePrimitives.checkPropertyRenameCall(com.google.javascript.rhino.Node call)
func (this_ *ProcessClosurePrimitives) priv_ProcessClosurePrimitives_checkPropertyRenameCall_4437(call pkg_rhino.Node) {
	_ = call
	var callee pkg_rhino.Node = call.GetFirstChild()
	_ = callee
	if !(this_.compiler.GetCodingConvention().IsPropertyRenameFunction(callee)) {
		return
	}
	var calleeName string = callee.GetQualifiedName()
	_ = calleeName
	switch call.GetChildCount() - int32(1) {
	case int32(1), int32(2):
	default:
		this_.compiler.Report_1(JSError_Make_3(call, StaticGet_ProcessClosurePrimitives_INVALID_RENAME_FUNCTION(), []string{calleeName, "Must be called with 1 or 2 arguments."}))
	}
	var propName pkg_rhino.Node = callee.GetNext()
	_ = propName
	if (propName == 0) || !(propName.IsStringLit()) {
		this_.compiler.Report_1(JSError_Make_3(call, StaticGet_ProcessClosurePrimitives_INVALID_RENAME_FUNCTION(), []string{calleeName, "The first argument must be a string literal."}))
	} else if j2gort.JStrIndexOf(propName.GetString(), ".") >= 0 {
		this_.compiler.Report_1(JSError_Make_3(call, StaticGet_ProcessClosurePrimitives_INVALID_RENAME_FUNCTION(), []string{calleeName, "The first argument must not be a property path."}))
	}
}

// Java class:
type ProcessClosurePrimitives__Anon_72 struct {
	cap_0_cssNames j2gort.Ext_java_util_Map
	idHash_        atomic.Int32
	cap_1_style    CssRenamingMap_Style
}

func (this_ *ProcessClosurePrimitives__Anon_72) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessClosurePrimitives__Anon_72) IsJavaNil_() bool { return this_ == nil }
func (*ProcessClosurePrimitives__Anon_72) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessClosurePrimitives$1", ""
}
func (this_ *ProcessClosurePrimitives__Anon_72) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessClosurePrimitives__Anon_72) HasIdentityEquals_()               {}
func (*ProcessClosurePrimitives__Anon_72) Implements_jscomp_CssRenamingMap() {}
func (*ProcessClosurePrimitives__Anon_72) Implements_jscomp_RenamingMap()    {}

func (this_ *ProcessClosurePrimitives__Anon_72) initFields_ProcessClosurePrimitives__Anon_72() {
	_ = this_
}

func NewProcessClosurePrimitives__Anon_72(cap_p_cssNames j2gort.Ext_java_util_Map, cap_p_style CssRenamingMap_Style) *ProcessClosurePrimitives__Anon_72 {
	this_ := &ProcessClosurePrimitives__Anon_72{}
	this_.cap_0_cssNames = cap_p_cssNames
	this_.cap_1_style = cap_p_style
	this_.initFields_ProcessClosurePrimitives__Anon_72()
	return this_
}

// Java method: .get(java.lang.String value)
func (this_ *ProcessClosurePrimitives__Anon_72) Get(value string) string {
	_ = value
	if j2gort.MapContainsKeyStr(this_.cap_0_cssNames, value) {
		return j2gort.CastString(j2gort.MapGetStr(this_.cap_0_cssNames, value))
	} else {
		return value
	}
}

// Java method: .getStyle()
func (this_ *ProcessClosurePrimitives__Anon_72) GetStyle() CssRenamingMap_Style {
	return this_.cap_1_style
}
func (*ProcessClosurePrimitives) AlwaysShouldTraverse_()      {}
func (*ProcessClosurePrimitives) CanReleaseSyntacticScopes_() {}
