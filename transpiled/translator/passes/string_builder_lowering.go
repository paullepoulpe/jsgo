// Clean-room string_builder_lowering pass: lowers non-escaping append-only local
// *j2gort.Ext_java_lang_StringBuilder variables to stack-allocated j2gort.LocalStringBuilder.

package passes

import (
	"cmp"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"slices"
)

// StringBuilderLoweringPass rewrites local *j2gort.Ext_java_lang_StringBuilder variables
// that never escape and never invoke UTF-16 index methods into stack-allocated j2gort.LocalStringBuilder.
type StringBuilderLoweringPass struct{}

func NewStringBuilderLoweringPass() *StringBuilderLoweringPass { return &StringBuilderLoweringPass{} }
func (p *StringBuilderLoweringPass) Name() string              { return "string_builder_lowering" }
func (p *StringBuilderLoweringPass) Description() string {
	return "Lower non-escaping append-only StringBuilder locals to stack-allocated j2gort.LocalStringBuilder"
}

func (p *StringBuilderLoweringPass) Apply(ctx *PassContext) (int, error) {
	if isJ2gortPkg(ctx.Pkg) || len(ctx.Files) == 0 {
		return 0, nil
	}
	var count int
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		parents, _ := buildParentMap(file)
		if n := lowerFileStringBuilders(file, ctx, parents); n > 0 {
			ctx.MarkFileModified(file)
			count += n
		}
	}
	return count, nil
}

type sbStep struct {
	method string
	args   []ast.Expr
	pos    token.Pos
}

type sbChainRewrite struct {
	outerNode ast.Node // *ast.ExprStmt or terminal *ast.CallExpr (.String() / .IsEmpty())
	steps     []sbStep
	terminal  string // "", "String", "IsEmpty", "Reset"
	termPos   token.Pos
}

type sbCandidate struct {
	v            *types.Var
	id           *ast.Ident
	declStmt     ast.Stmt
	initMethod   string // "", "InitCap", "InitStr"
	initArg      ast.Expr
	initSteps    []sbStep
	blankStmts   []ast.Stmt
	defCount     int
	hasAppendOrS bool
	disqualified bool
	uses         []sbChainRewrite
}

var sbOneArgAppendMethods = map[string]string{
	"Append_1_string":       "AppendStr",
	"Append_1_uint16":       "AppendChar",
	"Append_1_int32":        "AppendInt",
	"Append_1_int64":        "AppendLong",
	"Append_1_bool":         "AppendBool",
	"Append_1_float64":      "AppendDouble",
	"Append_1_float32":      "AppendFloat",
	"Append_1_slice_uint16": "AppendChars",
	"AppendCodePoint":       "AppendCodePoint",
}

var sbBoxedAppendMethods = map[string]string{
	"BoxString":      "AppendStr",
	"ValueOfInt":     "AppendInt",
	"BoxInt":         "AppendInt",
	"ValueOfLong":    "AppendLong",
	"BoxLong":        "AppendLong",
	"ValueOfBoolean": "AppendBool",
	"BoxBoolean":     "AppendBool",
	"ValueOfChar":    "AppendChar",
	"BoxChar":        "AppendChar",
	"ValueOfDouble":  "AppendDouble",
	"BoxDouble":      "AppendDouble",
	"ValueOfFloat":   "AppendFloat",
	"BoxFloat":       "AppendFloat",
}

