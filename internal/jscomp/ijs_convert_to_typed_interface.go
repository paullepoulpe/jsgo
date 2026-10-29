// Port of com/google/javascript/jscomp/ijs/ConvertToTypedInterface.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/ijs/FileInfo.java @ closure-compiler v20261021

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

// ConvertToTypedInterface: the goal of this pass is to shrink the AST,
// preserving only typing, not behavior.
//
// To do this, it does things like removing function/method bodies, rvalues
// that are not needed, expressions that are not declarations, etc.
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

var convertToTypedInterface_CALLS_TO_PRESERVE = jutil.NewLinkedHashSet[string](
	"Polymer",
	"goog.addSingletonGetter",
	"goog.define",
	"goog.forwardDeclare",
	"goog.module",
	"goog.module.declareLegacyNamespace",
	"goog.declareModuleId",
	"goog.provide",
	"goog.require",
	"goog.requireType")

// Java: ConvertToTypedInterface(AbstractCompiler)
func NewConvertToTypedInterface(compiler AbstractCompiler) *ConvertToTypedInterface {
	c := &ConvertToTypedInterface{}
	c.compiler = compiler
	return c
}

// Java: ConvertToTypedInterface.maybeReport(AbstractCompiler, Node, DiagnosticType, String...)
func ConvertToTypedInterface_MaybeReport(
	compiler AbstractCompiler, node rhino.Node, diagnostic *DiagnosticType, fillers ...string) {
	sourceName := *NodeUtil_GetSourceName(node)
	if strings.HasSuffix(sourceName, "_test.js") ||
		strings.HasSuffix(sourceName, "_test.closure.js") ||
		strings.HasSuffix(sourceName, "_test.tsx.cl.js") {
		// Allow _test.js files and their tsickle generated
		// equivalents to avoid emitting errors at .i.js generation time.
		// We expect these files to not be consumed by any other downstream libraries.
		return
	}
	compiler.Report(JSError_MakeWithNode(node, diagnostic, fillers...))
}

// Java: ConvertToTypedInterface.maybeWarnForConstWithoutExplicitType(AbstractCompiler, PotentialDeclaration)
func ConvertToTypedInterface_MaybeWarnForConstWithoutExplicitType(
	compiler AbstractCompiler, decl *PotentialDeclaration) {
	if decl.IsConstToBeInferred() &&
		!decl.GetLhs().IsFromExterns() &&
		!JsdocUtil_IsPrivate(decl.GetJsDoc()) {

		nameNode := decl.GetLhs()
		if nameNode.GetJSType() == nil {
			ConvertToTypedInterface_MaybeReport(compiler, nameNode, ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE)
		} else {
			ConvertToTypedInterface_MaybeReport(
				compiler,
				nameNode,
				ConvertToTypedInterface_CONSTANT_WITH_SUGGESTED_TYPE,
				nameNode.GetJSType().ToAnnotationString(rhino.JSType_Nullability_EXPLICIT))
		}
	}
}

// Java: ConvertToTypedInterface.process(Node, Node)
func (c *ConvertToTypedInterface) Process(externs rhino.Node, root rhino.Node) {
	for script := root.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		c.ProcessFile(script)
	}
}

// Java: ConvertToTypedInterface.processFile(Node)
func (c *ConvertToTypedInterface) ProcessFile(scriptNode rhino.Node) {
	jutil.CheckArgument(scriptNode.IsScript())
	sourceFileName := scriptNode.GetSourceFileName()
	if AbstractCompiler_IsFillFileName(*sourceFileName) {
		scriptNode.Detach()
		return
	}

	scriptJsDoc := scriptNode.GetJSDocInfo()
	if scriptJsDoc != nil && scriptJsDoc.IsClosureUnawareCode() {
		// If we are generating type summary files, then those files won't contain the method content
		// from within the closure-unaware section, and the entire file effectively becomes
		// closure-aware again (as it is generated code that describes a module shape).
		scriptJsDocBuilder := scriptJsDoc.ToBuilder()
		scriptJsDocBuilder.RemoveClosureUnawareCode()
		scriptNode.SetJSDocInfo(scriptJsDocBuilder.Build())
	}

	currentFile := NewFileInfo(*sourceFileName)
	NodeTraversal_Traverse(c.compiler, scriptNode, NewConvertToTypedInterface_RemoveNonDeclarations())
	NodeTraversal_Traverse(c.compiler, scriptNode, NewConvertToTypedInterface_PropagateConstJsdoc(currentFile))
	NewConvertToTypedInterface_SimplifyDeclarations(c.compiler, currentFile).SimplifyAll()
}

// Java: ConvertToTypedInterface.findNameDeclaration(Scope, Node)
func ConvertToTypedInterface_FindNameDeclaration(scope *Scope, rhs rhino.Node) *Var {
	if !rhs.IsName() {
		return nil
	}
	return scope.GetVar(rhs.GetString())
}

