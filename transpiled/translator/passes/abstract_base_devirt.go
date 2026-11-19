// Clean-room abstract_base_devirt pass: devirtualizes non-overridden exported methods on
// j2go abstract-class interfaces (such as rhino.JSType, ObjectType, PrototypeObjectType,
// jscomp.AbstractScope, AbstractVar) through their direct *<I>Base struct pointer so the
// Go compiler can inline leaf predicates and field getters instead of emitting indirect
// interface calls.

package passes

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

// AbstractBaseDevirtPass devirtualizes calls to non-overridden *<I>Base methods on
// abstract-class interfaces via Get<I>Base() and hoists repeated local base lookups
// once the receiver variable is unconditionally proven non-nil.
type AbstractBaseDevirtPass struct{}

// NewAbstractBaseDevirtPass constructs a new abstract_base_devirt pass instance.
func NewAbstractBaseDevirtPass() *AbstractBaseDevirtPass {
	return &AbstractBaseDevirtPass{}
}

// Name returns "abstract_base_devirt".
func (p *AbstractBaseDevirtPass) Name() string {
	return "abstract_base_devirt"
}

// Description returns a summary of the pass.
func (p *AbstractBaseDevirtPass) Description() string {
	return "Devirtualize non-overridden abstract-class interface calls through *<I>Base"
}

// Apply runs devirtAbstractBaseCalls across target files in ctx.
func (p *AbstractBaseDevirtPass) Apply(ctx *PassContext) (int, error) {
	return devirtAbstractBaseCalls(ctx)
}

type abstractBaseMeta struct {
	valid       bool
	ifaceNamed  *types.Named
	baseNamed   *types.Named
	getterName  string
	safeMethods map[string]bool
}

type abstractDevirtResolver struct {
	ctx               *PassContext
	extraPkgs         []*types.Package
	crossPkgOverrides map[string]map[string]bool
	metaByIface       map[*types.Named]*abstractBaseMeta
}

func newAbstractDevirtResolver(ctx *PassContext) *abstractDevirtResolver {
	r := &abstractDevirtResolver{
		ctx:         ctx,
		metaByIface: make(map[*types.Named]*abstractBaseMeta),
	}
	if mImp, ok := ctx.Importer.(*ModuleImporter); ok && mImp != nil {
		r.extraPkgs = mImp.CachedPackages()
		r.crossPkgOverrides = mImp.CrossPkgAbstractOverrides()
	}
	return r
}

func isExcludedFrameworkBaseMethod(mName string) bool {
	if !ast.IsExported(mName) {
		return true
	}
	if strings.HasPrefix(mName, "Implements_") || strings.HasPrefix(mName, "Slice_") {
		return true
	}
	switch mName {
	case "BindSelf_", "IsJavaNil_", "JavaClassInfo_", "IdentityHash_", "HasIdentityEquals_", "TCHandlePtr_":
		return true
	}
	return false
}

func isAbstractBaseGetterFunc(fn *types.Func) bool {
	if fn == nil {
		return false
	}
	name := fn.Name()
	if !strings.HasPrefix(name, "Get") || !strings.HasSuffix(name, "Base") {
		return false
	}
	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Params().Len() != 0 || sig.Results().Len() != 1 {
		return false
	}
	retPtr, ok := types.Unalias(sig.Results().At(0).Type()).(*types.Pointer)
	if !ok {
		return false
	}
	retNamed := namedOrigin(retPtr.Elem())
	if retNamed == nil || retNamed.Obj() == nil {
		return false
	}
	return "Get"+retNamed.Obj().Name() == name
}

