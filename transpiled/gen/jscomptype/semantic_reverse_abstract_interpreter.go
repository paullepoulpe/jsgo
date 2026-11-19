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
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java (generated mechanically by j2go).

package jscomptype

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_semantic_reverse_abstract_interpreter_1 SemanticReverseAbstractInterpreter_TypePairOp
var lambdaSingleton_semantic_reverse_abstract_interpreter_2 SemanticReverseAbstractInterpreter_TypePairOp
var lambdaSingleton_semantic_reverse_abstract_interpreter_3 SemanticReverseAbstractInterpreter_TypePairOp
var lambdaSingleton_semantic_reverse_abstract_interpreter_4 SemanticReverseAbstractInterpreter_TypePairOp
var lambdaSingleton_semantic_reverse_abstract_interpreter_5 SemanticReverseAbstractInterpreter_TypePairOp
var lambdaSingleton_semantic_reverse_abstract_interpreter_6 SemanticReverseAbstractInterpreter_TypePairOp
var lambdaSingleton_semantic_reverse_abstract_interpreter_7 SemanticReverseAbstractInterpreter_TypePairOp
var lambdaSingleton_semantic_reverse_abstract_interpreter_8 SemanticReverseAbstractInterpreter_TypePairOp
var lambdaSingleton_semantic_reverse_abstract_interpreter_9 SemanticReverseAbstractInterpreter_TypePairOp
var lambdaSingleton_semantic_reverse_abstract_interpreter_10 SemanticReverseAbstractInterpreter_TypePairOp
var lambdaSingleton_semantic_reverse_abstract_interpreter_11 SemanticReverseAbstractInterpreter_TypePairOp

func init() {
	lambdaSingleton_semantic_reverse_abstract_interpreter_1 = NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(func(mrp0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
		return SemanticReverseAbstractInterpreter_Eq(mrp0_)
	})
	lambdaSingleton_semantic_reverse_abstract_interpreter_2 = NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(func(mrp0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
		return SemanticReverseAbstractInterpreter_Ne(mrp0_)
	})
	lambdaSingleton_semantic_reverse_abstract_interpreter_3 = NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(func(mrp0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
		return SemanticReverseAbstractInterpreter_Ne(mrp0_)
	})
	lambdaSingleton_semantic_reverse_abstract_interpreter_4 = NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(func(mrp0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
		return SemanticReverseAbstractInterpreter_Eq(mrp0_)
	})
	lambdaSingleton_semantic_reverse_abstract_interpreter_5 = NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(func(mrp0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
		return SemanticReverseAbstractInterpreter_Sheq(mrp0_)
	})
	lambdaSingleton_semantic_reverse_abstract_interpreter_6 = NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(func(mrp0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
		return SemanticReverseAbstractInterpreter_Shne(mrp0_)
	})
	lambdaSingleton_semantic_reverse_abstract_interpreter_7 = NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(func(mrp0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
		return SemanticReverseAbstractInterpreter_Shne(mrp0_)
	})
	lambdaSingleton_semantic_reverse_abstract_interpreter_8 = NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(func(mrp0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
		return SemanticReverseAbstractInterpreter_Sheq(mrp0_)
	})
	lambdaSingleton_semantic_reverse_abstract_interpreter_9 = NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(func(mrp0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
		return SemanticReverseAbstractInterpreter_Ineq(mrp0_)
	})
	lambdaSingleton_semantic_reverse_abstract_interpreter_10 = NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(func(mrp0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
		return SemanticReverseAbstractInterpreter_Sheq(mrp0_)
	})
	lambdaSingleton_semantic_reverse_abstract_interpreter_11 = NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(func(mrp0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
		return SemanticReverseAbstractInterpreter_Shne(mrp0_)
	})
}

// Java class: com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter
type SemanticReverseAbstractInterpreter struct {
	ChainableReverseAbstractInterpreterBase
}

func (this_ *SemanticReverseAbstractInterpreter) BindSelf_(s any) {
	_ = s
	this_.ChainableReverseAbstractInterpreterBase.BindSelf_(s)
}
func (this_ *SemanticReverseAbstractInterpreter) IsJavaNil_() bool { return this_ == nil }
func (*SemanticReverseAbstractInterpreter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter", "SemanticReverseAbstractInterpreter"
}
func (this_ *SemanticReverseAbstractInterpreter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ChainableReverseAbstractInterpreterBase.IdentityHash_()
}
func (*SemanticReverseAbstractInterpreter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter", "com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter", func(v_ any) bool {
		return j2gort.InstanceOf[*SemanticReverseAbstractInterpreter](v_)
	})
}
func (*SemanticReverseAbstractInterpreter) Implements_jscomptype_ChainableReverseAbstractInterpreter() {
}
func (*SemanticReverseAbstractInterpreter) Implements_jscomptype_ReverseAbstractInterpreter() {}

func (this_ *SemanticReverseAbstractInterpreter) initFields_SemanticReverseAbstractInterpreter() {
	_ = this_
}

// Java constructor: SemanticReverseAbstractInterpreter(com.google.javascript.rhino.jstype.JSTypeRegistry typeRegistry)
func NewSemanticReverseAbstractInterpreter(typeRegistry *pkg_rhino.JSTypeRegistry) *SemanticReverseAbstractInterpreter {
	this_ := &SemanticReverseAbstractInterpreter{}
	this_.BindSelf_(this_)
	InitSemanticReverseAbstractInterpreter(this_, typeRegistry)
	return this_
}

func InitSemanticReverseAbstractInterpreter(this_ *SemanticReverseAbstractInterpreter, typeRegistry *pkg_rhino.JSTypeRegistry) {
	_ = typeRegistry
	InitChainableReverseAbstractInterpreter(&this_.ChainableReverseAbstractInterpreterBase, typeRegistry)
	this_.initFields_SemanticReverseAbstractInterpreter()
}

// Java static method: SemanticReverseAbstractInterpreter.eq(com.google.javascript.rhino.jstype.JSType.TypePair p)
func SemanticReverseAbstractInterpreter_Eq(p pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
	_ = p
	if j2gort.IfaceEquals[pkg_rhino.JSType](p.FieldGet_typeA(), nil) || j2gort.IfaceEquals[pkg_rhino.JSType](p.FieldGet_typeB(), nil) {
		return (pkg_rhino.JSType_TypePair{})
	}
	return p.FieldGet_typeA().GetTypesUnderEquality(p.FieldGet_typeB())
}

