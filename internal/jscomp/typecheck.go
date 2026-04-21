// Port of com/google/javascript/jscomp/TypeCheck.java
// @ closure-compiler v20260415 (subset).
//
// TypeCheck implements the NodeTraversal.Callback contract and
// dispatches per-node in Visit — the same post-order visitor shape
// Java's class uses (extends AbstractPostOrderCallback, overrides
// visit(NodeTraversal, Node, Node)). Each visit-method covers one
// node kind and is a direct analogue of Java's visitCall /
// visitNew / visitAssign / visitGetProp / visitFunction /
// visitClassField dispatch.
//
// Diagnostics currently emitted:
//   - "initializing variable" JSC_TYPE_MISMATCH on var/let/const
//   - "assignment" / "assignment to property X of Y" JSC_TYPE_MISMATCH
//   - "actual parameter N of F does not match formal parameter"
//     JSC_TYPE_MISMATCH on call + `new` sites
//   - "inconsistent return type" JSC_TYPE_MISMATCH on return statements
//   - JSC_INEXISTENT_PROPERTY on structural receivers with closed
//     property sets

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

type TypeCheck struct {
	creator *TypedScopeCreator
	global  *TypedScope
	path    string
	opts    *CompilerOptions
	diags   []diag.Diagnostic

	// SuppressedNames mirrors Java's TypeCheck skipping of
	// @implements / @interface-derived diagnostics when the relevant
	// symbol slot was already rejected by an earlier VarCheck pass
	// (Java implements this via the slot lookup naturally failing; we
	// thread the rejected-name set in explicitly).
	SuppressedNames map[string]bool

	// inexistentReported dedups JSC_INEXISTENT_PROPERTY emissions
	// across the script — Java fires the diagnostic only once per
	// (receiver, property) pair even when the program references
	// the missing slot multiple times. Keeps message quality high
	// without tracking per-guard refinement state.
	inexistentReported map[string]bool

	// fileSuppressions carries any @fileoverview @suppress {...}
	// names captured pre-traversal. visitGetProp / visitReturn /
	// etc. consult it to skip diagnostics the author opted out
	// of file-wide.
	fileSuppressions map[string]bool

	// flows maps each function / script root to its computed
	// DataFlowAnalysis, built lazily as the traversal enters the
	// function. statementAtNode maps a CFG-statement AST node to its
	// owning flow so TypeCheck can consult the IN FlowScope when
	// checking expressions inside it.
	flows       map[*rhino.Node]*DataFlowAnalysis
	currentFlow *DataFlowAnalysis
	currentCfg  *ControlFlowGraph
	flowStack   []*DataFlowAnalysis
	cfgStack    []*ControlFlowGraph
	currentStmt *rhino.Node
	stmtStack   []*rhino.Node
}

func NewTypeCheck(path string, creator *TypedScopeCreator, global *TypedScope) *TypeCheck {
	return NewTypeCheckWithOptions(path, creator, global, NewCompilerOptions())
}

// NewTypeCheckWithOptions is the per-options constructor, mirroring
// Java's TypeCheck(AbstractCompiler, …). The options bag routes
// per-DiagnosticGroup CheckLevel so strict-only validators stay
// silent unless `strictCheckTypes` is set to WARNING / ERROR.
func NewTypeCheckWithOptions(path string, creator *TypedScopeCreator, global *TypedScope, opts *CompilerOptions) *TypeCheck {
	if opts == nil {
		opts = NewCompilerOptions()
	}
	return &TypeCheck{
		creator: creator,
		global:  global,
		path:    path,
		opts:    opts,
		flows:   map[*rhino.Node]*DataFlowAnalysis{},
	}
}

// Check drives a NodeTraversal over the script and returns every
// diagnostic collected during Visit. Matches Java's
// TypeCheck#process(…) entry-point shape at a narrow scope.
//
// Before walking, we run ControlFlowAnalysis + DataFlowAnalysis on
// the top-level SCRIPT to produce per-statement FlowScopes. Nested
// FUNCTION bodies get their own flow built on demand.
//
// After the main visit walk we invoke the validate-override and
// duplicate-declaration helpers. In Java these are private
// checkPropertyVisibility / validateMethodOverride /
// checkDuplicateVarDeclaration methods on TypeCheck — we keep them
// in separate Go files (superclasscheck.go / checkdupvardecl.go) for
// file-level readability but call them from TypeCheck so the
// pipeline stage count matches DefaultPassConfig (one TypeCheck pass).
func (c *TypeCheck) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.installFlow(script, c.global)
	c.collectPropertyWrites(script)
	c.fileSuppressions = collectFileSuppressions(script)
	trav := NewNodeTraversal(c, c.creator, c.global)
	trav.Traverse(script)
	sc := NewSuperclassCheck(c.path, c.creator, c.global)
	sc.SuppressedNames = c.SuppressedNames
	c.diags = append(c.diags, sc.Check(script)...)
	c.diags = append(c.diags, NewCheckDupVarDecl(c.path, c.global).Check(script)...)
	return c.diags
}

// collectFileSuppressions walks the script's JSDoc blocks for any
// @fileoverview / stand-alone @suppress declarations and returns
// the union of their suppression names. Mirrors Java's
// CompilerPass.skipForSuppressedPath filter — diagnostics whose
// code matches a file-level suppression are dropped silently.
func collectFileSuppressions(script *rhino.Node) map[string]bool {
	out := map[string]bool{}
	var walk func(n *rhino.Node)
	walk = func(n *rhino.Node) {
		if n == nil {
			return
		}
		if info := n.JSDocInfo(); info != nil {
			for _, name := range []string{"missingProperties", "checkTypes", "duplicate", "undefinedVars"} {
				if info.HasSuppression(name) {
					out[name] = true
				}
			}
		}
		for _, c := range n.Children() {
			walk(c)
		}
	}
	walk(script)
	return out
}

// collectPropertyWrites walks the script once before TypeCheck
// runs, registering every `X.prop = …` where X resolves to a
// nominal ObjectType. Java's TypedScopeCreator does this during
// its own traversal; jsgo runs it here so every later
// visitGetProp read has the full property surface available.
// Only extends open (non-struct / non-structural) types — sealed
// receivers keep their declared surface so the INEXISTENT /
// ILLEGAL checks stay accurate.
func (c *TypeCheck) collectPropertyWrites(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.ASSIGN {
		kids := n.Children()
		if len(kids) >= 2 && kids[0].Token() == rhino.GETPROP {
			propKids := kids[0].Children()
			if len(propKids) >= 2 {
				recv := propKids[0]
				propName := propKids[1].String_()
				if recv.Token() == rhino.THIS {
					// Subclass constructor body: `this.X = V` inside
					// a `@constructor @extends {Foo}` function adds
					// X to Foo's surface too. Java's type-inference
					// propagates the write up the extends chain for
					// read-site validation.
					if enc := findEnclosingFunction(n); enc != nil {
						if info := enc.JSDocInfo(); info != nil && info.IsConstructor() {
							if base := info.BaseType(); base != nil {
								if super := c.global.LookupNominalType(base.TypeName); super != nil {
									if !super.HasProperty(propName) {
										super.DefineProperty(propName, jstype.Unknown)
									}
								}
							}
						}
					}
				} else {
					infer := NewTypeInference(c.global)
					recvType := stripNullUndef(infer.InferType(recv))
					defineOnOpenObject(recvType, propName)
				}
			}
		}
	}
	for _, ch := range n.Children() {
		c.collectPropertyWrites(ch)
	}
}

// installFlow builds the CFG + DataFlowAnalysis for `root` (a
// FUNCTION or SCRIPT) and stores it on the TypeCheck instance so
// Visit can look up IN FlowScopes for each statement.
func (c *TypeCheck) installFlow(root *rhino.Node, scope *TypedScope) *DataFlowAnalysis {
	if existing, ok := c.flows[root]; ok {
		return existing
	}
	cfg := BuildControlFlowGraph(root)
	pass := NewTypeInferencePass(scope)
	d := NewDataFlowAnalysis(cfg, EmptyFlowScope(), pass.Transfer)
	d.Run()
	c.flows[root] = d
	return d
}

// enclosingStmt returns the nearest ancestor of `n` that's a CFG
// statement (matches the Token set visited by BuildControlFlowGraph).
// Used so an expression node can find its owning statement's IN
// FlowScope.
func enclosingStmt(n *rhino.Node) *rhino.Node {
	for cur := n; cur != nil; cur = cur.Parent() {
		switch cur.Token() {
		case rhino.VAR, rhino.LET, rhino.CONST, rhino.EXPR_RESULT,
			rhino.IF, rhino.RETURN:
			return cur
		}
	}
	return nil
}

// ShouldTraverse is the pre-order hook. For BLOCK / IF / HOOK nodes
// we intercept the default depth-first descent so we can push
// branch-specific (and post-guard) refinements onto the active scope
// before walking. The rest of the AST is walked normally and emits
// diagnostics in Visit.
func (c *TypeCheck) ShouldTraverse(t *NodeTraversal, n, parent *rhino.Node) bool {
	if n == nil {
		return true
	}
	// Nested function bodies need their own CFG + FlowScope built
	// BEFORE descent so Visit() callbacks inside the body can
	// consult the pre-computed IN scopes. Must live here (pre-order)
	// rather than in Visit (post-order) — by the time Visit fires
	// for the FUNCTION, the children have already been walked.
	if n.Token() == rhino.FUNCTION {
		if local, ok := c.creator.Functions[n]; ok && local != nil {
			c.installFlow(n, local)
		}
	}
	switch n.Token() {
	case rhino.BLOCK, rhino.SCRIPT:
		c.traverseBlock(t, n)
		return false
	case rhino.IF:
		c.traverseIf(t, n)
		return false
	case rhino.HOOK:
		c.traverseHook(t, n)
		return false
	case rhino.AND, rhino.OR, rhino.COALESCE:
		c.traverseShortCircuit(t, n)
		return false
	case rhino.WHILE, rhino.DO_WHILE:
		c.traverseLoop(t, n)
		return false
	}
	return true
}

// traverseLoop walks a WHILE / DO_WHILE loop. The loop body runs
// only when the condition is truthy, so we push the truthy
// refinement before descending into the body. Matches how Java's
// TypeInference refines a FlowScope at the loop header's true
// back-edge.
func (c *TypeCheck) traverseLoop(t *NodeTraversal, n *rhino.Node) {
	kids := n.Children()
	if len(kids) < 2 {
		for _, k := range kids {
			t.TraverseInto(k, n)
		}
		return
	}
	cond, body := kids[0], kids[1]
	t.TraverseInto(cond, n)
	scope := t.CurrentScope()
	if refs := buildRefinements(scope, cond, true); len(refs) > 0 {
		scope.PushRefinement(refs)
		t.TraverseInto(body, n)
		scope.PopRefinement()
	} else {
		t.TraverseInto(body, n)
	}
}

// traverseShortCircuit handles the AND / OR / COALESCE binary
// operators. The RHS is evaluated only when the LHS is truthy
// (for &&) or falsy (for || / ??), so we walk the RHS under a
// refinement derived from the LHS. Matches Java's
// SemanticReverseAbstractInterpreter behaviour at the expression
// level.
func (c *TypeCheck) traverseShortCircuit(t *NodeTraversal, n *rhino.Node) {
	kids := n.Children()
	if len(kids) != 2 {
		for _, k := range kids {
			t.TraverseInto(k, n)
		}
		return
	}
	left, right := kids[0], kids[1]
	t.TraverseInto(left, n)
	scope := t.CurrentScope()
	whenTruthy := n.Token() == rhino.AND
	if refs := buildRefinements(scope, left, whenTruthy); len(refs) > 0 {
		scope.PushRefinement(refs)
		t.TraverseInto(right, n)
		scope.PopRefinement()
	} else {
		t.TraverseInto(right, n)
	}
}

