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

// ChainableReverseAbstractInterpreter is a chainable reverse abstract
// interpreter providing basic functionality.
//
// Java: abstract class ChainableReverseAbstractInterpreter implements
// ReverseAbstractInterpreter. The Go interface declares every overridable
// method; ChainableReverseAbstractInterpreterBase holds the fields and the
// base implementations and dispatches virtual calls through self
// (PORTING.md §4).
type ChainableReverseAbstractInterpreter interface {
	ReverseAbstractInterpreter

	// Java: ChainableReverseAbstractInterpreter.append(ChainableReverseAbstractInterpreter)
	Append(lastLink ChainableReverseAbstractInterpreter) ChainableReverseAbstractInterpreter
	// Java: ChainableReverseAbstractInterpreter.getFirst()
	GetFirst() ChainableReverseAbstractInterpreter
	// Java: ChainableReverseAbstractInterpreter.firstPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
	FirstPreciserScopeKnowingConditionOutcome(
		condition *rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope
	// Java: ChainableReverseAbstractInterpreter.nextPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
	NextPreciserScopeKnowingConditionOutcome(
		condition *rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope
	// Java: ChainableReverseAbstractInterpreter.getTypeIfRefinable(Node, FlowScope)
	GetTypeIfRefinable(node *rhino.Node, scope FlowScope) rhino.JSType
	// Java: ChainableReverseAbstractInterpreter.declareNameInScope(FlowScope, Node, JSType)
	DeclareNameInScope(scope FlowScope, node *rhino.Node, type_ rhino.JSType) FlowScope
	// Java: ChainableReverseAbstractInterpreter.getRestrictedByTypeOfResult(JSType, String, boolean)
	GetRestrictedByTypeOfResult(type_ rhino.JSType, value string, resultEqualsValue bool) rhino.JSType
	// Java: ChainableReverseAbstractInterpreter.getNativeType(JSTypeNative)
	GetNativeType(typeId rhino.JSTypeNative) rhino.JSType
	// Java: ChainableReverseAbstractInterpreter.getNativeTypeForTypeOf(String)
	GetNativeTypeForTypeOf(value string) rhino.JSType

	// GetChainableReverseAbstractInterpreterBase gives access to the fields
	// through the interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetChainableReverseAbstractInterpreterBase() *ChainableReverseAbstractInterpreterBase
}

// ChainableReverseAbstractInterpreterBase holds
// ChainableReverseAbstractInterpreter's fields and base methods.
type ChainableReverseAbstractInterpreterBase struct {
	self ChainableReverseAbstractInterpreter

	typeRegistry *rhino.JSTypeRegistry
	firstLink    ChainableReverseAbstractInterpreter
	nextLink     ChainableReverseAbstractInterpreter

	// For when goog.isArray or Array.isArray returns true.
	restrictToArrayVisitor rhino.Visitor

	// For when goog.isArray or Array.isArray returns false.
	restrictToNotArrayVisitor rhino.Visitor
}

// InitChainableReverseAbstractInterpreterBase constructs an interpreter,
// which is the only link in a chain. Interpreters can be appended using
// Append.
//
// Java: ChainableReverseAbstractInterpreter(JSTypeRegistry)
func (b *ChainableReverseAbstractInterpreterBase) InitChainableReverseAbstractInterpreterBase(
	self ChainableReverseAbstractInterpreter, typeRegistry *rhino.JSTypeRegistry) {
	b.self = self
	// Field initializers.
	b.restrictToArrayVisitor = newChainableReverseAbstractInterpreter_restrictToArrayVisitor(self)
	b.restrictToNotArrayVisitor = newChainableReverseAbstractInterpreter_restrictToNotArrayVisitor(self)
	// Constructor body.
	b.typeRegistry = typeRegistry
	b.firstLink = self
	b.nextLink = nil
}

// GetChainableReverseAbstractInterpreterBase gives access to the fields
// through the interface.
func (b *ChainableReverseAbstractInterpreterBase) GetChainableReverseAbstractInterpreterBase() *ChainableReverseAbstractInterpreterBase {
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
	jutil.CheckArgument(lastLink.GetChainableReverseAbstractInterpreterBase().nextLink == nil)
	b.nextLink = lastLink
	lastLink.GetChainableReverseAbstractInterpreterBase().firstLink = b.firstLink
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
	condition *rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	return b.firstLink.GetPreciserScopeKnowingConditionOutcome(
		condition, blindScope, outcome)
}

// NextPreciserScopeKnowingConditionOutcome delegates the calculation of the
// preciser scope to the next link. If there is no next link, returns the
// blind scope.
//
// Java: ChainableReverseAbstractInterpreter.nextPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
func (b *ChainableReverseAbstractInterpreterBase) NextPreciserScopeKnowingConditionOutcome(
	condition *rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
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
func (b *ChainableReverseAbstractInterpreterBase) GetTypeIfRefinable(node *rhino.Node, scope FlowScope) rhino.JSType {
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
			propVarType = b.self.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
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
	scope FlowScope, node *rhino.Node, type_ rhino.JSType) FlowScope {
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
			origType = b.self.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		return scope.InferQualifiedSlot(node, qualifiedName, origType, type_, false)
	case rhino.Token_THIS:
		// "this" references aren't currently modeled in the CFG.
		return scope
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Node cannot be refined. \n" + node.ToStringTree()))
	}
}

// ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor is a
// class common to all visitors that need to restrict the type based on
// typeof-like conditions.
//
// Java: abstract class ChainableReverseAbstractInterpreter.RestrictByTypeOfResultVisitor
// implements Visitor<JSType> (an inner class: it uses the enclosing
// interpreter's getNativeType).
type ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor interface {
	rhino.Visitor

	// CaseTopType abstracts away the similarities between visiting the
	// unknown type and the all type. topType is UNKNOWN_TYPE or ALL_TYPE;
	// returns the restricted type.
	//
	// Java: RestrictByTypeOfResultVisitor.caseTopType(JSType)
	CaseTopType(topType rhino.JSType) rhino.JSType
}

// ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase
// holds the base implementations; self is the concrete visitor and outer
// the enclosing interpreter.
type ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase struct {
	self  ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor
	outer ChainableReverseAbstractInterpreter
}

// Java: ChainableReverseAbstractInterpreter.RestrictByTypeOfResultVisitor()
func (v *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) InitChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase(
	self ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor, outer ChainableReverseAbstractInterpreter) {
	v.self = self
	v.outer = outer
}

// Java: RestrictByTypeOfResultVisitor.caseAllType()
func (v *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseAllType() any {
	return v.self.CaseTopType(v.outer.GetNativeType(rhino.JSTypeNative_ALL_TYPE))
}

// Java: RestrictByTypeOfResultVisitor.caseUnknownType()
func (v *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseUnknownType() any {
	return v.self.CaseTopType(v.outer.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE))
}

// Java: RestrictByTypeOfResultVisitor.caseUnionType(UnionType)
func (v *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseUnionType(type_ *rhino.UnionType) any {
	var restricted rhino.JSType
	for _, alternate := range type_.GetAlternates() {
		restrictedAlternate := rhino.JSType_VisitForJSType(alternate, v.self)
		if restrictedAlternate != nil {
			if restricted == nil {
				restricted = restrictedAlternate
			} else {
				restricted = restrictedAlternate.GetLeastSupertype(restricted)
			}
		}
	}
	return restricted
}

// Java: RestrictByTypeOfResultVisitor.caseNoType(NoType)
func (v *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseNoType(type_ *rhino.NoType) any {
	return type_
}

// Java: RestrictByTypeOfResultVisitor.caseEnumElementType(EnumElementType)
func (v *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseEnumElementType(
	enumElementType *rhino.EnumElementType) any {
	// NOTE(nicksantos): This is a white lie. Suppose we have:
	// /** @enum {string|number} */ var MyEnum = ...;
	// if (goog.isNumber(myEnumInstance)) {
	//   /* what is myEnumInstance here? */
	// }
	// There is no type that represents {MyEnum - string}. What we really
	// need is a notion of "enum subtyping", so that we could dynamically
	// create a subtype of MyEnum restricted by string. In any case,
	// this should catch the common case.
	type_ := rhino.JSType_VisitForJSType(enumElementType.GetPrimitiveType(), v.self)
	if type_ != nil && enumElementType.GetPrimitiveType().Equals(type_) {
		return enumElementType
	} else {
		return type_
	}
}

// Java: RestrictByTypeOfResultVisitor.caseTemplatizedType(TemplatizedType)
func (v *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseTemplatizedType(type_ *rhino.TemplatizedType) any {
	return v.self.CaseObjectType(type_)
}

// Java: RestrictByTypeOfResultVisitor.caseTemplateType(TemplateType)
func (v *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseTemplateType(templateType *rhino.TemplateType) any {
	return v.self.CaseObjectType(templateType)
}

// Java: RestrictByTypeOfResultVisitor.caseNamedType(NamedType)
func (v *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseNamedType(type_ *rhino.NamedType) any {
	return v.self.CaseProxyObjectType(type_)
}

// Java: RestrictByTypeOfResultVisitor.caseProxyObjectType(ProxyObjectType)
func (v *ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase) CaseProxyObjectType(type_ rhino.ProxyObjectType) any {
	return type_.VisitReferenceType(v.self)
}

// ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor is a
// class common to all visitors that need to restrict the type based on some
// typeof-like condition being true. All base cases return null. It is up to
// the subclasses to override the appropriate ones.
//
// Java: abstract class ChainableReverseAbstractInterpreter.RestrictByTrueTypeOfResultVisitor
// extends RestrictByTypeOfResultVisitor
type ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor interface {
	ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor
}

// ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase
// holds RestrictByTrueTypeOfResultVisitor's base implementations.
type ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase struct {
	ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase
}

// Java: ChainableReverseAbstractInterpreter.RestrictByTrueTypeOfResultVisitor()
func (v *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) InitChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase(
	self ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor, outer ChainableReverseAbstractInterpreter) {
	v.InitChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase(self, outer)
}

// Java: RestrictByTrueTypeOfResultVisitor.caseNoObjectType()
func (v *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseNoObjectType() any {
	return nil
}

// Java: RestrictByTrueTypeOfResultVisitor.caseBooleanType()
func (v *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseBooleanType() any {
	return nil
}

// Java: RestrictByTrueTypeOfResultVisitor.caseFunctionType(FunctionType)
func (v *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseFunctionType(type_ rhino.FunctionType) any {
	return nil
}

// Java: RestrictByTrueTypeOfResultVisitor.caseNullType()
func (v *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseNullType() any {
	return nil
}

// Java: RestrictByTrueTypeOfResultVisitor.caseNumberType()
func (v *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseNumberType() any {
	return nil
}

// Java: RestrictByTrueTypeOfResultVisitor.caseBigIntType()
func (v *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseBigIntType() any {
	return nil
}

// Java: RestrictByTrueTypeOfResultVisitor.caseObjectType(ObjectType)
func (v *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseObjectType(type_ rhino.ObjectType) any {
	return nil
}

// Java: RestrictByTrueTypeOfResultVisitor.caseStringType()
func (v *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseStringType() any {
	return nil
}

// Java: RestrictByTrueTypeOfResultVisitor.caseSymbolType(SymbolType)
func (v *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseSymbolType(type_ rhino.SymbolType) any {
	return nil
}

// Java: RestrictByTrueTypeOfResultVisitor.caseVoidType()
func (v *ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase) CaseVoidType() any {
	return nil
}

// ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor is
// a class common to all visitors that need to restrict the type based on
// some typeof-like condition being false. All base cases return their type.
// It is up to the subclasses to override the appropriate ones.
//
// Java: abstract class ChainableReverseAbstractInterpreter.RestrictByFalseTypeOfResultVisitor
// extends RestrictByTypeOfResultVisitor
type ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor interface {
	ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor
}

// ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase
// holds RestrictByFalseTypeOfResultVisitor's base implementations.
type ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase struct {
	ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase
}

// Java: ChainableReverseAbstractInterpreter.RestrictByFalseTypeOfResultVisitor()
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) InitChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase(
	self ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor, outer ChainableReverseAbstractInterpreter) {
	v.InitChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase(self, outer)
}

// Java: RestrictByFalseTypeOfResultVisitor.caseTopType(JSType)
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseTopType(topType rhino.JSType) rhino.JSType {
	return topType
}

// Java: RestrictByFalseTypeOfResultVisitor.caseNoObjectType()
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseNoObjectType() any {
	return v.outer.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
}

// Java: RestrictByFalseTypeOfResultVisitor.caseBooleanType()
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseBooleanType() any {
	return v.outer.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
}

// Java: RestrictByFalseTypeOfResultVisitor.caseFunctionType(FunctionType)
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseFunctionType(type_ rhino.FunctionType) any {
	return type_
}

// Java: RestrictByFalseTypeOfResultVisitor.caseNullType()
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseNullType() any {
	return v.outer.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
}

// Java: RestrictByFalseTypeOfResultVisitor.caseNumberType()
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseNumberType() any {
	return v.outer.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: RestrictByFalseTypeOfResultVisitor.caseBigIntType()
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseBigIntType() any {
	return v.outer.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
}

// Java: RestrictByFalseTypeOfResultVisitor.caseObjectType(ObjectType)
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseObjectType(type_ rhino.ObjectType) any {
	return type_
}

// Java: RestrictByFalseTypeOfResultVisitor.caseStringType()
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseStringType() any {
	return v.outer.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
}

// Java: RestrictByFalseTypeOfResultVisitor.caseSymbolType(SymbolType)
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseSymbolType(type_ rhino.SymbolType) any {
	return type_
}

// Java: RestrictByFalseTypeOfResultVisitor.caseVoidType()
func (v *ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase) CaseVoidType() any {
	return v.outer.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
}

// ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor: see
// ChainableReverseAbstractInterpreter.GetRestrictedByTypeOfResult.
//
// Java: private class ChainableReverseAbstractInterpreter.RestrictByOneTypeOfResultVisitor
// extends RestrictByTypeOfResultVisitor
type ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor struct {
	ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase

	// A value known to be equal or not equal to the result of the typeOf
	// operation.
	value string

	// true if the typeOf result is known to equal value; false if it is
	// known not to equal value.
	resultEqualsValue bool
}

// Java: RestrictByOneTypeOfResultVisitor(String, boolean)
func newChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor(
	outer ChainableReverseAbstractInterpreter, value string, resultEqualsValue bool) *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor {
	v := &ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor{}
	v.InitChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitorBase(v, outer)
	v.value = value
	v.resultEqualsValue = resultEqualsValue
	return v
}

// MatchesExpectation computes whether the given result of a typeof operator
// matches expectations, i.e. whether a type that gives such a result should
// be kept.
//
// Java: RestrictByOneTypeOfResultVisitor.matchesExpectation(String)
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) MatchesExpectation(result string) bool {
	return (result == v.value) == v.resultEqualsValue
}

// Java: RestrictByOneTypeOfResultVisitor.caseTopType(JSType)
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseTopType(topType rhino.JSType) rhino.JSType {
	result := topType
	if v.resultEqualsValue {
		typeByName := v.outer.GetNativeTypeForTypeOf(v.value)
		if typeByName != nil {
			result = typeByName
		}
	}
	return result
}

// Java: RestrictByOneTypeOfResultVisitor.caseNoObjectType()
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseNoObjectType() any {
	if (v.value == "object" || v.value == "function") == v.resultEqualsValue {
		return v.outer.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
	}
	return nil
}

// Java: RestrictByOneTypeOfResultVisitor.caseBooleanType()
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseBooleanType() any {
	if v.MatchesExpectation("boolean") {
		return v.outer.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	}
	return nil
}

// Java: RestrictByOneTypeOfResultVisitor.caseFunctionType(FunctionType)
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseFunctionType(type_ rhino.FunctionType) any {
	if v.MatchesExpectation("function") {
		return type_
	}
	return nil
}

// Java: RestrictByOneTypeOfResultVisitor.caseNullType()
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseNullType() any {
	if v.MatchesExpectation("object") {
		return v.outer.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
	}
	return nil
}

// Java: RestrictByOneTypeOfResultVisitor.caseNumberType()
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseNumberType() any {
	if v.MatchesExpectation("number") {
		return v.outer.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	}
	return nil
}

// Java: RestrictByOneTypeOfResultVisitor.caseBigIntType()
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseBigIntType() any {
	if v.MatchesExpectation("bigint") {
		return v.outer.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
	}
	return nil
}

// Java: RestrictByOneTypeOfResultVisitor.caseObjectType(ObjectType)
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseObjectType(type_ rhino.ObjectType) any {
	if v.value == "function" {
		ctorType := v.outer.GetNativeType(rhino.JSTypeNative_FUNCTION_TYPE)
		if v.resultEqualsValue {
			// Objects are restricted to "Function", subtypes are left
			return ctorType.GetGreatestSubtype(type_)
		} else {
			// Only filter out subtypes of "function"
			if type_.IsSubtypeOf(ctorType) {
				return nil
			}
			return type_
		}
	}
	if v.MatchesExpectation("object") {
		return type_
	}
	return nil
}

// Java: RestrictByOneTypeOfResultVisitor.caseStringType()
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseStringType() any {
	if v.MatchesExpectation("string") {
		return v.outer.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	}
	return nil
}

// Java: RestrictByOneTypeOfResultVisitor.caseSymbolType(SymbolType)
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseSymbolType(type_ rhino.SymbolType) any {
	if v.MatchesExpectation("symbol") {
		return type_
	}
	return nil
}

// Java: RestrictByOneTypeOfResultVisitor.caseVoidType()
func (v *ChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor) CaseVoidType() any {
	if v.MatchesExpectation("undefined") {
		return v.outer.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	}
	return nil
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
			result := b.self.GetNativeTypeForTypeOf(value)
			if result == nil {
				return b.self.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
			}
			return result
		} else {
			return nil
		}
	}
	return rhino.JSType_VisitForJSType(
		type_, newChainableReverseAbstractInterpreter_RestrictByOneTypeOfResultVisitor(b.self, value, resultEqualsValue))
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
		return b.self.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	case "boolean":
		return b.self.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	case "string":
		return b.self.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	case "symbol":
		return b.self.GetNativeType(rhino.JSTypeNative_SYMBOL_TYPE)
	case "undefined":
		return b.self.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	case "object":
		// NOTE: This is broader than it needs to be if it's from goog.typeof, but (a) it's more
		// consistent with common usage of the native builtin typeof, (b) it's more consistent
		// with
		// TypeScript, and (c) it's more useful than simply not narrowing.
		return b.typeRegistry.CreateUnionType(
			b.self.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE), b.self.GetNativeType(rhino.JSTypeNative_NULL_TYPE))
	case "function":
		return b.self.GetNativeType(rhino.JSTypeNative_FUNCTION_TYPE)
	default:
		return nil
	}
}

// chainableReverseAbstractInterpreter_restrictToArrayVisitor is the
// anonymous RestrictByTrueTypeOfResultVisitor of restrictToArrayVisitor:
// for when goog.isArray or Array.isArray returns true.
type chainableReverseAbstractInterpreter_restrictToArrayVisitor struct {
	ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase
}

func newChainableReverseAbstractInterpreter_restrictToArrayVisitor(
	outer ChainableReverseAbstractInterpreter) *chainableReverseAbstractInterpreter_restrictToArrayVisitor {
	v := &chainableReverseAbstractInterpreter_restrictToArrayVisitor{}
	v.InitChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase(v, outer)
	return v
}

// Java: ChainableReverseAbstractInterpreter.restrictToArrayVisitor.caseTopType(JSType)
func (v *chainableReverseAbstractInterpreter_restrictToArrayVisitor) CaseTopType(topType rhino.JSType) rhino.JSType {
	if topType.IsAllType() {
		return v.outer.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)
	}
	return topType
}

// Java: ChainableReverseAbstractInterpreter.restrictToArrayVisitor.caseObjectType(ObjectType)
func (v *chainableReverseAbstractInterpreter_restrictToArrayVisitor) CaseObjectType(type_ rhino.ObjectType) any {
	arrayType := v.outer.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)
	if arrayType.IsSubtypeOf(type_) {
		return arrayType
	}
	return nil
}

// chainableReverseAbstractInterpreter_restrictToNotArrayVisitor is the
// anonymous RestrictByFalseTypeOfResultVisitor of restrictToNotArrayVisitor:
// for when goog.isArray or Array.isArray returns false.
type chainableReverseAbstractInterpreter_restrictToNotArrayVisitor struct {
	ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase
}

func newChainableReverseAbstractInterpreter_restrictToNotArrayVisitor(
	outer ChainableReverseAbstractInterpreter) *chainableReverseAbstractInterpreter_restrictToNotArrayVisitor {
	v := &chainableReverseAbstractInterpreter_restrictToNotArrayVisitor{}
	v.InitChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase(v, outer)
	return v
}

// Java: ChainableReverseAbstractInterpreter.restrictToNotArrayVisitor.caseObjectType(ObjectType)
func (v *chainableReverseAbstractInterpreter_restrictToNotArrayVisitor) CaseObjectType(type_ rhino.ObjectType) any {
	if type_.IsSubtypeOf(v.outer.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE)) {
		return nil
	}
	return type_
}
