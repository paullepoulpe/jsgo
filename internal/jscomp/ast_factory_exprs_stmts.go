// Port of com/google/javascript/jscomp/AstFactory.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"iter"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ExprResult returns a new EXPR_RESULT node.
//
// Statements have no type information, so this is functionally the same as
// calling IR.exprResult(expr). It exists so that a pass can be consistent
// about always using AstFactory to create new nodes.
//
// Java: AstFactory.exprResult(Node)
func (a *AstFactory) ExprResult(expr rhino.Node) rhino.Node {
	// TODO(bradfordcsmith): This method should not be calling .srcref()
	return rhino.IR_ExprResult(expr).Srcref(expr)
}

// CreateEmpty returns a new EMPTY node.
//
// EMPTY Nodes have no type information, so this is functionally the same as
// calling IR.empty(). It exists so that a pass can be consistent about
// always using AstFactory to create new nodes.
//
// Java: AstFactory.createEmpty()
func (a *AstFactory) CreateEmpty() rhino.Node {
	return rhino.IR_Empty()
}

// CreateBlock returns a new BLOCK node.
//
// Blocks have no type information, so this is functionally the same as
// calling IR.block(statements). It exists so that a pass can be consistent
// about always using AstFactory to create new nodes.
//
// Java: AstFactory.createBlock(Node...)
func (a *AstFactory) CreateBlock(statements ...rhino.Node) rhino.Node {
	return rhino.IR_BlockOfNodes(statements...)
}

// CreateIf returns a new IF node.
//
// Blocks have no type information, so this is functionally the same as
// calling IR.ifNode(cond, then). It exists so that a pass can be consistent
// about always using AstFactory to create new nodes.
//
// Java: AstFactory.createIf(Node, Node)
func (a *AstFactory) CreateIf(cond rhino.Node, then rhino.Node) rhino.Node {
	return rhino.IR_IfNode(cond, then)
}

// CreateIfWithElse returns a new IF node.
//
// Blocks have no type information, so this is functionally the same as
// calling IR.ifNode(cond, then, elseNode). It exists so that a pass can be
// consistent about always using AstFactory to create new nodes.
//
// Java: AstFactory.createIf(Node, Node, Node)
func (a *AstFactory) CreateIfWithElse(cond rhino.Node, then rhino.Node, elseNode rhino.Node) rhino.Node {
	return rhino.IR_IfNodeWithElse(cond, then, elseNode)
}

// CreateFor returns a new FOR node.
//
// Blocks have no type information, so this is functionally the same as
// calling IR.forNode(init, cond, incr, body). It exists so that a pass can
// be consistent about always using AstFactory to create new nodes.
//
// Java: AstFactory.createFor(Node, Node, Node, Node)
func (a *AstFactory) CreateFor(init rhino.Node, cond rhino.Node, incr rhino.Node, body rhino.Node) rhino.Node {
	return rhino.IR_ForNode(init, cond, incr, body)
}

// CreateBreak returns a new BREAK node.
//
// Breaks have no type information, so this is functionally the same as
// calling IR.breakNode(). It exists so that a pass can be consistent about
// always using AstFactory to create new nodes.
//
// Java: AstFactory.createBreak()
func (a *AstFactory) CreateBreak() rhino.Node {
	return rhino.IR_BreakNode()
}

// CreateLabel returns a new LABEL node.
//
// Breaks have no type information, so this is functionally the same as
// calling IR.label(). It exists so that a pass can be consistent about
// always using AstFactory to create new nodes.
//
// Java: AstFactory.createLabel(Node, Node)
func (a *AstFactory) CreateLabel(label rhino.Node, stmt rhino.Node) rhino.Node {
	return rhino.IR_Label(label, stmt)
}

// CreateCatch returns a new CATCH node.
//
// CATCH have no type information, so this is functionally the same as
// calling IR.catchNode(). It exists so that a pass can be consistent about
// always using AstFactory to create new nodes.
//
// Java: AstFactory.createCatch(Node, Node)
func (a *AstFactory) CreateCatch(error_ rhino.Node, block rhino.Node) rhino.Node {
	return rhino.IR_CatchNode(error_, block)
}

