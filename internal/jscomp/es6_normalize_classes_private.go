// Port of com/google/javascript/jscomp/Es6NormalizeClasses.java @ closure-compiler v20261021

/*
 * Copyright 2025 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InsertPrivateDeclarations inserts declarations for the PrivateMap objects
// before the class statement if the class contains private members.
//
// Java: Es6NormalizeClasses.insertPrivateDeclarations(NodeTraversal, ClassRecord, Node)
func (e *Es6NormalizeClasses) InsertPrivateDeclarations(
	t *NodeTraversal, record *es6NormalizeClasses_ClassRecord, classInsertionPoint rhino.Node) {
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, e.compiler)
	// Create const PRIVATE_MAP$uid = new $jscomp.PrivateMap() declaration for instance private
	// members.
	if record.hasInstancePrivateMembers {
		e.compiler.GetRuntimeJsLibManager().InjectLibForField("$jscomp.PrivateMap")
		mapVar :=
			e.astFactory.
				CreateSingleConstNameDeclaration(
					record.GetPrivateMapVarName(e.compiler, t),
					e.astFactory.CreateNewNode(e.astFactory.CreateQNameWithUnknownTypeOfField(e.jscompPrivateMap))).
				SrcrefTreeIfMissing(record.classNode.Node())
		mapVar.InsertBefore(classInsertionPoint)

		// Create const PRIVATE_PROTO$uid = Object.create(null, ...) declaration for instance private
		// methods and accessors.
		if record.hasInstancePrivateMethodsOrAccessors {
			// Non-static private instance methods and accessors are gathered into descriptor
			// definitions for the shared private prototype descriptor (PRIVATE_PROTO).
			descriptorsLit := e.astFactory.CreateObjectLit()
			for _, pm := range record.privateMembers {
				// Process instance (non-static) private members.
				if !pm.isStatic {
					// Process private instance methods.
					if fn, ok := es6NormalizeClasses_unmovedFunction(pm.methodNode.Node()); ok &&
						pm.kind == es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD {
						e.RewriteSuperReferencesInPrivateMember(fn, pm.methodNode.Node().IsStaticMember(), record)
						fn.Detach()
						pm.methodNode.Detach()
						valProp := e.astFactory.CreateStringKey("value", fn.Node())
						descriptorLit := e.astFactory.CreateObjectLit(valProp)
						memberKey := e.astFactory.CreateStringKey(pm.PropertyName(), descriptorLit)
						descriptorsLit.AddChildToBack(memberKey)
					} else if pm.kind == es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR {
						// Process private instance accessors (getters/setters).
						descriptorLit := e.astFactory.CreateObjectLit()
						if getFn, ok := es6NormalizeClasses_unmovedFunction(pm.getterNode.Node()); ok {
							e.RewriteSuperReferencesInPrivateMember(getFn, pm.getterNode.Node().IsStaticMember(), record)
							getFn.Detach()
							pm.getterNode.Detach()
							if e.RewriteThisToThisDotSelf(getFn) {
								record.needsInstanceSelf = true
							}
							getProp := e.astFactory.CreateStringKey("get", getFn.Node())
							descriptorLit.AddChildToBack(getProp)
						}
						if setFn, ok := es6NormalizeClasses_unmovedFunction(pm.setterNode.Node()); ok {
							e.RewriteSuperReferencesInPrivateMember(setFn, pm.setterNode.Node().IsStaticMember(), record)
							setFn.Detach()
							pm.setterNode.Detach()
							if e.RewriteThisToThisDotSelf(setFn) {
								record.needsInstanceSelf = true
							}
							setProp := e.astFactory.CreateStringKey("set", setFn.Node())
							descriptorLit.AddChildToBack(setProp)
						}
						memberKey := e.astFactory.CreateStringKey(pm.PropertyName(), descriptorLit)
						descriptorsLit.AddChildToBack(memberKey)
					}
				}
			}

			// Generates: const PRIVATE_PROTO$uid = Object.create(null, { method: { value: function(...)
			// {} } });
			createProp :=
				e.astFactory.CreateGetPropWithUnknownType(
					e.astFactory.CreateNameWithUnknownType("Object"), "create")
			protoCall :=
				e.astFactory.CreateCallWithUnknownType(
					createProp, e.astFactory.CreateNull(), descriptorsLit)
			protoVar :=
				e.astFactory.
					CreateSingleConstNameDeclaration(
						record.GetPrivateProtoVarName(e.compiler, t), protoCall).
					SrcrefTreeIfMissing(record.classNode.Node())
			protoVar.InsertBefore(classInsertionPoint)
			e.compiler.ReportChangeToEnclosingScope(protoVar)
			t.ReportCodeChange()
		}
	}
	// Create const STATIC_PRIVATE_MAP$uid = new $jscomp.PrivateMap() declaration for static private
	// members.
	if record.hasStaticPrivateMembers {
		e.compiler.GetRuntimeJsLibManager().InjectLibForField("$jscomp.PrivateMap")
		staticMapVar :=
			e.astFactory.
				CreateSingleConstNameDeclaration(
					record.GetStaticPrivateMapVarName(e.compiler, t),
					e.astFactory.CreateNewNode(e.astFactory.CreateQNameWithUnknownTypeOfField(e.jscompPrivateMap))).
				SrcrefTreeIfMissing(record.classNode.Node())
		staticMapVar.InsertBefore(classInsertionPoint)
	}
}

// FindPrivateMember finds the private member declaration in the enclosing
// class stack by its name (including the '#' prefix).
//
// Java: Es6NormalizeClasses.findPrivateMember(String)
func (e *Es6NormalizeClasses) FindPrivateMember(name string) (*es6NormalizeClasses_ClassRecord, *es6NormalizeClasses_ClassRecord_PrivateMember, bool) {
	for i := len(e.classStack) - 1; i >= 0; i-- {
		record := e.classStack[i]
		if pm := record.GetPrivateMember(name); pm != nil {
			return record, pm, true
		}
	}
	return nil, nil, false
}

// RewritePrivateIn rewrites private identifier brand checks (e.g.
// `#field in obj`) to use `PRIVATE_MAP.has(obj)`.
//
// Java: Es6NormalizeClasses.rewritePrivateIn(NodeTraversal, Node)
func (e *Es6NormalizeClasses) RewritePrivateIn(t *NodeTraversal, in ast.Binary) {
	inNode := in.Node()
	jutil.CheckArgumentMsg1(inNode.IsIn(), "%s", inNode)
	// The parser marks a name on the left of `in` as private: `#field in obj`.
	left, right := ast.MustName(in.Left().Node()), in.Right()
	record, pm, ok := e.FindPrivateMember(left.Ident())
	if !ok {
		e.MaybeReportPrivatePropertiesCannotConvertYet(left.Node())
		return
	}

	right.Detach()
	var mapVarName string
	if pm.isStatic {
		mapVarName = record.GetStaticPrivateMapVarName(e.compiler, t)
	} else {
		mapVarName = record.GetPrivateMapVarName(e.compiler, t)
	}

	mapNameNode := e.astFactory.CreateNameWithUnknownType(mapVarName)
	hasCallee := e.astFactory.CreateGetPropWithUnknownType(mapNameNode, "has")
	hasCall := e.astFactory.CreateCallWithUnknownType(hasCallee, right.Node())

	inNode.ReplaceWith(hasCall.SrcrefTree(inNode))
	t.ReportCodeChange()
}

// RewritePrivateGetProp rewrites private property accesses (e.g.
// `obj.#field`) to use `PRIVATE_MAP.get(obj).field`.
//
// For non-static private methods:
//
//   - Direct calls: `obj.#method(args)` are rewritten to
//     `PRIVATE_MAP.get(obj).method.call(obj, args)`.
//   - Method access (tear-off): `obj.#method` is rewritten to
//     `PRIVATE_MAP.get(obj).method`.
//
// TODO(b/236744850): Consider memoizing PRIVATE_MAP.get(obj) into a local
// temporary variable when multiple private accesses occur on the same
// receiver within a single statement or method scope to reduce WeakMap
// lookup overhead.
//
// Java: Es6NormalizeClasses.rewritePrivateGetProp(NodeTraversal, Node)
func (e *Es6NormalizeClasses) RewritePrivateGetProp(t *NodeTraversal, getProp ast.NormalOrOptChainGetProp) {
	getPropNode := getProp.Node()
	record, pm, ok := e.FindPrivateMember(getProp.Property())
	if !ok {
		e.MaybeReportPrivatePropertiesCannotConvertYet(getPropNode)
		return
	}

	// Private methods are read-only and cannot be reassigned (e.g., `this.#method = fn`).
	// Report an error if a private method is used as an L-value.
	if NodeUtil_IsLValue(getPropNode) && pm.kind == es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD {
		e.compiler.Report(
			JSError_MakeWithNode(getPropNode, Es6NormalizeClasses_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT, pm.PropertyName()))
		return
	}

	e.RewritePrivateMemberAccess(t, getProp, record, pm, getPropNode.GetParent())
}

// RewritePrivateMemberAccess rewrites a private member access or
// invocation.
//
// Handles three primary cases:
//
//  1. Direct member calls with side-effecting receivers (e.g.
//     `getObj()?.#method(a)` or `getObj()?.#field(a)`):
//
//     let TMP;
//     (TMP = getObj(), TMP == null ? void 0 : PRIVATE_MAP.get(TMP).member.call(TMP, a))
//
//  2. Direct member calls with side-effect-free receivers (e.g.
//     `this.#method(a)` or `obj?.#field?.(a)`):
//
//     obj == null ? void 0 : PRIVATE_MAP.get(obj).member.call(obj, a)
//
//  3. Indirect reads, field accesses, and method tear-offs (e.g.
//     `const fn = obj?.#method` or `obj?.#field`):
//
//     const fn = obj == null ? void 0 : PRIVATE_MAP.get(obj).member;
//
// Java: Es6NormalizeClasses.rewritePrivateMemberAccess(NodeTraversal, Node, ClassRecord, ClassRecord.PrivateMember, Node)
func (e *Es6NormalizeClasses) RewritePrivateMemberAccess(
	t *NodeTraversal,
	getProp ast.NormalOrOptChainGetProp,
	record *es6NormalizeClasses_ClassRecord,
	pm *es6NormalizeClasses_ClassRecord_PrivateMember,
	parent rhino.Node) {
	getPropNode := getProp.Node()
	var mapVarName string
	if pm.isStatic {
		mapVarName = record.GetStaticPrivateMapVarName(e.compiler, t)
	} else {
		mapVarName = record.GetPrivateMapVarName(e.compiler, t)
	}
	memberPropName := pm.PropertyName()

	receiver := getProp.Object().Detach().Node()
	mapNameNode := e.astFactory.CreateNameWithUnknownType(mapVarName)
	getCallee := e.astFactory.CreateGetPropWithUnknownType(mapNameNode, "get")

	_, isOptChain := getProp.AsOptChainGetProp()
	var endOfSegment rhino.Node
	if isOptChain {
		endOfSegment = NodeUtil_GetEndOfOptChainSegment(getPropNode)
	}

	call, isCall := ast.AsNormalOrOptChainCall(parent)
	isDirectCall := isCall && call.Callee().Is(getPropNode)

	if isDirectCall {
		_, isOptChainCall := call.AsOptChainCall()
		isOptionalCall := isOptChainCall && parent.IsOptionalChainStart()

		parent.SetIsOptionalChainStart(false)

		if e.ReceiverNeedsTempVar(receiver) {
			// If evaluating receiver twice could introduce side effects or alter execution state,
			// declare a temporary variable (e.g. `let TMP;`) and evaluate the call using a comma
			// expression:
			// `(TMP = receiver, PRIVATE_MAP.get(TMP).member.call(TMP, ...args))`
			tempName := e.GenerateUniqueTmpVarName(t)
			insertionBlock := e.GetDeclarationInsertionPoint(getPropNode, t)
			decl := e.astFactory.CreateSingleLetNameDeclaration(tempName)
			insertionBlock.AddChildToFront(decl.SrcrefTree(getPropNode))
			NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)

			tempNameNode1 := e.astFactory.CreateNameWithUnknownType(tempName).Srcref(getPropNode)
			assign := e.astFactory.CreateAssign(tempNameNode1, receiver).Srcref(getPropNode)
			receiver = e.astFactory.CreateNameWithUnknownType(tempName).Srcref(getPropNode)

			tempNameNode2 := e.astFactory.CreateNameWithUnknownType(tempName).Srcref(getPropNode)
			tempNameNode3 := e.astFactory.CreateNameWithUnknownType(tempName).Srcref(getPropNode)

			getCall :=
				e.astFactory.CreateCallWithUnknownType(getCallee, tempNameNode2).Srcref(getPropNode)
			methodGetProp := e.astFactory.CreateGetPropWithUnknownType(getCall, memberPropName)
			callCallee := e.CreateCallCallee(methodGetProp, isOptionalCall, getPropNode)
			if callCallee.IsOptChainGetProp() {
				parent.SetToken(rhino.Token_OPTCHAIN_CALL)
			} else {
				parent.SetToken(rhino.Token_CALL)
			}

			callCallee.SrcrefTreeIfMissing(getPropNode)
			tempNameNode3.SrcrefTree(getPropNode)

			getPropNode.Detach()
			parent.AddChildToFront(callCallee)
			tempNameNode3.InsertAfter(callCallee)

			resultExpr :=
				e.WrapOptionalChainStartIfNecessary(receiver, parent, endOfSegment, isOptChain)
			placeholder := rhino.IR_Empty()
			resultExpr.ReplaceWith(placeholder)
			comma := e.astFactory.CreateComma(assign, resultExpr).Srcref(getPropNode)
			placeholder.ReplaceWith(comma)

			t.ReportCodeChange()
		} else {
			// Direct private member call with a side-effect-free receiver (e.g. `this` or a simple
			// variable):
			// Rewritten to `PRIVATE_MAP.get(receiver).member.call(receiver, ...args)`.
			getCall := e.astFactory.CreateCallWithUnknownType(getCallee, receiver)
			methodGetProp := e.astFactory.CreateGetPropWithUnknownType(getCall, memberPropName)
			receiverArg := receiver.CloneTree()
			receiverArg.SrcrefTree(getPropNode)

			getPropNode.Detach()
			callCallee := e.CreateCallCallee(methodGetProp, isOptionalCall, getPropNode)
			if callCallee.IsOptChainGetProp() {
				parent.SetToken(rhino.Token_OPTCHAIN_CALL)
			} else {
				parent.SetToken(rhino.Token_CALL)
			}
			callCallee.SrcrefTreeIfMissing(getPropNode)
			parent.AddChildToFront(callCallee)
			receiverArg.InsertAfter(callCallee)
			e.WrapOptionalChainStartIfNecessary(receiver, parent, endOfSegment, isOptChain)
			t.ReportCodeChange()
		}
	} else {
		// Indirect read, field access, or method tear-off (e.g., `const fn = this.#method` or
		// `this.#field`):
		// Rewritten to `PRIVATE_MAP.get(receiver).propertyName`.
		assignComma := e.ExtractReceiverTmpIfNecessary(t, receiver, getPropNode, isOptChain)
		if !assignComma.IsNil() {
			receiver = assignComma.GetFirstChild().CloneTree()
		}
		getCall := e.astFactory.CreateCallWithUnknownType(getCallee, receiver)
		methodGetProp := e.astFactory.CreateGetPropWithUnknownType(getCall, memberPropName)
		getPropNode.ReplaceWith(methodGetProp.SrcrefTree(getPropNode))
		result :=
			e.WrapOptionalChainStartIfNecessary(receiver, methodGetProp, endOfSegment, isOptChain)

		if !assignComma.IsNil() {
			placeholder := rhino.IR_Empty()
			result.ReplaceWith(placeholder)
			comma := e.astFactory.CreateComma(assignComma, result).Srcref(getPropNode)
			placeholder.ReplaceWith(comma)
		}

		t.ReportCodeChange()
	}
}

// ExtractReceiverTmpIfNecessary extracts a temporary variable declaration
// for a side-effecting receiver in an optional chain.
//
// Java: Es6NormalizeClasses.extractReceiverTmpIfNecessary(NodeTraversal, Node, Node, boolean)
func (e *Es6NormalizeClasses) ExtractReceiverTmpIfNecessary(
	t *NodeTraversal, receiver rhino.Node, getPropNode rhino.Node, isOptChain bool) rhino.Node {
	if !isOptChain || !e.ReceiverNeedsTempVar(receiver) {
		return rhino.NilNode
	}
	tempName := e.GenerateUniqueTmpVarName(t)
	insertionBlock := e.GetDeclarationInsertionPoint(getPropNode, t)
	decl := e.astFactory.CreateSingleLetNameDeclaration(tempName)
	insertionBlock.AddChildToFront(decl.SrcrefTree(getPropNode))
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)

	tempNameNode := e.astFactory.CreateNameWithUnknownType(tempName).Srcref(getPropNode)
	return e.astFactory.CreateAssign(tempNameNode, receiver).Srcref(getPropNode)
}

// WrapOptionalChainStartIfNecessary wraps an optional chained private
// member access with a nullish short-circuit check on the receiver,
// climbing to the end of the current optional chain segment.
//
// Java: Es6NormalizeClasses.wrapOptionalChainStartIfNecessary(Node, Node, Node, boolean)
func (e *Es6NormalizeClasses) WrapOptionalChainStartIfNecessary(
	receiver rhino.Node, accessExpr rhino.Node, endOfSegment rhino.Node, isOptChain bool) rhino.Node {
	if !isOptChain {
		return accessExpr
	}
	var segmentRoot rhino.Node
	if !endOfSegment.IsNil() && !endOfSegment.GetParent().IsNil() {
		segmentRoot = endOfSegment
	} else {
		segmentRoot = accessExpr
	}
	if segmentRoot != accessExpr {
		NodeUtil_ConvertToNonOptionalChainSegmentDownTo(segmentRoot, accessExpr)
	}

	placeholder := rhino.IR_Empty()
	if !segmentRoot.GetParent().IsNil() {
		segmentRoot.ReplaceWith(placeholder)
	}
	// We must short-circuit before calling PRIVATE_MAP.get(receiver), since PRIVATE_MAP.get throws
	// on nullish receivers.
	hook :=
		e.astFactory.
			CreateOptionalChainShortCircuit(receiver.CloneTree(), segmentRoot).
			SrcrefTreeIfMissing(segmentRoot)
	if !placeholder.GetParent().IsNil() {
		placeholder.ReplaceWith(hook)
	}
	hookParent := hook.GetParent()
	if !hookParent.IsNil() && (hookParent.IsCall() || hookParent.IsOptChainCall()) {
		hook.SetIsParenthesized(true)
		if hookParent.GetFirstChild() == hook {
			hookParent.PutBooleanProp(rhino.Node_FREE_CALL, true)
		}
	}
	return hook
}

// CreateCallCallee creates the `.call` property access node for a private
// method invocation.
//
// Java: Es6NormalizeClasses.createCallCallee(Node, boolean, Node)
func (e *Es6NormalizeClasses) CreateCallCallee(methodGetProp rhino.Node, isOptionalCall bool, getPropNode rhino.Node) rhino.Node {
	if isOptionalCall {
		return e.astFactory.CreateStartOptChainGetprop(methodGetProp, "call", AstFactory_Type_(getPropNode))
	}
	return e.astFactory.CreateGetPropWithUnknownType(methodGetProp, "call")
}

func es6NonClassOrNonArrowFnPred(n rhino.Node) bool {
	return !n.IsClass() && (!n.IsFunction() || n.IsArrowFunction())
}

// RewriteThisToThisDotSelf rewrites top-level `this` references inside a
// private accessor body to `this.$self`.
//
// Returns true if any `this` reference was rewritten.
//
// Java: Es6NormalizeClasses.rewriteThisToThisDotSelf(Node)
func (e *Es6NormalizeClasses) RewriteThisToThisDotSelf(fn ast.Function) bool {
	return e.rewriteThisToThisDotSelfInSubtree(fn.Body().Node())
}

func (e *Es6NormalizeClasses) rewriteThisToThisDotSelfInSubtree(node rhino.Node) bool {
	hasThis := false
	if es6NonClassOrNonArrowFnPred(node) {
		for c := node.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			if e.rewriteThisToThisDotSelfInSubtree(c) {
				hasThis = true
			}
			c = next
		}
	}
	if node.IsThis() {
		hasThis = true
		target :=
			e.astFactory.
				CreateGetPropWithUnknownType(
					e.astFactory.CreateThis(AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)), "$self").
				Srcref(node)
		node.ReplaceWith(target)
		e.compiler.ReportChangeToEnclosingScope(target)
	}
	return hasThis
}

// Java: Es6NormalizeClasses.hasThisReference(Node)
func (e *Es6NormalizeClasses) HasThisReference(fn ast.Function) bool {
	return es6SubtreeHasThis(fn.Body().Node())
}

func es6SubtreeHasThis(node rhino.Node) bool {
	if node.IsThis() {
		return true
	}
	if es6NonClassOrNonArrowFnPred(node) {
		for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if es6SubtreeHasThis(c) {
				return true
			}
		}
	}
	return false
}

// RewriteSuperReferencesInPrivateMember rewrites `super` references inside
// private methods or accessors into explicit target calls or Reflect calls.
//
// Java: Es6NormalizeClasses.rewriteSuperReferencesInPrivateMember(Node, ClassRecord)
func (e *Es6NormalizeClasses) RewriteSuperReferencesInPrivateMember(fn ast.Function, isStatic bool, record *es6NormalizeClasses_ClassRecord) {
	if record.normalizedNames.SuperClassNode().IsNil() {
		return
	}
	fnBody := fn.Body().Node()
	if !NodeUtil_ReferencesSuper(fnBody) {
		return
	}
	e.rewriteSuperRefsInSubtree(fnBody, fnBody, isStatic, record)
}

func (e *Es6NormalizeClasses) rewriteSuperRefsInSubtree(
	node rhino.Node, fnBody rhino.Node, isStatic bool, record *es6NormalizeClasses_ClassRecord) {
	if !node.IsClass() {
		for c := node.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			e.rewriteSuperRefsInSubtree(c, fnBody, isStatic, record)
			c = next
		}
	}
	if !node.IsSuper() {
		return
	}
	parent := node.GetParent()
	get, ok := ast.AsNormalGet(parent)
	if !ok {
		return
	}
	var target rhino.Node
	if isStatic {
		target = record.CreateNewSuperclassReferenceNode().SrcrefTree(node)
	} else {
		target =
			e.astFactory.
				CreatePrototypeAccess(record.CreateNewSuperclassReferenceNode()).
				SrcrefTree(node)
	}

	grandParent := parent.GetParent()
	if call, ok := ast.AsCall(grandParent); ok && call.Callee().Is(parent) {
		// Method call: super.method(...)
		node.ReplaceWith(target)
		if !isStatic || !e.compiler.GetOptions().GetAssumeStaticInheritanceIsNotUsed() {
			callCallee :=
				e.astFactory.CreateGetPropWithUnknownType(parent.Detach(), "call")
			grandParent.AddChildToFront(callCallee.SrcrefTreeIfMissing(parent))
			thisArg := e.astFactory.CreateThis(AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)).Srcref(node)
			thisArg.InsertAfter(callCallee)
		}
	} else if assign, ok := ast.AsAssign(grandParent); ok && assign.Target().Is(parent) {
		// Setter write: super.prop = value
		rhs := assign.Value().Detach().Node()
		propKey := e.superPropertyKey(get, record, node)

		reflectSet :=
			e.astFactory.CreateCall(
				e.astFactory.CreateGetPropWithUnknownType(
					e.astFactory.CreateName("Reflect", AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)), "set"),
				AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN),
				target,
				propKey,
				rhs,
				e.astFactory.CreateThis(AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)).Srcref(node))
		grandParent.ReplaceWith(reflectSet.SrcrefTreeIfMissing(grandParent))
	} else {
		// Getter read: super.prop
		propKey := e.superPropertyKey(get, record, node)

		reflectGet :=
			e.astFactory.CreateCall(
				e.astFactory.CreateGetPropWithUnknownType(
					e.astFactory.CreateName("Reflect", AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)), "get"),
				AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN),
				target,
				propKey,
				e.astFactory.CreateThis(AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)).Srcref(node))
		parent.ReplaceWith(reflectGet.SrcrefTreeIfMissing(parent))
	}
	e.compiler.ReportChangeToEnclosingScope(fnBody)
}

func (e *Es6NormalizeClasses) superPropertyKey(
	get ast.NormalGet, record *es6NormalizeClasses_ClassRecord, superNode rhino.Node) rhino.Node {
	if getProp, ok := get.AsGetProp(); ok {
		return e.astFactory.CreateCall(
			e.astFactory.CreateName(
				NodeUtil_JSC_PROPERTY_NAME_FN, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)),
			AstFactory_Type_OfColor(colors.StandardColors_STRING),
			e.astFactory.CreateString(getProp.Property()),
			record.CreateNewSuperclassReferenceNode().SrcrefTree(superNode))
	}
	return ast.MustGetElem(get.Node()).Index().Detach().Node()
}

// AddStaticPrivateMembers attaches static private methods and accessors to
// the static private object.
//
// Java: Es6NormalizeClasses.addStaticPrivateMembers(ClassRecord, String, Node)
func (e *Es6NormalizeClasses) AddStaticPrivateMembers(
	record *es6NormalizeClasses_ClassRecord, privName string, staticInitTempParent rhino.Node) {
	for _, pm := range record.privateMembers {
		// Only process static private members for static initialization.
		if pm.isStatic {
			// Assign static private method definitions to the static private object (e.g.,
			// PRIVATE$uid.method = function() {}).
			if fn, ok := es6NormalizeClasses_unmovedFunction(pm.methodNode.Node()); ok &&
				pm.kind == es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD {
				e.RewriteSuperReferencesInPrivateMember(fn, pm.methodNode.Node().IsStaticMember(), record)
				fn.Detach()
				pm.methodNode.Detach()
				lhs :=
					e.astFactory.CreateGetPropWithUnknownType(
						e.astFactory.CreateConstantName(privName, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)),
						pm.PropertyName())
				assign := e.astFactory.CreateAssignStatement(lhs, fn.Node())
				assign.SrcrefTreeIfMissing(pm.methodNode.Node())
				staticInitTempParent.AddChildToBack(assign)
			} else if pm.kind == es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR {
				// Assign static private accessor definitions to the static private object via
				// Object.defineProperty.
				descriptorLit := e.astFactory.CreateObjectLit()
				if getFn, ok := es6NormalizeClasses_unmovedFunction(pm.getterNode.Node()); ok {
					e.RewriteSuperReferencesInPrivateMember(getFn, pm.getterNode.Node().IsStaticMember(), record)
					getFn.Detach()
					pm.getterNode.Detach()
					e.RewriteThisToThisDotSelf(getFn)
					descriptorLit.AddChildToBack(e.astFactory.CreateStringKey("get", getFn.Node()))
				}
				if setFn, ok := es6NormalizeClasses_unmovedFunction(pm.setterNode.Node()); ok {
					e.RewriteSuperReferencesInPrivateMember(setFn, pm.setterNode.Node().IsStaticMember(), record)
					setFn.Detach()
					pm.setterNode.Detach()
					e.RewriteThisToThisDotSelf(setFn)
					descriptorLit.AddChildToBack(e.astFactory.CreateStringKey("set", setFn.Node()))
				}
				defineProp :=
					e.astFactory.CreateGetPropWithUnknownType(
						e.astFactory.CreateNameWithUnknownType("Object"), "defineProperty")
				call :=
					e.astFactory.CreateCallWithUnknownType(
						defineProp,
						e.astFactory.CreateConstantName(privName, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)),
						e.astFactory.CreateString(pm.PropertyName()),
						descriptorLit)
				declNode := pm.GetPrimaryDeclarationNode()
				exprResult := e.astFactory.ExprResult(call)
				if !declNode.IsNil() {
					exprResult.SrcrefTreeIfMissing(declNode)
				}
				staticInitTempParent.AddChildToBack(exprResult)
			}
		}
	}
}

// MaybeReportPrivatePropertiesCannotConvertYet reports that transpilation
// of private class properties is not yet supported when attempting to
// rewrite or normalize private identifiers. See b/236744850.
//
// Java: Es6NormalizeClasses.maybeReportPrivatePropertiesCannotConvertYet(Node)
func (e *Es6NormalizeClasses) MaybeReportPrivatePropertiesCannotConvertYet(n rhino.Node) {
	if e.transpilePrivateClassFields && n.IsPrivateIdentifier() {
		TranspilationUtil_CannotConvertYet(e.compiler, n, "Private elements")
	}
}

// ReceiverNeedsTempVar returns true if evaluating the receiver expression
// twice could introduce side effects or alter state. Side-effect-free
// receivers like `this` or simple variable names do not need a temp
// variable.
//
// Java: Es6NormalizeClasses.receiverNeedsTempVar(Node)
func (e *Es6NormalizeClasses) ReceiverNeedsTempVar(receiver rhino.Node) bool {
	return !receiver.IsThis() && !receiver.IsName()
}

// GetDeclarationInsertionPoint finds the enclosing block or script where
// temporary variable declarations (e.g. `let $jscomp$tmp$...`) can be
// safely inserted.
//
// Java: Es6NormalizeClasses.getDeclarationInsertionPoint(Node, NodeTraversal)
func (e *Es6NormalizeClasses) GetDeclarationInsertionPoint(n rhino.Node, t *NodeTraversal) rhino.Node {
	enclosingStmt := NodeUtil_GetEnclosingStatement(n)
	if !enclosingStmt.IsNil() && !enclosingStmt.GetParent().IsNil() {
		return enclosingStmt.GetParent()
	}
	return t.GetScopeRoot()
}

// Es6NormalizeClasses_IsPrivateMember returns true if the node represents a
// private member declaration (field, method, or accessor) with a private
// identifier starting with '#'.
//
// Java: Es6NormalizeClasses.isPrivateMember(Node)
func Es6NormalizeClasses_IsPrivateMember(n rhino.Node) bool {
	return n.IsPrivateIdentifier() || (n.HasChildren() && n.GetFirstChild().IsPrivateIdentifier())
}
