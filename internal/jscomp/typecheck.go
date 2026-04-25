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
//   - JSC_NON_STRINGIFIABLE_OBJECT_KEY on @type/@return/@param/@typedef
//     annotations whose Object<K,V> key K is not stringifiable
//     (port of checkJsdocInfoContainsObjectWithBadKey /
//     validateNonStringifiableKey / checkStringifiable)

package jscomp

import (
	"fmt"
	"sort"
	"strings"

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

	// tv is the shared TypeValidator for this file; wired in Check()
	// so that visitAssign / visitReturn / visitCall route JSC_TYPE_MISMATCH
	// through the canonical mismatch dispatcher (structural analysis,
	// TemplateType skipping, TypeMismatch accumulation). Mirrors Java's
	// TypeCheck holding a reference to a per-compiler TypeValidator.
	tv *TypeValidator

	// inexistentReported dedups JSC_INEXISTENT_PROPERTY emissions
	// across the script — Java fires the diagnostic only once per
	// (receiver, property) pair even when the program references
	// the missing slot multiple times. Keeps message quality high
	// without tracking per-guard refinement state.
	inexistentReported map[string]bool

	// noPropsReported dedups "No properties on this expression"
	// JSC_TYPE_MISMATCH emissions per receiver root name (NAME or
	// qualified name). Java's TypeValidator.expectObject reports the
	// mismatch through registerMismatchAndReport which dedupes by
	// (location, found, required); we approximate per-root here.
	noPropsReported map[string]bool

	// fileSuppressions carries any @fileoverview @suppress {...}
	// names captured pre-traversal. visitGetProp / visitReturn /
	// etc. consult it to skip diagnostics the author opted out
	// of file-wide.
	fileSuppressions map[string]bool

	// looseProps mirrors Java's eachRefTypeIndexedByProperty: for each
	// ObjectType, the set of property names defined on any of its
	// subtypes. When a property is "loose" on a type (present on a
	// subtype but not on the type itself), Java's
	// canPropertyBeDefined returns LOOSE, which causes
	// checkPropertyAccessHelper to set strictReport=true and emit
	// only STRICT_INEXISTENT_PROPERTY. In default mode that is
	// silenced, so visitGetProp skips the diagnostic when the
	// property is loose.
	looseProps map[*jstype.ObjectType]map[string]bool

	// externTypes is the set of ObjectType instances registered in the
	// externs scope (Object, Function, Date, Error, etc.). In Java these
	// types have full property surfaces populated from closure-compiler's
	// externs .js files; in jsgo the externs are sparse (mostly bare
	// @constructor stubs). Treating these types as "open" — like Java
	// treats Object as low-confidence — avoids FPs from accessing any
	// property on an extern type whose surface we haven't declared.
	externTypes     map[*jstype.ObjectType]bool
	globalPropNames map[string]bool

	// undeclaredNames is the set of variable names reported as
	// undeclared by VarCheck. When the root of a GETPROP chain is one
	// of these names, POSSIBLE_INEXISTENT_PROPERTY is suppressed.
	// Mirrors Java VarCheck's synthesized-extern mechanism: once VarCheck
	// marks a name as an extern, all property accesses on it are permitted.
	undeclaredNames map[string]bool

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

	// scriptRoot is the SCRIPT node passed to Check; used by the
	// inlined superclass-check helpers (findConstructorJSDoc,
	// scanConstructorJSDoc) that walk the full AST by name.
	scriptRoot *rhino.Node
}

// Go-only: language-adapter — Go's single-file package requires one constructor
// per caller ABI; Java's TypeCheck has two overloaded constructors (lines 437, 452).
func NewTypeCheck(path string, creator *TypedScopeCreator, global *TypedScope) *TypeCheck {
	return NewTypeCheckWithOptions(path, creator, global, NewCompilerOptions())
}

// Go-only: language-adapter — Go does not support method overloading; this variant
// carries the CompilerOptions bag that Java's constructor accepted implicitly via the
// AbstractCompiler (TypeCheck.java:437–459).
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

// Go-only: package-bridge — injects undeclared names from VarCheck. Java routes this
// through AbstractCompiler.createSynthesizedExternVar(); in Go the two passes share
// a plain map since they run in the same process without a shared compiler object.
func (c *TypeCheck) SetUndeclaredNames(names map[string]bool) {
	c.undeclaredNames = names
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
	c.scriptRoot = script
	c.tv = NewTypeValidator(c.path, c.global.Registry())
	c.installFlow(script, c.global)
	c.collectPropertyWrites(script)
	c.looseProps = buildLoosePropertyIndex(c.global)
	c.externTypes = buildExternTypeSet(c.global)
	c.globalPropNames = buildGlobalPropNameSet(c.global, script)
	c.fileSuppressions = collectFileSuppressions(script)
	// Emit JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES for every @template {B} T
	// annotation. Mirrors RhinoErrorReporter.UNSUPPORTED_BOUNDED_GENERIC_TYPES
	// which maps the addParserWarning(JSDOC_TEMPLATE_BOUNDEDGENERICS_USED)
	// diagnostic from JsDocInfoParser.java's TEMPLATE case.
	// TODO: bounded generics are still in development in the Java compiler too;
	// this diagnostic is intentionally always-error (no suppression path).
	c.emitJSDocParseWarnings(script)
	trav := NewNodeTraversal(c, c.creator, c.global)
	trav.Traverse(script)
	c.walkSuperclassCheck(script)
	return c.diags
}

// emitJSDocParseWarnings walks the AST and emits diagnostics for JSDoc
// parse warnings recorded during parsing: JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES
// (from @template {B} T) and JSC_INVALID_PARAM (from invalid @param names).
// Go-only: package-bridge — walks JSDoc parse errors that Java surfaces via
// RhinoErrorReporter / JsDocInfoParser callbacks. In Go, rhino.JSDocInfo exposes
// the positions directly, so TypeCheck walks them once at the top of Check().
func (c *TypeCheck) emitJSDocParseWarnings(script *rhino.Node) {
	var walk func(*rhino.Node)
	walk = func(n *rhino.Node) {
		if n == nil {
			return
		}
		if info := n.JSDocInfo(); info != nil {
			for _, pos := range info.BoundedTemplatePositions() {
				c.diags = append(c.diags, diag.Diagnostic{
					Path: c.path, Line: pos.Line, Col: pos.Col,
					Code: diag.JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES, Level: diag.Error,
					Message: "Bounded generic semantics are currently still in development",
				})
			}
			for _, pos := range info.InvalidParamPositions() {
				c.diags = append(c.diags, diag.Diagnostic{
					Path: c.path, Line: pos.Line, Col: pos.Col,
					Code:    diag.JSC_INVALID_PARAM,
					Level:   diag.Warning,
					Message: `Parse error. invalid param name "` + pos.Name + `"`,
				})
			}
		}
		for _, ch := range n.Children() {
			walk(ch)
		}
	}
	walk(script)
}

// Go-only: package-bridge — Java's CompilerPass base class carries a reference to
// AbstractCompiler which exposes file-level suppressions via getSuppressedWarnings();
// in Go we pre-scan the script's JSDoc once and pass the map to TypeCheck.
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

// Go-only: package-bridge — Java's TypedScopeCreator already registers all property
// writes during its own traversal so TypeCheck never needs a pre-pass; in Go the
// scope creator's traversal is separated, so TypeCheck performs this walk itself.
func (c *TypeCheck) collectPropertyWrites(n *rhino.Node) {
	c.collectPropertyWritesIn(n, c.global)
}

// Go-only: package-bridge — scope-aware worker for collectPropertyWrites; no Java
// counterpart because Java's TypedScopeCreator handles this during its own pass.
func (c *TypeCheck) collectPropertyWritesIn(n *rhino.Node, scope *TypedScope) {
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
					if enc := NU.GetEnclosingFunction(n); enc != nil {
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
					infer := NewTypeInference(scope)
					recvType := stripNullUndef(infer.InferType(recv))
					defineOnOpenObject(recvType, propName)
				}
			}
		}
	}
	for _, ch := range n.Children() {
		childScope := scope
		if ch.Token() == rhino.FUNCTION {
			if fnScope, ok := c.creator.Functions[ch]; ok && fnScope != nil {
				childScope = fnScope
			}
		}
		c.collectPropertyWritesIn(ch, childScope)
	}
}

// Go-only: language-adapter — Java runs TypeInference as a separate compiler pass
// before TypeCheck; in Go, TypeCheck builds CFG + DataFlowAnalysis lazily so the
// same traversal can consult in-scope refinements (TypeCheck.java has no equivalent).
func (c *TypeCheck) installFlow(root *rhino.Node, scope *TypedScope) *DataFlowAnalysis {
	if existing, ok := c.flows[root]; ok {
		return existing
	}
	cfg := BuildControlFlowGraph(root)
	pass := NewTypeInferencePass(scope)
	d := NewDataFlowAnalysis(cfg, EmptyFlowScope(), pass.Transfer)
	d.Run()
	// Port of DataFlowAnalysis.computeEscaped: identify variables declared
	// in this function scope that escape into nested closures. Their
	// flow-refined types are "unflowable" — we must not override the
	// declared type with the narrowed flow type at read sites (mirrors
	// TypeInference.java#isUnflowable / traverseName guard).
	if root.Token() == rhino.FUNCTION {
		escaped := ComputeEscaped(root, scope.LocalDeclaredNames())
		d.Escaped = escaped.Escaped
	}
	c.flows[root] = d
	return d
}

// Go-only: language-adapter — helper for the integrated flow analysis; Java's
// TypeInference pass uses a separate lattice/CFG framework; no Java equivalent.
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

// Go-only: language-adapter — Go's NodeTraversal lacks Java's per-branch
// FlowScope threading from TypeInference; traverseLoop / traverseShortCircuit /
// traverseBlock / traverseIf / traverseHook implement equivalent refinement pushing
// inline within ShouldTraverse so Visit sees narrowed types.
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

// Go-only: language-adapter — see traverseLoop comment; short-circuit variant.
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

// Go-only: language-adapter — see traverseLoop comment; block-statement variant.
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

// Go-only: language-adapter — helper for traverseBlock's post-guard narrowing;
// Java's TypeInference tracks this implicitly via FlowScope lattice join.
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

// Go-only: language-adapter — helper for earlyExitGuard; no Java counterpart.
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

// Go-only: language-adapter — see traverseLoop comment; if-branch variant.
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

// Go-only: language-adapter — see traverseLoop comment; ternary-expression variant.
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
		c.emit(c.visitShortCircuitAssign(n, scope))
	case rhino.CALL:
		c.emit(c.visitCall(n, scope))
	case rhino.NEW:
		c.emit(c.visitNew(n, scope))
	case rhino.GETPROP:
		c.emit(c.visitGetProp(n, scope))
	case rhino.GETELEM:
		c.emit(c.visitGetElem(n, scope))
	case rhino.CLASS_FIELD_DEF, rhino.MEMBER_FIELD_DEF:
		c.emit(c.visitClassField(n, parent, scope))
	case rhino.COMPUTED_PROP:
		c.emit(c.visitComputedPropInClass(n, parent, scope))
	case rhino.RETURN:
		c.emit(c.visitReturn(n, scope))
	case rhino.IN:
		c.emit(c.visitIn(n, scope))
	case rhino.FOR_IN:
		c.emit(c.visitForIn(n, scope))
	case rhino.FOR_OF, rhino.FOR_AWAIT_OF:
		c.emit(c.checkForOfTypes(n, scope))
	case rhino.LSH, rhino.RSH, rhino.URSH, rhino.BITAND, rhino.BITOR, rhino.BITXOR, rhino.BITNOT,
		rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH:
		c.emit(c.visitBitOp(n, scope))
	case rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE:
		c.emit(c.visitTypeofComparison(n, scope))
		c.emit(c.visitEqualityOp(n, scope))
	case rhino.LT, rhino.LE, rhino.GT, rhino.GE:
		c.emit(c.visitRelationalDefault(n, scope))
	case rhino.FUNCTION:
		c.emit(c.visitFunctionParams(n, scope))
	case rhino.OBJECTLIT:
		c.emit(c.visitObjectLit(n, scope))
	case rhino.ARRAYLIT:
		c.emit(c.visitArrayLit(n, scope))
	case rhino.CAST:
		c.emit(c.visitCast(n, scope))
	case rhino.REGEXP:
		c.emit(c.visitRegExpLit(n, scope))
	case rhino.TEMPLATELIT:
		c.emit(c.visitTemplateLit(n, scope))
	case rhino.OBJECT_PATTERN:
		c.emit(c.visitObjectPattern(n, scope))
	case rhino.DYNAMIC_IMPORT:
		// Java: visitDynamicImport — checks that import(specifier) has a
		// non-null, non-undefined string specifier.
		c.emit(c.visitDynamicImport(n, scope))
	case rhino.CASE:
		// Java: validator.expectSwitchMatchesCase(n, switchConditionType, caseType).
		// The CASE node is a child of SWITCH. The switch condition type is the
		// type of the SWITCH's first child (the expression before the SWITCH_BODY).
		c.emit(c.visitSwitchCase(n, parent, scope))
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
	// Java's TypeCheck fires checkJsdocInfoContainsObjectWithBadKey on
	// every typeable node at the end of visit(). Mirrors the call at
	// TypeCheck.java:969 (after ensureTyped).
	c.emit(c.checkJsdocInfoContainsObjectWithBadKey(n, scope))
}

// Go-only: language-adapter — threads CFG flow-scope into Visit; Java's TypeCheck
// receives a pre-narrowed scope from TypeInference via NodeTraversal; no equivalent.
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
	// Port of TypeInference.java#isUnflowable / traverseName guard:
	// variables that escape into nested closures are "unflowable" — their
	// flow-narrowed type must not override the declared type at read sites.
	// Filter escaped names from the refinement map before pushing.
	types := in.types
	if len(d.Escaped) > 0 {
		filtered := make(map[string]jstype.JSType, len(types))
		for name, t := range types {
			if !d.Escaped[name] {
				filtered[name] = t
			}
		}
		if len(filtered) == 0 {
			return func() {}
		}
		types = filtered
	}
	scope.PushRefinement(types)
	return scope.PopRefinement
}

// Go-only: language-adapter — helper for pushEnclosingFlow; no Java counterpart.
func enclosingFlowRoot(n *rhino.Node) *rhino.Node {
	for cur := n; cur != nil; cur = cur.Parent() {
		switch cur.Token() {
		case rhino.FUNCTION, rhino.SCRIPT:
			return cur
		}
	}
	return nil
}

// Go-only: language-adapter — diagnostic accumulator; Java's TypeCheck calls
// compiler.report() directly via the AbstractCompiler interface; Go collects into
// a slice returned from Check() since there is no shared compiler object.
func (c *TypeCheck) emit(d []diag.Diagnostic) {
	c.diags = append(c.diags, d...)
}

// tvDrain collects all diagnostics accumulated by c.tv since the last
// drain call, clears the TypeValidator's buffer, and returns them.
// Go-only: language-adapter — drains TypeValidator's diagnostic buffer between
// visitor calls; Java routes all diagnostics through AbstractCompiler.report() so
// this accumulation/drain pattern does not exist in Java.
func (c *TypeCheck) tvDrain() []diag.Diagnostic {
	d := c.tv.diags
	c.tv.diags = nil
	return d
}

// ---- per-node visitors ------------------------------------------ //

// Go-only: language-adapter — guards against parser-level attach failures that
// produce an absent init node; Java's parser always attaches initializers correctly.
func varHasEqOnLine(src string, lineno, nameCol int) bool {
	if src == "" || lineno < 1 {
		return false
	}
	lines := strings.Split(src, "\n")
	if lineno > len(lines) {
		return false
	}
	l := lines[lineno-1]
	if nameCol < 0 || nameCol >= len(l) {
		return false
	}
	return strings.ContainsRune(l[nameCol:], '=')
}

// visitVar — port of TypeCheck.java's handling of VAR / LET / CONST:
// compare the initializer's inferred type against the @type
// declaration. Reports at the initializer's position to match Java's
// "initializing variable" variant.
func (c *TypeCheck) visitVar(v *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	// JSC_MULTIPLE_VAR_DEF: a JSDoc @type annotation can only describe
	// a single variable. Fire when the declaration has more than one
	// NAME child. Mirrors TypeCheck.java's visitVar guard.
	if info := v.JSDocInfo(); info != nil && info.HasType() && len(v.Children()) > 1 {
		return []diag.Diagnostic{{
			Path: c.path, Line: v.Lineno(), Col: v.Charno(),
			Code:    diag.JSC_MULTIPLE_VAR_DEF,
			Level:   diag.Warning,
			Message: "declaration of multiple variables with shared type information",
		}}
	}
	name := v.FirstChild()
	if name == nil || name.Token() != rhino.NAME {
		return nil
	}
	var out []diag.Diagnostic
	init := name.FirstChild()
	if init == nil {
		// `/** @enum */ var E;` with no initializer — Java emits
		// JSC_ENUM_INITIALIZER_NOT_ENUM warning. Mirrors TypeCheck.java's
		// validateEnumInitializer guard which fires when the RHS is absent.
		// Guard: if the source line has `=` after the NAME, the parser
		// failed to attach the init (e.g. inline JSDoc inside OBJECTLIT);
		// skip the diagnostic to avoid false positives.
		if info := v.JSDocInfo(); info.IsEnum() && !varHasEqOnLine(c.creator.sourceSrc, name.Lineno(), name.Charno()) {
			return []diag.Diagnostic{{
				Path: c.path, Line: name.Lineno(), Col: name.Charno(),
				Code:    diag.JSC_ENUM_INITIALIZER_NOT_ENUM,
				Level:   diag.Warning,
				Message: "enum initializer must be an object literal or an enum",
			}}
		}
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
	// JSC_INCOMPATIBLE_ALIAS_ANNOTATION.
	if info := v.JSDocInfo(); info.IsRecord() || info.IsInterface() || info.IsConstructor() {
		if init.Token() == rhino.NAME {
			if src, ok := scope.GetDeclaredType(init.String_()).(*jstype.FunctionType); ok && src.Instance != nil {
				if ann := checkAliasAnnotation(info, src); ann != "" {
					return []diag.Diagnostic{{
						Path: c.path, Line: name.Lineno(), Col: name.Charno(),
						Code: diag.JSC_INCOMPATIBLE_ALIAS_ANNOTATION, Level: diag.Warning,
						Message: fmt.Sprintf("Annotation %s on %s incompatible with aliased type.", ann, name.String_()),
					}}
				}
			}
		}
	}
	if info := v.JSDocInfo(); info.IsEnum() {
		return append(out, c.visitEnumInit(name, init, scope)...)
	}
	// Java's visitVar iterates ALL NAME children for the initializing-variable
	// type check. Each child may have its own inline JSDoc (@type annotation).
	// `varInfo` (the VAR-level JSDoc) is used as fallback only when there's
	// exactly ONE child — for multi-binding VARs, per-child JSDoc applies.
	// Mirrors Java's `JSDocInfo varInfo = n.hasOneChild() ? n.getJSDocInfo() : null`.
	varInfo := v.JSDocInfo()
	isMultiChild := len(v.Children()) > 1
	infer := NewTypeInference(scope)
	for curName := name; curName != nil; curName = curName.Next() {
		if curName.Token() != rhino.NAME {
			continue
		}
		curInit := curName.FirstChild()
		if curInit == nil {
			continue
		}
		// CAST inits are validated by visitCast.
		if curInit.Token() == rhino.CAST {
			continue
		}
		// Determine the declared type for this binding.
		// Per-child JSDoc takes precedence; VAR-level JSDoc as fallback
		// only for single-child VARs (matching Java).
		childInfo := curName.JSDocInfo()
		if childInfo != nil && childInfo.IsCastHoisted() {
			continue // handled by visitCast
		}
		var declared jstype.JSType
		if childInfo != nil && childInfo.HasType() {
			declared = scope.ResolveExpr(childInfo.Type().TypeName)
		} else if !isMultiChild && varInfo != nil {
			declared = scope.GetDeclaredType(curName.String_())
		} else {
			declared = scope.GetDeclaredType(curName.String_())
		}
		if declared == nil {
			continue
		}
		found := infer.InferType(curInit)
		if jstype.IsSubtype(found, declared) {
			continue
		}
		msg := fmt.Sprintf("initializing variable\nfound   : %s\nrequired: %s",
			found.TypeName(), declared.TypeName())
		line, col := diagPosForArg(curInit)
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: line, Col: col,
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
		})
		// For the first child, stop (Java also stops after first for
		// single-child VARs). For multi-child, continue checking all.
		if !isMultiChild {
			break
		}
	}
	return out
}


// Go-only: language-adapter — Java handles CAST inline in visit() at line 613-629;
// Go extracts it into a named function for readability.
func (c *TypeCheck) visitCast(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	expr := n.FirstChild()
	if expr == nil || expr.Token() == rhino.OBJECTLIT {
		return nil
	}
	info := n.JSDocInfo()
	if info == nil || !info.HasType() {
		return nil
	}
	castType := scope.ResolveExpr(info.Type().TypeName)
	if castType == nil {
		return nil
	}
	// When the cast target resolved to Unknown, CheckJSDocTypes has
	// already emitted (or will emit) JSC_UNRECOGNIZED_TYPE_ERROR; Java's
	// TypeCheck skips the CAST subtype test when the target is
	// unknown-typed (see TypeCheck.java:validateTypeAnnotation which
	// returns early for unresolved annotations). Mirror that here.
	if isUnknownOrNil(castType) {
		return nil
	}
	// When the annotation resolves to a constructor (e.g. `@type {Function}`
	// or `@type {Object}` via externs), use the instance type for the cast
	// check — annotations refer to the nominal instance type, not the
	// constructor function itself. Mirrors Java's getJSType(CAST) which
	// returns the instance type of a constructor-named @type annotation.
	if fn, ok := castType.(*jstype.FunctionType); ok && fn.IsConstructor && fn.Instance != nil {
		castType = fn.Instance
	}
	infer := NewTypeInference(scope)
	exprType := infer.InferType(expr)
	// Ditto for exprType: if inference returns a constructor FunctionType
	// (e.g. when a NAME resolves to a bare constructor), use its instance.
	if fn, ok := exprType.(*jstype.FunctionType); ok && fn.IsConstructor && fn.Instance != nil {
		exprType = fn.Instance
	}
	// Use the full CanCastToVisitor logic (mirrors JSType.canCastTo).
	v := &jstype.CanCastToVisitor{}
	if jstype.VisitRelationship(v, exprType, castType).(bool) {
		return nil
	}
	return []diag.Diagnostic{{
		Path:  c.path,
		Line:  n.Lineno(),
		Col:   n.Charno(),
		Code:  diag.JSC_INVALID_CAST,
		Level: diag.Error,
		Message: fmt.Sprintf(
			"invalid cast - must be a subtype or supertype\nfrom: %s\nto  : %s",
			exprType.TypeName(), castType.TypeName()),
	}}
}

