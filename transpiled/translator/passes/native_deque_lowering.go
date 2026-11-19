// Clean-room native_deque_lowering pass: lowers non-escaping ArrayDeque, Deque, and Queue
// local variables and unexported struct fields (initialized from NewExt_java_util_ArrayDeque*
// or ArrayDequeFrom*) across a package to compact ring-buffer j2gort.Deque[T] (or *j2gort.Deque[T]
// when nullable or pointer-aliased).

package passes

import (
	"cmp"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"slices"
	"strings"
)

// NativeDequeLoweringPass rewrites non-escaping ArrayDeque / Deque / Queue locals and unexported
// struct fields to j2gort.Deque[T] (32-byte ring-buffer value struct with 0 heap allocations when
// empty) or *j2gort.Deque[T] (when nullable or pointer-aliased).
type NativeDequeLoweringPass struct{}

func NewNativeDequeLoweringPass() *NativeDequeLoweringPass { return &NativeDequeLoweringPass{} }
func (p *NativeDequeLoweringPass) Name() string            { return "native_deque_lowering" }
func (p *NativeDequeLoweringPass) Description() string {
	return "Lower non-escaping ArrayDeque/Deque/Queue locals and unexported fields to j2gort.Deque[T] / *j2gort.Deque[T]"
}

func (p *NativeDequeLoweringPass) Apply(ctx *PassContext) (int, error) {
	if isJ2gortPkg(ctx.Pkg) || len(ctx.Files) == 0 {
		return 0, nil
	}
	if !j2gortHasExport(ctx, "NewDeque") {
		return 0, nil
	}
	parentsByFile := make(map[*ast.File]map[ast.Node]ast.Node, len(ctx.Files))
	for _, f := range ctx.Files {
		parentsByFile[f], _ = buildParentMap(f)
	}

	eligible := analyzeDequeCandidates(ctx, parentsByFile)
	if len(eligible) == 0 {
		return 0, nil
	}
	return applyDequeRewrites(ctx, eligible, parentsByFile), nil
}

type dequeInitKind uint8

const (
	dequeInitNil dequeInitKind = iota
	dequeInitEmpty
	dequeInitCap
	dequeInitCopyDeque
	dequeInitCopySlice
	dequeInitCopyOrdSet
	dequeInitCopyOrdSetPtr
	dequeInitCopyLocalOrdSet
	dequeInitCopyCustomOrdSet
	dequeInitCopyCustomOrdSetPtr
	dequeInitCopyIdentityOrdSet
	dequeInitCopyIdentityOrdSetPtr
	dequeInitCopyCol
	dequeInitPtrAlias
	dequeInitIIFE
)

type dequeOpKind uint8

const (
	dequeOpDeclField dequeOpKind = iota
	dequeOpDeclParam
	dequeOpDeclVarSpec
	dequeOpDeclShortVar
	dequeOpAssignInit
	dequeOpParamArg
	dequeOpNilCheck
	dequeOpFormatSink
	dequeOpPrecondMsg    // CheckStateMsg(false, d) -> CheckStateStr(false, d.String())
	dequeOpJoinStrings   // String.join(delim, d) -> DequeJoinStrings(&d, delim)
	dequeOpMethodDirect  // .Size(), .IsEmpty(), .Clear(), .String()
	dequeOpWriteVoid     // .Push(x), .AddFirst(x), .AddLast(x), DequePushNode(d, n)
	dequeOpWriteBool     // .Add_1(x), .Add(x), .Offer*(x), DequeAddLastNode(d, n)
	dequeOpReadThrow     // .Pop(), .Remove_0(), .RemoveFirst(), .RemoveLast(), .GetFirst(), .Element(), .GetLast(), DequePopNode, DequeRemoveLastNode
	dequeOpReadNullable  // .Peek*(), .Poll*(), DequePollFirstNode, DequePeekFirstNode, DequePeekLastNode
	dequeOpContainsOrRem // .Contains(x), .Remove_1_any(x), .RemoveFirstOccurrence(x), .RemoveLastOccurrence(x)
	dequeOpAddAll        // .AddAll_1(col), .AddAll(col)
	dequeOpIterator      // .Iterator(), .DescendingIterator()
	dequeOpIterCursor    // cur := NewIterCursor(d); for cur.HasNext() { ... }
	dequeOpRawUnbox      // unbox Ident replacement inside cursor loop body
)

type dequeIIFEBranch struct {
	ret       *ast.ReturnStmt
	spec      dequeInitSpec
	otherCand *dequeCandidate
}

type dequeInitSpec struct {
	kind      dequeInitKind
	rawExpr   ast.Expr
	argExpr   ast.Expr
	srcCall   *ast.CallExpr
	otherCand *dequeCandidate
	iifeFl    *ast.FuncLit
	iifeBrs   []dequeIIFEBranch
}

type dequeRewriteOp struct {
	kind        dequeOpKind
	cand        *dequeCandidate
	otherCand   *dequeCandidate
	initSpec    dequeInitSpec
	file        *ast.File
	spanNode    ast.Node
	accessExpr  ast.Expr
	argExpr     ast.Expr
	argElemType types.Type
	call        *ast.CallExpr
	outerExpr   ast.Expr
	stmt        ast.Stmt
	forStmt     *ast.ForStmt
	methodName  string
	elemIdent   *ast.Ident
	dropStmts   int
	unwrapsBox  bool
	isSliceArg  bool
	discardRet  bool
	isNegated   bool
}

type dequeCandidate struct {
	v               *types.Var
	isField         bool
	isParam         bool
	paramFn         *types.Func
	paramIdx        int
	inPkgCalls      int
	isPtr           bool
	forcePtr        bool
	ownerStruct     *types.Named
	recvTArgs       *types.TypeList
	declFile        *ast.File
	defCount        int
	hasInit         bool
	hasNilAssign    bool
	hasNilCheck     bool
	disqualified    bool
	elemType        types.Type
	writeTypes      []types.Type
	readTypes       []types.Type
	siteConstraints []siteTypeConstraint
	spellSites      []typeSpellSite
	j2gortSites     []typeSpellSite
	loopBodies      []*ast.BlockStmt
	ops             []dequeRewriteOp
	parent          *dequeCandidate
}