// Java static method: SemanticReverseAbstractInterpreter.ne(com.google.javascript.rhino.jstype.JSType.TypePair p)
func SemanticReverseAbstractInterpreter_Ne(p pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
	_ = p
	if j2gort.IfaceEquals[pkg_rhino.JSType](p.FieldGet_typeA(), nil) || j2gort.IfaceEquals[pkg_rhino.JSType](p.FieldGet_typeB(), nil) {
		return (pkg_rhino.JSType_TypePair{})
	}
	return p.FieldGet_typeA().GetTypesUnderInequality(p.FieldGet_typeB())
}

// Java static method: SemanticReverseAbstractInterpreter.sheq(com.google.javascript.rhino.jstype.JSType.TypePair p)
func SemanticReverseAbstractInterpreter_Sheq(p pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
	_ = p
	if j2gort.IfaceEquals[pkg_rhino.JSType](p.FieldGet_typeA(), nil) || j2gort.IfaceEquals[pkg_rhino.JSType](p.FieldGet_typeB(), nil) {
		return (pkg_rhino.JSType_TypePair{})
	}
	return p.FieldGet_typeA().GetJSTypeBase().GetTypesUnderShallowEquality(p.FieldGet_typeB())
}

// Java static method: SemanticReverseAbstractInterpreter.shne(com.google.javascript.rhino.jstype.JSType.TypePair p)
func SemanticReverseAbstractInterpreter_Shne(p pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
	_ = p
	if j2gort.IfaceEquals[pkg_rhino.JSType](p.FieldGet_typeA(), nil) || j2gort.IfaceEquals[pkg_rhino.JSType](p.FieldGet_typeB(), nil) {
		return (pkg_rhino.JSType_TypePair{})
	}
	return p.FieldGet_typeA().GetTypesUnderShallowInequality(p.FieldGet_typeB())
}

// Java static method: SemanticReverseAbstractInterpreter.ineq(com.google.javascript.rhino.jstype.JSType.TypePair p)
func SemanticReverseAbstractInterpreter_Ineq(p pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
	_ = p
	return pkg_rhino.NewJSType_TypePair(func() pkg_rhino.JSType {
		if !(j2gort.IfaceEquals[pkg_rhino.JSType](p.FieldGet_typeA(), nil)) {
			return p.FieldGet_typeA().RestrictByNotUndefined()
		}
		return nil
	}(), func() pkg_rhino.JSType {
		if !(j2gort.IfaceEquals[pkg_rhino.JSType](p.FieldGet_typeB(), nil)) {
			return p.FieldGet_typeB().RestrictByNotUndefined()
		}
		return nil
	}())
}