// ConvertToTypedInterface_RemoveNonDeclarations is the private static nested
// callback ConvertToTypedInterface.RemoveNonDeclarations.
type ConvertToTypedInterface_RemoveNonDeclarations struct {
	_ byte // Go-only: language-adapter — distinct identity
}

var convertToTypedInterface_RemoveNonDeclarations_GOOG_SCOPE = rhino.QualifiedName_Of("goog.scope")

// Java: ConvertToTypedInterface.RemoveNonDeclarations()
func NewConvertToTypedInterface_RemoveNonDeclarations() *ConvertToTypedInterface_RemoveNonDeclarations {
	return &ConvertToTypedInterface_RemoveNonDeclarations{}
}

// Java: ConvertToTypedInterface.RemoveNonDeclarations.shouldTraverse(NodeTraversal, Node, Node)
func (r *ConvertToTypedInterface_RemoveNonDeclarations) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		if !ClassUtil_IsConstructor(n) || !ClassUtil_HasNamedClass(n) {
			body := n.GetLastChild()
			if !body.IsBlock() || body.HasChildren() {
				t.ReportCodeChangeWithNode(body)
				body.ReplaceWith(rhino.IR_Block().Srcref(body))
				NodeUtil_MarkFunctionsDeleted(body, t.GetCompiler())
			}
		}
		return true
	case rhino.Token_MEMBER_FIELD_DEF:
		if ClassUtil_IsMemberFieldDefInsideClassWithName(n) {
			return true
		}
		// We will remove fields in anonymous classes
		NodeUtil_DeleteNode(n, t.GetCompiler())
		return false
	case rhino.Token_EXPR_RESULT:
		expr := n.GetFirstChild()
		switch expr.GetToken() {
		case rhino.Token_CALL:
			callee := expr.GetFirstChild()
			jutil.CheckState(!convertToTypedInterface_RemoveNonDeclarations_GOOG_SCOPE.Matches(callee))
			if qname := callee.GetQualifiedName(); qname != nil && convertToTypedInterface_CALLS_TO_PRESERVE.Contains(*qname) {
				return true
			}
			NodeUtil_DeleteNode(n, t.GetCompiler())
			return false
		case rhino.Token_ASSIGN:
			if ConvertToTypedInterface_ShouldPreserveAssignment(expr, t) {
				return true
			}
			NodeUtil_DeleteNode(n, t.GetCompiler())
			return false
		case rhino.Token_GETPROP:
			if !expr.IsQualifiedName() || expr.GetJSDocInfo() == nil {
				NodeUtil_DeleteNode(n, t.GetCompiler())
				return false
			}
			return true
		case rhino.Token_GETELEM:
			if ConvertToTypedInterface_IsSymbolProp(expr.GetSecondChild()) && expr.GetJSDocInfo() != nil {
				return true
			}

			NodeUtil_DeleteNode(n, t.GetCompiler())
			return false
		default:
			NodeUtil_DeleteNode(n, t.GetCompiler())
			return false
		}
	case rhino.Token_COMPUTED_PROP:
		if ClassUtil_IsComputedMemberInsideClassWithName(n) {
			return true
		}
		if !NodeUtil_IsLhsByDestructuring(n.GetSecondChild()) {
			NodeUtil_DeleteNode(n, t.GetCompiler())
		}
		return false
	case rhino.Token_COMPUTED_FIELD_DEF:
		if ClassUtil_IsComputedMemberInsideClassWithName(n) {
			return true
		}
		NodeUtil_DeleteNode(n, t.GetCompiler())
		return false
	case rhino.Token_THROW, rhino.Token_RETURN, rhino.Token_BREAK, rhino.Token_CONTINUE, rhino.Token_DEBUGGER, rhino.Token_EMPTY:
		if NodeUtil_IsStatementParent(parent) {
			NodeUtil_DeleteNode(n, t.GetCompiler())
		}
		return false
	case rhino.Token_LABEL, rhino.Token_IF, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_WHILE:
		// First child can't have declaration. Statement itself will be removed post-order.
		NodeUtil_DeleteNode(n.GetFirstChild(), t.GetCompiler())
		return true
	case rhino.Token_TRY, rhino.Token_DO:
		// Second child can't have declarations. Statement itself will be removed post-order.
		NodeUtil_DeleteNode(n.GetSecondChild(), t.GetCompiler())
		return true
	case rhino.Token_FOR:
		NodeUtil_DeleteNode(n.GetSecondChild(), t.GetCompiler())
		// fall-through
		fallthrough
	case rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_FOR_IN:
		NodeUtil_DeleteNode(n.GetSecondChild(), t.GetCompiler())
		initializer := n.RemoveFirstChild()
		if initializer.IsVar() {
			n.GetLastChild().AddChildToFront(initializer)
		}
		return true
	case rhino.Token_CONST, rhino.Token_LET:
		if !t.InGlobalScope() && !t.InModuleScope() {
			NodeUtil_RemoveChild(parent, n)
			t.ReportCodeChangeWithNode(parent)
			return false
		}
		return true
	case rhino.Token_VAR:
		if !t.InGlobalHoistScope() && !t.InModuleHoistScope() {
			NodeUtil_RemoveChild(parent, n)
			t.ReportCodeChangeWithNode(parent)
			return false
		}
		return true
	case rhino.Token_MODULE_BODY, rhino.Token_CLASS, rhino.Token_DEFAULT_CASE, rhino.Token_BLOCK, rhino.Token_EXPORT, rhino.Token_IMPORT:
		return true
	default:
		jutil.CheckStateMsg1(!NodeUtil_IsStatement(n), "%s", n.GetToken())
		return true
	}
}

