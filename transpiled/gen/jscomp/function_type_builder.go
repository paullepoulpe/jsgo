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
	reusableParamBuilder         *pkg_rhino.FunctionParamBuilder                    `gson:"reusableParamBuilder"`
	reusableFnBuilder            *pkg_rhino.FunctionType_Builder                    `gson:"reusableFnBuilder"`
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
	contentsSourceNode           pkg_rhino.Node                                     `gson:"contentsSourceNode"`
	idHash_                      atomic.Int32
	contentsHasNonEmptyReturns   bool                       `gson:"contentsHasNonEmptyReturns"`
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

var static_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_FunctionTypeBuilder()
	return static_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES
}
func StaticSet_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_FunctionTypeBuilder()
	static_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES = v
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
	static_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
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

	this_.syntacticFnName = j2gort.NullString

	this_.templateTypeNames = StaticGet_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES()
	this_.constructorTemplateTypeNames = StaticGet_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES()

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
	this_.reusableParamBuilder = pkg_rhino.NewFunctionParamBuilder_1(this_.typeRegistry)
	this_.reusableFnBuilder = pkg_rhino.FunctionType_S_Builder(this_.typeRegistry)
}

// Java method: FunctionTypeBuilder.createExtendedTypeValidator()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_createExtendedTypeValidator_2981() *FunctionTypeBuilder_ExtendedTypeValidator {
	return NewFunctionTypeBuilder_ExtendedTypeValidator(this_.errorRoot, this_.compiler, this_.FormatFnName())
}

// Java method: FunctionTypeBuilder.createImplementedTypeValidator()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_createImplementedTypeValidator_2982() *FunctionTypeBuilder_ImplementedTypeValidator {
	return NewFunctionTypeBuilder_ImplementedTypeValidator(this_.errorRoot, this_.compiler, this_.FormatFnName())
}

// Java method: FunctionTypeBuilder.reset(java.lang.String fnName, com.google.javascript.rhino.Node errorRoot, com.google.javascript.jscomp.TypedScope scope)
func (this_ *FunctionTypeBuilder) Reset(fnName string, errorRoot pkg_rhino.Node, scope *TypedScope) {
	_ = fnName
	_ = errorRoot
	_ = scope
	j2gort.CheckNotNullNode(errorRoot)
	this_.fnName = j2gort.Ext_com_google_common_base_Strings_NullToEmpty(fnName)
	this_.errorRoot = errorRoot
	this_.contentsSourceNode = pkg_rhino.Node(0)
	this_.contentsHasNonEmptyReturns = false
	this_.syntacticFnName = j2gort.NullString
	this_.returnType = nil
	this_.returnTypeInferred = false
	this_.implementedInterfaces = nil
	this_.extendedInterfaces = nil
	this_.baseType = nil
	this_.thisType = nil
	this_.isClass = false
	this_.isConstructor = false
	this_.makesStructs = false
	this_.makesUnrestricted = false
	this_.makesDicts = false
	this_.isInterface = false
	this_.isRecord = false
	this_.isAbstract = false
	this_.isKnownAmbiguous = false
	this_.parameters = nil
	this_.closurePrimitiveId = pkg_rhino.ClosurePrimitive(0)
	this_.templateTypeNames = StaticGet_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES()
	this_.constructorTemplateTypeNames = StaticGet_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES()
	this_.declarationScope = nil
	this_.templateScope = pkg_rhino.ToIface_StaticTypedScope(scope)
}

