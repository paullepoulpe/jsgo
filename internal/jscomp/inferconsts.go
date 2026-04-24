// Port of com/google/javascript/jscomp/InferConsts.java @ closure-compiler v20260415

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InferConsts attaches the CONST_VAR annotation to any variable that is:
//  1. Provably well-defined and assigned once in its lifetime.
//  2. Annotated @const.
//  3. Declared with the const keyword.
//
// These three are considered semantically equivalent. A variable in a loop
// is never considered const.
//
// Note that criterion (1) is only used for normal code, not externs.
type InferConsts struct{}

// Go-only: NewInferConsts mirrors InferConsts(AbstractCompiler compiler).
// Java has no factory function; the object is constructed directly.
func NewInferConsts() *InferConsts {
	return &InferConsts{}
}

// Check runs the pass over script, annotating NAME nodes.
// Returns no diagnostics — this is a pure annotation pass.
func (ic *InferConsts) Check(script *rhino.Node) {
	// Collect assignment counts for every name in the script, so we can
	// determine whether a binding is assigned exactly once.
	counts := icCountAssignments(script)
	icAnnotate(script, counts)
}

// Go-only: icCountAssignments walks n and returns a map from identifier name to the
// number of times it appears as an assignment target (lvalue) in the tree.
// Declaration sites that also initialise (var x = 1, let x = 1, const x = 1)
// count as one assignment. Plain reassignments (x = 2) add another. The
// function-declaration form (function f(){}) counts the hoisted name as one
// assignment.
func icCountAssignments(root *rhino.Node) map[string]int {
	counts := map[string]int{}
	icWalkForCounts(root, counts, false)
	return counts
}

// Go-only: icWalkForCounts is the recursive worker for icCountAssignments.
// inLoop tracks whether the current position is inside a loop body, so that
// a declaration inside a loop does not qualify as single-assignment.
func icWalkForCounts(n *rhino.Node, counts map[string]int, inLoop bool) {
	if n == nil {
		return
	}

	loopNext := inLoop
	switch n.Token() {
	case rhino.WHILE, rhino.DO_WHILE, rhino.FOR, rhino.FOR_IN, rhino.FOR_OF:
		loopNext = true
	}

	switch n.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		for _, child := range n.Children() {
			if child.Token() == rhino.NAME {
				name := child.String_()
				if name != "" {
					if inLoop {
						// Inside a loop the binding is effectively re-executed every
						// iteration; treat as multi-assignment.
						counts[name] += 2
					} else if child.HasChildren() {
						// Initialised declaration counts as one assignment.
						counts[name]++
					} else {
						// Uninitialised declaration (var x;) — no value assigned yet.
						// We still record the name with a zero delta so it exists.
						if _, seen := counts[name]; !seen {
							counts[name] = 0
						}
					}
				}
			}
			// Recurse into initialisers for nested references.
			icWalkForCounts(child.FirstChild(), counts, loopNext)
		}
		return

	case rhino.FUNCTION:
		kids := n.Children()
		if len(kids) >= 1 {
			nameNode := kids[0]
			if nameNode.Token() == rhino.NAME && nameNode.String_() != "" {
				// Function declarations bind the name in the enclosing scope.
				parent := n.Parent()
				if parent != nil && (parent.Token() == rhino.BLOCK || parent.Token() == rhino.SCRIPT) {
					name := nameNode.String_()
					if inLoop {
						counts[name] += 2
					} else {
						counts[name]++
					}
				}
			}
		}
		// Recurse into the function body (new scope — do not carry inLoop).
		for _, child := range n.Children() {
			icWalkForCounts(child, counts, false)
		}
		return

	case rhino.CATCH:
		kids := n.Children()
		if len(kids) >= 1 {
			nameNode := kids[0]
			if nameNode.Token() == rhino.NAME && nameNode.String_() != "" {
				name := nameNode.String_()
				if inLoop {
					counts[name] += 2
				} else {
					counts[name]++
				}
			}
		}

	case rhino.PARAM_LIST:
		// Parameters are declared once per call — not loop-sensitive from
		// the InferConsts perspective (they get a fresh binding each call).
		for _, child := range n.Children() {
			if child.Token() == rhino.NAME && child.String_() != "" {
				counts[child.String_()]++
			}
		}

	case rhino.ASSIGN:
		// x = expr — the LHS name is an additional assignment.
		kids := n.Children()
		if len(kids) >= 1 {
			lhs := kids[0]
			if lhs.Token() == rhino.NAME && lhs.String_() != "" {
				counts[lhs.String_()]++
			}
		}
	}

	for _, child := range n.Children() {
		icWalkForCounts(child, counts, loopNext)
	}
}

