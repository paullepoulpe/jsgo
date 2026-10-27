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
	"strings"

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

// Java: RemoveUnusedCode.Continuation(Node, Scope)
func NewRemoveUnusedCode_Continuation(_ *RemoveUnusedCode, node rhino.Node, scope *Scope) RemoveUnusedCode_Continuation {
	return RemoveUnusedCode_Continuation{node: node, scope: scope}
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

// RemoveUnusedCode_Removable represents a portion of the AST that can be
// removed.
type RemoveUnusedCode_Removable = *RemoveUnusedCode_RemovableBase

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
	remFlagHasPropertyName uint8 = 1 << iota
	remFlagProtoDotPropRef
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

func (b *RemoveUnusedCode_RemovableBase) initFromBuilder(remKind uint8, targetNode rhino.Node, builder *RemoveUnusedCode_RemovableBuilder) {
	b.cont0 = builder.cont0
	b.targetNode = targetNode
	b.propertyNameID = builder.propertyNameID
	b.extraContIdx = builder.extraContIdx
	b.numConts = builder.numConts
	b.remKind = remKind
	b.flags = builder.flags
}

// GetRemoveUnusedCode_RemovableBase gives access to the base.
func (b *RemoveUnusedCode_RemovableBase) GetRemoveUnusedCode_RemovableBase() *RemoveUnusedCode_RemovableBase {
	return b
}

// Java: RemoveUnusedCode.Removable.getPropertyName()
func (b *RemoveUnusedCode_RemovableBase) GetPropertyName() string {
	if b.remKind == remKindIndirectAssign {
		jutil.CheckStateMsg1(b.targetNode.IsGetProp(), "%s", b.targetNode)
		return b.targetNode.GetString()
	}
	if b.flags&remFlagHasPropertyName == 0 {
		panic(jutil.NewNullPointerException())
	}
	return rhino.RhinoStringPool_StringFromID(b.propertyNameID)
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
	if b.remKind == remKindIndirectAssign {
		return b.targetNode.IsGetProp()
	}
	return b.flags&remFlagHasPropertyName != 0
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
	if !b.IsNamedPropertyAssignment() {
		return false
	}
	return b.GetPropertyName() == "prototype"
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
	if b.remKind == remKindAnonProtoPropAssign {
		return true
	}
	return b.IsPrototypeDotPropertyReference() || b.IsClassOrPrototypeNamedProperty()
}

// Java: RemoveUnusedCode.Removable.isThisDotPropertyReference()
func (b *RemoveUnusedCode_RemovableBase) IsThisDotPropertyReference() bool {
	if b.remKind == remKindIndirectAssign {
		return RemoveUnusedCode_IsThisDotProperty(b.targetNode)
	}
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
			undefinedNode := NodeUtil_NewUndefinedNode(b.targetNode)
			r.ReplaceNodeWith(b.targetNode, undefinedNode)
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
			if name, ok := ast.MustClass(b.targetNode).Name(); ok {
				name.Node().ReplaceWith(rhino.IR_Empty().Srcref(name.Node()))
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
		jutil.CheckStateMsg1(objDotPrototype.Property() == "prototype", "%s", objDotPrototype.Node())

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
		} else {
			if parent.IsComma() || parent.IsAnd() || parent.IsOr() {
				if parent.GetFirstChild() == b.targetNode {
					rhs := rhino.CheckNotNull(b.targetNode.GetNext())
					compiler.ReportChangeToEnclosingScope(parent)
					parent.ReplaceWith(rhs.Detach())
				} else {
					lhs := parent.GetFirstChild()
					compiler.ReportChangeToEnclosingScope(parent)
					parent.ReplaceWith(lhs.Detach())
				}
			} else {
				compiler.ReportChangeToEnclosingScope(parent)
				b.targetNode.ReplaceWith(rhino.IR_Number(0))
			}
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

// RemoveUnusedCode_RemovableBuilder is the inner class RemovableBuilder.
type RemoveUnusedCode_RemovableBuilder struct {
	outer          *RemoveUnusedCode
	cont0          RemoveUnusedCode_Continuation
	propertyNameID uint32
	extraContIdx   uint32
	numConts       uint16
	flags          uint8
}

// Java: RemoveUnusedCode.RemovableBuilder()
func NewRemoveUnusedCode_RemovableBuilder(outer *RemoveUnusedCode) *RemoveUnusedCode_RemovableBuilder {
	return &RemoveUnusedCode_RemovableBuilder{outer: outer}
}

// Java: RemoveUnusedCode.RemovableBuilder.addContinuation(Continuation)
func (b *RemoveUnusedCode_RemovableBuilder) AddContinuation(continuation RemoveUnusedCode_Continuation) *RemoveUnusedCode_RemovableBuilder {
	switch b.numConts {
	case 0:
		b.cont0 = continuation
	case 1:
		b.extraContIdx = uint32(len(b.outer.extraContSlab))
		b.outer.extraContSlab = append(b.outer.extraContSlab, continuation)
	default:
		b.outer.extraContSlab = append(b.outer.extraContSlab, continuation)
	}
	b.numConts++
	return b
}

// Java: RemoveUnusedCode.RemovableBuilder.setIsPrototypeDotPropertyReference(boolean)
func (b *RemoveUnusedCode_RemovableBuilder) SetIsPrototypeDotPropertyReference(value bool) *RemoveUnusedCode_RemovableBuilder {
	if value {
		b.flags |= remFlagProtoDotPropRef
	} else {
		b.flags &^= remFlagProtoDotPropRef
	}
	return b
}

// Java: RemoveUnusedCode.RemovableBuilder.setIsThisDotPropertyReference(boolean)
func (b *RemoveUnusedCode_RemovableBuilder) SetIsThisDotPropertyReference(value bool) *RemoveUnusedCode_RemovableBuilder {
	if value {
		b.flags |= remFlagThisDotPropRef
	} else {
		b.flags &^= remFlagThisDotPropRef
	}
	return b
}

func (b *RemoveUnusedCode_RemovableBuilder) setPropertyName(name string) {
	b.propertyNameID = uint32(rhino.RhinoStringPool_InternID(name))
	b.flags |= remFlagHasPropertyName
}

func (b *RemoveUnusedCode_RemovableBuilder) setPropertyNameNode(n rhino.Node) {
	id := n.GetStringID()
	if id == 0 {
		s := n.GetString()
		if s != "" {
			id = uint32(rhino.RhinoStringPool_InternID(s))
		}
	}
	b.propertyNameID = id
	b.flags |= remFlagHasPropertyName
}

// Java: RemoveUnusedCode.RemovableBuilder.buildIndirectAssign(Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildIndirectAssign(root rhino.Node, targetNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_IndirectAssign(b.outer, b, root, targetNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildPolyfill(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildPolyfill(polyfillNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_Polyfill(b.outer, b, polyfillNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassDeclaration(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassDeclaration(classNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_ClassDeclaration(b.outer, b, classNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNamedClassExpression(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNamedClassExpression(classNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_NamedClassExpression(b.outer, b, classNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassOrPrototypeNamedProperty(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassOrPrototypeNamedProperty(propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	if !(propertyNode.IsMemberFunctionDef() ||
		propertyNode.IsMemberFieldDef() ||
		NodeUtil_IsGetOrSetKey(propertyNode) ||
		(propertyNode.IsStringKey() && !propertyNode.IsQuotedStringKey())) {
		jutil.CheckArgumentMsg1(false, "%s", propertyNode)
	}
	b.setPropertyNameNode(propertyNode)
	return NewRemoveUnusedCode_ClassOrPrototypeNamedProperty(b.outer, b, propertyNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildObjectDefinePropertiesDefinition(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildObjectDefinePropertiesDefinition(propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	b.setPropertyNameNode(propertyNode)
	return NewRemoveUnusedCode_ObjectDefinePropertiesDefinition(b.outer, b, propertyNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildFunctionDeclaration(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildFunctionDeclaration(functionNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_FunctionDeclaration(b.outer, b, functionNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNameDeclarationStatement(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNameDeclarationStatement(declarationStatement rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_NameDeclarationStatement(b.outer, b, declarationStatement)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNamedPropertyAssign(Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNamedPropertyAssign(assignNode rhino.Node, propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return b.BuildNamedPropertyAssignWithVarInfo(assignNode, propertyNode, nil)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildNamedPropertyAssign(Node, Node, VarInfo)
func (b *RemoveUnusedCode_RemovableBuilder) BuildNamedPropertyAssignWithVarInfo(assignNode rhino.Node, propertyNode rhino.Node, varInfo *RemoveUnusedCode_RealVarInfo) *RemoveUnusedCode_RemovableBase {
	b.setPropertyNameNode(propertyNode)
	return NewRemoveUnusedCode_Assign(b.outer, b, assignNode, RemoveUnusedCode_Kind_NAMED_PROPERTY, propertyNode, varInfo)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildComputedPropertyAssign(Node, Node, VarInfo)
func (b *RemoveUnusedCode_RemovableBuilder) BuildComputedPropertyAssign(assignNode rhino.Node, propertyNode rhino.Node, varInfo *RemoveUnusedCode_RealVarInfo) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_Assign(b.outer, b, assignNode, RemoveUnusedCode_Kind_COMPUTED_PROPERTY, propertyNode, varInfo)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildVariableAssign(Node, VarInfo)
func (b *RemoveUnusedCode_RemovableBuilder) BuildVariableAssign(assignNode rhino.Node, varInfo *RemoveUnusedCode_RealVarInfo) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_Assign(b.outer, b, assignNode, RemoveUnusedCode_Kind_VARIABLE /* propertyNode= */, rhino.NilNode, varInfo)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassSetupCall(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassSetupCall(callNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return b.BuildClassSetupCallWithClassDefiningCall(callNode /* classDefiningCall= */, false)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildClassSetupCall(Node, boolean)
func (b *RemoveUnusedCode_RemovableBuilder) BuildClassSetupCallWithClassDefiningCall(callNode rhino.Node, classDefiningCall bool) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_ClassSetupCall(b.outer, b, callNode, classDefiningCall)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildVanillaForNameDeclaration(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildVanillaForNameDeclaration(nameNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_VanillaForNameDeclaration(b.outer, b, nameNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildAnonymousPrototypeNamedPropertyAssign(Node, String)
func (b *RemoveUnusedCode_RemovableBuilder) BuildAnonymousPrototypeNamedPropertyAssign(
	assignNode rhino.Node, propertyName string) *RemoveUnusedCode_RemovableBase {
	b.setPropertyName(propertyName)
	return NewRemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(b.outer, b, assignNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildIncOrDepOp(Node, Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildIncOrDepOp(incOrDecOp rhino.Node, propertyNode rhino.Node, toPreseve rhino.Node) *RemoveUnusedCode_RemovableBase {
	b.setPropertyNameNode(propertyNode)
	return NewRemoveUnusedCode_IncOrDecOp(b.outer, b, incOrDecOp, toPreseve)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildUnusedReadReference(Node, Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildUnusedReadReference(referenceNode rhino.Node, propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	b.setPropertyNameNode(propertyNode)
	return NewRemoveUnusedCode_UnusedReadReference(b.outer, b, referenceNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildInstanceofName(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildInstanceofName(instanceofNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_InstanceofName(b.outer, b, instanceofNode)
}

// Java: RemoveUnusedCode.RemovableBuilder.buildWeakUsageCall(Node)
func (b *RemoveUnusedCode_RemovableBuilder) BuildWeakUsageCall(weakUsageCall rhino.Node) *RemoveUnusedCode_RemovableBase {
	return NewRemoveUnusedCode_WeakUsageCall(b.outer, b, weakUsageCall)
}

// Java: RemoveUnusedCode.UnusedReadReference(RemovableBuilder, Node)
func NewRemoveUnusedCode_UnusedReadReference(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, referenceNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	if !(RemoveUnusedCode_IsThisDotProperty(referenceNode) || RemoveUnusedCode_IsDotPrototypeDotProperty(referenceNode)) {
		jutil.CheckStateMsg1(false, "%s", referenceNode)
	}
	u := outer.removables.New()
	u.initFromBuilder(remKindUnusedReadRef, referenceNode, builder)
	return u
}

// Java: RemoveUnusedCode.InstanceofName(RemovableBuilder, Node)
func NewRemoveUnusedCode_InstanceofName(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, instanceofNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	if !instanceofNode.IsInstanceOf() {
		jutil.CheckArgumentMsg1(false, "%s", instanceofNode)
	}
	i := outer.removables.New()
	i.initFromBuilder(remKindInstanceofName, instanceofNode, builder)
	return i
}

// Java: RemoveUnusedCode.WeakUsageCall(RemovableBuilder, Node)
func NewRemoveUnusedCode_WeakUsageCall(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, weakCallNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	jutil.CheckArgumentMsg1(weakCallNode.IsCall(), "%s", weakCallNode)
	w := outer.removables.New()
	w.initFromBuilder(remKindWeakUsageCall, weakCallNode, builder)
	return w
}

// Java: RemoveUnusedCode.IncOrDecOp(RemovableBuilder, Node, Node)
func NewRemoveUnusedCode_IncOrDecOp(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, incOrDecNode rhino.Node, toPreserve rhino.Node) *RemoveUnusedCode_RemovableBase {
	incDec, ok := ast.AsIncDec(incOrDecNode)
	jutil.CheckArgumentMsg1(ok, "%s", incOrDecNode)
	arg := incDec.Target().Node()
	jutil.CheckStateMsg1(RemoveUnusedCode_IsThisDotProperty(arg) || RemoveUnusedCode_IsDotPrototypeDotProperty(arg), "%s", arg)
	o := outer.removables.New()
	o.initFromBuilder(remKindIncOrDecOp, incOrDecNode, builder)
	o.auxNode = toPreserve
	return o
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

// Java: RemoveUnusedCode.IndirectAssign(RemovableBuilder, Node, Node)
func NewRemoveUnusedCode_IndirectAssign(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, root rhino.Node, targetNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	rootParent := root.GetParent()
	if !(rootParent.IsDestructuringPattern() || rootParent.IsParamList()) {
		jutil.CheckStateMsg1(false, "%s", rootParent)
	}
	if !(targetNode.IsName() || targetNode.IsGetProp()) {
		jutil.CheckStateMsg1(false, "%s", targetNode)
	}
	a := outer.removables.New()
	a.initFromBuilder(remKindIndirectAssign, targetNode, builder)
	a.auxNode = root
	return a
}

// Java: RemoveUnusedCode.Polyfill(RemovableBuilder, Node)
func NewRemoveUnusedCode_Polyfill(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, polyfillNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	p := outer.removables.New()
	p.initFromBuilder(remKindPolyfill, polyfillNode, builder)
	// The callee: $jscomp.patch or (collapsed) $jscomp$patch.
	callee := ast.MustExprStmt(polyfillNode).Expr().MustCall().Callee()
	if strings.Contains(callee.Node().GetString(), "patch") {
		p.flags |= remFlagSubBit
	}
	return p
}

// Java: RemoveUnusedCode.ClassDeclaration(RemovableBuilder, Node)
func NewRemoveUnusedCode_ClassDeclaration(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, classDeclarationNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	c := outer.removables.New()
	c.initFromBuilder(remKindClassDecl, classDeclarationNode, builder)
	return c
}

// Java: RemoveUnusedCode.NamedClassExpression(RemovableBuilder, Node)
func NewRemoveUnusedCode_NamedClassExpression(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, classNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	c := outer.removables.New()
	c.initFromBuilder(remKindNamedClassExpr, classNode, builder)
	return c
}

// Java: RemoveUnusedCode.ClassOrPrototypeNamedProperty(RemovableBuilder, Node)
func NewRemoveUnusedCode_ClassOrPrototypeNamedProperty(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	p := outer.removables.New()
	p.initFromBuilder(remKindClassOrProtoProp, propertyNode, builder)
	return p
}

// Java: RemoveUnusedCode.ObjectDefinePropertiesDefinition(RemovableBuilder, Node)
func NewRemoveUnusedCode_ObjectDefinePropertiesDefinition(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, propertyNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	d := outer.removables.New()
	d.initFromBuilder(remKindObjectDefinePropsDef, propertyNode, builder)
	return d
}

// Java: RemoveUnusedCode.FunctionDeclaration(RemovableBuilder, Node)
func NewRemoveUnusedCode_FunctionDeclaration(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, functionDeclarationNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	f := outer.removables.New()
	f.initFromBuilder(remKindFunctionDecl, functionDeclarationNode, builder)
	return f
}

// Java: RemoveUnusedCode.NameDeclarationStatement(RemovableBuilder, Node)
func NewRemoveUnusedCode_NameDeclarationStatement(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, declarationStatement rhino.Node) *RemoveUnusedCode_RemovableBase {
	if !NodeUtil_IsNameDeclaration(declarationStatement) {
		jutil.CheckArgumentMsg1(false, "%s", declarationStatement)
	}
	s := outer.removables.New()
	s.initFromBuilder(remKindNameDeclStmt, declarationStatement, builder)
	if _, ok := s.declaredName().Init(); !ok {
		s.flags |= remFlagUninitVarDecl
	}
	return s
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

// RemoveUnusedCode_Kind is the nested enum Kind.
type RemoveUnusedCode_Kind int32

const (
	// X = something;
	RemoveUnusedCode_Kind_VARIABLE RemoveUnusedCode_Kind = iota
	// X.propertyName = something;
	// X.prototype.propertyName = something;
	RemoveUnusedCode_Kind_NAMED_PROPERTY
	// X[expression] = something;
	// X.prototype[expression] = something;
	RemoveUnusedCode_Kind_COMPUTED_PROPERTY
)

var removeUnusedCode_Kind_names = [...]string{"VARIABLE", "NAMED_PROPERTY", "COMPUTED_PROPERTY"}

// Java: RemoveUnusedCode.Kind.name()
func (k RemoveUnusedCode_Kind) String() string { return removeUnusedCode_Kind_names[k] }

// Java: RemoveUnusedCode.Assign(RemovableBuilder, Node, Kind, Node, VarInfo)
func NewRemoveUnusedCode_Assign(
	outer *RemoveUnusedCode,
	builder *RemoveUnusedCode_RemovableBuilder,
	assignNode rhino.Node,
	kind RemoveUnusedCode_Kind,
	propertyNode rhino.Node,
	varInfo *RemoveUnusedCode_RealVarInfo) *RemoveUnusedCode_RemovableBase {
	if !NodeUtil_IsAssignmentOp(assignNode) {
		jutil.CheckArgumentMsg1(false, "%s", assignNode)
	}
	var rk uint8
	if kind == RemoveUnusedCode_Kind_VARIABLE {
		if !propertyNode.IsNil() {
			jutil.CheckArgumentMsg1(
				false,
				"got property node for simple variable assignment: %s",
				propertyNode)
		}
		if varInfo == nil {
			jutil.CheckArgumentMsg1(false, "missing VarInfo for variable assignment: %s", propertyNode)
		}
		rk = remKindAssignVar
	} else {
		if propertyNode.IsNil() {
			jutil.CheckArgumentMsg(false, "missing property node")
		}
		if kind == RemoveUnusedCode_Kind_NAMED_PROPERTY {
			if !propertyNode.IsGetProp() {
				jutil.CheckArgumentMsg1(
					false, "property name is not a GETPROP: %s", propertyNode)
			}
			rk = remKindAssignNamedProp
		} else {
			rk = remKindAssignComputedProp
		}
	}
	a := outer.removables.New()
	a.initFromBuilder(rk, assignNode, builder)
	a.varInfo = varInfo
	return a
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

// Java: RemoveUnusedCode.AnonymousPrototypeNamedPropertyAssign(RemovableBuilder, Node)
func NewRemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, assignNode rhino.Node) *RemoveUnusedCode_RemovableBase {
	if builder.flags&remFlagHasPropertyName == 0 {
		panic(jutil.NewNullPointerException())
	}
	jutil.CheckArgumentMsg1(assignNode.IsAssign(), "%s", assignNode)
	a := outer.removables.New()
	a.initFromBuilder(remKindAnonProtoPropAssign, assignNode, builder)
	return a
}

// Java: RemoveUnusedCode.ClassSetupCall(RemovableBuilder, Node, boolean)
func NewRemoveUnusedCode_ClassSetupCall(outer *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, callNode rhino.Node, classDefiningCall bool) *RemoveUnusedCode_RemovableBase {
	c := outer.removables.New()
	c.initFromBuilder(remKindClassSetupCall, callNode, builder)
	if classDefiningCall {
		c.flags |= remFlagSubBit
	}
	return c
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
