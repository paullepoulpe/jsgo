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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: AstValidator.validateImport(Node)
func (v *AstValidator) ValidateImport(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_MODULES, n)
	v.ValidateNodeType(rhino.Token_IMPORT, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)

	if n.GetFirstChild().IsName() {
		v.ValidateName(n.GetFirstChild())
	} else {
		v.ValidateNodeType(rhino.Token_EMPTY, n.GetFirstChild())
	}

	secondChild := n.GetSecondChild()
	switch secondChild.GetToken() {
	case rhino.Token_IMPORT_SPECS:
		v.ValidateImportSpecifiers(secondChild)
	case rhino.Token_IMPORT_STAR:
		v.ValidateNonEmptyString(secondChild)
	default:
		v.ValidateNodeType(rhino.Token_EMPTY, secondChild)
	}

	v.ValidateStringLit(n.GetChildAtIndex(2))
}

// Java: AstValidator.validateImportSpecifiers(Node)
func (v *AstValidator) ValidateImportSpecifiers(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_IMPORT_SPECS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateImportSpecifier(c)
	}
}

// Java: AstValidator.validateImportSpecifier(Node)
func (v *AstValidator) ValidateImportSpecifier(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_IMPORT_SPEC, n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 2)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateName(c)
	}
}

// Java: AstValidator.validateExport(Node, boolean)
func (v *AstValidator) ValidateExport(n rhino.Node, isAmbient bool) {
	v.ValidateFeature(parser.FeatureSet_Feature_MODULES, n)
	v.ValidateNodeType(rhino.Token_EXPORT, n)
	if n.GetBooleanProp(rhino.Node_EXPORT_ALL_FROM) { // export * from "mod"
		v.ValidateProperties(n)
		v.ValidateChildCountWithExpected(n, 2)
		v.ValidateNodeType(rhino.Token_EMPTY, n.GetFirstChild())
		v.ValidateStringLit(n.GetSecondChild())
	} else if n.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) { // export default foo = 2
		v.ValidateProperties(n)
		v.ValidateChildCountWithExpected(n, 1)
		v.ValidateExpression(n.GetFirstChild())
	} else {
		v.ValidateProperties(n)
		v.ValidateChildCountIn(n, 1, 2)
		if n.GetFirstChild().IsExportSpecs() {
			v.ValidateExportSpecifiers(n.GetFirstChild())
		} else {
			v.ValidateStatementWithAmbient(n.GetFirstChild(), isAmbient)
		}
		if n.HasTwoChildren() {
			v.ValidateStringLit(n.GetSecondChild())
		}
	}
}

// Java: AstValidator.validateExportSpecifiers(Node)
func (v *AstValidator) ValidateExportSpecifiers(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_EXPORT_SPECS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateExportSpecifier(c)
	}
}

// Java: AstValidator.validateExportSpecifier(Node)
func (v *AstValidator) ValidateExportSpecifier(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_EXPORT_SPEC, n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 2)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateName(c)
	}
}

// Java: AstValidator.validateInterface(Node)
func (v *AstValidator) ValidateInterface(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_INTERFACE, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	name := n.GetFirstChild()
	v.ValidateName(name)
	superTypes := name.GetNext()
	if superTypes.IsEmpty() {
		v.ValidateProperties(superTypes)
		v.ValidateChildless(superTypes)
	} else {
		v.ValidateInterfaceExtends(superTypes)
	}
	v.ValidateInterfaceMembers(n.GetLastChild())
}

// Java: AstValidator.validateInterfaceExtends(Node)
func (v *AstValidator) ValidateInterfaceExtends(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_INTERFACE_EXTENDS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateNamedType(c)
	}
}

// Java: AstValidator.validateInterfaceMembers(Node)
func (v *AstValidator) ValidateInterfaceMembers(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_INTERFACE_MEMBERS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateInterfaceMember(c)
	}
}

