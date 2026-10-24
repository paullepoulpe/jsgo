// Port of com/google/javascript/jscomp/TypedCodeGenerator.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypedCodeGenerator is a code generator that outputs type annotations for
// functions and constructors.
type TypedCodeGenerator struct {
	CodeGeneratorBase

	registry         *rhino.JSTypeRegistry
	jsDocInfoPrinter *JSDocInfoPrinter
}

// Java: TypedCodeGenerator(CodeConsumer, CompilerOptions, JSTypeRegistry)
func NewTypedCodeGenerator(consumer CodeConsumer, options *CompilerOptions, registry *rhino.JSTypeRegistry) *TypedCodeGenerator {
	t := &TypedCodeGenerator{}
	t.InitCodeGeneratorBaseWithOptions(t, consumer, options)
	jutil.CheckNotNull(registry)
	t.registry = registry
	t.jsDocInfoPrinter = NewJSDocInfoPrinter(options.GetUseOriginalNamesInOutput())
	return t
}

// Java: TypedCodeGenerator.add(Node, Context)
func (t *TypedCodeGenerator) AddNodeWithContext(n *rhino.Node, context CodeGenerator_Context) {
	t.MaybeAddTypeAnnotation(n)
	t.CodeGeneratorBase.AddNodeWithContext(n, context)
}

// Java: TypedCodeGenerator.maybeAddTypeAnnotation(Node)
func (t *TypedCodeGenerator) MaybeAddTypeAnnotation(n *rhino.Node) {
	parent := n.GetParent()
	if parent == nil {
		// root node cannot have a type annotation.
		return
	}
	// Generate type annotations only for statements and class member
	// functions.
	if parent.IsBlock() || parent.IsScript() || parent.IsClassMembers() {
		if n.IsClass() || n.IsFunction() || n.IsMemberFunctionDef() {
			t.Add(t.GetTypeAnnotation(n))
		} else if n.IsExprResult() &&
			n.GetFirstChild().IsAssign() {
			assign := n.GetFirstChild()
			if NodeUtil_IsNamespaceDecl(assign.GetFirstChild()) {
				t.Add(t.jsDocInfoPrinter.Print(assign.GetJSDocInfo()))
			} else {
				rhs := assign.GetLastChild()
				t.Add(t.GetTypeAnnotation(rhs))
			}
		} else if NodeUtil_IsNameDeclaration(n) && n.GetFirstFirstChild() != nil {
			// All namespace declarations except `const x = {};` are signified
			// by @const JSDoc.
			if NodeUtil_IsNamespaceDecl(n.GetFirstChild()) && n.GetJSDocInfo() != nil {
				t.Add(t.jsDocInfoPrinter.Print(n.GetJSDocInfo()))
			} else {
				t.Add(t.GetTypeAnnotation(n.GetFirstFirstChild()))
			}
		}
	}
}

// Java: TypedCodeGenerator.getTypeAnnotation(Node)
func (t *TypedCodeGenerator) GetTypeAnnotation(node *rhino.Node) string {
	if node.IsMemberFunctionDef() {
		// For a member function the type information is actually on the
		// function it contains, so just generate the type annotation for
		// that.
		return t.GetMemberFunctionAnnotation(node.GetOnlyChild())
	} else if node.IsClass() {
		return t.GetClassAnnotation(node.GetJSType())
	} else if node.IsFunction() {
		return t.GetFunctionAnnotation(node)
	} else {
		nodeOriginallyHadJSDoc := NodeUtil_GetBestJSDocInfo(node) != nil
		if !nodeOriginallyHadJSDoc {
			// For nodes that don't inherently define a type, ony generate
			// JSDoc if they originally had some.
			return ""
		}

		type_ := node.GetJSType()
		if type_ == nil {
			return ""
		} else if type_.IsFunctionType() {
			return t.GetFunctionAnnotation(node)
		} else if type_.IsEnumType() {
			return "/** @enum {" +
				type_.ToMaybeObjectType().
					GetEnumeratedTypeOfEnumObject().
					ToAnnotationString(rhino.JSType_Nullability_EXPLICIT) +
				"} */\n"
		} else if !type_.IsUnknownType() &&
			!type_.IsEmptyType() &&
			!type_.IsVoidType() &&
			!type_.IsFunctionPrototypeType() {
			return "/** @type {" + node.GetJSType().ToAnnotationString(rhino.JSType_Nullability_EXPLICIT) + "} */\n"
		} else {
			return ""
		}
	}
}

