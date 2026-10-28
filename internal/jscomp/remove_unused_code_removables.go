// Port of com/google/javascript/jscomp/RemoveUnusedCode.java @ closure-compiler v20261021

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

// RemoveUnusedCode_Continuation: our progress in a traversal can be
// expressed completely as the current node and scope. The continuation lets
// us save that information so that we can continue the traversal later.
type RemoveUnusedCode_Continuation struct {
	node  rhino.Node
	scope *Scope
}

// Java: RemoveUnusedCode.Continuation.apply()
func (c RemoveUnusedCode_Continuation) Apply(outer *RemoveUnusedCode) {
	if c.node.IsFunction() {
		// Calling traverseNode here would create infinite recursion for a function declaration
		outer.TraverseFunction(c.node, c.scope)
	} else {
		outer.TraverseNode(c.node, c.scope)
	}
}

const (
	remKindUnusedReadRef uint8 = iota + 1
	remKindInstanceofName
	remKindWeakUsageCall
	remKindIncOrDecOp
	remKindIndirectAssign
	remKindPolyfill
	remKindClassDecl
	remKindNamedClassExpr
	remKindClassOrProtoProp
	remKindObjectDefinePropsDef
	remKindFunctionDecl
	remKindNameDeclStmt
	remKindAssignVar
	remKindAssignNamedProp
	remKindAssignComputedProp
	remKindAnonProtoPropAssign
	remKindClassSetupCall
	remKindVanillaForNameDecl
)

const (
	remFlagProtoDotPropRef uint8 = 1 << iota
	remFlagThisDotPropRef
	remFlagContsApplied
	remFlagRemoved
	remFlagInPropList
	remFlagUninitVarDecl
	remFlagSubBit // isPatch for Polyfill, classDefiningCall for ClassSetupCall
)

// RemoveUnusedCode_RemovableBase holds all Removable state in a single 64-byte
// cache line (8 words, 6 pointers).
type RemoveUnusedCode_RemovableBase struct {
	propNext *RemoveUnusedCode_RemovableBase

	cont0 RemoveUnusedCode_Continuation

	// Primary AST node for the removable construct.
	targetNode rhino.Node

	// Auxiliary AST node: root for IndirectAssign, toPreserve for IncOrDecOp.
	auxNode rhino.Node

	// Associated VarInfo for Assign.
	varInfo *RemoveUnusedCode_RealVarInfo

	propertyNameID uint32
	extraContIdx   uint32
	numConts       uint16
	remKind        uint8
	flags          uint8
}

func (b *RemoveUnusedCode_RemovableBase) addContinuation(r *RemoveUnusedCode, node rhino.Node, scope *Scope) {
	c := RemoveUnusedCode_Continuation{node: node, scope: scope}
	switch b.numConts {
	case 0:
		b.cont0 = c
	case 1:
		b.extraContIdx = uint32(len(r.extraContSlab))
		r.extraContSlab = append(r.extraContSlab, c)
	default:
		r.extraContSlab = append(r.extraContSlab, c)
	}
	b.numConts++
}

// Remove removes the associated nodes from the AST, unless they've already
// been removed.
//
// Java: RemoveUnusedCode.Removable.remove(AbstractCompiler)
func (b *RemoveUnusedCode_RemovableBase) Remove(r *RemoveUnusedCode) {
	if b.flags&remFlagRemoved == 0 {
		b.flags |= remFlagRemoved
		b.RemoveInternal(r)
	}
}

// Java: RemoveUnusedCode.Removable.applyContinuations()
func (b *RemoveUnusedCode_RemovableBase) ApplyContinuations(r *RemoveUnusedCode) {
	if b.flags&remFlagContsApplied == 0 {
		b.flags |= remFlagContsApplied
		// Enqueue the continuation for processing.
		// Don't invoke the continuation immediately, because that can lead to concurrent
		// modification of data structures.
		if b.numConts >= 1 {
			r.enqueueContinuation(b.cont0)
			b.cont0 = RemoveUnusedCode_Continuation{}
		}
		if b.numConts >= 2 {
			end := b.extraContIdx + uint32(b.numConts-1)
			extras := r.extraContSlab[b.extraContIdx:end]
			for _, c := range extras {
				r.enqueueContinuation(c)
			}
			clear(extras)
		}
		b.numConts = 0
	}
}