// Java: AstValidator.validateInterfaceMember(Node)
func (v *AstValidator) ValidateInterfaceMember(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_MEMBER_FUNCTION_DEF:
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		v.ValidateFunctionSignature(n.GetFirstChild())
	case rhino.Token_MEMBER_VARIABLE_DEF:
		v.ValidateProperties(n)
		v.ValidateChildless(n)
	case rhino.Token_INDEX_SIGNATURE:
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		child := n.GetFirstChild()
		v.ValidateProperties(child)
		v.ValidateChildless(child)
	case rhino.Token_CALL_SIGNATURE:
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
	default:
		v.Violation("Interface contained member of invalid type "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateEnum(Node)
func (v *AstValidator) ValidateEnum(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_ENUM, n)
	v.ValidateName(n.GetFirstChild())
	v.ValidateEnumMembers(n.GetLastChild())
}

// Java: AstValidator.validateEnumMembers(Node)
func (v *AstValidator) ValidateEnumMembers(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_ENUM_MEMBERS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateEnumStringKey(c)
	}
}

// Java: AstValidator.validateEnumStringKey(Node)
func (v *AstValidator) ValidateEnumStringKey(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_STRING_KEY, n)
	v.ValidateObjectLiteralKeyName(n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 0)
}

// ValidateClassDeclaration: in a class declaration, unlike a class
// expression, the class name is required.
//
// Java: AstValidator.validateClassDeclaration(Node, boolean)
func (v *AstValidator) ValidateClassDeclaration(n rhino.Node, isAmbient bool) {
	v.ValidateClassHelper(n, isAmbient)
	v.ValidateName(n.GetFirstChild())
}

// Java: AstValidator.validateClass(Node)
func (v *AstValidator) ValidateClass(n rhino.Node) {
	v.ValidateClassHelper(n, false)
}

// Java: AstValidator.validateClassHelper(Node, boolean)
func (v *AstValidator) ValidateClassHelper(n rhino.Node, isAmbient bool) {
	v.ValidateFeature(parser.FeatureSet_Feature_CLASSES, n)
	v.ValidateNodeType(rhino.Token_CLASS, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)

	name := n.GetFirstChild()
	if name.IsEmpty() {
		v.ValidateProperties(name)
		v.ValidateChildless(name)
	} else {
		v.ValidateName(name)
	}

	superClass := name.GetNext()
	if superClass.IsEmpty() {
		v.ValidateProperties(superClass)
		v.ValidateChildless(superClass)
	} else {
		v.ValidateExpression(superClass)
	}

	v.ValidateClassMembers(n.GetLastChild(), isAmbient)
}

// Java: AstValidator.validateClassMembers(Node, boolean)
func (v *AstValidator) ValidateClassMembers(n rhino.Node, isAmbient bool) {
	v.ValidateNodeType(rhino.Token_CLASS_MEMBERS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateClassMember(c, isAmbient)
	}
}

