// Port of com/google/javascript/jscomp/AstValidator.java @ closure-compiler v20261028

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: AstValidator.validateExpression(Node)
func (v *AstValidator) ValidateExpression(n rhino.Node) {
	v.validateTypeInformation(n)
	v.validateRequiredInlinings(n)

	switch n.GetToken() {
	case rhino.Token_NEW_TARGET:
		v.validateFeature(parser.FeatureSet_Feature_NEW_TARGET, n)
		v.validateChildlessNode(n)
	case rhino.Token_IMPORT_META:
		v.validateFeature(parser.FeatureSet_Feature_IMPORT_META, n)
		v.validateChildlessNode(n)
	case rhino.Token_FALSE, rhino.Token_NULL, rhino.Token_THIS, rhino.Token_TRUE:
		v.validateChildlessNode(n)

	case rhino.Token_DELPROP, rhino.Token_POS, rhino.Token_NEG, rhino.Token_NOT, rhino.Token_TYPEOF,
		rhino.Token_VOID, rhino.Token_BITNOT:
		v.validatePropsAndCount(n, 1)
		v.ValidateExpression(ast.MustUnary(n).Operand().Node())
	case rhino.Token_CAST:
		v.validatePropsAndCount(n, 1)
		v.ValidateExpression(ast.MustCast(n).Expr().Node())
	case rhino.Token_DYNAMIC_IMPORT:
		v.validateFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT, n)
		v.validatePropsAndCount(n, 1)
		v.ValidateExpression(ast.MustDynamicImport(n).Specifier().Node())
	case rhino.Token_INC, rhino.Token_DEC:
		v.validateIncDecOp(ast.MustIncDec(n))

	case rhino.Token_ASSIGN:
		v.validateAssignmentExpression(ast.MustAssign(n))
	case rhino.Token_ASSIGN_EXPONENT:
		v.validateFeature(parser.FeatureSet_Feature_EXPONENT_OP, n)
		v.validateCompoundAssignmentExpression(ast.MustAssignOp(n))
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
		v.validateCompoundAssignmentExpression(ast.MustAssignOp(n))
	case rhino.Token_ASSIGN_COALESCE:
		v.validateFeature(parser.FeatureSet_Feature_NULL_COALESCE_OP, n)
		v.validateFeature(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT, n)
		v.validateCompoundAssignmentExpression(ast.MustAssignOp(n))
	case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND:
		v.validateFeature(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT, n)
		v.validateCompoundAssignmentExpression(ast.MustAssignOp(n))

	case rhino.Token_HOOK:
		v.validateTrinaryOp(ast.MustHook(n))

	case rhino.Token_STRINGLIT:
		v.validateStringLit(n)
	case rhino.Token_NUMBER:
		v.validateNumber(ast.MustNumber(n))
	case rhino.Token_BIGINT:
		v.validateBigInt(ast.MustBigInt(n))
	case rhino.Token_NAME:
		v.validateName(n)

	case rhino.Token_EXPONENT:
		v.validateFeature(parser.FeatureSet_Feature_EXPONENT_OP, n)
		v.validateBinaryOp(ast.MustBinary(n))
	case rhino.Token_COALESCE:
		v.validateFeature(parser.FeatureSet_Feature_NULL_COALESCE_OP, n)
		v.validateBinaryOp(ast.MustBinary(n))
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
		v.validateBinaryOp(ast.MustBinary(n))

	case rhino.Token_GETELEM:
		v.validateGetElem(ast.MustGetElem(n))
	case rhino.Token_OPTCHAIN_GETELEM:
		v.validateOptChainGetElem(ast.MustOptChainGetElem(n))
	case rhino.Token_GETPROP:
		v.validateGetProp(ast.MustGetProp(n))
	case rhino.Token_OPTCHAIN_GETPROP:
		v.validateOptChainGetProp(ast.MustOptChainGetProp(n))

	case rhino.Token_ARRAYLIT:
		v.validateArrayLit(ast.MustArrayLit(n))
	case rhino.Token_OBJECTLIT:
		v.validateObjectLit(ast.MustObjectLit(n))
	case rhino.Token_REGEXP:
		v.validateRegExpLit(ast.MustRegExp(n))

	case rhino.Token_CALL:
		v.validateCall(ast.MustCall(n))
	case rhino.Token_OPTCHAIN_CALL:
		v.validateOptChainCall(ast.MustOptChainCall(n))
	case rhino.Token_NEW:
		v.validateNew(ast.MustNew(n))

	case rhino.Token_FUNCTION:
		v.validateRequiredInlinings(n)
		v.validateFunctionExpression(n)
	case rhino.Token_CLASS:
		v.validateClass(ast.MustClass(n))

	case rhino.Token_TEMPLATELIT:
		v.validateTemplateLit(n)
	case rhino.Token_TAGGED_TEMPLATELIT:
		v.validateTaggedTemplateLit(ast.MustTaggedTemplateLit(n))

	case rhino.Token_YIELD:
		v.validateYield(ast.MustYield(n))
	case rhino.Token_AWAIT:
		v.validateAwait(ast.MustAwait(n))

	default:
		v.Violation("Expected expression but was "+n.GetToken().String(), n)
	}
}