// IsVariableAssignment: true if this object represents assignment to a
// variable.
//
// Java: RemoveUnusedCode.Removable.isVariableAssignment()
func (b *RemoveUnusedCode_RemovableBase) IsVariableAssignment() bool {
	switch b.remKind {
	case remKindAssignVar, remKindClassDecl, remKindFunctionDecl, remKindNameDeclStmt, remKindVanillaForNameDecl:
		return true
	case remKindIndirectAssign:
		return b.targetNode.IsName()
	default:
		return false
	}
}

// IsNamedProperty: true if this object represents a named property, either
// assignment or declaration.
//
// Java: RemoveUnusedCode.Removable.isNamedProperty()
func (b *RemoveUnusedCode_RemovableBase) IsNamedProperty() bool {
	return b.propertyNameID != 0
}

// IsNamedPropertyAssignment: true if this object represents assignment to a
// named property.
//
// This does not include class or object literal member declarations.
//
// Java: RemoveUnusedCode.Removable.isNamedPropertyAssignment()
func (b *RemoveUnusedCode_RemovableBase) IsNamedPropertyAssignment() bool {
	switch b.remKind {
	case remKindAssignNamedProp:
		return true
	case remKindIndirectAssign:
		return b.targetNode.IsGetProp()
	default:
		return false
	}
}

// Java: RemoveUnusedCode.Removable.isAssignedValueLocal()
func (b *RemoveUnusedCode_RemovableBase) IsAssignedValueLocal() bool {
	switch b.remKind {
	case remKindClassDecl, remKindFunctionDecl:
		return true
	case remKindNameDeclStmt, remKindVanillaForNameDecl:
		name := b.declaredName()
		initialValue, ok := name.Init()
		if !ok {
			return true
		}
		valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(name.Node(), initialValue.Node())
		return NodeUtil_EvaluatesToLocalValue(valueNode)
	case remKindAssignVar, remKindAssignNamedProp, remKindAssignComputedProp:
		return !b.GetLocalAssignedValue().IsNil()
	default:
		return false
	}
}

// GetLocalAssignedValue returns the Node representing the local value that
// is being assigned or `null` if the value is non-local or cannot be
// determined.
//
// Java: RemoveUnusedCode.Removable.getLocalAssignedValue()
func (b *RemoveUnusedCode_RemovableBase) GetLocalAssignedValue() rhino.Node {
	switch b.remKind {
	case remKindClassDecl, remKindFunctionDecl:
		return b.targetNode
	case remKindNameDeclStmt, remKindVanillaForNameDecl:
		name := b.declaredName()
		initialValue, ok := name.Init()
		if !ok {
			return rhino.NilNode
		}
		valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(name.Node(), initialValue.Node())
		if NodeUtil_EvaluatesToLocalValue(valueNode) {
			return valueNode
		}
		return rhino.NilNode
	case remKindAssignVar, remKindAssignNamedProp, remKindAssignComputedProp:
		if NodeUtil_IsExpressionResultUsed(b.targetNode) {
			return rhino.NilNode
		}
		assign := ast.MustAssignment(b.targetNode)
		valueNode := RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(
			assign.Target().Node(), assign.Value().Node())
		if NodeUtil_EvaluatesToLocalValue(valueNode) {
			return valueNode
		}
		return rhino.NilNode
	default:
		return rhino.NilNode
	}
}

// declaredName returns the name a NameDeclarationStatement (the only one of
// its statement) or a VanillaForNameDeclaration (the name itself) declares.
func (b *RemoveUnusedCode_RemovableBase) declaredName() ast.Name {
	if b.remKind == remKindVanillaForNameDecl {
		return ast.MustName(b.targetNode)
	}
	declarator, ok := ast.MustNameDeclaration(b.targetNode).OnlyDeclarator()
	jutil.CheckState(ok)
	return declarator.MustName()
}

// IsPrototypeAssignment: is this a direct assignment to `varName.prototype`?
//
// Java: RemoveUnusedCode.Removable.isPrototypeAssignment()
func (b *RemoveUnusedCode_RemovableBase) IsPrototypeAssignment() bool {
	return b.IsNamedPropertyAssignment() && b.propertyNameID == rucPropPrototypeID
}

// IsPrototypeDotPropertyReference: is this an assignment to a property on a
// prototype object?
//
// Java: RemoveUnusedCode.Removable.isPrototypeDotPropertyReference()
func (b *RemoveUnusedCode_RemovableBase) IsPrototypeDotPropertyReference() bool {
	return b.flags&remFlagProtoDotPropRef != 0
}