// traverseBlock walks a block's statements in order, threading
// post-guard refinements forward — when a statement is an `if (cond)
// { return; }` (or `throw`), the negated cond's refinement is pushed
// for the remainder of the block. Mirrors how Java's TypeInference
// applies a refined FlowScope to siblings after an early-exit guard.
func (c *TypeCheck) traverseBlock(t *NodeTraversal, block *rhino.Node) {
	scope := t.CurrentScope()
	pushed := 0
	for _, child := range block.Children() {
		t.TraverseInto(child, block)
		if guard, ok := earlyExitGuard(child); ok {
			if refs := buildRefinements(scope, guard, false); len(refs) > 0 {
				scope.PushRefinement(refs)
				pushed++
			}
		}
		// Post-while narrowing — the loop exits when `cond` is
		// false, so after `while (cond) { … }` the negated
		// refinement holds for the remainder of the block. Only
		// fire when the loop body has no `break` that could
		// skip the natural exit condition (simple heuristic —
		// we accept any body and let Java's conservative rule
		// catch corner cases).
		if child.Token() == rhino.WHILE {
			if kids := child.Children(); len(kids) >= 1 {
				if refs := buildRefinements(scope, kids[0], false); len(refs) > 0 {
					scope.PushRefinement(refs)
					pushed++
				}
			}
		}
	}
	for i := 0; i < pushed; i++ {
		scope.PopRefinement()
	}
	// Block / Script don't emit diagnostics in Visit themselves, but
	// the contract is to call Visit so callbacks downstream that care
	// see the node. Calling here is a no-op for TypeCheck.
	c.Visit(t, block, block.Parent())
}

// earlyExitGuard reports whether `n` is `if (cond) { return | throw
// | break | continue }` with no else, returning the cond expression
// for use in negated refinement of the remaining siblings.
func earlyExitGuard(n *rhino.Node) (*rhino.Node, bool) {
	if n == nil || n.Token() != rhino.IF {
		return nil, false
	}
	kids := n.Children()
	if len(kids) != 2 {
		return nil, false // has an else branch — no guarantee siblings refine
	}
	if !blockExits(kids[1]) {
		return nil, false
	}
	return kids[0], true
}

// blockExits reports whether `n` (a single statement or a BLOCK)
// definitely terminates execution of the enclosing function via
// return / throw on every path.
func blockExits(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	if definitelyReturns(n) {
		return true
	}
	// THROW — our parser doesn't currently surface a THROW token, but
	// finishLoopBody swallows `throw expr;`. The presence of "throw"
	// at AST time is rare; leaving the conservative no-throw check.
	return false
}

// traverseIf walks `if (cond) thenBlock [else elseBlock]` with a
// refinement overlay active on the current scope while descending
// into each branch. Mirrors how Java's TypeInferencePass refines a
// FlowScope at branch points.
func (c *TypeCheck) traverseIf(t *NodeTraversal, n *rhino.Node) {
	kids := n.Children()
	if len(kids) < 2 {
		for _, k := range kids {
			t.TraverseInto(k, n)
		}
		return
	}
	cond := kids[0]
	thenBlock := kids[1]
	var elseBlock *rhino.Node
	if len(kids) >= 3 {
		elseBlock = kids[2]
	}
	t.TraverseInto(cond, n)
	scope := t.CurrentScope()
	if refs := buildRefinements(scope, cond, true); len(refs) > 0 {
		scope.PushRefinement(refs)
		t.TraverseInto(thenBlock, n)
		scope.PopRefinement()
	} else {
		t.TraverseInto(thenBlock, n)
	}
	if elseBlock != nil {
		if refs := buildRefinements(scope, cond, false); len(refs) > 0 {
			scope.PushRefinement(refs)
			t.TraverseInto(elseBlock, n)
			scope.PopRefinement()
		} else {
			t.TraverseInto(elseBlock, n)
		}
	}
}

// traverseHook is the HOOK analogue of traverseIf — refines per
// branch so `cond ? cond.foo() : fallback` resolves cond.foo()
// against the truthy form of cond.
func (c *TypeCheck) traverseHook(t *NodeTraversal, n *rhino.Node) {
	kids := n.Children()
	if len(kids) < 3 {
		for _, k := range kids {
			t.TraverseInto(k, n)
		}
		return
	}
	cond, thenE, elseE := kids[0], kids[1], kids[2]
	t.TraverseInto(cond, n)
	scope := t.CurrentScope()
	if refs := buildRefinements(scope, cond, true); len(refs) > 0 {
		scope.PushRefinement(refs)
		t.TraverseInto(thenE, n)
		scope.PopRefinement()
	} else {
		t.TraverseInto(thenE, n)
	}
	if refs := buildRefinements(scope, cond, false); len(refs) > 0 {
		scope.PushRefinement(refs)
		t.TraverseInto(elseE, n)
		scope.PopRefinement()
	} else {
		t.TraverseInto(elseE, n)
	}
}

// Visit is the post-order dispatch. Matches Java's visit(t,n,parent)
// override on AbstractPostOrderCallback; each branch mirrors a
// TypeCheck.java visitXxx helper. Before per-node dispatch, we push
// the IN FlowScope of the enclosing statement onto the scope so
// expression-level inference sees the flow-narrowed types.
func (c *TypeCheck) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	scope := t.CurrentScope()
	popFlow := c.pushEnclosingFlow(n, scope)
	defer popFlow()
	switch n.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		c.emit(c.visitVar(n, scope))
	case rhino.ASSIGN:
		c.emit(c.visitAssign(n, scope))
	case rhino.ASSIGN_AND, rhino.ASSIGN_OR, rhino.ASSIGN_COALESCE:
		// Logical compound assigns (`&&=` / `||=` / `??=`) enforce the
		// same LHS-declared-vs-RHS-value contract as bare `=`. Java's
		// TypeValidator routes them through the same
		// expectCanAssignToWithMessage path — we reuse visitAssign so
		// both the NAME and GETPROP arms fire.
		c.emit(c.visitAssign(n, scope))
	case rhino.CALL:
		c.emit(c.visitCall(n, scope))
	case rhino.NEW:
		c.emit(c.visitNew(n, scope))
	case rhino.GETPROP:
		c.emit(c.visitGetProp(n, scope))
	case rhino.GETELEM:
		c.emit(c.visitGetElem(n, scope))
	case rhino.CAST:
		c.emit(c.visitCast(n, scope))
	case rhino.CLASS_FIELD_DEF:
		c.emit(c.visitClassField(n, parent, scope))
	case rhino.RETURN:
		c.emit(c.visitReturn(n, scope))
	case rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE:
		// JSC_DETERMINISTIC_TEST is emitted inline from Java's
		// TypeCheck.visitBinaryOp; the helper lives in
		// checkdeterministictest.go purely as file-organization.
		c.emit(VisitDeterministicTest(c.path, scope, n))
	case rhino.IN:
		c.emit(c.visitIn(n, scope))
	case rhino.FOR_IN:
		c.emit(c.visitForIn(n, scope))
	case rhino.LSH, rhino.RSH, rhino.URSH, rhino.BITAND, rhino.BITOR, rhino.BITXOR:
		c.emit(c.visitBitOp(n, scope))
	case rhino.LT, rhino.LE, rhino.GT, rhino.GE:
		c.emit(c.visitRelationalDefault(n, scope))
	case rhino.FUNCTION:
		c.emit(c.visitFunctionParams(n, scope))
	case rhino.ARRAYLIT:
		c.emit(c.visitArrayLit(n, scope))
	case rhino.OBJECTLIT:
		c.emit(c.visitObjectLit(n, scope))
	}
	// Strict-mode operand checks — OFF by default, matching Java.
	// These fire only when `--jscomp_warning=strictCheckTypes` is set.
	// Compound assigns (`+=`, `-=`, …) route to the same validator as
	// their bare counterparts; Java's TypeValidator uses one predicate
	// per operator family regardless of the assignment variant.
	if c.opts.IsOn(GroupStrictCheckTypes) {
		switch n.Token() {
		case rhino.ADD, rhino.ASSIGN_ADD:
			c.emit(c.visitAdd(n, scope))
		case rhino.SUB, rhino.MUL, rhino.DIV, rhino.MOD,
			rhino.ASSIGN_SUB, rhino.ASSIGN_MUL, rhino.ASSIGN_DIV, rhino.ASSIGN_MOD:
			c.emit(c.visitArithmetic(n, scope))
		case rhino.BITAND, rhino.BITOR, rhino.BITXOR,
			rhino.LSH, rhino.RSH, rhino.URSH,
			rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR,
			rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH:
			c.emit(c.visitBitwise(n, scope))
		case rhino.LT, rhino.LE, rhino.GT, rhino.GE:
			c.emit(c.visitRelativeComparison(n, scope))
		}
	}
}

// pushEnclosingFlow looks up the CFG / FlowScope active for `n`'s
// enclosing statement and pushes it as a refinement on the current
// scope. Returns a pop function the caller defers. Safe no-op if
// the node isn't inside a tracked CFG.
func (c *TypeCheck) pushEnclosingFlow(n *rhino.Node, scope *TypedScope) func() {
	stmt := enclosingStmt(n)
	if stmt == nil {
		return func() {}
	}
	// Look up the closest FUNCTION / SCRIPT ancestor so we know
	// which flow analysis owns this statement.
	fn := enclosingFlowRoot(n)
	d := c.flows[fn]
	if d == nil {
		return func() {}
	}
	cfgNode := d.Cfg.Get(stmt)
	if cfgNode == nil {
		return func() {}
	}
	in := d.InOfAst(cfgNode)
	if in == nil || len(in.types) == 0 {
		return func() {}
	}
	scope.PushRefinement(in.types)
	return scope.PopRefinement
}

// enclosingFlowRoot returns the nearest FUNCTION or SCRIPT ancestor
// of `n` — the CFG owner for that subtree.
func enclosingFlowRoot(n *rhino.Node) *rhino.Node {
	for cur := n; cur != nil; cur = cur.Parent() {
		switch cur.Token() {
		case rhino.FUNCTION, rhino.SCRIPT:
			return cur
		}
	}
	return nil
}

func (c *TypeCheck) emit(d []diag.Diagnostic) {
	c.diags = append(c.diags, d...)
}

// ---- per-node visitors ------------------------------------------ //

// visitVar — port of TypeCheck.java's handling of VAR / LET / CONST:
// compare the initializer's inferred type against the @type
// declaration. Reports at the initializer's position to match Java's
// "initializing variable" variant.
func (c *TypeCheck) visitVar(v *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	var out []diag.Diagnostic
	// CheckObjectKeys.java (folded in): scan @type expressions for
	// Object<K,V> entries with non-stringifiable keys before the
	// regular init-subtype check. Applies even when the decl has no
	// initializer.
	out = append(out, c.visitObjectKeyTypes(v, scope)...)
	infer := NewTypeInference(scope)
	// Walk every NAME child so multi-binding forms
	// `var a = 1, /** @type {T} */ b = 2` each get their own
	// initializer subtype check. Java's TypeCheck iterates via
	// Node.getFirstChild + getNext which we model through the slice.
	for _, name := range v.Children() {
		if name.Token() != rhino.NAME {
			continue
		}
		init := name.FirstChild()
		if init == nil {
			continue
		}
		// `/** @enum */ var E = {FOO: 1, ...}` — the VAR's declared
		// type is the EnumType, but the init is an OBJECTLIT that
		// infers to a structural {...}. Skip the normal subtype
		// check and run the enum-specific diagnostics instead.
		if info := v.JSDocInfo(); info.IsEnum() {
			out = append(out, c.visitEnumInit(name, init, scope)...)
			continue
		}
		declared := scope.GetDeclaredType(name.String_())
		if declared == nil {
			continue
		}
		found := infer.InferType(init)
		if jstype.IsSubtype(found, declared) {
			continue
		}
		msg := fmt.Sprintf("initializing variable\nfound   : %s\nrequired: %s",
			found.TypeName(), declared.TypeName())
		line, col := diagPosForArg(init)
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: line, Col: col,
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
		})
	}
	return out
}

