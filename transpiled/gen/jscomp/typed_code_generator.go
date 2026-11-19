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
// Port of com/google/javascript/jscomp/TypedCodeGenerator.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.TypedCodeGenerator
type TypedCodeGenerator struct {
	CodeGeneratorBase
	registry                              *pkg_rhino.JSTypeRegistry `gson:"registry"`
	f_TypedCodeGenerator_jsDocInfoPrinter *JSDocInfoPrinter         `gson:"jsDocInfoPrinter"`
}

func (this_ *TypedCodeGenerator) BindSelf_(s any) {
	_ = s
	this_.CodeGeneratorBase.BindSelf_(s)
}
func (this_ *TypedCodeGenerator) IsJavaNil_() bool { return this_ == nil }
func (*TypedCodeGenerator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedCodeGenerator", "TypedCodeGenerator"
}
func (this_ *TypedCodeGenerator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodeGeneratorBase.IdentityHash_()
}
func (*TypedCodeGenerator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedCodeGenerator", "com.google.javascript.jscomp.TypedCodeGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedCodeGenerator](v_)
	})
}
func (*TypedCodeGenerator) Implements_jscomp_CodeGenerator() {}

func (this_ *TypedCodeGenerator) initFields_TypedCodeGenerator() {
	_ = this_
}

// Java constructor: TypedCodeGenerator(com.google.javascript.jscomp.CodeConsumer consumer, com.google.javascript.jscomp.CompilerOptions options, com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewTypedCodeGenerator(consumer CodeConsumer, options *CompilerOptions, registry *pkg_rhino.JSTypeRegistry) *TypedCodeGenerator {
	this_ := &TypedCodeGenerator{}
	this_.BindSelf_(this_)
	InitTypedCodeGenerator(this_, consumer, options, registry)
	return this_
}

func InitTypedCodeGenerator(this_ *TypedCodeGenerator, consumer CodeConsumer, options *CompilerOptions, registry *pkg_rhino.JSTypeRegistry) {
	_ = consumer
	_ = options
	_ = registry
	InitCodeGenerator_2(&this_.CodeGeneratorBase, consumer, options)
	this_.initFields_TypedCodeGenerator()
	j2gort.CheckNotNullPtr(registry)
	this_.registry = registry
	this_.f_TypedCodeGenerator_jsDocInfoPrinter = NewJSDocInfoPrinter_1(options.GetUseOriginalNamesInOutput())
}

// Java method: TypedCodeGenerator.add(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.CodeGenerator.Context context)
func (this_ *TypedCodeGenerator) Add_2(n pkg_rhino.Node, context CodeGenerator_Context) {
	_ = n
	_ = context
	this_.priv_TypedCodeGenerator_maybeAddTypeAnnotation_5695(n)
	this_.CodeGeneratorBase.Add_2(n, context)
}

// Java method: TypedCodeGenerator.maybeAddTypeAnnotation(com.google.javascript.rhino.Node n)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_maybeAddTypeAnnotation_5695(n pkg_rhino.Node) {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent == 0 {
		return
	}
	if (parent.IsBlock() || parent.IsScript()) || parent.IsClassMembers() {
		if (n.IsClass() || n.IsFunction()) || n.IsMemberFunctionDef() {
			this_.Add_1_string(this_.priv_TypedCodeGenerator_getTypeAnnotation_5696(n))
		} else if n.IsExprResult() && n.GetFirstChild().IsAssign() {
			var assign pkg_rhino.Node = n.GetFirstChild()
			_ = assign
			if NodeUtil_IsNamespaceDecl(assign.GetFirstChild()) {
				this_.Add_1_string(this_.f_TypedCodeGenerator_jsDocInfoPrinter.Print(assign.GetJSDocInfo()))
			} else {
				var rhs pkg_rhino.Node = assign.GetLastChild()
				_ = rhs
				this_.Add_1_string(this_.priv_TypedCodeGenerator_getTypeAnnotation_5696(rhs))
			}
		} else if NodeUtil_IsNameDeclaration(n) && (n.GetFirstFirstChild() != 0) {
			if NodeUtil_IsNamespaceDecl(n.GetFirstChild()) && (n.GetJSDocInfo() != nil) {
				this_.Add_1_string(this_.f_TypedCodeGenerator_jsDocInfoPrinter.Print(n.GetJSDocInfo()))
			} else {
				this_.Add_1_string(this_.priv_TypedCodeGenerator_getTypeAnnotation_5696(n.GetFirstFirstChild()))
			}
		}
	}
}

