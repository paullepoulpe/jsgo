// Clean-room native_set_lowering pass: lowers non-escaping LinkedHashSet / HashSet locals,
// unexported struct fields (including nullable/lazy pointer fields), and Set-returning package
// functions / private methods across call stacks to j2gort.OrderedSet[T], *j2gort.OrderedSet[T],
// or stack-allocated j2gort.LocalOrderedSet[T].

package passes

import (
	"cmp"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"
)

// NativeSetLoweringPass rewrites non-escaping LinkedHashSet/HashSet locals, unexported struct
// fields, and Set-returning functions/private methods to j2gort.OrderedSet[T], *j2gort.OrderedSet[T]
// (for nullable/lazy fields and nullable locals), or j2gort.LocalOrderedSet[T] (for standalone stack locals).
type NativeSetLoweringPass struct{}

func NewNativeSetLoweringPass() *NativeSetLoweringPass { return &NativeSetLoweringPass{} }
func (p *NativeSetLoweringPass) Name() string          { return "native_set_lowering" }
func (p *NativeSetLoweringPass) Description() string {
	return "Lower non-escaping HashSet/LinkedHashSet locals, unexported fields, and Set-returning functions to OrderedSet[T] / *OrderedSet[T] / LocalOrderedSet[T]"
}

func (p *NativeSetLoweringPass) Apply(ctx *PassContext) (int, error) {
	if isJ2gortPkg(ctx.Pkg) || len(ctx.Files) == 0 {
		return 0, nil
	}
	parentsByFile := make(map[*ast.File]map[ast.Node]ast.Node, len(ctx.Files))
	for _, f := range ctx.Files {
		parentsByFile[f], _ = buildParentMap(f)
	}

	eligible := analyzeSetCandidates(ctx, parentsByFile)
	if len(eligible) == 0 {
		return 0, nil
	}
	return applySetRewrites(ctx, eligible, parentsByFile), nil
}

type setFlavor uint8

const (
	setFlavorOrdered setFlavor = iota
	setFlavorCustom
	setFlavorIdentity
)

type setInitKind uint8

const (
	setInitNil setInitKind = iota
	setInitEmpty
	setInitCap
	setInitCopySet
	setInitCopyCol
	setInitDirectCall
	setInitIIFE
)

type setOpKind uint8

const (
	setOpDeclField setOpKind = iota
	setOpDeclParam
	setOpDeclVarSpec
	setOpDeclShortVar
	setOpAssignInit
	setOpFuncDecl
	setOpFuncCallRename
	setOpReturnUnwrap
	setOpParamArg       // f(..., s, ...) -> f(..., &s, ...) / f(..., s, ...)
	setOpFormatSink     // StrVal(s), StrConcat(..., s), ToStringHelper.Add(..., s)
	setOpMethodDirect   // .Size(), .IsEmpty(), .Clear(), .HashCode()
	setOpMethodUnbox1   // .Add_1(x)->.Add(u), .Contains(x)->.Contains(u), .Remove_1_any(x)->.Remove(u)
	setOpStaticHelper   // j2gort.SetAddStr(s, k) -> s.Add(k), SetContainsNode -> s.Contains(n), etc.
	setOpGetFirst       // enclosingUnbox(s.GetFirst()) -> s.GetFirst()
	setOpAddAllSet      // s.AddAll_1(otherSet) -> s.AddAllSet(otherSet) / s.AddAllSetPtr(otherSet)
	setOpAddAllCol      // s.AddAll_1(col) -> j2gort.OrderedSetAddAll*(&s, col)
	setOpEqualsSet      // s.Equals(otherSet) -> s.EqualsSet(otherSet) / s.EqualsSetPtr(otherSet)
	setOpImmSetCopy      // ImmutableSet_CopyOf_1_*(s) -> s.ToImmutableSet()
	setOpArrayDequeCopy  // NewExt_java_util_ArrayDeque_1_Ext_java_util_Collection(s) -> j2gort.ArrayDequeFromOrderedSet*(s)
	setOpNativeDequeCopy // DequeFromCol[T](s) / NodeDequeFromCol[N](s) -> j2gort.DequeFromOrderedSet*(s) / NodeDequeFromOrderedSet*(s)
	setOpIterCursor      // cur := NewIterCursor(s); for cur.HasNext() { ... } -> cur := s.Iter(); ...
	setOpRawUnbox        // loop element unbox rewrite
)

type setIIFEBranch struct {
	ret       *ast.ReturnStmt
	initKind  setInitKind
	argExpr   ast.Expr
	otherCand *setCandidate
}

type setIIFESpec struct {
	fnLit    *ast.FuncLit
	branches []setIIFEBranch
}

type setRewriteOp struct {
	kind          setOpKind
	cand          *setCandidate
	otherCand     *setCandidate
	initKind      setInitKind
	iifeSpec      *setIIFESpec
	file          *ast.File
	spanNode      ast.Node
	accessExpr    ast.Expr
	argExpr       ast.Expr
	call          *ast.CallExpr
	outerExpr     ast.Expr
	stmt          ast.Stmt
	forStmt       *ast.ForStmt
	methodsSel    string
	elemIdent     *ast.Ident
	dropStmts     int
	unwrapsBox    bool
	needsNullWrap bool // true when value-receiver method is called on unguarded *OrderedSet[T]
}

type setCandidate struct {
	v                *types.Var  // non-nil for local / field / param
	fn               *types.Func // non-nil for Set-returning function / private method
	paramFn          *types.Func // non-nil for Set parameter of unexported function / private method
	paramIdx         int
	fnDecl           *ast.FuncDecl
	isField          bool
	isParam          bool
	isFunc           bool
	isExportedFn     bool
	isPtr            bool // lowered to *j2gort.*OrderedSet[T] (nullable/lazy field, nullable local, or param)
	useLocalSet      bool // lowered to j2gort.LocalOrderedSet[T] (standalone non-nullable stack local)
	usesOrdSetOnly   bool // uses *OrderedSetFromSet / AddAllSet / ArrayDequeFrom*OrderedSet / EqualsSet / param
	needsPtrRtHelper bool // uses New*OrderedSetPtrCap / *OrderedSetPtrFrom* / AddAllSetPtr / ArrayDequeFrom*OrderedSetPtr
	hasHashKind      bool
	hasIdentityKind  bool
	flavor           setFlavor
	ownerStruct      *types.Named
	recvTArgs        *types.TypeList
	declFile         *ast.File
	defCount         int
	inPkgCalls       int
	nonSelfCalls     int
	paramCallIdents  []*ast.Ident
	hasInit          bool
	hasNilAssign     bool
	hasNilCheck      bool
	hasElemOp        bool
	disqualified     bool
	elemType         types.Type
	elemTypes        []types.Type
	siteConstraints  []siteTypeConstraint
	spellSites       []typeSpellSite
	j2gortSites      []typeSpellSite
	loopBodies       []*ast.BlockStmt
	ops              []setRewriteOp
	parent           *setCandidate
}