func (r *abstractDevirtResolver) metadataForIface(ifaceNamed *types.Named) *abstractBaseMeta {
	if ifaceNamed == nil {
		return nil
	}
	ifaceNamed = ifaceNamed.Origin()
	if cached, ok := r.metaByIface[ifaceNamed]; ok {
		return cached
	}
	meta := &abstractBaseMeta{ifaceNamed: ifaceNamed}
	r.metaByIface[ifaceNamed] = meta

	if ifaceNamed.TypeParams() != nil && ifaceNamed.TypeParams().Len() > 0 {
		return meta
	}
	obj := ifaceNamed.Obj()
	if obj == nil || obj.Pkg() == nil || isJ2gortPkg(obj.Pkg()) {
		return meta
	}
	ifacePkg := obj.Pkg()
	ifaceType, ok := ifaceNamed.Underlying().(*types.Interface)
	if !ok || ifaceType.NumMethods() == 0 {
		return meta
	}

	baseName := obj.Name() + "Base"
	getterName := "Get" + baseName
	if ifacePkg.Scope() == nil {
		return meta
	}
	baseTN, ok := ifacePkg.Scope().Lookup(baseName).(*types.TypeName)
	if !ok || baseTN == nil {
		return meta
	}
	baseNamed := namedOrigin(baseTN.Type())
	if baseNamed == nil || (baseNamed.TypeParams() != nil && baseNamed.TypeParams().Len() > 0) {
		return meta
	}
	if _, isStruct := baseNamed.Underlying().(*types.Struct); !isStruct {
		return meta
	}
	basePtr := types.NewPointer(baseNamed)
	if !types.Implements(basePtr, ifaceType) {
		return meta
	}

	// Verify both I and *<I>Base declare `Get<I>Base() *<I>Base`.
	ifaceGetterObj, _, _ := types.LookupFieldOrMethod(ifaceNamed, false, ifacePkg, getterName)
	ifaceGetterFn, ok := ifaceGetterObj.(*types.Func)
	if !ok {
		return meta
	}
	igSig, ok := ifaceGetterFn.Type().(*types.Signature)
	if !ok || igSig.Params().Len() != 0 || igSig.Results().Len() != 1 || !types.Identical(igSig.Results().At(0).Type(), basePtr) {
		return meta
	}
	baseGetterObj, _, _ := types.LookupFieldOrMethod(basePtr, true, ifacePkg, getterName)
	baseGetterFn, ok := baseGetterObj.(*types.Func)
	if !ok {
		return meta
	}
	bgSig, ok := baseGetterFn.Type().(*types.Signature)
	if !ok || bgSig.Params().Len() != 0 || bgSig.Results().Len() != 1 || !types.Identical(bgSig.Results().At(0).Type(), basePtr) {
		return meta
	}

	// Collect all distinct concrete struct types implementing I across visible scopes.
	var implStructs []*types.Named
	seenStructs := make(map[*types.Named]bool)
	scanPkgScope := func(p *types.Package) {
		if p == nil || p.Scope() == nil || isJ2gortPkg(p) {
			return
		}
		scope := p.Scope()
		for _, name := range scope.Names() {
			tn, ok := scope.Lookup(name).(*types.TypeName)
			if !ok || tn == nil || tn.IsAlias() {
				continue
			}
			sn := namedOrigin(tn.Type())
			if sn == nil || sn == baseNamed || seenStructs[sn] {
				continue
			}
			if sn.TypeParams() != nil && sn.TypeParams().Len() > 0 {
				continue
			}
			if _, isSt := sn.Underlying().(*types.Struct); !isSt {
				continue
			}
			if types.Implements(types.NewPointer(sn), ifaceType) {
				seenStructs[sn] = true
				implStructs = append(implStructs, sn)
			}
		}
	}
	scanPkgScope(ifacePkg)
	scanPkgScope(r.ctx.Pkg)
	for _, imp := range r.ctx.Pkg.Imports() {
		scanPkgScope(imp)
	}
	for _, ep := range r.extraPkgs {
		scanPkgScope(ep)
	}

	var crossOverridden map[string]bool
	if r.crossPkgOverrides != nil {
		crossOverridden = r.crossPkgOverrides[ifacePkg.Name()+"."+obj.Name()]
	}

	safe := make(map[string]bool)
	for i := 0; i < ifaceType.NumMethods(); i++ {
		im := ifaceType.Method(i)
		mName := im.Name()
		if mName == getterName || isExcludedFrameworkBaseMethod(mName) {
			continue
		}
		if crossOverridden != nil && crossOverridden[mName] {
			continue
		}
		bObj, _, _ := types.LookupFieldOrMethod(basePtr, true, ifacePkg, mName)
		bFn, ok := bObj.(*types.Func)
		if !ok || !bFn.Exported() || isAbstractBaseGetterFunc(bFn) {
			continue
		}
		bSig, ok := bFn.Type().(*types.Signature)
		if !ok || bSig.Recv() == nil {
			continue
		}
		if _, isPtrRecv := types.Unalias(bSig.Recv().Type()).(*types.Pointer); !isPtrRecv {
			continue
		}
		if !types.Identical(bSig, im.Type()) {
			continue
		}
		overridden := false
		for _, sn := range implStructs {
			sObj, _, _ := types.LookupFieldOrMethod(types.NewPointer(sn), true, sn.Obj().Pkg(), mName)
			if sObj != bFn {
				overridden = true
				break
			}
		}
		if !overridden {
			safe[mName] = true
		}
	}

	meta.valid = len(safe) > 0
	meta.baseNamed = baseNamed
	meta.getterName = getterName
	meta.safeMethods = safe
	return meta
}