// Java method: SemanticReverseAbstractInterpreter.getPreciserScopeKnowingConditionOutcome(com.google.javascript.rhino.Node condition, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.rhino.Outcome outcome)
func (this_ *SemanticReverseAbstractInterpreter) GetPreciserScopeKnowingConditionOutcome(condition pkg_rhino.Node, blindScope FlowScope, outcome *pkg_rhino.Outcome) FlowScope {
	_ = condition
	_ = blindScope
	_ = outcome
	var operatorToken pkg_rhino.Token = condition.GetToken()
	_ = operatorToken
	switch (operatorToken).Ordinal() {
	case int32(4) /* Token.EQ */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */, int32(76) /* Token.CASE */ :
		var left pkg_rhino.Node = pkg_rhino.Node(0)
		_ = left
		var right pkg_rhino.Node = pkg_rhino.Node(0)
		_ = right
		if operatorToken == pkg_rhino.StaticGet_Token_CASE() {
			left = condition.GetParent().GetPrevious()
			if !(condition.GetParent().IsSwitchBody()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(condition.GetParent()))
			}
			right = condition.GetFirstChild()
		} else {
			left = condition.GetFirstChild()
			right = condition.GetLastChild()
		}
		var typeOfNode pkg_rhino.Node = pkg_rhino.Node(0)
		_ = typeOfNode
		var stringNode pkg_rhino.Node = pkg_rhino.Node(0)
		_ = stringNode
		if left.IsTypeOf() && right.IsStringLit() {
			typeOfNode = left
			stringNode = right
		} else if right.IsTypeOf() && left.IsStringLit() {
			typeOfNode = right
			stringNode = left
		}
		if (typeOfNode != 0) && (stringNode != 0) {
			var operandNode pkg_rhino.Node = typeOfNode.GetFirstChild()
			_ = operandNode
			var operandType pkg_rhino.JSType = this_.GetTypeIfRefinable(operandNode, blindScope)
			_ = operandType
			if operandType != nil && !operandType.IsJavaNil_() {
				var resultEqualsValue bool = (((operatorToken == pkg_rhino.StaticGet_Token_EQ()) || (operatorToken == pkg_rhino.StaticGet_Token_SHEQ())) || (operatorToken == pkg_rhino.StaticGet_Token_CASE()))
				_ = resultEqualsValue
				if !(outcome.IsTruthy()) {
					resultEqualsValue = !(resultEqualsValue)
				}
				return this_.priv_SemanticReverseAbstractInterpreter_caseTypeOf_7861(operandNode, operandType, stringNode.GetString(), resultEqualsValue, blindScope)
			}
		}
	default:
	}
	switch (operatorToken).Ordinal() {
	case int32(69) /* Token.AND */ :
		if outcome.IsTruthy() {
			return this_.priv_SemanticReverseAbstractInterpreter_caseAndOrNotShortCircuiting_7856(condition.GetFirstChild(), condition.GetLastChild(), blindScope, pkg_rhino.StaticGet_Outcome_TRUE())
		} else {
			return this_.priv_SemanticReverseAbstractInterpreter_caseAndOrMaybeShortCircuiting_7857(condition.GetFirstChild(), condition.GetLastChild(), blindScope, pkg_rhino.StaticGet_Outcome_TRUE())
		}
	case int32(68) /* Token.OR */ :
		if !(outcome.IsTruthy()) {
			return this_.priv_SemanticReverseAbstractInterpreter_caseAndOrNotShortCircuiting_7856(condition.GetFirstChild(), condition.GetLastChild(), blindScope, pkg_rhino.StaticGet_Outcome_FALSE())
		} else {
			return this_.priv_SemanticReverseAbstractInterpreter_caseAndOrMaybeShortCircuiting_7857(condition.GetFirstChild(), condition.GetLastChild(), blindScope, pkg_rhino.StaticGet_Outcome_FALSE())
		}
	case int32(4) /* Token.EQ */ :
		if outcome.IsTruthy() {
			return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7854(condition, blindScope, lambdaSingleton_semantic_reverse_abstract_interpreter_1)
		} else {
			return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7854(condition, blindScope, lambdaSingleton_semantic_reverse_abstract_interpreter_2)
		}
	case int32(5) /* Token.NE */ :
		if outcome.IsTruthy() {
			return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7854(condition, blindScope, lambdaSingleton_semantic_reverse_abstract_interpreter_3)
		} else {
			return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7854(condition, blindScope, lambdaSingleton_semantic_reverse_abstract_interpreter_4)
		}
	case int32(40) /* Token.SHEQ */ :
		if outcome.IsTruthy() {
			return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7854(condition, blindScope, lambdaSingleton_semantic_reverse_abstract_interpreter_5)
		} else {
			return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7854(condition, blindScope, lambdaSingleton_semantic_reverse_abstract_interpreter_6)
		}
	case int32(41) /* Token.SHNE */ :
		if outcome.IsTruthy() {
			return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7854(condition, blindScope, lambdaSingleton_semantic_reverse_abstract_interpreter_7)
		} else {
			return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7854(condition, blindScope, lambdaSingleton_semantic_reverse_abstract_interpreter_8)
		}
	case int32(32) /* Token.NAME */, int32(26) /* Token.GETPROP */ :
		return this_.priv_SemanticReverseAbstractInterpreter_caseNameOrGetProp_7860(condition, blindScope, outcome)
	case int32(51) /* Token.ASSIGN */ :
		return this_.FirstPreciserScopeKnowingConditionOutcome(condition.GetFirstChild(), this_.FirstPreciserScopeKnowingConditionOutcome(condition.GetSecondChild(), blindScope, outcome), outcome)
	case int32(19) /* Token.NOT */ :
		return this_.FirstPreciserScopeKnowingConditionOutcome(condition.GetFirstChild(), blindScope, outcome.Not())
	case int32(7) /* Token.LE */, int32(6) /* Token.LT */, int32(9) /* Token.GE */, int32(8) /* Token.GT */ :
		if outcome.IsTruthy() {
			return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7854(condition, blindScope, lambdaSingleton_semantic_reverse_abstract_interpreter_9)
		}
	case int32(45) /* Token.INSTANCEOF */ :
		return this_.priv_SemanticReverseAbstractInterpreter_caseInstanceOf_7862(condition.GetFirstChild(), condition.GetLastChild(), blindScope, outcome)
	case int32(44) /* Token.IN */ :
		if outcome.IsTruthy() && condition.GetFirstChild().IsStringLit() {
			return this_.priv_SemanticReverseAbstractInterpreter_caseIn_7863(condition.GetLastChild(), condition.GetFirstChild().GetString(), blindScope)
		}
	case int32(76) /* Token.CASE */ :
		var left_2 pkg_rhino.Node = condition.GetParent().GetPrevious()
		_ = left_2
		var right_2 pkg_rhino.Node = condition.GetFirstChild()
		_ = right_2
		if outcome.IsTruthy() {
			return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7855(left_2, right_2, blindScope, lambdaSingleton_semantic_reverse_abstract_interpreter_10)
		} else {
			return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7855(left_2, right_2, blindScope, lambdaSingleton_semantic_reverse_abstract_interpreter_11)
		}
	case int32(28) /* Token.CALL */ :
		var left_3 pkg_rhino.Node = condition.GetFirstChild()
		_ = left_3
		var leftName string = left_3.GetQualifiedName()
		_ = leftName
		if ("Array.isArray" == leftName) && (left_3.GetNext() != 0) {
			return this_.priv_SemanticReverseAbstractInterpreter_caseIsArray_7853(left_3.GetNext(), blindScope, outcome)
		}
	default:
	}
	return this_.NextPreciserScopeKnowingConditionOutcome(condition, blindScope, outcome)
}

// Java method: SemanticReverseAbstractInterpreter.caseIsArray(com.google.javascript.rhino.Node value, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.rhino.Outcome outcome)
func (this_ *SemanticReverseAbstractInterpreter) priv_SemanticReverseAbstractInterpreter_caseIsArray_7853(value pkg_rhino.Node, blindScope FlowScope, outcome *pkg_rhino.Outcome) FlowScope {
	_ = value
	_ = blindScope
	_ = outcome
	var type_ pkg_rhino.JSType = this_.GetTypeIfRefinable(value, blindScope)
	_ = type_
	if type_ != nil && !type_.IsJavaNil_() {
		var visitor pkg_rhino.Visitor
		_ = visitor
		if outcome.IsTruthy() {
			visitor = this_.GetChainableReverseAbstractInterpreterBase().restrictToArrayVisitor
		} else {
			visitor = this_.GetChainableReverseAbstractInterpreterBase().restrictToNotArrayVisitor
		}
		return this_.priv_SemanticReverseAbstractInterpreter_maybeRestrictName_7858(blindScope, value, type_, j2gort.AsIface[pkg_rhino.JSType](type_.Visit_1(visitor)))
	}
	return blindScope
}

// Java method: SemanticReverseAbstractInterpreter.caseEquality(com.google.javascript.rhino.Node condition, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter.TypePairOp merging)
func (this_ *SemanticReverseAbstractInterpreter) priv_SemanticReverseAbstractInterpreter_caseEquality_7854(condition pkg_rhino.Node, blindScope FlowScope, merging SemanticReverseAbstractInterpreter_TypePairOp) FlowScope {
	_ = condition
	_ = blindScope
	_ = merging
	return this_.priv_SemanticReverseAbstractInterpreter_caseEquality_7855(condition.GetFirstChild(), condition.GetLastChild(), blindScope, merging)
}