func (c *setCandidate) root() *setCandidate {
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

func unionSetCandidates(a, b *setCandidate) {
	ra, rb := a.root(), b.root()
	if ra != rb {
		rb.parent = ra
	}
}

func (c *setCandidate) addElemType(t types.Type, accessExpr ast.Expr, info *types.Info) bool {
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
			c.hasElemOp = true
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
	c.elemTypes = append(c.elemTypes, t)
	c.hasElemOp = true
	return true
}

func analyzeSetCandidates(ctx *PassContext, parentsByFile map[*ast.File]map[ast.Node]ast.Node) []*setCandidate {
	info := ctx.Info
	hasPtrHelpers := j2gortHasExport(ctx, "NewOrderedSetPtr")
	hasLocalOrderedSet := j2gortHasExport(ctx, "LocalOrderedSet")

	varCands := make(map[*types.Var]*setCandidate)
	funcCands := make(map[*types.Func]*setCandidate)
	paramCandsByFn := make(map[*types.Func][]*setCandidate)
	var allCands []*setCandidate
	structFields := make(map[*types.Named][]*setCandidate)
	var edges [][2]*setCandidate

	// Pass 1a: Collect candidate Set-returning functions/private methods, Set parameters on unexported
	// functions/private methods, and unexported Set struct fields.
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
						if v := varFromObj(info.Defs[fld.Names[0]]); v != nil && v.IsField() && !v.Exported() && !v.Embedded() && isCandidateSetType(v.Type()) {
							declHash := isConcreteHashSetType(v.Type())
							cand := &setCandidate{
								v: v, isField: true, ownerStruct: named, declFile: file, hasHashKind: declHash,
								spellSites:  []typeSpellSite{{file: file, pos: fld.Pos(), owner: named}},
								j2gortSites: []typeSpellSite{{file: file, pos: fld.Pos()}},
							}
							cand.ops = append(cand.ops, setRewriteOp{kind: setOpDeclField, cand: cand, file: file, spanNode: fld})
							varCands[v] = cand
							allCands = append(allCands, cand)
							structFields[named] = append(structFields[named], cand)
						}
					}
				}

			case *ast.FuncDecl:
				if d.Body == nil {
					continue
				}
				fnObj, ok := info.Defs[d.Name].(*types.Func)
				if !ok || fnObj == nil {
					continue
				}
				sig, ok := fnObj.Type().(*types.Signature)
				if !ok {
					continue
				}
				// Collect Set parameters on unexported package functions, priv_ methods, and void non-generic exported package functions.
				canCollectSetParams := false
				if d.Type.Params != nil {
					if !fnObj.Exported() && (sig.Recv() == nil || strings.HasPrefix(fnObj.Name(), "priv_")) {
						canCollectSetParams = true
					} else if fnObj.Exported() && sig.Recv() == nil &&
						(d.Type.TypeParams == nil || len(d.Type.TypeParams.List) == 0) && isNonGenericFunc(fnObj) &&
						(sig.Results() == nil || sig.Results().Len() == 0) &&
						d.Name.Name != "main" && d.Name.Name != "init" && allFuncParamsNamed(d.Type) &&
						(ctx.Pkg == nil || ctx.Pkg.Scope() == nil || ctx.Pkg.Scope().Lookup("ordSet_"+d.Name.Name) == nil) {
						canCollectSetParams = true
					}
				}
				if canCollectSetParams {
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
								if pv := varFromObj(info.Defs[fld.Names[0]]); pv != nil && isCandidateSetType(pv.Type()) {
									pc := &setCandidate{
										v: pv, isParam: true, paramFn: fnOrig, paramIdx: argPos,
										fnDecl: d, isExportedFn: fnObj.Exported(),
										isPtr: true, hasInit: true, usesOrdSetOnly: true,
										hasHashKind: isConcreteHashSetType(pv.Type()),
										ownerStruct: recvOwner, recvTArgs: recvTArgs, declFile: file,
										spellSites:  []typeSpellSite{{file: file, pos: fld.Pos(), owner: recvOwner, recvTArgs: recvTArgs}},
										j2gortSites: []typeSpellSite{{file: file, pos: fld.Pos()}},
									}
									pc.ops = append(pc.ops, setRewriteOp{kind: setOpDeclParam, cand: pc, file: file, spanNode: fld})
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
				if (d.Type.TypeParams != nil && len(d.Type.TypeParams.List) > 0) || !isNonGenericFunc(fnObj) {
					continue
				}
				if sig.Results() == nil || sig.Results().Len() != 1 || d.Type.Results == nil || len(d.Type.Results.List) != 1 || len(d.Type.Results.List[0].Names) > 0 {
					continue
				}
				retType := sig.Results().At(0).Type()
				if !isCandidateSetReturnType(retType, fnObj.Exported()) {
					continue
				}
				if sig.Recv() != nil {
					if fnObj.Exported() || !strings.HasPrefix(fnObj.Name(), "priv_") {
						continue
					}
				} else if fnObj.Exported() {
					if d.Name.Name == "main" || d.Name.Name == "init" || (ctx.Pkg != nil && ctx.Pkg.Scope() != nil && ctx.Pkg.Scope().Lookup("ordSet_"+d.Name.Name) != nil) {
						continue
					}
					if !allFuncParamsNamed(d.Type) {
						continue
					}
				}
				cand := &setCandidate{
					fn: fnObj, fnDecl: d, isFunc: true, isExportedFn: fnObj.Exported(), declFile: file, defCount: 1, hasInit: true,
					hasHashKind: isConcreteHashSetType(retType),
					spellSites:  []typeSpellSite{{file: file, pos: d.Type.Results.Pos()}},
					j2gortSites: []typeSpellSite{{file: file, pos: d.Type.Results.Pos()}},
				}
				cand.ops = append(cand.ops, setRewriteOp{kind: setOpFuncDecl, cand: cand, file: file, spanNode: d})
				funcCands[fnObj] = cand
				allCands = append(allCands, cand)
			}
		}
	}

	// Pass 1b: Collect local Set variables initialized from NewExt_java_util_*Set, nil, IIFE, or Set-returning candidate calls.
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		addLocal := func(id *ast.Ident, initExpr ast.Expr, span ast.Node, opK setOpKind) {
			v := varFromObj(info.Defs[id])
			if id.Name == "_" || v == nil || v.IsField() || varCands[v] != nil || !isCandidateSetType(v.Type()) {
				return
			}
			var initK setInitKind
			var argExpr ast.Expr
			var otherFn *setCandidate
			var iife *setIIFESpec
			var hasNil, initIdent, initHash bool
			if initExpr == nil || isNilIdent(initExpr, info) {
				if opK != setOpDeclVarSpec {
					return
				}
				initK = setInitNil
				hasNil = true
			} else if k, a, f, idK, hsK, ok := classifySetInitExpr(initExpr, funcCands, info); ok {
				initK, argExpr, otherFn, initIdent, initHash = k, a, f, idK, hsK
			} else if spec, specHasNil, fns, idK, hsK, ok := classifySetIIFEInit(initExpr, funcCands, info); ok {
				initK = setInitIIFE
				iife = spec
				hasNil = specHasNil
				initIdent = idK
				initHash = hsK
				for _, fnCand := range fns {
					otherFn = fnCand
				}
			} else {
				return
			}
			recvOwner, recvTArgs := enclosingMethodOwnerAndTArgs(span, parents, info)
			cand := &setCandidate{
				v: v, declFile: file, hasInit: !hasNil, hasNilAssign: hasNil,
				ownerStruct: recvOwner, recvTArgs: recvTArgs,
				hasIdentityKind: initIdent,
				hasHashKind:     initHash || isConcreteHashSetType(v.Type()),
				spellSites:      []typeSpellSite{{file: file, pos: span.Pos(), owner: recvOwner, recvTArgs: recvTArgs}},
				j2gortSites:     []typeSpellSite{{file: file, pos: span.Pos()}},
			}
			if otherFn != nil {
				edges = append(edges, [2]*setCandidate{cand, otherFn})
			}
			if iife != nil {
				for _, br := range iife.branches {
					if br.otherCand != nil {
						edges = append(edges, [2]*setCandidate{cand, br.otherCand})
					}
				}
			}
			cand.ops = append(cand.ops, setRewriteOp{
				kind: opK, cand: cand, initKind: initK, iifeSpec: iife, file: file, spanNode: span, argExpr: argExpr, call: asCallExpr(initExpr),
			})
			varCands[v] = cand
			allCands = append(allCands, cand)
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
						addLocal(spec.Names[0], initExpr, spec, setOpDeclVarSpec)
					}
				}
			case *ast.AssignStmt:
				if s.Tok == token.DEFINE && len(s.Lhs) == 1 && len(s.Rhs) == 1 && enclosingFuncBody(s, parents) != nil {
					if id, ok := s.Lhs[0].(*ast.Ident); ok {
						addLocal(id, s.Rhs[0], s, setOpDeclShortVar)
					}
				}
			}
			return true
		})
	}

	if len(allCands) == 0 {
		return nil
	}

	// Pass 1c: Upgrade setInitCopyCol initializers whose source is a lowered Set var/field/call to setInitCopySet.
	upgradeCopyBranch := func(cand *setCandidate, initK *setInitKind, argExpr *ast.Expr, otherCand **setCandidate) {
		if *initK != setInitCopyCol {
			return
		}
		unwrapped := unwrapCollectionAdaptersDown(*argExpr, info)
		if srcCand := resolveSetSourceExpr(unwrapped, varCands, funcCands, info); srcCand != nil && srcCand != cand {
			*initK = setInitCopySet
			*argExpr = unwrapped
			*otherCand = srcCand
			cand.usesOrdSetOnly = true
			srcCand.usesOrdSetOnly = true
			edges = append(edges, [2]*setCandidate{cand, srcCand})
		}
	}
	for _, cand := range allCands {
		if cand.isFunc || len(cand.ops) == 0 {
			continue
		}
		op0 := &cand.ops[0]
		upgradeCopyBranch(cand, &op0.initKind, &op0.argExpr, &op0.otherCand)
		if op0.iifeSpec != nil {
			for i := range op0.iifeSpec.branches {
				br := &op0.iifeSpec.branches[i]
				upgradeCopyBranch(cand, &br.initKind, &br.argExpr, &br.otherCand)
			}
		}
	}

	if len(structFields) > 0 {
		verifySetStructInvariants(ctx, parentsByFile, varCands, funcCands, structFields, &edges)
	}

	// Validate return statements inside every candidate Set-returning function/method.
	// order-independent: validates each function candidate independently
	for _, fc := range funcCands {
		if !ctx.ShouldTransform(fc.declFile) {
			fc.disqualified = true
			continue
		}
		validateSetFuncBodyReturns(fc, varCands, funcCands, parentsByFile[fc.declFile], info, &edges)
	}

	varUseCounts := make(map[*types.Var]int, len(info.Uses))
	// order-independent: commutative per-variable integer use-count accumulation
	for _, obj := range info.Uses {
		if v := varFromObj(obj); v != nil {
			varUseCounts[v]++
		}
	}
	funcsMutating := make(map[*types.Var]map[*types.Func]bool)
	markMut := func(cand *setCandidate, node ast.Node, parents map[ast.Node]ast.Node) {
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

	// Pass 2: Validate every use of every var/field candidate AND every call to every func candidate.
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
			if cand := varCands[varFromObj(info.Uses[id])]; cand != nil && !cand.disqualified {
				if !ctx.ShouldTransform(file) {
					cand.disqualified = true
				} else {
					validateSetVarUse(id, file, cand, varCands, funcCands, paramCandsByFn, ctx, parents, varUseCounts, markMut, &edges)
				}
			}
			if fnObj, ok := info.Uses[id].(*types.Func); ok {
				fnOrig := funcOrigin(fnObj)
				if fc := funcCands[fnOrig]; fc != nil && !fc.disqualified {
					if !ctx.ShouldTransform(file) {
						fc.disqualified = true
					} else {
						validateSetFuncUse(id, file, fc, varCands, funcCands, paramCandsByFn, ctx, parents, varUseCounts, &edges)
					}
				}
				if pcs := paramCandsByFn[fnOrig]; len(pcs) > 0 {
					validateSetParamFuncCall(id, pcs, varCands, ctx.ShouldTransform(file), parents, info)
				}
			}
			return true
		})
	}

	// Finalize pointer vs value representation on var/field candidates.
	for _, cand := range allCands {
		if cand.isFunc || cand.isParam {
			continue
		}
		if !cand.hasInit || cand.hasNilAssign || cand.hasNilCheck {
			cand.isPtr = true
		}
	}

	// Check pointer constraints on ops (e.g. direct Set-returning function calls return value OrderedSet, not *OrderedSet).
	for _, cand := range allCands {
		if cand.disqualified || cand.defCount != 1 || !ctx.ShouldTransform(cand.declFile) {
			cand.disqualified = true
			continue
		}
		if (cand.isFunc && cand.isExportedFn && cand.inPkgCalls == 0) ||
			(cand.isParam && cand.inPkgCalls == 0) ||
			(cand.isParam && cand.isExportedFn && cand.nonSelfCalls == 0) {
			cand.disqualified = true
			continue
		}
		if cand.v != nil && anySetLoopMutatesCandidate(cand, funcsMutating[cand.v], info) {
			cand.disqualified = true
			continue
		}
		checkInitPtrHelper := func(k setInitKind, other *setCandidate) bool {
			switch k {
			case setInitDirectCall:
				return !cand.isPtr
			case setInitCap, setInitCopyCol:
				if cand.isPtr {
					cand.needsPtrRtHelper = true
				}
			case setInitCopySet:
				if cand.isPtr || (other != nil && other.isPtr && other.v == nil) {
					cand.needsPtrRtHelper = true
				}
				if cand.isPtr && other != nil && !other.isPtr && other.v == nil {
					return false
				}
			}
			return true
		}
		for i := range cand.ops {
			op := &cand.ops[i]
			switch op.kind {
			case setOpDeclVarSpec, setOpDeclShortVar, setOpAssignInit:
				if op.initKind == setInitIIFE && op.iifeSpec != nil {
					for _, br := range op.iifeSpec.branches {
						if !checkInitPtrHelper(br.initKind, br.otherCand) {
							cand.disqualified = true
						}
					}
				} else if !checkInitPtrHelper(op.initKind, op.otherCand) {
					cand.disqualified = true
				}
			case setOpAddAllSet, setOpEqualsSet:
				if (op.otherCand != nil && op.otherCand.isPtr) || (op.kind == setOpEqualsSet && cand.isPtr) {
					cand.needsPtrRtHelper = true
				}
			case setOpArrayDequeCopy, setOpNativeDequeCopy:
				if cand.isPtr {
					cand.needsPtrRtHelper = true
				}
			case setOpMethodDirect, setOpMethodUnbox1, setOpGetFirst, setOpImmSetCopy:
				if cand.isPtr && op.methodsSel != "Add" && op.methodsSel != "Remove" && op.methodsSel != "Clear" {
					if !isGuardedNonNilSetUse(op.outerExpr, op.accessExpr, cand, parentsByFile[op.file], funcCands, info) {
						op.needsNullWrap = true
						cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: op.file, pos: op.outerExpr.Pos()})
					}
				}
			}
		}
		if cand.needsPtrRtHelper && !hasPtrHelpers {
			cand.disqualified = true
		}
	}

	// Re-verify Set-returning functions do not return a nullable pointer candidate.
	for _, fc := range allCands {
		if !fc.isFunc || fc.disqualified {
			continue
		}
		ast.Inspect(fc.fnDecl.Body, func(n ast.Node) bool {
			if _, isLit := n.(*ast.FuncLit); isLit {
				return false
			}
			if ret, ok := n.(*ast.ReturnStmt); ok && len(ret.Results) == 1 {
				inner := unwrapCollectionAdaptersDown(ret.Results[0], info)
				if srcCand := resolveSetSourceExpr(inner, varCands, funcCands, info); srcCand == nil || srcCand.isPtr {
					fc.disqualified = true
				}
			}
			return !fc.disqualified
		})
	}

	for _, e := range edges {
		unionSetCandidates(e[0], e[1])
	}

	groups := make(map[*setCandidate][]*setCandidate)
	for _, cand := range allCands {
		r := cand.root()
		groups[r] = append(groups[r], cand)
	}

	var out []*setCandidate
	// order-independent: out is immediately sorted deterministically by token.Pos below
	for _, members := range groups {
		groupOK := true
		hasIdentity := false
		hasHash := false
		var elemTypes []types.Type
		var memberOwners []*types.Named
		for _, m := range members {
			if m.disqualified {
				groupOK = false
				break
			}
			if m.hasIdentityKind {
				hasIdentity = true
			}
			if m.hasHashKind {
				hasHash = true
			}
			if m.ownerStruct != nil {
				memberOwners = append(memberOwners, m.ownerStruct)
			}
			elemTypes = append(elemTypes, m.elemTypes...)
		}
		if !groupOK || (hasIdentity && hasHash) || (!hasIdentity && !hasHash) {
			continue
		}
		if len(elemTypes) == 0 {
			for _, m := range members {
				for _, sc := range m.siteConstraints {
					elemTypes = append(elemTypes, sc.siteType)
				}
			}
		}
		if len(elemTypes) == 0 {
			continue
		}
		elemT, unifyOK := resolveKeyOrSetElemGroupType(elemTypes, ctx.Pkg)
		if !unifyOK {
			continue
		}
		flavor, elemOK := classifySetElemType(elemT, hasIdentity, memberOwners, ctx)
		if !elemOK {
			continue
		}
		for _, m := range members {
			if typeContainsTypeParam(elemT) && m.isField && !typeOnlyUsesOwnerTypeParams(elemT, m.ownerStruct) {
				groupOK = false
				break
			}
			if flavor == setFlavorIdentity && m.isFunc && m.isExportedFn {
				if sig, ok := m.fn.Type().(*types.Signature); ok && sig.Results().Len() == 1 && !isInterfaceType(sig.Results().At(0).Type()) {
					groupOK = false
					break
				}
			}
			for _, op := range m.ops {
				if op.kind == setOpNativeDequeCopy && strings.HasPrefix(op.methodsSel, "NodeDeque") && flavor == setFlavorCustom {
					groupOK = false
				}
				if hasIdentity {
					switch op.kind {
					case setOpImmSetCopy, setOpAddAllCol:
						groupOK = false
					case setOpDeclVarSpec, setOpDeclShortVar, setOpAssignInit:
						if op.initKind == setInitCopyCol {
							groupOK = false
						}
						if op.initKind == setInitIIFE && op.iifeSpec != nil {
							for _, br := range op.iifeSpec.branches {
								if br.initKind == setInitCopyCol {
									groupOK = false
								}
							}
						}
					}
				}
			}
			if !groupOK {
				break
			}
			for _, sc := range m.siteConstraints {
				inst, ok := instantiateOwnerType(elemT, m.ownerStruct, sc.targs)
				if !ok || (!types.Identical(inst, sc.siteType) && !(isNonEmptyInterfaceType(inst) && (isNonEmptyInterfaceType(sc.siteType) || isPtrToStructType(sc.siteType)) && types.AssignableTo(sc.siteType, inst))) {
					groupOK = false
					break
				}
			}
			if !groupOK || !canSpellTypeAtAllSites(elemT, m.spellSites, ctx) || !canSpellJ2gortAtAllSites(m.j2gortSites, ctx) {
				groupOK = false
				break
			}
		}
		if !groupOK {
			continue
		}
		for _, m := range members {
			m.elemType = elemT
			m.flavor = flavor
			if flavor == setFlavorOrdered && hasLocalOrderedSet && len(members) == 1 && m.v != nil && !m.isField && !m.isParam && !m.isPtr && !m.usesOrdSetOnly {
				m.useLocalSet = true
			}
			out = append(out, m)
		}
	}

	// Ensure multi-param exported functions only lower if all of their Set parameters are eligible;
	// otherwise disqualify the partial parameters together with their connected components.
	if len(paramCandsByFn) > 0 {
		eligSet := make(map[*setCandidate]bool, len(out))
		for _, c := range out {
			eligSet[c] = true
		}
		badRoots := make(map[*setCandidate]bool)
		for _, pcs := range paramCandsByFn {
			if len(pcs) == 0 || !pcs[0].isExportedFn {
				continue
			}
			anyIn := false
			allIn := true
			for _, pc := range pcs {
				if eligSet[pc] {
					anyIn = true
				} else {
					allIn = false
				}
			}
			if anyIn && !allIn {
				for _, pc := range pcs {
					badRoots[pc.root()] = true
				}
			}
		}
		if len(badRoots) > 0 {
			filtered := out[:0]
			for _, c := range out {
				if !badRoots[c.root()] {
					filtered = append(filtered, c)
				}
			}
			out = filtered
		}
	}

	slices.SortFunc(out, func(a, b *setCandidate) int { return cmp.Compare(a.pos(), b.pos()) })
	return out
}

