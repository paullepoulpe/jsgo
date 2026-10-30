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
	"math"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var TypeCheck_GOOG_INHERITS = rhino.QualifiedName_Of("goog.inherits")

// CheckCallConventions validates class-defining calls. Because JS has no
// 'native' syntax for defining classes, we need to do this manually.
//
// Java: TypeCheck.checkCallConventions(NodeTraversal, Node)
func (tc *TypeCheck) CheckCallConventions(t *NodeTraversal, call ast.NormalOrOptChainCall) {
	n := call.Node()
	relationship := tc.compiler.GetCodingConvention().GetClassesDefinedByCall(n)
	scope := t.GetTypedScope()
	if relationship != nil {
		superClass := rhino.StaticTypedScope_LookupQualifiedNameString(scope, relationship.superclassName)
		superClassInstanceType := TypeValidator_GetInstanceOfCtor(superClass)
		subClass := rhino.StaticTypedScope_LookupQualifiedNameString(scope, relationship.subclassName)
		subClassInstance := TypeValidator_GetInstanceOfCtor(subClass)
		if relationship.type_ == CodingConvention_SubclassType_INHERITS &&
			superClassInstanceType != nil &&
			!superClassInstanceType.IsEmptyType() &&
			subClassInstance != nil &&
			!subClassInstance.IsEmptyType() {
			if callee := call.Callee().Node(); callee.IsQualifiedName() &&
				TypeCheck_GOOG_INHERITS.Matches(callee) &&
				subClass.ToMaybeFunctionType() != nil &&
				!subClass.ToMaybeFunctionType().GetSource().IsNil() &&
				subClass.ToMaybeFunctionType().GetSource().IsClass() {
				tc.compiler.Report(JSError_MakeWithNode(n, TypeCheck_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS))
			}
			tc.validator.ExpectSuperType(n, superClassInstanceType, subClassInstance)
		}
	}
}

// VisitCall visits a CALL node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited.
//
// Java: TypeCheck.visitCall(NodeTraversal, Node)
func (tc *TypeCheck) VisitCall(t *NodeTraversal, call ast.NormalOrOptChainCall) {
	n := call.Node()
	tc.CheckCallConventions(t, call)

	callee := call.Callee()
	childType := tc.GetJSType(callee.Node()).RestrictByNotNullOrUndefined()

	if !childType.CanBeCalled() {
		tc.Report(n, TypeCheck_NOT_CALLABLE, childType.String())
		tc.EnsureTyped(n)
		return
	}

	// A couple of types can be called as if they were functions.
	// If it is a function type, then validate parameters.
	if childType.IsFunctionType() {
		functionType := childType.ToMaybeFunctionType()

		// Non-native constructors should not be called directly
		// unless they specify a return type
		if functionType.IsConstructor() &&
			!functionType.IsNativeObjectType() &&
			(functionType.GetReturnType().IsUnknownType() ||
				functionType.GetReturnType().IsVoidType()) &&
			!callee.Node().IsSuper() {
			tc.Report(n, TypeCheck_CONSTRUCTOR_NOT_CALLABLE, childType.String())
		}

		// Functions with explicit 'this' types must be called in a GETPROP or GETELEM.
		_, calleeIsGet := callee.AsNormalOrOptChainGet()
		if functionType.IsOrdinaryFunction() && !calleeIsGet {
			receiverType := functionType.GetTypeOfThis()
			if receiverType.IsUnknownType() ||
				receiverType.IsAllType() ||
				receiverType.IsVoidType() ||
				(receiverType.IsObjectType() && receiverType.ToObjectType().IsNativeObjectType()) {
				// Allow these special cases.
			} else {
				tc.Report(n, TypeCheck_EXPECTED_THIS_TYPE, functionType.String())
			}
		}

		tc.VisitArgumentList(call.AsCallOrNew(), functionType)
		tc.EnsureTypedWithType(n, functionType.GetReturnType())
	} else {
		tc.EnsureTyped(n)
	}

	// TODO(nicksantos): Add something to check for calls of RegExp objects,
	// which is not supported by IE. Either say something about the return type
	// or warn about the non-portability of the call or both.
}