// GetFunctionAnnotation: fnNode is a node for a function for which to
// generate a type annotation.
//
// Java: TypedCodeGenerator.getFunctionAnnotation(Node)
func (t *TypedCodeGenerator) GetFunctionAnnotation(fnNode *rhino.Node) string {
	type_ := fnNode.GetJSType()
	jutil.CheckState(fnNode.IsFunction() || type_.IsFunctionType())

	if type_ == nil || type_.IsUnknownType() {
		return ""
	}

	funType := type_.ToMaybeFunctionType()
	if type_.Equals(t.registry.GetNativeType(rhino.JSTypeNative_FUNCTION_TYPE)) {
		return "/** @type {!Function} */\n"
	}
	var sb strings.Builder
	sb.WriteString("/**\n")
	var paramNode *rhino.Node
	// We need to use the child nodes of the function as the nodes for the
	// parameters of the function type do not have the real parameter names.
	// FUNCTION
	//   NAME
	//   PARAM_LIST
	//     NAME param1
	//     NAME param2
	if fnNode != nil && fnNode.IsFunction() {
		paramNode = NodeUtil_GetFunctionParameters(fnNode).GetFirstChild()
	}

	// Param types
	t.AppendFunctionParamAnnotations(&sb, funType, paramNode)

	// Return type
	retType := funType.GetReturnType()
	if retType != nil &&
		!retType.IsEmptyType() && // There is no annotation for the empty type.
		!funType.IsInterface() && // Interfaces never return a value.
		!(funType.IsConstructor() && retType.IsVoidType()) {
		sb.WriteString(" * ")
		TypedCodeGenerator_AppendAnnotation(&sb, "return", retType.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT))
		sb.WriteString("\n")
	}

	// This function could be defining an ES5-style class or interface.
	// If it isn't but still requires a type for `this`, then we need to
	// explicitly add an annotation for that.
	if funType.IsConstructor() {
		// This function is defining an ES5-style class, so include the class
		// annotations here.
		t.AppendClassAnnotations(&sb, funType)
		sb.WriteString(" * @constructor\n")
	} else if funType.IsInterface() {
		t.AppendInterfaceAnnotations(&sb, funType)
	} else {
		thisType := funType.GetTypeOfThis()
		if thisType != nil && !thisType.IsUnknownType() && !thisType.IsVoidType() {
			if fnNode == nil || !thisType.Equals(typedCodeGenerator_objectTypeOrNull(t.FindMethodOwner(fnNode))) {
				sb.WriteString(" * ")
				TypedCodeGenerator_AppendAnnotation(&sb, "this", thisType.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT))
				sb.WriteString("\n")
			}
		}
	}

	t.AppendTemplateAnnotations(&sb, funType.GetTypeParameters())

	sb.WriteString(" */\n")
	return sb.String()
}

// typedCodeGenerator_objectTypeOrNull passes an ObjectType to
// Object.equals(Object), keeping null as null.
//
// Go-only: language-adapter — a nil interface must stay an untyped nil.
func typedCodeGenerator_objectTypeOrNull(t rhino.ObjectType) any {
	if t == nil {
		return nil
	}
	return t
}

// GetMemberFunctionAnnotation: fnNode is a function node child of a
// MEMBER_FUNCTION_DEF.
//
// Java: TypedCodeGenerator.getMemberFunctionAnnotation(Node)
func (t *TypedCodeGenerator) GetMemberFunctionAnnotation(fnNode *rhino.Node) string {
	jutil.CheckStateMsg1(fnNode.IsFunction() && fnNode.GetParent().IsMemberFunctionDef(), "%s", fnNode)
	type_ := fnNode.GetJSType()

	if type_ == nil || type_.IsUnknownType() {
		return ""
	}

	funType := type_.ToMaybeFunctionType()
	var sb strings.Builder
	sb.WriteString("/**\n")

	// We need to use the child nodes of the function as the nodes for the
	// parameters of the function type do not have the real parameter names.
	// FUNCTION
	//   NAME
	//   PARAM_LIST
	//     NAME param1
	//     NAME param2
	paramNode := NodeUtil_GetFunctionParameters(fnNode).GetFirstChild()

	// Param types
	t.AppendFunctionParamAnnotations(&sb, funType, paramNode)

	if NodeUtil_IsEs6Constructor(fnNode) {
		t.AppendTemplateAnnotations(&sb, funType.GetConstructorOnlyTemplateParameters())
		// no return type for the constructor
	} else {
		t.AppendTemplateAnnotations(&sb, funType.GetTypeParameters())
		// Return type
		retType := funType.GetReturnType()
		if retType != nil && !retType.IsEmptyType() {
			// There is no annotation for the empty type.
			sb.WriteString(" * ")
			TypedCodeGenerator_AppendAnnotation(&sb, "return", retType.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT))
			sb.WriteString("\n")
		}
	}

	sb.WriteString(" */\n")
	return sb.String()
}