// visitEnumInit validates a `/** @enum [{T}] */ var E = init`
// declaration. Three cases mirror Java's EnumType initialisation:
//
//  1. init is an OBJECTLIT — per-entry subtype check against the
//     declared Underlying; mismatch emits JSC_TYPE_MISMATCH at
//     the entry value with "assignment to property <k> of enum{E}".
//  2. init is a NAME referring to another enum — alias; already
//     folded in scope.go. No diagnostic here.
//  3. init is anything else — emit "initializing variable" with
//     the enum's EnumType as required. Matches how Java's
//     TypeCheck handles a non-enum RHS under an @enum decl.
func (c *TypeCheck) visitEnumInit(name, init *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	enum, ok := scope.GetDeclaredType(name.String_()).(*jstype.EnumType)
	if !ok {
		return nil
	}
	// Recover the DECLARED underlying from the @enum JSDoc (defaults
	// to Number when `@enum` has no brace). Aliasing may have
	// replaced the slot's pointer with the source enum, whose own
	// Underlying can differ — we need the original-site declaration
	// to diagnose mismatched-alias correctly.
	declaredUnderlying := jstype.JSType(jstype.Number)
	if info := name.Parent().JSDocInfo(); info != nil {
		if expr := info.EnumParameterType(); expr != nil {
			declaredUnderlying = scope.ResolveExpr(expr.TypeName)
		}
	}
	_ = declaredUnderlying // referenced below when the alias branch fires
	infer := NewTypeInference(scope)
	if init.Token() == rhino.OBJECTLIT {
		var out []diag.Diagnostic
		for _, kv := range init.Children() {
			if kv.Token() != rhino.STRING_KEY {
				continue
			}
			value := kv.FirstChild()
			if value == nil {
				continue
			}
			found := infer.InferType(value)
			if jstype.IsSubtype(found, enum.Underlying) {
				continue
			}
			// Report at the STRING_KEY's column — Java's enum
			// validator points at the property name, not the
			// initializer's RHS value.
			out = append(out, diag.Diagnostic{
				Path: c.path, Line: kv.Lineno(), Col: kv.Charno(),
				Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
				Message: fmt.Sprintf(
					"assignment to property %s of enum{%s}\nfound   : %s\nrequired: %s",
					kv.String_(), enum.Name, found.TypeName(), enum.Underlying.TypeName()),
			})
		}
		return out
	}
	if init.Token() == rhino.NAME {
		if other, ok := scope.GetDeclaredType(init.String_()).(*jstype.EnumType); ok {
			// Aliased enum — Java emits JSC_TYPE_MISMATCH
			// "incompatible enum element types" when the alias's
			// explicit `@enum {T}` underlying disagrees with the
			// source enum's. Use the declared underlying captured
			// above (the scope-creator has already redirected
			// scope[name] to the source enum, so enum.Underlying
			// reflects the source, not the alias declaration).
			if !jstype.IsSubtype(other.Underlying, declaredUnderlying) {
				return []diag.Diagnostic{{
					Path: c.path, Line: init.Lineno(), Col: init.Charno(),
					Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
					Message: fmt.Sprintf(
						"incompatible enum element types\nfound   : %s\nrequired: %s",
						other.Underlying.TypeName(), declaredUnderlying.TypeName()),
				}}
			}
			_ = enum
			return nil
		}
	}
	found := infer.InferType(init)
	return []diag.Diagnostic{{
		Path: c.path, Line: init.Lineno(), Col: init.Charno(),
		Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
		Message: fmt.Sprintf("initializing variable\nfound   : %s\nrequired: enum{%s}",
			found.TypeName(), enum.Name),
	}}
}

// visitAssign handles bare `IDENT = …` and `obj.prop = …` assignment
// forms. Reports at the LHS position (matching Java's "assignment"
// and "assignment to property X of Y" variants).
func (c *TypeCheck) visitAssign(a *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := a.Children()
	if len(kids) < 2 {
		return nil
	}
	lhs, rhs := kids[0], kids[1]
	infer := NewTypeInference(scope)
	switch lhs.Token() {
	case rhino.NAME:
		// Assignment-target lookup must use the raw declared type,
		// not the flow-refined one — otherwise a prior
		// `var x = null;` narrows x's declared type to null and
		// every subsequent assignment false-positives. The refined
		// type is still correct for reads (RHS inference).
		declared := scope.GetRawDeclaredType(lhs.String_())
		if declared == nil {
			return nil
		}
		found := infer.InferType(rhs)
		if jstype.IsSubtype(found, declared) {
			return nil
		}
		msg := fmt.Sprintf("assignment\nfound   : %s\nrequired: %s",
			found.TypeName(), declared.TypeName())
		return []diag.Diagnostic{{
			Path: c.path, Line: lhs.Lineno(), Col: lhs.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
		}}
	case rhino.GETPROP:
		propKids := lhs.Children()
		if len(propKids) < 2 {
			return nil
		}
		propName := propKids[1].String_()
		recv := propKids[0]
		// ES6 classes reject `ClassName.prototype = …` reassignments.
		// Java fires at the outer ASSIGN's (or wrapping EXPR_RESULT's)
		// column, so we report at the statement's position.
		if propName == "prototype" && recv.Token() == rhino.NAME {
			if fn, ok := scope.GetDeclaredType(recv.String_()).(*jstype.FunctionType); ok && fn.IsES6Class {
				line, col := lhs.Lineno(), lhs.Charno()
				if p := a.Parent(); p != nil && p.Token() == rhino.EXPR_RESULT {
					line, col = p.Lineno(), p.Charno()
				}
				return []diag.Diagnostic{{
					Path: c.path, Line: line, Col: col,
					Code: diag.JSC_REASSIGN_CLASS_PROTOTYPE, Level: diag.Error,
					Message: "Reassigning a class prototype is not allowed",
				}}
			}
		}
		var declared jstype.JSType
		var ownerName string
		if recv.Token() == rhino.NAME {
			if fn, ok := scope.GetDeclaredType(recv.String_()).(*jstype.FunctionType); ok && fn.Instance != nil {
				if t := fn.Instance.GetStaticPropertyType(propName); t != nil {
					declared = t
					ownerName = recv.String_()
				}
			}
		}
		if declared == nil && recv.Token() == rhino.NAME {
			// `/** @type {T} */ goog.foo = …` namespace-property
			// declaration made via recordNamespaceAssign.
			if t := scope.LookupNamespaceProp(recv.String_(), propName); t != nil {
				declared = t
				ownerName = recv.String_()
			}
		}
		if declared == nil {
			obj, ok := infer.InferType(recv).(*jstype.ObjectType)
			if !ok {
				return nil
			}
			declared = obj.GetPropertyType(propName)
			if declared == nil {
				return nil
			}
			if obj.IsPropertyInferred(propName) {
				// Slot was inferred from an RHS literal rather than
				// declared via @type — Java allows open reassignment
				// for inference-only property surfaces.
				return nil
			}
			ownerName = obj.TypeName()
		}
		found := infer.InferType(rhs)
		if jstype.IsSubtype(found, declared) {
			return nil
		}
		msg := fmt.Sprintf("assignment to property %s of %s\nfound   : %s\nrequired: %s",
			propName, ownerName, found.TypeName(), declared.TypeName())
		// Report at the enclosing EXPR_RESULT's column — Java
		// points at the statement's start, so parenthesized
		// receivers line up (`(new F).x_ = 3` → col 0 not the
		// `.x_` column). Falls back to the LHS when the ASSIGN
		// is nested inside a larger expression.
		line, col := lhs.Lineno(), lhs.Charno()
		if p := a.Parent(); p != nil && p.Token() == rhino.EXPR_RESULT {
			line, col = p.Lineno(), p.Charno()
		}
		return []diag.Diagnostic{{
			Path: c.path, Line: line, Col: col,
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
		}}
	}
	return nil
}

// visitNew — port of TypeCheck.java's visitNew. Reports
// JSC_NOT_A_CONSTRUCTOR when the callee's type isn't a constructor
// function; otherwise delegates to visitCall for the same per-arg
// check.
func (c *TypeCheck) visitNew(node *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := node.Children()
	if len(kids) < 1 {
		return nil
	}
	callee := kids[0]
	infer := NewTypeInference(scope)
	t := infer.InferType(callee)
	// Only fire on a definitely non-constructor callee. Unknown /
	// union-with-constructor / All stay quiet to avoid
	// false-positives on incomplete JSDoc.
	if notCtor, found := isDefinitelyNonConstructor(t); notCtor {
		return []diag.Diagnostic{{
			Path: c.path, Line: node.Lineno(), Col: node.Charno(),
			Code: diag.JSC_NOT_A_CONSTRUCTOR, Level: diag.Error,
			Message: "cannot instantiate non-constructor, found type: " + found.TypeName(),
		}}
	}
	return c.visitCall(node, scope)
}

// isSafeCalleeForArity reports whether a CALL's callee shape lets
// us confidently enforce JSC_WRONG_ARGUMENT_COUNT. Java fires on
// every call-site; jsgo narrows to patterns where our FunctionType
// carries accurate arity info — bare identifiers, namespace
// properties (via LookupNamespaceProp), and static-method GETPROPs
// on a known constructor. Prototype-method calls stay silent
// because override-aware arity isn't ported yet.
func isSafeCalleeForArity(callee *rhino.Node, scope *TypedScope) bool {
	if callee == nil {
		return false
	}
	if callee.Token() == rhino.NAME {
		return true
	}
	if callee.Token() == rhino.GETPROP {
		kids := callee.Children()
		if len(kids) < 2 {
			return false
		}
		recv := kids[0]
		prop := kids[1]
		if recv.Token() != rhino.NAME {
			return false
		}
		// Namespace-property case: `ns.foo()` with ns declared as a
		// plain VAR object literal and foo registered via
		// recordNamespaceAssign.
		if scope.LookupNamespaceProp(recv.String_(), prop.String_()) != nil {
			return true
		}
		// Static-method case: `C.foo()` where C is a FunctionType.
		if fn, ok := scope.GetDeclaredType(recv.String_()).(*jstype.FunctionType); ok && fn.Instance != nil {
			if fn.Instance.GetStaticPropertyType(prop.String_()) != nil {
				return true
			}
		}
	}
	return false
}

// isCallableBuiltinCtor reports whether `callee` names one of the
// ECMAScript built-in constructors that double as coercion functions
// (`String(x)`, `Number(x)`, etc.). Java's type registry marks these
// as callable-without-new; firing JSC_CONSTRUCTOR_NOT_CALLABLE on
// them would be a false positive. Scoped to bare NAME references so
// user code that shadows the global still trips the real check.
func isCallableBuiltinCtor(callee *rhino.Node) bool {
	if callee == nil || callee.Token() != rhino.NAME {
		return false
	}
	switch callee.String_() {
	case "String", "Number", "Boolean", "Symbol", "Error",
		"Array", "Object", "RegExp", "Date",
		"TypeError", "RangeError", "SyntaxError", "ReferenceError",
		"EvalError", "URIError":
		return true
	}
	return false
}

// isSuperCall reports whether `callee` is a `super(...)` reference
// inside a derived-class constructor body. The rhino parser emits
// `super` as a NAME node with the literal identifier "super" — Java
// handles super-calls through a dedicated visitSuper branch rather
// than the constructor-callability check, so we mirror that here.
func isSuperCall(callee *rhino.Node) bool {
	return callee != nil && callee.Token() == rhino.NAME && callee.String_() == "super"
}