func (c *setCandidate) pos() token.Pos {
	if c.v != nil {
		return c.v.Pos()
	}
	return c.fn.Pos()
}

func isGuardedNonNilSetUse(
	useNode ast.Node,
	accessExpr ast.Expr,
	cand *setCandidate,
	parents map[ast.Node]ast.Node,
	funcCands map[*types.Func]*setCandidate,
	info *types.Info,
) bool {
	if cand.v == nil {
		return true
	}
	return isGuardedNonNilVarUse(useNode, accessExpr, cand.v, parents, info, func(rhs ast.Expr) bool {
		if _, _, _, _, _, ok := classifySetInitExpr(rhs, funcCands, info); ok {
			return true
		}
		if _, hasNil, _, _, _, ok := classifySetIIFEInit(rhs, funcCands, info); ok && !hasNil {
			return true
		}
		return false
	})
}

func verifySetStructInvariants(
	ctx *PassContext,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
	varCands map[*types.Var]*setCandidate,
	funcCands map[*types.Func]*setCandidate,
	structFields map[*types.Named][]*setCandidate,
	edges *[][2]*setCandidate,
) {
	info := ctx.Info
	disqualify := func(named *types.Named) {
		for _, c := range structFields[named] {
			c.disqualified = true
		}
	}

	isNonNilSetInitRhs := func(rhs ast.Expr) (*setCandidate, bool) {
		if _, _, otherFn, _, _, ok := classifySetInitExpr(rhs, funcCands, info); ok {
			return otherFn, true
		}
		if _, hasNil, _, _, _, ok := classifySetIIFEInit(rhs, funcCands, info); ok && !hasNil {
			return nil, true
		}
		return nil, false
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
							if cand := varCands[varFromObj(info.Uses[sel.Sel])]; cand != nil && cand.isField {
								if _, initOK := isNonNilSetInitRhs(as.Rhs[0]); initOK {
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

	// order-independent: sets per-candidate boolean hasInit flag
	for named, cands := range structFields {
		ctors := ctorFuncs[named]
		if len(ctors) == 0 {
			continue
		}
		for _, cand := range cands {
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
								if id, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == recvVar && varFromObj(info.Uses[sel.Sel]) == cand.v {
									if otherFn, ok := isNonNilSetInitRhs(s.Rhs[0]); ok {
										found = true
										if otherFn != nil {
											*edges = append(*edges, [2]*setCandidate{cand, otherFn})
										}
									}
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
										// Constructor chaining (`InitFoo_2` -> `InitFoo_4(this_, ...)`): defer to target ctor.
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

	watched := make(map[*types.Named]bool, len(structFields))
	// order-independent: populates membership set of candidate struct types
	for named := range structFields {
		watched[named] = true
	}
	checkStructValueCopies(ctx, parentsByFile, watched, func(fv *types.Var) {
		if c := varCands[fv]; c != nil {
			c.disqualified = true
		}
	}, disqualify)
}

func validateSetFuncBodyReturns(
	fc *setCandidate,
	varCands map[*types.Var]*setCandidate,
	funcCands map[*types.Func]*setCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	edges *[][2]*setCandidate,
) {
	hasReturn := false
	var walk func(ast.Node) bool
	walk = func(n ast.Node) bool {
		if fc.disqualified || n == nil {
			return false
		}
		if _, isLit := n.(*ast.FuncLit); isLit {
			return false
		}
		ret, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		hasReturn = true
		if len(ret.Results) != 1 {
			fc.disqualified = true
			return false
		}
		origRes := ret.Results[0]
		inner := unwrapCollectionAdaptersDown(origRes, info)
		srcCand := resolveSetSourceExpr(inner, varCands, funcCands, info)
		if srcCand == nil || (srcCand.v != nil && (srcCand.isField || enclosingFuncObj(inner, parents, info) != fc.fn)) {
			fc.disqualified = true
			return false
		}
		*edges = append(*edges, [2]*setCandidate{fc, srcCand})
		if inner != stripParens(origRes) {
			fc.ops = append(fc.ops, setRewriteOp{
				kind: setOpReturnUnwrap, cand: fc, file: fc.declFile, spanNode: origRes, outerExpr: origRes, accessExpr: inner,
			})
		}
		return false
	}
	ast.Inspect(fc.fnDecl.Body, walk)
	if !hasReturn {
		fc.disqualified = true
	}
}

func validateSetParamFuncCall(
	id *ast.Ident,
	pcs []*setCandidate,
	varCands map[*types.Var]*setCandidate,
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
	// Unwrap generic instantiation `[N, E]` if present on callee.
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
		arg := unwrapCollectionAdaptersDown(call.Args[pc.paramIdx], info)
		srcCand := resolveSetVarAccess(arg, varCands, info)
		if srcCand == nil || srcCand.disqualified {
			pc.disqualified = true
			continue
		}
		pc.inPkgCalls++
		if srcCand != pc {
			pc.nonSelfCalls++
		}
		if pc.isExportedFn {
			pc.paramCallIdents = append(pc.paramCallIdents, id)
		}
	}
}

func validateSetVarUse(
	id *ast.Ident,
	file *ast.File,
	cand *setCandidate,
	varCands map[*types.Var]*setCandidate,
	funcCands map[*types.Func]*setCandidate,
	paramCandsByFn map[*types.Func][]*setCandidate,
	ctx *PassContext,
	parents map[ast.Node]ast.Node,
	varUseCounts map[*types.Var]int,
	markMut func(*setCandidate, ast.Node, map[ast.Node]ast.Node),
	edges *[][2]*setCandidate,
) {
	info := ctx.Info
	var accessExpr ast.Expr = id
	if cand.isField {
		sel, ok := parents[id].(*ast.SelectorExpr)
		if !ok || sel.Sel != id || !isIndirectStructFieldSel(sel, info) {
			cand.disqualified = true
			return
		}
		accessExpr = sel
	}

	if as, ok := parents[accessExpr].(*ast.AssignStmt); ok {
		for _, lhs := range as.Lhs {
			if stripParens(lhs) == accessExpr {
				if cand.isParam || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 || (cand.isField && !isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info)) {
					cand.disqualified = true
					return
				}
				rhs := as.Rhs[0]
				if isNilIdent(rhs, info) {
					cand.hasNilAssign = true
					cand.ops = append(cand.ops, setRewriteOp{
						kind: setOpAssignInit, cand: cand, initKind: setInitNil, file: file, spanNode: as, stmt: as, accessExpr: accessExpr,
					})
					if !cand.isField || !isInInitFieldsMethod(as, cand.ownerStruct, parents, info) {
						markMut(cand, as, parents)
					}
					return
				}
				if initK, argExpr, otherFn, idK, hsK, initOK := classifySetInitExpr(rhs, funcCands, info); initOK {
					if idK {
						cand.hasIdentityKind = true
					}
					if hsK {
						cand.hasHashKind = true
					}
					var srcCand *setCandidate
					if initK == setInitCopyCol {
						unwrapped := unwrapCollectionAdaptersDown(argExpr, info)
						if sc := resolveSetSourceExpr(unwrapped, varCands, funcCands, info); sc != nil && sc != cand {
							initK = setInitCopySet
							argExpr = unwrapped
							srcCand = sc
							cand.usesOrdSetOnly = true
							sc.usesOrdSetOnly = true
							*edges = append(*edges, [2]*setCandidate{cand, sc})
						}
					}
					if otherFn != nil {
						srcCand = otherFn
						*edges = append(*edges, [2]*setCandidate{cand, otherFn})
					}
					cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct, recvTArgs: cand.recvTArgs})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
					cand.ops = append(cand.ops, setRewriteOp{
						kind: setOpAssignInit, cand: cand, otherCand: srcCand, initKind: initK, file: file, spanNode: as, stmt: as, accessExpr: accessExpr, argExpr: argExpr, call: asCallExpr(rhs),
					})
					if !cand.isField || !isInInitFieldsMethod(as, cand.ownerStruct, parents, info) {
						markMut(cand, as, parents)
					}
					return
				}
				if iife, specHasNil, _, idK, hsK, ok := classifySetIIFEInit(rhs, funcCands, info); ok {
					if specHasNil {
						cand.hasNilAssign = true
					}
					if idK {
						cand.hasIdentityKind = true
					}
					if hsK {
						cand.hasHashKind = true
					}
					for i := range iife.branches {
						br := &iife.branches[i]
						if br.initKind == setInitCopyCol {
							unwrapped := unwrapCollectionAdaptersDown(br.argExpr, info)
							if sc := resolveSetSourceExpr(unwrapped, varCands, funcCands, info); sc != nil && sc != cand {
								br.initKind = setInitCopySet
								br.argExpr = unwrapped
								br.otherCand = sc
								cand.usesOrdSetOnly = true
								sc.usesOrdSetOnly = true
								*edges = append(*edges, [2]*setCandidate{cand, sc})
							}
						}
						if br.otherCand != nil {
							*edges = append(*edges, [2]*setCandidate{cand, br.otherCand})
						}
					}
					cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct, recvTArgs: cand.recvTArgs})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
					cand.ops = append(cand.ops, setRewriteOp{
						kind: setOpAssignInit, cand: cand, initKind: setInitIIFE, iifeSpec: iife, file: file, spanNode: as, stmt: as, accessExpr: accessExpr,
					})
					if !cand.isField || !isInInitFieldsMethod(as, cand.ownerStruct, parents, info) {
						markMut(cand, as, parents)
					}
					return
				}
				cand.disqualified = true
				return
			}
		}
		if !cand.isField && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && stripParens(as.Rhs[0]) == id {
			if lhsId, ok := stripParens(as.Lhs[0]).(*ast.Ident); ok && lhsId.Name == "_" {
				return
			}
		}
	}

	validateSetValueConsumer(accessExpr, true, file, cand, varCands, funcCands, paramCandsByFn, parents, info, varUseCounts, markMut, edges)
}