// Java method: FunctionTypeBuilder.clearForPool()
func (this_ *FunctionTypeBuilder) ClearForPool() {
	this_.fnName = j2gort.NullString
	this_.errorRoot = pkg_rhino.Node(0)
	this_.contentsSourceNode = pkg_rhino.Node(0)
	this_.contentsHasNonEmptyReturns = false
	this_.syntacticFnName = j2gort.NullString
	this_.returnType = nil
	this_.implementedInterfaces = nil
	this_.extendedInterfaces = nil
	this_.baseType = nil
	this_.thisType = nil
	this_.parameters = nil
	this_.closurePrimitiveId = pkg_rhino.ClosurePrimitive(0)
	this_.templateTypeNames = StaticGet_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES()
	this_.constructorTemplateTypeNames = StaticGet_FunctionTypeBuilder_EMPTY_TEMPLATE_TYPES()
	this_.declarationScope = nil
	this_.templateScope = nil
	this_.reusableParamBuilder.Reset()
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

// Java method: FunctionTypeBuilder.setAstContents(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node sourceNode, boolean hasNonEmptyReturns)
func (this_ *FunctionTypeBuilder) SetAstContents(sourceNode pkg_rhino.Node, hasNonEmptyReturns bool) *FunctionTypeBuilder {
	_ = sourceNode
	_ = hasNonEmptyReturns
	this_.contentsSourceNode = sourceNode
	this_.contentsHasNonEmptyReturns = hasNonEmptyReturns
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
			this_.parameters = pkg_rhino.StaticGet_FunctionType_EMPTY_PARAMETERS()
		}
	} else {
		var paramBuilder *pkg_rhino.FunctionParamBuilder = this_.reusableParamBuilder
		_ = paramBuilder
		paramBuilder.Reset()
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
					{

						ba_3 := warnedAboutArgList
						bb_4 := this_.priv_FunctionTypeBuilder_addParameter_3009(paramBuilder, this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE()), warnedAboutArgList, ((this_.codingConvention.IsOptionalParameter(currentParam) || oldParamsListHitOptArgs) || currentParam.IsDefaultValue()), this_.codingConvention.IsVarArgsParameter(currentParam))
						warnedAboutArgList = ba_3 || bb_4
					}

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
		if !(returnTypeExpr).IsJavaNil_() {
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
		this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED(), []string{"@struct", this_.FormatFnName()})
	} else if this_.makesDicts && !(this_.isConstructor) {
		this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED(), []string{"@dict", this_.FormatFnName()})
	}
	return this_
}

