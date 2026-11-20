// Clean-room native_map_lowering pass: lowers non-escaping HashMap, LinkedHashMap, and
// IdentityHashMap locals and unexported struct fields across a package to zero-boxing
// value-type j2gort.OrderedMap[K, V], j2gort.CustomOrderedMap[K, V], or j2gort.IdentityOrderedMap[K, V].

package passes

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"slices"
	"strings"
)

// NativeMapLoweringPass rewrites non-escaping HashMap, LinkedHashMap, and IdentityHashMap
// local variables and unexported struct fields to value-type j2gort.OrderedMap[K, V]
// (for comparable identity/primitive/string/Node/comparable-TypeParam keys),
// j2gort.IdentityOrderedMap[K, V] (for IdentityHashMap with Go == identity keys), or
// j2gort.CustomOrderedMap[K, V] (for FastCustomKey types such as *Var, AbstractVar, JSType, and ScopedName).
type NativeMapLoweringPass struct{}

func NewNativeMapLoweringPass() *NativeMapLoweringPass { return &NativeMapLoweringPass{} }
func (p *NativeMapLoweringPass) Name() string          { return "native_map_lowering" }
func (p *NativeMapLoweringPass) Description() string {
	return "Lower non-escaping HashMap/LinkedHashMap/IdentityHashMap locals and unexported fields to j2gort.OrderedMap[K, V], CustomOrderedMap[K, V], or IdentityOrderedMap[K, V]"
}

func (p *NativeMapLoweringPass) Apply(ctx *PassContext) (int, error) {
	if isJ2gortPkg(ctx.Pkg) || len(ctx.Files) == 0 {
		return 0, nil
	}
	parentsByFile := make(map[*ast.File]map[ast.Node]ast.Node, len(ctx.Files))
	for _, f := range ctx.Files {
		parentsByFile[f], _ = buildParentMap(f)
	}

	count := 0
	eligible := analyzeMapCandidates(ctx, parentsByFile)
	if len(eligible) > 0 {
		count += applyMapRewrites(ctx, eligible, parentsByFile)
	}
	count += applyCanonicalVarIdentityOrderedMapUpgrade(ctx)
	tlCount, err := applyTimelineMapAndSliceLowering(ctx)
	if err != nil {
		return count, err
	}
	count += tlCount
	count += applyOrderedMapUnboxedIntLowering(ctx, parentsByFile)
	count += applyOrderedMapValueSliceLowering(ctx, parentsByFile)
	return count, nil
}

type mapFlavor uint8

const (
	mapFlavorOrdered mapFlavor = iota
	mapFlavorCustom
	mapFlavorIdentity
	mapFlavorPlain
)

type mapInitKind uint8

const (
	mapInitNil mapInitKind = iota
	mapInitEmpty
	mapInitCap
	mapInitCopyMap
	mapInitCopyJavaMap
	mapInitPtrAlias
	mapInitSharedPtr
	mapInitIIFE
	mapInitCastAccessorCall
)

type mapOpKind uint8

const (
	mapOpDeclField mapOpKind = iota
	mapOpDeclFieldGetter
	mapOpDeclParam
	mapOpDeclCastAccessor
	mapOpDeclVarSpec
	mapOpDeclShortVar
	mapOpAssignInit
	mapOpAssignUnionField // s.unionAnyField = localMapCand -> s.unionAnyField = j2gort.PtrToAny(localMapCand)
	mapOpPassParamArg     // fn(..., m, ...) -> fn(..., &m, ...) / fn(..., m, ...)
	mapOpNilCheck         // m == nil / m != nil (preserved natively on *OrderedMap / *CustomOrderedMap, or .IsEmpty() on lazy value OrderedMap)
	mapOpFormatSink       // StrVal(m), StrConcat(..., m), ToStringHelper.Add(..., m), CheckStateFmt(..., PtrToAny(m))
	mapOpMethodDirect     // .Size(), .IsEmpty(), .Clear(), .String()
	mapOpMethodRefGet     // bound := m; return NewLambda_Function(func(k any) any { return bound.Get(k) })
	mapOpFirstValue       // m.Values().Iterator().Next() -> m.FirstValue()
	mapOpGet              // m.Get(k), MapGetStr(m, k), MapGetNode(m, k)
	mapOpGetAny           // m.Get(rawAny) -> j2gort.*OrderedMapGetAny(&m, rawAny)
	mapOpGetOrDefault     // m.GetOrDefault(k, def), MapGetOrDefaultStr/Node(m, k, def)
	mapOpContainsKey      // m.ContainsKey(k), MapContainsKeyStr/Node(m, k)
	mapOpPut              // m.Put_2(k, v), MapPutStr/Node(m, k, v), m.PutIfAbsent(k, v)
	mapOpRemove           // m.Remove_1(k), MapRemoveStr/Node(m, k)
	mapOpPutAllMap        // m.PutAll(otherCand) -> m.PutAllMap(otherCand)
	mapOpPutAllJavaMap    // m.PutAll(extMap) -> m.PutAllJavaMap(extMap)
	mapOpComputeIfAbsent  // ComputeIfAbsent / MapComputeIfAbsent* / MapTryGetCompute* IIFE
	mapOpForEach          // m.ForEach_1(biConsumer) -> m.ForEach(func(k K, v V) { ... }) or m.ForEachBiConsumer(...)
	mapOpEqualsMap        // m.Equals(otherCand) -> OrderedMapEquals / CustomOrderedMapEquals / EqualsMap
	mapOpHashCode         // m.HashCode() -> m.HashCodeMap(isIdent) / m.HashCode()
	mapOpImmMapCopy       // ImmutableMap_CopyOf(m) -> m.ToImmutableMap()
	mapOpKeysToImmSet     // ImmutableSet_CopyOf(m.KeySet_0()) -> m.KeysToImmutableSet()
	mapOpValsToImmList    // ImmutableList_CopyOf(m.Values()) -> m.ValuesToImmutableList()
	mapOpLiveView         // this_.f.Values() -> this_.f.LiveValues(), this_.f.KeySet_0() -> this_.f.LiveKeySet()
	mapOpEntrySetReadOnly // m.EntrySet().Stream() -> m.ToJavaLinkedHashMap().EntrySet().Stream()
	mapOpIterCursor       // NewMapKeyCursor / NewMapKeyStrCursor / NewMapValCursor / NewIterCursor(KeySet/Values)
	mapOpEntryCursor      // NewMapEntryCursor(m) -> m.IterEntries()
	mapOpEntryGetKey      // cur.GetKey() / cur.KeyStr() / Node(cur.KeyNode()) -> cur.Key()
	mapOpEntryGetVal      // cur.GetValue() -> cur.Value()
	mapOpEntrySetVal      // cur.SetValue(v) -> cur.SetValue(unboxedV)
	mapOpRawUnbox         // unbox Ident replacement inside cursor/lambda body
)

type mapIIFEBranch struct {
	ret       *ast.ReturnStmt
	initKind  mapInitKind
	otherExpr ast.Expr
	otherCand *mapCandidate
}

type mapIIFESpec struct {
	fnLit    *ast.FuncLit
	branches []mapIIFEBranch
}

type mapValWriteSample struct {
	t            types.Type
	nonNilAlloc  bool
	isNilLiteral bool
}

type computeClosureSpec struct {
	fnExpr      ast.Expr
	boundInit   *ast.AssignStmt // optional `bound := this_` in method-ref IIFE
	fnLit       *ast.FuncLit    // non-nil when inline typed closure can be emitted
	retUnwraps  map[*ast.ReturnStmt]bool
	hasParamVar bool
	paramName   string
	paramPos    token.Pos
}

type forEachClosureSpec struct {
	argExpr   ast.Expr
	boundInit *ast.AssignStmt
	fnLit     *ast.FuncLit
	kName     string
	kPos      token.Pos
	vName     string
	vPos      token.Pos
}

type mapRewriteOp struct {
	kind          mapOpKind
	cand          *mapCandidate
	otherCand     *mapCandidate
	initKind      mapInitKind
	iifeSpec      *mapIIFESpec
	file          *ast.File
	spanNode      ast.Node
	accessExpr    ast.Expr
	otherExpr     ast.Expr
	keyExpr       ast.Expr
	valExpr       ast.Expr
	call          *ast.CallExpr
	outerExpr     ast.Expr
	stmt          ast.Stmt
	forStmt       *ast.ForStmt
	methodsSel    string
	elemIdent     *ast.Ident
	dropStmts     int
	keyArgIdx     int
	valArgIdx     int
	unwrapsKeyBox bool
	unwrapsValBox bool
	valIsNilLit   bool
	discardRet    bool
	isNegated     bool
	isPutIfAbsent bool
	viaPtrGetter  bool
	computeSpec   *computeClosureSpec
	forEachSpec   *forEachClosureSpec
	methodRefSpec *methodRefGetSpec
}

type methodRefGetSpec struct {
	boundAs  *ast.AssignStmt
	boundId  *ast.Ident
	getCall  *ast.CallExpr
	paramArg ast.Expr
}

type mapCandidate struct {
	v                  *types.Var
	isField            bool
	isParam            bool
	paramFn            *types.Func
	paramIdx           int
	isCastAccessor     bool
	castFn             *types.Func
	castFnDecl         *ast.FuncDecl
	unionFieldVar      *types.Var
	inPkgCastCalls     int
	hasUnionMapWrite   bool
	isPtrAlias         bool
	isSharedPtrLocal   bool
	isNullable         bool
	forceNullable      bool
	isLazyValueMap     bool
	hasNonCtorAssign   bool
	hasLocalLiveView   bool
	ownerStruct        *types.Named
	declFile           *ast.File
	defCount           int
	initAssignCount    int
	hasInit            bool
	hasNilAssign       bool
	hasNilCheck        bool
	hasHashKind        bool
	hasIdentityKind    bool
	isIdentityMap      bool
	flavor             mapFlavor
	disqualified       bool
	keyType            types.Type
	valType            types.Type
	keyTypes           []types.Type
	valReadTypes       []types.Type
	valWrites          []mapValWriteSample
	keySiteConstraints []siteTypeConstraint
	valSiteConstraints []siteTypeConstraint
	spellSites         []typeSpellSite
	j2gortSites        []typeSpellSite
	loopBodies         []*ast.BlockStmt
	ops                []mapRewriteOp
	parent             *mapCandidate
}

func (c *mapCandidate) isPtr() bool {
	return c.isPtrAlias || c.isNullable || c.forceNullable
}

func (c *mapCandidate) pos() token.Pos {
	if c.v != nil {
		return c.v.Pos()
	}
	if c.castFn != nil {
		return c.castFn.Pos()
	}
	return 0
}

func (c *mapCandidate) root() *mapCandidate {
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

func unionMapCandidates(a, b *mapCandidate) {
	ra, rb := a.root(), b.root()
	if ra != rb {
		rb.parent = ra
	}
}

func (c *mapCandidate) addKeyType(t types.Type, accessExpr ast.Expr, info *types.Info) bool {
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
			c.keySiteConstraints = append(c.keySiteConstraints, *sc)
			return true
		}
		t = lifted
	}
	c.keyTypes = append(c.keyTypes, t)
	return true
}

func (c *mapCandidate) addValReadType(t types.Type, accessExpr ast.Expr, info *types.Info) bool {
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
			c.valSiteConstraints = append(c.valSiteConstraints, *sc)
			return true
		}
		t = lifted
	}
	c.valReadTypes = append(c.valReadTypes, t)
	return true
}

func (c *mapCandidate) addValWrite(t types.Type, nonNilAlloc, isNilLit bool, accessExpr ast.Expr, info *types.Info) bool {
	if isNilLit {
		c.valWrites = append(c.valWrites, mapValWriteSample{isNilLiteral: true})
		return true
	}
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
			c.valSiteConstraints = append(c.valSiteConstraints, *sc)
			return true
		}
		t = lifted
	}
	c.valWrites = append(c.valWrites, mapValWriteSample{t: t, nonNilAlloc: nonNilAlloc})
	return true
}

func analyzeMapCandidates(ctx *PassContext, parentsByFile map[*ast.File]map[ast.Node]ast.Node) []*mapCandidate {
	info := ctx.Info
	varUseCounts := make(map[*types.Var]int, len(info.Uses))
	// order-independent: commutative per-variable integer use-count accumulation
	for _, obj := range info.Uses {
		if v := varFromObj(obj); v != nil {
			varUseCounts[v]++
		}
	}

	// Pre-detect mTmp variables bound inside j2go ComputeIfAbsent fast-path IIFEs so they are
	// handled atomically with their enclosing map candidate rather than as standalone locals.
	computeIIFETemps := collectComputeIIFEMapTemps(ctx.Files, parentsByFile, varUseCounts, info)

	varCands := make(map[*types.Var]*mapCandidate)
	castCandsByFn := make(map[*types.Func]*mapCandidate)
	castCandByUnionField := make(map[*types.Var]*mapCandidate)
	consumedCastCalls := make(map[*ast.CallExpr]bool)
	getterFieldCands := make(map[*types.Func]*mapCandidate)
	getterBodyByField := make(map[*types.Var]*ast.BlockStmt)
	paramCandsByFn := make(map[*types.Func][]*mapCandidate)
	var allCands []*mapCandidate
	var structOrder []*types.Named
	structFields := make(map[*types.Named][]*mapCandidate)
	var edges [][2]*mapCandidate
	var sharedPtrEdges [][2]*mapCandidate
	hasPtrHelpers := j2gortHasExport(ctx, "NewOrderedMapPtr")
	hasMethodRefHelpers := j2gortHasExport(ctx, "OrderedMapGetNodeAny")

	// Pass 1a: Collect unexported Map struct fields.
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
					if v == nil || !v.IsField() || v.Exported() || v.Embedded() {
						continue
					}
					isCand, declIdent, declHash := classifyCandidateMapType(v.Type())
					if !isCand {
						continue
					}
					cand := &mapCandidate{
						v: v, isField: true, ownerStruct: named, declFile: file,
						hasIdentityKind: declIdent, hasHashKind: declHash,
						spellSites:  []typeSpellSite{{file: file, pos: fld.Pos(), owner: named}},
						j2gortSites: []typeSpellSite{{file: file, pos: fld.Pos()}},
					}
					cand.ops = append(cand.ops, mapRewriteOp{kind: mapOpDeclField, cand: cand, file: file, spanNode: fld})
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

	// Pass 1a': Collect private 0-arg cast helpers on `any` union struct fields (`priv_...() Ext_java_util_Map`).
	if hasPtrHelpers && hasMethodRefHelpers {
		multiCastFields := make(map[*types.Var]bool)
		for _, file := range ctx.Files {
			if !ctx.ShouldTransform(file) {
				continue
			}
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv == nil || fn.Body == nil || !strings.HasPrefix(fn.Name.Name, "priv_") {
					continue
				}
				owner, unionField, fnOrig, declIdent, declHash, ok := matchUnionFieldCastAccessor(fn, info)
				if !ok {
					continue
				}
				if prev := castCandByUnionField[unionField]; prev != nil {
					prev.disqualified = true
					delete(castCandsByFn, prev.castFn)
					delete(castCandByUnionField, unionField)
					multiCastFields[unionField] = true
					continue
				}
				if multiCastFields[unionField] {
					continue
				}
				cc := &mapCandidate{
					isField: true, isCastAccessor: true, castFn: fnOrig, castFnDecl: fn, unionFieldVar: unionField,
					isNullable: true, forceNullable: true, hasInit: true, defCount: 1,
					hasIdentityKind: declIdent, hasHashKind: declHash,
					ownerStruct: owner, declFile: file,
					spellSites:  []typeSpellSite{{file: file, pos: fn.Pos(), owner: owner}},
					j2gortSites: []typeSpellSite{{file: file, pos: fn.Pos()}},
				}
				cc.ops = append(cc.ops, mapRewriteOp{kind: mapOpDeclCastAccessor, cand: cc, file: file, spanNode: fn})
				castCandsByFn[fnOrig] = cc
				castCandByUnionField[unionField] = cc
				allCands = append(allCands, cc)
				if len(structFields[owner]) == 0 {
					structOrder = append(structOrder, owner)
				}
				structFields[owner] = append(structFields[owner], cc)
			}
		}
	}

	// Pass 1a'': Collect 1-line direct struct-field map getters (e.g. `GetAssociatedProps()` in leaf `jscomp`
	// or private `priv_*` getters in any package) returning `this_.<field>`.
	if hasPtrHelpers {
		discoverMapFieldGetters(ctx, varCands, getterFieldCands, getterBodyByField)
		discoverMapFuncParams(ctx, varCands, paramCandsByFn, &allCands)
	}

	// Pass 1b: Collect local Map variables.
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		addLocal := func(id *ast.Ident, initExpr ast.Expr, span ast.Node, opK mapOpKind) {
			v := varFromObj(info.Defs[id])
			if id.Name == "_" || v == nil || v.IsField() || computeIIFETemps[v] {
				return
			}
			isCand, declIdent, declHash := classifyCandidateMapType(v.Type())
			if !isCand {
				return
			}
			if initExpr != nil && len(castCandsByFn) > 0 {
				if call, isCall := unwrapMapAdaptersDown(initExpr, info).(*ast.CallExpr); isCall && len(call.Args) == 0 && !call.Ellipsis.IsValid() {
					if callee, _ := calleeFunc(call, info); callee != nil {
						if cc := castCandsByFn[funcOrigin(callee)]; cc != nil && !cc.disqualified {
							cand := &mapCandidate{
								v: v, declFile: file, hasInit: true,
								isNullable: true, forceNullable: true,
								hasIdentityKind: declIdent || cc.hasIdentityKind,
								hasHashKind:     declHash || cc.hasHashKind,
								spellSites:      []typeSpellSite{{file: file, pos: span.Pos()}},
								j2gortSites:     []typeSpellSite{{file: file, pos: span.Pos()}},
							}
							cand.ops = append(cand.ops, mapRewriteOp{
								kind: opK, cand: cand, otherCand: cc, initKind: mapInitCastAccessorCall,
								file: file, spanNode: span, otherExpr: call, call: call,
							})
							consumedCastCalls[call] = true
							edges = append(edges, [2]*mapCandidate{cand, cc})
							varCands[v] = cand
							allCands = append(allCands, cand)
							return
						}
					}
				}
			}
			if initExpr != nil && hasPtrHelpers && len(varCands) > 0 && !isMethodRefIIFEBoundAssign(span, parents, info) {
				if sel, isSel := unwrapMapAdaptersDown(initExpr, info).(*ast.SelectorExpr); isSel && isIndirectStructFieldSel(sel, info) {
					if fieldCand := varCands[varFromObj(info.Uses[sel.Sel])]; fieldCand != nil && fieldCand.isField && !fieldCand.isCastAccessor && !fieldCand.disqualified {
						cand := &mapCandidate{
							v: v, declFile: file, hasInit: true,
							isSharedPtrLocal: true, isNullable: true, forceNullable: true,
							hasIdentityKind: declIdent || fieldCand.hasIdentityKind,
							hasHashKind:     declHash || fieldCand.hasHashKind,
							spellSites:      []typeSpellSite{{file: file, pos: span.Pos()}},
							j2gortSites:     []typeSpellSite{{file: file, pos: span.Pos()}},
						}
						cand.ops = append(cand.ops, mapRewriteOp{
							kind: opK, cand: cand, otherCand: fieldCand, initKind: mapInitSharedPtr,
							file: file, spanNode: span, otherExpr: sel,
						})
						edges = append(edges, [2]*mapCandidate{cand, fieldCand})
						sharedPtrEdges = append(sharedPtrEdges, [2]*mapCandidate{cand, fieldCand})
						varCands[v] = cand
						allCands = append(allCands, cand)
						return
					}
				}
			}
			var initK mapInitKind
			var argExpr ast.Expr
			var iife *mapIIFESpec
			var initIdent, initHash, hasNil bool
			if initExpr == nil || isNilIdent(initExpr, info) {
				if opK != mapOpDeclVarSpec {
					return
				}
				initK = mapInitNil
				hasNil = true
			} else if k, a, idK, hsK, ok := classifyMapInitExpr(initExpr, info); ok {
				initK, argExpr, initIdent, initHash = k, a, idK, hsK
			} else if spec, specHasNil, idK, hsK, ok := classifyMapIIFEInit(initExpr, info); ok {
				initK = mapInitIIFE
				iife = spec
				hasNil = specHasNil
				initIdent = idK
				initHash = hsK
			} else {
				return
			}
			cand := &mapCandidate{
				v: v, declFile: file, hasInit: !hasNil,
				hasNilAssign:    hasNil,
				hasIdentityKind: declIdent || initIdent,
				hasHashKind:     declHash || initHash,
				spellSites:      []typeSpellSite{{file: file, pos: span.Pos()}},
				j2gortSites:     []typeSpellSite{{file: file, pos: span.Pos()}},
			}
			cand.ops = append(cand.ops, mapRewriteOp{
				kind: opK, cand: cand, initKind: initK, iifeSpec: iife, file: file, spanNode: span, otherExpr: argExpr, call: asCallExpr(initExpr),
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
						addLocal(spec.Names[0], initExpr, spec, mapOpDeclVarSpec)
					}
				}
			case *ast.AssignStmt:
				if s.Tok == token.DEFINE && len(s.Lhs) == 1 && len(s.Rhs) == 1 && enclosingFuncBody(s, parents) != nil {
					if id, ok := s.Lhs[0].(*ast.Ident); ok {
						addLocal(id, s.Rhs[0], s, mapOpDeclShortVar)
					}
				}
			}
			return true
		})
	}

	if len(allCands) == 0 {
		return nil
	}

	// Upgrade local copy-map initializers whose source is a lowered candidate.
	upgradeCopy := func(cand *mapCandidate, initK *mapInitKind, otherExpr *ast.Expr, otherCand **mapCandidate) {
		if *initK != mapInitCopyJavaMap {
			return
		}
		unwrapped := unwrapMapAdaptersDown(*otherExpr, info)
		if srcCand := resolveMapVarAccess(unwrapped, varCands, info); srcCand != nil && srcCand != cand {
			*initK = mapInitCopyMap
			*otherExpr = unwrapped
			*otherCand = srcCand
			edges = append(edges, [2]*mapCandidate{cand, srcCand})
		}
	}
	for _, cand := range allCands {
		if len(cand.ops) == 0 {
			continue
		}
		op0 := &cand.ops[0]
		upgradeCopy(cand, &op0.initKind, &op0.otherExpr, &op0.otherCand)
		if op0.iifeSpec != nil {
			for i := range op0.iifeSpec.branches {
				br := &op0.iifeSpec.branches[i]
				upgradeCopy(cand, &br.initKind, &br.otherExpr, &br.otherCand)
			}
		}
	}

	if len(structOrder) > 0 {
		verifyMapStructInvariants(ctx, parentsByFile, varCands, structOrder, structFields)
	}

	funcsMutating := make(map[*mapCandidate]map[*types.Func]bool)
	markMut := func(cand *mapCandidate, node ast.Node, parents map[ast.Node]ast.Node) {
		if cand == nil {
			return
		}
		if fn := enclosingFuncObj(node, parents, info); fn != nil {
			if funcsMutating[cand] == nil {
				funcsMutating[cand] = make(map[*types.Func]bool)
			}
			funcsMutating[cand][fn] = true
		}
	}

	// Pass 2: Validate every use of every map candidate across the package.
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
			usedVar := varFromObj(info.Uses[id])
			if cand := varCands[usedVar]; cand != nil && !cand.disqualified {
				if !ctx.ShouldTransform(file) {
					cand.disqualified = true
				} else if body := getterBodyByField[usedVar]; body != nil && id.Pos() >= body.Pos() && id.Pos() <= body.End() {
					// Internal `return this_.<field>` inside its lowered 1-line getter method.
				} else {
					validateMapVarUse(id, file, cand, varCands, castCandByUnionField, paramCandsByFn, ctx, parents, varUseCounts, markMut, &edges)
				}
			}
			if cc := castCandByUnionField[usedVar]; cc != nil && !cc.disqualified {
				validateUnionFieldUse(id, cc, varCands, ctx.ShouldTransform(file), parents, info)
			}
			if fnObj, ok := info.Uses[id].(*types.Func); ok {
				fnOrig := funcOrigin(fnObj)
				if cc := castCandsByFn[fnOrig]; cc != nil && !cc.disqualified {
					validateCastAccessorCallSite(id, file, cc, varCands, consumedCastCalls, ctx, parents, varUseCounts, markMut, &edges)
				}
				if gfCand := getterFieldCands[fnOrig]; gfCand != nil && !gfCand.disqualified {
					validateMapFieldGetterCallSite(id, file, gfCand, varCands, paramCandsByFn, ctx, parents, varUseCounts, markMut, &edges)
				}
				if pcs := paramCandsByFn[fnOrig]; len(pcs) > 0 {
					validateMapParamFuncCallSite(id, file, pcs, varCands, getterFieldCands, ctx, parents)
				}
			}
			return true
		})
	}

	for _, e := range edges {
		unionMapCandidates(e[0], e[1])
	}

	sharedPeers := make(map[*mapCandidate][]*mapCandidate)
	for _, cand := range allCands {
		for i := range cand.ops {
			if cand.ops[i].initKind == mapInitSharedPtr && cand.ops[i].otherCand != nil {
				sharedPtrEdges = append(sharedPtrEdges, [2]*mapCandidate{cand, cand.ops[i].otherCand})
			}
		}
	}
	for _, se := range sharedPtrEdges {
		a, b := se[0], se[1]
		if a.isPtrAlias || b.isPtrAlias {
			a.disqualified = true
			b.disqualified = true
			continue
		}
		if a.disqualified || b.disqualified {
			continue
		}
		a.forceNullable = true
		b.forceNullable = true
		if !slices.Contains(sharedPeers[a], b) {
			sharedPeers[a] = append(sharedPeers[a], b)
		}
		if !slices.Contains(sharedPeers[b], a) {
			sharedPeers[b] = append(sharedPeers[b], a)
		}
	}
	if len(sharedPtrEdges) > 0 {
		changedMut := true
		for changedMut {
			changedMut = false
			for _, se := range sharedPtrEdges {
				a, b := se[0], se[1]
				for fn := range funcsMutating[a] {
					if funcsMutating[b] == nil {
						funcsMutating[b] = make(map[*types.Func]bool)
					}
					if !funcsMutating[b][fn] {
						funcsMutating[b][fn] = true
						changedMut = true
					}
				}
				for fn := range funcsMutating[b] {
					if funcsMutating[a] == nil {
						funcsMutating[a] = make(map[*types.Func]bool)
					}
					if !funcsMutating[a][fn] {
						funcsMutating[a][fn] = true
						changedMut = true
					}
				}
			}
		}
	}

	for _, cand := range allCands {
		if cand.disqualified || cand.defCount != 1 || !ctx.ShouldTransform(cand.declFile) {
			cand.disqualified = true
			continue
		}
		if cand.isCastAccessor && (!cand.hasUnionMapWrite || cand.inPkgCastCalls == 0) {
			cand.disqualified = true
			continue
		}
		if cand.isPtrAlias && (!cand.hasInit || cand.hasNilAssign || cand.hasNilCheck) {
			cand.disqualified = true
			continue
		}
		if !cand.isPtrAlias && (!cand.hasInit || cand.hasNilAssign || cand.hasNilCheck || cand.forceNullable) {
			cand.isNullable = true
		}
		if cand.isNullable && !hasPtrHelpers {
			for i := range cand.ops {
				k := cand.ops[i].initKind
				if k == mapInitCap || k == mapInitCopyMap || k == mapInitCopyJavaMap || k == mapInitCastAccessorCall {
					cand.disqualified = true
					break
				}
				if k == mapInitIIFE && cand.ops[i].iifeSpec != nil {
					for _, br := range cand.ops[i].iifeSpec.branches {
						if br.initKind == mapInitCap || br.initKind == mapInitCopyMap || br.initKind == mapInitCopyJavaMap {
							cand.disqualified = true
							break
						}
					}
				}
			}
			if cand.disqualified {
				continue
			}
		}
		if anyMapLoopMutatesCandidate(cand, sharedPeers[cand], funcsMutating[cand], info) {
			cand.disqualified = true
			continue
		}
		if cand.hasLocalLiveView && (cand.isPtr() || !cand.hasInit || cand.initAssignCount > 0 || cand.hasNonCtorAssign || cand.hasNilAssign || cand.hasNilCheck) {
			cand.disqualified = true
			continue
		}
	}

	var groupRoots []*mapCandidate
	groups := make(map[*mapCandidate][]*mapCandidate)
	for _, cand := range allCands {
		r := cand.root()
		if len(groups[r]) == 0 {
			groupRoots = append(groupRoots, r)
		}
		groups[r] = append(groups[r], cand)
	}

	var out []*mapCandidate
	for _, r := range groupRoots {
		members := groups[r]
		if !resolveAndValidateMapGroup(members, parentsByFile, ctx) {
			continue
		}
		out = append(out, members...)
	}
	slices.SortFunc(out, func(a, b *mapCandidate) int { return cmp.Compare(a.pos(), b.pos()) })
	return out
}

func resolveAndValidateMapGroup(members []*mapCandidate, parentsByFile map[*ast.File]map[ast.Node]ast.Node, ctx *PassContext) bool {
	hasIdentity := false
	hasHash := false
	var keyTypes, valReads []types.Type
	var valWrites []mapValWriteSample
	hasComputeOp := false

	for _, m := range members {
		if m.disqualified {
			return false
		}
		if m.hasIdentityKind {
			hasIdentity = true
		}
		if m.hasHashKind {
			hasHash = true
		}
		keyTypes = append(keyTypes, m.keyTypes...)
		valReads = append(valReads, m.valReadTypes...)
		valWrites = append(valWrites, m.valWrites...)
		for i := range m.ops {
			if m.ops[i].kind == mapOpComputeIfAbsent {
				hasComputeOp = true
			}
		}
	}
	// Never mix IdentityHashMap and HashMap/LinkedHashMap in the same equivalence group.
	if hasIdentity && hasHash {
		return false
	}
	if !hasIdentity && !hasHash {
		return false
	}
	if len(keyTypes) == 0 {
		for _, m := range members {
			for _, sc := range m.keySiteConstraints {
				keyTypes = append(keyTypes, sc.siteType)
			}
		}
	}
	if len(keyTypes) == 0 {
		return false
	}
	keyT, keyUnifyOK := resolveKeyOrSetElemGroupType(keyTypes, ctx.Pkg)
	if !keyUnifyOK {
		return false
	}
	var memberOwners []*types.Named
	for _, m := range members {
		if m.ownerStruct != nil {
			memberOwners = append(memberOwners, m.ownerStruct)
		}
	}
	flavor, keyOK := classifyMapKeyType(keyT, hasIdentity, memberOwners, ctx)
	if !keyOK {
		return false
	}
	if len(valReads) == 0 {
		hasConcreteWrite := false
		for _, w := range valWrites {
			if !w.isNilLiteral {
				hasConcreteWrite = true
				break
			}
		}
		if !hasConcreteWrite {
			for _, m := range members {
				for _, sc := range m.valSiteConstraints {
					valReads = append(valReads, sc.siteType)
				}
			}
		}
	}
	valT, valOK := resolveMapGroupValType(valReads, valWrites)
	if !valOK {
		return false
	}
	if flavor == mapFlavorCustom && len(members) == 1 && members[0].isField &&
		isCanonicalVarIdentityMapField(members[0].ownerStruct, members[0].v.Name(), keyT, valT) {
		flavor = mapFlavorOrdered
	}
	if hasComputeOp && !isComparableOrJavaIface(valT) {
		return false
	}
	// IdentityHashMap copy from general Java Map or KeysToImmutableSet on IdentityHashMap is disallowed.
	for _, m := range members {
		if typeContainsTypeParam(keyT) && m.isField && !typeOnlyUsesOwnerTypeParams(keyT, m.ownerStruct) {
			return false
		}
		if typeContainsTypeParam(valT) && m.isField && !typeOnlyUsesOwnerTypeParams(valT, m.ownerStruct) {
			return false
		}
		for _, sc := range m.keySiteConstraints {
			instK, ok := instantiateOwnerType(keyT, m.ownerStruct, sc.targs)
			if !ok || (!types.Identical(instK, sc.siteType) && !(isNonEmptyInterfaceType(instK) && (isNonEmptyInterfaceType(sc.siteType) || isPtrToStructType(sc.siteType)) && types.AssignableTo(sc.siteType, instK))) {
				return false
			}
		}
		for _, sc := range m.valSiteConstraints {
			instV, ok := instantiateOwnerType(valT, m.ownerStruct, sc.targs)
			if !ok || (!types.Identical(instV, sc.siteType) && !(isNonEmptyInterfaceType(instV) && types.AssignableTo(sc.siteType, instV))) {
				return false
			}
		}
		for i := range m.ops {
			if m.ops[i].isPtrAliasTargetInvalid() {
				return false
			}
			if hasIdentity {
				switch m.ops[i].kind {
				case mapOpImmMapCopy, mapOpKeysToImmSet, mapOpPutAllJavaMap:
					return false
				case mapOpDeclVarSpec, mapOpDeclShortVar, mapOpAssignInit:
					if m.ops[i].initKind == mapInitCopyJavaMap {
						return false
					}
					if m.ops[i].initKind == mapInitIIFE && m.ops[i].iifeSpec != nil {
						for _, br := range m.ops[i].iifeSpec.branches {
							if br.initKind == mapInitCopyJavaMap {
								return false
							}
						}
					}
				}
			}
		}
		if !canSpellTypeAtAllSites(keyT, m.spellSites, ctx) ||
			!canSpellTypeAtAllSites(valT, m.spellSites, ctx) ||
			!canSpellJ2gortAtAllSites(m.j2gortSites, ctx) {
			return false
		}
	}

	for _, m := range members {
		m.keyType = keyT
		m.valType = valT
		m.flavor = flavor
		m.isIdentityMap = hasIdentity
	}
	if flavor == mapFlavorOrdered {
		if len(members) == 1 {
			m := members[0]
			if m.isField && !m.isPtrAlias && !m.isCastAccessor && m.isNullable && !m.forceNullable &&
				canPromoteLazyFieldToValueOrderedMap(m, parentsByFile, ctx.Info) {
				m.isNullable = false
				m.isLazyValueMap = true
			}
		}
		if canPromoteGroupToPlainMap(members, parentsByFile, ctx) {
			for _, m := range members {
				m.flavor = mapFlavorPlain
			}
		}
	}
	return true
}

func (op *mapRewriteOp) isPtrAliasTargetInvalid() bool {
	return op.initKind == mapInitPtrAlias && (op.otherCand == nil || (op.otherCand.isPtr() && !op.otherCand.isParam))
}

