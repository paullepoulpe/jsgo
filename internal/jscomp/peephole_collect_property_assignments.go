// Port of com/google/javascript/jscomp/PeepholeCollectPropertyAssignments.java @ closure-compiler v20261021

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
	"math"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// Go-only: this pass uses the typed AST view (internal/ast) to read the
// statement that creates the literal (name.Init(), assign.Target()) and the
// property assignments that follow it (lhs.Object(), elem.Index()...)
// instead of child positions. isPropertyAssignmentToName returns the parts
// of the assignment it recognized, which collectProperty and its helpers
// use instead of reading them again; Java's getName and getValue are one
// function. Code that treats children alike (the statements of a SCRIPT or
// a BLOCK, the members of an object literal grouped by shape, the
// forward-reference walk) and the node builders use the rhino API.

// PeepholeCollectPropertyAssignments is a pass that looks for assignments to
// properties of an object or array immediately following its creation using
// the abbreviated syntax.
//
// E.g. `var a = [];a[0] = 0` is optimized to `var a = [0]` and similarly for
// the object constructor.
type PeepholeCollectPropertyAssignments struct {
	AbstractPeepholeOptimizationBase
}

// Java: PeepholeCollectPropertyAssignments()
func NewPeepholeCollectPropertyAssignments() *PeepholeCollectPropertyAssignments {
	return &PeepholeCollectPropertyAssignments{}
}

var peepholeCollectPropertyAssignmentsTokens = peepholeTokenSet(
	rhino.Token_SCRIPT,
	rhino.Token_BLOCK,
)

func (p *PeepholeCollectPropertyAssignments) peepholeConfig() peepholeOptConfig {
	return peepholeOptConfig{
		tokens:     peepholeCollectPropertyAssignmentsTokens,
		isStandard: true,
	}
}

func (p *PeepholeCollectPropertyAssignments) cloneForWorker() AbstractPeepholeOptimization {
	return NewPeepholeCollectPropertyAssignments()
}

// Java: PeepholeCollectPropertyAssignments.optimizeSubtree(Node)
func (p *PeepholeCollectPropertyAssignments) OptimizeSubtree(subtree rhino.Node) rhino.Node {
	if !subtree.IsScript() && !subtree.IsBlock() {
		return subtree
	}

	codeChanged := false

	// Look for variable declarations or simple assignments
	// and start processing there.
	for child := subtree.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if !NodeUtil_IsNameDeclaration(child) && !NodeUtil_IsExprAssign(child) {
			continue
		}
		if _, ok := PeepholeCollectPropertyAssignments_IsPropertyAssignmentToName(child.GetNext()); !ok {
			// Quick check to see if there's anything to collapse.
			continue
		}

		name, value, ok := PeepholeCollectPropertyAssignments_GetNameAndValue(child)
		if !ok {
			// The assignment target is not a simple name, or there is no
			// initializer.
			continue
		}
		if !PeepholeCollectPropertyAssignments_IsInterestingValue(value) {
			// Not an Object or Array literal.
			continue
		}

		identID := name.Node().GetStringID()
		for propertyCandidate := child.GetNext(); !propertyCandidate.IsNil(); propertyCandidate = child.GetNext() {
			// This does not infinitely loop because collectProperty always
			// removes propertyCandidate from its parent when it returns true.
			if !p.collectPropertyByID(propertyCandidate, identID, value) {
				break
			}
			codeChanged = true
		}
	}

	if codeChanged {
		p.ReportChangeToEnclosingScope(subtree)
	}
	return subtree
}