// CreateTryFinally returns a new TRY-FINALLY node.
//
// TRY-FINALLY have no type information, so this is functionally the same as
// calling IR.tryFinally(). It exists so that a pass can be consistent about
// always using AstFactory to create new nodes.
//
// Java: AstFactory.createTryFinally(Node, Node)
func (a *AstFactory) CreateTryFinally(tryBlock rhino.Node, finallyBlock rhino.Node) rhino.Node {
	return rhino.IR_TryFinally(tryBlock, finallyBlock)
}

// CreateTryCatchFinally returns a new TRY-CATCH-FINALLY node.
//
// TRY-CATCH-FINALLY have no type information, so this is functionally the
// same as calling IR.tryCatchFinally(). It exists so that a pass can be
// consistent about always using AstFactory to create new nodes.
//
// Java: AstFactory.createTryCatchFinally(Node, Node, Node)
func (a *AstFactory) CreateTryCatchFinally(
	tryBlock rhino.Node, catchNode rhino.Node, finallyBlock rhino.Node) rhino.Node {
	jutil.CheckState(tryBlock.IsBlock())
	jutil.CheckState(catchNode.IsCatch())
	jutil.CheckState(finallyBlock.IsBlock())
	return rhino.IR_TryCatchFinally(tryBlock, catchNode, finallyBlock)
}

// CreateThrow returns a new THROW node.
//
// THROW have no type information, so this is functionally the same as
// calling IR.throwNode(). It exists so that a pass can be consistent about
// always using AstFactory to create new nodes.
//
// Java: AstFactory.createThrow(Node)
func (a *AstFactory) CreateThrow(expr rhino.Node) rhino.Node {
	return rhino.IR_ThrowNode(expr)
}

// CreateReturn returns a new return statement.
//
// Return statements have no type information, so this is functionally the
// same as calling IR.return(value). It exists so that a pass can be
// consistent about always using AstFactory to create new nodes.
//
// Java: AstFactory.createReturn(Node)
func (a *AstFactory) CreateReturn(value rhino.Node) rhino.Node {
	return rhino.IR_ReturnNodeWithExpr(value)
}

// CreateYield returns a new yield expression.
//
// type_ is the type we expect to get back after the yield; value is the
// value to yield.
//
// Java: AstFactory.createYield(AstFactory.Type, Node)
func (a *AstFactory) CreateYield(type_ AstFactory_Type, value rhino.Node) rhino.Node {
	result := rhino.IR_YieldNode(value)
	a.SetJSTypeOrColor(type_, result)
	return result
}

// CreateAwait returns a new await expression.
//
// type_ is the type we expect to get back after the await; value is the
// value to await.
//
// Java: AstFactory.createAwait(AstFactory.Type, Node)
func (a *AstFactory) CreateAwait(type_ AstFactory_Type, value rhino.Node) rhino.Node {
	result := rhino.IR_Await(value)
	a.SetJSTypeOrColor(type_, result)
	return result
}

// Java: AstFactory.createNot(Node)
func (a *AstFactory) CreateNot(child rhino.Node) rhino.Node {
	result := rhino.IR_Not(child)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_BOOLEAN_TYPE, colors.StandardColors_BOOLEAN), result)
	return result
}

// Java: AstFactory.createSingleNameDeclaration(Token, String, Node)
func (a *AstFactory) CreateSingleNameDeclaration(tokenType rhino.Token, name string, value rhino.Node) rhino.Node {
	switch tokenType {
	case rhino.Token_LET:
		return a.CreateSingleLetNameDeclarationWithValue(name, value)
	case rhino.Token_VAR:
		return a.CreateSingleVarNameDeclarationWithValue(name, value)
	case rhino.Token_CONST:
		return a.CreateSingleConstNameDeclaration(name, value)
	default:
		panic(jutil.NewUnsupportedOperationExceptionMsg("Unexpeted token type: " + tokenType.String()))
	}
}

// CreateSingleLetNameDeclaration creates a new `let` declaration for a
// single variable name with a void type and no JSDoc.
//
// e.g. `let variableName`
//
// Java: AstFactory.createSingleLetNameDeclaration(String)
func (a *AstFactory) CreateSingleLetNameDeclaration(variableName string) rhino.Node {
	return rhino.IR_Let(
		a.CreateName(variableName,
			AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_VOID_TYPE, colors.StandardColors_NULL_OR_VOID)))
}