// isDefinitelyNonConstructor reports whether `t` is guaranteed not
// to be a constructor. Returns the offending type so the diagnostic
// can show what was found. A FunctionType without IsConstructor
// counts (plain functions); a FunctionType marked as @interface /
// @record (Instance.IsStructural) counts too, because Java forbids
// instantiating an interface. Non-function primitives count;
// Unknown / All / unions stay quiet.
func isDefinitelyNonConstructor(t jstype.JSType) (bool, jstype.JSType) {
	if t == nil {
		return false, nil
	}
	if isUnknownOrNil(t) {
		return false, nil
	}
	// Strip `T | null / undefined` on the callee so a nullable
	// callable doesn't skip the check entirely. Mirrors Java's
	// restrictByNotNullOrUndefined before visitNew's isConstructor
	// probe.
	if fn, ok := unwrapCallableUnion(t); ok {
		t = fn
	} else if _, ok := t.(*jstype.UnionType); ok {
		return false, nil
	}
	if fn, ok := t.(*jstype.FunctionType); ok {
		// @interface / @record produces a structural Instance — Java
		// tracks this separately from IsConstructor and uses the
		// structural bit to reject `new`.
		if fn.Instance != nil && fn.Instance.IsStructural {
			return true, fn
		}
		if fn.IsConstructor {
			return false, nil
		}
		return true, fn
	}
	// Primitives are obviously not constructors.
	switch t {
	case jstype.Number, jstype.String, jstype.Boolean, jstype.Null, jstype.Void:
		return true, t
	}
	return false, nil
}

// visitCall handles both CALL and NEW nodes. Reports per-argument
// JSC_TYPE_MISMATCH at the argument's position, matching
// TypeCheck.java's visitCall / visitNew.
func (c *TypeCheck) visitCall(node *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := node.Children()
	if len(kids) < 1 {
		return nil
	}
	callee := kids[0]
	args := kids[1:]
	infer := NewTypeInference(scope)
	calleeType := infer.InferType(callee)
	fn, ok := calleeType.(*jstype.FunctionType)
	if !ok {
		// Unwrap `FunctionType|null` unions so a nullable callable
		// parameter (`/** @param {?function(T)} */`) still matches
		// the arity-check branch. Java's TypeCheck.validateCall
		// strips null from the callee's resolved type before the
		// isFunctionType check.
		if stripped, ok2 := unwrapCallableUnion(calleeType); ok2 {
			fn = stripped
			ok = true
		}
	}
	if !ok {
		// Only fire on bare NAME / literal callees so GETELEM
		// (well-known symbol access, element reads) doesn't
		// false-positive. Java's visitCall also considers the
		// callee's declared type rather than the element-read
		// fallback result.
		if node.Token() == rhino.CALL && isDefinitelyNotCallable(calleeType) && callee.Token() != rhino.GETELEM {
			return []diag.Diagnostic{{
				Path: c.path, Line: node.Lineno(), Col: node.Charno(),
				Code: diag.JSC_NOT_FUNCTION_TYPE, Level: diag.Error,
				Message: fmt.Sprintf("%s expressions are not callable", calleeType.TypeName()),
			}}
		}
		return nil
	}
	var out []diag.Diagnostic
	calleeName := describeCallee(callee)
	// `@constructor` functions cannot be called without `new` —
	// Java's TypeCheck.visitCall fires JSC_CONSTRUCTOR_NOT_CALLABLE
	// when the callee is constructor-shaped and the call-site is a
	// CALL (not NEW). Interfaces / records are structural and their
	// FunctionType has IsConstructor=false, so this branch is
	// limited to genuine class ctors.
	//
	// Exclusions: the ECMAScript "callable constructors"
	// (`String`, `Number`, `Boolean`, `Symbol`, `Error`, `Array`,
	// `Object`, `RegExp`, `Date`) legally double as coercion
	// functions; Java's registry marks them callable-without-new
	// and jsgo's inference.go intercepts these names at the call
	// level to return the primitive result. `super(...)` inside a
	// derived ctor body also resolves through the parent class's
	// FunctionType but is structurally always a CALL — Java handles
	// it through a separate visitSuper branch.
	if node.Token() == rhino.CALL && fn.IsConstructor && fn.Instance != nil && !fn.Instance.IsStructural &&
		!isCallableBuiltinCtor(callee) && !isSuperCall(callee) {
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: node.Lineno(), Col: node.Charno(),
			Code: diag.JSC_CONSTRUCTOR_NOT_CALLABLE, Level: diag.Error,
			Message: fmt.Sprintf(`Constructor (typeof %s) should be called with the "new" keyword`,
				fn.Instance.TypeName()),
		})
	}
	// Arity check — fire JSC_WRONG_ARGUMENT_COUNT when the call
	// site's arg count falls outside [MinArity, MaxArity] as
	// computed from the function's @param annotations and closure
	// naming conventions (opt_<name>, var_args). Mirrors
	// TypeCheck.visitCall's arity branch in Java.
	//
	// Scoped to NAME callees — GETPROP callees (e.g.
	// `obj.method()`) often resolve through override chains where
	// a subclass's @override-only JSDoc drops the parent's
	// optional-param markers. Firing there produces more FPs
	// than wins until the override-aware inference lands.
	if node.Token() == rhino.CALL && len(fn.Params) > 0 && isSafeCalleeForArity(callee, scope) {
		if got := len(args); got < fn.MinArity || got > fn.MaxArity {
			maxStr := fmt.Sprintf("%d", fn.MaxArity)
			if fn.MaxArity >= jstype.MaxArgsUnlimited {
				maxStr = "unlimited"
			}
			out = append(out, diag.Diagnostic{
				Path: c.path, Line: node.Lineno(), Col: node.Charno(),
				Code: diag.JSC_WRONG_ARGUMENT_COUNT, Level: diag.Error,
				Message: fmt.Sprintf(
					"Function %s: called with %d argument(s). Function requires at least %d argument(s) and no more than %s argument(s).",
					calleeName, got, fn.MinArity, maxStr),
			})
		}
	}
	// Arity-mismatch: still type-check the overlapping prefix so a
	// `f(0)` against `function f(s: string)` reports the type error
	// even when other params are dropped. Java fires both
	// JSC_WRONG_ARGUMENT_COUNT and the per-arg JSC_TYPE_MISMATCH.
	n := len(args)
	if len(fn.Params) < n {
		n = len(fn.Params)
	}
	for i := 0; i < n; i++ {
		arg := args[i]
		param := fn.Params[i]
		found := infer.InferType(arg)
		// Optional params (`T=`) admit `undefined` as a legal argument.
		// In Java the FunctionType stores the pipe-widened
		// `T|undefined`; our resolver keeps `T` and flags the slot
		// via MinArity (vs MaxArity for the optional range). Variadic
		// params (`...T`) also live past MinArity but must still
		// reject undefined — those callers should send at least one
		// concrete element — so we only widen when the slot is
		// non-variadic (MaxArity finite).
		if i >= fn.MinArity && fn.MaxArity < jstype.MaxArgsUnlimited && found == jstype.Void {
			continue
		}
		if jstype.IsSubtype(found, param) {
			continue
		}
		line, col := diagPosForArg(arg)
		msg := fmt.Sprintf("actual parameter %d of %s does not match formal parameter\nfound   : %s\nrequired: %s",
			i+1, calleeName, found.TypeName(), param.TypeName())
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: line, Col: col,
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
		})
	}
	return out
}

// visitGetProp — port of TypeCheck.java's visitGetProp. Fires:
//
//  1. JSC_ILLEGAL_PROPERTY_ACCESS when the receiver is a `@dict`
//     (dot-access on a dict is forbidden),
//  2. JSC_INEXISTENT_PROPERTY on NAMED structural receivers whose
//     member map doesn't declare propName.
//
// Anonymous inline records are widened through unions in Java and
// firing on them false-positives.
func (c *TypeCheck) visitGetProp(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	propNode := kids[1]
	propName := propNode.String_()
	infer := NewTypeInference(scope)
	rawRecv := infer.InferType(kids[0])
	// `null.X` / `undefined.X` are runtime TypeErrors — Java's
	// TypeCheck fires "No properties on this expression" at the
	// property name. Report before stripNullUndef so a union
	// with a non-null member (the ordinary nullable receiver
	// case) still flows through the INEXISTENT_PROPERTY rules
	// below.
	if rawRecv == jstype.Null || rawRecv == jstype.Void {
		return []diag.Diagnostic{{
			Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
			Message: fmt.Sprintf(
				"No properties on this expression\nfound   : %s\nrequired: Object",
				rawRecv.TypeName()),
		}}
	}
	// Closure's legacy-mode default makes `Foo` nullable — so a
	// receiver of type `Foo|null` should still be property-checked
	// against Foo. Strip the null/void members of a union so the
	// INEXISTENT_PROPERTY / ILLEGAL_PROPERTY_ACCESS checks still fire.
	recvType := stripNullUndef(rawRecv)
	// Enum member access: `E.FOO` with FOO absent from E's member
	// set fires JSC_INEXISTENT_ENUM_ELEMENT — unless the name is a
	// static property attached via `E.method = function(){…}`
	// (tracked on the scope's namespace-prop map).
	if enum, ok := recvType.(*jstype.EnumType); ok {
		if enum.Members[propName] {
			return nil
		}
		if kids[0].Token() == rhino.NAME {
			if scope.LookupNamespaceProp(kids[0].String_(), propName) != nil {
				return nil
			}
		}
		return []diag.Diagnostic{{
			Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
			Code: diag.JSC_INEXISTENT_ENUM_ELEMENT, Level: diag.Error,
			Message: fmt.Sprintf("element %s does not exist on this enum", propName),
		}}
	}
	if isDictReceiver(recvType) {
		return []diag.Diagnostic{{
			Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
			Code: diag.JSC_ILLEGAL_PROPERTY_ACCESS, Level: diag.Error,
			Message: "Cannot do '.' access on a dict",
		}}
	}
	// Union-typed receiver: Java's TypeCheck.visitGetProp emits
	// JSC_POSSIBLE_INEXISTENT_PROPERTY when at least one variant of
	// the union declares the property and another lacks it (ignoring
	// null / undefined variants, which already flag as TYPE_MISMATCH
	// at the dereference site). The null/void members were stripped
	// above; any remaining variant-without-the-property is a genuine
	// partial-surface union.
	if u, ok := recvType.(*jstype.UnionType); ok {
		return c.possibleInexistentPropertyDiag(u, propName, propNode)
	}
	obj := dereferencedObject(recvType)
	if obj == nil {
		return nil
	}
	if obj.Name == "" {
		return nil
	}
	if obj.HasProperty(propName) {
		return nil
	}
	if c.fileSuppressions["missingProperties"] {
		return nil
	}
	// Writing a new property to a @struct instance outside its own
	// constructor is JSC_ILLEGAL_PROPERTY_CREATION — Java flags the
	// sealed-surface violation distinctly from plain reads.
	if obj.IsStruct && isAssignLHS(n) && !insideOwnConstructor(n, obj) {
		return []diag.Diagnostic{{
			Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
			Code: diag.JSC_ILLEGAL_PROPERTY_CREATION, Level: diag.Error,
			Message: "Cannot add a property to a struct instance after it is constructed. (If you already declared the property, make sure to give it a type.)",
		}}
	}
	// Dedup across the script.
	key := obj.TypeName() + "." + propName
	if c.inexistentReported[key] {
		return nil
	}
	// Broader INEXISTENT_PROPERTY surface — port of TypeCheck.java's
	// non-structural branch. Java fires on any nominal ObjectType
	// whose prototype chain definitively lacks the property; we keep
	// the structural / @struct / known-surface cases, but also fall
	// through for plain @constructor nominal types provided the
	// receiver isn't an instance write that Java would treat as an
	// implicit surface extension.
	//
	// False-positive guard: builtin stubs (e.g. the minimal DOM
	// `Window` / `HTMLElement` externs we ship) have empty member
	// maps — firing on every access there produces 40+ regressions
	// with zero signal, so we defer to HasEmptyMemberMap() for those
	// chains. Real user code with a declared surface stays covered.
	if !obj.IsStructural && !obj.IsStruct && !isKnownSurfaceBuiltin(obj.Name) {
		if obj.HasEmptyMemberMap() {
			return nil
		}
		// Instance writes outside the constructor implicitly extend
		// a non-struct class's surface in Java's default mode; only
		// reads trip the diagnostic.
		if isAssignLHS(n) {
			return nil
		}
	}
	msg := fmt.Sprintf("Property %s never defined on %s", propName, obj.TypeName())
	if c.inexistentReported == nil {
		c.inexistentReported = map[string]bool{}
	}
	c.inexistentReported[key] = true
	return []diag.Diagnostic{{
		Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
		Code: diag.JSC_INEXISTENT_PROPERTY, Level: diag.Error, Message: msg,
	}}
}

