// Clean-room j2gort_peephole pass and shared j2gort AST/type analysis helpers.
// Eliminates redundant mechanical j2gort boxing/casting/null-check wrappers while
// strictly preserving NullString ("\xff\x00j2go_null"), uint32 Node(0), and typed-nil
// interface semantics.

package passes

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"slices"
	"strconv"
	"strings"
)

const j2gortNullString = "\xff\x00j2go_null"

// J2gortPeepholePass eliminates redundant j2gort cast, interface-adaptation, boxing/unboxing,
// and non-null string/pointer/interface wrapper calls emitted mechanically by j2go.
type J2gortPeepholePass struct{}

// NewJ2gortPeepholePass constructs a new j2gort_peephole pass instance.
func NewJ2gortPeepholePass() *J2gortPeepholePass {
	return &J2gortPeepholePass{}
}

// Name returns "j2gort_peephole".
func (p *J2gortPeepholePass) Name() string {
	return "j2gort_peephole"
}

// Description returns a summary of the pass.
func (p *J2gortPeepholePass) Description() string {
	return "Eliminate redundant j2gort Cast/AsIface/Box/StrVal/CheckNotNull*/StringEquals wrappers"
}

type peepholeExports struct {
	hasCheckNotNullString bool
	hasCheckNotNullPtr    bool
	hasCheckNotNullSlice  bool
	hasCheckNotNullNode   bool
}

// Apply runs j2gort_peephole across target files in ctx.
func (p *J2gortPeepholePass) Apply(ctx *PassContext) (int, error) {
	if isJ2gortPkg(ctx.Pkg) {
		return 0, nil
	}
	exports := peepholeExports{
		hasCheckNotNullString: j2gortHasExport(ctx, "CheckNotNullString"),
		hasCheckNotNullPtr:    j2gortHasExport(ctx, "CheckNotNullPtr"),
		hasCheckNotNullSlice:  j2gortHasExport(ctx, "CheckNotNullSlice"),
		hasCheckNotNullNode:   j2gortHasExport(ctx, "CheckNotNullNode"),
	}
	total := 0
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		if n := p.rewriteLazyPreconditions(file, ctx.Pkg, ctx.Info); n > 0 {
			ctx.MarkFileModified(file)
			total += n
		}
	}
	var safeBitSetParams map[bitSetParamKey]bool
	computedSafeParams := false
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		if !computedSafeParams && fileHasBitSetCtorCall(file) {
			safeBitSetParams = collectSafeBitSetParamFuncs(ctx.Files, ctx.Pkg, ctx.Info)
			computedSafeParams = true
		}
		var fileSafeParams map[bitSetParamKey]bool
		if computedSafeParams {
			fileSafeParams = safeBitSetParams
		}
		n := p.applyFile(file, ctx.Pkg, ctx.Info, fileSafeParams, exports)
		if n > 0 {
			ctx.MarkFileModified(file)
			total += n
		}
	}
	lihCount, err := p.rewriteLinkedIdentityHashMap(ctx)
	if err != nil {
		return total, err
	}
	total += lihCount
	total += p.elideZeroedInitFieldsWrites(ctx)
	total += p.rewriteNodeQualifiedNameJoinDot(ctx)
	total += p.rewriteNodeClearJSTypeAndJSDocFastPaths(ctx)
	total += p.rewriteNodeHasNonNegativeLength(ctx)
	total += p.rewriteIfCondIIFEs(ctx)
	total += p.rewriteStmtValueIIFEs(ctx)
	total += p.hoistStatelessNodeLambdas(ctx)
	nuHasCount, err := p.rewriteNodeUtilHasNodeLambdas(ctx)
	if err != nil {
		return total, err
	}
	total += nuHasCount
	total += p.rewriteStackNodeTraversalBuilders(ctx)
	return total, nil
}

func (p *J2gortPeepholePass) applyFile(file *ast.File, pkg *types.Package, info *types.Info, safeBitSetParams map[bitSetParamKey]bool, exports peepholeExports) int {
	parents, postOrder := buildParentMap(file)
	facts := collectPeepholeFileFacts(file, info, parents)
	jsEqMeta := resolveJSTypeEqualsMeta(pkg)
	count := 0

	applyReplacement := func(parent ast.Node, call *ast.CallExpr, replacement ast.Expr) bool {
		if replacement == nil {
			return false
		}
		if _, isExprStmt := parent.(*ast.ExprStmt); isExprStmt && !isValidExprStmtTarget(replacement) {
			return false
		}
		if needsParensInContext(replacement, call, parent, parents) {
			pe := &ast.ParenExpr{X: replacement}
			if tv, ok := info.Types[replacement]; ok {
				info.Types[pe] = tv
			}
			replacement = pe
		}
		return replaceExprInParent(parent, call, replacement, parents)
	}

	for _, node := range postOrder {
		call, ok := node.(*ast.CallExpr)
		if !ok || call.Ellipsis.IsValid() {
			continue
		}
		parent := parents[call]
		if parent == nil {
			continue
		}

		// Check `ToIface_<Name>(x)` adapters (same-package, `<pkg>.ToIface_*`, or `j2gort.ToIface_*`).
		if targetIface, isTI := isToIfaceCall(call, info); isTI {
			inner := stripParens(call.Args[0])
			if es, isExprStmt := parent.(*ast.ExprStmt); isExprStmt && es.X == call {
				if isValidExprStmtTarget(inner) && replaceExprInParent(parent, call, inner, parents) {
					count++
				}
				continue
			}
			if repl, ok := p.simplifyToIfaceCall(call, inner, targetIface, parent, parents, pkg, facts, info); ok && applyReplacement(parent, call, repl) {
				count++
			}
			continue
		}

		// Specialize `<JSType>.Equals(<JSType>)` -> `<*JSTypeBase>.EqualsJSType(<JSType>)` when available.
		if jsEqMeta != nil {
			if repl, ok := p.tryRewriteJSTypeEquals(call, file, pkg, facts, jsEqMeta, info); ok && applyReplacement(parent, call, repl) {
				count++
				continue
			}
		}

		fnName, fnObj, _, isJ2 := isJ2gortCall(call, info)
		if !isJ2 {
			continue
		}

		if fnName == "StringEquals" && len(call.Args) == 2 {
			if p.tryRewriteStringEquals(call, parent, parents, info) {
				count++
			}
			continue
		}
		if fnName == "IfaceEquals" && len(call.Args) == 2 {
			if p.tryRewriteIfaceEquals(call, parent, parents, pkg, info) {
				count++
			}
			continue
		}
		if fnName == "StrConcat" {
			if p.tryRewriteStrConcat(call, fnObj, parent, parents, info) {
				count++
			}
			continue
		}
		if len(call.Args) != 1 {
			continue
		}
		inner := stripParens(call.Args[0])

		// Discarded expression-statement cast/adaptation wrappers.
		if es, isExprStmt := parent.(*ast.ExprStmt); isExprStmt && es.X == call {
			switch fnName {
			case "Cast", "UncheckedCast", "CastPtr", "AsIface", "CastString", "CastNode",
				"ToAny", "PtrToAny", "SliceToAny", "BoxString", "BoxNode":
				if isValidExprStmtTarget(inner) && replaceExprInParent(parent, call, inner, parents) {
					count++
				}
				continue
			}
		}

		if repl, ok := p.simplifySingleArgCall(fnName, fnObj, call, inner, parent, parents, pkg, facts, exports, info); ok && applyReplacement(parent, call, repl) {
			count++
		}
	}

	count += p.rewriteSliceIteratorLoops(file, parents, facts, info)
	count += p.rewriteNodeGetAncestorsLoops(file, pkg, info)
	count += p.rewriteLazyLogFileStmts(file, pkg, info)
	count += p.rewriteStackBitSets(file, pkg, safeBitSetParams, info)
	count += p.rewriteStackUnionTypeBuilders(file, pkg, info)
	count += p.rewriteAstNodeDeserializerRelease(file, pkg, info)

	return count
}

func (p *J2gortPeepholePass) simplifySingleArgCall(
	fnName string,
	fnObj *types.Func,
	call *ast.CallExpr,
	inner ast.Expr,
	parent ast.Node,
	parents map[ast.Node]ast.Node,
	pkg *types.Package,
	facts *peepholeFileFacts,
	exports peepholeExports,
	info *types.Info,
) (ast.Expr, bool) {
	callType := info.Types[call].Type
	if callType == nil {
		return nil, false
	}

	switch fnName {
	case "CastNode":
		if innerCall, ok := inner.(*ast.CallExpr); ok && len(innerCall.Args) == 1 && !innerCall.Ellipsis.IsValid() {
			if innerFn, _, _, ok := isJ2gortCall(innerCall, info); ok && innerFn == "BoxNode" {
				x := stripParens(innerCall.Args[0])
				if hasExactIntrinsicType(x, callType, info) {
					return x, true
				}
			}
		}

	case "CastString":
		if innerCall, ok := inner.(*ast.CallExpr); ok && len(innerCall.Args) == 1 && !innerCall.Ellipsis.IsValid() {
			if innerFn, _, _, ok := isJ2gortCall(innerCall, info); ok && innerFn == "BoxString" {
				s := stripParens(innerCall.Args[0])
				if hasExactOrDefaultStringType(s, info) {
					return s, true
				}
			}
		}

	case "UnboxBoolean", "UnboxInt", "UnboxLong", "UnboxDouble", "UnboxFloat", "UnboxChar", "UnboxByte", "UnboxShort":
		return unwrapBoxRoundTrip(fnName, inner, info)

	case "CastPtr", "Cast", "UncheckedCast":
		_, isPtr := types.Unalias(callType).(*types.Pointer)
		_, isSlice := types.Unalias(callType).(*types.Slice)
		_, isStruct := types.Unalias(callType).Underlying().(*types.Struct)
		if isPtr || (fnName != "CastPtr" && (isSlice || isStruct)) {
			if hasExactIntrinsicType(inner, callType, info) {
				return inner, true
			}
			if innerCall, ok := inner.(*ast.CallExpr); ok && len(innerCall.Args) == 1 && !innerCall.Ellipsis.IsValid() {
				if innerFn, _, _, ok := isJ2gortCall(innerCall, info); ok {
					wantWrap := "PtrToAny"
					if isSlice {
						wantWrap = "SliceToAny"
					} else if isStruct {
						wantWrap = "ToAny"
					}
					if innerFn == "ToAny" || innerFn == wantWrap {
						if x := stripParens(innerCall.Args[0]); hasExactIntrinsicType(x, callType, info) {
							return x, true
						}
					}
				}
			}
		}
		if innerCall, ok := inner.(*ast.CallExpr); ok && len(innerCall.Args) == 1 && !innerCall.Ellipsis.IsValid() {
			if innerFn, _, _, ok := isJ2gortCall(innerCall, info); ok {
				sameCast := innerFn == fnName || (fnName == "CastPtr" && (innerFn == "Cast" || innerFn == "UncheckedCast"))
				if sameCast && types.Identical(callType, info.Types[innerCall].Type) {
					return innerCall, true
				}
			}
		}
		if fnName != "CastPtr" && isNonEmptyInterfaceType(callType) {
			innerTV, ok := info.Types[inner]
			if ok && innerTV.Type != nil && !isInvalidOrUntyped(innerTV.Type) &&
				types.AssignableTo(innerTV.Type, callType) &&
				isGuaranteedNeverTypedNil(inner, facts, info) &&
				exprContextPreservesInterfaceType(call, callType, innerTV.Type, parent, parents, pkg, info) {
				return inner, true
			}
		}

	case "AsIface":
		if !isInterfaceType(callType) {
			return nil, false
		}
		if innerCall, ok := inner.(*ast.CallExpr); ok && len(innerCall.Args) == 1 && !innerCall.Ellipsis.IsValid() {
			_, innerIsJ2AsIface := false, false
			if innerFn, _, _, ok := isJ2gortCall(innerCall, info); ok && innerFn == "AsIface" {
				innerIsJ2AsIface = true
			}
			_, innerIsUpcast := isUpcastCall(innerCall, info)
			if innerIsJ2AsIface || innerIsUpcast {
				innerType := info.Types[innerCall].Type
				if innerType != nil {
					if types.Identical(callType, innerType) {
						return innerCall, true
					}
					y := stripParens(innerCall.Args[0])
					yTV, okY := info.Types[y]
					if okY && yTV.Type != nil && !isInvalidOrUntyped(yTV.Type) &&
						isInterfaceType(innerType) &&
						types.AssignableTo(innerType, callType) &&
						types.AssignableTo(yTV.Type, innerType) {
						if isGuaranteedNeverTypedNil(y, facts, info) &&
							exprContextPreservesInterfaceType(call, callType, yTV.Type, parent, parents, pkg, info) {
							return y, true
						}
						if replUp, okUp := tryMakeUpcastCall(call, callType, y, yTV.Type, info); okUp {
							return replUp, true
						}
						collapsed := &ast.CallExpr{
							Fun:    call.Fun,
							Lparen: call.Lparen,
							Args:   []ast.Expr{y},
							Rparen: call.Rparen,
						}
						info.Types[collapsed] = info.Types[call]
						return collapsed, true
					}
				}
			}
		}

		innerTV, ok := info.Types[inner]
		if ok && innerTV.Type != nil && !isInvalidOrUntyped(innerTV.Type) &&
			types.AssignableTo(innerTV.Type, callType) &&
			isGuaranteedNeverTypedNil(inner, facts, info) &&
			exprContextPreservesInterfaceType(call, callType, innerTV.Type, parent, parents, pkg, info) {
			return inner, true
		}
		if ok && !innerTV.IsType() && !innerTV.IsBuiltin() && !innerTV.IsNil() && innerTV.Type != nil && !isInvalidOrUntyped(innerTV.Type) {
			if replUp, okUp := tryMakeUpcastCall(call, callType, inner, innerTV.Type, info); okUp {
				return replUp, true
			}
		}

	case "StrVal", "StrValString":
		if hasExactOrDefaultStringType(inner, info) && isGuaranteedNonNullString(inner, info) {
			return inner, true
		}

	case "CheckNotNull", "CheckNotNullPtr", "CheckNotNullString", "CheckNotNullSlice", "CheckNotNullNode", "CheckNotNullIface", "CheckNotNullGeneric":
		if innerCall, ok := inner.(*ast.CallExpr); ok && len(innerCall.Args) == 1 && !innerCall.Ellipsis.IsValid() {
			if innerFn, _, _, ok := isJ2gortCall(innerCall, info); ok && isCheckNotNullFuncName(innerFn) {
				if hasExactIntrinsicType(innerCall, callType, info) {
					return innerCall, true
				}
			}
		}
		_, isExprStmt := parent.(*ast.ExprStmt)
		if !isExprStmt || isValidExprStmtTarget(inner) {
			if types.Identical(callType, types.Typ[types.String]) {
				if hasExactOrDefaultStringType(inner, info) && isGuaranteedNonNullString(inner, info) {
					return inner, true
				}
			} else if isNamedUint32NodeType(callType) {
				if hasExactIntrinsicType(inner, callType, info) && isGuaranteedNonZeroNode(inner, facts, info) {
					return inner, true
				}
			} else if hasExactIntrinsicType(inner, callType, info) && isGuaranteedNonNilForCheckNotNull(inner, facts, info) {
				return inner, true
			}
		}
		if fnName == "CheckNotNull" {
			if spec, ok := specializeCheckNotNullCall(call, inner, callType, fnObj, isExprStmt, exports, info); ok {
				return spec, true
			}
		}

	case "BoxString":
		if hasExactOrDefaultStringType(inner, info) &&
			isGuaranteedNonNullString(inner, info) &&
			isInEmptyInterfaceSlot(call, parent, parents, info) {
			return inner, true
		}

	case "PtrToAny", "SliceToAny":
		if innerTV, ok := info.Types[inner]; ok && innerTV.Type != nil && !isInvalidOrUntyped(innerTV.Type) {
			_, isPtr := types.Unalias(innerTV.Type).(*types.Pointer)
			_, isSlice := types.Unalias(innerTV.Type).(*types.Slice)
			if ((fnName == "PtrToAny" && isPtr) || (fnName == "SliceToAny" && isSlice)) &&
				isProvablyNonNilExpr(inner, facts, info) &&
				isInEmptyInterfaceSlot(call, parent, parents, info) {
				return inner, true
			}
		}

	case "ToAny":
		if (isGuaranteedSafeForToAny(inner, info) || isProvablyNonNilExpr(inner, facts, info)) &&
			isInEmptyInterfaceSlot(call, parent, parents, info) {
			return inner, true
		}

	case "ImmListToSliceNullable", "ImmListToNodeSliceNullable", "ImmListCopyToSlice", "ImmListCopyNodeToSlice":
		if sl, okS := types.Unalias(callType).(*types.Slice); okS && sl.Elem() != nil {
			srcInner := unwrapCollectionAdaptersDown(inner, info)
			if gCall, okC := srcInner.(*ast.CallExpr); okC {
				if extGInfo, extElemT, recvExpr, okE := matchExternalSliceGetterCall(gCall, pkg, info); okE && types.Identical(sl.Elem(), extElemT) {
					pkgQual := sliceIterJ2gortQualifier(call)
					if pkgQual != "" {
						pos := call.Pos()
						j2Sel := func(name string) *ast.SelectorExpr {
							return &ast.SelectorExpr{X: &ast.Ident{NamePos: pos, Name: pkgQual}, Sel: &ast.Ident{NamePos: pos, Name: name}}
						}
						recvT := info.Types[recvExpr].Type
						target := recvExpr
						nonNil := isProvablyNonNilExpr(recvExpr, facts, info)
						if recvT != nil && isNonEmptyInterfaceType(recvT) {
							checked := recvExpr
							if !nonNil {
								chkName := "CheckNotNull"
								if hasZeroArgBoolMethod(recvT, pkg, "IsJavaNil_") {
									chkName = "CheckNotNullIface"
								}
								checked = &ast.CallExpr{Fun: j2Sel(chkName), Args: []ast.Expr{recvExpr}}
							}
							if extGInfo.directOnIface {
								target = checked
							} else if extGInfo.baseGetter != "" {
								target = &ast.CallExpr{Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(checked), Sel: &ast.Ident{NamePos: pos, Name: extGInfo.baseGetter}}}
							}
						} else if recvT != nil && !nonNil {
							if _, isPtr := types.Unalias(recvT).(*types.Pointer); isPtr {
								target = &ast.CallExpr{Fun: j2Sel("CheckNotNullPtr"), Args: []ast.Expr{recvExpr}}
							}
						}
						sliceCall := ast.Expr(&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(target), Sel: &ast.Ident{NamePos: pos, Name: extGInfo.sliceHelper}},
						})
						if fnName == "ImmListCopyToSlice" || fnName == "ImmListCopyNodeToSlice" {
							sliceCall = &ast.CallExpr{Fun: j2Sel("CheckNotNullSlice"), Args: []ast.Expr{sliceCall}}
						}
						return sliceCall, true
					}
				}
			}
		}
	}

	return nil, false
}

func (p *J2gortPeepholePass) simplifyToIfaceCall(
	call *ast.CallExpr,
	inner ast.Expr,
	targetIface types.Type,
	parent ast.Node,
	parents map[ast.Node]ast.Node,
	pkg *types.Package,
	facts *peepholeFileFacts,
	info *types.Info,
) (ast.Expr, bool) {
	innerTV, ok := info.Types[inner]
	if !ok || innerTV.IsType() || innerTV.IsBuiltin() || innerTV.IsNil() || innerTV.Type == nil || isInvalidOrUntyped(innerTV.Type) {
		return nil, false
	}
	if !types.AssignableTo(innerTV.Type, targetIface) {
		return nil, false
	}
	if types.Identical(innerTV.Type, targetIface) {
		return inner, true
	}
	if isInterfaceType(innerTV.Type) || isProvablyNonNilExpr(inner, facts, info) {
		if exprContextPreservesInterfaceType(call, targetIface, innerTV.Type, parent, parents, pkg, info) {
			return inner, true
		}
	}
	return nil, false
}

func (p *J2gortPeepholePass) tryRewriteStringEquals(
	call *ast.CallExpr,
	parent ast.Node,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	if _, isExprStmt := parent.(*ast.ExprStmt); isExprStmt {
		return false
	}

	lhs := stripParens(call.Args[0])
	if !hasExactOrDefaultStringType(lhs, info) || !isGuaranteedNonNullString(lhs, info) {
		return false
	}

	rawB := stripParens(call.Args[1])
	var rhs ast.Expr
	if bCall, ok := rawB.(*ast.CallExpr); ok && len(bCall.Args) == 1 && !bCall.Ellipsis.IsValid() {
		if bFn, _, _, ok := isJ2gortCall(bCall, info); ok && bFn == "BoxString" {
			cand := stripParens(bCall.Args[0])
			if hasExactOrDefaultStringType(cand, info) {
				rhs = cand
			}
		}
	}
	if rhs == nil && hasExactOrDefaultStringType(rawB, info) {
		rhs = rawB
	}
	if rhs == nil {
		return false
	}

	targetNode := ast.Expr(call)
	targetParent := parent
	op := token.EQL

	cur := ast.Expr(call)
	curParent := parent
	for {
		pe, isParen := curParent.(*ast.ParenExpr)
		if !isParen || pe.X != cur {
			break
		}
		cur = pe
		curParent = parents[pe]
	}
	if unary, isNot := curParent.(*ast.UnaryExpr); isNot && unary.Op == token.NOT && unary.X == cur {
		if outer := parents[unary]; outer != nil {
			if _, isES := outer.(*ast.ExprStmt); !isES {
				targetNode = unary
				targetParent = outer
				op = token.NEQ
			}
		}
	}

	var cmp ast.Expr = &ast.BinaryExpr{
		X:     wrapBinaryOperandIfNeeded(lhs, op),
		OpPos: call.Pos(),
		Op:    op,
		Y:     wrapBinaryOperandIfNeeded(rhs, op),
	}
	info.Types[cmp] = types.TypeAndValue{Type: types.Typ[types.Bool]}
	if binaryNeedsParens(targetParent, op) {
		pe := &ast.ParenExpr{X: cmp}
		info.Types[pe] = types.TypeAndValue{Type: types.Typ[types.Bool]}
		cmp = pe
	}
	return replaceExprInParent(targetParent, targetNode, cmp, parents)
}

type strConcatItem struct {
	isConst  bool
	constVal string
	origExpr ast.Expr
	pos      token.Pos
}

func (p *J2gortPeepholePass) tryRewriteStrConcat(
	call *ast.CallExpr,
	strConcatFn *types.Func,
	parent ast.Node,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	if _, isExprStmt := parent.(*ast.ExprStmt); isExprStmt {
		return false
	}
	funSel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkgIdent, ok := stripParens(funSel.X).(*ast.Ident)
	if !ok {
		return false
	}

	pos := call.Pos()
	if len(call.Args) == 0 {
		lit := makeTypedStringBasicLit("", pos, info)
		return replaceExprInParent(parent, call, lit, parents)
	}

	items := make([]strConcatItem, 0, len(call.Args))
	foldedChanged := false
	for _, a := range call.Args {
		if cv, ok := constStrValForConcat(a, info); ok {
			if !isPlainNonNullStringBasicLit(a, cv) {
				foldedChanged = true
			}
			if len(items) > 0 && items[len(items)-1].isConst {
				items[len(items)-1].constVal += cv
				foldedChanged = true
			} else {
				items = append(items, strConcatItem{isConst: true, constVal: cv, pos: a.Pos()})
			}
			continue
		}
		items = append(items, strConcatItem{origExpr: a, pos: a.Pos()})
	}

	var strValFn *types.Func
	if strConcatFn != nil && strConcatFn.Pkg() != nil && strConcatFn.Pkg().Scope() != nil {
		if obj, ok := strConcatFn.Pkg().Scope().Lookup("StrVal").(*types.Func); ok {
			strValFn = obj
		}
	}

	makeStrValExpr := func(arg ast.Expr, argPos token.Pos) *ast.CallExpr {
		selIdent := &ast.Ident{NamePos: argPos, Name: "StrVal"}
		if strValFn != nil {
			info.Uses[selIdent] = strValFn
		}
		pkgRef := &ast.Ident{NamePos: argPos, Name: pkgIdent.Name}
		if u := info.Uses[pkgIdent]; u != nil {
			info.Uses[pkgRef] = u
		}
		svSel := &ast.SelectorExpr{X: pkgRef, Sel: selIdent}
		svCall := &ast.CallExpr{Fun: svSel, Args: []ast.Expr{stripParens(arg)}}
		info.Types[svCall] = types.TypeAndValue{Type: types.Typ[types.String]}
		return svCall
	}

	allPureString := true
	for _, it := range items {
		if !it.isConst && !hasExactOrDefaultStringType(it.origExpr, info) {
			allPureString = false
		}
	}

	if len(items) == 1 && items[0].isConst {
		return replaceExprInParent(parent, call, makeTypedStringBasicLit(items[0].constVal, pos, info), parents)
	}

	if allPureString && len(items) >= 1 && len(items) <= 4 {
		inAdd := len(items) > 1
		buildOperand := func(it strConcatItem) ast.Expr {
			if it.isConst {
				return makeTypedStringBasicLit(it.constVal, it.pos, info)
			}
			inner := stripParens(it.origExpr)
			if isGuaranteedNonNullString(inner, info) {
				if inAdd {
					return wrapBinaryOperandIfNeeded(inner, token.ADD)
				}
				return inner
			}
			return makeStrValExpr(inner, it.pos)
		}

		repl := buildOperand(items[0])
		for i := 1; i < len(items); i++ {
			bin := &ast.BinaryExpr{X: repl, OpPos: items[i].pos, Op: token.ADD, Y: buildOperand(items[i])}
			info.Types[bin] = types.TypeAndValue{Type: types.Typ[types.String]}
			repl = bin
		}
		if needsParensInContext(repl, call, parent, parents) {
			pe := &ast.ParenExpr{X: repl}
			info.Types[pe] = info.Types[repl]
			repl = pe
		}
		return replaceExprInParent(parent, call, repl, parents)
	}

	if !foldedChanged {
		return false
	}
	newArgs := make([]ast.Expr, len(items))
	for i, it := range items {
		if it.isConst {
			newArgs[i] = makeTypedStringBasicLit(it.constVal, it.pos, info)
		} else {
			newArgs[i] = it.origExpr
		}
	}
	newCall := &ast.CallExpr{
		Fun:    call.Fun,
		Lparen: call.Lparen,
		Args:   newArgs,
		Rparen: call.Rparen,
	}
	info.Types[newCall] = info.Types[call]
	return replaceExprInParent(parent, call, newCall, parents)
}

func makeTypedStringBasicLit(val string, pos token.Pos, info *types.Info) *ast.BasicLit {
	lit := &ast.BasicLit{
		ValuePos: pos,
		Kind:     token.STRING,
		Value:    strconv.Quote(val),
	}
	if info != nil {
		info.Types[lit] = types.TypeAndValue{
			Type:  types.Typ[types.String],
			Value: constant.MakeString(val),
		}
	}
	return lit
}

func isPlainNonNullStringBasicLit(expr ast.Expr, val string) bool {
	if val == j2gortNullString {
		return false
	}
	lit, ok := expr.(*ast.BasicLit)
	return ok && lit.Kind == token.STRING
}

func isPureConstExpr(expr ast.Expr) bool {
	switch e := stripParens(expr).(type) {
	case *ast.BasicLit, *ast.Ident:
		return true
	case *ast.SelectorExpr:
		return isPureConstExpr(e.X)
	case *ast.UnaryExpr:
		return e.Op != token.ARROW && isPureConstExpr(e.X)
	case *ast.BinaryExpr:
		return isPureConstExpr(e.X) && isPureConstExpr(e.Y)
	case *ast.CallExpr:
		return len(e.Args) == 1 && !e.Ellipsis.IsValid() && isPureConstExpr(e.Args[0])
	}
	return false
}

func constStrValForConcat(expr ast.Expr, info *types.Info) (string, bool) {
	if info == nil || expr == nil || !isPureConstExpr(expr) {
		return "", false
	}
	tv, ok := info.Types[expr]
	if !ok || tv.IsType() || tv.IsBuiltin() || tv.IsNil() || tv.Type == nil || tv.Value == nil {
		return "", false
	}
	if hasExactOrDefaultStringType(expr, info) && tv.Value.Kind() == constant.String {
		s := constant.StringVal(tv.Value)
		if s == j2gortNullString {
			return "null", true
		}
		return s, true
	}
	b, ok := types.Unalias(tv.Type).(*types.Basic)
	if !ok {
		return "", false
	}
	switch b.Kind() {
	case types.Bool, types.UntypedBool:
		if tv.Value.Kind() == constant.Bool {
			return strconv.FormatBool(constant.BoolVal(tv.Value)), true
		}
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.UntypedInt:
		if tv.Value.Kind() == constant.Int {
			if v, exact := constant.Int64Val(tv.Value); exact {
				return strconv.FormatInt(v, 10), true
			}
		}
	case types.Uint16:
		if tv.Value.Kind() == constant.Int {
			if v, exact := constant.Int64Val(tv.Value); exact && v >= 0 && v < 0x80 {
				return string(byte(v)), true
			}
		}
	}
	return "", false
}

func wrapBinaryOperandIfNeeded(e ast.Expr, outerOp token.Token) ast.Expr {
	if be, ok := e.(*ast.BinaryExpr); ok && be.Op.Precedence() <= outerOp.Precedence() {
		return &ast.ParenExpr{X: e}
	}
	return e
}

var unboxWantKind = map[string]types.BasicKind{
	"UnboxBoolean": types.Bool,
	"UnboxInt":     types.Int32,
	"UnboxLong":    types.Int64,
	"UnboxDouble":  types.Float64,
	"UnboxFloat":   types.Float32,
	"UnboxChar":    types.Uint16,
	"UnboxByte":    types.Int8,
	"UnboxShort":   types.Int16,
}

func unwrapBoxRoundTrip(unboxName string, inner ast.Expr, info *types.Info) (ast.Expr, bool) {
	wantKind, ok := unboxWantKind[unboxName]
	if !ok {
		return nil, false
	}
	innerCall, ok := inner.(*ast.CallExpr)
	if !ok || len(innerCall.Args) != 1 || innerCall.Ellipsis.IsValid() {
		return nil, false
	}
	fnName, _, _, ok := isJ2gortCall(innerCall, info)
	if !ok || fnName != "ValueOf"+strings.TrimPrefix(unboxName, "Unbox") {
		return nil, false
	}
	x := stripParens(innerCall.Args[0])
	if hasExactIntrinsicType(x, types.Typ[wantKind], info) {
		return x, true
	}
	if wantKind == types.Bool {
		if tv, ok := info.Types[x]; ok && tv.Type != nil {
			if b, ok := types.Unalias(tv.Type).(*types.Basic); ok && (b.Kind() == types.Bool || b.Kind() == types.UntypedBool) {
				return x, true
			}
		}
	}
	return nil, false
}

func isJ2gortPkg(pkg *types.Package) bool {
	if pkg == nil {
		return false
	}
	p := pkg.Path()
	return pkg.Name() == "j2gort" || p == "j2gort" || strings.HasSuffix(p, "/j2gort")
}

func isJ2gortCall(call *ast.CallExpr, info *types.Info) (string, *types.Func, ast.Expr, bool) {
	if call == nil || info == nil {
		return "", nil, nil, false
	}
	fn, explicitTypeArg := calleeFunc(call, info)
	if fn == nil || !isJ2gortPkg(fn.Pkg()) {
		return "", nil, nil, false
	}
	if sig, ok := fn.Type().(*types.Signature); !ok || sig.Recv() != nil {
		return "", nil, nil, false
	}
	return fn.Name(), fn, explicitTypeArg, true
}

func hasExactOrDefaultStringType(expr ast.Expr, info *types.Info) bool {
	if info == nil {
		return false
	}
	tv, ok := info.Types[expr]
	if !ok || tv.IsType() || tv.IsBuiltin() || tv.IsNil() || tv.Type == nil {
		return false
	}
	b, ok := types.Unalias(tv.Type).(*types.Basic)
	return ok && (b.Kind() == types.String || b.Kind() == types.UntypedString)
}

func isGuaranteedNonNullString(expr ast.Expr, info *types.Info) bool {
	if info == nil || expr == nil {
		return false
	}
	e := stripParens(expr)
	tv, ok := info.Types[e]
	if !ok || tv.Type == nil || !isStringType(tv.Type) {
		return false
	}

	if tv.Value != nil && tv.Value.Kind() == constant.String {
		return constant.StringVal(tv.Value) != j2gortNullString
	}

	if bin, ok := e.(*ast.BinaryExpr); ok && bin.Op == token.ADD {
		if isGuaranteedNonNullString(bin.X, info) && isGuaranteedNonNullString(bin.Y, info) {
			if tvX, ok := info.Types[bin.X]; ok && tvX.Value != nil && tvX.Value.Kind() == constant.String && constant.StringVal(tvX.Value) != "" {
				return true
			}
			if tvY, ok := info.Types[bin.Y]; ok && tvY.Value != nil && tvY.Value.Kind() == constant.String && constant.StringVal(tvY.Value) != "" {
				return true
			}
		}
	}

	call, ok := e.(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return false
	}

	if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
		switch fnName {
		case "StrVal", "StrConcat", "StringFormat", "CheckNotNullString", "BoolToStr",
			"StrValString", "StrValInt32", "StrValInt64", "StrValBool", "StrValChar":
			return true
		case "CheckNotNull":
			return isStringType(tv.Type)
		}
		return false
	}

	fn, _ := calleeFunc(call, info)
	if fn == nil || fn.Pkg() == nil {
		return false
	}
	pkgPath := fn.Pkg().Path()
	name := fn.Name()

	switch {
	case pkgPath == "strconv":
		switch name {
		case "Itoa", "FormatInt", "FormatUint", "FormatBool", "FormatFloat", "Quote":
			return true
		}
	case pkgPath == "fmt" && name == "Sprintf":
		return true
	case strings.HasSuffix(pkgPath, "/javalang"):
		switch name {
		case "StringFromChar", "StringFromUTF16", "ToJavaString",
			"IntegerToString", "IntegerToHexString", "IntegerToBinaryString", "IntegerToOctalString",
			"LongToString", "LongToHexString", "DoubleToString", "FloatToString", "BooleanToString":
			return true
		}
	}

	if sig, ok := fn.Type().(*types.Signature); ok && sig.Recv() != nil && name == "String" && len(call.Args) == 0 {
		recvType := types.Unalias(sig.Recv().Type())
		if ptr, ok := recvType.(*types.Pointer); ok {
			recvType = types.Unalias(ptr.Elem())
		}
		if named, ok := recvType.(*types.Named); ok && named.Obj() != nil {
			if pkgPath == "strings" && named.Obj().Name() == "Builder" {
				return true
			}
			if isJ2gortPkg(named.Obj().Pkg()) && named.Obj().Name() == "Ext_java_lang_StringBuilder" {
				return true
			}
		}
	}

	return false
}

func isCheckNotNullFuncName(name string) bool {
	switch name {
	case "CheckNotNull", "CheckNotNullPtr", "CheckNotNullString", "CheckNotNullSlice", "CheckNotNullNode", "CheckNotNullIface", "CheckNotNullGeneric":
		return true
	}
	return false
}

func isNamedUint32NodeType(t types.Type) bool {
	if t == nil {
		return false
	}
	named, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return false
	}
	b, ok := named.Underlying().(*types.Basic)
	return ok && b.Kind() == types.Uint32
}

func isNamedRhinoNodeType(t types.Type) bool {
	if t == nil {
		return false
	}
	named, ok := types.Unalias(t).(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Name() != "Node" {
		return false
	}
	b, ok := named.Underlying().(*types.Basic)
	return ok && b.Kind() == types.Uint32
}

func lookupJ2gortSiblingFunc(fnObj *types.Func, name string) *types.Func {
	if fnObj == nil || fnObj.Pkg() == nil || fnObj.Pkg().Scope() == nil {
		return nil
	}
	f, _ := fnObj.Pkg().Scope().Lookup(name).(*types.Func)
	return f
}

func extractExplicitSingleTypeArg(fun ast.Expr) ast.Expr {
	switch ix := stripParens(fun).(type) {
	case *ast.IndexExpr:
		return ix.Index
	case *ast.IndexListExpr:
		if len(ix.Indices) == 1 {
			return ix.Indices[0]
		}
	}
	return nil
}

func buildSpecializedJ2gortCall(
	call *ast.CallExpr,
	targetName string,
	targetObj *types.Func,
	explicitTypeArg ast.Expr,
	arg ast.Expr,
	retType types.Type,
	info *types.Info,
) (*ast.CallExpr, bool) {
	baseFun := stripParens(call.Fun)
	switch ix := baseFun.(type) {
	case *ast.IndexExpr:
		baseFun = stripParens(ix.X)
	case *ast.IndexListExpr:
		baseFun = stripParens(ix.X)
	}

	var newFun ast.Expr
	var newSelIdent *ast.Ident
	switch fn := baseFun.(type) {
	case *ast.SelectorExpr:
		pkgIdent, ok := stripParens(fn.X).(*ast.Ident)
		if !ok {
			return nil, false
		}
		pkgCopy := &ast.Ident{NamePos: pkgIdent.NamePos, Name: pkgIdent.Name}
		if info != nil {
			if u, hasU := info.Uses[pkgIdent]; hasU {
				info.Uses[pkgCopy] = u
			}
		}
		newSelIdent = &ast.Ident{NamePos: fn.Sel.NamePos, Name: targetName}
		newFun = &ast.SelectorExpr{X: pkgCopy, Sel: newSelIdent}
	case *ast.Ident:
		newSelIdent = &ast.Ident{NamePos: fn.NamePos, Name: targetName}
		newFun = newSelIdent
	default:
		return nil, false
	}

	if explicitTypeArg != nil && (targetObj == nil || !isNonGenericFunc(targetObj)) {
		newFun = &ast.IndexExpr{
			X:      newFun,
			Lbrack: explicitTypeArg.Pos(),
			Index:  explicitTypeArg,
			Rbrack: explicitTypeArg.End(),
		}
	}

	newCall := &ast.CallExpr{
		Fun:    newFun,
		Lparen: call.Lparen,
		Args:   []ast.Expr{arg},
		Rparen: call.Rparen,
	}
	if info != nil {
		if targetObj != nil && newSelIdent != nil {
			info.Uses[newSelIdent] = targetObj
		}
		info.Types[newCall] = types.TypeAndValue{Type: retType}
	}
	return newCall, true
}

func specializeCheckNotNullCall(
	call *ast.CallExpr,
	inner ast.Expr,
	callType types.Type,
	fnObj *types.Func,
	isExprStmt bool,
	exports peepholeExports,
	info *types.Info,
) (*ast.CallExpr, bool) {
	if call == nil || inner == nil || callType == nil || info == nil {
		return nil, false
	}
	innerTV, ok := info.Types[inner]
	if !ok || innerTV.IsType() || innerTV.IsBuiltin() {
		return nil, false
	}

	hasHelper := func(name string, exportFlag bool) (*types.Func, bool) {
		targetObj := lookupJ2gortSiblingFunc(fnObj, name)
		return targetObj, exportFlag || targetObj != nil
	}
	retMatches := func(retT types.Type) bool {
		return isExprStmt || types.Identical(retT, callType)
	}

	// Untyped nil with explicit pointer or slice type argument: CheckNotNull[*U](nil) / CheckNotNull[[]U](nil).
	if innerTV.IsNil() {
		tArg := extractExplicitSingleTypeArg(call.Fun)
		if tArg == nil {
			return nil, false
		}
		switch ut := types.Unalias(callType).(type) {
		case *types.Pointer:
			if st, okS := stripParens(tArg).(*ast.StarExpr); okS {
				if targetObj, avail := hasHelper("CheckNotNullPtr", exports.hasCheckNotNullPtr); avail && (targetObj == nil || !isNonGenericFunc(targetObj)) {
					return buildSpecializedJ2gortCall(call, "CheckNotNullPtr", targetObj, st.X, inner, ut, info)
				}
			}
		case *types.Slice:
			if at, okA := stripParens(tArg).(*ast.ArrayType); okA && at.Len == nil {
				if targetObj, avail := hasHelper("CheckNotNullSlice", exports.hasCheckNotNullSlice); avail && (targetObj == nil || !isNonGenericFunc(targetObj)) {
					return buildSpecializedJ2gortCall(call, "CheckNotNullSlice", targetObj, at.Elt, inner, ut, info)
				}
			}
		}
		return nil, false
	}

	if innerTV.Type == nil {
		return nil, false
	}

	// 1. string -> CheckNotNullString(x)
	if hasExactOrDefaultStringType(inner, info) && retMatches(types.Typ[types.String]) {
		if targetObj, avail := hasHelper("CheckNotNullString", exports.hasCheckNotNullString); avail {
			return buildSpecializedJ2gortCall(call, "CheckNotNullString", targetObj, nil, inner, types.Typ[types.String], info)
		}
	}

	// 2. *U (unnamed/alias pointer type, not a defined named pointer type) -> CheckNotNullPtr(x)
	if _, isPtr := types.Unalias(innerTV.Type).(*types.Pointer); isPtr && !isInvalidOrUntyped(innerTV.Type) && retMatches(innerTV.Type) {
		if targetObj, avail := hasHelper("CheckNotNullPtr", exports.hasCheckNotNullPtr); avail {
			return buildSpecializedJ2gortCall(call, "CheckNotNullPtr", targetObj, nil, inner, innerTV.Type, info)
		}
	}

	// 3. []U (unnamed/alias slice type, not a defined named slice type) -> CheckNotNullSlice(x)
	if _, isSlice := types.Unalias(innerTV.Type).(*types.Slice); isSlice && !isInvalidOrUntyped(innerTV.Type) && retMatches(innerTV.Type) {
		if targetObj, avail := hasHelper("CheckNotNullSlice", exports.hasCheckNotNullSlice); avail {
			return buildSpecializedJ2gortCall(call, "CheckNotNullSlice", targetObj, nil, inner, innerTV.Type, info)
		}
	}

	// 4. rhino.Node (named uint32 type Node where 0 is null) -> CheckNotNullNode(x)
	if isNamedRhinoNodeType(innerTV.Type) && retMatches(innerTV.Type) {
		if targetObj, avail := hasHelper("CheckNotNullNode", exports.hasCheckNotNullNode); avail {
			if hasExactIntrinsicType(inner, innerTV.Type, info) {
				return buildSpecializedJ2gortCall(call, "CheckNotNullNode", targetObj, nil, inner, innerTV.Type, info)
			}
			if tArg := extractExplicitSingleTypeArg(call.Fun); tArg != nil {
				return buildSpecializedJ2gortCall(call, "CheckNotNullNode", targetObj, tArg, inner, innerTV.Type, info)
			}
		}
	}

	return nil, false
}

func isToIfaceCall(call *ast.CallExpr, info *types.Info) (types.Type, bool) {
	if call == nil || info == nil || len(call.Args) != 1 || call.Ellipsis.IsValid() {
		return nil, false
	}
	fn, _ := calleeFunc(call, info)
	if fn == nil {
		return nil, false
	}
	name := fn.Name()
	if !strings.HasPrefix(name, "ToIface_") || len(name) <= len("ToIface_") {
		return nil, false
	}
	origSig, ok := fn.Origin().Type().(*types.Signature)
	if !ok || origSig.Recv() != nil || origSig.Variadic() || origSig.Params().Len() != 1 || origSig.Results().Len() != 1 {
		return nil, false
	}
	if origSig.TypeParams() != nil && origSig.TypeParams().Len() > 1 {
		return nil, false
	}
	callType := info.Types[call].Type
	if !isInterfaceType(callType) {
		return nil, false
	}
	return callType, true
}

type peepholeFileFacts struct {
	nonNilVars      map[*types.Var]bool
	nonZeroNodeVars map[*types.Var]bool
}

func collectPeepholeFileFacts(file *ast.File, info *types.Info, parents map[ast.Node]ast.Node) *peepholeFileFacts {
	facts := &peepholeFileFacts{
		nonNilVars:      make(map[*types.Var]bool),
		nonZeroNodeVars: make(map[*types.Var]bool),
	}
	if file == nil || info == nil {
		return facts
	}

	defCounts := make(map[*types.Var]int)
	mutatedVars := make(map[*types.Var]bool)
	markMutatedIdent := func(e ast.Expr) {
		if e == nil {
			return
		}
		if id, ok := stripParens(e).(*ast.Ident); ok && id.Name != "_" {
			if v, ok := info.Uses[id].(*types.Var); ok && v != nil {
				mutatedVars[v] = true
			}
		}
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.AssignStmt:
			for _, lhs := range s.Lhs {
				if s.Tok == token.DEFINE {
					if id, ok := stripParens(lhs).(*ast.Ident); ok && id.Name != "_" {
						if v, ok := info.Defs[id].(*types.Var); ok && v != nil {
							defCounts[v]++
						}
					}
				}
				markMutatedIdent(lhs)
			}
		case *ast.ValueSpec:
			for _, id := range s.Names {
				if id.Name != "_" {
					if v, ok := info.Defs[id].(*types.Var); ok && v != nil {
						defCounts[v]++
					}
				}
			}
		case *ast.IncDecStmt:
			markMutatedIdent(s.X)
		case *ast.UnaryExpr:
			if s.Op == token.AND {
				markMutatedIdent(s.X)
			}
		case *ast.RangeStmt:
			if s.Tok == token.ASSIGN {
				markMutatedIdent(s.Key)
				markMutatedIdent(s.Value)
			}
		}
		return true
	})

	recordLocalVarInit := func(vObj *types.Var, rhs ast.Expr) {
		if vObj == nil || defCounts[vObj] != 1 || mutatedVars[vObj] || exprUsesVar(rhs, vObj, info) {
			return
		}
		vt := types.Unalias(vObj.Type())
		_, isPtr := vt.(*types.Pointer)
		_, isSlice := vt.(*types.Slice)
		if isPtr || isSlice || isInterfaceType(vt) {
			if isGuaranteedNonNilForCheckNotNull(rhs, facts, info) {
				facts.nonNilVars[vObj] = true
			}
		} else if isNamedUint32NodeType(vt) && isGuaranteedNonZeroNode(rhs, facts, info) {
			facts.nonZeroNodeVars[vObj] = true
		}
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			if fn.Recv != nil && len(fn.Recv.List) == 1 && len(fn.Recv.List[0].Names) == 1 {
				recvIdent := fn.Recv.List[0].Names[0]
				if recvIdent.Name == "this_" {
					if vObj, ok := info.Defs[recvIdent].(*types.Var); ok && vObj != nil && !mutatedVars[vObj] {
						if _, isPtr := types.Unalias(vObj.Type()).(*types.Pointer); isPtr {
							facts.nonNilVars[vObj] = true
						}
					}
				}
			}
		case *ast.AssignStmt:
			if fn.Tok == token.DEFINE && len(fn.Lhs) == 1 && len(fn.Rhs) == 1 && enclosingFuncBody(fn, parents) != nil {
				if lhsIdent, ok := fn.Lhs[0].(*ast.Ident); ok && lhsIdent.Name != "_" {
					if vObj, ok := info.Defs[lhsIdent].(*types.Var); ok {
						recordLocalVarInit(vObj, fn.Rhs[0])
					}
				}
			}
		case *ast.DeclStmt:
			if enclosingFuncBody(fn, parents) != nil {
				if gd, ok := fn.Decl.(*ast.GenDecl); ok && gd.Tok == token.VAR && len(gd.Specs) == 1 {
					if spec, ok := gd.Specs[0].(*ast.ValueSpec); ok && len(spec.Names) == 1 && len(spec.Values) == 1 && spec.Names[0].Name != "_" {
						if vObj, ok := info.Defs[spec.Names[0]].(*types.Var); ok {
							recordLocalVarInit(vObj, spec.Values[0])
						}
					}
				}
			}
		}
		return true
	})

	return facts
}

func isGuaranteedNonZeroNode(expr ast.Expr, facts *peepholeFileFacts, info *types.Info) bool {
	if info == nil || expr == nil {
		return false
	}
	e := stripParens(expr)
	tv, ok := info.Types[e]
	if !ok || tv.Type == nil || !isNamedUint32NodeType(tv.Type) {
		return false
	}
	if tv.Value != nil && tv.Value.Kind() == constant.Int {
		return constant.Sign(tv.Value) != 0
	}
	if id, ok := e.(*ast.Ident); ok && facts != nil {
		if vObj, ok := info.Uses[id].(*types.Var); ok && facts.nonZeroNodeVars[vObj] {
			return true
		}
	}
	if call, ok := e.(*ast.CallExpr); ok && len(call.Args) == 1 && !call.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			return fnName == "CheckNotNullNode" || fnName == "CheckNotNull"
		}
	}
	return false
}

func isProvablyNonNilExpr(expr ast.Expr, facts *peepholeFileFacts, info *types.Info) bool {
	return isGuaranteedNonNilForCheckNotNull(expr, facts, info)
}

func isGuaranteedNeverTypedNil(expr ast.Expr, facts *peepholeFileFacts, info *types.Info) bool {
	if isGuaranteedNonNilForCheckNotNull(expr, facts, info) {
		return true
	}
	if info == nil || expr == nil {
		return false
	}
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return false
	}
	if _, isTI := isToIfaceCall(call, info); isTI {
		return true
	}
	if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
		switch fnName {
		case "AsIface", "Cast", "UncheckedCast":
			return isInterfaceType(info.Types[call].Type)
		}
	}
	return false
}

func isGuaranteedNonNilForCheckNotNull(expr ast.Expr, facts *peepholeFileFacts, info *types.Info) bool {
	if info == nil || expr == nil {
		return false
	}
	e := stripParens(expr)
	if id, ok := e.(*ast.Ident); ok && facts != nil {
		if vObj, ok := info.Uses[id].(*types.Var); ok && facts.nonNilVars[vObj] {
			return true
		}
	}
	if u, ok := e.(*ast.UnaryExpr); ok && u.Op == token.AND {
		if tv := info.Types[e]; tv.Type != nil {
			_, isPtr := types.Unalias(tv.Type).(*types.Pointer)
			return isPtr
		}
	}
	if _, isComp := e.(*ast.CompositeLit); isComp {
		if tv := info.Types[e]; tv.Type != nil {
			switch types.Unalias(tv.Type).Underlying().(type) {
			case *types.Slice, *types.Map:
				return true
			}
		}
	}

	if st, ok := e.(*ast.StarExpr); ok {
		if innerCall, isCall := stripParens(st.X).(*ast.CallExpr); isCall {
			if fnName, _, _, isJ2 := isJ2gortCall(innerCall, info); isJ2 {
				switch fnName {
				case "SliceBuilderAdd", "SliceBuilderAddAll", "SliceBuilderAddAllNode", "SliceBuilderAddAllSlice":
					return true
				}
			}
		}
	}

	call, ok := e.(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return false
	}
	if len(call.Args) == 0 {
		if fl, isLit := stripParens(call.Fun).(*ast.FuncLit); isLit && fl.Body != nil {
			allNonNil := true
			hasRet := false
			ast.Inspect(fl.Body, func(n ast.Node) bool {
				if !allNonNil || n == nil {
					return false
				}
				if nested, okN := n.(*ast.FuncLit); okN && nested != fl {
					return false
				}
				if ret, okR := n.(*ast.ReturnStmt); okR {
					hasRet = true
					if len(ret.Results) != 1 || !isGuaranteedNonNilForCheckNotNull(ret.Results[0], facts, info) {
						allNonNil = false
					}
					return false
				}
				return true
			})
			if hasRet && allNonNil {
				return true
			}
		}
	}
	if id, ok := stripParens(call.Fun).(*ast.Ident); ok {
		if b, ok := info.Uses[id].(*types.Builtin); ok {
			if (b.Name() == "new" && len(call.Args) == 1) || (b.Name() == "make" && len(call.Args) >= 1) {
				return true
			}
		}
	}
	if targetIface, isTI := isToIfaceCall(call, info); isTI {
		inner := stripParens(call.Args[0])
		innerT := info.Types[inner].Type
		return innerT != nil && types.AssignableTo(innerT, targetIface) && isGuaranteedNonNilForCheckNotNull(inner, facts, info)
	}
	if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
		if strings.HasPrefix(fnName, "NewExt_") || strings.HasPrefix(fnName, "ValueOf") || isLambdaFactoryFuncName(fnName) {
			if tv := info.Types[call]; tv.Type != nil {
				_, isBasic := types.Unalias(tv.Type).(*types.Basic)
				return !isBasic
			}
		}
		switch fnName {
		case "ConstructShimThrowable", "CheckNotNullPtr", "CheckNotNullSlice", "CheckNotNullIface",
			"MakeNonNilSlice", "ImmListCopyToSlice", "ImmListCopyNodeToSlice", "ImmListCopyArrayToSlice",
			"ImmListOfVarargsToSlice", "ImmSliceCloneNonNull", "ImmSliceSubList", "ImmSliceReverse",
			"CloneNonNilSlice", "CopyStrColToSlice", "CopyNodeColToSlice", "CopyCollectionToSlice", "CopyIteratorToSlice",
			"SliceBuilderAdd", "SliceBuilderAddAll", "SliceBuilderAddAllNode", "SliceBuilderAddAllSlice",
			"SliceUnmodifiableList", "SliceUnmodifiableListLive", "SliceReadOnlyIterable",
			"NewOrderedMapPtr", "NewOrderedMapPtrCap", "AcquireOrderedMapPtr", "AcquireOrderedMapPtrCap", "OrderedMapPtrFrom", "OrderedMapPtrFromJavaMap", "OrderedMapPtrFromNodeKeyJavaMap",
			"NewCustomOrderedMapPtr", "NewCustomOrderedMapPtrCap", "AcquireCustomOrderedMapPtr", "AcquireCustomOrderedMapPtrCap", "CustomOrderedMapPtrFrom", "CustomOrderedMapPtrFromJavaMap",
			"NewIdentityOrderedMapPtr", "NewIdentityOrderedMapPtrCap", "AcquireIdentityOrderedMapPtr", "AcquireIdentityOrderedMapPtrCap", "IdentityOrderedMapPtrFrom", "IdentityOrderedMapPtrFromJavaMap",
			"NewOrderedSetPtr", "NewOrderedSetPtrCap", "OrderedSetPtrFromSet", "OrderedSetPtrFromStrCol", "OrderedSetPtrFromNodeCol", "OrderedSetPtrFromCol",
			"NewCustomOrderedSetPtr", "NewCustomOrderedSetPtrCap", "CustomOrderedSetPtrFromSet", "CustomOrderedSetPtrFromCol",
			"NewIdentityOrderedSetPtr", "NewIdentityOrderedSetPtrCap", "IdentityOrderedSetPtrFromSet",
			"NewDequePtr", "NewDequePtrCap", "DequePtrFromDeque", "DequePtrFromSlice", "DequePtrFromOrderedSet", "DequePtrFromOrderedSetPtr", "DequePtrFromLocalOrderedSet",
			"DequePtrFromCustomOrderedSet", "DequePtrFromCustomOrderedSetPtr", "DequePtrFromIdentityOrderedSet", "DequePtrFromIdentityOrderedSetPtr", "DequePtrFromCol":
			return true
		case "CheckNotNull":
			return facts != nil
		case "AsIface", "Cast", "UncheckedCast", "CastPtr":
			if len(call.Args) == 1 {
				inner := stripParens(call.Args[0])
				callT, innerT := info.Types[call].Type, info.Types[inner].Type
				return callT != nil && innerT != nil && types.AssignableTo(innerT, callT) && isGuaranteedNonNilForCheckNotNull(inner, facts, info)
			}
		}
		return false
	}

	if fn, _ := calleeFunc(call, info); fn != nil && fn.Pkg() != nil {
		if sig, ok := fn.Type().(*types.Signature); ok && sig.Recv() == nil {
			name := fn.Name()
			if isLambdaFactoryFuncName(name) {
				return true
			}
			if strings.HasPrefix(name, "New") && len(name) > 3 && name[3] >= 'A' && name[3] <= 'Z' {
				_, isPtr := types.Unalias(info.Types[call].Type).(*types.Pointer)
				return isPtr
			}
		}
	}
	return false
}

func isGuaranteedNonNilAlloc(expr ast.Expr, info *types.Info) bool {
	return isGuaranteedNonNilForCheckNotNull(expr, nil, info)
}

func isGuaranteedSafeForToAny(expr ast.Expr, info *types.Info) bool {
	if info == nil || expr == nil {
		return false
	}
	e := stripParens(expr)
	tv, ok := info.Types[e]
	if !ok || tv.IsType() || tv.IsBuiltin() || tv.IsNil() || tv.Type == nil {
		return false
	}

	if b, ok := types.Unalias(tv.Type).(*types.Basic); ok {
		switch b.Kind() {
		case types.Bool, types.UntypedBool:
			return true
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr,
			types.Float32, types.Float64:
			return tv.Value == nil || hasExactIntrinsicType(e, tv.Type, info)
		case types.String, types.UntypedString:
			return isGuaranteedNonNullString(e, info)
		default:
			return false
		}
	}

	if named, ok := types.Unalias(tv.Type).(*types.Named); ok {
		if b, isBasic := named.Underlying().(*types.Basic); isBasic && b.Kind() != types.String && (b.Info()&types.IsUntyped) == 0 {
			return tv.Value == nil || hasExactIntrinsicType(e, tv.Type, info)
		}
	}

	if isGuaranteedNonNilAlloc(e, info) {
		return true
	}
	if call, ok := e.(*ast.CallExpr); ok && !call.Ellipsis.IsValid() {
		if _, isTI := isToIfaceCall(call, info); isTI {
			return true
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "BoxString", "BoxNode", "ToAny", "PtrToAny", "SliceToAny", "AsIface":
				return true
			}
		}
	}
	return false
}

func ifaceUnderlying(t types.Type) *types.Interface {
	if t == nil {
		return nil
	}
	ut := types.Unalias(t)
	if _, isTP := ut.(*types.TypeParam); isTP {
		return nil
	}
	iface, _ := ut.Underlying().(*types.Interface)
	return iface
}

func isInterfaceType(t types.Type) bool { return ifaceUnderlying(t) != nil }
func isNonEmptyInterfaceType(t types.Type) bool {
	u := ifaceUnderlying(t)
	return u != nil && u.NumMethods() > 0
}
func isEmptyInterfaceType(t types.Type) bool {
	u := ifaceUnderlying(t)
	return u != nil && u.NumMethods() == 0
}

func isNonGenericFunc(fn *types.Func) bool {
	if fn == nil {
		return false
	}
	sig, ok := fn.Origin().Type().(*types.Signature)
	return ok && (sig.TypeParams() == nil || sig.TypeParams().Len() == 0) &&
		(sig.RecvTypeParams() == nil || sig.RecvTypeParams().Len() == 0)
}

func hasExplicitTypeArgs(fun ast.Expr) bool {
	switch stripParens(fun).(type) {
	case *ast.IndexExpr, *ast.IndexListExpr:
		return true
	}
	return false
}

func typeContainsTypeParam(t types.Type) bool {
	if t == nil {
		return false
	}
	switch ut := types.Unalias(t).(type) {
	case *types.TypeParam:
		return true
	case *types.Pointer:
		return typeContainsTypeParam(ut.Elem())
	case *types.Slice:
		return typeContainsTypeParam(ut.Elem())
	case *types.Array:
		return typeContainsTypeParam(ut.Elem())
	case *types.Map:
		return typeContainsTypeParam(ut.Key()) || typeContainsTypeParam(ut.Elem())
	case *types.Chan:
		return typeContainsTypeParam(ut.Elem())
	case *types.Named:
		if targs := ut.TypeArgs(); targs != nil {
			for i := 0; i < targs.Len(); i++ {
				if typeContainsTypeParam(targs.At(i)) {
					return true
				}
			}
		}
	case *types.Signature:
		for i := 0; i < ut.Params().Len(); i++ {
			if typeContainsTypeParam(ut.Params().At(i).Type()) {
				return true
			}
		}
		for i := 0; i < ut.Results().Len(); i++ {
			if typeContainsTypeParam(ut.Results().At(i).Type()) {
				return true
			}
		}
	}
	return false
}

func compositeSlotExpectedType(elem ast.Expr, parent ast.Node, parents map[ast.Node]ast.Node, info *types.Info) types.Type {
	switch p := parent.(type) {
	case *ast.CompositeLit:
		compT := info.Types[p].Type
		if compT == nil {
			return nil
		}
		switch ut := types.Unalias(compT).Underlying().(type) {
		case *types.Slice:
			return ut.Elem()
		case *types.Array:
			return ut.Elem()
		case *types.Struct:
			for i, e := range p.Elts {
				if e == elem && i < ut.NumFields() {
					return ut.Field(i).Type()
				}
			}
		}
	case *ast.KeyValueExpr:
		comp, ok := parents[p].(*ast.CompositeLit)
		if !ok || info.Types[comp].Type == nil {
			return nil
		}
		switch ut := types.Unalias(info.Types[comp].Type).Underlying().(type) {
		case *types.Map:
			if p.Key == elem {
				return ut.Key()
			}
			if p.Value == elem {
				return ut.Elem()
			}
		case *types.Slice:
			if p.Value == elem {
				return ut.Elem()
			}
		case *types.Array:
			if p.Value == elem {
				return ut.Elem()
			}
		case *types.Struct:
			if p.Value == elem {
				if keyId, ok := stripParens(p.Key).(*ast.Ident); ok {
					if fObj, ok := info.Uses[keyId].(*types.Var); ok && fObj != nil {
						return fObj.Type()
					}
				}
			}
		}
	}
	return nil
}

// expectedParentSlotType unwraps outer ParenExprs around expr and resolves the expected
// static destination type imposed by expr's immediate parent slot (or blankAssign=true for `_ = ...`).
func expectedParentSlotType(
	expr ast.Expr,
	parent ast.Node,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) (slotType types.Type, blankAssign bool, unwrappedExpr ast.Expr, unwrappedParent ast.Node) {
	for {
		pe, isParen := parent.(*ast.ParenExpr)
		if !isParen || pe.X != expr {
			break
		}
		expr = pe
		parent = parents[pe]
	}

	switch p := parent.(type) {
	case *ast.CallExpr:
		if p.Fun == expr || p.Ellipsis.IsValid() {
			return nil, false, expr, parent
		}
		if funTV, ok := info.Types[p.Fun]; ok && funTV.IsType() && len(p.Args) == 1 {
			return funTV.Type, false, expr, parent
		}
		fn, _ := calleeFunc(p, info)
		if fn == nil {
			return nil, false, expr, parent
		}
		sig, ok := fn.Type().(*types.Signature)
		if !ok {
			return nil, false, expr, parent
		}
		for i, a := range p.Args {
			if a == expr {
				if !isNonGenericFunc(fn) && !hasExplicitTypeArgs(p.Fun) {
					origSig, ok := fn.Origin().Type().(*types.Signature)
					if !ok || typeContainsTypeParam(formalParamTypeAt(origSig, i)) {
						return nil, false, expr, parent
					}
				}
				return formalParamTypeAt(sig, i), false, expr, parent
			}
		}

	case *ast.ValueSpec:
		if p.Type != nil {
			return info.Types[p.Type].Type, false, expr, parent
		}

	case *ast.AssignStmt:
		if p.Tok == token.ASSIGN && len(p.Lhs) == len(p.Rhs) {
			for i, r := range p.Rhs {
				if r == expr {
					if lhsId, ok := stripParens(p.Lhs[i]).(*ast.Ident); ok && lhsId.Name == "_" {
						return nil, true, expr, parent
					}
					return info.Types[p.Lhs[i]].Type, false, expr, parent
				}
			}
		}

	case *ast.ReturnStmt:
		resSig := enclosingFuncSignature(p, parents, info)
		if resSig != nil && resSig.Results().Len() == len(p.Results) {
			for i, r := range p.Results {
				if r == expr {
					return resSig.Results().At(i).Type(), false, expr, parent
				}
			}
		}

	case *ast.CompositeLit, *ast.KeyValueExpr:
		return compositeSlotExpectedType(expr, p, parents, info), false, expr, parent

	case *ast.SendStmt:
		if p.Value == expr {
			if ch, ok := types.Unalias(info.Types[p.Chan].Type).Underlying().(*types.Chan); ok {
				return ch.Elem(), false, expr, parent
			}
		}
	}
	return nil, false, expr, parent
}

func isInEmptyInterfaceSlot(
	expr ast.Expr,
	parent ast.Node,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	slotT, blankAssign, _, _ := expectedParentSlotType(expr, parent, parents, info)
	return blankAssign || isEmptyInterfaceType(slotT)
}

func exprContextPreservesInterfaceType(
	call *ast.CallExpr,
	targetIface types.Type,
	xType types.Type,
	parent ast.Node,
	parents map[ast.Node]ast.Node,
	pkg *types.Package,
	info *types.Info,
) bool {
	if types.Identical(xType, targetIface) {
		return true
	}
	slotT, blankAssign, expr, unwrappedParent := expectedParentSlotType(call, parent, parents, info)
	if blankAssign {
		return true
	}
	if sel, ok := unwrappedParent.(*ast.SelectorExpr); ok && sel.X == expr {
		obj, _, _ := types.LookupFieldOrMethod(xType, false, pkg, sel.Sel.Name)
		_, isMethod := obj.(*types.Func)
		return isMethod
	}
	return slotT != nil && isInterfaceType(slotT) && types.AssignableTo(xType, slotT)
}

func formalParamTypeAt(sig *types.Signature, argIdx int) types.Type {
	nParams := sig.Params().Len()
	if nParams == 0 {
		return nil
	}
	if sig.Variadic() && argIdx >= nParams-1 {
		if sliceParam, ok := types.Unalias(sig.Params().At(nParams - 1).Type()).(*types.Slice); ok {
			return sliceParam.Elem()
		}
		return nil
	}
	if argIdx < nParams {
		return sig.Params().At(argIdx).Type()
	}
	return nil
}

func enclosingFuncSignature(n ast.Node, parents map[ast.Node]ast.Node, info *types.Info) *types.Signature {
	for cur := n; ; {
		p := parents[cur]
		if p == nil {
			return nil
		}
		switch fn := p.(type) {
		case *ast.FuncDecl:
			if obj, ok := info.Defs[fn.Name].(*types.Func); ok {
				sig, _ := obj.Type().(*types.Signature)
				return sig
			}
			return nil
		case *ast.FuncLit:
			if sig, ok := info.Types[fn.Type].Type.(*types.Signature); ok {
				return sig
			}
			sig, _ := info.Types[fn].Type.(*types.Signature)
			return sig
		}
		cur = p
	}
}

func (p *J2gortPeepholePass) rewriteSliceIteratorLoops(
	file *ast.File,
	parents map[ast.Node]ast.Node,
	facts *peepholeFileFacts,
	info *types.Info,
) int {
	varUseCounts := make(map[*types.Var]int, len(info.Uses))
	// order-independent: commutative per-variable integer use-count accumulation
	for _, obj := range info.Uses {
		if v := varFromObj(obj); v != nil {
			varUseCounts[v]++
		}
	}

	makeRangeStmt := func(forStmt *ast.ForStmt, itVar *types.Var, iterCall *ast.CallExpr, sliceExpr ast.Expr, elemT types.Type, extraBlankUses int) *ast.RangeStmt {
		if forStmt.Post != nil || forStmt.Cond == nil || forStmt.Body == nil || len(forStmt.Body.List) == 0 {
			return nil
		}
		if !isCursorMethodCall(forStmt.Cond, itVar, "HasNext", info) {
			return nil
		}
		if sliceElemsMutatedInBody(forStmt.Body, sliceExpr, info) {
			return nil
		}
		elemIdent, dropStmts, itBodyUses, ok := matchSliceIterElemStep(forStmt.Body.List, itVar, elemT, varUseCounts, parents, info)
		if !ok || varUseCounts[itVar] != 1+extraBlankUses+itBodyUses {
			return nil
		}

		rangeTarget := sliceExpr
		if !isProvablyNonNilExpr(sliceExpr, facts, info) && !isCheckNotNullSliceCall(sliceExpr, info) {
			pkgQual := sliceIterJ2gortQualifier(iterCall)
			if pkgQual == "" {
				return nil
			}
			rangeTarget = &ast.CallExpr{
				Fun:  &ast.SelectorExpr{X: &ast.Ident{NamePos: iterCall.Pos(), Name: pkgQual}, Sel: &ast.Ident{NamePos: iterCall.Pos(), Name: "CheckNotNullSlice"}},
				Args: []ast.Expr{sliceExpr},
			}
		}

		valIdent := &ast.Ident{NamePos: forStmt.Pos(), Name: "_"}
		if elemIdent != nil {
			valIdent = &ast.Ident{NamePos: elemIdent.Pos(), Name: elemIdent.Name}
		}
		var keyExpr ast.Expr = &ast.Ident{NamePos: forStmt.Pos(), Name: "_"}
		tok := token.DEFINE
		if valIdent.Name == "_" {
			keyExpr = nil
			valIdent = nil
			tok = token.ILLEGAL
		}
		return &ast.RangeStmt{
			For:   forStmt.Pos(),
			Key:   keyExpr,
			Value: valIdent,
			Tok:   tok,
			X:     rangeTarget,
			Body: &ast.BlockStmt{
				Lbrace: forStmt.Body.Lbrace,
				List:   forStmt.Body.List[dropStmts:],
				Rbrace: forStmt.Body.Rbrace,
			},
		}
	}

	count := 0
	var walkBlock func(ast.Node) bool
	walkBlock = func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BlockStmt:
			for i := 0; i < len(x.List); i++ {
				// Pattern 1: `for it := j2gort.SliceIterator(s); it.HasNext(); { ... }`
				if fs, isFor := x.List[i].(*ast.ForStmt); isFor && fs.Init != nil {
					if initAs, ok := fs.Init.(*ast.AssignStmt); ok && initAs.Tok == token.DEFINE && len(initAs.Lhs) == 1 && len(initAs.Rhs) == 1 {
						if itId, ok := initAs.Lhs[0].(*ast.Ident); ok && itId.Name != "_" {
							if itVar := varFromObj(info.Defs[itId]); itVar != nil {
								if iterCall, sliceExpr, elemT, okS := matchSliceIteratorInitExpr(initAs.Rhs[0], info); okS {
									if rs := makeRangeStmt(fs, itVar, iterCall, sliceExpr, elemT, 0); rs != nil {
										parents[rs] = x
										parents[rs.Body] = rs
										x.List[i] = rs
										count++
										continue
									}
								}
							}
						}
					}
				}

				// Pattern 2: `var it = j2gort.SliceIterator(s); [_ = it;] [pure local decls...] for it.HasNext() { ... }`
				if itId, rhsExpr, okD := matchLocalSingleInitStmt(x.List[i]); okD && itId.Name != "_" {
					if itVar := varFromObj(info.Defs[itId]); itVar != nil {
						if iterCall, sliceExpr, elemT, okS := matchSliceIteratorInitExpr(rhsExpr, info); okS {
							scanIdx := i + 1
							blankIdx := -1
							blankUses := 0
							if scanIdx < len(x.List) && isBlankAssignOfVar(x.List[scanIdx], itVar, info) {
								blankIdx = scanIdx
								blankUses = 1
								scanIdx++
							}
							var keptMid []ast.Stmt
							for scanIdx < len(x.List) && isPureSetupStmtBeforeIterLoop(x.List[scanIdx], itVar, info) {
								keptMid = append(keptMid, x.List[scanIdx])
								scanIdx++
							}
							_ = blankIdx
							if scanIdx < len(x.List) {
								if fs, isFor := x.List[scanIdx].(*ast.ForStmt); isFor && fs.Init == nil {
									if rs := makeRangeStmt(fs, itVar, iterCall, sliceExpr, elemT, blankUses); rs != nil {
										parents[rs] = x
										parents[rs.Body] = rs
										replStmts := append(keptMid, ast.Stmt(rs))
										x.List = slices.Replace(x.List, i, scanIdx+1, replStmts...)
										count++
										continue
									}
								}
							}
						}
					}
				}
			}
		}
		return true
	}
	ast.Inspect(file, walkBlock)
	return count
}

func matchLocalSingleInitStmt(stmt ast.Stmt) (*ast.Ident, ast.Expr, bool) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if s.Tok == token.DEFINE && len(s.Lhs) == 1 && len(s.Rhs) == 1 {
			if id, ok := s.Lhs[0].(*ast.Ident); ok {
				return id, s.Rhs[0], true
			}
		}
	case *ast.DeclStmt:
		if gd, ok := s.Decl.(*ast.GenDecl); ok && gd.Tok == token.VAR && len(gd.Specs) == 1 {
			if vs, ok := gd.Specs[0].(*ast.ValueSpec); ok && len(vs.Names) == 1 && len(vs.Values) == 1 {
				return vs.Names[0], vs.Values[0], true
			}
		}
	}
	return nil, nil, false
}

func isBlankAssignOfVar(stmt ast.Stmt, v *types.Var, info *types.Info) bool {
	as, ok := stmt.(*ast.AssignStmt)
	if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
		return false
	}
	lhsId, okL := stripParens(as.Lhs[0]).(*ast.Ident)
	rhsId, okR := stripParens(as.Rhs[0]).(*ast.Ident)
	return okL && okR && lhsId.Name == "_" && varFromObj(info.Uses[rhsId]) == v
}

func isPureSetupStmtBeforeIterLoop(stmt ast.Stmt, itVar *types.Var, info *types.Info) bool {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if s.Tok == token.ASSIGN && len(s.Lhs) == 1 && len(s.Rhs) == 1 {
			if lhsId, ok := stripParens(s.Lhs[0]).(*ast.Ident); ok && lhsId.Name == "_" {
				_, okR := stripParens(s.Rhs[0]).(*ast.Ident)
				return okR && !exprUsesVar(s.Rhs[0], itVar, info)
			}
		}
		if s.Tok == token.DEFINE && len(s.Lhs) == 1 && len(s.Rhs) == 1 {
			return isPureConstExpr(s.Rhs[0]) && !exprUsesVar(s.Rhs[0], itVar, info)
		}
	case *ast.DeclStmt:
		if gd, ok := s.Decl.(*ast.GenDecl); ok && gd.Tok == token.VAR && len(gd.Specs) == 1 {
			if vs, ok := gd.Specs[0].(*ast.ValueSpec); ok && len(vs.Names) == 1 && len(vs.Values) <= 1 {
				if len(vs.Values) == 0 {
					return true
				}
				return isPureConstExpr(vs.Values[0]) && !exprUsesVar(vs.Values[0], itVar, info)
			}
		}
	}
	return false
}

func matchSliceIteratorInitExpr(expr ast.Expr, info *types.Info) (*ast.CallExpr, ast.Expr, types.Type, bool) {
	cur := stripParens(expr)
	for {
		call, ok := cur.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
			return nil, nil, nil, false
		}
		if _, isTI := isToIfaceCall(call, info); isTI {
			cur = stripParens(call.Args[0])
			continue
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "AsIface", "Cast", "UncheckedCast", "CheckNotNull", "CheckNotNullIface":
				cur = stripParens(call.Args[0])
				continue
			case "SliceIterator", "NewTypedSliceIter":
				sl, okS := types.Unalias(info.Types[call.Args[0]].Type).(*types.Slice)
				if !okS || sl.Elem() == nil || isInvalidOrUntyped(sl.Elem()) {
					return nil, nil, nil, false
				}
				return call, stripParens(call.Args[0]), sl.Elem(), true
			}
		}
		return nil, nil, nil, false
	}
}

func sliceIterJ2gortQualifier(iterCall *ast.CallExpr) string {
	fun := stripParens(iterCall.Fun)
	if idx, ok := fun.(*ast.IndexExpr); ok {
		fun = stripParens(idx.X)
	}
	if sel, ok := fun.(*ast.SelectorExpr); ok {
		if id, ok := stripParens(sel.X).(*ast.Ident); ok {
			return id.Name
		}
	}
	return ""
}

func isCheckNotNullSliceCall(expr ast.Expr, info *types.Info) bool {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
		return false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	return isJ2 && fnName == "CheckNotNullSlice"
}

func matchSliceIterElemStep(
	stmts []ast.Stmt,
	itVar *types.Var,
	elemT types.Type,
	varUseCounts map[*types.Var]int,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) (*ast.Ident, int, int, bool) {
	if len(stmts) == 0 {
		return nil, 0, 0, false
	}
	// Allow `_ = it.Next()` / `_ = it.NextTyped()` or discarded ExprStmt `it.Next()`!
	if es, isES := stmts[0].(*ast.ExprStmt); isES {
		if isCursorMethodCall(es.X, itVar, "Next", info) || isCursorMethodCall(es.X, itVar, "NextTyped", info) {
			return nil, 1, 1, true
		}
	}
	lhsId, rhsExpr, ok := matchLocalSingleInitStmt(stmts[0])
	if !ok || lhsId.Name == "_" {
		return nil, 0, 0, false
	}
	lhsVar := varFromObj(info.Defs[lhsId])
	if lhsVar == nil {
		return nil, 0, 0, false
	}

	// Direct `x := it.NextTyped()` or `x := <unbox>(it.Next())`.
	if isCursorMethodCall(rhsExpr, itVar, "NextTyped", info) && types.Identical(lhsVar.Type(), elemT) {
		return lhsId, 1, 1, true
	}
	if innerNext, okN := peelIterNextUnbox(rhsExpr, elemT, info); okN && isCursorMethodCall(innerNext, itVar, "Next", info) && types.Identical(lhsVar.Type(), elemT) {
		return lhsId, 1, 1, true
	}

	// Two-step `raw := it.Next(); x := <unbox>(raw)`.
	if isCursorMethodCall(rhsExpr, itVar, "Next", info) && len(stmts) >= 2 && varUseCounts[lhsVar] == 1 {
		if unboxedId, unboxedT, okU := matchSingleDefineUnbox(stmts[1], lhsVar, parents, info); okU && types.Identical(unboxedT, elemT) {
			return unboxedId, 2, 1, true
		}
	}

	return nil, 0, 0, false
}

func peelIterNextUnbox(expr ast.Expr, elemT types.Type, info *types.Info) (ast.Expr, bool) {
	e := stripParens(expr)
	if types.Identical(info.Types[e].Type, elemT) {
		if call, ok := e.(*ast.CallExpr); ok && len(call.Args) == 1 && !call.Ellipsis.IsValid() {
			if _, isTI := isToIfaceCall(call, info); isTI {
				return stripParens(call.Args[0]), true
			}
			if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
				switch fnName {
				case "Cast", "CastPtr", "UncheckedCast", "AsIface", "CastString", "CastNode":
					return stripParens(call.Args[0]), true
				case "UnboxBoolean", "UnboxInt", "UnboxLong", "UnboxDouble", "UnboxFloat", "UnboxChar", "UnboxByte", "UnboxShort":
					if castCall, okC := stripParens(call.Args[0]).(*ast.CallExpr); okC && len(castCall.Args) == 1 && !castCall.Ellipsis.IsValid() {
						if cFn, _, _, isCJ2 := isJ2gortCall(castCall, info); isCJ2 && (cFn == "CastPtr" || cFn == "Cast") {
							return stripParens(castCall.Args[0]), true
						}
					}
				}
			}
		}
		// Named `Node(j2gort.CastNode[uint32](it.Next()))` conversion.
		if conv, ok := e.(*ast.CallExpr); ok && len(conv.Args) == 1 && !conv.Ellipsis.IsValid() && info.Types[conv.Fun].IsType() {
			if inner, okI := peelIterNextUnbox(conv.Args[0], info.Types[conv.Args[0]].Type, info); okI {
				return inner, true
			}
		}
	}
	return nil, false
}

func sliceElemsMutatedInBody(body *ast.BlockStmt, sliceExpr ast.Expr, info *types.Info) bool {
	mut := false
	ast.Inspect(body, func(n ast.Node) bool {
		if mut {
			return false
		}
		if as, ok := n.(*ast.AssignStmt); ok {
			for _, lhs := range as.Lhs {
				if idx, isIdx := stripParens(lhs).(*ast.IndexExpr); isIdx {
					if types.Identical(info.Types[idx.X].Type, info.Types[sliceExpr].Type) {
						mut = true
						return false
					}
				}
			}
		}
		return true
	})
	return mut
}

// --- Wave 18: IfaceEquals lowering, local BitSet stack allocation, and AstNode.Release_() hook ---

func (p *J2gortPeepholePass) tryRewriteIfaceEquals(
	call *ast.CallExpr,
	parent ast.Node,
	parents map[ast.Node]ast.Node,
	pkg *types.Package,
	info *types.Info,
) bool {
	if _, isExprStmt := parent.(*ast.ExprStmt); isExprStmt {
		return false
	}

	arg0 := stripParens(call.Args[0])
	arg1 := stripParens(call.Args[1])
	nil0 := isUntypedNilExpr(arg0, info)
	nil1 := isUntypedNilExpr(arg1, info)

	if nil0 && nil1 {
		lit := &ast.Ident{NamePos: call.Pos(), Name: "true"}
		info.Types[lit] = types.TypeAndValue{Type: types.Typ[types.Bool], Value: constant.MakeBool(true)}
		return replaceExprInParent(parent, call, lit, parents)
	}

	makeIsNilDisjunction := func(x ast.Expr) *ast.BinaryExpr {
		eqNil := &ast.BinaryExpr{
			X:  cloneSimplePureExpr(x),
			Op: token.EQL,
			Y:  &ast.Ident{Name: "nil"},
		}
		isJavaNil := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   cloneSimplePureExpr(x),
				Sel: &ast.Ident{Name: "IsJavaNil_"},
			},
		}
		orExpr := &ast.BinaryExpr{
			X:  eqNil,
			Op: token.LOR,
			Y:  isJavaNil,
		}
		info.Types[eqNil] = types.TypeAndValue{Type: types.Typ[types.Bool]}
		info.Types[isJavaNil] = types.TypeAndValue{Type: types.Typ[types.Bool]}
		info.Types[orExpr] = types.TypeAndValue{Type: types.Typ[types.Bool]}
		return orExpr
	}

	if nil0 != nil1 {
		x := arg0
		if nil0 {
			x = arg1
		}
		if !isSimplePureExpr(x, info) || !isValidIfaceEqualsOperand(x, pkg, info) {
			return false
		}

		targetNode := ast.Expr(call)
		targetParent := parent
		negated := false

		cur := ast.Expr(call)
		curParent := parent
		for {
			pe, isParen := curParent.(*ast.ParenExpr)
			if !isParen || pe.X != cur {
				break
			}
			cur = pe
			curParent = parents[pe]
		}
		if unary, isNot := curParent.(*ast.UnaryExpr); isNot && unary.Op == token.NOT && unary.X == cur {
			if outer := parents[unary]; outer != nil {
				if _, isES := outer.(*ast.ExprStmt); !isES {
					targetNode = unary
					targetParent = outer
					negated = true
				}
			}
		}

		var cmp ast.Expr
		var topOp token.Token
		if negated {
			neqNil := &ast.BinaryExpr{
				X:  cloneSimplePureExpr(x),
				Op: token.NEQ,
				Y:  &ast.Ident{Name: "nil"},
			}
			notJavaNil := &ast.UnaryExpr{
				Op: token.NOT,
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   cloneSimplePureExpr(x),
						Sel: &ast.Ident{Name: "IsJavaNil_"},
					},
				},
			}
			andExpr := &ast.BinaryExpr{
				X:  neqNil,
				Op: token.LAND,
				Y:  notJavaNil,
			}
			info.Types[neqNil] = types.TypeAndValue{Type: types.Typ[types.Bool]}
			info.Types[notJavaNil] = types.TypeAndValue{Type: types.Typ[types.Bool]}
			info.Types[andExpr] = types.TypeAndValue{Type: types.Typ[types.Bool]}
			cmp = andExpr
			topOp = token.LAND
		} else {
			cmp = makeIsNilDisjunction(x)
			topOp = token.LOR
		}

		if binaryNeedsParensForNode(targetNode, targetParent, topOp) {
			pe := &ast.ParenExpr{X: cmp}
			info.Types[pe] = types.TypeAndValue{Type: types.Typ[types.Bool]}
			cmp = pe
		}
		return replaceExprInParent(targetParent, targetNode, cmp, parents)
	}

	if !isSimplePureExpr(arg0, info) || !isSimplePureExpr(arg1, info) ||
		!isValidIfaceEqualsOperand(arg0, pkg, info) || !isValidIfaceEqualsOperand(arg1, pkg, info) {
		return false
	}
	t0 := info.Types[arg0].Type
	t1 := info.Types[arg1].Type
	if !types.Comparable(t0) || !types.Comparable(t1) || (!types.AssignableTo(t0, t1) && !types.AssignableTo(t1, t0)) {
		return false
	}

	eqDirect := &ast.BinaryExpr{
		X:  cloneSimplePureExpr(arg0),
		Op: token.EQL,
		Y:  cloneSimplePureExpr(arg1),
	}
	bothNil := &ast.ParenExpr{
		X: &ast.BinaryExpr{
			X:  &ast.ParenExpr{X: makeIsNilDisjunction(arg0)},
			Op: token.LAND,
			Y:  &ast.ParenExpr{X: makeIsNilDisjunction(arg1)},
		},
	}
	fullOr := ast.Expr(&ast.BinaryExpr{
		X:  eqDirect,
		Op: token.LOR,
		Y:  bothNil,
	})
	info.Types[eqDirect] = types.TypeAndValue{Type: types.Typ[types.Bool]}
	info.Types[bothNil] = types.TypeAndValue{Type: types.Typ[types.Bool]}
	info.Types[fullOr] = types.TypeAndValue{Type: types.Typ[types.Bool]}

	if binaryNeedsParensForNode(call, parent, token.LOR) {
		pe := &ast.ParenExpr{X: fullOr}
		info.Types[pe] = types.TypeAndValue{Type: types.Typ[types.Bool]}
		fullOr = pe
	}
	return replaceExprInParent(parent, call, fullOr, parents)
}

func isUntypedNilExpr(expr ast.Expr, info *types.Info) bool {
	if expr == nil || info == nil {
		return false
	}
	tv, ok := info.Types[stripParens(expr)]
	return ok && tv.IsNil()
}

func isSimplePureExpr(expr ast.Expr, info *types.Info) bool {
	if expr == nil || info == nil {
		return false
	}
	switch e := stripParens(expr).(type) {
	case *ast.Ident:
		v, ok := info.Uses[e].(*types.Var)
		return ok && v != nil
	case *ast.SelectorExpr:
		if pkgId, ok := stripParens(e.X).(*ast.Ident); ok {
			if _, isPkg := info.Uses[pkgId].(*types.PkgName); isPkg {
				v, okV := info.Uses[e.Sel].(*types.Var)
				return okV && v != nil
			}
		}
		seln, ok := info.Selections[e]
		return ok && seln.Kind() == types.FieldVal && isSimplePureExpr(e.X, info)
	}
	return false
}

func cloneSimplePureExpr(expr ast.Expr) ast.Expr {
	switch e := stripParens(expr).(type) {
	case *ast.Ident:
		return &ast.Ident{Name: e.Name}
	case *ast.SelectorExpr:
		return &ast.SelectorExpr{
			X:   cloneSimplePureExpr(e.X),
			Sel: &ast.Ident{Name: e.Sel.Name},
		}
	}
	return expr
}

func isValidIfaceEqualsOperand(expr ast.Expr, pkg *types.Package, info *types.Info) bool {
	tv, ok := info.Types[expr]
	if !ok || tv.IsType() || tv.IsBuiltin() || tv.IsNil() || tv.Type == nil || isInvalidOrUntyped(tv.Type) {
		return false
	}
	ut := types.Unalias(tv.Type)
	if _, isTP := ut.(*types.TypeParam); isTP {
		return false
	}
	_, isPtr := ut.(*types.Pointer)
	if !isPtr && !isNonEmptyInterfaceType(ut) {
		return false
	}
	return hasZeroArgBoolMethod(tv.Type, pkg, "IsJavaNil_")
}

func binaryNeedsParensForNode(targetNode ast.Expr, parent ast.Node, op token.Token) bool {
	switch p := parent.(type) {
	case *ast.UnaryExpr, *ast.StarExpr, *ast.SelectorExpr, *ast.IndexExpr, *ast.IndexListExpr, *ast.SliceExpr, *ast.TypeAssertExpr:
		return true
	case *ast.CallExpr:
		return p.Fun == targetNode
	case *ast.BinaryExpr:
		return p.Op.Precedence() >= op.Precedence()
	}
	return false
}

type bitSetParamKey struct {
	fn       *types.Func
	paramIdx int
}

func fileHasBitSetCtorCall(file *ast.File) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}
		if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel != nil {
			switch sel.Sel.Name {
			case "NewExt_java_util_BitSet", "NewExt_java_util_BitSet_0", "NewExt_java_util_BitSet_1":
				found = true
				return false
			}
		}
		return true
	})
	return found
}

func isExtBitSetPtrType(t types.Type) bool {
	if t == nil {
		return false
	}
	ptr, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := types.Unalias(ptr.Elem()).(*types.Named)
	return ok && named.Obj() != nil && isJ2gortPkg(named.Obj().Pkg()) && named.Obj().Name() == "Ext_java_util_BitSet"
}

func isSafeBitSetArgMethodName(name string) bool {
	switch name {
	case "Or", "And", "AndNot", "Xor", "Intersects":
		return true
	}
	return false
}

func collectSafeBitSetParamFuncs(files []*ast.File, pkg *types.Package, info *types.Info) map[bitSetParamKey]bool {
	safe := make(map[bitSetParamKey]bool)
	if pkg == nil || info == nil {
		return safe
	}

	deps := make(map[bitSetParamKey][]bitSetParamKey)
	for _, f := range files {
		parents, _ := buildParentMap(f)
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Body == nil || fd.Name == nil {
				continue
			}
			fnObj, ok := info.Defs[fd.Name].(*types.Func)
			if !ok || fnObj == nil {
				continue
			}
			fnOrig := fnObj.Origin()
			sig, ok := fnOrig.Type().(*types.Signature)
			if !ok || sig.Params().Len() == 0 {
				continue
			}
			for pIdx := 0; pIdx < sig.Params().Len(); pIdx++ {
				pVar := sig.Params().At(pIdx)
				if !isExtBitSetPtrType(pVar.Type()) {
					continue
				}
				// Find the AST *types.Var for pIdx from fd.Type.Params.
				astPVar := findASTParamVarAt(fd.Type.Params, pIdx, info)
				if astPVar == nil {
					continue
				}
				calleeDeps, okP := checkBitSetParamInBody(fd.Body, astPVar, pkg, parents, info)
				if okP {
					k := bitSetParamKey{fn: fnOrig, paramIdx: pIdx}
					safe[k] = true
					deps[k] = calleeDeps
				}
			}
		}
	}

	// Fixed-point pruning: remove any (fn, paramIdx) that transitively calls an unsafe param.
	for {
		changed := false
		// order-independent: monotone fixed-point pruning converges to unique greatest fixed point
		for k := range safe {
			for _, dep := range deps[k] {
				if !safe[dep] {
					delete(safe, k)
					changed = true
					break
				}
			}
		}
		if !changed {
			break
		}
	}
	return safe
}

func findASTParamVarAt(params *ast.FieldList, targetIdx int, info *types.Info) *types.Var {
	if params == nil {
		return nil
	}
	idx := 0
	for _, field := range params.List {
		if len(field.Names) == 0 {
			if idx == targetIdx {
				return nil
			}
			idx++
			continue
		}
		for _, name := range field.Names {
			if idx == targetIdx {
				v, _ := info.Defs[name].(*types.Var)
				return v
			}
			idx++
		}
	}
	return nil
}

func checkBitSetParamInBody(
	body *ast.BlockStmt,
	targetVar *types.Var,
	pkg *types.Package,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) ([]bitSetParamKey, bool) {
	var calleeDeps []bitSetParamKey
	okAll := true

	ast.Inspect(body, func(n ast.Node) bool {
		if !okAll || n == nil {
			return false
		}
		id, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if info.Defs[id] == targetVar {
			okAll = false
			return false
		}
		if info.Uses[id] != targetVar {
			return true
		}

		dep, hasDep, okUse := checkSingleBitSetVarUse(id, pkg, parents, info, true)
		if !okUse {
			okAll = false
			return false
		}
		if hasDep {
			calleeDeps = append(calleeDeps, dep)
		}
		return true
	})

	if !okAll {
		return nil, false
	}
	return calleeDeps, true
}

func checkSingleBitSetVarUse(
	id *ast.Ident,
	pkg *types.Package,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	allowDiagnosticPanicFormat bool,
) (bitSetParamKey, bool, bool) {
	cur := ast.Expr(id)
	p := parents[cur]
	for {
		pe, isParen := p.(*ast.ParenExpr)
		if !isParen || pe.X != cur {
			break
		}
		cur = pe
		p = parents[pe]
	}

	switch parent := p.(type) {
	case *ast.AssignStmt:
		if parent.Tok == token.ASSIGN && len(parent.Lhs) == 1 && len(parent.Rhs) == 1 && parent.Rhs[0] == cur {
			if lhsId, ok := stripParens(parent.Lhs[0]).(*ast.Ident); ok && lhsId.Name == "_" {
				return bitSetParamKey{}, false, true
			}
		}
		return bitSetParamKey{}, false, false

	case *ast.SelectorExpr:
		if parent.X != cur || parent.Sel == nil || parent.Sel.Name == "Clone" {
			return bitSetParamKey{}, false, false
		}
		selCur := ast.Expr(parent)
		selParent := parents[selCur]
		for {
			pe, isParen := selParent.(*ast.ParenExpr)
			if !isParen || pe.X != selCur {
				break
			}
			selCur = pe
			selParent = parents[pe]
		}
		call, ok := selParent.(*ast.CallExpr)
		if !ok || call.Fun != selCur {
			return bitSetParamKey{}, false, false
		}
		return bitSetParamKey{}, false, true

	case *ast.CallExpr:
		if parent.Fun == cur || parent.Ellipsis.IsValid() {
			return bitSetParamKey{}, false, false
		}
		// Case 1: `<other>.Or(<v>)`, `<other>.And(<v>)`, `<other>.AndNot(<v>)`, `<other>.Xor(<v>)`, `<other>.Intersects(<v>)`.
		if len(parent.Args) == 1 && parent.Args[0] == cur {
			if funSel, ok := stripParens(parent.Fun).(*ast.SelectorExpr); ok && funSel.Sel != nil && isSafeBitSetArgMethodName(funSel.Sel.Name) {
				if isExtBitSetPtrType(info.Types[funSel.X].Type) {
					return bitSetParamKey{}, false, true
				}
			}
			// Allow `j2gort.PtrToAny(<v>)` ONLY inside an unconditional `CheckStateFmt(false, ...)` panic path in helper params.
			if allowDiagnosticPanicFormat {
				if fnName, _, _, isJ2 := isJ2gortCall(parent, info); isJ2 && fnName == "PtrToAny" {
					outerCur := ast.Expr(parent)
					outerP := parents[outerCur]
					for {
						pe, isParen := outerP.(*ast.ParenExpr)
						if !isParen || pe.X != outerCur {
							break
						}
						outerCur = pe
						outerP = parents[pe]
					}
					if outerCall, okC := outerP.(*ast.CallExpr); okC {
						if outerFn, _, _, isOuterJ2 := isJ2gortCall(outerCall, info); isOuterJ2 && outerFn == "CheckStateFmt" && len(outerCall.Args) >= 2 {
							if condId, okB := stripParens(outerCall.Args[0]).(*ast.Ident); okB && condId.Name == "false" {
								return bitSetParamKey{}, false, true
							}
						}
					}
				}
			}
		}

		// Case 2: Concrete same-package function/method call.
		calleeFn, _ := calleeFunc(parent, info)
		if calleeFn == nil || calleeFn.Pkg() != pkg {
			return bitSetParamKey{}, false, false
		}
		origFn := calleeFn.Origin()
		sig, ok := origFn.Type().(*types.Signature)
		if !ok || sig.Variadic() {
			return bitSetParamKey{}, false, false
		}
		if sig.Recv() != nil && isInterfaceType(sig.Recv().Type()) {
			return bitSetParamKey{}, false, false
		}
		for i, a := range parent.Args {
			if a == cur {
				if i < sig.Params().Len() && isExtBitSetPtrType(sig.Params().At(i).Type()) {
					return bitSetParamKey{fn: origFn, paramIdx: i}, true, true
				}
				return bitSetParamKey{}, false, false
			}
		}
	}
	return bitSetParamKey{}, false, false
}

func (p *J2gortPeepholePass) rewriteStackBitSets(
	file *ast.File,
	pkg *types.Package,
	safeBitSetParams map[bitSetParamKey]bool,
	info *types.Info,
) int {
	if !fileHasBitSetCtorCall(file) {
		return 0
	}
	parents, _ := buildParentMap(file)
	count := 0

	var walk func(ast.Node) bool
	walk = func(n ast.Node) bool {
		block, ok := n.(*ast.BlockStmt)
		if !ok {
			return true
		}
		for i := 0; i < len(block.List); i++ {
			stmt := block.List[i]
			vIdent, rhsExpr, okD := matchLocalSingleInitStmt(stmt)
			if !okD || vIdent == nil || vIdent.Name == "_" {
				continue
			}
			vObj := varFromObj(info.Defs[vIdent])
			if vObj == nil || !isExtBitSetPtrType(vObj.Type()) {
				continue
			}
			ctorCall, okC := stripParens(rhsExpr).(*ast.CallExpr)
			if !okC || ctorCall.Ellipsis.IsValid() {
				continue
			}
			fnName, fnObj, _, isJ2 := isJ2gortCall(ctorCall, info)
			if !isJ2 || fnObj == nil || fnObj.Pkg() == nil || fnObj.Pkg().Scope() == nil {
				continue
			}
			var initName string
			switch {
			case (fnName == "NewExt_java_util_BitSet" || fnName == "NewExt_java_util_BitSet_0") && len(ctorCall.Args) == 0:
				initName = "InitExt_java_util_BitSet_0"
			case fnName == "NewExt_java_util_BitSet_1" && len(ctorCall.Args) == 1:
				initName = "InitExt_java_util_BitSet_1"
			default:
				continue
			}
			if _, okInit := fnObj.Pkg().Scope().Lookup(initName).(*types.Func); !okInit {
				continue
			}
			pkgQual := sliceIterJ2gortQualifier(ctorCall)
			if pkgQual == "" {
				continue
			}
			funcBody := enclosingFuncBody(stmt, parents)
			if funcBody == nil || !isLocalBitSetVarEligible(funcBody, vIdent, vObj, pkg, safeBitSetParams, parents, info) {
				continue
			}

			bsValName := uniqueBlockLocalName(block, vIdent.Name+"_bsVal")
			pos := stmt.Pos()

			// 1. `var <v>_bsVal j2gort.Ext_java_util_BitSet`
			valDecl := &ast.DeclStmt{
				Decl: &ast.GenDecl{
					TokPos: pos,
					Tok:    token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{{NamePos: pos, Name: bsValName}},
							Type: &ast.SelectorExpr{
								X:   &ast.Ident{NamePos: pos, Name: pkgQual},
								Sel: &ast.Ident{NamePos: pos, Name: "Ext_java_util_BitSet"},
							},
						},
					},
				},
			}

			// 2. `j2gort.InitExt_java_util_BitSet_0(&<v>_bsVal)` or `InitExt_java_util_BitSet_1(&<v>_bsVal, <arg>)`
			initArgs := []ast.Expr{
				&ast.UnaryExpr{
					OpPos: pos,
					Op:    token.AND,
					X:     &ast.Ident{NamePos: pos, Name: bsValName},
				},
			}
			if len(ctorCall.Args) == 1 {
				initArgs = append(initArgs, stripParens(ctorCall.Args[0]))
			}
			initStmt := &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{NamePos: pos, Name: pkgQual},
						Sel: &ast.Ident{NamePos: pos, Name: initName},
					},
					Args: initArgs,
				},
			}

			// 3. `var <v> *j2gort.Ext_java_util_BitSet = &<v>_bsVal`
			ptrDecl := &ast.DeclStmt{
				Decl: &ast.GenDecl{
					TokPos: pos,
					Tok:    token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{{NamePos: pos, Name: vIdent.Name}},
							Type: &ast.StarExpr{
								Star: pos,
								X: &ast.SelectorExpr{
									X:   &ast.Ident{NamePos: pos, Name: pkgQual},
									Sel: &ast.Ident{NamePos: pos, Name: "Ext_java_util_BitSet"},
								},
							},
							Values: []ast.Expr{
								&ast.UnaryExpr{
									OpPos: pos,
									Op:    token.AND,
									X:     &ast.Ident{NamePos: pos, Name: bsValName},
								},
							},
						},
					},
				},
			}

			block.List = slices.Replace(block.List, i, i+1, ast.Stmt(valDecl), ast.Stmt(initStmt), ast.Stmt(ptrDecl))
			i += 2
			count++
		}
		return true
	}
	ast.Inspect(file, walk)
	return count
}

func isLocalBitSetVarEligible(
	funcBody *ast.BlockStmt,
	declIdent *ast.Ident,
	vObj *types.Var,
	pkg *types.Package,
	safeBitSetParams map[bitSetParamKey]bool,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	eligible := true
	inNestedClosure := 0

	var inspect func(ast.Node)
	inspect = func(n ast.Node) {
		if !eligible || n == nil {
			return
		}
		if fl, isFL := n.(*ast.FuncLit); isFL {
			inNestedClosure++
			if fl.Type != nil {
				ast.Inspect(fl.Type, func(tn ast.Node) bool { return true })
			}
			if fl.Body != nil {
				for _, s := range fl.Body.List {
					inspect(s)
				}
			}
			inNestedClosure--
			return
		}
		if id, isId := n.(*ast.Ident); isId {
			if info.Defs[id] == vObj {
				if id != declIdent || inNestedClosure > 0 {
					eligible = false
				}
				return
			}
			if info.Uses[id] == vObj {
				if inNestedClosure > 0 {
					eligible = false
					return
				}
				dep, hasDep, okUse := checkSingleBitSetVarUse(id, pkg, parents, info, false)
				if !okUse || (hasDep && !safeBitSetParams[dep]) {
					eligible = false
				}
				return
			}
		}
		ast.Inspect(n, func(child ast.Node) bool {
			if !eligible {
				return false
			}
			if child == n {
				return true
			}
			if _, isFL := child.(*ast.FuncLit); isFL {
				inspect(child)
				return false
			}
			if id, isId := child.(*ast.Ident); isId {
				inspect(id)
				return false
			}
			return true
		})
	}

	inspect(funcBody)
	return eligible
}

func uniqueBlockLocalName(block *ast.BlockStmt, candidate string) string {
	used := make(map[string]bool)
	ast.Inspect(block, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			used[id.Name] = true
		}
		return true
	})
	name := candidate
	for used[name] {
		name += "_"
	}
	return name
}

func (p *J2gortPeepholePass) rewriteAstNodeDeserializerRelease(
	file *ast.File,
	pkg *types.Package,
	info *types.Info,
) int {
	if !fileHasAstNodeParseFromCall(file) {
		return 0
	}
	count := 0

	var walk func(ast.Node) bool
	walk = func(n ast.Node) bool {
		block, ok := n.(*ast.BlockStmt)
		if !ok {
			return true
		}
		for i := 0; i < len(block.List); i++ {
			stmt := block.List[i]

			// Case B: `var protoVar = j2gort.Ext_..._AstNode_ParseFrom*(...)` or `protoVar := ...`
			if vIdent, rhsExpr, okD := matchLocalSingleInitStmt(stmt); okD && vIdent != nil && vIdent.Name != "_" {
				if parseCall, okP := matchAstNodeParseFromCall(rhsExpr, pkg, info); okP {
					if blockHasReleaseForIdent(block.List[i+1:], vIdent.Name) {
						continue
					}
					vObj := varFromObj(info.Defs[vIdent])
					if vObj != nil && astNodeVarEscapesBlock(block.List[i+1:], vObj, info) {
						continue
					}
					deferStmt := makeDeferReleaseStmt(vIdent.Name, parseCall.Pos())
					block.List = slices.Insert(block.List, i+1, ast.Stmt(deferStmt))
					i++
					count++
					continue
				}
			}

			// Case A: Inline `AstNode_ParseFrom*(...)` as `Args[0]` of a deserializer visit call in `stmt`
			// (e.g. `var scriptNode pkg_rhino.Node = this_.priv_Runner_visit_7484(j2gort.Ext_..._AstNode_ParseFrom_*(...), ...)`).
			if outerCall, parseCall, okInline := matchInlineAstNodeParseFromStmt(stmt, pkg, info); okInline {
				pos := parseCall.Pos()
				tmpName := uniqueBlockLocalName(block, "astRoot_")
				assignStmt := &ast.AssignStmt{
					Lhs:    []ast.Expr{&ast.Ident{NamePos: pos, Name: tmpName}},
					TokPos: pos,
					Tok:    token.DEFINE,
					Rhs:    []ast.Expr{parseCall},
				}
				deferStmt := makeDeferReleaseStmt(tmpName, pos)
				outerCall.Args[0] = &ast.Ident{NamePos: pos, Name: tmpName}
				block.List = slices.Insert(block.List, i, ast.Stmt(assignStmt), ast.Stmt(deferStmt))
				i += 2
				count++
			}
		}
		return true
	}
	ast.Inspect(file, walk)
	return count
}

func fileHasAstNodeParseFromCall(file *ast.File) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}
		if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel != nil &&
			strings.HasPrefix(sel.Sel.Name, "Ext_com_google_javascript_jscomp_serialization_AstNode_ParseFrom") {
			found = true
			return false
		}
		return true
	})
	return found
}

func matchAstNodeParseFromCall(expr ast.Expr, pkg *types.Package, info *types.Info) (*ast.CallExpr, bool) {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return nil, false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 || !strings.HasPrefix(fnName, "Ext_com_google_javascript_jscomp_serialization_AstNode_ParseFrom") {
		return nil, false
	}
	retT := info.Types[call].Type
	if retT == nil || !hasZeroArgVoidMethod(retT, pkg, "Release_") {
		return nil, false
	}
	return call, true
}

func hasZeroArgVoidMethod(t types.Type, pkg *types.Package, name string) bool {
	obj, _, _ := types.LookupFieldOrMethod(t, false, pkg, name)
	fn, ok := obj.(*types.Func)
	if !ok {
		return false
	}
	sig, ok := fn.Type().(*types.Signature)
	return ok && sig.Params().Len() == 0 && sig.Results().Len() == 0
}

func matchInlineAstNodeParseFromStmt(stmt ast.Stmt, pkg *types.Package, info *types.Info) (*ast.CallExpr, *ast.CallExpr, bool) {
	var rhs ast.Expr
	switch s := stmt.(type) {
	case *ast.DeclStmt:
		if gd, ok := s.Decl.(*ast.GenDecl); ok && gd.Tok == token.VAR && len(gd.Specs) == 1 {
			if vs, ok := gd.Specs[0].(*ast.ValueSpec); ok && len(vs.Names) == 1 && len(vs.Values) == 1 {
				rhs = vs.Values[0]
			}
		}
	case *ast.AssignStmt:
		if len(s.Lhs) == 1 && len(s.Rhs) == 1 {
			rhs = s.Rhs[0]
		}
	case *ast.ReturnStmt:
		if len(s.Results) == 1 {
			rhs = s.Results[0]
		}
	}
	if rhs == nil {
		return nil, nil, false
	}
	outerCall, ok := stripParens(rhs).(*ast.CallExpr)
	if !ok || len(outerCall.Args) == 0 || outerCall.Ellipsis.IsValid() {
		return nil, nil, false
	}
	// Ensure outerCall's Fun evaluation has no side effects before Args[0] so extracting Args[0] preserves evaluation order.
	switch fn := stripParens(outerCall.Fun).(type) {
	case *ast.Ident:
	case *ast.SelectorExpr:
		if !isSimplePureExpr(fn.X, info) {
			return nil, nil, false
		}
	default:
		return nil, nil, false
	}
	// Do not extract if outerCall returns the AstNode itself!
	if types.Identical(info.Types[outerCall].Type, info.Types[outerCall.Args[0]].Type) {
		return nil, nil, false
	}
	parseCall, ok := matchAstNodeParseFromCall(outerCall.Args[0], pkg, info)
	if !ok {
		return nil, nil, false
	}
	return outerCall, parseCall, true
}

func blockHasReleaseForIdent(stmts []ast.Stmt, varName string) bool {
	found := false
	for _, s := range stmts {
		ast.Inspect(s, func(n ast.Node) bool {
			if found {
				return false
			}
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) != 0 {
				return true
			}
			sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
			if !ok || sel.Sel == nil || sel.Sel.Name != "Release_" {
				return true
			}
			if id, ok := stripParens(sel.X).(*ast.Ident); ok && id.Name == varName {
				found = true
				return false
			}
			return true
		})
		if found {
			return true
		}
	}
	return false
}

func astNodeVarEscapesBlock(stmts []ast.Stmt, vObj *types.Var, info *types.Info) bool {
	escapes := false
	for _, s := range stmts {
		ast.Inspect(s, func(n ast.Node) bool {
			if escapes || n == nil {
				return false
			}
			switch x := n.(type) {
			case *ast.FuncLit:
				if exprUsesVarInNode(x, vObj, info) {
					escapes = true
				}
				return false
			case *ast.ReturnStmt:
				for _, r := range x.Results {
					if id, ok := stripParens(r).(*ast.Ident); ok && info.Uses[id] == vObj {
						escapes = true
						return false
					}
				}
			case *ast.AssignStmt:
				for i, r := range x.Rhs {
					if id, ok := stripParens(r).(*ast.Ident); ok && info.Uses[id] == vObj {
						if i < len(x.Lhs) {
							if lhsId, okL := stripParens(x.Lhs[i]).(*ast.Ident); okL && lhsId.Name == "_" {
								continue
							}
						}
						escapes = true
						return false
					}
				}
			}
			return true
		})
		if escapes {
			return true
		}
	}
	return false
}

func exprUsesVarInNode(node ast.Node, target *types.Var, info *types.Info) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		if id, ok := n.(*ast.Ident); ok && info.Uses[id] == target {
			found = true
			return false
		}
		return true
	})
	return found
}

func makeDeferReleaseStmt(varName string, pos token.Pos) *ast.DeferStmt {
	return &ast.DeferStmt{
		Defer: pos,
		Call: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{NamePos: pos, Name: varName},
				Sel: &ast.Ident{NamePos: pos, Name: "Release_"},
			},
		},
	}
}

func fileHasUnionTypeBuilderCtorCall(file *ast.File) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}
		if id, ok := n.(*ast.Ident); ok {
			if id.Name == "UnionType_S_Builder" || strings.HasPrefix(id.Name, "NewUnionType_Builder") {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

func isSamePkgUnionTypeBuilderPtr(t types.Type, pkg *types.Package) bool {
	if t == nil || pkg == nil {
		return false
	}
	ptr, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := types.Unalias(ptr.Elem()).(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() != pkg || named.Obj().Name() != "UnionType_Builder" {
		return false
	}
	_, isStruct := named.Underlying().(*types.Struct)
	return isStruct
}

func unionBuilderCallResultSafe(
	call *ast.CallExpr,
	builderPtrT types.Type,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	retT := info.Types[call].Type
	if retT != nil && isEmptyInterfaceType(retT) {
		return false
	}
	if retT == nil || !types.Identical(retT, builderPtrT) {
		return true
	}
	cur := ast.Expr(call)
	p := parents[cur]
	for {
		pe, isParen := p.(*ast.ParenExpr)
		if !isParen || pe.X != cur {
			break
		}
		cur = pe
		p = parents[pe]
	}
	switch parent := p.(type) {
	case *ast.ExprStmt:
		return parent.X == cur
	case *ast.AssignStmt:
		if parent.Tok == token.ASSIGN && len(parent.Lhs) == 1 && len(parent.Rhs) == 1 && parent.Rhs[0] == cur {
			if lhsId, ok := stripParens(parent.Lhs[0]).(*ast.Ident); ok && lhsId.Name == "_" {
				return true
			}
		}
		return false
	case *ast.SelectorExpr:
		if parent.X != cur || parent.Sel == nil {
			return false
		}
		selection, okSel := info.Selections[parent]
		if !okSel || selection.Kind() != types.MethodVal {
			return false
		}
		selCur := ast.Expr(parent)
		selParent := parents[selCur]
		for {
			pe, isParen := selParent.(*ast.ParenExpr)
			if !isParen || pe.X != selCur {
				break
			}
			selCur = pe
			selParent = parents[pe]
		}
		nextCall, ok := selParent.(*ast.CallExpr)
		if !ok || nextCall.Fun != selCur {
			return false
		}
		return unionBuilderCallResultSafe(nextCall, builderPtrT, parents, info)
	}
	return false
}

func checkSingleUnionTypeBuilderVarUse(
	id *ast.Ident,
	builderPtrT types.Type,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	cur := ast.Expr(id)
	p := parents[cur]
	for {
		pe, isParen := p.(*ast.ParenExpr)
		if !isParen || pe.X != cur {
			break
		}
		cur = pe
		p = parents[pe]
	}

	switch parent := p.(type) {
	case *ast.AssignStmt:
		if parent.Tok == token.ASSIGN && len(parent.Lhs) == 1 && len(parent.Rhs) == 1 && parent.Rhs[0] == cur {
			if lhsId, ok := stripParens(parent.Lhs[0]).(*ast.Ident); ok && lhsId.Name == "_" {
				return true
			}
		}
		return false

	case *ast.SelectorExpr:
		if parent.X != cur || parent.Sel == nil {
			return false
		}
		selection, okSel := info.Selections[parent]
		if !okSel || selection.Kind() != types.MethodVal {
			return false
		}
		selCur := ast.Expr(parent)
		selParent := parents[selCur]
		for {
			pe, isParen := selParent.(*ast.ParenExpr)
			if !isParen || pe.X != selCur {
				break
			}
			selCur = pe
			selParent = parents[pe]
		}
		call, ok := selParent.(*ast.CallExpr)
		if !ok || call.Fun != selCur {
			return false
		}
		return unionBuilderCallResultSafe(call, builderPtrT, parents, info)
	}
	return false
}

func isLocalUnionTypeBuilderVarEligible(
	funcBody *ast.BlockStmt,
	declIdent *ast.Ident,
	vObj *types.Var,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	eligible := true
	inNestedClosure := 0
	builderPtrT := vObj.Type()

	var inspect func(ast.Node)
	inspect = func(n ast.Node) {
		if !eligible || n == nil {
			return
		}
		if fl, isFL := n.(*ast.FuncLit); isFL {
			inNestedClosure++
			if fl.Body != nil {
				for _, s := range fl.Body.List {
					inspect(s)
				}
			}
			inNestedClosure--
			return
		}
		if id, isId := n.(*ast.Ident); isId {
			if info.Defs[id] == vObj {
				if id != declIdent || inNestedClosure > 0 {
					eligible = false
				}
				return
			}
			if info.Uses[id] == vObj {
				if inNestedClosure > 0 || !checkSingleUnionTypeBuilderVarUse(id, builderPtrT, parents, info) {
					eligible = false
				}
				return
			}
		}
		ast.Inspect(n, func(child ast.Node) bool {
			if !eligible {
				return false
			}
			if child == n {
				return true
			}
			if _, isFL := child.(*ast.FuncLit); isFL {
				inspect(child)
				return false
			}
			if id, isId := child.(*ast.Ident); isId {
				inspect(id)
				return false
			}
			return true
		})
	}

	inspect(funcBody)
	return eligible
}

func (p *J2gortPeepholePass) rewriteStackUnionTypeBuilders(
	file *ast.File,
	pkg *types.Package,
	info *types.Info,
) int {
	if pkg == nil || pkg.Scope() == nil || !fileHasUnionTypeBuilderCtorCall(file) {
		return 0
	}
	parents, _ := buildParentMap(file)
	count := 0

	var walk func(ast.Node) bool
	walk = func(n ast.Node) bool {
		block, ok := n.(*ast.BlockStmt)
		if !ok {
			return true
		}
		for i := 0; i < len(block.List); i++ {
			stmt := block.List[i]
			vIdent, rhsExpr, okD := matchLocalSingleInitStmt(stmt)
			if !okD || vIdent == nil || vIdent.Name == "_" {
				continue
			}
			vObj := varFromObj(info.Defs[vIdent])
			if vObj == nil || !isSamePkgUnionTypeBuilderPtr(vObj.Type(), pkg) {
				continue
			}
			ctorCall, okC := stripParens(rhsExpr).(*ast.CallExpr)
			if !okC || ctorCall.Ellipsis.IsValid() {
				continue
			}
			calleeFn, _ := calleeFunc(ctorCall, info)
			if calleeFn == nil || calleeFn.Pkg() != pkg {
				continue
			}
			calleeSig, okSig := calleeFn.Type().(*types.Signature)
			if !okSig || calleeSig.Recv() != nil || calleeSig.Variadic() {
				continue
			}

			var initName string
			switch {
			case strings.HasPrefix(calleeFn.Name(), "NewUnionType_Builder"):
				initName = "Init" + strings.TrimPrefix(calleeFn.Name(), "New")
			case calleeFn.Name() == "UnionType_S_Builder" && len(ctorCall.Args) == 1:
				initName = "InitUnionType_Builder_1_JSTypeRegistry"
			default:
				continue
			}

			initFn, okInit := pkg.Scope().Lookup(initName).(*types.Func)
			if !okInit {
				continue
			}
			initSig, okInitSig := initFn.Type().(*types.Signature)
			if !okInitSig || initSig.Recv() != nil || initSig.Variadic() ||
				initSig.Results().Len() != 0 || initSig.Params().Len() != 1+len(ctorCall.Args) ||
				!types.Identical(initSig.Params().At(0).Type(), vObj.Type()) {
				continue
			}
			argsOK := true
			for k, arg := range ctorCall.Args {
				argT := info.Types[arg].Type
				if argT == nil || !types.AssignableTo(argT, initSig.Params().At(k+1).Type()) {
					argsOK = false
					break
				}
			}
			if !argsOK {
				continue
			}

			funcBody := enclosingFuncBody(stmt, parents)
			if funcBody == nil || !isLocalUnionTypeBuilderVarEligible(funcBody, vIdent, vObj, parents, info) {
				continue
			}

			ubValName := uniqueBlockLocalName(block, vIdent.Name+"_ubVal")
			pos := stmt.Pos()

			// 1. `var <v>_ubVal UnionType_Builder`
			valDecl := &ast.DeclStmt{
				Decl: &ast.GenDecl{
					TokPos: pos,
					Tok:    token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{{NamePos: pos, Name: ubValName}},
							Type:  &ast.Ident{NamePos: pos, Name: "UnionType_Builder"},
						},
					},
				},
			}

			// 2. `InitUnionType_Builder...(&<v>_ubVal, <args...>)`
			initArgs := make([]ast.Expr, 0, 1+len(ctorCall.Args))
			initArgs = append(initArgs, &ast.UnaryExpr{
				OpPos: pos,
				Op:    token.AND,
				X:     &ast.Ident{NamePos: pos, Name: ubValName},
			})
			for _, arg := range ctorCall.Args {
				initArgs = append(initArgs, stripParens(arg))
			}
			initStmt := &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun:  &ast.Ident{NamePos: pos, Name: initName},
					Args: initArgs,
				},
			}

			// 3. `var <v> *UnionType_Builder = &<v>_ubVal`
			ptrDecl := &ast.DeclStmt{
				Decl: &ast.GenDecl{
					TokPos: pos,
					Tok:    token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{{NamePos: pos, Name: vIdent.Name}},
							Type: &ast.StarExpr{
								Star: pos,
								X:    &ast.Ident{NamePos: pos, Name: "UnionType_Builder"},
							},
							Values: []ast.Expr{
								&ast.UnaryExpr{
									OpPos: pos,
									Op:    token.AND,
									X:     &ast.Ident{NamePos: pos, Name: ubValName},
								},
							},
						},
					},
				},
			}

			block.List = slices.Replace(block.List, i, i+1, ast.Stmt(valDecl), ast.Stmt(initStmt), ast.Stmt(ptrDecl))
			i += 2
			count++
		}
		return true
	}
	ast.Inspect(file, walk)
	return count
}

// --- Wave 20C Task 4: Unbox base.LinkedIdentityHashMap[K, V] onto flat j2gort.IdentityOrderedMap[K, V] ---

func (p *J2gortPeepholePass) rewriteLinkedIdentityHashMap(ctx *PassContext) (int, error) {
	if ctx.Pkg == nil || ctx.Pkg.Name() != "base" || len(ctx.Files) == 0 {
		return 0, nil
	}
	info := ctx.Info
	var targetFile *ast.File
	var targetSpec *ast.TypeSpec
	var targetStruct *ast.StructType
	var mapNamed *types.Named
	var kName, vName string
	var innerMapTag *ast.BasicLit
	var idHashField *ast.Field

	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, sp := range gd.Specs {
				ts, ok := sp.(*ast.TypeSpec)
				if !ok || ts.Assign.IsValid() || ts.Name == nil || ts.Name.Name != "LinkedIdentityHashMap" {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok || st.Fields == nil || (len(st.Fields.List) != 3 && len(st.Fields.List) != 4) {
					continue
				}
				if len(st.Fields.List) == 4 {
					f3 := st.Fields.List[3]
					if len(f3.Names) != 1 || f3.Names[0].Name != "idHash_" {
						continue
					}
				}
				tpNames := extractFlatParamNames(ts.TypeParams)
				if len(tpNames) != 2 {
					continue
				}
				f0, f1, f2 := st.Fields.List[0], st.Fields.List[1], st.Fields.List[2]
				if len(f0.Names) != 1 || f0.Names[0].Name != "innerMap" ||
					len(f1.Names) != 1 || f1.Names[0].Name != "head" ||
					len(f2.Names) != 1 || f2.Names[0].Name != "tail" {
					continue
				}
				tn, ok := info.Defs[ts.Name].(*types.TypeName)
				if !ok || tn == nil {
					continue
				}
				named := namedOrigin(tn.Type())
				if named == nil {
					continue
				}
				f0Var := varFromObj(info.Defs[f0.Names[0]])
				if f0Var == nil {
					continue
				}
				iomNamed, ok := types.Unalias(f0Var.Type()).(*types.Named)
				if !ok || iomNamed.Obj() == nil || !isJ2gortPkg(iomNamed.Obj().Pkg()) || iomNamed.Obj().Name() != "IdentityOrderedMap" {
					continue
				}
				targs := iomNamed.TypeArgs()
				if targs == nil || targs.Len() != 2 {
					continue
				}
				entryNamed := ptrNamedOrigin(targs.At(1))
				if entryNamed == nil || entryNamed.Obj() == nil || entryNamed.Obj().Name() != "LinkedIdentityHashMap_Entry" {
					continue
				}
				targetFile = file
				targetSpec = ts
				targetStruct = st
				mapNamed = named
				kName, vName = tpNames[0], tpNames[1]
				innerMapTag = f0.Tag
				if len(st.Fields.List) == 4 {
					idHashField = st.Fields.List[3]
				}
				break
			}
		}
		if targetFile != nil {
			break
		}
	}
	if targetFile == nil || mapNamed == nil || targetSpec == nil {
		return 0, nil
	}

	j2Qual, ok := j2gortQualifierForFile(targetFile, ctx)
	if !ok || j2Qual == "" {
		return 0, nil
	}
	if !j2gortHasExport(ctx, "NewIdentityOrderedMap") {
		return 0, nil
	}

	// Ensure every method on *LinkedIdentityHashMap in ctx.Files that directly accesses innerMap, head, or tail
	// is in our exact supported method set so no custom linked-list mutation can be missed.
	for _, file := range ctx.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Body == nil || len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) == 0 {
				continue
			}
			rv := varFromObj(info.Defs[fn.Recv.List[0].Names[0]])
			if rv == nil || ptrNamedOrigin(rv.Type()) != mapNamed {
				continue
			}
			touchesRepr := false
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if touchesRepr || n == nil {
					return false
				}
				if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel != nil {
					switch sel.Sel.Name {
					case "innerMap", "head", "tail":
						if fv := varFromObj(info.Uses[sel.Sel]); fv != nil && fv.IsField() {
							touchesRepr = true
							return false
						}
					}
				}
				return true
			})
			if !touchesRepr {
				continue
			}
			if file != targetFile {
				return 0, nil
			}
			switch fn.Name.Name {
			case "initFields_LinkedIdentityHashMap", "Get", "Put", "ForEach",
				"Remove", "Clear", "ContainsKey", "Size", "IsEmpty":
			default:
				return 0, nil
			}
		}
	}

	// Rewrite all methods touching innerMap/head/tail in targetFile.
	type methodRewrite struct {
		fn      *ast.FuncDecl
		newBody *ast.BlockStmt
	}
	var rewrites []methodRewrite

	for _, decl := range targetFile.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || fn.Body == nil || len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 {
			continue
		}
		recvId := fn.Recv.List[0].Names[0]
		rv := varFromObj(info.Defs[recvId])
		if rv == nil || ptrNamedOrigin(rv.Type()) != mapNamed || recvId.Name == "_" {
			continue
		}
		recvTpNames := extractRecvTypeParamNames(fn.Recv.List[0].Type)
		mk, mv := kName, vName
		if len(recvTpNames) == 2 {
			mk, mv = recvTpNames[0], recvTpNames[1]
		}
		pNames := extractFlatParamNames(fn.Type.Params)
		recv := recvId.Name
		var bodySrc string
		switch fn.Name.Name {
		case "initFields_LinkedIdentityHashMap":
			if len(pNames) != 0 {
				return 0, nil
			}
			bodySrc = fmt.Sprintf("_ = %[1]s\n%[1]s.innerMap = %[2]s.NewIdentityOrderedMap[%[3]s, %[4]s]()", recv, j2Qual, mk, mv)
		case "Get":
			if len(pNames) != 1 {
				return 0, nil
			}
			bodySrc = fmt.Sprintf(
				"if v_, ok_ := %[1]s.innerMap.GetOk(%[2]s); ok_ {\nreturn v_\n}\nif _, ok_ := any((*%[4]s)(nil)).(*string); ok_ {\nreturn any(%[3]s.NullString).(%[4]s)\n}\nvar z_ %[4]s\nreturn z_",
				recv, pNames[0], j2Qual, mv,
			)
		case "Put":
			if len(pNames) != 2 {
				return 0, nil
			}
			bodySrc = fmt.Sprintf(
				"if old_, ok_ := %[1]s.innerMap.PutOk(%[2]s, %[3]s); ok_ {\nreturn old_\n}\nif _, ok_ := any((*%[5]s)(nil)).(*string); ok_ {\nreturn any(%[4]s.NullString).(%[5]s)\n}\nvar z_ %[5]s\nreturn z_",
				recv, pNames[0], pNames[1], j2Qual, mv,
			)
		case "ForEach":
			if len(pNames) != 1 {
				return 0, nil
			}
			bodySrc = fmt.Sprintf(
				"cur_ := %[1]s.innerMap.IterEntries()\nfor cur_.NextEntry() {\n%[2]s.Accept(%[3]s.ToAny(cur_.Key()), %[3]s.ToAny(cur_.Value()))\n}",
				recv, pNames[0], j2Qual,
			)
		case "Remove":
			if len(pNames) != 1 {
				return 0, nil
			}
			bodySrc = fmt.Sprintf(
				"if old_, ok_ := %[1]s.innerMap.RemoveOk(%[2]s); ok_ {\nreturn old_\n}\nif _, ok_ := any((*%[4]s)(nil)).(*string); ok_ {\nreturn any(%[3]s.NullString).(%[4]s)\n}\nvar z_ %[4]s\nreturn z_",
				recv, pNames[0], j2Qual, mv,
			)
		case "Clear":
			if len(pNames) != 0 {
				return 0, nil
			}
			bodySrc = fmt.Sprintf("%s.innerMap.Clear()", recv)
		case "ContainsKey":
			if len(pNames) != 1 {
				return 0, nil
			}
			bodySrc = fmt.Sprintf("return %s.innerMap.ContainsKey(%s)", recv, pNames[0])
		case "Size":
			if len(pNames) != 0 {
				return 0, nil
			}
			bodySrc = fmt.Sprintf("return %s.innerMap.Size()", recv)
		case "IsEmpty":
			if len(pNames) != 0 {
				return 0, nil
			}
			bodySrc = fmt.Sprintf("return %s.innerMap.IsEmpty()", recv)
		default:
			continue
		}
		pf, err := parser.ParseFile(ctx.Fset, "", "package p\nfunc f() {\n"+bodySrc+"\n}\n", 0)
		if err != nil {
			return 0, err
		}
		rewrites = append(rewrites, methodRewrite{fn: fn, newBody: pf.Decls[0].(*ast.FuncDecl).Body})
	}

	if len(rewrites) == 0 {
		return 0, nil
	}

	// Update LinkedIdentityHashMap[K, V] struct fields to single unboxed `innerMap j2gort.IdentityOrderedMap[K, V]` (+ optional `idHash_`).
	newFields := []*ast.Field{{
		Names: []*ast.Ident{{Name: "innerMap"}},
		Type: &ast.IndexListExpr{
			X:       &ast.SelectorExpr{X: &ast.Ident{Name: j2Qual}, Sel: &ast.Ident{Name: "IdentityOrderedMap"}},
			Indices: []ast.Expr{&ast.Ident{Name: kName}, &ast.Ident{Name: vName}},
		},
		Tag: innerMapTag,
	}}
	if idHashField != nil {
		newFields = append(newFields, idHashField)
	}
	targetStruct.Fields.List = newFields
	for _, rw := range rewrites {
		rw.fn.Body = rw.newBody
	}
	ctx.MarkFileModified(targetFile)
	return 1 + len(rewrites), nil
}

func extractRecvTypeParamNames(expr ast.Expr) []string {
	star, ok := stripParens(expr).(*ast.StarExpr)
	if !ok {
		return nil
	}
	switch ix := stripParens(star.X).(type) {
	case *ast.IndexExpr:
		if id, ok := stripParens(ix.Index).(*ast.Ident); ok {
			return []string{id.Name}
		}
	case *ast.IndexListExpr:
		out := make([]string, 0, len(ix.Indices))
		for _, idx := range ix.Indices {
			id, ok := stripParens(idx).(*ast.Ident)
			if !ok {
				return nil
			}
			out = append(out, id.Name)
		}
		return out
	}
	return nil
}

// --- Wave 20C Task 2.4: Elide Redundant Zero-Value Field Writes in Fresh-Only initFields_* Methods ---

func (p *J2gortPeepholePass) elideZeroedInitFieldsWrites(ctx *PassContext) int {
	if len(ctx.Files) == 0 || ctx.Info == nil || ctx.Pkg == nil {
		return 0
	}
	info := ctx.Info

	type initFieldsCand struct {
		file       *ast.File
		fnDecl     *ast.FuncDecl
		fnOrig     *types.Func
		recvVar    *types.Var
		ownerNamed *types.Named
		zeroIdxs   []int
	}
	var cands []initFieldsCand
	candByFn := make(map[*types.Func]int)

	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Body == nil || !strings.HasPrefix(fn.Name.Name, "initFields_") {
				continue
			}
			if len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 {
				continue
			}
			if (fn.Type.Params != nil && len(fn.Type.Params.List) > 0) || (fn.Type.Results != nil && len(fn.Type.Results.List) > 0) {
				continue
			}
			recvId := fn.Recv.List[0].Names[0]
			recvVar := varFromObj(info.Defs[recvId])
			fnObj, okF := info.Defs[fn.Name].(*types.Func)
			if !okF || fnObj == nil || recvVar == nil || recvId.Name == "_" {
				continue
			}
			owner := ptrNamedOrigin(recvVar.Type())
			if owner == nil {
				continue
			}
			zeroIdxs := findRedundantZeroAssignStmts(fn.Body, recvVar, owner, info)
			if len(zeroIdxs) == 0 {
				continue
			}
			fnOrig := funcOrigin(fnObj)
			candByFn[fnOrig] = len(cands)
			cands = append(cands, initFieldsCand{
				file:       file,
				fnDecl:     fn,
				fnOrig:     fnOrig,
				recvVar:    recvVar,
				ownerNamed: owner,
				zeroIdxs:   zeroIdxs,
			})
		}
	}

	if len(cands) == 0 {
		return 0
	}

	// Collect concrete structs in ctx.Pkg embedding each *Base struct so standalone *Base types
	// with only cross-package subclasses are conservatively excluded unless embedded in ctx.Pkg.
	embeddedBaseHasLocalConcrete := make(map[*types.Named]bool)
	for _, file := range ctx.Files {
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, sp := range gd.Specs {
				ts, ok := sp.(*ast.TypeSpec)
				if !ok || ts.Assign.IsValid() || ts.Name == nil || strings.HasSuffix(ts.Name.Name, "Base") {
					continue
				}
				tn, ok := info.Defs[ts.Name].(*types.TypeName)
				if !ok || tn == nil {
					continue
				}
				n := namedOrigin(tn.Type())
				if n == nil {
					continue
				}
				for _, emb := range collectDirectEmbeddedNamed(n) {
					embeddedBaseHasLocalConcrete[emb] = true
				}
			}
		}
	}

	// Catalog all package-level Init* functions `func Init*(p0 *T, ...)`.
	type initFuncMeta struct {
		fnOrig     *types.Func
		fnDecl     *ast.FuncDecl
		p0Var      *types.Var
		targetType *types.Named
		callCount  int
		nonFresh   bool
		callers    []*types.Func // Init* callers chaining into this Init*
	}
	var initFuncs []*initFuncMeta
	initByOrig := make(map[*types.Func]*initFuncMeta)

	for _, file := range ctx.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || fn.Body == nil || !strings.HasPrefix(fn.Name.Name, "Init") {
				continue
			}
			if fn.Type.Params == nil || len(fn.Type.Params.List) == 0 || len(fn.Type.Params.List[0].Names) == 0 {
				continue
			}
			fnObj, ok := info.Defs[fn.Name].(*types.Func)
			if !ok || fnObj == nil {
				continue
			}
			p0Var := varFromObj(info.Defs[fn.Type.Params.List[0].Names[0]])
			if p0Var == nil {
				continue
			}
			target := ptrNamedOrigin(p0Var.Type())
			if target == nil {
				continue
			}
			fnOrig := funcOrigin(fnObj)
			nonFresh := strings.HasSuffix(target.Obj().Name(), "Base") && !embeddedBaseHasLocalConcrete[target]
			im := &initFuncMeta{
				fnOrig:     fnOrig,
				fnDecl:     fn,
				p0Var:      p0Var,
				targetType: target,
				nonFresh:   nonFresh,
			}
			initFuncs = append(initFuncs, im)
			initByOrig[fnOrig] = im
		}
	}

	parentsByFile := make(map[*ast.File]map[ast.Node]ast.Node, len(ctx.Files))
	for _, f := range ctx.Files {
		parentsByFile[f], _ = buildParentMap(f)
	}

	candCallCount := make([]int, len(cands))
	candDisq := make([]bool, len(cands))
	candInitCallers := make([][]*types.Func, len(cands))

	// Inspect every reference to Init* functions and candidate initFields_* methods across ctx.Files.
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		ast.Inspect(file, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			usedFn, ok := info.Uses[id].(*types.Func)
			if !ok || usedFn == nil {
				return true
			}
			uOrig := funcOrigin(usedFn)

			if im := initByOrig[uOrig]; im != nil {
				call := directCallFromCalleeIdent(id, parents)
				if call == nil || call.Ellipsis.IsValid() || len(call.Args) < 1 {
					im.nonFresh = true
					return true
				}
				im.callCount++
				enclFn := enclosingFuncObj(call, parents, info)
				enclOrig := funcOrigin(enclFn)
				if callerIM := initByOrig[enclOrig]; callerIM != nil {
					if isChainedFreshInitArg(call.Args[0], call, callerIM.fnDecl.Body, callerIM.p0Var, parents, info) {
						im.callers = append(im.callers, callerIM.fnOrig)
						return true
					}
				}
				fnBody := enclosingFuncBody(call, parents)
				if fnBody != nil && isFreshZeroedPtrArgAtCall(call.Args[0], call, fnBody, parents, info) {
					return true
				}
				im.nonFresh = true
				return true
			}

			if cIdx, isCand := candByFn[uOrig]; isCand {
				sel, okS := parents[id].(*ast.SelectorExpr)
				if !okS || sel.Sel != id {
					candDisq[cIdx] = true
					return true
				}
				call, okC := parents[sel].(*ast.CallExpr)
				if !okC || call.Fun != sel || len(call.Args) != 0 {
					candDisq[cIdx] = true
					return true
				}
				if _, isES := parents[call].(*ast.ExprStmt); !isES {
					candDisq[cIdx] = true
					return true
				}
				candCallCount[cIdx]++
				enclFn := enclosingFuncObj(call, parents, info)
				enclOrig := funcOrigin(enclFn)
				if callerIM := initByOrig[enclOrig]; callerIM != nil {
					if rId, okId := stripParens(sel.X).(*ast.Ident); okId && varFromObj(info.Uses[rId]) == callerIM.p0Var &&
						noFieldMutationBeforeCall(callerIM.fnDecl.Body, callerIM.p0Var, call, parents, info) {
						candInitCallers[cIdx] = append(candInitCallers[cIdx], callerIM.fnOrig)
						return true
					}
				}
				fnBody := enclosingFuncBody(call, parents)
				if fnBody != nil && isFreshZeroedPtrArgAtCall(sel.X, call, fnBody, parents, info) {
					return true
				}
				candDisq[cIdx] = true
			}
			return true
		})
	}

	// Require at least 1 in-package call site for every fresh Init* and propagate nonFresh transitively.
	for _, im := range initFuncs {
		if im.callCount == 0 {
			im.nonFresh = true
		}
	}
	for changed := true; changed; {
		changed = false
		for _, im := range initFuncs {
			if im.nonFresh {
				continue
			}
			for _, callerOrig := range im.callers {
				if callerIM := initByOrig[callerOrig]; callerIM == nil || callerIM.nonFresh {
					im.nonFresh = true
					changed = true
					break
				}
			}
		}
	}

	totalElided := 0
	for i, c := range cands {
		if candDisq[i] || candCallCount[i] == 0 {
			continue
		}
		allFresh := true
		for _, callerOrig := range candInitCallers[i] {
			if callerIM := initByOrig[callerOrig]; callerIM == nil || callerIM.nonFresh {
				allFresh = false
				break
			}
		}
		if !allFresh {
			continue
		}
		dropSet := make(map[int]bool, len(c.zeroIdxs))
		for _, idx := range c.zeroIdxs {
			dropSet[idx] = true
		}
		newList := make([]ast.Stmt, 0, len(c.fnDecl.Body.List)-len(c.zeroIdxs))
		for idx, st := range c.fnDecl.Body.List {
			if !dropSet[idx] {
				newList = append(newList, st)
			}
		}
		c.fnDecl.Body.List = newList
		ctx.MarkFileModified(c.file)
		totalElided += len(c.zeroIdxs)
	}
	return totalElided
}

func findRedundantZeroAssignStmts(body *ast.BlockStmt, recvVar *types.Var, owner *types.Named, info *types.Info) []int {
	st, ok := owner.Underlying().(*types.Struct)
	if !ok {
		return nil
	}
	directFields := make(map[*types.Var]bool, st.NumFields())
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if !f.Embedded() {
			directFields[f.Origin()] = true
		}
	}
	seenField := make(map[*types.Var]bool)
	var zeroIdxs []int
	for idx, stmt := range body.List {
		if isBlankDiscardStmt(stmt) {
			continue
		}
		as, ok := stmt.(*ast.AssignStmt)
		if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
			// Stop scanning past any arbitrary control flow or method call that might mutate fields.
			break
		}
		sel, ok := stripParens(as.Lhs[0]).(*ast.SelectorExpr)
		if !ok {
			break
		}
		rId, ok := stripParens(sel.X).(*ast.Ident)
		if !ok || varFromObj(info.Uses[rId]) != recvVar {
			break
		}
		fv := varFromObj(info.Uses[sel.Sel])
		if fv == nil || !directFields[fv] {
			break
		}
		if !seenField[fv] && isExactGoZeroExprForType(as.Rhs[0], fv.Type(), info) && !exprHasSideEffects(as.Rhs[0], info) {
			zeroIdxs = append(zeroIdxs, idx)
		}
		seenField[fv] = true
	}
	return zeroIdxs
}

func isExactGoZeroExprForType(expr ast.Expr, fieldType types.Type, info *types.Info) bool {
	if expr == nil || fieldType == nil || info == nil {
		return false
	}
	e := stripParens(expr)
	tv, ok := info.Types[e]
	if !ok {
		return false
	}
	switch u := types.Unalias(fieldType).Underlying().(type) {
	case *types.Pointer, *types.Slice, *types.Map, *types.Interface, *types.Chan, *types.Signature:
		return tv.IsNil()
	case *types.Basic:
		if tv.Value == nil || !isPureConstExpr(e) {
			return false
		}
		switch {
		case (u.Info() & types.IsBoolean) != 0:
			return tv.Value.Kind() == constant.Bool && !constant.BoolVal(tv.Value)
		case (u.Info() & types.IsInteger) != 0:
			return tv.Value.Kind() == constant.Int && constant.Sign(tv.Value) == 0
		case (u.Info() & types.IsString) != 0:
			return tv.Value.Kind() == constant.String && constant.StringVal(tv.Value) == ""
		}
	}
	return false
}

func exprHasSideEffects(expr ast.Expr, info *types.Info) bool {
	if isPureConstExpr(expr) {
		return false
	}
	e := stripParens(expr)
	if id, ok := e.(*ast.Ident); ok {
		_, isNil := info.Uses[id].(*types.Nil)
		return !isNil
	}
	return true
}

func directCallFromCalleeIdent(id *ast.Ident, parents map[ast.Node]ast.Node) *ast.CallExpr {
	cur := ast.Node(id)
	if sel, ok := parents[cur].(*ast.SelectorExpr); ok && sel.Sel == id {
		cur = sel
	}
	switch ix := parents[cur].(type) {
	case *ast.IndexExpr:
		if ix.X == cur {
			cur = ix
		}
	case *ast.IndexListExpr:
		if ix.X == cur {
			cur = ix
		}
	}
	for {
		pe, ok := parents[cur].(*ast.ParenExpr)
		if !ok || pe.X != cur {
			break
		}
		cur = pe
	}
	if call, ok := parents[cur].(*ast.CallExpr); ok && call.Fun == cur {
		return call
	}
	return nil
}

func isChainedFreshInitArg(
	arg0 ast.Expr,
	call *ast.CallExpr,
	body *ast.BlockStmt,
	p0Var *types.Var,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	a := stripParens(arg0)
	validTarget := false
	if id, ok := a.(*ast.Ident); ok && varFromObj(info.Uses[id]) == p0Var {
		validTarget = true
	} else if un, ok := a.(*ast.UnaryExpr); ok && un.Op == token.AND {
		if sel, ok := stripParens(un.X).(*ast.SelectorExpr); ok {
			if id, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == p0Var {
				if fv := varFromObj(info.Uses[sel.Sel]); fv != nil && fv.Embedded() {
					validTarget = true
				}
			}
		}
	}
	if !validTarget {
		return false
	}
	return noFieldMutationBeforeCall(body, p0Var, call, parents, info)
}

func noFieldMutationBeforeCall(
	body *ast.BlockStmt,
	targetVar *types.Var,
	call *ast.CallExpr,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	if body == nil {
		return false
	}
	// call must be a top-level ExprStmt in body (not inside a loop or branch).
	es, ok := parents[call].(*ast.ExprStmt)
	if !ok || parents[es] != body {
		return false
	}
	for _, st := range body.List {
		if st == es {
			return true
		}
		if isHarmlessPreInitStmt(st, targetVar, info) {
			continue
		}
		return false
	}
	return false
}

func isHarmlessPreInitStmt(st ast.Stmt, targetVar *types.Var, info *types.Info) bool {
	if isBlankDiscardStmt(st) {
		return true
	}
	switch s := st.(type) {
	case *ast.ExprStmt:
		call, ok := stripParens(s.X).(*ast.CallExpr)
		if !ok {
			return false
		}
		if id, ok := stripParens(call.Fun).(*ast.Ident); ok && len(call.Args) == 0 && strings.HasPrefix(id.Name, "clinit_") {
			return true
		}
		if sel, ok := stripParens(call.Fun).(*ast.SelectorExpr); ok && sel.Sel.Name == "BindSelf_" && len(call.Args) == 1 {
			if rId, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[rId]) == targetVar {
				return true
			}
		}
		// Allow superclass Init*(&this_.Base, ...) before this_.initFields_Sub() because embedded Base is a disjoint sub-struct!
		if callee, _ := calleeFunc(call, info); callee != nil && strings.HasPrefix(callee.Name(), "Init") && len(call.Args) >= 1 {
			if un, ok := stripParens(call.Args[0]).(*ast.UnaryExpr); ok && un.Op == token.AND {
				if sel, ok := stripParens(un.X).(*ast.SelectorExpr); ok {
					if rId, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[rId]) == targetVar {
						if fv := varFromObj(info.Uses[sel.Sel]); fv != nil && fv.Embedded() {
							return true
						}
					}
				}
			}
		}
	case *ast.AssignStmt:
		// Allow `this_.outer_ = ...` or `this_.Base.outer_ = ...`.
		if s.Tok == token.ASSIGN && len(s.Lhs) == 1 && len(s.Rhs) == 1 {
			if sel, ok := stripParens(s.Lhs[0]).(*ast.SelectorExpr); ok && sel.Sel.Name == "outer_" {
				if rId, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[rId]) == targetVar {
					return true
				}
				if embSel, ok := stripParens(sel.X).(*ast.SelectorExpr); ok {
					if rId, ok := stripParens(embSel.X).(*ast.Ident); ok && varFromObj(info.Uses[rId]) == targetVar {
						if fv := varFromObj(info.Uses[embSel.Sel]); fv != nil && fv.Embedded() {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func isFreshZeroedPtrArgAtCall(
	arg0 ast.Expr,
	call *ast.CallExpr,
	fnBody *ast.BlockStmt,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	es, ok := parents[call].(*ast.ExprStmt)
	if !ok {
		return false
	}
	blk, ok := parents[es].(*ast.BlockStmt)
	if !ok {
		return false
	}
	// Ensure blk is not inside a loop relative to fnBody so `var val T` is zeroed on every execution.
	for cur := ast.Node(blk); cur != nil && cur != fnBody; cur = parents[cur] {
		switch cur.(type) {
		case *ast.ForStmt, *ast.RangeStmt:
			return false
		}
	}
	callIdx := slices.Index(blk.List, ast.Stmt(es))
	if callIdx <= 0 {
		return false
	}

	a := stripParens(arg0)
	// Case A: `&val` where `var val T` (with 0 initializers) is declared earlier in `blk`.
	if un, ok := a.(*ast.UnaryExpr); ok && un.Op == token.AND {
		valId, ok := stripParens(un.X).(*ast.Ident)
		if !ok {
			return false
		}
		valVar := varFromObj(info.Uses[valId])
		if valVar == nil || valVar.IsField() {
			return false
		}
		declIdx := -1
		for i := 0; i < callIdx; i++ {
			ds, ok := blk.List[i].(*ast.DeclStmt)
			if !ok {
				continue
			}
			gd, ok := ds.Decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR || len(gd.Specs) != 1 {
				continue
			}
			vs, ok := gd.Specs[0].(*ast.ValueSpec)
			if ok && len(vs.Names) == 1 && len(vs.Values) == 0 && varFromObj(info.Defs[vs.Names[0]]) == valVar {
				declIdx = i
				break
			}
		}
		if declIdx < 0 {
			return false
		}
		for i := declIdx + 1; i < callIdx; i++ {
			if !isHarmlessPreInitStmt(blk.List[i], valVar, info) {
				return false
			}
		}
		return true
	}

	// Case B: `this_` where `this_ := &T{}` or `this_ := <slab>.Alloc()` (or `allocOpt_*` if/else) earlier in `blk`.
	ptrId, ok := a.(*ast.Ident)
	if !ok {
		return false
	}
	ptrVar := varFromObj(info.Uses[ptrId])
	if ptrVar == nil || ptrVar.IsField() {
		return false
	}
	initIdx := -1
	for i := 0; i < callIdx; i++ {
		switch s := blk.List[i].(type) {
		case *ast.AssignStmt:
			if s.Tok == token.DEFINE && len(s.Lhs) == 1 && len(s.Rhs) == 1 {
				if lhsId, ok := s.Lhs[0].(*ast.Ident); ok && varFromObj(info.Defs[lhsId]) == ptrVar && isFreshZeroedPtrRhs(s.Rhs[0], info) {
					initIdx = i
				}
			}
		case *ast.DeclStmt:
			gd, ok := s.Decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR || len(gd.Specs) != 1 {
				continue
			}
			vs, ok := gd.Specs[0].(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || varFromObj(info.Defs[vs.Names[0]]) != ptrVar {
				continue
			}
			if len(vs.Values) == 1 && isFreshZeroedPtrRhs(vs.Values[0], info) {
				initIdx = i
			} else if len(vs.Values) == 0 && i+1 < callIdx {
				if isAllocOptIfElseStmt(blk.List[i+1], ptrVar, info) {
					initIdx = i + 1
					i++
				}
			}
		}
		if initIdx >= 0 {
			break
		}
	}
	if initIdx < 0 {
		return false
	}
	for i := initIdx + 1; i < callIdx; i++ {
		if !isHarmlessPreInitStmt(blk.List[i], ptrVar, info) {
			return false
		}
	}
	return true
}

func isFreshZeroedPtrRhs(rhs ast.Expr, info *types.Info) bool {
	r := stripParens(rhs)
	if un, ok := r.(*ast.UnaryExpr); ok && un.Op == token.AND {
		comp, ok := stripParens(un.X).(*ast.CompositeLit)
		return ok && len(comp.Elts) == 0
	}
	if call, ok := r.(*ast.CallExpr); ok && len(call.Args) == 0 && !call.Ellipsis.IsValid() {
		if sel, ok := stripParens(call.Fun).(*ast.SelectorExpr); ok && sel.Sel != nil && sel.Sel.Name == "Alloc" {
			if seln, ok := info.Selections[sel]; ok {
				if recvN := derefNamedOrigin(seln.Recv()); recvN != nil && recvN.Obj() != nil && isJ2gortPkg(recvN.Obj().Pkg()) {
					nm := recvN.Obj().Name()
					return nm == "Slab" || nm == "PooledSlab"
				}
			}
		}
	}
	return false
}

func isAllocOptIfElseStmt(stmt ast.Stmt, ptrVar *types.Var, info *types.Info) bool {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok || ifStmt.Init != nil || ifStmt.Body == nil || len(ifStmt.Body.List) != 1 {
		return false
	}
	elseBlk, ok := ifStmt.Else.(*ast.BlockStmt)
	if !ok || len(elseBlk.List) != 1 {
		return false
	}
	checkSingleAssign := func(s ast.Stmt) bool {
		as, ok := s.(*ast.AssignStmt)
		if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
			return false
		}
		id, ok := stripParens(as.Lhs[0]).(*ast.Ident)
		return ok && varFromObj(info.Uses[id]) == ptrVar && isFreshZeroedPtrRhs(as.Rhs[0], info)
	}
	return checkSingleAssign(ifStmt.Body.List[0]) && checkSingleAssign(elseBlk.List[0])
}

func (p *J2gortPeepholePass) rewriteNodeQualifiedNameJoinDot(ctx *PassContext) int {
	if ctx == nil || ctx.Pkg == nil || ctx.Pkg.Name() != "rhino" || ctx.Info == nil || ctx.Importer == nil {
		return 0
	}
	total := 0
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		nsQual, nsPkg := findNodestoreImportForFile(file, ctx)
		if nsQual == "" || nsPkg == nil || nsPkg.Scope() == nil {
			continue
		}
		hasJoinDot := nsPkg.Scope().Lookup("JoinDot") != nil
		hasJoinDot3 := nsPkg.Scope().Lookup("JoinDot3") != nil
		hasJoinDot4 := nsPkg.Scope().Lookup("JoinDot4") != nil
		if !hasJoinDot && !hasJoinDot3 && !hasJoinDot4 {
			continue
		}
		parents, postOrder := buildParentMap(file)
		fileCount := 0
		for _, n := range postOrder {
			bin, ok := n.(*ast.BinaryExpr)
			if !ok || bin.Op != token.ADD {
				continue
			}
			// Only match top-level `+` chains so nested `+` operands fold once as a single chain.
			topExpr := ast.Expr(bin)
			for {
				pe, isParen := parents[topExpr].(*ast.ParenExpr)
				if !isParen || pe.X != topExpr {
					break
				}
				topExpr = pe
			}
			if parentBin, isPB := parents[topExpr].(*ast.BinaryExpr); isPB && parentBin.Op == token.ADD {
				continue
			}
			parts, ok := matchDotSeparatedStringAddChain(bin, ctx.Info)
			if !ok {
				continue
			}
			parent := parents[topExpr]
			if parent == nil {
				continue
			}
			repl, ok := buildJoinDotCallChain(nsQual, topExpr.Pos(), parts, hasJoinDot, hasJoinDot3, hasJoinDot4)
			if !ok {
				continue
			}
			if replaceExprInParent(parent, topExpr, repl, parents) {
				fileCount++
			}
		}
		if fileCount > 0 {
			ctx.MarkFileModified(file)
			total += fileCount
		}
	}
	return total
}

func buildJoinDotCallChain(nsQual string, pos token.Pos, parts []ast.Expr, hasJoinDot, hasJoinDot3, hasJoinDot4 bool) (*ast.CallExpr, bool) {
	if len(parts) < 2 {
		return nil, false
	}
	mkCall := func(helper string, args ...ast.Expr) *ast.CallExpr {
		return &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{NamePos: pos, Name: nsQual},
				Sel: &ast.Ident{NamePos: pos, Name: helper},
			},
			Args: args,
		}
	}
	var acc *ast.CallExpr
	i := 0
	switch {
	case len(parts) >= 4 && hasJoinDot4:
		acc = mkCall("JoinDot4", parts[0], parts[1], parts[2], parts[3])
		i = 4
	case len(parts) >= 3 && hasJoinDot3:
		acc = mkCall("JoinDot3", parts[0], parts[1], parts[2])
		i = 3
	case len(parts) >= 2 && hasJoinDot:
		acc = mkCall("JoinDot", parts[0], parts[1])
		i = 2
	default:
		return nil, false
	}
	for i < len(parts) {
		rem := len(parts) - i
		switch {
		case rem >= 3 && hasJoinDot4:
			acc = mkCall("JoinDot4", acc, parts[i], parts[i+1], parts[i+2])
			i += 3
		case rem >= 2 && hasJoinDot3:
			acc = mkCall("JoinDot3", acc, parts[i], parts[i+1])
			i += 2
		case rem >= 1 && hasJoinDot:
			acc = mkCall("JoinDot", acc, parts[i])
			i++
		default:
			return nil, false
		}
	}
	return acc, true
}

func findNodestoreImportForFile(file *ast.File, ctx *PassContext) (string, *types.Package) {
	if file == nil || ctx == nil || ctx.Importer == nil {
		return "", nil
	}
	for _, spec := range file.Imports {
		if spec.Path == nil || len(spec.Path.Value) < 3 {
			continue
		}
		if spec.Name != nil && (spec.Name.Name == "_" || spec.Name.Name == ".") {
			continue
		}
		impPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		impPkg, err := ctx.Importer.Import(impPath)
		if err != nil || impPkg == nil || impPkg.Name() != "nodestore" {
			continue
		}
		qual := impPkg.Name()
		if spec.Name != nil && spec.Name.Name != "" {
			qual = spec.Name.Name
		}
		return qual, impPkg
	}
	return "", nil
}

func matchDotSeparatedStringAddChain(bin *ast.BinaryExpr, info *types.Info) ([]ast.Expr, bool) {
	var flat []ast.Expr
	cur := ast.Expr(bin)
	for {
		b, isB := stripParens(cur).(*ast.BinaryExpr)
		if !isB || b.Op != token.ADD {
			flat = append(flat, stripParens(cur))
			break
		}
		flat = append(flat, stripParens(b.Y))
		cur = b.X
	}
	slices.Reverse(flat)
	if len(flat) < 3 || len(flat)%2 != 1 {
		return nil, false
	}
	var parts []ast.Expr
	for i, e := range flat {
		if i%2 == 1 {
			lit, isLit := e.(*ast.BasicLit)
			if !isLit || lit.Kind != token.STRING {
				return nil, false
			}
			s, err := strconv.Unquote(lit.Value)
			if err != nil || s != "." {
				return nil, false
			}
			continue
		}
		if !hasExactOrDefaultStringType(e, info) {
			return nil, false
		}
		parts = append(parts, e)
	}
	return parts, true
}

func (p *J2gortPeepholePass) rewriteNodeClearJSTypeAndJSDocFastPaths(ctx *PassContext) int {
	if ctx == nil || ctx.Pkg == nil || ctx.Pkg.Name() != "rhino" || ctx.Pkg.Scope() == nil || ctx.Info == nil || ctx.Importer == nil {
		return 0
	}
	rhinoNodeObj := ctx.Pkg.Scope().Lookup("Node")
	if rhinoNodeObj == nil || rhinoNodeObj.Type() == nil {
		return 0
	}
	rhinoNodeType := types.Unalias(rhinoNodeObj.Type())
	hasNodeNullPtrEx := ctx.Pkg.Scope().Lookup("nodeNullPtrEx") != nil

	total := 0
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		nsQual, nsPkg := findNodestoreImportForFile(file, ctx)
		if nsQual == "" || nsPkg == nil || nsPkg.Scope() == nil {
			continue
		}
		nsNodeObj := nsPkg.Scope().Lookup("Node")
		if nsNodeObj == nil || nsNodeObj.Type() == nil {
			continue
		}
		nsNodeType := types.Unalias(nsNodeObj.Type())
		hasClearCast := peepholeTypeHasMethod(nsNodeType, false, nsPkg, "ClearJSTypeBeforeCast")
		hasClearDoc := peepholeTypeHasMethod(nsNodeType, false, nsPkg, "ClearJSDoc")
		if !hasClearCast && !hasClearDoc {
			continue
		}
		hasSetCastObj := peepholeTypeHasMethod(nsNodeType, false, nsPkg, "SetJSTypeBeforeCastObj")
		hasSetDocObj := peepholeTypeHasMethod(nsNodeType, false, nsPkg, "SetJSDocObj")

		fileCount := 0

		// 1. Rewrite direct `nodestore.Node(x).SetJSTypeBeforeCastObj(nil)` -> `.ClearJSTypeBeforeCast()`
		// and `nodestore.Node(x).SetJSDocObj(nil)` -> `.ClearJSDoc()`.
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || call.Ellipsis.IsValid() || len(call.Args) != 1 {
				return true
			}
			sel, isSel := stripParens(call.Fun).(*ast.SelectorExpr)
			if !isSel || sel.Sel == nil {
				return true
			}
			var clearName string
			switch sel.Sel.Name {
			case "SetJSTypeBeforeCastObj":
				if hasClearCast {
					clearName = "ClearJSTypeBeforeCast"
				}
			case "SetJSDocObj":
				if hasClearDoc {
					clearName = "ClearJSDoc"
				}
			}
			if clearName == "" {
				return true
			}
			argId, isId := stripParens(call.Args[0]).(*ast.Ident)
			if !isId || argId.Name != "nil" {
				return true
			}
			if used := ctx.Info.Uses[argId]; used != nil && used != types.Universe.Lookup("nil") {
				return true
			}
			recvT := ctx.Info.TypeOf(sel.X)
			if recvT == nil || !types.Identical(types.Unalias(recvT), nsNodeType) {
				return true
			}
			sel.Sel.Name = clearName
			call.Args = nil
			fileCount++
			return true
		})

		// 2. Add typed-nil / nil fast-path guards on `(this_ Node) SetJSTypeBeforeCast(arg)`
		// and `(this_ Node) SetJSDocInfo(arg)` before boxing `arg` into `any`.
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Body == nil || fn.Name == nil ||
				len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 {
				continue
			}
			var clearMethod, nodeSetHelper, nsSetObjMethod string
			isTypeBeforeCast := false
			switch fn.Name.Name {
			case "SetJSTypeBeforeCast":
				if !hasClearCast {
					continue
				}
				clearMethod = "ClearJSTypeBeforeCast"
				nodeSetHelper = "nodeSet_typeBeforeCast"
				if hasSetCastObj {
					nsSetObjMethod = "SetJSTypeBeforeCastObj"
				}
				isTypeBeforeCast = true
			case "SetJSDocInfo":
				if !hasClearDoc {
					continue
				}
				clearMethod = "ClearJSDoc"
				nodeSetHelper = "nodeSet_jsDocInfo"
				if hasSetDocObj {
					nsSetObjMethod = "SetJSDocObj"
				}
			default:
				continue
			}
			if peepholeASTContainsSelector(fn.Body, clearMethod) {
				continue
			}
			recvId := fn.Recv.List[0].Names[0]
			if recvId.Name == "" || recvId.Name == "_" {
				continue
			}
			recvVar := varFromObj(ctx.Info.Defs[recvId])
			if recvVar == nil || !types.Identical(types.Unalias(recvVar.Type()), rhinoNodeType) {
				continue
			}
			if fn.Type.Params == nil || len(fn.Type.Params.List) != 1 || len(fn.Type.Params.List[0].Names) != 1 {
				continue
			}
			paramId := fn.Type.Params.List[0].Names[0]
			if paramId.Name == "" || paramId.Name == "_" {
				continue
			}
			paramVar := varFromObj(ctx.Info.Defs[paramId])
			if paramVar == nil {
				continue
			}
			if peepholeVarAssignedInBlock(fn.Body, recvVar, ctx.Info) || peepholeVarAssignedInBlock(fn.Body, paramVar, ctx.Info) {
				continue
			}
			retStmtSrc := "return"
			if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
				if len(fn.Type.Results.List) != 1 || len(fn.Type.Results.List[0].Names) != 0 {
					continue
				}
				resT := ctx.Info.TypeOf(fn.Type.Results.List[0].Type)
				if resT == nil || !types.Identical(types.Unalias(resT), rhinoNodeType) {
					continue
				}
				retStmtSrc = "return " + recvId.Name
			}

			paramType := paramVar.Type()
			var condSrc string
			if types.IsInterface(paramType) && peepholeTypeHasMethod(paramType, true, ctx.Pkg, "IsJavaNil_") {
				condSrc = fmt.Sprintf("%[1]s == nil || %[1]s.IsJavaNil_()", paramId.Name)
			} else if types.AssignableTo(types.Typ[types.UntypedNil], paramType) {
				condSrc = fmt.Sprintf("%s == nil", paramId.Name)
			} else {
				continue
			}

			helperIdx := -1
			for idx, st := range fn.Body.List {
				es, isES := st.(*ast.ExprStmt)
				if !isES {
					continue
				}
				c, isC := stripParens(es.X).(*ast.CallExpr)
				if !isC || c.Ellipsis.IsValid() || len(c.Args) != 1 {
					continue
				}
				s, isS := stripParens(c.Fun).(*ast.SelectorExpr)
				if !isS || s.Sel == nil || s.Sel.Name != nodeSetHelper {
					continue
				}
				rId, isR := stripParens(s.X).(*ast.Ident)
				aId, isA := stripParens(c.Args[0]).(*ast.Ident)
				if isR && isA && rId.Name == recvId.Name && aId.Name == paramId.Name {
					helperIdx = idx
					break
				}
			}

			var guardLines []string
			if isTypeBeforeCast && helperIdx >= 0 && hasNodeNullPtrEx && !peepholeASTContainsIdent(fn.Body, "nodeNullPtrEx") {
				guardLines = append(guardLines, fmt.Sprintf("if %s == 0 {\n\tpanic(nodeNullPtrEx)\n}", recvId.Name))
			}
			guardLines = append(guardLines, fmt.Sprintf(
				"if %s {\n\t%s.Node(%s).%s()\n\t%s\n}",
				condSrc, nsQual, recvId.Name, clearMethod, retStmtSrc,
			))

			pf, err := parser.ParseFile(ctx.Fset, "", "package p\nfunc f() {\n"+strings.Join(guardLines, "\n")+"\n}\n", 0)
			if err != nil {
				continue
			}
			guardStmts := pf.Decls[0].(*ast.FuncDecl).Body.List

			if helperIdx >= 0 && nsSetObjMethod != "" {
				setPF, setErr := parser.ParseFile(ctx.Fset, "", fmt.Sprintf(
					"package p\nfunc f() {\n%s.Node(%s).%s(%s)\n}\n",
					nsQual, recvId.Name, nsSetObjMethod, paramId.Name,
				), 0)
				if setErr == nil {
					fn.Body.List[helperIdx] = setPF.Decls[0].(*ast.FuncDecl).Body.List[0]
				}
			}

			fn.Body.List = slices.Insert(fn.Body.List, 0, guardStmts...)
			fileCount++
		}

		if fileCount > 0 {
			ctx.MarkFileModified(file)
			total += fileCount
		}
	}
	return total
}

func peepholeTypeHasMethod(t types.Type, addressable bool, pkg *types.Package, name string) bool {
	if t == nil {
		return false
	}
	obj, _, _ := types.LookupFieldOrMethod(t, addressable, pkg, name)
	_, ok := obj.(*types.Func)
	return ok
}

func peepholeASTContainsSelector(root ast.Node, selName string) bool {
	if root == nil {
		return false
	}
	found := false
	ast.Inspect(root, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel != nil && sel.Sel.Name == selName {
			found = true
			return false
		}
		return true
	})
	return found
}

func peepholeASTContainsIdent(root ast.Node, name string) bool {
	if root == nil {
		return false
	}
	found := false
	ast.Inspect(root, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if id, ok := n.(*ast.Ident); ok && id.Name == name {
			found = true
			return false
		}
		return true
	})
	return found
}

func peepholeVarAssignedInBlock(body *ast.BlockStmt, target *types.Var, info *types.Info) bool {
	if body == nil || target == nil || info == nil {
		return false
	}
	assigned := false
	ast.Inspect(body, func(n ast.Node) bool {
		if assigned || n == nil {
			return false
		}
		switch x := n.(type) {
		case *ast.AssignStmt:
			for _, lhs := range x.Lhs {
				if id, ok := stripParens(lhs).(*ast.Ident); ok && varFromObj(info.Uses[id]) == target {
					assigned = true
					return false
				}
			}
		case *ast.IncDecStmt:
			if id, ok := stripParens(x.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == target {
				assigned = true
				return false
			}
		case *ast.UnaryExpr:
			if x.Op == token.AND {
				if id, ok := stripParens(x.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == target {
					assigned = true
					return false
				}
			}
		}
		return true
	})
	return assigned
}

// --- Wave 29C Task 1: Zero-Allocation Node.GetAncestors() Enhanced-For Loop Lowering ---

func fileHasGetAncestorsCall(file *ast.File) bool {
	if file == nil {
		return false
	}
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel != nil && sel.Sel.Name == "GetAncestors" {
			found = true
			return false
		}
		return true
	})
	return found
}

func hasZeroArgSameTypeMethod(t types.Type, pkg *types.Package, name string) bool {
	if t == nil {
		return false
	}
	obj, _, _ := types.LookupFieldOrMethod(t, false, pkg, name)
	fn, ok := obj.(*types.Func)
	if !ok {
		return false
	}
	sig, ok := fn.Type().(*types.Signature)
	return ok && sig.Params().Len() == 0 && sig.Results().Len() == 1 && types.Identical(sig.Results().At(0).Type(), t)
}

func matchNodeGetAncestorsCall(expr ast.Expr, pkg *types.Package, info *types.Info) (ast.Expr, types.Type, bool) {
	if expr == nil || info == nil {
		return nil, nil, false
	}
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() || len(call.Args) != 0 {
		return nil, nil, false
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || sel.Sel.Name != "GetAncestors" {
		return nil, nil, false
	}
	nodeT := info.Types[sel.X].Type
	if !isNamedUint32NodeType(nodeT) || !hasZeroArgSameTypeMethod(nodeT, pkg, "GetParent") {
		return nil, nil, false
	}
	return stripParens(sel.X), nodeT, true
}

func unwrapIterCursorTarget(expr ast.Expr, info *types.Info) ast.Expr {
	cur := stripParens(expr)
	for {
		call, ok := cur.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
			return cur
		}
		if _, isTI := isToIfaceCall(call, info); isTI {
			cur = stripParens(call.Args[0])
			continue
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "CheckNotNullPtr", "CheckNotNull", "CheckNotNullIface", "AsIface", "Cast", "CastPtr", "UncheckedCast":
				cur = stripParens(call.Args[0])
				continue
			}
		}
		return cur
	}
}

func matchNewIterCursorCall(expr ast.Expr, info *types.Info) (*ast.CallExpr, *types.Func, ast.Expr, bool) {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
		return nil, nil, nil, false
	}
	fnName, fnObj, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 || fnName != "NewIterCursor" {
		return nil, nil, nil, false
	}
	return call, fnObj, unwrapIterCursorTarget(call.Args[0], info), true
}

func peelCursorNodeReadExpr(expr ast.Expr, curVar *types.Var, info *types.Info) bool {
	cur := stripParens(expr)
	for {
		if isCursorMethodCall(cur, curVar, "NextNode", info) || isCursorMethodCall(cur, curVar, "Next", info) {
			return true
		}
		call, ok := cur.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
			return false
		}
		if tv, okT := info.Types[call.Fun]; okT && tv.IsType() {
			t := types.Unalias(tv.Type)
			if isNamedUint32NodeType(t) || types.Identical(t, types.Typ[types.Uint32]) {
				cur = stripParens(call.Args[0])
				continue
			}
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "UnboxNode", "CastNode", "Cast", "UncheckedCast":
				cur = stripParens(call.Args[0])
				continue
			}
		}
		return false
	}
}

func peelRawVarNodeReadExpr(expr ast.Expr, rawVar *types.Var, info *types.Info) bool {
	cur := stripParens(expr)
	for {
		if id, ok := cur.(*ast.Ident); ok && varFromObj(info.Uses[id]) == rawVar {
			return true
		}
		call, ok := cur.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
			return false
		}
		if tv, okT := info.Types[call.Fun]; okT && tv.IsType() {
			t := types.Unalias(tv.Type)
			if isNamedUint32NodeType(t) || types.Identical(t, types.Typ[types.Uint32]) {
				cur = stripParens(call.Args[0])
				continue
			}
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "UnboxNode", "CastNode", "Cast", "UncheckedCast":
				cur = stripParens(call.Args[0])
				continue
			}
		}
		return false
	}
}

func matchAncestorsIterElemStep(
	stmts []ast.Stmt,
	curVar *types.Var,
	nodeT types.Type,
	varUseCounts map[*types.Var]int,
	info *types.Info,
) (elemIdent *ast.Ident, elemIsUint32 bool, dropStmts int, curBodyUses int, ok bool) {
	if len(stmts) == 0 {
		return nil, false, 0, 0, false
	}
	if es, isES := stmts[0].(*ast.ExprStmt); isES {
		if peelCursorNodeReadExpr(es.X, curVar, info) {
			return nil, false, 1, 1, true
		}
		return nil, false, 0, 0, false
	}
	if as, isAS := stmts[0].(*ast.AssignStmt); isAS && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 {
		if lhsId, isId := stripParens(as.Lhs[0]).(*ast.Ident); isId && lhsId.Name == "_" && peelCursorNodeReadExpr(as.Rhs[0], curVar, info) {
			return nil, false, 1, 1, true
		}
	}

	lhsId, rhsExpr, okInit := matchLocalSingleInitStmt(stmts[0])
	if !okInit || lhsId == nil {
		return nil, false, 0, 0, false
	}
	if lhsId.Name == "_" {
		if peelCursorNodeReadExpr(rhsExpr, curVar, info) {
			return nil, false, 1, 1, true
		}
		return nil, false, 0, 0, false
	}
	lhsVar := varFromObj(info.Defs[lhsId])
	if lhsVar == nil {
		return nil, false, 0, 0, false
	}

	if peelCursorNodeReadExpr(rhsExpr, curVar, info) {
		if types.Identical(lhsVar.Type(), nodeT) {
			return lhsId, false, 1, 1, true
		}
		if types.Identical(lhsVar.Type(), types.Typ[types.Uint32]) {
			return lhsId, true, 1, 1, true
		}
	}

	if isCursorMethodCall(rhsExpr, curVar, "Next", info) && len(stmts) >= 2 && varUseCounts[lhsVar] == 1 {
		ancId, ancRhs, ok2 := matchLocalSingleInitStmt(stmts[1])
		if ok2 && ancId != nil && ancId.Name != "_" {
			if ancVar := varFromObj(info.Defs[ancId]); ancVar != nil && peelRawVarNodeReadExpr(ancRhs, lhsVar, info) {
				if types.Identical(ancVar.Type(), nodeT) {
					return ancId, false, 2, 1, true
				}
				if types.Identical(ancVar.Type(), types.Typ[types.Uint32]) {
					return ancId, true, 2, 1, true
				}
			}
		}
	}

	return nil, false, 0, 0, false
}

func (p *J2gortPeepholePass) rewriteNodeGetAncestorsLoops(
	file *ast.File,
	pkg *types.Package,
	info *types.Info,
) int {
	if !fileHasGetAncestorsCall(file) || info == nil {
		return 0
	}
	parents, _ := buildParentMap(file)
	varUseCounts := make(map[*types.Var]int, len(info.Uses))
	// order-independent: commutative per-variable integer use-count accumulation
	for _, obj := range info.Uses {
		if v := varFromObj(obj); v != nil {
			varUseCounts[v]++
		}
	}

	buildLoweredAncestorsStmts := func(
		forStmt *ast.ForStmt,
		curVar *types.Var,
		curBlankUses int,
		iterCall *ast.CallExpr,
		iterFn *types.Func,
		nodeExpr ast.Expr,
		nodeT types.Type,
		fallbackScope ast.Node,
	) ([]ast.Stmt, bool) {
		if forStmt.Post != nil || forStmt.Cond == nil || forStmt.Body == nil || len(forStmt.Body.List) == 0 {
			return nil, false
		}
		if !isCursorMethodCall(forStmt.Cond, curVar, "HasNext", info) {
			return nil, false
		}
		elemIdent, elemIsUint32, dropStmts, curBodyUses, okStep := matchAncestorsIterElemStep(
			forStmt.Body.List, curVar, nodeT, varUseCounts, info,
		)
		if !okStep || varUseCounts[curVar] != 1+curBlankUses+curBodyUses {
			return nil, false
		}
		pkgQual := sliceIterJ2gortQualifier(iterCall)
		if pkgQual == "" {
			return nil, false
		}

		scopeRoot := ast.Node(enclosingFuncBody(forStmt, parents))
		if scopeRoot == nil {
			scopeRoot = fallbackScope
		}
		curAncName := uniqueNodeLocalName(scopeRoot, "curAnc_")
		pos := iterCall.Pos()

		// 1. `curAnc_ := <nodeExpr>`
		initStmt := &ast.AssignStmt{
			Lhs:    []ast.Expr{&ast.Ident{NamePos: pos, Name: curAncName}},
			TokPos: pos,
			Tok:    token.DEFINE,
			Rhs:    []ast.Expr{nodeExpr},
		}

		// 2. `if curAnc_ == 0 { j2gort.ThrowNPE() }`
		var npeStmt ast.Stmt
		hasThrowNPE := iterFn != nil && iterFn.Pkg() != nil && iterFn.Pkg().Scope() != nil &&
			iterFn.Pkg().Scope().Lookup("ThrowNPE") != nil
		hasCheckNode := iterFn != nil && iterFn.Pkg() != nil && iterFn.Pkg().Scope() != nil &&
			iterFn.Pkg().Scope().Lookup("CheckNotNullNode") != nil
		if hasThrowNPE {
			npeStmt = &ast.IfStmt{
				If: pos,
				Cond: &ast.BinaryExpr{
					X:     &ast.Ident{NamePos: pos, Name: curAncName},
					OpPos: pos,
					Op:    token.EQL,
					Y:     &ast.BasicLit{ValuePos: pos, Kind: token.INT, Value: "0"},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.Ident{NamePos: pos, Name: pkgQual},
									Sel: &ast.Ident{NamePos: pos, Name: "ThrowNPE"},
								},
							},
						},
					},
				},
			}
		} else if hasCheckNode {
			npeStmt = &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{NamePos: pos, Name: pkgQual},
						Sel: &ast.Ident{NamePos: pos, Name: "CheckNotNullNode"},
					},
					Args: []ast.Expr{&ast.Ident{NamePos: pos, Name: curAncName}},
				},
			}
		} else {
			return nil, false
		}

		// 3. `curAnc_ = curAnc_.GetParent()` (Node.getAncestors() starts at this_.getParent()!)
		stepParentStmt := func(pPos token.Pos) *ast.AssignStmt {
			return &ast.AssignStmt{
				Lhs:    []ast.Expr{&ast.Ident{NamePos: pPos, Name: curAncName}},
				TokPos: pPos,
				Tok:    token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{NamePos: pPos, Name: curAncName},
							Sel: &ast.Ident{NamePos: pPos, Name: "GetParent"},
						},
					},
				},
			}
		}

		// 4. `for curAnc_ != 0 { anc := curAnc_; curAnc_ = curAnc_.GetParent(); <rest...> }`
		forPos := forStmt.Pos()
		var newBodyList []ast.Stmt
		if elemIdent != nil {
			var rhs ast.Expr = &ast.Ident{NamePos: elemIdent.Pos(), Name: curAncName}
			if elemIsUint32 {
				rhs = &ast.CallExpr{
					Fun:  &ast.Ident{NamePos: elemIdent.Pos(), Name: "uint32"},
					Args: []ast.Expr{rhs},
				}
			}
			newBodyList = append(newBodyList, &ast.AssignStmt{
				Lhs:    []ast.Expr{&ast.Ident{NamePos: elemIdent.Pos(), Name: elemIdent.Name}},
				TokPos: elemIdent.Pos(),
				Tok:    token.DEFINE,
				Rhs:    []ast.Expr{rhs},
			})
		}
		newBodyList = append(newBodyList, stepParentStmt(forPos))
		newBodyList = append(newBodyList, forStmt.Body.List[dropStmts:]...)

		newFor := &ast.ForStmt{
			For: forPos,
			Cond: &ast.BinaryExpr{
				X:     &ast.Ident{NamePos: forPos, Name: curAncName},
				OpPos: forPos,
				Op:    token.NEQ,
				Y:     &ast.BasicLit{ValuePos: forPos, Kind: token.INT, Value: "0"},
			},
			Body: &ast.BlockStmt{
				Lbrace: forStmt.Body.Lbrace,
				List:   newBodyList,
				Rbrace: forStmt.Body.Rbrace,
			},
		}

		return []ast.Stmt{initStmt, npeStmt, stepParentStmt(pos), newFor}, true
	}

	matchCursorAndForInSlice := func(stmts []ast.Stmt, startIdx int, expectedTargetVar *types.Var) (
		forStmt *ast.ForStmt,
		curVar *types.Var,
		curBlankUses int,
		iterCall *ast.CallExpr,
		iterFn *types.Func,
		nodeExpr ast.Expr,
		nodeT types.Type,
		endIdx int,
		ok bool,
	) {
		if startIdx >= len(stmts) {
			return
		}
		curId, rhsExpr, okD := matchLocalSingleInitStmt(stmts[startIdx])
		if !okD || curId == nil || curId.Name == "_" {
			return
		}
		cv := varFromObj(info.Defs[curId])
		if cv == nil {
			return
		}
		ic, ifn, targetExpr, okIC := matchNewIterCursorCall(rhsExpr, info)
		if !okIC {
			return
		}
		var nExpr ast.Expr
		var nType types.Type
		if expectedTargetVar != nil {
			tId, isId := targetExpr.(*ast.Ident)
			if !isId || varFromObj(info.Uses[tId]) != expectedTargetVar {
				return
			}
		} else {
			var okGA bool
			nExpr, nType, okGA = matchNodeGetAncestorsCall(targetExpr, pkg, info)
			if !okGA {
				return
			}
		}
		idx := startIdx + 1
		bUses := 0
		if idx < len(stmts) && isBlankAssignOfVar(stmts[idx], cv, info) {
			bUses = 1
			idx++
		}
		if idx >= len(stmts) {
			return
		}
		fs, isFor := stmts[idx].(*ast.ForStmt)
		if !isFor || fs.Init != nil {
			return
		}
		return fs, cv, bUses, ic, ifn, nExpr, nType, idx + 1, true
	}

	count := 0
	var walkNode func(ast.Node)
	replaceStmtList := func(container ast.Node, list []ast.Stmt) []ast.Stmt {
		for i := 0; i < len(list); i++ {
			// Layout 1: `for cur := j2gort.NewIterCursor(expr.GetAncestors()); cur.HasNext(); { ... }`
			if fs, isFor := list[i].(*ast.ForStmt); isFor && fs.Init != nil {
				if initAs, okAs := fs.Init.(*ast.AssignStmt); okAs && initAs.Tok == token.DEFINE && len(initAs.Lhs) == 1 && len(initAs.Rhs) == 1 {
					if curId, okId := initAs.Lhs[0].(*ast.Ident); okId && curId.Name != "_" {
						if cv := varFromObj(info.Defs[curId]); cv != nil {
							if ic, ifn, targetExpr, okIC := matchNewIterCursorCall(initAs.Rhs[0], info); okIC {
								if nExpr, nType, okGA := matchNodeGetAncestorsCall(targetExpr, pkg, info); okGA {
									walkNode(fs.Body)
									if lowered, okL := buildLoweredAncestorsStmts(fs, cv, 0, ic, ifn, nExpr, nType, container); okL {
										list[i] = &ast.BlockStmt{List: lowered}
										count++
										continue
									}
								}
							}
						}
					}
				}
			}

			// Layout 3: `var ancestors = expr.GetAncestors(); [_ = ancestors;] [pure setup...] { cur := NewIterCursor(ancestors); for ... }` (or inline in same list)
			if ancId, ancRhs, okAnc := matchLocalSingleInitStmt(list[i]); okAnc && ancId != nil && ancId.Name != "_" {
				if ancVar := varFromObj(info.Defs[ancId]); ancVar != nil {
					if nExpr, nType, okGA := matchNodeGetAncestorsCall(ancRhs, pkg, info); okGA {
						scanIdx := i + 1
						ancBlankUses := 0
						if scanIdx < len(list) && isBlankAssignOfVar(list[scanIdx], ancVar, info) {
							ancBlankUses = 1
							scanIdx++
						}
						var keptMid []ast.Stmt
						for scanIdx < len(list) && isPureSetupStmtBeforeIterLoop(list[scanIdx], ancVar, info) {
							keptMid = append(keptMid, list[scanIdx])
							scanIdx++
						}
						if scanIdx < len(list) && varUseCounts[ancVar] == 1+ancBlankUses {
							// Subcase 3a: Inner `{ cur := ...; for ... }` block
							if innerBlk, isBlk := list[scanIdx].(*ast.BlockStmt); isBlk {
								if fs, cv, bUses, ic, ifn, _, _, endInner, okM := matchCursorAndForInSlice(innerBlk.List, 0, ancVar); okM && endInner == len(innerBlk.List) {
									walkNode(fs.Body)
									if lowered, okL := buildLoweredAncestorsStmts(fs, cv, bUses, ic, ifn, nExpr, nType, container); okL {
										repl := append(keptMid, ast.Stmt(&ast.BlockStmt{Lbrace: innerBlk.Lbrace, List: lowered, Rbrace: innerBlk.Rbrace}))
										list = slices.Replace(list, i, scanIdx+1, repl...)
										i += len(keptMid)
										count++
										continue
									}
								}
							}
							// Subcase 3b: Same-block `cur := ...; for ...`
							if fs, cv, bUses, ic, ifn, _, _, endSame, okM := matchCursorAndForInSlice(list, scanIdx, ancVar); okM {
								walkNode(fs.Body)
								if lowered, okL := buildLoweredAncestorsStmts(fs, cv, bUses, ic, ifn, nExpr, nType, container); okL {
									if len(keptMid) == 0 && i == 0 && endSame == len(list) {
										list = lowered
										count++
										return list
									}
									repl := append(keptMid, ast.Stmt(&ast.BlockStmt{List: lowered}))
									list = slices.Replace(list, i, endSame, repl...)
									i += len(keptMid)
									count++
									continue
								}
							}
						}
					}
				}
			}

			// Layout 2: `cur := j2gort.NewIterCursor(expr.GetAncestors()); [_ = cur;] for cur.HasNext() { ... }`
			if fs, cv, bUses, ic, ifn, nExpr, nType, endIdx, okM := matchCursorAndForInSlice(list, i, nil); okM {
				walkNode(fs.Body)
				if lowered, okL := buildLoweredAncestorsStmts(fs, cv, bUses, ic, ifn, nExpr, nType, container); okL {
					if i == 0 && endIdx == len(list) {
						list = lowered
						count++
						return list
					}
					repl := ast.Stmt(&ast.BlockStmt{List: lowered})
					list = slices.Replace(list, i, endIdx, repl)
					count++
					continue
				}
			}

			walkNode(list[i])
		}
		return list
	}

	walkNode = func(n ast.Node) {
		if n == nil {
			return
		}
		switch x := n.(type) {
		case *ast.BlockStmt:
			x.List = replaceStmtList(x, x.List)
		case *ast.CaseClause:
			x.Body = replaceStmtList(x, x.Body)
		default:
			ast.Inspect(x, func(child ast.Node) bool {
				if child == nil || child == x {
					return true
				}
				switch child.(type) {
				case *ast.BlockStmt, *ast.CaseClause:
					walkNode(child)
					return false
				}
				return true
			})
		}
	}

	walkNode(file)
	return count
}

func uniqueNodeLocalName(scopeRoot ast.Node, candidate string) string {
	used := make(map[string]bool)
	if scopeRoot != nil {
		ast.Inspect(scopeRoot, func(n ast.Node) bool {
			if id, ok := n.(*ast.Ident); ok {
				used[id.Name] = true
			}
			return true
		})
	}
	name := candidate
	for used[name] {
		name += "_"
	}
	return name
}

// --- Wave 29C Task 2: Lazy Guard for Preconditions.checkState / checkArgument Format & Message Calls ---

func isLazyPreconditionFuncName(name string) bool {
	switch name {
	case "CheckStateFmt", "CheckStateFmt1", "CheckStateFmt2", "CheckStateFmt3", "CheckStateFmt4",
		"CheckArgumentFmt", "CheckArgumentFmt1", "CheckArgumentFmt2", "CheckArgumentFmt3", "CheckArgumentFmt4",
		"CheckStateMsg", "CheckArgumentMsg":
		return true
	}
	return false
}

func fileHasLazyPreconditionCall(file *ast.File) bool {
	if file == nil {
		return false
	}
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel != nil && isLazyPreconditionFuncName(sel.Sel.Name) {
			found = true
			return false
		}
		return true
	})
	return found
}

func isLiteralBoolExpr(expr ast.Expr, want bool, info *types.Info) bool {
	e := stripParens(expr)
	if id, ok := e.(*ast.Ident); ok {
		if want && id.Name == "true" {
			if info == nil || info.Uses[id] == nil || info.Uses[id] == types.Universe.Lookup("true") {
				return true
			}
		}
		if !want && id.Name == "false" {
			if info == nil || info.Uses[id] == nil || info.Uses[id] == types.Universe.Lookup("false") {
				return true
			}
		}
	}
	if info != nil {
		if tv, ok := info.Types[e]; ok && tv.Value != nil && tv.Value.Kind() == constant.Bool {
			return constant.BoolVal(tv.Value) == want
		}
	}
	return false
}

type lazyPreconditionCondFacts struct {
	nonNilVars      map[*types.Var]bool
	nonZeroNodeVars map[*types.Var]bool
}

func collectPreconditionCondFacts(cond ast.Expr, baseFacts *peepholeFileFacts, info *types.Info) *lazyPreconditionCondFacts {
	cf := &lazyPreconditionCondFacts{
		nonNilVars:      make(map[*types.Var]bool),
		nonZeroNodeVars: make(map[*types.Var]bool),
	}
	if baseFacts != nil {
		// order-independent: copying boolean variable fact sets
		for k, v := range baseFacts.nonNilVars {
			if v {
				cf.nonNilVars[k] = true
			}
		}
		// order-independent: copying boolean variable fact sets
		for k, v := range baseFacts.nonZeroNodeVars {
			if v {
				cf.nonZeroNodeVars[k] = true
			}
		}
	}
	if cond == nil || info == nil {
		return cf
	}

	// 1. Walk unconditionally evaluated expressions in cond: any Node method call (other than IsJavaNil_/RawNodeId_)
	//    or pointer struct field dereference would already have panicked with NPE during cond evaluation.
	var walkUncond func(ast.Expr)
	walkUncond = func(expr ast.Expr) {
		if expr == nil {
			return
		}
		switch e := stripParens(expr).(type) {
		case *ast.UnaryExpr:
			walkUncond(e.X)
		case *ast.BinaryExpr:
			walkUncond(e.X)
			if e.Op != token.LAND && e.Op != token.LOR {
				walkUncond(e.Y)
			}
		case *ast.SelectorExpr:
			walkUncond(e.X)
			if seln, ok := info.Selections[e]; ok && seln.Kind() == types.FieldVal {
				if rId, okId := stripParens(e.X).(*ast.Ident); okId {
					if v := varFromObj(info.Uses[rId]); v != nil {
						if _, isPtr := types.Unalias(v.Type()).(*types.Pointer); isPtr {
							cf.nonNilVars[v] = true
						}
					}
				}
			}
		case *ast.CallExpr:
			if sel, okSel := stripParens(e.Fun).(*ast.SelectorExpr); okSel && sel.Sel != nil {
				walkUncond(sel.X)
				mName := sel.Sel.Name
				if mName != "IsJavaNil_" && mName != "RawNodeId_" {
					if rId, okId := stripParens(sel.X).(*ast.Ident); okId {
						if v := varFromObj(info.Uses[rId]); v != nil && isNamedUint32NodeType(v.Type()) {
							cf.nonZeroNodeVars[v] = true
						}
					}
				}
			} else {
				walkUncond(e.Fun)
			}
			for _, a := range e.Args {
				walkUncond(a)
			}
		}
	}
	walkUncond(cond)

	// 2. Collect facts guaranteed specifically when cond evaluates to true (conjunctions, != 0, != nil, Is*/Has* calls).
	var walkWhenTrue func(ast.Expr)
	walkWhenTrue = func(expr ast.Expr) {
		if expr == nil {
			return
		}
		switch e := stripParens(expr).(type) {
		case *ast.BinaryExpr:
			switch e.Op {
			case token.LAND:
				walkWhenTrue(e.X)
				walkWhenTrue(e.Y)
				walkUncond(e.Y)
			case token.NEQ:
				markCondNeqFact(e.X, e.Y, cf, info)
				markCondNeqFact(e.Y, e.X, cf, info)
			}
		case *ast.UnaryExpr:
			if e.Op == token.NOT {
				if bin, okB := stripParens(e.X).(*ast.BinaryExpr); okB {
					switch bin.Op {
					case token.EQL:
						markCondNeqFact(bin.X, bin.Y, cf, info)
						markCondNeqFact(bin.Y, bin.X, cf, info)
					case token.LOR:
						// !(a || b) == (!a && !b): both a and b are evaluated!
						walkUncond(bin.X)
						walkUncond(bin.Y)
					}
				}
			}
		case *ast.CallExpr:
			if sel, okSel := stripParens(e.Fun).(*ast.SelectorExpr); okSel && sel.Sel != nil {
				mName := sel.Sel.Name
				if mName != "IsJavaNil_" && mName != "RawNodeId_" {
					if rId, okId := stripParens(sel.X).(*ast.Ident); okId {
						if v := varFromObj(info.Uses[rId]); v != nil && isNamedUint32NodeType(v.Type()) {
							cf.nonZeroNodeVars[v] = true
						}
					}
				}
			}
		}
	}
	walkWhenTrue(cond)

	return cf
}

func markCondNeqFact(lhs, rhs ast.Expr, cf *lazyPreconditionCondFacts, info *types.Info) {
	id, ok := stripParens(lhs).(*ast.Ident)
	if !ok {
		return
	}
	v := varFromObj(info.Uses[id])
	if v == nil {
		return
	}
	r := stripParens(rhs)
	if isNamedUint32NodeType(v.Type()) && isZeroConstExpr(r, info) {
		cf.nonZeroNodeVars[v] = true
		return
	}
	if isUntypedNilExpr(r, info) {
		_, isPtr := types.Unalias(v.Type()).(*types.Pointer)
		if isPtr || isInterfaceType(v.Type()) {
			cf.nonNilVars[v] = true
		}
	}
}

func isZeroConstExpr(expr ast.Expr, info *types.Info) bool {
	e := stripParens(expr)
	if lit, ok := e.(*ast.BasicLit); ok && lit.Kind == token.INT && lit.Value == "0" {
		return true
	}
	if info != nil {
		if tv, ok := info.Types[e]; ok && tv.Value != nil && tv.Value.Kind() == constant.Int && constant.Sign(tv.Value) == 0 {
			return true
		}
	}
	return false
}

func isNodeNonZeroInCond(expr ast.Expr, cf *lazyPreconditionCondFacts, baseFacts *peepholeFileFacts, info *types.Info) bool {
	if isGuaranteedNonZeroNode(expr, baseFacts, info) {
		return true
	}
	if cf != nil && info != nil {
		if id, ok := stripParens(expr).(*ast.Ident); ok {
			if v := varFromObj(info.Uses[id]); v != nil && cf.nonZeroNodeVars[v] {
				return true
			}
		}
	}
	return false
}

func isPtrOrValSafeForFieldRead(expr ast.Expr, cf *lazyPreconditionCondFacts, baseFacts *peepholeFileFacts, info *types.Info) bool {
	if expr == nil || info == nil {
		return false
	}
	e := stripParens(expr)
	t := info.Types[e].Type
	if t == nil {
		return false
	}
	if _, isPtr := types.Unalias(t).(*types.Pointer); !isPtr {
		// Value-typed struct: reading a field cannot dereference a nil pointer unless the selector hops through an embedded pointer.
		return true
	}
	if isGuaranteedNonNilForCheckNotNull(e, baseFacts, info) {
		return true
	}
	if cf != nil {
		if id, ok := e.(*ast.Ident); ok {
			if v := varFromObj(info.Uses[id]); v != nil && cf.nonNilVars[v] {
				return true
			}
		}
	}
	return false
}

func selectionHasIntermediatePtrDeref(seln *types.Selection) bool {
	if seln == nil || len(seln.Index()) <= 1 {
		return false
	}
	curT := seln.Recv()
	for _, fIdx := range seln.Index()[:len(seln.Index())-1] {
		ut := types.Unalias(curT)
		if p, ok := ut.(*types.Pointer); ok {
			ut = types.Unalias(p.Elem())
		}
		st, ok := ut.Underlying().(*types.Struct)
		if !ok || fIdx < 0 || fIdx >= st.NumFields() {
			return true
		}
		ft := st.Field(fIdx).Type()
		if _, isPtr := types.Unalias(ft).(*types.Pointer); isPtr {
			return true
		}
		curT = ft
	}
	return false
}

func isSafeLazyPreconditionArg(
	expr ast.Expr,
	cf *lazyPreconditionCondFacts,
	baseFacts *peepholeFileFacts,
	pkg *types.Package,
	info *types.Info,
) bool {
	if expr == nil || info == nil {
		return false
	}
	switch e := stripParens(expr).(type) {
	case *ast.BasicLit:
		return true

	case *ast.Ident:
		switch info.Uses[e].(type) {
		case *types.Var, *types.Const, *types.Nil:
			return true
		}
		return false

	case *ast.CompositeLit:
		compT := info.Types[e].Type
		if compT == nil {
			return false
		}
		switch types.Unalias(compT).Underlying().(type) {
		case *types.Slice, *types.Array, *types.Struct:
		default:
			return false
		}
		for _, elt := range e.Elts {
			val := elt
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				val = kv.Value
			}
			if !isSafeLazyPreconditionArg(val, cf, baseFacts, pkg, info) {
				return false
			}
		}
		return true

	case *ast.UnaryExpr:
		switch e.Op {
		case token.ADD, token.SUB, token.NOT, token.XOR:
			return isSafeLazyPreconditionArg(e.X, cf, baseFacts, pkg, info)
		case token.AND:
			inner := stripParens(e.X)
			if id, ok := inner.(*ast.Ident); ok {
				_, isVar := info.Uses[id].(*types.Var)
				return isVar
			}
			if sel, ok := inner.(*ast.SelectorExpr); ok {
				return isSafeLazyPreconditionArg(sel, cf, baseFacts, pkg, info)
			}
		}
		return false

	case *ast.BinaryExpr:
		switch e.Op {
		case token.ADD, token.SUB, token.MUL,
			token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ,
			token.LAND, token.LOR, token.AND, token.OR, token.XOR, token.SHL, token.SHR:
			return isSafeLazyPreconditionArg(e.X, cf, baseFacts, pkg, info) &&
				isSafeLazyPreconditionArg(e.Y, cf, baseFacts, pkg, info)
		}
		return false

	case *ast.SelectorExpr:
		if pkgId, ok := stripParens(e.X).(*ast.Ident); ok {
			if _, isPkg := info.Uses[pkgId].(*types.PkgName); isPkg {
				switch info.Uses[e.Sel].(type) {
				case *types.Var, *types.Const:
					return true
				}
				return false
			}
		}
		seln, ok := info.Selections[e]
		if !ok || seln.Kind() != types.FieldVal || selectionHasIntermediatePtrDeref(seln) {
			return false
		}
		return isPtrOrValSafeForFieldRead(e.X, cf, baseFacts, info) &&
			isSafeLazyPreconditionArg(e.X, cf, baseFacts, pkg, info)

	case *ast.CallExpr:
		if e.Ellipsis.IsValid() {
			return false
		}
		// 1. Safe scalar / interface type conversion `T(x)`.
		if tv, okT := info.Types[e.Fun]; okT && tv.IsType() && len(e.Args) == 1 {
			switch types.Unalias(tv.Type).Underlying().(type) {
			case *types.Basic, *types.Interface:
				return isSafeLazyPreconditionArg(e.Args[0], cf, baseFacts, pkg, info)
			}
			return false
		}

		// 2. `ToIface_*` interface adapter.
		if _, isTI := isToIfaceCall(e, info); isTI {
			return isSafeLazyPreconditionArg(e.Args[0], cf, baseFacts, pkg, info)
		}

		// 3. Pure boxing / conversion helpers in `j2gort`.
		if fnName, _, _, isJ2 := isJ2gortCall(e, info); isJ2 {
			switch fnName {
			case "BoxString", "BoxNode", "PtrToAny", "SliceToAny", "ToAny", "AsIface",
				"ValueOfInt", "ValueOfLong", "ValueOfBoolean", "ValueOfDouble", "ValueOfFloat",
				"ValueOfByte", "ValueOfShort", "ValueOfChar",
				"BoxNullableVal", "StrVal", "StrValString", "StrValInt32", "StrValInt64", "StrValBool", "StrValChar":
				return len(e.Args) == 1 && isSafeLazyPreconditionArg(e.Args[0], cf, baseFacts, pkg, info)
			case "StrConcat":
				for _, a := range e.Args {
					if !isSafeLazyPreconditionArg(a, cf, baseFacts, pkg, info) {
						return false
					}
				}
				return true
			case "CheckNotNullNode":
				return len(e.Args) == 1 && isNodeNonZeroInCond(e.Args[0], cf, baseFacts, info) &&
					isSafeLazyPreconditionArg(e.Args[0], cf, baseFacts, pkg, info)
			case "CheckNotNullString":
				return len(e.Args) == 1 && isGuaranteedNonNullString(e.Args[0], info) &&
					isSafeLazyPreconditionArg(e.Args[0], cf, baseFacts, pkg, info)
			case "CheckNotNullPtr":
				return len(e.Args) == 1 && isPtrOrValSafeForFieldRead(e.Args[0], cf, baseFacts, info) &&
					isSafeLazyPreconditionArg(e.Args[0], cf, baseFacts, pkg, info)
			}
			return false
		}

		// 4. Pure enum `.String()` / `.ToString()` / `.Name()` on `x.GetToken()` (with optional `CheckNotNullEnum`),
		//    and pure `Node` inspection methods on `x` when `x` is proven non-zero when `cond` is true.
		if sel, okSel := stripParens(e.Fun).(*ast.SelectorExpr); okSel && sel.Sel != nil && len(e.Args) == 0 {
			mName := sel.Sel.Name
			switch mName {
			case "String", "ToString", "Name":
				enumRecv := stripParens(sel.X)
				unwrappedCheckNotNullEnum := false
				if cCall, okC := enumRecv.(*ast.CallExpr); okC && len(cCall.Args) == 1 && !cCall.Ellipsis.IsValid() {
					if cFn, _, _, isCJ2 := isJ2gortCall(cCall, info); isCJ2 && cFn == "CheckNotNullEnum" {
						enumRecv = stripParens(cCall.Args[0])
						unwrappedCheckNotNullEnum = true
					}
				}
				if tokCall, okTC := enumRecv.(*ast.CallExpr); okTC && len(tokCall.Args) == 0 && !tokCall.Ellipsis.IsValid() {
					if tokSel, okTS := stripParens(tokCall.Fun).(*ast.SelectorExpr); okTS && tokSel.Sel != nil && tokSel.Sel.Name == "GetToken" {
						nodeX := stripParens(tokSel.X)
						if isNamedUint32NodeType(info.Types[nodeX].Type) &&
							isNodeNonZeroInCond(nodeX, cf, baseFacts, info) &&
							isSafeLazyPreconditionArg(nodeX, cf, baseFacts, pkg, info) {
							return true
						}
					}
				}
				// Also allow `x.ToString()` directly on a Node `x` proven non-zero when `cond` is true!
				if !unwrappedCheckNotNullEnum && (mName == "ToString" || mName == "String") &&
					isNamedUint32NodeType(info.Types[enumRecv].Type) &&
					isNodeNonZeroInCond(enumRecv, cf, baseFacts, info) &&
					isSafeLazyPreconditionArg(enumRecv, cf, baseFacts, pkg, info) {
					return true
				}

			case "ToString_0", "ToStringTree", "GetToken", "GetParent", "GetFirstChild", "GetLastChild", "GetNext", "GetLineno", "GetCharno", "GetChildCount":
				nodeX := stripParens(sel.X)
				if isNamedUint32NodeType(info.Types[nodeX].Type) &&
					isNodeNonZeroInCond(nodeX, cf, baseFacts, info) &&
					isSafeLazyPreconditionArg(nodeX, cf, baseFacts, pkg, info) {
					return true
				}
			}
		}
	}
	return false
}

func (p *J2gortPeepholePass) rewriteLazyPreconditions(
	file *ast.File,
	pkg *types.Package,
	info *types.Info,
) int {
	if !fileHasLazyPreconditionCall(file) || info == nil {
		return 0
	}
	parents, _ := buildParentMap(file)
	facts := collectPeepholeFileFacts(file, info, parents)
	count := 0

	tryRewriteStmt := func(stmt ast.Stmt) (ast.Stmt, bool) {
		es, isES := stmt.(*ast.ExprStmt)
		if !isES {
			return nil, false
		}
		call, isCall := stripParens(es.X).(*ast.CallExpr)
		if !isCall || len(call.Args) < 2 {
			return nil, false
		}
		fnName, _, _, isJ2 := isJ2gortCall(call, info)
		if !isJ2 || !isLazyPreconditionFuncName(fnName) {
			return nil, false
		}
		cond := stripParens(call.Args[0])
		// Skip literal false (100% Run-2 idempotency!) and literal true.
		if isLiteralBoolExpr(cond, false, info) || isLiteralBoolExpr(cond, true, info) {
			return nil, false
		}

		condFacts := collectPreconditionCondFacts(cond, facts, info)
		for _, arg := range call.Args[1:] {
			if !isSafeLazyPreconditionArg(arg, condFacts, facts, pkg, info) {
				return nil, false
			}
		}

		pos := call.Pos()
		falseIdent := &ast.Ident{NamePos: call.Args[0].Pos(), Name: "false"}
		if uFalse := types.Universe.Lookup("false"); uFalse != nil {
			info.Uses[falseIdent] = uFalse
		}
		info.Types[falseIdent] = types.TypeAndValue{
			Type:  types.Typ[types.Bool],
			Value: constant.MakeBool(false),
		}
		call.Args[0] = falseIdent

		parenCond := &ast.ParenExpr{Lparen: pos, X: cond, Rparen: pos}
		info.Types[parenCond] = types.TypeAndValue{Type: types.Typ[types.Bool]}
		negCond := &ast.UnaryExpr{OpPos: pos, Op: token.NOT, X: parenCond}
		info.Types[negCond] = types.TypeAndValue{Type: types.Typ[types.Bool]}

		ifStmt := &ast.IfStmt{
			If:   pos,
			Cond: negCond,
			Body: &ast.BlockStmt{
				Lbrace: pos,
				List:   []ast.Stmt{es},
				Rbrace: call.End(),
			},
		}
		return ifStmt, true
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BlockStmt:
			for i, st := range x.List {
				if repl, ok := tryRewriteStmt(st); ok {
					x.List[i] = repl
					count++
				}
			}
		case *ast.CaseClause:
			for i, st := range x.Body {
				if repl, ok := tryRewriteStmt(st); ok {
					x.Body[i] = repl
					count++
				}
			}
		}
		return true
	})

	return count
}

// --- Wave 29C Task 3: Lazy IsLogging() Guard on diagnostic.LogFile Statements ---

func isLogFileLogMethodName(name string) bool {
	return name == "Log" || name == "LogJson" ||
		strings.HasPrefix(name, "Log_") || strings.HasPrefix(name, "LogJson_")
}

func fileHasLogFileMethodCall(file *ast.File) bool {
	if file == nil {
		return false
	}
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel != nil && isLogFileLogMethodName(sel.Sel.Name) {
			found = true
			return false
		}
		return true
	})
	return found
}

func hasZeroArgBoolMethodOrAddr(t types.Type, pkg *types.Package, name string) bool {
	if t == nil {
		return false
	}
	if hasZeroArgBoolMethod(t, pkg, name) {
		return true
	}
	obj, _, _ := types.LookupFieldOrMethod(t, true, pkg, name)
	fn, ok := obj.(*types.Func)
	if !ok {
		return false
	}
	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Params().Len() != 0 || sig.Results().Len() != 1 {
		return false
	}
	b, ok := types.Unalias(sig.Results().At(0).Type()).(*types.Basic)
	return ok && b.Kind() == types.Bool
}

func isDiagnosticLogFileType(t types.Type, pkg *types.Package) bool {
	if t == nil || !hasZeroArgBoolMethodOrAddr(t, pkg, "IsLogging") {
		return false
	}
	named := derefNamedOrigin(t)
	if named != nil && named.Obj() != nil && named.Obj().Pkg() != nil {
		p := named.Obj().Pkg()
		tName := named.Obj().Name()
		if tName == "CompilerDecisionsLog" {
			return true
		}
		if p.Name() == "diagnostic" || strings.HasSuffix(p.Path(), "/diagnostic") {
			if tName == "LogFile" || hasZeroArgVoidMethod(t, pkg, "Implements_diagnostic_LogFile") {
				return true
			}
		}
	}
	return hasZeroArgVoidMethod(t, pkg, "Implements_diagnostic_LogFile") ||
		hasZeroArgVoidMethod(t, pkg, "Implements_jscomp_CompilerDecisionsLog")
}

func isBaseGetterMethodName(name string) bool {
	return len(name) > len("GetBase") && strings.HasPrefix(name, "Get") && strings.HasSuffix(name, "Base")
}

func isPureLogFileRecvExpr(expr ast.Expr, info *types.Info) bool {
	if expr == nil || info == nil {
		return false
	}
	switch e := stripParens(expr).(type) {
	case *ast.Ident:
		v, ok := info.Uses[e].(*types.Var)
		return ok && v != nil
	case *ast.SelectorExpr:
		if pkgId, ok := stripParens(e.X).(*ast.Ident); ok {
			if _, isPkg := info.Uses[pkgId].(*types.PkgName); isPkg {
				v, okV := info.Uses[e.Sel].(*types.Var)
				return okV && v != nil
			}
		}
		seln, ok := info.Selections[e]
		return ok && seln.Kind() == types.FieldVal && isPureLogFileRecvExpr(e.X, info)
	case *ast.CallExpr:
		if len(e.Args) != 0 || e.Ellipsis.IsValid() {
			return false
		}
		sel, ok := stripParens(e.Fun).(*ast.SelectorExpr)
		return ok && sel.Sel != nil && isBaseGetterMethodName(sel.Sel.Name) && isPureLogFileRecvExpr(sel.X, info)
	}
	return false
}

func clonePureLogFileRecvExpr(expr ast.Expr) ast.Expr {
	switch e := stripParens(expr).(type) {
	case *ast.Ident:
		return &ast.Ident{Name: e.Name}
	case *ast.SelectorExpr:
		return &ast.SelectorExpr{
			X:   clonePureLogFileRecvExpr(e.X),
			Sel: &ast.Ident{Name: e.Sel.Name},
		}
	case *ast.CallExpr:
		if sel, ok := stripParens(e.Fun).(*ast.SelectorExpr); ok && sel.Sel != nil {
			return &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   clonePureLogFileRecvExpr(sel.X),
					Sel: &ast.Ident{Name: sel.Sel.Name},
				},
			}
		}
	}
	return cloneSimplePureExpr(expr)
}

func sameSimplePureExpr(a, b ast.Expr, info *types.Info) bool {
	if a == nil || b == nil || info == nil {
		return false
	}
	ea := stripParens(a)
	eb := stripParens(b)
	switch xa := ea.(type) {
	case *ast.Ident:
		xb, ok := eb.(*ast.Ident)
		if !ok {
			return false
		}
		ua, ub := info.Uses[xa], info.Uses[xb]
		if ua != nil && ub != nil {
			return ua == ub
		}
		return xa.Name == xb.Name
	case *ast.SelectorExpr:
		xb, ok := eb.(*ast.SelectorExpr)
		if !ok || xa.Sel == nil || xb.Sel == nil || xa.Sel.Name != xb.Sel.Name {
			return false
		}
		return sameSimplePureExpr(xa.X, xb.X, info)
	case *ast.CallExpr:
		xb, ok := eb.(*ast.CallExpr)
		if !ok || len(xa.Args) != 0 || len(xb.Args) != 0 || xa.Ellipsis.IsValid() || xb.Ellipsis.IsValid() {
			return false
		}
		sa, okA := stripParens(xa.Fun).(*ast.SelectorExpr)
		sb, okB := stripParens(xb.Fun).(*ast.SelectorExpr)
		if !okA || !okB || sa.Sel == nil || sb.Sel == nil || !isBaseGetterMethodName(sa.Sel.Name) || sa.Sel.Name != sb.Sel.Name {
			return false
		}
		return sameSimplePureExpr(sa.X, sb.X, info)
	}
	return false
}

func condHasPositiveIsLoggingForRecv(cond, recv ast.Expr, info *types.Info) bool {
	if cond == nil {
		return false
	}
	switch e := stripParens(cond).(type) {
	case *ast.BinaryExpr:
		if e.Op == token.LAND {
			return condHasPositiveIsLoggingForRecv(e.X, recv, info) ||
				condHasPositiveIsLoggingForRecv(e.Y, recv, info)
		}
	case *ast.CallExpr:
		if len(e.Args) == 0 && !e.Ellipsis.IsValid() {
			if sel, ok := stripParens(e.Fun).(*ast.SelectorExpr); ok && sel.Sel != nil && sel.Sel.Name == "IsLogging" {
				return sameSimplePureExpr(sel.X, recv, info)
			}
		}
	}
	return false
}

func isAlreadyGuardedByIsLogging(stmt ast.Stmt, recv ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) bool {
	cur := ast.Node(stmt)
	for cur != nil {
		p := parents[cur]
		if p == nil {
			return false
		}
		switch par := p.(type) {
		case *ast.FuncDecl, *ast.FuncLit:
			return false
		case *ast.IfStmt:
			if par.Body == cur && condHasPositiveIsLoggingForRecv(par.Cond, recv, info) {
				return true
			}
		}
		cur = p
	}
	return false
}

func isSafeBoundLambdaIIFE(call *ast.CallExpr, info *types.Info) bool {
	if len(call.Args) != 0 || call.Ellipsis.IsValid() {
		return false
	}
	fl, ok := stripParens(call.Fun).(*ast.FuncLit)
	if !ok || fl.Body == nil || len(fl.Body.List) == 0 {
		return false
	}
	for i := 0; i < len(fl.Body.List)-1; i++ {
		_, rhs, okD := matchLocalSingleInitStmt(fl.Body.List[i])
		if !okD || !isPureLogFileRecvExpr(rhs, info) {
			return false
		}
	}
	ret, okR := fl.Body.List[len(fl.Body.List)-1].(*ast.ReturnStmt)
	if !okR || len(ret.Results) != 1 {
		return false
	}
	retCall, okC := stripParens(ret.Results[0]).(*ast.CallExpr)
	if !okC || retCall.Ellipsis.IsValid() {
		return false
	}
	if fnName, _, _, isJ2 := isJ2gortCall(retCall, info); isJ2 && isLambdaFactoryFuncName(fnName) {
		return true
	}
	if callee, _ := calleeFunc(retCall, info); callee != nil && isLambdaFactoryFuncName(callee.Name()) {
		return true
	}
	return false
}

func isPureLogArgExpr(expr ast.Expr, info *types.Info) bool {
	if expr == nil || info == nil {
		return false
	}
	switch e := stripParens(expr).(type) {
	case *ast.BasicLit:
		return true

	case *ast.Ident:
		switch info.Uses[e].(type) {
		case *types.Var, *types.Const, *types.Nil:
			return true
		}
		return false

	case *ast.SelectorExpr:
		if pkgId, ok := stripParens(e.X).(*ast.Ident); ok {
			if _, isPkg := info.Uses[pkgId].(*types.PkgName); isPkg {
				switch info.Uses[e.Sel].(type) {
				case *types.Var, *types.Const:
					return true
				}
				return false
			}
		}
		seln, ok := info.Selections[e]
		return ok && seln.Kind() == types.FieldVal && isPureLogArgExpr(e.X, info)

	case *ast.IndexExpr:
		return isPureLogArgExpr(e.X, info) && isPureLogArgExpr(e.Index, info)

	case *ast.CompositeLit:
		compT := info.Types[e].Type
		if compT == nil {
			return false
		}
		switch types.Unalias(compT).Underlying().(type) {
		case *types.Slice, *types.Array, *types.Struct:
		default:
			return false
		}
		for _, elt := range e.Elts {
			val := elt
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				val = kv.Value
			}
			if !isPureLogArgExpr(val, info) {
				return false
			}
		}
		return true

	case *ast.UnaryExpr:
		switch e.Op {
		case token.ADD, token.SUB, token.NOT, token.XOR, token.AND:
			return isPureLogArgExpr(e.X, info)
		}
		return false

	case *ast.BinaryExpr:
		switch e.Op {
		case token.ADD, token.SUB, token.MUL,
			token.AND, token.OR, token.XOR, token.SHL, token.SHR, token.AND_NOT,
			token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ,
			token.LAND, token.LOR:
			return isPureLogArgExpr(e.X, info) && isPureLogArgExpr(e.Y, info)
		}
		return false

	case *ast.CallExpr:
		if e.Ellipsis.IsValid() {
			return false
		}
		if isSafeBoundLambdaIIFE(e, info) {
			return true
		}
		if tv, okT := info.Types[e.Fun]; okT && tv.IsType() && len(e.Args) == 1 {
			switch types.Unalias(tv.Type).Underlying().(type) {
			case *types.Basic, *types.Interface, *types.Pointer, *types.Slice:
				return isPureLogArgExpr(e.Args[0], info)
			}
			return false
		}
		if _, isTI := isToIfaceCall(e, info); isTI {
			return isPureLogArgExpr(e.Args[0], info)
		}
		if _, isUp := isUpcastCall(e, info); isUp {
			return isPureLogArgExpr(e.Args[0], info)
		}
		allArgsPure := func() bool {
			for _, a := range e.Args {
				if !isPureLogArgExpr(a, info) {
					return false
				}
			}
			return true
		}
		lambdaArgsPure := func() bool {
			for _, a := range e.Args {
				if _, isFL := stripParens(a).(*ast.FuncLit); isFL {
					continue
				}
				if !isPureLogArgExpr(a, info) {
					return false
				}
			}
			return true
		}
		if fnName, _, _, isJ2 := isJ2gortCall(e, info); isJ2 {
			if isLambdaFactoryFuncName(fnName) {
				return lambdaArgsPure()
			}
			if strings.HasPrefix(fnName, "Box") ||
				strings.HasPrefix(fnName, "ValueOf") ||
				fnName == "PtrToAny" || fnName == "SliceToAny" || fnName == "ToAny" ||
				strings.HasPrefix(fnName, "AsIface") ||
				strings.HasPrefix(fnName, "Cast") || fnName == "UncheckedCast" ||
				strings.HasPrefix(fnName, "CheckNotNull") ||
				strings.HasPrefix(fnName, "StrVal") ||
				strings.HasPrefix(fnName, "ToIface_") ||
				strings.HasPrefix(fnName, "Upcast_") ||
				strings.HasPrefix(fnName, "Ext_java_lang_String_Format") ||
				strings.HasPrefix(fnName, "Ext_java_lang_String_S_Format") ||
				strings.HasPrefix(fnName, "StringFormat") ||
				fnName == "LenientFormat" || fnName == "StrConcat" {
				return allArgsPure()
			}
			return false
		}

		// Package-level formatting / adapter / record / lambda constructors.
		if callee, _ := calleeFunc(e, info); callee != nil {
			if sig, okS := callee.Type().(*types.Signature); okS && sig.Recv() == nil {
				cName := callee.Name()
				if isLambdaFactoryFuncName(cName) {
					return lambdaArgsPure()
				}
				if (strings.HasPrefix(cName, "ToIface_") || strings.HasPrefix(cName, "Upcast_")) && len(e.Args) == 1 {
					return isPureLogArgExpr(e.Args[0], info)
				}
				if strings.HasPrefix(cName, "StaticGet_") && len(e.Args) == 0 && sig.Results() != nil && sig.Results().Len() == 1 {
					return true
				}
				if strings.HasPrefix(cName, "DotFormatter_ToDot") ||
					strings.Contains(cName, "RemovalLogRecord_For") ||
					cName == "NewStripCode_StripCodeConfigRecord" ||
					strings.HasPrefix(cName, "NewSerializeTypesToPointers_Stream") ||
					(callee.Pkg() != nil && callee.Pkg().Name() == "javalang" && strings.HasPrefix(cName, "StringFormat")) {
					return allArgsPure()
				}
			}
		}

		// Pure formatting / graphviz / string / base-getter methods on pure receivers.
		if sel, okSel := stripParens(e.Fun).(*ast.SelectorExpr); okSel && sel.Sel != nil {
			mName := sel.Sel.Name
			zeroArgPure := mName == "ToGraphvizGraph" ||
				mName == "GetLocation" ||
				mName == "GetFullName" ||
				mName == "GetString" ||
				mName == "GetQualifiedName" ||
				mName == "String" ||
				isBaseGetterMethodName(mName)
			if (zeroArgPure && len(e.Args) == 0) || strings.HasPrefix(mName, "ToString") {
				return isPureLogArgExpr(sel.X, info) && allArgsPure()
			}
		}
	}
	return false
}

func isSafeLazyLogFileArg(expr ast.Expr, info *types.Info) bool {
	return isPureLogArgExpr(expr, info)
}

func extractPureRecvSelectorChain(expr ast.Expr, recvVar *types.Var, info *types.Info) ([]string, bool) {
	var revFields []string
	cur := stripParens(expr)
	for {
		sel, ok := cur.(*ast.SelectorExpr)
		if !ok || sel.Sel == nil {
			break
		}
		if s, okS := info.Selections[sel]; okS && s.Kind() != types.FieldVal {
			return nil, false
		}
		revFields = append(revFields, sel.Sel.Name)
		cur = stripParens(sel.X)
	}
	if len(revFields) == 0 {
		return nil, false
	}
	id, ok := cur.(*ast.Ident)
	if !ok || info.Uses[id] != recvVar {
		return nil, false
	}
	for l, r := 0, len(revFields)-1; l < r; l, r = l+1, r-1 {
		revFields[l], revFields[r] = revFields[r], revFields[l]
	}
	return revFields, true
}

type privLogHelperSpec struct {
	fieldChain  []string
	hasNilGuard bool
	hasJavaNil  bool
}

func matchPrivLogHelperDecl(fn *ast.FuncDecl, pkg *types.Package, info *types.Info) (*types.Func, privLogHelperSpec, bool) {
	var zero privLogHelperSpec
	if fn == nil || fn.Recv == nil || fn.Name == nil || fn.Body == nil || !strings.HasPrefix(fn.Name.Name, "priv_") {
		return nil, zero, false
	}
	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
		return nil, zero, false
	}
	if len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 {
		return nil, zero, false
	}
	rId := fn.Recv.List[0].Names[0]
	if rId == nil || rId.Name == "_" {
		return nil, zero, false
	}
	fnObj, okF := info.Defs[fn.Name].(*types.Func)
	recvVar, okV := info.Defs[rId].(*types.Var)
	if !okF || fnObj == nil || !okV || recvVar == nil {
		return nil, zero, false
	}
	idx := 0
	for idx < len(fn.Body.List) {
		as, isA := fn.Body.List[idx].(*ast.AssignStmt)
		if !isA || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
			break
		}
		bId, isB := stripParens(as.Lhs[0]).(*ast.Ident)
		_, isR := stripParens(as.Rhs[0]).(*ast.Ident)
		if !isB || bId.Name != "_" || !isR {
			break
		}
		idx++
	}
	if idx != len(fn.Body.List)-1 {
		return nil, zero, false
	}
	ifSt, isIf := fn.Body.List[idx].(*ast.IfStmt)
	if !isIf || ifSt.Init != nil || ifSt.Else != nil || ifSt.Body == nil {
		return nil, zero, false
	}

	matchIsLoggingCall := func(e ast.Expr) (ast.Expr, []string, bool) {
		c, isC := stripParens(e).(*ast.CallExpr)
		if !isC || len(c.Args) != 0 || c.Ellipsis.IsValid() {
			return nil, nil, false
		}
		sel, isS := stripParens(c.Fun).(*ast.SelectorExpr)
		if !isS || sel.Sel == nil || sel.Sel.Name != "IsLogging" {
			return nil, nil, false
		}
		logExpr := stripParens(sel.X)
		if !isDiagnosticLogFileType(info.Types[logExpr].Type, pkg) {
			return nil, nil, false
		}
		chain, okC := extractPureRecvSelectorChain(logExpr, recvVar, info)
		if !okC {
			return nil, nil, false
		}
		return logExpr, chain, true
	}

	if _, chain, okDirect := matchIsLoggingCall(ifSt.Cond); okDirect {
		return funcOrigin(fnObj), privLogHelperSpec{fieldChain: chain}, true
	}

	bin, isBin := stripParens(ifSt.Cond).(*ast.BinaryExpr)
	if !isBin || bin.Op != token.LAND {
		return nil, zero, false
	}
	logExpr, chain, okRight := matchIsLoggingCall(bin.Y)
	if !okRight {
		return nil, zero, false
	}

	matchNeqNil := func(e ast.Expr) bool {
		b, okB := stripParens(e).(*ast.BinaryExpr)
		if !okB || b.Op != token.NEQ {
			return false
		}
		if nId, okN := stripParens(b.Y).(*ast.Ident); okN && nId.Name == "nil" {
			return sameSimplePureExpr(b.X, logExpr, info)
		}
		return false
	}
	matchNotJavaNil := func(e ast.Expr) bool {
		u, okU := stripParens(e).(*ast.UnaryExpr)
		if !okU || u.Op != token.NOT {
			return false
		}
		c, okC := stripParens(u.X).(*ast.CallExpr)
		if !okC || len(c.Args) != 0 || c.Ellipsis.IsValid() {
			return false
		}
		s, okS := stripParens(c.Fun).(*ast.SelectorExpr)
		if !okS || s.Sel == nil || s.Sel.Name != "IsJavaNil_" {
			return false
		}
		return sameSimplePureExpr(s.X, logExpr, info)
	}

	if matchNeqNil(bin.X) {
		return funcOrigin(fnObj), privLogHelperSpec{fieldChain: chain, hasNilGuard: true}, true
	}
	if leftAnd, okLA := stripParens(bin.X).(*ast.BinaryExpr); okLA && leftAnd.Op == token.LAND {
		if matchNeqNil(leftAnd.X) && matchNotJavaNil(leftAnd.Y) {
			return funcOrigin(fnObj), privLogHelperSpec{fieldChain: chain, hasNilGuard: true, hasJavaNil: true}, true
		}
	}
	return nil, zero, false
}

func argsContainSupplierLambda(args []ast.Expr, info *types.Info) bool {
	found := false
	for _, arg := range args {
		ast.Inspect(arg, func(n ast.Node) bool {
			if found || n == nil {
				return false
			}
			c, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if isSafeBoundLambdaIIFE(c, info) {
				found = true
				return false
			}
			if fnName, _, _, isJ2 := isJ2gortCall(c, info); isJ2 && isLambdaFactoryFuncName(fnName) {
				found = true
				return false
			}
			if callee, _ := calleeFunc(c, info); callee != nil && isLambdaFactoryFuncName(callee.Name()) {
				found = true
				return false
			}
			return true
		})
		if found {
			return true
		}
	}
	return false
}

func buildLogFieldExprFromRecv(baseRecv ast.Expr, fieldChain []string, pos token.Pos) ast.Expr {
	cur := clonePureLogFileRecvExpr(baseRecv)
	for _, f := range fieldChain {
		cur = &ast.SelectorExpr{
			X:   cur,
			Sel: &ast.Ident{NamePos: pos, Name: f},
		}
	}
	return cur
}

func (p *J2gortPeepholePass) rewriteLazyLogFileStmts(
	file *ast.File,
	pkg *types.Package,
	info *types.Info,
) int {
	if !fileHasLogFileMethodCall(file) || info == nil {
		return 0
	}
	parents, _ := buildParentMap(file)
	count := 0

	var privHelpers map[*types.Func]privLogHelperSpec

	tryWrapLogStmt := func(stmt ast.Stmt) (ast.Stmt, bool) {
		es, isES := stmt.(*ast.ExprStmt)
		if !isES {
			return nil, false
		}
		call, isCall := stripParens(es.X).(*ast.CallExpr)
		if !isCall || call.Ellipsis.IsValid() || len(call.Args) == 0 {
			return nil, false
		}
		sel, isSel := stripParens(call.Fun).(*ast.SelectorExpr)
		if !isSel || sel.Sel == nil {
			return nil, false
		}

		pos := es.Pos()
		if isLogFileLogMethodName(sel.Sel.Name) {
			recv := stripParens(sel.X)
			if !isPureLogFileRecvExpr(recv, info) || !isDiagnosticLogFileType(info.Types[recv].Type, pkg) {
				return nil, false
			}
			if isAlreadyGuardedByIsLogging(es, recv, parents, info) {
				return nil, false
			}
			if len(call.Args) == 1 {
				switch stripParens(call.Args[0]).(type) {
				case *ast.BasicLit, *ast.Ident:
					return nil, false
				}
			}
			for _, arg := range call.Args {
				if !isPureLogArgExpr(arg, info) {
					return nil, false
				}
			}

			ifStmt := &ast.IfStmt{
				If: pos,
				Cond: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   clonePureLogFileRecvExpr(recv),
						Sel: &ast.Ident{NamePos: pos, Name: "IsLogging"},
					},
				},
				Body: &ast.BlockStmt{
					Lbrace: pos,
					List:   []ast.Stmt{es},
					Rbrace: es.End(),
				},
			}
			parents[ifStmt.Body] = ifStmt
			parents[es] = ifStmt.Body
			return ifStmt, true
		}

		if len(privHelpers) > 0 && strings.HasPrefix(sel.Sel.Name, "priv_") {
			callee, _ := calleeFunc(call, info)
			if callee == nil {
				return nil, false
			}
			spec, okH := privHelpers[funcOrigin(callee)]
			if !okH {
				return nil, false
			}
			callerRecv := stripParens(sel.X)
			if !isPureLogFileRecvExpr(callerRecv, info) || !argsContainSupplierLambda(call.Args, info) {
				return nil, false
			}
			logRecvExpr := buildLogFieldExprFromRecv(callerRecv, spec.fieldChain, pos)
			if isAlreadyGuardedByIsLogging(es, logRecvExpr, parents, info) {
				return nil, false
			}
			for _, arg := range call.Args {
				if !isSafeLazyLogFileArg(arg, info) {
					return nil, false
				}
			}

			isLoggingCall := ast.Expr(&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   buildLogFieldExprFromRecv(callerRecv, spec.fieldChain, pos),
					Sel: &ast.Ident{NamePos: pos, Name: "IsLogging"},
				},
			})
			guardCond := isLoggingCall
			if spec.hasNilGuard {
				neqNil := ast.Expr(&ast.BinaryExpr{
					X:  buildLogFieldExprFromRecv(callerRecv, spec.fieldChain, pos),
					Op: token.NEQ,
					Y:  &ast.Ident{NamePos: pos, Name: "nil"},
				})
				lhs := neqNil
				if spec.hasJavaNil {
					notJavaNil := &ast.UnaryExpr{
						Op: token.NOT,
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   buildLogFieldExprFromRecv(callerRecv, spec.fieldChain, pos),
								Sel: &ast.Ident{NamePos: pos, Name: "IsJavaNil_"},
							},
						},
					}
					lhs = &ast.ParenExpr{
						X: &ast.BinaryExpr{
							X:  neqNil,
							Op: token.LAND,
							Y:  notJavaNil,
						},
					}
				}
				guardCond = &ast.BinaryExpr{
					X:  lhs,
					Op: token.LAND,
					Y:  isLoggingCall,
				}
			}

			ifStmt := &ast.IfStmt{
				If:   pos,
				Cond: guardCond,
				Body: &ast.BlockStmt{
					Lbrace: pos,
					List:   []ast.Stmt{es},
					Rbrace: es.End(),
				},
			}
			parents[ifStmt.Body] = ifStmt
			parents[es] = ifStmt.Body
			return ifStmt, true
		}

		return nil, false
	}

	var walkLogNode func(ast.Node)
	walkLogNode = func(n ast.Node) {
		if n == nil {
			return
		}
		switch x := n.(type) {
		case *ast.BlockStmt:
			for i, st := range x.List {
				walkLogNode(st)
				if repl, ok := tryWrapLogStmt(st); ok {
					x.List[i] = repl
					parents[repl] = x
					count++
				}
			}
		case *ast.CaseClause:
			for i, st := range x.Body {
				walkLogNode(st)
				if repl, ok := tryWrapLogStmt(st); ok {
					x.Body[i] = repl
					parents[repl] = x
					count++
				}
			}
		default:
			ast.Inspect(x, func(child ast.Node) bool {
				if child == nil || child == x {
					return true
				}
				switch child.(type) {
				case *ast.BlockStmt, *ast.CaseClause:
					walkLogNode(child)
					return false
				}
				return true
			})
		}
	}

	// Pass 1: wrap direct <recv>.Log*(...) statements (including inside priv_* helpers).
	walkLogNode(file)

	// Pass 2: discover 1-hop IsLogging()-guarded priv_* void helpers and wrap unguarded Supplier call sites.
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fnOrig, spec, okM := matchPrivLogHelperDecl(fn, pkg, info); okM && fnOrig != nil {
			if privHelpers == nil {
				privHelpers = make(map[*types.Func]privLogHelperSpec)
			}
			privHelpers[fnOrig] = spec
		}
	}
	if len(privHelpers) > 0 {
		walkLogNode(file)
	}

	return count
}

// rewriteIfCondIIFEs hoists 0-arg single-return closure IIFEs out of `if` conditions
// (such as `if func() int32 { ... }() > int32(20000)` in generic methods like
// `DataFlowAnalysisBase[N, L].Analyze()`) into hygienic lexical blocks so Go's
// compiler never emits a heap/indirect closure call for `if` conditions.
func (p *J2gortPeepholePass) rewriteIfCondIIFEs(ctx *PassContext) int {
	if ctx == nil || ctx.Fset == nil {
		return 0
	}
	total := 0
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		if n := hoistFileIfCondIIFEs(file, ctx.Fset); n > 0 {
			ctx.MarkFileModified(file)
			total += n
		}
	}
	return total
}

func hoistFileIfCondIIFEs(file *ast.File, fset *token.FileSet) int {
	if file == nil || fset == nil || !fileHasIfAndFuncLit(file) {
		return 0
	}
	count := 0
	var walk func(ast.Node)
	walk = func(n ast.Node) {
		if n == nil {
			return
		}
		switch x := n.(type) {
		case *ast.BlockStmt:
			for i := 0; i < len(x.List); i++ {
				walk(x.List[i])
				if ifSt, ok := x.List[i].(*ast.IfStmt); ok {
					if replBlk, okH := tryHoistIfCondIIFE(ifSt, fset); okH {
						x.List[i] = replBlk
						count++
					}
				}
			}
		case *ast.CaseClause:
			for i := 0; i < len(x.Body); i++ {
				walk(x.Body[i])
				if ifSt, ok := x.Body[i].(*ast.IfStmt); ok {
					if replBlk, okH := tryHoistIfCondIIFE(ifSt, fset); okH {
						x.Body[i] = replBlk
						count++
					}
				}
			}
		case *ast.CommClause:
			for i := 0; i < len(x.Body); i++ {
				walk(x.Body[i])
				if ifSt, ok := x.Body[i].(*ast.IfStmt); ok {
					if replBlk, okH := tryHoistIfCondIIFE(ifSt, fset); okH {
						x.Body[i] = replBlk
						count++
					}
				}
			}
		case *ast.IfStmt:
			walk(x.Body)
			if x.Else != nil {
				walk(x.Else)
				if elseIf, ok := x.Else.(*ast.IfStmt); ok {
					if replBlk, okH := tryHoistIfCondIIFE(elseIf, fset); okH {
						x.Else = replBlk
						count++
					}
				}
			}
		default:
			ast.Inspect(x, func(child ast.Node) bool {
				if child == nil || child == x {
					return true
				}
				switch child.(type) {
				case *ast.BlockStmt, *ast.CaseClause, *ast.CommClause, *ast.IfStmt:
					walk(child)
					return false
				}
				return true
			})
		}
	}
	walk(file)
	return count
}

func fileHasIfAndFuncLit(file *ast.File) bool {
	hasIf, hasFL := false, false
	ast.Inspect(file, func(n ast.Node) bool {
		if hasIf && hasFL {
			return false
		}
		switch n.(type) {
		case *ast.IfStmt:
			hasIf = true
		case *ast.FuncLit:
			hasFL = true
		}
		return true
	})
	return hasIf && hasFL
}

func tryHoistIfCondIIFE(ifSt *ast.IfStmt, fset *token.FileSet) (*ast.BlockStmt, bool) {
	if ifSt == nil || ifSt.Init != nil || ifSt.Cond == nil || ifSt.Body == nil {
		return nil, false
	}
	call, fl, replaceSlot, ok := matchLeftmostIfCondIIFE(&ifSt.Cond)
	if !ok {
		return nil, false
	}

	var typeBuf bytes.Buffer
	if err := printer.Fprint(&typeBuf, fset, fl.Type.Results.List[0].Type); err != nil {
		return nil, false
	}
	clonedRetType, err := parser.ParseExpr(typeBuf.String())
	if err != nil {
		return nil, false
	}

	usedIdents := make(map[string]bool)
	ast.Inspect(ifSt, func(n ast.Node) bool {
		if id, isId := n.(*ast.Ident); isId && id.Name != "" {
			usedIdents[id.Name] = true
		}
		return true
	})
	idx := 1
	tempName := "iifeCond_1"
	for usedIdents[tempName] {
		idx++
		tempName = fmt.Sprintf("iifeCond_%d", idx)
	}

	tempIdent := &ast.Ident{NamePos: call.End(), Name: tempName}
	if stripParens(ifSt.Cond) == call {
		ifSt.Cond = tempIdent
	} else {
		replaceSlot(tempIdent)
	}

	nStmts := len(fl.Body.List)
	ret := fl.Body.List[nStmts-1].(*ast.ReturnStmt)
	innerStmts := make([]ast.Stmt, 0, nStmts)
	innerStmts = append(innerStmts, fl.Body.List[:nStmts-1]...)
	innerStmts = append(innerStmts, &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: tempName}},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{ret.Results[0]},
	})

	varDecl := &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{{Name: tempName}},
					Type:  clonedRetType,
				},
			},
		},
	}
	innerBlk := &ast.BlockStmt{List: innerStmts}
	return &ast.BlockStmt{
		List: []ast.Stmt{varDecl, innerBlk, ifSt},
	}, true
}

func matchLeftmostIfCondIIFE(condPtr *ast.Expr) (*ast.CallExpr, *ast.FuncLit, func(ast.Expr), bool) {
	curPtr := condPtr
	for curPtr != nil && *curPtr != nil {
		switch e := (*curPtr).(type) {
		case *ast.ParenExpr:
			curPtr = &e.X
		case *ast.CallExpr:
			if fl, ok := matchSafeHoistableIIFE(e); ok {
				slot := curPtr
				return e, fl, func(repl ast.Expr) { *slot = repl }, true
			}
			return nil, nil, nil, false
		case *ast.UnaryExpr:
			switch e.Op {
			case token.NOT, token.ADD, token.SUB, token.XOR:
				curPtr = &e.X
			default:
				return nil, nil, nil, false
			}
		case *ast.BinaryExpr:
			curPtr = &e.X
		default:
			return nil, nil, nil, false
		}
	}
	return nil, nil, nil, false
}

func matchSafeHoistableIIFE(call *ast.CallExpr) (*ast.FuncLit, bool) {
	if call == nil || len(call.Args) != 0 || call.Ellipsis.IsValid() {
		return nil, false
	}
	fl, ok := stripParens(call.Fun).(*ast.FuncLit)
	if !ok || fl == nil || fl.Type == nil || fl.Body == nil {
		return nil, false
	}
	if fl.Type.Params != nil && fl.Type.Params.NumFields() != 0 {
		return nil, false
	}
	if fl.Type.Results == nil || len(fl.Type.Results.List) != 1 ||
		len(fl.Type.Results.List[0].Names) != 0 || fl.Type.Results.List[0].Type == nil {
		return nil, false
	}
	n := len(fl.Body.List)
	if n < 2 {
		return nil, false
	}
	ret, ok := fl.Body.List[n-1].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return nil, false
	}
	safe := true
	for _, st := range fl.Body.List[:n-1] {
		ast.Inspect(st, func(node ast.Node) bool {
			if !safe || node == nil {
				return false
			}
			switch node.(type) {
			case *ast.ReturnStmt, *ast.DeferStmt, *ast.BranchStmt, *ast.LabeledStmt, *ast.GoStmt, *ast.FuncLit:
				safe = false
				return false
			}
			return true
		})
		if !safe {
			return nil, false
		}
	}
	ast.Inspect(ret.Results[0], func(node ast.Node) bool {
		if !safe || node == nil {
			return false
		}
		if _, isFL := node.(*ast.FuncLit); isFL {
			safe = false
			return false
		}
		return true
	})
	if !safe {
		return nil, false
	}
	return fl, true
}

func isLambdaFactoryFuncName(name string) bool {
	return strings.HasPrefix(name, "NewLambda_") || strings.HasPrefix(name, "NewNodeLambda_")
}

// --- Wave 32B: Stateless NewNodeLambda_* Singleton Hoisting ---

func (p *J2gortPeepholePass) hoistStatelessNodeLambdas(ctx *PassContext) int {
	if ctx == nil || ctx.Pkg == nil || ctx.Info == nil || ctx.Fset == nil {
		return 0
	}
	info := ctx.Info
	existingTopNames := make(map[string]bool)
	if ctx.Pkg.Scope() != nil {
		for _, nm := range ctx.Pkg.Scope().Names() {
			existingTopNames[nm] = true
		}
	}
	for _, f := range ctx.Files {
		for _, d := range f.Decls {
			switch x := d.(type) {
			case *ast.FuncDecl:
				if x.Recv == nil && x.Name != nil {
					existingTopNames[x.Name.Name] = true
				}
			case *ast.GenDecl:
				for _, sp := range x.Specs {
					switch s := sp.(type) {
					case *ast.ValueSpec:
						for _, id := range s.Names {
							existingTopNames[id.Name] = true
						}
					case *ast.TypeSpec:
						if s.Name != nil {
							existingTopNames[s.Name.Name] = true
						}
					}
				}
			}
		}
	}

	total := 0
	for fileIdx, file := range ctx.Files {
		if !ctx.ShouldTransform(file) || !fileHasNewNodeLambdaIdent(file) {
			continue
		}
		stem := peepholeFileStem(ctx.Fset.Position(file.Pos()).Filename)
		if stem == "" {
			stem = fmt.Sprintf("f%d", fileIdx+1)
		}

		type hoistedSpec struct {
			varName  string
			typeExpr ast.Expr
			initCall ast.Expr
		}
		var hoisted []hoistedSpec
		nextIdx := 1

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name == nil || fn.Body == nil {
				continue
			}
			if fn.Name.Name == "init" || strings.HasPrefix(fn.Name.Name, "clinit") {
				continue
			}
			parents, postOrder := buildParentMap(fn.Body)
			for _, node := range postOrder {
				call, isC := node.(*ast.CallExpr)
				if !isC || call.Ellipsis.IsValid() || len(call.Args) != 1 {
					continue
				}
				parent := parents[call]
				if parent == nil {
					continue
				}
				if _, isES := parent.(*ast.ExprStmt); isES {
					continue
				}
				calleeFn, _ := calleeFunc(call, info)
				if calleeFn == nil || !strings.HasPrefix(calleeFn.Name(), "NewNodeLambda_") {
					continue
				}
				sig, okS := calleeFn.Type().(*types.Signature)
				if !okS || sig.Recv() != nil || sig.Results() == nil || sig.Results().Len() != 1 {
					continue
				}
				fl, isFL := stripParens(call.Args[0]).(*ast.FuncLit)
				if !isFL || fl == nil || fl.Type == nil || fl.Body == nil {
					continue
				}
				if !isStatelessFuncLitForHoist(fl, ctx.Pkg, info) {
					continue
				}
				retType := info.Types[call].Type
				typeExpr, okT := buildNamedTypeExprForFile(retType, file, ctx.Pkg, ctx.Importer)
				if !okT {
					continue
				}
				clonedCall, okClone := cloneExprViaParse(ctx.Fset, call)
				if !okClone {
					continue
				}
				singletonName := fmt.Sprintf("lambdaSingleton_peep_%s_%d", stem, nextIdx)
				for existingTopNames[singletonName] {
					nextIdx++
					singletonName = fmt.Sprintf("lambdaSingleton_peep_%s_%d", stem, nextIdx)
				}
				nextIdx++
				existingTopNames[singletonName] = true

				replIdent := &ast.Ident{NamePos: call.Pos(), Name: singletonName}
				if !replaceExprInParent(parent, call, replIdent, parents) {
					continue
				}
				info.Types[replIdent] = types.TypeAndValue{Type: retType}
				hoisted = append(hoisted, hoistedSpec{
					varName:  singletonName,
					typeExpr: typeExpr,
					initCall: clonedCall,
				})
				total++
			}
		}

		if len(hoisted) > 0 {
			var injectedDecls []ast.Decl
			initStmts := make([]ast.Stmt, 0, len(hoisted))
			for _, h := range hoisted {
				injectedDecls = append(injectedDecls, &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{{Name: h.varName}},
							Type:  h.typeExpr,
						},
					},
				})
				initStmts = append(initStmts, &ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: h.varName}},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{h.initCall},
				})
			}
			injectedDecls = append(injectedDecls, &ast.FuncDecl{
				Name: &ast.Ident{Name: "init"},
				Type: &ast.FuncType{Params: &ast.FieldList{}},
				Body: &ast.BlockStmt{List: initStmts},
			})
			insAt := 0
			for insAt < len(file.Decls) {
				if gd, isG := file.Decls[insAt].(*ast.GenDecl); isG && gd.Tok == token.IMPORT {
					insAt++
					continue
				}
				break
			}
			file.Decls = slices.Insert(file.Decls, insAt, injectedDecls...)
			ctx.MarkFileModified(file)
		}
	}
	return total
}

func fileHasNewNodeLambdaIdent(file *ast.File) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if id, ok := n.(*ast.Ident); ok && strings.HasPrefix(id.Name, "NewNodeLambda_") {
			found = true
			return false
		}
		return true
	})
	return found
}

func peepholeFileStem(filename string) string {
	if idx := strings.LastIndexAny(filename, "/\\"); idx >= 0 {
		filename = filename[idx+1:]
	}
	filename = strings.TrimSuffix(filename, ".go")
	var b strings.Builder
	for i := 0; i < len(filename); i++ {
		ch := filename[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' {
			b.WriteByte(ch)
		} else {
			b.WriteByte('_')
		}
	}
	return b.String()
}

func isStatelessFuncLitForHoist(fl *ast.FuncLit, pkg *types.Package, info *types.Info) bool {
	if fl == nil || pkg == nil || info == nil {
		return false
	}
	hasTypeParam := false
	localDefs := make(map[types.Object]bool)
	ast.Inspect(fl, func(n ast.Node) bool {
		if hasTypeParam || n == nil {
			return false
		}
		if e, isE := n.(ast.Expr); isE {
			if tv, ok := info.Types[e]; ok && typeContainsTypeParam(tv.Type) {
				hasTypeParam = true
				return false
			}
		}
		if id, isId := n.(*ast.Ident); isId {
			if def := info.Defs[id]; def != nil {
				localDefs[def] = true
			}
		}
		return true
	})
	if hasTypeParam {
		return false
	}

	flParents, _ := buildParentMap(fl)
	safe := true
	ast.Inspect(fl, func(n ast.Node) bool {
		if !safe || n == nil {
			return false
		}
		id, isId := n.(*ast.Ident)
		if !isId || id.Name == "_" {
			return true
		}
		if sel, isSel := flParents[id].(*ast.SelectorExpr); isSel && sel.Sel == id {
			return true
		}
		if info.Defs[id] != nil {
			return true
		}
		use := info.Uses[id]
		if use == nil {
			safe = false
			return false
		}
		switch u := use.(type) {
		case *types.PkgName, *types.Builtin, *types.Nil:
			return true
		case *types.Func:
			if u.Pkg() != nil && u.Parent() == u.Pkg().Scope() {
				return true
			}
		case *types.Const:
			if u.Parent() == types.Universe || (u.Pkg() != nil && u.Parent() == u.Pkg().Scope()) {
				return true
			}
		case *types.TypeName:
			if _, isTP := u.Type().(*types.TypeParam); isTP {
				safe = false
				return false
			}
			if u.Parent() == types.Universe || (u.Pkg() != nil && u.Parent() == u.Pkg().Scope()) {
				return true
			}
		case *types.Var:
			if localDefs[u] {
				return true
			}
			if u.IsField() {
				if kv, isKV := flParents[id].(*ast.KeyValueExpr); isKV && kv.Key == id {
					return true
				}
				safe = false
				return false
			}
			if u.Pkg() != nil && u.Parent() == u.Pkg().Scope() {
				return true
			}
		case *types.Label:
			if localDefs[u] {
				return true
			}
		}
		safe = false
		return false
	})
	return safe
}

func buildNamedTypeExprForFile(t types.Type, file *ast.File, pkg *types.Package, imp types.Importer) (ast.Expr, bool) {
	if t == nil || file == nil || pkg == nil {
		return nil, false
	}
	named, ok := types.Unalias(t).(*types.Named)
	if !ok || named.Obj() == nil || (named.TypeArgs() != nil && named.TypeArgs().Len() > 0) {
		return nil, false
	}
	objPkg := named.Obj().Pkg()
	if objPkg == nil || objPkg == pkg {
		return &ast.Ident{Name: named.Obj().Name()}, true
	}
	for _, sp := range file.Imports {
		if sp.Name != nil && (sp.Name.Name == "_" || sp.Name.Name == ".") {
			continue
		}
		p, err := strconv.Unquote(sp.Path.Value)
		if err != nil || p != objPkg.Path() {
			continue
		}
		qual := specLocalName(sp, imp)
		if qual != "" {
			return &ast.SelectorExpr{
				X:   &ast.Ident{Name: qual},
				Sel: &ast.Ident{Name: named.Obj().Name()},
			}, true
		}
	}
	return nil, false
}

func cloneExprViaParse(fset *token.FileSet, expr ast.Expr) (ast.Expr, bool) {
	if fset == nil || expr == nil {
		return nil, false
	}
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, expr); err != nil {
		return nil, false
	}
	parsed, err := parser.ParseExpr(buf.String())
	if err != nil {
		return nil, false
	}
	return parsed, true
}

// --- Wave 32B: Unboxed Stack-Closure NodeUtil_Has Lowering ---

func (p *J2gortPeepholePass) rewriteNodeUtilHasNodeLambdas(ctx *PassContext) (int, error) {
	if ctx == nil || ctx.Pkg == nil || ctx.Pkg.Scope() == nil || ctx.Info == nil || ctx.Fset == nil {
		return 0, nil
	}
	info := ctx.Info
	hasObj, ok := ctx.Pkg.Scope().Lookup("NodeUtil_Has").(*types.Func)
	if !ok || hasObj == nil {
		return 0, nil
	}
	hasSig, okS := hasObj.Type().(*types.Signature)
	if !okS || hasSig.Recv() != nil || hasSig.Variadic() ||
		hasSig.Params() == nil || hasSig.Params().Len() != 3 ||
		hasSig.Results() == nil || hasSig.Results().Len() != 1 ||
		!types.Identical(hasSig.Results().At(0).Type(), types.Typ[types.Bool]) {
		return 0, nil
	}
	nodeT := hasSig.Params().At(0).Type()
	predT := hasSig.Params().At(1).Type()
	travT := hasSig.Params().At(2).Type()
	if !isNamedUint32NodeType(nodeT) || !types.Identical(predT, travT) {
		return 0, nil
	}
	predNamed, okPN := types.Unalias(predT).(*types.Named)
	if !okPN || predNamed.Obj() == nil || predNamed.Obj().Name() != "Ext_com_google_common_base_Predicate" ||
		!isJ2gortPkg(predNamed.Obj().Pkg()) {
		return 0, nil
	}

	var hasFile *ast.File
	var hasDecl *ast.FuncDecl
	for _, f := range ctx.Files {
		for _, d := range f.Decls {
			fn, isFn := d.(*ast.FuncDecl)
			if isFn && fn.Recv == nil && fn.Name != nil && fn.Name.Name == "NodeUtil_Has" && fn.Body != nil {
				hasFile = f
				hasDecl = fn
				break
			}
		}
		if hasDecl != nil {
			break
		}
	}
	if hasFile == nil || hasDecl == nil || !ctx.ShouldTransform(hasFile) ||
		hasDecl.Type == nil || hasDecl.Type.Params == nil || len(hasDecl.Type.Params.List) < 2 {
		return 0, nil
	}
	j2Qual := extractPredTestNodeQualifier(hasDecl.Body)
	if j2Qual == "" || !bodyHasIdentNamed(hasDecl.Body, "GetFirstChild") || !bodyHasIdentNamed(hasDecl.Body, "GetNext") {
		return 0, nil
	}
	var nodeTypeBuf, predTypeBuf bytes.Buffer
	if err := printer.Fprint(&nodeTypeBuf, ctx.Fset, hasDecl.Type.Params.List[0].Type); err != nil {
		return 0, nil
	}
	if err := printer.Fprint(&predTypeBuf, ctx.Fset, hasDecl.Type.Params.List[1].Type); err != nil {
		return 0, nil
	}
	nodeTypeSrc := nodeTypeBuf.String()
	predTypeSrc := predTypeBuf.String()

	const (
		helperPredAlwaysTrue = "nodeUtilHasPredAlwaysTrueFn_"
		helperPred           = "nodeUtilHasPredFn_"
		helperTrav           = "nodeUtilHasTravFn_"
		helperBoth           = "nodeUtilHasBothFn_"
	)
	needHelper := make(map[string]bool)
	total := 0

	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) || !bodyHasIdentNamed(file, "NodeUtil_Has") {
			continue
		}
		fileModified := false
		for _, decl := range file.Decls {
			fn, isFn := decl.(*ast.FuncDecl)
			if !isFn || fn.Name == nil || fn.Body == nil {
				continue
			}
			if fn.Name.Name == "NodeUtil_Has" || fn.Name.Name == "init" ||
				strings.HasPrefix(fn.Name.Name, "clinit") || strings.HasPrefix(fn.Name.Name, "nodeUtilHas") {
				continue
			}
			parents, _ := buildParentMap(fn.Body)
			n := rewriteNodeUtilHasInTree(fn.Body, fn.Body, hasObj, parents, info, needHelper)
			if n > 0 {
				fileModified = true
				total += n
			}
		}
		if fileModified {
			ctx.MarkFileModified(file)
		}
	}

	if total == 0 {
		return 0, nil
	}

	type helperTemplate struct {
		name string
		code string
	}
	templates := []helperTemplate{
		{
			name: helperPredAlwaysTrue,
			code: fmt.Sprintf(`func %[1]s(node %[2]s, pred func(uint32) bool) bool {
	if pred(uint32(node)) {
		return true
	}
	for c := node.GetFirstChild(); c != 0; c = c.GetNext() {
		if %[1]s(c, pred) {
			return true
		}
	}
	return false
}`, helperPredAlwaysTrue, nodeTypeSrc),
		},
		{
			name: helperPred,
			code: fmt.Sprintf(`func %[1]s(node %[2]s, pred func(uint32) bool, traverseChildrenPred %[3]s) bool {
	if pred(uint32(node)) {
		return true
	}
	if !%[4]s.PredTestNode(traverseChildrenPred, uint32(node)) {
		return false
	}
	for c := node.GetFirstChild(); c != 0; c = c.GetNext() {
		if %[1]s(c, pred, traverseChildrenPred) {
			return true
		}
	}
	return false
}`, helperPred, nodeTypeSrc, predTypeSrc, j2Qual),
		},
		{
			name: helperTrav,
			code: fmt.Sprintf(`func %[1]s(node %[2]s, pred %[3]s, traverseChildrenPred func(uint32) bool) bool {
	if %[4]s.PredTestNode(pred, uint32(node)) {
		return true
	}
	if !traverseChildrenPred(uint32(node)) {
		return false
	}
	for c := node.GetFirstChild(); c != 0; c = c.GetNext() {
		if %[1]s(c, pred, traverseChildrenPred) {
			return true
		}
	}
	return false
}`, helperTrav, nodeTypeSrc, predTypeSrc, j2Qual),
		},
		{
			name: helperBoth,
			code: fmt.Sprintf(`func %[1]s(node %[2]s, pred func(uint32) bool, traverseChildrenPred func(uint32) bool) bool {
	if pred(uint32(node)) {
		return true
	}
	if !traverseChildrenPred(uint32(node)) {
		return false
	}
	for c := node.GetFirstChild(); c != 0; c = c.GetNext() {
		if %[1]s(c, pred, traverseChildrenPred) {
			return true
		}
	}
	return false
}`, helperBoth, nodeTypeSrc),
		},
	}

	var specs []funcInsertSpec
	for _, tmpl := range templates {
		if !needHelper[tmpl.name] || fileHasTopLevelFuncDecl(hasFile, tmpl.name) {
			continue
		}
		pf, err := parser.ParseFile(ctx.Fset, "", "package p\n"+tmpl.code+"\n", 0)
		if err != nil {
			return total, err
		}
		specs = append(specs, funcInsertSpec{after: hasDecl, decl: pf.Decls[0].(*ast.FuncDecl)})
	}
	if len(specs) > 0 {
		spliceSynthesizedDecls(hasFile, specs)
		ctx.MarkFileModified(hasFile)
	}
	return total, nil
}

func extractPredTestNodeQualifier(body *ast.BlockStmt) string {
	qual := ""
	ast.Inspect(body, func(n ast.Node) bool {
		if qual != "" || n == nil {
			return false
		}
		sel, ok := n.(*ast.SelectorExpr)
		if !ok || sel.Sel == nil || sel.Sel.Name != "PredTestNode" {
			return true
		}
		if id, okI := stripParens(sel.X).(*ast.Ident); okI {
			qual = id.Name
			return false
		}
		return true
	})
	return qual
}

func bodyHasIdentNamed(root ast.Node, want string) bool {
	found := false
	ast.Inspect(root, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if id, ok := n.(*ast.Ident); ok && id.Name == want {
			found = true
			return false
		}
		return true
	})
	return found
}

func fileHasTopLevelFuncDecl(file *ast.File, name string) bool {
	for _, d := range file.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name != nil && fn.Name.Name == name {
			return true
		}
	}
	return false
}

func matchGuavaNodeLambdaExpr(expr ast.Expr, info *types.Info) *ast.FuncLit {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() || len(call.Args) != 1 {
		return nil
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 || fnName != "NewNodeLambda_Ext_com_google_common_base_Predicate" {
		return nil
	}
	fl, okFL := stripParens(call.Args[0]).(*ast.FuncLit)
	if !okFL || fl == nil || fl.Type == nil || fl.Body == nil {
		return nil
	}
	return fl
}

func isGuavaPredicatesAlwaysTrueCall(expr ast.Expr, info *types.Info) bool {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() || len(call.Args) != 0 {
		return false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	return isJ2 && fnName == "Ext_com_google_common_base_Predicates_AlwaysTrue"
}

type localNodeLambdaBinding struct {
	declIdx     int
	defIdent    *ast.Ident
	vObj        *types.Var
	fl          *ast.FuncLit
	discardIdxs []int
}

func resolveSingleUseBlockNodeLambda(
	argExpr ast.Expr,
	blk *ast.BlockStmt,
	useStmtIdx int,
	hasCall *ast.CallExpr,
	fnBody *ast.BlockStmt,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) (*ast.FuncLit, map[int]bool) {
	if fl := matchGuavaNodeLambdaExpr(argExpr, info); fl != nil {
		return fl, nil
	}
	id, isId := stripParens(argExpr).(*ast.Ident)
	if !isId || blk == nil || fnBody == nil {
		return nil, nil
	}
	vObj := varFromObj(info.Uses[id])
	if vObj == nil {
		return nil, nil
	}

	// Ensure `hasCall` is not nested inside a closure, loop, or goroutine within `blk.List[useStmtIdx]`.
	useStmt := blk.List[useStmtIdx]
	for cur := ast.Node(hasCall); cur != nil && cur != useStmt; cur = parents[cur] {
		switch cur.(type) {
		case *ast.FuncLit, *ast.ForStmt, *ast.RangeStmt, *ast.GoStmt, *ast.DeferStmt:
			return nil, nil
		}
	}

	declIdx := -1
	var defIdent *ast.Ident
	var fl *ast.FuncLit
	for i := 0; i < useStmtIdx; i++ {
		dId, rhs, okD := matchLocalSingleInitStmt(blk.List[i])
		if !okD || dId == nil || varFromObj(info.Defs[dId]) != vObj {
			continue
		}
		candFL := matchGuavaNodeLambdaExpr(rhs, info)
		if candFL == nil {
			return nil, nil
		}
		declIdx = i
		defIdent = dId
		fl = candFL
		break
	}
	if declIdx < 0 || fl == nil {
		return nil, nil
	}

	removeIdxs := map[int]bool{declIdx: true}
	validUses := true
	ast.Inspect(fnBody, func(n ast.Node) bool {
		if !validUses || n == nil {
			return false
		}
		uId, okU := n.(*ast.Ident)
		if !okU || uId == defIdent || uId == id {
			return true
		}
		if varFromObj(info.Uses[uId]) != vObj && varFromObj(info.Defs[uId]) != vObj {
			return true
		}
		// Allow `_ = <var>` blank discard statements in `blk.List` between `declIdx+1` and `len(blk.List)`.
		for k := declIdx + 1; k < len(blk.List); k++ {
			if k != useStmtIdx && isBlankDiscardOfIdent(blk.List[k], uId) {
				removeIdxs[k] = true
				return true
			}
		}
		validUses = false
		return false
	})
	if !validUses {
		return nil, nil
	}

	// Collect outer local variables captured by `fl` so we verify zero mutation/shadowing between `declIdx` and `hasCall`.
	flLocalDefs := make(map[types.Object]bool)
	capturedVars := make(map[*types.Var]bool)
	capturedNames := make(map[string]bool)
	ast.Inspect(fl, func(n ast.Node) bool {
		if uId, okU := n.(*ast.Ident); okU {
			if def := info.Defs[uId]; def != nil {
				flLocalDefs[def] = true
			}
		}
		return true
	})
	ast.Inspect(fl, func(n ast.Node) bool {
		uId, okU := n.(*ast.Ident)
		if !okU || info.Defs[uId] != nil {
			return true
		}
		if uObj := info.Uses[uId]; uObj != nil && !flLocalDefs[uObj] {
			capturedNames[uId.Name] = true
			if uv := varFromObj(uObj); uv != nil && !uv.IsField() {
				capturedVars[uv] = true
			}
		}
		return true
	})

	checkNoMutOrShadow := func(root ast.Node) bool {
		ok := true
		ast.Inspect(root, func(n ast.Node) bool {
			if !ok || n == nil || n == hasCall {
				return false
			}
			switch s := n.(type) {
			case *ast.ReturnStmt, *ast.BranchStmt, *ast.GoStmt, *ast.DeferStmt:
				if s != useStmt {
					ok = false
					return false
				}
			case *ast.Ident:
				if def := info.Defs[s]; def != nil && capturedNames[s.Name] {
					ok = false
					return false
				}
			case *ast.AssignStmt:
				for _, lhs := range s.Lhs {
					if lhsId, isL := stripParens(lhs).(*ast.Ident); isL {
						if lv := varFromObj(info.Uses[lhsId]); lv != nil && capturedVars[lv] {
							ok = false
							return false
						}
					}
				}
			case *ast.IncDecStmt:
				if xId, isX := stripParens(s.X).(*ast.Ident); isX {
					if xv := varFromObj(info.Uses[xId]); xv != nil && capturedVars[xv] {
						ok = false
						return false
					}
				}
			case *ast.UnaryExpr:
				if s.Op == token.AND {
					if xId, isX := stripParens(s.X).(*ast.Ident); isX {
						if xv := varFromObj(info.Uses[xId]); xv != nil && capturedVars[xv] {
							ok = false
							return false
						}
					}
				}
			}
			return true
		})
		return ok
	}

	for k := declIdx + 1; k < useStmtIdx; k++ {
		if removeIdxs[k] {
			continue
		}
		if !checkNoMutOrShadow(blk.List[k]) {
			return nil, nil
		}
	}
	if !checkNoMutOrShadow(useStmt) {
		return nil, nil
	}
	return fl, removeIdxs
}

func isBlankDiscardOfIdent(stmt ast.Stmt, rhsId *ast.Ident) bool {
	as, ok := stmt.(*ast.AssignStmt)
	return ok && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 &&
		isIdentNamed(as.Lhs[0], "_") && stripParens(as.Rhs[0]) == rhsId
}

func rewriteNodeUtilHasInTree(
	node ast.Node,
	fnBody *ast.BlockStmt,
	hasObj *types.Func,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	needHelper map[string]bool,
) int {
	if node == nil {
		return 0
	}
	count := 0
	if blk, isBlk := node.(*ast.BlockStmt); isBlk {
		toRemove := make(map[int]bool)
		for i := 0; i < len(blk.List); i++ {
			if toRemove[i] {
				continue
			}
			st := blk.List[i]
			// Recurse into nested blocks first.
			ast.Inspect(st, func(child ast.Node) bool {
				if child == nil || child == st {
					return true
				}
				if subBlk, okB := child.(*ast.BlockStmt); okB {
					count += rewriteNodeUtilHasInTree(subBlk, fnBody, hasObj, parents, info, needHelper)
					return false
				}
				return true
			})

			// Find direct `NodeUtil_Has(...)` calls owned by `st` (not inside nested `BlockStmt`s).
			var hasCalls []*ast.CallExpr
			ast.Inspect(st, func(child ast.Node) bool {
				if child == nil {
					return false
				}
				if child != st {
					if _, okB := child.(*ast.BlockStmt); okB {
						return false
					}
				}
				c, isC := child.(*ast.CallExpr)
				if !isC || c.Ellipsis.IsValid() || len(c.Args) != 3 {
					return true
				}
				callee, _ := calleeFunc(c, info)
				if callee != nil && funcOrigin(callee) == hasObj {
					if _, isId := stripParens(c.Fun).(*ast.Ident); isId {
						hasCalls = append(hasCalls, c)
					}
				}
				return true
			})

			for _, c := range hasCalls {
				predFL, predRem := resolveSingleUseBlockNodeLambda(c.Args[1], blk, i, c, fnBody, parents, info)
				travAlwaysTrue := isGuavaPredicatesAlwaysTrueCall(c.Args[2], info)
				var travFL *ast.FuncLit
				var travRem map[int]bool
				if !travAlwaysTrue {
					travFL, travRem = resolveSingleUseBlockNodeLambda(c.Args[2], blk, i, c, fnBody, parents, info)
				}
				if predFL == nil && travFL == nil {
					continue
				}
				var hName string
				var newArgs []ast.Expr
				switch {
				case predFL != nil && travAlwaysTrue:
					hName = "nodeUtilHasPredAlwaysTrueFn_"
					newArgs = []ast.Expr{c.Args[0], predFL}
				case predFL != nil && travFL != nil:
					hName = "nodeUtilHasBothFn_"
					newArgs = []ast.Expr{c.Args[0], predFL, travFL}
				case predFL != nil:
					hName = "nodeUtilHasPredFn_"
					newArgs = []ast.Expr{c.Args[0], predFL, c.Args[2]}
				default:
					hName = "nodeUtilHasTravFn_"
					newArgs = []ast.Expr{c.Args[0], c.Args[1], travFL}
				}
				funId := stripParens(c.Fun).(*ast.Ident)
				funId.Name = hName
				c.Args = newArgs
				needHelper[hName] = true
				for idx := range predRem {
					toRemove[idx] = true
				}
				for idx := range travRem {
					toRemove[idx] = true
				}
				count++
			}
		}
		if len(toRemove) > 0 {
			filtered := make([]ast.Stmt, 0, len(blk.List)-len(toRemove))
			for idx, st := range blk.List {
				if !toRemove[idx] {
					filtered = append(filtered, st)
				}
			}
			blk.List = filtered
		}
		return count
	}

	ast.Inspect(node, func(child ast.Node) bool {
		if child == nil || child == node {
			return true
		}
		if subBlk, okB := child.(*ast.BlockStmt); okB {
			count += rewriteNodeUtilHasInTree(subBlk, fnBody, hasObj, parents, info, needHelper)
			return false
		}
		return true
	})
	return count
}

// --- Wave 32B: Stack-Lower Immediate Fluent NewNodeTraversal_Builder Chains ---

type ntbSetterStep struct {
	methodName string
	args       []ast.Expr
	pos        token.Pos
}

type ntbPrivHelperSpec struct {
	recvOwner   *types.Named
	recvVarName string
	needClinit  bool
	ctorArgs    []ast.Expr
	setters     []ntbSetterStep
}

func (p *J2gortPeepholePass) rewriteStackNodeTraversalBuilders(ctx *PassContext) int {
	if ctx == nil || ctx.Pkg == nil || ctx.Pkg.Scope() == nil || ctx.Info == nil || ctx.Fset == nil {
		return 0
	}
	info := ctx.Info
	pkg := ctx.Pkg

	ntbTypeName, okTN := pkg.Scope().Lookup("NodeTraversal_Builder").(*types.TypeName)
	if !okTN || ntbTypeName == nil {
		return 0
	}
	ntbNamed := namedOrigin(ntbTypeName.Type())
	if ntbNamed == nil {
		return 0
	}
	if _, isSt := ntbNamed.Underlying().(*types.Struct); !isSt {
		return 0
	}
	ntbPtrType := types.NewPointer(ntbNamed)

	initFn, okInit := pkg.Scope().Lookup("InitNodeTraversal_Builder").(*types.Func)
	if !okInit || initFn == nil {
		return 0
	}
	initSig, okIS := initFn.Type().(*types.Signature)
	if !okIS || initSig.Recv() != nil || initSig.Variadic() ||
		initSig.Results().Len() != 0 || initSig.Params() == nil || initSig.Params().Len() < 1 ||
		!types.Identical(initSig.Params().At(0).Type(), ntbPtrType) {
		return 0
	}

	hasClinitFn := false
	if clFn, okCl := pkg.Scope().Lookup("clinit_NodeTraversal").(*types.Func); okCl && clFn != nil {
		if clSig, okCS := clFn.Type().(*types.Signature); okCS && clSig.Recv() == nil &&
			clSig.Params().Len() == 0 && clSig.Results().Len() == 0 {
			hasClinitFn = true
		}
	}

	sBuilderCallsClinit := false
	safeSetters := map[string]bool{
		"SetCompiler":                          true,
		"SetCallback":                          true,
		"SetCallback_1_NodeTraversal_Callback": true,
		"SetScopeCreator":                      true,
		"SetObeyDestructuringAndDefaultValueExecutionOrder": true,
		"SetObeyDestructuring":                              true,
	}
	safeTerminals := map[string]bool{
		"Traverse":          true,
		"TraverseRoots":     true,
		"TraverseAtScope":   true,
		"TraverseWithScope": true,
	}

	// Verify actual *NodeTraversal_Builder method bodies in ctx.Files so escaping setters/terminals are always rejected.
	for _, f := range ctx.Files {
		for _, d := range f.Decls {
			fn, isFn := d.(*ast.FuncDecl)
			if !isFn || fn.Name == nil || fn.Body == nil {
				continue
			}
			if fn.Recv == nil {
				if fn.Name.Name == "NodeTraversal_S_Builder" && hasClinitFn && bodyHasIdentNamed(fn.Body, "clinit_NodeTraversal") {
					sBuilderCallsClinit = true
				}
				continue
			}
			if len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 {
				continue
			}
			rId := fn.Recv.List[0].Names[0]
			rv := varFromObj(info.Defs[rId])
			if rv == nil || !types.Identical(rv.Type(), ntbPtrType) {
				continue
			}
			mName := fn.Name.Name
			fnObj, _ := info.Defs[fn.Name].(*types.Func)
			if fnObj == nil {
				continue
			}
			mSig, okMS := fnObj.Type().(*types.Signature)
			if !okMS || mSig.Variadic() {
				delete(safeSetters, mName)
				delete(safeTerminals, mName)
				continue
			}
			if strings.HasPrefix(mName, "Set") && mSig.Results().Len() == 1 && types.Identical(mSig.Results().At(0).Type(), ntbPtrType) {
				if isNonEscapingNtbSetterBody(fn.Body, rId.Name, rv, info) {
					safeSetters[mName] = true
				} else {
					delete(safeSetters, mName)
				}
			} else if strings.HasPrefix(mName, "Traverse") && mSig.Results().Len() == 0 {
				if isNonEscapingNtbTerminalBody(fn.Body, rId.Name, rv, info) {
					safeTerminals[mName] = true
				} else {
					delete(safeTerminals, mName)
				}
			}
		}
	}

	matchRootBuilderCall := func(expr ast.Expr) (ctorArgs []ast.Expr, needClinit bool, ok bool) {
		c, isC := stripParens(expr).(*ast.CallExpr)
		if !isC || c.Ellipsis.IsValid() {
			return nil, false, false
		}
		calleeFn, _ := calleeFunc(c, info)
		if calleeFn == nil || calleeFn.Pkg() != pkg {
			return nil, false, false
		}
		cSig, okCS := calleeFn.Type().(*types.Signature)
		if !okCS || cSig.Recv() != nil || cSig.Variadic() ||
			cSig.Results().Len() != 1 || !types.Identical(cSig.Results().At(0).Type(), ntbPtrType) {
			return nil, false, false
		}
		switch calleeFn.Name() {
		case "NewNodeTraversal_Builder":
			needClinit = false
		case "NodeTraversal_S_Builder":
			needClinit = sBuilderCallsClinit
		default:
			return nil, false, false
		}
		if initSig.Params().Len() != 1+len(c.Args) {
			return nil, false, false
		}
		for k, arg := range c.Args {
			argT := info.Types[arg].Type
			if argT == nil || !types.AssignableTo(argT, initSig.Params().At(k+1).Type()) {
				return nil, false, false
			}
		}
		return c.Args, needClinit, true
	}

	peelBuilderSetters := func(expr ast.Expr) (rootExpr ast.Expr, setters []ntbSetterStep, ok bool) {
		var rev []ntbSetterStep
		cur := stripParens(expr)
		for depth := 0; depth < 16 && cur != nil; depth++ {
			c, isC := cur.(*ast.CallExpr)
			if !isC || c.Ellipsis.IsValid() {
				return nil, nil, false
			}
			sel, isS := stripParens(c.Fun).(*ast.SelectorExpr)
			if !isS || sel.Sel == nil || !safeSetters[sel.Sel.Name] {
				break
			}
			calleeFn, _ := calleeFunc(c, info)
			if calleeFn == nil || calleeFn.Pkg() != pkg {
				return nil, nil, false
			}
			mSig, okMS := calleeFn.Type().(*types.Signature)
			if !okMS || mSig.Recv() == nil || !types.Identical(mSig.Recv().Type(), ntbPtrType) ||
				mSig.Results().Len() != 1 || !types.Identical(mSig.Results().At(0).Type(), ntbPtrType) {
				return nil, nil, false
			}
			rev = append(rev, ntbSetterStep{
				methodName: sel.Sel.Name,
				args:       c.Args,
				pos:        c.Pos(),
			})
			cur = stripParens(sel.X)
		}
		slices.Reverse(rev)
		return cur, rev, true
	}

	total := 0
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		// Discover same-file 0-arg `priv_*` helper methods returning a pure fluent `*NodeTraversal_Builder` chain.
		var privSpecs map[*types.Func]ntbPrivHelperSpec
		for _, d := range file.Decls {
			fn, isFn := d.(*ast.FuncDecl)
			if !isFn || fn.Recv == nil || fn.Name == nil || fn.Body == nil || !strings.HasPrefix(fn.Name.Name, "priv_") {
				continue
			}
			if len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 {
				continue
			}
			rId := fn.Recv.List[0].Names[0]
			rv := varFromObj(info.Defs[rId])
			if rv == nil || rId.Name == "_" {
				continue
			}
			ownerNamed := ptrNamedOrigin(rv.Type())
			if ownerNamed == nil {
				continue
			}
			fnObj, _ := info.Defs[fn.Name].(*types.Func)
			if fnObj == nil {
				continue
			}
			fSig, okFS := fnObj.Type().(*types.Signature)
			if !okFS || fSig.Params().Len() != 0 || fSig.Results().Len() != 1 ||
				!types.Identical(fSig.Results().At(0).Type(), ntbPtrType) {
				continue
			}
			var retExpr ast.Expr
			validBody := true
			for _, st := range fn.Body.List {
				if isBlankDiscardStmt(st) {
					continue
				}
				if ret, isR := st.(*ast.ReturnStmt); isR && retExpr == nil && len(ret.Results) == 1 {
					retExpr = ret.Results[0]
				} else {
					validBody = false
					break
				}
			}
			if !validBody || retExpr == nil || bodyContainsFuncLit(retExpr) {
				continue
			}
			rootExpr, innerSetters, okP := peelBuilderSetters(retExpr)
			if !okP {
				continue
			}
			ctorArgs, needClinit, okR := matchRootBuilderCall(rootExpr)
			if !okR {
				continue
			}
			if privSpecs == nil {
				privSpecs = make(map[*types.Func]ntbPrivHelperSpec)
			}
			privSpecs[funcOrigin(fnObj)] = ntbPrivHelperSpec{
				recvOwner:   ownerNamed,
				recvVarName: rId.Name,
				needClinit:  needClinit,
				ctorArgs:    ctorArgs,
				setters:     innerSetters,
			}
		}

		tryRewriteStmt := func(stmt ast.Stmt, enclosingBlk *ast.BlockStmt) (*ast.BlockStmt, bool) {
			es, isES := stmt.(*ast.ExprStmt)
			if !isES {
				return nil, false
			}
			termCall, isC := stripParens(es.X).(*ast.CallExpr)
			if !isC || termCall.Ellipsis.IsValid() {
				return nil, false
			}
			termSel, isS := stripParens(termCall.Fun).(*ast.SelectorExpr)
			if !isS || termSel.Sel == nil || !safeTerminals[termSel.Sel.Name] {
				return nil, false
			}
			termFn, _ := calleeFunc(termCall, info)
			if termFn == nil || termFn.Pkg() != pkg {
				return nil, false
			}
			tSig, okTS := termFn.Type().(*types.Signature)
			if !okTS || tSig.Recv() == nil || !types.Identical(tSig.Recv().Type(), ntbPtrType) || tSig.Results().Len() != 0 {
				return nil, false
			}

			rootExpr, outerSetters, okP := peelBuilderSetters(termSel.X)
			if !okP {
				return nil, false
			}

			var ctorArgs []ast.Expr
			var allSetters []ntbSetterStep
			needClinit := false

			if cArgs, nCl, okR := matchRootBuilderCall(rootExpr); okR {
				ctorArgs = cArgs
				needClinit = nCl
				allSetters = outerSetters
			} else if len(privSpecs) > 0 {
				pCall, isPC := stripParens(rootExpr).(*ast.CallExpr)
				if !isPC || pCall.Ellipsis.IsValid() || len(pCall.Args) != 0 {
					return nil, false
				}
				pSel, isPS := stripParens(pCall.Fun).(*ast.SelectorExpr)
				if !isPS || pSel.Sel == nil {
					return nil, false
				}
				callerRecvId, isId := stripParens(pSel.X).(*ast.Ident)
				if !isId || callerRecvId.Name == "_" {
					return nil, false
				}
				pCallee, _ := calleeFunc(pCall, info)
				if pCallee == nil {
					return nil, false
				}
				spec, okSp := privSpecs[funcOrigin(pCallee)]
				if !okSp {
					return nil, false
				}
				callerVar := varFromObj(info.Uses[callerRecvId])
				if callerVar == nil || ptrNamedOrigin(callerVar.Type()) != spec.recvOwner {
					return nil, false
				}
				needClinit = spec.needClinit
				clonedCtor := make([]ast.Expr, len(spec.ctorArgs))
				for k, a := range spec.ctorArgs {
					cl, okCl := cloneExprReplacingIdent(ctx.Fset, a, spec.recvVarName, callerRecvId.Name)
					if !okCl {
						return nil, false
					}
					clonedCtor[k] = cl
				}
				ctorArgs = clonedCtor
				for _, stp := range spec.setters {
					clArgs := make([]ast.Expr, len(stp.args))
					for k, a := range stp.args {
						cl, okCl := cloneExprReplacingIdent(ctx.Fset, a, spec.recvVarName, callerRecvId.Name)
						if !okCl {
							return nil, false
						}
						clArgs[k] = cl
					}
					allSetters = append(allSetters, ntbSetterStep{
						methodName: stp.methodName,
						args:       clArgs,
						pos:        es.Pos(),
					})
				}
				allSetters = append(allSetters, outerSetters...)
			} else {
				return nil, false
			}

			var allArgExprs []ast.Expr
			allArgExprs = append(allArgExprs, ctorArgs...)
			for _, stp := range allSetters {
				allArgExprs = append(allArgExprs, stp.args...)
			}
			allArgExprs = append(allArgExprs, termCall.Args...)

			ntbName := pickFreshNtbLocalName(enclosingBlk, allArgExprs)
			pos := es.Pos()

			var stmts []ast.Stmt
			if needClinit {
				stmts = append(stmts, &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.Ident{NamePos: pos, Name: "clinit_NodeTraversal"},
					},
				})
			}

			// `var ntb_ NodeTraversal_Builder`
			stmts = append(stmts, &ast.DeclStmt{
				Decl: &ast.GenDecl{
					TokPos: pos,
					Tok:    token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{{NamePos: pos, Name: ntbName}},
							Type:  &ast.Ident{NamePos: pos, Name: "NodeTraversal_Builder"},
						},
					},
				},
			})

			// `InitNodeTraversal_Builder(&ntb_, ctorArgs...)`
			initCallArgs := make([]ast.Expr, 0, 1+len(ctorArgs))
			initCallArgs = append(initCallArgs, &ast.UnaryExpr{
				OpPos: pos,
				Op:    token.AND,
				X:     &ast.Ident{NamePos: pos, Name: ntbName},
			})
			for _, a := range ctorArgs {
				initCallArgs = append(initCallArgs, stripParens(a))
			}
			stmts = append(stmts, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun:  &ast.Ident{NamePos: pos, Name: "InitNodeTraversal_Builder"},
					Args: initCallArgs,
				},
			})

			makeAddrParens := func(p token.Pos) *ast.ParenExpr {
				return &ast.ParenExpr{
					Lparen: p,
					X: &ast.UnaryExpr{
						OpPos: p,
						Op:    token.AND,
						X:     &ast.Ident{NamePos: p, Name: ntbName},
					},
					Rparen: p,
				}
			}

			// `(&ntb_).<Setter>(args...)`
			for _, stp := range allSetters {
				stmts = append(stmts, &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   makeAddrParens(stp.pos),
							Sel: &ast.Ident{NamePos: stp.pos, Name: stp.methodName},
						},
						Args: stp.args,
					},
				})
			}

			// `(&ntb_).<TerminalTraverseMethod>(termArgs...)`
			stmts = append(stmts, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   makeAddrParens(termCall.Pos()),
						Sel: &ast.Ident{NamePos: termSel.Sel.Pos(), Name: termSel.Sel.Name},
					},
					Args: termCall.Args,
				},
			})

			return &ast.BlockStmt{
				Lbrace: pos,
				List:   stmts,
				Rbrace: es.End(),
			}, true
		}

		fileCount := 0
		var walk func(ast.Node, *ast.BlockStmt)
		walk = func(n ast.Node, curBlk *ast.BlockStmt) {
			if n == nil {
				return
			}
			switch x := n.(type) {
			case *ast.BlockStmt:
				for i := 0; i < len(x.List); i++ {
					walk(x.List[i], x)
					if repl, ok := tryRewriteStmt(x.List[i], x); ok {
						x.List[i] = repl
						fileCount++
					}
				}
			case *ast.CaseClause:
				for i := 0; i < len(x.Body); i++ {
					walk(x.Body[i], curBlk)
					if repl, ok := tryRewriteStmt(x.Body[i], curBlk); ok {
						x.Body[i] = repl
						fileCount++
					}
				}
			case *ast.CommClause:
				for i := 0; i < len(x.Body); i++ {
					walk(x.Body[i], curBlk)
					if repl, ok := tryRewriteStmt(x.Body[i], curBlk); ok {
						x.Body[i] = repl
						fileCount++
					}
				}
			default:
				ast.Inspect(x, func(child ast.Node) bool {
					if child == nil || child == x {
						return true
					}
					switch child.(type) {
					case *ast.BlockStmt, *ast.CaseClause, *ast.CommClause:
						walk(child, curBlk)
						return false
					}
					return true
				})
			}
		}

		for _, d := range file.Decls {
			if fn, isFn := d.(*ast.FuncDecl); isFn && fn.Body != nil {
				walk(fn.Body, fn.Body)
			}
		}
		if fileCount > 0 {
			ctx.MarkFileModified(file)
			total += fileCount
		}
	}
	return total
}

func isNonEscapingNtbSetterBody(body *ast.BlockStmt, recvName string, recvVar *types.Var, info *types.Info) bool {
	if body == nil || len(body.List) == 0 {
		return false
	}
	sawReturnRecv := false
	for _, st := range body.List {
		if isBlankDiscardStmt(st) {
			continue
		}
		switch s := st.(type) {
		case *ast.AssignStmt:
			if s.Tok != token.ASSIGN || len(s.Lhs) != 1 || len(s.Rhs) != 1 {
				return false
			}
			sel, isS := stripParens(s.Lhs[0]).(*ast.SelectorExpr)
			if !isS || sel.Sel == nil || !isIdentNamed(sel.X, recvName) {
				return false
			}
			if exprMentionsVar(s.Rhs[0], recvVar, info) {
				return false
			}
		case *ast.ReturnStmt:
			if len(s.Results) != 1 || !isIdentNamed(s.Results[0], recvName) {
				return false
			}
			sawReturnRecv = true
		default:
			return false
		}
	}
	return sawReturnRecv
}

func isNonEscapingNtbTerminalBody(body *ast.BlockStmt, recvName string, recvVar *types.Var, info *types.Info) bool {
	if body == nil {
		return false
	}
	parents, _ := buildParentMap(body)
	safe := true
	ast.Inspect(body, func(n ast.Node) bool {
		if !safe || n == nil {
			return false
		}
		id, isId := n.(*ast.Ident)
		if !isId || varFromObj(info.Uses[id]) != recvVar {
			return true
		}
		par := parents[id]
		for {
			pe, isP := par.(*ast.ParenExpr)
			if !isP {
				break
			}
			par = parents[pe]
		}
		if as, isAs := par.(*ast.AssignStmt); isAs && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && isIdentNamed(as.Lhs[0], "_") {
			return true
		}
		if sel, isS := par.(*ast.SelectorExpr); isS && stripParens(sel.X) == id && sel.Sel != nil {
			if tv, okTV := info.Selections[sel]; okTV && tv.Kind() == types.FieldVal {
				if un, isU := parents[sel].(*ast.UnaryExpr); isU && un.Op == token.AND {
					safe = false
					return false
				}
				return true
			}
		}
		safe = false
		return false
	})
	return safe
}

func exprMentionsVar(expr ast.Expr, target *types.Var, info *types.Info) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if id, ok := n.(*ast.Ident); ok && varFromObj(info.Uses[id]) == target {
			found = true
			return false
		}
		return true
	})
	return found
}

func bodyContainsFuncLit(root ast.Node) bool {
	found := false
	ast.Inspect(root, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if _, isFL := n.(*ast.FuncLit); isFL {
			found = true
			return false
		}
		return true
	})
	return found
}

func cloneExprReplacingIdent(fset *token.FileSet, expr ast.Expr, oldName, newName string) (ast.Expr, bool) {
	cloned, ok := cloneExprViaParse(fset, expr)
	if !ok {
		return nil, false
	}
	if oldName == newName || oldName == "" {
		return cloned, true
	}
	parents, _ := buildParentMap(cloned)
	ast.Inspect(cloned, func(n ast.Node) bool {
		id, isId := n.(*ast.Ident)
		if !isId || id.Name != oldName {
			return true
		}
		if sel, isS := parents[id].(*ast.SelectorExpr); isS && sel.Sel == id {
			return true
		}
		if kv, isKV := parents[id].(*ast.KeyValueExpr); isKV && kv.Key == id {
			return true
		}
		id.Name = newName
		return true
	})
	return cloned, true
}

func pickFreshNtbLocalName(enclosingBlk *ast.BlockStmt, exprs []ast.Expr) string {
	used := make(map[string]bool)
	record := func(root ast.Node) {
		if root == nil {
			return
		}
		ast.Inspect(root, func(n ast.Node) bool {
			if id, isId := n.(*ast.Ident); isId && id.Name != "" {
				used[id.Name] = true
			}
			return true
		})
	}
	record(enclosingBlk)
	for _, e := range exprs {
		record(e)
	}
	if !used["ntb_"] {
		return "ntb_"
	}
	for i := 1; ; i++ {
		cand := fmt.Sprintf("ntb_%d_", i)
		if !used[cand] {
			return cand
		}
	}
}

type jsTypeEqualsMeta struct {
	rhinoPkg      *types.Package
	jsTypeNamed   *types.Named
	jsTypeIface   *types.Interface
	jsTypeBasePtr *types.Pointer
	baseEqualsFn  *types.Func
	hasToIface    bool
	hasUpcast     bool
}

func resolveJSTypeEqualsMeta(pkg *types.Package) *jsTypeEqualsMeta {
	if pkg == nil {
		return nil
	}
	var rhinoPkg *types.Package
	var findRhino func(cur *types.Package, visited map[*types.Package]bool)
	findRhino = func(cur *types.Package, visited map[*types.Package]bool) {
		if rhinoPkg != nil || cur == nil || visited[cur] {
			return
		}
		visited[cur] = true
		if cur.Name() == "rhino" && cur.Scope() != nil && cur.Scope().Lookup("JSType") != nil && cur.Scope().Lookup("JSTypeBase") != nil {
			rhinoPkg = cur
			return
		}
		for _, imp := range cur.Imports() {
			findRhino(imp, visited)
		}
	}
	findRhino(pkg, make(map[*types.Package]bool))
	if rhinoPkg == nil || rhinoPkg.Scope() == nil {
		return nil
	}
	jtTN, ok1 := rhinoPkg.Scope().Lookup("JSType").(*types.TypeName)
	jtbTN, ok2 := rhinoPkg.Scope().Lookup("JSTypeBase").(*types.TypeName)
	if !ok1 || !ok2 || jtTN == nil || jtbTN == nil {
		return nil
	}
	jtNamed := namedOrigin(jtTN.Type())
	jtbNamed := namedOrigin(jtbTN.Type())
	if jtNamed == nil || jtbNamed == nil {
		return nil
	}
	jtIface, okI := jtNamed.Underlying().(*types.Interface)
	if !okI {
		return nil
	}
	jtbPtr := types.NewPointer(jtbNamed)

	// Dynamically check whether *JSTypeBase exports EqualsJSType(JSType) bool.
	eqObj, _, _ := types.LookupFieldOrMethod(jtbPtr, true, rhinoPkg, "EqualsJSType")
	eqFn, okFn := eqObj.(*types.Func)
	if !okFn || eqFn == nil || !eqFn.Exported() {
		return nil
	}
	eqSig, okSig := eqFn.Type().(*types.Signature)
	if !okSig || eqSig.Params().Len() != 1 || eqSig.Results().Len() != 1 ||
		!types.Identical(eqSig.Params().At(0).Type(), jtNamed) ||
		!types.Identical(eqSig.Results().At(0).Type(), types.Typ[types.Bool]) {
		return nil
	}

	bEqObj, _, _ := types.LookupFieldOrMethod(jtbPtr, true, rhinoPkg, "Equals")
	bEqFn, _ := bEqObj.(*types.Func)
	_, hasToIface := rhinoPkg.Scope().Lookup("ToIface_JSType").(*types.Func)
	_, hasUpcast := rhinoPkg.Scope().Lookup("Upcast_JSType").(*types.Func)

	return &jsTypeEqualsMeta{
		rhinoPkg:      rhinoPkg,
		jsTypeNamed:   jtNamed,
		jsTypeIface:   jtIface,
		jsTypeBasePtr: jtbPtr,
		baseEqualsFn:  bEqFn,
		hasToIface:    hasToIface,
		hasUpcast:     hasUpcast,
	}
}

func (p *J2gortPeepholePass) tryRewriteJSTypeEquals(
	call *ast.CallExpr,
	file *ast.File,
	pkg *types.Package,
	facts *peepholeFileFacts,
	meta *jsTypeEqualsMeta,
	info *types.Info,
) (ast.Expr, bool) {
	if call == nil || len(call.Args) != 1 || call.Ellipsis.IsValid() || meta == nil {
		return nil, false
	}
	selExpr, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || selExpr.Sel == nil || selExpr.Sel.Name != "Equals" {
		return nil, false
	}
	sel := info.Selections[selExpr]
	if sel == nil || sel.Kind() != types.MethodVal {
		return nil, false
	}
	recvTV, ok := info.Types[selExpr.X]
	if !ok || recvTV.IsType() || recvTV.IsBuiltin() || recvTV.IsNil() || recvTV.Type == nil || isInvalidOrUntyped(recvTV.Type) {
		return nil, false
	}
	recvType := types.Unalias(recvTV.Type)
	if _, isTP := recvType.(*types.TypeParam); isTP {
		return nil, false
	}
	if !types.Implements(recvType, meta.jsTypeIface) {
		return nil, false
	}
	if !isInterfaceType(recvType) && meta.baseEqualsFn != nil && sel.Obj() != meta.baseEqualsFn {
		return nil, false
	}

	arg := stripParens(call.Args[0])
	if argCall, isAC := arg.(*ast.CallExpr); isAC && len(argCall.Args) == 1 && !argCall.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(argCall, info); isJ2 && (fnName == "PtrToAny" || fnName == "ToAny") {
			inner := stripParens(argCall.Args[0])
			if innerTV, okI := info.Types[inner]; okI && innerTV.Type != nil && !isInvalidOrUntyped(innerTV.Type) {
				ut := types.Unalias(innerTV.Type)
				if _, isTP := ut.(*types.TypeParam); !isTP && types.Implements(ut, meta.jsTypeIface) {
					arg = inner
				}
			}
		}
	}

	argTV, ok := info.Types[arg]
	if !ok || argTV.IsType() || argTV.IsBuiltin() || argTV.IsNil() || argTV.Type == nil || isInvalidOrUntyped(argTV.Type) {
		return nil, false
	}
	argType := types.Unalias(argTV.Type)
	if _, isTP := argType.(*types.TypeParam); isTP {
		return nil, false
	}
	if !types.Implements(argType, meta.jsTypeIface) {
		return nil, false
	}

	var recvBase ast.Expr
	if types.Identical(recvType, meta.jsTypeBasePtr) {
		recvBase = wrapPrimaryTarget(selExpr.X)
	} else {
		recvBase = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   wrapPrimaryTarget(selExpr.X),
				Sel: &ast.Ident{NamePos: selExpr.Sel.Pos(), Name: "GetJSTypeBase"},
			},
		}
	}

	var argJSType ast.Expr
	if isInterfaceType(argType) || (types.AssignableTo(argType, meta.jsTypeNamed) && isProvablyNonNilExpr(arg, facts, info)) {
		argJSType = arg
	} else {
		adapterName := ""
		if meta.hasToIface {
			adapterName = "ToIface_JSType"
		} else if meta.hasUpcast {
			adapterName = "Upcast_JSType"
		}
		if adapterName == "" {
			return nil, false
		}
		var adapterFun ast.Expr
		if pkg == meta.rhinoPkg {
			adapterFun = &ast.Ident{NamePos: arg.Pos(), Name: adapterName}
		} else {
			alias := findImportedPkgAliasInFile(file, meta.rhinoPkg, info)
			if alias == "" {
				return nil, false
			}
			adapterFun = &ast.SelectorExpr{
				X:   &ast.Ident{NamePos: arg.Pos(), Name: alias},
				Sel: &ast.Ident{NamePos: arg.Pos(), Name: adapterName},
			}
		}
		argJSType = &ast.CallExpr{
			Fun:  adapterFun,
			Args: []ast.Expr{arg},
		}
	}

	repl := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   wrapPrimaryTarget(recvBase),
			Sel: &ast.Ident{NamePos: selExpr.Sel.Pos(), Name: "EqualsJSType"},
		},
		Lparen: call.Lparen,
		Args:   []ast.Expr{argJSType},
		Rparen: call.Rparen,
	}
	info.Types[repl] = types.TypeAndValue{Type: types.Typ[types.Bool]}
	return repl, true
}

func findImportedPkgAliasInFile(file *ast.File, targetPkg *types.Package, info *types.Info) string {
	if file == nil || targetPkg == nil {
		return ""
	}
	wantPath := targetPkg.Path()
	for _, spec := range file.Imports {
		if spec.Path == nil {
			continue
		}
		impPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil || impPath != wantPath {
			continue
		}
		if spec.Name != nil {
			if spec.Name.Name == "_" || spec.Name.Name == "." {
				continue
			}
			return spec.Name.Name
		}
		if info != nil {
			return targetPkg.Name()
		}
	}
	return ""
}

func isUpcastCall(call *ast.CallExpr, info *types.Info) (types.Type, bool) {
	if call == nil || info == nil || len(call.Args) != 1 || call.Ellipsis.IsValid() {
		return nil, false
	}
	fn, _ := calleeFunc(call, info)
	if fn == nil {
		return nil, false
	}
	name := fn.Name()
	if !strings.HasPrefix(name, "Upcast_") || len(name) <= len("Upcast_") {
		return nil, false
	}
	origSig, ok := fn.Origin().Type().(*types.Signature)
	if !ok || origSig.Recv() != nil || origSig.Variadic() || origSig.Params().Len() != 1 || origSig.Results().Len() != 1 {
		return nil, false
	}
	if origSig.TypeParams() == nil || origSig.TypeParams().Len() != 1 {
		return nil, false
	}
	callType := info.Types[call].Type
	if !isInterfaceType(callType) {
		return nil, false
	}
	return callType, true
}

func tryMakeUpcastCall(
	asIfaceCall *ast.CallExpr,
	callType types.Type,
	expr ast.Expr,
	exprType types.Type,
	info *types.Info,
) (ast.Expr, bool) {
	if asIfaceCall == nil || callType == nil || expr == nil || exprType == nil || info == nil {
		return nil, false
	}
	utExpr := types.Unalias(exprType)
	if isEmptyInterfaceType(utExpr) {
		return nil, false
	}
	if _, isTP := utExpr.(*types.TypeParam); isTP {
		return nil, false
	}

	targetNamed, ok := types.Unalias(callType).(*types.Named)
	if !ok || targetNamed == nil || targetNamed.Obj() == nil {
		return nil, false
	}
	if targetNamed.TypeArgs() != nil && targetNamed.TypeArgs().Len() > 0 {
		return nil, false
	}
	if targetNamed.TypeParams() != nil && targetNamed.TypeParams().Len() > 0 {
		return nil, false
	}
	targetIface, ok := targetNamed.Underlying().(*types.Interface)
	if !ok || !types.Implements(utExpr, targetIface) {
		return nil, false
	}
	targetPkg := targetNamed.Obj().Pkg()
	if targetPkg == nil || targetPkg.Scope() == nil {
		return nil, false
	}
	upcastName := "Upcast_" + targetNamed.Obj().Name()
	upcastFn, ok := targetPkg.Scope().Lookup(upcastName).(*types.Func)
	if !ok || upcastFn == nil || !upcastFn.Exported() {
		return nil, false
	}
	upSig, ok := upcastFn.Type().(*types.Signature)
	if !ok || upSig.Recv() != nil || upSig.Variadic() || upSig.Params().Len() != 1 || upSig.Results().Len() != 1 ||
		upSig.TypeParams() == nil || upSig.TypeParams().Len() != 1 ||
		!types.Identical(upSig.Results().At(0).Type(), targetNamed) {
		return nil, false
	}
	constraintIface, ok := upSig.TypeParams().At(0).Constraint().Underlying().(*types.Interface)
	if !ok || !types.Satisfies(utExpr, constraintIface) {
		return nil, false
	}

	idxExpr, ok := stripParens(asIfaceCall.Fun).(*ast.IndexExpr)
	if !ok {
		return nil, false
	}
	pos := asIfaceCall.Pos()
	var funExpr ast.Expr
	switch tIdx := stripParens(idxExpr.Index).(type) {
	case *ast.Ident:
		id := &ast.Ident{NamePos: pos, Name: upcastName}
		info.Uses[id] = upcastFn
		funExpr = id
	case *ast.SelectorExpr:
		pkgId, okP := stripParens(tIdx.X).(*ast.Ident)
		if !okP {
			return nil, false
		}
		pObj, okU := info.Uses[pkgId].(*types.PkgName)
		if !okU || pObj.Imported() != targetPkg {
			return nil, false
		}
		newPkgId := &ast.Ident{NamePos: pos, Name: pkgId.Name}
		info.Uses[newPkgId] = pObj
		selId := &ast.Ident{NamePos: pos, Name: upcastName}
		info.Uses[selId] = upcastFn
		funExpr = &ast.SelectorExpr{X: newPkgId, Sel: selId}
	default:
		return nil, false
	}

	repl := &ast.CallExpr{
		Fun:    funExpr,
		Lparen: asIfaceCall.Lparen,
		Args:   []ast.Expr{expr},
		Rparen: asIfaceCall.Rparen,
	}
	info.Types[repl] = types.TypeAndValue{Type: callType}
	return repl, true
}

func (p *J2gortPeepholePass) rewriteNodeHasNonNegativeLength(ctx *PassContext) int {
	if ctx == nil || ctx.Pkg == nil || ctx.Info == nil || isJ2gortPkg(ctx.Pkg) {
		return 0
	}
	info := ctx.Info
	total := 0

	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) || !peepholeASTContainsSelector(file, "GetLength") {
			continue
		}
		parents, postOrder := buildParentMap(file)
		facts := collectPeepholeFileFacts(file, info, parents)
		consumed := make(map[ast.Node]bool)
		fileCount := 0

		for _, node := range postOrder {
			bin, ok := node.(*ast.BinaryExpr)
			if !ok || consumed[bin] {
				continue
			}
			parent := parents[bin]
			if parent == nil {
				continue
			}

			var recvExpr ast.Expr
			var rhinoPkg *types.Package
			var negate bool
			matched := false

			switch bin.Op {
			case token.GEQ:
				if isZeroInt32OrUntypedConst(bin.Y, info) {
					if rx, rPkg, okM := matchRhinoNodeGetLengthCall(bin.X, facts, info); okM {
						recvExpr, rhinoPkg, negate, matched = rx, rPkg, false, true
					}
				}
			case token.LEQ:
				if isZeroInt32OrUntypedConst(bin.X, info) {
					if rx, rPkg, okM := matchRhinoNodeGetLengthCall(bin.Y, facts, info); okM {
						recvExpr, rhinoPkg, negate, matched = rx, rPkg, false, true
					}
				}
			case token.LSS:
				if isZeroInt32OrUntypedConst(bin.Y, info) {
					if rx, rPkg, okM := matchRhinoNodeGetLengthCall(bin.X, facts, info); okM {
						recvExpr, rhinoPkg, negate, matched = rx, rPkg, true, true
					}
				}
			case token.GTR:
				if isZeroInt32OrUntypedConst(bin.X, info) {
					if rx, rPkg, okM := matchRhinoNodeGetLengthCall(bin.Y, facts, info); okM {
						recvExpr, rhinoPkg, negate, matched = rx, rPkg, true, true
					}
				}
			}
			if !matched {
				continue
			}

			nsPkg := findPeepholeNodestorePkg(file, ctx, rhinoPkg)
			if !nodestoreNodeHasNonNegativeLength(nsPkg) {
				continue
			}
			nsQual, okImp := ensurePeepholeNodestoreImport(file, nsPkg, ctx)
			if !okImp || nsQual == "" {
				continue
			}

			pos := bin.Pos()
			nsConv := &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{NamePos: pos, Name: nsQual},
					Sel: &ast.Ident{NamePos: pos, Name: "Node"},
				},
				Args: []ast.Expr{recvExpr},
			}
			hasCall := ast.Expr(&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   nsConv,
					Sel: &ast.Ident{NamePos: pos, Name: "HasNonNegativeLength"},
				},
			})

			cur := ast.Expr(bin)
			curParent := parent
			for {
				pe, isP := curParent.(*ast.ParenExpr)
				if !isP || pe.X != cur {
					break
				}
				cur = pe
				curParent = parents[pe]
			}
			if _, isES := curParent.(*ast.ExprStmt); isES {
				continue
			}

			if negate {
				if un, isNot := curParent.(*ast.UnaryExpr); isNot && un.Op == token.NOT && un.X == cur {
					if outer := parents[un]; outer != nil {
						if _, isES := outer.(*ast.ExprStmt); !isES {
							if replaceExprInParent(outer, un, hasCall, parents) {
								consumed[bin] = true
								consumed[cur] = true
								consumed[un] = true
								fileCount++
								continue
							}
						}
					}
				}
				repl := ast.Expr(&ast.UnaryExpr{OpPos: pos, Op: token.NOT, X: hasCall})
				switch parent.(type) {
				case *ast.SelectorExpr, *ast.CallExpr, *ast.IndexExpr, *ast.IndexListExpr, *ast.SliceExpr, *ast.TypeAssertExpr, *ast.StarExpr, *ast.UnaryExpr:
					repl = &ast.ParenExpr{X: repl}
				}
				if replaceExprInParent(parent, bin, repl, parents) {
					consumed[bin] = true
					fileCount++
				}
			} else {
				if un, isNot := curParent.(*ast.UnaryExpr); isNot && un.Op == token.NOT && un.X == cur && cur != bin {
					if replaceExprInParent(un, cur, hasCall, parents) {
						consumed[bin] = true
						consumed[cur] = true
						fileCount++
						continue
					}
				}
				if replaceExprInParent(parent, bin, hasCall, parents) {
					consumed[bin] = true
					fileCount++
				}
			}
		}

		if fileCount > 0 {
			ctx.MarkFileModified(file)
			total += fileCount
		}
	}

	return total
}

func isZeroInt32OrUntypedConst(expr ast.Expr, info *types.Info) bool {
	if expr == nil || info == nil || !isPureConstExpr(expr) || !isZeroConstExpr(expr, info) {
		return false
	}
	tv, ok := info.Types[stripParens(expr)]
	if !ok || tv.Type == nil {
		return false
	}
	b, ok := types.Unalias(tv.Type).(*types.Basic)
	if !ok {
		return false
	}
	return b.Info()&types.IsInteger != 0
}

func matchRhinoNodeGetLengthCall(expr ast.Expr, facts *peepholeFileFacts, info *types.Info) (ast.Expr, *types.Package, bool) {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() || len(call.Args) != 0 {
		return nil, nil, false
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || sel.Sel.Name != "GetLength" {
		return nil, nil, false
	}
	callee, _ := calleeFunc(call, info)
	if callee == nil || callee.Name() != "GetLength" {
		return nil, nil, false
	}
	sig, ok := callee.Type().(*types.Signature)
	if !ok || sig.Params().Len() != 0 || sig.Results() == nil || sig.Results().Len() != 1 ||
		!types.Identical(sig.Results().At(0).Type(), types.Typ[types.Int32]) {
		return nil, nil, false
	}
	recv := stripParens(sel.X)
	recvT := info.TypeOf(recv)
	if recvT == nil || !isNamedUint32NodeType(recvT) {
		return nil, nil, false
	}
	named, ok := types.Unalias(recvT).(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Name() != "Node" || named.Obj().Pkg() == nil ||
		named.Obj().Pkg().Name() != "rhino" {
		return nil, nil, false
	}
	if !isSideEffectFreePeepholeNodeExpr(recv, facts, info) {
		return nil, nil, false
	}
	return recv, named.Obj().Pkg(), true
}

func isSideEffectFreePeepholeNodeExpr(expr ast.Expr, facts *peepholeFileFacts, info *types.Info) bool {
	if expr == nil || info == nil {
		return false
	}
	switch e := stripParens(expr).(type) {
	case *ast.Ident:
		_, isVar := info.Uses[e].(*types.Var)
		return isVar
	case *ast.SelectorExpr:
		seln, ok := info.Selections[e]
		if !ok || seln.Kind() != types.FieldVal {
			return false
		}
		if seln.Indirect() {
			return isGuaranteedNonNilForCheckNotNull(e.X, facts, info) && isSideEffectFreeStructBaseExpr(e.X, facts, info)
		}
		return isSideEffectFreeStructBaseExpr(e.X, facts, info)
	case *ast.CallExpr:
		if len(e.Args) == 1 && !e.Ellipsis.IsValid() {
			if fnName, _, _, isJ2 := isJ2gortCall(e, info); isJ2 && (fnName == "CheckNotNullNode" || fnName == "CastNode") {
				return isSideEffectFreePeepholeNodeExpr(e.Args[0], facts, info)
			}
		}
	}
	return false
}

func isSideEffectFreeStructBaseExpr(expr ast.Expr, facts *peepholeFileFacts, info *types.Info) bool {
	if expr == nil || info == nil {
		return false
	}
	switch e := stripParens(expr).(type) {
	case *ast.Ident:
		_, isVar := info.Uses[e].(*types.Var)
		return isVar
	case *ast.SelectorExpr:
		seln, ok := info.Selections[e]
		if !ok || seln.Kind() != types.FieldVal {
			return false
		}
		if seln.Indirect() && !isGuaranteedNonNilForCheckNotNull(e.X, facts, info) {
			return false
		}
		return isSideEffectFreeStructBaseExpr(e.X, facts, info)
	}
	return false
}

func findPeepholeNodestorePkg(file *ast.File, ctx *PassContext, rhinoPkg *types.Package) *types.Package {
	if _, pkg := findNodestoreImportForFile(file, ctx); pkg != nil {
		return pkg
	}
	if pkg := findNodestorePackage(ctx, rhinoPkg); pkg != nil {
		return pkg
	}
	if ctx != nil && ctx.Importer != nil && rhinoPkg != nil {
		rp := rhinoPkg.Path()
		var cands []string
		if strings.HasSuffix(rp, "/gen/rhino") {
			cands = append(cands, strings.TrimSuffix(rp, "/gen/rhino")+"/shims/nodestore")
		}
		if strings.HasSuffix(rp, "/rhino") {
			base := strings.TrimSuffix(rp, "/rhino")
			cands = append(cands, base+"/nodestore", base+"/shims/nodestore")
		}
		for _, c := range cands {
			if p, err := ctx.Importer.Import(c); err == nil && p != nil && p.Name() == "nodestore" {
				return p
			}
		}
	}
	return nil
}

func nodestoreNodeHasNonNegativeLength(nsPkg *types.Package) bool {
	if nsPkg == nil || nsPkg.Scope() == nil {
		return false
	}
	obj := nsPkg.Scope().Lookup("Node")
	if obj == nil || obj.Type() == nil {
		return false
	}
	b, ok := obj.Type().Underlying().(*types.Basic)
	if !ok || b.Kind() != types.Uint32 {
		return false
	}
	mObj, _, _ := types.LookupFieldOrMethod(obj.Type(), false, nsPkg, "HasNonNegativeLength")
	mFn, ok := mObj.(*types.Func)
	if !ok {
		return false
	}
	sig, ok := mFn.Type().(*types.Signature)
	return ok && sig.Params().Len() == 0 && sig.Results() != nil && sig.Results().Len() == 1 &&
		types.Identical(sig.Results().At(0).Type(), types.Typ[types.Bool])
}

func ensurePeepholeNodestoreImport(file *ast.File, nsPkg *types.Package, ctx *PassContext) (string, bool) {
	if file == nil || nsPkg == nil || ctx == nil || ctx.Pkg == nsPkg {
		return "", false
	}
	if qual, existing := findNodestoreImportForFile(file, ctx); qual != "" && existing == nsPkg {
		return qual, true
	}
	alias := "nodestore"
	if (ctx.Pkg != nil && ctx.Pkg.Scope() != nil && ctx.Pkg.Scope().Lookup(alias) != nil) ||
		peepholeASTDeclaresIdent(file, alias) {
		alias = "pkg_nodestore"
		if (ctx.Pkg != nil && ctx.Pkg.Scope() != nil && ctx.Pkg.Scope().Lookup(alias) != nil) ||
			peepholeASTDeclaresIdent(file, alias) {
			return "", false
		}
	}
	return ensurePkgImportInFile(file, nsPkg.Path(), alias, ctx)
}

func peepholeASTDeclaresIdent(file *ast.File, name string) bool {
	if file == nil {
		return false
	}
	declared := false
	ast.Inspect(file, func(n ast.Node) bool {
		if declared || n == nil {
			return false
		}
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name != nil && x.Name.Name == name {
				declared = true
				return false
			}
			if peepholeFieldListHasName(x.Recv, name) || peepholeFieldListHasName(x.Type.Params, name) || peepholeFieldListHasName(x.Type.Results, name) {
				declared = true
				return false
			}
		case *ast.FuncLit:
			if x.Type != nil && (peepholeFieldListHasName(x.Type.Params, name) || peepholeFieldListHasName(x.Type.Results, name)) {
				declared = true
				return false
			}
		case *ast.ValueSpec:
			for _, id := range x.Names {
				if id.Name == name {
					declared = true
					return false
				}
			}
		case *ast.TypeSpec:
			if x.Name != nil && x.Name.Name == name {
				declared = true
				return false
			}
		case *ast.AssignStmt:
			if x.Tok == token.DEFINE {
				for _, lhs := range x.Lhs {
					if id, ok := stripParens(lhs).(*ast.Ident); ok && id.Name == name {
						declared = true
						return false
					}
				}
			}
		case *ast.RangeStmt:
			if x.Tok == token.DEFINE {
				if id, ok := x.Key.(*ast.Ident); ok && id.Name == name {
					declared = true
					return false
				}
				if id, ok := x.Value.(*ast.Ident); ok && id.Name == name {
					declared = true
					return false
				}
			}
		}
		return true
	})
	return declared
}

func peepholeFieldListHasName(fl *ast.FieldList, name string) bool {
	if fl == nil {
		return false
	}
	for _, f := range fl.List {
		for _, id := range f.Names {
			if id.Name == name {
				return true
			}
		}
	}
	return false
}

// rewriteStmtValueIIFEs unwraps statement-level zero-arg single-unnamed-result closure IIFEs
// (`var x T = func() T { ... }()`, `x := func() T { ... }()`, `x = func() T { ... }()`, and
// `return func() T { ... }()`), eliminating non-inlined closure calls from Java switch
// expressions and multi-statement value blocks while remaining 100% idempotent.
func (p *J2gortPeepholePass) rewriteStmtValueIIFEs(ctx *PassContext) int {
	if ctx == nil || ctx.Fset == nil {
		return 0
	}
	total := 0
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		if n := unwrapFileStmtValueIIFEs(file, ctx.Fset, ctx.Pkg, ctx.Info); n > 0 {
			ctx.MarkFileModified(file)
			total += n
		}
	}
	return total
}

func fileHasFuncLit(file *ast.File) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if _, ok := n.(*ast.FuncLit); ok {
			found = true
			return false
		}
		return true
	})
	return found
}

func unwrapFileStmtValueIIFEs(file *ast.File, fset *token.FileSet, pkg *types.Package, info *types.Info) int {
	if file == nil || fset == nil || !fileHasFuncLit(file) {
		return 0
	}

	var facts *peepholeFileFacts
	if info != nil {
		parents, _ := buildParentMap(file)
		facts = collectPeepholeFileFacts(file, info, parents)
	}

	usedIdents := make(map[string]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && id.Name != "" {
			usedIdents[id.Name] = true
		}
		return true
	})
	nextSwitchIdx := 1
	allocSwitchLabel := func(encLabels map[string]bool) string {
		for {
			lbl := fmt.Sprintf("iifeSwitch_%d", nextSwitchIdx)
			nextSwitchIdx++
			if !usedIdents[lbl] && !encLabels[lbl] {
				usedIdents[lbl] = true
				encLabels[lbl] = true
				return lbl
			}
		}
	}
	nextArgIdx := 1
	allocArgIdent := func() string {
		for {
			nm := fmt.Sprintf("iifeArg_%d", nextArgIdx)
			nextArgIdx++
			if !usedIdents[nm] {
				usedIdents[nm] = true
				return nm
			}
		}
	}
	synthPureIdents := make(map[string]bool)

	count := 0
	var walkFunc func(fnRecv *ast.FieldList, fnType *ast.FuncType, fnBody *ast.BlockStmt)
	walkFunc = func(fnRecv *ast.FieldList, fnType *ast.FuncType, fnBody *ast.BlockStmt) {
		if fnBody == nil {
			return
		}
		// Walk nested FuncLits bottom-up first so inner IIFEs unwrap before enclosing IIFEs.
		ast.Inspect(fnBody, func(n ast.Node) bool {
			if n == nil {
				return false
			}
			if fl, ok := n.(*ast.FuncLit); ok {
				walkFunc(nil, fl.Type, fl.Body)
				return false
			}
			return true
		})

		var fnRecvVar *types.Var
		if info != nil && fnRecv != nil && len(fnRecv.List) == 1 && len(fnRecv.List[0].Names) == 1 {
			if rId := fnRecv.List[0].Names[0]; rId != nil && rId.Name != "" && rId.Name != "_" {
				if rv, okV := info.Defs[rId].(*types.Var); okV && rv != nil {
					if _, isPtr := types.Unalias(rv.Type()).(*types.Pointer); isPtr && !isVarReassignedInBody(fnBody, rv, info) {
						fnRecvVar = rv
					}
				}
			}
		}

		encLabels := collectDirectLabelsInBlock(fnBody)
		encHasDefer := hasDirectDeferInBlock(fnBody)
		encHasGoto := hasDirectGotoInBlock(fnBody)
		canUnwrapReturn := !encHasDefer && fnType != nil &&
			fnType.Results != nil && len(fnType.Results.List) == 1 &&
			len(fnType.Results.List[0].Names) == 0 && fnType.Results.List[0].Type != nil

		var rewriteStmtList func(stmts []ast.Stmt, isTopBody bool) ([]ast.Stmt, bool)
		var walkNodeStmts func(n ast.Node, isTopBody bool)

		walkNodeStmts = func(n ast.Node, isTopBody bool) {
			if n == nil {
				return
			}
			switch x := n.(type) {
			case *ast.FuncLit:
				return
			case *ast.BlockStmt:
				for _, st := range x.List {
					walkNodeStmts(st, false)
				}
				if upd, changed := rewriteStmtList(x.List, isTopBody); changed {
					x.List = upd
				}
			case *ast.CaseClause:
				for _, st := range x.Body {
					walkNodeStmts(st, false)
				}
				if upd, changed := rewriteStmtList(x.Body, false); changed {
					x.Body = upd
				}
			case *ast.CommClause:
				for _, st := range x.Body {
					walkNodeStmts(st, false)
				}
				if upd, changed := rewriteStmtList(x.Body, false); changed {
					x.Body = upd
				}
			default:
				ast.Inspect(x, func(child ast.Node) bool {
					if child == nil || child == x {
						return true
					}
					switch child.(type) {
					case *ast.FuncLit:
						return false
					case *ast.BlockStmt, *ast.CaseClause, *ast.CommClause:
						walkNodeStmts(child, false)
						return false
					}
					return true
				})
			}
		}

		topRecvNonNil := false
		rewriteStmtList = func(stmts []ast.Stmt, isTopBody bool) ([]ast.Stmt, bool) {
			var out []ast.Stmt
			changed := false
			for i, st := range stmts {
				recvNonNilHere := topRecvNonNil
				repl, okR := tryUnwrapStmtValueIIFE(
					st, fnType, canUnwrapReturn, encLabels, allocSwitchLabel, fset, pkg, info,
				)
				nSites := 1
				if !okR {
					repl, nSites = tryUnwrapStmtCallArgIIFEs(
						st, encHasGoto, encLabels, allocSwitchLabel, allocArgIdent,
						synthPureIdents, fnRecvVar, recvNonNilHere, facts, fset, pkg, info,
					)
					okR = nSites > 0
				}
				if isTopBody && !topRecvNonNil && fnRecvVar != nil && stmtUnconditionallyDerefsVarFirst(st, fnRecvVar, info) {
					topRecvNonNil = true
				}
				if !okR {
					if changed {
						out = append(out, st)
					}
					continue
				}
				if !changed {
					out = make([]ast.Stmt, 0, len(stmts)+len(repl))
					out = append(out, stmts[:i]...)
					changed = true
				}
				out = append(out, repl...)
				count += nSites
			}
			return out, changed
		}

		walkNodeStmts(fnBody, true)
	}

	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Body != nil {
			walkFunc(fn.Recv, fn.Type, fn.Body)
		}
	}
	return count
}

func tryUnwrapStmtValueIIFE(
	stmt ast.Stmt,
	encFnType *ast.FuncType,
	canUnwrapReturn bool,
	encLabels map[string]bool,
	allocSwitchLabel func(map[string]bool) string,
	fset *token.FileSet,
	pkg *types.Package,
	info *types.Info,
) ([]ast.Stmt, bool) {
	switch s := stmt.(type) {
	case *ast.DeclStmt:
		gd, ok := s.Decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR || len(gd.Specs) != 1 {
			return nil, false
		}
		vs, ok := gd.Specs[0].(*ast.ValueSpec)
		if !ok || len(vs.Names) != 1 || len(vs.Values) != 1 {
			return nil, false
		}
		lhsIdent := vs.Names[0]
		if lhsIdent == nil || lhsIdent.Name == "" || lhsIdent.Name == "_" {
			return nil, false
		}
		fl, bodyLabels, ok := matchUnwrappableValueIIFECall(vs.Values[0], encLabels)
		if !ok || iifeASTContainsIdentNamed(fl.Body, lhsIdent.Name) {
			return nil, false
		}
		retTypeAST := fl.Type.Results.List[0].Type
		if vs.Type != nil {
			if !iifeTypesIdentical(vs.Type, retTypeAST, lhsIdent, true, fset, info) {
				return nil, false
			}
		} else {
			if info != nil {
				if tv, okT := info.Types[retTypeAST]; !okT || tv.Type == nil {
					return nil, false
				}
			}
			clonedType, okC := cloneIIFETypeExpr(fset, retTypeAST)
			if !okC {
				return nil, false
			}
			vs.Type = clonedType
		}
		for lbl := range bodyLabels {
			encLabels[lbl] = true
		}
		vs.Values = nil
		bodyStmt := buildUnwrappedIIFEAssignBlock(fl.Body, lhsIdent.Name, encLabels, allocSwitchLabel)
		return []ast.Stmt{s, bodyStmt}, true

	case *ast.AssignStmt:
		if len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			return nil, false
		}
		lhsIdent, ok := s.Lhs[0].(*ast.Ident)
		if !ok || lhsIdent == nil || lhsIdent.Name == "" || lhsIdent.Name == "_" {
			return nil, false
		}
		switch s.Tok {
		case token.DEFINE:
			if info != nil && info.Defs[lhsIdent] == nil {
				return nil, false
			}
			fl, bodyLabels, ok := matchUnwrappableValueIIFECall(s.Rhs[0], encLabels)
			if !ok || iifeASTContainsIdentNamed(fl.Body, lhsIdent.Name) {
				return nil, false
			}
			retTypeAST := fl.Type.Results.List[0].Type
			if info != nil {
				if tv, okT := info.Types[retTypeAST]; !okT || tv.Type == nil {
					return nil, false
				}
			}
			clonedType, okC := cloneIIFETypeExpr(fset, retTypeAST)
			if !okC {
				return nil, false
			}
			for lbl := range bodyLabels {
				encLabels[lbl] = true
			}
			varDecl := &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{lhsIdent},
							Type:  clonedType,
						},
					},
				},
			}
			bodyStmt := buildUnwrappedIIFEAssignBlock(fl.Body, lhsIdent.Name, encLabels, allocSwitchLabel)
			return []ast.Stmt{varDecl, bodyStmt}, true

		case token.ASSIGN:
			if info == nil {
				return nil, false
			}
			lhsVar, okV := info.Uses[lhsIdent].(*types.Var)
			if !okV || lhsVar == nil || lhsVar.IsField() || lhsVar.Parent() == nil {
				return nil, false
			}
			if pkg != nil && lhsVar.Parent() == pkg.Scope() {
				return nil, false
			}
			fl, bodyLabels, ok := matchUnwrappableValueIIFECall(s.Rhs[0], encLabels)
			if !ok {
				return nil, false
			}
			retTV, okT := info.Types[fl.Type.Results.List[0].Type]
			if !okT || retTV.Type == nil || !types.Identical(lhsVar.Type(), retTV.Type) {
				return nil, false
			}
			if !iifeLocalVarUnshadowedInBody(fl.Body, lhsIdent.Name, lhsVar, info) {
				return nil, false
			}
			for lbl := range bodyLabels {
				encLabels[lbl] = true
			}
			bodyStmt := buildUnwrappedIIFEAssignBlock(fl.Body, lhsIdent.Name, encLabels, allocSwitchLabel)
			return []ast.Stmt{bodyStmt}, true
		}

	case *ast.ReturnStmt:
		if !canUnwrapReturn || len(s.Results) != 1 || encFnType == nil {
			return nil, false
		}
		fl, bodyLabels, ok := matchUnwrappableValueIIFECall(s.Results[0], encLabels)
		if !ok {
			return nil, false
		}
		encRetAST := encFnType.Results.List[0].Type
		flRetAST := fl.Type.Results.List[0].Type
		if !iifeTypesIdentical(encRetAST, flRetAST, nil, false, fset, info) {
			return nil, false
		}
		for lbl := range bodyLabels {
			encLabels[lbl] = true
		}
		return []ast.Stmt{&ast.BlockStmt{List: fl.Body.List}}, true
	}
	return nil, false
}

func matchUnwrappableValueIIFECall(expr ast.Expr, encLabels map[string]bool) (*ast.FuncLit, map[string]bool, bool) {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call == nil || len(call.Args) != 0 || call.Ellipsis.IsValid() {
		return nil, nil, false
	}
	fl, ok := stripParens(call.Fun).(*ast.FuncLit)
	if !ok || fl == nil || fl.Type == nil || fl.Body == nil || len(fl.Body.List) == 0 {
		return nil, nil, false
	}
	if fl.Type.TypeParams != nil && len(fl.Type.TypeParams.List) > 0 {
		return nil, nil, false
	}
	if fl.Type.Params != nil && fl.Type.Params.NumFields() != 0 {
		return nil, nil, false
	}
	if fl.Type.Results == nil || len(fl.Type.Results.List) != 1 ||
		len(fl.Type.Results.List[0].Names) != 0 || fl.Type.Results.List[0].Type == nil {
		return nil, nil, false
	}

	safe := true
	hasReturn := false
	var bodyLabels map[string]bool
	ast.Inspect(fl.Body, func(n ast.Node) bool {
		if !safe || n == nil {
			return false
		}
		if nested, isFL := n.(*ast.FuncLit); isFL && nested != fl {
			return false
		}
		switch x := n.(type) {
		case *ast.DeferStmt, *ast.GoStmt:
			safe = false
			return false
		case *ast.BranchStmt:
			if x.Tok == token.GOTO {
				safe = false
				return false
			}
		case *ast.LabeledStmt:
			if x.Label == nil || x.Label.Name == "" || encLabels[x.Label.Name] || bodyLabels[x.Label.Name] {
				safe = false
				return false
			}
			if bodyLabels == nil {
				bodyLabels = make(map[string]bool)
			}
			bodyLabels[x.Label.Name] = true
		case *ast.ReturnStmt:
			if len(x.Results) != 1 {
				safe = false
				return false
			}
			hasReturn = true
		}
		return true
	})
	if !safe || !hasReturn {
		return nil, nil, false
	}
	return fl, bodyLabels, true
}

func buildUnwrappedIIFEAssignBlock(
	body *ast.BlockStmt,
	lhsName string,
	encLabels map[string]bool,
	allocSwitchLabel func(map[string]bool) string,
) ast.Stmt {
	if iifeStmtsAllReturnsTailOnly(body.List) {
		rewriteIIFEDirectReturns(body, lhsName, "")
		return &ast.BlockStmt{List: body.List}
	}
	lbl := allocSwitchLabel(encLabels)
	rewriteIIFEDirectReturns(body, lhsName, lbl)
	return &ast.LabeledStmt{
		Label: &ast.Ident{Name: lbl},
		Stmt: &ast.SwitchStmt{
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.CaseClause{
						Body: body.List,
					},
				},
			},
		},
	}
}

func rewriteIIFEDirectReturns(body *ast.BlockStmt, lhsName, breakLabel string) {
	var rewriteList func(stmts []ast.Stmt) []ast.Stmt
	var walk func(n ast.Node)

	walk = func(n ast.Node) {
		if n == nil {
			return
		}
		switch x := n.(type) {
		case *ast.FuncLit:
			return
		case *ast.BlockStmt:
			for _, st := range x.List {
				walk(st)
			}
			x.List = rewriteList(x.List)
		case *ast.CaseClause:
			for _, st := range x.Body {
				walk(st)
			}
			x.Body = rewriteList(x.Body)
		case *ast.CommClause:
			for _, st := range x.Body {
				walk(st)
			}
			x.Body = rewriteList(x.Body)
		default:
			ast.Inspect(x, func(child ast.Node) bool {
				if child == nil || child == x {
					return true
				}
				switch child.(type) {
				case *ast.FuncLit:
					return false
				case *ast.BlockStmt, *ast.CaseClause, *ast.CommClause:
					walk(child)
					return false
				}
				return true
			})
		}
	}

	rewriteList = func(stmts []ast.Stmt) []ast.Stmt {
		var out []ast.Stmt
		changed := false
		for i, st := range stmts {
			ret, isRet := st.(*ast.ReturnStmt)
			if !isRet {
				if changed {
					out = append(out, st)
				}
				continue
			}
			if !changed {
				out = make([]ast.Stmt, 0, len(stmts)+2)
				out = append(out, stmts[:i]...)
				changed = true
			}
			assign := &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{NamePos: ret.Return, Name: lhsName}},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{ret.Results[0]},
			}
			out = append(out, assign)
			if breakLabel != "" {
				out = append(out, &ast.BranchStmt{
					TokPos: ret.Return,
					Tok:    token.BREAK,
					Label:  &ast.Ident{Name: breakLabel},
				})
			}
		}
		if changed {
			return out
		}
		return stmts
	}

	walk(body)
}

func iifeStmtsAllReturnsTailOnly(stmts []ast.Stmt) bool {
	if len(stmts) == 0 {
		return true
	}
	for i := 0; i < len(stmts)-1; i++ {
		if iifeNodeHasDirectReturn(stmts[i]) {
			return false
		}
	}
	last := stmts[len(stmts)-1]
	if !iifeNodeHasDirectReturn(last) {
		return true
	}
	switch s := last.(type) {
	case *ast.ReturnStmt:
		return len(s.Results) == 1
	case *ast.BlockStmt:
		return iifeStmtsAllReturnsTailOnly(s.List)
	case *ast.IfStmt:
		if !iifeStmtsAllReturnsTailOnly(s.Body.List) {
			return false
		}
		if s.Else == nil {
			return true
		}
		switch el := s.Else.(type) {
		case *ast.BlockStmt:
			return iifeStmtsAllReturnsTailOnly(el.List)
		case *ast.IfStmt:
			return iifeStmtsAllReturnsTailOnly([]ast.Stmt{el})
		default:
			return false
		}
	case *ast.SwitchStmt:
		if s.Body == nil || switchHasDirectFallthrough(s.Body) {
			return false
		}
		for _, raw := range s.Body.List {
			cc, ok := raw.(*ast.CaseClause)
			if !ok || !iifeStmtsAllReturnsTailOnly(cc.Body) {
				return false
			}
		}
		return true
	case *ast.TypeSwitchStmt:
		if s.Body == nil {
			return false
		}
		for _, raw := range s.Body.List {
			cc, ok := raw.(*ast.CaseClause)
			if !ok || !iifeStmtsAllReturnsTailOnly(cc.Body) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func switchHasDirectFallthrough(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found || n == nil || n == body {
			return n == body
		}
		switch x := n.(type) {
		case *ast.FuncLit, *ast.SwitchStmt:
			return false
		case *ast.BranchStmt:
			if x.Tok == token.FALLTHROUGH {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

func iifeNodeHasDirectReturn(node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		switch n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.ReturnStmt:
			found = true
			return false
		}
		return true
	})
	return found
}

func collectDirectLabelsInBlock(body *ast.BlockStmt) map[string]bool {
	labels := make(map[string]bool)
	ast.Inspect(body, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch x := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.LabeledStmt:
			if x.Label != nil && x.Label.Name != "" {
				labels[x.Label.Name] = true
			}
		}
		return true
	})
	return labels
}

func hasDirectDeferInBlock(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		switch n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.DeferStmt:
			found = true
			return false
		}
		return true
	})
	return found
}

func iifeASTContainsIdentNamed(node ast.Node, name string) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if id, ok := n.(*ast.Ident); ok && id.Name == name {
			found = true
			return false
		}
		return true
	})
	return found
}

func iifeLocalVarUnshadowedInBody(body *ast.BlockStmt, name string, targetVar *types.Var, info *types.Info) bool {
	ok := true
	ast.Inspect(body, func(n ast.Node) bool {
		if !ok || n == nil {
			return false
		}
		if id, isId := n.(*ast.Ident); isId && id.Name == name {
			if info.Defs[id] != nil || info.Uses[id] != targetVar {
				ok = false
				return false
			}
		}
		return true
	})
	return ok
}

func cloneIIFETypeExpr(fset *token.FileSet, typeExpr ast.Expr) (ast.Expr, bool) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, typeExpr); err != nil {
		return nil, false
	}
	cloned, err := parser.ParseExpr(buf.String())
	if err != nil {
		return nil, false
	}
	return cloned, true
}

func iifeTypesIdentical(lhsTypeAST, rhsTypeAST ast.Expr, lhsDefIdent *ast.Ident, useDef bool, fset *token.FileSet, info *types.Info) bool {
	if info != nil {
		var lhsT types.Type
		if useDef && lhsDefIdent != nil {
			if obj := info.Defs[lhsDefIdent]; obj != nil {
				lhsT = obj.Type()
			}
		}
		if lhsT == nil && lhsTypeAST != nil {
			lhsT = info.Types[lhsTypeAST].Type
		}
		rhsT := info.Types[rhsTypeAST].Type
		if lhsT != nil && rhsT != nil {
			return types.Identical(lhsT, rhsT)
		}
		return false
	}
	var b1, b2 bytes.Buffer
	if err := printer.Fprint(&b1, fset, lhsTypeAST); err != nil {
		return false
	}
	if err := printer.Fprint(&b2, fset, rhsTypeAST); err != nil {
		return false
	}
	return b1.String() == b2.String()
}

func hasDirectGotoInBlock(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		switch x := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.BranchStmt:
			if x.Tok == token.GOTO {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

func isDirectFuncLitCallExpr(expr ast.Expr) bool {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call == nil {
		return false
	}
	_, isFL := stripParens(call.Fun).(*ast.FuncLit)
	return isFL
}

type iifePrefixCheckState struct {
	vars           map[*types.Var]bool
	addrTakenNames map[string]bool
	hasFieldAccess bool
}

func newIIFEPrefixCheckState() *iifePrefixCheckState {
	return &iifePrefixCheckState{
		vars:           make(map[*types.Var]bool),
		addrTakenNames: make(map[string]bool),
	}
}

func (s *iifePrefixCheckState) clone() *iifePrefixCheckState {
	c := &iifePrefixCheckState{
		vars:           make(map[*types.Var]bool, len(s.vars)),
		addrTakenNames: make(map[string]bool, len(s.addrTakenNames)),
		hasFieldAccess: s.hasFieldAccess,
	}
	for v := range s.vars {
		c.vars[v] = true
	}
	for nm := range s.addrTakenNames {
		c.addrTakenNames[nm] = true
	}
	return c
}

func tryUnwrapStmtCallArgIIFEs(
	stmt ast.Stmt,
	encHasGoto bool,
	encLabels map[string]bool,
	allocSwitchLabel func(map[string]bool) string,
	allocArgIdent func() string,
	synthPureIdents map[string]bool,
	fnRecvVar *types.Var,
	recvNonNilHere bool,
	facts *peepholeFileFacts,
	fset *token.FileSet,
	pkg *types.Package,
	info *types.Info,
) ([]ast.Stmt, int) {
	if info == nil || fset == nil {
		return nil, 0
	}

	var rootCall *ast.CallExpr
	baseState := newIIFEPrefixCheckState()
	canScopeWrap := false

	switch s := stmt.(type) {
	case *ast.ExprStmt:
		c, ok := stripParens(s.X).(*ast.CallExpr)
		if !ok || isDirectFuncLitCallExpr(s.X) {
			return nil, 0
		}
		rootCall = c
		canScopeWrap = true

	case *ast.DeclStmt:
		if encHasGoto {
			return nil, 0
		}
		gd, ok := s.Decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR || len(gd.Specs) != 1 {
			return nil, 0
		}
		vs, ok := gd.Specs[0].(*ast.ValueSpec)
		if !ok || len(vs.Names) != 1 || len(vs.Values) != 1 {
			return nil, 0
		}
		c, ok := stripParens(vs.Values[0]).(*ast.CallExpr)
		if !ok || isDirectFuncLitCallExpr(vs.Values[0]) {
			return nil, 0
		}
		rootCall = c

	case *ast.AssignStmt:
		if len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			return nil, 0
		}
		c, ok := stripParens(s.Rhs[0]).(*ast.CallExpr)
		if !ok || isDirectFuncLitCallExpr(s.Rhs[0]) {
			return nil, 0
		}
		switch s.Tok {
		case token.DEFINE:
			if encHasGoto {
				return nil, 0
			}
			lhsId, okId := s.Lhs[0].(*ast.Ident)
			if !okId || lhsId == nil || lhsId.Name == "" {
				return nil, 0
			}
			rootCall = c
		case token.ASSIGN:
			if !isSafeIIFEAssignLhs(s.Lhs[0], baseState, fnRecvVar, recvNonNilHere, facts, pkg, info) {
				return nil, 0
			}
			rootCall = c
			canScopeWrap = true
		default:
			return nil, 0
		}

	case *ast.ReturnStmt:
		if len(s.Results) != 1 {
			return nil, 0
		}
		c, ok := stripParens(s.Results[0]).(*ast.CallExpr)
		if !ok || isDirectFuncLitCallExpr(s.Results[0]) {
			return nil, 0
		}
		rootCall = c
		canScopeWrap = true

	default:
		return nil, 0
	}

	if encHasGoto && !canScopeWrap {
		return nil, 0
	}

	var preStmts []ast.Stmt
	unwrapped := 0
	for {
		curState := baseState.clone()
		slotPtr, fl, bodyLabels, clonedType, okM := findFirstUnwrappableCallArgIIFE(
			rootCall, curState, encLabels, synthPureIdents, fnRecvVar, recvNonNilHere, facts, fset, pkg, info,
		)
		if !okM {
			break
		}
		for lbl := range bodyLabels {
			encLabels[lbl] = true
		}
		tmpName := allocArgIdent()
		synthPureIdents[tmpName] = true
		argPos := (*slotPtr).Pos()
		varDecl := &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{{NamePos: argPos, Name: tmpName}},
						Type:  clonedType,
					},
				},
			},
		}
		bodyStmt := buildUnwrappedIIFEAssignBlock(fl.Body, tmpName, encLabels, allocSwitchLabel)
		*slotPtr = &ast.Ident{NamePos: argPos, Name: tmpName}
		preStmts = append(preStmts, varDecl, bodyStmt)
		unwrapped++
	}

	if unwrapped == 0 {
		return nil, 0
	}
	if encHasGoto {
		return []ast.Stmt{&ast.BlockStmt{List: append(preStmts, stmt)}}, unwrapped
	}
	return append(preStmts, stmt), unwrapped
}

func isSafeIIFEAssignLhs(
	lhs ast.Expr,
	state *iifePrefixCheckState,
	fnRecvVar *types.Var,
	recvNonNilHere bool,
	facts *peepholeFileFacts,
	pkg *types.Package,
	info *types.Info,
) bool {
	switch e := stripParens(lhs).(type) {
	case *ast.Ident:
		if e.Name == "_" {
			return true
		}
		v, ok := info.Uses[e].(*types.Var)
		if !ok || v == nil || v.IsField() || v.Parent() == nil {
			return false
		}
		if pkg != nil && v.Parent() == pkg.Scope() {
			return false
		}
		state.vars[v] = true
		return true

	case *ast.SelectorExpr:
		sel := info.Selections[e]
		if sel == nil || sel.Kind() != types.FieldVal || len(sel.Index()) != 1 {
			return false
		}
		baseId, ok := stripParens(e.X).(*ast.Ident)
		if !ok || !isIIFEIdentProvenNonNilLocal(baseId, state, fnRecvVar, recvNonNilHere, facts, pkg, info) {
			return false
		}
		return true
	}
	return false
}

func findFirstUnwrappableCallArgIIFE(
	call *ast.CallExpr,
	state *iifePrefixCheckState,
	encLabels map[string]bool,
	synthPureIdents map[string]bool,
	fnRecvVar *types.Var,
	recvNonNilHere bool,
	facts *peepholeFileFacts,
	fset *token.FileSet,
	pkg *types.Package,
	info *types.Info,
) (*ast.Expr, *ast.FuncLit, map[string]bool, ast.Expr, bool) {
	if call == nil || call.Ellipsis.IsValid() || len(call.Args) == 0 {
		return nil, nil, nil, nil, false
	}
	if !isSafeCalleeBeforeIIFEArg(call.Fun, state, fnRecvVar, recvNonNilHere, facts, pkg, info) {
		return nil, nil, nil, nil, false
	}

	for i := range call.Args {
		arg := call.Args[i]
		if fl, bodyLabels, ok := matchUnwrappableValueIIFECall(arg, encLabels); ok {
			retTypeAST := fl.Type.Results.List[0].Type
			if tv, okT := info.Types[retTypeAST]; !okT || tv.Type == nil || isInvalidOrUntyped(tv.Type) {
				return nil, nil, nil, nil, false
			}
			if iifeBodyMayMutatePrefix(fl.Body, state, info) {
				return nil, nil, nil, nil, false
			}
			clonedType, okC := cloneIIFETypeExpr(fset, retTypeAST)
			if !okC {
				return nil, nil, nil, nil, false
			}
			return &call.Args[i], fl, bodyLabels, clonedType, true
		}

		if innerCall, isC := stripParens(arg).(*ast.CallExpr); isC && !isDirectFuncLitCallExpr(arg) {
			subState := state.clone()
			if slot, fl, bodyLabels, clonedType, okSub := findFirstUnwrappableCallArgIIFE(
				innerCall, subState, encLabels, synthPureIdents, fnRecvVar, recvNonNilHere, facts, fset, pkg, info,
			); okSub {
				return slot, fl, bodyLabels, clonedType, true
			}
		}

		if !isSafePureArgBeforeIIFE(arg, state, synthPureIdents, fnRecvVar, recvNonNilHere, facts, pkg, info) {
			return nil, nil, nil, nil, false
		}
	}
	return nil, nil, nil, nil, false
}

func stripIndexExprs(expr ast.Expr) ast.Expr {
	cur := stripParens(expr)
	for {
		switch x := cur.(type) {
		case *ast.IndexExpr:
			cur = stripParens(x.X)
		case *ast.IndexListExpr:
			cur = stripParens(x.X)
		default:
			return cur
		}
	}
}

func isSafeCalleeBeforeIIFEArg(
	funExpr ast.Expr,
	state *iifePrefixCheckState,
	fnRecvVar *types.Var,
	recvNonNilHere bool,
	facts *peepholeFileFacts,
	pkg *types.Package,
	info *types.Info,
) bool {
	base := stripIndexExprs(funExpr)
	switch f := base.(type) {
	case *ast.Ident:
		switch obj := info.Uses[f].(type) {
		case *types.Func, *types.Builtin, *types.TypeName:
			return true
		case *types.Var:
			if obj == nil || obj.IsField() || obj.Parent() == nil || (pkg != nil && obj.Parent() == pkg.Scope()) {
				return false
			}
			state.vars[obj] = true
			return true
		}
		return false

	case *ast.SelectorExpr:
		if pkgId, ok := stripParens(f.X).(*ast.Ident); ok {
			if _, isPkg := info.Uses[pkgId].(*types.PkgName); isPkg {
				switch info.Uses[f.Sel].(type) {
				case *types.Func, *types.TypeName:
					return true
				}
				return false
			}
		}
		sel := info.Selections[f]
		if sel == nil || sel.Kind() != types.MethodVal {
			return false
		}
		mFn, okF := sel.Obj().(*types.Func)
		if !okF || mFn == nil {
			return false
		}
		sig, okS := mFn.Type().(*types.Signature)
		if !okS || sig.Recv() == nil {
			return false
		}
		_, isPtrRecv := types.Unalias(sig.Recv().Type()).(*types.Pointer)
		recv := stripParens(f.X)
		recvTV, okT := info.Types[recv]
		if !okT || recvTV.Type == nil || isInvalidOrUntyped(recvTV.Type) {
			return false
		}
		ut := types.Unalias(recvTV.Type)

		switch r := recv.(type) {
		case *ast.Ident:
			rv, okV := info.Uses[r].(*types.Var)
			if !okV || rv == nil || rv.IsField() || rv.Parent() == nil || (pkg != nil && rv.Parent() == pkg.Scope()) {
				return false
			}
			switch ut.(type) {
			case *types.Pointer:
				if isPtrRecv && len(sel.Index()) == 1 {
					state.vars[rv] = true
					return true
				}
				return isIIFEIdentProvenNonNilLocal(r, state, fnRecvVar, recvNonNilHere, facts, pkg, info) && !sel.Indirect()
			case *types.Struct, *types.Basic, *types.Slice:
				if len(sel.Index()) == 1 {
					state.vars[rv] = true
					return true
				}
			case *types.Interface:
				return isIIFEIdentProvenNonNilLocal(r, state, fnRecvVar, recvNonNilHere, facts, pkg, info)
			}

		case *ast.SelectorExpr:
			fSel := info.Selections[r]
			if fSel == nil || fSel.Kind() != types.FieldVal || len(fSel.Index()) != 1 {
				return false
			}
			baseId, okB := stripParens(r.X).(*ast.Ident)
			if !okB || !isIIFEIdentProvenNonNilLocal(baseId, state, fnRecvVar, recvNonNilHere, facts, pkg, info) {
				return false
			}
			switch ut.(type) {
			case *types.Pointer:
				if isPtrRecv && len(sel.Index()) == 1 {
					state.hasFieldAccess = true
					return true
				}
			case *types.Struct, *types.Basic, *types.Slice:
				if len(sel.Index()) == 1 {
					state.hasFieldAccess = true
					return true
				}
			}
		}
	}
	return false
}

func isIIFEIdentProvenNonNilLocal(
	id *ast.Ident,
	state *iifePrefixCheckState,
	fnRecvVar *types.Var,
	recvNonNilHere bool,
	facts *peepholeFileFacts,
	pkg *types.Package,
	info *types.Info,
) bool {
	if id == nil {
		return false
	}
	v, ok := info.Uses[id].(*types.Var)
	if !ok || v == nil || v.IsField() || v.Parent() == nil || (pkg != nil && v.Parent() == pkg.Scope()) {
		return false
	}
	nonNil := (recvNonNilHere && v == fnRecvVar) || (facts != nil && isProvablyNonNilExpr(id, facts, info))
	if !nonNil {
		return false
	}
	state.vars[v] = true
	return true
}

func isSafePureArgBeforeIIFE(
	expr ast.Expr,
	state *iifePrefixCheckState,
	synthPureIdents map[string]bool,
	fnRecvVar *types.Var,
	recvNonNilHere bool,
	facts *peepholeFileFacts,
	pkg *types.Package,
	info *types.Info,
) bool {
	switch e := stripParens(expr).(type) {
	case *ast.BasicLit:
		return true

	case *ast.Ident:
		if synthPureIdents[e.Name] {
			return true
		}
		switch obj := info.Uses[e].(type) {
		case *types.Const, *types.Nil:
			return true
		case *types.Var:
			if obj == nil || obj.IsField() || obj.Parent() == nil || (pkg != nil && obj.Parent() == pkg.Scope()) {
				return false
			}
			state.vars[obj] = true
			return true
		}
		return false

	case *ast.UnaryExpr:
		switch e.Op {
		case token.ADD, token.SUB, token.NOT, token.XOR:
			return isSafePureArgBeforeIIFE(e.X, state, synthPureIdents, fnRecvVar, recvNonNilHere, facts, pkg, info)
		case token.AND:
			id, ok := stripParens(e.X).(*ast.Ident)
			if !ok || id == nil {
				return false
			}
			v, okV := info.Uses[id].(*types.Var)
			if !okV || v == nil || v.IsField() || v.Parent() == nil || (pkg != nil && v.Parent() == pkg.Scope()) {
				return false
			}
			state.vars[v] = true
			state.addrTakenNames[id.Name] = true
			return true
		}
		return false

	case *ast.BinaryExpr:
		switch e.Op {
		case token.QUO, token.REM:
			return false
		}
		return isSafePureArgBeforeIIFE(e.X, state, synthPureIdents, fnRecvVar, recvNonNilHere, facts, pkg, info) &&
			isSafePureArgBeforeIIFE(e.Y, state, synthPureIdents, fnRecvVar, recvNonNilHere, facts, pkg, info)

	case *ast.SelectorExpr:
		if pkgId, ok := stripParens(e.X).(*ast.Ident); ok {
			if _, isPkg := info.Uses[pkgId].(*types.PkgName); isPkg {
				_, isConst := info.Uses[e.Sel].(*types.Const)
				return isConst
			}
		}
		sel := info.Selections[e]
		if sel == nil || sel.Kind() != types.FieldVal || len(sel.Index()) != 1 {
			return false
		}
		baseId, ok := stripParens(e.X).(*ast.Ident)
		if !ok || !isIIFEIdentProvenNonNilLocal(baseId, state, fnRecvVar, recvNonNilHere, facts, pkg, info) {
			return false
		}
		state.hasFieldAccess = true
		return true

	case *ast.CallExpr:
		if len(e.Args) != 1 || e.Ellipsis.IsValid() {
			return false
		}
		funTV, ok := info.Types[e.Fun]
		if !ok || !funTV.IsType() || funTV.Type == nil || isInvalidOrUntyped(funTV.Type) {
			return false
		}
		switch types.Unalias(funTV.Type).Underlying().(type) {
		case *types.Basic:
			return isSafePureArgBeforeIIFE(e.Args[0], state, synthPureIdents, fnRecvVar, recvNonNilHere, facts, pkg, info)
		case *types.Pointer:
			argTV, okA := info.Types[e.Args[0]]
			if !okA || argTV.Type == nil {
				return false
			}
			if _, isSlice := types.Unalias(argTV.Type).Underlying().(*types.Slice); isSlice {
				return false
			}
			return isSafePureArgBeforeIIFE(e.Args[0], state, synthPureIdents, fnRecvVar, recvNonNilHere, facts, pkg, info)
		}
	}
	return false
}

func iifeBodyMayMutatePrefix(body *ast.BlockStmt, state *iifePrefixCheckState, info *types.Info) bool {
	if len(state.vars) == 0 && len(state.addrTakenNames) == 0 && !state.hasFieldAccess {
		return false
	}
	unsafe := false
	ast.Inspect(body, func(n ast.Node) bool {
		if unsafe || n == nil {
			return false
		}
		if id, isId := n.(*ast.Ident); isId && state.addrTakenNames[id.Name] {
			unsafe = true
			return false
		}
		switch x := n.(type) {
		case *ast.AssignStmt:
			for _, lhs := range x.Lhs {
				if id, ok := stripParens(lhs).(*ast.Ident); ok {
					if x.Tok == token.DEFINE && info.Defs[id] != nil {
						continue
					}
					if uv, okV := info.Uses[id].(*types.Var); okV && state.vars[uv] {
						unsafe = true
						return false
					}
				} else if state.hasFieldAccess {
					unsafe = true
					return false
				}
			}
		case *ast.IncDecStmt:
			if id, ok := stripParens(x.X).(*ast.Ident); ok {
				if uv, okV := info.Uses[id].(*types.Var); okV && state.vars[uv] {
					unsafe = true
					return false
				}
			} else if state.hasFieldAccess {
				unsafe = true
				return false
			}
		case *ast.UnaryExpr:
			if x.Op == token.AND {
				if id, ok := stripParens(x.X).(*ast.Ident); ok {
					if uv, okV := info.Uses[id].(*types.Var); okV && state.vars[uv] {
						unsafe = true
						return false
					}
				}
			}
		case *ast.CallExpr:
			if state.hasFieldAccess {
				if _, isFL := n.(*ast.FuncLit); !isFL && !isKnownPureLambdaOrAdapterCall(x, info) {
					unsafe = true
					return false
				}
			}
		case *ast.FuncLit:
			if state.hasFieldAccess {
				// Uncalled lambda literal (returned as value) does not execute during IIFE evaluation,
				// though we still check variable captures for addrTakenNames/vars if mutated inside.
				for nm := range state.addrTakenNames {
					if iifeASTContainsIdentNamed(x.Body, nm) {
						unsafe = true
						return false
					}
				}
				return false
			}
		}
		return true
	})
	return unsafe
}

func isKnownPureLambdaOrAdapterCall(call *ast.CallExpr, info *types.Info) bool {
	if call == nil || call.Ellipsis.IsValid() {
		return false
	}
	if funTV, ok := info.Types[call.Fun]; ok && funTV.IsType() {
		return true
	}
	if _, ok := isToIfaceCall(call, info); ok {
		return true
	}
	if _, ok := isUpcastCall(call, info); ok {
		return true
	}
	if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
		if strings.HasPrefix(fnName, "NewLambda_") || fnName == "AsIface" ||
			fnName == "ToAny" || fnName == "PtrToAny" || fnName == "SliceToAny" ||
			fnName == "BoxString" || fnName == "BoxNode" {
			return true
		}
	}
	if id, ok := stripIndexExprs(call.Fun).(*ast.Ident); ok {
		if fn, okF := info.Uses[id].(*types.Func); okF && fn != nil && strings.HasPrefix(fn.Name(), "NewLambda_") {
			return true
		}
	}
	return false
}