func validateSetFuncUse(
	id *ast.Ident,
	file *ast.File,
	fc *setCandidate,
	varCands map[*types.Var]*setCandidate,
	funcCands map[*types.Func]*setCandidate,
	paramCandsByFn map[*types.Func][]*setCandidate,
	ctx *PassContext,
	parents map[ast.Node]ast.Node,
	varUseCounts map[*types.Var]int,
	edges *[][2]*setCandidate,
) {
	info := ctx.Info
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
	if call == nil {
		fc.disqualified = true
		return
	}
	fc.inPkgCalls++
	if fc.isExportedFn {
		fc.ops = append(fc.ops, setRewriteOp{kind: setOpFuncCallRename, cand: fc, file: file, spanNode: id, elemIdent: id})
	}

	// Check if call directly initializes a candidate var/field (`var s = Foo(...)` or `this_.f = Foo(...)` or inside IIFE).
	curP := unwrapParensUp(call, parents)
	switch par := parents[curP].(type) {
	case *ast.ValueSpec:
		if len(par.Names) == 1 && len(par.Values) == 1 && par.Values[0] == curP {
			if vc := varCands[varFromObj(info.Defs[par.Names[0]])]; vc != nil {
				return
			}
		}
	case *ast.AssignStmt:
		if len(par.Lhs) == 1 && len(par.Rhs) == 1 && par.Rhs[0] == curP {
			if par.Tok == token.DEFINE {
				if lhsId, ok := par.Lhs[0].(*ast.Ident); ok && varCands[varFromObj(info.Defs[lhsId])] != nil {
					return
				}
			} else if par.Tok == token.ASSIGN {
				if idL, ok := stripParens(par.Lhs[0]).(*ast.Ident); ok && varCands[varFromObj(info.Uses[idL])] != nil {
					return
				}
				if sel, ok := stripParens(par.Lhs[0]).(*ast.SelectorExpr); ok && varCands[varFromObj(info.Uses[sel.Sel])] != nil {
					return
				}
			}
		}
	}

	nopMut := func(*setCandidate, ast.Node, map[ast.Node]ast.Node) {}
	validateSetValueConsumer(call, false, file, fc, varCands, funcCands, paramCandsByFn, parents, info, varUseCounts, nopMut, edges)
}