// Java: AstValidator.validateClassMember(Node, boolean)
func (v *AstValidator) ValidateClassMember(n rhino.Node, isAmbient bool) {
	switch n.GetToken() {
	case rhino.Token_MEMBER_FUNCTION_DEF:
		v.ValidateFeature(parser.FeatureSet_Feature_MEMBER_DECLARATIONS, n)
		if n.IsPrivateIdentifier() {
			v.ValidateFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS, n)
		}
		v.ValidateObjectLiteralKeyName(n)
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		v.ValidateMemberFunction(n, isAmbient)
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		v.ValidateFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER, n)
		if n.IsPrivateIdentifier() {
			v.ValidateFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS, n)
		}
		v.ValidateObjectLiteralKeyName(n)
		v.ValidateObjectLitKey(n)
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		v.ValidateMemberFunction(n, isAmbient)
	case rhino.Token_MEMBER_VARIABLE_DEF:
		v.ValidateProperties(n)
		v.ValidateChildless(n)
	case rhino.Token_COMPUTED_PROP:
		v.ValidateComputedPropClassMethod(n)
	case rhino.Token_MEMBER_FIELD_DEF:
		v.ValidateClassField(n)
	case rhino.Token_COMPUTED_FIELD_DEF:
		v.ValidateComputedPropClassField(n)
	case rhino.Token_INDEX_SIGNATURE:
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		child := n.GetFirstChild()
		v.ValidateProperties(child)
		v.ValidateChildless(child)
	case rhino.Token_CALL_SIGNATURE:
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
	case rhino.Token_BLOCK:
		v.ValidateFeature(parser.FeatureSet_Feature_CLASS_STATIC_BLOCK, n)
		v.ValidateBlock(n)
	case rhino.Token_EMPTY:
		// Empty is allowed too.
	default:
		v.Violation("Class contained member of invalid type "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateMemberFunction(Node, boolean)
func (v *AstValidator) ValidateMemberFunction(n rhino.Node, isAmbient bool) {
	function := n.GetFirstChild()
	if isAmbient {
		v.ValidateFunctionSignature(function)
	} else {
		v.ValidateFunctionExpression(function)
	}
}

// Java: AstValidator.validateClassField(Node)
func (v *AstValidator) ValidateClassField(n rhino.Node) {
	if n.IsPrivateIdentifier() {
		v.ValidateFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS, n)
	} else {
		v.ValidateFeature(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS, n)
	}
	v.ValidateNonEmptyString(n)
	if n.HasChildren() {
		v.ValidateExpression(n.GetFirstChild())
	}
}

// Java: AstValidator.validateComputedPropClassField(Node)
func (v *AstValidator) ValidateComputedPropClassField(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS, n)
	v.ValidateExpression(n.GetFirstChild())
	if !n.GetSecondChild().IsNil() {
		v.ValidateExpression(n.GetSecondChild())
	}
}

// Java: AstValidator.validateComputedPropClassMethod(Node)
func (v *AstValidator) ValidateComputedPropClassMethod(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.ValidateNodeType(rhino.Token_COMPUTED_PROP, n)
	v.ValidateExpression(n.GetFirstChild())
	if n.GetBooleanProp(rhino.Node_COMPUTED_PROP_VARIABLE) {
		v.ValidateProperties(n)
		v.ValidateChildCountWithExpected(n, 1)
	} else {
		v.ValidateProperties(n)
		v.ValidateChildCountWithExpected(n, 2)
		v.ValidateFunctionExpression(n.GetLastChild())
		if n.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER) {
			v.ValidateObjectLitComputedPropGetKey(n)
			v.ValidateFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER, n)
		} else if n.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER) {
			v.ValidateObjectLitComputedPropSetKey(n)
			v.ValidateFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER, n)
		}
	}
}

// Java: AstValidator.validateFunctionStatement(Node)
func (v *AstValidator) ValidateFunctionStatement(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_FUNCTION, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateName(n.GetFirstChild())
	v.ValidateParameters(n.GetSecondChild())
	v.ValidateFunctionBody(n.GetLastChild(), false)
	v.ValidateFunctionFeatures(n)
	if n.GetParent().IsBlock() && !n.GetGrandparent().IsFunction() {
		// e.g. if (true) { function f() {} }
		v.ValidateFeature(parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION, n)
	}
}

// Java: AstValidator.validateFunctionExpression(Node)
func (v *AstValidator) ValidateFunctionExpression(n rhino.Node) {
	v.ValidateFunctionExpressionHelper(n, false)
}

// Java: AstValidator.validateFunctionSignature(Node)
func (v *AstValidator) ValidateFunctionSignature(n rhino.Node) {
	v.ValidateFunctionExpressionHelper(n, true)
}

