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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ProcessCommonJSModules_RewriteModule traverses a file and rewrite all
// references to imported names directly to the targeted module name.
//
// If a file is a CommonJS module, rewrite export statements. Typically
// exports create an alias - the rewriting tries to avoid such aliases.
type ProcessCommonJSModules_RewriteModule struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	outer *ProcessCommonJSModules // the enclosing instance

	allowFullRewrite          bool
	exports                   []*ProcessCommonJSModules_ExportInfo
	imports                   []rhino.Node
	rewrittenClassExpressions []rhino.Node
	functionsToHoist          []rhino.Node
	defaultExportIsConst      bool
}

// Java: ProcessCommonJSModules.RewriteModule(boolean, ImmutableCollection<ExportInfo>, boolean)
func NewProcessCommonJSModules_RewriteModule(
	outer *ProcessCommonJSModules,
	allowFullRewrite bool,
	exports []*ProcessCommonJSModules_ExportInfo,
	defaultExportIsConst bool) *ProcessCommonJSModules_RewriteModule {
	r := &ProcessCommonJSModules_RewriteModule{outer: outer}
	r.imports = []rhino.Node{}
	r.rewrittenClassExpressions = []rhino.Node{}
	r.functionsToHoist = []rhino.Node{}
	r.allowFullRewrite = allowFullRewrite
	r.exports = exports
	r.defaultExportIsConst = defaultExportIsConst
	return r
}

// Java: ProcessCommonJSModules.RewriteModule.visit(NodeTraversal, Node, Node)
func (r *ProcessCommonJSModules_RewriteModule) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	compiler := r.outer.compiler
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		// Class names can't be changed during the middle of a traversal. Unlike functions,
		// the name can be the EMPTY token rather than just a zero length string.
		for _, clazz := range r.rewrittenClassExpressions {
			clazz.GetFirstChild().ReplaceWith(rhino.IR_Empty().Srcref(clazz.GetFirstChild()))
			t.ReportCodeChange()
		}

		ci := compiler.GetInput(n.GetInputId())
		moduleName := *ProcessCommonJSModules_GetModuleName(ci)

		// If a function is the direct module export, move it to the top.
		for i := 1; i < len(r.functionsToHoist); i++ {
			if r.functionsToHoist[i].
				GetFirstFirstChild().
				MatchesQualifiedName(r.outer.GetBasePropertyImport(moduleName)) {
				fncVar := r.functionsToHoist[i]
				r.functionsToHoist = slices.Delete(r.functionsToHoist, i, i+1)
				r.functionsToHoist = slices.Insert(r.functionsToHoist, 0, fncVar)
				break
			}
		}

		// Hoist functions in reverse order so that they maintain the same relative
		// order after hoisting.
		for i := len(r.functionsToHoist) - 1; i >= 0; i-- {
			functionExpr := r.functionsToHoist[i]
			scopeRoot := t.GetClosestHoistScopeRoot()
			insertionPoint := scopeRoot.GetFirstChild()
			if insertionPoint.IsNil() ||
				!(insertionPoint.IsVar() &&
					insertionPoint.GetFirstChild().GetString() == moduleName) {
				insertionPoint = rhino.NilNode
			}

			if insertionPoint.IsNil() {
				if scopeRoot.GetFirstChild() != functionExpr {
					scopeRoot.AddChildToFront(functionExpr.Detach())
				}
			} else if insertionPoint != functionExpr && insertionPoint.GetNext() != functionExpr {
				functionExpr.Detach().InsertAfter(insertionPoint)
			}
		}

		for _, export := range r.exports {
			r.VisitExport(t, export)
		}

		for _, require := range r.imports {
			r.VisitRequireCall(t, require)
		}
	case rhino.Token_CALL:
		if r.outer.IsCommonJsImport(n) {
			r.imports = append(r.imports, n)
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		// Multiple declarations need split apart so that they can be refactored into
		// property assignments or removed altogether. Don't split declarations for
		// ES module export calls as it breaks AST and ES modules shouldn't be affected at all
		// by this pass.
		if n.HasMoreThanOneChild() && !NodeUtil_IsAnyFor(parent) && !parent.IsExport() {
			vars := r.SplitMultipleDeclarations(n)
			t.ReportCodeChange()
			for _, var_ := range vars {
				r.Visit(t, var_.GetFirstChild(), var_)
			}
		}

		// UMD Inlining can shadow global variables - these are just removed.
		//
		// var exports = exports;
		if processCommonJSModules_isSelfAlias(ast.MustNameDeclaration(n)) {
			n.Detach()
			t.ReportCodeChange()
			return
		}
	case rhino.Token_NAME:
		// If this is a name declaration with multiple names, it will be split apart when
		// the parent is visited and then revisit the children.
		if NodeUtil_IsNameDeclaration(n.GetParent()) &&
			n.GetParent().HasMoreThanOneChild() &&
			!NodeUtil_IsAnyFor(n.GetGrandparent()) {
			break
		}

		qName := n.GetQualifiedName()
		if qName == nil {
			break
		}
		nameDeclaration := t.GetScope().GetVar(*qName)
		if nameDeclaration != nil {
			if NodeUtil_IsLhsByDestructuring(n) {
				r.MaybeUpdateName(t, n, nameDeclaration)
			} else if !nameDeclaration.GetNode().IsNil() &&
				jutil.ObjectsEquals(nameDeclaration.GetNode().GetInputId(), n.GetInputId()) {
				// Avoid renaming a shadowed global
				//
				// var angular = angular;  // value is global ref
				enclosingDeclaration :=
					NodeUtil_GetEnclosingNode(
						n, func(node rhino.Node) bool { return node == nameDeclaration.GetNameNode() })

				if enclosingDeclaration.IsNil() ||
					enclosingDeclaration == n ||
					nameDeclaration.GetScope().(*Scope) != t.GetScope() {
					r.MaybeUpdateName(t, n, nameDeclaration)
				}
			}
			// Replace loose "module" references with an object literal
		} else if r.allowFullRewrite &&
			n.GetString() == processCommonJSModules_MODULE &&
			!(n.GetParent().IsGetProp() ||
				n.GetParent().IsGetElem() ||
				n.GetParent().IsTypeOf()) {
			objectLit := rhino.IR_Objectlit().Srcref(n)
			n.ReplaceWith(objectLit)
			t.ReportCodeChangeWithNode(objectLit)
		}
	case rhino.Token_GETPROP:
		if n.MatchesQualifiedName(processCommonJSModules_MODULE+".id") ||
			n.MatchesQualifiedName(processCommonJSModules_MODULE+".filename") {
			v := t.GetScope().GetVar(processCommonJSModules_MODULE)
			if v == nil || v.IsExtern() {
				n.ReplaceWith(rhino.IR_String(t.GetInput().GetPath().String()).Srcref(n))
			}
		} else if object, ok := ast.MustGetProp(n).Object().AsName(); r.allowFullRewrite &&
			ok && object.Ident() == processCommonJSModules_MODULE &&
			!n.MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) {
			v := t.GetScope().GetVar(processCommonJSModules_MODULE)
			if v == nil || v.IsExtern() {
				object.Node().ReplaceWith(rhino.IR_Objectlit().Srcref(object.Node()))
			}
		}
	case rhino.Token_TYPEOF:
		if operand, ok := ast.MustUnary(n).Operand().AsName(); r.allowFullRewrite &&
			ok && (operand.Ident() == processCommonJSModules_MODULE ||
			operand.Ident() == processCommonJSModules_EXPORTS) {
			v := t.GetScope().GetVar(operand.Ident())
			if v == nil || v.IsExtern() {
				n.ReplaceWith(rhino.IR_String("object"))
			}
		}
	default:
	}

	r.FixTypeAnnotationsForNode(t, n)
}

