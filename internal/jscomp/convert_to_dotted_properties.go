// Port of com/google/javascript/jscomp/ConvertToDottedProperties.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast): a handler per
// node kind instead of one Visit that switches on the token.

// ConvertToDottedProperties converts property accesses from quoted string or
// bracket access syntax to dot or unquoted string syntax, where possible. Dot
// syntax is more compact.
type ConvertToDottedProperties struct {
	compiler AbstractCompiler
}

// Java: ConvertToDottedProperties(AbstractCompiler)
func NewConvertToDottedProperties(compiler AbstractCompiler) *ConvertToDottedProperties {
	c := &ConvertToDottedProperties{}
	c.compiler = compiler
	return c
}

// Java: ConvertToDottedProperties.process(Node, Node)
func (c *ConvertToDottedProperties) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_TraverseVisitor(c.compiler, root, &ast.Visitor[*NodeTraversal]{
		ComputedProp:     c.visitComputedProp,
		ComputedFieldDef: c.visitComputedFieldDef,
		GetterDef:        func(_ *NodeTraversal, x ast.GetterDef) { c.unquote(x.Node(), x.Key()) },
		SetterDef:        func(_ *NodeTraversal, x ast.SetterDef) { c.unquote(x.Node(), x.Key()) },
		StringKey:        func(_ *NodeTraversal, x ast.StringKey) { c.unquote(x.Node(), x.Key()) },
		GetElem: func(_ *NodeTraversal, x ast.GetElem) {
			c.convertGetElem(x.Node(), x.Object(), x.Index(), rhino.IR_Getprop)
		},
		OptChainGetElem: func(_ *NodeTraversal, x ast.OptChainGetElem) {
			getProp := rhino.IR_ContinueOptChainGetprop
			if x.Node().IsOptionalChainStart() {
				getProp = rhino.IR_StartOptChainGetprop
			}
			c.convertGetElem(x.Node(), x.Object(), x.Index(), getProp)
		},
	})
	GatherGetterAndSetterProperties_Update(c.compiler, externs, root)
}

// dottableKey returns the name of a computed key that can be written as a
// plain key: ['x'] can become x. Not ['constructor'], since
// ['constructor']() and constructor() are different, nor, in an object
// literal, ['__proto__'], since {['__proto__']: v} defines an own property
// and {__proto__: v} sets the prototype.
func dottableKey(key ast.Expr, parent rhino.Node) (string, bool) {
	lit, ok := key.AsStringLit()
	if !ok {
		return "", false
	}
	name := lit.Value()
	if !NodeUtil_IsValidPropertyName(parser.FeatureSet_ES3, name) ||
		name == "constructor" ||
		(parent.IsObjectLit() && name == "__proto__") {
		return "", false
	}
	return name, true
}

// visitComputedProp turns ['x']: v into x: v, and a method, getter or
// setter with a computed key into one with a plain key.
//
// Java: ConvertToDottedProperties.visit(NodeTraversal, Node, Node) (COMPUTED_PROP)
func (c *ConvertToDottedProperties) visitComputedProp(t *NodeTraversal, x ast.ComputedProp) {
	n := x.Node()
	name, ok := dottableKey(x.Key(), n.GetParent())
	if !ok {
		return
	}
	value, hasValue := x.Value() // read before the key is detached
	x.Key().Node().Detach()
	if hasValue {
		value.Node().Detach()
	}

	var temp rhino.Node
	fn, isFunction := value.AsFunction()
	switch {
	case isFunction && n.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER):
		temp = rhino.IR_GetterDef(name, fn.Node())
	case isFunction && n.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER):
		temp = rhino.IR_SetterDef(name, fn.Node())
	case isFunction && n.GetParent().IsClassMembers():
		temp = rhino.IR_MemberFunctionDef(name, fn.Node())
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_MEMBER_DECLARATIONS, c.compiler)
	default:
		// TODO - further optimize a function value to a member function
		temp = rhino.IR_StringKeyWithValue(name, value.Node())
	}
	c.replaceKey(n, temp)
}

// visitComputedFieldDef turns the class field ['x'] = v into x = v.
//
// Java: ConvertToDottedProperties.visit(NodeTraversal, Node, Node) (COMPUTED_FIELD_DEF)
func (c *ConvertToDottedProperties) visitComputedFieldDef(_ *NodeTraversal, x ast.ComputedFieldDef) {
	n := x.Node()
	name, ok := dottableKey(x.Key(), n.GetParent())
	if !ok {
		return
	}
	value, hasValue := x.Value() // read before the key is detached
	x.Key().Node().Detach()
	var temp rhino.Node
	if hasValue {
		value.Node().Detach()
		temp = rhino.IR_MemberFieldDef(name, value.Node())
	} else {
		// An uninitialized field, e.g. `class C { ['x']; }`, becomes a 0-child
		// MEMBER_FIELD_DEF, matching how the parser represents `class C { x; }`.
		temp = rhino.Node_NewStringWithToken(rhino.Token_MEMBER_FIELD_DEF, name)
	}
	c.replaceKey(n, temp)
}

// replaceKey replaces the member n with temp.
func (c *ConvertToDottedProperties) replaceKey(n, temp rhino.Node) {
	temp.SetStaticMember(n.IsStaticMember())
	n.ReplaceWith(temp)
	c.compiler.ReportChangeToEnclosingScope(temp)
}

// unquote turns 'x': v into x: v (and the same for getters and setters).
//
// Java: ConvertToDottedProperties.visit(NodeTraversal, Node, Node) (GETTER_DEF, SETTER_DEF, STRING_KEY)
func (c *ConvertToDottedProperties) unquote(n rhino.Node, key string) {
	if NodeUtil_IsValidPropertyName(parser.FeatureSet_ES3, key) && n.GetBooleanProp(rhino.Node_QUOTED_PROP) {
		n.PutBooleanProp(rhino.Node_QUOTED_PROP, false)
		c.compiler.ReportChangeToEnclosingScope(n)
	}
}

// convertGetElem turns object['x'] (n) into object.x, made by getProp.
//
// Java: ConvertToDottedProperties.visit(NodeTraversal, Node, Node) (GETELEM, OPTCHAIN_GETELEM)
func (c *ConvertToDottedProperties) convertGetElem(n rhino.Node, object, index ast.Expr,
	getProp func(object rhino.Node, name string) rhino.Node) {
	lit, ok := index.AsStringLit()
	if !ok || !NodeUtil_IsValidPropertyName(parser.FeatureSet_ES3, lit.Value()) {
		return
	}
	object.Node().Detach()
	lit.Node().Detach()
	newGetProp := getProp(object.Node(), lit.Value())
	n.ReplaceWith(newGetProp)
	c.compiler.ReportChangeToEnclosingScope(newGetProp)
}