func (r *abstractDevirtResolver) matchCall(call *ast.CallExpr, info *types.Info) (*ast.SelectorExpr, string, bool) {
	if call == nil || call.Ellipsis.IsValid() {
		return nil, "", false
	}
	selExpr, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selExpr.Sel == nil {
		return nil, "", false
	}
	mName := selExpr.Sel.Name
	if isExcludedFrameworkBaseMethod(mName) {
		return nil, "", false
	}
	sel := info.Selections[selExpr]
	if sel == nil || sel.Kind() != types.MethodVal {
		return nil, "", false
	}
	recvTV, ok := info.Types[selExpr.X]
	if !ok || recvTV.IsType() || recvTV.IsBuiltin() || recvTV.IsNil() || recvTV.Type == nil || isInvalidOrUntyped(recvTV.Type) {
		return nil, "", false
	}
	ut := types.Unalias(recvTV.Type)
	if _, isTP := ut.(*types.TypeParam); isTP {
		return nil, "", false
	}
	if !isNonEmptyInterfaceType(ut) {
		return nil, "", false
	}

	// 1. Direct named abstract-class interface receiver (e.g. JSType, ObjectType, AbstractScope).
	if recvNamed, ok := ut.(*types.Named); ok {
		if meta := r.metadataForIface(recvNamed); meta != nil && meta.valid && meta.safeMethods[mName] {
			return selExpr, meta.getterName, true
		}
	}

	// 2. Sub-interface embedding an abstract-class interface (e.g. JSType_WithSourceRef).
	ifaceU, ok := ut.Underlying().(*types.Interface)
	if !ok {
		return nil, "", false
	}
	var candOrder []*types.Named
	candSeen := make(map[*types.Named]bool)
	var collectEmbedded func(it *types.Interface)
	collectEmbedded = func(it *types.Interface) {
		if it == nil {
			return
		}
		for i := 0; i < it.NumEmbeddeds(); i++ {
			et := types.Unalias(it.EmbeddedType(i))
			if en, ok := et.(*types.Named); ok {
				orig := en.Origin()
				if !candSeen[orig] {
					candSeen[orig] = true
					candOrder = append(candOrder, orig)
				}
				if eu, okU := orig.Underlying().(*types.Interface); okU {
					collectEmbedded(eu)
				}
			} else if eu, okU := et.(*types.Interface); okU {
				collectEmbedded(eu)
			}
		}
	}
	collectEmbedded(ifaceU)
	for _, en := range candOrder {
		if meta := r.metadataForIface(en); meta != nil && meta.valid && meta.safeMethods[mName] {
			return selExpr, meta.getterName, true
		}
	}
	return nil, "", false
}

func devirtAbstractBaseCalls(ctx *PassContext) (int, error) {
	if ctx == nil || ctx.Pkg == nil || ctx.Info == nil || isJ2gortPkg(ctx.Pkg) || len(ctx.Files) == 0 {
		return 0, nil
	}
	resolver := newAbstractDevirtResolver(ctx)
	total := 0
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		n := devirtAbstractBaseFile(file, resolver, ctx.Info)
		if n > 0 {
			ctx.MarkFileModified(file)
			total += n
		}
	}
	return total, nil
}

