// Port of com/google/javascript/jscomp/AstValidator.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: AstValidator.validateExpression(Node)
func (v *AstValidator) ValidateExpression(n rhino.Node) {
	v.ValidateTypeInformation(n)
	v.ValidateRequiredInlinings(n)

	switch n.GetToken() {
	// Childless expressions
	case rhino.Token_NEW_TARGET:
		v.ValidateFeature(parser.FeatureSet_Feature_NEW_TARGET, n)
		v.ValidateProperties(n)
		v.ValidateChildless(n)
	case rhino.Token_IMPORT_META:
		v.ValidateFeature(parser.FeatureSet_Feature_IMPORT_META, n)
		v.ValidateProperties(n)
		v.ValidateChildless(n)
	case rhino.Token_FALSE, rhino.Token_NULL, rhino.Token_THIS, rhino.Token_TRUE:
		v.ValidateProperties(n)
		v.ValidateChildless(n)

	// General unary ops
	case rhino.Token_DELPROP, rhino.Token_POS, rhino.Token_NEG, rhino.Token_NOT, rhino.Token_TYPEOF,
		rhino.Token_VOID, rhino.Token_BITNOT, rhino.Token_CAST:
		v.ValidateUnaryOp(n)
	case rhino.Token_INC, rhino.Token_DEC:
		v.ValidateIncDecOp(n)

	// Assignments
	case rhino.Token_ASSIGN:
		v.ValidateAssignmentExpression(n)
	case rhino.Token_ASSIGN_EXPONENT:
		v.ValidateFeature(parser.FeatureSet_Feature_EXPONENT_OP, n)
		v.ValidateCompoundAssignmentExpression(n)
	case rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD:
		v.ValidateCompoundAssignmentExpression(n)
	case rhino.Token_ASSIGN_COALESCE:
		v.ValidateFeature(parser.FeatureSet_Feature_NULL_COALESCE_OP, n)
		v.ValidateFeature(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT, n)
		v.ValidateCompoundAssignmentExpression(n)
	case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND:
		v.ValidateFeature(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT, n)
		v.ValidateCompoundAssignmentExpression(n)

	case rhino.Token_HOOK:
		v.ValidateTrinaryOp(n)

	// Node types that require special handling
	case rhino.Token_STRINGLIT:
		v.ValidateStringLit(n)
	case rhino.Token_NUMBER:
		v.ValidateNumber(n)
	case rhino.Token_BIGINT:
		v.ValidateBigInt(n)
	case rhino.Token_NAME:
		v.ValidateName(n)

	// General binary ops
	case rhino.Token_EXPONENT:
		v.ValidateFeature(parser.FeatureSet_Feature_EXPONENT_OP, n)
		v.ValidateBinaryOp(n)
	case rhino.Token_COALESCE:
		v.ValidateFeature(parser.FeatureSet_Feature_NULL_COALESCE_OP, n)
		v.ValidateBinaryOp(n)
	case rhino.Token_COMMA,
		rhino.Token_OR,
		rhino.Token_AND,
		rhino.Token_BITOR,
		rhino.Token_BITXOR,
		rhino.Token_BITAND,
		rhino.Token_EQ,
		rhino.Token_NE,
		rhino.Token_SHEQ,
		rhino.Token_SHNE,
		rhino.Token_LT,
		rhino.Token_GT,
		rhino.Token_LE,
		rhino.Token_GE,
		rhino.Token_INSTANCEOF,
		rhino.Token_IN,
		rhino.Token_LSH,
		rhino.Token_RSH,
		rhino.Token_URSH,
		rhino.Token_SUB,
		rhino.Token_ADD,
		rhino.Token_MUL,
		rhino.Token_MOD,
		rhino.Token_DIV:
		v.ValidateBinaryOp(n)

	case rhino.Token_GETELEM:
		v.ValidateGetElem(n)
	case rhino.Token_OPTCHAIN_GETELEM:
		v.ValidateOptChainGetElem(n)

	case rhino.Token_GETPROP:
		v.ValidateGetProp(n)
	case rhino.Token_OPTCHAIN_GETPROP:
		v.ValidateOptChainGetProp(n)

	case rhino.Token_ARRAYLIT:
		v.ValidateArrayLit(n)
	case rhino.Token_OBJECTLIT:
		v.ValidateObjectLit(n)
	case rhino.Token_REGEXP:
		v.ValidateRegExpLit(n)

	case rhino.Token_CALL:
		v.ValidateCall(n)
	case rhino.Token_OPTCHAIN_CALL:
		v.ValidateOptChainCall(n)
	case rhino.Token_NEW:
		v.ValidateNew(n)

	case rhino.Token_FUNCTION:
		v.ValidateRequiredInlinings(n)
		v.ValidateFunctionExpression(n)
	case rhino.Token_CLASS:
		v.ValidateClass(n)

	case rhino.Token_TEMPLATELIT:
		v.ValidateTemplateLit(n)
	case rhino.Token_TAGGED_TEMPLATELIT:
		v.ValidateTaggedTemplateLit(n)

	case rhino.Token_YIELD:
		v.ValidateYield(n)
	case rhino.Token_AWAIT:
		v.ValidateAwait(n)
	case rhino.Token_DYNAMIC_IMPORT:
		v.ValidateFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT, n)
		v.ValidateUnaryOp(n)

	default:
		v.Violation("Expected expression but was "+n.GetToken().String(), n)
	}
}

