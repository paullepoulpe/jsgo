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

func (r *RemoveUnusedCode) traverseNameReference(n rhino.Node, scope *Scope) {
	// The only cases that should reach this point are parameter declarations and references
	// to names. The name node does not have children in these cases.
	if n.HasChildren() {
		jutil.CheckState(false)
	}
	// (A switch on the parent's token: this runs for every name, and
	// RoleOf would read the parent's children too.)
	switch n.GetParent().GetToken() {
	case rhino.Token_PARAM_LIST:
		// the parameter declaration is not a read of the name
		return
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		// var|let|const name;
		// are handled at a higher level.
		jutil.CheckState(false)
	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		// function name() {}
		// class name() {}
		// are handled at a higher level.
		if role := ast.RoleOf(n); role == ast.RoleFunctionName || role == ast.RoleClassName {
			jutil.CheckState(false)
		}
	}
	varInfo := r.TraverseNameNode(n, scope)
	if varInfo.IsRemovable() {
		varInfo.pinAt(rucPinReferenceFound, n)
	}
}

// TraverseNode traverses everything in the current scope and marks
// variables that are referenced.
//
// During traversal, we identify subtrees that will only be referenced if
// their enclosing variables are referenced. Instead of traversing those
// subtrees, we create a continuation for them, and traverse them lazily.
//
// Java: RemoveUnusedCode.traverseNode(Node, Scope)
func (r *RemoveUnusedCode) TraverseNode(n rhino.Node, scope *Scope) {
	if n.GetFirstChild().IsNil() {
		if n.GetToken() == rhino.Token_NAME {
			r.traverseNameReference(n, scope)
		}
		return
	}
	switch n.GetToken() {
	case rhino.Token_CATCH:
		r.TraverseCatch(n, scope)
	case rhino.Token_FUNCTION:
		// If this function is a removable var, then create a continuation
		// for it instead of traversing immediately.
		if NodeUtil_IsFunctionDeclaration(n) {
			name, _ := ast.MustFunction(n).Name() // a declaration has one
			varInfo := r.TraverseNameNode(name.Node(), scope)
			varInfo.AddRemovable(r.newFunctionDeclRemovable(n, scope))
			if n.GetParent().IsExport() {
				varInfo.pin(rucPinExportedClass)
			}
		} else {
			r.TraverseFunction(n, scope)
		}
	case rhino.Token_ASSIGN:
		r.TraverseAssign(n, scope)
	case rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD:
		r.TraverseCompoundAssign(n, scope)
	case rhino.Token_INC, rhino.Token_DEC:
		r.TraverseIncrementOrDecrementOp(n, scope)
	case rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL:
		r.TraverseCall(n, scope)
	case rhino.Token_SWITCH_BODY:
		if nodeUtilBlockScopeMayHaveDecls(n, n) {
			r.TraverseChildren(n, r.scopeCreator.CreateScopeOfScope(n, scope))
		} else {
			r.TraverseChildren(n, scope)
		}
	case rhino.Token_BLOCK:
		// This case if for if there are let and const variables in block scopes.
		// Otherwise other variables will be hoisted up into the global scope and already be
		// handled.
		if scanVarsBlockCreatesScope(n) &&
			(NodeUtil_IsFunctionBlock(n) || NodeUtil_IsClassStaticBlock(n) || syntacticScopeHasBlockDecls(n)) {
			r.TraverseChildren(n, r.scopeCreator.CreateScopeOfScope(n, scope))
		} else {
			r.TraverseChildren(n, scope)
		}
	case rhino.Token_MODULE_BODY:
		r.TraverseChildren(n, r.scopeCreator.CreateScopeOfScope(n, scope))
	case rhino.Token_CLASS:
		r.TraverseClass(n, scope)
	case rhino.Token_CLASS_MEMBERS:
		r.TraverseClassMembers(n, scope)
	case rhino.Token_ARRAY_PATTERN, rhino.Token_PARAM_LIST:
		r.TraverseIndirectAssignmentList(n, scope)
	case rhino.Token_OBJECT_PATTERN:
		r.TraverseObjectPattern(n, scope)
	case rhino.Token_OBJECTLIT:
		r.TraverseObjectLiteral(n, scope)
	case rhino.Token_FOR:
		r.TraverseVanillaFor(n, scope)
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		r.TraverseEnhancedFor(n, scope)
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		// for-loop cases are handled by custom traversal methods.
		if !NodeUtil_IsStatement(n) {
			jutil.CheckState(false)
		}
		r.TraverseDeclarationStatement(n, scope)
	case rhino.Token_INSTANCEOF:
		r.TraverseInstanceof(n, scope)
	case rhino.Token_NAME:
		r.traverseNameReference(n, scope)
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		r.TraverseNormalOrOptChainGetProp(n, scope)
	default:
		r.TraverseChildren(n, scope)
	}
}