// PeepholeCollectPropertyAssignments_GetNameAndValue returns the name and
// the value of n, a declaration `var name = value` (with one declarator) or
// an assignment `name = value;`. ok is false if what n declares or assigns
// is not a simple name (a pattern, a property), or if the name is declared
// without a value.
//
// Go-only: Java's getName and getValue, as one function; getValue was only
// called once getName had returned a name.
//
// Java: PeepholeCollectPropertyAssignments.getName(Node)
// Java: PeepholeCollectPropertyAssignments.getValue(Node)
func PeepholeCollectPropertyAssignments_GetNameAndValue(n rhino.Node) (ast.Name, ast.Expr, bool) {
	if decl, ok := ast.AsNameDeclaration(n); ok {
		declarator, ok := decl.OnlyDeclarator()
		jutil.CheckState(ok)
		name, ok := declarator.AsName()
		if !ok {
			return ast.Name{}, ast.Expr{}, false
		}
		value, ok := name.Init()
		return name, value, ok
	}
	if stmt, ok := ast.AsExprStmt(n); ok {
		if assign, ok := stmt.Expr().AsAssign(); ok {
			name, ok := assign.Target().AsName()
			return name, assign.Value(), ok
		}
	}
	panic(jutil.NewIllegalStateException())
}

// Java: PeepholeCollectPropertyAssignments.isInterestingValue(Node)
func PeepholeCollectPropertyAssignments_IsInterestingValue(n ast.Expr) bool {
	return n.Node().IsObjectLit() || n.Node().IsArrayLit()
}

// peepholeCollectPropertyAssignments_propertyAssignment is the statement
// `obj.property = rhs;` or `obj[index] = rhs;`, where obj is a variable:
// what isPropertyAssignmentToName recognizes.
//
// Go-only: Java reads these parts again from the statement where it needs
// them.
type peepholeCollectPropertyAssignments_propertyAssignment struct {
	stmt ast.ExprStmt  // the whole statement
	lhs  ast.NormalGet // obj.property or obj[index]
	obj  ast.Name      // the variable whose property is assigned
	rhs  ast.Expr      // the assigned value
}

// PeepholeCollectPropertyAssignments_IsPropertyAssignmentToName returns the
// parts of propertyCandidate if it is an assignment to a property of a
// variable; ok is false otherwise (also when propertyCandidate is the nil
// node).
//
// Java: PeepholeCollectPropertyAssignments.isPropertyAssignmentToName(Node)
func PeepholeCollectPropertyAssignments_IsPropertyAssignmentToName(propertyCandidate rhino.Node) (peepholeCollectPropertyAssignments_propertyAssignment, bool) {
	var none peepholeCollectPropertyAssignments_propertyAssignment
	// Must be an assignment...
	stmt, ok := ast.AsExprStmt(propertyCandidate)
	if !ok {
		return none, false
	}
	assign, ok := stmt.Expr().AsAssign()
	if !ok {
		return none, false
	}

	// to a property...
	lhs, ok := assign.Target().AsNormalGet()
	if !ok {
		return none, false
	}

	// of a variable.
	obj, ok := lhs.Object().AsName()
	if !ok {
		return none, false
	}
	return peepholeCollectPropertyAssignments_propertyAssignment{
		stmt: stmt, lhs: lhs, obj: obj, rhs: assign.Value(),
	}, true
}

// Java: PeepholeCollectPropertyAssignments.collectProperty(Node, String, Node)
func (p *PeepholeCollectPropertyAssignments) CollectProperty(propertyCandidate rhino.Node, name string, value ast.Expr) bool {
	id := uint32(strpool.Lookup(name))
	if id == 0 {
		return false
	}
	return p.collectPropertyByID(propertyCandidate, id, value)
}

func (p *PeepholeCollectPropertyAssignments) collectPropertyByID(propertyCandidate rhino.Node, nameID uint32, value ast.Expr) bool {
	prop, ok := PeepholeCollectPropertyAssignments_IsPropertyAssignmentToName(propertyCandidate)
	if !ok {
		return false
	}

	// Must be an assignment to the recent variable...
	if nameID != prop.obj.Node().GetStringID() {
		return false
	}

	rhs := prop.rhs.Node()
	// with a value that cannot change the values of the variables,
	if p.MayHaveSideEffects(rhs) ||
		NodeUtil_CanBeSideEffected(rhs) {
		return false
	}
	// and does not have a reference to a variable initialized after it.
	if !NodeUtil_IsLiteralValue(rhs, true) &&
		peepholeCollectPropertyAssignments_mightContainForwardReferenceByID(rhs, nameID) {
		return false
	}

	if arrayLiteral, ok := value.AsArrayLit(); ok {
		return PeepholeCollectPropertyAssignments_CollectArrayProperty(arrayLiteral, prop)
	}
	if objectLiteral, ok := value.AsObjectLit(); ok {
		return p.CollectObjectProperty(objectLiteral, prop)
	}
	panic(jutil.NewIllegalStateException())
}