// processCommonJSModules_isSelfAlias reports whether declaration starts by
// declaring a name as itself: `var exports = exports;`.
//
// Go-only: Java tests the first declarator and its child positionally.
func processCommonJSModules_isSelfAlias(declaration ast.NameDeclaration) bool {
	first, _ := declaration.FirstDeclarator()
	name, ok := first.AsName()
	if !ok {
		return false
	}
	init, ok := name.Init()
	if !ok {
		return false
	}
	value, ok := init.AsName()
	return ok && value.Ident() == name.Ident()
}

// Java: ProcessCommonJSModules.RewriteModule.fixTypeAnnotationsForNode(NodeTraversal, Node)
func (r *ProcessCommonJSModules_RewriteModule) FixTypeAnnotationsForNode(t *NodeTraversal, n rhino.Node) {
	info := n.GetJSDocInfo()
	if info != nil {
		for it := info.TypeExprIter(); it.Next(); {
			r.FixTypeNode(t, it.Node())
		}
	}
}

// VisitRequireCall visits require calls. Rewrite require statements to be a
// direct reference to name of require module. By this point all references
// to the import alias should have already been renamed.
//
// Java: ProcessCommonJSModules.RewriteModule.visitRequireCall(NodeTraversal, Node)
func (r *ProcessCommonJSModules_RewriteModule) VisitRequireCall(t *NodeTraversal, require rhino.Node) {
	moduleName := r.outer.GetImportedModuleName(t, require)
	// The first argument: the path, or the `null` of __webpack_require__.bind(null, path).
	firstArg, _ := ast.MustCall(require).FirstArg()
	moduleRef :=
		NodeUtil_NewQName(r.outer.compiler, r.outer.GetBasePropertyImportWithRequireCall(moduleName, require)).
			SrcrefTree(firstArg.Node())
	require.ReplaceWith(moduleRef)

	t.ReportCodeChange()
}

