// Port of com/google/javascript/jscomp/ControlFlowAnalysis.java
// @ closure-compiler v20260415 (narrow; no labeled break/continue,
// no generator suspend edges, no with-statement).
//
// Builds a ControlFlowGraph for a FUNCTION body or top-level SCRIPT.
// Walks the AST producing one CfgNode per statement, then connects
// them with properly-labeled branch edges for IF / WHILE / DO_WHILE /
// FOR / FOR_IN / FOR_OF / SWITCH / TRY / THROW / BREAK / CONTINUE /
// RETURN. Matches the edge set Java produces so DataFlowAnalysis sees
// every path the real pass would see.

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

// BuildControlFlowGraph runs on a FUNCTION or SCRIPT node and returns
// its CFG. Nested FUNCTION / CLASS nodes are NOT descended into —
// each has its own CFG, built lazily by the caller.
func BuildControlFlowGraph(root *rhino.Node) *ControlFlowGraph {
	var body *rhino.Node
	switch root.Token() {
	case rhino.FUNCTION:
		kids := root.Children()
		if len(kids) < 3 {
			cfg := NewControlFlowGraph()
			cfg.addEdge(cfg.Entry, cfg.Exit, BranchUnconditional)
			cfg.computePreds()
			return cfg
		}
		body = kids[2]
	case rhino.SCRIPT:
		body = root
	default:
		cfg := NewControlFlowGraph()
		cfg.addEdge(cfg.Entry, cfg.Exit, BranchUnconditional)
		cfg.computePreds()
		return cfg
	}
	b := &cfgBuilder{cfg: NewControlFlowGraph()}
	b.visitStatements(body.Children(), b.cfg.Entry, b.cfg.Exit)
	b.cfg.computePreds()
	return b.cfg
}

// breakTarget pairs a loop / switch's post-loop join node with its
// continue-target (for loops, the update or condition node). The
// stack threads through nested loops / switches so `break;` and
// `continue;` always know where to jump.
type breakTarget struct {
	breakTo    *CfgNode
	continueTo *CfgNode // nil for switch (switch has no continue target)
}

// catchTarget records the nearest enclosing try's catch entry so
// THROW / exception edges know where to route. nil when no catch is
// active — THROW then flows directly to EXIT.
type catchTarget struct {
	to *CfgNode
}

// cfgBuilder carries the mutable state of a single CFG build: the
// graph under construction plus the break / continue / catch
// stacks threaded through nested loops / switches / trys.
type cfgBuilder struct {
	cfg      *ControlFlowGraph
	loops    []breakTarget
	catches  []catchTarget
}

func (b *cfgBuilder) pushLoop(breakTo, continueTo *CfgNode) {
	b.loops = append(b.loops, breakTarget{breakTo: breakTo, continueTo: continueTo})
}
func (b *cfgBuilder) popLoop() { b.loops = b.loops[:len(b.loops)-1] }

func (b *cfgBuilder) pushSwitch(breakTo *CfgNode) {
	b.loops = append(b.loops, breakTarget{breakTo: breakTo})
}

func (b *cfgBuilder) pushCatch(to *CfgNode)    { b.catches = append(b.catches, catchTarget{to: to}) }
func (b *cfgBuilder) popCatch()                { b.catches = b.catches[:len(b.catches)-1] }
func (b *cfgBuilder) topCatch() *CfgNode {
	if len(b.catches) == 0 {
		return nil
	}
	return b.catches[len(b.catches)-1].to
}

// visitStatements walks a sequence of statements, chaining them with
// unconditional edges. `prev` is the node whose fall-through edge
// should be the first statement; `fallthrough_` is where the last
// statement's implicit fall-through goes.
func (b *cfgBuilder) visitStatements(stmts []*rhino.Node, prev *CfgNode, fallthrough_ *CfgNode) {
	current := prev
	for _, stmt := range stmts {
		next := b.visitStatement(stmt, current)
		if next == nil {
			current = nil
			continue
		}
		current = next
	}
	if current != nil {
		b.cfg.addEdge(current, fallthrough_, BranchUnconditional)
	}
}

