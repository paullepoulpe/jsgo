// Port of com/google/javascript/jscomp/CheckGlobalThis.java @ closure-compiler v20260415
//
// Checks for certain uses of the `this` keyword that are considered unsafe
// because they are likely to reference the global `this` object unintentionally.
//
// A use of `this` is considered unsafe if it's on the left side of an
// assignment or a property access, and not inside one of the following:
//   - a class static initialization block
//   - a prototype method
//   - a function annotated with @constructor
//   - a function annotated with @this
//   - a function where there's no logical place to put a @this annotation

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckGlobalThis implements the JSC_USED_GLOBAL_THIS diagnostic.
// Mirrors Java's CheckGlobalThis implements NodeTraversal.Callback.
type CheckGlobalThis struct {
	path string
	// assignLhsChild: if non-nil, the node being traversed is a descendant
	// of the first child of an ASSIGN node. assignLhsChild's parent is this
	// ASSIGN node. Mirrors Java's @Nullable Node assignLhsChild field.
	assignLhsChild *rhino.Node
	diags          []diag.Diagnostic
}

// NewCheckGlobalThis mirrors CheckGlobalThis(AbstractCompiler compiler) in Java.
func NewCheckGlobalThis(path string) *CheckGlobalThis {
	return &CheckGlobalThis{path: path}
}

// Check runs the traversal and returns any JSC_USED_GLOBAL_THIS diagnostics.
func (c *CheckGlobalThis) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.assignLhsChild = nil
	trav := NewNodeTraversal(c, nil, nil)
	trav.Traverse(script)
	return c.diags
}

// ShouldTraverse mirrors CheckGlobalThis.shouldTraverse in Java.
// Since this pass reports errors only when a global `this` keyword is
// encountered, there is no reason to traverse non global contexts.
func (c *CheckGlobalThis) ShouldTraverse(_ *NodeTraversal, n, parent *rhino.Node) bool {
	if n.Token() == rhino.FUNCTION {
		// Arrow functions automatically get the "this" value from the
		// enclosing scope. e.g. the "this" in
		//   Foo.prototype.getBar = () => this.bar;
		// is the global "this", not an instance of Foo.
		if n.GetBooleanProp(rhino.PropArrowFn) {
			return true
		}

		// Don't traverse functions that are constructors or have the @this
		// or @override annotation.
		jsDoc := getBestJSDocInfo(n)
		if jsDoc != nil &&
			(jsDoc.IsConstructor() ||
				jsDoc.IsInterface() ||
				jsDoc.HasThisType() ||
				jsDoc.IsOverride()) {
			return false
		}

		// Don't traverse functions unless they would normally
		// be able to have a @this annotation associated with them. e.g.,
		//   var a = function() { }; // or
		//   function a() {} // or
		//   a.x = function() {}; // or
		//   var a = {x: function() {}};
		if parent == nil {
			return false
		}
		pType := parent.Token()
		if !(pType == rhino.BLOCK ||
			pType == rhino.SCRIPT ||
			pType == rhino.NAME ||
			pType == rhino.ASSIGN ||
			// object literal keys
			pType == rhino.STRING_KEY) {
			return false
		}

		// In our AST, class bodies are BLOCK nodes (children of CLASS).
		// A class method's FUNCTION node has parent=BLOCK, grandparent=CLASS.
		// `this` inside class methods is always the instance or class, not global.
		// Java's NodeUtil.isClassStaticBlock and isMemberFieldDef handle the
		// ES2022 class-body cases; we handle ES6 class methods here.
		if pType == rhino.BLOCK {
			grandparent := parent.Parent()
			if grandparent != nil && grandparent.Token() == rhino.CLASS {
				return false
			}
		}

		// Don't traverse functions that are getting lent to a prototype.
		grandparent := parent.Parent()
		if mayBeObjectLitKey(parent) && grandparent != nil {
			maybeLends := grandparent.JSDocInfo()
			if maybeLends != nil &&
				maybeLends.HasLendsName() &&
				strings.HasSuffix(maybeLends.LendsName().TypeName, ".prototype") {
				return false
			}
		}
	}

	// Don't traverse class static blocks, 'this' is never the global 'this'.
	// NodeUtil.isClassStaticBlock: not in our token set, always returns false.
	// (Java: if (NodeUtil.isClassStaticBlock(n)) return false;)

	// `this` in class field definitions refers to the instance (for non-static
	// fields) or the class itself (for static fields), never the global 'this'.
	// (Java: if ((n.isMemberFieldDef() || n.isComputedFieldDef()) && n.getParent().isClassMembers())
	//   return false;)
	// Not in our token set — safe to skip; would only affect ES2022 class fields.

	if parent != nil && parent.Token() == rhino.ASSIGN {
		lhs := parent.FirstChild()

		if n == lhs {
			// Always traverse the left side of the assignment. To handle
			// nested assignments properly (e.g., (a = this).property = c;),
			// assignLhsChild should not be overridden.
			if c.assignLhsChild == nil {
				c.assignLhsChild = lhs
			}
		} else {
			// Only traverse the right side if it's not an assignment to a
			// prototype property or subproperty.
			if isNormalGet(lhs) {
				if lhs.Token() == rhino.GETPROP && getPropName(lhs) == "prototype" {
					return false
				}
				llhs := lhs.FirstChild()
				if llhs != nil && llhs.Token() == rhino.GETPROP && getPropName(llhs) == "prototype" {
					return false
				}
			}
		}
	}

	return true
}

