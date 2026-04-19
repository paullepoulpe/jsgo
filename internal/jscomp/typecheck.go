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
	flows         map[*rhino.Node]*DataFlowAnalysis
	currentFlow   *DataFlowAnalysis
	currentCfg    *ControlFlowGraph
	flowStack     []*DataFlowAnalysis
	cfgStack      []*ControlFlowGraph
	currentStmt   *rhino.Node
	stmtStack     []*rhino.Node
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
func (c *TypeCheck) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.installFlow(script, c.global)
	c.collectPropertyWrites(script)
	c.fileSuppressions = collectFileSuppressions(script)
	trav := NewNodeTraversal(c, c.creator, c.global)
	trav.Traverse(script)
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
	}
	return true
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
	case rhino.CALL:
		c.emit(c.visitCall(n, scope))
	case rhino.NEW:
		c.emit(c.visitNew(n, scope))
	case rhino.GETPROP:
		c.emit(c.visitGetProp(n, scope))
	case rhino.GETELEM:
		c.emit(c.visitGetElem(n, scope))
	case rhino.CLASS_FIELD_DEF:
		c.emit(c.visitClassField(n, parent, scope))
	case rhino.RETURN:
		c.emit(c.visitReturn(n, scope))
	case rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE:
		c.emit(c.visitEquality(n, scope))
	case rhino.IN:
		c.emit(c.visitIn(n, scope))
	case rhino.FOR_IN:
		c.emit(c.visitForIn(n, scope))
	case rhino.LSH, rhino.RSH, rhino.URSH, rhino.BITAND, rhino.BITOR, rhino.BITXOR:
		c.emit(c.visitBitOp(n, scope))
	case rhino.FUNCTION:
		c.emit(c.visitFunctionParams(n, scope))
	}
	// Strict-mode operand checks — OFF by default, matching Java.
	// These fire only when `--jscomp_warning=strictCheckTypes` is set.
	if c.opts.IsOn(GroupStrictCheckTypes) {
		switch n.Token() {
		case rhino.ADD:
			c.emit(c.visitAdd(n, scope))
		case rhino.SUB, rhino.MUL, rhino.DIV, rhino.MOD:
			c.emit(c.visitArithmetic(n, scope))
		case rhino.BITAND, rhino.BITOR, rhino.BITXOR,
			rhino.LSH, rhino.RSH, rhino.URSH:
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
	name := v.FirstChild()
	if name == nil || name.Token() != rhino.NAME {
		return nil
	}
	var out []diag.Diagnostic
	init := name.FirstChild()
	if init == nil {
		return out
	}
	// `/** @enum */ var E = {FOO: 1, ...}` — the VAR's declared
	// type is the EnumType, but the init is an OBJECTLIT that
	// infers to a structural {...}. Skip the normal subtype
	// check (the RHS is the enum's *source of truth*, not a
	// value being assigned into it) and run the enum-specific
	// diagnostics below instead. Mirrors how Java's TypeCheck
	// routes `@enum` initializers through EnumType's own
	// validateEnumInit helper rather than the assignment rule.
	if info := v.JSDocInfo(); info.IsEnum() {
		return append(out, c.visitEnumInit(name, init, scope)...)
	}
	declared := scope.GetDeclaredType(name.String_())
	if declared == nil {
		return out
	}
	infer := NewTypeInference(scope)
	found := infer.InferType(init)
	if jstype.IsSubtype(found, declared) {
		return out
	}
	msg := fmt.Sprintf("initializing variable\nfound   : %s\nrequired: %s",
		found.TypeName(), declared.TypeName())
	line, col := diagPosForArg(init)
	return append(out, diag.Diagnostic{
		Path: c.path, Line: line, Col: col,
		Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
	})
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

// isDefinitelyNonConstructor reports whether `t` is guaranteed not
// to be a constructor. Returns the offending type so the diagnostic
// can show what was found. A FunctionType without IsConstructor
// counts (plain functions); non-function primitives count; Unknown
// / All / unions stay quiet.
func isDefinitelyNonConstructor(t jstype.JSType) (bool, jstype.JSType) {
	if t == nil {
		return false, nil
	}
	if isUnknownOrNil(t) {
		return false, nil
	}
	if _, ok := t.(*jstype.UnionType); ok {
		return false, nil
	}
	if fn, ok := t.(*jstype.FunctionType); ok {
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
	if node.Token() == rhino.CALL && len(fn.Params) > 0 && callee.Token() == rhino.NAME {
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
	// set fires JSC_INEXISTENT_ENUM_ELEMENT. Runs before the
	// ObjectType branch because an EnumType isn't an ObjectType.
	if enum, ok := recvType.(*jstype.EnumType); ok {
		if enum.Members[propName] {
			return nil
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
	// Broader INEXISTENT_PROPERTY: fire on structural (@record /
	// @interface) receivers, on @struct nominal classes (sealed
	// property surface by design), and on the known-surface
	// builtin types we ship concrete slot maps for. For plain
	// @constructor types without @struct, Java still fires on
	// `this.prop` READS inside a prototype method — the
	// missing-this-slot case — but allows writes on instance
	// values outside the constructor to implicitly extend the
	// surface.
	if !obj.IsStructural && !obj.IsStruct && !isKnownSurfaceBuiltin(obj.Name) {
		if !isThisReceiverRead(n) {
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
	if info == nil || info.ReturnType() == nil {
		return nil
	}
	declared := scope.ResolveExpr(info.ReturnType().TypeName)
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
		return []diag.Diagnostic{{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
			Message: "only arrays or objects can be accessed\nfound   : " +
				rawRecv.TypeName() + "\nrequired: Object",
		}}
	}
	recvType := stripNullUndef(rawRecv)
	if isStructReceiver(recvType) {
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

// visitEquality ports the JSC_DETERMINISTIC_TEST branch of
// TypeCheck.visitBinaryOp: when the two operands of `==` / `!=`
// can never compare equal — disjoint primitive slot types —
// the comparison always yields a known value and Java emits at
// the whole expression's column. `===` / `!==` use a separate
// code (not yet ported).
// visitFunctionParams fires JSC_INEXISTENT_PARAM when a @param
// JSDoc names a parameter that doesn't appear in the function's
// actual parameter list. Java reports at the FUNCTION keyword's
// column.
func (c *TypeCheck) visitFunctionParams(fn *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	info := fn.JSDocInfo()
	if info == nil {
		return nil
	}
	names := info.ParameterNames()
	if len(names) == 0 {
		return nil
	}
	kids := fn.Children()
	if len(kids) < 2 {
		return nil
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
	var out []diag.Diagnostic
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

// visitIn fires JSC_IN_USED_WITH_STRUCT on `key in struct` — Java's
// TypeCheck rejects the `in` operator on @struct receivers because
// the surface is sealed and membership is a compile-time concept.
func (c *TypeCheck) visitIn(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	rhs := kids[1]
	infer := NewTypeInference(scope)
	recv := stripNullUndef(infer.InferType(rhs))
	if !isStructReceiver(recv) {
		return nil
	}
	return []diag.Diagnostic{{
		Path: c.path, Line: rhs.Lineno(), Col: rhs.Charno(),
		Code: diag.JSC_IN_USED_WITH_STRUCT, Level: diag.Error,
		Message: "Cannot use the IN operator with structs",
	}}
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

func (c *TypeCheck) visitEquality(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	infer := NewTypeInference(scope)
	l := infer.InferType(kids[0])
	r := infer.InferType(kids[1])
	det, equalValue := disjointForEquality(l, r)
	if !det {
		return nil
	}
	// Compute the final boolean for this operator. `equalValue`
	// is the result under `==` / `===`; `!=` / `!==` flip it.
	value := equalValue
	if n.Token() == rhino.NE || n.Token() == rhino.SHNE {
		value = !value
	}
	valueStr := "false"
	if value {
		valueStr = "true"
	}
	return []diag.Diagnostic{{
		Path: c.path, Line: n.Lineno(), Col: n.Charno(),
		Code: diag.JSC_DETERMINISTIC_TEST, Level: diag.Error,
		Message: fmt.Sprintf(
			"condition always evaluates to %s\nleft : %s\nright: %s",
			valueStr, l.TypeName(), r.TypeName()),
	}}
}

// disjointForEquality reports whether types `a` and `b` always
// produce the same result under `==` / `===` / `!=` / `!==`. JS's
// `==` coerces across number / string / boolean; function ↔
// string / number, symbol ↔ anything-else etc. are always-false.
// Single-inhabitant types (null, undefined) compared against
// themselves always succeed. Conservative: stay silent on
// Unknown / All / unions that might overlap.
//
// Returns (deterministic, valueIfEqualOperator):
//   - `deterministic` true when the pair is disjoint OR both
//     sides are the same single-inhabitant type.
//   - `valueIfEqualOperator` is the comparison result assuming
//     the operator is `==` or `===` (`true` for same-singleton,
//     `false` for disjoint). Callers flip for `!=` / `!==`.
func disjointForEquality(a, b jstype.JSType) (bool, bool) {
	if a == nil || b == nil {
		return false, false
	}
	if isUnknownOrNil(a) || isUnknownOrNil(b) {
		return false, false
	}
	if _, ok := a.(*jstype.UnionType); ok {
		return false, false
	}
	if _, ok := b.(*jstype.UnionType); ok {
		return false, false
	}
	ka := classifyEquality(a)
	kb := classifyEquality(b)
	if ka == eqUnknown || kb == eqUnknown {
		return false, false
	}
	if ka == kb {
		// Same single-inhabitant kind → always equal.
		if ka == eqNull || ka == eqVoid {
			return true, true
		}
		return false, false
	}
	// `null == undefined` is true; keep them quiet.
	if (ka == eqNull && kb == eqVoid) || (ka == eqVoid && kb == eqNull) {
		return false, false
	}
	// number / string / boolean coerce to each other under `==`.
	if isPrimitiveKind(ka) && isPrimitiveKind(kb) {
		return false, false
	}
	// Everything else — function/object vs primitive, or distinct
	// function-vs-object — is always false under `==`.
	return true, false
}

func isPrimitiveKind(k equalityKind) bool {
	switch k {
	case eqNumber, eqString, eqBoolean:
		return true
	}
	return false
}

type equalityKind int

const (
	eqUnknown equalityKind = iota
	eqNumber
	eqString
	eqBoolean
	eqNull
	eqVoid
	eqFunction
	eqObject
)

func classifyEquality(t jstype.JSType) equalityKind {
	switch t {
	case jstype.Number:
		return eqNumber
	case jstype.String:
		return eqString
	case jstype.Boolean:
		return eqBoolean
	case jstype.Null:
		return eqNull
	case jstype.Void:
		return eqVoid
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return eqFunction
	}
	if obj, ok := t.(*jstype.ObjectType); ok {
		// Nominal object types compare with each other, not with
		// the primitive slots. Includes externs-declared Function /
		// Date / Array / etc. and user-declared @constructor
		// instances.
		return eqObjectKind(obj)
	}
	if tt, ok := t.(*jstype.TemplatizedType); ok && tt.Base != nil {
		return eqObjectKind(tt.Base)
	}
	return eqUnknown
}

// eqObjectKind collapses every distinct nominal ObjectType to a
// single `eqObject` bucket for coarse disjointness checks — good
// enough for `number == Date` / `Function == string`-style
// always-false comparisons without having to enumerate every
// extern.
func eqObjectKind(_ *jstype.ObjectType) equalityKind { return eqObject }

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
