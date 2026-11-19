// Clean-room immutable_list_lowering pass: lowers Guava ImmutableList locals, unexported
// struct fields (such as FunctionTypeBase.parameterList and TemplateTypeMap.templateKeys),
// and non-escaping local ImmutableList.Builder variables to typed Go []T slices.

package passes

import (
	"cmp"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strconv"
	"strings"
)

// ImmutableListLoweringPass lowers unexported struct fields and local variables of type
// j2gort.Ext_com_google_common_collect_ImmutableList (plus paired local
// *j2gort.Ext_com_google_common_collect_ImmutableList_Builder builders and same-package
// field getters) to native Go []T slices.
type ImmutableListLoweringPass struct{}

func NewImmutableListLoweringPass() *ImmutableListLoweringPass { return &ImmutableListLoweringPass{} }
func (p *ImmutableListLoweringPass) Name() string              { return "immutable_list_lowering" }
func (p *ImmutableListLoweringPass) Description() string {
	return "Lower ImmutableList fields/locals and local ImmutableList.Builder variables to native []T slices"
}

func (p *ImmutableListLoweringPass) Apply(ctx *PassContext) (int, error) {
	if isJ2gortPkg(ctx.Pkg) || len(ctx.Files) == 0 {
		return 0, nil
	}
	if !j2gortHasExport(ctx, "ImmListFromTypedSlice") {
		return 0, nil
	}
	parentsByFile := make(map[*ast.File]map[ast.Node]ast.Node, len(ctx.Files))
	for _, f := range ctx.Files {
		parentsByFile[f], _ = buildParentMap(f)
	}

	eligible, getters := analyzeImmCandidates(ctx, parentsByFile)
	if len(eligible) == 0 {
		return 0, nil
	}
	return applyImmRewrites(ctx, eligible, getters, parentsByFile), nil
}

type immInitKind uint8

const (
	immInitNil immInitKind = iota
	immInitBuilderEmpty
	immInitBuilderCap
	immInitOf0
	immInitOfFixed
	immInitOfVarargs
	immInitCopyCol
	immInitCopyLoweredSlice // CopyOf(ImmListFromSlice(arr)) -> ImmSliceCloneNonNull(arr)
	immInitCopyGetterSlice  // CopyOf(recv.GetFoo()) -> CheckNotNullSlice(recv...Slice_GetFoo())
	immInitCopyCand         // CopyOf(otherImmCand) -> CheckNotNullSlice(other) or other
	immInitCopyArray
	immInitFromCand       // dst = srcImmCand
	immInitFromBuilder    // dst = builder.BuildImmList_()
	immInitSubListCand    // dst = srcImmCand.SubList(from, to) -> ImmSliceSubList(src, from, to)
	immInitReverseCand    // dst = srcImmCand.Reverse() -> ImmSliceReverse(src)
	immInitAsListCand     // dst = srcImmCand.AsList() -> CheckNotNullSlice(src)
	immInitGetterFastPath // dst = recv.GetParameters() -> recv...Slice_GetParameters()
	immInitExternal       // dst = arbitraryImmutableListExpr -> ImmListToSliceNullable[T](expr)
	immInitIIFE
)

type immOpKind uint8

const (
	immOpDeclField immOpKind = iota
	immOpDeclVarSpec
	immOpDeclShortVar
	immOpAssignInit
	immOpNilCheck
	immOpSize
	immOpIsEmpty
	immOpGet
	immOpListGetStr
	immOpGetFirst
	immOpGetLast
	immOpContains
	immOpIndexOf
	immOpLastIndexOf
	immOpSubList
	immOpReverse
	immOpAsList
	immOpEquals
	immOpHashCode
	immOpIterator
	immOpRangeIterSlice
	immOpIterCursor
	immOpRawUnbox
	immOpBuilderAddStmt
	immOpBuilderAddExpr
	immOpBuilderAddAllStmt
	immOpBuilderAddAllExpr
	immOpBuilderBuildEscape
	immOpEscapeWrap
)

type immIIFEBranch struct {
	ret       *ast.ReturnStmt
	spec      immProducerSpec
	otherCand *immCandidate
}

type immProducerSpec struct {
	kind       immInitKind
	rawExpr    ast.Expr
	argExpr    ast.Expr
	srcCall    *ast.CallExpr
	elemExprs  []ast.Expr
	elemUnwrap []bool
	otherCand  *immCandidate
	getterInfo *immGetterInfo
	getterRecv ast.Expr
	iifeFl     *ast.FuncLit
	iifeBrs    []immIIFEBranch
}

type immGetterInfo struct {
	fnDecl        *ast.FuncDecl
	fnFile        *ast.File
	fnObj         *types.Func
	ownerStruct   *types.Named
	polyIface     *types.Named
	polyIfaceAST  *ast.InterfaceType
	polyIfaceFile *ast.File
	fieldCand     *immCandidate
	methodName    string
	sliceHelper   string
	baseGetter    string // e.g. "GetFunctionTypeBase"
	directOnIface bool   // true when interface itself declares Slice_<Getter>()
	isExternal    bool   // true when Slice_<Getter>() already exists on receiver type
}

type immRewriteOp struct {
	kind         immOpKind
	cand         *immCandidate
	otherCand    *immCandidate
	prod         immProducerSpec
	getterInfo   *immGetterInfo
	getterRecv   ast.Expr
	file         *ast.File
	spanNode     ast.Node
	accessExpr   ast.Expr
	argExpr      ast.Expr
	argExpr2     ast.Expr
	call         *ast.CallExpr
	outerExpr    ast.Expr
	stmt         ast.Stmt
	rangeStmt    *ast.RangeStmt
	forStmt      *ast.ForStmt
	elemIdent    *ast.Ident
	dropStmts    int
	unwrapsBox   bool
	isGuarded    bool
	isNegated    bool
	returnsSlice bool // true when SubList/Reverse/AsList directly feeds a lowered candidate
}

type immCandidate struct {
	v                *types.Var
	isField          bool
	isBuilder        bool
	isExternalGetter bool
	fixedElemType    types.Type
	ownerStruct      *types.Named
	declFile         *ast.File
	declPos          token.Pos
	declBlock        *ast.BlockStmt
	defCount         int
	hasInit          bool
	hasNilAssign     bool
	mayBeNil         bool
	disqualified     bool
	elemType         types.Type
	writeTypes       []types.Type
	readTypes        []types.Type
	siteConstraints  []siteTypeConstraint
	spellSites       []typeSpellSite
	j2gortSites      []typeSpellSite
	firstBuildPos    token.Pos
	maxMutPos        token.Pos
	ops              []immRewriteOp
	parent           *immCandidate
}

func (c *immCandidate) sortPos() token.Pos {
	if c.v != nil {
		return c.v.Pos()
	}
	return c.declPos
}

func (c *immCandidate) root() *immCandidate {
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

func unionImmCandidates(a, b *immCandidate) {
	ra, rb := a.root(), b.root()
	if ra != rb {
		rb.parent = ra
	}
}

func (c *immCandidate) addWriteType(t types.Type, accessExpr ast.Expr, info *types.Info) bool {
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
	}
	c.writeTypes = append(c.writeTypes, t)
	return true
}

func (c *immCandidate) addReadType(t types.Type, accessExpr ast.Expr, info *types.Info) bool {
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
	}
	c.readTypes = append(c.readTypes, t)
	return true
}

func analyzeImmCandidates(
	ctx *PassContext,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
) ([]*immCandidate, []*immGetterInfo) {
	info := ctx.Info
	varUseCounts := make(map[*types.Var]int, len(info.Uses))
	// order-independent: commutative per-variable integer use-count accumulation
	for _, obj := range info.Uses {
		if v := varFromObj(obj); v != nil {
			varUseCounts[v]++
		}
	}

	varCands := make(map[*types.Var]*immCandidate)
	var allCands []*immCandidate
	var structOrder []*types.Named
	structFields := make(map[*types.Named][]*immCandidate)
	var edges [][2]*immCandidate
	var mayBeNilDeps [][2]*immCandidate // [dst, src]: if src.mayBeNil -> dst.mayBeNil

	// Pass 1a: Collect unexported ImmutableList struct fields.
	for _, file := range ctx.Files {
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, s := range gd.Specs {
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
					if v == nil || !v.IsField() || v.Exported() || v.Embedded() || !isExactImmutableListType(v.Type()) {
						continue
					}
					cand := &immCandidate{
						v: v, isField: true, ownerStruct: named, declFile: file, declPos: fld.Pos(),
						spellSites:  []typeSpellSite{{file: file, pos: fld.Pos(), owner: named}},
						j2gortSites: []typeSpellSite{{file: file, pos: fld.Pos()}},
					}
					cand.ops = append(cand.ops, immRewriteOp{kind: immOpDeclField, cand: cand, file: file, spanNode: fld})
					varCands[v] = cand
					allCands = append(allCands, cand)
					if len(structFields[named]) == 0 {
						structOrder = append(structOrder, named)
					}
					structFields[named] = append(structFields[named], cand)
				}
			}
		}
	}

	// Pass 1b: Collect local ImmutableList.Builder variables first (so ImmutableList locals
	// initialized from `builder.BuildImmList_()` can link directly to them!).
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		addBuilder := func(id *ast.Ident, initExpr ast.Expr, span ast.Node, opK immOpKind) {
			v := varFromObj(info.Defs[id])
			if id.Name == "_" || v == nil || v.IsField() || initExpr == nil || !isImmListBuilderPtrType(v.Type()) {
				return
			}
			initK, capExpr, ok := classifyImmBuilderInit(initExpr, info)
			if !ok {
				return
			}
			cand := &immCandidate{
				v: v, isBuilder: true, declFile: file, declPos: span.Pos(),
				declBlock:   enclosingBlockStmt(span, parents),
				hasInit:     true,
				spellSites:  []typeSpellSite{{file: file, pos: span.Pos()}},
				j2gortSites: []typeSpellSite{{file: file, pos: span.Pos()}},
			}
			cand.ops = append(cand.ops, immRewriteOp{
				kind: opK, cand: cand, prod: immProducerSpec{kind: initK, rawExpr: initExpr, argExpr: capExpr},
				file: file, spanNode: span,
			})
			varCands[v] = cand
			allCands = append(allCands, cand)
		}

		ast.Inspect(file, func(n ast.Node) bool {
			switch s := n.(type) {
			case *ast.DeclStmt:
				if gd, ok := s.Decl.(*ast.GenDecl); ok && gd.Tok == token.VAR && len(gd.Specs) == 1 && enclosingFuncBody(s, parents) != nil {
					if spec, ok := gd.Specs[0].(*ast.ValueSpec); ok && len(spec.Names) == 1 && len(spec.Values) == 1 {
						addBuilder(spec.Names[0], spec.Values[0], spec, immOpDeclVarSpec)
					}
				}
			case *ast.AssignStmt:
				if s.Tok == token.DEFINE && len(s.Lhs) == 1 && len(s.Rhs) == 1 && enclosingFuncBody(s, parents) != nil {
					if id, ok := s.Lhs[0].(*ast.Ident); ok {
						addBuilder(id, s.Rhs[0], s, immOpDeclShortVar)
					}
				}
			}
			return true
		})
	}

	// Discover simple same-package getters `func (this_ *Owner) GetFoo() ImmutableList/List { return this_.field }`
	// for non-generic structs so in-package callers can use zero-alloc `Slice_GetFoo()`.
	getters := discoverImmFieldGetters(ctx, varCands)
	getterByMethod := make(map[*types.Func]*immGetterInfo, len(getters))
	getterByIfaceName := make(map[string]*immGetterInfo)
	for _, g := range getters {
		getterByMethod[g.fnObj] = g
		if g.baseGetter != "" {
			getterByIfaceName[g.methodName+"\x00"+g.baseGetter] = g
		}
	}

	// Pass 1c: Collect local ImmutableList variables.
	var localInitDefer []func()
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		addImmLocal := func(id *ast.Ident, initExpr ast.Expr, span ast.Node, opK immOpKind) {
			v := varFromObj(info.Defs[id])
			if id.Name == "_" || v == nil || v.IsField() || !isExactImmutableListType(v.Type()) {
				return
			}
			cand := &immCandidate{
				v: v, declFile: file, declPos: span.Pos(),
				declBlock:   enclosingBlockStmt(span, parents),
				spellSites:  []typeSpellSite{{file: file, pos: span.Pos()}},
				j2gortSites: []typeSpellSite{{file: file, pos: span.Pos()}},
			}
			varCands[v] = cand
			allCands = append(allCands, cand)

			f := file
			localInitDefer = append(localInitDefer, func() {
				if initExpr == nil {
					cand.hasNilAssign = true
					cand.ops = append(cand.ops, immRewriteOp{
						kind: opK, cand: cand, prod: immProducerSpec{kind: immInitNil}, file: f, spanNode: span,
					})
					return
				}
				prod, prodMayNil, ok := classifyImmProducer(initExpr, nil, f, cand, varCands, getterByMethod, getterByIfaceName, ctx.Pkg, info, &edges, &mayBeNilDeps)
				if !ok {
					cand.disqualified = true
					return
				}
				if prodMayNil {
					cand.hasNilAssign = true
				} else {
					cand.hasInit = true
				}
				cand.ops = append(cand.ops, immRewriteOp{
					kind: opK, cand: cand, prod: prod, otherCand: prod.otherCand, file: f, spanNode: span,
				})
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
						addImmLocal(spec.Names[0], initExpr, spec, immOpDeclVarSpec)
					}
				}
			case *ast.AssignStmt:
				if s.Tok == token.DEFINE && len(s.Lhs) == 1 && len(s.Rhs) == 1 && enclosingFuncBody(s, parents) != nil {
					if id, ok := s.Lhs[0].(*ast.Ident); ok {
						addImmLocal(id, s.Rhs[0], s, immOpDeclShortVar)
					}
				}
			}
			return true
		})
	}
	for _, fn := range localInitDefer {
		fn()
	}

	if len(structOrder) > 0 {
		verifyImmStructFieldInits(ctx, varCands, structOrder, structFields, getterByMethod, getterByIfaceName)
	}

	// Build a set of AST expressions already consumed as part of a candidate initializer
	// (e.g. `builder.BuildImmList_()` or `otherCand` in `dst = otherCand`).
	consumedExprs := make(map[ast.Expr]bool)
	recordConsumedProd := func(p *immProducerSpec) {
		var rec func(*immProducerSpec)
		rec = func(ps *immProducerSpec) {
			switch ps.kind {
			case immInitFromBuilder:
				if ps.rawExpr != nil {
					consumedExprs[stripParens(ps.rawExpr)] = true
				}
			case immInitFromCand, immInitCopyCand, immInitCopyGetterSlice, immInitGetterFastPath,
				immInitSubListCand, immInitReverseCand, immInitAsListCand:
				if ps.rawExpr != nil {
					consumedExprs[stripParens(ps.rawExpr)] = true
				}
				if ps.argExpr != nil {
					consumedExprs[stripParens(ps.argExpr)] = true
				}
			case immInitIIFE:
				for i := range ps.iifeBrs {
					rec(&ps.iifeBrs[i].spec)
				}
			}
		}
		rec(p)
	}
	for _, cand := range allCands {
		for i := range cand.ops {
			recordConsumedProd(&cand.ops[i].prod)
		}
	}

	// Pass 2a: Pre-scan assignment RHSs to candidates so `dst = builder.BuildImmList_()` and
	// `dst = srcCand` are marked in `consumedExprs` before consumer validation!
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		ast.Inspect(file, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			cand := varCands[varFromObj(info.Uses[id])]
			if cand == nil || cand.disqualified || cand.isBuilder {
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
				if cand.isField && !isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info) {
					cand.disqualified = true
					return true
				}
				prod, prodMayNil, ok := classifyImmProducer(as.Rhs[0], accessExpr, file, cand, varCands, getterByMethod, getterByIfaceName, ctx.Pkg, info, &edges, &mayBeNilDeps)
				if !ok {
					cand.disqualified = true
					return true
				}
				if prodMayNil {
					cand.hasNilAssign = true
				}
				recordConsumedProd(&prod)
				cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
				cand.ops = append(cand.ops, immRewriteOp{
					kind: immOpAssignInit, cand: cand, prod: prod, otherCand: prod.otherCand,
					file: file, spanNode: as, stmt: as, accessExpr: accessExpr,
				})
				consumedExprs[accessExpr] = true
			}
			return true
		})
	}

	// Pass 2b: Validate all uses of every candidate and direct fast-path getter calls
	// (both same-package synthesized `Slice_<Getter>` and cross-package/existing `Slice_<Getter>`).
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.Ident:
				if cand := varCands[varFromObj(info.Defs[x])]; cand != nil {
					cand.defCount++
					if !ctx.ShouldTransform(file) {
						cand.disqualified = true
					}
				}
				if cand := varCands[varFromObj(info.Uses[x])]; cand != nil && !cand.disqualified {
					if !ctx.ShouldTransform(file) {
						cand.disqualified = true
					} else if cand.isBuilder {
						validateImmBuilderUse(x, file, cand, varCands, getterByMethod, getterByIfaceName, ctx.Pkg, parents, info, &edges)
					} else {
						validateImmListUse(x, file, cand, varCands, consumedExprs, parents, info, varUseCounts, &edges)
					}
				}
			case *ast.CallExpr:
				if !ctx.ShouldTransform(file) || consumedExprs[x] {
					return true
				}
				if gInfo, recvExpr, ok := matchImmGetterFastCall(x, getterByMethod, getterByIfaceName, info); ok && !gInfo.fieldCand.disqualified {
					tryAttachDirectGetterConsumer(x, gInfo, recvExpr, file, varCands, parents, info, varUseCounts, &edges)
				} else if extGInfo, extElemT, recvExpr, okE := matchExternalSliceGetterCall(x, ctx.Pkg, info); okE {
					extCand := &immCandidate{
						isExternalGetter: true,
						fixedElemType:    extElemT,
						declFile:         file,
						declPos:          x.Pos(),
						defCount:         1,
						hasNilAssign:     true,
						mayBeNil:         true,
						writeTypes:       []types.Type{extElemT},
						spellSites:       []typeSpellSite{{file: file, pos: x.Pos()}},
						j2gortSites:      []typeSpellSite{{file: file, pos: x.Pos()}},
					}
					if validateImmConsumerChain(x, extGInfo, recvExpr, false, file, extCand, varCands, parents, info, varUseCounts, &edges) && !extCand.disqualified && len(extCand.ops) > 0 {
						allCands = append(allCands, extCand)
					}
				}
			}
			return true
		})
	}

	if len(allCands) == 0 {
		return nil, nil
	}

	// Propagate mayBeNil across assignment/copy dependencies to fixed point.
	for _, cand := range allCands {
		if !cand.isBuilder && (!cand.hasInit || cand.hasNilAssign) {
			cand.mayBeNil = true
		}
	}
	for changed := true; changed; {
		changed = false
		for _, dep := range mayBeNilDeps {
			dst, src := dep[0], dep[1]
			if !dst.mayBeNil && src.mayBeNil {
				dst.mayBeNil = true
				changed = true
			}
		}
	}

	for _, e := range edges {
		unionImmCandidates(e[0], e[1])
	}

	for _, cand := range allCands {
		if cand.disqualified || cand.defCount != 1 || !ctx.ShouldTransform(cand.declFile) {
			cand.disqualified = true
			continue
		}
		if cand.isBuilder {
			if !cand.firstBuildPos.IsValid() || (cand.maxMutPos.IsValid() && cand.maxMutPos > cand.firstBuildPos) {
				cand.disqualified = true
				continue
			}
		}
	}

	var groupRoots []*immCandidate
	groups := make(map[*immCandidate][]*immCandidate)
	for _, cand := range allCands {
		r := cand.root()
		if len(groups[r]) == 0 {
			groupRoots = append(groupRoots, r)
		}
		groups[r] = append(groups[r], cand)
	}

	var out []*immCandidate
	eligibleSet := make(map[*immCandidate]bool)
	for _, r := range groupRoots {
		members := groups[r]
		if !resolveAndValidateImmGroup(members, ctx) {
			continue
		}
		for _, m := range members {
			eligibleSet[m] = true
			out = append(out, m)
		}
	}

	var activeGetters []*immGetterInfo
	for _, g := range getters {
		if eligibleSet[g.fieldCand] && ctx.ShouldTransform(g.fnFile) {
			activeGetters = append(activeGetters, g)
		}
	}

	slices.SortFunc(out, func(a, b *immCandidate) int { return cmp.Compare(a.sortPos(), b.sortPos()) })
	return out, activeGetters
}

