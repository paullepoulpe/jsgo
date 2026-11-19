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
// Port of com/google/javascript/jscomp/InferJSDocInfo.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.InferJSDocInfo
type InferJSDocInfo struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *InferJSDocInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *InferJSDocInfo) IsJavaNil_() bool { return this_ == nil }
func (*InferJSDocInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InferJSDocInfo", "InferJSDocInfo"
}
func (this_ *InferJSDocInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*InferJSDocInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InferJSDocInfo", "com.google.javascript.jscomp.InferJSDocInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*InferJSDocInfo](v_)
	})
}
func (*InferJSDocInfo) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*InferJSDocInfo) Implements_jscomp_CompilerPass()                            {}
func (*InferJSDocInfo) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *InferJSDocInfo) initFields_InferJSDocInfo() {
	_ = this_
}

// Java constructor: InferJSDocInfo(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewInferJSDocInfo(compiler AbstractCompiler) *InferJSDocInfo {
	this_ := &InferJSDocInfo{}
	InitInferJSDocInfo(this_, compiler)
	return this_
}

func InitInferJSDocInfo(this_ *InferJSDocInfo, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_InferJSDocInfo()
	this_.compiler = compiler
}

// Java method: InferJSDocInfo.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InferJSDocInfo) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if externs != 0 {
		NodeTraversal_Traverse(this_.compiler, externs, this_)
	}
	if root != 0 {
		NodeTraversal_Traverse(this_.compiler, root, this_)
	}
}

// Java method: InferJSDocInfo.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InferJSDocInfo) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		this_.priv_InferJSDocInfo_inferJSDocForName_3198(n, parent)
		return
	case int32(100) /* Token.STRING_KEY */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(129) /* Token.COMPUTED_PROP */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		this_.priv_InferJSDocInfo_inferJSDocForObjectKeyOrClassField_3199(n, parent)
		return
	case int32(26) /* Token.GETPROP */ :
		this_.priv_InferJSDocInfo_inferJSDocForProperty_3200(n, parent)
		return
	default:
		return
	}
}

// Java method: InferJSDocInfo.inferJSDocForName(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InferJSDocInfo) priv_InferJSDocInfo_inferJSDocForName_3198(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	if parent == 0 {
		return
	}
	var typeDoc *pkg_rhino.JSDocInfo = nil
	_ = typeDoc
	var aliasedType pkg_rhino.JSType = nil
	_ = aliasedType
	var inferredType pkg_rhino.JSType = nil
	_ = inferredType
	if NodeUtil_IsNameDeclaration(parent) {
		var nameInfo *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
		_ = nameInfo
		if nameInfo != nil {
			typeDoc = nameInfo
		} else {
			typeDoc = parent.GetJSDocInfo()
		}
		inferredType = n.GetJSType()
		var value pkg_rhino.Node = n.GetFirstChild()
		_ = value
		if (value != 0) && value.IsQualifiedName() {
			aliasedType = value.GetJSType()
		} else {
			aliasedType = nil
		}
	} else if NodeUtil_IsFunctionDeclaration(parent) || NodeUtil_IsClassDeclaration(parent) {
		typeDoc = parent.GetJSDocInfo()
		inferredType = parent.GetJSType()
		aliasedType = nil
	} else if parent.IsAssign() && n.IsFirstChildOf(parent) {
		typeDoc = parent.GetJSDocInfo()
		inferredType = n.GetJSType()
		var value_2 pkg_rhino.Node = n.GetNext()
		_ = value_2
		if value_2.IsQualifiedName() {
			aliasedType = value_2.GetJSType()
		} else {
			aliasedType = nil
		}
	} else {
		return
	}
	if typeDoc == nil {
		return
	}
	var objType pkg_rhino.ObjectType = InferJSDocInfo_Dereferenced(inferredType)
	_ = objType
	if InferJSDocInfo_ShouldAttachJSDocToNominalTypeOrShape(objType, aliasedType) {
		InferJSDocInfo_AttachJSDocInfoToNominalTypeOrShape(objType, typeDoc, n)
	}
}