type devirtCallSite struct {
	selExpr    *ast.SelectorExpr
	getterName string
	recvVar    *types.Var // non-nil when stripParens(selExpr.X) is a bare *ast.Ident variable
	topStmtIdx int        // index in enclosing function's top-level BlockStmt (-1 if outside)
}

func devirtAbstractBaseFile(file *ast.File, resolver *abstractDevirtResolver, info *types.Info) int {
	count := 0

	// Rewrite any top-level package var initializers outside FuncDecl bodies inline.
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Body != nil {
				count += devirtFuncTree(d.Type, d.Body, resolver, info)
			}
		case *ast.GenDecl:
			ast.Inspect(d, func(n ast.Node) bool {
				if fl, ok := n.(*ast.FuncLit); ok && fl.Body != nil {
					count += devirtFuncTree(fl.Type, fl.Body, resolver, info)
					return false
				}
				if call, ok := n.(*ast.CallExpr); ok {
					if selExpr, getter, okM := resolver.matchCall(call, info); okM {
						rewriteSelectorToInlineBaseGet(selExpr, getter)
						count++
					}
				}
				return true
			})
		}
	}
	return count
}

func rewriteSelectorToInlineBaseGet(selExpr *ast.SelectorExpr, getterName string) {
	pos := selExpr.Sel.Pos()
	selExpr.X = &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   wrapPrimaryTarget(selExpr.X),
			Sel: &ast.Ident{NamePos: pos, Name: getterName},
		},
	}
}

func devirtFuncTree(fnType *ast.FuncType, body *ast.BlockStmt, resolver *abstractDevirtResolver, info *types.Info) int {
	if body == nil {
		return 0
	}
	count := 0

	// First recursively process nested FuncLits for their own parameters/locals,
	// while recording which calls have already been rewritten so the outer pass skips them.
	handledSels := make(map[*ast.SelectorExpr]bool)
	var visitNestedLits func(root ast.Node)
	visitNestedLits = func(root ast.Node) {
		ast.Inspect(root, func(n ast.Node) bool {
			if n == nil || n == root {
				return true
			}
			if fl, ok := n.(*ast.FuncLit); ok && fl.Body != nil {
				count += devirtSingleFunc(fl.Type, fl.Body, resolver, info, handledSels)
				return false
			}
			return true
		})
	}
	visitNestedLits(body)
	count += devirtSingleFunc(fnType, body, resolver, info, handledSels)
	return count
}

func devirtSingleFunc(
	fnType *ast.FuncType,
	body *ast.BlockStmt,
	resolver *abstractDevirtResolver,
	info *types.Info,
	handledSels map[*ast.SelectorExpr]bool,
) int {
	// Collect all devirtualizable calls inside body not yet handled by an inner FuncLit
	// (an inner FuncLit only hoists its own local/param variables; any unhandled call inside
	// a nested FuncLit referencing an outer variable or expression is still collected here).
	var sites []devirtCallSite
	for stmtIdx, stmt := range body.List {
		ast.Inspect(stmt, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			selExpr, getter, okM := resolver.matchCall(call, info)
			if !okM || handledSels[selExpr] {
				return true
			}
			var rv *types.Var
			if id, okId := stripParens(selExpr.X).(*ast.Ident); okId {
				rv = varFromObj(info.Uses[id])
			}
			sites = append(sites, devirtCallSite{
				selExpr:    selExpr,
				getterName: getter,
				recvVar:    rv,
				topStmtIdx: stmtIdx,
			})
			return true
		})
	}
	if len(sites) == 0 {
		return 0
	}

	// Attempt safe block-level base hoisting for unmodified parameters / top-level locals
	// when the function body has no `goto` statements.
	hoistedForSite := make([]string, len(sites))
	if !blockHasGotoStmt(body) {
		tryPlanBaseHoists(fnType, body, sites, hoistedForSite, info)
	}

	for i, s := range sites {
		handledSels[s.selExpr] = true
		if hoistedForSite[i] != "" {
			s.selExpr.X = &ast.Ident{NamePos: s.selExpr.X.Pos(), Name: hoistedForSite[i]}
		} else {
			rewriteSelectorToInlineBaseGet(s.selExpr, s.getterName)
		}
	}
	return len(sites)
}

