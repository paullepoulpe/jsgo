/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/FunctionTypeBuilder.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.FunctionTypeBuilder
type FunctionTypeBuilder struct {
	fnName                       string                                             `gson:"fnName"`
	compiler                     AbstractCompiler                                   `gson:"compiler"`
	codingConvention             CodingConvention                                   `gson:"codingConvention"`
	typeRegistry                 *pkg_rhino.JSTypeRegistry                          `gson:"typeRegistry"`
	contents                     FunctionTypeBuilder_FunctionContents               `gson:"contents"`
	syntacticFnName              string                                             `gson:"syntacticFnName"`
	returnType                   pkg_rhino.JSType                                   `gson:"returnType"`
	implementedInterfaces        j2gort.Ext_java_util_List                          `gson:"implementedInterfaces"`
	extendedInterfaces           j2gort.Ext_java_util_List                          `gson:"extendedInterfaces"`
	baseType                     pkg_rhino.ObjectType                               `gson:"baseType"`
	thisType                     pkg_rhino.JSType                                   `gson:"thisType"`
	parameters                   []pkg_rhino.FunctionType_Parameter                 `gson:"parameters"`
	templateTypeNames            j2gort.Ext_com_google_common_collect_ImmutableList `gson:"templateTypeNames"`
	constructorTemplateTypeNames j2gort.Ext_com_google_common_collect_ImmutableList `gson:"constructorTemplateTypeNames"`
	declarationScope             *TypedScope                                        `gson:"declarationScope"`
	templateScope                pkg_rhino.StaticTypedScope                         `gson:"templateScope"`
	errorRoot                    pkg_rhino.Node                                     `gson:"errorRoot"`
	idHash_                      atomic.Int32
	returnTypeInferred           bool                       `gson:"returnTypeInferred"`
	isClass                      bool                       `gson:"isClass"`
	isConstructor                bool                       `gson:"isConstructor"`
	makesStructs                 bool                       `gson:"makesStructs"`
	makesUnrestricted            bool                       `gson:"makesUnrestricted"`
	makesDicts                   bool                       `gson:"makesDicts"`
	isInterface                  bool                       `gson:"isInterface"`
	isRecord                     bool                       `gson:"isRecord"`
	isAbstract                   bool                       `gson:"isAbstract"`
	isKnownAmbiguous             bool                       `gson:"isKnownAmbiguous"`
	closurePrimitiveId           pkg_rhino.ClosurePrimitive `gson:"closurePrimitiveId"`
}

func (this_ *FunctionTypeBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionTypeBuilder) IsJavaNil_() bool { return this_ == nil }
func (*FunctionTypeBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionTypeBuilder", "FunctionTypeBuilder"
}
func (this_ *FunctionTypeBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionTypeBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionTypeBuilder", "com.google.javascript.jscomp.FunctionTypeBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionTypeBuilder](v_)
	})
}

var static_FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF *DiagnosticType

func StaticGet_FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF() *DiagnosticType {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF
}
func StaticSet_FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF(v *DiagnosticType) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF = v
}

var static_FunctionTypeBuilder_EXTENDS_NON_OBJECT *DiagnosticType

func StaticGet_FunctionTypeBuilder_EXTENDS_NON_OBJECT() *DiagnosticType {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_EXTENDS_NON_OBJECT
}
func StaticSet_FunctionTypeBuilder_EXTENDS_NON_OBJECT(v *DiagnosticType) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_EXTENDS_NON_OBJECT = v
}

var static_FunctionTypeBuilder_RESOLVED_TAG_EMPTY *DiagnosticType

func StaticGet_FunctionTypeBuilder_RESOLVED_TAG_EMPTY() *DiagnosticType {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_RESOLVED_TAG_EMPTY
}
func StaticSet_FunctionTypeBuilder_RESOLVED_TAG_EMPTY(v *DiagnosticType) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_RESOLVED_TAG_EMPTY = v
}

var static_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED *DiagnosticType

func StaticGet_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED() *DiagnosticType {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED
}
func StaticSet_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED(v *DiagnosticType) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED = v
}

var static_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST *DiagnosticType

func StaticGet_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST() *DiagnosticType {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST
}
func StaticSet_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST(v *DiagnosticType) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST = v
}

var static_FunctionTypeBuilder_OPTIONAL_ARG_AT_END *DiagnosticType

func StaticGet_FunctionTypeBuilder_OPTIONAL_ARG_AT_END() *DiagnosticType {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_OPTIONAL_ARG_AT_END
}
func StaticSet_FunctionTypeBuilder_OPTIONAL_ARG_AT_END(v *DiagnosticType) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_OPTIONAL_ARG_AT_END = v
}

var static_FunctionTypeBuilder_INEXISTENT_PARAM *DiagnosticType

func StaticGet_FunctionTypeBuilder_INEXISTENT_PARAM() *DiagnosticType {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_INEXISTENT_PARAM
}
func StaticSet_FunctionTypeBuilder_INEXISTENT_PARAM(v *DiagnosticType) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_INEXISTENT_PARAM = v
}

var static_FunctionTypeBuilder_TYPE_REDEFINITION *DiagnosticType

func StaticGet_FunctionTypeBuilder_TYPE_REDEFINITION() *DiagnosticType {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_TYPE_REDEFINITION
}
func StaticSet_FunctionTypeBuilder_TYPE_REDEFINITION(v *DiagnosticType) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_TYPE_REDEFINITION = v
}

var static_FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS *DiagnosticType

func StaticGet_FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS() *DiagnosticType {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS
}
func StaticSet_FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS(v *DiagnosticType) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS = v
}

var static_FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND *DiagnosticType

func StaticGet_FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND() *DiagnosticType {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND
}
func StaticSet_FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND(v *DiagnosticType) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND = v
}

var static_FunctionTypeBuilder_ALL_DIAGNOSTICS *DiagnosticGroup

func StaticGet_FunctionTypeBuilder_ALL_DIAGNOSTICS() *DiagnosticGroup {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_ALL_DIAGNOSTICS
}
func StaticSet_FunctionTypeBuilder_ALL_DIAGNOSTICS(v *DiagnosticGroup) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_ALL_DIAGNOSTICS = v
}

var clinitDone_FunctionTypeBuilder uint32
var clinitState_FunctionTypeBuilder j2gort.ClinitState

func clinit_FunctionTypeBuilder() {
	if atomic.LoadUint32(&clinitDone_FunctionTypeBuilder) != 0 {
		return
	}
	clinitSlow_FunctionTypeBuilder()
}

func clinitSlow_FunctionTypeBuilder() {
	if !j2gort.EnterClinit(&clinitState_FunctionTypeBuilder, &clinitDone_FunctionTypeBuilder) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FunctionTypeBuilder, &clinitDone_FunctionTypeBuilder)
	static_FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF = DiagnosticType_Warning("JSC_EXTENDS_WITHOUT_TYPEDEF", "@extends used without @constructor or @interface for {0}")
	static_FunctionTypeBuilder_EXTENDS_NON_OBJECT = DiagnosticType_Warning("JSC_EXTENDS_NON_OBJECT", "{0} @extends non-object type {1}")
	static_FunctionTypeBuilder_RESOLVED_TAG_EMPTY = DiagnosticType_Warning("JSC_RESOLVED_TAG_EMPTY", "Could not resolve type in {0} tag of {1}")
	static_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED = DiagnosticType_Warning("JSC_CONSTRUCTOR_REQUIRED", "{0} used without @constructor for {1}")
	static_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST = DiagnosticType_Warning("JSC_VAR_ARGS_MUST_BE_LAST", "variable length argument must be last")
	static_FunctionTypeBuilder_OPTIONAL_ARG_AT_END = DiagnosticType_Warning("JSC_OPTIONAL_ARG_AT_END", "optional arguments must be at the end")
	static_FunctionTypeBuilder_INEXISTENT_PARAM = DiagnosticType_Warning("JSC_INEXISTENT_PARAM", "parameter {0} does not appear in {1}''s parameter list")
	static_FunctionTypeBuilder_TYPE_REDEFINITION = DiagnosticType_Warning("JSC_TYPE_REDEFINITION", "attempted re-definition of type {0}\nfound   : {1}\nexpected: {2}")
	static_FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS = DiagnosticType_Warning("JSC_TEMPLATE_TRANSFORMATION_ON_CLASS", "Template type transformation {0} not allowed on classes or interfaces")
	static_FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND = DiagnosticType_Error("JSC_TEMPLATE_TYPE_ILLEGAL_BOUND", "Illegal upper bound ''{0}'' on template type parameter {1}")
	static_FunctionTypeBuilder_ALL_DIAGNOSTICS = NewDiagnosticGroup_1_DiagnosticType_2([]*DiagnosticType{StaticGet_FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF(), StaticGet_FunctionTypeBuilder_EXTENDS_NON_OBJECT(), StaticGet_FunctionTypeBuilder_RESOLVED_TAG_EMPTY(), StaticGet_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED(), StaticGet_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST(), StaticGet_FunctionTypeBuilder_OPTIONAL_ARG_AT_END(), StaticGet_FunctionTypeBuilder_INEXISTENT_PARAM(), StaticGet_FunctionTypeBuilder_TYPE_REDEFINITION(), StaticGet_FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS(), StaticGet_FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND(), StaticGet_TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS()})
}

func (this_ *FunctionTypeBuilder) initFields_FunctionTypeBuilder() {
	_ = this_
	this_.fnName = j2gort.NullString
	this_.contents = FunctionTypeBuilder_UnknownFunctionContents_Get()
	this_.syntacticFnName = j2gort.NullString

	this_.templateTypeNames = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	this_.constructorTemplateTypeNames = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()

}