// Java method: SemanticReverseAbstractInterpreter.caseEquality(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter.TypePairOp merging)
func (this_ *SemanticReverseAbstractInterpreter) priv_SemanticReverseAbstractInterpreter_caseEquality_7855(left pkg_rhino.Node, right pkg_rhino.Node, blindScope FlowScope, merging SemanticReverseAbstractInterpreter_TypePairOp) FlowScope {
	_ = left
	_ = right
	_ = blindScope
	_ = merging
	var leftType pkg_rhino.JSType = this_.GetTypeIfRefinable(left, blindScope)
	_ = leftType
	var leftIsRefineable bool = false
	_ = leftIsRefineable
	if leftType != nil && !leftType.IsJavaNil_() {
		leftIsRefineable = true
	} else {
		leftIsRefineable = false
		leftType = left.GetJSType()
	}
	var rightType pkg_rhino.JSType = this_.GetTypeIfRefinable(right, blindScope)
	_ = rightType
	var rightIsRefineable bool = false
	_ = rightIsRefineable
	if rightType != nil && !rightType.IsJavaNil_() {
		rightIsRefineable = true
	} else {
		rightIsRefineable = false
		rightType = right.GetJSType()
	}
	var merged pkg_rhino.JSType_TypePair = merging.Apply(pkg_rhino.NewJSType_TypePair(leftType, rightType))
	_ = merged
	if merged != (pkg_rhino.JSType_TypePair{}) {
		return this_.priv_SemanticReverseAbstractInterpreter_maybeRestrictTwoNames_7859(blindScope, left, leftType, func() pkg_rhino.JSType {
			if leftIsRefineable {
				return merged.FieldGet_typeA()
			}
			return nil
		}(), right, rightType, func() pkg_rhino.JSType {
			if rightIsRefineable {
				return merged.FieldGet_typeB()
			}
			return nil
		}())
	}
	return blindScope
}

// Java method: SemanticReverseAbstractInterpreter.caseAndOrNotShortCircuiting(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.rhino.Outcome outcome)
func (this_ *SemanticReverseAbstractInterpreter) priv_SemanticReverseAbstractInterpreter_caseAndOrNotShortCircuiting_7856(left pkg_rhino.Node, right pkg_rhino.Node, blindScope FlowScope, outcome *pkg_rhino.Outcome) FlowScope {
	_ = left
	_ = right
	_ = blindScope
	_ = outcome
	var leftType pkg_rhino.JSType = this_.GetTypeIfRefinable(left, blindScope)
	_ = leftType
	var leftIsRefineable bool = false
	_ = leftIsRefineable
	if leftType != nil && !leftType.IsJavaNil_() {
		leftIsRefineable = true
	} else {
		leftIsRefineable = false
		leftType = left.GetJSType()
		blindScope = this_.FirstPreciserScopeKnowingConditionOutcome(left, blindScope, outcome)
	}
	var restrictedLeftType pkg_rhino.JSType
	_ = restrictedLeftType
	if leftType == nil || leftType.IsJavaNil_() {
		restrictedLeftType = nil
	} else {
		restrictedLeftType = leftType.GetRestrictedTypeGivenOutcome(outcome)
	}
	if restrictedLeftType == nil || restrictedLeftType.IsJavaNil_() {
		return this_.FirstPreciserScopeKnowingConditionOutcome(right, blindScope, outcome)
	}
	blindScope = this_.priv_SemanticReverseAbstractInterpreter_maybeRestrictName_7858(blindScope, left, leftType, func() pkg_rhino.JSType {
		if leftIsRefineable {
			return restrictedLeftType
		}
		return nil
	}())
	var rightType pkg_rhino.JSType = this_.GetTypeIfRefinable(right, blindScope)
	_ = rightType
	var rightIsRefineable bool = false
	_ = rightIsRefineable
	if rightType != nil && !rightType.IsJavaNil_() {
		rightIsRefineable = true
	} else {
		rightIsRefineable = false
		rightType = right.GetJSType()
		blindScope = this_.FirstPreciserScopeKnowingConditionOutcome(right, blindScope, outcome)
	}
	if outcome.IsTruthy() {
		var restrictedRightType pkg_rhino.JSType
		_ = restrictedRightType
		if rightType == nil || rightType.IsJavaNil_() {
			restrictedRightType = nil
		} else {
			restrictedRightType = rightType.GetRestrictedTypeGivenOutcome(outcome)
		}
		return this_.priv_SemanticReverseAbstractInterpreter_maybeRestrictName_7858(blindScope, right, rightType, func() pkg_rhino.JSType {
			if rightIsRefineable {
				return restrictedRightType
			}
			return nil
		}())
	}
	return blindScope
}

// Java method: SemanticReverseAbstractInterpreter.caseAndOrMaybeShortCircuiting(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.rhino.Outcome outcome)
func (this_ *SemanticReverseAbstractInterpreter) priv_SemanticReverseAbstractInterpreter_caseAndOrMaybeShortCircuiting_7857(left pkg_rhino.Node, right pkg_rhino.Node, blindScope FlowScope, outcome *pkg_rhino.Outcome) FlowScope {
	_ = left
	_ = right
	_ = blindScope
	_ = outcome
	var refinements j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = refinements
	blindScope = NewSemanticReverseAbstractInterpreter_RefinementTrackingFlowScope(blindScope, refinements)
	var leftScope FlowScope = this_.FirstPreciserScopeKnowingConditionOutcome(left, blindScope, outcome.Not())
	_ = leftScope
	var leftVar pkg_rhino.StaticTypedSlot
	_ = leftVar
	if refinements.Size() == int32(1) {
		leftVar = j2gort.AsIface[pkg_rhino.StaticTypedSlot](leftScope.GetSlot(j2gort.CastString(j2gort.AsIface[j2gort.Ext_java_util_Iterator](refinements.Iterator()).Next())))
	} else {
		leftVar = nil
	}
	if leftVar == nil || leftVar.IsJavaNil_() {
		return SemanticReverseAbstractInterpreter_Unwrap(blindScope)
	}
	refinements.Clear()
	var rightScope FlowScope = NewSemanticReverseAbstractInterpreter_RefinementTrackingFlowScope(this_.FirstPreciserScopeKnowingConditionOutcome(left, blindScope, outcome), refinements)
	_ = rightScope
	rightScope = this_.FirstPreciserScopeKnowingConditionOutcome(right, rightScope, outcome.Not())
	var rightVar pkg_rhino.StaticTypedSlot
	_ = rightVar
	if refinements.Size() == int32(1) {
		rightVar = j2gort.AsIface[pkg_rhino.StaticTypedSlot](rightScope.GetSlot(j2gort.CastString(j2gort.AsIface[j2gort.Ext_java_util_Iterator](refinements.Iterator()).Next())))
	} else {
		rightVar = nil
	}
	if (rightVar == nil || rightVar.IsJavaNil_()) || !(func() bool {
		sr_1 := leftVar.GetName()
		se_2 := rightVar.GetName()
		return j2gort.CheckNotNullString(sr_1) == se_2
	}()) {
		return SemanticReverseAbstractInterpreter_Unwrap(blindScope)
	}
	var type_ pkg_rhino.JSType = leftVar.GetType().GetLeastSupertype(rightVar.GetType())
	_ = type_
	return j2gort.AsIface[FlowScope](SemanticReverseAbstractInterpreter_Unwrap(blindScope).InferSlotType(leftVar.GetName(), type_))
}