// Java: ConvertToTypedInterface.RemoveNonDeclarations.visit(NodeTraversal, Node, Node)
func (r *ConvertToTypedInterface_RemoveNonDeclarations) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_TRY, rhino.Token_LABEL, rhino.Token_DEFAULT_CASE, rhino.Token_CASE, rhino.Token_DO, rhino.Token_WHILE,
		rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_IF:
		if n.HasParent() {
			children := n.RemoveChildren()
			parent.AddChildrenAfter(children, n)
			// We don't need the special valid-AST-preserving behavior of `NodeUtil.removeChild()`
			// here. We always want to remove exactly and only `n`, so we use `n.detach()`.
			//
			// Also, the shouldTraverse() method intentionally puts the AST in an invalid state
			// by moving children to parents where they are not valid temporarily.
			// `NodeUtil.removeChild()` would throw an exception here if it noticed the invalid AST
			// state.
			n.Detach()
			t.ReportCodeChange()
		}
	case rhino.Token_SWITCH:
		// shouldTraverse() removed the switch condition already, so we just need to handle the
		// cases.
		if n.HasParent() {
			if !n.HasOneChild() {
				jutil.CheckStateMsg1(false, "malfrmed sWITCH %s", n.ToStringTree())
			}
			children := n.GetFirstChild().RemoveChildren()
			parent.AddChildrenAfter(children, n)
			n.Detach()
			t.ReportCodeChange()
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		ConvertToTypedInterface_RemoveNonDeclarations_SplitNameDeclarationsAndRemoveDestructuring(n, t)
	case rhino.Token_BLOCK:
		if !parent.IsFunction() {
			parent.AddChildrenAfter(n.RemoveChildren(), n)
			n.Detach()
			t.ReportCodeChangeWithNode(parent)
		}
	default:
	}
}

// ConvertToTypedInterface_RemoveNonDeclarations_SplitNameDeclarationsAndRemoveDestructuring
// does three simplifications to const/let/var nodes.
//
//   - Splits them so that each declaration is a separate statement.
//   - Removes non-import and non-alias destructuring statements, which we
//     assume are not type declarations.
//   - Moves inline JSDoc annotations onto the declaration nodes.
//
// Java: ConvertToTypedInterface.RemoveNonDeclarations.splitNameDeclarationsAndRemoveDestructuring(Node, NodeTraversal)
func ConvertToTypedInterface_RemoveNonDeclarations_SplitNameDeclarationsAndRemoveDestructuring(n rhino.Node, t *NodeTraversal) {
	jutil.CheckArgument(NodeUtil_IsNameDeclaration(n))
	sharedJsdoc := n.GetJSDocInfo()
	isExport := n.GetParent().IsExport()
	statement := n
	if isExport {
		statement = n.GetParent()
	}
	for n.HasChildren() {
		lhsToSplit := n.GetLastChild()
		nameJsdoc := lhsToSplit.GetJSDocInfo()
		lhsToSplit.SetJSDocInfo(nil)
		mergedJsdoc := JsdocUtil_MergeJsdocs(sharedJsdoc, nameJsdoc)
		if n.HasOneChild() {
			n.SetJSDocInfo(mergedJsdoc)
			return
		}
		// A name declaration with more than one LHS is split into separate declarations.
		var rhs rhino.Node
		if lhsToSplit.HasChildren() {
			rhs = lhsToSplit.RemoveFirstChild()
		}
		newDeclaration :=
			NodeUtil_NewDeclaration(lhsToSplit.Detach(), rhs, n.GetToken()).Srcref(n)
		newDeclaration.SetJSDocInfo(mergedJsdoc)
		if isExport {
			newDeclaration = rhino.IR_Export(newDeclaration).Srcref(statement)
		}
		newDeclaration.InsertAfter(statement)
		t.ReportCodeChange()
	}
}

// IjsErrors holds IJS-related errors occurring outside the
// ConvertToTypedInterface pass.
type IjsErrors struct{}

var (
	IjsErrors_BAD_IJS_FILE_NAME = DiagnosticType_Error("JSC_BAD_IJS_FILE_NAME", "--ijs file name must end in .i.js: {0}")

	IjsErrors_CONFLICTING_IJS_FILE = DiagnosticType_Warning(
		"JSC_CONFLICTING_IJS_FILE",
		"Found --ijs file in a conflict with --js or --weakdep file. This usually means the same"+
			" file appears both in the srcs and the transitive deps of a build target.\n"+
			"File 1: {0}\n"+
			"File 2: {1}")
)