// Java method: InferJSDocInfo.inferJSDocForObjectKeyOrClassField(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InferJSDocInfo) priv_InferJSDocInfo_inferJSDocForObjectKeyOrClassField_3199(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	var typeDoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = typeDoc
	if typeDoc == nil {
		return
	}
	var owningType pkg_rhino.ObjectType = nil
	_ = owningType
	if parent.IsClassMembers() {
		var ctorType *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(parent.GetParent().GetJSType())
		_ = ctorType
		if ctorType == nil {
			return
		}
		if n.IsStaticMember() {
			owningType = pkg_rhino.ToIface_ObjectType(ctorType)
		} else if n.IsMemberFieldDef() {
			owningType = ctorType.GetInstanceType()
		} else {
			owningType = ctorType.GetPrototype()
		}
	} else {
		owningType = InferJSDocInfo_Dereferenced(parent.GetJSType())
	}
	if owningType == nil || owningType.IsJavaNil_() {
		return
	}
	if n.IsComputedFieldDef() || n.IsComputedProp() {
		var keyType pkg_rhino.JSType = n.GetFirstChild().GetJSType()
		_ = keyType
		if (keyType == nil || keyType.IsJavaNil_()) || !(keyType.GetJSTypeBase().IsKnownSymbolValueType()) {
			return
		}
		var propName pkg_rhino.Property_Key = pkg_rhino.NewProperty_SymbolKey(keyType.ToMaybeKnownSymbolType())
		_ = propName
		if owningType.GetObjectTypeBase().HasOwnProperty_1_Property_Key(propName) && (owningType.GetObjectTypeBase().GetPropertyJSDocInfo_1_Property_Key(propName) == nil) {
			owningType.SetPropertyJSDocInfo_2_Property_Key_JSDocInfo(propName, typeDoc)
		}
	} else {
		var propName_2 string = n.GetString()
		_ = propName_2
		if owningType.GetObjectTypeBase().HasOwnProperty_1_string(propName_2) && (owningType.GetObjectTypeBase().GetPropertyJSDocInfo_1_string(propName_2) == nil) {
			owningType.SetPropertyJSDocInfo_2_string_JSDocInfo(propName_2, typeDoc)
		}
	}
}

// Java method: InferJSDocInfo.inferJSDocForProperty(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InferJSDocInfo) priv_InferJSDocInfo_inferJSDocForProperty_3200(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	var typeDoc *pkg_rhino.JSDocInfo = nil
	_ = typeDoc
	var aliasedType pkg_rhino.JSType = nil
	_ = aliasedType
	if parent.IsAssign() && n.IsFirstChildOf(parent) {
		typeDoc = parent.GetJSDocInfo()
		var rhs pkg_rhino.Node = n.GetNext()
		_ = rhs
		if rhs.IsQualifiedName() {
			aliasedType = rhs.GetJSType()
		} else {
			aliasedType = nil
		}
	} else if parent.IsExprResult() {
		typeDoc = n.GetJSDocInfo()
		aliasedType = nil
	} else {
		return
	}
	if (typeDoc == nil) || !(typeDoc.ContainsDeclaration()) {
		return
	}
	var lhsType pkg_rhino.ObjectType = InferJSDocInfo_Dereferenced(n.GetFirstChild().GetJSType())
	_ = lhsType
	if lhsType == nil || lhsType.IsJavaNil_() {
		return
	}
	base_lhsType_ := lhsType.GetObjectTypeBase()
	var propName string = n.GetString()
	_ = propName
	if base_lhsType_.HasOwnProperty_1_string(propName) && (base_lhsType_.GetPropertyJSDocInfo_1_string(propName) == nil) {
		lhsType.SetPropertyJSDocInfo_2_string_JSDocInfo(propName, typeDoc)
	}
	var propType pkg_rhino.ObjectType = InferJSDocInfo_Dereferenced(base_lhsType_.GetPropertyType_1_string(propName))
	_ = propType
	if InferJSDocInfo_ShouldAttachJSDocToNominalTypeOrShape(propType, aliasedType) {
		InferJSDocInfo_AttachJSDocInfoToNominalTypeOrShape(propType, typeDoc, n)
	}
}

