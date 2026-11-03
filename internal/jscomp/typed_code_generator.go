// Port of com/google/javascript/jscomp/TypedCodeGenerator.java @ closure-compiler v20261028

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
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func (t *CodeGenerator) initTypedWithOptions(options *CompilerOptions, registry *rhino.JSTypeRegistry) {
	t.initWithOptions(options)
	jutil.CheckNotNull(registry)
	t.registry = registry
	t.emitTypeAnnotations = true
}

// InitTypedCodeGenerator initializes a CodeGenerator to emit type annotations
// for functions and constructors (Java's TypedCodeGenerator subclass).
//
// Java: TypedCodeGenerator(CodeConsumer, CompilerOptions, JSTypeRegistry)
func (t *CodeGenerator) InitTypedCodeGenerator(consumer *CodeConsumer, options *CompilerOptions, registry *rhino.JSTypeRegistry) {
	if consumer != nil {
		t.cc = *consumer
	}
	t.initTypedWithOptions(options, registry)
}

// Java: TypedCodeGenerator(CodeConsumer, CompilerOptions, JSTypeRegistry)
func NewTypedCodeGenerator(consumer *CodeConsumer, options *CompilerOptions, registry *rhino.JSTypeRegistry) *CodeGenerator {
	t := &CodeGenerator{}
	t.InitTypedCodeGenerator(consumer, options, registry)
	return t
}

// Java: TypedCodeGenerator.maybeAddTypeAnnotation(Node)
func (t *CodeGenerator) MaybeAddTypeAnnotation(n rhino.Node) {
	parent := n.GetParent()
	// Generate type annotations only for statements and class member functions.
	if parent.IsNil() || (!parent.IsBlock() && !parent.IsScript() && !parent.IsClassMembers()) {
		return
	}
	if n.IsClass() || n.IsFunction() || n.IsMemberFunctionDef() {
		t.Add(t.GetTypeAnnotation(n))
	} else if stmt, ok := ast.AsExprStmt(n); ok {
		if assign, isAssign := stmt.Expr().AsAssign(); isAssign {
			if NodeUtil_IsNamespaceDecl(assign.Target().Node()) {
				t.Add(t.jsDocInfoPrinter.Print(assign.Node().GetJSDocInfo()))
			} else {
				t.Add(t.GetTypeAnnotation(assign.Value().Node()))
			}
		}
	} else if decl, ok := ast.AsNameDeclaration(n); ok {
		if first, hasFirst := decl.FirstDeclarator(); hasFirst {
			if name, isName := first.AsName(); isName {
				if init, hasInit := name.Init(); hasInit {
					// All namespace declarations except `const x = {};` are signified
					// by @const JSDoc.
					if NodeUtil_IsNamespaceDecl(name.Node()) && n.GetJSDocInfo() != nil {
						t.Add(t.jsDocInfoPrinter.Print(n.GetJSDocInfo()))
					} else {
						t.Add(t.GetTypeAnnotation(init.Node()))
					}
				}
			} else {
				t.Add(t.GetTypeAnnotation(first.MustDestructuringLhs().Pattern().Node()))
			}
		}
	}
}

// Java: TypedCodeGenerator.getTypeAnnotation(Node)
func (t *CodeGenerator) GetTypeAnnotation(node rhino.Node) string {
	if mfd, ok := ast.AsMemberFunctionDef(node); ok {
		// For a member function the type information is actually on the
		// function it contains, so just generate the type annotation for
		// that.
		return t.getMemberFunctionAnnotation(mfd.Function())
	} else if cls, ok := ast.AsClass(node); ok {
		return t.GetClassAnnotation(cls.Node().GetJSType())
	} else if fn, ok := ast.AsFunction(node); ok {
		return t.GetFunctionAnnotation(fn.Node())
	}

	if NodeUtil_GetBestJSDocInfo(node) == nil {
		// For nodes that don't inherently define a type, only generate
		// JSDoc if they originally had some.
		return ""
	}
	type_ := node.GetJSType()
	if type_ == nil {
		return ""
	}
	if type_.IsFunctionType() {
		return t.GetFunctionAnnotation(node)
	}
	if type_.IsEnumType() {
		return "/** @enum {" +
			type_.ToMaybeObjectType().
				GetEnumeratedTypeOfEnumObject().
				ToAnnotationString(rhino.JSType_Nullability_EXPLICIT) +
			"} */\n"
	}
	if !type_.IsUnknownType() &&
		!type_.IsEmptyType() &&
		!type_.IsVoidType() &&
		!type_.IsFunctionPrototypeType() {
		return "/** @type {" + type_.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT) + "} */\n"
	}
	return ""
}

