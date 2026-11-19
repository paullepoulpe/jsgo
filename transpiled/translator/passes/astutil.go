// Clean-room AST inspection, parent mapping, and safe node replacement helpers.

package passes

import (
	"go/ast"
	"go/token"
	"go/types"
)

// buildParentMap records the immediate AST parent for every node under root,
// and also returns nodes in post-order (children before parents).
func buildParentMap(root ast.Node) (map[ast.Node]ast.Node, []ast.Node) {
	parents := make(map[ast.Node]ast.Node)
	var stack []ast.Node
	var postOrder []ast.Node

	ast.Inspect(root, func(n ast.Node) bool {
		if n == nil {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			postOrder = append(postOrder, top)
			return false
		}
		if len(stack) > 0 {
			parents[n] = stack[len(stack)-1]
		}
		stack = append(stack, n)
		return true
	})
	return parents, postOrder
}

// stripParens unwraps any outer *ast.ParenExpr nodes.
func stripParens(e ast.Expr) ast.Expr {
	for {
		p, ok := e.(*ast.ParenExpr)
		if !ok {
			return e
		}
		e = p.X
	}
}

// calleeFunc extracts the *types.Func called by call, unwrapping generic instantiation
// (*ast.IndexExpr / *ast.IndexListExpr), and also returns the explicit single type argument
// expression when present.
func calleeFunc(call *ast.CallExpr, info *types.Info) (*types.Func, ast.Expr) {
	if call == nil || info == nil {
		return nil, nil
	}
	fun := stripParens(call.Fun)
	var explicitTypeArg ast.Expr
	switch ix := fun.(type) {
	case *ast.IndexExpr:
		fun = stripParens(ix.X)
		explicitTypeArg = ix.Index
	case *ast.IndexListExpr:
		fun = stripParens(ix.X)
		if len(ix.Indices) == 1 {
			explicitTypeArg = ix.Indices[0]
		}
	}
	switch fn := fun.(type) {
	case *ast.Ident:
		if f, ok := info.Uses[fn].(*types.Func); ok {
			return f, explicitTypeArg
		}
	case *ast.SelectorExpr:
		if sel, ok := info.Selections[fn]; ok {
			if f, ok := sel.Obj().(*types.Func); ok {
				return f, explicitTypeArg
			}
		}
		if f, ok := info.Uses[fn.Sel].(*types.Func); ok {
			return f, explicitTypeArg
		}
	}
	return nil, nil
}