// Java static method: InferJSDocInfo.dereferenced(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type)
func InferJSDocInfo_Dereferenced(type_ pkg_rhino.JSType) pkg_rhino.ObjectType {
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return nil
	} else {
		return type_.GetJSTypeBase().Dereference()
	}
}

// Java static method: InferJSDocInfo.shouldAttachJSDocToNominalTypeOrShape(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType type, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType aliasedType)
func InferJSDocInfo_ShouldAttachJSDocToNominalTypeOrShape(type_ pkg_rhino.ObjectType, aliasedType pkg_rhino.JSType) bool {
	_ = type_
	_ = aliasedType
	return (((type_ != nil && !type_.IsJavaNil_()) && (type_.GetJSDocInfo() == nil)) && !(type_.GetObjectTypeBase().GetJSTypeBase().EqualsJSType(aliasedType)))
}

// Java static method: InferJSDocInfo.attachJSDocInfoToNominalTypeOrShape(com.google.javascript.rhino.jstype.ObjectType objType, com.google.javascript.rhino.JSDocInfo docInfo, com.google.javascript.rhino.Node nameNode)
func InferJSDocInfo_AttachJSDocInfoToNominalTypeOrShape(objType pkg_rhino.ObjectType, docInfo *pkg_rhino.JSDocInfo, nameNode pkg_rhino.Node) {
	_ = objType
	_ = docInfo
	_ = nameNode
	base_objType_ := objType.GetObjectTypeBase()
	if base_objType_.IsConstructor() || base_objType_.IsInterface() {
		if !(InferJSDocInfo_IsReferenceNameOf(objType, nameNode)) {
			return
		}
		objType.SetJSDocInfo(docInfo)
		pkg_rhino.JSType_ToMaybeFunctionType(pkg_rhino.Upcast_JSType(objType)).GetInstanceType().SetJSDocInfo(docInfo)
	} else if base_objType_.IsEnumType() {
		var elementType pkg_rhino.ObjectType = pkg_rhino.ToIface_ObjectType(objType.ToMaybeEnumType().GetElementsType())
		_ = elementType
		if !(InferJSDocInfo_IsReferenceNameOf(elementType, nameNode)) {
			return
		}
		objType.SetJSDocInfo(docInfo)
		elementType.SetJSDocInfo(docInfo)
	} else if !(base_objType_.IsNativeObjectType()) && base_objType_.IsFunctionType() {
		objType.SetJSDocInfo(docInfo)
	}
}

// Java static method: InferJSDocInfo.isReferenceNameOf(com.google.javascript.rhino.jstype.ObjectType type, com.google.javascript.rhino.Node nameNode)
func InferJSDocInfo_IsReferenceNameOf(type_ pkg_rhino.ObjectType, nameNode pkg_rhino.Node) bool {
	_ = type_
	_ = nameNode
	if !(type_.GetObjectTypeBase().HasReferenceName()) {
		return false
	}
	var refName string = type_.GetReferenceName()
	_ = refName
	if refName == j2gort.NullString {
		return false
	}
	if nameNode.IsName() {
		{

			sr_1 := refName
			se_2 := nameNode.GetString()
			return j2gort.CheckNotNullString(sr_1) == se_2
		}

	}
	var h int32 = nameNode.HashQualifiedName()
	_ = h
	if (h == int32(0)) || (h != pkg_rhino.Node_HashQualifiedNameString(refName)) {
		return false
	}
	{

		sr_3 := refName
		se_4 := nameNode.GetQualifiedName()
		return j2gort.CheckNotNullString(sr_3) == se_4
	}

}
func (*InferJSDocInfo) AlwaysShouldTraverse_()      {}
func (*InferJSDocInfo) CanReleaseSyntacticScopes_() {}