// Java constructor: FunctionTypeBuilder(java.lang.String fnName, com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node errorRoot, com.google.javascript.jscomp.TypedScope scope)
func NewFunctionTypeBuilder(fnName string, compiler AbstractCompiler, errorRoot pkg_rhino.Node, scope *TypedScope) *FunctionTypeBuilder {
	clinit_FunctionTypeBuilder()
	this_ := &FunctionTypeBuilder{}
	InitFunctionTypeBuilder(this_, fnName, compiler, errorRoot, scope)
	return this_
}

func InitFunctionTypeBuilder(this_ *FunctionTypeBuilder, fnName string, compiler AbstractCompiler, errorRoot pkg_rhino.Node, scope *TypedScope) {
	_ = fnName
	_ = compiler
	_ = errorRoot
	_ = scope
	this_.initFields_FunctionTypeBuilder()
	j2gort.CheckNotNullNode(errorRoot)
	this_.fnName = j2gort.Ext_com_google_common_base_Strings_NullToEmpty(fnName)
	this_.codingConvention = compiler.GetCodingConvention()
	this_.typeRegistry = compiler.GetTypeRegistry()
	this_.errorRoot = errorRoot
	this_.compiler = compiler
	this_.templateScope = pkg_rhino.ToIface_StaticTypedScope(scope)
}

// Java method: FunctionTypeBuilder.createExtendedTypeValidator()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_createExtendedTypeValidator_2960() *FunctionTypeBuilder_ExtendedTypeValidator {
	return NewFunctionTypeBuilder_ExtendedTypeValidator(this_.errorRoot, this_.compiler, this_.FormatFnName())
}

// Java method: FunctionTypeBuilder.createImplementedTypeValidator()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_createImplementedTypeValidator_2961() *FunctionTypeBuilder_ImplementedTypeValidator {
	return NewFunctionTypeBuilder_ImplementedTypeValidator(this_.errorRoot, this_.compiler, this_.FormatFnName())
}

// Java method: FunctionTypeBuilder.formatFnName()
func (this_ *FunctionTypeBuilder) FormatFnName() string {
	if j2gort.CheckNotNullString(this_.fnName) == "" {
		return "<anonymous>"
	} else {
		return this_.fnName
	}
}

// Java method: FunctionTypeBuilder.setSyntacticFunctionName(java.lang.String syntacticFnName)
func (this_ *FunctionTypeBuilder) SetSyntacticFunctionName(syntacticFnName string) *FunctionTypeBuilder {
	_ = syntacticFnName
	this_.syntacticFnName = j2gort.Ext_com_google_common_base_Strings_NullToEmpty(syntacticFnName)
	return this_
}

// Java method: FunctionTypeBuilder.setContents(com.google.javascript.jscomp.FunctionTypeBuilder.@org.jspecify.annotations.Nullable FunctionContents contents)
func (this_ *FunctionTypeBuilder) SetContents(contents FunctionTypeBuilder_FunctionContents) *FunctionTypeBuilder {
	_ = contents
	if contents != nil && !contents.IsJavaNil_() {
		this_.contents = contents
	}
	return this_
}

// Java method: FunctionTypeBuilder.setDeclarationScope(com.google.javascript.jscomp.TypedScope declarationScope)
func (this_ *FunctionTypeBuilder) SetDeclarationScope(declarationScope *TypedScope) *FunctionTypeBuilder {
	_ = declarationScope
	this_.declarationScope = declarationScope
	return this_
}

// Java method: FunctionTypeBuilder.inferFromOverriddenFunction(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable FunctionType oldType, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node paramsParent)
func (this_ *FunctionTypeBuilder) InferFromOverriddenFunction(oldType *pkg_rhino.FunctionType, paramsParent pkg_rhino.Node) *FunctionTypeBuilder {
	_ = oldType
	_ = paramsParent
	if oldType == nil {
		return this_
	}
	this_.templateTypeNames = oldType.GetTemplateTypeMap().GetTemplateKeys()
	this_.returnType = oldType.GetReturnType()
	this_.returnTypeInferred = oldType.IsReturnTypeInferred()
	if paramsParent == 0 {
		this_.parameters = oldType.GetParameterArray()
		if this_.parameters == nil {
			this_.parameters = pkg_rhino.NewFunctionParamBuilder_1(this_.typeRegistry).BuildArray()
		}
	} else {
		var paramBuilder *pkg_rhino.FunctionParamBuilder = pkg_rhino.NewFunctionParamBuilder_1(this_.typeRegistry)
		_ = paramBuilder
		var oldParams []pkg_rhino.FunctionType_Parameter = oldType.GetParameterArray()
		_ = oldParams
		var oldParamsLen int32 = int32(len(oldParams))
		_ = oldParamsLen
		var oldParamsIdx int32 = int32(0)
		_ = oldParamsIdx
		var warnedAboutArgList bool = false
		_ = warnedAboutArgList
		var oldParamsListHitOptArgs bool = false
		_ = oldParamsListHitOptArgs
		{
			var currentParam pkg_rhino.Node = paramsParent.GetFirstChild()
			_ = currentParam
			for ; currentParam != 0; currentParam = currentParam.GetNext() {
				if oldParamsIdx < oldParamsLen {
					var oldParam pkg_rhino.FunctionType_Parameter = oldParams[func() int32 { old_1 := oldParamsIdx; nw_2 := old_1 + 1; oldParamsIdx = nw_2; return old_1 }()]
					_ = oldParam
					oldParamsListHitOptArgs = ((oldParamsListHitOptArgs || oldParam.IsVariadic()) || oldParam.IsOptional())
					var isOptionalArg bool = oldParam.IsOptional()
					_ = isOptionalArg
					var isVarArgs bool = oldParam.IsVariadic()
					_ = isVarArgs
					if (currentParam.GetNext() != 0) && isVarArgs {
						isVarArgs = false
						isOptionalArg = true
					}
					if currentParam.IsDefaultValue() {
						isOptionalArg = true
					}
					paramBuilder.NewParameterFrom(pkg_rhino.FunctionType_Parameter_Create(oldParam.GetJSType(), isOptionalArg, isVarArgs))
				} else {
					warnedAboutArgList = func() bool {
						ba_3 := warnedAboutArgList
						bb_4 := this_.priv_FunctionTypeBuilder_addParameter_2986(paramBuilder, this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE()), warnedAboutArgList, ((this_.codingConvention.IsOptionalParameter(currentParam) || oldParamsListHitOptArgs) || currentParam.IsDefaultValue()), this_.codingConvention.IsVarArgsParameter(currentParam))
						return ba_3 || bb_4
					}()
				}
			}
		}
		for oldParamsIdx < oldParamsLen {
			paramBuilder.NewOptionalParameterFrom(oldParams[func() int32 { old_5 := oldParamsIdx; nw_6 := old_5 + 1; oldParamsIdx = nw_6; return old_5 }()])
		}
		this_.parameters = paramBuilder.BuildArray()
	}
	return this_
}

// Java method: FunctionTypeBuilder.inferReturnType(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, boolean fromInlineDoc)
func (this_ *FunctionTypeBuilder) InferReturnType(info *pkg_rhino.JSDocInfo, fromInlineDoc bool) *FunctionTypeBuilder {
	_ = info
	_ = fromInlineDoc
	if info != nil {
		var returnTypeExpr pkg_rhino.JSTypeExpression
		_ = returnTypeExpr
		if fromInlineDoc {
			returnTypeExpr = info.GetType()
		} else {
			returnTypeExpr = info.GetReturnType()
		}
		if returnTypeExpr != (pkg_rhino.JSTypeExpression{}) {
			this_.returnType = returnTypeExpr.Evaluate(this_.templateScope, this_.typeRegistry)
			this_.returnTypeInferred = false
		}
	}
	return this_
}

// Java method: FunctionTypeBuilder.usingClassSyntax()
func (this_ *FunctionTypeBuilder) UsingClassSyntax() *FunctionTypeBuilder {
	this_.isClass = true
	return this_
}

// Java method: FunctionTypeBuilder.inferKind(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info)
func (this_ *FunctionTypeBuilder) InferKind(info *pkg_rhino.JSDocInfo) *FunctionTypeBuilder {
	_ = info
	if info != nil {
		if !(NodeUtil_IsMethodDeclaration(this_.errorRoot)) {
			this_.isConstructor = info.IsConstructor()
			this_.isInterface = info.IsInterface()
			this_.isRecord = info.UsesImplicitMatch()
			this_.makesStructs = info.MakesStructs()
			this_.makesUnrestricted = info.MakesUnrestricted()
			this_.makesDicts = info.MakesDicts()
		}
		this_.isAbstract = info.IsAbstract()
	}
	if this_.isClass {
		this_.isConstructor = !(this_.isInterface)
		this_.makesStructs = ((info == nil) || (!(this_.makesDicts) && !(info.MakesUnrestricted())))
	}
	if this_.makesStructs && !(this_.isConstructor || this_.isInterface) {
		this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED(), []string{"@struct", this_.FormatFnName()})
	} else if this_.makesDicts && !(this_.isConstructor) {
		this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED(), []string{"@dict", this_.FormatFnName()})
	}
	return this_
}

// Java method: FunctionTypeBuilder.maybeUseNativeClassTemplateNames(com.google.javascript.rhino.JSDocInfo info)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_maybeUseNativeClassTemplateNames_2970(info *pkg_rhino.JSDocInfo) bool {
	_ = info
	var nativeKeys j2gort.Ext_com_google_common_collect_ImmutableList = this_.typeRegistry.MaybeGetTemplateTypesOfBuiltin(pkg_rhino.ToIface_StaticScope(this_.declarationScope), this_.fnName)
	_ = nativeKeys
	if (nativeKeys != nil) && (info.GetTemplateTypeNames().Size() == nativeKeys.Size()) {
		this_.templateTypeNames = nativeKeys
		return true
	}
	return false
}