// Java: RemoveUnusedCode.Removable.isClassOrPrototypeNamedProperty()
func (b *RemoveUnusedCode_RemovableBase) IsClassOrPrototypeNamedProperty() bool {
	return b.remKind == remKindClassOrProtoProp && !b.targetNode.IsStaticMember()
}

// Java: RemoveUnusedCode.Removable.isPrototypeProperty()
func (b *RemoveUnusedCode_RemovableBase) IsPrototypeProperty() bool {
	return b.remKind == remKindAnonProtoPropAssign ||
		b.IsPrototypeDotPropertyReference() ||
		b.IsClassOrPrototypeNamedProperty()
}

// Java: RemoveUnusedCode.Removable.isThisDotPropertyReference()
func (b *RemoveUnusedCode_RemovableBase) IsThisDotPropertyReference() bool {
	return b.flags&remFlagThisDotPropRef != 0
}

// Java: RemoveUnusedCode.Removable.isObjectDefinePropertiesDefinition()
func (b *RemoveUnusedCode_RemovableBase) IsObjectDefinePropertiesDefinition() bool {
	return b.remKind == remKindObjectDefinePropsDef
}

// IsStaticProperty.
//
// TODO(b/134610338): Combine this method with `isPrototypeProperty`.
//
// Java: RemoveUnusedCode.Removable.isStaticProperty()
func (b *RemoveUnusedCode_RemovableBase) IsStaticProperty() bool {
	switch b.remKind {
	case remKindClassOrProtoProp:
		return b.targetNode.IsStaticMember()
	case remKindAssignNamedProp:
		if b.varInfo != nil && b.varInfo.HasFunctionOrClassLiteralValue() {
			// varName.propertyName = something
			target := ast.MustAssignment(b.targetNode).Target().MustGetProp()
			return target.Object().Token() == rhino.Token_NAME
		}
		return false
	default:
		return false
	}
}

// PreventsRemovalOfVariableWithNonLocalValueOrPrototype: would a nonlocal
// or nonliteral value prevent removal of a variable associated with this
// Removable?
//
// Java: RemoveUnusedCode.Removable.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (b *RemoveUnusedCode_RemovableBase) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	switch b.remKind {
	case remKindInstanceofName, remKindClassSetupCall, remKindAssignNamedProp, remKindAssignComputedProp:
		return true
	case remKindIndirectAssign:
		if getProp, ok := ast.AsGetProp(b.targetNode); ok {
			// varName.propertyName or varName.prototype.propertyName
			getPropLhs := getProp.Object()
			_, isNameDotPrototype := removeUnusedCode_nameDotPrototype(getPropLhs.Node())
			return getPropLhs.Token() == rhino.Token_NAME || isNameDotPrototype
		}
		return false
	default:
		return false
	}
}

