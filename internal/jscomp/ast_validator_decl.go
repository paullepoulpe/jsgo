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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: AstValidator.validateImport(Node)
func (v *AstValidator) validateImport(imp ast.Import) {
	n := imp.Node()
	v.validateFeature(parser.FeatureSet_Feature_MODULES, n)
	valid := v.validatePropsAndArity(n)

	if def, hasDefault := imp.Default(); hasDefault {
		if def.Node().IsName() {
			v.validateName(def.Node())
		} else {
			v.validateNodeType(rhino.Token_EMPTY, def.Node())
		}
	}

	if bindings, hasBindings := imp.Bindings(); hasBindings {
		if specs, ok := bindings.AsImportSpecs(); ok {
			v.validateImportSpecifiers(specs)
		} else if star, ok := bindings.AsImportStar(); ok {
			v.validateNonEmptyString(star.Node())
		} else {
			v.validateNodeType(rhino.Token_EMPTY, bindings.Node())
		}
	}

	if valid {
		v.validateStringLit(imp.Source().Node())
	} else {
		v.validateStringLit(n.GetChildAtIndex(2))
	}
}

// Java: AstValidator.validateImportSpecifiers(Node)
func (v *AstValidator) validateImportSpecifiers(specs ast.ImportSpecs) {
	for spec := range specs.Specs() {
		v.validateImportSpecifier(spec.Node())
	}
}

// Java: AstValidator.validateImportSpecifier(Node)
func (v *AstValidator) validateImportSpecifier(n rhino.Node) {
	isSpec := v.validateNodeType(rhino.Token_IMPORT_SPEC, n)
	if v.validatePropsAndCount(n, 2) && isSpec {
		spec := ast.MustImportSpec(n)
		v.validateName(spec.Imported().Node())
		v.validateName(spec.Local().Node())
		return
	}
	for c := range n.Children() {
		v.validateName(c)
	}
}

// Java: AstValidator.validateExport(Node, boolean)
func (v *AstValidator) validateExport(exp ast.Export, isAmbient bool) {
	n := exp.Node()
	v.validateFeature(parser.FeatureSet_Feature_MODULES, n)
	switch {
	case exp.IsAllFrom(): // export * from "mod"
		v.validatePropsAndCount(n, 2)
		v.validateNodeType(rhino.Token_EMPTY, exp.Body().Node())
		source, _ := exp.Source()
		v.validateStringLit(source.Node())
	case exp.IsDefault(): // export default foo = 2
		v.validatePropsAndCount(n, 1)
		v.ValidateExpression(exp.Body().Node())
	default:
		v.validateProperties(n)
		valid := v.validateChildCountIn(n, 1, 2)
		body := exp.Body()
		if specs, ok := body.AsExportSpecs(); ok {
			v.validateExportSpecifiers(specs)
		} else {
			v.validateStatementWithAmbient(body.Node(), isAmbient)
		}
		if valid {
			if source, ok := exp.Source(); ok {
				v.validateStringLit(source.Node())
			}
		}
	}
}

// Java: AstValidator.validateExportSpecifiers(Node)
func (v *AstValidator) validateExportSpecifiers(specs ast.ExportSpecs) {
	for spec := range specs.Specs() {
		v.validateExportSpecifier(spec.Node())
	}
}

// Java: AstValidator.validateExportSpecifier(Node)
func (v *AstValidator) validateExportSpecifier(n rhino.Node) {
	isSpec := v.validateNodeType(rhino.Token_EXPORT_SPEC, n)
	if v.validatePropsAndCount(n, 2) && isSpec {
		spec := ast.MustExportSpec(n)
		v.validateName(spec.Local().Node())
		v.validateName(spec.Exported().Node())
		return
	}
	for c := range n.Children() {
		v.validateName(c)
	}
}

// Java: AstValidator.validateInterface(Node)
func (v *AstValidator) validateInterface(n rhino.Node) {
	v.validateNodeType(rhino.Token_INTERFACE, n)
	v.validatePropsAndArity(n)
	name := n.GetFirstChild()
	v.validateName(name)
	superTypes := name.GetNext()
	if superTypes.IsEmpty() {
		v.validateChildlessNode(superTypes)
	} else {
		v.validateInterfaceExtends(superTypes)
	}
	v.validateInterfaceMembers(n.GetLastChild())
}

