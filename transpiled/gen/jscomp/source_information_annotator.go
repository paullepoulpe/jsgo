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
// Port of com/google/javascript/jscomp/SourceInformationAnnotator.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.SourceInformationAnnotator
type SourceInformationAnnotator struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	sourceFileToCheck string `gson:"sourceFileToCheck"`
}

func (this_ *SourceInformationAnnotator) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceInformationAnnotator) IsJavaNil_() bool { return this_ == nil }
func (*SourceInformationAnnotator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceInformationAnnotator", "SourceInformationAnnotator"
}
func (this_ *SourceInformationAnnotator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*SourceInformationAnnotator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceInformationAnnotator", "com.google.javascript.jscomp.SourceInformationAnnotator", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceInformationAnnotator](v_)
	})
}
func (*SourceInformationAnnotator) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*SourceInformationAnnotator) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *SourceInformationAnnotator) initFields_SourceInformationAnnotator() {
	_ = this_
	this_.sourceFileToCheck = j2gort.NullString
}

// Java constructor: SourceInformationAnnotator(java.lang.@org.jspecify.annotations.Nullable String sourceFileToCheck)
func NewSourceInformationAnnotator(sourceFileToCheck string) *SourceInformationAnnotator {
	this_ := &SourceInformationAnnotator{}
	InitSourceInformationAnnotator(this_, sourceFileToCheck)
	return this_
}

func InitSourceInformationAnnotator(this_ *SourceInformationAnnotator, sourceFileToCheck string) {
	_ = sourceFileToCheck
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_SourceInformationAnnotator()
	this_.sourceFileToCheck = sourceFileToCheck
}

// Java static method: SourceInformationAnnotator.create()
func SourceInformationAnnotator_Create() *SourceInformationAnnotator {
	return NewSourceInformationAnnotator(j2gort.NullString)
}

// Java static method: SourceInformationAnnotator.createWithAnnotationChecks(java.lang.String sourceFile)
func SourceInformationAnnotator_CreateWithAnnotationChecks(sourceFile string) *SourceInformationAnnotator {
	_ = sourceFile
	return NewSourceInformationAnnotator(sourceFile)
}

// Java method: SourceInformationAnnotator.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SourceInformationAnnotator) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.sourceFileToCheck != j2gort.NullString {
		{
			var iifeCond_1 bool
			{

				sr_1 := this_.sourceFileToCheck
				se_2 := n.GetSourceFileName()
				iifeCond_1 = j2gort.CheckNotNullString(sr_1) == se_2
			}
			if !(iifeCond_1) {
				j2gort.CheckState(false)
			}
		}
	}
	if SourceInformationAnnotator_IsStringNodeRequiringOriginalName(n) {
		SourceInformationAnnotator_SetOriginalName_2_Node_Node(n, n)
		return
	}
	if n.IsFunction() {
		var functionName string = NodeUtil_GetNearestFunctionName(n)
		_ = functionName
		if functionName != j2gort.NullString {
			SourceInformationAnnotator_SetOriginalName_2_Node_string(n, functionName)
		}
	}
}

// Java static method: SourceInformationAnnotator.isStringNodeRequiringOriginalName(com.google.javascript.rhino.Node node)
func SourceInformationAnnotator_IsStringNodeRequiringOriginalName(node pkg_rhino.Node) bool {
	_ = node
	switch node.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(32) /* Token.NAME */ :
		return true
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(100) /* Token.STRING_KEY */ :
		return (node.GetParent().IsObjectLit() && !(node.IsQuotedStringKey()))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: SourceInformationAnnotator.setOriginalName(com.google.javascript.rhino.Node n, java.lang.String name)
func SourceInformationAnnotator_SetOriginalName_2_Node_string(n pkg_rhino.Node, name string) {
	_ = n
	_ = name
	if !(j2gort.CheckNotNullString(name) == "") && (n.GetOriginalName() == j2gort.NullString) {
		n.SetOriginalName(name)
	}
}

// Java static method: SourceInformationAnnotator.setOriginalName(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node name)
func SourceInformationAnnotator_SetOriginalName_2_Node_Node(n pkg_rhino.Node, name pkg_rhino.Node) {
	_ = n
	_ = name
	if !(j2gort.CheckNotNullString(name.GetString()) == "") && (n.GetOriginalName() == j2gort.NullString) {
		n.SetOriginalNameFromName(name)
	}
}
func (*SourceInformationAnnotator) AlwaysShouldTraverse_()      {}
func (*SourceInformationAnnotator) CanReleaseSyntacticScopes_() {}
