// Port of com/google/javascript/jscomp/ijs/ConvertToTypedInterface.java @ closure-compiler v20261028
// Port of com/google/javascript/jscomp/ijs/FileInfo.java @ closure-compiler v20261028
// Port of com/google/javascript/jscomp/ijs/ClassUtil.java @ closure-compiler v20261028

/*
 * Copyright 2016 The Closure Compiler Authors.
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

/*
 * Copyright 2017 The Closure Compiler Authors.
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

// ConvertToTypedInterface shrinks an AST into a `.i.js` type summary,
// preserving typing rather than runtime behavior by stripping function/method
// bodies, non-declaration statements, and unneeded initializer expressions.
//
// This is conceptually similar to the ijar tool[1] that bazel uses to shrink
// jars into minimal versions that can be used equivalently for compilation
// of downstream dependencies.
//
// [1] https://github.com/bazelbuild/bazel/blob/master/third_party/ijar/README.txt
type ConvertToTypedInterface struct {
	compiler AbstractCompiler
}

var ConvertToTypedInterface_CONSTANT_WITH_SUGGESTED_TYPE = DiagnosticType_Warning(
	"JSC_CONSTANT_WITH_SUGGESTED_TYPE",
	"Constants in top-level should have types explicitly specified.\n"+
		"You may want specify this type as:\t@const '{'{0}'}'")

var ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE = DiagnosticType_Warning(
	"JSC_CONSTANT_WITHOUT_EXPLICIT_TYPE",
	"Constants in top-level should have types explicitly specified.")

var ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE = DiagnosticType_Warning(
	"JSC_GOOG_SCOPE_HIDDEN_TYPE",
	"Found a goog.scope local type declaration.\n"+
		"If you can't yet migrate this file to goog.module, use a /** @private */"+
		" declaration within the goog.provide namespace.")

var (
	IjsErrors_BAD_IJS_FILE_NAME = DiagnosticType_Error(
		"JSC_BAD_IJS_FILE_NAME", "--ijs file name must end in .i.js: {0}")

	IjsErrors_CONFLICTING_IJS_FILE = DiagnosticType_Warning(
		"JSC_CONFLICTING_IJS_FILE",
		"Found --ijs file in a conflict with --js or --weakdep file. This usually means the same"+
			" file appears both in the srcs and the transitive deps of a build target.\n"+
			"File 1: {0}\n"+
			"File 2: {1}")
)

var (
	ijsPolymerQName                    = rhino.QualifiedName_Of("Polymer")
	ijsGoogAddSingletonGetterQName     = rhino.QualifiedName_Of("goog.addSingletonGetter")
	ijsGoogModuleQName                 = rhino.QualifiedName_Of("goog.module")
	ijsGoogDeclareLegacyNamespaceQName = rhino.QualifiedName_Of("goog.module.declareLegacyNamespace")
	ijsGoogDeclareModuleIdQName        = rhino.QualifiedName_Of("goog.declareModuleId")
)

func isCalleeNodeToPreserve(callee rhino.Node) bool {
	if callee.IsName() {
		return ijsPolymerQName.Matches(callee)
	}
	if !callee.IsGetProp() {
		return false
	}
	return ijsGoogModuleQName.Matches(callee) ||
		ijsGoogProvideQName.Matches(callee) ||
		ijsGoogRequireQName.Matches(callee) ||
		ijsGoogRequireTypeQName.Matches(callee) ||
		ijsGoogForwardDeclareQName.Matches(callee) ||
		ijsGoogDeclareModuleIdQName.Matches(callee) ||
		ijsGoogDeclareLegacyNamespaceQName.Matches(callee) ||
		ijsGoogDefineQName.Matches(callee) ||
		ijsGoogAddSingletonGetterQName.Matches(callee)
}

func NewConvertToTypedInterface(compiler AbstractCompiler) *ConvertToTypedInterface {
	return &ConvertToTypedInterface{compiler: compiler}
}

func maybeReportIjsDiagnostic(
	compiler AbstractCompiler, node rhino.Node, diagnostic *DiagnosticType, fillers ...string) {
	if sourceName, ok := NodeUtil_SourceName(node); ok &&
		(strings.HasSuffix(sourceName, "_test.js") ||
			strings.HasSuffix(sourceName, "_test.closure.js") ||
			strings.HasSuffix(sourceName, "_test.tsx.cl.js")) {
		// Allow _test.js files and their tsickle generated
		// equivalents to avoid emitting errors at .i.js generation time.
		// We expect these files to not be consumed by any other downstream libraries.
		return
	}
	compiler.Report(JSError_MakeWithNode(node, diagnostic, fillers...))
}

func maybeWarnForConstWithoutExplicitType(compiler AbstractCompiler, decl *potentialDeclaration) {
	if decl.isConstToBeInferred() &&
		!decl.lhs.IsFromExterns() &&
		!isPrivateJSDoc(decl.jsdoc()) {

		nameNode := decl.lhs
		if nameNode.GetJSType() == nil {
			maybeReportIjsDiagnostic(compiler, nameNode, ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE)
		} else {
			maybeReportIjsDiagnostic(
				compiler,
				nameNode,
				ConvertToTypedInterface_CONSTANT_WITH_SUGGESTED_TYPE,
				nameNode.GetJSType().ToAnnotationString(rhino.JSType_Nullability_EXPLICIT))
		}
	}
}