// Java: RemoveUnusedCode.traverseInstanceof(Node, Scope)
func (r *RemoveUnusedCode) TraverseInstanceof(instanceofNode rhino.Node, scope *Scope) {
	instanceof, ok := ast.AsBinaryOf(instanceofNode, rhino.Token_INSTANCEOF)
	if !ok {
		jutil.CheckArgumentMsg1(false, "%s", instanceofNode)
	}
	r.TraverseNode(instanceof.Left().Node(), scope)
	if rhs, ok := instanceof.Right().AsName(); ok {
		varInfo := r.TraverseNameNode(rhs.Node(), scope)
		varInfo.AddRemovable(r.newInstanceofNameRemovable(instanceofNode))
	} else {
		r.TraverseNode(instanceof.Right().Node(), scope)
	}
}

// TraverseIncrementOrDecrementOp.
//
// TODO(b/137380742): Combine with `traverseCompoundAssign`.
//
// Java: RemoveUnusedCode.traverseIncrementOrDecrementOp(Node, Scope)
func (r *RemoveUnusedCode) TraverseIncrementOrDecrementOp(incOrDecOp rhino.Node, scope *Scope) {
	incDec, ok := ast.AsIncDec(incOrDecOp)
	jutil.CheckArgumentMsg1(ok, "%s", incOrDecOp)
	arg := incDec.Target().Node()
	if NodeUtil_IsExpressionResultUsed(incOrDecOp) {
		// If expression result is used, then this expression is definitely not removable.
		r.TraverseNode(arg, scope)
	} else if getProp, ok := incDec.Target().AsGetProp(); ok {
		getPropObj := getProp.Object().Node()

		if r.ConsiderForAccessorSideEffects(arg, AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER) {
			r.TraverseNode(getPropObj, scope) // Don't re-traverse the GETPROP as a read.
		} else if getPropObj.IsThis() {
			// this.propName++
			r.ConsiderForIndependentRemoval(r.newIncOrDecOpRemovable(incOrDecOp, arg, rhino.NilNode, remFlagThisDotPropRef))
		} else if owner, ok := removeUnusedCode_prototypeOwner(getPropObj); ok {
			// someExpression.prototype.propName++
			exprObj := owner.Node()
			if exprObj.IsName() {
				// varName.prototype.propName++
				varInfo := r.TraverseNameNode(exprObj, scope)
				varInfo.AddRemovable(r.newIncOrDecOpRemovable(incOrDecOp, arg, rhino.NilNode, remFlagProtoDotPropRef))
			} else {
				// (someExpression).prototype.propName++
				var toPreserve rhino.Node
				var rem *RemoveUnusedCode_RemovableBase
				if r.MayHaveSideEffects(exprObj) {
					toPreserve = exprObj
					r.TraverseNode(exprObj, scope)
					rem = r.newIncOrDecOpRemovable(incOrDecOp, arg, toPreserve, remFlagProtoDotPropRef)
				} else {
					rem = r.newIncOrDecOpRemovable(incOrDecOp, arg, rhino.NilNode, remFlagProtoDotPropRef)
					rem.addContinuation(r, exprObj, scope)
				}
				r.ConsiderForIndependentRemoval(rem)
			}
		} else {
			// someExpression.propName++ is not removable except in the cases covered above
			r.TraverseNode(arg, scope)
		}
	} else {
		// TODO(bradfordcsmith): varName++ should be removable if varName is otherwise unused
		r.TraverseNode(arg, scope)
	}
}

