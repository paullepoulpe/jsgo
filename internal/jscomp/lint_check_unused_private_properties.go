// Port of com/google/javascript/jscomp/lint/CheckUnusedPrivateProperties.java @ closure-compiler v20260415

/*
 * Copyright 2015 The Closure Compiler Authors.
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

// CheckUnusedPrivateProperties looks for properties that are never read.
// These can be properties created using "this", or static properties of
// constructors or interfaces. Explicitly ignored is the possibility that these
// properties may be indirectly referenced using "for-in" or "Object.keys".
type CheckUnusedPrivateProperties struct {
	compiler                  AbstractCompiler
	used                      *jutil.LinkedHashSet[string]
	candidates                []*rhino.Node
	constructorsAndInterfaces *jutil.LinkedHashSet[string]
}

var CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY = DiagnosticType_Disabled("JSC_UNUSED_PRIVATE_PROPERTY", "Private property {0} is never read")

// NewCheckUnusedPrivateProperties is new CheckUnusedPrivateProperties(compiler).
//
// Java: CheckUnusedPrivateProperties(AbstractCompiler)
func NewCheckUnusedPrivateProperties(compiler AbstractCompiler) *CheckUnusedPrivateProperties {
	return &CheckUnusedPrivateProperties{
		compiler:                  compiler,
		used:                      jutil.NewLinkedHashSet[string](),
		constructorsAndInterfaces: jutil.NewLinkedHashSet[string](),
	}
}

// Java: CheckUnusedPrivateProperties.process(Node, Node)
func (c *CheckUnusedPrivateProperties) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckUnusedPrivateProperties.reportUnused(NodeTraversal)
func (c *CheckUnusedPrivateProperties) ReportUnused(t *NodeTraversal) {
	for _, n := range c.candidates {
		propName := c.GetPropName(n)
		if !c.used.Contains(propName) {
			t.Report(n, CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY, propName)
		}
	}
}

// Java: CheckUnusedPrivateProperties.getPropName(Node)
func (c *CheckUnusedPrivateProperties) GetPropName(n *rhino.Node) string {
	switch n.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_MEMBER_FUNCTION_DEF:
		return n.GetString()
	default:
	}
	panic(jutil.NewRuntimeException("Unexpected node type: " + n.String()))
}

// Java: CheckUnusedPrivateProperties.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckUnusedPrivateProperties) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if n.IsScript() {
		c.used.Clear()
		// Even if a constructor is unused , we still want to prevent the construction of the class
		// outside the file.  i.e. a container of static methods.
		c.used.Add("constructor")
		clear(c.candidates)
		c.candidates = c.candidates[:0]
	}
	return true
}

// Java: CheckUnusedPrivateProperties.visit(NodeTraversal, Node, Node)
func (c *CheckUnusedPrivateProperties) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		// exiting the script, report any privates not used in the file.
		c.ReportUnused(t)
	case rhino.Token_GETPROP:
		propName := n.GetString()
		if CheckUnusedPrivateProperties_IsPinningPropertyUse(n) || !c.IsCandidatePropertyDefinition(n) {
			c.used.Add(propName)
		} else {
			// Only consider "private" properties.
			if c.IsCheckablePrivatePropDecl(n) {
				c.candidates = append(c.candidates, n)
			}
		}
	case rhino.Token_MEMBER_FUNCTION_DEF:
		// Only consider "private" methods.
		if c.IsCheckablePrivatePropDecl(n) {
			c.candidates = append(c.candidates, n)
		}
	case rhino.Token_OBJECTLIT:
		// Assume any object literal definition might be a reflection on the
		// class property.
		for c2 := n.GetFirstChild(); c2 != nil; c2 = c2.GetNext() {
			if c2.IsStringKey() || c2.IsGetterDef() || c2.IsSetterDef() || c2.IsMemberFunctionDef() {
				c.used.Add(c2.GetString())
			}
		}
	case rhino.Token_CALL:
		// Look for properties referenced through a property rename function.
		target := n.GetFirstChild()
		if n.HasMoreThanOneChild() &&
			c.compiler.GetCodingConvention().IsPropertyRenameFunction(target) {
			propName := target.GetNext()
			if propName.IsStringLit() {
				c.used.Add(propName.GetString())
			}
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

// Java: CheckUnusedPrivateProperties.recordConstructorOrInterface(Node)
func (c *CheckUnusedPrivateProperties) RecordConstructorOrInterface(n *rhino.Node) {
	className := NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n))
	if className != nil {
		// If className is null, then this pass won't report any diagnostics on static members of this
		// class.
		c.constructorsAndInterfaces.Add(*className)
	}
}

// Java: CheckUnusedPrivateProperties.isPrivatePropDecl(Node)
func (c *CheckUnusedPrivateProperties) IsPrivatePropDecl(n *rhino.Node) bool {
	// TODO(johnlenz): add support private by convention property definitions without JSDoc.
	// TODO(johnlenz): add support for checking protected properties in final classes
	// TODO(johnlenz): add support for checking "package" properties when checking an entire
	// library.
	info := NodeUtil_GetBestJSDocInfo(n)
	return info != nil && info.GetVisibility() == rhino.JSDocInfo_Visibility_PRIVATE
}

// Java: CheckUnusedPrivateProperties.isCheckablePrivatePropDecl(Node)
func (c *CheckUnusedPrivateProperties) IsCheckablePrivatePropDecl(n *rhino.Node) bool {
	// TODO(tbreisacher): Look for uses of the typedef/interface in type expressions; warn if there
	// are no uses.
	info := NodeUtil_GetBestJSDocInfo(n)
	return c.IsPrivatePropDecl(n) && !info.HasTypedefType() && !info.IsInterface()
}

// Java: CheckUnusedPrivateProperties.isCandidatePropertyDefinition(Node)
func (c *CheckUnusedPrivateProperties) IsCandidatePropertyDefinition(n *rhino.Node) bool {
	jutil.CheckStateMsg1(n.IsGetProp(), "%s", n)
	target := n.GetFirstChild()
	return target.IsThis() ||
		c.containsQualifiedName(target) ||
		(target.IsGetProp() && target.GetString() == "prototype")
}

// containsQualifiedName is
// constructorsAndInterfaces.contains(target.getQualifiedName()).
//
// Go-only: language-adapter — getQualifiedName() may be null (*string here),
// and a Java set of non-null strings never contains null.
func (c *CheckUnusedPrivateProperties) containsQualifiedName(target *rhino.Node) bool {
	qname := target.GetQualifiedName()
	return qname != nil && c.constructorsAndInterfaces.Contains(*qname)
}

// CheckUnusedPrivateProperties_IsPinningPropertyUse returns whether the
// property is used in a way that prevents its removal.
//
// Java: CheckUnusedPrivateProperties.isPinningPropertyUse(Node)
func CheckUnusedPrivateProperties_IsPinningPropertyUse(n *rhino.Node) bool {
	// Rather than looking for cases that are uses, we assume all references are
	// pinning uses unless they are:
	//  - a simple assignment (x.a = 1)
	//  - a compound assignment or increment (x++, x += 1) whose result is
	//    otherwise unused

	parent := n.GetParent()
	if n == parent.GetFirstChild() {
		if parent.IsExprResult() {
			// A stub declaration "this.x;" isn't a pinning use.
			return false
		} else if parent.IsAssign() {
			// A simple assignment doesn't pin the property.
			return false
		} else if NodeUtil_IsAssignmentOp(parent) || parent.IsInc() || parent.IsDec() {
			// In general, compound assignments are both reads and writes, but
			// if the property is never otherwise read we can consider it simply
			// a write.
			// However if the assign expression is used as part of a larger
			// expression, we must consider it a read. For example:
			//    x = (y.a += 1);
			return NodeUtil_IsExpressionResultUsed(parent)
		}
	}
	return true
}