func (c *dequeCandidate) root() *dequeCandidate {
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

func unionDequeCandidates(a, b *dequeCandidate) {
	ra, rb := a.root(), b.root()
	if ra != rb {
		rb.parent = ra
	}
}

func (c *dequeCandidate) addWriteType(t types.Type, accessExpr ast.Expr, info *types.Info) bool {
	if t == nil || isInvalidOrUntyped(t) || isEmptyInterfaceType(t) {
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
	} else if !c.isField && c.ownerStruct != nil && c.recvTArgs != nil {
		lifted, ok := liftTypeViaTArgs(t, c.ownerStruct, c.recvTArgs)
		if !ok {
			c.disqualified = true
			return false
		}
		t = lifted
	}
	c.writeTypes = append(c.writeTypes, t)
	return true
}

func (c *dequeCandidate) addReadType(t types.Type, accessExpr ast.Expr, info *types.Info) bool {
	if t == nil || isInvalidOrUntyped(t) || isEmptyInterfaceType(t) {
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
	} else if !c.isField && c.ownerStruct != nil && c.recvTArgs != nil {
		lifted, ok := liftTypeViaTArgs(t, c.ownerStruct, c.recvTArgs)
		if !ok {
			c.disqualified = true
			return false
		}
		t = lifted
	}
	c.readTypes = append(c.readTypes, t)
	return true
}

func analyzeDequeCandidates(
	ctx *PassContext,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
) []*dequeCandidate {
	info := ctx.Info
	varUseCounts := make(map[*types.Var]int, len(info.Uses))
	// order-independent: commutative per-variable integer use-count accumulation
	for _, obj := range info.Uses {
		if v := varFromObj(obj); v != nil {
			varUseCounts[v]++
		}
	}

	varCands := make(map[*types.Var]*dequeCandidate)
	paramCandsByFn := make(map[*types.Func][]*dequeCandidate)
	var allCands []*dequeCandidate
	var structOrder []*types.Named
	structFields := make(map[*types.Named][]*dequeCandidate)
	var edges [][2]*dequeCandidate
	var ptrAliasEdges [][2]*dequeCandidate

	// Pass 1a: Collect unexported ArrayDeque / Deque / Queue struct fields and private helper parameters.
	for _, file := range ctx.Files {
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
						v := varFromObj(info.Defs[fld.Names[0]])
						if v == nil || !v.IsField() || v.Exported() || v.Embedded() || !isCandidateDequeType(v.Type()) {
							continue
						}
						cand := &dequeCandidate{
							v: v, isField: true, ownerStruct: named, declFile: file,
							spellSites:  []typeSpellSite{{file: file, pos: fld.Pos(), owner: named}},
							j2gortSites: []typeSpellSite{{file: file, pos: fld.Pos()}},
						}
						cand.ops = append(cand.ops, dequeRewriteOp{kind: dequeOpDeclField, cand: cand, file: file, spanNode: fld})
						varCands[v] = cand
						allCands = append(allCands, cand)
						if len(structFields[named]) == 0 {
							structOrder = append(structOrder, named)
						}
						structFields[named] = append(structFields[named], cand)
					}
				}
			case *ast.FuncDecl:
				if d.Body == nil || d.Type.Params == nil || (d.Type.TypeParams != nil && len(d.Type.TypeParams.List) > 0) {
					continue
				}
				fnObj, ok := info.Defs[d.Name].(*types.Func)
				if !ok || fnObj == nil || fnObj.Exported() {
					continue
				}
				sig, ok := fnObj.Type().(*types.Signature)
				if !ok || (sig.Recv() != nil && !strings.HasPrefix(fnObj.Name(), "priv_")) {
					continue
				}
				fnOrig := funcOrigin(fnObj)
				var recvOwner *types.Named
				var recvTArgs *types.TypeList
				if sig.Recv() != nil {
					recvOwner, recvTArgs = derefNamedAndTypeArgs(sig.Recv().Type())
				}
				argPos := 0
				for _, fld := range d.Type.Params.List {
					if len(fld.Names) == 1 && fld.Names[0].Name != "_" {
						if _, isEllipsis := fld.Type.(*ast.Ellipsis); !isEllipsis {
							if pv := varFromObj(info.Defs[fld.Names[0]]); pv != nil && isCandidateDequeType(pv.Type()) {
								pc := &dequeCandidate{
									v: pv, isParam: true, paramFn: fnOrig, paramIdx: argPos,
									isPtr: true, forcePtr: true, hasInit: true,
									ownerStruct: recvOwner, recvTArgs: recvTArgs, declFile: file,
									spellSites:  []typeSpellSite{{file: file, pos: fld.Pos(), owner: recvOwner, recvTArgs: recvTArgs}},
									j2gortSites: []typeSpellSite{{file: file, pos: fld.Pos()}},
								}
								pc.ops = append(pc.ops, dequeRewriteOp{kind: dequeOpDeclParam, cand: pc, file: file, spanNode: fld})
								varCands[pv] = pc
								paramCandsByFn[fnOrig] = append(paramCandsByFn[fnOrig], pc)
								allCands = append(allCands, pc)
							}
						}
					}
					if len(fld.Names) == 0 {
						argPos++
					} else {
						argPos += len(fld.Names)
					}
				}
			}
		}
	}

	// Pass 1b: Collect local ArrayDeque / Deque / Queue variables (two-phase so local declarations
	// can reference struct fields or earlier locals via CopyCol / PtrAlias).
	type pendingLocalInit struct {
		cand     *dequeCandidate
		initExpr ast.Expr
		file     *ast.File
		span     ast.Node
		opK      dequeOpKind
	}
	var pendingLocals []pendingLocalInit
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		collectLocal := func(id *ast.Ident, initExpr ast.Expr, span ast.Node, opK dequeOpKind) {
			v := varFromObj(info.Defs[id])
			if id.Name == "_" || v == nil || v.IsField() || !isCandidateDequeType(v.Type()) {
				return
			}
			if (initExpr == nil || isNilIdent(initExpr, info)) && opK != dequeOpDeclVarSpec {
				return
			}
			ownerStruct, recvTArgs := enclosingMethodOwnerAndTArgs(span, parents, info)
			cand := &dequeCandidate{
				v: v, declFile: file, ownerStruct: ownerStruct, recvTArgs: recvTArgs,
				spellSites:  []typeSpellSite{{file: file, pos: span.Pos(), owner: ownerStruct, recvTArgs: recvTArgs}},
				j2gortSites: []typeSpellSite{{file: file, pos: span.Pos()}},
			}
			varCands[v] = cand
			allCands = append(allCands, cand)
			pendingLocals = append(pendingLocals, pendingLocalInit{
				cand: cand, initExpr: initExpr, file: file, span: span, opK: opK,
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
						collectLocal(spec.Names[0], initExpr, spec, dequeOpDeclVarSpec)
					}
				}
			case *ast.AssignStmt:
				if s.Tok == token.DEFINE && len(s.Lhs) == 1 && len(s.Rhs) == 1 && enclosingFuncBody(s, parents) != nil {
					if id, ok := s.Lhs[0].(*ast.Ident); ok {
						collectLocal(id, s.Rhs[0], s, dequeOpDeclShortVar)
					}
				}
			}
			return true
		})
	}

	if len(allCands) == 0 {
		return nil
	}

	consumedExprs := make(map[ast.Expr]bool)
	recordConsumedInit := func(spec *dequeInitSpec) {
		var rec func(*dequeInitSpec)
		rec = func(s *dequeInitSpec) {
			switch s.kind {
			case dequeInitPtrAlias, dequeInitCopyDeque:
				if s.rawExpr != nil {
					consumedExprs[stripParens(s.rawExpr)] = true
				}
				if s.argExpr != nil {
					consumedExprs[stripParens(s.argExpr)] = true
				}
			case dequeInitIIFE:
				for i := range s.iifeBrs {
					rec(&s.iifeBrs[i].spec)
				}
			}
		}
		rec(spec)
	}

	for _, pl := range pendingLocals {
		if pl.initExpr == nil {
			pl.cand.hasNilAssign = true
			pl.cand.ops = append(pl.cand.ops, dequeRewriteOp{
				kind: pl.opK, cand: pl.cand, initSpec: dequeInitSpec{kind: dequeInitNil},
				file: pl.file, spanNode: pl.span,
			})
			continue
		}
		spec, mayNil, ok := classifyDequeInit(pl.initExpr, nil, pl.file, pl.cand, varCands, info, &edges, &ptrAliasEdges)
		if !ok {
			// Not initialized from ArrayDeque (e.g. LinkedList or PriorityQueue) — remove from varCands cleanly.
			pl.cand.disqualified = true
			delete(varCands, pl.cand.v)
			continue
		}
		if mayNil {
			pl.cand.hasNilAssign = true
		} else {
			pl.cand.hasInit = true
		}
		recordConsumedInit(&spec)
		pl.cand.ops = append(pl.cand.ops, dequeRewriteOp{
			kind: pl.opK, cand: pl.cand, initSpec: spec, otherCand: spec.otherCand,
			file: pl.file, spanNode: pl.span,
		})
	}

	if len(structOrder) > 0 {
		verifyDequeStructInvariants(ctx, parentsByFile, varCands, structOrder, structFields)
	}

	// Pass 2a: Pre-scan assignment RHSs `dst = ...` so candidate-to-candidate aliases and copies
	// are registered in `consumedExprs` before general consumer validation.
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		ast.Inspect(file, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			cand := varCands[varFromObj(info.Uses[id])]
			if cand == nil || cand.disqualified {
				return true
			}
			var accessExpr ast.Expr = id
			if cand.isField {
				sel, ok := parents[id].(*ast.SelectorExpr)
				if !ok || sel.Sel != id || !isIndirectStructFieldSel(sel, info) {
					return true
				}
				accessExpr = sel
			}
			if as, ok := parents[accessExpr].(*ast.AssignStmt); ok && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && stripParens(as.Lhs[0]) == accessExpr {
				if cand.isParam || (cand.isField && !isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info)) {
					cand.disqualified = true
					return true
				}
				spec, mayNil, ok := classifyDequeInit(as.Rhs[0], accessExpr, file, cand, varCands, info, &edges, &ptrAliasEdges)
				if !ok {
					cand.disqualified = true
					return true
				}
				if mayNil {
					cand.hasNilAssign = true
				}
				recordConsumedInit(&spec)
				cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
				cand.ops = append(cand.ops, dequeRewriteOp{
					kind: dequeOpAssignInit, cand: cand, initSpec: spec, otherCand: spec.otherCand,
					file: file, spanNode: as, stmt: as, accessExpr: accessExpr,
				})
				consumedExprs[accessExpr] = true
			}
			return true
		})
	}

	funcsMutating := make(map[*types.Var]map[*types.Func]bool)
	markMut := func(cand *dequeCandidate, node ast.Node, parents map[ast.Node]ast.Node) {
		if cand.v == nil {
			return
		}
		if fn := enclosingFuncObj(node, parents, info); fn != nil {
			if funcsMutating[cand.v] == nil {
				funcsMutating[cand.v] = make(map[*types.Func]bool)
			}
			funcsMutating[cand.v][fn] = true
		}
	}

	// Pass 2b: Validate all uses of every candidate and parameter-bearing private callees.
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		ast.Inspect(file, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if cand := varCands[varFromObj(info.Defs[id])]; cand != nil {
				cand.defCount++
				if !ctx.ShouldTransform(file) {
					cand.disqualified = true
				}
			}
			if fnObj, ok := info.Uses[id].(*types.Func); ok {
				if pcs := paramCandsByFn[funcOrigin(fnObj)]; len(pcs) > 0 {
					validateDequeParamFuncCall(id, pcs, varCands, ctx.ShouldTransform(file), parents, info)
				}
			}
			if cand := varCands[varFromObj(info.Uses[id])]; cand != nil && !cand.disqualified {
				if !ctx.ShouldTransform(file) {
					cand.disqualified = true
				} else {
					validateDequeUse(id, file, cand, varCands, paramCandsByFn, consumedExprs, ctx, parents, info, varUseCounts, markMut, &edges)
				}
			}
			return true
		})
	}

	// Finalize pointer vs value representation and propagate across pointer-alias edges.
	for _, cand := range allCands {
		if !cand.hasInit || cand.hasNilAssign || cand.hasNilCheck || cand.forcePtr {
			cand.isPtr = true
		}
	}
	for changed := true; changed; {
		changed = false
		for _, pe := range ptrAliasEdges {
			a, b := pe[0], pe[1]
			if a.isPtr != b.isPtr {
				a.isPtr = true
				b.isPtr = true
				changed = true
			}
		}
	}

	for _, e := range edges {
		unionDequeCandidates(e[0], e[1])
	}

	for _, cand := range allCands {
		if cand.disqualified || cand.defCount != 1 || !ctx.ShouldTransform(cand.declFile) {
			cand.disqualified = true
			continue
		}
		if cand.isParam && cand.inPkgCalls == 0 {
			cand.disqualified = true
			continue
		}
		if anyDequeLoopMutatesCandidate(cand, funcsMutating[cand.v], info) {
			cand.disqualified = true
			continue
		}
	}

	var groupRoots []*dequeCandidate
	groups := make(map[*dequeCandidate][]*dequeCandidate)
	for _, cand := range allCands {
		r := cand.root()
		if len(groups[r]) == 0 {
			groupRoots = append(groupRoots, r)
		}
		groups[r] = append(groups[r], cand)
	}

	var out []*dequeCandidate
	for _, r := range groupRoots {
		members := groups[r]
		if !resolveAndValidateDequeGroup(members, ctx) {
			continue
		}
		out = append(out, members...)
	}

	slices.SortFunc(out, func(a, b *dequeCandidate) int { return cmp.Compare(a.v.Pos(), b.v.Pos()) })
	return out
}

func resolveAndValidateDequeGroup(members []*dequeCandidate, ctx *PassContext) bool {
	var reads, writes []types.Type
	for _, m := range members {
		if m.disqualified {
			return false
		}
		writes = append(writes, m.writeTypes...)
		reads = append(reads, m.readTypes...)
	}
	var all []types.Type
	all = append(all, writes...)
	all = append(all, reads...)
	if len(all) == 0 {
		for _, m := range members {
			for _, sc := range m.siteConstraints {
				all = append(all, sc.siteType)
			}
		}
	}
	if len(all) == 0 {
		return false
	}

	elemT, ok := resolveImmElemType(all, reads, writes)
	if !ok || elemT == nil || isEmptyInterfaceType(elemT) || isInvalidOrUntyped(elemT) {
		return false
	}

	for _, m := range members {
		if typeContainsTypeParam(elemT) && (m.isField || m.isParam) && !typeOnlyUsesOwnerTypeParams(elemT, m.ownerStruct) {
			return false
		}
		for _, sc := range m.siteConstraints {
			inst, ok := instantiateOwnerType(elemT, m.ownerStruct, sc.targs)
			if !ok || (!types.Identical(inst, sc.siteType) && !(isNonEmptyInterfaceType(inst) && types.AssignableTo(sc.siteType, inst))) {
				return false
			}
		}
		if !canSpellTypeAtAllSites(elemT, m.spellSites, ctx) || !canSpellJ2gortAtAllSites(m.j2gortSites, ctx) {
			return false
		}
	}

	for _, m := range members {
		m.elemType = elemT
	}
	return true
}