// Java method: SemanticReverseAbstractInterpreter.maybeRestrictName(com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.rhino.Node node, com.google.javascript.rhino.jstype.JSType originalType, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType restrictedType)
func (this_ *SemanticReverseAbstractInterpreter) priv_SemanticReverseAbstractInterpreter_maybeRestrictName_7858(blindScope FlowScope, node pkg_rhino.Node, originalType pkg_rhino.JSType, restrictedType pkg_rhino.JSType) FlowScope {
	_ = blindScope
	_ = node
	_ = originalType
	_ = restrictedType
	if (restrictedType != nil && !restrictedType.IsJavaNil_()) && !(restrictedType == originalType || ((restrictedType == nil || restrictedType.IsJavaNil_()) && (originalType == nil || originalType.IsJavaNil_()))) {
		return this_.DeclareNameInScope(blindScope, node, restrictedType)
	}
	return blindScope
}

// Java method: SemanticReverseAbstractInterpreter.maybeRestrictTwoNames(com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.rhino.Node left, com.google.javascript.rhino.jstype.JSType originalLeftType, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType restrictedLeftType, com.google.javascript.rhino.Node right, com.google.javascript.rhino.jstype.JSType originalRightType, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType restrictedRightType)
func (this_ *SemanticReverseAbstractInterpreter) priv_SemanticReverseAbstractInterpreter_maybeRestrictTwoNames_7859(blindScope FlowScope, left pkg_rhino.Node, originalLeftType pkg_rhino.JSType, restrictedLeftType pkg_rhino.JSType, right pkg_rhino.Node, originalRightType pkg_rhino.JSType, restrictedRightType pkg_rhino.JSType) FlowScope {
	_ = blindScope
	_ = left
	_ = originalLeftType
	_ = restrictedLeftType
	_ = right
	_ = originalRightType
	_ = restrictedRightType
	var shouldRefineLeft bool = ((restrictedLeftType != nil && !restrictedLeftType.IsJavaNil_()) && !(restrictedLeftType == originalLeftType || ((restrictedLeftType == nil || restrictedLeftType.IsJavaNil_()) && (originalLeftType == nil || originalLeftType.IsJavaNil_()))))
	_ = shouldRefineLeft
	var shouldRefineRight bool = ((restrictedRightType != nil && !restrictedRightType.IsJavaNil_()) && !(restrictedRightType == originalRightType || ((restrictedRightType == nil || restrictedRightType.IsJavaNil_()) && (originalRightType == nil || originalRightType.IsJavaNil_()))))
	_ = shouldRefineRight
	if shouldRefineLeft || shouldRefineRight {
		var informed FlowScope = blindScope
		_ = informed
		if shouldRefineLeft {
			informed = this_.DeclareNameInScope(informed, left, restrictedLeftType)
		}
		if shouldRefineRight {
			informed = this_.DeclareNameInScope(informed, right, restrictedRightType)
		}
		return informed
	}
	return blindScope
}

// Java method: SemanticReverseAbstractInterpreter.caseNameOrGetProp(com.google.javascript.rhino.Node name, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.rhino.Outcome outcome)
func (this_ *SemanticReverseAbstractInterpreter) priv_SemanticReverseAbstractInterpreter_caseNameOrGetProp_7860(name pkg_rhino.Node, blindScope FlowScope, outcome *pkg_rhino.Outcome) FlowScope {
	_ = name
	_ = blindScope
	_ = outcome
	var type_ pkg_rhino.JSType = this_.GetTypeIfRefinable(name, blindScope)
	_ = type_
	if type_ != nil && !type_.IsJavaNil_() {
		var restrictedType pkg_rhino.JSType = type_.GetRestrictedTypeGivenOutcome(outcome)
		_ = restrictedType
		return this_.priv_SemanticReverseAbstractInterpreter_maybeRestrictName_7858(blindScope, name, type_, restrictedType)
	}
	return blindScope
}

// Java method: SemanticReverseAbstractInterpreter.caseTypeOf(com.google.javascript.rhino.Node node, com.google.javascript.rhino.jstype.JSType type, java.lang.String value, boolean resultEqualsValue, com.google.javascript.jscomp.type.FlowScope blindScope)
func (this_ *SemanticReverseAbstractInterpreter) priv_SemanticReverseAbstractInterpreter_caseTypeOf_7861(node pkg_rhino.Node, type_ pkg_rhino.JSType, value string, resultEqualsValue bool, blindScope FlowScope) FlowScope {
	_ = node
	_ = type_
	_ = value
	_ = resultEqualsValue
	_ = blindScope
	return this_.priv_SemanticReverseAbstractInterpreter_maybeRestrictName_7858(blindScope, node, type_, this_.GetRestrictedByTypeOfResult(type_, value, resultEqualsValue))
}