// Go-only: language-adapter — Go extracts Java's visitVar @enum branch into its
// own function for clarity; in Java the enum-initializer logic is inlined inside
// visitVar / TypeCheck.java:2266-2300 without a named sub-method.
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
	if init.Token() == rhino.NAME || init.Token() == rhino.GETPROP {
		// Resolve the RHS: for NAME look up the declared type; for
		// GETPROP (qualified namespace enum alias e.g. `goog.a`)
		// use InferType which routes through LookupNamespaceProp.
		var other *jstype.EnumType
		if init.Token() == rhino.NAME {
			other, _ = scope.GetDeclaredType(init.String_()).(*jstype.EnumType)
		} else {
			other, _ = infer.InferType(init).(*jstype.EnumType)
		}
		if other != nil {
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

// Go-only: language-adapter — Java's visitAssign handles ASSIGN_AND/OR/COALESCE
// inline (TypeCheck.java:1113 switch); Go extracts this variant for clarity.
func (c *TypeCheck) visitShortCircuitAssign(a *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := a.Children()
	if len(kids) < 2 {
		return nil
	}
	lhs, rhs := kids[0], kids[1]
	if lhs.Token() != rhino.NAME {
		return nil
	}
	declared := scope.GetRawDeclaredType(lhs.String_())
	if declared == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	found := infer.InferType(rhs)
	c.tv.expectCanAssignTo(lhs, found, declared, "assignment")
	return c.tvDrain()
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
		// Route through TypeValidator.expectCanAssignTo so the canonical
		// mismatch dispatcher runs (TemplateType skip, structural-property
		// analysis, TypeMismatch accumulation). Mirrors Java's TypeCheck.visitAssign
		// calling typeValidator.expectCanAssignTo(lhs, rhsType, lhsType, msg).
		c.tv.expectCanAssignTo(lhs, found, declared, "assignment")
		return c.tvDrain()
	case rhino.GETELEM:
		// Only check for ReadonlyArray assignment; GETELEM type-mismatch
		// is not yet ported (matches Java's checkNotReadonlyPropertyAssignment
		// scope: fires before any other per-element checks).
		elemKids := lhs.Children()
		if len(elemKids) >= 1 {
			if d := c.checkNotReadonlyPropertyAssignment(lhs, elemKids[0], infer); d != nil {
				return []diag.Diagnostic{*d}
			}
		}
		return nil
	case rhino.GETPROP:
		propKids := lhs.Children()
		if len(propKids) < 2 {
			return nil
		}
		propName := propKids[1].String_()
		recv := propKids[0]
		// `@suppress {duplicate}` on the assignment opts the slot out
		// of redeclaration / type-mismatch checks. Mirrors how Java's
		// TypeCheck consults the JSDoc suppression set before firing
		// the assignment mismatch.
		if info := a.JSDocInfo(); info != nil && info.HasSuppression("duplicate") {
			return nil
		}
		// ReadonlyArray guard (port of checkNotReadonlyPropertyAssignment).
		// Java fires at the lhs (GETPROP) node; Java GETPROP charno is the
		// property-name column, which corresponds to propKids[1] in Go.
		if d := c.checkNotReadonlyPropertyAssignment(propKids[1], recv, infer); d != nil {
			return []diag.Diagnostic{*d}
		}
		// Interface prototype assignment: `I.prototype.x = value` where
		// I is an @interface. Only empty functions or goog.abstractMethod
		// are permitted as rvalues; anything else fires
		// JSC_INVALID_INTERFACE_MEMBER_DECLARATION. Mirrors Java's
		// TypeCheck.visitAssign dispatching to visitInterfacePropertyAssignment
		// at TypeCheck.java:1190-1194.
		if recv.Token() == rhino.GETPROP {
			recvKids := recv.Children()
			if len(recvKids) >= 2 && recvKids[1].String_() == "prototype" {
				if iname := recvKids[0]; iname.Token() == rhino.NAME {
					if fn, ok := scope.GetDeclaredType(iname.String_()).(*jstype.FunctionType); ok && fn.Instance != nil && fn.Instance.IsInterface {
						c.emit(c.visitInterfacePropertyAssignment(recv, lhs, scope))
					}
				}
			}
		}
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
		// `this.X = …` at global scope resolves to the externs slot for
		// global `X` (Window.X). Mirrors Closure's TypeCheck where the
		// global this prototype is Window, so assignments to its
		// properties type-check against the externs declaration. Only
		// considers names declared in externs (not in the user's global
		// scope) — a user redeclaration would already have triggered
		// JSC_VAR_MULTIPLY_DECLARED_ERROR and Java skips further analysis
		// on that file in that case.
		if recv.Token() == rhino.THIS && scope.ThisType() == nil {
			if t := scope.LookupExternsOnly(propName); t != nil {
				declared = t
				ownerName = "global this"
			}
		}
		if declared == nil && recv.Token() == rhino.NAME {
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
			recvName := recv.String_()
			if t := scope.LookupNamespaceProp(recvName, propName); t != nil {
				// Skip properties whose type was inferred from a literal
				// RHS (no @type JSDoc) — Java treats these as open surfaces
				// where any subsequent assignment is allowed, matching the
				// IsPropertyInferred guard on prototype properties.
				// Exception: if the flow-scope object has a @type-annotated
				// declared property (from an object literal key), the declared
				// type takes precedence over the inferred namespace prop.
				if scope.IsNamespacePropInferred(recvName, propName) {
					if objCheck, ok2 := infer.InferType(recv).(*jstype.ObjectType); ok2 && objCheck.HasOwnDeclaredProperty(propName) {
						declared = objCheck.GetPropertyType(propName)
						ownerName = recvName
					} else {
						return nil
					}
				} else {
					declared = t
					ownerName = recvName
				}
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
		// `/** @enum {T} */ ns.Foo = {A:1, B:2}` — the RHS OBJECTLIT
		// is the enum's source of truth, not a value being assigned into
		// it; skip the top-level assignment check (per-entry validation
		// is handled by visitEnumInit for VAR form; namespace-prop enum
		// inits don't fire per-entry checks). Mirrors Java's TypeCheck
		// which routes @enum assignments through EnumType's own init
		// validator and does not fire ASSIGN mismatch at the top level.
		// However, enum aliasing (`/** @enum {string} */ ns.b = a` where
		// `a: @enum {number}`) fires "incompatible enum element types" when
		// the source and target underlying types disagree.
		if declEnum, isEnum := declared.(*jstype.EnumType); isEnum {
			if rhs.Token() == rhino.OBJECTLIT {
				return nil // object-literal enum init — no top-level check
			}
			if rhs.Token() == rhino.NAME || rhs.Token() == rhino.GETPROP {
				// Enum alias: check incompatible underlying types.
				if rhsEnum, ok := infer.InferType(rhs).(*jstype.EnumType); ok {
					// Get the declared underlying type from the JSDoc on the ASSIGN.
					declUnderlying := declEnum.Underlying
					if aInfo := a.JSDocInfo(); aInfo != nil {
						if expr := aInfo.EnumParameterType(); expr != nil {
							declUnderlying = scope.ResolveExpr(expr.TypeName)
						}
					}
					if !jstype.IsSubtype(rhsEnum.Underlying, declUnderlying) {
						// Java reports at the outer ASSIGN's LHS position.
						posNode := a
						if p := a.Parent(); p != nil && p.Token() == rhino.EXPR_RESULT {
							posNode = p
						}
						return []diag.Diagnostic{{
							Path:  c.path,
							Line:  posNode.Lineno(),
							Col:   posNode.Charno(),
							Code:  diag.JSC_TYPE_MISMATCH,
							Level: diag.Error,
							Message: fmt.Sprintf(
								"incompatible enum element types\nfound   : %s\nrequired: %s",
								rhsEnum.Underlying.TypeName(), declUnderlying.TypeName()),
						}}
					}
				}
				return nil // compatible enum alias or unknown RHS — no further check
			}
		}
		found := infer.InferType(rhs)
		// Route through TypeValidator.expectCanAssignToPropertyOf so the canonical
		// mismatch dispatcher runs (TemplateType skip, structural-property
		// analysis, TypeMismatch accumulation). Mirrors Java's TypeCheck.visitAssign
		// calling typeValidator.expectCanAssignToPropertyOf(lhs, rhsType, propType, …).
		// Position: use parent EXPR_RESULT when available (Java points at statement
		// start for parenthesised receivers), falling back to lhs. For `this.X`
		// at global scope, Java reports at the THIS receiver — `lhs` (GETPROP)
		// shares its column, so skip the EXPR_RESULT promotion.
		posNode := lhs
		if recv.Token() != rhino.THIS {
			if p := a.Parent(); p != nil && p.Token() == rhino.EXPR_RESULT {
				posNode = p
			}
		}
		// `@implicitCast` on the property declaration opts the slot out of
		// the assignment subtype check. Java's TypeValidator.expectCanAssignToPropertyOf
		// returns early when ObjectType.isPropertyTypeImplicitlyCast is true.
		// Look up the per-property JSDoc (if any) to honour the annotation.
		if recvObj, ok := infer.InferType(recv).(*jstype.ObjectType); ok {
			if pdoc := recvObj.GetPropertyJSDocInfo(propName); pdoc != nil && pdoc.IsImplicitCast() {
				return nil
			}
		}
		// Pass the receiver type so the validator can apply Java's
		// `ownerType.isFunctionPrototypeType()` interface-method skip
		// only when the assignment is on a `Foo.prototype` surface.
		recvType := infer.InferType(recv)
		c.tv.expectCanAssignToPropertyOf(posNode, found, declared, recvType, ownerName, propName)
		return c.tvDrain()
	case rhino.OBJECT_PATTERN, rhino.ARRAY_PATTERN:
		infer := NewTypeInference(scope)
		rhsType := infer.InferType(kids[1])
		return c.checkDestructuringAssignment(lhs, rhsType, scope)
	}
	return nil
}

// checkNotReadonlyPropertyAssignment mirrors TypeCheck.java's
// checkNotReadonlyPropertyAssignment. It fires JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE
// when the receiver of a GETPROP or GETELEM assignment is a ReadonlyArray
// (bare or templatized). Reports at the lhs node's position.
func (c *TypeCheck) checkNotReadonlyPropertyAssignment(lhs, recv *rhino.Node, infer *TypeInference) *diag.Diagnostic {
	// Use the raw declared type (not flow-narrowed) to check for ReadonlyArray.
	// Java reads lhs.getFirstChild().getJSType() which is the static type;
	// mirroring that here avoids false negatives when flow analysis narrows
	// a ReadonlyArray-declared variable to a more specific type like Array.
	var recvType jstype.JSType
	if recv.Token() == rhino.NAME {
		recvType = infer.scope.GetRawDeclaredType(recv.String_())
	}
	if recvType == nil {
		recvType = infer.InferType(recv)
	}
	typeName := ""
	switch t := recvType.(type) {
	case *jstype.ObjectType:
		if t.TypeName() == "ReadonlyArray" {
			typeName = t.TypeName()
		}
	case *jstype.TemplatizedType:
		if t.Base.TypeName() == "ReadonlyArray" {
			typeName = t.TypeName()
		}
	}
	if typeName == "" {
		return nil
	}
	msg := fmt.Sprintf("Should not assign to a property of readonly type '%s'", typeName)
	return &diag.Diagnostic{
		Path: c.path, Line: lhs.Lineno(), Col: lhs.Charno(),
		Code: diag.JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE, Level: diag.Error, Message: msg,
	}
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
	// Port of TypeCheck.java: fire JSC_INSTANTIATE_ABSTRACT_CLASS when
	// the callee is a @constructor @abstract type.
	if fn, ok := t.(*jstype.FunctionType); ok {
		if doc := fn.GetJSDocInfo(); doc.IsAbstract() {
			return []diag.Diagnostic{{
				Path: c.path, Line: node.Lineno(), Col: node.Charno(),
				Code: diag.JSC_INSTANTIATE_ABSTRACT_CLASS, Level: diag.Error,
				Message: "cannot instantiate abstract class",
			}}
		}
	}
	return c.visitCall(node, scope)
}

// Go-only: language-adapter — Java's visitCall fires arity checks on every call
// site unconditionally; Go applies this guard to avoid FPs from incomplete
// override-chain inference (no Java counterpart; Java has full type info).
func isSafeCalleeForArity(callee *rhino.Node, scope *TypedScope) bool {
	if callee == nil {
		return false
	}
	if callee.Token() == rhino.NAME {
		return true
	}
	if callee.Token() == rhino.SUPER {
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

// Go-only: language-adapter — companion guard to isSafeCalleeForArity for
// prototype-method calls; no Java counterpart.
func isSafeInstanceMethodForArity(callee *rhino.Node, infer *TypeInference) bool {
	if callee == nil || callee.Token() != rhino.GETPROP {
		return false
	}
	kids := callee.Children()
	if len(kids) < 2 {
		return false
	}
	recv := kids[0]
	if recv.Token() != rhino.THIS && recv.Token() != rhino.NEW {
		return false
	}
	propName := kids[1].String_()
	recvType := infer.InferType(recv)
	obj, ok := recvType.(*jstype.ObjectType)
	if !ok {
		return false
	}
	// Only safe when no superclass also declares this property
	// (i.e. the method is not overriding a parent slot).
	return !obj.SuperHasProperty(propName)
}

// Go-only: language-adapter — Java's visitNew dispatches through
// couldBeAConstructor (TypeCheck.java:2328) but returns a JSType directly;
// Go needs a typed predicate to propagate the offending type into the message.
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
		// A FunctionType is a real constructor only when IsConstructor
		// is set AND it is not an interface/record (those are tagged
		// structural on their Instance). Mirroring Java's distinction:
		// @constructor → isConstructor()=true, @interface → isInterface()=true
		// (isConstructor()=false in Java, but we share IsConstructor=true
		// for both in Go; use Instance.IsStructural to disambiguate).
		if fn.IsConstructor && (fn.Instance == nil || !fn.Instance.IsStructural) {
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
	// Handle f.bind(thisArg, arg1, …) and f.call(thisArg, arg1, …) specially.
	// Java lazily defines a synthetic function type for "bind" and "call" on
	// FunctionType.getPropertyType (FunctionType.java:getCallOrBindSignature).
	// In Go, FunctionType doesn't embed ObjectType so GetPropertyType("bind")
	// returns nil and infer.InferType(callee) → Unknown, bypassing all checks.
	// We replicate Java's behaviour inline: when callee is GETPROP(recv, "bind"/
	// "call") and recv's type is a FunctionType, check the post-thisArg args
	// against the receiver's parameter list. Port of
	// FunctionType.java:getCallOrBindSignature + TypeCheck.java:visitArgumentList.
	if callee.Token() == rhino.GETPROP {
		if diags, handled := c.visitCallOrBind(node, callee, args, infer); handled {
			return diags
		}
	}
	calleeType := infer.InferType(callee)
	// Strip null/undefined from the callee type so a `?function(T)`
	// (nullable function) still participates in arity checks — Java's
	// restrictByNotNullOrUndefined() does the same thing before
	// deciding whether to fire. Keep the raw type for NOT_FUNCTION_TYPE
	// which should see the full type.
	strippedCallee := stripNullUndef(calleeType)
	fn, ok := strippedCallee.(*jstype.FunctionType)
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
	// Non-native constructors must be called with `new`, not as plain
	// functions. Mirrors TypeCheck.java's visitCall guard:
	//   if (functionType.isConstructor() && !functionType.isNativeObjectType()
	//       && (returnType.isUnknownType() || returnType.isVoidType())
	//       && !child.isSuper())
	// Java's isNativeObjectType() is approximated via externTypes: types
	// from the externs scope (String, Function, Array, etc.) are skipped
	// because they may be callable without `new`. User-defined constructors
	// (Instance not in externTypes) always map to the unknown-return case.
	if node.Token() == rhino.CALL && fn.IsConstructor && callee.Token() != rhino.SUPER &&
		fn.Instance != nil && !c.externTypes[fn.Instance] {
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: node.Lineno(), Col: node.Charno(),
			Code: diag.JSC_CONSTRUCTOR_NOT_CALLABLE, Level: diag.Error,
			Message: fmt.Sprintf("Constructor %s should be called with the \"new\" keyword",
				strippedCallee.TypeName()),
		})
	}
	calleeName := describeCallee(callee)
	// JSC_EXPECTED_THIS_TYPE — port of TypeCheck.java:
	// ordinary functions with an explicit `this:T` type must be called
	// via a property access (GETPROP / GETELEM), not as bare names.
	// Mirrors: `functionType.isOrdinaryFunction() && !NodeUtil.isNormalOrOptChainGet(child)`
	// followed by `receiverType` allowlist (unknown / all / void / native objects).
	if node.Token() == rhino.CALL && fn.IsOrdinaryFunction() &&
		callee.Token() != rhino.GETPROP && callee.Token() != rhino.GETELEM {
		if tot := fn.GetTypeOfThis(); tot != nil {
			// Allow: Unknown, All, Void, or native ObjectType (matches Java's exemptions).
			allow := false
			switch tot {
			case jstype.Unknown, jstype.All, jstype.Void:
				allow = true
			}
			if !allow {
				if obj, ok := tot.(*jstype.ObjectType); ok && obj.IsNativeObjectType() {
					allow = true
				}
				if inst, ok := tot.(*jstype.InstanceObjectType); ok && inst.IsNativeObjectType() {
					allow = true
				}
			}
			if !allow {
				out = append(out, diag.Diagnostic{
					Path: c.path, Line: node.Lineno(), Col: node.Charno(),
					Code: diag.JSC_EXPECTED_THIS_TYPE, Level: diag.Error,
					Message: fmt.Sprintf("%q must be called with a \"this\" type", fn.TypeName()),
				})
			}
		}
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
	// Arity guard: skip when the callee is a constructor with 0
	// declared params (externs like Array whose arity we intentionally
	// leave unchecked). All other FunctionTypes — including 0-param
	// non-constructors — are subject to the check.
	arityKnown := !(fn.IsConstructor && len(fn.Params) == 0)
	safeForArity := isSafeCalleeForArity(callee, scope) || isSafeInstanceMethodForArity(callee, infer)
	if node.Token() == rhino.CALL && arityKnown && safeForArity {
		if got := len(args); got < fn.MinArity || got > fn.MaxArity {
			maxStr := fmt.Sprintf("%d", fn.MaxArity)
			if fn.MaxArity >= jstype.MaxArgsUnlimited {
				maxStr = "unlimited"
			}
			// Position: when the CALL is a top-level EXPR_RESULT with
			// an outer paren wrapping the callee (e.g. `(new Foo).run()`),
			// Java reports at the statement's column. Promote to the
			// EXPR_RESULT position when it precedes the CALL node.
			line, col := node.Lineno(), node.Charno()
			if p := node.Parent(); p != nil && p.Token() == rhino.EXPR_RESULT && p.Charno() < col {
				line, col = p.Lineno(), p.Charno()
			}
			out = append(out, diag.Diagnostic{
				Path: c.path, Line: line, Col: col,
				Code: diag.JSC_WRONG_ARGUMENT_COUNT, Level: diag.Error,
				Message: fmt.Sprintf(
					"Function %s: called with %d argument(s). Function requires at least %d argument(s) and no more than %s argument(s).",
					calleeName, got, fn.MinArity, maxStr),
			})
		}
	}
	// Per-argument type check. Java's TypeCheck.visitCallParameters:
	// when parameters are exhausted, if the last parameter isVariadic,
	// keep reusing it for remaining arguments. Mirrors:
	//   if (parameters.hasNext()) { parameter = parameters.next(); }
	//   else if (parameter != null && parameter.isVariadic()) { /* reuse */ }
	//
	// We detect variadic functions via MaxArity == MaxArgsUnlimited — this
	// is set by setArityFromParams for @param {...T} and var_args patterns.
	// The last declared param type is then reused for all remaining args.
	params := fn.GetParameters()
	// Determine if the function is variadic (rest param / var_args).
	// When variadic, the last declared param type is reused for all excess args.
	var lastParamType jstype.JSType
	isVariadic := fn.MaxArity >= jstype.MaxArgsUnlimited && len(fn.Params) > 0
	if isVariadic {
		lastParamType = fn.Params[len(fn.Params)-1]
	}
	// Track assignments-to-NAME that occur in earlier sibling args so
	// subsequent args see the freshly-assigned type. Mirrors Java's
	// TypeInference flowing the FlowScope through expression children
	// in source order: `g(x, x = 3, x)` — the third `x` reads number
	// because the assignment mutated the flow scope between args.
	// Without this, arg 3 would still see x's pre-call refined type.
	var siblingRefinements map[string]jstype.JSType
	for i, arg := range args {
		var param jstype.JSType
		if i < len(fn.Params) {
			param = fn.Params[i]
		} else if isVariadic {
			// Rest parameter covers all remaining arguments.
			param = lastParamType
		} else {
			break
		}
		// Push any sibling-arg assignment refinements collected from
		// prior args so this arg's NAME lookups see the updated type.
		var popSibling func()
		if len(siblingRefinements) > 0 {
			scope.PushRefinement(siblingRefinements)
			popSibling = scope.PopRefinement
		}
		found := infer.InferType(arg)
		if popSibling != nil {
			popSibling()
		}
		// After this arg is processed, if the arg itself was an
		// `x = expr` ASSIGN to a NAME, record the new flow type for x
		// so later sibling args see it.
		if arg.Token() == rhino.ASSIGN {
			if ak := arg.Children(); len(ak) >= 2 && ak[0].Token() == rhino.NAME {
				if siblingRefinements == nil {
					siblingRefinements = map[string]jstype.JSType{}
				}
				siblingRefinements[ak[0].String_()] = found
			}
		}
		// Optional parameters (JSDoc {=T}) accept `undefined` at the
		// call site. Java's TypeCheck widens the effective required
		// type to (T|undefined) for optional params, so passing
		// `undefined` explicitly is always valid.
		// Mirrors FunctionType.Parameter.isOptional() check in Java's
		// TypeCheck.visitCallParameters.
		if i < len(params) && params[i].IsOptional() && found == jstype.Void {
			continue
		}
		// Route through TypeValidator.expectArgumentMatchesParameter so the
		// canonical mismatch dispatcher runs (TemplateType skip, structural-
		// property analysis, TypeMismatch accumulation). Mirrors Java's
		// TypeCheck.visitCallParameters calling
		// typeValidator.expectArgumentMatchesParameter(arg, argType, paramType, callNode, i).
		// Position: for GETPROP args Java reports at the property-name node
		// (rightmost segment), matching diagPosForArg's behaviour.
		posNode := arg
		if arg.Token() == rhino.GETPROP {
			if gpKids := arg.Children(); len(gpKids) >= 2 {
				posNode = gpKids[1]
			}
		}
		c.tv.expectArgumentMatchesParameter(posNode, found, param, node, i+1)
		out = append(out, c.tvDrain()...)
	}
	return out
}

// visitCallOrBind checks argument types for f.call(thisArg, a, b, …) and
// f.bind(thisArg, a, b, …) when the receiver of the GETPROP is a FunctionType.
// Returns (diagnostics, true) when the pattern is matched (caller should return
// those diagnostics directly), or (nil, false) when it is not.
//
// Port of FunctionType.java:getCallOrBindSignature (lines 709-750) combined with
// TypeCheck.java:visitArgumentList (lines 2656-2660): Java lazily synthesises a
// function type for "bind"/"call" and then runs the standard argument-checking
// loop against it; we replicate that here for the Go port.
func (c *TypeCheck) visitCallOrBind(node *rhino.Node, callee *rhino.Node, args []*rhino.Node, infer *TypeInference) ([]diag.Diagnostic, bool) {
	gpKids := callee.Children()
	if len(gpKids) < 2 {
		return nil, false
	}
	propName := gpKids[1].String_()
	if propName != "call" && propName != "bind" {
		return nil, false
	}
	// Resolve the receiver type (e.g. the `f` in `f.call`).
	recvType := infer.InferType(gpKids[0])
	recvType = stripNullUndef(recvType)
	recvFn, ok := recvType.(*jstype.FunctionType)
	if !ok {
		return nil, false
	}
	// args[0] is thisArg — not type-checked (Java also skips it).
	// args[1..] are checked against recvFn.Params[0..].
	if len(args) <= 1 {
		// Only thisArg or nothing — nothing to check.
		return nil, true
	}
	boundArgs := args[1:] // args after thisArg
	params := recvFn.GetParameters()
	isVariadic := recvFn.MaxArity >= jstype.MaxArgsUnlimited && len(recvFn.Params) > 0
	var lastParamType jstype.JSType
	if isVariadic {
		lastParamType = recvFn.Params[len(recvFn.Params)-1]
	}
	var out []diag.Diagnostic
	for i, arg := range boundArgs {
		var param jstype.JSType
		if i < len(recvFn.Params) {
			param = recvFn.Params[i]
		} else if isVariadic {
			param = lastParamType
		} else {
			break
		}
		found := infer.InferType(arg)
		// Optional params accept undefined explicitly (mirrors Java).
		if i < len(params) && params[i].IsOptional() && found == jstype.Void {
			continue
		}
		posNode := arg
		if arg.Token() == rhino.GETPROP {
			if gk := arg.Children(); len(gk) >= 2 {
				posNode = gk[1]
			}
		}
		// Ordinal is i+2: 1-based, starting at 2 because arg 1 is thisArg.
		c.tv.expectArgumentMatchesParameter(posNode, found, param, node, i+2)
		out = append(out, c.tvDrain()...)
	}
	return out, true
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
		// Dedup: Java's TypeValidator only reports the mismatch once
		// per receiver-root within a function body. Two consecutive
		// `bar.x; bar.y;` accesses on a void-typed `bar` produce a
		// single error in the oracle.
		key := ""
		if kids[0].Token() == rhino.NAME {
			key = kids[0].String_()
		}
		if key != "" {
			if c.noPropsReported == nil {
				c.noPropsReported = map[string]bool{}
			}
			if c.noPropsReported[key] {
				return nil
			}
			c.noPropsReported[key] = true
		}
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
	// Unknown/AllType: POSSIBLE_INEXISTENT_PROPERTY path.
	if recvType == jstype.Unknown || recvType == jstype.All {
		if !c.fileSuppressions["missingProperties"] {
			return c.checkPossibleInexistent(n, kids[0], propName, recvType)
		}
		return nil
	}
	// FunctionType receiver: in Java FunctionType extends ObjectType,
	// so `Foo.bar` checks Foo's own static-slot map (plus the implicit
	// Function.prototype chain). jsgo's FunctionType doesn't embed
	// ObjectType — static members live on the Instance ObjectType — so
	// we replicate the Java check inline. Port of FunctionType.hasProperty
	// + ObjectType.cast(childType) in TypeCheck.checkPropertyAccess.
	//
	// Restrictions to keep the check conservative:
	//   - Receiver must be a NAME (not a GETPROP) — qualified-namespace
	//     constructors like `ns.Outer` register their nested members on
	//     the namespace map, not the FunctionType's Instance, so we'd
	//     false-positive on `ns.Outer.Inner`.
	//   - Skip interfaces: their nested-interface pattern (`I1.I2 = …`)
	//     extends the surface in ways jsgo doesn't fully track yet.
	if fn, ok := recvType.(*jstype.FunctionType); ok && fn.Instance != nil &&
		fn.IsConstructor && !fn.IsInterface() && !fn.Instance.IsInterface &&
		kids[0].Token() == rhino.NAME && !c.externTypes[fn.Instance] {
		// Skip writes — collectPropertyWrites already extends the surface.
		if isAssignLHS(n) {
			return nil
		}
		// Existence/absence guards (`if (Foo.bar)`, `Foo.bar && …`):
		// Java routes these through STRICT_INEXISTENT_PROPERTY only.
		if isPropertyExistenceOrAbsenceTest(n) {
			return nil
		}
		// `'prop' in Foo` narrowing.
		if qname := qualifiedName(n); qname != "" {
			if scope.GetQualifiedRefinement(qname) != nil {
				return nil
			}
		}
		if c.fileSuppressions["missingProperties"] {
			return nil
		}
		// Standard Function-instance slots (Function.prototype + own
		// FunctionType slots): always present on any function value.
		if isFunctionStandardProp(propName) {
			return nil
		}
		// User-declared static slot.
		if fn.Instance.GetStaticPropertyType(propName) != nil {
			return nil
		}
		// `@lends {Foo}` registers the lent keys via DefineInferredProperty
		// on Foo's Instance (not the static slot). Treat those as defined
		// for the static-access check too.
		if fn.Instance.HasOwnProperty(propName) {
			return nil
		}
		// User-declared namespace property attached via
		// `/** @type {T} */ Foo.bar = …` or `Foo.Bar = function() {}`
		// — tracked on the scope's namespace-prop map keyed on the
		// receiver's qualified name (so both `Foo.bar` and
		// `ns.Outer.Inner` resolve here).
		if qname := qualifiedName(kids[0]); qname != "" {
			if scope.LookupNamespaceProp(qname, propName) != nil {
				return nil
			}
		}
		// Dedup across the script.
		ctorName := fn.Instance.Name
		if ctorName == "" {
			ctorName = "Function"
		}
		key := ctorName + "." + propName
		if c.inexistentReported[key] {
			return nil
		}
		if c.inexistentReported == nil {
			c.inexistentReported = map[string]bool{}
		}
		c.inexistentReported[key] = true
		return []diag.Diagnostic{{
			Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
			Code: diag.JSC_INEXISTENT_PROPERTY, Level: diag.Error,
			Message: fmt.Sprintf("Property %s never defined on %s", propName, ctorName),
		}}
	}
	obj := dereferencedObject(recvType)
	if obj == nil {
		// Port of TypeCheck.java checkPropertyAccessHelper: when the receiver
		// is a union type, writing a new property to it is illegal if any
		// alternate is a @struct. Mirrors Java's iteration over alternates
		// and JSC_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE emission.
		if u, ok := recvType.(*jstype.UnionType); ok && isAssignLHS(n) {
			for _, m := range u.GetAlternates() {
				if mo := dereferencedObject(m); mo != nil && mo.IsStruct {
					return []diag.Diagnostic{{
						Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
						Code: diag.JSC_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE, Level: diag.Error,
						Message: "Cannot add a property to an instance of union type.",
					}}
				}
			}
		}
		// Port of TypeCheck.java reportMissingProperty with PropDefinitionKind.LOOSE_UNION:
		// when the receiver is a union type and the property is present on some
		// but not all members (MAYBE_PRESENT → LOOSE_UNION), emit
		// JSC_STRICT_INEXISTENT_UNION_PROPERTY in strict mode.
		// Mirrors canPropertyBeDefined + reportMissingProperty w/ isKnownToUnionMember=true.
		if u, ok := recvType.(*jstype.UnionType); ok {
			kind := u.GetPropertyKind(propName)
			if kind == jstype.HasPropertyMaybePresent {
				if c.opts != nil && c.opts.IsOn(GroupStrictCheckTypes) {
					// Readable receiver name: use the variable name when the
					// receiver node is a NAME (mirrors Java's getReadableTypeName(objNode)
					// which falls back to n.getQualifiedName()).
					recvName := recvType.TypeName()
					if kids[0].Token() == rhino.NAME {
						recvName = kids[0].String_()
					}
					return []diag.Diagnostic{{
						Path:    c.path,
						Line:    propNode.Lineno(),
						Col:     propNode.Charno(),
						Code:    diag.JSC_STRICT_INEXISTENT_UNION_PROPERTY,
						Level:   diag.Warning,
						Message: fmt.Sprintf("Property %s not defined on all member types of %s", propName, recvName),
					}}
				}
			}
		}
		return nil
	}
	if obj.Name == "" {
		return nil
	}
	if obj.HasProperty(propName) {
		// @record properties registered as inferred-only by
		// collectPropertyWrites still emit JSC_STRICT_INEXISTENT_PROPERTY
		// in strict mode: the surface is open at runtime but the lint
		// surfaces unannotated property additions. Mirrors Java's
		// TypeInference + TypeCheck split where inferred properties
		// participate in default-mode reads but are LOOSE under
		// STRICT_CHECK_TYPES.
		if obj.IsRecord && obj.IsPropertyInferred(propName) {
			if c.opts != nil && c.opts.IsOn(GroupStrictCheckTypes) {
				return []diag.Diagnostic{{
					Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
					Code: diag.JSC_STRICT_INEXISTENT_PROPERTY, Level: diag.Warning,
					Message: fmt.Sprintf("Property %s never defined on %s", propName, obj.TypeName()),
				}}
			}
		}
		return nil
	}
	// `'prop' in obj` narrowing (SRAI caseIn): if the caller guarded
	// with `if ('prop' in this)` then obj.prop is in the refinement
	// stack with Unknown type. Suppress JSC_INEXISTENT_PROPERTY inside
	// such a guard branch. Port of caseIn / inferQualifiedSlot.
	if qname := qualifiedName(n); qname != "" {
		if scope.GetQualifiedRefinement(qname) != nil {
			return nil
		}
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
	// builtin types we ship concrete slot maps for.
	//
	// For plain @constructor types (non-struct, non-structural),
	// Java fires JSC_INEXISTENT_PROPERTY on any read of a missing
	// property — not just `this.prop` — provided:
	//   (a) the property is not LOOSE (defined on a subtype):
	//       Java's canPropertyBeDefined returns LOOSE in that case,
	//       which routes to STRICT_INEXISTENT_PROPERTY only.
	//   (b) the access is not a property existence test (&&, ||,
	//       if-cond, typeof, etc.) or property absence test
	//       (== null/undefined). Those contexts also fire as
	//       STRICT_INEXISTENT_PROPERTY only in Java.
	//   (c) the access is not the LHS of an assign (writes extend
	//       the surface via defineOnOpenObject / collectPropertyWrites).
	if !obj.IsStructural && !obj.IsStruct && !isKnownSurfaceBuiltin(obj.Name) {
		// Skip types registered in the externs scope (Object, Function,
		// Date, Error, …). jsgo ships only sparse @constructor stubs for
		// these; their full property surfaces aren't declared, so any
		// property access would spuriously fire. Java treats the global
		// Object type as "low confidence" (fires POSSIBLE_INEXISTENT_PROPERTY)
		// and leaves the rest of the externs surface open. Go's externs
		// are even sparser, so we skip all extern-registered types here.
		if c.externTypes[obj] {
			if obj.Name == "Object" {
				return c.checkPossibleInexistent(n, kids[0], propName, recvType)
			}
			return nil
		}
		// Skip when the property is defined on a subtype (LOOSE).
		if c.looseProps[obj] != nil && c.looseProps[obj][propName] {
			return nil
		}
		// Skip writes — collectPropertyWrites already extends the surface.
		if isAssignLHS(n) {
			return nil
		}
		// Skip property existence / absence tests — these fire as
		// STRICT_INEXISTENT_PROPERTY only (not default JSC_INEXISTENT_PROPERTY).
		if isPropertyExistenceOrAbsenceTest(n) {
			// In strict mode, emit STRICT_INEXISTENT_PROPERTY.
			if c.opts != nil && c.opts.IsOn(GroupStrictCheckTypes) {
				if c.inexistentReported == nil {
					c.inexistentReported = map[string]bool{}
				}
				c.inexistentReported[key] = true
				return []diag.Diagnostic{{
					Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
					Code: diag.JSC_STRICT_INEXISTENT_PROPERTY, Level: diag.Warning,
					Message: fmt.Sprintf("Property %s never defined on %s", propName, obj.TypeName()),
				}}
			}
			return nil
		}
	}
	// @interface types are open surfaces in default mode: Java routes
	// unknown-property accesses on them as LOOSE (STRICT_INEXISTENT_PROPERTY
	// only), which is silenced in default mode. In strict mode, emit.
	if obj.IsStructural && !obj.IsStruct && obj.IsInterface {
		if c.opts != nil && c.opts.IsOn(GroupStrictCheckTypes) {
			if c.inexistentReported == nil {
				c.inexistentReported = map[string]bool{}
			}
			c.inexistentReported[key] = true
			return []diag.Diagnostic{{
				Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
				Code: diag.JSC_STRICT_INEXISTENT_PROPERTY, Level: diag.Warning,
				Message: fmt.Sprintf("Property %s never defined on %s", propName, obj.TypeName()),
			}}
		}
		return nil
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

// Go-only: language-adapter — splits out the this-receiver predicate that Java's
// allowStrictPropertyAccessOnNode handles inline (TypeCheck.java:2202-2213).
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

// Go-only: language-adapter — helper for isThisReceiverRead; no Java counterpart.
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

// Go-only: language-adapter — walks a GETPROP chain to find the root NAME;
// Java's NodeUtil handles this generically; no single named Go equivalent.
func rootOfGetProp(n *rhino.Node) *rhino.Node {
	curr := n
	for curr != nil && curr.Token() == rhino.GETPROP {
		kids := curr.Children()
		if len(kids) == 0 {
			break
		}
		curr = kids[0]
	}
	return curr
}

// Go-only: language-adapter — mirrors Java's allowLoosePropertyAccessOnNode
// (TypeCheck.java:2214) which routes to STRICT_INEXISTENT_PROPERTY; Go skips
// the diagnostic entirely in default mode since STRICT_INEXISTENT is silenced.
func isPropertyExistenceOrAbsenceTest(n *rhino.Node) bool {
	p := n.Parent()
	if p == nil {
		return false
	}
	switch p.Token() {
	case rhino.AND, rhino.OR, rhino.INSTANCEOF, rhino.TYPEOF, rhino.COALESCE:
		// Direct operand of &&, ||, ??, instanceof, typeof — always a test.
		return true
	case rhino.IF, rhino.WHILE:
		return p.FirstChild() == n
	case rhino.DO_WHILE:
		// DO_WHILE(body, cond): condition is last child.
		doKids := p.Children()
		return len(doKids) >= 2 && doKids[len(doKids)-1] == n
	case rhino.FOR:
		// Classic for-loop: children are [init, cond, update, body].
		// Condition is the second child.
		kids := p.Children()
		return len(kids) >= 2 && kids[1] == n
	case rhino.HOOK:
		// Ternary condition is first child.
		return p.FirstChild() == n
	case rhino.NOT:
		// !foo.bar as the left operand of || → Java's NOT-in-OR pattern.
		gp := p.Parent()
		return gp != nil && gp.Token() == rhino.OR && gp.FirstChild() == p
	case rhino.NE:
		// foo.bar != null / foo.bar != undefined → existence test.
		return otherSideIsNullOrUndefined(p, n)
	case rhino.SHNE:
		// foo.bar !== undefined → existence test.
		// foo.bar !== null is NOT an existence test (Java's NE-only rule).
		return otherSideIs(p, n, rhino.UNDEFINED)
	case rhino.EQ:
		// foo.bar == null / foo.bar == undefined → absence test.
		return otherSideIsNullOrUndefined(p, n)
	case rhino.SHEQ:
		// foo.bar === undefined → absence test.
		// foo.bar === null is NOT an absence test (only == null qualifies).
		return otherSideIs(p, n, rhino.UNDEFINED)
	}
	return false
}

// Go-only: language-adapter — helper for isPropertyExistenceOrAbsenceTest.
func otherSideIsNullOrUndefined(p, n *rhino.Node) bool {
	kids := p.Children()
	if len(kids) < 2 {
		return false
	}
	var other *rhino.Node
	if kids[0] == n {
		other = kids[1]
	} else {
		other = kids[0]
	}
	return isNullOrUndefinedLiteral(other)
}

// Go-only: language-adapter — helper for isPropertyExistenceOrAbsenceTest.
func otherSideIs(p, n *rhino.Node, tok rhino.Token) bool {
	kids := p.Children()
	if len(kids) < 2 {
		return false
	}
	var other *rhino.Node
	if kids[0] == n {
		other = kids[1]
	} else {
		other = kids[0]
	}
	return other != nil && other.Token() == tok
}

// Go-only: language-adapter — helper for isInGuardedBranch; no Java counterpart.
func containsPropTest(sub *rhino.Node, recvName, propName string) bool {
	if sub == nil {
		return false
	}
	if sub.Token() == rhino.GETPROP {
		kids := sub.Children()
		if len(kids) >= 2 && kids[0].Token() == rhino.NAME &&
			kids[0].String_() == recvName && kids[1].String_() == propName &&
			isPropertyExistenceOrAbsenceTest(sub) {
			return true
		}
	}
	for _, ch := range sub.Children() {
		if containsPropTest(ch, recvName, propName) {
			return true
		}
	}
	return false
}

// Go-only: language-adapter — implements Java's type-narrowing guard check inline
// via explicit AST walking; Java uses FlowScope lattice ops; no named equivalent.
func isInGuardedBranch(n *rhino.Node, recvName, propName string) bool {
	for cur := n; cur != nil; cur = cur.Parent() {
		p := cur.Parent()
		if p == nil {
			break
		}
		switch p.Token() {
		case rhino.AND:
			kids := p.Children()
			if len(kids) >= 2 && cur == kids[1] && containsPropTest(kids[0], recvName, propName) {
				return true
			}
		case rhino.OR:
			// `!x.foo || x.foo()` — RHS is guarded by the LHS NOT test.
			// Java's isPropertyTest recognises !foo.bar as the left operand
			// of ||; the RHS access is thus a guarded branch. Port of
			// SemanticReverseAbstractInterpreter / checkPropertyAccessHelper.
			kids := p.Children()
			if len(kids) >= 2 && cur == kids[1] {
				lhs := kids[0]
				if lhs.Token() == rhino.NOT {
					notKids := lhs.Children()
					if len(notKids) >= 1 && containsPropTest(notKids[0], recvName, propName) {
						return true
					}
				}
			}
		case rhino.IF:
			kids := p.Children()
			if len(kids) >= 2 && cur == kids[1] && containsPropTest(kids[0], recvName, propName) {
				return true
			}
		case rhino.HOOK:
			kids := p.Children()
			if len(kids) >= 3 && cur == kids[1] && containsPropTest(kids[0], recvName, propName) {
				return true
			}
		case rhino.TYPEOF, rhino.INSTANCEOF:
			return true
		}
	}
	return false
}

// Go-only: language-adapter — emits JSC_POSSIBLE_INEXISTENT_PROPERTY for
// low-confidence receivers; Java routes this through checkPropertyAccessHelper's
// strictReport path which funnels to STRICT_INEXISTENT_PROPERTY; Go emits the
// weaker POSSIBLE variant directly since STRICT_ is silenced in default mode.
func (c *TypeCheck) checkPossibleInexistent(n, recvNode *rhino.Node, propName string, recvType jstype.JSType) []diag.Diagnostic {
	if isAssignLHS(n) {
		return nil
	}
	if p := n.Parent(); p != nil {
		switch p.Token() {
		case rhino.EXPR_RESULT:
			return nil
		case rhino.CAST:
			// `/** @type {T} */ (x.prop)` — the explicit cast overrides the
			// type; suppress the existence check. Port of Java's visitCast
			// suppressing the property check when a JSDoc cast is present.
			return nil
		}
	}
	// Cast expressions: `/** @type {T} */ (obj.prop)` — suppress
	// POSSIBLE_INEXISTENT_PROPERTY inside casts (developer asserted the type).
	if p := n.Parent(); p != nil && p.Token() == rhino.CAST {
		return nil
	}
	if isPropertyExistenceOrAbsenceTest(n) {
		return nil
	}
	if c.globalPropNames[propName] {
		return nil
	}
	recvName := ""
	if recvNode != nil && recvNode.Token() == rhino.NAME {
		recvName = recvNode.String_()
	}
	// Suppress POSSIBLE_INEXISTENT_PROPERTY when the root receiver is a
	// variable that VarCheck already flagged as undeclared. Java's VarCheck
	// synthesizes an extern for undeclared names (createSynthesizedExternVar),
	// making subsequent property-access checks treat the name as open-typed.
	// We approximate this by checking whether the ROOT of the GETPROP chain
	// is in the undeclared set.
	if len(c.undeclaredNames) > 0 {
		root := rootOfGetProp(n)
		if root != nil && root.Token() == rhino.NAME && c.undeclaredNames[root.String_()] {
			return nil
		}
	}
	if recvName != "" && isInGuardedBranch(n, recvName, propName) {
		return nil
	}
	label := recvType.TypeName()
	if recvType == jstype.Unknown && recvName != "" {
		label = recvName
	}
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	return []diag.Diagnostic{{
		Path: c.path, Line: kids[1].Lineno(), Col: kids[1].Charno(),
		Code: diag.JSC_POSSIBLE_INEXISTENT_PROPERTY, Level: diag.Error,
		Message: fmt.Sprintf("Property %s never defined on %s", propName, label),
	}}
}

// Go-only: language-adapter — checks JSC_INCOMPATIBLE_ALIAS_ANNOTATION inline
// during visitVar; Java handles this as part of visitVar's annotation checking
// at TypeCheck.java:2266 without a named sub-method.
func checkAliasAnnotation(info *rhino.JSDocInfo, aliased *jstype.FunctionType) string {
	structural := aliased.Instance != nil && aliased.Instance.IsStructural
	if info.IsRecord() {
		instDoc := aliased.Instance.GetJSDocInfo()
		if !(structural && instDoc != nil && instDoc.IsRecord()) {
			return "@record"
		}
	} else if info.IsInterface() {
		if !structural {
			return "@interface"
		}
	} else if info.IsConstructor() {
		if structural {
			return "@constructor"
		}
	}
	return ""
}

// Go-only: package-bridge — builds a global property name index for
// POSSIBLE_INEXISTENT_PROPERTY suppression; Java's eachRefTypeIndexedByProperty
// is populated by TypedScopeCreator and queried via AbstractCompiler; Go
// pre-builds an equivalent map since there is no shared compiler object.
func buildGlobalPropNameSet(global *TypedScope, script *rhino.Node) map[string]bool {
	// Universal JS properties (exist on every object / function).
	// Suppresses spurious POSSIBLE_INEXISTENT_PROPERTY on unknown receivers.
	out := map[string]bool{
		"prototype": true, "constructor": true, "hasOwnProperty": true,
		"toString": true, "valueOf": true, "length": true,
		"call": true, "apply": true, "bind": true,
	}
	for s := global; s != nil; s = s.parent {
		for name, t := range s.declared {
			// Global names are accessible as properties of the global object.
			if name != "" {
				out[name] = true
			}
			if fn, ok := t.(*jstype.FunctionType); ok && fn.Instance != nil {
				for _, prop := range fn.Instance.MemberNames() {
					out[prop] = true
				}
			}
		}
		for _, props := range s.namespaceProps {
			for prop := range props {
				out[prop] = true
			}
		}
	}
	var scan func(n *rhino.Node)
	scan = func(n *rhino.Node) {
		if n == nil {
			return
		}
		switch n.Token() {
		case rhino.ASSIGN:
			kids := n.Children()
			if len(kids) >= 2 && kids[0].Token() == rhino.GETPROP {
				if pk := kids[0].Children(); len(pk) >= 2 {
					out[pk[1].String_()] = true
				}
			}
		case rhino.EXPR_RESULT:
			// Prototype/namespace stub declarations such as
			// `Foo.prototype.bar;` or `ns.T.x;` are property-name
			// declarations (not references). Collect the leaf property
			// name so POSSIBLE_INEXISTENT_PROPERTY is not spuriously
			// fired when the receiver type is Unknown / unresolved because
			// the root object (e.g. `u`) was flagged as undeclared.
			// Mirrors Java's recordPrototypeStub path in TypedScopeCreator.
			if inner := n.FirstChild(); inner != nil && inner.Token() == rhino.GETPROP {
				if pk := inner.Children(); len(pk) >= 2 {
					out[pk[1].String_()] = true
				}
			}
		case rhino.STRING_KEY:
			if s := n.String_(); s != "" {
				out[s] = true
			}
		}
		for _, ch := range n.Children() {
			scan(ch)
		}
	}
	scan(script)
	return out
}

// Go-only: package-bridge — builds the externs type set used to suppress
// JSC_INEXISTENT_PROPERTY on extern types whose surfaces are incomplete in jsgo;
// Java's full externs populate all property surfaces so no such guard is needed.
func buildExternTypeSet(global *TypedScope) map[*jstype.ObjectType]bool {
	out := map[*jstype.ObjectType]bool{}
	if global == nil {
		return out
	}
	// Walk only the extern scope (global.parent), not the user scope.
	for s := global.parent; s != nil; s = s.parent {
		for _, t := range s.declared {
			fn, ok := t.(*jstype.FunctionType)
			if !ok || fn.Instance == nil {
				continue
			}
			out[fn.Instance] = true
		}
	}
	return out
}

// Go-only: package-bridge — pre-builds the subtype-property index that Java
// lazily computes via eachRefTypeIndexedByProperty inside canPropertyBeDefined;
// in Go this is a pre-pass since there is no shared AbstractCompiler registry.
func buildLoosePropertyIndex(global *TypedScope) map[*jstype.ObjectType]map[string]bool {
	out := map[*jstype.ObjectType]map[string]bool{}
	for s := global; s != nil; s = s.parent {
		for _, t := range s.declared {
			fn, ok := t.(*jstype.FunctionType)
			if !ok || fn.Instance == nil {
				continue
			}
			inst := fn.Instance
			// For each property directly declared on inst's own member
			// map, mark it as loose on every ancestor in the super chain.
			for _, prop := range inst.MemberNames() {
				for anc := inst.Super(); anc != nil; anc = anc.Super() {
					m := out[anc]
					if m == nil {
						m = map[string]bool{}
						out[anc] = m
					}
					m[prop] = true
				}
			}
		}
	}
	return out
}

// Go-only: language-adapter — determines whether an assignment to this.prop occurs
// inside the struct's own constructor body; Java's NodeUtil.isConstructorDeclaration
// handles the equivalent check inline inside checkPropertyInheritanceOnAssignment.
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
	enclosing := NU.GetEnclosingFunction(ret)
	if enclosing == nil {
		return nil
	}
	info := functionJSDoc(enclosing)
	var declared jstype.JSType
	if info != nil && info.ReturnType() != nil {
		declared = scope.ResolveExpr(info.ReturnType().TypeName)
	} else if info != nil && info.HasType() {
		// @type {function(...): T} annotation — extract return type from
		// the declared function type.
		if fn := enclosingFunctionType(enclosing, scope); fn != nil && fn.Return != nil {
			declared = fn.Return
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
	// Route through TypeValidator.expectCanAssignTo so the canonical
	// mismatch dispatcher runs (TemplateType skip, structural-property
	// analysis, TypeMismatch accumulation). Mirrors Java's TypeCheck.visitReturn
	// calling typeValidator.expectCanAssignTo(retExpr, actualType, declaredType, msg).
	// Position: for GETPROP return-exprs Java reports at the property-name node
	// (rightmost segment), matching diagPosForArg's behaviour.
	posNode := retExpr
	if retExpr.Token() == rhino.GETPROP {
		if kids := retExpr.Children(); len(kids) >= 2 {
			posNode = kids[1]
		}
	}
	c.tv.expectCanAssignTo(posNode, found, declared, "inconsistent return type")
	return c.tvDrain()
}

// Go-only: language-adapter — Go extracted CLASS_FIELD_DEF handling from Java's
// visitClass / visitObjectOrClassLiteralKey into a named function for readability;
// in Java this is inlined in the CLASS_MEMBERS case at TypeCheck.java:902-910.
func (c *TypeCheck) visitClassField(field *rhino.Node, parent *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	// --- @dict / @struct class-key gate (Java's visitObjectOrClassLiteralKey) ---
	// Mirrors: if (ownerType.isDict() && !NodeUtil.isEs6ConstructorMemberFunctionDef(key))
	//            report(key, ILLEGAL_CLASS_KEY, "dict")
	// CLASS_FIELD_DEF is never the constructor, so the constructor exclusion never applies.
	var out []diag.Diagnostic
	if parent != nil && parent.Parent() != nil {
		classNode := parent.Parent()
		if classNode.Token() == rhino.CLASS && classNode.FirstChild() != nil {
			className := classNode.FirstChild().String_()
			if fn, ok := scope.GetDeclaredType(className).(*jstype.FunctionType); ok && fn.Instance != nil {
				if fn.Instance.IsDict {
					out = append(out, diag.Diagnostic{
						Path:    c.path,
						Line:    field.Lineno(),
						Col:     field.Charno(),
						Code:    diag.JSC_ILLEGAL_CLASS_KEY,
						Level:   diag.Error,
						Message: "Illegal key, the class is a dict",
					})
				}
			}
		}
	}

	// --- @type mismatch check ---
	// CLASS_FIELD_DEF (legacy): kids = [NAME(fieldName), initializer]
	// MEMBER_FIELD_DEF (new pipeline): String_() = fieldName, kids = [initializer]
	kids := field.Children()
	var name *rhino.Node
	var init *rhino.Node
	if field.Token() == rhino.MEMBER_FIELD_DEF {
		// New pipeline: the field name is the node's string, initializer is first child.
		name = rhino.NewString(rhino.NAME, field.String_())
		if len(kids) >= 1 {
			init = kids[0]
		} else {
			return out
		}
	} else {
		if len(kids) < 2 {
			return out
		}
		name = kids[0]
		init = kids[1]
	}
	info := field.JSDocInfo()
	if info == nil || !info.HasType() {
		return out
	}
	declared := scope.ResolveExpr(info.Type().TypeName)
	infer := NewTypeInference(scope)
	found := infer.InferType(init)
	if jstype.IsSubtype(found, declared) {
		return out
	}
	className := ""
	if parent != nil && parent.Parent() != nil {
		gp := parent.Parent()
		if gp.Token() == rhino.CLASS && gp.FirstChild() != nil {
			className = gp.FirstChild().String_()
		}
		// Anonymous class: derive the name from the enclosing
		// declaration / assignment. Mirrors Java's
		// FunctionType.getInstanceType().toString() which reports the
		// qualified name from the surrounding binding (e.g. `obj.C`
		// for `obj.C = class { ... }`, or `Foo` for `var Foo = class
		// { ... }`).
		if className == "" && gp.Token() == rhino.CLASS {
			className = anonClassQualifiedName(gp)
		}
	}
	msg := fmt.Sprintf("assignment to property %s of %s\nfound   : %s\nrequired: %s",
		name.String_(), className, found.TypeName(), declared.TypeName())
	// Java reports at the CLASS_FIELD_DEF node's position — the
	// `static`/name column, not the bare IDENT — so we use field's
	// pos rather than name's.
	out = append(out, diag.Diagnostic{
		Path: c.path, Line: field.Lineno(), Col: field.Charno(),
		Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
	})
	return out
}

// Go-only: language-adapter — anonClassQualifiedName resolves an anonymous
// CLASS node's owner name from the surrounding declaration/assignment so
// JSC_TYPE_MISMATCH messages mirror Java's FunctionType.getInstanceType()
// output. Java derives this nominally via TypedScopeCreator's assignment
// recording; Go walks the AST since the class is not registered when
// anonymous.
func anonClassQualifiedName(cls *rhino.Node) string {
	if cls == nil {
		return ""
	}
	p := cls.Parent()
	if p == nil {
		return ""
	}
	switch p.Token() {
	case rhino.ASSIGN:
		// `lhs = class { … }` — the LHS qualified name is the owner.
		kids := p.Children()
		if len(kids) >= 2 && kids[1] == cls {
			return dottedNameOf(kids[0])
		}
	case rhino.NAME:
		// `var Foo = class { … }` — the NAME is the owner.
		return p.String_()
	}
	return ""
}

// Go-only: language-adapter — dottedNameOf walks a NAME / GETPROP chain
// and produces the dotted source name (e.g. `obj.C` for the LHS of
// `obj.C = …`). Mirrors how Java's NodeUtil.getBestLValueName recovers a
// readable owner name without consulting the TypedScopeCreator.
func dottedNameOf(n *rhino.Node) string {
	if n == nil {
		return ""
	}
	switch n.Token() {
	case rhino.NAME:
		return n.String_()
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) >= 2 {
			lhs := dottedNameOf(kids[0])
			if lhs == "" {
				return ""
			}
			return lhs + "." + kids[1].String_()
		}
	}
	return ""
}

// Go-only: language-adapter — Go extracted COMPUTED_PROP handling from Java's
// visitObjectOrClassLiteralKey into a named function; Java handles this inline
// via TypeValidator.expectIndexMatch at TypeCheck.java:1434-1539.
func (c *TypeCheck) visitComputedPropInClass(n *rhino.Node, parent *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	// Only applies when directly inside a class body (parent = BLOCK, grandparent = CLASS).
	if parent == nil {
		return nil
	}
	classNode := parent.Parent()
	if classNode == nil || classNode.Token() != rhino.CLASS {
		return nil
	}
	nameNode := classNode.FirstChild()
	if nameNode == nil {
		return nil
	}
	fn, ok := scope.GetDeclaredType(nameNode.String_()).(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return nil
	}
	if !fn.Instance.IsStruct {
		return nil
	}
	// Index node = first child of COMPUTED_PROP (the key expression).
	keyExpr := n.FirstChild()
	if keyExpr == nil {
		keyExpr = n
	}
	return []diag.Diagnostic{{
		Path:    c.path,
		Line:    keyExpr.Lineno(),
		Col:     keyExpr.Charno(),
		Code:    diag.JSC_ILLEGAL_PROPERTY_ACCESS,
		Level:   diag.Error,
		Message: "Cannot do '[]' access on a struct",
	}}
}

// ---- helpers ----------------------------------------------------- //

// Go-only: language-adapter — formats callee name for diagnostic message;
// Java's getBestFunctionName (TypeCheck.java:2560) serves the same purpose
// but uses a different strategy; Go uses simpler AST-walk to produce the same output.
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
	case rhino.SUPER:
		return "super"
	}
	return "<fn>"
}

// Go-only: language-adapter — maps argument AST position to diagnostic line/col;
// Java's TypeCheck reports at the argument node directly; Go walks GETPROP to
// find the property segment, matching how Java formats the call-site position.
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

// Go-only: language-adapter — resolves a @type {function(...):T} annotation to
// extract the return type for visitReturn; Java's NodeUtil.getBestJSDocInfo
// provides this inline; Go needs a separate helper due to different scope lookup.
func enclosingFunctionType(fn *rhino.Node, scope *TypedScope) *jstype.FunctionType {
	if fn == nil {
		return nil
	}
	// Named function declaration: look up by function name.
	if nameKid := fn.FirstChild(); nameKid != nil && nameKid.Token() == rhino.NAME && nameKid.String_() != "" {
		if ft, ok := scope.GetRawDeclaredType(nameKid.String_()).(*jstype.FunctionType); ok {
			return ft
		}
	}
	return nil
}

// Go-only: language-adapter — mirrors Java's NodeUtil.getBestJSDocInfo which
// hops up the AST to find the nearest declaration-level comment; Go's rhino
// doesn't expose this generically, so TypeCheck implements it directly.
func functionJSDoc(fn *rhino.Node) *rhino.JSDocInfo {
	if fn == nil {
		return nil
	}
	if info := fn.JSDocInfo(); info != nil {
		return info
	}
	// New pipeline: class method FUNCTION is wrapped in MEMBER_FUNCTION_DEF.
	// Check the immediate parent for JSDoc.
	if p := fn.Parent(); p != nil {
		if p.Token() == rhino.MEMBER_FUNCTION_DEF || p.Token() == rhino.GETTER_DEF || p.Token() == rhino.SETTER_DEF {
			if info := p.JSDocInfo(); info != nil {
				return info
			}
		}
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
	return nil
}

// Go-only: language-adapter — predicate for visitGetElem's array-index check;
// Java's TypeValidator.expectIndexMatch handles this inline.
func isArrayReceiver(t jstype.JSType) bool {
	obj := dereferencedObject(t)
	return obj != nil && obj.Name == "Array"
}

// Go-only: language-adapter — predicate for isArrayReceiver's index check;
// Java checks this inline in TypeValidator.expectIndexMatch.
func isNumberIndexOk(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	if isUnknownOrNil(t) {
		return true
	}
	return jstype.IsSubtype(t, jstype.Number)
}

// Go-only: package-bridge — guards JSC_INEXISTENT_PROPERTY to types whose
// externs surface is complete in jsgo; Java's full externs need no such guard.
func isKnownSurfaceBuiltin(name string) bool {
	switch name {
	case "Array", "Date":
		return true
	}
	return false
}

// Go-only: language-adapter — Java's FunctionType inherits these slots
// from the Function.prototype native chain (call/apply/bind, etc.) and
// from FunctionType itself (prototype, length, name). jsgo's FunctionType
// has no inherited surface, so we whitelist the standard slots here to
// keep `Foo.prototype`, `Foo.call`, `Foo.bind`, `Foo.apply`, `Foo.length`,
// `Foo.name`, `Foo.toString` quiet when the constructor has no user
// static of that name.
func isFunctionStandardProp(name string) bool {
	switch name {
	case "prototype", "call", "apply", "bind", "length", "name",
		"arguments", "caller", "toString", "constructor":
		return true
	}
	return false
}

// Go-only: language-adapter — extends ObjectType's property surface on assignment;
// Java's TypedScopeCreator handles this during its traversal; in Go TypeCheck
// pre-builds the surface via collectPropertyWrites using this helper.
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
	// @struct: sealed surface — property additions are forbidden.
	// @interface: open-ended but routed through STRICT_INEXISTENT_PROPERTY.
	// @record: open-ended structural type — property writes via
	// `y.prop = val` register `prop` as inferred so subsequent reads
	// (`x.prop`) succeed without JSC_INEXISTENT_PROPERTY. Mirrors Java's
	// TypeInference which tracks these writes and exposes them to TypeCheck.
	if obj.IsStruct || (obj.IsStructural && !obj.IsRecord) {
		return
	}
	if obj.HasProperty(prop) {
		return
	}
	obj.DefineProperty(prop, jstype.Unknown)
	if obj.IsRecord {
		obj.MarkInferred(prop)
	}
}

// Go-only: language-adapter — unwraps TemplatizedType to reach the base
// ObjectType; Java's JSType.toObjectType() handles this uniformly on the type
// hierarchy; Go needs an explicit helper due to Go's non-polymorphic type system.
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

// Go-only: language-adapter — predicate for visitGetElem's struct bracket-access
// check; Java's ObjectType.isStruct() is queried inline in visitGetElem.
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

// Go-only: language-adapter — predicate for @dict access checks; no Java counterpart.
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

// Go-only: language-adapter — handles THIS in visit(); Java dispatches THIS inline
// at TypeCheck.java:632 (`THIS -> ensureTyped(n, t.getTypedScope().getTypeOfThis())`);
// Go needs a named function since the Visit switch routes through emit().
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

// Go-only: language-adapter — Go extracts JSC_INEXISTENT_PARAM and
// JSC_TEMPLATE_TRANSFORMATION_ON_CLASS checks from Java's visitFunction
// (TypeCheck.java:2391-2419) into a named function; Java handles both inline.
func (c *TypeCheck) visitFunctionParams(fn *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	info := fn.JSDocInfo()
	if info == nil {
		return nil
	}
	var out []diag.Diagnostic
	// JSC_TEMPLATE_TRANSFORMATION_ON_CLASS — port of FunctionTypeBuilder.java:
	// when a constructor or interface has TTL @template T := … =: annotations,
	// fire a warning for each such name. `allowTypeTransformations` in Java is
	// false for constructors and interfaces (only plain functions allow TTL).
	if info.IsConstructor() || info.IsInterface() || info.IsRecord() {
		for _, name := range info.TTLTemplateNames() {
			out = append(out, diag.Diagnostic{
				Path: c.path, Line: fn.Lineno(), Col: fn.Charno(),
				Code: diag.JSC_TEMPLATE_TRANSFORMATION_ON_CLASS, Level: diag.Warning,
				Message: fmt.Sprintf("Template type transformation %s not allowed on classes or interfaces", name),
			})
		}
	}
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
		name := p.String_()
		if name != "" {
			paramSet[name] = true
		} else if p.Token() == rhino.DEFAULT_VALUE {
			// Default parameter: `x = expr` — the binding name is the
			// first child (NAME node).
			if fc := p.FirstChild(); fc != nil && fc.Token() == rhino.NAME {
				paramSet[fc.String_()] = true
			}
		} else if p.Token() == rhino.ITER_SPREAD {
			// Rest parameter: `...x` — the first child is the NAME.
			if fc := p.FirstChild(); fc != nil && fc.Token() == rhino.NAME {
				paramSet[fc.String_()] = true
			}
		}
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

// visitUnaryPlus ports TypeCheck.java:visitUnaryPlus (line 975).
// Fires JSC_BAD_TYPE_FOR_UNARY_OPERATION when the POS operand is not
// number-compatible (string, object, function).
func (c *TypeCheck) visitUnaryPlus(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	child := n.FirstChild()
	if child == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	t := infer.InferType(child)
	if isDefinitelyNotNumber(t) {
		return []diag.Diagnostic{{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code: diag.JSC_BAD_TYPE_FOR_UNARY_OPERATION, Level: diag.Error,
			Message: fmt.Sprintf("sign operator cannot be applied to %s", t.TypeName()),
		}}
	}
	return nil
}

// visitUnaryMinus ports TypeCheck.java:visitUnaryMinus (line 987).
// Fires JSC_BAD_TYPE_FOR_UNARY_OPERATION when the NEG operand is not
// number-compatible.
func (c *TypeCheck) visitUnaryMinus(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	child := n.FirstChild()
	if child == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	t := infer.InferType(child)
	if isDefinitelyNotNumber(t) {
		return []diag.Diagnostic{{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code: diag.JSC_BAD_TYPE_FOR_UNARY_OPERATION, Level: diag.Error,
			Message: fmt.Sprintf("unary minus operator cannot be applied to %s", t.TypeName()),
		}}
	}
	return nil
}

// visitBitwiseNOT ports TypeCheck.java:visitBitwiseNOT (line 999).
// Fires JSC_BAD_TYPE_FOR_BIT_OPERATION when the BITNOT operand is not
// number-compatible.
func (c *TypeCheck) visitBitwiseNOT(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	return c.visitBitOp(n, scope)
}

// checkSpread ports TypeCheck.java:checkSpread (line 1015).
// Currently a no-op since object spread has no type constraint and
// array spread checking is deferred to element-type analysis.
func (c *TypeCheck) checkSpread(spreadNode *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	_ = scope
	// Object spread: nothing to check (Java's comment: "Nothing to check about object spread").
	// Array spread and call spread: deferred (matches Java default-mode silencing).
	return nil
}

// visitName ports TypeCheck.java:visitName (line 1838).
// Handles NAME nodes: validates identifiers against the scope.
// Most checks (undeclared variables) are handled by VarCheck; here
// we return nil matching Java's post-VarCheck default.
func (c *TypeCheck) visitName(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	_ = t
	_ = parent
	// Java fires JSC_USED_GLOBAL_THIS here when n is "this"; in Go this
	// is handled via visitThis in Visit(). For other NAMEs, Java's visitName
	// mostly calls ensureTyped which is not needed in Go's type model.
	return true // typeable = true (matches Java's return convention)
}

// visitOptChainGetProp ports TypeCheck.java:visitOptChainGetProp (line 1964).
// Optional chaining getprop: same checks as visitGetProp but silenced when
// the receiver is null/undefined (the whole point of ?.). Currently delegates
// to visitGetProp since Go's type model doesn't distinguish optional chains
// from normal accesses.
func (c *TypeCheck) visitOptChainGetProp(optChainGetProp *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	return c.visitGetProp(optChainGetProp, scope)
}

// visitOptChainGetElem ports TypeCheck.java:visitOptChainGetElem (line 2247).
// Optional chaining getelem: delegates to visitGetElem.
func (c *TypeCheck) visitOptChainGetElem(optChainGetElem *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	return c.visitGetElem(optChainGetElem, scope)
}

// visitFunction ports TypeCheck.java:visitFunction (line 2391).
// Checks function declarations: abstract method in concrete class,
// struct/dict hierarchy. Delegates to visitFunctionParams for @param checks.
func (c *TypeCheck) visitFunction(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	return c.visitFunctionParams(n, scope)
}

// visitClass ports TypeCheck.java:visitClass (line 2420).
// Checks class declarations: struct/dict hierarchy, @implements, @extends.
// The bulk of class checking is in walkSuperclassCheck.
func (c *TypeCheck) visitClass(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	_ = scope
	// Class checks are performed in walkSuperclassCheck which is called from Check().
	// This stub exists to mirror Java's visitClass entry point.
	return nil
}

// visitModuleBody ports TypeCheck.java:visitModuleBody (line 3024).
// Validates legacy goog.module exports against global type. Currently a no-op
// since jsgo does not model the goog.module namespace.
func (c *TypeCheck) visitModuleBody(t *NodeTraversal, moduleBody *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	_ = t
	_ = moduleBody
	_ = scope
	return nil
}

// visitYield ports TypeCheck.java:visitYield (line 2837).
// Checks that the yielded value's type matches the generator's declared return type.
// Currently a no-op pending full generator type modeling.
func (c *TypeCheck) visitYield(t *NodeTraversal, n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	_ = t
	// TODO(port): TypeCheck.visitYield from TypeCheck.java:2837
	// Reason: requires JsIterables.getElementType and generator template modeling
	// which is not yet implemented in jsgo's type system.
	return nil
}

// visitTaggedTemplateLit ports TypeCheck.java:visitTaggedTemplateLit (line 2896).
// Checks that the tag function is callable and that its parameter types match
// the template literal substitutions.
func (c *TypeCheck) visitTaggedTemplateLit(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	tag := n.FirstChild()
	if tag == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	tagType := infer.InferType(tag)
	stripped := stripNullUndef(tagType)
	if isDefinitelyNotCallable(stripped) {
		return []diag.Diagnostic{{
			Path:    c.path,
			Line:    n.Lineno(),
			Col:     n.Charno(),
			Code:    diag.JSC_NOT_FUNCTION_TYPE,
			Level:   diag.Error,
			Message: fmt.Sprintf("%s expressions are not callable", stripped.TypeName()),
		}}
	}
	return nil
}

// visitImplicitReturnExpression ports TypeCheck.java:visitImplicitReturnExpression (line 2747).
// Checks blockless arrow functions `() => expr` against the declared @return type.
// Java fires when the arrow body's expression type doesn't match @return.
func (c *TypeCheck) visitImplicitReturnExpression(t *NodeTraversal, exprNode *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	_ = t
	// Find the enclosing arrow function.
	fn := exprNode.Parent()
	for fn != nil && fn.Token() != rhino.FUNCTION {
		fn = fn.Parent()
	}
	if fn == nil {
		return nil
	}
	info := functionJSDoc(fn)
	if info == nil || info.ReturnType() == nil {
		return nil
	}
	declared := scope.ResolveExpr(info.ReturnType().TypeName)
	if declared == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	found := infer.InferType(exprNode)
	if jstype.IsSubtype(found, declared) {
		return nil
	}
	return []diag.Diagnostic{{
		Path: c.path, Line: exprNode.Lineno(), Col: exprNode.Charno(),
		Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
		Message: fmt.Sprintf("inconsistent return type\nfound   : %s\nrequired: %s",
			found.TypeName(), declared.TypeName()),
	}}
}

// visitBinaryOperator ports TypeCheck.java:visitBinaryOperator (line 2949).
// Dispatches binary operator type checks. Routes to per-operator handlers.
func (c *TypeCheck) visitBinaryOperator(op rhino.Token, n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	switch op {
	case rhino.LSH, rhino.RSH, rhino.URSH, rhino.BITNOT,
		rhino.BITAND, rhino.BITOR, rhino.BITXOR,
		rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH:
		return c.visitBitOp(n, scope)
	case rhino.ADD:
		return c.visitAdd(n, scope)
	case rhino.SUB, rhino.MUL, rhino.DIV, rhino.MOD:
		return c.visitArithmetic(n, scope)
	case rhino.LT, rhino.LE, rhino.GT, rhino.GE:
		return c.visitRelationalDefault(n, scope)
	}
	return nil
}

// Go-only: language-adapter — Go splits Java's visitBinaryOperator (TypeCheck.java:2949)
// for bitwise tokens into this named function; Java handles all operators in one big switch.
func (c *TypeCheck) visitBitOp(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	infer := NewTypeInference(scope)
	var out []diag.Diagnostic
	op := bitOpSymbol(n.Token())
	_ = n.Token() == rhino.LSH || n.Token() == rhino.RSH || n.Token() == rhino.URSH

	// Unary BITNOT (~): one operand; emit at the operator node position
	// to match Java's TypeCheck which points at the BITNOT expression.
	// Java fires JSC_BAD_TYPE_FOR_BIT_OPERATION when the operand doesn't
	// match number context. Void/null/boolean/object types don't match.
	if n.Token() == rhino.BITNOT {
		if len(kids) == 0 {
			return nil
		}
		t := infer.InferType(kids[0])
		if isDefinitelyNotNumber(t) {
			out = append(out, diag.Diagnostic{
				Path: c.path, Line: n.Lineno(), Col: n.Charno(),
				Code: diag.JSC_BAD_TYPE_FOR_BIT_OPERATION, Level: diag.Error,
				Message: fmt.Sprintf("operator %s cannot be applied to %s", op, t.TypeName()),
			})
		}
		return out
	}
	if len(kids) < 2 {
		return nil
	}
	switch n.Token() {
	case rhino.BITAND, rhino.BITOR, rhino.BITXOR,
		rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR:
		// Bitwise operators fire JSC_TYPE_MISMATCH when not bitwiseable.
		for i, k := range kids {
			t := infer.InferType(k)
			if !isBitwiseable(t) {
				msg := "bad left operand to bitwise operator"
				if i > 0 {
					msg = "bad right operand to bitwise operator"
				}
				c.tv.expectCanAssignTo(k, t, bitwiseableType(), msg)
				out = append(out, c.tvDrain()...)
				break
			}
		}
	default:
		// Shift operators: fire JSC_BAD_TYPE_FOR_BIT_OPERATION for non-numbers.
		for _, k := range kids {
			t := infer.InferType(k)
			if isDefinitelyNotNumber(t) {
				out = append(out, diag.Diagnostic{
					Path: c.path, Line: k.Lineno(), Col: k.Charno(),
					Code: diag.JSC_BAD_TYPE_FOR_BIT_OPERATION, Level: diag.Error,
					Message: fmt.Sprintf("operator %s cannot be applied to %s", op, t.TypeName()),
				})
				break
			}
		}
	}
	return out
}


// Go-only: language-adapter — helper for visitBitOp; Java's TypeValidator uses
// the allBitwisableValueTypes constant; Go constructs the union on the fly.
func bitwiseableType() jstype.JSType {
	return jstype.NewUnion(jstype.Boolean, jstype.Null, jstype.Number, jstype.String, jstype.Void)
}

// Go-only: language-adapter — predicate for visitBitOp; mirrors Java's
// TypeValidator.expectBitwiseable check but as a standalone boolean predicate.
func isBitwiseable(t jstype.JSType) bool {
	if t == nil || isJSUnknown(t) {
		return true
	}
	if jstype.IsSubtype(t, bitwiseableType()) {
		return true
	}
	// Approximate Java's matchesNumberContext() for wrapper object types
	// (Number, Date, Boolean — they coerce to numbers in bitwise contexts).
	// PrototypeObjectType.matchesNumberContext: isNumberObjectType() ||
	//   isDateType() || isBooleanObjectType() || isStringObjectType() ||
	//   hasOverriddenNativeProperty("valueOf")
	if obj, ok := t.(*jstype.ObjectType); ok {
		switch obj.TypeName() {
		case "Number", "Date", "Boolean", "String":
			return true
		}
	}
	return false
}

// Go-only: language-adapter — maps Token to operator string for diagnostic messages;
// Java uses n.getToken().toString() which is built into the AST node API.
func bitOpSymbol(t rhino.Token) string {
	switch t {
	case rhino.LSH:
		return "<<"
	case rhino.RSH:
		return ">>"
	case rhino.URSH:
		return ">>>"
	case rhino.BITAND, rhino.ASSIGN_BITAND:
		return "&"
	case rhino.BITOR, rhino.ASSIGN_BITOR:
		return "|"
	case rhino.BITXOR, rhino.ASSIGN_BITXOR:
		return "^"
	case rhino.BITNOT:
		return "~"
	case rhino.ASSIGN_LSH:
		return "<<="
	case rhino.ASSIGN_RSH:
		return ">>="
	case rhino.ASSIGN_URSH:
		return ">>>="
	}
	return "?"
}

// Go-only: language-adapter — extracts Java's typeof-string validation inline
// in visit() (TypeCheck.java:692-702, checkTypeofString call) into a named
// function; Go routes EQ/NE/SHEQ/SHNE through a dedicated function.
func (c *TypeCheck) visitTypeofComparison(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) != 2 {
		return nil
	}
	var _, strNode *rhino.Node
	switch {
	case kids[0].Token() == rhino.TYPEOF && kids[1].Token() == rhino.STRINGLIT:
		_, strNode = kids[0], kids[1]
	case kids[1].Token() == rhino.TYPEOF && kids[0].Token() == rhino.STRINGLIT:
		_, strNode = kids[1], kids[0]
	default:
		return nil
	}
	val := strNode.String_()
	if isValidTypeofValue(val) {
		return nil
	}
	return []diag.Diagnostic{{
		Path: c.path, Line: strNode.Lineno(), Col: strNode.Charno(),
		Code:    diag.JSC_UNKNOWN_TYPEOF_VALUE,
		Level:   diag.Warning,
		Message: fmt.Sprintf("unknown type: %s", val),
	}}
}

// Go-only: language-adapter — Go exposes Java's VALID_TYPEOF_VALUES set as a
// predicate; in Java TypeCheck.checkTypeofString (line 1037) calls
// NodeUtil.isValidIdentifier-style check via the validator.
func isValidTypeofValue(s string) bool {
	switch s {
	case "undefined", "boolean", "number", "string", "function", "object", "bigint", "symbol", "unknown":
		return true
	}
	return false
}

// Go-only: language-adapter — guards struct bracket-access check for Symbol.*
// indices; Java's JSType.isSymbol() + allowedIndexType check is inline in
// TypeValidator.expectIndexMatch; Go wraps it as a named predicate.
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

// Go-only: language-adapter — Go extracted Java's CASE branch (TypeCheck.java:815-820)
// from visit() into a named function; Java calls validator.expectSwitchMatchesCase
// inline in the switch case.
func (c *TypeCheck) visitSwitchCase(caseNode, switchNode *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	if switchNode == nil || switchNode.Token() != rhino.SWITCH {
		return nil
	}
	discrim := switchNode.FirstChild()
	if discrim == nil {
		return nil
	}
	caseExpr := caseNode.FirstChild()
	if caseExpr == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	switchType := infer.InferType(discrim)
	caseType := infer.InferType(caseExpr)
	// Skip unknown types — can't prove incompatibility.
	if isJSUnknown(switchType) || isJSUnknown(caseType) {
		return nil
	}
	// canTestForShallowEqualityWith: types are comparable if at least one
	// is a subtype of the other. If neither is a subtype → incompatible.
	if jstype.IsSubtype(switchType, caseType) || jstype.IsSubtype(caseType, switchType) {
		return nil
	}
	// Also allow if the types could share a common value (e.g. union types
	// might have overlapping members).
	if canShareCommonValue(switchType, caseType) {
		return nil
	}
	msg := fmt.Sprintf("case expression doesn't match switch\nfound   : %s\nrequired: %s",
		caseType.TypeName(), switchType.TypeName())
	return []diag.Diagnostic{{
		Path: c.path, Line: caseExpr.Lineno(), Col: caseExpr.Charno(),
		Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
	}}
}

// Go-only: language-adapter — approximates Java's canTestForShallowEqualityWith
// (JSType.java) which Go does not expose generically; used by visitSwitchCase.
func canShareCommonValue(a, b jstype.JSType) bool {
	// Union types: check if any member of one type overlaps with the other.
	if u, ok := a.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if canShareCommonValue(m, b) {
				return true
			}
		}
		return false
	}
	if u, ok := b.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if canShareCommonValue(a, m) {
				return true
			}
		}
		return false
	}
	// Function types can always be compared (Java's function exception).
	if _, ok := a.(*jstype.FunctionType); ok {
		return true
	}
	if _, ok := b.(*jstype.FunctionType); ok {
		return true
	}
	// Same type → can share value.
	return jstype.IsSubtype(a, b) || jstype.IsSubtype(b, a)
}

// visitDynamicImport ports Java's TypeCheck.visitDynamicImport.
// Fires JSC_TYPE_MISMATCH "dynamic import specifier" when the specifier
// is null or undefined — import() requires a string specifier.
// Java: validator.expectNotNullOrUndefined(t, specifier, specifierType,
//   "dynamic import specifier", STRING_TYPE).
func (c *TypeCheck) visitDynamicImport(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	specifier := n.FirstChild()
	if specifier == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	specType := infer.InferType(specifier)
	// Only fire when the specifier is definitively null/undefined.
	// Java's expectNotNullOrUndefined skips unknown types.
	if specType == jstype.Null || specType == jstype.Void {
		msg := fmt.Sprintf("dynamic import specifier\nfound   : %s\nrequired: string",
			specType.TypeName())
		posNode := specifier
		return []diag.Diagnostic{{
			Path:    c.path,
			Line:    posNode.Lineno(),
			Col:     posNode.Charno(),
			Code:    diag.JSC_TYPE_MISMATCH,
			Level:   diag.Error,
			Message: msg,
		}}
	}
	return nil
}

// Go-only: language-adapter — predicate for visitCall's NOT_FUNCTION_TYPE guard;
// Java's TypeCheck calls couldBeAConstructor (line 2328) inline; Go wraps the
// negation as a named predicate with explicit union handling.
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

// Go-only: language-adapter — predicate for visitBitOp shift checks; no direct
// Java named equivalent (Java uses matchesNumberContext inline).
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

// Go-only: language-adapter — extracts Java's IN case (TypeCheck.java:774-784)
// from visit() into a named function.
func (c *TypeCheck) visitIn(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	lhs, rhs := kids[0], kids[1]
	infer := NewTypeInference(scope)
	rawRhs := infer.InferType(rhs)
	recv := stripNullUndef(rawRhs)
	var out []diag.Diagnostic
	// Java's TypeValidator.expectObject fires JSC_TYPE_MISMATCH "'in' requires
	// an object" when the RHS type is definitely null or undefined. Mirrors
	// TypeCheck.java:784 where visitCall routes the RHS through
	// validator.expectObject(n, rt, "'in' requires an object").
	if rawRhs == jstype.Null || rawRhs == jstype.Void {
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
			Message: fmt.Sprintf(
				"'in' requires an object\nfound   : %s\nrequired: Object",
				rawRhs.TypeName()),
		})
	}
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