// CreateSingleLetNameDeclarationWithValue creates a new `let` declaration
// statement for a single variable name.
//
// Takes the type for the variable name from the value node.
//
// e.g. `let variableName = value;`
//
// Java: AstFactory.createSingleLetNameDeclaration(String, Node)
func (a *AstFactory) CreateSingleLetNameDeclarationWithValue(variableName string, value rhino.Node) rhino.Node {
	return rhino.IR_LetWithValue(a.CreateName(variableName, AstFactory_Type_(value)), value)
}

// CreateSingleVarNameDeclaration creates a new `var` declaration statement
// for a single variable name with void type and no JSDoc.
//
// e.g. `var variableName`
//
// Java: AstFactory.createSingleVarNameDeclaration(String)
func (a *AstFactory) CreateSingleVarNameDeclaration(variableName string) rhino.Node {
	return rhino.IR_Var(
		a.CreateName(variableName,
			AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_VOID_TYPE, colors.StandardColors_NULL_OR_VOID)))
}

// CreateSingleVarNameDeclarationWithValue creates a new `var` declaration
// statement for a single variable name.
//
// Takes the type for the variable name from the value node.
//
// e.g. `var variableName = value;`
//
// Java: AstFactory.createSingleVarNameDeclaration(String, Node)
func (a *AstFactory) CreateSingleVarNameDeclarationWithValue(variableName string, value rhino.Node) rhino.Node {
	return rhino.IR_VarWithValue(a.CreateName(variableName, AstFactory_Type_(value)), value)
}

// CreateSingleConstNameDeclaration creates a new `const` declaration
// statement for a single variable name.
//
// Takes the type for the variable name from the value node.
//
// e.g. `const variableName = value;`
//
// Java: AstFactory.createSingleConstNameDeclaration(String, Node)
func (a *AstFactory) CreateSingleConstNameDeclaration(variableName string, value rhino.Node) rhino.Node {
	nameNode := a.CreateConstantName(variableName, AstFactory_Type_OfJSTypeAndColor(value.GetJSType(), value.GetColor()))
	return rhino.IR_ConstNode(nameNode, value)
}

// CreateSingleConstObjectPatternDeclaration creates a new `const`
// declaration statement for an object pattern.
//
// e.g. `const {Foo} = value;`
//
// Java: AstFactory.createSingleConstObjectPatternDeclaration(Node, Node)
func (a *AstFactory) CreateSingleConstObjectPatternDeclaration(objectPattern rhino.Node, value rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(objectPattern.IsObjectPattern(), "not an object pattern: %s", objectPattern)
	return rhino.IR_ConstNode(objectPattern, value)
}

// CreateArgumentsAliasDeclaration creates a statement declaring a const
// alias for "arguments".
//
// e.g. `const argsAlias = arguments;`
//
// Java: AstFactory.createArgumentsAliasDeclaration(String)
func (a *AstFactory) CreateArgumentsAliasDeclaration(aliasName string) rhino.Node {
	return a.CreateSingleConstNameDeclaration(aliasName, a.CreateArgumentsReference())
}

// Java: AstFactory.createDelProp(Node)
func (a *AstFactory) CreateDelProp(target rhino.Node) rhino.Node {
	result := rhino.IR_Delprop(target)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_BOOLEAN_TYPE, colors.StandardColors_BOOLEAN), result)
	return result
}

// Java: AstFactory.createStringKey(String, Node)
func (a *AstFactory) CreateStringKey(key string, value rhino.Node) rhino.Node {
	result := rhino.IR_StringKeyWithValue(key, value)
	a.SetJSTypeOrColor(AstFactory_Type_OfJSTypeAndColor(value.GetJSType(), value.GetColor()), result)
	return result
}

// Java: AstFactory.createComputedProperty(Node, Node)
func (a *AstFactory) CreateComputedProperty(key rhino.Node, value rhino.Node) rhino.Node {
	result := rhino.IR_ComputedProp(key, value)
	a.SetJSTypeOrColor(AstFactory_Type_OfJSTypeAndColor(value.GetJSType(), value.GetColor()), result)
	return result
}

