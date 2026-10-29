// Port of com/google/javascript/jscomp/VariableReferenceCheck.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// VariableReferenceCheck checks variables to see if they are referenced before
// their declaration, or if they are redeclared in a way that is suspicious
// (i.e. not dictated by control structures). This is a more aggressive
// version of VarCheck, but it lacks the cross-chunk checks.
type VariableReferenceCheck struct {
	compiler AbstractCompiler

	checkUnusedLocals bool

	// NOTE(nicksantos): It's a lot faster to use a shared Set that
	// we clear after each method call, because the Set never gets too big.
	blocksWithDeclarations []*BasicBlock
}

// Java: VariableReferenceCheck.EARLY_REFERENCE
var VariableReferenceCheck_EARLY_REFERENCE = DiagnosticType_Warning(
	"JSC_REFERENCE_BEFORE_DECLARE", "Variable referenced before declaration: {0}")

// Java: VariableReferenceCheck.EARLY_EXPORTS_REFERENCE
var VariableReferenceCheck_EARLY_EXPORTS_REFERENCE = DiagnosticType_Error(
	"JSC_EXPORTS_REFERENCE_BEFORE_ASSIGN",
	"Illegal reference to `exports` before assignment `exports = ...`")

// Java: VariableReferenceCheck.REDECLARED_VARIABLE
var VariableReferenceCheck_REDECLARED_VARIABLE = DiagnosticType_Warning("JSC_REDECLARED_VARIABLE", "Redeclared variable: {0}")

// Java: VariableReferenceCheck.EARLY_REFERENCE_ERROR
var VariableReferenceCheck_EARLY_REFERENCE_ERROR = DiagnosticType_Error(
	"JSC_REFERENCE_BEFORE_DECLARE_ERROR",
	"Illegal variable reference before declaration: {0}")

// Java: VariableReferenceCheck.REASSIGNED_CONSTANT
var VariableReferenceCheck_REASSIGNED_CONSTANT = DiagnosticType_Error("JSC_REASSIGNED_CONSTANT", "Constant reassigned: {0}")

// Java: VariableReferenceCheck.REDECLARED_VARIABLE_ERROR
var VariableReferenceCheck_REDECLARED_VARIABLE_ERROR = DiagnosticType_Error("JSC_REDECLARED_VARIABLE_ERROR", "Illegal redeclared variable: {0}")

// Java: VariableReferenceCheck.DECLARATION_NOT_DIRECTLY_IN_BLOCK
var VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK = DiagnosticType_Error(
	"JSC_DECLARATION_NOT_DIRECTLY_IN_BLOCK",
	"Block-scoped declaration not directly within block: {0}")

// Java: VariableReferenceCheck.UNUSED_LOCAL_ASSIGNMENT
var VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT = DiagnosticType_Disabled(
	"JSC_UNUSED_LOCAL_ASSIGNMENT", "Value assigned to local variable {0} is never read")

// Java: VariableReferenceCheck.GOOG_REQUIRE
var variableReferenceCheck_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")

// Java: VariableReferenceCheck.GOOG_REQUIRE_TYPE
var variableReferenceCheck_GOOG_REQUIRE_TYPE = rhino.QualifiedName_Of("goog.requireType")

// Java: VariableReferenceCheck.GOOG_FORWARD_DECLARE
var variableReferenceCheck_GOOG_FORWARD_DECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")

// These types do not permit a block-scoped declaration inside them without an
// explicit block. e.g. if (b) let x;
// This list omits Token.LABEL intentionally. It's handled differently in
// IRFactory.
//
// Java: VariableReferenceCheck.BLOCKLESS_DECLARATION_FORBIDDEN_STATEMENTS
func variableReferenceCheck_isBlocklessDeclarationForbidden(t rhino.Token) bool {
	switch t {
	case rhino.Token_IF, rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_WHILE:
		return true
	default:
		return false
	}
}

// Java: VariableReferenceCheck.GOOG_SCOPE
var variableReferenceCheck_GOOG_SCOPE = rhino.QualifiedName_Of("goog.scope")