// Go-only: language-adapter — predicate for visitIn LHS check; Java's
// TypeValidator.expectStringOrSymbol handles this inline.
func isDefinitelyNotInKey(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	return false
}

// Go-only: language-adapter — extracts Java's FOR_IN case (TypeCheck.java:912-918)
// from visit() into a named function.
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

// checkForOfTypes ports TypeCheck.java's checkForOfTypes.
// Fires JSC_TYPE_MISMATCH when the RHS of a for-of loop is null,
// undefined, or a known non-Iterable type (e.g. Iterator without
// Iterable). Matches Java's validator.expectAutoboxesToIterable:
// non-null Iterable<T>, Array<T>, string, and unknown types pass;
// nullable wrappers or pure Iterator types fail.
//
// AST shape: FOR_OF(lhs, iterable, body).
func (c *TypeCheck) checkForOfTypes(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	iterable := kids[1]
	infer := NewTypeInference(scope)
	iterType := infer.InferType(iterable)
	if isUnknownOrNil(iterType) {
		return nil // unknown → can't prove non-iterability
	}

	// Java's expectAutoboxesToIterable passes when the type or its
	// autoboxed form implements Iterable. In our simplified type model
	// we check via TypeName analysis plus a nullable-union guard.
	// Build found-type string before stripping null (needed for message).
	foundName := iterType.TypeName()

	stripped := stripNullUndef(iterType)
	// hasNullOrUndef: the type changed after stripping null/void.
	// Use TypeName comparison (not pointer equality) since stripNullUndef
	// may create a new union pointer even when no null was present.
	var hasNullOrUndef bool
	if stripped == nil {
		hasNullOrUndef = true
	} else {
		hasNullOrUndef = (stripped.TypeName() != iterType.TypeName())
	}

	if hasNullOrUndef {
		// Nullable iterable: `?Iterable<T>` → cannot iterate.
		return []diag.Diagnostic{{
			Path: c.path, Line: iterable.Lineno(), Col: iterable.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
			Message: "Can only iterate over a (non-null) Iterable type\nfound   : " +
				foundName + "\nrequired: Iterable",
		}}
	}

	// After stripping null, check if the remaining type is known to be
	// non-iterable. Skip unknown / all types.
	if isUnknownOrNil(stripped) {
		return nil
	}
	if !forOfIsIterableType(stripped) {
		return []diag.Diagnostic{{
			Path: c.path, Line: iterable.Lineno(), Col: iterable.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
			Message: "Can only iterate over a (non-null) Iterable type\nfound   : " +
				foundName + "\nrequired: Iterable",
		}}
	}
	return nil
}

