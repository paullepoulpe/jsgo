// Port of com/google/javascript/jscomp/PeepholeFoldConstants.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// TryFoldCtorCall tries to fold away unnecessary object instantiation. e.g.
// this[new String('eval')] -> this.eval
//
// Java: PeepholeFoldConstants.tryFoldCtorCall(Node)
func (p *PeepholeFoldConstants) TryFoldCtorCall(n rhino.Node) rhino.Node {
	jutil.CheckArgument(n.IsNew())

	// we can remove this for GETELEM calls (anywhere else?)
	if PeepholeFoldConstants_InForcedStringContext(n) {
		return p.TryFoldInForcedStringContext(n)
	}
	return n
}

// TryFoldUselessObjectDotDefinePropertiesCall removes useless calls:
// Object.defineProperties(o, {}) -> o
//
// Java: PeepholeFoldConstants.tryFoldUselessObjectDotDefinePropertiesCall(Node)
func (p *PeepholeFoldConstants) TryFoldUselessObjectDotDefinePropertiesCall(n rhino.Node) rhino.Node {
	call, ok := ast.AsNormalOrOptChainCall(n)
	jutil.CheckArgument(ok)

	// Go-only: the callee's property is compared first, a cheap test that
	// fails for almost every call, before NodeUtil's full check.
	callee, ok := call.Callee().AsGetProp()
	if !ok || callee.Node().GetStringID() != foldIDDefineProperties || !NodeUtil_IsObjectDefinePropertiesDefinition(n) {
		return n
	}
	// Object.defineProperties(destObj, srcObj): the check above proves both
	// arguments.
	destObj, _ := call.FirstArg()
	srcObj, _ := call.LastArg()
	if lit, ok := srcObj.AsObjectLit(); ok && !lit.HasMembers() {
		parent := n.GetParent()
		n.ReplaceWith(destObj.Node().Detach())
		p.ReportChangeToEnclosingScope(parent)
	}
	return n
}

// PeepholeFoldConstants_InForcedStringContext returns whether this node must
// be coerced to a string.
//
// Java: PeepholeFoldConstants.inForcedStringContext(Node)
func PeepholeFoldConstants_InForcedStringContext(n rhino.Node) bool {
	parent := n.GetParent()
	if getElem, ok := ast.AsGetElem(parent); ok && getElem.Index().Is(n) {
		return true
	}

	// we can fold in the case "" + new String("")
	return parent.IsAdd()
}

// Java: PeepholeFoldConstants.tryFoldInForcedStringContext(Node)
func (p *PeepholeFoldConstants) TryFoldInForcedStringContext(n rhino.Node) rhino.Node {
	// For now, we only know how to fold ctors.
	newNode, ok := ast.AsNew(n)
	jutil.CheckArgument(ok)

	objectType, ok := newNode.Callee().AsName()
	if !ok || objectType.Node().GetStringID() != foldIDString {
		return n
	}

	var stringValue string
	hasStringValue := true
	if value, ok := newNode.FirstArg(); ok {
		stringValue, hasStringValue = p.GetSideEffectFreeStringValue(value.Node())
	}
	// else: new String() -> ""

	if !hasStringValue {
		return n
	}

	parent := n.GetParent()
	newString := rhino.IR_String(stringValue)

	n.ReplaceWith(newString)
	newString.SrcrefIfMissing(parent)
	p.ReportChangeToEnclosingScope(parent)

	return newString
}

// TryFoldGetElem tries, for element access using GETLEM/OPTCHAIN_GETELEM on
// object literals, arrays or strings, to fold the prop access. e.g. folds
// array-element [1, 2, 3][1];
//
// Java: PeepholeFoldConstants.tryFoldGetElem(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldGetElem(n rhino.Node) rhino.Node {
	get, ok := ast.AsNormalOrOptChainGetElem(n)
	jutil.CheckArgument(ok)
	left, right := get.Object(), get.Index()

	if obj, ok := left.AsObjectLit(); ok {
		if key, ok := right.AsStringLit(); ok {
			return p.TryFoldObjectPropAccess(n, obj, key.Value())
		}
	} else if arr, ok := left.AsArrayLit(); ok {
		return p.TryFoldArrayAccess(n, arr, right.Node())
	} else if str, ok := left.AsStringLit(); ok {
		return p.TryFoldStringArrayAccess(n, str, right.Node())
	}
	return n
}

