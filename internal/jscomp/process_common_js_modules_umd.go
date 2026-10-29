// Port of com/google/javascript/jscomp/ProcessCommonJSModules.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ProcessCommonJSModules_IsCommonJsDynamicImportCallback recognizes webpack dynamic import callbacks.
//
// Java: ProcessCommonJSModules.isCommonJsDynamicImportCallback(Node, ModuleLoader.ResolutionMode)
func ProcessCommonJSModules_IsCommonJsDynamicImportCallback(
	n rhino.Node, resolutionMode ModuleLoader_ResolutionMode) bool {
	if n.IsNil() || resolutionMode != ModuleLoader_ResolutionMode_WEBPACK {
		return false
	}
	return n.IsFunction() && isWebpackRequireEnsureCallback(n)
}

var processCommonJSModules_PROMISE_ALL = rhino.QualifiedName_Of("Promise.all")

func isWebpackRequireEnsureCallback(fnc rhino.Node) bool {
	jutil.CheckArgument(fnc.IsFunction())
	if fnc.GetParent().IsNil() {
		return false
	}

	callParent, ok := ast.AsCall(fnc.GetParent())
	if !ok {
		return false
	}
	then, ok := callParent.Callee().AsGetProp()
	if !ok {
		return false
	}
	ensureCall, ok := then.Object().AsCall()
	if !ok {
		return false
	}
	callParentTarget := ensureCall.Callee().Node()

	if callParentTarget.MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE+".e") &&
		then.Property() == "then" {
		return true
	}
	if processCommonJSModules_PROMISE_ALL.Matches(callParentTarget) {
		firstArg, _ := ensureCall.FirstArg()
		array, ok := ast.AsArrayLit(firstArg.Node())
		if !ok {
			return false
		}
		allElementsAreDynamicImports := false
		for arrayItem := range array.Elements() {
			call, ok := ast.AsCall(arrayItem.Node())
			if !ok {
				return false
			}
			if _, ok := call.OnlyArg(); !ok ||
				!call.Callee().Node().MatchesQualifiedName(processCommonJSModules_WEBPACK_REQUIRE+".e") {
				return false
			}
			allElementsAreDynamicImports = true
		}
		return allElementsAreDynamicImports
	}
	return false
}

type commonJsUmdPattern struct {
	ifRoot       rhino.Node
	activeBranch rhino.Node
}

func (p *ProcessCommonJSModules) removeIIFEWrapper(root rhino.Node) bool {
	jutil.CheckState(root.IsScript())
	exprStmt, ok := processCommonJSModules_onlyExprStatement(ast.MustScript(root))
	if !ok {
		return false
	}

	expr := exprStmt.Expr()
	if not, ok := expr.AsUnaryOf(rhino.Token_NOT); ok {
		expr = not.Operand()
	}

	iife, ok := expr.AsCall()
	if !ok {
		return false
	}
	call := iife.Node()

	var fnc rhino.Node
	if function, ok := iife.Callee().AsFunction(); ok {
		fnc = function.Node()
	} else if callMethod, ok := iife.Callee().AsGetProp(); ok && callMethod.Property() == "call" {
		function, ok := callMethod.Object().AsFunction()
		if !ok {
			return false
		}
		fnc = function.Node()

		thisArg, ok := iife.FirstArg()
		if !ok || !(thisArg.Node().IsThis() || thisArg.Node().MatchesQualifiedName(processCommonJSModules_EXPORTS)) {
			return false
		}
	} else {
		return false
	}

	if NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fnc) {
		return false
	}

	ci := p.compiler.GetInput(root.GetInputId())
	modulePath := ci.GetPath()
	if modulePath == nil {
		return false
	}

	iifeLabel := modulePath.ToModuleName() + "_iifeWrapper"

	mutator := NewFunctionToBlockMutator(p.compiler, p.compiler.GetUniqueNameIdSupplier())
	block := mutator.MutateWithoutRenaming(iifeLabel, fnc, call, "", false, false)
	root.RemoveChildren()
	root.AddChildrenToFront(block.RemoveChildren())
	p.reportNestedScopesDeleted(fnc)
	p.compiler.ReportChangeToEnclosingScope(root)

	return true
}