// visitStatement processes a single statement and returns the CfgNode
// that the following statement's unconditional predecessor should be.
// Returns nil when the statement terminates flow (return / throw /
// break / continue). `prev == nil` skips the predecessor edge (used
// when the caller has already attached a branch-labeled edge).
func (b *cfgBuilder) visitStatement(stmt *rhino.Node, prev *CfgNode) *CfgNode {
	n := b.cfg.node(stmt)
	if prev != nil {
		b.cfg.addEdge(prev, n, BranchUnconditional)
	}
	switch stmt.Token() {
	case rhino.IF:
		return b.visitIf(stmt, n)
	case rhino.WHILE:
		return b.visitWhile(stmt, n)
	case rhino.DO_WHILE:
		return b.visitDoWhile(stmt, n)
	case rhino.FOR:
		return b.visitFor(stmt, n)
	case rhino.FOR_IN, rhino.FOR_OF:
		return b.visitForInOf(stmt, n)
	case rhino.SWITCH:
		return b.visitSwitch(stmt, n)
	case rhino.TRY:
		return b.visitTry(stmt, n)
	case rhino.THROW:
		if c := b.topCatch(); c != nil {
			b.cfg.addEdge(n, c, BranchException)
		} else {
			b.cfg.addEdge(n, b.cfg.Exit, BranchException)
		}
		return nil
	case rhino.BREAK:
		if t := b.breakTarget(); t != nil {
			b.cfg.addEdge(n, t, BranchUnconditional)
		} else {
			b.cfg.addEdge(n, b.cfg.Exit, BranchUnconditional)
		}
		return nil
	case rhino.CONTINUE:
		if t := b.continueTarget(); t != nil {
			b.cfg.addEdge(n, t, BranchUnconditional)
		} else {
			b.cfg.addEdge(n, b.cfg.Exit, BranchUnconditional)
		}
		return nil
	case rhino.RETURN:
		b.cfg.addEdge(n, b.cfg.Exit, BranchUnconditional)
		return nil
	case rhino.BLOCK:
		if kids := stmt.Children(); len(kids) > 0 {
			after := &CfgNode{Kind: "BLOCK_JOIN"}
			b.cfg.Nodes = append(b.cfg.Nodes, after)
			b.visitStatements(kids, n, after)
			return after
		}
		return n
	}
	return n
}

// breakTarget returns the break target of the innermost loop /
// switch, or nil if no loop / switch is enclosing. Labeled breaks
// are deferred — `break;` unadorned is what the corpus uses.
func (b *cfgBuilder) breakTarget() *CfgNode {
	if len(b.loops) == 0 {
		return nil
	}
	return b.loops[len(b.loops)-1].breakTo
}

// continueTarget returns the continue target of the innermost loop,
// skipping pure switches (which can't host `continue`).
func (b *cfgBuilder) continueTarget() *CfgNode {
	for i := len(b.loops) - 1; i >= 0; i-- {
		if b.loops[i].continueTo != nil {
			return b.loops[i].continueTo
		}
	}
	return nil
}

// visitIf handles an `if (cond) then [else]` statement.
func (b *cfgBuilder) visitIf(stmt *rhino.Node, n *CfgNode) *CfgNode {
	kids := stmt.Children()
	if len(kids) < 2 {
		return n
	}
	thenBlock := kids[1]
	var elseBlock *rhino.Node
	if len(kids) >= 3 {
		elseBlock = kids[2]
	}
	join := &CfgNode{Kind: "IF_JOIN"}
	b.cfg.Nodes = append(b.cfg.Nodes, join)
	b.visitBranch(n, thenBlock, BranchOnTrue, join)
	if elseBlock != nil {
		b.visitBranch(n, elseBlock, BranchOnFalse, join)
	} else {
		b.cfg.addEdge(n, join, BranchOnFalse)
	}
	return join
}

