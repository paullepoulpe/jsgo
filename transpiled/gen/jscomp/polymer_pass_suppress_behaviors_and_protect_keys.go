/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PolymerPassSuppressBehaviorsAndProtectKeys.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.PolymerPassSuppressBehaviorsAndProtectKeys
type PolymerPassSuppressBehaviorsAndProtectKeys struct {
	NodeTraversal_ExternsSkippingCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *PolymerPassSuppressBehaviorsAndProtectKeys) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerPassSuppressBehaviorsAndProtectKeys) IsJavaNil_() bool { return this_ == nil }
func (*PolymerPassSuppressBehaviorsAndProtectKeys) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerPassSuppressBehaviorsAndProtectKeys", "PolymerPassSuppressBehaviorsAndProtectKeys"
}
func (this_ *PolymerPassSuppressBehaviorsAndProtectKeys) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_ExternsSkippingCallbackBase.IdentityHash_()
}
func (*PolymerPassSuppressBehaviorsAndProtectKeys) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerPassSuppressBehaviorsAndProtectKeys", "com.google.javascript.jscomp.PolymerPassSuppressBehaviorsAndProtectKeys", func(v_ any) bool {
		return j2gort.InstanceOf[*PolymerPassSuppressBehaviorsAndProtectKeys](v_)
	})
}
func (*PolymerPassSuppressBehaviorsAndProtectKeys) Implements_jscomp_NodeTraversal_ExternsSkippingCallback() {
}
func (*PolymerPassSuppressBehaviorsAndProtectKeys) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *PolymerPassSuppressBehaviorsAndProtectKeys) initFields_PolymerPassSuppressBehaviorsAndProtectKeys() {
	_ = this_
}

// Java constructor: PolymerPassSuppressBehaviorsAndProtectKeys(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewPolymerPassSuppressBehaviorsAndProtectKeys(compiler AbstractCompiler) *PolymerPassSuppressBehaviorsAndProtectKeys {
	this_ := &PolymerPassSuppressBehaviorsAndProtectKeys{}
	InitPolymerPassSuppressBehaviorsAndProtectKeys(this_, compiler)
	return this_
}

func InitPolymerPassSuppressBehaviorsAndProtectKeys(this_ *PolymerPassSuppressBehaviorsAndProtectKeys, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_ExternsSkippingCallback(&this_.NodeTraversal_ExternsSkippingCallbackBase)
	this_.initFields_PolymerPassSuppressBehaviorsAndProtectKeys()
	this_.compiler = compiler
}

// Java method: PolymerPassSuppressBehaviorsAndProtectKeys.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *PolymerPassSuppressBehaviorsAndProtectKeys) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if PolymerPassSuppressBehaviorsAndProtectKeys_IsBehavior(n) {
		if !(NodeUtil_IsNameDeclaration(n)) && !(n.IsAssign()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR(), []string{}))
			return
		}
		this_.priv_PolymerPassSuppressBehaviorsAndProtectKeys_traverseBehaviorsAndProtectKeys_4423(n)
		var newDocs *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(n.GetJSDocInfo())
		_ = newDocs
		newDocs.RecordNoCollapse()
		n.SetJSDocInfo(newDocs.Build_0())
		var behaviorValue pkg_rhino.Node = n.GetSecondChild()
		_ = behaviorValue
		if NodeUtil_IsNameDeclaration(n) {
			behaviorValue = n.GetFirstFirstChild()
		}
		this_.priv_PolymerPassSuppressBehaviorsAndProtectKeys_suppressBehavior_4424(behaviorValue, n)
	}
}

// Java method: PolymerPassSuppressBehaviorsAndProtectKeys.traverseBehaviorsAndProtectKeys(com.google.javascript.rhino.Node n)
func (this_ *PolymerPassSuppressBehaviorsAndProtectKeys) priv_PolymerPassSuppressBehaviorsAndProtectKeys_traverseBehaviorsAndProtectKeys_4423(n pkg_rhino.Node) {
	_ = n
	if n.IsObjectLit() {
		PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(n, this_.compiler)
		PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys(n)
	}
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_PolymerPassSuppressBehaviorsAndProtectKeys_traverseBehaviorsAndProtectKeys_4423(child)
		}
	}
}

