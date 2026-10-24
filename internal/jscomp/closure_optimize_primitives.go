// Port of com/google/javascript/jscomp/ClosureOptimizePrimitives.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/dtoa"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ClosureOptimizePrimitives is a compiler pass that converts primitive
// calls:
//
// Converts goog.object.create(key1, val1, key2, val2, ...) where all of the
// keys are literals into object literals.
//
// Converts goog.object.createSet(key1, key2, ...) into an object literal
// with the given keys, where all the values are true.
//
// Converts goog.reflect.objectProperty(propName, object) to
// JSCompiler_renameProperty
type ClosureOptimizePrimitives struct {
	// Reference to the JS compiler
	compiler AbstractCompiler

	// Whether we can use Es6 syntax
	canUseEs6Syntax bool
}

var ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER = DiagnosticType_Warning("JSC_DUPLICATE_SET_MEMBER", "Found duplicate value ''{0}'' in set")

var closureOptimizePrimitives_GOOG_OBJECT_CREATE = rhino.QualifiedName_Of("goog.object.create")
var closureOptimizePrimitives_GOOG_OBJECT_CREATESET = rhino.QualifiedName_Of("goog.object.createSet")

// ClosureOptimizePrimitives_FindPrimitives identifies all calls to closure
// primitive functions
type ClosureOptimizePrimitives_FindPrimitives struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	outer *ClosureOptimizePrimitives
}

// Java: ClosureOptimizePrimitives.FindPrimitives()
func newClosureOptimizePrimitives_FindPrimitives(outer *ClosureOptimizePrimitives) *ClosureOptimizePrimitives_FindPrimitives {
	return &ClosureOptimizePrimitives_FindPrimitives{outer: outer}
}

// Java: ClosureOptimizePrimitives.FindPrimitives.visit(NodeTraversal, Node, Node)
func (f *ClosureOptimizePrimitives_FindPrimitives) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsCall() {
		fn := n.GetFirstChild()
		// TODO(user): Once goog.object becomes a goog.module, remove "goog$object$create" and
		// related checks.
		if f.outer.compiler.GetCodingConvention().IsPropertyRenameFunction(fn) {
			f.outer.ProcessRenamePropertyCall(n)
		} else if fn.MatchesName("goog$object$create") ||
			fn.MatchesName("module$contents$goog$object_create") ||
			closureOptimizePrimitives_GOOG_OBJECT_CREATE.Matches(fn) {
			f.outer.ProcessObjectCreateCall(n)
		} else if fn.MatchesName("goog$object$createSet") ||
			fn.MatchesName("module$contents$goog$object_createSet") ||
			closureOptimizePrimitives_GOOG_OBJECT_CREATESET.Matches(fn) {
			f.outer.ProcessObjectCreateSetCall(n)
		}
	}
}

// NewClosureOptimizePrimitives: compiler is the AbstractCompiler.
//
// Java: ClosureOptimizePrimitives(AbstractCompiler, boolean)
func NewClosureOptimizePrimitives(compiler AbstractCompiler, canUseEs6Syntax bool) *ClosureOptimizePrimitives {
	c := &ClosureOptimizePrimitives{}
	c.compiler = compiler
	c.canUseEs6Syntax = canUseEs6Syntax
	return c
}

// Java: ClosureOptimizePrimitives.process(Node, Node)
func (c *ClosureOptimizePrimitives) Process(externs *rhino.Node, root *rhino.Node) {
	pass := newClosureOptimizePrimitives_FindPrimitives(c)
	NodeTraversal_Traverse(c.compiler, root, pass)
}

// ProcessObjectCreateCall converts all of the given call nodes to object
// literals that are safe to do so.
//
// Java: ClosureOptimizePrimitives.processObjectCreateCall(Node)
func (c *ClosureOptimizePrimitives) ProcessObjectCreateCall(callNode *rhino.Node) {
	curParam := callNode.GetSecondChild()
	if c.CanOptimizeObjectCreate(curParam) {
		objNode := rhino.IR_Objectlit().Srcref(callNode)
		for curParam != nil {
			keyNode := curParam
			valueNode := curParam.GetNext()
			curParam = valueNode.GetNext()

			keyNode.Detach()
			valueNode.Detach()

			c.AddKeyValueToObjLit(objNode, keyNode, valueNode, NodeUtil_GetEnclosingScript(callNode))
		}
		callNode.ReplaceWith(objNode)
		c.compiler.ReportChangeToEnclosingScope(objNode)
	}
}

// ProcessRenamePropertyCall converts all of the given call nodes to object
// literals that are safe to do so.
//
// Java: ClosureOptimizePrimitives.processRenamePropertyCall(Node)
func (c *ClosureOptimizePrimitives) ProcessRenamePropertyCall(callNode *rhino.Node) {
	nameNode := callNode.GetFirstChild()
	if nameNode.MatchesQualifiedName(NodeUtil_JSC_PROPERTY_NAME_FN) {
		return
	}

	newTarget := rhino.IR_Name(NodeUtil_JSC_PROPERTY_NAME_FN).Srcref(nameNode)
	newTarget.SetOriginalName(nameNode.GetOriginalQualifiedName())

	nameNode.ReplaceWith(newTarget)
	callNode.PutBooleanProp(rhino.Node_FREE_CALL, true)
	c.compiler.ReportChangeToEnclosingScope(callNode)
}