func validateSetValueConsumer(
	accessExpr ast.Expr,
	addressable bool,
	file *ast.File,
	cand *setCandidate,
	varCands map[*types.Var]*setCandidate,
	funcCands map[*types.Func]*setCandidate,
	paramCandsByFn map[*types.Func][]*setCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
	markMut func(*setCandidate, ast.Node, map[ast.Node]ast.Node),
	edges *[][2]*setCandidate,
) {
	// Allow `s == nil` / `s != nil` on var/field candidates (marks candidate as nullable pointer `*OrderedSet[T]`).
	pTop := unwrapParensUp(accessExpr, parents)
	if !cand.isFunc {
		if bin, ok := parents[pTop].(*ast.BinaryExpr); ok && (bin.Op == token.EQL || bin.Op == token.NEQ) {
			other := bin.Y
			if bin.Y == pTop {
				other = bin.X
			}
			if isNilIdent(other, info) {
				cand.hasNilCheck = true
				return
			}
		}
		if outerSink, isPtrWrapped, ok := matchFormattingSetSink(pTop, parents, info); ok {
			if isPtrWrapped {
				// `j2gort.PtrToAny(s)` inside formatting sink: if `s` becomes `*OrderedSet[T]`, `PtrToAny(s)` compiles natively!
				// If `s` is non-nullable, `setOpFormatSink` replaces `PtrToAny(s)` with `&s`.
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: accessExpr.Pos()})
				cand.ops = append(cand.ops, setRewriteOp{
					kind: setOpFormatSink, cand: cand, file: file, spanNode: outerSink, outerExpr: outerSink, accessExpr: accessExpr, unwrapsBox: true,
				})
				return
			}
			if addressable {
				cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: accessExpr.Pos()})
				cand.ops = append(cand.ops, setRewriteOp{
					kind: setOpFormatSink, cand: cand, file: file, spanNode: accessExpr, outerExpr: accessExpr, accessExpr: accessExpr,
				})
				return
			}
		}
	}

	adaptedTop := unwrapCollectionAdapters(accessExpr, parents, info)

	// Case A: Returned from enclosing candidate Set-returning function (`return unsafe` / `return AsIface(unsafe)`).
	if ret, ok := parents[adaptedTop].(*ast.ReturnStmt); ok && len(ret.Results) == 1 && ret.Results[0] == adaptedTop {
		if enclFn := enclosingFuncObj(ret, parents, info); enclFn != nil && funcCands[funcOrigin(enclFn)] != nil {
			if !cand.isField {
				return
			}
		}
	}

	// Case B: Consumed by a j2gort static helper or collection constructor/adapter.
	if staticCall, ok := parents[adaptedTop].(*ast.CallExpr); ok && !staticCall.Ellipsis.IsValid() && staticCall.Fun != adaptedTop {
		if fnName, _, _, isJ2 := isJ2gortCall(staticCall, info); isJ2 {
			switch {
			case len(staticCall.Args) == 2 && staticCall.Args[0] == adaptedTop:
				if validateStaticSetOp(fnName, staticCall, accessExpr, addressable, file, cand, parents, info, markMut) {
					return
				}

			case len(staticCall.Args) == 1 && staticCall.Args[0] == adaptedTop:
				switch fnName {
				case "NewIterCursor":
					if !addressable {
						cand.disqualified = true
						return
					}
					validateSetIterCursorLoop(staticCall, accessExpr, file, cand, parents, info, varUseCounts)
					return

				case "Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection",
					"Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_lang_Iterable":
					if _, isES := parents[staticCall].(*ast.ExprStmt); !isES {
						cand.ops = append(cand.ops, setRewriteOp{
							kind: setOpImmSetCopy, cand: cand, file: file, spanNode: staticCall, outerExpr: staticCall, accessExpr: accessExpr,
						})
						return
					}

				case "NewExt_java_util_ArrayDeque_1_Ext_java_util_Collection":
					if _, isES := parents[staticCall].(*ast.ExprStmt); !isES {
						cand.usesOrdSetOnly = true
						cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: staticCall.Pos()})
						cand.ops = append(cand.ops, setRewriteOp{
							kind: setOpArrayDequeCopy, cand: cand, file: file, spanNode: staticCall, outerExpr: staticCall, accessExpr: accessExpr,
						})
						return
					}

				case "DequeFromCol", "DequePtrFromCol", "NodeDequeFromCol", "NodeDequePtrFromCol":
					if _, isES := parents[staticCall].(*ast.ExprStmt); !isES {
						isNodeDeq := strings.HasPrefix(fnName, "NodeDeque")
						reqExport := "DequeFromOrderedSet"
						if isNodeDeq {
							reqExport = "NodeDequeFromOrderedSet"
						}
						var j2Pkg *types.Package
						if callee, _ := calleeFunc(staticCall, info); callee != nil {
							j2Pkg = callee.Pkg()
						}
						if j2Pkg != nil && j2Pkg.Scope() != nil && j2Pkg.Scope().Lookup(reqExport) != nil {
							deqNamed, deqTArgs := derefNamedAndTypeArgs(info.Types[staticCall].Type)
							if deqNamed != nil && deqTArgs != nil && deqTArgs.Len() == 1 {
								deqElemT := deqTArgs.At(0)
								if !cand.addElemType(deqElemT, accessExpr, info) {
									cand.disqualified = true
									return
								}
								cand.usesOrdSetOnly = true
								cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: staticCall.Pos()})
								cand.ops = append(cand.ops, setRewriteOp{
									kind: setOpNativeDequeCopy, cand: cand, file: file, spanNode: staticCall,
									outerExpr: staticCall, accessExpr: accessExpr,
									methodsSel: strings.TrimSuffix(fnName, "Col"),
								})
								return
							}
						}
					}

				case "NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection",
					"NewExt_java_util_HashSet_1_Ext_java_util_Collection",
					"Ext_com_google_common_collect_Sets_NewLinkedHashSet_1_Ext_java_lang_Iterable",
					"Ext_com_google_common_collect_Sets_NewHashSet_1_Ext_java_lang_Iterable":
					if isSetInitTargetCandidate(staticCall, varCands, parents, info) {
						return
					}
				}
			}
		}

		// Check if `cand` is passed as a Set parameter to an unexported function / priv_ method.
		if callee, _ := calleeFunc(staticCall, info); callee != nil {
			if pcs := paramCandsByFn[funcOrigin(callee)]; len(pcs) > 0 {
				for _, pc := range pcs {
					if pc.paramIdx < len(staticCall.Args) && staticCall.Args[pc.paramIdx] == adaptedTop {
						if !addressable {
							cand.disqualified = true
							pc.disqualified = true
							return
						}
						if pc.ownerStruct != nil {
							callSel, selOK := stripParens(staticCall.Fun).(*ast.SelectorExpr)
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
						cand.usesOrdSetOnly = true
						pc.usesOrdSetOnly = true
						*edges = append(*edges, [2]*setCandidate{cand, pc})
						cand.ops = append(cand.ops, setRewriteOp{
							kind: setOpParamArg, cand: cand, otherCand: pc, file: file, spanNode: adaptedTop,
							outerExpr: adaptedTop, accessExpr: accessExpr,
						})
						markMut(cand, staticCall, parents)
						return
					}
				}
			}
		}

		// Check if `cand` is the source argument of `dstSet.AddAll_1(AsIface(cand))` or `otherSet.Equals(...)`.
		if len(staticCall.Args) == 1 && staticCall.Args[0] == adaptedTop {
			if mSel, ok := stripParens(staticCall.Fun).(*ast.SelectorExpr); ok && (mSel.Sel.Name == "AddAll_1" || mSel.Sel.Name == "Equals") {
				if dstCand := resolveSetVarAccess(unwrapCheckNotNullArg(mSel.X, info), varCands, info); dstCand != nil && (dstCand != cand || (mSel.Sel.Name == "Equals" && cand.isField)) {
					return
				}
			}
		}
	}
	if ptrWrap, ok := parents[pTop].(*ast.CallExpr); ok && len(ptrWrap.Args) == 1 && !ptrWrap.Ellipsis.IsValid() && ptrWrap.Args[0] == pTop {
		if fnName, _, _, isJ2 := isJ2gortCall(ptrWrap, info); isJ2 && (fnName == "PtrToAny" || fnName == "ToAny") {
			wTop := unwrapParensUp(ptrWrap, parents)
			if eqCall, ok := parents[wTop].(*ast.CallExpr); ok && len(eqCall.Args) == 1 && !eqCall.Ellipsis.IsValid() && eqCall.Args[0] == wTop {
				if mSel, ok := stripParens(eqCall.Fun).(*ast.SelectorExpr); ok && mSel.Sel.Name == "Equals" {
					if dstCand := resolveSetVarAccess(unwrapCheckNotNullArg(mSel.X, info), varCands, info); dstCand != nil && (dstCand != cand || cand.isField) {
						return
					}
				}
			}
		}
	}

	// Case C: Direct method call `s.<Method>(...)`.
	topExpr := unwrapEnclosingCheckNotNulls(accessExpr, parents, info)
	methodSel, selOK := parents[topExpr].(*ast.SelectorExpr)
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

	switch methodSel.Sel.Name {
	case "Add_1", "AddLast", "Contains", "Remove_1_any", "Remove":
		isContains := methodSel.Sel.Name == "Contains"
		if (!isContains && !addressable) || len(mCall.Args) != 1 || (methodSel.Sel.Name == "AddLast" && !inExprStmt) || (isContains && inExprStmt) {
			cand.disqualified = true
			return
		}
		sel := "Add"
		if isContains {
			sel = "Contains"
		} else if methodSel.Sel.Name != "Add_1" && methodSel.Sel.Name != "AddLast" {
			sel = "Remove"
		}
		if unboxed, elemT, ok := unwrapSetOrMapElemArg(mCall.Args[0], info); ok && cand.addElemType(elemT, accessExpr, info) {
			cand.ops = append(cand.ops, setRewriteOp{
				kind: setOpMethodUnbox1, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
				accessExpr: accessExpr, call: mCall, methodsSel: sel, unwrapsBox: unboxed != stripParens(mCall.Args[0]),
			})
			if !isContains {
				markMut(cand, mCall, parents)
			}
		} else {
			cand.disqualified = true
		}

	case "AddAll_1":
		if !addressable || len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		unwrappedArg := unwrapCollectionAdaptersDown(mCall.Args[0], info)
		if srcCand := resolveSetSourceExpr(unwrappedArg, varCands, funcCands, info); srcCand != nil && srcCand != cand {
			cand.usesOrdSetOnly = true
			srcCand.usesOrdSetOnly = true
			*edges = append(*edges, [2]*setCandidate{cand, srcCand})
			cand.ops = append(cand.ops, setRewriteOp{
				kind: setOpAddAllSet, cand: cand, otherCand: srcCand, file: file, spanNode: mCall, outerExpr: mCall,
				accessExpr: accessExpr, argExpr: unwrappedArg,
			})
			markMut(cand, mCall, parents)
			return
		}
		cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: mCall.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, setRewriteOp{
			kind: setOpAddAllCol, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, argExpr: mCall.Args[0],
		})
		markMut(cand, mCall, parents)

	case "Size", "IsEmpty", "HashCode":
		if inExprStmt || len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		cand.ops = append(cand.ops, setRewriteOp{
			kind: setOpMethodDirect, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, methodsSel: methodSel.Sel.Name,
		})

	case "Equals":
		if inExprStmt || len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		otherInner := unwrapEqualsOtherArg(mCall.Args[0], info)
		otherCand := resolveSetVarAccess(otherInner, varCands, info)
		if otherCand == nil || (otherCand == cand && !cand.isField) {
			cand.disqualified = true
			return
		}
		cand.usesOrdSetOnly = true
		otherCand.usesOrdSetOnly = true
		*edges = append(*edges, [2]*setCandidate{cand, otherCand})
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, setRewriteOp{
			kind: setOpEqualsSet, cand: cand, otherCand: otherCand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, argExpr: otherInner,
		})

	case "Clear":
		if !addressable || !inExprStmt || len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		cand.ops = append(cand.ops, setRewriteOp{
			kind: setOpMethodDirect, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, methodsSel: "Clear",
		})
		markMut(cand, mCall, parents)

	case "GetFirst":
		outerExpr, expectedT, ok := enclosingUnboxExpr(mCall, parents, info)
		_, outerES := parents[outerExpr].(*ast.ExprStmt)
		if !ok || outerES || len(mCall.Args) != 0 || !cand.addElemType(expectedT, accessExpr, info) {
			cand.disqualified = true
			return
		}
		cand.ops = append(cand.ops, setRewriteOp{
			kind: setOpGetFirst, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr, accessExpr: accessExpr,
		})

	default:
		cand.disqualified = true
	}
}

func validateSetIterCursorLoop(
	cursorCall *ast.CallExpr,
	accessExpr ast.Expr,
	file *ast.File,
	cand *setCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
) {
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
	if lbl, isLbl := nextStmt.(*ast.LabeledStmt); isLbl {
		nextStmt = lbl.Stmt
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
		cand.addElemType(types.Typ[types.String], accessExpr, info)
	case isCursorMethodCall(rhs0, curVar, "NextNode", info):
		cand.addElemType(types.Typ[types.Uint32], accessExpr, info)
	case isNodeConvOfCursorNextNode(rhs0, curVar, info):
		cand.addElemType(info.Types[rhs0].Type, accessExpr, info)
	case isCursorMethodCall(rhs0, curVar, "Next", info):
		rawVar := varFromObj(info.Defs[lhs0])
		if rawVar == nil {
			cand.disqualified = true
			return
		}
		if len(forStmt.Body.List) >= 2 && varUseCounts[rawVar] == 1 {
			if unboxedId, elemT, ok := matchSingleDefineUnbox(forStmt.Body.List[1], rawVar, parents, info); ok && cand.addElemType(elemT, accessExpr, info) {
				elemId = unboxedId
				drop = 2
				break
			}
		}
		if varUseCounts[rawVar] == 0 {
			break
		}
		valid := true
		for _, st := range forStmt.Body.List[1:] {
			ast.Inspect(st, func(n ast.Node) bool {
				id, ok := n.(*ast.Ident)
				if !valid || !ok || varFromObj(info.Uses[id]) != rawVar {
					return valid
				}
				if as, isAs := parents[id].(*ast.AssignStmt); isAs && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && as.Rhs[0] == id {
					if lhsBlank, ok := stripParens(as.Lhs[0]).(*ast.Ident); ok && lhsBlank.Name == "_" {
						return true
					}
				}
				outerExpr, expectedT, ok := enclosingUnboxExpr(id, parents, info)
				_, isES := parents[outerExpr].(*ast.ExprStmt)
				if !ok || isES || !cand.addElemType(expectedT, accessExpr, info) {
					valid = false
					return false
				}
				if outerExpr != id {
					cand.ops = append(cand.ops, setRewriteOp{kind: setOpRawUnbox, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr, elemIdent: id})
				}
				return true
			})
		}
		if !valid {
			cand.disqualified = true
			return
		}
	default:
		cand.disqualified = true
		return
	}

	cand.loopBodies = append(cand.loopBodies, forStmt.Body)
	cand.ops = append(cand.ops, setRewriteOp{
		kind: setOpIterCursor, cand: cand, file: file, spanNode: forStmt,
		stmt: initAs, forStmt: forStmt, outerExpr: cursorCall, accessExpr: accessExpr, elemIdent: elemId, dropStmts: drop,
	})
}