func resolveAndValidateImmGroup(members []*immCandidate, ctx *PassContext) bool {
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
		if m.fixedElemType != nil {
			if !types.Identical(elemT, m.fixedElemType) {
				return false
			}
			for _, rt := range m.readTypes {
				if !types.Identical(rt, m.fixedElemType) {
					return false
				}
			}
		}
		if typeContainsTypeParam(elemT) && (!m.isField || !typeOnlyUsesOwnerTypeParams(elemT, m.ownerStruct)) {
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

func resolveImmElemType(all, reads, writes []types.Type) (types.Type, bool) {
	t0 := all[0]
	allIdentical := true
	for _, t := range all[1:] {
		if !types.Identical(t0, t) {
			allIdentical = false
			break
		}
	}
	if allIdentical {
		return t0, true
	}

	// ImmutableList forbids null elements (checked via CheckNotNullElem on every write), so any
	// subtype assignable to a common non-empty interface I can be stored safely in []I with zero
	// typed-nil interface hazard.
	var candidates []types.Type
	if len(reads) > 0 {
		candidates = reads
	} else {
		candidates = writes
	}
	var targetIface types.Type
	for _, ct := range candidates {
		if isNonEmptyInterfaceType(ct) {
			if targetIface == nil || types.AssignableTo(targetIface, ct) {
				targetIface = ct
			}
		}
	}
	if targetIface == nil || !isNonEmptyInterfaceType(targetIface) {
		return nil, false
	}
	for _, rt := range reads {
		if !types.Identical(rt, targetIface) {
			return nil, false
		}
	}
	for _, wt := range writes {
		if !types.Identical(wt, targetIface) && !types.AssignableTo(wt, targetIface) {
			return nil, false
		}
	}
	return targetIface, true
}

func discoverImmFieldGetters(ctx *PassContext, varCands map[*types.Var]*immCandidate) []*immGetterInfo {
	info := ctx.Info
	type polyIfaceSpec struct {
		named    *types.Named
		ifaceObj *types.Interface
		astNode  *ast.InterfaceType
		file     *ast.File
	}

	// Count concrete methods per name in ctx.Pkg to ensure no subclass overrides `<Getter>()`,
	// and collect polymorphic interface AST declarations in ctx.Files.
	methodDeclsByName := make(map[string][]*types.Func)
	existingFuncNames := make(map[string]bool)
	anyMethodNameExists := make(map[string]bool)
	var polyIfaces []polyIfaceSpec

	for _, file := range ctx.Files {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Name == nil {
					continue
				}
				if fnObj, ok := info.Defs[d.Name].(*types.Func); ok {
					orig := funcOrigin(fnObj)
					if orig != nil && d.Recv != nil {
						anyMethodNameExists[d.Name.Name] = true
						methodDeclsByName[d.Name.Name] = append(methodDeclsByName[d.Name.Name], orig)
						if recvNamed := ptrNamedOrigin(orig.Type().(*types.Signature).Recv().Type()); recvNamed != nil {
							existingFuncNames[recvNamed.Obj().Name()+"."+d.Name.Name] = true
						}
					}
				}
			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}
				for _, s := range d.Specs {
					ts, ok := s.(*ast.TypeSpec)
					if !ok || ts.Assign.IsValid() {
						continue
					}
					it, itOK := ts.Type.(*ast.InterfaceType)
					tn, tnOK := info.Defs[ts.Name].(*types.TypeName)
					if !itOK || !tnOK || it.Methods == nil {
						continue
					}
					named := namedOrigin(tn.Type())
					if named == nil || named.TypeParams().Len() > 0 {
						continue
					}
					if ifaceObj, ok := named.Underlying().(*types.Interface); ok {
						polyIfaces = append(polyIfaces, polyIfaceSpec{
							named:    named,
							ifaceObj: ifaceObj,
							astNode:  it,
							file:     file,
						})
					}
				}
			}
		}
	}

	// Map `*OwnerStruct` -> `Get<OwnerStructName>` if an interface in ctx.Pkg exposes it,
	// and collect all package structs to verify polymorphic interface implementors.
	baseGetterExists := make(map[string]bool)
	var pkgStructs []*types.Named
	scope := ctx.Pkg.Scope()
	for _, name := range scope.Names() {
		tn, ok := scope.Lookup(name).(*types.TypeName)
		if !ok {
			continue
		}
		if named := namedOrigin(tn.Type()); named != nil && named.TypeParams().Len() == 0 {
			if _, isSt := named.Underlying().(*types.Struct); isSt {
				pkgStructs = append(pkgStructs, named)
			}
		}
		iface, ok := tn.Type().Underlying().(*types.Interface)
		if !ok {
			continue
		}
		for i := 0; i < iface.NumMethods(); i++ {
			m := iface.Method(i)
			anyMethodNameExists[m.Name()] = true
			if strings.HasPrefix(m.Name(), "Get") {
				if sig, ok := m.Type().(*types.Signature); ok && sig.Params().Len() == 0 && sig.Results().Len() == 1 {
					if rn := ptrNamedOrigin(sig.Results().At(0).Type()); rn != nil && m.Name() == "Get"+rn.Obj().Name() {
						baseGetterExists[rn.Obj().Name()] = true
					}
				}
			}
		}
	}

	addedPolyMethod := make(map[string]bool)
	var out []*immGetterInfo
	for _, file := range ctx.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 || fn.Body == nil || len(fn.Body.List) == 0 {
				continue
			}
			if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
				continue
			}
			if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 || len(fn.Type.Results.List[0].Names) > 0 {
				continue
			}
			ret, ok := fn.Body.List[len(fn.Body.List)-1].(*ast.ReturnStmt)
			if !ok || len(ret.Results) != 1 {
				continue
			}
			if len(fn.Body.List) > 1 && !prefixHasNoEarlyExitOrDefer(fn.Body.List[:len(fn.Body.List)-1]) {
				continue
			}
			recvVar := varFromObj(info.Defs[fn.Recv.List[0].Names[0]])
			if recvVar == nil {
				continue
			}
			owner := ptrNamedOrigin(recvVar.Type())
			if owner == nil || owner.TypeParams().Len() > 0 {
				continue
			}
			resT := info.Types[fn.Type.Results.List[0].Type].Type
			if !isImmEscapeIface(resT) {
				continue
			}
			retExpr := unwrapCollectionAdaptersDown(ret.Results[0], info)
			sel, ok := retExpr.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			recvId, ok := stripParens(sel.X).(*ast.Ident)
			if !ok || varFromObj(info.Uses[recvId]) != recvVar {
				continue
			}
			fieldCand := varCands[varFromObj(info.Uses[sel.Sel])]
			if fieldCand == nil || !fieldCand.isField || fieldCand.ownerStruct != owner {
				continue
			}
			fnObj, ok := info.Defs[fn.Name].(*types.Func)
			if !ok {
				continue
			}
			sliceHelper := "Slice_" + fn.Name.Name
			if existingFuncNames[owner.Obj().Name()+"."+sliceHelper] {
				continue
			}

			// Verify no other concrete struct embedding `owner` overrides `fn.Name.Name`.
			overridden := false
			for _, otherFn := range methodDeclsByName[fn.Name.Name] {
				if otherFn == funcOrigin(fnObj) {
					continue
				}
				otherRecv := ptrNamedOrigin(otherFn.Type().(*types.Signature).Recv().Type())
				if otherRecv != nil && structEmbedsNamed(otherRecv, owner) {
					overridden = true
					break
				}
			}
			if overridden {
				continue
			}

			baseGetter := ""
			if baseGetterExists[owner.Obj().Name()] {
				baseGetter = "Get" + owner.Obj().Name()
			}

			var polyNamed *types.Named
			var polyAST *ast.InterfaceType
			var polyFile *ast.File
			directOnIface := false
			if baseGetter != "" && !anyMethodNameExists[sliceHelper] {
				for _, pi := range polyIfaces {
					if !ctx.ShouldTransform(pi.file) {
						continue
					}
					if !explicitIfaceHasMethod(pi.astNode, baseGetter) || !explicitIfaceHasMethod(pi.astNode, fn.Name.Name) {
						continue
					}
					if obj, _, _ := types.LookupFieldOrMethod(pi.named, false, ctx.Pkg, sliceHelper); obj != nil {
						continue
					}
					if addedPolyMethod[pi.named.Obj().Name()+"."+sliceHelper] {
						continue
					}
					allEmbed := true
					for _, st := range pkgStructs {
						if types.Implements(types.NewPointer(st), pi.ifaceObj) && !structEmbedsNamed(st, owner) {
							allEmbed = false
							break
						}
					}
					if allEmbed {
						polyNamed = pi.named
						polyAST = pi.astNode
						polyFile = pi.file
						directOnIface = true
						addedPolyMethod[pi.named.Obj().Name()+"."+sliceHelper] = true
						fieldCand.spellSites = append(fieldCand.spellSites, typeSpellSite{file: pi.file, pos: pi.astNode.Pos(), owner: owner})
						break
					}
				}
			}

			fieldCand.spellSites = append(fieldCand.spellSites, typeSpellSite{file: file, pos: fn.Pos(), owner: owner})
			if len(fn.Body.List) > 1 {
				fieldCand.j2gortSites = append(fieldCand.j2gortSites, typeSpellSite{file: file, pos: fn.Pos()})
			}
			out = append(out, &immGetterInfo{
				fnDecl:        fn,
				fnFile:        file,
				fnObj:         funcOrigin(fnObj),
				ownerStruct:   owner,
				polyIface:     polyNamed,
				polyIfaceAST:  polyAST,
				polyIfaceFile: polyFile,
				fieldCand:     fieldCand,
				methodName:    fn.Name.Name,
				sliceHelper:   sliceHelper,
				baseGetter:    baseGetter,
				directOnIface: directOnIface,
			})
		}
	}
	return out
}