// ValidatePseudoExpression validates an expression or expresison-like
// construct.
//
// An expression-like construct (pseudoexpression) is an AST fragment that
// is valid in some, but not all, of the same contexts as true expressions.
// For example, a VANILLA_FOR permits EMPTY as its condition and increment
// expressions, even though EMPTY is not valid as an expression in general.
//
// allowedPseudoexpressions allows the caller to specify which
// pseudoexpressions are valid for their context. If n is a
// pseudoexpression, it will be considered invalid unless its token is in
// this set.
//
// Java: AstValidator.validatePseudoExpression(Node, Token...)
func (v *AstValidator) ValidatePseudoExpression(n rhino.Node, allowedPseudoexpressions ...rhino.Token) {
	switch n.GetToken() {
	case rhino.Token_EMPTY:
		jutil.CheckArgumentMsg1(slices.Contains(allowedPseudoexpressions, rhino.Token_EMPTY), "Unexpected pseudoexpression %s", n)
		v.ValidateProperties(n)
		v.ValidateChildless(n)
	case rhino.Token_ITER_SPREAD:
		jutil.CheckArgumentMsg1(
			slices.Contains(allowedPseudoexpressions, rhino.Token_ITER_SPREAD), "Unexpected pseudoexpression %s", n)
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		v.ValidateFeature(parser.FeatureSet_Feature_SPREAD_EXPRESSIONS, n)
		v.ValidateExpression(n.GetFirstChild())
	// The only kinds of potential pseudo-expressions we recognize are EMPTY
	// and ITER_SPREAD. So if the given node is neither, validate that it's a
	// (non-pseudo) legitimate expression
	default:
		v.ValidateExpression(n)
	}
}

// Java: AstValidator.validateYield(Node)
func (v *AstValidator) ValidateYield(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_GENERATORS, n)
	v.ValidateNodeType(rhino.Token_YIELD, n)
	v.ValidateProperties(n)
	v.ValidateChildCountIn(n, 0, 1)
	if n.HasChildren() {
		v.ValidateExpression(n.GetFirstChild())
	}
	v.ValidateYieldWithinGeneratorFunction(n)
}

// Java: AstValidator.validateYieldWithinGeneratorFunction(Node)
func (v *AstValidator) ValidateYieldWithinGeneratorFunction(n rhino.Node) {
	parentFunction := NodeUtil_GetEnclosingFunction(n)
	if parentFunction.IsNil() || !parentFunction.IsGeneratorFunction() {
		v.Violation("'yield' expression is not within a generator function", n)
	} else if v.IsInParameterListOfFunction(n, parentFunction) {
		v.Violation("'yield' expression is not allowed in a parameter list", n)
	}
}

// Java: AstValidator.validateAwait(Node)
func (v *AstValidator) ValidateAwait(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS, n)
	v.ValidateNodeType(rhino.Token_AWAIT, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateAwaitWithinAsyncFunction(n)
}

// Java: AstValidator.AWAIT_NOT_WITHIN_ASYNC_FUNCTION
const astValidator_AWAIT_NOT_WITHIN_ASYNC_FUNCTION = "'await' expression is not within an async function"

// Java: AstValidator.AWAIT_NOT_ALLOWED_IN_PARAMETER_LIST
const astValidator_AWAIT_NOT_ALLOWED_IN_PARAMETER_LIST = "'await' expression is not allowed in a parameter list"