// possibleInexistentPropertyDiag handles the union-typed receiver
// branch of TypeCheck.visitGetProp. Java emits
// JSC_POSSIBLE_INEXISTENT_PROPERTY when at least one union variant
// lacks the property (and that variant isn't null / undefined). When
// every variant has the property, stay silent. When every variant
// lacks it, degrade to plain JSC_INEXISTENT_PROPERTY keyed on the
// first nominal offender — mirroring Java's best-effort single-arm
// fallback.
func (c *TypeCheck) possibleInexistentPropertyDiag(
	u *jstype.UnionType, propName string, propNode *rhino.Node,
) []diag.Diagnostic {
	if c.fileSuppressions["missingProperties"] {
		return nil
	}
	haveWith := 0
	haveWithout := 0
	var sampleWithout *jstype.ObjectType
	for _, m := range u.Members() {
		// null / undefined members were stripped upstream but guard
		// anyway so future callers can pass raw unions safely.
		if m == jstype.Null || m == jstype.Void {
			continue
		}
		// Primitive / unknown / stub-chain (HasEmptyMemberMap) variants
		// are treated as "has it" — anything non-nominal can't be
		// definitively said to lack the property, and firing on bare
		// stubs is the false-positive trap documented above.
		obj := dereferencedObject(m)
		if obj == nil || obj.HasProperty(propName) || obj.HasEmptyMemberMap() {
			haveWith++
			continue
		}
		haveWithout++
		if sampleWithout == nil {
			sampleWithout = obj
		}
	}
	if haveWithout == 0 {
		return nil
	}
	if haveWith == 0 {
		// haveWithout > 0 guarantees sampleWithout is non-nil.
		return c.emitInexistentDiag(propNode,
			sampleWithout.TypeName()+"."+propName,
			diag.JSC_INEXISTENT_PROPERTY, diag.Error,
			fmt.Sprintf("Property %s never defined on %s",
				propName, sampleWithout.TypeName()))
	}
	return c.emitInexistentDiag(propNode,
		"union:"+u.TypeName()+"."+propName,
		diag.JSC_POSSIBLE_INEXISTENT_PROPERTY, diag.Warning,
		fmt.Sprintf("Property %s not defined on all member types of %s",
			propName, u.TypeName()))
}

// emitInexistentDiag dedups on `key` via c.inexistentReported and
// returns a single diagnostic when this is the first time the pair
// has been reported for the compile. Lazily initialises the dedup
// map for callers that allocate a TypeCheck without it.
func (c *TypeCheck) emitInexistentDiag(
	propNode *rhino.Node, key, code string, level diag.Level, msg string,
) []diag.Diagnostic {
	if c.inexistentReported[key] {
		return nil
	}
	if c.inexistentReported == nil {
		c.inexistentReported = map[string]bool{}
	}
	c.inexistentReported[key] = true
	return []diag.Diagnostic{{
		Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
		Code: code, Level: level, Message: msg,
	}}
}

// isThisReceiverRead reports whether `n` is a `this.prop` READ —
// i.e. the GETPROP's receiver is `this` and the parent isn't an
// ASSIGN with `n` as LHS. Java fires INEXISTENT_PROPERTY on this
// pattern even for bare @constructor classes (non-struct).
func isThisReceiverRead(n *rhino.Node) bool {
	kids := n.Children()
	if len(kids) < 1 {
		return false
	}
	if kids[0].Token() != rhino.THIS {
		return false
	}
	return !isAssignLHS(n)
}

// isAssignLHS reports whether `n` (typically a GETPROP) sits on the
// LHS of an enclosing ASSIGN. The ILLEGAL_PROPERTY_CREATION check
// only fires for struct-instance writes, not reads.
func isAssignLHS(n *rhino.Node) bool {
	p := n.Parent()
	if p == nil {
		return false
	}
	if p.Token() != rhino.ASSIGN {
		return false
	}
	return p.FirstChild() == n
}

// insideOwnConstructor reports whether `n` is lexically inside a
// @constructor function whose Instance matches `obj`. Java allows
// `this.newProp = x` inside the struct's own ctor body even though
// the surface is otherwise sealed.
func insideOwnConstructor(n *rhino.Node, obj *jstype.ObjectType) bool {
	for cur := n.Parent(); cur != nil; cur = cur.Parent() {
		if cur.Token() != rhino.FUNCTION {
			continue
		}
		info := cur.JSDocInfo()
		if info == nil || !info.IsConstructor() {
			continue
		}
		nameKid := cur.FirstChild()
		if nameKid == nil || nameKid.String_() == "" {
			continue
		}
		if nameKid.String_() == obj.Name {
			return true
		}
	}
	return false
}

// visitReturn validates `return <expr>;` against the enclosing
// function's @return annotation. The enclosing function is found by
// walking Node parent links — mirroring how Java's TypeCheck tracks
// the current function through NodeUtil.getEnclosingFunction.
func (c *TypeCheck) visitReturn(ret *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	enclosing := findEnclosingFunction(ret)
	if enclosing == nil {
		return nil
	}
	info := functionJSDoc(enclosing)
	if info == nil {
		return nil
	}
	var declared jstype.JSType
	if info.ReturnType() != nil {
		declared = scope.ResolveExpr(info.ReturnType().TypeName)
	} else if info.HasType() && info.Type() != nil {
		// No explicit `@return`, but a `@type {function(...):R}` on
		// the function declaration carries the return type. Resolve
		// the callable-expression form and read its Return slot —
		// mirrors how Java's TypeInference unpacks the declared
		// FunctionType for a function declared via @type.
		if ft, ok := scope.ResolveExpr(info.Type().TypeName).(*jstype.FunctionType); ok {
			declared = ft.Return
		}
	}
	if declared == nil {
		return nil
	}
	retExpr := ret.FirstChild()
	if retExpr == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	found := infer.InferType(retExpr)
	if jstype.IsSubtype(found, declared) {
		return nil
	}
	msg := fmt.Sprintf("inconsistent return type\nfound   : %s\nrequired: %s",
		found.TypeName(), declared.TypeName())
	line, col := diagPosForArg(retExpr)
	return []diag.Diagnostic{{
		Path: c.path, Line: line, Col: col,
		Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
	}}
}

// visitClassField — port of the class-field branch inside Java's
// visitClass. Emits "assignment to property X of C" at the field
// name's position when the initializer type doesn't match the
// declared @type.
func (c *TypeCheck) visitClassField(field *rhino.Node, parent *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := field.Children()
	if len(kids) < 2 {
		return nil
	}
	name := kids[0]
	init := kids[1]
	info := field.JSDocInfo()
	if info == nil || !info.HasType() {
		return nil
	}
	declared := scope.ResolveExpr(info.Type().TypeName)
	infer := NewTypeInference(scope)
	found := infer.InferType(init)
	if jstype.IsSubtype(found, declared) {
		return nil
	}
	className := ""
	if parent != nil && parent.Parent() != nil {
		gp := parent.Parent()
		if gp.Token() == rhino.CLASS && gp.FirstChild() != nil {
			className = gp.FirstChild().String_()
		}
	}
	msg := fmt.Sprintf("assignment to property %s of %s\nfound   : %s\nrequired: %s",
		name.String_(), className, found.TypeName(), declared.TypeName())
	// Java reports at the CLASS_FIELD_DEF node's position — the
	// `static`/name column, not the bare IDENT — so we use field's
	// pos rather than name's.
	return []diag.Diagnostic{{
		Path: c.path, Line: field.Lineno(), Col: field.Charno(),
		Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
	}}
}

// ---- helpers ----------------------------------------------------- //

// describeCallee renders a NAME or GETPROP chain as a dotted string
// for the "actual parameter N of X" diagnostic, matching how Java
// formats the callee in TypeCheck.java's visitCall message.
func describeCallee(n *rhino.Node) string {
	switch n.Token() {
	case rhino.NAME:
		return n.String_()
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) >= 2 {
			return describeCallee(kids[0]) + "." + kids[1].String_()
		}
	case rhino.THIS:
		return "this"
	}
	return "<fn>"
}

// diagPosForArg returns the (line, col) pair jsgo reports for a
// call-site argument diagnostic. Java's TypeCheck uses the
// rightmost segment of a GETPROP chain (`a.b.c` reports at `c`);
// for everything else the node's own position is used.
func diagPosForArg(n *rhino.Node) (int, int) {
	if n != nil && n.Token() == rhino.GETPROP {
		kids := n.Children()
		if len(kids) >= 2 {
			prop := kids[1]
			return prop.Lineno(), prop.Charno()
		}
	}
	return n.Lineno(), n.Charno()
}

// findEnclosingFunction walks Node parent links to the nearest
// FUNCTION ancestor. Mirrors NodeUtil.getEnclosingFunction in Java.
func findEnclosingFunction(n *rhino.Node) *rhino.Node {
	for cur := n.Parent(); cur != nil; cur = cur.Parent() {
		if cur.Token() == rhino.FUNCTION {
			return cur
		}
	}
	return nil
}

// functionJSDoc returns the JSDocInfo attached to `fn`, or — for
// anonymous function-expression forms like
// `/** @return {T} */ X.prototype.m = function(){…}` — the JSDoc
// attached to the enclosing VAR / ASSIGN / EXPR_RESULT statement.
// Mirrors how Java's NodeUtil.getBestJSDocInfo hops up the tree
// looking for the nearest declaration-level comment.
func functionJSDoc(fn *rhino.Node) *rhino.JSDocInfo {
	if fn == nil {
		return nil
	}
	if info := fn.JSDocInfo(); info != nil {
		return info
	}
	// Nested inside `X.prototype.m = fn` — climb via ASSIGN → EXPR_RESULT.
	for cur := fn.Parent(); cur != nil; cur = cur.Parent() {
		if info := cur.JSDocInfo(); info != nil {
			return info
		}
		switch cur.Token() {
		case rhino.ASSIGN, rhino.EXPR_RESULT, rhino.VAR, rhino.LET, rhino.CONST,
			rhino.NAME:
			continue
		}
		break
	}
	return nil
}

// ---- GETELEM + struct/dict access checks ------------------------ //

// visitGetElem — port of TypeCheck.java's visitGetElem (the
// diagnostics half). Emits:
//
//  1. JSC_ILLEGAL_PROPERTY_ACCESS on a `@struct` receiver,
//  2. JSC_TYPE_MISMATCH "restricted index type" when the receiver is
//     an Array (or Array<T>) and the index is not number-compatible.
//
// The "only arrays or objects can be accessed" diagnostic from Java
// lives under STRICT_CHECK_TYPES and is off by default — skipped
// here to match the oracle.
func (c *TypeCheck) visitGetElem(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	recv := kids[0]
	idx := kids[1]
	infer := NewTypeInference(scope)
	rawRecv := infer.InferType(recv)
	// `null[…]` / `undefined[…]` are runtime TypeErrors — Java fires
	// "only arrays or objects can be accessed" at the whole GETELEM
	// position (not at the index). Mirrors the expectNotNullOrUndefined
	// branch of TypeCheck.visitGetElem.
	if rawRecv == jstype.Null || rawRecv == jstype.Void {
		// Dedup per (receiver-name, void/null type) so a sequence of
		// `bar[0]; bar[1]; bar[2];` fires the diagnostic once at the
		// first occurrence, matching Java's TypeCheck behaviour.
		key := "<null-recv>"
		if recv.Token() == rhino.NAME {
			key = recv.String_() + "." + rawRecv.TypeName()
		}
		if c.inexistentReported == nil {
			c.inexistentReported = map[string]bool{}
		}
		if c.inexistentReported[key] {
			return nil
		}
		c.inexistentReported[key] = true
		return []diag.Diagnostic{{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
			Message: "only arrays or objects can be accessed\nfound   : " +
				rawRecv.TypeName() + "\nrequired: Object",
		}}
	}
	recvType := stripNullUndef(rawRecv)
	if isStructReceiver(recvType) && !isWellKnownSymbolIndex(idx) {
		return []diag.Diagnostic{{
			Path: c.path, Line: idx.Lineno(), Col: idx.Charno(),
			Code: diag.JSC_ILLEGAL_PROPERTY_ACCESS, Level: diag.Error,
			Message: "Cannot do '[]' access on a struct",
		}}
	}
	if isArrayReceiver(recvType) {
		idxType := infer.InferType(idx)
		if !isNumberIndexOk(idxType) {
			return []diag.Diagnostic{{
				Path: c.path, Line: idx.Lineno(), Col: idx.Charno(),
				Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
				Message: "restricted index type\nfound   : " + idxType.TypeName() +
					"\nrequired: number",
				}}
		}
	}
	// Object<K,V> receiver with a non-stringifiable K — CheckObjectKeys
	// also scans GETELEM/COMPUTED_PROP accesses where the receiver's
	// declared type is an `Object<…>` instantiation whose key type
	// won't coerce cleanly through Object.prototype.toString. Mirrors
	// TypeCheck.java's validateNonStringifiableKey call from
	// visitGetElem / visitGetProp.
	if keyT := objectTemplateKey(recvType); keyT != nil && !isStringifiableJSType(keyT) {
		return []diag.Diagnostic{{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code: diag.JSC_NON_STRINGIFIABLE_OBJECT_KEY, Level: diag.Error,
			Message: "Object type \"" + recvType.TypeName() +
				"\" contains non-stringifiable key and it may lead to an error. " +
				"Please use ES6 Map instead or implement your own Map structure.",
		}}
	}
	return nil
}

