// Clean-room node_traversal_lowering pass: pools *NodeTraversal via sync.Pool, embeds
// inline [12]Node / [12]AbstractScope scope stacks, reuses an inline default
// SyntacticScopeCreator, emits AST-verified AlwaysShouldTraverse_ / NoOpVisit_ callback
// markers across jscomp, and rewrites traverseBranch / traverseChildren with FirstNextTok_()
// and exact leaf unrolling.

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

// NodeTraversalLoweringPass lowers jscomp.NodeTraversal to eliminate per-traversal heap
// allocations, devirtualize trivial ShouldTraverse_3 / Visit_3 callbacks, and unroll leaf
// traversal using single-lookup FirstNextTok_().
type NodeTraversalLoweringPass struct{}

// NewNodeTraversalLoweringPass constructs a new node_traversal_lowering pass instance.
func NewNodeTraversalLoweringPass() *NodeTraversalLoweringPass {
	return &NodeTraversalLoweringPass{}
}

// Name returns "node_traversal_lowering".
func (p *NodeTraversalLoweringPass) Name() string {
	return "node_traversal_lowering"
}

// Description returns a summary of the pass.
func (p *NodeTraversalLoweringPass) Description() string {
	return "Pool jscomp.NodeTraversal, devirtualize AlwaysShouldTraverse_/NoOpVisit_ callbacks, and unroll FirstNextTok_ leaf traversal"
}

// Apply runs node_traversal_lowering on package jscomp if NodeTraversal or Scope is present.
func (p *NodeTraversalLoweringPass) Apply(ctx *PassContext) (int, error) {
	if ctx.Pkg == nil || ctx.Pkg.Name() != "jscomp" || ctx.Pkg.Scope() == nil {
		return 0, nil
	}
	scope := ctx.Pkg.Scope()
	hasNTPrereqs := true
	for _, req := range []string{
		"NodeTraversal",
		"NodeTraversal_Builder",
		"InitNodeTraversal",
		"clinit_NodeTraversal",
		"SyntacticScopeCreator",
		"InitSyntacticScopeCreator_1",
		"StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER",
	} {
		if scope.Lookup(req) == nil {
			hasNTPrereqs = false
			break
		}
	}

	total := 0

	if hasNTPrereqs {
		// 1. Emit AST-verified AlwaysShouldTraverse_(), NoOpVisit_(), and CanReleaseSyntacticScopes_() markers across jscomp.
		markerSites, err := p.emitCallbackMarkers(ctx)
		if err != nil {
			return total, err
		}
		total += markerSites

		// 2. Transform NodeTraversal struct, InitNodeTraversal, pooling, and traverseBranch / traverseChildren.
		for _, file := range ctx.Files {
			if !ctx.ShouldTransform(file) {
				continue
			}
			ntStruct, nodeTypeExpr, nodeTypeSrc, ok := findNodeTraversalStruct(file, ctx.Fset)
			if !ok {
				continue
			}
			sites, err := p.transformNodeTraversalFile(ctx, file, ntStruct, nodeTypeExpr, nodeTypeSrc)
			if err != nil {
				return total, err
			}
			if sites > 0 {
				ctx.MarkFileModified(file)
				total += sites
			}
		}
	}

	// 3. Zero-alloc Scope/TypedScope GetVarIterable() / GetAllSymbols() cursor loop lowering via scopeVarCursor_.
	scopeSites, err := p.applyScopeVarCursorLowering(ctx)
	if err != nil {
		return total, err
	}
	total += scopeSites

	return total, nil
}

type cbStructCandidate struct {
	file               *ast.File
	ts                 *ast.TypeSpec
	st                 *ast.StructType
	originNamed        *types.Named
	ptrT               *types.Pointer
	isConcrete         bool
	alwaysTrav         bool
	noOpVisit          bool
	canRelScope        bool
	canRelScopeDynamic bool
	unwrapCapField     string
	embeddedOrig       []*types.Named
}

func syntacticScopeCreatorHasReleaseSlabs(ctx *PassContext) bool {
	if ctx == nil || ctx.Pkg == nil || ctx.Pkg.Scope() == nil {
		return false
	}
	obj := ctx.Pkg.Scope().Lookup("SyntacticScopeCreator")
	if obj == nil || obj.Type() == nil {
		return false
	}
	return hasMethodNamed(types.NewPointer(obj.Type()), ctx.Pkg, "ReleaseSlabs_")
}

func (p *NodeTraversalLoweringPass) emitCallbackMarkers(ctx *PassContext) (int, error) {
	if ctx.Info == nil || ctx.Info.Defs == nil || ctx.Info.Uses == nil {
		return 0, nil
	}
	cbObj := ctx.Pkg.Scope().Lookup("NodeTraversal_Callback")
	if cbObj == nil {
		return 0, nil
	}
	callbackIface, ok := cbObj.Type().Underlying().(*types.Interface)
	if !ok {
		return 0, nil
	}
	var postOrderCbIface *types.Interface
	var postOrderCbNamed *types.Named
	if poObj := ctx.Pkg.Scope().Lookup("NodeTraversal_AbstractPostOrderCallbackInterface"); poObj != nil && poObj.Type() != nil {
		postOrderCbNamed, _ = poObj.Type().(*types.Named)
		postOrderCbIface, _ = poObj.Type().Underlying().(*types.Interface)
	}
	sscHasRelease := syntacticScopeCreatorHasReleaseSlabs(ctx)

	// Build *types.Func -> *ast.FuncDecl map for all method and top-level function declarations across ctx.Files.
	methodDecls := make(map[*types.Func]*ast.FuncDecl)
	topFuncDecls := make(map[*types.Func]*ast.FuncDecl)
	var allConcreteStructs []*types.Named
	for _, file := range ctx.Files {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Name == nil || d.Body == nil {
					continue
				}
				fn, ok := ctx.Info.Defs[d.Name].(*types.Func)
				if !ok || fn == nil {
					continue
				}
				if d.Recv != nil {
					methodDecls[fn] = d
					if orig := fn.Origin(); orig != nil {
						methodDecls[orig] = d
					}
				} else {
					topFuncDecls[fn] = d
					if orig := fn.Origin(); orig != nil {
						topFuncDecls[orig] = d
					}
				}
			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}
				for _, spec := range d.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok || ts.Assign.IsValid() || ts.Name == nil || strings.HasSuffix(ts.Name.Name, "Base") {
						continue
					}
					if _, isSt := ts.Type.(*ast.StructType); !isSt {
						continue
					}
					if tn, ok := ctx.Info.Defs[ts.Name].(*types.TypeName); ok && tn != nil {
						if nm, ok := tn.Type().(*types.Named); ok && nm != nil {
							allConcreteStructs = append(allConcreteStructs, nm)
						}
					}
				}
			}
		}
	}

	pkgScopeVarFreeIfaces := computePkgScopeVarFreeInterfaces(ctx.Pkg, allConcreteStructs)
	canEmitDynamicRelHelper := ctx.Pkg.Scope().Lookup("canReleaseSyntacticScopesForCallback_") != nil
	if !canEmitDynamicRelHelper {
		for _, f := range ctx.Files {
			if !ctx.ShouldTransform(f) {
				continue
			}
			if _, _, _, ok := findNodeTraversalStruct(f, ctx.Fset); ok {
				canEmitDynamicRelHelper = true
				break
			}
		}
	}

	var candidates []*cbStructCandidate
	byOrigin := make(map[*types.Named]*cbStructCandidate)

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
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				tn, ok := ctx.Info.Defs[ts.Name].(*types.TypeName)
				if !ok || tn == nil {
					continue
				}
				named, ok := tn.Type().(*types.Named)
				if !ok || named == nil {
					continue
				}
				originNamed := named.Origin()
				instNamed := named
				if tp := named.TypeParams(); tp != nil && tp.Len() > 0 && named.TypeArgs() == nil {
					targs := make([]types.Type, tp.Len())
					for i := 0; i < tp.Len(); i++ {
						targs[i] = tp.At(i)
					}
					inst, err := types.Instantiate(nil, named, targs, false)
					if err != nil {
						continue
					}
					instNamed, ok = inst.(*types.Named)
					if !ok {
						continue
					}
				}
				ptrT := types.NewPointer(instNamed)
				isConcrete := !strings.HasSuffix(ts.Name.Name, "Base")
				var alwaysTrav, noOpVisit, canRelScope, canRelScopeDynamic bool
				var unwrapCapField string
				if isConcrete && types.Implements(ptrT, callbackIface) {
					alwaysTrav = canAlwaysShouldTraverse(ptrT, ctx.Pkg, ctx.Info, methodDecls)
					noOpVisit = canNoOpVisit(ptrT, ctx.Pkg, ctx.Info, methodDecls)
					if postOrderCbNamed != nil {
						unwrapCapField = matchPostOrderCallbackAdapterStruct(instNamed, ptrT, postOrderCbNamed, ctx.Pkg, ctx.Info, methodDecls)
					}
					if sscHasRelease && unwrapCapField == "" {
						canRelScope = canReleaseSyntacticScopes(ptrT, instNamed, ctx.Pkg, ctx.Info, methodDecls, topFuncDecls, pkgScopeVarFreeIfaces)
						if !canRelScope && canEmitDynamicRelHelper {
							canRelScopeDynamic = matchCombinedCompilerPassStruct(instNamed, callbackIface, ctx.Pkg, pkgScopeVarFreeIfaces)
						}
					}
				} else if isConcrete && postOrderCbIface != nil && types.Implements(ptrT, postOrderCbIface) &&
					!isSynthesizedPostOrderLambdaName(ts.Name.Name) {
					noOpVisit = canNoOpPostOrderInterfaceVisit(ptrT, ctx.Pkg, ctx.Info, methodDecls)
					if sscHasRelease {
						canRelScope = canReleaseSyntacticScopes(ptrT, instNamed, ctx.Pkg, ctx.Info, methodDecls, topFuncDecls, pkgScopeVarFreeIfaces)
					}
				}
				cand := &cbStructCandidate{
					file:               file,
					ts:                 ts,
					st:                 st,
					originNamed:        originNamed,
					ptrT:               ptrT,
					isConcrete:         isConcrete,
					alwaysTrav:         alwaysTrav,
					noOpVisit:          noOpVisit,
					canRelScope:        canRelScope,
					canRelScopeDynamic: canRelScopeDynamic,
					unwrapCapField:     unwrapCapField,
					embeddedOrig:       collectDirectEmbeddedNamed(instNamed),
				}
				candidates = append(candidates, cand)
				byOrigin[originNamed] = cand
			}
		}
	}

	// Soundness propagation: if struct E is embedded (directly or transitively) inside
	// another struct S in ctx.Pkg where *S does NOT also qualify for AlwaysShouldTraverse_
	// (or NoOpVisit_ or CanReleaseSyntacticScopes_), E must NOT receive the marker method
	// because Go method promotion would otherwise promote E's marker onto *S.
	changed := true
	for changed {
		changed = false
		for _, outer := range candidates {
			for _, embOrig := range outer.embeddedOrig {
				embCand := byOrigin[embOrig]
				if embCand == nil {
					continue
				}
				if embCand.alwaysTrav && (!outer.isConcrete || !outer.alwaysTrav) {
					embCand.alwaysTrav = false
					changed = true
				}
				if embCand.noOpVisit && (!outer.isConcrete || !outer.noOpVisit) {
					embCand.noOpVisit = false
					changed = true
				}
				if embCand.canRelScope && (!outer.isConcrete || !outer.canRelScope) {
					embCand.canRelScope = false
					changed = true
				}
				if embCand.canRelScopeDynamic {
					embCand.canRelScopeDynamic = false
					changed = true
				}
				if embCand.unwrapCapField != "" {
					embCand.unwrapCapField = ""
					changed = true
				}
			}
		}
	}

	totalSites := 0
	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		fileSites := 0
		for _, cand := range candidates {
			if cand.file != file || !cand.isConcrete ||
				(!cand.alwaysTrav && !cand.noOpVisit && !cand.canRelScope && !cand.canRelScopeDynamic && cand.unwrapCapField == "") {
				continue
			}
			recvType := formatPtrReceiverType(cand.ts)
			if recvType == "" {
				continue
			}
			if cand.alwaysTrav &&
				!hasMethodNamed(cand.ptrT, ctx.Pkg, "AlwaysShouldTraverse_") &&
				!fileHasPtrMethod(file, cand.ts.Name.Name, "AlwaysShouldTraverse_") {
				decls, err := parseDeclsSnippet(ctx.Fset, fmt.Sprintf("package p\nfunc (%s) AlwaysShouldTraverse_() {}\n", recvType))
				if err != nil {
					return totalSites, err
				}
				file.Decls = append(file.Decls, decls...)
				fileSites++
			}
			if cand.noOpVisit &&
				!hasMethodNamed(cand.ptrT, ctx.Pkg, "NoOpVisit_") &&
				!fileHasPtrMethod(file, cand.ts.Name.Name, "NoOpVisit_") {
				decls, err := parseDeclsSnippet(ctx.Fset, fmt.Sprintf("package p\nfunc (%s) NoOpVisit_() {}\n", recvType))
				if err != nil {
					return totalSites, err
				}
				file.Decls = append(file.Decls, decls...)
				fileSites++
			}
			if cand.canRelScope &&
				!hasMethodNamed(cand.ptrT, ctx.Pkg, "CanReleaseSyntacticScopes_") &&
				!fileHasPtrMethod(file, cand.ts.Name.Name, "CanReleaseSyntacticScopes_") {
				decls, err := parseDeclsSnippet(ctx.Fset, fmt.Sprintf("package p\nfunc (%s) CanReleaseSyntacticScopes_() {}\n", recvType))
				if err != nil {
					return totalSites, err
				}
				file.Decls = append(file.Decls, decls...)
				fileSites++
			}
			if cand.canRelScopeDynamic &&
				!hasMethodNamed(cand.ptrT, ctx.Pkg, "CanReleaseSyntacticScopesDynamic_") &&
				!fileHasPtrMethod(file, cand.ts.Name.Name, "CanReleaseSyntacticScopesDynamic_") {
				decls, err := parseDeclsSnippet(ctx.Fset, fmt.Sprintf(
					"package p\nfunc (this_ %s) CanReleaseSyntacticScopesDynamic_() bool {\n\tif this_ == nil || len(this_.callbacks) == 0 {\n\t\treturn false\n\t}\n\tfor _, w_ := range this_.callbacks {\n\t\tif w_ == nil || !canReleaseSyntacticScopesForCallback_(w_.callback) {\n\t\t\treturn false\n\t\t}\n\t}\n\treturn true\n}\n",
					recvType))
				if err != nil {
					return totalSites, err
				}
				file.Decls = append(file.Decls, decls...)
				fileSites++
			}
			if cand.unwrapCapField != "" &&
				!hasMethodNamed(cand.ptrT, ctx.Pkg, "UnwrapCallback_") &&
				!fileHasPtrMethod(file, cand.ts.Name.Name, "UnwrapCallback_") {
				decls, err := parseDeclsSnippet(ctx.Fset, fmt.Sprintf(
					"package p\nfunc (this_ %s) UnwrapCallback_() any {\n\tif this_ == nil {\n\t\treturn nil\n\t}\n\treturn this_.%s\n}\n",
					recvType, cand.unwrapCapField))
				if err != nil {
					return totalSites, err
				}
				file.Decls = append(file.Decls, decls...)
				fileSites++
			}
		}
		if fileSites > 0 {
			ctx.MarkFileModified(file)
			totalSites += fileSites
		}
	}

	if sscHasRelease && postOrderCbNamed != nil {
		lambdaSites, err := p.rewritePostOrderCallbackLambdas(ctx, methodDecls, topFuncDecls)
		if err != nil {
			return totalSites, err
		}
		totalSites += lambdaSites
	}

	return totalSites, nil
}

func collectDirectEmbeddedNamed(named *types.Named) []*types.Named {
	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return nil
	}
	var out []*types.Named
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if !f.Embedded() {
			continue
		}
		ft := f.Type()
		if pt, ok := ft.(*types.Pointer); ok {
			ft = pt.Elem()
		}
		if emb, ok := ft.(*types.Named); ok {
			out = append(out, emb.Origin())
		}
	}
	return out
}

func formatPtrReceiverType(ts *ast.TypeSpec) string {
	if ts == nil || ts.Name == nil {
		return ""
	}
	if ts.TypeParams == nil || len(ts.TypeParams.List) == 0 {
		return "*" + ts.Name.Name
	}
	var params []string
	for _, field := range ts.TypeParams.List {
		for _, id := range field.Names {
			if id == nil || id.Name == "" {
				return ""
			}
			params = append(params, id.Name)
		}
	}
	if len(params) == 0 {
		return ""
	}
	return fmt.Sprintf("*%s[%s]", ts.Name.Name, strings.Join(params, ", "))
}

func hasMethodNamed(ptrT *types.Pointer, pkg *types.Package, methodName string) bool {
	obj, _, _ := types.LookupFieldOrMethod(ptrT, true, pkg, methodName)
	return obj != nil
}

func fileHasPtrMethod(file *ast.File, typeName, methodName string) bool {
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv == nil || len(fd.Recv.List) != 1 || fd.Name == nil || fd.Name.Name != methodName {
			continue
		}
		if recvStructTypeName(fd.Recv.List[0].Type) == typeName {
			return true
		}
	}
	return false
}

func recvStructTypeName(expr ast.Expr) string {
	star, ok := stripParens(expr).(*ast.StarExpr)
	if !ok {
		return ""
	}
	base := stripParens(star.X)
	switch x := base.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.IndexExpr:
		if id, ok := stripParens(x.X).(*ast.Ident); ok {
			return id.Name
		}
	case *ast.IndexListExpr:
		if id, ok := stripParens(x.X).(*ast.Ident); ok {
			return id.Name
		}
	}
	return ""
}

func canAlwaysShouldTraverse(
	ptrT *types.Pointer,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
) bool {
	if ov, _, _ := types.LookupFieldOrMethod(ptrT, true, pkg, "ov_ShouldTraverse_3"); ov != nil {
		return false
	}
	sel := types.NewMethodSet(ptrT).Lookup(pkg, "ShouldTraverse_3")
	if sel == nil {
		return false
	}
	fn, ok := sel.Obj().(*types.Func)
	if !ok || fn == nil {
		return false
	}
	fd := methodDecls[fn]
	if fd == nil && fn.Origin() != nil {
		fd = methodDecls[fn.Origin()]
	}
	if fd == nil || fd.Body == nil {
		return false
	}
	return isPureReturnTrueBody(fd.Body, info)
}

func canNoOpVisit(
	ptrT *types.Pointer,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
) bool {
	if ov, _, _ := types.LookupFieldOrMethod(ptrT, true, pkg, "ov_Visit_3"); ov != nil {
		return false
	}
	sel := types.NewMethodSet(ptrT).Lookup(pkg, "Visit_3")
	if sel == nil {
		return false
	}
	fn, ok := sel.Obj().(*types.Func)
	if !ok || fn == nil {
		return false
	}
	fd := methodDecls[fn]
	if fd == nil && fn.Origin() != nil {
		fd = methodDecls[fn.Origin()]
	}
	if fd == nil || fd.Body == nil {
		return false
	}
	return isPureNoOpBody(fd.Body, info)
}

func bodyHasOverrideOrSelfRef(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		if id, ok := n.(*ast.Ident); ok {
			if id.Name == "self_" || strings.HasPrefix(id.Name, "ov_") {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

func isBlankVarAssignOrEmpty(stmt ast.Stmt, info *types.Info) bool {
	switch s := stmt.(type) {
	case *ast.EmptyStmt:
		return true
	case *ast.AssignStmt:
		if s.Tok != token.ASSIGN || len(s.Lhs) == 0 || len(s.Lhs) != len(s.Rhs) {
			return false
		}
		for i := range s.Lhs {
			lhsID, ok := stripParens(s.Lhs[i]).(*ast.Ident)
			if !ok || lhsID.Name != "_" {
				return false
			}
			rhsID, ok := stripParens(s.Rhs[i]).(*ast.Ident)
			if !ok {
				return false
			}
			if info != nil && info.Uses != nil {
				v, ok := info.Uses[rhsID].(*types.Var)
				if !ok || v == nil {
					return false
				}
				if v.Pkg() != nil && v.Parent() == v.Pkg().Scope() {
					return false
				}
			}
		}
		return true
	default:
		return false
	}
}

func isPureReturnTrueBody(body *ast.BlockStmt, info *types.Info) bool {
	if body == nil || len(body.List) == 0 || bodyHasOverrideOrSelfRef(body) {
		return false
	}
	for _, stmt := range body.List[:len(body.List)-1] {
		if !isBlankVarAssignOrEmpty(stmt, info) {
			return false
		}
	}
	ret, ok := body.List[len(body.List)-1].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return false
	}
	id, ok := stripParens(ret.Results[0]).(*ast.Ident)
	if !ok || id.Name != "true" {
		return false
	}
	return info != nil && info.Uses != nil && info.Uses[id] == types.Universe.Lookup("true")
}

func isPureNoOpBody(body *ast.BlockStmt, info *types.Info) bool {
	if body == nil || bodyHasOverrideOrSelfRef(body) {
		return false
	}
	for i, stmt := range body.List {
		if i == len(body.List)-1 {
			if ret, ok := stmt.(*ast.ReturnStmt); ok && len(ret.Results) == 0 {
				continue
			}
		}
		if !isBlankVarAssignOrEmpty(stmt, info) {
			return false
		}
	}
	return true
}

func canNoOpPostOrderInterfaceVisit(
	ptrT *types.Pointer,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
) bool {
	if ov, _, _ := types.LookupFieldOrMethod(ptrT, true, pkg, "ov_Visit"); ov != nil {
		return false
	}
	sel := types.NewMethodSet(ptrT).Lookup(pkg, "Visit")
	if sel == nil {
		return false
	}
	fn, ok := sel.Obj().(*types.Func)
	if !ok || fn == nil {
		return false
	}
	fd := methodDecls[fn]
	if fd == nil && fn.Origin() != nil {
		fd = methodDecls[fn.Origin()]
	}
	if fd == nil || fd.Body == nil {
		return false
	}
	return isPureNoOpBody(fd.Body, info)
}

func isSynthesizedPostOrderLambdaName(name string) bool {
	switch name {
	case "lambda_NodeTraversal_AbstractPostOrderCallbackInterface",
		"releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface",
		"releasableNoOpLambda_NodeTraversal_AbstractPostOrderCallbackInterface":
		return true
	}
	return false
}

func matchPostOrderCallbackAdapterStruct(
	instNamed *types.Named,
	ptrT *types.Pointer,
	postOrderCbNamed *types.Named,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
) string {
	if instNamed == nil || ptrT == nil || postOrderCbNamed == nil || instNamed.Obj() == nil {
		return ""
	}
	if !strings.HasPrefix(instNamed.Obj().Name(), "NodeTraversal__Anon_") {
		return ""
	}
	if structHasAnyOvHookField(instNamed, make(map[*types.Named]bool)) {
		return ""
	}
	st, ok := instNamed.Underlying().(*types.Struct)
	if !ok {
		return ""
	}
	var capField string
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if f.Embedded() {
			embT := types.Unalias(f.Type())
			if pt, isP := embT.(*types.Pointer); isP {
				embT = types.Unalias(pt.Elem())
			}
			embN, isN := embT.(*types.Named)
			if !isN || embN.Obj() == nil || embN.Obj().Name() != "NodeTraversal_AbstractPostOrderCallbackBase" {
				return ""
			}
			continue
		}
		if f.Name() == "outer_" || f.Name() == "idHash_" {
			continue
		}
		if fn, isN := types.Unalias(f.Type()).(*types.Named); isN && fn.Origin() == postOrderCbNamed.Origin() && capField == "" {
			capField = f.Name()
			continue
		}
		return ""
	}
	if capField == "" {
		return ""
	}
	sel := types.NewMethodSet(ptrT).Lookup(pkg, "Visit_3")
	if sel == nil {
		return ""
	}
	fn, ok := sel.Obj().(*types.Func)
	if !ok || fn == nil {
		return ""
	}
	fd := methodDecls[fn]
	if fd == nil && fn.Origin() != nil {
		fd = methodDecls[fn.Origin()]
	}
	if fd == nil || fd.Body == nil || len(fd.Body.List) == 0 {
		return ""
	}
	for _, stmt := range fd.Body.List[:len(fd.Body.List)-1] {
		if !isBlankVarAssignOrEmpty(stmt, info) {
			return ""
		}
	}
	es, ok := fd.Body.List[len(fd.Body.List)-1].(*ast.ExprStmt)
	if !ok {
		return ""
	}
	call, ok := stripParens(es.X).(*ast.CallExpr)
	if !ok || len(call.Args) != 3 {
		return ""
	}
	visitSel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || visitSel.Sel == nil || visitSel.Sel.Name != "Visit" {
		return ""
	}
	capSel, ok := stripParens(visitSel.X).(*ast.SelectorExpr)
	if !ok || capSel.Sel == nil || capSel.Sel.Name != capField {
		return ""
	}
	recvName := singleParamOrRecvName(fd.Recv)
	if recvName == "" || !isIdentNamed(capSel.X, recvName) {
		return ""
	}
	return capField
}

func matchCombinedCompilerPassStruct(
	instNamed *types.Named,
	callbackIface *types.Interface,
	pkg *types.Package,
	pkgScopeVarFreeIfaces map[*types.Named]bool,
) bool {
	if instNamed == nil || callbackIface == nil || pkg == nil || instNamed.Obj() == nil {
		return false
	}
	if instNamed.Obj().Pkg() != pkg || instNamed.Obj().Name() != "CombinedCompilerPass" {
		return false
	}
	if tp := instNamed.TypeParams(); tp != nil && tp.Len() > 0 {
		return false
	}
	if structHasAnyOvHookField(instNamed, make(map[*types.Named]bool)) {
		return false
	}
	st, ok := instNamed.Underlying().(*types.Struct)
	if !ok {
		return false
	}
	foundCallbacks := false
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if f.Embedded() {
			return false
		}
		if f.Name() == "callbacks" {
			sl, isSl := types.Unalias(f.Type()).(*types.Slice)
			if !isSl {
				return false
			}
			pt, isPtr := types.Unalias(sl.Elem()).(*types.Pointer)
			if !isPtr {
				return false
			}
			wrapNamed, isN := types.Unalias(pt.Elem()).(*types.Named)
			if !isN || wrapNamed.Obj() == nil || wrapNamed.Obj().Pkg() != pkg ||
				wrapNamed.Obj().Name() != "CombinedCompilerPass_CallbackWrapper" {
				return false
			}
			if tp := wrapNamed.TypeParams(); tp != nil && tp.Len() > 0 {
				return false
			}
			if structHasAnyOvHookField(wrapNamed, make(map[*types.Named]bool)) {
				return false
			}
			wrapSt, isWS := wrapNamed.Underlying().(*types.Struct)
			if !isWS {
				return false
			}
			foundInnerCb := false
			for j := 0; j < wrapSt.NumFields(); j++ {
				wf := wrapSt.Field(j)
				if wf.Embedded() {
					return false
				}
				if wf.Name() == "callback" {
					cbNamed, isCN := types.Unalias(wf.Type()).(*types.Named)
					if !isCN || cbNamed.Obj() == nil || cbNamed.Obj().Pkg() != pkg ||
						cbNamed.Obj().Name() != "NodeTraversal_Callback" || !types.Identical(cbNamed.Underlying(), callbackIface) {
						return false
					}
					foundInnerCb = true
					continue
				}
				if wf.Name() == "scopedCallback" {
					scNamed, isSN := types.Unalias(wf.Type()).(*types.Named)
					if isSN && scNamed.Obj() != nil && scNamed.Obj().Pkg() == pkg &&
						scNamed.Obj().Name() == "NodeTraversal_ScopedCallback" {
						continue
					}
				}
				if !isScopeVarFreeType(wf.Type(), pkg, make(map[*types.Named]bool), 0, pkgScopeVarFreeIfaces) {
					return false
				}
			}
			if !foundInnerCb {
				return false
			}
			foundCallbacks = true
			continue
		}
		if !isScopeVarFreeType(f.Type(), pkg, make(map[*types.Named]bool), 0, pkgScopeVarFreeIfaces) {
			return false
		}
	}
	return foundCallbacks
}