// VisitArgumentList visits the parameters of a CALL or a NEW node.
//
// Java: TypeCheck.visitArgumentList(Node, FunctionType)
func (tc *TypeCheck) VisitArgumentList(call ast.CallOrNew, functionType rhino.FunctionType) {
	tc.CheckArgumentsMatchParameters(call.Node(), functionType, functionType.GetParameters(), 0)
}

// CheckArgumentsMatchParameters checks that a list of arguments match a list
// of formal parameters.
//
// If given a TAGGED_TEMPLATE_LIT, the given Iterator should only contain the
// parameters corresponding to the actual template lit sub arguments, skipping
// over the first parameter.
//
// firstParameterIndex is the index of the first parameter in the given
// Iterator in the function type's parameter list.
//
// Go-only: the Java gets the arguments as an Iterator over
// NodeUtil.getInvocationArgsAsIterable(call); here they are read from call,
// a CallOrNew (its arguments) or a TaggedTemplateLit (the substitutions of
// its template), and checkArgument is the body of the Java loop. The
// parameters are a slice (the elements the Java Iterator has left to
// yield).
//
// Java: TypeCheck.checkArgumentsMatchParameters(Node, FunctionType, Iterator<Node>, Iterator<Parameter>, int)
func (tc *TypeCheck) CheckArgumentsMatchParameters(
	call rhino.Node,
	functionType rhino.FunctionType,
	parameters []rhino.FunctionType_Parameter,
	firstParameterIndex int) {

	spreadArgumentCount := 0
	normalArgumentCount := firstParameterIndex
	checkArgumentTypeAgainstParameter := true
	var parameter rhino.FunctionType_Parameter
	nextParameter := 0
	checkArgument := func(argument ast.Arg) {
		// Count normal & spread arguments.
		if _, isSpread := argument.AsIterSpread(); isSpread {
			// we have some form of this case
			// someCall(arg1, arg2, ...firstSpreadExpression, argN, ...secondSpreadExpression)
			spreadArgumentCount++
			// Once we see a spread parameter, we can no longer match up arguments with parameters.
			checkArgumentTypeAgainstParameter = false
		} else {
			normalArgumentCount++
		}

		// Get the next parameter, if we're still matching parameters and arguments.
		if checkArgumentTypeAgainstParameter {
			if nextParameter < len(parameters) {
				parameter = parameters[nextParameter]
				nextParameter++
			} else if !parameter.IsZero() && parameter.IsVariadic() {
				// use varargs for all remaining parameters
			} else {
				// else we ran out of parameters and will report that after this loop
				parameter = rhino.FunctionType_Parameter{}
				checkArgumentTypeAgainstParameter = false
			}
		}

		if checkArgumentTypeAgainstParameter {
			tc.validator.ExpectArgumentMatchesParameter(
				argument.Node(), tc.GetJSType(argument.Node()), parameter.GetJSType(), call, normalArgumentCount)
		}
	}

	var callee ast.Expr
	if tagged, ok := ast.AsTaggedTemplateLit(call); ok {
		callee = tagged.Tag()
		for part := range tagged.Template().Parts() {
			if sub, ok := part.AsTemplateLitSub(); ok {
				checkArgument(sub.Expr().AsArg())
			}
		}
	} else {
		invocation := ast.MustCallOrNew(call)
		callee = invocation.Callee()
		for argument := range invocation.Args() {
			checkArgument(argument)
		}
	}

	minArity := functionType.GetMinArity()
	maxArity := functionType.GetMaxArity()

	if spreadArgumentCount > 0 {
		if normalArgumentCount > maxArity {
			// We cannot reliably check whether the total argument count is wrong, but we can at
			// least tell if there are more arguments than the function can handle even ignoring the
			// spreads.
			maxArityText := ""
			if maxArity != math.MaxInt32 {
				maxArityText = " and no more than " + strconv.Itoa(maxArity) + " argument(s)"
			}
			tc.Report(
				call,
				TypeCheck_WRONG_ARGUMENT_COUNT,
				tc.typeRegistry.GetReadableTypeNameNoDeref(callee.Node()),
				"at least "+strconv.Itoa(normalArgumentCount),
				strconv.Itoa(minArity),
				maxArityText)
		}
	} else {
		if minArity > normalArgumentCount || maxArity < normalArgumentCount {
			maxArityText := ""
			if maxArity != math.MaxInt32 {
				maxArityText = " and no more than " + strconv.Itoa(maxArity) + " argument(s)"
			}
			tc.Report(
				call,
				TypeCheck_WRONG_ARGUMENT_COUNT,
				tc.typeRegistry.GetReadableTypeNameNoDeref(callee.Node()),
				strconv.Itoa(normalArgumentCount),
				strconv.Itoa(minArity),
				maxArityText)
		}
	}
}

