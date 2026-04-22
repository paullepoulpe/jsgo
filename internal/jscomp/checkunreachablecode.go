// Port of com/google/javascript/jscomp/CheckUnreachableCode.java @ closure-compiler v20260415
//
// Use ControlFlowGraph and GraphReachability to inform user about unreachable code.
//
// Java's pass extends AbstractCfgCallback which fires enterScopeWithCfg per
// FUNCTION / SCRIPT root, builds a CFG via ControlFlowAnalysis, runs
// GraphReachability from the ENTRY node, then uses shouldTraverse to detect
// any CFG node whose reachability annotation is not REACHABLE. When found it
// reports JSC_UNREACHABLE_CODE and calls GraphReachability.recompute(n) so
// sub-nodes of the newly-unreachable section don't cascade into further
// warnings.
//
// Go port strategy: we walk the AST finding each valid CFG root
// (FUNCTION / SCRIPT). Per root we build the CFG, compute reachability,
// then walk the CFG node set to report unreachable statement-level nodes.
// We suppress cascades the same way Java does — once a node is reported,
// we re-mark it reachable (recompute) so its successors stop being flagged.
//
// TODO: widen cfg to carry a per-node annotation slot so GraphReachability
// can be written as a generic annotating BFS matching the Java API exactly
// (GraphNode.setAnnotation / getAnnotation). Currently we inline the BFS as
// a reachable-set map which is behaviourally equivalent.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckUnreachableCode implements the JSC_UNREACHABLE_CODE diagnostic.
// Mirrors the class-level structure of Java's CheckUnreachableCode which
// holds only a reference to the compiler.
type CheckUnreachableCode struct {
	path  string
	diags []diag.Diagnostic
}

// NewCheckUnreachableCode constructs the pass. `path` is the source file
// path used in diagnostic positions, matching how Java passes `compiler`
// (which carries the current source file context).
func NewCheckUnreachableCode(path string) *CheckUnreachableCode {
	return &CheckUnreachableCode{path: path}
}

// Check walks `script`, finds each valid CFG root (FUNCTION or SCRIPT),
// builds its CFG, computes reachability, and reports unreachable nodes.
// Mirrors the AbstractCfgCallback lifecycle: enterScopeWithCfg per root,
// then shouldTraverse per node.
func (c *CheckUnreachableCode) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.walkRoot(script)
	return c.diags
}

// walkRoot handles one CFG root (FUNCTION or SCRIPT): builds CFG,
// computes reachability, checks nodes, then descends into nested
// functions. Matches the AbstractCfgCallback's per-scope lifecycle.
func (c *CheckUnreachableCode) walkRoot(root *rhino.Node) {
	if root == nil {
		return
	}
	// isValidCfgRoot: Java's NodeUtil.isValidCfgRoot accepts FUNCTION and SCRIPT.
	if root.Token() == rhino.FUNCTION || root.Token() == rhino.SCRIPT {
		cfg := BuildControlFlowGraph(root)
		// initScope: run GraphReachability.compute from the entry node,
		// using isReachable to skip definitely-unreachable conditional branches.
		reachable := c.computeReachability(cfg)
		// shouldTraverse pass: walk statement nodes in the CFG, report
		// unreachable ones, re-expand their successors into reachable after
		// reporting (matching Java's GraphReachability.recompute(n) call).
		c.checkReachability(cfg, reachable)
	}
	// Descend into children to find nested FUNCTION nodes — each has its own CFG.
	// For FUNCTION roots we skip the function node itself since it was just processed,
	// but we still need to descend into its body to find nested functions.
	for _, child := range root.Children() {
		c.walkChildren(child)
	}
}

// walkChildren descends into nested FUNCTION nodes so each gets its own CFG pass.
func (c *CheckUnreachableCode) walkChildren(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.FUNCTION {
		// This is a new CFG root — process it.
		cfg := BuildControlFlowGraph(n)
		reachable := c.computeReachability(cfg)
		c.checkReachability(cfg, reachable)
		// Descend into the function's children to find nested functions.
		for _, child := range n.Children() {
			c.walkChildren(child)
		}
		return
	}
	for _, child := range n.Children() {
		c.walkChildren(child)
	}
}