// TraverseCompoundAssign.
//
// TODO(b/137380742): Combine with `traverseIncrementOrDecrement`.
//
// Java: RemoveUnusedCode.traverseCompoundAssign(Node, Scope)
func (r *RemoveUnusedCode) TraverseCompoundAssign(compoundAssignNode rhino.Node, scope *Scope) {
	// We'll allow removal of compound assignment to a `this` property as long as the result of the
	// assignment is unused.
	// e.g. `this.prop += 3;`

	// NOTE: Some history here, as there were questions about "why is 'this' special".  The "remove
	// unused properties" is not a general property removal algorithm.  It only removes unreferenced
	// properties that are part of class definitions.  "SomeClass.prop += 3" and
	// "SomeClass.prototype.prop += 3" could so be candidates but they aren't considered here.

	assign := ast.MustAssignOp(compoundAssignNode)
	targetNode := assign.Target().Node()
	valueNode := assign.Value().Node()
	if getProp, ok := assign.Target().AsGetProp(); ok {
		if r.ConsiderForAccessorSideEffects(targetNode, AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER) {
			r.TraverseNode(getProp.Object().Node(), scope) // Don't re-traverse the GETPROP as a read.
			r.TraverseNode(valueNode, scope)
		} else if getProp.Object().Token() == rhino.Token_THIS &&
			!NodeUtil_IsExpressionResultUsed(compoundAssignNode) {
			rem := r.newNamedPropAssignRemovable(compoundAssignNode, targetNode, nil, remFlagThisDotPropRef)
			r.traverseRemovableAssignValue(rem, valueNode, scope)
			r.ConsiderForIndependentRemoval(rem)
		} else {
			r.TraverseNode(targetNode, scope)
			r.TraverseNode(valueNode, scope)
		}
	} else {
		r.TraverseNode(targetNode, scope)
		r.TraverseNode(valueNode, scope)
	}
}

// Java: RemoveUnusedCode.traverseNameNode(Node, Scope)
func (r *RemoveUnusedCode) TraverseNameNode(n rhino.Node, scope *Scope) *RemoveUnusedCode_RealVarInfo {
	nameID := n.GetStringID()
	if r.polyBits.Get(int(nameID)) {
		r.considerPolyfillRef(nameID, n)
	}
	v := scope.GetVarByID(nameID)
	if v == nil {
		jutil.CheckNotNullMsg1(v, "%s", n)
	}
	return r.TraverseVar(v)
}

func isSpecialCallCandidate(callee rhino.Node) bool {
	switch callee.GetToken() {
	case rhino.Token_GETPROP:
		switch callee.GetStringID() {
		case rucCalleeInheritsID, rucCalleeMixinID, rucCalleeDefinePropertiesID,
			rucCalleeObjectPropertyID, rucCalleeAddSingletonGetterID, rucCalleePolyfillID,
			rucCalleePatchID, rucCalleePolyfillTypedArrayMethodID, rucCalleeWeakUsageID:
			return true
		}
		return false
	case rhino.Token_NAME:
		switch callee.GetStringID() {
		case rucCalleeJscPropNameFnID,
			rucCalleeGoogReflectObjectPropertyID,
			rucCalleeJscompPolyfillID,
			rucCalleeJscompPatchID,
			rucCalleeJscompPolyfillTypedArrayMethodID,
			rucCalleeGoogWeakUsageID,
			rucCalleeGoogAddSingletonGetterID:
			return true
		}
		name := callee.GetString()
		return strings.HasSuffix(name, "$inherits") || strings.HasSuffix(name, "$mixin")
	default:
		return false
	}
}