type hoistVarDeclInfo struct {
	v        *types.Var
	name     string
	declIdx  int // -1 for function parameter, >= 0 for body.List[declIdx]
	initExpr ast.Expr
}

func tryPlanBaseHoists(
	fnType *ast.FuncType,
	body *ast.BlockStmt,
	sites []devirtCallSite,
	hoistedForSite []string,
	info *types.Info,
) {
	// 1. Discover candidate variables declared in fnType.Params or top-level body.List.
	var candOrder []*types.Var
	candInfo := make(map[*types.Var]*hoistVarDeclInfo)
	addCand := func(id *ast.Ident, declIdx int, initExpr ast.Expr) {
		if id == nil || id.Name == "_" || id.Name == "" {
			return
		}
		v := varFromObj(info.Defs[id])
		if v == nil || v.IsField() || candInfo[v] != nil {
			return
		}
		candInfo[v] = &hoistVarDeclInfo{
			v:        v,
			name:     id.Name,
			declIdx:  declIdx,
			initExpr: initExpr,
		}
		candOrder = append(candOrder, v)
	}

	if fnType != nil && fnType.Params != nil {
		for _, fld := range fnType.Params.List {
			for _, nm := range fld.Names {
				addCand(nm, -1, nil)
			}
		}
	}
	for i, stmt := range body.List {
		switch s := stmt.(type) {
		case *ast.DeclStmt:
			gd, ok := s.Decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR || len(gd.Specs) != 1 {
				continue
			}
			vs, ok := gd.Specs[0].(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || len(vs.Values) > 1 {
				continue
			}
			var initExpr ast.Expr
			if len(vs.Values) == 1 {
				initExpr = vs.Values[0]
			}
			addCand(vs.Names[0], i, initExpr)
		case *ast.AssignStmt:
			if s.Tok == token.DEFINE && len(s.Lhs) == 1 && len(s.Rhs) == 1 {
				if id, ok := s.Lhs[0].(*ast.Ident); ok {
					addCand(id, i, s.Rhs[0])
				}
			}
		}
	}
	if len(candOrder) == 0 {
		return
	}

	// 2. Scan body once to collect all used identifier names and disqualify any variable
	// that is reassigned, incremented, redeclared, or has its address taken.
	usedNames := make(map[string]bool)
	if fnType != nil {
		ast.Inspect(fnType, func(n ast.Node) bool {
			if id, ok := n.(*ast.Ident); ok {
				usedNames[id.Name] = true
			}
			return true
		})
	}
	defCounts := make(map[*types.Var]int, len(candOrder))
	mutated := make(map[*types.Var]bool)

	ast.Inspect(body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			usedNames[x.Name] = true
			if v := varFromObj(info.Defs[x]); v != nil && candInfo[v] != nil {
				defCounts[v]++
				if defCounts[v] > 1 {
					mutated[v] = true
				}
			}
		case *ast.AssignStmt:
			for idx, lhs := range x.Lhs {
				id, ok := stripParens(lhs).(*ast.Ident)
				if !ok || id.Name == "_" {
					continue
				}
				if x.Tok == token.DEFINE {
					if v := varFromObj(info.Defs[id]); v == nil {
						if uv := varFromObj(info.Uses[id]); uv != nil {
							mutated[uv] = true
						}
					} else if ci := candInfo[v]; ci != nil && (len(x.Lhs) != 1 || body.List[ci.declIdx] != x || idx != 0) {
						mutated[v] = true
					}
					continue
				}
				if uv := varFromObj(info.Uses[id]); uv != nil {
					mutated[uv] = true
				}
			}
		case *ast.IncDecStmt:
			if id, ok := stripParens(x.X).(*ast.Ident); ok {
				if uv := varFromObj(info.Uses[id]); uv != nil {
					mutated[uv] = true
				}
			}
		case *ast.UnaryExpr:
			if x.Op == token.AND {
				if id, ok := stripParens(x.X).(*ast.Ident); ok {
					if uv := varFromObj(info.Uses[id]); uv != nil {
						mutated[uv] = true
					}
				}
			}
		}
		return true
	})

	// Also verify that no nested scope shadows v.name after its declaration.
	shadowed := make(map[*types.Var]bool)
	ast.Inspect(body, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if defObj := info.Defs[id]; defObj != nil {
			for _, cv := range candOrder {
				if cv.Name() == id.Name && defObj != cv {
					shadowed[cv] = true
				}
			}
		}
		return true
	})

	// 3. Group call site indices by candidate variable.
	siteIdxsByVar := make(map[*types.Var][]int)
	for i, s := range sites {
		if s.recvVar == nil {
			continue
		}
		ci := candInfo[s.recvVar]
		if ci == nil || mutated[s.recvVar] || shadowed[s.recvVar] {
			continue
		}
		siteIdxsByVar[s.recvVar] = append(siteIdxsByVar[s.recvVar], i)
	}

	insertBefore := make(map[int][]ast.Stmt)
	anyHoist := false

	for _, cv := range candOrder {
		idxs := siteIdxsByVar[cv]
		if len(idxs) < 2 {
			continue
		}
		ci := candInfo[cv]
		getter := sites[idxs[0]].getterName
		sameGetter := true
		for _, si := range idxs[1:] {
			if sites[si].getterName != getter {
				sameGetter = false
				break
			}
		}
		if !sameGetter {
			continue
		}

		hoistIdx := findEarliestSafeBaseHoistIdx(ci, body.List, info)
		if hoistIdx < 0 || hoistIdx >= len(body.List) {
			continue
		}

		afterCount := 0
		for _, si := range idxs {
			if sites[si].topStmtIdx >= hoistIdx {
				afterCount++
			}
		}
		if afterCount < 2 {
			continue
		}

		cacheName := pickFreshDevirtIdent(usedNames, "base_"+ci.name+"_")
		usedNames[cacheName] = true
		pos := body.List[hoistIdx].Pos()
		hoistStmt := &ast.AssignStmt{
			Lhs:    []ast.Expr{&ast.Ident{NamePos: pos, Name: cacheName}},
			TokPos: pos,
			Tok:    token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{NamePos: pos, Name: ci.name},
						Sel: &ast.Ident{NamePos: pos, Name: getter},
					},
				},
			},
		}
		insertBefore[hoistIdx] = append(insertBefore[hoistIdx], hoistStmt)
		anyHoist = true
		for _, si := range idxs {
			if sites[si].topStmtIdx >= hoistIdx {
				hoistedForSite[si] = cacheName
			}
		}
	}

	if !anyHoist {
		return
	}

	newList := make([]ast.Stmt, 0, len(body.List)+len(candOrder))
	for i, stmt := range body.List {
		if pre := insertBefore[i]; len(pre) > 0 {
			newList = append(newList, pre...)
		}
		newList = append(newList, stmt)
	}
	body.List = newList
}