func anySetLoopMutatesCandidate(cand *setCandidate, mutatingFuncs map[*types.Func]bool, info *types.Info) bool {
	for _, body := range cand.loopBodies {
		mutated := false
		ast.Inspect(body, func(n ast.Node) bool {
			if mutated {
				return false
			}
			switch x := n.(type) {
			case *ast.CallExpr:
				if len(mutatingFuncs) > 0 {
					if fn, _ := calleeFunc(x, info); fn != nil && mutatingFuncs[funcOrigin(fn)] {
						mutated = true
						return false
					}
				}
				if fnName, _, _, isJ2 := isJ2gortCall(x, info); isJ2 {
					switch fnName {
					case "SetAddStr", "SetRemoveStr", "SetAddNode", "SetRemoveNode":
						if len(x.Args) >= 1 && exprReferencesVar(x.Args[0], cand.v, info) {
							mutated = true
						}
					}
				}
				if sel, ok := stripParens(x.Fun).(*ast.SelectorExpr); ok {
					switch sel.Sel.Name {
					case "Add_1", "AddLast", "AddAll_1", "Remove", "Remove_1_any", "Clear":
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

func applySetRewrites(ctx *PassContext, eligible []*setCandidate, parentsByFile map[*ast.File]map[ast.Node]ast.Node) int {
	// Snapshot untouched FuncDecls of exported Set-parameter functions BEFORE expr/stmt rewrites mutate their bodies,
	// so external callers retain an exact copy of the original exported function signature and body.
	type expParamFnTwin struct {
		fnDecl    *ast.FuncDecl
		file      *ast.File
		origName  string
		origClone *ast.FuncDecl
		callIds   []*ast.Ident
	}
	var expParamTwins []*expParamFnTwin
	seenExpParamFn := make(map[*ast.FuncDecl]*expParamFnTwin)
	for _, cand := range eligible {
		if !cand.isParam || !cand.isExportedFn || cand.fnDecl == nil || cand.declFile == nil {
			continue
		}
		tw := seenExpParamFn[cand.fnDecl]
		if tw == nil {
			cl, err := cloneFuncDeclFresh(ctx.Fset, cand.fnDecl)
			if err != nil {
				continue
			}
			tw = &expParamFnTwin{
				fnDecl:    cand.fnDecl,
				file:      cand.declFile,
				origName:  cand.fnDecl.Name.Name,
				origClone: cl,
			}
			seenExpParamFn[cand.fnDecl] = tw
			expParamTwins = append(expParamTwins, tw)
		}
		tw.callIds = append(tw.callIds, cand.paramCallIdents...)
	}

	var exprOps, stmtOps, funcOps []setRewriteOp
	for _, cand := range eligible {
		if cand.isParam && cand.isExportedFn && cand.fnDecl != nil && seenExpParamFn[cand.fnDecl] == nil {
			continue
		}
		for _, op := range cand.ops {
			switch op.kind {
			case setOpFuncDecl:
				funcOps = append(funcOps, op)
			case setOpDeclField, setOpDeclParam, setOpDeclVarSpec, setOpDeclShortVar, setOpAssignInit, setOpIterCursor:
				stmtOps = append(stmtOps, op)
			default:
				exprOps = append(exprOps, op)
			}
		}
	}
	spanOrder := func(a, b setRewriteOp) int {
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
		if applySingleSetRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}
	for _, op := range stmtOps {
		if applySingleSetRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}
	for _, op := range funcOps {
		if applySingleSetRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}
	for _, tw := range expParamTwins {
		twinName := "ordSet_" + tw.origName
		tw.fnDecl.Name.Name = twinName
		seenId := make(map[*ast.Ident]bool, len(tw.callIds))
		for _, cid := range tw.callIds {
			if seenId[cid] {
				continue
			}
			seenId[cid] = true
			cid.Name = twinName
			count++
		}
		tw.file.Decls = append(tw.file.Decls, tw.origClone)
		ctx.MarkFileModified(tw.file)
		count++
	}
	return count
}

func applySingleSetRewrite(ctx *PassContext, op setRewriteOp, parents map[ast.Node]ast.Node) bool {
	info := ctx.Info
	j2Sel := func(name string, pos token.Pos) *ast.SelectorExpr {
		q, _ := j2gortQualifierForFile(op.file, ctx)
		return &ast.SelectorExpr{X: &ast.Ident{NamePos: pos, Name: q}, Sel: &ast.Ident{NamePos: pos, Name: name}}
	}
	siteElemT := op.cand.elemType
	if op.cand.isField {
		siteElemT = instantiateOwnerTypeAtExpr(siteElemT, op.cand.ownerStruct, op.accessExpr, info)
	} else if op.cand.ownerStruct != nil && op.cand.recvTArgs != nil {
		if inst, ok := instantiateOwnerType(siteElemT, op.cand.ownerStruct, op.cand.recvTArgs); ok {
			siteElemT = inst
		}
	}
	makeElemType := func() ast.Expr {
		elt, _ := typeExprForFile(siteElemT, op.file, ctx.Pkg, ctx.Importer)
		return elt
	}
	setBaseName := "OrderedSet"
	nonLocalBaseName := "OrderedSet"
	switch op.cand.flavor {
	case setFlavorCustom:
		setBaseName = "CustomOrderedSet"
		nonLocalBaseName = "CustomOrderedSet"
	case setFlavorIdentity:
		setBaseName = "IdentityOrderedSet"
		nonLocalBaseName = "IdentityOrderedSet"
	default:
		if op.cand.useLocalSet {
			setBaseName = "LocalOrderedSet"
		}
	}
	makeOrdSetType := func(pos token.Pos) ast.Expr {
		base := ast.Expr(&ast.IndexExpr{X: j2Sel(setBaseName, pos), Index: makeElemType()})
		if op.cand.isPtr {
			return &ast.StarExpr{Star: pos, X: base}
		}
		return base
	}
	methodTarget := func(pos token.Pos) ast.Expr {
		if op.needsNullWrap {
			return &ast.CallExpr{Fun: j2Sel("CheckNotNullPtr", pos), Args: []ast.Expr{op.accessExpr}}
		}
		return wrapPrimaryTarget(op.accessExpr)
	}
	makeSingleInitRhs := func(initK setInitKind, argExpr ast.Expr, otherCand *setCandidate, rhs ast.Expr, inHeader bool) ast.Expr {
		pos := rhs.Pos()
		switch initK {
		case setInitNil:
			return &ast.Ident{NamePos: pos, Name: "nil"}
		case setInitEmpty:
			if op.cand.isPtr {
				var addr ast.Expr = &ast.UnaryExpr{
					OpPos: pos, Op: token.AND,
					X: &ast.CompositeLit{Type: &ast.IndexExpr{X: j2Sel(nonLocalBaseName, pos), Index: makeElemType()}},
				}
				if inHeader {
					addr = &ast.ParenExpr{X: addr}
				}
				return addr
			}
			return &ast.CallExpr{Fun: &ast.IndexExpr{X: j2Sel("New"+setBaseName, pos), Index: makeElemType()}}
		case setInitCap:
			if op.cand.isPtr {
				return &ast.CallExpr{Fun: &ast.IndexExpr{X: j2Sel("New"+nonLocalBaseName+"PtrCap", pos), Index: makeElemType()}, Args: []ast.Expr{argExpr}}
			}
			return &ast.CallExpr{Fun: &ast.IndexExpr{X: j2Sel("New"+setBaseName+"Cap", pos), Index: makeElemType()}, Args: []ast.Expr{argExpr}}
		case setInitCopySet:
			if op.cand.isPtr {
				srcPtr := argExpr
				if otherCand != nil && !otherCand.isPtr {
					srcPtr = &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: wrapPrimaryTarget(argExpr)}
				}
				return &ast.CallExpr{Fun: j2Sel(nonLocalBaseName+"PtrFromSet", pos), Args: []ast.Expr{srcPtr}}
			}
			srcVal := argExpr
			if otherCand != nil && otherCand.isPtr {
				srcVal = &ast.StarExpr{Star: pos, X: &ast.CallExpr{Fun: j2Sel("CheckNotNullPtr", pos), Args: []ast.Expr{argExpr}}}
			}
			return &ast.CallExpr{Fun: j2Sel(nonLocalBaseName+"FromSet", pos), Args: []ast.Expr{srcVal}}
		case setInitCopyCol:
			pfx := setBaseName + "From"
			if op.cand.isPtr {
				pfx = nonLocalBaseName + "PtrFrom"
			}
			switch {
			case op.cand.flavor == setFlavorOrdered && types.Identical(siteElemT, types.Typ[types.String]):
				return &ast.CallExpr{Fun: j2Sel(pfx+"StrCol", pos), Args: []ast.Expr{argExpr}}
			case op.cand.flavor == setFlavorOrdered && isNamedUint32NodeType(siteElemT):
				return &ast.CallExpr{Fun: &ast.IndexExpr{X: j2Sel(pfx+"NodeCol", pos), Index: makeElemType()}, Args: []ast.Expr{argExpr}}
			default:
				return &ast.CallExpr{Fun: &ast.IndexExpr{X: j2Sel(pfx+"Col", pos), Index: makeElemType()}, Args: []ast.Expr{argExpr}}
			}
		case setInitDirectCall:
			return rhs
		}
		return rhs
	}
	makeInitRhs := func(rhs ast.Expr, inHeader bool) ast.Expr {
		if op.initKind == setInitIIFE && op.iifeSpec != nil {
			fl := op.iifeSpec.fnLit
			fl.Type.Results.List[0].Type = makeOrdSetType(fl.Type.Results.Pos())
			for _, br := range op.iifeSpec.branches {
				br.ret.Results[0] = makeSingleInitRhs(br.initKind, br.argExpr, br.otherCand, br.ret.Results[0], false)
			}
			return rhs
		}
		return makeSingleInitRhs(op.initKind, op.argExpr, op.otherCand, rhs, inHeader)
	}

	switch op.kind {
	case setOpDeclField, setOpDeclParam:
		fld := op.spanNode.(*ast.Field)
		fld.Type = makeOrdSetType(fld.Pos())
		return true

	case setOpParamArg:
		parent := parents[op.outerExpr]
		var repl ast.Expr = op.accessExpr
		if !op.cand.isPtr {
			repl = &ast.UnaryExpr{OpPos: op.outerExpr.Pos(), Op: token.AND, X: wrapPrimaryTarget(op.accessExpr)}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpDeclVarSpec:
		spec := op.spanNode.(*ast.ValueSpec)
		spec.Type = makeOrdSetType(spec.Pos())
		if len(spec.Values) == 1 {
			if op.initKind == setInitNil {
				if !op.cand.isPtr {
					return false
				}
			} else {
				spec.Values[0] = makeInitRhs(spec.Values[0], false)
			}
		}
		return true

	case setOpDeclShortVar, setOpAssignInit:
		if op.initKind == setInitNil {
			return false
		}
		as := op.spanNode.(*ast.AssignStmt)
		as.Rhs[0] = makeInitRhs(as.Rhs[0], inControlHeader(as, parents))
		return true

	case setOpFuncCallRename:
		op.elemIdent.Name = "ordSet_" + op.cand.fn.Name()
		return true

	case setOpReturnUnwrap:
		parent := parents[op.outerExpr]
		return parent != nil && replaceExprInParent(parent, op.outerExpr, op.accessExpr, parents)

	case setOpFormatSink:
		if op.unwrapsBox && op.cand.isPtr {
			// `j2gort.PtrToAny(this_.reservedNames)` already compiles cleanly when `this_.reservedNames` is `*OrderedSet[T]`!
			return false
		}
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		var repl ast.Expr
		if op.cand.isPtr {
			repl = &ast.CallExpr{Fun: j2Sel("PtrToAny", pos), Args: []ast.Expr{op.accessExpr}}
		} else {
			repl = &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: wrapPrimaryTarget(op.accessExpr)}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpRawUnbox:
		parent := parents[op.outerExpr]
		return parent != nil && replaceExprInParent(parent, op.outerExpr, &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name}, parents)

	case setOpMethodDirect:
		parent := parents[op.outerExpr]
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: methodTarget(op.outerExpr.Pos()), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: op.methodsSel}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpMethodUnbox1:
		parent := parents[op.outerExpr]
		unboxed := unwrapBoxArgExpr(op.call.Args[0], op.unwrapsBox)
		repl := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: methodTarget(op.outerExpr.Pos()), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: op.methodsSel}},
			Args: []ast.Expr{unboxed},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpStaticHelper:
		parent := parents[op.outerExpr]
		repl := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: methodTarget(op.outerExpr.Pos()), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: op.methodsSel}},
			Args: []ast.Expr{op.argExpr},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpGetFirst:
		parent := parents[op.outerExpr]
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: methodTarget(op.outerExpr.Pos()), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "GetFirst"}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpAddAllSet:
		parent := parents[op.outerExpr]
		selName := "AddAllSet"
		if op.otherCand != nil && op.otherCand.isPtr {
			selName = "AddAllSetPtr"
		}
		repl := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: selName}},
			Args: []ast.Expr{op.argExpr},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpAddAllCol:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		var fn ast.Expr
		switch {
		case op.cand.flavor == setFlavorOrdered && types.Identical(siteElemT, types.Typ[types.String]):
			fn = j2Sel(setBaseName+"AddAllStr", pos)
		case op.cand.flavor != setFlavorCustom && isNamedUint32NodeType(siteElemT):
			fn = j2Sel(setBaseName+"AddAllNode", pos)
		default:
			fn = &ast.IndexExpr{X: j2Sel(setBaseName+"AddAllCol", pos), Index: makeElemType()}
		}
		addr := op.accessExpr
		if !op.cand.isPtr {
			addr = &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: wrapPrimaryTarget(op.accessExpr)}
		}
		repl := &ast.CallExpr{Fun: fn, Args: []ast.Expr{addr, op.argExpr}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpEqualsSet:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		if op.cand.isPtr || (op.otherCand != nil && op.otherCand.isPtr) {
			lhs := op.accessExpr
			if !op.cand.isPtr {
				lhs = &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: wrapPrimaryTarget(op.accessExpr)}
			}
			rhs := op.argExpr
			if op.otherCand != nil && !op.otherCand.isPtr {
				rhs = &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: wrapPrimaryTarget(op.argExpr)}
			}
			repl := &ast.CallExpr{
				Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(lhs), Sel: &ast.Ident{NamePos: pos, Name: "EqualsSetPtr"}},
				Args: []ast.Expr{rhs},
			}
			return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)
		}
		repl := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: methodTarget(pos), Sel: &ast.Ident{NamePos: pos, Name: "EqualsSet"}},
			Args: []ast.Expr{op.argExpr},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpImmSetCopy:
		parent := parents[op.outerExpr]
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: methodTarget(op.outerExpr.Pos()), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "ToImmutableSet"}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpArrayDequeCopy:
		parent := parents[op.outerExpr]
		helper := "ArrayDequeFrom" + nonLocalBaseName
		if op.cand.isPtr {
			helper = "ArrayDequeFrom" + nonLocalBaseName + "Ptr"
		}
		repl := &ast.CallExpr{
			Fun:  j2Sel(helper, op.outerExpr.Pos()),
			Args: []ast.Expr{op.accessExpr},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpNativeDequeCopy:
		parent := parents[op.outerExpr]
		helper := op.methodsSel + nonLocalBaseName
		if op.cand.isPtr {
			helper += "Ptr"
		}
		repl := &ast.CallExpr{
			Fun:  j2Sel(helper, op.outerExpr.Pos()),
			Args: []ast.Expr{op.accessExpr},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case setOpIterCursor:
		initAs := op.stmt.(*ast.AssignStmt)
		curId := initAs.Lhs[0].(*ast.Ident)
		initAs.Rhs[0] = &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "Iter"}},
		}
		if len(op.forStmt.Body.List) < op.dropStmts {
			return false
		}
		nextCall := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: &ast.Ident{NamePos: curId.Pos(), Name: curId.Name}, Sel: &ast.Ident{NamePos: op.elemIdent.Pos(), Name: "Next"}},
		}
		defStmt := &ast.AssignStmt{
			Lhs: []ast.Expr{&ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name}},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{nextCall},
		}
		op.forStmt.Body.List = append([]ast.Stmt{defStmt}, op.forStmt.Body.List[op.dropStmts:]...)
		return true

	case setOpFuncDecl:
		fn := op.cand.fnDecl
		ordSetType := makeOrdSetType(fn.Type.Results.Pos())
		if !op.cand.isExportedFn {
			fn.Type.Results.List[0].Type = ordSetType
			return true
		}
		origFuncName := fn.Name.Name
		twinName := "ordSet_" + origFuncName
		wrapperParams, callArgs := cloneFuncParamsAndCallArgs(fn.Type.Params)
		wrapperResults := &ast.FieldList{List: []*ast.Field{{Type: stripTypePos(fn.Type.Results.List[0].Type)}}}
		fn.Name.Name = twinName
		fn.Type.Results.List[0].Type = ordSetType

		twinCall := &ast.CallExpr{Fun: &ast.Ident{Name: twinName}, Args: callArgs}
		wrapperDecl := &ast.FuncDecl{
			Name: &ast.Ident{Name: origFuncName},
			Type: &ast.FuncType{
				Params:  &ast.FieldList{List: wrapperParams},
				Results: wrapperResults,
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CallExpr{Fun: &ast.SelectorExpr{X: twinCall, Sel: &ast.Ident{Name: "ToJavaSet"}}},
						},
					},
				},
			},
		}
		op.file.Decls = append(op.file.Decls, wrapperDecl)
		return true
	}
	return false
}

