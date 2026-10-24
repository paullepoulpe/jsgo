// Port of com/google/javascript/jscomp/J2clConstantHoisterPass.java @ closure-compiler v20261021

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type j2clFieldAssignments struct {
	first  *rhino.Node
	second *rhino.Node
	count  int32
}

type j2clConstantHoisterCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *J2clConstantHoisterPass
}

func (c *j2clConstantHoisterCollector) Visit(t *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	p := c.outer
	if parent != nil && NodeUtil_IsLValue(node) {
		if q := node.GetQualifiedName(); q != nil {
			key := *q
			entry := p.fieldAssignments[key]
			switch entry.count {
			case 0:
				entry.first = parent
				entry.count = 1
				p.fieldOrder = append(p.fieldOrder, key)
			case 1:
				entry.second = parent
				entry.count = 2
			default:
				entry.count++
			}
			p.fieldAssignments[key] = entry
		}
	}
	if j2clConstantHoisterPass_isHoistableFunction(t, node) {
		p.hoistableFunctions.Add(node)
	}
}

// J2clConstantHoisterPass is an optimization pass for J2CL-generated code to
// hoist some constant assignments out clinit method to declaration phase so
// they could be used by other optimization passes for static evaluation.
//
// Java: public class J2clConstantHoisterPass implements CompilerPass
type J2clConstantHoisterPass struct {
	compiler           AbstractCompiler
	fieldAssignments   map[string]j2clFieldAssignments
	fieldOrder         []string
	hoistableFunctions *jutil.LinkedHashSet[*rhino.Node]
	collector          j2clConstantHoisterCollector
}

// Java: J2clConstantHoisterPass(AbstractCompiler)
func NewJ2clConstantHoisterPass(compiler AbstractCompiler) *J2clConstantHoisterPass {
	p := &J2clConstantHoisterPass{compiler: compiler}
	p.collector.outer = p
	return p
}

// Java: J2clConstantHoisterPass.process(Node, Node)
func (p *J2clConstantHoisterPass) Process(externs *rhino.Node, root *rhino.Node) {
	if !J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
		return
	}

	if p.fieldAssignments == nil {
		p.fieldAssignments = make(map[string]j2clFieldAssignments, 64)
		p.hoistableFunctions = jutil.NewLinkedHashSet[*rhino.Node]()
	} else {
		clear(p.fieldAssignments)
		p.fieldOrder = p.fieldOrder[:0]
		p.hoistableFunctions.Clear()
	}

	NodeTraversal_Traverse(p.compiler, root, &p.collector)

	var pair [2]*rhino.Node
	for _, key := range p.fieldOrder {
		entry := p.fieldAssignments[key]
		if entry.count == 2 {
			pair[0] = entry.first
			pair[1] = entry.second
			p.maybeHoistClassField(pair[:], p.hoistableFunctions)
		}
	}
}

// j2clConstantHoisterPass_isHoistableFunction returns whether the specified
// rValue is a function which does not receive any variables from its
// containing scope, and is thus 'hoistable'.
//
// Java: J2clConstantHoisterPass.isHoistableFunction(NodeTraversal, Node)
func j2clConstantHoisterPass_isHoistableFunction(t *NodeTraversal, node *rhino.Node) bool {
	// TODO(michaelthomas): This could be improved slightly by not assuming that any variable in the
	// outer scope is used in the function.
	return node.IsFunction() && t.GetScope().GetVarCount() == 0
}

