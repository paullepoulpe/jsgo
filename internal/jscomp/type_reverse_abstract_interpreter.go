// Port of com/google/javascript/jscomp/type/ChainableReverseAbstractInterpreter.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/type/ClosureReverseAbstractInterpreter.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	raiGoogID     = uint32(strpool.Intern("goog"))
	raiIsObjectID = uint32(strpool.Intern("isObject"))
	raiArrayID    = uint32(strpool.Intern("Array"))
	raiIsArrayID  = uint32(strpool.Intern("isArray"))
)

type raiChainMode uint8

const (
	raiModeSemantic raiChainMode = iota
	raiModeClosure
	raiModeClosureThenSemantic
	raiModeSemanticThenClosure
)

// ReverseAbstractInterpreter is a concrete pointer to the unified reverse
// abstract interpreter, avoiding interface dispatch on every conditional edge
// during dataflow type inference.
type ReverseAbstractInterpreter = *SemanticReverseAbstractInterpreter

// ChainableReverseAbstractInterpreter aliases ReverseAbstractInterpreter for
// caller compatibility without linked-list interface indirection.
type ChainableReverseAbstractInterpreter = *SemanticReverseAbstractInterpreter

// ClosureReverseAbstractInterpreter aliases SemanticReverseAbstractInterpreter
// configured in Closure narrowing mode.
type ClosureReverseAbstractInterpreter = SemanticReverseAbstractInterpreter

// SemanticReverseAbstractInterpreter is the unified concrete reverse abstract
// interpreter supporting semantic JS narrowing, Closure-specific narrowing
// (goog.isObject), and chained combinations of both without dynamic dispatch.
type SemanticReverseAbstractInterpreter struct {
	typeRegistry   *rhino.JSTypeRegistry
	mode           raiChainMode
	activeTrackers []*sraiRefinementTracker
}

// NewSemanticReverseAbstractInterpreter creates a semantic reverse abstract
// interpreter.
//
// Java: SemanticReverseAbstractInterpreter(JSTypeRegistry)
func NewSemanticReverseAbstractInterpreter(typeRegistry *rhino.JSTypeRegistry) *SemanticReverseAbstractInterpreter {
	return &SemanticReverseAbstractInterpreter{
		typeRegistry: typeRegistry,
		mode:         raiModeSemantic,
	}
}

// NewClosureReverseAbstractInterpreter creates a reverse abstract interpreter
// for Closure patterns such as goog.isObject.
//
// Java: ClosureReverseAbstractInterpreter(JSTypeRegistry)
func NewClosureReverseAbstractInterpreter(typeRegistry *rhino.JSTypeRegistry) *ClosureReverseAbstractInterpreter {
	return &SemanticReverseAbstractInterpreter{
		typeRegistry: typeRegistry,
		mode:         raiModeClosure,
	}
}

func (s *SemanticReverseAbstractInterpreter) hasNextLink() bool {
	return s.mode == raiModeClosureThenSemantic || s.mode == raiModeSemanticThenClosure
}

// Append combines s with lastLink into a single concrete interpreter mode.
//
// Java: ChainableReverseAbstractInterpreter.append(ChainableReverseAbstractInterpreter)
func (s *SemanticReverseAbstractInterpreter) Append(
	lastLink ChainableReverseAbstractInterpreter) ChainableReverseAbstractInterpreter {
	jutil.CheckArgument(lastLink != nil && !s.hasNextLink() && !lastLink.hasNextLink())
	var combined raiChainMode
	switch {
	case s.mode == raiModeClosure && lastLink.mode == raiModeSemantic:
		combined = raiModeClosureThenSemantic
	case s.mode == raiModeSemantic && lastLink.mode == raiModeClosure:
		combined = raiModeSemanticThenClosure
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Unsupported ReverseAbstractInterpreter chain"))
	}
	s.mode = combined
	lastLink.mode = combined
	return s
}

// GetFirst returns this interpreter.
//
// Java: ChainableReverseAbstractInterpreter.getFirst()
func (s *SemanticReverseAbstractInterpreter) GetFirst() ChainableReverseAbstractInterpreter {
	return s
}