func prefixHasNoEarlyExitOrDefer(stmts []ast.Stmt) bool {
	clean := true
	for _, st := range stmts {
		ast.Inspect(st, func(n ast.Node) bool {
			if !clean || n == nil {
				return false
			}
			if _, isLit := n.(*ast.FuncLit); isLit {
				return false
			}
			switch x := n.(type) {
			case *ast.ReturnStmt, *ast.DeferStmt, *ast.LabeledStmt:
				clean = false
				return false
			case *ast.BranchStmt:
				if x.Tok == token.GOTO {
					clean = false
					return false
				}
			}
			return true
		})
		if !clean {
			return false
		}
	}
	return true
}

func explicitIfaceHasMethod(it *ast.InterfaceType, methodName string) bool {
	if it == nil || it.Methods == nil {
		return false
	}
	for _, f := range it.Methods.List {
		for _, id := range f.Names {
			if id.Name == methodName {
				return true
			}
		}
	}
	return false
}

func structEmbedsNamed(sub, target *types.Named) bool {
	if sub == nil || target == nil {
		return false
	}
	if sub == target {
		return true
	}
	st, ok := sub.Underlying().(*types.Struct)
	if !ok {
		return false
	}
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if !f.Embedded() {
			continue
		}
		if emb := derefNamedOrigin(f.Type()); emb != nil && structEmbedsNamed(emb, target) {
			return true
		}
	}
	return false
}

func matchImmGetterFastCall(
	call *ast.CallExpr,
	getterByMethod map[*types.Func]*immGetterInfo,
	getterByIfaceName map[string]*immGetterInfo,
	info *types.Info,
) (*immGetterInfo, ast.Expr, bool) {
	if call == nil || len(call.Args) != 0 || call.Ellipsis.IsValid() {
		return nil, nil, false
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok {
		return nil, nil, false
	}
	fn, _ := calleeFunc(call, info)
	if fn == nil {
		return nil, nil, false
	}
	fnOrig := funcOrigin(fn)
	if g := getterByMethod[fnOrig]; g != nil {
		return g, sel.X, true
	}
	// Interface receiver (e.g. `subtype.GetParameters()` where `subtype` has interface type `FunctionType`
	// exposing `GetFunctionTypeBase() *FunctionTypeBase`).
	recvT := info.Types[sel.X].Type
	if recvT == nil || !isNonEmptyInterfaceType(recvT) {
		return nil, nil, false
	}
	iface, ok := types.Unalias(recvT).Underlying().(*types.Interface)
	if !ok {
		return nil, nil, false
	}
	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)
		if strings.HasPrefix(m.Name(), "Get") {
			if g := getterByIfaceName[sel.Sel.Name+"\x00"+m.Name()]; g != nil {
				return g, sel.X, true
			}
		}
	}
	return nil, nil, false
}

func matchExternalSliceGetterCall(
	call *ast.CallExpr,
	pkg *types.Package,
	info *types.Info,
) (*immGetterInfo, types.Type, ast.Expr, bool) {
	if call == nil || len(call.Args) != 0 || call.Ellipsis.IsValid() {
		return nil, nil, nil, false
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || strings.HasPrefix(sel.Sel.Name, "Slice_") {
		return nil, nil, nil, false
	}
	fn, _ := calleeFunc(call, info)
	if fn == nil {
		return nil, nil, nil, false
	}
	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Recv() == nil {
		return nil, nil, nil, false
	}
	retT := info.Types[call].Type
	if !isExactImmutableListType(retT) && !isExactJavaListType(retT) {
		return nil, nil, nil, false
	}
	recvT := info.Types[sel.X].Type
	if recvT == nil || isInvalidOrUntyped(recvT) || isEmptyInterfaceType(recvT) {
		return nil, nil, nil, false
	}
	sliceHelper := "Slice_" + sel.Sel.Name
	checkSliceMethod := func(targetT types.Type) (types.Type, bool) {
		obj, _, _ := types.LookupFieldOrMethod(targetT, false, pkg, sliceHelper)
		mFn, ok := obj.(*types.Func)
		if !ok || !mFn.Exported() {
			return nil, false
		}
		mSig, ok := mFn.Type().(*types.Signature)
		if !ok || mSig.Params().Len() != 0 || mSig.Results().Len() != 1 {
			return nil, false
		}
		sl, ok := types.Unalias(mSig.Results().At(0).Type()).(*types.Slice)
		if !ok || sl.Elem() == nil || isInvalidOrUntyped(sl.Elem()) || isEmptyInterfaceType(sl.Elem()) {
			return nil, false
		}
		return sl.Elem(), true
	}

	if elemT, ok := checkSliceMethod(recvT); ok {
		return &immGetterInfo{
			methodName:    sel.Sel.Name,
			sliceHelper:   sliceHelper,
			directOnIface: isNonEmptyInterfaceType(recvT),
			isExternal:    true,
		}, elemT, sel.X, true
	}

	// Fallback if interface exposes `Get<Base>() *<Base>` where `*<Base>` has `Slice_<Getter>()`.
	if isNonEmptyInterfaceType(recvT) {
		if iface, ok := types.Unalias(recvT).Underlying().(*types.Interface); ok {
			for i := 0; i < iface.NumMethods(); i++ {
				m := iface.Method(i)
				if !strings.HasPrefix(m.Name(), "Get") {
					continue
				}
				mSig, ok := m.Type().(*types.Signature)
				if !ok || mSig.Params().Len() != 0 || mSig.Results().Len() != 1 {
					continue
				}
				resPtr := mSig.Results().At(0).Type()
				if rn := ptrNamedOrigin(resPtr); rn != nil && m.Name() == "Get"+rn.Obj().Name() {
					if elemT, okS := checkSliceMethod(resPtr); okS {
						return &immGetterInfo{
							methodName:  sel.Sel.Name,
							sliceHelper: sliceHelper,
							baseGetter:  m.Name(),
							isExternal:  true,
						}, elemT, sel.X, true
					}
				}
			}
		}
	}
	return nil, nil, nil, false
}

func isExactJavaListType(t types.Type) bool {
	if t == nil {
		return false
	}
	named, ok := types.Unalias(t).(*types.Named)
	return ok && named.Obj() != nil && isJ2gortPkg(named.Obj().Pkg()) && named.Obj().Name() == "Ext_java_util_List"
}

func verifyImmStructFieldInits(
	ctx *PassContext,
	varCands map[*types.Var]*immCandidate,
	structOrder []*types.Named,
	structFields map[*types.Named][]*immCandidate,
	getterByMethod map[*types.Func]*immGetterInfo,
	getterByIfaceName map[string]*immGetterInfo,
) {
	info := ctx.Info
	isGuaranteedNonNilRhs := func(rhs ast.Expr) bool {
		return isNonNilImmInitExpr(rhs, varCands, getterByMethod, getterByIfaceName, info)
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
							if cand := varCands[varFromObj(info.Uses[sel.Sel])]; cand != nil && cand.isField && isGuaranteedNonNilRhs(as.Rhs[0]) {
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
								if id, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == recvVar && varFromObj(info.Uses[sel.Sel]) == cand.v && isGuaranteedNonNilRhs(s.Rhs[0]) {
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
}

func isNonNilImmInitExpr(
	expr ast.Expr,
	varCands map[*types.Var]*immCandidate,
	getterByMethod map[*types.Func]*immGetterInfo,
	getterByIfaceName map[string]*immGetterInfo,
	info *types.Info,
) bool {
	e := stripParens(expr)
	if isNilIdent(e, info) {
		return false
	}
	if call, ok := e.(*ast.CallExpr); ok && !call.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			if fnName == "Ext_com_google_common_collect_ImmutableList_Of_0" ||
				strings.HasPrefix(fnName, "Ext_com_google_common_collect_ImmutableList_Of_") ||
				strings.HasPrefix(fnName, "Ext_com_google_common_collect_ImmutableList_CopyOf_") {
				return true
			}
		}
		if sel, ok := stripParens(call.Fun).(*ast.SelectorExpr); ok && (sel.Sel.Name == "BuildImmList_" || sel.Sel.Name == "Build") && len(call.Args) == 0 {
			return true
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
						if len(ret.Results) != 1 || !isNonNilImmInitExpr(ret.Results[0], varCands, getterByMethod, getterByIfaceName, info) {
							allOK = false
						}
						return false
					}
					return true
				})
				return hasRet && allOK
			}
		}
		_ = getterByMethod
		_ = getterByIfaceName
		_ = varCands
	}
	return false
}

func classifyImmProducer(
	expr ast.Expr,
	accessExpr ast.Expr,
	file *ast.File,
	cand *immCandidate,
	varCands map[*types.Var]*immCandidate,
	getterByMethod map[*types.Func]*immGetterInfo,
	getterByIfaceName map[string]*immGetterInfo,
	pkg *types.Package,
	info *types.Info,
	edges *[][2]*immCandidate,
	mayBeNilDeps *[][2]*immCandidate,
) (immProducerSpec, bool, bool) {
	e := stripParens(expr)
	if isNilIdent(e, info) {
		return immProducerSpec{kind: immInitNil, rawExpr: expr}, true, true
	}

	// Candidate-to-candidate copy (`dst = src`).
	if srcCand := resolveImmListCandAccess(e, varCands, info); srcCand != nil {
		if srcCand != cand {
			*edges = append(*edges, [2]*immCandidate{cand, srcCand})
			*mayBeNilDeps = append(*mayBeNilDeps, [2]*immCandidate{cand, srcCand})
		}
		return immProducerSpec{kind: immInitFromCand, rawExpr: expr, argExpr: e, otherCand: srcCand}, false, true
	}

	// Ternary / multi-branch IIFE `func() Ext_com_google_common_collect_ImmutableList { ... }()`.
	if call, ok := e.(*ast.CallExpr); ok && len(call.Args) == 0 && !call.Ellipsis.IsValid() {
		if fl, ok := stripParens(call.Fun).(*ast.FuncLit); ok && fl.Body != nil && (fl.Type.Params == nil || len(fl.Type.Params.List) == 0) && fl.Type.Results != nil && len(fl.Type.Results.List) == 1 {
			var brs []immIIFEBranch
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
				subSpec, subNil, subOK := classifyImmProducer(ret.Results[0], accessExpr, file, cand, varCands, getterByMethod, getterByIfaceName, pkg, info, edges, mayBeNilDeps)
				if !subOK {
					valid = false
					return false
				}
				if subNil {
					mayNil = true
				}
				brs = append(brs, immIIFEBranch{ret: ret, spec: subSpec, otherCand: subSpec.otherCand})
				return false
			})
			if valid && len(brs) > 0 {
				cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: call.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: call.Pos()})
				return immProducerSpec{kind: immInitIIFE, rawExpr: expr, iifeFl: fl, iifeBrs: brs}, mayNil, true
			}
		}
	}

	// Unwrap outer AsIface[ImmutableList] / Cast[ImmutableList] around `Build()`, `SubList()`, etc.
	unwrapped := unwrapImmListCastDown(e, info)

	// Builder `.BuildImmList_()` / `.Build()` (including chained `b.AddSelf_(...).BuildImmList_()`).
	if bCall, ok := unwrapped.(*ast.CallExpr); ok && len(bCall.Args) == 0 && !bCall.Ellipsis.IsValid() {
		if sel, ok := stripParens(bCall.Fun).(*ast.SelectorExpr); ok && (sel.Sel.Name == "BuildImmList_" || sel.Sel.Name == "Build") {
			bId := unwrapBuilderChainRoot(sel.X, info)
			if bId != nil {
				if bCand := varCands[varFromObj(info.Uses[bId])]; bCand != nil && bCand.isBuilder && !bCand.disqualified {
					if !recordBuilderBuildSite(bCand, bId, bCall) {
						return immProducerSpec{}, false, false
					}
					*edges = append(*edges, [2]*immCandidate{cand, bCand})
					arg := ast.Expr(bId)
					if stripParens(unwrapCheckNotNullArg(sel.X, info)) != bId {
						arg = bCall
					}
					return immProducerSpec{kind: immInitFromBuilder, rawExpr: expr, argExpr: arg, otherCand: bCand}, false, true
				}
			}
		}
	}

	// Candidate `.SubList(from, to)` / `.Reverse()` / `.AsList()` feeding a candidate directly!
	if mCall, ok := unwrapped.(*ast.CallExpr); ok && !mCall.Ellipsis.IsValid() {
		if sel, ok := stripParens(mCall.Fun).(*ast.SelectorExpr); ok {
			recvInner := unwrapCheckNotNullArg(sel.X, info)
			if srcCand := resolveImmListCandAccess(recvInner, varCands, info); srcCand != nil && srcCand != cand && !srcCand.disqualified {
				switch {
				case sel.Sel.Name == "SubList" && len(mCall.Args) == 2:
					*edges = append(*edges, [2]*immCandidate{cand, srcCand})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
					return immProducerSpec{
						kind: immInitSubListCand, rawExpr: expr, argExpr: recvInner, srcCall: mCall,
						elemExprs: []ast.Expr{mCall.Args[0], mCall.Args[1]}, otherCand: srcCand,
					}, false, true
				case sel.Sel.Name == "Reverse" && len(mCall.Args) == 0:
					*edges = append(*edges, [2]*immCandidate{cand, srcCand})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
					return immProducerSpec{kind: immInitReverseCand, rawExpr: expr, argExpr: recvInner, otherCand: srcCand}, false, true
				case sel.Sel.Name == "AsList" && len(mCall.Args) == 0:
					*edges = append(*edges, [2]*immCandidate{cand, srcCand})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
					return immProducerSpec{kind: immInitAsListCand, rawExpr: expr, argExpr: recvInner, otherCand: srcCand}, false, true
				}
			}
		}
	}

	// Same-package or cross-package getter fast path `recv.GetParameters()` / `recv.GetAlternates()`.
	if gCall, ok := unwrapped.(*ast.CallExpr); ok {
		if gInfo, recvExpr, okG := matchImmGetterFastCall(gCall, getterByMethod, getterByIfaceName, info); okG && !gInfo.fieldCand.disqualified {
			*edges = append(*edges, [2]*immCandidate{cand, gInfo.fieldCand})
			*mayBeNilDeps = append(*mayBeNilDeps, [2]*immCandidate{cand, gInfo.fieldCand})
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: gCall.Pos()})
			return immProducerSpec{
				kind: immInitGetterFastPath, rawExpr: expr, argExpr: gCall,
				otherCand: gInfo.fieldCand, getterInfo: gInfo, getterRecv: recvExpr,
			}, false, true
		}
		if extGInfo, extElemT, recvExpr, okE := matchExternalSliceGetterCall(gCall, pkg, info); okE {
			if !cand.addWriteType(extElemT, accessExpr, info) {
				return immProducerSpec{}, false, false
			}
			cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: gCall.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: gCall.Pos()})
			return immProducerSpec{
				kind: immInitGetterFastPath, rawExpr: expr, argExpr: gCall,
				getterInfo: extGInfo, getterRecv: recvExpr,
			}, true, true
		}
	}

	// Static j2gort `ImmutableList.Of_*` and `ImmutableList.CopyOf_*` producers.
	if sCall, ok := unwrapped.(*ast.CallExpr); ok && !sCall.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(sCall, info); isJ2 {
			switch {
			case fnName == "Ext_com_google_common_collect_ImmutableList_Of_0" && len(sCall.Args) == 0:
				cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: sCall.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
				return immProducerSpec{kind: immInitOf0, rawExpr: expr}, false, true

			case fnName == "Ext_com_google_common_collect_ImmutableList_Of_13" && len(sCall.Args) == 13:
				elems := make([]ast.Expr, 12)
				unwraps := make([]bool, 12)
				for i := 0; i < 12; i++ {
					ub, elemT, okE := unwrapImmElemWriteArg(sCall.Args[i], info)
					if !okE || !cand.addWriteType(elemT, accessExpr, info) {
						return immProducerSpec{}, false, false
					}
					elems[i] = sCall.Args[i]
					unwraps[i] = ub != stripParens(sCall.Args[i])
				}
				cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: sCall.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
				return immProducerSpec{
					kind: immInitOfVarargs, rawExpr: expr, argExpr: sCall.Args[12], srcCall: sCall, elemExprs: elems, elemUnwrap: unwraps,
				}, false, true

			case strings.HasPrefix(fnName, "Ext_com_google_common_collect_ImmutableList_Of_"):
				arity, err := strconv.Atoi(strings.TrimPrefix(fnName, "Ext_com_google_common_collect_ImmutableList_Of_"))
				if err == nil && arity >= 1 && arity <= 12 && len(sCall.Args) == arity {
					elems := make([]ast.Expr, arity)
					unwraps := make([]bool, arity)
					for i, a := range sCall.Args {
						ub, elemT, okE := unwrapImmElemWriteArg(a, info)
						if !okE || !cand.addWriteType(elemT, accessExpr, info) {
							return immProducerSpec{}, false, false
						}
						elems[i] = a
						unwraps[i] = ub != stripParens(a)
					}
					cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: sCall.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
					return immProducerSpec{kind: immInitOfFixed, rawExpr: expr, srcCall: sCall, elemExprs: elems, elemUnwrap: unwraps}, false, true
				}

			case (fnName == "Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection" ||
				fnName == "Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable" ||
				fnName == "Ext_com_google_common_collect_ImmutableList_CopyOf_1_Collection" ||
				fnName == "Ext_com_google_common_collect_ImmutableList_CopyOf_1_Iterable") && len(sCall.Args) == 1:
				srcInner := unwrapCollectionAdaptersDown(sCall.Args[0], info)
				if srcCand := resolveImmListCandAccess(srcInner, varCands, info); srcCand != nil && srcCand != cand {
					*edges = append(*edges, [2]*immCandidate{cand, srcCand})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
					return immProducerSpec{kind: immInitCopyCand, rawExpr: expr, argExpr: srcInner, otherCand: srcCand}, false, true
				}
				if gCall, isCall := srcInner.(*ast.CallExpr); isCall {
					if gInfo, recvExpr, okG := matchImmGetterFastCall(gCall, getterByMethod, getterByIfaceName, info); okG && !gInfo.fieldCand.disqualified {
						*edges = append(*edges, [2]*immCandidate{cand, gInfo.fieldCand})
						cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
						return immProducerSpec{
							kind: immInitCopyGetterSlice, rawExpr: expr, argExpr: gCall,
							otherCand: gInfo.fieldCand, getterInfo: gInfo, getterRecv: recvExpr,
						}, false, true
					}
					if extGInfo, extElemT, recvExpr, okE := matchExternalSliceGetterCall(gCall, pkg, info); okE {
						if !cand.addWriteType(extElemT, accessExpr, info) {
							return immProducerSpec{}, false, false
						}
						cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: sCall.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
						cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
						return immProducerSpec{
							kind: immInitCopyGetterSlice, rawExpr: expr, argExpr: gCall,
							getterInfo: extGInfo, getterRecv: recvExpr,
						}, false, true
					}
				}
				if sliceArg, elemT, okS := matchImmListFromSliceCall(srcInner, info); okS {
					if !cand.addWriteType(elemT, accessExpr, info) {
						return immProducerSpec{}, false, false
					}
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
					return immProducerSpec{kind: immInitCopyLoweredSlice, rawExpr: expr, argExpr: sliceArg}, false, true
				}
				cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: sCall.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
				return immProducerSpec{kind: immInitCopyCol, rawExpr: expr, argExpr: sCall.Args[0], srcCall: sCall}, false, true

			case (fnName == "Ext_com_google_common_collect_ImmutableList_CopyOf_1_slice_any" ||
				fnName == "Ext_com_google_common_collect_ImmutableList_CopyOf_1_ObjectArray") && len(sCall.Args) == 1:
				cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: sCall.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
				return immProducerSpec{kind: immInitCopyArray, rawExpr: expr, argExpr: sCall.Args[0], srcCall: sCall}, false, true
			}
		}
	}

	// Arbitrary external expression of type `Ext_com_google_common_collect_ImmutableList`
	// (such as `registry.CreateParametersWithVarArgs(...)` in `FunctionTypeBase`).
	if isExactImmutableListType(info.Types[expr].Type) {
		cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: expr.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: expr.Pos()})
		return immProducerSpec{kind: immInitExternal, rawExpr: expr, argExpr: expr}, true, true
	}

	return immProducerSpec{}, false, false
}