// ConvertToTypedInterface_PropagateConstJsdoc is the private static nested
// class ConvertToTypedInterface.PropagateConstJsdoc.
type ConvertToTypedInterface_PropagateConstJsdoc struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	currentFile *FileInfo
}

var processConstJsdocCallback_GOOG_DEFINE = rhino.QualifiedName_Of("goog.define")
var processConstJsdocCallback_GOOG_PROVIDE = rhino.QualifiedName_Of("goog.provide")
var processConstJsdocCallback_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")
var processConstJsdocCallback_CJS_REQUIRE = rhino.QualifiedName_Of("require")

// Java: ConvertToTypedInterface.PropagateConstJsdoc(FileInfo)
func NewConvertToTypedInterface_PropagateConstJsdoc(currentFile *FileInfo) *ConvertToTypedInterface_PropagateConstJsdoc {
	return &ConvertToTypedInterface_PropagateConstJsdoc{currentFile: currentFile}
}

// Java: ProcessConstJsdocCallback.visit(NodeTraversal, Node, Node)
func (p *ConvertToTypedInterface_PropagateConstJsdoc) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		if NodeUtil_IsStatementParent(parent) {
			p.currentFile.RecordNameDeclaration(ast.MustClass(n).NameNode())
		}
	case rhino.Token_MEMBER_FIELD_DEF:
		if !NodeUtil_GetRValueOfLValue(n).IsNil() {
			p.ProcessDeclarationWithRhs(t, n)
		}
		p.currentFile.RecordMemberFieldDef(n)
	case rhino.Token_FUNCTION:
		if NodeUtil_IsStatementParent(parent) {
			p.currentFile.RecordNameDeclaration(ast.MustFunction(n).NameNode().Node())
		} else if ClassUtil_IsClassMethod(n) && ClassUtil_HasNamedClass(n) {
			p.currentFile.RecordMethod(n)
		}
	case rhino.Token_EXPR_RESULT:
		expr := ast.MustExprStmt(n).Expr().Node()
		switch expr.GetToken() {
		case rhino.Token_CALL:
			call := ast.MustCall(expr)
			callee := call.Callee().Node()
			if processConstJsdocCallback_GOOG_PROVIDE.Matches(callee) {
				lastArg, _ := call.LastArg()
				p.currentFile.MarkProvided(lastArg.Node().GetString())
			} else if processConstJsdocCallback_GOOG_REQUIRE.Matches(callee) || processConstJsdocCallback_CJS_REQUIRE.Matches(callee) {
				lastArg, _ := call.LastArg()
				p.currentFile.RecordImport(lastArg.Node().GetString())
			} else if processConstJsdocCallback_GOOG_DEFINE.Matches(callee) {
				p.currentFile.RecordDefine(expr)
			}
		case rhino.Token_ASSIGN:
			assign := ast.MustAssign(expr)
			p.RecordDeclaration(t, assign.Target().Node(), assign.Value().Node())
		case rhino.Token_GETPROP:
			p.currentFile.RecordNameDeclaration(expr)
		case rhino.Token_GETELEM:
		default:
			panic(jutil.NewIllegalArgumentExceptionMsg("Unexpected declaration: " + expr.String()))
		}
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		jutil.CheckStateMsg1(n.HasOneChild(), "%s", n)
		p.RecordDeclaration(t, n.GetFirstChild(), n.GetFirstChild().GetLastChild())
	case rhino.Token_STRING_KEY:
		if parent.IsObjectLit() && n.HasOneChild() {
			p.ProcessDeclarationWithRhs(t, n)
			p.currentFile.RecordStringKeyDeclaration(n)
		}
	default:
	}
}

// Java: ProcessConstJsdocCallback.recordDeclaration(NodeTraversal, Node, Node)
func (p *ConvertToTypedInterface_PropagateConstJsdoc) RecordDeclaration(t *NodeTraversal, lhs rhino.Node, rhs rhino.Node) {
	if !rhs.IsNil() &&
		rhs.IsCall() &&
		processConstJsdocCallback_GOOG_DEFINE.Matches(rhs.GetFirstChild()) &&
		lhs.IsQualifiedName() {
		p.currentFile.RecordDefine(rhs)
	} else if lhs.IsGetElem() {
		return
	} else {
		p.RecordNameDeclaration(lhs, rhs)
		if !lhs.IsDestructuringLhs() && !rhs.IsNil() {
			p.ProcessDeclarationWithRhs(t, lhs)
		}
	}
}