// replaceExprInParent replaces oldExpr with newExpr in parent and updates parents[newExpr] = parent.
func replaceExprInParent(parent ast.Node, oldExpr, newExpr ast.Expr, parents map[ast.Node]ast.Node) bool {
	setPtr := func(slot *ast.Expr) bool {
		if *slot == oldExpr {
			*slot = newExpr
			return true
		}
		return false
	}
	setSlice := func(slice []ast.Expr) bool {
		for i, e := range slice {
			if e == oldExpr {
				slice[i] = newExpr
				return true
			}
		}
		return false
	}

	ok := false
	switch p := parent.(type) {
	case *ast.ParenExpr:
		ok = setPtr(&p.X)
	case *ast.SelectorExpr:
		ok = setPtr(&p.X)
	case *ast.IndexExpr:
		ok = setPtr(&p.X) || setPtr(&p.Index)
	case *ast.IndexListExpr:
		ok = setPtr(&p.X) || setSlice(p.Indices)
	case *ast.SliceExpr:
		ok = setPtr(&p.X) || setPtr(&p.Low) || setPtr(&p.High) || setPtr(&p.Max)
	case *ast.TypeAssertExpr:
		ok = setPtr(&p.X)
	case *ast.CallExpr:
		ok = setPtr(&p.Fun) || setSlice(p.Args)
	case *ast.StarExpr:
		ok = setPtr(&p.X)
	case *ast.UnaryExpr:
		ok = setPtr(&p.X)
	case *ast.BinaryExpr:
		ok = setPtr(&p.X) || setPtr(&p.Y)
	case *ast.KeyValueExpr:
		ok = setPtr(&p.Key) || setPtr(&p.Value)
	case *ast.CompositeLit:
		ok = setSlice(p.Elts)
	case *ast.ExprStmt:
		ok = setPtr(&p.X)
	case *ast.SendStmt:
		ok = setPtr(&p.Chan) || setPtr(&p.Value)
	case *ast.IncDecStmt:
		ok = setPtr(&p.X)
	case *ast.AssignStmt:
		ok = setSlice(p.Lhs) || setSlice(p.Rhs)
	case *ast.GoStmt:
		if c, isCall := newExpr.(*ast.CallExpr); isCall && p.Call == oldExpr {
			p.Call = c
			ok = true
		}
	case *ast.DeferStmt:
		if c, isCall := newExpr.(*ast.CallExpr); isCall && p.Call == oldExpr {
			p.Call = c
			ok = true
		}
	case *ast.ReturnStmt:
		ok = setSlice(p.Results)
	case *ast.IfStmt:
		ok = setPtr(&p.Cond)
	case *ast.SwitchStmt:
		ok = setPtr(&p.Tag)
	case *ast.CaseClause:
		ok = setSlice(p.List)
	case *ast.ForStmt:
		ok = setPtr(&p.Cond)
	case *ast.RangeStmt:
		ok = setPtr(&p.Key) || setPtr(&p.Value) || setPtr(&p.X)
	case *ast.ValueSpec:
		ok = setSlice(p.Values)
	}

	if ok && parents != nil {
		parents[newExpr] = parent
	}
	return ok
}

func enclosingFuncBody(n ast.Node, parents map[ast.Node]ast.Node) *ast.BlockStmt {
	for cur := n; ; {
		p := parents[cur]
		if p == nil {
			return nil
		}
		switch fn := p.(type) {
		case *ast.FuncDecl:
			return fn.Body
		case *ast.FuncLit:
			return fn.Body
		}
		cur = p
	}
}

func binaryNeedsParens(parent ast.Node, op token.Token) bool {
	switch p := parent.(type) {
	case *ast.UnaryExpr, *ast.StarExpr, *ast.SelectorExpr, *ast.IndexExpr, *ast.SliceExpr, *ast.TypeAssertExpr:
		return true
	case *ast.BinaryExpr:
		return p.Op.Precedence() >= op.Precedence()
	}
	return false
}

func isStringType(t types.Type) bool {
	if t == nil {
		return false
	}
	b, ok := types.Unalias(t).(*types.Basic)
	return ok && (b.Kind() == types.String || b.Kind() == types.UntypedString)
}