// Java method: FunctionTypeBuilder.inferInheritance(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType classExtendsType)
func (this_ *FunctionTypeBuilder) InferInheritance(info *pkg_rhino.JSDocInfo, classExtendsType pkg_rhino.ObjectType) *FunctionTypeBuilder {
	_ = info
	_ = classExtendsType
	if (info != nil) && info.HasBaseType() {
		if this_.isConstructor || this_.isInterface {
			var infoBaseType pkg_rhino.JSType = info.GetBaseType().Evaluate(this_.templateScope, this_.typeRegistry)
			_ = infoBaseType
			if !(FunctionTypeBuilder_AreCompatibleExtendsTypes(infoBaseType.ToMaybeObjectType(), classExtendsType)) {
				this_.isKnownAmbiguous = true
			}
			if infoBaseType.SetValidator(j2gort.ToIface_Ext_com_google_common_base_Predicate(this_.priv_FunctionTypeBuilder_createExtendedTypeValidator_2960())) {
				this_.baseType = infoBaseType.ToObjectType()
			} else {
				this_.isKnownAmbiguous = true
			}
		} else {
			this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF(), []string{this_.FormatFnName()})
			this_.isKnownAmbiguous = true
		}
	} else if (classExtendsType != nil && !classExtendsType.IsJavaNil_()) && (this_.isConstructor || this_.isInterface) {
		this_.baseType = classExtendsType
	}
	if (info != nil) && (info.GetImplementedInterfaceCount() > int32(0)) {
		if this_.isConstructor {
			var implCount int32 = info.GetImplementedInterfaceCount()
			_ = implCount
			this_.implementedInterfaces = j2gort.NewExt_java_util_ArrayList_1_int32(implCount)
			{
				var i int32 = int32(0)
				_ = i
				for ; i < implCount; i = i + 1 {
					var t pkg_rhino.JSTypeExpression = info.GetImplementedInterfaceAt(i)
					_ = t
					var maybeInterType pkg_rhino.JSType = t.Evaluate(this_.templateScope, this_.typeRegistry)
					_ = maybeInterType
					if maybeInterType.SetValidator(j2gort.ToIface_Ext_com_google_common_base_Predicate(this_.priv_FunctionTypeBuilder_createImplementedTypeValidator_2961())) {
						this_.implementedInterfaces.Add_1(j2gort.AsIface[pkg_rhino.ObjectType](maybeInterType))
					}
				}
			}
		} else if this_.isInterface {
			this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_TypeCheck_CONFLICTING_IMPLEMENTED_TYPE(), []string{this_.FormatFnName()})
		} else {
			this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED(), []string{"@implements", this_.FormatFnName()})
		}
	}
	if this_.isInterface {
		this_.extendedInterfaces = j2gort.NewExt_java_util_ArrayList()
		if info != nil {
			{
				var i_2 int32 = int32(0)
				_ = i_2
				var extCount int32 = info.GetExtendedInterfacesCount()
				_ = extCount
				for ; i_2 < extCount; i_2 = i_2 + 1 {
					var t_2 pkg_rhino.JSTypeExpression = info.GetExtendedInterfaceAt(i_2)
					_ = t_2
					var maybeInterfaceType pkg_rhino.JSType = t_2.Evaluate(this_.templateScope, this_.typeRegistry)
					_ = maybeInterfaceType
					if maybeInterfaceType != nil && !maybeInterfaceType.IsJavaNil_() {
						var isValid bool = maybeInterfaceType.SetValidator(j2gort.ToIface_Ext_com_google_common_base_Predicate(this_.priv_FunctionTypeBuilder_createExtendedTypeValidator_2960()))
						_ = isValid
						if isValid && !(j2gort.IfaceEquals[pkg_rhino.ObjectType](maybeInterfaceType.ToMaybeObjectType(), nil)) {
							this_.extendedInterfaces.Add_1(maybeInterfaceType.ToMaybeObjectType())
						}
					}
					if (classExtendsType != nil && !classExtendsType.IsJavaNil_()) && maybeInterfaceType.IsSubtypeOf_1(j2gort.AsIface[pkg_rhino.JSType](classExtendsType)) {
						classExtendsType = nil
					}
				}
			}
		}
		if (classExtendsType != nil && !classExtendsType.IsJavaNil_()) && classExtendsType.SetValidator(j2gort.ToIface_Ext_com_google_common_base_Predicate(this_.priv_FunctionTypeBuilder_createExtendedTypeValidator_2960())) {
			this_.extendedInterfaces.Add_1(classExtendsType)
		}
	}
	return this_
}

// Java static method: FunctionTypeBuilder.areCompatibleExtendsTypes(com.google.javascript.rhino.jstype.ObjectType annotated, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType extendsClause)
func FunctionTypeBuilder_AreCompatibleExtendsTypes(annotated pkg_rhino.ObjectType, extendsClause pkg_rhino.ObjectType) bool {
	clinit_FunctionTypeBuilder()
	_ = annotated
	_ = extendsClause
	if (extendsClause == nil || extendsClause.IsJavaNil_()) || (annotated == extendsClause || ((annotated == nil || annotated.IsJavaNil_()) && (extendsClause == nil || extendsClause.IsJavaNil_()))) {
		return true
	}
	return (annotated.IsTemplatizedType() && j2gort.IfaceEquals[pkg_rhino.ObjectType](annotated.ToMaybeTemplatizedType().GetReferencedType(), extendsClause))
}

// Java method: FunctionTypeBuilder.inferThisType(com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.jstype.JSType type)
func (this_ *FunctionTypeBuilder) InferThisType_2(info *pkg_rhino.JSDocInfo, type_ pkg_rhino.JSType) *FunctionTypeBuilder {
	_ = info
	_ = type_
	this_.InferThisType_1(info)
	if this_.thisType == nil || this_.thisType.IsJavaNil_() {
		var objType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(type_)
		_ = objType
		if (objType != nil && !objType.IsJavaNil_()) && ((info == nil) || !(info.HasType())) {
			this_.thisType = j2gort.AsIface[pkg_rhino.JSType](objType)
		}
	}
	return this_
}

// Java method: FunctionTypeBuilder.inferThisType(com.google.javascript.rhino.JSDocInfo info)
func (this_ *FunctionTypeBuilder) InferThisType_1(info *pkg_rhino.JSDocInfo) *FunctionTypeBuilder {
	_ = info
	if (info != nil) && info.HasThisType() {
		var maybeThisType pkg_rhino.JSType = info.GetThisType().Evaluate(this_.templateScope, this_.typeRegistry).RestrictByNotNullOrUndefined()
		_ = maybeThisType
		if maybeThisType != nil && !maybeThisType.IsJavaNil_() {
			this_.thisType = maybeThisType
		}
	}
	return this_
}

// Java method: FunctionTypeBuilder.inferParameterTypes(com.google.javascript.rhino.JSDocInfo info)
func (this_ *FunctionTypeBuilder) InferParameterTypes_1(info *pkg_rhino.JSDocInfo) *FunctionTypeBuilder {
	_ = info
	var lp pkg_rhino.Node = pkg_rhino.IR_ParamList([]pkg_rhino.Node{})
	_ = lp
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(info.GetParameterNames()))
		for cur_7.HasNext() {
			name := cur_7.NextStr()
			_ = name
			lp.AddChildToBack(pkg_rhino.IR_Name(name))
		}
	}
	return this_.InferParameterTypes_2(lp, info)
}