// VisitImplicitReturnExpression visits an arrow function expression body.
//
// Java: TypeCheck.visitImplicitReturnExpression(NodeTraversal, Node)
func (tc *TypeCheck) VisitImplicitReturnExpression(t *NodeTraversal, exprNode rhino.Node) {
	enclosingFunction := t.GetEnclosingFunction()
	jsType := tc.GetJSType(enclosingFunction)
	if jsType.IsFunctionType() {
		functionType := jsType.ToMaybeFunctionType()

		expectedReturnType := functionType.GetReturnType()
		// if no return type is specified, undefined must be returned
		// (it's a void function)
		if expectedReturnType == nil {
			expectedReturnType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		} else if enclosingFunction.IsAsyncFunction() {
			// Unwrap the async function's declared return type.
			expectedReturnType = Promises_CreateAsyncReturnableType(tc.typeRegistry, expectedReturnType)
		}

		// Fetch the returned value's type
		actualReturnType := tc.GetJSType(exprNode)

		tc.validator.ExpectCanAssignTo(
			exprNode, actualReturnType, expectedReturnType, "inconsistent return type")
	}
}

// VisitReturn visits a RETURN node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited.
//
// Java: TypeCheck.visitReturn(NodeTraversal, Node)
func (tc *TypeCheck) VisitReturn(t *NodeTraversal, ret ast.Return) {
	enclosingFunction := t.GetEnclosingFunction()

	jsType := tc.GetJSType(enclosingFunction)

	if jsType.IsFunctionType() {
		functionType := jsType.ToMaybeFunctionType()

		returnType := functionType.GetReturnType()
		// if no return type is specified, undefined must be returned
		// (it's a void function)
		if returnType == nil {
			returnType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		} else if enclosingFunction.IsGeneratorFunction() {
			// Unwrap the template variable from a generator function's declared return type.
			// e.g. if returnType is "Generator<string, number, void>", the generator should return
			// "number".
			returnType = JsIterables_GetReturnElementType(returnType, tc.typeRegistry)

			if enclosingFunction.IsAsyncGeneratorFunction() {
				// Can return x|IThenable<x> in an AsyncGenerator<x>, no await needed. Note that we must
				// first wrap the type in IThenable as createAsyncReturnableType will map a non-IThenable
				// to `?`.
				returnType =
					Promises_CreateAsyncReturnableType(
						tc.typeRegistry, Promises_WrapInIThenable(tc.typeRegistry, returnType))
			}
		} else if enclosingFunction.IsAsyncFunction() {
			// e.g. `!Promise<string>` => `string|!IThenable<string>`
			// We transform the expected return type rather than the actual return type so that the
			// extual return type is always reported to the user. This was felt to be clearer.
			returnType = Promises_CreateAsyncReturnableType(tc.typeRegistry, returnType)
		} else if returnType.IsVoidType() && functionType.IsConstructor() {
			// Allow constructors to use empty returns for flow control.
			if _, hasValue := ret.Value(); !hasValue {
				return
			}

			// Allow constructors to return its own instance type
			returnType = functionType.GetInstanceType()
		}

		// fetching the returned value's type
		var valueNode rhino.Node
		var actualReturnType rhino.JSType
		if value, ok := ret.Value(); ok {
			valueNode = value.Node()
			actualReturnType = tc.GetJSType(valueNode)
		} else {
			actualReturnType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
			valueNode = ret.Node()
		}

		// verifying
		tc.validator.ExpectCanAssignTo(
			valueNode, actualReturnType, returnType, "inconsistent return type")
	}
}