func lowerFileStringBuilders(file *ast.File, ctx *PassContext, parents map[ast.Node]ast.Node) int {
	info := ctx.Info
	cands := make(map[*types.Var]*sbCandidate)

	tryRegisterDecl := func(id *ast.Ident, initExpr ast.Expr, declStmt ast.Stmt) {
		v := varFromObj(info.Defs[id])
		if id.Name == "_" || v == nil || v.IsField() || !isStringBuilderPtrType(v.Type()) {
			return
		}
		initMethod, initArg, initSteps, ok := parseSBInitChain(initExpr, info)
		if !ok {
			return
		}
		needsInitStmt := initMethod != "" || len(initSteps) > 0
		if needsInitStmt && (!canReplaceStmtIn(parents[declStmt]) || exprMentionsIdentName(initExpr, id.Name)) {
			return
		}
		cands[v] = &sbCandidate{
			v: v, id: id, declStmt: declStmt,
			initMethod: initMethod, initArg: initArg, initSteps: initSteps,
			hasAppendOrS: len(initSteps) > 0,
		}
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.DeclStmt:
			if gd, ok := s.Decl.(*ast.GenDecl); ok && gd.Tok == token.VAR && len(gd.Specs) == 1 && enclosingFuncBody(s, parents) != nil {
				if spec, ok := gd.Specs[0].(*ast.ValueSpec); ok && len(spec.Names) == 1 && len(spec.Values) == 1 {
					tryRegisterDecl(spec.Names[0], spec.Values[0], s)
				}
			}
		case *ast.AssignStmt:
			if s.Tok == token.DEFINE && len(s.Lhs) == 1 && len(s.Rhs) == 1 && enclosingFuncBody(s, parents) != nil {
				if id, ok := s.Lhs[0].(*ast.Ident); ok {
					tryRegisterDecl(id, s.Rhs[0], s)
				}
			}
		}
		return true
	})
	if len(cands) == 0 {
		return 0
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			if cand := cands[varFromObj(info.Defs[id])]; cand != nil {
				cand.defCount++
			}
			if cand := cands[varFromObj(info.Uses[id])]; cand != nil && !cand.disqualified {
				validateSBCandidateUse(id, cand, parents, info)
			}
		}
		return true
	})

	var eligible []*sbCandidate
	// order-independent: eligible is immediately sorted deterministically by token.Pos below
	for _, cand := range cands {
		if !cand.disqualified && cand.defCount == 1 && cand.hasAppendOrS &&
			canSpellJ2gortAtAllSites([]typeSpellSite{{file: file, pos: cand.declStmt.Pos()}}, ctx) {
			eligible = append(eligible, cand)
		}
	}
	if len(eligible) == 0 {
		return 0
	}
	slices.SortFunc(eligible, func(a, b *sbCandidate) int { return cmp.Compare(a.v.Pos(), b.v.Pos()) })

	type boundUse struct {
		cand *sbCandidate
		rw   sbChainRewrite
	}
	var exprUses, stmtUses []boundUse
	for _, cand := range eligible {
		for _, u := range cand.uses {
			if _, isStmt := u.outerNode.(*ast.ExprStmt); isStmt {
				stmtUses = append(stmtUses, boundUse{cand: cand, rw: u})
			} else {
				exprUses = append(exprUses, boundUse{cand: cand, rw: u})
			}
		}
	}
	spanCmp := func(a, b boundUse) int {
		na, nb := a.rw.outerNode, b.rw.outerNode
		if c := cmp.Compare(na.End()-na.Pos(), nb.End()-nb.Pos()); c != 0 {
			return c
		}
		return cmp.Compare(na.Pos(), nb.Pos())
	}
	slices.SortFunc(exprUses, spanCmp)
	slices.SortFunc(stmtUses, spanCmp)

	j2Qual, ok := j2gortQualifierForFile(file, ctx)
	if !ok {
		return 0
	}
	sites := 0
	for _, bu := range exprUses {
		oldExpr := bu.rw.outerNode.(ast.Expr)
		if parent := parents[oldExpr]; parent != nil && replaceExprInParent(parent, oldExpr, buildSBChainExpr(bu.cand.id.Name, oldExpr.Pos(), bu.rw), parents) {
			sites++
		}
	}
	for _, bu := range stmtUses {
		es := bu.rw.outerNode.(*ast.ExprStmt)
		es.X = buildSBChainExpr(bu.cand.id.Name, es.Pos(), bu.rw)
		parents[es.X] = es
		sites++
	}
	for _, cand := range eligible {
		for _, bs := range cand.blankStmts {
			replaceStmtInParent(parents[bs], bs, nil, parents)
		}
		if applySBDeclRewrite(cand, j2Qual, parents) {
			sites++
		}
	}
	return sites
}