func (c *ConvertToTypedInterface) Process(externs rhino.Node, root rhino.Node) {
	rootView, ok := ast.AsRoot(root)
	if !ok {
		return
	}
	for item := range ast.ReadAhead(rootView.FirstItem()) {
		if script, isScript := item.AsScript(); isScript {
			sourceFileName, _ := script.Node().SourceFileName()
			if AbstractCompiler_IsFillFileName(sourceFileName) {
				script.Detach()
			}
		}
	}
	if c.compiler != nil && traverseScriptsParallel(c.compiler, root, parallelScriptOptions[struct{}]{
		enableTypeRegistryReads: true,
		initWorker:              func(int, *Compiler) struct{} { return struct{}{} },
		runScript: func(_ struct{}, _ int, _ rhino.Node, scriptNode rhino.Node, rc *Compiler) {
			processIjsScript(rc, ast.MustScript(scriptNode))
		},
	}) {
		return
	}
	for item := range ast.ReadAhead(rootView.FirstItem()) {
		if script, isScript := item.AsScript(); isScript {
			processIjsScript(c.compiler, script)
		}
	}
}

func processIjsScript(compiler AbstractCompiler, script ast.Script) {
	scriptNode := script.Node()
	sourceFileName, _ := scriptNode.SourceFileName()

	if scriptJsDoc := scriptNode.GetJSDocInfo(); scriptJsDoc != nil && scriptJsDoc.IsClosureUnawareCode() {
		// If we are generating type summary files, then those files won't contain the method content
		// from within the closure-unaware section, and the entire file effectively becomes
		// closure-aware again (as it is generated code that describes a module shape).
		scriptJsDocBuilder := scriptJsDoc.ToBuilder()
		scriptJsDocBuilder.RemoveClosureUnawareCode()
		scriptNode.SetJSDocInfo(scriptJsDocBuilder.Build())
	}

	currentFile := newIjsFileInfo(sourceFileName)
	NodeTraversal_Traverse(compiler, scriptNode, ijsRemoveNonDeclarations{})
	NodeTraversal_Traverse(compiler, scriptNode, &ijsPropagateConstJsdoc{currentFile: currentFile})
	(&ijsSimplifyDeclarations{compiler: compiler, currentFile: currentFile}).simplifyAll()
}

func findNameDeclaration(scope *Scope, rhs rhino.Node) *Var {
	if scope == nil || !rhs.IsName() {
		return nil
	}
	return scope.GetVarByID(rhs.GetStringID())
}

type ijsRemoveNonDeclarations struct{}

var ijsRemoveVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_TRY)
	s.Add(rhino.Token_LABEL)
	s.Add(rhino.Token_DEFAULT_CASE)
	s.Add(rhino.Token_CASE)
	s.Add(rhino.Token_DO)
	s.Add(rhino.Token_WHILE)
	s.Add(rhino.Token_FOR)
	s.Add(rhino.Token_FOR_IN)
	s.Add(rhino.Token_FOR_OF)
	s.Add(rhino.Token_FOR_AWAIT_OF)
	s.Add(rhino.Token_IF)
	s.Add(rhino.Token_SWITCH)
	s.Add(rhino.Token_VAR)
	s.Add(rhino.Token_LET)
	s.Add(rhino.Token_CONST)
	s.Add(rhino.Token_BLOCK)
	return s
}()

func (ijsRemoveNonDeclarations) visitTokens() *ast.TokenSet {
	return &ijsRemoveVisitTokens
}

var ijsGoogScopeQName = rhino.QualifiedName_Of("goog.scope")

func preserveOrDelete(n rhino.Node, t *NodeTraversal, keep bool) bool {
	if !keep {
		NodeUtil_DeleteNode(n, t.GetCompiler())
	}
	return keep
}