func resolveMapGroupValType(reads []types.Type, writes []mapValWriteSample) (types.Type, bool) {
	var concreteWrites []types.Type
	hasNilWrite := false
	for _, w := range writes {
		if w.isNilLiteral {
			hasNilWrite = true
			continue
		}
		concreteWrites = append(concreteWrites, w.t)
	}
	var all []types.Type
	all = append(all, concreteWrites...)
	all = append(all, reads...)
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
		if isEmptyInterfaceType(t0) {
			return nil, false
		}
		if hasNilWrite && !isNilableMapValType(t0) {
			return nil, false
		}
		return t0, true
	}

	// Allow non-empty interface V when all reads expect V and every non-nil write is either V,
	// a sub-interface of V, or a provably non-nil allocation assignable to V (zero typed-nil hazard).
	var targetIface types.Type
	for _, rt := range reads {
		if !isNonEmptyInterfaceType(rt) {
			return nil, false
		}
		if targetIface == nil {
			targetIface = rt
		} else if !types.Identical(targetIface, rt) {
			return nil, false
		}
	}
	if targetIface == nil {
		return nil, false
	}
	for _, w := range writes {
		if w.isNilLiteral {
			continue
		}
		if types.Identical(w.t, targetIface) {
			continue
		}
		if isNonEmptyInterfaceType(w.t) && types.AssignableTo(w.t, targetIface) {
			continue
		}
		if w.nonNilAlloc && types.AssignableTo(w.t, targetIface) {
			continue
		}
		return nil, false
	}
	return targetIface, true
}

func isNilableMapValType(t types.Type) bool {
	if t == nil {
		return false
	}
	ut := types.Unalias(t)
	if _, isTP := ut.(*types.TypeParam); isTP {
		return false
	}
	switch ut.Underlying().(type) {
	case *types.Pointer, *types.Slice:
		return true
	case *types.Interface:
		return isNonEmptyInterfaceType(t)
	}
	return false
}

func isComparableOrJavaIface(t types.Type) bool {
	if t == nil || isInvalidOrUntyped(t) {
		return false
	}
	if types.Comparable(t) {
		return true
	}
	return isNonEmptyInterfaceType(t)
}

func satisfiesComparableConstraint(t types.Type) bool {
	if t == nil || isInvalidOrUntyped(t) {
		return false
	}
	if types.Comparable(t) {
		return true
	}
	if cmpObj := types.Universe.Lookup("comparable"); cmpObj != nil {
		if cmpIface, ok := cmpObj.Type().Underlying().(*types.Interface); ok {
			return types.Satisfies(t, cmpIface)
		}
	}
	return false
}

func isPtrToStructType(t types.Type) bool {
	if t == nil {
		return false
	}
	ptr, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return false
	}
	_, isStruct := types.Unalias(ptr.Elem()).Underlying().(*types.Struct)
	return isStruct
}

// resolveKeyOrSetElemGroupType resolves observed key/set-element types across a candidate group.
// In addition to exact types.Identical match, it unifies interface subtyping / pointer-to-struct
// subtypes to a common observed non-empty Java interface with IsJavaNil_() (e.g. ObjectType + JSType -> JSType,
// or *TemplateType + JSType -> JSType).
func resolveKeyOrSetElemGroupType(observed []types.Type, pkg *types.Package) (types.Type, bool) {
	if len(observed) == 0 {
		return nil, false
	}
	t0 := observed[0]
	allIdentical := true
	for _, kt := range observed[1:] {
		if !types.Identical(t0, kt) {
			allIdentical = false
			break
		}
	}
	if allIdentical {
		return t0, !isEmptyInterfaceType(t0) && !isInvalidOrUntyped(t0)
	}
	var chosen types.Type
	for _, candIface := range observed {
		ut := types.Unalias(candIface)
		if _, isNamed := ut.(*types.Named); !isNamed || !isNonEmptyInterfaceType(ut) || !hasZeroArgBoolMethod(candIface, pkg, "IsJavaNil_") {
			continue
		}
		allAssignable := true
		for _, kt := range observed {
			if types.Identical(kt, candIface) {
				continue
			}
			if (isNonEmptyInterfaceType(kt) || isPtrToStructType(kt)) && types.AssignableTo(kt, candIface) {
				continue
			}
			allAssignable = false
			break
		}
		if !allAssignable {
			continue
		}
		if chosen == nil {
			chosen = candIface
		} else if !types.Identical(chosen, candIface) {
			return nil, false
		}
	}
	if chosen != nil {
		return chosen, true
	}
	return nil, false
}

func hasValueStructImplInPkg(ifaceType types.Type, ctx *PassContext) bool {
	iface, ok := types.Unalias(ifaceType).Underlying().(*types.Interface)
	if !ok || iface.NumMethods() == 0 {
		return true
	}
	checkScope := func(sc *types.Scope) bool {
		if sc == nil {
			return false
		}
		for _, name := range sc.Names() {
			tn, ok := sc.Lookup(name).(*types.TypeName)
			if !ok || tn.IsAlias() {
				continue
			}
			t := tn.Type()
			named, ok := t.(*types.Named)
			if !ok {
				continue
			}
			if isInterfaceType(named) {
				continue
			}
			if named.TypeParams() != nil && named.TypeParams().Len() > 0 {
				targs := make([]types.Type, named.TypeParams().Len())
				for i := 0; i < len(targs); i++ {
					targs[i] = named.TypeParams().At(i)
				}
				if inst, err := types.Instantiate(nil, named, targs, false); err == nil {
					t = inst
				}
			}
			switch named.Underlying().(type) {
			case *types.Struct, *types.Slice, *types.Map, *types.Signature:
				if types.Implements(t, iface) {
					return true
				}
			case *types.Array:
				if !types.Comparable(t) && types.Implements(t, iface) {
					return true
				}
			}
		}
		return false
	}
	if named, ok := types.Unalias(ifaceType).(*types.Named); ok && named.Obj() != nil && named.Obj().Pkg() != nil {
		if checkScope(named.Obj().Pkg().Scope()) {
			return true
		}
	}
	if ctx != nil && ctx.Pkg != nil && checkScope(ctx.Pkg.Scope()) {
		return true
	}
	return false
}

// allNamedIfaceImplsInPkgUseIdentity reports whether ifaceType is a package-local Java interface
// with >=1 pointer-to-struct implementation in ctx.Pkg and EVERY implementation in ctx.Pkg uses
// Go identity equality (HasIdentityEquals_ / Enum / no custom Equals).
func allNamedIfaceImplsInPkgUseIdentity(ifaceType types.Type, ctx *PassContext) bool {
	if ctx == nil || ctx.Pkg == nil || ctx.Pkg.Scope() == nil {
		return false
	}
	namedIface, ok := types.Unalias(ifaceType).(*types.Named)
	if !ok || namedIface.Obj() == nil || namedIface.Obj().Pkg() != ctx.Pkg || !hasZeroArgBoolMethod(ifaceType, ctx.Pkg, "IsJavaNil_") {
		return false
	}
	targetIface := ifaceType
	if namedIface.Origin().TypeParams() != nil && namedIface.Origin().TypeParams().Len() > 0 {
		orig := namedIface.Origin()
		targs := make([]types.Type, orig.TypeParams().Len())
		for i := 0; i < len(targs); i++ {
			targs[i] = orig.TypeParams().At(i)
		}
		if inst, err := types.Instantiate(nil, orig, targs, false); err == nil {
			targetIface = inst
		}
	}
	iface, ok := targetIface.Underlying().(*types.Interface)
	if !ok || iface.NumMethods() == 0 {
		return false
	}
	implCount := 0
	sc := ctx.Pkg.Scope()
	for _, name := range sc.Names() {
		tn, ok := sc.Lookup(name).(*types.TypeName)
		if !ok || tn.IsAlias() {
			continue
		}
		named, ok := tn.Type().(*types.Named)
		if !ok || isInterfaceType(named) {
			continue
		}
		var candT types.Type = named
		if named.TypeParams() != nil && named.TypeParams().Len() > 0 {
			if namedIface.Origin().TypeParams() == nil || named.TypeParams().Len() != namedIface.Origin().TypeParams().Len() {
				// Cannot align arity conservatively; if method names overlap, reject.
				for i := 0; i < iface.NumMethods(); i++ {
					mName := iface.Method(i).Name()
					if strings.HasPrefix(mName, "Implements_") {
						if obj, _, _ := types.LookupFieldOrMethod(types.NewPointer(named), false, ctx.Pkg, mName); obj != nil {
							return false
						}
					}
				}
				continue
			}
			targs := make([]types.Type, named.TypeParams().Len())
			for i := 0; i < len(targs); i++ {
				targs[i] = namedIface.Origin().TypeParams().At(i)
			}
			inst, err := types.Instantiate(nil, named, targs, false)
			if err != nil {
				return false
			}
			candT = inst
		}
		ptrT := types.NewPointer(candT)
		if types.Implements(candT, iface) {
			return false
		}
		if !types.Implements(ptrT, iface) {
			continue
		}
		if _, isStruct := named.Underlying().(*types.Struct); !isStruct {
			return false
		}
		if hasNonIdentityCustomEquals(ptrT, ctx.Pkg) {
			return false
		}
		implCount++
	}
	return implCount > 0
}

func hasNonIdentityCustomEquals(t types.Type, pkg *types.Package) bool {
	if t == nil {
		return false
	}
	ut := types.Unalias(t)
	if b, ok := ut.Underlying().(*types.Basic); ok {
		info := b.Info()
		return (info&types.IsString) == 0 && (info&types.IsBoolean) == 0 && (info&types.IsInteger) == 0
	}
	hasM := func(name string) bool {
		obj, _, _ := types.LookupFieldOrMethod(t, true, pkg, name)
		_, ok := obj.(*types.Func)
		return ok
	}
	if hasM("HasIdentityEquals_") || hasM("Implements_j2gort_Ext_java_lang_Enum") {
		return false
	}
	return hasM("Equals")
}

func hasUnsafeTPInstantiation(tp *types.TypeParam, owners []*types.Named, ctx *PassContext, badArg func(types.Type) bool) bool {
	if ctx == nil || ctx.Info == nil || len(ctx.Info.Instances) == 0 {
		return false
	}
	tpOrigName := ""
	if tp.Obj() != nil {
		tpOrigName = tp.Obj().Name()
	}
	ownerSet := make(map[*types.Named]bool, len(owners))
	for _, ow := range owners {
		if ow != nil {
			ownerSet[ow.Origin()] = true
		}
	}
	// order-independent: searches for existence of any disqualifying concrete type argument across package instantiations
	for id, inst := range ctx.Info.Instances {
		if inst.TypeArgs == nil || inst.TypeArgs.Len() == 0 {
			continue
		}
		obj := ctx.Info.Uses[id]
		if obj == nil {
			continue
		}
		var tparams *types.TypeParamList
		matched := false
		switch o := obj.(type) {
		case *types.TypeName:
			if named := namedOrigin(o.Type()); named != nil {
				tparams = named.TypeParams()
				if len(ownerSet) == 0 || ownerSet[named] {
					matched = true
				}
			}
		case *types.Func:
			fnOrig := funcOrigin(o)
			if sig, ok := fnOrig.Type().(*types.Signature); ok && sig.TypeParams() != nil {
				tparams = sig.TypeParams()
				if len(ownerSet) == 0 {
					matched = true
				} else {
					for ow := range ownerSet { // order-independent: boolean match check on constructor name prefix
						owBase := strings.TrimSuffix(ow.Obj().Name(), "Base")
						if strings.HasPrefix(fnOrig.Name(), "New"+ow.Obj().Name()) ||
							strings.HasPrefix(fnOrig.Name(), "Init"+ow.Obj().Name()) ||
							strings.HasPrefix(fnOrig.Name(), "New"+owBase) ||
							strings.HasPrefix(fnOrig.Name(), "Init"+owBase) ||
							strings.HasPrefix(fnOrig.Name(), owBase+"_") {
							matched = true
							break
						}
					}
				}
			}
		}
		if !matched || tparams == nil || tparams.Len() != inst.TypeArgs.Len() {
			continue
		}
		for i := 0; i < tparams.Len(); i++ {
			candTP := tparams.At(i)
			if candTP == tp || (len(ownerSet) > 0 && candTP.Obj() != nil && candTP.Obj().Name() == tpOrigName) {
				targ := types.Unalias(inst.TypeArgs.At(i))
				if _, stillTP := targ.(*types.TypeParam); stillTP {
					continue
				}
				if badArg(inst.TypeArgs.At(i)) {
					return true
				}
			}
		}
	}
	return false
}

func classifyMapKeyType(t types.Type, isIdentity bool, owners []*types.Named, ctx *PassContext) (mapFlavor, bool) {
	var pkg *types.Package
	if ctx != nil {
		pkg = ctx.Pkg
	}
	if t == nil || isInvalidOrUntyped(t) || isEmptyInterfaceType(t) || !satisfiesComparableConstraint(t) {
		return 0, false
	}
	identFlavor := mapFlavorOrdered
	if j2gortHasExport(ctx, "IdentityOrderedMap") {
		identFlavor = mapFlavorIdentity
	}
	ut := types.Unalias(t)
	if tp, isTP := ut.(*types.TypeParam); isTP {
		if isIdentity {
			unsafeIdent := hasUnsafeTPInstantiation(tp, owners, ctx, func(targ types.Type) bool {
				if !satisfiesComparableConstraint(targ) {
					return true
				}
				u := types.Unalias(targ)
				if b, ok := u.Underlying().(*types.Basic); ok {
					return (b.Info() & types.IsString) != 0
				}
				if isNonEmptyInterfaceType(u) {
					return !hasZeroArgBoolMethod(targ, pkg, "IsJavaNil_") || hasValueStructImplInPkg(u, ctx)
				}
				return !isPtrToStructType(targ)
			})
			if unsafeIdent {
				return 0, false
			}
			return identFlavor, true
		}
		if implementsFastCustomKey(t, pkg) {
			return mapFlavorCustom, true
		}
		if hasUnsafeTPInstantiation(tp, owners, ctx, func(targ types.Type) bool {
			return !satisfiesComparableConstraint(targ) || hasNonIdentityCustomEquals(targ, pkg) || (isNonEmptyInterfaceType(targ) && hasValueStructImplInPkg(targ, ctx))
		}) {
			return 0, false
		}
		return mapFlavorOrdered, true
	}
	if isIdentity {
		if b, ok := ut.Underlying().(*types.Basic); ok {
			info := b.Info()
			if (info&types.IsBoolean) != 0 || (info&types.IsInteger) != 0 {
				return identFlavor, true
			}
			return 0, false
		}
		if isPtrToStructType(ut) {
			return identFlavor, true
		}
		if _, isNamed := ut.(*types.Named); isNamed && isNonEmptyInterfaceType(ut) && hasZeroArgBoolMethod(t, pkg, "IsJavaNil_") && !hasValueStructImplInPkg(ut, ctx) {
			return identFlavor, true
		}
		return 0, false
	}

	if b, ok := ut.Underlying().(*types.Basic); ok {
		info := b.Info()
		if (info&types.IsString) != 0 || (info&types.IsBoolean) != 0 || (info&types.IsInteger) != 0 {
			return mapFlavorOrdered, true
		}
		return 0, false
	}
	if ptr, ok := ut.(*types.Pointer); ok {
		if _, isStruct := types.Unalias(ptr.Elem()).Underlying().(*types.Struct); !isStruct {
			return 0, false
		}
		hasM := func(name string) bool {
			obj, _, _ := types.LookupFieldOrMethod(t, true, pkg, name)
			_, ok := obj.(*types.Func)
			return ok
		}
		if hasM("HasIdentityEquals_") || hasM("Implements_j2gort_Ext_java_lang_Enum") || !hasM("Equals") {
			return mapFlavorOrdered, true
		}
		if implementsFastCustomKey(t, pkg) {
			return mapFlavorCustom, true
		}
		return 0, false
	}
	if named, isNamed := ut.(*types.Named); isNamed {
		if _, isStruct := named.Underlying().(*types.Struct); isStruct {
			if hasZeroArgBoolMethod(t, pkg, "HasIdentityEquals_") || (named.Obj() != nil && named.Obj().Name() == "ColorId") {
				return mapFlavorOrdered, true
			}
			if implementsFastCustomKey(t, pkg) {
				return mapFlavorCustom, true
			}
			return 0, false
		}
		if isNonEmptyInterfaceType(ut) && !hasValueStructImplInPkg(ut, ctx) {
			if implementsFastCustomKey(t, pkg) {
				return mapFlavorCustom, true
			}
			if allNamedIfaceImplsInPkgUseIdentity(ut, ctx) {
				return mapFlavorOrdered, true
			}
		}
	}
	return 0, false
}

func hasZeroArgBoolMethod(t types.Type, pkg *types.Package, name string) bool {
	obj, _, _ := types.LookupFieldOrMethod(t, false, pkg, name)
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

func implementsFastCustomKey(t types.Type, pkg *types.Package) bool {
	if !hasZeroArgBoolMethod(t, pkg, "IsJavaNil_") {
		return false
	}
	hcObj, _, _ := types.LookupFieldOrMethod(t, false, pkg, "HashCode")
	hcFn, ok := hcObj.(*types.Func)
	if !ok {
		return false
	}
	hcSig, ok := hcFn.Type().(*types.Signature)
	if !ok || hcSig.Params().Len() != 0 || hcSig.Results().Len() != 1 {
		return false
	}
	if b, ok := types.Unalias(hcSig.Results().At(0).Type()).(*types.Basic); !ok || b.Kind() != types.Int32 {
		return false
	}
	eqObj, _, _ := types.LookupFieldOrMethod(t, false, pkg, "Equals")
	eqFn, ok := eqObj.(*types.Func)
	if !ok {
		return false
	}
	eqSig, ok := eqFn.Type().(*types.Signature)
	if !ok || eqSig.Params().Len() != 1 || eqSig.Results().Len() != 1 {
		return false
	}
	if !isEmptyInterfaceType(eqSig.Params().At(0).Type()) {
		return false
	}
	b, ok := types.Unalias(eqSig.Results().At(0).Type()).(*types.Basic)
	return ok && b.Kind() == types.Bool
}

func verifyMapStructInvariants(
	ctx *PassContext,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
	varCands map[*types.Var]*mapCandidate,
	structOrder []*types.Named,
	structFields map[*types.Named][]*mapCandidate,
) {
	info := ctx.Info
	disqualify := func(named *types.Named) {
		for _, c := range structFields[named] {
			c.disqualified = true
		}
	}

	isNonNilStructFieldInitRhs := func(rhs ast.Expr) bool {
		if _, _, _, _, ok := classifyMapInitExpr(rhs, info); ok {
			return true
		}
		if _, hasNil, _, _, ok := classifyMapIIFEInit(rhs, info); ok && !hasNil {
			return true
		}
		if sel, ok := stripParens(rhs).(*ast.SelectorExpr); ok && isIndirectStructFieldSel(sel, info) {
			if srcCand := varCands[varFromObj(info.Uses[sel.Sel])]; srcCand != nil && srcCand.isField {
				return true
			}
		}
		if id, ok := stripParens(rhs).(*ast.Ident); ok {
			if srcCand := varCands[varFromObj(info.Uses[id])]; srcCand != nil && srcCand.isParam {
				return true
			}
		}
		return false
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
							if cand := varCands[varFromObj(info.Uses[sel.Sel])]; cand != nil && cand.isField && isNonNilStructFieldInitRhs(as.Rhs[0]) {
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
								if id, ok := stripParens(sel.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == recvVar && varFromObj(info.Uses[sel.Sel]) == cand.v && isNonNilStructFieldInitRhs(s.Rhs[0]) {
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

func validateMapVarUse(
	id *ast.Ident,
	file *ast.File,
	cand *mapCandidate,
	varCands map[*types.Var]*mapCandidate,
	castCandByUnionField map[*types.Var]*mapCandidate,
	paramCandsByFn map[*types.Func][]*mapCandidate,
	ctx *PassContext,
	parents map[ast.Node]ast.Node,
	varUseCounts map[*types.Var]int,
	markMut func(*mapCandidate, ast.Node, map[ast.Node]ast.Node),
	edges *[][2]*mapCandidate,
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

	// Assignment to candidate (`m = ...` or `this_.f = ...`).
	if as, ok := parents[accessExpr].(*ast.AssignStmt); ok {
		for _, lhs := range as.Lhs {
			if stripParens(lhs) == accessExpr {
				if as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
					cand.disqualified = true
					return
				}
				if cand.isField && !isPureFieldSelector(accessExpr.(*ast.SelectorExpr), info) {
					cand.disqualified = true
					return
				}
				rhs := as.Rhs[0]
				if isNilIdent(rhs, info) {
					if cand.isPtrAlias {
						cand.disqualified = true
						return
					}
					cand.initAssignCount++
					cand.hasNilAssign = true
					cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
					cand.ops = append(cand.ops, mapRewriteOp{
						kind: mapOpAssignInit, cand: cand, initKind: mapInitNil, file: file, spanNode: as, stmt: as, accessExpr: accessExpr,
					})
					if !cand.isField || !isInInitFieldsOrCtor(as, cand.ownerStruct, parents, info) {
						cand.hasNonCtorAssign = true
						markMut(cand, as, parents)
					}
					return
				}
				if initK, argExpr, initIdent, initHash, initOK := classifyMapInitExpr(rhs, info); initOK {
					if cand.isPtrAlias {
						cand.disqualified = true
						return
					}
					cand.initAssignCount++
					if initIdent {
						cand.hasIdentityKind = true
					}
					if initHash {
						cand.hasHashKind = true
					}
					var srcCand *mapCandidate
					if initK == mapInitCopyJavaMap {
						unwrapped := unwrapMapAdaptersDown(argExpr, info)
						if sc := resolveMapVarAccess(unwrapped, varCands, info); sc != nil && (sc != cand || (cand.isField && !sameAccessPath(accessExpr, unwrapped, cand.v, info))) {
							initK = mapInitCopyMap
							argExpr = unwrapped
							srcCand = sc
							*edges = append(*edges, [2]*mapCandidate{cand, sc})
						}
					}
					cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
					cand.ops = append(cand.ops, mapRewriteOp{
						kind: mapOpAssignInit, cand: cand, otherCand: srcCand, initKind: initK, file: file, spanNode: as, stmt: as, accessExpr: accessExpr, otherExpr: argExpr, call: asCallExpr(rhs),
					})
					if !cand.isField || !isInInitFieldsOrCtor(as, cand.ownerStruct, parents, info) {
						cand.hasNonCtorAssign = true
						markMut(cand, as, parents)
					}
					return
				}
				if iife, specHasNil, initIdent, initHash, iifeOK := classifyMapIIFEInit(rhs, info); iifeOK {
					if cand.isPtrAlias {
						cand.disqualified = true
						return
					}
					cand.initAssignCount++
					if specHasNil {
						cand.hasNilAssign = true
					}
					if initIdent {
						cand.hasIdentityKind = true
					}
					if initHash {
						cand.hasHashKind = true
					}
					for i := range iife.branches {
						br := &iife.branches[i]
						if br.initKind == mapInitCopyJavaMap {
							unwrapped := unwrapMapAdaptersDown(br.otherExpr, info)
							if sc := resolveMapVarAccess(unwrapped, varCands, info); sc != nil && (sc != cand || (cand.isField && !sameAccessPath(accessExpr, unwrapped, cand.v, info))) {
								br.initKind = mapInitCopyMap
								br.otherExpr = unwrapped
								br.otherCand = sc
								*edges = append(*edges, [2]*mapCandidate{cand, sc})
							}
						}
					}
					cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
					cand.ops = append(cand.ops, mapRewriteOp{
						kind: mapOpAssignInit, cand: cand, initKind: mapInitIIFE, iifeSpec: iife, file: file, spanNode: as, stmt: as, accessExpr: accessExpr,
					})
					if !cand.isField || !isInInitFieldsOrCtor(as, cand.ownerStruct, parents, info) {
						cand.hasNonCtorAssign = true
						markMut(cand, as, parents)
					}
					return
				}
				// Pointer-alias struct field initialization (`this_.resultMap = this_.result.reachingDef`
				// or constructor/init map parameter handoff `this_.idToColor = idToColor`).
				if cand.isField && cand.hasInit && !cand.hasNilAssign && !cand.hasNilCheck && cand.initAssignCount == 0 {
					if rhsSel, ok := stripParens(rhs).(*ast.SelectorExpr); ok && isIndirectStructFieldSel(rhsSel, info) {
						if srcCand := varCands[varFromObj(info.Uses[rhsSel.Sel])]; srcCand != nil && srcCand.isField && srcCand != cand && !srcCand.isPtrAlias {
							cand.initAssignCount++
							cand.isPtrAlias = true
							*edges = append(*edges, [2]*mapCandidate{cand, srcCand})
							cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
							cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
							cand.ops = append(cand.ops, mapRewriteOp{
								kind: mapOpAssignInit, cand: cand, otherCand: srcCand, initKind: mapInitPtrAlias,
								file: file, spanNode: as, stmt: as, accessExpr: accessExpr, otherExpr: rhsSel,
							})
							return
						}
					}
					if rhsId, ok := stripParens(rhs).(*ast.Ident); ok {
						if srcCand := varCands[varFromObj(info.Uses[rhsId])]; srcCand != nil && srcCand.isParam && !srcCand.isPtrAlias && !srcCand.isSharedPtrLocal && srcCand != cand {
							cand.initAssignCount++
							cand.isPtrAlias = true
							*edges = append(*edges, [2]*mapCandidate{cand, srcCand})
							cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
							cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
							cand.ops = append(cand.ops, mapRewriteOp{
								kind: mapOpAssignInit, cand: cand, otherCand: srcCand, initKind: mapInitPtrAlias,
								file: file, spanNode: as, stmt: as, accessExpr: accessExpr, otherExpr: rhsId,
							})
							return
						}
					}
				}
				// Shared nullable pointer assignment between struct field and local lazy-init variable
				// (e.g. `var map_ = this_.properties; if map_ == nil { map_ = New...(4); ...; this_.properties = map_ }`
				// or `var m = NewHashMap(16); ...; r.map_ = m`).
				if j2gortHasExport(ctx, "NewOrderedMapPtr") {
					rhsUnwrapped := unwrapMapAdaptersDown(rhs, info)
					if cand.isField && !cand.isCastAccessor && !cand.isPtrAlias {
						if rhsId, ok := rhsUnwrapped.(*ast.Ident); ok {
							if srcCand := varCands[varFromObj(info.Uses[rhsId])]; srcCand != nil && !srcCand.isField && !srcCand.isParam && !srcCand.isPtrAlias && srcCand != cand {
								cand.initAssignCount++
								cand.forceNullable = true
								srcCand.forceNullable = true
								*edges = append(*edges, [2]*mapCandidate{cand, srcCand})
								cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
								cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
								cand.ops = append(cand.ops, mapRewriteOp{
									kind: mapOpAssignInit, cand: cand, otherCand: srcCand, initKind: mapInitSharedPtr,
									file: file, spanNode: as, stmt: as, accessExpr: accessExpr, otherExpr: rhsId,
								})
								if !isInInitFieldsOrCtor(as, cand.ownerStruct, parents, info) {
									cand.hasNonCtorAssign = true
									markMut(cand, as, parents)
								}
								return
							}
						}
					} else if cand.isSharedPtrLocal && !cand.isPtrAlias {
						if rhsSel, ok := rhsUnwrapped.(*ast.SelectorExpr); ok && isIndirectStructFieldSel(rhsSel, info) {
							if srcCand := varCands[varFromObj(info.Uses[rhsSel.Sel])]; srcCand != nil && srcCand.isField && !srcCand.isCastAccessor && !srcCand.isPtrAlias {
								cand.initAssignCount++
								cand.forceNullable = true
								srcCand.forceNullable = true
								*edges = append(*edges, [2]*mapCandidate{cand, srcCand})
								cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: as.Pos()})
								cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
								cand.ops = append(cand.ops, mapRewriteOp{
									kind: mapOpAssignInit, cand: cand, otherCand: srcCand, initKind: mapInitSharedPtr,
									file: file, spanNode: as, stmt: as, accessExpr: accessExpr, otherExpr: rhsSel,
								})
								markMut(cand, as, parents)
								return
							}
						}
					}
				}
				cand.disqualified = true
				return
			}
		}
	}

	// Blank-identifier suppression `_ = m` for locals/params.
	if !cand.isField {
		if as, ok := parents[id].(*ast.AssignStmt); ok && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && stripParens(as.Rhs[0]) == id {
			if lhsId, ok := stripParens(as.Lhs[0]).(*ast.Ident); ok && lhsId.Name == "_" {
				return
			}
		}
	}

	validateMapValueConsumer(accessExpr, file, cand, varCands, castCandByUnionField, paramCandsByFn, ctx, parents, info, varUseCounts, markMut, edges)
}

func validateMapValueConsumer(
	accessExpr ast.Expr,
	file *ast.File,
	cand *mapCandidate,
	varCands map[*types.Var]*mapCandidate,
	castCandByUnionField map[*types.Var]*mapCandidate,
	paramCandsByFn map[*types.Func][]*mapCandidate,
	ctx *PassContext,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
	markMut func(*mapCandidate, ast.Node, map[ast.Node]ast.Node),
	edges *[][2]*mapCandidate,
) {
	// Check if `accessExpr` is the RHS of a pointer-alias struct field init `this_.resultMap = this_.result.reachingDef`
	// or a union `any` field handoff `this_.refsForNode = refsForNodeMap`.
	pTop := unwrapParensUp(accessExpr, parents)
	if as, ok := parents[pTop].(*ast.AssignStmt); ok && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && as.Rhs[0] == pTop {
		if lhsSel, ok := stripParens(as.Lhs[0]).(*ast.SelectorExpr); ok && isPureFieldSelector(lhsSel, info) {
			if dstCand := varCands[varFromObj(info.Uses[lhsSel.Sel])]; dstCand != nil && dstCand.isField && dstCand.hasInit && dstCand != cand && (cand.isField || cand.isParam) && !cand.isPtrAlias && !cand.isSharedPtrLocal && !dstCand.hasNilAssign && !dstCand.hasNilCheck {
				*edges = append(*edges, [2]*mapCandidate{cand, dstCand})
				return
			}
		}
	}

	adaptedTop := unwrapMapAdaptersUp(accessExpr, parents, info)
	if cand.isField && !cand.isCastAccessor && !cand.isPtrAlias {
		switch p := parents[adaptedTop].(type) {
		case *ast.ValueSpec:
			if len(p.Names) == 1 && len(p.Values) == 1 && p.Values[0] == adaptedTop {
				if dstCand := varCands[varFromObj(info.Defs[p.Names[0]])]; dstCand != nil && dstCand.isSharedPtrLocal {
					*edges = append(*edges, [2]*mapCandidate{cand, dstCand})
					return
				}
			}
		case *ast.AssignStmt:
			if len(p.Lhs) == 1 && len(p.Rhs) == 1 && p.Rhs[0] == adaptedTop {
				if p.Tok == token.DEFINE {
					if lhsId, ok := p.Lhs[0].(*ast.Ident); ok {
						if dstCand := varCands[varFromObj(info.Defs[lhsId])]; dstCand != nil && dstCand.isSharedPtrLocal {
							*edges = append(*edges, [2]*mapCandidate{cand, dstCand})
							return
						}
					}
				} else if p.Tok == token.ASSIGN {
					if lhsId, ok := stripParens(p.Lhs[0]).(*ast.Ident); ok {
						if dstCand := varCands[varFromObj(info.Uses[lhsId])]; dstCand != nil && dstCand.isSharedPtrLocal {
							*edges = append(*edges, [2]*mapCandidate{cand, dstCand})
							return
						}
					}
				}
			}
		}
	} else if !cand.isField && !cand.isParam && !cand.isPtrAlias && j2gortHasExport(ctx, "NewOrderedMapPtr") {
		if as, ok := parents[adaptedTop].(*ast.AssignStmt); ok && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && as.Rhs[0] == adaptedTop {
			if lhsSel, ok := stripParens(as.Lhs[0]).(*ast.SelectorExpr); ok && isPureFieldSelector(lhsSel, info) {
				if dstCand := varCands[varFromObj(info.Uses[lhsSel.Sel])]; dstCand != nil && dstCand.isField && !dstCand.isCastAccessor && !dstCand.isPtrAlias && dstCand != cand {
					*edges = append(*edges, [2]*mapCandidate{cand, dstCand})
					return
				}
			}
		}
	}
	if !cand.isField && len(castCandByUnionField) > 0 {
		if as, ok := parents[adaptedTop].(*ast.AssignStmt); ok && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && as.Rhs[0] == adaptedTop {
			if lhsSel, ok := stripParens(as.Lhs[0]).(*ast.SelectorExpr); ok && isPureFieldSelector(lhsSel, info) {
				if cc := castCandByUnionField[varFromObj(info.Uses[lhsSel.Sel])]; cc != nil && !cc.disqualified {
					cand.forceNullable = true
					cc.hasUnionMapWrite = true
					*edges = append(*edges, [2]*mapCandidate{cand, cc})
					cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: as.Pos()})
					cand.ops = append(cand.ops, mapRewriteOp{
						kind: mapOpAssignUnionField, cand: cand, otherCand: cc, file: file,
						spanNode: adaptedTop, outerExpr: adaptedTop, accessExpr: accessExpr,
					})
					return
				}
			}
		}
	}

	// Check nil comparisons `m == nil` / `m != nil`.
	if bin, ok := parents[pTop].(*ast.BinaryExpr); ok && (bin.Op == token.EQL || bin.Op == token.NEQ) {
		var otherSide ast.Expr
		if bin.X == pTop {
			otherSide = bin.Y
		} else if bin.Y == pTop {
			otherSide = bin.X
		}
		if otherSide != nil && isNilIdent(otherSide, info) && !cand.isPtrAlias {
			cand.hasNilCheck = true
			cand.ops = append(cand.ops, mapRewriteOp{
				kind: mapOpNilCheck, cand: cand, file: file, spanNode: bin, outerExpr: bin,
				accessExpr: accessExpr, isNegated: bin.Op == token.NEQ,
			})
			return
		}
	}

	if outerSink, isPtrWrapped, ok := matchFormattingSetSink(pTop, parents, info); ok {
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: accessExpr.Pos()})
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpFormatSink, cand: cand, file: file, spanNode: outerSink, outerExpr: outerSink, accessExpr: accessExpr, unwrapsValBox: isPtrWrapped,
		})
		return
	}

	// Check if `adaptedTop` is the map binding `mTmp := <adaptedTop>` of a j2go ComputeIfAbsent IIFE.
	if validateComputeIfAbsentIIFE(adaptedTop, accessExpr, file, cand, parents, info, varUseCounts, markMut) {
		return
	}

	// Check if `adaptedTop` is `bound := <adaptedTop>` inside a `Map::get` method-reference IIFE.
	if j2gortHasExport(ctx, "OrderedMapGetNodeAny") && validateMethodRefGetIIFE(adaptedTop, accessExpr, file, cand, parents, info, varUseCounts) {
		return
	}

	// Check static j2gort helpers, candidate map parameters, and constructors/adapters consuming `adaptedTop`.
	if staticCall, ok := parents[adaptedTop].(*ast.CallExpr); ok && !staticCall.Ellipsis.IsValid() && staticCall.Fun != adaptedTop {
		if len(paramCandsByFn) > 0 {
			if callee, _ := calleeFunc(staticCall, info); callee != nil {
				if pcs := paramCandsByFn[funcOrigin(callee)]; len(pcs) > 0 {
					for _, pc := range pcs {
						if pc.paramIdx < len(staticCall.Args) && staticCall.Args[pc.paramIdx] == adaptedTop && !pc.disqualified {
							*edges = append(*edges, [2]*mapCandidate{cand, pc})
							cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: staticCall.Pos()})
							cand.ops = append(cand.ops, mapRewriteOp{
								kind: mapOpPassParamArg, cand: cand, otherCand: pc, file: file,
								spanNode: adaptedTop, outerExpr: adaptedTop, accessExpr: accessExpr,
							})
							markMut(cand, staticCall, parents)
							return
						}
					}
				}
			}
		}
		if fnName, _, _, isJ2 := isJ2gortCall(staticCall, info); isJ2 {
			switch {
			case len(staticCall.Args) >= 2 && staticCall.Args[0] == adaptedTop:
				if validateStaticMapOp(fnName, staticCall, accessExpr, file, cand, parents, info, varUseCounts, markMut) {
					return
				}

			case len(staticCall.Args) == 1 && staticCall.Args[0] == adaptedTop:
				switch fnName {
				case "NewMapKeyCursor", "NewMapKeyStrCursor", "NewMapValCursor":
					if unwrapEnclosingCheckNotNulls(accessExpr, parents, info) != adaptedTop {
						cand.disqualified = true
						return
					}
					validateMapIterCursorLoop(fnName, staticCall, accessExpr, file, cand, parents, info, varUseCounts)
					return

				case "NewMapEntryCursor":
					if unwrapEnclosingCheckNotNulls(accessExpr, parents, info) != adaptedTop {
						cand.disqualified = true
						return
					}
					validateMapEntryCursorLoop(staticCall, accessExpr, file, cand, parents, info, varUseCounts)
					return

				case "Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map":
					if _, isES := parents[staticCall].(*ast.ExprStmt); !isES {
						cand.ops = append(cand.ops, mapRewriteOp{
							kind: mapOpImmMapCopy, cand: cand, file: file, spanNode: staticCall, outerExpr: staticCall, accessExpr: accessExpr,
						})
						return
					}

				case "NewExt_java_util_LinkedHashMap_1_Ext_java_util_Map",
					"NewExt_java_util_HashMap_1_Ext_java_util_Map",
					"Ext_com_google_common_collect_Maps_NewLinkedHashMap_1_Ext_java_util_Map",
					"Ext_com_google_common_collect_Maps_NewHashMap_1_Ext_java_util_Map":
					if isMapInitTargetCandidate(staticCall, varCands, parents, info) {
						return
					}
				}
			}
		}

		// Source argument of `dstMap.PutAll(AsIface(cand))`.
		if len(staticCall.Args) == 1 && staticCall.Args[0] == adaptedTop {
			if mSel, ok := stripParens(staticCall.Fun).(*ast.SelectorExpr); ok && mSel.Sel.Name == "PutAll" {
				recvExpr := unwrapCheckNotNullArg(mSel.X, info)
				if dstCand := resolveMapVarAccess(recvExpr, varCands, info); dstCand != nil && (dstCand != cand || (cand.isField && !sameAccessPath(accessExpr, recvExpr, cand.v, info))) {
					return
				}
			}
		}
		// Source argument of `otherMap.Equals(PtrToAny(cand))`.
		if len(staticCall.Args) == 1 && staticCall.Args[0] == adaptedTop {
			if mSel, ok := stripParens(staticCall.Fun).(*ast.SelectorExpr); ok && mSel.Sel.Name == "Equals" {
				if dstCand := resolveMapVarAccess(unwrapCheckNotNullArg(mSel.X, info), varCands, info); dstCand != nil && (dstCand != cand || cand.isField) {
					return
				}
			}
		}
	}
	// Also check `otherMap.Equals(PtrToAny(cand))` when wrapped in `PtrToAny` / `ToAny`.
	if ptrWrap, ok := parents[pTop].(*ast.CallExpr); ok && len(ptrWrap.Args) == 1 && !ptrWrap.Ellipsis.IsValid() && ptrWrap.Args[0] == pTop {
		if fnName, _, _, isJ2 := isJ2gortCall(ptrWrap, info); isJ2 && (fnName == "PtrToAny" || fnName == "ToAny") {
			wTop := unwrapParensUp(ptrWrap, parents)
			if eqCall, ok := parents[wTop].(*ast.CallExpr); ok && len(eqCall.Args) == 1 && !eqCall.Ellipsis.IsValid() && eqCall.Args[0] == wTop {
				if mSel, ok := stripParens(eqCall.Fun).(*ast.SelectorExpr); ok && mSel.Sel.Name == "Equals" {
					if dstCand := resolveMapVarAccess(unwrapCheckNotNullArg(mSel.X, info), varCands, info); dstCand != nil && (dstCand != cand || cand.isField) {
						return
					}
				}
			}
		}
	}

	// Direct method call `m.<Method>(...)`.
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
	case "Size", "IsEmpty", "String":
		if inExprStmt || len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		if methodSel.Sel.Name == "String" {
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		}
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpMethodDirect, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, methodsSel: methodSel.Sel.Name,
		})

	case "Clear":
		if !inExprStmt || len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpMethodDirect, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, methodsSel: "Clear",
		})
		markMut(cand, mCall, parents)

	case "ContainsKey":
		if inExprStmt || len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		unboxedK, keyT, ok := unwrapMapKeyArg(mCall.Args[0], cand.hasIdentityKind, info)
		if !ok || !cand.addKeyType(keyT, accessExpr, info) {
			cand.disqualified = true
			return
		}
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpContainsKey, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, call: mCall, keyArgIdx: 0, unwrapsKeyBox: unboxedK != stripParens(mCall.Args[0]),
		})

	case "Get":
		if inExprStmt || len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		unboxedK, keyT, okK := unwrapMapKeyArg(mCall.Args[0], cand.hasIdentityKind, info)
		outerExpr, valT, okV := enclosingMapValUnboxExpr(mCall, parents, info)
		_, outerES := parents[outerExpr].(*ast.ExprStmt)
		if (!okK || isEmptyInterfaceType(keyT)) && isEmptyInterfaceType(info.Types[mCall.Args[0]].Type) && j2gortHasExport(ctx, "OrderedMapGetAny") {
			if !okV || outerES || (!isEmptyInterfaceType(valT) && !cand.addValReadType(valT, accessExpr, info)) {
				cand.disqualified = true
				return
			}
			cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
			cand.ops = append(cand.ops, mapRewriteOp{
				kind: mapOpGetAny, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
				accessExpr: accessExpr, call: mCall, keyArgIdx: 0,
			})
			return
		}
		if !okK || !okV || outerES || !cand.addKeyType(keyT, accessExpr, info) || !cand.addValReadType(valT, accessExpr, info) {
			cand.disqualified = true
			return
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outerExpr.Pos()})
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpGet, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr,
			accessExpr: accessExpr, call: mCall, keyArgIdx: 0, unwrapsKeyBox: unboxedK != stripParens(mCall.Args[0]),
		})

	case "GetOrDefault":
		if inExprStmt || len(mCall.Args) != 2 {
			cand.disqualified = true
			return
		}
		unboxedK, keyT, okK := unwrapMapKeyArg(mCall.Args[0], cand.hasIdentityKind, info)
		unboxedDef, defT, defNonNil, defIsNil, okD := unwrapMapValArg(mCall.Args[1], info)
		outerExpr, valT, okV := enclosingMapValUnboxExpr(mCall, parents, info)
		_, outerES := parents[outerExpr].(*ast.ExprStmt)
		if !okK || !okD || !okV || outerES || !cand.addKeyType(keyT, accessExpr, info) || !cand.addValReadType(valT, accessExpr, info) || !cand.addValWrite(defT, defNonNil, defIsNil, accessExpr, info) {
			cand.disqualified = true
			return
		}
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpGetOrDefault, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr,
			accessExpr: accessExpr, call: mCall, keyArgIdx: 0, valArgIdx: 1,
			unwrapsKeyBox: unboxedK != stripParens(mCall.Args[0]),
			unwrapsValBox: unboxedDef != stripParens(mCall.Args[1]),
			valIsNilLit:   defIsNil,
		})

	case "Put_2", "Put", "PutIfAbsent":
		if len(mCall.Args) != 2 {
			cand.disqualified = true
			return
		}
		unboxedK, keyT, okK := unwrapMapKeyArg(mCall.Args[0], cand.hasIdentityKind, info)
		unboxedV, valWT, nonNilAlloc, isNilLit, okV := unwrapMapValArg(mCall.Args[1], info)
		if !okK || !okV || !cand.addKeyType(keyT, accessExpr, info) || !cand.addValWrite(valWT, nonNilAlloc, isNilLit, accessExpr, info) {
			cand.disqualified = true
			return
		}
		outerExpr := ast.Expr(mCall)
		discard := inExprStmt || isBlankAssignRhs(mCall, parents)
		if !discard {
			var readT types.Type
			var okR bool
			outerExpr, readT, okR = enclosingMapValUnboxExpr(mCall, parents, info)
			if !okR || !cand.addValReadType(readT, accessExpr, info) {
				cand.disqualified = true
				return
			}
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outerExpr.Pos()})
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpPut, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr,
			accessExpr: accessExpr, call: mCall, keyArgIdx: 0, valArgIdx: 1,
			unwrapsKeyBox: unboxedK != stripParens(mCall.Args[0]),
			unwrapsValBox: unboxedV != stripParens(mCall.Args[1]),
			valIsNilLit:   isNilLit,
			discardRet:    discard,
			isPutIfAbsent: methodSel.Sel.Name == "PutIfAbsent",
		})
		markMut(cand, mCall, parents)

	case "Remove_1", "Remove":
		if len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		unboxedK, keyT, okK := unwrapMapKeyArg(mCall.Args[0], cand.hasIdentityKind, info)
		if !okK || !cand.addKeyType(keyT, accessExpr, info) {
			cand.disqualified = true
			return
		}
		outerExpr := ast.Expr(mCall)
		discard := inExprStmt || isBlankAssignRhs(mCall, parents)
		if !discard {
			var readT types.Type
			var okR bool
			outerExpr, readT, okR = enclosingMapValUnboxExpr(mCall, parents, info)
			if !okR || !cand.addValReadType(readT, accessExpr, info) {
				cand.disqualified = true
				return
			}
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outerExpr.Pos()})
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpRemove, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr,
			accessExpr: accessExpr, call: mCall, keyArgIdx: 0,
			unwrapsKeyBox: unboxedK != stripParens(mCall.Args[0]),
			discardRet:    discard,
		})
		markMut(cand, mCall, parents)

	case "PutAll":
		if !inExprStmt || len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		unwrappedArg := unwrapMapAdaptersDown(mCall.Args[0], info)
		if srcCand := resolveMapVarAccess(unwrappedArg, varCands, info); srcCand != nil && (srcCand != cand || (cand.isField && !sameAccessPath(accessExpr, unwrappedArg, cand.v, info))) {
			*edges = append(*edges, [2]*mapCandidate{cand, srcCand})
			cand.ops = append(cand.ops, mapRewriteOp{
				kind: mapOpPutAllMap, cand: cand, otherCand: srcCand, file: file, spanNode: mCall, outerExpr: mCall,
				accessExpr: accessExpr, otherExpr: unwrappedArg,
			})
			markMut(cand, mCall, parents)
			return
		}
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpPutAllJavaMap, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, otherExpr: mCall.Args[0],
		})
		markMut(cand, mCall, parents)

	case "ComputeIfAbsent":
		if len(mCall.Args) != 2 {
			cand.disqualified = true
			return
		}
		unboxedK, keyT, okK := unwrapMapKeyArg(mCall.Args[0], cand.hasIdentityKind, info)
		if !okK || !cand.addKeyType(keyT, accessExpr, info) {
			cand.disqualified = true
			return
		}
		validateDirectComputeIfAbsent(mCall, accessExpr, 0, unboxedK != stripParens(mCall.Args[0]), mCall.Args[1], inExprStmt, file, cand, parents, info, varUseCounts, markMut)

	case "ForEach_1", "ForEach":
		if !inExprStmt || len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		validateMapForEach(mCall, accessExpr, file, cand, parents, info, varUseCounts)

	case "Equals":
		if inExprStmt || len(mCall.Args) != 1 {
			cand.disqualified = true
			return
		}
		otherInner := unwrapEqualsOtherArg(mCall.Args[0], info)
		otherCand := resolveMapVarAccess(otherInner, varCands, info)
		if otherCand == nil || (otherCand == cand && !cand.isField) {
			cand.disqualified = true
			return
		}
		*edges = append(*edges, [2]*mapCandidate{cand, otherCand})
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: mCall.Pos()})
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpEqualsMap, cand: cand, otherCand: otherCand, file: file, spanNode: mCall, outerExpr: mCall,
			accessExpr: accessExpr, otherExpr: otherInner,
		})

	case "HashCode":
		if inExprStmt || len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpHashCode, cand: cand, file: file, spanNode: mCall, outerExpr: mCall, accessExpr: accessExpr,
		})

	case "KeySet_0", "KeySet", "Values":
		if inExprStmt || len(mCall.Args) != 0 {
			cand.disqualified = true
			return
		}
		isKeys := methodSel.Sel.Name != "Values"
		validateMapCollectionView(mCall, isKeys, accessExpr, file, cand, ctx, parents, info, varUseCounts)

	case "EntrySet_0", "EntrySet":
		if inExprStmt || len(mCall.Args) != 0 || !isReadOnlyEntrySetStreamUse(mCall, parents, info) {
			cand.disqualified = true
			return
		}
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpEntrySetReadOnly, cand: cand, file: file, spanNode: mCall, outerExpr: mCall, accessExpr: accessExpr,
		})

	default:
		cand.disqualified = true
	}
}