func computePkgScopeVarFreeInterfaces(pkg *types.Package, concreteStructs []*types.Named) map[*types.Named]bool {
	out := make(map[*types.Named]bool)
	if pkg == nil || pkg.Scope() == nil {
		return out
	}
	sc := pkg.Scope()
	for _, name := range sc.Names() {
		tn, ok := sc.Lookup(name).(*types.TypeName)
		if !ok || tn == nil || tn.IsAlias() {
			continue
		}
		ifaceNamed, ok := tn.Type().(*types.Named)
		if !ok || ifaceNamed == nil || (ifaceNamed.TypeParams() != nil && ifaceNamed.TypeParams().Len() > 0) {
			continue
		}
		iface, ok := ifaceNamed.Underlying().(*types.Interface)
		if !ok || iface.NumMethods() == 0 {
			continue
		}
		// Require a j2go package marker method Implements_<pkg>_<Name>() and verify no method returns Scope/Var/any.
		wantMarker := "Implements_" + pkg.Name() + "_" + name
		hasMarker := false
		returnsSuspicious := false
		for m := 0; m < iface.NumMethods(); m++ {
			fn := iface.Method(m)
			if fn.Name() == wantMarker {
				hasMarker = true
			}
			if sig, ok := fn.Type().(*types.Signature); ok && sig.Results() != nil {
				for r := 0; r < sig.Results().Len(); r++ {
					rt := types.Unalias(sig.Results().At(r).Type())
					if _, isIf := rt.(*types.Interface); isIf {
						if rn, isN := rt.(*types.Named); !isN || rn.Origin() != ifaceNamed.Origin() {
							returnsSuspicious = true
						}
					} else if !isScopeVarFreeType(rt, pkg, make(map[*types.Named]bool), 0, nil) {
						returnsSuspicious = true
					}
				}
			}
		}
		if !hasMarker || returnsSuspicious {
			continue
		}
		implCount := 0
		allImplsFree := true
		for _, cs := range concreteStructs {
			if cs == nil || (cs.TypeParams() != nil && cs.TypeParams().Len() > 0) {
				continue
			}
			csPtr := types.NewPointer(cs)
			if !types.Implements(csPtr, iface) {
				continue
			}
			implCount++
			visiting := map[*types.Named]bool{cs.Origin(): true}
			if !checkStructFieldsScopeVarFree(cs, csPtr, true, pkg, visiting, 0, nil) {
				allImplsFree = false
				break
			}
		}
		if implCount > 0 && allImplsFree {
			out[ifaceNamed.Origin()] = true
		}
	}
	return out
}

func canReleaseSyntacticScopes(
	ptrT *types.Pointer,
	instNamed *types.Named,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
	topFuncDecls map[*types.Func]*ast.FuncDecl,
	pkgScopeVarFreeIfaces map[*types.Named]bool,
) bool {
	if ptrT == nil || instNamed == nil || pkg == nil || info == nil {
		return false
	}
	if structHasAnyOvHookField(instNamed, make(map[*types.Named]bool)) {
		return false
	}
	// Fast path 1: None of the methods on *S ever invoke a scope-producing method on *NodeTraversal
	// or allow *NodeTraversal to escape to external functions/fields.
	if callbackNeverTouchesSyntacticScope(ptrT, pkg, info, methodDecls) {
		return true
	}
	// Fast path 1b: Methods on *S only query t.GetScope() / Scope / Var ephemerally (read-only primitive/Node
	// queries, nil checks, or local GetVarIterable() cursor loops) without storing or escaping *Scope or *Var.
	if callbackUsesScopeAndVarEphemerally(ptrT, instNamed, pkg, info, methodDecls, topFuncDecls) {
		return callbackAndOuterMethodsNoGlobalEscape(ptrT, instNamed, pkg, info, methodDecls)
	}
	// Fast path 2: Every reachable field on S (including embedded base structs and outer_ chain)
	// has a go/types type incapable of storing *Scope or *Var, and no method on S (or its outer_
	// chain) spawns goroutines or assigns to package-level variables.
	if isScopeVarFreeCallbackStruct(instNamed, ptrT, pkg, pkgScopeVarFreeIfaces) {
		return callbackAndOuterMethodsNoGlobalEscape(ptrT, instNamed, pkg, info, methodDecls)
	}
	// Fast path 3: Pass struct has transient per-scope fields that are provably nil-ed via `defer`
	// before the scope method returns, and all remaining fields are scope/var-free.
	if callbackClearsTransientScopeFieldsOnExit(ptrT, instNamed, pkg, info, methodDecls, pkgScopeVarFreeIfaces) {
		return callbackAndOuterMethodsNoGlobalEscape(ptrT, instNamed, pkg, info, methodDecls)
	}
	return false
}

func structHasAnyOvHookField(named *types.Named, seen map[*types.Named]bool) bool {
	if named == nil {
		return false
	}
	orig := named.Origin()
	if seen[orig] {
		return false
	}
	seen[orig] = true
	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return false
	}
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if strings.HasPrefix(f.Name(), "ov_") {
			return true
		}
		if f.Embedded() {
			ft := types.Unalias(f.Type())
			if pt, ok := ft.(*types.Pointer); ok {
				ft = types.Unalias(pt.Elem())
			}
			if emb, ok := ft.(*types.Named); ok && structHasAnyOvHookField(emb, seen) {
				return true
			}
		}
	}
	return false
}

func isScopeFreeNodeTraversalMethod(name string) bool {
	switch name {
	case "GetCompiler", "GetSourceName", "GetInput", "GetInputId", "GetChunk",
		"GetCurrentNode", "GetCurrentScript", "GetEnclosingFunction",
		"GetScopeRoot", "GetScopeDepth", "IsHoistScope",
		"GetClosestHoistScopeRoot", "GetClosestScopeRootNodeBindingThisOrSuper",
		"InGlobalScope", "InModuleScope", "InGlobalOrModuleScope",
		"InFunctionBlockScope", "InGlobalHoistScope", "InModuleHoistScope",
		"Report", "Report_3", "Report_4",
		"ReportCodeChange", "ReportCodeChange_0", "ReportCodeChange_1":
		return true
	}
	return false
}

func isPtrToNodeTraversalType(t types.Type, pkg *types.Package) bool {
	if t == nil || pkg == nil {
		return false
	}
	pt, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := types.Unalias(pt.Elem()).(*types.Named)
	return ok && named.Obj() != nil && named.Obj().Pkg() == pkg && named.Obj().Name() == "NodeTraversal"
}

func callbackNeverTouchesSyntacticScope(
	ptrT *types.Pointer,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
) bool {
	ms := types.NewMethodSet(ptrT)
	if ms.Len() == 0 {
		return false
	}
	for i := 0; i < ms.Len(); i++ {
		fn, ok := ms.At(i).Obj().(*types.Func)
		if !ok || fn == nil {
			continue
		}
		if strings.HasPrefix(fn.Name(), "Implements_") ||
			fn.Name() == "BindSelf_" || fn.Name() == "IsJavaNil_" ||
			fn.Name() == "JavaClassInfo_" || fn.Name() == "IdentityHash_" ||
			fn.Name() == "HasIdentityEquals_" || fn.Name() == "SetOuter_" || fn.Name() == "GetOuter_" ||
			fn.Name() == "AlwaysShouldTraverse_" || fn.Name() == "NoOpVisit_" || fn.Name() == "CanReleaseSyntacticScopes_" ||
			strings.HasPrefix(fn.Name(), "GetNodeTraversal_") {
			continue
		}
		fd := methodDecls[fn]
		if fd == nil && fn.Origin() != nil {
			fd = methodDecls[fn.Origin()]
		}
		if fd == nil || fd.Body == nil {
			return false
		}
		if !methodUsesNodeTraversalScopeFree(fd, ptrT, pkg, info) {
			return false
		}
	}
	return true
}

func methodUsesNodeTraversalScopeFree(
	fd *ast.FuncDecl,
	cbPtrT *types.Pointer,
	pkg *types.Package,
	info *types.Info,
) bool {
	var recvVar *types.Var
	if fd.Recv != nil && len(fd.Recv.List) == 1 && len(fd.Recv.List[0].Names) == 1 {
		recvVar = varFromObj(info.Defs[fd.Recv.List[0].Names[0]])
	}
	ntParams := make(map[*types.Var]bool)
	if fd.Type != nil && fd.Type.Params != nil {
		for _, fld := range fd.Type.Params.List {
			for _, id := range fld.Names {
				if v := varFromObj(info.Defs[id]); v != nil && isPtrToNodeTraversalType(v.Type(), pkg) {
					ntParams[v] = true
				}
			}
		}
	}
	if len(ntParams) == 0 {
		return true
	}
	parents, _ := buildParentMap(fd.Body)
	okAll := true
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		if !okAll || n == nil {
			return false
		}
		id, isID := n.(*ast.Ident)
		if !isID {
			return true
		}
		v := varFromObj(info.Uses[id])
		if v == nil || !ntParams[v] {
			return true
		}
		// Walk up enclosing ParenExprs.
		var expr ast.Node = id
		par := parents[expr]
		for {
			pe, isParen := par.(*ast.ParenExpr)
			if !isParen || pe.X != expr {
				break
			}
			expr = pe
			par = parents[pe]
		}
		// Case 1: `_ = t`
		if as, isAs := par.(*ast.AssignStmt); isAs && as.Tok == token.ASSIGN && len(as.Lhs) == len(as.Rhs) {
			for idx, rhs := range as.Rhs {
				if rhs == expr && isIdentNamed(as.Lhs[idx], "_") {
					return true
				}
			}
			okAll = false
			return false
		}
		// Case 2: `t.<ScopeFreeMethod>(...)`
		if sel, isSel := par.(*ast.SelectorExpr); isSel && sel.X == expr && sel.Sel != nil {
			selExpr := ast.Node(sel)
			selPar := parents[selExpr]
			for {
				pe, isParen := selPar.(*ast.ParenExpr)
				if !isParen || pe.X != selExpr {
					break
				}
				selExpr = pe
				selPar = parents[pe]
			}
			if call, isCall := selPar.(*ast.CallExpr); isCall && call.Fun == selExpr && isScopeFreeNodeTraversalMethod(sel.Sel.Name) {
				return true
			}
			okAll = false
			return false
		}
		// Case 3: passed to `this_.<SubMethod>(...)` or `this_.self_.<SubMethod>(...)` on cbPtrT (also checked in MethodSet!)
		if call, isCall := par.(*ast.CallExpr); isCall && call.Fun != expr {
			if recvVar != nil && isSelfOrThisMethodCallOnCallback(call, recvVar, cbPtrT, pkg, info) {
				return true
			}
		}
		okAll = false
		return false
	})
	return okAll
}

func isSelfOrThisMethodCallOnCallback(
	call *ast.CallExpr,
	recvVar *types.Var,
	cbPtrT *types.Pointer,
	pkg *types.Package,
	info *types.Info,
) bool {
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel == nil {
		return false
	}
	target := stripParens(sel.X)
	isRecvOrSelf := false
	if id, ok := target.(*ast.Ident); ok && varFromObj(info.Uses[id]) == recvVar {
		isRecvOrSelf = true
	} else if innerSel, ok := target.(*ast.SelectorExpr); ok && innerSel.Sel != nil && innerSel.Sel.Name == "self_" {
		if id, ok := stripParens(innerSel.X).(*ast.Ident); ok && varFromObj(info.Uses[id]) == recvVar {
			isRecvOrSelf = true
		}
	}
	if !isRecvOrSelf {
		return false
	}
	return types.NewMethodSet(cbPtrT).Lookup(pkg, sel.Sel.Name) != nil
}

func isScopeVarFreeCallbackStruct(
	instNamed *types.Named,
	cbPtrT *types.Pointer,
	pkg *types.Package,
	pkgScopeVarFreeIfaces map[*types.Named]bool,
) bool {
	// Generic callback structs instantiated over type parameters could hold *Scope/*Var in T.
	if tp := instNamed.TypeParams(); tp != nil && tp.Len() > 0 {
		return false
	}
	visiting := map[*types.Named]bool{instNamed.Origin(): true}
	return checkStructFieldsScopeVarFree(instNamed, cbPtrT, true, pkg, visiting, 0, pkgScopeVarFreeIfaces)
}

func checkStructFieldsScopeVarFree(
	named *types.Named,
	cbPtrT *types.Pointer,
	allowEmbeddedSelf bool,
	pkg *types.Package,
	visiting map[*types.Named]bool,
	depth int,
	pkgScopeVarFreeIfaces map[*types.Named]bool,
) bool {
	if depth > 4 || named == nil {
		return false
	}
	if orig := named.Origin(); orig != nil && orig.Obj() != nil &&
		!structHasAnyOvHookField(named, make(map[*types.Named]bool)) &&
		isKnownScopeVarFreeStruct(orig.Obj().Pkg(), orig.Obj().Name(), pkg) {
		return true
	}
	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return false
	}
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		ft := f.Type()
		if allowEmbeddedSelf && f.Name() == "self_" {
			if iface, isIf := ft.Underlying().(*types.Interface); isIf && cbPtrT != nil && types.Implements(cbPtrT, iface) {
				continue
			}
		}
		if allowEmbeddedSelf && f.Embedded() {
			embT := types.Unalias(ft)
			if pt, isPtr := embT.(*types.Pointer); isPtr {
				embT = types.Unalias(pt.Elem())
			}
			if embNamed, isN := embT.(*types.Named); isN {
				orig := embNamed.Origin()
				if visiting[orig] {
					continue
				}
				visiting[orig] = true
				okEmb := checkStructFieldsScopeVarFree(embNamed, cbPtrT, true, pkg, visiting, depth+1, pkgScopeVarFreeIfaces)
				delete(visiting, orig)
				if !okEmb {
					return false
				}
				continue
			}
		}
		if !isScopeVarFreeType(ft, pkg, visiting, depth, pkgScopeVarFreeIfaces) {
			return false
		}
	}
	return true
}

func isScopeVarFreeType(
	t types.Type,
	pkg *types.Package,
	visiting map[*types.Named]bool,
	depth int,
	pkgScopeVarFreeIfaces map[*types.Named]bool,
) bool {
	if t == nil || depth > 4 {
		return false
	}
	ut := types.Unalias(t)
	switch x := ut.(type) {
	case *types.Basic:
		return x.Kind() != types.Invalid && x.Kind() != types.UnsafePointer
	case *types.Pointer:
		return isScopeVarFreeType(x.Elem(), pkg, visiting, depth, pkgScopeVarFreeIfaces)
	case *types.Slice:
		return isScopeVarFreeType(x.Elem(), pkg, visiting, depth, pkgScopeVarFreeIfaces)
	case *types.Array:
		return isScopeVarFreeType(x.Elem(), pkg, visiting, depth, pkgScopeVarFreeIfaces)
	case *types.Map:
		return isScopeVarFreeType(x.Key(), pkg, visiting, depth, pkgScopeVarFreeIfaces) &&
			isScopeVarFreeType(x.Elem(), pkg, visiting, depth, pkgScopeVarFreeIfaces)
	case *types.Signature, *types.Chan, *types.TypeParam:
		return false
	case *types.Interface:
		return false
	case *types.Named:
		orig := x.Origin()
		if orig == nil || orig.Obj() == nil {
			return false
		}
		objPkg := orig.Obj().Pkg()
		name := orig.Obj().Name()

		switch under := x.Underlying().(type) {
		case *types.Basic:
			return under.Kind() != types.Invalid && under.Kind() != types.UnsafePointer
		case *types.Interface:
			if objPkg == pkg && name == "ControlFlowGraph" {
				targs := x.TypeArgs()
				return targs != nil && targs.Len() == 1 && isScopeVarFreeType(targs.At(0), pkg, visiting, depth+1, pkgScopeVarFreeIfaces)
			}
			if pkgScopeVarFreeIfaces != nil && pkgScopeVarFreeIfaces[orig] {
				return true
			}
			return isKnownScopeVarFreeInterface(objPkg, name, pkg)
		case *types.Struct:
			if objPkg != nil {
				pkgPath := objPkg.Path()
				if pkgPath == "sync" {
					switch name {
					case "Mutex", "RWMutex", "Once", "WaitGroup":
						return true
					default:
						return false
					}
				}
				if pkgPath == "sync/atomic" {
					switch name {
					case "Int32", "Int64", "Uint32", "Uint64", "Uintptr", "Bool":
						return true
					default:
						return false
					}
				}
				if isJ2gortPkg(objPkg) {
					switch name {
					case "OrderedMap", "IdentityOrderedMap", "CustomOrderedMap",
						"OrderedSet", "IdentityOrderedSet", "CustomOrderedSet",
						"Deque", "Slab", "PooledSlab":
						targs := x.TypeArgs()
						if targs == nil || targs.Len() == 0 {
							return false
						}
						for i := 0; i < targs.Len(); i++ {
							if !isScopeVarFreeType(targs.At(i), pkg, visiting, depth, pkgScopeVarFreeIfaces) {
								return false
							}
						}
						return true
					case "ClinitState", "LocalStringBuilder", "BitSet", "EnumBase",
						"BoxBoolean", "BoxInt", "BoxLong", "BoxDouble", "BoxFloat", "BoxShort", "BoxByte", "BoxChar":
						return true
					default:
						return false
					}
				}
			}

			if isDeniedScopeOrVarStruct(objPkg, name, pkg) {
				return false
			}
			if isKnownScopeVarFreeStruct(objPkg, name, pkg) {
				return true
			}
			if tp := orig.TypeParams(); tp != nil && tp.Len() > 0 && x.TypeArgs() == nil {
				return false
			}
			if visiting[orig] {
				return true
			}
			visiting[orig] = true
			okFields := checkStructFieldsScopeVarFree(x, nil, false, pkg, visiting, depth+1, pkgScopeVarFreeIfaces)
			delete(visiting, orig)
			return okFields
		default:
			return false
		}
	default:
		return false
	}
}

type ephemValKind uint8

const (
	ephemNone ephemValKind = iota
	ephemNT
	ephemScope
	ephemVar
	ephemVarIterable
	ephemIterCursor
	ephemIterRawElem
	ephemInvalid
)

func isScopeTypeInPkg(t types.Type, pkg *types.Package) bool {
	if t == nil || pkg == nil {
		return false
	}
	ut := types.Unalias(t)
	if pt, isP := ut.(*types.Pointer); isP {
		ut = types.Unalias(pt.Elem())
	}
	nm, isN := ut.(*types.Named)
	if !isN || nm.Obj() == nil || nm.Obj().Pkg() != pkg {
		return false
	}
	switch nm.Obj().Name() {
	case "Scope", "AbstractScope", "TypedScope":
		return true
	}
	return false
}

func isVarTypeInPkg(t types.Type, pkg *types.Package) bool {
	if t == nil || pkg == nil {
		return false
	}
	ut := types.Unalias(t)
	if pt, isP := ut.(*types.Pointer); isP {
		ut = types.Unalias(pt.Elem())
	}
	nm, isN := ut.(*types.Named)
	if !isN || nm.Obj() == nil || nm.Obj().Pkg() != pkg {
		return false
	}
	switch nm.Obj().Name() {
	case "Var", "AbstractVar", "TypedVar":
		return true
	}
	return false
}

func isJ2gortIterCursorType(t types.Type) bool {
	if t == nil {
		return false
	}
	ut := types.Unalias(t)
	if pt, isP := ut.(*types.Pointer); isP {
		ut = types.Unalias(pt.Elem())
	}
	nm, isN := ut.(*types.Named)
	if !isN || nm.Obj() == nil || !isJ2gortPkg(nm.Obj().Pkg()) {
		return false
	}
	switch nm.Obj().Name() {
	case "IterCursor", "JIterCursor":
		return true
	}
	return false
}

func isScopeVarCursorTypeInPkg(t types.Type, pkg *types.Package) bool {
	if t == nil || pkg == nil {
		return false
	}
	ut := types.Unalias(t)
	if pt, isP := ut.(*types.Pointer); isP {
		ut = types.Unalias(pt.Elem())
	}
	nm, isN := ut.(*types.Named)
	if !isN || nm.Obj() == nil || nm.Obj().Pkg() != pkg {
		return false
	}
	orig := nm.Origin()
	return orig != nil && orig.Obj() != nil && orig.Obj().Name() == "scopeVarCursor_"
}