// AppendFunctionParamAnnotations generates @param annotations.
//
// sb is where annotations will be appended; funType is the function type;
// paramNode is where parameter names will be taken from.
//
// Java: TypedCodeGenerator.appendFunctionParamAnnotations(StringBuilder, FunctionType, Node)
func (t *TypedCodeGenerator) AppendFunctionParamAnnotations(
	sb *strings.Builder, funType rhino.FunctionType, paramNode *rhino.Node) {
	minArity := funType.GetMinArity()
	maxArity := funType.GetMaxArity()
	formals := funType.GetParameters()
	for i := 0; i < len(formals); i++ {
		sb.WriteString(" * ")
		TypedCodeGenerator_AppendAnnotation(sb, "param", t.GetParameterJSDocType(formals, i, minArity, maxArity))
		parameterName := t.GetParameterJSDocName(paramNode, i)
		sb.WriteString(" ")
		sb.WriteString(parameterName)
		sb.WriteString("\n")
		if paramNode != nil {
			paramNode = paramNode.GetNext()
		}
	}
}

// Java: TypedCodeGenerator.getClassAnnotation(JSType)
func (t *TypedCodeGenerator) GetClassAnnotation(classType rhino.JSType) string {
	if classType == nil || classType.IsUnknownType() {
		return ""
	}

	jutil.CheckStateMsg1(classType.IsFunctionType(), "%s", classType)

	funType := classType.ToMaybeFunctionType()
	var sb strings.Builder

	if funType.IsInterface() {
		t.AppendInterfaceAnnotations(&sb, funType)
	} else {
		jutil.CheckStateMsg1(funType.IsConstructor(), "%s", funType)
		t.AppendClassAnnotations(&sb, funType)
	}

	t.AppendTemplateAnnotations(&sb, funType.GetTypeParameters())

	jsdocContent := sb.String()

	// For simple class, it's possible we didn't end up generating any JSDoc
	// at all.
	if jsdocContent == "" {
		return jsdocContent
	} else {
		return "/**\n" + jsdocContent + " */\n"
	}
}

// Java: TypedCodeGenerator.appendTemplateAnnotations(StringBuilder, Collection<? extends JSType>)
func (t *TypedCodeGenerator) AppendTemplateAnnotations(
	sb *strings.Builder, typeParams []*rhino.TemplateType) {
	if len(typeParams) != 0 {
		sb.WriteString(" * @template ")
		for i, typeParam := range typeParams {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(t.FormatTypeVar(typeParam))
		}
		sb.WriteString("\n")
	}
}

// GetParameterJSDocName returns the name of the parameter to be used in
// JSDoc, generating one for destructuring parameters.
//
// paramNode is a child node of a parameter list; paramIndex is the position
// of the child in the list. Returns the name to use in JSDoc.
//
// Java: TypedCodeGenerator.getParameterJSDocName(Node, int)
func (t *TypedCodeGenerator) GetParameterJSDocName(paramNode *rhino.Node, paramIndex int) string {
	var nameNode *rhino.Node
	if paramNode != nil {
		jutil.CheckArgumentMsg1(paramNode.GetParent().IsParamList(), "%s", paramNode)
		if paramNode.IsRest() {
			// use `restParam` of `...restParam`
			// restParam might still be a destructuring pattern
			paramNode = paramNode.GetOnlyChild()
		} else if paramNode.IsDefaultValue() {
			// use `defaultParam` of `defaultParam = something`
			// defaultParam might still be a destructuring pattern
			paramNode = paramNode.GetFirstChild()
		}
		if paramNode.IsName() {
			nameNode = paramNode
		} else {
			jutil.CheckStateMsg1(paramNode.IsObjectPattern() || paramNode.IsArrayPattern(), "%s", paramNode)
			nameNode = nil // must generate a fake name
		}
	}
	if nameNode == nil {
		return "p" + strconv.Itoa(paramIndex)
	} else {
		jutil.CheckStateMsg1(nameNode.IsName(), "%s", nameNode)
		return nameNode.GetString()
	}
}

// Java: TypedCodeGenerator.formatTypeVar(JSType)
func (t *TypedCodeGenerator) FormatTypeVar(var_ rhino.JSType) string {
	return var_.ToAnnotationString(rhino.JSType_Nullability_IMPLICIT)
}

// TODO(dimvar): it's awkward that we print @constructor after the
// extends/implements; we should print it first, like users write it. Same
// for @interface and @record.
//
// Java: TypedCodeGenerator.appendClassAnnotations(StringBuilder, FunctionType)
func (t *TypedCodeGenerator) AppendClassAnnotations(sb *strings.Builder, funType rhino.FunctionType) {
	superConstructor := funType.GetInstanceType().GetSuperClassConstructor()
	if superConstructor != nil {
		superInstance := superConstructor.GetInstanceType()
		if superInstance.String() != "Object" {
			sb.WriteString(" * ")
			TypedCodeGenerator_AppendAnnotation(sb, "extends", superInstance.ToAnnotationString(rhino.JSType_Nullability_IMPLICIT))
			sb.WriteString("\n")
		}
	}
	// Avoid duplicates, add implemented type to a set first
	interfaces := jutil.NewTreeSet(jutil.StringCompareTo)
	for _, interfaze := range funType.GetAncestorInterfaces() {
		interfaces.Add(interfaze.ToAnnotationString(rhino.JSType_Nullability_IMPLICIT))
	}
	for interfaze := range interfaces.All() {
		sb.WriteString(" * ")
		TypedCodeGenerator_AppendAnnotation(sb, "implements", interfaze)
		sb.WriteString("\n")
	}
}