func validateMapCollectionView(
	viewCall *ast.CallExpr,
	isKeys bool,
	accessExpr ast.Expr,
	file *ast.File,
	cand *mapCandidate,
	ctx *PassContext,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
) {
	adaptedView := unwrapCollectionAdapters(viewCall, parents, info)
	if !isKeys && j2gortHasExport(ctx, "OrderedMapGetNodeAny") {
		if iterSel, ok := parents[adaptedView].(*ast.SelectorExpr); ok && iterSel.X == adaptedView && iterSel.Sel.Name == "Iterator" {
			if iterCall, ok := parents[iterSel].(*ast.CallExpr); ok && iterCall.Fun == iterSel && len(iterCall.Args) == 0 && !iterCall.Ellipsis.IsValid() {
				adaptedIter := unwrapIteratorCast(iterCall, parents, info)
				if nextSel, ok := parents[adaptedIter].(*ast.SelectorExpr); ok && nextSel.X == adaptedIter && nextSel.Sel.Name == "Next" {
					if nextCall, ok := parents[nextSel].(*ast.CallExpr); ok && nextCall.Fun == nextSel && len(nextCall.Args) == 0 && !nextCall.Ellipsis.IsValid() {
						outerExpr, valT, okV := enclosingMapValUnboxExpr(nextCall, parents, info)
						if _, isES := parents[outerExpr].(*ast.ExprStmt); !isES && okV && cand.addValReadType(valT, accessExpr, info) {
							cand.ops = append(cand.ops, mapRewriteOp{
								kind: mapOpFirstValue, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr, accessExpr: accessExpr,
							})
							return
						}
					}
				}
			}
		}
	}
	if outerCall, ok := parents[adaptedView].(*ast.CallExpr); ok && len(outerCall.Args) == 1 && !outerCall.Ellipsis.IsValid() && outerCall.Args[0] == adaptedView {
		if fnName, _, _, isJ2 := isJ2gortCall(outerCall, info); isJ2 {
			switch fnName {
			case "NewIterCursor":
				cursorKind := "NewMapValCursor"
				if isKeys {
					cursorKind = "NewMapKeyCursor"
				}
				validateMapIterCursorLoop(cursorKind, outerCall, accessExpr, file, cand, parents, info, varUseCounts)
				return
			case "Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection",
				"Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_lang_Iterable":
				if isKeys {
					if _, isES := parents[outerCall].(*ast.ExprStmt); !isES {
						cand.ops = append(cand.ops, mapRewriteOp{
							kind: mapOpKeysToImmSet, cand: cand, file: file, spanNode: outerCall, outerExpr: outerCall, accessExpr: accessExpr,
						})
						return
					}
				}
			case "Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection",
				"Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable":
				if !isKeys {
					if _, isES := parents[outerCall].(*ast.ExprStmt); !isES {
						cand.ops = append(cand.ops, mapRewriteOp{
							kind: mapOpValsToImmList, cand: cand, file: file, spanNode: outerCall, outerExpr: outerCall, accessExpr: accessExpr,
						})
						return
					}
				}
			}
		}
	}
	// Live collection/keyset view is allowed on struct fields (where the enclosing struct is already heap-allocated)
	// and on single-assignment non-nullable local value maps when consumed by a non-escaping read-only operation.
	sel := "LiveValues"
	if isKeys {
		sel = "LiveKeySet"
	}
	if cand.isField {
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpLiveView, cand: cand, file: file, spanNode: viewCall, outerExpr: viewCall,
			accessExpr: accessExpr, methodsSel: sel,
		})
		return
	}
	if !cand.isParam && !cand.isCastAccessor && !cand.isPtrAlias &&
		j2gortOrderedMapHasMethod(ctx, sel) &&
		isNonEscapingReadOnlyLiveViewConsumer(adaptedView, parents, info) {
		cand.hasLocalLiveView = true
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpLiveView, cand: cand, file: file, spanNode: viewCall, outerExpr: viewCall,
			accessExpr: accessExpr, methodsSel: sel,
		})
		return
	}
	cand.disqualified = true
}

func validateStaticMapOp(
	fnName string,
	staticCall *ast.CallExpr,
	accessExpr ast.Expr,
	file *ast.File,
	cand *mapCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
	markMut func(*mapCandidate, ast.Node, map[ast.Node]ast.Node),
) bool {
	isNode := strings.HasSuffix(fnName, "Node")
	checkKey := func(arg ast.Expr) bool {
		if isNode {
			t := info.Types[arg].Type
			return isNamedUint32NodeType(t) && cand.addKeyType(t, accessExpr, info)
		}
		return hasExactOrDefaultStringType(arg, info) && cand.addKeyType(types.Typ[types.String], accessExpr, info)
	}
	_, inExprStmt := parents[staticCall].(*ast.ExprStmt)

	switch fnName {
	case "MapContainsKeyStr", "MapContainsKeyNode":
		if inExprStmt || len(staticCall.Args) != 2 || !checkKey(staticCall.Args[1]) {
			return false
		}
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpContainsKey, cand: cand, file: file, spanNode: staticCall, outerExpr: staticCall,
			accessExpr: accessExpr, call: staticCall, keyArgIdx: 1,
		})
		return true

	case "MapGetStr", "MapGetNode":
		if inExprStmt || len(staticCall.Args) != 2 || !checkKey(staticCall.Args[1]) {
			return false
		}
		outerExpr, valT, okV := enclosingMapValUnboxExpr(staticCall, parents, info)
		_, outerES := parents[outerExpr].(*ast.ExprStmt)
		if !okV || outerES || !cand.addValReadType(valT, accessExpr, info) {
			return false
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outerExpr.Pos()})
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpGet, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr,
			accessExpr: accessExpr, call: staticCall, keyArgIdx: 1,
		})
		return true

	case "MapGetOrDefaultStr", "MapGetOrDefaultNode":
		if inExprStmt || len(staticCall.Args) != 3 || !checkKey(staticCall.Args[1]) {
			return false
		}
		unboxedDef, defT, defNonNil, defIsNil, okD := unwrapMapValArg(staticCall.Args[2], info)
		outerExpr, valT, okV := enclosingMapValUnboxExpr(staticCall, parents, info)
		_, outerES := parents[outerExpr].(*ast.ExprStmt)
		if !okD || !okV || outerES || !cand.addValReadType(valT, accessExpr, info) || !cand.addValWrite(defT, defNonNil, defIsNil, accessExpr, info) {
			return false
		}
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpGetOrDefault, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr,
			accessExpr: accessExpr, call: staticCall, keyArgIdx: 1, valArgIdx: 2,
			unwrapsValBox: unboxedDef != stripParens(staticCall.Args[2]),
			valIsNilLit:   defIsNil,
		})
		return true

	case "MapPutStr", "MapPutNode":
		if len(staticCall.Args) != 3 || !checkKey(staticCall.Args[1]) {
			return false
		}
		unboxedV, valWT, nonNilAlloc, isNilLit, okV := unwrapMapValArg(staticCall.Args[2], info)
		if !okV || !cand.addValWrite(valWT, nonNilAlloc, isNilLit, accessExpr, info) {
			return false
		}
		outerExpr := ast.Expr(staticCall)
		discard := inExprStmt || isBlankAssignRhs(staticCall, parents)
		if !discard {
			var readT types.Type
			var okR bool
			outerExpr, readT, okR = enclosingMapValUnboxExpr(staticCall, parents, info)
			if !okR || !cand.addValReadType(readT, accessExpr, info) {
				return false
			}
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outerExpr.Pos()})
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpPut, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr,
			accessExpr: accessExpr, call: staticCall, keyArgIdx: 1, valArgIdx: 2,
			unwrapsValBox: unboxedV != stripParens(staticCall.Args[2]),
			valIsNilLit:   isNilLit,
			discardRet:    discard,
		})
		markMut(cand, staticCall, parents)
		return true

	case "MapRemoveStr", "MapRemoveNode":
		if len(staticCall.Args) != 2 || !checkKey(staticCall.Args[1]) {
			return false
		}
		outerExpr := ast.Expr(staticCall)
		discard := inExprStmt || isBlankAssignRhs(staticCall, parents)
		if !discard {
			var readT types.Type
			var okR bool
			outerExpr, readT, okR = enclosingMapValUnboxExpr(staticCall, parents, info)
			if !okR || !cand.addValReadType(readT, accessExpr, info) {
				return false
			}
		}
		cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outerExpr.Pos()})
		cand.ops = append(cand.ops, mapRewriteOp{
			kind: mapOpRemove, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr,
			accessExpr: accessExpr, call: staticCall, keyArgIdx: 1,
			discardRet: discard,
		})
		markMut(cand, staticCall, parents)
		return true

	case "MapComputeIfAbsentStr", "MapComputeIfAbsentNode":
		if len(staticCall.Args) != 3 || !checkKey(staticCall.Args[1]) {
			return false
		}
		validateDirectComputeIfAbsent(staticCall, accessExpr, 1, false, staticCall.Args[2], inExprStmt, file, cand, parents, info, varUseCounts, markMut)
		return !cand.disqualified
	}
	return false
}

func collectComputeIIFEMapTemps(
	files []*ast.File,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
	varUseCounts map[*types.Var]int,
	info *types.Info,
) map[*types.Var]bool {
	out := make(map[*types.Var]bool)
	for _, file := range files {
		parents := parentsByFile[file]
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) != 0 {
				return true
			}
			fl, ok := stripParens(call.Fun).(*ast.FuncLit)
			if !ok || fl.Body == nil || (len(fl.Body.List) != 3 && len(fl.Body.List) != 4) {
				return true
			}
			s0, ok0 := fl.Body.List[0].(*ast.AssignStmt)
			s1, ok1 := fl.Body.List[1].(*ast.AssignStmt)
			ifs, ok2 := fl.Body.List[2].(*ast.IfStmt)
			if !ok0 || !ok1 || !ok2 || s0.Tok != token.DEFINE || s1.Tok != token.DEFINE || len(s0.Lhs) != 1 || len(s1.Lhs) != 1 {
				return true
			}
			mId, _ := s0.Lhs[0].(*ast.Ident)
			mVar := varFromObj(info.Defs[mId])
			if mVar == nil || varUseCounts[mVar] != 2 {
				return true
			}
			initAs, ok := ifs.Init.(*ast.AssignStmt)
			if !ok || len(initAs.Rhs) != 1 {
				return true
			}
			tryCall, ok := stripParens(initAs.Rhs[0]).(*ast.CallExpr)
			if !ok || len(tryCall.Args) != 2 {
				return true
			}
			if fnName, _, _, isJ2 := isJ2gortCall(tryCall, info); isJ2 && strings.HasPrefix(fnName, "MapTryGetCompute") {
				if id0, ok := stripParens(tryCall.Args[0]).(*ast.Ident); ok && varFromObj(info.Uses[id0]) == mVar {
					_ = parents
					out[mVar] = true
				}
			}
			return true
		})
	}
	return out
}

func validateComputeIfAbsentIIFE(
	adaptedTop ast.Expr,
	accessExpr ast.Expr,
	file *ast.File,
	cand *mapCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
	markMut func(*mapCandidate, ast.Node, map[ast.Node]ast.Node),
) bool {
	s0, ok := parents[adaptedTop].(*ast.AssignStmt)
	if !ok || s0.Tok != token.DEFINE || len(s0.Lhs) != 1 || len(s0.Rhs) != 1 || s0.Rhs[0] != adaptedTop {
		return false
	}
	body, ok := parents[s0].(*ast.BlockStmt)
	if !ok || len(body.List) < 3 || body.List[0] != s0 {
		return false
	}
	fl, ok := parents[body].(*ast.FuncLit)
	if !ok || fl.Type.Params != nil && len(fl.Type.Params.List) > 0 {
		return false
	}
	flTop := unwrapParensUp(fl, parents)
	iifeCall, ok := parents[flTop].(*ast.CallExpr)
	if !ok || iifeCall.Fun != flTop || len(iifeCall.Args) != 0 {
		return false
	}

	mId, _ := s0.Lhs[0].(*ast.Ident)
	mVar := varFromObj(info.Defs[mId])
	s1, ok1 := body.List[1].(*ast.AssignStmt)
	ifs, ok2 := body.List[2].(*ast.IfStmt)
	if mVar == nil || varUseCounts[mVar] != 2 || !ok1 || !ok2 || s1.Tok != token.DEFINE || len(s1.Lhs) != 1 || len(s1.Rhs) != 1 {
		return false
	}
	kId, _ := s1.Lhs[0].(*ast.Ident)
	kVar := varFromObj(info.Defs[kId])
	if kVar == nil || varUseCounts[kVar] != 2 {
		return false
	}
	initAs, ok := ifs.Init.(*ast.AssignStmt)
	if !ok || initAs.Tok != token.DEFINE || len(initAs.Lhs) != 2 || len(initAs.Rhs) != 1 {
		return false
	}
	tryCall, ok := stripParens(initAs.Rhs[0]).(*ast.CallExpr)
	if !ok || len(tryCall.Args) != 2 {
		return false
	}
	tryFn, _, _, isJ2 := isJ2gortCall(tryCall, info)
	if !isJ2 || (tryFn != "MapTryGetComputeAny" && tryFn != "MapTryGetComputeStr" && tryFn != "MapTryGetComputeNode") {
		return false
	}
	if id0, ok := stripParens(tryCall.Args[0]).(*ast.Ident); !ok || varFromObj(info.Uses[id0]) != mVar {
		return false
	}
	if id1, ok := stripParens(tryCall.Args[1]).(*ast.Ident); !ok || varFromObj(info.Uses[id1]) != kVar {
		return false
	}

	rawKeyExpr := s1.Rhs[0]
	unboxedK, keyT, okK := unwrapMapKeyArg(rawKeyExpr, cand.hasIdentityKind, info)
	if !okK || !cand.addKeyType(keyT, accessExpr, info) {
		cand.disqualified = true
		return true
	}

	var slowCall *ast.CallExpr
	_, inExprStmt := parents[iifeCall].(*ast.ExprStmt)
	if len(body.List) == 4 && fl.Type.Results != nil && len(fl.Type.Results.List) == 1 {
		retStmt, ok := body.List[3].(*ast.ReturnStmt)
		if !ok || len(retStmt.Results) != 1 {
			cand.disqualified = true
			return true
		}
		retT := info.Types[iifeCall].Type
		if !cand.addValReadType(retT, accessExpr, info) {
			cand.disqualified = true
			return true
		}
		slowCall = unwrapCastAroundSlowCompute(retStmt.Results[0], info)
	} else if len(body.List) == 3 && (fl.Type.Results == nil || len(fl.Type.Results.List) == 0) && inExprStmt {
		if ifs.Body == nil || len(ifs.Body.List) != 1 {
			cand.disqualified = true
			return true
		}
		if es, ok := ifs.Body.List[0].(*ast.ExprStmt); ok {
			slowCall, _ = stripParens(es.X).(*ast.CallExpr)
		}
	}
	if slowCall == nil {
		cand.disqualified = true
		return true
	}

	fnExpr, ok := extractComputeFnFromSlowCall(slowCall, mVar, kVar, info)
	if !ok {
		cand.disqualified = true
		return true
	}
	spec, ok := buildComputeClosureSpec(fnExpr, accessExpr, file, cand, parents, info, varUseCounts)
	if !ok {
		cand.disqualified = true
		return true
	}

	cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: iifeCall.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
	cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: iifeCall.Pos()})
	cand.ops = append(cand.ops, mapRewriteOp{
		kind: mapOpComputeIfAbsent, cand: cand, file: file, spanNode: iifeCall, outerExpr: iifeCall,
		accessExpr: accessExpr, keyExpr: rawKeyExpr, unwrapsKeyBox: unboxedK != stripParens(rawKeyExpr),
		discardRet: inExprStmt, computeSpec: spec,
	})
	markMut(cand, iifeCall, parents)
	return true
}

func validateDirectComputeIfAbsent(
	call *ast.CallExpr,
	accessExpr ast.Expr,
	keyArgIdx int,
	unwrapsKey bool,
	fnExpr ast.Expr,
	inExprStmt bool,
	file *ast.File,
	cand *mapCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
	markMut func(*mapCandidate, ast.Node, map[ast.Node]ast.Node),
) {
	outerExpr := ast.Expr(call)
	discard := inExprStmt || isBlankAssignRhs(call, parents)
	if !discard {
		var readT types.Type
		var okR bool
		outerExpr, readT, okR = enclosingMapValUnboxExpr(call, parents, info)
		if !okR || !cand.addValReadType(readT, accessExpr, info) {
			cand.disqualified = true
			return
		}
	}
	spec, ok := buildComputeClosureSpec(fnExpr, accessExpr, file, cand, parents, info, varUseCounts)
	if !ok {
		cand.disqualified = true
		return
	}
	cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: outerExpr.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
	cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: outerExpr.Pos()})
	cand.ops = append(cand.ops, mapRewriteOp{
		kind: mapOpComputeIfAbsent, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr,
		accessExpr: accessExpr, call: call, keyArgIdx: keyArgIdx, unwrapsKeyBox: unwrapsKey,
		discardRet: discard, computeSpec: spec,
	})
	markMut(cand, call, parents)
}

func unwrapCastAroundSlowCompute(expr ast.Expr, info *types.Info) *ast.CallExpr {
	c, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || c.Ellipsis.IsValid() {
		return nil
	}
	if fnName, _, _, isJ2 := isJ2gortCall(c, info); isJ2 && len(c.Args) == 1 {
		switch fnName {
		case "CastPtr", "Cast", "UncheckedCast", "AsIface", "CastString", "CastNode":
			inner, _ := stripParens(c.Args[0]).(*ast.CallExpr)
			return inner
		}
	}
	return c
}

func extractComputeFnFromSlowCall(slowCall *ast.CallExpr, mVar, kVar *types.Var, info *types.Info) (ast.Expr, bool) {
	if slowCall == nil || slowCall.Ellipsis.IsValid() {
		return nil, false
	}
	if fnName, _, _, isJ2 := isJ2gortCall(slowCall, info); isJ2 && (fnName == "MapComputeIfAbsentStr" || fnName == "MapComputeIfAbsentNode") && len(slowCall.Args) == 3 {
		id0, ok0 := stripParens(slowCall.Args[0]).(*ast.Ident)
		id1, ok1 := stripParens(slowCall.Args[1]).(*ast.Ident)
		if ok0 && ok1 && varFromObj(info.Uses[id0]) == mVar && varFromObj(info.Uses[id1]) == kVar {
			return slowCall.Args[2], true
		}
		return nil, false
	}
	if sel, ok := stripParens(slowCall.Fun).(*ast.SelectorExpr); ok && sel.Sel.Name == "ComputeIfAbsent" && len(slowCall.Args) == 2 {
		id0, ok0 := stripParens(sel.X).(*ast.Ident)
		id1, ok1 := stripParens(slowCall.Args[0]).(*ast.Ident)
		if ok0 && ok1 && varFromObj(info.Uses[id0]) == mVar && varFromObj(info.Uses[id1]) == kVar {
			return slowCall.Args[1], true
		}
	}
	return nil, false
}

func buildComputeClosureSpec(
	fnExpr ast.Expr,
	accessExpr ast.Expr,
	file *ast.File,
	cand *mapCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
) (*computeClosureSpec, bool) {
	spec := &computeClosureSpec{fnExpr: fnExpr}
	lambdaCall, boundInit := matchNewLambdaCall(fnExpr, "NewLambda_Ext_java_util_function_Function", info)
	if lambdaCall == nil {
		return spec, true
	}
	fl, ok := stripParens(lambdaCall.Args[0]).(*ast.FuncLit)
	if !ok || fl.Body == nil || fl.Type.Params == nil || len(fl.Type.Params.List) != 1 || len(fl.Type.Params.List[0].Names) != 1 {
		return spec, true
	}
	pIdent := fl.Type.Params.List[0].Names[0]
	pVar := varFromObj(info.Defs[pIdent])
	if pIdent.Name != "_" && pVar != nil && varUseCounts[pVar] > 0 {
		var paramUnboxes []mapRewriteOp
		validParam := true
		ast.Inspect(fl.Body, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !validParam || !ok || varFromObj(info.Uses[id]) != pVar {
				return validParam
			}
			if isBlankAssignTarget(id, parents) {
				return true
			}
			outer, expectedK, okU := enclosingUnboxExpr(id, parents, info)
			_, isES := parents[outer].(*ast.ExprStmt)
			if !okU || isES || isEmptyInterfaceType(expectedK) || !cand.addKeyType(expectedK, accessExpr, info) {
				validParam = false
				return false
			}
			if outer != id {
				paramUnboxes = append(paramUnboxes, mapRewriteOp{
					kind: mapOpRawUnbox, cand: cand, file: file, spanNode: outer, outerExpr: outer, elemIdent: id,
				})
			}
			return true
		})
		if !validParam {
			return spec, true
		}
		cand.ops = append(cand.ops, paramUnboxes...)
		spec.hasParamVar = true
		spec.paramName = pIdent.Name
		spec.paramPos = pIdent.Pos()
	}

	retUnwraps := make(map[*ast.ReturnStmt]bool)
	hasRet := false
	validRet := true
	var walk func(ast.Node) bool
	walk = func(n ast.Node) bool {
		if !validRet || n == nil {
			return false
		}
		if nested, isLit := n.(*ast.FuncLit); isLit && nested != fl {
			return false
		}
		ret, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		hasRet = true
		if len(ret.Results) != 1 {
			validRet = false
			return false
		}
		unboxedV, valWT, nonNilAlloc, isNilLit, okV := unwrapMapValArg(ret.Results[0], info)
		if !okV || !cand.addValWrite(valWT, nonNilAlloc, isNilLit, accessExpr, info) {
			validRet = false
			return false
		}
		retUnwraps[ret] = unboxedV != stripParens(ret.Results[0])
		return false
	}
	ast.Inspect(fl.Body, walk)
	if !hasRet || !validRet {
		return spec, true
	}

	spec.boundInit = boundInit
	spec.fnLit = fl
	spec.retUnwraps = retUnwraps
	return spec, true
}