// Java method: PolymerPassSuppressBehaviorsAndProtectKeys.suppressBehavior(com.google.javascript.rhino.Node behaviorValue, com.google.javascript.rhino.Node reportNode)
func (this_ *PolymerPassSuppressBehaviorsAndProtectKeys) priv_PolymerPassSuppressBehaviorsAndProtectKeys_suppressBehavior_4424(behaviorValue pkg_rhino.Node, reportNode pkg_rhino.Node) {
	_ = behaviorValue
	_ = reportNode
	if behaviorValue == 0 {
		this_.compiler.Report_1(JSError_Make_3(reportNode, StaticGet_PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR(), []string{}))
		return
	}
	if behaviorValue.IsArrayLit() {
		{
			var child pkg_rhino.Node = behaviorValue.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				this_.priv_PolymerPassSuppressBehaviorsAndProtectKeys_suppressBehavior_4424(child, behaviorValue)
			}
		}
	} else if behaviorValue.IsObjectLit() {
		this_.priv_PolymerPassSuppressBehaviorsAndProtectKeys_stripPropertyTypes_4425(behaviorValue)
		this_.priv_PolymerPassSuppressBehaviorsAndProtectKeys_addBehaviorSuppressions_4427(behaviorValue)
	}
}

// Java static method: PolymerPassSuppressBehaviorsAndProtectKeys.isBehavior(com.google.javascript.rhino.Node value)
func PolymerPassSuppressBehaviorsAndProtectKeys_IsBehavior(value pkg_rhino.Node) bool {
	_ = value
	return ((value.GetJSDocInfo() != nil) && value.GetJSDocInfo().IsPolymerBehavior())
}

// Java method: PolymerPassSuppressBehaviorsAndProtectKeys.stripPropertyTypes(com.google.javascript.rhino.Node behaviorValue)
func (this_ *PolymerPassSuppressBehaviorsAndProtectKeys) priv_PolymerPassSuppressBehaviorsAndProtectKeys_stripPropertyTypes_4425(behaviorValue pkg_rhino.Node) {
	_ = behaviorValue
	var properties []*PolymerPass_MemberDefinition = j2gort.ImmListToSliceNullable[*PolymerPass_MemberDefinition](PolymerPassStaticUtils_ExtractProperties(behaviorValue, StaticGet_PolymerClassDefinition_DefinitionType_ObjectLiteral(), this_.compiler, pkg_rhino.Node(0)))
	_ = properties
	{

		for _, property := range properties {

			_ = property
			property.name.SetJSDocInfo(nil)
		}
	}
}

// Java method: PolymerPassSuppressBehaviorsAndProtectKeys.suppressDefaultValues(com.google.javascript.rhino.Node behaviorValue)
func (this_ *PolymerPassSuppressBehaviorsAndProtectKeys) priv_PolymerPassSuppressBehaviorsAndProtectKeys_suppressDefaultValues_4426(behaviorValue pkg_rhino.Node) {
	_ = behaviorValue
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(PolymerPassStaticUtils_ExtractProperties(behaviorValue, StaticGet_PolymerClassDefinition_DefinitionType_ObjectLiteral(), this_.compiler, pkg_rhino.Node(0))))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			property := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_4)
			_ = property
			if !(property.value.IsObjectLit()) {
				continue
			}
			var defaultValue pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(property.value, "value")
			_ = defaultValue
			if (defaultValue == 0) || !(defaultValue.IsFunction()) {
				continue
			}
			var defaultValueKey pkg_rhino.Node = defaultValue.GetParent()
			_ = defaultValueKey
			this_.priv_PolymerPassSuppressBehaviorsAndProtectKeys_addSuppressionJsDocAnnotations_4428(defaultValueKey)
		}
	}
}

// Java method: PolymerPassSuppressBehaviorsAndProtectKeys.addBehaviorSuppressions(com.google.javascript.rhino.Node behaviorValue)
func (this_ *PolymerPassSuppressBehaviorsAndProtectKeys) priv_PolymerPassSuppressBehaviorsAndProtectKeys_addBehaviorSuppressions_4427(behaviorValue pkg_rhino.Node) {
	_ = behaviorValue
	{
		var keyNode pkg_rhino.Node = behaviorValue.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			if keyNode.GetFirstChild().IsFunction() {
				this_.priv_PolymerPassSuppressBehaviorsAndProtectKeys_addSuppressionJsDocAnnotations_4428(keyNode)
			}
		}
	}
	this_.priv_PolymerPassSuppressBehaviorsAndProtectKeys_suppressDefaultValues_4426(behaviorValue)
}

// Java method: PolymerPassSuppressBehaviorsAndProtectKeys.addSuppressionJsDocAnnotations(com.google.javascript.rhino.Node node)
func (this_ *PolymerPassSuppressBehaviorsAndProtectKeys) priv_PolymerPassSuppressBehaviorsAndProtectKeys_addSuppressionJsDocAnnotations_4428(node pkg_rhino.Node) {
	_ = node
	var suppressDoc *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(node.GetJSDocInfo())
	_ = suppressDoc
	suppressDoc.RecordSuppressions_1(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_3("checkTypes", "globalThis", "visibility")))
	node.SetJSDocInfo(suppressDoc.Build_0())
}
func (*PolymerPassSuppressBehaviorsAndProtectKeys) CanReleaseSyntacticScopes_() {}