func exprMentionsScopeIterSource(e ast.Expr, ephemIterVars map[*types.Var]bool, info *types.Info) bool {
	found := false
	ast.Inspect(e, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		switch x := n.(type) {
		case *ast.SelectorExpr:
			if x.Sel != nil {
				switch x.Sel.Name {
				case "GetVarIterable", "GetAllSymbols", "iterVars_":
					found = true
					return false
				}
			}
		case *ast.Ident:
			if v := varFromObj(info.Uses[x]); v != nil && ephemIterVars[v] {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

func isScopeReadOnlyPrimitiveMethod(name string) bool {
	switch name {
	case "IsGlobal", "IsLocal", "IsBlockScope", "IsStaticBlockScope", "IsFunctionBlockScope",
		"IsFunctionScope", "IsModuleScope", "IsMemberFieldDefScope", "IsComputedFieldDefRhsScope",
		"IsCatchScope", "IsCfgRootScope", "IsHoistScope",
		"GetDepth", "GetVarCount", "GetRootNode", "GetRootNode_0",
		"HasSlot", "HasOwnSlot", "CanDeclare", "AddVarNamesTo", "IsJavaNil_":
		return true
	}
	return false
}

func isScopeToScopeMethod(name string) bool {
	switch name {
	case "GetParent", "GetParentScope", "GetClosestHoistScope", "GetClosestContainerScope",
		"GetClosestCfgRootScope", "GetGlobalScope", "Untyped":
		return true
	}
	return false
}

func isScopeToVarMethod(name string) bool {
	switch name {
	case "GetVar", "LookupVar", "GetSlot", "GetOwnSlot", "GetOwnImplicitSlot", "GetOwnDeclaredVar", "GetArgumentsVar":
		return true
	}
	return false
}

func isVarReadOnlyPrimitiveMethod(name string) bool {
	switch name {
	case "GetName", "GetNode", "GetNameNode", "GetParentNode", "GetInitialValue",
		"GetScopeRoot", "GetJSDocInfo", "GetInput", "GetInputName", "GetSourceFile", "GetIndex",
		"DeclarationType", "IsGlobal", "IsLocal", "IsExtern", "IsNoShadow",
		"IsVar", "IsLet", "IsConst", "IsParam", "IsDefaultParam",
		"IsCatch", "IsClass", "IsFunction", "IsImport", "IsImplicitGoogNamespace",
		"GetImplicitGoogNamespaceStrength", "IsGoogModuleExports",
		"IsArguments", "IsThis", "IsDeclaredOrInferredConst", "IsInferredConst",
		"IsDefine", "IsBleedingFunction", "IsJavaNil_":
		return true
	}
	return false
}

func isVarToVarMethod(name string) bool {
	return name == "GetDeclaration"
}

func callbackUsesScopeAndVarEphemerally(
	ptrT *types.Pointer,
	instNamed *types.Named,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
	topFuncDecls map[*types.Func]*ast.FuncDecl,
) bool {
	ms := types.NewMethodSet(ptrT)
	if ms.Len() == 0 {
		return false
	}
	outerNamed := findOuterStructNamed(instNamed)
	for i := 0; i < ms.Len(); i++ {
		fn, ok := ms.At(i).Obj().(*types.Func)
		if !ok || fn == nil || isCallbackBoilerplateMethodName(fn.Name()) {
			continue
		}
		fd := methodDecls[fn]
		if fd == nil && fn.Origin() != nil {
			fd = methodDecls[fn.Origin()]
		}
		if fd == nil || fd.Body == nil {
			return false
		}
		if !funcDeclUsesScopeAndVarEphemerally(fd, ptrT, outerNamed, pkg, info, methodDecls, topFuncDecls, make(map[*ast.FuncDecl]bool), 0) {
			return false
		}
	}
	return true
}

func isCallbackBoilerplateMethodName(name string) bool {
	return strings.HasPrefix(name, "Implements_") ||
		name == "BindSelf_" || name == "IsJavaNil_" ||
		name == "JavaClassInfo_" || name == "IdentityHash_" ||
		name == "HasIdentityEquals_" || name == "SetOuter_" || name == "GetOuter_" ||
		name == "AlwaysShouldTraverse_" || name == "NoOpVisit_" ||
		name == "CanReleaseSyntacticScopes_" || name == "CanReleaseSyntacticScopesDynamic_" ||
		name == "UnwrapCallback_" || name == "ReleaseSlabs_" || strings.HasPrefix(name, "GetNodeTraversal_") ||
		name == "GetGuardedCallbackBase" || name == "iterVars_"
}

func funcDeclUsesScopeAndVarEphemerally(
	fd *ast.FuncDecl,
	cbPtrT *types.Pointer,
	outerNamed *types.Named,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
	topFuncDecls map[*types.Func]*ast.FuncDecl,
	seen map[*ast.FuncDecl]bool,
	depth int,
) bool {
	if fd == nil || fd.Body == nil || depth > 6 {
		return false
	}
	if seen[fd] {
		return true
	}
	seen[fd] = true
	defer delete(seen, fd)

	var recvVar *types.Var
	if fd.Recv != nil && len(fd.Recv.List) == 1 && len(fd.Recv.List[0].Names) == 1 {
		recvVar = varFromObj(info.Defs[fd.Recv.List[0].Names[0]])
	}
	ntParams := make(map[*types.Var]bool)
	scopeParams := make(map[*types.Var]bool)
	varParams := make(map[*types.Var]bool)
	if fd.Type != nil && fd.Type.Params != nil {
		for _, fld := range fd.Type.Params.List {
			for _, id := range fld.Names {
				v := varFromObj(info.Defs[id])
				if v == nil {
					continue
				}
				switch {
				case isPtrToNodeTraversalType(v.Type(), pkg):
					ntParams[v] = true
				case isScopeTypeInPkg(v.Type(), pkg):
					scopeParams[v] = true
				case isVarTypeInPkg(v.Type(), pkg):
					varParams[v] = true
				}
			}
		}
	}
	if len(ntParams) == 0 && len(scopeParams) == 0 && len(varParams) == 0 {
		return true
	}
	if fd.Type != nil && fd.Type.Results != nil {
		for _, rf := range fd.Type.Results.List {
			if rt := info.TypeOf(rf.Type); rt != nil {
				if isScopeTypeInPkg(rt, pkg) || isVarTypeInPkg(rt, pkg) {
					return false
				}
			}
		}
	}
	return blockUsesScopeAndVarEphemerally(fd.Body, recvVar, cbPtrT, outerNamed, ntParams, scopeParams, varParams, pkg, info, methodDecls, topFuncDecls, seen, depth)
}

func blockUsesScopeAndVarEphemerally(
	body *ast.BlockStmt,
	recvVar *types.Var,
	cbPtrT *types.Pointer,
	outerNamed *types.Named,
	ntParams map[*types.Var]bool,
	scopeParams map[*types.Var]bool,
	varParams map[*types.Var]bool,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
	topFuncDecls map[*types.Func]*ast.FuncDecl,
	seen map[*ast.FuncDecl]bool,
	depth int,
) bool {
	if body == nil || depth > 6 {
		return false
	}

	// Collect candidate local variables defined in body that may hold ephemeral Scope / Var / IterCursor / rawElem.
	ephemScopeVars := make(map[*types.Var]bool)
	ephemVarVars := make(map[*types.Var]bool)
	ephemIterVars := make(map[*types.Var]bool)
	ephemRawVars := make(map[*types.Var]bool)
	// order-independent: copying tracked Scope parameter set
	for v := range scopeParams {
		ephemScopeVars[v] = true
	}
	// order-independent: copying tracked Var parameter set
	for v := range varParams {
		ephemVarVars[v] = true
	}

	type cursorWrite struct {
		v   *types.Var
		rhs ast.Expr
	}
	type nextDef struct {
		v     *types.Var
		curID *ast.Ident
	}
	candCursors := make(map[*types.Var]bool)
	var cursorWrites []cursorWrite
	var nextDefs []nextDef

	hasGo := false
	ast.Inspect(body, func(n ast.Node) bool {
		if hasGo || n == nil {
			return false
		}
		if _, isGo := n.(*ast.GoStmt); isGo {
			hasGo = true
			return false
		}
		if _, isFL := n.(*ast.FuncLit); isFL {
			// Do not register inner closure vars; any capture of outer tracked vars inside FuncLit will be caught below.
			return false
		}
		switch s := n.(type) {
		case *ast.ValueSpec:
			for i, nm := range s.Names {
				if nm == nil || nm.Name == "_" {
					continue
				}
				v := varFromObj(info.Defs[nm])
				if v == nil {
					continue
				}
				switch {
				case isScopeTypeInPkg(v.Type(), pkg):
					ephemScopeVars[v] = true
				case isVarTypeInPkg(v.Type(), pkg):
					ephemVarVars[v] = true
				case isJ2gortIterCursorType(v.Type()) || isScopeVarCursorTypeInPkg(v.Type(), pkg):
					candCursors[v] = true
					if i < len(s.Values) {
						cursorWrites = append(cursorWrites, cursorWrite{v: v, rhs: s.Values[i]})
					}
				}
			}
		case *ast.AssignStmt:
			if len(s.Lhs) == len(s.Rhs) {
				for i := range s.Lhs {
					id, ok := stripParens(s.Lhs[i]).(*ast.Ident)
					if !ok || id.Name == "_" {
						continue
					}
					var v *types.Var
					if s.Tok == token.DEFINE {
						v = varFromObj(info.Defs[id])
					} else if s.Tok == token.ASSIGN {
						v = varFromObj(info.Uses[id])
					}
					if v == nil {
						continue
					}
					if s.Tok == token.DEFINE {
						switch {
						case isScopeTypeInPkg(v.Type(), pkg):
							ephemScopeVars[v] = true
						case isVarTypeInPkg(v.Type(), pkg):
							ephemVarVars[v] = true
						case isJ2gortIterCursorType(v.Type()) || isScopeVarCursorTypeInPkg(v.Type(), pkg):
							candCursors[v] = true
							cursorWrites = append(cursorWrites, cursorWrite{v: v, rhs: s.Rhs[i]})
						default:
							if call, isC := stripParens(s.Rhs[i]).(*ast.CallExpr); isC && len(call.Args) == 0 {
								if sel, isS := stripParens(call.Fun).(*ast.SelectorExpr); isS && sel.Sel != nil && sel.Sel.Name == "Next" {
									if rId, isI := stripParens(sel.X).(*ast.Ident); isI {
										nextDefs = append(nextDefs, nextDef{v: v, curID: rId})
									}
								}
							}
						}
					} else if candCursors[v] {
						cursorWrites = append(cursorWrites, cursorWrite{v: v, rhs: s.Rhs[i]})
					}
				}
			}
		}
		return true
	})
	if hasGo {
		return false
	}

	for {
		progressed := false
		for _, cw := range cursorWrites {
			if !ephemIterVars[cw.v] && exprMentionsScopeIterSource(cw.rhs, ephemIterVars, info) {
				ephemIterVars[cw.v] = true
				progressed = true
			}
		}
		if !progressed {
			break
		}
	}
	for _, nd := range nextDefs {
		if iv := varFromObj(info.Uses[nd.curID]); iv != nil && ephemIterVars[iv] && !ephemVarVars[nd.v] {
			ephemRawVars[nd.v] = true
		}
	}

	// Ensure no tracked var is ever referenced inside a nested FuncLit.
	capturedInClosure := false
	ast.Inspect(body, func(n ast.Node) bool {
		if capturedInClosure || n == nil {
			return false
		}
		fl, isFL := n.(*ast.FuncLit)
		if !isFL {
			return true
		}
		ast.Inspect(fl.Body, func(m ast.Node) bool {
			if capturedInClosure || m == nil {
				return false
			}
			if id, isID := m.(*ast.Ident); isID {
				if v := varFromObj(info.Uses[id]); v != nil {
					if ntParams[v] || ephemScopeVars[v] || ephemVarVars[v] || ephemIterVars[v] || ephemRawVars[v] {
						capturedInClosure = true
						return false
					}
				}
			}
			return true
		})
		return false
	})
	if capturedInClosure {
		return false
	}

	memo := make(map[ast.Expr]ephemValKind)
	var eval func(ast.Expr) ephemValKind
	eval = func(e ast.Expr) ephemValKind {
		if e == nil {
			return ephemNone
		}
		if k, ok := memo[e]; ok {
			return k
		}
		var res ephemValKind
		switch x := e.(type) {
		case *ast.ParenExpr:
			res = eval(x.X)

		case *ast.Ident:
			if v := varFromObj(info.Uses[x]); v != nil {
				switch {
				case ntParams[v]:
					res = ephemNT
				case ephemScopeVars[v]:
					res = ephemScope
				case ephemVarVars[v]:
					res = ephemVar
				case ephemIterVars[v]:
					res = ephemIterCursor
				case ephemRawVars[v]:
					res = ephemIterRawElem
				default:
					res = ephemNone
				}
			} else {
				res = ephemNone
			}

		case *ast.BasicLit:
			res = ephemNone

		case *ast.UnaryExpr:
			xk := eval(x.X)
			if xk == ephemNone {
				res = ephemNone
			} else {
				res = ephemInvalid
			}

		case *ast.BinaryExpr:
			lk := eval(x.X)
			rk := eval(x.Y)
			if lk == ephemNone && rk == ephemNone {
				res = ephemNone
			} else if (x.Op == token.EQL || x.Op == token.NEQ) &&
				(((lk == ephemScope || lk == ephemVar) && ntIsNilIdent(x.Y, info)) ||
					((rk == ephemScope || rk == ephemVar) && ntIsNilIdent(x.X, info)) ||
					(lk == ephemScope && rk == ephemScope) ||
					(lk == ephemVar && rk == ephemVar)) {
				res = ephemNone
			} else {
				res = ephemInvalid
			}

		case *ast.CallExpr:
			if x.Ellipsis.IsValid() {
				res = ephemInvalid
				break
			}
			res = evalEphemeralCallExpr(x, eval, recvVar, cbPtrT, outerNamed, pkg, info, methodDecls, topFuncDecls, seen, depth)

		case *ast.SelectorExpr:
			xk := eval(x.X)
			if xk == ephemNone {
				res = ephemNone
			} else {
				// SelectorExpr on an ephemeral value is ONLY legal as Fun of a CallExpr (handled in CallExpr above without calling eval(x.Fun)).
				res = ephemInvalid
			}

		case *ast.IndexExpr:
			if eval(x.X) == ephemNone && eval(x.Index) == ephemNone {
				res = ephemNone
			} else {
				res = ephemInvalid
			}

		case *ast.SliceExpr:
			if eval(x.X) == ephemNone && eval(x.Low) == ephemNone && eval(x.High) == ephemNone && eval(x.Max) == ephemNone {
				res = ephemNone
			} else {
				res = ephemInvalid
			}

		case *ast.TypeAssertExpr:
			if eval(x.X) == ephemNone {
				res = ephemNone
			} else {
				res = ephemInvalid
			}

		case *ast.StarExpr:
			if eval(x.X) == ephemNone {
				res = ephemNone
			} else {
				res = ephemInvalid
			}

		case *ast.CompositeLit:
			res = ephemNone
			for _, elt := range x.Elts {
				if kv, isKV := elt.(*ast.KeyValueExpr); isKV {
					if eval(kv.Key) != ephemNone || eval(kv.Value) != ephemNone {
						res = ephemInvalid
						break
					}
				} else if eval(elt) != ephemNone {
					res = ephemInvalid
					break
				}
			}

		case *ast.FuncLit:
			// Already verified no tracked vars are captured inside FuncLit; check zero-arg IIFE returning primitive/Node.
			res = ephemNone

		default:
			res = ephemNone
		}
		memo[e] = res
		return res
	}

	okAll := true
	var checkStmt func(ast.Stmt)
	checkStmt = func(st ast.Stmt) {
		if !okAll || st == nil {
			return
		}
		switch s := st.(type) {
		case *ast.BlockStmt:
			for _, child := range s.List {
				checkStmt(child)
			}
		case *ast.EmptyStmt, *ast.BranchStmt:
		case *ast.ExprStmt:
			k := eval(s.X)
			if k == ephemInvalid {
				okAll = false
			}
		case *ast.ReturnStmt:
			for _, r := range s.Results {
				if eval(r) != ephemNone {
					okAll = false
					return
				}
			}
		case *ast.DeclStmt:
			gd, ok := s.Decl.(*ast.GenDecl)
			if !ok {
				return
			}
			for _, sp := range gd.Specs {
				vs, ok := sp.(*ast.ValueSpec)
				if !ok {
					continue
				}
				if len(vs.Values) == 0 {
					continue
				}
				if len(vs.Names) != len(vs.Values) {
					for _, val := range vs.Values {
						if eval(val) != ephemNone {
							okAll = false
							return
						}
					}
					continue
				}
				for i, lhs := range vs.Names {
					rk := eval(vs.Values[i])
					if rk == ephemInvalid {
						okAll = false
						return
					}
					if lhs.Name == "_" {
						continue
					}
					lv := varFromObj(info.Defs[lhs])
					if !checkEphemeralVarAssign(lv, rk, ephemScopeVars, ephemVarVars, ephemIterVars, ephemRawVars) {
						okAll = false
						return
					}
				}
			}
		case *ast.AssignStmt:
			if len(s.Lhs) != len(s.Rhs) {
				for _, lhs := range s.Lhs {
					if eval(lhs) != ephemNone {
						okAll = false
						return
					}
				}
				for _, rhs := range s.Rhs {
					if eval(rhs) != ephemNone {
						okAll = false
						return
					}
				}
				return
			}
			for i := range s.Lhs {
				rk := eval(s.Rhs[i])
				if rk == ephemInvalid {
					okAll = false
					return
				}
				lhs := stripParens(s.Lhs[i])
				if id, isID := lhs.(*ast.Ident); isID {
					if id.Name == "_" {
						continue
					}
					var lv *types.Var
					if s.Tok == token.DEFINE {
						lv = varFromObj(info.Defs[id])
					} else {
						lv = varFromObj(info.Uses[id])
					}
					if s.Tok != token.ASSIGN && s.Tok != token.DEFINE && rk != ephemNone {
						okAll = false
						return
					}
					if !checkEphemeralVarAssign(lv, rk, ephemScopeVars, ephemVarVars, ephemIterVars, ephemRawVars) {
						okAll = false
						return
					}
					continue
				}
				if eval(s.Lhs[i]) != ephemNone || rk != ephemNone {
					okAll = false
					return
				}
			}
		case *ast.IfStmt:
			checkStmt(s.Init)
			if eval(s.Cond) != ephemNone {
				okAll = false
				return
			}
			checkStmt(s.Body)
			checkStmt(s.Else)
		case *ast.ForStmt:
			checkStmt(s.Init)
			if eval(s.Cond) != ephemNone {
				okAll = false
				return
			}
			checkStmt(s.Post)
			checkStmt(s.Body)
		case *ast.RangeStmt:
			if eval(s.Key) != ephemNone || eval(s.Value) != ephemNone || eval(s.X) != ephemNone {
				okAll = false
				return
			}
			checkStmt(s.Body)
		case *ast.SwitchStmt:
			checkStmt(s.Init)
			if eval(s.Tag) != ephemNone {
				okAll = false
				return
			}
			checkStmt(s.Body)
		case *ast.TypeSwitchStmt:
			checkStmt(s.Init)
			checkStmt(s.Assign)
			checkStmt(s.Body)
		case *ast.CaseClause:
			for _, e := range s.List {
				if eval(e) != ephemNone {
					okAll = false
					return
				}
			}
			for _, sub := range s.Body {
				checkStmt(sub)
			}
		case *ast.IncDecStmt:
			if eval(s.X) != ephemNone {
				okAll = false
			}
		case *ast.LabeledStmt:
			checkStmt(s.Stmt)
		case *ast.DeferStmt:
			if eval(s.Call) != ephemNone {
				okAll = false
			}
		default:
			okAll = false
		}
	}

	checkStmt(body)
	return okAll
}

func checkEphemeralVarAssign(
	lv *types.Var,
	rk ephemValKind,
	ephemScopeVars, ephemVarVars, ephemIterVars, ephemRawVars map[*types.Var]bool,
) bool {
	if lv == nil {
		return false
	}
	switch rk {
	case ephemNone:
		return !ephemScopeVars[lv] && !ephemVarVars[lv] && !ephemIterVars[lv] && !ephemRawVars[lv]
	case ephemScope:
		return ephemScopeVars[lv]
	case ephemVar:
		return ephemVarVars[lv]
	case ephemIterCursor:
		return ephemIterVars[lv]
	case ephemIterRawElem:
		return ephemRawVars[lv]
	default:
		return false
	}
}

func ntIsNilIdent(e ast.Expr, info *types.Info) bool {
	id, ok := stripParens(e).(*ast.Ident)
	return ok && id.Name == "nil" && info.Uses[id] == types.Universe.Lookup("nil")
}

func evalEphemeralCallExpr(
	call *ast.CallExpr,
	eval func(ast.Expr) ephemValKind,
	recvVar *types.Var,
	cbPtrT *types.Pointer,
	outerNamed *types.Named,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
	topFuncDecls map[*types.Func]*ast.FuncDecl,
	seen map[*ast.FuncDecl]bool,
	depth int,
) ephemValKind {
	argsAllNone := func() bool {
		for _, a := range call.Args {
			if eval(a) != ephemNone {
				return false
			}
		}
		return true
	}

	fun := stripParens(call.Fun)

	// 1. Method call `recv.Sel(args...)` where `recv` may be ephemeral.
	if sel, isSel := fun.(*ast.SelectorExpr); isSel && sel.Sel != nil {
		// Check if sel.X is a package identifier (`j2gort.CastPtr[...]`, `j2gort.NewIterCursor`, etc.) — handled below when xk == ephemNone.
		xk := eval(sel.X)
		switch xk {
		case ephemInvalid, ephemVarIterable, ephemIterRawElem:
			return ephemInvalid
		case ephemNT:
			if isScopeFreeNodeTraversalMethod(sel.Sel.Name) && argsAllNone() {
				return ephemNone
			}
			if (sel.Sel.Name == "GetScope" || sel.Sel.Name == "GetClosestContainerScope" || sel.Sel.Name == "GetClosestHoistScope") && len(call.Args) == 0 {
				return ephemScope
			}
			return ephemInvalid
		case ephemScope:
			if isScopeReadOnlyPrimitiveMethod(sel.Sel.Name) && argsAllNone() {
				return ephemNone
			}
			if sel.Sel.Name == "HasSameContainerScope" && len(call.Args) == 1 && eval(call.Args[0]) == ephemScope {
				return ephemNone
			}
			if isScopeToScopeMethod(sel.Sel.Name) && len(call.Args) == 0 {
				return ephemScope
			}
			if isScopeToVarMethod(sel.Sel.Name) && argsAllNone() {
				return ephemVar
			}
			if (sel.Sel.Name == "GetVarIterable" || sel.Sel.Name == "GetAllSymbols") && len(call.Args) == 0 {
				return ephemVarIterable
			}
			if sel.Sel.Name == "iterVars_" && len(call.Args) == 0 {
				return ephemIterCursor
			}
			return ephemInvalid
		case ephemVar:
			if isVarReadOnlyPrimitiveMethod(sel.Sel.Name) && len(call.Args) == 0 {
				return ephemNone
			}
			if sel.Sel.Name == "GetScope" && len(call.Args) == 0 {
				return ephemScope
			}
			if isVarToVarMethod(sel.Sel.Name) && len(call.Args) == 0 {
				return ephemVar
			}
			return ephemInvalid
		case ephemIterCursor:
			if sel.Sel.Name == "HasNext" && len(call.Args) == 0 {
				return ephemNone
			}
			if sel.Sel.Name == "Next" && len(call.Args) == 0 {
				if tv, ok := info.Types[call]; ok && tv.Type != nil && isVarTypeInPkg(tv.Type, pkg) {
					return ephemVar
				}
				return ephemIterRawElem
			}
			return ephemInvalid
		}
	} else {
		// Indexed generic call like `j2gort.CastPtr[Var](...)` or `ToIface_AbstractScope(...)` or direct function call `Foo(...)`.
		switch ix := fun.(type) {
		case *ast.IndexExpr:
			if eval(ix.X) != ephemNone {
				return ephemInvalid
			}
		case *ast.IndexListExpr:
			if eval(ix.X) != ephemNone {
				return ephemInvalid
			}
		default:
			if eval(fun) != ephemNone {
				return ephemInvalid
			}
		}
	}

	if argsAllNone() {
		// Even when all args are ephemNone, ensure the call itself doesn't return a raw *Scope/*Var unless tracked!
		if tv, ok := info.Types[call]; ok && tv.Type != nil {
			if isScopeTypeInPkg(tv.Type, pkg) {
				return ephemScope
			}
			if isVarTypeInPkg(tv.Type, pkg) {
				return ephemVar
			}
		}
		return ephemNone
	}

	callee := resolveCalledTypesFunc(fun, info)
	if callee != nil && isJ2gortPkg(callee.Pkg()) && len(call.Args) == 1 {
		ak := eval(call.Args[0])
		switch callee.Name() {
		case "CheckNotNull":
			if ak == ephemVarIterable {
				return ephemVarIterable
			}
		case "NewIterCursor":
			if ak == ephemVarIterable {
				return ephemIterCursor
			}
		case "CastPtr", "AsIface", "CheckNotNullPtr", "CheckNotNullIface":
			if tv, ok := info.Types[call]; ok && tv.Type != nil {
				if isScopeTypeInPkg(tv.Type, pkg) && ak == ephemScope {
					return ephemScope
				}
				if isVarTypeInPkg(tv.Type, pkg) && (ak == ephemVar || ak == ephemIterRawElem) {
					return ephemVar
				}
			}
		}
		return ephemInvalid
	}

	if callee != nil && callee.Pkg() == pkg && len(call.Args) == 1 {
		ak := eval(call.Args[0])
		if callee.Name() == "ToIface_AbstractScope" && ak == ephemScope {
			return ephemScope
		}
		if callee.Name() == "ToIface_AbstractVar" && ak == ephemVar {
			return ephemVar
		}
	}

	// Allow forwarding `ephemNT` / `ephemScope` / `ephemVar` (with all other args `ephemNone`) to:
	//  - `this_.<SubMethod>` / `this_.self_.<SubMethod>` on cbPtrT
	//  - `this_.outer_.<OuterMethod>` on *outerNamed
	//  - package-local top-level function `Foo(...)`
	hasEphemArg := false
	for _, a := range call.Args {
		ak := eval(a)
		if ak == ephemNT || ak == ephemScope || ak == ephemVar {
			hasEphemArg = true
		} else if ak != ephemNone {
			return ephemInvalid
		}
	}
	if !hasEphemArg || callee == nil || callee.Pkg() != pkg {
		return ephemInvalid
	}
	if tv, ok := info.Types[call]; ok && tv.Type != nil {
		if isScopeTypeInPkg(tv.Type, pkg) || isVarTypeInPkg(tv.Type, pkg) {
			return ephemInvalid
		}
	}

	// Top-level package helper function:
	if callee.Type() != nil {
		if sig, ok := callee.Type().(*types.Signature); ok && sig.Recv() == nil {
			targetFD := topFuncDecls[callee]
			if targetFD == nil && callee.Origin() != nil {
				targetFD = topFuncDecls[callee.Origin()]
			}
			if targetFD != nil && funcDeclUsesScopeAndVarEphemerally(targetFD, nil, nil, pkg, info, methodDecls, topFuncDecls, seen, depth+1) {
				return ephemNone
			}
			return ephemInvalid
		}
	}

	// Method call on `this_`, `this_.self_`, or `this_.outer_`:
	sel, isSel := fun.(*ast.SelectorExpr)
	if !isSel || sel.Sel == nil || recvVar == nil {
		return ephemInvalid
	}
	if cbPtrT != nil && isSelfOrThisMethodCallOnCallback(call, recvVar, cbPtrT, pkg, info) {
		msSel := types.NewMethodSet(cbPtrT).Lookup(pkg, sel.Sel.Name)
		if msSel == nil {
			return ephemInvalid
		}
		mFn, _ := msSel.Obj().(*types.Func)
		if mFn == nil {
			return ephemInvalid
		}
		mFD := methodDecls[mFn]
		if mFD == nil && mFn.Origin() != nil {
			mFD = methodDecls[mFn.Origin()]
		}
		if mFD != nil && funcDeclUsesScopeAndVarEphemerally(mFD, cbPtrT, outerNamed, pkg, info, methodDecls, topFuncDecls, seen, depth+1) {
			return ephemNone
		}
		return ephemInvalid
	}

	// `this_.outer_.<Method>(...)` (or direct `this_.<Method>(...)` when inside a FuncLit of `*outerNamed`),
	// OR `this_.<field>.<Method>(...)` where `this_.<field>` is a concrete package-local struct pointer `*F`
	// with no override hooks and no global escape (e.g., `Normalize_NormalizeStatements` delegating to
	// `this_.makeDeclaredNamesUnique` and `this_.rewriteLogicalAssignmentOperatorsHelper`).
	var targetNamed *types.Named
	checkNoGlobalEscape := false
	if id, isID := stripParens(sel.X).(*ast.Ident); isID && varFromObj(info.Uses[id]) == recvVar && cbPtrT == nil && outerNamed != nil {
		targetNamed = outerNamed
	} else if fieldSel, isFS := stripParens(sel.X).(*ast.SelectorExpr); isFS && fieldSel.Sel != nil {
		if id, isID := stripParens(fieldSel.X).(*ast.Ident); isID && varFromObj(info.Uses[id]) == recvVar {
			if fieldSel.Sel.Name == "outer_" && outerNamed != nil {
				targetNamed = outerNamed
			} else if fieldSel.Sel.Name != "self_" {
				if tv, ok := info.Types[fieldSel]; ok && tv.Type != nil {
					if pt, isP := types.Unalias(tv.Type).(*types.Pointer); isP {
						if fnm, isN := types.Unalias(pt.Elem()).(*types.Named); isN && fnm.Obj() != nil && fnm.Obj().Pkg() == pkg &&
							(fnm.TypeParams() == nil || fnm.TypeParams().Len() == 0) {
							targetNamed = fnm
							checkNoGlobalEscape = true
						}
					}
				}
			}
		}
	}
	if targetNamed != nil && !structHasAnyOvHookField(targetNamed, make(map[*types.Named]bool)) {
		outerPtr := types.NewPointer(targetNamed)
		if !checkNoGlobalEscape || callbackAndOuterMethodsNoGlobalEscape(outerPtr, targetNamed, pkg, info, methodDecls) {
			msSel := types.NewMethodSet(outerPtr).Lookup(pkg, sel.Sel.Name)
			if msSel != nil {
				if mFn, _ := msSel.Obj().(*types.Func); mFn != nil {
					mFD := methodDecls[mFn]
					if mFD == nil && mFn.Origin() != nil {
						mFD = methodDecls[mFn.Origin()]
					}
					nextOuter := findOuterStructNamed(targetNamed)
					if mFD != nil && funcDeclUsesScopeAndVarEphemerally(mFD, outerPtr, nextOuter, pkg, info, methodDecls, topFuncDecls, seen, depth+1) {
						return ephemNone
					}
				}
			}
		}
	}
	return ephemInvalid
}

func resolveCalledTypesFunc(fun ast.Expr, info *types.Info) *types.Func {
	f := stripParens(fun)
	switch x := f.(type) {
	case *ast.IndexExpr:
		f = stripParens(x.X)
	case *ast.IndexListExpr:
		f = stripParens(x.X)
	}
	switch x := f.(type) {
	case *ast.Ident:
		fn, _ := info.Uses[x].(*types.Func)
		return fn
	case *ast.SelectorExpr:
		if x.Sel != nil {
			fn, _ := info.Uses[x.Sel].(*types.Func)
			return fn
		}
	}
	return nil
}

func callbackClearsTransientScopeFieldsOnExit(
	ptrT *types.Pointer,
	instNamed *types.Named,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
	pkgScopeVarFreeIfaces map[*types.Named]bool,
) bool {
	if instNamed == nil || instNamed.Obj() == nil || (instNamed.TypeParams() != nil && instNamed.TypeParams().Len() > 0) {
		return false
	}
	name := instNamed.Obj().Name()
	var transientFields map[string]bool
	var cleanerMethodPrefix string
	switch name {
	case "FlowSensitiveInlineVariables":
		transientFields = make(map[string]bool, 8)
		for _, cf := range []string{
			"cfg",
			"candidates",
			"reachingDef",
			"reachingUses",
			"slab_Candidate",
			"reusableReachingDef",
			"reusableReachingUses",
			"reusableEscaped",
		} {
			if ntStructHasFieldNamed(instNamed, cf) {
				transientFields[cf] = true
			}
		}
		cleanerMethodPrefix = "EnterScope"
	case "DeadAssignmentsElimination":
		transientFields = map[string]bool{
			"liveness": true,
		}
		cleanerMethodPrefix = "priv_DeadAssignmentsElimination_eliminateDeadAssignments_"
	default:
		return false
	}

	// 1. All non-transient fields on instNamed (including embedded structs) must be scope/var-free.
	st, ok := instNamed.Underlying().(*types.Struct)
	if !ok {
		return false
	}
	visiting := map[*types.Named]bool{instNamed.Origin(): true}
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if transientFields[f.Name()] {
			continue
		}
		if f.Name() == "self_" {
			if iface, isIf := f.Type().Underlying().(*types.Interface); isIf && types.Implements(ptrT, iface) {
				continue
			}
		}
		if f.Embedded() {
			embT := types.Unalias(f.Type())
			if pt, isP := embT.(*types.Pointer); isP {
				embT = types.Unalias(pt.Elem())
			}
			if embNamed, isN := embT.(*types.Named); isN {
				visiting[embNamed.Origin()] = true
				okEmb := checkStructFieldsScopeVarFree(embNamed, ptrT, true, pkg, visiting, 1, pkgScopeVarFreeIfaces)
				delete(visiting, embNamed.Origin())
				if !okEmb {
					return false
				}
				continue
			}
		}
		if !isScopeVarFreeType(f.Type(), pkg, visiting, 0, pkgScopeVarFreeIfaces) {
			return false
		}
	}

	// 2. Verify cleaner method has a top-level `defer func() { ... recv.<field> = nil ... }()` nil-ing EVERY
	// transient field before any non-nil assignment to those fields, and NO other method assigns non-nil to them!
	ms := types.NewMethodSet(ptrT)
	foundCleaner := false
	for i := 0; i < ms.Len(); i++ {
		fn, ok := ms.At(i).Obj().(*types.Func)
		if !ok || fn == nil || isCallbackBoilerplateMethodName(fn.Name()) {
			continue
		}
		fd := methodDecls[fn]
		if fd == nil && fn.Origin() != nil {
			fd = methodDecls[fn.Origin()]
		}
		if fd == nil || fd.Body == nil {
			return false
		}
		recvName := singleParamOrRecvName(fd.Recv)
		if recvName == "" {
			return false
		}
		isCleaner := fn.Name() == cleanerMethodPrefix || strings.HasPrefix(fn.Name(), cleanerMethodPrefix)
		if isCleaner {
			if !methodDeferClearsTransientFieldsBeforeWrite(fd.Body, recvName, transientFields, info) {
				return false
			}
			foundCleaner = true
		} else if strings.HasPrefix(fn.Name(), "initFields_") {
			if methodAssignsNonNilNonEmptyCtorToAnyField(fd.Body, recvName, transientFields, info) {
				return false
			}
		} else {
			if methodAssignsNonNilToAnyField(fd.Body, recvName, transientFields, info) {
				return false
			}
		}
	}
	return foundCleaner
}

func ntStructHasFieldNamed(named *types.Named, fieldName string) bool {
	if named == nil {
		return false
	}
	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return false
	}
	for i := 0; i < st.NumFields(); i++ {
		if st.Field(i).Name() == fieldName {
			return true
		}
	}
	return false
}

func deferIfStmtClearsRecvField(ifSt *ast.IfStmt, recvName string, fields map[string]bool, info *types.Info) string {
	if ifSt == nil || ifSt.Init != nil || ifSt.Else != nil || ifSt.Body == nil {
		return ""
	}
	if fieldName := matchRecvFieldNonNilGuardCond(ifSt.Cond, recvName, fields, info); fieldName != "" {
		for _, st := range ifSt.Body.List {
			if _, isRet := st.(*ast.ReturnStmt); isRet {
				return ""
			}
			if as, isAs := st.(*ast.AssignStmt); isAs && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 && ntIsNilIdent(as.Rhs[0], info) {
				if sel, isS := stripParens(as.Lhs[0]).(*ast.SelectorExpr); isS && sel.Sel != nil &&
					isIdentNamed(sel.X, recvName) && sel.Sel.Name == fieldName {
					return fieldName
				}
			}
			if fCleared := matchZeroArgClearCallOnRecvField(st, recvName, fields); fCleared == fieldName {
				return fieldName
			}
		}
	}
	// Also recognize `if !(this_.reusableEscaped.IsEmpty()) { this_.reusableEscaped.Clear() }`.
	if fieldName := matchNegatedIsEmptyGuardCond(ifSt.Cond, recvName, fields); fieldName != "" {
		for _, st := range ifSt.Body.List {
			if _, isRet := st.(*ast.ReturnStmt); isRet {
				return ""
			}
			if fCleared := matchZeroArgClearCallOnRecvField(st, recvName, fields); fCleared == fieldName {
				return fieldName
			}
		}
	}
	return ""
}

func matchNegatedIsEmptyGuardCond(cond ast.Expr, recvName string, fields map[string]bool) string {
	un, ok := stripParens(cond).(*ast.UnaryExpr)
	if !ok || un.Op != token.NOT {
		return ""
	}
	call, okC := stripParens(un.X).(*ast.CallExpr)
	if !okC || len(call.Args) != 0 {
		return ""
	}
	mSel, okM := stripParens(call.Fun).(*ast.SelectorExpr)
	if !okM || mSel.Sel == nil || mSel.Sel.Name != "IsEmpty" {
		return ""
	}
	fSel, okF := stripParens(mSel.X).(*ast.SelectorExpr)
	if !okF || fSel.Sel == nil || !isIdentNamed(fSel.X, recvName) || !fields[fSel.Sel.Name] {
		return ""
	}
	return fSel.Sel.Name
}

func matchZeroArgClearCallOnRecvField(st ast.Stmt, recvName string, fields map[string]bool) string {
	es, isEs := st.(*ast.ExprStmt)
	if !isEs {
		return ""
	}
	call, isC := stripParens(es.X).(*ast.CallExpr)
	if !isC || len(call.Args) != 0 {
		return ""
	}
	mSel, isMS := stripParens(call.Fun).(*ast.SelectorExpr)
	if !isMS || mSel.Sel == nil {
		return ""
	}
	fSel, isFS := stripParens(mSel.X).(*ast.SelectorExpr)
	if !isFS || fSel.Sel == nil || !isIdentNamed(fSel.X, recvName) || !fields[fSel.Sel.Name] {
		return ""
	}
	switch mSel.Sel.Name {
	case "Reset", "Release":
		if strings.HasPrefix(fSel.Sel.Name, "slab_") {
			return fSel.Sel.Name
		}
	case "ClearScopeState", "Clear":
		return fSel.Sel.Name
	}
	return ""
}

func matchReusableInitOrResetIfStmt(st ast.Stmt, recvName string, fields map[string]bool, info *types.Info) string {
	ifSt, ok := st.(*ast.IfStmt)
	if !ok || ifSt.Init != nil || ifSt.Body == nil || ifSt.Else == nil {
		return ""
	}
	bin, okB := stripParens(ifSt.Cond).(*ast.BinaryExpr)
	if !okB || bin.Op != token.EQL {
		return ""
	}
	var target ast.Expr
	if ntIsNilIdent(bin.Y, info) {
		target = bin.X
	} else if ntIsNilIdent(bin.X, info) {
		target = bin.Y
	}
	fSel, okF := stripParens(target).(*ast.SelectorExpr)
	if !okF || fSel.Sel == nil || !isIdentNamed(fSel.X, recvName) || !fields[fSel.Sel.Name] {
		return ""
	}
	elseBlk, okE := ifSt.Else.(*ast.BlockStmt)
	if !okE || len(ifSt.Body.List) != 1 || len(elseBlk.List) != 1 {
		return ""
	}
	as, okA := ifSt.Body.List[0].(*ast.AssignStmt)
	if !okA || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
		return ""
	}
	lhsSel, okLS := stripParens(as.Lhs[0]).(*ast.SelectorExpr)
	if !okLS || lhsSel.Sel == nil || !isIdentNamed(lhsSel.X, recvName) || lhsSel.Sel.Name != fSel.Sel.Name {
		return ""
	}
	es, okES := elseBlk.List[0].(*ast.ExprStmt)
	if !okES {
		return ""
	}
	call, okC := stripParens(es.X).(*ast.CallExpr)
	if !okC {
		return ""
	}
	rSel, okRS := stripParens(call.Fun).(*ast.SelectorExpr)
	if !okRS || rSel.Sel == nil || rSel.Sel.Name != "ResetForScope" {
		return ""
	}
	rfSel, okRFS := stripParens(rSel.X).(*ast.SelectorExpr)
	if !okRFS || rfSel.Sel == nil || !isIdentNamed(rfSel.X, recvName) || rfSel.Sel.Name != fSel.Sel.Name {
		return ""
	}
	return fSel.Sel.Name
}

func bodyAliasesReusableFieldToClearedField(body *ast.BlockStmt, recvName, reusableField, aliasField string) bool {
	foundAlias := false
	ast.Inspect(body, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
			return true
		}
		lSel, okL := stripParens(as.Lhs[0]).(*ast.SelectorExpr)
		rSel, okR := stripParens(as.Rhs[0]).(*ast.SelectorExpr)
		if okL && okR && lSel.Sel != nil && rSel.Sel != nil &&
			isIdentNamed(lSel.X, recvName) && lSel.Sel.Name == aliasField &&
			isIdentNamed(rSel.X, recvName) && rSel.Sel.Name == reusableField {
			foundAlias = true
		}
		return true
	})
	return foundAlias
}

