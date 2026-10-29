// Port of com/google/javascript/jscomp/InlineVariables.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	ivaImmutableKnown uint8 = 1 << iota
	ivaImmutableVal
	ivaAliasedVarKnown
)

// inlineVariables_initialValueAnalysis is information about a value used to
// initialize a variable.
type inlineVariables_initialValueAnalysis struct {
	outer *inlineVariables_standardVarExpert

	aliasedVar *Var
	value      rhino.Node
	flags      uint8
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.InitialValueAnalysis.isImmutableValue()
func (a *inlineVariables_initialValueAnalysis) IsImmutableValue() bool {
	if a.flags&ivaImmutableKnown != 0 {
		return a.flags&ivaImmutableVal != 0
	}
	imm := !a.value.IsNil() && NodeUtil_IsImmutableValue(a.value)
	a.flags |= ivaImmutableKnown
	if imm {
		a.flags |= ivaImmutableVal
	}
	return imm
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.InitialValueAnalysis.isAlias()
func (a *inlineVariables_initialValueAnalysis) IsAlias() bool {
	return a.GetAliasedVar() != nil
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.InitialValueAnalysis.getAliasedVar()
func (a *inlineVariables_initialValueAnalysis) GetAliasedVar() *Var {
	if a.flags&ivaAliasedVarKnown != 0 {
		return a.aliasedVar
	}
	a.flags |= ivaAliasedVarKnown
	if !a.value.IsNil() && a.value.IsName() {
		aliasedID := a.value.GetStringID()
		// Never consider this variable to be an alias of itself.
		if aliasedID != a.outer.v.nameID {
			a.aliasedVar = a.outer.v.scope.GetVarByID(aliasedID)
		}
	}
	return a.aliasedVar
}

// checkScopeVarEligibility handles the fast bailouts of
// Java's InliningBehavior.createVarExpert before initializing a
// StandardVarExpert.
func (b *InlineVariables_InliningBehavior) checkScopeVarEligibility(
	v *Var,
	referenceInfo *ReferenceCollection,
	isLocal bool,
) (isDeclaredOrInferredConstant bool, bailoutState inlineVariables_handledState, bailedOut bool) {
	if referenceInfo == nil {
		// If we couldn't collect any reference info, don't try to inline and assume it's unsafe
		// to inline aliases of this variable, too.
		return false, handledVarUnsafe, true
	}

	if v.IsExtern() {
		// TODO(bradfordcsmith): Extern variables are generally unsafe to inline.
		return false, handledVarUnsafe, true
	}

	isDeclaredOrInferredConstant = inlineVariables_isDeclaredOrInferredConst(v)
	if !isDeclaredOrInferredConstant && b.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY {
		// If we're only inlining constants, then we shouldn't inline an alias.
		return false, handledVarUnsafe, true
	}

	cc := b.codingConvention
	if cc == nil {
		cc = b.outer.compiler.GetCodingConvention()
		b.codingConvention = cc
	}

	if cc.IsExportedWithLocal(v.GetName() /* local= */, isLocal) {
		// If the variable is exported, it might be assigned a new value by code we cannot see,
		// so aliases to it are creating snapshots of its state.
		// We cannot inline this variable or its aliases.
		return isDeclaredOrInferredConstant, handledVarUnsafe, true
	}

	if cc.IsPropertyRenameFunction(v.nameNode) {
		// It's not terribly likely that anything creates an alias of our special property rename
		// function, but its value never changes, so it should be safe to inline aliases to it.
		return isDeclaredOrInferredConstant, handledVarSafe, true
	}

	return isDeclaredOrInferredConstant, handledVarUnseen, false
}

// IsSafeUseOfArguments: true if `argumentsNode` is a use of `arguments` we
// know won't modify any parameter values.
//
// In sloppy mode it is possible to change the value of a function parameter
// by assigning to the corresponding entry in `arguments`.
//
// This method checks for just a few common uses that are known to be safe.
// However, we may soon remove this check entirely because unsafe uses
// aren't worth worrying about. See the comment on
// varsInThisScopeMayBeModifiedUsingArguments.
//
// Java: InlineVariables.InliningBehavior.isSafeUseOfArguments(Node)
func (b *InlineVariables_InliningBehavior) IsSafeUseOfArguments(argumentsNode rhino.Node) bool {
	jutil.CheckArgument(argumentsNode.IsName() && argumentsNode.GetStringID() == implicitArgumentsID)
	// `arguments[i]` that is only read, not assigned
	// or `fn.apply(thisArg, arguments)`
	return b.IsTargetOfPropertyRead(argumentsNode) ||
		b.IsSecondArgumentToDotApplyMethod(argumentsNode)
}

// IsTargetOfPropertyRead: true if `n` is the object in a property access
// expression (`obj[prop]` or `obj.prop` or an optional chain version of one
// of those) and the property is only being read, not modified.
//
// Java: InlineVariables.InliningBehavior.isTargetOfPropertyRead(Node)
func (b *InlineVariables_InliningBehavior) IsTargetOfPropertyRead(n rhino.Node) bool {
	get, ok := ast.AsNormalOrOptChainGet(n.GetParent())
	return ok && get.Object().Is(n) && !NodeUtil_IsLValue(get.Node())
}

// IsSecondArgumentToDotApplyMethod: true if `n` is being used in a call like
// this: `fn.apply(thisArg, n)`.
//
// Java: InlineVariables.InliningBehavior.isSecondArgumentToDotApplyMethod(Node)
func (b *InlineVariables_InliningBehavior) IsSecondArgumentToDotApplyMethod(n rhino.Node) bool {
	call, ok := ast.AsNormalOrOptChainCall(n.GetParent())
	if !ok {
		return false
	}
	callee, ok := call.Callee().AsNormalOrOptChainGetProp()
	if !ok || callee.Node().GetStringID() != inlineVariablesApplyID {
		return false
	}
	thisArg, ok := call.FirstArg()
	if !ok {
		return false
	}
	secondArg, ok := thisArg.Next()
	return ok && secondArg.Is(n)
}

// Inline does the actual work of inlining a single declaration into a
// single reference.
//
// Java: InlineVariables.InliningBehavior.inline(Reference, Reference, Reference)
func (b *InlineVariables_InliningBehavior) Inline(decl *Reference, init *Reference, ref *Reference) {
	value := init.GetAssignedValue()
	jutil.CheckState(!value.IsNil())
	// Check for function declarations before the value is moved in the AST.
	isFunctionDeclaration := NodeUtil_IsFunctionDeclaration(value)
	if isFunctionDeclaration {
		// We're inlining this function declaration because there was only one reference to it,
		// probably the rhs of an assignment statement. Since we're eliminating the only
		// reference, we can get rid of it. See b/277538101.
		fn := ast.MustFunction(value)
		fn.RemoveName()

		// In addition to changing the containing scope, inlining function declarations also changes
		// the function name scope from the containing scope to the inner scope.
		b.outer.compiler.ReportChangeToChangeScope(value)
		b.outer.compiler.ReportChangeToEnclosingScope(value.GetParent())
	}
	b.InlineValue(ref.GetNode(), value.Detach())
	if decl != init {
		expressRoot := init.GetGrandparent()
		jutil.CheckState(expressRoot.IsExprResult())
		NodeUtil_RemoveChild(expressRoot.GetParent(), expressRoot)
	}
	// Function declarations have already been removed.
	if !isFunctionDeclaration {
		b.RemoveDeclaration(decl)
	}
}

// InlineWellDefinedVariable inlines an immutable variable into all of its
// references.
//
// Java: InlineVariables.InliningBehavior.inlineWellDefinedVariable(Var, Node, List<Reference>)
func (b *InlineVariables_InliningBehavior) InlineWellDefinedVariable(v *Var, value rhino.Node, refSet []*Reference) {
	for _, r := range refSet {
		if r.GetNode() == v.GetNameNode() {
			b.RemoveDeclaration(r)
		} else if r.IsSimpleAssignmentToName() {
			/*
			 * This is the initialization.
			 *
			 * <p>Replace the entire assignment with just the value, and use the original value node
			 * in case it contains references to variables that still require inlining.
			 */
			b.InlineValue(r.GetParent(), value.Detach())
		} else {
			toRemove := r.GetNode()
			if toRemove.CanCloneLeafInPlace() && value.CanCloneLeafInPlace() {
				b.outer.compiler.ReportChangeToEnclosingScope(toRemove)
				castColor := toRemove.GetColor()
				hasCastColor := castColor != nil && toRemove.IsColorFromTypeCast()
				toRemove.OverwriteLeafFromClone(value)
				if hasCastColor {
					toRemove.SetColor(castColor)
					toRemove.SetColorFromTypeCast()
				}
			} else {
				clonedValue := value.CloneTree()
				NodeUtil_MarkNewScopesChanged(clonedValue, b.outer.compiler)
				b.InlineValue(toRemove, clonedValue)
			}
		}
	}
}

// RemoveDeclaration removes the given VAR declaration.
//
// Java: InlineVariables.InliningBehavior.removeDeclaration(Reference)
func (b *InlineVariables_InliningBehavior) RemoveDeclaration(decl *Reference) {
	varNode, ok := ast.AsNameDeclaration(decl.GetParent())
	jutil.CheckStateMsg1(ok, "%s", decl.GetParent())
	grandparent := decl.GetGrandparent()

	b.outer.compiler.ReportChangeToEnclosingScope(decl.GetNode())
	decl.GetNode().Detach()
	// Remove var node if empty
	if !varNode.HasDeclarators() {
		NodeUtil_RemoveChild(grandparent, varNode.Node())
	}
}

// Java: InlineVariables.InliningBehavior.inlineValue(Node, Node)
func (b *InlineVariables_InliningBehavior) InlineValue(toRemove rhino.Node, toInsert rhino.Node) {
	b.outer.compiler.ReportChangeToEnclosingScope(toRemove)

	// Help type-based optimizations by propagating more specific types from type assertions
	if toRemove.GetColor() != nil && toRemove.IsColorFromTypeCast() {
		toInsert.SetColor(toRemove.GetColor())
		toInsert.SetColorFromTypeCast()
	}
	toRemove.ReplaceWith(toInsert)
	NodeUtil_MarkFunctionsDeleted(toRemove, b.outer.compiler)
}

// CanInline returns true if the provided reference and declaration can be
// safely inlined according to our criteria.
//
// Java: InlineVariables.InliningBehavior.canInline(Reference, Reference, Reference, Node)
func (b *InlineVariables_InliningBehavior) CanInline(
	declaration *Reference, initialization *Reference, reference *Reference, initValue rhino.Node) bool {
	// If the value is read more than once, skip it.
	// VAR declarations and EXPR_RESULT don't need the value, but other
	// ASSIGN expressions parents do.
	if declaration != initialization && !initialization.GetGrandparent().IsExprResult() {
		return false
	}

	// Be very conservative and do not cross control structures or scope boundaries
	if declaration.GetBasicBlock() != initialization.GetBasicBlock() ||
		declaration.GetBasicBlock() != reference.GetBasicBlock() {
		return false
	}

	// Do not inline into a call node. This would change
	// the context in which it was being called. For example,
	//   var a = b.c;
	//   a();
	// should not be inlined, because it calls a in the context of b
	// rather than the context of the window.
	//   var a = b.c;
	//   f(a)
	// is OK.
	jutil.CheckState(!initValue.IsNil())
	call, isCall := ast.AsCall(reference.GetParent())
	if initValue.IsGetProp() && isCall && call.Callee().Is(reference.GetNode()) {
		return false
	}

	if initValue.IsFunction() && isCall {
		convention := b.outer.compiler.GetCodingConvention()
		// Bug 2388531: Don't inline subclass definitions into class defining
		// calls as this confused class removing logic.
		relationship := convention.GetClassesDefinedByCall(call.Node())
		if relationship != nil {
			return false
		}

		// issue 668: Don't inline singleton getter methods
		// calls as this confused class removing logic.
		if convention.GetSingletonGetterClassName(call.Node()) != nil {
			return false
		}
	}

	if initialization.GetScope() != declaration.GetScope() ||
		!initialization.GetScope().Contains(reference.GetScope()) {
		return false
	}

	return b.CanMoveAggressively(initValue) || b.CanMoveModerately(initialization, reference)
}

// CanMoveAggressively: if the value is a literal, we can cross more
// boundaries to inline it.
//
// Java: InlineVariables.InliningBehavior.canMoveAggressively(Node)
func (b *InlineVariables_InliningBehavior) CanMoveAggressively(value rhino.Node) bool {
	// Function expressions and other mutable objects can move within
	// the same basic block.
	return NodeUtil_IsLiteralValue(value, true) || value.IsFunction()
}

// InlineVariables_InliningBehavior_MAX_MODERATE_MOTION_STEPS is the maximum
// number of AST steps LocalVarMotion will traverse before giving up. Prevents
// quadratic traversal in large monolithic scopes with many definitions
// (b/545160233).
//
// 5000 steps is large enough to span massive side-effect-free AST structures
// like protobuf descriptor arrays without regressing code size, while
// bounding worst-case traversal time in huge monolithic scopes to fractions
// of a second.
//
// Java: InlineVariables.InliningBehavior.MAX_MODERATE_MOTION_STEPS
const InlineVariables_InliningBehavior_MAX_MODERATE_MOTION_STEPS = 5000

// CanMoveModerately: if the value of a variable is not constant, then it may
// read or modify state. Therefore it cannot be moved past anything else that
// may modify the value being read or read values that are modified.
//
// Java: InlineVariables.InliningBehavior.canMoveModerately(Reference, Reference)
func (b *InlineVariables_InliningBehavior) CanMoveModerately(initialization *Reference, reference *Reference) bool {
	// Check if declaration can be inlined without passing
	// any side-effect causing nodes.
	it := &b.motionScratch
	if NodeUtil_IsNameDeclaration(initialization.GetParent()) {
		it.initForVar(
			b.outer.compiler,
			initialization.GetNode(),        // NAME
			initialization.GetParent(),      // VAR/LET/CONST
			initialization.GetGrandparent()) // VAR/LET/CONST container
	} else if initialization.GetParent().IsAssign() {
		jutil.CheckState(initialization.GetGrandparent().IsExprResult())
		it.initForAssign(
			b.outer.compiler,
			initialization.GetNode(),                    // NAME
			initialization.GetParent(),                  // ASSIGN
			initialization.GetGrandparent(),             // EXPR_RESULT
			initialization.GetGrandparent().GetParent()) // EXPR container
	} else {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected initialization parent\n" + initialization.GetParent().ToStringTree()))
	}
	targetName := reference.GetNode()
	steps := 0
	for it.HasNext() {
		curNode := it.Next()
		if curNode == targetName {
			return true
		}
		steps++
		if steps > InlineVariables_InliningBehavior_MAX_MODERATE_MOTION_STEPS {
			return false
		}
	}

	return false
}

// IsValidDeclaration returns true if the reference is a normal VAR or
// FUNCTION declaration.
//
// Java: InlineVariables.InliningBehavior.isValidDeclaration(Reference)
func (b *InlineVariables_InliningBehavior) IsValidDeclaration(declaration *Reference) bool {
	return (NodeUtil_IsNameDeclaration(declaration.GetParent()) &&
		!NodeUtil_IsLoopStructure(declaration.GetGrandparent())) ||
		NodeUtil_IsFunctionDeclaration(declaration.GetParent())
}

// IsValidInitialization returns whether there is a initial value.
//
// Java: InlineVariables.InliningBehavior.isValidInitialization(Reference)
func (b *InlineVariables_InliningBehavior) IsValidInitialization(initialization *Reference) bool {
	if initialization == nil {
		return false
	} else if initialization.IsDeclaration() {
		// The reference is a FUNCTION declaration or normal VAR declaration
		// with a value.
		if !NodeUtil_IsFunctionDeclaration(initialization.GetParent()) &&
			!initialization.GetNode().HasChildren() {
			return false
		}
	} else {
		assign, ok := ast.AsAssign(initialization.GetParent())
		jutil.CheckState(ok && assign.Target().Is(initialization.GetNode()))
	}

	return true
}

// IsValidReference returns true if the reference is a candidate for
// inlining.
//
// Java: InlineVariables.InliningBehavior.isValidReference(Reference)
func (b *InlineVariables_InliningBehavior) IsValidReference(reference *Reference) bool {
	return !reference.IsDeclaration() && !reference.IsLvalue()
}

// Java: InlineVariables.hasNoInlineAnnotation(Var)
func InlineVariables_HasNoInlineAnnotation(var_ *Var) bool {
	n := var_.nameNode
	if n.IsNil() {
		return false
	}
	if info := n.GetJSDocInfo(); info != nil {
		return info.IsNoInline()
	}
	parent := n.GetParent()
	if parent.IsNil() {
		return false
	}
	switch parent.GetToken() {
	case rhino.Token_PARAM_LIST:
		return false
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_FUNCTION, rhino.Token_CLASS:
		if parent.GetFirstChild() == n {
			if info := parent.GetJSDocInfo(); info != nil {
				return info.IsNoInline()
			}
			if gp := parent.GetParent(); !gp.IsNil() && gp.IsExport() {
				info := gp.GetJSDocInfo()
				return info != nil && info.IsNoInline()
			}
			return false
		}
	}
	jsDocInfo := NodeUtil_GetBestJSDocInfo(n)
	return jsDocInfo != nil && jsDocInfo.IsNoInline()
}