func validateImmBuilderUse(
	id *ast.Ident,
	file *ast.File,
	cand *immCandidate,
	varCands map[*types.Var]*immCandidate,
	getterByMethod map[*types.Func]*immGetterInfo,
	getterByIfaceName map[string]*immGetterInfo,
	pkg *types.Package,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	edges *[][2]*immCandidate,
) {
	// Blank-identifier suppression `_ = builder`.
	if isBlankAssignTarget(id, parents) {
		return
	}
	curTarget := unwrapEnclosingCheckNotNulls(id, parents, info)
	for {
		sel, ok := parents[curTarget].(*ast.SelectorExpr)
		if !ok || sel.X != curTarget {
			cand.disqualified = true
			return
		}
		mCall, ok := parents[sel].(*ast.CallExpr)
		if !ok || mCall.Fun != sel || mCall.Ellipsis.IsValid() {
			cand.disqualified = true
			return
		}

		switch sel.Sel.Name {
		case "AddVoid_", "AddStrVoid_", "Add_1_any", "Add_1_Object", "AddSelf_", "AddStrSelf_":
			if len(mCall.Args) != 1 || enclosesLoopAfterDecl(mCall, cand.declBlock, parents, true) {
				cand.disqualified = true
				return
			}
			var elemT types.Type
			var unwraps bool
			if sel.Sel.Name == "AddStrVoid_" || sel.Sel.Name == "AddStrSelf_" {
				if !hasExactOrDefaultStringType(mCall.Args[0], info) {
					cand.disqualified = true
					return
				}
				elemT = types.Typ[types.String]
			} else {
				ub, et, okE := unwrapImmElemWriteArg(mCall.Args[0], info)
				if !okE {
					cand.disqualified = true
					return
				}
				elemT = et
				unwraps = ub != stripParens(mCall.Args[0])
			}
			if !cand.addWriteType(elemT, nil, info) {
				cand.disqualified = true
				return
			}
			if mCall.Pos() > cand.maxMutPos {
				cand.maxMutPos = mCall.Pos()
			}
			outer := unwrapBuilderSelfCastUp(mCall, parents, info)
			if es, isES := parents[outer].(*ast.ExprStmt); isES && canReplaceStmtIn(parents[es]) {
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
				cand.ops = append(cand.ops, immRewriteOp{
					kind: immOpBuilderAddStmt, cand: cand, file: file, spanNode: es, stmt: es,
					accessExpr: id, call: mCall, unwrapsBox: unwraps,
				})
				return
			}
			if isChainedBuilderStep(outer, parents) {
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
				cand.ops = append(cand.ops, immRewriteOp{
					kind: immOpBuilderAddExpr, cand: cand, file: file, spanNode: outer, outerExpr: outer,
					accessExpr: id, call: mCall, unwrapsBox: unwraps,
				})
				curTarget = outer
				continue
			}
			cand.disqualified = true
			return

		case "AddAllVoid_", "AddAllSelf_", "AddAll_1_Ext_java_lang_Iterable", "AddAll_1_Iterable":
			if len(mCall.Args) != 1 || enclosesLoopAfterDecl(mCall, cand.declBlock, parents, true) {
				cand.disqualified = true
				return
			}
			if mCall.Pos() > cand.maxMutPos {
				cand.maxMutPos = mCall.Pos()
			}
			var otherC *immCandidate
			var gInfo *immGetterInfo
			var gRecv ast.Expr
			srcInner := unwrapCollectionAdaptersDown(mCall.Args[0], info)
			if sc := resolveImmListCandAccess(srcInner, varCands, info); sc != nil && sc != cand {
				otherC = sc
				*edges = append(*edges, [2]*immCandidate{cand, sc})
			} else if gCall, isCall := srcInner.(*ast.CallExpr); isCall {
				if gi, rx, okG := matchImmGetterFastCall(gCall, getterByMethod, getterByIfaceName, info); okG && !gi.fieldCand.disqualified {
					gInfo = gi
					gRecv = rx
					otherC = gi.fieldCand
					*edges = append(*edges, [2]*immCandidate{cand, gi.fieldCand})
				} else if extGI, extElemT, rx, okE := matchExternalSliceGetterCall(gCall, pkg, info); okE {
					if !cand.addWriteType(extElemT, nil, info) {
						cand.disqualified = true
						return
					}
					gInfo = extGI
					gRecv = rx
				}
			}
			outer := unwrapBuilderSelfCastUp(mCall, parents, info)
			cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: mCall.Pos()})
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
			if es, isES := parents[outer].(*ast.ExprStmt); isES && canReplaceStmtIn(parents[es]) {
				cand.ops = append(cand.ops, immRewriteOp{
					kind: immOpBuilderAddAllStmt, cand: cand, otherCand: otherC,
					getterInfo: gInfo, getterRecv: gRecv, file: file, spanNode: es, stmt: es,
					accessExpr: id, argExpr: srcInner, call: mCall,
				})
				return
			}
			if isChainedBuilderStep(outer, parents) {
				cand.ops = append(cand.ops, immRewriteOp{
					kind: immOpBuilderAddAllExpr, cand: cand, otherCand: otherC,
					getterInfo: gInfo, getterRecv: gRecv, file: file, spanNode: outer, outerExpr: outer,
					accessExpr: id, argExpr: srcInner, call: mCall,
				})
				curTarget = outer
				continue
			}
			cand.disqualified = true
			return

		case "BuildImmList_", "Build":
			if len(mCall.Args) != 0 || enclosesLoopAfterDecl(mCall, cand.declBlock, parents, false) {
				cand.disqualified = true
				return
			}
			if !recordBuilderBuildSite(cand, id, mCall) {
				cand.disqualified = true
				return
			}
			unwrappedUp := unwrapImmListCastUp(mCall, parents, info)
			if isConsumedByImmCandidateInit(unwrappedUp, varCands, parents, info) {
				return
			}
			if mSel2, ok2 := parents[unwrappedUp].(*ast.SelectorExpr); ok2 && mSel2.X == unwrappedUp {
				if c2, okC := parents[mSel2].(*ast.CallExpr); okC && c2.Fun == mSel2 && len(c2.Args) == 0 && !c2.Ellipsis.IsValid() {
					if _, inES := parents[c2].(*ast.ExprStmt); !inES {
						if mSel2.Sel.Name == "Size" {
							cand.ops = append(cand.ops, immRewriteOp{
								kind: immOpSize, cand: cand, file: file, spanNode: c2, outerExpr: c2, accessExpr: id, isGuarded: true,
							})
							return
						}
						if mSel2.Sel.Name == "IsEmpty" {
							cand.ops = append(cand.ops, immRewriteOp{
								kind: immOpIsEmpty, cand: cand, file: file, spanNode: c2, outerExpr: c2, accessExpr: id, isGuarded: true,
							})
							return
						}
					}
				}
			}
			escapeOuter := unwrapEscapeAdaptersUp(unwrappedUp, parents, info)
			if _, inES := parents[escapeOuter].(*ast.ExprStmt); inES {
				cand.disqualified = true
				return
			}
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: escapeOuter.Pos()})
			cand.ops = append(cand.ops, immRewriteOp{
				kind: immOpBuilderBuildEscape, cand: cand, file: file, spanNode: escapeOuter, outerExpr: escapeOuter, accessExpr: id,
			})
			return

		default:
			cand.disqualified = true
			return
		}
	}
}