// Java: AstValidator.validateAwaitWithinAsyncFunction(Node)
func (v *AstValidator) ValidateAwaitWithinAsyncFunction(n rhino.Node) {
	parentFunction := NodeUtil_GetEnclosingFunction(n)

	if parentFunction.IsNil() {
		// Top-level await is only allowed in modules.
		if !NodeUtil_GetEnclosingScript(n).GetBooleanProp(rhino.Node_ES6_MODULE) {
			v.Violation(astValidator_AWAIT_NOT_WITHIN_ASYNC_FUNCTION, n)
		}
		return
	}

	if !parentFunction.IsAsyncFunction() {
		v.Violation(astValidator_AWAIT_NOT_WITHIN_ASYNC_FUNCTION, n)
	} else if v.IsInParameterListOfFunction(n, parentFunction) {
		v.Violation(astValidator_AWAIT_NOT_ALLOWED_IN_PARAMETER_LIST, n)
	}
}

// Java: AstValidator.isInParameterListOfFunction(Node, Node)
func (v *AstValidator) IsInParameterListOfFunction(child rhino.Node, functionNode rhino.Node) bool {
	paramList := rhino.CheckNotNullMsg1(functionNode.GetSecondChild(), "%s", functionNode)
	for parent := child.GetParent(); parent != functionNode; parent = parent.GetParent() {
		rhino.CheckNotNullMsg2(parent, "%s not contained in function %s", child, functionNode)
		if parent == paramList {
			return true
		}
	}
	return false
}

// Java: AstValidator.validateTaggedTemplateLit(Node)
func (v *AstValidator) ValidateTaggedTemplateLit(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_TEMPLATE_LITERALS, n)
	v.ValidateNodeType(rhino.Token_TAGGED_TEMPLATELIT, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateTemplateLit(n.GetLastChild())
}

// Java: AstValidator.validateTemplateLit(Node)
func (v *AstValidator) ValidateTemplateLit(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_TEMPLATE_LITERALS, n)
	v.ValidateNodeType(rhino.Token_TEMPLATELIT, n)
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if child.IsTemplateLitString() {
			v.ValidateTemplateLitString(child)
		} else {
			v.ValidateTemplateLitSub(child)
		}
	}
}

// Java: AstValidator.validateTemplateLitString(Node)
func (v *AstValidator) ValidateTemplateLitString(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_TEMPLATELIT_STRING, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	_ = n.GetRawString()
}

// Java: AstValidator.validateTemplateLitSub(Node)
func (v *AstValidator) ValidateTemplateLitSub(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_TEMPLATELIT_SUB, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
}

// Java: AstValidator.validateCall(Node)
func (v *AstValidator) ValidateCall(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_CALL, n)
	v.ValidateProperties(n)
	v.ValidateMinimumChildCount(n, 1)
	callee := n.GetFirstChild()
	if callee.IsSuper() {
		v.ValidateSuper(callee)
	} else {
		v.ValidateExpression(callee)
	}
	for c := callee.GetNext(); !c.IsNil(); c = c.GetNext() {
		v.ValidatePseudoExpression(c, rhino.Token_ITER_SPREAD)
	}
}

// Java: AstValidator.validateOptChainCall(Node)
func (v *AstValidator) ValidateOptChainCall(node rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING, node)
	v.ValidateNodeType(rhino.Token_OPTCHAIN_CALL, node)
	v.ValidateProperties(node)
	v.ValidateMinimumChildCount(node, 1)
	callee := node.GetFirstChild()
	v.ValidateExpression(callee)
	for argument := callee.GetNext(); !argument.IsNil(); argument = argument.GetNext() {
		v.ValidatePseudoExpression(argument, rhino.Token_ITER_SPREAD)
	}
	v.ValidateFirstNodeOfOptChain(node)
}