// Java method: FunctionTypeBuilder.inferParameterTypes(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node paramsParent, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info)
func (this_ *FunctionTypeBuilder) InferParameterTypes_2(paramsParent pkg_rhino.Node, info *pkg_rhino.JSDocInfo) *FunctionTypeBuilder {
	_ = paramsParent
	_ = info
	if paramsParent == 0 {
		if info == nil {
			return this_
		} else {
			return this_.InferParameterTypes_1(info)
		}
	}
	var oldParameters []pkg_rhino.FunctionType_Parameter = this_.parameters
	_ = oldParameters
	var oldParamLen int32
	_ = oldParamLen
	if oldParameters != nil {
		oldParamLen = int32(len(oldParameters))
	} else {
		oldParamLen = int32(0)
	}
	var oldParamIdx int32 = int32(0)
	_ = oldParamIdx
	var oldParameterType pkg_rhino.FunctionType_Parameter
	_ = oldParameterType
	if oldParamIdx < oldParamLen {
		oldParameterType = oldParameters[func() int32 { old_8 := oldParamIdx; nw_9 := old_8 + 1; oldParamIdx = nw_9; return old_8 }()]
	} else {
		oldParameterType = (pkg_rhino.FunctionType_Parameter{})
	}
	var builder *pkg_rhino.FunctionParamBuilder = pkg_rhino.NewFunctionParamBuilder_1(this_.typeRegistry)
	_ = builder
	var warnedAboutArgList bool = false
	_ = warnedAboutArgList
	var allJsDocParams *j2gort.OrderedSet[string]
	_ = allJsDocParams
	if (info == nil) || (info.GetParameterCount() == int32(0)) {
		allJsDocParams = nil
	} else {
		allJsDocParams = j2gort.OrderedSetPtrFromStrCol(j2gort.AsIface[j2gort.Ext_java_util_Collection](info.GetParameterNames()))
	}
	var isVarArgs bool = false
	_ = isVarArgs
	var paramIndex int32 = int32(0)
	_ = paramIndex
	{
		var param pkg_rhino.Node = paramsParent.GetFirstChild()
		_ = param
		for ; param != 0; param = param.GetNext() {
			var isOptionalParam bool = false
			_ = isOptionalParam
			var paramLhs pkg_rhino.Node = pkg_rhino.Node(0)
			_ = paramLhs
			if param.IsRest() {
				isVarArgs = true
				paramLhs = param.GetOnlyChild()
			} else if param.IsDefaultValue() {
				paramLhs = j2gort.CheckNotNullNode(param.GetFirstChild())
				isOptionalParam = true
			} else {
				isVarArgs = this_.priv_FunctionTypeBuilder_isVarArgsParameterByConvention_2983(param)
				isOptionalParam = this_.priv_FunctionTypeBuilder_isOptionalParameterByConvention_2982(param)
				paramLhs = param
			}
			var paramName string = j2gort.NullString
			_ = paramName
			if paramLhs.IsName() {
				paramName = paramLhs.GetString()
			} else {
				if !(paramLhs.IsDestructuringPattern()) {
					j2gort.CheckState(false)
				}
				if info != nil {
					paramName = info.GetParameterNameAt(paramIndex)
				}
			}
			if allJsDocParams != nil {
				allJsDocParams.Remove(paramName)
			}
			var parameterType pkg_rhino.JSType = nil
			_ = parameterType
			if (info != nil) && info.HasParameterType(paramName) {
				var parameterTypeExpression pkg_rhino.JSTypeExpression = info.GetParameterType(paramName)
				_ = parameterTypeExpression
				parameterType = parameterTypeExpression.Evaluate(this_.templateScope, this_.typeRegistry)
				isOptionalParam = (isOptionalParam || parameterTypeExpression.IsOptionalArg())
				isVarArgs = (isVarArgs || parameterTypeExpression.IsVarArgs())
			} else if (paramLhs.GetJSDocInfo() != nil) && paramLhs.GetJSDocInfo().HasType() {
				var parameterTypeExpression_2 pkg_rhino.JSTypeExpression = paramLhs.GetJSDocInfo().GetType()
				_ = parameterTypeExpression_2
				parameterType = parameterTypeExpression_2.Evaluate(this_.templateScope, this_.typeRegistry)
				isOptionalParam = parameterTypeExpression_2.IsOptionalArg()
				isVarArgs = parameterTypeExpression_2.IsVarArgs()
			} else if (oldParameterType != (pkg_rhino.FunctionType_Parameter{})) && !(j2gort.IfaceEquals[pkg_rhino.JSType](oldParameterType.GetJSType(), nil)) {
				parameterType = oldParameterType.GetJSType()
				isOptionalParam = oldParameterType.IsOptional()
				isVarArgs = oldParameterType.IsVariadic()
			} else {
				parameterType = this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
			}
			warnedAboutArgList = func() bool {
				ba_10 := warnedAboutArgList
				bb_11 := this_.priv_FunctionTypeBuilder_addParameter_2986(builder, parameterType, warnedAboutArgList, isOptionalParam, isVarArgs)
				return ba_10 || bb_11
			}()
			if oldParamIdx < oldParamLen {
				oldParameterType = oldParameters[func() int32 { old_12 := oldParamIdx; nw_13 := old_12 + 1; oldParamIdx = nw_13; return old_12 }()]
			} else {
				oldParameterType = (pkg_rhino.FunctionType_Parameter{})
			}
			paramIndex = paramIndex + 1
		}
	}
	if !(isVarArgs) {
		for (oldParameterType != (pkg_rhino.FunctionType_Parameter{})) && !(isVarArgs) {
			builder.NewParameterFrom(oldParameterType)
			if oldParamIdx < oldParamLen {
				oldParameterType = oldParameters[func() int32 { old_14 := oldParamIdx; nw_15 := old_14 + 1; oldParamIdx = nw_15; return old_14 }()]
			} else {
				oldParameterType = (pkg_rhino.FunctionType_Parameter{})
			}
		}
	}
	if allJsDocParams != nil {
		{
			cur_16 := allJsDocParams.Iter()
			for cur_16.HasNext() {
				inexistentName := cur_16.
					Next()
				_ = inexistentName
				this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_FunctionTypeBuilder_INEXISTENT_PARAM(), []string{inexistentName, this_.FormatFnName()})
			}
		}
	}
	this_.parameters = builder.BuildArray()
	return this_
}

// Java method: FunctionTypeBuilder.registerTemplates(java.lang.Iterable<com.google.javascript.rhino.jstype.TemplateType> templates, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node scopeRoot)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_registerTemplates_2976(templates j2gort.Ext_java_lang_Iterable, scopeRoot pkg_rhino.Node) {
	_ = templates
	_ = scopeRoot
	if !(j2gort.Ext_com_google_common_collect_Iterables_IsEmpty(templates)) {
		this_.templateScope = this_.typeRegistry.CreateScopeWithTemplates_2_StaticTypedScope_Ext_java_lang_Iterable(this_.templateScope, templates)
		if scopeRoot != 0 {
			this_.typeRegistry.RegisterTemplateTypeNamesInScope_2_Ext_java_lang_Iterable_Node(templates, scopeRoot)
		}
	}
}

// Java method: FunctionTypeBuilder.registerTemplatesWithOwner(com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.TemplateType> templates, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable TemplateTypeMap ownerMap, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node scopeRoot)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_registerTemplatesWithOwner_2977(templates j2gort.Ext_com_google_common_collect_ImmutableList, ownerMap *pkg_rhino.TemplateTypeMap, scopeRoot pkg_rhino.Node) {
	_ = templates
	_ = ownerMap
	_ = scopeRoot
	var hasOwn bool = !(templates.IsEmpty())
	_ = hasOwn
	var hasOwner bool = ((ownerMap != nil) && !(ownerMap.IsEmpty()))
	_ = hasOwner
	if hasOwn || hasOwner {
		this_.templateScope = this_.typeRegistry.CreateScopeWithTemplates_3(this_.templateScope, templates, ownerMap)
		if scopeRoot != 0 {
			{
				var i int32 = int32(0)
				_ = i
				var n int32 = templates.Size()
				_ = n
				for ; i < n; i = i + 1 {
					this_.typeRegistry.RegisterSingleTemplateTypeNameInScope(j2gort.CastPtr[pkg_rhino.TemplateType](templates.Get(i)), scopeRoot)
				}
			}
			if hasOwner {
				this_.typeRegistry.RegisterTemplateTypeNamesInScope_2_TemplateTypeMap_Node(ownerMap, scopeRoot)
			}
		}
	}
}

// Java method: FunctionTypeBuilder.inferConstructorParameters(com.google.javascript.rhino.Node argsParent, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info)
func (this_ *FunctionTypeBuilder) InferConstructorParameters(argsParent pkg_rhino.Node, info *pkg_rhino.JSDocInfo) *FunctionTypeBuilder {
	_ = argsParent
	_ = info
	if info != nil {
		this_.priv_FunctionTypeBuilder_setConstructorTemplateTypeNames_2981(j2gort.AsIface[j2gort.Ext_java_util_List](this_.priv_FunctionTypeBuilder_buildTemplateTypesFromJSDocInfo_2984(info, true)), argsParent.GetParent())
	}
	this_.InferParameterTypes_2(argsParent, info)
	return this_
}

// Java method: FunctionTypeBuilder.inferImplicitConstructorParameters(com.google.javascript.rhino.jstype.FunctionType.Parameter[] parameters)
func (this_ *FunctionTypeBuilder) InferImplicitConstructorParameters(parameters []pkg_rhino.FunctionType_Parameter) *FunctionTypeBuilder {
	_ = parameters
	this_.parameters = parameters
	return this_
}

// Java method: FunctionTypeBuilder.inferClosurePrimitive(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info)
func (this_ *FunctionTypeBuilder) InferClosurePrimitive(info *pkg_rhino.JSDocInfo) *FunctionTypeBuilder {
	_ = info
	if (info != nil) && info.HasClosurePrimitiveId() {
		this_.closurePrimitiveId = pkg_rhino.ClosurePrimitive_FromStringId(info.GetClosurePrimitiveId())
	}
	return this_
}

// Java method: FunctionTypeBuilder.setConstructorTemplateTypeNames(java.util.List<com.google.javascript.rhino.jstype.TemplateType> templates, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node ctor)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_setConstructorTemplateTypeNames_2981(templates j2gort.Ext_java_util_List, ctor pkg_rhino.Node) {
	_ = templates
	_ = ctor
	if !(templates.IsEmpty()) {
		this_.constructorTemplateTypeNames = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](templates))
		if this_.templateTypeNames.IsEmpty() {
			this_.templateTypeNames = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](templates))
		} else {
			this_.templateTypeNames = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}).AddAllSelf_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.templateTypeNames)).AddAllSelf_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.constructorTemplateTypeNames)).BuildImmList_()
		}
		this_.priv_FunctionTypeBuilder_registerTemplates_2976(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](templates), ctor)
	}
}

// Java method: FunctionTypeBuilder.isOptionalParameterByConvention(com.google.javascript.rhino.Node param)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_isOptionalParameterByConvention_2982(param pkg_rhino.Node) bool {
	_ = param
	if param.IsDestructuringPattern() {
		return false
	}
	return this_.codingConvention.IsOptionalParameter(param)
}

// Java method: FunctionTypeBuilder.isVarArgsParameterByConvention(com.google.javascript.rhino.Node param)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_isVarArgsParameterByConvention_2983(param pkg_rhino.Node) bool {
	_ = param
	if param.IsDestructuringPattern() {
		return false
	}
	return this_.codingConvention.IsVarArgsParameter(param)
}