func processCommonJSModules_onlyExprStatement(script ast.Script) (ast.ExprStmt, bool) {
	item, ok := script.FirstItem()
	for ok && item.Node().IsEmpty() {
		item, ok = item.Next()
	}
	if !ok || !item.IsLast() {
		return ast.ExprStmt{}, false
	}
	return ast.AsExprStmt(item.Node())
}

func (p *ProcessCommonJSModules) removeWebpackModuleShim(root rhino.Node) {
	jutil.CheckState(root.IsScript())
	exprStmt, ok := processCommonJSModules_onlyExprStatement(ast.MustScript(root))
	if !ok {
		return
	}

	iife, ok := exprStmt.Expr().AsCall()
	if !ok {
		return
	}
	call := iife.Node()

	function, ok := iife.Callee().AsFunction()
	if !ok {
		return
	}
	fnc := function.Node()

	var param rhino.Node
	paramNumber := 0
	for typedParam := range function.Params().Params() {
		paramNumber++
		if name, ok := typedParam.AsName(); ok && name.Ident() == processCommonJSModules_MODULE {
			param = name.Node()
			break
		}
	}
	if param.IsNil() {
		return
	}

	// Java quirk kept: the argument is read by child index, which counts the
	// callee; for a free call the index is one lower, so it reads the argument
	// before the one of the `module` parameter.
	argIndex := paramNumber + 1
	if call.GetBooleanProp(rhino.Node_FREE_CALL) {
		argIndex = paramNumber
	}
	arg := call.GetChildAtIndex(argIndex)
	if arg.IsNil() {
		return
	}

	argCall, ok := ast.AsCall(arg)
	if !ok {
		return
	}
	argCallTarget := argCall.Callee().Node()
	if moduleArg, ok := argCall.FirstArg(); ok &&
		p.isCommonJsImport(argCallTarget) &&
		moduleArg.Node().MatchesName(processCommonJSModules_MODULE) {
		importPath := p.getCommonJsImportPath(argCallTarget)
		modulePath := p.compiler.
			GetInput(root.GetInputId()).
			GetPath().
			ResolveJsModule(
				importPath, processCommonJSModules_sourceFileName(arg), arg.GetLineno(), arg.GetCharno())
		if modulePath == nil {
			return
		}
		if strings.Contains(modulePath.String(), "/buildin/module.js") {
			arg.Detach()
			param.Detach()
			p.compiler.ReportChangeToChangeScope(fnc)
			p.compiler.ReportChangeToEnclosingScope(fnc)
		}
	}
}

type commonJsUmdTestInfo struct {
	enclosingIf  rhino.Node
	activeBranch rhino.Node
}

func hasUmdTestReference(node rhino.Node) bool {
	switch node.GetToken() {
	case rhino.Token_NAME:
		s := node.GetString()
		if s == processCommonJSModules_MODULE || s == "define" {
			return true
		}
	case rhino.Token_GETPROP:
		if processCommonJSModules_WINDOW_DEFINE.Matches(node) {
			return true
		}
	case rhino.Token_STRINGLIT:
		if node.GetParent().IsIn() && node.GetString() == "amd" {
			return true
		}
	}
	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if hasUmdTestReference(child) {
			return true
		}
	}
	return false
}

func getOutermostUmdTest(n rhino.Node) (commonJsUmdTestInfo, bool) {
	if n.IsNil() || NodeUtil_IsTopLevel(n) || n.IsFunction() {
		return commonJsUmdTestInfo{}, false
	}
	parent := n.GetParent()
	if parent.IsNil() {
		return commonJsUmdTestInfo{}, false
	}

	if cond, ok := processCommonJSModules_ifOrHookCondition(parent); ok {
		if umdTestInfo, hasOuter := getOutermostUmdTest(parent); hasOuter {
			if ifs, ok := ast.AsIf(parent); ok &&
				ifs.Then().Is(n) &&
				parent.GetParent().IsBlock() &&
				parent.GetNext().IsNil() &&
				umdTestInfo.activeBranch == parent.GetParent() {
				return commonJsUmdTestInfo{enclosingIf: umdTestInfo.enclosingIf, activeBranch: n}, true
			}
			return umdTestInfo, true
		}

		_, isTrue := cond.AsTrue()
		if isTrue || hasUmdTestReference(cond.Node()) {
			return commonJsUmdTestInfo{enclosingIf: parent, activeBranch: n}, true
		}
		return commonJsUmdTestInfo{}, false
	}

	return getOutermostUmdTest(parent)
}