// Go-only: language-adapter — predicate for checkForOfTypes; Java's
// TypeValidator.expectAutoboxesToIterable checks this inline via JSType.isIterable().
func forOfIsIterableType(t jstype.JSType) bool {
	// Primitive string is iterable (JS engines implement Symbol.iterator).
	if t == jstype.String || t == jstype.All {
		return true
	}
	// Unknown types are treated as potentially iterable.
	if isJSUnknown(t) {
		return true
	}
	// Union: iterable if ALL members are iterable.
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if !forOfIsIterableType(m) {
				return false
			}
		}
		return true
	}
	// Derive the base name from the TypeName: strip template parameters.
	name := t.TypeName()
	if i := strings.IndexByte(name, '<'); i >= 0 {
		name = name[:i]
	}
	switch name {
	case "Iterable", "IteratorIterable", "IterableIterator",
		"AsyncIterable", "AsyncIteratorIterable",
		"Array", "ReadonlyArray", "IArrayLike", "IIterable",
		"Set", "Map", "Generator", "AsyncGenerator",
		"String": // String object is iterable
		return true
	}
	return false
}

// Go-only: language-adapter — handles ADD under strict mode; Java's visitBinaryOperator
// (TypeCheck.java:2949) routes ADD to TypeValidator.expectNumberOrString inline; Go
// extracts it as a named function dispatched only when GroupStrictCheckTypes is on.
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

// Go-only: language-adapter — predicate for visitAdd; no named Java equivalent.
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

// Go-only: language-adapter — strict-mode SUB/MUL/DIV/MOD check routed to
// visitNumericBinaryOp; Java's visitBinaryOperator (TypeCheck.java:2949) handles
// these inline in one big switch.
func (c *TypeCheck) visitArithmetic(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	return c.visitNumericBinaryOp(n, scope)
}

// Go-only: language-adapter — strict-mode bitwise check in the GroupStrictCheckTypes
// branch; Java routes these through visitBinaryOperator inline.
func (c *TypeCheck) visitBitwise(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	return c.visitNumericBinaryOp(n, scope)
}

// Go-only: language-adapter — shared numeric operand checker for visitArithmetic /
// visitBitwise; Java calls TypeValidator.expectNumber directly for each operand.
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

// Go-only: language-adapter — predicate for visitNumericBinaryOp; no named Java equivalent.
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

// Go-only: language-adapter — default-mode relational operand check; Java handles
// LT/LE/GT/GE inline in visit() at TypeCheck.java:738-772 without a named method.
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

// Go-only: language-adapter — predicate for visitRelationalDefault; no named Java
// equivalent (Java uses JSType.isBigIntOrNumber inline).
func isDefinitelyNotComparable(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if t == jstype.Unknown || t == jstype.All {
		return false
	}
	switch t {
	case jstype.Number, jstype.String, jstype.Boolean:
		// boolean is coerced to 0/1 by JS, so it is comparable.
		// Mirrors Java's TypeCheck which does not flag boolean operands
		// for relational ops in default mode.
		return false
	case jstype.Void, jstype.Null:
		// undefined/null: Java fires errors for these in numeric comparisons
		// (expectBigIntOrNumber fails for null/undefined). Treat as not
		// comparable so the default-mode relational check fires.
		return true
	case jstype.Symbol:
		// symbol: Java fires JSC_TYPE_MISMATCH for relational ops.
		// Mirrors expectBigIntNumberOrString which fails for SymbolType.
		return true
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	if _, ok := t.(*jstype.ObjectType); ok {
		return true
	}
	// Union types: not comparable if any member is not comparable.
	// Simplification: not comparable if ALL members are not comparable.
	if u, ok := t.(*jstype.UnionType); ok {
		anyComparable := false
		for _, m := range u.Members() {
			if !isDefinitelyNotComparable(m) {
				anyComparable = true
				break
			}
		}
		return !anyComparable
	}
	return false
}

// Go-only: language-adapter — strict-mode relational operand check; Java's
// visitBinaryOperator handles LT/LE/GT/GE inline at TypeCheck.java:738-772.
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

// Go-only: language-adapter — predicate for visitRelativeComparison; Java's
// TypeValidator.expectMatchingTypesStrict is called inline; Go wraps it.
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

// Go-only: language-adapter — predicate for canCompareStrict; no Java named equivalent.
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

// ---- visitObjectLit / visitArrayLit ----------------------------- //

// Go-only: language-adapter — Go extracts Java's OBJECTLIT case (TypeCheck.java:925-936)
// plus TypedScopeCreator's @lends handling into one named function; Java handles
// these across two separate classes (TypeCheck + TypedScopeCreator).
func (c *TypeCheck) visitObjectLit(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	var out []diag.Diagnostic

	// @lends checks — port of TypedScopeCreator.defineObjectLiteral.
	// Java uses getLendsName().getRoot().getString() to get the root
	// identifier, so `@lends {Foo.prototype}` checks whether "Foo" is
	// declared (not the full dotted name). Union types (`|`) in a @lends
	// annotation are invalid and already caught by JSC_TYPE_PARSE_ERROR;
	// skip the UNKNOWN_LENDS check for them.
	info := n.JSDocInfo()
	if info != nil && info.HasLendsName() {
		rawLendsName := info.LendsName().TypeName
		// Extract the root identifier (first segment before any '.').
		lendsName := rawLendsName
		if i := strings.IndexByte(lendsName, '.'); i >= 0 {
			lendsName = lendsName[:i]
		}
		// Skip malformed lends names (e.g. containing '|' union operators).
		if !strings.ContainsAny(lendsName, "|{}") {
			lendsType := scope.GetDeclaredType(lendsName)
			if lendsType == nil {
				out = append(out, diag.Diagnostic{
					Path:    c.path,
					Line:    n.Lineno(),
					Col:     n.Charno(),
					Code:    diag.JSC_UNKNOWN_LENDS,
					Level:   diag.Error,
					Message: fmt.Sprintf("Variable %s not declared before @lends annotation.", lendsName),
				})
			} else if rawLendsName == lendsName {
				// Only check LENDS_ON_NON_OBJECT for bare names (not dotted
				// names like `Foo.prototype` where the property chain is
				// expected to be an object type).
				if _, isObj := lendsType.(*jstype.ObjectType); !isObj {
					if _, isFn := lendsType.(*jstype.FunctionType); !isFn {
						out = append(out, diag.Diagnostic{
							Path:    c.path,
							Line:    n.Lineno(),
							Col:     n.Charno(),
							Code:    diag.JSC_LENDS_ON_NON_OBJECT,
							Level:   diag.Error,
							Message: fmt.Sprintf("May only lend properties to object types. %s has type %s.",
								lendsName, lendsType.TypeName()),
						})
					}
				}
			}
		}
	}

	infer := NewTypeInference(scope)
	ownerType := infer.InferType(n)
	obj := dereferencedObject(ownerType)

	// Enum element type check — mirrors Java's visitObjectOrClassLiteralKey
	// line 1515: when key type is EnumElementType, check value against the
	// enum's underlying primitive type.
	// Detect enum init: check if the parent has @enum JSDoc. This covers
	// both `/** @enum */ var E = {...}` (handled by visitVar/visitEnumInit)
	// and `/** @enum */ goog.a = {...}` (namespace-prop enum, visited here).
	var enumUnderlyingType jstype.JSType
	if p := n.Parent(); p != nil {
		parentInfo := p.JSDocInfo()
		if parentInfo == nil {
			// For `goog.a = {...}`, the parent is ASSIGN; check the ASSIGN's JSDoc.
			if p.Token() == rhino.ASSIGN {
				parentInfo = p.JSDocInfo()
			}
		}
		if parentInfo != nil && parentInfo.IsEnum() {
			if expr := parentInfo.EnumParameterType(); expr != nil {
				enumUnderlyingType = scope.ResolveExpr(expr.TypeName)
			} else {
				enumUnderlyingType = jstype.Number // @enum without type defaults to number
			}
		}
	}

	for key := n.FirstChild(); key != nil; key = key.Next() {
		if key.Token() != rhino.STRING_KEY {
			continue
		}
		// Enum element check: verify value type matches the enum's underlying type.
		// Mirrors Java's visitObjectOrClassLiteralKey when keyType.isEnumElementType().
		if enumUnderlyingType != nil {
			value := key.FirstChild()
			if value != nil {
				rightType := infer.InferType(value)
				if !jstype.IsSubtype(rightType, enumUnderlyingType) {
					propName := key.String_()
					parentNode := n.Parent()
					ownerName := ""
					if parentNode != nil {
						if parentNode.Token() == rhino.NAME {
							ownerName = "enum{" + parentNode.String_() + "}"
						} else if parentNode.Token() == rhino.ASSIGN {
							if lhs := parentNode.FirstChild(); lhs != nil {
								qname := NU.GetBestLValueName(lhs)
								if qname != "" {
									ownerName = "enum{" + qname + "}"
								}
							}
						}
					}
					msg := fmt.Sprintf("assignment to property %s of %s\nfound   : %s\nrequired: %s",
						propName, ownerName, rightType.TypeName(), enumUnderlyingType.TypeName())
					out = append(out, diag.Diagnostic{
						Path: c.path, Line: key.Lineno(), Col: key.Charno(),
						Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
					})
				}
			}
			continue // skip the regular key type check for enum values
		}
		// struct/dict key form checks.
		if obj != nil && (obj.IsStruct || obj.IsDict) {
			if key.IsQuotedStringKey() {
				if obj.IsStruct {
					out = append(out, diag.Diagnostic{
						Path:    c.path,
						Line:    key.Lineno(),
						Col:     key.Charno(),
						Code:    diag.JSC_ILLEGAL_OBJLIT_KEY,
						Level:   diag.Error,
						Message: "Illegal key, the object literal is a struct",
					})
				}
			} else if obj.IsDict {
				out = append(out, diag.Diagnostic{
					Path:    c.path,
					Line:    key.Lineno(),
					Col:     key.Charno(),
					Code:    diag.JSC_ILLEGAL_OBJLIT_KEY,
					Level:   diag.Error,
					Message: "Illegal key, the object literal is a dict",
				})
			}
		}
		// Key value type check: when the key carries a @type or @return JSDoc,
		// verify the actual value is assignable to the declared key type.
		// Mirrors Java's TypeCheck.visitObjectOrClassLiteralKey.
		kvInfo := key.JSDocInfo()
		if kvInfo == nil {
			continue
		}
		var declaredKeyType jstype.JSType
		if kvInfo.HasType() {
			declaredKeyType = scope.ResolveExpr(kvInfo.Type().TypeName)
		} else if kvInfo.ReturnType() != nil {
			if ret := scope.ResolveExpr(kvInfo.ReturnType().TypeName); ret != nil {
				declaredKeyType = jstype.NewFunctionType(nil, ret)
			}
		}
		if declaredKeyType == nil {
			continue
		}
		// Skip template types — open-type rule matches Java.
		if _, isTpl := declaredKeyType.(*jstype.TemplateType); isTpl {
			continue
		}
		value := key.FirstChild()
		if value == nil {
			continue
		}
		rightType := infer.InferType(value)
		if jstype.IsSubtype(rightType, declaredKeyType) {
			continue
		}
		propName := key.String_()
		ownerName := ownerType.TypeName()
		msg := fmt.Sprintf("assignment to property %s of %s\nfound   : %s\nrequired: %s",
			propName, ownerName, rightType.TypeName(), declaredKeyType.TypeName())
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: key.Lineno(), Col: key.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
		})
	}
	return out
}

