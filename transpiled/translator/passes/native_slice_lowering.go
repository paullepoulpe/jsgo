// Clean-room native_slice_lowering pass: lowers non-escaping ArrayList local variables
// and unexported struct fields (including nullable/lazy fields and linear builder handoffs)
// across a package to native Go slices []T.

package passes

import (
	"cmp"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"slices"
	"strconv"
	"strings"
)

// NativeSliceLoweringPass rewrites package-private ArrayList local variables and unexported
// struct fields to native Go slices []T whenever whole-package go/types analysis proves
// zero live aliasing/escape, guarded nullability, and a sound element type T.
type NativeSliceLoweringPass struct{}

func NewNativeSliceLoweringPass() *NativeSliceLoweringPass { return &NativeSliceLoweringPass{} }
func (p *NativeSliceLoweringPass) Name() string            { return "native_slice_lowering" }
func (p *NativeSliceLoweringPass) Description() string {
	return "Lower non-escaping ArrayList locals and unexported struct fields to native Go slices []T"
}

func (p *NativeSliceLoweringPass) Apply(ctx *PassContext) (int, error) {
	if isJ2gortPkg(ctx.Pkg) || len(ctx.Files) == 0 {
		return 0, nil
	}
	parentsByFile := make(map[*ast.File]map[ast.Node]ast.Node, len(ctx.Files))
	for _, f := range ctx.Files {
		parentsByFile[f], _ = buildParentMap(f)
	}

	count := 0
	eligible := analyzeSliceCandidates(ctx, parentsByFile)
	if len(eligible) > 0 {
		count += applySliceRewrites(ctx, eligible, parentsByFile)
	}
	count += applyOrderedMapValueSliceLowering(ctx, parentsByFile)
	return count, nil
}

type sliceInitKind uint8

const (
	sliceInitEmpty sliceInitKind = iota
	sliceInitCap
	sliceInitCopySlice
	sliceInitCopyTypedSlice
	sliceInitCopyCol
	sliceInitCopyIter
	sliceInitMoveCand
)

type opKind uint8

const (
	opDeclField opKind = iota
	opDeclVarSpec
	opDeclShortVar
	opAssignInit
	opAssignNil
	opAssignReverseSelf
	opAccFuncDecl
	opAccCallRename
	opExprAccParamArg
	opExprGet
	opExprSize
	opExprIsEmpty
	opExprRawUnbox
	opExprImmListCopy
	opExprImmCopyToSlice
	opExprUnmodList
	opExprUnmodListLive
	opExprReadOnlyIter
	opExprSliceIter
	opExprToJavaList
	opExprWrapTypedArrayList
	opExprFormatJavaStr
	opExprGetFirst
	opExprGetLast
	opExprRemoveAt
	opStmtAdd
	opStmtInsertAt
	opStmtSet
	opStmtSortNat
	opStmtSortCmp
	opStmtRemoveAt
	opStmtRemoveFirst
	opStmtRemoveLast
	opStmtAddAllSlice
	opStmtAddAllCol
	opStmtClear
	opStmtIterSlice
	opStmtIterCursor
	opStmtForEach
)

type sliceRewriteOp struct {
	kind          opKind
	initKind      sliceInitKind
	cand          *sliceCandidate
	otherCand     *sliceCandidate
	forEachSpec   *setForEachSpec
	file          *ast.File
	spanNode      ast.Node
	accessExpr    ast.Expr
	otherExpr     ast.Expr
	call          *ast.CallExpr
	outerExpr     ast.Expr
	stmt          ast.Stmt
	forStmt       *ast.ForStmt
	block         *ast.BlockStmt
	label         *ast.Ident
	elemIdent     *ast.Ident
	dropStmts     int
	unwrapsBox    bool
	isReverseIter bool
}

type typeSpellSite struct {
	file       *ast.File
	pos        token.Pos
	accessExpr ast.Expr
	owner      *types.Named
	recvTArgs  *types.TypeList
}

type derefGuardCheck struct {
	node       ast.Node
	accessExpr ast.Expr
	file       *ast.File
}

type unmodLocalCheck struct {
	callNode ast.Node
	usePos   token.Pos
	file     *ast.File
}

type handoffEdge struct {
	dst          *sliceCandidate
	src          *sliceCandidate
	propagateNil bool
}

type sliceCandidate struct {
	v                     *types.Var
	accFn                 *types.Func
	accFnDecl             *ast.FuncDecl
	accParamField         *ast.Field
	accParamIdx           int
	isField               bool
	isAccParam            bool
	isExportedAccFn       bool
	ownerStruct           *types.Named
	declFile              *ast.File
	defCount              int
	inPkgAccCalls         int
	hasInitFields         bool
	initInAllCtors        bool
	initToNilAtDecl       bool
	hasNilAssign          bool
	mayBeNil              bool
	hasAdd                bool
	hasCopyInit           bool
	hasInsertAt           bool
	hasSet                bool
	hasRemove             bool
	hasClear              bool
	hasIterCall           bool
	hasLiveUnmod          bool
	hasWrapTypedArrayList bool
	isMoveInitDst         bool
	disqualified          bool
	elemType              types.Type
	writeTypes            []types.Type
	readTypes             []types.Type
	siteConstraints       []siteTypeConstraint
	spellSites            []typeSpellSite
	j2gortSites           []typeSpellSite
	clearPositions        []token.Pos
	mutNodes              []ast.Node
	unmodChecks           []unmodLocalCheck
	derefChecks           []derefGuardCheck
	loopBodies            []*ast.BlockStmt
	ops                   []sliceRewriteOp
	parent                *sliceCandidate
}

func (c *sliceCandidate) root() *sliceCandidate {
	r := c
	for r.parent != nil {
		r = r.parent
	}
	for cur := c; cur != r; {
		next := cur.parent
		cur.parent = r
		cur = next
	}
	return r
}

func unionSliceCandidates(a, b *sliceCandidate) {
	ra, rb := a.root(), b.root()
	if ra != rb {
		rb.parent = ra
	}
}

func (c *sliceCandidate) addWriteType(t types.Type, accessExpr ast.Expr, info *types.Info) bool {
	if t == nil || isInvalidOrUntyped(t) {
		c.disqualified = true
		return false
	}
	if c.isField && c.ownerStruct != nil {
		lifted, sc, ok := liftSiteTypeForOwner(t, accessExpr, c.ownerStruct, info)
		if !ok {
			c.disqualified = true
			return false
		}
		if sc != nil {
			c.siteConstraints = append(c.siteConstraints, *sc)
			return true
		}
		t = lifted
	}
	c.writeTypes = append(c.writeTypes, t)
	return true
}

func (c *sliceCandidate) addReadType(t types.Type, accessExpr ast.Expr, info *types.Info) bool {
	if t == nil || isInvalidOrUntyped(t) {
		c.disqualified = true
		return false
	}
	if c.isField && c.ownerStruct != nil {
		lifted, sc, ok := liftSiteTypeForOwner(t, accessExpr, c.ownerStruct, info)
		if !ok {
			c.disqualified = true
			return false
		}
		if sc != nil {
			c.siteConstraints = append(c.siteConstraints, *sc)
			return true
		}
		t = lifted
	}
	c.readTypes = append(c.readTypes, t)
	return true
}

func verifyMemberSiteConstraints(declT types.Type, m *sliceCandidate) bool {
	if typeContainsTypeParam(declT) {
		if m.isField && !typeOnlyUsesOwnerTypeParams(declT, m.ownerStruct) {
			return false
		}
	}
	for _, sc := range m.siteConstraints {
		inst, ok := instantiateOwnerType(declT, m.ownerStruct, sc.targs)
		if !ok {
			return false
		}
		if !types.Identical(inst, sc.siteType) && !(isNonEmptyInterfaceType(inst) && (isNonEmptyInterfaceType(sc.siteType) || isPtrToStructType(sc.siteType)) && types.AssignableTo(sc.siteType, inst)) {
			return false
		}
	}
	return true
}

func resolveGroupElemType(members []*sliceCandidate) (types.Type, bool) {
	var all, reads, writes []types.Type
	for _, m := range members {
		reads = append(reads, m.readTypes...)
		writes = append(writes, m.writeTypes...)
		all = append(all, m.writeTypes...)
		all = append(all, m.readTypes...)
	}
	if len(all) == 0 {
		for _, m := range members {
			for _, sc := range m.siteConstraints {
				all = append(all, sc.siteType)
			}
		}
	}
	if len(all) == 0 {
		return nil, false
	}
	t0 := all[0]
	allIdentical := true
	for _, t := range all[1:] {
		if !types.Identical(t0, t) {
			allIdentical = false
			break
		}
	}
	if allIdentical {
		for _, m := range members {
			if !verifyMemberSiteConstraints(t0, m) {
				return nil, false
			}
		}
		return t0, true
	}
	// Allow interface element type I when all reads expect I and all writes are interface subtypes of I
	// (e.g., FunctionType written into []JSType in UnionType_Builder, which has zero typed-nil hazard).
	var targetIface types.Type
	for _, rt := range reads {
		if isNonEmptyInterfaceType(rt) {
			if targetIface == nil {
				targetIface = rt
			} else if !types.Identical(targetIface, rt) {
				return nil, false
			}
		} else {
			return nil, false
		}
	}
	if targetIface == nil {
		return nil, false
	}
	for _, wt := range writes {
		if !types.Identical(wt, targetIface) && !((isNonEmptyInterfaceType(wt) || isPtrToStructType(wt)) && types.AssignableTo(wt, targetIface)) {
			return nil, false
		}
	}
	for _, m := range members {
		if !verifyMemberSiteConstraints(targetIface, m) {
			return nil, false
		}
	}
	return targetIface, true
}

func analyzeSliceCandidates(ctx *PassContext, parentsByFile map[*ast.File]map[ast.Node]ast.Node) []*sliceCandidate {
	info := ctx.Info
	hasCopyHelpers := j2gortHasExport(ctx, "CloneNonNilSlice")
	candidates := make(map[*types.Var]*sliceCandidate)
	accCandsByFn := make(map[*types.Func][]*sliceCandidate)
	structFields := make(map[*types.Named][]*sliceCandidate)
	var handoffs []handoffEdge

	type pendingLocalSlice struct {
		cand     *sliceCandidate
		initExpr ast.Expr
		file     *ast.File
		span     ast.Node
		kind     opKind
	}
	var pendingLocals []pendingLocalSlice

	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}
				for _, s := range d.Specs {
					ts, ok := s.(*ast.TypeSpec)
					if !ok || ts.Assign.IsValid() {
						continue
					}
					st, stOK := ts.Type.(*ast.StructType)
					tn, tnOK := info.Defs[ts.Name].(*types.TypeName)
					if !stOK || !tnOK || st.Fields == nil {
						continue
					}
					named := namedOrigin(tn.Type())
					if named == nil {
						continue
					}
					for _, fld := range st.Fields.List {
						if len(fld.Names) != 1 {
							continue
						}
						if v := varFromObj(info.Defs[fld.Names[0]]); v != nil && v.IsField() && !v.Exported() && !v.Embedded() && isCandidateListType(v.Type()) {
							cand := &sliceCandidate{
								v: v, isField: true, ownerStruct: named, declFile: file,
								spellSites: []typeSpellSite{{file: file, pos: fld.Pos(), owner: named}},
							}
							cand.ops = append(cand.ops, sliceRewriteOp{kind: opDeclField, cand: cand, file: file, spanNode: fld})
							candidates[v] = cand
							structFields[named] = append(structFields[named], cand)
						}
					}
				}

			case *ast.FuncDecl:
				// Collect single slice-accumulator parameter on non-generic functions and priv_ methods.
				if d.Body == nil || d.Type.Params == nil || (d.Type.TypeParams != nil && len(d.Type.TypeParams.List) > 0) {
					continue
				}
				fnObj, ok := info.Defs[d.Name].(*types.Func)
				if !ok || fnObj == nil || !isNonGenericFunc(fnObj) {
					continue
				}
				sig, ok := fnObj.Type().(*types.Signature)
				if !ok || sig.Variadic() || (d.Type.Results != nil && len(d.Type.Results.List) > 1) {
					continue
				}
				if sig.Recv() != nil {
					if fnObj.Exported() || !strings.HasPrefix(fnObj.Name(), "priv_") {
						continue
					}
				} else if fnObj.Exported() {
					if d.Name.Name == "main" || d.Name.Name == "init" || !allFuncParamsNamed(d.Type) {
						continue
					}
					if ctx.Pkg != nil && ctx.Pkg.Scope() != nil && ctx.Pkg.Scope().Lookup("sliceAcc_"+d.Name.Name) != nil {
						continue
					}
				}
				type matchedAccParam struct {
					fld *ast.Field
					pv  *types.Var
					idx int
				}
				var matched []matchedAccParam
				argIdx := 0
				for _, fld := range d.Type.Params.List {
					if len(fld.Names) == 1 && fld.Names[0].Name != "_" {
						if pv := varFromObj(info.Defs[fld.Names[0]]); pv != nil && isCandidateAccParamType(pv.Type()) {
							matched = append(matched, matchedAccParam{fld: fld, pv: pv, idx: argIdx})
						}
					}
					if len(fld.Names) == 0 {
						argIdx++
					} else {
						argIdx += len(fld.Names)
					}
				}
				if len(matched) > 0 && (!fnObj.Exported() || len(matched) == 1) {
					fnOrig := funcOrigin(fnObj)
					for _, m := range matched {
						span := ast.Node(m.fld)
						if fnObj.Exported() {
							span = d
						}
						ac := &sliceCandidate{
							v: m.pv, accFn: fnOrig, accFnDecl: d, accParamField: m.fld,
							accParamIdx: m.idx, isAccParam: true, isExportedAccFn: fnObj.Exported(),
							declFile:    file,
							spellSites:  []typeSpellSite{{file: file, pos: m.fld.Pos()}},
							j2gortSites: []typeSpellSite{{file: file, pos: m.fld.Pos()}},
						}
						ac.ops = append(ac.ops, sliceRewriteOp{kind: opAccFuncDecl, cand: ac, file: file, spanNode: span})
						candidates[m.pv] = ac
						accCandsByFn[fnOrig] = append(accCandsByFn[fnOrig], ac)
					}
				}
			}
		}

		collectLocal := func(id *ast.Ident, initExpr ast.Expr, span ast.Node, kind opKind) {
			v := varFromObj(info.Defs[id])
			if id.Name == "_" || v == nil || v.IsField() || candidates[v] != nil || !isCandidateListType(v.Type()) {
				return
			}
			if initExpr == nil || isNilIdent(initExpr, info) {
				if kind != opDeclVarSpec {
					return
				}
				cand := &sliceCandidate{
					v: v, declFile: file, initToNilAtDecl: true,
					spellSites: []typeSpellSite{{file: file, pos: span.Pos()}},
				}
				cand.ops = append(cand.ops, sliceRewriteOp{kind: kind, cand: cand, file: file, spanNode: span})
				candidates[v] = cand
				return
			}
			cand := &sliceCandidate{
				v: v, declFile: file,
				spellSites: []typeSpellSite{{file: file, pos: span.Pos()}},
			}
			candidates[v] = cand
			pendingLocals = append(pendingLocals, pendingLocalSlice{
				cand: cand, initExpr: initExpr, file: file, span: span, kind: kind,
			})
		}

		ast.Inspect(file, func(n ast.Node) bool {
			switch s := n.(type) {
			case *ast.DeclStmt:
				if gd, ok := s.Decl.(*ast.GenDecl); ok && gd.Tok == token.VAR && len(gd.Specs) == 1 && enclosingFuncBody(s, parents) != nil {
					if spec, ok := gd.Specs[0].(*ast.ValueSpec); ok && len(spec.Names) == 1 && len(spec.Values) <= 1 {
						var initExpr ast.Expr
						if len(spec.Values) == 1 {
							initExpr = spec.Values[0]
						}
						collectLocal(spec.Names[0], initExpr, spec, opDeclVarSpec)
					}
				}
			case *ast.AssignStmt:
				if s.Tok == token.DEFINE && len(s.Lhs) == 1 && len(s.Rhs) == 1 && enclosingFuncBody(s, parents) != nil {
					if id, ok := s.Lhs[0].(*ast.Ident); ok {
						collectLocal(id, s.Rhs[0], s, opDeclShortVar)
					}
				}
			}
			return true
		})
	}

	consumedCopyInitSrcs := make(map[ast.Expr]bool)
	for _, pl := range pendingLocals {
		initK, argExpr, otherCand, typedSliceElem, ok := classifySliceInitExpr(pl.initExpr, pl.cand, candidates, hasCopyHelpers, ctx, info)
		if !ok {
			if moveExpr, moveSrc, moveOK := matchFieldMoveInitToLocal(pl.initExpr, pl.span, pl.cand, candidates, parentsByFile[pl.file], info); moveOK {
				pl.cand.isMoveInitDst = true
				consumedCopyInitSrcs[moveExpr] = true
				handoffs = append(handoffs, handoffEdge{dst: pl.cand, src: moveSrc, propagateNil: true})
				pl.cand.ops = append(pl.cand.ops, sliceRewriteOp{
					kind: pl.kind, initKind: sliceInitMoveCand, cand: pl.cand, otherCand: moveSrc,
					file: pl.file, spanNode: pl.span, otherExpr: moveExpr,
				})
				continue
			}
			pl.cand.disqualified = true
			delete(candidates, pl.cand.v)
			continue
		}
		if initK != sliceInitEmpty && !(initK == sliceInitCap && isNonNegativeConstInt(argExpr, info)) {
			pl.cand.j2gortSites = append(pl.cand.j2gortSites, typeSpellSite{file: pl.file, pos: pl.initExpr.Pos()})
		}
		if initK == sliceInitCopySlice && otherCand != nil {
			pl.cand.hasCopyInit = true
			consumedCopyInitSrcs[stripParens(argExpr)] = true
			otherCand.derefChecks = append(otherCand.derefChecks, derefGuardCheck{node: pl.span, accessExpr: argExpr, file: pl.file})
			handoffs = append(handoffs, handoffEdge{dst: pl.cand, src: otherCand, propagateNil: false})
		} else if initK == sliceInitCopyTypedSlice && typedSliceElem != nil {
			pl.cand.hasCopyInit = true
			if !pl.cand.addWriteType(typedSliceElem, nil, info) {
				pl.cand.disqualified = true
				delete(candidates, pl.cand.v)
				continue
			}
		} else if initK == sliceInitCopyCol || initK == sliceInitCopyIter {
			pl.cand.hasCopyInit = true
		}
		pl.cand.ops = append(pl.cand.ops, sliceRewriteOp{
			kind: pl.kind, initKind: initK, cand: pl.cand, otherCand: otherCand,
			file: pl.file, spanNode: pl.span, otherExpr: argExpr,
		})
	}

	if len(candidates) == 0 {
		return nil
	}
	if len(structFields) > 0 {
		verifyStructInvariants(ctx, parentsByFile, candidates, structFields, hasCopyHelpers)
	}

	varUseCounts := make(map[*types.Var]int, len(info.Uses))
	varMaxUsePos := make(map[*types.Var]token.Pos, len(info.Uses))
	// order-independent: commutative per-variable integer count and max-Pos accumulation
	for id, obj := range info.Uses {
		if v := varFromObj(obj); v != nil {
			varUseCounts[v]++
			if id.Pos() > varMaxUsePos[v] {
				varMaxUsePos[v] = id.Pos()
			}
		}
	}

	funcsMutating := make(map[*types.Var]map[*types.Func]bool)
	markMut := func(cand *sliceCandidate, node ast.Node, parents map[ast.Node]ast.Node) {
		cand.mutNodes = append(cand.mutNodes, node)
		if fn := enclosingFuncObj(node, parents, info); fn != nil {
			if funcsMutating[cand.v] == nil {
				funcsMutating[cand.v] = make(map[*types.Func]bool)
			}
			funcsMutating[cand.v][fn] = true
		}
	}

	var builderChecks []builderHandoffSpec

	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		ast.Inspect(file, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if cand := candidates[varFromObj(info.Defs[id])]; cand != nil {
				cand.defCount++
				if !ctx.ShouldTransform(file) {
					cand.disqualified = true
				}
			}
			if cand := candidates[varFromObj(info.Uses[id])]; cand != nil && !cand.disqualified {
				if !ctx.ShouldTransform(file) {
					cand.disqualified = true
				} else {
					validateSliceCandidateUse(id, file, cand, candidates, accCandsByFn, consumedCopyInitSrcs, hasCopyHelpers, ctx, parents, varUseCounts, varMaxUsePos, markMut, &handoffs, &builderChecks)
				}
			}
			if fnObj, ok := info.Uses[id].(*types.Func); ok {
				for _, ac := range accCandsByFn[funcOrigin(fnObj)] {
					if !ac.disqualified {
						validateAccFuncCallSite(id, file, ac, candidates, ctx.ShouldTransform(file), parents, info)
					}
				}
			}
			return true
		})
	}

	for _, bs := range builderChecks {
		if !verifyBuilderSingleUseInPkg(ctx, parentsByFile, bs) {
			bs.src.disqualified = true
			bs.dst.disqualified = true
		}
	}

	for _, h := range handoffs {
		unionSliceCandidates(h.dst, h.src)
	}

	// Propagate mayBeNil across candidates and ownership-handoff edges to fixed point.
	// order-independent: sets per-candidate boolean mayBeNil flag
	for _, cand := range candidates {
		if cand.isAccParam {
			continue
		}
		if cand.isField {
			if cand.hasNilAssign || (!cand.hasInitFields && !cand.initInAllCtors) {
				cand.mayBeNil = true
			}
		} else if cand.initToNilAtDecl || cand.hasNilAssign {
			cand.mayBeNil = true
		}
	}
	for changed := true; changed; {
		changed = false
		for _, h := range handoffs {
			if h.propagateNil && h.src.mayBeNil && !h.dst.mayBeNil {
				h.dst.mayBeNil = true
				changed = true
			}
		}
	}

	// Require slice-accumulator parameters (`*[]T`) to perform at least one direct or transitive mutation.
	accMutates := make(map[*sliceCandidate]bool)
	for _, cand := range candidates {
		if cand.isAccParam && (cand.hasAdd || cand.hasInsertAt || cand.hasSet || cand.hasRemove || cand.hasClear) {
			accMutates[cand] = true
		}
	}
	for changed := true; changed; {
		changed = false
		for _, cand := range candidates {
			if !cand.isAccParam || accMutates[cand] {
				continue
			}
			for _, op := range cand.ops {
				if op.kind == opExprAccParamArg && op.otherCand != nil && accMutates[op.otherCand] {
					accMutates[cand] = true
					changed = true
					break
				}
			}
		}
	}
	for _, cand := range candidates {
		for _, op := range cand.ops {
			if op.kind == opExprAccParamArg && op.otherCand != nil && accMutates[op.otherCand] && op.call != nil {
				markMut(cand, op.call, parentsByFile[op.file])
			}
		}
	}

	// Validate each candidate individually before group resolution.
	// order-independent: sets per-candidate boolean disqualified flag
	for v, cand := range candidates {
		if cand.disqualified || cand.defCount != 1 || !ctx.ShouldTransform(cand.declFile) {
			cand.disqualified = true
			continue
		}
		if cand.isAccParam && (cand.inPkgAccCalls == 0 || (cand.isExportedAccFn && !accMutates[cand])) {
			cand.disqualified = true
			continue
		}
		if cand.hasLiveUnmod && cand.hasNilAssign {
			cand.disqualified = true
			continue
		}
		if cand.isMoveInitDst && (cand.hasAdd || cand.hasInsertAt || cand.hasSet || cand.hasRemove || cand.hasClear || cand.hasLiveUnmod || cand.hasNilAssign) {
			cand.disqualified = true
			continue
		}
		if cand.hasIterCall && (cand.hasClear || cand.hasRemove || cand.hasInsertAt || cand.hasSet) {
			cand.disqualified = true
			continue
		}
		if anyLoopMutatesCandidate(cand, funcsMutating[v], info) {
			cand.disqualified = true
			continue
		}
		if len(cand.unmodChecks) > 0 && !verifyLocalUnmodNoPostMutation(cand, parentsByFile) {
			hasReassignInit := slices.ContainsFunc(cand.ops, func(op sliceRewriteOp) bool { return op.kind == opAssignInit })
			hasValueHandoff := slices.ContainsFunc(handoffs, func(h handoffEdge) bool { return h.propagateNil && (h.src == cand || h.dst == cand) })
			if !cand.isField && !cand.isAccParam && !cand.initToNilAtDecl && !cand.hasNilAssign && !cand.isMoveInitDst && !cand.hasWrapTypedArrayList && !hasReassignInit && !hasValueHandoff && j2gortHasExport(ctx, "SliceUnmodifiableListLive") {
				cand.hasLiveUnmod = true
				for i := range cand.ops {
					if cand.ops[i].kind == opExprUnmodList {
						cand.ops[i].kind = opExprUnmodListLive
					}
				}
			} else {
				cand.disqualified = true
				continue
			}
		}
		if cand.mayBeNil {
			for _, dc := range cand.derefChecks {
				if !isGuardedNonNilSliceUse(dc.node, dc.accessExpr, cand, parentsByFile[dc.file], hasCopyHelpers, ctx, info) {
					cand.disqualified = true
					break
				}
			}
		}
	}

	groups := make(map[*sliceCandidate][]*sliceCandidate)
	// order-independent: groups are resolved independently and out is sorted by token.Pos below
	for _, cand := range candidates {
		r := cand.root()
		groups[r] = append(groups[r], cand)
	}

	var out []*sliceCandidate
	// order-independent: out is immediately sorted deterministically by token.Pos below
	for _, members := range groups {
		groupOK := true
		groupHasPopulate := false
		groupHasWrapTyped := false
		groupHasIterCall := false
		groupHasMoveInitDst := false
		groupHasStructMut := false
		groupHasAnyMut := false
		for _, m := range members {
			if m.disqualified {
				groupOK = false
				break
			}
			if m.hasAdd || m.hasCopyInit {
				groupHasPopulate = true
			}
			if m.hasIterCall {
				groupHasIterCall = true
			}
			if m.isMoveInitDst {
				groupHasMoveInitDst = true
			}
			if m.hasClear || m.hasRemove || m.hasInsertAt || m.hasSet {
				groupHasStructMut = true
			}
			if m.isAccParam && (m.hasAdd || m.hasInsertAt || m.hasSet || m.hasRemove || m.hasClear) {
				groupHasAnyMut = true
			}
			if m.hasWrapTypedArrayList {
				groupHasWrapTyped = true
				if m.hasLiveUnmod || len(m.unmodChecks) > 0 || m.hasIterCall {
					groupOK = false
					break
				}
			}
		}
		if (groupHasIterCall && groupHasStructMut) || (groupHasMoveInitDst && groupHasAnyMut) {
			continue
		}
		if groupHasWrapTyped && len(members) != 1 {
			continue
		}
		if !groupOK || !groupHasPopulate {
			continue
		}
		elemT, ok := resolveGroupElemType(members)
		if !ok || isEmptyInterfaceType(elemT) {
			continue
		}
		for _, m := range members {
			if !canSpellTypeAtAllSites(elemT, m.spellSites, ctx) || !canSpellJ2gortAtAllSites(m.j2gortSites, ctx) {
				groupOK = false
				break
			}
			if m.hasClear && typeNeedsGCZeroing(elemT) {
				if slices.ContainsFunc(m.clearPositions, func(pos token.Pos) bool { return !isUnshadowedBuiltin("clear", pos, ctx.Pkg) }) {
					groupOK = false
					break
				}
			}
		}
		if !groupOK {
			continue
		}
		for _, m := range members {
			m.elemType = elemT
			out = append(out, m)
		}
	}
	slices.SortFunc(out, func(a, b *sliceCandidate) int { return cmp.Compare(a.v.Pos(), b.v.Pos()) })
	return out
}