// validatePseudoExpression validates an expression or expression-like
// construct (EMPTY or ITER_SPREAD when explicitly allowed in context).
//
// Java: AstValidator.validatePseudoExpression(Node, Token...)
func (v *AstValidator) validatePseudoExpression(n rhino.Node, allowedPseudoexpressions ...rhino.Token) {
	switch n.GetToken() {
	case rhino.Token_EMPTY:
		jutil.CheckArgumentMsg1(
			slices.Contains(allowedPseudoexpressions, rhino.Token_EMPTY), "Unexpected pseudoexpression %s", n)
		v.validateChildlessNode(n)
	case rhino.Token_ITER_SPREAD:
		jutil.CheckArgumentMsg1(
			slices.Contains(allowedPseudoexpressions, rhino.Token_ITER_SPREAD), "Unexpected pseudoexpression %s", n)
		v.validatePropsAndArity(n)
		v.validateFeature(parser.FeatureSet_Feature_SPREAD_EXPRESSIONS, n)
		v.ValidateExpression(ast.MustIterSpread(n).Expr().Node())
	default:
		v.ValidateExpression(n)
	}
}

// Java: AstValidator.validateYield(Node)
func (v *AstValidator) validateYield(y ast.Yield) {
	n := y.Node()
	v.validateFeature(parser.FeatureSet_Feature_GENERATORS, n)
	v.validateProperties(n)
	v.validateChildCountIn(n, 0, 1)
	if val, ok := y.Value(); ok {
		v.ValidateExpression(val.Node())
	}
	v.validateYieldWithinGeneratorFunction(n)
}

// Java: AstValidator.validateYieldWithinGeneratorFunction(Node)
func (v *AstValidator) validateYieldWithinGeneratorFunction(n rhino.Node) {
	fn, ok := ast.AsFunction(NodeUtil_GetEnclosingFunction(n))
	if !ok || !fn.IsGenerator() {
		v.Violation("'yield' expression is not within a generator function", n)
	} else if v.isInParameterListOfFunction(n, fn) {
		v.Violation("'yield' expression is not allowed in a parameter list", n)
	}
}

// Java: AstValidator.validateAwait(Node)
func (v *AstValidator) validateAwait(aw ast.Await) {
	n := aw.Node()
	v.validateFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS, n)
	v.validatePropsAndArity(n)
	v.ValidateExpression(aw.Value().Node())
	v.validateAwaitWithinAsyncFunction(n)
}

const (
	astValidator_AWAIT_NOT_WITHIN_ASYNC_FUNCTION     = "'await' expression is not within an async function"
	astValidator_AWAIT_NOT_ALLOWED_IN_PARAMETER_LIST = "'await' expression is not allowed in a parameter list"
)