func classifySetInitExpr(expr ast.Expr, funcCands map[*types.Func]*setCandidate, info *types.Info) (kind setInitKind, argExpr ast.Expr, otherFn *setCandidate, isIdent, isHash, ok bool) {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return 0, nil, nil, false, false, false
	}
	if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
		switch fnName {
		case "NewExt_java_util_LinkedHashSet", "NewExt_java_util_HashSet",
			"Ext_com_google_common_collect_Sets_NewLinkedHashSet_0", "Ext_com_google_common_collect_Sets_NewHashSet_0":
			if len(call.Args) == 0 {
				return setInitEmpty, nil, nil, false, true, true
			}
		case "Ext_com_google_common_collect_Sets_NewIdentityHashSet":
			if len(call.Args) == 0 {
				return setInitEmpty, nil, nil, true, false, true
			}
		case "Ext_java_util_Collections_NewSetFromMap":
			if len(call.Args) == 1 {
				if innerCall, isCall := unwrapMapAdaptersDown(call.Args[0], info).(*ast.CallExpr); isCall && !innerCall.Ellipsis.IsValid() {
					if innerFn, _, _, innerJ2 := isJ2gortCall(innerCall, info); innerJ2 {
						switch innerFn {
						case "NewExt_java_util_IdentityHashMap", "Ext_com_google_common_collect_Maps_NewIdentityHashMap":
							if len(innerCall.Args) == 0 {
								return setInitEmpty, nil, nil, true, false, true
							}
						case "NewExt_java_util_IdentityHashMap_1_int32":
							if len(innerCall.Args) == 1 && isIntegerExpr(innerCall.Args[0], info) {
								return setInitCap, innerCall.Args[0], nil, true, false, true
							}
						}
					}
				}
			}
		case "NewExt_java_util_LinkedHashSet_1_int32", "NewExt_java_util_HashSet_1_int32",
			"Ext_com_google_common_collect_Sets_NewLinkedHashSetWithExpectedSize",
			"Ext_com_google_common_collect_Sets_NewHashSetWithExpectedSize":
			if len(call.Args) == 1 && isIntegerExpr(call.Args[0], info) {
				return setInitCap, call.Args[0], nil, false, true, true
			}
		case "NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection",
			"NewExt_java_util_HashSet_1_Ext_java_util_Collection",
			"Ext_com_google_common_collect_Sets_NewLinkedHashSet_1_Ext_java_lang_Iterable",
			"Ext_com_google_common_collect_Sets_NewHashSet_1_Ext_java_lang_Iterable":
			if len(call.Args) == 1 {
				return setInitCopyCol, call.Args[0], nil, false, true, true
			}
		}
	}
	if fn, _ := calleeFunc(call, info); fn != nil {
		if fc := funcCands[funcOrigin(fn)]; fc != nil {
			return setInitDirectCall, call, fc, false, false, true
		}
	}
	return 0, nil, nil, false, false, false
}

func classifySetIIFEInit(expr ast.Expr, funcCands map[*types.Func]*setCandidate, info *types.Info) (spec *setIIFESpec, hasNil bool, fns []*setCandidate, isIdent, isHash, ok bool) {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 0 || call.Ellipsis.IsValid() {
		return nil, false, nil, false, false, false
	}
	fl, ok := stripParens(call.Fun).(*ast.FuncLit)
	if !ok || fl.Body == nil || (fl.Type.Params != nil && len(fl.Type.Params.List) > 0) || fl.Type.Results == nil || len(fl.Type.Results.List) != 1 {
		return nil, false, nil, false, false, false
	}
	spec = &setIIFESpec{fnLit: fl}
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
		res := ret.Results[0]
		if isNilIdent(res, info) {
			hasNil = true
			spec.branches = append(spec.branches, setIIFEBranch{ret: ret, initKind: setInitNil})
			return false
		}
		initK, argExpr, otherFn, idK, hsK, initOK := classifySetInitExpr(res, funcCands, info)
		if !initOK {
			valid = false
			return false
		}
		if idK {
			isIdent = true
		}
		if hsK {
			isHash = true
		}
		if otherFn != nil {
			fns = append(fns, otherFn)
		}
		spec.branches = append(spec.branches, setIIFEBranch{ret: ret, initKind: initK, argExpr: argExpr, otherCand: otherFn})
		return false
	})
	if !valid || len(spec.branches) == 0 {
		return nil, false, nil, false, false, false
	}
	return spec, hasNil, fns, isIdent, isHash, true
}