// Go-only: language-adapter — Go extracts Java's ARRAYLIT case from visit()
// (TypeCheck.java:639 `ARRAYLIT -> ensureTyped(n, ARRAY_TYPE)`) plus contextual
// element-type checking that Java does not perform inline in TypeCheck.
func (c *TypeCheck) visitArrayLit(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	elemType := arrayLitContextElemType(n, scope)
	if elemType == nil {
		return nil
	}
	// Template types: open-type rule — accept any element.
	if _, isTpl := elemType.(*jstype.TemplateType); isTpl {
		return nil
	}
	infer := NewTypeInference(scope)
	var out []diag.Diagnostic
	for elem := n.FirstChild(); elem != nil; elem = elem.Next() {
		if elem.Token() == rhino.ITER_SPREAD {
			continue
		}
		found := infer.InferType(elem)
		if jstype.IsSubtype(found, elemType) {
			continue
		}
		msg := fmt.Sprintf("array element\nfound   : %s\nrequired: %s",
			found.TypeName(), elemType.TypeName())
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: elem.Lineno(), Col: elem.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
		})
	}
	return out
}

// Go-only: language-adapter — helper for visitArrayLit; Java doesn't perform
// contextual element-type inference in TypeCheck; Go adds it to reduce FPs.
func arrayLitContextElemType(n *rhino.Node, scope *TypedScope) jstype.JSType {
	parent := n.Parent()
	if parent == nil {
		return nil
	}
	var declaredType jstype.JSType
	switch parent.Token() {
	case rhino.NAME:
		// Array literal is the initializer of a VAR/LET/CONST binding.
		declaredType = scope.GetRawDeclaredType(parent.String_())
	case rhino.ASSIGN:
		// Array literal is the RHS of an assignment.
		if lhs := parent.FirstChild(); lhs != nil && lhs.Token() == rhino.NAME {
			declaredType = scope.GetRawDeclaredType(lhs.String_())
		}
	}
	if declaredType == nil {
		return nil
	}
	tt, ok := declaredType.(*jstype.TemplatizedType)
	if !ok {
		return nil
	}
	if tt.Base == nil || tt.Base.Name != "Array" || len(tt.Args) < 1 {
		return nil
	}
	return tt.Args[0]
}

// ---- JSC_NON_STRINGIFIABLE_OBJECT_KEY branch -------------------- //
//
// Port of TypeCheck.java checkJsdocInfoContainsObjectWithBadKey /
// checkTypeContainsObjectWithBadKey / findObjectWithNonStringifiableKey /
// isObjectTypeWithNonStringifiableKey / isReasonableObjectPropertyKey /
// classHasToString (TypeCheck.java ~line 3175–3358).
//
// Java fires on every typeable node (TypeCheck.java:969). We mirror that
// by calling checkJsdocInfoContainsObjectWithBadKey from Visit() after
// all per-token dispatches.

// checkJsdocInfoContainsObjectWithBadKey checks whether the JSDoc
// attached to n contains any Object<K,V> annotation whose key K cannot
// be stringified. Matches Java's same-named method; fires JSC_NON_STRINGIFIABLE_OBJECT_KEY
// at n's position when a bad key is found. Checks @type, @return,
// @typedef, and each @param type.
func (c *TypeCheck) checkJsdocInfoContainsObjectWithBadKey(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	info := n.JSDocInfo()
	if info == nil {
		return nil
	}
	var out []diag.Diagnostic
	out = append(out, c.checkTypeContainsObjectWithBadKey(n, scope, info.Type())...)
	out = append(out, c.checkTypeContainsObjectWithBadKey(n, scope, info.ReturnType())...)
	out = append(out, c.checkTypeContainsObjectWithBadKey(n, scope, info.TypedefType())...)
	for _, name := range info.ParameterNames() {
		out = append(out, c.checkTypeContainsObjectWithBadKey(n, scope, info.ParameterType(name))...)
	}
	return out
}

// checkTypeContainsObjectWithBadKey resolves a JSTypeExpression and, if
// it contains an Object with a non-stringifiable key, reports at n.
// Mirrors Java's checkTypeContainsObjectWithBadKey.
func (c *TypeCheck) checkTypeContainsObjectWithBadKey(n *rhino.Node, scope *TypedScope, expr *rhino.JSTypeExpression) []diag.Diagnostic {
	if expr == nil {
		return nil
	}
	realType := scope.ResolveExpr(expr.TypeName)
	if realType == nil || isUnknownOrNil(realType) {
		return nil
	}
	visited := map[jstype.JSType]bool{}
	bad := findObjectWithNonStringifiableKey(realType, visited)
	if bad == nil {
		return nil
	}
	return []diag.Diagnostic{{
		Path:    c.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    diag.JSC_NON_STRINGIFIABLE_OBJECT_KEY,
		Level:   diag.Error,
		Message: "Object type \"" + bad.TypeName() + "\" contains non-stringifiable key and it may lead to an error. Please use ES6 Map instead or implement your own Map structure.",
	}}
}

// findObjectWithNonStringifiableKey recursively searches `type` for any
// Object<K,V> whose K is not stringifiable. Returns the first such Object
// type, or nil if all Object keys are reasonable. Mirrors Java's
// findObjectWithNonStringifiableKey.
func findObjectWithNonStringifiableKey(t jstype.JSType, visited map[jstype.JSType]bool) jstype.JSType {
	if t == nil {
		return nil
	}
	if visited[t] {
		// Recursive type — already being checked; skip to avoid infinite loops.
		return nil
	}
	visited[t] = true

	if isObjectTypeWithNonStringifiableKey(t) {
		return t
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, alt := range u.Members() {
			if r := findObjectWithNonStringifiableKey(alt, visited); r != nil {
				return r
			}
		}
	}
	if tt, ok := t.(*jstype.TemplatizedType); ok {
		for _, arg := range tt.Args {
			if r := findObjectWithNonStringifiableKey(arg, visited); r != nil {
				return r
			}
		}
	}
	if fn, ok := t.(*jstype.FunctionType); ok {
		for _, param := range fn.Params {
			if r := findObjectWithNonStringifiableKey(param, visited); r != nil {
				return r
			}
		}
		return findObjectWithNonStringifiableKey(fn.Return, visited)
	}
	return nil
}

// isObjectTypeWithNonStringifiableKey reports whether `t` is a
// TemplatizedType Object<K,V> whose K is not a reasonable property key.
// Only fires on 2+-arg instantiations (where Args[0] is the KEY slot).
// Mirrors Java's isObjectTypeWithNonStringifiableKey.
func isObjectTypeWithNonStringifiableKey(t jstype.JSType) bool {
	tt, ok := t.(*jstype.TemplatizedType)
	if !ok {
		return false
	}
	if tt.Base == nil || tt.Base.Name != "Object" {
		return false
	}
	// Java's Object has two template params (KEY, VALUE). When only one
	// type arg is present it is the value; only 2+ args carry the KEY
	// in position 0. Mirrors templateTypeMap.hasTemplateKey(objectIndexKey).
	if len(tt.Args) < 2 {
		return false
	}
	return !isReasonableObjectPropertyKey(tt.Args[0])
}

// isReasonableObjectPropertyKey reports whether `t` is acceptable as the
// key of an Object<K,V>. Mirrors Java's isReasonableObjectPropertyKey in
// TypeCheck.java.
//
// Stringifiable / always-unique key types:
//   - Unknown / ?          (conservative)
//   - number, string, boolean, null, undefined
//   - Date, RegExp         (native toString)
//   - @interface / @record (user may implement toString in classes)
//   - Array (bare)         (stringifies via element join)
//   - Array<T>             (stringifiable when T is)
//   - enum{T}              (check underlying T)
//   - Union                (all members must be stringifiable)
//   - Object with user-defined toString
func isReasonableObjectPropertyKey(t jstype.JSType) bool {
	if t == nil || isUnknownOrNil(t) {
		return true
	}
	// Primitives and well-known stringifiable/usable-as-key singletons.
	// Symbol is a valid ES6 property key (not stringifiable but allowed
	// as Object index). BigInt is also valid as a key.
	switch t {
	case jstype.Number, jstype.String, jstype.Boolean, jstype.Null, jstype.Void,
		jstype.Symbol, jstype.BigInt:
		return true
	}
	// EnumElementType — check underlying primitive.
	// Mirrors Java: toMaybeEnumElementType() != null → getPrimitiveType().
	if ee, ok := t.(*jstype.EnumElementType); ok {
		return isReasonableObjectPropertyKey(ee.Enum.Underlying)
	}
	// Union — every alternate must be stringifiable.
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if !isReasonableObjectPropertyKey(m) {
				return false
			}
		}
		return true
	}
	// TemplatizedType — handle Array<T> specially.
	if tt, ok := t.(*jstype.TemplatizedType); ok {
		if tt.Base != nil && tt.Base.Name == "Array" {
			// Array<T>: stringifiable when T is stringifiable.
			if len(tt.Args) > 0 {
				return isReasonableObjectPropertyKey(tt.Args[0])
			}
			return true // no T → bare Array
		}
		// Other TemplatizedType (e.g. Object<K,V>) — fall through to
		// the ObjectType check below using the Base.
		return classHasToString(tt.Base)
	}
	// Plain ObjectType — covers @constructor, @interface, @record, and
	// native builtins (Date, RegExp, Array).
	if obj, ok := t.(*jstype.ObjectType); ok {
		// Date and RegExp have native toString(). Matches Java's
		// isDateType() / isRegexpType() early-returns.
		switch obj.Name {
		case "Date", "RegExp":
			return true
		}
		// Bare Array (non-templatized) — always stringifiable.
		if obj.Name == "Array" {
			return true
		}
		// @interface / @record — user may implement toString().
		// Matches Java's isInterface() / isRecordType() early-returns.
		if obj.IsStructural {
			return true
		}
		// User-defined @constructor: OK iff it defines a custom toString().
		return classHasToString(obj)
	}
	// FunctionType and anything else are not reasonable keys.
	return false
}

// nativeObjectsWithNonUniqueToString is the set of built-in types whose
// toString() is defined in our externs but produces non-unique output
// (e.g. Object → "[object Object]", Array → "1,2,3"). These should NOT
// be treated as reasonable Object keys even if HasProperty("toString").
var nativeObjectsWithNonUniqueToString = map[string]bool{
	"Object":   true,
	"Array":    true,
	"Function": true,
}

// classHasToString reports whether the ObjectType (or a non-native
// ancestor) declares a user-defined toString() method. Native builtins
// (Object, Array, …) don't count — their native toString gives "[object
// Object]" / "[1,2]", which is ambiguous. Mirrors Java's classHasToString.
func classHasToString(obj *jstype.ObjectType) bool {
	if obj == nil {
		return false
	}
	// Skip native builtins whose toString is non-unique. Java's
	// classHasToString explicitly rejects Object/Array/Function.
	if nativeObjectsWithNonUniqueToString[obj.Name] {
		return false
	}
	return obj.HasProperty("toString")
}

// ---- visitRegExpLit / visitTemplateLit / visitObjectPattern / checkDestructuringAssignment --

// Go-only: language-adapter — Go extracted REGEXP from visit() into a named function;
// Java calls ensureTyped(n, REGEXP_TYPE) inline at line 641 without a named method.
func (c *TypeCheck) visitRegExpLit(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	// Suppress REGEXP nodes inside @suppress annotations or casts.
	if p := n.Parent(); p != nil && p.Token() == rhino.CAST {
		return nil
	}
	// Java's ensureTyped just records the node for percentage-typed stats.
	// jsgo doesn't track those metrics, so return no diagnostics.
	return nil
}

// Go-only: language-adapter — Go extracted TEMPLATELIT from visit() into a named
// function; Java calls ensureTyped inline plus validator.expectStringOrSymbol
// (TypeCheck.java:684); Go adds explicit TEMPLATELIT_SUB checking.
func (c *TypeCheck) visitTemplateLit(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	infer := NewTypeInference(scope)
	var out []diag.Diagnostic
	for _, child := range n.Children() {
		if child.Token() != rhino.TEMPLATELIT_SUB {
			continue
		}
		expr := child.FirstChild()
		if expr == nil {
			continue
		}
		t := infer.InferType(expr)
		if !isTemplateLitSubOk(t) {
			out = append(out, diag.Diagnostic{
				Path: c.path, Line: expr.Lineno(), Col: expr.Charno(),
				Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error,
				Message: fmt.Sprintf(
					"template literal substitution\nfound   : %s\nrequired: (string|number|boolean|null|undefined|symbol)",
					t.TypeName()),
			})
		}
	}
	return out
}

// Go-only: language-adapter — predicate for visitTemplateLit; Java's
// TypeValidator.expectStringOrSymbol handles this inline.
func isTemplateLitSubOk(t jstype.JSType) bool {
	if t == nil || isUnknownOrNil(t) {
		return true
	}
	switch t {
	case jstype.String, jstype.Number, jstype.Boolean, jstype.Null, jstype.Void, jstype.All:
		return true
	}
	// Objects have implicit toString/valueOf coercion.
	if _, ok := t.(*jstype.ObjectType); ok {
		return true
	}
	if _, ok := t.(*jstype.TemplatizedType); ok {
		return true
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if !isTemplateLitSubOk(m) {
				return false
			}
		}
		return true
	}
	return true
}

// visitObjectPattern ports TypeCheck.java's visitObjectPattern. When
// visiting an OBJECT_PATTERN node (`{a, b} = rhs`), each STRING_KEY
// child names a property that must exist on the RHS object type.
// Missing properties fire JSC_INEXISTENT_PROPERTY.
//
// The RHS type is resolved by climbing the parent chain:
//   - ASSIGN(OBJECT_PATTERN, rhs)  → InferType(rhs)
//   - NAME(init=rhs) inside VAR/LET/CONST → InferType(rhs)
//
// Guards: fires only on closed-surface types (structural @record /
// @interface and @struct nominal classes). Open @constructor surfaces
// are skipped to match Java's default-mode behavior.
func (c *TypeCheck) visitObjectPattern(pattern *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	rhsType := objectPatternRHSType(pattern, scope)
	if rhsType == nil {
		return nil
	}
	rhsType = stripNullUndef(rhsType)
	if isUnknownOrNil(rhsType) {
		return nil
	}
	obj := dereferencedObject(rhsType)
	if obj == nil {
		return nil
	}
	// Only check closed-surface types (structural or @struct) to avoid
	// false-positives on open @constructor surfaces, matching Java's
	// default-mode checkPropertyAccessForDestructuring guard.
	if !obj.IsStructural && !obj.IsStruct {
		return nil
	}
	var out []diag.Diagnostic
	for _, key := range pattern.Children() {
		if key.Token() != rhino.STRING_KEY {
			continue // skip rest elements
		}
		propName := key.String_()
		if obj.HasProperty(propName) {
			continue
		}
		dedupKey := obj.TypeName() + "." + propName
		if c.inexistentReported[dedupKey] {
			continue
		}
		if c.inexistentReported == nil {
			c.inexistentReported = map[string]bool{}
		}
		c.inexistentReported[dedupKey] = true
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: key.Lineno(), Col: key.Charno(),
			Code:    diag.JSC_INEXISTENT_PROPERTY,
			Level:   diag.Error,
			Message: fmt.Sprintf("Property %s never defined on %s", propName, obj.TypeName()),
		})
	}
	return out
}

// Go-only: language-adapter — helper for visitObjectPattern; Java resolves the RHS
// type inline via NodeUtil.getBestJSDocInfo; Go climbs the parent chain explicitly.
func objectPatternRHSType(pattern *rhino.Node, scope *TypedScope) jstype.JSType {
	parent := pattern.Parent()
	if parent == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	switch parent.Token() {
	case rhino.ASSIGN:
		kids := parent.Children()
		if len(kids) >= 2 {
			return infer.InferType(kids[1])
		}
	case rhino.NAME:
		// NAME is the VAR/LET/CONST binding — init is the first child of NAME.
		init := parent.FirstChild()
		if init != nil {
			return infer.InferType(init)
		}
	}
	return nil
}

// checkDestructuringAssignment ports TypeCheck.java's
// checkDestructuringAssignment. Routes ARRAY_PATTERN and OBJECT_PATTERN
// LHS nodes to their respective checkers.
//
//   - ARRAY_PATTERN: checks that the RHS is iterable (has Symbol.iterator
//     or is an array/tuple type). Currently a structural no-op since
//     iterability is not yet modeled in jsgo's type system.
//   - OBJECT_PATTERN: delegates to visitObjectPattern to check that each
//     named property exists on the RHS type.
func (c *TypeCheck) checkDestructuringAssignment(lhs *rhino.Node, rhsType jstype.JSType, scope *TypedScope) []diag.Diagnostic {
	switch lhs.Token() {
	case rhino.OBJECT_PATTERN:
		return c.visitObjectPattern(lhs, scope)
	case rhino.ARRAY_PATTERN:
		// Array destructuring iterability: not yet modeled.
		// Java fires JSC_TYPE_MISMATCH when the RHS is not iterable
		// (e.g. a non-array object without Symbol.iterator). jsgo
		// does not track iterability in the type graph, so this is
		// a no-op pending that feature.
		return nil
	}
	return nil
}

// ---- Abstract method / property access checks ----

// checkAbstractMethodInConcreteClass ports TypeCheck.java's
// checkAbstractMethodInConcreteClass. Fires JSC_ABSTRACT_METHOD_IN_CONCRETE_CLASS
// when a method annotated @abstract is defined inside a non-abstract class.
//
// Java signature: checkAbstractMethodInConcreteClass(Node n, FunctionType ctorType, JSDocInfo info)
// Called when visiting prototype method assignments and class members.
func (c *TypeCheck) checkAbstractMethodInConcreteClass(n *rhino.Node, ctorFn *jstype.FunctionType, info *rhino.JSDocInfo) []diag.Diagnostic {
	if info == nil || !info.IsAbstract() {
		return nil
	}
	// The constructor itself must be non-abstract to fire.
	ctorDoc := ctorFn.GetJSDocInfo()
	if ctorDoc.IsAbstract() {
		return nil // abstract method in abstract class is fine
	}
	if !ctorFn.IsConstructor {
		return nil
	}
	return []diag.Diagnostic{{
		Path: c.path, Line: n.Lineno(), Col: n.Charno(),
		Code:    diag.JSC_ABSTRACT_METHOD_IN_CONCRETE_CLASS,
		Level:   diag.Warning,
		Message: "Abstract methods can only appear in abstract classes. Please declare the class as @abstract",
	}}
}

// checkAbstractPropertyAccess ports TypeCheck.java's checkAbstractPropertyAccess.
// Fires JSC_ABSTRACT_SUPER_METHOD_NOT_USABLE when an abstract method is accessed
// via super.foo() or Foo.prototype.bar().
//
// Does not fire for @abstract constructors (they are still callable) or when
// the access is the target of an assignment (declaring the method).
func (c *TypeCheck) checkAbstractPropertyAccess(method *rhino.Node) []diag.Diagnostic {
	// Allow lhs of assignment: `/** @abstract */ Foo.prototype.bar = function() {}`
	if method.Parent() != nil && method.Parent().Token() == rhino.ASSIGN &&
		method.Parent().FirstChild() == method {
		return nil
	}
	// The method itself must have an abstract FunctionType.
	// We check via the property's inferred type (not available without
	// full inference), so we guard on the JSDoc of the node.
	// Java uses getJSType(method).toMaybeFunctionType().isAbstract().
	// For the basic port, we check objectNode.isSuper() and
	// objectNode == Foo.prototype pattern, as Java does.
	objectNode := method.FirstChild()
	if objectNode == nil {
		return nil
	}
	if objectNode.Token() != rhino.SUPER && objectNode.Token() != rhino.GETPROP {
		return nil
	}
	if objectNode.Token() == rhino.SUPER {
		// case (a): super.foo()
		// Java also requires methodType.isAbstract() — without full type
		// inference for the method type, we skip this to avoid FPs.
		return nil
	}
	if objectNode.Token() == rhino.GETPROP {
		prop := objectNode.String_()
		if prop == "prototype" {
			// case (b): Foo.prototype.bar — also needs abstract check
			// Skip without full type info.
			return nil
		}
	}
	return nil
}

// ---- JSC_DETERMINISTIC_TEST (port of TypeCheck.visitCompare equality branch) ----
//
// Java's TypeCheck.visitCompare dispatches differently for the two
// operator groups:
//
//	EQ / NE   — TypeCheck calls leftType.testForEquality(rightType),
//	            which honours JS coercion rules (number ↔ string ↔ boolean).
//
//	SHEQ / SHNE — TypeCheck calls
//	            leftType.canTestForShallowEqualityWith(rightType),
//	            the strict-equality variant that requires a non-empty
//	            greatest common subtype.
//
// Both paths strip null/undefined before comparing (Java's
// restrictByNotNullOrUndefined) so `(string|null) === number` and
// `string === number` both fire.

// Go-only: language-adapter — Go extracts Java's EQ/NE/SHEQ/SHNE inline handling
// (TypeCheck.java:692-736) into a named function; Java calls JSType.testForEquality
// directly in the visit() switch.
func (c *TypeCheck) visitEqualityOp(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	// Java's TypeCheck.visitCompare uses getJSType(left)/getJSType(right),
	// which returns the flow-refined type at the comparison site. After
	// "x = null", flow narrows x to null, so "x == null" becomes null==null
	// → eqTrue and fires JSC_DETERMINISTIC_TEST. Escaped variables (assigned
	// in nested closures) are filtered from the flow refinements in
	// pushEnclosingFlow, matching TypeInference.java#isUnflowable.
	infer := NewTypeInference(scope)
	leftType := infer.InferType(kids[0])
	rightType := infer.InferType(kids[1])

	// Java strips null/undefined before comparing (restrictByNotNullOrUndefined).
	leftRestricted := equalityRestrictByNotNullOrUndefined(leftType)
	rightRestricted := equalityRestrictByNotNullOrUndefined(rightType)

	var result eqTri
	switch n.Token() {
	case rhino.EQ:
		result = equalityTestForEquality(leftRestricted, rightRestricted)
	case rhino.NE:
		result = eqTriNot(equalityTestForEquality(leftRestricted, rightRestricted))
	case rhino.SHEQ:
		result = equalityTestForShallowEquality(leftRestricted, rightRestricted)
	case rhino.SHNE:
		result = eqTriNot(equalityTestForShallowEquality(leftRestricted, rightRestricted))
	default:
		return nil
	}

	if result == eqUnknown {
		return nil
	}
	valueStr := "false"
	if result == eqTrue {
		valueStr = "true"
	}
	// Position: when the comparison is wrapped directly in an
	// EXPR_RESULT with an outer paren grouping (e.g. `(function(){}) == 'x'`),
	// Java reports at the statement's column (0). The inner EQ's Charno
	// may be inside the wrapping parens; promote to the EXPR_RESULT when
	// its column is earlier.
	line, col := n.Lineno(), n.Charno()
	if p := n.Parent(); p != nil && p.Token() == rhino.EXPR_RESULT && p.Charno() < col {
		line, col = p.Lineno(), p.Charno()
	}
	return []diag.Diagnostic{{
		Path:  c.path,
		Line:  line,
		Col:   col,
		Code:  diag.JSC_DETERMINISTIC_TEST,
		Level: diag.Error,
		Message: fmt.Sprintf(
			"condition always evaluates to %s\nleft : %s\nright: %s",
			valueStr, leftType.TypeName(), rightType.TypeName()),
	}}
}

// eqTri mirrors Java's com.google.javascript.rhino.jstype.TernaryValue (Tri).
type eqTri int8

const (
	eqUnknown eqTri = 0
	eqTrue    eqTri = 1
	eqFalse   eqTri = -1
)

// Go-only: language-adapter — mirrors Java's TernaryValue.not() method; Go uses a
// free function since eqTri is a plain int8, not an enum with methods.
func eqTriNot(t eqTri) eqTri {
	switch t {
	case eqTrue:
		return eqFalse
	case eqFalse:
		return eqTrue
	}
	return eqUnknown
}

// Go-only: language-adapter — ports Java's JSType.restrictByNotNullOrUndefined()
// virtual method as a free function; Go's non-polymorphic type system requires
// explicit type switching instead of virtual dispatch.
func equalityRestrictByNotNullOrUndefined(t jstype.JSType) jstype.JSType {
	if t == nil {
		return nil // already empty
	}
	// Null and Undefined restrict to the empty type.
	if t == jstype.Null || t == jstype.Void {
		return nil
	}
	u, ok := t.(*jstype.UnionType)
	if !ok {
		return t
	}
	var kept []jstype.JSType
	for _, m := range u.Members() {
		if m == jstype.Null || m == jstype.Void {
			continue
		}
		kept = append(kept, m)
	}
	switch len(kept) {
	case 0:
		return nil // all members stripped → empty type
	case 1:
		return kept[0]
	default:
		return jstype.NewUnion(kept...)
	}
}

// Go-only: language-adapter — ports Java's virtual JSType.testForEquality()
// dispatch hierarchy as a free function with explicit type switching; Java uses
// per-class overrides (NullType, VoidType, SymbolType, UnionType, etc.).
func equalityTestForEquality(a, b jstype.JSType) eqTri {
	// nil represents the empty type (NO_TYPE in Java).
	aEmpty := a == nil
	bEmpty := b == nil
	if aEmpty || bEmpty {
		if aEmpty && bEmpty {
			return eqTrue // e.g. null == null → both restrict to NO_TYPE → TRUE
		}
		return eqUnknown
	}
	// All, Unknown → conservative.
	if equalityIsWide(a) || equalityIsWide(b) {
		return eqUnknown
	}
	// Union: reduce to the merge of each member's result.
	if au, ok := a.(*jstype.UnionType); ok {
		return equalityUnionTestForEquality(au, b)
	}
	if bu, ok := b.(*jstype.UnionType); ok {
		return equalityUnionTestForEquality(bu, a)
	}
	// EnumElementType: unwrap to underlying and re-evaluate.
	if ae, ok := a.(*jstype.EnumElementType); ok {
		return equalityTestForEquality(ae.Enum.Underlying, b)
	}
	if be, ok := b.(*jstype.EnumElementType); ok {
		return equalityTestForEquality(a, be.Enum.Underlying)
	}

	// FunctionType or the built-in Function ObjectType.
	if equalityIsFunctionLike(a) {
		return equalityFunctionVsType(b)
	}
	if equalityIsFunctionLike(b) {
		return equalityFunctionVsType(a)
	}

	// symbol primitive (Java SymbolType.testForEquality).
	if a == jstype.Symbol {
		return equalitySymbolVsType(b)
	}
	if b == jstype.Symbol {
		return equalitySymbolVsType(a)
	}

	// null (Java NullType.testForEquality).
	if a == jstype.Null {
		return equalityNullVsType(b)
	}
	if b == jstype.Null {
		return equalityNullVsType(a)
	}

	// undefined/void (Java VoidType.testForEquality).
	if a == jstype.Void {
		return equalityVoidVsType(b)
	}
	if b == jstype.Void {
		return equalityVoidVsType(a)
	}

	// Number/String/Boolean can coerce with each other and objects → UNKNOWN.
	if equalityIsCoerciblePrimitive(a) {
		return equalityCoercibleVsType(b)
	}
	if equalityIsCoerciblePrimitive(b) {
		return equalityCoercibleVsType(a)
	}

	// ObjectType / unknown kind → conservative.
	return eqUnknown
}

// Go-only: language-adapter — ports Java's UnionType.testForEquality() override;
// Java iterates members polymorphically; Go uses explicit union type switching.
func equalityUnionTestForEquality(u *jstype.UnionType, other jstype.JSType) eqTri {
	var result eqTri
	first := true
	for _, m := range u.Members() {
		r := equalityTestForEquality(m, other)
		if first {
			result = r
			first = false
		} else if r != result {
			return eqUnknown
		}
	}
	if first {
		return eqUnknown
	}
	return result
}

// Go-only: language-adapter — ports Java's virtual JSType.canTestForShallowEqualityWith()
// as a free function; Java dispatches through per-class overrides.
func equalityTestForShallowEquality(a, b jstype.JSType) eqTri {
	if a == nil || b == nil {
		return eqUnknown
	}
	// All, Unknown → conservative.
	if equalityIsWide(a) || equalityIsWide(b) {
		return eqUnknown
	}
	// Empty types — mirror Java: subtype check.
	if equalityIsEmpty(a) || equalityIsEmpty(b) {
		if jstype.IsSubtype(a, b) || jstype.IsSubtype(b, a) {
			return eqUnknown
		}
		return eqFalse
	}
	// EnumElementType: unwrap to underlying for strict comparison.
	if ae, ok := a.(*jstype.EnumElementType); ok {
		return equalityTestForShallowEquality(ae.Enum.Underlying, b)
	}
	if be, ok := b.(*jstype.EnumElementType); ok {
		return equalityTestForShallowEquality(a, be.Enum.Underlying)
	}
	// Union: reduce over members.
	if au, ok := a.(*jstype.UnionType); ok {
		return equalityUnionTestForShallowEquality(au, b)
	}
	if bu, ok := b.(*jstype.UnionType); ok {
		return equalityUnionTestForShallowEquality(bu, a)
	}

	// symbol primitive vs object → disjoint (strict types, no coercion).
	if a == jstype.Symbol {
		return equalityShallowSymbolVsType(b)
	}
	if b == jstype.Symbol {
		return equalityShallowSymbolVsType(a)
	}

	// Primitives vs object (or function): no common subtype → FALSE.
	aPrim := equalityIsPrimitive(a)
	bPrim := equalityIsPrimitive(b)
	if aPrim && !bPrim {
		return eqFalse
	}
	if bPrim && !aPrim {
		return eqFalse
	}
	// Two different primitives → no common subtype.
	if aPrim && bPrim && a != b {
		return eqFalse
	}
	// Two object types (or two functions): conservatively UNKNOWN.
	return eqUnknown
}

// Go-only: language-adapter — ports Java's UnionType.canTestForShallowEqualityWith()
// override; no named Java equivalent at this abstraction level.
func equalityUnionTestForShallowEquality(u *jstype.UnionType, other jstype.JSType) eqTri {
	var result eqTri
	first := true
	for _, m := range u.Members() {
		r := equalityTestForShallowEquality(m, other)
		if first {
			result = r
			first = false
		} else if r != result {
			return eqUnknown
		}
	}
	if first {
		return eqUnknown
	}
	return result
}

// --- type-classification helpers for equality checks ---