// objectTemplateKey returns the K argument of an `Object<K,V>`-style
// templatized receiver, or nil if `t` isn't such a type. Unions
// collapse via a single-member probe: if every non-null/undef member
// of the union names the same Object<K,_>, the first K is returned.
// Mirrors the shape TypeValidator.getObjectKeyType reads in Java.
func objectTemplateKey(t jstype.JSType) jstype.JSType {
	if t == nil {
		return nil
	}
	if tt, ok := t.(*jstype.TemplatizedType); ok && tt.Base != nil &&
		tt.Base.Name == "Object" && len(tt.Args) >= 1 {
		return tt.Args[0]
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if m == jstype.Null || m == jstype.Void {
				continue
			}
			if k := objectTemplateKey(m); k != nil {
				return k
			}
		}
	}
	return nil
}

// visitArrayLit is the call-through hook for ARRAYLIT nodes — port
// of TypeCheck.java's visitArrayLit, which does nothing for untyped
// literals and defers per-element subtype validation to the
// initializer / assignment path. Kept as a seam so future strict-
// mode extensions (tuple types, spread-element checks) can hang off
// the same dispatch site without another Visit switch edit.
func (c *TypeCheck) visitArrayLit(_ *rhino.Node, _ *TypedScope) []diag.Diagnostic {
	return nil
}

// visitObjectLit is the call-through hook for OBJECTLIT nodes — port
// of TypeCheck.java's visitObjectLit, which validates computed keys
// against any declared index signature and defers per-value subtype
// checks to the surrounding assignment. The structural-shape
// inference happens in TypeInference.inferObjectLit; this method
// exists so Visit has an OBJECTLIT seam mirroring Java's dispatch
// surface.
func (c *TypeCheck) visitObjectLit(_ *rhino.Node, _ *TypedScope) []diag.Diagnostic {
	return nil
}

// isArrayReceiver reports whether `t` dereferences to the builtin
// `Array` instance (or Array<T>). Mirrors the Array-test branch in
// TypeValidator.expectIndexMatch.
func isArrayReceiver(t jstype.JSType) bool {
	obj := dereferencedObject(t)
	return obj != nil && obj.Name == "Array"
}

// isNumberIndexOk reports whether `t` is acceptable as an Array
// index. Java accepts number / unknown / union-with-number; any
// definite non-number fires. Uses jstype.IsSubtype so
// EnumElementType<number> (which is a subtype of Number via its
// Underlying rule) still passes.
func isNumberIndexOk(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	if isUnknownOrNil(t) {
		return true
	}
	return jstype.IsSubtype(t, jstype.Number)
}

// isKnownSurfaceBuiltin reports whether `name` is one of the
// built-in constructors whose property surface jsgo ships
// concretely enough to fire JSC_INEXISTENT_PROPERTY on missing
// members. Array is the canonical case — our default scope
// populates length/push/pop/…, so `[].foo` cleanly fires.
func isKnownSurfaceBuiltin(name string) bool {
	switch name {
	case "Array":
		return true
	}
	return false
}

// defineOnOpenObject registers `prop` on every non-struct,
// non-structural ObjectType reachable through `t` (direct,
// Templatized, or union member). Used so `X.prop = …` assignments
// extend the instance's property surface, matching Java's
// "add to slot map on assign" behaviour.
func defineOnOpenObject(t jstype.JSType, prop string) {
	if t == nil {
		return
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			defineOnOpenObject(m, prop)
		}
		return
	}
	obj := dereferencedObject(t)
	if obj == nil || obj.Name == "" {
		return
	}
	if obj.IsStruct || obj.IsStructural {
		return
	}
	if obj.HasProperty(prop) {
		return
	}
	obj.DefineProperty(prop, jstype.Unknown)
}

// dereferencedObject returns the ObjectType behind a templatized
// receiver type. Returns nil for primitives / Unknown / unions /
// anonymous types. Used by the property-access diagnostics to
// reach the underlying nominal slot (IsStruct / IsDict /
// IsStructural). Mirrors JSType.toObjectType in Java.
func dereferencedObject(t jstype.JSType) *jstype.ObjectType {
	if t == nil {
		return nil
	}
	if tt, ok := t.(*jstype.TemplatizedType); ok {
		t = tt.Base
	}
	if obj, ok := t.(*jstype.ObjectType); ok {
		return obj
	}
	return nil
}

// isStructReceiver reports whether `t` is (or contains) a @struct
// ObjectType. Java's `x['prop']` check fires whenever ANY member
// of a union is @struct — the union promises a struct is possible
// at runtime, and bypassing the access check there would let the
// program sneak in a bad lookup via the union's other arm.
func isStructReceiver(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if obj := dereferencedObject(t); obj != nil {
		return obj.IsStruct
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if isStructReceiver(m) {
				return true
			}
		}
	}
	return false
}

// isDictReceiver is the @dict mirror of isStructReceiver.
func isDictReceiver(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if obj := dereferencedObject(t); obj != nil {
		return obj.IsDict
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if isDictReceiver(m) {
				return true
			}
		}
	}
	return false
}

// ---- strict operand validators ---------------------------------- //
//
// The methods below implement the Java `TypeValidator`
// expectNumber / expectNumberOrString / expectMatchingTypes helpers
// for binary operators. They fire only when the caller gated them
// behind `GroupStrictCheckTypes` — matching how Java's
// `DiagnosticGroups.STRICT_CHECK_TYPES` elevates the relevant
// DiagnosticTypes from OFF to WARNING / ERROR.

// visitThis fires JSC_USED_GLOBAL_THIS when a `this` keyword
// appears outside a method / constructor / class context (i.e.
// the global object). Mirrors TypeCheck.java's visitThis branch.
func (c *TypeCheck) visitThis(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	if scope.ThisType() != nil {
		return nil
	}
	// Walk ancestors looking for a FUNCTION / CLASS that would
	// bind `this` at runtime. Arrow functions don't bind this,
	// but @this JSDoc does — both checked via the enclosing
	// function's JSDocInfo.
	for cur := n.Parent(); cur != nil; cur = cur.Parent() {
		if cur.Token() == rhino.FUNCTION {
			// Non-arrow FUNCTION binds this implicitly; we stay
			// quiet even when we don't know the bound type.
			if !cur.GetBooleanProp(rhino.PropArrowFn) {
				return nil
			}
		}
	}
	return []diag.Diagnostic{{
		Path: c.path, Line: n.Lineno(), Col: n.Charno(),
		Code: diag.JSC_USED_GLOBAL_THIS, Level: diag.Error,
		Message: "dangerous use of the global this object",
	}}
}

// visitFunctionParams fires JSC_INEXISTENT_PARAM when a @param
// JSDoc names a parameter that doesn't appear in the function's
// actual parameter list. Java reports at the FUNCTION keyword's
// column.
func (c *TypeCheck) visitFunctionParams(fn *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	info := fn.JSDocInfo()
	if info == nil {
		return nil
	}
	// CheckObjectKeys also scans @param / @return types on the
	// function's JSDoc. A signature like
	// `@param {!Object<Foo, number>} m` fires at the FUNCTION
	// position just like a VAR-level non-stringifiable key does.
	out := c.visitFunctionSigObjectKeys(fn, info, scope)
	names := info.ParameterNames()
	if len(names) == 0 {
		return out
	}
	kids := fn.Children()
	if len(kids) < 2 {
		return out
	}
	params := kids[1]
	paramSet := map[string]bool{}
	for _, p := range params.Children() {
		paramSet[p.String_()] = true
	}
	fnName := kids[0].String_()
	if fnName == "" {
		fnName = "<anonymous>"
	}
	for _, name := range names {
		if paramSet[name] {
			continue
		}
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: fn.Lineno(), Col: fn.Charno(),
			Code: diag.JSC_INEXISTENT_PARAM, Level: diag.Error,
			Message: fmt.Sprintf("parameter %s does not appear in %s's parameter list", name, fnName),
		})
	}
	return out
}

// visitFunctionSigObjectKeys scans the @param / @return type
// expressions on `fn`'s JSDoc and fires
// JSC_NON_STRINGIFIABLE_OBJECT_KEY for any embedded
// `Object<K,…>` whose K isn't stringifiable. Mirrors the
// CheckObjectKeys walk over Java's JSDocInfo.getParameterType /
// getReturnType — the per-declaration scan in visitObjectKeyTypes
// handles @type only, which misses function signatures entirely.
func (c *TypeCheck) visitFunctionSigObjectKeys(fn *rhino.Node, info *rhino.JSDocInfo, scope *TypedScope) []diag.Diagnostic {
	var out []diag.Diagnostic
	line, col := fn.Lineno(), fn.Charno()
	fire := func(keyExpr string) {
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: line, Col: col,
			Code: diag.JSC_NON_STRINGIFIABLE_OBJECT_KEY, Level: diag.Error,
			Message: fmt.Sprintf("Object type %q contains non-stringifiable key "+
				"and it may lead to an error. Please use ES6 Map instead or "+
				"implement your own Map structure.", "Object<"+keyExpr+",?>"),
		})
	}
	for _, name := range info.ParameterNames() {
		if pt := info.ParameterType(name); pt != nil {
			walkObjectKeyExpr(pt.TypeName, scope, fire)
		}
	}
	if rt := info.ReturnType(); rt != nil {
		walkObjectKeyExpr(rt.TypeName, scope, fire)
	}
	return out
}

// visitObjectKeyTypes walks every `Object<K, V>` occurrence inside
// the @type/@param/@return expression attached to `v` and fires
// JSC_NON_STRINGIFIABLE_OBJECT_KEY when K isn't a stringifiable
// primitive (string / number / boolean / symbol / void / null, or
// unions / nullable variants thereof).
//
// Ports com/google/javascript/jscomp/CheckObjectKeys.java: Java walks
// the resolved JSType tree; jsgo walks the raw JSDoc text since the
// checker operates against TypeName strings. The diagnostic surfaces
// at the declaring node's position (the VAR / LET / CONST keyword).
//
// Only the on-by-default TypeCheck pathway is ported; Java's
// separate CheckObjectKeys traversal is folded in here because its
// sole emission path is JSC_NON_STRINGIFIABLE_OBJECT_KEY and the
// scope-aware resolution logic already lives in this pass.
func (c *TypeCheck) visitObjectKeyTypes(v *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	info := v.JSDocInfo()
	if info == nil || info.Type() == nil {
		return nil
	}
	line, col := v.Lineno(), v.Charno()
	// Fall back to the first NAME child for position when the decl
	// node itself has no pos (parser only tags names, not the VAR
	// shell). Matches where Java anchors the warning — at the
	// "name-of-variable" token.
	if line == 0 && col == 0 {
		for _, name := range v.Children() {
			if name.Token() == rhino.NAME {
				line, col = name.Lineno(), name.Charno()
				break
			}
		}
	}
	var out []diag.Diagnostic
	walkObjectKeyExpr(info.Type().TypeName, scope, func(keyExpr string) {
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: line, Col: col,
			Code: diag.JSC_NON_STRINGIFIABLE_OBJECT_KEY, Level: diag.Error,
			Message: fmt.Sprintf("Object type %q contains non-stringifiable key "+
				"and it may lead to an error. Please use ES6 Map instead or "+
				"implement your own Map structure.", "Object<"+keyExpr+",?>"),
		})
	})
	return out
}