// Java: RemoveUnusedCode.traverseCall(Node, Scope)
func (r *RemoveUnusedCode) TraverseCall(callNode rhino.Node, scope *Scope) {
	call := ast.MustNormalOrOptChainCall(callNode)
	callee := call.Callee().Node()
	if !isSpecialCallCandidate(callee) {
		r.TraverseChildren(callNode, scope)
		return
	}

	if r.codingConvention.IsPropertyRenameFunction(callee) {
		if firstArg, ok := call.FirstArg(); ok && firstArg.Node().IsStringLit() {
			r.markPropertyIDPinned(firstArg.Node().GetStringID())
		}
		r.TraverseChildren(callNode, scope)
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(callNode) {
		// TODO(bradfordcsmith): Should also handle Object.create() and Object.defineProperty().
		r.TraverseObjectDefinePropertiesCall(callNode, scope)
	} else if r.removeUnusedPolyfills && RemoveUnusedCode_IsJscompPolyfill(callee) {
		firstArg, _ := call.FirstArg() // isJscompPolyfill: a string
		polyfillName := firstArg.MustStringLit().Value()
		if calleeID := callee.GetStringID(); calleeID == rucCalleePolyfillTypedArrayMethodID ||
			calleeID == rucCalleeJscompPolyfillTypedArrayMethodID {
			polyfillName = "TypedArray.prototype." + polyfillName
		}
		info := r.createPolyfillInfo(call, scope, polyfillName)
		r.polyBits.Set(int(info.keyID))
		idx := r.ensurePropEntry(info.keyID)
		entry := &r.propEntries[idx]
		if entry.polyHead == nil {
			entry.polyHead = info
			entry.polyTail = info
		} else {
			entry.polyTail.nextByKey = info
			entry.polyTail = info
		}
		r.polyfillOrder = append(r.polyfillOrder, info)
		// Only traverse the callee (to mark it as used).  The arguments may be traversed later.
		r.TraverseNode(callee, scope)
	} else if weakName, ok := removeUnusedCode_onlyNameArg(call); ok && NodeUtil_IsGoogWeakUsageCall(callNode) {
		// goog.weakUsage() should have exactly one argument, and it should be either a name or a
		// qualified name (this is checked in ProcessClosurePrimitives.java).
		//
		// If it is a qualified name, then we do not attempt to remove it at this time (see condition
		// above). We rely on CollapseProperties to turn qualified names into simple names where
		// possible.

		// Mark this call as removable if the var is not otherwise referenced.
		varInfo := r.TraverseNameNode(weakName.Node(), scope)
		varInfo.AddRemovable(r.newWeakUsageCallRemovable(callNode))

		// We need to traverse the goog.weakUsage function itself (to mark it as used, in case our
		// usage of it is not removed).
		r.TraverseNode(callee, scope)
	} else {
		parent := callNode.GetParent()
		var classVarName string
		hasClassVarName := false
		classDefiningCall := false

		// A call that is a statement unto itself or the left side of a comma expression might be
		// a call to a known method for doing class setup
		// e.g. $jscomp.inherits(Class, BaseClass) or goog.addSingletonGetter(Class)
		// Such methods never have meaningful return values, so we won't look for them in other
		// contexts
		if parent.IsExprResult() || (parent.IsComma() && ast.RoleOf(callNode) == ast.RoleBinaryLeft) {
			subclassRelationship :=
				r.codingConvention.GetClassesDefinedByCall(callNode)
			if subclassRelationship != nil {
				// e.g. goog.inherits(DerivedClass, BaseClass);
				// NOTE: DerivedClass and BaseClass must be QNames. Otherwise getClassesDefinedByCall()
				// will return null.
				classVarName = subclassRelationship.subclassName
				hasClassVarName = true
				classDefiningCall = true
			} else if singletonName := r.codingConvention.GetSingletonGetterClassName(callNode); singletonName != nil {
				// Look for calls to addSingletonGetter calls.
				classVarName = *singletonName
				hasClassVarName = true
			}
		}

		var classVar *Var
		if hasClassVarName && NodeUtil_IsValidSimpleName(classVarName) {
			if id := rhino.RhinoStringPool_LookupID(classVarName); id != 0 {
				classVar = scope.GetVarByID(id)
			}
			if classVar == nil {
				jutil.CheckNotNullMsg1(classVar, "%s", classVarName)
			}
		}

		if classVar == nil || !classVar.IsGlobal() {
			// The call we are traversing does not modify a class definition,
			// or the class is not specified with a simple variable name,
			// or the variable name is not global.
			// TODO(bradfordcsmith): It would be more correct to check whether the class name
			// references a known constructor and expand to allow QNames.
			r.TraverseChildren(callNode, scope)
		} else {
			rem := r.newClassSetupCallRemovable(callNode, classDefiningCall)
			for child := callNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
				rem.addContinuation(r, child, scope)
			}
			r.TraverseVar(classVar).AddRemovable(rem)
		}
	}
}

// removeUnusedCode_onlyNameArg returns the argument of a call that has
// exactly one, a name.
func removeUnusedCode_onlyNameArg(call ast.NormalOrOptChainCall) (ast.Name, bool) {
	arg, ok := call.OnlyArg()
	if !ok {
		return ast.Name{}, false
	}
	return arg.AsName()
}