func isInIfTest(n rhino.Node) bool {
	for !n.IsNil() && !NodeUtil_IsTopLevel(n) && !n.IsFunction() {
		parent := n.GetParent()
		if parent.IsNil() {
			return false
		}
		if cond, ok := processCommonJSModules_ifOrHookCondition(parent); ok && cond.Is(n) {
			return true
		}
		n = parent
	}
	return false
}

func processCommonJSModules_ifOrHookCondition(n rhino.Node) (ast.Expr, bool) {
	if ifs, ok := ast.AsIf(n); ok {
		return ifs.Cond(), true
	}
	if hook, ok := ast.AsHook(n); ok {
		return hook.Cond(), true
	}
	return ast.Expr{}, false
}

func processCommonJSModules_isIfWithoutElse(n rhino.Node) bool {
	ifs, ok := ast.AsIf(n)
	if !ok {
		return false
	}
	_, hasElse := ifs.Else()
	return !hasElse
}

func (f *commonJsFindImportsAndExports) replaceUmdPatterns() bool {
	compiler := f.outer.compiler
	needsRetraverse := false
	for _, umdPattern := range f.umdPatterns {
		if NodeUtil_GetEnclosingScript(umdPattern.ifRoot).IsNil() {
			f.outer.reportNestedScopesDeleted(umdPattern.ifRoot)
			continue
		}

		parent := umdPattern.ifRoot.GetParent()
		newNode := umdPattern.activeBranch

		if newNode.IsNil() {
			umdPattern.ifRoot.Detach()
			f.outer.reportNestedScopesDeleted(umdPattern.ifRoot)
			compiler.ReportChangeToEnclosingScope(parent)
			needsRetraverse = true
			continue
		}

		if block, ok := ast.AsBlock(umdPattern.activeBranch); ok {
			if only, ok := block.OnlyStatement(); ok {
				newNode = only.Detach().Node()
			} else {
				newNode.Detach()
			}
		} else {
			newNode.Detach()
		}
		needsRetraverse = true
		umdPattern.ifRoot.ReplaceWith(newNode)
		f.outer.reportNestedScopesDeleted(umdPattern.ifRoot)
		if changeScope := ChangeTracker_GetEnclosingChangeScopeRoot(newNode); !changeScope.IsNil() {
			compiler.ReportChangeToEnclosingScope(newNode)
		}

		block := parent
		if block.IsExprResult() {
			block = block.GetParent()
		}

		// Detect UMD Factory Patterns and inline the functions
		if block.IsBlock() && block.GetParent().IsFunction() &&
			block.GetGrandparent().IsCall() &&
			parent.HasOneChild() {
			enclosingFnCall := block.GetGrandparent()
			fn := block.GetParent()

			enclosingScript := NodeUtil_GetEnclosingScript(enclosingFnCall)
			if enclosingScript.IsNil() {
				continue
			}
			ci := compiler.GetInput(enclosingScript.GetInputId())
			modulePath := ci.GetPath()
			if modulePath == nil {
				continue
			}
			needsRetraverse = true
			factoryLabel := modulePath.ToModuleName() + "_factory" + compiler.GetUniqueIdSupplier().GetUniqueId(ci)

			mutator := NewFunctionToBlockMutator(compiler, compiler.GetUniqueNameIdSupplier())
			newStatements := mutator.MutateWithoutRenaming(factoryLabel, fn, enclosingFnCall, "", false, false)

			if inlinedFn, exprStmt, ok := processCommonJSModules_inlinedFactory(newStatements); ok {
				expr := exprStmt.Expr()
				var call rhino.Node
				var assignedName string
				var exprAssign ast.Assign
				if assign, ok := expr.AsAssign(); ok {
					if valueCall, ok := assign.Value().AsCall(); ok {
						exprAssign = assign
						call = valueCall.Node()
						assignedName = modulePath.ToModuleName() + "_iife" + compiler.GetUniqueIdSupplier().GetUniqueId(ci)
					}
				} else if exprCall, ok := expr.AsCall(); ok {
					call = exprCall.Node()
				}

				if !call.IsNil() {
					newStatements = mutator.MutateWithoutRenaming(
						factoryLabel, inlinedFn.Node(), call, assignedName, false, false)
					if assignedName != "" {
						targetQName, _ := exprAssign.Target().Node().GetQualifiedNameString()
						newName := rhino.IR_Var(
							NodeUtil_NewNameWithOriginalName(compiler, assignedName, fn, targetQName)).
							SrcrefTree(fn)
						if first, ok := processCommonJSModules_firstAssignTo(newStatements, assignedName); ok {
							newName.GetFirstChild().AddChildToFront(first.Value().Detach().Node())
							newStatements.GetFirstChild().ReplaceWith(newName)
						} else {
							newStatements.AddChildToFront(newName)
						}
						exprAssign.Value().Node().ReplaceWith(newName.GetFirstChild().CloneNode())
						newStatements.AddChildToBack(exprStmt.Detach().Node())
					}
				}
			}

			callRoot := enclosingFnCall.GetParent()
			if callRoot.IsNot() {
				callRoot = callRoot.GetParent()
			}
			if callRoot.IsExprResult() {
				callRootParent := callRoot.GetParent()
				callRootParent.AddChildrenAfter(newStatements.RemoveChildren(), callRoot)
				callRoot.Detach()
				f.outer.reportNestedScopesChanged(callRootParent)
				compiler.ReportChangeToEnclosingScope(callRootParent)
				f.outer.reportNestedScopesDeleted(enclosingFnCall)
			} else {
				umdPattern.ifRoot.ReplaceWith(newNode)
				compiler.ReportChangeToEnclosingScope(newNode)
				f.outer.reportNestedScopesDeleted(umdPattern.ifRoot)
			}
		}
	}
	return needsRetraverse
}