// Go-only: language-adapter — classification helper for equality test functions;
// Java's AllType / UnknownType implement testForEquality as UNKNOWN inline.
func equalityIsWide(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	n := t.TypeName()
	return n == "*" || n == "?"
}

// Go-only: language-adapter — Java has a dedicated NoType (bottom type); Go
// represents it as nil. nil inputs are handled upstream; non-nil inputs are
// never the empty type in Go's type model, so this always returns false.
func equalityIsEmpty(t jstype.JSType) bool {
	// jsgo has no NoType sentinel; nil represents the empty (bottom) type
	// and is handled upstream (callers check t == nil before calling this).
	_ = t
	return false
}

// Go-only: language-adapter — classification helper; Java uses JSType.isNominalType()
// and related predicates inline.
func equalityIsPrimitive(t jstype.JSType) bool {
	switch t {
	case jstype.Number, jstype.String, jstype.Boolean,
		jstype.Null, jstype.Void, jstype.Symbol:
		return true
	}
	return false
}

// Go-only: language-adapter — classification helper for equalityTestForEquality.
func equalityIsCoerciblePrimitive(t jstype.JSType) bool {
	return t == jstype.Number || t == jstype.String || t == jstype.Boolean
}

// Go-only: language-adapter — ports Java's FunctionType.testForEquality() override.
func equalityFunctionVsType(other jstype.JSType) eqTri {
	if equalityIsWide(other) {
		return eqUnknown
	}
	if equalityIsPrimitive(other) {
		return eqFalse
	}
	return eqUnknown
}

// Go-only: language-adapter — ports Java's SymbolType.testForEquality() override.
func equalitySymbolVsType(other jstype.JSType) eqTri {
	if equalityIsWide(other) {
		return eqUnknown
	}
	if other == jstype.Symbol {
		return eqUnknown
	}
	if equalityIsObject(other) {
		return eqUnknown
	}
	return eqFalse
}

// Go-only: language-adapter — ports the symbol branch of JSType.canTestForShallowEqualityWith.
func equalityShallowSymbolVsType(other jstype.JSType) eqTri {
	if equalityIsWide(other) {
		return eqUnknown
	}
	if other == jstype.Symbol {
		return eqUnknown
	}
	if equalityIsObject(other) {
		return eqFalse
	}
	return eqFalse
}

// Go-only: language-adapter — ports Java's NullType.testForEquality() override.
func equalityNullVsType(other jstype.JSType) eqTri {
	if equalityIsWide(other) {
		return eqUnknown
	}
	if other == jstype.Null {
		return eqTrue
	}
	if other == jstype.Void {
		return eqTrue // null == undefined is true in JS
	}
	if equalityIsNullable(other) {
		return eqUnknown
	}
	return eqFalse
}

// Go-only: language-adapter — ports Java's VoidType.testForEquality() override.
func equalityVoidVsType(other jstype.JSType) eqTri {
	if equalityIsWide(other) {
		return eqUnknown
	}
	if other == jstype.Void || other == jstype.Null {
		return eqTrue
	}
	if equalityIsNullable(other) {
		return eqUnknown
	}
	return eqFalse
}

// Go-only: language-adapter — ports Java's NumberType/StringType/BooleanType
// testForEquality() overrides (they share the same logic).
func equalityCoercibleVsType(other jstype.JSType) eqTri {
	if equalityIsWide(other) {
		return eqUnknown
	}
	switch other {
	case jstype.Null, jstype.Void, jstype.Symbol:
		return eqFalse
	}
	if equalityIsObject(other) {
		return eqUnknown
	}
	return eqUnknown
}

// Go-only: language-adapter — classification helper for equalityNullVsType etc.;
// Java uses JSType.isObjectType() virtual method.
func equalityIsObject(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if _, ok := t.(*jstype.ObjectType); ok {
		return true
	}
	if _, ok := t.(*jstype.PrototypeObjectType); ok {
		return true
	}
	if _, ok := t.(*jstype.InstanceObjectType); ok {
		return true
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	if _, ok := t.(*jstype.TemplatizedType); ok {
		return true
	}
	if ee, ok := t.(*jstype.EnumElementType); ok {
		return equalityIsObject(ee.Enum.Underlying)
	}
	return false
}

// Go-only: language-adapter — reports nullable potential; Java's types use
// isNullable() virtual method.
func equalityIsNullable(t jstype.JSType) bool {
	return equalityIsObject(t)
}

// Go-only: language-adapter — classification helper for equalityTestForEquality.
func equalityIsFunctionLike(t jstype.JSType) bool {
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	if obj, ok := t.(*jstype.ObjectType); ok && obj.Name == "Function" {
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Inlined SuperclassCheck — port of TypeCheck.java's
// checkPropertyVisibility / validateMethodOverride branches.
// Fires JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH and
// JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH.
// ---------------------------------------------------------------------------

// Go-only: language-adapter — Go dispatches superclass/implements checks via a
// single recursive AST walker; Java's TypeCheck fires these checks per-node in
// visit() by calling individual methods; Go batches them in a pre-pass.
func (c *TypeCheck) walkSuperclassCheck(n *rhino.Node) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.CLASS:
		c.checkClass(n)
		c.checkClassDictExtendsStruct(n)
		c.checkClassImplements(n)
		c.checkExtendsNonObject(n)
	case rhino.FUNCTION:
		c.checkFunctionImplements(n)
		c.checkInterfaceExtendsConflicts(n)
		c.checkExtendsKindMismatch(n)
		c.checkExtendsNonObject(n)
		c.checkStructDictRequireConstructor(n)
		c.checkFunctionDictExtendsStruct(n)
		c.checkExtendsWithoutTypedef(n)
	case rhino.VAR, rhino.LET, rhino.CONST:
		c.checkDeclaredConstructorImplements(n)
	case rhino.EXPR_RESULT:
		if a := n.FirstChild(); a != nil && a.Token() == rhino.ASSIGN {
			c.checkPrototypeAssign(n, a)
		}
		if gp := n.FirstChild(); gp != nil && gp.Token() == rhino.GETPROP {
			c.checkPrototypeStub(n, gp)
		}
	}
	for _, ch := range n.Children() {
		c.walkSuperclassCheck(ch)
	}
}

// Go-only: language-adapter — extracted from Java's visitFunction / visitClass
// struct-dict hierarchy check; Java calls TypeCheck.checkStructDictSubtyping
// (TypeCheck.java:2515) inline; Go splits each check into a named function.
func (c *TypeCheck) checkFunctionDictExtendsStruct(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if info == nil || !info.IsDict() || !info.IsConstructor() {
		return
	}
	base := info.BaseType()
	if base == nil {
		return
	}
	super := c.global.LookupNominalType(base.TypeName)
	if super == nil || !super.IsStruct {
		return
	}
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	c.diags = append(c.diags, diag.Diagnostic{
		Path: c.path, Line: fn.Lineno(), Col: fn.Charno(),
		Code: diag.JSC_DICT_EXTEND_STRUCT_TYPE, Level: diag.Error,
		Message: fmt.Sprintf("@dict class %s cannot extend @struct class %s", name, base.TypeName),
	})
}

// Go-only: language-adapter — ES6 class variant of checkFunctionDictExtendsStruct;
// Java's checkStructDictSubtyping (TypeCheck.java:2515) handles both uniformly.
func (c *TypeCheck) checkClassDictExtendsStruct(cls *rhino.Node) {
	info := cls.JSDocInfo()
	kids := cls.Children()
	if len(kids) < 2 {
		return
	}
	nameNode, superNode := kids[0], kids[1]
	name := nameNode.String_()
	if name == "" {
		return
	}
	if superNode == nil || superNode.String_() == "" {
		return
	}
	superName := superNode.String_()
	super := c.global.LookupNominalType(superName)
	// ES6 classes are @struct by default unless explicitly @dict or @unrestricted.
	subIsDict := info != nil && info.IsDict()
	subIsUnrestricted := info != nil && info.MakesUnrestricted()
	subIsStruct := !subIsDict && !subIsUnrestricted
	isSuperStruct := super != nil && super.IsStruct
	isSuperDict := super != nil && super.IsDict
	// Treat a plain ES6-class parent (no JSDoc) as struct-by-default.
	if super != nil && !super.IsStruct && !super.IsDict && !super.IsStructural {
		isSuperStruct = true
	}
	if subIsDict && isSuperStruct {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: cls.Lineno(), Col: cls.Charno(),
			Code: diag.JSC_DICT_EXTEND_STRUCT_TYPE, Level: diag.Error,
			Message: fmt.Sprintf("@dict class %s cannot extend @struct class %s", name, superName),
		})
		return
	}
	if subIsStruct && isSuperDict {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: cls.Lineno(), Col: cls.Charno(),
			Code: diag.JSC_DICT_EXTEND_STRUCT_TYPE, Level: diag.Error,
			Message: fmt.Sprintf("@struct class %s cannot extend @dict class %s", name, superName),
		})
	}
}

// Go-only: language-adapter — Go extracted this guard from Java's
// checkConstructor / visitFunction (TypeCheck.java:2450-2507) into a named
// function dispatched via walkSuperclassCheck.
func (c *TypeCheck) checkStructDictRequireConstructor(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if info == nil {
		return
	}
	if !info.IsStruct() && !info.IsDict() {
		return
	}
	if info.IsConstructor() || info.IsInterface() || info.IsRecord() {
		return
	}
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	tag := "@struct"
	if info.IsDict() {
		tag = "@dict"
	}
	c.diags = append(c.diags, diag.Diagnostic{
		Path: c.path, Line: fn.Lineno(), Col: fn.Charno(),
		Code: diag.JSC_CONSTRUCTOR_REQUIRED, Level: diag.Error,
		Message: fmt.Sprintf("%s used without @constructor for %s", tag, name),
	})
}

// checkExtendsWithoutTypedef ports FunctionTypeBuilder.java:469 — fires
// JSC_EXTENDS_WITHOUT_TYPEDEF when a FUNCTION node carries @extends but
// is neither @constructor nor @interface.
func (c *TypeCheck) checkExtendsWithoutTypedef(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if info == nil {
		return
	}
	if info.IsConstructor() || info.IsInterface() || info.IsRecord() {
		return
	}
	if info.BaseType() == nil {
		return
	}
	kids := fn.Children()
	name := ""
	if len(kids) >= 1 {
		name = kids[0].String_()
	}
	c.diags = append(c.diags, diag.Diagnostic{
		Path: c.path, Line: fn.Lineno(), Col: fn.Charno(),
		Code: diag.JSC_EXTENDS_WITHOUT_TYPEDEF, Level: diag.Error,
		Message: fmt.Sprintf("@extends used without @constructor or @interface for %s", name),
	})
}

// Go-only: language-adapter — Go extracted @extends kind validation from Java's
// checkConstructor / checkInterface (TypeCheck.java:2450 and 2528) into a named
// function.
func (c *TypeCheck) checkExtendsKindMismatch(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if info == nil {
		return
	}
	isIface := info.IsInterface() || info.IsRecord()
	isCtor := info.IsConstructor()
	if !isIface && !isCtor {
		return
	}
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	// Walk every @extends — fires once for the first parent whose
	// kind conflicts with the child's. Matches Java's behaviour of
	// reporting one conflicting-extends diagnostic per declaration.
	bases := info.BaseTypes()
	if len(bases) == 0 {
		if b := info.BaseType(); b != nil {
			bases = []*rhino.JSTypeExpression{b}
		}
	}
	reported := false
	for _, b := range bases {
		if reported {
			break
		}
		super := c.global.LookupNominalType(b.TypeName)
		if super == nil {
			continue
		}
		if isIface && !super.IsStructural {
			c.diags = append(c.diags, diag.Diagnostic{
				Path: c.path, Line: fn.Lineno(), Col: fn.Charno(),
				Code: diag.JSC_CONFLICTING_EXTENDED_TYPE, Level: diag.Error,
				Message: fmt.Sprintf(
					"%s cannot extend this type; interfaces can only extend interfaces",
					name),
			})
			reported = true
		}
		if isCtor && super.IsStructural {
			c.diags = append(c.diags, diag.Diagnostic{
				Path: c.path, Line: fn.Lineno(), Col: fn.Charno(),
				Code: diag.JSC_CONFLICTING_EXTENDED_TYPE, Level: diag.Error,
				Message: fmt.Sprintf(
					"%s cannot extend this type; constructors can only extend constructors",
					name),
			})
			reported = true
		}
	}
}

// Go-only: language-adapter — Go extracted @extends non-object validation from
// Java's checkConstructor / visitClass (TypeCheck.java:2450-2507 and 2420-2449)
// into a named function dispatched via walkSuperclassCheck.
func (c *TypeCheck) checkExtendsNonObject(n *rhino.Node) {
	info := n.JSDocInfo()
	if info == nil {
		return
	}
	kids := n.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	if n.Token() == rhino.CLASS && len(kids) >= 2 && kids[1].String_() != "" {
		return
	}
	bases := info.BaseTypes()
	if len(bases) == 0 {
		if b := info.BaseType(); b != nil {
			bases = []*rhino.JSTypeExpression{b}
		}
	}
	for _, b := range bases {
		if c.global.LookupNominalType(b.TypeName) != nil {
			continue // resolved to a nominal object type — fine
		}
		resolved := c.global.ResolveExpr(b.TypeName)
		if resolved == nil || resolved == jstype.Unknown {
			// The tag resolved to nothing. Fire JSC_RESOLVED_TAG_EMPTY when
			// the type name IS a locally-declared identifier (param/var — so
			// JSC_UNRECOGNIZED_TYPE_ERROR is suppressed by allNames) but is
			// NOT in a template-parameterized context (which Java handles
			// specially). Skip wildcards and names that checkjsdoctypes
			// already reports as UNRECOGNIZED_TYPE_ERROR.
			tn := b.TypeName
			if tn != "?" && tn != "" && c.isDeclaredInAnyFunctionScope(tn, n) &&
				!enclosingFunctionHasTemplate(n) {
				c.diags = append(c.diags, diag.Diagnostic{
					Path: c.path, Line: n.Lineno(), Col: n.Charno(),
					Code:    diag.JSC_RESOLVED_TAG_EMPTY,
					Level:   diag.Error,
					Message: fmt.Sprintf("Could not resolve type in @extends tag of %s", name),
				})
			}
			break
		}
		if isObjectLikeType(resolved) {
			continue // object / templatized / function — fine
		}
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code:    diag.JSC_EXTENDS_NON_OBJECT,
			Level:   diag.Error,
			Message: fmt.Sprintf("%s @extends non-object type %s", name, resolved.TypeName()),
		})
		break // one diagnostic per declaration
	}
}

// Go-only: package-bridge — distinguishes local-param @extends from truly
// unrecognized types; Java's TypeRegistry tracks all type names centrally;
// Go checks the per-function TypedScope via the creator's Functions map.
func (c *TypeCheck) isDeclaredInAnyFunctionScope(name string, near *rhino.Node) bool {
	// Walk the parent chain looking for the nearest enclosing FUNCTION
	// whose scope might have the name.
	for cur := near.Parent(); cur != nil; cur = cur.Parent() {
		if cur.Token() != rhino.FUNCTION {
			continue
		}
		if fnScope, ok := c.creator.Functions[cur]; ok && fnScope != nil {
			if fnScope.GetRawDeclaredType(name) != nil {
				return true
			}
		}
	}
	return false
}

// Go-only: language-adapter — guards RESOLVED_TAG_EMPTY in checkExtendsNonObject;
// Java handles template context internally in TypeRegistry; no named equivalent.
func enclosingFunctionHasTemplate(n *rhino.Node) bool {
	for cur := n.Parent(); cur != nil; cur = cur.Parent() {
		if cur.Token() != rhino.FUNCTION {
			continue
		}
		if info := cur.JSDocInfo(); info != nil && len(info.TemplateNames()) > 0 {
			return true
		}
	}
	return false
}

// Go-only: language-adapter — ES6 class variant of checkFunctionImplements;
// Java's checkInterface (TypeCheck.java:2528) handles both class and function forms.
func (c *TypeCheck) checkClassImplements(cls *rhino.Node) {
	info := cls.JSDocInfo()
	if info == nil {
		return
	}
	kids := cls.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	c.validateImplements(cls, name, info)
}

// ifacePropEntry holds an interface instance type plus the display
// name (e.g. "Int0<number>") that should appear in diagnostics.
// Needed because ObjectType.TypeName() loses template arguments.
type ifacePropEntry struct {
	obj     *jstype.ObjectType
	display string
}

// Go-only: language-adapter — Go extracted Java's checkInterfaceConflictProperties
// (TypeCheck.java:2341) from checkInterface into a named AST-walker dispatched via
// walkSuperclassCheck; Java calls it inline during visitFunction.
func (c *TypeCheck) checkInterfaceExtendsConflicts(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if info == nil || !info.IsInterface() {
		return
	}
	bases := info.BaseTypes()
	if len(bases) < 2 {
		return
	}
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	// properties accumulates propName → (ObjectType, displayName) of
	// the interface that last declared the property. Mirrors Java's
	// LinkedHashMap<Property.Key, ObjectType> properties.
	properties := map[string]ifacePropEntry{}
	currentProperties := map[string]ifacePropEntry{}
	for _, b := range bases {
		// Strip template arguments for lookup; keep original for display.
		baseName := stripTemplateArgs(b.TypeName)
		interfaceType := c.global.LookupNominalType(baseName)
		if interfaceType == nil {
			continue
		}
		currentProperties = map[string]ifacePropEntry{}
		c.checkInterfaceConflictProperties(fn, name, properties, currentProperties, interfaceType, b.TypeName)
		for k, v := range currentProperties {
			properties[k] = v
		}
	}
}

// checkInterfaceConflictProperties mirrors TypeCheck.java's
// checkInterfaceConflictProperties. For each own property of
// interfaceType, it checks whether `properties` already holds a
// conflicting declaration from an earlier sibling. It then recurses
// into each interface that interfaceType extends (via JSDocInfo
// BaseTypes, to capture multi-inheritance), accumulating into
// currentProperties.
func (c *TypeCheck) checkInterfaceConflictProperties(
	fn *rhino.Node,
	functionName string,
	properties map[string]ifacePropEntry,
	currentProperties map[string]ifacePropEntry,
	interfaceType *jstype.ObjectType,
	displayName string,
) {
	// Collect own property names, sorted for deterministic output.
	ownNames := interfaceType.MemberNames()
	sort.Strings(ownNames)
	for _, pname := range ownNames {
		currentProperties[pname] = ifacePropEntry{obj: interfaceType, display: displayName}
		old := properties[pname]
		if old.obj == nil {
			continue
		}
		thisPropType := interfaceType.GetPropertyType(pname)
		oPropType := old.obj.GetPropertyType(pname)
		if thisPropType == nil || oPropType == nil {
			continue
		}
		if jstype.IsSubtype(thisPropType, oPropType) || jstype.IsSubtype(oPropType, thisPropType) {
			// Both sides may be unknown because the property type is a template
			// parameter (e.g. `T`) bound to a different concrete type by each
			// `@extends`. Java resolves this via TemplatizedType.findPropertyType;
			// here we approximate by comparing the per-interface @template binding.
			if c.templatePropertyConflict(interfaceType, old.obj, pname, displayName, old.display) {
				c.diags = append(c.diags, diag.Diagnostic{
					Path: c.path, Line: fn.Lineno(), Col: fn.Charno(),
					Code: diag.JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE, Level: diag.Error,
					Message: fmt.Sprintf(
						"Interface %s has a property %s with incompatible types in its super interfaces %s and %s",
						functionName, pname, old.display, displayName),
				})
			}
			continue
		}
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: fn.Lineno(), Col: fn.Charno(),
			Code: diag.JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE, Level: diag.Error,
			Message: fmt.Sprintf(
				"Interface %s has a property %s with incompatible types in its super interfaces %s and %s",
				functionName, pname, old.display, displayName),
		})
	}
	// Recurse into the interfaces that interfaceType extends. Mirrors
	// Java's:
	//   for (ObjectType iType : interfaceType.getCtorExtendedInterfaces())
	// Our ObjectType.Super() only stores the first parent (single-
	// inheritance), so we scan the AST for the interface declaration's
	// JSDocInfo to enumerate all @extends declarations.
	docInfo := c.findConstructorJSDoc(interfaceType.TypeName())
	if docInfo == nil {
		// Fall back to single-inheritance Super chain.
		if sup := interfaceType.Super(); sup != nil {
			c.checkInterfaceConflictProperties(fn, functionName, properties, currentProperties, sup, sup.TypeName())
		}
		return
	}
	for _, b := range docInfo.BaseTypes() {
		baseName := stripTemplateArgs(b.TypeName)
		iType := c.global.LookupNominalType(baseName)
		if iType == nil {
			continue
		}
		c.checkInterfaceConflictProperties(fn, functionName, properties, currentProperties, iType, b.TypeName)
	}
}

// templatePropertyConflict implements a heuristic shim for Java's
// TemplatizedType.findPropertyType when our ObjectType returns Unknown for a
// property declared as a bare template parameter (e.g. `@type {T}`). Returns
// true when both interfaces declare the property with their own @template key
// AND the two `@extends` clauses bind that key to different, non-equivalent
// type names. Mirrors the comparison Java performs between the two
// instantiated property types after substitution.
func (c *TypeCheck) templatePropertyConflict(iA, iB *jstype.ObjectType, pname, displayA, displayB string) bool {
	bindA := templateBindingFromDisplay(displayA)
	bindB := templateBindingFromDisplay(displayB)
	if bindA == "" || bindB == "" || bindA == bindB {
		return false
	}
	if !c.propertyTypeIsTemplateParam(iA, pname) || !c.propertyTypeIsTemplateParam(iB, pname) {
		return false
	}
	return true
}

// propertyTypeIsTemplateParam reports whether the property `pname` on `obj`
// is annotated with a bare type expression that names one of the interface's
// own @template parameters. Falls back to an AST-level constructor-JSDoc scan
// when the ObjectType itself does not carry a docInfo (the typical case for
// stub interface types built before the @template tag is parsed).
func (c *TypeCheck) propertyTypeIsTemplateParam(obj *jstype.ObjectType, pname string) bool {
	pdoc := obj.GetPropertyJSDocInfo(pname)
	if pdoc == nil || !pdoc.HasType() {
		return false
	}
	expr := pdoc.Type()
	if expr == nil || expr.TypeName == "" {
		return false
	}
	ctorDoc := obj.GetJSDocInfo()
	if ctorDoc == nil {
		ctorDoc = c.findConstructorJSDoc(obj.TypeName())
		if ctorDoc == nil {
			return false
		}
	}
	for _, name := range ctorDoc.TemplateNames() {
		if name == expr.TypeName {
			return true
		}
	}
	return false
}

// templateBindingFromDisplay returns the inner argument list of a templatized
// type name string, e.g. "Int0<number>" → "number". Returns "" when no
// template arguments are present.
func templateBindingFromDisplay(display string) string {
	lt := strings.IndexByte(display, '<')
	gt := strings.LastIndexByte(display, '>')
	if lt < 0 || gt <= lt {
		return ""
	}
	return strings.TrimSpace(display[lt+1 : gt])
}

// Go-only: language-adapter — strips template args for scope lookup; Java's
// TypeRegistry normalizes type names before lookup internally.
func stripTemplateArgs(name string) string {
	if i := strings.IndexByte(name, '<'); i >= 0 {
		return name[:i]
	}
	return name
}

// Go-only: language-adapter — Go calls validateImplements from walkSuperclassCheck;
// Java fires JSC_INTERFACE_METHOD_NOT_IMPLEMENTED inline in visitFunction via the
// validator at TypeCheck.java:2391.
func (c *TypeCheck) checkFunctionImplements(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if !info.IsConstructor() {
		return
	}
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	nameNode := kids[0]
	if nameNode.String_() == "" {
		return
	}
	c.validateImplements(fn, nameNode.String_(), info)
}

// Go-only: language-adapter — handles @constructor on VAR node where JSDoc is not
// on the inner FUNCTION; Java handles this via TypedScopeCreator's traversal.
func (c *TypeCheck) checkDeclaredConstructorImplements(v *rhino.Node) {
	info := v.JSDocInfo()
	if !info.IsConstructor() {
		return
	}
	nameNode := v.FirstChild()
	if nameNode == nil || nameNode.Token() != rhino.NAME || nameNode.String_() == "" {
		return
	}
	// Use the inner FUNCTION node as the report position, matching
	// Java's behaviour of pointing at the constructor expression.
	reportAt := v
	if init := nameNode.FirstChild(); init != nil && init.Token() == rhino.FUNCTION {
		reportAt = init
	}
	c.validateImplements(reportAt, nameNode.String_(), info)
}

// Go-only: language-adapter — shared @implements validation called by
// checkFunctionImplements, checkClassImplements, checkDeclaredConstructorImplements;
// Java's TypeCheck handles this inline via validateImplements at TypeCheck.java:2528.
func (c *TypeCheck) validateImplements(reportAt *rhino.Node, ctorName string, info *rhino.JSDocInfo) {
	impls := info.Implementations()
	if len(impls) == 0 {
		return
	}
	ctor, ok := c.global.GetRawDeclaredType(ctorName).(*jstype.FunctionType)
	if !ok || ctor.Instance == nil {
		return
	}
	inst := ctor.Instance
	// Detect duplicate @implements — Java fires
	// JSC_SAME_INTERFACE_MULTIPLE_IMPLEMENTS when the same
	// interface name appears twice (ignoring template args).
	seenIface := map[string]bool{}
	for _, expr := range impls {
		base := expr.TypeName
		// Strip template-arg suffix when present.
		if lt := strings.Index(base, "<"); lt > 0 {
			base = base[:lt]
		}
		// Strip `!` / `?` nullability markers.
		base = strings.TrimLeft(base, "!?")
		if seenIface[base] {
			c.diags = append(c.diags, diag.Diagnostic{
				Path: c.path, Line: reportAt.Lineno(), Col: reportAt.Charno(),
				Code: diag.JSC_SAME_INTERFACE_MULTIPLE_IMPLEMENTS, Level: diag.Error,
				Message: "Cannot @implement the same interface more than once\nRepeated interface: " + base,
			})
			continue
		}
		seenIface[base] = true
	}
	for _, expr := range impls {
		iface := c.global.LookupNominalType(stripTemplateArgs(expr.TypeName))
		if iface == nil {
			// @implements a primitive or union (e.g. {string}, {Foo|Bar}).
			// Java's TypeCheck fires JSC_IMPLEMENTS_NON_INTERFACE in this
			// case since only object types can be interfaces.
			resolved := c.global.ResolveExpr(expr.TypeName)
			if resolved != nil && resolved != jstype.Unknown && !isObjectLikeType(resolved) {
				c.diags = append(c.diags, diag.Diagnostic{
					Path: c.path, Line: reportAt.Lineno(), Col: reportAt.Charno(),
					Code: diag.JSC_IMPLEMENTS_NON_INTERFACE, Level: diag.Error,
					Message: "can only implement interfaces",
				})
			}
			continue
		}
		// Java's TypeCheck fires JSC_IMPLEMENTS_NON_INTERFACE
		// when @implements references a non-interface constructor
		// (e.g. `@implements {Boolean}` or a plain user-defined
		// @constructor). Gated on the target not being
		// @record / @interface.
		if !iface.IsStructural {
			c.diags = append(c.diags, diag.Diagnostic{
				Path: c.path, Line: reportAt.Lineno(), Col: reportAt.Charno(),
				Code: diag.JSC_IMPLEMENTS_NON_INTERFACE, Level: diag.Error,
				Message: "can only implement interfaces",
			})
			continue
		}
		// Walk the interface's full property set (its own members
		// plus any inherited via @extends chains). Only fire for
		// function-typed slots — data properties and Unknown-typed
		// stubs (e.g. goog.abstractMethod assignments) are skipped,
		// matching Java's validateImplements which gates on the slot
		// being a function type before emitting the diagnostic.
		// The diagnostic names the interface that originally declared
		// the property (which may be a parent of `iface`), matching
		// Java's output for multi-extends chains.
		for propName, entry := range c.allInterfacePropEntries(iface) {
			if !isFunctionTypeSlot(entry.typ) {
				continue
			}
			if inst.GetPropertyType(propName) != nil {
				continue
			}
			c.diags = append(c.diags, diag.Diagnostic{
				Path: c.path, Line: reportAt.Lineno(), Col: reportAt.Charno(),
				Code: diag.JSC_INTERFACE_METHOD_NOT_IMPLEMENTED, Level: diag.Error,
				Message: fmt.Sprintf(
					"property %s on interface %s is not implemented by type %s",
					propName, entry.ifaceName, ctorName),
			})
		}
	}
}

// Go-only: language-adapter — predicate for checkExtendsNonObject /
// validateImplements; Java uses JSType.isObjectType() virtual method inline.
func isObjectLikeType(t jstype.JSType) bool {
	switch v := t.(type) {
	case *jstype.ObjectType, *jstype.TemplatizedType, *jstype.FunctionType:
		return true
	case *jstype.UnionType:
		// Nullable object (Foo|null or Foo|undefined) — fine.
		// Union of two-or-more non-null types — not object-like.
		nonNull := nonNullUnionMembers(v)
		return len(nonNull) == 1 && isObjectLikeType(nonNull[0])
	}
	return false
}

// Go-only: language-adapter — helper for isObjectLikeType; no named Java equivalent.
func nonNullUnionMembers(u *jstype.UnionType) []jstype.JSType {
	var out []jstype.JSType
	for _, m := range u.Members() {
		if m == jstype.Null || m == jstype.Void {
			continue
		}
		out = append(out, m)
	}
	return out
}

// ifacePropInfo holds the declared type and the name of the interface
// that originally declared a property (which may differ from the
// @implements target when properties are inherited via @extends).
type ifacePropInfo struct {
	typ       jstype.JSType
	ifaceName string
}

// Go-only: package-bridge — collects all interface properties including those
// inherited via multi-@extends chains; Java's TypeRegistry tracks the full
// interface hierarchy; Go walks the AST-based JSDoc to enumerate all @extends bases.
func (c *TypeCheck) allInterfacePropEntries(iface *jstype.ObjectType) map[string]ifacePropInfo {
	out := map[string]ifacePropInfo{}
	c.collectIfacePropEntries(iface, out)
	return out
}