// walkObjectKeyExpr scans `expr` for every `Object<K, V>` occurrence
// (including nested under nullable / unions / function / templatized
// wrappers) and invokes `emit` with the raw K expression when K is
// *not* a stringifiable type. Mirrors Java's
// CheckObjectKeys.checkType / JSType.isStringifiable recursion.
func walkObjectKeyExpr(expr string, scope *TypedScope, emit func(string)) {
	for {
		idx := findTopLevelObjectLt(expr)
		if idx < 0 {
			return
		}
		close := matchingAngleClose(expr, idx)
		if close < 0 {
			return
		}
		inside := expr[idx+1 : close]
		// Split into K, V (comma at nesting depth 0). Walk K for
		// nested Object<…> too.
		k, v := splitFirstTopLevelComma(inside)
		if k != "" && !isStringifiableTypeExpr(k, scope) {
			emit(k)
		}
		// Recurse into both args so `Object<string, Object<!Foo, T>>`
		// fires on the inner Object as well (Java's recursive walk).
		walkObjectKeyExpr(k, scope, emit)
		walkObjectKeyExpr(v, scope, emit)
		expr = expr[close+1:]
	}
}

// findTopLevelObjectLt returns the index of the `<` in an
// `Object<...>` occurrence — scanning for the bare name `Object`
// (not a qualified name suffix) followed immediately by `<`.
// Returns -1 if none found.
func findTopLevelObjectLt(expr string) int {
	for i := 0; i < len(expr); i++ {
		if i+6 > len(expr) {
			return -1
		}
		if expr[i] != 'O' {
			continue
		}
		// Match "Object" as a whole word.
		if expr[i:i+6] != "Object" {
			continue
		}
		// Preceding char must not be part of an identifier or dot.
		if i > 0 {
			p := expr[i-1]
			if p == '.' || p == '_' || p == '$' ||
				(p >= 'a' && p <= 'z') || (p >= 'A' && p <= 'Z') ||
				(p >= '0' && p <= '9') {
				continue
			}
		}
		// Must be followed by `<` — allow whitespace between.
		j := i + 6
		for j < len(expr) && (expr[j] == ' ' || expr[j] == '\t') {
			j++
		}
		if j < len(expr) && expr[j] == '<' {
			return j
		}
	}
	return -1
}

// matchingAngleClose returns the index of the `>` that matches the
// `<` at position `open` in `expr`, or -1 if unbalanced.
func matchingAngleClose(expr string, open int) int {
	depth := 0
	paren := 0
	curly := 0
	for i := open; i < len(expr); i++ {
		switch expr[i] {
		case '<':
			depth++
		case '>':
			depth--
			if depth == 0 {
				return i
			}
		case '(':
			paren++
		case ')':
			paren--
		case '{':
			curly++
		case '}':
			curly--
		}
		_ = paren
		_ = curly
	}
	return -1
}

// splitFirstTopLevelComma returns the prefix up to the first comma
// at nesting depth 0 and the suffix after it. Used to peel K off
// `Object<K, V>`.
func splitFirstTopLevelComma(s string) (string, string) {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '<', '(', '{', '[':
			depth++
		case '>', ')', '}', ']':
			depth--
		case ',':
			if depth == 0 {
				return trimSpaces(s[:i]), trimSpaces(s[i+1:])
			}
		}
	}
	return trimSpaces(s), ""
}

// trimSpaces is a byte-level strings.TrimSpace; the file avoids
// importing strings per the parser's local style.
func trimSpaces(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 {
		c := s[len(s)-1]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}

// isStringifiableTypeExpr reports whether the resolved form of `expr`
// is safe to coerce to a string (i.e. a legal `Object<…>` key).
// Java's JSType.isStringifiable: primitives (string, number, boolean,
// symbol), plus unions where every member is stringifiable, plus
// `?`/`*` (unknown/all accept). Nullable (`?T`) and non-null (`!T`)
// prefixes are transparent. Anything else — Object, Array, function,
// enum-of-object, user class — is NON stringifiable.
func isStringifiableTypeExpr(expr string, scope *TypedScope) bool {
	expr = trimSpaces(expr)
	if expr == "" {
		return true
	}
	// `?` / `*` (unknown / all) — Java's CheckObjectKeys skips.
	if expr == "?" || expr == "*" {
		return true
	}
	// Strip nullable / non-null prefix.
	if expr[0] == '?' || expr[0] == '!' {
		return isStringifiableTypeExpr(expr[1:], scope)
	}
	// Strip outer parens.
	if expr[0] == '(' && expr[len(expr)-1] == ')' && balancedOuterExpr(expr) {
		return isStringifiableTypeExpr(expr[1:len(expr)-1], scope)
	}
	// Union — every member must be stringifiable.
	if parts := splitTopLevelUnionExpr(expr); len(parts) > 1 {
		for _, p := range parts {
			if !isStringifiableTypeExpr(p, scope) {
				return false
			}
		}
		return true
	}
	// Primitive name — Java's allowlist.
	switch expr {
	case "string", "number", "boolean", "symbol", "void", "undefined", "null":
		return true
	}
	// Templatized types like `Array<T>` / `Object<K,V>` — the outer
	// constructor itself is never stringifiable (Java's rule: the
	// instance coerces via Object.prototype.toString which returns
	// `[object Object]` and collides across keys).
	if findTopLevelObjectLt(expr) >= 0 {
		return false
	}
	// Named type lookup via the scope. An EnumType over a
	// stringifiable primitive underlying is itself stringifiable
	// (Java's EnumType.isStringifiable). Anything else — nominal
	// object, interface, class, function — is not.
	resolved := scope.ResolveExpr(expr)
	if resolved == nil {
		return false
	}
	return isStringifiableJSType(resolved)
}

// balancedOuterExpr reports whether `expr` is wrapped in a single
// matched pair of outer parens. Local helper so TypeCheck doesn't
// pull in scope.go's internal implementation.
func balancedOuterExpr(expr string) bool {
	if len(expr) < 2 || expr[0] != '(' || expr[len(expr)-1] != ')' {
		return false
	}
	depth := 0
	for i := 0; i < len(expr); i++ {
		switch expr[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 && i != len(expr)-1 {
				return false
			}
		}
	}
	return depth == 0
}

// splitTopLevelUnionExpr returns the `|`-separated members of a
// top-level union expression. Parens / angles / braces raise the
// depth so `(A|B)|C` splits into two parts, not three.
func splitTopLevelUnionExpr(expr string) []string {
	var out []string
	depth := 0
	last := 0
	for i := 0; i < len(expr); i++ {
		switch expr[i] {
		case '<', '(', '{', '[':
			depth++
		case '>', ')', '}', ']':
			depth--
		case '|':
			if depth == 0 {
				out = append(out, trimSpaces(expr[last:i]))
				last = i + 1
			}
		}
	}
	out = append(out, trimSpaces(expr[last:]))
	if len(out) == 1 {
		return nil
	}
	return out
}

// isStringifiableJSType walks a resolved JSType and reports whether
// every alternative coerces cleanly to a string key. Mirrors Java's
// JSType.isStringifiable dispatch.
func isStringifiableJSType(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	switch v := t.(type) {
	case *jstype.UnionType:
		for _, m := range v.Members() {
			if !isStringifiableJSType(m) {
				return false
			}
		}
		return true
	case *jstype.EnumType:
		return isStringifiableJSType(v.Underlying)
	}
	if t == jstype.Unknown {
		return true
	}
	switch t.TypeName() {
	case "string", "number", "boolean", "symbol", "null", "undefined", "void":
		return true
	}
	return false
}

// visitCast fires JSC_INVALID_CAST when a `/** @type {T} */(expr)`
// annotation casts between two types that have no subtype / supertype
// relationship. Mirrors Java's TypeValidator.expectCanCast, which
// rejects casts where neither `sourceType.isSubtypeOf(targetType)`
// nor `targetType.isSubtypeOf(sourceType)` holds. Unknown on either
// side is always accepted (Java treats `?` as "anything permitted").
func (c *TypeCheck) visitCast(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	info := n.JSDocInfo()
	if info == nil || info.Type() == nil {
		return nil
	}
	kid := n.FirstChild()
	if kid == nil {
		return nil
	}
	to := scope.ResolveExpr(info.Type().TypeName)
	if to == nil {
		return nil
	}
	from := NewTypeInference(scope).InferType(kid)
	if from == nil {
		return nil
	}
	if isUnknownLike(to) || isUnknownLike(from) {
		return nil
	}
	if jstype.IsSubtype(from, to) || jstype.IsSubtype(to, from) {
		return nil
	}
	return []diag.Diagnostic{{
		Path: c.path, Line: n.Lineno(), Col: n.Charno(),
		Code: diag.JSC_INVALID_CAST, Level: diag.Error,
		Message: fmt.Sprintf("invalid cast - must be a subtype or supertype\nfrom: %s\nto  : %s",
			from.TypeName(), to.TypeName()),
	}}
}

// isUnknownLike reports whether `t` is the Unknown or all (`*`) type.
// Java's TypeValidator bypasses cast checks whenever either side is
// one of these wildcards, since they subsume every concrete type.
func isUnknownLike(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	return t == jstype.Unknown || t.TypeName() == "*" || t.TypeName() == "?"
}

// visitBitOp fires JSC_BAD_TYPE_FOR_BIT_OPERATION when either
// operand of a bitwise operator (`<<`, `>>`, `>>>`, `&`, `|`,
// `^`) is definitely not a number. Java's TypeCheck rejects
// void / null / boolean / object operands statically.
func (c *TypeCheck) visitBitOp(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	infer := NewTypeInference(scope)
	var out []diag.Diagnostic
	op := bitOpSymbol(n.Token())
	for _, k := range kids {
		t := infer.InferType(k)
		if !isDefinitelyNotNumber(t) {
			continue
		}
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: k.Lineno(), Col: k.Charno(),
			Code: diag.JSC_BAD_TYPE_FOR_BIT_OPERATION, Level: diag.Error,
			Message: fmt.Sprintf("operator %s cannot be applied to %s", op, t.TypeName()),
		})
		break // report the first bad operand only, matching Java
	}
	return out
}

// bitOpSymbol returns the textual operator for a bit-op token.
func bitOpSymbol(t rhino.Token) string {
	switch t {
	case rhino.LSH:
		return "<<"
	case rhino.RSH:
		return ">>"
	case rhino.URSH:
		return ">>>"
	case rhino.BITAND:
		return "&"
	case rhino.BITOR:
		return "|"
	case rhino.BITXOR:
		return "^"
	}
	return "?"
}

// isWellKnownSymbolIndex reports whether `idx` is a GETPROP on the
// \`Symbol\` builtin (e.g. \`Symbol.iterator\`, \`Symbol.asyncIterator\`).
// Java's TypeCheck allows well-known symbols to index a @struct
// receiver since they're the canonical protocol entry points.
func isWellKnownSymbolIndex(idx *rhino.Node) bool {
	if idx == nil || idx.Token() != rhino.GETPROP {
		return false
	}
	kids := idx.Children()
	if len(kids) < 2 {
		return false
	}
	recv := kids[0]
	return recv.Token() == rhino.NAME && recv.String_() == "Symbol"
}

