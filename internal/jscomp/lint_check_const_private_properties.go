// Port of com/google/javascript/jscomp/lint/CheckConstPrivateProperties.java @ closure-compiler v20260415

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckConstPrivateProperties looks for properties that are not modified and
// ensures they use the @const annotation.
type CheckConstPrivateProperties struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler                  AbstractCompiler
	candidates                []*rhino.Node
	modified                  *jutil.LinkedHashSet[string]
	constructorsAndInterfaces *jutil.LinkedHashSet[string]
}

var CheckConstPrivateProperties_MISSING_CONST_PROPERTY = DiagnosticType_Disabled(
	"JSC_MISSING_CONST_PROPERTY",
	"Private property {0} is never modified, use the @const annotation")

// NewCheckConstPrivateProperties is new CheckConstPrivateProperties(compiler).
//
// Java: CheckConstPrivateProperties(AbstractCompiler)
func NewCheckConstPrivateProperties(compiler AbstractCompiler) *CheckConstPrivateProperties {
	return &CheckConstPrivateProperties{
		compiler:                  compiler,
		modified:                  jutil.NewLinkedHashSet[string](),
		constructorsAndInterfaces: jutil.NewLinkedHashSet[string](),
	}
}

// Java: CheckConstPrivateProperties.process(Node, Node)
func (c *CheckConstPrivateProperties) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// ReportMissingConst reports the property definitions that should use the
// @const annotation.
//
// Java: CheckConstPrivateProperties.reportMissingConst(NodeTraversal)
func (c *CheckConstPrivateProperties) ReportMissingConst(t *NodeTraversal) {
	for _, n := range c.candidates {
		propName := n.GetString()
		if !c.modified.Contains(propName) {
			t.Report(n, CheckConstPrivateProperties_MISSING_CONST_PROPERTY, propName)
		}
	}
	clear(c.candidates)
	c.candidates = c.candidates[:0]
	c.modified.Clear()
}

// Java: CheckConstPrivateProperties.visit(NodeTraversal, Node, Node)
func (c *CheckConstPrivateProperties) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		// Exiting the script, report any non-const privates not modified in the file.
		c.ReportMissingConst(t)
	case rhino.Token_GETELEM, rhino.Token_GETPROP:
		var propNode *rhino.Node
		if n.IsGetProp() {
			propNode = n
		} else if n.GetLastChild().IsStringLit() {
			propNode = n.GetLastChild()
		} else {
			return
		}

		// Only consider non-const @private class properties as candidates
		if c.IsCandidatePropertyDefinition(n) {
			c.candidates = append(c.candidates, propNode)
		} else if c.IsModificationOp(n) {
			// Mark any other modification operation as a modified property, to deal with lambdas, etc
			c.modified.Add(propNode.GetString())
		}
	case rhino.Token_FUNCTION:
		info := NodeUtil_GetBestJSDocInfo(n)
		if info != nil && (info.IsConstructor() || info.IsInterface()) {
			c.RecordConstructorOrInterface(n)
		}
	case rhino.Token_CLASS:
		c.RecordConstructorOrInterface(n)
	default:
	}
}

// Java: CheckConstPrivateProperties.recordConstructorOrInterface(Node)
func (c *CheckConstPrivateProperties) RecordConstructorOrInterface(n *rhino.Node) {
	className := NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n))
	if className != nil {
		// If className is null, then this pass won't report any diagnostics on static members of this
		// class.
		c.constructorsAndInterfaces.Add(*className)
	}
}

// IsCandidatePropertyDefinition returns whether the given node is a @private
// property declaration that is not marked constant.
//
// Java: CheckConstPrivateProperties.isCandidatePropertyDefinition(Node)
func (c *CheckConstPrivateProperties) IsCandidatePropertyDefinition(n *rhino.Node) bool {
	if !NodeUtil_IsLhsOfAssign(n) {
		return false
	}

	target := n.GetFirstChild()
	// Check whether the given property access is on 'this' or a static property on a class.
	if !(target.IsThis() || c.containsQualifiedName(target)) {
		return false
	}

	info := NodeUtil_GetBestJSDocInfo(n)
	return info != nil &&
		info.GetVisibility() == rhino.JSDocInfo_Visibility_PRIVATE &&
		!info.IsConstant() &&
		!info.HasTypedefType() &&
		!info.HasEnumParameterType() &&
		!info.IsInterface() &&
		!c.IsFunctionProperty(n)
}

// containsQualifiedName is
// constructorsAndInterfaces.contains(target.getQualifiedName()).
//
// Go-only: language-adapter — getQualifiedName() may be null (*string here),
// and a Java set of non-null strings never contains null.
func (c *CheckConstPrivateProperties) containsQualifiedName(target *rhino.Node) bool {
	qname := target.GetQualifiedName()
	return qname != nil && c.constructorsAndInterfaces.Contains(*qname)
}

// IsFunctionProperty returns whether the given property declaration is
// assigned to a function.
//
// Java: CheckConstPrivateProperties.isFunctionProperty(Node)
func (c *CheckConstPrivateProperties) IsFunctionProperty(n *rhino.Node) bool {
	// TODO(dylandavidson): getAssignedValue does not support GETELEM.
	if n.IsGetElem() {
		return false
	}
	assignedValue := NodeUtil_GetAssignedValue(n)
	return assignedValue != nil && assignedValue.IsFunction()
}

// IsModificationOp returns whether the given property is modified in any way
// (assignment, increment/decrement, or 'delete' property).
//
// Java: CheckConstPrivateProperties.isModificationOp(Node)
func (c *CheckConstPrivateProperties) IsModificationOp(n *rhino.Node) bool {
	return NodeUtil_IsLValue(n) || n.GetParent().IsDelProp()
}
