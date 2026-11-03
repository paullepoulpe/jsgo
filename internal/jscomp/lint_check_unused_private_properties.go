// Port of com/google/javascript/jscomp/lint/CheckUnusedPrivateProperties.java @ closure-compiler v20261028

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var checkUnusedPrivatePropertiesVisitTokens = makeVisitTokenSet(
	rhino.Token_SCRIPT,
	rhino.Token_GETPROP,
	rhino.Token_MEMBER_FUNCTION_DEF,
	rhino.Token_OBJECTLIT,
	rhino.Token_CALL,
	rhino.Token_FUNCTION,
	rhino.Token_CLASS,
)

var _ nodeTraversal_visitTokens = (*CheckUnusedPrivateProperties)(nil)

func (*CheckUnusedPrivateProperties) visitTokens() *ast.TokenSet {
	return &checkUnusedPrivatePropertiesVisitTokens
}

// CheckUnusedPrivateProperties looks for properties that are never read.
// These can be properties created using "this", or static properties of
// constructors or interfaces. Explicitly ignored is the possibility that these
// properties may be indirectly referenced using "for-in" or "Object.keys".
type CheckUnusedPrivateProperties struct {
	compiler                  AbstractCompiler
	used                      map[string]bool
	candidates                []rhino.Node
	constructorsAndInterfaces map[string]bool
}

var CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY = DiagnosticType_Disabled(
	"JSC_UNUSED_PRIVATE_PROPERTY",
	"Private property {0} is never read",
)

// NewCheckUnusedPrivateProperties constructs a CheckUnusedPrivateProperties pass.
//
// Java: CheckUnusedPrivateProperties(AbstractCompiler)
func NewCheckUnusedPrivateProperties(compiler AbstractCompiler) *CheckUnusedPrivateProperties {
	return &CheckUnusedPrivateProperties{
		compiler:                  compiler,
		used:                      make(map[string]bool),
		constructorsAndInterfaces: make(map[string]bool),
	}
}

// Java: CheckUnusedPrivateProperties.process(Node, Node)
func (c *CheckUnusedPrivateProperties) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

func (c *CheckUnusedPrivateProperties) reportUnused(t *NodeTraversal) {
	for _, n := range c.candidates {
		name := candidatePropName(n)
		if !c.used[name] {
			t.Report(n, CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY, name)
		}
	}
}

func candidatePropName(n rhino.Node) string {
	if gp, ok := ast.AsGetProp(n); ok {
		return gp.Property()
	}
	if m, ok := ast.AsMemberFunctionDef(n); ok {
		return m.Key()
	}
	panic(jutil.NewRuntimeException("Unexpected node type: " + n.String()))
}

// Java: CheckUnusedPrivateProperties.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckUnusedPrivateProperties) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		clear(c.used)
		// Even if a constructor is unused, we still want to prevent the construction of the class
		// outside the file (e.g. a container of static methods).
		c.used["constructor"] = true
		clear(c.candidates)
		c.candidates = c.candidates[:0]
	}
	return true
}

// Java: CheckUnusedPrivateProperties.visit(NodeTraversal, Node, Node)
func (c *CheckUnusedPrivateProperties) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		// Exiting the script: report any privates not used in the file.
		c.reportUnused(t)
	case rhino.Token_GETPROP:
		gp := ast.MustGetProp(n)
		if isPinningPropertyUse(n) || !c.isCandidatePropertyDefinition(gp) {
			c.used[gp.Property()] = true
		} else if c.isCheckablePrivatePropDecl(n) {
			c.candidates = append(c.candidates, n)
		}
	case rhino.Token_MEMBER_FUNCTION_DEF:
		// Only consider "private" methods.
		if c.isCheckablePrivatePropDecl(n) {
			c.candidates = append(c.candidates, n)
		}
	case rhino.Token_OBJECTLIT:
		// Assume any object literal definition might be a reflection on the class property.
		for member := range ast.MustObjectLit(n).Members() {
			if pkm, ok := member.AsPlainKeyMember(); ok {
				c.used[pkm.Key()] = true
			}
		}
	case rhino.Token_CALL:
		// Look for properties referenced through a property rename function.
		call := ast.MustCall(n)
		if firstArg, ok := call.FirstArg(); ok &&
			c.compiler.GetCodingConvention().IsPropertyRenameFunction(call.Callee().Node()) {
			if propName, ok := ast.AsStringLit(firstArg.Node()); ok {
				c.used[propName.Value()] = true
			}
		}
	case rhino.Token_FUNCTION:
		if info := NodeUtil_GetBestJSDocInfo(n); info != nil && (info.IsConstructor() || info.IsInterface()) {
			c.recordConstructorOrInterface(n)
		}
	case rhino.Token_CLASS:
		c.recordConstructorOrInterface(n)
	}
}

func (c *CheckUnusedPrivateProperties) recordConstructorOrInterface(n rhino.Node) {
	if className := NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n)); className != nil {
		// If className is null, this pass won't report diagnostics on static members of this class.
		c.constructorsAndInterfaces[*className] = true
	}
}

func (c *CheckUnusedPrivateProperties) isPrivatePropDecl(n rhino.Node) bool {
	// TODO(johnlenz): add support private by convention property definitions without JSDoc.
	// TODO(johnlenz): add support for checking protected properties in final classes
	// TODO(johnlenz): add support for checking "package" properties when checking an entire
	// library.
	info := NodeUtil_GetBestJSDocInfo(n)
	return info != nil && info.GetVisibility() == rhino.JSDocInfo_Visibility_PRIVATE
}

func (c *CheckUnusedPrivateProperties) isCheckablePrivatePropDecl(n rhino.Node) bool {
	// TODO(tbreisacher): Look for uses of the typedef/interface in type expressions; warn if there
	// are no uses.
	info := NodeUtil_GetBestJSDocInfo(n)
	return c.isPrivatePropDecl(n) && !info.HasTypedefType() && !info.IsInterface()
}

func (c *CheckUnusedPrivateProperties) isCandidatePropertyDefinition(gp ast.GetProp) bool {
	target := gp.Object().Node()
	if target.IsThis() {
		return true
	}
	if qname, ok := target.GetQualifiedNameString(); ok && c.constructorsAndInterfaces[qname] {
		return true
	}
	if targetProp, ok := gp.Object().AsGetProp(); ok && targetProp.Property() == "prototype" {
		return true
	}
	return false
}

// isPinningPropertyUse reports whether n is used in a way that prevents its removal.
func isPinningPropertyUse(n rhino.Node) bool {
	// Rather than looking for cases that are uses, we assume all references are
	// pinning uses unless they are:
	//  - a stub declaration (`this.x;`)
	//  - a simple assignment (`x.a = 1`)
	//  - a compound assignment or increment (`x++`, `x += 1`) whose result is otherwise unused
	parent := n.GetParent()
	switch ast.RoleIn(n, parent) {
	case ast.RoleExprStmtExpr, ast.RoleAssignTarget:
		return false
	case ast.RoleAssignOpTarget, ast.RoleIncDecTarget:
		return NodeUtil_IsExpressionResultUsed(parent)
	}
	return true
}