// Java: J2clConstantHoisterPass.maybeHoistClassField(Collection<Node>, Collection<Node>)
func (p *J2clConstantHoisterPass) maybeHoistClassField(
	assignments []*rhino.Node, hoistableFunctions *jutil.LinkedHashSet[*rhino.Node]) {
	// The field is only assigned twice:
	if len(assignments) != 2 {
		return
	}

	first := assignments[0]
	second := assignments[1]

	// One of them is the top level declaration and the other is the assignment in clinit.
	topLevelDeclaration := second
	if j2clConstantHoisterPass_isClassFieldDeclaration(first) {
		topLevelDeclaration = first
	}
	clinitAssignment := second
	if j2clConstantHoisterPass_isClinitFieldAssignment(first) {
		clinitAssignment = first
	}
	if !j2clConstantHoisterPass_isClassFieldDeclaration(topLevelDeclaration) ||
		!j2clConstantHoisterPass_isClinitFieldAssignment(clinitAssignment) {
		return
	}

	// And it is assigned to a literal value; hence could be used in static eval and safe to move:
	assignmentRhs := clinitAssignment.GetSecondChild()
	if !NodeUtil_IsLiteralValue(assignmentRhs /* includeFunctions= */, true) ||
		(assignmentRhs.IsFunction() && !hoistableFunctions.Contains(assignmentRhs)) {
		return
	}

	// And the assignment are in the same script:
	if NodeUtil_GetEnclosingScript(clinitAssignment) != NodeUtil_GetEnclosingScript(topLevelDeclaration) {
		return
	}

	// At this point the only case some could observe the declaration value is the when you have
	// cycle between clinits and the field is accessed before initialization; which is almost always
	// a bug and GWT never assumed this state is observable in its optimization, yet nobody
	// complained. So it is safe to upgrade it to a constant.
	p.hoistConstantLikeField(clinitAssignment, topLevelDeclaration)
}

// Java: J2clConstantHoisterPass.hoistConstantLikeField(Node, Node)
func (p *J2clConstantHoisterPass) hoistConstantLikeField(clinitAssignment *rhino.Node, topLevelDeclaration *rhino.Node) {
	clinitAssignedValue := clinitAssignment.GetSecondChild()
	declarationInClass := topLevelDeclaration.GetFirstChild()
	declarationAssignedValue := declarationInClass.GetFirstChild()

	clinitChangeScope := ChangeTracker_GetEnclosingChangeScopeRoot(clinitAssignment)
	// Remove the clinit initialization
	NodeUtil_RemoveChild(clinitAssignment.GetParent(), clinitAssignment)

	clinitAssignedValue.Detach()
	p.compiler.ReportChangeToChangeScope(clinitChangeScope)
	if declarationAssignedValue == nil {
		// Add the value from clinit to the stub declaration
		declarationInClass.AddChildToFront(clinitAssignedValue)
		p.compiler.ReportChangeToEnclosingScope(topLevelDeclaration)
	} else if declarationAssignedValue != clinitAssignedValue {
		jutil.CheckState(NodeUtil_IsLiteralValue(declarationAssignedValue /* includeFunctions= */, false))
		// Replace the assignment in declaration with the value from clinit
		declarationAssignedValue.ReplaceWith(clinitAssignedValue)
		p.compiler.ReportChangeToEnclosingScope(topLevelDeclaration)
	}
	declarationInClass.SetDeclaredConstantVar(true)
}

// j2clConstantHoisterPass_isClassFieldDeclaration matches literal value
// declarations `var foo = 3;` or stub declarations like `var bar;`.
//
// Java: J2clConstantHoisterPass.isClassFieldDeclaration(Node)
func j2clConstantHoisterPass_isClassFieldDeclaration(node *rhino.Node) bool {
	return node.GetParent().IsScript() &&
		node.IsVar() &&
		(!node.GetFirstChild().HasChildren() ||
			(node.GetFirstFirstChild() != nil &&
				NodeUtil_IsLiteralValue(
					node.GetFirstFirstChild() /* includeFunctions= */, false)))
}

// Java: J2clConstantHoisterPass.isClinitFieldAssignment(Node)
func j2clConstantHoisterPass_isClinitFieldAssignment(node *rhino.Node) bool {
	return node.GetParent().IsExprResult() &&
		node.GetGrandparent().IsBlock() &&
		j2clConstantHoisterPass_isClinitMethod(node.GetGrandparent().GetParent())
}

// TODO(goktug): Create a utility to share this logic and start using getQualifiedOriginalName.
//
// Java: J2clConstantHoisterPass.isClinitMethod(Node)
func j2clConstantHoisterPass_isClinitMethod(fnNode *rhino.Node) bool {
	if !fnNode.IsFunction() {
		return false
	}

	fnName := NodeUtil_GetName(fnNode)
	return fnName != nil && j2clConstantHoisterPass_isClinitMethodName(fnName)
}

// Java: J2clConstantHoisterPass.isClinitMethodName(String)
func j2clConstantHoisterPass_isClinitMethodName(methodName *string) bool {
	return methodName != nil &&
		(strings.HasSuffix(*methodName, "$$0clinit") || strings.HasSuffix(*methodName, ".$clinit"))
}
