// Port of com/google/javascript/jscomp/J2clConstantHoisterPass.java @ closure-compiler v20261104

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type j2clFieldAssignments struct {
	first  rhino.Node
	second rhino.Node
	count  int32
}

type j2clConstantHoisterCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *J2clConstantHoisterPass
}

func (c *j2clConstantHoisterCollector) Visit(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	p := c.outer
	if !parent.IsNil() && NodeUtil_IsLValue(node) {
		if key, ok := node.GetQualifiedNameString(); ok {
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
	if isHoistableJ2clFunction(t, node) {
		p.hoistableFunctions[node] = true
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
	hoistableFunctions map[rhino.Node]bool
	collector          j2clConstantHoisterCollector
}

// Java: J2clConstantHoisterPass(AbstractCompiler)
func NewJ2clConstantHoisterPass(compiler AbstractCompiler) *J2clConstantHoisterPass {
	p := &J2clConstantHoisterPass{compiler: compiler}
	p.collector.outer = p
	return p
}

// Java: J2clConstantHoisterPass.process(Node, Node)
func (p *J2clConstantHoisterPass) Process(externs rhino.Node, root rhino.Node) {
	if !J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
		return
	}

	if p.fieldAssignments == nil {
		p.fieldAssignments = make(map[string]j2clFieldAssignments, 64)
		p.hoistableFunctions = make(map[rhino.Node]bool)
	} else {
		clear(p.fieldAssignments)
		p.fieldOrder = p.fieldOrder[:0]
		clear(p.hoistableFunctions)
	}

	NodeTraversal_Traverse(p.compiler, root, &p.collector)

	for _, key := range p.fieldOrder {
		entry := p.fieldAssignments[key]
		if entry.count == 2 {
			p.maybeHoistClassField(entry.first, entry.second, p.hoistableFunctions)
		}
	}
}

// isHoistableJ2clFunction returns whether the specified rValue is a function
// which does not receive any variables from its containing scope, and is thus
// 'hoistable'.
//
// Java: J2clConstantHoisterPass.isHoistableFunction(NodeTraversal, Node)
func isHoistableJ2clFunction(t *NodeTraversal, node rhino.Node) bool {
	// TODO(michaelthomas): This could be improved slightly by not assuming that any variable in the
	// outer scope is used in the function.
	return node.IsFunction() && t.GetScope().GetVarCount() == 0
}

// Java: J2clConstantHoisterPass.maybeHoistClassField(Collection<Node>, Collection<Node>)
func (p *J2clConstantHoisterPass) maybeHoistClassField(
	first, second rhino.Node, hoistableFunctions map[rhino.Node]bool) {
	// One of the two assignments is the top-level declaration and the other is the assignment in clinit.
	topLevelDeclaration := second
	declName, hasDecl := asJ2clClassFieldDeclaration(first)
	if hasDecl {
		topLevelDeclaration = first
	} else {
		declName, hasDecl = asJ2clClassFieldDeclaration(second)
	}

	clinitAssign, hasClinit := asJ2clClinitFieldAssignment(first)
	if !hasClinit {
		clinitAssign, hasClinit = asJ2clClinitFieldAssignment(second)
	}
	if !hasDecl || !hasClinit {
		return
	}

	// And it is assigned to a literal value; hence could be used in static eval and safe to move:
	assignmentRhs := clinitAssign.Value().Node()
	if !NodeUtil_IsLiteralValue(assignmentRhs /* includeFunctions= */, true) ||
		(assignmentRhs.IsFunction() && !hoistableFunctions[assignmentRhs]) {
		return
	}

	// And the assignments are in the same script:
	if NodeUtil_GetEnclosingScript(clinitAssign.Node()) != NodeUtil_GetEnclosingScript(topLevelDeclaration) {
		return
	}

	// At this point the only case some could observe the declaration value is the when you have
	// cycle between clinits and the field is accessed before initialization; which is almost always
	// a bug and GWT never assumed this state is observable in its optimization, yet nobody
	// complained. So it is safe to upgrade it to a constant.
	p.hoistConstantLikeField(clinitAssign, declName, topLevelDeclaration)
}

// Java: J2clConstantHoisterPass.hoistConstantLikeField(Node, Node)
func (p *J2clConstantHoisterPass) hoistConstantLikeField(
	clinitAssign ast.Assign, declName ast.Name, topLevelDeclaration rhino.Node) {
	clinitAssignment := clinitAssign.Node()
	clinitAssignedValue := clinitAssign.Value().Node()
	declarationInClass := declName.Node()

	clinitChangeScope := ChangeTracker_GetEnclosingChangeScopeRoot(clinitAssignment)
	// Remove the clinit initialization
	NodeUtil_RemoveChild(clinitAssignment.GetParent(), clinitAssignment)

	clinitAssignedValue.Detach()
	p.compiler.ReportChangeToChangeScope(clinitChangeScope)
	if init, hasInit := declName.Init(); !hasInit {
		// Add the value from clinit to the stub declaration
		declarationInClass.AddChildToFront(clinitAssignedValue)
		p.compiler.ReportChangeToEnclosingScope(topLevelDeclaration)
	} else if declarationAssignedValue := init.Node(); declarationAssignedValue != clinitAssignedValue {
		jutil.CheckState(NodeUtil_IsLiteralValue(declarationAssignedValue /* includeFunctions= */, false))
		// Replace the assignment in declaration with the value from clinit
		declarationAssignedValue.ReplaceWith(clinitAssignedValue)
		p.compiler.ReportChangeToEnclosingScope(topLevelDeclaration)
	}
	declarationInClass.SetDeclaredConstantVar(true)
}

// asJ2clClassFieldDeclaration matches literal value declarations `var foo = 3;`
// or stub declarations like `var bar;`.
//
// Java: J2clConstantHoisterPass.isClassFieldDeclaration(Node)
func asJ2clClassFieldDeclaration(node rhino.Node) (ast.Name, bool) {
	if !node.IsVar() || !node.GetParent().IsScript() {
		return ast.Name{}, false
	}
	first, ok := ast.MustNameDeclaration(node).FirstDeclarator()
	if !ok {
		return ast.Name{}, false
	}
	name, ok := first.AsName()
	if !ok {
		return ast.Name{}, false
	}
	if init, hasInit := name.Init(); hasInit && !NodeUtil_IsLiteralValue(init.Node() /* includeFunctions= */, false) {
		return ast.Name{}, false
	}
	return name, true
}

// Java: J2clConstantHoisterPass.isClinitFieldAssignment(Node)
func asJ2clClinitFieldAssignment(node rhino.Node) (ast.Assign, bool) {
	assign, ok := ast.AsAssign(node)
	if !ok {
		return ast.Assign{}, false
	}
	if !node.GetParent().IsExprResult() ||
		!node.GetGrandparent().IsBlock() ||
		!isJ2clClinitMethod(node.GetGrandparent().GetParent()) {
		return ast.Assign{}, false
	}
	return assign, true
}

// TODO(goktug): Create a utility to share this logic and start using getQualifiedOriginalName.
//
// Java: J2clConstantHoisterPass.isClinitMethod(Node)
func isJ2clClinitMethod(fnNode rhino.Node) bool {
	if !fnNode.IsFunction() {
		return false
	}
	fnName, ok := NodeUtil_GetNameString(fnNode)
	return ok && isJ2clClinitMethodName(fnName)
}

// Java: J2clConstantHoisterPass.isClinitMethodName(String)
func isJ2clClinitMethodName(methodName string) bool {
	return strings.HasSuffix(methodName, "$$0clinit") || strings.HasSuffix(methodName, ".$clinit")
}