func processCommonJSModules_inlinedFactory(newStatements rhino.Node) (ast.Function, ast.ExprStmt, bool) {
	block, ok := ast.AsBlock(newStatements)
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	first, ok := block.FirstStatement()
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	second, ok := first.Next()
	if !ok || !second.IsLast() {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	declaration, ok := first.AsNameDeclarationOf(rhino.Token_VAR)
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	firstDeclarator, _ := declaration.FirstDeclarator()
	name, ok := firstDeclarator.AsName()
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	init, ok := name.Init()
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	inlinedFn, ok := init.AsFunction()
	if !ok {
		return ast.Function{}, ast.ExprStmt{}, false
	}
	exprStmt, ok := second.AsExprStmt()
	return inlinedFn, exprStmt, ok
}

func processCommonJSModules_firstAssignTo(block rhino.Node, name string) (ast.Assign, bool) {
	exprStmt, ok := ast.AsExprStmt(block.GetFirstChild())
	if !ok {
		return ast.Assign{}, false
	}
	assign, ok := exprStmt.Expr().AsAssign()
	if !ok {
		return ast.Assign{}, false
	}
	target, ok := assign.Target().AsName()
	return assign, ok && target.Ident() == name
}

func (p *ProcessCommonJSModules) reportNestedScopesDeleted(n rhino.Node) {
	if n.IsFunction() {
		p.compiler.ReportFunctionDeleted(n)
	}
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		p.reportNestedScopesDeleted(child)
	}
}

func (p *ProcessCommonJSModules) reportNestedScopesChanged(n rhino.Node) {
	if n.IsFunction() {
		p.compiler.ReportChangeToChangeScope(n)
	}
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		p.reportNestedScopesChanged(child)
	}
}

func hasUmdPattern(umdPatterns []commonJsUmdPattern, n rhino.Node) bool {
	for i := range umdPatterns {
		if umdPatterns[i].ifRoot == n {
			return true
		}
	}
	return false
}