// Java: PeepholeCollectPropertyAssignments.collectArrayProperty(Node, Node)
func PeepholeCollectPropertyAssignments_CollectArrayProperty(arrayLiteral ast.ArrayLit, prop peepholeCollectPropertyAssignments_propertyAssignment) bool {
	sizeOfArrayAtStart := arrayLiteral.Node().GetChildCount()
	maxIndexAssigned := sizeOfArrayAtStart - 1

	elem, ok := prop.lhs.AsGetElem()
	if !ok {
		return false
	}
	// The left hand side must have a numeric index
	property, ok := elem.Index().AsNumber()
	if !ok {
		return false
	}
	// that is a valid array index
	dindex := property.Value()
	if !(dindex >= 0) || // Handles NaN and negatives.
		math.IsInf(dindex, 0) || dindex > 0x7fffffff {
		return false
	}
	index := int(jutil.D2I(dindex))
	if dindex != float64(index) {
		return false
	}
	// that would not make the array so sparse that they take more space
	// when rendered than x[9]=1.
	if maxIndexAssigned+4 < index {
		return false
	}
	// If the array literal contains a spread element, its runtime length and element
	// offsets are unknown, so subsequent indexed property assignments cannot safely
	// be appended to or folded into the literal.
	if PeepholeCollectPropertyAssignments_ArrayLiteralHasSpread(arrayLiteral) {
		return false
	}
	if index > maxIndexAssigned {
		for maxIndexAssigned < index-1 {
			// Pad the array if it is sparse.
			// So if array is [0] and integer 3 is assigned at index is 2, then
			// we want to produce [0,,2].
			emptyNode := rhino.IR_Empty().Srcref(arrayLiteral.Node())
			arrayLiteral.Node().AddChildToBack(emptyNode)
			maxIndexAssigned++
		}
		arrayLiteral.AppendElement(prop.rhs.Detach().AsArrayElem())
	} else {
		// An out of order assignment.  Allow it if it's a hole.
		currentValue := arrayLiteral.Node().GetChildAtIndex(index)
		if !currentValue.IsEmpty() {
			// We've already collected a value for this index.
			return false
		}
		currentValue.ReplaceWith(prop.rhs.Detach().Node())
	}

	prop.stmt.Detach()
	return true
}

// Java: PeepholeCollectPropertyAssignments.arrayLiteralHasSpread(Node)
func PeepholeCollectPropertyAssignments_ArrayLiteralHasSpread(arrayLiteral ast.ArrayLit) bool {
	script := NodeUtil_GetEnclosingScript(arrayLiteral.Node())
	var scriptFeatures *parser.FeatureSet
	if !script.IsNil() {
		scriptFeatures = NodeUtil_GetFeatureSetOfScript(script)
	}
	if scriptFeatures != nil && !scriptFeatures.Has(parser.FeatureSet_Feature_SPREAD_EXPRESSIONS) {
		return false
	}
	for elem := range arrayLiteral.Elements() {
		if _, ok := elem.AsIterSpread(); ok {
			return true
		}
	}
	return false
}

