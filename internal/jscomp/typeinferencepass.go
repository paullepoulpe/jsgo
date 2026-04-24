// Port of com/google/javascript/jscomp/TypeInferencePass.java @ closure-compiler v20260415
//
// TypeInferencePass in Java orchestrates two traversal phases
// (scope-building then type-propagation) via NodeTraversal callbacks.
// jsgo collapses that into a single transfer-function struct handed to
// DataFlowAnalysis.Run; expression-level type inference lives in the
// existing inference.go (port of TypeInference.java).
//
// Supported refinements (match buildRefinements in refinement.go):
//   - NAME truthy / falsy
//   - NAME == / != / === / !== null | undefined
//   - typeof NAME === 'prim'
//   - NAME instanceof Ctor
//   - !cond, cond1 && cond2, cond1 || cond2 (falsy branch)
//
// Assignment refinement:
//   - VAR/LET/CONST x = expr   → OUT.x = inferTypeInFlow(expr, IN)
//   - ASSIGN x = expr          → OUT.x = inferTypeInFlow(expr, IN)

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// TypeInferencePass holds the declaring scope (so it can resolve
// original declared types) and produces the TransferFunction used
// by DataFlowAnalysis.
type TypeInferencePass struct {
	scope *TypedScope
}

// Go-only: Go constructor replacing Java's package-private TypeInferencePass(compiler, rai, scopeCreator).
// jsgo collapses the full orchestration into a single transfer-function object.
func NewTypeInferencePass(scope *TypedScope) *TypeInferencePass {
	return &TypeInferencePass{scope: scope}
}

// Go-only: single Transfer entry-point; Java spreads this logic across inferAllScopes/inferScope
// and the inner DataFlowAnalysis.flowThrough callback inside TypeInference.java.
// Transfer is the method reference handed to DataFlowAnalysis.Run.
func (p *TypeInferencePass) Transfer(node *CfgNode, in *FlowScope) map[Branch]*FlowScope {
	if node == nil || node.AST == nil {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	ast := node.AST
	switch ast.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		return p.transferDecl(ast, in)
	case rhino.EXPR_RESULT:
		if a := ast.FirstChild(); a != nil && a.Token().IsAssign() {
			return p.transferAssign(a, in)
		}
	case rhino.IF:
		return p.transferIf(ast, in)
	case rhino.WHILE:
		return p.transferWhile(ast, in)
	case rhino.DO_WHILE:
		return p.transferDoWhile(ast, in)
	case rhino.FOR:
		return p.transferFor(ast, in)
	}
	return map[Branch]*FlowScope{BranchUnconditional: in}
}