// Java method: FunctionTypeBuilder.maybeUseNativeClassTemplateNames(com.google.javascript.rhino.JSDocInfo info)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_maybeUseNativeClassTemplateNames_2993(info *pkg_rhino.JSDocInfo) bool {
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
			if !(FunctionTypeBuilder_AreCompatibleExtendsTypes(infoBaseType.GetJSTypeBase().ToMaybeObjectType(), classExtendsType)) {
				this_.isKnownAmbiguous = true
			}
			if infoBaseType.SetValidator(j2gort.ToIface_Ext_com_google_common_base_Predicate(this_.priv_FunctionTypeBuilder_createExtendedTypeValidator_2981())) {
				this_.baseType = infoBaseType.GetJSTypeBase().ToObjectType()
			} else {
				this_.isKnownAmbiguous = true
			}
		} else {
			this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF(), []string{this_.FormatFnName()})
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
					if maybeInterType.SetValidator(j2gort.ToIface_Ext_com_google_common_base_Predicate(this_.priv_FunctionTypeBuilder_createImplementedTypeValidator_2982())) {
						this_.implementedInterfaces.Add_1(j2gort.AsIface[pkg_rhino.ObjectType](maybeInterType))
					}
				}
			}
		} else if this_.isInterface {
			this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_TypeCheck_CONFLICTING_IMPLEMENTED_TYPE(), []string{this_.FormatFnName()})
		} else {
			this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_FunctionTypeBuilder_CONSTRUCTOR_REQUIRED(), []string{"@implements", this_.FormatFnName()})
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
						var isValid bool = maybeInterfaceType.SetValidator(j2gort.ToIface_Ext_com_google_common_base_Predicate(this_.priv_FunctionTypeBuilder_createExtendedTypeValidator_2981()))
						_ = isValid
						if isValid && !(j2gort.IfaceEquals[pkg_rhino.ObjectType](maybeInterfaceType.GetJSTypeBase().ToMaybeObjectType(), nil)) {
							this_.extendedInterfaces.Add_1(maybeInterfaceType.GetJSTypeBase().ToMaybeObjectType())
						}
					}
					if (classExtendsType != nil && !classExtendsType.IsJavaNil_()) && maybeInterfaceType.GetJSTypeBase().IsSubtypeOf_1(pkg_rhino.Upcast_JSType(classExtendsType)) {
						classExtendsType = nil
					}
				}
			}
		}
		if (classExtendsType != nil && !classExtendsType.IsJavaNil_()) && classExtendsType.SetValidator(j2gort.ToIface_Ext_com_google_common_base_Predicate(this_.priv_FunctionTypeBuilder_createExtendedTypeValidator_2981())) {
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
	return (annotated.GetObjectTypeBase().IsTemplatizedType() && j2gort.IfaceEquals[pkg_rhino.ObjectType](annotated.ToMaybeTemplatizedType().GetReferencedType(), extendsClause))
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
			this_.thisType = pkg_rhino.Upcast_JSType(objType)
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
		var i int32 = int32(0)
		_ = i
		var paramCount int32 = info.GetParameterCount()
		_ = paramCount
		for ; i < paramCount; i = i + 1 {
			lp.AddChildToBack(pkg_rhino.IR_Name(info.GetParameterNameAt(i)))
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
		oldParameterType = oldParameters[func() int32 { old_7 := oldParamIdx; nw_8 := old_7 + 1; oldParamIdx = nw_8; return old_7 }()]
	} else {
		oldParameterType = (pkg_rhino.FunctionType_Parameter{})
	}
	var builder *pkg_rhino.FunctionParamBuilder = this_.reusableParamBuilder
	_ = builder
	builder.Reset()
	var warnedAboutArgList bool = false
	_ = warnedAboutArgList
	var jsDocParamCount int32
	_ = jsDocParamCount
	if info == nil {
		jsDocParamCount = int32(0)
	} else {
		jsDocParamCount = info.GetParameterCount()
	}
	var matchedJsDocParamsMask int64 = int64(0)
	_ = matchedJsDocParamsMask
	var matchedJsDocParamsExtra []bool
	_ = matchedJsDocParamsExtra
	if jsDocParamCount > int32(64) {
		matchedJsDocParamsExtra = make([]bool, jsDocParamCount)
	} else {
		matchedJsDocParamsExtra = nil
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
				isVarArgs = this_.priv_FunctionTypeBuilder_isVarArgsParameterByConvention_3006(param)
				isOptionalParam = this_.priv_FunctionTypeBuilder_isOptionalParameterByConvention_3005(param)
				paramLhs = param
			}
			var paramName string = j2gort.NullString
			_ = paramName
			var jsDocParamIdx int32 = int32(-1)
			_ = jsDocParamIdx
			if paramLhs.IsName() {
				paramName = paramLhs.GetString()
				if jsDocParamCount > int32(0) {
					jsDocParamIdx = info.GetParameterIndex(paramName)
				}
			} else {
				if !(paramLhs.IsDestructuringPattern()) {
					j2gort.CheckState(false)
				}
				if paramIndex < jsDocParamCount {
					paramName = info.GetParameterNameAt(paramIndex)
					jsDocParamIdx = paramIndex
				}
			}
			var jsDocParamTypeExpr pkg_rhino.JSTypeExpression = (pkg_rhino.JSTypeExpression{})
			_ = jsDocParamTypeExpr
			if jsDocParamIdx >= int32(0) {
				if jsDocParamIdx < int32(64) {
					matchedJsDocParamsMask = (matchedJsDocParamsMask | (javalang.Shl64(int64(1), int(jsDocParamIdx))))
				} else {
					{
						a_9 := matchedJsDocParamsExtra
						i_10 := jsDocParamIdx
						a_9[i_10] = true
					}
				}
				jsDocParamTypeExpr = info.GetParameterTypeAt(jsDocParamIdx)
			}
			var parameterType pkg_rhino.JSType = nil
			_ = parameterType
			if !(jsDocParamTypeExpr).IsJavaNil_() {
				parameterType = jsDocParamTypeExpr.Evaluate(this_.templateScope, this_.typeRegistry)
				isOptionalParam = (isOptionalParam || jsDocParamTypeExpr.IsOptionalArg())
				isVarArgs = (isVarArgs || jsDocParamTypeExpr.IsVarArgs())
			} else if (paramLhs.GetJSDocInfo() != nil) && paramLhs.GetJSDocInfo().HasType() {
				var parameterTypeExpression pkg_rhino.JSTypeExpression = paramLhs.GetJSDocInfo().GetType()
				_ = parameterTypeExpression
				parameterType = parameterTypeExpression.Evaluate(this_.templateScope, this_.typeRegistry)
				isOptionalParam = parameterTypeExpression.IsOptionalArg()
				isVarArgs = parameterTypeExpression.IsVarArgs()
			} else if (oldParameterType != (pkg_rhino.FunctionType_Parameter{})) && !(j2gort.IfaceEquals[pkg_rhino.JSType](oldParameterType.GetJSType(), nil)) {
				parameterType = oldParameterType.GetJSType()
				isOptionalParam = oldParameterType.IsOptional()
				isVarArgs = oldParameterType.IsVariadic()
			} else {
				parameterType = this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
			}
			{

				ba_11 := warnedAboutArgList
				bb_12 := this_.priv_FunctionTypeBuilder_addParameter_3009(builder, parameterType, warnedAboutArgList, isOptionalParam, isVarArgs)
				warnedAboutArgList = ba_11 || bb_12
			}

			if oldParamIdx < oldParamLen {
				oldParameterType = oldParameters[func() int32 { old_13 := oldParamIdx; nw_14 := old_13 + 1; oldParamIdx = nw_14; return old_13 }()]
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
				oldParameterType = oldParameters[func() int32 { old_15 := oldParamIdx; nw_16 := old_15 + 1; oldParamIdx = nw_16; return old_15 }()]
			} else {
				oldParameterType = (pkg_rhino.FunctionType_Parameter{})
			}
		}
	}
	if jsDocParamCount > int32(0) {
		var expectedMask int64
		_ = expectedMask
		if jsDocParamCount >= int32(64) {
			expectedMask = int64(-1)
		} else {
			expectedMask = ((javalang.Shl64(int64(1), int(jsDocParamCount))) - int64(1))
		}
		if (matchedJsDocParamsMask != expectedMask) || (matchedJsDocParamsExtra != nil) {
			{
				var i int32 = int32(0)
				_ = i
				for ; i < jsDocParamCount; i = i + 1 {
					var matched bool
					_ = matched
					if i < int32(64) {
						matched = ((matchedJsDocParamsMask & (javalang.Shl64(int64(1), int(i)))) != int64(0))
					} else {
						matched = matchedJsDocParamsExtra[i]
					}
					if !(matched) {
						this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_FunctionTypeBuilder_INEXISTENT_PARAM(), []string{info.GetParameterNameAt(i), this_.FormatFnName()})
					}
				}
			}
		}
	}
	this_.parameters = builder.BuildArray()
	return this_
}