// Java: VariableReferenceCheck(AbstractCompiler)
func NewVariableReferenceCheck(compiler AbstractCompiler) *VariableReferenceCheck {
	return &VariableReferenceCheck{
		compiler: compiler,
		checkUnusedLocals: compiler.GetOptions().
			Enables(DiagnosticGroup_ForType(VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT)),
	}
}

func (v *VariableReferenceCheck) addBlockWithDeclaration(block *BasicBlock) {
	if !slices.Contains(v.blocksWithDeclarations, block) {
		v.blocksWithDeclarations = append(v.blocksWithDeclarations, block)
	}
}

// Java: VariableReferenceCheck.process(Node, Node)
func (v *VariableReferenceCheck) Process(externs rhino.Node, root rhino.Node) {
	sc := NewSyntacticScopeCreator(v.compiler)
	callback := NewReferenceCollector(
		v.compiler, NewVariableReferenceCheck_ReferenceCheckingBehavior(v), sc)
	callback.enableLocalReclamation(nil)
	callback.Process(externs, root)
	clear(v.blocksWithDeclarations)
	v.blocksWithDeclarations = v.blocksWithDeclarations[:0]
}

// VariableReferenceCheck_ReferenceCheckingBehavior is behavior that checks
// variables for redeclaration or early references just after they go out of
// scope.
type VariableReferenceCheck_ReferenceCheckingBehavior struct {
	v *VariableReferenceCheck

	varsInFunctionBody map[string]struct{}
}

// Java: VariableReferenceCheck.ReferenceCheckingBehavior()
func NewVariableReferenceCheck_ReferenceCheckingBehavior(v *VariableReferenceCheck) *VariableReferenceCheck_ReferenceCheckingBehavior {
	return &VariableReferenceCheck_ReferenceCheckingBehavior{v: v}
}

// Java: VariableReferenceCheck.ReferenceCheckingBehavior.afterExitScope(NodeTraversal, ReferenceMap)
func (b *VariableReferenceCheck_ReferenceCheckingBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	if len(t.scopes) > 0 && t.scopes[len(t.scopes)-1].scope == nil {
		if NodeUtil_IsFunctionBlock(t.GetScopeRoot()) {
			clear(b.varsInFunctionBody)
		}
		return
	}
	// Check all vars after finishing a scope
	scope := t.GetScope()
	if scope.IsFunctionBlockScope() {
		clear(b.varsInFunctionBody)
		for _, v := range scope.varList {
			if v != nil {
				if b.varsInFunctionBody == nil {
					b.varsInFunctionBody = make(map[string]struct{}, len(scope.varList))
				}
				b.varsInFunctionBody[v.GetName()] = struct{}{}
			}
		}
	}
	for _, v := range scope.varList {
		if v == nil {
			continue
		}
		referenceCollection := referenceMap.GetReferences(v)
		// TODO(moz): Figure out why this could be null
		if referenceCollection != nil {
			if scope.GetRootNode().IsFunction() && v.IsDefaultParam() {
				b.CheckDefaultParam(v, scope, b.varsInFunctionBody)
			}
			if scope.GetRootNode().IsFunction() {
				b.CheckShadowParam(v, scope, referenceCollection.references)
			}
			b.CheckVar(v, referenceCollection.references)
		}
	}
	if scope.hasOwnImplicitVar(AbstractScope_ImplicitVar_EXPORTS) {
		b.CheckGoogModuleExports(scope.MakeImplicitVar(AbstractScope_ImplicitVar_EXPORTS), referenceMap)
	}
}

