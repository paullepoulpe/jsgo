// Port of com/google/javascript/jscomp/TypeInference.java (flow-
// transfer-function half) @ closure-compiler v20260415 (narrow).
//
// The `TypeInference` pass in Java extends DataFlowAnalysis<Node,
// FlowScope> and supplies a transfer function that computes a new
// FlowScope from the old one + the current CFG node. jsgo's port
// keeps the same shape: TypeInferencePass is a transfer function
// passed to DataFlowAnalysis.Run, while expression-level type
// inference lives in the existing inference.go.
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

func NewTypeInferencePass(scope *TypedScope) *TypeInferencePass {
	return &TypeInferencePass{scope: scope}
}

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
	}
	return map[Branch]*FlowScope{BranchUnconditional: in}
}

// transferDecl handles VAR/LET/CONST with an initializer. The
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
	// Don't overwrite a constructor/interface binding created by the
	// scope-creator's `/** @constructor */ var Foo = function(){}`
	// branch — the init (a bare FUNCTION) infers to a plain
	// FunctionType that would otherwise shadow the richer
	// ConstructorType. Java's TypeInference keeps the declared slot
	// when the declaration is a named constructor / interface.
	if declared := p.scope.GetRawDeclaredType(name.String_()); declared != nil {
		if fn, ok := declared.(*jstype.FunctionType); ok && fn.IsConstructor {
			return map[Branch]*FlowScope{BranchUnconditional: in}
		}
	}
	return map[Branch]*FlowScope{BranchUnconditional: in.Bind(name.String_(), t)}
}

// transferAssign handles ASSIGN (and compound-assign variants) where
// the LHS is a bare NAME. GETPROP assignments don't refine the flow
// scope at this narrow tier.
func (p *TypeInferencePass) transferAssign(a *rhino.Node, in *FlowScope) map[Branch]*FlowScope {
	kids := a.Children()
	if len(kids) < 2 {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	lhs, rhs := kids[0], kids[1]
	if lhs.Token() != rhino.NAME {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	t := p.inferExpr(rhs, in)
	if t == nil {
		return map[Branch]*FlowScope{BranchUnconditional: in}
	}
	return map[Branch]*FlowScope{BranchUnconditional: in.Bind(lhs.String_(), t)}
}

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
