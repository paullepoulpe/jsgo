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
// Port of com/google/javascript/jscomp/PeepholeCollectPropertyAssignments.java (generated mechanically by j2go).

package jscomp

import (
	"math"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.PeepholeCollectPropertyAssignments
type PeepholeCollectPropertyAssignments struct {
	AbstractPeepholeOptimizationBase
}

func (this_ *PeepholeCollectPropertyAssignments) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *PeepholeCollectPropertyAssignments) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeCollectPropertyAssignments) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeCollectPropertyAssignments", "PeepholeCollectPropertyAssignments"
}
func (this_ *PeepholeCollectPropertyAssignments) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*PeepholeCollectPropertyAssignments) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeCollectPropertyAssignments", "com.google.javascript.jscomp.PeepholeCollectPropertyAssignments", func(v_ any) bool {
		return j2gort.InstanceOf[*PeepholeCollectPropertyAssignments](v_)
	})
}
func (*PeepholeCollectPropertyAssignments) Implements_jscomp_AbstractPeepholeOptimization() {}

func (this_ *PeepholeCollectPropertyAssignments) initFields_PeepholeCollectPropertyAssignments() {
	_ = this_
}

// Java constructor: PeepholeCollectPropertyAssignments()
func NewPeepholeCollectPropertyAssignments() *PeepholeCollectPropertyAssignments {
	this_ := &PeepholeCollectPropertyAssignments{}
	this_.BindSelf_(this_)
	InitPeepholeCollectPropertyAssignments(this_)
	return this_
}

func InitPeepholeCollectPropertyAssignments(this_ *PeepholeCollectPropertyAssignments) {
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_PeepholeCollectPropertyAssignments()
}

// Java method: PeepholeCollectPropertyAssignments.optimizeSubtree(com.google.javascript.rhino.Node subtree)
func (this_ *PeepholeCollectPropertyAssignments) OptimizeSubtree(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	if !(subtree.IsScript()) && !(subtree.IsBlock()) {
		return subtree
	}
	var codeChanged bool = false
	_ = codeChanged
	{
		var child pkg_rhino.Node = subtree.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if !(NodeUtil_IsNameDeclaration(child)) && !(NodeUtil_IsExprAssign(child)) {
				continue
			}
			if !(PeepholeCollectPropertyAssignments_IsPropertyAssignmentToName(child.GetNext())) {
				continue
			}
			if !(child.HasOneChild()) {
				j2gort.CheckState(false)
			}
			var name pkg_rhino.Node = PeepholeCollectPropertyAssignments_GetName(child)
			_ = name
			if !(name.IsName()) {
				continue
			}
			var value pkg_rhino.Node = PeepholeCollectPropertyAssignments_GetValue(child)
			_ = value
			if (value == 0) || !(PeepholeCollectPropertyAssignments_IsInterestingValue(value)) {
				continue
			}
			var propertyCandidate pkg_rhino.Node = pkg_rhino.Node(0)
			_ = propertyCandidate
			for (func() pkg_rhino.Node {
				var asgn_1 pkg_rhino.Node = child.GetNext()
				propertyCandidate = asgn_1
				return asgn_1
			}()) != 0 {
				if !(this_.priv_PeepholeCollectPropertyAssignments_collectProperty_4153(propertyCandidate, name.GetString(), value)) {
					break
				}
				codeChanged = true
			}
		}
	}
	if codeChanged {
		this_.ReportChangeToEnclosingScope(subtree)
	}
	return subtree
}