// Java method: TypedCodeGenerator.getTypeAnnotation(com.google.javascript.rhino.Node node)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_getTypeAnnotation_5696(node pkg_rhino.Node) string {
	_ = node
	if node.IsMemberFunctionDef() {
		return this_.priv_TypedCodeGenerator_getMemberFunctionAnnotation_5698(node.GetOnlyChild())
	} else if node.IsClass() {
		return this_.priv_TypedCodeGenerator_getClassAnnotation_5700(node.GetJSType())
	} else if node.IsFunction() {
		return this_.priv_TypedCodeGenerator_getFunctionAnnotation_5697(node)
	} else {
		var nodeOriginallyHadJSDoc bool = (NodeUtil_GetBestJSDocInfo(node) != nil)
		_ = nodeOriginallyHadJSDoc
		if !(nodeOriginallyHadJSDoc) {
			return ""
		}
		var type_ pkg_rhino.JSType = node.GetJSType()
		_ = type_
		if type_ == nil || type_.IsJavaNil_() {
			return ""
		} else if type_.GetJSTypeBase().IsFunctionType() {
			return this_.priv_TypedCodeGenerator_getFunctionAnnotation_5697(node)
		} else if type_.GetJSTypeBase().IsEnumType() {
			return ("/** @enum {" + j2gort.StrValString(type_.GetJSTypeBase().ToMaybeObjectType().GetEnumeratedTypeOfEnumObject().GetJSTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_EXPLICIT())) + "} */\n")
		} else if ((!(type_.GetJSTypeBase().IsUnknownType()) && !(type_.GetJSTypeBase().IsEmptyType())) && !(type_.GetJSTypeBase().IsVoidType())) && !(type_.IsFunctionPrototypeType()) {
			return ("/** @type {" + j2gort.StrValString(node.GetJSType().GetJSTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_EXPLICIT())) + "} */\n")
		} else {
			return ""
		}
	}
}

// Java method: TypedCodeGenerator.getFunctionAnnotation(com.google.javascript.rhino.Node fnNode)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_getFunctionAnnotation_5697(fnNode pkg_rhino.Node) string {
	_ = fnNode
	var type_ pkg_rhino.JSType = fnNode.GetJSType()
	_ = type_
	if !(fnNode.IsFunction() || type_.GetJSTypeBase().IsFunctionType()) {
		j2gort.CheckState(false)
	}
	if (type_ == nil || type_.IsJavaNil_()) || type_.GetJSTypeBase().IsUnknownType() {
		return ""
	}
	var funType *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
	_ = funType
	if type_.GetJSTypeBase().EqualsJSType(this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_TYPE())) {
		return "/** @type {!Function} */\n"
	}
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder_1_string("/**\n")
	_ = sb
	var paramNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = paramNode
	if (fnNode != 0) && fnNode.IsFunction() {
		paramNode = NodeUtil_GetFunctionParameters(fnNode).GetFirstChild()
	}
	this_.priv_TypedCodeGenerator_appendFunctionParamAnnotations_5699(sb, funType, paramNode)
	var retType pkg_rhino.JSType = funType.GetReturnType()
	_ = retType
	if (((retType != nil && !retType.IsJavaNil_()) && !(retType.GetJSTypeBase().IsEmptyType())) && !(funType.IsInterface())) && !(funType.IsConstructor() && retType.GetJSTypeBase().IsVoidType()) {
		sb.Append_1_string(" * ")
		TypedCodeGenerator_AppendAnnotation(sb, "return", retType.GetJSTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_EXPLICIT()))
		sb.Append_1_string("\n")
	}
	if funType.IsConstructor() {
		this_.priv_TypedCodeGenerator_appendClassAnnotations_5704(sb, funType)
		sb.Append_1_string(" * @constructor\n")
	} else if funType.IsInterface() {
		this_.priv_TypedCodeGenerator_appendInterfaceAnnotations_5705(sb, funType)
	} else {
		var thisType pkg_rhino.JSType = funType.GetTypeOfThis()
		_ = thisType
		if ((thisType != nil && !thisType.IsJavaNil_()) && !(thisType.GetJSTypeBase().IsUnknownType())) && !(thisType.GetJSTypeBase().IsVoidType()) {
			if (fnNode == 0) || !(thisType.GetJSTypeBase().EqualsJSType(this_.priv_TypedCodeGenerator_findMethodOwner_5706(fnNode))) {
				sb.Append_1_string(" * ")
				TypedCodeGenerator_AppendAnnotation(sb, "this", thisType.GetJSTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_EXPLICIT()))
				sb.Append_1_string("\n")
			}
		}
	}
	this_.priv_TypedCodeGenerator_appendTemplateAnnotations_5701(sb, j2gort.AsIface[j2gort.Ext_java_util_Collection](funType.GetTypeParameters()))
	sb.Append_1_string(" */\n")
	return sb.String()
}