// Java: VariableReferenceCheck.ReferenceCheckingBehavior.checkDefaultParam(Var, Scope, Set<String>)
func (b *VariableReferenceCheck_ReferenceCheckingBehavior) CheckDefaultParam(
	param *Var, scope *Scope, varsInFunctionBody map[string]struct{}) {
	NodeTraversal_Traverse(
		b.v.compiler,
		param.GetParentNode().GetSecondChild(),
		/*
		 * Do a shallow check since cases like: {@code
		 *   function f(y = () => x, x = 5) { return y(); }
		 * } is legal. We are going to miss cases like: {@code
		 *   function f(y = (() => x)(), x = 5) { return y(); }
		 * } but this should be rare.
		 */
		&variableReferenceCheck_defaultParamCallback{b: b, scope: scope, varsInFunctionBody: varsInFunctionBody})
}

var variableReferenceCheckDefaultParamVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_NAME)
	return s
}()

// variableReferenceCheck_defaultParamCallback is the anonymous
// AbstractShallowCallback of checkDefaultParam.
type variableReferenceCheck_defaultParamCallback struct {
	NodeTraversal_AbstractShallowCallbackBase

	b                  *VariableReferenceCheck_ReferenceCheckingBehavior
	scope              *Scope
	varsInFunctionBody map[string]struct{}
}

var _ nodeTraversal_visitTokens = (*variableReferenceCheck_defaultParamCallback)(nil)

func (*variableReferenceCheck_defaultParamCallback) visitTokens() *ast.TokenSet {
	return &variableReferenceCheckDefaultParamVisitTokens
}

// Java: VariableReferenceCheck.ReferenceCheckingBehavior.checkDefaultParam(...).new AbstractShallowCallback() {...}.visit(NodeTraversal, Node, Node)
func (cb *variableReferenceCheck_defaultParamCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !NodeUtil_IsReferenceName(n) {
		return
	}
	refName := n.GetString()
	if _, ok := cb.varsInFunctionBody[refName]; ok && !cb.scope.HasSlot(refName) {
		cb.b.v.compiler.Report(JSError_MakeWithNode(n, VariableReferenceCheck_EARLY_REFERENCE_ERROR, refName))
	}
}

// Java: VariableReferenceCheck.ReferenceCheckingBehavior.checkShadowParam(Var, Scope, List<Reference>)
func (b *VariableReferenceCheck_ReferenceCheckingBehavior) CheckShadowParam(
	v *Var, functionScope *Scope, references []*Reference) {
	maybeParam := functionScope.GetVarByID(v.nameID)
	if maybeParam != nil && maybeParam.IsParam() && maybeParam.scope == functionScope {
		for _, r := range references {
			if (r.IsVarDeclaration() || r.IsHoistedFunction()) &&
				r.GetNode() != v.GetNameNode() {
				b.v.compiler.Report(JSError_MakeWithNode(r.GetNode(), VariableReferenceCheck_REDECLARED_VARIABLE, v.GetName()))
			}
		}
	}
}

// Java: VariableReferenceCheck.ReferenceCheckingBehavior.checkGoogModuleExports(Var, ReferenceMap)
func (b *VariableReferenceCheck_ReferenceCheckingBehavior) CheckGoogModuleExports(exportsVar *Var, referenceMap ReferenceMap) {
	references := referenceMap.GetReferences(exportsVar)
	if references == nil || references.IsNeverAssigned() {
		return
	}

	for _, reference := range references.references {
		if reference.IsLvalue() {
			break
		}
		b.v.CheckEarlyReference(exportsVar, reference, reference.GetNode())
	}
}