// Java: AstValidator.validateFunctionExpressionHelper(Node, boolean)
func (v *AstValidator) ValidateFunctionExpressionHelper(n rhino.Node, isAmbient bool) {
	v.ValidateNodeType(rhino.Token_FUNCTION, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)

	v.ValidateParameters(n.GetSecondChild())

	name := n.GetFirstChild()
	body := n.GetLastChild()
	if n.IsArrowFunction() {
		v.ValidateEmptyName(name)
		if body.IsBlock() {
			v.ValidateBlock(body)
		} else {
			v.ValidateExpression(body)
		}
	} else {
		v.ValidateOptionalName(name)
		v.ValidateFunctionBody(body, isAmbient)
	}
	v.ValidateFunctionFeatures(n)
}

// Java: AstValidator.validateFunctionFeatures(Node)
func (v *AstValidator) ValidateFunctionFeatures(n rhino.Node) {
	if n.IsArrowFunction() {
		v.ValidateFeature(parser.FeatureSet_Feature_ARROW_FUNCTIONS, n)
	}
	if n.IsGeneratorFunction() {
		v.ValidateFeature(parser.FeatureSet_Feature_GENERATORS, n)
	}
	if n.IsAsyncFunction() {
		v.ValidateFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS, n)
	}
	if n.IsAsyncFunction() && n.IsGeneratorFunction() {
		v.ValidateFeature(parser.FeatureSet_Feature_ASYNC_GENERATORS, n)
	}
}

// Java: AstValidator.validateFunctionBody(Node, boolean)
func (v *AstValidator) ValidateFunctionBody(n rhino.Node, noBlock bool) {
	if noBlock {
		v.ValidateNodeType(rhino.Token_EMPTY, n)
	} else {
		v.ValidateBlock(n)
	}
}

// Java: AstValidator.validateParameters(Node)
func (v *AstValidator) ValidateParameters(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_PARAM_LIST, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.IsRest() {
			v.ValidateRestParameters(rhino.Token_PARAM_LIST, c)
		} else if c.IsDefaultValue() {
			v.ValidateFeature(parser.FeatureSet_Feature_DEFAULT_PARAMETERS, c)
			v.ValidateDefaultValue(rhino.Token_PARAM_LIST, c)
		} else {
			if c.IsName() {
				v.ValidateName(c)
			} else if c.IsArrayPattern() {
				v.ValidateArrayPattern(rhino.Token_PARAM_LIST, c)
			} else {
				v.ValidateObjectPattern(rhino.Token_PARAM_LIST, c)
			}
		}
	}
}

// Java: AstValidator.validateDefaultValue(Token, Node)
func (v *AstValidator) ValidateDefaultValue(contextType rhino.Token, n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateLHS(contextType, n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateRestParameters(Token, Node)
func (v *AstValidator) ValidateRestParameters(contextType rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_REST_PARAMETERS, n)
	v.ValidateRest(contextType, n)
}

// Java: AstValidator.validateArrayPatternRest(Token, Node)
func (v *AstValidator) ValidateArrayPatternRest(contextType rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_ARRAY_PATTERN_REST, n)
	v.ValidateRest(contextType, n)
}

// Java: AstValidator.validateObjectPatternRest(Token, Node)
func (v *AstValidator) ValidateObjectPatternRest(contextType rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_OBJECT_PATTERN_REST, n)
	v.ValidateRest(contextType, n)
}

// ValidateRest: contextType is a Token constant value indicating that n
// should be validated appropriately for a descendant of a Node of this
// type.
//
// Java: AstValidator.validateRest(Token, Node)
func (v *AstValidator) ValidateRest(contextType rhino.Token, n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
	default:
		v.Violation("Unexpected node type.", n)
		return
	}
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateLHS(contextType, n.GetFirstChild())
	if !n.GetNext().IsNil() {
		v.Violation("Rest parameters must come after all other parameters.", n)
	}
}

