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

import "github.com/paullepoulpe/jsgo/internal/rhino"

// ClosureReverseAbstractInterpreter is a reverse abstract interpreter (RAI)
// for specific closure patterns such as goog.isObject.
type ClosureReverseAbstractInterpreter struct {
	ChainableReverseAbstractInterpreterBase

	// For when goog.isObject returns true. This includes functions, but not
	// null.
	restrictToObjectVisitor rhino.Visitor[rhino.JSType]

	// For when goog.isObject returns false.
	restrictToNotObjectVisitor rhino.Visitor[rhino.JSType]
}

// Java: ClosureReverseAbstractInterpreter(JSTypeRegistry)
func NewClosureReverseAbstractInterpreter(typeRegistry *rhino.JSTypeRegistry) *ClosureReverseAbstractInterpreter {
	c := &ClosureReverseAbstractInterpreter{}
	c.InitChainableReverseAbstractInterpreterBase(c, typeRegistry)
	c.restrictToObjectVisitor = newClosureReverseAbstractInterpreter_restrictToObjectVisitor(c)
	c.restrictToNotObjectVisitor = newClosureReverseAbstractInterpreter_restrictToNotObjectVisitor(c)
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
			if receiver.IsName() && receiver.GetString() == "goog" && callee.GetString() == "isObject" {
				return c.restrictIsObject(param, paramType, blindScope, outcome.IsTruthy())
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
		type_ = rhino.VisitJSType[rhino.JSType](type_, c.restrictToObjectVisitor)
	} else {
		type_ = rhino.VisitJSType[rhino.JSType](type_, c.restrictToNotObjectVisitor)
	}
	if type_ != nil {
		return c.DeclareNameInScope(blindScope, parameter, type_)
	}
	return blindScope
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
func (v *closureReverseAbstractInterpreter_restrictToObjectVisitor) CaseObjectType(type_ rhino.ObjectType) rhino.JSType {
	return type_
}

// Java: ClosureReverseAbstractInterpreter.restrictToObjectVisitor.caseFunctionType(FunctionType)
func (v *closureReverseAbstractInterpreter_restrictToObjectVisitor) CaseFunctionType(type_ rhino.FunctionType) rhino.JSType {
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
func (v *closureReverseAbstractInterpreter_restrictToNotObjectVisitor) CaseObjectType(type_ rhino.ObjectType) rhino.JSType {
	return nil
}

// Java: ClosureReverseAbstractInterpreter.restrictToNotObjectVisitor.caseFunctionType(FunctionType)
func (v *closureReverseAbstractInterpreter_restrictToNotObjectVisitor) CaseFunctionType(type_ rhino.FunctionType) rhino.JSType {
	return nil
}