func methodAssignsNonNilNonEmptyCtorToAnyField(body *ast.BlockStmt, recvName string, fields map[string]bool, info *types.Info) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		as, isAs := n.(*ast.AssignStmt)
		if !isAs || len(as.Lhs) != len(as.Rhs) {
			return true
		}
		for i, lhs := range as.Lhs {
			sel, isS := stripParens(lhs).(*ast.SelectorExpr)
			if !isS || sel.Sel == nil || !isIdentNamed(sel.X, recvName) || !fields[sel.Sel.Name] {
				continue
			}
			rhs := stripParens(as.Rhs[i])
			if ntIsNilIdent(rhs, info) {
				continue
			}
			// Allow zero-argument empty container constructors in initFields_* (e.g. NewExt_java_util_HashSet()).
			if call, isC := rhs.(*ast.CallExpr); isC && len(call.Args) == 0 {
				continue
			}
			found = true
			return false
		}
		return true
	})
	return found
}

func matchRecvFieldNonNilGuardCond(cond ast.Expr, recvName string, fields map[string]bool, info *types.Info) string {
	bin, ok := stripParens(cond).(*ast.BinaryExpr)
	if !ok {
		return ""
	}
	if bin.Op == token.NEQ {
		var target ast.Expr
		if ntIsNilIdent(bin.Y, info) {
			target = bin.X
		} else if ntIsNilIdent(bin.X, info) {
			target = bin.Y
		}
		if sel, isS := stripParens(target).(*ast.SelectorExpr); isS && sel.Sel != nil &&
			isIdentNamed(sel.X, recvName) && fields[sel.Sel.Name] {
			return sel.Sel.Name
		}
		return ""
	}
	if bin.Op == token.LAND {
		fn := matchRecvFieldNonNilGuardCond(bin.X, recvName, fields, info)
		if fn == "" {
			return ""
		}
		un, isU := stripParens(bin.Y).(*ast.UnaryExpr)
		if !isU || un.Op != token.NOT {
			return ""
		}
		call, isC := stripParens(un.X).(*ast.CallExpr)
		if !isC || len(call.Args) != 0 {
			return ""
		}
		sel, isS := stripParens(call.Fun).(*ast.SelectorExpr)
		if !isS || sel.Sel == nil || sel.Sel.Name != "IsJavaNil_" {
			return ""
		}
		fSel, isFS := stripParens(sel.X).(*ast.SelectorExpr)
		if !isFS || fSel.Sel == nil || !isIdentNamed(fSel.X, recvName) || fSel.Sel.Name != fn {
			return ""
		}
		return fn
	}
	return ""
}

func methodDeferClearsTransientFieldsBeforeWrite(
	body *ast.BlockStmt,
	recvName string,
	fields map[string]bool,
	info *types.Info,
) bool {
	if body == nil {
		return false
	}
	clearedInDefer := make(map[string]bool, len(fields))
	clearedViaScopeState := make(map[string]bool, len(fields))
	for _, st := range body.List {
		if ds, isDef := st.(*ast.DeferStmt); isDef && ds.Call != nil {
			if fl, isFL := stripParens(ds.Call.Fun).(*ast.FuncLit); isFL && fl.Body != nil && len(ds.Call.Args) == 0 {
				for _, dst := range fl.Body.List {
					if ifSt, isIf := dst.(*ast.IfStmt); isIf {
						if fName := deferIfStmtClearsRecvField(ifSt, recvName, fields, info); fName != "" {
							clearedInDefer[fName] = true
							for _, bst := range ifSt.Body.List {
								if fc := matchZeroArgClearCallOnRecvField(bst, recvName, fields); fc == fName {
									clearedViaScopeState[fName] = true
								}
							}
						}
						continue
					}
					if fCleared := matchZeroArgClearCallOnRecvField(dst, recvName, fields); fCleared != "" {
						clearedInDefer[fCleared] = true
						clearedViaScopeState[fCleared] = true
						continue
					}
					as, isAs := dst.(*ast.AssignStmt)
					if !isAs || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 || !ntIsNilIdent(as.Rhs[0], info) {
						continue
					}
					if sel, isS := stripParens(as.Lhs[0]).(*ast.SelectorExpr); isS && sel.Sel != nil &&
						isIdentNamed(sel.X, recvName) && fields[sel.Sel.Name] {
						clearedInDefer[sel.Sel.Name] = true
					}
				}
				// Patch 86 reusable alias fields (FlowSensitiveInlineVariables):
				// `reachingDef` aliases `reusableReachingDef` and calls `reachingDef.ClearScopeState()` on exit;
				// `reachingUses` aliases `reusableReachingUses` and calls `reachingUses.ClearScopeState()` on exit.
				if fields["reusableReachingDef"] && clearedViaScopeState["reachingDef"] &&
					bodyAliasesReusableFieldToClearedField(body, recvName, "reusableReachingDef", "reachingDef") {
					clearedInDefer["reusableReachingDef"] = true
				}
				if fields["reusableReachingUses"] && clearedViaScopeState["reachingUses"] &&
					bodyAliasesReusableFieldToClearedField(body, recvName, "reusableReachingUses", "reachingUses") {
					clearedInDefer["reusableReachingUses"] = true
				}
				allCleared := true
				// order-independent: universal conjunction over fields with no side effects
				for f := range fields {
					if !clearedInDefer[f] {
						allCleared = false
						break
					}
				}
				if allCleared {
					return true
				}
			}
		}
		// Allow Patch 86 `if this_.liveness == nil { this_.liveness = New... } else { this_.liveness.ResetForScope(...) }`
		// immediately preceding its `defer func() { this_.liveness.ClearScopeState() }()`.
		if matchReusableInitOrResetIfStmt(st, recvName, fields, info) != "" {
			continue
		}
		if stmtAssignsNonNilToAnyField(st, recvName, fields, info) {
			return false
		}
	}
	return false
}

func methodAssignsNonNilToAnyField(body *ast.BlockStmt, recvName string, fields map[string]bool, info *types.Info) bool {
	return stmtAssignsNonNilToAnyField(body, recvName, fields, info)
}

func stmtAssignsNonNilToAnyField(root ast.Node, recvName string, fields map[string]bool, info *types.Info) bool {
	found := false
	ast.Inspect(root, func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		as, isAs := n.(*ast.AssignStmt)
		if !isAs || len(as.Lhs) != len(as.Rhs) {
			return true
		}
		for i, lhs := range as.Lhs {
			sel, isS := stripParens(lhs).(*ast.SelectorExpr)
			if !isS || sel.Sel == nil || !isIdentNamed(sel.X, recvName) || !fields[sel.Sel.Name] {
				continue
			}
			if !ntIsNilIdent(as.Rhs[i], info) {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

func (p *NodeTraversalLoweringPass) rewritePostOrderCallbackLambdas(
	ctx *PassContext,
	methodDecls map[*types.Func]*ast.FuncDecl,
	topFuncDecls map[*types.Func]*ast.FuncDecl,
) (int, error) {
	ctorObj := ctx.Pkg.Scope().Lookup("NewLambda_NodeTraversal_AbstractPostOrderCallbackInterface")
	ctorFn, ok := ctorObj.(*types.Func)
	if !ok || ctorFn == nil {
		return 0, nil
	}

	var ntFile *ast.File
	var nodeTypeSrc string
	for _, f := range ctx.Files {
		for _, decl := range f.Decls {
			fd, isFD := decl.(*ast.FuncDecl)
			if isFD && fd.Recv == nil && fd.Name != nil && fd.Name.Name == "NewLambda_NodeTraversal_AbstractPostOrderCallbackInterface" {
				ntFile = f
				if fd.Type != nil && fd.Type.Params != nil && len(fd.Type.Params.List) == 1 {
					if ft, isFT := fd.Type.Params.List[0].Type.(*ast.FuncType); isFT && ft.Params != nil && len(ft.Params.List) >= 2 {
						var buf bytes.Buffer
						if err := printer.Fprint(&buf, ctx.Fset, ft.Params.List[len(ft.Params.List)-1].Type); err == nil {
							nodeTypeSrc = buf.String()
						}
					}
				}
				break
			}
		}
	}
	if ntFile == nil || nodeTypeSrc == "" || !ctx.ShouldTransform(ntFile) {
		return 0, nil
	}

	needRel := false
	needNoOpRel := false
	sites := 0

	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		fileModified := false
		for _, decl := range file.Decls {
			enclosingFD, isFD := decl.(*ast.FuncDecl)
			if !isFD || enclosingFD.Body == nil {
				continue
			}
			var enclosingRecvVar *types.Var
			var enclosingOwner *types.Named
			if enclosingFD.Recv != nil && len(enclosingFD.Recv.List) == 1 && len(enclosingFD.Recv.List[0].Names) == 1 {
				enclosingRecvVar = varFromObj(ctx.Info.Defs[enclosingFD.Recv.List[0].Names[0]])
				if enclosingRecvVar != nil {
					if pt, isP := types.Unalias(enclosingRecvVar.Type()).(*types.Pointer); isP {
						if nm, isN := types.Unalias(pt.Elem()).(*types.Named); isN {
							enclosingOwner = nm
						}
					}
				}
			}

			ast.Inspect(enclosingFD.Body, func(n ast.Node) bool {
				call, isC := n.(*ast.CallExpr)
				if !isC || call.Ellipsis.IsValid() || len(call.Args) != 1 {
					return true
				}
				id, isID := stripParens(call.Fun).(*ast.Ident)
				if !isID || id.Name != "NewLambda_NodeTraversal_AbstractPostOrderCallbackInterface" {
					return true
				}
				if ctx.Info.Uses[id] != ctorFn {
					return true
				}
				fl, isFL := stripParens(call.Args[0]).(*ast.FuncLit)
				if !isFL || fl.Type == nil || fl.Type.Params == nil || fl.Body == nil {
					return true
				}
				ntParams := make(map[*types.Var]bool)
				badParam := false
				for _, fld := range fl.Type.Params.List {
					for _, pid := range fld.Names {
						v := varFromObj(ctx.Info.Defs[pid])
						if v == nil {
							continue
						}
						if isPtrToNodeTraversalType(v.Type(), ctx.Pkg) {
							ntParams[v] = true
						} else if isScopeTypeInPkg(v.Type(), ctx.Pkg) || isVarTypeInPkg(v.Type(), ctx.Pkg) {
							badParam = true
						}
					}
				}
				if badParam {
					return true
				}
				if len(ntParams) > 0 && !blockUsesScopeAndVarEphemerally(
					fl.Body, enclosingRecvVar, nil, enclosingOwner, ntParams, nil, nil,
					ctx.Pkg, ctx.Info, methodDecls, topFuncDecls, make(map[*ast.FuncDecl]bool), 0,
				) {
					return true
				}
				if isPureNoOpBody(fl.Body, ctx.Info) {
					id.Name = "newReleasableNoOpLambda_NodeTraversal_AbstractPostOrderCallbackInterface"
					needNoOpRel = true
				} else {
					id.Name = "newReleasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface"
					needRel = true
				}
				fileModified = true
				sites++
				return true
			})
		}
		if fileModified {
			ctx.MarkFileModified(file)
		}
	}

	if needRel && !fileHasTopLevelFunc(ntFile, "newReleasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface") {
		decls, err := parseDeclsSnippet(ctx.Fset, fmt.Sprintf(`package p

type releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface struct {
	fn func(*NodeTraversal, %[1]s, %[1]s)
}

func (l *releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface) IsJavaNil_() bool {
	return l == nil
}
func (*releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface) Implements_jscomp_NodeTraversal_AbstractPostOrderCallbackInterface() {
}
func (*releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface) CanReleaseSyntacticScopes_() {
}
func (l *releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface) Visit(p0_ *NodeTraversal, p1_ %[1]s, p2_ %[1]s) {
	l.fn(p0_, p1_, p2_)
}
func newReleasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface(fn func(*NodeTraversal, %[1]s, %[1]s)) NodeTraversal_AbstractPostOrderCallbackInterface {
	if fn == nil {
		return nil
	}
	return &releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface{fn: fn}
}
`, nodeTypeSrc))
		if err != nil {
			return sites, err
		}
		ntFile.Decls = append(ntFile.Decls, decls...)
		ctx.MarkFileModified(ntFile)
	}

	if needNoOpRel && !fileHasTopLevelFunc(ntFile, "newReleasableNoOpLambda_NodeTraversal_AbstractPostOrderCallbackInterface") {
		decls, err := parseDeclsSnippet(ctx.Fset, fmt.Sprintf(`package p

type releasableNoOpLambda_NodeTraversal_AbstractPostOrderCallbackInterface struct {
	fn func(*NodeTraversal, %[1]s, %[1]s)
}

func (l *releasableNoOpLambda_NodeTraversal_AbstractPostOrderCallbackInterface) IsJavaNil_() bool {
	return l == nil
}
func (*releasableNoOpLambda_NodeTraversal_AbstractPostOrderCallbackInterface) Implements_jscomp_NodeTraversal_AbstractPostOrderCallbackInterface() {
}
func (*releasableNoOpLambda_NodeTraversal_AbstractPostOrderCallbackInterface) CanReleaseSyntacticScopes_() {
}
func (*releasableNoOpLambda_NodeTraversal_AbstractPostOrderCallbackInterface) NoOpVisit_() {
}
func (l *releasableNoOpLambda_NodeTraversal_AbstractPostOrderCallbackInterface) Visit(p0_ *NodeTraversal, p1_ %[1]s, p2_ %[1]s) {
	l.fn(p0_, p1_, p2_)
}
func newReleasableNoOpLambda_NodeTraversal_AbstractPostOrderCallbackInterface(fn func(*NodeTraversal, %[1]s, %[1]s)) NodeTraversal_AbstractPostOrderCallbackInterface {
	if fn == nil {
		return nil
	}
	return &releasableNoOpLambda_NodeTraversal_AbstractPostOrderCallbackInterface{fn: fn}
}
`, nodeTypeSrc))
		if err != nil {
			return sites, err
		}
		ntFile.Decls = append(ntFile.Decls, decls...)
		ctx.MarkFileModified(ntFile)
	}

	return sites, nil
}

func isDeniedScopeOrVarStruct(objPkg *types.Package, name string, jscompPkg *types.Package) bool {
	if objPkg != jscompPkg {
		return false
	}
	switch name {
	case "Scope", "AbstractScopeBase", "Var", "AbstractVarBase",
		"SyntacticScopeCreator", "SyntacticScopeCreator_ScopeScanner", "MemoizedScopeCreator",
		"Reference", "ReferenceCollection", "ReferenceCollector", "LiveVariablesAnalysis":
		return true
	}
	return false
}

func isKnownScopeVarFreeInterface(objPkg *types.Package, name string, jscompPkg *types.Package) bool {
	if objPkg == nil {
		return false
	}
	if objPkg == jscompPkg {
		switch name {
		case "AbstractCompiler", "CodingConvention", "ErrorHandler", "WarningsGuard",
			"Region", "SourceExcerptProvider", "ModuleMetadataMap_ModuleMetadata":
			return true
		}
		return false
	}
	switch objPkg.Name() {
	case "rhino":
		return name == "ErrorReporter" || name == "JSType"
	case "diagnostic":
		return name == "LogFile"
	}
	return false
}

func isKnownScopeVarFreeStruct(objPkg *types.Package, name string, jscompPkg *types.Package) bool {
	if objPkg == nil {
		return false
	}
	if objPkg == jscompPkg {
		switch name {
		case "Compiler", "CompilerOptions", "CompilerInput", "JSChunk", "JSChunkGraph",
			"DiagnosticType", "DiagnosticGroup", "CheckLevel",
			"AstFactory", "AstAnalyzer", "UniqueIdSupplier", "SourceFile",
			"ModuleMetadataMap", "RuntimeJsLibManager", "AccessorSummary",
			"NodeTraversal", "NodeTraversal_Builder",
			"NodeTraversal_AbstractCfgCallbackBase",
			"FlowSensitiveInlineVariables_SideEffectPredicate":
			return true
		}
		return false
	}
	switch objPkg.Name() {
	case "rhino":
		switch name {
		case "JSTypeRegistry", "InputId", "JSDocInfo", "QualifiedName":
			return true
		}
	case "colors":
		switch name {
		case "Color", "ColorRegistry", "ColorId":
			return true
		}
	case "parser":
		switch name {
		case "FeatureSet", "FeatureSet_Feature":
			return true
		}
	}
	return false
}

func callbackAndOuterMethodsNoGlobalEscape(
	ptrT *types.Pointer,
	instNamed *types.Named,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
) bool {
	seenStruct := make(map[*types.Named]bool)
	cur := instNamed
	for cur != nil && !seenStruct[cur.Origin()] {
		seenStruct[cur.Origin()] = true
		if !methodSetNoGlobalEscape(types.NewPointer(cur), pkg, info, methodDecls) {
			return false
		}
		cur = findOuterStructNamed(cur)
	}
	_ = ptrT
	return true
}

func findOuterStructNamed(named *types.Named) *types.Named {
	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return nil
	}
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if f.Name() == "outer_" {
			if pt, ok := types.Unalias(f.Type()).(*types.Pointer); ok {
				if on, ok := types.Unalias(pt.Elem()).(*types.Named); ok {
					return on
				}
			}
		}
	}
	return nil
}

func methodSetNoGlobalEscape(
	ptrT *types.Pointer,
	pkg *types.Package,
	info *types.Info,
	methodDecls map[*types.Func]*ast.FuncDecl,
) bool {
	pkgScope := pkg.Scope()
	ms := types.NewMethodSet(ptrT)
	for i := 0; i < ms.Len(); i++ {
		fn, ok := ms.At(i).Obj().(*types.Func)
		if !ok || fn == nil {
			continue
		}
		fd := methodDecls[fn]
		if fd == nil && fn.Origin() != nil {
			fd = methodDecls[fn.Origin()]
		}
		if fd == nil || fd.Body == nil {
			continue
		}
		bad := false
		ast.Inspect(fd.Body, func(n ast.Node) bool {
			if bad || n == nil {
				return false
			}
			switch s := n.(type) {
			case *ast.GoStmt:
				bad = true
				return false
			case *ast.IncDecStmt:
				if id := rootBaseIdent(s.X); id != nil {
					if v := varFromObj(info.Uses[id]); v != nil && v.Parent() == pkgScope {
						bad = true
						return false
					}
				}
			case *ast.AssignStmt:
				for _, lhs := range s.Lhs {
					if id := rootBaseIdent(lhs); id != nil {
						if v := varFromObj(info.Uses[id]); v != nil && v.Parent() == pkgScope {
							bad = true
							return false
						}
					}
				}
			case *ast.CallExpr:
				if id, isID := stripParens(s.Fun).(*ast.Ident); isID && strings.HasPrefix(id.Name, "StaticSet_") {
					bad = true
					return false
				}
			}
			return true
		})
		if bad {
			return false
		}
	}
	return true
}

func rootBaseIdent(expr ast.Expr) *ast.Ident {
	cur := stripParens(expr)
	for cur != nil {
		switch x := cur.(type) {
		case *ast.Ident:
			return x
		case *ast.SelectorExpr:
			cur = stripParens(x.X)
		case *ast.IndexExpr:
			cur = stripParens(x.X)
		case *ast.StarExpr:
			cur = stripParens(x.X)
		default:
			return nil
		}
	}
	return nil
}

func findNodeTraversalStruct(file *ast.File, fset *token.FileSet) (*ast.StructType, ast.Expr, string, bool) {
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Name == nil || ts.Name.Name != "NodeTraversal" {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok || st.Fields == nil {
				return nil, nil, "", false
			}
			required := map[string]bool{
				"compiler":     false,
				"callback":     false,
				"scopeCreator": false,
				"obeyDestructuringAndDefaultValueExecutionOrder": false,
				"scopeNodes": false,
				"scopeObjs":  false,
				"scopeDepth": false,
			}
			var nodeTypeExpr ast.Expr
			for _, f := range st.Fields.List {
				for _, name := range f.Names {
					if _, exists := required[name.Name]; exists {
						required[name.Name] = true
					}
					if name.Name == "scopeNodes" {
						arr, isSlice := f.Type.(*ast.ArrayType)
						if !isSlice || arr.Len != nil {
							return nil, nil, "", false
						}
						nodeTypeExpr = arr.Elt
					}
					if name.Name == "scopeObjs" {
						arr, isSlice := f.Type.(*ast.ArrayType)
						if !isSlice || arr.Len != nil {
							return nil, nil, "", false
						}
					}
				}
			}
			// order-independent: conjunctive boolean check over required field names
			for _, present := range required {
				if !present {
					return nil, nil, "", false
				}
			}
			if nodeTypeExpr == nil {
				return nil, nil, "", false
			}
			var buf bytes.Buffer
			if err := printer.Fprint(&buf, fset, nodeTypeExpr); err != nil {
				return nil, nil, "", false
			}
			return st, nodeTypeExpr, buf.String(), true
		}
	}
	return nil, nil, "", false
}