// Java: PeepholeCollectPropertyAssignments.collectObjectProperty(Node, Node)
func (p *PeepholeCollectPropertyAssignments) CollectObjectProperty(objectLiteral ast.ObjectLit, prop peepholeCollectPropertyAssignments_propertyAssignment) bool {
	rhs := prop.rhs.Node()

	// The property must be statically known.
	// Go-only: Java checks that the lhs is a GETPROP or a GETELEM, and returns
	// false for anything else; prop.lhs is one of them.
	var propertyName string
	var propertyNameNode rhino.Node
	var propertyNameID uint32
	var hasPropertyName bool
	elem, isGetElem := prop.lhs.AsGetElem()
	if isGetElem {
		property := elem.Index()
		if number, ok := property.AsNumber(); ok {
			propertyName, hasPropertyName = p.GetSideEffectFreeStringValue(number.Node())
			if hasPropertyName {
				propertyNameID = uint32(strpool.Lookup(propertyName))
			}
		} else if str, ok := property.AsStringLit(); ok {
			propertyNameNode = str.Node()
			propertyNameID, hasPropertyName = propertyNameNode.GetStringID(), true
		} else {
			return false
		}
	} else if get, ok := prop.lhs.AsGetProp(); ok {
		propertyNameNode = get.Node()
		propertyNameID, hasPropertyName = propertyNameNode.GetStringID(), true
	}

	// Check if the new property already exists in the object literal
	// Note: Duplicate keys are invalid in strict mode
	var existingProperty rhino.Node
	for currentProperty := objectLiteral.Node().GetFirstChild(); !currentProperty.IsNil(); currentProperty = currentProperty.GetNext() {
		if currentProperty.IsStringKey() || currentProperty.IsMemberFunctionDef() {
			currentValue := currentProperty.GetFirstChild()
			// Compare the current property name with the new property name
			if hasPropertyName && propertyNameID != 0 && currentProperty.GetStringID() == propertyNameID {
				existingProperty = currentProperty
				// Check if the current value and the new value are side-effect
				isCurrentValueSideEffect := NodeUtil_CanBeSideEffected(currentValue)
				isNewValueSideEffect := NodeUtil_CanBeSideEffected(rhs)
				// If they are side-effect free then replace the current value with the new one
				if isCurrentValueSideEffect || isNewValueSideEffect {
					return false
				}
				// Break the loop if the property exists
				break
			}
		} else if currentProperty.IsGetterDef() || currentProperty.IsSetterDef() {
			if hasPropertyName && propertyNameID != 0 && currentProperty.GetStringID() == propertyNameID {
				return false
			}
		} else if computed, ok := ast.AsComputedProp(currentProperty); ok {
			key, ok := computed.Key().AsStringLit()
			if !ok || (hasPropertyName && propertyNameID != 0 && key.Node().GetStringID() == propertyNameID) {
				return false
			}
		}
	}

	if !propertyNameNode.IsNil() {
		propertyName = propertyNameNode.GetString()
	}
	newProperty := rhino.IR_StringKey(propertyName).SrcrefIfMissing(prop.stmt.Node())
	// Preserve the quotedness of a property reference
	if isGetElem {
		newProperty.SetQuotedStringKey()
	}
	newValue := rhs.Detach()
	newProperty.AddChildToBack(newValue)

	if !existingProperty.IsNil() {
		p.DeleteNode(existingProperty)
	}
	// If the property does not already exist we can safely add it
	objectLiteral.Node().AddChildToBack(newProperty)
	prop.stmt.Detach()
	return true
}

// Java: PeepholeCollectPropertyAssignments.mightContainForwardReference(Node, String)
func PeepholeCollectPropertyAssignments_MightContainForwardReference(node rhino.Node, varName string) bool {
	id := uint32(strpool.Lookup(varName))
	if id == 0 {
		return false
	}
	return peepholeCollectPropertyAssignments_mightContainForwardReferenceByID(node, id)
}

func peepholeCollectPropertyAssignments_mightContainForwardReferenceByID(node rhino.Node, varNameID uint32) bool {
	if node.IsName() {
		return varNameID == node.GetStringID()
	}
	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if peepholeCollectPropertyAssignments_mightContainForwardReferenceByID(child, varNameID) {
			return true
		}
	}
	return false
}