// TryFoldGetProp tries, for prop access using GETPROP/OPTCHAIN_GETPROP on
// object literals, to fold their property access. For prop access on arrays,
// only tries to fold array-length. e.g [1, 2, 3].length ==> 3, [x, y].length
// ==> 2
//
// Java: PeepholeFoldConstants.tryFoldGetProp(Node)
func (p *PeepholeFoldConstants) TryFoldGetProp(n rhino.Node) rhino.Node {
	get, ok := ast.AsNormalOrOptChainGetProp(n)
	jutil.CheckArgument(ok)
	left := get.Object()

	switch left.Token() {
	case rhino.Token_OBJECTLIT:
		obj, _ := left.AsObjectLit() // the case proves the kind
		return p.TryFoldObjectPropAccess(n, obj, get.Property())
	case rhino.Token_ARRAYLIT, rhino.Token_STRINGLIT:
		if get.Node().GetStringID() != foldIDLength || NodeUtil_IsLValue(n) {
			return n
		}
		var knownLength int
		if arr, ok := left.AsArrayLit(); ok {
			if p.MayHaveSideEffects(arr.Node()) {
				// Nope, can't fold this, without handling the side-effects.
				return n
			}
			for elem := range arr.Elements() {
				if _, ok := elem.AsIterSpread(); ok {
					return n
				}
			}
			knownLength = arr.Node().GetChildCount()
		} else if str, ok := left.AsStringLit(); ok {
			knownLength = jutil.StringLength(str.Value())
		}

		lengthNode := rhino.IR_Number(float64(knownLength))
		p.ReportChangeToEnclosingScope(n)
		n.ReplaceWith(lengthNode)

		return lengthNode
	}

	return n
}

// PeepholeFoldConstants_IsFoldableElementIndexKey returns true if right
// represents a valid key for folding element reads. For string keys, only
// canonical decimal array index spellings (e.g. "0", "1") match. For number
// keys, any number-typed expression is accepted so callers can emit warnings
// for non-integers or fold out-of-bounds index reads to undefined.
//
// Java: PeepholeFoldConstants.isFoldableElementIndexKey(Node, int)
func PeepholeFoldConstants_IsFoldableElementIndexKey(right rhino.Node, intIndex int32) bool {
	if right.IsStringLit() {
		return intIndex >= 0 && right.GetString() == strconv.Itoa(int(intIndex))
	}
	return NodeUtil_IsNumericResult(right)
}

// Java: PeepholeFoldConstants.tryFoldArrayAccess(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldArrayAccess(n rhino.Node, left ast.ArrayLit, right rhino.Node) rhino.Node {
	// If GETPROP/GETELEM is used as assignment target the array literal is
	// acting as a temporary we can't fold it here:
	//    "[][0] += 1"
	if NodeUtil_IsLValue(n) {
		return n
	}

	index, ok := p.GetSideEffectFreeNumberValue(right)
	if !p.isReasonableDoubleVal(index, ok) {
		// Sometimes people like to use complex expressions to index into
		// arrays, or strings to index into array methods.
		return n
	}

	intIndex := jutil.D2I(index)
	if !PeepholeFoldConstants_IsFoldableElementIndexKey(right, intIndex) {
		return n
	}

	if !base.JSCompDoubles_IsExactInt32(index) {
		// Ideally this should be caught in the check passes.
		p.Report(PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR, right)
		return n
	}
	var elem rhino.Node
	i := int32(0)
	for current := range left.Elements() {
		if _, ok := current.AsIterSpread(); ok {
			// The only time we can fold getelems with spread is for spread arrays literals, and
			// `tryFlattenArray` already flattens those.
			return n
		}
		if i != intIndex {
			if p.MayHaveSideEffects(current.Node()) {
				return n
			}
		} else {
			elem = current.Node()
		}
		i++
	}

	if elem.IsNil() { // If the index was out of bounds
		elem = NodeUtil_NewUndefinedNode(left.Node())
	} else if elem.IsEmpty() {
		elem = NodeUtil_NewUndefinedNode(elem)
	} else {
		elem.Detach()
	}

	// Replace the entire GETELEM with the value
	n.ReplaceWith(elem)
	p.ReportChangeToEnclosingScope(elem)
	return elem
}