func (b *RemoveUnusedCode_RemovableBase) RemoveInternal(r *RemoveUnusedCode) {
	compiler := r.compiler
	switch b.remKind {
	case remKindUnusedReadRef:
		if !RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			if RemoveUnusedCode_IsThisDotProperty(b.targetNode) {
				r.RemoveExpressionCompletely(b.targetNode)
			} else {
				if !RemoveUnusedCode_IsDotPrototypeDotProperty(b.targetNode) {
					jutil.CheckStateMsg1(false, "%s", b.targetNode)
				}
				// objExpression.prototype.propertyName
				owner, _ := removeUnusedCode_prototypeOwner(ast.MustNormalOrOptChainGetProp(b.targetNode).Object().Node())
				objExpression := owner.Node()
				if r.MayHaveSideEffects(objExpression) {
					r.ReplaceNodeWith(b.targetNode, objExpression.Detach())
				} else {
					r.RemoveExpressionCompletely(b.targetNode)
				}
			}
		}
	case remKindInstanceofName:
		if !RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			lhs := ast.MustBinary(b.targetNode).Left().Node()
			falseNode := rhino.IR_FalseNode().Srcref(b.targetNode)
			if r.MayHaveSideEffects(lhs) {
				r.ReplaceNodeWith(b.targetNode, rhino.IR_Comma(lhs.Detach(), falseNode).Srcref(b.targetNode))
			} else {
				r.ReplaceNodeWith(b.targetNode, falseNode)
			}
		}
	case remKindWeakUsageCall:
		if !RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			r.ReplaceNodeWith(b.targetNode, NodeUtil_NewUndefinedNode(b.targetNode))
		}
	case remKindIncOrDecOp:
		if RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			return
		}
		arg := ast.MustIncDec(b.targetNode).Target()
		jutil.CheckStateMsg1(arg.Token() == rhino.Token_GETPROP, "%s", arg.Node())
		if b.auxNode.IsNil() {
			r.RemoveExpressionCompletely(b.targetNode)
		} else {
			r.ReplaceNodeWith(b.targetNode, b.auxNode.Detach())
		}
	case remKindIndirectAssign:
		if !RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			b.removeIndirectAssignRoot(r)
		}
	case remKindPolyfill, remKindClassDecl, remKindClassOrProtoProp, remKindObjectDefinePropsDef, remKindFunctionDecl:
		NodeUtil_DeleteNode(b.targetNode, compiler)
	case remKindNamedClassExpr:
		if !RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			cls := ast.MustClass(b.targetNode)
			if _, ok := cls.Name(); ok {
				cls.RemoveName()
				compiler.ReportChangeToEnclosingScope(b.targetNode)
			}
		}
	case remKindNameDeclStmt:
		value, hasValue := b.declaredName().Init()
		valueNode := value.Node()
		if hasValue && r.MayHaveSideEffects(valueNode) {
			compiler.ReportChangeToEnclosingScope(b.targetNode)
			valueNode.Detach()
			b.targetNode.ReplaceWith(rhino.IR_ExprResult(valueNode).Srcref(valueNode))
		} else {
			NodeUtil_DeleteNode(b.targetNode, compiler)
		}
	case remKindAssignVar, remKindAssignNamedProp, remKindAssignComputedProp:
		if RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			return
		}
		parent := b.targetNode.GetParent()
		compiler.ReportChangeToEnclosingScope(parent)
		assign := ast.MustAssignment(b.targetNode)
		rhs := assign.Value().Node()
		mustPreserveRhs :=
			r.MayHaveSideEffects(rhs) || NodeUtil_IsExpressionResultUsed(b.targetNode)
		// The key of varName[key] = value.
		getElem, isGetElem := assign.Target().AsGetElem()
		mustPreserveGetElmExpr := isGetElem && r.MayHaveSideEffects(getElem.Index().Node())

		if mustPreserveRhs && mustPreserveGetElmExpr {
			replacement := rhino.IR_Comma(getElem.Index().Detach().Node(), rhs.Detach()).Srcref(b.targetNode)
			r.ReplaceNodeWith(b.targetNode, replacement)
		} else if mustPreserveGetElmExpr {
			r.ReplaceNodeWith(b.targetNode, getElem.Index().Detach().Node())
		} else if mustPreserveRhs {
			r.ReplaceNodeWith(b.targetNode, rhs.Detach())
		} else {
			r.RemoveExpressionCompletely(b.targetNode)
		}
	case remKindAnonProtoPropAssign:
		if RemoveUnusedCode_AlreadyRemoved(b.targetNode) {
			return
		}
		parent := b.targetNode.GetParent()
		compiler.ReportChangeToEnclosingScope(parent)
		assign := ast.MustAssignment(b.targetNode)
		rhs := assign.Value().Node()

		// objExpression.prototype.propertyName = rhs
		lhs, ok := assign.Target().AsGetProp()
		jutil.CheckStateMsg1(ok, "%s", assign.Target().Node())
		objDotPrototype, ok := lhs.Object().AsGetProp()
		jutil.CheckStateMsg1(ok, "%s", lhs.Object().Node())
		objExpression := objDotPrototype.Object().Node()
		jutil.CheckStateMsg1(objDotPrototype.Node().GetStringID() == rucPropPrototypeID, "%s", objDotPrototype.Node())

		mustPreserveRhs :=
			r.MayHaveSideEffects(rhs) || NodeUtil_IsExpressionResultUsed(b.targetNode)
		mustPreserveObjExpression := r.MayHaveSideEffects(objExpression)

		if mustPreserveRhs && mustPreserveObjExpression {
			replacement := rhino.IR_Comma(objExpression.Detach(), rhs.Detach()).Srcref(b.targetNode)
			r.ReplaceNodeWith(b.targetNode, replacement)
		} else if mustPreserveObjExpression {
			r.ReplaceNodeWith(b.targetNode, objExpression.Detach())
		} else if mustPreserveRhs {
			r.ReplaceNodeWith(b.targetNode, rhs.Detach())
		} else {
			r.RemoveExpressionCompletely(b.targetNode)
		}
	case remKindClassSetupCall:
		parent := b.targetNode.GetParent()
		classDefiningCall := b.flags&remFlagSubBit != 0

		var replacement rhino.Node
		b.targetNode.RemoveFirstChild()
		for arg := b.targetNode.GetLastChild(); !arg.IsNil(); arg = b.targetNode.GetLastChild() {
			arg.Detach()
			if !classDefiningCall && r.MayHaveSideEffects(arg) {
				if replacement.IsNil() {
					replacement = arg
				} else {
					replacement = rhino.IR_Comma(arg, replacement).Srcref(b.targetNode)
				}
			} else {
				NodeUtil_MarkFunctionsDeleted(arg, compiler)
			}
		}

		if !replacement.IsNil() {
			r.ReplaceNodeWith(b.targetNode, replacement)
		} else if parent.IsExprResult() {
			NodeUtil_DeleteNode(parent, compiler)
		} else if parent.IsComma() || parent.IsAnd() || parent.IsOr() {
			bin := ast.MustBinary(parent)
			var keep rhino.Node
			if bin.Left().Is(b.targetNode) {
				keep = bin.Right().Detach().Node()
			} else {
				keep = bin.Left().Detach().Node()
			}
			compiler.ReportChangeToEnclosingScope(parent)
			parent.ReplaceWith(keep)
		} else {
			compiler.ReportChangeToEnclosingScope(parent)
			b.targetNode.ReplaceWith(rhino.IR_Number(0))
		}
	case remKindVanillaForNameDecl:
		declaration := rhino.CheckNotNull(b.targetNode.GetParent())
		compiler.ReportChangeToEnclosingScope(declaration)
		if b.targetNode.GetPrevious().IsNil() && b.targetNode.GetNext().IsNil() {
			declaration.ReplaceWith(rhino.IR_Empty().Srcref(declaration))
		} else {
			b.targetNode.Detach()
		}
		NodeUtil_MarkFunctionsDeleted(b.targetNode, compiler)
	}
}