func (p *NodeTraversalLoweringPass) transformNodeTraversalFile(
	ctx *PassContext,
	file *ast.File,
	ntStruct *ast.StructType,
	nodeTypeExpr ast.Expr,
	nodeTypeSrc string,
) (int, error) {
	fset := ctx.Fset
	sites := 0
	sscHasRelease := syntacticScopeCreatorHasReleaseSlabs(ctx)

	// 2a. Add inline buffer fields, defaultScopeCreator_, and alwaysTraverse_/noOpVisit_/canReleaseScopes_ flags to NodeTraversal struct.
	if !structHasNamedField(ntStruct, "scopeNodesBuf_") {
		canRelFieldLine := ""
		if sscHasRelease {
			canRelFieldLine = "\n\tcanReleaseScopes_    bool"
		}
		extraFields, err := parseStructFieldsSnippet(fset, fmt.Sprintf(`package p
type t struct {
	scopeNodesBuf_       [12]%s
	scopeObjsBuf_        [12]AbstractScope
	defaultScopeCreator_ SyntacticScopeCreator
	alwaysTraverse_      bool
	noOpVisit_           bool%s
}
`, nodeTypeSrc, canRelFieldLine))
		if err != nil {
			return sites, err
		}
		ntStruct.Fields.List = append(ntStruct.Fields.List, extraFields...)
		sites++
	} else {
		var missingFieldsSrc strings.Builder
		if !structHasNamedField(ntStruct, "alwaysTraverse_") {
			missingFieldsSrc.WriteString("\talwaysTraverse_ bool\n\tnoOpVisit_ bool\n")
		}
		if sscHasRelease && !structHasNamedField(ntStruct, "canReleaseScopes_") {
			missingFieldsSrc.WriteString("\tcanReleaseScopes_ bool\n")
		}
		if missingFieldsSrc.Len() > 0 {
			extraFields, err := parseStructFieldsSnippet(fset, "package p\ntype t struct {\n"+missingFieldsSrc.String()+"}\n")
			if err != nil {
				return sites, err
			}
			ntStruct.Fields.List = append(ntStruct.Fields.List, extraFields...)
			sites++
		}
	}

	// 2b. Rewrite (*NodeTraversal).priv_NodeTraversal_ensureScopeCapacity_* to use scopeNodesBuf_[:12] / scopeObjsBuf_[:12].
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil || !isPtrMethodOnType(fd, "NodeTraversal") {
			continue
		}
		if !hasPrivMethodPrefix(fd.Name.Name, "priv_NodeTraversal_ensureScopeCapacity_") {
			continue
		}
		if ntContainsIdent(fd.Body, "scopeNodesBuf_") {
			break
		}
		recvName := singleParamOrRecvName(fd.Recv)
		if recvName == "" {
			break
		}
		newBody, err := parseFuncBodySnippet(fset, fmt.Sprintf(`package p
func f() {
	if %[1]s.scopeDepth == int32(len(%[1]s.scopeNodes)) {
		if %[1]s.scopeDepth == 0 {
			%[1]s.scopeNodes = %[1]s.scopeNodesBuf_[:12]
			%[1]s.scopeObjs = %[1]s.scopeObjsBuf_[:12]
			return
		}
		newLen := int(%[1]s.scopeDepth) << 1
		newScopeNodes := make([]%[2]s, newLen)
		copy(newScopeNodes, %[1]s.scopeNodes)
		%[1]s.scopeNodes = newScopeNodes
		newScopeObjs := make([]AbstractScope, newLen)
		copy(newScopeObjs, %[1]s.scopeObjs)
		%[1]s.scopeObjs = newScopeObjs
	}
}
`, recvName, nodeTypeSrc))
		if err != nil {
			return sites, err
		}
		fd.Body = newBody
		sites++
		break
	}

	// 2c. Rewrite InitNodeTraversal to reuse this_.defaultScopeCreator_ and cache alwaysTraverse_ / noOpVisit_ / canReleaseScopes_.
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || fd.Body == nil || fd.Name == nil || fd.Name.Name != "InitNodeTraversal" {
			continue
		}
		paramNames := flatParamNames(fd.Type.Params)
		if len(paramNames) != 2 {
			break
		}
		thisName, builderName := paramNames[0], paramNames[1]
		rewroteInit := false

		if !ntContainsIdent(fd.Body, "defaultScopeCreator_") {
			for i, stmt := range fd.Body.List {
				ifStmt, ok := stmt.(*ast.IfStmt)
				if !ok || ifStmt.Body == nil || !ntContainsIdent(ifStmt.Body, "NewSyntacticScopeCreator_1") {
					continue
				}
				replacementIf, err := parseSingleIfSnippet(fset, fmt.Sprintf(`package p
func f() {
	if %[2]s.scopeCreator == nil || %[2]s.scopeCreator.IsJavaNil_() {
		if %[1]s.defaultScopeCreator_.reusableScanner == nil {
			InitSyntacticScopeCreator_1(&%[1]s.defaultScopeCreator_, %[1]s.compiler)
		} else {
			handler_ := StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER()
			%[1]s.defaultScopeCreator_.compiler = %[1]s.compiler
			%[1]s.defaultScopeCreator_.redeclarationHandler = handler_
			%[1]s.defaultScopeCreator_.treatProvidesAsRedeclarations = false
			%[1]s.defaultScopeCreator_.scannerInUse = false
			%[1]s.defaultScopeCreator_.reusableScanner.compiler = %[1]s.compiler
			%[1]s.defaultScopeCreator_.reusableScanner.redeclarationHandler = handler_
			%[1]s.defaultScopeCreator_.reusableScanner.treatProvidesAsRedeclarations = false
			%[1]s.defaultScopeCreator_.reusableScanner.Reset()
		}
		%[1]s.scopeCreator = &%[1]s.defaultScopeCreator_
	}
}
`, thisName, builderName))
				if err != nil {
					return sites, err
				}
				replacementIf.Else = ifStmt.Else
				fd.Body.List[i] = replacementIf
				rewroteInit = true
				break
			}
		}

		if ntContainsIdent(fd.Body, "callback") && !ntContainsIdent(fd.Body, "unwrapNodeTraversalCallback_") {
			// Strip any older non-unwrapped flag cache statements before inserting the unwrapped flag block.
			kept := make([]ast.Stmt, 0, len(fd.Body.List))
			for _, st := range fd.Body.List {
				if isLegacyInitNodeTraversalFlagStmt(st, thisName) {
					continue
				}
				kept = append(kept, st)
			}
			fd.Body.List = kept

			var flagLines strings.Builder
			fmt.Fprintf(&flagLines, "\tcbTarget_ := unwrapNodeTraversalCallback_(%[1]s.callback)\n", thisName)
			fmt.Fprintf(&flagLines, "\t_, %[1]s.alwaysTraverse_ = any(%[1]s.callback).(interface{ AlwaysShouldTraverse_() })\n", thisName)
			fmt.Fprintf(&flagLines, "\tif !%[1]s.alwaysTraverse_ {\n\t\t_, %[1]s.alwaysTraverse_ = cbTarget_.(interface{ AlwaysShouldTraverse_() })\n\t}\n", thisName)
			fmt.Fprintf(&flagLines, "\t_, %[1]s.noOpVisit_ = any(%[1]s.callback).(interface{ NoOpVisit_() })\n", thisName)
			fmt.Fprintf(&flagLines, "\tif !%[1]s.noOpVisit_ {\n\t\t_, %[1]s.noOpVisit_ = cbTarget_.(interface{ NoOpVisit_() })\n\t}\n", thisName)
			if sscHasRelease {
				fmt.Fprintf(&flagLines, "\t%[1]s.canReleaseScopes_ = canReleaseSyntacticScopesForCallback_(%[1]s.callback)\n", thisName)
			}
			flagBody, err := parseFuncBodySnippet(fset, "package p\nfunc f() {\n"+flagLines.String()+"}\n")
			if err != nil {
				return sites, err
			}
			insertAt := len(fd.Body.List)
			if insertAt > 0 {
				if ret, ok := fd.Body.List[insertAt-1].(*ast.ReturnStmt); ok && len(ret.Results) == 0 {
					insertAt--
				}
			}
			newStmts := make([]ast.Stmt, 0, len(fd.Body.List)+len(flagBody.List))
			newStmts = append(newStmts, fd.Body.List[:insertAt]...)
			newStmts = append(newStmts, flagBody.List...)
			newStmts = append(newStmts, fd.Body.List[insertAt:]...)
			fd.Body.List = newStmts
			rewroteInit = true
		} else if sscHasRelease && ntContainsIdent(fd.Body, "unwrapNodeTraversalCallback_") && !ntContainsIdent(fd.Body, "canReleaseScopes_") {
			flagBody, err := parseFuncBodySnippet(fset, fmt.Sprintf("package p\nfunc f() {\n\t%[1]s.canReleaseScopes_ = canReleaseSyntacticScopesForCallback_(%[1]s.callback)\n}\n", thisName))
			if err != nil {
				return sites, err
			}
			insertAt := len(fd.Body.List)
			if insertAt > 0 {
				if ret, ok := fd.Body.List[insertAt-1].(*ast.ReturnStmt); ok && len(ret.Results) == 0 {
					insertAt--
				}
			}
			newStmts := make([]ast.Stmt, 0, len(fd.Body.List)+len(flagBody.List))
			newStmts = append(newStmts, fd.Body.List[:insertAt]...)
			newStmts = append(newStmts, flagBody.List...)
			newStmts = append(newStmts, fd.Body.List[insertAt:]...)
			fd.Body.List = newStmts
			rewroteInit = true
		} else if sscHasRelease && ntContainsIdent(fd.Body, "canReleaseScopes_") && !ntContainsIdent(fd.Body, "canReleaseSyntacticScopesForCallback_") {
			for i, st := range fd.Body.List {
				if isLegacyInitNodeTraversalFlagStmt(st, thisName) && ntContainsIdent(st, "canReleaseScopes_") {
					flagBody, err := parseFuncBodySnippet(fset, fmt.Sprintf("package p\nfunc f() {\n\t%[1]s.canReleaseScopes_ = canReleaseSyntacticScopesForCallback_(%[1]s.callback)\n}\n", thisName))
					if err != nil {
						return sites, err
					}
					fd.Body.List[i] = flagBody.List[0]
					rewroteInit = true
					break
				}
			}
		}

		if rewroteInit {
			sites++
		}
		break
	}

	if !fileHasTopLevelFunc(file, "unwrapNodeTraversalCallback_") {
		uwDecls, err := parseDeclsSnippet(fset, `package p

func unwrapNodeTraversalCallback_(cb NodeTraversal_Callback) any {
	if uw_, ok_ := any(cb).(interface{ UnwrapCallback_() any }); ok_ {
		return uw_.UnwrapCallback_()
	}
	return cb
}
`)
		if err != nil {
			return sites, err
		}
		file.Decls = append(file.Decls, uwDecls...)
	}

	if !fileHasTopLevelFunc(file, "canReleaseSyntacticScopesForCallback_") {
		crDecls, err := parseDeclsSnippet(fset, `package p

func canReleaseSyntacticScopesForCallback_(cb NodeTraversal_Callback) bool {
	u_ := unwrapNodeTraversalCallback_(cb)
	if _, ok_ := u_.(interface{ CanReleaseSyntacticScopes_() }); ok_ {
		return true
	}
	if d_, ok_ := u_.(interface{ CanReleaseSyntacticScopesDynamic_() bool }); ok_ {
		return d_.CanReleaseSyntacticScopesDynamic_()
	}
	return false
}
`)
		if err != nil {
			return sites, err
		}
		file.Decls = append(file.Decls, crDecls...)
		if sites == 0 {
			sites++
		}
	}

	// Discover private traversal methods on (*NodeTraversal).
	var (
		traverseMethod          string
		traverseRootsMethod     string
		traverseWithScopeMethod string
		traverseAtScopeMethod   string
		traverseScopeRootMethod string
	)
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || !isPtrMethodOnType(fd, "NodeTraversal") {
			continue
		}
		name := fd.Name.Name
		switch {
		case hasPrivMethodPrefix(name, "priv_NodeTraversal_traverse_"):
			traverseMethod = name
		case hasPrivMethodPrefix(name, "priv_NodeTraversal_traverseRoots_"):
			traverseRootsMethod = name
		case hasPrivMethodPrefix(name, "priv_NodeTraversal_traverseWithScope_"):
			traverseWithScopeMethod = name
		case hasPrivMethodPrefix(name, "priv_NodeTraversal_traverseAtScope_"):
			traverseAtScopeMethod = name
		case hasPrivMethodPrefix(name, "priv_NodeTraversal_traverseScopeRoot_"):
			traverseScopeRootMethod = name
		}
	}

	if traverseMethod == "" {
		return sites, nil
	}

	releaseScopesLine := ""
	if sscHasRelease {
		releaseScopesLine = "\tif t_.canReleaseScopes_ {\n\t\tt_.defaultScopeCreator_.ReleaseSlabs_()\n\t}\n"
	}
	releaseBodySrc := fmt.Sprintf(`package p
func f() {
	savedNodes_ := t_.scopeNodes
	savedObjs_ := t_.scopeObjs
	if cap(savedNodes_) > len(t_.scopeNodesBuf_) {
		clear(savedNodes_)
	}
	if cap(savedObjs_) > len(t_.scopeObjsBuf_) {
		clear(savedObjs_)
	}
	savedScanner_ := t_.defaultScopeCreator_.reusableScanner
	if savedScanner_ != nil {
		savedScanner_.compiler = nil
		savedScanner_.redeclarationHandler = nil
		savedScanner_.treatProvidesAsRedeclarations = false
		savedScanner_.Reset()
	}
%s	*t_ = NodeTraversal{}
	t_.defaultScopeCreator_.reusableScanner = savedScanner_
	if cap(savedNodes_) > len(t_.scopeNodesBuf_) {
		t_.scopeNodes = savedNodes_[:cap(savedNodes_)]
	}
	if cap(savedObjs_) > len(t_.scopeObjsBuf_) {
		t_.scopeObjs = savedObjs_[:cap(savedObjs_)]
	}
	nodeTraversalPool_.Put(t_)
}
`, releaseScopesLine)

	// 2d. Add nodeTraversalPool_, acquireNodeTraversal_, and releaseNodeTraversal_, and pool Traverse* entry points.
	if !fileHasTopLevelFunc(file, "acquireNodeTraversal_") {
		syncQual := ensureStdImport(file, "sync")
		poolDecls, err := parseDeclsSnippet(fset, fmt.Sprintf(`package p

var nodeTraversalPool_ = %[1]s.Pool{
	New: func() any { return &NodeTraversal{} },
}

func acquireNodeTraversal_(compiler AbstractCompiler, callback NodeTraversal_Callback, scopeCreator ScopeCreator, obey bool) *NodeTraversal {
	clinit_NodeTraversal()
	t_ := nodeTraversalPool_.Get().(*NodeTraversal)
	savedNodes_ := t_.scopeNodes
	savedObjs_ := t_.scopeObjs
	b_ := NodeTraversal_Builder{
		compiler:                                       compiler,
		callback:                                       callback,
		scopeCreator:                                   scopeCreator,
		obeyDestructuringAndDefaultValueExecutionOrder: obey,
	}
	InitNodeTraversal(t_, &b_)
	if cap(savedNodes_) > len(t_.scopeNodesBuf_) {
		t_.scopeNodes = savedNodes_[:cap(savedNodes_)]
	} else {
		t_.scopeNodes = t_.scopeNodesBuf_[:len(t_.scopeNodesBuf_)]
	}
	if cap(savedObjs_) > len(t_.scopeObjsBuf_) {
		t_.scopeObjs = savedObjs_[:cap(savedObjs_)]
	} else {
		t_.scopeObjs = t_.scopeObjsBuf_[:len(t_.scopeObjsBuf_)]
	}
	return t_
}

func releaseNodeTraversal_(t_ *NodeTraversal) {}
`, syncQual))
		if err != nil {
			return sites, err
		}
		relBody, err := parseFuncBodySnippet(fset, releaseBodySrc)
		if err != nil {
			return sites, err
		}
		poolDecls[len(poolDecls)-1].(*ast.FuncDecl).Body = relBody
		file.Decls = append(file.Decls, poolDecls...)
		sites++
	} else {
		// Upgrade existing releaseNodeTraversal_ if needed (e.g., replacing savedCreator_ leak or adding canReleaseScopes_).
		for _, decl := range file.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv != nil || fd.Body == nil || fd.Name == nil || fd.Name.Name != "releaseNodeTraversal_" {
				continue
			}
			if ntContainsIdent(fd.Body, "savedCreator_") || (sscHasRelease && !ntContainsIdent(fd.Body, "canReleaseScopes_")) {
				relBody, err := parseFuncBodySnippet(fset, releaseBodySrc)
				if err != nil {
					return sites, err
				}
				fd.Body = relBody
				sites++
			}
			break
		}
	}

	scbWrapperSafe := verifyTraverseScopeRootsWrapperSafe(file)
	if traverseScopeRootMethod != "" {
		canRelBodyLine := "\t_ = canRel_\n"
		if sscHasRelease {
			canRelBodyLine = "\tif canRel_ {\n\t\tt_.canReleaseScopes_ = true\n\t}\n"
		}
		helperSrc := fmt.Sprintf(`package p

func traverseScopeRootPooled_(compiler AbstractCompiler, callback NodeTraversal_Callback, scopeCreator ScopeCreator, scopeRoot %[1]s, canRel_ bool) {
	t_ := acquireNodeTraversal_(compiler, callback, scopeCreator, false)
%[3]s	defer releaseNodeTraversal_(t_)
	t_.%[2]s(scopeRoot)
}
`, nodeTypeSrc, traverseScopeRootMethod, canRelBodyLine)
		if !fileHasTopLevelFunc(file, "traverseScopeRootPooled_") {
			helperDecls, err := parseDeclsSnippet(fset, helperSrc)
			if err != nil {
				return sites, err
			}
			file.Decls = append(file.Decls, helperDecls...)
		} else {
			for i, decl := range file.Decls {
				fd, ok := decl.(*ast.FuncDecl)
				if !ok || fd.Recv != nil || fd.Name == nil || fd.Name.Name != "traverseScopeRootPooled_" {
					continue
				}
				pNames := flatParamNames(fd.Type.Params)
				if len(pNames) != 5 || (sscHasRelease && !ntContainsIdent(fd.Body, "canReleaseScopes_")) {
					helperDecls, err := parseDeclsSnippet(fset, helperSrc)
					if err != nil {
						return sites, err
					}
					file.Decls[i] = helperDecls[0]
					sites++
				}
				break
			}
		}
	}

	mscHasRelease := memoizedScopeCreatorHasReleaseSlabs(ctx)

	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil || fd.Name == nil {
			continue
		}

		if fd.Recv == nil && fd.Name.Name == "NodeTraversal_TraverseScopeRoots" && traverseScopeRootMethod != "" {
			srSites, err := rewriteTraverseScopeRootsFunc(
				fset, fd, traverseScopeRootMethod,
				sscHasRelease && mscHasRelease,
				sscHasRelease && scbWrapperSafe,
			)
			if err != nil {
				return sites, err
			}
			sites += srSites
			continue
		}

		if ntContainsIdent(fd.Body, "acquireNodeTraversal_") {
			continue
		}

		// Static NodeTraversal_Traverse* functions.
		if fd.Recv == nil {
			params := flatParamNames(fd.Type.Params)
			var bodySnippet string
			switch {
			case fd.Name.Name == "NodeTraversal_Traverse" && len(params) == 3 && traverseMethod != "":
				// (compiler, root, cb)
				bodySnippet = fmt.Sprintf(`package p
func f() {
	t_ := acquireNodeTraversal_(%[1]s, %[3]s, nil, false)
	defer releaseNodeTraversal_(t_)
	t_.%[4]s(%[2]s)
}`, params[0], params[1], params[2], traverseMethod)

			case fd.Name.Name == "NodeTraversal_TraverseRoots" && len(params) == 4 && traverseRootsMethod != "":
				// (compiler, cb, externs, root)
				bodySnippet = fmt.Sprintf(`package p
func f() {
	t_ := acquireNodeTraversal_(%[1]s, %[2]s, nil, false)
	defer releaseNodeTraversal_(t_)
	t_.%[5]s(%[3]s, %[4]s)
}`, params[0], params[1], params[2], params[3], traverseRootsMethod)

			case fd.Name.Name == "NodeTraversal_TraverseAtScope" && len(params) == 3 && traverseAtScopeMethod != "":
				// (compiler, cb, scope)
				bodySnippet = fmt.Sprintf(`package p
func f() {
	t_ := acquireNodeTraversal_(%[1]s, %[2]s, nil, false)
	defer releaseNodeTraversal_(t_)
	t_.%[4]s(%[3]s)
}`, params[0], params[1], params[2], traverseAtScopeMethod)
			}

			if bodySnippet != "" {
				newBody, err := parseFuncBodySnippet(fset, bodySnippet)
				if err != nil {
					return sites, err
				}
				fd.Body = newBody
				sites++
			}
			continue
		}

		// (*NodeTraversal_Builder).Traverse* methods.
		if isPtrMethodOnType(fd, "NodeTraversal_Builder") {
			recvName := singleParamOrRecvName(fd.Recv)
			params := flatParamNames(fd.Type.Params)
			if recvName == "" {
				continue
			}
			var callLine string
			switch {
			case fd.Name.Name == "Traverse" && len(params) == 1 && traverseMethod != "":
				callLine = fmt.Sprintf("t_.%s(%s)", traverseMethod, params[0])
			case fd.Name.Name == "TraverseRoots" && len(params) == 2 && traverseRootsMethod != "":
				callLine = fmt.Sprintf("t_.%s(%s, %s)", traverseRootsMethod, params[0], params[1])
			case fd.Name.Name == "TraverseWithScope" && len(params) == 2 && traverseWithScopeMethod != "":
				callLine = fmt.Sprintf("t_.%s(%s, %s)", traverseWithScopeMethod, params[0], params[1])
			case fd.Name.Name == "TraverseAtScope" && len(params) == 1 && traverseAtScopeMethod != "":
				callLine = fmt.Sprintf("t_.%s(%s)", traverseAtScopeMethod, params[0])
			}
			if callLine != "" {
				newBody, err := parseFuncBodySnippet(fset, fmt.Sprintf(`package p
func f() {
	t_ := acquireNodeTraversal_(%[1]s.compiler, %[1]s.callback, %[1]s.scopeCreator, %[1]s.obeyDestructuringAndDefaultValueExecutionOrder)
	defer releaseNodeTraversal_(t_)
	%[2]s
}`, recvName, callLine))
				if err != nil {
					return sites, err
				}
				fd.Body = newBody
				sites++
			}
		}
	}

	// 2e. Rewrite priv_NodeTraversal_traverseBranch_* and priv_NodeTraversal_traverseChildren_*
	// using FirstNextTok_() / BlockOf + exact leaf unrolling when Node exposes FirstNextTok_().
	branchSites, err := p.rewriteTraverseBranchAndChildren(ctx, file, nodeTypeExpr, nodeTypeSrc)
	if err != nil {
		return sites, err
	}
	sites += branchSites

	return sites, nil
}

func memoizedScopeCreatorHasReleaseSlabs(ctx *PassContext) bool {
	if ctx == nil || ctx.Pkg == nil || ctx.Pkg.Scope() == nil {
		return false
	}
	obj := ctx.Pkg.Scope().Lookup("MemoizedScopeCreator")
	if obj == nil || obj.Type() == nil {
		return false
	}
	return hasMethodNamed(types.NewPointer(obj.Type()), ctx.Pkg, "ReleaseSlabs_")
}

func isLegacyInitNodeTraversalFlagStmt(st ast.Stmt, thisName string) bool {
	as, ok := st.(*ast.AssignStmt)
	if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 2 || len(as.Rhs) != 1 || !isIdentNamed(as.Lhs[0], "_") {
		return false
	}
	sel, ok := stripParens(as.Lhs[1]).(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || !isIdentNamed(sel.X, thisName) {
		return false
	}
	switch sel.Sel.Name {
	case "alwaysTraverse_", "noOpVisit_", "canReleaseScopes_":
		return true
	}
	return false
}

func verifyTraverseScopeRootsWrapperSafe(file *ast.File) bool {
	if file == nil {
		return false
	}
	found := false
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Name == nil || !strings.HasPrefix(ts.Name.Name, "NodeTraversal__Local_") ||
				!strings.HasSuffix(ts.Name.Name, "_TraverseScopeRootsCallback") {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok || st.Fields == nil {
				return false
			}
			found = true
			for _, f := range st.Fields.List {
				for _, nm := range f.Names {
					switch nm.Name {
					case "scopeNode", "insideScopeNode", "cap_0_traverseNested", "idHash_":
						continue
					default:
						if strings.HasPrefix(nm.Name, "cap_") && strings.HasSuffix(nm.Name, "_cb") && isIdentNamed(f.Type, "NodeTraversal_Callback") {
							continue
						}
						return false
					}
				}
			}
		}
	}
	return found
}