// CheckVar: if the variable is declared more than once in a basic block,
// generate a warning. Also check if a variable is used in a given scope
// before it is declared, which suggest a likely error. Relies on the fact
// that references is in parse-tree order.
//
// Java: VariableReferenceCheck.ReferenceCheckingBehavior.checkVar(Var, List<Reference>)
func (b *VariableReferenceCheck_ReferenceCheckingBehavior) CheckVar(v *Var, references []*Reference) {
	c := b.v
	clear(c.blocksWithDeclarations)
	c.blocksWithDeclarations = c.blocksWithDeclarations[:0]
	hasSeenDeclaration := false
	hasErrors := false
	isRead := false
	var unusedAssignment *Reference

	hoistedFn := c.LookForHoistedFunction(references)
	if hoistedFn != nil {
		hasSeenDeclaration = true
	}

	for _, reference := range references {
		if reference == hoistedFn {
			continue
		}

		referenceNode := reference.GetNode()
		basicBlock := reference.GetBasicBlock()
		isDeclaration := reference.IsDeclaration()
		isAssignment := isDeclaration || reference.IsLvalue()

		if isDeclaration {
			// Checks for declarations
			hasSeenDeclaration = true
			hasErrors = c.CheckRedeclaration(v, reference, referenceNode, hoistedFn, basicBlock)
			// Add the current basic block after checking redeclarations
			c.addBlockWithDeclaration(basicBlock)
			c.CheckBlocklessDeclaration(v, reference, referenceNode)

			if reference.GetGrandparent().IsExport() {
				isRead = true
			}
		} else {
			// Checks for references
			if !hasSeenDeclaration {
				hasErrors = c.CheckEarlyReference(v, reference, referenceNode)
			}

			if !hasErrors && v.IsConst() && reference.IsLvalue() {
				c.compiler.Report(JSError_MakeWithNode(referenceNode, VariableReferenceCheck_REASSIGNED_CONSTANT, v.GetName()))
			}

			// Check for temporal dead zone of let / const declarations in for-in and for-of loops
			// TODO(b/111441110): Fix this check. it causes spurious warnings on `b = a` in
			//   for (const [a, b = a] of []) {}
			if (v.IsLet() || v.IsConst()) &&
				v.scope == reference.GetScope() &&
				NodeUtil_IsEnhancedFor(reference.GetScope().GetRootNode()) {
				c.compiler.Report(JSError_MakeWithNode(referenceNode, VariableReferenceCheck_EARLY_REFERENCE_ERROR, v.GetName()))
			}
		}

		if isAssignment {
			decl := references[0]
			declNode := decl.GetNode()
			gp := declNode.GetGrandparent()
			lhsOfForInLoop := gp.IsForIn() && gp.GetFirstFirstChild() == declNode

			if decl.GetScope().IsLocal() &&
				(decl.IsVarDeclaration() || decl.IsLetDeclaration() || decl.IsConstDeclaration()) &&
				!decl.GetNode().IsFromExterns() &&
				!lhsOfForInLoop {
				unusedAssignment = reference
			}
			if (reference.GetParent().IsDec() || reference.GetParent().IsInc() ||
				NodeUtil_IsCompoundAssignmentOp(reference.GetParent())) &&
				NodeUtil_IsExpressionResultUsed(reference.GetNode()) {
				isRead = true
			}
		} else {
			isRead = true
		}
	}

	if c.checkUnusedLocals && unusedAssignment != nil && !isRead && !hasErrors {
		c.CheckForUnusedLocalVar(v, unusedAssignment)
	}
}

// LookForHoistedFunction returns the reference to the hoisted function, if
// the variable is one.
//
// Java: VariableReferenceCheck.lookForHoistedFunction(List<Reference>)
func (v *VariableReferenceCheck) LookForHoistedFunction(references []*Reference) *Reference {
	for _, reference := range references {
		if reference.IsHoistedFunction() {
			v.addBlockWithDeclaration(reference.GetBasicBlock())
			return reference
		}
	}
	return nil
}

// Java: VariableReferenceCheck.checkBlocklessDeclaration(Var, Reference, Node)
func (v *VariableReferenceCheck) CheckBlocklessDeclaration(va *Var, reference *Reference, referenceNode rhino.Node) {
	if !reference.IsVarDeclaration() && reference.GetGrandparent().IsAddedBlock() &&
		variableReferenceCheck_isBlocklessDeclarationForbidden(
			reference.GetGrandparent().GetParent().GetToken()) {
		v.compiler.Report(JSError_MakeWithNode(referenceNode, VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK, va.GetName()))
	}
}

