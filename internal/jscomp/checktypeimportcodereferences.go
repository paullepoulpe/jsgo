// Port of com/google/javascript/jscomp/CheckTypeImportCodeReferences.java @ closure-compiler v20260415

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckTypeImportCodeReferences implements the JSC_TYPE_IMPORT_CODE_REFERENCE
// diagnostic: fires when a name declared via goog.requireType() is referenced
// in a non-type-annotation position (i.e. in executable code).
//
// Mirrors CheckTypeImportCodeReferences.java's AbstractPostOrderCallback walk.
// Java: extends AbstractPostOrderCallback implements CompilerPass.
type CheckTypeImportCodeReferences struct {
	path  string
	diags []diag.Diagnostic

	// requireTypeNames maps the declaration NAME node pointer to the
	// name string. Populated by ctirCollect in a first pass over the
	// AST; consulted during visit to distinguish the declaration
	// site (excluded) from code references (flagged).
	requireTypeNames map[*rhino.Node]string
}

// NewCheckTypeImportCodeReferences mirrors the Java constructor
// CheckTypeImportCodeReferences(AbstractCompiler compiler).
// Go-only: language-adapter — Go constructor shim (Java uses class constructor directly; Go has no class constructors).
func NewCheckTypeImportCodeReferences(path string) *CheckTypeImportCodeReferences {
	return &CheckTypeImportCodeReferences{
		path:             path,
		requireTypeNames: map[*rhino.Node]string{},
	}
}

// Check runs the two-phase walk and returns collected diagnostics.
// Mirrors Java's process(Node externs, Node root): calls NodeTraversal.traverse(compiler, root, this).
// Phase 1 (ctirCollect) builds the requireTypeNames index from
// var/let/const declarations of the form `name = goog.requireType(…)`.
// Phase 2 (post-order visit) flags code references.
func (c *CheckTypeImportCodeReferences) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.requireTypeNames = map[*rhino.Node]string{}
	c.ctirCollect(script)
	c.walk(script)
	return c.diags
}

// ctirCollect is a pre-order walk that populates requireTypeNames.
// It mirrors the information Java's TypedScopeCreator stores in the
// Var object — specifically Var.getNameNode() and the call chain
// check performed in visit(). By collecting first we avoid needing
// scope infrastructure.
// TODO(port): TypedScopeCreator / Var.getNameNode() — scope-free substitute (TypedScopeCreator not yet ported; collecting declarations manually).
func (c *CheckTypeImportCodeReferences) ctirCollect(n *rhino.Node) {
	if n == nil {
		return
	}
	// Look for: var/let/const <nameBinding> = goog.requireType(…)
	// Java's visit() checks NodeUtil.isNameDeclaration(decl) &&
	// decl.hasOneChild() && decl.getFirstChild().hasChildren() &&
	// decl.getFirstChild().getLastChild().isCall().
	if isNameDeclaration(n) {
		kids := n.Children()
		if len(kids) == 1 {
			binding := kids[0] // VAR's single child
			bindingKids := binding.Children()
			if len(bindingKids) > 0 {
				callNode := bindingKids[len(bindingKids)-1]
				if callNode.Token() == rhino.CALL && matchesGoogRequireType(callNode) {
					// Simple form: var Foo = goog.requireType('…')
					// binding is NAME with Foo as its string.
					if binding.Token() == rhino.NAME {
						c.requireTypeNames[binding] = binding.String_()
					}
					// Destructuring form: var {Foo} = goog.requireType('…')
					// binding is OBJECT_PATTERN / DESTRUCTURING_LHS;
					// each STRING_KEY child carries the bound name.
					// Java: declNameNode.getParent().isStringKey() branch.
					for _, skNode := range binding.Children() {
						if skNode.Token() == rhino.STRING_KEY {
							// The STRING_KEY's first child is the bound NAME.
							if nameChild := skNode.FirstChild(); nameChild != nil && nameChild.Token() == rhino.NAME {
								c.requireTypeNames[nameChild] = nameChild.String_()
							}
						}
					}
				}
			}
		}
	}
	for _, child := range n.Children() {
		c.ctirCollect(child)
	}
}

// walk is the post-order traversal driver. Mirrors AbstractPostOrderCallback
// behaviour: visits children before the current node, then calls visit().
// TODO(port): NodeTraversal.traverse(AbstractCompiler, Node, Callback) — NodeTraversal not yet ported.
func (c *CheckTypeImportCodeReferences) walk(n *rhino.Node) {
	if n == nil {
		return
	}
	// Post-order: children first.
	for _, child := range n.Children() {
		c.walk(child)
	}
	c.visit(n)
}

// visit checks a single node. Mirrors Java's
// visit(NodeTraversal t, Node n, Node parent):
// checks that NAME references to goog.requireType'd bindings are not
// used in code (non-type-annotation) positions.
func (c *CheckTypeImportCodeReferences) visit(n *rhino.Node) {
	if n.Token() != rhino.NAME {
		return
	}
	name := n.String_()
	if name == "" {
		return
	}
	// Check if n is itself one of the declaration nodes (Java: declNameNode == n).
	if _, isDeclNode := c.requireTypeNames[n]; isDeclNode {
		return
	}
	// Check if any requireType binding matches this name.
	// Java uses t.getScope().getVar(n.getString()).getNameNode() to get
	// the declaration node and then checks it against requireTypeNames.
	// Since we index by declaration node we can't do that lookup in O(1)
	// without a reverse name→declNode map. Build the check by name match.
	for declNode, declName := range c.requireTypeNames {
		if declName == name {
			// Additional guard: skip if n IS the declaration node
			// (already handled above, but double-check for the
			// destructuring path where declNode is a child of STRING_KEY).
			if declNode == n {
				return
			}
			c.diags = append(c.diags, diag.Diagnostic{
				Path:    c.path,
				Line:    n.Lineno(),
				Col:     n.Charno(),
				Code:    diag.JSC_TYPE_IMPORT_CODE_REFERENCE,
				Level:   diag.Error,
				Message: fmt.Sprintf("Cannot reference goog.requireType()'d name %s outside of a type annotation.", name),
			})
			return
		}
	}
}

// isNameDeclaration reports whether n is a var/let/const statement.
// Mirrors Java's NodeUtil.isNameDeclaration(Node) (Node.java family of token
// checks). The canonical Go port is rhino.Node.IsVar/IsLet/IsConst; this
// helper is local to CheckTypeImportCodeReferences for naming clarity.
func isNameDeclaration(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	return n.IsVar() || n.IsLet() || n.IsConst()
}

// matchesGoogRequireType reports whether callNode is a call to
// goog.requireType(…). Mirrors Java's
// callNode.getFirstChild().matchesQualifiedName(GOOG_REQUIRE_TYPE)
// where GOOG_REQUIRE_TYPE = IR.getprop(IR.name("goog"), "requireType").
// Delegates to rhino.Node.MatchesQualifiedNameStr (port of Java
// Node.matchesQualifiedName(String)).
func matchesGoogRequireType(callNode *rhino.Node) bool {
	if callNode == nil || callNode.Token() != rhino.CALL {
		return false
	}
	callee := callNode.FirstChild()
	if callee == nil {
		return false
	}
	return callee.MatchesQualifiedNameStr("goog.requireType")
}