// ValidateNameDeclarationHelper: statement is the enclosing statement. Will
// not always match the declaration Token.
//
// Java: AstValidator.validateNameDeclarationHelper(Node, Token, Node)
func (v *AstValidator) ValidateNameDeclarationHelper(statement rhino.Node, declaration rhino.Token, n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateMinimumChildCount(n, 1)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateNameDeclarationChild(statement, declaration, c)
	}
	if declaration == rhino.Token_LET {
		v.ValidateFeature(parser.FeatureSet_Feature_LET_DECLARATIONS, n)
	} else if declaration == rhino.Token_CONST {
		v.ValidateFeature(parser.FeatureSet_Feature_CONST_DECLARATIONS, n)
	}
}

// Java: AstValidator.validateNameDeclarationChild(Node, Token, Node)
func (v *AstValidator) ValidateNameDeclarationChild(statement rhino.Node, declaration rhino.Token, n rhino.Node) {
	inEnhancedFor := NodeUtil_IsEnhancedFor(statement)
	inForIn := statement.IsForIn()
	var minValues int
	var maxValues int
	if inForIn && declaration == rhino.Token_VAR {
		// ECMASCRIPT5 sloppy mode allows for-in initializers.
		minValues = 0
		maxValues = 1
	} else if inEnhancedFor {
		minValues = 0
		maxValues = 0
	} else if n.IsDestructuringLhs() || declaration == rhino.Token_CONST {
		minValues = 1
		maxValues = 1
	} else {
		minValues = 0
		maxValues = 1
	}

	if n.IsName() {
		// Don't use validateName here since this NAME node may have a child.
		v.ValidateNonEmptyString(n)
		v.ValidateProperties(n)
		v.ValidateChildCountIn(n, minValues, maxValues)

		if n.HasChildren() {
			v.ValidateExpression(n.GetFirstChild())
		}
	} else if n.IsDestructuringLhs() {
		v.ValidateProperties(n)
		v.ValidateChildCountIn(n, 1+minValues, 1+maxValues)

		c := n.GetFirstChild()
		switch c.GetToken() {
		case rhino.Token_ARRAY_PATTERN:
			v.ValidateArrayPattern(declaration, c)
		case rhino.Token_OBJECT_PATTERN:
			v.ValidateObjectPattern(declaration, c)
		default:
			v.Violation("Invalid destructuring lhs first child for "+declaration.String()+" node", n)
		}

		if n.HasTwoChildren() {
			v.ValidateExpression(n.GetSecondChild())
		}
	} else {
		v.Violation("Invalid child for "+declaration.String()+" node", n)
	}
}

// ValidateLHS: contextType is a Token constant value indicating that n
// should be validated appropriately for a descendant of a Node of this
// type.
//
// Java: AstValidator.validateLHS(Token, Node)
func (v *AstValidator) ValidateLHS(contextType rhino.Token, n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		v.ValidateName(n)
	case rhino.Token_ARRAY_PATTERN:
		v.ValidateArrayPattern(contextType, n)
	case rhino.Token_OBJECT_PATTERN:
		v.ValidateObjectPattern(contextType, n)
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		v.ValidateGetPropGetElemInLHS(contextType, n)
	case rhino.Token_CAST:
		v.ValidateLHS(contextType, n.GetOnlyChild())
	default:
		v.Violation("Invalid child for "+contextType.String()+" node", n)
	}
}

// Java: AstValidator.validateArrayPattern(Token, Node)
func (v *AstValidator) ValidateArrayPattern(type_ rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, n)
	v.ValidateNodeType(rhino.Token_ARRAY_PATTERN, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_DEFAULT_VALUE:
			v.ValidateDefaultValue(type_, c)
		case rhino.Token_ITER_REST:
			v.ValidateArrayPatternRest(type_, c)
		case rhino.Token_EMPTY:
			v.ValidateProperties(c)
			v.ValidateChildless(c)
		default:
			v.ValidateLHS(type_, c)
		}
	}
}