func verifyStructInvariants(
	ctx *PassContext,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
	candidates map[*types.Var]*sliceCandidate,
	structFields map[*types.Named][]*sliceCandidate,
	hasCopyHelpers bool,
) {
	info := ctx.Info
	disqualify := func(named *types.Named) {
		for _, c := range structFields[named] {
			c.disqualified = true
		}
	}
	isNonNilSliceFieldRhs := func(rhs ast.Expr, cand *sliceCandidate) bool {
		_, _, _, _, ok := classifySliceInitExpr(rhs, cand, candidates, hasCopyHelpers, ctx, info)
		return ok
	}

	// Collect fields initialized via NewExt_java_util_ArrayList* inside helper methods on *S.
	methodInits := make(map[*types.Func]map[*types.Var]bool)
	ctorFuncs := make(map[*types.Named][]*ast.FuncDecl)
	newFuncs := make(map[*types.Named][]*ast.FuncDecl)

	for _, file := range ctx.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			if fn.Recv != nil && len(fn.Recv.List) == 1 && len(fn.Recv.List[0].Names) == 1 {
				recvVar := varFromObj(info.Defs[fn.Recv.List[0].Names[0]])
				if recvVar == nil {
					continue
				}
				named := ptrNamedOrigin(recvVar.Type())
				if named == nil || len(structFields[named]) == 0 {
					continue
				}
				fnObj, _ := info.Defs[fn.Name].(*types.Func)
				fnOrig := funcOrigin(fnObj)
				isInitFields := fn.Name.Name == "initFields_"+named.Obj().Name() || fn.Name.Name == "initFields" ||
					strings.HasPrefix(fn.Name.Name, "initFields_")
				for _, stmt := range fn.Body.List {
					as, ok := stmt.(*ast.AssignStmt)
					if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
						continue
					}
					if sel, ok := stripParens(as.Lhs[0]).(*ast.SelectorExpr); ok {
						if recvId, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[recvId]) == recvVar {
							if cand := candidates[varFromObj(info.Uses[sel.Sel])]; cand != nil && cand.isField {
								if isNonNilSliceFieldRhs(as.Rhs[0], cand) {
									if isInitFields {
										cand.hasInitFields = true
									}
									if fnOrig != nil {
										if methodInits[fnOrig] == nil {
											methodInits[fnOrig] = make(map[*types.Var]bool)
										}
										methodInits[fnOrig][cand.v] = true
									}
								}
							}
						}
					}
				}
			} else if fn.Recv == nil && strings.HasPrefix(fn.Name.Name, "Init") && len(fn.Type.Params.List) >= 1 {
				p0 := fn.Type.Params.List[0]
				if len(p0.Names) == 1 {
					if v0 := varFromObj(info.Defs[p0.Names[0]]); v0 != nil {
						if named := ptrNamedOrigin(v0.Type()); named != nil && len(structFields[named]) > 0 {
							if strings.HasPrefix(fn.Name.Name, "Init"+named.Obj().Name()) ||
								strings.HasPrefix(fn.Name.Name, "Init"+strings.TrimSuffix(named.Obj().Name(), "Base")) {
								ctorFuncs[named] = append(ctorFuncs[named], fn)
							}
						}
					}
				}
			} else if fn.Recv == nil && strings.HasPrefix(fn.Name.Name, "New") && fn.Type.Results != nil && len(fn.Type.Results.List) == 1 {
				if fnObj, ok := info.Defs[fn.Name].(*types.Func); ok && fnObj != nil {
					if sig, ok := fnObj.Type().(*types.Signature); ok && sig.Results().Len() == 1 {
						if named := ptrNamedOrigin(sig.Results().At(0).Type()); named != nil && len(structFields[named]) > 0 &&
							strings.HasPrefix(fn.Name.Name, "New"+named.Obj().Name()) {
							newFuncs[named] = append(newFuncs[named], fn)
						}
					}
				}
			}
		}
	}

	// Check whether all Init<Struct>* functions initialize each field (directly, via handoff, via helper method, or via constructor chaining).
	// order-independent: sets per-candidate boolean initInAllCtors flag
	for named, cands := range structFields {
		ctors := ctorFuncs[named]
		useNew := false
		if len(ctors) == 0 {
			ctors = newFuncs[named]
			useNew = true
		}
		if len(ctors) == 0 {
			continue
		}
		for _, cand := range cands {
			if cand.hasInitFields {
				cand.initInAllCtors = true
				continue
			}
			allInit := true
			for _, fn := range ctors {
				var recvVar *types.Var
				if !useNew {
					recvVar = varFromObj(info.Defs[fn.Type.Params.List[0].Names[0]])
				}
				found := false
				for _, stmt := range fn.Body.List {
					switch s := stmt.(type) {
					case *ast.AssignStmt:
						if s.Tok == token.DEFINE && useNew && len(s.Lhs) == 1 {
							if id, ok := s.Lhs[0].(*ast.Ident); ok {
								recvVar = varFromObj(info.Defs[id])
							}
						}
						if s.Tok == token.ASSIGN && len(s.Lhs) == 1 && len(s.Rhs) == 1 && recvVar != nil {
							if sel, ok := stripParens(s.Lhs[0]).(*ast.SelectorExpr); ok {
								if id, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == recvVar && varFromObj(info.Uses[sel.Sel]) == cand.v {
									if isNonNilSliceFieldRhs(s.Rhs[0], cand) {
										found = true
									} else if rhsCand := resolveCandidateAccess(unwrapUnmodifiableListCall(s.Rhs[0], info), candidates, info); rhsCand != nil {
										found = true
									}
								}
							}
						}
					case *ast.ExprStmt:
						if call, ok := stripParens(s.X).(*ast.CallExpr); ok && recvVar != nil {
							if sel, ok := stripParens(call.Fun).(*ast.SelectorExpr); ok {
								if id, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == recvVar {
									if callee, _ := calleeFunc(call, info); callee != nil && methodInits[funcOrigin(callee)][cand.v] {
										found = true
									}
								}
							} else if callee, _ := calleeFunc(call, info); callee != nil && len(call.Args) >= 1 {
								if id0, ok := stripParens(call.Args[0]).(*ast.Ident); ok && varFromObj(info.Uses[id0]) == recvVar {
									baseObjName := strings.TrimSuffix(named.Obj().Name(), "Base")
									if strings.HasPrefix(callee.Name(), "Init"+named.Obj().Name()) ||
										strings.HasPrefix(callee.Name(), "Init"+baseObjName) {
										found = true
									}
								}
							}
						}
					}
				}
				if !found {
					allInit = false
					break
				}
			}
			cand.initInAllCtors = allInit
		}
	}

	watched := make(map[*types.Named]bool, len(structFields))
	// order-independent: populates membership set of candidate struct types
	for named := range structFields {
		watched[named] = true
	}
	checkStructValueCopies(ctx, parentsByFile, watched, func(fv *types.Var) {
		if c := candidates[fv]; c != nil {
			c.disqualified = true
		}
	}, disqualify)
}

func checkStructValueCopies(
	ctx *PassContext,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
	watched map[*types.Named]bool,
	disqualifyField func(*types.Var),
	disqualifyStruct func(*types.Named),
) {
	info := ctx.Info
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CompositeLit:
				if tv := info.Types[x]; tv.Type != nil && len(x.Elts) > 0 {
					if named := derefNamedOrigin(tv.Type); named != nil && watched[named] {
						for _, elt := range x.Elts {
							kv, isKV := elt.(*ast.KeyValueExpr)
							if !isKV {
								disqualifyStruct(named)
								break
							}
							if kId, ok := stripParens(kv.Key).(*ast.Ident); ok {
								if fv := varFromObj(info.Uses[kId]); fv != nil {
									disqualifyField(fv)
								}
							}
						}
					}
				}
			case *ast.MapType:
				if kt := info.Types[x.Key].Type; kt != nil {
					if named := namedOrigin(kt); named != nil && watched[named] {
						disqualifyStruct(named)
					}
				}
			case ast.Expr:
				tv := info.Types[x]
				if tv.IsType() || tv.Type == nil {
					return true
				}
				named := namedOrigin(tv.Type)
				if named == nil || !watched[named] {
					return true
				}
				if fn := enclosingFuncObj(x, parents, info); fn != nil &&
					(fn.Name() == "releaseNodeTraversal_" || fn.Name() == "ReleaseSlabs_" || fn.Name() == "release_" || strings.HasPrefix(fn.Name(), "allocScratch_")) {
					return true
				}
				cur := ast.Node(x)
				for {
					pe, isParen := parents[cur].(*ast.ParenExpr)
					if !isParen || pe.X != cur {
						break
					}
					cur = pe
				}
				switch par := parents[cur].(type) {
				case *ast.UnaryExpr:
					if par.Op == token.AND && par.X == cur {
						return true
					}
				case *ast.SelectorExpr:
					if par.X == cur {
						return true
					}
				case *ast.ValueSpec:
					if len(par.Values) == 0 {
						for _, nm := range par.Names {
							if nm == cur {
								return true
							}
						}
					}
				}
				disqualifyStruct(named)
			}
			return true
		})
	}
}

type builderHandoffSpec struct {
	dst         *sliceCandidate
	src         *sliceCandidate
	builderType *types.Named
	initFunc    *types.Func
}

func validateAccFuncCallSite(
	id *ast.Ident,
	file *ast.File,
	ac *sliceCandidate,
	candidates map[*types.Var]*sliceCandidate,
	canTransform bool,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) {
	if !canTransform {
		ac.disqualified = true
		return
	}
	var call *ast.CallExpr
	switch p := parents[id].(type) {
	case *ast.CallExpr:
		if p.Fun == id && !p.Ellipsis.IsValid() {
			call = p
		}
	case *ast.SelectorExpr:
		if p.Sel == id {
			if c, ok := parents[p].(*ast.CallExpr); ok && c.Fun == p && !c.Ellipsis.IsValid() {
				call = c
			}
		}
	}
	if call == nil || ac.accParamIdx >= len(call.Args) {
		ac.disqualified = true
		return
	}
	argInner := unwrapCollectionAdaptersDown(call.Args[ac.accParamIdx], info)
	srcCand := resolveCandidateAccess(argInner, candidates, info)
	if srcCand == nil || srcCand.disqualified {
		ac.disqualified = true
		return
	}
	if enclosingFuncObj(call, parents, info) != ac.accFn {
		ac.inPkgAccCalls++
	}
	if ac.isExportedAccFn {
		ac.ops = append(ac.ops, sliceRewriteOp{kind: opAccCallRename, cand: ac, file: file, spanNode: id, elemIdent: id})
	}
}