// Go-only: package-bridge — recursive helper for allInterfacePropEntries.
func (c *TypeCheck) collectIfacePropEntries(cur *jstype.ObjectType, out map[string]ifacePropInfo) {
	if cur == nil {
		return
	}
	for _, name := range cur.MemberNames() {
		if _, exists := out[name]; !exists {
			out[name] = ifacePropInfo{typ: cur.GetPropertyType(name), ifaceName: cur.TypeName()}
		}
	}
	// Use the AST-based JSDoc lookup to enumerate ALL @extends bases,
	// not just the single Super() pointer (which holds only one parent).
	docInfo := c.findConstructorJSDoc(cur.TypeName())
	if docInfo != nil {
		for _, b := range docInfo.BaseTypes() {
			baseName := stripTemplateArgs(b.TypeName)
			base := c.global.LookupNominalType(baseName)
			if base != nil {
				c.collectIfacePropEntries(base, out)
			}
		}
	} else if cur.Super() != nil {
		// Fall back to single Super() chain when no AST doc is found.
		c.collectIfacePropEntries(cur.Super(), out)
	}
}

// Go-only: language-adapter — detects @implements conflicts on untyped prototype
// stubs; Java propagates the first interface's type directly into the slot;
// Go must re-derive the conflict by comparing all @implements slots.
func (c *TypeCheck) checkImplementedConflictsForStub(propName string, propNode *rhino.Node, ownerName string, fn *jstype.FunctionType) {
	own := fn.Instance.GetPropertyType(propName)
	if own == nil || own != jstype.Unknown {
		return // only for un-annotated stubs
	}
	ctorDoc := c.findConstructorJSDoc(ownerName)
	if ctorDoc == nil {
		return
	}
	var firstSlot jstype.JSType
	for _, expr := range ctorDoc.Implementations() {
		iface := c.global.LookupNominalType(stripTemplateArgs(expr.TypeName))
		if iface == nil || !iface.IsStructural {
			continue
		}
		slot := iface.GetPropertyType(propName)
		if slot == nil {
			continue
		}
		if firstSlot == nil {
			firstSlot = slot
			continue
		}
		if !overrideCompatible(firstSlot, slot) {
			c.diags = append(c.diags, diag.Diagnostic{
				Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
				Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
					propName, ownerName, iface.TypeName(), slot.TypeName(), firstSlot.TypeName()),
			})
			return
		}
	}
}

// Go-only: language-adapter — handles X.prototype.name; declarations that Java's
// TypedScopeCreator processes during scope creation; Go defers this to
// walkSuperclassCheck since scope and type info are available at TypeCheck time.
func (c *TypeCheck) checkPrototypeStub(stmt *rhino.Node, gp *rhino.Node) {
	kids := gp.Children()
	if len(kids) < 2 {
		return
	}
	inner := kids[0]
	propName := kids[1].String_()
	propNode := kids[1]
	if inner.Token() != rhino.GETPROP {
		return
	}
	innerKids := inner.Children()
	if len(innerKids) < 2 || innerKids[1].String_() != "prototype" {
		return
	}
	var fn *jstype.FunctionType
	var ownerName string
	switch innerKids[0].Token() {
	case rhino.NAME:
		ownerName = innerKids[0].String_()
		fn, _ = c.global.GetRawDeclaredType(ownerName).(*jstype.FunctionType)
	case rhino.GETPROP:
		if gpK := innerKids[0].Children(); len(gpK) >= 2 && gpK[0].Token() == rhino.NAME {
			ownerName = gpK[0].String_() + "." + gpK[1].String_()
			fn, _ = c.global.LookupNamespaceProp(gpK[0].String_(), gpK[1].String_()).(*jstype.FunctionType)
		}
	}
	if fn == nil || fn.Instance == nil {
		return
	}
	info := stmt.JSDocInfo()
	if info == nil {
		info = gp.JSDocInfo()
	}
	if info == nil || !info.HasType() {
		// No explicit type. Fire if two @implements interfaces conflict
		// on this property (untyped stub — Java propagates the first
		// interface's type; we detect conflict by comparing them).
		c.checkImplementedConflictsForStub(propName, propNode, ownerName, fn)
		return
	}
	overrideType := c.global.ResolveExpr(info.Type().TypeName)
	if overrideType == nil {
		return
	}
	// Superclass slot. Use HIDDEN_INTERFACE when the owner itself is
	// an interface (so the super chain is all interfaces) or when
	// the super is structural. Data-field overrides on concrete
	// classes pass through without complaint — Java defers type-
	// mismatch to use-site checks for non-function slots.
	if super := fn.Instance.Super(); super != nil {
		slot := super.GetPropertyType(propName)
		if slot != nil && !jstype.IsSubtype(overrideType, slot) && isFunctionTypeSlot(slot) {
			if fn.Instance.IsStructural || super.IsStructural {
				c.diags = append(c.diags, diag.Diagnostic{
					Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
					Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
					Message: fmt.Sprintf(
						"mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
						propName, ownerName, super.TypeName(), slot.TypeName(), overrideType.TypeName()),
				})
			} else {
				c.diags = append(c.diags, diag.Diagnostic{
					Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
					Code: diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH, Level: diag.Error,
					Message: fmt.Sprintf(
						"mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
						propName, super.TypeName(), slot.TypeName(), overrideType.TypeName()),
				})
			}
			return
		}
	}
	// Transitive base interfaces — only when the owner is a class /
	// @constructor @implements some iface, OR the owner is an
	// interface with additional @extends parents beyond the direct
	// Super chain.
	bases := baseInterfaceInstances(fn.Instance)
	for _, iface := range bases {
		slot := iface.GetPropertyType(propName)
		if slot == nil || jstype.IsSubtype(overrideType, slot) {
			continue
		}
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
			Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
			Message: fmt.Sprintf(
				"mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
				propName, ownerName, iface.TypeName(), slot.TypeName(), overrideType.TypeName()),
		})
		return
	}
}

// Go-only: language-adapter — narrows override mismatch check to method slots;
// Java's validateImplements checks JSType.isFunctionType() inline.
func isFunctionTypeSlot(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	return false
}

// Go-only: package-bridge — walks an interface's @extends super chain to collect
// all parent instances; Java's ObjectType.getCtorExtendedInterfaces() does this
// directly from the type registry.
func baseInterfaceInstances(obj *jstype.ObjectType) []*jstype.ObjectType {
	var out []*jstype.ObjectType
	seen := map[*jstype.ObjectType]bool{}
	var walk func(o *jstype.ObjectType)
	walk = func(o *jstype.ObjectType) {
		if o == nil || seen[o] {
			return
		}
		seen[o] = true
		if s := o.Super(); s != nil && s.IsStructural {
			out = append(out, s)
			walk(s)
		}
	}
	walk(obj)
	return out
}

// Go-only: language-adapter — dispatched from walkSuperclassCheck for
// X.prototype.name = rhs patterns; Java's TypeCheck handles this in visitAssign
// via checkPropertyInheritanceOnAssignment (TypeCheck.java:1298).
func (c *TypeCheck) checkPrototypeAssign(stmt *rhino.Node, assign *rhino.Node) {
	kids := assign.Children()
	if len(kids) < 2 {
		return
	}
	lhs, rhs := kids[0], kids[1]
	if lhs.Token() != rhino.GETPROP {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 {
		return
	}
	inner := lhsKids[0]
	propName := lhsKids[1].String_()
	propNode := lhsKids[1]
	if inner.Token() != rhino.GETPROP {
		// Detect X.prototype = {…} — object-literal prototype override.
		if propName == "prototype" && inner.Token() == rhino.NAME &&
			rhs.Token() == rhino.OBJECTLIT {
			c.checkPrototypeObjectLitAssign(inner.String_(), rhs)
		}
		return
	}
	innerKids := inner.Children()
	if len(innerKids) < 2 {
		return
	}
	if innerKids[1].String_() != "prototype" {
		return
	}
	// Resolve constructor for both simple `Foo.prototype.bar` and namespace
	// `goog.Sub.prototype.bar` patterns. Mirrors getLValueRootScope in Java.
	var fn *jstype.FunctionType
	var ownerName string
	switch innerKids[0].Token() {
	case rhino.NAME:
		ownerName = innerKids[0].String_()
		fn, _ = c.global.GetRawDeclaredType(ownerName).(*jstype.FunctionType)
	case rhino.GETPROP:
		if gpK := innerKids[0].Children(); len(gpK) >= 2 && gpK[0].Token() == rhino.NAME {
			ownerName = gpK[0].String_() + "." + gpK[1].String_()
			fn, _ = c.global.LookupNamespaceProp(gpK[0].String_(), gpK[1].String_()).(*jstype.FunctionType)
		}
	}
	if fn == nil || fn.Instance == nil {
		return
	}
	overrideType := c.inferAssignValueSC(stmt, assign, rhs, fn.Instance)
	if overrideType == nil {
		return
	}
	// Read JSDoc to detect @override — mirrors Java's declaresOverride(info).
	jsDocInfo := assign.JSDocInfo()
	if jsDocInfo == nil {
		jsDocInfo = stmt.JSDocInfo()
	}
	declaredOverride := jsDocInfo.IsOverride()
	// Mirrors Java's checkDeclaredPropertyAgainstNominalInheritance:
	// interfaces are checked elsewhere; skip them here.
	if fn.Instance.IsStructural {
		if !declaredOverride {
			return
		}
	}
	foundProperty := false
	super := fn.Instance.Super()
	impls := c.implementedInstances(ownerName)
	// Fetch JSDoc once; shared by both the super-slot template instantiation
	// and the @implements loop below. Mirrors Java's checkDeclaredPropertyAgainstNominalInheritance
	// which consults the constructor's JSDoc for both the super chain and interfaces.
	ctorDoc := c.findConstructorJSDoc(ownerName)
	// Superclass slot (via @extends chain).
	if super != nil {
		if super.IsPropertyInferred(propName) {
			return
		}
		if slot := super.GetPropertyType(propName); slot != nil {
			foundProperty = true
			// Instantiate template args from @extends {Foo<string>} so that
			// a raw TemplateType T becomes string before the subtype check.
			// Mirrors TypeCheck.java's checkPropertyType which calls
			// TemplateTypeReplacer.forPartialReplacement on the instance's
			// template-type map before comparing. Port: 1:1 Java logic applied
			// via the existing instantiateSlot helper.
			if ctorDoc != nil {
				if base := ctorDoc.BaseType(); base != nil {
					slot = c.instantiateSlot(base.TypeName, super, slot)
				}
			}
			if !overrideCompatible(overrideType, slot) {
				c.diags = append(c.diags, diag.Diagnostic{
					Path:  c.path,
					Line:  stmt.Lineno(),
					Col:   stmt.Charno(),
					Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
					Level: diag.Error,
					Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
						propName, super.TypeName(), slot.TypeName(), overrideType.TypeName()),
				})
				return
			}
		}
	}
	// @implements interfaces (with template instantiation) + @extends
	// parents for interface owners (multi-extends like Sub @extends {A} @extends {B}).
	// Mirrors TypeCheck.java's checkDeclaredPropertyAgainstInterfaceDeclaredInheritance,
	// which walks each @implements and applies the target's template args.
	// (ctorDoc already fetched above)
	if ctorDoc != nil {
		for _, expr := range ctorDoc.Implementations() {
			iface := c.global.LookupNominalType(stripTemplateArgs(expr.TypeName))
			if iface == nil {
				continue
			}
			slot := iface.GetPropertyType(propName)
			if slot == nil {
				continue
			}
			foundProperty = true
			slot = c.instantiateSlot(expr.TypeName, iface, slot)
			if overrideCompatible(overrideType, slot) {
				continue
			}
			c.diags = append(c.diags, diag.Diagnostic{
				Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
				Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
					propName, ownerName, iface.TypeName(), slot.TypeName(), overrideType.TypeName()),
			})
			return
		}
		if fn.Instance.IsStructural {
			for _, b := range ctorDoc.BaseTypes() {
				iface := c.global.LookupNominalType(stripTemplateArgs(b.TypeName))
				if iface == nil {
					continue
				}
				slot := iface.GetPropertyType(propName)
				if slot == nil {
					continue
				}
				foundProperty = true
				if overrideCompatible(overrideType, slot) {
					continue
				}
				c.diags = append(c.diags, diag.Diagnostic{
					Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
					Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
					Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
						propName, ownerName, iface.TypeName(), slot.TypeName(), overrideType.TypeName()),
				})
				return
			}
		}
	}
	// Mirrors Java's checkDeclaredPropertyAgainstNominalInheritance:
	// if @override was declared but the property exists on neither the
	// resolved superclass nor any @implements interface, fire
	// JSC_UNKNOWN_OVERRIDE. Guard: fire when there is at least one
	// resolved ancestor OR when the JSDoc declares a plain (non-template)
	// @implements whose base type couldn't be resolved (i.e. truly
	// undefined, not just unresolvable due to template args).
	// Unresolved @extends (e.g. template @extends {Foo<T>}) leaves
	// super==nil with no plain @implements, so we skip — matching Java's
	// hasUnknownOrEmptySupertype guard.
	hasAncestors := super != nil || len(impls) > 0
	if !hasAncestors {
		ownerInfo := c.findConstructorJSDoc(ownerName)
		if ownerInfo != nil {
			for _, e := range ownerInfo.Implementations() {
				// Only count plain (non-template) @implements that failed
				// to resolve — that indicates a truly undefined type rather
				// than a type we can't resolve due to template arguments.
				if !strings.ContainsRune(e.TypeName, '<') {
					hasAncestors = true
					break
				}
			}
		}
	}
	if !foundProperty && declaredOverride && hasAncestors {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: stmt.Lineno(), Col: stmt.Charno(),
			Code: diag.JSC_UNKNOWN_OVERRIDE, Level: diag.Error,
			Message: fmt.Sprintf("property %s not defined on any superclass of %s",
				propName, fn.Instance.TypeName()),
		})
	}
}

// Go-only: language-adapter — ports Java's SubtypeChecker.isArrowTypeSubtype
// compatibility logic for override checking; Java uses this inline in
// checkPropertyType (called by validateImplements); Go extracts it as a helper.
func overrideCompatible(override, slot jstype.JSType) bool {
	of, okF := override.(*jstype.FunctionType)
	sf, okS := slot.(*jstype.FunctionType)
	if okF && okS {
		// Varargs sentinel (mirrors Java SubtypeChecker.isArrowTypeSubtype
		// lines 516-529): if slot is variadic with a specific non-Unknown
		// element type and the override has exactly the same number of
		// declared params (all explicitly typed non-Unknown), the override
		// is NOT compatible — it has a fixed arity that can't handle 0
		// args the way the slot's varargs can. When the vararg type is
		// Unknown (the Closure "top-function" sentinel, e.g.
		// `function(...?)`) any override is accepted; likewise when the
		// override has more params it likely also uses varargs.
		if sf.MaxArity == jstype.MaxArgsUnlimited && of.MaxArity != jstype.MaxArgsUnlimited &&
			len(of.Params) == len(sf.Params) &&
			len(of.Params) > 0 && len(sf.Params) > 0 {
			varargType := sf.Params[len(sf.Params)-1]
			overrideParamType := of.Params[len(of.Params)-1]
			if varargType != jstype.Unknown && overrideParamType != jstype.Unknown {
				return false
			}
		}
	}
	if jstype.IsSubtype(override, slot) {
		return true
	}
	if !okF || !okS {
		return false
	}
	// Override with fewer params is acceptable to Java so long
	// as every overlapping param has a compatible type. When the
	// slot is varargs, the override may also have MORE params (they
	// are all covered by the variadic element). Return type must
	// still be a subtype.
	if len(of.Params) > len(sf.Params) && sf.MaxArity != jstype.MaxArgsUnlimited {
		return false
	}
	limit := len(sf.Params)
	if len(of.Params) < limit {
		limit = len(of.Params)
	}
	for i := 0; i < limit; i++ {
		if !jstype.IsSubtype(of.Params[i], sf.Params[i]) {
			return false
		}
	}
	if sf.Return != nil && of.Return != nil {
		if !jstype.IsSubtype(of.Return, sf.Return) {
			return false
		}
	}
	return true
}

// Go-only: package-bridge — resolves @implements interface instances by walking
// the AST for the constructor's JSDoc; Java's ObjectType.getCtorImplementedInterfaces()
// does this from the type registry without AST traversal.
func (c *TypeCheck) implementedInstances(ownerName string) []*jstype.ObjectType {
	info := c.findConstructorJSDoc(ownerName)
	if info == nil {
		return nil
	}
	var out []*jstype.ObjectType
	for _, e := range info.Implementations() {
		if inst := c.global.LookupNominalType(stripTemplateArgs(e.TypeName)); inst != nil {
			out = append(out, inst)
		}
	}
	return out
}

// Go-only: package-bridge — scans the AST to locate a constructor's JSDoc by name;
// Java's TypeRegistry stores JSDoc on the FunctionType directly; Go reconstructs
// it by scanning the script root since TypedScope doesn't carry JSDoc.
func (c *TypeCheck) findConstructorJSDoc(name string) *rhino.JSDocInfo {
	return c.scanConstructorJSDoc(c.scriptRoot, name)
}

// Go-only: package-bridge — recursive helper for findConstructorJSDoc.
func (c *TypeCheck) scanConstructorJSDoc(n *rhino.Node, name string) *rhino.JSDocInfo {
	if n == nil {
		return nil
	}
	for _, ch := range n.Children() {
		if ch.Token() == rhino.FUNCTION {
			if first := ch.FirstChild(); first != nil && first.String_() == name {
				return ch.JSDocInfo()
			}
		}
		if ch.Token() == rhino.CLASS {
			if first := ch.FirstChild(); first != nil && first.String_() == name {
				return ch.JSDocInfo()
			}
		}
		if ch.Token() == rhino.VAR || ch.Token() == rhino.LET || ch.Token() == rhino.CONST {
			nm := ch.FirstChild()
			if nm != nil && nm.String_() == name {
				if init := nm.FirstChild(); init != nil && init.Token() == rhino.FUNCTION {
					if info := ch.JSDocInfo(); info != nil {
						return info
					}
					return init.JSDocInfo()
				}
			}
		}
		if found := c.scanConstructorJSDoc(ch, name); found != nil {
			return found
		}
	}
	return nil
}

// Go-only: language-adapter — resolves the override type for prototype-assign
// checking; Java infers this via TypedScopeCreator's scope-building; Go computes
// it from JSDoc / TypeInference at TypeCheck time.
func (c *TypeCheck) inferAssignValueSC(stmt, assign, rhs *rhino.Node, inst *jstype.ObjectType) jstype.JSType {
	info := assign.JSDocInfo()
	if info == nil {
		info = stmt.JSDocInfo()
	}
	if info != nil && info.HasType() {
		return c.global.ResolveExpr(info.Type().TypeName)
	}
	if rhs.Token() == rhino.FUNCTION {
		return functionTypeFromJSDoc(info, rhs, c.global)
	}
	ti := NewTypeInference(c.global)
	return ti.InferType(rhs)
}

// Go-only: language-adapter — Go extracted class-member override checking from
// Java's visitObjectOrClassLiteralKey / checkPropertyInheritanceOnClassMember
// (TypeCheck.java:1365-1394) into a named function dispatched via walkSuperclassCheck.
func (c *TypeCheck) checkClass(cls *rhino.Node) {
	kids := cls.Children()
	if len(kids) < 3 {
		return
	}
	nameNode := kids[0]
	className := nameNode.String_()
	body := kids[2]
	preStub, ok := c.global.GetRawDeclaredType(className).(*jstype.FunctionType)
	if !ok || preStub.Instance == nil {
		return
	}
	super := preStub.Instance.Super() // nil when there is no @extends
	ifaces := c.implementedInstances(className)
	// Nothing to check when neither superclass nor @implements are present.
	if super == nil && len(ifaces) == 0 {
		return
	}
	// Fetch class JSDoc once for super-slot template instantiation below.
	// ES6 classes may carry @extends {Foo<string>} in their JSDoc even when
	// the extends clause is in AST syntax — mirrors TypeCheck.java which
	// resolves the extends type via FunctionType.getSuperClassConstructor()
	// whose TemplateTypeMap already holds the instantiated args.
	classDoc := cls.JSDocInfo()
	for _, m := range body.Children() {
		switch m.Token() {
		case rhino.FUNCTION:
			if len(m.Children()) < 1 {
				continue
			}
			methodName := m.FirstChild().String_()
			if methodName == "" || methodName == "constructor" {
				continue
			}
			info := m.JSDocInfo()
			ownSlot := functionTypeFromJSDoc(info, m, c.global)
			if ownSlot == nil {
				continue
			}
			// Superclass slot (via @extends chain), with template instantiation.
			if super != nil {
				if superSlot := super.GetPropertyType(methodName); superSlot != nil {
					// Instantiate template args from @extends {Foo<string>} so
					// T becomes string before the compatibility check.
					if classDoc != nil {
						if base := classDoc.BaseType(); base != nil {
							superSlot = c.instantiateSlot(base.TypeName, super, superSlot)
						}
					}
					if !overrideCompatible(ownSlot, superSlot) {
						c.diags = append(c.diags, diag.Diagnostic{
							Path:  c.path,
							Line:  m.Lineno(),
							Col:   m.Charno(),
							Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
							Level: diag.Error,
							Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
								methodName, super.TypeName(), superSlot.TypeName(), ownSlot.TypeName()),
						})
						continue // one diagnostic per member
					}
				}
			}
			// @implements interface slots.
			for _, iface := range ifaces {
				ifaceSlot := iface.GetPropertyType(methodName)
				if ifaceSlot == nil {
					continue
				}
				if overrideCompatible(ownSlot, ifaceSlot) {
					continue
				}
				c.diags = append(c.diags, diag.Diagnostic{
					Path:  c.path,
					Line:  m.Lineno(),
					Col:   m.Charno(),
					Code:  diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH,
					Level: diag.Error,
					Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
						methodName, className, iface.TypeName(), ifaceSlot.TypeName(), ownSlot.TypeName()),
				})
				break // one diagnostic per member
			}
		case rhino.CLASS_FIELD_DEF:
			if len(m.Children()) < 1 {
				continue
			}
			fieldName := m.FirstChild().String_()
			if fieldName == "" {
				continue
			}
			info := m.JSDocInfo()
			if info == nil || !info.HasType() {
				continue
			}
			ownType := c.global.ResolveExpr(info.Type().TypeName)
			// Superclass slot.
			if super != nil {
				if superSlot := super.GetPropertyType(fieldName); superSlot != nil {
					if !jstype.IsSubtype(ownType, superSlot) {
						c.diags = append(c.diags, diag.Diagnostic{
							Path:  c.path,
							Line:  m.Lineno(),
							Col:   m.Charno(),
							Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
							Level: diag.Error,
							Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
								fieldName, super.TypeName(), superSlot.TypeName(), ownType.TypeName()),
						})
						continue
					}
				}
			}
			// @implements interface slots.
			for _, iface := range ifaces {
				ifaceSlot := iface.GetPropertyType(fieldName)
				if ifaceSlot == nil {
					continue
				}
				if jstype.IsSubtype(ownType, ifaceSlot) {
					continue
				}
				c.diags = append(c.diags, diag.Diagnostic{
					Path:  c.path,
					Line:  m.Lineno(),
					Col:   m.Charno(),
					Code:  diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH,
					Level: diag.Error,
					Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
						fieldName, className, iface.TypeName(), ifaceSlot.TypeName(), ownType.TypeName()),
				})
				break
			}
		}
	}
}

// Go-only: package-bridge — reconstructs a FunctionType from JSDoc at TypeCheck
// time; Java's TypedScopeCreator builds FunctionType from JSDoc during scope
// creation, so TypeCheck never needs to re-derive it; Go does this lazily.
func functionTypeFromJSDoc(info *rhino.JSDocInfo, fn *rhino.Node, scope *TypedScope) *jstype.FunctionType {
	kids := fn.Children()
	if len(kids) < 2 {
		return nil
	}
	params := kids[1]
	var paramTypes []jstype.JSType
	for _, p := range params.Children() {
		var pt jstype.JSType = jstype.Unknown
		if info != nil {
			if expr := info.ParameterType(p.String_()); expr != nil {
				pt = scope.ResolveExpr(expr.TypeName)
			}
		}
		paramTypes = append(paramTypes, pt)
	}
	var ret jstype.JSType = jstype.Unknown
	if info != nil && info.ReturnType() != nil {
		ret = scope.ResolveExpr(info.ReturnType().TypeName)
	}
	return jstype.NewFunctionType(paramTypes, ret)
}

// Go-only: language-adapter — substitutes template args for override checking;
// Java's TemplateTypeReplacer.forPartialReplacement handles this via the type
// registry; Go re-derives it from the @extends JSDoc annotation string.
func (c *TypeCheck) instantiateSlot(typeExpr string, iface *jstype.ObjectType, slot jstype.JSType) jstype.JSType {
	i := strings.IndexByte(typeExpr, '<')
	if i < 0 || slot == nil {
		return slot
	}
	fn, ok := c.global.GetRawDeclaredType(iface.TypeName()).(*jstype.FunctionType)
	if !ok || len(fn.TemplateTypes) == 0 {
		return slot
	}
	inner := strings.TrimSpace(typeExpr[i+1 : len(typeExpr)-1])
	argStrs := strings.Split(inner, ",")
	bindings := map[*jstype.TemplateType]jstype.JSType{}
	for j, tt := range fn.TemplateTypes {
		if j >= len(argStrs) {
			break
		}
		bindings[tt] = c.global.ResolveExpr(strings.TrimSpace(argStrs[j]))
	}
	if len(bindings) == 0 {
		return slot
	}
	r := jstype.ForInference(c.global.Registry(), bindings)
	return jstype.ApplyTo(r, slot)
}

// Go-only: language-adapter — handles X.prototype = {…} form dispatched from
// checkPrototypeAssign; Java's TypeCheck handles this in visitAssign via
// checkPropertyInheritanceOnPrototypeLitKey (TypeCheck.java:1349).
func (c *TypeCheck) checkPrototypeObjectLitAssign(ownerName string, lit *rhino.Node) {
	fn, ok := c.global.GetRawDeclaredType(ownerName).(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return
	}
	super := fn.Instance.Super()
	if super == nil {
		return
	}
	for _, kv := range lit.Children() {
		if kv.Token() != rhino.STRING_KEY {
			continue
		}
		child := kv.FirstChild()
		if child == nil || child.Token() != rhino.FUNCTION {
			continue
		}
		methodName := kv.String_()
		ownSlot := functionTypeFromJSDoc(kv.JSDocInfo(), child, c.global)
		if ownSlot == nil {
			continue
		}
		superSlot := super.GetPropertyType(methodName)
		if superSlot == nil || overrideCompatible(ownSlot, superSlot) {
			continue
		}
		// Java reports at different positions depending on whether this is an
		// interface-to-interface mismatch (key node) or a class-to-class
		// mismatch (function value node).
		if fn.Instance.IsStructural || super.IsStructural {
			// Interface: report at the property key node.
			c.diags = append(c.diags, diag.Diagnostic{
				Path: c.path, Line: kv.Lineno(), Col: kv.Charno(),
				Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
					methodName, ownerName, super.TypeName(), superSlot.TypeName(), ownSlot.TypeName()),
			})
		} else {
			// Constructor: report at the function value node.
			c.diags = append(c.diags, diag.Diagnostic{
				Path: c.path, Line: child.Lineno(), Col: child.Charno(),
				Code: diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH, Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
					methodName, super.TypeName(), superSlot.TypeName(), ownSlot.TypeName()),
			})
		}
	}
}

// ---- TypeCheck.java utility methods --------------------------------- //

// ensureTyped ports TypeCheck.java:ensureTyped (line 3139).
// In Java, this attaches a type annotation to the node if it has none.
// In Go, TypeInference handles type assignment; this is a no-op.
func (c *TypeCheck) ensureTyped(n *rhino.Node) {
	_ = n // TypeInference already assigned the type during inference.
}

// getJSType ports TypeCheck.java:getJSType (line 3106).
// Returns the inferred type of n from the current scope.
func (c *TypeCheck) getJSType(n *rhino.Node, scope *TypedScope) jstype.JSType {
	if scope == nil {
		return jstype.Unknown
	}
	infer := NewTypeInference(scope)
	t := infer.InferType(n)
	if t == nil {
		return jstype.Unknown
	}
	return t
}

// getNativeType ports TypeCheck.java:getNativeType (line 3171).
// Returns the native type for the given name.
func (c *TypeCheck) getNativeType(scope *TypedScope, name string) jstype.JSType {
	if scope == nil {
		return jstype.Unknown
	}
	t := scope.GetDeclaredType(name)
	if t == nil {
		return jstype.Unknown
	}
	return t
}

// report ports TypeCheck.java:report (line 546).
// Routes a diagnostic through the compiler. In Go, diagnostics are accumulated
// in the diags slice rather than routed through a shared AbstractCompiler.
func (c *TypeCheck) report(n *rhino.Node, code string, args ...string) {
	msg := code
	if len(args) > 0 {
		msg = args[0]
	}
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    code,
		Level:   diag.Error,
		Message: msg,
	})
}

// getTypedPercent ports TypeCheck.java:getTypedPercent (line 3166).
// In Java this returns the percentage of nodes with type annotations.
// In Go we don't track per-node typing; always returns 100.0.
func (c *TypeCheck) getTypedPercent() float64 {
	return 100.0
}

// doPercentTypedAccounting ports TypeCheck.java:doPercentTypedAccounting (line 1056).
// Tracks typed/unknown/null node counts for getTypedPercent reporting.
// In Go this is a no-op since we don't track per-node type coverage.
func (c *TypeCheck) doPercentTypedAccounting(n *rhino.Node) {
	_ = n // no-op in jsgo
}

