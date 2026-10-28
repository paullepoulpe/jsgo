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
)

// ClosureReverseAbstractInterpreter is a reverse abstract interpreter (RAI)
// for specific closure patterns such as goog.isObject.
type ClosureReverseAbstractInterpreter struct {
	ChainableReverseAbstractInterpreterBase
}

// Java: ClosureReverseAbstractInterpreter(JSTypeRegistry)
func NewClosureReverseAbstractInterpreter(typeRegistry *rhino.JSTypeRegistry) *ClosureReverseAbstractInterpreter {
	c := &ClosureReverseAbstractInterpreter{}
	c.initChainableBase(c, typeRegistry)
	return c
}

// Java: ClosureReverseAbstractInterpreter.getPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
func (c *ClosureReverseAbstractInterpreter) GetPreciserScopeKnowingConditionOutcome(
	condition rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	if call, ok := ast.AsCall(condition); ok {
		if param, ok := call.OnlyArg(); ok {
			if callee, ok := call.Callee().AsGetProp(); ok && param.Node().IsQualifiedName() {
				paramType := c.GetTypeIfRefinable(param.Node(), blindScope)
				if receiver, ok := callee.Object().AsName(); ok && receiver.Ident() == "goog" && callee.Property() == "isObject" {
					return c.restrictIsObject(param.Node(), paramType, blindScope, outcome.IsTruthy())
				}
			}
		}
	}
	return c.NextPreciserScopeKnowingConditionOutcome(condition, blindScope, outcome)
}

func (c *ClosureReverseAbstractInterpreter) restrictIsObject(
	parameter rhino.Node,
	type_ rhino.JSType,
	blindScope FlowScope,
	outcome bool) FlowScope {
	if type_ == nil {
		if outcome {
			type_ = c.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)
		}
	} else if outcome {
		type_ = c.restrictToObject(type_)
	} else {
		type_ = c.restrictToNotObject(type_)
	}
	if type_ != nil {
		return c.DeclareNameInScope(blindScope, parameter, type_)
	}
	return blindScope
}

// restrictToObject restricts t when goog.isObject returns true (includes
// functions, excludes null).
//
// Java: ClosureReverseAbstractInterpreter.restrictToObjectVisitor
func (c *ClosureReverseAbstractInterpreter) restrictToObject(t rhino.JSType) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.AllType, *rhino.UnknownType:
		return c.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := c.restrictToObject(prim)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.UnionType:
		var restricted rhino.JSType
		for _, alt := range x.GetAlternates() {
			if ra := c.restrictToObject(alt); ra != nil {
				if restricted == nil {
					restricted = ra
				} else {
					restricted = ra.GetLeastSupertype(restricted)
				}
			}
		}
		return restricted
	case *rhino.NamedType:
		return c.restrictToObject(x.GetReferencedTypeInternal())
	case *rhino.ProxyObjectTypeBase:
		return c.restrictToObject(x.GetReferencedTypeInternal())
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
			return c.restrictToObject(proxy.GetReferencedTypeInternal())
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
func (c *ClosureReverseAbstractInterpreter) restrictToNotObject(t rhino.JSType) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.AllType:
		return c.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
	case *rhino.UnknownType:
		return c.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
	case *rhino.NoObjectTypeBase:
		return c.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
	case *rhino.BooleanType:
		return c.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	case *rhino.NullType:
		return c.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
	case *rhino.NumberType:
		return c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	case *rhino.BigIntType:
		return c.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
	case *rhino.StringType:
		return c.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	case *rhino.VoidType:
		return c.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	case *rhino.SymbolTypeBase, *rhino.KnownSymbolType:
		return t
	case *rhino.FunctionTypeBase, *rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase,
		*rhino.RecordType, *rhino.EnumType, *rhino.NoResolvedType, *rhino.TemplatizedType, *rhino.TemplateType:
		return nil
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := c.restrictToNotObject(prim)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.UnionType:
		var restricted rhino.JSType
		for _, alt := range x.GetAlternates() {
			if ra := c.restrictToNotObject(alt); ra != nil {
				if restricted == nil {
					restricted = ra
				} else {
					restricted = ra.GetLeastSupertype(restricted)
				}
			}
		}
		return restricted
	case *rhino.NamedType:
		return c.restrictToNotObject(x.GetReferencedTypeInternal())
	case *rhino.ProxyObjectTypeBase:
		return c.restrictToNotObject(x.GetReferencedTypeInternal())
	default:
		if _, ok := t.(rhino.NoObjectType); ok {
			return c.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
		}
		if _, ok := t.(rhino.FunctionType); ok {
			return nil
		}
		if proxy, ok := t.(rhino.ProxyObjectType); ok {
			return c.restrictToNotObject(proxy.GetReferencedTypeInternal())
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