func validateMapForEach(
	mCall *ast.CallExpr,
	accessExpr ast.Expr,
	file *ast.File,
	cand *mapCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
) {
	argExpr := mCall.Args[0]
	spec := &forEachClosureSpec{argExpr: argExpr}
	if lambdaCall, boundInit := matchNewLambdaCall(argExpr, "NewLambda_Ext_java_util_function_BiConsumer", info); lambdaCall != nil {
		if fl, ok := stripParens(lambdaCall.Args[0]).(*ast.FuncLit); ok && fl.Body != nil && fl.Type.Params != nil {
			cand.loopBodies = append(cand.loopBodies, fl.Body)
			pIdents := collectFuncParamIdents(fl.Type.Params)
			if len(pIdents) == 2 {
				kId, vId := pIdents[0], pIdents[1]
				kVar, vVar := varFromObj(info.Defs[kId]), varFromObj(info.Defs[vId])
				var unboxOps []mapRewriteOp
				valid := true
				ast.Inspect(fl.Body, func(n ast.Node) bool {
					id, ok := n.(*ast.Ident)
					if !valid || !ok {
						return valid
					}
					u := varFromObj(info.Uses[id])
					if u != nil && u == kVar {
						if isBlankAssignTarget(id, parents) {
							return true
						}
						outer, kt, okU := enclosingUnboxExpr(id, parents, info)
						_, isES := parents[outer].(*ast.ExprStmt)
						if !okU || isES || !cand.addKeyType(kt, accessExpr, info) {
							valid = false
							return false
						}
						if outer != id {
							unboxOps = append(unboxOps, mapRewriteOp{kind: mapOpRawUnbox, cand: cand, file: file, spanNode: outer, outerExpr: outer, elemIdent: id})
						}
					} else if u != nil && u == vVar {
						if isBlankAssignTarget(id, parents) {
							return true
						}
						outer, vt, okU := enclosingMapValUnboxExpr(id, parents, info)
						_, isES := parents[outer].(*ast.ExprStmt)
						if !okU || isES || !cand.addValReadType(vt, accessExpr, info) {
							valid = false
							return false
						}
						if outer != id {
							unboxOps = append(unboxOps, mapRewriteOp{kind: mapOpRawUnbox, cand: cand, file: file, spanNode: outer, outerExpr: outer, elemIdent: id})
						}
					}
					return true
				})
				_ = varUseCounts
				if valid {
					cand.ops = append(cand.ops, unboxOps...)
					spec.boundInit = boundInit
					spec.fnLit = fl
					spec.kName = nonBlankParamName(kId.Name, "_k")
					spec.kPos = kId.Pos()
					spec.vName = nonBlankParamName(vId.Name, "_v")
					spec.vPos = vId.Pos()
					cand.spellSites = append(cand.spellSites, typeSpellSite{file: file, pos: mCall.Pos(), accessExpr: accessExpr, owner: cand.ownerStruct})
				}
			}
		}
	}
	cand.ops = append(cand.ops, mapRewriteOp{
		kind: mapOpForEach, cand: cand, file: file, spanNode: mCall, outerExpr: mCall,
		accessExpr: accessExpr, forEachSpec: spec,
	})
}

func validateMapIterCursorLoop(
	cursorFn string,
	cursorCall *ast.CallExpr,
	accessExpr ast.Expr,
	file *ast.File,
	cand *mapCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
) {
	isVals := cursorFn == "NewMapValCursor"
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

	addElemT := func(t types.Type) bool { return cand.addKeyType(t, accessExpr, info) }
	unboxElem := enclosingUnboxExpr
	if isVals {
		addElemT = func(t types.Type) bool { return cand.addValReadType(t, accessExpr, info) }
		unboxElem = enclosingMapValUnboxExpr
	}

	elemId := lhs0
	drop := 1
	rhs0 := stripParens(s0.Rhs[0])
	switch {
	case isCursorMethodCall(rhs0, curVar, "NextStr", info):
		if !addElemT(types.Typ[types.String]) {
			cand.disqualified = true
			return
		}
	case isCursorMethodCall(rhs0, curVar, "NextNode", info):
		if !addElemT(types.Typ[types.Uint32]) {
			cand.disqualified = true
			return
		}
	case isNodeConvOfCursorNextNode(rhs0, curVar, info):
		if !addElemT(info.Types[rhs0].Type) {
			cand.disqualified = true
			return
		}
	case isCursorMethodCall(rhs0, curVar, "Next", info):
		rawVar := varFromObj(info.Defs[lhs0])
		if rawVar == nil {
			cand.disqualified = true
			return
		}
		if len(forStmt.Body.List) >= 2 && varUseCounts[rawVar] == 1 {
			if unboxedId, elemT, ok := matchSingleMapDefineUnbox(forStmt.Body.List[1], rawVar, isVals, parents, info); ok && addElemT(elemT) {
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
				if isBlankAssignTarget(id, parents) {
					return true
				}
				outerExpr, expectedT, ok := unboxElem(id, parents, info)
				_, isES := parents[outerExpr].(*ast.ExprStmt)
				if !ok || isES || !addElemT(expectedT) {
					valid = false
					return false
				}
				if outerExpr != id {
					cand.ops = append(cand.ops, mapRewriteOp{kind: mapOpRawUnbox, cand: cand, file: file, spanNode: outerExpr, outerExpr: outerExpr, elemIdent: id})
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

	iterSel := "IterKeys"
	if isVals {
		iterSel = "IterValues"
	}
	cand.loopBodies = append(cand.loopBodies, forStmt.Body)
	cand.ops = append(cand.ops, mapRewriteOp{
		kind: mapOpIterCursor, cand: cand, file: file, spanNode: forStmt,
		stmt: initAs, forStmt: forStmt, outerExpr: cursorCall, accessExpr: accessExpr,
		methodsSel: iterSel, elemIdent: elemId, dropStmts: drop,
	})
}

func validateMapEntryCursorLoop(
	cursorCall *ast.CallExpr,
	accessExpr ast.Expr,
	file *ast.File,
	cand *mapCandidate,
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
	if curId == nil || curId.Name == "_" || curVar == nil || block == nil {
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
	if !ok || forStmt.Init != nil || forStmt.Post != nil || forStmt.Cond == nil || forStmt.Body == nil || !isCursorMethodCall(forStmt.Cond, curVar, "NextEntry", info) {
		cand.disqualified = true
		return
	}

	bodyUses := 0
	valid := true
	var entryOps []mapRewriteOp
	ast.Inspect(forStmt.Body, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if !valid || !ok || varFromObj(info.Uses[id]) != curVar {
			return valid
		}
		bodyUses++
		sel, ok := parents[id].(*ast.SelectorExpr)
		if !ok || sel.X != id {
			valid = false
			return false
		}
		mCall, ok := parents[sel].(*ast.CallExpr)
		if !ok || mCall.Fun != sel || mCall.Ellipsis.IsValid() {
			valid = false
			return false
		}
		_, inES := parents[mCall].(*ast.ExprStmt)

		switch sel.Sel.Name {
		case "KeyStr":
			if inES || len(mCall.Args) != 0 || !cand.addKeyType(types.Typ[types.String], accessExpr, info) {
				valid = false
				return false
			}
			entryOps = append(entryOps, mapRewriteOp{
				kind: mapOpEntryGetKey, cand: cand, file: file, spanNode: mCall, outerExpr: mCall, elemIdent: id,
			})
		case "KeyNode":
			if inES || len(mCall.Args) != 0 {
				valid = false
				return false
			}
			outer := ast.Expr(mCall)
			keyT := types.Type(types.Typ[types.Uint32])
			pTop := unwrapParensUp(mCall, parents)
			if conv, ok := parents[pTop].(*ast.CallExpr); ok && len(conv.Args) == 1 && conv.Args[0] == pTop && info.Types[conv.Fun].IsType() && isNamedUint32NodeType(info.Types[conv.Fun].Type) {
				outer = conv
				keyT = info.Types[conv.Fun].Type
			}
			if !cand.addKeyType(keyT, accessExpr, info) {
				valid = false
				return false
			}
			entryOps = append(entryOps, mapRewriteOp{
				kind: mapOpEntryGetKey, cand: cand, file: file, spanNode: outer, outerExpr: outer, elemIdent: id,
			})
		case "GetKey":
			if inES || len(mCall.Args) != 0 {
				valid = false
				return false
			}
			outer, keyT, okK := enclosingUnboxExpr(mCall, parents, info)
			_, outerES := parents[outer].(*ast.ExprStmt)
			if !okK || outerES || !cand.addKeyType(keyT, accessExpr, info) {
				valid = false
				return false
			}
			entryOps = append(entryOps, mapRewriteOp{
				kind: mapOpEntryGetKey, cand: cand, file: file, spanNode: outer, outerExpr: outer, elemIdent: id,
			})
		case "GetValue":
			if inES || len(mCall.Args) != 0 {
				valid = false
				return false
			}
			outer, valT, okV := enclosingMapValUnboxExpr(mCall, parents, info)
			_, outerES := parents[outer].(*ast.ExprStmt)
			if !okV || outerES || !cand.addValReadType(valT, accessExpr, info) {
				valid = false
				return false
			}
			entryOps = append(entryOps, mapRewriteOp{
				kind: mapOpEntryGetVal, cand: cand, file: file, spanNode: outer, outerExpr: outer, elemIdent: id,
			})
		case "SetValue":
			if len(mCall.Args) != 1 {
				valid = false
				return false
			}
			unboxedV, valWT, nonNilAlloc, isNilLit, okV := unwrapMapValArg(mCall.Args[0], info)
			if !okV || !cand.addValWrite(valWT, nonNilAlloc, isNilLit, accessExpr, info) {
				valid = false
				return false
			}
			outer := ast.Expr(mCall)
			discard := inES || isBlankAssignRhs(mCall, parents)
			if !discard {
				var readT types.Type
				var okR bool
				outer, readT, okR = enclosingMapValUnboxExpr(mCall, parents, info)
				if !okR || !cand.addValReadType(readT, accessExpr, info) {
					valid = false
					return false
				}
			}
			entryOps = append(entryOps, mapRewriteOp{
				kind: mapOpEntrySetVal, cand: cand, file: file, spanNode: outer, outerExpr: outer,
				call: mCall, elemIdent: id, valArgIdx: 0,
				unwrapsValBox: unboxedV != stripParens(mCall.Args[0]),
				valIsNilLit:   isNilLit,
			})
		default:
			valid = false
			return false
		}
		return true
	})

	if !valid || varUseCounts[curVar] != 1+bodyUses {
		cand.disqualified = true
		return
	}
	cand.loopBodies = append(cand.loopBodies, forStmt.Body)
	cand.ops = append(cand.ops, entryOps...)
	cand.ops = append(cand.ops, mapRewriteOp{
		kind: mapOpEntryCursor, cand: cand, file: file, spanNode: initAs,
		stmt: initAs, outerExpr: cursorCall, accessExpr: accessExpr,
	})
}

func anyMapLoopMutatesCandidate(cand *mapCandidate, peers []*mapCandidate, mutatingFuncs map[*types.Func]bool, info *types.Info) bool {
	refsSingle := func(c *mapCandidate, expr ast.Expr) bool {
		if c.v != nil && exprReferencesVar(expr, c.v, info) {
			return true
		}
		if c.castFn != nil {
			if call, ok := unwrapMapAdaptersDown(expr, info).(*ast.CallExpr); ok && len(call.Args) == 0 {
				if callee, _ := calleeFunc(call, info); callee != nil && funcOrigin(callee) == c.castFn {
					return true
				}
			}
		}
		return false
	}
	refsCand := func(expr ast.Expr) bool {
		if refsSingle(cand, expr) {
			return true
		}
		for _, p := range peers {
			if refsSingle(p, expr) {
				return true
			}
		}
		return false
	}
	lhsRefsCandOrPeer := func(lhs ast.Expr) bool {
		if (cand.v != nil && exprReferencesVar(lhs, cand.v, info)) ||
			(cand.unionFieldVar != nil && exprReferencesVar(lhs, cand.unionFieldVar, info)) {
			return true
		}
		for _, p := range peers {
			if (p.v != nil && exprReferencesVar(lhs, p.v, info)) ||
				(p.unionFieldVar != nil && exprReferencesVar(lhs, p.unionFieldVar, info)) {
				return true
			}
		}
		return false
	}
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
					case "MapPutStr", "MapPutNode", "MapRemoveStr", "MapRemoveNode", "MapComputeIfAbsentStr", "MapComputeIfAbsentNode":
						if len(x.Args) >= 1 && refsCand(x.Args[0]) {
							mutated = true
						}
					}
				}
				if sel, ok := stripParens(x.Fun).(*ast.SelectorExpr); ok {
					switch sel.Sel.Name {
					case "Put_2", "Put", "PutIfAbsent", "PutAll", "Remove_1", "Remove", "Clear", "ComputeIfAbsent", "ComputeIfPresent", "Compute", "Merge":
						if refsCand(sel.X) {
							mutated = true
						}
					}
				}
			case *ast.AssignStmt:
				for _, lhs := range x.Lhs {
					if lhsRefsCandOrPeer(lhs) {
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

func applyMapRewrites(ctx *PassContext, eligible []*mapCandidate, parentsByFile map[*ast.File]map[ast.Node]ast.Node) int {
	var exprOps, stmtOps []mapRewriteOp
	for _, cand := range eligible {
		for _, op := range cand.ops {
			switch op.kind {
			case mapOpDeclField, mapOpDeclCastAccessor, mapOpDeclVarSpec, mapOpDeclShortVar, mapOpAssignInit, mapOpIterCursor, mapOpEntryCursor, mapOpMethodRefGet:
				stmtOps = append(stmtOps, op)
			default:
				exprOps = append(exprOps, op)
			}
		}
	}
	spanOrder := func(a, b mapRewriteOp) int {
		if c := cmp.Compare(a.spanNode.End()-a.spanNode.Pos(), b.spanNode.End()-b.spanNode.Pos()); c != 0 {
			return c
		}
		return cmp.Compare(a.spanNode.Pos(), b.spanNode.Pos())
	}
	slices.SortFunc(exprOps, spanOrder)
	slices.SortFunc(stmtOps, spanOrder)

	count := 0
	for _, op := range exprOps {
		if applySingleMapRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}
	for _, op := range stmtOps {
		if applySingleMapRewrite(ctx, op, parentsByFile[op.file]) {
			ctx.MarkFileModified(op.file)
			count++
		}
	}
	return count
}

func applySingleMapRewrite(ctx *PassContext, op mapRewriteOp, parents map[ast.Node]ast.Node) bool {
	j2Sel := func(name string, pos token.Pos) *ast.SelectorExpr {
		q, _ := j2gortQualifierForFile(op.file, ctx)
		return &ast.SelectorExpr{X: &ast.Ident{NamePos: pos, Name: q}, Sel: &ast.Ident{NamePos: pos, Name: name}}
	}
	siteKeyT := instantiateOwnerTypeAtExpr(op.cand.keyType, op.cand.ownerStruct, op.accessExpr, ctx.Info)
	siteValT := instantiateOwnerTypeAtExpr(op.cand.valType, op.cand.ownerStruct, op.accessExpr, ctx.Info)
	makeKeyType := func() ast.Expr {
		kt, _ := typeExprForFile(siteKeyT, op.file, ctx.Pkg, ctx.Importer)
		return kt
	}
	makeValType := func() ast.Expr {
		vt, _ := typeExprForFile(siteValT, op.file, ctx.Pkg, ctx.Importer)
		return vt
	}
	typeBaseName := func() string {
		switch op.cand.flavor {
		case mapFlavorCustom:
			return "CustomOrderedMap"
		case mapFlavorIdentity:
			return "IdentityOrderedMap"
		case mapFlavorPlain:
			return "PlainMap"
		default:
			return "OrderedMap"
		}
	}
	makeValueMapType := func(pos token.Pos) ast.Expr {
		return &ast.IndexListExpr{
			X:       j2Sel(typeBaseName(), pos),
			Indices: []ast.Expr{makeKeyType(), makeValType()},
		}
	}
	accessIsPtr := op.cand.isPtr() || op.viaPtrGetter
	makeMapType := func(pos token.Pos) ast.Expr {
		mt := makeValueMapType(pos)
		if op.cand.isPtr() {
			return &ast.StarExpr{Star: pos, X: mt}
		}
		return mt
	}
	mapAddrTarget := func(expr ast.Expr, pos token.Pos) ast.Expr {
		if accessIsPtr {
			return expr
		}
		return &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: wrapPrimaryTarget(expr)}
	}
	otherAddrTarget := func(expr ast.Expr, oc *mapCandidate, pos token.Pos) ast.Expr {
		if oc != nil && oc.isPtr() {
			return expr
		}
		return &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: wrapPrimaryTarget(expr)}
	}
	isStrVal := types.Identical(siteValT, types.Typ[types.String])

	makeSingleInitRhs := func(initK mapInitKind, argExpr ast.Expr, oc *mapCandidate, rhs ast.Expr, inHeader bool) ast.Expr {
		pos := rhs.Pos()
		base := typeBaseName()
		switch initK {
		case mapInitNil:
			if op.cand.isPtr() {
				return &ast.Ident{NamePos: pos, Name: "nil"}
			}
			lit := ast.Expr(&ast.CompositeLit{Type: makeValueMapType(pos)})
			if inHeader {
				lit = &ast.ParenExpr{X: lit}
			}
			return lit
		case mapInitEmpty:
			if op.cand.flavor == mapFlavorPlain {
				return makePlainMapInitExpr(ctx, j2Sel, makeValueMapType, makeKeyType, makeValType, op.cand.isPtr(), false, nil, pos, inHeader)
			}
			if op.cand.isPtr() {
				lit := ast.Expr(&ast.UnaryExpr{
					OpPos: pos,
					Op:    token.AND,
					X:     &ast.CompositeLit{Type: makeValueMapType(pos)},
				})
				if inHeader {
					lit = &ast.ParenExpr{X: lit}
				}
				return lit
			}
			return &ast.CallExpr{Fun: &ast.IndexListExpr{X: j2Sel("New"+base, pos), Indices: []ast.Expr{makeKeyType(), makeValType()}}}
		case mapInitCap:
			if op.cand.flavor == mapFlavorPlain {
				return makePlainMapInitExpr(ctx, j2Sel, makeValueMapType, makeKeyType, makeValType, op.cand.isPtr(), true, argExpr, pos, inHeader)
			}
			fnName := "New" + base + "Cap"
			if op.cand.isPtr() {
				fnName = "New" + base + "PtrCap"
			}
			return &ast.CallExpr{
				Fun:  &ast.IndexListExpr{X: j2Sel(fnName, pos), Indices: []ast.Expr{makeKeyType(), makeValType()}},
				Args: []ast.Expr{argExpr},
			}
		case mapInitCopyMap:
			if op.cand.isPtr() {
				return &ast.CallExpr{Fun: j2Sel(base+"PtrFrom", pos), Args: []ast.Expr{otherAddrTarget(argExpr, oc, pos)}}
			}
			if oc != nil && oc.isPtr() {
				return &ast.CallExpr{Fun: j2Sel(base+"FromPtr", pos), Args: []ast.Expr{argExpr}}
			}
			return &ast.CallExpr{Fun: j2Sel(base+"From", pos), Args: []ast.Expr{argExpr}}
		case mapInitCopyJavaMap:
			fnName := base + "FromJavaMap"
			if op.cand.isPtr() {
				fnName = base + "PtrFromJavaMap"
			}
			if op.cand.flavor == mapFlavorOrdered && isNamedUint32NodeType(siteKeyT) {
				fnName = "OrderedMapFromNodeKeyJavaMap"
				if op.cand.isPtr() {
					fnName = "OrderedMapPtrFromNodeKeyJavaMap"
				}
			}
			return &ast.CallExpr{
				Fun:  &ast.IndexListExpr{X: j2Sel(fnName, pos), Indices: []ast.Expr{makeKeyType(), makeValType()}},
				Args: []ast.Expr{argExpr},
			}
		case mapInitPtrAlias, mapInitSharedPtr:
			if oc != nil && oc.isPtr() {
				return argExpr
			}
			return &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: wrapPrimaryTarget(argExpr)}
		case mapInitCastAccessorCall:
			return argExpr
		}
		return rhs
	}
	makeInitRhs := func(rhs ast.Expr, inHeader bool) ast.Expr {
		if op.initKind == mapInitIIFE && op.iifeSpec != nil {
			fl := op.iifeSpec.fnLit
			fl.Type.Results.List[0].Type = makeMapType(fl.Type.Results.Pos())
			for _, br := range op.iifeSpec.branches {
				br.ret.Results[0] = makeSingleInitRhs(br.initKind, br.otherExpr, br.otherCand, br.ret.Results[0], false)
			}
			return rhs
		}
		return makeSingleInitRhs(op.initKind, op.otherExpr, op.otherCand, rhs, inHeader)
	}

	switch op.kind {
	case mapOpDeclField, mapOpDeclParam:
		fld := op.spanNode.(*ast.Field)
		fld.Type = makeMapType(fld.Pos())
		return true

	case mapOpDeclFieldGetter:
		fn := op.spanNode.(*ast.FuncDecl)
		pos := fn.Type.Results.Pos()
		fn.Type.Results.List[0].Type = &ast.StarExpr{Star: pos, X: makeValueMapType(pos)}
		ret := fn.Body.List[len(fn.Body.List)-1].(*ast.ReturnStmt)
		checkedRecv := &ast.CallExpr{
			Fun:  j2Sel("CheckNotNullPtr", pos),
			Args: []ast.Expr{&ast.Ident{NamePos: pos, Name: op.elemIdent.Name}},
		}
		fieldAccess := ast.Expr(&ast.SelectorExpr{
			X:   checkedRecv,
			Sel: &ast.Ident{NamePos: pos, Name: op.methodsSel},
		})
		if !op.cand.isPtr() {
			fieldAccess = &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: fieldAccess}
		}
		ret.Results[0] = fieldAccess
		return true

	case mapOpPassParamArg:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := mapAddrTarget(op.accessExpr, pos)
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpDeclCastAccessor:
		fn := op.spanNode.(*ast.FuncDecl)
		pos := fn.Type.Results.Pos()
		fn.Type.Results.List[0].Type = makeMapType(pos)
		ret := fn.Body.List[0].(*ast.ReturnStmt)
		origCall := stripParens(ret.Results[0]).(*ast.CallExpr)
		ret.Results[0] = &ast.CallExpr{
			Fun:  &ast.IndexExpr{X: j2Sel("CastPtr", pos), Index: makeValueMapType(pos)},
			Args: []ast.Expr{origCall.Args[0]},
		}
		return true

	case mapOpAssignUnionField:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		repl := &ast.CallExpr{Fun: j2Sel("PtrToAny", pos), Args: []ast.Expr{op.accessExpr}}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpMethodRefGet:
		spec := op.methodRefSpec
		pos := spec.boundAs.Pos()
		if op.cand.isPtr() {
			spec.boundAs.Rhs[0] = &ast.CallExpr{
				Fun:  j2Sel("CheckNotNullPtr", pos),
				Args: []ast.Expr{op.accessExpr},
			}
		} else {
			if sel, ok := stripParens(op.accessExpr).(*ast.SelectorExpr); ok {
				checkedRecv := &ast.CallExpr{Fun: j2Sel("CheckNotNullPtr", pos), Args: []ast.Expr{sel.X}}
				spec.boundAs.Rhs[0] = &ast.UnaryExpr{
					OpPos: pos,
					Op:    token.AND,
					X:     &ast.SelectorExpr{X: checkedRecv, Sel: &ast.Ident{NamePos: sel.Sel.Pos(), Name: sel.Sel.Name}},
				}
			} else {
				spec.boundAs.Rhs[0] = &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: wrapPrimaryTarget(op.accessExpr)}
			}
		}
		helper := "OrderedMapGetAny"
		switch {
		case op.cand.flavor == mapFlavorOrdered && isNamedUint32NodeType(siteKeyT) && j2gortHasExport(ctx, "OrderedMapGetNodeAny"):
			helper = "OrderedMapGetNodeAny"
		case op.cand.flavor == mapFlavorCustom && j2gortHasExport(ctx, "CustomOrderedMapGetAny"):
			helper = "CustomOrderedMapGetAny"
		case op.cand.flavor == mapFlavorIdentity && j2gortHasExport(ctx, "IdentityOrderedMapGetAny"):
			helper = "IdentityOrderedMapGetAny"
		}
		getParent := parents[spec.getCall]
		pArg := spec.paramArg
		if pId, ok := stripParens(pArg).(*ast.Ident); ok {
			pArg = &ast.Ident{NamePos: pId.Pos(), Name: pId.Name}
		}
		replGet := &ast.CallExpr{
			Fun:    j2Sel(helper, spec.getCall.Pos()),
			Lparen: spec.getCall.Lparen,
			Args:   []ast.Expr{&ast.Ident{NamePos: spec.getCall.Lparen, Name: spec.boundId.Name}, pArg},
			Rparen: spec.getCall.Rparen,
		}
		return getParent != nil && replaceExprInParent(getParent, spec.getCall, replGet, parents)

	case mapOpFirstValue:
		parent := parents[op.outerExpr]
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "FirstValue"}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpDeclVarSpec:
		spec := op.spanNode.(*ast.ValueSpec)
		spec.Type = makeMapType(spec.Pos())
		if len(spec.Values) == 1 {
			if op.initKind == mapInitNil {
				spec.Values = nil
			} else {
				spec.Values[0] = makeInitRhs(spec.Values[0], false)
			}
		}
		return true

	case mapOpDeclShortVar, mapOpAssignInit:
		if op.initKind == mapInitNil && op.cand.isPtr() {
			return false
		}
		as := op.spanNode.(*ast.AssignStmt)
		as.Rhs[0] = makeInitRhs(as.Rhs[0], inControlHeader(as, parents))
		return true

	case mapOpNilCheck:
		if accessIsPtr {
			return false
		}
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		nilSel := "IsNil"
		if op.cand.isLazyValueMap {
			nilSel = "IsEmpty"
		}
		var repl ast.Expr = &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: nilSel}},
		}
		if op.isNegated {
			repl = &ast.UnaryExpr{OpPos: pos, Op: token.NOT, X: repl}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpFormatSink:
		if op.unwrapsValBox && accessIsPtr {
			return false
		}
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		var repl ast.Expr
		if accessIsPtr {
			repl = &ast.CallExpr{Fun: j2Sel("PtrToAny", pos), Args: []ast.Expr{op.accessExpr}}
		} else {
			repl = &ast.UnaryExpr{OpPos: pos, Op: token.AND, X: wrapPrimaryTarget(op.accessExpr)}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpRawUnbox:
		parent := parents[op.outerExpr]
		return parent != nil && replaceExprInParent(parent, op.outerExpr, &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name}, parents)

	case mapOpMethodDirect, mapOpLiveView:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		target := op.accessExpr
		if op.methodsSel == "String" && op.cand.isPtr() {
			target = &ast.CallExpr{Fun: j2Sel("CheckNotNullPtr", pos), Args: []ast.Expr{op.accessExpr}}
		}
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(target), Sel: &ast.Ident{NamePos: pos, Name: op.methodsSel}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpEntrySetReadOnly:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		toJavaCall := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: "ToJavaLinkedHashMap"}},
		}
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: toJavaCall, Sel: &ast.Ident{NamePos: pos, Name: "EntrySet"}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpContainsKey:
		parent := parents[op.outerExpr]
		keyArg := unwrapBoxArgExpr(op.call.Args[op.keyArgIdx], op.unwrapsKeyBox)
		repl := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "ContainsKey"}},
			Args: []ast.Expr{keyArg},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpGetAny:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		helper := "OrderedMapGetAny"
		switch {
		case op.cand.flavor == mapFlavorOrdered && isNamedUint32NodeType(siteKeyT) && j2gortHasExport(ctx, "OrderedMapGetNodeAny"):
			helper = "OrderedMapGetNodeAny"
		case op.cand.flavor == mapFlavorCustom && j2gortHasExport(ctx, "CustomOrderedMapGetAny"):
			helper = "CustomOrderedMapGetAny"
		case op.cand.flavor == mapFlavorIdentity && j2gortHasExport(ctx, "IdentityOrderedMapGetAny"):
			helper = "IdentityOrderedMapGetAny"
		}
		repl := &ast.CallExpr{
			Fun:  j2Sel(helper, pos),
			Args: []ast.Expr{mapAddrTarget(op.accessExpr, pos), op.call.Args[op.keyArgIdx]},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpGet:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		keyArg := unwrapBoxArgExpr(op.call.Args[op.keyArgIdx], op.unwrapsKeyBox)
		var repl ast.Expr
		if isStrVal {
			fnName := typeBaseName() + "GetStr"
			if accessIsPtr {
				fnName = typeBaseName() + "PtrGetStr"
			}
			repl = &ast.CallExpr{
				Fun:  j2Sel(fnName, pos),
				Args: []ast.Expr{op.accessExpr, keyArg},
			}
		} else {
			repl = &ast.CallExpr{
				Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: "Get"}},
				Args: []ast.Expr{keyArg},
			}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpGetOrDefault:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		keyArg := unwrapBoxArgExpr(op.call.Args[op.keyArgIdx], op.unwrapsKeyBox)
		defArg := unwrapMapValExprAtApply(op.call.Args[op.valArgIdx], op.unwrapsValBox, op.valIsNilLit)
		repl := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: "GetOrDefault"}},
			Args: []ast.Expr{keyArg, defArg},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpPut:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		keyArg := unwrapBoxArgExpr(op.call.Args[op.keyArgIdx], op.unwrapsKeyBox)
		valArg := unwrapMapValExprAtApply(op.call.Args[op.valArgIdx], op.unwrapsValBox, op.valIsNilLit)
		var repl ast.Expr
		if isStrVal && (!op.discardRet || op.isPutIfAbsent) {
			helperSuffix := "PutStr"
			if op.isPutIfAbsent {
				helperSuffix = "PutIfAbsentStr"
			}
			repl = &ast.CallExpr{
				Fun:  j2Sel(typeBaseName()+helperSuffix, pos),
				Args: []ast.Expr{mapAddrTarget(op.accessExpr, pos), keyArg, valArg},
			}
		} else {
			sel := "Put"
			if op.isPutIfAbsent {
				sel = "PutIfAbsent"
			} else if op.cand.flavor == mapFlavorPlain && op.discardRet {
				if _, isES := parent.(*ast.ExprStmt); isES && j2gortTypeHasMethod(ctx, "PlainMap", "PutVoid") {
					sel = "PutVoid"
				}
			}
			repl = &ast.CallExpr{
				Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: sel}},
				Args: []ast.Expr{keyArg, valArg},
			}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpRemove:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		keyArg := unwrapBoxArgExpr(op.call.Args[op.keyArgIdx], op.unwrapsKeyBox)
		var repl ast.Expr
		if isStrVal && !op.discardRet {
			repl = &ast.CallExpr{
				Fun:  j2Sel(typeBaseName()+"RemoveStr", pos),
				Args: []ast.Expr{mapAddrTarget(op.accessExpr, pos), keyArg},
			}
		} else {
			repl = &ast.CallExpr{
				Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: "Remove"}},
				Args: []ast.Expr{keyArg},
			}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpPutAllMap:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		selName := "PutAllMap"
		if op.otherCand != nil && op.otherCand.isPtr() {
			selName = "PutAllMapPtr"
		}
		repl := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: selName}},
			Args: []ast.Expr{op.otherExpr},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpPutAllJavaMap:
		parent := parents[op.outerExpr]
		repl := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "PutAllJavaMap"}},
			Args: []ast.Expr{op.otherExpr},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpComputeIfAbsent:
		parent := parents[op.outerExpr]
		rawK := op.keyExpr
		if rawK == nil {
			rawK = op.call.Args[op.keyArgIdx]
		}
		keyArg := unwrapBoxArgExpr(rawK, op.unwrapsKeyBox)
		clearNodePositions(keyArg)
		nullV := makeMapNullValExpr(siteValT, makeValType, j2Sel, token.NoPos)
		addr := mapAddrTarget(op.accessExpr, token.NoPos)
		clearNodePositions(addr)
		var repl ast.Expr
		if op.computeSpec != nil && op.computeSpec.fnLit != nil {
			fl := op.computeSpec.fnLit
			for ret, unwraps := range op.computeSpec.retUnwraps { // order-independent: mutates disjoint ReturnStmt AST nodes in-place
				ret.Results[0] = unwrapMapValExprAtApply(ret.Results[0], unwraps, false)
			}
			pName := "_"
			if op.computeSpec.hasParamVar {
				pName = op.computeSpec.paramName
			}
			typedLit := &ast.FuncLit{
				Type: &ast.FuncType{
					Params:  &ast.FieldList{List: []*ast.Field{{Names: []*ast.Ident{{Name: pName}}, Type: makeKeyType()}}},
					Results: &ast.FieldList{List: []*ast.Field{{Type: makeValType()}}},
				},
				Body: fl.Body,
			}
			computeFnName := typeBaseName() + "ComputeIfAbsent"
			computeArgs := []ast.Expr{addr, keyArg, nullV, typedLit}
			if op.cand.flavor == mapFlavorPlain {
				if helperName, takesNull, okH := resolvePlainMapComputeHelper(findJ2gortTypesPkg(ctx), true, isStrVal); okH {
					computeFnName = helperName
					if !takesNull {
						computeArgs = []ast.Expr{addr, keyArg, typedLit}
					}
				}
			}
			computeCall := &ast.CallExpr{
				Fun:  j2Sel(computeFnName, token.NoPos),
				Args: computeArgs,
			}
			if op.computeSpec.boundInit != nil {
				clearNodePositions(op.computeSpec.boundInit)
				if op.discardRet {
					repl = &ast.CallExpr{
						Fun: &ast.FuncLit{
							Type: &ast.FuncType{Params: &ast.FieldList{}},
							Body: &ast.BlockStmt{List: []ast.Stmt{op.computeSpec.boundInit, &ast.ExprStmt{X: computeCall}}},
						},
					}
				} else {
					repl = &ast.CallExpr{
						Fun: &ast.FuncLit{
							Type: &ast.FuncType{
								Params:  &ast.FieldList{},
								Results: &ast.FieldList{List: []*ast.Field{{Type: makeValType()}}},
							},
							Body: &ast.BlockStmt{List: []ast.Stmt{op.computeSpec.boundInit, &ast.ReturnStmt{Results: []ast.Expr{computeCall}}}},
						},
					}
				}
			} else {
				repl = computeCall
			}
		} else {
			clearNodePositions(op.computeSpec.fnExpr)
			computeFnName := typeBaseName() + "ComputeIfAbsentFn"
			computeArgs := []ast.Expr{addr, keyArg, nullV, op.computeSpec.fnExpr}
			if op.cand.flavor == mapFlavorPlain {
				if helperName, takesNull, okH := resolvePlainMapComputeHelper(findJ2gortTypesPkg(ctx), false, isStrVal); okH {
					computeFnName = helperName
					if !takesNull {
						computeArgs = []ast.Expr{addr, keyArg, op.computeSpec.fnExpr}
					}
				}
			}
			repl = &ast.CallExpr{
				Fun:  j2Sel(computeFnName, token.NoPos),
				Args: computeArgs,
			}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpForEach:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		if op.forEachSpec != nil && op.forEachSpec.fnLit != nil {
			typedLit := &ast.FuncLit{
				Type: &ast.FuncType{
					Params: &ast.FieldList{List: []*ast.Field{
						{Names: []*ast.Ident{{NamePos: op.forEachSpec.kPos, Name: op.forEachSpec.kName}}, Type: makeKeyType()},
						{Names: []*ast.Ident{{NamePos: op.forEachSpec.vPos, Name: op.forEachSpec.vName}}, Type: makeValType()},
					}},
				},
				Body: op.forEachSpec.fnLit.Body,
			}
			feCall := &ast.CallExpr{
				Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: "ForEach"}},
				Args: []ast.Expr{typedLit},
			}
			if op.forEachSpec.boundInit != nil {
				if es, isES := parent.(*ast.ExprStmt); isES && canReplaceStmtIn(parents[es]) {
					blk := &ast.BlockStmt{List: []ast.Stmt{op.forEachSpec.boundInit, &ast.ExprStmt{X: feCall}}}
					return replaceStmtInParent(parents[es], es, []ast.Stmt{blk}, parents)
				}
				iife := &ast.CallExpr{
					Fun: &ast.FuncLit{
						Type: &ast.FuncType{Params: &ast.FieldList{}},
						Body: &ast.BlockStmt{List: []ast.Stmt{op.forEachSpec.boundInit, &ast.ExprStmt{X: feCall}}},
					},
				}
				return parent != nil && replaceExprInParent(parent, op.outerExpr, iife, parents)
			}
			return parent != nil && replaceExprInParent(parent, op.outerExpr, feCall, parents)
		}
		repl := &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: "ForEachBiConsumer"}},
			Args: []ast.Expr{op.forEachSpec.argExpr},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpEqualsMap:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		eitherPtr := accessIsPtr || (op.otherCand != nil && op.otherCand.isPtr())
		lhs := op.accessExpr
		rhs := op.otherExpr
		if eitherPtr {
			lhs = mapAddrTarget(op.accessExpr, pos)
			rhs = otherAddrTarget(op.otherExpr, op.otherCand, pos)
		}
		identLit := &ast.Ident{NamePos: pos, Name: "false"}
		if op.cand.isIdentityMap {
			identLit.Name = "true"
		}
		var repl ast.Expr
		if isComparableOrJavaIface(siteValT) {
			switch op.cand.flavor {
			case mapFlavorCustom:
				fnName := "CustomOrderedMapEquals"
				if eitherPtr {
					fnName = "CustomOrderedMapPtrEquals"
				}
				repl = &ast.CallExpr{Fun: j2Sel(fnName, pos), Args: []ast.Expr{lhs, rhs}}
			case mapFlavorIdentity:
				fnName := "IdentityOrderedMapEquals"
				if eitherPtr {
					fnName = "IdentityOrderedMapPtrEquals"
				}
				repl = &ast.CallExpr{Fun: j2Sel(fnName, pos), Args: []ast.Expr{lhs, rhs}}
			default:
				fnName := "OrderedMapEquals"
				if eitherPtr {
					fnName = "OrderedMapPtrEquals"
				}
				repl = &ast.CallExpr{Fun: j2Sel(fnName, pos), Args: []ast.Expr{lhs, rhs, identLit}}
			}
		} else {
			selName := "EqualsMap"
			if op.otherCand != nil && op.otherCand.isPtr() {
				selName = "EqualsMapPtr"
			}
			args := []ast.Expr{op.otherExpr}
			if op.cand.flavor == mapFlavorOrdered {
				args = append(args, identLit)
			}
			repl = &ast.CallExpr{
				Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: selName}},
				Args: args,
			}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpHashCode:
		parent := parents[op.outerExpr]
		pos := op.outerExpr.Pos()
		var repl ast.Expr
		if op.cand.flavor == mapFlavorCustom || op.cand.flavor == mapFlavorIdentity {
			repl = &ast.CallExpr{
				Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: "HashCode"}},
			}
		} else {
			identName := "false"
			if op.cand.isIdentityMap {
				identName = "true"
			}
			repl = &ast.CallExpr{
				Fun:  &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: pos, Name: "HashCodeMap"}},
				Args: []ast.Expr{&ast.Ident{NamePos: pos, Name: identName}},
			}
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpImmMapCopy:
		parent := parents[op.outerExpr]
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "ToImmutableMap"}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpKeysToImmSet:
		parent := parents[op.outerExpr]
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "KeysToImmutableSet"}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpValsToImmList:
		parent := parents[op.outerExpr]
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "ValuesToImmutableList"}},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpIterCursor:
		initAs := op.stmt.(*ast.AssignStmt)
		curId := initAs.Lhs[0].(*ast.Ident)
		initAs.Rhs[0] = &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: op.methodsSel}},
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

	case mapOpEntryCursor:
		initAs := op.stmt.(*ast.AssignStmt)
		initAs.Rhs[0] = &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: wrapPrimaryTarget(op.accessExpr), Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "IterEntries"}},
		}
		return true

	case mapOpEntryGetKey:
		parent := parents[op.outerExpr]
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name},
				Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "Key"},
			},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpEntryGetVal:
		parent := parents[op.outerExpr]
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name},
				Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "Value"},
			},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)

	case mapOpEntrySetVal:
		parent := parents[op.outerExpr]
		valArg := unwrapMapValExprAtApply(op.call.Args[op.valArgIdx], op.unwrapsValBox, op.valIsNilLit)
		repl := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{NamePos: op.elemIdent.Pos(), Name: op.elemIdent.Name},
				Sel: &ast.Ident{NamePos: op.outerExpr.Pos(), Name: "SetValue"},
			},
			Args: []ast.Expr{valArg},
		}
		return parent != nil && replaceExprInParent(parent, op.outerExpr, repl, parents)
	}
	return false
}