func (b *RemoveUnusedCode_RemovableBase) removeIndirectAssignRoot(r *RemoveUnusedCode) {
	root := b.auxNode
	rootParent := root.GetParent()

	switch rootParent.GetToken() {
	case rhino.Token_ARRAY_PATTERN:
		r.ReplaceNodeWith(root, rhino.IR_Empty().Srcref(root))
		for maybeEmpty := rootParent.GetLastChild(); !maybeEmpty.IsNil() && maybeEmpty.IsEmpty(); maybeEmpty = rootParent.GetLastChild() {
			maybeEmpty.Detach()
		}
		r.compiler.ReportChangeToEnclosingScope(rootParent)
	case rhino.Token_PARAM_LIST:
		// function f(name = value): keep the name, unless the parameter can go.
		param, ok := ast.AsDefaultValue(root)
		if !ok {
			return
		}
		r.compiler.ReportChangeToEnclosingScope(rootParent)
		name, ok := param.Target().AsName()
		jutil.CheckState(ok)
		if param.AsParam().IsLast() &&
			r.removeGlobals &&
			r.CanRemoveParameters(rootParent) {
			root.Detach()
		} else {
			root.ReplaceWith(name.Detach().Node())
		}
		NodeUtil_MarkFunctionsDeleted(root, r.compiler)
	case rhino.Token_OBJECT_PATTERN:
		NodeUtil_DeleteNode(root, r.compiler)
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected parent of indirect assignment: " + rootParent.ToStringTree()))
	}
}