// Java: AstValidator.validateAwaitWithinAsyncFunction(Node)
func (v *AstValidator) validateAwaitWithinAsyncFunction(n rhino.Node) {
	fn, ok := ast.AsFunction(NodeUtil_GetEnclosingFunction(n))
	if !ok {
		// Top-level await is only allowed in modules.
		if !NodeUtil_GetEnclosingScript(n).GetBooleanProp(rhino.Node_ES6_MODULE) {
			v.Violation(astValidator_AWAIT_NOT_WITHIN_ASYNC_FUNCTION, n)
		}
		return
	}
	if !fn.IsAsync() {
		v.Violation(astValidator_AWAIT_NOT_WITHIN_ASYNC_FUNCTION, n)
	} else if v.isInParameterListOfFunction(n, fn) {
		v.Violation(astValidator_AWAIT_NOT_ALLOWED_IN_PARAMETER_LIST, n)
	}
}

// Java: AstValidator.isInParameterListOfFunction(Node, Node)
func (v *AstValidator) isInParameterListOfFunction(child rhino.Node, fn ast.Function) bool {
	paramList := rhino.CheckNotNullMsg1(fn.Params().Node(), "%s", fn.Node())
	for parent := child.GetParent(); !fn.Is(parent); parent = parent.GetParent() {
		rhino.CheckNotNullMsg2(parent, "%s not contained in function %s", child, fn.Node())
		if parent == paramList {
			return true
		}
	}
	return false
}

// Java: AstValidator.validateTaggedTemplateLit(Node)
func (v *AstValidator) validateTaggedTemplateLit(tagged ast.TaggedTemplateLit) {
	n := tagged.Node()
	v.validateFeature(parser.FeatureSet_Feature_TEMPLATE_LITERALS, n)
	if v.validatePropsAndArity(n) {
		v.ValidateExpression(tagged.Tag().Node())
		v.validateTemplateLit(tagged.Template().Node())
		return
	}
	v.ValidateExpression(n.GetFirstChild())
	v.validateTemplateLit(n.GetLastChild())
}

// Java: AstValidator.validateTemplateLit(Node)
func (v *AstValidator) validateTemplateLit(n rhino.Node) {
	v.validateFeature(parser.FeatureSet_Feature_TEMPLATE_LITERALS, n)
	if v.validateNodeType(rhino.Token_TEMPLATELIT, n) {
		for part := range ast.MustTemplateLit(n).Parts() {
			if str, ok := part.AsTemplateLitString(); ok {
				v.validateTemplateLitString(str)
			} else {
				v.validateTemplateLitSub(part.Node())
			}
		}
		return
	}
	for child := range n.Children() {
		if str, ok := ast.AsTemplateLitString(child); ok {
			v.validateTemplateLitString(str)
		} else {
			v.validateTemplateLitSub(child)
		}
	}
}

// Java: AstValidator.validateTemplateLitString(Node)
func (v *AstValidator) validateTemplateLitString(str ast.TemplateLitString) {
	n := str.Node()
	v.validatePropsAndArity(n)
	_ = str.Raw()
}

// Java: AstValidator.validateTemplateLitSub(Node)
func (v *AstValidator) validateTemplateLitSub(n rhino.Node) {
	if v.validateNodeType(rhino.Token_TEMPLATELIT_SUB, n) {
		sub := ast.MustTemplateLitSub(n)
		v.validatePropsAndArity(n)
		v.ValidateExpression(sub.Expr().Node())
		return
	}
	v.validatePropsAndArity(n)
	v.ValidateExpression(n.GetFirstChild())
}

// Java: AstValidator.validateCall(Node)
func (v *AstValidator) validateCall(call ast.Call) {
	n := call.Node()
	v.validateProperties(n)
	v.validateMinimumChildCount(n, 1)
	callee := call.Callee().Node()
	if callee.IsSuper() {
		v.validateSuper(callee)
	} else {
		v.ValidateExpression(callee)
	}
	for arg := range call.Args() {
		v.validatePseudoExpression(arg.Node(), rhino.Token_ITER_SPREAD)
	}
}