// Go-only: icAnnotate walks n and sets DeclaredConstantVar / InferredConstantVar on
// every NAME node that is a declaration binding site.
func icAnnotate(n *rhino.Node, counts map[string]int) {
	if n == nil {
		return
	}

	switch n.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		for _, child := range n.Children() {
			if child.Token() == rhino.NAME {
				icConsiderNameNode(child, n, counts)
			}
		}

	case rhino.FUNCTION:
		kids := n.Children()
		if len(kids) >= 1 {
			nameNode := kids[0]
			if nameNode.Token() == rhino.NAME && nameNode.String_() != "" {
				parent := n.Parent()
				if parent != nil && (parent.Token() == rhino.BLOCK || parent.Token() == rhino.SCRIPT) {
					icConsiderNameNode(nameNode, n, counts)
				}
			}
		}

	case rhino.CATCH:
		kids := n.Children()
		if len(kids) >= 1 {
			nameNode := kids[0]
			if nameNode.Token() == rhino.NAME && nameNode.String_() != "" {
				icConsiderNameNode(nameNode, n, counts)
			}
		}

	case rhino.PARAM_LIST:
		for _, child := range n.Children() {
			if child.Token() == rhino.NAME && child.String_() != "" {
				icConsiderNameNode(child, n, counts)
			}
		}
	}

	for _, child := range n.Children() {
		icAnnotate(child, counts)
	}
}

// Go-only: icConsiderNameNode mirrors InferConsts.considerVar for a single NAME node.
// Renamed with ic prefix to avoid collisions; Java's method is considerVar.
// nameNode is the binding NAME; declParent is its declaration parent
// (VAR/LET/CONST/FUNCTION/CATCH/PARAM_LIST). counts maps name → total
// assignment count in the script.
func icConsiderNameNode(nameNode, declParent *rhino.Node, counts map[string]int) {
	// Check @const JSDoc: the JSDoc may be on the nameNode or on the
	// declaration statement (var/let/const). Mirrors v.getJSDocInfo() in Java
	// which walks up to the var/let/const node.
	docInfo := nameNode.JSDocInfo()
	if docInfo == nil {
		docInfo = declParent.JSDocInfo()
		if docInfo == nil {
			// Also check the enclosing EXPR_RESULT / statement for @const on
			// standalone assignment forms.
			p := declParent.Parent()
			if p != nil {
				docInfo = p.JSDocInfo()
			}
		}
	}

	// Criterion 2: @const JSDoc annotation.
	if docInfo != nil && docInfo.IsConstant() {
		nameNode.SetDeclaredConstantVar(true)
	}

	// Criterion 3: const keyword.
	if declParent.Token() == rhino.CONST {
		nameNode.SetDeclaredConstantVar(true)
	}

	// Criterion 1: inferred const — assigned exactly once.
	if icIsInferredConst(nameNode, declParent, counts) {
		nameNode.SetInferredConstantVar(true)
	}
}

// Go-only: icIsInferredConst mirrors InferConsts.isInferredConst.
// Prefixed with ic to avoid name collision; Java method is isInferredConst.
// Returns true when the variable can be inferred to be assigned exactly once.
func icIsInferredConst(nameNode, declParent *rhino.Node, counts map[string]int) bool {
	name := nameNode.String_()
	if name == "" {
		return false
	}

	switch declParent.Token() {
	case rhino.LET:
		// Non-destructuring let names must be assigned at declaration.
		// Also require single-assignment in the whole script.
		if !nameNode.HasChildren() {
			return false
		}
		return counts[name] <= 1

	case rhino.CONST:
		// const always has exactly one assignment (the declaration).
		return true

	case rhino.CATCH:
		// catch bindings are assigned once (the caught exception).
		return true

	case rhino.CLASS:
		// Class declarations bind the name once.
		return true

	case rhino.PARAM_LIST:
		// Parameters cannot be referenced before declaration.
		return true

	case rhino.FUNCTION:
		// Function hoisting means no references before declaration.
		return true

	case rhino.VAR:
		// var hoisting: require the declaration also initialises the name
		// (first reference is an assigning declaration) and is well-defined
		// (single assignment in the whole scope).
		if !nameNode.HasChildren() {
			return false
		}
		return counts[name] <= 1

	default:
		return false
	}
}