// Go-only: Java TypeInference.java flowThrough() handles all node types in one method;
// jsgo splits by node type for clarity.  transferDecl handles VAR/LET/CONST with an initializer. The
// declared type flows through unchanged; only names whose
// initializer evaluates to a type sharper than the declared form
// end up in the FlowScope (declared-type fallback lives in the
// caller's GetDeclaredType lookup).
func (p *TypeInferencePass) transferDecl(decl *rhino.Node, in *FlowScope) map[Branch]*FlowScope {
	name := decl.FirstChild()
	if name == nil || name.Token() != rhino.NAME {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	init := name.FirstChild()
	if init == nil {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	t := p.inferExpr(init, in)
	if t == nil {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	// Don't overwrite a richer declared type with a weaker inferred
	// type: constructor / interface bindings shadow their FUNCTION
	// init's bare FunctionType, and `@enum` bindings shadow their
	// OBJECTLIT init's structural {...}. Matches how Java's
	// TypeInference keeps the declared slot for these.
	if declared := p.scope.GetRawDeclaredType(name.String_()); declared != nil {
		switch d := declared.(type) {
		case *jstype.FunctionType:
			if d.IsConstructor {
				return map[Branch]*FlowScope{BranchUnconditional: in}
			}
			// FunctionType whose params or return carry a declared
			// shape (not all-Unknown) — preserve it. Covers \`var
			// bar = function(a){}\` with @param on the VAR.
			if functionHasDeclaredShape(d) {
				return map[Branch]*FlowScope{BranchUnconditional: in}
			}
		case *jstype.EnumType:
			return map[Branch]*FlowScope{BranchUnconditional: in}
		case *jstype.TemplatizedType:
			// Preserve the declared generic binding — the init's
			// bare Array / Map type would drop the template args,
			// breaking \`arr[0]\`-style reads that rely on them.
			return map[Branch]*FlowScope{BranchUnconditional: in}
		}
	}
	return map[Branch]*FlowScope{BranchUnconditional: in.Bind(name.String_(), t)}
}

// Go-only: helper extracted from transferDecl; no direct Java equivalent.
// functionHasDeclaredShape reports whether any of \`fn\`'s params
// or return carry a non-Unknown type. Used to avoid clobbering a
// VAR-declared FunctionType whose signature comes from outer
// @param/@return JSDoc.
func functionHasDeclaredShape(fn *jstype.FunctionType) bool {
	if fn == nil {
		return false
	}
	for _, p := range fn.Params {
		if p != nil && p != jstype.Unknown {
			return true
		}
	}
	if fn.Return != nil && fn.Return != jstype.Unknown && fn.Return != jstype.Void {
		return true
	}
	return false
}

// Go-only: split from Transfer for ASSIGN nodes; mirrors Java TypeInference.updateScopeForTypeChange.
// transferAssign handles ASSIGN (and compound-assign variants) where
// the LHS is a bare NAME or a GETPROP qualified name.
//
// GETPROP handling: `t.x = expr` narrows the qualified name "t.x" in
// the output FlowScope via inferQualifiedSlot. This matches Java's
// TypeInference.updateScopeForTypeChange which calls
// scope.inferQualifiedSlot for GETPROP LHS nodes so that subsequent
// reads of t.x see the narrowed type rather than the declared union.
func (p *TypeInferencePass) transferAssign(a *rhino.Node, in *FlowScope) map[Branch]*FlowScope {
	kids := a.Children()
	if len(kids) < 2 {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	lhs, rhs := kids[0], kids[1]
	switch lhs.Token() {
	case rhino.NAME:
		t := p.inferExpr(rhs, in)
		if t == nil {
			return map[Branch]*FlowScope{BranchUnconditional: in}
		}
		// Don't widen past the declared type. If the RHS-inferred type
		// is a proper superset of the declared slot (e.g. declared
		// `number` being fed `null|number`), the declared type wins —
		// TypeCheck's subtype check on the ASSIGN itself would fire if
		// the RHS truly couldn't fit. Matches Java's refuse-to-widen
		// rule on named-slot assignments.
		if declared := p.scope.GetRawDeclaredType(lhs.String_()); declared != nil {
			if !jstype.IsSubtype(t, declared) {
				return map[Branch]*FlowScope{BranchUnconditional: in}
			}
		}
		return map[Branch]*FlowScope{BranchUnconditional: in.inferSlotType(lhs.String_(), t)}
	case rhino.GETPROP:
		// Qualified-name assignment: `obj.prop = expr` narrows
		// "obj.prop" in the output scope. Port of TypeInference.java
		// updateScopeForTypeChange for the GETPROP branch which calls
		// scope.inferQualifiedSlot(qualifiedName, inferredType).
		//
		// Guards (all must pass before narrowing):
		//   1. No JSDoc type annotation on the ASSIGN — e.g.
		//      `/** @constructor */ ns.Foo = ...` should not be narrowed
		//      since TypedScopeCreator controls the declared slot type.
		//   2. No static namespace declaration for this prop via
		//      LookupNamespaceProp (catches @type-annotated ns props).
		// Together these avoid clobbering constructor/enum declarations
		// with the RHS-inferred type, which would corrupt property checks.
		info := a.JSDocInfo()
		if info != nil && (info.IsConstructor() || info.IsInterface() ||
			info.IsRecord() || info.IsEnum() || info.HasType()) {
			return map[Branch]*FlowScope{BranchUnconditional: in}
		}
		lhsKids := lhs.Children()
		if len(lhsKids) >= 2 {
			// GETPROP property key is a STRINGLIT node in this port's
			// AST (the parser builds: gp.AddChildToBack(rhino.NewString(
			// rhino.STRINGLIT, ident))). Use String_() regardless of
			// the token type.
			prop := lhsKids[1].String_()
			// `X.prototype = …` is handled by TypedScopeCreator
			// (recordPrototypeAssign). Don't record it in the FlowScope
			// or the constructor's instance type (which may be @struct
			// or @dict) would be pinned to `X.prototype`, causing
			// spurious struct/dict errors on subsequent property additions
			// to the prototype. Matches Java's TypeInference which skips
			// the prototype slot in updateScopeForTypeChange.
			if prop == "prototype" {
				return map[Branch]*FlowScope{BranchUnconditional: in}
			}
			if lhsKids[0].Token() == rhino.NAME {
				owner := lhsKids[0].String_()
				if p.scope.LookupNamespaceProp(owner, prop) != nil {
					// Has a static namespace declaration — don't narrow.
					return map[Branch]*FlowScope{BranchUnconditional: in}
				}
			}
		}
		qname := qualifiedName(lhs)
		if qname == "" {
			return map[Branch]*FlowScope{BranchUnconditional: in}
		}
		t := p.inferExpr(rhs, in)
		if t == nil {
			return map[Branch]*FlowScope{BranchUnconditional: in}
		}
		// Don't narrow to All / Unknown — those types carry no information
		// and would discard the declared type of a slot that was
		// previously narrower. E.g. `this.a = unknown_param` should not
		// override `this.a : string` with `?`. Mirrors Java's
		// updateScopeForTypeChange guard: it skips the slot update when
		// the inferred type is wider than the declared type.
		if t == jstype.All || t == jstype.Unknown {
			return map[Branch]*FlowScope{BranchUnconditional: in}
		}
		return map[Branch]*FlowScope{BranchUnconditional: in.inferQualifiedSlot(qname, t)}
	default:
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
}

// Go-only: split from Transfer for IF nodes; Java uses ReverseAbstractInterpreter inline.
// transferIf produces two branch-specific FlowScopes: the truthy
// branch's refinements on BranchOnTrue, the falsy branch's on
// BranchOnFalse. The AST child `kids[0]` is the condition; the then
// and else blocks are walked separately by the CFG structure.
func (p *TypeInferencePass) transferIf(ast *rhino.Node, in *FlowScope) map[Branch]*FlowScope {
	kids := ast.Children()
	if len(kids) < 2 {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	cond := kids[0]
	trueRefs := p.branchRefinements(cond, in, true)
	falseRefs := p.branchRefinements(cond, in, false)
	return map[Branch]*FlowScope{
		BranchOnTrue:  in.BindAll(trueRefs),
		BranchOnFalse: in.BindAll(falseRefs),
	}
}

// Go-only: split from Transfer for WHILE nodes; mirrors Java TypeInference.flowThrough WHILE case.
// transferWhile produces branch-specific FlowScopes for a WHILE loop.
// Children are [cond, body]. The loop body executes when the condition
// is truthy (BranchOnTrue), and the loop exits when it is falsy
// (BranchOnFalse). Mirrors how Java's TypeInference.flowThrough treats
// WHILE nodes as conditional branches in the DataFlowAnalysis.
func (p *TypeInferencePass) transferWhile(ast *rhino.Node, in *FlowScope) map[Branch]*FlowScope {
	kids := ast.Children()
	if len(kids) < 1 {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	cond := kids[0]
	trueRefs := p.branchRefinements(cond, in, true)
	falseRefs := p.branchRefinements(cond, in, false)
	return map[Branch]*FlowScope{
		BranchOnTrue:  in.BindAll(trueRefs),
		BranchOnFalse: in.BindAll(falseRefs),
		BranchSyn:     in, // back-edge carries unrestricted scope
	}
}

// Go-only: split from Transfer for DO_WHILE nodes; mirrors Java TypeInference.flowThrough DO_WHILE case.
// transferDoWhile produces branch-specific FlowScopes for a DO_WHILE
// loop node. Children are [body, cond]. The condition is at the bottom
// of the loop, so the CFG node for DO_WHILE is the condition site:
// BranchOnTrue loops back into the body, BranchOnFalse exits.
// Port of TypeInference.java's DO_WHILE case in flowThrough.
func (p *TypeInferencePass) transferDoWhile(ast *rhino.Node, in *FlowScope) map[Branch]*FlowScope {
	kids := ast.Children()
	if len(kids) < 2 {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	cond := kids[len(kids)-1] // condition is the last child
	trueRefs := p.branchRefinements(cond, in, true)
	falseRefs := p.branchRefinements(cond, in, false)
	return map[Branch]*FlowScope{
		BranchOnTrue:  in.BindAll(trueRefs),
		BranchOnFalse: in.BindAll(falseRefs),
	}
}

// Go-only: split from Transfer for FOR nodes; mirrors Java TypeInference.flowThrough FOR case.
// transferFor produces branch-specific FlowScopes for a classic FOR
// loop. Children are [init, cond, update, body]. The FOR node acts as
// the condition-check site in the CFG: BranchOnTrue enters the body,
// BranchOnFalse exits. An absent (empty) condition is always truthy
// (infinite loop); in that case we pass through unconditionally.
// Port of TypeInference.java's FOR case in flowThrough.
func (p *TypeInferencePass) transferFor(ast *rhino.Node, in *FlowScope) map[Branch]*FlowScope {
	kids := ast.Children()
	if len(kids) < 2 {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	cond := kids[1]
	// A nil condition (omitted: `for (;;)`) means always-true;
	// pass through without branch refinement.
	if cond == nil {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	trueRefs := p.branchRefinements(cond, in, true)
	falseRefs := p.branchRefinements(cond, in, false)
	return map[Branch]*FlowScope{
		BranchOnTrue:  in.BindAll(trueRefs),
		BranchOnFalse: in.BindAll(falseRefs),
		BranchSyn:     in, // update-node back-edge
	}
}

// Go-only: helper that wraps buildRefinements with flow-scope push/pop; no direct Java counterpart.
// branchRefinements wraps buildRefinements with an IN-aware scope so
// typeof / instanceof / null-check narrowing see the current flow
// types rather than just the declared ones.
func (p *TypeInferencePass) branchRefinements(cond *rhino.Node, in *FlowScope, whenTruthy bool) map[string]jstype.JSType {
	// Temporarily push the IN FlowScope as a refinement layer so
	// buildRefinements's scope.GetDeclaredType returns the flow
	// type. Pop on exit.
	if in != nil && len(in.types) > 0 {
		p.scope.PushRefinement(in.types)
		defer p.scope.PopRefinement()
	}
	return buildRefinements(p.scope, cond, whenTruthy)
}

// Go-only: delegates to inference.go (TypeInference.java port); no direct method in TypeInferencePass.java.
// inferExpr runs expression-level TypeInference against the given
// FlowScope — the flow-narrowed name map takes precedence over the
// declared type. Same contract as the old TypeInference.InferType
// but it consults `in` first.
func (p *TypeInferencePass) inferExpr(n *rhino.Node, in *FlowScope) jstype.JSType {
	if in != nil && len(in.types) > 0 {
		p.scope.PushRefinement(in.types)
		defer p.scope.PopRefinement()
	}
	ti := NewTypeInference(p.scope)
	return ti.InferType(n)
}