// Java: AstValidator.validateOptChainCall(Node)
func (v *AstValidator) validateOptChainCall(call ast.OptChainCall) {
	n := call.Node()
	v.validateFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING, n)
	v.validateProperties(n)
	v.validateMinimumChildCount(n, 1)
	v.ValidateExpression(call.Callee().Node())
	for arg := range call.Args() {
		v.validatePseudoExpression(arg.Node(), rhino.Token_ITER_SPREAD)
	}
	v.validateFirstNodeOfOptChain(call.AsOptChain())
}

// Java: AstValidator.validateSuper(Node)
func (v *AstValidator) validateSuper(superNode rhino.Node) {
	v.validateFeature(parser.FeatureSet_Feature_SUPER, superNode)
	v.validateChildlessNode(superNode)
	v.validateTypeInformation(superNode)
	superParent := superNode.GetParent()
	methodNode := NodeUtil_GetEnclosingNonArrowFunction(superParent)

	if get, ok := ast.AsNormalGet(superParent); ok && get.Object().Is(superNode) {
		// `super.prop` or `super['prop']`
		if !v.allowsSuperPropertyReference(superParent) {
			v.Violation(
				"super property references are only allowed in methods, class static blocks and class fields.",
				superNode)
		}
		return
	}

	if call, ok := ast.AsCall(superParent); ok && call.Callee().Is(superNode) {
		// super() constructor call
		if methodNode.IsNil() || !NodeUtil_IsEs6Constructor(methodNode) {
			v.Violation("super constructor call is only allowed in a constructor method", superNode)
			return
		}
		classNode := methodNode.
			GetParent(). // MEMBER_FUNCTION_DEF
			GetParent(). // CLASS_MEMBERS
			GetParent()  // CLASS
		if cls, ok := ast.AsClass(classNode); ok {
			if _, hasExtends := cls.Extends(); !hasExtends {
				v.Violation("super constructor call in a class that extends nothing", superNode)
			}
		} else if classNode.GetSecondChild().IsEmpty() {
			v.Violation("super constructor call in a class that extends nothing", superNode)
		}
		return
	}

	v.Violation("`super` is a syntax error here", superNode)
}

// allowsSuperPropertyReference checks if a super property reference is inside
// a method, class static block, or class field.
//
// Java: AstValidator.allowsSuperPropertyReference(Node)
func (v *AstValidator) allowsSuperPropertyReference(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		return false
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		return true
	case rhino.Token_FUNCTION:
		if NodeUtil_IsMethodDeclaration(n) {
			return true
		}
		if !ast.MustFunction(n).IsArrow() {
			return false
		}
		return v.allowsSuperPropertyReference(n.GetParent())
	case rhino.Token_BLOCK:
		if NodeUtil_IsClassStaticBlock(n) {
			return true
		}
		return v.allowsSuperPropertyReference(n.GetParent())
	default:
		return v.allowsSuperPropertyReference(n.GetParent())
	}
}

// Java: AstValidator.validateObjectSpread(Node)
func (v *AstValidator) validateObjectSpread(spread ast.ObjectSpread) {
	n := spread.Node()
	v.validatePropsAndArity(n)
	v.validateFeature(parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD, n)
	v.ValidateExpression(spread.Expr().Node())
}

// Java: AstValidator.validateNew(Node)
func (v *AstValidator) validateNew(newExpr ast.New) {
	n := newExpr.Node()
	v.validateProperties(n)
	v.validateMinimumChildCount(n, 1)
	v.ValidateExpression(newExpr.Callee().Node())
	for arg := range newExpr.Args() {
		v.validatePseudoExpression(arg.Node(), rhino.Token_ITER_SPREAD)
	}
}

// Java: AstValidator.validateGetPropGetElemInLHS(Token, Node)
func (v *AstValidator) validateGetPropGetElemInLHS(contextType rhino.Token, n rhino.Node) {
	switch contextType {
	case rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR, rhino.Token_PARAM_LIST:
		v.Violation("Invalid child for "+contextType.String()+" node", n)
		return
	}
	switch n.GetToken() {
	case rhino.Token_GETPROP:
		v.validateGetProp(ast.MustGetProp(n))
	case rhino.Token_GETELEM:
		v.validateGetElem(ast.MustGetElem(n))
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Expected GETPROP or GETELEM but instead got node " + n.GetToken().String()))
	}
}