func pickFreshDevirtIdent(used map[string]bool, prefix string) string {
	if !used[prefix] {
		return prefix
	}
	for i := 2; ; i++ {
		cand := fmt.Sprintf("%s%d_", prefix, i)
		if !used[cand] {
			return cand
		}
	}
}

func blockHasGotoStmt(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		if br, ok := n.(*ast.BranchStmt); ok && br.Tok == token.GOTO {
			found = true
			return false
		}
		return true
	})
	return found
}

func findEarliestSafeBaseHoistIdx(ci *hoistVarDeclInfo, stmts []ast.Stmt, info *types.Info) int {
	startIdx := 0
	if ci.declIdx >= 0 {
		startIdx = ci.declIdx + 1
		if ci.initExpr != nil && isGuaranteedNonNilInitForDevirt(ci.initExpr, info) {
			return startIdx
		}
	}
	for i := startIdx; i < len(stmts); i++ {
		st := stmts[i]
		if stmtUnconditionallyDerefsVarFirst(st, ci.v, info) {
			return i
		}
		if ifSt, ok := st.(*ast.IfStmt); ok && ifSt.Init == nil && ifSt.Else == nil && blockUnconditionallyReturnsOrPanics(ifSt.Body, info) {
			if exprProvesVarNonNilWhenFalse(ifSt.Cond, ci.v, info) {
				return i + 1
			}
		}
	}
	return -1
}

