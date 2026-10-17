// Port of com/google/javascript/jscomp/ConvertToDottedProperties.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ConvertToDottedProperties converts property accesses from quoted string or
// bracket access syntax to dot or unquoted string syntax, where possible. Dot
// syntax is more compact.
type ConvertToDottedProperties struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

// Java: ConvertToDottedProperties(AbstractCompiler)
func NewConvertToDottedProperties(compiler AbstractCompiler) *ConvertToDottedProperties {
	c := &ConvertToDottedProperties{}
	c.compiler = compiler
	return c
}

// Java: ConvertToDottedProperties.process(Node, Node)
func (c *ConvertToDottedProperties) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
	GatherGetterAndSetterProperties_Update(c.compiler, externs, root)
}

// Java: ConvertToDottedProperties.visit(NodeTraversal, Node, Node)
func (c *ConvertToDottedProperties) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
		leftElem := n.GetFirstChild()
		rightElem := leftElem.GetNext()

		// not convert property named constructor or __proto__ in object literals.
		// ['constructor']() and constructor() are different.
		// {['__proto__']: v} defines an own property; {__proto__: v} sets prototype.
		if leftElem.IsStringLit() &&
			NodeUtil_IsValidPropertyName(parser.FeatureSet_ES3, leftElem.GetString()) &&
			leftElem.GetString() != "constructor" &&
			!(parent.IsObjectLit() && leftElem.GetString() == "__proto__") {
			leftElem.Detach()
			if rightElem != nil {
				// A COMPUTED_FIELD_DEF has no second child when the field has no initializer.
				rightElem.Detach()
			}
			var temp *rhino.Node
			if n.IsComputedProp() {
				if rightElem.IsFunction() {
					if n.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER) {

						temp = rhino.IR_GetterDef(leftElem.GetString(), rightElem)

					} else if n.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER) {
						temp = rhino.IR_SetterDef(leftElem.GetString(), rightElem)
					} else {
						if n.GetParent().IsClassMembers() {

							temp = rhino.IR_MemberFunctionDef(leftElem.GetString(), rightElem)
							NodeUtil_AddFeatureToScript(
								t.GetCurrentScript(), parser.FeatureSet_Feature_MEMBER_DECLARATIONS, c.compiler)
						} else {
							// TODO - further optimize this code to a member function
							temp = rhino.IR_StringKeyWithValue(leftElem.GetString(), rightElem)
						}
					}

				} else {
					temp = rhino.IR_StringKeyWithValue(leftElem.GetString(), rightElem)
				}
			} else if rightElem == nil {
				// An uninitialized field, e.g. `class C { ['x']; }`, becomes a 0-child
				// MEMBER_FIELD_DEF, matching how the parser represents `class C { x; }`.
				temp = rhino.Node_NewStringWithToken(rhino.Token_MEMBER_FIELD_DEF, leftElem.GetString())
			} else {
				temp = rhino.IR_MemberFieldDef(leftElem.GetString(), rightElem)
			}

			temp.SetStaticMember(n.IsStaticMember())
			n.ReplaceWith(temp)
			c.compiler.ReportChangeToEnclosingScope(temp)
		}
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_STRING_KEY:
		if NodeUtil_IsValidPropertyName(parser.FeatureSet_ES3, n.GetString()) {
			if n.GetBooleanProp(rhino.Node_QUOTED_PROP) {
				n.PutBooleanProp(rhino.Node_QUOTED_PROP, false)
				c.compiler.ReportChangeToEnclosingScope(n)
			}
		}
	case rhino.Token_OPTCHAIN_GETELEM, rhino.Token_GETELEM:
		left := n.GetFirstChild()
		right := left.GetNext()
		if right.IsStringLit() &&
			NodeUtil_IsValidPropertyName(parser.FeatureSet_ES3, right.GetString()) {
			left.Detach()
			right.Detach()
			var newGetProp *rhino.Node
			if n.IsGetElem() {
				newGetProp = rhino.IR_Getprop(left, right.GetString())
			} else if n.IsOptionalChainStart() {
				newGetProp = rhino.IR_StartOptChainGetprop(left, right.GetString())
			} else {
				newGetProp = rhino.IR_ContinueOptChainGetprop(left, right.GetString())
			}
			n.ReplaceWith(newGetProp)
			c.compiler.ReportChangeToEnclosingScope(newGetProp)
		}
	default:
	}
}