// Java: AstValidator.validateAssignmentExpression(Node)
func (v *AstValidator) validateAssignmentExpression(assign ast.Assign) {
	n := assign.Node()
	if v.validatePropsAndArity(n) {
		v.validateLHS(rhino.Token_ASSIGN, assign.Target().Node())
		v.ValidateExpression(assign.Value().Node())
		return
	}
	v.validateLHS(rhino.Token_ASSIGN, n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateCompoundAssignmentExpression(Node)
func (v *AstValidator) validateCompoundAssignmentExpression(assign ast.AssignOp) {
	n := assign.Node()
	contextType := assign.Token()
	if v.validatePropsAndArity(n) {
		v.validateAssignmentOpTarget(assign.Target().Node(), contextType)
		v.ValidateExpression(assign.Value().Node())
		return
	}
	v.validateAssignmentOpTarget(n.GetFirstChild(), contextType)
	v.ValidateExpression(n.GetLastChild())
}

// validateAssignmentOpTarget validates the LHS of a compound assignment op,
// inc, or dec (stricter than validateLHS).
//
// Java: AstValidator.validateAssignmentOpTarget(Node, Token)
func (v *AstValidator) validateAssignmentOpTarget(lhs rhino.Node, contextType rhino.Token) {
	switch lhs.GetToken() {
	case rhino.Token_NAME:
		v.validateName(lhs)
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		v.validateGetPropGetElemInLHS(contextType, lhs)
	case rhino.Token_CAST:
		v.validatePropsAndCount(lhs, 1)
		v.validateAssignmentOpTarget(ast.MustCast(lhs).Expr().Node(), contextType)
	default:
		v.Violation("Invalid child for "+contextType.String()+" node", lhs)
	}
}

// Java: AstValidator.validateGetElem(Node)
func (v *AstValidator) validateGetElem(elem ast.GetElem) {
	n := elem.Node()
	if v.validatePropsAndCount(n, 2) {
		v.validatePropertyReferenceTarget(elem.Object().Node())
		v.ValidateExpression(elem.Index().Node())
		return
	}
	v.validatePropertyReferenceTarget(n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateOptChainGetElem(Node)
func (v *AstValidator) validateOptChainGetElem(elem ast.OptChainGetElem) {
	n := elem.Node()
	v.validateFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING, n)
	if v.validatePropsAndCount(n, 2) {
		v.ValidateExpression(elem.Object().Node())
		v.ValidateExpression(elem.Index().Node())
	} else {
		v.ValidateExpression(n.GetFirstChild())
		v.ValidateExpression(n.GetLastChild())
	}
	v.validateFirstNodeOfOptChain(elem.AsOptChain())
}

// Java: AstValidator.validateGetProp(Node)
func (v *AstValidator) validateGetProp(prop ast.GetProp) {
	n := prop.Node()
	v.validatePropertyReferenceTarget(prop.Object().Node())
	v.validatePropsAndArity(n)
	v.validateNonEmptyString(n)
}

// Java: AstValidator.validateOptChainGetProp(Node)
func (v *AstValidator) validateOptChainGetProp(prop ast.OptChainGetProp) {
	n := prop.Node()
	v.validateFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING, n)
	v.ValidateExpression(prop.Object().Node())
	v.validateFirstNodeOfOptChain(prop.AsOptChain())
	v.validatePropsAndArity(n)
	v.validateNonEmptyString(n)
}

// Java: AstValidator.validatePropertyReferenceTarget(Node)
func (v *AstValidator) validatePropertyReferenceTarget(objectNode rhino.Node) {
	if objectNode.IsSuper() {
		v.validateSuper(objectNode)
	} else {
		v.ValidateExpression(objectNode)
	}
}

// Java: AstValidator.validateRegExpLit(Node)
func (v *AstValidator) validateRegExpLit(re ast.RegExp) {
	n := re.Node()
	v.validateFeature(parser.FeatureSet_Feature_REGEXP_SYNTAX, n)
	v.validateProperties(n)
	if v.validateChildCountIn(n, 1, 2) {
		v.validateStringLit(re.Pattern().Node())
		if flags, ok := re.Flags(); ok {
			v.validateStringLit(flags.Node())
		}
		return
	}
	for c := range n.Children() {
		v.validateStringLit(c)
	}
}

// Java: AstValidator.validateStringLit(Node)
func (v *AstValidator) validateStringLit(n rhino.Node) {
	v.validateNodeType(rhino.Token_STRINGLIT, n)
	v.validatePropsAndArity(n)
	_ = n.GetString()
}

// Java: AstValidator.validateNumber(Node)
func (v *AstValidator) validateNumber(num ast.Number) {
	v.validatePropsAndArity(num.Node())
	_ = num.Value()
}

// Java: AstValidator.validateBigInt(Node)
func (v *AstValidator) validateBigInt(bi ast.BigInt) {
	n := bi.Node()
	v.validatePropsAndArity(n)
	_ = n.GetBigInt()
}

// Java: AstValidator.validateArrayLit(Node)
func (v *AstValidator) validateArrayLit(arr ast.ArrayLit) {
	if first, ok := arr.FirstElement(); ok {
		// Upstream Java AstValidator.validateArrayLit has an unconditional
		// break after the first element; preserve exact parity.
		v.validatePseudoExpression(first.Node(), rhino.Token_EMPTY, rhino.Token_ITER_SPREAD)
	}
}

// Java: AstValidator.validateObjectLit(Node)
func (v *AstValidator) validateObjectLit(obj ast.ObjectLit) {
	for member := range obj.Members() {
		v.validateObjectLitKey(member.Node())
	}
}

// Java: AstValidator.validateObjectLitKey(Node)
func (v *AstValidator) validateObjectLitKey(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETTER_DEF:
		v.validateObjectLitGetKey(ast.MustGetterDef(n))
	case rhino.Token_SETTER_DEF:
		v.validateObjectLitSetKey(ast.MustSetterDef(n))
	case rhino.Token_STRING_KEY:
		v.validateObjectLitStringKey(ast.MustStringKey(n))
	case rhino.Token_MEMBER_FUNCTION_DEF:
		def := ast.MustMemberFunctionDef(n)
		v.validateClassMember(n, false)
		if def.IsStatic() {
			v.Violation("Keys in an object literal should not be static.", n)
		}
	case rhino.Token_COMPUTED_PROP:
		v.validateObjectLitComputedPropKey(ast.MustComputedProp(n))
	case rhino.Token_OBJECT_SPREAD:
		v.validateObjectSpread(ast.MustObjectSpread(n))
	default:
		v.Violation("Expected object literal key expression but was "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateObjectLitGetKey(Node)
func (v *AstValidator) validateObjectLitGetKey(def ast.GetterDef) {
	n := def.Node()
	v.validateFeature(parser.FeatureSet_Feature_GETTER, n)
	v.validatePropsAndArity(n)
	v.validateObjectLiteralKeyName(n)
	v.validateAccessorFunction(n, def.Function().Node(), false)
}

// Java: AstValidator.validateObjectLitSetKey(Node)
func (v *AstValidator) validateObjectLitSetKey(def ast.SetterDef) {
	n := def.Node()
	v.validateFeature(parser.FeatureSet_Feature_SETTER, n)
	v.validatePropsAndArity(n)
	v.validateObjectLiteralKeyName(n)
	v.validateAccessorFunction(n, def.Function().Node(), true)
}

// Java: AstValidator.validateObjectLitStringKey(Node)
func (v *AstValidator) validateObjectLitStringKey(sk ast.StringKey) {
	n := sk.Node()
	v.validateObjectLiteralKeyName(n)
	v.validatePropsAndCount(n, 1)
	v.ValidateExpression(sk.Value().Node())
	if n.GetBooleanProp(rhino.Node_IS_SHORTHAND_PROPERTY) {
		v.validateFeature(parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES, n)
	}
}

// Java: AstValidator.validateObjectLitComputedPropKey(Node)
func (v *AstValidator) validateObjectLitComputedPropKey(prop ast.ComputedProp) {
	n := prop.Node()
	// TODO: b/499356003 - clean up this hack once computed getter/setters have
	// their own Feature separate from COMPUTED_PROPERTIES.
	if !(n.GetIsInClosureUnawareSubtree() && NodeUtil_IsGetOrSetKey(n)) {
		v.validateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	}
	if v.validatePropsAndArity(n) {
		v.ValidateExpression(prop.Key().Node())
		val, _ := prop.Value()
		v.ValidateExpression(val.Node())
		return
	}
	if first := n.GetFirstChild(); !first.IsNil() {
		v.ValidateExpression(first)
	}
	if n.HasMoreThanOneChild() {
		v.ValidateExpression(n.GetLastChild())
	}
}

// Java: AstValidator.validateObjectLitComputedPropGetKey(Node)
func (v *AstValidator) validateObjectLitComputedPropGetKey(prop ast.ComputedProp) {
	n := prop.Node()
	v.validateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.validateFeature(parser.FeatureSet_Feature_GETTER, n)
	var fnNode rhino.Node
	if v.validatePropsAndArity(n) {
		val, _ := prop.Value()
		fnNode = val.Node()
	} else if n.HasMoreThanOneChild() {
		fnNode = n.GetLastChild()
	}
	v.validateAccessorFunction(n, fnNode, false)
}

// Java: AstValidator.validateObjectLitComputedPropSetKey(Node)
func (v *AstValidator) validateObjectLitComputedPropSetKey(prop ast.ComputedProp) {
	n := prop.Node()
	v.validateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.validateFeature(parser.FeatureSet_Feature_SETTER, n)
	var fnNode rhino.Node
	if v.validatePropsAndArity(n) {
		val, _ := prop.Value()
		fnNode = val.Node()
	} else if n.HasMoreThanOneChild() {
		fnNode = n.GetLastChild()
	}
	v.validateAccessorFunction(n, fnNode, true)
}

func (v *AstValidator) validateAccessorFunction(keyNode, fnNode rhino.Node, isSetter bool) {
	if fnNode.IsNil() {
		return
	}
	fn, ok := ast.AsFunction(fnNode)
	if !ok {
		v.validateNodeType(rhino.Token_FUNCTION, fnNode)
		return
	}
	v.validateFunctionExpression(fnNode)
	if fn.NameNode().Ident() != "" {
		v.Violation("Expected unnamed function expression.", keyNode)
	}
	if isSetter {
		if _, ok := fn.Params().OnlyParam(); !ok {
			v.Violation("set methods must have exactly one parameter.", keyNode)
		}
	} else if fn.Params().HasParams() {
		v.Violation("get methods must not have parameters.", keyNode)
	}
}

// Java: AstValidator.validateObjectLiteralKeyName(Node)
func (v *AstValidator) validateObjectLiteralKeyName(n rhino.Node) {
	if n.IsQuotedStringKey() {
		_ = n.GetString()
	} else {
		v.validateNonEmptyString(n)
	}
}

// Java: AstValidator.validateIncDecOp(Node)
func (v *AstValidator) validateIncDecOp(inc ast.IncDec) {
	n := inc.Node()
	v.validatePropsAndCount(n, 1)
	v.validateAssignmentOpTarget(inc.Target().Node(), inc.Token())
}

// Java: AstValidator.validateBinaryOp(Node)
func (v *AstValidator) validateBinaryOp(bin ast.Binary) {
	n := bin.Node()
	if v.validatePropsAndCount(n, 2) {
		v.ValidateExpression(bin.Left().Node())
		v.ValidateExpression(bin.Right().Node())
		return
	}
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateTrinaryOp(Node)
func (v *AstValidator) validateTrinaryOp(hook ast.Hook) {
	v.validatePropsAndCount(hook.Node(), 3)
	v.ValidateExpression(hook.Cond().Node())
	v.ValidateExpression(hook.Then().Node())
	v.ValidateExpression(hook.Else().Node())
}
