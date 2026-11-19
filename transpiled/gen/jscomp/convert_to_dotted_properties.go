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
// Port of com/google/javascript/jscomp/ConvertToDottedProperties.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ConvertToDottedProperties
type ConvertToDottedProperties struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *ConvertToDottedProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *ConvertToDottedProperties) IsJavaNil_() bool { return this_ == nil }
func (*ConvertToDottedProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConvertToDottedProperties", "ConvertToDottedProperties"
}
func (this_ *ConvertToDottedProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ConvertToDottedProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConvertToDottedProperties", "com.google.javascript.jscomp.ConvertToDottedProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*ConvertToDottedProperties](v_)
	})
}
func (*ConvertToDottedProperties) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ConvertToDottedProperties) Implements_jscomp_CompilerPass()                            {}
func (*ConvertToDottedProperties) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *ConvertToDottedProperties) initFields_ConvertToDottedProperties() {
	_ = this_
}

// Java constructor: ConvertToDottedProperties(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewConvertToDottedProperties(compiler AbstractCompiler) *ConvertToDottedProperties {
	this_ := &ConvertToDottedProperties{}
	InitConvertToDottedProperties(this_, compiler)
	return this_
}

func InitConvertToDottedProperties(this_ *ConvertToDottedProperties, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ConvertToDottedProperties()
	this_.compiler = compiler
}

// Java method: ConvertToDottedProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ConvertToDottedProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	GatherGetterAndSetterProperties_Update(this_.compiler, externs, root)
}

// Java method: ConvertToDottedProperties.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ConvertToDottedProperties) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(129) /* Token.COMPUTED_PROP */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		var leftElem pkg_rhino.Node = n.GetFirstChild()
		_ = leftElem
		var rightElem pkg_rhino.Node = leftElem.GetNext()
		_ = rightElem
		if ((leftElem.IsStringLit() && NodeUtil_IsValidPropertyName(pkg_parser.StaticGet_FeatureSet_ES3(), leftElem.GetString())) && !(j2gort.CheckNotNullString(leftElem.GetString()) == "constructor")) && !(parent.IsObjectLit() && (j2gort.CheckNotNullString(leftElem.GetString()) == "__proto__")) {
			leftElem.Detach()
			if rightElem != 0 {
				rightElem.Detach()
			}
			var temp pkg_rhino.Node = pkg_rhino.Node(0)
			_ = temp
			if n.IsComputedProp() {
				if rightElem.IsFunction() {
					if n.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER()) {
						temp = pkg_rhino.IR_GetterDef(leftElem.GetString(), rightElem)
					} else if n.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER()) {
						temp = pkg_rhino.IR_SetterDef(leftElem.GetString(), rightElem)
					} else {
						if n.GetParent().IsClassMembers() {
							temp = pkg_rhino.IR_MemberFunctionDef(leftElem.GetString(), rightElem)
							NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_MEMBER_DECLARATIONS(), this_.compiler)
						} else {
							temp = pkg_rhino.IR_StringKey_2(leftElem.GetString(), rightElem)
						}
					}
				} else {
					temp = pkg_rhino.IR_StringKey_2(leftElem.GetString(), rightElem)
				}
			} else if rightElem == 0 {
				temp = pkg_rhino.Node_NewString_2(pkg_rhino.StaticGet_Token_MEMBER_FIELD_DEF(), leftElem.GetString())
			} else {
				temp = pkg_rhino.IR_MemberFieldDef(leftElem.GetString(), rightElem)
			}
			temp.SetStaticMember(n.IsStaticMember())
			n.ReplaceWith(temp)
			this_.compiler.ReportChangeToEnclosingScope(temp)
		}
	case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(100) /* Token.STRING_KEY */ :
		if NodeUtil_IsValidPropertyName(pkg_parser.StaticGet_FeatureSet_ES3(), n.GetString()) {
			if n.GetBooleanProp(pkg_rhino.StaticGet_Node_QUOTED_PROP()) {
				n.PutBooleanProp(pkg_rhino.StaticGet_Node_QUOTED_PROP(), false)
				this_.compiler.ReportChangeToEnclosingScope(n)
			}
		}
	case int32(30) /* Token.OPTCHAIN_GETELEM */, int32(27) /* Token.GETELEM */ :
		var left pkg_rhino.Node = n.GetFirstChild()
		_ = left
		var right pkg_rhino.Node = left.GetNext()
		_ = right
		if right.IsStringLit() && NodeUtil_IsValidPropertyName(pkg_parser.StaticGet_FeatureSet_ES3(), right.GetString()) {
			left.Detach()
			right.Detach()
			var newGetProp pkg_rhino.Node
			_ = newGetProp
			if n.IsGetElem() {
				newGetProp = pkg_rhino.IR_Getprop_2(left, right.GetString())
			} else if n.IsOptionalChainStart() {
				newGetProp = pkg_rhino.IR_StartOptChainGetprop(left, right.GetString())
			} else {
				newGetProp = pkg_rhino.IR_ContinueOptChainGetprop(left, right.GetString())
			}
			n.ReplaceWith(newGetProp)
			this_.compiler.ReportChangeToEnclosingScope(newGetProp)
		}
	default:
	}
}
func (*ConvertToDottedProperties) AlwaysShouldTraverse_()      {}
func (*ConvertToDottedProperties) CanReleaseSyntacticScopes_() {}
