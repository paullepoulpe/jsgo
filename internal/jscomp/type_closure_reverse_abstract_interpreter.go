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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ClosureReverseAbstractInterpreter is a reverse abstract interpreter (RAI)
// for specific closure patterns such as goog.isObject.
type ClosureReverseAbstractInterpreter struct {
	ChainableReverseAbstractInterpreterBase

	// For when goog.isObject returns true. This includes functions, but not
	// null.
	restrictToObjectVisitor rhino.Visitor

	// For when goog.isObject returns false.
	restrictToNotObjectVisitor rhino.Visitor

	// Functions used to restrict types.
	restricters *jutil.LinkedHashMap[string, func(*ClosureReverseAbstractInterpreter_TypeRestriction) rhino.JSType]
}

// Java: ClosureReverseAbstractInterpreter(JSTypeRegistry)
func NewClosureReverseAbstractInterpreter(typeRegistry *rhino.JSTypeRegistry) *ClosureReverseAbstractInterpreter {
	c := &ClosureReverseAbstractInterpreter{}
	c.InitChainableReverseAbstractInterpreterBase(c, typeRegistry)
	// Field initializers.
	c.restrictToObjectVisitor = newClosureReverseAbstractInterpreter_restrictToObjectVisitor(c)
	c.restrictToNotObjectVisitor = newClosureReverseAbstractInterpreter_restrictToNotObjectVisitor(c)
	// Constructor body: ImmutableMap.of("isObject", p -> ...)
	c.restricters = jutil.NewLinkedHashMap[string, func(*ClosureReverseAbstractInterpreter_TypeRestriction) rhino.JSType]()
	c.restricters.Put(
		"isObject",
		func(p *ClosureReverseAbstractInterpreter_TypeRestriction) rhino.JSType {
			if p.type_ == nil {
				if p.outcome {
					return c.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)
				}
				return nil
			}

			var visitor rhino.Visitor
			if p.outcome {
				visitor = c.restrictToObjectVisitor
			} else {
				visitor = c.restrictToNotObjectVisitor
			}
			return rhino.JSType_VisitForJSType(p.type_, visitor)
		})
	return c
}

// Java: ClosureReverseAbstractInterpreter.getPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
func (c *ClosureReverseAbstractInterpreter) GetPreciserScopeKnowingConditionOutcome(
	condition rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	if condition.IsCall() && condition.HasTwoChildren() {
		callee := condition.GetFirstChild()
		param := condition.GetLastChild()
		if callee.IsGetProp() && param.IsQualifiedName() {
			paramType := c.GetTypeIfRefinable(param, blindScope)
			receiver := callee.GetFirstChild()
			if receiver.IsName() && receiver.GetString() == "goog" {
				restricter := c.restricters.Get(callee.GetString())
				if restricter != nil {
					return c.RestrictParameter(param, paramType, blindScope, restricter, outcome.IsTruthy())
				}
			}
		}
	}
	return c.NextPreciserScopeKnowingConditionOutcome(condition, blindScope, outcome)
}

// Java: ClosureReverseAbstractInterpreter.restrictParameter(Node, JSType, FlowScope, Function<TypeRestriction, JSType>, boolean)
func (c *ClosureReverseAbstractInterpreter) RestrictParameter(
	parameter rhino.Node,
	type_ rhino.JSType,
	blindScope FlowScope,
	restriction func(*ClosureReverseAbstractInterpreter_TypeRestriction) rhino.JSType,
	outcome bool) FlowScope {
	// restricting
	type_ = restriction(newClosureReverseAbstractInterpreter_TypeRestriction(type_, outcome))

	// changing the scope
	if type_ != nil {
		return c.DeclareNameInScope(blindScope, parameter, type_)
	} else {
		return blindScope
	}
}

// ClosureReverseAbstractInterpreter_TypeRestriction is
// ClosureReverseAbstractInterpreter.TypeRestriction.
type ClosureReverseAbstractInterpreter_TypeRestriction struct {
	type_   rhino.JSType
	outcome bool
}

// Java: TypeRestriction(JSType, boolean)
func newClosureReverseAbstractInterpreter_TypeRestriction(
	type_ rhino.JSType, outcome bool) *ClosureReverseAbstractInterpreter_TypeRestriction {
	t := &ClosureReverseAbstractInterpreter_TypeRestriction{}
	t.type_ = type_
	t.outcome = outcome
	return t
}

// closureReverseAbstractInterpreter_restrictToObjectVisitor is the anonymous
// RestrictByTrueTypeOfResultVisitor of restrictToObjectVisitor.
type closureReverseAbstractInterpreter_restrictToObjectVisitor struct {
	ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase
}

func newClosureReverseAbstractInterpreter_restrictToObjectVisitor(
	outer *ClosureReverseAbstractInterpreter) *closureReverseAbstractInterpreter_restrictToObjectVisitor {
	v := &closureReverseAbstractInterpreter_restrictToObjectVisitor{}
	v.InitChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase(v, outer)
	return v
}

// Java: ClosureReverseAbstractInterpreter.restrictToObjectVisitor.caseTopType(JSType)
func (v *closureReverseAbstractInterpreter_restrictToObjectVisitor) CaseTopType(topType rhino.JSType) rhino.JSType {
	return v.outer.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
}

// Java: ClosureReverseAbstractInterpreter.restrictToObjectVisitor.caseObjectType(ObjectType)
func (v *closureReverseAbstractInterpreter_restrictToObjectVisitor) CaseObjectType(type_ rhino.ObjectType) any {
	return type_
}

// Java: ClosureReverseAbstractInterpreter.restrictToObjectVisitor.caseFunctionType(FunctionType)
func (v *closureReverseAbstractInterpreter_restrictToObjectVisitor) CaseFunctionType(type_ rhino.FunctionType) any {
	return type_
}

// closureReverseAbstractInterpreter_restrictToNotObjectVisitor is the
// anonymous RestrictByFalseTypeOfResultVisitor of
// restrictToNotObjectVisitor.
type closureReverseAbstractInterpreter_restrictToNotObjectVisitor struct {
	ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase
}

func newClosureReverseAbstractInterpreter_restrictToNotObjectVisitor(
	outer *ClosureReverseAbstractInterpreter) *closureReverseAbstractInterpreter_restrictToNotObjectVisitor {
	v := &closureReverseAbstractInterpreter_restrictToNotObjectVisitor{}
	v.InitChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase(v, outer)
	return v
}

// Java: ClosureReverseAbstractInterpreter.restrictToNotObjectVisitor.caseObjectType(ObjectType)
func (v *closureReverseAbstractInterpreter_restrictToNotObjectVisitor) CaseObjectType(type_ rhino.ObjectType) any {
	return nil
}

// Java: ClosureReverseAbstractInterpreter.restrictToNotObjectVisitor.caseFunctionType(FunctionType)
func (v *closureReverseAbstractInterpreter_restrictToNotObjectVisitor) CaseFunctionType(type_ rhino.FunctionType) any {
	return nil
}