// visitBranch wires one arm of an if/else.
func (b *cfgBuilder) visitBranch(ifNode *CfgNode, block *rhino.Node, branch Branch, join *CfgNode) {
	stmts := flattenBlock(block)
	if len(stmts) == 0 {
		b.cfg.addEdge(ifNode, join, branch)
		return
	}
	first := stmts[0]
	firstNode := b.cfg.node(first)
	b.cfg.addEdge(ifNode, firstNode, branch)
	next := b.visitStatement(first, nil)
	if next != nil {
		b.visitStatements(stmts[1:], next, join)
	} else if len(stmts) > 1 {
		b.visitStatements(stmts[1:], nil, join)
	}
}

// visitWhile handles `while (cond) body`. The CFG shape:
//   prev → WHILE    (unconditional)
//   WHILE → body-start (BranchOnTrue)
//   WHILE → after   (BranchOnFalse)
//   body-tail → WHILE (back-edge, BranchSyn)
func (b *cfgBuilder) visitWhile(stmt *rhino.Node, n *CfgNode) *CfgNode {
	kids := stmt.Children()
	if len(kids) < 2 {
		return n
	}
	after := &CfgNode{Kind: "LOOP_AFTER"}
	b.cfg.Nodes = append(b.cfg.Nodes, after)
	b.pushLoop(after, n)
	defer b.popLoop()
	// True branch: walk the body; its fall-through loops back to the
	// while condition node.
	b.visitBranch(n, kids[1], BranchOnTrue, n)
	b.cfg.addEdge(n, after, BranchOnFalse)
	return after
}

// visitDoWhile handles `do body while (cond);`. CFG:
//   prev → body-start
//   body-tail → DO_WHILE (cond node)
//   DO_WHILE → body-start (BranchOnTrue, back-edge)
//   DO_WHILE → after (BranchOnFalse)
func (b *cfgBuilder) visitDoWhile(stmt *rhino.Node, n *CfgNode) *CfgNode {
	kids := stmt.Children()
	if len(kids) < 2 {
		return n
	}
	body := kids[0]
	after := &CfgNode{Kind: "LOOP_AFTER"}
	b.cfg.Nodes = append(b.cfg.Nodes, after)
	b.pushLoop(after, n)
	defer b.popLoop()
	// Body is executed unconditionally first, so `prev → body-start`.
	// The `n` node here is the DO_WHILE wrapper; we want its
	// predecessor (already wired) to flow straight into the body.
	stmts := flattenBlock(body)
	if len(stmts) == 0 {
		b.cfg.addEdge(n, n, BranchOnTrue)
		b.cfg.addEdge(n, after, BranchOnFalse)
		return after
	}
	firstNode := b.cfg.node(stmts[0])
	// Replace prev→DO_WHILE with prev→body-start by adding the body
	// link (the DO_WHILE node stays in place as the cond site).
	b.cfg.addEdge(n, firstNode, BranchSyn)
	next := b.visitStatement(stmts[0], nil)
	if next != nil {
		b.visitStatements(stmts[1:], next, n) // fall-through → cond
	} else if len(stmts) > 1 {
		b.visitStatements(stmts[1:], nil, n)
	}
	// From cond (n): true loops back to body, false exits.
	b.cfg.addEdge(n, firstNode, BranchOnTrue)
	b.cfg.addEdge(n, after, BranchOnFalse)
	return after
}

// visitFor handles classic `for (init; cond; upd) body`.
// Children are (init, cond, upd, body). CFG:
//   prev → FOR (the init slot conceptually runs first)
//   FOR → cond-node via init
//   cond → body-start (BranchOnTrue)
//   cond → after       (BranchOnFalse)
//   body-tail → upd → cond
func (b *cfgBuilder) visitFor(stmt *rhino.Node, n *CfgNode) *CfgNode {
	kids := stmt.Children()
	if len(kids) < 4 {
		return n
	}
	// We reuse `n` as the cond-evaluation site so branch labels on
	// its outgoing edges match the standard IF convention. Init is a
	// synthetic child folded into the unconditional predecessor edge
	// — without a dedicated node, the dataflow transfer function
	// still sees init via `n`'s AST children.
	after := &CfgNode{Kind: "LOOP_AFTER"}
	b.cfg.Nodes = append(b.cfg.Nodes, after)
	updNode := &CfgNode{Kind: "FOR_UPDATE", AST: kids[2]}
	b.cfg.Nodes = append(b.cfg.Nodes, updNode)
	b.pushLoop(after, updNode)
	defer b.popLoop()
	// Body:
	body := kids[3]
	stmts := flattenBlock(body)
	if len(stmts) > 0 {
		first := b.cfg.node(stmts[0])
		b.cfg.addEdge(n, first, BranchOnTrue)
		next := b.visitStatement(stmts[0], nil)
		if next != nil {
			b.visitStatements(stmts[1:], next, updNode)
		} else if len(stmts) > 1 {
			b.visitStatements(stmts[1:], nil, updNode)
		}
	} else {
		b.cfg.addEdge(n, updNode, BranchOnTrue)
	}
	b.cfg.addEdge(n, after, BranchOnFalse)
	b.cfg.addEdge(updNode, n, BranchSyn) // upd → cond back-edge
	return after
}