// VisitExport visits export statements. Export statements can be either a
// direct assignment: module.exports = foo or a property assignment:
// module.exports.foo = foo; exports.foo = foo;
//
// Java: ProcessCommonJSModules.RewriteModule.visitExport(NodeTraversal, ExportInfo)
func (r *ProcessCommonJSModules_RewriteModule) VisitExport(t *NodeTraversal, export *ProcessCommonJSModules_ExportInfo) {
	compiler := r.outer.compiler
	root := r.outer.GetBaseQualifiedNameNode(export.node)
	rValue := NodeUtil_GetRValueOfLValue(root)

	// For object literal assignments to module.exports, convert them to
	// individual property assignments.
	//
	//     module.exports = { foo: bar};
	//
	// becomes
	//
	//     module.exports = {};
	//     module.exports.foo = bar;
	if processCommonJSModules_MODULE_EXPORTS.Matches(root) {
		if !rValue.IsNil() &&
			rValue.IsObjectLit() &&
			root.GetParent().IsAssign() &&
			root.GetGrandparent().IsExprResult() {
			if r.ExpandObjectLitAssignment(t, root, export.scope) {
				return
			}
		}
	}

	moduleName := *ProcessCommonJSModules_GetModuleName(t.GetInput())
	moduleInitialization := t.GetScope().GetVar(moduleName)

	// If this is an assignment to module.exports or exports, renaming
	// has already handled this case. Remove the export.
	var rValueVar *Var
	if !rValue.IsNil() && rValue.IsQualifiedName() {
		rValueVar = export.scope.GetVar(*rValue.GetQualifiedName())
	}

	if root.GetParent().IsAssign() &&
		root.GetGrandparent().IsExprResult() &&
		(!root.GetNext().IsNil() && (root.GetNext().IsName() || root.GetNext().IsGetProp())) &&
		root.GetGrandparent().IsExprResult() &&
		rValueVar != nil &&
		(NodeUtil_GetEnclosingScript(rValueVar.GetNameNode()).IsNil() ||
			(rValueVar.GetNameNode().HasParent() && !rValueVar.IsParam())) &&
		export.isInSupportedScope &&
		(rValueVar.GetNameNode().GetParent().IsNil() ||
			!NodeUtil_IsLhsByDestructuring(rValueVar.GetNameNode())) {
		root.GetGrandparent().Detach()
		t.ReportCodeChange()
		return
	}

	moduleName = moduleName + "." + processCommonJSModules_EXPORT_PROPERTY_NAME

	updatedExport :=
		NodeUtil_NewQNameWithOriginalName(compiler, moduleName, export.node, *export.node.GetQualifiedName())
	updatedExport.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
	exportIsConst :=
		r.defaultExportIsConst &&
			updatedExport.MatchesQualifiedName(
				r.outer.GetBasePropertyImport(*ProcessCommonJSModules_GetModuleName(t.GetInput()))) &&
			root == export.node &&
			NodeUtil_IsLValue(export.node)

	var changeScope rhino.Node

	if processCommonJSModules_MODULE_EXPORTS.Matches(root) &&
		!rValue.IsNil() &&
		export.scope.GetVar("module.exports") == nil &&
		root.GetParent().IsAssign() &&
		root.GetGrandparent().IsExprResult() &&
		moduleInitialization == nil {
		// Rewrite "module.exports = foo;" to "var moduleName = {default: foo};"
		parent := root.GetParent()
		exportName := rhino.IR_ExprResult(rhino.IR_Assign(updatedExport, rValue.Detach()))
		if exportIsConst {
			info := rhino.JSDocInfo_Builder_()
			info.RecordConstancy()
			exportName.GetFirstChild().SetJSDocInfo(info.Build())
		}
		parent.GetParent().ReplaceWith(exportName.SrcrefTree(root.GetParent()))
		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(parent)
	} else if !root.GetNext().IsNil() &&
		root.GetNext().IsName() &&
		rValueVar != nil &&
		rValueVar.IsGlobal() &&
		export.isInSupportedScope &&
		(rValueVar.GetNameNode().GetParent().IsNil() ||
			!NodeUtil_IsLhsByDestructuring(rValueVar.GetNameNode())) {
		// This is a where a module export assignment is used in a complex expression.
		// Before: `SOME_VALUE !== undefined && module.exports = SOME_VALUE`
		// After: `SOME_VALUE !== undefined && module$name`
		parent := root.GetParent()
		root.Detach()
		parent.ReplaceWith(root)
		if root == export.node {
			root = updatedExport
		}
		export.node.ReplaceWith(updatedExport)
		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(root)
	} else {
		// Other references to "module.exports" are just replaced with the module name.
		export.node.ReplaceWith(updatedExport)
		if updatedExport.GetParent().IsAssign() && exportIsConst {
			infoBuilder :=
				rhino.JSDocInfo_Builder_MaybeCopyFrom(updatedExport.GetParent().GetJSDocInfo())
			infoBuilder.RecordConstancy()
			updatedExport.GetParent().SetJSDocInfo(infoBuilder.Build())
		}

		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(updatedExport)
	}

	if !changeScope.IsNil() {
		compiler.ReportChangeToChangeScope(changeScope)
	}
}