// computeReachability implements GraphReachability.compute: BFS/DFS from
// the ENTRY node, using isReachable to prune edges where the branch
// condition is provably false. Returns the set of reachable CfgNodes.
//
// Mirrors Java's GraphReachability<Node, Branch> with the custom
// isReachable predicate from CheckUnreachableCode.isReachable.
func (c *CheckUnreachableCode) computeReachability(cfg *ControlFlowGraph) map[*CfgNode]bool {
	reachable := map[*CfgNode]bool{}
	if cfg == nil || cfg.Entry == nil {
		return reachable
	}
	stack := []*CfgNode{cfg.Entry}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if reachable[n] {
			continue
		}
		reachable[n] = true
		for _, e := range n.Successors {
			if !reachable[e.To] && c.isReachableEdge(n, e.Branch) {
				stack = append(stack, e.To)
			}
		}
	}
	return reachable
}

// isReachableEdge mirrors Java's CheckUnreachableCode.isReachable(EdgeTuple).
// For unconditional branches it always returns true. For conditional branches
// it evaluates whether the branch condition is a compile-time boolean constant
// that makes this particular arm provably unreachable.
//
// Port of:
//
//	private static boolean isReachable(EdgeTuple<Node, Branch> input) {
//	  Branch branch = input.edge;
//	  if (!branch.isConditional()) { return true; }
//	  Node predecessor = input.sourceNode;
//	  Node condition = NodeUtil.getConditionExpression(predecessor);
//	  if (condition != null) {
//	    Tri val = NodeUtil.getBooleanValue(condition);
//	    if (val != Tri.UNKNOWN) {
//	      return val.toBoolean(true) == (branch == Branch.ON_TRUE);
//	    }
//	  }
//	  return true;
//	}
//
// TODO(user): Handle more complicated expression like true == true, etc....
func (c *CheckUnreachableCode) isReachableEdge(pred *CfgNode, branch Branch) bool {
	if !isConditionalBranch(branch) {
		return true
	}
	if pred.AST == nil {
		return true
	}
	condition := getConditionExpression(pred.AST)
	if condition == nil {
		return true
	}
	val, known := getBooleanValueKnown(condition)
	if !known {
		return true
	}
	// The edge is reachable iff the static value matches the branch direction.
	return val == (branch == BranchOnTrue)
}

// isConditionalBranch reports whether `b` is a conditional branch (ON_TRUE
// or ON_FALSE). Mirrors Java's Branch.isConditional().
func isConditionalBranch(b Branch) bool {
	return b == BranchOnTrue || b == BranchOnFalse
}

// getConditionExpression returns the condition expression node for a
// statement that can have a conditional branch in the CFG. Mirrors
// Java's NodeUtil.getConditionExpression.
//
// IF(cond, then, else?) → cond (child[0])
// WHILE(cond, body)      → cond (child[0])
// DO_WHILE(body, cond)   → cond (child[1])  [note reversed order vs WHILE]
// FOR(init, cond, upd, body) → cond (child[1])
// FOR_IN / FOR_OF        → nil (no static condition)
func getConditionExpression(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	kids := n.Children()
	switch n.Token() {
	case rhino.IF, rhino.WHILE:
		if len(kids) >= 1 {
			return kids[0]
		}
	case rhino.DO_WHILE:
		// DO_WHILE children: [body, condition]
		if len(kids) >= 2 {
			return kids[1]
		}
	case rhino.FOR:
		// FOR children: [init, cond, update, body]
		if len(kids) >= 2 {
			return kids[1]
		}
	}
	return nil
}