// Java static method: PeepholeCollectPropertyAssignments.getName(com.google.javascript.rhino.Node n)
func PeepholeCollectPropertyAssignments_GetName(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if NodeUtil_IsNameDeclaration(n) {
		return n.GetFirstChild()
	} else if NodeUtil_IsExprAssign(n) {
		return n.GetFirstFirstChild()
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
}

// Java static method: PeepholeCollectPropertyAssignments.getValue(com.google.javascript.rhino.Node n)
func PeepholeCollectPropertyAssignments_GetValue(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if NodeUtil_IsNameDeclaration(n) {
		return n.GetFirstFirstChild()
	} else if NodeUtil_IsExprAssign(n) {
		return n.GetFirstChild().GetLastChild()
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
}

// Java static method: PeepholeCollectPropertyAssignments.isInterestingValue(com.google.javascript.rhino.Node n)
func PeepholeCollectPropertyAssignments_IsInterestingValue(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsObjectLit() || n.IsArrayLit())
}

// Java static method: PeepholeCollectPropertyAssignments.isPropertyAssignmentToName(com.google.javascript.rhino.Node propertyCandidate)
func PeepholeCollectPropertyAssignments_IsPropertyAssignmentToName(propertyCandidate pkg_rhino.Node) bool {
	_ = propertyCandidate
	if propertyCandidate == 0 {
		return false
	}
	if !(NodeUtil_IsExprAssign(propertyCandidate)) {
		return false
	}
	var expr pkg_rhino.Node = propertyCandidate.GetFirstChild()
	_ = expr
	var lhs pkg_rhino.Node = expr.GetFirstChild()
	_ = lhs
	if !(NodeUtil_IsNormalGet(lhs)) {
		return false
	}
	var obj pkg_rhino.Node = lhs.GetFirstChild()
	_ = obj
	return obj.IsName()
}

// Java method: PeepholeCollectPropertyAssignments.collectProperty(com.google.javascript.rhino.Node propertyCandidate, java.lang.String name, com.google.javascript.rhino.Node value)
func (this_ *PeepholeCollectPropertyAssignments) priv_PeepholeCollectPropertyAssignments_collectProperty_4153(propertyCandidate pkg_rhino.Node, name string, value pkg_rhino.Node) bool {
	_ = propertyCandidate
	_ = name
	_ = value
	if !(PeepholeCollectPropertyAssignments_IsPropertyAssignmentToName(propertyCandidate)) {
		return false
	}
	var lhs pkg_rhino.Node = propertyCandidate.GetFirstFirstChild()
	_ = lhs
	{
		var iifeCond_1 bool
		{

			sr_2 := name
			se_3 := lhs.GetFirstChild().GetString()
			iifeCond_1 = j2gort.CheckNotNullString(sr_2) == se_3
		}
		if !(iifeCond_1) {
			return false
		}
	}
	var rhs pkg_rhino.Node = lhs.GetNext()
	_ = rhs
	if this_.MayHaveSideEffects(rhs) || NodeUtil_CanBeSideEffected_1(rhs) {
		return false
	}
	if !(NodeUtil_IsLiteralValue(rhs, true)) && PeepholeCollectPropertyAssignments_MightContainForwardReference(rhs, name) {
		return false
	}
	switch value.TokenSwitchOrd_() {
	case int32(46) /* Token.ARRAYLIT */ :
		if !(PeepholeCollectPropertyAssignments_CollectArrayProperty(value, propertyCandidate)) {
			return false
		}
	case int32(47) /* Token.OBJECTLIT */ :
		if !(this_.priv_PeepholeCollectPropertyAssignments_collectObjectProperty_4154(value, propertyCandidate)) {
			return false
		}
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
	}
	return true
}

// Java static method: PeepholeCollectPropertyAssignments.collectArrayProperty(com.google.javascript.rhino.Node arrayLiteral, com.google.javascript.rhino.Node propertyCandidate)
func PeepholeCollectPropertyAssignments_CollectArrayProperty(arrayLiteral pkg_rhino.Node, propertyCandidate pkg_rhino.Node) bool {
	_ = arrayLiteral
	_ = propertyCandidate
	var assignment pkg_rhino.Node = propertyCandidate.GetFirstChild()
	_ = assignment
	var sizeOfArrayAtStart int32 = arrayLiteral.GetChildCount()
	_ = sizeOfArrayAtStart
	var maxIndexAssigned int32 = (sizeOfArrayAtStart - int32(1))
	_ = maxIndexAssigned
	var lhs pkg_rhino.Node = assignment.GetFirstChild()
	_ = lhs
	var rhs pkg_rhino.Node = lhs.GetNext()
	_ = rhs
	if !(lhs.IsGetElem()) {
		return false
	}
	var obj pkg_rhino.Node = lhs.GetFirstChild()
	_ = obj
	var property pkg_rhino.Node = obj.GetNext()
	_ = property
	if !(property.IsNumber()) {
		return false
	}
	var dindex float64 = property.GetDouble()
	_ = dindex
	if (!(dindex >= float64(int32(0))) || math.IsInf(dindex, 0)) || (dindex > float64(int64(2147483647))) {
		return false
	}
	var index int32 = javalang.D2I(dindex)
	_ = index
	if dindex != float64(index) {
		return false
	}
	if (maxIndexAssigned + int32(4)) < index {
		return false
	}
	if PeepholeCollectPropertyAssignments_ArrayLiteralHasSpread(arrayLiteral) {
		return false
	}
	if index > maxIndexAssigned {
		for maxIndexAssigned < (index - int32(1)) {
			var emptyNode pkg_rhino.Node = pkg_rhino.IR_Empty().Srcref(arrayLiteral)
			_ = emptyNode
			arrayLiteral.AddChildToBack(emptyNode)
			maxIndexAssigned = maxIndexAssigned + 1
		}
		arrayLiteral.AddChildToBack(rhs.Detach())
	} else {
		var currentValue pkg_rhino.Node = arrayLiteral.GetChildAtIndex(index)
		_ = currentValue
		if !(currentValue.IsEmpty()) {
			return false
		}
		currentValue.ReplaceWith(rhs.Detach())
	}
	propertyCandidate.Detach()
	return true
}

// Java static method: PeepholeCollectPropertyAssignments.arrayLiteralHasSpread(com.google.javascript.rhino.Node arrayLiteral)
func PeepholeCollectPropertyAssignments_ArrayLiteralHasSpread(arrayLiteral pkg_rhino.Node) bool {
	_ = arrayLiteral
	var script pkg_rhino.Node = NodeUtil_GetEnclosingScript(arrayLiteral)
	_ = script
	var scriptFeatures *pkg_parser.FeatureSet
	_ = scriptFeatures
	if script != 0 {
		scriptFeatures = NodeUtil_GetFeatureSetOfScript(script)
	} else {
		scriptFeatures = nil
	}
	if (scriptFeatures != nil) && !(scriptFeatures.Has(pkg_parser.StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS())) {
		return false
	}
	{
		var child pkg_rhino.Node = arrayLiteral.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsSpread() {
				return true
			}
		}
	}
	return false
}