// CreateGetterDef creates a getter definition to be inserted into either a
// class body or object literal.
//
// `get name() { return value; }`
//
// Java: AstFactory.createGetterDef(String, Node)
func (a *AstFactory) CreateGetterDef(name string, value rhino.Node) rhino.Node {
	returnType := value.GetJSType()
	// Name is stored on the GETTER_DEF node. The function has no name.
	functionNode :=
		a.CreateZeroArgFunction( /* name= */ "", rhino.IR_BlockOfNode(a.CreateReturn(value)), returnType)
	getterNode := rhino.Node_NewStringWithToken(rhino.Token_GETTER_DEF, name)
	getterNode.AddChildToFront(functionNode)
	return getterNode
}

// Java: AstFactory.createIn(Node, Node)
func (a *AstFactory) CreateIn(left rhino.Node, right rhino.Node) rhino.Node {
	result := rhino.IR_In(left, right)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_BOOLEAN_TYPE, colors.StandardColors_BOOLEAN), result)
	return result
}

// Java: AstFactory.createComma(Node, Node)
func (a *AstFactory) CreateComma(left rhino.Node, right rhino.Node) rhino.Node {
	result := rhino.IR_Comma(left, right)
	a.SetJSTypeOrColor(AstFactory_Type_OfJSTypeAndColor(right.GetJSType(), right.GetColor()), result)
	return result
}

// Java: AstFactory.createCommas(Node, Node, Node...)
func (a *AstFactory) CreateCommas(first rhino.Node, second rhino.Node, rest ...rhino.Node) rhino.Node {
	result := a.CreateComma(first, second)
	for _, next := range rest {
		result = a.CreateComma(result, next)
	}
	return result
}

// Java: AstFactory.createAnd(Node, Node)
func (a *AstFactory) CreateAnd(left rhino.Node, right rhino.Node) rhino.Node {
	result := rhino.IR_And(left, right)
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		leftType := jutil.CheckNotNullMsg1(left.GetJSType(), "%s", left)
		rightType := jutil.CheckNotNullMsg1(right.GetJSType(), "%s", right)
		result.SetJSType(a.registry.CreateUnionType(leftType, rightType))
	case AstFactory_TypeMode_COLOR:
		leftColor := jutil.CheckNotNullMsg1(left.GetColor(), "%s", left)
		rightColor := jutil.CheckNotNullMsg1(right.GetColor(), "%s", right)
		result.SetColor(colors.Color_CreateUnion2(leftColor, rightColor))
	case AstFactory_TypeMode_NONE:
	}
	return result
}

// Java: AstFactory.createOr(Node, Node)
func (a *AstFactory) CreateOr(left rhino.Node, right rhino.Node) rhino.Node {
	result := rhino.IR_Or(left, right)
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		leftType := jutil.CheckNotNullMsg1(left.GetJSType(), "%s", left)
		rightType := jutil.CheckNotNullMsg1(right.GetJSType(), "%s", right)
		result.SetJSType(a.registry.CreateUnionType(leftType, rightType))
	case AstFactory_TypeMode_COLOR:
		leftColor := jutil.CheckNotNullMsg1(left.GetColor(), "%s", left)
		rightColor := jutil.CheckNotNullMsg1(right.GetColor(), "%s", right)
		result.SetColor(colors.Color_CreateUnion2(leftColor, rightColor))
	case AstFactory_TypeMode_NONE:
	}
	return result
}

// Java: AstFactory.createAdd(Node, Node)
func (a *AstFactory) CreateAdd(left rhino.Node, right rhino.Node) rhino.Node {
	result := rhino.IR_Add(left, right)
	// Note: this result type could be made tighter if it proves useful for optimizations later on
	// like setting the string type if both operands are strings.
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		result.SetJSType(a.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER_STRING))
	case AstFactory_TypeMode_COLOR:
		result.SetColor(AstFactory_bigintNumberStringColor())
	case AstFactory_TypeMode_NONE:
	}
	return result
}

// Java: AstFactory.createSub(Node, Node)
func (a *AstFactory) CreateSub(left rhino.Node, right rhino.Node) rhino.Node {
	result := rhino.IR_Sub(left, right)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER), result)
	return result
}

// Java: AstFactory.createInc(Node, boolean)
func (a *AstFactory) CreateInc(operand rhino.Node, isPost bool) rhino.Node {
	result := rhino.IR_Inc(operand, isPost)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER), result)
	return result
}