func isGuaranteedNonNilInitForDevirt(expr ast.Expr, info *types.Info) bool {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() {
		return false
	}
	if fnName, _, _, isJ2 := isJ2gortCall(call, info); isJ2 {
		switch fnName {
		case "CheckNotNull", "CheckNotNullIface", "CheckNotNullPtr":
			return true
		}
	}
	return false
}

func blockUnconditionallyReturnsOrPanics(blk *ast.BlockStmt, info *types.Info) bool {
	if blk == nil || len(blk.List) == 0 {
		return false
	}
	last := blk.List[len(blk.List)-1]
	switch s := last.(type) {
	case *ast.ReturnStmt:
		return true
	case *ast.ExprStmt:
		if call, ok := stripParens(s.X).(*ast.CallExpr); ok {
			if id, okId := stripParens(call.Fun).(*ast.Ident); okId && id.Name == "panic" {
				if b, isB := info.Uses[id].(*types.Builtin); isB && b.Name() == "panic" {
					return true
				}
			}
		}
	}
	return false
}

func stmtUnconditionallyDerefsVarFirst(stmt ast.Stmt, v *types.Var, info *types.Info) bool {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		return exprDerefsVarFirst(s.X, v, info)
	case *ast.DeclStmt:
		gd, ok := s.Decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR || len(gd.Specs) != 1 {
			return false
		}
		vs, ok := gd.Specs[0].(*ast.ValueSpec)
		if !ok || len(vs.Values) == 0 {
			return false
		}
		return exprDerefsVarFirst(vs.Values[0], v, info)
	case *ast.AssignStmt:
		if len(s.Lhs) == 1 && len(s.Rhs) == 1 {
			if _, ok := stripParens(s.Lhs[0]).(*ast.Ident); ok {
				return exprDerefsVarFirst(s.Rhs[0], v, info)
			}
		}
	case *ast.ReturnStmt:
		if len(s.Results) >= 1 {
			return exprDerefsVarFirst(s.Results[0], v, info)
		}
	case *ast.IfStmt:
		if s.Init == nil {
			return exprDerefsVarFirst(s.Cond, v, info)
		}
	case *ast.SwitchStmt:
		if s.Init == nil && s.Tag != nil {
			return exprDerefsVarFirst(s.Tag, v, info)
		}
	}
	return false
}

func exprDerefsVarFirst(expr ast.Expr, v *types.Var, info *types.Info) bool {
	if expr == nil {
		return false
	}
	switch e := stripParens(expr).(type) {
	case *ast.UnaryExpr:
		return exprDerefsVarFirst(e.X, v, info)
	case *ast.BinaryExpr:
		return exprDerefsVarFirst(e.X, v, info)
	case *ast.SelectorExpr:
		return exprDerefsVarFirst(e.X, v, info)
	case *ast.CallExpr:
		if sel, ok := stripParens(e.Fun).(*ast.SelectorExpr); ok && sel.Sel != nil {
			if id, okId := stripParens(sel.X).(*ast.Ident); okId && varFromObj(info.Uses[id]) == v {
				if sel.Sel.Name != "IsJavaNil_" && allPureArgsForDevirt(e.Args) {
					return true
				}
			}
			if exprDerefsVarFirst(sel.X, v, info) {
				return true
			}
		}
	}
	return false
}

func allPureArgsForDevirt(args []ast.Expr) bool {
	for _, a := range args {
		if !isPureArgForDevirt(a) {
			return false
		}
	}
	return true
}

func isPureArgForDevirt(expr ast.Expr) bool {
	switch e := stripParens(expr).(type) {
	case *ast.BasicLit, *ast.Ident:
		return true
	case *ast.UnaryExpr:
		return e.Op != token.ARROW && isPureArgForDevirt(e.X)
	case *ast.BinaryExpr:
		return isPureArgForDevirt(e.X) && isPureArgForDevirt(e.Y)
	}
	return false
}