func (ijsRemoveNonDeclarations) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		if !isIjsConstructor(n) || !hasNamedClass(n) {
			fnBody := ast.MustFunction(n).Body()
			if block, ok := fnBody.AsBlock(); !ok || block.HasStatements() {
				body := fnBody.Node()
				t.ReportCodeChangeWithNode(body)
				body.ReplaceWith(ast.NewBlock().Node().Srcref(body))
				NodeUtil_MarkFunctionsDeleted(body, t.GetCompiler())
			}
		}
		return true
	case rhino.Token_MEMBER_FIELD_DEF:
		// Remove fields in anonymous classes.
		return preserveOrDelete(n, t, isMemberFieldInsideNamedClass(n))
	case rhino.Token_EXPR_RESULT:
		expr := ast.MustExprStmt(n).Expr().Node()
		switch expr.GetToken() {
		case rhino.Token_CALL:
			callee := ast.MustCall(expr).Callee().Node()
			if ijsGoogScopeQName.Matches(callee) {
				panic(jutil.NewIllegalStateException())
			}
			return preserveOrDelete(n, t, isCalleeNodeToPreserve(callee))
		case rhino.Token_ASSIGN:
			return preserveOrDelete(n, t, shouldPreserveAssignment(expr, t))
		case rhino.Token_GETPROP:
			return preserveOrDelete(n, t, expr.IsQualifiedName() && expr.GetJSDocInfo() != nil)
		case rhino.Token_GETELEM:
			elem := ast.MustGetElem(expr)
			return preserveOrDelete(
				n,
				t,
				isSymbolProp(elem.Index().Node()) &&
					elem.Object().Node().IsQualifiedName() &&
					expr.GetJSDocInfo() != nil)
		default:
			NodeUtil_DeleteNode(n, t.GetCompiler())
			return false
		}
	case rhino.Token_COMPUTED_PROP:
		if isComputedMemberInsideClass(n) {
			return true
		}
		if val, ok := ast.MustComputedProp(n).Value(); !ok || !NodeUtil_IsLhsByDestructuring(val.Node()) {
			NodeUtil_DeleteNode(n, t.GetCompiler())
		}
		return false
	case rhino.Token_COMPUTED_FIELD_DEF:
		return preserveOrDelete(n, t, isComputedMemberInsideClass(n))
	case rhino.Token_THROW, rhino.Token_RETURN, rhino.Token_BREAK, rhino.Token_CONTINUE, rhino.Token_DEBUGGER, rhino.Token_EMPTY:
		if !parent.IsNil() && NodeUtil_IsStatementParent(parent) {
			NodeUtil_DeleteNode(n, t.GetCompiler())
		}
		return false
	case rhino.Token_LABEL:
		NodeUtil_DeleteNode(ast.MustLabel(n).Name().Node(), t.GetCompiler())
		return true
	case rhino.Token_IF:
		NodeUtil_DeleteNode(ast.MustIf(n).Cond().Node(), t.GetCompiler())
		return true
	case rhino.Token_SWITCH:
		NodeUtil_DeleteNode(ast.MustSwitch(n).Discriminant().Node(), t.GetCompiler())
		return true
	case rhino.Token_CASE:
		NodeUtil_DeleteNode(ast.MustCase(n).Value().Node(), t.GetCompiler())
		return true
	case rhino.Token_WHILE:
		NodeUtil_DeleteNode(ast.MustWhile(n).Cond().Node(), t.GetCompiler())
		return true
	case rhino.Token_TRY:
		NodeUtil_DeleteNode(ast.MustTry(n).Catches().Node(), t.GetCompiler())
		return true
	case rhino.Token_DO:
		NodeUtil_DeleteNode(ast.MustDo(n).Cond().Node(), t.GetCompiler())
		return true
	case rhino.Token_FOR:
		var children [4]rhino.Node
		i := 0
		for child := range n.Children() {
			if i < 4 {
				children[i] = child
			}
			i++
		}
		if i == 4 {
			compiler := t.GetCompiler()
			initNode, condNode, incrNode, bodyNode := children[0], children[1], children[2], children[3]
			NodeUtil_DeleteNode(condNode, compiler)
			NodeUtil_DeleteNode(incrNode, compiler)
			initNode.Detach()
			if initNode.IsVar() {
				bodyNode.AddChildToFront(initNode)
			} else {
				NodeUtil_MarkFunctionsDeleted(initNode, compiler)
			}
		}
		return true
	case rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_FOR_IN:
		loop := ast.MustEnhancedFor(n)
		target, collection, body := loop.Target(), loop.Collection(), loop.Body()
		NodeUtil_DeleteNode(collection.Node(), t.GetCompiler())
		if initNode := target.Detach().Node(); initNode.IsVar() {
			body.Node().AddChildToFront(initNode)
		} else {
			NodeUtil_MarkFunctionsDeleted(initNode, t.GetCompiler())
		}
		return true
	case rhino.Token_CONST, rhino.Token_LET:
		if !t.InGlobalScope() && !t.InModuleScope() {
			NodeUtil_RemoveChild(parent, n)
			NodeUtil_MarkFunctionsDeleted(n, t.GetCompiler())
			t.ReportCodeChangeWithNode(parent)
			return false
		}
		return true
	case rhino.Token_VAR:
		if !t.InGlobalHoistScope() && !t.InModuleHoistScope() {
			NodeUtil_RemoveChild(parent, n)
			NodeUtil_MarkFunctionsDeleted(n, t.GetCompiler())
			t.ReportCodeChangeWithNode(parent)
			return false
		}
		return true
	case rhino.Token_BLOCK:
		if !parent.IsNil() && parent.IsClassMembers() {
			// ES2022 `static { ... }` class initialization blocks hold only runtime code.
			NodeUtil_DeleteNode(n, t.GetCompiler())
			return false
		}
		return true
	case rhino.Token_MODULE_BODY, rhino.Token_CLASS, rhino.Token_DEFAULT_CASE, rhino.Token_EXPORT, rhino.Token_IMPORT:
		return true
	default:
		if NodeUtil_IsStatement(n) {
			panic(jutil.NewIllegalStateExceptionMsg(n.GetToken().String()))
		}
		return true
	}
}

func (ijsRemoveNonDeclarations) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_TRY, rhino.Token_LABEL, rhino.Token_DEFAULT_CASE, rhino.Token_CASE, rhino.Token_DO, rhino.Token_WHILE,
		rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_IF:
		if n.HasParent() {
			children := n.RemoveChildren()
			parent.AddChildrenAfter(children, n)
			// ShouldTraverse intentionally leaves control-flow containers in a temporary
			// non-standard layout before splicing hoisted VAR declarations into parent.
			n.Detach()
			t.ReportCodeChange()
		}
	case rhino.Token_SWITCH:
		// ShouldTraverse already deleted the switch discriminant, leaving only SWITCH_BODY.
		if n.HasParent() {
			if !n.HasOneChild() {
				panic(jutil.NewIllegalStateExceptionMsg("malformed SWITCH: " + n.ToStringTree()))
			}
			var bodyNode rhino.Node
			for child := range n.Children() {
				bodyNode = child
			}
			children := ast.MustSwitchBody(bodyNode).Node().RemoveChildren()
			parent.AddChildrenAfter(children, n)
			n.Detach()
			t.ReportCodeChange()
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		splitNameDeclarations(ast.MustNameDeclaration(n), t)
	case rhino.Token_BLOCK:
		if !parent.IsNil() && parent.IsClassMembers() {
			NodeUtil_DeleteNode(n, t.GetCompiler())
		} else if !parent.IsNil() && !parent.IsFunction() {
			parent.AddChildrenAfter(n.RemoveChildren(), n)
			n.Detach()
			t.ReportCodeChangeWithNode(parent)
		}
	}
}