// Java: AstFactory.createLessThan(Node, Node)
func (a *AstFactory) CreateLessThan(left rhino.Node, right rhino.Node) rhino.Node {
	result := rhino.IR_Lt(left, right)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_BOOLEAN_TYPE, colors.StandardColors_BOOLEAN), result)
	return result
}

// Java: AstFactory.createBitwiseAnd(Node, Node)
func (a *AstFactory) CreateBitwiseAnd(left rhino.Node, right rhino.Node) rhino.Node {
	result := rhino.IR_BitwiseAnd(left, right)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER), result)
	return result
}

// Java: AstFactory.createRightShift(Node, Node)
func (a *AstFactory) CreateRightShift(left rhino.Node, right rhino.Node) rhino.Node {
	result := rhino.IR_RightShift(left, right)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER), result)
	return result
}

// CreateAssignStatement creates a statement `lhs = rhs;`.
//
// Java: AstFactory.createAssignStatement(Node, Node)
func (a *AstFactory) CreateAssignStatement(lhs rhino.Node, rhs rhino.Node) rhino.Node {
	return a.ExprResult(a.CreateAssign(lhs, rhs))
}

// CreateAssign creates an assignment expression `lhs = rhs`
//
// Java: AstFactory.createAssign(Node, Node)
func (a *AstFactory) CreateAssign(lhs rhino.Node, rhs rhino.Node) rhino.Node {
	result := rhino.IR_Assign(lhs, rhs)
	a.SetJSTypeOrColor(AstFactory_Type_(rhs), result)
	return result
}

// CreateAssignToName creates an assignment expression `lhs = rhs`
//
// Java: AstFactory.createAssign(String, Node)
func (a *AstFactory) CreateAssignToName(lhsName string, rhs rhino.Node) rhino.Node {
	name := a.CreateName(lhsName, AstFactory_Type_(rhs))
	return a.CreateAssign(name, rhs)
}

// CreateObjectLit creates an object-literal with zero or more elements, `{}`.
//
// The type of the literal, if assigned, may be a supertype of the known
// properties.
//
// Java: AstFactory.createObjectLit(Node...)
func (a *AstFactory) CreateObjectLit(elements ...rhino.Node) rhino.Node {
	result := rhino.IR_Objectlit(elements...)
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		result.SetJSType(a.registry.CreateAnonymousObjectType(nil))
	case AstFactory_TypeMode_COLOR:
		result.SetColor(colors.StandardColors_TOP_OBJECT)
	case AstFactory_TypeMode_NONE:
	}
	return result
}

// CreateObjectLitWithType creates an object-literal with zero or more
// elements and a specific type.
//
// Java: AstFactory.createObjectLit(AstFactory.Type, Node...)
func (a *AstFactory) CreateObjectLitWithType(type_ AstFactory_Type, elements ...rhino.Node) rhino.Node {
	result := rhino.IR_Objectlit(elements...)
	a.SetJSTypeOrColor(type_, result)
	return result
}

// Java: AstFactory.createQuotedStringKey(String, Node)
func (a *AstFactory) CreateQuotedStringKey(key string, value rhino.Node) rhino.Node {
	result := rhino.IR_StringKeyWithValue(key, value)
	result.SetQuotedStringKey()
	return result
}

// CreateEmptyFunction creates an empty function `function() {}`
//
// Java: AstFactory.createEmptyFunction(AstFactory.Type)
func (a *AstFactory) CreateEmptyFunction(type_ AstFactory_Type) rhino.Node {
	result := NodeUtil_EmptyFunction()
	if a.IsAddingTypes() {
		jutil.CheckArgumentMsg1(type_.GetJSType(a.registry).IsFunctionType(), "%s", type_)
	}
	a.SetJSTypeOrColor(type_, result)
	return result
}

// CreateEmptyGeneratorFunction creates an empty function `function*() {}`
//
// Java: AstFactory.createEmptyGeneratorFunction(AstFactory.Type)
func (a *AstFactory) CreateEmptyGeneratorFunction(type_ AstFactory_Type) rhino.Node {
	result := a.CreateEmptyFunction(type_)
	result.SetIsGeneratorFunction(true)
	return result
}

