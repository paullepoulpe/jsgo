// Port of com/google/javascript/jscomp/ClosureOptimizePrimitives.java @ closure-compiler v20261104

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/dtoa"
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
	compiler        AbstractCompiler
	canUseEs6Syntax bool
}

var ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER = DiagnosticType_Warning(
	"JSC_DUPLICATE_SET_MEMBER",
	"Found duplicate value ''{0}'' in set",
)

var (
	googObjectCreateQName                = rhino.QualifiedName_Of("goog.object.create")
	googObjectCreateSetQName             = rhino.QualifiedName_Of("goog.object.createSet")
	closureOptimizePrimitivesVisitTokens = makeVisitTokenSet(rhino.Token_CALL)
)

// closurePrimitiveFinder identifies all calls to closure primitive functions.
type closurePrimitiveFinder struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	outer *ClosureOptimizePrimitives
}

func (*closurePrimitiveFinder) noVisitLeaves() {}

func (*closurePrimitiveFinder) visitTokens() *ast.TokenSet {
	return &closureOptimizePrimitivesVisitTokens
}

func (f *closurePrimitiveFinder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	call, ok := ast.AsCall(n)
	if !ok {
		return
	}
	fn := call.Callee().Node()
	// TODO(user): Once goog.object becomes a goog.module, remove "goog$object$create" and
	// related checks.
	if f.outer.compiler.GetCodingConvention().IsPropertyRenameFunction(fn) {
		f.outer.processRenamePropertyCall(call)
	} else if fn.MatchesName("goog$object$create") ||
		fn.MatchesName("module$contents$goog$object_create") ||
		googObjectCreateQName.Matches(fn) {
		f.outer.processObjectCreateCall(call)
	} else if fn.MatchesName("goog$object$createSet") ||
		fn.MatchesName("module$contents$goog$object_createSet") ||
		googObjectCreateSetQName.Matches(fn) {
		f.outer.processObjectCreateSetCall(call)
	}
}

// NewClosureOptimizePrimitives constructs the ClosureOptimizePrimitives pass.
//
// Java: ClosureOptimizePrimitives(AbstractCompiler, boolean)
func NewClosureOptimizePrimitives(compiler AbstractCompiler, canUseEs6Syntax bool) *ClosureOptimizePrimitives {
	return &ClosureOptimizePrimitives{
		compiler:        compiler,
		canUseEs6Syntax: canUseEs6Syntax,
	}
}

// Java: ClosureOptimizePrimitives.process(Node, Node)
func (c *ClosureOptimizePrimitives) Process(externs rhino.Node, root rhino.Node) {
	pass := closurePrimitiveFinder{outer: c}
	NodeTraversal_Traverse(c.compiler, root, &pass)
}

// processObjectCreateCall converts safe goog.object.create calls to object literals.
func (c *ClosureOptimizePrimitives) processObjectCreateCall(call ast.Call) {
	if !c.canOptimizeObjectCreate(call) {
		return
	}
	callNode := call.Node()
	objNode := rhino.IR_Objectlit().Srcref(callNode)
	scriptNode := NodeUtil_GetEnclosingScript(callNode)
	curParam, hasParam := call.FirstArg()
	for hasParam {
		keyNode := curParam.Node()
		valArg, _ := curParam.Next()
		valueNode := valArg.Node()
		curParam, hasParam = valArg.Next()

		keyNode.Detach()
		valueNode.Detach()

		c.addKeyValueToObjLit(objNode, keyNode, valueNode, scriptNode)
	}
	callNode.ReplaceWith(objNode)
	c.compiler.ReportChangeToEnclosingScope(objNode)
}