// RemoveUnusedCode_IsJscompPolyfill checks whether this is a recognizable
// call to $jscomp.polyfill.
//
// Java: RemoveUnusedCode.isJscompPolyfill(Node)
func RemoveUnusedCode_IsJscompPolyfill(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		// Need to work correctly after CollapseProperties.
		id := n.GetStringID()
		return (id == rucCalleeJscompPolyfillID ||
			id == rucCalleeJscompPatchID ||
			id == rucCalleeJscompPolyfillTypedArrayMethodID) &&
			n.GetNext().IsStringLit()
	case rhino.Token_GETPROP:
		// Need to work correctly without CollapseProperties.
		id := n.GetStringID()
		if id != rucCalleePolyfillID && id != rucCalleePatchID && id != rucCalleePolyfillTypedArrayMethodID {
			return false
		}
		object, ok := ast.MustGetProp(n).Object().AsName()
		return ok && object.Node().GetStringID() == rucJscompIdentID && n.GetNext().IsStringLit()
	default:
		return false
	}
}

// Java: RemoveUnusedCode.traverseCatch(Node, Scope)
func (r *RemoveUnusedCode) TraverseCatch(catchNode rhino.Node, scope *Scope) {
	catch := ast.MustCatch(catchNode)
	// The binding is missing (an EMPTY node) in catch {} (ES2019).
	if param, ok := catch.Param(); ok {
		if exceptionName, ok := param.AsName(); ok {
			exceptionVarInfo := r.TraverseNameNode(exceptionName.Node(), scope)
			exceptionVarInfo.pin(rucPinCatchVariable)
		} else {
			r.TraverseNode(param.Node(), scope)
		}
	}
	r.TraverseNode(catch.Body().Node(), scope)
}

// Java: RemoveUnusedCode.traverseEnhancedFor(Node, Scope)
func (r *RemoveUnusedCode) TraverseEnhancedFor(enhancedFor rhino.Node, scope *Scope) {
	forScope := scope
	if syntacticScopeHasBlockDecls(enhancedFor) {
		forScope = r.scopeCreator.CreateScopeOfScope(enhancedFor, scope)
	}
	// for (iterationTarget in|of collection) body;
	loop := ast.MustEnhancedFor(enhancedFor)
	iterationTarget := loop.Target().Node()
	if name, ok := loop.Target().AsName(); ok {
		// using previously-declared loop variable. e.g.
		// `for (varName of collection) {}`
		varInfo := r.TraverseNameNode(name.Node(), forScope)
		varInfo.pin(rucPinEnhancedForVar)
	} else if decl, ok := loop.Target().AsNameDeclaration(); ok {
		// loop has const/var/let declaration
		d, ok := decl.OnlyDeclarator()
		jutil.CheckState(ok)
		declNode := d.Node()
		if name, ok := d.AsName(); !ok {
			// e.g.
			// `for (const [a, b] of pairList) {}`
			// destructuring is handled at a lower level
			// Note that destructuring assignments are always considered to set an unknown value
			// equivalent to what we set for the var name case above and below.
			// It isn't necessary to set the variable names as not removable, though, because the
			// thing that isn't removable is the destructuring pattern itself, which we never remove.
			// TODO(bradfordcsmith): The need to explain all the above shows this should be reworked.
			r.TraverseNode(declNode, forScope)
		} else {
			// e.g.
			// `for (const varName of collection) {}`
			_, hasInit := name.Init()
			jutil.CheckState(!hasInit)
			// We can never remove the loop variable of a for-in or for-of loop, because it's
			// essential to loop syntax.
			varInfo := r.TraverseNameNode(name.Node(), forScope)
			varInfo.pin(rucPinEnhancedForVar)
		}
	} else {
		// using some general LHS value e.g.
		// `for ([a, b] of collection) {}` destructuring with existing vars
		// `for (a.x of collection) {}` using a property as the loop var
		// TODO(bradfordcsmith): This should be considered a write if it's a property reference.
		r.TraverseNode(iterationTarget, forScope)
	}
	r.TraverseNode(loop.Collection().Node(), forScope)
	r.TraverseNode(loop.AsLoop().Body().Node(), forScope)
}