// GetFunctionAnnotation: fnNode is a node for a function for which to
// generate a type annotation.
//
// Java: TypedCodeGenerator.getFunctionAnnotation(Node)
func (t *CodeGenerator) GetFunctionAnnotation(fnNode rhino.Node) string {
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
	var param ast.Param
	var hasParam bool
	// We need to use the child nodes of the function as the nodes for the
	// parameters of the function type do not have the real parameter names.
	if fn, ok := ast.AsFunction(fnNode); ok {
		param, hasParam = fn.Params().FirstParam()
	}

	// Param types
	t.appendFunctionParamAnnotations(&sb, funType, param, hasParam)

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
			var owner rhino.ObjectType
			if !fnNode.IsNil() {
				owner = t.FindMethodOwner(fnNode)
			}
			if owner == nil || !thisType.Equals(owner) {
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

func (t *CodeGenerator) getMemberFunctionAnnotation(fn ast.Function) string {
	type_ := fn.Node().GetJSType()
	if type_ == nil || type_.IsUnknownType() {
		return ""
	}

	funType := type_.ToMaybeFunctionType()
	var sb strings.Builder
	sb.WriteString("/**\n")

	param, hasParam := fn.Params().FirstParam()
	t.appendFunctionParamAnnotations(&sb, funType, param, hasParam)

	if NodeUtil_IsEs6Constructor(fn.Node()) {
		t.AppendTemplateAnnotations(&sb, funType.GetConstructorOnlyTemplateParameters())
		// no return type for the constructor
	} else {
		t.AppendTemplateAnnotations(&sb, funType.GetTypeParameters())
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

func (t *CodeGenerator) appendFunctionParamAnnotations(
	sb *strings.Builder, funType rhino.FunctionType, param ast.Param, hasParam bool) {
	minArity := funType.GetMinArity()
	maxArity := funType.GetMaxArity()
	formals := funType.GetParameters()
	for i := 0; i < len(formals); i++ {
		sb.WriteString(" * ")
		TypedCodeGenerator_AppendAnnotation(sb, "param", t.GetParameterJSDocType(formals, i, minArity, maxArity))
		sb.WriteString(" ")
		sb.WriteString(getParameterJSDocName(param, hasParam, i))
		sb.WriteString("\n")
		if hasParam {
			param, hasParam = param.Next()
		}
	}
}

// Java: TypedCodeGenerator.getClassAnnotation(JSType)
func (t *CodeGenerator) GetClassAnnotation(classType rhino.JSType) string {
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
		return ""
	}
	return "/**\n" + jsdocContent + " */\n"
}

// Java: TypedCodeGenerator.appendTemplateAnnotations(StringBuilder, Collection<? extends JSType>)
func (t *CodeGenerator) AppendTemplateAnnotations(
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

func getParameterJSDocName(param ast.Param, hasParam bool, paramIndex int) string {
	if hasParam {
		if rest, ok := param.AsIterRest(); ok {
			if name, isName := rest.Target().AsName(); isName {
				return name.Ident()
			}
		} else if dv, ok := param.AsDefaultValue(); ok {
			if name, isName := dv.Target().AsName(); isName {
				return name.Ident()
			}
		} else if name, ok := param.AsName(); ok {
			return name.Ident()
		}
	}
	return "p" + strconv.Itoa(paramIndex)
}

// Java: TypedCodeGenerator.formatTypeVar(JSType)
func (t *CodeGenerator) FormatTypeVar(var_ rhino.JSType) string {
	return var_.ToAnnotationString(rhino.JSType_Nullability_IMPLICIT)
}

// TODO(dimvar): it's awkward that we print @constructor after the
// extends/implements; we should print it first, like users write it. Same
// for @interface and @record.
//
// Java: TypedCodeGenerator.appendClassAnnotations(StringBuilder, FunctionType)
func (t *CodeGenerator) AppendClassAnnotations(sb *strings.Builder, funType rhino.FunctionType) {
	superConstructor := funType.GetInstanceType().GetSuperClassConstructor()
	if superConstructor != nil {
		superInstance := superConstructor.GetInstanceType()
		if superInstance.String() != "Object" {
			sb.WriteString(" * ")
			TypedCodeGenerator_AppendAnnotation(sb, "extends", superInstance.ToAnnotationString(rhino.JSType_Nullability_IMPLICIT))
			sb.WriteString("\n")
		}
	}
	// Avoid duplicates, sort and deduplicate implemented interfaces.
	appendAncestorInterfaceAnnotations(sb, funType, "implements")
}

// Java: TypedCodeGenerator.appendInterfaceAnnotations(StringBuilder, FunctionType)
func (t *CodeGenerator) AppendInterfaceAnnotations(sb *strings.Builder, funType rhino.FunctionType) {
	appendAncestorInterfaceAnnotations(sb, funType, "extends")
	if funType.IsStructuralInterface() {
		sb.WriteString(" * @record\n")
	} else {
		sb.WriteString(" * @interface\n")
	}
}

func appendAncestorInterfaceAnnotations(sb *strings.Builder, funType rhino.FunctionType, tag string) {
	ancestors := funType.GetAncestorInterfaces()
	if len(ancestors) == 0 {
		return
	}
	interfaces := make([]string, len(ancestors))
	for i, interfaze := range ancestors {
		interfaces[i] = interfaze.ToAnnotationString(rhino.JSType_Nullability_IMPLICIT)
	}
	slices.SortFunc(interfaces, jutil.StringCompareTo)
	interfaces = slices.Compact(interfaces)
	for _, interfaze := range interfaces {
		sb.WriteString(" * ")
		TypedCodeGenerator_AppendAnnotation(sb, tag, interfaze)
		sb.WriteString("\n")
	}
}

// TODO(sdh): This whole method could be deleted if we don't mind adding
// additional @this annotations where they're not actually necessary.

// FindMethodOwner: given a method definition node, returns the ObjectType
// corresponding to the class the method is defined on, or null if it is not
// a prototype method.
//
// Java: TypedCodeGenerator.findMethodOwner(Node)
func (t *CodeGenerator) FindMethodOwner(n rhino.Node) rhino.ObjectType {
	if n.IsNil() {
		return nil
	}
	parent := n.GetParent()
	var ctor rhino.FunctionType
	if assign, ok := ast.AsAssign(parent); ok {
		if gp, isGp := assign.Target().AsGetProp(); isGp && NodeUtil_IsPrototypeProperty(gp.Node()) {
			// TODO(johnlenz): handle non-global types
			protoObj := ast.MustGetProp(gp.Object().Node()).Object().Node()
			if qname, hasQName := protoObj.GetQualifiedNameString(); hasQName {
				if type_ := t.registry.GetGlobalType(qname); type_ != nil {
					if objType := type_.ToMaybeObjectType(); objType != nil {
						ctor = objType.GetConstructor()
					}
				}
			}
		}
	} else if cls, ok := ast.AsClass(parent); ok {
		// TODO(sdh): test this case once the type checker understands ES6
		// classes
		if clsType := cls.Node().GetJSType(); clsType != nil {
			ctor = clsType.ToMaybeFunctionType()
		}
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
func (t *CodeGenerator) GetParameterJSDocType(
	parameters []rhino.FunctionType_Parameter, index int, minArgs int, maxArgs int) string {
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
func (t *CodeGenerator) RestrictByUndefined(type_ rhino.JSType) rhino.JSType {
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