// Java method: TypedCodeGenerator.getMemberFunctionAnnotation(com.google.javascript.rhino.Node fnNode)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_getMemberFunctionAnnotation_5698(fnNode pkg_rhino.Node) string {
	_ = fnNode
	if !(fnNode.IsFunction() && fnNode.GetParent().IsMemberFunctionDef()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(fnNode))
	}
	var type_ pkg_rhino.JSType = fnNode.GetJSType()
	_ = type_
	if (type_ == nil || type_.IsJavaNil_()) || type_.GetJSTypeBase().IsUnknownType() {
		return ""
	}
	var funType *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
	_ = funType
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder_1_string("/**\n")
	_ = sb
	var paramNode pkg_rhino.Node = NodeUtil_GetFunctionParameters(fnNode).GetFirstChild()
	_ = paramNode
	this_.priv_TypedCodeGenerator_appendFunctionParamAnnotations_5699(sb, funType, paramNode)
	if NodeUtil_IsEs6Constructor(fnNode) {
		this_.priv_TypedCodeGenerator_appendTemplateAnnotations_5701(sb, j2gort.AsIface[j2gort.Ext_java_util_Collection](funType.GetConstructorOnlyTemplateParameters()))
	} else {
		this_.priv_TypedCodeGenerator_appendTemplateAnnotations_5701(sb, j2gort.AsIface[j2gort.Ext_java_util_Collection](funType.GetTypeParameters()))
		var retType pkg_rhino.JSType = funType.GetReturnType()
		_ = retType
		if (retType != nil && !retType.IsJavaNil_()) && !(retType.GetJSTypeBase().IsEmptyType()) {
			sb.Append_1_string(" * ")
			TypedCodeGenerator_AppendAnnotation(sb, "return", retType.GetJSTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_EXPLICIT()))
			sb.Append_1_string("\n")
		}
	}
	sb.Append_1_string(" */\n")
	return sb.String()
}

// Java method: TypedCodeGenerator.appendFunctionParamAnnotations(java.lang.StringBuilder sb, com.google.javascript.rhino.jstype.FunctionType funType, com.google.javascript.rhino.Node paramNode)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_appendFunctionParamAnnotations_5699(sb *j2gort.Ext_java_lang_StringBuilder, funType *pkg_rhino.FunctionType, paramNode pkg_rhino.Node) {
	_ = sb
	_ = funType
	_ = paramNode
	var minArity int32 = funType.GetMinArity()
	_ = minArity
	var maxArity int32 = funType.GetMaxArity()
	_ = maxArity
	var formals []pkg_rhino.FunctionType_Parameter = funType.GetParameterArray()
	_ = formals
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(formals)); i = i + 1 {
			sb.Append_1_string(" * ")
			TypedCodeGenerator_AppendAnnotation(sb, "param", this_.priv_TypedCodeGenerator_getParameterJSDocType_5707(formals, i, minArity, maxArity))
			var parameterName string = this_.priv_TypedCodeGenerator_getParameterJSDocName_5702(paramNode, i)
			_ = parameterName
			j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string(" ")).Append_1_string(parameterName)).Append_1_string("\n")
			if paramNode != 0 {
				paramNode = paramNode.GetNext()
			}
		}
	}
}

// Java method: TypedCodeGenerator.getClassAnnotation(com.google.javascript.rhino.jstype.JSType classType)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_getClassAnnotation_5700(classType pkg_rhino.JSType) string {
	_ = classType
	if (classType == nil || classType.IsJavaNil_()) || classType.GetJSTypeBase().IsUnknownType() {
		return ""
	}
	if !(classType.GetJSTypeBase().IsFunctionType()) {
		j2gort.CheckStateMsg(false, classType)
	}
	var funType *pkg_rhino.FunctionType = classType.ToMaybeFunctionType()
	_ = funType
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	if funType.IsInterface() {
		this_.priv_TypedCodeGenerator_appendInterfaceAnnotations_5705(sb, funType)
	} else {
		if !(funType.IsConstructor()) {
			j2gort.CheckStateMsg(false, j2gort.PtrToAny(funType))
		}
		this_.priv_TypedCodeGenerator_appendClassAnnotations_5704(sb, funType)
	}
	this_.priv_TypedCodeGenerator_appendTemplateAnnotations_5701(sb, j2gort.AsIface[j2gort.Ext_java_util_Collection](funType.GetTypeParameters()))
	var jsdocContent string = sb.String()
	_ = jsdocContent
	if j2gort.CheckNotNullString(jsdocContent) == "" {
		return jsdocContent
	} else {
		return ("/**\n" + j2gort.StrValString(jsdocContent) + " */\n")
	}
}