// Java method: FunctionTypeBuilder.registerTemplates(java.lang.Iterable<com.google.javascript.rhino.jstype.TemplateType> templates, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node scopeRoot)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_registerTemplates_2999(templates j2gort.Ext_java_lang_Iterable, scopeRoot pkg_rhino.Node) {
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
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_registerTemplatesWithOwner_3000(templates j2gort.Ext_com_google_common_collect_ImmutableList, ownerMap *pkg_rhino.TemplateTypeMap, scopeRoot pkg_rhino.Node) {
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
		this_.priv_FunctionTypeBuilder_setConstructorTemplateTypeNames_3004(j2gort.AsIface[j2gort.Ext_java_util_List](this_.priv_FunctionTypeBuilder_buildTemplateTypesFromJSDocInfo_3007(info, true)), argsParent.GetParent())
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
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_setConstructorTemplateTypeNames_3004(templates j2gort.Ext_java_util_List, ctor pkg_rhino.Node) {
	_ = templates
	_ = ctor
	if !(templates.IsEmpty()) {
		this_.constructorTemplateTypeNames = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](templates))
		if this_.templateTypeNames.IsEmpty() {
			this_.templateTypeNames = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](templates))
		} else {
			this_.templateTypeNames = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}).AddAllSelf_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.templateTypeNames)).AddAllSelf_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.constructorTemplateTypeNames)).BuildImmList_()
		}
		this_.priv_FunctionTypeBuilder_registerTemplates_2999(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](templates), ctor)
	}
}

// Java method: FunctionTypeBuilder.isOptionalParameterByConvention(com.google.javascript.rhino.Node param)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_isOptionalParameterByConvention_3005(param pkg_rhino.Node) bool {
	_ = param
	if param.IsDestructuringPattern() {
		return false
	}
	return this_.codingConvention.IsOptionalParameter(param)
}

// Java method: FunctionTypeBuilder.isVarArgsParameterByConvention(com.google.javascript.rhino.Node param)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_isVarArgsParameterByConvention_3006(param pkg_rhino.Node) bool {
	_ = param
	if param.IsDestructuringPattern() {
		return false
	}
	return this_.codingConvention.IsVarArgsParameter(param)
}