// Java method: FunctionTypeBuilder.buildTemplateTypesFromJSDocInfo(com.google.javascript.rhino.JSDocInfo info, boolean allowTypeTransformations)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_buildTemplateTypesFromJSDocInfo_2984(info *pkg_rhino.JSDocInfo, allowTypeTransformations bool) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = info
	_ = allowTypeTransformations
	var infoTypeKeys j2gort.Ext_com_google_common_collect_ImmutableMap = info.GetTemplateTypes()
	_ = infoTypeKeys
	var infoTypeTransformations j2gort.Ext_com_google_common_collect_ImmutableMap = info.GetTypeTransformations()
	_ = infoTypeTransformations
	if infoTypeKeys.IsEmpty() && infoTypeTransformations.IsEmpty() {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var unboundedTemplates j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = unboundedTemplates
	{
		cur_17 := j2gort.NewMapKeyStrCursor(infoTypeKeys)
		for cur_17.HasNext() {
			templateKey := cur_17.NextStr()
			_ = templateKey
			unboundedTemplates.Add_1(j2gort.PtrToAny(this_.typeRegistry.CreateTemplateType_1(templateKey)))
		}
	}
	this_.templateScope = this_.typeRegistry.CreateScopeWithTemplates_2_StaticTypedScope_Ext_java_lang_Iterable(this_.templateScope, unboundedTemplates)
	var templates []*pkg_rhino.TemplateType = []*pkg_rhino.TemplateType{}
	_ = templates
	var templatesToBounds j2gort.CustomOrderedMap[*pkg_rhino.TemplateType, pkg_rhino.JSType] = j2gort.NewCustomOrderedMap[*pkg_rhino.TemplateType, pkg_rhino.JSType]()
	_ = templatesToBounds
	{
		cur_18 := j2gort.NewMapEntryCursor(infoTypeKeys)
		for cur_18.NextEntry() {
			var expr pkg_rhino.JSTypeExpression = j2gort.Cast[pkg_rhino.JSTypeExpression](cur_18.GetValue())
			_ = expr
			var typeBound pkg_rhino.JSType = this_.typeRegistry.EvaluateTypeExpression(j2gort.Cast[pkg_rhino.JSTypeExpression](cur_18.GetValue()), this_.templateScope)
			_ = typeBound
			if expr.IsExplicitUnknownTemplateBound() {
				this_.priv_FunctionTypeBuilder_reportError_2994(StaticGet_FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND(), []string{j2gort.StrVal(typeBound), cur_18.KeyStr()})
			}
			var template *pkg_rhino.TemplateType = this_.typeRegistry.GetType_2(j2gort.AsIface[pkg_rhino.StaticScope](this_.templateScope), cur_18.KeyStr()).ToMaybeTemplateType()
			_ = template
			if template != nil {
				templatesToBounds.Put(template, typeBound)
			} else {
				templatesToBounds.Put(this_.typeRegistry.CreateTemplateType_2(cur_18.KeyStr(), typeBound), typeBound)
			}
		}
	}
	{
		cur_19 := templatesToBounds.IterEntries()
		for cur_19.NextEntry() {
			var template_2 *pkg_rhino.TemplateType = cur_19.Key()
			_ = template_2
			var bound pkg_rhino.JSType = cur_19.Value()
			_ = bound
			template_2.SetBound(bound)
			templates = append(templates, j2gort.CheckNotNullElem(template_2))
		}
	}
	{
		cur_20 := j2gort.NewMapEntryCursor(infoTypeTransformations)
		for cur_20.NextEntry() {
			if allowTypeTransformations {
				templates = append(templates, j2gort.CheckNotNullElem(this_.typeRegistry.CreateTemplateTypeWithTransformation(cur_20.KeyStr(), j2gort.CastNode[pkg_rhino.Node](cur_20.GetValue()))))
			} else {
				this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS(), []string{cur_20.KeyStr()})
			}
		}
	}
	var builtTemplates []*pkg_rhino.TemplateType = templates
	_ = builtTemplates
	{

		for _, template_3 := range builtTemplates {

			_ = template_3
			if template_3.ContainsCycle() {
				this_.priv_FunctionTypeBuilder_reportError_2994(StaticGet_RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR(), []string{("Cycle detected in inheritance chain of type " + j2gort.StrValString(template_3.GetReferenceName()))})
			}
		}
	}
	return j2gort.ImmListFromTypedSlice(builtTemplates)
}

// Java method: FunctionTypeBuilder.inferTemplateTypeName(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType ownerType)
func (this_ *FunctionTypeBuilder) InferTemplateTypeName(info *pkg_rhino.JSDocInfo, ownerType pkg_rhino.JSType) *FunctionTypeBuilder {
	_ = info
	_ = ownerType
	if (info != nil) && !(this_.priv_FunctionTypeBuilder_maybeUseNativeClassTemplateNames_2970(info)) {
		var templates j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_FunctionTypeBuilder_buildTemplateTypesFromJSDocInfo_2984(info, !(this_.isConstructor || this_.isInterface))
		_ = templates
		if !(templates.IsEmpty()) {
			this_.templateTypeNames = templates
		}
	}
	var ownerTypeMap *pkg_rhino.TemplateTypeMap
	_ = ownerTypeMap
	if ownerType != nil && !ownerType.IsJavaNil_() {
		ownerTypeMap = ownerType.GetTemplateTypeMap()
	} else {
		ownerTypeMap = nil
	}
	if !(this_.templateTypeNames.IsEmpty()) || ((ownerTypeMap != nil) && !(ownerTypeMap.IsEmpty())) {
		this_.priv_FunctionTypeBuilder_registerTemplatesWithOwner_2977(this_.templateTypeNames, ownerTypeMap, this_.contents.GetSourceNode())
	}
	return this_
}

// Java method: FunctionTypeBuilder.addParameter(com.google.javascript.rhino.jstype.FunctionParamBuilder builder, com.google.javascript.rhino.jstype.JSType paramType, boolean warnedAboutArgList, boolean isOptional, boolean isVarArgs)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_addParameter_2986(builder *pkg_rhino.FunctionParamBuilder, paramType pkg_rhino.JSType, warnedAboutArgList bool, isOptional bool, isVarArgs bool) bool {
	_ = builder
	_ = paramType
	_ = warnedAboutArgList
	_ = isOptional
	_ = isVarArgs
	var emittedWarning bool = false
	_ = emittedWarning
	if isOptional {
		if !(builder.AddOptionalParams([]pkg_rhino.JSType{paramType})) && !(warnedAboutArgList) {
			this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST(), []string{})
			emittedWarning = true
		}
	} else if isVarArgs {
		if !(builder.AddVarArgs(paramType)) && !(warnedAboutArgList) {
			this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST(), []string{})
			emittedWarning = true
		}
	} else {
		if !(builder.AddRequiredParams([]pkg_rhino.JSType{paramType})) && !(warnedAboutArgList) {
			if builder.HasVarArgs() {
				this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST(), []string{})
			} else {
				this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_FunctionTypeBuilder_OPTIONAL_ARG_AT_END(), []string{})
			}
			emittedWarning = true
		}
	}
	return emittedWarning
}

// Java method: FunctionTypeBuilder.provideDefaultReturnType()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_provideDefaultReturnType_2987() {
	if (this_.contents.GetSourceNode() != 0) && this_.contents.GetSourceNode().IsAsyncGeneratorFunction() {
		var generatorType pkg_rhino.ObjectType = this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ASYNC_GENERATOR_TYPE())
		_ = generatorType
		this_.returnType = pkg_rhino.ToIface_JSType(this_.typeRegistry.CreateTemplatizedType_2_ObjectType_JSType(generatorType, []pkg_rhino.JSType{this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())}))
		return
	} else if (this_.contents.GetSourceNode() != 0) && this_.contents.GetSourceNode().IsGeneratorFunction() {
		var generatorType_2 pkg_rhino.ObjectType = this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_GENERATOR_TYPE())
		_ = generatorType_2
		this_.returnType = pkg_rhino.ToIface_JSType(this_.typeRegistry.CreateTemplatizedType_2_ObjectType_JSType(generatorType_2, []pkg_rhino.JSType{this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())}))
		return
	}
	var inferredReturnType pkg_rhino.JSType = this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	_ = inferredReturnType
	if (!(this_.contents.MayHaveNonEmptyReturns()) && !(this_.contents.MayHaveSingleThrow())) && !(this_.contents.MayBeFromExterns()) {
		inferredReturnType = this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
		this_.returnTypeInferred = true
	}
	if (this_.contents.GetSourceNode() != 0) && this_.contents.GetSourceNode().IsAsyncFunction() {
		var promiseType pkg_rhino.ObjectType = this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_PROMISE_TYPE())
		_ = promiseType
		this_.returnType = pkg_rhino.ToIface_JSType(this_.typeRegistry.CreateTemplatizedType_2_ObjectType_JSType(promiseType, []pkg_rhino.JSType{inferredReturnType}))
	} else {
		this_.returnType = inferredReturnType
	}
}

// Java method: FunctionTypeBuilder.buildAndRegister()
func (this_ *FunctionTypeBuilder) BuildAndRegister() *pkg_rhino.FunctionType {
	if this_.returnType == nil || this_.returnType.IsJavaNil_() {
		this_.priv_FunctionTypeBuilder_provideDefaultReturnType_2987()
		j2gort.CheckNotNullIface(this_.returnType)
	}
	if this_.parameters == nil {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "All Function types must have params and a return type")))
	}
	var fnType *pkg_rhino.FunctionType = nil
	_ = fnType
	if this_.isConstructor {
		fnType = this_.priv_FunctionTypeBuilder_getOrCreateConstructor_2991()
	} else if this_.isInterface {
		fnType = this_.priv_FunctionTypeBuilder_getOrCreateInterface_2992()
	} else {
		var fb_ pkg_rhino.FunctionType_Builder
		this_.priv_FunctionTypeBuilder_createDefaultBuilder_2990(&fb_)
		(&fb_).WithParameterArray(this_.parameters)
		(&fb_).WithReturnType_2(this_.returnType, this_.returnTypeInferred)
		(&fb_).WithTypeOfThis(this_.thisType)
		(&fb_).WithIsAbstract(this_.isAbstract)
		(&fb_).WithClosurePrimitiveId(this_.closurePrimitiveId)
		fnType = j2gort.CastPtr[pkg_rhino.FunctionType]((&fb_).Build())
		this_.priv_FunctionTypeBuilder_maybeSetBaseType_2989(fnType)
	}
	if (this_.implementedInterfaces != nil) && fnType.IsConstructor() {
		fnType.SetImplementedInterfaces(this_.implementedInterfaces)
	}
	if this_.extendedInterfaces != nil {
		fnType.SetExtendedInterfaces(this_.extendedInterfaces)
	}
	if this_.isRecord {
		fnType.SetImplicitMatch(true)
	}
	return fnType
}