// Java: ProcessConstJsdocCallback.recordNameDeclaration(Node, Node)
func (p *ConvertToTypedInterface_PropagateConstJsdoc) RecordNameDeclaration(lhs rhino.Node, rhs rhino.Node) {
	jutil.CheckArgument(NodeUtil_IsNameDeclaration(lhs.GetParent()) || lhs.GetParent().IsAssign())
	isImport := PotentialDeclaration_IsImportRhs(rhs)
	isAlias := PotentialDeclaration_IsAliasDeclaration(lhs, rhs)
	NodeUtil_VisitLhsNodesInNode(
		lhs.GetParent(),
		func(name rhino.Node) {
			if isAlias || isImport {
				p.currentFile.RecordAliasDeclaration(name)
			} else {
				p.currentFile.RecordNameDeclaration(name)
			}
		})
}

// Java: ProcessConstJsdocCallback.processDeclarationWithRhs(NodeTraversal, Node)
func (p *ConvertToTypedInterface_PropagateConstJsdoc) ProcessDeclarationWithRhs(t *NodeTraversal, lhs rhino.Node) {
	jutil.CheckArgumentMsg1(lhs.IsQualifiedName() || lhs.IsStringKey() || lhs.IsMemberFieldDef(), "%s", lhs)
	jutil.CheckStateMsg1(!NodeUtil_GetRValueOfLValue(lhs).IsNil(), "%s", lhs)
	if PotentialDeclaration_IsConstToBeInferred(lhs) {
		p.ProcessConstWithRhs(t, lhs)
	}
}

// Java: ConvertToTypedInterface.PropagateConstJsdoc.processConstWithRhs(NodeTraversal, Node)
func (p *ConvertToTypedInterface_PropagateConstJsdoc) ProcessConstWithRhs(t *NodeTraversal, nameNode rhino.Node) {
	jutil.CheckArgumentMsg1(
		nameNode.IsQualifiedName() ||
			nameNode.IsStringKey() ||
			nameNode.IsDestructuringLhs() ||
			nameNode.IsMemberFieldDef(),
		"%s",
		nameNode)
	jsdocNode := NodeUtil_GetBestJSDocInfoNode(nameNode)
	originalJsdoc := jsdocNode.GetJSDocInfo()
	rhs := NodeUtil_GetRValueOfLValue(nameNode)
	newJsdoc := JsdocUtil_GetJSDocForRhs(rhs, originalJsdoc)
	if newJsdoc == nil && ClassUtil_IsThisPropInsideClassWithName(nameNode) {
		decl := ConvertToTypedInterface_FindNameDeclaration(t.GetScope(), rhs)
		newJsdoc = JsdocUtil_GetJSDocForName(decl, originalJsdoc)
	}
	if newJsdoc != nil {
		jsdocNode.SetJSDocInfo(newJsdoc)
		t.ReportCodeChange()
	}
}

// ConvertToTypedInterface_SimplifyDeclarations is the private static nested
// class ConvertToTypedInterface.SimplifyDeclarations.
type ConvertToTypedInterface_SimplifyDeclarations struct {
	compiler    AbstractCompiler
	currentFile *FileInfo
}

// ConvertToTypedInterface_SimplifyDeclarations_TypingLevel lists the levels
// of JSDoc, starting from those most likely to be on the canonical
// declaration.
//
// Java: ConvertToTypedInterface.SimplifyDeclarations.TypingLevel
type ConvertToTypedInterface_SimplifyDeclarations_TypingLevel int32

const (
	ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_TYPED_JSDOC_DECLARATION ConvertToTypedInterface_SimplifyDeclarations_TypingLevel = iota
	ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_UNTYPED_JSDOC_DECLARATION
	ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_NO_JSDOC
)

var convertToTypedInterface_SimplifyDeclarations_TypingLevel_names = [...]string{
	"TYPED_JSDOC_DECLARATION", "UNTYPED_JSDOC_DECLARATION", "NO_JSDOC",
}

// String is Enum.name().
func (l ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) String() string {
	return convertToTypedInterface_SimplifyDeclarations_TypingLevel_names[l]
}

// Java: ConvertToTypedInterface.SimplifyDeclarations.countDots(String)
func ConvertToTypedInterface_SimplifyDeclarations_CountDots(name string) int {
	// '.' is ASCII: counting it in the WTF-8 bytes counts the UTF-16 units alike.
	return strings.Count(name, ".")
}

// ConvertToTypedInterface_SimplifyDeclarations_SHORT_TO_LONG is
// comparing(SimplifyDeclarations::countDots).
var ConvertToTypedInterface_SimplifyDeclarations_SHORT_TO_LONG = func(a, b string) int {
	return ConvertToTypedInterface_SimplifyDeclarations_CountDots(a) - ConvertToTypedInterface_SimplifyDeclarations_CountDots(b)
}

// convertToTypedInterface_SimplifyDeclarations_typingLevel is the key
// extractor of DECLARATIONS_FIRST.
func convertToTypedInterface_SimplifyDeclarations_typingLevel(decl *PotentialDeclaration) ConvertToTypedInterface_SimplifyDeclarations_TypingLevel {
	jsdoc := decl.GetJsDoc()
	if jsdoc == nil {
		return ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_NO_JSDOC
	}
	if !jsdoc.HasTypeNodes() {
		return ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_UNTYPED_JSDOC_DECLARATION
	}
	return ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_TYPED_JSDOC_DECLARATION
}