// TryFoldSpread folds any occurrences of spread of array literals e.g.
// ...[1,2,3] to 1,2,3
//
// Java: PeepholeFoldConstants.tryFoldSpread(Node)
func (p *PeepholeFoldConstants) TryFoldSpread(spread rhino.Node) rhino.Node {
	// (Go: the Java checks isSpread(); OptimizeSubtree only calls this for
	// an ITER_SPREAD, in an array literal or the arguments of a call.)
	iterSpread, ok := ast.AsIterSpread(spread)
	jutil.CheckState(ok)
	parent := spread.GetParent()
	child, ok := iterSpread.Expr().AsArrayLit()
	if ok && PeepholeFoldConstants_IsSafeToFlattenArrayLitSpread(child.Node()) {
		parent.AddChildrenAfter(child.Node().RemoveChildren(), spread)
		spread.Detach()
		p.ReportChangeToEnclosingScope(parent)
	}
	return parent
}

// TryFlattenArrayOrObjectLit flattens array- or object-literals that contain
// spreads of other literals.
//
// Does not recurse into nested spreads because this method is already called
// as part of a postorder traversal and nested spreads will already have been
// flattened.
//
// Example: `[0, ...[1, 2, 3], 4, ...[5]]` => `[0, 1, 2, 3, 4, 5]`
//
// Go-only: array and object literals (and their spreads) are handled
// alike, so this uses the rhino API.
//
// Java: PeepholeFoldConstants.tryFlattenArrayOrObjectLit(Node)
func (p *PeepholeFoldConstants) TryFlattenArrayOrObjectLit(parentLit rhino.Node) rhino.Node {
	for child := parentLit.GetFirstChild(); !child.IsNil(); {
		// We have to store the next element here because nodes may be inserted below.
		spread := child
		child = child.GetNext()

		if !spread.IsSpread() {
			continue
		}

		innerLit := spread.GetOnlyChild()
		if parentLit.GetToken() != innerLit.GetToken() {
			continue // We only want to inline arrays into arrays and objects into objects.
		}

		if parentLit.IsObjectLit() && !PeepholeFoldConstants_IsSafeToFlattenObjectLitSpread(innerLit) {
			continue
		}

		if parentLit.IsArrayLit() && !PeepholeFoldConstants_IsSafeToFlattenArrayLitSpread(innerLit) {
			continue
		}

		parentLit.AddChildrenAfter(innerLit.RemoveChildren(), spread)
		spread.Detach()
		p.ReportChangeToEnclosingScope(parentLit)
	}
	return parentLit
}