func (b *RemoveUnusedCode_RemovableBase) String() string {
	switch b.remKind {
	case remKindUnusedReadRef:
		return "UnusedReadReference:" + b.targetNode.String()
	case remKindInstanceofName:
		return "InstanceofName:" + b.targetNode.String()
	case remKindWeakUsageCall:
		return "weakUsageCall:" + b.targetNode.String()
	case remKindIncOrDecOp:
		return "IncOrDecOp:" + b.targetNode.String()
	case remKindPolyfill:
		if b.flags&remFlagSubBit != 0 {
			return "Patch:" + b.targetNode.String()
		}
		return "Polyfill:" + b.targetNode.String()
	case remKindClassDecl:
		return "ClassDeclaration:" + b.targetNode.String()
	case remKindNamedClassExpr:
		return "NamedClassExpression:" + b.targetNode.String()
	case remKindClassOrProtoProp:
		return "ClassOrPrototypeNamedProperty:" + b.targetNode.String()
	case remKindFunctionDecl:
		return "FunctionDeclaration:" + b.targetNode.String()
	case remKindNameDeclStmt:
		return "NameDeclStmt:" + b.targetNode.String()
	case remKindAssignVar, remKindAssignNamedProp, remKindAssignComputedProp:
		return "Assign:" + b.targetNode.String()
	case remKindAnonProtoPropAssign:
		return "AnonymousPrototypeNamedPropertyAssign:" + b.targetNode.String()
	case remKindClassSetupCall:
		return "ClassSetupCall:" + b.targetNode.String()
	default:
		return "Removable:" + b.targetNode.String()
	}
}

func (r *RemoveUnusedCode) newUnusedReadRefRemovable(referenceNode rhino.Node, flags uint8) *RemoveUnusedCode_RemovableBase {
	if !(RemoveUnusedCode_IsThisDotProperty(referenceNode) || RemoveUnusedCode_IsDotPrototypeDotProperty(referenceNode)) {
		jutil.CheckStateMsg1(false, "%s", referenceNode)
	}
	u := r.removables.New()
	u.remKind = remKindUnusedReadRef
	u.targetNode = referenceNode
	u.propertyNameID = referenceNode.GetStringID()
	u.flags = flags
	return u
}

func (r *RemoveUnusedCode) newInstanceofNameRemovable(instanceofNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	jutil.CheckArgumentMsg1(instanceofNode.IsInstanceOf(), "%s", instanceofNode)
	i := r.removables.New()
	i.remKind = remKindInstanceofName
	i.targetNode = instanceofNode
	return i
}

func (r *RemoveUnusedCode) newWeakUsageCallRemovable(weakCallNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	jutil.CheckArgumentMsg1(weakCallNode.IsCall(), "%s", weakCallNode)
	w := r.removables.New()
	w.remKind = remKindWeakUsageCall
	w.targetNode = weakCallNode
	return w
}

func (r *RemoveUnusedCode) newIncOrDecOpRemovable(incOrDecNode, propNode, toPreserve rhino.Node, flags uint8) *RemoveUnusedCode_RemovableBase {
	incDec, ok := ast.AsIncDec(incOrDecNode)
	jutil.CheckArgumentMsg1(ok, "%s", incOrDecNode)
	arg := incDec.Target().Node()
	jutil.CheckStateMsg1(RemoveUnusedCode_IsThisDotProperty(arg) || RemoveUnusedCode_IsDotPrototypeDotProperty(arg), "%s", arg)
	o := r.removables.New()
	o.remKind = remKindIncOrDecOp
	o.targetNode = incOrDecNode
	o.auxNode = toPreserve
	o.propertyNameID = propNode.GetStringID()
	o.flags = flags
	return o
}

func (r *RemoveUnusedCode) newIndirectAssignRemovable(root, targetNode rhino.Node, scope *Scope) *RemoveUnusedCode_RemovableBase {
	rootParent := root.GetParent()
	if !(rootParent.IsDestructuringPattern() || rootParent.IsParamList()) {
		jutil.CheckStateMsg1(false, "%s", rootParent)
	}
	if !(targetNode.IsName() || targetNode.IsGetProp()) {
		jutil.CheckStateMsg1(false, "%s", targetNode)
	}
	a := r.removables.New()
	a.remKind = remKindIndirectAssign
	a.targetNode = targetNode
	a.auxNode = root
	if targetNode.IsGetProp() {
		a.propertyNameID = targetNode.GetStringID()
		if RemoveUnusedCode_IsThisDotProperty(targetNode) {
			a.flags = remFlagThisDotPropRef
		}
	}
	a.addContinuation(r, root, scope)
	return a
}

func (r *RemoveUnusedCode) newPolyfillRemovable(call ast.NormalOrOptChainCall, scope *Scope) *RemoveUnusedCode_RemovableBase {
	p := r.removables.New()
	p.remKind = remKindPolyfill
	p.targetNode = call.Node().GetParent()
	// The callee: $jscomp.patch or (collapsed) $jscomp$patch.
	if calleeID := call.Callee().Node().GetStringID(); calleeID == rucCalleePatchID || calleeID == rucCalleeJscompPatchID {
		p.flags = remFlagSubBit
	}
	for arg := range call.Args() {
		p.addContinuation(r, arg.Node(), scope)
	}
	return p
}