func classifyCandidateMapType(t types.Type) (isCand, isIdent, isHash bool) {
	if t == nil {
		return false, false, false
	}
	ut := types.Unalias(t)
	if ptr, ok := ut.(*types.Pointer); ok {
		ut = types.Unalias(ptr.Elem())
	}
	named, ok := ut.(*types.Named)
	if !ok || named.Obj() == nil || !isJ2gortPkg(named.Obj().Pkg()) {
		return false, false, false
	}
	switch named.Obj().Name() {
	case "Ext_java_util_LinkedHashMap", "Ext_java_util_HashMap", "Ext_java_util_HashMapBase":
		return true, false, true
	case "Ext_java_util_IdentityHashMap":
		return true, true, false
	case "Ext_java_util_Map", "Ext_java_util_SequencedMap":
		return true, false, false
	}
	return false, false, false
}

func classifyMapInitExpr(expr ast.Expr, info *types.Info) (kind mapInitKind, argExpr ast.Expr, isIdent, isHash, ok bool) {
	expr = unwrapMapAdaptersDown(expr, info)
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return 0, nil, false, false, false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 {
		return 0, nil, false, false, false
	}
	switch fnName {
	case "NewExt_java_util_LinkedHashMap", "NewExt_java_util_HashMap",
		"Ext_com_google_common_collect_Maps_NewLinkedHashMap_0", "Ext_com_google_common_collect_Maps_NewHashMap_0":
		if len(call.Args) == 0 {
			return mapInitEmpty, nil, false, true, true
		}
	case "NewExt_java_util_IdentityHashMap", "Ext_com_google_common_collect_Maps_NewIdentityHashMap":
		if len(call.Args) == 0 {
			return mapInitEmpty, nil, true, false, true
		}
	case "NewExt_java_util_LinkedHashMap_1_int32", "NewExt_java_util_HashMap_1_int32",
		"Ext_com_google_common_collect_Maps_NewLinkedHashMapWithExpectedSize",
		"Ext_com_google_common_collect_Maps_NewHashMapWithExpectedSize":
		if len(call.Args) == 1 && isIntegerExpr(call.Args[0], info) {
			return mapInitCap, call.Args[0], false, true, true
		}
	case "NewExt_java_util_IdentityHashMap_1_int32":
		if len(call.Args) == 1 && isIntegerExpr(call.Args[0], info) {
			return mapInitCap, call.Args[0], true, false, true
		}
	case "NewExt_java_util_LinkedHashMap_1_Ext_java_util_Map", "NewExt_java_util_HashMap_1_Ext_java_util_Map",
		"Ext_com_google_common_collect_Maps_NewLinkedHashMap_1_Ext_java_util_Map",
		"Ext_com_google_common_collect_Maps_NewHashMap_1_Ext_java_util_Map":
		if len(call.Args) == 1 {
			return mapInitCopyJavaMap, call.Args[0], false, true, true
		}
	}
	return 0, nil, false, false, false
}

func classifyMapIIFEInit(expr ast.Expr, info *types.Info) (spec *mapIIFESpec, hasNil, isIdent, isHash, ok bool) {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 0 || call.Ellipsis.IsValid() {
		return nil, false, false, false, false
	}
	fl, ok := stripParens(call.Fun).(*ast.FuncLit)
	if !ok || fl.Body == nil || (fl.Type.Params != nil && len(fl.Type.Params.List) > 0) || fl.Type.Results == nil || len(fl.Type.Results.List) != 1 {
		return nil, false, false, false, false
	}
	spec = &mapIIFESpec{fnLit: fl}
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
			spec.branches = append(spec.branches, mapIIFEBranch{ret: ret, initKind: mapInitNil})
			return false
		}
		initK, argExpr, idK, hsK, initOK := classifyMapInitExpr(res, info)
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
		spec.branches = append(spec.branches, mapIIFEBranch{ret: ret, initKind: initK, otherExpr: argExpr})
		return false
	})
	if !valid || len(spec.branches) == 0 {
		return nil, false, false, false, false
	}
	return spec, hasNil, isIdent, isHash, true
}

func unwrapMapKeyArg(arg ast.Expr, isIdentityMap bool, info *types.Info) (ast.Expr, types.Type, bool) {
	if isIdentityMap {
		e := stripParens(arg)
		if call, ok := e.(*ast.CallExpr); ok && len(call.Args) == 1 && !call.Ellipsis.IsValid() {
			if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
				inner := stripParens(call.Args[0])
				tv, hasTV := info.Types[inner]
				switch fnName {
				case "BoxString":
					return inner, types.Typ[types.String], hasTV && hasExactOrDefaultStringType(inner, info)
				case "BoxNode":
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
				case "ToAny":
					if hasTV && !tv.IsNil() && tv.Type != nil {
						dt := types.Default(tv.Type)
						return inner, dt, !isInvalidOrUntyped(dt) && !isEmptyInterfaceType(dt)
					}
					return nil, nil, false
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
	return unwrapSetOrMapElemArg(arg, info)
}

func unwrapSetOrMapElemArg(arg ast.Expr, info *types.Info) (ast.Expr, types.Type, bool) {
	if unboxed, elemT, ok := unwrapBoxedAddArg(arg, info); ok {
		return unboxed, elemT, true
	}
	e := stripParens(arg)
	if call, ok := e.(*ast.CallExpr); ok && len(call.Args) == 1 && !call.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && fnName == "ToAny" {
			inner := stripParens(call.Args[0])
			if tv, hasTV := info.Types[inner]; hasTV && !tv.IsNil() && tv.Type != nil {
				dt := types.Default(tv.Type)
				_, isTP := types.Unalias(dt).(*types.TypeParam)
				if !isInvalidOrUntyped(dt) && (isNonEmptyInterfaceType(dt) || isTP) {
					return inner, dt, true
				}
			}
		}
	}
	return nil, nil, false
}

func unwrapMapValArg(arg ast.Expr, info *types.Info) (unboxed ast.Expr, valT types.Type, nonNilAlloc, isNilLit, ok bool) {
	e := stripParens(arg)
	if isNilIdent(e, info) {
		return e, nil, false, true, true
	}
	if call, isCall := e.(*ast.CallExpr); isCall && len(call.Args) == 1 && !call.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			inner := stripParens(call.Args[0])
			tv, hasTV := info.Types[inner]
			switch fnName {
			case "BoxString":
				if hasTV && hasExactOrDefaultStringType(inner, info) {
					return inner, types.Typ[types.String], false, false, true
				}
				return nil, nil, false, false, false
			case "BoxNode":
				if hasTV && !tv.IsNil() && tv.Type != nil && !isInvalidOrUntyped(tv.Type) {
					return inner, tv.Type, false, false, true
				}
				return nil, nil, false, false, false
			case "PtrToAny":
				if hasTV && !tv.IsNil() && tv.Type != nil {
					_, isPtr := types.Unalias(tv.Type).(*types.Pointer)
					return inner, tv.Type, isGuaranteedNonNilAlloc(inner, info), false, isPtr
				}
				return nil, nil, false, false, false
			case "SliceToAny":
				if hasTV && !tv.IsNil() && tv.Type != nil {
					_, isSl := types.Unalias(tv.Type).(*types.Slice)
					return inner, tv.Type, isGuaranteedNonNilAlloc(inner, info), false, isSl
				}
				return nil, nil, false, false, false
			case "ToAny", "BoxNullableVal":
				if hasTV && !tv.IsNil() && tv.Type != nil {
					dt := types.Default(tv.Type)
					if !isInvalidOrUntyped(dt) && !isEmptyInterfaceType(dt) {
						return inner, dt, isGuaranteedNonNilAlloc(inner, info), false, true
					}
				}
				return nil, nil, false, false, false
			}
		}
	}
	tv, hasTV := info.Types[e]
	if !hasTV || tv.IsType() || tv.IsBuiltin() || tv.IsNil() || tv.Type == nil {
		return nil, nil, false, false, false
	}
	dt := types.Default(tv.Type)
	if isInvalidOrUntyped(dt) || isEmptyInterfaceType(dt) {
		return nil, nil, false, false, false
	}
	return e, dt, isGuaranteedNonNilAlloc(e, info), false, true
}

func unwrapMapValExprAtApply(rawArg ast.Expr, unwrapsBox, isNilLit bool) ast.Expr {
	if isNilLit {
		return &ast.Ident{NamePos: rawArg.Pos(), Name: "nil"}
	}
	return unwrapBoxArgExpr(rawArg, unwrapsBox)
}

func enclosingMapValUnboxExpr(rawExpr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) (ast.Expr, types.Type, bool) {
	cur := unwrapParensUp(rawExpr, parents)
	call, ok := parents[cur].(*ast.CallExpr)
	if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() || call.Args[0] != cur {
		return rawExpr, types.Universe.Lookup("any").Type(), true
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 {
		return rawExpr, types.Universe.Lookup("any").Type(), true
	}
	switch fnName {
	case "CastPtr", "Cast", "UncheckedCast", "AsIface", "CastString", "CastNode", "CastVal", "BoxByteString":
		t := info.Types[call].Type
		return call, t, !isInvalidOrUntyped(t)
	}
	return rawExpr, types.Universe.Lookup("any").Type(), true
}

func matchSingleMapDefineUnbox(stmt ast.Stmt, rawVar *types.Var, isVals bool, parents map[ast.Node]ast.Node, info *types.Info) (*ast.Ident, types.Type, bool) {
	if !isVals {
		return matchSingleDefineUnbox(stmt, rawVar, parents, info)
	}
	as, ok := stmt.(*ast.AssignStmt)
	if !ok || as.Tok != token.DEFINE || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
		return nil, nil, false
	}
	lhs, ok := as.Lhs[0].(*ast.Ident)
	if !ok || lhs.Name == "_" || info.Defs[lhs] == nil {
		return nil, nil, false
	}
	var found *ast.Ident
	count := 0
	ast.Inspect(as.Rhs[0], func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && varFromObj(info.Uses[id]) == rawVar {
			found = id
			count++
		}
		return true
	})
	if count != 1 || found == nil {
		return nil, nil, false
	}
	outer, elemT, ok := enclosingMapValUnboxExpr(found, parents, info)
	if !ok || stripParens(as.Rhs[0]) != stripParens(outer) {
		return nil, nil, false
	}
	return lhs, elemT, true
}

func makeMapNullValExpr(valT types.Type, makeValType func() ast.Expr, j2Sel func(string, token.Pos) *ast.SelectorExpr, pos token.Pos) ast.Expr {
	ut := types.Unalias(valT)
	if _, isTP := ut.(*types.TypeParam); isTP {
		return &ast.StarExpr{Star: pos, X: &ast.CallExpr{Fun: &ast.Ident{NamePos: pos, Name: "new"}, Args: []ast.Expr{makeValType()}}}
	}
	if b, ok := ut.Underlying().(*types.Basic); ok {
		switch {
		case (b.Info() & types.IsString) != 0:
			return j2Sel("NullString", pos)
		case (b.Info() & types.IsBoolean) != 0:
			return &ast.Ident{NamePos: pos, Name: "false"}
		case (b.Info() & types.IsNumeric) != 0:
			if _, isNamed := ut.(*types.Named); isNamed {
				return &ast.CallExpr{Fun: makeValType(), Args: []ast.Expr{&ast.BasicLit{ValuePos: pos, Kind: token.INT, Value: "0"}}}
			}
			return &ast.BasicLit{ValuePos: pos, Kind: token.INT, Value: "0"}
		}
	}
	return &ast.Ident{NamePos: pos, Name: "nil"}
}

func matchNewLambdaCall(expr ast.Expr, wantFn string, info *types.Info) (*ast.CallExpr, *ast.AssignStmt) {
	e := stripParens(expr)
	if call, ok := e.(*ast.CallExpr); ok && !call.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && fnName == wantFn && len(call.Args) == 1 {
			return call, nil
		}
		// Method-reference IIFE: `func() Ext_java_util_function_* { bound := <expr>; return NewLambda_*(...) }()`
		if len(call.Args) == 0 {
			if fl, ok := stripParens(call.Fun).(*ast.FuncLit); ok && fl.Body != nil && len(fl.Body.List) == 2 {
				bAs, ok0 := fl.Body.List[0].(*ast.AssignStmt)
				ret, ok1 := fl.Body.List[1].(*ast.ReturnStmt)
				if ok0 && ok1 && bAs.Tok == token.DEFINE && len(bAs.Lhs) == 1 && len(bAs.Rhs) == 1 && len(ret.Results) == 1 {
					if inner, ok := stripParens(ret.Results[0]).(*ast.CallExpr); ok && len(inner.Args) == 1 && !inner.Ellipsis.IsValid() {
						if fnName, _, _, isJ2 := isJ2gortCall(inner, info); isJ2 && fnName == wantFn {
							return inner, bAs
						}
					}
				}
			}
		}
	}
	return nil, nil
}

func collectFuncParamIdents(params *ast.FieldList) []*ast.Ident {
	var out []*ast.Ident
	for _, f := range params.List {
		out = append(out, f.Names...)
	}
	return out
}

func nonBlankParamName(name, fallback string) string {
	if name == "" || name == "_" {
		return fallback
	}
	return name
}

func unwrapParensUp(expr ast.Expr, parents map[ast.Node]ast.Node) ast.Expr {
	for {
		pe, ok := parents[expr].(*ast.ParenExpr)
		if !ok || pe.X != expr {
			return expr
		}
		expr = pe
	}
}

func unwrapMapAdaptersUp(expr ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) ast.Expr {
	for {
		if pe, ok := parents[expr].(*ast.ParenExpr); ok && pe.X == expr {
			expr = pe
			continue
		}
		call, ok := parents[expr].(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() || call.Args[0] != expr {
			return expr
		}
		if _, isTI := isToIfaceCall(call, info); isTI && isMapIface(info.Types[call].Type) {
			expr = call
			continue
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "CheckNotNull", "CheckNotNullPtr", "CheckNotNullIface":
				expr = call
				continue
			case "AsIface", "Cast", "UncheckedCast":
				if isMapIface(info.Types[call].Type) {
					expr = call
					continue
				}
			}
		}
		return expr
	}
}

func unwrapMapAdaptersDown(expr ast.Expr, info *types.Info) ast.Expr {
	cur := stripParens(expr)
	for {
		call, ok := cur.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
			return cur
		}
		if _, isTI := isToIfaceCall(call, info); isTI && isMapIface(info.Types[call].Type) {
			cur = stripParens(call.Args[0])
			continue
		}
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "CheckNotNull", "CheckNotNullPtr", "CheckNotNullIface":
				cur = stripParens(call.Args[0])
				continue
			case "AsIface", "Cast", "UncheckedCast":
				if isMapIface(info.Types[call].Type) {
					cur = stripParens(call.Args[0])
					continue
				}
			}
		}
		return cur
	}
}

func unwrapEqualsOtherArg(expr ast.Expr, info *types.Info) ast.Expr {
	cur := unwrapMapAdaptersDown(expr, info)
	if call, ok := cur.(*ast.CallExpr); ok && len(call.Args) == 1 && !call.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && (fnName == "PtrToAny" || fnName == "ToAny") {
			return unwrapMapAdaptersDown(call.Args[0], info)
		}
	}
	return cur
}

func isMapIface(t types.Type) bool {
	named, ok := types.Unalias(t).(*types.Named)
	if !ok || named.Obj() == nil || !isJ2gortPkg(named.Obj().Pkg()) {
		return false
	}
	switch named.Obj().Name() {
	case "Ext_java_util_Map", "Ext_java_util_SequencedMap", "Ext_java_util_HashMap":
		return true
	}
	return false
}

func resolveMapVarAccess(expr ast.Expr, varCands map[*types.Var]*mapCandidate, info *types.Info) *mapCandidate {
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

func isMapInitTargetCandidate(initCall *ast.CallExpr, varCands map[*types.Var]*mapCandidate, parents map[ast.Node]ast.Node, info *types.Info) bool {
	cur := unwrapParensUp(initCall, parents)
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
			switch lhs := stripParens(p.Lhs[0]).(type) {
			case *ast.Ident:
				return varCands[varFromObj(info.Uses[lhs])] != nil
			case *ast.SelectorExpr:
				return varCands[varFromObj(info.Uses[lhs.Sel])] != nil
			}
		}
	}
	return false
}

func isBlankAssignRhs(expr ast.Expr, parents map[ast.Node]ast.Node) bool {
	cur := unwrapParensUp(expr, parents)
	as, ok := parents[cur].(*ast.AssignStmt)
	if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 || as.Rhs[0] != cur {
		return false
	}
	lhsId, ok := stripParens(as.Lhs[0]).(*ast.Ident)
	return ok && lhsId.Name == "_"
}

func isBlankAssignTarget(id *ast.Ident, parents map[ast.Node]ast.Node) bool {
	as, ok := parents[id].(*ast.AssignStmt)
	if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 || as.Rhs[0] != id {
		return false
	}
	lhsBlank, ok := stripParens(as.Lhs[0]).(*ast.Ident)
	return ok && lhsBlank.Name == "_"
}

func clearNodePositions(root ast.Node) {
	if root == nil {
		return
	}
	ast.Inspect(root, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			x.NamePos = token.NoPos
		case *ast.BasicLit:
			x.ValuePos = token.NoPos
		case *ast.CallExpr:
			x.Lparen = token.NoPos
			x.Rparen = token.NoPos
		case *ast.UnaryExpr:
			x.OpPos = token.NoPos
		case *ast.StarExpr:
			x.Star = token.NoPos
		case *ast.ParenExpr:
			x.Lparen = token.NoPos
			x.Rparen = token.NoPos
		case *ast.AssignStmt:
			x.TokPos = token.NoPos
		}
		return true
	})
}

func isInInitFieldsOrCtor(node ast.Node, owner *types.Named, parents map[ast.Node]ast.Node, info *types.Info) bool {
	if isInInitFieldsMethod(node, owner, parents, info) {
		return true
	}
	if owner == nil || owner.Obj() == nil {
		return false
	}
	for cur := node; cur != nil; cur = parents[cur] {
		fn, ok := cur.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Recv == nil && (strings.HasPrefix(fn.Name.Name, "Init"+owner.Obj().Name()) ||
			strings.HasPrefix(fn.Name.Name, "Init"+strings.TrimSuffix(owner.Obj().Name(), "Base"))) {
			return true
		}
		return false
	}
	return false
}

func isMethodRefIIFEBoundAssign(span ast.Node, parents map[ast.Node]ast.Node, info *types.Info) bool {
	as, ok := span.(*ast.AssignStmt)
	if !ok || as.Tok != token.DEFINE || parents == nil {
		return false
	}
	blk, ok := parents[as].(*ast.BlockStmt)
	if !ok || len(blk.List) != 2 || blk.List[0] != as {
		return false
	}
	ret, ok := blk.List[1].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return false
	}
	call, ok := stripParens(ret.Results[0]).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	return isJ2 && strings.HasPrefix(fnName, "NewLambda_")
}

func validateMethodRefGetIIFE(
	adaptedTop ast.Expr,
	accessExpr ast.Expr,
	file *ast.File,
	cand *mapCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
	varUseCounts map[*types.Var]int,
) bool {
	boundAs, ok := parents[adaptedTop].(*ast.AssignStmt)
	if !ok || boundAs.Tok != token.DEFINE || len(boundAs.Lhs) != 1 || len(boundAs.Rhs) != 1 || boundAs.Rhs[0] != adaptedTop {
		return false
	}
	boundId, ok := boundAs.Lhs[0].(*ast.Ident)
	if !ok || boundId.Name == "_" {
		return false
	}
	boundVar := varFromObj(info.Defs[boundId])
	if boundVar == nil || varUseCounts[boundVar] != 1 {
		return false
	}
	iifeBody, ok := parents[boundAs].(*ast.BlockStmt)
	if !ok || len(iifeBody.List) != 2 || iifeBody.List[0] != boundAs {
		return false
	}
	iifeLit, ok := parents[iifeBody].(*ast.FuncLit)
	if !ok || (iifeLit.Type.Params != nil && len(iifeLit.Type.Params.List) != 0) || iifeLit.Type.Results == nil || len(iifeLit.Type.Results.List) != 1 {
		return false
	}
	iifeTop := unwrapParensUp(iifeLit, parents)
	iifeCall, ok := parents[iifeTop].(*ast.CallExpr)
	if !ok || iifeCall.Fun != iifeTop || len(iifeCall.Args) != 0 || iifeCall.Ellipsis.IsValid() {
		return false
	}
	if _, isES := parents[iifeCall].(*ast.ExprStmt); isES {
		return false
	}
	ret1, ok := iifeBody.List[1].(*ast.ReturnStmt)
	if !ok || len(ret1.Results) != 1 {
		return false
	}
	lamCall, ok := stripParens(ret1.Results[0]).(*ast.CallExpr)
	if !ok || len(lamCall.Args) != 1 || lamCall.Ellipsis.IsValid() {
		return false
	}
	fnName, _, _, isJ2 := isJ2gortCall(lamCall, info)
	if !isJ2 || fnName != "NewLambda_Ext_java_util_function_Function" {
		return false
	}
	innerLit, ok := stripParens(lamCall.Args[0]).(*ast.FuncLit)
	if !ok || innerLit.Body == nil || len(innerLit.Body.List) != 1 || innerLit.Type.Params == nil || len(innerLit.Type.Params.List) != 1 {
		return false
	}
	pFld := innerLit.Type.Params.List[0]
	if len(pFld.Names) != 1 || pFld.Names[0].Name == "_" {
		return false
	}
	pVar := varFromObj(info.Defs[pFld.Names[0]])
	if pVar == nil || !isEmptyInterfaceType(pVar.Type()) {
		return false
	}
	innerRet, ok := innerLit.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(innerRet.Results) != 1 {
		return false
	}
	getCall, ok := stripParens(innerRet.Results[0]).(*ast.CallExpr)
	if !ok || len(getCall.Args) != 1 || getCall.Ellipsis.IsValid() {
		return false
	}
	getSel, ok := stripParens(getCall.Fun).(*ast.SelectorExpr)
	if !ok || getSel.Sel.Name != "Get" {
		return false
	}
	bUse, ok := stripParens(getSel.X).(*ast.Ident)
	if !ok || varFromObj(info.Uses[bUse]) != boundVar {
		return false
	}
	pUse, ok := stripParens(getCall.Args[0]).(*ast.Ident)
	if !ok || varFromObj(info.Uses[pUse]) != pVar {
		return false
	}
	if cand.isField {
		sel, ok := stripParens(accessExpr).(*ast.SelectorExpr)
		if !ok || !isPureFieldSelector(sel, info) {
			return false
		}
		if cand.hasNonCtorAssign {
			cand.forceNullable = true
		}
	} else {
		cand.forceNullable = true
	}
	cand.j2gortSites = append(cand.j2gortSites, typeSpellSite{file: file, pos: boundAs.Pos()})
	cand.ops = append(cand.ops, mapRewriteOp{
		kind: mapOpMethodRefGet, cand: cand, file: file, spanNode: boundAs, accessExpr: accessExpr,
		methodRefSpec: &methodRefGetSpec{
			boundAs:  boundAs,
			boundId:  boundId,
			getCall:  getCall,
			paramArg: getCall.Args[0],
		},
	})
	return true
}

func matchUnionFieldCastAccessor(fn *ast.FuncDecl, info *types.Info) (*types.Named, *types.Var, *types.Func, bool, bool, bool) {
	if fn.Type.TypeParams != nil && len(fn.Type.TypeParams.List) > 0 {
		return nil, nil, nil, false, false, false
	}
	if fn.Recv == nil || len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 {
		return nil, nil, nil, false, false, false
	}
	if (fn.Type.Params != nil && len(fn.Type.Params.List) != 0) || fn.Type.Results == nil || len(fn.Type.Results.List) != 1 || len(fn.Type.Results.List[0].Names) > 0 {
		return nil, nil, nil, false, false, false
	}
	fnObj, ok := info.Defs[fn.Name].(*types.Func)
	if !ok || fnObj == nil || fnObj.Exported() {
		return nil, nil, nil, false, false, false
	}
	sig, ok := fnObj.Type().(*types.Signature)
	if !ok || sig.Recv() == nil || sig.RecvTypeParams() != nil {
		return nil, nil, nil, false, false, false
	}
	recvVar := varFromObj(info.Defs[fn.Recv.List[0].Names[0]])
	if recvVar == nil {
		return nil, nil, nil, false, false, false
	}
	owner := ptrNamedOrigin(recvVar.Type())
	if owner == nil || (owner.TypeParams() != nil && owner.TypeParams().Len() > 0) {
		return nil, nil, nil, false, false, false
	}
	isCand, declIdent, declHash := classifyCandidateMapType(sig.Results().At(0).Type())
	if !isCand || len(fn.Body.List) != 1 {
		return nil, nil, nil, false, false, false
	}
	ret, ok := fn.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return nil, nil, nil, false, false, false
	}
	castCall, ok := stripParens(ret.Results[0]).(*ast.CallExpr)
	if !ok || len(castCall.Args) != 1 || castCall.Ellipsis.IsValid() {
		return nil, nil, nil, false, false, false
	}
	fnName, _, _, isJ2 := isJ2gortCall(castCall, info)
	if !isJ2 || (fnName != "AsIface" && fnName != "Cast" && fnName != "UncheckedCast" && fnName != "CastPtr") {
		return nil, nil, nil, false, false, false
	}
	sel, ok := stripParens(castCall.Args[0]).(*ast.SelectorExpr)
	if !ok || !isIndirectStructFieldSel(sel, info) {
		return nil, nil, nil, false, false, false
	}
	rId, ok := stripParens(sel.X).(*ast.Ident)
	if !ok || varFromObj(info.Uses[rId]) != recvVar {
		return nil, nil, nil, false, false, false
	}
	fVar := varFromObj(info.Uses[sel.Sel])
	if fVar == nil || !fVar.IsField() || fVar.Exported() || fVar.Embedded() || !isEmptyInterfaceType(fVar.Type()) {
		return nil, nil, nil, false, false, false
	}
	return owner, fVar, funcOrigin(fnObj), declIdent, declHash, true
}

func validateUnionFieldUse(
	id *ast.Ident,
	cc *mapCandidate,
	varCands map[*types.Var]*mapCandidate,
	canTransform bool,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) {
	if !canTransform {
		cc.disqualified = true
		return
	}
	// Skip the field read inside `cc.castFnDecl.Body` itself.
	if cc.castFnDecl != nil && cc.castFnDecl.Body != nil && id.Pos() >= cc.castFnDecl.Body.Pos() && id.Pos() <= cc.castFnDecl.Body.End() {
		return
	}
	if kv, isKV := parents[id].(*ast.KeyValueExpr); isKV && kv.Key == id {
		if isNilIdent(kv.Value, info) || isNonMapConcretePtrAnyExpr(kv.Value, info) {
			return
		}
		cc.disqualified = true
		return
	}
	sel, ok := parents[id].(*ast.SelectorExpr)
	if !ok || sel.Sel != id || !isIndirectStructFieldSel(sel, info) {
		cc.disqualified = true
		return
	}
	pTop := unwrapParensUp(sel, parents)
	if as, ok := parents[pTop].(*ast.AssignStmt); ok {
		for _, lhs := range as.Lhs {
			if stripParens(lhs) == sel {
				if as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 || !isPureFieldSelector(sel, info) {
					cc.disqualified = true
					return
				}
				rhs := as.Rhs[0]
				if isNilIdent(rhs, info) || isNonMapConcretePtrAnyExpr(rhs, info) {
					return
				}
				if srcCand := resolveMapVarAccess(unwrapMapAdaptersDown(rhs, info), varCands, info); srcCand != nil && !srcCand.isField && !srcCand.disqualified {
					return
				}
				cc.disqualified = true
				return
			}
		}
	}

	// Allowed read patterns on the `any` union field outside the cast accessor:
	// 1. `sel == nil` / `sel != nil`
	if bin, ok := parents[pTop].(*ast.BinaryExpr); ok && (bin.Op == token.EQL || bin.Op == token.NEQ) {
		other := bin.Y
		if bin.Y == pTop {
			other = bin.X
		}
		if isNilIdent(other, info) {
			return
		}
	}
	// 2. Formatting sinks (`CheckStateFmt(..., sel)`, `StrVal(sel)`, etc.).
	if _, _, ok := matchFormattingSetSink(pTop, parents, info); ok {
		return
	}
	// 3. `CheckNotNull(sel)` (as statement), `InstanceOf[*OtherStruct](sel)`, `TryCast[*OtherStruct](sel)`,
	//    `CastPtr[OtherStruct](sel)`, `RefEquals(sel, otherStructPtr)`.
	if call, ok := parents[pTop].(*ast.CallExpr); ok && !call.Ellipsis.IsValid() && call.Fun != pTop {
		if fnName, _, explicitTArg, isJ2 := isJ2gortCall(call, info); isJ2 {
			switch fnName {
			case "CheckNotNull":
				if len(call.Args) == 1 && call.Args[0] == pTop {
					if _, isES := parents[call].(*ast.ExprStmt); isES {
						return
					}
				}
			case "InstanceOf", "TryCast":
				if len(call.Args) == 1 && call.Args[0] == pTop && explicitTArg != nil && isNonMapStructPtrType(info.Types[explicitTArg].Type) {
					return
				}
			case "CastPtr":
				if len(call.Args) == 1 && call.Args[0] == pTop && explicitTArg != nil && isNonMapStructValueType(info.Types[explicitTArg].Type) {
					return
				}
			case "RefEquals":
				if len(call.Args) == 2 {
					var other ast.Expr
					if call.Args[0] == pTop {
						other = call.Args[1]
					} else if call.Args[1] == pTop {
						other = call.Args[0]
					}
					if other != nil && isNonMapStructPtrType(info.Types[other].Type) {
						return
					}
				}
			}
		}
	}
	cc.disqualified = true
}