// CreateFunction creates a function `function name(paramList) { body }`
//
// name is the name, empty string if no name; paramList a PARAM_LIST node;
// body a BLOCK node; type_ the type to apply to the function itself.
//
// Java: AstFactory.createFunction(String, Node, Node, AstFactory.Type)
func (a *AstFactory) CreateFunction(name string, paramList rhino.Node, body rhino.Node, type_ AstFactory_Type) rhino.Node {
	nameNode := a.CreateName(name, type_)
	result := rhino.IR_Function(nameNode, paramList, body)
	if a.IsAddingTypes() {
		jutil.CheckArgumentMsg1(type_.GetJSType(a.registry).IsFunctionType(), "%s", type_)
	}
	a.SetJSTypeOrColor(type_, result)
	return result
}

// Java: AstFactory.createParamList(String...)
func (a *AstFactory) CreateParamList(parameterNames ...string) rhino.Node {
	paramList := rhino.IR_ParamList()
	for _, parameterName := range parameterNames {
		paramList.AddChildToBack(a.CreateNameWithUnknownType(parameterName))
	}
	return paramList
}

// Java: AstFactory.createZeroArgFunction(String, Node, JSType)
func (a *AstFactory) CreateZeroArgFunction(name string, body rhino.Node, returnType rhino.JSType) rhino.Node {
	var functionType rhino.FunctionType
	if a.IsAddingTypes() {
		functionType = a.registry.CreateFunctionType(returnType).ToMaybeFunctionType()
	}
	return a.CreateFunction(
		name, rhino.IR_ParamList(), body,
		AstFactory_Type_OfJSTypeAndColor(astFactoryFunctionTypeAsJSType(functionType), colors.StandardColors_TOP_OBJECT))
}

// astFactoryFunctionTypeAsJSType widens a FunctionType to a JSType.
//
// Go-only: language-adapter — keeps a nil FunctionType a nil JSType.
func astFactoryFunctionTypeAsJSType(t rhino.FunctionType) rhino.JSType {
	if t == nil {
		return nil
	}
	return t
}

// Java: AstFactory.createZeroArgGeneratorFunction(String, Node, JSType)
func (a *AstFactory) CreateZeroArgGeneratorFunction(name string, body rhino.Node, returnType rhino.JSType) rhino.Node {
	result := a.CreateZeroArgFunction(name, body, returnType)
	result.SetIsGeneratorFunction(true)
	return result
}

// Java: AstFactory.createZeroArgArrowFunctionForExpression(Node)
func (a *AstFactory) CreateZeroArgArrowFunctionForExpression(expression rhino.Node) rhino.Node {
	result := rhino.IR_ArrowFunction(rhino.IR_Name(""), rhino.IR_ParamList(), expression)
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		// It feels like we should be adding type-of-this here, but it should remain unknown,
		// because you're allowed to supply any kind of value of `this` when calling an arrow
		// function. It will just be ignored in favor of the `this` in the scope where the
		// arrow was defined.
		var b rhino.FunctionType_Builder
		b.Init(a.registry)
		functionType :=
			b.WithReturnType(expression.GetJSTypeRequired()).
				WithParameters().
				BuildAndResolve()
		result.SetJSType(functionType)
	case AstFactory_TypeMode_COLOR:
		result.SetColor(colors.StandardColors_TOP_OBJECT)
	case AstFactory_TypeMode_NONE:
	}
	return result
}

// Java: AstFactory.createMemberFunctionDef(String, Node)
func (a *AstFactory) CreateMemberFunctionDef(name string, function rhino.Node) rhino.Node {
	// A function used for a member function definition must have an empty name,
	// because the name string goes on the MEMBER_FUNCTION_DEF node.
	jutil.CheckArgumentMsg1(function.GetFirstChild().GetString() == "", "%s", function)
	result := rhino.IR_MemberFunctionDef(name, function)
	a.SetJSTypeOrColor(AstFactory_Type_(function), result)
	return result
}

// Java: AstFactory.createSheq(Node, Node)
func (a *AstFactory) CreateSheq(expr1 rhino.Node, expr2 rhino.Node) rhino.Node {
	result := rhino.IR_Sheq(expr1, expr2)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_BOOLEAN_TYPE, colors.StandardColors_BOOLEAN), result)
	return result
}