// ExpandObjectLitAssignment: since CommonJS modules may have only a single
// export, it's common to see the export be an object pattern. We want to
// expand this to individual property assignments. If any individual
// property assignment has been renamed, it will be removed.
//
// We need to keep assignments which aren't names
//
// module.exports = { foo: bar, baz: function() {} }
//
// becomes
//
// module.exports.foo = bar; // removed later module.exports.baz = function() {};
//
// Java: ProcessCommonJSModules.RewriteModule.expandObjectLitAssignment(NodeTraversal, Node, Scope)
func (r *ProcessCommonJSModules_RewriteModule) ExpandObjectLitAssignment(t *NodeTraversal, export rhino.Node, scope *Scope) bool {
	jutil.CheckState(export.GetParent().IsAssign())
	insertionRef := export.GetGrandparent()
	jutil.CheckState(insertionRef.IsExprResult())
	insertionParent := insertionRef.GetParent()
	rhino.CheckNotNull(insertionParent)

	rValue := ast.MustObjectLit(NodeUtil_GetRValueOfLValue(export))
	member, ok := rValue.FirstMember()

	removedNodes := false
	for ok {
		if !processCommonJSModules_isNamespaceMember(member) {
			member, ok = member.Next()
			continue
		}
		key := member.Node()

		lhs := rhino.IR_Getprop(export.CloneTree(), key.GetString())
		var value rhino.Node
		if stringKey, ok := member.AsStringKey(); ok {
			value = stringKey.Value().Detach().Node()
		} else {
			value = ast.MustMemberFunctionDef(key).Function().Detach().Node()
		}

		expr := rhino.IR_ExprResult(rhino.IR_Assign(lhs, value)).SrcrefTreeIfMissing(key)
		expr.InsertAfter(insertionRef)
		newExport := NewProcessCommonJSModules_ExportInfo(lhs.GetFirstChild(), scope)
		r.VisitExport(t, newExport)

		// Export statements can be removed in visitExport
		if !expr.IsNil() && expr.HasParent() {
			insertionRef = expr
		}

		// Read the next member after visitExport, before detaching this one (Java's order).
		current := member
		member, ok = member.Next()
		current.Detach()
		removedNodes = true
	}

	if !rValue.HasMembers() {
		export.GetGrandparent().Detach()
		return true
	}

	if removedNodes {
		t.ReportCodeChangeWithNode(rValue.Node())
	}
	return false
}

// MaybeUpdateName: given a name reference, check to see if it needs
// renamed.
//
// We handle 3 main cases: 1. References to an import alias. These are
// replaced with a direct reference to the imported module. 2. Names which
// are exported. These are rewritten to be the export assignment directly. 3.
// Global names: If a name is global to the script, add a suffix so it
// doesn't collide with any other global.
//
// Rewriting case 1 is safe to perform on all files. Cases 2 and 3 can only
// be done if this file is a commonjs module.
//
// Java: ProcessCommonJSModules.RewriteModule.maybeUpdateName(NodeTraversal, Node, Var)
func (r *ProcessCommonJSModules_RewriteModule) MaybeUpdateName(t *NodeTraversal, n rhino.Node, var_ *Var) {
	jutil.CheckNotNull(var_)
	jutil.CheckState(n.IsName() || n.IsGetProp())
	jutil.CheckState(n.HasParent())
	importedModuleName := r.GetModuleImportName(t, var_.GetNode())
	name := *n.GetQualifiedName()

	// Check if the name refers to a alias for a require('foo') import.
	if importedModuleName != nil && n != var_.GetNode() {
		// Reference the imported name directly, rather than the alias
		r.UpdateNameReference(t, n, name, *importedModuleName, false, false)

	} else if r.allowFullRewrite {
		exportedName := r.GetExportedName(t, n, var_)

		// We need to exclude the alias created by the require import. We assume dead
		// code elimination will remove these later.
		if (n != var_.GetNode() || n.GetParent().IsClass()) &&
			exportedName == nil &&
			(var_ == nil ||
				var_.GetNameNode().GetParent().IsNil() ||
				!NodeUtil_IsLhsByDestructuring(var_.GetNameNode())) {
			// The name is actually the export reference itself.
			// This will be handled later by visitExports.
			if n.GetParent().IsClass() && n.GetParent().GetFirstChild() == n {
				r.rewrittenClassExpressions = append(r.rewrittenClassExpressions, n.GetParent())
			}

			return
		}

		// Check if the name is used as an export
		if importedModuleName == nil &&
			exportedName != nil &&
			*exportedName != name &&
			!var_.IsParam() &&
			(var_.GetNameNode().GetParent().IsNil() ||
				!NodeUtil_IsLhsByDestructuring(var_.GetNameNode())) {
			exportPropIsConst :=
				r.defaultExportIsConst &&
					r.outer.GetBasePropertyImport(*ProcessCommonJSModules_GetModuleName(t.GetInput())) == *exportedName &&
					r.outer.GetBaseQualifiedNameNode(n) == n &&
					NodeUtil_IsLValue(n)
			r.UpdateNameReference(t, n, name, *exportedName, true, exportPropIsConst)

			// If it's a global name, rename it to prevent conflicts with other scripts
		} else if var_.IsGlobal() {
			currentModuleName := *ProcessCommonJSModules_GetModuleName(t.GetInput())

			if currentModuleName == name {
				return
			}

			// refs to 'exports' are handled separately.
			if processCommonJSModules_EXPORTS == name {
				return
			}

			// closure_test_suite looks for test*() functions
			if r.outer.compiler.GetOptions().ShouldExportTestFunctions() &&
				strings.HasPrefix(currentModuleName, "test") {
				return
			}

			newName := name + "$$" + currentModuleName
			r.UpdateNameReference(t, n, name, newName, false, false)
		}
	}
}