// Java method: FunctionTypeBuilder.maybeSetBaseType(com.google.javascript.rhino.jstype.FunctionType fnType)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_maybeSetBaseType_2989(fnType *pkg_rhino.FunctionType) {
	_ = fnType
	if !(fnType.HasInstanceType()) || (this_.baseType == nil || this_.baseType.IsJavaNil_()) {
		return
	}
	fnType.SetPrototypeBasedOn_1(this_.baseType)
	fnType.GetInstanceType().MergeSupertypeTemplateTypes(this_.baseType)
}

// Java method: FunctionTypeBuilder.createDefaultBuilder()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_createDefaultBuilder_2990(fb_ *pkg_rhino.FunctionType_Builder) {
	pkg_rhino.InitFunctionType_Builder(fb_, this_.typeRegistry)
	fb_.WithName(this_.fnName)
	fb_.WithSourceNode(this_.contents.GetSourceNode())
	fb_.WithTemplateKeys_1_Ext_com_google_common_collect_ImmutableList(this_.templateTypeNames)
	fb_.SetIsKnownAmbiguous(this_.isKnownAmbiguous)
	fb_.SetGoogModuleId(TypedScopeCreator_ContainingGoogModuleIdOf(this_.declarationScope))
}

// Java method: FunctionTypeBuilder.getOrCreateConstructor()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_getOrCreateConstructor_2991() *pkg_rhino.FunctionType {
	var fb_ pkg_rhino.FunctionType_Builder
	this_.priv_FunctionTypeBuilder_createDefaultBuilder_2990(&fb_)
	(&fb_).ForConstructor()
	(&fb_).WithParameterArray(this_.parameters)
	(&fb_).WithReturnType_1(this_.returnType)
	(&fb_).WithConstructorTemplateKeys(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.constructorTemplateTypeNames))
	(&fb_).WithIsAbstract(this_.isAbstract)
	var fnType *pkg_rhino.FunctionType = j2gort.CastPtr[pkg_rhino.FunctionType]((&fb_).Build())
	_ = fnType
	if this_.makesStructs {
		fnType.SetStruct()
	} else if this_.makesDicts {
		fnType.SetDict()
	} else if this_.makesUnrestricted {
		fnType.SetExplicitUnrestricted()
	}
	var existingType pkg_rhino.JSType = this_.typeRegistry.GetType_2(pkg_rhino.ToIface_StaticScope(this_.declarationScope), this_.fnName)
	_ = existingType
	if existingType != nil && !existingType.IsJavaNil_() {
		var isInstanceObject bool = existingType.IsInstanceType()
		_ = isInstanceObject
		if isInstanceObject || (j2gort.CheckNotNullString(this_.fnName) == "Function") {
			var existingFn *pkg_rhino.FunctionType
			_ = existingFn
			if isInstanceObject {
				existingFn = existingType.ToObjectType().GetConstructor()
			} else {
				existingFn = this_.typeRegistry.GetNativeFunctionType(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_FUNCTION_TYPE())
			}
			if existingFn.GetSource() == 0 {
				existingFn.SetSource(this_.contents.GetSourceNode())
			}
			if !(existingFn.HasEqualCallType(fnType)) {
				this_.priv_FunctionTypeBuilder_reportWarning_2993(StaticGet_FunctionTypeBuilder_TYPE_REDEFINITION(), []string{this_.FormatFnName(), fnType.String(), existingFn.String()})
			}
			return existingFn
		} else {
		}
	}
	this_.priv_FunctionTypeBuilder_maybeSetBaseType_2989(fnType)
	if !(j2gort.CheckNotNullString(this_.syntacticFnName) == "") && !(javalang.StringStartsWith(j2gort.CheckNotNullString(this_.syntacticFnName), "this.")) {
		this_.typeRegistry.DeclareTypeForExactScope(pkg_rhino.ToIface_StaticScope(this_.declarationScope), this_.syntacticFnName, j2gort.AsIface[pkg_rhino.JSType](fnType.GetInstanceType()))
	}
	return fnType
}

// Java method: FunctionTypeBuilder.getOrCreateInterface()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_getOrCreateInterface_2992() *pkg_rhino.FunctionType {
	var fnType *pkg_rhino.FunctionType = nil
	_ = fnType
	var type_ pkg_rhino.JSType = this_.typeRegistry.GetType_2(pkg_rhino.ToIface_StaticScope(this_.declarationScope), this_.syntacticFnName)
	_ = type_
	if (type_ != nil && !type_.IsJavaNil_()) && type_.IsInstanceType() {
		var ctor *pkg_rhino.FunctionType = type_.ToMaybeObjectType().GetConstructor()
		_ = ctor
		if ctor.IsInterface() {
			fnType = ctor
			fnType.SetSource(this_.contents.GetSourceNode())
		}
	}
	if fnType == nil {
		var fb_ pkg_rhino.FunctionType_Builder
		this_.priv_FunctionTypeBuilder_createDefaultBuilder_2990(&fb_)
		(&fb_).ForInterface()
		(&fb_).WithParameters_0()
		fnType = j2gort.CastPtr[pkg_rhino.FunctionType]((&fb_).Build())
		if this_.makesStructs {
			fnType.SetStruct()
		}
		if !(j2gort.CheckNotNullString(this_.fnName) == "") {
			this_.typeRegistry.DeclareTypeForExactScope(pkg_rhino.ToIface_StaticScope(this_.declarationScope), this_.syntacticFnName, j2gort.AsIface[pkg_rhino.JSType](fnType.GetInstanceType()))
		}
		this_.priv_FunctionTypeBuilder_maybeSetBaseType_2989(fnType)
	}
	return fnType
}

// Java method: FunctionTypeBuilder.reportWarning(com.google.javascript.jscomp.DiagnosticType warning, java.lang.String[] args)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_reportWarning_2993(warning *DiagnosticType, args []string) {
	_ = warning
	_ = args
	this_.compiler.Report_1(JSError_Make_3(this_.errorRoot, warning, args))
}

// Java method: FunctionTypeBuilder.reportError(com.google.javascript.jscomp.DiagnosticType error, java.lang.String[] args)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_reportError_2994(error_ *DiagnosticType, args []string) {
	_ = error_
	_ = args
	this_.compiler.Report_1(JSError_Make_3(this_.errorRoot, error_, args))
}

// Java static method: FunctionTypeBuilder.isFunctionTypeDeclaration(com.google.javascript.rhino.JSDocInfo info)
func FunctionTypeBuilder_IsFunctionTypeDeclaration(info *pkg_rhino.JSDocInfo) bool {
	clinit_FunctionTypeBuilder()
	_ = info
	return ((((((info.GetParameterCount() > int32(0)) || info.HasReturnType()) || info.HasThisType()) || info.IsConstructor()) || info.IsInterface()) || info.IsAbstract())
}

// Java static method: FunctionTypeBuilder.hasMoreTagsToResolve(com.google.javascript.rhino.jstype.ObjectType objectType)
func FunctionTypeBuilder_HasMoreTagsToResolve(objectType pkg_rhino.ObjectType) bool {
	clinit_FunctionTypeBuilder()
	_ = objectType
	if !(objectType.IsUnknownType()) {
		j2gort.CheckArgument(false)
	}
	var ctor *pkg_rhino.FunctionType = objectType.GetConstructor()
	_ = ctor
	if ctor != nil {
		{
			cur_23 := j2gort.NewIterCursor(j2gort.CheckNotNull(ctor.GetExtendedInterfaces()))
			for cur_23.HasNext() {
				rawElem_24 := cur_23.Next()
				interfaceType := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_24)
				_ = interfaceType
				if !(interfaceType.IsResolved()) {
					return true
				}
			}
		}
	}
	if !(j2gort.IfaceEquals[pkg_rhino.ObjectType](objectType.GetImplicitPrototype(), nil)) {
		return !(objectType.GetImplicitPrototype().IsResolved())
	}
	return false
}

// Java polymorphic class interface: com.google.javascript.jscomp.FunctionTypeBuilder.ValidatorBase
type FunctionTypeBuilder_ValidatorBase interface {
	Implements_jscomp_FunctionTypeBuilder_ValidatorBase()
	IsJavaNil_() bool
	GetFunctionTypeBuilder_ValidatorBaseBase() *FunctionTypeBuilder_ValidatorBaseBase
	j2gort.Ext_com_google_common_base_Predicate
	ReportWarning(*DiagnosticType, []string)
	ReportError(*DiagnosticType, []string)
}

func ToIface_FunctionTypeBuilder_ValidatorBase[T interface {
	FunctionTypeBuilder_ValidatorBase
	comparable
}](v T) FunctionTypeBuilder_ValidatorBase {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.FunctionTypeBuilder.ValidatorBase
type FunctionTypeBuilder_ValidatorBaseBase struct {
	self_     FunctionTypeBuilder_ValidatorBase
	compiler  AbstractCompiler `gson:"compiler"`
	errorRoot pkg_rhino.Node   `gson:"errorRoot"`
	idHash_   atomic.Int32
}

func (this_ *FunctionTypeBuilder_ValidatorBaseBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(FunctionTypeBuilder_ValidatorBase)
}
func (this_ *FunctionTypeBuilder_ValidatorBaseBase) IsJavaNil_() bool { return this_ == nil }
func (*FunctionTypeBuilder_ValidatorBaseBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionTypeBuilder$ValidatorBase", "ValidatorBase"
}
func (this_ *FunctionTypeBuilder_ValidatorBaseBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionTypeBuilder.ValidatorBase", "com.google.javascript.jscomp.FunctionTypeBuilder$ValidatorBase", func(v_ any) bool {
		return j2gort.InstanceOf[FunctionTypeBuilder_ValidatorBase](v_)
	})
}
func (b *FunctionTypeBuilder_ValidatorBaseBase) GetFunctionTypeBuilder_ValidatorBaseBase() *FunctionTypeBuilder_ValidatorBaseBase {
	return b
}
func (*FunctionTypeBuilder_ValidatorBaseBase) Implements_jscomp_FunctionTypeBuilder_ValidatorBase() {}
func (*FunctionTypeBuilder_ValidatorBaseBase) Implements_j2gort_Ext_com_google_common_base_Predicate() {
}
func (*FunctionTypeBuilder_ValidatorBaseBase) Implements_j2gort_Ext_java_util_function_Predicate() {}