// Java method: FunctionTypeBuilder.buildTemplateTypesFromJSDocInfo(com.google.javascript.rhino.JSDocInfo info, boolean allowTypeTransformations)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_buildTemplateTypesFromJSDocInfo_3007(info *pkg_rhino.JSDocInfo, allowTypeTransformations bool) j2gort.Ext_com_google_common_collect_ImmutableList {
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
				this_.priv_FunctionTypeBuilder_reportError_3017(StaticGet_FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND(), []string{j2gort.StrVal(typeBound), cur_18.KeyStr()})
			}
			var template *pkg_rhino.TemplateType = this_.typeRegistry.GetType_2(pkg_rhino.Upcast_StaticScope(this_.templateScope), cur_18.KeyStr()).ToMaybeTemplateType()
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
				this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS(), []string{cur_20.KeyStr()})
			}
		}
	}
	var builtTemplates []*pkg_rhino.TemplateType = templates
	_ = builtTemplates
	{

		for _, template_3 := range builtTemplates {

			_ = template_3
			if template_3.ContainsCycle() {
				this_.priv_FunctionTypeBuilder_reportError_3017(StaticGet_RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR(), []string{("Cycle detected in inheritance chain of type " + j2gort.StrValString(template_3.GetReferenceName()))})
			}
		}
	}
	return j2gort.ImmListFromTypedSlice(builtTemplates)
}

// Java method: FunctionTypeBuilder.inferTemplateTypeName(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType ownerType)
func (this_ *FunctionTypeBuilder) InferTemplateTypeName(info *pkg_rhino.JSDocInfo, ownerType pkg_rhino.JSType) *FunctionTypeBuilder {
	_ = info
	_ = ownerType
	if (info != nil) && !(this_.priv_FunctionTypeBuilder_maybeUseNativeClassTemplateNames_2993(info)) {
		var templates j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_FunctionTypeBuilder_buildTemplateTypesFromJSDocInfo_3007(info, !(this_.isConstructor || this_.isInterface))
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
		this_.priv_FunctionTypeBuilder_registerTemplatesWithOwner_3000(this_.templateTypeNames, ownerTypeMap, this_.contentsSourceNode)
	}
	return this_
}

// Java method: FunctionTypeBuilder.addParameter(com.google.javascript.rhino.jstype.FunctionParamBuilder builder, com.google.javascript.rhino.jstype.JSType paramType, boolean warnedAboutArgList, boolean isOptional, boolean isVarArgs)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_addParameter_3009(builder *pkg_rhino.FunctionParamBuilder, paramType pkg_rhino.JSType, warnedAboutArgList bool, isOptional bool, isVarArgs bool) bool {
	_ = builder
	_ = paramType
	_ = warnedAboutArgList
	_ = isOptional
	_ = isVarArgs
	var emittedWarning bool = false
	_ = emittedWarning
	if isOptional {
		if !(builder.AddOptionalParam(paramType)) && !(warnedAboutArgList) {
			this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST(), []string{})
			emittedWarning = true
		}
	} else if isVarArgs {
		if !(builder.AddVarArgs(paramType)) && !(warnedAboutArgList) {
			this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST(), []string{})
			emittedWarning = true
		}
	} else {
		if !(builder.AddRequiredParam(paramType)) && !(warnedAboutArgList) {
			if builder.HasVarArgs() {
				this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST(), []string{})
			} else {
				this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_FunctionTypeBuilder_OPTIONAL_ARG_AT_END(), []string{})
			}
			emittedWarning = true
		}
	}
	return emittedWarning
}

// Java static method: FunctionTypeBuilder.mayHaveSingleThrow(com.google.javascript.rhino.Node n)
func FunctionTypeBuilder_MayHaveSingleThrow(n pkg_rhino.Node) bool {
	clinit_FunctionTypeBuilder()
	_ = n
	var block pkg_rhino.Node = n.GetLastChild()
	_ = block
	return (block.HasOneChild() && block.GetFirstChild().IsThrow())
}