// Java: TypedCodeGenerator.appendInterfaceAnnotations(StringBuilder, FunctionType)
func (t *TypedCodeGenerator) AppendInterfaceAnnotations(sb *strings.Builder, funType rhino.FunctionType) {
	interfaces := jutil.NewTreeSet(jutil.StringCompareTo)
	for _, interfaceType := range funType.GetAncestorInterfaces() {
		interfaces.Add(interfaceType.ToAnnotationString(rhino.JSType_Nullability_IMPLICIT))
	}
	for interfaze := range interfaces.All() {
		sb.WriteString(" * ")
		TypedCodeGenerator_AppendAnnotation(sb, "extends", interfaze)
		sb.WriteString("\n")
	}
	if funType.IsStructuralInterface() {
		sb.WriteString(" * @record\n")
	} else {
		sb.WriteString(" * @interface\n")
	}
}

// TODO(sdh): This whole method could be deleted if we don't mind adding
// additional @this annotations where they're not actually necessary.

// FindMethodOwner: given a method definition node, returns the ObjectType
// corresponding to the class the method is defined on, or null if it is not
// a prototype method.
//
// Java: TypedCodeGenerator.findMethodOwner(Node)
func (t *TypedCodeGenerator) FindMethodOwner(n *rhino.Node) rhino.ObjectType {
	if n == nil {
		return nil
	}
	parent := n.GetParent()
	var ctor rhino.FunctionType
	if parent.IsAssign() {
		target := parent.GetFirstChild()
		if NodeUtil_IsPrototypeProperty(target) {
			// TODO(johnlenz): handle non-global types
			type_ := t.registry.GetGlobalType(*target.GetFirstFirstChild().GetQualifiedName())
			if type_ != nil {
				ctor = type_.(rhino.ObjectType).GetConstructor()
			} else {
				ctor = nil
			}
		}
	} else if parent.IsClass() {
		// TODO(sdh): test this case once the type checker understands ES6
		// classes
		ctor = parent.GetJSType().ToMaybeFunctionType()
	}
	if ctor != nil {
		return ctor.GetInstanceType()
	}
	return nil
}

// Java: TypedCodeGenerator.appendAnnotation(StringBuilder, String, String)
func TypedCodeGenerator_AppendAnnotation(sb *strings.Builder, name string, type_ string) {
	sb.WriteString("@")
	sb.WriteString(name)
	sb.WriteString(" {")
	sb.WriteString(type_)
	sb.WriteString("}")
}

// GetParameterJSDocType creates a JSDoc-suitable String representation of
// the type of a parameter.
//
// Java: TypedCodeGenerator.getParameterJSDocType(List<FunctionType.Parameter>, int, int, int)
func (t *TypedCodeGenerator) GetParameterJSDocType(
	parameters []*rhino.FunctionType_Parameter, index int, minArgs int, maxArgs int) string {
	type_ := parameters[index].GetJSType()
	if index < minArgs {
		return type_.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT)
	}
	isRestArgument := maxArgs == math.MaxInt32 && index == len(parameters)-1
	if isRestArgument {
		return "..." + t.RestrictByUndefined(type_).ToAnnotationString(rhino.JSType_Nullability_EXPLICIT)
	}
	return t.RestrictByUndefined(type_).ToAnnotationString(rhino.JSType_Nullability_EXPLICIT) + "="
}

// RestrictByUndefined removes undefined from a union type.
//
// Java: TypedCodeGenerator.restrictByUndefined(JSType)
func (t *TypedCodeGenerator) RestrictByUndefined(type_ rhino.JSType) rhino.JSType {
	// If not voidable, there's nothing to do. If not nullable then the
	// easiest thing is to simply remove both null and undefined. If
	// nullable, then add null back into the union after removing null and
	// undefined.
	if !type_.IsVoidable() {
		return type_
	}
	restricted := type_.RestrictByNotNullOrUndefined()
	if type_.IsNullable() {
		nullType := t.registry.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
		return t.registry.CreateUnionTypeOfList([]rhino.JSType{restricted, nullType})
	}
	// The bottom type cannot appear in a jsdoc
	if restricted.IsEmptyType() {
		return type_
	}
	return restricted
}