func rewriteTraverseScopeRootsFunc(
	fset *token.FileSet,
	fd *ast.FuncDecl,
	traverseScopeRootMethod string,
	canRelMemoized bool,
	canRelScopeRoot bool,
) (int, error) {
	if fd == nil || fd.Body == nil {
		return 0, nil
	}
	sites := 0
	params := flatParamNames(fd.Type.Params)
	wantCanRelVar := len(params) == 4 && (canRelMemoized || canRelScopeRoot)

	// 1. Dynamically compute `canRel_` from unwrapped `cb` and release local MemoizedScopeCreator when safe.
	if wantCanRelVar && !ntContainsIdent(fd.Body, "canReleaseSyntacticScopesForCallback_") {
		cbParam := params[2]
		upgradedInPlace := false
		for i, st := range fd.Body.List {
			as, isAs := st.(*ast.AssignStmt)
			if isAs && as.Tok == token.DEFINE && len(as.Lhs) == 2 && isIdentNamed(as.Lhs[0], "_") && isIdentNamed(as.Lhs[1], "canRel_") &&
				ntContainsIdent(as, "CanReleaseSyntacticScopes_") {
				upgBlock, err := parseFuncBodySnippet(fset, fmt.Sprintf("package p\nfunc f() {\n\tcanRel_ := canReleaseSyntacticScopesForCallback_(%[1]s)\n}\n", cbParam))
				if err != nil {
					return sites, err
				}
				fd.Body.List[i] = upgBlock.List[0]
				upgradedInPlace = true
				sites++
				break
			}
		}

		if !upgradedInPlace {
			scVarName, insertIdx := findLocalMemoizedScopeCreatorVarInBlock(fd.Body)

			// Remove any legacy `if _, canRel_ := any(cb).(interface{ CanReleaseSyntacticScopes_() }); canRel_ { defer scopeCreator.ReleaseSlabs_() }`
			for i, st := range fd.Body.List {
				if ifSt, isIf := st.(*ast.IfStmt); isIf && ntContainsIdent(ifSt, "CanReleaseSyntacticScopes_") {
					fd.Body.List = slices.Delete(fd.Body.List, i, i+1)
					if insertIdx > i {
						insertIdx--
					}
					break
				}
			}

			if insertIdx < 0 {
				insertIdx = 0
				for insertIdx < len(fd.Body.List) {
					if es, isES := fd.Body.List[insertIdx].(*ast.ExprStmt); isES {
						if call, isC := stripParens(es.X).(*ast.CallExpr); isC && isIdentNamed(call.Fun, "clinit_NodeTraversal") {
							insertIdx++
							continue
						}
					}
					if as, isAs := fd.Body.List[insertIdx].(*ast.AssignStmt); isAs && len(as.Lhs) == 1 && isIdentNamed(as.Lhs[0], "_") {
						insertIdx++
						continue
					}
					break
				}
			}

			var relSrc string
			if canRelMemoized && scVarName != "" {
				relSrc = fmt.Sprintf(`package p
func f() {
	canRel_ := canReleaseSyntacticScopesForCallback_(%[1]s)
	if canRel_ {
		defer %[2]s.ReleaseSlabs_()
	}
}
`, cbParam, scVarName)
			} else {
				relSrc = fmt.Sprintf(`package p
func f() {
	canRel_ := canReleaseSyntacticScopesForCallback_(%[1]s)
	_ = canRel_
}
`, cbParam)
			}
			relBlock, err := parseFuncBodySnippet(fset, relSrc)
			if err != nil {
				return sites, err
			}
			newStmts := make([]ast.Stmt, 0, len(fd.Body.List)+len(relBlock.List))
			newStmts = append(newStmts, fd.Body.List[:insertIdx]...)
			newStmts = append(newStmts, relBlock.List...)
			newStmts = append(newStmts, fd.Body.List[insertIdx:]...)
			fd.Body.List = newStmts
			sites++
		}
	}

	canRelArgName := "false"
	if canRelScopeRoot && wantCanRelVar && ntContainsIdent(fd.Body, "canReleaseSyntacticScopesForCallback_") {
		canRelArgName = "canRel_"
	}

	// 2. Rewrite chained `NodeTraversal_S_Builder()...Build().priv_NodeTraversal_traverseScopeRoot_*(scopeNode)`
	// calls (or upgrade 4-arg `traverseScopeRootPooled_` calls) to 5-arg `traverseScopeRootPooled_(..., canRelArg)`.
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		blk, ok := n.(*ast.BlockStmt)
		if !ok {
			return true
		}
		for i, st := range blk.List {
			if es, isES := st.(*ast.ExprStmt); isES {
				if call, isC := stripParens(es.X).(*ast.CallExpr); isC && isIdentNamed(call.Fun, "traverseScopeRootPooled_") && len(call.Args) == 4 {
					call.Args = append(call.Args, &ast.Ident{Name: canRelArgName})
					sites++
					continue
				}
			}
			compExpr, cbExpr, scExpr, rootExpr, matched := matchBuilderTraverseScopeRootStmt(st, traverseScopeRootMethod)
			if !matched {
				continue
			}
			blk.List[i] = &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun:  &ast.Ident{Name: "traverseScopeRootPooled_"},
					Args: []ast.Expr{compExpr, cbExpr, scExpr, rootExpr, &ast.Ident{Name: canRelArgName}},
				},
			}
			sites++
		}
		return true
	})

	return sites, nil
}

func matchBuilderTraverseScopeRootStmt(
	stmt ast.Stmt,
	traverseScopeRootMethod string,
) (compExpr, cbExpr, scExpr, rootExpr ast.Expr, ok bool) {
	es, isES := stmt.(*ast.ExprStmt)
	if !isES {
		return nil, nil, nil, nil, false
	}
	travCall, isCall := stripParens(es.X).(*ast.CallExpr)
	if !isCall || travCall.Ellipsis.IsValid() || len(travCall.Args) != 1 {
		return nil, nil, nil, nil, false
	}
	travSel, isSel := stripParens(travCall.Fun).(*ast.SelectorExpr)
	if !isSel || travSel.Sel == nil || travSel.Sel.Name != traverseScopeRootMethod {
		return nil, nil, nil, nil, false
	}
	buildCall, isBC := stripParens(travSel.X).(*ast.CallExpr)
	if !isBC || len(buildCall.Args) != 0 {
		return nil, nil, nil, nil, false
	}
	buildSel, isBS := stripParens(buildCall.Fun).(*ast.SelectorExpr)
	if !isBS || buildSel.Sel == nil || buildSel.Sel.Name != "Build" {
		return nil, nil, nil, nil, false
	}

	cur := stripParens(buildSel.X)
	for depth := 0; depth < 6; depth++ {
		call, isC := cur.(*ast.CallExpr)
		if !isC || call.Ellipsis.IsValid() {
			return nil, nil, nil, nil, false
		}
		if id, isID := stripParens(call.Fun).(*ast.Ident); isID {
			if (id.Name == "NodeTraversal_S_Builder" || id.Name == "NewNodeTraversal_Builder") && len(call.Args) == 0 {
				if compExpr == nil || cbExpr == nil {
					return nil, nil, nil, nil, false
				}
				if scExpr == nil {
					scExpr = &ast.Ident{Name: "nil"}
				}
				return compExpr, cbExpr, scExpr, travCall.Args[0], true
			}
			return nil, nil, nil, nil, false
		}
		sel, isS := stripParens(call.Fun).(*ast.SelectorExpr)
		if !isS || sel.Sel == nil || len(call.Args) != 1 {
			return nil, nil, nil, nil, false
		}
		switch sel.Sel.Name {
		case "SetCompiler":
			compExpr = call.Args[0]
		case "SetCallback_1_NodeTraversal_Callback":
			cbExpr = call.Args[0]
		case "SetScopeCreator":
			scExpr = call.Args[0]
		default:
			return nil, nil, nil, nil, false
		}
		cur = stripParens(sel.X)
	}
	return nil, nil, nil, nil, false
}

func findLocalMemoizedScopeCreatorVarInBlock(body *ast.BlockStmt) (string, int) {
	if body == nil {
		return "", -1
	}
	for i, st := range body.List {
		ds, ok := st.(*ast.DeclStmt)
		if !ok {
			continue
		}
		gd, ok := ds.Decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR || len(gd.Specs) != 1 {
			continue
		}
		vs, ok := gd.Specs[0].(*ast.ValueSpec)
		if !ok || len(vs.Names) != 1 || vs.Names[0].Name == "_" || len(vs.Values) != 1 {
			continue
		}
		star, ok := stripParens(vs.Type).(*ast.StarExpr)
		if !ok || !isIdentNamed(star.X, "MemoizedScopeCreator") {
			continue
		}
		name := vs.Names[0].Name
		next := skipSpecificBlankUse(body.List, i+1, name)
		return name, next
	}
	return "", -1
}

func nodestoreSupportsBlockOf(nodestorePkg *types.Package) bool {
	if nodestorePkg == nil || nodestorePkg.Scope() == nil {
		return false
	}
	sc := nodestorePkg.Scope()
	if sc.Lookup("Node") == nil || sc.Lookup("BlockShift") == nil {
		return false
	}
	nbObj := sc.Lookup("NodeBlock")
	if nbObj == nil || nbObj.Type() == nil {
		return false
	}
	boObj := sc.Lookup("BlockOf")
	if _, ok := boObj.(*types.Func); !ok {
		return false
	}
	fntObj, _, _ := types.LookupFieldOrMethod(types.NewPointer(nbObj.Type()), false, nodestorePkg, "FirstNextTok")
	fntFn, ok := fntObj.(*types.Func)
	if !ok {
		return false
	}
	sig, ok := fntFn.Type().(*types.Signature)
	return ok && sig.Params().Len() == 1 && sig.Results().Len() == 3
}

func (p *NodeTraversalLoweringPass) rewriteTraverseBranchAndChildren(
	ctx *PassContext,
	file *ast.File,
	nodeTypeExpr ast.Expr,
	nodeTypeSrc string,
) (int, error) {
	if ctx == nil || ctx.Pkg == nil || ctx.Info == nil || nodeTypeExpr == nil {
		return 0, nil
	}
	nodeType := ctx.Info.TypeOf(nodeTypeExpr)
	if !nodeSupportsFirstNextTokLowering(nodeType, ctx.Pkg) {
		return 0, nil
	}

	var rhinoPkg *types.Package
	if nodeNamed, ok := types.Unalias(nodeType).(*types.Named); ok && nodeNamed.Obj() != nil {
		rhinoPkg = nodeNamed.Obj().Pkg()
	}
	nodestorePkg := findNodestorePackage(ctx, rhinoPkg)
	hasBlockOf := nodestoreSupportsBlockOf(nodestorePkg)

	var (
		traverseBranchFD    *ast.FuncDecl
		traverseChildrenFD  *ast.FuncDecl
		createsBlockScopeFD *ast.FuncDecl
		hasPushScope        bool
		hasPopScope         bool
	)
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Name == nil || fd.Body == nil {
			continue
		}
		name := fd.Name.Name
		if isPtrMethodOnType(fd, "NodeTraversal") {
			switch {
			case hasPrivMethodPrefix(name, "priv_NodeTraversal_traverseBranch_"):
				traverseBranchFD = fd
			case hasPrivMethodPrefix(name, "priv_NodeTraversal_traverseChildren_"):
				traverseChildrenFD = fd
			case hasPrivMethodPrefix(name, "priv_NodeTraversal_createsBlockScope_"):
				createsBlockScopeFD = fd
			case hasPrivMethodPrefix(name, "priv_NodeTraversal_pushScope_"):
				hasPushScope = true
			case hasPrivMethodPrefix(name, "priv_NodeTraversal_popScope_"):
				hasPopScope = true
			}
			continue
		}
		if fd.Recv == nil && (hasPrivMethodPrefix(name, "priv_NodeTraversal_createsBlockScope_") ||
			hasPrivMethodPrefix(name, "NodeTraversal_S_createsBlockScope_")) {
			createsBlockScopeFD = fd
		}
	}

	if traverseBranchFD == nil || createsBlockScopeFD == nil || !hasPushScope || !hasPopScope {
		return 0, nil
	}

	branchHasFirstNext := ntContainsIdent(traverseBranchFD.Body, "FirstNextTok_")
	branchHasBlockOf := ntContainsIdent(traverseBranchFD.Body, "BlockOf")
	childrenHasFirstNext := traverseChildrenFD == nil || ntContainsIdent(traverseChildrenFD.Body, "FirstNextTok_")
	childrenHasBlockOf := traverseChildrenFD == nil || ntContainsIdent(traverseChildrenFD.Body, "BlockOf")

	branchDone := branchHasFirstNext && (!hasBlockOf || branchHasBlockOf)
	childrenDone := (childrenHasFirstNext || childrenHasBlockOf) && (!hasBlockOf || childrenHasBlockOf)
	if branchDone && childrenDone {
		return 0, nil
	}

	nsQual := ""
	if hasBlockOf {
		var ok bool
		nsQual, ok = ensurePkgImportInFile(file, nodestorePkg.Path(), "nodestore", ctx)
		if !ok {
			nsQual = ""
		}
	}

	blockScopeToks, ok := extractCreatesBlockScopeTokens(createsBlockScopeFD, ctx.Info)
	if !ok || len(blockScopeToks) == 0 {
		return 0, nil
	}

	sites := 0
	if !branchHasFirstNext {
		specialToks, rewrittenStmts, err := buildLoweredTraverseBranchStmts(
			ctx.Fset,
			ctx.Info,
			traverseBranchFD,
			createsBlockScopeFD.Name.Name,
			traverseChildrenFD,
			nsQual,
			nodeTypeSrc,
		)
		if err != nil {
			return 0, err
		}
		if rewrittenStmts == nil || len(specialToks) == 0 {
			return 0, nil
		}

		if !fileHasTopLevelVar(file, "nodeTraversalSpecialOrBlockTok_") {
			tableDecls, err := buildNodeTraversalTokTablesDecls(ctx.Fset, blockScopeToks, specialToks)
			if err != nil {
				return 0, err
			}
			file.Decls = append(file.Decls, tableDecls...)
		}

		traverseBranchFD.Body.List = rewrittenStmts
		sites++
	} else if nsQual != "" && !branchHasBlockOf && fileHasTopLevelVar(file, "nodeTraversalSpecialOrBlockTok_") {
		recvName := singleParamOrRecvName(traverseBranchFD.Recv)
		if recvName != "" {
			rewrote, err := rewriteChildLoopsInBlock(
				ctx.Fset,
				traverseBranchFD.Body,
				recvName,
				traverseBranchFD.Name.Name,
				"",
				nsQual,
				nodeTypeSrc,
			)
			if err != nil {
				return sites, err
			}
			if rewrote > 0 {
				sites++
			}
		}
	}

	if traverseChildrenFD != nil && !childrenDone && fileHasTopLevelVar(file, "nodeTraversalSpecialOrBlockTok_") {
		recvName := singleParamOrRecvName(traverseChildrenFD.Recv)
		if recvName != "" {
			rewrote, err := rewriteChildLoopsInBlock(
				ctx.Fset,
				traverseChildrenFD.Body,
				recvName,
				traverseBranchFD.Name.Name,
				"",
				nsQual,
				nodeTypeSrc,
			)
			if err != nil {
				return sites, err
			}
			if rewrote > 0 {
				sites++
			}
		}
	}

	return sites, nil
}

func nodeSupportsFirstNextTokLowering(nodeType types.Type, pkg *types.Package) bool {
	if nodeType == nil {
		return false
	}
	basic, ok := nodeType.Underlying().(*types.Basic)
	if !ok || (basic.Info()&types.IsInteger) == 0 {
		return false
	}
	fntObj, _, _ := types.LookupFieldOrMethod(nodeType, false, pkg, "FirstNextTok_")
	fntFn, ok := fntObj.(*types.Func)
	if !ok {
		return false
	}
	fntSig, ok := fntFn.Type().(*types.Signature)
	if !ok || fntSig.Params().Len() != 0 || fntSig.Results().Len() != 3 {
		return false
	}
	if !types.Identical(fntSig.Results().At(0).Type(), nodeType) ||
		!types.Identical(fntSig.Results().At(1).Type(), nodeType) {
		return false
	}
	if b, ok := fntSig.Results().At(2).Type().Underlying().(*types.Basic); !ok || b.Kind() != types.Uint8 {
		return false
	}
	for _, reqMethod := range []string{"TokenSwitchOrd_", "GetFirstChild", "GetNext"} {
		obj, _, _ := types.LookupFieldOrMethod(nodeType, false, pkg, reqMethod)
		if _, ok := obj.(*types.Func); !ok {
			return false
		}
	}
	return true
}

func extractCreatesBlockScopeTokens(fd *ast.FuncDecl, info *types.Info) ([]int, bool) {
	params := flatParamNames(fd.Type.Params)
	if len(params) != 1 || fd.Body == nil {
		return nil, false
	}
	nodeParam := params[0]
	stmts := fd.Body.List
	idx := skipPreambleStmts(stmts)
	if idx >= len(stmts) {
		return nil, false
	}

	var switchStmt *ast.SwitchStmt
	var trailingStmts []ast.Stmt

	switch s := stmts[idx].(type) {
	case *ast.SwitchStmt:
		switchStmt = s
		trailingStmts = stmts[idx+1:]
	case *ast.ReturnStmt:
		// Support Java switch-expression lowered as `return func() bool { switch n.TokenSwitchOrd_() { ... } }()`
		if idx == len(stmts)-1 && len(s.Results) == 1 {
			if call, ok := stripParens(s.Results[0]).(*ast.CallExpr); ok && len(call.Args) == 0 {
				if fl, ok := stripParens(call.Fun).(*ast.FuncLit); ok && fl.Body != nil {
					innerIdx := skipPreambleStmts(fl.Body.List)
					if innerIdx < len(fl.Body.List) {
						if sw, ok := fl.Body.List[innerIdx].(*ast.SwitchStmt); ok {
							switchStmt = sw
							trailingStmts = fl.Body.List[innerIdx+1:]
						}
					}
				}
			}
		}
	}

	if switchStmt == nil || switchStmt.Init != nil || !isTokenSwitchOrdCallOn(switchStmt.Tag, nodeParam) {
		return nil, false
	}

	hasDefaultFalse := false
	hasDefault := false
	var toks []int
	for _, rawClause := range switchStmt.Body.List {
		cc, ok := rawClause.(*ast.CaseClause)
		if !ok {
			return nil, false
		}
		if cc.List == nil {
			hasDefault = true
			if len(cc.Body) == 0 {
				continue
			}
			if len(cc.Body) == 1 && isReturnBoolConst(cc.Body[0], false, info) {
				hasDefaultFalse = true
				continue
			}
			return nil, false
		}
		for _, expr := range cc.List {
			ord, ok := evalSwitchCaseOrdinal(expr, info)
			if !ok {
				return nil, false
			}
			toks = append(toks, ord)
		}
	}

	if !hasDefaultFalse {
		// If default clause was empty or absent, the statement immediately after switch must be `return false`.
		if len(trailingStmts) != 1 || !isReturnBoolConst(trailingStmts[0], false, info) {
			return nil, false
		}
		_ = hasDefault
	} else if len(trailingStmts) > 0 {
		// Allow optional unreachable `return false` or `panic(...)` after exhaustive switch with default.
		if len(trailingStmts) != 1 || (!isReturnBoolConst(trailingStmts[0], false, info) && !isUnreachableSwitchPanicStmt(trailingStmts[0], info)) {
			return nil, false
		}
	}

	return sortDedupInts(toks), len(toks) > 0
}

func buildLoweredTraverseBranchStmts(
	fset *token.FileSet,
	info *types.Info,
	fd *ast.FuncDecl,
	createsBlockScopeName string,
	traverseChildrenFD *ast.FuncDecl,
	nsQual string,
	nodeTypeSrc string,
) ([]int, []ast.Stmt, error) {
	recvName := singleParamOrRecvName(fd.Recv)
	params := flatParamNames(fd.Type.Params)
	if recvName == "" || len(params) != 2 || fd.Body == nil {
		return nil, nil, nil
	}
	nParam, parentParam := params[0], params[1]
	stmts := fd.Body.List

	// 1. Leading switch n.TokenSwitchOrd_()
	switchIdx := skipPreambleStmts(stmts)
	if switchIdx >= len(stmts) {
		return nil, nil, nil
	}
	origSwitch, ok := stmts[switchIdx].(*ast.SwitchStmt)
	if !ok || origSwitch.Init != nil || !isTokenSwitchOrdCallOn(origSwitch.Tag, nParam) {
		return nil, nil, nil
	}
	specialToks, ok := verifyAndExtractTraverseBranchSwitch(origSwitch, recvName, info)
	if !ok || len(specialToks) == 0 {
		return nil, nil, nil
	}

	// 2. this_.currentNode = n; if !this_.callback.ShouldTraverse_3(this_, n, parent) { return }
	curNodeIdx := skipBlankAssignStmts(stmts, switchIdx+1)
	if curNodeIdx >= len(stmts) || !isCurrentNodeAssign(stmts[curNodeIdx], recvName, nParam) {
		return nil, nil, nil
	}
	shouldTravIdx := skipBlankAssignStmts(stmts, curNodeIdx+1)
	if shouldTravIdx >= len(stmts) || !isShouldTraverseEarlyReturnIf(stmts[shouldTravIdx], recvName, nParam, parentParam) {
		return nil, nil, nil
	}

	// 3.Trailing this_.currentNode = n; this_.callback.Visit_3(this_, n, parent)
	endLen := len(stmts)
	if endLen > 0 {
		if ret, ok := stmts[endLen-1].(*ast.ReturnStmt); ok && len(ret.Results) == 0 {
			endLen--
		}
	}
	if endLen < shouldTravIdx+4 {
		return nil, nil, nil
	}
	visitCallIdx := endLen - 1
	visitCurNodeIdx := endLen - 2
	if !isCurrentNodeAssign(stmts[visitCurNodeIdx], recvName, nParam) ||
		!isCallbackVisitStmt(stmts[visitCallIdx], recvName, nParam, parentParam) {
		return nil, nil, nil
	}

	// 4. Locate createsBlockScope initialization in middle section stmts[shouldTravIdx+1 : visitCurNodeIdx].
	cbsInitIdx := -1
	var cbsCallSrc string
	for i := shouldTravIdx + 1; i < visitCurNodeIdx; i++ {
		if src, ok := matchCreatesBlockScopeInit(fset, stmts[i], recvName, nParam, createsBlockScopeName); ok {
			cbsInitIdx = i
			cbsCallSrc = src
			break
		}
	}
	if cbsInitIdx < 0 {
		return nil, nil, nil
	}

	// Verify pushScope and popScope calls appear in the middle section.
	hasPushCall := false
	hasPopCall := false
	for i := cbsInitIdx + 1; i < visitCurNodeIdx; i++ {
		ast.Inspect(stmts[i], func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
			if !ok || !isIdentNamed(sel.X, recvName) || sel.Sel == nil {
				return true
			}
			if hasPrivMethodPrefix(sel.Sel.Name, "priv_NodeTraversal_pushScope_") {
				hasPushCall = true
			}
			if hasPrivMethodPrefix(sel.Sel.Name, "priv_NodeTraversal_popScope_") {
				hasPopCall = true
			}
			return true
		})
	}
	if !hasPushCall || !hasPopCall {
		return nil, nil, nil
	}

	// 5. Build cloned middle statements [cbsInitIdx+1 : visitCurNodeIdx] and rewrite child traversal loop inside them.
	midBlock := &ast.BlockStmt{
		List: append([]ast.Stmt(nil), stmts[cbsInitIdx+1:visitCurNodeIdx]...),
	}
	traverseChildrenName := ""
	if traverseChildrenFD != nil && traverseChildrenFD.Name != nil {
		traverseChildrenName = traverseChildrenFD.Name.Name
	}
	rewroteLoops, err := rewriteChildLoopsInBlock(fset, midBlock, recvName, fd.Name.Name, traverseChildrenName, nsQual, nodeTypeSrc)
	if err != nil {
		return nil, nil, err
	}
	if rewroteLoops != 1 {
		return nil, nil, nil
	}

	// Construct replacement fragments!
	entryBlock, err := parseFuncBodySnippet(fset, fmt.Sprintf(`package p
func f() {
	_, _, tok_ := %[1]s.FirstNextTok_()
	if tok_ == 0xFF {
		_ = %[1]s.TokenSwitchOrd_()
	}
}
`, nParam))
	if err != nil {
		return nil, nil, err
	}

	newSwitch := *origSwitch
	tagExpr, err := parser.ParseExpr("int32(tok_)")
	if err != nil {
		return nil, nil, err
	}
	newSwitch.Tag = tagExpr

	guardedShouldTrav := &ast.IfStmt{
		Cond: mustParseExpr(fmt.Sprintf("!%s.alwaysTraverse_", recvName)),
		Body: &ast.BlockStmt{List: []ast.Stmt{stmts[shouldTravIdx]}},
	}

	cbsGateBlock, err := parseFuncBodySnippet(fset, fmt.Sprintf(`package p
func f() {
	var createsBlockScope bool
	if nodeTraversalBlockScopeTok_[tok_] {
		createsBlockScope = %s
	}
}
`, cbsCallSrc))
	if err != nil {
		return nil, nil, err
	}

	guardedVisit := &ast.IfStmt{
		Cond: mustParseExpr(fmt.Sprintf("!%s.noOpVisit_", recvName)),
		Body: &ast.BlockStmt{List: []ast.Stmt{stmts[visitCallIdx]}},
	}

	var out []ast.Stmt
	out = append(out, stmts[:switchIdx]...)
	out = append(out, entryBlock.List...)
	out = append(out, &newSwitch)
	out = append(out, stmts[switchIdx+1:shouldTravIdx]...) // includes this_.currentNode = n (kept unconditionally!)
	out = append(out, guardedShouldTrav)
	out = append(out, stmts[shouldTravIdx+1:cbsInitIdx]...)
	out = append(out, cbsGateBlock.List...)
	out = append(out, midBlock.List...)
	out = append(out, stmts[visitCurNodeIdx]) // this_.currentNode = n (kept unconditionally!)
	out = append(out, guardedVisit)
	out = append(out, stmts[endLen:]...)

	return specialToks, out, nil
}

func verifyAndExtractTraverseBranchSwitch(sw *ast.SwitchStmt, recvName string, info *types.Info) ([]int, bool) {
	requiredPrefixes := [...]string{
		"priv_NodeTraversal_handleScript_",
		"priv_NodeTraversal_handleFunction_",
		"priv_NodeTraversal_handleModule_",
		"priv_NodeTraversal_handleClass_",
		"priv_NodeTraversal_handleClassMembers_",
		"priv_NodeTraversal_handleDestructuringOrDefaultValue_",
	}
	var seenPrefixes [len(requiredPrefixes)]bool
	var toks []int
	for _, raw := range sw.Body.List {
		cc, ok := raw.(*ast.CaseClause)
		if !ok {
			return nil, false
		}
		if cc.List == nil {
			continue
		}
		for _, expr := range cc.List {
			ord, ok := evalSwitchCaseOrdinal(expr, info)
			if !ok {
				return nil, false
			}
			toks = append(toks, ord)
		}
		matchedKnownHandler := false
		unknownPrivCall := false
		for _, st := range cc.Body {
			ast.Inspect(st, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
				if !ok || !isIdentNamed(sel.X, recvName) || sel.Sel == nil {
					return true
				}
				mName := sel.Sel.Name
				matchedThis := false
				for idx, prefix := range requiredPrefixes {
					if hasPrivMethodPrefix(mName, prefix) {
						seenPrefixes[idx] = true
						matchedKnownHandler = true
						matchedThis = true
					}
				}
				if !matchedThis && strings.HasPrefix(mName, "priv_NodeTraversal_") {
					unknownPrivCall = true
				}
				return true
			})
		}
		if !matchedKnownHandler || unknownPrivCall {
			return nil, false
		}
	}
	for _, seen := range seenPrefixes {
		if !seen {
			return nil, false
		}
	}
	return sortDedupInts(toks), len(toks) > 0
}

