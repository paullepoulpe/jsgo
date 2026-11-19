// Clean-room dead_cast_cleanup pass: eliminates redundant T(x) and AsIface[T](x) self-conversions.

package passes

import (
	"go/ast"
	"go/token"
	"go/types"
)

// DeadCastCleanupPass eliminates redundant self-type conversions `T(x)` and `j2gort.AsIface[T](x)`
// where `go/types` proves `x` already has exact standalone type `T` (and for `AsIface`, cannot be typed-nil).
type DeadCastCleanupPass struct{}

// NewDeadCastCleanupPass constructs a new dead_cast_cleanup pass instance.
func NewDeadCastCleanupPass() *DeadCastCleanupPass {
	return &DeadCastCleanupPass{}
}

// Name returns "dead_cast_cleanup".
func (p *DeadCastCleanupPass) Name() string {
	return "dead_cast_cleanup"
}

// Description returns a summary of the pass.
func (p *DeadCastCleanupPass) Description() string {
	return "Eliminate redundant self-type conversions T(x) and AsIface[T](x) where x already has exact type T"
}

// Apply runs dead_cast_cleanup across target files in ctx.
func (p *DeadCastCleanupPass) Apply(ctx *PassContext) (int, error) {
	total := 0
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		n := p.applyFile(file, ctx.Pkg, ctx.Info)
		if n > 0 {
			ctx.MarkFileModified(file)
			total += n
		}
	}
	return total, nil
}

func (p *DeadCastCleanupPass) applyFile(file *ast.File, pkg *types.Package, info *types.Info) int {
	parents, postOrder := buildParentMap(file)
	var facts *peepholeFileFacts
	count := 0

	for _, node := range postOrder {
		call, ok := node.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
			continue
		}

		parent := parents[call]
		if parent == nil {
			continue
		}

		arg := call.Args[0]
		inner := stripParens(arg)
		canStrip := false

		if funTV, ok := info.Types[call.Fun]; ok && funTV.IsType() && funTV.Type != nil {
			targetType := funTV.Type
			if !isInvalidOrUntyped(targetType) && hasExactIntrinsicType(arg, targetType, info) {
				canStrip = true
			}
		} else if !isJ2gortPkg(pkg) {
			if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && fnName == "AsIface" {
				callType := info.Types[call].Type
				if isInterfaceType(callType) && hasExactIntrinsicType(inner, callType, info) {
					if facts == nil {
						facts = collectPeepholeFileFacts(file, info, parents)
					}
					if isGuaranteedNeverTypedNil(inner, facts, info) {
						canStrip = true
					}
				}
			}
		}

		if !canStrip {
			continue
		}

		// Expression-statement guard: only function calls and channel receives are valid Go ExprStmts.
		if _, isExprStmt := parent.(*ast.ExprStmt); isExprStmt && !isValidExprStmtTarget(inner) {
			continue
		}

		replacement := inner
		if needsParensInContext(replacement, call, parent, parents) {
			replacement = &ast.ParenExpr{X: replacement}
		}

		if replaceExprInParent(parent, call, replacement, parents) {
			count++
		}
	}
	return count
}

// hasExactIntrinsicType reports whether expr has exact standalone type targetType even when
// the surrounding conversion `targetType(expr)` is removed.
func hasExactIntrinsicType(expr ast.Expr, targetType types.Type, info *types.Info) bool {
	tv, ok := info.Types[expr]
	if !ok || tv.IsType() || tv.IsBuiltin() || tv.IsNil() || tv.Type == nil {
		return false
	}
	if isInvalidOrUntyped(tv.Type) || !types.Identical(targetType, tv.Type) {
		return false
	}

	// For constant expressions (`tv.Value != nil`), `info.Types[c].Type` records the converted
	// target type rather than untyped; only strip `T(...)` if the operand carries type `T` on its own.
	if tv.Value != nil {
		inner := stripParens(expr)
		switch x := inner.(type) {
		case *ast.Ident:
			cObj, ok := info.Uses[x].(*types.Const)
			return ok && !isInvalidOrUntyped(cObj.Type()) && types.Identical(targetType, cObj.Type())
		case *ast.SelectorExpr:
			cObj, ok := info.Uses[x.Sel].(*types.Const)
			return ok && !isInvalidOrUntyped(cObj.Type()) && types.Identical(targetType, cObj.Type())
		case *ast.CallExpr:
			innerFunTV, ok := info.Types[x.Fun]
			return ok && innerFunTV.IsType() && !isInvalidOrUntyped(innerFunTV.Type) && types.Identical(targetType, innerFunTV.Type)
		default:
			return false
		}
	}
	return true
}

func isInvalidOrUntyped(t types.Type) bool {
	if t == nil {
		return true
	}
	b, ok := types.Unalias(t).(*types.Basic)
	if !ok {
		return false
	}
	return b.Kind() == types.Invalid || (b.Info()&types.IsUntyped) != 0
}

func isValidExprStmtTarget(expr ast.Expr) bool {
	switch u := expr.(type) {
	case *ast.CallExpr:
		return true
	case *ast.UnaryExpr:
		return u.Op == token.ARROW
	default:
		return false
	}
}

func needsParensInContext(inner ast.Expr, oldCall *ast.CallExpr, parent ast.Node, parents map[ast.Node]ast.Node) bool {
	if _, isComp := inner.(*ast.CompositeLit); isComp && inControlHeader(oldCall, parents) {
		return true
	}

	switch innerExpr := inner.(type) {
	case *ast.UnaryExpr, *ast.StarExpr:
		switch p := parent.(type) {
		case *ast.SelectorExpr, *ast.IndexExpr, *ast.IndexListExpr, *ast.SliceExpr, *ast.TypeAssertExpr,
			*ast.UnaryExpr, *ast.StarExpr:
			return true
		case *ast.CallExpr:
			return p.Fun == oldCall
		}
	case *ast.BinaryExpr:
		switch p := parent.(type) {
		case *ast.SelectorExpr, *ast.IndexExpr, *ast.IndexListExpr, *ast.SliceExpr, *ast.TypeAssertExpr,
			*ast.UnaryExpr, *ast.StarExpr:
			return true
		case *ast.CallExpr:
			return p.Fun == oldCall
		case *ast.BinaryExpr:
			if p.Op.Precedence() > innerExpr.Op.Precedence() {
				return true
			}
			if p.Op.Precedence() == innerExpr.Op.Precedence() && p.Y == oldCall {
				return true
			}
		}
	}
	return false
}

func inControlHeader(node ast.Node, parents map[ast.Node]ast.Node) bool {
	for cur := node; ; {
		p := parents[cur]
		if p == nil {
			return false
		}
		switch s := p.(type) {
		case *ast.ParenExpr, *ast.CallExpr, *ast.IndexExpr, *ast.IndexListExpr, *ast.BlockStmt:
			return false
		case *ast.IfStmt:
			return cur != s.Body && cur != s.Else
		case *ast.ForStmt:
			return cur != s.Body
		case *ast.RangeStmt:
			return cur != s.Body
		case *ast.SwitchStmt:
			return cur != s.Body
		case *ast.TypeSwitchStmt:
			return cur != s.Body
		}
		cur = p
	}
}