// VisitYield visits a YIELD node.
//
// Java: TypeCheck.visitYield(NodeTraversal, Node)
func (tc *TypeCheck) VisitYield(t *NodeTraversal, yield ast.Yield) {
	n := yield.Node()
	jsType := tc.GetJSType(t.GetEnclosingFunction())

	declaredYieldType := tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	if jsType.IsFunctionType() {
		functionType := jsType.ToMaybeFunctionType()
		returnType := functionType.GetReturnType()
		declaredYieldType = JsIterables_GetElementType(returnType, tc.typeRegistry)

		if t.GetEnclosingFunction().IsAsyncGeneratorFunction() {
			// Can yield x|IThenable<x> in an AsyncGenerator<x>, no await needed. Note that we must
			// first wrap the type in IThenable as createAsyncReturnableType will map a non-IThenable to
			// `?`.
			declaredYieldType =
				Promises_CreateAsyncReturnableType(
					tc.typeRegistry, Promises_WrapInIThenable(tc.typeRegistry, declaredYieldType))
		}
	}

	// fetching the yielded value's type
	var valueNode rhino.Node
	var actualYieldType rhino.JSType
	if value, ok := yield.Value(); ok {
		valueNode = value.Node()
		actualYieldType = tc.GetJSType(valueNode)
	} else {
		actualYieldType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		valueNode = n
	}

	if n.IsYieldAll() {
		if t.GetEnclosingFunction().IsAsyncGeneratorFunction() {
			maybeActualYieldType :=
				tc.validator.ExpectAutoboxesToIterableOrAsyncIterable(
					n, actualYieldType, "Expression yield* expects an iterable or async iterable")
			if maybeActualYieldType == nil {
				// don't do any further typechecking of the yield* type.
				return
			}
			actualYieldType = maybeActualYieldType
		} else {
			if !tc.validator.ExpectAutoboxesToIterable(
				n, actualYieldType, "Expression yield* expects an iterable") {
				// don't do any further typechecking of the yield* type.
				return
			}
			templateType := tc.typeRegistry.GetIterableValueTemplate()
			actualYieldType =
				actualYieldType.Autobox().GetTemplateTypeMap().GetResolvedTemplateType(templateType)
		}
	}

	// verifying
	tc.validator.ExpectCanAssignTo(
		valueNode,
		actualYieldType,
		declaredYieldType,
		"Yielded type does not match declared return type.")
}

// Java: TypeCheck.visitTaggedTemplateLit(Node)
func (tc *TypeCheck) VisitTaggedTemplateLit(tagged ast.TaggedTemplateLit) {
	n := tagged.Node()
	tag := tagged.Tag().Node()
	tagType := tag.GetJSType().RestrictByNotNullOrUndefined()

	if !tagType.CanBeCalled() {
		tc.Report(n, TypeCheck_NOT_CALLABLE, tagType.String())
		return
	} else if !tagType.IsFunctionType() {
		// A few types, like the unknown, regexp, and bottom types, can be called as if they are
		// functions. Return if we have one of those types that is not actually a known function.
		return
	}

	tagFnType := tagType.ToMaybeFunctionType()
	parameters := tagFnType.GetParameters()

	// The tag function gets an array of all the template lit substitutions as its first argument,
	// but there's no actual AST node representing that array so we typecheck it separately from
	// the other tag arguments.

	// Validate that the tag function takes at least one parameter
	if len(parameters) == 0 {
		tc.Report(
			n,
			TypeCheck_WRONG_ARGUMENT_COUNT,
			tc.typeRegistry.GetReadableTypeNameNoDeref(tag),
			strconv.Itoa(NodeUtil_GetInvocationArgsCount(n)),
			"0",
			" and no more than 0 argument(s)")
		return
	}

	// Validate that the first parameter is a supertype of ITemplateArray
	firstParameter := parameters[0]
	parameters = parameters[1:]
	parameterType := firstParameter.GetJSType().RestrictByNotNullOrUndefined()
	if parameterType != nil {
		tc.validator.ExpectITemplateArraySupertype(
			tag, parameterType, "Invalid type for the first parameter of tag function")
	}

	// Validate the remaining parameters (the template literal substitutions)
	tc.CheckArgumentsMatchParameters(n, tagFnType, parameters, 1)
}