func verifyDequeStructInvariants(
	ctx *PassContext,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
	varCands map[*types.Var]*dequeCandidate,
	structOrder []*types.Named,
	structFields map[*types.Named][]*dequeCandidate,
) {
	info := ctx.Info
	disqualify := func(named *types.Named) {
		for _, c := range structFields[named] {
			c.disqualified = true
		}
	}

	methodInits := make(map[*types.Func]map[*types.Var]bool)
	ctorFuncs := make(map[*types.Named][]*ast.FuncDecl)

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
				isInitFields := fn.Name.Name == "initFields_"+named.Obj().Name() || fn.Name.Name == "initFields" || strings.HasPrefix(fn.Name.Name, "initFields_")
				for _, stmt := range fn.Body.List {
					as, ok := stmt.(*ast.AssignStmt)
					if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
						continue
					}
					if sel, ok := stripParens(as.Lhs[0]).(*ast.SelectorExpr); ok {
						if recvId, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[recvId]) == recvVar {
							if cand := varCands[varFromObj(info.Uses[sel.Sel])]; cand != nil && cand.isField && isNonNilDequeInitExpr(as.Rhs[0], info) {
								if isInitFields {
									cand.hasInit = true
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
			}
		}
	}

	for _, named := range structOrder {
		ctors := ctorFuncs[named]
		if len(ctors) == 0 {
			continue
		}
		for _, cand := range structFields[named] {
			if cand.hasInit {
				continue
			}
			allInit := true
			for _, fn := range ctors {
				recvVar := varFromObj(info.Defs[fn.Type.Params.List[0].Names[0]])
				found := false
				for _, stmt := range fn.Body.List {
					switch s := stmt.(type) {
					case *ast.AssignStmt:
						if s.Tok == token.ASSIGN && len(s.Lhs) == 1 && len(s.Rhs) == 1 {
							if sel, ok := stripParens(s.Lhs[0]).(*ast.SelectorExpr); ok {
								if id, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == recvVar && varFromObj(info.Uses[sel.Sel]) == cand.v && isNonNilDequeInitExpr(s.Rhs[0], info) {
									found = true
								}
							}
						}
					case *ast.ExprStmt:
						if call, ok := stripParens(s.X).(*ast.CallExpr); ok {
							if sel, ok := stripParens(call.Fun).(*ast.SelectorExpr); ok {
								if id, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == recvVar {
									if callee, _ := calleeFunc(call, info); callee != nil && methodInits[funcOrigin(callee)][cand.v] {
										found = true
									}
								}
							} else if callee, _ := calleeFunc(call, info); callee != nil && len(call.Args) >= 1 {
								if id0, ok := stripParens(call.Args[0]).(*ast.Ident); ok && varFromObj(info.Uses[id0]) == recvVar {
									if strings.HasPrefix(callee.Name(), "Init"+named.Obj().Name()) {
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
			if allInit {
				cand.hasInit = true
			}
		}
	}

	watched := make(map[*types.Named]bool, len(structOrder))
	for _, named := range structOrder {
		watched[named] = true
	}
	checkStructValueCopies(ctx, parentsByFile, watched, func(fv *types.Var) {
		if c := varCands[fv]; c != nil {
			c.disqualified = true
		}
	}, disqualify)
}

func isNonNilDequeInitExpr(expr ast.Expr, info *types.Info) bool {
	e := unwrapDequeAdaptersDown(expr, info)
	if isNilIdent(e, info) {
		return false
	}
	call, ok := e.(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return false
	}
	if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
		switch fnName {
		case "NewExt_java_util_ArrayDeque", "NewExt_java_util_ArrayDeque_0",
			"NewExt_java_util_ArrayDeque_1_int32", "NewExt_java_util_ArrayDeque_1_Ext_java_util_Collection",
			"ArrayDequeFromOrderedSet", "ArrayDequeFromOrderedSetPtr", "ArrayDequeFromLocalOrderedSet",
			"ArrayDequeFromCustomOrderedSet", "ArrayDequeFromCustomOrderedSetPtr",
			"ArrayDequeFromIdentityOrderedSet", "ArrayDequeFromIdentityOrderedSetPtr":
			return true
		}
	}
	if len(call.Args) == 0 {
		if fl, ok := stripParens(call.Fun).(*ast.FuncLit); ok && fl.Body != nil {
			allOK := true
			hasRet := false
			ast.Inspect(fl.Body, func(n ast.Node) bool {
				if !allOK || n == nil {
					return false
				}
				if nested, isLit := n.(*ast.FuncLit); isLit && nested != fl {
					return false
				}
				if ret, isRet := n.(*ast.ReturnStmt); isRet {
					hasRet = true
					if len(ret.Results) != 1 || !isNonNilDequeInitExpr(ret.Results[0], info) {
						allOK = false
					}
					return false
				}
				return true
			})
			return hasRet && allOK
		}
	}
	return false
}

func classifyDequeInit(
	expr ast.Expr,
	accessExpr ast.Expr,
	file *ast.File,
	cand *dequeCandidate,
	varCands map[*types.Var]*dequeCandidate,
	info *types.Info,
	edges *[][2]*dequeCandidate,
	ptrAliasEdges *[][2]*dequeCandidate,
) (dequeInitSpec, bool, bool) {
	e := stripParens(expr)
	if isNilIdent(e, info) {
		return dequeInitSpec{kind: dequeInitNil, rawExpr: expr}, true, true
	}

	unwrapped := unwrapDequeAdaptersDown(e, info)

	// Pointer alias (`oldBlockStack = this_.blockStack` or `stack1 = resultStack`).
	if srcCand := resolveDequeCandAccess(unwrapped, varCands, info); srcCand != nil && !srcCand.disqualified {
		cand.forcePtr = true
		srcCand.forcePtr = true
		if srcCand != cand {
			*edges = append(*edges, [2]*dequeCandidate{cand, srcCand})
			*ptrAliasEdges = append(*ptrAliasEdges, [2]*dequeCandidate{cand, srcCand})
		}
		return dequeInitSpec{kind: dequeInitPtrAlias, rawExpr: expr, argExpr: unwrapped, otherCand: srcCand}, srcCand.hasNilAssign, true
	}

	// Ternary / multi-branch IIFE.
	if call, ok := e.(*ast.CallExpr); ok && len(call.Args) == 0 && !call.Ellipsis.IsValid() {
		if fl, ok := stripParens(call.Fun).(*ast.FuncLit); ok && fl.Body != nil && (fl.Type.Params == nil || len(fl.Type.Params.List) == 0) && fl.Type.Results != nil && len(fl.Type.Results.List) == 1 {
			var brs []dequeIIFEBranch
			mayNil := false
			valid := true
			ast.Inspect(fl.Body, func(n ast.Node) bool {
				if !valid || n == nil {
					return false
				}
				if nested, isLit := n.(*ast.FuncLit); isLit && nested != fl {
					return false
				}
				ret, isRet := n.(*ast.ReturnStmt)
				if !isRet {
					return true
				}
				if len(ret.Results) != 1 {
					valid = false
					return false
				}
				subSpec, subNil, subOK := classifyDequeInit(ret.Results[0], accessExpr, file, cand, varCands, info, edges, ptrAliasEdges)
				if !subOK {
					valid = false
					return false
				}
				if subNil {
					mayNil = true
				}
				brs = append(brs, dequeIIFEBranch{ret: ret, spec: subSpec, otherCand: subSpec.otherCand})
				return false
			})
			if valid && len(brs) > 0 {
				cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: call.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: call.Pos()})
				return dequeInitSpec{kind: dequeInitIIFE, rawExpr: expr, iifeFl: fl, iifeBrs: brs}, mayNil, true
			}
		}
	}

	call, ok := unwrapped.(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return dequeInitSpec{}, false, false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 {
		return dequeInitSpec{}, false, false
	}

	switch fnName {
	case "NewExt_java_util_ArrayDeque", "NewExt_java_util_ArrayDeque_0":
		if len(call.Args) == 0 {
			cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: call.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: call.Pos()})
			return dequeInitSpec{kind: dequeInitEmpty, rawExpr: expr}, false, true
		}

	case "NewExt_java_util_ArrayDeque_1_int32":
		if len(call.Args) == 1 && isIntegerExpr(call.Args[0], info) {
			cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: call.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: call.Pos()})
			return dequeInitSpec{kind: dequeInitCap, rawExpr: expr, argExpr: call.Args[0], srcCall: call}, false, true
		}

	case "NewExt_java_util_ArrayDeque_1_Ext_java_util_Collection":
		if len(call.Args) == 1 {
			colInner := unwrapDequeAdaptersDown(call.Args[0], info)
			if srcCand := resolveDequeCandAccess(colInner, varCands, info); srcCand != nil && srcCand != cand && !srcCand.disqualified {
				*edges = append(*edges, [2]*dequeCandidate{cand, srcCand})
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: call.Pos()})
				return dequeInitSpec{kind: dequeInitCopyDeque, rawExpr: expr, argExpr: colInner, otherCand: srcCand}, false, true
			}
			if sliceArg, elemT, okS := matchWrappedSliceArg(colInner, info); okS {
				if !cand.addWriteType(elemT, accessExpr, info) {
					return dequeInitSpec{}, false, false
				}
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: call.Pos()})
				return dequeInitSpec{kind: dequeInitCopySlice, rawExpr: expr, argExpr: sliceArg}, false, true
			}
			cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: call.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: call.Pos()})
			return dequeInitSpec{kind: dequeInitCopyCol, rawExpr: expr, argExpr: call.Args[0], srcCall: call}, false, true
		}

	case "ArrayDequeFromOrderedSet", "ArrayDequeFromOrderedSetPtr", "ArrayDequeFromLocalOrderedSet",
		"ArrayDequeFromCustomOrderedSet", "ArrayDequeFromCustomOrderedSetPtr",
		"ArrayDequeFromIdentityOrderedSet", "ArrayDequeFromIdentityOrderedSetPtr":
		if len(call.Args) == 1 {
			elemT, okT := extractFirstGenericTypeArg(info.Types[call.Args[0]].Type)
			if !okT || !cand.addWriteType(elemT, accessExpr, info) {
				return dequeInitSpec{}, false, false
			}
			initK := dequeInitCopyOrdSet
			switch fnName {
			case "ArrayDequeFromOrderedSetPtr":
				initK = dequeInitCopyOrdSetPtr
			case "ArrayDequeFromLocalOrderedSet":
				initK = dequeInitCopyLocalOrdSet
			case "ArrayDequeFromCustomOrderedSet":
				initK = dequeInitCopyCustomOrdSet
			case "ArrayDequeFromCustomOrderedSetPtr":
				initK = dequeInitCopyCustomOrdSetPtr
			case "ArrayDequeFromIdentityOrderedSet":
				initK = dequeInitCopyIdentityOrdSet
			case "ArrayDequeFromIdentityOrderedSetPtr":
				initK = dequeInitCopyIdentityOrdSetPtr
			}
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: call.Pos()})
			return dequeInitSpec{kind: initK, rawExpr: expr, argExpr: call.Args[0], srcCall: call}, false, true
		}
	}

	return dequeInitSpec{}, false, false
}

