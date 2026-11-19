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
// Port of com/google/javascript/jscomp/CheckAccessControls.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_check_access_controls_1 j2gort.Ext_com_google_common_base_Supplier
var lambdaSingleton_check_access_controls_2 j2gort.Ext_java_util_function_Supplier

func init() {
	lambdaSingleton_check_access_controls_1 = j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any { return j2gort.NewExt_java_util_LinkedHashMap() })
	lambdaSingleton_check_access_controls_2 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return ""
	})
}

// Java class: com.google.javascript.jscomp.CheckAccessControls
type CheckAccessControls struct {
	compiler                  AbstractCompiler                                  `gson:"compiler"`
	typeRegistry              *pkg_rhino.JSTypeRegistry                         `gson:"typeRegistry"`
	currentClassStack         j2gort.Ext_java_util_Deque                        `gson:"currentClassStack"`
	constPropertyInits        j2gort.Ext_com_google_common_collect_Table        `gson:"constPropertyInits"`
	defaultVisibilityForFiles j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"defaultVisibilityForFiles"`
	deprecationDepth          int32                                             `gson:"deprecationDepth"`
	idHash_                   atomic.Int32
}

func (this_ *CheckAccessControls) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckAccessControls) IsJavaNil_() bool { return this_ == nil }
func (*CheckAccessControls) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckAccessControls", "CheckAccessControls"
}
func (this_ *CheckAccessControls) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckAccessControls) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckAccessControls", "com.google.javascript.jscomp.CheckAccessControls", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckAccessControls](v_)
	})
}
func (*CheckAccessControls) Implements_jscomp_NodeTraversal_Callback() {}
func (*CheckAccessControls) Implements_jscomp_CompilerPass()           {}

var static_CheckAccessControls_DEPRECATED_NAME *DiagnosticType

func StaticGet_CheckAccessControls_DEPRECATED_NAME() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_DEPRECATED_NAME
}
func StaticSet_CheckAccessControls_DEPRECATED_NAME(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_DEPRECATED_NAME = v
}

var static_CheckAccessControls_DEPRECATED_NAME_REASON *DiagnosticType

func StaticGet_CheckAccessControls_DEPRECATED_NAME_REASON() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_DEPRECATED_NAME_REASON
}
func StaticSet_CheckAccessControls_DEPRECATED_NAME_REASON(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_DEPRECATED_NAME_REASON = v
}

var static_CheckAccessControls_DEPRECATED_PROP *DiagnosticType

func StaticGet_CheckAccessControls_DEPRECATED_PROP() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_DEPRECATED_PROP
}
func StaticSet_CheckAccessControls_DEPRECATED_PROP(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_DEPRECATED_PROP = v
}

var static_CheckAccessControls_DEPRECATED_PROP_REASON *DiagnosticType

func StaticGet_CheckAccessControls_DEPRECATED_PROP_REASON() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_DEPRECATED_PROP_REASON
}
func StaticSet_CheckAccessControls_DEPRECATED_PROP_REASON(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_DEPRECATED_PROP_REASON = v
}

var static_CheckAccessControls_DEPRECATED_CLASS *DiagnosticType

func StaticGet_CheckAccessControls_DEPRECATED_CLASS() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_DEPRECATED_CLASS
}
func StaticSet_CheckAccessControls_DEPRECATED_CLASS(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_DEPRECATED_CLASS = v
}

var static_CheckAccessControls_DEPRECATED_CLASS_REASON *DiagnosticType

func StaticGet_CheckAccessControls_DEPRECATED_CLASS_REASON() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_DEPRECATED_CLASS_REASON
}
func StaticSet_CheckAccessControls_DEPRECATED_CLASS_REASON(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_DEPRECATED_CLASS_REASON = v
}

var static_CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS *DiagnosticType

func StaticGet_CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS
}
func StaticSet_CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS = v
}

var static_CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS *DiagnosticType

func StaticGet_CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS
}
func StaticSet_CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS = v
}

var static_CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS *DiagnosticType

func StaticGet_CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS
}
func StaticSet_CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS = v
}

var static_CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS *DiagnosticType

func StaticGet_CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS
}
func StaticSet_CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS = v
}

var static_CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY *DiagnosticType

func StaticGet_CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY
}
func StaticSet_CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY = v
}

var static_CheckAccessControls_PRIVATE_OVERRIDE *DiagnosticType

func StaticGet_CheckAccessControls_PRIVATE_OVERRIDE() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_PRIVATE_OVERRIDE
}
func StaticSet_CheckAccessControls_PRIVATE_OVERRIDE(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_PRIVATE_OVERRIDE = v
}

var static_CheckAccessControls_EXTEND_FINAL_CLASS *DiagnosticType

func StaticGet_CheckAccessControls_EXTEND_FINAL_CLASS() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_EXTEND_FINAL_CLASS
}
func StaticSet_CheckAccessControls_EXTEND_FINAL_CLASS(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_EXTEND_FINAL_CLASS = v
}

var static_CheckAccessControls_VISIBILITY_MISMATCH *DiagnosticType

func StaticGet_CheckAccessControls_VISIBILITY_MISMATCH() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_VISIBILITY_MISMATCH
}
func StaticSet_CheckAccessControls_VISIBILITY_MISMATCH(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_VISIBILITY_MISMATCH = v
}

var static_CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE *DiagnosticType

func StaticGet_CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE
}
func StaticSet_CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE = v
}

var static_CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN *DiagnosticType

func StaticGet_CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN
}
func StaticSet_CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN = v
}

var static_CheckAccessControls_CONST_PROPERTY_DELETED *DiagnosticType

func StaticGet_CheckAccessControls_CONST_PROPERTY_DELETED() *DiagnosticType {
	clinit_CheckAccessControls()
	return static_CheckAccessControls_CONST_PROPERTY_DELETED
}
func StaticSet_CheckAccessControls_CONST_PROPERTY_DELETED(v *DiagnosticType) {
	clinit_CheckAccessControls()
	static_CheckAccessControls_CONST_PROPERTY_DELETED = v
}

var clinitDone_CheckAccessControls uint32
var clinitState_CheckAccessControls j2gort.ClinitState

func clinit_CheckAccessControls() {
	if atomic.LoadUint32(&clinitDone_CheckAccessControls) != 0 {
		return
	}
	clinitSlow_CheckAccessControls()
}

func clinitSlow_CheckAccessControls() {
	if !j2gort.EnterClinit(&clinitState_CheckAccessControls, &clinitDone_CheckAccessControls) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckAccessControls, &clinitDone_CheckAccessControls)
	static_CheckAccessControls_DEPRECATED_NAME = DiagnosticType_Disabled("JSC_DEPRECATED_VAR", "Variable {0} has been deprecated.")
	static_CheckAccessControls_DEPRECATED_NAME_REASON = DiagnosticType_Disabled("JSC_DEPRECATED_VAR_REASON", "Variable {0} has been deprecated: {1}")
	static_CheckAccessControls_DEPRECATED_PROP = DiagnosticType_Disabled("JSC_DEPRECATED_PROP", "Property {0} of type {1} has been deprecated.")
	static_CheckAccessControls_DEPRECATED_PROP_REASON = DiagnosticType_Disabled("JSC_DEPRECATED_PROP_REASON", "Property {0} of type {1} has been deprecated: {2}")
	static_CheckAccessControls_DEPRECATED_CLASS = DiagnosticType_Disabled("JSC_DEPRECATED_CLASS", "Class {0} has been deprecated.")
	static_CheckAccessControls_DEPRECATED_CLASS_REASON = DiagnosticType_Disabled("JSC_DEPRECATED_CLASS_REASON", "Class {0} has been deprecated: {1}")
	static_CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS = DiagnosticType_Error("JSC_BAD_PACKAGE_PROPERTY_ACCESS", "Access to package-private property {0} of {1} not allowed here.")
	static_CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS = DiagnosticType_Error("JSC_BAD_PRIVATE_GLOBAL_ACCESS", "Access to private variable {0} not allowed outside file {1}.")
	static_CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS = DiagnosticType_Warning("JSC_BAD_PRIVATE_PROPERTY_ACCESS", "Access to private property {0} of {1} not allowed here.")
	static_CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS = DiagnosticType_Warning("JSC_BAD_PROTECTED_PROPERTY_ACCESS", "Access to protected property {0} of {1} not allowed here.")
	static_CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY = DiagnosticType_Error("JSC_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY", "Overridden property {0} in file with fileoverview visibility {1} must explicitly redeclare superclass visibility")
	static_CheckAccessControls_PRIVATE_OVERRIDE = DiagnosticType_Warning("JSC_PRIVATE_OVERRIDE", "Overriding private property of {0}.")
	static_CheckAccessControls_EXTEND_FINAL_CLASS = DiagnosticType_Error("JSC_EXTEND_FINAL_CLASS", "{0} is not allowed to extend final class {1}.")
	static_CheckAccessControls_VISIBILITY_MISMATCH = DiagnosticType_Warning("JSC_VISIBILITY_MISMATCH", "Overriding {0} property of {1} with {2} property.")
	static_CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE = DiagnosticType_Warning("JSC_CONSTANT_PROPERTY_REASSIGNED_VALUE", pkg_base.JSCompStrings_Lines_1_slice_string([]string{"constant property {0} assigned a value more than once", "Initialized at {1}"}))
	static_CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN = DiagnosticType_Warning("JSC_FINAL_PROPERTY_OVERRIDDEN", pkg_base.JSCompStrings_Lines_1_slice_string([]string{"@final method or property {0} overridden", "Initialized at {1}"}))
	static_CheckAccessControls_CONST_PROPERTY_DELETED = DiagnosticType_Warning("JSC_CONSTANT_PROPERTY_DELETED", "constant property {0} cannot be deleted")
}

func (this_ *CheckAccessControls) initFields_CheckAccessControls() {
	_ = this_

	this_.currentClassStack = j2gort.NewExt_java_util_LinkedList()
	this_.constPropertyInits = j2gort.Ext_com_google_common_collect_Tables_NewCustomTable(j2gort.NewExt_java_util_LinkedHashMap(), lambdaSingleton_check_access_controls_1)
}

// Java constructor: CheckAccessControls(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckAccessControls(compiler AbstractCompiler) *CheckAccessControls {
	clinit_CheckAccessControls()
	this_ := &CheckAccessControls{}
	InitCheckAccessControls(this_, compiler)
	return this_
}

func InitCheckAccessControls(this_ *CheckAccessControls, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckAccessControls()
	this_.compiler = compiler
	this_.typeRegistry = compiler.GetTypeRegistry()
}

// Java method: CheckAccessControls.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckAccessControls) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var val_1 CollectFileOverviewVisibility
	InitCollectFileOverviewVisibility(&val_1, this_.compiler)
	var collectPass *CollectFileOverviewVisibility = &val_1
	_ = collectPass
	collectPass.Process_2(externs, root)
	this_.defaultVisibilityForFiles = collectPass.GetFileOverviewVisibilityMap()
	NodeTraversal_Traverse(this_.compiler, externs, this_)
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckAccessControls.enterAccessControlScope(com.google.javascript.rhino.Node root)
func (this_ *CheckAccessControls) priv_CheckAccessControls_enterAccessControlScope_989(root pkg_rhino.Node) {
	_ = root
	var scopeType pkg_rhino.ObjectType = this_.priv_CheckAccessControls_bestInstanceTypeForMethodOrCtor_991(root)
	_ = scopeType
	if CheckAccessControls_IsMarkedDeprecated(root) {
		this_.deprecationDepth = this_.deprecationDepth + 1
	}
	this_.currentClassStack.Push(scopeType)
}

// Java method: CheckAccessControls.exitAccessControlScope(com.google.javascript.rhino.Node root)
func (this_ *CheckAccessControls) priv_CheckAccessControls_exitAccessControlScope_990(root pkg_rhino.Node) {
	_ = root
	if CheckAccessControls_IsMarkedDeprecated(root) {
		this_.deprecationDepth = this_.deprecationDepth - 1
	}
	this_.currentClassStack.Pop()
}

// Java static method: CheckAccessControls.primaryAccessControlScopeRootFor(com.google.javascript.rhino.Node node)
func CheckAccessControls_PrimaryAccessControlScopeRootFor(node pkg_rhino.Node) pkg_rhino.Node {
	clinit_CheckAccessControls()
	_ = node
	if CheckAccessControls_IsExtendsTarget(node) {
		return node.GetParent()
	} else if CheckAccessControls_IsFunctionOrClass(node) {
		return node
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: CheckAccessControls.bestInstanceTypeForMethodOrCtor(com.google.javascript.rhino.Node n)
func (this_ *CheckAccessControls) priv_CheckAccessControls_bestInstanceTypeForMethodOrCtor_991(n pkg_rhino.Node) pkg_rhino.ObjectType {
	_ = n
	if !(CheckAccessControls_IsFunctionOrClass(n)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if NodeUtil_IsFunctionDeclaration(n) || n.IsClass() {
		return CheckAccessControls_InstanceTypeFor(n.GetJSType())
	}
	switch parent.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		return CheckAccessControls_InstanceTypeFor(n.GetJSType())
	case int32(51) /* Token.ASSIGN */ :
		var lValue pkg_rhino.Node = parent.GetFirstChild()
		_ = lValue
		if NodeUtil_IsNormalGet(lValue) {
			var lValueType pkg_rhino.JSType = lValue.GetJSType()
			_ = lValueType
			if (lValueType != nil && !lValueType.IsJavaNil_()) && (lValueType.IsConstructor() || lValueType.IsInterface()) {
				return CheckAccessControls_InstanceTypeFor(lValueType)
			} else if NodeUtil_IsPrototypeProperty(lValue) {
				return CheckAccessControls_InstanceTypeFor(NodeUtil_GetPrototypeClassName(lValue).GetJSType())
			} else {
				return CheckAccessControls_InstanceTypeFor(lValue.GetFirstChild().GetJSType())
			}
		} else {
			return CheckAccessControls_InstanceTypeFor(lValue.GetJSType())
		}
	case int32(100) /* Token.STRING_KEY */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(129) /* Token.COMPUTED_PROP */ :
		var grandparent pkg_rhino.Node = parent.GetParent()
		_ = grandparent
		var greatGrandparent pkg_rhino.Node = grandparent.GetParent()
		_ = greatGrandparent
		if grandparent.IsObjectLit() {
			if grandparent.GetJSType().IsFunctionPrototypeType() {
				return CheckAccessControls_InstanceTypeFor(grandparent.GetJSType())
			} else {
				return nil
			}
		} else if greatGrandparent.IsClass() {
			return CheckAccessControls_InstanceTypeFor(greatGrandparent.GetJSType())
		} else {
			panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.BoxNode(greatGrandparent))))
		}
	default:
		return nil
	}
	panic("unreachable")
}

// Java static method: CheckAccessControls.instanceTypeFor(com.google.javascript.rhino.jstype.JSType type)
func CheckAccessControls_InstanceTypeFor(type_ pkg_rhino.JSType) pkg_rhino.ObjectType {
	clinit_CheckAccessControls()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return nil
	} else if type_.IsUnionType() {
		return nil
	} else if type_.IsInstanceType() || type_.IsUnknownType() {
		return type_.ToMaybeObjectType()
	} else if type_.IsConstructor() || type_.IsInterface() {
		return type_.ToMaybeFunctionType().GetInstanceType()
	} else if type_.IsFunctionType() {
		return nil
	} else if type_.IsFunctionPrototypeType() {
		return CheckAccessControls_InstanceTypeFor(pkg_rhino.ToIface_JSType(type_.ToMaybeObjectType().GetOwnerFunction()))
	}
	return type_.ToMaybeObjectType()
}

// Java method: CheckAccessControls.shouldTraverse(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *CheckAccessControls) ShouldTraverse_3(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = traversal
	_ = node
	_ = parent
	var accessControlRoot pkg_rhino.Node = CheckAccessControls_PrimaryAccessControlScopeRootFor(node)
	_ = accessControlRoot
	if accessControlRoot != 0 {
		this_.priv_CheckAccessControls_enterAccessControlScope_989(accessControlRoot)
	}
	return true
}

// Java method: CheckAccessControls.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *CheckAccessControls) Visit_3(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = node
	_ = parent
	var identifierBehaviour CheckAccessControls_IdentifierBehaviour = CheckAccessControls_IdentifierBehaviour_Select(node)
	_ = identifierBehaviour
	var propRef *CheckAccessControls_PropertyReference = this_.priv_CheckAccessControls_createPropertyReference_1019(node)
	_ = propRef
	this_.priv_CheckAccessControls_checkDeprecation_992(node, propRef, identifierBehaviour, traversal)
	this_.priv_CheckAccessControls_checkVisibility_993(node, propRef, identifierBehaviour, traversal.GetScope())
	this_.priv_CheckAccessControls_checkConstantProperty_1001(propRef, identifierBehaviour)
	this_.priv_CheckAccessControls_checkFinalClassOverrides_1000(node)
	var accessControlRoot pkg_rhino.Node = CheckAccessControls_PrimaryAccessControlScopeRootFor(node)
	_ = accessControlRoot
	if accessControlRoot != 0 {
		this_.priv_CheckAccessControls_exitAccessControlScope_990(accessControlRoot)
	}
}

// Java method: CheckAccessControls.checkDeprecation(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.CheckAccessControls.@org.jspecify.annotations.Nullable PropertyReference propRef, com.google.javascript.jscomp.CheckAccessControls.IdentifierBehaviour identifierBehaviour, com.google.javascript.jscomp.NodeTraversal traversal)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkDeprecation_992(node pkg_rhino.Node, propRef *CheckAccessControls_PropertyReference, identifierBehaviour CheckAccessControls_IdentifierBehaviour, traversal *NodeTraversal) {
	_ = node
	_ = propRef
	_ = identifierBehaviour
	_ = traversal
	switch (identifierBehaviour).Ordinal() {
	case int32(1) /* IdentifierBehaviour.ES5_CLASS_INVOCATION */, int32(3) /* IdentifierBehaviour.ES6_CLASS_INVOCATION */, int32(4) /* IdentifierBehaviour.ES6_CLASS_NAMESPACE */ :
		this_.priv_CheckAccessControls_checkTypeDeprecation_994(traversal, node)
	case int32(0) /* IdentifierBehaviour.NON_CONSTRUCTOR */ :
		this_.priv_CheckAccessControls_checkNameDeprecation_995(traversal, node)
	default:
	}
	if (propRef != nil) && !(func() bool {
		re_2 := identifierBehaviour
		ae_3 := StaticGet_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE()
		if re_2 == 0 {
			j2gort.ThrowNPE()
		}
		return re_2 == ae_3
	}()) {
		this_.priv_CheckAccessControls_checkPropertyDeprecation_996(traversal, propRef)
	}
}

// Java method: CheckAccessControls.checkVisibility(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.CheckAccessControls.@org.jspecify.annotations.Nullable PropertyReference propRef, com.google.javascript.jscomp.CheckAccessControls.IdentifierBehaviour identifierBehaviour, com.google.javascript.jscomp.Scope scope)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkVisibility_993(node pkg_rhino.Node, propRef *CheckAccessControls_PropertyReference, identifierBehaviour CheckAccessControls_IdentifierBehaviour, scope *Scope) {
	_ = node
	_ = propRef
	_ = identifierBehaviour
	_ = scope
	{
		var iifeCond_1 bool
		{

			re_4 := identifierBehaviour
			ae_5 := StaticGet_CheckAccessControls_IdentifierBehaviour_ES6_CLASS_INVOCATION()
			if re_4 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_4 == ae_5
		}
		if iifeCond_1 {
			this_.priv_CheckAccessControls_checkEs6ConstructorInvocationVisibility_1008(node)
		}
	}
	{
		var iifeCond_1 bool
		{

			re_6 := identifierBehaviour
			ae_7 := StaticGet_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE()
			if re_6 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_6 == ae_7
		}
		if !(iifeCond_1) {
			this_.priv_CheckAccessControls_checkNameVisibility_997(scope, node)
		}
	}
	if (propRef != nil) && !(func() bool {
		re_8 := identifierBehaviour
		ae_9 := StaticGet_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE()
		if re_8 == 0 {
			j2gort.ThrowNPE()
		}
		return re_8 == ae_9
	}()) {
		this_.priv_CheckAccessControls_checkPropertyVisibility_1007(propRef)
	}
}

// Java method: CheckAccessControls.checkTypeDeprecation(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkTypeDeprecation_994(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if !(this_.priv_CheckAccessControls_shouldEmitDeprecationWarning_1014(t, n)) {
		return
	}
	var instanceType pkg_rhino.ObjectType = n.GetJSType().ToMaybeFunctionType().GetInstanceType()
	_ = instanceType
	var deprecationInfo string = CheckAccessControls_GetTypeDeprecationInfo(j2gort.AsIface[pkg_rhino.JSType](instanceType))
	_ = deprecationInfo
	if deprecationInfo == j2gort.NullString {
		return
	}
	var message *DiagnosticType
	_ = message
	if j2gort.CheckNotNullString(deprecationInfo) == "" {
		message = StaticGet_CheckAccessControls_DEPRECATED_CLASS()
	} else {
		message = StaticGet_CheckAccessControls_DEPRECATED_CLASS_REASON()
	}
	this_.compiler.Report_1(JSError_Make_3(n, message, []string{instanceType.String(), deprecationInfo}))
}

// Java method: CheckAccessControls.checkNameDeprecation(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkNameDeprecation_995(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if !(n.IsName()) {
		return
	}
	if !(this_.priv_CheckAccessControls_shouldEmitDeprecationWarning_1014(t, n)) {
		return
	}
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString()))
	_ = var_
	var docInfo *pkg_rhino.JSDocInfo
	_ = docInfo
	if var_ == nil {
		docInfo = nil
	} else {
		docInfo = var_.GetJSDocInfo()
	}
	if (docInfo != nil) && docInfo.IsDeprecated() {
		if docInfo.GetDeprecationReason() != j2gort.NullString {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_CheckAccessControls_DEPRECATED_NAME_REASON(), []string{n.GetString(), docInfo.GetDeprecationReason()}))
		} else {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_CheckAccessControls_DEPRECATED_NAME(), []string{n.GetString()}))
		}
	}
}

// Java method: CheckAccessControls.checkPropertyDeprecation(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkPropertyDeprecation_996(t *NodeTraversal, propRef *CheckAccessControls_PropertyReference) {
	_ = t
	_ = propRef
	if !(this_.priv_CheckAccessControls_shouldEmitDeprecationWarning_1015(t, propRef)) {
		return
	}
	if propRef.SourceNode().GetParent().IsNew() {
		return
	}
	var objectType pkg_rhino.ObjectType = CheckAccessControls_CastToObject(j2gort.AsIface[pkg_rhino.JSType](CheckAccessControls_Dereference(j2gort.AsIface[pkg_rhino.JSType](propRef.ReceiverType()))))
	_ = objectType
	var propertyName string = propRef.Name()
	_ = propertyName
	if objectType != nil && !objectType.IsJavaNil_() {
		var deprecationInfo string = CheckAccessControls_GetPropertyDeprecationInfo(objectType, propertyName)
		_ = deprecationInfo
		if deprecationInfo != j2gort.NullString {
			if !(j2gort.CheckNotNullString(deprecationInfo) == "") {
				this_.compiler.Report_1(JSError_Make_3(propRef.SourceNode(), StaticGet_CheckAccessControls_DEPRECATED_PROP_REASON(), []string{propertyName, propRef.GetReadableTypeNameOrDefault(), deprecationInfo}))
			} else {
				this_.compiler.Report_1(JSError_Make_3(propRef.SourceNode(), StaticGet_CheckAccessControls_DEPRECATED_PROP(), []string{propertyName, propRef.GetReadableTypeNameOrDefault()}))
			}
		}
	}
}

// Java method: CheckAccessControls.checkNameVisibility(com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node name)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkNameVisibility_997(scope *Scope, name pkg_rhino.Node) {
	_ = scope
	_ = name
	if !(name.IsName()) {
		return
	}
	var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(name.GetString()))
	_ = var_
	if var_ == nil {
		return
	}
	var v pkg_rhino.JSDocInfo_Visibility = AccessControlUtils_GetEffectiveNameVisibility(name, var_, this_.defaultVisibilityForFiles)
	_ = v
	switch (v).Ordinal() {
	case int32(1) /* Visibility.PACKAGE */ :
		if !(this_.priv_CheckAccessControls_isPackageAccessAllowed_998(var_, name)) {
			this_.compiler.Report_1(JSError_Make_3(name, StaticGet_CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS(), []string{name.GetString(), var_.GetSourceFile().GetName()}))
		}
	case int32(0) /* Visibility.PRIVATE */ :
		if !(CheckAccessControls_IsPrivateAccessAllowed(var_, name)) {
			this_.compiler.Report_1(JSError_Make_3(name, StaticGet_CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS(), []string{name.GetString(), var_.GetSourceFile().GetName()}))
		}
	default:
	}
}

// Java static method: CheckAccessControls.isPrivateAccessAllowed(com.google.javascript.jscomp.Var var, com.google.javascript.rhino.Node name)
func CheckAccessControls_IsPrivateAccessAllowed(var_ *Var, name pkg_rhino.Node) bool {
	clinit_CheckAccessControls()
	_ = var_
	_ = name
	var varSrc pkg_rhino.StaticSourceFile = var_.GetSourceFile()
	_ = varSrc
	var refSrc pkg_rhino.StaticSourceFile = name.GetStaticSourceFile()
	_ = refSrc
	return (((varSrc == nil || varSrc.IsJavaNil_()) || (refSrc == nil || refSrc.IsJavaNil_())) || (varSrc.GetName() == refSrc.GetName()))
}

// Java method: CheckAccessControls.isPackageAccessAllowed(com.google.javascript.jscomp.Var var, com.google.javascript.rhino.Node name)
func (this_ *CheckAccessControls) priv_CheckAccessControls_isPackageAccessAllowed_998(var_ *Var, name pkg_rhino.Node) bool {
	_ = var_
	_ = name
	var varSrc pkg_rhino.StaticSourceFile = var_.GetSourceFile()
	_ = varSrc
	var refSrc pkg_rhino.StaticSourceFile = name.GetStaticSourceFile()
	_ = refSrc
	if (varSrc == nil || varSrc.IsJavaNil_()) && (refSrc == nil || refSrc.IsJavaNil_()) {
		return false
	}
	var codingConvention CodingConvention = this_.compiler.GetCodingConvention()
	_ = codingConvention
	var srcPackage string = codingConvention.GetPackageName(varSrc)
	_ = srcPackage
	var refPackage string = codingConvention.GetPackageName(refSrc)
	_ = refPackage
	return (((srcPackage != j2gort.NullString) && (refPackage != j2gort.NullString)) && (srcPackage == refPackage))
}

// Java method: CheckAccessControls.checkPropertyOverrideVisibilityIsSame(com.google.javascript.rhino.JSDocInfo.Visibility overriding, com.google.javascript.rhino.JSDocInfo.Visibility overridden, com.google.javascript.rhino.JSDocInfo.@org.jspecify.annotations.Nullable Visibility fileOverview, com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkPropertyOverrideVisibilityIsSame_999(overriding pkg_rhino.JSDocInfo_Visibility, overridden pkg_rhino.JSDocInfo_Visibility, fileOverview pkg_rhino.JSDocInfo_Visibility, propRef *CheckAccessControls_PropertyReference) {
	_ = overriding
	_ = overridden
	_ = fileOverview
	_ = propRef
	if (((overriding == pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()) && (overriding != overridden)) && (fileOverview != 0)) && (fileOverview != pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()) {
		this_.compiler.Report_1(JSError_Make_3(propRef.SourceNode(), StaticGet_CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY(), []string{propRef.Name(), (fileOverview).Name()}))
	}
}

// Java static method: CheckAccessControls.getOverridingPropertyVisibility(com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef)
func CheckAccessControls_GetOverridingPropertyVisibility(propRef *CheckAccessControls_PropertyReference) pkg_rhino.JSDocInfo_Visibility {
	clinit_CheckAccessControls()
	_ = propRef
	var overridingInfo *pkg_rhino.JSDocInfo = propRef.GetJSDocInfo()
	_ = overridingInfo
	if (overridingInfo == nil) || !(overridingInfo.IsOverride()) {
		return pkg_rhino.JSDocInfo_Visibility(0)
	} else {
		return overridingInfo.GetVisibility()
	}
}

// Java method: CheckAccessControls.checkFinalClassOverrides(com.google.javascript.rhino.Node ctor)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkFinalClassOverrides_1000(ctor pkg_rhino.Node) {
	_ = ctor
	if !(CheckAccessControls_IsFunctionOrClass(ctor)) || NodeUtil_IsEs6ConstructorMemberFunctionDef(ctor.GetParent()) {
		return
	}
	var ctorType *pkg_rhino.FunctionType = ctor.GetJSType().ToMaybeFunctionType()
	_ = ctorType
	if (ctorType == nil) || !(ctorType.IsConstructor()) {
		return
	}
	var finalParentClass pkg_rhino.ObjectType = CheckAccessControls_GetSuperClassInstanceIfFinal(ctorType)
	_ = finalParentClass
	if finalParentClass != nil && !finalParentClass.IsJavaNil_() {
		this_.compiler.Report_1(JSError_Make_3(ctor, StaticGet_CheckAccessControls_EXTEND_FINAL_CLASS(), []string{ctorType.GetDisplayName(), finalParentClass.GetDisplayName()}))
	}
}

// Java method: CheckAccessControls.checkConstantProperty(com.google.javascript.jscomp.CheckAccessControls.@org.jspecify.annotations.Nullable PropertyReference propRef, com.google.javascript.jscomp.CheckAccessControls.IdentifierBehaviour identifierBehaviour)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkConstantProperty_1001(propRef *CheckAccessControls_PropertyReference, identifierBehaviour CheckAccessControls_IdentifierBehaviour) {
	_ = propRef
	_ = identifierBehaviour
	if (propRef == nil) || func() bool {
		re_10 := identifierBehaviour
		ae_11 := StaticGet_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE()
		if re_10 == 0 {
			j2gort.ThrowNPE()
		}
		return re_10 == ae_11
	}() {
		return
	}
	var objectType pkg_rhino.ObjectType = CheckAccessControls_Dereference(j2gort.AsIface[pkg_rhino.JSType](propRef.ReceiverType()))
	_ = objectType
	var propertyName string = propRef.Name()
	_ = propertyName
	var sourceNode pkg_rhino.Node = propRef.SourceNode()
	_ = sourceNode
	var constness CheckAccessControls_Constancy = this_.priv_CheckAccessControls_isPropertyDeclaredConstant_1017(objectType, propertyName)
	_ = constness
	{
		var iifeCond_1 bool
		{

			re_12 := constness
			ae_13 := StaticGet_CheckAccessControls_Constancy_MUTABLE()
			if re_12 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_12 == ae_13
		}
		if iifeCond_1 {
			return
		}
	}
	if sourceNode.IsFromExterns() && propRef.Declaration() {
		this_.priv_CheckAccessControls_recordConstPropertyInit_1003(propRef, objectType, constness)
		return
	}
	if !(propRef.Mutation()) {
		return
	}
	if propRef.IsDeletion() {
		this_.compiler.Report_1(JSError_Make_3(sourceNode, StaticGet_CheckAccessControls_CONST_PROPERTY_DELETED(), []string{propertyName}))
		return
	}
	if objectType.IsFunctionType() && !(objectType.ToMaybeFunctionType().IsConstructor()) {
		return
	}
	if objectType.IsStructuralType() && !(propRef.Declaration()) {
		this_.compiler.Report_1(JSError_Make_3(sourceNode, StaticGet_CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE(), []string{propertyName, "unknown location due to structural typing"}))
		return
	}
	var init *CheckAccessControls_ConstantDeclaration = this_.priv_CheckAccessControls_getConstPropertyInit_1002(propRef, objectType)
	_ = init
	if init != nil {
		var diagnostic *DiagnosticType
		_ = diagnostic
		{
			var iifeCond_1 bool
			{

				re_14 := init.annotation
				ae_15 := StaticGet_CheckAccessControls_Constancy_FINAL()
				if re_14 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_14 == ae_15
			}
			if iifeCond_1 {
				diagnostic = StaticGet_CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN()
			} else {
				diagnostic = StaticGet_CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE()
			}
		}
		this_.compiler.Report_1(JSError_Make_3(sourceNode, diagnostic, []string{propertyName, init.node.GetLocation()}))
	}
	this_.priv_CheckAccessControls_recordConstPropertyInit_1003(propRef, objectType, constness)
}

// Java method: CheckAccessControls.getConstPropertyInit(com.google.javascript.jscomp.CheckAccessControls.PropertyReference ref, com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *CheckAccessControls) priv_CheckAccessControls_getConstPropertyInit_1002(ref *CheckAccessControls_PropertyReference, type_ pkg_rhino.ObjectType) *CheckAccessControls_ConstantDeclaration {
	_ = ref
	_ = type_
	var name string = ref.Name()
	_ = name
	for type_ != nil && !type_.IsJavaNil_() {
		var init *CheckAccessControls_ConstantDeclaration = j2gort.CastPtr[CheckAccessControls_ConstantDeclaration](this_.constPropertyInits.Get(type_, j2gort.BoxString(name)))
		_ = init
		if init != nil {
			return init
		}
		var canonicalInit *CheckAccessControls_ConstantDeclaration = j2gort.CastPtr[CheckAccessControls_ConstantDeclaration](this_.constPropertyInits.Get(CheckAccessControls_GetCanonicalInstance(type_), j2gort.BoxString(name)))
		_ = canonicalInit
		if canonicalInit != nil {
			return canonicalInit
		}
		type_ = type_.GetImplicitPrototype()
	}
	return nil
}

// Java method: CheckAccessControls.recordConstPropertyInit(com.google.javascript.jscomp.CheckAccessControls.PropertyReference ref, com.google.javascript.rhino.jstype.ObjectType type, com.google.javascript.jscomp.CheckAccessControls.Constancy annotation)
func (this_ *CheckAccessControls) priv_CheckAccessControls_recordConstPropertyInit_1003(ref *CheckAccessControls_PropertyReference, type_ pkg_rhino.ObjectType, annotation CheckAccessControls_Constancy) {
	_ = ref
	_ = type_
	_ = annotation
	this_.constPropertyInits.Row(type_).PutIfAbsent(j2gort.BoxString(ref.Name()), NewCheckAccessControls_ConstantDeclaration(ref.SourceNode(), annotation))
	if type_.IsInstanceType() {
		var prototype pkg_rhino.ObjectType = type_.GetImplicitPrototype()
		_ = prototype
		if (prototype != nil && !prototype.IsJavaNil_()) && prototype.HasProperty_1_string(ref.Name()) {
			this_.constPropertyInits.Row(prototype).PutIfAbsent(j2gort.BoxString(ref.Name()), NewCheckAccessControls_ConstantDeclaration(ref.SourceNode(), annotation))
		}
	}
}

// Java static method: CheckAccessControls.getCanonicalInstance(com.google.javascript.rhino.jstype.ObjectType obj)
func CheckAccessControls_GetCanonicalInstance(obj pkg_rhino.ObjectType) pkg_rhino.ObjectType {
	clinit_CheckAccessControls()
	_ = obj
	var ctor *pkg_rhino.FunctionType = obj.GetConstructor()
	_ = ctor
	if ctor == nil {
		return obj
	} else {
		return ctor.GetInstanceType()
	}
}

// Java static method: CheckAccessControls.dereference(com.google.javascript.rhino.jstype.JSType type)
func CheckAccessControls_Dereference(type_ pkg_rhino.JSType) pkg_rhino.ObjectType {
	clinit_CheckAccessControls()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return nil
	} else {
		return type_.Dereference()
	}
}

// Java method: CheckAccessControls.typeOrUnknown(com.google.javascript.rhino.jstype.JSType type)
func (this_ *CheckAccessControls) priv_CheckAccessControls_typeOrUnknown_1004(type_ pkg_rhino.JSType) pkg_rhino.JSType {
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	} else {
		return type_
	}
}

// Java method: CheckAccessControls.typeOrUnknown(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *CheckAccessControls) priv_CheckAccessControls_typeOrUnknown_1005(type_ pkg_rhino.ObjectType) pkg_rhino.ObjectType {
	_ = type_
	return j2gort.AsIface[pkg_rhino.ObjectType](this_.priv_CheckAccessControls_typeOrUnknown_1004(j2gort.AsIface[pkg_rhino.JSType](type_)))
}

// Java method: CheckAccessControls.boxedOrUnknown(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type)
func (this_ *CheckAccessControls) priv_CheckAccessControls_boxedOrUnknown_1006(type_ pkg_rhino.JSType) pkg_rhino.ObjectType {
	_ = type_
	return this_.priv_CheckAccessControls_typeOrUnknown_1005(CheckAccessControls_Dereference(type_))
}

// Java method: CheckAccessControls.checkPropertyVisibility(com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkPropertyVisibility_1007(propRef *CheckAccessControls_PropertyReference) {
	_ = propRef
	if NodeUtil_IsEs6ConstructorMemberFunctionDef(propRef.SourceNode()) {
		return
	}
	var rawReferenceType pkg_rhino.JSType = this_.priv_CheckAccessControls_typeOrUnknown_1005(propRef.ReceiverType()).Autobox()
	_ = rawReferenceType
	var referenceType pkg_rhino.ObjectType = CheckAccessControls_CastToObject(rawReferenceType)
	_ = referenceType
	var propertyName string = propRef.Name()
	_ = propertyName
	var definingSource pkg_rhino.StaticSourceFile = AccessControlUtils_GetDefiningSource(propRef.SourceNode(), referenceType, propertyName)
	_ = definingSource
	var isOverride bool = (propRef.IsDocumentedDeclaration() || propRef.Override())
	_ = isOverride
	var objectType pkg_rhino.ObjectType = AccessControlUtils_GetObjectType(referenceType, isOverride, propertyName)
	_ = objectType
	var fileOverviewVisibility pkg_rhino.JSDocInfo_Visibility = j2gort.Cast[pkg_rhino.JSDocInfo_Visibility](this_.defaultVisibilityForFiles.Get(definingSource))
	_ = fileOverviewVisibility
	var visibility pkg_rhino.JSDocInfo_Visibility = CheckAccessControls_GetEffectivePropertyVisibility(propRef, referenceType, this_.defaultVisibilityForFiles)
	_ = visibility
	if isOverride {
		var overriding pkg_rhino.JSDocInfo_Visibility = CheckAccessControls_GetOverridingPropertyVisibility(propRef)
		_ = overriding
		if overriding != 0 {
			this_.priv_CheckAccessControls_checkPropertyOverrideVisibilityIsSame_999(overriding, visibility, fileOverviewVisibility, propRef)
		}
	}
	var reportType pkg_rhino.JSType = rawReferenceType
	_ = reportType
	if objectType != nil && !objectType.IsJavaNil_() {
		var node pkg_rhino.Node = objectType.GetOwnPropertyDefSite_1_string(propertyName)
		_ = node
		if node == 0 {
			return
		}
		reportType = j2gort.AsIface[pkg_rhino.JSType](objectType)
		definingSource = node.GetStaticSourceFile()
	} else if fileOverviewVisibility == 0 {
		return
	}
	var referenceSource pkg_rhino.StaticSourceFile = propRef.SourceNode().GetStaticSourceFile()
	_ = referenceSource
	if isOverride {
		var sameInput bool = ((referenceSource != nil && !referenceSource.IsJavaNil_()) && func() bool {
			sr_16 := referenceSource.GetName()
			se_17 := definingSource.GetName()
			return j2gort.CheckNotNullString(sr_16) == se_17
		}())
		_ = sameInput
		this_.priv_CheckAccessControls_checkPropertyOverrideVisibility_1009(propRef, visibility, fileOverviewVisibility, reportType, sameInput)
	} else {
		this_.priv_CheckAccessControls_checkPropertyAccessVisibility_1010(propRef, visibility, reportType, referenceSource, definingSource)
	}
}

// Java method: CheckAccessControls.checkEs6ConstructorInvocationVisibility(com.google.javascript.rhino.Node target)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkEs6ConstructorInvocationVisibility_1008(target pkg_rhino.Node) {
	_ = target
	var ctorType *pkg_rhino.FunctionType = target.GetJSType().ToMaybeFunctionType()
	_ = ctorType
	var prototypeType pkg_rhino.ObjectType = ctorType.GetPrototype()
	_ = prototypeType
	var classDefinition pkg_rhino.Node = ctorType.GetSource()
	_ = classDefinition
	var definingSource pkg_rhino.StaticSourceFile
	_ = definingSource
	if classDefinition == 0 {
		definingSource = nil
	} else {
		definingSource = AccessControlUtils_GetDefiningSource(classDefinition, prototypeType, "constructor")
	}
	var fauxCtorRef *CheckAccessControls_PropertyReference = CheckAccessControls_PropertyReference_S_Builder().SetSourceNode(target).SetName("constructor").SetReceiverType(prototypeType).SetMutation(false).SetDeclaration(false).SetOverride(false).SetReadableTypeName(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.BoxString(ctorType.GetInstanceType().String())
	})).Build()
	_ = fauxCtorRef
	var annotatedCtorVisibility pkg_rhino.JSDocInfo_Visibility = CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty(fauxCtorRef, prototypeType, j2gort.Cast[pkg_rhino.JSDocInfo_Visibility](this_.defaultVisibilityForFiles.Get(definingSource)))
	_ = annotatedCtorVisibility
	var effectiveCtorVisibility pkg_rhino.JSDocInfo_Visibility
	_ = effectiveCtorVisibility
	{
		var iifeCond_1 bool
		{

			re_18 := annotatedCtorVisibility
			ae_19 := pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()
			if re_18 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_18 == ae_19
		}
		if iifeCond_1 {
			effectiveCtorVisibility = pkg_rhino.StaticGet_JSDocInfo_Visibility_PUBLIC()
		} else {
			effectiveCtorVisibility = annotatedCtorVisibility
		}
	}
	this_.priv_CheckAccessControls_checkPropertyAccessVisibility_1010(fauxCtorRef, effectiveCtorVisibility, pkg_rhino.ToIface_JSType(ctorType), target.GetStaticSourceFile(), definingSource)
}

// Java method: CheckAccessControls.checkPropertyOverrideVisibility(com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef, com.google.javascript.rhino.JSDocInfo.Visibility visibility, com.google.javascript.rhino.JSDocInfo.Visibility fileOverviewVisibility, com.google.javascript.rhino.jstype.JSType objectType, boolean sameInput)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkPropertyOverrideVisibility_1009(propRef *CheckAccessControls_PropertyReference, visibility pkg_rhino.JSDocInfo_Visibility, fileOverviewVisibility pkg_rhino.JSDocInfo_Visibility, objectType pkg_rhino.JSType, sameInput bool) {
	_ = propRef
	_ = visibility
	_ = fileOverviewVisibility
	_ = objectType
	_ = sameInput
	var overridingVisibility pkg_rhino.JSDocInfo_Visibility
	_ = overridingVisibility
	if propRef.Override() {
		overridingVisibility = propRef.GetJSDocInfo().GetVisibility()
	} else {
		overridingVisibility = pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()
	}
	if (visibility == pkg_rhino.StaticGet_JSDocInfo_Visibility_PRIVATE()) && !(sameInput) {
		this_.compiler.Report_1(JSError_Make_3(propRef.SourceNode(), StaticGet_CheckAccessControls_PRIVATE_OVERRIDE(), []string{objectType.String()}))
	} else if !(CheckAccessControls_CanOverrideVisibility(visibility, overridingVisibility)) && (fileOverviewVisibility == 0) {
		this_.compiler.Report_1(JSError_Make_3(propRef.SourceNode(), StaticGet_CheckAccessControls_VISIBILITY_MISMATCH(), []string{(visibility).Name(), objectType.String(), (overridingVisibility).Name()}))
	}
}

// Java static method: CheckAccessControls.canOverrideVisibility(com.google.javascript.rhino.JSDocInfo.Visibility superclassVisibility, com.google.javascript.rhino.JSDocInfo.Visibility subclassVisibility)
func CheckAccessControls_CanOverrideVisibility(superclassVisibility pkg_rhino.JSDocInfo_Visibility, subclassVisibility pkg_rhino.JSDocInfo_Visibility) bool {
	clinit_CheckAccessControls()
	_ = superclassVisibility
	_ = subclassVisibility
	return (((superclassVisibility).Ordinal() - (subclassVisibility).Ordinal()) <= int32(0))
}

// Java method: CheckAccessControls.checkPropertyAccessVisibility(com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef, com.google.javascript.rhino.JSDocInfo.Visibility visibility, com.google.javascript.rhino.jstype.JSType objectType, com.google.javascript.rhino.StaticSourceFile referenceSource, com.google.javascript.rhino.StaticSourceFile definingSource)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkPropertyAccessVisibility_1010(propRef *CheckAccessControls_PropertyReference, visibility pkg_rhino.JSDocInfo_Visibility, objectType pkg_rhino.JSType, referenceSource pkg_rhino.StaticSourceFile, definingSource pkg_rhino.StaticSourceFile) {
	_ = propRef
	_ = visibility
	_ = objectType
	_ = referenceSource
	_ = definingSource
	if ((referenceSource != nil && !referenceSource.IsJavaNil_()) && (definingSource != nil && !definingSource.IsJavaNil_())) && func() bool {
		sr_20 := referenceSource.GetName()
		se_21 := definingSource.GetName()
		return j2gort.CheckNotNullString(sr_20) == se_21
	}() {
		return
	}
	var ownerType pkg_rhino.ObjectType = CheckAccessControls_InstanceTypeFor(objectType)
	_ = ownerType
	switch (visibility).Ordinal() {
	case int32(1) /* Visibility.PACKAGE */ :
		this_.priv_CheckAccessControls_checkPackagePropertyVisibility_1011(propRef, referenceSource, definingSource)
	case int32(0) /* Visibility.PRIVATE */ :
		this_.priv_CheckAccessControls_checkPrivatePropertyVisibility_1012(propRef, ownerType)
	case int32(2) /* Visibility.PROTECTED */ :
		this_.priv_CheckAccessControls_checkProtectedPropertyVisibility_1013(propRef, ownerType)
	default:
	}
}

// Java method: CheckAccessControls.checkPackagePropertyVisibility(com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef, com.google.javascript.rhino.StaticSourceFile referenceSource, com.google.javascript.rhino.StaticSourceFile definingSource)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkPackagePropertyVisibility_1011(propRef *CheckAccessControls_PropertyReference, referenceSource pkg_rhino.StaticSourceFile, definingSource pkg_rhino.StaticSourceFile) {
	_ = propRef
	_ = referenceSource
	_ = definingSource
	var codingConvention CodingConvention = this_.compiler.GetCodingConvention()
	_ = codingConvention
	var refPackage string = codingConvention.GetPackageName(referenceSource)
	_ = refPackage
	var defPackage string = codingConvention.GetPackageName(definingSource)
	_ = defPackage
	if ((refPackage == j2gort.NullString) || (defPackage == j2gort.NullString)) || !(j2gort.CheckNotNullString(refPackage) == defPackage) {
		this_.compiler.Report_1(JSError_Make_3(propRef.SourceNode(), StaticGet_CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS(), []string{propRef.Name(), propRef.GetReadableTypeNameOrDefault()}))
	}
}

// Java method: CheckAccessControls.checkPrivatePropertyVisibility(com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType ownerType)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkPrivatePropertyVisibility_1012(propRef *CheckAccessControls_PropertyReference, ownerType pkg_rhino.ObjectType) {
	_ = propRef
	_ = ownerType
	var readableTypeName string
	_ = readableTypeName
	if (ownerType == nil || ownerType.IsJavaNil_()) || ownerType.Equals(propRef.ReceiverType()) {
		readableTypeName = propRef.GetReadableTypeNameOrDefault()
	} else {
		readableTypeName = ownerType.String()
	}
	this_.compiler.Report_1(JSError_Make_3(propRef.SourceNode(), StaticGet_CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS(), []string{propRef.Name(), readableTypeName}))
}

// Java method: CheckAccessControls.checkProtectedPropertyVisibility(com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType ownerType)
func (this_ *CheckAccessControls) priv_CheckAccessControls_checkProtectedPropertyVisibility_1013(propRef *CheckAccessControls_PropertyReference, ownerType pkg_rhino.ObjectType) {
	_ = propRef
	_ = ownerType
	if ownerType != nil && !ownerType.IsJavaNil_() {
		{
			cur_22 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.currentClassStack))
			for cur_22.HasNext() {
				rawElem_23 := cur_22.Next()
				scopeType := j2gort.AsIface[pkg_rhino.JSType](rawElem_23)
				_ = scopeType
				if scopeType == nil || scopeType.IsJavaNil_() {
					continue
				} else if scopeType.IsSubtypeOf_1(j2gort.AsIface[pkg_rhino.JSType](ownerType)) {
					return
				}
			}
		}
	}
	this_.compiler.Report_1(JSError_Make_3(propRef.SourceNode(), StaticGet_CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS(), []string{propRef.Name(), propRef.GetReadableTypeNameOrDefault()}))
}

// Java method: CheckAccessControls.shouldEmitDeprecationWarning(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckAccessControls) priv_CheckAccessControls_shouldEmitDeprecationWarning_1014(t *NodeTraversal, n pkg_rhino.Node) bool {
	_ = t
	_ = n
	if t.InGlobalScope() {
		if !(NodeUtil_IsInvocationTarget(n)) && !(n.IsNew()) {
			return false
		}
	}
	return !(this_.priv_CheckAccessControls_canAccessDeprecatedTypes_1016(t))
}

// Java method: CheckAccessControls.shouldEmitDeprecationWarning(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef)
func (this_ *CheckAccessControls) priv_CheckAccessControls_shouldEmitDeprecationWarning_1015(t *NodeTraversal, propRef *CheckAccessControls_PropertyReference) bool {
	_ = t
	_ = propRef
	if t.InGlobalScope() {
		if !(NodeUtil_IsInvocationTarget(propRef.SourceNode())) {
			return false
		}
	}
	if propRef.Mutation() {
		return false
	}
	var jsdoc *pkg_rhino.JSDocInfo = propRef.GetJSDocInfo()
	_ = jsdoc
	if (propRef.Declaration() && (jsdoc != nil)) && jsdoc.IsDeprecated() {
		return false
	}
	return !(this_.priv_CheckAccessControls_canAccessDeprecatedTypes_1016(t))
}

// Java method: CheckAccessControls.canAccessDeprecatedTypes(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CheckAccessControls) priv_CheckAccessControls_canAccessDeprecatedTypes_1016(t *NodeTraversal) bool {
	_ = t
	var scopeRoot pkg_rhino.Node = t.GetClosestHoistScopeRoot()
	_ = scopeRoot
	if NodeUtil_IsFunctionBlock(scopeRoot) {
		scopeRoot = scopeRoot.GetParent()
	}
	var scopeRootParent pkg_rhino.Node = scopeRoot.GetParent()
	_ = scopeRootParent
	return (((this_.deprecationDepth > int32(0)) || (CheckAccessControls_GetTypeDeprecationInfo(this_.priv_CheckAccessControls_getTypeOfThis_1018(scopeRoot)) != j2gort.NullString)) || (((scopeRootParent != 0) && scopeRootParent.IsAssign()) && (CheckAccessControls_GetTypeDeprecationInfo(j2gort.AsIface[pkg_rhino.JSType](this_.priv_CheckAccessControls_bestInstanceTypeForMethodOrCtor_991(scopeRoot))) != j2gort.NullString)))
}

// Java static method: CheckAccessControls.isMarkedDeprecated(com.google.javascript.rhino.Node n)
func CheckAccessControls_IsMarkedDeprecated(n pkg_rhino.Node) bool {
	clinit_CheckAccessControls()
	_ = n
	return (CheckAccessControls_GetDeprecationReason(NodeUtil_GetBestJSDocInfo(n)) != j2gort.NullString)
}

// Java static method: CheckAccessControls.getTypeDeprecationInfo(com.google.javascript.rhino.jstype.JSType type)
func CheckAccessControls_GetTypeDeprecationInfo(type_ pkg_rhino.JSType) string {
	clinit_CheckAccessControls()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return j2gort.NullString
	}
	return CheckAccessControls_GetDeprecationReason(type_.GetJSDocInfo())
}

// Java static method: CheckAccessControls.getDeprecationReason(com.google.javascript.rhino.JSDocInfo info)
func CheckAccessControls_GetDeprecationReason(info *pkg_rhino.JSDocInfo) string {
	clinit_CheckAccessControls()
	_ = info
	if (info != nil) && info.IsDeprecated() {
		if info.GetDeprecationReason() != j2gort.NullString {
			return info.GetDeprecationReason()
		}
		return ""
	}
	return j2gort.NullString
}

// Java method: CheckAccessControls.isPropertyDeclaredConstant(com.google.javascript.rhino.jstype.ObjectType objectType, java.lang.String prop)
func (this_ *CheckAccessControls) priv_CheckAccessControls_isPropertyDeclaredConstant_1017(objectType pkg_rhino.ObjectType, prop string) CheckAccessControls_Constancy {
	_ = objectType
	_ = prop
	for ; objectType != nil && !objectType.IsJavaNil_(); objectType = objectType.GetImplicitPrototype() {
		var docInfo *pkg_rhino.JSDocInfo = objectType.GetOwnPropertyJSDocInfo_1_string(prop)
		_ = docInfo
		if docInfo == nil {
			continue
		} else if docInfo.IsFinal() {
			return StaticGet_CheckAccessControls_Constancy_FINAL()
		} else if docInfo.IsConstant() {
			return StaticGet_CheckAccessControls_Constancy_OTHER_CONSTANT()
		}
	}
	return StaticGet_CheckAccessControls_Constancy_MUTABLE()
}

// Java static method: CheckAccessControls.getPropertyDeprecationInfo(com.google.javascript.rhino.jstype.ObjectType type, java.lang.String prop)
func CheckAccessControls_GetPropertyDeprecationInfo(type_ pkg_rhino.ObjectType, prop string) string {
	clinit_CheckAccessControls()
	_ = type_
	_ = prop
	var depReason string = CheckAccessControls_GetDeprecationReason(type_.GetOwnPropertyJSDocInfo_1_string(prop))
	_ = depReason
	if depReason != j2gort.NullString {
		return depReason
	}
	var implicitProto pkg_rhino.ObjectType = type_.GetImplicitPrototype()
	_ = implicitProto
	if implicitProto != nil && !implicitProto.IsJavaNil_() {
		return CheckAccessControls_GetPropertyDeprecationInfo(implicitProto, prop)
	}
	return j2gort.NullString
}

// Java static method: CheckAccessControls.getSuperClassInstanceIfFinal(com.google.javascript.rhino.jstype.FunctionType subCtor)
func CheckAccessControls_GetSuperClassInstanceIfFinal(subCtor *pkg_rhino.FunctionType) pkg_rhino.ObjectType {
	clinit_CheckAccessControls()
	_ = subCtor
	var ctor *pkg_rhino.FunctionType = subCtor.GetSuperClassConstructor()
	_ = ctor
	var doc *pkg_rhino.JSDocInfo
	_ = doc
	if ctor == nil {
		doc = nil
	} else {
		doc = ctor.GetJSDocInfo()
	}
	if (doc != nil) && doc.IsFinal() {
		return ctor.GetInstanceType()
	}
	return nil
}

// Java static method: CheckAccessControls.castToObject(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type)
func CheckAccessControls_CastToObject(type_ pkg_rhino.JSType) pkg_rhino.ObjectType {
	clinit_CheckAccessControls()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return nil
	} else {
		return type_.ToMaybeObjectType()
	}
}

// Java static method: CheckAccessControls.isFunctionOrClass(com.google.javascript.rhino.Node n)
func CheckAccessControls_IsFunctionOrClass(n pkg_rhino.Node) bool {
	clinit_CheckAccessControls()
	_ = n
	return (n.IsFunction() || n.IsClass())
}

// Java static method: CheckAccessControls.isExtendsTarget(com.google.javascript.rhino.Node node)
func CheckAccessControls_IsExtendsTarget(node pkg_rhino.Node) bool {
	clinit_CheckAccessControls()
	_ = node
	var parent pkg_rhino.Node = node.GetParent()
	_ = parent
	return (parent.IsClass() && node.IsSecondChildOf(parent))
}

// Java method: CheckAccessControls.getTypeOfThis(com.google.javascript.rhino.Node scopeRoot)
func (this_ *CheckAccessControls) priv_CheckAccessControls_getTypeOfThis_1018(scopeRoot pkg_rhino.Node) pkg_rhino.JSType {
	_ = scopeRoot
	if scopeRoot.IsRoot() || scopeRoot.IsScript() {
		return j2gort.AsIface[pkg_rhino.JSType](CheckAccessControls_CastToObject(scopeRoot.GetJSType()))
	} else if scopeRoot.IsModuleBody() {
		return nil
	} else if NodeUtil_IsClassStaticBlock(scopeRoot) {
		var classNode pkg_rhino.Node = NodeUtil_GetEnclosingClass(scopeRoot)
		_ = classNode
		return classNode.GetJSType()
	}
	if !(scopeRoot.IsFunction()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(scopeRoot))
	}
	var nodeType pkg_rhino.JSType = scopeRoot.GetJSType()
	_ = nodeType
	if (nodeType != nil && !nodeType.IsJavaNil_()) && nodeType.IsFunctionType() {
		return nodeType.ToMaybeFunctionType().GetTypeOfThis()
	} else {
		return nil
	}
}

// Java method: CheckAccessControls.createPropertyReference(com.google.javascript.rhino.Node sourceNode)
func (this_ *CheckAccessControls) priv_CheckAccessControls_createPropertyReference_1019(sourceNode pkg_rhino.Node) *CheckAccessControls_PropertyReference {
	_ = sourceNode
	var parent pkg_rhino.Node = sourceNode.GetParent()
	_ = parent
	var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(sourceNode)
	_ = jsdoc
	var builder CheckAccessControls_PropertyReference_Builder = CheckAccessControls_PropertyReference_S_Builder()
	_ = builder
	switch sourceNode.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */ :
		{
			var isLValue bool = NodeUtil_IsLValue(sourceNode)
			_ = isLValue
			builder.SetName(sourceNode.GetString()).SetReceiverType(this_.priv_CheckAccessControls_boxedOrUnknown_1006(sourceNode.GetFirstChild().GetJSType())).SetMutation((isLValue || sourceNode.GetParent().IsDelProp())).SetDeclaration((parent.IsExprResult() || (((jsdoc != nil) && jsdoc.IsConstant()) && isLValue))).SetOverride(((jsdoc != nil) && isLValue)).SetReadableTypeName(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
				return j2gort.BoxString(this_.typeRegistry.GetReadableTypeName(sourceNode.GetFirstChild()))
			}))
		}
	case int32(100) /* Token.STRING_KEY */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */ :
		{
			{
				var ctorType pkg_rhino.JSType = nil
				_ = ctorType
				switch parent.TokenSwitchOrd_() {
				case int32(47) /* Token.OBJECTLIT */ :
					if !(parent.GetJSType().IsLiteralObject()) {
						return nil
					}
					builder.SetName(sourceNode.GetString()).SetReceiverType(this_.priv_CheckAccessControls_typeOrUnknown_1005(pkg_rhino.ObjectType_Cast(parent.GetJSType()))).SetMutation(true).SetDeclaration(true).SetOverride(false).SetReadableTypeName(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
						return j2gort.BoxString(this_.typeRegistry.GetReadableTypeName(parent))
					}))
				case int32(103) /* Token.OBJECT_PATTERN */ :
					builder.SetName(sourceNode.GetString()).SetReceiverType(this_.priv_CheckAccessControls_typeOrUnknown_1005(pkg_rhino.ObjectType_Cast(parent.GetJSType()))).SetMutation(false).SetDeclaration(false).SetOverride(false).SetReadableTypeName(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
						return j2gort.BoxString(this_.typeRegistry.GetReadableTypeName(parent))
					}))
				case int32(106) /* Token.CLASS_MEMBERS */ :
					builder.SetName(sourceNode.GetString()).SetReceiverType(this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())).SetMutation(!(sourceNode.IsMemberFieldDef() && !(sourceNode.HasChildren()))).SetDeclaration(true).SetOverride((jsdoc != nil)).SetReadableTypeName(lambdaSingleton_check_access_controls_2)
					ctorType = parent.GetParent().GetJSType()
					if (ctorType != nil && !ctorType.IsJavaNil_()) && ctorType.IsFunctionType() {
						var ctorFunctionType *pkg_rhino.FunctionType = ctorType.ToMaybeFunctionType()
						_ = ctorFunctionType
						var owningType pkg_rhino.ObjectType = nil
						_ = owningType
						if sourceNode.IsStaticMember() {
							owningType = pkg_rhino.ToIface_ObjectType(ctorFunctionType)
						} else if sourceNode.IsMemberFieldDef() {
							owningType = ctorFunctionType.GetInstanceType()
						} else {
							owningType = ctorFunctionType.GetPrototype()
						}
						builder.SetReceiverType(owningType)
					}
				default:
					panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
				}
			}
		}
	default:
		return nil
	}
	return builder.SetSourceNode(sourceNode).Build()
}

// Java static method: CheckAccessControls.getEffectivePropertyVisibility(com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef, com.google.javascript.rhino.jstype.ObjectType referenceType, com.google.common.collect.ImmutableMap<com.google.javascript.rhino.StaticSourceFile,com.google.javascript.rhino.JSDocInfo.Visibility> fileVisibilityMap)
func CheckAccessControls_GetEffectivePropertyVisibility(propRef *CheckAccessControls_PropertyReference, referenceType pkg_rhino.ObjectType, fileVisibilityMap j2gort.Ext_com_google_common_collect_ImmutableMap) pkg_rhino.JSDocInfo_Visibility {
	clinit_CheckAccessControls()
	_ = propRef
	_ = referenceType
	_ = fileVisibilityMap
	var propertyName string = propRef.Name()
	_ = propertyName
	var isOverride bool = propRef.Override()
	_ = isOverride
	var definingSource pkg_rhino.StaticSourceFile = AccessControlUtils_GetDefiningSource(propRef.SourceNode(), referenceType, propertyName)
	_ = definingSource
	var fileOverviewVisibility pkg_rhino.JSDocInfo_Visibility = j2gort.Cast[pkg_rhino.JSDocInfo_Visibility](fileVisibilityMap.Get(definingSource))
	_ = fileOverviewVisibility
	var objectType pkg_rhino.ObjectType = AccessControlUtils_GetObjectType(referenceType, isOverride, propertyName)
	_ = objectType
	if isOverride {
		var overridden pkg_rhino.JSDocInfo_Visibility = AccessControlUtils_GetOverriddenPropertyVisibility(objectType, propertyName)
		_ = overridden
		return AccessControlUtils_GetEffectiveVisibilityForOverriddenProperty(overridden, fileOverviewVisibility, propertyName)
	} else {
		return CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty(propRef, objectType, fileOverviewVisibility)
	}
}

// Java static method: CheckAccessControls.getEffectiveVisibilityForNonOverriddenProperty(com.google.javascript.jscomp.CheckAccessControls.PropertyReference propRef, com.google.javascript.rhino.jstype.ObjectType objectType, com.google.javascript.rhino.JSDocInfo.@org.jspecify.annotations.Nullable Visibility fileOverviewVisibility)
func CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty(propRef *CheckAccessControls_PropertyReference, objectType pkg_rhino.ObjectType, fileOverviewVisibility pkg_rhino.JSDocInfo_Visibility) pkg_rhino.JSDocInfo_Visibility {
	clinit_CheckAccessControls()
	_ = propRef
	_ = objectType
	_ = fileOverviewVisibility
	var propertyName string = propRef.Name()
	_ = propertyName
	var raw pkg_rhino.JSDocInfo_Visibility = pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()
	_ = raw
	if objectType != nil && !objectType.IsJavaNil_() {
		var jsdoc *pkg_rhino.JSDocInfo = objectType.GetOwnPropertyJSDocInfo_1_string(propertyName)
		_ = jsdoc
		if jsdoc != nil {
			raw = jsdoc.GetVisibility()
		}
	}
	var type_ pkg_rhino.JSType = propRef.GetJSType()
	_ = type_
	var createdFromGoogProvide bool = ((type_ != nil && !type_.IsJavaNil_()) && type_.IsLiteralObject())
	_ = createdFromGoogProvide
	if ((raw != pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()) || (fileOverviewVisibility == 0)) || createdFromGoogProvide {
		return raw
	} else {
		return fileOverviewVisibility
	}
}

// Java enum: com.google.javascript.jscomp.CheckAccessControls.Constancy
type CheckAccessControls_Constancy uint8

const (
	enum_CheckAccessControls_Constancy_FINAL          CheckAccessControls_Constancy = 1
	enum_CheckAccessControls_Constancy_OTHER_CONSTANT CheckAccessControls_Constancy = 2
	enum_CheckAccessControls_Constancy_MUTABLE        CheckAccessControls_Constancy = 3
)

var enumBaseTable_CheckAccessControls_Constancy = [...]j2gort.EnumBase{
	{Name_: "FINAL", Ordinal_: 0},
	{Name_: "OTHER_CONSTANT", Ordinal_: 1},
	{Name_: "MUTABLE", Ordinal_: 2},
}
var enumValues_CheckAccessControls_Constancy = []CheckAccessControls_Constancy{1, 2, 3}

func StaticGet_CheckAccessControls_Constancy_FINAL() CheckAccessControls_Constancy {
	return enum_CheckAccessControls_Constancy_FINAL
}
func StaticGet_CheckAccessControls_Constancy_OTHER_CONSTANT() CheckAccessControls_Constancy {
	return enum_CheckAccessControls_Constancy_OTHER_CONSTANT
}
func StaticGet_CheckAccessControls_Constancy_MUTABLE() CheckAccessControls_Constancy {
	return enum_CheckAccessControls_Constancy_MUTABLE
}

func CheckAccessControls_Constancy_Values() []CheckAccessControls_Constancy {
	return j2gort.CloneSlice(enumValues_CheckAccessControls_Constancy)
}

func CheckAccessControls_Constancy_ValueOf(name string) CheckAccessControls_Constancy {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CheckAccessControls_Constancy {
		if enumBaseTable_CheckAccessControls_Constancy[i_].Name_ == name {
			return CheckAccessControls_Constancy(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CheckAccessControls.Constancy." + name))
}

func (this_ CheckAccessControls_Constancy) IsJavaNil_() bool                              { return this_ == 0 }
func (CheckAccessControls_Constancy) HasIdentityEquals_()                                 {}
func (CheckAccessControls_Constancy) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CheckAccessControls_Constancy) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CheckAccessControls_Constancy) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CheckAccessControls_Constancy) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckAccessControls$Constancy", "Constancy"
}
func (this_ CheckAccessControls_Constancy) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CheckAccessControls_Constancy[this_-1]
}
func (this_ CheckAccessControls_Constancy) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CheckAccessControls_Constancy) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CheckAccessControls_Constancy) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckAccessControls_Constancy[this_-1].Name_
}
func (this_ CheckAccessControls_Constancy) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CheckAccessControls_Constancy[this_-1].Name_
}
func (this_ CheckAccessControls_Constancy) Equals(other_ any) bool {
	o_, ok_ := other_.(CheckAccessControls_Constancy)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CheckAccessControls_Constancy) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CheckAccessControls_Constancy[this_-1].HashCode()
}
func (this_ CheckAccessControls_Constancy) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckAccessControls_Constancy[this_-1].CompareTo(other_)
}
func (this_ CheckAccessControls_Constancy) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CheckAccessControls_Constancy[this_-1].Name_)
}
func init() {
	vals_ := []any{CheckAccessControls_Constancy(1), CheckAccessControls_Constancy(2), CheckAccessControls_Constancy(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CheckAccessControls.Constancy", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CheckAccessControls$Constancy", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckAccessControls.Constancy", "com.google.javascript.jscomp.CheckAccessControls$Constancy", func(v_ any) bool {
		return j2gort.InstanceOf[CheckAccessControls_Constancy](v_)
	})
}

// Java class: com.google.javascript.jscomp.CheckAccessControls.ConstantDeclaration
type CheckAccessControls_ConstantDeclaration struct {
	node       pkg_rhino.Node `gson:"node"`
	idHash_    atomic.Int32
	annotation CheckAccessControls_Constancy `gson:"annotation"`
}

func (this_ *CheckAccessControls_ConstantDeclaration) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckAccessControls_ConstantDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*CheckAccessControls_ConstantDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckAccessControls$ConstantDeclaration", "ConstantDeclaration"
}
func (this_ *CheckAccessControls_ConstantDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckAccessControls_ConstantDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckAccessControls.ConstantDeclaration", "com.google.javascript.jscomp.CheckAccessControls$ConstantDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckAccessControls_ConstantDeclaration](v_)
	})
}

func (this_ *CheckAccessControls_ConstantDeclaration) initFields_CheckAccessControls_ConstantDeclaration() {
	_ = this_
}

// Java constructor: ConstantDeclaration(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.CheckAccessControls.Constancy annotation)
func NewCheckAccessControls_ConstantDeclaration(node pkg_rhino.Node, annotation CheckAccessControls_Constancy) *CheckAccessControls_ConstantDeclaration {
	this_ := &CheckAccessControls_ConstantDeclaration{}
	InitCheckAccessControls_ConstantDeclaration(this_, node, annotation)
	return this_
}

func InitCheckAccessControls_ConstantDeclaration(this_ *CheckAccessControls_ConstantDeclaration, node pkg_rhino.Node, annotation CheckAccessControls_Constancy) {
	_ = node
	_ = annotation
	this_.initFields_CheckAccessControls_ConstantDeclaration()
	this_.node = node
	this_.annotation = annotation
}

// Java enum: com.google.javascript.jscomp.CheckAccessControls.IdentifierBehaviour
type CheckAccessControls_IdentifierBehaviour uint8

const (
	enum_CheckAccessControls_IdentifierBehaviour_NON_CONSTRUCTOR      CheckAccessControls_IdentifierBehaviour = 1
	enum_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_INVOCATION CheckAccessControls_IdentifierBehaviour = 2
	enum_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE  CheckAccessControls_IdentifierBehaviour = 3
	enum_CheckAccessControls_IdentifierBehaviour_ES6_CLASS_INVOCATION CheckAccessControls_IdentifierBehaviour = 4
	enum_CheckAccessControls_IdentifierBehaviour_ES6_CLASS_NAMESPACE  CheckAccessControls_IdentifierBehaviour = 5
)

var enumBaseTable_CheckAccessControls_IdentifierBehaviour = [...]j2gort.EnumBase{
	{Name_: "NON_CONSTRUCTOR", Ordinal_: 0},
	{Name_: "ES5_CLASS_INVOCATION", Ordinal_: 1},
	{Name_: "ES5_CLASS_NAMESPACE", Ordinal_: 2},
	{Name_: "ES6_CLASS_INVOCATION", Ordinal_: 3},
	{Name_: "ES6_CLASS_NAMESPACE", Ordinal_: 4},
}
var enumValues_CheckAccessControls_IdentifierBehaviour = []CheckAccessControls_IdentifierBehaviour{1, 2, 3, 4, 5}

func StaticGet_CheckAccessControls_IdentifierBehaviour_NON_CONSTRUCTOR() CheckAccessControls_IdentifierBehaviour {
	return enum_CheckAccessControls_IdentifierBehaviour_NON_CONSTRUCTOR
}
func StaticGet_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_INVOCATION() CheckAccessControls_IdentifierBehaviour {
	return enum_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_INVOCATION
}
func StaticGet_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE() CheckAccessControls_IdentifierBehaviour {
	return enum_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE
}
func StaticGet_CheckAccessControls_IdentifierBehaviour_ES6_CLASS_INVOCATION() CheckAccessControls_IdentifierBehaviour {
	return enum_CheckAccessControls_IdentifierBehaviour_ES6_CLASS_INVOCATION
}
func StaticGet_CheckAccessControls_IdentifierBehaviour_ES6_CLASS_NAMESPACE() CheckAccessControls_IdentifierBehaviour {
	return enum_CheckAccessControls_IdentifierBehaviour_ES6_CLASS_NAMESPACE
}

func CheckAccessControls_IdentifierBehaviour_Values() []CheckAccessControls_IdentifierBehaviour {
	return j2gort.CloneSlice(enumValues_CheckAccessControls_IdentifierBehaviour)
}

func CheckAccessControls_IdentifierBehaviour_ValueOf(name string) CheckAccessControls_IdentifierBehaviour {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CheckAccessControls_IdentifierBehaviour {
		if enumBaseTable_CheckAccessControls_IdentifierBehaviour[i_].Name_ == name {
			return CheckAccessControls_IdentifierBehaviour(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CheckAccessControls.IdentifierBehaviour." + name))
}

// Java static method: IdentifierBehaviour.select(com.google.javascript.rhino.Node target)
func CheckAccessControls_IdentifierBehaviour_Select(target pkg_rhino.Node) CheckAccessControls_IdentifierBehaviour {
	_ = target
	var type_ pkg_rhino.JSType = target.GetJSType()
	_ = type_
	if (type_ == nil || type_.IsJavaNil_()) || !(type_.IsFunctionType()) {
		return StaticGet_CheckAccessControls_IdentifierBehaviour_NON_CONSTRUCTOR()
	}
	var ctorType *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
	_ = ctorType
	if !(ctorType.IsConstructor()) {
		return StaticGet_CheckAccessControls_IdentifierBehaviour_NON_CONSTRUCTOR()
	}
	var isInvocation bool = (NodeUtil_IsInvocationTarget(target) || CheckAccessControls_IsExtendsTarget(target))
	_ = isInvocation
	var isEs6 bool = ((ctorType.GetSource() != 0) && ctorType.GetSource().IsClass())
	_ = isEs6
	if !(isEs6) {
		if isInvocation {
			return StaticGet_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_INVOCATION()
		} else {
			return StaticGet_CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE()
		}
	} else {
		if isInvocation {
			return StaticGet_CheckAccessControls_IdentifierBehaviour_ES6_CLASS_INVOCATION()
		} else {
			return StaticGet_CheckAccessControls_IdentifierBehaviour_ES6_CLASS_NAMESPACE()
		}
	}
}

func (this_ CheckAccessControls_IdentifierBehaviour) IsJavaNil_() bool                              { return this_ == 0 }
func (CheckAccessControls_IdentifierBehaviour) HasIdentityEquals_()                                 {}
func (CheckAccessControls_IdentifierBehaviour) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CheckAccessControls_IdentifierBehaviour) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CheckAccessControls_IdentifierBehaviour) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CheckAccessControls_IdentifierBehaviour) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckAccessControls$IdentifierBehaviour", "IdentifierBehaviour"
}
func (this_ CheckAccessControls_IdentifierBehaviour) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CheckAccessControls_IdentifierBehaviour[this_-1]
}
func (this_ CheckAccessControls_IdentifierBehaviour) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CheckAccessControls_IdentifierBehaviour) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CheckAccessControls_IdentifierBehaviour) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckAccessControls_IdentifierBehaviour[this_-1].Name_
}
func (this_ CheckAccessControls_IdentifierBehaviour) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CheckAccessControls_IdentifierBehaviour[this_-1].Name_
}
func (this_ CheckAccessControls_IdentifierBehaviour) Equals(other_ any) bool {
	o_, ok_ := other_.(CheckAccessControls_IdentifierBehaviour)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CheckAccessControls_IdentifierBehaviour) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CheckAccessControls_IdentifierBehaviour[this_-1].HashCode()
}
func (this_ CheckAccessControls_IdentifierBehaviour) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckAccessControls_IdentifierBehaviour[this_-1].CompareTo(other_)
}
func (this_ CheckAccessControls_IdentifierBehaviour) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CheckAccessControls_IdentifierBehaviour[this_-1].Name_)
}
func init() {
	vals_ := []any{CheckAccessControls_IdentifierBehaviour(1), CheckAccessControls_IdentifierBehaviour(2), CheckAccessControls_IdentifierBehaviour(3), CheckAccessControls_IdentifierBehaviour(4), CheckAccessControls_IdentifierBehaviour(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CheckAccessControls.IdentifierBehaviour", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CheckAccessControls$IdentifierBehaviour", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckAccessControls.IdentifierBehaviour", "com.google.javascript.jscomp.CheckAccessControls$IdentifierBehaviour", func(v_ any) bool {
		return j2gort.InstanceOf[CheckAccessControls_IdentifierBehaviour](v_)
	})
}

// Java record: com.google.javascript.jscomp.CheckAccessControls.PropertyReference
type CheckAccessControls_PropertyReference struct {
	name             string                                 `gson:"name"`
	receiverType     pkg_rhino.ObjectType                   `gson:"receiverType"`
	readableTypeName j2gort.Ext_java_util_function_Supplier `gson:"readableTypeName"`
	sourceNode       pkg_rhino.Node                         `gson:"sourceNode"`
	mutation         bool                                   `gson:"mutation"`
	declaration      bool                                   `gson:"declaration"`
	override         bool                                   `gson:"override"`
}

func (this_ *CheckAccessControls_PropertyReference) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckAccessControls_PropertyReference) IsJavaNil_() bool { return this_ == nil }
func (*CheckAccessControls_PropertyReference) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckAccessControls$PropertyReference", "PropertyReference"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckAccessControls.PropertyReference", "com.google.javascript.jscomp.CheckAccessControls$PropertyReference", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckAccessControls_PropertyReference](v_)
	})
}
func (*CheckAccessControls_PropertyReference) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *CheckAccessControls_PropertyReference) initFields_CheckAccessControls_PropertyReference() {
	_ = this_
}

// Java constructor: PropertyReference(com.google.javascript.rhino.Node sourceNode, java.lang.String name, com.google.javascript.rhino.jstype.ObjectType receiverType, boolean mutation, boolean declaration, boolean override, java.util.function.Supplier<java.lang.String> readableTypeName)
func NewCheckAccessControls_PropertyReference(sourceNode pkg_rhino.Node, name string, receiverType pkg_rhino.ObjectType, mutation bool, declaration bool, override bool, readableTypeName j2gort.Ext_java_util_function_Supplier) *CheckAccessControls_PropertyReference {
	this_ := &CheckAccessControls_PropertyReference{}
	InitCheckAccessControls_PropertyReference(this_, sourceNode, name, receiverType, mutation, declaration, override, readableTypeName)
	return this_
}

func InitCheckAccessControls_PropertyReference(this_ *CheckAccessControls_PropertyReference, sourceNode pkg_rhino.Node, name string, receiverType pkg_rhino.ObjectType, mutation bool, declaration bool, override bool, readableTypeName j2gort.Ext_java_util_function_Supplier) {
	_ = sourceNode
	_ = name
	_ = receiverType
	_ = mutation
	_ = declaration
	_ = override
	_ = readableTypeName
	this_.initFields_CheckAccessControls_PropertyReference()
	j2gort.CheckNotNullNode(sourceNode)
	j2gort.CheckNotNullString(name)
	j2gort.CheckNotNullIface(receiverType)
	j2gort.CheckNotNull(readableTypeName)
	this_.sourceNode = sourceNode
	this_.name = name
	this_.receiverType = receiverType
	this_.mutation = mutation
	this_.declaration = declaration
	this_.override = override
	this_.readableTypeName = readableTypeName
}

// Java static method: PropertyReference.builder()
func CheckAccessControls_PropertyReference_S_Builder() CheckAccessControls_PropertyReference_Builder {
	return NewAutoBuilder_CheckAccessControls_PropertyReference_Builder()
}

// Java method: PropertyReference.getParentNode()
func (this_ *CheckAccessControls_PropertyReference) GetParentNode() pkg_rhino.Node {
	return this_.SourceNode().GetParent()
}

// Java method: PropertyReference.getJSType()
func (this_ *CheckAccessControls_PropertyReference) GetJSType() pkg_rhino.JSType {
	return this_.SourceNode().GetJSType()
}

// Java method: PropertyReference.getJSDocInfo()
func (this_ *CheckAccessControls_PropertyReference) GetJSDocInfo() *pkg_rhino.JSDocInfo {
	return NodeUtil_GetBestJSDocInfo(this_.SourceNode())
}

// Java method: PropertyReference.isDocumentedDeclaration()
func (this_ *CheckAccessControls_PropertyReference) IsDocumentedDeclaration() bool {
	return (this_.Declaration() && (this_.GetJSDocInfo() != nil))
}

// Java method: PropertyReference.isDeletion()
func (this_ *CheckAccessControls_PropertyReference) IsDeletion() bool {
	return this_.SourceNode().GetParent().IsDelProp()
}

// Java method: PropertyReference.getReadableTypeNameOrDefault()
func (this_ *CheckAccessControls_PropertyReference) GetReadableTypeNameOrDefault() string {
	var preferred string = j2gort.CastString(this_.ReadableTypeName().Get())
	_ = preferred
	if j2gort.CheckNotNullString(preferred) == "" {
		return this_.ReceiverType().String()
	} else {
		return preferred
	}
}

func (this_ *CheckAccessControls_PropertyReference) SourceNode() pkg_rhino.Node {
	return this_.sourceNode
}
func (this_ *CheckAccessControls_PropertyReference) Name() string { return this_.name }
func (this_ *CheckAccessControls_PropertyReference) ReceiverType() pkg_rhino.ObjectType {
	return this_.receiverType
}
func (this_ *CheckAccessControls_PropertyReference) Mutation() bool    { return this_.mutation }
func (this_ *CheckAccessControls_PropertyReference) Declaration() bool { return this_.declaration }
func (this_ *CheckAccessControls_PropertyReference) Override() bool    { return this_.override }
func (this_ *CheckAccessControls_PropertyReference) ReadableTypeName() j2gort.Ext_java_util_function_Supplier {
	return this_.readableTypeName
}
func (this_ *CheckAccessControls_PropertyReference) Equals(other_ any) bool {
	o_, ok_ := other_.(*CheckAccessControls_PropertyReference)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.sourceNode != o_.sourceNode {
		return false
	}
	if this_.name != o_.name {
		return false
	}
	if !javalang.ObjectsEquals(this_.receiverType, o_.receiverType) {
		return false
	}
	if this_.mutation != o_.mutation {
		return false
	}
	if this_.declaration != o_.declaration {
		return false
	}
	if this_.override != o_.override {
		return false
	}
	if !javalang.ObjectsEquals(this_.readableTypeName, o_.readableTypeName) {
		return false
	}
	return true
}
func (this_ *CheckAccessControls_PropertyReference) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + int32(this_.sourceNode)
	{
		hs_24 := this_.name
		h_ = 31 * h_
		if hs_24 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_24)
		}
	}
	{
		hi_25 := this_.receiverType
		h_ = 31 * h_
		if hi_25 != nil && !hi_25.IsJavaNil_() {
			h_ += hi_25.HashCode()
		}
	}
	if this_.mutation {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	if this_.declaration {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	if this_.override {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	h_ = 31*h_ + j2gort.ObjectHash(this_.readableTypeName)
	return h_
}
func (this_ *CheckAccessControls_PropertyReference) String() string {
	_ = fmt.Sprint
	return "PropertyReference[" + "sourceNode=" + j2gort.StrVal(j2gort.BoxNode(this_.sourceNode)) + ", " + "name=" + j2gort.StrVal(j2gort.BoxString(this_.name)) + ", " + "receiverType=" + j2gort.StrVal(this_.receiverType) + ", " + "mutation=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.mutation)) + ", " + "declaration=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.declaration)) + ", " + "override=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.override)) + ", " + "readableTypeName=" + j2gort.StrVal(this_.readableTypeName) + "]"
}

// Java interface: com.google.javascript.jscomp.CheckAccessControls.PropertyReference.Builder
type CheckAccessControls_PropertyReference_Builder interface {
	Implements_jscomp_CheckAccessControls_PropertyReference_Builder()
	IsJavaNil_() bool
	SetSourceNode(pkg_rhino.Node) CheckAccessControls_PropertyReference_Builder
	SetName(string) CheckAccessControls_PropertyReference_Builder
	SetReceiverType(pkg_rhino.ObjectType) CheckAccessControls_PropertyReference_Builder
	SetMutation(bool) CheckAccessControls_PropertyReference_Builder
	SetDeclaration(bool) CheckAccessControls_PropertyReference_Builder
	SetOverride(bool) CheckAccessControls_PropertyReference_Builder
	SetReadableTypeName(j2gort.Ext_java_util_function_Supplier) CheckAccessControls_PropertyReference_Builder
	Build() *CheckAccessControls_PropertyReference
}

func ToIface_CheckAccessControls_PropertyReference_Builder[T interface {
	CheckAccessControls_PropertyReference_Builder
	comparable
}](v T) CheckAccessControls_PropertyReference_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckAccessControls.PropertyReference.Builder", "com.google.javascript.jscomp.CheckAccessControls$PropertyReference$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[CheckAccessControls_PropertyReference_Builder](v_)
	})
}
func (*CheckAccessControls) CanReleaseSyntacticScopes_() {}
