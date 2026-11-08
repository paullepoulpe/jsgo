// Port of com/google/javascript/jscomp/Normalize.java @ closure-compiler v20261104

/*
 * Copyright 2008 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

type normalizeLocalDedupAction struct {
	nameNode    rhino.Node
	varParent   rhino.Node
	grandparent rhino.Node
}

type normalizeParallelScriptResult struct {
	dedups          []normalizeLocalDedupAction
	constNodes      []rhino.Node
	skippedDupNames []rhino.Node
}

var (
	propagateConstVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_NAME)
		return s
	}()

	normalizeStatementsVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_NAME)
		s.Add(rhino.Token_IMPORT_STAR)
		s.Add(rhino.Token_WHILE)
		s.Add(rhino.Token_DESTRUCTURING_LHS)
		s.Add(rhino.Token_ASSIGN_OR)
		s.Add(rhino.Token_ASSIGN_AND)
		s.Add(rhino.Token_ASSIGN_COALESCE)
		s.Add(rhino.Token_ARRAYLIT)
		s.Add(rhino.Token_CALL)
		s.Add(rhino.Token_NEW)
		s.Add(rhino.Token_OBJECTLIT)
		s.Add(rhino.Token_OPTCHAIN_CALL)
		s.Add(rhino.Token_PARAM_LIST)
		return s
	}()
)

type normalizeConstConvCache struct {
	alwaysFalse bool
	slots       [256]struct {
		id  uint32
		val bool
	}
}

func newNormalizeConstConvCache(conv CodingConvention) normalizeConstConvCache {
	_, isDefault := conv.(*CodingConventions_DefaultCodingConvention)
	return normalizeConstConvCache{alwaysFalse: isDefault}
}

func (c *normalizeConstConvCache) isConstant(conv CodingConvention, id uint32, n rhino.Node) bool {
	if c.alwaysFalse {
		return false
	}
	if id != 0 {
		slot := &c.slots[uint8(id^(id>>8))]
		if slot.id == id {
			return slot.val
		}
		val := conv.IsConstant(n.GetString())
		slot.id = id
		slot.val = val
		return val
	}
	return conv.IsConstant(n.GetString())
}

func (*Normalize_PropagateConstantPropertyOverVars) onlyNameLeaves()                {}
func (*Normalize_PropagateConstantPropertyOverVars) nodeTraversalNoShouldTraverse() {}
func (*Normalize_PropagateConstantPropertyOverVars) visitTokens() *ast.TokenSet {
	return &propagateConstVisitTokens
}
func (*ValidityCheck_VerifyConstants) onlyNameLeaves()                {}
func (*ValidityCheck_VerifyConstants) nodeTraversalNoShouldTraverse() {}

// Normalize_PropagateConstantPropertyOverVars propagates constant
// annotations and IS_CONSTANT_NAME property over the Var graph.
//
// Also invokes t.getScope() on every scope, for use with the
// DuplicateDeclarationHandler.
type Normalize_PropagateConstantPropertyOverVars struct {
	compiler       AbstractCompiler
	convention     CodingConvention
	assertOnChange bool
	varConst       map[*Var]bool
	convCache      normalizeConstConvCache
	parallelRes    *normalizeParallelScriptResult
}

// Java: Normalize.PropagateConstantPropertyOverVars(AbstractCompiler, boolean)
func NewNormalize_PropagateConstantPropertyOverVars(
	compiler AbstractCompiler, forbidChanges bool) *Normalize_PropagateConstantPropertyOverVars {
	conv := compiler.GetCodingConvention()
	p := &Normalize_PropagateConstantPropertyOverVars{
		compiler:       compiler,
		convention:     conv,
		assertOnChange: forbidChanges,
		convCache:      newNormalizeConstConvCache(conv),
	}
	return p
}

// Java: Normalize.PropagateConstantPropertyOverVars.enterScope(NodeTraversal)
func (p *Normalize_PropagateConstantPropertyOverVars) EnterScope(t *NodeTraversal) {
	// Cause the scope to be created, which will cause duplicate
	// to be found.
	if normalizeScopeMayHaveDeclarations(t.GetScopeRoot()) {
		t.GetScope()
	}
}

// Java: Normalize.PropagateConstantPropertyOverVars.exitScope(NodeTraversal)
func (p *Normalize_PropagateConstantPropertyOverVars) ExitScope(t *NodeTraversal) {
	// Nothing to do.
}

// Java: Normalize.PropagateConstantPropertyOverVars.shouldTraverse(NodeTraversal, Node, Node)
func (p *Normalize_PropagateConstantPropertyOverVars) ShouldTraverse(
	t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// Java: Normalize.PropagateConstantPropertyOverVars.visit(NodeTraversal, Node, Node)
func (p *Normalize_PropagateConstantPropertyOverVars) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Note: Constant properties annotations are not propagated.
	if !n.IsName() || n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
		return
	}
	if p.parallelRes != nil && len(p.parallelRes.skippedDupNames) > 0 && slices.Contains(p.parallelRes.skippedDupNames, n) {
		return
	}
	id := n.GetStringID()
	if id <= strpool.EmptyID {
		return
	}

	// Find the JSDocInfo for a top-level variable
	var_ := t.LookupVarByID(id)
	var shouldBeConstant bool
	if var_ != nil {
		var ok bool
		shouldBeConstant, ok = p.varConst[var_]
		if !ok {
			if var_.IsConst() || Normalize_IsFunctionExpressionNameNode(var_.GetNode()) {
				shouldBeConstant = true
			} else if info := var_.GetJSDocInfo(); info != nil && info.IsConstant() {
				shouldBeConstant = true
			} else if p.convCache.isConstant(p.convention, id, n) {
				shouldBeConstant = true
			}
			if p.varConst == nil {
				p.varConst = make(map[*Var]bool, 64)
			}
			p.varConst[var_] = shouldBeConstant
		}
	} else {
		shouldBeConstant = p.convCache.isConstant(p.convention, id, n)
	}
	if shouldBeConstant {
		if p.assertOnChange {
			panic(jutil.NewIllegalStateExceptionMsg(
				"Should be const but not marked as const.\n" +
					"  name: " +
					n.GetString() +
					"\n" +
					"  parent:" +
					n.GetParent().ToStringTree()))
		}
		if p.parallelRes != nil {
			p.parallelRes.constNodes = append(p.parallelRes.constNodes, n)
		} else {
			n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}
	}
}

func (*Normalize_NormalizeStatements) visitTokens() *ast.TokenSet {
	return &normalizeStatementsVisitTokens
}

// Normalize_NormalizeStatements simplifies the AST:
//   - VAR declarations split, so they represent exactly one child
//     declaration.
//   - WHILEs are converted to FORs.
//   - FOR loop are initializers are moved out of the FOR structure.
//   - LABEL node of children other than LABEL, BLOCK, WHILE, FOR, or DO are
//     moved into a block.
//   - Add constant annotations based on coding convention.
type Normalize_NormalizeStatements struct {
	compiler                                AbstractCompiler
	convention                              CodingConvention
	astFactory                              *AstFactory
	assertOnChange                          bool
	isNormalizedObfuscated                  bool
	rewriteLogicalAssignmentOperatorsHelper *RewriteLogicalAssignmentOperatorsHelper
	renamer                                 contextualRenamerState
	convCache                               normalizeConstConvCache
}

// NewNormalize_NormalizeStatements: prefer using Normalize as a compiler
// pass instead of using this class directly.
//
// Java: Normalize.NormalizeStatements(AbstractCompiler, boolean, MakeDeclaredNamesUnique)
func NewNormalize_NormalizeStatements(
	compiler AbstractCompiler,
	assertOnChange bool) *Normalize_NormalizeStatements {
	conv := compiler.GetCodingConvention()
	s := &Normalize_NormalizeStatements{
		compiler:               compiler,
		convention:             conv,
		assertOnChange:         assertOnChange,
		isNormalizedObfuscated: compiler.GetLifeCycleStage().IsNormalizedObfuscated(),
		astFactory:             compiler.CreateAstFactory(),
		convCache:              newNormalizeConstConvCache(conv),
	}
	s.rewriteLogicalAssignmentOperatorsHelper =
		NewRewriteLogicalAssignmentOperatorsHelper(
			compiler, s.astFactory, compiler.GetUniqueIdSupplier())
	return s
}

// normalizeScopeMayHaveDeclarations reports whether the scope rooted at n
// may declare names; EnterScope does not create the others.
func normalizeScopeMayHaveDeclarations(n rhino.Node) bool {
	if fn, ok := ast.AsFunction(n); ok {
		// Its parameters, and the name of a function expression.
		if fn.Params().HasParams() {
			return true
		}
		_, named := fn.Name()
		return named && NodeUtil_IsFunctionExpression(n)
	}
	if _, ok := ast.AsBlock(n); ok {
		if NodeUtil_IsFunctionBlock(n) || NodeUtil_IsClassStaticBlock(n) {
			return !syntacticFunctionBlockHasNoDecls(n)
		}
		// (The children are statements, or the catch of a try.)
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			switch c.GetToken() {
			case rhino.Token_LET, rhino.Token_CONST, rhino.Token_FUNCTION, rhino.Token_CLASS,
				rhino.Token_CATCH, rhino.Token_EXPORT, rhino.Token_LABEL, rhino.Token_BLOCK:
				return true
			}
		}
		return false
	}
	if cls, ok := ast.AsClass(n); ok {
		// The name of a class expression.
		_, named := cls.Name()
		return named && NodeUtil_IsClassExpression(n)
	}
	if loop, ok := ast.AsEnhancedFor(n); ok {
		return normalizeMayDeclareInLoopScope(loop.Target().Node())
	}
	if f, ok := ast.AsFor(n); ok {
		init, ok := f.Init()
		return ok && normalizeMayDeclareInLoopScope(init.Node())
	}
	switch n.GetToken() {
	case rhino.Token_COMPUTED_FIELD_DEF, rhino.Token_MEMBER_FIELD_DEF:
		return false
	}
	return true
}

// normalizeMayDeclareInLoopScope reports whether n, the init of a for or the
// target of a for-in/of, may declare names in the scope of the loop.
func normalizeMayDeclareInLoopScope(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_FUNCTION:
		return true
	}
	return false
}

// Java: Normalize.NormalizeStatements.reportCodeChange(String, Node)
func (s *Normalize_NormalizeStatements) ReportCodeChange(changeDescription string, n rhino.Node) {
	if s.assertOnChange {
		panic(jutil.NewIllegalStateExceptionMsg("Normalize constraints violated:\n" + changeDescription))
	}
	s.compiler.ReportChangeToEnclosingScope(n)
}

// Java: Normalize.NormalizeStatements.enterScope(NodeTraversal)
func (s *Normalize_NormalizeStatements) EnterScope(t *NodeTraversal) {
	s.renamer.enterScope(t)
}

// Java: Normalize.NormalizeStatements.exitScope(NodeTraversal)
func (s *Normalize_NormalizeStatements) ExitScope(t *NodeTraversal) {
	s.renamer.exitScope(t)
}

// Java: Normalize.NormalizeStatements.shouldTraverse(NodeTraversal, Node, Node)
func (s *Normalize_NormalizeStatements) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if !normalizePreOrderTokens[uint8(n.GetToken())] || n.GetFirstChild().IsNil() {
		return true
	}

	oldParent := n.GetParent()
	s.DoStatementNormalizations(t, n)
	jutil.CheckStateMsg(
		n.GetParent() == oldParent,
		"n's parent changed during normalization. It needs to stay the same so that we properly"+
			" visit any added nodes.")
	return true
}

// Visit runs post-order normalizations for n's token.
//
// Java: Normalize.NormalizeStatements.visit(NodeTraversal, Node, Node)
func (s *Normalize_NormalizeStatements) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		s.visitName(t, ast.MustName(n), parent)
	case rhino.Token_IMPORT_STAR:
		s.visitImportStar(t, ast.MustImportStar(n))
	case rhino.Token_WHILE:
		s.visitWhile(ast.MustWhile(n))
	case rhino.Token_DESTRUCTURING_LHS:
		s.NormalizeDestructuringLhs(ast.MustDestructuringLhs(n))
	case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_COALESCE:
		s.rewriteLogicalAssignmentOperatorsHelper.VisitLogicalAssignmentOperator(t, n)
	case rhino.Token_ARRAYLIT, rhino.Token_CALL, rhino.Token_NEW,
		rhino.Token_OBJECTLIT, rhino.Token_OPTCHAIN_CALL, rhino.Token_PARAM_LIST:
		n.SetTrailingComma(false)
	}
}

// visitName renames x if it is shadowed and marks it constant when it is.
//
// Java: Normalize.NormalizeStatements.visit(NodeTraversal, Node, Node) (NAME)
func (s *Normalize_NormalizeStatements) visitName(t *NodeTraversal, x ast.Name, parent rhino.Node) {
	n := x.Node()
	s.renamer.renameNode(t, n, s.assertOnChange)
	if n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
		return
	}
	if !s.isNormalizedObfuscated && s.convCache.isConstant(s.convention, n.GetStringID(), n) {
		jutil.CheckStateMsg1(!s.assertOnChange, "Not marked as constant when it should be: %s", n)
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		return
	}
	if parent.IsNil() {
		parent = n.GetParent()
	}
	if parent.IsFunction() && n == parent.GetFirstChild() && NodeUtil_IsFunctionExpression(parent) {
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
}

// visitImportStar renames x if it is shadowed.
func (s *Normalize_NormalizeStatements) visitImportStar(t *NodeTraversal, x ast.ImportStar) {
	s.renamer.renameNode(t, x.Node(), s.assertOnChange)
}

// visitWhile turns while (cond) body into for (; cond; ) body.
//
// Java: Normalize.NormalizeStatements.visit(NodeTraversal, Node, Node) (WHILE)
func (s *Normalize_NormalizeStatements) visitWhile(x ast.While) {
	n := x.Node()
	cond := x.Cond().Node()
	n.SetToken(rhino.Token_FOR)
	empty := rhino.IR_Empty()
	empty.SrcrefIfMissing(n)
	empty.InsertBefore(cond)
	empty.CloneNode().InsertAfter(cond)
	s.ReportCodeChange("WHILE node", n)
}

// AnnotateConstantsByConvention marks names and properties that are
// constants by convention.
//
// Java: Normalize.NormalizeStatements.annotateConstantsByConvention(Node)
func (s *Normalize_NormalizeStatements) AnnotateConstantsByConvention(x ast.Name) {
	// Need to check that variables have not been renamed, to determine whether
	// coding conventions still apply.
	if s.isNormalizedObfuscated {
		return
	}

	n := x.Node()
	if !n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) && s.convCache.isConstant(s.convention, n.GetStringID(), n) {
		jutil.CheckStateMsg1(!s.assertOnChange, "Not marked as constant when it should be: %s", n)
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
}

// AnnotateFunctionExpressionNameAsConstant annotates function names on
// function expressions to be unconditionally constant
//
// Java: Normalize.NormalizeStatements.annotateFunctionExpressionNameAsConstant(Node)
func (s *Normalize_NormalizeStatements) AnnotateFunctionExpressionNameAsConstant(x ast.Name) {
	if normalizeIsFunctionExpressionName(x) {
		x.Node().PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
}

// Normalize_DuplicateDeclarationHandler is the ScopeCreator duplicate
// declaration handler.
type Normalize_DuplicateDeclarationHandler struct {
	outer       *Normalize
	parallelRes *normalizeParallelScriptResult

	hasOkDuplicateDeclaration map[*Var]struct{}
	okDupBuf                  [8]*Var
	okDupLen                  int
}

// Java: Normalize.DuplicateDeclarationHandler()
func NewNormalize_DuplicateDeclarationHandler(outer *Normalize) *Normalize_DuplicateDeclarationHandler {
	return &Normalize_DuplicateDeclarationHandler{
		outer: outer,
	}
}

// OnRedeclaration removes duplicate VAR declarations discovered during scope
// creation.
//
// Java: Normalize.DuplicateDeclarationHandler.onRedeclaration(Scope, String, Node, CompilerInput)
func (h *Normalize_DuplicateDeclarationHandler) OnRedeclaration(
	s *Scope, name string, n rhino.Node, input *CompilerInput) {
	jutil.CheckState(n.IsName())
	parent := n.GetParent()
	v := s.GetVarByNode(n)

	if s.IsGlobal() {
		// We allow variables to be duplicate declared if one
		// declaration appears in source and the other in externs.
		// This deals with issues where a browser built-in is declared
		// in one browser but not in another.
		if v.IsExtern() && !input.IsExtern() {
			for i := 0; i < h.okDupLen; i++ {
				if h.okDupBuf[i] == v {
					goto duplicateNotOk
				}
			}
			if h.okDupLen < len(h.okDupBuf) {
				h.okDupBuf[h.okDupLen] = v
				h.okDupLen++
				return
			}
			if h.hasOkDuplicateDeclaration == nil {
				h.hasOkDuplicateDeclaration = make(map[*Var]struct{}, 32)
			}
			if _, exists := h.hasOkDuplicateDeclaration[v]; !exists {
				h.hasOkDuplicateDeclaration[v] = struct{}{}
				return
			}
		}
	}
duplicateNotOk:

	if parent.IsFunction() {
		if v.GetParentNode().IsVar() {
			s.Undeclare(v)
			s.Declare(name, n, v.GetInput())
			h.recordOrReplaceVar(
				v.GetNameNode(), v.GetParentNode(), v.GetParentNode().GetParent())
		}
	} else if parent.IsVar() {
		jutil.CheckState(parent.HasOneChild())
		h.recordOrReplaceVar(n, parent, parent.GetParent())
	}
}

func (h *Normalize_DuplicateDeclarationHandler) recordOrReplaceVar(
	n rhino.Node, parent rhino.Node, grandparent rhino.Node) {
	if res := h.parallelRes; res != nil {
		name, ok := ast.AsName(n)
		jutil.CheckState(ok)
		if _, hasInit := name.Init(); !hasInit {
			jutil.CheckStateMsg1(NodeUtil_IsStatementBlock(grandparent), "%s", grandparent)
			res.skippedDupNames = append(res.skippedDupNames, n)
		}
		res.dedups = append(res.dedups, normalizeLocalDedupAction{
			nameNode:    n,
			varParent:   parent,
			grandparent: grandparent,
		})
		return
	}
	h.ReplaceVarWithAssignment(n, parent, grandparent)
}

// ReplaceVarWithAssignment removes the parent VAR. There is only one case
// that need to be handled: "var a = b;" which is replaced with "a = b"
//
// Cases we don't need to handle are VARs with multiple children, which have
// already been split into separate declarations, so there is no need to
// handle that here; "for (var a;;);"/"for (var a of b)"/"for (var a in b)",
// which have been moved out of the loop; and "LABEL: var x;" which has been
// put in a BLOCK
//
// The result of this is that in each case the parent node is replaced which
// is generally dangerous in a traversal but is fine here with the scope
// creator, as the next node of interest is the parent's next sibling.
//
// Java: Normalize.DuplicateDeclarationHandler.replaceVarWithAssignment(Node, Node, Node)
func (h *Normalize_DuplicateDeclarationHandler) ReplaceVarWithAssignment(
	n rhino.Node, parent rhino.Node, grandparent rhino.Node) {
	name, ok := ast.AsName(n)
	jutil.CheckState(ok)
	if value, ok := name.Init(); ok {
		// The  *  is being initialize, preserve the new value.
		n.Detach()
		// Convert "var name = value" to "name = value"
		value.Node().Detach()
		replacement := h.outer.astFactory.CreateAssign(n, value.Node())
		replacement.SetJSDocInfo(parent.GetJSDocInfo())
		replacement.SrcrefIfMissing(parent)
		var statement rhino.Node
		if parent.CanResetInPlace() && replacement.CanResetInPlace() {
			statement = parent.MutateAttachedToExprResultInPlace(replacement)
		} else {
			statement = NodeUtil_NewExpr(replacement)
			parent.ReplaceWith(statement)
		}
		h.outer.ReportCodeChange("Duplicate VAR declaration", statement)
	} else {
		// It is an empty reference. Remove it.
		jutil.CheckStateMsg1(NodeUtil_IsStatementBlock(grandparent), "%s", grandparent)
		parent.Detach()
		h.outer.ReportCodeChange("Duplicate VAR declaration", grandparent)
	}
}