func validateDequeParamFuncCall(
	id *ast.Ident,
	pcs []*dequeCandidate,
	varCands map[*types.Var]*dequeCandidate,
	canTransform bool,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) {
	disqualifyAll := func() {
		for _, pc := range pcs {
			pc.disqualified = true
		}
	}
	if !canTransform {
		disqualifyAll()
		return
	}
	var cur ast.Node = id
	if sel, ok := parents[cur].(*ast.SelectorExpr); ok && sel.Sel == cur {
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
	call, ok := parents[cur].(*ast.CallExpr)
	if !ok || call.Fun != cur || call.Ellipsis.IsValid() {
		disqualifyAll()
		return
	}
	for _, pc := range pcs {
		if pc.disqualified {
			continue
		}
		if pc.paramIdx >= len(call.Args) {
			pc.disqualified = true
			continue
		}
		arg := unwrapDequeAdaptersDown(call.Args[pc.paramIdx], info)
		srcCand := resolveDequeCandAccess(arg, varCands, info)
		if srcCand == nil || srcCand.disqualified {
			pc.disqualified = true
			continue
		}
		pc.inPkgCalls++
	}
}

func isConstFalseExpr(expr ast.Expr, info *types.Info) bool {
	if tv, ok := info.Types[expr]; ok && tv.Value != nil && tv.Value.Kind() == constant.Bool {
		return !constant.BoolVal(tv.Value)
	}
	return false
}

func validateDequeUse(
	id *ast.Ident,
	file *ast.File,
	cand *dequeCandidate,
	varCands map[*types.Var]*dequeCandidate,
	paramCandsByFn map[*types.Func][]*dequeCandidate,
	consumedExprs map[ast.Expr]bool,
	ctx *PassContext,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
	markMut func(*dequeCandidate, ast.Node, map[ast.Node]ast.Node),
	edges *[][2]*dequeCandidate,
) {
	var accessExpr ast.Expr = id
	if cand.isField {
		sel, ok := parents[id].(*ast.SelectorExpr)
		if !ok || sel.Sel != id || !isIndirectStructFieldSel(sel, info) {
			cand.disqualified = true
			return
		}
		accessExpr = sel
	}

	if consumedExprs[accessExpr] {
		return
	}

	if as, ok := parents[accessExpr].(*ast.AssignStmt); ok {
		for _, lhs := range as.Lhs {
			if stripParens(lhs) == accessExpr {
				cand.disqualified = true
				return
			}
		}
	}
	if !cand.isField && isBlankAssignTarget(id, parents) {
		return
	}

	pTop := unwrapParensUp(accessExpr, parents)
	if un, ok := parents[pTop].(*ast.UnaryExpr); ok && un.Op == token.AND && un.X == pTop {
		cand.disqualified = true
		return
	}
	if ta, ok := parents[pTop].(*ast.TypeAssertExpr); ok && ta.X == pTop {
		cand.disqualified = true
		return
	}

	// Equality / nil checks: `d == nil` / `d != nil` are preserved natively via `*j2gort.Deque[T]`.
	if bin, ok := parents[pTop].(*ast.BinaryExpr); ok && (bin.Op == token.EQL || bin.Op == token.NEQ) {
		var otherSide ast.Expr
		if bin.X == pTop {
			otherSide = bin.Y
		} else if bin.Y == pTop {
			otherSide = bin.X
		}
		if otherSide != nil {
			if isNilIdent(otherSide, info) {
				cand.hasNilCheck = true
				cand.ops = append(cand.ops, dequeRewriteOp{
					kind: dequeOpNilCheck, cand: cand, file: file, spanNode: bin, outerExpr: bin, accessExpr: accessExpr,
				})
				return
			}
			cand.disqualified = true
			return
		}
	}

	// Check parameter-forwarding calls and static j2gort calls (`NewIterCursor`, `Deque*Node`, `StrVal`, `CheckStateMsg`, `String.join`).
	adaptedCol := unwrapDequeAdaptersUp(accessExpr, parents, info)
	if sCall, ok := parents[adaptedCol].(*ast.CallExpr); ok && !sCall.Ellipsis.IsValid() && sCall.Fun != adaptedCol && len(sCall.Args) >= 1 {
		if callee, _ := calleeFunc(sCall, info); callee != nil {
			if pcs := paramCandsByFn[funcOrigin(callee)]; len(pcs) > 0 {
				for _, pc := range pcs {
					if pc.paramIdx < len(sCall.Args) && sCall.Args[pc.paramIdx] == adaptedCol {
						addressable := true
						if cand.isField {
							addressable = isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info)
						}
						if !addressable {
							cand.disqualified = true
							pc.disqualified = true
							return
						}
						if pc.ownerStruct != nil {
							callSel, selOK := stripParens(sCall.Fun).(*ast.SelectorExpr)
							if !selOK {
								cand.disqualified = true
								pc.disqualified = true
								return
							}
							callOwner, callTArgs := derefNamedAndTypeArgs(info.Types[callSel.X].Type)
							var candTArgs *types.TypeList
							if cand.isField {
								candTArgs = ownerTypeArgsAtAccess(accessExpr, cand.ownerStruct, info)
							} else {
								candTArgs = cand.recvTArgs
							}
							if callOwner != pc.ownerStruct || cand.ownerStruct != pc.ownerStruct || !sameTypeList(callTArgs, candTArgs) {
								cand.disqualified = true
								pc.disqualified = true
								return
							}
						}
						*edges = append(*edges, [2]*dequeCandidate{cand, pc})
						cand.ops = append(cand.ops, dequeRewriteOp{
							kind: dequeOpParamArg, cand: cand, otherCand: pc, file: file, spanNode: adaptedCol,
							outerExpr: adaptedCol, accessExpr: accessExpr,
						})
						markMut(cand, sCall, parents)
						return
					}
				}
			}
		}

		if fnName, _, _, isJ2 := isJ2gortCall(sCall, info); isJ2 {
			if sCall.Args[0] == adaptedCol {
				switch fnName {
				case "NewIterCursor", "NewTypedIterCursor":
					if len(sCall.Args) == 1 && unwrapEnclosingCheckNotNulls(accessExpr, parents, info) == adaptedCol {
						if validateDequeIterCursorLoop(sCall, accessExpr, file, cand, parents, info, varUseCounts) {
							return
						}
					}
				case "DequePushNode":
					if len(sCall.Args) == 2 {
						es, isES := parents[sCall].(*ast.ExprStmt)
						nodeArg, nodeT, okN := unwrapDequeNodeWriteArg(sCall.Args[1], info)
						if isES && okN && canReplaceStmtIn(parents[es]) && cand.addWriteType(nodeT, accessExpr, info) {
							markMut(cand, sCall, parents)
							cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
							cand.ops = append(cand.ops, dequeRewriteOp{
								kind: dequeOpWriteVoid, cand: cand, file: file, spanNode: es, stmt: es,
								accessExpr: accessExpr, argExpr: nodeArg, call: sCall, methodName: "Push",
							})
							return
						}
					}
				case "DequeAddLastNode":
					if len(sCall.Args) == 2 {
						nodeArg, nodeT, okN := unwrapDequeNodeWriteArg(sCall.Args[1], info)
						if okN && cand.addWriteType(nodeT, accessExpr, info) {
							markMut(cand, sCall, parents)
							cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
							if es, isES := parents[sCall].(*ast.ExprStmt); isES && canReplaceStmtIn(parents[es]) {
								cand.ops = append(cand.ops, dequeRewriteOp{
									kind: dequeOpWriteVoid, cand: cand, file: file, spanNode: es, stmt: es,
									accessExpr: accessExpr, argExpr: nodeArg, call: sCall, methodName: "AddLast",
								})
								return
							}
							cand.ops = append(cand.ops, dequeRewriteOp{
								kind: dequeOpWriteBool, cand: cand, file: file, spanNode: sCall, outerExpr: sCall,
								accessExpr: accessExpr, argExpr: nodeArg, call: sCall, methodName: "OfferLast",
							})
							return
						}
					}
				case "DequePopNode", "DequeRemoveLastNode":
					if len(sCall.Args) == 1 {
						markMut(cand, sCall, parents)
						mName := "Pop"
						if fnName == "DequeRemoveLastNode" {
							mName = "RemoveLast"
						}
						if _, isES := parents[sCall].(*ast.ExprStmt); isES {
							cand.ops = append(cand.ops, dequeRewriteOp{
								kind: dequeOpReadThrow, cand: cand, file: file, spanNode: sCall, outerExpr: sCall,
								accessExpr: accessExpr, methodName: mName, discardRet: true,
							})
							return
						}
						outer, nodeT := unwrapEnclosingNodeConvUp(sCall, parents, info)
						if _, outerES := parents[outer].(*ast.ExprStmt); !outerES && cand.addReadType(nodeT, accessExpr, info) {
							cand.ops = append(cand.ops, dequeRewriteOp{
								kind: dequeOpReadThrow, cand: cand, file: file, spanNode: outer, outerExpr: outer,
								accessExpr: accessExpr, methodName: mName,
							})
							return
						}
					}
				case "DequePollFirstNode", "DequePeekFirstNode", "DequePeekLastNode":
					if len(sCall.Args) == 1 {
						mName := "PeekFirst"
						if fnName == "DequePollFirstNode" {
							markMut(cand, sCall, parents)
							mName = "PollFirst"
						} else if fnName == "DequePeekLastNode" {
							mName = "PeekLast"
						}
						if _, isES := parents[sCall].(*ast.ExprStmt); isES {
							cand.ops = append(cand.ops, dequeRewriteOp{
								kind: dequeOpReadNullable, cand: cand, file: file, spanNode: sCall, outerExpr: sCall,
								accessExpr: accessExpr, methodName: mName, discardRet: true,
							})
							return
						}
						outer, nodeT := unwrapEnclosingNodeConvUp(sCall, parents, info)
						if _, outerES := parents[outer].(*ast.ExprStmt); !outerES && cand.addReadType(nodeT, accessExpr, info) {
							cand.ops = append(cand.ops, dequeRewriteOp{
								kind: dequeOpReadNullable, cand: cand, file: file, spanNode: outer, outerExpr: outer,
								accessExpr: accessExpr, methodName: mName,
							})
							return
						}
					}
				case "StrVal":
					if len(sCall.Args) == 1 {
						if _, isES := parents[sCall].(*ast.ExprStmt); !isES {
							cand.ops = append(cand.ops, dequeRewriteOp{
								kind: dequeOpFormatSink, cand: cand, file: file, spanNode: sCall, outerExpr: sCall, accessExpr: accessExpr,
							})
							return
						}
					}
				}
			} else if len(sCall.Args) == 2 && sCall.Args[1] == adaptedCol {
				switch fnName {
				case "CheckStateMsg", "CheckArgumentMsg":
					strHelper := "CheckStateStr"
					if fnName == "CheckArgumentMsg" {
						strHelper = "CheckArgumentStr"
					}
					if isConstFalseExpr(sCall.Args[0], info) {
						cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
						cand.ops = append(cand.ops, dequeRewriteOp{
							kind: dequeOpPrecondMsg, cand: cand, file: file, spanNode: sCall, outerExpr: sCall,
							accessExpr: accessExpr, call: sCall, methodName: strHelper,
						})
						return
					}
					if es, isES := parents[sCall].(*ast.ExprStmt); isES && canReplaceStmtIn(parents[es]) {
						cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
						cand.ops = append(cand.ops, dequeRewriteOp{
							kind: dequeOpPrecondMsg, cand: cand, file: file, spanNode: es, stmt: es,
							accessExpr: accessExpr, call: sCall, methodName: strHelper,
						})
						return
					}
				case "Ext_java_lang_String_Join_2_any_Ext_java_lang_Iterable":
					if _, isES := parents[sCall].(*ast.ExprStmt); !isES && j2gortHasExport(ctx, "DequeJoinStrings") && cand.addReadType(types.Typ[types.String], accessExpr, info) {
						cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
						cand.ops = append(cand.ops, dequeRewriteOp{
							kind: dequeOpJoinStrings, cand: cand, file: file, spanNode: sCall, outerExpr: sCall,
							accessExpr: accessExpr, argExpr: sCall.Args[0], call: sCall,
						})
						return
					}
				}
			}
		}
	}

	topExpr := unwrapEnclosingCheckNotNulls(accessExpr, parents, info)
	mSel, ok := parents[topExpr].(*ast.SelectorExpr)
	if !ok || mSel.X != topExpr {
		cand.disqualified = true
		return
	}
	mCall, ok := parents[mSel].(*ast.CallExpr)
	if !ok || mCall.Fun != mSel || mCall.Ellipsis.IsValid() {
		cand.disqualified = true
		return
	}
	es, inES := parents[mCall].(*ast.ExprStmt)

	switch mSel.Sel.Name {
	case "Size", "IsEmpty", "String":
		if inES || len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		cand.ops = append(cand.ops, dequeRewriteOp{
			kind: dequeOpMethodDirect, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, methodName: mSel.Sel.Name,
		})
		return

	case "Clear":
		if !inES || len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		markMut(cand, mCall, parents)
		cand.ops = append(cand.ops, dequeRewriteOp{
			kind: dequeOpMethodDirect, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, methodName: "Clear",
		})
		return

	case "Push", "AddFirst", "AddLast":
		if !inES || len(mCall.Args) != 1 || !canReplaceStmtIn(parents[es]) {
			cand.disqualified = true
			return
		}
		ub, elemT, okE := unwrapBoxedAddArg(mCall.Args[0], info)
		if !okE || !cand.addWriteType(elemT, accessExpr, info) {
			cand.disqualified = true
			return
		}
		markMut(cand, mCall, parents)
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, dequeRewriteOp{
			kind: dequeOpWriteVoid, cand: cand, file: file, spanNode: es, stmt: es,
			accessExpr: accessExpr, call: mCall, methodName: mSel.Sel.Name, unwrapsBox: ub != stripParens(mCall.Args[0]),
		})
		return

	case "Add_1", "Add", "Offer", "OfferFirst", "OfferLast":
		if len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		ub, elemT, okE := unwrapBoxedAddArg(mCall.Args[0], info)
		if !okE || !cand.addWriteType(elemT, accessExpr, info) {
			cand.disqualified = true
			return
		}
		markMut(cand, mCall, parents)
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		unw := ub != stripParens(mCall.Args[0])
		if inES && canReplaceStmtIn(parents[es]) {
			voidMethod := "AddLast"
			if mSel.Sel.Name == "OfferFirst" {
				voidMethod = "AddFirst"
			}
			cand.ops = append(cand.ops, dequeRewriteOp{
				kind: dequeOpWriteVoid, cand: cand, file: file, spanNode: es, stmt: es,
				accessExpr: accessExpr, call: mCall, methodName: voidMethod, unwrapsBox: unw,
			})
			return
		}
		boolMethod := mSel.Sel.Name
		if boolMethod == "Add_1" {
			boolMethod = "Add"
		}
		cand.ops = append(cand.ops, dequeRewriteOp{
			kind: dequeOpWriteBool, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, call: mCall, methodName: boolMethod, unwrapsBox: unw,
		})
		return

	case "Pop", "Remove_0", "RemoveFirst", "RemoveLast", "GetFirst", "Element", "GetLast":
		if len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		mName := mSel.Sel.Name
		if mName == "Remove_0" {
			mName = "RemoveFirst"
		}
		if mName == "Pop" || mName == "RemoveFirst" || mName == "RemoveLast" {
			markMut(cand, mCall, parents)
		}
		if inES {
			cand.ops = append(cand.ops, dequeRewriteOp{
				kind: dequeOpReadThrow, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
				accessExpr: accessExpr, methodName: mName, discardRet: true,
			})
			return
		}
		outer, elemT, okU := enclosingUnboxExpr(mCall, parents, info)
		if !okU {
			cand.disqualified = true
			return
		}
		if _, outerES := parents[outer].(*ast.ExprStmt); outerES || !cand.addReadType(elemT, accessExpr, info) {
			cand.disqualified = true
			return
		}
		cand.ops = append(cand.ops, dequeRewriteOp{
			kind: dequeOpReadThrow, cand: cand, file: file, spanNode: outer, outerExpr: outer,
			accessExpr: accessExpr, methodName: mName,
		})
		return

	case "Peek", "Peek_0", "PeekFirst", "PeekLast", "Poll", "Poll_0", "PollFirst", "PollLast":
		if len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		mName := mSel.Sel.Name
		switch mName {
		case "Peek_0":
			mName = "Peek"
		case "Poll_0":
			mName = "Poll"
		}
		if strings.HasPrefix(mName, "Poll") {
			markMut(cand, mCall, parents)
		}
		if inES {
			cand.ops = append(cand.ops, dequeRewriteOp{
				kind: dequeOpReadNullable, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
				accessExpr: accessExpr, methodName: mName, discardRet: true,
			})
			return
		}
		outer, elemT, okU := enclosingUnboxExpr(mCall, parents, info)
		if !okU {
			cand.disqualified = true
			return
		}
		if _, outerES := parents[outer].(*ast.ExprStmt); outerES || !cand.addReadType(elemT, accessExpr, info) {
			cand.disqualified = true
			return
		}
		cand.ops = append(cand.ops, dequeRewriteOp{
			kind: dequeOpReadNullable, cand: cand, file: file, spanNode: outer, outerExpr: outer,
			accessExpr: accessExpr, methodName: mName,
		})
		return

	case "Contains", "Remove_1_any", "RemoveFirstOccurrence", "RemoveLastOccurrence":
		if len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		if mSel.Sel.Name != "Contains" {
			markMut(cand, mCall, parents)
		}
		mName := mSel.Sel.Name
		if mName == "Remove_1_any" {
			mName = "RemoveFirstOccurrence"
		}
		var argElemT types.Type
		unw := false
		if ub, et, okE := unwrapBoxedAddArg(mCall.Args[0], info); okE && !isNilIdent(mCall.Args[0], info) {
			argElemT = et
			unw = ub != stripParens(mCall.Args[0])
		}
		cand.ops = append(cand.ops, dequeRewriteOp{
			kind: dequeOpContainsOrRem, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, call: mCall, methodName: mName, argElemType: argElemT, unwrapsBox: unw,
		})
		return

	case "AddAll_1", "AddAll":
		if len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		markMut(cand, mCall, parents)
		colInner := unwrapDequeAdaptersDown(mCall.Args[0], info)
		var otherC *dequeCandidate
		var argExpr ast.Expr = mCall.Args[0]
		isSliceArg := false
		if sc := resolveDequeCandAccess(colInner, varCands, info); sc != nil && !sc.disqualified {
			otherC = sc
			argExpr = colInner
			if sc != cand {
				*edges = append(*edges, [2]*dequeCandidate{cand, sc})
			}
		} else if sliceArg, elemT, okS := matchWrappedSliceArg(colInner, info); okS {
			if !cand.addWriteType(elemT, accessExpr, info) {
				cand.disqualified = true
				return
			}
			argExpr = sliceArg
			isSliceArg = true
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, dequeRewriteOp{
			kind: dequeOpAddAll, cand: cand, otherCand: otherC, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, argExpr: argExpr, call: mCall, isSliceArg: isSliceArg,
		})
		return

	case "Iterator", "DescendingIterator":
		if inES || len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		outer := unwrapIteratorCast(mCall, parents, info)
		if mSel.Sel.Name == "Iterator" {
			topIterCheck := unwrapEnclosingCheckNotNulls(outer, parents, info)
			if cursorCall, isCall := parents[topIterCheck].(*ast.CallExpr); isCall && len(cursorCall.Args) == 1 && !cursorCall.Ellipsis.IsValid() && cursorCall.Args[0] == topIterCheck {
				if fnName, _, _, isJ2 := isJ2gortCall(cursorCall, info); isJ2 && (fnName == "NewIterCursor" || fnName == "NewTypedIterCursor") {
					if validateDequeIterCursorLoop(cursorCall, accessExpr, file, cand, parents, info, varUseCounts) {
						return
					}
				}
			}
		}
		if _, outerES := parents[outer].(*ast.ExprStmt); outerES {
			cand.disqualified = true
			return
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outer.Pos()})
		cand.ops = append(cand.ops, dequeRewriteOp{
			kind: dequeOpIterator, cand: cand, file: file, spanNode: outer, outerExpr: outer,
			accessExpr: accessExpr, methodName: mSel.Sel.Name,
		})
		return
	}

	cand.disqualified = true
}