func isNonMapStructValueType(t types.Type) bool {
	if t == nil {
		return false
	}
	named, ok := types.Unalias(t).(*types.Named)
	if !ok || named.Obj() == nil || isJ2gortPkg(named.Obj().Pkg()) {
		return false
	}
	_, isSt := named.Underlying().(*types.Struct)
	return isSt
}

func isNonMapStructPtrType(t types.Type) bool {
	ptr, ok := types.Unalias(t).(*types.Pointer)
	return ok && isNonMapStructValueType(ptr.Elem())
}

func isNonMapConcretePtrAnyExpr(expr ast.Expr, info *types.Info) bool {
	e := stripParens(expr)
	if call, ok := e.(*ast.CallExpr); ok && len(call.Args) == 1 && !call.Ellipsis.IsValid() {
		if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 && (fnName == "PtrToAny" || fnName == "ToAny" || fnName == "CheckNotNull" || fnName == "CheckNotNullPtr") {
			e = stripParens(call.Args[0])
		}
	}
	return isNonMapStructPtrType(info.Types[e].Type)
}

func validateCastAccessorCallSite(
	id *ast.Ident,
	file *ast.File,
	cc *mapCandidate,
	varCands map[*types.Var]*mapCandidate,
	consumedCastCalls map[*ast.CallExpr]bool,
	ctx *PassContext,
	parents map[ast.Node]ast.Node,
	varUseCounts map[*types.Var]int,
	markMut func(*mapCandidate, ast.Node, map[ast.Node]ast.Node),
	edges *[][2]*mapCandidate,
) {
	if !ctx.ShouldTransform(file) {
		cc.disqualified = true
		return
	}
	sel, ok := parents[id].(*ast.SelectorExpr)
	if !ok || sel.Sel != id {
		cc.disqualified = true
		return
	}
	call, ok := parents[sel].(*ast.CallExpr)
	if !ok || call.Fun != sel || len(call.Args) != 0 || call.Ellipsis.IsValid() {
		cc.disqualified = true
		return
	}
	cc.inPkgCastCalls++
	if consumedCastCalls[call] {
		return
	}
	cc.spellSites = append(cc.spellSites, typeSpellSite{file: file, pos: call.Pos(), owner: cc.ownerStruct})
	cc.j2gortSites = append(cc.j2gortSites, typeSpellSite{file: file, pos: call.Pos()})
	validateMapValueConsumer(call, file, cc, varCands, nil, nil, ctx, parents, ctx.Info, varUseCounts, markMut, edges)
}

func isReadOnlyEntrySetStreamUse(viewCall *ast.CallExpr, parents map[ast.Node]ast.Node, info *types.Info) bool {
	adaptedView := unwrapCollectionAdapters(viewCall, parents, info)
	streamSel, ok := parents[adaptedView].(*ast.SelectorExpr)
	if !ok || streamSel.X != adaptedView || streamSel.Sel == nil || streamSel.Sel.Name != "Stream" {
		return false
	}
	streamCall, ok := parents[streamSel].(*ast.CallExpr)
	return ok && streamCall.Fun == streamSel && len(streamCall.Args) == 0 && !streamCall.Ellipsis.IsValid()
}

func discoverMapFieldGetters(
	ctx *PassContext,
	varCands map[*types.Var]*mapCandidate,
	getterFieldCands map[*types.Func]*mapCandidate,
	getterBodyByField map[*types.Var]*ast.BlockStmt,
) {
	if ctx == nil || ctx.Pkg == nil || ctx.Pkg.Scope() == nil || ctx.Info == nil {
		return
	}
	info := ctx.Info
	isLeafJscomp := ctx.Pkg.Name() == "jscomp"

	// Collect all interface method names in ctx.Pkg and all structs + method counts to guard exported getters.
	ifaceMethods := make(map[string]bool)
	var pkgStructs []*types.Named
	methodCountByName := make(map[string]int)
	sc := ctx.Pkg.Scope()
	for _, nm := range sc.Names() {
		tn, ok := sc.Lookup(nm).(*types.TypeName)
		if !ok || tn.IsAlias() {
			continue
		}
		if named, isN := tn.Type().(*types.Named); isN {
			if st, isSt := named.Underlying().(*types.Struct); isSt && (named.TypeParams() == nil || named.TypeParams().Len() == 0) {
				_ = st
				pkgStructs = append(pkgStructs, named)
			}
			if iface, isIf := named.Underlying().(*types.Interface); isIf {
				for i := 0; i < iface.NumMethods(); i++ {
					ifaceMethods[iface.Method(i).Name()] = true
				}
			}
		}
	}
	for _, file := range ctx.Files {
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv != nil && fn.Name != nil {
				methodCountByName[fn.Name.Name]++
			}
		}
	}

	multiGetterField := make(map[*types.Var]bool)
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Body == nil || len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 {
				continue
			}
			if (fn.Type.TypeParams != nil && len(fn.Type.TypeParams.List) > 0) ||
				(fn.Type.Params != nil && len(fn.Type.Params.List) > 0) ||
				fn.Type.Results == nil || len(fn.Type.Results.List) != 1 || len(fn.Type.Results.List[0].Names) > 0 {
				continue
			}
			fnObj, ok := info.Defs[fn.Name].(*types.Func)
			if !ok || fnObj == nil {
				continue
			}
			if fnObj.Exported() {
				if !isLeafJscomp || ifaceMethods[fn.Name.Name] || methodCountByName[fn.Name.Name] != 1 {
					continue
				}
			}
			sig, ok := fnObj.Type().(*types.Signature)
			if !ok || sig.Recv() == nil || sig.RecvTypeParams() != nil {
				continue
			}
			recvId := fn.Recv.List[0].Names[0]
			if recvId.Name == "" || recvId.Name == "_" {
				continue
			}
			recvVar := varFromObj(info.Defs[recvId])
			if recvVar == nil {
				continue
			}
			owner := ptrNamedOrigin(recvVar.Type())
			if owner == nil || (owner.TypeParams() != nil && owner.TypeParams().Len() > 0) {
				continue
			}
			if fnObj.Exported() {
				embeddedElsewhere := false
				for _, other := range pkgStructs {
					if other != owner && structEmbedsNamed(other, owner) {
						embeddedElsewhere = true
						break
					}
				}
				if embeddedElsewhere {
					continue
				}
			}
			isCand, _, _ := classifyCandidateMapType(sig.Results().At(0).Type())
			if !isCand || len(fn.Body.List) == 0 {
				continue
			}
			allBlankPrefix := true
			for _, st := range fn.Body.List[:len(fn.Body.List)-1] {
				if !isBlankDiscardStmt(st) {
					allBlankPrefix = false
					break
				}
			}
			if !allBlankPrefix {
				continue
			}
			ret, ok := fn.Body.List[len(fn.Body.List)-1].(*ast.ReturnStmt)
			if !ok || len(ret.Results) != 1 {
				continue
			}
			sel, ok := stripParens(ret.Results[0]).(*ast.SelectorExpr)
			if !ok || !isIndirectStructFieldSel(sel, info) {
				continue
			}
			rId, ok := stripParens(sel.X).(*ast.Ident)
			if !ok || varFromObj(info.Uses[rId]) != recvVar {
				continue
			}
			fVar := varFromObj(info.Uses[sel.Sel])
			fieldCand := varCands[fVar]
			if fieldCand == nil || !fieldCand.isField || fieldCand.ownerStruct != owner || multiGetterField[fVar] {
				continue
			}
			if getterBodyByField[fVar] != nil {
				delete(getterBodyByField, fVar)
				multiGetterField[fVar] = true
				for k, v := range getterFieldCands {
					if v == fieldCand {
						delete(getterFieldCands, k)
					}
				}
				fieldCand.disqualified = true
				continue
			}
			fnOrig := funcOrigin(fnObj)
			getterFieldCands[fnOrig] = fieldCand
			getterBodyByField[fVar] = fn.Body
			fieldCand.spellSites = append(fieldCand.spellSites, typeSpellSite{file: file, pos: fn.Pos(), owner: owner})
			fieldCand.j2gortSites = append(fieldCand.j2gortSites, typeSpellSite{file: file, pos: fn.Pos()})
			fieldCand.ops = append(fieldCand.ops, mapRewriteOp{
				kind: mapOpDeclFieldGetter, cand: fieldCand, file: file, spanNode: fn,
				elemIdent: recvId, methodsSel: fVar.Name(),
			})
		}
	}
}

func validateMapFieldGetterCallSite(
	id *ast.Ident,
	file *ast.File,
	gfCand *mapCandidate,
	varCands map[*types.Var]*mapCandidate,
	paramCandsByFn map[*types.Func][]*mapCandidate,
	ctx *PassContext,
	parents map[ast.Node]ast.Node,
	varUseCounts map[*types.Var]int,
	markMut func(*mapCandidate, ast.Node, map[ast.Node]ast.Node),
	edges *[][2]*mapCandidate,
) {
	if !ctx.ShouldTransform(file) {
		gfCand.disqualified = true
		return
	}
	sel, ok := parents[id].(*ast.SelectorExpr)
	if !ok || sel.Sel != id {
		gfCand.disqualified = true
		return
	}
	call, ok := parents[sel].(*ast.CallExpr)
	if !ok || call.Fun != sel || len(call.Args) != 0 || call.Ellipsis.IsValid() {
		gfCand.disqualified = true
		return
	}
	gfCand.spellSites = append(gfCand.spellSites, typeSpellSite{file: file, pos: call.Pos(), owner: gfCand.ownerStruct})
	gfCand.j2gortSites = append(gfCand.j2gortSites, typeSpellSite{file: file, pos: call.Pos()})
	prevOps := len(gfCand.ops)
	validateMapValueConsumer(call, file, gfCand, varCands, nil, paramCandsByFn, ctx, parents, ctx.Info, varUseCounts, markMut, edges)
	for i := prevOps; i < len(gfCand.ops); i++ {
		gfCand.ops[i].viaPtrGetter = true
	}
}

func discoverMapFuncParams(
	ctx *PassContext,
	varCands map[*types.Var]*mapCandidate,
	paramCandsByFn map[*types.Func][]*mapCandidate,
	allCands *[]*mapCandidate,
) {
	if ctx == nil || ctx.Pkg == nil || ctx.Info == nil {
		return
	}
	info := ctx.Info
	isLeafJscomp := ctx.Pkg.Name() == "jscomp"

	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil || fn.Name == nil || fn.Type.Params == nil {
				continue
			}
			if fn.Type.TypeParams != nil && len(fn.Type.TypeParams.List) > 0 {
				continue
			}
			fnObj, ok := info.Defs[fn.Name].(*types.Func)
			if !ok || fnObj == nil {
				continue
			}
			sig, ok := fnObj.Type().(*types.Signature)
			if !ok || sig.Variadic() || sig.RecvTypeParams() != nil {
				continue
			}
			if fnObj.Exported() {
				if !isLeafJscomp || fn.Recv != nil ||
					(!strings.HasPrefix(fn.Name.Name, "Init") && !strings.HasPrefix(fn.Name.Name, "New")) {
					continue
				}
			}
			fnOrig := funcOrigin(fnObj)
			flatIdx := 0
			for _, fld := range fn.Type.Params.List {
				if len(fld.Names) != 1 {
					flatIdx += max(1, len(fld.Names))
					continue
				}
				pId := fld.Names[0]
				pVar := varFromObj(info.Defs[pId])
				if pId.Name != "_" && pVar != nil {
					if isCand, declIdent, declHash := classifyCandidateMapType(pVar.Type()); isCand {
						pc := &mapCandidate{
							v: pVar, isParam: true, paramFn: fnOrig, paramIdx: flatIdx,
							declFile: file, hasInit: true, forceNullable: true,
							hasIdentityKind: declIdent, hasHashKind: declHash,
							spellSites:  []typeSpellSite{{file: file, pos: fld.Pos()}},
							j2gortSites: []typeSpellSite{{file: file, pos: fld.Pos()}},
						}
						pc.ops = append(pc.ops, mapRewriteOp{kind: mapOpDeclParam, cand: pc, file: file, spanNode: fld})
						varCands[pVar] = pc
						paramCandsByFn[fnOrig] = append(paramCandsByFn[fnOrig], pc)
						*allCands = append(*allCands, pc)
					}
				}
				flatIdx++
			}
		}
	}
}

func validateMapParamFuncCallSite(
	id *ast.Ident,
	file *ast.File,
	pcs []*mapCandidate,
	varCands map[*types.Var]*mapCandidate,
	getterFieldCands map[*types.Func]*mapCandidate,
	ctx *PassContext,
	parents map[ast.Node]ast.Node,
) {
	disqualifyAll := func() {
		for _, pc := range pcs {
			pc.disqualified = true
		}
	}
	if !ctx.ShouldTransform(file) {
		disqualifyAll()
		return
	}
	var funNode ast.Node = id
	if sel, ok := parents[id].(*ast.SelectorExpr); ok && sel.Sel == id {
		funNode = sel
	}
	call, ok := parents[funNode].(*ast.CallExpr)
	if !ok || call.Fun != funNode || call.Ellipsis.IsValid() {
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
		argUnwrapped := unwrapMapAdaptersDown(call.Args[pc.paramIdx], ctx.Info)
		srcCand := resolveMapVarAccess(argUnwrapped, varCands, ctx.Info)
		if srcCand == nil {
			if gCall, isC := argUnwrapped.(*ast.CallExpr); isC && len(gCall.Args) == 0 && !gCall.Ellipsis.IsValid() {
				if callee, _ := calleeFunc(gCall, ctx.Info); callee != nil {
					srcCand = getterFieldCands[funcOrigin(callee)]
				}
			}
		}
		if srcCand == nil || srcCand.disqualified {
			pc.disqualified = true
		}
	}
}

func canPromoteLazyFieldToValueOrderedMap(
	m *mapCandidate,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	if !m.isField || m.isPtrAlias || m.isCastAccessor || m.v == nil {
		return false
	}
	lazyInitPutSpans := make(map[ast.Node]bool)
	nonNilInits := 0

	for i := range m.ops {
		op := &m.ops[i]
		parents := parentsByFile[op.file]
		switch op.kind {
		case mapOpDeclField, mapOpNilCheck:
			continue
		case mapOpRemove, mapOpComputeIfAbsent:
			return false
		case mapOpMethodDirect:
			if op.methodsSel == "Clear" {
				return false
			}
		case mapOpPut:
			if op.isPutIfAbsent {
				return false
			}
		case mapOpAssignInit:
			if op.initKind == mapInitNil {
				if !isInInitFieldsMethod(op.stmt, m.ownerStruct, parents, info) {
					return false
				}
				continue
			}
			if op.initKind != mapInitEmpty && op.initKind != mapInitCap {
				return false
			}
			putSpan, ok := findImmediateLazyPutAfterInit(op.stmt, op.accessExpr, m, parents, info)
			if !ok {
				return false
			}
			lazyInitPutSpans[putSpan] = true
			nonNilInits++
		}
	}
	if nonNilInits == 0 {
		return false
	}

	// Verify all other operations on m are guarded by `m != nil` or precede-return `if m == nil { ... return }`.
	for i := range m.ops {
		op := &m.ops[i]
		switch op.kind {
		case mapOpDeclField, mapOpNilCheck, mapOpAssignInit:
			continue
		}
		if lazyInitPutSpans[op.spanNode] {
			continue
		}
		if op.accessExpr == nil || !isGuardedByFieldNilCheck(op.spanNode, op.accessExpr, m.v, parentsByFile[op.file], info) {
			return false
		}
	}
	return true
}

func findImmediateLazyPutAfterInit(
	initStmt ast.Stmt,
	initAccess ast.Expr,
	m *mapCandidate,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) (ast.Node, bool) {
	body, ok := parents[initStmt].(*ast.BlockStmt)
	if !ok || len(body.List) == 0 || body.List[0] != initStmt {
		return nil, false
	}
	ifStmt, ok := parents[body].(*ast.IfStmt)
	if !ok || ifStmt.Body != body || ifStmt.Else != nil || !isFieldEqNilCond(ifStmt.Cond, initAccess, m.v, token.EQL, info) {
		return nil, false
	}
	var nextStmt ast.Stmt
	if len(body.List) >= 2 {
		nextStmt = body.List[1]
	} else {
		outerBlock, ok := parents[ifStmt].(*ast.BlockStmt)
		if !ok {
			return nil, false
		}
		idx := slices.Index(outerBlock.List, ast.Stmt(ifStmt))
		if idx < 0 || idx+1 >= len(outerBlock.List) {
			return nil, false
		}
		nextStmt = outerBlock.List[idx+1]
	}
	es, ok := nextStmt.(*ast.ExprStmt)
	if !ok {
		return nil, false
	}
	for i := range m.ops {
		op := &m.ops[i]
		if op.kind == mapOpPut && !op.isPutIfAbsent && (op.spanNode == es || op.spanNode == es.X) &&
			sameAccessPath(initAccess, op.accessExpr, m.v, info) {
			return op.spanNode, true
		}
	}
	return nil, false
}

func isFieldEqNilCond(cond ast.Expr, wantAccess ast.Expr, fieldVar *types.Var, wantOp token.Token, info *types.Info) bool {
	bin, ok := stripParens(cond).(*ast.BinaryExpr)
	if !ok || bin.Op != wantOp {
		return false
	}
	var fieldSide, nilSide ast.Expr
	if isNilIdent(bin.Y, info) {
		fieldSide, nilSide = bin.X, bin.Y
	} else if isNilIdent(bin.X, info) {
		fieldSide, nilSide = bin.Y, bin.X
	}
	return nilSide != nil && sameAccessPath(wantAccess, fieldSide, fieldVar, info)
}

func isGuardedByFieldNilCheck(
	node ast.Node,
	accessExpr ast.Expr,
	fieldVar *types.Var,
	parents map[ast.Node]ast.Node,
	info *types.Info,
) bool {
	for cur := node; cur != nil; cur = parents[cur] {
		switch p := parents[cur].(type) {
		case *ast.FuncLit:
			return false
		case *ast.IfStmt:
			if p.Body == cur && isFieldEqNilCond(p.Cond, accessExpr, fieldVar, token.NEQ, info) {
				return true
			}
			if p.Else == cur && isFieldEqNilCond(p.Cond, accessExpr, fieldVar, token.EQL, info) {
				return true
			}
		case *ast.BlockStmt:
			curStmt, isStmt := cur.(ast.Stmt)
			if !isStmt {
				continue
			}
			idx := slices.Index(p.List, curStmt)
			for i := 0; i < idx; i++ {
				if prevIf, ok := p.List[i].(*ast.IfStmt); ok && prevIf.Else == nil &&
					isFieldEqNilCond(prevIf.Cond, accessExpr, fieldVar, token.EQL, info) &&
					blockAlwaysTerminatesForMapGuard(prevIf.Body) {
					return true
				}
			}
		}
	}
	return false
}

func blockAlwaysTerminatesForMapGuard(b *ast.BlockStmt) bool {
	if b == nil || len(b.List) == 0 {
		return false
	}
	switch last := b.List[len(b.List)-1].(type) {
	case *ast.ReturnStmt:
		return true
	case *ast.BranchStmt:
		return last.Tok == token.RETURN || last.Tok == token.BREAK || last.Tok == token.CONTINUE
	case *ast.ExprStmt:
		if call, ok := stripParens(last.X).(*ast.CallExpr); ok {
			if id, ok := stripParens(call.Fun).(*ast.Ident); ok && id.Name == "panic" {
				return true
			}
		}
	case *ast.IfStmt:
		if els, ok := last.Else.(*ast.BlockStmt); ok {
			return blockAlwaysTerminatesForMapGuard(last.Body) && blockAlwaysTerminatesForMapGuard(els)
		}
	}
	return false
}

func j2gortOrderedMapHasMethod(ctx *PassContext, methodName string) bool {
	if ctx == nil || ctx.Importer == nil {
		return false
	}
	for _, f := range ctx.Files {
		for _, spec := range f.Imports {
			if spec.Path == nil || len(spec.Path.Value) < 3 {
				continue
			}
			p := spec.Path.Value[1 : len(spec.Path.Value)-1]
			if impPkg, err := ctx.Importer.Import(p); err == nil && isJ2gortPkg(impPkg) {
				if impPkg.Scope() == nil {
					return false
				}
				tn, ok := impPkg.Scope().Lookup("OrderedMap").(*types.TypeName)
				if !ok || tn == nil {
					return false
				}
				named, ok := tn.Type().(*types.Named)
				if !ok || named == nil || named.TypeParams() == nil || named.TypeParams().Len() != 2 {
					return false
				}
				inst, err := types.Instantiate(nil, named, []types.Type{types.Typ[types.String], types.Typ[types.String]}, false)
				if err != nil {
					return false
				}
				return types.NewMethodSet(types.NewPointer(inst)).Lookup(impPkg, methodName) != nil
			}
		}
	}
	return false
}

// isNonEscapingReadOnlyLiveViewConsumer reports whether adaptedView (wrapping m.KeySet_0() or m.Values())
// is consumed synchronously without escaping (e.g. `.ContainsAll(...)`, `.Contains(...)`, `.IsEmpty()`,
// `.Size()`, `.Stream()`, or `Sets_Difference(a, b).ImmutableCopy()`).
func isNonEscapingReadOnlyLiveViewConsumer(adaptedView ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) bool {
	if adaptedView == nil || parents == nil {
		return false
	}
	switch p := parents[adaptedView].(type) {
	case *ast.SelectorExpr:
		if p.X != adaptedView || p.Sel == nil {
			return false
		}
		call, ok := parents[p].(*ast.CallExpr)
		if !ok || call.Fun != p || call.Ellipsis.IsValid() {
			return false
		}
		switch p.Sel.Name {
		case "ContainsAll", "Contains":
			return len(call.Args) == 1
		case "IsEmpty", "Size", "Stream":
			return len(call.Args) == 0
		}
	case *ast.CallExpr:
		if p.Fun == adaptedView || p.Ellipsis.IsValid() {
			return false
		}
		if sel, ok := stripParens(p.Fun).(*ast.SelectorExpr); ok && sel.Sel != nil && sel.Sel.Name == "ContainsAll" && len(p.Args) == 1 && p.Args[0] == adaptedView {
			return isCollectionOrSetIface(info.Types[sel.X].Type)
		}
		if fnName, _, _, isJ2 := isJ2gortCall(p, info); isJ2 && len(p.Args) == 2 && (p.Args[0] == adaptedView || p.Args[1] == adaptedView) {
			switch fnName {
			case "Ext_com_google_common_collect_Sets_Difference",
				"Ext_com_google_common_collect_Sets_Intersection",
				"Ext_com_google_common_collect_Sets_Union":
				topSetView := unwrapParensUp(p, parents)
				if mSel, ok := parents[topSetView].(*ast.SelectorExpr); ok && mSel.X == topSetView && mSel.Sel != nil && mSel.Sel.Name == "ImmutableCopy" {
					if mCall, okC := parents[mSel].(*ast.CallExpr); okC && mCall.Fun == mSel && len(mCall.Args) == 0 {
						return true
					}
				}
			}
		}
	}
	return false
}

// applyTimelineMapAndSliceLowering lowers jscomp.Timeline's eventsByTime and node-keyed eventsByValue
// maps to j2gort.OrderedMap[string, *Timeline_Event] and j2gort.OrderedMap[pkg_rhino.Node, *Timeline_Event]
// (with lazy Ext_java_util_LinkedHashMap fallback for arbitrary non-Node callers), chunk-allocates/recycles
// Timeline_Event records, lowers GetSince to []pkg_rhino.Node + WrapTypedArrayList, and rewrites
// (*Timeline).Add(BoxNode(n)) / Remove(BoxNode(n)) call sites to zero-boxing AddNode_(n) / RemoveNode_(n).
func applyTimelineMapAndSliceLowering(ctx *PassContext) (int, error) {
	if ctx == nil || ctx.Pkg == nil || ctx.Pkg.Name() != "jscomp" || ctx.Info == nil {
		return 0, nil
	}
	if !j2gortHasExport(ctx, "OrderedMap") || !j2gortHasExport(ctx, "WrapTypedArrayList") ||
		!j2gortHasExport(ctx, "BoxNode") || !j2gortHasExport(ctx, "NullString") || !j2gortHasExport(ctx, "ThrowNPE") {
		return 0, nil
	}
	info := ctx.Info

	var tlNamed, tlEventNamed *types.Named
	var tlStructAST, tlEventStructAST *ast.StructType
	var tlFile, tlEventFile *ast.File
	for _, file := range ctx.Files {
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Assign.IsValid() || ts.Name == nil {
					continue
				}
				stAST, stOK := ts.Type.(*ast.StructType)
				tn, tnOK := info.Defs[ts.Name].(*types.TypeName)
				if !stOK || !tnOK || stAST.Fields == nil {
					continue
				}
				named := namedOrigin(tn.Type())
				if named == nil {
					continue
				}
				switch ts.Name.Name {
				case "Timeline":
					tlNamed = named
					tlStructAST = stAST
					tlFile = file
				case "Timeline_Event":
					tlEventNamed = named
					tlEventStructAST = stAST
					tlEventFile = file
				}
			}
		}
	}
	if tlNamed == nil || tlEventNamed == nil || tlFile == nil || tlFile != tlEventFile {
		return 0, nil
	}

	// Resolve pkg_rhino.Node type from imported rhino package or BoxNode call sites.
	var rhinoPkg *types.Package
	for _, imp := range ctx.Pkg.Imports() {
		if imp != nil && imp.Name() == "rhino" && imp.Scope() != nil {
			if obj := imp.Scope().Lookup("Node"); obj != nil && isNamedUint32NodeType(obj.Type()) {
				rhinoPkg = imp
				break
			}
		}
	}
	if rhinoPkg == nil {
		return 0, nil
	}

	count := 0
	tlAlreadyLowered := hasASTFieldNamed(tlStructAST, "eventsByNode_")

	if !tlAlreadyLowered && ctx.ShouldTransform(tlFile) && !ctx.ModifiedFiles[tlFile] {
		var byTimeField, byValField, headField *ast.Field
		for _, fld := range tlStructAST.Fields.List {
			for _, nm := range fld.Names {
				switch nm.Name {
				case "eventsByTime":
					byTimeField = fld
				case "eventsByValue":
					byValField = fld
				case "headEvent":
					headField = fld
				}
			}
		}
		if byTimeField != nil && len(byTimeField.Names) == 1 &&
			byValField != nil && len(byValField.Names) == 1 &&
			headField != nil && hasASTFieldNamed(tlEventStructAST, "nextEvent") &&
			hasASTFieldNamed(tlEventStructAST, "previousEvent") && hasASTFieldNamed(tlEventStructAST, "value") {

			// Verify unexported Timeline & Timeline_Event fields are only accessed inside tlFile.
			watchedFields := make(map[*types.Var]bool)
			for _, stNamed := range []*types.Named{tlNamed, tlEventNamed} {
				if st, ok := stNamed.Underlying().(*types.Struct); ok {
					for i := 0; i < st.NumFields(); i++ {
						f := st.Field(i)
						if f.Name() != "idHash_" {
							watchedFields[f] = true
						}
					}
				}
			}

			var initFn, addFn, removeFn, markFn, getSinceFn, addEventFn *ast.FuncDecl
			var addEventObj *types.Func
			for _, decl := range tlFile.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Name == nil || fn.Body == nil || fn.Recv == nil || len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 {
					continue
				}
				rId := fn.Recv.List[0].Names[0]
				if rId.Name == "" || rId.Name == "_" {
					continue
				}
				rv := varFromObj(info.Defs[rId])
				if rv == nil || ptrNamedOrigin(rv.Type()) != tlNamed {
					continue
				}
				switch {
				case fn.Name.Name == "initFields_Timeline":
					initFn = fn
				case fn.Name.Name == "Add":
					addFn = fn
				case fn.Name.Name == "Remove":
					removeFn = fn
				case fn.Name.Name == "Mark":
					markFn = fn
				case fn.Name.Name == "GetSince":
					getSinceFn = fn
				case strings.HasPrefix(fn.Name.Name, "priv_Timeline_addEvent"):
					addEventFn = fn
					addEventObj, _ = info.Defs[fn.Name].(*types.Func)
				}
			}

			encapsulated := initFn != nil && addFn != nil && removeFn != nil && markFn != nil &&
				getSinceFn != nil && addEventFn != nil && addEventObj != nil
			if encapsulated {
				for _, f := range ctx.Files {
					if f == tlFile {
						continue
					}
					ast.Inspect(f, func(n ast.Node) bool {
						if !encapsulated || n == nil {
							return false
						}
						if id, ok := n.(*ast.Ident); ok {
							if fv := varFromObj(info.Uses[id]); fv != nil && watchedFields[fv] {
								encapsulated = false
								return false
							}
							if fnU, okF := info.Uses[id].(*types.Func); okF && fnU == addEventObj {
								encapsulated = false
								return false
							}
						}
						return true
					})
					if !encapsulated {
						break
					}
				}
			}

			initRecv := ""
			if initFn != nil {
				initRecv = initFn.Recv.List[0].Names[0].Name
			}
			addParams := extractFlatParamNames(safeFuncParams(addFn))
			remParams := extractFlatParamNames(safeFuncParams(removeFn))
			markParams := extractFlatParamNames(safeFuncParams(markFn))
			sinceParams := extractFlatParamNames(safeFuncParams(getSinceFn))
			aeParams := extractFlatParamNames(safeFuncParams(addEventFn))

			if encapsulated && len(addParams) == 1 && len(remParams) == 1 &&
				len(markParams) == 1 && len(sinceParams) == 1 && len(aeParams) == 2 {
				rhinoQual, okR := ensurePkgImportInFile(tlFile, rhinoPkg.Path(), "pkg_rhino", ctx)
				j2Qual, okJ := ensureJ2gortImportInFile(tlFile, ctx)
				if okR && okJ && rhinoQual != "" && j2Qual != "" {
					byTimeField.Type = &ast.IndexListExpr{
						X: &ast.SelectorExpr{X: &ast.Ident{Name: j2Qual}, Sel: &ast.Ident{Name: "OrderedMap"}},
						Indices: []ast.Expr{
							&ast.Ident{Name: "string"},
							&ast.StarExpr{X: &ast.Ident{Name: "Timeline_Event"}},
						},
					}
					byValField.Type = &ast.StarExpr{
						X: &ast.SelectorExpr{X: &ast.Ident{Name: j2Qual}, Sel: &ast.Ident{Name: "Ext_java_util_LinkedHashMap"}},
					}
					tlStructAST.Fields.List = append(tlStructAST.Fields.List,
						&ast.Field{
							Names: []*ast.Ident{{Name: "eventsByNode_"}},
							Type: &ast.IndexListExpr{
								X: &ast.SelectorExpr{X: &ast.Ident{Name: j2Qual}, Sel: &ast.Ident{Name: "OrderedMap"}},
								Indices: []ast.Expr{
									&ast.SelectorExpr{X: &ast.Ident{Name: rhinoQual}, Sel: &ast.Ident{Name: "Node"}},
									&ast.StarExpr{X: &ast.Ident{Name: "Timeline_Event"}},
								},
							},
						},
						&ast.Field{
							Names: []*ast.Ident{{Name: "eventBuf_"}},
							Type:  &ast.ArrayType{Elt: &ast.Ident{Name: "Timeline_Event"}},
						},
						&ast.Field{
							Names: []*ast.Ident{{Name: "freeEvent_"}},
							Type:  &ast.StarExpr{X: &ast.Ident{Name: "Timeline_Event"}},
						},
					)
					tlEventStructAST.Fields.List = append(tlEventStructAST.Fields.List,
						&ast.Field{
							Names: []*ast.Ident{{Name: "node_"}},
							Type:  &ast.SelectorExpr{X: &ast.Ident{Name: rhinoQual}, Sel: &ast.Ident{Name: "Node"}},
						},
						&ast.Field{
							Names: []*ast.Ident{{Name: "kind_"}},
							Type:  &ast.Ident{Name: "uint8"},
						},
					)

					setBody := func(fn *ast.FuncDecl, bodySrc string) error {
						pf, err := parser.ParseFile(ctx.Fset, "", "package p\nfunc f() {\n"+bodySrc+"\n}\n", 0)
						if err != nil {
							return err
						}
						fn.Body.List = pf.Decls[0].(*ast.FuncDecl).Body.List
						return nil
					}

					if err := setBody(initFn, fmt.Sprintf(`%[1]s.eventsByTime = %[2]s.NewOrderedMap[string, *Timeline_Event]()
%[1]s.eventsByNode_ = %[2]s.NewOrderedMap[%[3]s.Node, *Timeline_Event]()
%[1]s.eventsByValue = nil
ev0_ := %[1]s.allocEvent_()
ev0_.kind_ = 1
%[1]s.headEvent = ev0_`, initRecv, j2Qual, rhinoQual)); err != nil {
						return count, err
					}

					addRecv := addFn.Recv.List[0].Names[0].Name
					if err := setBody(addFn, fmt.Sprintf(`if n_, ok_ := %[2]s.(%[3]s.Node); ok_ {
	%[1]s.AddNode_(n_)
	return
}
if %[1]s == nil || %[2]s == nil {
	%[4]s.ThrowNPE()
}
if %[1]s.eventsByValue == nil {
	%[1]s.eventsByValue = %[4]s.NewExt_java_util_LinkedHashMap()
}
%[1]s.%[5]s(%[2]s, %[1]s.eventsByValue)`, addRecv, addParams[0], rhinoQual, j2Qual, addEventFn.Name.Name)); err != nil {
						return count, err
					}

					remRecv := removeFn.Recv.List[0].Names[0].Name
					if err := setBody(removeFn, fmt.Sprintf(`if n_, ok_ := %[2]s.(%[3]s.Node); ok_ {
	%[1]s.RemoveNode_(n_)
	return
}
if %[1]s == nil {
	%[4]s.ThrowNPE()
}
if %[2]s == nil || %[1]s.eventsByValue == nil {
	return
}
event, _ := %[1]s.eventsByValue.Remove_1(%[2]s).(*Timeline_Event)
if event != nil {
	%[1]s.unlinkAndFreeEvent_(event)
}`, remRecv, remParams[0], rhinoQual, j2Qual)); err != nil {
						return count, err
					}

					markRecv := markFn.Recv.List[0].Names[0].Name
					if err := setBody(markFn, fmt.Sprintf(`if %[1]s == nil || %[2]s == %[3]s.NullString {
	%[3]s.ThrowNPE()
}
event := %[1]s.eventsByTime.Get(%[2]s)
if %[1]s.headEvent == event {
	return
}
if event != nil {
	event.previousEvent.nextEvent = event.nextEvent
	event.nextEvent.previousEvent = event.previousEvent
	event.nextEvent = nil
} else {
	event = %[1]s.allocEvent_()
	event.kind_ = 1
	%[1]s.eventsByTime.Put(%[2]s, event)
}
event.previousEvent = %[1]s.headEvent
%[1]s.headEvent.nextEvent = event
%[1]s.headEvent = event`, markRecv, markParams[0], j2Qual)); err != nil {
						return count, err
					}

					sinceRecv := getSinceFn.Recv.List[0].Names[0].Name
					if err := setBody(getSinceFn, fmt.Sprintf(`if %[1]s == nil || %[2]s == %[3]s.NullString {
	%[3]s.ThrowNPE()
}
firstEvent := %[1]s.eventsByTime.Get(%[2]s)
if firstEvent == nil {
	return nil
}
if %[1]s.eventsByValue == nil {
	values := []%[4]s.Node{}
	for event := firstEvent; event != nil; event = event.nextEvent {
		if event.kind_ == 2 {
			values = append(values, event.node_)
		}
	}
	return %[3]s.WrapTypedArrayList(values)
}
values := %[3]s.NewExt_java_util_ArrayList()
for event := firstEvent; event != nil; event = event.nextEvent {
	if event.kind_ == 2 {
		values.Add_1(%[3]s.BoxNode(event.node_))
	} else if event.kind_ == 0 {
		values.Add_1(event.value)
	}
}
return values`, sinceRecv, sinceParams[0], j2Qual, rhinoQual)); err != nil {
						return count, err
					}

					aeRecv := addEventFn.Recv.List[0].Names[0].Name
					if err := setBody(addEventFn, fmt.Sprintf(`event, _ := %[3]s.Get(%[2]s).(*Timeline_Event)
if %[1]s.headEvent == event {
	return
}
if event != nil {
	event.previousEvent.nextEvent = event.nextEvent
	event.nextEvent.previousEvent = event.previousEvent
	event.nextEvent = nil
} else {
	event = %[1]s.allocEvent_()
	event.value = %[2]s
	event.kind_ = 0
	%[3]s.Put_2(%[2]s, event)
}
event.previousEvent = %[1]s.headEvent
%[1]s.headEvent.nextEvent = event
%[1]s.headEvent = event`, aeRecv, aeParams[0], aeParams[1])); err != nil {
						return count, err
					}

					helpersPF, err := parser.ParseFile(ctx.Fset, "", fmt.Sprintf(`package p

func (this_ *Timeline) allocEvent_() *Timeline_Event {
	if ev := this_.freeEvent_; ev != nil {
		this_.freeEvent_ = ev.nextEvent
		ev.nextEvent = nil
		return ev
	}
	if len(this_.eventBuf_) == cap(this_.eventBuf_) {
		this_.eventBuf_ = make([]Timeline_Event, 0, 64)
	}
	i := len(this_.eventBuf_)
	this_.eventBuf_ = this_.eventBuf_[:i+1]
	return &this_.eventBuf_[i]
}

func (this_ *Timeline) unlinkAndFreeEvent_(event *Timeline_Event) {
	if event.nextEvent != nil {
		event.nextEvent.previousEvent = event.previousEvent
	} else {
		this_.headEvent = event.previousEvent
	}
	event.previousEvent.nextEvent = event.nextEvent
	event.previousEvent = nil
	event.value = nil
	event.node_ = 0
	event.kind_ = 0
	event.nextEvent = this_.freeEvent_
	this_.freeEvent_ = event
}

func (this_ *Timeline) AddNode_(n %[2]s.Node) {
	if this_ == nil || n == 0 {
		%[1]s.ThrowNPE()
	}
	event := this_.eventsByNode_.Get(n)
	if this_.headEvent == event {
		return
	}
	if event != nil {
		event.previousEvent.nextEvent = event.nextEvent
		event.nextEvent.previousEvent = event.previousEvent
		event.nextEvent = nil
	} else {
		event = this_.allocEvent_()
		event.node_ = n
		event.kind_ = 2
		this_.eventsByNode_.Put(n, event)
	}
	event.previousEvent = this_.headEvent
	this_.headEvent.nextEvent = event
	this_.headEvent = event
}

func (this_ *Timeline) RemoveNode_(n %[2]s.Node) {
	if this_ == nil {
		%[1]s.ThrowNPE()
	}
	if n == 0 {
		return
	}
	event := this_.eventsByNode_.Remove(n)
	if event != nil {
		this_.unlinkAndFreeEvent_(event)
	}
}
`, j2Qual, rhinoQual), 0)
					if err != nil {
						return count, err
					}
					tlFile.Decls = append(tlFile.Decls, helpersPF.Decls...)
					ctx.MarkFileModified(tlFile)
					tlAlreadyLowered = true
					count += 4
				}
			}
		}
	}

	if !tlAlreadyLowered {
		return count, nil
	}

	// Rewrite (*Timeline).Add(j2gort.BoxNode(n)) -> (*Timeline).AddNode_(n)
	// and (*Timeline).Remove(j2gort.BoxNode(n)) -> (*Timeline).RemoveNode_(n) across ctx.Files.
	for _, f := range ctx.Files {
		if !ctx.ShouldTransform(f) || (f != tlFile && ctx.ModifiedFiles[f]) {
			continue
		}
		fileChanged := false
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || call.Ellipsis.IsValid() || len(call.Args) != 1 {
				return true
			}
			sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
			if !ok || sel.Sel == nil || (sel.Sel.Name != "Add" && sel.Sel.Name != "Remove") {
				return true
			}
			if ptrNamedOrigin(info.Types[sel.X].Type) != tlNamed {
				return true
			}
			argCall, ok := stripParens(call.Args[0]).(*ast.CallExpr)
			if !ok || argCall.Ellipsis.IsValid() || len(argCall.Args) != 1 {
				return true
			}
			fnName, _, _, isJ2 := isJ2gortCall(argCall, info)
			if !isJ2 || fnName != "BoxNode" {
				return true
			}
			nodeT := info.Types[argCall.Args[0]].Type
			if namedOrigin(nodeT) == nil || namedOrigin(nodeT).Obj() == nil || namedOrigin(nodeT).Obj().Pkg() != rhinoPkg || namedOrigin(nodeT).Obj().Name() != "Node" {
				return true
			}
			if sel.Sel.Name == "Add" {
				sel.Sel = &ast.Ident{NamePos: sel.Sel.Pos(), Name: "AddNode_"}
			} else {
				sel.Sel = &ast.Ident{NamePos: sel.Sel.Pos(), Name: "RemoveNode_"}
			}
			call.Args[0] = argCall.Args[0]
			fileChanged = true
			count++
			return true
		})
		if fileChanged {
			ctx.MarkFileModified(f)
		}
	}

	return count, nil
}