// splitNameDeclarations splits multi-declarator var/let/const statements into
// separate single-declarator statements, merges inline JSDoc onto the
// declaration node, and removes empty non-import destructuring patterns.
func splitNameDeclarations(decl ast.NameDeclaration, t *NodeTraversal) {
	n := decl.Node()
	sharedJsdoc := n.GetJSDocInfo()
	isExport := ast.RoleOf(n) == ast.RoleExportBody
	statement := n
	if isExport {
		statement = n.GetParent()
	}
	for decl.HasDeclarators() {
		d, _ := decl.LastDeclarator()
		lhsToSplit := d.Node()
		nameJsdoc := lhsToSplit.GetJSDocInfo()
		lhsToSplit.SetJSDocInfo(nil)
		mergedJsdoc := mergeIjsJSDocs(sharedJsdoc, nameJsdoc)
		if decl.NumDeclarators() == 1 {
			n.SetJSDocInfo(mergedJsdoc)
			deleteIfEmptyNonImportDestructuring(d, statement, t)
			return
		}
		// A name declaration with more than one declarator is split into separate declarations.
		nameDecl, ok := d.AsName()
		if !ok {
			d.Node().Detach()
			panic(jutil.NewIllegalStateException())
		}
		var rhs rhino.Node
		if init, hasInit := nameDecl.Init(); hasInit {
			rhs = init.Detach().Node()
		}
		newDeclNode := NodeUtil_NewDeclaration(nameDecl.Detach().Node(), rhs, decl.Token()).Srcref(n)
		newDeclNode.SetJSDocInfo(mergedJsdoc)
		newStmt := newDeclNode
		if isExport {
			newStmt = ast.NewExport(ast.MustExportBody(newDeclNode)).Node().Srcref(statement)
		}
		newStmt.InsertAfter(statement)
		t.ReportCodeChange()
	}
}

func deleteIfEmptyNonImportDestructuring(d ast.Declarator, stmt rhino.Node, t *NodeTraversal) {
	dlhs, ok := d.AsDestructuringLhs()
	if !ok {
		return
	}
	hasBinding := false
	NodeUtil_VisitLhsNodesInNode(dlhs.Node(), func(rhino.Node) { hasBinding = true })
	if hasBinding {
		return
	}
	init, _ := dlhs.Init()
	if !isImportRhs(init.Node()) {
		NodeUtil_DeleteNode(stmt, t.GetCompiler())
	}
}

type ijsPropagateConstJsdoc struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	currentFile *ijsFileInfo
}

var ijsPropagateVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_CLASS)
	s.Add(rhino.Token_MEMBER_FIELD_DEF)
	s.Add(rhino.Token_FUNCTION)
	s.Add(rhino.Token_EXPR_RESULT)
	s.Add(rhino.Token_VAR)
	s.Add(rhino.Token_CONST)
	s.Add(rhino.Token_LET)
	s.Add(rhino.Token_STRING_KEY)
	return s
}()

func (*ijsPropagateConstJsdoc) visitTokens() *ast.TokenSet {
	return &ijsPropagateVisitTokens
}

var (
	ijsGoogDefineQName  = rhino.QualifiedName_Of("goog.define")
	ijsGoogProvideQName = rhino.QualifiedName_Of("goog.provide")
	ijsCjsRequireQName  = rhino.QualifiedName_Of("require")
)