// Java: AstValidator.validateSuper(Node)
func (v *AstValidator) ValidateSuper(superNode rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_SUPER, superNode)
	v.ValidateProperties(superNode)
	v.ValidateChildless(superNode)
	v.ValidateTypeInformation(superNode)
	superParent := superNode.GetParent()
	methodNode := NodeUtil_GetEnclosingNonArrowFunction(superParent)

	if NodeUtil_IsNormalGet(superParent) && superNode.IsFirstChildOf(superParent) {
		// `super.prop` or `super['prop']`
		if !v.AllowsSuperPropertyReference(superParent) {
			v.Violation(
				"super property references are only allowed in methods, class static blocks and class"+
					" fields.",
				superNode)
		}
	} else if superParent.IsCall() && superNode.IsFirstChildOf(superParent) {
		// super() constructor call
		if methodNode.IsNil() || !NodeUtil_IsEs6Constructor(methodNode) {
			v.Violation("super constructor call is only allowed in a constructor method", superNode)
		} else {
			extendsNode := methodNode.
				GetParent().     // MEMBER_FUNCTION_DEF
				GetParent().     // CLASS_METHODS
				GetParent().     // CLASS
				GetSecondChild() // extends clause
			if extendsNode.IsEmpty() {
				v.Violation("super constructor call in a class that extends nothing", superNode)
			}
		}
	} else {
		v.Violation("`super` is a syntax error here", superNode)
	}
}

// AllowsSuperPropertyReference checks if a super property reference is in
// a method, class static block or class field.
//
// Java: AstValidator.allowsSuperPropertyReference(Node)
func (v *AstValidator) AllowsSuperPropertyReference(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		return false
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		return true
	case rhino.Token_FUNCTION:
		if NodeUtil_IsMethodDeclaration(n) {
			return true
		} else if !n.IsArrowFunction() {
			return false
		}
		return v.AllowsSuperPropertyReference(n.GetParent())
	case rhino.Token_BLOCK:
		if NodeUtil_IsClassStaticBlock(n) {
			return true
		}
		return v.AllowsSuperPropertyReference(n.GetParent())
	default:
		return v.AllowsSuperPropertyReference(n.GetParent())
	}
}

// Java: AstValidator.validateObjectSpread(Node)
func (v *AstValidator) ValidateObjectSpread(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateFeature(parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD, n)
	v.ValidateExpression(n.GetFirstChild())
}

// Java: AstValidator.validateNew(Node)
func (v *AstValidator) ValidateNew(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_NEW, n)
	v.ValidateProperties(n)
	v.ValidateMinimumChildCount(n, 1)

	v.ValidateExpression(n.GetFirstChild())
	for c := n.GetSecondChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidatePseudoExpression(c, rhino.Token_ITER_SPREAD)
	}
}

// Java: AstValidator.validateGetPropGetElemInLHS(Token, Node)
func (v *AstValidator) ValidateGetPropGetElemInLHS(contextType rhino.Token, n rhino.Node) {
	if contextType == rhino.Token_CONST ||
		contextType == rhino.Token_LET ||
		contextType == rhino.Token_VAR ||
		contextType == rhino.Token_PARAM_LIST {
		v.Violation("Invalid child for "+contextType.String()+" node", n)
		return
	}
	switch n.GetToken() {
	case rhino.Token_GETPROP:
		v.ValidateGetProp(n)
	case rhino.Token_GETELEM:
		v.ValidateGetElem(n)
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Expected GETPROP or GETELEM but instead got node " + n.GetToken().String()))
	}
}