// Java method: SemanticReverseAbstractInterpreter.caseInstanceOf(com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right, com.google.javascript.jscomp.type.FlowScope blindScope, com.google.javascript.rhino.Outcome outcome)
func (this_ *SemanticReverseAbstractInterpreter) priv_SemanticReverseAbstractInterpreter_caseInstanceOf_7862(left pkg_rhino.Node, right pkg_rhino.Node, blindScope FlowScope, outcome *pkg_rhino.Outcome) FlowScope {
	_ = left
	_ = right
	_ = blindScope
	_ = outcome
	var leftType pkg_rhino.JSType = this_.GetTypeIfRefinable(left, blindScope)
	_ = leftType
	if leftType == nil || leftType.IsJavaNil_() {
		return blindScope
	}
	var rightType pkg_rhino.JSType = right.GetJSType()
	_ = rightType
	var targetType pkg_rhino.ObjectType = this_.GetChainableReverseAbstractInterpreterBase().typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	_ = targetType
	if (rightType != nil && !rightType.IsJavaNil_()) && rightType.GetJSTypeBase().IsFunctionType() {
		targetType = pkg_rhino.ToIface_ObjectType(rightType.ToMaybeFunctionType())
	}
	var visitor pkg_rhino.Visitor = nil
	_ = visitor
	if outcome.IsTruthy() {
		visitor = NewSemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor(this_, targetType)
	} else {
		visitor = NewSemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor(this_, targetType)
	}
	return this_.priv_SemanticReverseAbstractInterpreter_maybeRestrictName_7858(blindScope, left, leftType, j2gort.AsIface[pkg_rhino.JSType](leftType.Visit_1(visitor)))
}

// Java method: SemanticReverseAbstractInterpreter.caseIn(com.google.javascript.rhino.Node object, java.lang.String propertyName, com.google.javascript.jscomp.type.FlowScope blindScope)
func (this_ *SemanticReverseAbstractInterpreter) priv_SemanticReverseAbstractInterpreter_caseIn_7863(object pkg_rhino.Node, propertyName string, blindScope FlowScope) FlowScope {
	_ = object
	_ = propertyName
	_ = blindScope
	var jsType pkg_rhino.JSType = object.GetJSType()
	_ = jsType
	if jsType != nil && !jsType.IsJavaNil_() {
		jsType = jsType.RestrictByNotNullOrUndefined()
	} else {
		jsType = nil
	}
	var hasProperty bool = false
	_ = hasProperty
	var objectType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(jsType)
	_ = objectType
	if objectType != nil && !objectType.IsJavaNil_() {
		hasProperty = objectType.GetObjectTypeBase().HasProperty_1_string(propertyName)
	}
	if !(hasProperty) {
		var qualifiedName string = object.GetQualifiedName()
		_ = qualifiedName
		if qualifiedName != j2gort.NullString {
			var propertyQualifiedName string = (j2gort.StrValString(qualifiedName) + "." + j2gort.StrValString(propertyName))
			_ = propertyQualifiedName
			if j2gort.IfaceEquals[pkg_rhino.StaticTypedSlot](j2gort.AsIface[pkg_rhino.StaticTypedSlot](blindScope.GetSlot(propertyQualifiedName)), nil) {
				var unknownType pkg_rhino.JSType = this_.GetChainableReverseAbstractInterpreterBase().typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
				_ = unknownType
				return j2gort.AsIface[FlowScope](blindScope.InferQualifiedSlot(object, propertyQualifiedName, unknownType, unknownType, false))
			}
		}
	}
	return blindScope
}

// Java static method: SemanticReverseAbstractInterpreter.unwrap(com.google.javascript.jscomp.type.FlowScope scope)
func SemanticReverseAbstractInterpreter_Unwrap(scope FlowScope) FlowScope {
	_ = scope
	for j2gort.InstanceOf[*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope](scope) {
		scope = (j2gort.CastPtr[SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope](scope)).delegate
	}
	return scope
}

// Java interface: com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter.TypePairOp
type SemanticReverseAbstractInterpreter_TypePairOp interface {
	Implements_jscomptype_SemanticReverseAbstractInterpreter_TypePairOp()
	IsJavaNil_() bool
	Apply(pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair
}

func ToIface_SemanticReverseAbstractInterpreter_TypePairOp[T interface {
	SemanticReverseAbstractInterpreter_TypePairOp
	comparable
}](v T) SemanticReverseAbstractInterpreter_TypePairOp {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SemanticReverseAbstractInterpreter_TypePairOp[T SemanticReverseAbstractInterpreter_TypePairOp](v T) SemanticReverseAbstractInterpreter_TypePairOp {
	var vi SemanticReverseAbstractInterpreter_TypePairOp = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter.TypePairOp", "com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter$TypePairOp", func(v_ any) bool {
		return j2gort.InstanceOf[SemanticReverseAbstractInterpreter_TypePairOp](v_)
	})
}

type lambda_SemanticReverseAbstractInterpreter_TypePairOp struct {
	fn func(pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair
}

func (l *lambda_SemanticReverseAbstractInterpreter_TypePairOp) IsJavaNil_() bool { return l == nil }
func NewLambda_SemanticReverseAbstractInterpreter_TypePairOp(fn func(pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair) SemanticReverseAbstractInterpreter_TypePairOp {
	if fn == nil {
		return nil
	}
	return &lambda_SemanticReverseAbstractInterpreter_TypePairOp{fn: fn}
}

func (*lambda_SemanticReverseAbstractInterpreter_TypePairOp) Implements_jscomptype_SemanticReverseAbstractInterpreter_TypePairOp() {
}
func (l *lambda_SemanticReverseAbstractInterpreter_TypePairOp) Apply(p0_ pkg_rhino.JSType_TypePair) pkg_rhino.JSType_TypePair {
	return l.fn(p0_)
}

// Java class: com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter.RestrictByTrueInstanceOfResultVisitor
type SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor struct {
	ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase
	outer_ *SemanticReverseAbstractInterpreter
	target pkg_rhino.ObjectType `gson:"target"`
}

func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) BindSelf_(s any) {
	_ = s
	this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase.BindSelf_(s)
}
func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) IsJavaNil_() bool {
	return this_ == nil
}
func (*SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter$RestrictByTrueInstanceOfResultVisitor", "RestrictByTrueInstanceOfResultVisitor"
}
func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase.IdentityHash_()
}
func (*SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter.RestrictByTrueInstanceOfResultVisitor", "com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter$RestrictByTrueInstanceOfResultVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[*SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor](v_)
	})
}
func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) SetOuter_(o *SemanticReverseAbstractInterpreter) {
	this_.outer_ = o
}
func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) GetOuter_() *SemanticReverseAbstractInterpreter {
	return this_.outer_
}
func (*SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor() {
}
func (*SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor() {
}
func (*SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) Implements_rhino_Visitor() {
}

func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) initFields_SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor() {
	_ = this_
}

// Java constructor: RestrictByTrueInstanceOfResultVisitor(com.google.javascript.rhino.jstype.ObjectType target)
func NewSemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor(outer_ *SemanticReverseAbstractInterpreter, target pkg_rhino.ObjectType) *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor {
	this_ := &SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitSemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor(this_, target)
	return this_
}

func InitSemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor(this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor, target pkg_rhino.ObjectType) {
	_ = target
	this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase.outer_ = this_.outer_
	InitChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitor(&this_.ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase)
	this_.initFields_SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor()
	this_.target = target
}

// Java method: RestrictByTrueInstanceOfResultVisitor.caseTopType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) CaseTopType(type_ pkg_rhino.JSType) pkg_rhino.JSType {
	_ = type_
	return this_.priv_RestrictByTrueInstanceOfResultVisitor_applyCommonRestriction_7865(type_)
}