func validateDequeIterCursorLoop(
	cursorCall *ast.CallExpr,
	accessExpr ast.Expr,
	file *ast.File,
	cand *dequeCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
) bool {
	initAs, ok := parents[cursorCall].(*ast.AssignStmt)
	if !ok || initAs.Tok != token.DEFINE || len(initAs.Lhs) != 1 || len(initAs.Rhs) != 1 || initAs.Rhs[0] != cursorCall {
		return false
	}
	curId, _ := initAs.Lhs[0].(*ast.Ident)
	curVar := varFromObj(info.Defs[curId])
	block, _ := parents[initAs].(*ast.BlockStmt)
	if curId == nil || curId.Name == "_" || curVar == nil || varUseCounts[curVar] != 2 || block == nil {
		return false
	}
	idx := slices.Index(block.List, ast.Stmt(initAs))
	if idx < 0 || idx+1 >= len(block.List) {
		return false
	}
	forStmt, ok := block.List[idx+1].(*ast.ForStmt)
	if !ok || forStmt.Init != nil || forStmt.Post != nil || forStmt.Cond == nil || forStmt.Body == nil || len(forStmt.Body.List) == 0 || !isCursorMethodCall(forStmt.Cond, curVar, "HasNext", info) {
		return false
	}
	s0, ok := forStmt.Body.List[0].(*ast.AssignStmt)
	if !ok || s0.Tok != token.DEFINE || len(s0.Lhs) != 1 || len(s0.Rhs) != 1 {
		return false
	}
	lhs0, _ := s0.Lhs[0].(*ast.Ident)
	if lhs0 == nil || lhs0.Name == "_" || info.Defs[lhs0] == nil {
		return false
	}

	elemId := lhs0
	drop := 1
	rhs0 := stripParens(s0.Rhs[0])
	switch {
	case isCursorMethodCall(rhs0, curVar, "NextStr", info):
		if !cand.addReadType(types.Typ[types.String], accessExpr, info) {
			return false
		}
	case isCursorMethodCall(rhs0, curVar, "NextNode", info):
		if !cand.addReadType(types.Typ[types.Uint32], accessExpr, info) {
			return false
		}
	case isNodeConvOfCursorNextNode(rhs0, curVar, info):
		if !cand.addReadType(info.Types[rhs0].Type, accessExpr, info) {
			return false
		}
	case isCursorMethodCall(rhs0, curVar, "Next", info):
		if !isEmptyInterfaceType(info.Types[rhs0].Type) {
			if !cand.addReadType(info.Types[rhs0].Type, accessExpr, info) {
				return false
			}
		} else {
			rawVar := varFromObj(info.Defs[lhs0])
			if rawVar == nil {
				return false
			}
			unboxedId, extraDrop, okL := validateDequeLoopRawElem(rawVar, accessExpr, file, forStmt.Body.List[1:], cand, parents, info, varUseCounts)
			if !okL {
				return false
			}
			if unboxedId != nil {
				elemId = unboxedId
				drop += extraDrop
			}
		}
	default:
		return false
	}

	cand.loopBodies = append(cand.loopBodies, forStmt.Body)
	cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: forStmt.Pos()})
	cand.ops = append(cand.ops, dequeRewriteOp{
		kind: dequeOpIterCursor, cand: cand, file: file, spanNode: forStmt,
		stmt: initAs, forStmt: forStmt, outerExpr: cursorCall, accessExpr: accessExpr,
		elemIdent: elemId, dropStmts: drop,
	})
	return true
}

