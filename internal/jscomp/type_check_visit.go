// Port of com/google/javascript/jscomp/TypeCheck.java @ closure-compiler v20261021

/*
 * Copyright 2006 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func (tc *TypeCheck) shouldTraverseOnlyScriptAndFunction() bool {
	return !tc.debugTypeLogger.isLogging
}

// Java: TypeCheck.shouldTraverse(NodeTraversal, Node, Node)
func (tc *TypeCheck) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Start logging types if we're in debug mode and this node was annotated.
	if tc.debugTypeLogger.isLogging {
		tc.debugTypeLogger.MaybeStartLoggingAt(n)
	}

	if n.IsScript() {
		if NodeUtil_IsFromTypeSummary(n) {
			// Errors in type summary files are suppressed, so no use traversing them.
			return false
		}
		if J2clSourceUtils_IsJ2clSourceOfNode(n) {
			tc.subtypingMode = rhino.JSType_SubtypingMode_IGNORE_NULL_UNDEFINED
		} else {
			tc.subtypingMode = rhino.JSType_SubtypingMode_NORMAL
		}
		tc.validator.SetSubtypingMode(tc.subtypingMode)
	} else if fn, ok := ast.AsFunction(n); ok {
		// normal type checking
		// Check for a bug in IE9 (quirks mode) and earlier where bleeding function names would
		// refer to the wrong variable (i.e. "var x; var y = function x() { use(x); }" would pass
		// the 'x' from the outer scope, rather than the local alias bled into the function).
		if fnName, ok := fn.Name(); ok {
			outerScope := t.GetTypedScope()
			var_ := outerScope.GetVarByID(fnName.Node().GetStringID())
			if var_ != nil &&
				var_.Scope().HasSameContainerScope(outerScope) &&
				// Ideally, we would want to check whether the type in the scope
				// differs from the type being defined, but then the extern
				// redeclarations of built-in types generates spurious warnings.
				!typeCheck_isFunctionType(var_.GetType()) &&
				!TypeValidator_HasDuplicateDeclarationSuppression(tc.compiler, var_.GetNameNode()) {
				tc.Report(n, TypeCheck_FUNCTION_MASKS_VARIABLE, var_.GetName())
			}
		}

		// TODO(user): Only traverse the function's body. The function's
		// name and arguments are traversed by the scope creator, and ideally
		// should not be traversed by the type checker.
	}
	return true
}

// typeCheck_isFunctionType is Java's `type instanceof FunctionType`.
//
// Go-only: language-adapter — instanceof on an interface type.
func typeCheck_isFunctionType(type_ rhino.JSType) bool {
	_, ok := type_.(rhino.FunctionType)
	return ok
}

// Visit is the meat of the type checking. It is basically one big switch,
// with each case representing one type of parse tree node. The individual
// cases are usually pretty straightforward.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited; parent the parent of the node n.
//
// Java: TypeCheck.visit(NodeTraversal, Node, Node)
func (tc *TypeCheck) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// To be explicitly set to false if the node is not typeable.
	typeable := true

	if tb := n.GetJSTypeBase(); tb != nil && !tb.IsWellFormedTemplateValid() {
		nodeType := tb.Self()
		switch nodeType.GetTypeClass() {
		case rhino.JSTypeClass_ENUM_ELEMENT,
			rhino.JSTypeClass_FUNCTION,
			rhino.JSTypeClass_NAMED,
			rhino.JSTypeClass_UNION,
			rhino.JSTypeClass_TEMPLATIZED,
			rhino.JSTypeClass_TEMPLATE:
			tc.validator.checkWellFormedTemplatizedType(n, nodeType)
		default:
			tb.MarkWellFormedTemplateValidIfResolved()
		}
	}

	switch n.GetToken() {
	case rhino.Token_CAST:
		expr := ast.MustCast(n).Expr().Node()
		exprType := tc.GetJSType(expr)
		castType := tc.GetJSType(n)
		isObjLit := expr.IsObjectLit()

		// TODO(johnlenz): determine if we can limit object literals in some
		// way.
		if !isObjLit {
			tc.validator.ExpectCanCast(n, castType, exprType)
		}
		tc.EnsureTypedWithType(n, castType)

		expr.SetJSTypeBeforeCast(exprType)
		if isObjLit || castType.RestrictByNotNullOrUndefined().IsSubtypeOf(exprType) {
			expr.SetJSType(castType)
		}
	case rhino.Token_NAME:
		typeable = tc.VisitName(t, ast.MustName(n), parent)
	case rhino.Token_COMMA:
		tc.EnsureTypedWithType(n, tc.GetJSType(ast.MustBinary(n).Right().Node()))
	case rhino.Token_THIS:
		tc.EnsureTypedWithType(n, t.GetTypedScope().GetTypeOfThis())
	case rhino.Token_NULL:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NULL_TYPE)
	case rhino.Token_NUMBER:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NUMBER_TYPE)
	case rhino.Token_BIGINT:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BIGINT_TYPE)
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		// Object literal keys are handled with OBJECTLIT
	case rhino.Token_ARRAYLIT:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_ARRAY_TYPE)
	case rhino.Token_REGEXP:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_REGEXP_TYPE)
	case rhino.Token_GETPROP:
		tc.VisitGetProp(t, ast.MustGetProp(n))
		assign, ok := ast.AsAssign(parent)
		typeable = !(ok && assign.Target().Is(n))
	case rhino.Token_OPTCHAIN_GETPROP:
		tc.VisitOptChainGetProp(ast.MustOptChainGetProp(n))
	case rhino.Token_OPTCHAIN_GETELEM:
		tc.VisitOptChainGetElem(ast.MustOptChainGetElem(n))
	case rhino.Token_GETELEM:
		tc.VisitGetElem(ast.MustNormalOrOptChainGetElem(n))
		// The type of GETELEM is always unknown, so no point counting that.
		// If that unknown leaks elsewhere (say by an assignment to another
		// variable), then it will be counted.
		typeable = false
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		tc.VisitVar(t, ast.MustNameDeclaration(n))
		typeable = false
	case rhino.Token_NEW:
		tc.VisitNew(ast.MustNew(n))
	case rhino.Token_OPTCHAIN_CALL:
		// We reuse the `visitCall` functionality for OptChain call nodes because we don't report
		// an error for regular calls when the callee is null or undefined. However, we make sure
		// `typeable` isn't explicitly unset as OptChain nodes are always typed during inference.
		tc.VisitCall(t, ast.MustNormalOrOptChainCall(n))
	case rhino.Token_CALL:
		tc.VisitCall(t, ast.MustNormalOrOptChainCall(n))
		typeable = !parent.IsExprResult()
	case rhino.Token_RETURN:
		tc.VisitReturn(t, ast.MustReturn(n))
		typeable = false
	case rhino.Token_YIELD:
		tc.VisitYield(t, ast.MustYield(n))
	case rhino.Token_DEC, rhino.Token_INC:
		incDec := ast.MustIncDec(n)
		tc.CheckPropCreation(incDec.Target())
		target := incDec.Target().Node()
		if tc.GetJSType(n).IsNumber() {
			tc.validator.ExpectNumber(target, tc.GetJSType(target), "increment/decrement")
			tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NUMBER_TYPE)
		} else {
			tc.validator.ExpectBigIntOrNumber(target, tc.GetJSType(target), "increment/decrement")
		}
	case rhino.Token_VOID:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_VOID_TYPE)
	case rhino.Token_STRINGLIT, rhino.Token_TYPEOF, rhino.Token_TEMPLATELIT, rhino.Token_TEMPLATELIT_STRING:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_STRING_TYPE)
	case rhino.Token_TAGGED_TEMPLATELIT:
		tc.VisitTaggedTemplateLit(ast.MustTaggedTemplateLit(n))
		tc.EnsureTyped(n)
	case rhino.Token_BITNOT:
		tc.VisitBitwiseNOT(ast.MustUnary(n))
	case rhino.Token_POS:
		tc.VisitUnaryPlus(ast.MustUnary(n))
	case rhino.Token_NEG:
		tc.VisitUnaryMinus(ast.MustUnary(n))
	case rhino.Token_EQ, rhino.Token_NE, rhino.Token_SHEQ, rhino.Token_SHNE:
		op := ast.MustBinary(n)
		left, right := op.Left().Node(), op.Right().Node()

		if left.IsTypeOf() {
			if right.IsStringLit() {
				tc.CheckTypeofString(right, right.GetString())
			}
		} else if right.IsTypeOf() && left.IsStringLit() {
			tc.CheckTypeofString(left, left.GetString())
		}

		leftType := tc.GetJSType(left)
		rightType := tc.GetJSType(right)

		// We do not want to warn about explicit comparisons to VOID. People
		// often do this if they think their type annotations screwed up.
		//
		// We do want to warn about cases where people compare things like
		// (Array|null) == (Function|null)
		// because it probably means they screwed up.
		//
		// This heuristic here is not perfect, but should catch cases we
		// care about without too many false negatives.
		leftTypeRestricted := leftType.RestrictByNotNullOrUndefined()
		rightTypeRestricted := rightType.RestrictByNotNullOrUndefined()

		result := base.Tri_UNKNOWN
		if n.IsEQ() || n.IsNE() {
			result = leftTypeRestricted.TestForEquality(rightTypeRestricted)
			if n.IsNE() {
				result = result.Not()
			}
		} else {
			// SHEQ or SHNE
			if !leftTypeRestricted.CanTestForShallowEqualityWith(rightTypeRestricted) {
				if n.IsSHEQ() {
					result = base.Tri_FALSE
				} else {
					result = base.Tri_TRUE
				}
			}
		}

		if result != base.Tri_UNKNOWN {
			tc.Report(
				n, TypeCheck_DETERMINISTIC_TEST, leftType.String(), rightType.String(), result.String())
		}
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BOOLEAN_TYPE)
	case rhino.Token_LT, rhino.Token_LE, rhino.Token_GT, rhino.Token_GE:
		op := ast.MustBinary(n)
		leftSide, rightSide := op.Left().Node(), op.Right().Node()
		leftType := tc.GetJSType(leftSide)
		rightType := tc.GetJSType(rightSide)
		if rightType.IsUnknownType() {
			// validate comparable left
			tc.validator.ExpectUnknownOrComparable(leftSide, leftType, "left side of comparison")
		} else if leftType.IsUnknownType() {
			// validate comparable right
			tc.validator.ExpectUnknownOrComparable(rightSide, rightType, "right side of comparison")
		} else if rightType.IsBigIntOrNumber() {
			// validate left operand for numeric comparison
			tc.validator.ExpectBigIntOrNumber(leftSide, leftType, "left side of numeric comparison")
		} else if leftType.IsBigIntOrNumber() {
			// validate right operand for numeric comparison
			tc.validator.ExpectBigIntOrNumber(rightSide, rightType, "right side of numeric comparison")
		} else {
			errorMsg := "expected matching types in comparison"
			tc.validator.ExpectMatchingTypesStrict(n, leftType, rightType, errorMsg)
			if !leftType.MatchesNumberContext() || !rightType.MatchesNumberContext() {
				// Whether the comparison is numeric will be determined at runtime
				// each time the expression is evaluated. Regardless, both operands
				// should match a string context.
				message := "left side of comparison"
				tc.validator.ExpectString(leftSide, leftType, message)
				tc.validator.ExpectNotNullOrUndefined(
					t, leftSide, leftType, message, tc.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
				message = "right side of comparison"
				tc.validator.ExpectString(rightSide, rightType, message)
				tc.validator.ExpectNotNullOrUndefined(
					t, rightSide, rightType, message, tc.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
			}
		}
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BOOLEAN_TYPE)
	case rhino.Token_IN:
		op := ast.MustBinary(n)
		left, right := op.Left().Node(), op.Right().Node()
		rightType := tc.GetJSType(right)
		tc.validator.ExpectStringOrSymbol(left, tc.GetJSType(left), "left side of 'in'")
		tc.validator.ExpectObject(n, rightType, "'in' requires an object")
		if rightType.IsStruct() {
			tc.Report(right, TypeCheck_IN_USED_WITH_STRUCT)
		}
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BOOLEAN_TYPE)
	case rhino.Token_INSTANCEOF:
		op := ast.MustBinary(n)
		left, right := op.Left().Node(), op.Right().Node()
		rightType := tc.GetJSType(right).RestrictByNotNullOrUndefined()
		tc.validator.ExpectAnyObject(left, tc.GetJSType(left), "deterministic instanceof yields false")
		tc.validator.ExpectActualObject(right, rightType, "instanceof requires an object")
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BOOLEAN_TYPE)
	case rhino.Token_ASSIGN, rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_COALESCE:
		tc.VisitAssign(t, ast.MustAssignment(n))
		typeable = false
	case rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_EXPONENT:
		assignOp := ast.MustAssignOp(n)
		tc.CheckPropCreation(assignOp.Target())
		tc.VisitBinaryOperator(assignOp.Token(), n, assignOp.Target().Node(), assignOp.Value().Node())
	case rhino.Token_LSH, rhino.Token_RSH, rhino.Token_URSH, rhino.Token_DIV, rhino.Token_MOD, rhino.Token_BITOR,
		rhino.Token_BITXOR, rhino.Token_BITAND, rhino.Token_SUB, rhino.Token_ADD, rhino.Token_MUL, rhino.Token_EXPONENT:
		op := ast.MustBinary(n)
		tc.VisitBinaryOperator(op.Token(), n, op.Left().Node(), op.Right().Node())
	case rhino.Token_TRUE, rhino.Token_FALSE, rhino.Token_NOT, rhino.Token_DELPROP:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BOOLEAN_TYPE)
	case rhino.Token_CASE:
		// The parent of a case is the body of its switch.
		switchNode := ast.MustSwitch(parent.GetParent())
		switchConditionType := tc.GetJSType(switchNode.Discriminant().Node())
		caseType := tc.GetJSType(ast.MustCase(n).Value().Node())
		tc.validator.ExpectSwitchMatchesCase(n, switchConditionType, caseType)
		typeable = false
	case rhino.Token_WITH:
		child := ast.MustWith(n).Object().Node()
		childType := tc.GetJSType(child)
		tc.validator.ExpectObject(child, childType, "with requires an object")
		typeable = false
	case rhino.Token_FUNCTION:
		tc.VisitFunction(n)
	case rhino.Token_CLASS:
		tc.VisitClass(ast.MustClass(n))
	case rhino.Token_MODULE_BODY:
		tc.VisitModuleBody(t, ast.MustModuleBody(n))
	// These nodes have no interesting type behavior.
	// These nodes require data flow analysis.
	case rhino.Token_PARAM_LIST,
		rhino.Token_STRING_KEY,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_COMPUTED_PROP,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_COMPUTED_FIELD_DEF,
		rhino.Token_LABEL,
		rhino.Token_LABEL_NAME,
		rhino.Token_SWITCH,
		rhino.Token_SWITCH_BODY,
		rhino.Token_BREAK,
		rhino.Token_CATCH,
		rhino.Token_TRY,
		rhino.Token_SCRIPT,
		rhino.Token_EXPORT,
		rhino.Token_EXPORT_SPEC,
		rhino.Token_EXPORT_SPECS,
		rhino.Token_IMPORT,
		rhino.Token_IMPORT_SPEC,
		rhino.Token_IMPORT_SPECS,
		rhino.Token_IMPORT_STAR,
		rhino.Token_EXPR_RESULT,
		rhino.Token_BLOCK,
		rhino.Token_ROOT,
		rhino.Token_EMPTY,
		rhino.Token_DEFAULT_CASE,
		rhino.Token_CONTINUE,
		rhino.Token_DEBUGGER,
		rhino.Token_THROW,
		rhino.Token_DO,
		rhino.Token_IF,
		rhino.Token_WHILE,
		rhino.Token_FOR,
		rhino.Token_TEMPLATELIT_SUB,
		rhino.Token_ITER_REST,
		rhino.Token_OBJECT_REST,
		rhino.Token_DESTRUCTURING_LHS:
		typeable = false
	case rhino.Token_DYNAMIC_IMPORT:
		tc.VisitDynamicImport(t, ast.MustDynamicImport(n))
	case rhino.Token_ARRAY_PATTERN:
		tc.EnsureTyped(n)
		tc.validator.ExpectAutoboxesToIterable(
			n, tc.GetJSType(n), "array pattern destructuring requires an Iterable")
	case rhino.Token_OBJECT_PATTERN:
		tc.VisitObjectPattern(ast.MustObjectPattern(n))
	case rhino.Token_DEFAULT_VALUE:
		defaultValue := ast.MustDefaultValue(n)
		target, value := defaultValue.Target(), defaultValue.Value().Node()
		valueType := tc.GetJSType(value)
		tc.CheckCanAssignToWithScope(
			t,
			n,
			target,
			valueType,
			/* info= */ nil,
			"default value has wrong type")

		// Every other usage of a destructuring pattern is checked while visiting the pattern,
		// but default values are different because they are a conditional assignment and the
		// pattern is not given the default value's type
		switch target.Token() {
		case rhino.Token_ARRAY_PATTERN:
			tc.validator.ExpectAutoboxesToIterable(
				value, valueType, "array pattern destructuring requires an Iterable")
		case rhino.Token_OBJECT_PATTERN:
			// Verify that the value is not null/undefined, since those can't be destructured.
			tc.validator.ExpectObject(
				value, valueType, "cannot destructure a 'null' or 'undefined' default value")
		}

		typeable = false
	case rhino.Token_CLASS_MEMBERS:
		classType := parent.GetJSType().ToMaybeFunctionType()
		typ := classType.GetInstanceType()
		for member := range ast.MustClassBody(n).Members() {
			tc.VisitObjectOrClassLiteralKey(member.Node(), parent, typ)
			if getterOrSetter, ok := member.AsGetterOrSetterDef(); ok {
				tc.CheckGetterOrSetterType(getterOrSetter, classType)
			}
		}
		typeable = false
	case rhino.Token_FOR_IN:
		obj := ast.MustForIn(n).Object().Node()
		if tc.GetJSType(obj).IsStruct() {
			tc.Report(obj, TypeCheck_IN_USED_WITH_STRUCT)
		}
		typeable = false
	case rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		tc.EnsureTyped(ast.MustForOfLoop(n).Iterable().Node())
		typeable = false
		// These nodes are typed during the type inference.
	case rhino.Token_SUPER, rhino.Token_NEW_TARGET, rhino.Token_IMPORT_META, rhino.Token_AWAIT, rhino.Token_AND,
		rhino.Token_HOOK, rhino.Token_OR, rhino.Token_COALESCE:
		tc.EnsureTyped(n)
	case rhino.Token_OBJECTLIT:
		// If this is an enum, then give that type to the objectlit as well.
		if pt := parent.GetJSType(); pt != nil && pt.IsEnumType() {
			tc.EnsureTypedWithType(n, pt)
		} else {
			tc.EnsureTyped(n)
		}
		typ := tc.GetJSType(n)
		for member := range ast.MustObjectLit(n).Members() {
			tc.VisitObjectOrClassLiteralKey(member.Node(), n, typ)
		}
	case rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD:
		tc.CheckSpread(n)
		typeable = false
	default:
		tc.Report(n, TypeCheck_UNEXPECTED_TOKEN, n.GetToken().String())
		tc.EnsureTyped(n)
	}

	if p := n.GetParent(); !p.IsNil() {
		switch p.GetToken() {
		case rhino.Token_FUNCTION:
			// Visit the body of blockless arrow functions
			if !n.IsBlock() && ast.MustFunction(p).Body().Is(n) {
				tc.VisitImplicitReturnExpression(t, n)
			}
		case rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
			// Visit the loop initializer of a for-of loop
			// We do this check here, instead of when visiting FOR_OF, in order to get the correct
			// TypedScope.
			if forOf := ast.MustForOfLoop(p); forOf.Target().Is(n) {
				tc.CheckForOfTypes(t, forOf)
			}
		}
	}

	// Don't count externs since the user's code may not even use that part.
	typeable = typeable && !tc.inExterns

	// Record typing logs and metrics.
	if typeable {
		if !tc.debugTypeLogger.parentNode.IsNil() {
			tc.debugTypeLogger.MaybeLogTypeOfNode(n)
		}
		tc.DoPercentTypedAccounting(n)
	}

	tc.CheckJsdocInfoContainsObjectWithBadKey(n)

	// If this is the node which started logging types, stop logging.
	if !tc.debugTypeLogger.parentNode.IsNil() {
		tc.debugTypeLogger.StopLoggingIfThisIsWhereWeStarted(n)
	}
}

// DoPercentTypedAccounting counts the given node in the typed statistics.
//
// n is a node that should be typed.
//
// Java: TypeCheck.doPercentTypedAccounting(Node)
func (tc *TypeCheck) DoPercentTypedAccounting(n rhino.Node) {
	type_ := n.GetJSType()
	if type_ == nil {
		tc.nullCount++
	} else if type_.IsUnknownType() {
		if tc.reportUnknownTypes && !n.GetParent().IsExprResult() {
			tc.compiler.Report(JSError_MakeWithNode(n, TypeCheck_UNKNOWN_EXPR_TYPE))
		}
		tc.unknownCount++
	} else {
		tc.typedCount++
	}
}