func validateSliceCandidateUse(
	id *ast.Ident,
	file *ast.File,
	cand *sliceCandidate,
	candidates map[*types.Var]*sliceCandidate,
	accCandsByFn map[*types.Func][]*sliceCandidate,
	consumedCopyInitSrcs map[ast.Expr]bool,
	hasCopyHelpers bool,
	ctx *PassContext,
	parents map[ast.Node]ast.Node,
	varUseCounts map[*types.Var]int,
	varMaxUsePos map[*types.Var]token.Pos,
	markMut func(*sliceCandidate, ast.Node, map[ast.Node]ast.Node),
	handoffs *[]handoffEdge,
	builderChecks *[]builderHandoffSpec,
) {
	info := ctx.Info
	var accessExpr ast.Expr = id
	if cand.isField {
		sel, ok := parents[id].(*ast.SelectorExpr)
		if !ok || sel.Sel != id {
			cand.disqualified = true
			return
		}
		accessExpr = sel
	}
	if consumedCopyInitSrcs[accessExpr] {
		return
	}

	// Check LHS assignments: `accessExpr = rhs` (or `_ = accessExpr`).
	if as, ok := parents[accessExpr].(*ast.AssignStmt); ok {
		for _, lhs := range as.Lhs {
			if stripParens(lhs) == accessExpr {
				if cand.isAccParam || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 || (cand.isField && !isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info)) {
					cand.disqualified = true
					return
				}
				rhs := as.Rhs[0]
				if isNilIdent(rhs, info) {
					cand.hasNilAssign = true
					cand.ops = append(cand.ops, sliceRewriteOp{kind: opAssignNil, cand: cand, file: file, spanNode: as, stmt: as})
					markMut(cand, as, parents)
					return
				}
				// Self-assignment `cand = Lists_Reverse(cand)` on a local slice!
				if !cand.isField && j2gortHasExport(ctx, "SliceReverseInPlace") && isListsReverseOfCandidate(rhs, cand, info) {
					cand.hasSet = true
					cand.derefChecks = append(cand.derefChecks, derefGuardCheck{node: as, accessExpr: accessExpr, file: file})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
					cand.ops = append(cand.ops, sliceRewriteOp{
						kind: opAssignReverseSelf, cand: cand, file: file, spanNode: as, stmt: as, accessExpr: accessExpr,
					})
					markMut(cand, as, parents)
					return
				}
				if initK, argExpr, otherCand, typedSliceElem, initOK := classifySliceInitExpr(rhs, cand, candidates, hasCopyHelpers, ctx, info); initOK {
					if initK != sliceInitEmpty && !(initK == sliceInitCap && isNonNegativeConstInt(argExpr, info)) {
						cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: rhs.Pos()})
					}
					if initK == sliceInitCopySlice && otherCand != nil {
						cand.hasCopyInit = true
						consumedCopyInitSrcs[stripParens(argExpr)] = true
						otherCand.derefChecks = append(otherCand.derefChecks, derefGuardCheck{node: as, accessExpr: argExpr, file: file})
						*handoffs = append(*handoffs, handoffEdge{dst: cand, src: otherCand, propagateNil: false})
					} else if initK == sliceInitCopyTypedSlice && typedSliceElem != nil {
						cand.hasCopyInit = true
						if !cand.addWriteType(typedSliceElem, accessExpr, info) {
							cand.disqualified = true
							return
						}
					} else if initK == sliceInitCopyCol || initK == sliceInitCopyIter {
						cand.hasCopyInit = true
					}
					cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
					cand.ops = append(cand.ops, sliceRewriteOp{
						kind: opAssignInit, initKind: initK, cand: cand, otherCand: otherCand,
						file: file, spanNode: as, stmt: as, accessExpr: accessExpr, otherExpr: argExpr,
					})
					if !cand.isField || !isInInitFieldsMethod(as, cand.ownerStruct, parents, info) {
						markMut(cand, as, parents)
					}
					return
				}
				// Linear ownership handoff: `destCand = srcCand` or `destCand = UnmodifiableList(srcCand)`.
				unwrappedRhs := unwrapUnmodifiableListCall(rhs, info)
				if srcCand := resolveCandidateAccess(unwrappedRhs, candidates, info); srcCand != nil && srcCand != cand {
					if isValidSliceHandoff(as, cand, srcCand, unwrappedRhs, parents, info, varMaxUsePos, builderChecks) {
						*handoffs = append(*handoffs, handoffEdge{dst: cand, src: srcCand, propagateNil: true})
						if unwrappedRhs != stripParens(rhs) {
							cand.ops = append(cand.ops, sliceRewriteOp{
								kind: opExprRawUnbox, cand: cand, file: file, spanNode: rhs, outerExpr: rhs, otherExpr: unwrappedRhs,
							})
						}
						if !cand.isField || !isInInitFieldsMethod(as, cand.ownerStruct, parents, info) {
							markMut(cand, as, parents)
						}
						return
					}
				}
				cand.disqualified = true
				return
			}
		}
		if as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && stripParens(as.Rhs[0]) == accessExpr {
			if lhsId, ok := stripParens(as.Lhs[0]).(*ast.Ident); ok && lhsId.Name == "_" {
				return
			}
			// RHS of a linear handoff `destCand = cand` is validated when visiting destCand; verify here too.
			if dstCand := resolveCandidateAccess(as.Lhs[0], candidates, info); dstCand != nil && dstCand != cand && !dstCand.disqualified {
				var dummy []builderHandoffSpec
				if isValidSliceHandoff(as, dstCand, cand, accessExpr, parents, info, varMaxUsePos, &dummy) {
					*handoffs = append(*handoffs, handoffEdge{dst: dstCand, src: cand, propagateNil: true})
					return
				}
			}
		}
	}

	// Allow `x == nil` and `x != nil` directly on lowered slices (not on accumulator params).
	curExpr := ast.Expr(accessExpr)
	for {
		pe, isParen := parents[curExpr].(*ast.ParenExpr)
		if !isParen || pe.X != curExpr {
			break
		}
		curExpr = pe
	}
	if !cand.isAccParam || !cand.isExportedAccFn {
		if bin, isBin := parents[curExpr].(*ast.BinaryExpr); isBin && (bin.Op == token.EQL || bin.Op == token.NEQ) {
			other := bin.Y
			if bin.Y == curExpr {
				other = bin.X
			}
			if isNilIdent(other, info) {
				return
			}
		}
	}
	if !cand.isAccParam {
		// Allow `javalang.ToJavaString(cand)` -> `j2gort.SliceFormatJavaString(cand)`.
		if strCall, ok := parents[curExpr].(*ast.CallExpr); ok && len(strCall.Args) == 1 && !strCall.Ellipsis.IsValid() && strCall.Args[0] == curExpr && j2gortHasExport(ctx, "SliceFormatJavaString") {
			if callee, _ := calleeFunc(strCall, info); callee != nil && callee.Pkg() != nil && callee.Pkg().Name() == "javalang" && callee.Name() == "ToJavaString" {
				if _, isES := parents[strCall].(*ast.ExprStmt); !isES {
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: strCall.Pos()})
					cand.ops = append(cand.ops, sliceRewriteOp{
						kind: opExprFormatJavaStr, cand: cand, file: file, spanNode: strCall, outerExpr: strCall, accessExpr: accessExpr,
					})
					return
				}
			}
		}

		// Allow read-only string-formatting sinks (e.g. ReferenceCollection.String() ToStringHelper).
		if isFormattingOnlySliceSink(curExpr, parents, info) {
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: accessExpr.Pos()})
			cand.ops = append(cand.ops, sliceRewriteOp{
				kind: opExprToJavaList, cand: cand, file: file, spanNode: accessExpr, outerExpr: accessExpr, accessExpr: accessExpr,
			})
			return
		}
	}

	recordDeref := func(node ast.Node) {
		if !cand.isAccParam {
			cand.derefChecks = append(cand.derefChecks, derefGuardCheck{node: node, accessExpr: accessExpr, file: file})
		}
	}

	// Check static j2gort / helper consumers wrapping `AsIface`/`ToIface_*`/`CheckNotNull*`.
	adaptedTop := unwrapCollectionAdapters(accessExpr, parents, info)
	if staticCall, ok := parents[adaptedTop].(*ast.CallExpr); ok && !staticCall.Ellipsis.IsValid() && staticCall.Fun != adaptedTop {
		// Check if passed as accumulator parameter to a candidate slice-accumulator function/method.
		if callee, _ := calleeFunc(staticCall, info); callee != nil {
			for _, ac := range accCandsByFn[funcOrigin(callee)] {
				if ac != nil && ac.accParamIdx < len(staticCall.Args) && staticCall.Args[ac.accParamIdx] == adaptedTop {
					if cand.isField && !isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info) {
						cand.disqualified = true
						ac.disqualified = true
						return
					}
					recordDeref(staticCall)
					*handoffs = append(*handoffs, handoffEdge{dst: cand, src: ac, propagateNil: false})
					cand.ops = append(cand.ops, sliceRewriteOp{
						kind: opExprAccParamArg, cand: cand, otherCand: ac, file: file, spanNode: adaptedTop,
						call: staticCall, outerExpr: adaptedTop, accessExpr: accessExpr,
					})
					return
				}
			}
		}
		if cand.isAccParam && cand.isExportedAccFn {
			cand.disqualified = true
			return
		}

		if fnName, _, _, isJ2 := isJ2gortCall(staticCall, info); isJ2 {
			switch {
			case len(staticCall.Args) == 1 && staticCall.Args[0] == adaptedTop:
				switch fnName {
				case "NewIterCursor", "NewTypedIterCursor":
					if topCheck := unwrapEnclosingCheckNotNulls(accessExpr, parents, info); topCheck == adaptedTop {
						recordDeref(staticCall)
						validateIterCursorLoop(staticCall, accessExpr, false, file, cand, parents, info, varUseCounts)
						return
					}
				case "Ext_com_google_common_collect_Lists_Reverse":
					if cand.isAccParam {
						break
					}
					// Pattern 1: RHS of self-reassignment `cand = Lists_Reverse(cand)` (handled when visiting LHS).
					revTop := unwrapParensUp(staticCall, parents)
					if as, isAs := parents[revTop].(*ast.AssignStmt); isAs && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && as.Rhs[0] == revTop && sameAccessPath(as.Lhs[0], accessExpr, cand.v, info) && !cand.isField && j2gortHasExport(ctx, "SliceReverseInPlace") {
						return
					}
					// Pattern 2: `cur := NewIterCursor(CheckNotNull(Lists_Reverse(cand))); for cur.HasNext() { ... }`.
					revAdapted := unwrapCollectionAdapters(staticCall, parents, info)
					if iterCall, isIter := parents[revAdapted].(*ast.CallExpr); isIter && len(iterCall.Args) == 1 && !iterCall.Ellipsis.IsValid() && iterCall.Args[0] == revAdapted {
						if iterFn, _, _, iterJ2 := isJ2gortCall(iterCall, info); iterJ2 && (iterFn == "NewIterCursor" || iterFn == "NewTypedIterCursor") && unwrapEnclosingCheckNotNulls(staticCall, parents, info) == revAdapted && isUnshadowedBuiltin("len", iterCall.Pos(), ctx.Pkg) {
							recordDeref(iterCall)
							validateIterCursorLoop(iterCall, accessExpr, true, file, cand, parents, info, varUseCounts)
							return
						}
					}
				case "Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection",
					"Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable":
					if !cand.isAccParam {
						if _, isES := parents[staticCall].(*ast.ExprStmt); !isES {
							recordDeref(staticCall)
							cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: staticCall.Pos()})
							cand.ops = append(cand.ops, sliceRewriteOp{
								kind: opExprImmListCopy, cand: cand, file: file, spanNode: staticCall, outerExpr: staticCall, accessExpr: accessExpr,
							})
							return
						}
					}
				case "ImmListCopyToSlice", "ImmListCopyNodeToSlice":
					if !cand.isAccParam {
						if _, isES := parents[staticCall].(*ast.ExprStmt); !isES {
							if sl, isSl := types.Unalias(info.Types[staticCall].Type).(*types.Slice); isSl && cand.addReadType(sl.Elem(), accessExpr, info) {
								recordDeref(staticCall)
								cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: staticCall.Pos()})
								cand.ops = append(cand.ops, sliceRewriteOp{
									kind: opExprImmCopyToSlice, cand: cand, file: file, spanNode: staticCall, outerExpr: staticCall, accessExpr: accessExpr,
								})
								return
							}
						}
					}
				case "Ext_java_util_Collections_UnmodifiableList":
					if cand.isAccParam {
						break
					}
					// Check if RHS of builder handoff `this_.f = UnmodifiableList(builder.f)`.
					uTop := unwrapParensUp(staticCall, parents)
					if as, isAs := parents[uTop].(*ast.AssignStmt); isAs && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && as.Rhs[0] == uTop {
						if dstCand := resolveCandidateAccess(as.Lhs[0], candidates, info); dstCand != nil && dstCand != cand {
							var dummy []builderHandoffSpec
							if isValidSliceHandoff(as, dstCand, cand, accessExpr, parents, info, varMaxUsePos, &dummy) {
								return
							}
						}
					}
					if _, isES := parents[staticCall].(*ast.ExprStmt); !isES {
						if !cand.isField {
							recordDeref(staticCall)
							cand.unmodChecks = append(cand.unmodChecks, unmodLocalCheck{callNode: staticCall, usePos: accessExpr.Pos(), file: file})
							cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: staticCall.Pos()})
							cand.ops = append(cand.ops, sliceRewriteOp{
								kind: opExprUnmodList, cand: cand, file: file, spanNode: staticCall, outerExpr: staticCall, accessExpr: accessExpr,
							})
							return
						}
						if isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info) && j2gortHasExport(ctx, "SliceUnmodifiableListLive") {
							cand.hasLiveUnmod = true
							recordDeref(staticCall)
							cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: staticCall.Pos()})
							cand.ops = append(cand.ops, sliceRewriteOp{
								kind: opExprUnmodListLive, cand: cand, file: file, spanNode: staticCall, outerExpr: staticCall, accessExpr: accessExpr,
							})
							return
						}
					}
				case "Ext_java_util_Collections_Sort_1":
					if !cand.isAccParam {
						if es, isES := parents[staticCall].(*ast.ExprStmt); isES && es.X == staticCall && j2gortHasExport(ctx, "SliceSortNatural") {
							cand.hasSet = true
							recordDeref(es)
							cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: staticCall.Pos()})
							cand.ops = append(cand.ops, sliceRewriteOp{
								kind: opStmtSortNat, cand: cand, file: file, spanNode: es, stmt: es, call: staticCall, accessExpr: accessExpr,
							})
							markMut(cand, staticCall, parents)
							return
						}
					}
				case "Ext_com_google_common_collect_Iterables_GetLast_1":
					outerExpr, expectedT, ok := enclosingUnboxExpr(staticCall, parents, info)
					_, outerES := parents[outerExpr].(*ast.ExprStmt)
					if ok && !outerES && cand.addReadType(expectedT, accessExpr, info) {
						recordDeref(outerExpr)
						cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: staticCall.Pos()})
						cand.ops = append(cand.ops, sliceRewriteOp{
							kind: opExprGetLast, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr, accessExpr: accessExpr,
						})
						return
					}
				case "NewExt_java_util_ArrayList_1_Ext_java_util_Collection",
					"Ext_com_google_common_collect_Lists_NewArrayList_1_Ext_java_lang_Iterable":
					if !cand.isAccParam && isSliceCopyTargetCandidate(staticCall, candidates, parents, info) {
						return
					}
				}
			case len(staticCall.Args) == 2 && staticCall.Args[0] == adaptedTop:
				switch fnName {
				case "Ext_java_util_Collections_Sort_2":
					if !cand.isAccParam {
						if es, isES := parents[staticCall].(*ast.ExprStmt); isES && es.X == staticCall && j2gortHasExport(ctx, "SliceSortComparator") {
							cand.hasSet = true
							recordDeref(es)
							cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: staticCall.Pos()})
							cand.ops = append(cand.ops, sliceRewriteOp{
								kind: opStmtSortCmp, cand: cand, file: file, spanNode: es, stmt: es, call: staticCall, accessExpr: accessExpr, otherExpr: staticCall.Args[1],
							})
							markMut(cand, staticCall, parents)
							return
						}
					}
				case "ListGetStr", "ListGetByteString", "ListGetVal":
					_, outerES := parents[staticCall].(*ast.ExprStmt)
					readT := info.Types[staticCall].Type
					if fnName == "ListGetStr" {
						readT = types.Typ[types.String]
					}
					if !outerES && !isInvalidOrUntyped(readT) && isIntegerExpr(staticCall.Args[1], info) && cand.addReadType(readT, accessExpr, info) {
						recordDeref(staticCall)
						cand.ops = append(cand.ops, sliceRewriteOp{
							kind: opExprGet, cand: cand, file: file, spanNode: staticCall, call: staticCall, outerExpr: staticCall, accessExpr: accessExpr, otherExpr: staticCall.Args[1],
						})
						return
					}
				}
			}
		}

		if cand.isAccParam {
			cand.disqualified = true
			return
		}

		// Read-only Iterable consumers: Guava Joiner.Join_1_Ext_java_lang_Iterable / AppendTo and rhino TypeStringBuilder.AppendAll.
		if isReadOnlyIterableConsumerCall(staticCall, adaptedTop, info) && j2gortHasExport(ctx, "SliceReadOnlyIterable") {
			recordDeref(staticCall)
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: accessExpr.Pos()})
			cand.ops = append(cand.ops, sliceRewriteOp{
				kind: opExprReadOnlyIter, cand: cand, file: file, spanNode: adaptedTop, outerExpr: adaptedTop, accessExpr: accessExpr,
			})
			return
		}

		// Also check if `cand` is the RHS argument of another candidate's `dst.AddAll_1(AsIface(cand))`!
		if len(staticCall.Args) == 1 && staticCall.Args[0] == adaptedTop {
			if mSel, ok := stripParens(staticCall.Fun).(*ast.SelectorExpr); ok && mSel.Sel.Name == "AddAll_1" {
				if _, isES := parents[staticCall].(*ast.ExprStmt); isES {
					if dstCand := resolveCandidateAccess(unwrapCheckNotNullArg(mSel.X, info), candidates, info); dstCand != nil && dstCand != cand && !dstCand.disqualified {
						*handoffs = append(*handoffs, handoffEdge{dst: dstCand, src: cand})
						return
					}
				}
			}
		}
	}

	if !cand.isField && !cand.isAccParam && j2gortHasExport(ctx, "WrapTypedArrayList") {
		if ret, isRet := parents[adaptedTop].(*ast.ReturnStmt); isRet && validateLocalReturnWrapTypedArrayList(ret, adaptedTop, cand, parents, info) {
			if unwrapEnclosingCheckNotNulls(accessExpr, parents, info) != unwrapParensUp(accessExpr, parents) {
				recordDeref(ret)
			}
			cand.hasWrapTypedArrayList = true
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: adaptedTop.Pos()})
			cand.ops = append(cand.ops, sliceRewriteOp{
				kind: opExprWrapTypedArrayList, cand: cand, file: file, spanNode: adaptedTop, outerExpr: adaptedTop, accessExpr: accessExpr,
			})
			return
		}
	}

	topExpr := unwrapEnclosingCheckNotNulls(accessExpr, parents, info)
	parent := parents[topExpr]

	methodSel, selOK := parent.(*ast.SelectorExpr)
	if !selOK || methodSel.X != topExpr {
		cand.disqualified = true
		return
	}
	mCall, callOK := parents[methodSel].(*ast.CallExpr)
	if !callOK || mCall.Fun != methodSel || mCall.Ellipsis.IsValid() {
		cand.disqualified = true
		return
	}
	_, inExprStmt := parents[mCall].(*ast.ExprStmt)

	if cand.isAccParam && cand.isExportedAccFn && methodSel.Sel.Name != "Add_1" && methodSel.Sel.Name != "AddLast" {
		cand.disqualified = true
		return
	}

	switch methodSel.Sel.Name {
	case "Add_1", "AddLast":
		es, ok := parents[mCall].(*ast.ExprStmt)
		if !ok || es.X != mCall || len(mCall.Args) != 1 || !isUnshadowedBuiltin("append", mCall.Pos(), ctx.Pkg) {
			cand.disqualified = true
			return
		}
		if cand.isField {
			sel := accessExpr.(*ast.SelectorExpr)
			if !isPureFieldSelector(sel, info) || hasClosureMutation(mCall.Args[0]) {
				cand.disqualified = true
				return
			}
		}
		if unboxed, elemT, ok := unwrapBoxedAddArg(mCall.Args[0], info); ok && cand.addWriteType(elemT, accessExpr, info) {
			cand.hasAdd = true
			recordDeref(es)
			markMut(cand, mCall, parents)
			cand.ops = append(cand.ops, sliceRewriteOp{
				kind: opStmtAdd, cand: cand, file: file, spanNode: es, stmt: es,
				call: mCall, accessExpr: accessExpr, unwrapsBox: unboxed != stripParens(mCall.Args[0]),
			})
		} else {
			cand.disqualified = true
		}

	case "Add_2", "AddFirst":
		es, ok := parents[mCall].(*ast.ExprStmt)
		isFirst := methodSel.Sel.Name == "AddFirst"
		if !ok || es.X != mCall || !j2gortHasExport(ctx, "SliceInsertAt") {
			cand.disqualified = true
			return
		}
		var idxExpr, elemArg ast.Expr
		if isFirst {
			if len(mCall.Args) != 1 {
				cand.disqualified = true
				return
			}
			idxExpr = &ast.BasicLit{ValuePos: mCall.Pos(), Kind: token.INT, Value: "0"}
			elemArg = mCall.Args[0]
		} else {
			if len(mCall.Args) != 2 || !isIntegerExpr(mCall.Args[0], info) {
				cand.disqualified = true
				return
			}
			idxExpr = mCall.Args[0]
			elemArg = mCall.Args[1]
		}
		if cand.isField && (!isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info) || hasClosureMutation(idxExpr) || hasClosureMutation(elemArg)) {
			cand.disqualified = true
			return
		}
		if unboxed, elemT, ok := unwrapBoxedAddArg(elemArg, info); ok && cand.addWriteType(elemT, accessExpr, info) {
			cand.hasAdd = true
			cand.hasInsertAt = true
			recordDeref(es)
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
			cand.ops = append(cand.ops, sliceRewriteOp{
				kind: opStmtInsertAt, cand: cand, file: file, spanNode: es, stmt: es,
				call: mCall, accessExpr: accessExpr, otherExpr: idxExpr, unwrapsBox: unboxed != stripParens(elemArg),
			})
			markMut(cand, mCall, parents)
		} else {
			cand.disqualified = true
		}

	case "Sort":
		if cand.isAccParam {
			cand.disqualified = true
			return
		}
		es, ok := parents[mCall].(*ast.ExprStmt)
		if !ok || es.X != mCall || len(mCall.Args) != 1 || !j2gortHasExport(ctx, "SliceSortComparator") {
			cand.disqualified = true
			return
		}
		if cand.isField && (!isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info) || hasClosureMutation(mCall.Args[0])) {
			cand.disqualified = true
			return
		}
		cand.hasSet = true
		recordDeref(es)
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, sliceRewriteOp{
			kind: opStmtSortCmp, cand: cand, file: file, spanNode: es, stmt: es, call: mCall, accessExpr: accessExpr, otherExpr: mCall.Args[0],
		})
		markMut(cand, mCall, parents)

	case "Set", "Set_2":
		es, ok := parents[mCall].(*ast.ExprStmt)
		if !ok || es.X != mCall || len(mCall.Args) != 2 || !isIntegerExpr(mCall.Args[0], info) {
			cand.disqualified = true
			return
		}
		if cand.isField && (!isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info) || hasClosureMutation(mCall.Args[1])) {
			cand.disqualified = true
			return
		}
		if unboxed, elemT, ok := unwrapBoxedAddArg(mCall.Args[1], info); ok && cand.addWriteType(elemT, accessExpr, info) {
			cand.hasSet = true
			recordDeref(es)
			cand.ops = append(cand.ops, sliceRewriteOp{
				kind: opStmtSet, cand: cand, file: file, spanNode: es, stmt: es,
				call: mCall, accessExpr: accessExpr, unwrapsBox: unboxed != stripParens(mCall.Args[1]),
			})
			markMut(cand, mCall, parents)
		} else {
			cand.disqualified = true
		}

	case "Remove_1_int32":
		if len(mCall.Args) != 1 || !isIntegerExpr(mCall.Args[0], info) {
			cand.disqualified = true
			return
		}
		if cand.isField && (!isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info) || hasClosureMutation(mCall.Args[0])) {
			cand.disqualified = true
			return
		}
		cand.hasRemove = true
		if es, ok := parents[mCall].(*ast.ExprStmt); ok && es.X == mCall {
			recordDeref(es)
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
			cand.ops = append(cand.ops, sliceRewriteOp{
				kind: opStmtRemoveAt, cand: cand, file: file, spanNode: es, stmt: es, call: mCall, accessExpr: accessExpr,
			})
			markMut(cand, mCall, parents)
			return
		}
		outer, elemT, ok := enclosingUnboxExpr(mCall, parents, info)
		if !ok || !cand.addReadType(elemT, accessExpr, info) {
			cand.disqualified = true
			return
		}
		if _, outerES := parents[outer].(*ast.ExprStmt); outerES {
			cand.disqualified = true
			return
		}
		recordDeref(outer)
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outer.Pos()})
		cand.ops = append(cand.ops, sliceRewriteOp{
			kind: opExprRemoveAt, cand: cand, file: file, spanNode: outer, outerExpr: outer, call: mCall, accessExpr: accessExpr,
		})
		markMut(cand, mCall, parents)

	case "RemoveLast", "RemoveFirst":
		es, ok := parents[mCall].(*ast.ExprStmt)
		isLast := methodSel.Sel.Name == "RemoveLast"
		needHelper := "SliceGetFirst"
		if isLast {
			needHelper = "SliceGetLast"
		}
		if !ok || es.X != mCall || len(mCall.Args) != 0 || !canReplaceStmtIn(parents[es]) || !j2gortHasExport(ctx, needHelper) {
			cand.disqualified = true
			return
		}
		if isLast && !isUnshadowedBuiltin("len", mCall.Pos(), ctx.Pkg) {
			cand.disqualified = true
			return
		}
		if !isLast && !j2gortHasExport(ctx, "SliceRemoveAt") {
			cand.disqualified = true
			return
		}
		if cand.isField && !isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info) {
			cand.disqualified = true
			return
		}
		cand.hasRemove = true
		recordDeref(es)
		if isLast {
			cand.clearPositions = append(cand.clearPositions, mCall.Pos())
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		kind := opStmtRemoveFirst
		if isLast {
			kind = opStmtRemoveLast
		}
		cand.ops = append(cand.ops, sliceRewriteOp{
			kind: kind, cand: cand, file: file, spanNode: es, stmt: es, call: mCall, accessExpr: accessExpr,
		})
		markMut(cand, mCall, parents)

	case "AddAll_1":
		if cand.isAccParam {
			cand.disqualified = true
			return
		}
		es, ok := parents[mCall].(*ast.ExprStmt)
		if !ok || es.X != mCall || len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		if cand.isField && (!isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info) || hasClosureMutation(mCall.Args[0])) {
			cand.disqualified = true
			return
		}
		unwrappedArg := unwrapCollectionAdaptersDown(mCall.Args[0], info)
		if srcCand := resolveCandidateAccess(unwrappedArg, candidates, info); srcCand != nil && srcCand != cand {
			if !isUnshadowedBuiltin("append", mCall.Pos(), ctx.Pkg) {
				cand.disqualified = true
				return
			}
			recordDeref(es)
			srcCand.derefChecks = append(srcCand.derefChecks, derefGuardCheck{node: es, accessExpr: unwrappedArg, file: file})
			*handoffs = append(*handoffs, handoffEdge{dst: cand, src: srcCand})
			cand.ops = append(cand.ops, sliceRewriteOp{
				kind: opStmtAddAllSlice, cand: cand, otherCand: srcCand, file: file, spanNode: es, stmt: es,
				call: mCall, accessExpr: accessExpr, otherExpr: unwrappedArg,
			})
			markMut(cand, mCall, parents)
			return
		}
		recordDeref(es)
		cand.hasAdd = true
		cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: mCall.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, sliceRewriteOp{
			kind: opStmtAddAllCol, cand: cand, file: file, spanNode: es, stmt: es, call: mCall, accessExpr: accessExpr,
		})
		markMut(cand, mCall, parents)

	case "Get":
		outerExpr, expectedT, ok := enclosingUnboxExpr(mCall, parents, info)
		_, outerES := parents[outerExpr].(*ast.ExprStmt)
		if !ok || outerES || len(mCall.Args) != 1 || !isIntegerExpr(mCall.Args[0], info) {
			cand.disqualified = true
			return
		}
		if cand.addReadType(expectedT, accessExpr, info) {
			recordDeref(outerExpr)
			cand.ops = append(cand.ops, sliceRewriteOp{kind: opExprGet, cand: cand, file: file, spanNode: outerExpr, call: mCall, outerExpr: outerExpr, accessExpr: accessExpr})
		}

	case "GetFirst", "GetLast":
		outerExpr, expectedT, ok := enclosingUnboxExpr(mCall, parents, info)
		_, outerES := parents[outerExpr].(*ast.ExprStmt)
		if !ok || outerES || len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		kind := opExprGetFirst
		if methodSel.Sel.Name == "GetLast" {
			kind = opExprGetLast
		}
		if cand.addReadType(expectedT, accessExpr, info) {
			recordDeref(outerExpr)
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
			cand.ops = append(cand.ops, sliceRewriteOp{kind: kind, cand: cand, file: file, spanNode: outerExpr, call: mCall, outerExpr: outerExpr, accessExpr: accessExpr})
		}

	case "Iterator":
		if inExprStmt || len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		outer := unwrapIteratorCast(mCall, parents, info)
		topIterCheck := unwrapEnclosingCheckNotNulls(outer, parents, info)
		if cursorCall, isCall := parents[topIterCheck].(*ast.CallExpr); isCall && len(cursorCall.Args) == 1 && !cursorCall.Ellipsis.IsValid() && cursorCall.Args[0] == topIterCheck {
			if fnName, _, _, isJ2 := isJ2gortCall(cursorCall, info); isJ2 && (fnName == "NewIterCursor" || fnName == "NewTypedIterCursor") {
				recordDeref(cursorCall)
				validateIterCursorLoop(cursorCall, accessExpr, false, file, cand, parents, info, varUseCounts)
				return
			}
		}
		if cand.isAccParam || !isSafeReadOnlySliceIteratorUse(outer, cand, ctx, parents, info) {
			cand.disqualified = true
			return
		}
		cand.hasIterCall = true
		recordDeref(outer)
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, sliceRewriteOp{
			kind: opExprSliceIter, cand: cand, file: file, spanNode: outer, call: mCall, outerExpr: outer, accessExpr: accessExpr,
		})

	case "Size":
		if inExprStmt || len(mCall.Args) != 0 || !isUnshadowedBuiltin("len", mCall.Pos(), ctx.Pkg) || !isUnshadowedUniverseType("int32", mCall.Pos(), ctx.Pkg) {
			cand.disqualified = true
			return
		}
		recordDeref(mCall)
		cand.ops = append(cand.ops, sliceRewriteOp{kind: opExprSize, cand: cand, file: file, spanNode: mCall, call: mCall, accessExpr: accessExpr})

	case "IsEmpty":
		if inExprStmt || len(mCall.Args) != 0 || !isUnshadowedBuiltin("len", mCall.Pos(), ctx.Pkg) {
			cand.disqualified = true
			return
		}
		recordDeref(mCall)
		cand.ops = append(cand.ops, sliceRewriteOp{kind: opExprIsEmpty, cand: cand, file: file, spanNode: mCall, call: mCall, accessExpr: accessExpr})

	case "Clear":
		es, ok := parents[mCall].(*ast.ExprStmt)
		if !ok || es.X != mCall || len(mCall.Args) != 0 || !canReplaceStmtIn(parents[es]) || (cand.isField && !isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info)) {
			cand.disqualified = true
			return
		}
		cand.hasClear = true
		recordDeref(es)
		cand.clearPositions = append(cand.clearPositions, mCall.Pos())
		cand.ops = append(cand.ops, sliceRewriteOp{kind: opStmtClear, cand: cand, file: file, spanNode: es, stmt: es, call: mCall, accessExpr: accessExpr})
		markMut(cand, mCall, parents)

	case "ForEach":
		if !inExprStmt || len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		es, boundAs, fnLit, paramName, paramPos, drop, elemTypes, unboxOps, ok := matchForEachConsumerLambda(mCall, parents, info, varUseCounts)
		if !ok || !canReplaceStmtIn(parents[es]) {
			cand.disqualified = true
			return
		}
		for _, et := range elemTypes {
			if !cand.addReadType(et, accessExpr, info) {
				cand.disqualified = true
				return
			}
		}
		for _, uo := range unboxOps {
			cand.ops = append(cand.ops, sliceRewriteOp{
				kind: opExprRawUnbox, cand: cand, file: file, spanNode: uo.outer, outerExpr: uo.outer, elemIdent: uo.id,
			})
		}
		recordDeref(es)
		cand.loopBodies = append(cand.loopBodies, fnLit.Body)
		cand.ops = append(cand.ops, sliceRewriteOp{
			kind: opStmtForEach, cand: cand, file: file, spanNode: es, stmt: es, accessExpr: accessExpr,
			forEachSpec: &setForEachSpec{
				exprStmt:  es,
				boundAs:   boundAs,
				fnLit:     fnLit,
				paramName: paramName,
				paramPos:  paramPos,
				dropStmts: drop,
			},
		})

	case "IterSlice_":
		rs, ok := parents[mCall].(*ast.RangeStmt)
		if !ok || rs.X != mCall || len(mCall.Args) != 0 || rs.Body == nil {
			cand.disqualified = true
			return
		}
		if rs.Key != nil {
			if kId, ok := rs.Key.(*ast.Ident); !ok || kId.Name != "_" {
				cand.disqualified = true
				return
			}
		}
		var elemId *ast.Ident
		drop := 0
		if rs.Value != nil {
			valId, ok := rs.Value.(*ast.Ident)
			if !ok || rs.Tok != token.DEFINE {
				cand.disqualified = true
				return
			}
			if valId.Name != "_" {
				rawVar := varFromObj(info.Defs[valId])
				var valid bool
				elemId, drop, valid = validateLoopRawElemVar(rawVar, accessExpr, file, rs.Body.List, cand, parents, info, varUseCounts)
				if rawVar == nil || !valid {
					cand.disqualified = true
					return
				}
			}
		}
		recordDeref(rs)
		cand.loopBodies = append(cand.loopBodies, rs.Body)
		cand.ops = append(cand.ops, sliceRewriteOp{kind: opStmtIterSlice, cand: cand, file: file, spanNode: rs, stmt: rs, accessExpr: accessExpr, elemIdent: elemId, dropStmts: drop})

	default:
		cand.disqualified = true
	}
}

func validateIterCursorLoop(cursorCall *ast.CallExpr, accessExpr ast.Expr, isReverse bool, file *ast.File, cand *sliceCandidate, parents map[ast.Node]ast.Node, info *types.Info, varUseCounts map[*types.Var]int) {
	initAs, ok := parents[cursorCall].(*ast.AssignStmt)
	if !ok || initAs.Tok != token.DEFINE || len(initAs.Lhs) != 1 || len(initAs.Rhs) != 1 || initAs.Rhs[0] != cursorCall {
		cand.disqualified = true
		return
	}
	curId, _ := initAs.Lhs[0].(*ast.Ident)
	curVar := varFromObj(info.Defs[curId])
	block, _ := parents[initAs].(*ast.BlockStmt)
	if curId == nil || curId.Name == "_" || curVar == nil || varUseCounts[curVar] != 2 || block == nil {
		cand.disqualified = true
		return
	}
	idx := slices.Index(block.List, ast.Stmt(initAs))
	if idx < 0 || idx+1 >= len(block.List) {
		cand.disqualified = true
		return
	}
	nextStmt := block.List[idx+1]
	var label *ast.Ident
	if lbl, isLbl := nextStmt.(*ast.LabeledStmt); isLbl {
		label, nextStmt = lbl.Label, lbl.Stmt
	}
	forStmt, ok := nextStmt.(*ast.ForStmt)
	if !ok || forStmt.Init != nil || forStmt.Post != nil || forStmt.Cond == nil || forStmt.Body == nil || len(forStmt.Body.List) == 0 || !isCursorMethodCall(forStmt.Cond, curVar, "HasNext", info) {
		cand.disqualified = true
		return
	}
	s0, ok := forStmt.Body.List[0].(*ast.AssignStmt)
	if !ok || s0.Tok != token.DEFINE || len(s0.Lhs) != 1 || len(s0.Rhs) != 1 {
		cand.disqualified = true
		return
	}
	lhs0, _ := s0.Lhs[0].(*ast.Ident)
	if lhs0 == nil || lhs0.Name == "_" || info.Defs[lhs0] == nil {
		cand.disqualified = true
		return
	}

	elemId := lhs0
	drop := 1
	rhs0 := stripParens(s0.Rhs[0])
	switch {
	case isCursorMethodCall(rhs0, curVar, "NextStr", info):
		cand.addReadType(types.Typ[types.String], accessExpr, info)
	case isCursorMethodCall(rhs0, curVar, "NextNode", info):
		cand.addReadType(types.Typ[types.Uint32], accessExpr, info)
	case isNodeConvOfCursorNextNode(rhs0, curVar, info):
		cand.addReadType(info.Types[rhs0].Type, accessExpr, info)
	case isCursorMethodCall(rhs0, curVar, "Next", info):
		if !isEmptyInterfaceType(info.Types[rhs0].Type) {
			if !cand.addReadType(info.Types[rhs0].Type, accessExpr, info) {
				cand.disqualified = true
				return
			}
		} else {
			rawVar := varFromObj(info.Defs[lhs0])
			unboxedId, extraDrop, valid := validateLoopRawElemVar(rawVar, accessExpr, file, forStmt.Body.List[1:], cand, parents, info, varUseCounts)
			if rawVar == nil || !valid {
				cand.disqualified = true
				return
			}
			if unboxedId != nil {
				elemId = unboxedId
			}
			drop += extraDrop
		}
	default:
		cand.disqualified = true
		return
	}

	cand.loopBodies = append(cand.loopBodies, forStmt.Body)
	cand.ops = append(cand.ops, sliceRewriteOp{
		kind: opStmtIterCursor, cand: cand, file: file, spanNode: forStmt,
		stmt: initAs, forStmt: forStmt, block: block, label: label, accessExpr: accessExpr,
		elemIdent: elemId, dropStmts: drop, isReverseIter: isReverse,
	})
}

func validateLoopRawElemVar(rawVar *types.Var, accessExpr ast.Expr, file *ast.File, stmts []ast.Stmt, cand *sliceCandidate, parents map[ast.Node]ast.Node, info *types.Info, varUseCounts map[*types.Var]int) (*ast.Ident, int, bool) {
	if rawVar == nil {
		return nil, 0, false
	}
	if len(stmts) >= 1 && varUseCounts[rawVar] == 1 {
		if elemId, elemT, ok := matchSingleDefineUnbox(stmts[0], rawVar, parents, info); ok {
			return elemId, 1, cand.addReadType(elemT, accessExpr, info)
		}
	}
	if varUseCounts[rawVar] == 0 {
		return nil, 0, true
	}
	valid := true
	for _, st := range stmts {
		ast.Inspect(st, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !valid || !ok || varFromObj(info.Uses[id]) != rawVar {
				return valid
			}
			if as, isAs := parents[id].(*ast.AssignStmt); isAs && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && as.Rhs[0] == id {
				if lhsId, ok := stripParens(as.Lhs[0]).(*ast.Ident); ok && lhsId.Name == "_" {
					return true
				}
			}
			outerExpr, expectedT, ok := enclosingUnboxExpr(id, parents, info)
			_, isES := parents[outerExpr].(*ast.ExprStmt)
			if !ok || isES || !cand.addReadType(expectedT, accessExpr, info) {
				valid = false
				return false
			}
			if outerExpr != id {
				cand.ops = append(cand.ops, sliceRewriteOp{kind: opExprRawUnbox, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr, elemIdent: id})
			}
			return true
		})
	}
	return nil, 0, valid
}

func matchSingleDefineUnbox(stmt ast.Stmt, rawVar *types.Var, parents map[ast.Node]ast.Node, info *types.Info) (*ast.Ident, types.Type, bool) {
	as, ok := stmt.(*ast.AssignStmt)
	if !ok || as.Tok != token.DEFINE || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
		return nil, nil, false
	}
	elemId, _ := as.Lhs[0].(*ast.Ident)
	if elemId == nil || elemId.Name == "_" || info.Defs[elemId] == nil {
		return nil, nil, false
	}
	var foundRaw *ast.Ident
	ast.Inspect(as.Rhs[0], func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && info.Uses[id] == rawVar {
			foundRaw = id
			return false
		}
		return true
	})
	if foundRaw == nil {
		return nil, nil, false
	}
	outerExpr, expectedT, ok := enclosingUnboxExpr(foundRaw, parents, info)
	if !ok || stripParens(as.Rhs[0]) != stripParens(outerExpr) {
		return nil, nil, false
	}
	return elemId, expectedT, true
}

func anyLoopMutatesCandidate(cand *sliceCandidate, mutatingFuncs map[*types.Func]bool, info *types.Info) bool {
	for _, body := range cand.loopBodies {
		mutated := false
		ast.Inspect(body, func(n ast.Node) bool {
			if mutated {
				return false
			}
			switch x := n.(type) {
			case *ast.CallExpr:
				if len(mutatingFuncs) > 0 {
					if fn, _ := calleeFunc(x, info); fn != nil && mutatingFuncs[fn] {
						mutated = true
						return false
					}
				}
				if fnName, _, _, isJ2 := isJ2gortCall(x, info); isJ2 {
					switch fnName {
					case "Ext_java_util_Collections_Sort_1", "Ext_java_util_Collections_Sort_2":
						if len(x.Args) >= 1 && exprReferencesVar(unwrapCollectionAdaptersDown(x.Args[0], info), cand.v, info) {
							mutated = true
						}
					}
				}
				if sel, ok := stripParens(x.Fun).(*ast.SelectorExpr); ok {
					switch sel.Sel.Name {
					case "Add_1", "Add_2", "AddFirst", "AddLast", "AddAll_1", "Clear", "Remove_1_int32", "RemoveFirst", "RemoveLast", "Sort":
						if exprReferencesVar(sel.X, cand.v, info) {
							mutated = true
						}
					}
				}
			case *ast.AssignStmt:
				for _, lhs := range x.Lhs {
					if exprReferencesVar(lhs, cand.v, info) {
						mutated = true
					}
				}
			}
			return !mutated
		})
		if mutated {
			return true
		}
	}
	return false
}