// ConvertToTypedInterface_SimplifyDeclarations_DECLARATIONS_FIRST is
// comparing(decl -> <its TypingLevel>), enums comparing by ordinal.
var ConvertToTypedInterface_SimplifyDeclarations_DECLARATIONS_FIRST = func(a, b *PotentialDeclaration) int {
	return int(convertToTypedInterface_SimplifyDeclarations_typingLevel(a)) - int(convertToTypedInterface_SimplifyDeclarations_typingLevel(b))
}

// Java: ConvertToTypedInterface.SimplifyDeclarations(AbstractCompiler, FileInfo)
func NewConvertToTypedInterface_SimplifyDeclarations(compiler AbstractCompiler, currentFile *FileInfo) *ConvertToTypedInterface_SimplifyDeclarations {
	s := &ConvertToTypedInterface_SimplifyDeclarations{}
	s.compiler = compiler
	s.currentFile = currentFile
	return s
}

// Java: ConvertToTypedInterface.SimplifyDeclarations.removeDuplicateDeclarations()
func (s *ConvertToTypedInterface_SimplifyDeclarations) RemoveDuplicateDeclarations() {
	for name, declList := range s.currentFile.GetDeclarations().All() {
		if strings.HasPrefix(name, "this.") {
			continue
		}
		slices.SortStableFunc(declList, ConvertToTypedInterface_SimplifyDeclarations_DECLARATIONS_FIRST)
		if len(declList) > 1 {
			// Don't remove the first declaration (at index 0)
			s.currentFile.GetDeclarations().Put(name, declList[:1])
			for _, decl := range declList[1:] {
				decl.Remove(s.compiler)
			}
		}
	}
}

// Java: ConvertToTypedInterface.SimplifyDeclarations.simplifyAll()
func (s *ConvertToTypedInterface_SimplifyDeclarations) SimplifyAll() {
	// Remove duplicate assignments to the same symbol
	s.RemoveDuplicateDeclarations()

	// Simplify all names in the top-level scope.
	seenNames := slices.Collect(s.currentFile.GetDeclarations().Keys())
	slices.SortStableFunc(seenNames, ConvertToTypedInterface_SimplifyDeclarations_SHORT_TO_LONG)

	for _, name := range seenNames {
		for _, decl := range s.currentFile.GetDeclarations().Get(name) {
			s.ProcessDeclaration(name, decl)
		}
	}
}

// Java: ConvertToTypedInterface.SimplifyDeclarations.processDeclaration(String, PotentialDeclaration)
func (s *ConvertToTypedInterface_SimplifyDeclarations) ProcessDeclaration(name string, decl *PotentialDeclaration) {
	if s.ShouldRemove(name, decl) {
		decl.Remove(s.compiler)
		return
	}

	if decl.BreakDownDestructure(s.compiler) {
		ConvertToTypedInterface_MaybeReport(s.compiler, decl.GetLhs(), ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE)
	}

	if fn, ok := ast.AsFunction(decl.GetRhs()); ok {
		s.ProcessFunction(fn)
	} else if cls, ok := ast.AsClass(decl.GetRhs()); ok {
		s.ProcessClass(cls)
	}
	s.SetUndeclaredToUnusableType(decl)
	s.ElidePrivateMemberType(decl)
	decl.Simplify(s.compiler)
}

// Java: ConvertToTypedInterface.SimplifyDeclarations.processClass(Node)
func (s *ConvertToTypedInterface_SimplifyDeclarations) ProcessClass(cls ast.Class) {
	for member := range ast.ReadAhead(cls.Body().FirstMember()) {
		if _, ok := member.AsFieldDef(); ok {
			// MEMBER_FIELD_DEF's are handled in ProcessConstJsdocCallback which is called in the
			// traversal by its subclass PropogateConstJsdoc.
			// If no jsdoc is present on the MEMBER_FIELD_DEF, we will be add a Jsdoc in
			// `setUndeclaredToUnusableType()`.
			// No further simplification is needed for MEMBER_FIELD_DEF's when we call
			// `simplifyAll()` so we will break.
		} else if _, ok := member.AsEmpty(); ok {
			// a lonely `;` in a class body can just be deleted.
			NodeUtil_DeleteNode(member.Node(), s.compiler)
		} else if method, ok := member.AsMethodDef(); ok {
			s.ProcessFunction(method.Function())
		} else if cp, ok := member.AsComputedProp(); ok {
			val, hasVal := cp.Value()
			fn, isFn := ast.AsFunction(val.Node())
			jutil.CheckStateMsg1(
				hasVal && isFn,
				"Non-function computed class member: %s",
				member.Node())
			s.ProcessFunction(fn)
		} else {
			panic(jutil.NewAssertionError(
				member.Token().String() + " should not be handled by processClass"))
		}
	}
}