// Java: AstValidator.validateInterfaceExtends(Node)
func (v *AstValidator) validateInterfaceExtends(n rhino.Node) {
	v.validateNodeType(rhino.Token_INTERFACE_EXTENDS, n)
	for c := range n.Children() {
		v.validateNamedType(c)
	}
}

// Java: AstValidator.validateInterfaceMembers(Node)
func (v *AstValidator) validateInterfaceMembers(n rhino.Node) {
	v.validateNodeType(rhino.Token_INTERFACE_MEMBERS, n)
	for c := range n.Children() {
		v.validateInterfaceMember(c)
	}
}

// Java: AstValidator.validateInterfaceMember(Node)
func (v *AstValidator) validateInterfaceMember(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_MEMBER_FUNCTION_DEF:
		v.validatePropsAndArity(n)
		v.validateFunctionSignature(ast.MustMemberFunctionDef(n).Function().Node())
	case rhino.Token_MEMBER_VARIABLE_DEF:
		v.validateChildlessNode(n)
	case rhino.Token_INDEX_SIGNATURE:
		v.validatePropsAndArity(n)
		v.validateChildlessNode(n.GetFirstChild())
	case rhino.Token_CALL_SIGNATURE:
		v.validatePropsAndArity(n)
	default:
		v.Violation("Interface contained member of invalid type "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateEnum(Node)
func (v *AstValidator) validateEnum(n rhino.Node) {
	v.validateNodeType(rhino.Token_ENUM, n)
	v.validateName(n.GetFirstChild())
	v.validateEnumMembers(n.GetLastChild())
}

// Java: AstValidator.validateEnumMembers(Node)
func (v *AstValidator) validateEnumMembers(n rhino.Node) {
	v.validateNodeType(rhino.Token_ENUM_MEMBERS, n)
	for c := range n.Children() {
		v.validateEnumStringKey(c)
	}
}

// Java: AstValidator.validateEnumStringKey(Node)
func (v *AstValidator) validateEnumStringKey(n rhino.Node) {
	v.validateNodeType(rhino.Token_STRING_KEY, n)
	v.validateObjectLiteralKeyName(n)
	v.validateChildlessNode(n)
}

// validateClassDeclaration validates a class declaration, where unlike a class
// expression the class name is required.
//
// Java: AstValidator.validateClassDeclaration(Node, boolean)
func (v *AstValidator) validateClassDeclaration(cls ast.Class, isAmbient bool) {
	v.validateClassHelper(cls, isAmbient)
	v.validateName(cls.NameNode())
}

// Java: AstValidator.validateClass(Node)
func (v *AstValidator) validateClass(cls ast.Class) {
	v.validateClassHelper(cls, false)
}

// Java: AstValidator.validateClassHelper(Node, boolean)
func (v *AstValidator) validateClassHelper(cls ast.Class, isAmbient bool) {
	n := cls.Node()
	v.validateFeature(parser.FeatureSet_Feature_CLASSES, n)
	v.validatePropsAndArity(n)

	if name, ok := cls.Name(); !ok {
		v.validateChildlessNode(cls.NameNode())
	} else {
		v.validateName(name.Node())
	}

	if superClass, ok := cls.Extends(); !ok {
		v.validateChildlessNode(cls.ExtendsNode())
	} else {
		v.ValidateExpression(superClass.Node())
	}

	v.validateClassMembers(cls.Body().Node(), isAmbient)
}

// Java: AstValidator.validateClassMembers(Node, boolean)
func (v *AstValidator) validateClassMembers(n rhino.Node, isAmbient bool) {
	if v.validateNodeType(rhino.Token_CLASS_MEMBERS, n) {
		for m := range ast.MustClassBody(n).Members() {
			v.validateClassMember(m.Node(), isAmbient)
		}
		return
	}
	for c := range n.Children() {
		v.validateClassMember(c, isAmbient)
	}
}

// Java: AstValidator.validateClassMember(Node, boolean)
func (v *AstValidator) validateClassMember(n rhino.Node, isAmbient bool) {
	switch n.GetToken() {
	case rhino.Token_MEMBER_FUNCTION_DEF:
		v.validateFeature(parser.FeatureSet_Feature_MEMBER_DECLARATIONS, n)
		if n.IsPrivateIdentifier() {
			v.validateFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS, n)
		}
		v.validateObjectLiteralKeyName(n)
		v.validatePropsAndArity(n)
		v.validateMemberFunction(ast.MustMethodDef(n), isAmbient)
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		v.validateFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER, n)
		if n.IsPrivateIdentifier() {
			v.validateFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS, n)
		}
		v.validateObjectLiteralKeyName(n)
		v.validateObjectLitKey(n)
		v.validatePropsAndArity(n)
		v.validateMemberFunction(ast.MustMethodDef(n), isAmbient)
	case rhino.Token_MEMBER_VARIABLE_DEF:
		v.validateChildlessNode(n)
	case rhino.Token_COMPUTED_PROP:
		v.validateComputedPropClassMethod(ast.MustComputedProp(n))
	case rhino.Token_MEMBER_FIELD_DEF:
		v.validateClassField(ast.MustMemberFieldDef(n))
	case rhino.Token_COMPUTED_FIELD_DEF:
		v.validateComputedPropClassField(ast.MustComputedFieldDef(n))
	case rhino.Token_INDEX_SIGNATURE:
		v.validatePropsAndArity(n)
		v.validateChildlessNode(n.GetFirstChild())
	case rhino.Token_CALL_SIGNATURE:
		v.validatePropsAndArity(n)
	case rhino.Token_BLOCK:
		v.validateFeature(parser.FeatureSet_Feature_CLASS_STATIC_BLOCK, n)
		v.validateBlock(n)
	case rhino.Token_EMPTY:
		// Empty semicolons in class bodies are permitted.
	default:
		v.Violation("Class contained member of invalid type "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateMemberFunction(Node, boolean)
func (v *AstValidator) validateMemberFunction(def ast.MethodDef, isAmbient bool) {
	fnNode := def.Function().Node()
	if isAmbient {
		v.validateFunctionSignature(fnNode)
	} else {
		v.validateFunctionExpression(fnNode)
	}
}

// Java: AstValidator.validateClassField(Node)
func (v *AstValidator) validateClassField(field ast.MemberFieldDef) {
	n := field.Node()
	if n.IsPrivateIdentifier() {
		v.validateFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS, n)
	} else {
		v.validateFeature(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS, n)
	}
	v.validateNonEmptyString(n)
	if val, ok := field.Value(); ok {
		v.ValidateExpression(val.Node())
	}
}

// Java: AstValidator.validateComputedPropClassField(Node)
func (v *AstValidator) validateComputedPropClassField(field ast.ComputedFieldDef) {
	v.validateFeature(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS, field.Node())
	v.ValidateExpression(field.Key().Node())
	if val, ok := field.Value(); ok {
		v.ValidateExpression(val.Node())
	}
}

// Java: AstValidator.validateComputedPropClassMethod(Node)
func (v *AstValidator) validateComputedPropClassMethod(prop ast.ComputedProp) {
	n := prop.Node()
	v.validateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.ValidateExpression(prop.Key().Node())
	if n.GetBooleanProp(rhino.Node_COMPUTED_PROP_VARIABLE) {
		v.validatePropsAndCount(n, 1)
		return
	}
	var fnNode rhino.Node
	if v.validatePropsAndCount(n, 2) {
		val, _ := prop.Value()
		fnNode = val.Node()
	} else {
		fnNode = n.GetLastChild()
	}
	v.validateFunctionExpression(fnNode)
	if prop.IsGetter() {
		v.validateObjectLitComputedPropGetKey(prop)
		v.validateFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER, n)
	} else if prop.IsSetter() {
		v.validateObjectLitComputedPropSetKey(prop)
		v.validateFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER, n)
	}
}