// applySliceRewrites executes all verified rewrites innermost-first (expressions before statements, func decls last).
func applySliceRewrites(ctx *PassContext, eligible []*sliceCandidate, parentsByFile map[*ast.File]map[ast.Node]ast.Node) int {
	var exprOps, stmtOps, funcOps []sliceRewriteOp
	for _, cand := range eligible {
		for _, op := range cand.ops {
			switch op.kind {
			case opAccFuncDecl:
				funcOps = append(funcOps, op)
			case opAccCallRename, opExprAccParamArg, opExprGet, opExprSize, opExprIsEmpty, opExprRawUnbox,
				opExprImmListCopy, opExprImmCopyToSlice, opExprUnmodList, opExprUnmodListLive, opExprReadOnlyIter,
				opExprSliceIter, opExprToJavaList, opExprWrapTypedArrayList, opExprFormatJavaStr, opExprGetFirst, opExprGetLast, opExprRemoveAt:
				exprOps = append(exprOps, op)
			default:
				stmtOps = append(stmtOps, op)
			}
		}
	}
	spanOrder := func(a, b sliceRewriteOp) int {
		if c := cmp.Compare(a.spanNode.End()-a.spanNode.Pos(), b.spanNode.End()-b.spanNode.Pos()); c != 0 {
			return c
		}
		return cmp.Compare(a.spanNode.Pos(), b.spanNode.Pos())
	}
	slices.SortFunc(exprOps, spanOrder)
	slices.SortFunc(stmtOps, spanOrder)
	slices.SortFunc(funcOps, spanOrder)

	count := 0
	for _, op := range exprOps {
		if applySingleSliceRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}
	for _, op := range stmtOps {
		if applySingleSliceRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}
	for _, op := range funcOps {
		if applySingleSliceRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}
	return count
}

func applySingleSliceRewrite(ctx *PassContext, op sliceRewriteOp, parents map[ast.Node]ast.Node) bool {
	info := ctx.Info
	siteElemType := func() types.Type {
		return instantiateOwnerTypeAtExpr(op.cand.elemType, op.cand.ownerStruct, op.accessExpr, info)
	}
	makeSliceType := func() *ast.ArrayType {
		elt, _ := typeExprForFile(siteElemType(), op.file, ctx.Pkg, ctx.Importer)
		return &ast.ArrayType{Elt: elt}
	}
	makeElemType := func() ast.Expr {
		elt, _ := typeExprForFile(siteElemType(), op.file, ctx.Pkg, ctx.Importer)
		return elt
	}
	j2Sel := func(name string, pos token.Pos) *ast.SelectorExpr {
		q, _ := j2gortQualifierForFile(op.file, ctx)
		return &ast.SelectorExpr{X: &ast.Ident{NamePos: pos, Name: q}, Sel: &ast.Ident{NamePos: pos, Name: name}}
	}
	sliceValExpr := func(pos token.Pos) ast.Expr {
		e := clonePureExpr(op.accessExpr)
		if op.cand.isAccParam {
			return &ast.StarExpr{Star: pos, X: e}
		}
		return e
	}
	sliceAddrExpr := func(pos token.Pos) ast.Expr {
		e := clonePureExpr(op.accessExpr)
		if op.cand.isAccParam {
			return e
		}
		return &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: wrapPrimaryTarget(e)}
	}
	makeInitExpr := func(rhs ast.Expr, inHeader bool) ast.Expr {
		pos := rhs.Pos()
		st := makeSliceType()
		liveCallArg := op.otherExpr
		if c, ok := stripParens(rhs).(*ast.CallExpr); ok && len(c.Args) == 1 {
			liveCallArg = c.Args[0]
		}
		switch op.initKind {
		case sliceInitCap:
			origCapArg := op.otherExpr
			if origCapArg == nil {
				origCapArg = liveCallArg
			}
			if isNonNegativeConstInt(origCapArg, info) {
				return &ast.CallExpr{Fun: &ast.Ident{Name: "make"}, Args: []ast.Expr{st, &ast.BasicLit{Kind: token.INT, Value: "0"}, liveCallArg}}
			}
			return &ast.CallExpr{
				Fun:  &ast.IndexExpr{X: j2Sel("MakeNonNilSlice", pos), Index: makeElemType()},
				Args: []ast.Expr{liveCallArg},
			}
		case sliceInitCopySlice, sliceInitCopyTypedSlice:
			return &ast.CallExpr{
				Fun:  j2Sel("CloneNonNilSlice", pos),
				Args: []ast.Expr{clonePureExpr(op.otherExpr)},
			}
		case sliceInitCopyCol:
			elemT := siteElemType()
			switch {
			case types.Identical(elemT, types.Typ[types.String]):
				return &ast.CallExpr{Fun: j2Sel("CopyStrColToSlice", pos), Args: []ast.Expr{liveCallArg}}
			case isNamedUint32NodeType(elemT):
				return &ast.CallExpr{Fun: &ast.IndexExpr{X: j2Sel("CopyNodeColToSlice", pos), Index: makeElemType()}, Args: []ast.Expr{liveCallArg}}
			default:
				return &ast.CallExpr{Fun: &ast.IndexExpr{X: j2Sel("CopyCollectionToSlice", pos), Index: makeElemType()}, Args: []ast.Expr{liveCallArg}}
			}
		case sliceInitCopyIter:
			return &ast.CallExpr{
				Fun:  &ast.IndexExpr{X: j2Sel("CopyIteratorToSlice", pos), Index: makeElemType()},
				Args: []ast.Expr{liveCallArg},
			}
		case sliceInitMoveCand:
			return clonePureExpr(op.otherExpr)
		default:
			var lit ast.Expr = &ast.CompositeLit{Type: st}
			if inHeader {
				lit = &ast.ParenExpr{X: lit}
			}
			return lit
		}
	}

	switch op.kind {
	case opDeclField:
		fld := op.spanNode.(*ast.Field)
		fld.Type = makeSliceType()
		return true

	case opDeclVarSpec:
		spec := op.spanNode.(*ast.ValueSpec)
		spec.Type = makeSliceType()
		if len(spec.Values) == 1 && !isNilIdent(spec.Values[0], info) {
			spec.Values[0] = makeInitExpr(spec.Values[0], false)
		}
		return true

	case opDeclShortVar, opAssignInit:
		as := op.spanNode.(*ast.AssignStmt)
		as.Rhs[0] = makeInitExpr(as.Rhs[0], inControlHeader(as, parents))
		return true

	case opAssignNil:
		return false

	case opAssignReverseSelf:
		as := op.spanNode.(*ast.AssignStmt)
		pos := as.Pos()
		as.Rhs[0] = &ast.CallExpr{
			Fun:  j2Sel("SliceReverseInPlace", pos),
			Args: []ast.Expr{clonePureExpr(op.accessExpr)},
		}
		return true

	case opAccCallRename:
		op.elemIdent.Name = "sliceAcc_" + op.cand.accFn.Name()
		return true

	case opExprAccParamArg:
		parent := parents[op.outerExpr]
		var repl ast.Expr = clonePureExpr(op.accessExpr)
		if !op.cand.isAccParam {
			repl = &ast.UnaryExpr{OpPos: op.outerExpr.Pos(), Op: token.AND, X: wrapPrimaryTarget(repl)}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case opExprRawUnbox:
		parent := parents[op.outerExpr]
		if parent == nil {
			return false
		}
		if op.otherExpr != nil {
			return replaceExprInParent(parent, op.outerExpr, clonePureExpr(op.otherExpr), parents)
		}
		return replaceExprInParent(parent, op.outerExpr, &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name}, parents)

	case opExprGet:
		parent := parents[op.outerExpr]
		idxArg := op.otherExpr
		if op.call != nil && len(op.call.Args) > 0 {
			idxArg = op.call.Args[len(op.call.Args)-1]
		}
		idxExpr := &ast.IndexExpr{X: wrapPrimaryTarget(sliceValExpr(op.outerExpr.Pos())), Index: idxArg}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, idxExpr, parents)

	case opExprImmListCopy, opExprImmCopyToSlice, opExprUnmodList, opExprReadOnlyIter, opExprSliceIter, opExprToJavaList, opExprWrapTypedArrayList, opExprFormatJavaStr, opExprGetFirst, opExprGetLast:
		parent := parents[op.outerExpr]
		if parent == nil {
			return false
		}
		helper := "ImmListFromSlice"
		switch op.kind {
		case opExprImmCopyToSlice:
			helper = "ImmSliceCloneNonNull"
		case opExprUnmodList:
			helper = "SliceUnmodifiableList"
		case opExprReadOnlyIter:
			helper = "SliceReadOnlyIterable"
		case opExprSliceIter:
			helper = "SliceIterator"
			if j2gortHasExport(ctx, "NewTypedSliceIter") && !exprExpectsUnmodifiableIterator(op.outerExpr, parents, info) {
				helper = "NewTypedSliceIter"
			}
		case opExprToJavaList:
			helper = "SliceToJavaList"
		case opExprWrapTypedArrayList:
			helper = "WrapTypedArrayList"
		case opExprFormatJavaStr:
			helper = "SliceFormatJavaString"
		case opExprGetFirst:
			helper = "SliceGetFirst"
		case opExprGetLast:
			helper = "SliceGetLast"
		}
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{Fun: j2Sel(helper, pos), Args: []ast.Expr{sliceValExpr(pos)}}
		return replaceExprInParent(parent, op.outerExpr, repl, parents)

	case opExprUnmodListLive:
		parent := parents[op.outerExpr]
		if parent == nil {
			return false
		}
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{Fun: j2Sel("SliceUnmodifiableListLive", pos), Args: []ast.Expr{sliceAddrExpr(pos)}}
		return replaceExprInParent(parent, op.outerExpr, repl, parents)

	case opExprRemoveAt:
		parent := parents[op.outerExpr]
		if parent == nil {
			return false
		}
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{Fun: j2Sel("SliceRemoveAtElem", pos), Args: []ast.Expr{sliceAddrExpr(pos), op.call.Args[0]}}
		return replaceExprInParent(parent, op.outerExpr, repl, parents)

	case opExprSize:
		parent := parents[op.call]
		pos := op.call.Pos()
		lenCall := &ast.CallExpr{Fun: &ast.Ident{NamePos: pos, Name: "len"}, Args: []ast.Expr{sliceValExpr(pos)}}
		conv := &ast.CallExpr{Fun: &ast.Ident{NamePos: pos, Name: "int32"}, Args: []ast.Expr{lenCall}}
		return parent != nil && replaceExprInParent(parent, op.call, conv, parents)

	case opExprIsEmpty:
		parent := parents[op.call]
		if parent == nil {
			return false
		}
		pos := op.call.Pos()
		targetNode, targetParent, tok := ast.Expr(op.call), parent, token.EQL
		cur := ast.Expr(op.call)
		for {
			pe, isParen := parents[cur].(*ast.ParenExpr)
			if !isParen || pe.X != cur {
				break
			}
			cur = pe
		}
		if u, isNot := parents[cur].(*ast.UnaryExpr); isNot && u.Op == token.NOT && u.X == cur && parents[u] != nil {
			if _, isES := parents[u].(*ast.ExprStmt); !isES {
				targetNode, targetParent, tok = u, parents[u], token.NEQ
			}
		}
		var cmpExpr ast.Expr = &ast.BinaryExpr{
			X:     &ast.CallExpr{Fun: &ast.Ident{NamePos: pos, Name: "len"}, Args: []ast.Expr{sliceValExpr(pos)}},
			OpPos: pos, Op: tok,
			Y: &ast.BasicLit{ValuePos: pos, Kind: token.INT, Value: "0"},
		}
		if binaryNeedsParens(targetParent, tok) {
			cmpExpr = &ast.ParenExpr{X: cmpExpr}
		}
		return replaceExprInParent(targetParent, targetNode, cmpExpr, parents)

	case opStmtAdd:
		unboxedArg := wrapSliceIfaceWriteElem(unwrapBoxArgExpr(op.call.Args[0], op.unwrapsBox), op.cand.elemType, makeElemType, j2Sel, ctx.Info)
		pos := op.call.Pos()
		as := &ast.AssignStmt{
			Lhs: []ast.Expr{sliceValExpr(pos)}, Tok: token.ASSIGN,
			Rhs: []ast.Expr{&ast.CallExpr{Fun: &ast.Ident{NamePos: pos, Name: "append"}, Args: []ast.Expr{sliceValExpr(pos), unboxedArg}}},
		}
		return replaceStmtInParent(parents[op.stmt], op.stmt, []ast.Stmt{as}, parents)

	case opStmtInsertAt:
		elemRaw := op.call.Args[len(op.call.Args)-1]
		unboxedArg := wrapSliceIfaceWriteElem(unwrapBoxArgExpr(elemRaw, op.unwrapsBox), op.cand.elemType, makeElemType, j2Sel, ctx.Info)
		pos := op.call.Pos()
		insCall := &ast.CallExpr{
			Fun:  j2Sel("SliceInsertAt", pos),
			Args: []ast.Expr{sliceValExpr(pos), op.otherExpr, unboxedArg},
		}
		as := &ast.AssignStmt{Lhs: []ast.Expr{sliceValExpr(pos)}, Tok: token.ASSIGN, Rhs: []ast.Expr{insCall}}
		return replaceStmtInParent(parents[op.stmt], op.stmt, []ast.Stmt{as}, parents)

	case opStmtSortNat:
		pos := op.call.Pos()
		sortCall := &ast.ExprStmt{X: &ast.CallExpr{Fun: j2Sel("SliceSortNatural", pos), Args: []ast.Expr{sliceValExpr(pos)}}}
		return replaceStmtInParent(parents[op.stmt], op.stmt, []ast.Stmt{sortCall}, parents)

	case opStmtSortCmp:
		pos := op.call.Pos()
		sortCall := &ast.ExprStmt{X: &ast.CallExpr{Fun: j2Sel("SliceSortComparator", pos), Args: []ast.Expr{sliceValExpr(pos), op.otherExpr}}}
		return replaceStmtInParent(parents[op.stmt], op.stmt, []ast.Stmt{sortCall}, parents)

	case opStmtSet:
		pos := op.call.Pos()
		unboxedArg := wrapSliceIfaceWriteElem(unwrapBoxArgExpr(op.call.Args[1], op.unwrapsBox), op.cand.elemType, makeElemType, j2Sel, ctx.Info)
		lhs := &ast.IndexExpr{X: wrapPrimaryTarget(sliceValExpr(pos)), Index: op.call.Args[0]}
		as := &ast.AssignStmt{Lhs: []ast.Expr{lhs}, Tok: token.ASSIGN, Rhs: []ast.Expr{unboxedArg}}
		return replaceStmtInParent(parents[op.stmt], op.stmt, []ast.Stmt{as}, parents)

	case opStmtRemoveAt:
		pos := op.call.Pos()
		rmCall := &ast.CallExpr{Fun: j2Sel("SliceRemoveAt", pos), Args: []ast.Expr{sliceValExpr(pos), op.call.Args[0]}}
		as := &ast.AssignStmt{Lhs: []ast.Expr{sliceValExpr(pos)}, Tok: token.ASSIGN, Rhs: []ast.Expr{rmCall}}
		return replaceStmtInParent(parents[op.stmt], op.stmt, []ast.Stmt{as}, parents)

	case opStmtRemoveFirst:
		pos := op.call.Pos()
		chk := &ast.AssignStmt{
			Lhs: []ast.Expr{&ast.Ident{NamePos: pos, Name: "_"}},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{&ast.CallExpr{Fun: j2Sel("SliceGetFirst", pos), Args: []ast.Expr{sliceValExpr(pos)}}},
		}
		rmCall := &ast.CallExpr{
			Fun:  j2Sel("SliceRemoveAt", pos),
			Args: []ast.Expr{sliceValExpr(pos), &ast.BasicLit{ValuePos: pos, Kind: token.INT, Value: "0"}},
		}
		as := &ast.AssignStmt{Lhs: []ast.Expr{sliceValExpr(pos)}, Tok: token.ASSIGN, Rhs: []ast.Expr{rmCall}}
		return replaceStmtInParent(parents[op.stmt], op.stmt, []ast.Stmt{chk, as}, parents)

	case opStmtRemoveLast:
		pos := op.call.Pos()
		chk := &ast.AssignStmt{
			Lhs: []ast.Expr{&ast.Ident{NamePos: pos, Name: "_"}},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{&ast.CallExpr{Fun: j2Sel("SliceGetLast", pos), Args: []ast.Expr{sliceValExpr(pos)}}},
		}
		makeLastIdx := func() ast.Expr {
			return &ast.BinaryExpr{
				X:     &ast.CallExpr{Fun: &ast.Ident{NamePos: pos, Name: "len"}, Args: []ast.Expr{sliceValExpr(pos)}},
				OpPos: pos, Op: token.SUB,
				Y: &ast.BasicLit{ValuePos: pos, Kind: token.INT, Value: "1"},
			}
		}
		reslice := &ast.AssignStmt{
			Lhs: []ast.Expr{sliceValExpr(pos)},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{&ast.SliceExpr{X: wrapPrimaryTarget(sliceValExpr(pos)), High: makeLastIdx()}},
		}
		stmts := []ast.Stmt{chk, reslice}
		if typeNeedsGCZeroing(op.cand.elemType) {
			clearStmt := &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun:  &ast.Ident{NamePos: pos, Name: "clear"},
					Args: []ast.Expr{&ast.SliceExpr{X: wrapPrimaryTarget(sliceValExpr(pos)), Low: makeLastIdx()}},
				},
			}
			stmts = []ast.Stmt{chk, clearStmt, reslice}
		}
		return replaceStmtInParent(parents[op.stmt], op.stmt, stmts, parents)

	case opStmtAddAllSlice:
		pos := op.call.Pos()
		as := &ast.AssignStmt{
			Lhs: []ast.Expr{clonePureExpr(op.accessExpr)}, Tok: token.ASSIGN,
			Rhs: []ast.Expr{&ast.CallExpr{
				Fun:      &ast.Ident{NamePos: pos, Name: "append"},
				Args:     []ast.Expr{clonePureExpr(op.accessExpr), clonePureExpr(op.otherExpr)},
				Ellipsis: pos,
			}},
		}
		return replaceStmtInParent(parents[op.stmt], op.stmt, []ast.Stmt{as}, parents)

	case opStmtAddAllCol:
		pos := op.call.Pos()
		var fn ast.Expr
		switch {
		case types.Identical(op.cand.elemType, types.Typ[types.String]):
			fn = j2Sel("AppendStrCol", pos)
		case isNamedUint32NodeType(op.cand.elemType):
			fn = j2Sel("AppendNodeCol", pos)
		default:
			fn = &ast.IndexExpr{X: j2Sel("AppendCollection", pos), Index: makeElemType()}
		}
		call := &ast.CallExpr{Fun: fn, Args: []ast.Expr{clonePureExpr(op.accessExpr), op.call.Args[0]}}
		as := &ast.AssignStmt{Lhs: []ast.Expr{clonePureExpr(op.accessExpr)}, Tok: token.ASSIGN, Rhs: []ast.Expr{call}}
		return replaceStmtInParent(parents[op.stmt], op.stmt, []ast.Stmt{as}, parents)

	case opStmtClear:
		pos := op.call.Pos()
		reslice := &ast.AssignStmt{
			Lhs: []ast.Expr{sliceValExpr(pos)}, Tok: token.ASSIGN,
			Rhs: []ast.Expr{&ast.SliceExpr{X: wrapPrimaryTarget(sliceValExpr(pos)), High: &ast.BasicLit{ValuePos: pos, Kind: token.INT, Value: "0"}}},
		}
		stmts := []ast.Stmt{reslice}
		if typeNeedsGCZeroing(op.cand.elemType) {
			clearStmt := &ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{NamePos: pos, Name: "clear"}, Args: []ast.Expr{sliceValExpr(pos)}}}
			stmts = []ast.Stmt{clearStmt, reslice}
		}
		return replaceStmtInParent(parents[op.stmt], op.stmt, stmts, parents)

	case opStmtIterSlice:
		rs := op.stmt.(*ast.RangeStmt)
		rs.X = sliceValExpr(rs.Pos())
		parents[rs.X] = rs
		if op.dropStmts > 0 && op.elemIdent != nil && len(rs.Body.List) >= op.dropStmts {
			rs.Body.Lbrace = rs.Body.List[op.dropStmts-1].End()
			rs.Value = &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name}
			rs.Body.List = rs.Body.List[op.dropStmts:]
		}
		return true

	case opStmtIterCursor:
		idx := slices.Index(op.block.List, op.stmt)
		if idx < 0 || idx+1 >= len(op.block.List) || len(op.forStmt.Body.List) < op.dropStmts {
			return false
		}
		bodyLbrace := op.forStmt.Body.List[op.dropStmts-1].End()
		var loopStmt ast.Stmt
		if op.isReverseIter {
			idxVarName := freshLoopIndexName(op.forStmt.Body, "_revIdx")
			initIdx := &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: idxVarName}},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{&ast.BinaryExpr{
					X:  &ast.CallExpr{Fun: &ast.Ident{Name: "len"}, Args: []ast.Expr{sliceValExpr(op.stmt.Pos())}},
					Op: token.SUB,
					Y:  &ast.BasicLit{Kind: token.INT, Value: "1"},
				}},
			}
			condIdx := &ast.BinaryExpr{
				X:  &ast.Ident{Name: idxVarName},
				Op: token.GEQ,
				Y:  &ast.BasicLit{Kind: token.INT, Value: "0"},
			}
			postIdx := &ast.IncDecStmt{X: &ast.Ident{Name: idxVarName}, Tok: token.DEC}
			elemDef := &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name}},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{&ast.IndexExpr{X: wrapPrimaryTarget(sliceValExpr(op.elemIdent.Pos())), Index: &ast.Ident{Name: idxVarName}}},
			}
			bodyStmts := append([]ast.Stmt{elemDef}, op.forStmt.Body.List[op.dropStmts:]...)
			loopStmt = &ast.ForStmt{
				For:  op.stmt.Pos(),
				Init: initIdx,
				Cond: condIdx,
				Post: postIdx,
				Body: &ast.BlockStmt{Lbrace: bodyLbrace, List: bodyStmts, Rbrace: op.forStmt.Body.Rbrace},
			}
		} else {
			loopStmt = &ast.RangeStmt{
				For:    op.stmt.Pos(),
				Key:    &ast.Ident{NamePos: op.elemIdent.Pos(), Name: "_"},
				Value:  &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name},
				TokPos: op.elemIdent.End(),
				Tok:    token.DEFINE,
				X:      sliceValExpr(op.stmt.Pos()),
				Body:   &ast.BlockStmt{Lbrace: bodyLbrace, List: op.forStmt.Body.List[op.dropStmts:], Rbrace: op.forStmt.Body.Rbrace},
			}
		}
		var replacement ast.Stmt = loopStmt
		if op.label != nil {
			replacement = &ast.LabeledStmt{Label: op.label, Colon: op.label.End(), Stmt: loopStmt}
		}
		op.block.List = append(append(append([]ast.Stmt(nil), op.block.List[:idx]...), replacement), op.block.List[idx+2:]...)
		parents[replacement] = op.block
		if len(op.block.List) == 1 {
			if outerBlock, ok := parents[op.block].(*ast.BlockStmt); ok {
				if op.label == nil {
					if rs, isRS := loopStmt.(*ast.RangeStmt); isRS {
						rs.For = op.block.Lbrace
					} else if fs, isFS := loopStmt.(*ast.ForStmt); isFS {
						fs.For = op.block.Lbrace
					}
				}
				replaceStmtInParent(outerBlock, op.block, op.block.List, parents)
			}
		}
		return true

	case opAccFuncDecl:
		fn := op.cand.accFnDecl
		ptrSliceType := &ast.StarExpr{X: makeSliceType()}
		if !op.cand.isExportedAccFn {
			op.cand.accParamField.Type = ptrSliceType
			return true
		}
		origFuncName := fn.Name.Name
		twinName := "sliceAcc_" + origFuncName
		wrapperParams, callArgs := cloneFuncParamsAndCallArgs(fn.Type.Params)
		// Mutate fn in-place into twinName so all inline comments inside fn.Body keep their exact original Pos,
		// and append the brand-new NoPos wrapperDecl at the end of op.file.Decls after all file comments.
		fieldIdx := slices.Index(fn.Type.Params.List, op.cand.accParamField)
		if fieldIdx < 0 || op.cand.accParamIdx >= len(callArgs) {
			return false
		}
		origParamName := op.cand.accParamField.Names[0].Name
		callArgs[op.cand.accParamIdx] = &ast.UnaryExpr{Op: token.AND, X: &ast.Ident{Name: "_sliceAcc_"}}

		var wrapperResults *ast.FieldList
		hasRet := fn.Type.Results != nil && len(fn.Type.Results.List) == 1
		if hasRet {
			wrapperResults = &ast.FieldList{List: []*ast.Field{{Type: stripTypePos(fn.Type.Results.List[0].Type)}}}
		}
		fn.Name.Name = twinName
		fn.Type.Params.List[fieldIdx].Type = ptrSliceType

		twinCall := &ast.CallExpr{Fun: &ast.Ident{Name: twinName}, Args: callArgs}
		accDecl := &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok:   token.VAR,
				Specs: []ast.Spec{&ast.ValueSpec{Names: []*ast.Ident{{Name: "_sliceAcc_"}}, Type: makeSliceType()}},
			},
		}
		elemBoxExpr := boxSliceElemForJavaList(&ast.Ident{Name: "_sliceElem_"}, siteElemType(), j2Sel)
		copyLoop := &ast.RangeStmt{
			Key:   &ast.Ident{Name: "_"},
			Value: &ast.Ident{Name: "_sliceElem_"},
			Tok:   token.DEFINE,
			X:     &ast.Ident{Name: "_sliceAcc_"},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.CallExpr{Fun: j2Sel("CheckNotNull", token.NoPos), Args: []ast.Expr{&ast.Ident{Name: origParamName}}},
								Sel: &ast.Ident{Name: "Add_1"},
							},
							Args: []ast.Expr{elemBoxExpr},
						},
					},
				},
			},
		}
		var wrapperStmts []ast.Stmt
		wrapperStmts = append(wrapperStmts, accDecl)
		if hasRet {
			wrapperStmts = append(wrapperStmts,
				&ast.AssignStmt{Lhs: []ast.Expr{&ast.Ident{Name: "_sliceRet_"}}, Tok: token.DEFINE, Rhs: []ast.Expr{twinCall}},
				copyLoop,
				&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "_sliceRet_"}}},
			)
		} else {
			wrapperStmts = append(wrapperStmts, &ast.ExprStmt{X: twinCall}, copyLoop)
		}
		wrapperDecl := &ast.FuncDecl{
			Name: &ast.Ident{Name: origFuncName},
			Type: &ast.FuncType{
				Params:  &ast.FieldList{List: wrapperParams},
				Results: wrapperResults,
			},
			Body: &ast.BlockStmt{List: wrapperStmts},
		}
		op.file.Decls = append(op.file.Decls, wrapperDecl)
		return true

	case opStmtForEach:
		spec := op.forEachSpec
		if spec == nil || len(spec.fnLit.Body.List) < spec.dropStmts {
			return false
		}
		pos := spec.exprStmt.Pos()
		rangeLoop := &ast.RangeStmt{
			For:    pos,
			Key:    &ast.Ident{NamePos: spec.paramPos, Name: "_"},
			Value:  &ast.Ident{NamePos: spec.paramPos, Name: spec.paramName},
			TokPos: spec.paramPos,
			Tok:    token.DEFINE,
			X:      sliceValExpr(pos),
			Body:   &ast.BlockStmt{List: spec.fnLit.Body.List[spec.dropStmts:]},
		}
		if spec.boundAs != nil {
			return replaceStmtInParent(parents[spec.exprStmt], spec.exprStmt, []ast.Stmt{&ast.BlockStmt{List: []ast.Stmt{spec.boundAs, rangeLoop}}}, parents)
		}
		return replaceStmtInParent(parents[spec.exprStmt], spec.exprStmt, []ast.Stmt{rangeLoop}, parents)
	}
	return false
}