// Java: ConvertToTypedInterface.SimplifyDeclarations.processFunction(Node)
func (s *ConvertToTypedInterface_SimplifyDeclarations) ProcessFunction(fn ast.Function) {
	s.ProcessFunctionParameters(fn.Params())
}

// Java: ConvertToTypedInterface.SimplifyDeclarations.processFunctionParameters(Node)
func (s *ConvertToTypedInterface_SimplifyDeclarations) ProcessFunctionParameters(paramList ast.ParamList) {
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

// Java: ConvertToTypedInterface.SimplifyDeclarations.rootName(String)
func ConvertToTypedInterface_SimplifyDeclarations_RootName(qualifiedName string) string {
	if prefix, _, ok := strings.Cut(qualifiedName, "."); ok {
		return prefix
	}
	return qualifiedName
}

// Java: ConvertToTypedInterface.SimplifyDeclarations.shouldRemove(String, PotentialDeclaration)
func (s *ConvertToTypedInterface_SimplifyDeclarations) ShouldRemove(name string, decl *PotentialDeclaration) bool {
	if decl.IsDetached() {
		return true
	}
	if strings.HasPrefix(ConvertToTypedInterface_SimplifyDeclarations_RootName(name), "$jscomp") {
		// These are created by goog.scope processing, but clash with each other
		// and should not be depended on.
		if (!decl.GetRhs().IsNil() && decl.GetRhs().IsClass()) ||
			(decl.GetJsDoc() != nil && decl.GetJsDoc().ContainsTypeDefinition()) {
			ConvertToTypedInterface_MaybeReport(s.compiler, decl.GetLhs(), ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE)
		}
		return true
	}
	// This looks like an update rather than a declaration in this file.
	return !strings.HasPrefix(name, "this.") &&
		!decl.IsDefiniteDeclaration(s.compiler) &&
		!decl.GetLhs().IsMemberFieldDef() &&
		!s.currentFile.IsPrefixProvided(name) &&
		!s.currentFile.IsStrictPrefixDeclared(name)
}

// Java: ConvertToTypedInterface.SimplifyDeclarations.setUndeclaredToUnusableType(PotentialDeclaration)
func (s *ConvertToTypedInterface_SimplifyDeclarations) SetUndeclaredToUnusableType(decl *PotentialDeclaration) {
	nameNode := decl.GetLhs()
	jsdoc := decl.GetJsDoc()
	if decl.ShouldPreserve() ||
		NodeUtil_IsNamespaceDecl(nameNode) ||
		(!decl.GetRhs().IsNil() && NodeUtil_IsCallTo(decl.GetRhs(), "Symbol")) ||
		(jsdoc != nil && jsdoc.ContainsDeclaration() && !decl.IsConstToBeInferred()) {
		return
	}
	ConvertToTypedInterface_MaybeWarnForConstWithoutExplicitType(s.compiler, decl)
	jsdocNode := NodeUtil_GetBestJSDocInfoNode(nameNode)
	jsdocNode.SetJSDocInfo(JsdocUtil_GetUnusableTypeJSDoc(jsdoc))
}

// Java: ConvertToTypedInterface.SimplifyDeclarations.elidePrivateMemberType(PotentialDeclaration)
func (s *ConvertToTypedInterface_SimplifyDeclarations) ElidePrivateMemberType(decl *PotentialDeclaration) {
	if s.ShouldElidePrivateMemberType(decl) {
		jsdoc := decl.GetJsDoc()
		jsdocNode := NodeUtil_GetBestJSDocInfoNode(decl.GetLhs())
		jsdocNode.SetJSDocInfo(JsdocUtil_ReplaceAnnotatedTypesWithUnusableType(jsdoc))
	}
}

// Java: ConvertToTypedInterface.SimplifyDeclarations.shouldElidePrivateMemberType(PotentialDeclaration)
func (s *ConvertToTypedInterface_SimplifyDeclarations) ShouldElidePrivateMemberType(decl *PotentialDeclaration) bool {
	// TODO: b/465855752 - roll out elison for JS files as well after depot cleanup.
	if !s.currentFile.IsFromTypeScript() {
		return false
	}
	jsdoc := decl.GetJsDoc()
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
	// a goog.module. (Ultimately we should probably ban marking a local name @private so that we
	// can drop this special case.)
	if decl.GetLhs().IsName() {
		return false
	}
	return JsdocUtil_HasAnnotatedType(jsdoc)
}

// Java: ConvertToTypedInterface.isSymbolProp(Node)
func ConvertToTypedInterface_IsSymbolProp(lhs rhino.Node) bool {
	prop, ok := ast.AsGetProp(lhs)
	return ok && prop.Object().Node().MatchesName("Symbol")
}

// Java: ConvertToTypedInterface.shouldPreserveAssignment(Node, NodeTraversal)
func ConvertToTypedInterface_ShouldPreserveAssignment(expr rhino.Node, t *NodeTraversal) bool {
	lhs := ast.MustAssignment(expr).Target().Node()
	// Ignore assignments in function bodies, unless they're also a constructor with a this. prop.
	if !t.InGlobalHoistScope() && !t.InModuleHoistScope() {
		return ClassUtil_IsThisPropInsideClassWithName(lhs)
	}

	// Well-known symbol properties, like Foo.prototype[Symbol.iterator] = function() {};
	if elem, ok := ast.AsGetElem(lhs); ok && ConvertToTypedInterface_IsSymbolProp(elem.Index().Node()) {
		return elem.Object().Node().IsQualifiedName()
	}
	// Assignments to names don't have global typechecking side-effects even within the 'hoist
	// scope'
	if lhs.IsName() {
		return t.InGlobalScope() || t.InModuleScope()
	}
	return lhs.IsQualifiedName()
}

// FileInfo is the class to keep track of what has been seen so far in a
// given file.
type FileInfo struct {
	providedNamespaces *jutil.LinkedHashSet[string]
	requiredLocalNames *jutil.LinkedHashSet[string]
	// MultimapBuilder.linkedHashKeys().arrayListValues().build()
	declarations     *jutil.LinkedHashMap[string, []*PotentialDeclaration]
	isFromTypeScript bool
}

// Java: FileInfo(String)
func NewFileInfo(sourceFileName string) *FileInfo {
	f := &FileInfo{
		providedNamespaces: jutil.NewLinkedHashSet[string](),
		requiredLocalNames: jutil.NewLinkedHashSet[string](),
		declarations:       jutil.NewLinkedHashMap[string, []*PotentialDeclaration](),
	}
	f.isFromTypeScript = strings.HasSuffix(sourceFileName, ".closure.js")
	return f
}

// Java: FileInfo.isFromTypeScript()
func (f *FileInfo) IsFromTypeScript() bool {
	return f.isFromTypeScript
}

// Java: FileInfo.recordNameDeclaration(Node)
func (f *FileInfo) RecordNameDeclaration(qualifiedNameNode rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromName(qualifiedNameNode))
}