// Java method: RestrictByTrueInstanceOfResultVisitor.caseUnknownType()
func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) CaseUnknownType() any {
	var funcTarget *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(pkg_rhino.Upcast_JSType(this_.target))
	_ = funcTarget
	if (funcTarget != nil) && funcTarget.HasInstanceType() {
		return funcTarget.GetInstanceType()
	}
	return this_.outer_.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
}

// Java method: RestrictByTrueInstanceOfResultVisitor.caseObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) CaseObjectType(type_ pkg_rhino.ObjectType) any {
	_ = type_
	return this_.priv_RestrictByTrueInstanceOfResultVisitor_applyCommonRestriction_7865(pkg_rhino.Upcast_JSType(type_))
}

// Java method: RestrictByTrueInstanceOfResultVisitor.caseUnionType(com.google.javascript.rhino.jstype.UnionType type)
func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) CaseUnionType(type_ *pkg_rhino.UnionType) any {
	_ = type_
	return this_.priv_RestrictByTrueInstanceOfResultVisitor_applyCommonRestriction_7865(pkg_rhino.ToIface_JSType(type_))
}

// Java method: RestrictByTrueInstanceOfResultVisitor.caseFunctionType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) CaseFunctionType(type_ *pkg_rhino.FunctionType) any {
	_ = type_
	return j2gort.AsIface[pkg_rhino.JSType](this_.CaseObjectType(pkg_rhino.ToIface_ObjectType(type_)))
}

// Java method: RestrictByTrueInstanceOfResultVisitor.applyCommonRestriction(com.google.javascript.rhino.jstype.JSType type)
func (this_ *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) priv_RestrictByTrueInstanceOfResultVisitor_applyCommonRestriction_7865(type_ pkg_rhino.JSType) pkg_rhino.JSType {
	_ = type_
	if this_.target.GetObjectTypeBase().IsUnknownType() {
		return type_
	}
	var funcTarget *pkg_rhino.FunctionType = this_.target.ToMaybeFunctionType()
	_ = funcTarget
	if funcTarget.HasInstanceType() {
		return type_.GetJSTypeBase().GetGreatestSubtype(pkg_rhino.Upcast_JSType(funcTarget.GetInstanceType()))
	}
	return nil
}

// Java class: com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter.RestrictByFalseInstanceOfResultVisitor
type SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor struct {
	ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase
	outer_ *SemanticReverseAbstractInterpreter
	target pkg_rhino.ObjectType `gson:"target"`
}

func (this_ *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) BindSelf_(s any) {
	_ = s
	this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase.BindSelf_(s)
}
func (this_ *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) IsJavaNil_() bool {
	return this_ == nil
}
func (*SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter$RestrictByFalseInstanceOfResultVisitor", "RestrictByFalseInstanceOfResultVisitor"
}
func (this_ *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase.IdentityHash_()
}
func (*SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter.RestrictByFalseInstanceOfResultVisitor", "com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter$RestrictByFalseInstanceOfResultVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[*SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor](v_)
	})
}
func (this_ *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) SetOuter_(o *SemanticReverseAbstractInterpreter) {
	this_.outer_ = o
}
func (this_ *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) GetOuter_() *SemanticReverseAbstractInterpreter {
	return this_.outer_
}
func (*SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor() {
}
func (*SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) Implements_jscomptype_ChainableReverseAbstractInterpreter_RestrictByTypeOfResultVisitor() {
}
func (*SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) Implements_rhino_Visitor() {
}

func (this_ *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) initFields_SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor() {
	_ = this_
}

// Java constructor: RestrictByFalseInstanceOfResultVisitor(com.google.javascript.rhino.jstype.ObjectType target)
func NewSemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor(outer_ *SemanticReverseAbstractInterpreter, target pkg_rhino.ObjectType) *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor {
	this_ := &SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitSemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor(this_, target)
	return this_
}

func InitSemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor(this_ *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor, target pkg_rhino.ObjectType) {
	_ = target
	this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase.outer_ = this_.outer_
	InitChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitor(&this_.ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase)
	this_.initFields_SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor()
	this_.target = target
}

// Java method: RestrictByFalseInstanceOfResultVisitor.caseObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) CaseObjectType(type_ pkg_rhino.ObjectType) any {
	_ = type_
	if this_.target.GetObjectTypeBase().IsUnknownType() {
		return type_
	}
	var funcTarget *pkg_rhino.FunctionType = this_.target.ToMaybeFunctionType()
	_ = funcTarget
	if funcTarget.HasInstanceType() {
		if type_.GetObjectTypeBase().IsSubtypeOf_1(pkg_rhino.Upcast_JSType(funcTarget.GetInstanceType())) {
			return nil
		}
		return type_
	}
	return nil
}

// Java method: RestrictByFalseInstanceOfResultVisitor.caseUnionType(com.google.javascript.rhino.jstype.UnionType type)
func (this_ *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) CaseUnionType(type_ *pkg_rhino.UnionType) any {
	_ = type_
	if this_.target.GetObjectTypeBase().IsUnknownType() {
		return j2gort.PtrToAny(type_)
	}
	var funcTarget *pkg_rhino.FunctionType = this_.target.ToMaybeFunctionType()
	_ = funcTarget
	if funcTarget.HasInstanceType() {
		return type_.GetRestrictedUnion(pkg_rhino.Upcast_JSType(funcTarget.GetInstanceType()))
	}
	return nil
}