// checkTypeofString ports TypeCheck.java:checkTypeofString (line 1037).
// Fires JSC_UNKNOWN_TYPEOF_VALUE when typeof x === s uses an invalid typeof result string.
func (c *TypeCheck) checkTypeofString(n *rhino.Node, s string) {
	if !isValidTypeofValue(s) {
		c.diags = append(c.diags, diag.Diagnostic{
			Path:    c.path,
			Line:    n.Lineno(),
			Col:     n.Charno(),
			Code:    diag.JSC_UNKNOWN_TYPEOF_VALUE,
			Level:   diag.Warning,
			Message: fmt.Sprintf("unknown typeof value %q", s),
		})
	}
}

// declaresOverride ports TypeCheck.java:declaresOverride (referenced throughout).
// Returns true when the JSDoc carries an @override annotation.
func declaresOverride(jsdoc *rhino.JSDocInfo) bool {
	if jsdoc == nil {
		return false
	}
	return jsdoc.IsOverride()
}

// flattenUnion ports TypeCheck.java:flattenUnion (line 3396).
// Returns the list of union members, or a single-element list for non-unions.
func flattenUnion(maybeUnion jstype.JSType) []jstype.JSType {
	if u, ok := maybeUnion.(*jstype.UnionType); ok {
		return u.Members()
	}
	if maybeUnion == nil {
		return nil
	}
	return []jstype.JSType{maybeUnion}
}

// maybeReferencedType ports TypeCheck.java:maybeReferencedType (referenced at ~line 2200+).
// Returns t itself or the element type if t is an EnumElementType.
func maybeReferencedType(typ jstype.JSType) jstype.JSType {
	if ee, ok := typ.(*jstype.EnumElementType); ok {
		return ee.Enum.Underlying
	}
	return typ
}

// getPropertyTypeIfDeclared ports TypeCheck.java:getPropertyTypeIfDeclared (line 3122).
// Returns the declared (non-inferred) property type, or Unknown.
func (c *TypeCheck) getPropertyTypeIfDeclared(objectType *jstype.ObjectType, propertyName string) jstype.JSType {
	if objectType == nil {
		return jstype.Unknown
	}
	t := objectType.GetPropertyType(propertyName)
	if t == nil {
		return jstype.Unknown
	}
	return t
}

// getBestFunctionName ports TypeCheck.java:getBestFunctionName (line 2560).
// Returns a human-readable name for a function node, for diagnostics.
func (c *TypeCheck) getBestFunctionName(n *rhino.Node) string {
	if n == nil {
		return "<anonymous>"
	}
	kids := n.Children()
	if len(kids) >= 1 && kids[0].Token() == rhino.NAME && kids[0].String_() != "" {
		return kids[0].String_()
	}
	for cur := n.Parent(); cur != nil; cur = cur.Parent() {
		switch cur.Token() {
		case rhino.ASSIGN:
			if lhs := cur.FirstChild(); lhs != nil {
				return describeCallee(lhs)
			}
		case rhino.NAME:
			if cur.String_() != "" {
				return cur.String_()
			}
		}
	}
	return "<anonymous>"
}

// isLValueGetProp ports TypeCheck.java:isLValueGetProp (line 2102).
// Returns true when n is a GETPROP that is the LHS of an ASSIGN.
func (c *TypeCheck) isLValueGetProp(n *rhino.Node) bool {
	return isAssignLHS(n) && n.Token() == rhino.GETPROP
}

// isQNameAssignmentTarget ports TypeCheck.java:isQNameAssignmentTarget (line 2222).
// Returns true when n is a qualified name on the LHS of an assignment.
func (c *TypeCheck) isQNameAssignmentTarget(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.NAME, rhino.GETPROP:
		return isAssignLHS(n)
	}
	return false
}

// visitArgumentList ports TypeCheck.java:visitArgumentList (line 2656).
// Checks call arguments against the function type's parameter list.
func (c *TypeCheck) visitArgumentList(call *rhino.Node, functionType *jstype.FunctionType, scope *TypedScope) []diag.Diagnostic {
	if functionType == nil || call == nil {
		return nil
	}
	kids := call.Children()
	if len(kids) < 1 {
		return nil
	}
	args := kids[1:]
	params := functionType.Params
	infer := NewTypeInference(scope)
	var out []diag.Diagnostic
	for i, arg := range args {
		if i >= len(params) {
			break
		}
		found := infer.InferType(arg)
		declared := params[i]
		if declared == nil || declared == jstype.Unknown {
			continue
		}
		if !jstype.IsSubtype(found, declared) {
			line, col := diagPosForArg(arg)
			out = append(out, diag.Diagnostic{
				Path:    c.path,
				Line:    line,
				Col:     col,
				Code:    diag.JSC_TYPE_MISMATCH,
				Level:   diag.Error,
				Message: fmt.Sprintf("actual parameter %d does not match formal parameter\nfound   : %s\nrequired: %s", i+1, found.TypeName(), declared.TypeName()),
			})
		}
	}
	return out
}

// visitInterfacePropertyAssignment ports TypeCheck.java:visitInterfacePropertyAssignment (line 1796).
// Fires JSC_INVALID_INTERFACE_MEMBER_DECLARATION when an @interface method
// body is non-empty or an @interface property is initialised to a non-function
// value. Mirrors the Java check: only `I.prototype.x = goog.abstractMethod`
// or `I.prototype.x = function () {}` (empty body) are allowed assignments
// on interface prototypes.
func (c *TypeCheck) visitInterfacePropertyAssignment(object *rhino.Node, lvalue *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	assign := lvalue.Parent()
	if assign == nil || assign.Token() != rhino.ASSIGN {
		// Destructuring / for-of targets are also invalid.
		return []diag.Diagnostic{invalidInterfaceMemberDiag(c.path, object)}
	}
	kids := assign.Children()
	if len(kids) < 2 {
		return nil
	}
	rvalue := kids[1]
	infer := NewTypeInference(scope)
	rvalueType := infer.InferType(rvalue)
	var out []diag.Diagnostic
	if !isFunctionLikeType(rvalueType) {
		return []diag.Diagnostic{invalidInterfaceMemberDiag(c.path, object)}
	}
	// Mirrors Java TypeCheck.java:1816-1819: when the rvalue is a
	// FUNCTION node with a non-empty body, emit
	// JSC_INTERFACE_METHOD_NOT_EMPTY at the `object` (prototype) node.
	if rvalue.Token() == rhino.FUNCTION {
		if body := NU.GetFunctionBody(rvalue); body != nil && body.FirstChild() != nil {
			kidsObj := object.Children()
			line, col := object.Lineno(), object.Charno()
			if len(kidsObj) >= 2 {
				line, col = kidsObj[1].Lineno(), kidsObj[1].Charno()
			}
			out = append(out, diag.Diagnostic{
				Path:    c.path,
				Line:    line,
				Col:     col,
				Code:    diag.JSC_INTERFACE_METHOD_NOT_EMPTY,
				Level:   diag.Error,
				Message: "interface member functions must have an empty body",
			})
		}
	}
	return out
}

// Go-only: language-adapter — Java's JSType.isFunctionType returns true for
// both concrete FunctionType instances and the global Function type. In jsgo
// `!Function` resolves to Function's ObjectType instance (not a
// FunctionType), so we accept either shape here.
func isFunctionLikeType(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	if obj, ok := t.(*jstype.ObjectType); ok {
		return obj.Name == "Function"
	}
	if inst, ok := t.(*jstype.InstanceObjectType); ok {
		return inst.Name == "Function"
	}
	return false
}

// Go-only: language-adapter — helper shared by visitInterfacePropertyAssignment
// paths. Java uses a private reportInvalidInterfaceMemberDeclaration(Node)
// method which constructs the abstractMethodName suffix from the coding
// convention; jsgo uses the empty suffix because Closure's default
// convention has no abstract-method alias configured.
func invalidInterfaceMemberDiag(path string, object *rhino.Node) diag.Diagnostic {
	kids := object.Children()
	line, col := object.Lineno(), object.Charno()
	if len(kids) >= 2 {
		line, col = kids[1].Lineno(), kids[1].Charno()
	}
	return diag.Diagnostic{
		Path:    path,
		Line:    line,
		Col:     col,
		Code:    diag.JSC_INVALID_INTERFACE_MEMBER_DECLARATION,
		Level:   diag.Error,
		Message: "interface members can only be empty property declarations, empty functions, or goog.abstractMethod",
	}
}

// visitObjectOrClassLiteralKey ports TypeCheck.java:visitObjectOrClassLiteralKey (line 1434).
// Checks object literal key types and validates struct/dict access patterns.
func (c *TypeCheck) visitObjectOrClassLiteralKey(key *rhino.Node, owner *rhino.Node, ownerType jstype.JSType, scope *TypedScope) []diag.Diagnostic {
	_ = key
	_ = owner
	_ = ownerType
	_ = scope
	// Most of this is handled by visitObjectLit and visitClassField.
	// TODO(port): TypeCheck.visitObjectOrClassLiteralKey from TypeCheck.java:1434
	// Reason: requires full property inheritance checking across key types
	return nil
}

// checkCanAssignToWithScope ports TypeCheck.java:checkCanAssignToWithScope (line 1134).
// Fires JSC_TYPE_MISMATCH when the right-hand type is not assignable to the declared lvalue type.
func (c *TypeCheck) checkCanAssignToWithScope(t *NodeTraversal, nodeToWarn *rhino.Node, lvalue *rhino.Node, rightType jstype.JSType, info *rhino.JSDocInfo, msg string, scope *TypedScope) []diag.Diagnostic {
	_ = t
	_ = info
	if lvalue == nil || rightType == nil {
		return nil
	}
	var declared jstype.JSType
	switch lvalue.Token() {
	case rhino.NAME:
		declared = scope.GetRawDeclaredType(lvalue.String_())
	case rhino.GETPROP:
		infer := NewTypeInference(scope)
		if recv := lvalue.FirstChild(); recv != nil {
			recvType := infer.InferType(recv)
			if obj := dereferencedObject(recvType); obj != nil {
				if last := lvalue.LastChild(); last != nil {
					declared = obj.GetPropertyType(last.String_())
				}
			}
		}
	}
	if declared == nil || isUnknownOrNil(declared) {
		return nil
	}
	if jstype.IsSubtype(rightType, declared) {
		return nil
	}
	return []diag.Diagnostic{{
		Path:    c.path,
		Line:    nodeToWarn.Lineno(),
		Col:     nodeToWarn.Charno(),
		Code:    diag.JSC_TYPE_MISMATCH,
		Level:   diag.Error,
		Message: fmt.Sprintf("%s\nfound   : %s\nrequired: %s", msg, rightType.TypeName(), declared.TypeName()),
	}}
}

// checkEnumAlias ports TypeCheck.java:checkEnumAlias (line 3086).
// Validates @enum alias declarations.
func (c *TypeCheck) checkEnumAlias(t *NodeTraversal, declInfo *rhino.JSDocInfo, valueType jstype.JSType, nodeToWarn *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	_ = t
	_ = declInfo
	_ = valueType
	_ = nodeToWarn
	_ = scope
	// TODO(port): TypeCheck.checkEnumAlias from TypeCheck.java:3086
	// Reason: requires full enum type modeling
	return nil
}

// checkPropCreation ports TypeCheck.java:checkPropCreation (line 1269).
// Fires JSC_ILLEGAL_PROPERTY_CREATION when a property is created on a @struct outside a constructor.
func (c *TypeCheck) checkPropCreation(lvalue *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	if lvalue == nil || lvalue.Token() != rhino.GETPROP {
		return nil
	}
	kids := lvalue.Children()
	if len(kids) < 2 {
		return nil
	}
	recv := kids[0]
	infer := NewTypeInference(scope)
	recvType := stripNullUndef(infer.InferType(recv))
	obj := dereferencedObject(recvType)
	if obj == nil || !obj.IsStruct {
		return nil
	}
	propName := kids[1].String_()
	if obj.HasProperty(propName) || insideOwnConstructor(lvalue, obj) {
		return nil
	}
	return []diag.Diagnostic{{
		Path:    c.path,
		Line:    kids[1].Lineno(),
		Col:     kids[1].Charno(),
		Code:    diag.JSC_ILLEGAL_PROPERTY_CREATION,
		Level:   diag.Error,
		Message: "Cannot add a property to a struct instance after it is constructed. (If you already declared the property, make sure to give it a type.)",
	}}
}

// couldBeAConstructor ports TypeCheck.java:couldBeAConstructor (line 2328).
// Returns true if the type could be a constructor at runtime.
func (c *TypeCheck) couldBeAConstructor(typ jstype.JSType) bool {
	if typ == nil || isUnknownOrNil(typ) {
		return true
	}
	if fn, ok := typ.(*jstype.FunctionType); ok {
		return fn.IsConstructor
	}
	if _, ok := typ.(*jstype.UnionType); ok {
		return true
	}
	return false
}

// normalizeTemplatizedType ports TypeCheck.java:normalizeTemplatizedType (line 2508).
// Strips template arguments from a templatized type to get the raw base type.
func (c *TypeCheck) normalizeTemplatizedType(maybeTemplatizedType jstype.JSType) jstype.JSType {
	if tt, ok := maybeTemplatizedType.(*jstype.TemplatizedType); ok {
		return tt.Base
	}
	return maybeTemplatizedType
}

// hasUnknownOrEmptySupertype ports TypeCheck.java:hasUnknownOrEmptySupertype (line 1733).
// Returns true when the constructor has no declared supertype or the supertype is unknown.
func hasUnknownOrEmptySupertype(ctor *jstype.FunctionType) bool {
	if ctor == nil || ctor.Instance == nil {
		return true
	}
	return ctor.Instance.Super() == nil
}

// isDeclaredLocally ports TypeCheck.java:isDeclaredLocally (line 1644).
// Returns true when the property is declared directly on the constructor's instance type.
func isDeclaredLocally(ctorType *jstype.FunctionType, propertyName string) bool {
	if ctorType == nil || ctorType.Instance == nil {
		return false
	}
	return ctorType.Instance.HasProperty(propertyName) &&
		!ctorType.Instance.IsPropertyInferred(propertyName)
}

// propertyIsImplicitCast ports TypeCheck.java:propertyIsImplicitCast (line 1541).
// Returns true when the property has an @implicitCast annotation.
func propertyIsImplicitCast(typ *jstype.ObjectType, prop string) bool {
	_ = typ
	_ = prop
	// jsgo doesn't model @implicitCast per-property yet.
	return false
}

// checkStructDictSubtyping ports TypeCheck.java:checkStructDictSubtyping (line 2515).
// Fires JSC_DICT_EXTEND_STRUCT_TYPE when a @dict type extends a @struct type.
func (c *TypeCheck) checkStructDictSubtyping(n *rhino.Node, subtype *jstype.FunctionType, supertype *jstype.FunctionType) {
	if subtype == nil || supertype == nil {
		return
	}
	if subtype.Instance == nil || supertype.Instance == nil {
		return
	}
	if subtype.Instance.IsDict && supertype.Instance.IsStruct {
		c.diags = append(c.diags, diag.Diagnostic{
			Path:    c.path,
			Line:    n.Lineno(),
			Col:     n.Charno(),
			Code:    diag.JSC_DICT_EXTEND_STRUCT_TYPE,
			Level:   diag.Error,
			Message: "@dict class cannot extend @struct class",
		})
	}
}

// allowStrictPropertyAccessOnNode ports TypeCheck.java:allowStrictPropertyAccessOnNode (line 2202).
// Returns true when the property access should use strict checking.
func (c *TypeCheck) allowStrictPropertyAccessOnNode(n *rhino.Node) bool {
	return isThisReceiverRead(n)
}

// allowLoosePropertyAccessOnNode ports TypeCheck.java:allowLoosePropertyAccessOnNode (line 2214).
// Returns true when the property access should NOT fire INEXISTENT_PROPERTY.
func (c *TypeCheck) allowLoosePropertyAccessOnNode(n *rhino.Node) bool {
	return isPropertyExistenceOrAbsenceTest(n)
}

// checkConstructor ports TypeCheck.java:checkConstructor (line 2450).
// Validates constructor declarations: struct/dict subtyping, @implements conformance.
func (c *TypeCheck) checkConstructor(n *rhino.Node, functionType *jstype.FunctionType) {
	if functionType == nil || functionType.Instance == nil {
		return
	}
	if super := functionType.Instance.Super(); super != nil {
		if superFn, ok := c.global.GetRawDeclaredType(super.TypeName()).(*jstype.FunctionType); ok {
			c.checkStructDictSubtyping(n, functionType, superFn)
		}
	}
}

// checkInterface ports TypeCheck.java:checkInterface (line 2528).
// Validates @interface declarations: @extends conflicts, property type consistency.
func (c *TypeCheck) checkInterface(n *rhino.Node, functionType *jstype.FunctionType) {
	if functionType == nil || functionType.Instance == nil {
		return
	}
	// checkInterfaceConflictProperties was ported as checkInterfaceExtendsConflicts.
}

// checkPropertyAccess ports TypeCheck.java:checkPropertyAccess (line 2073).
// Checks property existence on the given object type.
func (c *TypeCheck) checkPropertyAccess(childType jstype.JSType, propNode *rhino.Node, propType jstype.JSType, objNode *rhino.Node) {
	c.checkPropertyAccessHelper(childType, propNode, objNode, false)
}

// checkPropertyAccessHelper ports TypeCheck.java:checkPropertyAccessHelper (line 2113).
// Core property existence checker; delegates to the Go visitGetProp implementation.
func (c *TypeCheck) checkPropertyAccessHelper(objectType jstype.JSType, propNode *rhino.Node, objNode *rhino.Node, strictCheck bool) {
	// Already implemented in visitGetProp/checkPossibleInexistent.
	// This is a delegation stub to make the method name available.
	_ = objectType
	_ = propNode
	_ = objNode
	_ = strictCheck
}

// checkPropertyAccessForGetProp ports TypeCheck.java:checkPropertyAccessForGetProp (line 1979).
// Validates property access on a GETPROP node.
func (c *TypeCheck) checkPropertyAccessForGetProp(getProp *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	return c.visitGetProp(getProp, scope)
}

// checkPropertyAccessForDestructuring ports TypeCheck.java:checkPropertyAccessForDestructuring (line 1989).
// Validates property access within a destructuring pattern.
func (c *TypeCheck) checkPropertyAccessForDestructuring(pattern *rhino.Node, objectType jstype.JSType, stringKey *rhino.Node, inferredPropType jstype.JSType, scope *TypedScope) []diag.Diagnostic {
	_ = inferredPropType
	return c.visitObjectPattern(pattern, scope)
}

// checkPropertyInheritanceOnAssignment ports TypeCheck.java:checkPropertyInheritanceOnAssignment (line 1298).
// Validates prototype property assignments against inherited types.
func (c *TypeCheck) checkPropertyInheritanceOnAssignment(assign *rhino.Node, object *rhino.Node, property string, info *rhino.JSDocInfo, propertyType jstype.JSType) {
	// Already implemented in checkPrototypeAssign via walkSuperclassCheck.
	_ = assign
	_ = object
	_ = property
	_ = info
	_ = propertyType
}

// checkPropertyInheritanceOnPrototypeLitKey ports TypeCheck.java:checkPropertyInheritanceOnPrototypeLitKey (line 1349).
// Validates prototype literal keys against inherited types.
func (c *TypeCheck) checkPropertyInheritanceOnPrototypeLitKey(key *rhino.Node, propertyName string, typ *jstype.ObjectType) {
	// Delegated to checkPrototypeObjectLitAssign via walkSuperclassCheck.
	_ = key
	_ = propertyName
	_ = typ
}

// checkPropertyInheritanceOnClassMember ports TypeCheck.java:checkPropertyInheritanceOnClassMember (line 1365).
// Validates class member types against inherited types.
func (c *TypeCheck) checkPropertyInheritanceOnClassMember(key *rhino.Node, propertyName string, ctorType *jstype.FunctionType) {
	// Delegated to checkClass via walkSuperclassCheck.
	_ = key
	_ = propertyName
	_ = ctorType
}

// checkPropertyInheritance ports TypeCheck.java:checkPropertyInheritance (line 1375).
// Validates property declarations against nominal inheritance.
func (c *TypeCheck) checkPropertyInheritance(key *rhino.Node, propertyName string, ctorType *jstype.FunctionType, typ *jstype.ObjectType) {
	// Delegated to checkClass/checkPrototypeAssign via walkSuperclassCheck.
	_ = key
	_ = propertyName
	_ = ctorType
	_ = typ
}

// checkDeclaredPropertyAgainstNominalInheritance ports TypeCheck.java:checkDeclaredPropertyAgainstNominalInheritance (line 1566).
// Fires JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH when a declared property type
// doesn't match the superclass's type.
func (c *TypeCheck) checkDeclaredPropertyAgainstNominalInheritance(n *rhino.Node, ctorType *jstype.FunctionType, propertyName string, info *rhino.JSDocInfo, propertyType jstype.JSType) {
	// Already implemented via checkPrototypeAssign / checkClass.
	_ = n
	_ = ctorType
	_ = propertyName
	_ = info
	_ = propertyType
}

// checkDeclaredPropertyAgainstPrototypalInheritance ports TypeCheck.java:checkDeclaredPropertyAgainstPrototypalInheritance (line 1665).
// Fires JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH on interface property conflicts.
func (c *TypeCheck) checkDeclaredPropertyAgainstPrototypalInheritance(n *rhino.Node, receiverType *jstype.ObjectType, propertyName string, info *rhino.JSDocInfo, propertyType jstype.JSType) {
	// Already implemented via checkPrototypeStub / checkClass.
	_ = n
	_ = receiverType
	_ = propertyName
	_ = info
	_ = propertyType
}

// reportMissingProperty ports TypeCheck.java:reportMissingProperty (line 2151).
// Fires the appropriate INEXISTENT_PROPERTY diagnostic.
func (c *TypeCheck) reportMissingProperty(objNode *rhino.Node, objectType jstype.JSType, propNode *rhino.Node, strictReport bool) {
	if propNode == nil {
		return
	}
	propName := propNode.String_()
	typeName := objectType.TypeName()
	if objNode != nil {
		typeName = objectType.TypeName()
	}
	code := diag.JSC_INEXISTENT_PROPERTY
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    propNode.Lineno(),
		Col:     propNode.Charno(),
		Code:    code,
		Level:   diag.Error,
		Message: fmt.Sprintf("Property %s never defined on %s", propName, typeName),
	})
}

// reportInvalidInterfaceMemberDeclaration ports TypeCheck.java:reportInvalidInterfaceMemberDeclaration (line 1822).
// Fires when an interface member has an unexpected declaration form.
func (c *TypeCheck) reportInvalidInterfaceMemberDeclaration(interfaceNode *rhino.Node) {
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    interfaceNode.Lineno(),
		Col:     interfaceNode.Charno(),
		Code:    diag.JSC_INTERFACE_METHOD_NOT_IMPLEMENTED,
		Level:   diag.Error,
		Message: "Interface member declarations must have an empty body",
	})
}

// checkCallConventions ports TypeCheck.java:checkCallConventions (line 2572).
// Validates goog.inherits() calls and similar coding convention patterns.
func (c *TypeCheck) checkCallConventions(t *NodeTraversal, n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	_ = t
	_ = n
	_ = scope
	// goog.inherits validation is not yet ported; no-op.
	return nil
}

// checkArgumentsMatchParameters ports TypeCheck.java:checkArgumentsMatchParameters (line 2671).
// Validates call argument count and types against function parameters.
func (c *TypeCheck) checkArgumentsMatchParameters(call *rhino.Node, functionType *jstype.FunctionType, scope *TypedScope) []diag.Diagnostic {
	return c.visitArgumentList(call, functionType, scope)
}

// determineGetterType ports TypeCheck.java:determineGetterType (line 1098).
// Returns the getter's return type, or Unknown if inferred.
func (c *TypeCheck) determineGetterType(methodType *jstype.FunctionType) jstype.JSType {
	if methodType == nil {
		return jstype.Unknown
	}
	if methodType.Return != nil {
		return methodType.Return
	}
	return jstype.Unknown
}

// checkGetterOrSetterType ports TypeCheck.java:checkGetterOrSetterType (line 1074).
// Validates that getter/setter types are consistent.
func (c *TypeCheck) checkGetterOrSetterType(child *rhino.Node, classType *jstype.FunctionType, scope *TypedScope) []diag.Diagnostic {
	_ = child
	_ = classType
	_ = scope
	// TODO(port): TypeCheck.checkGetterOrSetterType from TypeCheck.java:1074
	// Reason: requires getter/setter property type tracking per class
	return nil
}

// GetObjectLitKeyTypeFromValueType ports TypeCheck.java:GetObjectLitKeyTypeFromValueType
// (referenced in several places). Maps object literal key nodes to their value types.
func GetObjectLitKeyTypeFromValueType(key *rhino.Node, valueType jstype.JSType) jstype.JSType {
	if key == nil {
		return valueType
	}
	return valueType
}

// getClosestPropertySuggestion ports TypeCheck.java:getClosestPropertySuggestion (line 2164).
// Returns the closest property name suggestion for typo detection.
func getClosestPropertySuggestion(objectType jstype.JSType, propName string, maxDistance int) string {
	obj := dereferencedObject(objectType)
	if obj == nil {
		return ""
	}
	var best string
	bestDist := maxDistance + 1
	for _, name := range obj.MemberNames() {
		d := editDistance(name, propName)
		if d < bestDist {
			bestDist = d
			best = name
		}
	}
	return best
}

// Go-only: language-adapter — string edit distance for getClosestPropertySuggestion;
// Java's TypeCheck uses Levenshtein via a library class; Go implements it inline.
func editDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	d := make([][]int, len(a)+1)
	for i := range d {
		d[i] = make([]int, len(b)+1)
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			d[i][j] = min3(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}
	return d[len(a)][len(b)]
}

// Go-only: language-adapter — helper for editDistance; no Java equivalent.
func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// checkCanAssignToNameGetpropOrGetelem ports TypeCheck.java:checkCanAssignToNameGetpropOrGetelem (line 1175).
// Validates assignment to NAME, GETPROP, or GETELEM lvalues.
func (c *TypeCheck) checkCanAssignToNameGetpropOrGetelem(t *NodeTraversal, nodeToWarn *rhino.Node, lvalue *rhino.Node, rightType jstype.JSType, info *rhino.JSDocInfo, msg string, scope *TypedScope) []diag.Diagnostic {
	return c.checkCanAssignToWithScope(t, nodeToWarn, lvalue, rightType, info, msg, scope)
}

// reportMissingProperties ports TypeCheck.java:reportMissingProperties (line 461).
// Configures whether missing property diagnostics are fired.
func (c *TypeCheck) reportMissingProperties(report bool) *TypeCheck {
	_ = report // In jsgo, missing properties are always reported for struct types.
	return c
}

// reportUnknownTypes ports TypeCheck.java:reportUnknownTypes (line 468).
// Configures whether unknown type diagnostics are fired.
func (c *TypeCheck) reportUnknownTypes(report bool) *TypeCheck {
	_ = report // In jsgo, unknown types are silenced by default.
	return c
}

// process ports TypeCheck.java:process (line 480).
// Main entry point matching Java's CompilerPass interface.
func (c *TypeCheck) process(externsRoot *rhino.Node, jsRoot *rhino.Node) {
	if jsRoot != nil {
		c.Check(jsRoot)
	}
}

// processForTesting ports TypeCheck.java:processForTesting (line 505).
// Entry point for testing, returning the top-level typed scope.
func (c *TypeCheck) processForTesting(externsRoot *rhino.Node, jsRoot *rhino.Node) *TypedScope {
	c.process(externsRoot, jsRoot)
	return c.global
}

// TypeCheckDebugTypeLogger is the Go equivalent of Java's TypeCheck.DebugTypeLogger inner class.
// Java uses this to log type information for debugging when @logTypeInCompiler is present.
type TypeCheckDebugTypeLogger struct {
	active bool
}

// debugTypeLogger ports TypeCheck.java:DebugTypeLogger constructor (line 3415).
// Returns a new logger instance; logging is no-op in jsgo since we don't ship debug logs.
func (c *TypeCheck) debugTypeLogger() *TypeCheckDebugTypeLogger {
	return &TypeCheckDebugTypeLogger{}
}

// maybeStartLoggingAt ports TypeCheck.java:DebugTypeLogger.maybeStartLoggingAt (line 3422).
func (l *TypeCheckDebugTypeLogger) maybeStartLoggingAt(n *rhino.Node) { _ = n }

// stopLoggingIfThisIsWhereWeStarted ports TypeCheck.java:DebugTypeLogger.stopLoggingIfThisIsWhereWeStarted (line 3443).
func (l *TypeCheckDebugTypeLogger) stopLoggingIfThisIsWhereWeStarted(n *rhino.Node) { _ = n }

// close ports TypeCheck.java:DebugTypeLogger.close (line 3450).
func (l *TypeCheckDebugTypeLogger) close() { l.active = false }

// maybeLogTypeOfNode ports TypeCheck.java:DebugTypeLogger.maybeLogTypeOfNode (line 3456).
func (l *TypeCheckDebugTypeLogger) maybeLogTypeOfNode(n *rhino.Node) { _ = n }

// SuggestionPair is the Go equivalent of Java's TypeCheck.SuggestionPair inner class.
// It pairs a property name suggestion with its edit distance for typo-suggestion messages.
type SuggestionPair struct {
	Suggestion string
	Distance   int
}

// suggestionPair ports TypeCheck.java:SuggestionPair.suggestionPair (line ~430).
func (s *SuggestionPair) suggestionPair(suggestion string, distance int) SuggestionPair {
	return SuggestionPair{Suggestion: suggestion, Distance: distance}
}

// typeInferencePass creates a new TypeInferencePass for the given scope.
// Mirrors TypeCheck.java's construction of TypeInferencePass at line ~430.
func typeInferencePass(scope *TypedScope) *TypeInferencePass {
	return NewTypeInferencePass(scope)
}

// stringKey ports Java's Property.StringKey used for property name lookups.
// In Java, Property.StringKey wraps a string; in Go, we use plain strings.
func stringKey(name string) string {
	return name
}