func validateDequeLoopRawElem(
	rawVar *types.Var,
	accessExpr ast.Expr,
	file *ast.File,
	stmts []ast.Stmt,
	cand *dequeCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
) (*ast.Ident, int, bool) {
	if len(stmts) > 0 && varUseCounts[rawVar] == 1 {
		if unboxedId, elemT, ok := matchSingleDefineUnbox(stmts[0], rawVar, parents, info); ok {
			if !cand.addReadType(elemT, accessExpr, info) {
				return nil, 0, false
			}
			return unboxedId, 1, true
		}
	}
	if varUseCounts[rawVar] == 0 {
		return nil, 0, true
	}
	var ops []dequeRewriteOp
	valid := true
	for _, st := range stmts {
		ast.Inspect(st, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !valid || !ok || varFromObj(info.Uses[id]) != rawVar {
				return valid
			}
			if isBlankAssignTarget(id, parents) {
				return true
			}
			outer, expectedT, okU := enclosingUnboxExpr(id, parents, info)
			_, isES := parents[outer].(*ast.ExprStmt)
			if !okU || isES || !cand.addReadType(expectedT, accessExpr, info) {
				valid = false
				return false
			}
			if outer != id {
				ops = append(ops, dequeRewriteOp{kind: dequeOpRawUnbox, cand: cand, file: file, spanNode: outer, outerExpr: outer, elemIdent: id})
			}
			return true
		})
	}
	if !valid {
		return nil, 0, false
	}
	cand.ops = append(cand.ops, ops...)
	return nil, 0, true
}

func applyDequeRewrites(
	ctx *PassContext,
	eligible []*dequeCandidate,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
) int {
	var exprOps, stmtOps []dequeRewriteOp
	for _, cand := range eligible {
		for _, op := range cand.ops {
			switch op.kind {
			case dequeOpDeclField, dequeOpDeclParam, dequeOpDeclVarSpec, dequeOpDeclShortVar, dequeOpAssignInit,
				dequeOpWriteVoid, dequeOpIterCursor:
				stmtOps = append(stmtOps, op)
			case dequeOpPrecondMsg:
				if op.stmt != nil {
					stmtOps = append(stmtOps, op)
				} else {
					exprOps = append(exprOps, op)
				}
			default:
				exprOps = append(exprOps, op)
			}
		}
	}

	spanOrder := func(a, b dequeRewriteOp) int {
		if c := cmp.Compare(a.spanNode.End()-a.spanNode.Pos(), b.spanNode.End()-b.spanNode.Pos()); c != 0 {
			return c
		}
		return cmp.Compare(a.spanNode.Pos(), b.spanNode.Pos())
	}
	slices.SortFunc(exprOps, spanOrder)
	slices.SortFunc(stmtOps, spanOrder)

	count := 0
	for _, op := range exprOps {
		if applySingleDequeRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}
	for _, op := range stmtOps {
		if applySingleDequeRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}
	return count
}