// Java: AstValidator.validateAssignmentExpression(Node)
func (v *AstValidator) ValidateAssignmentExpression(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateLHS(n.GetToken(), n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateCompoundAssignmentExpression(Node)
func (v *AstValidator) ValidateCompoundAssignmentExpression(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	contextType := n.GetToken()
	lhs := n.GetFirstChild()
	v.ValidateAssignmentOpTarget(lhs, contextType)
	v.ValidateExpression(n.GetLastChild())
}

// ValidateAssignmentOpTarget validates the lhs of a compound assignment op,
// inc, or dec
//
// This check is stricter than validateLhs.
//
// Java: AstValidator.validateAssignmentOpTarget(Node, Token)
func (v *AstValidator) ValidateAssignmentOpTarget(lhs rhino.Node, contextType rhino.Token) {
	switch lhs.GetToken() {
	case rhino.Token_NAME:
		v.ValidateName(lhs)
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		v.ValidateGetPropGetElemInLHS(contextType, lhs)
	case rhino.Token_CAST:
		v.ValidateProperties(lhs)
		v.ValidateChildCountWithExpected(lhs, 1)
		v.ValidateAssignmentOpTarget(lhs.GetFirstChild(), contextType)
	default:
		v.Violation("Invalid child for "+contextType.String()+" node", lhs)
	}
}

// Java: AstValidator.validateGetElem(Node)
func (v *AstValidator) ValidateGetElem(n rhino.Node) {
	jutil.CheckArgumentMsg1(n.IsGetElem(), "%s", n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 2)
	v.ValidatePropertyReferenceTarget(n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateOptChainGetElem(Node)
func (v *AstValidator) ValidateOptChainGetElem(node rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING, node)
	jutil.CheckArgumentMsg1(node.IsOptChainGetElem(), "%s", node)
	v.ValidateProperties(node)
	v.ValidateChildCountWithExpected(node, 2)
	v.ValidateExpression(node.GetFirstChild())
	v.ValidateExpression(node.GetLastChild())
	v.ValidateFirstNodeOfOptChain(node)
}

// Java: AstValidator.validateGetProp(Node)
func (v *AstValidator) ValidateGetProp(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_GETPROP, n)
	v.ValidatePropertyReferenceTarget(n.GetFirstChild())
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateNonEmptyString(n)
}

// Java: AstValidator.validateOptChainGetProp(Node)
func (v *AstValidator) ValidateOptChainGetProp(node rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING, node)
	v.ValidateNodeType(rhino.Token_OPTCHAIN_GETPROP, node)
	v.ValidateExpression(node.GetFirstChild())
	v.ValidateFirstNodeOfOptChain(node)
	v.ValidateProperties(node)
	v.ValidateChildCount(node)
	v.ValidateNonEmptyString(node)
}

// Java: AstValidator.validatePropertyReferenceTarget(Node)
func (v *AstValidator) ValidatePropertyReferenceTarget(objectNode rhino.Node) {
	if objectNode.IsSuper() {
		v.ValidateSuper(objectNode)
	} else {
		v.ValidateExpression(objectNode)
	}
}

// Java: AstValidator.validateRegExpLit(Node)
func (v *AstValidator) ValidateRegExpLit(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_REGEXP_SYNTAX, n)
	v.ValidateNodeType(rhino.Token_REGEXP, n)
	v.ValidateProperties(n)
	v.ValidateChildCountIn(n, 1, 2)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateStringLit(c)
	}
}

// Java: AstValidator.validateStringLit(Node)
func (v *AstValidator) ValidateStringLit(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_STRINGLIT, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	_ = n.GetString()
}

// Java: AstValidator.validateNumber(Node)
func (v *AstValidator) ValidateNumber(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_NUMBER, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	_ = n.GetDouble()
}

// Java: AstValidator.validateBigInt(Node)
func (v *AstValidator) ValidateBigInt(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_BIGINT, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	_ = n.GetBigInt()
}

// Java: AstValidator.validateArrayLit(Node)
func (v *AstValidator) ValidateArrayLit(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_ARRAYLIT, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		// Array-literals may have empty slots.
		v.ValidatePseudoExpression(c, rhino.Token_EMPTY, rhino.Token_ITER_SPREAD)
		break
	}
}

// Java: AstValidator.validateObjectLit(Node)
func (v *AstValidator) ValidateObjectLit(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_OBJECTLIT, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateObjectLitKey(c)
	}
}

// Java: AstValidator.validateObjectLitKey(Node)
func (v *AstValidator) ValidateObjectLitKey(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETTER_DEF:
		v.ValidateObjectLitGetKey(n)
	case rhino.Token_SETTER_DEF:
		v.ValidateObjectLitSetKey(n)
	case rhino.Token_STRING_KEY:
		v.ValidateObjectLitStringKey(n)
	case rhino.Token_MEMBER_FUNCTION_DEF:
		v.ValidateClassMember(n, false)
		if n.IsStaticMember() {
			v.Violation("Keys in an object literal should not be static.", n)
		}
	case rhino.Token_COMPUTED_PROP:
		v.ValidateObjectLitComputedPropKey(n)
	case rhino.Token_OBJECT_SPREAD:
		v.ValidateObjectSpread(n)
	default:
		v.Violation("Expected object literal key expression but was "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateObjectLitGetKey(Node)
func (v *AstValidator) ValidateObjectLitGetKey(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_GETTER, n)
	v.ValidateNodeType(rhino.Token_GETTER_DEF, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateObjectLiteralKeyName(n)
	function := n.GetFirstChild()
	v.ValidateFunctionExpression(function)
	// objlit get functions must be nameless, and must have zero parameters.
	if function.GetFirstChild().GetString() != "" {
		v.Violation("Expected unnamed function expression.", n)
	}
	functionParams := function.GetSecondChild()
	if functionParams.HasChildren() {
		v.Violation("get methods must not have parameters.", n)
	}
}

// Java: AstValidator.validateObjectLitSetKey(Node)
func (v *AstValidator) ValidateObjectLitSetKey(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_SETTER, n)
	v.ValidateNodeType(rhino.Token_SETTER_DEF, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateObjectLiteralKeyName(n)
	function := n.GetFirstChild()
	v.ValidateFunctionExpression(function)
	// objlit set functions must be nameless, and must have 1 parameter.
	if function.GetFirstChild().GetString() != "" {
		v.Violation("Expected unnamed function expression.", n)
	}
	functionParams := function.GetSecondChild()
	if !functionParams.HasOneChild() {
		v.Violation("set methods must have exactly one parameter.", n)
	}
}

// Java: AstValidator.validateObjectLitStringKey(Node)
func (v *AstValidator) ValidateObjectLitStringKey(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_STRING_KEY, n)
	v.ValidateObjectLiteralKeyName(n)

	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 1)
	v.ValidateExpression(n.GetFirstChild())
	if n.GetBooleanProp(rhino.Node_IS_SHORTHAND_PROPERTY) {
		v.ValidateFeature(parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES, n)
	}
}