// Java method: RestrictByFalseInstanceOfResultVisitor.caseFunctionType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) CaseFunctionType(type_ *pkg_rhino.FunctionType) any {
	_ = type_
	return j2gort.AsIface[pkg_rhino.JSType](this_.CaseObjectType(pkg_rhino.ToIface_ObjectType(type_)))
}

// Java class: com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter.RefinementTrackingFlowScope
type SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope struct {
	delegate    FlowScope                `gson:"delegate"`
	refinements j2gort.Ext_java_util_Set `gson:"refinements"`
	idHash_     atomic.Int32
}

func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) BindSelf_(s any) {
	_ = s
}
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) IsJavaNil_() bool {
	return this_ == nil
}
func (*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter$RefinementTrackingFlowScope", "RefinementTrackingFlowScope"
}
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter.RefinementTrackingFlowScope", "com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter$RefinementTrackingFlowScope", func(v_ any) bool {
		return j2gort.InstanceOf[*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope](v_)
	})
}
func (*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) Implements_jscomptype_FlowScope() {
}
func (*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) Implements_rhino_StaticTypedScope() {
}
func (*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) Implements_graph_LatticeElement() {
}
func (*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) Implements_rhino_StaticScope() {
}
func (*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) Implements_graph_Annotation() {
}

func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) initFields_SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope() {
	_ = this_
}

// Java constructor: RefinementTrackingFlowScope(com.google.javascript.jscomp.type.FlowScope delegate, java.util.Set<java.lang.String> refinements)
func NewSemanticReverseAbstractInterpreter_RefinementTrackingFlowScope(delegate FlowScope, refinements j2gort.Ext_java_util_Set) *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope {
	this_ := &SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope{}
	InitSemanticReverseAbstractInterpreter_RefinementTrackingFlowScope(this_, delegate, refinements)
	return this_
}

func InitSemanticReverseAbstractInterpreter_RefinementTrackingFlowScope(this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope, delegate FlowScope, refinements j2gort.Ext_java_util_Set) {
	_ = delegate
	_ = refinements
	this_.initFields_SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope()
	this_.delegate = delegate
	this_.refinements = refinements
}

// Java method: RefinementTrackingFlowScope.withSyntacticScope(com.google.javascript.rhino.jstype.StaticTypedScope scope)
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) WithSyntacticScope(scope pkg_rhino.StaticTypedScope) FlowScope {
	_ = scope
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: RefinementTrackingFlowScope.inferSlotType(java.lang.String symbol, com.google.javascript.rhino.jstype.JSType type)
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) InferSlotType(symbol string, type_ pkg_rhino.JSType) any {
	_ = symbol
	_ = type_
	j2gort.SetAddStr(this_.refinements, symbol)
	return this_.priv_RefinementTrackingFlowScope_wrap_7866(j2gort.AsIface[FlowScope](this_.delegate.InferSlotType(symbol, type_)))
}

// Java method: RefinementTrackingFlowScope.inferQualifiedSlot(com.google.javascript.rhino.Node node, java.lang.String symbol, com.google.javascript.rhino.jstype.JSType bottomType, com.google.javascript.rhino.jstype.JSType inferredType, boolean declare)
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) InferQualifiedSlot(node pkg_rhino.Node, symbol string, bottomType pkg_rhino.JSType, inferredType pkg_rhino.JSType, declare bool) any {
	_ = node
	_ = symbol
	_ = bottomType
	_ = inferredType
	_ = declare
	j2gort.SetAddStr(this_.refinements, symbol)
	return this_.priv_RefinementTrackingFlowScope_wrap_7866(j2gort.AsIface[FlowScope](this_.delegate.InferQualifiedSlot(node, symbol, bottomType, inferredType, declare)))
}

// Java method: RefinementTrackingFlowScope.wrap(com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) priv_RefinementTrackingFlowScope_wrap_7866(scope FlowScope) FlowScope {
	_ = scope
	if !(scope == this_.delegate || ((scope == nil || scope.IsJavaNil_()) && (this_.delegate == nil || this_.delegate.IsJavaNil_()))) {
		return NewSemanticReverseAbstractInterpreter_RefinementTrackingFlowScope(scope, this_.refinements)
	} else {
		return this_
	}
}

// Java method: RefinementTrackingFlowScope.getDeclarationScope()
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetDeclarationScope() any {
	return j2gort.AsIface[pkg_rhino.StaticTypedScope](this_.delegate.GetDeclarationScope())
}

// Java method: RefinementTrackingFlowScope.getRootNode()
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetRootNode_0() pkg_rhino.Node {
	return this_.delegate.GetRootNode_0()
}

// Java method: RefinementTrackingFlowScope.getParentScope()
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetParentScope() any {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: RefinementTrackingFlowScope.getSlot(java.lang.String name)
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetSlot(name string) any {
	_ = name
	return j2gort.AsIface[pkg_rhino.StaticTypedSlot](this_.delegate.GetSlot(name))
}

// Java method: RefinementTrackingFlowScope.getOwnSlot(java.lang.String name)
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetOwnSlot(name string) any {
	_ = name
	return j2gort.AsIface[pkg_rhino.StaticTypedSlot](this_.delegate.GetOwnSlot(name))
}

// Java method: RefinementTrackingFlowScope.getTypeOfThis()
func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetTypeOfThis() pkg_rhino.JSType {
	return this_.delegate.GetTypeOfThis()
}

func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) LookupQualifiedName_1_QualifiedName(p0_ pkg_rhino.QualifiedName) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedName_1_QualifiedName(this_, p0_)
}

func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) LookupQualifiedName_1_string(p0_ string) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedName_1_string(this_, p0_)
}

func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) LookupQualifiedNameNode(p0_ pkg_rhino.Node) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedNameNode(this_, p0_)
}

func (this_ *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetTopmostScopeOfEventualDeclaration(p0_ string) pkg_rhino.StaticScope {
	return pkg_rhino.IfaceDefault_StaticScope_GetTopmostScopeOfEventualDeclaration(this_, p0_)
}