func matchFormattingSetSink(pTop ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) (ast.Expr, bool, bool) {
	cur := pTop
	isPtrWrapped := false
	if wrapCall, ok := parents[cur].(*ast.CallExpr); ok && len(wrapCall.Args) == 1 && !wrapCall.Ellipsis.IsValid() && wrapCall.Args[0] == cur {
		if fnName, _, _, isJ2 := isJ2gortCall(wrapCall, info); isJ2 && (fnName == "PtrToAny" || fnName == "ToAny") {
			isPtrWrapped = true
			cur = unwrapParensUp(wrapCall, parents)
		}
	}
	outerSink := cur
	if isFormattingOnlySliceSink(cur, parents, info) {
		return outerSink, isPtrWrapped, true
	}
	if call, ok := parents[cur].(*ast.CallExpr); ok && call.Fun != cur {
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "CheckStateFmt", "CheckArgumentFmt":
				for i, a := range call.Args {
					if a == cur && i >= 2 {
						return outerSink, isPtrWrapped, true
					}
				}
			case "CheckNotNullFmt":
				for i, a := range call.Args {
					if a == cur && i >= 2 {
						return outerSink, isPtrWrapped, true
					}
				}
			}
		}
	}
	return nil, false, false
}

func isCandidateSetType(t types.Type) bool {
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
	case "Ext_java_util_LinkedHashSet", "Ext_java_util_HashSet",
		"Ext_java_util_Set", "Ext_java_util_SequencedSet", "Ext_java_util_Collection", "Ext_java_lang_Iterable":
		return true
	}
	return false
}

func isCandidateSetReturnType(t types.Type, exportedFn bool) bool {
	if !isCandidateSetType(t) {
		return false
	}
	if !exportedFn {
		return true
	}
	// For exported wrapper `return ordSet_Foo(...).ToJavaSet()`, `*Ext_java_util_LinkedHashSet` must be assignable to t.
	if ptr, ok := types.Unalias(t).(*types.Pointer); ok {
		named, ok := types.Unalias(ptr.Elem()).(*types.Named)
		return ok && named.Obj() != nil && named.Obj().Name() == "Ext_java_util_LinkedHashSet"
	}
	return isInterfaceType(t)
}

func isConcreteHashSetType(t types.Type) bool {
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
	case "Ext_java_util_LinkedHashSet", "Ext_java_util_HashSet":
		return true
	}
	return false
}

func classifySetElemType(t types.Type, hasIdentity bool, owners []*types.Named, ctx *PassContext) (setFlavor, bool) {
	mf, ok := classifyMapKeyType(t, hasIdentity, owners, ctx)
	if !ok {
		return 0, false
	}
	switch mf {
	case mapFlavorCustom:
		if !j2gortHasExport(ctx, "CustomOrderedSet") {
			return 0, false
		}
		return setFlavorCustom, true
	case mapFlavorIdentity:
		if !j2gortHasExport(ctx, "IdentityOrderedSet") {
			return setFlavorOrdered, true
		}
		return setFlavorIdentity, true
	default:
		return setFlavorOrdered, true
	}
}

func resolveSetVarAccess(expr ast.Expr, varCands map[*types.Var]*setCandidate, info *types.Info) *setCandidate {
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

func resolveSetSourceExpr(expr ast.Expr, varCands map[*types.Var]*setCandidate, funcCands map[*types.Func]*setCandidate, info *types.Info) *setCandidate {
	if vc := resolveSetVarAccess(expr, varCands, info); vc != nil {
		return vc
	}
	if call, ok := stripParens(expr).(*ast.CallExpr); ok && !call.Ellipsis.IsValid() {
		if fn, _ := calleeFunc(call, info); fn != nil {
			return funcCands[funcOrigin(fn)]
		}
	}
	return nil
}

func isSetInitTargetCandidate(initCall *ast.CallExpr, varCands map[*types.Var]*setCandidate, parents map[ast.Node]ast.Node, info *types.Info) bool {
	cur := unwrapParensUp(initCall, parents)
	// Also allow `return NewExt_java_util_LinkedHashSet(...)` inside a candidate's ternary IIFE initializer!
	if ret, isRet := parents[cur].(*ast.ReturnStmt); isRet && len(ret.Results) == 1 && ret.Results[0] == cur {
		for p := parents[ret]; p != nil; p = parents[p] {
			if fl, ok := p.(*ast.FuncLit); ok {
				flTop := unwrapParensUp(fl, parents)
				if iifeCall, ok := parents[flTop].(*ast.CallExpr); ok && iifeCall.Fun == flTop && len(iifeCall.Args) == 0 {
					cur = unwrapParensUp(iifeCall, parents)
				}
				break
			}
			if _, ok := p.(*ast.FuncDecl); ok {
				break
			}
		}
	}
	switch p := parents[cur].(type) {
	case *ast.ValueSpec:
		return len(p.Names) == 1 && len(p.Values) == 1 && p.Values[0] == cur && varCands[varFromObj(info.Defs[p.Names[0]])] != nil
	case *ast.AssignStmt:
		if len(p.Lhs) != 1 || len(p.Rhs) != 1 || p.Rhs[0] != cur {
			return false
		}
		if p.Tok == token.DEFINE {
			id, ok := p.Lhs[0].(*ast.Ident)
			return ok && varCands[varFromObj(info.Defs[id])] != nil
		}
		if p.Tok == token.ASSIGN {
			if id, ok := stripParens(p.Lhs[0]).(*ast.Ident); ok {
				return varCands[varFromObj(info.Uses[id])] != nil
			}
			sel, ok := stripParens(p.Lhs[0]).(*ast.SelectorExpr)
			return ok && varCands[varFromObj(info.Uses[sel.Sel])] != nil
		}
	}
	return false
}

func isIndirectStructFieldSel(sel *ast.SelectorExpr, info *types.Info) bool {
	seln, ok := info.Selections[sel]
	if !ok || seln.Kind() != types.FieldVal {
		return false
	}
	if seln.Indirect() {
		return true
	}
	if _, isPtr := types.Unalias(seln.Recv()).(*types.Pointer); isPtr {
		return true
	}
	if id, ok := stripParens(sel.X).(*ast.Ident); ok {
		_, isVar := info.Uses[id].(*types.Var)
		return isVar
	}
	return false
}

func allFuncParamsNamed(ft *ast.FuncType) bool {
	if ft == nil || ft.Params == nil {
		return true
	}
	for _, f := range ft.Params.List {
		if len(f.Names) == 0 {
			return false
		}
		for _, n := range f.Names {
			if n.Name == "_" || n.Name == "" {
				return false
			}
		}
	}
	return true
}

func cloneFuncParamsAndCallArgs(params *ast.FieldList) ([]*ast.Field, []ast.Expr) {
	if params == nil || len(params.List) == 0 {
		return nil, nil
	}
	fields := make([]*ast.Field, len(params.List))
	var args []ast.Expr
	for i, f := range params.List {
		names := make([]*ast.Ident, len(f.Names))
		for j, n := range f.Names {
			names[j] = &ast.Ident{Name: n.Name}
			args = append(args, &ast.Ident{Name: n.Name})
		}
		fields[i] = &ast.Field{Names: names, Type: stripTypePos(f.Type)}
	}
	return fields, args
}

func stripTypePos(e ast.Expr) ast.Expr {
	switch x := e.(type) {
	case *ast.Ident:
		return &ast.Ident{Name: x.Name}
	case *ast.SelectorExpr:
		return &ast.SelectorExpr{X: stripTypePos(x.X), Sel: &ast.Ident{Name: x.Sel.Name}}
	case *ast.StarExpr:
		return &ast.StarExpr{X: stripTypePos(x.X)}
	case *ast.ArrayType:
		var lenExpr ast.Expr
		if x.Len != nil {
			lenExpr = stripTypePos(x.Len)
		}
		return &ast.ArrayType{Len: lenExpr, Elt: stripTypePos(x.Elt)}
	case *ast.Ellipsis:
		return &ast.Ellipsis{Elt: stripTypePos(x.Elt)}
	case *ast.MapType:
		return &ast.MapType{Key: stripTypePos(x.Key), Value: stripTypePos(x.Value)}
	case *ast.IndexExpr:
		return &ast.IndexExpr{X: stripTypePos(x.X), Index: stripTypePos(x.Index)}
	case *ast.IndexListExpr:
		ixs := make([]ast.Expr, len(x.Indices))
		for i, ix := range x.Indices {
			ixs[i] = stripTypePos(ix)
		}
		return &ast.IndexListExpr{X: stripTypePos(x.X), Indices: ixs}
	}
	return e
}

func asCallExpr(expr ast.Expr) *ast.CallExpr {
	c, _ := stripParens(expr).(*ast.CallExpr)
	return c
}

var staticSetOpSpecs = map[string]struct {
	method string
	isNode bool
	mut    bool
}{
	"SetAddStr":       {"Add", false, true},
	"SetRemoveStr":    {"Remove", false, true},
	"SetContainsStr":  {"Contains", false, false},
	"SetAddNode":      {"Add", true, true},
	"SetRemoveNode":   {"Remove", true, true},
	"SetContainsNode": {"Contains", true, false},
}

func validateStaticSetOp(
	fnName string,
	staticCall *ast.CallExpr,
	accessExpr ast.Expr,
	addressable bool,
	file *ast.File,
	cand *setCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	markMut func(*setCandidate, ast.Node, map[ast.Node]ast.Node),
) bool {
	spec, ok := staticSetOpSpecs[fnName]
	if !ok {
		return false
	}
	_, isES := parents[staticCall].(*ast.ExprStmt)
	arg := staticCall.Args[1]
	var elemT types.Type
	if spec.isNode {
		if t := info.Types[arg].Type; isNamedUint32NodeType(t) {
			elemT = t
		}
	} else if hasExactOrDefaultStringType(arg, info) {
		elemT = types.Typ[types.String]
	}
	if (spec.mut && !addressable) || (!spec.mut && isES) || elemT == nil || !cand.addElemType(elemT, accessExpr, info) {
		cand.disqualified = true
		return true
	}
	cand.ops = append(cand.ops, setRewriteOp{
		kind: setOpStaticHelper, cand: cand, file: file, spanNode: staticCall, outerExpr: staticCall,
		accessExpr: accessExpr, argExpr: arg, methodsSel: spec.method,
	})
	if spec.mut {
		markMut(cand, staticCall, parents)
	}
	return true
}