func exprProvesVarNonNilWhenFalse(expr ast.Expr, v *types.Var, info *types.Info) bool {
	if expr == nil {
		return false
	}
	e := stripParens(expr)
	if isVarJavaNilGuardDisjunction(e, v, info) {
		return true
	}
	switch x := e.(type) {
	case *ast.BinaryExpr:
		if x.Op == token.LOR {
			return exprProvesVarNonNilWhenFalse(x.X, v, info) || exprProvesVarNonNilWhenFalse(x.Y, v, info)
		}
	case *ast.UnaryExpr:
		if x.Op == token.NOT {
			return exprProvesVarNonNilWhenTrue(x.X, v, info)
		}
	case *ast.CallExpr:
		if isNonNilCheckMethodCallOnVar(x, v, info) {
			return true
		}
	}
	return false
}

func exprProvesVarNonNilWhenTrue(expr ast.Expr, v *types.Var, info *types.Info) bool {
	if expr == nil {
		return false
	}
	e := stripParens(expr)
	if isVarJavaNonNilConjunction(e, v, info) {
		return true
	}
	switch x := e.(type) {
	case *ast.BinaryExpr:
		if x.Op == token.LAND {
			return exprProvesVarNonNilWhenTrue(x.X, v, info) || exprProvesVarNonNilWhenTrue(x.Y, v, info)
		}
	case *ast.UnaryExpr:
		if x.Op == token.NOT {
			return exprProvesVarNonNilWhenFalse(x.X, v, info)
		}
	case *ast.CallExpr:
		if isNonNilCheckMethodCallOnVar(x, v, info) {
			return true
		}
	}
	return false
}

func isNonNilCheckMethodCallOnVar(call *ast.CallExpr, v *types.Var, info *types.Info) bool {
	if call == nil || call.Ellipsis.IsValid() {
		return false
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || sel.Sel.Name == "IsJavaNil_" {
		return false
	}
	if s := info.Selections[sel]; s == nil || s.Kind() != types.MethodVal {
		return false
	}
	id, ok := stripParens(sel.X).(*ast.Ident)
	return ok && varFromObj(info.Uses[id]) == v
}

// isVarJavaNilGuardDisjunction matches `v == nil || v.IsJavaNil_()`.
func isVarJavaNilGuardDisjunction(expr ast.Expr, v *types.Var, info *types.Info) bool {
	bin, ok := stripParens(expr).(*ast.BinaryExpr)
	if !ok || bin.Op != token.LOR {
		return false
	}
	return isVarNilComparison(bin.X, token.EQL, v, info) && isVarIsJavaNilCall(bin.Y, v, info)
}

// isVarJavaNonNilConjunction matches `v != nil && !v.IsJavaNil_()`.
func isVarJavaNonNilConjunction(expr ast.Expr, v *types.Var, info *types.Info) bool {
	bin, ok := stripParens(expr).(*ast.BinaryExpr)
	if !ok || bin.Op != token.LAND {
		return false
	}
	if !isVarNilComparison(bin.X, token.NEQ, v, info) {
		return false
	}
	un, ok := stripParens(bin.Y).(*ast.UnaryExpr)
	return ok && un.Op == token.NOT && isVarIsJavaNilCall(un.X, v, info)
}

func isVarNilComparison(expr ast.Expr, wantOp token.Token, v *types.Var, info *types.Info) bool {
	bin, ok := stripParens(expr).(*ast.BinaryExpr)
	if !ok || bin.Op != wantOp {
		return false
	}
	x, y := stripParens(bin.X), stripParens(bin.Y)
	if id, ok := x.(*ast.Ident); ok && varFromObj(info.Uses[id]) == v && isNilIdent(y, info) {
		return true
	}
	if id, ok := y.(*ast.Ident); ok && varFromObj(info.Uses[id]) == v && isNilIdent(x, info) {
		return true
	}
	return false
}

func isVarIsJavaNilCall(expr ast.Expr, v *types.Var, info *types.Info) bool {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 0 || call.Ellipsis.IsValid() {
		return false
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || sel.Sel.Name != "IsJavaNil_" {
		return false
	}
	id, ok := stripParens(sel.X).(*ast.Ident)
	return ok && varFromObj(info.Uses[id]) == v
}