func rewriteChildLoopsInBlock(
	fset *token.FileSet,
	block *ast.BlockStmt,
	recvName string,
	traverseBranchMethod string,
	traverseChildrenMethod string,
	nsQual string,
	nodeTypeSrc string,
) (int, error) {
	if block == nil {
		return 0, nil
	}
	rewrote := 0
	var newList []ast.Stmt
	stmts := block.List
	for i := 0; i < len(stmts); i++ {
		// Pattern 1: direct `for child := parent.GetFirstChild(); child != 0; { ... }`
		if fs, ok := stmts[i].(*ast.ForStmt); ok {
			if parentVar, ok := matchForStmtWithInit(fs, recvName, traverseBranchMethod); ok {
				repl, err := buildUnrolledChildLoopStmts(fset, recvName, parentVar, traverseBranchMethod, nsQual, nodeTypeSrc)
				if err != nil {
					return rewrote, err
				}
				newList = append(newList, repl...)
				rewrote++
				continue
			}
			// Pattern 5: upgrade pre-lowered FirstNextTok_() child loop to BlockOf-cached child loop.
			if nsQual != "" {
				if parentVar, ok := matchExistingUnrolledChildLoop(fs, recvName, traverseBranchMethod); ok {
					repl, err := buildUnrolledChildLoopStmts(fset, recvName, parentVar, traverseBranchMethod, nsQual, nodeTypeSrc)
					if err != nil {
						return rewrote, err
					}
					newList = append(newList, repl...)
					rewrote++
					continue
				}
			}
		}

		// Pattern 2: block `{ var child ... = parent.GetFirstChild(); _ = child; for child != 0 { ... } }`
		if bs, ok := stmts[i].(*ast.BlockStmt); ok {
			if parentVar, ok := matchScopedChildLoopBlock(bs, recvName, traverseBranchMethod); ok {
				repl, err := buildUnrolledChildLoopStmts(fset, recvName, parentVar, traverseBranchMethod, nsQual, nodeTypeSrc)
				if err != nil {
					return rewrote, err
				}
				newList = append(newList, repl...)
				rewrote++
				continue
			}
		}

		// Pattern 3: `var child ... = parent.GetFirstChild()` (+ optional `_ = child`) followed by `for child != 0 { ... }`
		if childVar, parentVar, ok := matchFirstChildLocalInit(stmts[i]); ok {
			nextIdx := skipSpecificBlankUse(stmts, i+1, childVar)
			if nextIdx < len(stmts) {
				if fs, isFor := stmts[nextIdx].(*ast.ForStmt); isFor && fs.Init == nil &&
					matchChildLoopForBody(fs, childVar, parentVar, recvName, traverseBranchMethod) {
					repl, err := buildUnrolledChildLoopStmts(fset, recvName, parentVar, traverseBranchMethod, nsQual, nodeTypeSrc)
					if err != nil {
						return rewrote, err
					}
					newList = append(newList, repl...)
					rewrote++
					i = nextIdx
					continue
				}
			}
		}

		// Pattern 4: direct call `this_.priv_NodeTraversal_traverseChildren_*(parentVar)` inside traverseBranch
		if traverseChildrenMethod != "" {
			if parentVar, ok := matchDirectMethod1ArgCallStmt(stmts[i], recvName, traverseChildrenMethod); ok {
				repl, err := buildUnrolledChildLoopStmts(fset, recvName, parentVar, traverseBranchMethod, nsQual, nodeTypeSrc)
				if err != nil {
					return rewrote, err
				}
				newList = append(newList, repl...)
				rewrote++
				continue
			}
		}

		newList = append(newList, stmts[i])
	}
	block.List = newList
	return rewrote, nil
}

func matchExistingUnrolledChildLoop(fs *ast.ForStmt, recvName, traverseBranchMethod string) (string, bool) {
	if fs == nil || fs.Init == nil || fs.Post != nil || fs.Body == nil || len(fs.Body.List) != 4 {
		return "", false
	}
	childVar, parentVar, ok := matchFirstChildLocalInit(fs.Init)
	if !ok || childVar != "child_" || !isNodeNonZeroCond(fs.Cond, "child_") {
		return "", false
	}
	as0, ok := fs.Body.List[0].(*ast.AssignStmt)
	if !ok || as0.Tok != token.DEFINE || len(as0.Lhs) != 3 || len(as0.Rhs) != 1 {
		return "", false
	}
	if recv, ok := matchMethod0ArgOnIdent(as0.Rhs[0], "FirstNextTok_"); !ok || recv != "child_" {
		return "", false
	}
	if !isTraverseBranchCallStmt(fs.Body.List[2], recvName, traverseBranchMethod, "child_", parentVar) {
		return "", false
	}
	if !isSimpleIdentAssign(fs.Body.List[3], "child_", "next_") {
		return "", false
	}
	return parentVar, true
}

func buildUnrolledChildLoopStmts(
	fset *token.FileSet,
	recvName string,
	parentVar string,
	traverseBranchMethod string,
	nsQual string,
	nodeTypeSrc string,
) ([]ast.Stmt, error) {
	if nsQual != "" && nodeTypeSrc != "" {
		body, err := parseFuncBodySnippet(fset, fmt.Sprintf(`package p
func f() {
	var curBlkIdx_ uint32 = 0xFFFFFFFF
	var curBlk_ *%[4]s.NodeBlock
	for child_ := %[2]s.GetFirstChild(); child_ != 0; {
		var cFirst_, next_ %[5]s
		var cTok_ uint8
		if idx_ := uint32(child_) >> %[4]s.BlockShift; idx_ == curBlkIdx_ {
			f_, nx_, t_ := curBlk_.FirstNextTok(%[4]s.Node(child_))
			cFirst_, next_, cTok_ = %[5]s(f_), %[5]s(nx_), t_
		} else {
			curBlk_ = %[4]s.BlockOf(%[4]s.Node(child_))
			curBlkIdx_ = idx_
			f_, nx_, t_ := curBlk_.FirstNextTok(%[4]s.Node(child_))
			cFirst_, next_, cTok_ = %[5]s(f_), %[5]s(nx_), t_
		}
		if cFirst_ == 0 && cTok_ != 0xFF && !nodeTraversalSpecialOrBlockTok_[cTok_] {
			if %[1]s.alwaysTraverse_ {
				%[1]s.currentNode = child_
				if !%[1]s.noOpVisit_ {
					%[1]s.callback.Visit_3(%[1]s, child_, %[2]s)
				}
				child_ = next_
				continue
			}
			%[1]s.currentNode = child_
			if !%[1]s.callback.ShouldTraverse_3(%[1]s, child_, %[2]s) {
				child_ = next_
				continue
			}
			if child_.GetFirstChild() == 0 {
				%[1]s.currentNode = child_
				if !%[1]s.noOpVisit_ {
					%[1]s.callback.Visit_3(%[1]s, child_, %[2]s)
				}
				child_ = next_
				continue
			}
			for gc_ := child_.GetFirstChild(); gc_ != 0; {
				var gNext_ %[5]s
				if gIdx_ := uint32(gc_) >> %[4]s.BlockShift; gIdx_ == curBlkIdx_ {
					_, gnx_, _ := curBlk_.FirstNextTok(%[4]s.Node(gc_))
					gNext_ = %[5]s(gnx_)
				} else {
					curBlk_ = %[4]s.BlockOf(%[4]s.Node(gc_))
					curBlkIdx_ = gIdx_
					_, gnx_, _ := curBlk_.FirstNextTok(%[4]s.Node(gc_))
					gNext_ = %[5]s(gnx_)
				}
				%[1]s.%[3]s(gc_, child_)
				gc_ = gNext_
			}
			%[1]s.currentNode = child_
			if !%[1]s.noOpVisit_ {
				%[1]s.callback.Visit_3(%[1]s, child_, %[2]s)
			}
			child_ = next_
			continue
		}
		%[1]s.%[3]s(child_, %[2]s)
		child_ = next_
	}
}
`, recvName, parentVar, traverseBranchMethod, nsQual, nodeTypeSrc))
		if err != nil {
			return nil, err
		}
		return body.List, nil
	}

	body, err := parseFuncBodySnippet(fset, fmt.Sprintf(`package p
func f() {
	for child_ := %[2]s.GetFirstChild(); child_ != 0; {
		cFirst_, next_, cTok_ := child_.FirstNextTok_()
		if cFirst_ == 0 && cTok_ != 0xFF && !nodeTraversalSpecialOrBlockTok_[cTok_] {
			if %[1]s.alwaysTraverse_ {
				%[1]s.currentNode = child_
				if !%[1]s.noOpVisit_ {
					%[1]s.callback.Visit_3(%[1]s, child_, %[2]s)
				}
				child_ = next_
				continue
			}
			%[1]s.currentNode = child_
			if !%[1]s.callback.ShouldTraverse_3(%[1]s, child_, %[2]s) {
				child_ = next_
				continue
			}
			if child_.GetFirstChild() == 0 {
				%[1]s.currentNode = child_
				if !%[1]s.noOpVisit_ {
					%[1]s.callback.Visit_3(%[1]s, child_, %[2]s)
				}
				child_ = next_
				continue
			}
			for gc_ := child_.GetFirstChild(); gc_ != 0; {
				gNext_ := gc_.GetNext()
				%[1]s.%[3]s(gc_, child_)
				gc_ = gNext_
			}
			%[1]s.currentNode = child_
			if !%[1]s.noOpVisit_ {
				%[1]s.callback.Visit_3(%[1]s, child_, %[2]s)
			}
			child_ = next_
			continue
		}
		%[1]s.%[3]s(child_, %[2]s)
		child_ = next_
	}
}
`, recvName, parentVar, traverseBranchMethod))
	if err != nil {
		return nil, err
	}
	return body.List, nil
}

func buildNodeTraversalTokTablesDecls(fset *token.FileSet, blockScopeToks, specialToks []int) ([]ast.Decl, error) {
	var blockParts []string
	for _, t := range blockScopeToks {
		blockParts = append(blockParts, fmt.Sprintf("%d: true", t))
	}
	combined := append([]int{0xFF}, blockScopeToks...)
	combined = append(combined, specialToks...)
	combined = sortDedupInts(combined)
	var combParts []string
	for _, t := range combined {
		combParts = append(combParts, fmt.Sprintf("%d: true", t))
	}
	src := fmt.Sprintf(`package p

var nodeTraversalBlockScopeTok_ = [256]bool{%s}

var nodeTraversalSpecialOrBlockTok_ = [256]bool{%s}
`, strings.Join(blockParts, ", "), strings.Join(combParts, ", "))
	return parseDeclsSnippet(fset, src)
}

func matchForStmtWithInit(fs *ast.ForStmt, recvName, traverseBranchMethod string) (string, bool) {
	if fs == nil || fs.Init == nil || fs.Post != nil {
		return "", false
	}
	childVar, parentVar, ok := matchFirstChildLocalInit(fs.Init)
	if !ok {
		return "", false
	}
	if !matchChildLoopForBody(fs, childVar, parentVar, recvName, traverseBranchMethod) {
		return "", false
	}
	return parentVar, true
}

func matchScopedChildLoopBlock(bs *ast.BlockStmt, recvName, traverseBranchMethod string) (string, bool) {
	if bs == nil || len(bs.List) < 2 {
		return "", false
	}
	childVar, parentVar, ok := matchFirstChildLocalInit(bs.List[0])
	if !ok {
		return "", false
	}
	idx := skipSpecificBlankUse(bs.List, 1, childVar)
	if idx != len(bs.List)-1 {
		return "", false
	}
	fs, ok := bs.List[idx].(*ast.ForStmt)
	if !ok || fs.Init != nil || fs.Post != nil {
		return "", false
	}
	if !matchChildLoopForBody(fs, childVar, parentVar, recvName, traverseBranchMethod) {
		return "", false
	}
	return parentVar, true
}

func matchChildLoopForBody(fs *ast.ForStmt, childVar, parentVar, recvName, traverseBranchMethod string) bool {
	if fs == nil || fs.Post != nil || fs.Body == nil || !isNodeNonZeroCond(fs.Cond, childVar) {
		return false
	}
	var core []ast.Stmt
	for _, s := range fs.Body.List {
		if isBlankIdentAssign(s) {
			continue
		}
		core = append(core, s)
	}
	if len(core) != 3 {
		return false
	}
	nextVar, ok := matchNextSiblingInit(core[0], childVar)
	if !ok || nextVar == "" || nextVar == childVar || nextVar == parentVar {
		return false
	}
	if !isTraverseBranchCallStmt(core[1], recvName, traverseBranchMethod, childVar, parentVar) {
		return false
	}
	return isSimpleIdentAssign(core[2], childVar, nextVar)
}

func matchFirstChildLocalInit(stmt ast.Stmt) (childVar, parentVar string, ok bool) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if s.Tok != token.DEFINE || len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			return "", "", false
		}
		id, isID := s.Lhs[0].(*ast.Ident)
		if !isID || id.Name == "" || id.Name == "_" {
			return "", "", false
		}
		pVar, isCall := matchMethod0ArgOnIdent(s.Rhs[0], "GetFirstChild")
		if !isCall {
			return "", "", false
		}
		return id.Name, pVar, true
	case *ast.DeclStmt:
		gd, isGD := s.Decl.(*ast.GenDecl)
		if !isGD || gd.Tok != token.VAR || len(gd.Specs) != 1 {
			return "", "", false
		}
		vs, isVS := gd.Specs[0].(*ast.ValueSpec)
		if !isVS || len(vs.Names) != 1 || len(vs.Values) != 1 || vs.Names[0].Name == "" || vs.Names[0].Name == "_" {
			return "", "", false
		}
		pVar, isCall := matchMethod0ArgOnIdent(vs.Values[0], "GetFirstChild")
		if !isCall {
			return "", "", false
		}
		return vs.Names[0].Name, pVar, true
	}
	return "", "", false
}

func matchNextSiblingInit(stmt ast.Stmt, childVar string) (string, bool) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if s.Tok != token.DEFINE || len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			return "", false
		}
		id, isID := s.Lhs[0].(*ast.Ident)
		if !isID || id.Name == "" || id.Name == "_" {
			return "", false
		}
		cVar, isCall := matchMethod0ArgOnIdent(s.Rhs[0], "GetNext")
		return id.Name, isCall && cVar == childVar
	case *ast.DeclStmt:
		gd, isGD := s.Decl.(*ast.GenDecl)
		if !isGD || gd.Tok != token.VAR || len(gd.Specs) != 1 {
			return "", false
		}
		vs, isVS := gd.Specs[0].(*ast.ValueSpec)
		if !isVS || len(vs.Names) != 1 || len(vs.Values) != 1 || vs.Names[0].Name == "" || vs.Names[0].Name == "_" {
			return "", false
		}
		cVar, isCall := matchMethod0ArgOnIdent(vs.Values[0], "GetNext")
		return vs.Names[0].Name, isCall && cVar == childVar
	}
	return "", false
}

func matchMethod0ArgOnIdent(expr ast.Expr, methodName string) (string, bool) {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 0 {
		return "", false
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || sel.Sel.Name != methodName {
		return "", false
	}
	id, ok := stripParens(sel.X).(*ast.Ident)
	if !ok || id.Name == "" {
		return "", false
	}
	return id.Name, true
}

func isNodeNonZeroCond(cond ast.Expr, varName string) bool {
	e := stripParens(cond)
	switch c := e.(type) {
	case *ast.BinaryExpr:
		if c.Op == token.NEQ {
			return (isIdentNamed(c.X, varName) && isZeroNodeExpr(c.Y)) ||
				(isZeroNodeExpr(c.X) && isIdentNamed(c.Y, varName))
		}
	case *ast.UnaryExpr:
		if c.Op != token.NOT {
			return false
		}
		inner := stripParens(c.X)
		if bin, ok := inner.(*ast.BinaryExpr); ok && bin.Op == token.EQL {
			return (isIdentNamed(bin.X, varName) && isZeroNodeExpr(bin.Y)) ||
				(isZeroNodeExpr(bin.X) && isIdentNamed(bin.Y, varName))
		}
		if recv, ok := matchMethod0ArgOnIdent(inner, "IsJavaNil_"); ok {
			return recv == varName
		}
	}
	return false
}

func isZeroNodeExpr(expr ast.Expr) bool {
	e := stripParens(expr)
	switch x := e.(type) {
	case *ast.BasicLit:
		return x.Kind == token.INT && x.Value == "0"
	case *ast.CallExpr:
		if len(x.Args) == 1 {
			return isZeroNodeExpr(x.Args[0])
		}
	}
	return false
}

func isTraverseBranchCallStmt(stmt ast.Stmt, recvName, methodName, childVar, parentVar string) bool {
	es, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return false
	}
	call, ok := stripParens(es.X).(*ast.CallExpr)
	if !ok || len(call.Args) != 2 {
		return false
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || !isIdentNamed(sel.X, recvName) || sel.Sel == nil || sel.Sel.Name != methodName {
		return false
	}
	return isIdentNamed(call.Args[0], childVar) && isIdentNamed(call.Args[1], parentVar)
}

func matchDirectMethod1ArgCallStmt(stmt ast.Stmt, recvName, methodName string) (string, bool) {
	es, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return "", false
	}
	call, ok := stripParens(es.X).(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return "", false
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || !isIdentNamed(sel.X, recvName) || sel.Sel == nil || sel.Sel.Name != methodName {
		return "", false
	}
	id, ok := stripParens(call.Args[0]).(*ast.Ident)
	if !ok || id.Name == "" {
		return "", false
	}
	return id.Name, true
}

func isSimpleIdentAssign(stmt ast.Stmt, lhsName, rhsName string) bool {
	as, ok := stmt.(*ast.AssignStmt)
	if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
		return false
	}
	return isIdentNamed(as.Lhs[0], lhsName) && isIdentNamed(as.Rhs[0], rhsName)
}

func matchCreatesBlockScopeInit(
	fset *token.FileSet,
	stmt ast.Stmt,
	recvName, nParam, methodName string,
) (string, bool) {
	var callExpr ast.Expr
	switch s := stmt.(type) {
	case *ast.DeclStmt:
		gd, ok := s.Decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR || len(gd.Specs) != 1 {
			return "", false
		}
		vs, ok := gd.Specs[0].(*ast.ValueSpec)
		if !ok || len(vs.Names) != 1 || vs.Names[0].Name != "createsBlockScope" || len(vs.Values) != 1 {
			return "", false
		}
		callExpr = vs.Values[0]
	case *ast.AssignStmt:
		if s.Tok != token.DEFINE || len(s.Lhs) != 1 || len(s.Rhs) != 1 || !isIdentNamed(s.Lhs[0], "createsBlockScope") {
			return "", false
		}
		callExpr = s.Rhs[0]
	default:
		return "", false
	}

	call, ok := stripParens(callExpr).(*ast.CallExpr)
	if !ok || len(call.Args) != 1 || !isIdentNamed(call.Args[0], nParam) {
		return "", false
	}
	switch fn := stripParens(call.Fun).(type) {
	case *ast.Ident:
		if fn.Name != methodName {
			return "", false
		}
	case *ast.SelectorExpr:
		if !isIdentNamed(fn.X, recvName) || fn.Sel == nil || fn.Sel.Name != methodName {
			return "", false
		}
	default:
		return "", false
	}

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, call); err != nil {
		return "", false
	}
	return buf.String(), true
}

func isCurrentNodeAssign(stmt ast.Stmt, recvName, nParam string) bool {
	as, ok := stmt.(*ast.AssignStmt)
	if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
		return false
	}
	sel, ok := stripParens(as.Lhs[0]).(*ast.SelectorExpr)
	if !ok || !isIdentNamed(sel.X, recvName) || sel.Sel == nil || sel.Sel.Name != "currentNode" {
		return false
	}
	return isIdentNamed(as.Rhs[0], nParam)
}

func isShouldTraverseEarlyReturnIf(stmt ast.Stmt, recvName, nParam, parentParam string) bool {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok || ifStmt.Init != nil || ifStmt.Body == nil || len(ifStmt.Body.List) != 1 {
		return false
	}
	if ifStmt.Else != nil {
		if eb, ok := ifStmt.Else.(*ast.BlockStmt); !ok || len(eb.List) != 0 {
			return false
		}
	}
	ret, ok := ifStmt.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 0 {
		return false
	}
	u, ok := stripParens(ifStmt.Cond).(*ast.UnaryExpr)
	if !ok || u.Op != token.NOT {
		return false
	}
	return isCallbackMethod3Call(u.X, recvName, "ShouldTraverse_3", nParam, parentParam)
}

func isCallbackVisitStmt(stmt ast.Stmt, recvName, nParam, parentParam string) bool {
	es, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return false
	}
	return isCallbackMethod3Call(es.X, recvName, "Visit_3", nParam, parentParam)
}

func isCallbackMethod3Call(expr ast.Expr, recvName, methodName, nParam, parentParam string) bool {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 3 {
		return false
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || sel.Sel.Name != methodName {
		return false
	}
	cbSel, ok := stripParens(sel.X).(*ast.SelectorExpr)
	if !ok || !isIdentNamed(cbSel.X, recvName) || cbSel.Sel == nil || cbSel.Sel.Name != "callback" {
		return false
	}
	return ntContainsIdent(call.Args[0], recvName) &&
		isIdentNamed(call.Args[1], nParam) &&
		isIdentNamed(call.Args[2], parentParam)
}

func isTokenSwitchOrdCallOn(expr ast.Expr, varName string) bool {
	recv, ok := matchMethod0ArgOnIdent(expr, "TokenSwitchOrd_")
	return ok && recv == varName
}

func isReturnBoolConst(stmt ast.Stmt, want bool, info *types.Info) bool {
	ret, ok := stmt.(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return false
	}
	id, ok := stripParens(ret.Results[0]).(*ast.Ident)
	if !ok {
		return false
	}
	wantName := "false"
	if want {
		wantName = "true"
	}
	if id.Name != wantName {
		return false
	}
	if info != nil && info.Uses != nil {
		return info.Uses[id] == types.Universe.Lookup(wantName)
	}
	return true
}

func evalSwitchCaseOrdinal(expr ast.Expr, info *types.Info) (int, bool) {
	if info != nil && info.Types != nil {
		if tv, ok := info.Types[expr]; ok && tv.Value != nil {
			v := constant.ToInt(tv.Value)
			if v.Kind() == constant.Int {
				if ord, exact := constant.Int64Val(v); exact && ord >= 0 && ord < 255 {
					return int(ord), true
				}
			}
		}
	}
	ord, ok := evalConstIntFallback(expr)
	if !ok || ord < 0 || ord >= 255 {
		return 0, false
	}
	return int(ord), true
}

func evalConstIntFallback(expr ast.Expr) (int64, bool) {
	switch e := stripParens(expr).(type) {
	case *ast.BasicLit:
		if e.Kind != token.INT {
			return 0, false
		}
		v, err := strconv.ParseInt(e.Value, 0, 64)
		return v, err == nil
	case *ast.CallExpr:
		if len(e.Args) == 1 {
			if id, ok := stripParens(e.Fun).(*ast.Ident); ok {
				switch id.Name {
				case "int", "int32", "int64", "uint8", "uint32":
					return evalConstIntFallback(e.Args[0])
				}
			}
		}
	}
	return 0, false
}

func skipPreambleStmts(stmts []ast.Stmt) int {
	i := 0
	for i < len(stmts) {
		if isBlankIdentAssign(stmts[i]) || isClinitNodeTraversalStmt(stmts[i]) {
			i++
			continue
		}
		break
	}
	return i
}

func skipBlankAssignStmts(stmts []ast.Stmt, start int) int {
	i := start
	for i < len(stmts) && isBlankIdentAssign(stmts[i]) {
		i++
	}
	return i
}

func skipSpecificBlankUse(stmts []ast.Stmt, start int, varName string) int {
	i := start
	for i < len(stmts) {
		as, ok := stmts[i].(*ast.AssignStmt)
		if ok && as.Tok == token.ASSIGN && len(as.Lhs) == 1 && len(as.Rhs) == 1 &&
			isIdentNamed(as.Lhs[0], "_") && isIdentNamed(as.Rhs[0], varName) {
			i++
			continue
		}
		break
	}
	return i
}

func isBlankIdentAssign(stmt ast.Stmt) bool {
	as, ok := stmt.(*ast.AssignStmt)
	if !ok || as.Tok != token.ASSIGN || len(as.Lhs) == 0 || len(as.Lhs) != len(as.Rhs) {
		return false
	}
	for i := range as.Lhs {
		if !isIdentNamed(as.Lhs[i], "_") {
			return false
		}
		if _, ok := stripParens(as.Rhs[i]).(*ast.Ident); !ok {
			return false
		}
	}
	return true
}

func isClinitNodeTraversalStmt(stmt ast.Stmt) bool {
	es, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return false
	}
	call, ok := stripParens(es.X).(*ast.CallExpr)
	if !ok || len(call.Args) != 0 {
		return false
	}
	return isIdentNamed(call.Fun, "clinit_NodeTraversal")
}

func isUnreachableSwitchPanicStmt(stmt ast.Stmt, info *types.Info) bool {
	es, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return false
	}
	call, ok := stripParens(es.X).(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return false
	}
	id, ok := stripParens(call.Fun).(*ast.Ident)
	if !ok || id.Name != "panic" {
		return false
	}
	if info != nil && info.Uses != nil {
		return info.Uses[id] == types.Universe.Lookup("panic")
	}
	return true
}

func isIdentNamed(expr ast.Expr, want string) bool {
	id, ok := stripParens(expr).(*ast.Ident)
	return ok && id.Name == want
}

func sortDedupInts(in []int) []int {
	if len(in) == 0 {
		return nil
	}
	cp := append([]int(nil), in...)
	slices.Sort(cp)
	return slices.Compact(cp)
}

func mustParseExpr(src string) ast.Expr {
	e, err := parser.ParseExpr(src)
	if err != nil {
		panic(err)
	}
	return e
}

func structHasNamedField(st *ast.StructType, target string) bool {
	if st == nil || st.Fields == nil {
		return false
	}
	for _, f := range st.Fields.List {
		for _, name := range f.Names {
			if name.Name == target {
				return true
			}
		}
	}
	return false
}

func isPtrMethodOnType(fd *ast.FuncDecl, typeName string) bool {
	if fd == nil || fd.Recv == nil || len(fd.Recv.List) != 1 {
		return false
	}
	star, ok := stripParens(fd.Recv.List[0].Type).(*ast.StarExpr)
	if !ok {
		return false
	}
	id, ok := stripParens(star.X).(*ast.Ident)
	return ok && id.Name == typeName
}