func validateSBCandidateUse(id *ast.Ident, cand *sbCandidate, parents map[ast.Node]ast.Node, info *types.Info) {
	if as, ok := parents[id].(*ast.AssignStmt); ok && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && stripParens(as.Rhs[0]) == id {
		if lhsId, ok := stripParens(as.Lhs[0]).(*ast.Ident); ok && lhsId.Name == "_" {
			if canReplaceStmtIn(parents[as]) {
				cand.blankStmts = append(cand.blankStmts, as)
			}
			return
		}
	}
	var steps []sbStep
	for cur := ast.Expr(id); ; {
		top := unwrapSBCastsUp(cur, parents, info)
		parent := parents[top]
		if es, isES := parent.(*ast.ExprStmt); isES && es.X == top {
			if len(steps) == 0 {
				cand.disqualified = true
				return
			}
			cand.hasAppendOrS = true
			cand.uses = append(cand.uses, sbChainRewrite{outerNode: es, steps: steps})
			return
		}
		sel, selOK := parent.(*ast.SelectorExpr)
		if !selOK || sel.X != top {
			cand.disqualified = true
			return
		}
		mCall, callOK := parents[sel].(*ast.CallExpr)
		if !callOK || mCall.Fun != sel || mCall.Ellipsis.IsValid() {
			cand.disqualified = true
			return
		}
		if step, ok := classifySBAppendCall(sel.Sel.Name, mCall.Args, sel.Sel.Pos(), info); ok {
			steps = append(steps, step)
			cur = mCall
			continue
		}
		switch sel.Sel.Name {
		case "String", "IsEmpty":
			if _, isES := parents[mCall].(*ast.ExprStmt); isES || len(mCall.Args) != 0 {
				cand.disqualified = true
				return
			}
			if sel.Sel.Name == "String" {
				cand.hasAppendOrS = true
			}
			cand.uses = append(cand.uses, sbChainRewrite{outerNode: mCall, steps: steps, terminal: sel.Sel.Name, termPos: sel.Sel.Pos()})
			return
		case "SetLength":
			es, isES := parents[mCall].(*ast.ExprStmt)
			if !isES || es.X != mCall || len(steps) > 0 || len(mCall.Args) != 1 || !isZeroConstInt(mCall.Args[0], info) {
				cand.disqualified = true
				return
			}
			cand.uses = append(cand.uses, sbChainRewrite{outerNode: es, terminal: "Reset", termPos: sel.Sel.Pos()})
			return
		default:
			cand.disqualified = true
			return
		}
	}
}

func parseSBInitChain(expr ast.Expr, info *types.Info) (string, ast.Expr, []sbStep, bool) {
	var revSteps []sbStep
	for cur := stripParens(expr); ; {
		for {
			call, ok := cur.(*ast.CallExpr)
			if !ok || !isSBCastOrNullCheck(call, info) {
				break
			}
			cur = stripParens(call.Args[0])
		}
		call, isCall := cur.(*ast.CallExpr)
		if !isCall || call.Ellipsis.IsValid() {
			return "", nil, nil, false
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			slices.Reverse(revSteps)
			switch {
			case fnName == "NewExt_java_lang_StringBuilder" && len(call.Args) == 0:
				return "", nil, revSteps, true
			case fnName == "NewExt_java_lang_StringBuilder_1_int32" && len(call.Args) == 1 && isIntegerExpr(call.Args[0], info):
				return "InitCap", call.Args[0], revSteps, true
			case fnName == "NewExt_java_lang_StringBuilder_1_string" && len(call.Args) == 1 && hasExactOrDefaultStringType(call.Args[0], info):
				return "InitStr", call.Args[0], revSteps, true
			}
			return "", nil, nil, false
		}
		sel, selOK := stripParens(call.Fun).(*ast.SelectorExpr)
		if !selOK {
			return "", nil, nil, false
		}
		step, stepOK := classifySBAppendCall(sel.Sel.Name, call.Args, sel.Sel.Pos(), info)
		if !stepOK {
			return "", nil, nil, false
		}
		revSteps = append(revSteps, step)
		cur = stripParens(sel.X)
	}
}

func classifySBAppendCall(method string, args []ast.Expr, pos token.Pos, info *types.Info) (sbStep, bool) {
	if mapped, ok := sbOneArgAppendMethods[method]; ok && len(args) == 1 {
		return sbStep{method: mapped, args: []ast.Expr{args[0]}, pos: pos}, true
	}
	if method == "Append_3_any_int32_int32" && len(args) == 3 {
		return sbStep{method: "AppendSub", args: []ast.Expr{args[0], args[1], args[2]}, pos: pos}, true
	}
	if (method == "Append_1_any" || method == "Append_1_any_9389") && len(args) == 1 {
		if call, ok := stripParens(args[0]).(*ast.CallExpr); ok && len(call.Args) == 1 && !call.Ellipsis.IsValid() {
			if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
				inner := stripParens(call.Args[0])
				if mapped, ok := sbBoxedAppendMethods[fnName]; ok && (fnName != "BoxString" || hasExactOrDefaultStringType(inner, info)) {
					return sbStep{method: mapped, args: []ast.Expr{inner}, pos: pos}, true
				}
			}
		}
		return sbStep{method: "AppendAny", args: []ast.Expr{args[0]}, pos: pos}, true
	}
	return sbStep{}, false
}