// Visit mirrors CheckGlobalThis.visit in Java.
func (c *CheckGlobalThis) Visit(_ *NodeTraversal, n, _ *rhino.Node) {
	if n.Token() == rhino.THIS && c.shouldReportThis(n) {
		c.diags = append(c.diags, diag.Diagnostic{
			Path:    c.path,
			Line:    n.Lineno(),
			Col:     n.Charno(),
			Code:    diag.JSC_USED_GLOBAL_THIS,
			// Elevated to ERROR because GLOBAL_THIS belongs to the
			// CHECK_TYPES group, which the oracle enables with
			// --jscomp_error=checkTypes. Mirrors Java's behaviour where
			// DiagnosticGroups.GLOBAL_THIS ⊂ DiagnosticGroups.CHECK_TYPES.
			Level:   diag.Error,
			Message: "dangerous use of the global 'this' object",
		})
	}
	if n == c.assignLhsChild {
		c.assignLhsChild = nil
	}
}

// shouldReportThis mirrors CheckGlobalThis.shouldReportThis in Java.
func (c *CheckGlobalThis) shouldReportThis(n *rhino.Node) bool {
	parent := n.Parent()
	if c.assignLhsChild != nil {
		// Always report a THIS on the left side of an assign.
		return true
	}

	// Also report a THIS with a property access.
	return parent != nil && isNormalGet(parent)
}

// getBestJSDocInfo mirrors NodeUtil.getBestJSDocInfo in Java.
// Climbs from a function node up through ASSIGN / VAR / NAME to find
// the nearest JSDocInfo, matching Java's "function declaration or
// expression whose comment lands on an enclosing node" heuristic.
func getBestJSDocInfo(fn *rhino.Node) *rhino.JSDocInfo {
	if fn == nil {
		return nil
	}
	if info := fn.JSDocInfo(); info != nil {
		return info
	}
	// Climb via ASSIGN → EXPR_RESULT / NAME / VAR / MEMBER_FUNCTION_DEF
	// to find an attached comment.
	for cur := fn.Parent(); cur != nil; cur = cur.Parent() {
		if info := cur.JSDocInfo(); info != nil {
			return info
		}
		switch cur.Token() {
		case rhino.ASSIGN, rhino.EXPR_RESULT, rhino.VAR, rhino.LET, rhino.CONST,
			rhino.NAME, rhino.MEMBER_FUNCTION_DEF, rhino.GETTER_DEF, rhino.SETTER_DEF,
			rhino.STRING_KEY:
			continue
		}
		break
	}
	return nil
}

// isNormalGet mirrors NodeUtil.isNormalGet in Java.
// Returns true for GETPROP (dotted) or GETELEM (computed) accesses.
func isNormalGet(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	t := n.Token()
	return t == rhino.GETPROP || t == rhino.GETELEM
}

// mayBeObjectLitKey mirrors NodeUtil.mayBeObjectLitKey in Java.
// Returns true for STRING_KEY nodes (object-literal property keys).
func mayBeObjectLitKey(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	return n.Token() == rhino.STRING_KEY
}

// getPropName returns the property name from a GETPROP node.
// In our Go AST, GETPROP stores [receiver, STRINGLIT(propName)] as children;
// the GETPROP node's own str field is empty. This mirrors the Java
// Node.getString() which for GETPROP returns the property name.
func getPropName(n *rhino.Node) string {
	if n == nil || n.Token() != rhino.GETPROP {
		return ""
	}
	kids := n.Children()
	if len(kids) < 2 {
		return ""
	}
	return kids[1].String_()
}