// Java method: TypedCodeGenerator.appendTemplateAnnotations(java.lang.StringBuilder sb, java.util.Collection<? extends com.google.javascript.rhino.jstype.JSType> typeParams)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_appendTemplateAnnotations_5701(sb *j2gort.Ext_java_lang_StringBuilder, typeParams j2gort.Ext_java_util_Collection) {
	_ = sb
	_ = typeParams
	if !(typeParams.IsEmpty()) {
		sb.Append_1_string(" * @template ")
		sb.Append_1_string(j2gort.CastString(typeParams.Stream().Map(func() j2gort.Ext_java_util_function_Function {
			bound_1 := this_
			return j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
				return j2gort.BoxString(bound_1.priv_TypedCodeGenerator_formatTypeVar_5703(j2gort.AsIface[pkg_rhino.JSType](mrp0_)))
			})
		}()).Collect_1(j2gort.Ext_java_util_stream_Collectors_Joining_1(","))))
		sb.Append_1_string("\n")
	}
}

// Java method: TypedCodeGenerator.getParameterJSDocName(com.google.javascript.rhino.Node paramNode, int paramIndex)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_getParameterJSDocName_5702(paramNode pkg_rhino.Node, paramIndex int32) string {
	_ = paramNode
	_ = paramIndex
	var nameNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = nameNode
	if paramNode != 0 {
		if !(paramNode.GetParent().IsParamList()) {
			j2gort.CheckArgumentMsg(false, j2gort.BoxNode(paramNode))
		}
		if paramNode.IsRest() {
			paramNode = paramNode.GetOnlyChild()
		} else if paramNode.IsDefaultValue() {
			paramNode = paramNode.GetFirstChild()
		}
		if paramNode.IsName() {
			nameNode = paramNode
		} else {
			if !(paramNode.IsObjectPattern() || paramNode.IsArrayPattern()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(paramNode))
			}
			nameNode = pkg_rhino.Node(0)
		}
	}
	if nameNode == 0 {
		return ("p" + j2gort.StrValInt32(paramIndex))
	} else {
		if !(nameNode.IsName()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(nameNode))
		}
		return nameNode.GetString()
	}
}

// Java method: TypedCodeGenerator.formatTypeVar(com.google.javascript.rhino.jstype.JSType var)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_formatTypeVar_5703(var_ pkg_rhino.JSType) string {
	_ = var_
	return var_.GetJSTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_IMPLICIT())
}

// Java method: TypedCodeGenerator.appendClassAnnotations(java.lang.StringBuilder sb, com.google.javascript.rhino.jstype.FunctionType funType)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_appendClassAnnotations_5704(sb *j2gort.Ext_java_lang_StringBuilder, funType *pkg_rhino.FunctionType) {
	_ = sb
	_ = funType
	var superConstructor *pkg_rhino.FunctionType = funType.GetInstanceType().GetSuperClassConstructor()
	_ = superConstructor
	if superConstructor != nil {
		var superInstance pkg_rhino.ObjectType = superConstructor.GetInstanceType()
		_ = superInstance
		if !(j2gort.CheckNotNullString(superInstance.GetObjectTypeBase().String()) == "Object") {
			sb.Append_1_string(" * ")
			TypedCodeGenerator_AppendAnnotation(sb, "extends", superInstance.GetObjectTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_IMPLICIT()))
			sb.Append_1_string("\n")
		}
	}
	var interfaces j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet()
	_ = interfaces
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(funType.GetAncestorInterfaces()))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			interfaze := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_3)
			_ = interfaze
			j2gort.SetAddStr(interfaces, interfaze.GetObjectTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_IMPLICIT()))
		}
	}
	{
		cur_4 := j2gort.NewIterCursor(interfaces)
		for cur_4.HasNext() {
			interfaze_2 := cur_4.NextStr()
			_ = interfaze_2
			sb.Append_1_string(" * ")
			TypedCodeGenerator_AppendAnnotation(sb, "implements", interfaze_2)
			sb.Append_1_string("\n")
		}
	}
}