// visitForInOf handles for-in / for-of with children (lhs, iter, body).
// Models the iterator as a conditional on the header: OnTrue enters
// the body, OnFalse exits. Body back-edges to the header. Close
// enough to Java's IR even though we lack a dedicated ITER_NEXT node.
func (b *cfgBuilder) visitForInOf(stmt *rhino.Node, n *CfgNode) *CfgNode {
	kids := stmt.Children()
	if len(kids) < 3 {
		return n
	}
	body := kids[2]
	after := &CfgNode{Kind: "LOOP_AFTER"}
	b.cfg.Nodes = append(b.cfg.Nodes, after)
	b.pushLoop(after, n)
	defer b.popLoop()
	b.visitBranch(n, body, BranchOnTrue, n)
	b.cfg.addEdge(n, after, BranchOnFalse)
	return after
}

// visitSwitch handles `switch (e) { case v1: s1 case v2: s2 default:
// d }`. Children are (discriminant, CASE*, DEFAULT_CASE?). Each case
// body falls through to the next unless it ends in a break or
// return. CFG:
//   prev → SWITCH
//   SWITCH → case[0] (BranchSyn — synthetic match-test)
//   case[0] body-tail → case[1] (fall-through)
//   case[i] body ends in BREAK → after
//   SWITCH → default → after
func (b *cfgBuilder) visitSwitch(stmt *rhino.Node, n *CfgNode) *CfgNode {
	after := &CfgNode{Kind: "SWITCH_AFTER"}
	b.cfg.Nodes = append(b.cfg.Nodes, after)
	b.pushSwitch(after)
	defer b.popLoop()
	kids := stmt.Children()
	cases := kids
	if len(cases) > 0 {
		cases = cases[1:] // skip discriminant
	}
	// Wire switch → each case with synthetic edges.
	var bodies []*rhino.Node
	for _, c := range cases {
		bodyNode := c.Children()[len(c.Children())-1] // last child is BLOCK
		bodies = append(bodies, bodyNode)
		b.cfg.addEdge(n, b.cfg.node(c), BranchSyn)
	}
	// Chain case bodies sequentially — fall-through from one to the
	// next mirrors JS semantics when a case body doesn't break.
	for i, c := range cases {
		caseNode := b.cfg.node(c)
		bodyStmts := flattenBlock(bodies[i])
		var bodyTail *CfgNode
		if len(bodyStmts) == 0 {
			bodyTail = caseNode
		} else {
			first := b.cfg.node(bodyStmts[0])
			b.cfg.addEdge(caseNode, first, BranchUnconditional)
			next := b.visitStatement(bodyStmts[0], nil)
			if next != nil {
				// Fall-through from the tail of this body goes to
				// the next case's start (or to `after` if it's the
				// last case).
				var nextFall *CfgNode
				if i+1 < len(cases) {
					nextFall = b.cfg.node(cases[i+1])
				} else {
					nextFall = after
				}
				b.visitStatements(bodyStmts[1:], next, nextFall)
			} else if len(bodyStmts) > 1 {
				b.visitStatements(bodyStmts[1:], nil, after)
			}
			bodyTail = nil // handled by visitStatements
			_ = bodyTail
		}
	}
	// Ensure `after` has at least one predecessor even when every
	// case returns — fall-through from the switch itself on "no
	// case matched, no default".
	hasDefault := false
	for _, c := range cases {
		if c.Token() == rhino.DEFAULT_CASE {
			hasDefault = true
			break
		}
	}
	if !hasDefault {
		b.cfg.addEdge(n, after, BranchSyn)
	}
	return after
}