func unwrapBuilderChainRoot(expr ast.Expr, info *types.Info) *ast.Ident {
	cur := stripParens(unwrapCheckNotNullArg(expr, info))
	for {
		if id, ok := cur.(*ast.Ident); ok {
			return id
		}
		call, ok := cur.(*ast.CallExpr)
		if !ok || call.Ellipsis.IsValid() {
			return nil
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && len(call.Args) == 1 {
			switch fnName {
			case "CastPtr", "Cast", "UncheckedCast":
				cur = stripParens(unwrapCheckNotNullArg(call.Args[0], info))
				continue
			}
		}
		sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
		if !ok {
			return nil
		}
		switch sel.Sel.Name {
		case "AddSelf_", "AddStrSelf_", "Add_1_any", "Add_1_Object",
			"AddAllSelf_", "AddAll_1_Ext_java_lang_Iterable", "AddAll_1_Iterable":
			cur = stripParens(unwrapCheckNotNullArg(sel.X, info))
			continue
		}
		return nil
	}
}

func validateImmListUse(
	id *ast.Ident,
	file *ast.File,
	cand *immCandidate,
	varCands map[*types.Var]*immCandidate,
	consumedExprs map[ast.Expr]bool,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
	edges *[][2]*immCandidate,
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

	// Disqualify if mutated/assigned via unrecognized LHS or `&x` address-taken.
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

	// Equality / nil checks: `s == nil` / `s != nil` are native; `s1 == s2` (interface identity) disqualifies!
	if bin, ok := parents[pTop].(*ast.BinaryExpr); ok && (bin.Op == token.EQL || bin.Op == token.NEQ) {
		var otherSide ast.Expr
		if bin.X == pTop {
			otherSide = bin.Y
		} else if bin.Y == pTop {
			otherSide = bin.X
		}
		if otherSide != nil {
			if isNilIdent(otherSide, info) {
				cand.ops = append(cand.ops, immRewriteOp{
					kind: immOpNilCheck, cand: cand, file: file, spanNode: bin, outerExpr: bin, accessExpr: accessExpr,
				})
				return
			}
			cand.disqualified = true
			return
		}
	}
	// Also disqualify `j2gort.IfaceEquals(s, ...)` reference-identity checks.
	if eqCall, ok := parents[pTop].(*ast.CallExpr); ok && !eqCall.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(eqCall, info); isJ2 && fnName == "IfaceEquals" {
			cand.disqualified = true
			return
		}
	}

	// Note: if `accessExpr` is inside `otherBuilder.AddAll*(...)` or `otherCand.Equals(...)`,
	// emitting `immOpEscapeWrap` as a fallback when no direct consumer matches is 100% composable:
	// - If the outer `AddAll*`/`Equals` candidate IS lowered, its larger-span op replaces the outer
	//   call/stmt using the saved `op.argExpr` (`accessExpr`), discarding the inner wrapper!
	// - If the outer candidate is NOT lowered, `ImmListFromTypedSlice(accessExpr)` adapts `[]T` cleanly!

	isGuarded := isImmAccessNonNil(accessExpr, cand, parents, info)
	if validateImmConsumerChain(accessExpr, nil, nil, isGuarded, file, cand, varCands, parents, info, varUseCounts, edges) {
		return
	}

	// Fallback: escape wrap via `j2gort.ImmListFromTypedSlice(s)`!
	escapeOuter := unwrapEscapeAdaptersUp(accessExpr, parents, info)
	if _, inES := parents[escapeOuter].(*ast.ExprStmt); inES {
		cand.disqualified = true
		return
	}
	cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: escapeOuter.Pos()})
	cand.ops = append(cand.ops, immRewriteOp{
		kind: immOpEscapeWrap, cand: cand, file: file, spanNode: escapeOuter, outerExpr: escapeOuter, accessExpr: accessExpr,
	})
}

func tryAttachDirectGetterConsumer(
	getterCall *ast.CallExpr,
	gInfo *immGetterInfo,
	recvExpr ast.Expr,
	file *ast.File,
	varCands map[*types.Var]*immCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
	edges *[][2]*immCandidate,
) {
	// If `recv.GetFoo()` feeds `builder.AddAll*(AsIface[Iterable](recv.GetFoo()))`, let builder's AddAll
	// rewrite it or let ImmListFromTypedSlice view handle it; here handle direct Size/IsEmpty/Get/Iterator/IterSlice_/NewIterCursor!
	validateImmConsumerChain(getterCall, gInfo, recvExpr, false, file, gInfo.fieldCand, varCands, parents, info, varUseCounts, edges)
}

func validateImmConsumerChain(
	baseExpr ast.Expr,
	gInfo *immGetterInfo,
	getterRecv ast.Expr,
	isGuarded bool,
	file *ast.File,
	cand *immCandidate,
	varCands map[*types.Var]*immCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
	edges *[][2]*immCandidate,
) bool {
	// Site type-arg context: for direct field access use `baseExpr`, for getter call `gInfo` owner has 0 type params.
	var siteAccess ast.Expr
	if gInfo == nil {
		siteAccess = baseExpr
	}

	adaptedCol := unwrapCollectionAdapters(baseExpr, parents, info)
	if sCall, ok := parents[adaptedCol].(*ast.CallExpr); ok && !sCall.Ellipsis.IsValid() && len(sCall.Args) >= 1 && sCall.Args[0] == adaptedCol {
		if fnName, _, _, isJ2 := isJ2gortCall(sCall, info); isJ2 {
			switch {
			case fnName == "NewIterCursor" && len(sCall.Args) == 1:
				if unwrapEnclosingCheckNotNulls(baseExpr, parents, info) == adaptedCol {
					return validateImmIterCursorLoop(sCall, baseExpr, siteAccess, gInfo, getterRecv, isGuarded, file, cand, parents, info, varUseCounts)
				}
			case fnName == "Ext_com_google_common_collect_Iterables_GetLast_1" && len(sCall.Args) == 1:
				outer, elemT, okU := enclosingImmElemUnboxExpr(sCall, parents, info)
				_, isES := parents[outer].(*ast.ExprStmt)
				if okU && !isES && cand.addReadType(elemT, siteAccess, info) {
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outer.Pos()})
					cand.ops = append(cand.ops, immRewriteOp{
						kind: immOpGetLast, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
						file: file, spanNode: outer, outerExpr: outer, accessExpr: baseExpr, isGuarded: isGuarded,
					})
					return true
				}
			case fnName == "ListGetStr" && len(sCall.Args) == 2:
				_, isES := parents[sCall].(*ast.ExprStmt)
				if !isES && cand.addReadType(types.Typ[types.String], siteAccess, info) {
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
					cand.ops = append(cand.ops, immRewriteOp{
						kind: immOpListGetStr, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
						file: file, spanNode: sCall, outerExpr: sCall, accessExpr: baseExpr, argExpr: sCall.Args[1], call: sCall, isGuarded: isGuarded,
					})
					return true
				}
			case fnName == "SetContainsStr" && len(sCall.Args) == 2:
				_, isES := parents[sCall].(*ast.ExprStmt)
				if !isES && cand.addReadType(types.Typ[types.String], siteAccess, info) {
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: sCall.Pos()})
					cand.ops = append(cand.ops, immRewriteOp{
						kind: immOpContains, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
						file: file, spanNode: sCall, outerExpr: sCall, accessExpr: baseExpr, argExpr: sCall.Args[1], call: sCall, isGuarded: isGuarded,
					})
					return true
				}
			}
		}
	}

	topExpr := unwrapEnclosingCheckNotNulls(baseExpr, parents, info)
	mSel, ok := parents[topExpr].(*ast.SelectorExpr)
	if !ok || mSel.X != topExpr {
		return false
	}
	mCall, ok := parents[mSel].(*ast.CallExpr)
	if !ok || mCall.Fun != mSel || mCall.Ellipsis.IsValid() {
		return false
	}
	_, inES := parents[mCall].(*ast.ExprStmt)

	switch mSel.Sel.Name {
	case "Size":
		if inES || len(mCall.Args) != 0 {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpSize, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: mCall, outerExpr: mCall, accessExpr: baseExpr, isGuarded: isGuarded,
		})
		return true

	case "IsEmpty":
		if inES || len(mCall.Args) != 0 {
			return false
		}
		outer := ast.Expr(mCall)
		neg := false
		pUp := unwrapParensUp(mCall, parents)
		if un, isNot := parents[pUp].(*ast.UnaryExpr); isNot && un.Op == token.NOT && un.X == pUp {
			if _, uES := parents[un].(*ast.ExprStmt); !uES {
				outer = un
				neg = true
			}
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outer.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpIsEmpty, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: outer, outerExpr: outer, accessExpr: baseExpr, isGuarded: isGuarded, isNegated: neg,
		})
		return true

	case "Get":
		if inES || len(mCall.Args) != 1 {
			return false
		}
		outer, elemT, okU := enclosingImmElemUnboxExpr(mCall, parents, info)
		if !okU {
			return false
		}
		if _, outerES := parents[outer].(*ast.ExprStmt); outerES {
			return false
		}
		if !cand.addReadType(elemT, siteAccess, info) {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outer.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpGet, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: outer, outerExpr: outer, accessExpr: baseExpr, argExpr: mCall.Args[0], call: mCall, isGuarded: isGuarded,
		})
		return true

	case "ListGetStr_":
		if inES || len(mCall.Args) != 1 || !cand.addReadType(types.Typ[types.String], siteAccess, info) {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpListGetStr, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: mCall, outerExpr: mCall, accessExpr: baseExpr, argExpr: mCall.Args[0], call: mCall, isGuarded: isGuarded,
		})
		return true

	case "GetFirst", "GetLast":
		if inES || len(mCall.Args) != 0 {
			return false
		}
		outer, elemT, okU := enclosingImmElemUnboxExpr(mCall, parents, info)
		if !okU {
			return false
		}
		if _, outerES := parents[outer].(*ast.ExprStmt); outerES {
			return false
		}
		if !cand.addReadType(elemT, siteAccess, info) {
			return false
		}
		opK := immOpGetFirst
		if mSel.Sel.Name == "GetLast" {
			opK = immOpGetLast
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outer.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: opK, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: outer, outerExpr: outer, accessExpr: baseExpr, isGuarded: isGuarded,
		})
		return true

	case "Contains", "SetContainsStr_":
		if inES || len(mCall.Args) != 1 {
			return false
		}
		if mSel.Sel.Name == "SetContainsStr_" && !cand.addReadType(types.Typ[types.String], siteAccess, info) {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpContains, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: mCall, outerExpr: mCall, accessExpr: baseExpr, argExpr: mCall.Args[0], call: mCall, isGuarded: isGuarded,
		})
		return true

	case "IndexOf_1", "IndexOf":
		if inES || len(mCall.Args) != 1 {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpIndexOf, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: mCall, outerExpr: mCall, accessExpr: baseExpr, argExpr: mCall.Args[0], call: mCall,
		})
		return true

	case "LastIndexOf_1", "LastIndexOf":
		if inES || len(mCall.Args) != 1 {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpLastIndexOf, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: mCall, outerExpr: mCall, accessExpr: baseExpr, argExpr: mCall.Args[0], call: mCall,
		})
		return true

	case "SubList":
		if inES || len(mCall.Args) != 2 {
			return false
		}
		outer := unwrapEscapeAdaptersUp(unwrapImmListCastUp(mCall, parents, info), parents, info)
		if _, outerES := parents[outer].(*ast.ExprStmt); outerES {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outer.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpSubList, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: outer, outerExpr: outer,
			accessExpr: baseExpr, argExpr: mCall.Args[0], argExpr2: mCall.Args[1], call: mCall,
		})
		return true

	case "Reverse":
		if inES || len(mCall.Args) != 0 {
			return false
		}
		outer := unwrapEscapeAdaptersUp(unwrapImmListCastUp(mCall, parents, info), parents, info)
		if _, outerES := parents[outer].(*ast.ExprStmt); outerES {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outer.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpReverse, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: outer, outerExpr: outer, accessExpr: baseExpr,
		})
		return true

	case "AsList":
		if inES || len(mCall.Args) != 0 {
			return false
		}
		outer := unwrapEscapeAdaptersUp(unwrapImmListCastUp(mCall, parents, info), parents, info)
		if _, outerES := parents[outer].(*ast.ExprStmt); outerES {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outer.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpAsList, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: outer, outerExpr: outer, accessExpr: baseExpr, isGuarded: isGuarded,
		})
		return true

	case "Equals":
		if inES || len(mCall.Args) != 1 {
			return false
		}
		var otherC *immCandidate
		otherInner := unwrapCollectionAdaptersDown(mCall.Args[0], info)
		if oc := resolveImmListCandAccess(otherInner, varCands, info); oc != nil && oc != cand && !oc.disqualified {
			otherC = oc
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpEquals, cand: cand, otherCand: otherC, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: baseExpr, argExpr: otherInner, call: mCall,
		})
		return true

	case "HashCode":
		if inES || len(mCall.Args) != 0 {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpHashCode, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: mCall, outerExpr: mCall, accessExpr: baseExpr,
		})
		return true

	case "Iterator", "Iterator_0":
		if inES || len(mCall.Args) != 0 {
			return false
		}
		outer := unwrapIteratorCast(mCall, parents, info)
		if _, outerES := parents[outer].(*ast.ExprStmt); outerES {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outer.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpIterator, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: outer, outerExpr: outer, accessExpr: baseExpr, isGuarded: isGuarded,
		})
		return true

	case "IterSlice_":
		if len(mCall.Args) != 0 {
			return false
		}
		rs, okR := parents[mCall].(*ast.RangeStmt)
		if !okR || rs.X != mCall || rs.Tok != token.DEFINE || rs.Value == nil || rs.Body == nil {
			return false
		}
		if rs.Key != nil {
			if kId, okK := rs.Key.(*ast.Ident); !okK || kId.Name != "_" {
				return false
			}
		}
		valId, okV := rs.Value.(*ast.Ident)
		if !okV || valId.Name == "_" {
			return false
		}
		rawVar := varFromObj(info.Defs[valId])
		if rawVar == nil {
			return false
		}
		elemId, drop, okL := validateImmLoopRawElem(rawVar, siteAccess, file, rs.Body.List, cand, parents, info, varUseCounts)
		if !okL {
			return false
		}
		if elemId == nil {
			elemId = valId
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: rs.Pos()})
		cand.ops = append(cand.ops, immRewriteOp{
			kind: immOpRangeIterSlice, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
			file: file, spanNode: rs, rangeStmt: rs, outerExpr: mCall, accessExpr: baseExpr,
			elemIdent: elemId, dropStmts: drop, isGuarded: isGuarded,
		})
		return true
	}

	return false
}

func validateImmIterCursorLoop(
	cursorCall *ast.CallExpr,
	baseExpr ast.Expr,
	siteAccess ast.Expr,
	gInfo *immGetterInfo,
	getterRecv ast.Expr,
	isGuarded bool,
	file *ast.File,
	cand *immCandidate,
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
		if !cand.addReadType(types.Typ[types.String], siteAccess, info) {
			return false
		}
	case isCursorMethodCall(rhs0, curVar, "NextNode", info):
		if !cand.addReadType(types.Typ[types.Uint32], siteAccess, info) {
			return false
		}
	case isNodeConvOfCursorNextNode(rhs0, curVar, info):
		if !cand.addReadType(info.Types[rhs0].Type, siteAccess, info) {
			return false
		}
	case isCursorMethodCall(rhs0, curVar, "Next", info):
		rawVar := varFromObj(info.Defs[lhs0])
		if rawVar == nil {
			return false
		}
		unboxedId, extraDrop, okL := validateImmLoopRawElem(rawVar, siteAccess, file, forStmt.Body.List[1:], cand, parents, info, varUseCounts)
		if !okL {
			return false
		}
		if unboxedId != nil {
			elemId = unboxedId
			drop += extraDrop
		}
	default:
		return false
	}

	cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: forStmt.Pos()})
	cand.ops = append(cand.ops, immRewriteOp{
		kind: immOpIterCursor, cand: cand, getterInfo: gInfo, getterRecv: getterRecv,
		file: file, spanNode: forStmt, stmt: initAs, forStmt: forStmt, outerExpr: cursorCall,
		accessExpr: baseExpr, elemIdent: elemId, dropStmts: drop, isGuarded: isGuarded,
	})
	return true
}