// Java: RemoveUnusedCode.traverseVanillaFor(Node, Scope)
func (r *RemoveUnusedCode) TraverseVanillaFor(forNode rhino.Node, scope *Scope) {
	forScope := scope
	if syntacticScopeHasBlockDecls(forNode) {
		forScope = r.scopeCreator.CreateScopeOfScope(forNode, scope)
	}
	// The missing parts of a for (an EMPTY node each) hold nothing to traverse.
	loop := ast.MustFor(forNode)
	if initialization, ok := loop.Init(); ok {
		if decl, ok := initialization.AsNameDeclaration(); ok {
			r.TraverseVanillaForNameDeclarations(decl.Node(), forScope)
		} else {
			r.TraverseNode(initialization.Node(), forScope)
		}
	}
	if condition, ok := loop.Cond(); ok {
		r.TraverseNode(condition.Node(), forScope)
	}
	if update, ok := loop.Update(); ok {
		r.TraverseNode(update.Node(), forScope)
	}
	r.TraverseNode(loop.Body().Node(), forScope)
}

// Java: RemoveUnusedCode.traverseVanillaForNameDeclarations(Node, Scope)
func (r *RemoveUnusedCode) TraverseVanillaForNameDeclarations(nameDeclaration rhino.Node, scope *Scope) {
	for child := range ast.MustNameDeclaration(nameDeclaration).Declarators() {
		if name, ok := child.AsName(); !ok {
			// TODO(bradfordcsmith): Customize handling of destructuring
			r.TraverseNode(child.Node(), scope)
		} else {
			nameNode := name.Node()
			value, hasValue := name.Init()
			valueNode := value.Node()
			varInfo := r.TraverseNameNode(nameNode, scope)
			if !hasValue {
				varInfo.AddRemovable(r.newVanillaForNameDeclRemovable(nameNode, false))
			} else if r.MayHaveSideEffects(valueNode) {
				// TODO(bradfordcsmith): Actually allow for removing the variable while keeping the
				// valueNode for its side-effects.
				varInfo.pin(rucPinForLoopInitSideEffects)
				r.TraverseNode(valueNode, scope)
			} else {
				rem := r.newVanillaForNameDeclRemovable(nameNode, true)
				rem.addContinuation(r, valueNode, scope)
				varInfo.AddRemovable(rem)
			}
		}
	}
}

// Java: RemoveUnusedCode.traverseDeclarationStatement(Node, Scope)
func (r *RemoveUnusedCode) TraverseDeclarationStatement(declarationStatement rhino.Node, scope *Scope) {
	// Normalization should ensure that declaration statements always have just one child.
	declarator, ok := ast.MustNameDeclaration(declarationStatement).OnlyDeclarator()
	jutil.CheckState(ok)
	if name, ok := declarator.AsName(); !ok {
		// Destructuring declarations are handled elsewhere.
		r.TraverseNode(declarator.Node(), scope)
	} else {
		value, hasValue := name.Init()
		valueNode := value.Node()
		varInfo := r.TraverseNameNode(name.Node(), scope)
		rem := r.newNameDeclStmtRemovable(declarationStatement, hasValue)
		if hasValue {
			if r.MayHaveSideEffects(valueNode) {
				r.TraverseNode(valueNode, scope)
			} else {
				rem.addContinuation(r, valueNode, scope)
			}
		}
		varInfo.AddRemovable(rem)
	}
}

// Java: RemoveUnusedCode.traverseIndirectAssignmentList(Node, Scope)
func (r *RemoveUnusedCode) TraverseIndirectAssignmentList(list rhino.Node, scope *Scope) {
	if !list.IsArrayPattern() && !list.IsParamList() {
		jutil.CheckStateMsg1(false, "%s", list)
	}

	for elem := list.GetFirstChild(); !elem.IsNil(); elem = elem.GetNext() {
		switch elem.GetToken() {
		case rhino.Token_EMPTY:
		case rhino.Token_ARRAY_PATTERN, rhino.Token_DEFAULT_VALUE, rhino.Token_GETELEM, rhino.Token_GETPROP, rhino.Token_NAME, rhino.Token_OBJECT_PATTERN:
			r.TraverseIndirectAssignment(elem, elem, scope)
		case rhino.Token_ITER_REST:
			r.TraverseIndirectAssignment(elem, ast.MustIterRest(elem).Target().Node(), scope)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of " + list.GetToken().String() + ": " + elem.ToStringTree()))
		}
	}
}