func (p *ijsPropagateConstJsdoc) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		if !parent.IsNil() && NodeUtil_IsStatementParent(parent) {
			if name, ok := ast.MustClass(n).Name(); ok && name.Ident() != "" {
				p.currentFile.recordNameDeclaration(name.Node())
			}
		}
	case rhino.Token_MEMBER_FIELD_DEF:
		field := ast.MustMemberFieldDef(n)
		if _, hasVal := field.Value(); hasVal {
			p.processDeclarationWithRhs(t, n)
		}
		p.currentFile.recordMemberFieldDef(field)
	case rhino.Token_FUNCTION:
		if !parent.IsNil() && NodeUtil_IsStatementParent(parent) {
			if name, ok := ast.MustFunction(n).Name(); ok && name.Ident() != "" {
				p.currentFile.recordNameDeclaration(name.Node())
			}
		} else if isClassMethod(n) && hasNamedClass(n) {
			p.currentFile.recordMethod(ast.MustMemberFunctionDef(parent))
		}
	case rhino.Token_EXPR_RESULT:
		expr := ast.MustExprStmt(n).Expr().Node()
		switch expr.GetToken() {
		case rhino.Token_CALL:
			call := ast.MustCall(expr)
			callee := call.Callee().Node()
			if ijsGoogProvideQName.Matches(callee) {
				if lastArg, ok := call.LastArg(); ok {
					if str, ok := lastArg.AsStringLit(); ok {
						p.currentFile.markProvided(str.Value())
					}
				}
			} else if ijsGoogRequireQName.Matches(callee) || ijsCjsRequireQName.Matches(callee) {
				if lastArg, ok := call.LastArg(); ok {
					if str, ok := lastArg.AsStringLit(); ok {
						p.currentFile.recordImport(str.Value())
					}
				}
			} else if ijsGoogDefineQName.Matches(callee) {
				p.currentFile.recordDefine(call)
			}
		case rhino.Token_ASSIGN:
			assign := ast.MustAssign(expr)
			p.recordDeclaration(t, assign.Target().Node(), assign.Value().Node())
		case rhino.Token_GETPROP:
			p.currentFile.recordNameDeclaration(expr)
		case rhino.Token_GETELEM:
		default:
			panic(jutil.NewIllegalArgumentExceptionMsg("Unexpected declaration: " + expr.String()))
		}
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		d, ok := ast.MustNameDeclaration(n).OnlyDeclarator()
		if !ok {
			panic(jutil.NewIllegalStateExceptionMsg("expected single declarator"))
		}
		var rhs rhino.Node
		if name, isName := d.AsName(); isName {
			init, _ := name.Init()
			rhs = init.Node()
		} else if dlhs, isDlhs := d.AsDestructuringLhs(); isDlhs {
			init, _ := dlhs.Init()
			rhs = init.Node()
		}
		p.recordDeclaration(t, d.Node(), rhs)
	case rhino.Token_STRING_KEY:
		if !parent.IsNil() && parent.IsObjectLit() {
			sk := ast.MustStringKey(n)
			if !sk.Value().Node().IsNil() {
				p.processDeclarationWithRhs(t, n)
				p.currentFile.recordStringKeyDeclaration(sk)
			}
		}
	}
}

func (p *ijsPropagateConstJsdoc) recordDeclaration(t *NodeTraversal, lhs, rhs rhino.Node) {
	if call, ok := ast.AsCall(rhs); ok &&
		ijsGoogDefineQName.Matches(call.Callee().Node()) &&
		lhs.IsQualifiedName() {
		p.currentFile.recordDefine(call)
	} else if !lhs.IsGetElem() {
		p.recordNameDeclaration(lhs, rhs)
		if !lhs.IsDestructuringLhs() && !rhs.IsNil() {
			p.processDeclarationWithRhs(t, lhs)
		}
	}
}

func (p *ijsPropagateConstJsdoc) recordNameDeclaration(lhs, rhs rhino.Node) {
	parent := lhs.GetParent()
	if parent.IsNil() || (!NodeUtil_IsNameDeclaration(parent) && !parent.IsAssign()) {
		panic(jutil.NewIllegalArgumentExceptionMsg("unexpected declaration parent"))
	}
	isImport := isImportRhs(rhs)
	isAlias := isAliasDeclaration(lhs, rhs)
	NodeUtil_VisitLhsNodesInNode(
		parent,
		func(name rhino.Node) {
			if !name.IsQualifiedName() {
				return
			}
			if isAlias || isImport {
				p.currentFile.recordAliasDeclaration(name)
			} else {
				p.currentFile.recordNameDeclaration(name)
			}
		})
}

func (p *ijsPropagateConstJsdoc) processDeclarationWithRhs(t *NodeTraversal, lhs rhino.Node) {
	if !lhs.IsQualifiedName() && !lhs.IsStringKey() && !lhs.IsMemberFieldDef() {
		panic(jutil.NewIllegalArgumentExceptionMsg("unexpected declaration lhs: " + lhs.String()))
	}
	if NodeUtil_GetRValueOfLValue(lhs).IsNil() {
		panic(jutil.NewIllegalStateExceptionMsg("missing rhs for declaration: " + lhs.String()))
	}
	if isConstNodeToBeInferred(lhs) {
		p.processConstWithRhs(t, lhs)
	}
}

func (p *ijsPropagateConstJsdoc) processConstWithRhs(t *NodeTraversal, nameNode rhino.Node) {
	jsdocNode := NodeUtil_GetBestJSDocInfoNode(nameNode)
	if jsdocNode.IsNil() {
		return
	}
	originalJsdoc := jsdocNode.GetJSDocInfo()
	rhs := NodeUtil_GetRValueOfLValue(nameNode)
	newJsdoc := jsdocForRhs(rhs, originalJsdoc)
	if newJsdoc == nil && rhs.IsName() && isThisPropInsideNamedClass(nameNode) {
		decl := findNameDeclaration(t.GetScope(), rhs)
		newJsdoc = jsdocForName(decl, originalJsdoc)
	}
	if newJsdoc != nil {
		jsdocNode.SetJSDocInfo(newJsdoc)
		t.ReportCodeChange()
	}
}

type ijsSimplifyDeclarations struct {
	compiler    AbstractCompiler
	currentFile *ijsFileInfo
}

// ijsTypingLevel orders declarations by JSDoc completeness so canonical typed
// declarations win when deduplicating symbols.
type ijsTypingLevel uint8

const (
	ijsTypingTypedJsdoc ijsTypingLevel = iota
	ijsTypingUntypedJsdoc
	ijsTypingNoJsdoc
)

func declTypingLevel(decl *potentialDeclaration) ijsTypingLevel {
	jsdoc := decl.jsdoc()
	if jsdoc == nil {
		return ijsTypingNoJsdoc
	}
	if !jsdoc.HasTypeNodes() {
		return ijsTypingUntypedJsdoc
	}
	return ijsTypingTypedJsdoc
}