func validateImmLoopRawElem(
	rawVar *types.Var,
	siteAccess ast.Expr,
	file *ast.File,
	stmts []ast.Stmt,
	cand *immCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
) (*ast.Ident, int, bool) {
	if len(stmts) > 0 && varUseCounts[rawVar] == 1 {
		if unboxedId, elemT, ok := matchSingleImmDefineUnbox(stmts[0], rawVar, parents, info); ok {
			if !cand.addReadType(elemT, siteAccess, info) {
				return nil, 0, false
			}
			return unboxedId, 1, true
		}
	}
	if varUseCounts[rawVar] == 0 {
		return nil, 0, true
	}
	var ops []immRewriteOp
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
			outer, expectedT, okU := enclosingImmElemUnboxExpr(id, parents, info)
			_, isES := parents[outer].(*ast.ExprStmt)
			if !okU || isES || !cand.addReadType(expectedT, siteAccess, info) {
				valid = false
				return false
			}
			if outer != id {
				ops = append(ops, immRewriteOp{kind: immOpRawUnbox, cand: cand, file: file, spanNode: outer, outerExpr: outer, elemIdent: id})
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

func applyImmRewrites(
	ctx *PassContext,
	eligible []*immCandidate,
	getters []*immGetterInfo,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
) int {
	var exprOps, stmtOps []immRewriteOp
	for _, cand := range eligible {
		for _, op := range cand.ops {
			switch op.kind {
			case immOpDeclField, immOpDeclVarSpec, immOpDeclShortVar, immOpAssignInit,
				immOpRangeIterSlice, immOpIterCursor, immOpBuilderAddStmt, immOpBuilderAddAllStmt:
				stmtOps = append(stmtOps, op)
			default:
				exprOps = append(exprOps, op)
			}
		}
	}

	spanOrder := func(a, b immRewriteOp) int {
		if c := cmp.Compare(a.spanNode.End()-a.spanNode.Pos(), b.spanNode.End()-b.spanNode.Pos()); c != 0 {
			return c
		}
		return cmp.Compare(a.spanNode.Pos(), b.spanNode.Pos())
	}
	slices.SortFunc(exprOps, spanOrder)
	slices.SortFunc(stmtOps, spanOrder)

	count := 0
	for _, op := range exprOps {
		if applySingleImmRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}
	for _, op := range stmtOps {
		if applySingleImmRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}

	// Emit synthesized `Slice_<Getter>() []T` methods on eligible structs (and their polymorphic interfaces)!
	for _, g := range getters {
		elemExpr, ok := typeExprForFile(g.fieldCand.elemType, g.fnFile, ctx.Pkg, ctx.Importer)
		if !ok {
			continue
		}
		recvName := "this_"
		if len(g.fnDecl.Recv.List[0].Names) == 1 && g.fnDecl.Recv.List[0].Names[0].Name != "" {
			recvName = g.fnDecl.Recv.List[0].Names[0].Name
		}
		retSliceStmt := &ast.ReturnStmt{Results: []ast.Expr{
			&ast.SelectorExpr{
				X:   &ast.Ident{Name: recvName},
				Sel: &ast.Ident{Name: g.fieldCand.v.Name()},
			},
		}}
		var bodyStmts []ast.Stmt
		if len(g.fnDecl.Body.List) > 1 {
			prefix := slices.Clone(g.fnDecl.Body.List[:len(g.fnDecl.Body.List)-1])
			bodyStmts = append(prefix, retSliceStmt)
			q, okQ := j2gortQualifierForFile(g.fnFile, ctx)
			if !okQ {
				continue
			}
			g.fnDecl.Body.List = []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: q},
							Sel: &ast.Ident{Name: "ImmListFromTypedSlice"},
						},
						Args: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.Ident{Name: recvName},
									Sel: &ast.Ident{Name: g.sliceHelper},
								},
							},
						},
					},
				}},
			}
		} else {
			bodyStmts = []ast.Stmt{retSliceStmt}
		}
		helperDecl := &ast.FuncDecl{
			Recv: &ast.FieldList{List: []*ast.Field{{
				Names: []*ast.Ident{{Name: recvName}},
				Type:  &ast.StarExpr{X: &ast.Ident{Name: g.ownerStruct.Obj().Name()}},
			}}},
			Name: &ast.Ident{Name: g.sliceHelper},
			Type: &ast.FuncType{
				Params:  &ast.FieldList{},
				Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.ArrayType{Elt: elemExpr}}}},
			},
			Body: &ast.BlockStmt{List: bodyStmts},
		}
		g.fnFile.Decls = append(g.fnFile.Decls, helperDecl)
		ctx.MarkFileModified(g.fnFile)
		count++

		if g.polyIfaceAST != nil && g.polyIfaceFile != nil {
			if ifaceElemExpr, okI := typeExprForFile(g.fieldCand.elemType, g.polyIfaceFile, ctx.Pkg, ctx.Importer); okI {
				g.polyIfaceAST.Methods.List = append(g.polyIfaceAST.Methods.List, &ast.Field{
					Names: []*ast.Ident{{Name: g.sliceHelper}},
					Type: &ast.FuncType{
						Params:  &ast.FieldList{},
						Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.ArrayType{Elt: ifaceElemExpr}}}},
					},
				})
				ctx.MarkFileModified(g.polyIfaceFile)
				count++
			}
		}
	}

	return count
}