// Java: AstFactory.createEq(Node, Node)
func (a *AstFactory) CreateEq(expr1 rhino.Node, expr2 rhino.Node) rhino.Node {
	result := rhino.IR_Eq(expr1, expr2)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_BOOLEAN_TYPE, colors.StandardColors_BOOLEAN), result)
	return result
}

// Java: AstFactory.createNe(Node, Node)
func (a *AstFactory) CreateNe(expr1 rhino.Node, expr2 rhino.Node) rhino.Node {
	result := rhino.IR_Ne(expr1, expr2)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_BOOLEAN_TYPE, colors.StandardColors_BOOLEAN), result)
	return result
}

// Java: AstFactory.createHook(Node, Node, Node)
func (a *AstFactory) CreateHook(condition rhino.Node, expr1 rhino.Node, expr2 rhino.Node) rhino.Node {
	result := rhino.IR_Hook(condition, expr1, expr2)
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		result.SetJSType(a.registry.CreateUnionType(expr1.GetJSType(), expr2.GetJSType()))
	case AstFactory_TypeMode_COLOR:
		result.SetColor(colors.Color_CreateUnion2(expr1.GetColor(), expr2.GetColor()))
	case AstFactory_TypeMode_NONE:
	}

	return result
}

// CreateOptionalChainShortCircuit creates a nullish short-circuit hook
// `receiver == null ? void 0 : resultExpr` used for optional chaining.
//
// Java: AstFactory.createOptionalChainShortCircuit(Node, Node)
func (a *AstFactory) CreateOptionalChainShortCircuit(receiver rhino.Node, resultExpr rhino.Node) rhino.Node {
	cond := a.CreateEq(receiver, a.CreateNull())
	return a.CreateHook(cond, a.CreateUndefinedValue(), resultExpr)
}

// Java: AstFactory.createArraylit(Node...)
func (a *AstFactory) CreateArraylit(elements ...rhino.Node) rhino.Node {
	return a.CreateArraylitOfIterable(slices.Values(elements))
}

// Java: AstFactory.createArraylit(Iterable<Node>)
func (a *AstFactory) CreateArraylitOfIterable(elements iter.Seq[rhino.Node]) rhino.Node {
	result := rhino.IR_ArraylitOfIterable(elements)
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		result.SetJSType(
			a.registry.CreateTemplatizedType(
				a.registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE),
				// TODO(nickreid): Use a reasonable template type. Remeber to consider SPREAD.
				a.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)))
	case AstFactory_TypeMode_COLOR:
		result.SetColor(a.colorRegistry.Get(colors.StandardColors_ARRAY_ID))
	case AstFactory_TypeMode_NONE:
	}
	return result
}

// CreateEmptyAsyncGeneratorWrapperArgument creates an empty generator
// function with the correct return type to be an argument to
// $jscomp.AsyncGeneratorWrapper.
//
// asyncGeneratorWrapperType is the specific type of the
// $jscomp.AsyncGeneratorWrapper with its template filled in. Should be the
// type on the node returned from createAsyncGeneratorWrapperReference.
//
// Java: AstFactory.createEmptyAsyncGeneratorWrapperArgument(JSType)
func (a *AstFactory) CreateEmptyAsyncGeneratorWrapperArgument(asyncGeneratorWrapperType rhino.JSType) rhino.Node {
	generatorType := AstFactory_NoTypeInformation()

	if a.IsAddingTypes() {
		if asyncGeneratorWrapperType.IsUnknownType() {
			// Not injecting libraries?
			generatorType =
				AstFactory_Type_OfJSType(
					a.registry.CreateFunctionType(
						a.ReplaceTemplate(
							a.GetNativeType(rhino.JSTypeNative_GENERATOR_TYPE),
							[]rhino.JSType{a.unknownType})))
		} else {
			// Generator<$jscomp.AsyncGeneratorWrapper$ActionRecord<number>>
			innerFunctionReturnType :=
				colorPool_getOnlyElement(
					asyncGeneratorWrapperType.ToMaybeFunctionType().GetParameters()).
					GetJSType()
			generatorType = AstFactory_Type_OfJSType(a.registry.CreateFunctionType(innerFunctionReturnType))
		}
	} else if a.IsAddingColors() {
		// colors don't model function types, so it's fine to fallback to the top object.
		generatorType = AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT)
	}

	return a.CreateEmptyGeneratorFunction(generatorType)
}