func compareDeclarationsFirst(a, b *potentialDeclaration) int {
	return int(declTypingLevel(a)) - int(declTypingLevel(b))
}

func compareShortToLongQName(a, b string) int {
	// '.' is ASCII: counting it in WTF-8 bytes counts UTF-16 units alike.
	return strings.Count(a, ".") - strings.Count(b, ".")
}

func (s *ijsSimplifyDeclarations) removeDuplicateDeclarations() {
	entries := s.currentFile.declarations.Entries()
	for i := range entries {
		name := entries[i].Key
		if strings.HasPrefix(name, "this.") {
			continue
		}
		declList := entries[i].Value
		if len(declList) > 1 {
			slices.SortStableFunc(declList, compareDeclarationsFirst)
			// Retain the highest-ranked declaration at index 0 and remove duplicates.
			entries[i].Value = declList[:1]
			for _, decl := range declList[1:] {
				decl.remove(s.compiler)
			}
		}
	}
}

func (s *ijsSimplifyDeclarations) simplifyAll() {
	s.removeDuplicateDeclarations()

	// Simplify all names in the top-level scope from shortest to longest qualified name.
	seenNames := s.currentFile.declarations.KeySlice()
	slices.SortStableFunc(seenNames, compareShortToLongQName)

	for _, name := range seenNames {
		for _, decl := range s.currentFile.declarations.Get(name) {
			s.processDeclaration(name, decl)
		}
	}
}

func (s *ijsSimplifyDeclarations) processDeclaration(name string, decl *potentialDeclaration) {
	if s.shouldRemove(name, decl) {
		decl.remove(s.compiler)
		return
	}

	if decl.breakDownDestructure(s.compiler) {
		maybeReportIjsDiagnostic(s.compiler, decl.lhs, ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE)
	}

	if fn, ok := ast.AsFunction(decl.rhs); ok {
		s.processFunction(fn)
	} else if cls, ok := ast.AsClass(decl.rhs); ok {
		s.processClass(cls)
	}
	s.setUndeclaredToUnusableType(decl)
	s.elidePrivateMemberType(decl)
	decl.simplify(s.compiler)
}

func (s *ijsSimplifyDeclarations) processClass(cls ast.Class) {
	for member := range ast.ReadAhead(cls.Body().FirstMember()) {
		if cfd, ok := member.AsComputedFieldDef(); ok {
			if val, hasVal := cfd.Value(); hasVal {
				if fn, isFn := val.AsFunction(); isFn {
					s.processFunction(fn)
				} else if !isTypedIjsRhs(val.Node()) {
					NodeUtil_DeleteNode(val.Node(), s.compiler)
				}
			}
		} else if _, ok := member.AsMemberFieldDef(); ok {
			// MEMBER_FIELD_DEF declarations are recorded in ijsPropagateConstJsdoc and
			// simplified individually in simplifyAll.
		} else if _, ok := member.AsEmpty(); ok {
			NodeUtil_DeleteNode(member.Node(), s.compiler)
		} else if _, ok := member.AsBlock(); ok {
			NodeUtil_DeleteNode(member.Node(), s.compiler)
		} else if method, ok := member.AsMethodDef(); ok {
			s.processFunction(method.Function())
		} else if cp, ok := member.AsComputedProp(); ok {
			val, hasVal := cp.Value()
			fn, isFn := ast.AsFunction(val.Node())
			if !hasVal || !isFn {
				panic(jutil.NewIllegalStateExceptionMsg("Non-function computed class member"))
			}
			s.processFunction(fn)
		} else {
			panic(jutil.NewAssertionError(
				member.Token().String() + " should not be handled by processClass"))
		}
	}
}

func (s *ijsSimplifyDeclarations) processFunction(fn ast.Function) {
	s.processFunctionParameters(fn.Params())
}

func (s *ijsSimplifyDeclarations) processFunctionParameters(paramList ast.ParamList) {
	for param := range paramList.Params() {
		if defVal, ok := param.AsDefaultValue(); ok {
			rhs := defVal.Value()
			if call, isCall := rhs.AsCall(); isCall &&
				RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME.Matches(
					call.Callee().Node()) {
				continue
			}
			rhsNode := rhs.Node()
			rhsNode.ReplaceWith(NodeUtil_NewUndefinedNode(rhsNode))
			s.compiler.ReportChangeToEnclosingScope(defVal.Node())
		}
	}
}

func rootOfQualifiedName(qualifiedName string) string {
	if prefix, _, ok := strings.Cut(qualifiedName, "."); ok {
		return prefix
	}
	return qualifiedName
}

func (s *ijsSimplifyDeclarations) shouldRemove(name string, decl *potentialDeclaration) bool {
	if decl.isDetached() {
		return true
	}
	if strings.HasPrefix(rootOfQualifiedName(name), "$jscomp") {
		// These are created by goog.scope processing, but clash with each other
		// and should not be depended on.
		if (!decl.rhs.IsNil() && decl.rhs.IsClass()) ||
			(decl.jsdoc() != nil && decl.jsdoc().ContainsTypeDefinition()) {
			maybeReportIjsDiagnostic(s.compiler, decl.lhs, ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE)
		}
		return true
	}
	// Looks like an update rather than a declaration in this file.
	return !strings.HasPrefix(name, "this.") &&
		!decl.isDefiniteDeclaration() &&
		!decl.lhs.IsMemberFieldDef() &&
		!s.currentFile.isPrefixProvided(name) &&
		!s.currentFile.isStrictPrefixDeclared(name)
}