// Java method: FunctionTypeBuilder.provideDefaultReturnType()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_provideDefaultReturnType_3010() {
	var srcNode pkg_rhino.Node = this_.contentsSourceNode
	_ = srcNode
	if (srcNode != 0) && srcNode.IsAsyncGeneratorFunction() {
		var generatorType pkg_rhino.ObjectType = this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ASYNC_GENERATOR_TYPE())
		_ = generatorType
		this_.returnType = pkg_rhino.ToIface_JSType(this_.typeRegistry.CreateTemplatizedType_2_ObjectType_JSType(generatorType, []pkg_rhino.JSType{this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())}))
		return
	} else if (srcNode != 0) && srcNode.IsGeneratorFunction() {
		var generatorType_2 pkg_rhino.ObjectType = this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_GENERATOR_TYPE())
		_ = generatorType_2
		this_.returnType = pkg_rhino.ToIface_JSType(this_.typeRegistry.CreateTemplatizedType_2_ObjectType_JSType(generatorType_2, []pkg_rhino.JSType{this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())}))
		return
	}
	var inferredReturnType pkg_rhino.JSType = this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	_ = inferredReturnType
	if (((srcNode != 0) && !(this_.contentsHasNonEmptyReturns)) && !(FunctionTypeBuilder_MayHaveSingleThrow(srcNode))) && !(srcNode.IsFromExterns()) {
		inferredReturnType = this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
		this_.returnTypeInferred = true
	}
	if (srcNode != 0) && srcNode.IsAsyncFunction() {
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
		this_.priv_FunctionTypeBuilder_provideDefaultReturnType_3010()
		j2gort.CheckNotNullIface(this_.returnType)
	}
	if this_.parameters == nil {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "All Function types must have params and a return type")))
	}
	var fnType *pkg_rhino.FunctionType = nil
	_ = fnType
	if this_.isConstructor {
		fnType = this_.priv_FunctionTypeBuilder_getOrCreateConstructor_3014()
	} else if this_.isInterface {
		fnType = this_.priv_FunctionTypeBuilder_getOrCreateInterface_3015()
	} else {
		fnType = j2gort.CastPtr[pkg_rhino.FunctionType](this_.priv_FunctionTypeBuilder_createDefaultBuilder_3013().WithParameterArray(this_.parameters).WithReturnType_2(this_.returnType, this_.returnTypeInferred).WithTypeOfThis(this_.thisType).WithIsAbstract(this_.isAbstract).WithClosurePrimitiveId(this_.closurePrimitiveId).Build())
		this_.priv_FunctionTypeBuilder_maybeSetBaseType_3012(fnType)
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
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_maybeSetBaseType_3012(fnType *pkg_rhino.FunctionType) {
	_ = fnType
	if !(fnType.HasInstanceType()) || (this_.baseType == nil || this_.baseType.IsJavaNil_()) {
		return
	}
	fnType.SetPrototypeBasedOn_1(this_.baseType)
	fnType.GetInstanceType().GetObjectTypeBase().MergeSupertypeTemplateTypes(this_.baseType)
}

// Java method: FunctionTypeBuilder.createDefaultBuilder()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_createDefaultBuilder_3013() *pkg_rhino.FunctionType_Builder {
	return this_.reusableFnBuilder.Reset().WithName(this_.fnName).WithSourceNode(this_.contentsSourceNode).WithTemplateKeys_1_Ext_com_google_common_collect_ImmutableList(this_.templateTypeNames).SetIsKnownAmbiguous(this_.isKnownAmbiguous).SetGoogModuleId(TypedScopeCreator_ContainingGoogModuleIdOf(this_.declarationScope))
}

// Java method: FunctionTypeBuilder.getOrCreateConstructor()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_getOrCreateConstructor_3014() *pkg_rhino.FunctionType {
	var fnType *pkg_rhino.FunctionType = j2gort.CastPtr[pkg_rhino.FunctionType](this_.priv_FunctionTypeBuilder_createDefaultBuilder_3013().ForConstructor().WithParameterArray(this_.parameters).WithReturnType_1(this_.returnType).WithConstructorTemplateKeys(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.constructorTemplateTypeNames)).WithIsAbstract(this_.isAbstract).Build())
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
				existingFn = existingType.GetJSTypeBase().ToObjectType().GetConstructor()
			} else {
				existingFn = this_.typeRegistry.GetNativeFunctionType(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_FUNCTION_TYPE())
			}
			if existingFn.GetSource() == 0 {
				existingFn.SetSource(this_.contentsSourceNode)
			}
			if !(existingFn.HasEqualCallType(fnType)) {
				this_.priv_FunctionTypeBuilder_reportWarning_3016(StaticGet_FunctionTypeBuilder_TYPE_REDEFINITION(), []string{this_.FormatFnName(), fnType.String(), existingFn.String()})
			}
			return existingFn
		} else {
		}
	}
	this_.priv_FunctionTypeBuilder_maybeSetBaseType_3012(fnType)
	if !(j2gort.CheckNotNullString(this_.syntacticFnName) == "") && !(javalang.StringStartsWith(j2gort.CheckNotNullString(this_.syntacticFnName), "this.")) {
		this_.typeRegistry.DeclareTypeForExactScope(pkg_rhino.ToIface_StaticScope(this_.declarationScope), this_.syntacticFnName, pkg_rhino.Upcast_JSType(fnType.GetInstanceType()))
	}
	return fnType
}