// unwrapCallableUnion extracts the single FunctionType member of a
// `FunctionType | null | undefined` union. Returns (nil, false) when
// the union mixes multiple functions, non-null non-function members,
// or isn't a union at all. Used by visitCall / visitNew to thread a
// nullable-callable parameter (/** @param {?function(...)} */)
// through the arity-check branch. Mirrors Java's
// TypeCheck.validateCall, which calls `JSType.restrictByNotNullOrUndefined`
// before the isFunctionType branch.
func unwrapCallableUnion(t jstype.JSType) (*jstype.FunctionType, bool) {
	u, ok := t.(*jstype.UnionType)
	if !ok {
		return nil, false
	}
	var fn *jstype.FunctionType
	for _, m := range u.Members() {
		switch x := m.(type) {
		case *jstype.FunctionType:
			if fn != nil {
				return nil, false
			}
			fn = x
		default:
			if m == jstype.Null || m == jstype.Void {
				continue
			}
			return nil, false
		}
	}
	if fn == nil {
		return nil, false
	}
	return fn, true
}

// isDefinitelyNotCallable reports whether `t` is a type that can
// never be a function at runtime. Primitives (number, string,
// boolean, null, void) qualify; a union does iff every member is
// non-callable. Unknown / All stay callable-ambiguous.
func isDefinitelyNotCallable(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if t == jstype.Unknown || t == jstype.All {
		return false
	}
	switch t {
	case jstype.Number, jstype.String, jstype.Boolean, jstype.Null, jstype.Void:
		return true
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if _, isFn := m.(*jstype.FunctionType); isFn {
				return false
			}
			if obj := dereferencedObject(m); obj != nil && obj.Name == "Function" {
				return false
			}
		}
		return true
	}
	return false
}

// isDefinitelyNotNumber reports whether `t` is a primitive type
// known to not be a number at runtime. Void, Null, Boolean, and
// String are the non-numeric primitives; Unknown/All still allow
// numeric values so they don't trigger.
func isDefinitelyNotNumber(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	switch t {
	case jstype.Void, jstype.Null, jstype.Boolean:
		return true
	}
	return false
}

// visitIn fires diagnostics for the IN operator:
//  1. JSC_IN_USED_WITH_STRUCT when the RHS is a @struct receiver.
//  2. JSC_TYPE_MISMATCH when the LHS is definitely not string /
//     symbol / number (Java's requirement for the 'in' operand).
func (c *TypeCheck) visitIn(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	lhs, rhs := kids[0], kids[1]
	infer := NewTypeInference(scope)
	recv := stripNullUndef(infer.InferType(rhs))
	var out []diag.Diagnostic
	lt := infer.InferType(lhs)
	if isDefinitelyNotInKey(lt) {
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
			Message: fmt.Sprintf(
				"left side of 'in'\nfound   : %s\nrequired: (string|symbol)",
				lt.TypeName()),
		})
	}
	if isStructReceiver(recv) {
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: rhs.Lineno(), Col: rhs.Charno(),
			Code: diag.JSC_IN_USED_WITH_STRUCT, Level: diag.Error,
			Message: "Cannot use the IN operator with structs",
		})
	}
	return out
}

// isDefinitelyNotInKey reports whether `t` is a type that's
// definitely not a valid `X in Y` LHS. Java's TypeCheck accepts
// String / Number / Symbol / Unknown / All / Null / Undefined;
// only firmly-nominal FunctionType / ObjectType references fail.
func isDefinitelyNotInKey(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	return false
}

// visitForIn fires JSC_IN_USED_WITH_STRUCT on `for (var k in struct)`.
// Same rule as visitIn, but the AST shape is FOR_IN(var, iterable, body).
func (c *TypeCheck) visitForIn(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	iter := kids[1]
	infer := NewTypeInference(scope)
	recv := stripNullUndef(infer.InferType(iter))
	if !isStructReceiver(recv) {
		return nil
	}
	return []diag.Diagnostic{{
		Path: c.path, Line: iter.Lineno(), Col: iter.Charno(),
		Code: diag.JSC_IN_USED_WITH_STRUCT, Level: diag.Error,
		Message: "Cannot use the IN operator with structs",
	}}
}

// visitEquality used to live here and fired JSC_DETERMINISTIC_TEST
// directly from TypeCheck. It's been hoisted to its own pass in
// checkdeterministictest.go, which keeps the always-true / always-
// false logic co-located with the rest of the static-comparison
// reasoning. TypeCheck's binary-op switch skips EQ / NE / SHEQ / SHNE
// so the two passes don't double-fire.

// visitAdd ports TypeValidator.expectNumberOrString on the ADD
// operator. JS `+` coerces to string whenever either side is
// string-typed; a mismatch is reported only when BOTH sides are
// incompatible (neither number nor string nor Unknown nor a union
// whose members all satisfy the predicate). Matches Java's
// conservative TypeValidator — Unknown / All silence.
func (c *TypeCheck) visitAdd(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	infer := NewTypeInference(scope)
	l := infer.InferType(kids[0])
	r := infer.InferType(kids[1])
	var out []diag.Diagnostic
	if !isAddOperandOk(l) {
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: kids[0].Lineno(), Col: kids[0].Charno(),
			Code: diag.JSC_INVALID_OPERAND_TYPE, Level: diag.Warning,
			Message: "left operand\nfound   : " + l.TypeName() +
				"\nrequired: (number|string)",
		})
	}
	if !isAddOperandOk(r) {
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: kids[1].Lineno(), Col: kids[1].Charno(),
			Code: diag.JSC_INVALID_OPERAND_TYPE, Level: diag.Warning,
			Message: "right operand\nfound   : " + r.TypeName() +
				"\nrequired: (number|string)",
		})
	}
	return out
}

// isAddOperandOk reports whether `t` is acceptable to JS `+`. Java's
// TypeValidator rule: number, string, boolean (ToNumber coercion),
// null / undefined (coerce to 0 / "undefined"), Unknown / All, or
// any union whose members all pass the predicate. Object-ish types
// implicitly coerce via ToPrimitive and Java doesn't flag them; we
// do the same.
func isAddOperandOk(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	if isUnknownOrNil(t) {
		return true
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if !isAddOperandOk(m) {
				return false
			}
		}
		return true
	}
	switch t {
	case jstype.Number, jstype.String, jstype.Boolean,
		jstype.Null, jstype.Void, jstype.All:
		return true
	}
	return true
}

// visitArithmetic ports TypeValidator.expectNumber for SUB / MUL /
// DIV / MOD. Both operands must be numeric-compatible; mismatch
// reports JSC_TYPE_MISMATCH at the offending operand.
func (c *TypeCheck) visitArithmetic(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	return c.visitNumericBinaryOp(n, scope)
}

// visitBitwise ports the BITAND/BITOR/BITXOR/LSH/RSH/URSH branch —
// same predicate as expectNumber.
func (c *TypeCheck) visitBitwise(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	return c.visitNumericBinaryOp(n, scope)
}

// visitNumericBinaryOp walks both operands and emits
// JSC_TYPE_MISMATCH "found: T\nrequired: number" when either side
// is definitely non-numeric. Stays quiet for Unknown / All /
// unions-with-number, matching Java's conservative default.
func (c *TypeCheck) visitNumericBinaryOp(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	infer := NewTypeInference(scope)
	var out []diag.Diagnostic
	for _, k := range kids {
		t := infer.InferType(k)
		if isNumericOperandOk(t) {
			continue
		}
		label := "operator " + n.Token().String()
		if n.Token() == rhino.SUB || n.Token() == rhino.MUL ||
			n.Token() == rhino.DIV || n.Token() == rhino.MOD {
			// Java uses "left operand" / "right operand" on the
			// arithmetic branches and "operator X" on bitwise. Infer
			// side from child index so the diagnostic message matches.
			if k == kids[0] {
				label = "left operand"
			} else {
				label = "right operand"
			}
		}
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: k.Lineno(), Col: k.Charno(),
			Code: diag.JSC_INVALID_OPERAND_TYPE, Level: diag.Warning,
			Message: label + "\nfound   : " + t.TypeName() + "\nrequired: number",
		})
	}
	return out
}

// isNumericOperandOk reports whether `t` is arithmetic-compatible.
// Java's TypeValidator treats null / undefined as coercible-to-0
// and silences diagnostics for them; we do the same to keep
// `var x; x + 1` quiet when x was declared but uninitialised.
func isNumericOperandOk(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	if isUnknownOrNil(t) {
		return true
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if !isNumericOperandOk(m) {
				return false
			}
		}
		return true
	}
	switch t {
	case jstype.Number, jstype.Boolean, jstype.Null, jstype.Void, jstype.All:
		return true
	}
	return false
}

// visitRelationalDefault fires JSC_TYPE_MISMATCH on the default-
// mode relational check: both \`<\` operands must be bigint-, number-
// or string-compatible. Non-primitives (Object / FunctionType) are
// definitively invalid. Matches Java's TypeCheck.visitBinaryOp
// default-mode relational branch.
func (c *TypeCheck) visitRelationalDefault(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	infer := NewTypeInference(scope)
	l := infer.InferType(kids[0])
	r := infer.InferType(kids[1])
	var out []diag.Diagnostic
	if isDefinitelyNotComparable(l) {
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: kids[0].Lineno(), Col: kids[0].Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
			Message: "left side of comparison\nfound   : " + l.TypeName() + "\nrequired: (bigint|number|string)",
		})
	}
	if isDefinitelyNotComparable(r) {
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: kids[1].Lineno(), Col: kids[1].Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
			Message: "right side of comparison\nfound   : " + r.TypeName() + "\nrequired: (bigint|number|string)",
		})
	}
	return out
}

// isDefinitelyNotComparable reports whether \`t\` is a type that
// definitely can't participate in a relational comparison. Objects
// and functions are rejected; primitives (including boolean —
// Java default mode permits true<0, coercing via ToNumber) and
// unknowns stay quiet.
func isDefinitelyNotComparable(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if t == jstype.Unknown || t == jstype.All {
		return false
	}
	switch t {
	case jstype.Number, jstype.String, jstype.Null, jstype.Void, jstype.Boolean:
		return false
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	if _, ok := t.(*jstype.ObjectType); ok {
		return true
	}
	return false
}

// visitRelativeComparison ports TypeValidator.expectMatchingTypes
// for LT / LE / GT / GE: both operands must be number-compatible
// OR both string-compatible. Mismatch fires at the right-hand side
// (Java's arbitrary-but-stable choice).
func (c *TypeCheck) visitRelativeComparison(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	infer := NewTypeInference(scope)
	l := infer.InferType(kids[0])
	r := infer.InferType(kids[1])
	if canCompareStrict(l, r) {
		return nil
	}
	return []diag.Diagnostic{{
		Path: c.path, Line: kids[1].Lineno(), Col: kids[1].Charno(),
		Code: diag.JSC_INVALID_OPERAND_TYPE, Level: diag.Warning,
		Message: "expected matching types in comparison\nfound   : " +
			r.TypeName() + "\nrequired: " + l.TypeName(),
	}}
}

// canCompareStrict reports whether two types may participate in a
// relative comparison. Java narrows to "both number-compatible OR
// both string-compatible", with Unknown / null / undefined always
// accepted to avoid false-positives on partially-annotated code.
func canCompareStrict(a, b jstype.JSType) bool {
	if a == nil || b == nil {
		return true
	}
	if isUnknownOrNil(a) || isUnknownOrNil(b) {
		return true
	}
	if isNumericOperandOk(a) && isNumericOperandOk(b) {
		return true
	}
	if isStringOperandOk(a) && isStringOperandOk(b) {
		return true
	}
	return false
}

// isStringOperandOk is the string-mirror of isNumericOperandOk —
// used by canCompareStrict when the LHS is string-shaped.
func isStringOperandOk(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	if isUnknownOrNil(t) {
		return true
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if !isStringOperandOk(m) {
				return false
			}
		}
		return true
	}
	return t == jstype.String || t == jstype.Null || t == jstype.Void || t == jstype.All
}