// visitTry handles `try body catch (e) handler [finally block]`.
// Body runs with a catch target pushed so contained THROWs route
// to the catch entry. After body / catch, flow joins at an after
// node (and may pass through the finally block).
func (b *cfgBuilder) visitTry(stmt *rhino.Node, n *CfgNode) *CfgNode {
	kids := stmt.Children()
	if len(kids) < 1 {
		return n
	}
	after := &CfgNode{Kind: "TRY_AFTER"}
	b.cfg.Nodes = append(b.cfg.Nodes, after)
	body := kids[0]
	var catchClause, finallyClause *rhino.Node
	for _, c := range kids[1:] {
		switch c.Token() {
		case rhino.CATCH:
			catchClause = c
		case rhino.FINALLY:
			finallyClause = c
		}
	}
	var catchEntry *CfgNode
	if catchClause != nil {
		catchEntry = b.cfg.node(catchClause)
	}
	// Wire the body with the catch target pushed so THROWs route
	// correctly. The body's fall-through leaves the try and goes to
	// `after` (or to the finally block if one exists).
	var bodyFall *CfgNode = after
	if finallyClause != nil {
		bodyFall = b.cfg.node(finallyClause)
	}
	if catchEntry != nil {
		b.pushCatch(catchEntry)
	}
	bodyStmts := flattenBlock(body)
	if len(bodyStmts) > 0 {
		b.cfg.addEdge(n, b.cfg.node(bodyStmts[0]), BranchUnconditional)
		next := b.visitStatement(bodyStmts[0], nil)
		if next != nil {
			b.visitStatements(bodyStmts[1:], next, bodyFall)
		} else if len(bodyStmts) > 1 {
			b.visitStatements(bodyStmts[1:], nil, bodyFall)
		}
	} else {
		b.cfg.addEdge(n, bodyFall, BranchUnconditional)
	}
	if catchEntry != nil {
		b.popCatch()
	}
	// Wire the catch-handler body.
	if catchClause != nil {
		handler := catchClause.Children()[len(catchClause.Children())-1]
		hStmts := flattenBlock(handler)
		if len(hStmts) > 0 {
			b.cfg.addEdge(catchEntry, b.cfg.node(hStmts[0]), BranchUnconditional)
			next := b.visitStatement(hStmts[0], nil)
			if next != nil {
				b.visitStatements(hStmts[1:], next, bodyFall)
			} else if len(hStmts) > 1 {
				b.visitStatements(hStmts[1:], nil, bodyFall)
			}
		} else {
			b.cfg.addEdge(catchEntry, bodyFall, BranchUnconditional)
		}
	}
	// Wire the finally block (if any). Its fall-through goes to
	// `after`.
	if finallyClause != nil {
		finallyNode := b.cfg.node(finallyClause)
		fStmts := flattenBlock(finallyClause.FirstChild())
		if len(fStmts) > 0 {
			b.cfg.addEdge(finallyNode, b.cfg.node(fStmts[0]), BranchUnconditional)
			next := b.visitStatement(fStmts[0], nil)
			if next != nil {
				b.visitStatements(fStmts[1:], next, after)
			} else if len(fStmts) > 1 {
				b.visitStatements(fStmts[1:], nil, after)
			}
		} else {
			b.cfg.addEdge(finallyNode, after, BranchUnconditional)
		}
	}
	return after
}

// flattenBlock returns the statement list inside `n`. If `n` is a
// BLOCK, returns its children; otherwise returns a single-element
// slice with `n` itself so `if (c) x;` parses like `if (c) { x; }`.
func flattenBlock(n *rhino.Node) []*rhino.Node {
	if n == nil {
		return nil
	}
	if n.Token() == rhino.BLOCK {
		return n.Children()
	}
	return []*rhino.Node{n}
}