// Java: AstValidator.validateFunctionStatement(Node)
func (v *AstValidator) validateFunctionStatement(fn ast.Function) {
	n := fn.Node()
	v.validatePropsAndArity(n)
	v.validateName(fn.NameNode().Node())
	v.validateParameters(fn.Params().Node())
	v.validateFunctionBody(fn.Body().Node(), false)
	v.validateFunctionFeatures(fn)
	if n.GetParent().IsBlock() && !n.GetGrandparent().IsFunction() {
		// e.g. if (true) { function f() {} }
		v.validateFeature(parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION, n)
	}
}

// Java: AstValidator.validateFunctionExpression(Node)
func (v *AstValidator) validateFunctionExpression(n rhino.Node) {
	v.validateFunctionExpressionHelper(n, false)
}

// Java: AstValidator.validateFunctionSignature(Node)
func (v *AstValidator) validateFunctionSignature(n rhino.Node) {
	v.validateFunctionExpressionHelper(n, true)
}

// Java: AstValidator.validateFunctionExpressionHelper(Node, boolean)
func (v *AstValidator) validateFunctionExpressionHelper(n rhino.Node, isAmbient bool) {
	isFn := v.validateNodeType(rhino.Token_FUNCTION, n)
	v.validatePropsAndArity(n)

	var name, params, body rhino.Node
	var isArrow bool
	if isFn {
		fn := ast.MustFunction(n)
		params = fn.Params().Node()
		name = fn.NameNode().Node()
		body = fn.Body().Node()
		isArrow = fn.IsArrow()
	} else {
		params = n.GetSecondChild()
		name = n.GetFirstChild()
		body = n.GetLastChild()
		isArrow = n.IsArrowFunction()
	}

	v.validateParameters(params)
	if isArrow {
		v.validateEmptyName(name)
		if body.IsBlock() {
			v.validateBlock(body)
		} else {
			v.ValidateExpression(body)
		}
	} else {
		v.validateOptionalName(name)
		v.validateFunctionBody(body, isAmbient)
	}
	if isFn {
		v.validateFunctionFeatures(ast.MustFunction(n))
	} else {
		v.validateFunctionNodeFeatures(n)
	}
}