// Java method: TypedCodeGenerator.appendInterfaceAnnotations(java.lang.StringBuilder sb, com.google.javascript.rhino.jstype.FunctionType funType)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_appendInterfaceAnnotations_5705(sb *j2gort.Ext_java_lang_StringBuilder, funType *pkg_rhino.FunctionType) {
	_ = sb
	_ = funType
	var interfaces j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet()
	_ = interfaces
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(funType.GetAncestorInterfaces()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			interfaceType := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_6)
			_ = interfaceType
			j2gort.SetAddStr(interfaces, interfaceType.GetObjectTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_IMPLICIT()))
		}
	}
	{
		cur_7 := j2gort.NewIterCursor(interfaces)
		for cur_7.HasNext() {
			interfaze := cur_7.NextStr()
			_ = interfaze
			sb.Append_1_string(" * ")
			TypedCodeGenerator_AppendAnnotation(sb, "extends", interfaze)
			sb.Append_1_string("\n")
		}
	}
	if funType.IsStructuralInterface() {
		sb.Append_1_string(" * @record\n")
	} else {
		sb.Append_1_string(" * @interface\n")
	}
}

// Java method: TypedCodeGenerator.findMethodOwner(com.google.javascript.rhino.Node n)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_findMethodOwner_5706(n pkg_rhino.Node) pkg_rhino.ObjectType {
	_ = n
	if n == 0 {
		return nil
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var ctor *pkg_rhino.FunctionType = nil
	_ = ctor
	if parent.IsAssign() {
		var target pkg_rhino.Node = parent.GetFirstChild()
		_ = target
		if NodeUtil_IsPrototypeProperty(target) {
			var type_ pkg_rhino.JSType = this_.registry.GetGlobalType(target.GetFirstFirstChild().GetQualifiedName())
			_ = type_
			if type_ != nil && !type_.IsJavaNil_() {
				ctor = (j2gort.AsIface[pkg_rhino.ObjectType](type_)).GetConstructor()
			} else {
				ctor = nil
			}
		}
	} else if parent.IsClass() {
		ctor = parent.GetJSType().ToMaybeFunctionType()
	}
	if ctor != nil {
		return ctor.GetInstanceType()
	} else {
		return nil
	}
}

// Java static method: TypedCodeGenerator.appendAnnotation(java.lang.StringBuilder sb, java.lang.String name, java.lang.String type)
func TypedCodeGenerator_AppendAnnotation(sb *j2gort.Ext_java_lang_StringBuilder, name string, type_ string) {
	_ = sb
	_ = name
	_ = type_
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string("@")).Append_1_string(name)).Append_1_string(" {")).Append_1_string(type_)).Append_1_string("}")
}

// Java method: TypedCodeGenerator.getParameterJSDocType(com.google.javascript.rhino.jstype.FunctionType.Parameter[] parameters, int index, int minArgs, int maxArgs)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_getParameterJSDocType_5707(parameters []pkg_rhino.FunctionType_Parameter, index int32, minArgs int32, maxArgs int32) string {
	_ = parameters
	_ = index
	_ = minArgs
	_ = maxArgs
	var type_ pkg_rhino.JSType = parameters[index].GetJSType()
	_ = type_
	if index < minArgs {
		return type_.GetJSTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_EXPLICIT())
	}
	var isRestArgument bool = ((maxArgs == int32(2147483647)) && (index == (int32(len(parameters)) - int32(1))))
	_ = isRestArgument
	if isRestArgument {
		return ("..." + j2gort.StrValString(this_.priv_TypedCodeGenerator_restrictByUndefined_5708(type_).GetJSTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_EXPLICIT())))
	}
	return (j2gort.StrValString(this_.priv_TypedCodeGenerator_restrictByUndefined_5708(type_).GetJSTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_EXPLICIT())) + "=")
}

// Java method: TypedCodeGenerator.restrictByUndefined(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypedCodeGenerator) priv_TypedCodeGenerator_restrictByUndefined_5708(type_ pkg_rhino.JSType) pkg_rhino.JSType {
	_ = type_
	if !(type_.IsVoidable()) {
		return type_
	}
	var restricted pkg_rhino.JSType = type_.RestrictByNotNullOrUndefined()
	_ = restricted
	if type_.IsNullable() {
		var nullType pkg_rhino.JSType = this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NULL_TYPE())
		_ = nullType
		return this_.registry.CreateUnionType_1_Ext_java_util_List(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(restricted, nullType)))
	}
	if restricted.GetJSTypeBase().IsEmptyType() {
		return type_
	} else {
		return restricted
	}
}