// UpdateNameReference: nameRef is the qualified name node; originalName of
// nameRef; newName for nameRef; requireFunctionExpressions whether named
// class or functions should be rewritten to variable assignments.
//
// Java: ProcessCommonJSModules.RewriteModule.updateNameReference(NodeTraversal, Node, String, String, boolean, boolean)
func (r *ProcessCommonJSModules_RewriteModule) UpdateNameReference(
	t *NodeTraversal,
	nameRef rhino.Node,
	originalName string,
	newName string,
	requireFunctionExpressions bool,
	qualifiedNameIsConst bool) {
	compiler := r.outer.compiler
	parent := nameRef.GetParent()
	rhino.CheckNotNull(parent)
	// checkNotNull(newName): a Go string cannot be null.
	newNameIsQualified := strings.IndexByte(newName, '.') >= 0
	newNameIsModuleExport :=
		newName == r.outer.GetBasePropertyImport(*ProcessCommonJSModules_GetModuleName(t.GetInput()))

	newNameDeclaration := t.GetScope().GetVar(newName)

	switch parent.GetToken() {
	case rhino.Token_CLASS:
		cls := ast.MustClass(parent)
		className, hasName := cls.Name()
		superClass, hasExtends := cls.Extends()
		if hasName && className.Is(nameRef) &&
			(newNameIsQualified || requireFunctionExpressions) {
			// Refactor a named class to a class expression
			// We can't remove the class name during a traversal, so save it for later
			r.rewrittenClassExpressions = append(r.rewrittenClassExpressions, parent)

			newNameRef := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				newNameRef.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}

			var expr rhino.Node
			if !newNameIsQualified && newNameDeclaration == nil {
				expr = rhino.IR_LetWithValue(newNameRef, rhino.IR_NullNode()).SrcrefTreeIfMissing(nameRef)
			} else {
				expr =
					rhino.IR_ExprResult(rhino.IR_Assign(newNameRef, rhino.IR_NullNode())).SrcrefTreeIfMissing(nameRef)
				info := rhino.JSDocInfo_Builder_MaybeCopyFrom(parent.GetJSDocInfo())
				parent.SetJSDocInfo(nil)
				if qualifiedNameIsConst {
					info.RecordConstancy()
				}
				expr.GetFirstChild().SetJSDocInfo(info.Build())
				r.FixTypeAnnotationsForNode(t, expr.GetFirstChild())
			}
			parent.ReplaceWith(expr)
			if expr.IsLet() {
				expr.GetFirstFirstChild().ReplaceWith(parent)
			} else {
				expr.GetFirstChild().GetSecondChild().ReplaceWith(parent)
			}
		} else if hasExtends && superClass.Is(nameRef) {
			newNameRef := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				newNameRef.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			nameRef.ReplaceWith(newNameRef)
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(jutil.StringPtr(originalName))
		}
	case rhino.Token_FUNCTION:
		if newNameIsQualified || requireFunctionExpressions {
			// Refactor a named function to a function expression
			if NodeUtil_IsFunctionExpression(parent) {
				// Don't refactor if the parent is a named function expression.
				// e.g. var foo = function foo() {};
				return
			}
			newNameRef := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				newNameRef.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			nameRef.SetString("")

			var expr rhino.Node
			if !newNameIsQualified && newNameDeclaration == nil {
				expr = rhino.IR_VarWithValue(newNameRef, rhino.IR_NullNode()).SrcrefTreeIfMissing(nameRef)
			} else {
				expr =
					rhino.IR_ExprResult(rhino.IR_Assign(newNameRef, rhino.IR_NullNode())).SrcrefTreeIfMissing(nameRef)
			}
			parent.ReplaceWith(expr)
			if expr.IsVar() {
				expr.GetFirstFirstChild().ReplaceWith(parent)
			} else {
				expr.GetFirstChild().GetSecondChild().ReplaceWith(parent)
				info := rhino.JSDocInfo_Builder_MaybeCopyFrom(parent.GetJSDocInfo())
				parent.SetJSDocInfo(nil)
				if qualifiedNameIsConst {
					info.RecordConstancy()
				}
				expr.GetFirstChild().SetJSDocInfo(info.Build())
				r.FixTypeAnnotationsForNode(t, expr.GetFirstChild())
			}
			r.functionsToHoist = append(r.functionsToHoist, expr)
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(jutil.StringPtr(originalName))
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		// Multiple declaration - needs split apart.
		if parent.HasMoreThanOneChild() && !NodeUtil_IsAnyFor(parent.GetParent()) {
			r.SplitMultipleDeclarations(parent)
			parent = nameRef.GetParent()
			newNameDeclaration = t.GetScope().GetVar(newName)
		}

		// nameRef is a declared name, with its value if it has one.
		init, hasInit := ast.MustName(nameRef).Init()
		if newNameIsQualified {
			// Var declarations without initialization can simply
			// be removed if they are being converted to a property.
			if !hasInit && parent.GetJSDocInfo() == nil {
				parent.Detach()
				break
			}

			// Refactor a var declaration to a getprop assignment
			getProp := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				getProp.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			info := parent.GetJSDocInfo()
			parent.SetJSDocInfo(nil)
			if hasInit {
				assign := rhino.IR_Assign(getProp, init.Detach().Node())
				assign.SetJSDocInfo(info)
				expr := rhino.IR_ExprResult(assign).SrcrefTreeIfMissing(nameRef)
				parent.ReplaceWith(expr)
				infoBuilder := rhino.JSDocInfo_Builder_MaybeCopyFrom(info)
				parent.SetJSDocInfo(nil)
				if qualifiedNameIsConst {
					infoBuilder.RecordConstancy()
				}
				assign.SetJSDocInfo(infoBuilder.Build())
				r.FixTypeAnnotationsForNode(t, assign)
			} else {
				getProp.SetJSDocInfo(info)
				parent.ReplaceWith(rhino.IR_ExprResult(getProp).Srcref(getProp))
			}
		} else if newNameDeclaration != nil && newNameDeclaration.GetNameNode() != nameRef {
			// Variable is already defined. Convert this to an assignment.
			// If the variable declaration has no initialization, we simply
			// remove the node. This can occur when the variable which is exported
			// is declared in an outer scope but assigned in an inner one.
			if !hasInit {
				parent.Detach()
				break
			}

			name := NodeUtil_NewNameWithOriginalName(compiler, newName, nameRef, originalName)
			assign := rhino.IR_Assign(name, init.Detach().Node())
			info := parent.GetJSDocInfo()
			if info != nil {
				parent.SetJSDocInfo(nil)
				assign.SetJSDocInfo(info)
			}

			parent.ReplaceWith(rhino.IR_ExprResult(assign).SrcrefTree(nameRef))
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(jutil.StringPtr(originalName))
		}
	default:
		// Whenever possible, reuse the existing reference
		if !newNameIsQualified && nameRef.IsName() {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(jutil.StringPtr(originalName))
		} else {
			var name rhino.Node
			if newNameIsQualified {
				name = NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			} else {
				name = NodeUtil_NewNameWithOriginalName(compiler, newName, nameRef, originalName)
			}

			if newNameIsModuleExport {
				name.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			info := nameRef.GetJSDocInfo()
			if info != nil {
				nameRef.SetJSDocInfo(nil)
				name.SetJSDocInfo(info)
			}
			name.SrcrefTree(nameRef)
			nameRef.ReplaceWith(name)
			if nameRef.HasChildren() {
				name.AddChildrenToFront(nameRef.RemoveChildren())
			}
		}
	}

	t.ReportCodeChange()
}

// GetExportedName determines whether the given name Node n is referenced in
// an export. Returns: if the name is not used in an export, its own name; if
// the name node is actually the export target itself, nil.
//
// Java: ProcessCommonJSModules.RewriteModule.getExportedName(NodeTraversal, Node, Var)
func (r *ProcessCommonJSModules_RewriteModule) GetExportedName(t *NodeTraversal, n rhino.Node, var_ *Var) *string {
	if var_ == nil || !jutil.ObjectsEquals(var_.GetNode().GetInputId(), n.GetInputId()) {
		return n.GetQualifiedName()
	}

	baseExportName := r.outer.GetBasePropertyImport(*ProcessCommonJSModules_GetModuleName(t.GetInput()))

	for _, export := range r.exports {
		exportBase := r.outer.GetBaseQualifiedNameNode(export.node)
		exportRValue := NodeUtil_GetRValueOfLValue(exportBase)

		if exportRValue.IsNil() {
			continue
		}

		exportedName := r.GetExportedNameNode(export)
		// We don't want to handle the export itself
		if export.isInSupportedScope &&
			(exportRValue == n ||
				((NodeUtil_IsClassExpression(exportRValue) ||
					NodeUtil_IsFunctionExpression(exportRValue)) &&
					exportedName == n)) {
			return nil
		}

		exportBaseQName := exportBase.GetQualifiedName()

		if exportRValue.IsObjectLit() {
			if exportBaseQName == nil || "module.exports" != *exportBaseQName {
				return n.GetQualifiedName()
			}

			for member := range ast.MustObjectLit(exportRValue).Members() {
				key, ok := member.AsStringKey()
				if !ok || key.IsQuoted() ||
					!NodeUtil_IsValidPropertyName(
						r.outer.compiler.GetOptions().GetLanguageIn().ToFeatureSet(), key.Key()) {
					continue
				}
				value := key.Value().Node()
				if !value.IsQualifiedName() {
					continue
				}
				if value == n {
					return nil
				}

				valVar := t.GetScope().GetVar(*value.GetQualifiedName())
				if valVar != nil && valVar.GetNameNode() == var_.GetNameNode() {
					// The key exports the name.
					if export.isInSupportedScope {
						return jutil.StringPtr(baseExportName + "." + key.Key())
					} else {
						return n.GetQualifiedName()
					}
				}
			}
		} else {
			if !export.isInSupportedScope {
				return n.GetQualifiedName()
			}
			if var_.GetNameNode() == exportedName {
				var exportPrefix string
				if strings.HasPrefix(*exportBaseQName, processCommonJSModules_MODULE) {
					exportPrefix = processCommonJSModules_MODULE + "." + processCommonJSModules_EXPORTS
				} else {
					exportPrefix = processCommonJSModules_EXPORTS
				}

				if jutil.StringLength(*exportBaseQName) == jutil.StringLength(exportPrefix) {
					return jutil.StringPtr(baseExportName)
				}

				return jutil.StringPtr(baseExportName + jutil.StringSubstringFrom(*exportBaseQName, jutil.StringLength(exportPrefix)))
			}
		}
	}
	return n.GetQualifiedName()
}

// Java: ProcessCommonJSModules.RewriteModule.getExportedNameNode(ExportInfo)
func (r *ProcessCommonJSModules_RewriteModule) GetExportedNameNode(info *ProcessCommonJSModules_ExportInfo) rhino.Node {
	qNameBase := r.outer.GetBaseQualifiedNameNode(info.node)
	rValue := NodeUtil_GetRValueOfLValue(qNameBase)

	if rValue.IsNil() {
		return rhino.NilNode
	}

	if NodeUtil_IsFunctionExpression(rValue) {
		return rValue.GetFirstChild()
	}

	if NodeUtil_IsClassExpression(rValue) && rValue.GetFirstChild() == qNameBase {
		return rValue.GetFirstChild()
	}

	if !rValue.IsName() {
		return rhino.NilNode
	}

	var_ := info.scope.GetVar(rValue.GetString())
	if var_ == nil {
		return rhino.NilNode
	}

	return var_.GetNameNode()
}

// GetModuleImportName determines if the given Node n is an alias created by
// a module import. Returns nil if it's not an alias or the imported module
// name.
//
// Java: ProcessCommonJSModules.RewriteModule.getModuleImportName(NodeTraversal, Node)
func (r *ProcessCommonJSModules_RewriteModule) GetModuleImportName(t *NodeTraversal, n rhino.Node) *string {
	var rValue rhino.Node
	propSuffix := ""
	parent := n.GetParent()
	// `{foo} = value` in `const {foo} = value`: the value is the destructuring's init.
	if key, ok := ast.AsStringKey(parent); ok {
		if lhs, ok := ast.AsDestructuringLhs(parent.GetGrandparent()); ok && lhs.Pattern().Is(parent.GetParent()) {
			if _, isObjectPattern := lhs.Pattern().AsObjectPattern(); isObjectPattern {
				init, _ := lhs.Init() // the nil node if there is none
				rValue = init.Node()
				propSuffix = "." + key.Key()
			}
		}
	}
	if propSuffix == "" && !parent.IsNil() {
		rValue = NodeUtil_GetRValueOfLValue(n)
	}

	if rValue.IsNil() {
		return nil
	}

	if rValue.IsCall() && r.outer.IsCommonJsImport(rValue) {
		return jutil.StringPtr(r.outer.GetBasePropertyImportWithRequireCall(r.outer.GetImportedModuleName(t, rValue), rValue) + propSuffix)
	} else if getProp, ok := ast.AsGetProp(rValue); ok && r.outer.IsCommonJsImport(getProp.Object().Node()) {
		// var foo = require('bar').foo;
		require := getProp.Object().Node()
		importName :=
			r.outer.GetBasePropertyImportWithRequireCall(r.outer.GetImportedModuleName(t, require), require)

		suffix := getProp.Property()

		return jutil.StringPtr(importName + "." + suffix + propSuffix)
	}

	return nil
}

// FixTypeNode updates any type references in JSDoc annotations to account
// for all the rewriting we've done.
//
// Java: ProcessCommonJSModules.RewriteModule.fixTypeNode(NodeTraversal, Node)
func (r *ProcessCommonJSModules_RewriteModule) FixTypeNode(t *NodeTraversal, typeNode rhino.Node) {
	if typeNode.IsStringLit() {
		name := typeNode.GetString()
		// Type nodes can be module paths.
		if ModuleLoader_IsPathIdentifier(name) {
			lastSlash := jutil.StringLastIndexOf(name, "/")
			endIndex := jutil.StringIndexOfCharFrom(name, '.', lastSlash)
			var localTypeName *string
			if endIndex == -1 {
				endIndex = jutil.StringLength(name)
			} else {
				localTypeName = jutil.StringPtr(jutil.StringSubstringFrom(name, endIndex))
			}

			moduleName := jutil.StringSubstring(name, 0, endIndex)
			globalModuleName := r.outer.GetImportedModuleNameWithImportPath(t, typeNode, moduleName)
			baseImportProperty := r.outer.GetBasePropertyImport(globalModuleName)
			if localTypeName == nil {
				typeNode.SetString(baseImportProperty)
			} else {
				typeNode.SetString(baseImportProperty + *localTypeName)
			}

		} else {
			// A type node can be a getprop. Any portion of the getprop
			// can be either an import alias or export alias. Check each
			// segment.
			wasRewritten := false
			endIndex := -1
			for endIndex < jutil.StringLength(name) {
				endIndex = jutil.StringIndexOfCharFrom(name, '.', endIndex+1)
				if endIndex == -1 {
					endIndex = jutil.StringLength(name)
				}
				baseName := jutil.StringSubstring(name, 0, endIndex)
				suffix := ""
				if endIndex < jutil.StringLength(name) {
					suffix = jutil.StringSubstringFrom(name, endIndex)
				}
				typeDeclaration := t.GetScope().GetVar(baseName)

				// Make sure we can find a variable declaration (and it's in this file)
				if typeDeclaration != nil && !typeDeclaration.GetNode().IsNil() &&
					jutil.ObjectsEquals(typeDeclaration.GetNode().GetInputId(), typeNode.GetInputId()) {
					importedModuleName := r.GetModuleImportName(t, typeDeclaration.GetNode())

					// If the name is an import alias, rewrite it to be a reference to the
					// module name directly
					if importedModuleName != nil {
						typeNode.SetString(*importedModuleName + suffix)
						typeNode.SetOriginalName(jutil.StringPtr(name))
						wasRewritten = true
						break
					} else if r.allowFullRewrite {
						// Names referenced in export statements can only be rewritten in
						// commonjs modules.
						exportedName := r.GetExportedName(t, typeNode, typeDeclaration)
						if exportedName != nil && *exportedName != name {
							typeNode.SetString(*exportedName + suffix)
							typeNode.SetOriginalName(jutil.StringPtr(name))
							wasRewritten = true
							break
						}
					}
				}
			}

			// If the name was neither an import alias or referenced in an export,
			// We still may need to rename it if it's global
			if !wasRewritten && r.allowFullRewrite {
				endIndex = jutil.StringIndexOfChar(name, '.')
				if endIndex == -1 {
					endIndex = jutil.StringLength(name)
				}
				baseName := jutil.StringSubstring(name, 0, endIndex)
				typeDeclaration := t.GetScope().GetVar(baseName)
				if typeDeclaration != nil && typeDeclaration.IsGlobal() {
					moduleName := *ProcessCommonJSModules_GetModuleName(t.GetInput())
					newName := baseName + "$$" + moduleName
					if endIndex < jutil.StringLength(name) {
						newName += jutil.StringSubstringFrom(name, endIndex)
					}

					typeNode.SetString(newName)
					typeNode.SetOriginalName(jutil.StringPtr(name))
				}
			}
		}
	}

	for child := typeNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		r.FixTypeNode(t, child)
	}
}

// Java: ProcessCommonJSModules.RewriteModule.splitMultipleDeclarations(Node)
func (r *ProcessCommonJSModules_RewriteModule) SplitMultipleDeclarations(var_ rhino.Node) []rhino.Node {
	declaration, ok := ast.AsNameDeclaration(var_)
	jutil.CheckState(ok)
	vars := []rhino.Node{}
	info := var_.GetJSDocInfo()
	// Move each declarator but the last into a declaration of its own.
	for first, _ := declaration.FirstDeclarator(); !first.IsLast(); first, _ = declaration.FirstDeclarator() {
		newVar := rhino.NewNodeWithChild(var_.GetToken(), first.Detach().Node())

		if info != nil {
			newVar.SetJSDocInfo(info.Clone())
		}

		newVar.Srcref(var_)
		newVar.InsertBefore(var_)
		vars = append(vars, newVar)
	}
	vars = append(vars, var_)
	return vars
}