// Java: AstValidator.validateFunctionFeatures(Node)
func (v *AstValidator) validateFunctionFeatures(fn ast.Function) {
	n := fn.Node()
	if fn.IsArrow() {
		v.validateFeature(parser.FeatureSet_Feature_ARROW_FUNCTIONS, n)
	}
	isGen := fn.IsGenerator()
	isAsync := fn.IsAsync()
	if isGen {
		v.validateFeature(parser.FeatureSet_Feature_GENERATORS, n)
	}
	if isAsync {
		v.validateFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS, n)
	}
	if isAsync && isGen {
		v.validateFeature(parser.FeatureSet_Feature_ASYNC_GENERATORS, n)
	}
}

func (v *AstValidator) validateFunctionNodeFeatures(n rhino.Node) {
	if n.IsArrowFunction() {
		v.validateFeature(parser.FeatureSet_Feature_ARROW_FUNCTIONS, n)
	}
	if n.IsGeneratorFunction() {
		v.validateFeature(parser.FeatureSet_Feature_GENERATORS, n)
	}
	if n.IsAsyncFunction() {
		v.validateFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS, n)
	}
	if n.IsAsyncFunction() && n.IsGeneratorFunction() {
		v.validateFeature(parser.FeatureSet_Feature_ASYNC_GENERATORS, n)
	}
}

// Java: AstValidator.validateFunctionBody(Node, boolean)
func (v *AstValidator) validateFunctionBody(n rhino.Node, noBlock bool) {
	if noBlock {
		v.validateNodeType(rhino.Token_EMPTY, n)
	} else {
		v.validateBlock(n)
	}
}

// Java: AstValidator.validateParameters(Node)
func (v *AstValidator) validateParameters(n rhino.Node) {
	if v.validateNodeType(rhino.Token_PARAM_LIST, n) {
		for p := range ast.MustParamList(n).Params() {
			v.validateParameter(p.Node())
		}
		return
	}
	for c := range n.Children() {
		v.validateParameter(c)
	}
}