// TraverseIndirectAssignment traverses an AST structure representing an
// assignment operation for which the target and value are far apart.
//
// Examples include destructurings and function parameters.
//
// root is the root of the assignment subtree; target the l-value expression
// being assigned to.
//
// Java: RemoveUnusedCode.traverseIndirectAssignment(Node, Node, Scope)
func (r *RemoveUnusedCode) TraverseIndirectAssignment(root rhino.Node, target rhino.Node, scope *Scope) {
	rootParent := root.GetParent()
	if !rootParent.IsDestructuringPattern() && !rootParent.IsParamList() {
		jutil.CheckArgumentMsg1(false, "%s", rootParent)
	}

	// Flatten out the case where the target is a default value. We always have to consider it.
	if defaultValue, ok := ast.AsDefaultValue(target); ok {
		target = defaultValue.Target().Node()
	}

	if target.IsGetProp() {
		r.ConsiderForAccessorSideEffects(target, AccessorSummary_PropertyAccessKind_SETTER_ONLY)
	}

	if r.MayHaveSideEffects(root) {
		// If anywhere in the assignment subtree has side-effects, it means that even if the target is
		// removable the subtree is not.
		r.TraverseNode(root, scope)
		// TODO(bradfordcsmith): Preserve side effects without preventing removal of variables and
		// properties. We could probably do this by subbing in an empty object pattern.
	} else if target.IsName() {
		varInfo := r.TraverseNameNode(target, scope)
		varInfo.AddRemovable(r.newIndirectAssignRemovable(root, target, scope))
	} else if r.IsNameDotPrototype(target) || RemoveUnusedCode_IsThisDotProperty(target) {
		r.ConsiderForIndependentRemoval(r.newIndirectAssignRemovable(root, target, scope))
	} else {
		// TODO(bradfordcsmith): Handle property assignments also
		// e.g. `({a: foo.bar, b: foo.baz}) = {a: 1, b: 2}`
		r.TraverseNode(root, scope)
	}
}

// Java: RemoveUnusedCode.traverseChildren(Node, Scope)
func (r *RemoveUnusedCode) TraverseChildren(n rhino.Node, scope *Scope) {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.GetToken() == rhino.Token_NAME {
				r.traverseNameReference(c, scope)
			}
			continue
		}
		r.TraverseNode(c, scope)
	}
}

// TraverseFunction traverses a function
//
// ES6 scopes of a function include the parameter scope and the body scope of
// the function.
//
// Note that CATCH blocks also create a new scope, but only for the catch
// variable. Declarations within the block actually belong to the enclosing
// scope. Because we don't remove catch variables, there's no need to treat
// CATCH blocks differently like we do functions.
//
// Java: RemoveUnusedCode.traverseFunction(Node, Scope)
func (r *RemoveUnusedCode) TraverseFunction(function rhino.Node, parentScope *Scope) {
	fn, ok := ast.AsFunction(function)
	if !ok {
		jutil.CheckStateMsg1(false, "%s", function)
	}
	body := fn.Body().Node()
	if !body.IsBlock() {
		jutil.CheckStateMsg1(false, "%s", body)
	}
	paramlist := fn.Params()

	// Checking the parameters
	fparamScope := r.scopeCreator.CreateScopeOfScope(function, parentScope)

	// Checking the function body
	fbodyScope := r.scopeCreator.CreateScopeOfScope(body, fparamScope)

	if name, ok := fn.Name(); ok {
		// var x = function funcName() {};
		// make sure funcName gets into the varInfoMap so it will be considered for removal.
		varInfo := r.TraverseNameNode(name.Node(), fparamScope)
		if NodeUtil_IsExpressionResultUsed(function) {
			// var f = function g() {};
			// The f is an alias for g, so g escapes from the scope where it is defined.
			varInfo.SetHasNonLocalOrNonLiteralValue()
		}
	}

	if paramlist.HasParams() {
		r.TraverseIndirectAssignmentList(paramlist.Node(), fparamScope)
		if r.removeGlobals {
			r.allFunctionParamScopes = append(r.allFunctionParamScopes, fparamScope)
		}
	}
	r.TraverseChildren(body, fbodyScope)
}