func (s *ijsSimplifyDeclarations) setUndeclaredToUnusableType(decl *potentialDeclaration) {
	nameNode := decl.lhs
	jsdoc := decl.jsdoc()
	if decl.shouldPreserve() ||
		NodeUtil_IsNamespaceDecl(nameNode) ||
		(!decl.rhs.IsNil() && NodeUtil_IsCallTo(decl.rhs, "Symbol")) ||
		(jsdoc != nil && jsdoc.ContainsDeclaration() && !decl.isConstToBeInferred()) {
		return
	}
	maybeWarnForConstWithoutExplicitType(s.compiler, decl)
	if jsdocNode := NodeUtil_GetBestJSDocInfoNode(nameNode); !jsdocNode.IsNil() {
		jsdocNode.SetJSDocInfo(unusableTypeJSDoc(jsdoc))
	}
}

func (s *ijsSimplifyDeclarations) elidePrivateMemberType(decl *potentialDeclaration) {
	if s.shouldElidePrivateMemberType(decl) {
		if jsdocNode := NodeUtil_GetBestJSDocInfoNode(decl.lhs); !jsdocNode.IsNil() {
			jsdocNode.SetJSDocInfo(replaceAnnotatedTypesWithUnusableType(decl.jsdoc()))
		}
	}
}

func (s *ijsSimplifyDeclarations) shouldElidePrivateMemberType(decl *potentialDeclaration) bool {
	// TODO: b/465855752 - roll out elision for JS files as well after depot cleanup.
	if !s.currentFile.isFromTypeScript {
		return false
	}
	jsdoc := decl.jsdoc()
	if jsdoc == nil || jsdoc.GetVisibility() != rhino.JSDocInfo_Visibility_PRIVATE {
		return false
	}
	// Preserve JSDoc declaring a new type name, to avoid loosening types of other public members.
	if jsdoc.IsConstructorOrInterface() ||
		jsdoc.HasEnumParameterType() ||
		jsdoc.HasTypedefType() {
		return false
	}
	// Simple name declarations marked @private are not really treated as private if exported from
	// a goog.module.
	if decl.lhs.IsName() {
		return false
	}
	return hasAnnotatedType(jsdoc)
}

func isSymbolProp(lhs rhino.Node) bool {
	prop, ok := ast.AsGetProp(lhs)
	return ok && prop.Object().Node().MatchesName("Symbol")
}

func shouldPreserveAssignment(expr rhino.Node, t *NodeTraversal) bool {
	lhs := ast.MustAssignment(expr).Target().Node()
	// Ignore assignments in function bodies, unless they're also a constructor with a `this.prop`.
	if !t.InGlobalHoistScope() && !t.InModuleHoistScope() {
		return isThisPropInsideNamedClass(lhs)
	}

	// Well-known symbol properties, like `Foo.prototype[Symbol.iterator] = function() {};`.
	if elem, ok := ast.AsGetElem(lhs); ok && isSymbolProp(elem.Index().Node()) {
		return elem.Object().Node().IsQualifiedName()
	}
	// Assignments to names don't have global typechecking side-effects even within the hoist scope.
	if lhs.IsName() {
		return t.InGlobalScope() || t.InModuleScope()
	}
	return lhs.IsQualifiedName()
}

// ijsFileInfo tracks provided namespaces, imports, and declarations seen in a single file.
type ijsFileInfo struct {
	providedNamespaces OrderedStringMap[struct{}]
	requiredLocalNames OrderedStringMap[struct{}]
	declarations       OrderedStringMap[[]*potentialDeclaration]
	isFromTypeScript   bool
}

func newIjsFileInfo(sourceFileName string) *ijsFileInfo {
	return &ijsFileInfo{
		isFromTypeScript: strings.HasSuffix(sourceFileName, ".closure.js"),
	}
}

func (f *ijsFileInfo) recordNameDeclaration(qualifiedNameNode rhino.Node) {
	f.recordDeclaration(potentialDeclFromName(qualifiedNameNode))
}

func (f *ijsFileInfo) recordMemberFieldDef(field ast.MemberFieldDef) {
	f.recordDeclaration(potentialDeclFromMemberField(field))
}

func (f *ijsFileInfo) recordMethod(method ast.MemberFunctionDef) {
	f.recordDeclaration(potentialDeclFromMethod(method))
}

func (f *ijsFileInfo) recordStringKeyDeclaration(sk ast.StringKey) {
	f.recordDeclaration(potentialDeclFromStringKey(sk))
}

func (f *ijsFileInfo) recordDefine(call ast.Call) {
	if decl := potentialDeclFromDefine(call); decl != nil {
		f.recordDeclaration(decl)
	}
}

func (f *ijsFileInfo) recordAliasDeclaration(nameNode rhino.Node) {
	f.recordDeclaration(potentialDeclFromAlias(nameNode))
}

func (f *ijsFileInfo) recordDeclaration(decl *potentialDeclaration) {
	name := decl.fullyQualifiedName
	f.declarations.Put(name, append(f.declarations.Get(name), decl))
}

func (f *ijsFileInfo) recordImport(localName string) {
	f.requiredLocalNames.Put(localName, struct{}{})
}