func (r *RemoveUnusedCode) newClassDeclRemovable(classNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	c := r.removables.New()
	c.remKind = remKindClassDecl
	c.targetNode = classNode
	return c
}

func (r *RemoveUnusedCode) newNamedClassExprRemovable(classNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	c := r.removables.New()
	c.remKind = remKindNamedClassExpr
	c.targetNode = classNode
	return c
}

func (r *RemoveUnusedCode) newClassOrProtoPropRemovable(propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	if !(propertyNode.IsMemberFunctionDef() ||
		propertyNode.IsMemberFieldDef() ||
		NodeUtil_IsGetOrSetKey(propertyNode) ||
		(propertyNode.IsStringKey() && !propertyNode.IsQuotedStringKey())) {
		jutil.CheckArgumentMsg1(false, "%s", propertyNode)
	}
	p := r.removables.New()
	p.remKind = remKindClassOrProtoProp
	p.targetNode = propertyNode
	p.propertyNameID = propertyNode.GetStringID()
	return p
}

func (r *RemoveUnusedCode) newObjectDefinePropsDefRemovable(propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	d := r.removables.New()
	d.remKind = remKindObjectDefinePropsDef
	d.targetNode = propertyNode
	d.propertyNameID = propertyNode.GetStringID()
	return d
}

func (r *RemoveUnusedCode) newFunctionDeclRemovable(functionNode rhino.Node, scope *Scope) *RemoveUnusedCode_RemovableBase {
	f := r.removables.New()
	f.remKind = remKindFunctionDecl
	f.targetNode = functionNode
	f.addContinuation(r, functionNode, scope)
	return f
}

func (r *RemoveUnusedCode) newNameDeclStmtRemovable(declarationStatement rhino.Node, hasInit bool) *RemoveUnusedCode_RemovableBase {
	jutil.CheckArgumentMsg1(NodeUtil_IsNameDeclaration(declarationStatement), "%s", declarationStatement)
	s := r.removables.New()
	s.remKind = remKindNameDeclStmt
	s.targetNode = declarationStatement
	if !hasInit {
		s.flags = remFlagUninitVarDecl
	}
	return s
}

func (r *RemoveUnusedCode) newVanillaForNameDeclRemovable(nameNode rhino.Node, hasInit bool) *RemoveUnusedCode_RemovableBase {
	d := r.removables.New()
	d.remKind = remKindVanillaForNameDecl
	d.targetNode = nameNode
	if !hasInit {
		d.flags = remFlagUninitVarDecl
	}
	return d
}

func (r *RemoveUnusedCode) newVarAssignRemovable(assignNode rhino.Node, varInfo *RemoveUnusedCode_RealVarInfo) *RemoveUnusedCode_RemovableBase {
	jutil.CheckArgumentMsg1(NodeUtil_IsAssignmentOp(assignNode), "%s", assignNode)
	jutil.CheckArgumentMsg(varInfo != nil, "missing VarInfo for variable assignment")
	a := r.removables.New()
	a.remKind = remKindAssignVar
	a.targetNode = assignNode
	a.varInfo = varInfo
	return a
}

func (r *RemoveUnusedCode) newNamedPropAssignRemovable(assignNode, propNode rhino.Node, varInfo *RemoveUnusedCode_RealVarInfo, flags uint8) *RemoveUnusedCode_RemovableBase {
	jutil.CheckArgumentMsg1(NodeUtil_IsAssignmentOp(assignNode), "%s", assignNode)
	jutil.CheckArgumentMsg1(propNode.IsGetProp(), "property name is not a GETPROP: %s", propNode)
	a := r.removables.New()
	a.remKind = remKindAssignNamedProp
	a.targetNode = assignNode
	a.varInfo = varInfo
	a.propertyNameID = propNode.GetStringID()
	a.flags = flags
	return a
}

func (r *RemoveUnusedCode) newComputedPropAssignRemovable(assignNode, propNode rhino.Node, varInfo *RemoveUnusedCode_RealVarInfo) *RemoveUnusedCode_RemovableBase {
	jutil.CheckArgumentMsg1(NodeUtil_IsAssignmentOp(assignNode), "%s", assignNode)
	jutil.CheckArgumentMsg(!propNode.IsNil(), "missing property node")
	a := r.removables.New()
	a.remKind = remKindAssignComputedProp
	a.targetNode = assignNode
	a.varInfo = varInfo
	return a
}