// GetPreciserScopeKnowingConditionOutcome calculates a precise version of the
// scope knowing the outcome of the condition.
//
// Java: ReverseAbstractInterpreter.getPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) GetPreciserScopeKnowingConditionOutcome(
	condition rhino.Node, blindScope *LinkedFlowScope, outcome rhino.Outcome) *LinkedFlowScope {
	switch s.mode {
	case raiModeSemantic:
		scope, _ := s.semanticPreciserScope(condition, blindScope, outcome)
		return scope
	case raiModeClosureThenSemantic:
		if scope, handled := s.closurePreciserScope(condition, blindScope, outcome); handled {
			return scope
		}
		scope, _ := s.semanticPreciserScope(condition, blindScope, outcome)
		return scope
	case raiModeSemanticThenClosure:
		if scope, handled := s.semanticPreciserScope(condition, blindScope, outcome); handled {
			return scope
		}
		scope, _ := s.closurePreciserScope(condition, blindScope, outcome)
		return scope
	case raiModeClosure:
		scope, _ := s.closurePreciserScope(condition, blindScope, outcome)
		return scope
	default:
		return blindScope
	}
}

// closurePreciserScope handles Closure-specific narrowing (goog.isObject)
// using pre-interned strpool IDs without allocating qualified name strings.
//
// Java: ClosureReverseAbstractInterpreter.getPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) closurePreciserScope(
	condition rhino.Node, blindScope *LinkedFlowScope, outcome rhino.Outcome) (*LinkedFlowScope, bool) {
	if call, ok := ast.AsCall(condition); ok {
		if param, ok := call.OnlyArg(); ok {
			if callee, ok := call.Callee().AsGetProp(); ok && callee.Node().GetStringID() == raiIsObjectID {
				obj := callee.Object().Node()
				if obj.IsName() && obj.GetStringID() == raiGoogID && param.Node().IsQualifiedName() {
					paramType := s.GetTypeIfRefinable(param.Node(), blindScope)
					return s.restrictIsObject(param.Node(), paramType, blindScope, outcome.IsTruthy()), true
				}
			}
		}
	}
	return blindScope, false
}

func (s *SemanticReverseAbstractInterpreter) restrictIsObject(
	parameter rhino.Node,
	type_ rhino.JSType,
	blindScope *LinkedFlowScope,
	outcome bool) *LinkedFlowScope {
	if type_ == nil {
		if outcome {
			type_ = s.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)
		}
	} else if outcome {
		type_ = s.restrictToObject(type_)
	} else {
		type_ = s.restrictToNotObject(type_)
	}
	if type_ != nil {
		return s.DeclareNameInScope(blindScope, parameter, type_)
	}
	return blindScope
}