// CanOptimizeObjectCreate returns whether the given call to
// goog.object.create can be converted to an object literal.
//
// Java: ClosureOptimizePrimitives.canOptimizeObjectCreate(Node)
func (c *ClosureOptimizePrimitives) CanOptimizeObjectCreate(firstParam *rhino.Node) bool {
	curParam := firstParam
	for curParam != nil {
		if !c.IsOptimizableKey(curParam) {
			return false
		}
		curParam = curParam.GetNext()

		// Check for an odd number of parameters.
		if curParam == nil {
			return false
		}
		curParam = curParam.GetNext()
	}
	return true
}

// ProcessObjectCreateSetCall converts all of the given call nodes to object
// literals that are safe to do so.
//
// Java: ClosureOptimizePrimitives.processObjectCreateSetCall(Node)
func (c *ClosureOptimizePrimitives) ProcessObjectCreateSetCall(callNode *rhino.Node) {
	curParam := callNode.GetSecondChild()
	if c.CanOptimizeObjectCreateSet(curParam) {
		objNode := rhino.IR_Objectlit().Srcref(callNode)
		for curParam != nil {
			keyNode := curParam
			valueNode := rhino.IR_TrueNode().Srcref(keyNode)

			curParam = curParam.GetNext()
			keyNode.Detach()

			c.AddKeyValueToObjLit(objNode, keyNode, valueNode, NodeUtil_GetEnclosingScript(callNode))
		}
		callNode.ReplaceWith(objNode)
		c.compiler.ReportChangeToEnclosingScope(objNode)
	}
}

// CanOptimizeObjectCreateSet returns whether the given call to
// goog.object.createSet can be converted to an object literal.
//
// Java: ClosureOptimizePrimitives.canOptimizeObjectCreateSet(Node)
func (c *ClosureOptimizePrimitives) CanOptimizeObjectCreateSet(firstParam *rhino.Node) bool {
	if firstParam != nil &&
		firstParam.GetNext() == nil &&
		!(firstParam.IsNumber() || firstParam.IsStringLit()) {
		// if there is only one argument, and it's an array, then the method uses the array elements
		// as keys. Don't optimize it to {[arr]: true}. We only special-case number and string
		// arguments in order to not regress ES5-out behavior
		return false
	}

	curParam := firstParam
	keys := jutil.NewLinkedHashSet[string]()
	for curParam != nil {
		// All keys must be strings or numbers, otherwise we can't optimize the call.
		if !c.IsOptimizableKey(curParam) {
			return false
		}
		if curParam.IsStringLit() || curParam.IsNumber() {
			var key string
			if curParam.IsStringLit() {
				key = curParam.GetString()
			} else {
				key = dtoa.DToA_NumberToString(curParam.GetDouble())
			}
			if !keys.Add(key) {
				c.compiler.Report(JSError_MakeWithNode(firstParam.GetPrevious(), ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER, key))
				return false
			}
		}
		curParam = curParam.GetNext()
	}
	return true
}

// Java: ClosureOptimizePrimitives.addKeyValueToObjLit(Node, Node, Node, Node)
func (c *ClosureOptimizePrimitives) AddKeyValueToObjLit(objNode *rhino.Node, keyNode *rhino.Node, valueNode *rhino.Node, scriptNode *rhino.Node) {
	if keyNode.IsNumber() || keyNode.IsStringLit() {
		if keyNode.IsNumber() {
			keyNode = rhino.IR_String(dtoa.DToA_NumberToString(keyNode.GetDouble())).Srcref(keyNode)
		}
		// It isn't valid for a `STRING_KEY` to be marked as parenthesized.
		keyNode.SetIsParenthesized(false)
		keyNode.SetToken(rhino.Token_STRING_KEY)
		keyNode.SetQuotedStringKey()
		objNode.AddChildToBack(rhino.IR_Propdef(keyNode, valueNode))
	} else {
		objNode.AddChildToBack(rhino.IR_ComputedProp(keyNode, valueNode).Srcref(keyNode))
		NodeUtil_AddFeatureToScript(scriptNode, parser.FeatureSet_Feature_COMPUTED_PROPERTIES, c.compiler)
	}
}

// Java: ClosureOptimizePrimitives.isOptimizableKey(Node)
func (c *ClosureOptimizePrimitives) IsOptimizableKey(curParam *rhino.Node) bool {
	if c.canUseEs6Syntax {
		return !NodeUtil_IsStatement(curParam)
	} else {
		// Not ES6, all keys must be strings or numbers.
		return curParam.IsStringLit() || curParam.IsNumber()
	}
}