// Java: PeepholeFoldConstants.tryFoldStringArrayAccess(Node, Node, Node)
//
// Go-only: left is typed, so the Java's checkState(left.isString()) is gone.
func (p *PeepholeFoldConstants) TryFoldStringArrayAccess(n rhino.Node, left ast.StringLit, right rhino.Node) rhino.Node {
	// If GETPROP/GETELEM is used as assignment target the array literal is
	// acting as a temporary we can't fold it here:
	//    "[][0] += 1"
	if NodeUtil_IsLValue(n) {
		return n
	}

	index, ok := p.GetSideEffectFreeNumberValue(right)
	if !p.isReasonableDoubleVal(index, ok) {
		// Sometimes people like to use complex expressions to index into
		// arrays, or strings to index into array methods.
		return n
	}

	intIndex := jutil.D2I(index)
	if !PeepholeFoldConstants_IsFoldableElementIndexKey(right, intIndex) {
		return n
	}

	if float64(intIndex) != index {
		p.Report(PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR, right)
		return n
	}

	value := left.Value()
	if intIndex < 0 || int(intIndex) >= jutil.StringLength(value) {
		undefined := NodeUtil_NewUndefinedNode(left.Node())
		n.ReplaceWith(undefined)
		p.ReportChangeToEnclosingScope(undefined)
		return undefined
	}

	// Note: For now skip the strings with unicode
	// characters as I don't understand the differences
	// between Java and JavaScript.
	for i := 0; i <= int(intIndex); i++ {
		if c := value[i]; c < 32 || c > 127 {
			return n
		}
	}
	elem := rhino.IR_String(value[intIndex : intIndex+1])

	// Replace the entire GETELEM with the value
	n.ReplaceWith(elem)
	p.ReportChangeToEnclosingScope(elem)
	return elem
}