// restrictToObject restricts t when goog.isObject returns true (includes
// functions, excludes null).
//
// Java: ClosureReverseAbstractInterpreter.restrictToObjectVisitor
func (s *SemanticReverseAbstractInterpreter) restrictToObject(t rhino.JSType) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.AllType, *rhino.UnknownType:
		return s.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := s.restrictToObject(prim)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.UnionType:
		var restricted rhino.JSType
		for _, alt := range x.GetAlternates() {
			if ra := s.restrictToObject(alt); ra != nil {
				if restricted == nil {
					restricted = ra
				} else {
					restricted = ra.GetLeastSupertype(restricted)
				}
			}
		}
		return restricted
	case *rhino.NamedType:
		return s.restrictToObject(x.GetReferencedTypeInternal())
	case *rhino.ProxyObjectTypeBase:
		return s.restrictToObject(x.GetReferencedTypeInternal())
	case *rhino.FunctionTypeBase, *rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase,
		*rhino.RecordType, *rhino.EnumType, *rhino.NoResolvedType, *rhino.TemplatizedType, *rhino.TemplateType:
		return t
	case *rhino.NoObjectTypeBase, *rhino.BooleanType, *rhino.NullType,
		*rhino.NumberType, *rhino.BigIntType, *rhino.StringType,
		*rhino.SymbolTypeBase, *rhino.KnownSymbolType, *rhino.VoidType:
		return nil
	default:
		if _, ok := t.(rhino.NoObjectType); ok {
			return nil
		}
		if _, ok := t.(rhino.FunctionType); ok {
			return t
		}
		if proxy, ok := t.(rhino.ProxyObjectType); ok {
			return s.restrictToObject(proxy.GetReferencedTypeInternal())
		}
		if _, ok := t.(rhino.ObjectType); ok {
			return t
		}
		if _, ok := t.(rhino.SymbolType); ok {
			return nil
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

// restrictToNotObject restricts t when goog.isObject returns false.
//
// Java: ClosureReverseAbstractInterpreter.restrictToNotObjectVisitor
func (s *SemanticReverseAbstractInterpreter) restrictToNotObject(t rhino.JSType) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.AllType:
		return s.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
	case *rhino.UnknownType:
		return s.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
	case *rhino.NoObjectTypeBase:
		return s.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
	case *rhino.BooleanType:
		return s.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	case *rhino.NullType:
		return s.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
	case *rhino.NumberType:
		return s.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	case *rhino.BigIntType:
		return s.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
	case *rhino.StringType:
		return s.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	case *rhino.VoidType:
		return s.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	case *rhino.SymbolTypeBase, *rhino.KnownSymbolType:
		return t
	case *rhino.FunctionTypeBase, *rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase,
		*rhino.RecordType, *rhino.EnumType, *rhino.NoResolvedType, *rhino.TemplatizedType, *rhino.TemplateType:
		return nil
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := s.restrictToNotObject(prim)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.UnionType:
		var restricted rhino.JSType
		for _, alt := range x.GetAlternates() {
			if ra := s.restrictToNotObject(alt); ra != nil {
				if restricted == nil {
					restricted = ra
				} else {
					restricted = ra.GetLeastSupertype(restricted)
				}
			}
		}
		return restricted
	case *rhino.NamedType:
		return s.restrictToNotObject(x.GetReferencedTypeInternal())
	case *rhino.ProxyObjectTypeBase:
		return s.restrictToNotObject(x.GetReferencedTypeInternal())
	default:
		if _, ok := t.(rhino.NoObjectType); ok {
			return s.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
		}
		if _, ok := t.(rhino.FunctionType); ok {
			return nil
		}
		if proxy, ok := t.(rhino.ProxyObjectType); ok {
			return s.restrictToNotObject(proxy.GetReferencedTypeInternal())
		}
		if _, ok := t.(rhino.ObjectType); ok {
			return nil
		}
		if _, ok := t.(rhino.SymbolType); ok {
			return t
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

// GetTypeIfRefinable returns the type of a node in the given scope if the
// node corresponds to a name whose type is capable of being refined.
//
// Returns the current type of the node if it can be refined, null otherwise.
//
// Java: ChainableReverseAbstractInterpreter.getTypeIfRefinable(Node, FlowScope)
func (s *SemanticReverseAbstractInterpreter) GetTypeIfRefinable(node rhino.Node, scope *LinkedFlowScope) rhino.JSType {
	switch node.GetToken() {
	case rhino.Token_NAME:
		if _, _, nameVarType, _, ok := scope.LookupSlotType(node.GetString()); ok {
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
		_, _, propVarType, _, _ := scope.LookupSlotType(qualifiedName)
		if propVarType == nil {
			propVarType = node.GetJSType()
		}
		if propVarType == nil {
			propVarType = s.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
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
func (s *SemanticReverseAbstractInterpreter) DeclareNameInScope(
	scope *LinkedFlowScope, node rhino.Node, type_ rhino.JSType) *LinkedFlowScope {
	switch node.GetToken() {
	case rhino.Token_NAME:
		name := node.GetString()
		s.recordRefinement(name)
		return scope.InferSlotType(name, type_)
	case rhino.Token_GETPROP:
		qualifiedName, ok := node.GetQualifiedNameString()
		if !ok {
			panic(jutil.NewNullPointerException())
		}

		origType := node.GetJSType()
		if origType == nil {
			origType = s.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		s.recordRefinement(qualifiedName)
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
// Java: ChainableReverseAbstractInterpreter.getRestrictedByTypeOfResult(JSType, String, boolean)
func (s *SemanticReverseAbstractInterpreter) GetRestrictedByTypeOfResult(
	type_ rhino.JSType, value string, resultEqualsValue bool) rhino.JSType {
	if type_ == nil {
		if resultEqualsValue {
			result := s.GetNativeTypeForTypeOf(value)
			if result == nil {
				return s.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
			}
			return result
		}
		return nil
	}
	return s.restrictByOneTypeOf(type_, value, resultEqualsValue)
}

// Java: ChainableReverseAbstractInterpreter.RestrictByOneTypeOfResultVisitor
func (s *SemanticReverseAbstractInterpreter) restrictByOneTypeOf(
	t rhino.JSType, value string, resultEqualsValue bool) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := s.restrictByOneTypeOf(prim, value, resultEqualsValue)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.AllType:
		return s.restrictTopByOneTypeOf(s.GetNativeType(rhino.JSTypeNative_ALL_TYPE), value, resultEqualsValue)
	case *rhino.UnknownType:
		return s.restrictTopByOneTypeOf(s.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE), value, resultEqualsValue)
	case *rhino.NoObjectTypeBase:
		if (value == "object" || value == "function") == resultEqualsValue {
			return s.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
		}
		return nil
	case *rhino.BooleanType:
		if (value == "boolean") == resultEqualsValue {
			return s.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
		}
		return nil
	case *rhino.FunctionTypeBase:
		if (value == "function") == resultEqualsValue {
			return x
		}
		return nil
	case *rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase, *rhino.RecordType,
		*rhino.EnumType, *rhino.NoResolvedType, *rhino.TemplatizedType, *rhino.TemplateType:
		return s.restrictObjectByOneTypeOf(t.(rhino.ObjectType), value, resultEqualsValue)
	case *rhino.NullType:
		if (value == "object") == resultEqualsValue {
			return s.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
		}
		return nil
	case *rhino.NamedType:
		return s.restrictByOneTypeOf(x.GetReferencedTypeInternal(), value, resultEqualsValue)
	case *rhino.ProxyObjectTypeBase:
		return s.restrictByOneTypeOf(x.GetReferencedTypeInternal(), value, resultEqualsValue)
	case *rhino.NumberType:
		if (value == "number") == resultEqualsValue {
			return s.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
		}
		return nil
	case *rhino.BigIntType:
		if (value == "bigint") == resultEqualsValue {
			return s.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
		}
		return nil
	case *rhino.StringType:
		if (value == "string") == resultEqualsValue {
			return s.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
		}
		return nil
	case *rhino.SymbolTypeBase, *rhino.KnownSymbolType:
		if (value == "symbol") == resultEqualsValue {
			return t
		}
		return nil
	case *rhino.VoidType:
		if (value == "undefined") == resultEqualsValue {
			return s.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		}
		return nil
	case *rhino.UnionType:
		var restricted rhino.JSType
		for _, alt := range x.GetAlternates() {
			if ra := s.restrictByOneTypeOf(alt, value, resultEqualsValue); ra != nil {
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
				return s.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
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
			return s.restrictByOneTypeOf(proxy.GetReferencedTypeInternal(), value, resultEqualsValue)
		}
		if obj, ok := t.(rhino.ObjectType); ok {
			return s.restrictObjectByOneTypeOf(obj, value, resultEqualsValue)
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

func (s *SemanticReverseAbstractInterpreter) restrictTopByOneTypeOf(
	topType rhino.JSType, value string, resultEqualsValue bool) rhino.JSType {
	if resultEqualsValue {
		if typeByName := s.GetNativeTypeForTypeOf(value); typeByName != nil {
			return typeByName
		}
	}
	return topType
}

func (s *SemanticReverseAbstractInterpreter) restrictObjectByOneTypeOf(
	type_ rhino.ObjectType, value string, resultEqualsValue bool) rhino.JSType {
	if value == "function" {
		ctorType := s.GetNativeType(rhino.JSTypeNative_FUNCTION_TYPE)
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
func (s *SemanticReverseAbstractInterpreter) GetNativeType(typeId rhino.JSTypeNative) rhino.JSType {
	return s.typeRegistry.GetNativeType(typeId)
}

// GetNativeTypeForTypeOf: if we definitely know what a type is based on the
// typeof result, return it. Otherwise, return null.
//
// Java: ChainableReverseAbstractInterpreter.getNativeTypeForTypeOf(String)
func (s *SemanticReverseAbstractInterpreter) GetNativeTypeForTypeOf(value string) rhino.JSType {
	switch value {
	case "number":
		return s.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	case "boolean":
		return s.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	case "string":
		return s.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	case "symbol":
		return s.GetNativeType(rhino.JSTypeNative_SYMBOL_TYPE)
	case "undefined":
		return s.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	case "object":
		// NOTE: This is broader than it needs to be if it's from goog.typeof, but (a) it's more
		// consistent with common usage of the native builtin typeof, (b) it's more consistent
		// with TypeScript, and (c) it's more useful than simply not narrowing.
		return s.typeRegistry.CreateUnionType(
			s.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE), s.GetNativeType(rhino.JSTypeNative_NULL_TYPE))
	case "function":
		return s.GetNativeType(rhino.JSTypeNative_FUNCTION_TYPE)
	default:
		return nil
	}
}

// restrictToArray restricts t when goog.isArray or Array.isArray returns true.
//
// Java: ChainableReverseAbstractInterpreter.restrictToArrayVisitor
func (s *SemanticReverseAbstractInterpreter) restrictToArray(t rhino.JSType) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.AllType:
		return s.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)
	case *rhino.UnknownType:
		return s.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := s.restrictToArray(prim)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.UnionType:
		var restricted rhino.JSType
		for _, alt := range x.GetAlternates() {
			if ra := s.restrictToArray(alt); ra != nil {
				if restricted == nil {
					restricted = ra
				} else {
					restricted = ra.GetLeastSupertype(restricted)
				}
			}
		}
		return restricted
	case *rhino.NamedType:
		return s.restrictToArray(x.GetReferencedTypeInternal())
	case *rhino.ProxyObjectTypeBase:
		return s.restrictToArray(x.GetReferencedTypeInternal())
	case *rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase, *rhino.RecordType,
		*rhino.EnumType, *rhino.NoResolvedType, *rhino.TemplatizedType, *rhino.TemplateType:
		arrayType := s.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)
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
			return s.restrictToArray(proxy.GetReferencedTypeInternal())
		}
		if obj, ok := t.(rhino.ObjectType); ok {
			arrayType := s.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)
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
func (s *SemanticReverseAbstractInterpreter) restrictToNotArray(t rhino.JSType) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.AllType:
		return s.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
	case *rhino.UnknownType:
		return s.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
	case *rhino.NoObjectTypeBase:
		return s.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
	case *rhino.BooleanType:
		return s.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	case *rhino.NullType:
		return s.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
	case *rhino.NumberType:
		return s.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	case *rhino.BigIntType:
		return s.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
	case *rhino.StringType:
		return s.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	case *rhino.VoidType:
		return s.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	case *rhino.FunctionTypeBase, *rhino.SymbolTypeBase, *rhino.KnownSymbolType:
		return t
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := s.restrictToNotArray(prim)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.UnionType:
		var restricted rhino.JSType
		for _, alt := range x.GetAlternates() {
			if ra := s.restrictToNotArray(alt); ra != nil {
				if restricted == nil {
					restricted = ra
				} else {
					restricted = ra.GetLeastSupertype(restricted)
				}
			}
		}
		return restricted
	case *rhino.NamedType:
		return s.restrictToNotArray(x.GetReferencedTypeInternal())
	case *rhino.ProxyObjectTypeBase:
		return s.restrictToNotArray(x.GetReferencedTypeInternal())
	case *rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase, *rhino.RecordType,
		*rhino.EnumType, *rhino.NoResolvedType, *rhino.TemplatizedType, *rhino.TemplateType:
		if t.IsSubtypeOf(s.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)) {
			return nil
		}
		return t
	default:
		if _, ok := t.(rhino.NoObjectType); ok {
			return s.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
		}
		if _, ok := t.(rhino.FunctionType); ok {
			return t
		}
		if proxy, ok := t.(rhino.ProxyObjectType); ok {
			return s.restrictToNotArray(proxy.GetReferencedTypeInternal())
		}
		if obj, ok := t.(rhino.ObjectType); ok {
			if obj.IsSubtypeOf(s.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)) {
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