// getBooleanValueKnown mirrors Java's NodeUtil.getBooleanValueKnown (Tri-valued).
// Returns (bool value, known) where known==false means Tri.UNKNOWN.
//
// TODO(user): Handle more complicated expression like true == true, etc....
func getBooleanValueKnown(n *rhino.Node) (bool, bool) {
	if n == nil {
		return false, false
	}
	switch n.Token() {
	case rhino.TRUE:
		return true, true
	case rhino.FALSE:
		return false, true
	case rhino.NULL:
		return false, true
	case rhino.VOID:
		// void <anything> is undefined, which is falsy
		return false, true
	case rhino.NOT:
		if child := n.FirstChild(); child != nil {
			v, known := getBooleanValueKnown(child)
			if known {
				return !v, true
			}
		}
	case rhino.AND:
		kids := n.Children()
		if len(kids) == 2 {
			lv, lk := getBooleanValueKnown(kids[0])
			if lk && !lv {
				return false, true // false && anything = false
			}
			rv, rk := getBooleanValueKnown(kids[1])
			if lk && lv && rk {
				return rv, true // true && rv
			}
		}
	case rhino.OR:
		kids := n.Children()
		if len(kids) == 2 {
			lv, lk := getBooleanValueKnown(kids[0])
			if lk && lv {
				return true, true // true || anything = true
			}
			rv, rk := getBooleanValueKnown(kids[1])
			if lk && !lv && rk {
				return rv, true // false || rv
			}
		}
	case rhino.NUMBER:
		// Non-zero numbers are truthy, zero is falsy.
		// Java's getBooleanValueKnown checks the double value.
		if str := n.String_(); str == "0" || str == "0.0" || str == "-0" {
			return false, true
		}
		if str := n.String_(); str != "" {
			return true, true
		}
	case rhino.STRINGLIT:
		return n.String_() != "", true
	}
	return false, false
}

// checkReachability walks the CFG nodes and reports those that are not
// reachable. When it reports a node, it calls recompute to expand that
// node's successors back into the reachable set — this prevents cascade
// warnings on code the user presumably intended to have (they just need
// to fix the dead code before it). Mirrors Java's shouldTraverse callback
// plus the GraphReachability.recompute(n) call.
func (c *CheckUnreachableCode) checkReachability(cfg *ControlFlowGraph, reachable map[*CfgNode]bool) {
	for _, cfgNode := range cfg.Nodes {
		// Only real AST statement nodes participate (not synthetic ENTRY/EXIT/JOIN).
		if cfgNode.AST == nil {
			continue
		}
		if reachable[cfgNode] {
			continue
		}
		n := cfgNode.AST
		// Only report error when there are some line number informations.
		// There are synthetic nodes with no line number informations, nodes
		// introduce by other passes (although not likely since this pass should
		// be executed early) or some rhino bug.
		// Java uses -1; Go uses 0 for "unset" since lineno is 1-indexed.
		if n.Lineno() <= 0 {
			continue
		}
		// Allow spurious semi-colons (EMPTY token in Java; not present in the
		// Go parser subset — the parser elides empty statements) and spurious breaks.
		// TODO: widen cfg to add EMPTY token to rhino.Token when the parser
		// produces empty-statement nodes.
		if n.Token() == rhino.BREAK {
			continue
		}
		c.diags = append(c.diags, diag.Diagnostic{
			Path:    c.path,
			Line:    n.Lineno(),
			Col:     n.Charno(),
			Code:    diag.JSC_UNREACHABLE_CODE,
			Level:   diag.Warning,
			Message: "unreachable code",
		})
		// From now on, we are going to assume the user fixed the error and not
		// give more warning related to code section reachable from this node.
		// Mirrors: new GraphReachability<>(cfg).recompute(n)
		c.recompute(cfg, cfgNode, reachable)
	}
}

// recompute marks `start` and all nodes reachable from it as reachable in
// the given set. This mirrors Java's GraphReachability.recompute(n) which
// does a BFS/DFS from `n` through the CFG's forward edges, re-annotating
// visited nodes as REACHABLE so they don't trigger further warnings.
func (c *CheckUnreachableCode) recompute(cfg *ControlFlowGraph, start *CfgNode, reachable map[*CfgNode]bool) {
	stack := []*CfgNode{start}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if reachable[n] {
			continue
		}
		reachable[n] = true
		for _, e := range n.Successors {
			if !reachable[e.To] {
				stack = append(stack, e.To)
			}
		}
	}
}