func (r *RemoveUnusedCode) newAnonProtoPropAssignRemovable(assignNode, propNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	jutil.CheckArgumentMsg1(assignNode.IsAssign(), "%s", assignNode)
	a := r.removables.New()
	a.remKind = remKindAnonProtoPropAssign
	a.targetNode = assignNode
	a.propertyNameID = propNode.GetStringID()
	a.flags = remFlagProtoDotPropRef
	return a
}

func (r *RemoveUnusedCode) newClassSetupCallRemovable(callNode rhino.Node, classDefiningCall bool) *RemoveUnusedCode_RemovableBase {
	c := r.removables.New()
	c.remKind = remKindClassSetupCall
	c.targetNode = callNode
	if classDefiningCall {
		c.flags = remFlagSubBit
	}
	return c
}

// Java: RemoveUnusedCode.traverseRemovableAssignValue(Node, RemovableBuilder, Scope)
func (r *RemoveUnusedCode) traverseRemovableAssignValue(rem *RemoveUnusedCode_RemovableBase, valueNode rhino.Node, scope *Scope) {
	if r.MayHaveSideEffects(valueNode) || NodeUtil_IsExpressionResultUsed(valueNode.GetParent()) {
		r.TraverseNode(valueNode, scope)
	} else {
		rem.addContinuation(r, valueNode, scope)
	}
}

// RemoveUnusedCode_IsThisDotProperty: true for `this.propertyName`.
//
// Java: RemoveUnusedCode.isThisDotProperty(Node)
func RemoveUnusedCode_IsThisDotProperty(n rhino.Node) bool {
	access, ok := ast.AsNormalOrOptChainGetProp(n)
	return ok && access.Object().Token() == rhino.Token_THIS
}

// RemoveUnusedCode_IsDotPrototypeDotProperty: true for
// `(something).prototype.propertyName`.
//
// Java: RemoveUnusedCode.isDotPrototypeDotProperty(Node)
func RemoveUnusedCode_IsDotPrototypeDotProperty(n rhino.Node) bool {
	access, ok := ast.AsNormalOrOptChainGetProp(n)
	return ok && RemoveUnusedCode_IsDotPrototype(access.Object().Node())
}

// RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode: targetNode is the
// node to which a value is being assigned; valueNode the value being
// assigned. If `valueNode` has the form `qualifiedName || defaultValue` and
// `qualifiedName` matches `targetNode`, return `defaultValue`. Otherwise
// return `valueNode`.
//
// Java: RemoveUnusedCode.maybeUnwrapQnameOrDefaultValueNode(Node, Node)
func RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(targetNode rhino.Node, valueNode rhino.Node) rhino.Node {
	if or, ok := ast.AsBinaryOf(valueNode, rhino.Token_OR); ok && targetNode.IsQualifiedName() {
		if or.Left().Node().MatchesQualifiedNameNode(targetNode) {
			return or.Right().Node()
		}
	}
	return valueNode
}

// Java: RemoveUnusedCode.mayHaveSideEffects(Node)
func (r *RemoveUnusedCode) MayHaveSideEffects(node rhino.Node) bool {
	switch node.GetToken() {
	case rhino.Token_NAME, rhino.Token_THIS, rhino.Token_NUMBER, rhino.Token_STRINGLIT,
		rhino.Token_TRUE, rhino.Token_FALSE, rhino.Token_NULL, rhino.Token_EMPTY,
		rhino.Token_FUNCTION, rhino.Token_BIGINT:
		return false
	case rhino.Token_VOID:
		return r.MayHaveSideEffects(ast.MustUnary(node).Operand().Node())
	case rhino.Token_CALL:
		if jsDocInfo := node.GetJSDocInfo(); jsDocInfo != nil && jsDocInfo.IsPureOrBreakMyCode() {
			return false
		}
	}
	return r.astAnalyzer.MayHaveSideEffects(node)
}

// Java: RemoveUnusedCode.alreadyRemoved(Node)
func RemoveUnusedCode_AlreadyRemoved(n rhino.Node) bool {
	for parent := n.GetParent(); !parent.IsNil(); parent = parent.GetParent() {
		if parent.IsRoot() {
			return false
		}
	}
	return true
}