func unwrapBoxArgExpr(rawArg ast.Expr, unwrapsBox bool) ast.Expr {
	arg := stripParens(rawArg)
	if unwrapsBox {
		if boxCall, ok := arg.(*ast.CallExpr); ok && len(boxCall.Args) == 1 {
			return stripParens(boxCall.Args[0])
		}
	}
	return arg
}

func wrapSliceIfaceWriteElem(
	elemExpr ast.Expr,
	sliceElemT types.Type,
	makeElemType func() ast.Expr,
	j2Sel func(string, token.Pos) *ast.SelectorExpr,
	info *types.Info,
) ast.Expr {
	if elemExpr == nil || sliceElemT == nil || !isNonEmptyInterfaceType(sliceElemT) {
		return elemExpr
	}
	if tv, ok := info.Types[elemExpr]; ok && tv.Type != nil && isPtrToStructType(tv.Type) && !isGuaranteedNonNilAlloc(elemExpr, info) {
		pos := elemExpr.Pos()
		return &ast.CallExpr{
			Fun:  &ast.IndexExpr{X: j2Sel("AsIface", pos), Index: makeElemType()},
			Args: []ast.Expr{elemExpr},
		}
	}
	return elemExpr
}

func resolveCandidateAccess(expr ast.Expr, candidates map[*types.Var]*sliceCandidate, info *types.Info) *sliceCandidate {
	switch x := stripParens(expr).(type) {
	case *ast.Ident:
		if c := candidates[varFromObj(info.Uses[x])]; c != nil && !c.isField {
			return c
		}
	case *ast.SelectorExpr:
		if c := candidates[varFromObj(info.Uses[x.Sel])]; c != nil && c.isField && isPureFieldSelector(x, info) {
			return c
		}
	}
	return nil
}

// matchFieldMoveInitToLocal matches `var dst = this_.field` (or `dst := this_.field`) when
// `this_.field = nil` immediately follows in the same BlockStmt with only pure, call-free
// non-interfering assignments/declarations in between (e.g. RemoveUnusedCode_VarInfo.SetIsExplicitlyNotRemovable).
func matchFieldMoveInitToLocal(
	initExpr ast.Expr,
	span ast.Node,
	dstCand *sliceCandidate,
	candidates map[*types.Var]*sliceCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) (ast.Expr, *sliceCandidate, bool) {
	if initExpr == nil || dstCand == nil || dstCand.isField || parents == nil {
		return nil, nil, false
	}
	srcExpr := stripParens(initExpr)
	srcCand := resolveCandidateAccess(srcExpr, candidates, info)
	if srcCand == nil || !srcCand.isField || srcCand == dstCand {
		return nil, nil, false
	}
	var declStmt ast.Stmt
	switch s := span.(type) {
	case *ast.AssignStmt:
		declStmt = s
	case *ast.ValueSpec:
		if gd, ok := parents[s].(*ast.GenDecl); ok {
			declStmt, _ = parents[gd].(*ast.DeclStmt)
		}
	}
	if declStmt == nil {
		return nil, nil, false
	}
	blk, ok := parents[declStmt].(*ast.BlockStmt)
	if !ok || blk == nil {
		return nil, nil, false
	}
	declIdx := slices.Index(blk.List, declStmt)
	if declIdx < 0 {
		return nil, nil, false
	}
	// Collect receiver chain variables of srcExpr so no intervening statement can mutate them.
	srcBaseVars := make(map[*types.Var]bool)
	ast.Inspect(srcExpr, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			if v := varFromObj(info.Uses[id]); v != nil && v != srcCand.v {
				srcBaseVars[v] = true
			}
		}
		return true
	})
	exprHasCallOrDisallowedVar := func(e ast.Expr) bool {
		bad := false
		ast.Inspect(e, func(n ast.Node) bool {
			if bad || n == nil {
				return false
			}
			switch x := n.(type) {
			case *ast.CallExpr, *ast.FuncLit:
				bad = true
				return false
			case *ast.Ident:
				v := varFromObj(info.Uses[x])
				if v == dstCand.v || v == srcCand.v {
					bad = true
					return false
				}
			}
			return true
		})
		return bad
	}
	for k := declIdx + 1; k < len(blk.List); k++ {
		st := blk.List[k]
		if as, ok := st.(*ast.AssignStmt); ok && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 {
			if sameAccessPath(as.Lhs[0], srcExpr, srcCand.v, info) && isNilIdent(as.Rhs[0], info) {
				return srcExpr, srcCand, true
			}
			if lhsId, okId := stripParens(as.Lhs[0]).(*ast.Ident); okId && lhsId.Name == "_" {
				if _, okR := stripParens(as.Rhs[0]).(*ast.Ident); okR {
					continue
				}
			}
			if exprHasCallOrDisallowedVar(as.Lhs[0]) || exprHasCallOrDisallowedVar(as.Rhs[0]) {
				return nil, nil, false
			}
			if lhsId, okId := stripParens(as.Lhs[0]).(*ast.Ident); okId && srcBaseVars[varFromObj(info.Uses[lhsId])] {
				return nil, nil, false
			}
			continue
		}
		if ds, ok := st.(*ast.DeclStmt); ok {
			gd, okG := ds.Decl.(*ast.GenDecl)
			if !okG || gd.Tok != token.VAR || len(gd.Specs) != 1 {
				return nil, nil, false
			}
			vs, okV := gd.Specs[0].(*ast.ValueSpec)
			if !okV || len(vs.Names) != 1 || len(vs.Values) > 1 {
				return nil, nil, false
			}
			if len(vs.Values) == 1 && exprHasCallOrDisallowedVar(vs.Values[0]) {
				return nil, nil, false
			}
			continue
		}
		return nil, nil, false
	}
	return nil, nil, false
}

func isValidSliceHandoff(
	as *ast.AssignStmt,
	dst, src *sliceCandidate,
	unwrappedRhs ast.Expr,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varMaxUsePos map[*types.Var]token.Pos,
	builderChecks *[]builderHandoffSpec,
) bool {
	if dst.isAccParam || src.isAccParam {
		return false
	}
	// Case 1: Local variable moved once (`this_.rest = list` or `dst = localSrc`) as its final use in the same block.
	if !src.isField {
		assignBlock, _ := parents[as].(*ast.BlockStmt)
		declBlock := localDeclBlock(src, parents)
		return assignBlock != nil && assignBlock == declBlock && varMaxUsePos[src.v] <= as.End()
	}
	// Case 2: Builder-to-product field handoff `this_.f = builder.g` (or `this_.f = UnmodifiableList(builder.g)`) inside `Init<DstStruct>*`.
	if dst.isField && src.isField && dst.ownerStruct != src.ownerStruct {
		fn := enclosingFuncObj(as, parents, info)
		if fn != nil && strings.HasPrefix(fn.Name(), "Init") {
			rhsSel, _ := stripParens(unwrappedRhs).(*ast.SelectorExpr)
			lhsSel, _ := stripParens(as.Lhs[0]).(*ast.SelectorExpr)
			if rhsSel != nil && lhsSel != nil {
				bId, okB := stripParens(rhsSel.X).(*ast.Ident)
				tId, okT := stripParens(lhsSel.X).(*ast.Ident)
				if okB && okT && info.Uses[bId] != info.Uses[tId] {
					*builderChecks = append(*builderChecks, builderHandoffSpec{
						dst: dst, src: src, builderType: src.ownerStruct, initFunc: fn,
					})
					return true
				}
			}
		}
	}
	return false
}