func (this_ *FunctionTypeBuilder_ValidatorBaseBase) initFields_FunctionTypeBuilder_ValidatorBase() {
	_ = this_
}

func InitFunctionTypeBuilder_ValidatorBase(this_ *FunctionTypeBuilder_ValidatorBaseBase, errorRoot pkg_rhino.Node, compiler AbstractCompiler) {
	_ = errorRoot
	_ = compiler
	this_.initFields_FunctionTypeBuilder_ValidatorBase()
	this_.errorRoot = errorRoot
	this_.compiler = compiler
}

// Java method: ValidatorBase.reportWarning(com.google.javascript.jscomp.DiagnosticType warning, java.lang.String[] args)
func (this_ *FunctionTypeBuilder_ValidatorBaseBase) ReportWarning(warning *DiagnosticType, args []string) {
	_ = warning
	_ = args
	this_.compiler.Report_1(JSError_Make_3(this_.errorRoot, warning, args))
}

// Java method: ValidatorBase.reportError(com.google.javascript.jscomp.DiagnosticType error, java.lang.String[] args)
func (this_ *FunctionTypeBuilder_ValidatorBaseBase) ReportError(error_ *DiagnosticType, args []string) {
	_ = error_
	_ = args
	this_.compiler.Report_1(JSError_Make_3(this_.errorRoot, error_, args))
}

func (this_ *FunctionTypeBuilder_ValidatorBaseBase) Equals(p0_ any) bool {
	return j2gort.RefEquals(this_.self_, p0_)
}

func (this_ *FunctionTypeBuilder_ValidatorBaseBase) Test(p0_ any) bool {
	return this_.self_.Apply(p0_)
}

// Java class: com.google.javascript.jscomp.FunctionTypeBuilder.ExtendedTypeValidator
type FunctionTypeBuilder_ExtendedTypeValidator struct {
	FunctionTypeBuilder_ValidatorBaseBase
	formattedFnName string `gson:"formattedFnName"`
}

func (this_ *FunctionTypeBuilder_ExtendedTypeValidator) BindSelf_(s any) {
	_ = s
	this_.FunctionTypeBuilder_ValidatorBaseBase.BindSelf_(s)
}
func (this_ *FunctionTypeBuilder_ExtendedTypeValidator) IsJavaNil_() bool { return this_ == nil }
func (*FunctionTypeBuilder_ExtendedTypeValidator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionTypeBuilder$ExtendedTypeValidator", "ExtendedTypeValidator"
}
func (this_ *FunctionTypeBuilder_ExtendedTypeValidator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.FunctionTypeBuilder_ValidatorBaseBase.IdentityHash_()
}
func (*FunctionTypeBuilder_ExtendedTypeValidator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionTypeBuilder.ExtendedTypeValidator", "com.google.javascript.jscomp.FunctionTypeBuilder$ExtendedTypeValidator", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionTypeBuilder_ExtendedTypeValidator](v_)
	})
}
func (*FunctionTypeBuilder_ExtendedTypeValidator) Implements_jscomp_FunctionTypeBuilder_ValidatorBase() {
}
func (*FunctionTypeBuilder_ExtendedTypeValidator) Implements_j2gort_Ext_com_google_common_base_Predicate() {
}
func (*FunctionTypeBuilder_ExtendedTypeValidator) Implements_j2gort_Ext_java_util_function_Predicate() {
}

func (this_ *FunctionTypeBuilder_ExtendedTypeValidator) initFields_FunctionTypeBuilder_ExtendedTypeValidator() {
	_ = this_
	this_.formattedFnName = j2gort.NullString
}

// Java constructor: ExtendedTypeValidator(com.google.javascript.rhino.Node errorRoot, com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String formattedFnName)
func NewFunctionTypeBuilder_ExtendedTypeValidator(errorRoot pkg_rhino.Node, compiler AbstractCompiler, formattedFnName string) *FunctionTypeBuilder_ExtendedTypeValidator {
	this_ := &FunctionTypeBuilder_ExtendedTypeValidator{}
	this_.BindSelf_(this_)
	InitFunctionTypeBuilder_ExtendedTypeValidator(this_, errorRoot, compiler, formattedFnName)
	return this_
}

func InitFunctionTypeBuilder_ExtendedTypeValidator(this_ *FunctionTypeBuilder_ExtendedTypeValidator, errorRoot pkg_rhino.Node, compiler AbstractCompiler, formattedFnName string) {
	_ = errorRoot
	_ = compiler
	_ = formattedFnName
	InitFunctionTypeBuilder_ValidatorBase(&this_.FunctionTypeBuilder_ValidatorBaseBase, errorRoot, compiler)
	this_.initFields_FunctionTypeBuilder_ExtendedTypeValidator()
	this_.formattedFnName = formattedFnName
}

// Java method: ExtendedTypeValidator.apply(com.google.javascript.rhino.jstype.JSType type)
func (this_ *FunctionTypeBuilder_ExtendedTypeValidator) Apply(arg_25 any) bool {
	type_ := j2gort.AsIface[pkg_rhino.JSType](arg_25)
	_ = type_
	var objectType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(type_)
	_ = objectType
	if objectType == nil || objectType.IsJavaNil_() {
		this_.ReportWarning(StaticGet_FunctionTypeBuilder_EXTENDS_NON_OBJECT(), []string{this_.formattedFnName, type_.String()})
		return false
	}
	if objectType.IsEmptyType() {
		this_.ReportWarning(StaticGet_FunctionTypeBuilder_RESOLVED_TAG_EMPTY(), []string{"@extends", this_.formattedFnName})
		return false
	}
	if objectType.IsUnknownType() {
		if FunctionTypeBuilder_HasMoreTagsToResolve(objectType) || type_.IsTemplateType() {
			return true
		} else {
			this_.ReportWarning(StaticGet_FunctionTypeBuilder_RESOLVED_TAG_EMPTY(), []string{"@extends", this_.formattedFnName})
			return false
		}
	}
	return true
}

// Java class: com.google.javascript.jscomp.FunctionTypeBuilder.ImplementedTypeValidator
type FunctionTypeBuilder_ImplementedTypeValidator struct {
	FunctionTypeBuilder_ValidatorBaseBase
	formattedFnName string `gson:"formattedFnName"`
}

func (this_ *FunctionTypeBuilder_ImplementedTypeValidator) BindSelf_(s any) {
	_ = s
	this_.FunctionTypeBuilder_ValidatorBaseBase.BindSelf_(s)
}
func (this_ *FunctionTypeBuilder_ImplementedTypeValidator) IsJavaNil_() bool { return this_ == nil }
func (*FunctionTypeBuilder_ImplementedTypeValidator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionTypeBuilder$ImplementedTypeValidator", "ImplementedTypeValidator"
}
func (this_ *FunctionTypeBuilder_ImplementedTypeValidator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.FunctionTypeBuilder_ValidatorBaseBase.IdentityHash_()
}
func (*FunctionTypeBuilder_ImplementedTypeValidator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionTypeBuilder.ImplementedTypeValidator", "com.google.javascript.jscomp.FunctionTypeBuilder$ImplementedTypeValidator", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionTypeBuilder_ImplementedTypeValidator](v_)
	})
}
func (*FunctionTypeBuilder_ImplementedTypeValidator) Implements_jscomp_FunctionTypeBuilder_ValidatorBase() {
}
func (*FunctionTypeBuilder_ImplementedTypeValidator) Implements_j2gort_Ext_com_google_common_base_Predicate() {
}
func (*FunctionTypeBuilder_ImplementedTypeValidator) Implements_j2gort_Ext_java_util_function_Predicate() {
}

func (this_ *FunctionTypeBuilder_ImplementedTypeValidator) initFields_FunctionTypeBuilder_ImplementedTypeValidator() {
	_ = this_
	this_.formattedFnName = j2gort.NullString
}

// Java constructor: ImplementedTypeValidator(com.google.javascript.rhino.Node errorRoot, com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String formattedFnName)
func NewFunctionTypeBuilder_ImplementedTypeValidator(errorRoot pkg_rhino.Node, compiler AbstractCompiler, formattedFnName string) *FunctionTypeBuilder_ImplementedTypeValidator {
	this_ := &FunctionTypeBuilder_ImplementedTypeValidator{}
	this_.BindSelf_(this_)
	InitFunctionTypeBuilder_ImplementedTypeValidator(this_, errorRoot, compiler, formattedFnName)
	return this_
}

func InitFunctionTypeBuilder_ImplementedTypeValidator(this_ *FunctionTypeBuilder_ImplementedTypeValidator, errorRoot pkg_rhino.Node, compiler AbstractCompiler, formattedFnName string) {
	_ = errorRoot
	_ = compiler
	_ = formattedFnName
	InitFunctionTypeBuilder_ValidatorBase(&this_.FunctionTypeBuilder_ValidatorBaseBase, errorRoot, compiler)
	this_.initFields_FunctionTypeBuilder_ImplementedTypeValidator()
	this_.formattedFnName = formattedFnName
}