func applySingleDequeRewrite(ctx *PassContext, op dequeRewriteOp, parents map[ast.Node]ast.Node) bool {
	j2Sel := func(name string, pos token.Pos) *ast.SelectorExpr {
		q, _ := j2gortQualifierForFile(op.file, ctx)
		return &ast.SelectorExpr{X: &ast.Ident{NamePos: pos, Name: q}, Sel: &ast.Ident{NamePos: pos, Name: name}}
	}
	siteElemT := op.cand.elemType
	if op.cand.ownerStruct != nil && typeContainsTypeParam(siteElemT) {
		if op.cand.recvTArgs != nil {
			if inst, ok := instantiateOwnerType(siteElemT, op.cand.ownerStruct, op.cand.recvTArgs); ok {
				siteElemT = inst
			}
		} else if op.accessExpr != nil {
			siteElemT = instantiateOwnerTypeAtExpr(siteElemT, op.cand.ownerStruct, op.accessExpr, ctx.Info)
		}
	}
	makeElemType := func() ast.Expr {
		et, _ := typeExprForFile(siteElemT, op.file, ctx.Pkg, ctx.Importer)
		return et
	}
	makeDequeValType := func(pos token.Pos) ast.Expr {
		return &ast.IndexExpr{X: j2Sel("Deque", pos), Lbrack: pos, Index: makeElemType(), Rbrack: pos}
	}
	makeDequeDeclType := func(pos token.Pos) ast.Expr {
		vt := makeDequeValType(pos)
		if op.cand.isPtr {
			return &ast.StarExpr{Star: pos, X: vt}
		}
		return vt
	}

	wrapWriteElem := func(rawArg ast.Expr, unwraps bool, pos token.Pos) ast.Expr {
		ub := unwrapBoxArgExpr(rawArg, unwraps)
		if isNamedUint32NodeType(siteElemT) && !isGuaranteedNonZeroNode(ub, nil, ctx.Info) {
			return &ast.CallExpr{Fun: j2Sel("CheckNotNullNode", pos), Args: []ast.Expr{ub}}
		}
		return ub
	}

	var buildInitExpr func(s dequeInitSpec, inHeader bool) ast.Expr
	buildInitExpr = func(s dequeInitSpec, inHeader bool) ast.Expr {
		pos := token.NoPos
		if s.rawExpr != nil {
			pos = s.rawExpr.Pos()
		}
		isPtr := op.cand.isPtr
		switch s.kind {
		case dequeInitNil:
			return &ast.Ident{NamePos: pos, Name: "nil"}
		case dequeInitEmpty:
			lit := &ast.CompositeLit{Type: makeDequeValType(pos)}
			if isPtr {
				return &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: lit}
			}
			if inHeader {
				return &ast.ParenExpr{X: lit}
			}
			return lit
		case dequeInitCap:
			capArg := s.argExpr
			if s.srcCall != nil && len(s.srcCall.Args) == 1 {
				capArg = s.srcCall.Args[0]
			}
			fnName := "NewDequeCap"
			if isPtr {
				fnName = "NewDequePtrCap"
			}
			return &ast.CallExpr{
				Fun:  &ast.IndexExpr{X: j2Sel(fnName, pos), Index: makeElemType()},
				Args: []ast.Expr{capArg},
			}
		case dequeInitCopyDeque:
			srcExpr := clonePureExpr(s.argExpr)
			if s.otherCand != nil && s.otherCand.isPtr {
				if isPtr {
					return &ast.CallExpr{Fun: j2Sel("DequePtrFromDeque", pos), Args: []ast.Expr{srcExpr}}
				}
				checked := &ast.CallExpr{Fun: j2Sel("CheckNotNullPtr", pos), Args: []ast.Expr{srcExpr}}
				return &ast.CallExpr{Fun: j2Sel("DequeFromDeque", pos), Args: []ast.Expr{&ast.StarExpr{Star: pos, X: checked}}}
			}
			if isPtr {
				return &ast.CallExpr{Fun: j2Sel("DequePtrFromDeque", pos), Args: []ast.Expr{&ast.UnaryExpr{OpPos: pos, Op: token.AND, X: srcExpr}}}
			}
			return &ast.CallExpr{Fun: j2Sel("DequeFromDeque", pos), Args: []ast.Expr{srcExpr}}
		case dequeInitCopySlice:
			fnName := "DequeFromSlice"
			if isPtr {
				fnName = "DequePtrFromSlice"
			}
			return &ast.CallExpr{Fun: j2Sel(fnName, pos), Args: []ast.Expr{s.argExpr}}
		case dequeInitCopyOrdSet, dequeInitCopyOrdSetPtr, dequeInitCopyLocalOrdSet,
			dequeInitCopyCustomOrdSet, dequeInitCopyCustomOrdSetPtr,
			dequeInitCopyIdentityOrdSet, dequeInitCopyIdentityOrdSetPtr:
			arg := s.argExpr
			if s.srcCall != nil && len(s.srcCall.Args) == 1 {
				arg = s.srcCall.Args[0]
			}
			fnName := "DequeFromOrderedSet"
			switch {
			case s.kind == dequeInitCopyOrdSet && isPtr:
				fnName = "DequePtrFromOrderedSet"
			case s.kind == dequeInitCopyOrdSetPtr && !isPtr:
				fnName = "DequeFromOrderedSetPtr"
			case s.kind == dequeInitCopyOrdSetPtr && isPtr:
				fnName = "DequePtrFromOrderedSetPtr"
			case s.kind == dequeInitCopyLocalOrdSet && !isPtr:
				fnName = "DequeFromLocalOrderedSet"
			case s.kind == dequeInitCopyLocalOrdSet && isPtr:
				fnName = "DequePtrFromLocalOrderedSet"
			case s.kind == dequeInitCopyCustomOrdSet && !isPtr:
				fnName = "DequeFromCustomOrderedSet"
			case s.kind == dequeInitCopyCustomOrdSet && isPtr:
				fnName = "DequePtrFromCustomOrderedSet"
			case s.kind == dequeInitCopyCustomOrdSetPtr && !isPtr:
				fnName = "DequeFromCustomOrderedSetPtr"
			case s.kind == dequeInitCopyCustomOrdSetPtr && isPtr:
				fnName = "DequePtrFromCustomOrderedSetPtr"
			case s.kind == dequeInitCopyIdentityOrdSet && !isPtr:
				fnName = "DequeFromIdentityOrderedSet"
			case s.kind == dequeInitCopyIdentityOrdSet && isPtr:
				fnName = "DequePtrFromIdentityOrderedSet"
			case s.kind == dequeInitCopyIdentityOrdSetPtr && !isPtr:
				fnName = "DequeFromIdentityOrderedSetPtr"
			case s.kind == dequeInitCopyIdentityOrdSetPtr && isPtr:
				fnName = "DequePtrFromIdentityOrderedSetPtr"
			}
			return &ast.CallExpr{Fun: j2Sel(fnName, pos), Args: []ast.Expr{arg}}
		case dequeInitCopyCol:
			colArg := s.argExpr
			if s.srcCall != nil && len(s.srcCall.Args) == 1 {
				colArg = s.srcCall.Args[0]
			}
			fnName := "DequeFromCol"
			if isPtr {
				fnName = "DequePtrFromCol"
			}
			return &ast.CallExpr{
				Fun:  &ast.IndexExpr{X: j2Sel(fnName, pos), Index: makeElemType()},
				Args: []ast.Expr{colArg},
			}
		case dequeInitPtrAlias:
			return clonePureExpr(s.argExpr)
		case dequeInitIIFE:
			s.iifeFl.Type.Results.List[0].Type = makeDequeDeclType(s.iifeFl.Type.Results.Pos())
			for _, br := range s.iifeBrs {
				br.ret.Results[0] = buildInitExpr(br.spec, false)
			}
			return s.rawExpr
		}
		return s.rawExpr
	}

	dequePtrExpr := func(pos token.Pos) ast.Expr {
		if op.cand.isPtr {
			return clonePureExpr(op.accessExpr)
		}
		return &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: clonePureExpr(op.accessExpr)}
	}

	switch op.kind {
	case dequeOpDeclField, dequeOpDeclParam:
		fld := op.spanNode.(*ast.Field)
		fld.Type = makeDequeDeclType(fld.Pos())
		return true

	case dequeOpParamArg:
		parent := parents[op.outerExpr]
		return parent != nil && replaceExprInParent(parent, op.outerExpr, dequePtrExpr(op.outerExpr.Pos()), parents)

	case dequeOpDeclVarSpec:
		spec := op.spanNode.(*ast.ValueSpec)
		spec.Type = makeDequeDeclType(spec.Pos())
		if len(spec.Values) == 1 {
			if op.initSpec.kind == dequeInitNil || (op.initSpec.kind == dequeInitEmpty && !op.cand.isPtr) {
				spec.Values = nil
			} else {
				spec.Values[0] = buildInitExpr(op.initSpec, false)
			}
		}
		return true

	case dequeOpDeclShortVar:
		as := op.spanNode.(*ast.AssignStmt)
		as.Rhs[0] = buildInitExpr(op.initSpec, inControlHeader(as, parents))
		return true

	case dequeOpAssignInit:
		if op.initSpec.kind == dequeInitNil {
			return false
		}
		as := op.spanNode.(*ast.AssignStmt)
		as.Rhs[0] = buildInitExpr(op.initSpec, inControlHeader(as, parents))
		return true

	case dequeOpNilCheck:
		return false

	case dequeOpRawUnbox:
		parent := parents[op.outerExpr]
		return parent != nil && replaceExprInParent(parent, op.outerExpr, &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name}, parents)

	case dequeOpFormatSink:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: "String"}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case dequeOpPrecondMsg:
		pos := op.call.Pos()
		strExpr := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: "String"}},
		}
		checkCall := &ast.CallExpr{
			Fun:  j2Sel(op.methodName, pos),
			Args: []ast.Expr{&ast.Ident{NamePos: pos, Name: "false"}, strExpr},
		}
		if op.stmt != nil {
			parent := parents[op.stmt]
			ifStmt := &ast.IfStmt{
				Cond: &ast.UnaryExpr{Op: token.NOT, X: &ast.ParenExpr{X: op.call.Args[0]}},
				Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{X: checkCall}}},
			}
			return replaceStmtInParent(parent, op.stmt, []ast.Stmt{ifStmt}, parents)
		}
		parent := parents[op.outerExpr]
		return parent != nil && replaceExprInParent(parent, op.outerExpr, checkCall, parents)

	case dequeOpJoinStrings:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{
			Fun:  j2Sel("DequeJoinStrings", pos),
			Args: []ast.Expr{dequePtrExpr(pos), op.argExpr},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case dequeOpMethodDirect:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: op.methodName}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case dequeOpWriteVoid:
		es := op.stmt
		parent := parents[es]
		pos := op.call.Pos()
		var elemArg ast.Expr
		if op.argExpr != nil {
			elemArg = wrapWriteElem(op.argExpr, false, pos)
		} else {
			elemArg = wrapWriteElem(op.call.Args[0], op.unwrapsBox, pos)
		}
		newCall := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: op.methodName}},
			Args: []ast.Expr{elemArg},
		}
		return replaceStmtInParent(parent, es, []ast.Stmt{&ast.ExprStmt{X: newCall}}, parents)

	case dequeOpWriteBool:
		parent := parents[op.outerExpr]
		pos := op.call.Pos()
		var elemArg ast.Expr
		if op.argExpr != nil {
			elemArg = wrapWriteElem(op.argExpr, false, pos)
		} else {
			elemArg = wrapWriteElem(op.call.Args[0], op.unwrapsBox, pos)
		}
		repl := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: op.methodName}},
			Args: []ast.Expr{elemArg},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case dequeOpReadThrow:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: op.methodName}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case dequeOpReadNullable:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		mName := op.methodName
		if !op.discardRet && isPrimitiveUnboxBasicType(siteElemT) {
			switch mName {
			case "Peek", "PeekFirst":
				mName = "PeekFirstUnbox"
			case "PeekLast":
				mName = "PeekLastUnbox"
			case "Poll", "PollFirst":
				mName = "PollFirstUnbox"
			case "PollLast":
				mName = "PollLastUnbox"
			}
		}
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: mName}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case dequeOpContainsOrRem:
		parent := parents[op.outerExpr]
		pos := op.call.Pos()
		rawArg := op.call.Args[0]
		mName := op.methodName + "Any"
		arg := rawArg
		if op.argElemType != nil && types.Identical(op.argElemType, siteElemT) {
			mName = op.methodName
			arg = unwrapBoxArgExpr(rawArg, op.unwrapsBox)
		}
		repl := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: mName}},
			Args: []ast.Expr{arg},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case dequeOpAddAll:
		parent := parents[op.outerExpr]
		pos := op.call.Pos()
		var repl ast.Expr
		if op.otherCand != nil && op.otherCand.elemType != nil {
			otherElemT := instantiateOwnerTypeAtExpr(op.otherCand.elemType, op.otherCand.ownerStruct, op.argExpr, ctx.Info)
			if types.Identical(siteElemT, otherElemT) {
				otherArg := clonePureExpr(op.argExpr)
				if !op.otherCand.isPtr {
					otherArg = &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: otherArg}
				}
				repl = &ast.CallExpr{
					Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: "AddAllDeque"}},
					Args: []ast.Expr{otherArg},
				}
			}
		}
		if repl == nil && op.isSliceArg {
			repl = &ast.CallExpr{
				Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: "AddAllSlice"}},
				Args: []ast.Expr{op.argExpr},
			}
		}
		if repl == nil {
			repl = &ast.CallExpr{
				Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: "AddAllCol"}},
				Args: []ast.Expr{op.call.Args[0]},
			}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case dequeOpIterator:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		fnName := "DequeIterator"
		if op.methodName == "DescendingIterator" {
			fnName = "DequeDescendingIterator"
		}
		repl := &ast.CallExpr{Fun: j2Sel(fnName, pos), Args: []ast.Expr{dequePtrExpr(pos)}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case dequeOpIterCursor:
		initAs := op.stmt.(*ast.AssignStmt)
		pos := op.outerExpr.Pos()
		initAs.Rhs[0] = &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(clonePureExpr(op.accessExpr)), Sel: &ast.Ident{NamePos: pos, Name: "Iter"}},
		}
		curId := initAs.Lhs[0].(*ast.Ident)
		nextStmt := &ast.AssignStmt{
			Lhs:    []ast.Expr{&ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name}},
			TokPos: op.elemIdent.Pos(),
			Tok:    token.DEFINE,
			Rhs: []ast.Expr{&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{NamePos: op.elemIdent.Pos(), Name: curId.Name},
					Sel: &ast.Ident{NamePos: op.elemIdent.Pos(), Name: "Next"},
				},
			}},
		}
		op.forStmt.Body.List = append([]ast.Stmt{nextStmt}, op.forStmt.Body.List[op.dropStmts:]...)
		return true
	}
	return false
}