func verifyBuilderSingleUseInPkg(ctx *PassContext, parentsByFile map[*ast.File]map[ast.Node]ast.Node, spec builderHandoffSpec) bool {
	info := ctx.Info
	// Ensure src field is only mutated on `this_` inside methods of builderType, and find build methods.
	buildMethods := make(map[*types.Func]bool)
	type builderMethodInfo struct {
		fnOrig  *types.Func
		recvVar *types.Var
		body    *ast.BlockStmt
	}
	var builderMethodList []builderMethodInfo
	for _, file := range ctx.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil || fn.Recv == nil || len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 {
				continue
			}
			recvVar := varFromObj(info.Defs[fn.Recv.List[0].Names[0]])
			if recvVar == nil || ptrNamedOrigin(recvVar.Type()) != spec.builderType {
				continue
			}
			fnObj, _ := info.Defs[fn.Name].(*types.Func)
			fnOrig := funcOrigin(fnObj)
			if fnOrig != nil {
				builderMethodList = append(builderMethodList, builderMethodInfo{fnOrig: fnOrig, recvVar: recvVar, body: fn.Body})
			}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if call, ok := n.(*ast.CallExpr); ok {
					for _, a := range call.Args {
						if id, ok := stripParens(a).(*ast.Ident); ok && varFromObj(info.Uses[id]) == recvVar && fnOrig != nil {
							buildMethods[fnOrig] = true
						}
					}
				}
				return true
			})
		}
	}
	// Propagate buildMethods transitively across builder helper/overload methods (`this_.Build_3(...)`).
	for changed := true; changed; {
		changed = false
		for _, bm := range builderMethodList {
			if buildMethods[bm.fnOrig] {
				continue
			}
			ast.Inspect(bm.body, func(n ast.Node) bool {
				if call, ok := n.(*ast.CallExpr); ok {
					if sel, ok := stripParens(call.Fun).(*ast.SelectorExpr); ok {
						if id, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == bm.recvVar {
							if callee, _ := calleeFunc(call, info); callee != nil && buildMethods[funcOrigin(callee)] {
								buildMethods[bm.fnOrig] = true
								changed = true
								return false
							}
						}
					}
				}
				return true
			})
		}
	}

	// Verify every local `*builderType` variable is never used on any control-flow path after calling a build method.
	maxUse := make(map[*types.Var]token.Pos)
	defIdents := make(map[*types.Var]*ast.Ident)
	// order-independent: commutative max token.Pos tracking per variable
	for id, obj := range info.Uses {
		if v := varFromObj(obj); v != nil {
			if ptrNamedOrigin(v.Type()) == spec.builderType {
				if id.Pos() > maxUse[v] {
					maxUse[v] = id.Pos()
				}
			}
		}
	}
	// order-independent: maps each unique *types.Var to its single defining *ast.Ident
	for id, obj := range info.Defs {
		if v := varFromObj(obj); v != nil && ptrNamedOrigin(v.Type()) == spec.builderType {
			defIdents[v] = id
		}
	}

	valid := true
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		ast.Inspect(file, func(n ast.Node) bool {
			if !valid {
				return false
			}
			// Check src field uses: outside initFunc, must be `recv.srcField` inside a method of builderType.
			if id, ok := n.(*ast.Ident); ok && varFromObj(info.Uses[id]) == spec.src.v {
				encl := enclosingFuncObj(id, parents, info)
				if funcOrigin(encl) != funcOrigin(spec.initFunc) {
					if encl == nil {
						valid = false
						return false
					}
					sig, _ := encl.Type().(*types.Signature)
					isBuilderMethod := false
					if sig != nil && sig.Recv() != nil {
						if ptrNamedOrigin(sig.Recv().Type()) == spec.builderType {
							isBuilderMethod = true
						}
					} else if strings.HasPrefix(encl.Name(), "Init"+spec.builderType.Obj().Name()) || strings.HasPrefix(encl.Name(), "New"+spec.builderType.Obj().Name()) {
						isBuilderMethod = true
					}
					if !isBuilderMethod {
						valid = false
						return false
					}
				}
			}
			// Check calls to buildMethods: if receiver is a local var `b`, `b` cannot be used after this call or in a loop.
			if call, ok := n.(*ast.CallExpr); ok {
				callee, _ := calleeFunc(call, info)
				if buildMethods[funcOrigin(callee)] {
					if sel, ok := stripParens(call.Fun).(*ast.SelectorExpr); ok {
						if bId, ok := stripParens(sel.X).(*ast.Ident); ok {
							if bVar := varFromObj(info.Uses[bId]); bVar != nil {
								encl := enclosingFuncObj(call, parents, info)
								isSelfBuildRedirect := false
								if encl != nil {
									if sig, _ := encl.Type().(*types.Signature); sig != nil && sig.Recv() != nil && varFromObj(sig.Recv()) == bVar {
										isSelfBuildRedirect = true
									}
								}
								if !isSelfBuildRedirect && builderVarUsedAfterBuildCall(bVar, call, defIdents[bVar], maxUse[bVar], parents, info) {
									valid = false
									return false
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return valid
}

func builderVarUsedAfterBuildCall(
	bVar *types.Var,
	buildCall *ast.CallExpr,
	defId *ast.Ident,
	maxPos token.Pos,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	if maxPos <= buildCall.End() {
		// Still ensure buildCall is not inside a loop relative to bVar's declaration block!
		var declBlk *ast.BlockStmt
		for p := ast.Node(defId); p != nil; p = parents[p] {
			if b, ok := p.(*ast.BlockStmt); ok {
				declBlk = b
				break
			}
		}
		if declBlk == nil {
			return false
		}
		for cur := ast.Node(buildCall); cur != nil && cur != declBlk; cur = parents[cur] {
			switch cur.(type) {
			case *ast.ForStmt, *ast.RangeStmt, *ast.FuncLit:
				return true
			}
		}
		return false
	}
	if defId == nil {
		return true
	}
	var declBlk *ast.BlockStmt
	for p := ast.Node(defId); p != nil; p = parents[p] {
		if b, ok := p.(*ast.BlockStmt); ok {
			declBlk = b
			break
		}
	}
	if declBlk == nil {
		return true
	}
	nodeUsesVarAfter := func(root ast.Node, afterPos token.Pos) bool {
		used := false
		ast.Inspect(root, func(n ast.Node) bool {
			if used || n == nil {
				return false
			}
			if id, ok := n.(*ast.Ident); ok && varFromObj(info.Uses[id]) == bVar && (!afterPos.IsValid() || id.Pos() > afterPos) {
				used = true
				return false
			}
			return true
		})
		return used
	}
	checkedInnerStmt := false
	for cur := ast.Node(buildCall); cur != nil && cur != declBlk; cur = parents[cur] {
		switch cur.(type) {
		case *ast.ForStmt, *ast.RangeStmt, *ast.FuncLit:
			return true
		}
		st, isSt := cur.(ast.Stmt)
		if !isSt {
			continue
		}
		if !checkedInnerStmt {
			checkedInnerStmt = true
			if _, isIf := st.(*ast.IfStmt); !isIf {
				if _, isSw := st.(*ast.SwitchStmt); !isSw {
					if nodeUsesVarAfter(st, buildCall.End()) {
						return true
					}
				}
			}
		}
		var list []ast.Stmt
		switch par := parents[st].(type) {
		case *ast.BlockStmt:
			list = par.List
		case *ast.CaseClause:
			list = par.Body
		case *ast.CommClause:
			list = par.Body
		}
		if len(list) > 0 {
			idx := slices.Index(list, st)
			if idx >= 0 {
				for _, later := range list[idx+1:] {
					if nodeUsesVarAfter(later, token.NoPos) {
						return true
					}
				}
			}
		}
	}
	return !checkedInnerStmt
}

func verifyLocalUnmodNoPostMutation(cand *sliceCandidate, parentsByFile map[*ast.File]map[ast.Node]ast.Node) bool {
	if len(cand.unmodChecks) == 0 {
		return true
	}
	var maxMutPos token.Pos
	for _, mn := range cand.mutNodes {
		if mn.Pos() > maxMutPos {
			maxMutPos = mn.Pos()
		}
	}
	for _, uc := range cand.unmodChecks {
		if maxMutPos > uc.usePos {
			return false
		}
		// Also ensure unmodifiableList(cand) is not inside a loop body that also mutates cand!
		parents := parentsByFile[uc.file]
		declBlk := localDeclBlock(cand, parents)
		for cur := uc.callNode; cur != nil && cur != declBlk; cur = parents[cur] {
			switch loop := cur.(type) {
			case *ast.ForStmt:
				for _, mn := range cand.mutNodes {
					if mn.Pos() >= loop.Pos() && mn.End() <= loop.End() {
						return false
					}
				}
			case *ast.RangeStmt:
				for _, mn := range cand.mutNodes {
					if mn.Pos() >= loop.Pos() && mn.End() <= loop.End() {
						return false
					}
				}
			}
		}
	}
	return true
}

func localDeclBlock(cand *sliceCandidate, parents map[ast.Node]ast.Node) *ast.BlockStmt {
	if len(cand.ops) == 0 {
		return nil
	}
	span := cand.ops[0].spanNode
	if spec, ok := span.(*ast.ValueSpec); ok {
		if gd, ok := parents[spec].(*ast.GenDecl); ok {
			if ds, ok := parents[gd].(*ast.DeclStmt); ok {
				b, _ := parents[ds].(*ast.BlockStmt)
				return b
			}
		}
	}
	b, _ := parents[span].(*ast.BlockStmt)
	return b
}

func isInInitFieldsMethod(node ast.Node, owner *types.Named, parents map[ast.Node]ast.Node, info *types.Info) bool {
	fn := enclosingFuncObj(node, parents, info)
	return fn != nil && (fn.Name() == "initFields_"+owner.Obj().Name() || fn.Name() == "initFields")
}

func isGuardedNonNilSliceUse(useNode ast.Node, accessExpr ast.Expr, cand *sliceCandidate, parents map[ast.Node]ast.Node, hasCopyHelpers bool, ctx *PassContext, info *types.Info) bool {
	return isGuardedNonNilVarUse(useNode, accessExpr, cand.v, parents, info, func(rhs ast.Expr) bool {
		_, _, _, _, ok := classifySliceInitExpr(rhs, cand, nil, hasCopyHelpers, ctx, info)
		return ok
	})
}

func isGuardedNonNilVarUse(
	useNode ast.Node,
	accessExpr ast.Expr,
	v *types.Var,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	isNonNilInitRhs func(ast.Expr) bool,
) bool {
	// 1. Check enclosing short-circuit `&&` / `||` expressions and `if` statements.
	for cur := useNode; cur != nil; cur = parents[cur] {
		if _, isFuncLit := cur.(*ast.FuncLit); isFuncLit {
			break
		}
		if bin, ok := parents[cur].(*ast.BinaryExpr); ok && cur == bin.Y {
			if bin.Op == token.LAND && condGuaranteesNonNil(bin.X, accessExpr, v, info) {
				return true
			}
			if bin.Op == token.LOR && condGuaranteesNil(bin.X, accessExpr, v, info) {
				return true
			}
		}
		if ifStmt, ok := parents[cur].(*ast.IfStmt); ok {
			if cur == ifStmt.Body && condGuaranteesNonNil(ifStmt.Cond, accessExpr, v, info) && !stmtAssignsNilBetween(ifStmt.Body.List, token.NoPos, useNode.Pos(), v, info) {
				return true
			}
			if cur == ifStmt.Else && condGuaranteesNil(ifStmt.Cond, accessExpr, v, info) && !nodeAssignsNilBefore(ifStmt.Else, useNode.Pos(), v, info) {
				return true
			}
		}
		if forStmt, ok := parents[cur].(*ast.ForStmt); ok && cur == forStmt.Body && forStmt.Cond != nil {
			if condGuaranteesNonNil(forStmt.Cond, accessExpr, v, info) && !stmtAssignsNilBetween(forStmt.Body.List, token.NoPos, useNode.Pos(), v, info) {
				return true
			}
		}
	}
	// 2. Check preceding statements in enclosing BlockStmt.
	for cur := useNode; cur != nil; cur = parents[cur] {
		if _, isFuncLit := cur.(*ast.FuncLit); isFuncLit {
			break
		}
		stmt, isStmt := cur.(ast.Stmt)
		block, isBlock := parents[cur].(*ast.BlockStmt)
		if !isStmt || !isBlock {
			continue
		}
		idx := slices.Index(block.List, stmt)
		for i := idx - 1; i >= 0; i-- {
			prev := block.List[i]
			if stmtAssignsNilBetween([]ast.Stmt{prev}, token.NoPos, useNode.Pos(), v, info) {
				return false
			}
			if stmtEstablishesNonNilWith(prev, accessExpr, v, info, isNonNilInitRhs) && !stmtAssignsNilBetween(block.List[i+1:idx+1], token.NoPos, useNode.Pos(), v, info) {
				return true
			}
		}
	}
	return false
}

func stmtEstablishesNonNilWith(stmt ast.Stmt, accessExpr ast.Expr, v *types.Var, info *types.Info, isNonNilInitRhs func(ast.Expr) bool) bool {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if s.Tok == token.ASSIGN && len(s.Lhs) == 1 && len(s.Rhs) == 1 && sameAccessPath(s.Lhs[0], accessExpr, v, info) {
			return isNonNilInitRhs != nil && isNonNilInitRhs(s.Rhs[0])
		}
	case *ast.IfStmt:
		if s.Init == nil && s.Else == nil && s.Body != nil && condGuaranteesNil(s.Cond, accessExpr, v, info) {
			if blockTerminates(s.Body) {
				return true
			}
			for _, bs := range s.Body.List {
				if stmtEstablishesNonNilWith(bs, accessExpr, v, info, isNonNilInitRhs) && !stmtAssignsNilBetween(s.Body.List, bs.End(), s.Body.End(), v, info) {
					return true
				}
			}
		}
	}
	return false
}

func blockTerminates(b *ast.BlockStmt) bool {
	if b == nil || len(b.List) == 0 {
		return false
	}
	switch last := b.List[len(b.List)-1].(type) {
	case *ast.ReturnStmt, *ast.BranchStmt:
		return true
	case *ast.ExprStmt:
		if call, ok := stripParens(last.X).(*ast.CallExpr); ok {
			if id, ok := stripParens(call.Fun).(*ast.Ident); ok && id.Name == "panic" {
				return true
			}
		}
	}
	return false
}

func condGuaranteesNonNil(cond, accessExpr ast.Expr, v *types.Var, info *types.Info) bool {
	c := stripParens(cond)
	if bin, ok := c.(*ast.BinaryExpr); ok {
		if bin.Op == token.LAND {
			return condGuaranteesNonNil(bin.X, accessExpr, v, info) || condGuaranteesNonNil(bin.Y, accessExpr, v, info)
		}
		if bin.Op == token.NEQ {
			return (sameAccessPath(bin.X, accessExpr, v, info) && isNilIdent(bin.Y, info)) ||
				(sameAccessPath(bin.Y, accessExpr, v, info) && isNilIdent(bin.X, info))
		}
	}
	return false
}

func condGuaranteesNil(cond, accessExpr ast.Expr, v *types.Var, info *types.Info) bool {
	c := stripParens(cond)
	if bin, ok := c.(*ast.BinaryExpr); ok && bin.Op == token.EQL {
		return (sameAccessPath(bin.X, accessExpr, v, info) && isNilIdent(bin.Y, info)) ||
			(sameAccessPath(bin.Y, accessExpr, v, info) && isNilIdent(bin.X, info))
	}
	return false
}

func sameAccessPath(a, b ast.Expr, v *types.Var, info *types.Info) bool {
	a, b = stripParens(a), stripParens(b)
	switch xa := a.(type) {
	case *ast.Ident:
		xb, ok := b.(*ast.Ident)
		return ok && varFromObj(info.Uses[xa]) == v && varFromObj(info.Uses[xb]) == v
	case *ast.SelectorExpr:
		xb, ok := b.(*ast.SelectorExpr)
		return ok && xa.Sel.Name == xb.Sel.Name && varFromObj(info.Uses[xa.Sel]) == varFromObj(info.Uses[xb.Sel]) && sameReceiverChain(xa.X, xb.X, info)
	}
	return false
}

func sameReceiverChain(a, b ast.Expr, info *types.Info) bool {
	a, b = stripParens(a), stripParens(b)
	switch xa := a.(type) {
	case *ast.Ident:
		xb, ok := b.(*ast.Ident)
		va := varFromObj(info.Uses[xa])
		return ok && va != nil && va == varFromObj(info.Uses[xb])
	case *ast.SelectorExpr:
		xb, ok := b.(*ast.SelectorExpr)
		va := varFromObj(info.Uses[xa.Sel])
		return ok && va != nil && va == varFromObj(info.Uses[xb.Sel]) && sameReceiverChain(xa.X, xb.X, info)
	}
	return false
}

func stmtAssignsNilBetween(stmts []ast.Stmt, afterPos, beforePos token.Pos, v *types.Var, info *types.Info) bool {
	for _, st := range stmts {
		if afterPos.IsValid() && st.Pos() < afterPos {
			continue
		}
		if beforePos.IsValid() && st.Pos() >= beforePos {
			break
		}
		if nodeAssignsNilBefore(st, beforePos, v, info) {
			return true
		}
	}
	return false
}

func nodeAssignsNilBefore(root ast.Node, beforePos token.Pos, v *types.Var, info *types.Info) bool {
	found := false
	ast.Inspect(root, func(n ast.Node) bool {
		if found || n == nil || (beforePos.IsValid() && n.Pos() >= beforePos) {
			return false
		}
		if as, ok := n.(*ast.AssignStmt); ok {
			for i, lhs := range as.Lhs {
				if exprReferencesVar(lhs, v, info) && (i >= len(as.Rhs) || isNilIdent(as.Rhs[i], info)) {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}

func isFormattingOnlySliceSink(expr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) bool {
	return isFormattingOnlySliceSinkDepth(expr, parents, info, 0)
}

func isFormattingOnlySliceSinkDepth(expr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info, depth int) bool {
	if depth > 4 {
		return false
	}
	cur := expr
	// Unwrap local `any` staging variable (`var refsView any = nil; ... refsView = this_.rest; ... Add_2_string_any("references", refsView)`).
	if as, ok := parents[cur].(*ast.AssignStmt); ok && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && as.Rhs[0] == cur {
		lhsId, okId := stripParens(as.Lhs[0]).(*ast.Ident)
		if !okId || lhsId.Name == "_" {
			return false
		}
		vObj := varFromObj(info.Uses[lhsId])
		fnBody := enclosingFuncBody(as, parents)
		if vObj == nil || vObj.IsField() || fnBody == nil || vObj.Pos() < fnBody.Pos() || vObj.Pos() > fnBody.End() || !isEmptyInterfaceType(vObj.Type()) {
			return false
		}
		for n := ast.Node(as); n != nil && n != fnBody; n = parents[n] {
			switch n.(type) {
			case *ast.ForStmt, *ast.RangeStmt, *ast.FuncLit:
				return false
			}
		}
		readCount := 0
		allSink := true
		ast.Inspect(fnBody, func(n ast.Node) bool {
			if !allSink || n == nil {
				return false
			}
			id, isId := n.(*ast.Ident)
			if !isId || varFromObj(info.Uses[id]) != vObj {
				return true
			}
			// Ignore LHS of `refsView = ...` and `_ = refsView`.
			if pAs, isAs := parents[id].(*ast.AssignStmt); isAs {
				if len(pAs.Lhs) == 1 && stripParens(pAs.Lhs[0]) == id {
					return true
				}
				if pAs.Tok == token.ASSIGN && len(pAs.Lhs) == 1 && len(pAs.Rhs) == 1 && stripParens(pAs.Rhs[0]) == id {
					if dId, okD := stripParens(pAs.Lhs[0]).(*ast.Ident); okD && dId.Name == "_" {
						return true
					}
				}
			}
			if id.Pos() <= as.End() {
				allSink = false
				return false
			}
			top := unwrapParensUp(id, parents)
			if !isFormattingOnlySliceSinkDepth(top, parents, info, depth+1) {
				allSink = false
				return false
			}
			readCount++
			return true
		})
		return allSink && readCount > 0
	}
	if ret, ok := parents[cur].(*ast.ReturnStmt); ok && len(ret.Results) == 1 && ret.Results[0] == cur {
		// Unwrap immediately-invoked ternary IIFE `func() any { ... return x }()`.
		var fnLit *ast.FuncLit
		for p := parents[ret]; p != nil; p = parents[p] {
			if fl, ok := p.(*ast.FuncLit); ok {
				fnLit = fl
				break
			}
			if _, ok := p.(*ast.FuncDecl); ok {
				break
			}
		}
		if fnLit == nil || fnLit.Type.Results == nil || len(fnLit.Type.Results.List) != 1 || !isEmptyInterfaceType(info.Types[fnLit.Type.Results.List[0].Type].Type) {
			return false
		}
		iifeExpr := ast.Expr(fnLit)
		for {
			pe, ok := parents[iifeExpr].(*ast.ParenExpr)
			if !ok || pe.X != iifeExpr {
				break
			}
			iifeExpr = pe
		}
		iifeCall, ok := parents[iifeExpr].(*ast.CallExpr)
		if !ok || iifeCall.Fun != iifeExpr || len(iifeCall.Args) != 0 {
			return false
		}
		cur = iifeCall
		for {
			pe, ok := parents[cur].(*ast.ParenExpr)
			if !ok || pe.X != cur {
				break
			}
			cur = pe
		}
	}

	call, ok := parents[cur].(*ast.CallExpr)
	if !ok || call.Fun == cur {
		return false
	}
	if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && (fnName == "StrVal" || fnName == "StrConcat") {
		return true
	}
	if fn, _ := calleeFunc(call, info); fn != nil && isJ2gortPkg(fn.Pkg()) && fn.Name() == "Add_2_string_any" {
		if sig, ok := fn.Type().(*types.Signature); ok && sig.Recv() != nil {
			if strings.Contains(sig.Recv().Type().String(), "ToStringHelper") && len(call.Args) == 2 && call.Args[1] == cur {
				return true
			}
		}
	}
	return false
}

func unwrapCollectionAdapters(expr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) ast.Expr {
	for {
		if pe, ok := parents[expr].(*ast.ParenExpr); ok && pe.X == expr {
			expr = pe
			continue
		}
		call, ok := parents[expr].(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() || call.Args[0] != expr {
			return expr
		}
		if _, isTI := isToIfaceCall(call, info); isTI && isCollectionOrSetIface(info.Types[call].Type) {
			expr = call
			continue
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "CheckNotNull", "CheckNotNullPtr", "CheckNotNullIface":
				expr = call
				continue
			case "AsIface", "Cast", "UncheckedCast":
				if isCollectionOrSetIface(info.Types[call].Type) {
					expr = call
					continue
				}
			}
		}
		return expr
	}
}

func unwrapCollectionAdaptersDown(expr ast.Expr, info *types.Info) ast.Expr {
	cur := stripParens(expr)
	for {
		call, ok := cur.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
			return cur
		}
		if _, isTI := isToIfaceCall(call, info); isTI && isCollectionOrSetIface(info.Types[call].Type) {
			cur = stripParens(call.Args[0])
			continue
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "CheckNotNull", "CheckNotNullPtr", "CheckNotNullIface":
				cur = stripParens(call.Args[0])
				continue
			case "AsIface", "Cast", "UncheckedCast":
				if isCollectionOrSetIface(info.Types[call].Type) {
					cur = stripParens(call.Args[0])
					continue
				}
			}
		}
		return cur
	}
}

func unwrapCheckNotNullArg(expr ast.Expr, info *types.Info) ast.Expr {
	cur := stripParens(expr)
	for {
		call, ok := cur.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
			return cur
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && strings.HasPrefix(fnName, "CheckNotNull") {
			cur = stripParens(call.Args[0])
			continue
		}
		return cur
	}
}

func unwrapIteratorCast(call *ast.CallExpr, parents map[ast.Node]ast.Node, info *types.Info) ast.Expr {
	cur := ast.Expr(call)
	for {
		if pe, ok := parents[cur].(*ast.ParenExpr); ok && pe.X == cur {
			cur = pe
			continue
		}
		if c, ok := parents[cur].(*ast.CallExpr); ok && len(c.Args) == 1 && !c.Ellipsis.IsValid() && c.Args[0] == cur {
			if fnName, _, _, isJ2 := isJ2gortCall(c, info); isJ2 && (fnName == "AsIface" || fnName == "Cast") {
				if named, ok := types.Unalias(info.Types[c].Type).(*types.Named); ok && named.Obj() != nil && isJ2gortPkg(named.Obj().Pkg()) {
					if named.Obj().Name() == "Ext_java_util_Iterator" || named.Obj().Name() == "Ext_com_google_common_collect_UnmodifiableIterator" {
						return c
					}
				}
			}
		}
		return call
	}
}

func isCollectionOrSetIface(t types.Type) bool {
	named, ok := types.Unalias(t).(*types.Named)
	if !ok || named.Obj() == nil || !isJ2gortPkg(named.Obj().Pkg()) {
		return false
	}
	switch named.Obj().Name() {
	case "Ext_java_util_Collection", "Ext_java_lang_Iterable", "Ext_java_util_SequencedCollection",
		"Ext_java_util_List", "Ext_java_util_Set", "Ext_java_util_SequencedSet":
		return true
	}
	return false
}

func isNilIdent(expr ast.Expr, info *types.Info) bool {
	id, ok := stripParens(expr).(*ast.Ident)
	if !ok || id.Name != "nil" {
		return false
	}
	_, isNil := info.Uses[id].(*types.Nil)
	return isNil
}

func isNonNegativeConstInt(expr ast.Expr, info *types.Info) bool {
	tv, ok := info.Types[expr]
	if !ok || tv.Value == nil || tv.Value.Kind() != constant.Int {
		return false
	}
	return constant.Sign(tv.Value) >= 0
}

func isCandidateListType(t types.Type) bool {
	if t == nil {
		return false
	}
	ut := types.Unalias(t)
	if ptr, ok := ut.(*types.Pointer); ok {
		ut = types.Unalias(ptr.Elem())
	}
	named, ok := ut.(*types.Named)
	if !ok || named.Obj() == nil || !isJ2gortPkg(named.Obj().Pkg()) {
		return false
	}
	switch named.Obj().Name() {
	case "Ext_java_util_ArrayList", "Ext_java_util_List", "Ext_java_util_SequencedCollection", "Ext_java_util_Collection", "Ext_java_lang_Iterable":
		return true
	}
	return false
}

func isCandidateAccParamType(t types.Type) bool {
	if t == nil {
		return false
	}
	ut := types.Unalias(t)
	if ptr, ok := ut.(*types.Pointer); ok {
		ut = types.Unalias(ptr.Elem())
	}
	named, ok := ut.(*types.Named)
	if !ok || named.Obj() == nil || !isJ2gortPkg(named.Obj().Pkg()) {
		return false
	}
	switch named.Obj().Name() {
	case "Ext_java_util_ArrayList", "Ext_java_util_List", "Ext_java_util_Collection":
		return true
	}
	return false
}

func classifySliceInitExpr(
	expr ast.Expr,
	selfCand *sliceCandidate,
	candidates map[*types.Var]*sliceCandidate,
	hasCopyHelpers bool,
	ctx *PassContext,
	info *types.Info,
) (kind sliceInitKind, argExpr ast.Expr, otherCand *sliceCandidate, typedSliceElem types.Type, ok bool) {
	call, isCall := stripParens(expr).(*ast.CallExpr)
	if !isCall || call.Ellipsis.IsValid() {
		return 0, nil, nil, nil, false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 {
		return 0, nil, nil, nil, false
	}
	switch fnName {
	case "NewExt_java_util_ArrayList", "Ext_com_google_common_collect_Lists_NewArrayList_0":
		if len(call.Args) == 0 {
			return sliceInitEmpty, nil, nil, nil, true
		}
	case "NewExt_java_util_ArrayList_1_int32", "Ext_com_google_common_collect_Lists_NewArrayListWithCapacity":
		if len(call.Args) == 1 && isIntegerExpr(call.Args[0], info) {
			if isNonNegativeConstInt(call.Args[0], info) && ctx != nil && !isUnshadowedBuiltin("make", expr.Pos(), ctx.Pkg) {
				return 0, nil, nil, nil, false
			}
			return sliceInitCap, call.Args[0], nil, nil, true
		}
	case "NewExt_java_util_ArrayList_1_Ext_java_util_Collection",
		"Ext_com_google_common_collect_Lists_NewArrayList_1_Ext_java_lang_Iterable":
		if hasCopyHelpers && len(call.Args) == 1 {
			inner := unwrapCollectionAdaptersDown(call.Args[0], info)
			if candidates != nil {
				if sc := resolveCandidateAccess(inner, candidates, info); sc != nil && sc != selfCand && !sc.isAccParam {
					return sliceInitCopySlice, inner, sc, nil, true
				}
			}
			if slExpr, elemT, okS := matchWrappedSliceArg(inner, info); okS {
				return sliceInitCopyTypedSlice, slExpr, nil, elemT, true
			}
			return sliceInitCopyCol, call.Args[0], nil, nil, true
		}
	case "Ext_com_google_common_collect_Lists_NewArrayList_1_Ext_java_util_Iterator":
		if hasCopyHelpers && len(call.Args) == 1 {
			return sliceInitCopyIter, call.Args[0], nil, nil, true
		}
	}
	return 0, nil, nil, nil, false
}

func isSliceCopyTargetCandidate(initCall *ast.CallExpr, candidates map[*types.Var]*sliceCandidate, parents map[ast.Node]ast.Node, info *types.Info) bool {
	cur := unwrapParensUp(initCall, parents)
	switch p := parents[cur].(type) {
	case *ast.ValueSpec:
		return len(p.Names) == 1 && len(p.Values) == 1 && p.Values[0] == cur && candidates[varFromObj(info.Defs[p.Names[0]])] != nil
	case *ast.AssignStmt:
		if len(p.Lhs) != 1 || len(p.Rhs) != 1 || p.Rhs[0] != cur {
			return false
		}
		if p.Tok == token.DEFINE {
			id, ok := p.Lhs[0].(*ast.Ident)
			return ok && candidates[varFromObj(info.Defs[id])] != nil
		}
		if p.Tok == token.ASSIGN {
			switch lhs := stripParens(p.Lhs[0]).(type) {
			case *ast.Ident:
				return candidates[varFromObj(info.Uses[lhs])] != nil
			case *ast.SelectorExpr:
				return candidates[varFromObj(info.Uses[lhs.Sel])] != nil
			}
		}
	}
	return false
}

func unwrapUnmodifiableListCall(expr ast.Expr, info *types.Info) ast.Expr {
	inner := unwrapCollectionAdaptersDown(expr, info)
	if call, ok := inner.(*ast.CallExpr); ok && len(call.Args) == 1 && !call.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && fnName == "Ext_java_util_Collections_UnmodifiableList" {
			return unwrapCollectionAdaptersDown(call.Args[0], info)
		}
	}
	return stripParens(expr)
}

func isListsReverseOfCandidate(expr ast.Expr, cand *sliceCandidate, info *types.Info) bool {
	inner := unwrapCollectionAdaptersDown(expr, info)
	call, ok := inner.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
		return false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 || fnName != "Ext_com_google_common_collect_Lists_Reverse" {
		return false
	}
	arg := unwrapCollectionAdaptersDown(call.Args[0], info)
	id, ok := arg.(*ast.Ident)
	return ok && varFromObj(info.Uses[id]) == cand.v
}

func isReadOnlyIterableConsumerCall(call *ast.CallExpr, adaptedTop ast.Expr, info *types.Info) bool {
	callee, _ := calleeFunc(call, info)
	if callee == nil || callee.Pkg() == nil {
		return false
	}
	if isJ2gortPkg(callee.Pkg()) {
		switch callee.Name() {
		case "Join_1_Ext_java_lang_Iterable":
			return len(call.Args) == 1 && call.Args[0] == adaptedTop
		case "AppendTo_2_Ext_java_lang_StringBuilder_Ext_java_lang_Iterable":
			return len(call.Args) == 2 && call.Args[1] == adaptedTop
		}
	}
	if callee.Name() == "AppendAll" && len(call.Args) == 2 && call.Args[0] == adaptedTop {
		if sig, ok := callee.Type().(*types.Signature); ok && sig.Recv() != nil {
			if recvNamed := ptrNamedOrigin(sig.Recv().Type()); recvNamed != nil && recvNamed.Obj() != nil && recvNamed.Obj().Name() == "TypeStringBuilder" {
				return true
			}
		}
	}
	return false
}

func freshLoopIndexName(body *ast.BlockStmt, base string) string {
	used := make(map[string]bool)
	ast.Inspect(body, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			used[id.Name] = true
		}
		return true
	})
	name := base
	for i := 2; used[name]; i++ {
		name = base + "_" + strconv.Itoa(i)
	}
	return name
}

func boxSliceElemForJavaList(elemExpr ast.Expr, elemT types.Type, j2Sel func(string, token.Pos) *ast.SelectorExpr) ast.Expr {
	ut := types.Unalias(elemT)
	if types.Identical(elemT, types.Typ[types.String]) {
		return &ast.CallExpr{Fun: j2Sel("BoxString", token.NoPos), Args: []ast.Expr{elemExpr}}
	}
	if isNamedUint32NodeType(elemT) {
		return &ast.CallExpr{Fun: j2Sel("BoxNode", token.NoPos), Args: []ast.Expr{elemExpr}}
	}
	switch b := ut.(type) {
	case *types.Pointer:
		return &ast.CallExpr{Fun: j2Sel("PtrToAny", token.NoPos), Args: []ast.Expr{elemExpr}}
	case *types.Slice:
		return &ast.CallExpr{Fun: j2Sel("SliceToAny", token.NoPos), Args: []ast.Expr{elemExpr}}
	case *types.Basic:
		switch b.Kind() {
		case types.Int32:
			return &ast.CallExpr{Fun: j2Sel("ValueOfInt", token.NoPos), Args: []ast.Expr{elemExpr}}
		case types.Int64:
			return &ast.CallExpr{Fun: j2Sel("ValueOfLong", token.NoPos), Args: []ast.Expr{elemExpr}}
		case types.Bool:
			return &ast.CallExpr{Fun: j2Sel("ValueOfBoolean", token.NoPos), Args: []ast.Expr{elemExpr}}
		case types.Float64:
			return &ast.CallExpr{Fun: j2Sel("ValueOfDouble", token.NoPos), Args: []ast.Expr{elemExpr}}
		case types.Float32:
			return &ast.CallExpr{Fun: j2Sel("ValueOfFloat", token.NoPos), Args: []ast.Expr{elemExpr}}
		case types.Uint16:
			return &ast.CallExpr{Fun: j2Sel("ValueOfChar", token.NoPos), Args: []ast.Expr{elemExpr}}
		case types.Int8:
			return &ast.CallExpr{Fun: j2Sel("ValueOfByte", token.NoPos), Args: []ast.Expr{elemExpr}}
		case types.Int16:
			return &ast.CallExpr{Fun: j2Sel("ValueOfShort", token.NoPos), Args: []ast.Expr{elemExpr}}
		}
	}
	if isNonEmptyInterfaceType(elemT) {
		return elemExpr
	}
	if types.Comparable(elemT) && !isInterfaceType(elemT) && types.NewMethodSet(elemT).Lookup(nil, "IsJavaNil_") != nil {
		return &ast.CallExpr{Fun: j2Sel("BoxNullableVal", token.NoPos), Args: []ast.Expr{elemExpr}}
	}
	return &ast.CallExpr{Fun: j2Sel("ToAny", token.NoPos), Args: []ast.Expr{elemExpr}}
}

func isNewArrayListCall(expr ast.Expr, info *types.Info) (ast.Expr, bool, bool) {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return nil, false, false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 {
		return nil, false, false
	}
	if (fnName == "NewExt_java_util_ArrayList" || fnName == "Ext_com_google_common_collect_Lists_NewArrayList_0") && len(call.Args) == 0 {
		return nil, false, true
	}
	if (fnName == "NewExt_java_util_ArrayList_1_int32" || fnName == "Ext_com_google_common_collect_Lists_NewArrayListWithCapacity") && len(call.Args) == 1 && isIntegerExpr(call.Args[0], info) {
		return call.Args[0], true, true
	}
	return nil, false, false
}

func primBoxKind(fnName string) (types.BasicKind, bool) {
	for _, pfx := range []string{"ValueOf", "Box"} {
		if strings.HasPrefix(fnName, pfx) {
			k, ok := unboxWantKind["Unbox"+strings.TrimPrefix(fnName, pfx)]
			return k, ok
		}
	}
	return 0, false
}

func unwrapBoxedAddArg(arg ast.Expr, info *types.Info) (ast.Expr, types.Type, bool) {
	e := stripParens(arg)
	if call, ok := e.(*ast.CallExpr); ok && len(call.Args) == 1 && !call.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			inner := stripParens(call.Args[0])
			tv, hasTV := info.Types[inner]
			switch fnName {
			case "BoxString":
				return inner, types.Typ[types.String], hasTV && hasExactOrDefaultStringType(inner, info)
			case "BoxNode", "BoxByteString":
				if hasTV && !tv.IsNil() && tv.Type != nil && !isInvalidOrUntyped(tv.Type) {
					return inner, tv.Type, true
				}
				return nil, nil, false
			case "PtrToAny":
				if hasTV && !tv.IsNil() && tv.Type != nil {
					_, ok := types.Unalias(tv.Type).(*types.Pointer)
					return inner, tv.Type, ok
				}
				return nil, nil, false
			case "SliceToAny":
				if hasTV && !tv.IsNil() && tv.Type != nil {
					_, ok := types.Unalias(tv.Type).(*types.Slice)
					return inner, tv.Type, ok
				}
				return nil, nil, false
			case "ToAny", "BoxNullableVal":
				if hasTV && !tv.IsNil() && tv.Type != nil {
					dt := types.Default(tv.Type)
					return inner, dt, !isInvalidOrUntyped(dt) && !isInterfaceType(dt)
				}
				return nil, nil, false
			default:
				if kind, ok := primBoxKind(fnName); ok {
					return inner, types.Typ[kind], hasTV && !tv.IsNil() && tv.Type != nil && types.AssignableTo(tv.Type, types.Typ[kind])
				}
			}
		}
	}
	tv, ok := info.Types[e]
	if !ok || tv.IsType() || tv.IsBuiltin() || tv.IsNil() || tv.Type == nil {
		return nil, nil, false
	}
	dt := types.Default(tv.Type)
	return e, dt, !isInvalidOrUntyped(dt)
}

func enclosingUnboxExpr(rawExpr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) (ast.Expr, types.Type, bool) {
	cur := rawExpr
	for {
		pe, ok := parents[cur].(*ast.ParenExpr)
		if !ok || pe.X != cur {
			break
		}
		cur = pe
	}
	call, ok := parents[cur].(*ast.CallExpr)
	if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() || call.Args[0] != cur {
		return rawExpr, types.Universe.Lookup("any").Type(), true
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 {
		return rawExpr, types.Universe.Lookup("any").Type(), true
	}
	switch fnName {
	case "CastPtr", "Cast", "CastVal", "UncheckedCast":
		topCast := ast.Expr(call)
		for {
			pe, ok := parents[topCast].(*ast.ParenExpr)
			if !ok || pe.X != topCast {
				break
			}
			topCast = pe
		}
		if outerCall, ok := parents[topCast].(*ast.CallExpr); ok && len(outerCall.Args) == 1 && !outerCall.Ellipsis.IsValid() && outerCall.Args[0] == topCast {
			if outerFn, _, _, ok := isJ2gortCall(outerCall, info); ok {
				if wantKind, isUnbox := unboxWantKind[outerFn]; isUnbox {
					return outerCall, types.Typ[wantKind], true
				}
			}
		}
		t := info.Types[call].Type
		return call, t, !isInvalidOrUntyped(t)
	case "AsIface", "CastString", "CastNode", "UnboxNode", "Int32Val",
		"UnboxInt", "UnboxLong", "UnboxBoolean", "UnboxDouble", "UnboxFloat", "UnboxChar", "UnboxByte", "UnboxShort":
		t := info.Types[call].Type
		return call, t, !isInvalidOrUntyped(t)
	}
	return rawExpr, types.Universe.Lookup("any").Type(), true
}

func unwrapEnclosingCheckNotNulls(expr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) ast.Expr {
	for {
		if pe, ok := parents[expr].(*ast.ParenExpr); ok && pe.X == expr {
			expr = pe
			continue
		}
		if call, ok := parents[expr].(*ast.CallExpr); ok && len(call.Args) == 1 && !call.Ellipsis.IsValid() && call.Args[0] == expr {
			if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && (fnName == "CheckNotNull" || fnName == "CheckNotNullPtr" || fnName == "CheckNotNullIface") {
				expr = call
				continue
			}
		}
		return expr
	}
}

func j2gortQualifierForFile(file *ast.File, ctx *PassContext) (string, bool) {
	for _, spec := range file.Imports {
		if spec.Name != nil && (spec.Name.Name == "_" || spec.Name.Name == ".") {
			continue
		}
		p, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		isJ2 := p == "j2gort" || strings.HasSuffix(p, "/j2gort")
		if !isJ2 && ctx.Importer != nil {
			if impPkg, impErr := ctx.Importer.Import(p); impErr == nil && isJ2gortPkg(impPkg) {
				isJ2 = true
			}
		}
		if isJ2 {
			if local := specLocalName(spec, ctx.Importer); local != "" {
				return local, true
			}
		}
	}
	return "", false
}

func canSpellJ2gortAtAllSites(sites []typeSpellSite, ctx *PassContext) bool {
	for _, site := range sites {
		qual, ok := j2gortQualifierForFile(site.file, ctx)
		if !ok {
			return false
		}
		if ctx.Pkg != nil && ctx.Pkg.Scope() != nil {
			_, obj := ctx.Pkg.Scope().Innermost(site.pos).LookupParent(qual, site.pos)
			if _, isPkg := obj.(*types.PkgName); !isPkg {
				return false
			}
		}
	}
	return true
}

func collectTypeParamsMap(t types.Type) map[string]*types.TypeParam {
	if !typeContainsTypeParam(t) {
		return nil
	}
	out := make(map[string]*types.TypeParam)
	var walk func(types.Type)
	walk = func(cur types.Type) {
		if cur == nil {
			return
		}
		switch ut := types.Unalias(cur).(type) {
		case *types.TypeParam:
			if ut.Obj() != nil {
				out[ut.Obj().Name()] = ut
			}
		case *types.Pointer:
			walk(ut.Elem())
		case *types.Slice:
			walk(ut.Elem())
		case *types.Array:
			walk(ut.Elem())
		case *types.Named:
			if targs := ut.TypeArgs(); targs != nil {
				for i := 0; i < targs.Len(); i++ {
					walk(targs.At(i))
				}
			}
		}
	}
	walk(t)
	return out
}

func canSpellTypeAtAllSites(t types.Type, sites []typeSpellSite, ctx *PassContext) bool {
	for _, site := range sites {
		siteT := t
		isStructFieldDecl := site.owner != nil && site.accessExpr == nil && site.recvTArgs == nil
		if site.owner != nil && site.recvTArgs != nil {
			if inst, ok := instantiateOwnerType(t, site.owner, site.recvTArgs); ok {
				siteT = inst
			}
		} else if site.owner != nil && site.accessExpr != nil {
			siteT = instantiateOwnerTypeAtExpr(t, site.owner, site.accessExpr, ctx.Info)
		}
		expr, ok := typeExprForFile(siteT, site.file, ctx.Pkg, ctx.Importer)
		if !ok || !isTypeExprUnshadowedWithTP(expr, site.pos, ctx.Pkg, collectTypeParamsMap(siteT), isStructFieldDecl) {
			return false
		}
	}
	return true
}

func typeExprForFile(t types.Type, file *ast.File, currentPkg *types.Package, imp types.Importer) (ast.Expr, bool) {
	if alias, isAlias := t.(*types.Alias); isAlias {
		return namedTypeExpr(alias.Obj(), alias.TypeArgs(), file, currentPkg, imp)
	}
	switch ut := t.(type) {
	case *types.Basic:
		if ut.Kind() == types.Invalid || (ut.Info()&types.IsUntyped) != 0 {
			return nil, false
		}
		return &ast.Ident{Name: ut.Name()}, true
	case *types.TypeParam:
		if ut.Obj() == nil {
			return nil, false
		}
		return &ast.Ident{Name: ut.Obj().Name()}, true
	case *types.Interface:
		if ut.Empty() {
			return &ast.Ident{Name: "any"}, true
		}
	case *types.Pointer:
		if elt, ok := typeExprForFile(ut.Elem(), file, currentPkg, imp); ok {
			return &ast.StarExpr{X: elt}, true
		}
	case *types.Slice:
		if elt, ok := typeExprForFile(ut.Elem(), file, currentPkg, imp); ok {
			return &ast.ArrayType{Elt: elt}, true
		}
	case *types.Named:
		return namedTypeExpr(ut.Obj(), ut.TypeArgs(), file, currentPkg, imp)
	}
	return nil, false
}

func namedTypeExpr(obj *types.TypeName, targs *types.TypeList, file *ast.File, currentPkg *types.Package, imp types.Importer) (ast.Expr, bool) {
	if obj == nil || (obj.Pkg() != nil && obj.Parent() != obj.Pkg().Scope()) {
		return nil, false
	}
	var base ast.Expr
	if obj.Pkg() == nil || obj.Pkg() == currentPkg {
		base = &ast.Ident{Name: obj.Name()}
	} else {
		if !obj.Exported() {
			return nil, false
		}
		qual := findImportedQualifier(file, obj.Pkg().Path(), imp)
		if qual == "" {
			return nil, false
		}
		base = &ast.SelectorExpr{X: &ast.Ident{Name: qual}, Sel: &ast.Ident{Name: obj.Name()}}
	}
	if targs != nil && targs.Len() > 0 {
		args := make([]ast.Expr, targs.Len())
		for i := 0; i < targs.Len(); i++ {
			ae, ok := typeExprForFile(targs.At(i), file, currentPkg, imp)
			if !ok {
				return nil, false
			}
			args[i] = ae
		}
		if len(args) == 1 {
			return &ast.IndexExpr{X: base, Index: args[0]}, true
		}
		return &ast.IndexListExpr{X: base, Indices: args}, true
	}
	return base, true
}

func findImportedQualifier(file *ast.File, targetPath string, imp types.Importer) string {
	for _, spec := range file.Imports {
		if spec.Name != nil && (spec.Name.Name == "_" || spec.Name.Name == ".") {
			continue
		}
		if p, err := strconv.Unquote(spec.Path.Value); err == nil && p == targetPath {
			if local := specLocalName(spec, imp); local != "" {
				return local
			}
		}
	}
	return ""
}

func isTypeExprUnshadowed(expr ast.Expr, pos token.Pos, pkg *types.Package) bool {
	return isTypeExprUnshadowedWithTP(expr, pos, pkg, nil, false)
}

func isTypeExprUnshadowedWithTP(expr ast.Expr, pos token.Pos, pkg *types.Package, tpMap map[string]*types.TypeParam, isStructFieldDecl bool) bool {
	if pkg == nil || pkg.Scope() == nil {
		return true
	}
	scope := pkg.Scope().Innermost(pos)
	ok := true
	var check func(ast.Expr)
	check = func(e ast.Expr) {
		if !ok || e == nil {
			return
		}
		switch x := e.(type) {
		case *ast.Ident:
			if wantTP := tpMap[x.Name]; wantTP != nil {
				if isStructFieldDecl {
					return
				}
				_, obj := scope.LookupParent(x.Name, pos)
				if tn, isTN := obj.(*types.TypeName); isTN && types.Identical(tn.Type(), wantTP) {
					return
				}
				if obj == nil {
					return
				}
				ok = false
				return
			}
			_, obj := scope.LookupParent(x.Name, pos)
			if _, isTN := obj.(*types.TypeName); isTN {
				return
			}
			ok = false
		case *ast.SelectorExpr:
			id, isId := x.X.(*ast.Ident)
			if !isId {
				ok = false
				return
			}
			_, obj := scope.LookupParent(id.Name, pos)
			_, isPkg := obj.(*types.PkgName)
			ok = ok && isPkg
		case *ast.StarExpr:
			check(x.X)
		case *ast.ArrayType:
			check(x.Elt)
		case *ast.IndexExpr:
			check(x.X)
			check(x.Index)
		case *ast.IndexListExpr:
			check(x.X)
			for _, ix := range x.Indices {
				check(ix)
			}
		default:
			ok = false
		}
	}
	check(expr)
	return ok
}

func isUnshadowedBuiltin(name string, pos token.Pos, pkg *types.Package) bool {
	if pkg == nil || pkg.Scope() == nil {
		return true
	}
	_, obj := pkg.Scope().Innermost(pos).LookupParent(name, pos)
	_, ok := obj.(*types.Builtin)
	return ok
}

func isUnshadowedUniverseType(name string, pos token.Pos, pkg *types.Package) bool {
	if pkg == nil || pkg.Scope() == nil {
		return true
	}
	_, obj := pkg.Scope().Innermost(pos).LookupParent(name, pos)
	return obj == types.Universe.Lookup(name)
}

func typeNeedsGCZeroing(t types.Type) bool {
	switch u := types.Unalias(t).Underlying().(type) {
	case *types.Basic:
		return (u.Info() & (types.IsBoolean | types.IsNumeric)) == 0
	case *types.Array:
		return typeNeedsGCZeroing(u.Elem())
	case *types.Struct:
		for i := 0; i < u.NumFields(); i++ {
			if typeNeedsGCZeroing(u.Field(i).Type()) {
				return true
			}
		}
		return false
	default:
		return true
	}
}

func isIntegerExpr(expr ast.Expr, info *types.Info) bool {
	tv, ok := info.Types[expr]
	if !ok || tv.Type == nil {
		return false
	}
	b, ok := types.Unalias(tv.Type).Underlying().(*types.Basic)
	return ok && (b.Info()&types.IsInteger) != 0
}

func isPureFieldSelector(sel *ast.SelectorExpr, info *types.Info) bool {
	seln, ok := info.Selections[sel]
	if !ok || seln.Kind() != types.FieldVal {
		return false
	}
	_, recvPtr := types.Unalias(seln.Recv()).(*types.Pointer)
	ptrBacked := seln.Indirect() || recvPtr
	if !ptrBacked {
		if id, ok := stripParens(sel.X).(*ast.Ident); ok {
			_, ptrBacked = info.Uses[id].(*types.Var)
		}
	}
	return ptrBacked && isPureSelectorChain(sel.X, info)
}

func isPureSelectorChain(expr ast.Expr, info *types.Info) bool {
	switch x := stripParens(expr).(type) {
	case *ast.Ident:
		_, ok := info.Uses[x].(*types.Var)
		return ok
	case *ast.SelectorExpr:
		seln, ok := info.Selections[x]
		return ok && seln.Kind() == types.FieldVal && isPureSelectorChain(x.X, info)
	}
	return false
}

func hasClosureMutation(expr ast.Expr) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			found = true
			return false
		}
		return !found
	})
	return found
}

func exprReferencesVar(expr ast.Expr, target *types.Var, info *types.Info) bool {
	cur := stripParens(expr)
	for {
		call, ok := cur.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
			break
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && strings.HasPrefix(fnName, "CheckNotNull") {
			cur = stripParens(call.Args[0])
			continue
		}
		break
	}
	switch x := cur.(type) {
	case *ast.Ident:
		return varFromObj(info.Uses[x]) == target
	case *ast.SelectorExpr:
		return varFromObj(info.Uses[x.Sel]) == target
	}
	return false
}

func isCursorMethodCall(expr ast.Expr, curVar *types.Var, methodName string, info *types.Info) bool {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 0 || call.Ellipsis.IsValid() {
		return false
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != methodName {
		return false
	}
	id, ok := stripParens(sel.X).(*ast.Ident)
	return ok && varFromObj(info.Uses[id]) == curVar
}

func isNodeConvOfCursorNextNode(expr ast.Expr, curVar *types.Var, info *types.Info) bool {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
		return false
	}
	funTV := info.Types[call.Fun]
	return funTV.IsType() && isNamedUint32NodeType(funTV.Type) && isCursorMethodCall(call.Args[0], curVar, "NextNode", info)
}

func varFromObj(obj types.Object) *types.Var {
	if v, ok := obj.(*types.Var); ok && v != nil {
		return v.Origin()
	}
	return nil
}

func enclosingFuncObj(n ast.Node, parents map[ast.Node]ast.Node, info *types.Info) *types.Func {
	for cur := n; cur != nil; cur = parents[cur] {
		if fn, ok := cur.(*ast.FuncDecl); ok {
			obj, _ := info.Defs[fn.Name].(*types.Func)
			return obj
		}
	}
	return nil
}

func clonePureExpr(expr ast.Expr) ast.Expr {
	switch e := stripParens(expr).(type) {
	case *ast.Ident:
		return &ast.Ident{NamePos: e.NamePos, Name: e.Name}
	case *ast.SelectorExpr:
		return &ast.SelectorExpr{X: clonePureExpr(e.X), Sel: &ast.Ident{NamePos: e.Sel.NamePos, Name: e.Sel.Name}}
	}
	return expr
}

func wrapPrimaryTarget(expr ast.Expr) ast.Expr {
	switch expr.(type) {
	case *ast.StarExpr, *ast.UnaryExpr, *ast.BinaryExpr:
		return &ast.ParenExpr{X: expr}
	}
	return expr
}

func canReplaceStmtIn(parent ast.Node) bool {
	switch parent.(type) {
	case *ast.BlockStmt, *ast.CaseClause, *ast.CommClause, *ast.LabeledStmt:
		return true
	}
	return false
}

func replaceStmtInParent(parent ast.Node, oldStmt ast.Stmt, newStmts []ast.Stmt, parents map[ast.Node]ast.Node) bool {
	splice := func(list []ast.Stmt) ([]ast.Stmt, bool) {
		idx := slices.Index(list, oldStmt)
		if idx < 0 {
			return list, false
		}
		return append(append(append([]ast.Stmt(nil), list[:idx]...), newStmts...), list[idx+1:]...), true
	}
	ok := false
	switch p := parent.(type) {
	case *ast.BlockStmt:
		p.List, ok = splice(p.List)
	case *ast.CaseClause:
		p.Body, ok = splice(p.Body)
	case *ast.CommClause:
		p.Body, ok = splice(p.Body)
	case *ast.LabeledStmt:
		if p.Stmt == oldStmt {
			if len(newStmts) == 1 {
				p.Stmt = newStmts[0]
			} else {
				p.Stmt = &ast.BlockStmt{List: newStmts}
			}
			ok = true
		}
	}
	if ok && parents != nil {
		for _, s := range newStmts {
			parents[s] = parent
		}
	}
	return ok
}

func validateLocalReturnWrapTypedArrayList(
	ret *ast.ReturnStmt,
	adaptedTop ast.Expr,
	cand *sliceCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	idx := slices.Index(ret.Results, adaptedTop)
	if idx < 0 || cand.v == nil {
		return false
	}
	// Find immediate enclosing function or FuncLit.
	var fnBody *ast.BlockStmt
	var sig *types.Signature
	for cur := ast.Node(ret); cur != nil; cur = parents[cur] {
		switch fn := cur.(type) {
		case *ast.FuncDecl:
			fnBody = fn.Body
			if fnObj, ok := info.Defs[fn.Name].(*types.Func); ok && fnObj != nil {
				sig, _ = fnObj.Type().(*types.Signature)
			}
		case *ast.FuncLit:
			fnBody = fn.Body
			sig, _ = types.Unalias(info.Types[fn].Type).(*types.Signature)
		}
		if fnBody != nil {
			break
		}
	}
	if fnBody == nil || sig == nil || sig.Results() == nil || idx >= sig.Results().Len() {
		return false
	}
	// Ensure cand.v is declared inside this same function/closure (not captured from an outer function).
	if cand.v.Pos() < fnBody.Pos() || cand.v.Pos() > fnBody.End() {
		return false
	}
	if !isTypedArrayListReturnTarget(sig.Results().At(idx).Type()) {
		return false
	}
	// Ensure no defer in fnBody touches cand.v after return expression evaluation.
	hasDeferRef := false
	ast.Inspect(fnBody, func(n ast.Node) bool {
		if hasDeferRef || n == nil {
			return false
		}
		if defStmt, ok := n.(*ast.DeferStmt); ok {
			ast.Inspect(defStmt, func(dn ast.Node) bool {
				if id, isId := dn.(*ast.Ident); isId && varFromObj(info.Uses[id]) == cand.v {
					hasDeferRef = true
					return false
				}
				return !hasDeferRef
			})
			return false
		}
		return true
	})
	return !hasDeferRef
}

func isTypedArrayListReturnTarget(t types.Type) bool {
	if t == nil {
		return false
	}
	ut := types.Unalias(t)
	if isEmptyInterfaceType(ut) {
		return true
	}
	named, ok := ut.(*types.Named)
	if !ok || named.Obj() == nil || !isJ2gortPkg(named.Obj().Pkg()) {
		return false
	}
	switch named.Obj().Name() {
	case "Ext_java_util_ArrayList", "Ext_java_util_AbstractList", "Ext_java_util_List",
		"Ext_java_util_SequencedCollection", "Ext_java_util_AbstractCollection",
		"Ext_java_util_Collection", "Ext_java_lang_Iterable", "Ext_java_util_RandomAccess":
		return true
	}
	return false
}

func exprExpectsUnmodifiableIterator(outer ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) bool {
	isUnmodIterType := func(t types.Type) bool {
		if t == nil {
			return false
		}
		named, ok := types.Unalias(t).(*types.Named)
		return ok && named.Obj() != nil && isJ2gortPkg(named.Obj().Pkg()) &&
			named.Obj().Name() == "Ext_com_google_common_collect_UnmodifiableIterator"
	}
	if tv, ok := info.Types[outer]; ok && isUnmodIterType(tv.Type) {
		// Only if outer itself was an explicit AsIface[UnmodifiableIterator] cast!
		if call, okC := stripParens(outer).(*ast.CallExpr); okC {
			if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && (fnName == "AsIface" || fnName == "Cast") {
				return true
			}
		}
	}
	cur := unwrapParensUp(outer, parents)
	switch p := parents[cur].(type) {
	case *ast.ValueSpec:
		if p.Type != nil && isUnmodIterType(info.Types[p.Type].Type) {
			return true
		}
	case *ast.AssignStmt:
		for i, rhs := range p.Rhs {
			if rhs == cur && i < len(p.Lhs) && isUnmodIterType(info.Types[p.Lhs[i]].Type) {
				return true
			}
		}
	case *ast.ReturnStmt:
		idx := slices.Index(p.Results, cur)
		if idx >= 0 {
			for anc := parents[p]; anc != nil; anc = parents[anc] {
				if fd, ok := anc.(*ast.FuncDecl); ok {
					if fnObj, ok := info.Defs[fd.Name].(*types.Func); ok && fnObj != nil {
						if sig, ok := fnObj.Type().(*types.Signature); ok && sig.Results() != nil && idx < sig.Results().Len() {
							return isUnmodIterType(sig.Results().At(idx).Type())
						}
					}
					break
				}
				if fl, ok := anc.(*ast.FuncLit); ok {
					if sig, ok := types.Unalias(info.Types[fl].Type).(*types.Signature); ok && sig.Results() != nil && idx < sig.Results().Len() {
						return isUnmodIterType(sig.Results().At(idx).Type())
					}
					break
				}
			}
		}
	case *ast.CallExpr:
		argIdx := slices.Index(p.Args, cur)
		if argIdx >= 0 {
			if sig, ok := types.Unalias(info.Types[p.Fun].Type).(*types.Signature); ok && !sig.Variadic() && argIdx < sig.Params().Len() {
				return isUnmodIterType(sig.Params().At(argIdx).Type())
			}
		}
	}
	return false
}

func isSafeReadOnlySliceIteratorUse(
	outer ast.Expr,
	cand *sliceCandidate,
	ctx *PassContext,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	cur := unwrapParensUp(outer, parents)

	// 1. Check if assigned to or initializing a local iterator variable: disallow if `.Remove()` is called on it.
	var itVar *types.Var
	switch p := parents[cur].(type) {
	case *ast.ValueSpec:
		if len(p.Names) == 1 && len(p.Values) == 1 && p.Values[0] == cur {
			itVar = varFromObj(info.Defs[p.Names[0]])
		}
	case *ast.AssignStmt:
		if len(p.Lhs) == 1 && len(p.Rhs) == 1 && p.Rhs[0] == cur {
			if lhsId, ok := stripParens(p.Lhs[0]).(*ast.Ident); ok {
				if p.Tok == token.DEFINE {
					itVar = varFromObj(info.Defs[lhsId])
				} else if p.Tok == token.ASSIGN {
					itVar = varFromObj(info.Uses[lhsId])
				}
			}
		}
	}
	if itVar != nil {
		fnBody := enclosingFuncBody(cur, parents)
		if fnBody == nil {
			return false
		}
		callsRemove := false
		ast.Inspect(fnBody, func(n ast.Node) bool {
			if callsRemove || n == nil {
				return false
			}
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) != 0 {
				return true
			}
			sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
			if !ok || sel.Sel == nil || sel.Sel.Name != "Remove" {
				return true
			}
			if id, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == itVar {
				callsRemove = true
				return false
			}
			return true
		})
		return !callsRemove
	}

	// 2. Check `return <recv>.<field>.Iterator()` on a struct method:
	//    Ensure the owner struct is not a mutable Collection/List/Set and no caller in ctx.Files calls `.Remove()` on its iterator.
	if ret, isRet := parents[cur].(*ast.ReturnStmt); isRet && slices.Contains(ret.Results, cur) && cand.isField && cand.ownerStruct != nil {
		ptrOwner := types.NewPointer(cand.ownerStruct)
		ms := types.NewMethodSet(ptrOwner)
		for i := 0; i < ms.Len(); i++ {
			mName := ms.At(i).Obj().Name()
			switch mName {
			case "Implements_j2gort_Ext_java_util_Collection",
				"Implements_j2gort_Ext_java_util_List",
				"Implements_j2gort_Ext_java_util_Set",
				"Implements_j2gort_Ext_java_util_Queue",
				"Implements_j2gort_Ext_java_util_Deque":
				return false
			}
		}
		enclFn := enclosingFuncObj(ret, parents, info)
		if enclFn == nil {
			return false
		}
		enclOrig := funcOrigin(enclFn)
		for _, f := range ctx.Files {
			callerRemoves := false
			ast.Inspect(f, func(n ast.Node) bool {
				if callerRemoves || n == nil {
					return false
				}
				fnBody, isBlk := n.(*ast.BlockStmt)
				if !isBlk {
					return true
				}
				var ownerIterVars []*types.Var
				for _, st := range fnBody.List {
					var lhsId *ast.Ident
					var rhsExpr ast.Expr
					switch s := st.(type) {
					case *ast.DeclStmt:
						if gd, ok := s.Decl.(*ast.GenDecl); ok && gd.Tok == token.VAR && len(gd.Specs) == 1 {
							if vs, ok := gd.Specs[0].(*ast.ValueSpec); ok && len(vs.Names) == 1 && len(vs.Values) == 1 {
								lhsId, rhsExpr = vs.Names[0], vs.Values[0]
							}
						}
					case *ast.AssignStmt:
						if len(s.Lhs) == 1 && len(s.Rhs) == 1 {
							if id, ok := stripParens(s.Lhs[0]).(*ast.Ident); ok {
								lhsId, rhsExpr = id, s.Rhs[0]
							}
						}
					}
					if lhsId != nil && rhsExpr != nil && exprCallsMethodOrigin(rhsExpr, enclOrig, cand.ownerStruct, info) {
						v := varFromObj(info.Defs[lhsId])
						if v == nil {
							v = varFromObj(info.Uses[lhsId])
						}
						if v != nil {
							ownerIterVars = append(ownerIterVars, v)
						}
					}
				}
				if len(ownerIterVars) > 0 {
					ast.Inspect(fnBody, func(m ast.Node) bool {
						if callerRemoves || m == nil {
							return false
						}
						c, ok := m.(*ast.CallExpr)
						if !ok || len(c.Args) != 0 {
							return true
						}
						sel, ok := stripParens(c.Fun).(*ast.SelectorExpr)
						if !ok || sel.Sel == nil || sel.Sel.Name != "Remove" {
							return true
						}
						if id, ok := stripParens(sel.X).(*ast.Ident); ok && slices.Contains(ownerIterVars, varFromObj(info.Uses[id])) {
							callerRemoves = true
							return false
						}
						return true
					})
				}
				return !callerRemoves
			})
			if callerRemoves {
				return false
			}
		}
		return true
	}

	return true
}

func exprCallsMethodOrigin(expr ast.Expr, targetFn *types.Func, ownerStruct *types.Named, info *types.Info) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if callee, _ := calleeFunc(call, info); callee != nil {
			if funcOrigin(callee) == targetFn {
				found = true
				return false
			}
			if sel, ok := stripParens(call.Fun).(*ast.SelectorExpr); ok && sel.Sel.Name == targetFn.Name() {
				if ptrNamedOrigin(info.Types[sel.X].Type) == ownerStruct {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}

type ordMapValSliceCand struct {
	v               *types.Var
	isField         bool
	ownerStruct     *types.Named
	declFile        *ast.File
	declTypeCarrier ast.Node // *ast.Field or *ast.ValueSpec
	defCount        int
	disqualified    bool
	matchedPutCalls map[*ast.CallExpr]bool
	readTypes       []types.Type
	writeTypes      []types.Type
	spellSites      []typeSpellSite
	j2gortSites     []typeSpellSite
	typeIndexRewr   []struct {
		file *ast.File
		ix   *ast.IndexListExpr
	}
	localVarSpecs []struct {
		file *ast.File
		vs   *ast.ValueSpec
	}
	sizeCalls []struct {
		file *ast.File
		call *ast.CallExpr
		id   *ast.Ident
	}
	isEmptyCalls []struct {
		file *ast.File
		call *ast.CallExpr
		id   *ast.Ident
	}
	getCalls []struct {
		file  *ast.File
		outer ast.Expr
		id    *ast.Ident
		idx   ast.Expr
	}
	getInitAdds []ordMapValSliceGetInitAdd
}

type ordMapValSliceGetInitAdd struct {
	file       *ast.File
	listId     *ast.Ident
	keyName    string
	ifStmt     *ast.IfStmt
	capExpr    ast.Expr // nil for 0-arg NewExt_java_util_ArrayList()
	addStmt    *ast.ExprStmt
	mapAccess  ast.Expr
	elemRawArg ast.Expr
	unwrapsBox bool
}

// applyOrderedMapValueSliceLowering lowers non-escaping OrderedMap / CustomOrderedMap / IdentityOrderedMap
// fields and locals whose value type is Ext_java_util_ArrayList / Ext_java_util_List (such as
// RemoveUnusedCode.removablesForPropertyNames) to native Go slice values OrderedMap[K, []T].
func applyOrderedMapValueSliceLowering(ctx *PassContext, parentsByFile map[*ast.File]map[ast.Node]ast.Node) int {
	info := ctx.Info
	if info == nil || len(ctx.Files) == 0 {
		return 0
	}

	isArrayListOrListType := func(t types.Type) bool {
		if t == nil {
			return false
		}
		named, ok := types.Unalias(t).(*types.Named)
		if !ok || named.Obj() == nil || !isJ2gortPkg(named.Obj().Pkg()) {
			return false
		}
		switch named.Obj().Name() {
		case "Ext_java_util_ArrayList", "Ext_java_util_List":
			return true
		}
		return false
	}

	isOrdMapListValType := func(t types.Type) bool {
		if t == nil {
			return false
		}
		named, ok := types.Unalias(t).(*types.Named)
		if !ok || named.Obj() == nil || !isJ2gortPkg(named.Obj().Pkg()) {
			return false
		}
		switch named.Obj().Name() {
		case "OrderedMap", "CustomOrderedMap", "IdentityOrderedMap":
		default:
			return false
		}
		targs := named.TypeArgs()
		return targs != nil && targs.Len() == 2 && isArrayListOrListType(targs.At(1))
	}

	isNewOrdMapListCall := func(expr ast.Expr) (*ast.IndexListExpr, bool) {
		call, ok := stripParens(expr).(*ast.CallExpr)
		if !ok || call.Ellipsis.IsValid() {
			return nil, false
		}
		fnName, _, _, isJ2 := isJ2gortCall(call, info)
		if !isJ2 {
			return nil, false
		}
		switch fnName {
		case "NewOrderedMap", "NewCustomOrderedMap", "NewIdentityOrderedMap":
			if len(call.Args) != 0 {
				return nil, false
			}
		case "NewOrderedMapCap", "NewCustomOrderedMapCap", "NewIdentityOrderedMapCap":
			if len(call.Args) != 1 || !isIntegerExpr(call.Args[0], info) {
				return nil, false
			}
		default:
			return nil, false
		}
		ile, ok := stripParens(call.Fun).(*ast.IndexListExpr)
		if !ok || len(ile.Indices) != 2 {
			return nil, false
		}
		return ile, true
	}

	cands := make(map[*types.Var]*ordMapValSliceCand)
	var allCands []*ordMapValSliceCand
	watchedStructs := make(map[*types.Named]bool)

	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) || ctx.ModifiedFiles[file] {
			continue
		}
		parents := parentsByFile[file]
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, s := range gd.Specs {
				ts, ok := s.(*ast.TypeSpec)
				if !ok || ts.Assign.IsValid() || (ts.TypeParams != nil && len(ts.TypeParams.List) > 0) {
					continue
				}
				st, stOK := ts.Type.(*ast.StructType)
				tn, tnOK := info.Defs[ts.Name].(*types.TypeName)
				if !stOK || !tnOK || st.Fields == nil {
					continue
				}
				named := namedOrigin(tn.Type())
				if named == nil || (named.TypeParams() != nil && named.TypeParams().Len() > 0) {
					continue
				}
				for _, fld := range st.Fields.List {
					if len(fld.Names) != 1 {
						continue
					}
					v := varFromObj(info.Defs[fld.Names[0]])
					if v == nil || !v.IsField() || v.Exported() || v.Embedded() || !isOrdMapListValType(v.Type()) {
						continue
					}
					c := &ordMapValSliceCand{
						v:               v,
						isField:         true,
						ownerStruct:     named,
						declFile:        file,
						declTypeCarrier: fld,
						matchedPutCalls: make(map[*ast.CallExpr]bool),
						spellSites:      []typeSpellSite{{file: file, pos: fld.Pos()}},
					}
					cands[v] = c
					allCands = append(allCands, c)
					watchedStructs[named] = true
				}
			}
		}

		ast.Inspect(file, func(n ast.Node) bool {
			switch s := n.(type) {
			case *ast.DeclStmt:
				gd, ok := s.Decl.(*ast.GenDecl)
				if !ok || gd.Tok != token.VAR || len(gd.Specs) != 1 || enclosingFuncBody(s, parents) == nil {
					return true
				}
				vs, ok := gd.Specs[0].(*ast.ValueSpec)
				if !ok || len(vs.Names) != 1 || len(vs.Values) != 1 {
					return true
				}
				v := varFromObj(info.Defs[vs.Names[0]])
				if v == nil || v.IsField() || vs.Names[0].Name == "_" || cands[v] != nil || !isOrdMapListValType(v.Type()) {
					return true
				}
				ile, initOK := isNewOrdMapListCall(vs.Values[0])
				if !initOK {
					return true
				}
				c := &ordMapValSliceCand{
					v:               v,
					declFile:        file,
					declTypeCarrier: vs,
					matchedPutCalls: make(map[*ast.CallExpr]bool),
					spellSites:      []typeSpellSite{{file: file, pos: vs.Pos()}},
				}
				c.typeIndexRewr = append(c.typeIndexRewr, struct {
					file *ast.File
					ix   *ast.IndexListExpr
				}{file: file, ix: ile})
				cands[v] = c
				allCands = append(allCands, c)

			case *ast.AssignStmt:
				if s.Tok != token.DEFINE || len(s.Lhs) != 1 || len(s.Rhs) != 1 || enclosingFuncBody(s, parents) == nil {
					return true
				}
				id, ok := s.Lhs[0].(*ast.Ident)
				if !ok || id.Name == "_" {
					return true
				}
				v := varFromObj(info.Defs[id])
				if v == nil || v.IsField() || cands[v] != nil || !isOrdMapListValType(v.Type()) {
					return true
				}
				ile, initOK := isNewOrdMapListCall(s.Rhs[0])
				if !initOK {
					return true
				}
				c := &ordMapValSliceCand{
					v:               v,
					declFile:        file,
					matchedPutCalls: make(map[*ast.CallExpr]bool),
					spellSites:      []typeSpellSite{{file: file, pos: s.Pos()}},
				}
				c.typeIndexRewr = append(c.typeIndexRewr, struct {
					file *ast.File
					ix   *ast.IndexListExpr
				}{file: file, ix: ile})
				cands[v] = c
				allCands = append(allCands, c)
			}
			return true
		})
	}

	if len(allCands) == 0 {
		return 0
	}

	if len(watchedStructs) > 0 {
		checkStructValueCopies(ctx, parentsByFile, watchedStructs, func(fv *types.Var) {
			if c := cands[fv]; c != nil {
				c.disqualified = true
			}
		}, func(named *types.Named) {
			for _, c := range allCands {
				if c.ownerStruct == named {
					c.disqualified = true
				}
			}
		})
	}

	type useSite struct {
		id   *ast.Ident
		file *ast.File
	}
	usesByVar := make(map[*types.Var][]useSite)
	for _, file := range ctx.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if c := cands[varFromObj(info.Defs[id])]; c != nil {
				c.defCount++
				if !ctx.ShouldTransform(file) || ctx.ModifiedFiles[file] {
					c.disqualified = true
				}
			}
			if v := varFromObj(info.Uses[id]); v != nil {
				usesByVar[v] = append(usesByVar[v], useSite{id: id, file: file})
			}
			return true
		})
	}

	matchLocalDecl := func(initExpr ast.Expr, parents map[ast.Node]ast.Node) (ast.Stmt, *ast.ValueSpec, *ast.Ident) {
		switch p := parents[initExpr].(type) {
		case *ast.ValueSpec:
			if len(p.Names) == 1 && len(p.Values) == 1 && p.Values[0] == initExpr {
				if gd, ok := parents[p].(*ast.GenDecl); ok && gd.Tok == token.VAR && len(gd.Specs) == 1 {
					if ds, ok := parents[gd].(*ast.DeclStmt); ok {
						return ds, p, p.Names[0]
					}
				}
			}
		case *ast.AssignStmt:
			if p.Tok == token.DEFINE && len(p.Lhs) == 1 && len(p.Rhs) == 1 && p.Rhs[0] == initExpr {
				if id, ok := p.Lhs[0].(*ast.Ident); ok {
					return p, nil, id
				}
			}
		}
		return nil, nil, nil
	}

	validateReadOnlySliceLocal := func(cand *ordMapValSliceCand, declFile *ast.File, vs *ast.ValueSpec, listId *ast.Ident) bool {
		if listId == nil || listId.Name == "_" {
			return false
		}
		lv := varFromObj(info.Defs[listId])
		if lv == nil {
			return false
		}
		if vs != nil && vs.Type != nil {
			cand.spellSites = append(cand.spellSites, typeSpellSite{file: declFile, pos: vs.Pos()})
			cand.localVarSpecs = append(cand.localVarSpecs, struct {
				file *ast.File
				vs   *ast.ValueSpec
			}{file: declFile, vs: vs})
		}
		for _, us := range usesByVar[lv] {
			if !ctx.ShouldTransform(us.file) || ctx.ModifiedFiles[us.file] {
				return false
			}
			up := parentsByFile[us.file]
			uid := us.id
			if as, ok := up[uid].(*ast.AssignStmt); ok && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && stripParens(as.Rhs[0]) == uid {
				if bId, ok := stripParens(as.Lhs[0]).(*ast.Ident); ok && bId.Name == "_" {
					continue
				}
			}
			top := unwrapParensUp(uid, up)
			if bin, ok := up[top].(*ast.BinaryExpr); ok && (bin.Op == token.EQL || bin.Op == token.NEQ) {
				other := bin.Y
				if bin.Y == top {
					other = bin.X
				}
				if isNilIdent(other, info) {
					continue
				}
			}
			topCheck := unwrapEnclosingCheckNotNulls(uid, up, info)
			sel, ok := up[topCheck].(*ast.SelectorExpr)
			if !ok || sel.X != topCheck {
				return false
			}
			mc, ok := up[sel].(*ast.CallExpr)
			if !ok || mc.Fun != sel || mc.Ellipsis.IsValid() {
				return false
			}
			_, inES := up[mc].(*ast.ExprStmt)
			switch sel.Sel.Name {
			case "Size":
				if inES || len(mc.Args) != 0 || !isUnshadowedBuiltin("len", mc.Pos(), ctx.Pkg) || !isUnshadowedUniverseType("int32", mc.Pos(), ctx.Pkg) {
					return false
				}
				cand.sizeCalls = append(cand.sizeCalls, struct {
					file *ast.File
					call *ast.CallExpr
					id   *ast.Ident
				}{file: us.file, call: mc, id: uid})
			case "IsEmpty":
				if inES || len(mc.Args) != 0 || !isUnshadowedBuiltin("len", mc.Pos(), ctx.Pkg) {
					return false
				}
				cand.isEmptyCalls = append(cand.isEmptyCalls, struct {
					file *ast.File
					call *ast.CallExpr
					id   *ast.Ident
				}{file: us.file, call: mc, id: uid})
			case "Get":
				if len(mc.Args) != 1 || !isIntegerExpr(mc.Args[0], info) {
					return false
				}
				outer, readT, ok := enclosingUnboxExpr(mc, up, info)
				_, outerES := up[outer].(*ast.ExprStmt)
				if !ok || outerES || isInvalidOrUntyped(readT) {
					return false
				}
				cand.readTypes = append(cand.readTypes, readT)
				cand.getCalls = append(cand.getCalls, struct {
					file  *ast.File
					outer ast.Expr
					id    *ast.Ident
					idx   ast.Expr
				}{file: us.file, outer: outer, id: uid, idx: mc.Args[0]})
			default:
				return false
			}
		}
		return true
	}

	tryMatchGetInitAdd := func(cand *ordMapValSliceCand, file *ast.File, mapAccess ast.Expr, getCall *ast.CallExpr, parents map[ast.Node]ast.Node) bool {
		if len(getCall.Args) != 1 {
			return false
		}
		keyId, ok := stripParens(getCall.Args[0]).(*ast.Ident)
		if !ok || keyId.Name == "_" {
			return false
		}
		keyVar := varFromObj(info.Uses[keyId])
		if keyVar == nil || keyVar.IsField() {
			return false
		}
		declStmt, vs, listId := matchLocalDecl(getCall, parents)
		if declStmt == nil || listId == nil || listId.Name == "_" {
			return false
		}
		listVar := varFromObj(info.Defs[listId])
		block, _ := parents[declStmt].(*ast.BlockStmt)
		if listVar == nil || block == nil {
			return false
		}
		idx := slices.Index(block.List, declStmt)
		if idx < 0 || idx+2 >= len(block.List) {
			return false
		}
		hasBlank := false
		ifIdx := idx + 1
		if asBlank, ok := block.List[ifIdx].(*ast.AssignStmt); ok && asBlank.Tok == token.ASSIGN && len(asBlank.Lhs) == 1 && len(asBlank.Rhs) == 1 {
			if bId, ok1 := stripParens(asBlank.Lhs[0]).(*ast.Ident); ok1 && bId.Name == "_" {
				if rId, ok2 := stripParens(asBlank.Rhs[0]).(*ast.Ident); ok2 && varFromObj(info.Uses[rId]) == listVar {
					hasBlank = true
					ifIdx++
				}
			}
		}
		addIdx := ifIdx + 1
		if addIdx >= len(block.List) {
			return false
		}
		ifStmt, ok := block.List[ifIdx].(*ast.IfStmt)
		if !ok || ifStmt.Init != nil || ifStmt.Else != nil || ifStmt.Body == nil || len(ifStmt.Body.List) != 2 {
			return false
		}
		// Match `if list == nil`
		bin, ok := stripParens(ifStmt.Cond).(*ast.BinaryExpr)
		if !ok || bin.Op != token.EQL {
			return false
		}
		condL, condR := stripParens(bin.X), stripParens(bin.Y)
		if isNilIdent(condL, info) {
			condL, condR = condR, condL
		}
		condId, ok := condL.(*ast.Ident)
		if !ok || varFromObj(info.Uses[condId]) != listVar || !isNilIdent(condR, info) {
			return false
		}
		// Body[0]: `list = NewExt_java_util_ArrayList()` or `NewExt_java_util_ArrayList_1_int32(cap)`
		initAs, ok := ifStmt.Body.List[0].(*ast.AssignStmt)
		if !ok || initAs.Tok != token.ASSIGN || len(initAs.Lhs) != 1 || len(initAs.Rhs) != 1 {
			return false
		}
		initLhs, ok := stripParens(initAs.Lhs[0]).(*ast.Ident)
		if !ok || varFromObj(info.Uses[initLhs]) != listVar {
			return false
		}
		allocCall, ok := stripParens(initAs.Rhs[0]).(*ast.CallExpr)
		if !ok || allocCall.Ellipsis.IsValid() {
			return false
		}
		allocFn, _, _, allocJ2 := isJ2gortCall(allocCall, info)
		if !allocJ2 {
			return false
		}
		var capExpr ast.Expr
		switch allocFn {
		case "NewExt_java_util_ArrayList":
			if len(allocCall.Args) != 0 {
				return false
			}
		case "NewExt_java_util_ArrayList_1_int32":
			if len(allocCall.Args) != 1 || !isIntegerExpr(allocCall.Args[0], info) || exprReferencesVar(allocCall.Args[0], keyVar, info) {
				return false
			}
			capExpr = allocCall.Args[0]
			if !isUnshadowedBuiltin("make", allocCall.Pos(), ctx.Pkg) {
				return false
			}
		default:
			return false
		}
		// Body[1]: `m.Put(keyId, list)`
		putEs, ok := ifStmt.Body.List[1].(*ast.ExprStmt)
		if !ok {
			return false
		}
		putCall, ok := stripParens(putEs.X).(*ast.CallExpr)
		if !ok || len(putCall.Args) != 2 || putCall.Ellipsis.IsValid() {
			return false
		}
		putSel, ok := stripParens(putCall.Fun).(*ast.SelectorExpr)
		if !ok || putSel.Sel.Name != "Put" || !sameAccessPath(putSel.X, mapAccess, cand.v, info) {
			return false
		}
		putKeyId, ok1 := stripParens(putCall.Args[0]).(*ast.Ident)
		putListId, ok2 := stripParens(putCall.Args[1]).(*ast.Ident)
		if !ok1 || !ok2 || varFromObj(info.Uses[putKeyId]) != keyVar || varFromObj(info.Uses[putListId]) != listVar {
			return false
		}
		// Next statement: `list.Add_1(elemExpr)`
		addEs, ok := block.List[addIdx].(*ast.ExprStmt)
		if !ok || !isUnshadowedBuiltin("append", addEs.Pos(), ctx.Pkg) {
			return false
		}
		addCall, ok := stripParens(addEs.X).(*ast.CallExpr)
		if !ok || len(addCall.Args) != 1 || addCall.Ellipsis.IsValid() {
			return false
		}
		addSel, ok := stripParens(addCall.Fun).(*ast.SelectorExpr)
		if !ok || addSel.Sel.Name != "Add_1" {
			return false
		}
		addRecvId, ok := stripParens(addSel.X).(*ast.Ident)
		if !ok || varFromObj(info.Uses[addRecvId]) != listVar {
			return false
		}
		elemArg := addCall.Args[0]
		if exprReferencesVar(elemArg, listVar, info) || exprReferencesVar(elemArg, cand.v, info) || hasClosureMutation(elemArg) {
			return false
		}
		unboxed, elemT, ok := unwrapBoxedAddArg(elemArg, info)
		if !ok || isInvalidOrUntyped(elemT) {
			return false
		}
		expectedUses := 4
		if hasBlank {
			expectedUses = 5
		}
		if len(usesByVar[listVar]) != expectedUses {
			return false
		}

		cand.matchedPutCalls[putCall] = true
		cand.writeTypes = append(cand.writeTypes, elemT)
		cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: ifStmt.Pos()})
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: addEs.Pos()})
		if vs != nil && vs.Type != nil {
			cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: vs.Pos()})
			cand.localVarSpecs = append(cand.localVarSpecs, struct {
				file *ast.File
				vs   *ast.ValueSpec
			}{file: file, vs: vs})
		}
		cand.getInitAdds = append(cand.getInitAdds, ordMapValSliceGetInitAdd{
			file:       file,
			listId:     listId,
			keyName:    keyId.Name,
			ifStmt:     ifStmt,
			capExpr:    capExpr,
			addStmt:    addEs,
			mapAccess:  mapAccess,
			elemRawArg: elemArg,
			unwrapsBox: unboxed != stripParens(elemArg),
		})
		return true
	}

	// Pre-pass over Get calls so matchedPutCalls is populated regardless of AST order.
	for _, cand := range allCands {
		if cand.disqualified {
			continue
		}
		for _, us := range usesByVar[cand.v] {
			if !ctx.ShouldTransform(us.file) || ctx.ModifiedFiles[us.file] {
				cand.disqualified = true
				break
			}
			up := parentsByFile[us.file]
			var accessExpr ast.Expr = us.id
			if cand.isField {
				sel, ok := up[us.id].(*ast.SelectorExpr)
				if !ok || sel.Sel != us.id || !isPureFieldSelector(sel, info) {
					cand.disqualified = true
					break
				}
				accessExpr = sel
			}
			if mSel, ok := up[accessExpr].(*ast.SelectorExpr); ok && mSel.X == accessExpr && mSel.Sel.Name == "Get" {
				if mCall, ok := up[mSel].(*ast.CallExpr); ok && mCall.Fun == mSel && !mCall.Ellipsis.IsValid() {
					_ = tryMatchGetInitAdd(cand, us.file, accessExpr, mCall, up)
				}
			}
		}
	}

	// Full validation pass over every use of every candidate map.
	for _, cand := range allCands {
		if cand.disqualified || cand.defCount != 1 {
			cand.disqualified = true
			continue
		}
		matchedGetCalls := make(map[*ast.CallExpr]bool, len(cand.getInitAdds))
		for _, gia := range cand.getInitAdds {
			up := parentsByFile[gia.file]
			if vs, ok := up[gia.listId].(*ast.ValueSpec); ok && len(vs.Values) == 1 {
				if gc, ok := vs.Values[0].(*ast.CallExpr); ok {
					matchedGetCalls[gc] = true
				}
			} else if as, ok := up[gia.listId].(*ast.AssignStmt); ok && len(as.Rhs) == 1 {
				if gc, ok := as.Rhs[0].(*ast.CallExpr); ok {
					matchedGetCalls[gc] = true
				}
			}
		}

		for _, us := range usesByVar[cand.v] {
			if cand.disqualified {
				break
			}
			up := parentsByFile[us.file]
			var accessExpr ast.Expr = us.id
			if cand.isField {
				sel, ok := up[us.id].(*ast.SelectorExpr)
				if !ok || sel.Sel != us.id || !isPureFieldSelector(sel, info) {
					cand.disqualified = true
					break
				}
				accessExpr = sel
			}
			if as, ok := up[accessExpr].(*ast.AssignStmt); ok {
				if as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 {
					if stripParens(as.Lhs[0]) == accessExpr {
						ile, initOK := isNewOrdMapListCall(as.Rhs[0])
						if !initOK {
							cand.disqualified = true
							break
						}
						cand.spellSites = append(cand.spellSites, typeSpellSite{file: us.file, pos: as.Pos()})
						cand.typeIndexRewr = append(cand.typeIndexRewr, struct {
							file *ast.File
							ix   *ast.IndexListExpr
						}{file: us.file, ix: ile})
						continue
					}
					if !cand.isField && stripParens(as.Rhs[0]) == us.id {
						if bId, ok := stripParens(as.Lhs[0]).(*ast.Ident); ok && bId.Name == "_" {
							continue
						}
					}
				}
			}

			mSel, ok := up[accessExpr].(*ast.SelectorExpr)
			if !ok || mSel.X != accessExpr {
				cand.disqualified = true
				break
			}
			mCall, ok := up[mSel].(*ast.CallExpr)
			if !ok || mCall.Fun != mSel || mCall.Ellipsis.IsValid() {
				cand.disqualified = true
				break
			}
			_, inES := up[mCall].(*ast.ExprStmt)

			switch mSel.Sel.Name {
			case "Size", "IsEmpty":
				if inES || len(mCall.Args) != 0 {
					cand.disqualified = true
				}
			case "ContainsKey":
				if inES || len(mCall.Args) != 1 {
					cand.disqualified = true
				}
			case "Clear":
				if !inES || len(mCall.Args) != 0 {
					cand.disqualified = true
				}
			case "Put":
				if !cand.matchedPutCalls[mCall] {
					cand.disqualified = true
				}
			case "Get":
				if matchedGetCalls[mCall] {
					continue
				}
				if inES || len(mCall.Args) != 1 {
					cand.disqualified = true
					break
				}
				ds, vs, listId := matchLocalDecl(mCall, up)
				if ds == nil || !validateReadOnlySliceLocal(cand, us.file, vs, listId) {
					cand.disqualified = true
				}
			case "Remove":
				if inES || len(mCall.Args) != 1 {
					cand.disqualified = true
					break
				}
				ds, vs, listId := matchLocalDecl(mCall, up)
				if ds == nil || !validateReadOnlySliceLocal(cand, us.file, vs, listId) {
					cand.disqualified = true
				}
			case "IterEntries":
				if inES || len(mCall.Args) != 0 {
					cand.disqualified = true
					break
				}
				ds, vs, curId := matchLocalDecl(mCall, up)
				if ds == nil || vs != nil || curId == nil || curId.Name == "_" {
					cand.disqualified = true
					break
				}
				curVar := varFromObj(info.Defs[curId])
				if curVar == nil {
					cand.disqualified = true
					break
				}
				for _, cus := range usesByVar[curVar] {
					if !ctx.ShouldTransform(cus.file) || ctx.ModifiedFiles[cus.file] {
						cand.disqualified = true
						break
					}
					cup := parentsByFile[cus.file]
					csel, ok := cup[cus.id].(*ast.SelectorExpr)
					if !ok || csel.X != cus.id {
						cand.disqualified = true
						break
					}
					ccall, ok := cup[csel].(*ast.CallExpr)
					if !ok || ccall.Fun != csel || len(ccall.Args) != 0 || ccall.Ellipsis.IsValid() {
						cand.disqualified = true
						break
					}
					switch csel.Sel.Name {
					case "NextEntry", "Key":
					case "Value":
						vds, vvs, listId := matchLocalDecl(ccall, cup)
						if vds == nil || !validateReadOnlySliceLocal(cand, cus.file, vvs, listId) {
							cand.disqualified = true
						}
					default:
						cand.disqualified = true
					}
					if cand.disqualified {
						break
					}
				}
			default:
				cand.disqualified = true
			}
		}
	}

	slices.SortFunc(allCands, func(a, b *ordMapValSliceCand) int { return cmp.Compare(a.v.Pos(), b.v.Pos()) })

	total := 0
	for _, cand := range allCands {
		if cand.disqualified || len(cand.getInitAdds) == 0 {
			continue
		}
		dummy := &sliceCandidate{readTypes: cand.readTypes, writeTypes: cand.writeTypes}
		elemT, ok := resolveGroupElemType([]*sliceCandidate{dummy})
		if !ok || isEmptyInterfaceType(elemT) {
			continue
		}
		if !canSpellTypeAtAllSites(elemT, cand.spellSites, ctx) || !canSpellJ2gortAtAllSites(cand.j2gortSites, ctx) {
			continue
		}

		makeSliceTypeFor := func(f *ast.File) *ast.ArrayType {
			elt, _ := typeExprForFile(elemT, f, ctx.Pkg, ctx.Importer)
			return &ast.ArrayType{Elt: elt}
		}
		makeElemTypeFor := func(f *ast.File) func() ast.Expr {
			return func() ast.Expr {
				elt, _ := typeExprForFile(elemT, f, ctx.Pkg, ctx.Importer)
				return elt
			}
		}
		j2SelFor := func(f *ast.File) func(string, token.Pos) *ast.SelectorExpr {
			return func(name string, pos token.Pos) *ast.SelectorExpr {
				q, _ := j2gortQualifierForFile(f, ctx)
				return &ast.SelectorExpr{X: &ast.Ident{NamePos: pos, Name: q}, Sel: &ast.Ident{NamePos: pos, Name: name}}
			}
		}

		// 1. Rewrite expression sites first (`list.Size()`, `list.IsEmpty()`, `unbox(list.Get(i))`).
		for _, sc := range cand.sizeCalls {
			up := parentsByFile[sc.file]
			pos := sc.call.Pos()
			repl := &ast.CallExpr{
				Fun: &ast.Ident{NamePos: pos, Name: "int32"},
				Args: []ast.Expr{
					&ast.CallExpr{
						Fun:  &ast.Ident{NamePos: pos, Name: "len"},
						Args: []ast.Expr{&ast.Ident{NamePos: sc.id.Pos(), Name: sc.id.Name}},
					},
				},
			}
			if replaceExprInParent(up[sc.call], sc.call, repl, up) {
				ctx.MarkFileModified(sc.file)
				total++
			}
		}
		for _, ec := range cand.isEmptyCalls {
			up := parentsByFile[ec.file]
			pos := ec.call.Pos()
			repl := &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun:  &ast.Ident{NamePos: pos, Name: "len"},
					Args: []ast.Expr{&ast.Ident{NamePos: ec.id.Pos(), Name: ec.id.Name}},
				},
				Op: token.EQL,
				Y:  &ast.BasicLit{Kind: token.INT, Value: "0"},
			}
			if replaceExprInParent(up[ec.call], ec.call, repl, up) {
				ctx.MarkFileModified(ec.file)
				total++
			}
		}
		for _, gc := range cand.getCalls {
			up := parentsByFile[gc.file]
			repl := &ast.IndexExpr{
				X:     &ast.Ident{NamePos: gc.id.Pos(), Name: gc.id.Name},
				Index: gc.idx,
			}
			if replaceExprInParent(up[gc.outer], gc.outer, repl, up) {
				ctx.MarkFileModified(gc.file)
				total++
			}
		}

		// 2. Rewrite `get-init-add` blocks (`if list == nil { list = make(...) }; list = append(list, elem); m.Put(k, list)`).
		for _, gia := range cand.getInitAdds {
			up := parentsByFile[gia.file]
			pos := gia.ifStmt.Pos()
			var initRhs ast.Expr
			if gia.capExpr != nil {
				initRhs = &ast.CallExpr{
					Fun: &ast.Ident{NamePos: pos, Name: "make"},
					Args: []ast.Expr{
						makeSliceTypeFor(gia.file),
						&ast.BasicLit{Kind: token.INT, Value: "0"},
						gia.capExpr,
					},
				}
			} else {
				initRhs = &ast.CompositeLit{Type: makeSliceTypeFor(gia.file)}
			}
			gia.ifStmt.Body.List = []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{NamePos: pos, Name: gia.listId.Name}},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{initRhs},
				},
			}
			unboxedElem := unwrapBoxArgExpr(gia.elemRawArg, gia.unwrapsBox)
			writeElem := wrapSliceIfaceWriteElem(unboxedElem, elemT, makeElemTypeFor(gia.file), j2SelFor(gia.file), info)
			addPos := gia.addStmt.Pos()
			appendAs := &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{NamePos: addPos, Name: gia.listId.Name}},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.Ident{NamePos: addPos, Name: "append"},
						Args: []ast.Expr{
							&ast.Ident{NamePos: addPos, Name: gia.listId.Name},
							writeElem,
						},
					},
				},
			}
			putStmt := &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   clonePureExpr(gia.mapAccess),
						Sel: &ast.Ident{NamePos: addPos, Name: "Put"},
					},
					Args: []ast.Expr{
						&ast.Ident{NamePos: addPos, Name: gia.keyName},
						&ast.Ident{NamePos: addPos, Name: gia.listId.Name},
					},
				},
			}
			if replaceStmtInParent(up[gia.addStmt], gia.addStmt, []ast.Stmt{appendAs, putStmt}, up) {
				ctx.MarkFileModified(gia.file)
				total++
			}
		}

		// 3. Rewrite map declaration type, `NewOrderedMap[K, V]` type arguments, and local slice ValueSpecs.
		switch d := cand.declTypeCarrier.(type) {
		case *ast.Field:
			if ile, ok := stripParens(d.Type).(*ast.IndexListExpr); ok && len(ile.Indices) == 2 {
				ile.Indices[1] = makeSliceTypeFor(cand.declFile)
				ctx.MarkFileModified(cand.declFile)
				total++
			}
		case *ast.ValueSpec:
			if d.Type != nil {
				if ile, ok := stripParens(d.Type).(*ast.IndexListExpr); ok && len(ile.Indices) == 2 {
					ile.Indices[1] = makeSliceTypeFor(cand.declFile)
					ctx.MarkFileModified(cand.declFile)
					total++
				}
			}
		}
		for _, tr := range cand.typeIndexRewr {
			tr.ix.Indices[1] = makeSliceTypeFor(tr.file)
			ctx.MarkFileModified(tr.file)
			total++
		}
		for _, lvs := range cand.localVarSpecs {
			lvs.vs.Type = makeSliceTypeFor(lvs.file)
			ctx.MarkFileModified(lvs.file)
			total++
		}
	}
	return total
}