func buildSBChainExpr(varName string, basePos token.Pos, rw sbChainRewrite) ast.Expr {
	var expr ast.Expr = &ast.Ident{NamePos: basePos, Name: varName}
	for _, st := range rw.steps {
		expr = &ast.CallExpr{Fun: &ast.SelectorExpr{X: expr, Sel: &ast.Ident{NamePos: st.pos, Name: st.method}}, Args: st.args}
	}
	if rw.terminal != "" {
		expr = &ast.CallExpr{Fun: &ast.SelectorExpr{X: expr, Sel: &ast.Ident{NamePos: rw.termPos, Name: rw.terminal}}}
	}
	return expr
}

func applySBDeclRewrite(cand *sbCandidate, j2Qual string, parents map[ast.Node]ast.Node) bool {
	pos := cand.declStmt.Pos()
	sbType := &ast.SelectorExpr{X: &ast.Ident{NamePos: pos, Name: j2Qual}, Sel: &ast.Ident{NamePos: pos, Name: "LocalStringBuilder"}}
	if cand.initMethod == "" && len(cand.initSteps) == 0 {
		switch s := cand.declStmt.(type) {
		case *ast.DeclStmt:
			spec := s.Decl.(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
			spec.Type, spec.Values = sbType, nil
			return true
		case *ast.AssignStmt:
			var lit ast.Expr = &ast.CompositeLit{Type: sbType}
			if inControlHeader(s, parents) {
				lit = &ast.ParenExpr{X: lit}
			}
			s.Rhs[0] = lit
			return true
		}
		return false
	}
	varDecl := &ast.DeclStmt{Decl: &ast.GenDecl{TokPos: pos, Tok: token.VAR, Specs: []ast.Spec{
		&ast.ValueSpec{Names: []*ast.Ident{{NamePos: cand.id.Pos(), Name: cand.id.Name}}, Type: sbType},
	}}}
	var chain ast.Expr = &ast.Ident{NamePos: pos, Name: cand.id.Name}
	if cand.initMethod != "" {
		chain = &ast.CallExpr{Fun: &ast.SelectorExpr{X: chain, Sel: &ast.Ident{NamePos: pos, Name: cand.initMethod}}, Args: []ast.Expr{cand.initArg}}
	}
	for _, st := range cand.initSteps {
		chain = &ast.CallExpr{Fun: &ast.SelectorExpr{X: chain, Sel: &ast.Ident{NamePos: st.pos, Name: st.method}}, Args: st.args}
	}
	return replaceStmtInParent(parents[cand.declStmt], cand.declStmt, []ast.Stmt{varDecl, &ast.ExprStmt{X: chain}}, parents)
}

func unwrapSBCastsUp(expr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) ast.Expr {
	for {
		if pe, ok := parents[expr].(*ast.ParenExpr); ok && pe.X == expr {
			expr = pe
			continue
		}
		if call, ok := parents[expr].(*ast.CallExpr); ok && call.Args != nil && len(call.Args) == 1 && call.Args[0] == expr && isSBCastOrNullCheck(call, info) {
			expr = call
			continue
		}
		return expr
	}
}

func isSBCastOrNullCheck(call *ast.CallExpr, info *types.Info) bool {
	if call.Ellipsis.IsValid() || len(call.Args) != 1 {
		return false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 {
		return false
	}
	if fnName == "CheckNotNull" || fnName == "CheckNotNullPtr" {
		return true
	}
	return (fnName == "CastPtr" || fnName == "Cast" || fnName == "UncheckedCast") && isStringBuilderPtrType(info.Types[call].Type)
}

func isStringBuilderPtrType(t types.Type) bool {
	ptr, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := types.Unalias(ptr.Elem()).(*types.Named)
	return ok && named.Obj() != nil && isJ2gortPkg(named.Obj().Pkg()) && named.Obj().Name() == "Ext_java_lang_StringBuilder"
}

func isZeroConstInt(expr ast.Expr, info *types.Info) bool {
	tv, ok := info.Types[expr]
	return ok && tv.Value != nil && tv.Value.Kind() == constant.Int && constant.Sign(tv.Value) == 0
}

func exprMentionsIdentName(expr ast.Expr, name string) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && id.Name == name {
			found = true
			return false
		}
		return !found
	})
	return found
}
