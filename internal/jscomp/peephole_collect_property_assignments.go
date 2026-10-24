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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	p := &PeepholeCollectPropertyAssignments{}
	p.InitAbstractPeepholeOptimizationBase(p) // the implicit super()
	return p
}

// Java: PeepholeCollectPropertyAssignments.optimizeSubtree(Node)
func (p *PeepholeCollectPropertyAssignments) OptimizeSubtree(subtree *rhino.Node) *rhino.Node {
	if !subtree.IsScript() && !subtree.IsBlock() {
		return subtree
	}

	codeChanged := false

	// Look for variable declarations or simple assignments
	// and start processing there.
	for child := subtree.GetFirstChild(); child != nil; child = child.GetNext() {
		if !NodeUtil_IsNameDeclaration(child) && !NodeUtil_IsExprAssign(child) {
			continue
		}
		if !PeepholeCollectPropertyAssignments_IsPropertyAssignmentToName(child.GetNext()) {
			// Quick check to see if there's anything to collapse.
			continue
		}

		jutil.CheckState(child.HasOneChild())
		name := PeepholeCollectPropertyAssignments_GetName(child)
		if !name.IsName() {
			// The assignment target is not a simple name.
			continue
		}
		value := PeepholeCollectPropertyAssignments_GetValue(child)
		if value == nil || !PeepholeCollectPropertyAssignments_IsInterestingValue(value) {
			// No initializer or not an Object or Array literal.
			continue
		}

		var propertyCandidate *rhino.Node
		for {
			propertyCandidate = child.GetNext()
			if propertyCandidate == nil {
				break
			}
			// This does not infinitely loop because collectProperty always
			// removes propertyCandidate from its parent when it returns true.
			if !p.CollectProperty(propertyCandidate, name.GetString(), value) {
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

// Java: PeepholeCollectPropertyAssignments.getName(Node)
func PeepholeCollectPropertyAssignments_GetName(n *rhino.Node) *rhino.Node {
	if NodeUtil_IsNameDeclaration(n) {
		return n.GetFirstChild()
	} else if NodeUtil_IsExprAssign(n) {
		return n.GetFirstFirstChild()
	}
	panic(jutil.NewIllegalStateException())
}

// Java: PeepholeCollectPropertyAssignments.getValue(Node)
func PeepholeCollectPropertyAssignments_GetValue(n *rhino.Node) *rhino.Node {
	if NodeUtil_IsNameDeclaration(n) {
		return n.GetFirstFirstChild()
	} else if NodeUtil_IsExprAssign(n) {
		return n.GetFirstChild().GetLastChild()
	}
	panic(jutil.NewIllegalStateException())
}

// Java: PeepholeCollectPropertyAssignments.isInterestingValue(Node)
func PeepholeCollectPropertyAssignments_IsInterestingValue(n *rhino.Node) bool {
	return n.IsObjectLit() || n.IsArrayLit()
}

// Java: PeepholeCollectPropertyAssignments.isPropertyAssignmentToName(Node)
func PeepholeCollectPropertyAssignments_IsPropertyAssignmentToName(propertyCandidate *rhino.Node) bool {
	if propertyCandidate == nil {
		return false
	}
	// Must be an assignment...
	if !NodeUtil_IsExprAssign(propertyCandidate) {
		return false
	}

	expr := propertyCandidate.GetFirstChild()

	// to a property...
	lhs := expr.GetFirstChild()
	if !NodeUtil_IsNormalGet(lhs) {
		return false
	}

	// of a variable.
	obj := lhs.GetFirstChild()
	return obj.IsName()
}

// Java: PeepholeCollectPropertyAssignments.collectProperty(Node, String, Node)
func (p *PeepholeCollectPropertyAssignments) CollectProperty(propertyCandidate *rhino.Node, name string, value *rhino.Node) bool {
	if !PeepholeCollectPropertyAssignments_IsPropertyAssignmentToName(propertyCandidate) {
		return false
	}

	lhs := propertyCandidate.GetFirstFirstChild()
	// Must be an assignment to the recent variable...
	if name != lhs.GetFirstChild().GetString() {
		return false
	}

	rhs := lhs.GetNext()
	// with a value that cannot change the values of the variables,
	if p.MayHaveSideEffects(rhs) ||
		NodeUtil_CanBeSideEffected(rhs) {
		return false
	}
	// and does not have a reference to a variable initialized after it.
	if !NodeUtil_IsLiteralValue(rhs, true) &&
		PeepholeCollectPropertyAssignments_MightContainForwardReference(rhs, name) {
		return false
	}

	switch value.GetToken() {
	case rhino.Token_ARRAYLIT:
		if !PeepholeCollectPropertyAssignments_CollectArrayProperty(value, propertyCandidate) {
			return false
		}
	case rhino.Token_OBJECTLIT:
		if !p.CollectObjectProperty(value, propertyCandidate) {
			return false
		}
	default:
		panic(jutil.NewIllegalStateException())
	}
	return true
}

// Java: PeepholeCollectPropertyAssignments.collectArrayProperty(Node, Node)
func PeepholeCollectPropertyAssignments_CollectArrayProperty(arrayLiteral *rhino.Node, propertyCandidate *rhino.Node) bool {
	assignment := propertyCandidate.GetFirstChild()
	sizeOfArrayAtStart := arrayLiteral.GetChildCount()
	maxIndexAssigned := sizeOfArrayAtStart - 1

	lhs := assignment.GetFirstChild()
	rhs := lhs.GetNext()
	if !lhs.IsGetElem() {
		return false
	}
	obj := lhs.GetFirstChild()
	property := obj.GetNext()
	// The left hand side must have a numeric index
	if !property.IsNumber() {
		return false
	}
	// that is a valid array index
	dindex := property.GetDouble()
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
			emptyNode := rhino.IR_Empty().Srcref(arrayLiteral)
			arrayLiteral.AddChildToBack(emptyNode)
			maxIndexAssigned++
		}
		arrayLiteral.AddChildToBack(rhs.Detach())
	} else {
		// An out of order assignment.  Allow it if it's a hole.
		currentValue := arrayLiteral.GetChildAtIndex(index)
		if !currentValue.IsEmpty() {
			// We've already collected a value for this index.
			return false
		}
		currentValue.ReplaceWith(rhs.Detach())
	}

	propertyCandidate.Detach()
	return true
}

// Java: PeepholeCollectPropertyAssignments.arrayLiteralHasSpread(Node)
func PeepholeCollectPropertyAssignments_ArrayLiteralHasSpread(arrayLiteral *rhino.Node) bool {
	script := NodeUtil_GetEnclosingScript(arrayLiteral)
	var scriptFeatures *parser.FeatureSet
	if script != nil {
		scriptFeatures = NodeUtil_GetFeatureSetOfScript(script)
	}
	if scriptFeatures != nil && !scriptFeatures.Has(parser.FeatureSet_Feature_SPREAD_EXPRESSIONS) {
		return false
	}
	for child := arrayLiteral.GetFirstChild(); child != nil; child = child.GetNext() {
		if child.IsSpread() {
			return true
		}
	}
	return false
}

// Java: PeepholeCollectPropertyAssignments.collectObjectProperty(Node, Node)
func (p *PeepholeCollectPropertyAssignments) CollectObjectProperty(objectLiteral *rhino.Node, propertyCandidate *rhino.Node) bool {
	assignment := propertyCandidate.GetFirstChild()
	lhs := assignment.GetFirstChild()
	rhs := lhs.GetNext()
	obj := lhs.GetFirstChild()
	jutil.CheckStateMsg1(lhs.IsGetProp() || lhs.IsGetElem(), "%s", lhs)

	// The property must be statically known.
	var propertyName *string
	if lhs.IsGetElem() {
		property := obj.GetNext()
		if property.IsNumber() {
			propertyName = p.GetSideEffectFreeStringValue(property)
		} else if property.IsStringLit() {
			s := property.GetString()
			propertyName = &s
		} else {
			return false
		}
	} else if lhs.IsGetProp() {
		s := lhs.GetString()
		propertyName = &s
	} else {
		return false
	}

	// Check if the new property already exists in the object literal
	// Note: Duplicate keys are invalid in strict mode
	var existingProperty *rhino.Node
	for currentProperty := objectLiteral.GetFirstChild(); currentProperty != nil; currentProperty = currentProperty.GetNext() {
		if currentProperty.IsStringKey() || currentProperty.IsMemberFunctionDef() {
			// Get the name of the current property
			currentPropertyName := currentProperty.GetString()
			// Get the value of the property
			currentValue := currentProperty.GetFirstChild()
			// Compare the current property name with the new property name
			if propertyName != nil && currentPropertyName == *propertyName {
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
			currentPropertyName := currentProperty.GetString()
			if propertyName != nil && currentPropertyName == *propertyName {
				return false
			}
		} else if currentProperty.IsComputedProp() {
			keyNode := currentProperty.GetFirstChild()
			if !keyNode.IsStringLit() || (propertyName != nil && keyNode.GetString() == *propertyName) {
				return false
			}
		}
	}

	newProperty := rhino.IR_StringKey(*propertyName).SrcrefIfMissing(propertyCandidate)
	// Preserve the quotedness of a property reference
	if lhs.IsGetElem() {
		newProperty.SetQuotedStringKey()
	}
	newValue := rhs.Detach()
	newProperty.AddChildToBack(newValue)

	if existingProperty != nil {
		p.DeleteNode(existingProperty)
	}
	// If the property does not already exist we can safely add it
	objectLiteral.AddChildToBack(newProperty)
	propertyCandidate.Detach()
	return true
}

// Java: PeepholeCollectPropertyAssignments.mightContainForwardReference(Node, String)
func PeepholeCollectPropertyAssignments_MightContainForwardReference(node *rhino.Node, varName string) bool {
	if node.IsName() {
		return varName == node.GetString()
	}
	for child := node.GetFirstChild(); child != nil; child = child.GetNext() {
		if PeepholeCollectPropertyAssignments_MightContainForwardReference(child, varName) {
			return true
		}
	}
	return false
}