func exprUsesVar(expr ast.Expr, target *types.Var, info *types.Info) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if found {
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

// namedOrigin returns the canonical uninstantiated *types.Named for t, or nil.
func namedOrigin(t types.Type) *types.Named {
	if t == nil {
		return nil
	}
	if n, ok := types.Unalias(t).(*types.Named); ok {
		return n.Origin()
	}
	return nil
}

// ptrNamedOrigin unwraps *T and returns T's canonical *types.Named, or nil.
func ptrNamedOrigin(t types.Type) *types.Named {
	if t == nil {
		return nil
	}
	if p, ok := types.Unalias(t).(*types.Pointer); ok {
		return namedOrigin(p.Elem())
	}
	return nil
}

// derefNamedOrigin unwraps an optional *T and returns T's canonical *types.Named, or nil.
func derefNamedOrigin(t types.Type) *types.Named {
	if t == nil {
		return nil
	}
	ut := types.Unalias(t)
	if p, ok := ut.(*types.Pointer); ok {
		ut = types.Unalias(p.Elem())
	}
	if n, ok := ut.(*types.Named); ok {
		return n.Origin()
	}
	return nil
}

// derefNamedAndTypeArgs unwraps an optional *T and returns T's canonical *types.Named and instantiated TypeArgs, if T is generic.
func derefNamedAndTypeArgs(t types.Type) (*types.Named, *types.TypeList) {
	if t == nil {
		return nil, nil
	}
	ut := types.Unalias(t)
	if p, ok := ut.(*types.Pointer); ok {
		ut = types.Unalias(p.Elem())
	}
	if n, ok := ut.(*types.Named); ok {
		orig := n.Origin()
		if orig.TypeParams() != nil && orig.TypeParams().Len() > 0 && n.TypeArgs() != nil && n.TypeArgs().Len() == orig.TypeParams().Len() {
			return orig, n.TypeArgs()
		}
	}
	return nil, nil
}

// enclosingMethodOwnerAndTArgs returns the generic receiver struct origin and receiver TypeArgs for the method enclosing node.
func enclosingMethodOwnerAndTArgs(node ast.Node, parents map[ast.Node]ast.Node, info *types.Info) (*types.Named, *types.TypeList) {
	for cur := node; cur != nil; cur = parents[cur] {
		if fnDecl, ok := cur.(*ast.FuncDecl); ok {
			if fnDecl.Recv != nil && info != nil {
				if fnObj, ok := info.Defs[fnDecl.Name].(*types.Func); ok && fnObj != nil {
					if sig, ok := fnObj.Type().(*types.Signature); ok && sig.Recv() != nil {
						return derefNamedAndTypeArgs(sig.Recv().Type())
					}
				}
			}
			return nil, nil
		}
	}
	return nil, nil
}

// liftTypeViaTArgs maps bijective receiver TypeParams in targs back to owner.TypeParams().
func liftTypeViaTArgs(t types.Type, owner *types.Named, targs *types.TypeList) (types.Type, bool) {
	if t == nil {
		return nil, false
	}
	if owner == nil || owner.TypeParams() == nil || owner.TypeParams().Len() == 0 {
		return t, true
	}
	tparams := owner.TypeParams()
	if targs == nil || targs.Len() != tparams.Len() {
		return nil, false
	}
	mapping := make(map[*types.TypeParam]types.Type, targs.Len())
	for i := 0; i < targs.Len(); i++ {
		tp, isTP := types.Unalias(targs.At(i)).(*types.TypeParam)
		if !isTP {
			return nil, false
		}
		if _, dup := mapping[tp]; dup {
			return nil, false
		}
		mapping[tp] = tparams.At(i)
	}
	mapped, ok := substTypeParams(t, mapping)
	if !ok || !typeOnlyUsesOwnerTypeParams(mapped, owner) {
		return nil, false
	}
	return mapped, true
}

// sameTypeList reports whether two TypeLists have identical length and element types.
func sameTypeList(a, b *types.TypeList) bool {
	if a == nil || b == nil || a.Len() != b.Len() {
		return a == nil && b == nil
	}
	for i := 0; i < a.Len(); i++ {
		if !types.Identical(a.At(i), b.At(i)) {
			return false
		}
	}
	return true
}

// funcOrigin returns fn.Origin() when fn != nil, or nil.
func funcOrigin(fn *types.Func) *types.Func {
	if fn == nil {
		return nil
	}
	return fn.Origin()
}

// ownerTypeArgsAtAccess extracts the instantiated type argument list of generic struct owner
// at a field selector accessExpr (supporting embedded generic base structs as well).
func ownerTypeArgsAtAccess(accessExpr ast.Expr, owner *types.Named, info *types.Info) *types.TypeList {
	if accessExpr == nil || owner == nil || info == nil || owner.TypeParams() == nil || owner.TypeParams().Len() == 0 {
		return nil
	}
	sel, ok := stripParens(accessExpr).(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	seln, ok := info.Selections[sel]
	if !ok || seln.Kind() != types.FieldVal || len(seln.Index()) == 0 {
		return nil
	}
	curT := seln.Recv()
	for _, fIdx := range seln.Index()[:len(seln.Index())-1] {
		ut := types.Unalias(curT)
		if p, ok := ut.(*types.Pointer); ok {
			ut = types.Unalias(p.Elem())
		}
		st, ok := ut.Underlying().(*types.Struct)
		if !ok || fIdx < 0 || fIdx >= st.NumFields() {
			return nil
		}
		curT = st.Field(fIdx).Type()
	}
	ut := types.Unalias(curT)
	if p, ok := ut.(*types.Pointer); ok {
		ut = types.Unalias(p.Elem())
	}
	if inst, ok := ut.(*types.Named); ok && inst.Origin() == owner {
		return inst.TypeArgs()
	}
	return nil
}

// substTypeParams substitutes TypeParams in t according to mapping.
func substTypeParams(t types.Type, mapping map[*types.TypeParam]types.Type) (types.Type, bool) {
	if t == nil || len(mapping) == 0 || !typeContainsTypeParam(t) {
		return t, true
	}
	switch ut := types.Unalias(t).(type) {
	case *types.TypeParam:
		if repl, ok := mapping[ut]; ok {
			return repl, true
		}
		return ut, true
	case *types.Pointer:
		elt, ok := substTypeParams(ut.Elem(), mapping)
		if !ok {
			return nil, false
		}
		if elt == ut.Elem() {
			return t, true
		}
		return types.NewPointer(elt), true
	case *types.Slice:
		elt, ok := substTypeParams(ut.Elem(), mapping)
		if !ok {
			return nil, false
		}
		if elt == ut.Elem() {
			return t, true
		}
		return types.NewSlice(elt), true
	case *types.Array:
		elt, ok := substTypeParams(ut.Elem(), mapping)
		if !ok {
			return nil, false
		}
		if elt == ut.Elem() {
			return t, true
		}
		return types.NewArray(elt, ut.Len()), true
	case *types.Named:
		targs := ut.TypeArgs()
		if targs == nil || targs.Len() == 0 {
			return t, true
		}
		newArgs := make([]types.Type, targs.Len())
		changed := false
		for i := 0; i < targs.Len(); i++ {
			na, ok := substTypeParams(targs.At(i), mapping)
			if !ok {
				return nil, false
			}
			newArgs[i] = na
			if na != targs.At(i) {
				changed = true
			}
		}
		if !changed {
			return t, true
		}
		inst, err := types.Instantiate(nil, ut.Origin(), newArgs, false)
		if err != nil {
			return nil, false
		}
		return inst, true
	default:
		return nil, false
	}
}

func typeOnlyUsesOwnerTypeParams(t types.Type, owner *types.Named) bool {
	if !typeContainsTypeParam(t) {
		return true
	}
	if owner == nil || owner.TypeParams() == nil || owner.TypeParams().Len() == 0 {
		return false
	}
	allowed := make(map[*types.TypeParam]bool, owner.TypeParams().Len())
	for i := 0; i < owner.TypeParams().Len(); i++ {
		allowed[owner.TypeParams().At(i)] = true
	}
	var check func(types.Type) bool
	check = func(cur types.Type) bool {
		if cur == nil {
			return true
		}
		switch ut := types.Unalias(cur).(type) {
		case *types.TypeParam:
			return allowed[ut]
		case *types.Pointer:
			return check(ut.Elem())
		case *types.Slice:
			return check(ut.Elem())
		case *types.Array:
			return check(ut.Elem())
		case *types.Named:
			if targs := ut.TypeArgs(); targs != nil {
				for i := 0; i < targs.Len(); i++ {
					if !check(targs.At(i)) {
						return false
					}
				}
			}
			return true
		default:
			return !typeContainsTypeParam(cur)
		}
	}
	return check(t)
}

// siteTypeConstraint records a concrete-instantiation observation on a generic struct field.
type siteTypeConstraint struct {
	siteType types.Type
	targs    *types.TypeList
}

// liftSiteTypeForOwner maps receiver TypeParams at accessExpr back to owner.TypeParams().
// If accessExpr has concrete (non-bijective TypeParam) type arguments, isConcrete=true is returned
// so the caller can defer verification until declT is chosen from generic method sites (or fall back
// to siteType if no generic method sites exist).
func liftSiteTypeForOwner(t types.Type, accessExpr ast.Expr, owner *types.Named, info *types.Info) (lifted types.Type, constraint *siteTypeConstraint, ok bool) {
	if t == nil {
		return nil, nil, false
	}
	if owner == nil || owner.TypeParams() == nil || owner.TypeParams().Len() == 0 {
		return t, nil, true
	}
	targs := ownerTypeArgsAtAccess(accessExpr, owner, info)
	tparams := owner.TypeParams()
	if targs == nil || targs.Len() != tparams.Len() {
		return t, nil, typeOnlyUsesOwnerTypeParams(t, owner)
	}
	mapping := make(map[*types.TypeParam]types.Type, targs.Len())
	allDistinctTP := true
	for i := 0; i < targs.Len(); i++ {
		tp, isTP := types.Unalias(targs.At(i)).(*types.TypeParam)
		if !isTP {
			allDistinctTP = false
			break
		}
		if _, dup := mapping[tp]; dup {
			allDistinctTP = false
			break
		}
		mapping[tp] = tparams.At(i)
	}
	if allDistinctTP {
		mapped, mapOK := substTypeParams(t, mapping)
		if !mapOK || !typeOnlyUsesOwnerTypeParams(mapped, owner) {
			return nil, nil, false
		}
		return mapped, nil, true
	}
	if typeContainsTypeParam(t) {
		return nil, nil, false
	}
	return nil, &siteTypeConstraint{siteType: t, targs: targs}, true
}

// instantiateOwnerType substitutes owner.TypeParams() with targs inside declT.
func instantiateOwnerType(declT types.Type, owner *types.Named, targs *types.TypeList) (types.Type, bool) {
	if declT == nil {
		return nil, false
	}
	if owner == nil || owner.TypeParams() == nil || owner.TypeParams().Len() == 0 || targs == nil || !typeContainsTypeParam(declT) {
		return declT, true
	}
	tparams := owner.TypeParams()
	if targs.Len() != tparams.Len() {
		return nil, false
	}
	mapping := make(map[*types.TypeParam]types.Type, tparams.Len())
	for i := 0; i < tparams.Len(); i++ {
		mapping[tparams.At(i)] = targs.At(i)
	}
	return substTypeParams(declT, mapping)
}

// instantiateOwnerTypeAtExpr substitutes owner.TypeParams() with the receiver type args at accessExpr.
func instantiateOwnerTypeAtExpr(declT types.Type, owner *types.Named, accessExpr ast.Expr, info *types.Info) types.Type {
	if declT == nil || owner == nil || accessExpr == nil || !typeContainsTypeParam(declT) {
		return declT
	}
	targs := ownerTypeArgsAtAccess(accessExpr, owner, info)
	if inst, ok := instantiateOwnerType(declT, owner, targs); ok {
		return inst
	}
	return declT
}

// j2gortHasExport reports whether the imported j2gort package declares exported symbol name.
func j2gortHasExport(ctx *PassContext, name string) bool {
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
				return impPkg.Scope() != nil && impPkg.Scope().Lookup(name) != nil
			}
		}
	}
	return false
}

// nodestoreHasExport reports whether the imported or sibling nodestore package declares exported symbol name.
func nodestoreHasExport(ctx *PassContext, name string) bool {
	nsPkg := findNodestorePackage(ctx, nil)
	return nsPkg != nil && nsPkg.Scope() != nil && nsPkg.Scope().Lookup(name) != nil
}