// Java method: FunctionTypeBuilder.getOrCreateInterface()
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_getOrCreateInterface_3015() *pkg_rhino.FunctionType {
	var fnType *pkg_rhino.FunctionType = nil
	_ = fnType
	var type_ pkg_rhino.JSType = this_.typeRegistry.GetType_2(pkg_rhino.ToIface_StaticScope(this_.declarationScope), this_.syntacticFnName)
	_ = type_
	if (type_ != nil && !type_.IsJavaNil_()) && type_.IsInstanceType() {
		var ctor *pkg_rhino.FunctionType = type_.GetJSTypeBase().ToMaybeObjectType().GetConstructor()
		_ = ctor
		if ctor.IsInterface() {
			fnType = ctor
			fnType.SetSource(this_.contentsSourceNode)
		}
	}
	if fnType == nil {
		fnType = j2gort.CastPtr[pkg_rhino.FunctionType](this_.priv_FunctionTypeBuilder_createDefaultBuilder_3013().ForInterface().WithParameters_0().Build())
		if this_.makesStructs {
			fnType.SetStruct()
		}
		if !(j2gort.CheckNotNullString(this_.fnName) == "") {
			this_.typeRegistry.DeclareTypeForExactScope(pkg_rhino.ToIface_StaticScope(this_.declarationScope), this_.syntacticFnName, pkg_rhino.Upcast_JSType(fnType.GetInstanceType()))
		}
		this_.priv_FunctionTypeBuilder_maybeSetBaseType_3012(fnType)
	}
	return fnType
}

// Java method: FunctionTypeBuilder.reportWarning(com.google.javascript.jscomp.DiagnosticType warning, java.lang.String[] args)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_reportWarning_3016(warning *DiagnosticType, args []string) {
	_ = warning
	_ = args
	this_.compiler.Report_1(JSError_Make_3(this_.errorRoot, warning, args))
}

// Java method: FunctionTypeBuilder.reportError(com.google.javascript.jscomp.DiagnosticType error, java.lang.String[] args)
func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_reportError_3017(error_ *DiagnosticType, args []string) {
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
	if !(objectType.GetObjectTypeBase().IsUnknownType()) {
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
				if !(interfaceType.GetObjectTypeBase().IsResolved()) {
					return true
				}
			}
		}
	}
	if !(j2gort.IfaceEquals[pkg_rhino.ObjectType](objectType.GetImplicitPrototype(), nil)) {
		return !(objectType.GetImplicitPrototype().GetObjectTypeBase().IsResolved())
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

func Upcast_FunctionTypeBuilder_ValidatorBase[T FunctionTypeBuilder_ValidatorBase](v T) FunctionTypeBuilder_ValidatorBase {
	var vi FunctionTypeBuilder_ValidatorBase = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
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
		this_.ReportWarning(StaticGet_FunctionTypeBuilder_EXTENDS_NON_OBJECT(), []string{this_.formattedFnName, type_.GetJSTypeBase().String()})
		return false
	}
	base_objectType_ := objectType.GetObjectTypeBase()
	if base_objectType_.IsEmptyType() {
		this_.ReportWarning(StaticGet_FunctionTypeBuilder_RESOLVED_TAG_EMPTY(), []string{"@extends", this_.formattedFnName})
		return false
	}
	if base_objectType_.IsUnknownType() {
		if FunctionTypeBuilder_HasMoreTagsToResolve(objectType) || type_.GetJSTypeBase().IsTemplateType() {
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
	} else if objectType.GetObjectTypeBase().IsEmptyType() {
		this_.ReportWarning(StaticGet_FunctionTypeBuilder_RESOLVED_TAG_EMPTY(), []string{"@implements", this_.formattedFnName})
		return false
	} else if objectType.GetObjectTypeBase().IsUnknownType() {
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