func safeFuncParams(fn *ast.FuncDecl) *ast.FieldList {
	if fn == nil || fn.Type == nil {
		return nil
	}
	return fn.Type.Params
}

// ---------------------------------------------------------------------------
// Wave 37B: Unboxed int32 Specialization for OrderedMap / CustomOrderedMap / IdentityOrderedMap
// Holding Non-Null Java Integer (*j2gort.BoxInt) Values
// ---------------------------------------------------------------------------

type unboxIntMapVarSpec struct {
	v            *types.Var
	cand         *unboxIntMapCand
	fn           *ast.FuncDecl
	file         *ast.File
	declStmt     ast.Stmt
	okName       string
	hasReassign  bool
	reassigns    []*ast.AssignStmt // `v = j2gort.ValueOfInt(expr)` -> `v, vOk_ = expr, true`
	nilBins      []*ast.BinaryExpr // `v == nil` / `v != nil`
	nilIsEq      []bool
	unboxCalls   []*ast.CallExpr // `j2gort.UnboxInt([CheckNotNullPtr](v))`
	ptrToAnyCall []*ast.CallExpr // `j2gort.PtrToAny(v)` -> `func() any { if !vOk_ { return nil }; return j2gort.ValueOfInt(v) }()`
	putCalls     []*ast.CallExpr // `m.Put(k, v)` immediately after `v = j2gort.ValueOfInt(...)`
}

type unboxIntMapCand struct {
	v             *types.Var
	isField       bool
	isPtr         bool
	declFile      *ast.File
	defCount      int
	hasActiveOp   bool
	disqualified  bool
	touchedFiles  map[*ast.File]bool
	indexValExprs []*ast.IndexListExpr // `OrderedMap[K, *j2gort.BoxInt]` / `NewOrderedMap[K, *j2gort.BoxInt]` -> replace Indices[1] with `int32`
	putCalls      []*ast.CallExpr      // `m.Put(k, j2gort.ValueOfInt(x))` / `PutIfAbsent` -> unwrap arg 1
	directGetUnb  []*unboxDirectGetOp  // `j2gort.UnboxInt([CheckNotNullPtr](m.Get(k)))`
	localGetVars  []*unboxIntMapVarSpec
}

type unboxDirectGetOp struct {
	file      *ast.File
	outerCall *ast.CallExpr
	getCall   *ast.CallExpr
}

func isJ2gortOrderedMapBoxIntType(t types.Type) (isPtr bool, ok bool) {
	if t == nil {
		return false, false
	}
	ut := types.Unalias(t)
	if p, pOK := ut.(*types.Pointer); pOK {
		isPtr = true
		ut = types.Unalias(p.Elem())
	}
	nm, nmOK := ut.(*types.Named)
	if !nmOK || nm.Obj() == nil || !isJ2gortPkg(nm.Obj().Pkg()) {
		return false, false
	}
	switch nm.Obj().Name() {
	case "OrderedMap", "CustomOrderedMap", "IdentityOrderedMap", "PlainMap":
	default:
		return false, false
	}
	targs := nm.TypeArgs()
	if targs == nil || targs.Len() != 2 {
		return false, false
	}
	return isPtr, isPtrToJ2gortBoxInt(targs.At(1))
}

func isPtrToJ2gortBoxInt(t types.Type) bool {
	if t == nil {
		return false
	}
	p, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return false
	}
	nm, ok := types.Unalias(p.Elem()).(*types.Named)
	return ok && nm.Obj() != nil && isJ2gortPkg(nm.Obj().Pkg()) && nm.Obj().Name() == "BoxInt"
}

func unwrapValueOfIntCall(e ast.Expr, info *types.Info) (ast.Expr, bool) {
	call, ok := stripParens(e).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() || len(call.Args) != 1 {
		return nil, false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 || fnName != "ValueOfInt" {
		return nil, false
	}
	arg := call.Args[0]
	if tv, okT := info.Types[arg]; !okT || !types.Identical(types.Unalias(tv.Type), types.Typ[types.Int32]) {
		return nil, false
	}
	return arg, true
}

func enclosingUnboxIntOfExpr(inner ast.Expr, parents map[ast.Node]ast.Node, info *types.Info) *ast.CallExpr {
	cur := unwrapParensUp(inner, parents)
	if chkCall, ok := parents[cur].(*ast.CallExpr); ok && len(chkCall.Args) == 1 && !chkCall.Ellipsis.IsValid() && chkCall.Args[0] == cur {
		if fnName, _, _, isJ2 := isJ2gortCall(chkCall, info); isJ2 && (fnName == "CheckNotNullPtr" || fnName == "CheckNotNull") {
			cur = unwrapParensUp(chkCall, parents)
		}
	}
	unbCall, ok := parents[cur].(*ast.CallExpr)
	if !ok || len(unbCall.Args) != 1 || unbCall.Ellipsis.IsValid() || unbCall.Args[0] != cur {
		return nil
	}
	fnName, _, _, isJ2 := isJ2gortCall(unbCall, info)
	if !isJ2 || fnName != "UnboxInt" {
		return nil
	}
	return unbCall
}

func applyOrderedMapUnboxedIntLowering(ctx *PassContext, parentsByFile map[*ast.File]map[ast.Node]ast.Node) int {
	if ctx == nil || ctx.Info == nil || len(ctx.Files) == 0 || !j2gortHasExport(ctx, "ThrowNPE") {
		return 0
	}
	info := ctx.Info
	cands := make(map[*types.Var]*unboxIntMapCand)
	var candList []*unboxIntMapCand

	addCand := func(v *types.Var, isField, isPtr bool, file *ast.File, ile *ast.IndexListExpr) *unboxIntMapCand {
		c := &unboxIntMapCand{
			v:            v,
			isField:      isField,
			isPtr:        isPtr,
			declFile:     file,
			touchedFiles: map[*ast.File]bool{file: true},
		}
		if ile != nil {
			if !isUnshadowedUniverseType("int32", ile.Indices[1].Pos(), ctx.Pkg) {
				c.disqualified = true
			}
			c.indexValExprs = append(c.indexValExprs, ile)
		}
		cands[v] = c
		candList = append(candList, c)
		return c
	}

	extractOrdMapILE := func(tExpr ast.Expr) *ast.IndexListExpr {
		if tExpr == nil {
			return nil
		}
		cur := stripParens(tExpr)
		if st, ok := cur.(*ast.StarExpr); ok {
			cur = stripParens(st.X)
		}
		if ile, ok := cur.(*ast.IndexListExpr); ok && len(ile.Indices) == 2 {
			return ile
		}
		return nil
	}

	extractInitILE := func(rhs ast.Expr) (*ast.IndexListExpr, bool) {
		cur := stripParens(rhs)
		if un, ok := cur.(*ast.UnaryExpr); ok && un.Op == token.AND {
			cur = stripParens(un.X)
		}
		switch x := cur.(type) {
		case *ast.CompositeLit:
			if len(x.Elts) == 0 {
				if ile := extractOrdMapILE(x.Type); ile != nil {
					return ile, true
				}
			}
		case *ast.CallExpr:
			if ile, ok := stripParens(x.Fun).(*ast.IndexListExpr); ok && len(ile.Indices) == 2 {
				fnName, _, _, isJ2 := isJ2gortCall(x, info)
				if isJ2 {
					switch fnName {
					case "NewOrderedMap", "NewOrderedMapCap", "NewOrderedMapPtrCap", "AcquireOrderedMap", "AcquireOrderedMapCap",
						"NewCustomOrderedMap", "NewCustomOrderedMapCap", "NewCustomOrderedMapPtrCap", "AcquireCustomOrderedMap", "AcquireCustomOrderedMapCap",
						"NewIdentityOrderedMap", "NewIdentityOrderedMapCap", "NewIdentityOrderedMapPtrCap", "AcquireIdentityOrderedMap", "AcquireIdentityOrderedMapCap",
						"NewPlainMap", "NewPlainMapVal", "NewPlainMapCap", "NewPlainMapPtr", "NewPlainMapPtrCap":
						return ile, true
					}
				}
			}
		}
		return nil, false
	}

	// Pass 1: Collect candidate OrderedMap[K, *j2gort.BoxInt] unexported fields and local variables.
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
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
				st, ok := ts.Type.(*ast.StructType)
				if !ok || st.Fields == nil {
					continue
				}
				for _, fld := range st.Fields.List {
					if len(fld.Names) != 1 {
						continue
					}
					v := varFromObj(info.Defs[fld.Names[0]])
					if v == nil || !v.IsField() || v.Exported() || v.Embedded() {
						continue
					}
					isPtr, okM := isJ2gortOrderedMapBoxIntType(v.Type())
					if !okM {
						continue
					}
					ile := extractOrdMapILE(fld.Type)
					if ile == nil {
						continue
					}
					addCand(v, true, isPtr, file, ile)
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
				if !ok || len(vs.Names) != 1 || len(vs.Values) > 1 || vs.Names[0].Name == "_" {
					return true
				}
				v := varFromObj(info.Defs[vs.Names[0]])
				if v == nil || v.IsField() || cands[v] != nil {
					return true
				}
				isPtr, okM := isJ2gortOrderedMapBoxIntType(v.Type())
				if !okM {
					return true
				}
				var declILE *ast.IndexListExpr
				if vs.Type != nil {
					declILE = extractOrdMapILE(vs.Type)
					if declILE == nil {
						return true
					}
				}
				c := addCand(v, false, isPtr, file, declILE)
				if len(vs.Values) == 1 {
					if isPtr && isNilIdent(vs.Values[0], info) {
						return true
					}
					initILE, okI := extractInitILE(vs.Values[0])
					if !okI || !isUnshadowedUniverseType("int32", initILE.Indices[1].Pos(), ctx.Pkg) {
						c.disqualified = true
						return true
					}
					c.indexValExprs = append(c.indexValExprs, initILE)
				}
			case *ast.AssignStmt:
				if s.Tok != token.DEFINE || len(s.Lhs) != 1 || len(s.Rhs) != 1 || enclosingFuncBody(s, parents) == nil {
					return true
				}
				id, ok := s.Lhs[0].(*ast.Ident)
				if !ok || id.Name == "_" {
					return true
				}
				v := varFromObj(info.Defs[id])
				if v == nil || v.IsField() || cands[v] != nil {
					return true
				}
				isPtr, okM := isJ2gortOrderedMapBoxIntType(v.Type())
				if !okM {
					return true
				}
				initILE, okI := extractInitILE(s.Rhs[0])
				if !okI {
					return true
				}
				addCand(v, false, isPtr, file, initILE)
			}
			return true
		})
	}

	if len(candList) == 0 {
		return 0
	}

	// Track local variables initialized via `var v *j2gort.BoxInt = m.Get(k)` or `v := m.Get(k)`.
	getVars := make(map[*types.Var]*unboxIntMapVarSpec)

	validateLocalGetDecl := func(mCand *unboxIntMapCand, getCall *ast.CallExpr, file *ast.File, parents map[ast.Node]ast.Node) bool {
		cur := unwrapParensUp(getCall, parents)
		var vId *ast.Ident
		var declStmt ast.Stmt
		switch p := parents[cur].(type) {
		case *ast.ValueSpec:
			if len(p.Names) != 1 || len(p.Values) != 1 || p.Values[0] != cur || p.Names[0].Name == "_" {
				return false
			}
			gd, okG := parents[p].(*ast.GenDecl)
			if !okG || gd.Tok != token.VAR || len(gd.Specs) != 1 {
				return false
			}
			ds, okD := parents[gd].(*ast.DeclStmt)
			if !okD {
				return false
			}
			if _, okBlk := parents[ds].(*ast.BlockStmt); !okBlk {
				return false
			}
			vId = p.Names[0]
			declStmt = ds
		case *ast.AssignStmt:
			if p.Tok != token.DEFINE || len(p.Lhs) != 1 || len(p.Rhs) != 1 || p.Rhs[0] != cur {
				return false
			}
			if _, okBlk := parents[p].(*ast.BlockStmt); !okBlk {
				return false
			}
			id, okI := p.Lhs[0].(*ast.Ident)
			if !okI || id.Name == "_" {
				return false
			}
			vId = id
			declStmt = p
		default:
			return false
		}
		vObj := varFromObj(info.Defs[vId])
		if vObj == nil || !isPtrToJ2gortBoxInt(vObj.Type()) {
			return false
		}
		fnDecl := enclosingFuncDeclNode(declStmt, parents)
		if fnDecl == nil || fnDecl.Body == nil {
			return false
		}
		okName := vId.Name + "Ok_"
		if slabASTContainsIdent(fnDecl, okName) || !isUnshadowedUniverseType("int32", declStmt.Pos(), ctx.Pkg) || !isUnshadowedUniverseType("true", declStmt.Pos(), ctx.Pkg) {
			return false
		}
		spec := &unboxIntMapVarSpec{
			v:        vObj,
			cand:     mCand,
			fn:       fnDecl,
			file:     file,
			declStmt: declStmt,
			okName:   okName,
		}
		getVars[vObj] = spec
		mCand.localGetVars = append(mCand.localGetVars, spec)
		mCand.hasActiveOp = true
		return true
	}

	// Pass 2: Validate all definitions and uses of every candidate map across ctx.Files.
	for _, file := range ctx.Files {
		parents := parentsByFile[file]
		ast.Inspect(file, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if c := cands[varFromObj(info.Defs[id])]; c != nil {
				c.defCount++
				c.touchedFiles[file] = true
				if !ctx.ShouldTransform(file) || ctx.ModifiedFiles[file] {
					c.disqualified = true
				}
			}
			c := cands[varFromObj(info.Uses[id])]
			if c == nil || c.disqualified {
				return true
			}
			c.touchedFiles[file] = true
			if !ctx.ShouldTransform(file) || ctx.ModifiedFiles[file] {
				c.disqualified = true
				return true
			}

			var accessExpr ast.Expr = id
			if c.isField {
				if kv, okKV := parents[id].(*ast.KeyValueExpr); okKV && kv.Key == id {
					if c.isPtr && isNilIdent(kv.Value, info) {
						return true
					}
					initILE, okI := extractInitILE(kv.Value)
					if !okI || !isUnshadowedUniverseType("int32", initILE.Indices[1].Pos(), ctx.Pkg) {
						c.disqualified = true
						return true
					}
					c.indexValExprs = append(c.indexValExprs, initILE)
					return true
				}
				sel, okS := parents[id].(*ast.SelectorExpr)
				if !okS || sel.Sel != id || !isIndirectStructFieldSel(sel, info) {
					c.disqualified = true
					return true
				}
				accessExpr = sel
			}

			// Blank discard `_ = m`.
			if as, okAs := parents[accessExpr].(*ast.AssignStmt); okAs && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 {
				if lhsId, okL := stripParens(as.Lhs[0]).(*ast.Ident); okL && lhsId.Name == "_" && stripParens(as.Rhs[0]) == accessExpr {
					return true
				}
				if stripParens(as.Lhs[0]) == accessExpr {
					if c.isPtr && isNilIdent(as.Rhs[0], info) {
						return true
					}
					initILE, okI := extractInitILE(as.Rhs[0])
					if !okI || !isUnshadowedUniverseType("int32", initILE.Indices[1].Pos(), ctx.Pkg) {
						c.disqualified = true
						return true
					}
					c.indexValExprs = append(c.indexValExprs, initILE)
					return true
				}
			}

			pTop := unwrapParensUp(accessExpr, parents)

			// Nil check `m == nil` / `m != nil` on pointer map.
			if c.isPtr {
				if bin, okB := parents[pTop].(*ast.BinaryExpr); okB && (bin.Op == token.EQL || bin.Op == token.NEQ) {
					other := bin.Y
					if bin.Y == pTop {
						other = bin.X
					}
					if isNilIdent(other, info) {
						return true
					}
				}
			}

			// Release helper `j2gort.ReleaseOrderedMap(&m)` / `j2gort.ReleaseOrderedMap(m)`.
			addrTop := pTop
			if !c.isPtr {
				if un, okU := parents[pTop].(*ast.UnaryExpr); okU && un.Op == token.AND && un.X == pTop {
					addrTop = unwrapParensUp(un, parents)
				}
			}
			if relCall, okC := parents[addrTop].(*ast.CallExpr); okC && len(relCall.Args) == 1 && !relCall.Ellipsis.IsValid() && relCall.Args[0] == addrTop {
				if fnName, _, _, isJ2 := isJ2gortCall(relCall, info); isJ2 {
					switch fnName {
					case "ReleaseOrderedMap", "ReleaseCustomOrderedMap", "ReleaseIdentityOrderedMap":
						if ile, okILE := stripParens(relCall.Fun).(*ast.IndexListExpr); okILE && len(ile.Indices) == 2 {
							if !isUnshadowedUniverseType("int32", ile.Indices[1].Pos(), ctx.Pkg) {
								c.disqualified = true
								return true
							}
							c.indexValExprs = append(c.indexValExprs, ile)
						}
						return true
					}
				}
			}

			// Method call on `m`.
			mSel, okSel := parents[pTop].(*ast.SelectorExpr)
			if !okSel || mSel.X != pTop || mSel.Sel == nil {
				c.disqualified = true
				return true
			}
			mCall, okCall := parents[mSel].(*ast.CallExpr)
			if !okCall || mCall.Fun != mSel || mCall.Ellipsis.IsValid() {
				c.disqualified = true
				return true
			}
			_, inExprStmt := parents[mCall].(*ast.ExprStmt)

			switch mSel.Sel.Name {
			case "Size", "IsEmpty", "Clear", "ReleaseToPool", "IterKeys", "LiveKeySet", "KeysToImmutableSet", "String":
				if len(mCall.Args) != 0 {
					c.disqualified = true
				}
			case "ContainsKey":
				if len(mCall.Args) != 1 {
					c.disqualified = true
				}
			case "Remove":
				if !inExprStmt || len(mCall.Args) != 1 {
					c.disqualified = true
				}
			case "Put", "PutVoid", "PutIfAbsent":
				if !inExprStmt || len(mCall.Args) != 2 {
					c.disqualified = true
					return true
				}
				if _, okV := unwrapValueOfIntCall(mCall.Args[1], info); okV {
					c.putCalls = append(c.putCalls, mCall)
					c.hasActiveOp = true
					return true
				}
				if argId, okI := stripParens(mCall.Args[1]).(*ast.Ident); okI {
					if gv := getVars[varFromObj(info.Uses[argId])]; gv != nil && gv.cand == c {
						gv.putCalls = append(gv.putCalls, mCall)
						c.hasActiveOp = true
						return true
					}
				}
				c.disqualified = true
			case "Get":
				if len(mCall.Args) != 1 {
					c.disqualified = true
					return true
				}
				if unbCall := enclosingUnboxIntOfExpr(mCall, parents, info); unbCall != nil {
					if !isUnshadowedUniverseType("int32", unbCall.Pos(), ctx.Pkg) {
						c.disqualified = true
						return true
					}
					c.directGetUnb = append(c.directGetUnb, &unboxDirectGetOp{
						file:      file,
						outerCall: unbCall,
						getCall:   mCall,
					})
					c.hasActiveOp = true
					return true
				}
				if !validateLocalGetDecl(c, mCall, file, parents) {
					c.disqualified = true
				}
			default:
				c.disqualified = true
			}
			return true
		})
	}

	// Pass 3: Validate every use of every tracked `getResultVar` (`var ptr *BoxInt = m.Get(k)`).
	for _, c := range candList {
		if c.disqualified || c.defCount != 1 || !c.hasActiveOp {
			c.disqualified = true
			continue
		}
		for f := range c.touchedFiles {
			if !ctx.ShouldTransform(f) || ctx.ModifiedFiles[f] {
				c.disqualified = true
				break
			}
		}
		if c.disqualified {
			continue
		}
		for _, gv := range c.localGetVars {
			parents := parentsByFile[gv.file]
			defSeen := 0
			ast.Inspect(gv.fn.Body, func(n ast.Node) bool {
				if c.disqualified || n == nil {
					return false
				}
				if _, isFL := n.(*ast.FuncLit); isFL {
					// Disqualify if gv.v is captured inside any nested closure.
					ast.Inspect(n, func(cn ast.Node) bool {
						if cid, okI := cn.(*ast.Ident); okI && varFromObj(info.Uses[cid]) == gv.v {
							c.disqualified = true
							return false
						}
						return !c.disqualified
					})
					return false
				}
				id, ok := n.(*ast.Ident)
				if !ok {
					return true
				}
				if varFromObj(info.Defs[id]) == gv.v {
					defSeen++
					return true
				}
				if varFromObj(info.Uses[id]) != gv.v {
					return true
				}

				pTop := unwrapParensUp(id, parents)

				// 1. `_ = v` blank discard or `v = j2gort.ValueOfInt(expr)` reassignment.
				if as, okAs := parents[pTop].(*ast.AssignStmt); okAs && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 {
					if lhsId, okL := stripParens(as.Lhs[0]).(*ast.Ident); okL && lhsId.Name == "_" && as.Rhs[0] == pTop {
						return true
					}
					if stripParens(as.Lhs[0]) == id {
						if _, okBlk := parents[as].(*ast.BlockStmt); !okBlk {
							c.disqualified = true
							return false
						}
						if _, okV := unwrapValueOfIntCall(as.Rhs[0], info); okV {
							gv.hasReassign = true
							gv.reassigns = append(gv.reassigns, as)
							return true
						}
						c.disqualified = true
						return false
					}
				}

				// 2. `v == nil` / `v != nil`.
				if bin, okB := parents[pTop].(*ast.BinaryExpr); okB && (bin.Op == token.EQL || bin.Op == token.NEQ) {
					other := bin.Y
					if bin.Y == pTop {
						other = bin.X
					}
					if isNilIdent(other, info) {
						gv.nilBins = append(gv.nilBins, bin)
						gv.nilIsEq = append(gv.nilIsEq, bin.Op == token.EQL)
						return true
					}
				}

				// 3. `j2gort.UnboxInt([CheckNotNullPtr](v))`.
				if unbCall := enclosingUnboxIntOfExpr(id, parents, info); unbCall != nil {
					gv.unboxCalls = append(gv.unboxCalls, unbCall)
					return true
				}

				// 4. `j2gort.PtrToAny(v)` (not directly inside UnboxInt).
				if ptrCall, okP := parents[pTop].(*ast.CallExpr); okP && len(ptrCall.Args) == 1 && !ptrCall.Ellipsis.IsValid() && ptrCall.Args[0] == pTop {
					if fnName, _, _, isJ2 := isJ2gortCall(ptrCall, info); isJ2 && fnName == "PtrToAny" {
						if !isUnshadowedUniverseType("any", ptrCall.Pos(), ctx.Pkg) || !isUnshadowedUniverseType("nil", ptrCall.Pos(), ctx.Pkg) {
							c.disqualified = true
							return false
						}
						gv.ptrToAnyCall = append(gv.ptrToAnyCall, ptrCall)
						return true
					}
				}

				// 5. `m.Put(k, v)` when preceded in the same BlockStmt by `v = j2gort.ValueOfInt(...)`.
				if putCall, okP := parents[pTop].(*ast.CallExpr); okP && slices.Contains(gv.putCalls, putCall) && putCall.Args[1] == pTop {
					es, okES := parents[putCall].(*ast.ExprStmt)
					blk, okBlk := parents[es].(*ast.BlockStmt)
					if okES && okBlk {
						idx := slices.Index(blk.List, ast.Stmt(es))
						if idx > 0 {
							if prevAs, okPrev := blk.List[idx-1].(*ast.AssignStmt); okPrev && slices.Contains(gv.reassigns, prevAs) {
								return true
							}
						}
					}
				}

				c.disqualified = true
				return false
			})
			if defSeen != 1 {
				c.disqualified = true
			}
		}
	}

	// Pass 4: Apply rewrites cleanly on all eligible candidates.
	count := 0
	for _, c := range candList {
		if c.disqualified {
			continue
		}
		for _, ile := range c.indexValExprs {
			ile.Indices[1] = &ast.Ident{NamePos: ile.Indices[1].Pos(), Name: "int32"}
		}
		for _, putCall := range c.putCalls {
			if inner, ok := unwrapValueOfIntCall(putCall.Args[1], info); ok {
				putCall.Args[1] = inner
				count++
			}
		}
		for _, dg := range c.directGetUnb {
			parents := parentsByFile[dg.file]
			j2Qual, okQ := j2gortQualifierForFile(dg.file, ctx)
			if !okQ {
				continue
			}
			getSel := stripParens(dg.getCall.Fun).(*ast.SelectorExpr)
			getSel.Sel = &ast.Ident{NamePos: getSel.Sel.Pos(), Name: "GetOk"}
			pos := dg.outerCall.Pos()
			iife := &ast.CallExpr{
				Fun: &ast.FuncLit{
					Type: &ast.FuncType{
						Func:    pos,
						Params:  &ast.FieldList{},
						Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{NamePos: pos, Name: "int32"}}}},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.IfStmt{
								Init: &ast.AssignStmt{
									Lhs: []ast.Expr{&ast.Ident{Name: "v_"}, &ast.Ident{Name: "ok_"}},
									Tok: token.DEFINE,
									Rhs: []ast.Expr{dg.getCall},
								},
								Cond: &ast.Ident{Name: "ok_"},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "v_"}}}},
								},
							},
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: &ast.Ident{Name: j2Qual}, Sel: &ast.Ident{Name: "ThrowNPE"}},
								},
							},
							&ast.ReturnStmt{Results: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}}},
						},
					},
				},
			}
			if parent := parents[dg.outerCall]; parent != nil && replaceExprInParent(parent, dg.outerCall, iife, parents) {
				count++
			}
		}

		for _, gv := range c.localGetVars {
			parents := parentsByFile[gv.file]
			j2Qual, okQ := j2gortQualifierForFile(gv.file, ctx)
			if !okQ {
				continue
			}
			vName := gv.v.Name()

			// Rewrite `var v *BoxInt = m.Get(k)` / `v := m.Get(k)` -> `v, vOk_ := m.GetOk(k); _ = vOk_`.
			var getCall *ast.CallExpr
			switch ds := gv.declStmt.(type) {
			case *ast.DeclStmt:
				vs := ds.Decl.(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				getCall = stripParens(vs.Values[0]).(*ast.CallExpr)
			case *ast.AssignStmt:
				getCall = stripParens(ds.Rhs[0]).(*ast.CallExpr)
			}
			getSel := stripParens(getCall.Fun).(*ast.SelectorExpr)
			getSel.Sel = &ast.Ident{NamePos: getSel.Sel.Pos(), Name: "GetOk"}
			newDef := &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{NamePos: gv.declStmt.Pos(), Name: vName},
					&ast.Ident{NamePos: gv.declStmt.Pos(), Name: gv.okName},
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{getCall},
			}
			blankOk := &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "_"}},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{&ast.Ident{Name: gv.okName}},
			}
			blk := parents[gv.declStmt].(*ast.BlockStmt)
			dIdx := slices.Index(blk.List, gv.declStmt)
			blk.List = slices.Replace(blk.List, dIdx, dIdx+1, ast.Stmt(newDef), ast.Stmt(blankOk))
			parents[newDef] = blk
			parents[blankOk] = blk
			count++

			// Rewrite `v = j2gort.ValueOfInt(expr)` -> `v, vOk_ = expr, true`.
			for _, as := range gv.reassigns {
				inner, _ := unwrapValueOfIntCall(as.Rhs[0], info)
				as.Lhs = []ast.Expr{
					&ast.Ident{NamePos: as.Pos(), Name: vName},
					&ast.Ident{NamePos: as.Pos(), Name: gv.okName},
				}
				as.Rhs = []ast.Expr{
					inner,
					&ast.Ident{NamePos: as.Pos(), Name: "true"},
				}
				count++
			}

			// Rewrite `v == nil` -> `!vOk_` and `v != nil` -> `vOk_`.
			for i, bin := range gv.nilBins {
				var repl ast.Expr = &ast.Ident{NamePos: bin.Pos(), Name: gv.okName}
				if gv.nilIsEq[i] {
					repl = &ast.UnaryExpr{OpPos: bin.Pos(), Op: token.NOT, X: repl}
				}
				if parent := parents[bin]; parent != nil && replaceExprInParent(parent, bin, repl, parents) {
					count++
				}
			}

			// Rewrite `j2gort.UnboxInt([CheckNotNullPtr](v))`.
			for _, unbCall := range gv.unboxCalls {
				pos := unbCall.Pos()
				iife := &ast.CallExpr{
					Fun: &ast.FuncLit{
						Type: &ast.FuncType{
							Func:    pos,
							Params:  &ast.FieldList{},
							Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{NamePos: pos, Name: "int32"}}}},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.IfStmt{
									Cond: &ast.UnaryExpr{Op: token.NOT, X: &ast.Ident{Name: gv.okName}},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.ExprStmt{
												X: &ast.CallExpr{
													Fun: &ast.SelectorExpr{X: &ast.Ident{Name: j2Qual}, Sel: &ast.Ident{Name: "ThrowNPE"}},
												},
											},
										},
									},
								},
								&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: vName}}},
							},
						},
					},
				}
				if parent := parents[unbCall]; parent != nil && replaceExprInParent(parent, unbCall, iife, parents) {
					count++
				}
			}

			// Rewrite `j2gort.PtrToAny(v)` -> `func() any { if !vOk_ { return nil }; return j2gort.ValueOfInt(v) }()`.
			for _, pCall := range gv.ptrToAnyCall {
				pos := pCall.Pos()
				iife := &ast.CallExpr{
					Fun: &ast.FuncLit{
						Type: &ast.FuncType{
							Func:    pos,
							Params:  &ast.FieldList{},
							Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{NamePos: pos, Name: "any"}}}},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.IfStmt{
									Cond: &ast.UnaryExpr{Op: token.NOT, X: &ast.Ident{Name: gv.okName}},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "nil"}}}},
									},
								},
								&ast.ReturnStmt{
									Results: []ast.Expr{
										&ast.CallExpr{
											Fun:  &ast.SelectorExpr{X: &ast.Ident{Name: j2Qual}, Sel: &ast.Ident{Name: "ValueOfInt"}},
											Args: []ast.Expr{&ast.Ident{Name: vName}},
										},
									},
								},
							},
						},
					},
				}
				if parent := parents[pCall]; parent != nil && replaceExprInParent(parent, pCall, iife, parents) {
					count++
				}
			}
		}

		for f := range c.touchedFiles {
			ctx.MarkFileModified(f)
		}
		count++
	}

	return count
}