func containsNamespacePrefix(fullyQualifiedName string, prefixNamespaces *OrderedStringMap[struct{}]) bool {
	if prefixNamespaces.Size() == 0 {
		return false
	}
	if prefixNamespaces.ContainsKey(fullyQualifiedName) {
		return true
	}
	for i := 0; i < len(fullyQualifiedName); i++ {
		if fullyQualifiedName[i] == '.' && prefixNamespaces.ContainsKey(fullyQualifiedName[:i]) {
			return true
		}
	}
	return false
}

func (f *ijsFileInfo) isPrefixProvided(fullyQualifiedName string) bool {
	return containsNamespacePrefix(fullyQualifiedName, &f.providedNamespaces)
}

func (f *ijsFileInfo) isStrictPrefixDeclared(fullyQualifiedName string) bool {
	if f.declarations.Size() == 0 {
		return false
	}
	for i := 0; i < len(fullyQualifiedName); i++ {
		if fullyQualifiedName[i] == '.' && f.declarations.ContainsKey(fullyQualifiedName[:i]) {
			return true
		}
	}
	return false
}

func (f *ijsFileInfo) markProvided(providedName string) {
	f.providedNamespaces.Put(providedName, struct{}{})
}

func ijsNodeQualifiedName(n rhino.Node) (string, bool) {
	if n.IsNil() {
		return "", false
	}
	nameNode := NodeUtil_GetNameNode(n)
	if nameNode.IsNil() {
		return "", false
	}
	return nameNode.GetQualifiedNameString()
}

// isThisPropInsideNamedClass reports whether maybeGetProp is a `this.prop`
// access inside a named class or constructor.
func isThisPropInsideNamedClass(maybeGetProp rhino.Node) bool {
	_, ok := classNameOfThisProp(maybeGetProp)
	return ok
}

func classNameOfThisProp(getprop rhino.Node) (string, bool) {
	gp, ok := ast.AsGetProp(getprop)
	if !ok || !gp.Object().Node().IsThis() {
		return "", false
	}
	function := NodeUtil_GetEnclosingFunction(getprop)
	if function.IsNil() {
		return "", false
	}
	className, ok := memberFunctionClassName(function)
	if !ok || className == "" {
		return "", false
	}
	return className, true
}

func isMemberFieldInsideNamedClass(fieldNode rhino.Node) bool {
	field, ok := ast.AsMemberFieldDef(fieldNode)
	if !ok {
		return false
	}
	_, hasName := memberFieldClassName(field)
	return hasName
}

func isComputedMemberInsideClass(fieldNode rhino.Node) bool {
	cm, ok := ast.AsComputedMember(fieldNode)
	return ok && ast.RoleOf(fieldNode) == ast.RoleClassBodyMembers && isSymbolProp(cm.Key().Node())
}

func memberFieldQualifiedName(field ast.MemberFieldDef) string {
	className, ok := memberFieldClassName(field)
	if !ok {
		panic(jutil.NewIllegalArgumentExceptionMsg("expected field inside named class"))
	}
	if field.IsStatic() {
		return className + "." + field.Key()
	}
	return className + ".prototype." + field.Key()
}

func memberFieldClassName(field ast.MemberFieldDef) (string, bool) {
	parent := field.Node().GetParent()
	if parent.IsNil() {
		return "", false
	}
	classNode, ok := ast.AsClass(parent.GetParent())
	if !ok {
		return "", false
	}
	return ijsNodeQualifiedName(classNode.Node())
}

func classMethodQualifiedName(method ast.MemberFunctionDef) string {
	fn := method.Function().Node()
	if !isClassMethod(fn) {
		panic(jutil.NewIllegalArgumentExceptionMsg("expected class method"))
	}
	className, ok := memberFunctionClassName(fn)
	if !ok || className == "" {
		panic(jutil.NewIllegalStateExceptionMsg("expected named class for method"))
	}
	if method.IsStatic() {
		return className + "." + method.Key()
	}
	return className + ".prototype." + method.Key()
}

func isClassMethod(functionNode rhino.Node) bool {
	if !functionNode.IsFunction() {
		return false
	}
	return ast.RoleOf(functionNode) == ast.RoleMemberFunctionDefFunction &&
		ast.RoleOf(functionNode.GetParent()) == ast.RoleClassBodyMembers
}

func hasNamedClass(functionNode rhino.Node) bool {
	if !functionNode.IsFunction() {
		return false
	}
	_, ok := memberFunctionClassName(functionNode)
	return ok
}

func memberFunctionClassName(functionNode rhino.Node) (string, bool) {
	if !functionNode.IsFunction() {
		return "", false
	}
	if isClassMethod(functionNode) {
		// ES6 class method: FUNCTION -> MEMBER_FUNCTION_DEF -> CLASS_MEMBERS -> CLASS.
		classMembers := functionNode.GetGrandparent()
		if classMembers.IsNil() {
			return "", false
		}
		classNode, ok := ast.AsClass(classMembers.GetParent())
		if !ok {
			return "", false
		}
		return ijsNodeQualifiedName(classNode.Node())
	}
	return ijsNodeQualifiedName(functionNode)
}

func isIjsConstructor(functionNode rhino.Node) bool {
	if isClassMethod(functionNode) {
		return NodeUtil_IsEs6Constructor(functionNode)
	}
	jsdoc := NodeUtil_GetBestJSDocInfo(functionNode)
	return jsdoc != nil && jsdoc.IsConstructorOrInterface()
}