// CheckRedeclaration returns if a redeclaration error has been found.
//
// Java: VariableReferenceCheck.checkRedeclaration(Var, Reference, Node, Reference, BasicBlock)
func (vc *VariableReferenceCheck) CheckRedeclaration(
	v *Var, reference *Reference, referenceNode rhino.Node, hoistedFn *Reference, basicBlock *BasicBlock) bool {
	letConstShadowsVar := v.GetParentNode().IsVar() &&
		(reference.IsLetDeclaration() || reference.IsConstDeclaration())
	isVarNodeSameAsReferenceNode := v.GetNode() == reference.GetNode()
	// We disallow redeclaration of caught exceptions
	shadowCatchVar := v.GetParentNode().IsCatch() && !isVarNodeSameAsReferenceNode

	if vc.IsRedeclaration(basicBlock) &&
		!VarCheck_HasDuplicateDeclarationSuppression(vc.compiler, referenceNode, v.GetNameNode()) {
		var diagnosticType *DiagnosticType
		warningNode := referenceNode
		shadowParam := v.IsParam() &&
			// TODOO - lharker: this won't handle destructuring block scoped declarations. Is this
			// an actual bug?
			NodeUtil_IsBlockScopedDeclaration(referenceNode.GetParent()) &&
			v.scope == reference.GetScope().parent

		isFunctionDecl := !v.GetParentNode().IsNil() &&
			v.GetParentNode().IsFunction() &&
			v.GetParentNode().GetFirstChild() == referenceNode

		if v.IsLet() ||
			v.IsConst() ||
			v.IsClass() ||
			letConstShadowsVar ||
			shadowCatchVar ||
			shadowParam ||
			v.IsImport() ||
			isFunctionDecl {
			// These cases are all hard errors that violate ES6 semantics
			diagnosticType = VariableReferenceCheck_REDECLARED_VARIABLE_ERROR
		} else if reference.GetNode().GetParent().IsCatch() {
			return false
		} else {
			// These diagnostics are for valid, but suspicious, code, and are suppressible.
			// For vars defined in the global scope, give the same error as VarCheck
			if v.scope.IsGlobal() {
				diagnosticType = VarCheck_VAR_MULTIPLY_DECLARED_ERROR
			} else {
				diagnosticType = VariableReferenceCheck_REDECLARED_VARIABLE
			}
			// Since we skip hoisted functions, we would have the wrong warning node in cases
			// where the redeclaration is a function declaration. Check for that case.
			if isVarNodeSameAsReferenceNode &&
				hoistedFn != nil &&
				v.GetName() == hoistedFn.GetNode().GetString() {
				warningNode = hoistedFn.GetNode()
			}
		}
		vc.compiler.Report(
			JSError_MakeWithNode(warningNode, diagnosticType, v.GetName(), VariableReferenceCheck_LocationOf(v.GetNode())))
		return true
	}

	if (letConstShadowsVar || shadowCatchVar) && v.scope == reference.GetScope() {
		vc.compiler.Report(JSError_MakeWithNode(referenceNode, VariableReferenceCheck_REDECLARED_VARIABLE_ERROR, v.GetName()))
		return true
	}

	return false
}

// IsRedeclaration returns whether the given block executes after any known
// declarations of the variable being visited.
//
// Java: VariableReferenceCheck.isRedeclaration(BasicBlock)
func (v *VariableReferenceCheck) IsRedeclaration(newDeclaration *BasicBlock) bool {
	for _, previousDeclaration := range v.blocksWithDeclarations {
		if previousDeclaration.ProvablyExecutesBefore(newDeclaration) {
			return true
		}
	}
	return false
}

// Java: VariableReferenceCheck.locationOf(Node)
func VariableReferenceCheck_LocationOf(n rhino.Node) string {
	if n.IsNil() {
		return "<unknown>"
	}
	return n.GetLocation()
}

