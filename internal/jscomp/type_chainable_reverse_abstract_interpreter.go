// Port of com/google/javascript/jscomp/type/ChainableReverseAbstractInterpreter.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReverseAbstractInterpreter defines what reversed abstract interpreters
// provide.
//
// Abstract interpretation is the process of interpreting a program at an
// abstracted level (such as at the type level) instead of the concrete level
// (the flow of values). This reversed abstract interpreter reverses the
// abstract interpretation process by knowing the outcome of some computation
// and calculating a preciser view of the world than the view without knowing
// the outcome of the computation.
type ReverseAbstractInterpreter interface {
	// GetPreciserScopeKnowingConditionOutcome calculates a precise version of
	// the scope knowing the outcome of the condition.
	//
	// condition is the condition's expression; blindScope is the scope
	// without knowledge about the outcome of the condition; outcome is the
	// outcome of the condition.
	//
	// Java: ReverseAbstractInterpreter.getPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
	GetPreciserScopeKnowingConditionOutcome(
		condition rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope
}

// ChainableReverseAbstractInterpreter is a chainable reverse abstract
// interpreter providing basic functionality.
//
// Java: abstract class ChainableReverseAbstractInterpreter implements
// ReverseAbstractInterpreter.
type ChainableReverseAbstractInterpreter interface {
	ReverseAbstractInterpreter

	// Java: ChainableReverseAbstractInterpreter.append(ChainableReverseAbstractInterpreter)
	Append(lastLink ChainableReverseAbstractInterpreter) ChainableReverseAbstractInterpreter
	// Java: ChainableReverseAbstractInterpreter.getFirst()
	GetFirst() ChainableReverseAbstractInterpreter

	getChainableBase() *ChainableReverseAbstractInterpreterBase
}

// ChainableReverseAbstractInterpreterBase holds
// ChainableReverseAbstractInterpreter's fields and base methods.
type ChainableReverseAbstractInterpreterBase struct {
	typeRegistry *rhino.JSTypeRegistry
	firstLink    ChainableReverseAbstractInterpreter
	nextLink     ChainableReverseAbstractInterpreter
}

// initChainableBase initializes an interpreter as the sole link in a chain.
// Interpreters can be appended using Append.
//
// Java: ChainableReverseAbstractInterpreter(JSTypeRegistry)
func (b *ChainableReverseAbstractInterpreterBase) initChainableBase(
	firstLink ChainableReverseAbstractInterpreter, typeRegistry *rhino.JSTypeRegistry) {
	b.typeRegistry = typeRegistry
	b.firstLink = firstLink
	b.nextLink = nil
}

func (b *ChainableReverseAbstractInterpreterBase) getChainableBase() *ChainableReverseAbstractInterpreterBase {
	return b
}

// Append appends a link to this, returning the updated last link.
//
// The pattern new X().append(new Y())...append(new Z()) forms a chain
// starting with X, then Y, then ... Z.
//
// lastLink is a chainable interpreter, with no next link; returns the
// updated last link.
//
// Java: ChainableReverseAbstractInterpreter.append(ChainableReverseAbstractInterpreter)
func (b *ChainableReverseAbstractInterpreterBase) Append(
	lastLink ChainableReverseAbstractInterpreter) ChainableReverseAbstractInterpreter {
	lastBase := lastLink.getChainableBase()
	jutil.CheckArgument(lastBase.nextLink == nil)
	b.nextLink = lastLink
	lastBase.firstLink = b.firstLink
	return lastLink
}

// GetFirst gets the first link of this chain.
//
// Java: ChainableReverseAbstractInterpreter.getFirst()
func (b *ChainableReverseAbstractInterpreterBase) GetFirst() ChainableReverseAbstractInterpreter {
	return b.firstLink
}

// FirstPreciserScopeKnowingConditionOutcome calculates the preciser scope
// starting with the first link.
//
// Java: ChainableReverseAbstractInterpreter.firstPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
func (b *ChainableReverseAbstractInterpreterBase) FirstPreciserScopeKnowingConditionOutcome(
	condition rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	return b.firstLink.GetPreciserScopeKnowingConditionOutcome(
		condition, blindScope, outcome)
}

// NextPreciserScopeKnowingConditionOutcome delegates the calculation of the
// preciser scope to the next link. If there is no next link, returns the
// blind scope.
//
// Java: ChainableReverseAbstractInterpreter.nextPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
func (b *ChainableReverseAbstractInterpreterBase) NextPreciserScopeKnowingConditionOutcome(
	condition rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	if b.nextLink != nil {
		return b.nextLink.GetPreciserScopeKnowingConditionOutcome(
			condition, blindScope, outcome)
	}
	return blindScope
}

// GetTypeIfRefinable returns the type of a node in the given scope if the
// node corresponds to a name whose type is capable of being refined.
//
// Returns the current type of the node if it can be refined, null
// otherwise.
//
// Java: ChainableReverseAbstractInterpreter.getTypeIfRefinable(Node, FlowScope)
func (b *ChainableReverseAbstractInterpreterBase) GetTypeIfRefinable(node rhino.Node, scope FlowScope) rhino.JSType {
	switch node.GetToken() {
	case rhino.Token_NAME:
		nameVar, _ := scope.GetSlot(node.GetString()).(rhino.StaticTypedSlot)
		if nameVar != nil {
			nameVarType := nameVar.GetType()
			if nameVarType == nil {
				nameVarType = node.GetJSType()
			}
			return nameVarType
		}
		return nil
	case rhino.Token_GETPROP:
		qualifiedName, ok := node.GetQualifiedNameString()
		if !ok {
			return nil
		}
		propVar, _ := scope.GetSlot(qualifiedName).(rhino.StaticTypedSlot)
		var propVarType rhino.JSType
		if propVar != nil {
			propVarType = propVar.GetType()
		}
		if propVarType == nil {
			propVarType = node.GetJSType()
		}
		if propVarType == nil {
			propVarType = b.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		return propVarType
	default:
	}
	return nil
}

// DeclareNameInScope declares a refined type in scope for the name
// represented by node. It must be possible to refine the type of the given
// node in the given scope, as determined by GetTypeIfRefinable. Returns an
// updated flow scope, which may be different from the passed-in scope if any
// changes occur.
//
// Java: ChainableReverseAbstractInterpreter.declareNameInScope(FlowScope, Node, JSType)
func (b *ChainableReverseAbstractInterpreterBase) DeclareNameInScope(
	scope FlowScope, node rhino.Node, type_ rhino.JSType) FlowScope {
	switch node.GetToken() {
	case rhino.Token_NAME:
		return scope.InferSlotType(node.GetString(), type_)
	case rhino.Token_GETPROP:
		qualifiedName, ok := node.GetQualifiedNameString()
		if !ok {
			panic(jutil.NewNullPointerException())
		}

		origType := node.GetJSType()
		if origType == nil {
			origType = b.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		return scope.InferQualifiedSlot(node, qualifiedName, origType, type_, false)
	case rhino.Token_THIS:
		// "this" references aren't currently modeled in the CFG.
		return scope
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Node cannot be refined. \n" + node.ToStringTree()))
	}
}

// GetRestrictedByTypeOfResult returns a version of type that is restricted
// by some knowledge about the result of the typeof operation.
//
// The behavior of the typeof operator can be summarized by the following
// table:
//
//	type                                      result
//	undefined                                 "undefined"
//	null                                      "object"
//	boolean                                   "boolean"
//	number                                    "number"
//	string                                    "string"
//	Object (which doesn't implement [[Call]]) "object"
//	Object (which implements [[Call]])        "function"
//
// type is the type to restrict; value is a value known to be equal or not
// equal to the result of the typeof operation; resultEqualsValue is true if
// the typeOf result is known to equal value, false if it is known not to
// equal value. Returns the restricted type or null if no version of the type
// matches the restriction.
//
// Java: ChainableReverseAbstractInterpreter.getRestrictedByTypeOfResult(JSType, String, boolean)
func (b *ChainableReverseAbstractInterpreterBase) GetRestrictedByTypeOfResult(
	type_ rhino.JSType, value string, resultEqualsValue bool) rhino.JSType {
	if type_ == nil {
		if resultEqualsValue {
			result := b.GetNativeTypeForTypeOf(value)
			if result == nil {
				return b.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
			}
			return result
		}
		return nil
	}
	return b.restrictByOneTypeOf(type_, value, resultEqualsValue)
}

// Java: ChainableReverseAbstractInterpreter.RestrictByOneTypeOfResultVisitor
func (b *ChainableReverseAbstractInterpreterBase) restrictByOneTypeOf(
	t rhino.JSType, value string, resultEqualsValue bool) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := b.restrictByOneTypeOf(prim, value, resultEqualsValue)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.AllType:
		return b.restrictTopByOneTypeOf(b.GetNativeType(rhino.JSTypeNative_ALL_TYPE), value, resultEqualsValue)
	case *rhino.UnknownType:
		return b.restrictTopByOneTypeOf(b.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE), value, resultEqualsValue)
	case *rhino.NoObjectTypeBase:
		if (value == "object" || value == "function") == resultEqualsValue {
			return b.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
		}
		return nil
	case *rhino.BooleanType:
		if (value == "boolean") == resultEqualsValue {
			return b.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
		}
		return nil
	case *rhino.FunctionTypeBase:
		if (value == "function") == resultEqualsValue {
			return x
		}
		return nil
	case *rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase, *rhino.RecordType,
		*rhino.EnumType, *rhino.NoResolvedType, *rhino.TemplatizedType, *rhino.TemplateType:
		return b.restrictObjectByOneTypeOf(t.(rhino.ObjectType), value, resultEqualsValue)
	case *rhino.NullType:
		if (value == "object") == resultEqualsValue {
			return b.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
		}
		return nil
	case *rhino.NamedType:
		return b.restrictByOneTypeOf(x.GetReferencedTypeInternal(), value, resultEqualsValue)
	case *rhino.ProxyObjectTypeBase:
		return b.restrictByOneTypeOf(x.GetReferencedTypeInternal(), value, resultEqualsValue)
	case *rhino.NumberType:
		if (value == "number") == resultEqualsValue {
			return b.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
		}
		return nil
	case *rhino.BigIntType:
		if (value == "bigint") == resultEqualsValue {
			return b.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
		}
		return nil
	case *rhino.StringType:
		if (value == "string") == resultEqualsValue {
			return b.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
		}
		return nil
	case *rhino.SymbolTypeBase, *rhino.KnownSymbolType:
		if (value == "symbol") == resultEqualsValue {
			return t
		}
		return nil
	case *rhino.VoidType:
		if (value == "undefined") == resultEqualsValue {
			return b.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		}
		return nil
	case *rhino.UnionType:
		var restricted rhino.JSType
		for _, alt := range x.GetAlternates() {
			if ra := b.restrictByOneTypeOf(alt, value, resultEqualsValue); ra != nil {
				if restricted == nil {
					restricted = ra
				} else {
					restricted = ra.GetLeastSupertype(restricted)
				}
			}
		}
		return restricted
	default:
		if _, ok := t.(rhino.NoObjectType); ok {
			if (value == "object" || value == "function") == resultEqualsValue {
				return b.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
			}
			return nil
		}
		if fn, ok := t.(rhino.FunctionType); ok {
			if (value == "function") == resultEqualsValue {
				return fn
			}
			return nil
		}
		if proxy, ok := t.(rhino.ProxyObjectType); ok {
			return b.restrictByOneTypeOf(proxy.GetReferencedTypeInternal(), value, resultEqualsValue)
		}
		if obj, ok := t.(rhino.ObjectType); ok {
			return b.restrictObjectByOneTypeOf(obj, value, resultEqualsValue)
		}
		if sym, ok := t.(rhino.SymbolType); ok {
			if (value == "symbol") == resultEqualsValue {
				return sym
			}
			return nil
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

func (b *ChainableReverseAbstractInterpreterBase) restrictTopByOneTypeOf(
	topType rhino.JSType, value string, resultEqualsValue bool) rhino.JSType {
	if resultEqualsValue {
		if typeByName := b.GetNativeTypeForTypeOf(value); typeByName != nil {
			return typeByName
		}
	}
	return topType
}

func (b *ChainableReverseAbstractInterpreterBase) restrictObjectByOneTypeOf(
	type_ rhino.ObjectType, value string, resultEqualsValue bool) rhino.JSType {
	if value == "function" {
		ctorType := b.GetNativeType(rhino.JSTypeNative_FUNCTION_TYPE)
		if resultEqualsValue {
			// Objects are restricted to "Function", subtypes are left
			return ctorType.GetGreatestSubtype(type_)
		}
		// Only filter out subtypes of "function"
		if type_.IsSubtypeOf(ctorType) {
			return nil
		}
		return type_
	}
	if (value == "object") == resultEqualsValue {
		return type_
	}
	return nil
}

// Java: ChainableReverseAbstractInterpreter.getNativeType(JSTypeNative)
func (b *ChainableReverseAbstractInterpreterBase) GetNativeType(typeId rhino.JSTypeNative) rhino.JSType {
	return b.typeRegistry.GetNativeType(typeId)
}

// GetNativeTypeForTypeOf: if we definitely know what a type is based on the
// typeof result, return it. Otherwise, return null.
//
// The typeof operation in JS is poorly defined, and this function works for
// both the native typeof and goog.typeOf. It should not be made public,
// because its semantics are informally defined, and would be wrong in the
// general case.
//
// Java: ChainableReverseAbstractInterpreter.getNativeTypeForTypeOf(String)
func (b *ChainableReverseAbstractInterpreterBase) GetNativeTypeForTypeOf(value string) rhino.JSType {
	switch value {
	case "number":
		return b.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	case "boolean":
		return b.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	case "string":
		return b.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	case "symbol":
		return b.GetNativeType(rhino.JSTypeNative_SYMBOL_TYPE)
	case "undefined":
		return b.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	case "object":
		// NOTE: This is broader than it needs to be if it's from goog.typeof, but (a) it's more
		// consistent with common usage of the native builtin typeof, (b) it's more consistent
		// with TypeScript, and (c) it's more useful than simply not narrowing.
		return b.typeRegistry.CreateUnionType(
			b.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE), b.GetNativeType(rhino.JSTypeNative_NULL_TYPE))
	case "function":
		return b.GetNativeType(rhino.JSTypeNative_FUNCTION_TYPE)
	default:
		return nil
	}
}

// restrictToArray restricts t when goog.isArray or Array.isArray returns true.
//
// Java: ChainableReverseAbstractInterpreter.restrictToArrayVisitor
func (b *ChainableReverseAbstractInterpreterBase) restrictToArray(t rhino.JSType) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.AllType:
		return b.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)
	case *rhino.UnknownType:
		return b.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := b.restrictToArray(prim)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.UnionType:
		var restricted rhino.JSType
		for _, alt := range x.GetAlternates() {
			if ra := b.restrictToArray(alt); ra != nil {
				if restricted == nil {
					restricted = ra
				} else {
					restricted = ra.GetLeastSupertype(restricted)
				}
			}
		}
		return restricted
	case *rhino.NamedType:
		return b.restrictToArray(x.GetReferencedTypeInternal())
	case *rhino.ProxyObjectTypeBase:
		return b.restrictToArray(x.GetReferencedTypeInternal())
	case *rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase, *rhino.RecordType,
		*rhino.EnumType, *rhino.NoResolvedType, *rhino.TemplatizedType, *rhino.TemplateType:
		arrayType := b.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)
		if arrayType.IsSubtypeOf(t) {
			return arrayType
		}
		return nil
	case *rhino.NoObjectTypeBase, *rhino.BooleanType, *rhino.FunctionTypeBase,
		*rhino.NullType, *rhino.NumberType, *rhino.BigIntType, *rhino.StringType,
		*rhino.SymbolTypeBase, *rhino.KnownSymbolType, *rhino.VoidType:
		return nil
	default:
		if _, ok := t.(rhino.NoObjectType); ok {
			return nil
		}
		if _, ok := t.(rhino.FunctionType); ok {
			return nil
		}
		if proxy, ok := t.(rhino.ProxyObjectType); ok {
			return b.restrictToArray(proxy.GetReferencedTypeInternal())
		}
		if obj, ok := t.(rhino.ObjectType); ok {
			arrayType := b.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)
			if arrayType.IsSubtypeOf(obj) {
				return arrayType
			}
			return nil
		}
		if _, ok := t.(rhino.SymbolType); ok {
			return nil
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

// restrictToNotArray restricts t when goog.isArray or Array.isArray returns false.
//
// Java: ChainableReverseAbstractInterpreter.restrictToNotArrayVisitor
func (b *ChainableReverseAbstractInterpreterBase) restrictToNotArray(t rhino.JSType) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.AllType:
		return b.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
	case *rhino.UnknownType:
		return b.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
	case *rhino.NoObjectTypeBase:
		return b.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
	case *rhino.BooleanType:
		return b.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	case *rhino.NullType:
		return b.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
	case *rhino.NumberType:
		return b.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	case *rhino.BigIntType:
		return b.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
	case *rhino.StringType:
		return b.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	case *rhino.VoidType:
		return b.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	case *rhino.FunctionTypeBase, *rhino.SymbolTypeBase, *rhino.KnownSymbolType:
		return t
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := b.restrictToNotArray(prim)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.UnionType:
		var restricted rhino.JSType
		for _, alt := range x.GetAlternates() {
			if ra := b.restrictToNotArray(alt); ra != nil {
				if restricted == nil {
					restricted = ra
				} else {
					restricted = ra.GetLeastSupertype(restricted)
				}
			}
		}
		return restricted
	case *rhino.NamedType:
		return b.restrictToNotArray(x.GetReferencedTypeInternal())
	case *rhino.ProxyObjectTypeBase:
		return b.restrictToNotArray(x.GetReferencedTypeInternal())
	case *rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase, *rhino.RecordType,
		*rhino.EnumType, *rhino.NoResolvedType, *rhino.TemplatizedType, *rhino.TemplateType:
		if t.IsSubtypeOf(b.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)) {
			return nil
		}
		return t
	default:
		if _, ok := t.(rhino.NoObjectType); ok {
			return b.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
		}
		if _, ok := t.(rhino.FunctionType); ok {
			return t
		}
		if proxy, ok := t.(rhino.ProxyObjectType); ok {
			return b.restrictToNotArray(proxy.GetReferencedTypeInternal())
		}
		if obj, ok := t.(rhino.ObjectType); ok {
			if obj.IsSubtypeOf(b.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)) {
				return nil
			}
			return obj
		}
		if _, ok := t.(rhino.SymbolType); ok {
			return t
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}