// VisitBinaryOperator unifies the type checking involved in the core binary
// operators and the corresponding assignment operators. The representation
// used internally is such that common code can handle both kinds of
// operators easily.
//
// op is the operator; n the node being checked.
//
// Go-only: n is an ast.Binary or an ast.AssignOp, which share no view, so
// the caller passes the operands: left and right, or target and value.
//
// Java: TypeCheck.visitBinaryOperator(Token, Node)
func (tc *TypeCheck) VisitBinaryOperator(op rhino.Token, n rhino.Node, left rhino.Node, right rhino.Node) {
	operatorType := tc.GetJSType(n)
	leftType := tc.GetJSType(left)
	rightType := tc.GetJSType(right)
	if operatorType.IsNoType() {
		// An attempt to mix bigint with other types
		tc.Report(
			n,
			TypeCheck_BINARY_OPERATION,
			stringPtrOrNull(NodeUtil_OpToStr(n.GetToken())),
			leftType.String(),
			rightType.String())
		return
	}
	switch op {
	case rhino.Token_ASSIGN_LSH, rhino.Token_ASSIGN_RSH, rhino.Token_LSH, rhino.Token_RSH, rhino.Token_ASSIGN_URSH,
		rhino.Token_URSH:
		if operatorType.IsNumber() {
			// TypeInference set the operator type to 'number', so we know bigint isn't involved.
			// NOTE: >>> and >>>= aren't valid operations for bigint, so TypeInference ignores the
			//     operand types for it and always sets the operator's type to 'number'. Thus,
			//     if one of the operands is a bigint, we'll end up reporting that as an error here.
			opStr := stringPtrOrNull(NodeUtil_OpToStr(n.GetToken()))
			if !leftType.MatchesNumberContext() {
				tc.Report(left, TypeCheck_BIT_OPERATION, opStr, leftType.String())
			} else {
				tc.validator.ExpectNumberStrict(n, leftType, "operator "+opStr)
			}
			if !rightType.MatchesNumberContext() {
				tc.Report(right, TypeCheck_BIT_OPERATION, opStr, rightType.String())
			} else {
				tc.validator.ExpectNumberStrict(n, rightType, "operator "+opStr)
			}
		} else {
			tc.validator.ExpectBigIntOrNumber(left, leftType, "left operand")
			tc.validator.ExpectBigIntOrNumber(right, rightType, "right operand")
		}
	case rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_DIV,
		rhino.Token_MOD,
		rhino.Token_MUL,
		rhino.Token_SUB,
		rhino.Token_EXPONENT:
		if operatorType.IsNumber() {
			// TypeInference set the operator type to 'number', so we know bigint isn't involved.
			tc.validator.ExpectNumber(left, leftType, "left operand")
			tc.validator.ExpectNumber(right, rightType, "right operand")
		} else {
			tc.validator.ExpectBigIntOrNumber(left, leftType, "left operand")
			tc.validator.ExpectBigIntOrNumber(right, rightType, "right operand")
		}
	case rhino.Token_ASSIGN_BITAND, rhino.Token_ASSIGN_BITXOR, rhino.Token_ASSIGN_BITOR, rhino.Token_BITAND,
		rhino.Token_BITXOR, rhino.Token_BITOR:
		if operatorType.IsNumber() {
			// This condition is meant to catch any old cases (where bigint isn't involved)
			tc.validator.ExpectBitwiseable(left, leftType, "bad left operand to bitwise operator")
			tc.validator.ExpectBitwiseable(right, rightType, "bad right operand to bitwise operator")
		} else {
			tc.validator.ExpectBigIntOrNumber(left, leftType, "bad left operand to bitwise operator")
			tc.validator.ExpectBigIntOrNumber(right, rightType, "bad right operand to bitwise operator")
		}
	case rhino.Token_ASSIGN_ADD, rhino.Token_ADD:
	default:
		tc.Report(n, TypeCheck_UNEXPECTED_TOKEN, op.String())
	}
	tc.EnsureTyped(n)
}