func applySingleImmRewrite(ctx *PassContext, op immRewriteOp, parents map[ast.Node]ast.Node) bool {
	j2Sel := func(name string, pos token.Pos) *ast.SelectorExpr {
		q, _ := j2gortQualifierForFile(op.file, ctx)
		return &ast.SelectorExpr{X: &ast.Ident{NamePos: pos, Name: q}, Sel: &ast.Ident{NamePos: pos, Name: name}}
	}
	var siteAccess ast.Expr
	if op.getterInfo == nil {
		siteAccess = op.accessExpr
	}
	siteElemT := instantiateOwnerTypeAtExpr(op.cand.elemType, op.cand.ownerStruct, siteAccess, ctx.Info)
	makeElemType := func() ast.Expr {
		et, _ := typeExprForFile(siteElemT, op.file, ctx.Pkg, ctx.Importer)
		return et
	}
	makeSliceType := func(pos token.Pos) ast.Expr {
		return &ast.ArrayType{Lbrack: pos, Elt: makeElemType()}
	}

	buildGetterSliceCall := func(gInfo *immGetterInfo, recv ast.Expr, pos token.Pos) ast.Expr {
		recvT := ctx.Info.Types[recv].Type
		target := recv
		nonNil := isImmRecvNonNil(recv, parents, ctx.Info)
		if recvT != nil && isNonEmptyInterfaceType(recvT) {
			checked := recv
			if !nonNil {
				fnName := "CheckNotNull"
				if hasZeroArgBoolMethod(recvT, ctx.Pkg, "IsJavaNil_") {
					fnName = "CheckNotNullIface"
				}
				checked = &ast.CallExpr{Fun: j2Sel(fnName, pos), Args: []ast.Expr{recv}}
			}
			if gInfo.directOnIface {
				target = checked
			} else if gInfo.baseGetter != "" {
				target = &ast.CallExpr{
					Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(checked), Sel: &ast.Ident{NamePos: pos, Name: gInfo.baseGetter}},
				}
			}
		} else if recvT != nil && !nonNil {
			if _, isPtr := types.Unalias(recvT).(*types.Pointer); isPtr {
				target = &ast.CallExpr{Fun: j2Sel("CheckNotNullPtr", pos), Args: []ast.Expr{recv}}
			}
		}
		return &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(target), Sel: &ast.Ident{NamePos: pos, Name: gInfo.sliceHelper}},
		}
	}

	sliceTarget := func(pos token.Pos) ast.Expr {
		if op.getterInfo != nil {
			return buildGetterSliceCall(op.getterInfo, op.getterRecv, pos)
		}
		return op.accessExpr
	}
	needsNilGuard := op.cand.mayBeNil && !op.isGuarded
	guardedSliceTarget := func(pos token.Pos) ast.Expr {
		st := sliceTarget(pos)
		if needsNilGuard {
			return &ast.CallExpr{Fun: j2Sel("CheckNotNullSlice", pos), Args: []ast.Expr{st}}
		}
		return st
	}

	wrapCheckedWriteElem := func(rawArg ast.Expr, unwraps bool, pos token.Pos) ast.Expr {
		ub := unwrapBoxArgExpr(rawArg, unwraps)
		if isGuaranteedNonNullImmElem(ub, siteElemT, ctx.Info) {
			return ub
		}
		if isStringType(siteElemT) {
			return &ast.CallExpr{Fun: j2Sel("CheckNotNullString", pos), Args: []ast.Expr{ub}}
		}
		return &ast.CallExpr{Fun: j2Sel("CheckNotNullElem", pos), Args: []ast.Expr{ub}}
	}

	var buildProducerExpr func(p immProducerSpec, liveRhs ast.Expr, inHeader bool) ast.Expr
	buildProducerExpr = func(p immProducerSpec, liveRhs ast.Expr, inHeader bool) ast.Expr {
		pos := token.NoPos
		if p.rawExpr != nil {
			pos = p.rawExpr.Pos()
		}
		wrapComp := func(lit *ast.CompositeLit) ast.Expr {
			if inHeader {
				return &ast.ParenExpr{X: lit}
			}
			return lit
		}
		switch p.kind {
		case immInitNil:
			return &ast.Ident{NamePos: pos, Name: "nil"}
		case immInitBuilderEmpty, immInitOf0:
			return wrapComp(&ast.CompositeLit{Type: makeSliceType(pos)})
		case immInitBuilderCap:
			if isNonNegativeConstInt(p.argExpr, ctx.Info) {
				return &ast.CallExpr{
					Fun:  &ast.Ident{NamePos: pos, Name: "make"},
					Args: []ast.Expr{makeSliceType(pos), &ast.BasicLit{ValuePos: pos, Kind: token.INT, Value: "0"}, p.argExpr},
				}
			}
			return &ast.CallExpr{
				Fun:  &ast.IndexExpr{X: j2Sel("MakeNonNilSlice", pos), Index: makeElemType()},
				Args: []ast.Expr{p.argExpr},
			}
		case immInitOfFixed:
			elts := make([]ast.Expr, len(p.elemExprs))
			for i, raw := range p.elemExprs {
				unw := p.elemUnwrap[i]
				if p.srcCall != nil && i < len(p.srcCall.Args) {
					if p.srcCall.Args[i] != raw {
						unw = false
					}
					raw = p.srcCall.Args[i]
				}
				elts[i] = wrapCheckedWriteElem(raw, unw, raw.Pos())
			}
			return wrapComp(&ast.CompositeLit{Type: makeSliceType(pos), Elts: elts})
		case immInitOfVarargs:
			elts := make([]ast.Expr, len(p.elemExprs))
			for i, raw := range p.elemExprs {
				unw := p.elemUnwrap[i]
				if p.srcCall != nil && i < len(p.srcCall.Args) {
					if p.srcCall.Args[i] != raw {
						unw = false
					}
					raw = p.srcCall.Args[i]
				}
				elts[i] = wrapCheckedWriteElem(raw, unw, raw.Pos())
			}
			tailArg := p.argExpr
			if p.srcCall != nil && len(p.srcCall.Args) > 12 {
				tailArg = p.srcCall.Args[12]
			}
			prefixLit := &ast.CompositeLit{Type: makeSliceType(pos), Elts: elts}
			return &ast.CallExpr{
				Fun:  &ast.IndexExpr{X: j2Sel("ImmListOfVarargsToSlice", pos), Index: makeElemType()},
				Args: []ast.Expr{prefixLit, tailArg},
			}
		case immInitCopyCol:
			fnName := "ImmListCopyToSlice"
			if isNamedUint32NodeType(siteElemT) {
				fnName = "ImmListCopyNodeToSlice"
			}
			colArg := p.argExpr
			if p.srcCall != nil && len(p.srcCall.Args) == 1 {
				colArg = p.srcCall.Args[0]
			}
			return &ast.CallExpr{
				Fun:  &ast.IndexExpr{X: j2Sel(fnName, pos), Index: makeElemType()},
				Args: []ast.Expr{colArg},
			}
		case immInitCopyLoweredSlice:
			return &ast.CallExpr{Fun: j2Sel("ImmSliceCloneNonNull", pos), Args: []ast.Expr{p.argExpr}}
		case immInitCopyGetterSlice:
			src := buildGetterSliceCall(p.getterInfo, p.getterRecv, pos)
			if p.getterInfo.isExternal || (p.otherCand != nil && p.otherCand.mayBeNil) {
				return &ast.CallExpr{Fun: j2Sel("CheckNotNullSlice", pos), Args: []ast.Expr{src}}
			}
			return src
		case immInitCopyCand, immInitAsListCand:
			if p.otherCand != nil && p.otherCand.mayBeNil {
				return &ast.CallExpr{Fun: j2Sel("CheckNotNullSlice", pos), Args: []ast.Expr{p.argExpr}}
			}
			return p.argExpr
		case immInitSubListCand:
			a0, a1 := p.elemExprs[0], p.elemExprs[1]
			if p.srcCall != nil && len(p.srcCall.Args) == 2 {
				a0, a1 = p.srcCall.Args[0], p.srcCall.Args[1]
			}
			return &ast.CallExpr{
				Fun:  j2Sel("ImmSliceSubList", pos),
				Args: []ast.Expr{p.argExpr, a0, a1},
			}
		case immInitReverseCand:
			return &ast.CallExpr{Fun: j2Sel("ImmSliceReverse", pos), Args: []ast.Expr{p.argExpr}}
		case immInitCopyArray:
			arrArg := p.argExpr
			if p.srcCall != nil && len(p.srcCall.Args) == 1 {
				arrArg = p.srcCall.Args[0]
			}
			return &ast.CallExpr{
				Fun:  &ast.IndexExpr{X: j2Sel("ImmListCopyArrayToSlice", pos), Index: makeElemType()},
				Args: []ast.Expr{arrArg},
			}
		case immInitFromCand:
			return p.argExpr
		case immInitFromBuilder:
			if bCall, isCall := stripParens(p.argExpr).(*ast.CallExpr); isCall {
				return builderBuildSliceArg(bCall, p.otherCand, pos)
			}
			return p.argExpr
		case immInitGetterFastPath:
			return buildGetterSliceCall(p.getterInfo, p.getterRecv, pos)
		case immInitExternal:
			extArg := p.argExpr
			if liveRhs != nil {
				extArg = liveRhs
			}
			fnName := "ImmListToSliceNullable"
			if isNamedUint32NodeType(siteElemT) {
				fnName = "ImmListToNodeSliceNullable"
			}
			return &ast.CallExpr{
				Fun:  &ast.IndexExpr{X: j2Sel(fnName, pos), Index: makeElemType()},
				Args: []ast.Expr{extArg},
			}
		case immInitIIFE:
			p.iifeFl.Type.Results.List[0].Type = makeSliceType(p.iifeFl.Type.Results.Pos())
			for _, br := range p.iifeBrs {
				br.ret.Results[0] = buildProducerExpr(br.spec, br.ret.Results[0], false)
			}
			return p.rawExpr
		}
		return p.rawExpr
	}

	liveCallArg := func(fallback ast.Expr) ast.Expr {
		if op.call != nil && len(op.call.Args) > 0 {
			return op.call.Args[len(op.call.Args)-1]
		}
		return fallback
	}

	switch op.kind {
	case immOpDeclField:
		fld := op.spanNode.(*ast.Field)
		fld.Type = makeSliceType(fld.Pos())
		return true

	case immOpDeclVarSpec:
		spec := op.spanNode.(*ast.ValueSpec)
		spec.Type = makeSliceType(spec.Pos())
		if len(spec.Values) == 1 {
			if op.prod.kind == immInitNil {
				spec.Values = nil
			} else {
				spec.Values[0] = buildProducerExpr(op.prod, spec.Values[0], false)
			}
		}
		return true

	case immOpDeclShortVar:
		as := op.spanNode.(*ast.AssignStmt)
		if op.prod.kind == immInitNil {
			as.Rhs[0] = &ast.CallExpr{Fun: &ast.ParenExpr{X: makeSliceType(as.Pos())}, Args: []ast.Expr{&ast.Ident{Name: "nil"}}}
		} else {
			as.Rhs[0] = buildProducerExpr(op.prod, as.Rhs[0], inControlHeader(as, parents))
		}
		return true

	case immOpAssignInit:
		if op.prod.kind == immInitNil {
			return false
		}
		as := op.spanNode.(*ast.AssignStmt)
		as.Rhs[0] = buildProducerExpr(op.prod, as.Rhs[0], inControlHeader(as, parents))
		return true

	case immOpNilCheck:
		return false

	case immOpRawUnbox:
		parent := parents[op.outerExpr]
		return parent != nil && replaceExprInParent(parent, op.outerExpr, &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name}, parents)

	case immOpSize:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		var repl ast.Expr
		if needsNilGuard {
			repl = &ast.CallExpr{Fun: j2Sel("ImmSliceSize", pos), Args: []ast.Expr{sliceTarget(pos)}}
		} else {
			repl = &ast.CallExpr{
				Fun:  &ast.Ident{NamePos: pos, Name: "int32"},
				Args: []ast.Expr{&ast.CallExpr{Fun: &ast.Ident{NamePos: pos, Name: "len"}, Args: []ast.Expr{sliceTarget(pos)}}},
			}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpIsEmpty:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		var repl ast.Expr
		if needsNilGuard {
			repl = &ast.CallExpr{Fun: j2Sel("ImmSliceIsEmpty", pos), Args: []ast.Expr{sliceTarget(pos)}}
			if op.isNegated {
				repl = &ast.UnaryExpr{OpPos: pos, Op: token.NOT, X: repl}
			}
		} else {
			tok := token.EQL
			if op.isNegated {
				tok = token.NEQ
			}
			bin := &ast.BinaryExpr{
				X:     &ast.CallExpr{Fun: &ast.Ident{NamePos: pos, Name: "len"}, Args: []ast.Expr{sliceTarget(pos)}},
				OpPos: pos,
				Op:    tok,
				Y:     &ast.BasicLit{ValuePos: pos, Kind: token.INT, Value: "0"},
			}
			repl = bin
			if binaryNeedsParens(parent, tok) {
				repl = &ast.ParenExpr{X: bin}
			}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpGet, immOpListGetStr:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.IndexExpr{
			X:      wrapPrimaryTarget(guardedSliceTarget(pos)),
			Lbrack: pos,
			Index:  liveCallArg(op.argExpr),
			Rbrack: pos,
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpGetFirst, immOpGetLast:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		fnName := "SliceGetFirst"
		if op.kind == immOpGetLast {
			fnName = "SliceGetLast"
		}
		repl := &ast.CallExpr{Fun: j2Sel(fnName, pos), Args: []ast.Expr{sliceTarget(pos)}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpContains:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{Fun: j2Sel("SliceContains", pos), Args: []ast.Expr{sliceTarget(pos), liveCallArg(op.argExpr)}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpIndexOf:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{Fun: j2Sel("SliceIndexOf", pos), Args: []ast.Expr{sliceTarget(pos), liveCallArg(op.argExpr)}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpLastIndexOf:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{Fun: j2Sel("SliceLastIndexOf", pos), Args: []ast.Expr{sliceTarget(pos), liveCallArg(op.argExpr)}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpSubList:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		a0, a1 := op.argExpr, op.argExpr2
		if op.call != nil && len(op.call.Args) == 2 {
			a0, a1 = op.call.Args[0], op.call.Args[1]
		}
		sub := &ast.CallExpr{Fun: j2Sel("ImmSliceSubList", pos), Args: []ast.Expr{sliceTarget(pos), a0, a1}}
		repl := &ast.CallExpr{Fun: j2Sel("ImmListFromTypedSlice", pos), Args: []ast.Expr{sub}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpReverse:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		rev := &ast.CallExpr{Fun: j2Sel("ImmSliceReverse", pos), Args: []ast.Expr{sliceTarget(pos)}}
		repl := &ast.CallExpr{Fun: j2Sel("ImmListFromTypedSlice", pos), Args: []ast.Expr{rev}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpAsList:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{Fun: j2Sel("ImmListFromTypedSlice", pos), Args: []ast.Expr{guardedSliceTarget(pos)}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpEquals:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		var repl ast.Expr
		if op.otherCand != nil && op.otherCand.elemType != nil {
			otherElemT := instantiateOwnerTypeAtExpr(op.otherCand.elemType, op.otherCand.ownerStruct, op.argExpr, ctx.Info)
			if types.Identical(siteElemT, otherElemT) {
				repl = &ast.CallExpr{Fun: j2Sel("ImmSliceEqualsTyped", pos), Args: []ast.Expr{sliceTarget(pos), clonePureExpr(op.argExpr)}}
			}
		}
		if repl == nil {
			repl = &ast.CallExpr{Fun: j2Sel("ImmSliceEquals", pos), Args: []ast.Expr{sliceTarget(pos), op.call.Args[0]}}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpHashCode:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{Fun: j2Sel("ImmSliceHashCode", pos), Args: []ast.Expr{sliceTarget(pos)}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpIterator:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{Fun: j2Sel("SliceIterator", pos), Args: []ast.Expr{sliceTarget(pos)}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpRangeIterSlice:
		rs := op.rangeStmt
		pos := op.outerExpr.Pos()
		rs.X = guardedSliceTarget(pos)
		rs.Key = &ast.Ident{NamePos: rs.Pos(), Name: "_"}
		rs.Value = &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name}
		if op.dropStmts > 0 && len(rs.Body.List) >= op.dropStmts {
			rs.Body.List = rs.Body.List[op.dropStmts:]
		}
		return true

	case immOpIterCursor:
		block, ok := parents[op.stmt].(*ast.BlockStmt)
		if !ok {
			return false
		}
		idx := slices.Index(block.List, op.stmt)
		if idx < 0 || idx+1 >= len(block.List) || block.List[idx+1] != op.forStmt {
			return false
		}
		pos := op.forStmt.Pos()
		bodyStmts := op.forStmt.Body.List
		if op.dropStmts > 0 && len(bodyStmts) >= op.dropStmts {
			bodyStmts = bodyStmts[op.dropStmts:]
		}
		newRange := &ast.RangeStmt{
			For:   pos,
			Key:   &ast.Ident{NamePos: pos, Name: "_"},
			Value: &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name},
			Tok:   token.DEFINE,
			X:     guardedSliceTarget(pos),
			Body:  &ast.BlockStmt{Lbrace: op.forStmt.Body.Lbrace, List: bodyStmts, Rbrace: op.forStmt.Body.Rbrace},
		}
		parents[newRange] = block
		parents[newRange.Body] = newRange
		block.List = slices.Replace(block.List, idx, idx+2, ast.Stmt(newRange))
		return true

	case immOpBuilderAddStmt:
		es := op.stmt
		parent := parents[es]
		pos := op.call.Pos()
		elem := wrapCheckedWriteElem(op.call.Args[0], op.unwrapsBox, pos)
		lhs := &ast.Ident{NamePos: op.accessExpr.Pos(), Name: op.cand.v.Name()}
		rhs := &ast.CallExpr{
			Fun:  &ast.Ident{NamePos: pos, Name: "append"},
			Args: []ast.Expr{&ast.Ident{NamePos: pos, Name: op.cand.v.Name()}, elem},
		}
		as := &ast.AssignStmt{Lhs: []ast.Expr{lhs}, TokPos: pos, Tok: token.ASSIGN, Rhs: []ast.Expr{rhs}}
		return replaceStmtInParent(parent, es, []ast.Stmt{as}, parents)

	case immOpBuilderAddExpr:
		parent := parents[op.outerExpr]
		pos := op.call.Pos()
		elem := wrapCheckedWriteElem(op.call.Args[0], op.unwrapsBox, pos)
		addr := builderChainPtrArg(op.call, op.cand, ctx.Info, pos)
		repl := &ast.CallExpr{Fun: j2Sel("SliceBuilderAdd", pos), Args: []ast.Expr{addr, elem}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpBuilderAddAllStmt:
		es := op.stmt
		parent := parents[es]
		pos := op.call.Pos()
		lhs := &ast.Ident{NamePos: op.accessExpr.Pos(), Name: op.cand.v.Name()}
		bRef := &ast.Ident{NamePos: pos, Name: op.cand.v.Name()}
		var rhs ast.Expr
		if op.getterInfo != nil {
			src := buildGetterSliceCall(op.getterInfo, op.getterRecv, pos)
			if op.getterInfo.isExternal || (op.otherCand != nil && op.otherCand.mayBeNil) {
				src = &ast.CallExpr{Fun: j2Sel("CheckNotNullSlice", pos), Args: []ast.Expr{src}}
			}
			rhs = &ast.CallExpr{
				Fun:      &ast.Ident{NamePos: pos, Name: "append"},
				Args:     []ast.Expr{bRef, src},
				Ellipsis: pos,
			}
		} else if op.otherCand != nil && op.otherCand.elemType != nil {
			otherElemT := instantiateOwnerTypeAtExpr(op.otherCand.elemType, op.otherCand.ownerStruct, op.argExpr, ctx.Info)
			if types.Identical(siteElemT, otherElemT) {
				src := clonePureExpr(op.argExpr)
				if op.otherCand.mayBeNil {
					src = &ast.CallExpr{Fun: j2Sel("CheckNotNullSlice", pos), Args: []ast.Expr{src}}
				}
				rhs = &ast.CallExpr{
					Fun:      &ast.Ident{NamePos: pos, Name: "append"},
					Args:     []ast.Expr{bRef, src},
					Ellipsis: pos,
				}
			}
		}
		if rhs == nil {
			fnName := "AppendImmCollection"
			if isNamedUint32NodeType(siteElemT) {
				fnName = "AppendImmNodeCollection"
			}
			rhs = &ast.CallExpr{
				Fun:  &ast.IndexExpr{X: j2Sel(fnName, pos), Index: makeElemType()},
				Args: []ast.Expr{bRef, op.call.Args[0]},
			}
		}
		as := &ast.AssignStmt{Lhs: []ast.Expr{lhs}, TokPos: pos, Tok: token.ASSIGN, Rhs: []ast.Expr{rhs}}
		return replaceStmtInParent(parent, es, []ast.Stmt{as}, parents)

	case immOpBuilderAddAllExpr:
		parent := parents[op.outerExpr]
		pos := op.call.Pos()
		addr := builderChainPtrArg(op.call, op.cand, ctx.Info, pos)
		if op.getterInfo != nil {
			src := buildGetterSliceCall(op.getterInfo, op.getterRecv, pos)
			repl := &ast.CallExpr{Fun: j2Sel("SliceBuilderAddAllSlice", pos), Args: []ast.Expr{addr, src}}
			return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)
		}
		fnName := "SliceBuilderAddAll"
		if isNamedUint32NodeType(siteElemT) {
			fnName = "SliceBuilderAddAllNode"
		}
		repl := &ast.CallExpr{Fun: j2Sel(fnName, pos), Args: []ast.Expr{addr, op.call.Args[0]}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpBuilderBuildEscape:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		bArg := builderBuildSliceArg(op.outerExpr, op.cand, pos)
		repl := &ast.CallExpr{Fun: j2Sel("ImmListFromTypedSlice", pos), Args: []ast.Expr{bArg}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case immOpEscapeWrap:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{Fun: j2Sel("ImmListFromTypedSlice", pos), Args: []ast.Expr{sliceTarget(pos)}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)
	}
	return false
}

func builderChainPtrArg(call *ast.CallExpr, cand *immCandidate, info *types.Info, pos token.Pos) ast.Expr {
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if ok {
		rx := unwrapCheckNotNullArg(sel.X, info)
		if innerCall, isCall := stripParens(rx).(*ast.CallExpr); isCall {
			// Chained call `b.AddSelf_(x).AddSelf_(y)` where innerCall already returns `*[]T`!
			return innerCall
		}
	}
	return &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: &ast.Ident{NamePos: pos, Name: cand.v.Name()}}
}

func builderBuildSliceArg(outerExpr ast.Expr, cand *immCandidate, pos token.Pos) ast.Expr {
	// If `b.AddSelf_(x).BuildImmList_()` was chained, dereference `*SliceBuilderAdd(&b, x)`!
	cur := stripParens(outerExpr)
	for {
		c, ok := cur.(*ast.CallExpr)
		if !ok {
			break
		}
		if sel, ok := stripParens(c.Fun).(*ast.SelectorExpr); ok && (sel.Sel.Name == "BuildImmList_" || sel.Sel.Name == "Build") {
			if rxCall, isCall := stripParens(sel.X).(*ast.CallExpr); isCall {
				return &ast.StarExpr{Star: pos, X: rxCall}
			}
			break
		}
		if len(c.Args) == 1 {
			cur = stripParens(c.Args[0])
			continue
		}
		break
	}
	return &ast.Ident{NamePos: pos, Name: cand.v.Name()}
}

func isChainedBuilderStep(outer ast.Expr, parents map[ast.Node]ast.Node) bool {
	sel, ok := parents[outer].(*ast.SelectorExpr)
	if !ok || sel.X != outer {
		return false
	}
	call, ok := parents[sel].(*ast.CallExpr)
	if !ok || call.Fun != sel {
		return false
	}
	switch sel.Sel.Name {
	case "AddSelf_", "AddStrSelf_", "Add_1_any", "Add_1_Object",
		"AddAllSelf_", "AddAll_1_Ext_java_lang_Iterable", "AddAll_1_Iterable",
		"BuildImmList_", "Build":
		return true
	}
	return false
}

func recordBuilderBuildSite(cand *immCandidate, bId *ast.Ident, buildCall *ast.CallExpr) bool {
	_ = bId
	if !cand.firstBuildPos.IsValid() || buildCall.Pos() < cand.firstBuildPos {
		cand.firstBuildPos = buildCall.Pos()
	}
	if cand.maxMutPos.IsValid() && cand.maxMutPos > cand.firstBuildPos {
		return false
	}
	return true
}

func enclosesLoopAfterDecl(node ast.Node, declBlock *ast.BlockStmt, parents map[ast.Node]ast.Node, allowLoopBeforeBuild bool) bool {
	if allowLoopBeforeBuild {
		// Builder Add/AddAll inside a loop BEFORE `.Build()` IS valid (e.g. `FunctionType.GetAllImplementedInterfaces()`
		// and `resolveTypes`!), because `maxMutPos <= firstBuildPos` ensures all loop mutations finish before `.Build()`!
		for cur := parents[node]; cur != nil && cur != declBlock; cur = parents[cur] {
			if _, isLit := cur.(*ast.FuncLit); isLit {
				return true
			}
		}
		return false
	}
	for cur := parents[node]; cur != nil && cur != declBlock; cur = parents[cur] {
		switch cur.(type) {
		case *ast.ForStmt, *ast.RangeStmt, *ast.FuncLit:
			return true
		}
	}
	return false
}

func enclosingBlockStmt(n ast.Node, parents map[ast.Node]ast.Node) *ast.BlockStmt {
	for cur := parents[n]; cur != nil; cur = parents[cur] {
		if b, ok := cur.(*ast.BlockStmt); ok {
			return b
		}
	}
	return nil
}

func isExactImmutableListType(t types.Type) bool {
	if t == nil {
		return false
	}
	named, ok := types.Unalias(t).(*types.Named)
	return ok && named.Obj() != nil && isJ2gortPkg(named.Obj().Pkg()) && named.Obj().Name() == "Ext_com_google_common_collect_ImmutableList"
}

func isImmListBuilderPtrType(t types.Type) bool {
	if t == nil {
		return false
	}
	ptr, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := types.Unalias(ptr.Elem()).(*types.Named)
	return ok && named.Obj() != nil && isJ2gortPkg(named.Obj().Pkg()) && named.Obj().Name() == "Ext_com_google_common_collect_ImmutableList_Builder"
}

func isImmEscapeIface(t types.Type) bool {
	if t == nil {
		return false
	}
	if isEmptyInterfaceType(t) {
		return true
	}
	named, ok := types.Unalias(t).(*types.Named)
	if !ok || named.Obj() == nil || !isJ2gortPkg(named.Obj().Pkg()) {
		return false
	}
	switch named.Obj().Name() {
	case "Ext_com_google_common_collect_ImmutableList",
		"Ext_com_google_common_collect_ImmutableCollection",
		"Ext_java_util_List",
		"Ext_java_util_SequencedCollection",
		"Ext_java_util_Collection",
		"Ext_java_lang_Iterable",
		"Ext_java_util_RandomAccess",
		"Ext_java_io_Serializable":
		return true
	}
	return false
}

func classifyImmBuilderInit(expr ast.Expr, info *types.Info) (immInitKind, ast.Expr, bool) {
	e := stripParens(expr)
	if un, ok := e.(*ast.UnaryExpr); ok && un.Op == token.AND {
		if cl, ok := stripParens(un.X).(*ast.CompositeLit); ok && len(cl.Elts) == 0 && isImmListBuilderPtrType(info.Types[e].Type) {
			return immInitBuilderEmpty, nil, true
		}
	}
	call, ok := e.(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return 0, nil, false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 {
		return 0, nil, false
	}
	switch fnName {
	case "NewExt_com_google_common_collect_ImmutableList_Builder",
		"NewExt_com_google_common_collect_ImmutableList_Builder_0",
		"Ext_com_google_common_collect_ImmutableList_S_Builder",
		"Ext_com_google_common_collect_ImmutableList_Builder_0":
		if len(call.Args) == 0 {
			return immInitBuilderEmpty, nil, true
		}
	case "Ext_com_google_common_collect_ImmutableList_BuilderWithExpectedSize",
		"Ext_com_google_common_collect_ImmutableList_BuilderWithExpectedSize_1",
		"NewExt_com_google_common_collect_ImmutableList_Builder_1_int32":
		if len(call.Args) == 1 && isIntegerExpr(call.Args[0], info) {
			return immInitBuilderCap, call.Args[0], true
		}
	}
	return 0, nil, false
}

func unwrapImmElemWriteArg(arg ast.Expr, info *types.Info) (ast.Expr, types.Type, bool) {
	return unwrapBoxedAddArg(arg, info)
}

func enclosingImmElemUnboxExpr(rawExpr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) (ast.Expr, types.Type, bool) {
	return enclosingUnboxExpr(rawExpr, parents, info)
}

func matchSingleImmDefineUnbox(stmt ast.Stmt, rawVar *types.Var, parents map[ast.Node]ast.Node, info *types.Info) (*ast.Ident, types.Type, bool) {
	return matchSingleDefineUnbox(stmt, rawVar, parents, info)
}

func isGuaranteedNonNullImmElem(expr ast.Expr, elemT types.Type, info *types.Info) bool {
	if elemT == nil {
		return false
	}
	ut := types.Unalias(elemT)
	if _, isTP := ut.(*types.TypeParam); isTP {
		return false
	}
	switch u := ut.Underlying().(type) {
	case *types.Basic:
		if (u.Info() & types.IsString) != 0 {
			return isGuaranteedNonNullString(expr, info)
		}
		if (u.Info() & (types.IsBoolean | types.IsNumeric)) != 0 {
			if isNamedUint32NodeType(elemT) {
				return false
			}
			return true
		}
	case *types.Struct:
		// Value struct (e.g. FunctionType_Parameter when lowered to a value struct): check if it has IsJavaNil_.
		if obj, _, _ := types.LookupFieldOrMethod(elemT, false, nil, "IsJavaNil_"); obj != nil {
			return false
		}
		return true
	case *types.Pointer:
		return isGuaranteedNonNilAlloc(expr, info)
	case *types.Interface:
		return isGuaranteedNonNilAlloc(expr, info)
	}
	return false
}

func resolveImmListCandAccess(expr ast.Expr, varCands map[*types.Var]*immCandidate, info *types.Info) *immCandidate {
	switch x := stripParens(expr).(type) {
	case *ast.Ident:
		if c := varCands[varFromObj(info.Uses[x])]; c != nil && !c.isField && !c.isBuilder {
			return c
		}
	case *ast.SelectorExpr:
		if c := varCands[varFromObj(info.Uses[x.Sel])]; c != nil && c.isField && isIndirectStructFieldSel(x, info) {
			return c
		}
	}
	return nil
}

func matchImmListFromSliceCall(expr ast.Expr, info *types.Info) (ast.Expr, types.Type, bool) {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
		return nil, nil, false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 || (fnName != "ImmListFromSlice" && fnName != "ImmListFromTypedSlice") {
		return nil, nil, false
	}
	sl, ok := types.Unalias(info.Types[call.Args[0]].Type).(*types.Slice)
	if !ok || isInvalidOrUntyped(sl.Elem()) || isEmptyInterfaceType(sl.Elem()) {
		return nil, nil, false
	}
	return call.Args[0], sl.Elem(), true
}

func unwrapImmListCastDown(expr ast.Expr, info *types.Info) ast.Expr {
	cur := stripParens(expr)
	for {
		call, ok := cur.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
			return cur
		}
		if _, isTI := isToIfaceCall(call, info); isTI && isExactImmutableListType(info.Types[call].Type) {
			cur = stripParens(call.Args[0])
			continue
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "AsIface", "Cast", "UncheckedCast", "CheckNotNull", "CheckNotNullIface":
				if isExactImmutableListType(info.Types[call].Type) {
					cur = stripParens(call.Args[0])
					continue
				}
			}
		}
		return cur
	}
}

func unwrapImmListCastUp(expr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) ast.Expr {
	for {
		if pe, ok := parents[expr].(*ast.ParenExpr); ok && pe.X == expr {
			expr = pe
			continue
		}
		call, ok := parents[expr].(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() || call.Args[0] != expr {
			return expr
		}
		if _, isTI := isToIfaceCall(call, info); isTI && isExactImmutableListType(info.Types[call].Type) {
			expr = call
			continue
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "AsIface", "Cast", "UncheckedCast":
				if isExactImmutableListType(info.Types[call].Type) {
					expr = call
					continue
				}
			}
		}
		return expr
	}
}

func unwrapBuilderSelfCastUp(expr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) ast.Expr {
	for {
		if pe, ok := parents[expr].(*ast.ParenExpr); ok && pe.X == expr {
			expr = pe
			continue
		}
		call, ok := parents[expr].(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() || call.Args[0] != expr {
			return expr
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "CastPtr", "Cast", "UncheckedCast":
				if isImmListBuilderPtrType(info.Types[call].Type) {
					expr = call
					continue
				}
			}
		}
		return expr
	}
}

func unwrapEscapeAdaptersUp(expr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) ast.Expr {
	for {
		if pe, ok := parents[expr].(*ast.ParenExpr); ok && pe.X == expr {
			expr = pe
			continue
		}
		call, ok := parents[expr].(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() || call.Args[0] != expr {
			return expr
		}
		// Peel `AsIface[List/Collection/Iterable/ImmutableList]` since `ImmListFromTypedSlice`
		// returns `Ext_com_google_common_collect_ImmutableList`, which satisfies all of them!
		if _, isTI := isToIfaceCall(call, info); isTI && isImmEscapeIface(info.Types[call].Type) && !isEmptyInterfaceType(info.Types[call].Type) {
			expr = call
			continue
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && fnName == "AsIface" {
			if isImmEscapeIface(info.Types[call].Type) && !isEmptyInterfaceType(info.Types[call].Type) {
				expr = call
				continue
			}
		}
		return expr
	}
}

func isConsumedByImmCandidateInit(expr ast.Expr, varCands map[*types.Var]*immCandidate, parents map[ast.Node]ast.Node, info *types.Info) bool {
	cur := unwrapParensUp(expr, parents)
	if ret, isRet := parents[cur].(*ast.ReturnStmt); isRet && len(ret.Results) == 1 && ret.Results[0] == cur {
		for p := parents[ret]; p != nil; p = parents[p] {
			if fl, ok := p.(*ast.FuncLit); ok {
				flTop := unwrapParensUp(fl, parents)
				if iifeCall, ok := parents[flTop].(*ast.CallExpr); ok && iifeCall.Fun == flTop && len(iifeCall.Args) == 0 {
					cur = unwrapParensUp(iifeCall, parents)
				}
				break
			}
		}
	}
	switch p := parents[cur].(type) {
	case *ast.ValueSpec:
		if len(p.Names) == 1 && len(p.Values) == 1 && p.Values[0] == cur {
			c := varCands[varFromObj(info.Defs[p.Names[0]])]
			return c != nil && !c.isBuilder && !c.disqualified
		}
	case *ast.AssignStmt:
		if len(p.Lhs) != 1 || len(p.Rhs) != 1 || p.Rhs[0] != cur {
			return false
		}
		if p.Tok == token.DEFINE {
			id, ok := p.Lhs[0].(*ast.Ident)
			if !ok {
				return false
			}
			c := varCands[varFromObj(info.Defs[id])]
			return c != nil && !c.isBuilder && !c.disqualified
		}
		if p.Tok == token.ASSIGN {
			switch lhs := stripParens(p.Lhs[0]).(type) {
			case *ast.Ident:
				c := varCands[varFromObj(info.Uses[lhs])]
				return c != nil && !c.isBuilder && !c.disqualified
			case *ast.SelectorExpr:
				c := varCands[varFromObj(info.Uses[lhs.Sel])]
				return c != nil && c.isField && !c.disqualified
			}
		}
	}
	return false
}

func isImmAccessNonNil(accessExpr ast.Expr, cand *immCandidate, parents map[ast.Node]ast.Node, info *types.Info) bool {
	if !cand.mayBeNil {
		return true
	}
	return isGuardedNonNilVarUse(accessExpr, accessExpr, cand.v, parents, info, func(rhs ast.Expr) bool {
		return !isNilIdent(rhs, info)
	})
}

func isImmRecvNonNil(recv ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) bool {
	if recv == nil || info == nil {
		return false
	}
	if isGuaranteedNonNilAlloc(recv, info) {
		return true
	}
	id, ok := stripParens(recv).(*ast.Ident)
	if !ok {
		return false
	}
	vObj := varFromObj(info.Uses[id])
	if vObj == nil {
		return false
	}
	// Check unmutated `this_` pointer receiver of enclosing method.
	if id.Name == "this_" {
		if _, isPtr := types.Unalias(vObj.Type()).(*types.Pointer); isPtr {
			for cur := ast.Node(recv); cur != nil; cur = parents[cur] {
				if fn, isFn := cur.(*ast.FuncDecl); isFn {
					if fn.Recv != nil && len(fn.Recv.List) == 1 && len(fn.Recv.List[0].Names) == 1 &&
						varFromObj(info.Defs[fn.Recv.List[0].Names[0]]) == vObj &&
						!varAssignedInNode(fn.Body, vObj, info) {
						return true
					}
					break
				}
			}
		}
	}
	if isGuardedNonNilVarUse(recv, recv, vObj, parents, info, func(rhs ast.Expr) bool {
		return isGuaranteedNonNilAlloc(rhs, info)
	}) {
		return true
	}
	// Check preceding `j2gort.CheckNotNull*(v)` statement in an enclosing BlockStmt.
	for cur := ast.Node(recv); cur != nil; cur = parents[cur] {
		blk, isBlk := parents[cur].(*ast.BlockStmt)
		if !isBlk {
			continue
		}
		stmtNode, isStmt := cur.(ast.Stmt)
		if !isStmt {
			continue
		}
		idx := slices.Index(blk.List, stmtNode)
		for i := idx - 1; i >= 0; i-- {
			if varAssignedInNode(blk.List[i], vObj, info) {
				break
			}
			if es, isES := blk.List[i].(*ast.ExprStmt); isES {
				if c, isCall := stripParens(es.X).(*ast.CallExpr); isCall && len(c.Args) == 1 && !c.Ellipsis.IsValid() {
					if fnName, _, _, isJ2 := isJ2gortCall(c, info); isJ2 && isCheckNotNullFuncName(fnName) {
						if argId, isId := stripParens(c.Args[0]).(*ast.Ident); isId && varFromObj(info.Uses[argId]) == vObj {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func varAssignedInNode(root ast.Node, target *types.Var, info *types.Info) bool {
	if root == nil || target == nil {
		return false
	}
	assigned := false
	ast.Inspect(root, func(n ast.Node) bool {
		if assigned || n == nil {
			return false
		}
		switch s := n.(type) {
		case *ast.AssignStmt:
			for _, lhs := range s.Lhs {
				if id, ok := stripParens(lhs).(*ast.Ident); ok && varFromObj(info.Uses[id]) == target {
					assigned = true
					return false
				}
			}
		case *ast.UnaryExpr:
			if s.Op == token.AND {
				if id, ok := stripParens(s.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == target {
					assigned = true
					return false
				}
			}
		}
		return true
	})
	return assigned
}