// CheckEarlyReference returns if an early reference has been found.
//
// Java: VariableReferenceCheck.checkEarlyReference(Var, Reference, Node)
func (vc *VariableReferenceCheck) CheckEarlyReference(v *Var, reference *Reference, referenceNode rhino.Node) bool {
	// Don't check the order of references in externs files.
	if referenceNode.IsFromExterns() || v.IsImplicitGoogNamespace() {
		return false
	}
	// Special case to deal with var goog = goog || {}. Note that
	// let x = x || {} is illegal, just like var y = x || {}; let x = y;
	if v.IsVar() {
		curr := reference.GetParent()
		for curr.IsOr() && curr.GetParent().GetFirstChild() == curr {
			curr = curr.GetParent()
		}
		if curr.IsName() && curr.GetStringID() == v.nameID {
			return false
		}
	}

	// RHS of public fields are not early references
	referenceScopeRoot := reference.GetScope().GetRootNode()
	if referenceScopeRoot.IsMemberFieldDef() && !referenceScopeRoot.IsStaticMember() {
		return false
	}

	// Only generate warnings for early references in the same function scope/global scope in
	// order to deal with possible forward declarations and recursion
	// e.g. don't warn on:
	//   function f() { return x; } f(); let x = 5;
	// We don't track where `f` is called, just where it's defined, and don't want to warn for
	//     function f() { return x; } let x = 5; f();
	// TODO(moz): See if we can remove the bypass for "goog"
	if reference.GetScope().HasSameContainerScope(v.scope) && v.GetName() != "goog" {
		var diagnostic *DiagnosticType
		if v.IsGoogModuleExports() {
			diagnostic = VariableReferenceCheck_EARLY_EXPORTS_REFERENCE
		} else if v.IsLet() || v.IsConst() || v.IsClass() || v.IsParam() {
			diagnostic = VariableReferenceCheck_EARLY_REFERENCE_ERROR
		} else {
			diagnostic = VariableReferenceCheck_EARLY_REFERENCE
		}
		vc.compiler.Report(JSError_MakeWithNode(reference.GetNode(), diagnostic, v.GetName()))
		return true
	}

	return false
}

// CheckForUnusedLocalVar: only check for unused local if not in a goog.scope
// function.
// TODO(tbreisacher): Consider moving UNUSED_LOCAL_ASSIGNMENT into its own
// check pass, so that we can run it after goog.scope processing, and get rid
// of the inGoogScope check.
//
// Java: VariableReferenceCheck.checkForUnusedLocalVar(Var, Reference)
func (vc *VariableReferenceCheck) CheckForUnusedLocalVar(v *Var, unusedAssignment *Reference) {
	if !v.IsLocal() {
		return
	}
	jsDoc := NodeUtil_GetBestJSDocInfo(unusedAssignment.GetNode())
	if jsDoc != nil && jsDoc.HasTypedefType() {
		return
	}

	inGoogScope := false
	s := v.scope
	if s.IsFunctionBlockScope() {
		function := s.GetRootNode().GetParent()
		callee := function.GetPrevious()
		inGoogScope = !callee.IsNil() && variableReferenceCheck_GOOG_SCOPE.Matches(callee)
	}

	if inGoogScope {
		// No warning.
		return
	}

	if s.IsModuleScope() {
		statement := NodeUtil_GetEnclosingStatement(v.GetNode())
		if NodeUtil_IsNameDeclaration(statement) {
			lhs := statement.GetFirstChild()
			rhs := lhs.GetFirstChild()
			if !rhs.IsNil() &&
				(NodeUtil_IsCallToQualifiedName(rhs, variableReferenceCheck_GOOG_FORWARD_DECLARE) ||
					NodeUtil_IsCallToQualifiedName(rhs, variableReferenceCheck_GOOG_REQUIRE_TYPE) ||
					NodeUtil_IsCallToQualifiedName(rhs, variableReferenceCheck_GOOG_REQUIRE) ||
					rhs.IsQualifiedName()) {
				// No warning. module imports will be caught by the unused-require check, and if the
				// right side is a qualified name then this is likely an alias used in type annotations.
				return
			}
		}
	}

	vc.compiler.Report(JSError_MakeWithNode(unusedAssignment.GetNode(), VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT, v.GetName()))
}