// Java: AstValidator.validateObjectLitComputedPropKey(Node)
func (v *AstValidator) ValidateObjectLitComputedPropKey(n rhino.Node) {
	// TODO: b/499356003 - clean up this hack. The right fix is to:
	// 1) model computed prop getter/setters with a separate Feature from
	// COMPUTED_PROPERTIES
	// 2) modify the transpilation passes to *not* remove the computed prop
	// getter/setter in object literals feature from the compiler featureset
	// because it can't actually transpile them. (The transpilation passes
	// report a suppressible error instead).
	if !(n.GetIsInClosureUnawareSubtree() && NodeUtil_IsGetOrSetKey(n)) {
		v.ValidateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	}
	v.ValidateNodeType(rhino.Token_COMPUTED_PROP, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateObjectLitComputedPropGetKey(Node)
func (v *AstValidator) ValidateObjectLitComputedPropGetKey(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.ValidateFeature(parser.FeatureSet_Feature_GETTER, n)
	v.ValidateNodeType(rhino.Token_COMPUTED_PROP, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	function := n.GetLastChild()
	v.ValidateFunctionExpression(function)
	// objlit get functions must be nameless, and must have zero parameters.
	if function.GetFirstChild().GetString() != "" {
		v.Violation("Expected unnamed function expression.", n)
	}
	functionParams := function.GetSecondChild()
	if functionParams.HasChildren() {
		v.Violation("get methods must not have parameters.", n)
	}
}

// Java: AstValidator.validateObjectLitComputedPropSetKey(Node)
func (v *AstValidator) ValidateObjectLitComputedPropSetKey(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.ValidateFeature(parser.FeatureSet_Feature_SETTER, n)
	v.ValidateNodeType(rhino.Token_COMPUTED_PROP, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	function := n.GetLastChild()
	v.ValidateFunctionExpression(function)
	// objlit set functions must be nameless, and must have 1 parameter.
	if function.GetFirstChild().GetString() != "" {
		v.Violation("Expected unnamed function expression.", n)
	}
	functionParams := function.GetSecondChild()
	if !functionParams.HasOneChild() {
		v.Violation("set methods must have exactly one parameter.", n)
	}
}

// Java: AstValidator.validateObjectLiteralKeyName(Node)
func (v *AstValidator) ValidateObjectLiteralKeyName(n rhino.Node) {
	if n.IsQuotedStringKey() {
		_ = n.GetString()
	} else {
		v.ValidateNonEmptyString(n)
	}
}

// Java: AstValidator.validateIncDecOp(Node)
func (v *AstValidator) ValidateIncDecOp(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 1)
	v.ValidateAssignmentOpTarget(n.GetFirstChild(), n.GetToken())
}

// Java: AstValidator.validateUnaryOp(Node)
func (v *AstValidator) ValidateUnaryOp(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 1)
	v.ValidateExpression(n.GetFirstChild())
}

// Java: AstValidator.validateBinaryOp(Node)
func (v *AstValidator) ValidateBinaryOp(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 2)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateTrinaryOp(Node)
func (v *AstValidator) ValidateTrinaryOp(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 3)
	first := n.GetFirstChild()
	v.ValidateExpression(first)
	v.ValidateExpression(first.GetNext())
	v.ValidateExpression(n.GetLastChild())
}