// Java: AstValidator.validateObjectPattern(Token, Node)
func (v *AstValidator) ValidateObjectPattern(type_ rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, n)
	v.ValidateNodeType(rhino.Token_OBJECT_PATTERN, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_STRING_KEY:
			v.ValidateObjectPatternStringKey(type_, c)
		case rhino.Token_OBJECT_REST:
			v.ValidateObjectPatternRest(type_, c)
		case rhino.Token_COMPUTED_PROP:
			v.ValidateObjectPatternComputedPropKey(type_, c)
		default:
			v.Violation("Invalid object pattern child for "+type_.String()+" node", n)
		}
	}
}

// Java: AstValidator.validateObjectPatternStringKey(Token, Node)
func (v *AstValidator) ValidateObjectPatternStringKey(type_ rhino.Token, n rhino.Node) {
	v.ValidateNodeType(rhino.Token_STRING_KEY, n)
	v.ValidateObjectLiteralKeyName(n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 1)

	c := n.GetFirstChild()
	switch c.GetToken() {
	case rhino.Token_DEFAULT_VALUE:
		v.ValidateDefaultValue(type_, c)
	default:
		v.ValidateLHS(type_, c)
	}
}

// Java: AstValidator.validateObjectPatternComputedPropKey(Token, Node)
func (v *AstValidator) ValidateObjectPatternComputedPropKey(type_ rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.ValidateNodeType(rhino.Token_COMPUTED_PROP, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
	if n.GetLastChild().IsDefaultValue() {
		v.ValidateDefaultValue(type_, n.GetLastChild())
	} else {
		v.ValidateLHS(n.GetLastChild().GetToken(), n.GetLastChild())
	}
}

// Java: AstValidator.validateNamedType(Node)
func (v *AstValidator) ValidateNamedType(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_NAMED_TYPE, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateName(n.GetFirstChild())
}

// Java: AstValidator.validateTypeAlias(Node)
func (v *AstValidator) ValidateTypeAlias(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_TYPE_ALIAS, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
}

// Java: AstValidator.validateAmbientDeclaration(Node)
func (v *AstValidator) ValidateAmbientDeclaration(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_DECLARE, n)
	v.ValidateAmbientDeclarationHelper(n.GetFirstChild())
}

// Java: AstValidator.validateAmbientDeclarationHelper(Node)
func (v *AstValidator) ValidateAmbientDeclarationHelper(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		v.ValidateNameDeclarationHelper(n.GetParent(), n.GetToken(), n)
	case rhino.Token_FUNCTION:
		v.ValidateFunctionSignature(n)
	case rhino.Token_CLASS:
		v.ValidateClassDeclaration(n, true)
	case rhino.Token_ENUM:
		v.ValidateEnum(n)
	case rhino.Token_NAMESPACE:
		v.ValidateNamespace(n, true)
	case rhino.Token_TYPE_ALIAS:
		v.ValidateTypeAlias(n)
	case rhino.Token_EXPORT:
		v.ValidateExport(n, true)
	default:
	}
}

// Java: AstValidator.validateNamespace(Node, boolean)
func (v *AstValidator) ValidateNamespace(n rhino.Node, isAmbient bool) {
	v.ValidateNodeType(rhino.Token_NAMESPACE, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateNamespaceName(n.GetFirstChild())
	v.ValidateNamespaceElements(n.GetLastChild(), isAmbient)
}

// Java: AstValidator.validateNamespaceName(Node)
func (v *AstValidator) ValidateNamespaceName(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		v.ValidateName(n)
	case rhino.Token_GETPROP:
		v.ValidateGetProp(n)
	default:
	}
}

// Java: AstValidator.validateNamespaceElements(Node, boolean)
func (v *AstValidator) ValidateNamespaceElements(n rhino.Node, isAmbient bool) {
	v.ValidateNodeType(rhino.Token_NAMESPACE_ELEMENTS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if isAmbient {
			v.ValidateAmbientDeclarationHelper(c)
		} else {
			v.ValidateStatement(c)
		}
	}
}