func (v *AstValidator) validateParameter(c rhino.Node) {
	switch c.GetToken() {
	case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
		v.validateRestParameters(rhino.Token_PARAM_LIST, c)
	case rhino.Token_DEFAULT_VALUE:
		v.validateFeature(parser.FeatureSet_Feature_DEFAULT_PARAMETERS, c)
		v.validateDefaultValue(rhino.Token_PARAM_LIST, ast.MustDefaultValue(c))
	case rhino.Token_NAME:
		v.validateName(c)
	case rhino.Token_ARRAY_PATTERN:
		v.validateArrayPattern(rhino.Token_PARAM_LIST, c)
	default:
		v.validateObjectPattern(rhino.Token_PARAM_LIST, c)
	}
}

// Java: AstValidator.validateDefaultValue(Token, Node)
func (v *AstValidator) validateDefaultValue(contextType rhino.Token, dv ast.DefaultValue) {
	n := dv.Node()
	if v.validatePropsAndArity(n) {
		v.validateLHS(contextType, dv.Target().Node())
		v.ValidateExpression(dv.Value().Node())
		return
	}
	v.validateLHS(contextType, n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateRestParameters(Token, Node)
func (v *AstValidator) validateRestParameters(contextType rhino.Token, n rhino.Node) {
	v.validateFeature(parser.FeatureSet_Feature_REST_PARAMETERS, n)
	v.validateRest(contextType, n)
}

// Java: AstValidator.validateArrayPatternRest(Token, Node)
func (v *AstValidator) validateArrayPatternRest(contextType rhino.Token, n rhino.Node) {
	v.validateFeature(parser.FeatureSet_Feature_ARRAY_PATTERN_REST, n)
	v.validateRest(contextType, n)
}

// Java: AstValidator.validateObjectPatternRest(Token, Node)
func (v *AstValidator) validateObjectPatternRest(contextType rhino.Token, n rhino.Node) {
	v.validateFeature(parser.FeatureSet_Feature_OBJECT_PATTERN_REST, n)
	v.validateRest(contextType, n)
}

// Java: AstValidator.validateRest(Token, Node)
func (v *AstValidator) validateRest(contextType rhino.Token, n rhino.Node) {
	rest, ok := ast.AsRest(n)
	if !ok {
		v.Violation("Unexpected node type.", n)
		return
	}
	v.validatePropsAndArity(n)
	v.validateLHS(contextType, rest.Target().Node())
	if !n.GetNext().IsNil() {
		v.Violation("Rest parameters must come after all other parameters.", n)
	}
}

// validateNameDeclaration validates a VAR, LET, or CONST declaration within
// its enclosing statement (which may be a FOR / FOR_IN / FOR_OF / DECLARE).
//
// Java: AstValidator.validateNameDeclarationHelper(Node, Token, Node)
func (v *AstValidator) validateNameDeclaration(statement rhino.Node, decl ast.NameDeclaration) {
	n := decl.Node()
	v.validateProperties(n)
	v.validateMinimumChildCount(n, 1)
	kind := decl.Token()
	for d := range decl.Declarators() {
		v.validateNameDeclarationChild(statement, kind, d.Node())
	}
	switch kind {
	case rhino.Token_LET:
		v.validateFeature(parser.FeatureSet_Feature_LET_DECLARATIONS, n)
	case rhino.Token_CONST:
		v.validateFeature(parser.FeatureSet_Feature_CONST_DECLARATIONS, n)
	}
}

// Java: AstValidator.validateNameDeclarationChild(Node, Token, Node)
func (v *AstValidator) validateNameDeclarationChild(statement rhino.Node, declaration rhino.Token, n rhino.Node) {
	var minValues, maxValues int
	switch {
	case statement.IsForIn() && declaration == rhino.Token_VAR:
		// ECMAScript 5 sloppy mode allows for-in var initializers.
		minValues, maxValues = 0, 1
	case NodeUtil_IsEnhancedFor(statement):
		minValues, maxValues = 0, 0
	case n.IsDestructuringLhs() || declaration == rhino.Token_CONST:
		minValues, maxValues = 1, 1
	default:
		minValues, maxValues = 0, 1
	}

	if nameNode, ok := ast.AsName(n); ok {
		// Don't use validateName here since this NAME node may have an initializer child.
		v.validateNonEmptyString(n)
		v.validateProperties(n)
		v.validateChildCountIn(n, minValues, maxValues)
		if init, ok := nameNode.Init(); ok {
			v.ValidateExpression(init.Node())
		}
		return
	}

	if lhs, ok := ast.AsDestructuringLhs(n); ok {
		v.validateProperties(n)
		v.validateChildCountIn(n, 1+minValues, 1+maxValues)
		pat := lhs.Pattern().Node()
		switch pat.GetToken() {
		case rhino.Token_ARRAY_PATTERN:
			v.validateArrayPattern(declaration, pat)
		case rhino.Token_OBJECT_PATTERN:
			v.validateObjectPattern(declaration, pat)
		default:
			v.Violation("Invalid destructuring lhs first child for "+declaration.String()+" node", n)
		}
		if n.HasTwoChildren() {
			init, _ := lhs.Init()
			v.ValidateExpression(init.Node())
		}
		return
	}

	v.Violation("Invalid child for "+declaration.String()+" node", n)
}

// Java: AstValidator.validateLHS(Token, Node)
func (v *AstValidator) validateLHS(contextType rhino.Token, n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		v.validateName(n)
	case rhino.Token_ARRAY_PATTERN:
		v.validateArrayPattern(contextType, n)
	case rhino.Token_OBJECT_PATTERN:
		v.validateObjectPattern(contextType, n)
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		v.validateGetPropGetElemInLHS(contextType, n)
	case rhino.Token_CAST:
		// GetOnlyChild intentionally throws IllegalStateException if CAST arity != 1.
		v.validateLHS(contextType, n.GetOnlyChild())
	default:
		v.Violation("Invalid child for "+contextType.String()+" node", n)
	}
}

// Java: AstValidator.validateArrayPattern(Token, Node)
func (v *AstValidator) validateArrayPattern(contextType rhino.Token, n rhino.Node) {
	v.validateFeature(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, n)
	if v.validateNodeType(rhino.Token_ARRAY_PATTERN, n) {
		for elem := range ast.MustArrayPattern(n).Elements() {
			v.validateArrayPatternElement(contextType, elem.Node())
		}
		return
	}
	for c := range n.Children() {
		v.validateArrayPatternElement(contextType, c)
	}
}

func (v *AstValidator) validateArrayPatternElement(contextType rhino.Token, c rhino.Node) {
	switch c.GetToken() {
	case rhino.Token_DEFAULT_VALUE:
		v.validateDefaultValue(contextType, ast.MustDefaultValue(c))
	case rhino.Token_ITER_REST:
		v.validateArrayPatternRest(contextType, c)
	case rhino.Token_EMPTY:
		v.validateChildlessNode(c)
	default:
		v.validateLHS(contextType, c)
	}
}

// Java: AstValidator.validateObjectPattern(Token, Node)
func (v *AstValidator) validateObjectPattern(contextType rhino.Token, n rhino.Node) {
	v.validateFeature(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, n)
	if v.validateNodeType(rhino.Token_OBJECT_PATTERN, n) {
		for prop := range ast.MustObjectPattern(n).Props() {
			v.validateObjectPatternProp(contextType, n, prop.Node())
		}
		return
	}
	for c := range n.Children() {
		v.validateObjectPatternProp(contextType, n, c)
	}
}

func (v *AstValidator) validateObjectPatternProp(contextType rhino.Token, patternNode, c rhino.Node) {
	switch c.GetToken() {
	case rhino.Token_STRING_KEY:
		v.validateObjectPatternStringKey(contextType, ast.MustStringKey(c))
	case rhino.Token_OBJECT_REST:
		v.validateObjectPatternRest(contextType, c)
	case rhino.Token_COMPUTED_PROP:
		v.validateObjectPatternComputedPropKey(contextType, ast.MustComputedProp(c))
	default:
		v.Violation("Invalid object pattern child for "+contextType.String()+" node", patternNode)
	}
}

// Java: AstValidator.validateObjectPatternStringKey(Token, Node)
func (v *AstValidator) validateObjectPatternStringKey(contextType rhino.Token, sk ast.StringKey) {
	n := sk.Node()
	v.validateObjectLiteralKeyName(n)
	v.validatePropsAndCount(n, 1)

	val := sk.Value().Node()
	if dv, ok := ast.AsDefaultValue(val); ok {
		v.validateDefaultValue(contextType, dv)
	} else {
		v.validateLHS(contextType, val)
	}
}

// Java: AstValidator.validateObjectPatternComputedPropKey(Token, Node)
func (v *AstValidator) validateObjectPatternComputedPropKey(contextType rhino.Token, prop ast.ComputedProp) {
	n := prop.Node()
	valid := v.validatePropsAndArity(n)
	v.validateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.ValidateExpression(prop.Key().Node())

	var target rhino.Node
	if valid {
		val, _ := prop.Value()
		target = val.Node()
	} else {
		target = n.GetLastChild()
	}
	if dv, ok := ast.AsDefaultValue(target); ok {
		v.validateDefaultValue(contextType, dv)
	} else {
		v.validateLHS(target.GetToken(), target)
	}
}

// Java: AstValidator.validateNamedType(Node)
func (v *AstValidator) validateNamedType(n rhino.Node) {
	v.validateNodeType(rhino.Token_NAMED_TYPE, n)
	v.validatePropsAndArity(n)
	v.validateName(n.GetFirstChild())
}

// Java: AstValidator.validateTypeAlias(Node)
func (v *AstValidator) validateTypeAlias(n rhino.Node) {
	v.validateNodeType(rhino.Token_TYPE_ALIAS, n)
	v.validatePropsAndArity(n)
}

// Java: AstValidator.validateAmbientDeclaration(Node)
func (v *AstValidator) validateAmbientDeclaration(n rhino.Node) {
	v.validateNodeType(rhino.Token_DECLARE, n)
	v.validateAmbientDeclarationHelper(n.GetFirstChild())
}

// Java: AstValidator.validateAmbientDeclarationHelper(Node)
func (v *AstValidator) validateAmbientDeclarationHelper(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		v.validateNameDeclaration(n.GetParent(), ast.MustNameDeclaration(n))
	case rhino.Token_FUNCTION:
		v.validateFunctionSignature(n)
	case rhino.Token_CLASS:
		v.validateClassDeclaration(ast.MustClass(n), true)
	case rhino.Token_ENUM:
		v.validateEnum(n)
	case rhino.Token_NAMESPACE:
		v.validateNamespace(n, true)
	case rhino.Token_TYPE_ALIAS:
		v.validateTypeAlias(n)
	case rhino.Token_EXPORT:
		v.validateExport(ast.MustExport(n), true)
	}
}

// Java: AstValidator.validateNamespace(Node, boolean)
func (v *AstValidator) validateNamespace(n rhino.Node, isAmbient bool) {
	v.validateNodeType(rhino.Token_NAMESPACE, n)
	v.validatePropsAndArity(n)
	v.validateNamespaceName(n.GetFirstChild())
	v.validateNamespaceElements(n.GetLastChild(), isAmbient)
}

// Java: AstValidator.validateNamespaceName(Node)
func (v *AstValidator) validateNamespaceName(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		v.validateName(n)
	case rhino.Token_GETPROP:
		v.validateGetProp(ast.MustGetProp(n))
	}
}

// Java: AstValidator.validateNamespaceElements(Node, boolean)
func (v *AstValidator) validateNamespaceElements(n rhino.Node, isAmbient bool) {
	v.validateNodeType(rhino.Token_NAMESPACE_ELEMENTS, n)
	for c := range n.Children() {
		if isAmbient {
			v.validateAmbientDeclarationHelper(c)
		} else {
			v.ValidateStatement(c)
		}
	}
}