// typeCheck_inputProvider is the lambda `(b) -> t.getInput()`, a
// CompilerInputProvider.
//
// Go-only: language-adapter — a lambda implementing a Go interface.
type typeCheck_inputProvider struct {
	t *NodeTraversal
}

func (p typeCheck_inputProvider) GetInput(*rhino.InputId) *CompilerInput {
	return p.t.GetInput()
}

// VisitModuleBody validates the implicit assignment to the global for a
// legacy goog.module.
//
// Java: TypeCheck.visitModuleBody(NodeTraversal, Node)
func (tc *TypeCheck) VisitModuleBody(t *NodeTraversal, body ast.ModuleBody) {
	moduleBody := body.Node()
	associatedModule :=
		ModuleImportResolver_GetModuleFromScopeRoot(
			tc.compiler.GetModuleMap(), typeCheck_inputProvider{t}, moduleBody)
	if !associatedModule.Metadata().IsLegacyGoogModule() {
		return
	}
	moduleName := rhino.QualifiedName_Of(*associatedModule.ClosureNamespace())
	// The first statement of a legacy goog.module is its goog.module call.
	googModuleCall, ok := body.FirstStatement()
	jutil.CheckState(ok)
	if moduleName.IsSimple() {
		globalVar := tc.topScope.GetVar(moduleName.Join())
		tc.validator.ExpectCanAssignTo(
			googModuleCall.Node(), moduleBody.GetJSType(), globalVar.GetType(), "legacy goog.module export")
	} else {
		parentType := tc.topScope.LookupQualifiedName(moduleName.GetOwner())
		var parentObjectType rhino.ObjectType
		if parentType != nil {
			parentObjectType = parentType.ToMaybeObjectType()
		}
		if parentObjectType == nil {
			return
		}
		tc.validator.expectCanAssignToPropertyOfQualifiedOwner(
			googModuleCall.Node(),
			moduleBody.GetJSType(),
			parentObjectType.GetPropertyType(moduleName.GetComponent()),
			parentObjectType,
			moduleName.GetOwner(),
			moduleName.GetComponent())
	}
}

// VisitDynamicImport validates that a dynamic import statement has a single
// child of type string.
//
// Java: TypeCheck.visitDynamicImport(NodeTraversal, Node)
func (tc *TypeCheck) VisitDynamicImport(t *NodeTraversal, dynamicImport ast.DynamicImport) {
	tc.EnsureTypedWithType(dynamicImport.Node(), tc.GetNativeType(rhino.JSTypeNative_PROMISE_TYPE))

	importSpecifier := dynamicImport.Specifier().Node()
	importSpecifierType := importSpecifier.GetJSType()
	if importSpecifierType == nil {
		tc.EnsureTypedWithNative(importSpecifier, rhino.JSTypeNative_STRING_TYPE)
	} else {
		tc.validator.ExpectNotNullOrUndefined(
			t,
			importSpecifier,
			importSpecifierType,
			"dynamic import specifier",
			tc.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	}
}