// Java method: ImplementedTypeValidator.apply(com.google.javascript.rhino.jstype.JSType type)
func (this_ *FunctionTypeBuilder_ImplementedTypeValidator) Apply(arg_26 any) bool {
	type_ := j2gort.AsIface[pkg_rhino.JSType](arg_26)
	_ = type_
	var objectType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(type_)
	_ = objectType
	if objectType == nil || objectType.IsJavaNil_() {
		this_.ReportError(StaticGet_TypeCheck_BAD_IMPLEMENTED_TYPE(), []string{this_.formattedFnName})
		return false
	} else if objectType.IsEmptyType() {
		this_.ReportWarning(StaticGet_FunctionTypeBuilder_RESOLVED_TAG_EMPTY(), []string{"@implements", this_.formattedFnName})
		return false
	} else if objectType.IsUnknownType() {
		if FunctionTypeBuilder_HasMoreTagsToResolve(objectType) {
			return true
		} else {
			this_.ReportWarning(StaticGet_FunctionTypeBuilder_RESOLVED_TAG_EMPTY(), []string{"@implements", this_.formattedFnName})
			return false
		}
	} else {
		return true
	}
}

// Java interface: com.google.javascript.jscomp.FunctionTypeBuilder.FunctionContents
type FunctionTypeBuilder_FunctionContents interface {
	Implements_jscomp_FunctionTypeBuilder_FunctionContents()
	IsJavaNil_() bool
	GetSourceNode() pkg_rhino.Node
	MayBeFromExterns() bool
	MayHaveNonEmptyReturns() bool
	MayHaveSingleThrow() bool
}

func ToIface_FunctionTypeBuilder_FunctionContents[T interface {
	FunctionTypeBuilder_FunctionContents
	comparable
}](v T) FunctionTypeBuilder_FunctionContents {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionTypeBuilder.FunctionContents", "com.google.javascript.jscomp.FunctionTypeBuilder$FunctionContents", func(v_ any) bool {
		return j2gort.InstanceOf[FunctionTypeBuilder_FunctionContents](v_)
	})
}

// Java class: com.google.javascript.jscomp.FunctionTypeBuilder.UnknownFunctionContents
type FunctionTypeBuilder_UnknownFunctionContents struct {
	idHash_ atomic.Int32
}

func (this_ *FunctionTypeBuilder_UnknownFunctionContents) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionTypeBuilder_UnknownFunctionContents) IsJavaNil_() bool { return this_ == nil }
func (*FunctionTypeBuilder_UnknownFunctionContents) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionTypeBuilder$UnknownFunctionContents", "UnknownFunctionContents"
}
func (this_ *FunctionTypeBuilder_UnknownFunctionContents) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionTypeBuilder_UnknownFunctionContents) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionTypeBuilder.UnknownFunctionContents", "com.google.javascript.jscomp.FunctionTypeBuilder$UnknownFunctionContents", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionTypeBuilder_UnknownFunctionContents](v_)
	})
}
func (*FunctionTypeBuilder_UnknownFunctionContents) Implements_jscomp_FunctionTypeBuilder_FunctionContents() {
}

var static_FunctionTypeBuilder_UnknownFunctionContents_singleton *FunctionTypeBuilder_UnknownFunctionContents

func StaticGet_FunctionTypeBuilder_UnknownFunctionContents_Singleton() *FunctionTypeBuilder_UnknownFunctionContents {
	clinit_FunctionTypeBuilder_UnknownFunctionContents()
	return static_FunctionTypeBuilder_UnknownFunctionContents_singleton
}
func StaticSet_FunctionTypeBuilder_UnknownFunctionContents_Singleton(v *FunctionTypeBuilder_UnknownFunctionContents) {
	clinit_FunctionTypeBuilder_UnknownFunctionContents()
	static_FunctionTypeBuilder_UnknownFunctionContents_singleton = v
}

var clinitDone_FunctionTypeBuilder_UnknownFunctionContents uint32
var clinitState_FunctionTypeBuilder_UnknownFunctionContents j2gort.ClinitState

func clinit_FunctionTypeBuilder_UnknownFunctionContents() {
	if atomic.LoadUint32(&clinitDone_FunctionTypeBuilder_UnknownFunctionContents) != 0 {
		return
	}
	clinitSlow_FunctionTypeBuilder_UnknownFunctionContents()
}

func clinitSlow_FunctionTypeBuilder_UnknownFunctionContents() {
	if !j2gort.EnterClinit(&clinitState_FunctionTypeBuilder_UnknownFunctionContents, &clinitDone_FunctionTypeBuilder_UnknownFunctionContents) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FunctionTypeBuilder_UnknownFunctionContents, &clinitDone_FunctionTypeBuilder_UnknownFunctionContents)
	static_FunctionTypeBuilder_UnknownFunctionContents_singleton = NewFunctionTypeBuilder_UnknownFunctionContents()
}

func (this_ *FunctionTypeBuilder_UnknownFunctionContents) initFields_FunctionTypeBuilder_UnknownFunctionContents() {
	_ = this_
}

// Java constructor: UnknownFunctionContents()
func NewFunctionTypeBuilder_UnknownFunctionContents() *FunctionTypeBuilder_UnknownFunctionContents {
	clinit_FunctionTypeBuilder_UnknownFunctionContents()
	this_ := &FunctionTypeBuilder_UnknownFunctionContents{}
	InitFunctionTypeBuilder_UnknownFunctionContents(this_)
	return this_
}

func InitFunctionTypeBuilder_UnknownFunctionContents(this_ *FunctionTypeBuilder_UnknownFunctionContents) {
	this_.initFields_FunctionTypeBuilder_UnknownFunctionContents()
}

// Java static method: UnknownFunctionContents.get()
func FunctionTypeBuilder_UnknownFunctionContents_Get() FunctionTypeBuilder_FunctionContents {
	clinit_FunctionTypeBuilder_UnknownFunctionContents()
	return ToIface_FunctionTypeBuilder_FunctionContents(StaticGet_FunctionTypeBuilder_UnknownFunctionContents_Singleton())
}

// Java method: UnknownFunctionContents.getSourceNode()
func (this_ *FunctionTypeBuilder_UnknownFunctionContents) GetSourceNode() pkg_rhino.Node {
	return pkg_rhino.Node(0)
}

// Java method: UnknownFunctionContents.mayBeFromExterns()
func (this_ *FunctionTypeBuilder_UnknownFunctionContents) MayBeFromExterns() bool {
	return true
}

// Java method: UnknownFunctionContents.mayHaveNonEmptyReturns()
func (this_ *FunctionTypeBuilder_UnknownFunctionContents) MayHaveNonEmptyReturns() bool {
	return true
}

// Java method: UnknownFunctionContents.mayHaveSingleThrow()
func (this_ *FunctionTypeBuilder_UnknownFunctionContents) MayHaveSingleThrow() bool {
	return true
}

// Java class: com.google.javascript.jscomp.FunctionTypeBuilder.AstFunctionContents
type FunctionTypeBuilder_AstFunctionContents struct {
	n                  pkg_rhino.Node `gson:"n"`
	idHash_            atomic.Int32
	hasNonEmptyReturns bool `gson:"hasNonEmptyReturns"`
}

func (this_ *FunctionTypeBuilder_AstFunctionContents) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionTypeBuilder_AstFunctionContents) IsJavaNil_() bool { return this_ == nil }
func (*FunctionTypeBuilder_AstFunctionContents) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionTypeBuilder$AstFunctionContents", "AstFunctionContents"
}
func (this_ *FunctionTypeBuilder_AstFunctionContents) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionTypeBuilder_AstFunctionContents) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionTypeBuilder.AstFunctionContents", "com.google.javascript.jscomp.FunctionTypeBuilder$AstFunctionContents", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionTypeBuilder_AstFunctionContents](v_)
	})
}
func (*FunctionTypeBuilder_AstFunctionContents) Implements_jscomp_FunctionTypeBuilder_FunctionContents() {
}

func (this_ *FunctionTypeBuilder_AstFunctionContents) initFields_FunctionTypeBuilder_AstFunctionContents() {
	_ = this_

}

// Java constructor: AstFunctionContents(com.google.javascript.rhino.Node n)
func NewFunctionTypeBuilder_AstFunctionContents(n pkg_rhino.Node) *FunctionTypeBuilder_AstFunctionContents {
	this_ := &FunctionTypeBuilder_AstFunctionContents{}
	InitFunctionTypeBuilder_AstFunctionContents(this_, n)
	return this_
}

func InitFunctionTypeBuilder_AstFunctionContents(this_ *FunctionTypeBuilder_AstFunctionContents, n pkg_rhino.Node) {
	_ = n
	this_.initFields_FunctionTypeBuilder_AstFunctionContents()
	this_.n = n
}

// Java method: AstFunctionContents.getSourceNode()
func (this_ *FunctionTypeBuilder_AstFunctionContents) GetSourceNode() pkg_rhino.Node {
	return this_.n
}

// Java method: AstFunctionContents.mayBeFromExterns()
func (this_ *FunctionTypeBuilder_AstFunctionContents) MayBeFromExterns() bool {
	return this_.n.IsFromExterns()
}

// Java method: AstFunctionContents.mayHaveNonEmptyReturns()
func (this_ *FunctionTypeBuilder_AstFunctionContents) MayHaveNonEmptyReturns() bool {
	return this_.hasNonEmptyReturns
}

// Java method: AstFunctionContents.recordNonEmptyReturn()
func (this_ *FunctionTypeBuilder_AstFunctionContents) RecordNonEmptyReturn() {
	this_.hasNonEmptyReturns = true
}

// Java method: AstFunctionContents.mayHaveSingleThrow()
func (this_ *FunctionTypeBuilder_AstFunctionContents) MayHaveSingleThrow() bool {
	var block pkg_rhino.Node = this_.n.GetLastChild()
	_ = block
	return (block.HasOneChild() && block.GetFirstChild().IsThrow())
}