func hasPrivMethodPrefix(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	rem := name[len(prefix):]
	if rem == "" {
		return true
	}
	for i := 0; i < len(rem); i++ {
		if rem[i] < '0' || rem[i] > '9' {
			return false
		}
	}
	return true
}

func singleParamOrRecvName(fl *ast.FieldList) string {
	names := flatParamNames(fl)
	if len(names) != 1 {
		return ""
	}
	return names[0]
}

func flatParamNames(fl *ast.FieldList) []string {
	if fl == nil {
		return nil
	}
	var out []string
	for _, f := range fl.List {
		for _, n := range f.Names {
			if n.Name == "" || n.Name == "_" {
				return nil
			}
			out = append(out, n.Name)
		}
	}
	return out
}

func ntContainsIdent(root ast.Node, name string) bool {
	found := false
	ast.Inspect(root, func(n ast.Node) bool {
		if found {
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

func fileHasTopLevelFunc(file *ast.File, name string) bool {
	for _, decl := range file.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Recv == nil && fd.Name != nil && fd.Name.Name == name {
			return true
		}
	}
	return false
}

func fileHasTopLevelVar(file *ast.File, name string) bool {
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, id := range vs.Names {
				if id != nil && id.Name == name {
					return true
				}
			}
		}
	}
	return false
}

func ensureStdImport(file *ast.File, pkgPath string) string {
	for _, spec := range file.Imports {
		if spec.Path == nil {
			continue
		}
		p, err := strconv.Unquote(spec.Path.Value)
		if err != nil || p != pkgPath {
			continue
		}
		if spec.Name != nil && spec.Name.Name != "_" && spec.Name.Name != "." {
			return spec.Name.Name
		}
		if spec.Name == nil {
			return pkgPath
		}
	}

	spec := &ast.ImportSpec{
		Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(pkgPath)},
	}
	file.Imports = append([]*ast.ImportSpec{spec}, file.Imports...)
	for _, decl := range file.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
			if len(gd.Specs) > 0 {
				spec.Path.ValuePos = gd.Specs[0].Pos()
			}
			gd.Specs = append([]ast.Spec{spec}, gd.Specs...)
			return pkgPath
		}
	}
	gd := &ast.GenDecl{Tok: token.IMPORT, Specs: []ast.Spec{spec}}
	file.Decls = append([]ast.Decl{gd}, file.Decls...)
	return pkgPath
}

func parseStructFieldsSnippet(fset *token.FileSet, src string) ([]*ast.Field, error) {
	af, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, err
	}
	gd, ok := af.Decls[0].(*ast.GenDecl)
	if !ok || len(gd.Specs) == 0 {
		return nil, fmt.Errorf("unexpected struct snippet AST")
	}
	ts, ok := gd.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil, fmt.Errorf("unexpected TypeSpec in snippet")
	}
	st, ok := ts.Type.(*ast.StructType)
	if !ok || st.Fields == nil {
		return nil, fmt.Errorf("unexpected StructType in snippet")
	}
	return st.Fields.List, nil
}

func parseFuncBodySnippet(fset *token.FileSet, src string) (*ast.BlockStmt, error) {
	af, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, err
	}
	fd, ok := af.Decls[0].(*ast.FuncDecl)
	if !ok || fd.Body == nil {
		return nil, fmt.Errorf("unexpected func snippet AST")
	}
	return fd.Body, nil
}

func parseSingleIfSnippet(fset *token.FileSet, src string) (*ast.IfStmt, error) {
	body, err := parseFuncBodySnippet(fset, src)
	if err != nil {
		return nil, err
	}
	if len(body.List) != 1 {
		return nil, fmt.Errorf("expected single IfStmt in snippet")
	}
	ifStmt, ok := body.List[0].(*ast.IfStmt)
	if !ok {
		return nil, fmt.Errorf("expected IfStmt in snippet")
	}
	return ifStmt, nil
}

func parseDeclsSnippet(fset *token.FileSet, src string) ([]ast.Decl, error) {
	af, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, err
	}
	return af.Decls, nil
}

// --- Wave 24C: Zero-Alloc Scope / TypedScope GetVarIterable() Cursor Loop Lowering ---

func (p *NodeTraversalLoweringPass) applyScopeVarCursorLowering(ctx *PassContext) (int, error) {
	if ctx == nil || ctx.Pkg == nil || ctx.Info == nil || ctx.Info.Defs == nil || ctx.Info.Uses == nil || ctx.Info.Types == nil {
		return 0, nil
	}
	if !j2gortHasExport(ctx, "OrderedMapValCursor") || !j2gortHasExport(ctx, "ThrowNPE") {
		return 0, nil
	}

	scopeNamed, varNamed, scopeFile, scopeOK := verifyInlineVarScopeLayout("Scope", "Var", ctx)
	if !scopeOK || scopeFile == nil {
		return 0, nil
	}
	canEmitCursorType := ctx.ShouldTransform(scopeFile) || fileHasTypeSpec(scopeFile, "scopeVarCursor_")
	if !canEmitCursorType {
		return 0, nil
	}
	canEmitScopeIter := ctx.ShouldTransform(scopeFile) || fileHasPtrMethodNamed(scopeFile, "Scope", "iterVars_")

	typedScopeNamed, typedVarNamed, typedScopeFile, typedScopeOK := verifyInlineVarScopeLayout("TypedScope", "TypedVar", ctx)
	canEmitTypedScopeIter := typedScopeOK && typedScopeFile != nil &&
		(ctx.ShouldTransform(typedScopeFile) || fileHasPtrMethodNamed(typedScopeFile, "TypedScope", "iterVars_"))

	sites := 0
	needScopeIter := false
	needTypedScopeIter := false

	for _, file := range ctx.Files {
		if !ctx.ShouldTransform(file) {
			continue
		}
		fileSites := 0
		for _, decl := range file.Decls {
			fd, isFD := decl.(*ast.FuncDecl)
			if !isFD || fd.Body == nil {
				continue
			}
			ast.Inspect(fd.Body, func(n ast.Node) bool {
				blk, isBlk := n.(*ast.BlockStmt)
				if !isBlk || len(blk.List) < 2 {
					return true
				}
				for i := 0; i < len(blk.List)-1; i++ {
					curDefId, curVar, scopeExpr, isTypedScope, okInit := matchScopeGetVarIterableCursorDefine(
						blk.List[i], scopeNamed, typedScopeNamed, canEmitScopeIter, canEmitTypedScopeIter, ctx.Info,
					)
					if !okInit {
						continue
					}
					forIdx := i + 1
					var curDiscardId *ast.Ident
					if forIdx < len(blk.List) {
						if dId := matchBlankDiscardIdent(blk.List[forIdx], curVar, ctx.Info); dId != nil {
							curDiscardId = dId
							forIdx++
						}
					}
					if forIdx >= len(blk.List) {
						continue
					}
					forStmt, isFor := blk.List[forIdx].(*ast.ForStmt)
					if !isFor || forStmt.Init != nil || forStmt.Post != nil || forStmt.Body == nil {
						continue
					}
					condCurId := matchZeroArgMethodCallIdent(forStmt.Cond, "HasNext", curVar, ctx.Info)
					if condCurId == nil {
						continue
					}
					wantVarNamed := varNamed
					if isTypedScope {
						wantVarNamed = typedVarNamed
					}
					if !rewriteScopeVarCursorLoopBody(fd.Body, forStmt, curDefId, curDiscardId, condCurId, curVar, wantVarNamed, ctx.Info) {
						continue
					}
					as := blk.List[i].(*ast.AssignStmt)
					as.Rhs[0] = &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.ParenExpr{X: scopeExpr},
							Sel: &ast.Ident{Name: "iterVars_"},
						},
					}
					fileSites++
					if isTypedScope {
						needTypedScopeIter = true
					} else {
						needScopeIter = true
					}
				}
				return true
			})
		}
		if fileSites > 0 {
			ctx.MarkFileModified(file)
			sites += fileSites
		}
	}

	if (needScopeIter || needTypedScopeIter) && !fileHasTypeSpec(scopeFile, "scopeVarCursor_") {
		j2Qual, ok := ensureJ2gortImportInFile(scopeFile, ctx)
		if !ok {
			return sites, fmt.Errorf("failed to ensure j2gort import in Scope file")
		}
		decls, err := parseDeclsSnippet(ctx.Fset, fmt.Sprintf(`package p

type scopeVarCursor_[V any] struct {
	vars   [6]V
	idx    int
	n      int
	useMap bool
	mapCur %[1]s.OrderedMapValCursor[string, V]
}

func (c *scopeVarCursor_[V]) HasNext() bool {
	if c.useMap {
		return c.mapCur.HasNext()
	}
	return c.idx < c.n
}

func (c *scopeVarCursor_[V]) Next() V {
	if c.useMap {
		return c.mapCur.Next()
	}
	v := c.vars[c.idx]
	c.idx++
	return v
}
`, j2Qual))
		if err != nil {
			return sites, err
		}
		scopeFile.Decls = append(scopeFile.Decls, decls...)
		ctx.MarkFileModified(scopeFile)
	}

	if needScopeIter && !fileHasPtrMethodNamed(scopeFile, "Scope", "iterVars_") {
		j2Qual, ok := ensureJ2gortImportInFile(scopeFile, ctx)
		if !ok {
			return sites, fmt.Errorf("failed to ensure j2gort import in Scope file")
		}
		decls, err := parseDeclsSnippet(ctx.Fset, buildScopeIterVarsMethodSrc("Scope", "*Var", j2Qual))
		if err != nil {
			return sites, err
		}
		scopeFile.Decls = append(scopeFile.Decls, decls...)
		ctx.MarkFileModified(scopeFile)
	}

	if needTypedScopeIter && typedScopeFile != nil && !fileHasPtrMethodNamed(typedScopeFile, "TypedScope", "iterVars_") {
		j2Qual, ok := ensureJ2gortImportInFile(typedScopeFile, ctx)
		if !ok {
			return sites, fmt.Errorf("failed to ensure j2gort import in TypedScope file")
		}
		decls, err := parseDeclsSnippet(ctx.Fset, buildScopeIterVarsMethodSrc("TypedScope", "*TypedVar", j2Qual))
		if err != nil {
			return sites, err
		}
		typedScopeFile.Decls = append(typedScopeFile.Decls, decls...)
		ctx.MarkFileModified(typedScopeFile)
	}

	return sites, nil
}

func buildScopeIterVarsMethodSrc(scopeTypeName, varPtrTypeSrc, j2Qual string) string {
	return fmt.Sprintf(`package p

func (this_ *%[1]s) iterVars_() scopeVarCursor_[%[2]s] {
	if this_ == nil {
		%[3]s.ThrowNPE()
	}
	if this_.vars != nil {
		return scopeVarCursor_[%[2]s]{useMap: true, mapCur: this_.vars.IterValues()}
	}
	if this_.var0 == nil {
		return scopeVarCursor_[%[2]s]{}
	}
	var c scopeVarCursor_[%[2]s]
	c.vars[0] = this_.var0
	if this_.var1 == nil {
		c.n = 1
		return c
	}
	c.vars[1] = this_.var1
	if this_.var2 == nil {
		c.n = 2
		return c
	}
	c.vars[2] = this_.var2
	if this_.var3 == nil {
		c.n = 3
		return c
	}
	c.vars[3] = this_.var3
	if this_.var4 == nil {
		c.n = 4
		return c
	}
	c.vars[4] = this_.var4
	if this_.var5 == nil {
		c.n = 5
		return c
	}
	c.vars[5] = this_.var5
	c.n = 6
	return c
}
`, scopeTypeName, varPtrTypeSrc, j2Qual)
}

func verifyInlineVarScopeLayout(scopeName, varName string, ctx *PassContext) (*types.Named, *types.Named, *ast.File, bool) {
	sc := ctx.Pkg.Scope()
	sObj, okS := sc.Lookup(scopeName).(*types.TypeName)
	vObj, okV := sc.Lookup(varName).(*types.TypeName)
	if !okS || !okV || sObj == nil || vObj == nil {
		return nil, nil, nil, false
	}
	scopeNamed, okSN := sObj.Type().(*types.Named)
	varNamed, okVN := vObj.Type().(*types.Named)
	if !okSN || !okVN || scopeNamed == nil || varNamed == nil {
		return nil, nil, nil, false
	}
	st, okSt := scopeNamed.Underlying().(*types.Struct)
	if !okSt {
		return nil, nil, nil, false
	}
	wantSlots := map[string]bool{"var0": false, "var1": false, "var2": false, "var3": false, "var4": false, "var5": false}
	hasVarsMap := false
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if _, isSlot := wantSlots[f.Name()]; isSlot {
			pt, isP := types.Unalias(f.Type()).(*types.Pointer)
			if !isP || namedOrigin(pt.Elem()) != varNamed.Origin() {
				return nil, nil, nil, false
			}
			wantSlots[f.Name()] = true
			continue
		}
		if f.Name() == "vars" {
			pt, isP := types.Unalias(f.Type()).(*types.Pointer)
			if !isP {
				return nil, nil, nil, false
			}
			om, isN := types.Unalias(pt.Elem()).(*types.Named)
			if !isN || om.Origin() == nil || om.Origin().Obj() == nil ||
				!isJ2gortPkg(om.Origin().Obj().Pkg()) || om.Origin().Obj().Name() != "OrderedMap" {
				return nil, nil, nil, false
			}
			targs := om.TypeArgs()
			if targs == nil || targs.Len() != 2 {
				return nil, nil, nil, false
			}
			kb, isB := types.Unalias(targs.At(0)).(*types.Basic)
			vp, isVP := types.Unalias(targs.At(1)).(*types.Pointer)
			if !isB || kb.Kind() != types.String || !isVP || namedOrigin(vp.Elem()) != varNamed.Origin() {
				return nil, nil, nil, false
			}
			hasVarsMap = true
		}
	}
	// order-independent: conjunctive check over var0..var5 inline slot presence
	for _, ok := range wantSlots {
		if !ok {
			return nil, nil, nil, false
		}
	}
	if !hasVarsMap {
		return nil, nil, nil, false
	}
	for _, f := range ctx.Files {
		if fileHasTypeSpec(f, scopeName) {
			return scopeNamed.Origin(), varNamed.Origin(), f, true
		}
	}
	return nil, nil, nil, false
}

func fileHasTypeSpec(file *ast.File, typeName string) bool {
	if file == nil {
		return false
	}
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, sp := range gd.Specs {
			if ts, ok := sp.(*ast.TypeSpec); ok && ts.Name != nil && ts.Name.Name == typeName {
				return true
			}
		}
	}
	return false
}

func fileHasPtrMethodNamed(file *ast.File, typeName, methodName string) bool {
	if file == nil {
		return false
	}
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if ok && fd.Name != nil && fd.Name.Name == methodName && isPtrMethodOnType(fd, typeName) {
			return true
		}
	}
	return false
}

func matchScopeGetVarIterableCursorDefine(
	stmt ast.Stmt,
	scopeNamed, typedScopeNamed *types.Named,
	allowScope, allowTypedScope bool,
	info *types.Info,
) (*ast.Ident, *types.Var, ast.Expr, bool, bool) {
	as, ok := stmt.(*ast.AssignStmt)
	if !ok || as.Tok != token.DEFINE || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
		return nil, nil, nil, false, false
	}
	curId, ok := as.Lhs[0].(*ast.Ident)
	if !ok || curId.Name == "" || curId.Name == "_" {
		return nil, nil, nil, false, false
	}
	curVar := varFromObj(info.Defs[curId])
	if curVar == nil || !isJ2gortIterCursorType(curVar.Type()) {
		return nil, nil, nil, false, false
	}
	outerCall, ok := stripParens(as.Rhs[0]).(*ast.CallExpr)
	if !ok || outerCall.Ellipsis.IsValid() || len(outerCall.Args) != 1 {
		return nil, nil, nil, false, false
	}
	callee := resolveCalledTypesFunc(outerCall.Fun, info)
	if callee == nil || !isJ2gortPkg(callee.Pkg()) || callee.Name() != "NewIterCursor" {
		return nil, nil, nil, false, false
	}
	iterArg := stripParens(outerCall.Args[0])
	if innerCall, isC := iterArg.(*ast.CallExpr); isC && !innerCall.Ellipsis.IsValid() && len(innerCall.Args) == 1 {
		if innerFn := resolveCalledTypesFunc(innerCall.Fun, info); innerFn != nil && isJ2gortPkg(innerFn.Pkg()) && innerFn.Name() == "CheckNotNull" {
			iterArg = stripParens(innerCall.Args[0])
		}
	}
	getIterCall, ok := iterArg.(*ast.CallExpr)
	if !ok || getIterCall.Ellipsis.IsValid() || len(getIterCall.Args) != 0 {
		return nil, nil, nil, false, false
	}
	sel, ok := stripParens(getIterCall.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || (sel.Sel.Name != "GetVarIterable" && sel.Sel.Name != "GetAllSymbols") {
		return nil, nil, nil, false, false
	}
	recvT := info.TypeOf(sel.X)
	if recvT == nil {
		return nil, nil, nil, false, false
	}
	pt, isP := types.Unalias(recvT).(*types.Pointer)
	if !isP {
		return nil, nil, nil, false, false
	}
	elemOrigin := namedOrigin(pt.Elem())
	if allowScope && scopeNamed != nil && elemOrigin == scopeNamed {
		return curId, curVar, sel.X, false, true
	}
	if allowTypedScope && typedScopeNamed != nil && elemOrigin == typedScopeNamed {
		return curId, curVar, sel.X, true, true
	}
	return nil, nil, nil, false, false
}

func matchBlankDiscardIdent(stmt ast.Stmt, targetVar *types.Var, info *types.Info) *ast.Ident {
	as, ok := stmt.(*ast.AssignStmt)
	if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 || !isIdentNamed(as.Lhs[0], "_") {
		return nil
	}
	id, ok := stripParens(as.Rhs[0]).(*ast.Ident)
	if !ok || varFromObj(info.Uses[id]) != targetVar {
		return nil
	}
	return id
}

func matchZeroArgMethodCallIdent(expr ast.Expr, methodName string, targetVar *types.Var, info *types.Info) *ast.Ident {
	call, ok := stripParens(expr).(*ast.CallExpr)
	if !ok || call.Ellipsis.IsValid() || len(call.Args) != 0 {
		return nil
	}
	sel, ok := stripParens(call.Fun).(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || sel.Sel.Name != methodName {
		return nil
	}
	id, ok := stripParens(sel.X).(*ast.Ident)
	if !ok || varFromObj(info.Uses[id]) != targetVar {
		return nil
	}
	return id
}

func rewriteScopeVarCursorLoopBody(
	funcBody *ast.BlockStmt,
	forStmt *ast.ForStmt,
	curDefId, curDiscardId, condCurId *ast.Ident,
	curVar *types.Var,
	wantVarNamed *types.Named,
	info *types.Info,
) bool {
	if forStmt == nil || forStmt.Body == nil || len(forStmt.Body.List) == 0 {
		return false
	}
	// Pattern 1: `rawElem := cur.Next()` (+ optional `_ = rawElem`) followed by `v := j2gort.CastPtr[Var](rawElem)` (or `var v *Var = ...`).
	if rawAssign, ok := forStmt.Body.List[0].(*ast.AssignStmt); ok &&
		rawAssign.Tok == token.DEFINE && len(rawAssign.Lhs) == 1 && len(rawAssign.Rhs) == 1 {
		if rawDefId, okID := rawAssign.Lhs[0].(*ast.Ident); okID && rawDefId.Name != "_" {
			rawVar := varFromObj(info.Defs[rawDefId])
			bodyNextCurId := matchZeroArgMethodCallIdent(rawAssign.Rhs[0], "Next", curVar, info)
			if rawVar != nil && bodyNextCurId != nil {
				castStmtIdx := 1
				var rawDiscardId *ast.Ident
				if castStmtIdx < len(forStmt.Body.List) {
					if dId := matchBlankDiscardIdent(forStmt.Body.List[castStmtIdx], rawVar, info); dId != nil {
						rawDiscardId = dId
						castStmtIdx++
					}
				}
				if castStmtIdx < len(forStmt.Body.List) {
					if rawUseId, rewriteCast := matchVarCastAssignFromRaw(forStmt.Body.List[castStmtIdx], rawVar, wantVarNamed, rawAssign.Rhs[0], info); rawUseId != nil {
						allowed := map[*ast.Ident]bool{
							curDefId:      true,
							condCurId:     true,
							bodyNextCurId: true,
							rawDefId:      true,
							rawUseId:      true,
						}
						if curDiscardId != nil {
							allowed[curDiscardId] = true
						}
						if rawDiscardId != nil {
							allowed[rawDiscardId] = true
						}
						if allVarRefsInSet(funcBody, curVar, rawVar, allowed, info) {
							rewriteCast()
							forStmt.Body.List = append(forStmt.Body.List[:0], forStmt.Body.List[castStmtIdx:]...)
							return true
						}
					}
				}
			}
		}
	}

	// Pattern 2: Direct `j2gort.CastPtr[Var](cur.Next())` uses inside `forStmt.Body` (not inside nested FuncLit).
	allowed := map[*ast.Ident]bool{
		curDefId:  true,
		condCurId: true,
	}
	if curDiscardId != nil {
		allowed[curDiscardId] = true
	}
	parents, _ := buildParentMap(forStmt.Body)
	type castRewrite struct {
		parent ast.Node
		oldExp ast.Expr
		newExp ast.Expr
	}
	var rewrites []castRewrite
	badBody := false
	ast.Inspect(forStmt.Body, func(n ast.Node) bool {
		if badBody || n == nil {
			return false
		}
		if fl, isFL := n.(*ast.FuncLit); isFL {
			if !allVarRefsInSet(fl.Body, curVar, nil, allowed, info) {
				badBody = true
			}
			return false
		}
		id, isID := n.(*ast.Ident)
		if !isID || varFromObj(info.Uses[id]) != curVar {
			return true
		}
		sel, isSel := parents[id].(*ast.SelectorExpr)
		if !isSel || sel.X != id || sel.Sel == nil || sel.Sel.Name != "Next" {
			badBody = true
			return false
		}
		nextCall, isCall := parents[sel].(*ast.CallExpr)
		if !isCall || nextCall.Fun != sel || nextCall.Ellipsis.IsValid() || len(nextCall.Args) != 0 {
			badBody = true
			return false
		}
		castNode := parents[nextCall]
		castCall, isCast := castNode.(*ast.CallExpr)
		if !isCast || !isJ2gortCastPtrToNamed(castCall, nextCall, wantVarNamed, info) {
			badBody = true
			return false
		}
		allowed[id] = true
		rewrites = append(rewrites, castRewrite{parent: parents[castCall], oldExp: castCall, newExp: nextCall})
		return true
	})
	if badBody || len(rewrites) == 0 || !allVarRefsInSet(funcBody, curVar, nil, allowed, info) {
		return false
	}
	for _, rw := range rewrites {
		if !replaceDirectChildExpr(rw.parent, rw.oldExp, rw.newExp) {
			return false
		}
	}
	return true
}

func matchVarCastAssignFromRaw(
	stmt ast.Stmt,
	rawVar *types.Var,
	wantVarNamed *types.Named,
	nextCallExpr ast.Expr,
	info *types.Info,
) (*ast.Ident, func()) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if (s.Tok != token.DEFINE && s.Tok != token.ASSIGN) || len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			return nil, nil
		}
		call, ok := stripParens(s.Rhs[0]).(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return nil, nil
		}
		rawId, ok := stripParens(call.Args[0]).(*ast.Ident)
		if !ok || varFromObj(info.Uses[rawId]) != rawVar {
			return nil, nil
		}
		if !isJ2gortCastPtrToNamed(call, call.Args[0], wantVarNamed, info) {
			return nil, nil
		}
		return rawId, func() { s.Rhs[0] = nextCallExpr }

	case *ast.DeclStmt:
		gd, ok := s.Decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR || len(gd.Specs) != 1 {
			return nil, nil
		}
		vs, ok := gd.Specs[0].(*ast.ValueSpec)
		if !ok || len(vs.Names) != 1 || len(vs.Values) != 1 {
			return nil, nil
		}
		call, ok := stripParens(vs.Values[0]).(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return nil, nil
		}
		rawId, ok := stripParens(call.Args[0]).(*ast.Ident)
		if !ok || varFromObj(info.Uses[rawId]) != rawVar {
			return nil, nil
		}
		if !isJ2gortCastPtrToNamed(call, call.Args[0], wantVarNamed, info) {
			return nil, nil
		}
		return rawId, func() { vs.Values[0] = nextCallExpr }
	}
	return nil, nil
}

func isJ2gortCastPtrToNamed(call *ast.CallExpr, arg ast.Expr, wantNamed *types.Named, info *types.Info) bool {
	if call == nil || call.Ellipsis.IsValid() || len(call.Args) != 1 || call.Args[0] != arg {
		return false
	}
	callee := resolveCalledTypesFunc(call.Fun, info)
	if callee == nil || !isJ2gortPkg(callee.Pkg()) || callee.Name() != "CastPtr" {
		return false
	}
	tv, ok := info.Types[call]
	if !ok || tv.Type == nil {
		return false
	}
	pt, isP := types.Unalias(tv.Type).(*types.Pointer)
	return isP && namedOrigin(pt.Elem()) == wantNamed
}

func allVarRefsInSet(root ast.Node, v1, v2 *types.Var, allowed map[*ast.Ident]bool, info *types.Info) bool {
	ok := true
	ast.Inspect(root, func(n ast.Node) bool {
		if !ok || n == nil {
			return false
		}
		id, isID := n.(*ast.Ident)
		if !isID {
			return true
		}
		if defV := varFromObj(info.Defs[id]); (defV != nil && (defV == v1 || (v2 != nil && defV == v2))) && !allowed[id] {
			ok = false
			return false
		}
		if useV := varFromObj(info.Uses[id]); (useV != nil && (useV == v1 || (v2 != nil && useV == v2))) && !allowed[id] {
			ok = false
			return false
		}
		return true
	})
	return ok
}

func replaceDirectChildExpr(parent ast.Node, oldExpr, newExpr ast.Expr) bool {
	switch p := parent.(type) {
	case *ast.AssignStmt:
		for i := range p.Rhs {
			if p.Rhs[i] == oldExpr {
				p.Rhs[i] = newExpr
				return true
			}
		}
	case *ast.ValueSpec:
		for i := range p.Values {
			if p.Values[i] == oldExpr {
				p.Values[i] = newExpr
				return true
			}
		}
	case *ast.CallExpr:
		for i := range p.Args {
			if p.Args[i] == oldExpr {
				p.Args[i] = newExpr
				return true
			}
		}
	case *ast.ExprStmt:
		if p.X == oldExpr {
			p.X = newExpr
			return true
		}
	case *ast.SelectorExpr:
		if p.X == oldExpr {
			p.X = newExpr
			return true
		}
	case *ast.ParenExpr:
		if p.X == oldExpr {
			p.X = newExpr
			return true
		}
	case *ast.ReturnStmt:
		for i := range p.Results {
			if p.Results[i] == oldExpr {
				p.Results[i] = newExpr
				return true
			}
		}
	}
	return false
}