// Java: FileInfo.recordMemberFieldDef(Node)
func (f *FileInfo) RecordMemberFieldDef(fieldNode rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromMemberFieldDef(fieldNode))
}

// Java: FileInfo.recordMethod(Node)
func (f *FileInfo) RecordMethod(functionNode rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromMethod(functionNode))
}

// Java: FileInfo.recordStringKeyDeclaration(Node)
func (f *FileInfo) RecordStringKeyDeclaration(stringKeyNode rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromStringKey(stringKeyNode))
}

// Java: FileInfo.recordDefine(Node)
func (f *FileInfo) RecordDefine(callNode rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromDefine(callNode))
}

// Java: FileInfo.recordAliasDeclaration(Node)
func (f *FileInfo) RecordAliasDeclaration(nameNode rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromAlias(nameNode))
}

// Java: FileInfo.getDeclarations()
func (f *FileInfo) GetDeclarations() *jutil.LinkedHashMap[string, []*PotentialDeclaration] {
	return f.declarations
}

// Java: FileInfo.recordDeclaration(PotentialDeclaration)
func (f *FileInfo) RecordDeclaration(decl *PotentialDeclaration) {
	name := decl.GetFullyQualifiedName()
	f.declarations.Put(name, append(f.declarations.Get(name), decl))
}

// Java: FileInfo.recordImport(String)
func (f *FileInfo) RecordImport(localName string) {
	f.requiredLocalNames.Add(localName)
}

// Java: FileInfo.isNameDeclared(String)
func (f *FileInfo) IsNameDeclared(fullyQualifiedName string) bool {
	return f.declarations.ContainsKey(fullyQualifiedName)
}

// Java: FileInfo.containsPrefix(String, Iterable<String>)
func FileInfo_ContainsPrefix(fullyQualifiedName string, prefixNamespaces *jutil.LinkedHashSet[string]) bool {
	for prefix := range prefixNamespaces.All() {
		if fullyQualifiedName == prefix || strings.HasPrefix(fullyQualifiedName, prefix+".") {
			return true
		}
	}
	return false
}

// Java: FileInfo.isPrefixProvided(String)
func (f *FileInfo) IsPrefixProvided(fullyQualifiedName string) bool {
	return FileInfo_ContainsPrefix(fullyQualifiedName, f.providedNamespaces)
}

// Java: FileInfo.isPrefixRequired(String)
func (f *FileInfo) IsPrefixRequired(fullyQualifiedName string) bool {
	return FileInfo_ContainsPrefix(fullyQualifiedName, f.requiredLocalNames)
}

// Java: FileInfo.isStrictPrefixDeclared(String)
func (f *FileInfo) IsStrictPrefixDeclared(fullyQualifiedName string) bool {
	for prefix := range f.declarations.Keys() {
		if strings.HasPrefix(fullyQualifiedName, prefix+".") {
			return true
		}
	}
	return false
}

// Java: FileInfo.markProvided(String)
func (f *FileInfo) MarkProvided(providedName string) {
	// checkNotNull(providedName): a Go string cannot be null.
	f.providedNamespaces.Add(providedName)
}