// TryFoldObjectPropAccess tries to fold the node `n` that's accessing an
// object literal's property. Bails out (skips folding and returns the same
// `n` node) with certainty when any of the following holds true:
//   - the access `n` is L-value
//   - the property references super
//   - the object has side-effects other than those preserved by folding the
//     property
//   - property accessed does not exist on the object (might exist on
//     prototype)
//
// Examples of folding:
//   - `({a() { return 1; }})?.a` ---> `(function() { return 1; })`
//   - `({a() { return 1; }}).a()` ---> `(function() { return 1; }())`
//
// Java: PeepholeFoldConstants.tryFoldObjectPropAccess(Node, Node, String)
//
// Go-only: left is typed, so the Java's isObjectLit() test is gone.
func (p *PeepholeFoldConstants) TryFoldObjectPropAccess(n rhino.Node, left ast.ObjectLit, name string) rhino.Node {
	jutil.CheckArgument(NodeUtil_IsNormalOrOptChainGet(n))

	if NodeUtil_IsLValue(n) {
		// If GETPROP/GETELEM is used as assignment target the object literal is
		// acting as a temporary we can't fold it here:
		//    "{a:x}.a += 1" is not "x += 1"
		jutil.CheckState(!NodeUtil_IsOptChainNode(n)) // optional chains can not be targets of assign
		return n
	}

	// find the last definition in the object literal
	nameID := uint32(strpool.Lookup(name))
	var key rhino.Node
	var value rhino.Node
	for member := range left.Members() {
		c := member.Node()
		switch member.Token() {
		case rhino.Token_SETTER_DEF:
		case rhino.Token_OBJECT_SPREAD:
			// Reset the search because spread could overwrite any previous result.
			key = rhino.NilNode
			value = rhino.NilNode
		case rhino.Token_COMPUTED_PROP:
			computed, _ := member.AsComputedProp() // the case proves the kind
			// don't handle computed properties unless the input is a simple string
			prop, ok := computed.Key().AsStringLit()
			if !ok {
				return n
			}
			if nameID != 0 && prop.Node().GetStringID() == nameID && !c.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER) {
				v, _ := computed.Value() // the nil node if it has none, as in the Java
				key = c
				value = v.Node()
			}
		case rhino.Token_GETTER_DEF, rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF:
			// (Go: each of these kinds names its key and its only child
			// differently, so they are read alike with the rhino API.)
			if nameID != 0 && c.GetStringID() == nameID {
				key = c
				value = key.GetFirstChild()
			}
		default:
			panic(jutil.NewIllegalStateException())
		}
	}

	// Didn't find a definition of the name in the object literal, it might
	// be coming from the Object prototype.
	if value.IsNil() {
		return n
	}

	/* `super` captures a hidden reference to the declaring objectlit, so we can't fold it away. */
	if NodeUtil_ReferencesSuper(value) {
		return n
	}

	/*
	 * Check to see if there are any side-effects to this object-literal.
	 *
	 * <p>We remove the value we're going to use because its side-effects will be preserved.
	 */
	tempValue := rhino.IR_NullNode()
	value.ReplaceWith(tempValue)
	hasSideEffectBesidesValue := p.MayHaveSideEffects(left.Node())
	tempValue.ReplaceWith(value)
	if hasSideEffectBesidesValue {
		return n
	}

	keyIsGetter := NodeUtil_IsGetOrSetKey(key)
	nIsInvoked := NodeUtil_IsInvocationTarget(n)

	if keyIsGetter || nIsInvoked {
		// When the code looks like:
		//   {x: f}.x();
		//   {get x() {...}}.x;
		//   {get ['x']() {...}}.x;

		/*
		 * It's not safe, in general, to convert that to just a function call, because the receiver
		 * value will be wrong.
		 *
		 * <p>However, there are some cases where it's ok which we check here.
		 */
		if !value.IsFunction() || NodeUtil_ReferencesOwnReceiver(value) {
			return n
		}
	}

	value.Detach()

	parent := n.GetParent()

	if NodeUtil_IsOptChainNode(parent) &&
		n.IsFirstChildOf(parent) &&
		!parent.IsOptionalChainStart() {

		/*
		 * If the chain continues after `n`, simply doing `n.replaceWith(value)` below would leave the
		 * subsequent nodes in the current chain segment optional, with their start `n` replaced.
		 *
		 * <p>So, we must ensure that all nodes in the chain's current segment are made non-optional.
		 *
		 * <p>This can happen for e.g.
		 *
		 * <ul>
		 *   <li>`({a() { return 1; }})?.a()` ---> `(function() { return 1; })()`. Here, parent
		 *       OPTCHAIN_CALL node must be converted to CALL.
		 *   <li>`({a() { return 1; }})?.a().b.c?.d` ---> `(function() { return 1; })().b.c?.d`. Here,
		 *       all nodes upto `({a() { return 1; }})?.a().b.c` must become non-optional
		 * </ul>
		 */
		endOfCurrentChain := NodeUtil_GetEndOfOptChainSegment(parent)
		NodeUtil_ConvertToNonOptionalChainSegment(endOfCurrentChain)
	}
	if keyIsGetter {
		value = rhino.IR_Call(value)
		value.PutBooleanProp(rhino.Node_FREE_CALL, true)
	} else if nIsInvoked {
		n.GetParent().PutBooleanProp(rhino.Node_FREE_CALL, true)
	}

	n.ReplaceWith(value)
	p.ReportChangeToEnclosingScope(value)
	p.MarkFunctionsDeleted(n)
	return n
}

// Java: PeepholeFoldConstants.isSafeToFlattenObjectLitSpread(Node)
func PeepholeFoldConstants_IsSafeToFlattenObjectLitSpread(innerLit rhino.Node) bool {
	for prop := innerLit.GetFirstChild(); !prop.IsNil(); prop = prop.GetNext() {
		// Object spread converts all properties to regular data properties. This step evaluates
		// any getters and omits any setters. See https://tc39.es/ecma262/#sec-copydataproperties.
		// For simplicity, we just back off and do not flatten spread in this case.
		if NodeUtil_IsGetOrSetKey(prop) {
			return false
		}
	}
	return true
}

// Java: PeepholeFoldConstants.isSafeToFlattenArrayLitSpread(Node)
func PeepholeFoldConstants_IsSafeToFlattenArrayLitSpread(innerLit rhino.Node) bool {
	for n := innerLit.GetFirstChild(); !n.IsNil(); n = n.GetNext() {
		if n.IsEmpty() {
			return false
		}
	}
	return true
}