// Java method: PeepholeCollectPropertyAssignments.collectObjectProperty(com.google.javascript.rhino.Node objectLiteral, com.google.javascript.rhino.Node propertyCandidate)
func (this_ *PeepholeCollectPropertyAssignments) priv_PeepholeCollectPropertyAssignments_collectObjectProperty_4154(objectLiteral pkg_rhino.Node, propertyCandidate pkg_rhino.Node) bool {
	_ = objectLiteral
	_ = propertyCandidate
	var assignment pkg_rhino.Node = propertyCandidate.GetFirstChild()
	_ = assignment
	var lhs pkg_rhino.Node = assignment.GetFirstChild()
	_ = lhs
	var rhs pkg_rhino.Node = lhs.GetNext()
	_ = rhs
	var obj pkg_rhino.Node = lhs.GetFirstChild()
	_ = obj
	if !(lhs.IsGetProp() || lhs.IsGetElem()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(lhs))
	}
	var propertyName string = j2gort.NullString
	_ = propertyName
	if lhs.IsGetElem() {
		var property pkg_rhino.Node = obj.GetNext()
		_ = property
		if property.IsNumber() {
			propertyName = this_.GetSideEffectFreeStringValue(property)
		} else if property.IsStringLit() {
			propertyName = property.GetString()
		} else {
			return false
		}
	} else if lhs.IsGetProp() {
		propertyName = lhs.GetString()
	} else {
		return false
	}
	var existingProperty pkg_rhino.Node = pkg_rhino.Node(0)
	_ = existingProperty
	{
		var currentProperty pkg_rhino.Node = objectLiteral.GetFirstChild()
		_ = currentProperty
		for ; currentProperty != 0; currentProperty = currentProperty.GetNext() {
			if currentProperty.IsStringKey() || currentProperty.IsMemberFunctionDef() {
				var currentPropertyName string = currentProperty.GetString()
				_ = currentPropertyName
				var currentValue pkg_rhino.Node = currentProperty.GetFirstChild()
				_ = currentValue
				if j2gort.CheckNotNullString(currentPropertyName) == propertyName {
					existingProperty = currentProperty
					var isCurrentValueSideEffect bool = NodeUtil_CanBeSideEffected_1(currentValue)
					_ = isCurrentValueSideEffect
					var isNewValueSideEffect bool = NodeUtil_CanBeSideEffected_1(rhs)
					_ = isNewValueSideEffect
					if isCurrentValueSideEffect || isNewValueSideEffect {
						return false
					}
					break
				}
			} else if currentProperty.IsGetterDef() || currentProperty.IsSetterDef() {
				var currentPropertyName_2 string = currentProperty.GetString()
				_ = currentPropertyName_2
				if j2gort.CheckNotNullString(currentPropertyName_2) == propertyName {
					return false
				}
			} else if currentProperty.IsComputedProp() {
				var keyNode pkg_rhino.Node = currentProperty.GetFirstChild()
				_ = keyNode
				if !(keyNode.IsStringLit()) || (j2gort.CheckNotNullString(keyNode.GetString()) == propertyName) {
					return false
				}
			}
		}
	}
	var newProperty pkg_rhino.Node = pkg_rhino.IR_StringKey_1(propertyName).SrcrefIfMissing(propertyCandidate)
	_ = newProperty
	if lhs.IsGetElem() {
		newProperty.SetQuotedStringKey()
	}
	var newValue pkg_rhino.Node = rhs.Detach()
	_ = newValue
	newProperty.AddChildToBack(newValue)
	if existingProperty != 0 {
		this_.DeleteNode(existingProperty)
	}
	objectLiteral.AddChildToBack(newProperty)
	propertyCandidate.Detach()
	return true
}

// Java static method: PeepholeCollectPropertyAssignments.mightContainForwardReference(com.google.javascript.rhino.Node node, java.lang.String varName)
func PeepholeCollectPropertyAssignments_MightContainForwardReference(node pkg_rhino.Node, varName string) bool {
	_ = node
	_ = varName
	if node.IsName() {
		return func() bool { sr_4 := varName; se_5 := node.GetString(); return j2gort.CheckNotNullString(sr_4) == se_5 }()
	}
	{
		var child pkg_rhino.Node = node.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if PeepholeCollectPropertyAssignments_MightContainForwardReference(child, varName) {
				return true
			}
		}
	}
	return false
}