func isCandidateDequeType(t types.Type) bool {
	if t == nil {
		return false
	}
	ut := types.Unalias(t)
	if ptr, isPtr := ut.(*types.Pointer); isPtr {
		named, ok := types.Unalias(ptr.Elem()).(*types.Named)
		return ok && named.Obj() != nil && isJ2gortPkg(named.Obj().Pkg()) && named.Obj().Name() == "Ext_java_util_ArrayDeque"
	}
	named, ok := ut.(*types.Named)
	if !ok || named.Obj() == nil || !isJ2gortPkg(named.Obj().Pkg()) {
		return false
	}
	switch named.Obj().Name() {
	case "Ext_java_util_Deque", "Ext_java_util_Queue":
		return true
	}
	return false
}

func isDequeAdapterIface(t types.Type) bool {
	if t == nil {
		return false
	}
	ut := types.Unalias(t)
	if ptr, isPtr := ut.(*types.Pointer); isPtr {
		named, ok := types.Unalias(ptr.Elem()).(*types.Named)
		return ok && named.Obj() != nil && isJ2gortPkg(named.Obj().Pkg()) && named.Obj().Name() == "Ext_java_util_ArrayDeque"
	}
	named, ok := ut.(*types.Named)
	if !ok || named.Obj() == nil || !isJ2gortPkg(named.Obj().Pkg()) {
		return false
	}
	switch named.Obj().Name() {
	case "Ext_java_util_Deque", "Ext_java_util_Queue",
		"Ext_java_util_SequencedCollection", "Ext_java_util_Collection", "Ext_java_lang_Iterable":
		return true
	}
	return false
}

func unwrapDequeAdaptersDown(expr ast.Expr, info *types.Info) ast.Expr {
	cur := stripParens(expr)
	for {
		call, ok := cur.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
			return cur
		}
		if _, isTI := isToIfaceCall(call, info); isTI && isDequeAdapterIface(info.Types[call].Type) {
			cur = stripParens(call.Args[0])
			continue
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "CheckNotNull", "CheckNotNullPtr", "CheckNotNullIface":
				cur = stripParens(call.Args[0])
				continue
			case "AsIface", "Cast", "CastPtr", "UncheckedCast":
				if isDequeAdapterIface(info.Types[call].Type) {
					cur = stripParens(call.Args[0])
					continue
				}
			}
		}
		return cur
	}
}

func unwrapDequeAdaptersUp(expr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) ast.Expr {
	for {
		if pe, ok := parents[expr].(*ast.ParenExpr); ok && pe.X == expr {
			expr = pe
			continue
		}
		call, ok := parents[expr].(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() || call.Args[0] != expr {
			return expr
		}
		if _, isTI := isToIfaceCall(call, info); isTI && isDequeAdapterIface(info.Types[call].Type) {
			expr = call
			continue
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "CheckNotNull", "CheckNotNullPtr", "CheckNotNullIface", "PtrToAny", "ToAny":
				expr = call
				continue
			case "AsIface", "Cast", "UncheckedCast":
				if isDequeAdapterIface(info.Types[call].Type) {
					expr = call
					continue
				}
			}
		}
		return expr
	}
}

func resolveDequeCandAccess(expr ast.Expr, varCands map[*types.Var]*dequeCandidate, info *types.Info) *dequeCandidate {
	switch x := stripParens(expr).(type) {
	case *ast.Ident:
		if c := varCands[varFromObj(info.Uses[x])]; c != nil && !c.isField {
			return c
		}
	case *ast.SelectorExpr:
		if c := varCands[varFromObj(info.Uses[x.Sel])]; c != nil && c.isField && isIndirectStructFieldSel(x, info) {
			return c
		}
	}
	return nil
}

func matchWrappedSliceArg(expr ast.Expr, info *types.Info) (ast.Expr, types.Type, bool) {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
		return nil, nil, false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 {
		return nil, nil, false
	}
	switch fnName {
	case "ImmListFromTypedSlice", "ImmListFromSlice", "SliceToJavaList", "SliceUnmodifiableList":
		sl, okS := types.Unalias(info.Types[call.Args[0]].Type).(*types.Slice)
		if okS && sl.Elem() != nil && !isInvalidOrUntyped(sl.Elem()) && !isEmptyInterfaceType(sl.Elem()) {
			return call.Args[0], sl.Elem(), true
		}
	}
	return nil, nil, false
}

func extractFirstGenericTypeArg(t types.Type) (types.Type, bool) {
	if t == nil {
		return nil, false
	}
	ut := types.Unalias(t)
	if ptr, ok := ut.(*types.Pointer); ok {
		ut = types.Unalias(ptr.Elem())
	}
	named, ok := ut.(*types.Named)
	if !ok || named.TypeArgs() == nil || named.TypeArgs().Len() == 0 {
		return nil, false
	}
	ta := named.TypeArgs().At(0)
	if ta == nil || isInvalidOrUntyped(ta) || isEmptyInterfaceType(ta) {
		return nil, false
	}
	return ta, true
}

func unwrapDequeNodeWriteArg(arg ast.Expr, info *types.Info) (ast.Expr, types.Type, bool) {
	e := stripParens(arg)
	if conv, ok := e.(*ast.CallExpr); ok && len(conv.Args) == 1 && !conv.Ellipsis.IsValid() && info.Types[conv.Fun].IsType() {
		inner := stripParens(conv.Args[0])
		innerT := info.Types[inner].Type
		if isNamedUint32NodeType(innerT) {
			return inner, innerT, true
		}
	}
	t := info.Types[e].Type
	if t != nil && (types.Identical(types.Unalias(t), types.Typ[types.Uint32]) || isNamedUint32NodeType(t)) {
		return e, t, true
	}
	return nil, nil, false
}

func unwrapEnclosingNodeConvUp(call *ast.CallExpr, parents map[ast.Node]ast.Node, info *types.Info) (ast.Expr, types.Type) {
	cur := unwrapParensUp(call, parents)
	if conv, ok := parents[cur].(*ast.CallExpr); ok && len(conv.Args) == 1 && !conv.Ellipsis.IsValid() && conv.Args[0] == cur && info.Types[conv.Fun].IsType() {
		convT := info.Types[conv].Type
		if isNamedUint32NodeType(convT) {
			return conv, convT
		}
	}
	return call, types.Typ[types.Uint32]
}

func isPrimitiveUnboxBasicType(t types.Type) bool {
	if t == nil || isNamedUint32NodeType(t) {
		return false
	}
	b, ok := types.Unalias(t).(*types.Basic)
	if !ok {
		return false
	}
	switch b.Kind() {
	case types.Bool, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint16, types.Float32, types.Float64:
		return true
	}
	return false
}

func anyDequeLoopMutatesCandidate(cand *dequeCandidate, mutatingFuncs map[*types.Func]bool, info *types.Info) bool {
	if len(cand.loopBodies) == 0 || len(mutatingFuncs) == 0 {
		return false
	}
	for _, body := range cand.loopBodies {
		mut := false
		ast.Inspect(body, func(n ast.Node) bool {
			if mut || n == nil {
				return false
			}
			switch x := n.(type) {
			case *ast.AssignStmt:
				for _, lhs := range x.Lhs {
					if usesDequeVar(lhs, cand.v, info) {
						mut = true
						return false
					}
				}
			case *ast.CallExpr:
				if sel, ok := stripParens(x.Fun).(*ast.SelectorExpr); ok {
					rx := unwrapEnclosingCheckNotNullsDown(sel.X, info)
					if usesDequeVar(rx, cand.v, info) && isDequeMutatingMethodName(sel.Sel.Name) {
						mut = true
						return false
					}
				}
				if fnName, _, _, isJ2 := isJ2gortCall(x, info); isJ2 && len(x.Args) >= 1 {
					switch fnName {
					case "DequePushNode", "DequeAddLastNode", "DequePopNode", "DequeRemoveLastNode", "DequePollFirstNode":
						if usesDequeVar(unwrapDequeAdaptersDown(x.Args[0], info), cand.v, info) {
							mut = true
							return false
						}
					}
				}
				if cand.isField {
					if fn, _ := calleeFunc(x, info); fn != nil && mutatingFuncs[funcOrigin(fn)] {
						mut = true
						return false
					}
				}
			}
			return true
		})
		if mut {
			return true
		}
	}
	return false
}

func unwrapEnclosingCheckNotNullsDown(expr ast.Expr, info *types.Info) ast.Expr {
	return unwrapCheckNotNullArg(expr, info)
}

func usesDequeVar(expr ast.Expr, target *types.Var, info *types.Info) bool {
	switch x := stripParens(expr).(type) {
	case *ast.Ident:
		return varFromObj(info.Uses[x]) == target
	case *ast.SelectorExpr:
		return varFromObj(info.Uses[x.Sel]) == target
	}
	return false
}

func isDequeMutatingMethodName(name string) bool {
	switch name {
	case "Push", "Pop", "AddFirst", "AddLast", "Add", "Add_1", "Offer", "OfferFirst", "OfferLast",
		"Poll", "Poll_0", "PollFirst", "PollLast", "Remove_0", "RemoveFirst", "RemoveLast",
		"Remove_1_any", "RemoveFirstOccurrence", "RemoveLastOccurrence", "AddAll", "AddAll_1", "Clear":
		return true
	}
	return false
}