// processRenamePropertyCall rewrites property-rename calls to JSCompiler_renameProperty.
func (c *ClosureOptimizePrimitives) processRenamePropertyCall(call ast.Call) {
	nameNode := call.Callee().Node()
	if nameNode.MatchesQualifiedName(NodeUtil_JSC_PROPERTY_NAME_FN) {
		return
	}

	newTarget := rhino.IR_Name(NodeUtil_JSC_PROPERTY_NAME_FN).Srcref(nameNode)
	if orig, ok := nameNode.OriginalQualifiedName(); ok {
		newTarget.SetOriginalName(orig)
	} else {
		newTarget.ClearOriginalName()
	}

	call.SetCallee(ast.MustExpr(newTarget))
	call.SetFreeCall(true)
	c.compiler.ReportChangeToEnclosingScope(call.Node())
}

// canOptimizeObjectCreate reports whether call to goog.object.create can be
// converted to an object literal.
func (c *ClosureOptimizePrimitives) canOptimizeObjectCreate(call ast.Call) bool {
	curParam, hasParam := call.FirstArg()
	for hasParam {
		if !c.isOptimizableKey(curParam.Node()) {
			return false
		}
		// Check for an odd number of parameters.
		valArg, hasVal := curParam.Next()
		if !hasVal {
			return false
		}
		curParam, hasParam = valArg.Next()
	}
	return true
}

// processObjectCreateSetCall converts safe goog.object.createSet calls to object literals.
func (c *ClosureOptimizePrimitives) processObjectCreateSetCall(call ast.Call) {
	if !c.canOptimizeObjectCreateSet(call) {
		return
	}
	callNode := call.Node()
	objNode := rhino.IR_Objectlit().Srcref(callNode)
	scriptNode := NodeUtil_GetEnclosingScript(callNode)
	for arg := range ast.ReadAhead(call.FirstArg()) {
		keyNode := arg.Detach().Node()
		valueNode := rhino.IR_TrueNode().Srcref(keyNode)
		c.addKeyValueToObjLit(objNode, keyNode, valueNode, scriptNode)
	}
	callNode.ReplaceWith(objNode)
	c.compiler.ReportChangeToEnclosingScope(objNode)
}

// canOptimizeObjectCreateSet reports whether call to goog.object.createSet can
// be converted to an object literal.
func (c *ClosureOptimizePrimitives) canOptimizeObjectCreateSet(call ast.Call) bool {
	if only, ok := call.OnlyArg(); ok && !only.Node().IsNumber() && !only.Node().IsStringLit() {
		// If there is only one argument, and it's an array, then the method uses the array elements
		// as keys. Don't optimize it to {[arr]: true}. We only special-case number and string
		// arguments in order to not regress ES5-out behavior.
		return false
	}

	keys := make(map[string]bool)
	for arg := range call.Args() {
		argNode := arg.Node()
		// All keys must be strings or numbers (or ES6 expressions), otherwise we can't optimize.
		if !c.isOptimizableKey(argNode) {
			return false
		}
		var key string
		var hasKey bool
		if str, ok := ast.AsStringLit(argNode); ok {
			key, hasKey = str.Value(), true
		} else if num, ok := ast.AsNumber(argNode); ok {
			key, hasKey = dtoa.DToA_NumberToString(num.Value()), true
		}
		if hasKey {
			if keys[key] {
				c.compiler.Report(JSError_MakeWithNode(call.Callee().Node(), ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER, key))
				return false
			}
			keys[key] = true
		}
	}
	return true
}

func (c *ClosureOptimizePrimitives) addKeyValueToObjLit(objNode, keyNode, valueNode, scriptNode rhino.Node) {
	if keyNode.IsNumber() || keyNode.IsStringLit() {
		if num, ok := ast.AsNumber(keyNode); ok {
			keyNode = rhino.IR_String(dtoa.DToA_NumberToString(num.Value())).Srcref(keyNode)
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

func (c *ClosureOptimizePrimitives) isOptimizableKey(curParam rhino.Node) bool {
	if c.canUseEs6Syntax {
		return !NodeUtil_IsStatement(curParam)
	}
	// Not ES6, all keys must be strings or numbers.
	return curParam.IsStringLit() || curParam.IsNumber()
}