func isSlabPooledOrderedMapField(owner *types.Named, fieldName string) bool {
	if owner == nil || owner.Obj() == nil {
		return false
	}
	switch owner.Obj().Name() {
	case "LinkedDirectedGraphBase":
		return fieldName == "nodes"
	case "Scope":
		return fieldName == "vars"
	case "GlobalNamespace":
		return fieldName == "nameMap" || fieldName == "nameMapByModule"
	case "GlobalNamespace_ModuleNameMap":
		return fieldName == "names"
	case "GlobalNamespace_Name":
		return fieldName == "refsForNodeMap"
	case "InlineFunctions":
		return fieldName == "fns" || fieldName == "anonFns"
	case "InlineFunctions_FunctionState":
		return fieldName == "references"
	}
	return false
}

func isLinkedDeclaredCollectionType(t types.Type) bool {
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
	case "Ext_java_util_LinkedHashMap", "Ext_java_util_SequencedMap",
		"Ext_java_util_LinkedHashSet", "Ext_java_util_SequencedSet":
		return true
	}
	return false
}

func isJ2gortArrayListOrListValType(t types.Type) bool {
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

func isUnorderedHashMapCtorCall(expr ast.Expr, info *types.Info) bool {
	expr = unwrapMapAdaptersDown(expr, info)
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return false
	}
	fnName, _, _, isJ2 := isJ2gortCall(call, info)
	if !isJ2 {
		return false
	}
	switch fnName {
	case "NewExt_java_util_HashMap", "Ext_com_google_common_collect_Maps_NewHashMap_0":
		return len(call.Args) == 0
	case "NewExt_java_util_HashMap_1_int32", "Ext_com_google_common_collect_Maps_NewHashMapWithExpectedSize":
		return len(call.Args) == 1 && isIntegerExpr(call.Args[0], info)
	}
	return false
}

func resolvePlainCtorSig(j2Pkg *types.Package, fnName string) (sig *types.Signature, returnsPtr bool, ok bool) {
	if j2Pkg == nil || j2Pkg.Scope() == nil {
		return nil, false, false
	}
	fn, okFn := j2Pkg.Scope().Lookup(fnName).(*types.Func)
	if !okFn || fn == nil || !fn.Exported() {
		return nil, false, false
	}
	s, okS := fn.Type().(*types.Signature)
	if !okS || s.Results() == nil || s.Results().Len() != 1 {
		return nil, false, false
	}
	_, isPtr := types.Unalias(s.Results().At(0).Type()).(*types.Pointer)
	return s, isPtr, true
}

func canMakePlainCtor(j2Pkg *types.Package, baseName string, wantPtr bool, hasCap bool, capExpr ast.Expr, ctx *PassContext) bool {
	if j2Pkg == nil {
		return false
	}
	checkCapConv := func(sig *types.Signature) bool {
		if sig.Params().Len() != 1 {
			return false
		}
		if capExpr == nil {
			return true
		}
		paramT := types.Unalias(sig.Params().At(0).Type())
		if types.Identical(paramT, types.Typ[types.Int]) {
			if _, okLit := unwrapInt32BasicLit(capExpr, ctx.Info); okLit {
				return true
			}
			if tv, okT := ctx.Info.Types[capExpr]; okT && tv.Type != nil && types.AssignableTo(tv.Type, paramT) {
				return true
			}
			return isUnshadowedUniverseType("int", capExpr.Pos(), ctx.Pkg)
		}
		return true
	}

	if wantPtr {
		if !hasCap {
			return true
		}
		if sig, isPtr, ok := resolvePlainCtorSig(j2Pkg, "New"+baseName+"PtrCap"); ok && isPtr && checkCapConv(sig) {
			return true
		}
		if sig, isPtr, ok := resolvePlainCtorSig(j2Pkg, "New"+baseName); ok && isPtr && checkCapConv(sig) {
			return true
		}
		return false
	}

	if !hasCap {
		if sig, isPtr, ok := resolvePlainCtorSig(j2Pkg, "New"+baseName); ok && !isPtr && (sig.Params().Len() == 0 || sig.Params().Len() == 1) {
			return true
		}
		if sig, isPtr, ok := resolvePlainCtorSig(j2Pkg, "New"+baseName+"Val"); ok && !isPtr && sig.Params().Len() == 0 {
			return true
		}
		return false
	}
	if sig, isPtr, ok := resolvePlainCtorSig(j2Pkg, "New"+baseName); ok && !isPtr && checkCapConv(sig) {
		return true
	}
	if sig, isPtr, ok := resolvePlainCtorSig(j2Pkg, "New"+baseName+"Cap"); ok && !isPtr && checkCapConv(sig) {
		return true
	}
	return false
}

func unwrapInt32BasicLit(expr ast.Expr, info *types.Info) (*ast.BasicLit, bool) {
	e := stripParens(expr)
	if lit, ok := e.(*ast.BasicLit); ok && lit.Kind == token.INT {
		return lit, true
	}
	call, ok := e.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 || call.Ellipsis.IsValid() {
		return nil, false
	}
	if tv, okT := info.Types[call.Fun]; okT && tv.IsType() && types.Identical(types.Unalias(tv.Type), types.Typ[types.Int32]) {
		if lit, okL := stripParens(call.Args[0]).(*ast.BasicLit); okL && lit.Kind == token.INT {
			return lit, true
		}
	}
	return nil, false
}

func convertPlainCapArg(capExpr ast.Expr, sig *types.Signature, info *types.Info, pos token.Pos) ast.Expr {
	if capExpr == nil {
		return &ast.BasicLit{ValuePos: pos, Kind: token.INT, Value: "0"}
	}
	if sig != nil && sig.Params().Len() == 1 && types.Identical(types.Unalias(sig.Params().At(0).Type()), types.Typ[types.Int]) {
		if lit, okLit := unwrapInt32BasicLit(capExpr, info); okLit {
			return &ast.BasicLit{ValuePos: lit.ValuePos, Kind: token.INT, Value: lit.Value}
		}
		if info != nil {
			if tv, okT := info.Types[capExpr]; okT && tv.Type != nil && types.AssignableTo(tv.Type, types.Typ[types.Int]) {
				return capExpr
			}
		}
		return &ast.CallExpr{
			Fun:  &ast.Ident{NamePos: pos, Name: "int"},
			Args: []ast.Expr{capExpr},
		}
	}
	return capExpr
}

func makePlainMapInitExpr(
	ctx *PassContext,
	j2Sel func(string, token.Pos) *ast.SelectorExpr,
	makeValueMapType func(token.Pos) ast.Expr,
	makeKeyType func() ast.Expr,
	makeValType func() ast.Expr,
	wantPtr bool,
	hasCap bool,
	capExpr ast.Expr,
	pos token.Pos,
	inHeader bool,
) ast.Expr {
	j2Pkg := findJ2gortTypesPkg(ctx)
	indices := []ast.Expr{makeKeyType(), makeValType()}
	if wantPtr {
		if !hasCap {
			if sig, isPtr, ok := resolvePlainCtorSig(j2Pkg, "NewPlainMapPtr"); ok && isPtr && sig.Params().Len() == 0 {
				return &ast.CallExpr{Fun: &ast.IndexListExpr{X: j2Sel("NewPlainMapPtr", pos), Indices: indices}}
			}
			lit := ast.Expr(&ast.UnaryExpr{
				OpPos: pos,
				Op:    token.AND,
				X:     &ast.CompositeLit{Type: makeValueMapType(pos)},
			})
			if inHeader {
				lit = &ast.ParenExpr{X: lit}
			}
			return lit
		}
		if sig, isPtr, ok := resolvePlainCtorSig(j2Pkg, "NewPlainMapPtrCap"); ok && isPtr && sig.Params().Len() == 1 {
			return &ast.CallExpr{
				Fun:  &ast.IndexListExpr{X: j2Sel("NewPlainMapPtrCap", pos), Indices: indices},
				Args: []ast.Expr{convertPlainCapArg(capExpr, sig, ctx.Info, pos)},
			}
		}
		sig, _, _ := resolvePlainCtorSig(j2Pkg, "NewPlainMap")
		return &ast.CallExpr{
			Fun:  &ast.IndexListExpr{X: j2Sel("NewPlainMap", pos), Indices: indices},
			Args: []ast.Expr{convertPlainCapArg(capExpr, sig, ctx.Info, pos)},
		}
	}

	if !hasCap {
		if sig, isPtr, ok := resolvePlainCtorSig(j2Pkg, "NewPlainMap"); ok && !isPtr {
			var args []ast.Expr
			if sig.Params().Len() == 1 {
				args = []ast.Expr{&ast.BasicLit{ValuePos: pos, Kind: token.INT, Value: "0"}}
			}
			return &ast.CallExpr{Fun: &ast.IndexListExpr{X: j2Sel("NewPlainMap", pos), Indices: indices}, Args: args}
		}
		return &ast.CallExpr{Fun: &ast.IndexListExpr{X: j2Sel("NewPlainMapVal", pos), Indices: indices}}
	}

	if sig, isPtr, ok := resolvePlainCtorSig(j2Pkg, "NewPlainMap"); ok && !isPtr && sig.Params().Len() == 1 {
		return &ast.CallExpr{
			Fun:  &ast.IndexListExpr{X: j2Sel("NewPlainMap", pos), Indices: indices},
			Args: []ast.Expr{convertPlainCapArg(capExpr, sig, ctx.Info, pos)},
		}
	}
	sig, _, _ := resolvePlainCtorSig(j2Pkg, "NewPlainMapCap")
	return &ast.CallExpr{
		Fun:  &ast.IndexListExpr{X: j2Sel("NewPlainMapCap", pos), Indices: indices},
		Args: []ast.Expr{convertPlainCapArg(capExpr, sig, ctx.Info, pos)},
	}
}

func resolvePlainMapComputeHelper(j2Pkg *types.Package, wantGoFunc bool, isStrVal bool) (helperName string, takesNullVal bool, ok bool) {
	if j2Pkg == nil || j2Pkg.Scope() == nil {
		return "", false, false
	}
	var candidates []string
	if wantGoFunc {
		candidates = []string{"PlainMapComputeIfAbsent", "PlainMapComputeIfAbsentFn"}
	} else {
		candidates = []string{"PlainMapComputeIfAbsentExtFn", "PlainMapComputeIfAbsentFn"}
	}
	for _, name := range candidates {
		fn, okFn := j2Pkg.Scope().Lookup(name).(*types.Func)
		if !okFn || fn == nil || !fn.Exported() {
			continue
		}
		sig, okS := fn.Type().(*types.Signature)
		if !okS || (sig.Params().Len() != 3 && sig.Params().Len() != 4) {
			continue
		}
		lastT := types.Unalias(sig.Params().At(sig.Params().Len() - 1).Type())
		_, lastIsGoFunc := lastT.(*types.Signature)
		if lastIsGoFunc != wantGoFunc {
			continue
		}
		takesNull := sig.Params().Len() == 4
		if isStrVal && !takesNull {
			// Without an explicit NullString sentinel parameter, Go's zero string "" would be
			// conflated with Java null.
			continue
		}
		return name, takesNull, true
	}
	return "", false, false
}

func canPromoteGroupToPlainMap(members []*mapCandidate, parentsByFile map[*ast.File]map[ast.Node]ast.Node, ctx *PassContext) bool {
	if len(members) == 0 || ctx == nil {
		return false
	}
	j2Pkg := findJ2gortTypesPkg(ctx)
	if j2Pkg == nil || j2Pkg.Scope() == nil || j2Pkg.Scope().Lookup("PlainMap") == nil || j2Pkg.Scope().Lookup("NewPlainMap") == nil {
		return false
	}
	inGroup := make(map[*mapCandidate]bool, len(members))
	for _, m := range members {
		if m == nil {
			return false
		}
		inGroup[m] = true
	}
	sawAnyUnorderedAlloc := false
	for _, m := range members {
		ok, sawAlloc := canPromoteMemberToPlainMap(m, len(members) > 1, inGroup, j2Pkg, parentsByFile, ctx)
		if !ok {
			return false
		}
		if sawAlloc {
			sawAnyUnorderedAlloc = true
		}
	}
	return sawAnyUnorderedAlloc
}

func canPromoteMemberToPlainMap(
	m *mapCandidate,
	multiMember bool,
	inGroup map[*mapCandidate]bool,
	j2Pkg *types.Package,
	parentsByFile map[*ast.File]map[ast.Node]ast.Node,
	ctx *PassContext,
) (ok bool, sawUnorderedAlloc bool) {
	_ = parentsByFile
	if m == nil || ctx == nil || m.v == nil || m.isParam || m.isCastAccessor || m.isPtrAlias || m.isIdentityMap {
		return false, false
	}
	if multiMember {
		if !m.isPtr() {
			return false, false
		}
	} else if m.isSharedPtrLocal {
		return false, false
	}
	if m.isField && isSlabPooledOrderedMapField(m.ownerStruct, m.v.Name()) {
		return false, false
	}
	if m.isField && isCanonicalVarIdentityMapField(m.ownerStruct, m.v.Name(), m.keyType, m.valType) {
		return false, false
	}
	if isLinkedDeclaredCollectionType(m.v.Type()) || isJ2gortArrayListOrListValType(m.valType) {
		return false, false
	}
	if m.hasNilCheck && !m.isPtr() && !m.isLazyValueMap && !j2gortTypeHasMethod(ctx, "PlainMap", "IsNil") {
		return false, false
	}

	isStrVal := types.Identical(m.valType, types.Typ[types.String])

	checkInitBranch := func(initK mapInitKind, initExpr ast.Expr, capArg ast.Expr, otherCand *mapCandidate) bool {
		switch initK {
		case mapInitNil:
			return true
		case mapInitEmpty:
			if !isUnorderedHashMapCtorCall(initExpr, ctx.Info) || !canMakePlainCtor(j2Pkg, "PlainMap", m.isPtr(), false, nil, ctx) {
				return false
			}
			sawUnorderedAlloc = true
			return true
		case mapInitCap:
			if !isUnorderedHashMapCtorCall(initExpr, ctx.Info) || !canMakePlainCtor(j2Pkg, "PlainMap", m.isPtr(), true, capArg, ctx) {
				return false
			}
			sawUnorderedAlloc = true
			return true
		case mapInitSharedPtr:
			return multiMember && otherCand != nil && inGroup[otherCand] && m.isPtr() && otherCand.isPtr()
		default:
			return false
		}
	}

	for i := range m.ops {
		op := &m.ops[i]
		switch op.kind {
		case mapOpDeclField:
			continue
		case mapOpDeclVarSpec:
			spec, ok := op.spanNode.(*ast.ValueSpec)
			if !ok {
				return false, false
			}
			if len(spec.Values) == 0 {
				continue
			}
			rhs := spec.Values[0]
			if op.initKind == mapInitIIFE && op.iifeSpec != nil {
				for _, br := range op.iifeSpec.branches {
					if !checkInitBranch(br.initKind, br.ret.Results[0], br.otherExpr, br.otherCand) {
						return false, false
					}
				}
			} else if !checkInitBranch(op.initKind, rhs, op.otherExpr, op.otherCand) {
				return false, false
			}
		case mapOpDeclShortVar, mapOpAssignInit:
			as, ok := op.spanNode.(*ast.AssignStmt)
			if !ok || len(as.Rhs) != 1 {
				return false, false
			}
			rhs := as.Rhs[0]
			if op.initKind == mapInitIIFE && op.iifeSpec != nil {
				for _, br := range op.iifeSpec.branches {
					if !checkInitBranch(br.initKind, br.ret.Results[0], br.otherExpr, br.otherCand) {
						return false, false
					}
				}
			} else if !checkInitBranch(op.initKind, rhs, op.otherExpr, op.otherCand) {
				return false, false
			}
		case mapOpNilCheck:
			if !m.isPtr() && !m.isLazyValueMap && !j2gortTypeHasMethod(ctx, "PlainMap", "IsNil") {
				return false, false
			}
		case mapOpMethodDirect:
			switch op.methodsSel {
			case "Size", "IsEmpty", "Clear":
			default:
				return false, false
			}
		case mapOpContainsKey, mapOpGetOrDefault, mapOpRawUnbox:
			continue
		case mapOpGet:
			if isStrVal {
				wantHelper := "PlainMapGetStr"
				if m.isPtr() {
					wantHelper = "PlainMapPtrGetStr"
				}
				if j2Pkg.Scope().Lookup(wantHelper) == nil {
					return false, false
				}
			}
		case mapOpPut:
			if isStrVal && (!op.discardRet || op.isPutIfAbsent) {
				wantHelper := "PlainMapPutStr"
				if op.isPutIfAbsent {
					wantHelper = "PlainMapPutIfAbsentStr"
				}
				if j2Pkg.Scope().Lookup(wantHelper) == nil {
					return false, false
				}
			}
		case mapOpRemove:
			if isStrVal && !op.discardRet {
				if j2Pkg.Scope().Lookup("PlainMapRemoveStr") == nil {
					return false, false
				}
			}
		case mapOpComputeIfAbsent:
			wantGoFunc := op.computeSpec != nil && op.computeSpec.fnLit != nil
			if _, _, okH := resolvePlainMapComputeHelper(j2Pkg, wantGoFunc, isStrVal); !okH {
				return false, false
			}
		default:
			return false, false
		}
	}

	return true, sawUnorderedAlloc
}

func isCanonicalVarIdentityMapField(owner *types.Named, fieldName string, keyT, valT types.Type) bool {
	if owner == nil || owner.Obj() == nil || keyT == nil || valT == nil {
		return false
	}
	kNamed := ptrNamedOrigin(keyT)
	vNamed := ptrNamedOrigin(valT)
	if kNamed == nil || kNamed.Obj() == nil || kNamed.Obj().Name() != "Var" || vNamed == nil || vNamed.Obj() == nil {
		return false
	}
	switch owner.Obj().Name() {
	case "ReferenceCollector":
		return fieldName == "referenceMap" && vNamed.Obj().Name() == "ReferenceCollection"
	case "InlineVariables_InliningBehavior":
		return fieldName == "currentScopeHandledVarAnalysesMap" && vNamed.Obj().Name() == "InlineVariables_InlineVarAnalysis"
	}
	return false
}

// applyCanonicalVarIdentityOrderedMapUpgrade upgrades already-lowered CustomOrderedMap[*Var, ...]
// fields on ReferenceCollector.referenceMap and InlineVariables_InliningBehavior.currentScopeHandledVarAnalysesMap
// to pointer-keyed OrderedMap[*Var, ...], and ensures ReferenceCollector.GetReferences safely resolves
// uncached implicit Vars (e.g. VariableReferenceCheck's scope.MakeImplicitVar(EXPORTS)) via their Scope.
func applyCanonicalVarIdentityOrderedMapUpgrade(ctx *PassContext) int {
	if ctx == nil || ctx.Info == nil || !j2gortHasExport(ctx, "OrderedMap") {
		return 0
	}
	info := ctx.Info

	type targetField struct {
		file      *ast.File
		owner     *types.Named
		fld       *ast.Field
		fldVar    *types.Var
		isCustom  bool
		isOrdered bool
		typeSel   *ast.SelectorExpr
	}
	var targets []*targetField
	var rcNamed, varNamed, scopeNamed *types.Named
	var rcFile *ast.File

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
				if named == nil || named.Obj() == nil {
					continue
				}
				switch named.Obj().Name() {
				case "ReferenceCollector":
					rcNamed = named
					rcFile = file
				case "Var":
					varNamed = named
				case "Scope":
					scopeNamed = named
				}
				if named.Obj().Name() != "ReferenceCollector" && named.Obj().Name() != "InlineVariables_InliningBehavior" {
					continue
				}
				for _, fld := range st.Fields.List {
					if len(fld.Names) != 1 {
						continue
					}
					v := varFromObj(info.Defs[fld.Names[0]])
					if v == nil || !v.IsField() {
						continue
					}
					ft := types.Unalias(v.Type())
					if ptr, isP := ft.(*types.Pointer); isP {
						ft = types.Unalias(ptr.Elem())
					}
					var typeSel *ast.SelectorExpr
					tExpr := stripParens(fld.Type)
					if stE, isSt := tExpr.(*ast.StarExpr); isSt {
						tExpr = stripParens(stE.X)
					}
					if ile, isIL := tExpr.(*ast.IndexListExpr); isIL {
						if sel, isSel := stripParens(ile.X).(*ast.SelectorExpr); isSel {
							typeSel = sel
						}
					}
					// Also detect ReferenceCollector.referenceMap freshly lowered to OrderedMap in the same pass run
					// (before ctx.Info is reloaded) so GetReferences implicit-var fallback lands on Run 1 (idempotent).
					if named.Obj().Name() == "ReferenceCollector" && v.Name() == "referenceMap" &&
						typeSel != nil && typeSel.Sel != nil && typeSel.Sel.Name == "OrderedMap" {
						targets = append(targets, &targetField{
							file:      file,
							owner:     named,
							fld:       fld,
							fldVar:    v,
							isOrdered: true,
							typeSel:   typeSel,
						})
						continue
					}
					fNamed, okN := ft.(*types.Named)
					if !okN || fNamed.Obj() == nil || !isJ2gortPkg(fNamed.Obj().Pkg()) {
						continue
					}
					mapTypeName := fNamed.Obj().Name()
					if mapTypeName != "CustomOrderedMap" && mapTypeName != "OrderedMap" {
						continue
					}
					targs := fNamed.TypeArgs()
					if targs == nil || targs.Len() != 2 {
						continue
					}
					if !isCanonicalVarIdentityMapField(named, v.Name(), targs.At(0), targs.At(1)) {
						continue
					}
					targets = append(targets, &targetField{
						file:      file,
						owner:     named,
						fld:       fld,
						fldVar:    v,
						isCustom:  mapTypeName == "CustomOrderedMap",
						isOrdered: mapTypeName == "OrderedMap",
						typeSel:   typeSel,
					})
				}
			}
		}
	}

	count := 0
	var rcHasOrderedMap bool
	for _, tf := range targets {
		if tf.isOrdered {
			if tf.owner == rcNamed {
				rcHasOrderedMap = true
			}
			continue
		}
		if !tf.isCustom || tf.typeSel == nil || !ctx.ShouldTransform(tf.file) {
			continue
		}
		// Scan every reference to tf.fldVar across all files to verify all required OrderedMap symbols exist.
		type fnRewrite struct {
			file    *ast.File
			sel     *ast.SelectorExpr
			newName string
		}
		var fnRewrites []fnRewrite
		canUpgrade := true
		needMethods := make(map[string]bool)
		needFuncs := make(map[string]bool)

		for _, file := range ctx.Files {
			parents, _ := buildParentMap(file)
			ast.Inspect(file, func(n ast.Node) bool {
				if !canUpgrade || n == nil {
					return false
				}
				id, ok := n.(*ast.Ident)
				if !ok || varFromObj(info.Uses[id]) != tf.fldVar {
					return true
				}
				if !ctx.ShouldTransform(file) {
					canUpgrade = false
					return false
				}
				sel, okS := parents[id].(*ast.SelectorExpr)
				if !okS || sel.Sel != id {
					canUpgrade = false
					return false
				}
				// Case 1: Assignment `this_.field = j2gort.<New|Acquire>CustomOrderedMap*[...](...)`
				if as, okA := parents[sel].(*ast.AssignStmt); okA && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && stripParens(as.Lhs[0]) == sel {
					if call, okC := stripParens(as.Rhs[0]).(*ast.CallExpr); okC {
						fun := stripParens(call.Fun)
						if ile, okIL := fun.(*ast.IndexListExpr); okIL {
							fun = stripParens(ile.X)
						}
						if cSel, okCS := fun.(*ast.SelectorExpr); okCS && cSel.Sel != nil {
							oldN := cSel.Sel.Name
							if strings.Contains(oldN, "CustomOrderedMap") {
								newN := strings.Replace(oldN, "CustomOrderedMap", "OrderedMap", 1)
								needFuncs[newN] = true
								fnRewrites = append(fnRewrites, fnRewrite{file: file, sel: cSel, newName: newN})
								return true
							}
						}
					}
					if isNilIdent(as.Rhs[0], info) {
						return true
					}
					canUpgrade = false
					return false
				}
				// Case 2: Method call `this_.field.<Method>(...)`
				if mSel, okM := parents[sel].(*ast.SelectorExpr); okM && mSel.X == sel && mSel.Sel != nil {
					needMethods[mSel.Sel.Name] = true
					return true
				}
				// Case 3: Address-of `&this_.field` passed to `j2gort.CustomOrderedMap*(...)`
				if un, okU := parents[sel].(*ast.UnaryExpr); okU && un.Op == token.AND && un.X == sel {
					if call, okC := parents[un].(*ast.CallExpr); okC && len(call.Args) >= 1 && call.Args[0] == un {
						fun := stripParens(call.Fun)
						if ile, okIL := fun.(*ast.IndexListExpr); okIL {
							fun = stripParens(ile.X)
						} else if ie, okI := fun.(*ast.IndexExpr); okI {
							fun = stripParens(ie.X)
						}
						if cSel, okCS := fun.(*ast.SelectorExpr); okCS && cSel.Sel != nil && strings.HasPrefix(cSel.Sel.Name, "CustomOrderedMap") {
							newN := strings.Replace(cSel.Sel.Name, "CustomOrderedMap", "OrderedMap", 1)
							needFuncs[newN] = true
							fnRewrites = append(fnRewrites, fnRewrite{file: file, sel: cSel, newName: newN})
							return true
						}
					}
					canUpgrade = false
					return false
				}
				// Case 4: Nil check `this_.field != nil` / `this_.field == nil`
				if bin, okB := parents[sel].(*ast.BinaryExpr); okB && (bin.Op == token.EQL || bin.Op == token.NEQ) {
					return true
				}
				canUpgrade = false
				return false
			})
		}
		if !canUpgrade {
			continue
		}
		for fnName := range needFuncs {
			if !j2gortHasExport(ctx, fnName) {
				canUpgrade = false
				break
			}
		}
		for mName := range needMethods {
			if !j2gortTypeHasMethod(ctx, "OrderedMap", mName) {
				canUpgrade = false
				break
			}
		}
		if !canUpgrade {
			continue
		}
		tf.typeSel.Sel = &ast.Ident{NamePos: tf.typeSel.Sel.Pos(), Name: "OrderedMap"}
		ctx.MarkFileModified(tf.file)
		count++
		for _, rw := range fnRewrites {
			rw.sel.Sel = &ast.Ident{NamePos: rw.sel.Sel.Pos(), Name: rw.newName}
			ctx.MarkFileModified(rw.file)
			count++
		}
		if tf.owner == rcNamed {
			rcHasOrderedMap = true
		}
	}

	// Ensure ReferenceCollector.GetReferences resolves uncached implicit Vars (index < 0, e.g. MakeImplicitVar(EXPORTS))
	// through their Scope's cached implicit slot when referenceMap is pointer-keyed OrderedMap[*Var, *ReferenceCollection].
	if rcHasOrderedMap && rcNamed != nil && rcFile != nil && varNamed != nil && scopeNamed != nil && ctx.ShouldTransform(rcFile) {
		count += ensureRefCollectorGetReferencesImplicitVarFallback(ctx, rcFile, rcNamed, varNamed, scopeNamed)
	}

	return count
}

func ensureRefCollectorGetReferencesImplicitVarFallback(
	ctx *PassContext,
	rcFile *ast.File,
	rcNamed, varNamed, scopeNamed *types.Named,
) int {
	// Verify Var has index (int32), scope (interface assignable from *Scope), name (string)
	// and Scope has implicitVars (any) and method GetOwnSlot(string) any.
	idxObj, _, _ := types.LookupFieldOrMethod(types.NewPointer(varNamed), true, ctx.Pkg, "index")
	scopeObj, _, _ := types.LookupFieldOrMethod(types.NewPointer(varNamed), true, ctx.Pkg, "scope")
	nameObj, _, _ := types.LookupFieldOrMethod(types.NewPointer(varNamed), true, ctx.Pkg, "name")
	impObj, _, _ := types.LookupFieldOrMethod(types.NewPointer(scopeNamed), true, ctx.Pkg, "implicitVars")
	slotObj, _, _ := types.LookupFieldOrMethod(types.NewPointer(scopeNamed), true, ctx.Pkg, "GetOwnSlot")
	idxV, okI := idxObj.(*types.Var)
	scopeV, okSc := scopeObj.(*types.Var)
	nameV, okN := nameObj.(*types.Var)
	impV, okImp := impObj.(*types.Var)
	slotFn, okSl := slotObj.(*types.Func)
	if !okI || !okSc || !okN || !okImp || !okSl ||
		!types.Identical(idxV.Type(), types.Typ[types.Int32]) ||
		!types.AssignableTo(types.NewPointer(scopeNamed), scopeV.Type()) ||
		!types.Identical(nameV.Type(), types.Typ[types.String]) ||
		impV.Type() == nil {
		return 0
	}
	slotSig, okSS := slotFn.Type().(*types.Signature)
	if !okSS || slotSig.Params().Len() != 1 || slotSig.Results().Len() != 1 ||
		!types.Identical(slotSig.Params().At(0).Type(), types.Typ[types.String]) {
		return 0
	}

	j2Qual, okQ := ensureJ2gortImportInFile(rcFile, ctx)
	if !okQ {
		return 0
	}

	for _, decl := range rcFile.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || fn.Name == nil || fn.Name.Name != "GetReferences" || fn.Body == nil ||
			len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 {
			continue
		}
		rId := fn.Recv.List[0].Names[0]
		rv := varFromObj(ctx.Info.Defs[rId])
		if rv == nil || ptrNamedOrigin(rv.Type()) != rcNamed || rId.Name == "" || rId.Name == "_" {
			continue
		}
		// Idempotency guard: already has implicitVars fallback.
		hasFallback := false
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if sel, okS := n.(*ast.SelectorExpr); okS && sel.Sel != nil && sel.Sel.Name == "implicitVars" {
				hasFallback = true
				return false
			}
			return true
		})
		if hasFallback || len(fn.Body.List) == 0 {
			return 0
		}
		// Match final `return j2gort.PtrToAny(this_.referenceMap.Get(v))`
		retIdx := len(fn.Body.List) - 1
		retStmt, okR := fn.Body.List[retIdx].(*ast.ReturnStmt)
		if !okR || len(retStmt.Results) != 1 {
			continue
		}
		ptrCall, okP := stripParens(retStmt.Results[0]).(*ast.CallExpr)
		if !okP || len(ptrCall.Args) != 1 {
			continue
		}
		ptrSel, okPS := stripParens(ptrCall.Fun).(*ast.SelectorExpr)
		if !okPS || ptrSel.Sel == nil || ptrSel.Sel.Name != "PtrToAny" {
			continue
		}
		getCall, okG := stripParens(ptrCall.Args[0]).(*ast.CallExpr)
		if !okG || len(getCall.Args) != 1 {
			continue
		}
		getSel, okGS := stripParens(getCall.Fun).(*ast.SelectorExpr)
		if !okGS || getSel.Sel == nil || getSel.Sel.Name != "Get" {
			continue
		}
		rmSel, okRM := stripParens(getSel.X).(*ast.SelectorExpr)
		if !okRM || rmSel.Sel == nil || rmSel.Sel.Name != "referenceMap" {
			continue
		}
		rmX, okRX := stripParens(rmSel.X).(*ast.Ident)
		vArg, okVA := stripParens(getCall.Args[0]).(*ast.Ident)
		if !okRX || !okVA || rmX.Name != rId.Name || vArg.Name == "" || vArg.Name == "_" {
			continue
		}

		src := fmt.Sprintf(`package p
func f() any {
	rc_ := %[1]s.referenceMap.Get(%[2]s)
	if rc_ == nil && %[2]s != nil && %[2]s.index < 0 {
		if sc_, ok_ := %[2]s.scope.(*Scope); ok_ && sc_ != nil && sc_.implicitVars != nil {
			if cv_, okV_ := sc_.GetOwnSlot(%[2]s.name).(*Var); okV_ && cv_ != nil && cv_ != %[2]s {
				rc_ = %[1]s.referenceMap.Get(cv_)
			}
		}
	}
	return %[3]s.PtrToAny(rc_)
}
`, rId.Name, vArg.Name, j2Qual)
		pf, err := parser.ParseFile(ctx.Fset, "", src, 0)
		if err != nil {
			return 0
		}
		fn.Body.List = slices.Replace(fn.Body.List, retIdx, retIdx+1, pf.Decls[0].(*ast.FuncDecl).Body.List...)
		ctx.MarkFileModified(rcFile)
		return 1
	}
	return 0
}
