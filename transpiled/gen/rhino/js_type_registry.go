/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/JSTypeRegistry.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.JSTypeRegistry
type JSTypeRegistry struct {
	iObjectIndexTemplateKey         *TemplateType                                                  `gson:"iObjectIndexTemplateKey"`
	iObjectElementTemplateKey       *TemplateType                                                  `gson:"iObjectElementTemplateKey"`
	iterableValueTemplate           *TemplateType                                                  `gson:"iterableValueTemplate"`
	iterableReturnTemplate          *TemplateType                                                  `gson:"iterableReturnTemplate"`
	iteratorIterableValueTemplate   *TemplateType                                                  `gson:"iteratorIterableValueTemplate"`
	iteratorIterableReturnTemplate  *TemplateType                                                  `gson:"iteratorIterableReturnTemplate"`
	iteratorLikeValueTemplate       *TemplateType                                                  `gson:"iteratorLikeValueTemplate"`
	iteratorLikeReturnTemplate      *TemplateType                                                  `gson:"iteratorLikeReturnTemplate"`
	iteratorValueTemplate           *TemplateType                                                  `gson:"iteratorValueTemplate"`
	iteratorReturnTemplate          *TemplateType                                                  `gson:"iteratorReturnTemplate"`
	iiterableResultValueTemplate    *TemplateType                                                  `gson:"iiterableResultValueTemplate"`
	asyncIterableValueTemplate      *TemplateType                                                  `gson:"asyncIterableValueTemplate"`
	asyncIterableReturnTemplate     *TemplateType                                                  `gson:"asyncIterableReturnTemplate"`
	asyncIteratorValueTemplate      *TemplateType                                                  `gson:"asyncIteratorValueTemplate"`
	asyncIteratorReturnTemplate     *TemplateType                                                  `gson:"asyncIteratorReturnTemplate"`
	generatorValueTemplate          *TemplateType                                                  `gson:"generatorValueTemplate"`
	generatorReturnTemplate         *TemplateType                                                  `gson:"generatorReturnTemplate"`
	iThenableTemplateKey            *TemplateType                                                  `gson:"iThenableTemplateKey"`
	promiseTemplateKey              *TemplateType                                                  `gson:"promiseTemplateKey"`
	arrayElementTemplateKey         *TemplateType                                                  `gson:"arrayElementTemplateKey"`
	readonlyArrayElementTemplateKey *TemplateType                                                  `gson:"readonlyArrayElementTemplateKey"`
	readonlyMapKeyTemplateKey       *TemplateType                                                  `gson:"readonlyMapKeyTemplateKey"`
	readonlyMapValueTemplateKey     *TemplateType                                                  `gson:"readonlyMapValueTemplateKey"`
	mapKeyTemplateKey               *TemplateType                                                  `gson:"mapKeyTemplateKey"`
	mapValueTemplateKey             *TemplateType                                                  `gson:"mapValueTemplateKey"`
	reporter                        ErrorReporter                                                  `gson:"-"`
	nativeTypes                     []JSType                                                       `gson:"nativeTypes"`
	defaultVarArgsParameterArray    []FunctionType_Parameter                                       `gson:"defaultVarArgsParameterArray"`
	unknownRequiredParameter        FunctionType_Parameter                                         `gson:"unknownRequiredParameter"`
	unknownRequiredParameterArrays  [][]FunctionType_Parameter                                     `gson:"unknownRequiredParameterArrays"`
	scopedNameTable                 *JSTypeRegistry_ScopedNameTable                                `gson:"scopedNameTable"`
	qualifiedScopedNameHashes       []int32                                                        `gson:"qualifiedScopedNameHashes"`
	closureNamespaces               j2gort.OrderedMap[string, JSTypeRegistry_ClosureNamespace]     `gson:"-"`
	nonLegacyClosureNamespaceRoots  j2gort.OrderedSet[string]                                      `gson:"-"`
	nonNullableTypeNames            *JSTypeRegistry_ScopedNameSet                                  `gson:"nonNullableTypeNames"`
	union2A                         []JSType                                                       `gson:"union2A"`
	union2B                         []JSType                                                       `gson:"union2B"`
	union2Res                       []JSType                                                       `gson:"union2Res"`
	unionBuilderPool                []*UnionType_Builder                                           `gson:"unionBuilderPool"`
	replacerPool                    []*TemplateTypeReplacer                                        `gson:"replacerPool"`
	forwardDeclaredTypes            j2gort.Ext_java_util_Set                                       `gson:"-"`
	nonRefTypesIndexedByProperty    j2gort.OrderedMap[string, *j2gort.Ext_java_util_LinkedHashSet] `gson:"-"`
	sentinelObjectLiteral           JSType                                                         `gson:"sentinelObjectLiteral"`
	propertiesOfSupertypesInUnions  *j2gort.Ext_java_util_LinkedHashSet                            `gson:"propertiesOfSupertypesInUnions"`
	droppedPropertiesOfUnions       j2gort.OrderedSet[string]                                      `gson:"droppedPropertiesOfUnions"`
	eachRefTypeIndexedByProperty    j2gort.OrderedMap[string, *j2gort.Ext_java_util_LinkedHashSet] `gson:"eachRefTypeIndexedByProperty"`
	emptyTemplateTypeMap            *TemplateTypeMap                                               `gson:"emptyTemplateTypeMap"`
	resolver                        *JSTypeResolver                                                `gson:"resolver"`
	qualifiedScopedNameHashCount    int32                                                          `gson:"qualifiedScopedNameHashCount"`
	nameTableGlobalRoot             Node                                                           `gson:"nameTableGlobalRoot"`
	nextTypeId                      int32                                                          `gson:"nextTypeId"`
	hashDepth                       int32                                                          `gson:"hashDepth"`
	unionBuilderPoolSize            int32                                                          `gson:"unionBuilderPoolSize"`
	replacerPoolSize                int32                                                          `gson:"replacerPoolSize"`
	idHash_                         atomic.Int32
	slab_UnionType                  j2gort.PooledSlab[UnionType]
	slab_FunctionType               j2gort.PooledSlab[FunctionType]
	slab_TemplatizedType            j2gort.PooledSlab[TemplatizedType]
}

func (this_ *JSTypeRegistry) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeRegistry) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeRegistry) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSTypeRegistry", "JSTypeRegistry"
}
func (this_ *JSTypeRegistry) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSTypeRegistry) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeRegistry", "com.google.javascript.rhino.jstype.JSTypeRegistry", func(v_ any) bool {
		return j2gort.InstanceOf[*JSTypeRegistry](v_)
	})
}

var static_JSTypeRegistry_DOT_SPLITTER *j2gort.Ext_com_google_common_base_Splitter

func StaticGet_JSTypeRegistry_DOT_SPLITTER() *j2gort.Ext_com_google_common_base_Splitter {
	clinit_JSTypeRegistry()
	return static_JSTypeRegistry_DOT_SPLITTER
}
func StaticSet_JSTypeRegistry_DOT_SPLITTER(v *j2gort.Ext_com_google_common_base_Splitter) {
	clinit_JSTypeRegistry()
	static_JSTypeRegistry_DOT_SPLITTER = v
}
func StaticGet_JSTypeRegistry_I_OBJECT_ELEMENT_TEMPLATE() string { return "IOBJECT_VALUE" }
func StaticGet_JSTypeRegistry_OBJECT_ELEMENT_TEMPLATE() string   { return "IOBJECT_VALUE" }
func StaticGet_JSTypeRegistry_UNION2_CACHE_MASK() int32          { return int32(16383) }

var clinitDone_JSTypeRegistry uint32
var clinitState_JSTypeRegistry j2gort.ClinitState

func clinit_JSTypeRegistry() {
	if atomic.LoadUint32(&clinitDone_JSTypeRegistry) != 0 {
		return
	}
	clinitSlow_JSTypeRegistry()
}

func clinitSlow_JSTypeRegistry() {
	if !j2gort.EnterClinit(&clinitState_JSTypeRegistry, &clinitDone_JSTypeRegistry) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSTypeRegistry, &clinitDone_JSTypeRegistry)
	static_JSTypeRegistry_DOT_SPLITTER = j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(46))
}

func (this_ *JSTypeRegistry) initFields_JSTypeRegistry() {
	_ = this_
	this_.scopedNameTable = NewJSTypeRegistry_ScopedNameTable()
	this_.closureNamespaces = j2gort.NewOrderedMap[string, JSTypeRegistry_ClosureNamespace]()
	this_.nonLegacyClosureNamespaceRoots = j2gort.NewOrderedSet[string]()
	this_.nameTableGlobalRoot = NewNode_1(StaticGet_Token_ROOT())
	this_.nonNullableTypeNames = NewJSTypeRegistry_ScopedNameSet()

	this_.unionBuilderPool = make([]*UnionType_Builder, int32(16))

	this_.replacerPool = make([]*TemplateTypeReplacer, int32(8))

	this_.nonRefTypesIndexedByProperty = j2gort.NewOrderedMap[string, *j2gort.Ext_java_util_LinkedHashSet]()
	this_.propertiesOfSupertypesInUnions = j2gort.NewExt_java_util_LinkedHashSet()
	this_.droppedPropertiesOfUnions = j2gort.NewOrderedSet[string]()
	this_.eachRefTypeIndexedByProperty = j2gort.NewOrderedMap[string, *j2gort.Ext_java_util_LinkedHashSet](

	// Java constructor: JSTypeRegistry(com.google.javascript.rhino.ErrorReporter reporter)
	)
}

func NewJSTypeRegistry_1(reporter ErrorReporter) *JSTypeRegistry {
	clinit_JSTypeRegistry()
	this_ := &JSTypeRegistry{}
	InitJSTypeRegistry_1(this_, reporter)
	return this_
}

func InitJSTypeRegistry_1(this_ *JSTypeRegistry, reporter ErrorReporter) {
	_ = reporter
	InitJSTypeRegistry_2(this_, reporter, j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()))
}

// Java constructor: JSTypeRegistry(com.google.javascript.rhino.ErrorReporter reporter, java.util.Set<java.lang.String> forwardDeclaredTypes)
func NewJSTypeRegistry_2(reporter ErrorReporter, forwardDeclaredTypes j2gort.Ext_java_util_Set) *JSTypeRegistry {
	clinit_JSTypeRegistry()
	this_ := &JSTypeRegistry{}
	InitJSTypeRegistry_2(this_, reporter, forwardDeclaredTypes)
	return this_
}

func InitJSTypeRegistry_2(this_ *JSTypeRegistry, reporter ErrorReporter, forwardDeclaredTypes j2gort.Ext_java_util_Set) {
	_ = reporter
	_ = forwardDeclaredTypes
	this_.initFields_JSTypeRegistry()
	this_.reporter = reporter
	this_.forwardDeclaredTypes = forwardDeclaredTypes
	this_.emptyTemplateTypeMap = TemplateTypeMap_CreateEmpty(this_)
	this_.resolver = JSTypeResolver_Create(this_)
	this_.nativeTypes = make([]JSType, int32(len(JSTypeNative_Values())))
	func() {
		var closer *JSTypeResolver_Closer = this_.resolver.OpenForDefinition()
		_ = closer
		defer j2gort.CloseResource(closer)
		this_.priv_JSTypeRegistry_initializeBuiltInTypes_9246()
		this_.priv_JSTypeRegistry_initializeRegistry_9247()
	}()
}

// Java method: JSTypeRegistry.getNextTypeId()
func (this_ *JSTypeRegistry) GetNextTypeId() int32 {
	return this_.nextTypeId
}

// Java method: JSTypeRegistry.acquireTemplateTypeReplacer(com.google.javascript.rhino.jstype.TemplateTypeMap bindings, boolean visitProperties, boolean useUnknownForMissingKeys, boolean useUnknownForMissingValues)
func (this_ *JSTypeRegistry) AcquireTemplateTypeReplacer(bindings *TemplateTypeMap, visitProperties bool, useUnknownForMissingKeys bool, useUnknownForMissingValues bool) *TemplateTypeReplacer {
	_ = bindings
	_ = visitProperties
	_ = useUnknownForMissingKeys
	_ = useUnknownForMissingValues
	var r *TemplateTypeReplacer = nil
	_ = r
	var n int32 = this_.replacerPoolSize
	_ = n
	if n > int32(0) {
		var top int32 = (n - int32(1))
		_ = top
		r = this_.replacerPool[top]
		{
			a_1 := this_.replacerPool
			i_2 := top
			a_1[i_2] = nil
		}
		this_.replacerPoolSize = top
	} else {
		r = TemplateTypeReplacer_CreatePooled(this_)
	}
	r.InitForPool(bindings, visitProperties, useUnknownForMissingKeys, useUnknownForMissingValues)
	return r
}

// Java method: JSTypeRegistry.releaseTemplateTypeReplacer(com.google.javascript.rhino.jstype.TemplateTypeReplacer r)
func (this_ *JSTypeRegistry) ReleaseTemplateTypeReplacer(r *TemplateTypeReplacer) {
	_ = r
	r.ResetForPool()
	var n int32 = this_.replacerPoolSize
	_ = n
	if n < int32(len(this_.replacerPool)) {
		{
			a_3 := this_.replacerPool
			i_4 := n
			a_3[i_4] = r
		}
		this_.replacerPoolSize = (n + int32(1))
	}
}

// Java method: JSTypeRegistry.replaceTemplateTypesPartial(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.jstype.TemplateTypeMap bindings)
func (this_ *JSTypeRegistry) ReplaceTemplateTypesPartial(type_ JSType, bindings *TemplateTypeMap) JSType {
	_ = type_
	_ = bindings
	if bindings.IsEmpty() || !(TemplateTypeReplacer_MayBeAffectedByPartialReplacement(type_)) {
		return type_
	}
	var replacer *TemplateTypeReplacer = this_.AcquireTemplateTypeReplacer(bindings, false, false, false)
	_ = replacer
	defer func() {
		this_.ReleaseTemplateTypeReplacer(replacer)
	}()
	{
		return replacer.Replace(type_)
	}
	panic("unreachable")
}

// Java method: JSTypeRegistry.acquireUnionTypeBuilder()
func (this_ *JSTypeRegistry) AcquireUnionTypeBuilder() *UnionType_Builder {
	var n int32 = this_.unionBuilderPoolSize
	_ = n
	if n > int32(0) {
		var top int32 = (n - int32(1))
		_ = top
		var b *UnionType_Builder = this_.unionBuilderPool[top]
		_ = b
		{
			a_5 := this_.unionBuilderPool
			i_6 := top
			a_5[i_6] = nil
		}
		this_.unionBuilderPoolSize = top
		return b
	}
	return UnionType_CreatePooledBuilder(this_)
}

// Java method: JSTypeRegistry.acquireUnionTypeRebuilder(com.google.javascript.rhino.jstype.UnionType rebuildTarget)
func (this_ *JSTypeRegistry) AcquireUnionTypeRebuilder(rebuildTarget *UnionType) *UnionType_Builder {
	_ = rebuildTarget
	var b *UnionType_Builder = this_.AcquireUnionTypeBuilder()
	_ = b
	b.SetRebuildTargetForPool(rebuildTarget)
	return b
}

// Java method: JSTypeRegistry.releaseUnionTypeBuilder(com.google.javascript.rhino.jstype.UnionType.Builder b)
func (this_ *JSTypeRegistry) ReleaseUnionTypeBuilder(b *UnionType_Builder) {
	_ = b
	b.ResetForPool()
	var n int32 = this_.unionBuilderPoolSize
	_ = n
	if n < int32(len(this_.unionBuilderPool)) {
		{
			a_7 := this_.unionBuilderPool
			i_8 := n
			a_7[i_8] = b
		}
		this_.unionBuilderPoolSize = (n + int32(1))
	}
}

// Java method: JSTypeRegistry.getSentinelObjectLiteral()
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_getSentinelObjectLiteral_9221() JSType {
	if this_.sentinelObjectLiteral == nil || this_.sentinelObjectLiteral.IsJavaNil_() {
		this_.sentinelObjectLiteral = Upcast_JSType(this_.CreateAnonymousObjectType(nil))
	}
	return this_.sentinelObjectLiteral
}

// Java method: JSTypeRegistry.getArrayElementKey()
func (this_ *JSTypeRegistry) GetArrayElementKey() *TemplateType {
	return this_.arrayElementTemplateKey
}

// Java method: JSTypeRegistry.getReadonlyArrayElementKey()
func (this_ *JSTypeRegistry) GetReadonlyArrayElementKey() *TemplateType {
	return this_.readonlyArrayElementTemplateKey
}

// Java method: JSTypeRegistry.getReadonlyMapKey()
func (this_ *JSTypeRegistry) GetReadonlyMapKey() *TemplateType {
	return this_.readonlyMapKeyTemplateKey
}

// Java method: JSTypeRegistry.getReadonlyMapValue()
func (this_ *JSTypeRegistry) GetReadonlyMapValue() *TemplateType {
	return this_.readonlyMapValueTemplateKey
}

// Java method: JSTypeRegistry.getObjectElementKey()
func (this_ *JSTypeRegistry) GetObjectElementKey() *TemplateType {
	return this_.iObjectElementTemplateKey
}

// Java method: JSTypeRegistry.getObjectIndexKey()
func (this_ *JSTypeRegistry) GetObjectIndexKey() *TemplateType {
	j2gort.CheckNotNullPtr(this_.iObjectIndexTemplateKey)
	return this_.iObjectIndexTemplateKey
}

// Java method: JSTypeRegistry.getIterableValueTemplate()
func (this_ *JSTypeRegistry) GetIterableValueTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.iterableValueTemplate)
}

// Java method: JSTypeRegistry.getIterableReturnTemplate()
func (this_ *JSTypeRegistry) GetIterableReturnTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.iterableReturnTemplate)
}

// Java method: JSTypeRegistry.getIteratorIterableValueTemplate()
func (this_ *JSTypeRegistry) GetIteratorIterableValueTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.iteratorIterableValueTemplate)
}

// Java method: JSTypeRegistry.getIteratorIterableReturnTemplate()
func (this_ *JSTypeRegistry) GetIteratorIterableReturnTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.iteratorIterableReturnTemplate)
}

// Java method: JSTypeRegistry.getIIterableResultValueTemplate()
func (this_ *JSTypeRegistry) GetIIterableResultValueTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.iiterableResultValueTemplate)
}

// Java method: JSTypeRegistry.getIteratorLikeValueTemplate()
func (this_ *JSTypeRegistry) GetIteratorLikeValueTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.iteratorLikeValueTemplate)
}

// Java method: JSTypeRegistry.getIteratorLikeReturnTemplate()
func (this_ *JSTypeRegistry) GetIteratorLikeReturnTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.iteratorLikeReturnTemplate)
}

// Java method: JSTypeRegistry.getIteratorValueTemplate()
func (this_ *JSTypeRegistry) GetIteratorValueTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.iteratorValueTemplate)
}

// Java method: JSTypeRegistry.getIteratorReturnTemplate()
func (this_ *JSTypeRegistry) GetIteratorReturnTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.iteratorReturnTemplate)
}

// Java method: JSTypeRegistry.getGeneratorValueTemplate()
func (this_ *JSTypeRegistry) GetGeneratorValueTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.generatorValueTemplate)
}

// Java method: JSTypeRegistry.getGeneratorReturnTemplate()
func (this_ *JSTypeRegistry) GetGeneratorReturnTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.generatorReturnTemplate)
}

// Java method: JSTypeRegistry.getAsyncIterableValueTemplate()
func (this_ *JSTypeRegistry) GetAsyncIterableValueTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.asyncIterableValueTemplate)
}

// Java method: JSTypeRegistry.getAsyncIterableReturnTemplate()
func (this_ *JSTypeRegistry) GetAsyncIterableReturnTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.asyncIterableReturnTemplate)
}

// Java method: JSTypeRegistry.getAsyncIteratorValueTemplate()
func (this_ *JSTypeRegistry) GetAsyncIteratorValueTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.asyncIteratorValueTemplate)
}

// Java method: JSTypeRegistry.getAsyncIteratorReturnTemplate()
func (this_ *JSTypeRegistry) GetAsyncIteratorReturnTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.asyncIteratorReturnTemplate)
}

// Java method: JSTypeRegistry.getIThenableTemplate()
func (this_ *JSTypeRegistry) GetIThenableTemplate() *TemplateType {
	return j2gort.CheckNotNullPtr(this_.iThenableTemplateKey)
}

// Java method: JSTypeRegistry.maybeGetTemplateTypesOfBuiltin(com.google.javascript.rhino.StaticScope scope, java.lang.String fnName)
func (this_ *JSTypeRegistry) MaybeGetTemplateTypesOfBuiltin(scope StaticScope, fnName string) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = scope
	_ = fnName
	var type_ JSType = this_.GetType_2(scope, fnName)
	_ = type_
	var objType ObjectType
	_ = objType
	if type_ == nil || type_.IsJavaNil_() {
		objType = nil
	} else {
		objType = type_.GetJSTypeBase().ToObjectType()
	}
	if (objType != nil && !objType.IsJavaNil_()) && objType.GetObjectTypeBase().IsNativeObjectType() {
		return objType.GetObjectTypeBase().GetTypeParameters()
	}
	return nil
}

// Java method: JSTypeRegistry.getErrorReporter()
func (this_ *JSTypeRegistry) GetErrorReporter() ErrorReporter {
	return this_.reporter
}

// Java method: JSTypeRegistry.initializeBuiltInTypes()
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_initializeBuiltInTypes_9246() {
	var booleanType *BooleanType = NewBooleanType(this_)
	_ = booleanType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_BOOLEAN_TYPE(), booleanType)
	var nullType *NullType = NewNullType(this_)
	_ = nullType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NULL_TYPE(), nullType)
	var bigIntType *BigIntType = NewBigIntType(this_)
	_ = bigIntType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_BIGINT_TYPE(), bigIntType)
	var numberType *NumberType = NewNumberType(this_)
	_ = numberType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NUMBER_TYPE(), numberType)
	var stringType *StringType = NewStringType(this_)
	_ = stringType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_STRING_TYPE(), stringType)
	var symbolType SymbolType = NewSymbolType_1(this_)
	_ = symbolType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_SYMBOL_TYPE(), Upcast_JSType(symbolType))
	var unknownType *UnknownType = NewUnknownType(this_, false)
	_ = unknownType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_UNKNOWN_TYPE(), unknownType)
	var checkedUnknownType *UnknownType = NewUnknownType(this_, true)
	_ = checkedUnknownType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE(), checkedUnknownType)
	var voidType *VoidType = NewVoidType(this_)
	_ = voidType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_VOID_TYPE(), voidType)
	var allType *AllType = NewAllType(this_)
	_ = allType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ALL_TYPE(), allType)
	this_.iObjectIndexTemplateKey = NewTemplateType_2(this_, "IOBJECT_KEY")
	this_.iObjectElementTemplateKey = NewTemplateType_2(this_, "IOBJECT_VALUE")
	var iArrayLikeTemplate *TemplateType = NewTemplateType_2(this_, "VALUE2")
	_ = iArrayLikeTemplate
	this_.arrayElementTemplateKey = NewTemplateType_2(this_, "T")
	this_.readonlyArrayElementTemplateKey = NewTemplateType_2(this_, "T")
	this_.readonlyMapKeyTemplateKey = NewTemplateType_2(this_, "K")
	this_.readonlyMapValueTemplateKey = NewTemplateType_2(this_, "V")
	this_.mapKeyTemplateKey = NewTemplateType_2(this_, "K")
	this_.mapValueTemplateKey = NewTemplateType_2(this_, "V")
	this_.iThenableTemplateKey = NewTemplateType_2(this_, "TYPE")
	this_.promiseTemplateKey = NewTemplateType_2(this_, "TYPE")
	var functionPrototype PrototypeObjectType = j2gort.AsIface[PrototypeObjectType](j2gort.AsIface[PrototypeObjectType_Builder](PrototypeObjectType_S_Builder(this_).SetName("Function.prototype")).GetPrototypeObjectType_BuilderBase().SetNative(true).Build())
	_ = functionPrototype
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_FUNCTION_PROTOTYPE(), Upcast_JSType(functionPrototype))
	var functionInstancePrototype PrototypeObjectType = j2gort.AsIface[PrototypeObjectType](j2gort.AsIface[PrototypeObjectType_Builder](PrototypeObjectType_S_Builder(this_).SetName("?.prototype")).GetPrototypeObjectType_BuilderBase().SetNative(true).Build())
	_ = functionInstancePrototype
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE(), Upcast_JSType(functionInstancePrototype))
	var fb_ FunctionType_Builder
	InitFunctionType_Builder(&fb_, this_)
	(&fb_).WithName("Function")
	(&fb_).ForConstructor()
	(&fb_).ForNativeType()
	(&fb_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateParametersWithVarArgs([]JSType{unknownType})))
	(&fb_).WithTypeOfThis(unknownType)
	(&fb_).WithReturnType_1(unknownType)
	var functionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_).Build())
	_ = functionType
	functionType.SetPrototype(Upcast_ObjectType(functionInstancePrototype), Node(0))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_FUNCTION_TYPE(), ToIface_JSType(functionType))
	var fb_1_ FunctionType_Builder
	InitFunctionType_Builder(&fb_1_, this_)
	(&fb_1_).WithName("Function")
	(&fb_1_).ForConstructor()
	(&fb_1_).ForNativeType()
	(&fb_1_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateParametersWithVarArgs([]JSType{allType})))
	(&fb_1_).WithTypeOfThis(ToIface_JSType(functionType))
	var functionFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_1_).Build())
	_ = functionFunctionType
	functionFunctionType.SetPrototype(Upcast_ObjectType(functionPrototype), Node(0))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_FUNCTION_FUNCTION_TYPE(), ToIface_JSType(functionFunctionType))
	var objectPrototype ObjectType = Upcast_ObjectType(j2gort.AsIface[PrototypeObjectType](j2gort.AsIface[PrototypeObjectType_Builder](PrototypeObjectType_S_Builder(this_).SetName("Object.prototype")).GetPrototypeObjectType_BuilderBase().SetNative(true).GetPrototypeObjectType_BuilderBase().SetImplicitPrototype(nil).Build()))
	_ = objectPrototype
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_OBJECT_PROTOTYPE(), Upcast_JSType(objectPrototype))
	var fb_2_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_2_, "Object")
	(&fb_2_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateOptionalParameters([]JSType{allType})))
	(&fb_2_).WithReturnsOwnInstanceType()
	(&fb_2_).WithTemplateKeys_1_TemplateType([]*TemplateType{this_.iObjectIndexTemplateKey, this_.iObjectElementTemplateKey})
	var objectFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_2_).Build())
	_ = objectFunctionType
	objectFunctionType.SetPrototype(objectPrototype, Node(0))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_OBJECT_FUNCTION_TYPE(), ToIface_JSType(objectFunctionType))
	var objectType ObjectType = objectFunctionType.GetInstanceType()
	_ = objectType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_OBJECT_TYPE(), Upcast_JSType(objectType))
	functionPrototype.ClearCachedValues()
	functionPrototype.GetPrototypeObjectTypeBase().SetImplicitPrototype(objectType)
	functionInstancePrototype.ClearCachedValues()
	functionInstancePrototype.GetPrototypeObjectTypeBase().SetImplicitPrototype(objectType)
	var iObjectFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("IObject", []*TemplateType{this_.iObjectIndexTemplateKey, this_.iObjectElementTemplateKey})
	_ = iObjectFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_I_OBJECT_FUNCTION_TYPE(), ToIface_JSType(iObjectFunctionType))
	var iObjectType ObjectType = iObjectFunctionType.GetInstanceType()
	_ = iObjectType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_I_OBJECT_TYPE(), Upcast_JSType(iObjectType))
	var noType *FunctionType = FunctionType_CreateNoType(this_)
	_ = noType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NO_TYPE(), ToIface_JSType(noType))
	var noObjectType *FunctionType = FunctionType_CreateNoObjectType(this_)
	_ = noObjectType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NO_OBJECT_TYPE(), ToIface_JSType(noObjectType))
	this_.iterableValueTemplate = NewTemplateType_2(this_, "T")
	this_.iterableReturnTemplate = NewTemplateType_2(this_, "TReturn")
	var iterableNextTemplate *TemplateType = NewTemplateType_2(this_, "TNext")
	_ = iterableNextTemplate
	var iterableFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("Iterable", []*TemplateType{this_.iterableValueTemplate, this_.iterableReturnTemplate, iterableNextTemplate})
	_ = iterableFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ITERABLE_FUNCTION_TYPE(), ToIface_JSType(iterableFunctionType))
	var iterableType ObjectType = iterableFunctionType.GetInstanceType()
	_ = iterableType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ITERABLE_TYPE(), Upcast_JSType(iterableType))
	this_.iteratorLikeValueTemplate = NewTemplateType_2(this_, "T")
	this_.iteratorLikeReturnTemplate = NewTemplateType_2(this_, "TReturn")
	var iteratorLikeNextTemplate *TemplateType = NewTemplateType_2(this_, "TNext")
	_ = iteratorLikeNextTemplate
	var iteratorLikeFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("IteratorLike", []*TemplateType{this_.iteratorLikeValueTemplate, this_.iteratorLikeReturnTemplate, iteratorLikeNextTemplate})
	_ = iteratorLikeFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ITERATOR_LIKE_FUNCTION_TYPE(), ToIface_JSType(iteratorLikeFunctionType))
	var iteratorLikeType ObjectType = iteratorLikeFunctionType.GetInstanceType()
	_ = iteratorLikeType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ITERATOR_LIKE_TYPE(), Upcast_JSType(iteratorLikeType))
	this_.iteratorValueTemplate = NewTemplateType_2(this_, "T")
	this_.iteratorReturnTemplate = NewTemplateType_2(this_, "TReturn")
	var iteratorNextTemplate *TemplateType = NewTemplateType_2(this_, "TNext")
	_ = iteratorNextTemplate
	this_.iteratorIterableValueTemplate = NewTemplateType_2(this_, "T")
	this_.iteratorIterableReturnTemplate = NewTemplateType_2(this_, "TReturn")
	var iteratorIterableNextTemplate *TemplateType = NewTemplateType_2(this_, "TNext")
	_ = iteratorIterableNextTemplate
	var iteratorIterableFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("IteratorIterable", []*TemplateType{this_.iteratorIterableValueTemplate, this_.iteratorIterableReturnTemplate, iteratorIterableNextTemplate})
	_ = iteratorIterableFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ITERATOR_ITERABLE_FUNCTION_TYPE(), ToIface_JSType(iteratorIterableFunctionType))
	iteratorIterableFunctionType.SetExtendedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(this_.CreateTemplatizedType_2_ObjectType_JSType(iterableType, []JSType{ToIface_JSType(this_.iteratorIterableValueTemplate), ToIface_JSType(this_.iteratorIterableReturnTemplate), iteratorIterableNextTemplate})), j2gort.PtrToAny(this_.CreateTemplatizedType_2_ObjectType_JSType(iteratorLikeType, []JSType{ToIface_JSType(this_.iteratorIterableValueTemplate), ToIface_JSType(this_.iteratorIterableReturnTemplate), iteratorIterableNextTemplate})))))
	var iteratorIterableType ObjectType = iteratorIterableFunctionType.GetInstanceType()
	_ = iteratorIterableType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ITERATOR_ITERABLE_TYPE(), Upcast_JSType(iteratorIterableType))
	var fb_3_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_3_, "Iterator")
	(&fb_3_).WithParameters_0()
	(&fb_3_).WithTemplateKeys_1_TemplateType([]*TemplateType{this_.iteratorValueTemplate, this_.iteratorReturnTemplate, iteratorNextTemplate})
	var iteratorFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_3_).Build())
	_ = iteratorFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ITERATOR_FUNCTION_TYPE(), ToIface_JSType(iteratorFunctionType))
	iteratorFunctionType.GetPrototype()
	iteratorFunctionType.SetImplementedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(this_.CreateTemplatizedType_2_ObjectType_JSType(iteratorIterableType, []JSType{ToIface_JSType(this_.iteratorValueTemplate), ToIface_JSType(this_.iteratorReturnTemplate), iteratorNextTemplate}))))
	var iteratorType ObjectType = iteratorFunctionType.GetInstanceType()
	_ = iteratorType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ITERATOR_TYPE(), Upcast_JSType(iteratorType))
	this_.iiterableResultValueTemplate = NewTemplateType_2(this_, "TYield")
	var iiterableResultFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeRecord_9337("IIterableResult", []*TemplateType{this_.iiterableResultValueTemplate})
	_ = iiterableResultFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_I_ITERABLE_RESULT_FUNCTION_TYPE(), ToIface_JSType(iiterableResultFunctionType))
	var iiterableResultType ObjectType = iiterableResultFunctionType.GetInstanceType()
	_ = iiterableResultType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_I_ITERABLE_RESULT_TYPE(), Upcast_JSType(iiterableResultType))
	var iArrayLikeFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeRecord_9337("IArrayLike", []*TemplateType{iArrayLikeTemplate})
	_ = iArrayLikeFunctionType
	iArrayLikeFunctionType.SetExtendedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(this_.CreateTemplatizedType_2_ObjectType_JSType(iObjectType, []JSType{numberType, iArrayLikeTemplate}))))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_I_ARRAY_LIKE_FUNCTION_TYPE(), ToIface_JSType(iArrayLikeFunctionType))
	var iArrayLikeType ObjectType = iArrayLikeFunctionType.GetInstanceType()
	_ = iArrayLikeType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_I_ARRAY_LIKE_TYPE(), Upcast_JSType(iArrayLikeType))
	var readonlyArrayFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeRecord_9337("ReadonlyArray", []*TemplateType{this_.readonlyArrayElementTemplateKey})
	_ = readonlyArrayFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_READONLY_ARRAY_FUNCTION_TYPE(), ToIface_JSType(readonlyArrayFunctionType))
	readonlyArrayFunctionType.SetExtendedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(this_.CreateTemplatizedType_2_ObjectType_JSType(iArrayLikeType, []JSType{ToIface_JSType(this_.readonlyArrayElementTemplateKey)})), j2gort.PtrToAny(this_.CreateTemplatizedType_2_ObjectType_JSType(iterableType, []JSType{ToIface_JSType(this_.readonlyArrayElementTemplateKey)})))))
	var readonlyArrayType ObjectType = readonlyArrayFunctionType.GetInstanceType()
	_ = readonlyArrayType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_READONLY_ARRAY_TYPE(), Upcast_JSType(readonlyArrayType))
	var fb_4_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_4_, "Array")
	(&fb_4_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateParametersWithVarArgs([]JSType{allType})))
	(&fb_4_).WithReturnsOwnInstanceType()
	(&fb_4_).WithTemplateKeys_1_TemplateType([]*TemplateType{this_.arrayElementTemplateKey})
	var arrayFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_4_).Build())
	_ = arrayFunctionType
	arrayFunctionType.GetPrototype()
	arrayFunctionType.SetImplementedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(this_.CreateTemplatizedType_2_ObjectType_JSType(readonlyArrayType, []JSType{ToIface_JSType(this_.arrayElementTemplateKey)}))))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ARRAY_FUNCTION_TYPE(), ToIface_JSType(arrayFunctionType))
	var arrayType ObjectType = arrayFunctionType.GetInstanceType()
	_ = arrayType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ARRAY_TYPE(), Upcast_JSType(arrayType))
	var arrayOfString ObjectType = ToIface_ObjectType(this_.CreateTemplatizedType_2_ObjectType_JSType(arrayType, []JSType{stringType}))
	_ = arrayOfString
	var fb_5_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_5_, "ITemplateArray")
	(&fb_5_).WithParameters_0()
	var iTemplateArrayFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_5_).Build())
	_ = iTemplateArrayFunctionType
	var iTemplateArrayPrototype ObjectType = Upcast_ObjectType(j2gort.AsIface[PrototypeObjectType](j2gort.AsIface[PrototypeObjectType_Builder](PrototypeObjectType_S_Builder(this_).SetName("ITemplateArray.prototype")).GetPrototypeObjectType_BuilderBase().SetNative(true).GetPrototypeObjectType_BuilderBase().SetImplicitPrototype(arrayOfString).Build()))
	_ = iTemplateArrayPrototype
	iTemplateArrayFunctionType.SetPrototype(iTemplateArrayPrototype, Node(0))
	iTemplateArrayFunctionType.GetInstanceType().GetObjectTypeBase().MergeSupertypeTemplateTypes(arrayOfString)
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_I_TEMPLATE_ARRAY_TYPE(), Upcast_JSType(iTemplateArrayFunctionType.GetInstanceType()))
	this_.generatorValueTemplate = NewTemplateType_2(this_, "T")
	this_.generatorReturnTemplate = NewTemplateType_2(this_, "TReturn")
	var generatorNextTemplate *TemplateType = NewTemplateType_2(this_, "TNext")
	_ = generatorNextTemplate
	var generatorFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("Generator", []*TemplateType{this_.generatorValueTemplate, this_.generatorReturnTemplate, generatorNextTemplate})
	_ = generatorFunctionType
	generatorFunctionType.SetExtendedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(this_.CreateTemplatizedType_2_ObjectType_JSType(iteratorIterableType, []JSType{ToIface_JSType(this_.generatorValueTemplate), ToIface_JSType(this_.generatorReturnTemplate), generatorNextTemplate}))))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_GENERATOR_FUNCTION_TYPE(), ToIface_JSType(generatorFunctionType))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_GENERATOR_TYPE(), Upcast_JSType(generatorFunctionType.GetInstanceType()))
	this_.asyncIteratorValueTemplate = NewTemplateType_2(this_, "T")
	this_.asyncIteratorReturnTemplate = NewTemplateType_2(this_, "TReturn")
	var asyncIteratorNextTemplate *TemplateType = NewTemplateType_2(this_, "TNext")
	_ = asyncIteratorNextTemplate
	var asyncIteratorFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("AsyncIterator", []*TemplateType{this_.asyncIteratorValueTemplate, this_.asyncIteratorReturnTemplate, asyncIteratorNextTemplate})
	_ = asyncIteratorFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ASYNC_ITERATOR_FUNCTION_TYPE(), ToIface_JSType(asyncIteratorFunctionType))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ASYNC_ITERATOR_TYPE(), Upcast_JSType(asyncIteratorFunctionType.GetInstanceType()))
	this_.asyncIterableValueTemplate = NewTemplateType_2(this_, "T")
	this_.asyncIterableReturnTemplate = NewTemplateType_2(this_, "TReturn")
	var asyncIterableNextTemplate *TemplateType = NewTemplateType_2(this_, "TNext")
	_ = asyncIterableNextTemplate
	var asyncIterableFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("AsyncIterable", []*TemplateType{this_.asyncIterableValueTemplate, this_.asyncIterableReturnTemplate, asyncIterableNextTemplate})
	_ = asyncIterableFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ASYNC_ITERABLE_FUNCTION_TYPE(), ToIface_JSType(asyncIterableFunctionType))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ASYNC_ITERABLE_TYPE(), Upcast_JSType(asyncIterableFunctionType.GetInstanceType()))
	var asyncIteratorIterableTemplate *TemplateType = NewTemplateType_2(this_, "T")
	_ = asyncIteratorIterableTemplate
	var asyncIteratorIterableReturnTemplate *TemplateType = NewTemplateType_2(this_, "TReturn")
	_ = asyncIteratorIterableReturnTemplate
	var asyncIteratorIterableNextTemplate *TemplateType = NewTemplateType_2(this_, "TNext")
	_ = asyncIteratorIterableNextTemplate
	var asyncIteratorIterableFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("AsyncIteratorIterable", []*TemplateType{asyncIteratorIterableTemplate, asyncIteratorIterableReturnTemplate, asyncIteratorIterableNextTemplate})
	_ = asyncIteratorIterableFunctionType
	asyncIteratorIterableFunctionType.SetExtendedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(this_.CreateTemplatizedType_2_ObjectType_JSType(asyncIteratorFunctionType.GetInstanceType(), []JSType{asyncIteratorIterableTemplate, asyncIteratorIterableReturnTemplate, asyncIteratorIterableNextTemplate})), j2gort.PtrToAny(this_.CreateTemplatizedType_2_ObjectType_JSType(asyncIterableFunctionType.GetInstanceType(), []JSType{asyncIteratorIterableTemplate, asyncIteratorIterableReturnTemplate, asyncIteratorIterableNextTemplate})))))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ASYNC_ITERATOR_ITERABLE_FUNCTION_TYPE(), ToIface_JSType(asyncIteratorIterableFunctionType))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE(), Upcast_JSType(asyncIteratorIterableFunctionType.GetInstanceType()))
	var asyncGeneratorValueTemplate *TemplateType = NewTemplateType_2(this_, "T")
	_ = asyncGeneratorValueTemplate
	var asyncGeneratorReturnTemplate *TemplateType = NewTemplateType_2(this_, "TReturn")
	_ = asyncGeneratorReturnTemplate
	var asyncGeneratorNextTemplate *TemplateType = NewTemplateType_2(this_, "TNext")
	_ = asyncGeneratorNextTemplate
	var asyncGeneratorFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("AsyncGenerator", []*TemplateType{asyncGeneratorValueTemplate, asyncGeneratorReturnTemplate, asyncGeneratorNextTemplate})
	_ = asyncGeneratorFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ASYNC_GENERATOR_FUNCTION_TYPE(), ToIface_JSType(asyncGeneratorFunctionType))
	asyncGeneratorFunctionType.SetExtendedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(this_.CreateTemplatizedType_2_ObjectType_JSType(asyncIteratorIterableFunctionType.GetInstanceType(), []JSType{asyncGeneratorValueTemplate, asyncGeneratorReturnTemplate, asyncGeneratorNextTemplate}))))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ASYNC_GENERATOR_TYPE(), Upcast_JSType(asyncGeneratorFunctionType.GetInstanceType()))
	var ithenableFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("IThenable", []*TemplateType{this_.iThenableTemplateKey})
	_ = ithenableFunctionType
	ithenableFunctionType.SetStruct()
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_I_THENABLE_FUNCTION_TYPE(), ToIface_JSType(ithenableFunctionType))
	var ithenableType ObjectType = ithenableFunctionType.GetInstanceType()
	_ = ithenableType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_I_THENABLE_TYPE(), Upcast_JSType(ithenableType))
	var thenableType JSType = this_.CreateRecordType(j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_2("then", unknownType)))
	_ = thenableType
	this_.IdentifyNonNullableName(nil, "Thenable")
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_THENABLE_TYPE(), thenableType)
	var promiseParameterType JSType = ToIface_JSType(this_.CreateFunctionType_2_JSType_JSType(unknownType, []JSType{ToIface_JSType(this_.CreateFunctionType_2_JSType_Ext_java_util_List(unknownType, j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateOptionalParameters([]JSType{this_.CreateUnionType_1_JSType([]JSType{ToIface_JSType(this_.promiseTemplateKey), ToIface_JSType(this_.CreateTemplatizedType_2_ObjectType_JSType(ithenableType, []JSType{ToIface_JSType(this_.promiseTemplateKey)})), thenableType, nullType})})))), ToIface_JSType(this_.CreateFunctionType_2_JSType_Ext_java_util_List(unknownType, j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateOptionalParameters([]JSType{allType}))))}))
	_ = promiseParameterType
	var fb_6_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_6_, "Promise")
	(&fb_6_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateParameters_1([]JSType{promiseParameterType})))
	(&fb_6_).WithTemplateKeys_1_TemplateType([]*TemplateType{this_.promiseTemplateKey})
	var promiseFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_6_).Build())
	_ = promiseFunctionType
	promiseFunctionType.SetImplementedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(this_.CreateTemplatizedType_2_ObjectType_JSType(ithenableType, []JSType{ToIface_JSType(this_.promiseTemplateKey)}))))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_PROMISE_FUNCTION_TYPE(), ToIface_JSType(promiseFunctionType))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_PROMISE_TYPE(), Upcast_JSType(promiseFunctionType.GetInstanceType()))
	var fb_7_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_7_, "Arguments")
	(&fb_7_).WithParameters_0()
	var argumentsFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_7_).Build())
	_ = argumentsFunctionType
	argumentsFunctionType.SetImplementedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(this_.CreateTemplatizedType_2_ObjectType_JSType(iArrayLikeType, []JSType{unknownType})), j2gort.PtrToAny(this_.CreateTemplatizedType_2_ObjectType_JSType(iterableType, []JSType{unknownType})))))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ARGUMENTS_FUNCTION_TYPE(), ToIface_JSType(argumentsFunctionType))
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_ARGUMENTS_TYPE(), Upcast_JSType(argumentsFunctionType.GetInstanceType()))
	var bigintNumberString JSType = this_.CreateUnionType_1_JSType([]JSType{bigIntType, numberType, stringType})
	_ = bigintNumberString
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_BIGINT_NUMBER_STRING(), bigintNumberString)
	var fb_8_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_8_, "BigInt")
	(&fb_8_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateParameters_1([]JSType{bigintNumberString})))
	(&fb_8_).WithReturnType_1(bigIntType)
	var bigIntObjectFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_8_).Build())
	_ = bigIntObjectFunctionType
	bigIntObjectFunctionType.GetPrototype()
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_BIGINT_OBJECT_FUNCTION_TYPE(), ToIface_JSType(bigIntObjectFunctionType))
	var bigIntObjectType ObjectType = bigIntObjectFunctionType.GetInstanceType()
	_ = bigIntObjectType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_BIGINT_OBJECT_TYPE(), Upcast_JSType(bigIntObjectType))
	var fb_9_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_9_, "Boolean")
	(&fb_9_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateOptionalParameters([]JSType{allType})))
	(&fb_9_).WithReturnType_1(booleanType)
	var booleanObjectFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_9_).Build())
	_ = booleanObjectFunctionType
	booleanObjectFunctionType.GetPrototype()
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE(), ToIface_JSType(booleanObjectFunctionType))
	var booleanObjectType ObjectType = booleanObjectFunctionType.GetInstanceType()
	_ = booleanObjectType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_BOOLEAN_OBJECT_TYPE(), Upcast_JSType(booleanObjectType))
	var fb_10_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_10_, "Date")
	(&fb_10_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateOptionalParameters([]JSType{unknownType, unknownType, unknownType, unknownType, unknownType, unknownType, unknownType})))
	(&fb_10_).WithReturnType_1(stringType)
	var dateFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_10_).Build())
	_ = dateFunctionType
	dateFunctionType.GetPrototype()
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_DATE_FUNCTION_TYPE(), ToIface_JSType(dateFunctionType))
	var dateType ObjectType = dateFunctionType.GetInstanceType()
	_ = dateType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_DATE_TYPE(), Upcast_JSType(dateType))
	var fb_11_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_11_, "Number")
	(&fb_11_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateOptionalParameters([]JSType{allType})))
	(&fb_11_).WithReturnType_1(numberType)
	var numberObjectFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_11_).Build())
	_ = numberObjectFunctionType
	numberObjectFunctionType.GetPrototype()
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NUMBER_OBJECT_FUNCTION_TYPE(), ToIface_JSType(numberObjectFunctionType))
	var numberObjectType ObjectType = numberObjectFunctionType.GetInstanceType()
	_ = numberObjectType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NUMBER_OBJECT_TYPE(), Upcast_JSType(numberObjectType))
	var fb_12_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_12_, "RegExp")
	(&fb_12_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateOptionalParameters([]JSType{allType, allType})))
	(&fb_12_).WithReturnsOwnInstanceType()
	var regexpFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_12_).Build())
	_ = regexpFunctionType
	regexpFunctionType.GetPrototype()
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_REGEXP_FUNCTION_TYPE(), ToIface_JSType(regexpFunctionType))
	var regexpType ObjectType = regexpFunctionType.GetInstanceType()
	_ = regexpType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_REGEXP_TYPE(), Upcast_JSType(regexpType))
	var fb_13_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_13_, "String")
	(&fb_13_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateOptionalParameters([]JSType{allType})))
	(&fb_13_).WithReturnType_1(stringType)
	var stringObjectFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_13_).Build())
	_ = stringObjectFunctionType
	stringObjectFunctionType.GetPrototype()
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_STRING_OBJECT_FUNCTION_TYPE(), ToIface_JSType(stringObjectFunctionType))
	var stringObjectType ObjectType = stringObjectFunctionType.GetInstanceType()
	_ = stringObjectType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_STRING_OBJECT_TYPE(), Upcast_JSType(stringObjectType))
	var fb_14_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_14_, "Symbol")
	(&fb_14_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateOptionalParameters([]JSType{allType})))
	(&fb_14_).WithReturnType_1(Upcast_JSType(symbolType))
	var symbolObjectFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_14_).Build())
	_ = symbolObjectFunctionType
	symbolObjectFunctionType.GetPrototype()
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_SYMBOL_OBJECT_FUNCTION_TYPE(), ToIface_JSType(symbolObjectFunctionType))
	var symbolObjectType ObjectType = symbolObjectFunctionType.GetInstanceType()
	_ = symbolObjectType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_SYMBOL_OBJECT_TYPE(), Upcast_JSType(symbolObjectType))
	var gbigintFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("gbigint", []*TemplateType{})
	_ = gbigintFunctionType
	gbigintFunctionType.GetPrototype()
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_GBIGINT_TYPE(), Upcast_JSType(gbigintFunctionType.GetInstanceType()))
	var nullVoid JSType = this_.CreateUnionType_2(nullType, voidType)
	_ = nullVoid
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NULL_VOID(), nullVoid)
	var numberStringBoolean JSType = this_.CreateUnionType_1_JSType([]JSType{numberType, stringType, booleanType})
	_ = numberStringBoolean
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NUMBER_STRING_BOOLEAN(), numberStringBoolean)
	var valueTypes JSType = this_.CreateUnionType_1_JSType([]JSType{numberType, stringType, booleanType, Upcast_JSType(symbolType)})
	_ = valueTypes
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_VALUE_TYPES(), valueTypes)
	var numberSymbol JSType = this_.CreateUnionType_2(numberType, Upcast_JSType(symbolType))
	_ = numberSymbol
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NUMBER_SYMBOL(), numberSymbol)
	var stringSymbol JSType = this_.CreateUnionType_2(stringType, Upcast_JSType(symbolType))
	_ = stringSymbol
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_STRING_SYMBOL(), stringSymbol)
	var numberString JSType = this_.CreateUnionType_2(numberType, stringType)
	_ = numberString
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NUMBER_STRING(), numberString)
	var bigintNumber JSType = this_.CreateUnionType_2(bigIntType, numberType)
	_ = bigintNumber
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_BIGINT_NUMBER(), bigintNumber)
	var bigintNumberObject JSType = this_.CreateUnionType_2(Upcast_JSType(bigIntObjectType), Upcast_JSType(numberObjectType))
	_ = bigintNumberObject
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_BIGINT_NUMBER_OBJECT(), bigintNumberObject)
	var bigintNumberStringObject JSType = this_.CreateUnionType_1_JSType([]JSType{Upcast_JSType(bigIntObjectType), Upcast_JSType(numberObjectType), Upcast_JSType(stringObjectType)})
	_ = bigintNumberStringObject
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_BIGINT_NUMBER_STRING_OBJECT(), bigintNumberStringObject)
	var numberStringSymbol JSType = this_.CreateUnionType_1_JSType([]JSType{numberType, stringType, Upcast_JSType(symbolType)})
	_ = numberStringSymbol
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NUMBER_STRING_SYMBOL(), numberStringSymbol)
	var numberAdditionSupertype JSType = this_.CreateUnionType_1_JSTypeNative([]JSTypeNative{StaticGet_JSTypeNative_VOID_TYPE(), StaticGet_JSTypeNative_NULL_TYPE(), StaticGet_JSTypeNative_NUMBER_TYPE(), StaticGet_JSTypeNative_NUMBER_OBJECT_TYPE(), StaticGet_JSTypeNative_BOOLEAN_TYPE(), StaticGet_JSTypeNative_BOOLEAN_OBJECT_TYPE()})
	_ = numberAdditionSupertype
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_NUMBER_ADDITION_SUPERTYPE(), numberAdditionSupertype)
	var leastFunctionType *FunctionType = this_.priv_JSTypeRegistry_createNativeFunctionTypeWithVarArgs_9296(ToIface_JSType(noType), []JSType{allType})
	_ = leastFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_LEAST_FUNCTION_TYPE(), ToIface_JSType(leastFunctionType))
	var fb_15_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_15_, "global this")
	(&fb_15_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateParameters_1([]JSType{allType})))
	(&fb_15_).WithReturnType_1(numberType)
	var globalThisCtor *FunctionType = j2gort.CastPtr[FunctionType]((&fb_15_).Build())
	_ = globalThisCtor
	var globalThis ObjectType = globalThisCtor.GetInstanceType()
	_ = globalThis
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_GLOBAL_THIS(), Upcast_JSType(globalThis))
	var greatestFunctionType *FunctionType = this_.priv_JSTypeRegistry_createNativeFunctionTypeWithVarArgs_9296(allType, []JSType{ToIface_JSType(noType)})
	_ = greatestFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_GREATEST_FUNCTION_TYPE(), ToIface_JSType(greatestFunctionType))
	this_.RegisterPropertyOnType("prototype", ToIface_JSType(objectFunctionType))
	var readonlyMapFunctionType *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336("ReadonlyMap", []*TemplateType{this_.readonlyMapKeyTemplateKey, this_.readonlyMapValueTemplateKey})
	_ = readonlyMapFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_READONLY_MAP_FUNCTION_TYPE(), ToIface_JSType(readonlyMapFunctionType))
	readonlyMapFunctionType.SetExtendedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(this_.CreateTemplatizedType_2_ObjectType_JSType(iterableType, []JSType{ToIface_JSType(this_.CreateTemplatizedType_2_ObjectType_JSType(arrayType, []JSType{this_.CreateUnionType_2(ToIface_JSType(this_.readonlyMapKeyTemplateKey), ToIface_JSType(this_.readonlyMapValueTemplateKey))}))}))))
	var readonlyMapType ObjectType = readonlyMapFunctionType.GetInstanceType()
	_ = readonlyMapType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_READONLY_MAP_TYPE(), Upcast_JSType(readonlyMapType))
	var fb_16_ FunctionType_Builder
	this_.priv_JSTypeRegistry_nativeConstructorBuilder_9335(&fb_16_, "Map")
	(&fb_16_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateOptionalParameters([]JSType{this_.CreateUnionType_1_JSType([]JSType{nullType, ToIface_JSType(this_.CreateTemplatizedType_2_ObjectType_JSType(iterableType, []JSType{ToIface_JSType(this_.CreateTemplatizedType_2_ObjectType_JSType(arrayType, []JSType{this_.CreateUnionType_2(ToIface_JSType(this_.mapKeyTemplateKey), ToIface_JSType(this_.mapValueTemplateKey))}))})), ToIface_JSType(this_.CreateTemplatizedType_2_ObjectType_JSType(arrayType, []JSType{ToIface_JSType(this_.CreateTemplatizedType_2_ObjectType_JSType(arrayType, []JSType{this_.CreateUnionType_2(ToIface_JSType(this_.mapKeyTemplateKey), ToIface_JSType(this_.mapValueTemplateKey))}))}))})})))
	(&fb_16_).WithTemplateKeys_1_TemplateType([]*TemplateType{this_.mapKeyTemplateKey, this_.mapValueTemplateKey})
	var mapFunctionType *FunctionType = j2gort.CastPtr[FunctionType]((&fb_16_).Build())
	_ = mapFunctionType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_MAP_FUNCTION_TYPE(), ToIface_JSType(mapFunctionType))
	mapFunctionType.GetPrototype()
	mapFunctionType.SetImplementedInterfaces(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(this_.CreateTemplatizedType_2_ObjectType_JSType(readonlyMapType, []JSType{ToIface_JSType(this_.mapKeyTemplateKey), ToIface_JSType(this_.mapValueTemplateKey)}))))
	var mapType ObjectType = mapFunctionType.GetInstanceType()
	_ = mapType
	this_.priv_JSTypeRegistry_registerNativeType_9263(StaticGet_JSTypeNative_MAP_TYPE(), Upcast_JSType(mapType))
}

// Java method: JSTypeRegistry.initializeRegistry()
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_initializeRegistry_9247() {
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_ARGUMENTS_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_ARRAY_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_READONLY_ARRAY_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_MAP_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_READONLY_MAP_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_ASYNC_ITERABLE_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_ASYNC_ITERATOR_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_ASYNC_GENERATOR_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_BIGINT_OBJECT_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_BIGINT_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_GBIGINT_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_BOOLEAN_OBJECT_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_BOOLEAN_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_I_ARRAY_LIKE_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_ITERABLE_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_ITERATOR_LIKE_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_ITERATOR_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_ITERATOR_ITERABLE_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_GENERATOR_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_DATE_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_I_OBJECT_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_I_ITERABLE_RESULT_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_I_TEMPLATE_ARRAY_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_I_THENABLE_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_NULL_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9257(this_.GetNativeType(StaticGet_JSTypeNative_NULL_TYPE()), "Null")
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_NUMBER_OBJECT_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_NUMBER_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_OBJECT_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_PROMISE_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_REGEXP_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_STRING_OBJECT_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_STRING_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_OBJECT_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9257(this_.GetNativeType(StaticGet_JSTypeNative_THENABLE_TYPE()), "Thenable")
	this_.priv_JSTypeRegistry_registerGlobalType_9256(this_.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE()))
	this_.priv_JSTypeRegistry_registerGlobalType_9257(this_.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE()), "Undefined")
	this_.priv_JSTypeRegistry_registerGlobalType_9257(this_.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE()), "void")
	this_.priv_JSTypeRegistry_registerGlobalType_9257(this_.GetNativeType(StaticGet_JSTypeNative_FUNCTION_TYPE()), "Function")
	this_.priv_JSTypeRegistry_registerGlobalType_9257(this_.GetNativeType(StaticGet_JSTypeNative_GLOBAL_THIS()), "Global")
}

// Java static method: JSTypeRegistry.getRootElementOfName(java.lang.String name)
func JSTypeRegistry_GetRootElementOfName(name string) string {
	clinit_JSTypeRegistry()
	_ = name
	var index int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
	_ = index
	if index != int32(-1) {
		return j2gort.JStrSub(name, int32(0), index)
	}
	return name
}

// Java static method: JSTypeRegistry.getLookupScope(com.google.javascript.rhino.StaticScope scope, java.lang.String name)
func JSTypeRegistry_GetLookupScope(scope StaticScope, name string) StaticScope {
	clinit_JSTypeRegistry()
	_ = scope
	_ = name
	if (scope != nil && !scope.IsJavaNil_()) && !(j2gort.IfaceEquals[StaticScope](j2gort.AsIface[StaticScope](scope.GetParentScope()), nil)) {
		return scope.GetTopmostScopeOfEventualDeclaration(JSTypeRegistry_GetRootElementOfName(name))
	}
	return scope
}

// Java method: JSTypeRegistry.getRootNodeForScope(com.google.javascript.rhino.StaticScope scope)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_getRootNodeForScope_9248(scope StaticScope) Node {
	_ = scope
	var root Node
	_ = root
	if scope != nil && !scope.IsJavaNil_() {
		root = scope.GetRootNode_0()
	} else {
		root = Node(0)
	}
	if ((root == 0) || root.IsRoot()) || root.IsScript() {
		return this_.nameTableGlobalRoot
	} else {
		return root
	}
}

// Java method: JSTypeRegistry.isDeclaredForScope(com.google.javascript.rhino.StaticScope scope, java.lang.String name)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_isDeclaredForScope_9249(scope StaticScope, name string) bool {
	_ = scope
	_ = name
	return !(j2gort.IfaceEquals[JSType](this_.priv_JSTypeRegistry_getTypeInternal_9250(scope, name), nil))
}

// Java static method: JSTypeRegistry.checkTypeName(java.lang.String typeName)
func JSTypeRegistry_CheckTypeName(typeName string) {
	clinit_JSTypeRegistry()
	_ = typeName
	if !(!(j2gort.JStrIndexOf(typeName, "<") >= 0)) {
		j2gort.CheckArgumentStr(false, "Type names cannot contain template annotations.")
	}
}

// Java method: JSTypeRegistry.getTypeInternal(com.google.javascript.rhino.StaticScope scope, java.lang.String name)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_getTypeInternal_9250(scope StaticScope, name string) JSType {
	_ = scope
	_ = name
	JSTypeRegistry_CheckTypeName(name)
	var syntheticTemplateScope *JSTypeRegistry_SyntheticTemplateScope
	_ = syntheticTemplateScope
	var ok_9 bool
	syntheticTemplateScope, ok_9 = j2gort.TryCast[*JSTypeRegistry_SyntheticTemplateScope](scope)
	if ok_9 {
		var type_ *TemplateType = syntheticTemplateScope.GetTemplateType(name)
		_ = type_
		if type_ != nil {
			return ToIface_JSType(type_)
		}
	}
	var declarationScope StaticScope = JSTypeRegistry_GetLookupScope(scope, name)
	_ = declarationScope
	var resolvedViaTable JSType = this_.priv_JSTypeRegistry_getTypeForScopeInternal_9255(declarationScope, name)
	_ = resolvedViaTable
	if resolvedViaTable != nil && !resolvedViaTable.IsJavaNil_() {
		return resolvedViaTable
	}
	var bestScope StaticScope
	_ = bestScope
	if declarationScope != nil && !declarationScope.IsJavaNil_() {
		bestScope = declarationScope
	} else {
		bestScope = scope
	}
	return this_.priv_JSTypeRegistry_resolveViaComponents_9251(bestScope, name)
}

// Java method: JSTypeRegistry.resolveViaComponents(com.google.javascript.rhino.StaticScope scope, java.lang.String qualifiedName)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_resolveViaComponents_9251(scope StaticScope, qualifiedName string) JSType {
	_ = scope
	_ = qualifiedName
	var resolutionScope StaticTypedScope
	_ = resolutionScope
	ok_10 := (j2gort.CheckNotNullString(qualifiedName) == "")
	if !ok_10 {
		var ok_11 bool
		resolutionScope, ok_11 = j2gort.TryCast[StaticTypedScope](scope)
		ok_10 = !ok_11
	}
	if ok_10 {
		return nil
	}
	if !(JSTypeRegistry_IsNameDefinedLocally(resolutionScope, JSTypeRegistry_GetRootElementOfName(qualifiedName))) {
		var resolvedViaClosureNamespace JSType = this_.ResolveViaClosureNamespace(qualifiedName)
		_ = resolvedViaClosureNamespace
		if resolvedViaClosureNamespace != nil && !resolvedViaClosureNamespace.IsJavaNil_() {
			return resolvedViaClosureNamespace
		}
	}
	return this_.priv_JSTypeRegistry_resolveViaProperties_9253(resolutionScope, qualifiedName)
}

// Java static method: JSTypeRegistry.isNameDefinedLocally(com.google.javascript.rhino.jstype.StaticTypedScope resolutionScope, java.lang.String reference)
func JSTypeRegistry_IsNameDefinedLocally(resolutionScope StaticTypedScope, reference string) bool {
	clinit_JSTypeRegistry()
	_ = resolutionScope
	_ = reference
	var slot StaticTypedSlot = j2gort.AsIface[StaticTypedSlot](resolutionScope.GetSlot(reference))
	_ = slot
	return (((slot != nil && !slot.IsJavaNil_()) && !(j2gort.IfaceEquals[StaticTypedScope](j2gort.AsIface[StaticTypedScope](slot.GetScope()), nil))) && !(j2gort.IfaceEquals[StaticTypedScope](j2gort.AsIface[StaticTypedScope](j2gort.AsIface[StaticTypedScope](slot.GetScope()).GetParentScope()), nil)))
}

// Java method: JSTypeRegistry.resolveViaClosureNamespace(java.lang.String reference)
func (this_ *JSTypeRegistry) ResolveViaClosureNamespace(reference string) JSType {
	_ = reference
	if this_.nonLegacyClosureNamespaceRoots.IsEmpty() || !(this_.nonLegacyClosureNamespaceRoots.Contains(JSTypeRegistry_GetRootElementOfName(reference))) {
		return nil
	}
	var prefix string = reference
	_ = prefix
	var firstComponentStart int32 = int32(-1)
	_ = firstComponentStart
	for {
		var namespace JSTypeRegistry_ClosureNamespace = this_.closureNamespaces.Get(prefix)
		_ = namespace
		if namespace != nil && !namespace.IsJavaNil_() {
			if namespace.IsLegacy() {
				return nil
			} else {
				return this_.priv_JSTypeRegistry_resolveViaPropertyGivenSlot_9254(namespace.Type(), namespace.DefinitionNode(), reference, firstComponentStart)
			}
		}
		var lastDot int32 = j2gort.JStrLastIndexOfCp(prefix, int32(uint16(46)))
		_ = lastDot
		if lastDot < int32(0) {
			return nil
		}
		firstComponentStart = (lastDot + int32(1))
		prefix = j2gort.JStrSub(reference, int32(0), lastDot)
	}
	panic("unreachable")
}

// Java method: JSTypeRegistry.resolveViaProperties(com.google.javascript.rhino.jstype.StaticTypedScope declarationScope, java.lang.String qualifiedName)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_resolveViaProperties_9253(declarationScope StaticTypedScope, qualifiedName string) JSType {
	_ = declarationScope
	_ = qualifiedName
	j2gort.CheckNotNullString(qualifiedName)
	if !(!(j2gort.CheckNotNullString(qualifiedName) == "")) {
		j2gort.CheckArgument(false)
	}
	var firstDot int32 = j2gort.JStrIndexOfCp(qualifiedName, int32(uint16(46)))
	_ = firstDot
	var rootName string
	_ = rootName
	if firstDot == int32(-1) {
		rootName = qualifiedName
	} else {
		rootName = j2gort.JStrSub(qualifiedName, int32(0), firstDot)
	}
	var slot StaticTypedSlot = j2gort.AsIface[StaticTypedSlot](declarationScope.GetOwnSlot(rootName))
	_ = slot
	if slot == nil || slot.IsJavaNil_() {
		return nil
	}
	var slotType JSType = slot.GetType()
	_ = slotType
	var iifeArg_1 int32
iifeSwitch_1:
	switch {
	default:

		if firstDot == int32(-1) {
			iifeArg_1 = int32(-1)
			break iifeSwitch_1
		}
		iifeArg_1 = (firstDot + int32(1))
		break iifeSwitch_1
	}
	return this_.priv_JSTypeRegistry_resolveViaPropertyGivenSlot_9254(slotType, Node(0), qualifiedName, iifeArg_1)
}

// Java method: JSTypeRegistry.resolveViaPropertyGivenSlot(com.google.javascript.rhino.jstype.JSType slotType, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node definitionNode, java.lang.String qualifiedName, int firstComponentStart)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_resolveViaPropertyGivenSlot_9254(slotType JSType, definitionNode Node, qualifiedName string, firstComponentStart int32) JSType {
	_ = slotType
	_ = definitionNode
	_ = qualifiedName
	_ = firstComponentStart
	if firstComponentStart < int32(0) {
		var typedefType JSType = JSTypeRegistry_ResolveTypeFromNodeIfTypedef(definitionNode)
		_ = typedefType
		if typedefType != nil && !typedefType.IsJavaNil_() {
			return typedefType
		}
	}
	if ((slotType == nil || slotType.IsJavaNil_()) || slotType.GetJSTypeBase().IsAllType()) || slotType.GetJSTypeBase().IsNoType() {
		return nil
	}
	if firstComponentStart >= int32(0) {
		var pos int32 = firstComponentStart
		_ = pos
		for {
			var nextDot int32 = j2gort.JStrIndexOfCpFrom(qualifiedName, int32(uint16(46)), pos)
			_ = nextDot
			var isLast bool = (nextDot == int32(-1))
			_ = isLast
			var component string
			_ = component
			if isLast {
				component = j2gort.JStrSubFrom(qualifiedName, pos)
			} else {
				component = j2gort.JStrSub(qualifiedName, pos, nextDot)
			}
			var parentObj ObjectType = ObjectType_Cast(slotType)
			_ = parentObj
			if ((parentObj == nil || parentObj.IsJavaNil_()) || (j2gort.CheckNotNullString(component) == "")) || !(parentObj.GetObjectTypeBase().HasOwnProperty_1_string(component)) {
				return nil
			}
			if isLast {
				var def Node = parentObj.GetObjectTypeBase().GetPropertyDefSite_1_string(component)
				_ = def
				var typedefType_2 JSType = JSTypeRegistry_ResolveTypeFromNodeIfTypedef(def)
				_ = typedefType_2
				if typedefType_2 != nil && !typedefType_2.IsJavaNil_() {
					return typedefType_2
				}
				slotType = parentObj.GetObjectTypeBase().GetPropertyType_1_string(component)
				break
			}
			slotType = parentObj.GetObjectTypeBase().GetPropertyType_1_string(component)
			pos = (nextDot + int32(1))
		}
	}
	var enumType *EnumType
	_ = enumType
	if slotType == nil || slotType.IsJavaNil_() {
		return nil
	} else {
		if slotType.GetJSTypeBase().IsFunctionType() && (slotType.GetJSTypeBase().IsConstructor() || slotType.GetJSTypeBase().IsInterface()) {
			return Upcast_JSType(slotType.ToMaybeFunctionType().GetInstanceType())
		} else {
			if slotType.GetJSTypeBase().IsNoObjectType() {
				return Upcast_JSType(this_.GetNativeObjectType(StaticGet_JSTypeNative_NO_OBJECT_TYPE()))
			} else {
				var ok_12 bool
				enumType, ok_12 = j2gort.TryCast[*EnumType](slotType)
				if ok_12 {
					return ToIface_JSType(enumType.GetElementsType())
				} else {
					return nil
				}
			}
		}
	}
}

// Java static method: JSTypeRegistry.resolveTypeFromNodeIfTypedef(com.google.javascript.rhino.Node node)
func JSTypeRegistry_ResolveTypeFromNodeIfTypedef(node Node) JSType {
	clinit_JSTypeRegistry()
	_ = node
	if node == 0 {
		return nil
	}
	return node.GetTypedefTypeProp()
}

// Java method: JSTypeRegistry.getTypeForScopeInternal(com.google.javascript.rhino.StaticScope scope, java.lang.String name)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_getTypeForScopeInternal_9255(scope StaticScope, name string) JSType {
	_ = scope
	_ = name
	var rootNode Node = this_.priv_JSTypeRegistry_getRootNodeForScope_9248(scope)
	_ = rootNode
	var type_ JSType = this_.scopedNameTable.Get(rootNode, name)
	_ = type_
	return type_
}

// Java method: JSTypeRegistry.registerGlobalType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_registerGlobalType_9256(type_ JSType) {
	_ = type_
	this_.priv_JSTypeRegistry_register_9258(nil, type_, type_.GetJSTypeBase(

	// Java method: JSTypeRegistry.registerGlobalType(com.google.javascript.rhino.jstype.JSType type, java.lang.String name)
	).String())
}

func (this_ *JSTypeRegistry) priv_JSTypeRegistry_registerGlobalType_9257(type_ JSType, name string) {
	_ = type_
	_ = name
	this_.priv_JSTypeRegistry_register_9258(nil, type_, name)
}

// Java method: JSTypeRegistry.register(com.google.javascript.rhino.@org.jspecify.annotations.Nullable StaticScope scope, com.google.javascript.rhino.jstype.JSType type, java.lang.String name)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_register_9258(scope StaticScope, type_ JSType, name string) {
	_ = scope
	_ = type_
	_ = name
	JSTypeRegistry_CheckTypeName(name)
	this_.priv_JSTypeRegistry_registerForScope_9259(JSTypeRegistry_GetLookupScope(scope, name), type_, name)
}

// Java method: JSTypeRegistry.registerForScope(com.google.javascript.rhino.StaticScope scope, com.google.javascript.rhino.jstype.JSType type, java.lang.String name)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_registerForScope_9259(scope StaticScope, type_ JSType, name string) {
	_ = scope
	_ = type_
	_ = name
	this_.scopedNameTable.Put(this_.priv_JSTypeRegistry_getRootNodeForScope_9248(scope), name, type_)
	if j2gort.JStrIndexOfCp(name, int32(uint16(46))) >= int32(0) {
		this_.priv_JSTypeRegistry_addQualifiedScopedNameHash_9260(Node_HashQualifiedNameString(name))
	}
}

// Java method: JSTypeRegistry.addQualifiedScopedNameHash(int h)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_addQualifiedScopedNameHash_9260(h int32) {
	_ = h
	if h == int32(0) {
		return
	}
	var table []int32 = this_.qualifiedScopedNameHashes
	_ = table
	if table == nil {
		table = make([]int32, int32(32))
		this_.qualifiedScopedNameHashes = table
	} else if (this_.qualifiedScopedNameHashCount * int32(4)) >= (int32(len(table)) * int32(3)) {
		var next []int32 = make([]int32, (int32(len(table)) * int32(2)))
		_ = next
		var nextMask int32 = (int32(len(next)) - int32(1))
		_ = nextMask
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = int32(len(table))
			_ = len_
			for ; i < len_; i = i + 1 {
				var existing int32 = table[i]
				_ = existing
				if existing != int32(0) {
					var s int32 = (existing & nextMask)
					_ = s
					for next[s] != int32(0) {
						s = ((s + int32(1)) & nextMask)
					}
					{
						a_13 := next
						i_14 := s
						a_13[i_14] = existing
					}
				}
			}
		}
		table = next
		this_.qualifiedScopedNameHashes = table
	}
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var slot int32 = (h & mask)
	_ = slot
	for {
		var existing_2 int32 = table[slot]
		_ = existing_2
		if existing_2 == int32(0) {
			{
				a_15 := table
				i_16 := slot
				a_15[i_16] = h
			}
			this_.qualifiedScopedNameHashCount = this_.qualifiedScopedNameHashCount + 1
			return
		}
		if existing_2 == h {
			return
		}
		slot = ((slot + int32(1)) & mask)
	}
}

// Java method: JSTypeRegistry.containsQualifiedScopedNameHash(int h)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_containsQualifiedScopedNameHash_9261(h int32) bool {
	_ = h
	var table []int32 = this_.qualifiedScopedNameHashes
	_ = table
	if (table == nil) || (h == int32(0)) {
		return false
	}
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var slot int32 = (h & mask)
	_ = slot
	for {
		var existing int32 = table[slot]
		_ = existing
		if existing == int32(0) {
			return false
		}
		if existing == h {
			return true
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: JSTypeRegistry.mayHaveQualifiedTypeForGetProp(com.google.javascript.rhino.@org.jspecify.annotations.Nullable StaticScope scope, java.lang.String rootName, com.google.javascript.rhino.Node getPropNode)
func (this_ *JSTypeRegistry) MayHaveQualifiedTypeForGetProp(scope StaticScope, rootName string, getPropNode Node) bool {
	_ = scope
	_ = rootName
	_ = getPropNode
	if j2gort.InstanceOf[*JSTypeRegistry_SyntheticTemplateScope](scope) {
		return true
	}
	if !(this_.nonLegacyClosureNamespaceRoots.IsEmpty()) && this_.nonLegacyClosureNamespaceRoots.Contains(rootName) {
		return true
	}
	var h int32 = getPropNode.HashQualifiedName()
	_ = h
	if h == int32(0) {
		return false
	}
	if this_.priv_JSTypeRegistry_containsQualifiedScopedNameHash_9261(h) {
		return true
	}
	var declarationScope StaticScope
	_ = declarationScope
	if (scope != nil && !scope.IsJavaNil_()) && !(j2gort.IfaceEquals[StaticScope](j2gort.AsIface[StaticScope](scope.GetParentScope()), nil)) {
		declarationScope = scope.GetTopmostScopeOfEventualDeclaration(rootName)
	} else {
		declarationScope = scope
	}
	var bestScope StaticScope
	_ = bestScope
	if declarationScope != nil && !declarationScope.IsJavaNil_() {
		bestScope = declarationScope
	} else {
		bestScope = scope
	}
	var resolutionScope StaticTypedScope
	_ = resolutionScope
	var ok_17 bool
	resolutionScope, ok_17 = j2gort.TryCast[StaticTypedScope](bestScope)
	if !ok_17 {
		return false
	}
	var slot StaticTypedSlot = j2gort.AsIface[StaticTypedSlot](resolutionScope.GetOwnSlot(rootName))
	_ = slot
	if slot == nil || slot.IsJavaNil_() {
		return false
	}
	var slotType JSType = slot.GetType()
	_ = slotType
	if ((slotType == nil || slotType.IsJavaNil_()) || slotType.GetJSTypeBase().IsAllType()) || slotType.GetJSTypeBase().IsNoType() {
		return false
	}
	return JSTypeRegistry_CanResolveGetPropViaSlot(slotType, getPropNode)
}

// Java static method: JSTypeRegistry.canResolveGetPropViaSlot(com.google.javascript.rhino.jstype.JSType rootSlotType, com.google.javascript.rhino.Node getPropNode)
func JSTypeRegistry_CanResolveGetPropViaSlot(rootSlotType JSType, getPropNode Node) bool {
	clinit_JSTypeRegistry()
	_ = rootSlotType
	_ = getPropNode
	var root Node = getPropNode.GetFirstChild()
	_ = root
	for root.IsGetProp() {
		root = root.GetFirstChild()
	}
	if root.IsName() && (j2gort.JStrIndexOfCp(root.GetString(), int32(uint16(46))) >= int32(0)) {
		return true
	}
	var ownerType JSType = JSTypeRegistry_StepGetPropOwnerType(rootSlotType, getPropNode.GetFirstChild())
	_ = ownerType
	if ownerType == nil || ownerType.IsJavaNil_() {
		return false
	}
	var ownerObj ObjectType = ObjectType_Cast(ownerType)
	_ = ownerObj
	var lastProp string = getPropNode.GetString()
	_ = lastProp
	return (((ownerObj != nil && !ownerObj.IsJavaNil_()) && !(j2gort.CheckNotNullString(lastProp) == "")) && ownerObj.GetObjectTypeBase().HasOwnProperty_1_string(lastProp))
}

// Java static method: JSTypeRegistry.stepGetPropOwnerType(com.google.javascript.rhino.jstype.JSType rootSlotType, com.google.javascript.rhino.Node node)
func JSTypeRegistry_StepGetPropOwnerType(rootSlotType JSType, node Node) JSType {
	clinit_JSTypeRegistry()
	_ = rootSlotType
	_ = node
	if !(node.IsGetProp()) {
		return rootSlotType
	}
	var parentType JSType = JSTypeRegistry_StepGetPropOwnerType(rootSlotType, node.GetFirstChild())
	_ = parentType
	if parentType == nil || parentType.IsJavaNil_() {
		return nil
	}
	var parentObj ObjectType = ObjectType_Cast(parentType)
	_ = parentObj
	var prop string = node.GetString()
	_ = prop
	if ((parentObj == nil || parentObj.IsJavaNil_()) || (j2gort.CheckNotNullString(prop) == "")) || !(parentObj.GetObjectTypeBase().HasOwnProperty_1_string(prop)) {
		return nil
	}
	return parentObj.GetObjectTypeBase().GetPropertyType_1_string(prop)
}

// Java method: JSTypeRegistry.registerNativeType(com.google.javascript.rhino.jstype.JSTypeNative typeId, com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_registerNativeType_9263(typeId JSTypeNative, type_ JSType) {
	_ = typeId
	_ = type_
	{
		a_18 := this_.nativeTypes
		i_19 := (typeId).Ordinal()
		a_18[i_19] = type_
	}
}

// Java static method: JSTypeRegistry.isObjectLiteralThatCanBeSkipped(com.google.javascript.rhino.jstype.JSType t)
func JSTypeRegistry_IsObjectLiteralThatCanBeSkipped(t JSType) bool {
	clinit_JSTypeRegistry()
	_ = t
	t = t.RestrictByNotNullOrUndefined()
	return (t.GetJSTypeBase().IsRecordType() || t.GetJSTypeBase().IsLiteralObject())
}

// Java method: JSTypeRegistry.registerDroppedPropertiesInUnion(com.google.javascript.rhino.jstype.RecordType subtype, com.google.javascript.rhino.jstype.RecordType supertype)
func (this_ *JSTypeRegistry) RegisterDroppedPropertiesInUnion(subtype *RecordType, supertype *RecordType) {
	_ = subtype
	_ = supertype
	var subMap *PropertyMap = subtype.GetOwnPropertyMapOrNull()
	_ = subMap
	if subMap == nil {
		return
	}
	var foundDroppedProperty bool = false
	_ = foundDroppedProperty
	if subMap.HasFlatProps() {
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = subMap.GetOwnPropertiesCount()
			_ = n
			for ; i < n; i = i + 1 {
				var k string = subMap.GetFlatPropAt(i).GetName()
				_ = k
				if !(supertype.HasProperty_1_string(k)) {
					foundDroppedProperty = true
					this_.droppedPropertiesOfUnions.Add(k)
				}
			}
		}
	} else {
		{
			cur_20 := j2gort.NewIterCursor(j2gort.CheckNotNull(subMap.GetOwnPropertyNames()))
			for cur_20.HasNext() {
				pname := cur_20.NextStr()
				_ = pname
				if !(supertype.HasProperty_1_string(pname)) {
					foundDroppedProperty = true
					this_.droppedPropertiesOfUnions.Add(pname)
				}
			}
		}
	}
	if foundDroppedProperty {
		var superMap *PropertyMap = supertype.GetOwnPropertyMapOrNull()
		_ = superMap
		if superMap != nil {
			superMap.AddOwnPropertyNamesTo(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.propertiesOfSupertypesInUnions))
		}
	}
}

// Java method: JSTypeRegistry.registerPropertyOnType(java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeRegistry) RegisterPropertyOnType(propertyName string, type_ JSType) {
	_ = propertyName
	_ = type_
	if type_.GetJSTypeBase().IsUnionType() {
		var union *UnionType = type_.ToMaybeUnionType()
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = union.GetAlternatesCount()
			_ = n
			for ; i < n; i = i + 1 {
				this_.RegisterPropertyOnType(propertyName, union.GetAlternateAt(i))
			}
		}
		return
	}
	if JSTypeRegistry_IsObjectLiteralThatCanBeSkipped(type_) {
		type_ = this_.priv_JSTypeRegistry_getSentinelObjectLiteral_9221()
	}
	var objType ObjectType
	_ = objType
	var ok_21 bool
	objType, ok_21 = j2gort.TryCast[ObjectType](type_)
	if ok_21 && (objType.GetObjectTypeBase().HasReferenceName()) {
		var set *j2gort.Ext_java_util_LinkedHashSet = this_.eachRefTypeIndexedByProperty.Get(propertyName)
		_ = set
		if set == nil {
			set = j2gort.NewExt_java_util_LinkedHashSet()
			this_.eachRefTypeIndexedByProperty.Put(propertyName, set)
		}
		set.Add_1(objType)
	} else {
		var set_2 *j2gort.Ext_java_util_LinkedHashSet = this_.nonRefTypesIndexedByProperty.Get(propertyName)
		_ = set_2
		if set_2 == nil {
			set_2 = j2gort.NewExt_java_util_LinkedHashSet()
			this_.nonRefTypesIndexedByProperty.Put(propertyName, set_2)
		}
		set_2.Add_1(type_)
	}
}

// Java method: JSTypeRegistry.canPropertyBeDefined(com.google.javascript.rhino.jstype.JSType type, java.lang.String propertyName)
func (this_ *JSTypeRegistry) CanPropertyBeDefined(type_ JSType, propertyName string) JSTypeRegistry_PropDefinitionKind {
	_ = type_
	_ = propertyName
	base_type__ := type_.GetJSTypeBase()
	if type_.IsStruct() {
		switch (base_type__.GetPropertyKind_1_string(propertyName)).Ordinal() {
		case int32(1) /* HasPropertyKind.KNOWN_PRESENT */ :
			return StaticGet_JSTypeRegistry_PropDefinitionKind_KNOWN()
		case int32(2) /* HasPropertyKind.MAYBE_PRESENT */ :
			return StaticGet_JSTypeRegistry_PropDefinitionKind_LOOSE_UNION()
		case int32(0) /* HasPropertyKind.ABSENT */ :
			return StaticGet_JSTypeRegistry_PropDefinitionKind_UNKNOWN()
		}
		panic("unreachable switch expression")
	} else {
		if !(base_type__.IsEmptyType()) && !(base_type__.IsUnknownType()) {
			switch (base_type__.GetPropertyKind_1_string(propertyName)).Ordinal() {
			case int32(1) /* HasPropertyKind.KNOWN_PRESENT */ :
				return StaticGet_JSTypeRegistry_PropDefinitionKind_KNOWN()
			case int32(2) /* HasPropertyKind.MAYBE_PRESENT */ :
				return StaticGet_JSTypeRegistry_PropDefinitionKind_LOOSE_UNION()
			case int32(0) /* HasPropertyKind.ABSENT */ :
			}
		}
		var nonRefTypes *j2gort.Ext_java_util_LinkedHashSet = this_.nonRefTypesIndexedByProperty.Get(propertyName)
		_ = nonRefTypes
		if nonRefTypes != nil {
			{
				cur_22 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(nonRefTypes))
				for cur_22.HasNext() {
					rawElem_23 := cur_22.Next()
					alternative := j2gort.AsIface[JSType](rawElem_23)
					_ = alternative
					var greatestSubtype JSType = alternative.GetJSTypeBase().GetGreatestSubtype(type_)
					_ = greatestSubtype
					if !(greatestSubtype.GetJSTypeBase().IsEmptyType()) {
						var maybeRecordType *RecordType = greatestSubtype.ToMaybeRecordType()
						_ = maybeRecordType
						if (maybeRecordType != nil) && maybeRecordType.IsSynthetic() {
							continue
						}
						return StaticGet_JSTypeRegistry_PropDefinitionKind_LOOSE()
					}
				}
			}
		}
		var refTypes *j2gort.Ext_java_util_LinkedHashSet = this_.eachRefTypeIndexedByProperty.Get(propertyName)
		_ = refTypes
		if refTypes != nil {
			{
				cur_24 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(refTypes))
				for cur_24.HasNext() {
					rawElem_25 := cur_24.Next()
					alternative_2 := j2gort.AsIface[ObjectType](rawElem_25)
					_ = alternative_2
					var greatestSubtype_2 JSType = alternative_2.GetObjectTypeBase().GetGreatestSubtype(type_)
					_ = greatestSubtype_2
					if !(greatestSubtype_2.GetJSTypeBase().IsEmptyType()) {
						var maybeRecordType_2 *RecordType = greatestSubtype_2.ToMaybeRecordType()
						_ = maybeRecordType_2
						if (maybeRecordType_2 != nil) && maybeRecordType_2.IsSynthetic() {
							continue
						}
						return StaticGet_JSTypeRegistry_PropDefinitionKind_LOOSE()
					}
				}
			}
		}
		if type_.ToMaybeRecordType() != nil {
			var rec *RecordType = type_.ToMaybeRecordType()
			_ = rec
			var recMap *PropertyMap = rec.GetOwnPropertyMapOrNull()
			_ = recMap
			var mayBeInUnion bool = false
			_ = mayBeInUnion
			if recMap != nil {
				if recMap.HasFlatProps() {
					{
						var i int32 = int32(0)
						_ = i
						var n int32 = recMap.GetOwnPropertiesCount()
						_ = n
						for ; i < n; i = i + 1 {
							if j2gort.SetContainsStr(this_.propertiesOfSupertypesInUnions, recMap.GetFlatPropAt(i).GetName()) {
								mayBeInUnion = true
								break
							}
						}
					}
				} else {
					{
						cur_26 := j2gort.NewIterCursor(j2gort.CheckNotNull(recMap.GetOwnPropertyNames()))
						for cur_26.HasNext() {
							pname := cur_26.NextStr()
							_ = pname
							if j2gort.SetContainsStr(this_.propertiesOfSupertypesInUnions, pname) {
								mayBeInUnion = true
								break
							}
						}
					}
				}
			}
			if mayBeInUnion && this_.droppedPropertiesOfUnions.Contains(propertyName) {
				return StaticGet_JSTypeRegistry_PropDefinitionKind_LOOSE()
			}
		}
	}
	return StaticGet_JSTypeRegistry_PropDefinitionKind_UNKNOWN()
}

// Java method: JSTypeRegistry.getEachReferenceTypeWithProperty(java.lang.String propertyName)
func (this_ *JSTypeRegistry) GetEachReferenceTypeWithProperty(propertyName string) j2gort.Ext_java_lang_Iterable {
	_ = propertyName
	var set *j2gort.Ext_java_util_LinkedHashSet = this_.eachRefTypeIndexedByProperty.Get(propertyName)
	_ = set
	if set != nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Collections_UnmodifiableSet(j2gort.AsIface[j2gort.Ext_java_util_Set](set)))
	} else {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
}

// Java method: JSTypeRegistry.findCommonSuperObject(com.google.javascript.rhino.jstype.ObjectType a, com.google.javascript.rhino.jstype.ObjectType b)
func (this_ *JSTypeRegistry) FindCommonSuperObject(a ObjectType, b ObjectType) ObjectType {
	_ = a
	_ = b
	var stackA j2gort.Ext_java_util_List = JSTypeRegistry_GetSuperStack(a)
	_ = stackA
	var stackB j2gort.Ext_java_util_List = JSTypeRegistry_GetSuperStack(b)
	_ = stackB
	var result ObjectType = this_.GetNativeObjectType(StaticGet_JSTypeNative_OBJECT_TYPE())
	_ = result
	for !(stackA.IsEmpty()) && !(stackB.IsEmpty()) {
		var currentA ObjectType = j2gort.AsIface[ObjectType](stackA.Remove_1_int32((stackA.Size() - int32(1))))
		_ = currentA
		var currentB ObjectType = j2gort.AsIface[ObjectType](stackB.Remove_1_int32((stackB.Size() - int32(1))))
		_ = currentB
		if currentA.GetObjectTypeBase().GetJSTypeBase().EqualsJSType(currentB) {
			result = currentA
		} else {
			return result
		}
	}
	return result
}

// Java static method: JSTypeRegistry.getSuperStack(com.google.javascript.rhino.jstype.ObjectType a)
func JSTypeRegistry_GetSuperStack(a ObjectType) j2gort.Ext_java_util_List {
	clinit_JSTypeRegistry()
	_ = a
	var stack []ObjectType = make([]ObjectType, 0, int32(5))
	_ = stack
	{
		var current ObjectType = a
		_ = current
		for ; current != nil && !current.IsJavaNil_(); current = current.GetImplicitPrototype() {
			stack = append(stack, current)
		}
	}
	return j2gort.WrapTypedArrayList(stack)
}

// Java method: JSTypeRegistry.declareType(com.google.javascript.rhino.StaticScope scope, java.lang.String name, com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeRegistry) DeclareType(scope StaticScope, name string, type_ JSType) bool {
	_ = scope
	_ = name
	_ = type_
	if !(!(j2gort.CheckNotNullString(name) == "")) {
		j2gort.CheckState(false)
	}
	if !(j2gort.IfaceEquals[JSType](this_.priv_JSTypeRegistry_getTypeForScopeInternal_9255(JSTypeRegistry_GetLookupScope(scope, name), name), nil)) {
		return false
	}
	this_.priv_JSTypeRegistry_register_9258(scope, type_, name)
	return true
}

// Java method: JSTypeRegistry.declareTypeForExactScope(com.google.javascript.rhino.StaticScope scope, java.lang.String name, com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeRegistry) DeclareTypeForExactScope(scope StaticScope, name string, type_ JSType) bool {
	_ = scope
	_ = name
	_ = type_
	if !(!(j2gort.CheckNotNullString(name) == "")) {
		j2gort.CheckState(false)
	}
	if !(j2gort.IfaceEquals[JSType](this_.priv_JSTypeRegistry_getTypeForScopeInternal_9255(scope, name), nil)) {
		return false
	}
	this_.priv_JSTypeRegistry_registerForScope_9259(scope, type_, name)
	return true
}

// Java method: JSTypeRegistry.overwriteDeclaredType(com.google.javascript.rhino.StaticScope scope, java.lang.String name, com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeRegistry) OverwriteDeclaredType(scope StaticScope, name string, type_ JSType) {
	_ = scope
	_ = name
	_ = type_
	if !(this_.priv_JSTypeRegistry_isDeclaredForScope_9249(scope, name)) {
		j2gort.CheckStateFmt(false, "missing name %s", j2gort.BoxString(name))
	}
	this_.priv_JSTypeRegistry_register_9258(scope, type_, name)
}

// Java method: JSTypeRegistry.isForwardDeclaredType(java.lang.String name)
func (this_ *JSTypeRegistry) IsForwardDeclaredType(name string) bool {
	_ = name
	return j2gort.SetContainsStr(this_.forwardDeclaredTypes, name)
}

// Java method: JSTypeRegistry.getReadableTypeName(com.google.javascript.rhino.Node n)
func (this_ *JSTypeRegistry) GetReadableTypeName(n Node) string {
	_ = n
	return this_.GetReadableJSTypeName(n, true)
}

// Java method: JSTypeRegistry.getReadableTypeNameNoDeref(com.google.javascript.rhino.Node n)
func (this_ *JSTypeRegistry) GetReadableTypeNameNoDeref(n Node) string {
	_ = n
	return this_.GetReadableJSTypeName(n, false)
}

// Java static method: JSTypeRegistry.getSimpleReadableJSTypeName(com.google.javascript.rhino.jstype.JSType type)
func JSTypeRegistry_GetSimpleReadableJSTypeName(type_ JSType) string {
	clinit_JSTypeRegistry()
	_ = type_
	if j2gort.InstanceOf[*AllType](type_) {
		return type_.GetJSTypeBase().String()
	} else if j2gort.InstanceOf[ValueType](type_) {
		return type_.GetJSTypeBase().String()
	} else if type_.IsFunctionPrototypeType() {
		return type_.GetJSTypeBase().String()
	} else if j2gort.InstanceOf[ObjectType](type_) {
		if !(JSTypeRegistry_IsReadableObjectType(type_.GetJSTypeBase().ToMaybeObjectType())) {
			return j2gort.NullString
		}
		var source Node = type_.GetJSTypeBase().ToObjectType().GetConstructor().GetSource()
		_ = source
		if source != 0 {
			if !(source.IsFunction() || source.IsClass()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(source))
			}
			var readable string = source.GetFirstChild().GetOriginalName()
			_ = readable
			if readable != j2gort.NullString {
				return readable
			}
		}
		return type_.GetJSTypeBase().String()
	} else if j2gort.InstanceOf[*UnionType](type_) {
		var unionType *UnionType = type_.ToMaybeUnionType()
		_ = unionType
		var union string = j2gort.NullString
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = unionType.GetAlternatesCount()
			_ = n
			for ; i < n; i = i + 1 {
				var name string = JSTypeRegistry_GetSimpleReadableJSTypeName(unionType.GetAlternateAt(i))
				_ = name
				if name == j2gort.NullString {
					return j2gort.NullString
				}
				if union == j2gort.NullString {
					union = ("(" + j2gort.StrValString(name))
				} else {
					union = (j2gort.StrValString(union) + "|" + j2gort.StrValString(name))
				}
			}
		}
		union = (j2gort.StrValString(union) + ")")
		return union
	}
	return j2gort.NullString
}

// Java method: JSTypeRegistry.getReadableJSTypeName(com.google.javascript.rhino.Node n, boolean dereference)
func (this_ *JSTypeRegistry) GetReadableJSTypeName(n Node, dereference bool) string {
	_ = n
	_ = dereference
	var type_ JSType = this_.priv_JSTypeRegistry_getJSTypeOrUnknown_9276(n)
	_ = type_
	if dereference {
		var autoboxed JSType = type_.Autobox()
		_ = autoboxed
		if !(autoboxed.GetJSTypeBase().IsNoType()) {
			type_ = autoboxed
		}
	}
	var name string = JSTypeRegistry_GetSimpleReadableJSTypeName(type_)
	_ = name
	if name != j2gort.NullString {
		return name
	}
	if n.IsGetProp() {
		var objectType ObjectType = this_.priv_JSTypeRegistry_getJSTypeOrUnknown_9276(n.GetFirstChild()).GetJSTypeBase().Dereference()
		_ = objectType
		if objectType != nil && !objectType.IsJavaNil_() {
			var propName string = n.GetString()
			_ = propName
			objectType = objectType.GetObjectTypeBase().GetClosestDefiningType(propName)
			if JSTypeRegistry_IsReadableObjectType(objectType) {
				return (j2gort.StrVal(objectType) + "." + j2gort.StrValString(propName))
			}
		}
	}
	if n.IsQualifiedName() {
		return n.GetQualifiedName()
	} else if type_.GetJSTypeBase().IsFunctionType() {
		return "function"
	} else {
		return type_.GetJSTypeBase().String()
	}
}

// Java static method: JSTypeRegistry.isReadableObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func JSTypeRegistry_IsReadableObjectType(type_ ObjectType) bool {
	clinit_JSTypeRegistry()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return false
	} else if type_.IsInstanceType() {
		return true
	} else if type_.GetObjectTypeBase().IsFunctionPrototypeType() {
		return true
	} else if type_.GetObjectTypeBase().IsEnumElementType() {
		return (type_.ToMaybeEnumElementType().GetConstructor() != nil)
	}
	return false
}

// Java method: JSTypeRegistry.getJSTypeOrUnknown(com.google.javascript.rhino.Node n)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_getJSTypeOrUnknown_9276(n Node) JSType {
	_ = n
	var jsType JSType = n.GetJSType()
	_ = jsType
	if jsType == nil || jsType.IsJavaNil_() {
		return this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
	} else {
		return jsType
	}
}

// Java method: JSTypeRegistry.getGlobalType(java.lang.String jsTypeName)
func (this_ *JSTypeRegistry) GetGlobalType(jsTypeName string) JSType {
	_ = jsTypeName
	return this_.GetType_2(nil, jsTypeName)
}

// Java method: JSTypeRegistry.getType(com.google.javascript.rhino.@org.jspecify.annotations.Nullable StaticScope scope, java.lang.String jsTypeName)
func (this_ *JSTypeRegistry) GetType_2(scope StaticScope, jsTypeName string) JSType {
	_ = scope
	_ = jsTypeName
	return this_.priv_JSTypeRegistry_getTypeInternal_9250(scope, jsTypeName)
}

// Java method: JSTypeRegistry.getType(com.google.javascript.rhino.jstype.StaticTypedScope scope, java.lang.String jsTypeName, java.lang.String sourceName, int lineno, int charno)
func (this_ *JSTypeRegistry) GetType_5(scope StaticTypedScope, jsTypeName string, sourceName string, lineno int32, charno int32) JSType {
	_ = scope
	_ = jsTypeName
	_ = sourceName
	_ = lineno
	_ = charno
	switch j2gort.CheckNotNullString(jsTypeName) {
	case "boolean":
		return this_.GetNativeType(StaticGet_JSTypeNative_BOOLEAN_TYPE())
	case "number":
		return this_.GetNativeType(StaticGet_JSTypeNative_NUMBER_TYPE())
	case "bigint":
		return this_.GetNativeType(StaticGet_JSTypeNative_BIGINT_TYPE())
	case "string":
		return this_.GetNativeType(StaticGet_JSTypeNative_STRING_TYPE())
	case "undefined", "void":
		return this_.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE())
	}
	var type_ JSType = nil
	_ = type_
	var thisType JSType = nil
	_ = thisType
	if (scope != nil && !scope.IsJavaNil_()) && !(j2gort.IfaceEquals[JSType](scope.GetTypeOfThis(), nil)) {
		thisType = Upcast_JSType(scope.GetTypeOfThis().GetJSTypeBase().ToObjectType())
	}
	if thisType != nil && !thisType.IsJavaNil_() {
		type_ = ToIface_JSType(thisType.GetTemplateTypeMap().GetLastTemplateTypeKeyByName(jsTypeName))
		if type_ != nil && !type_.IsJavaNil_() {
			if !(type_.GetJSTypeBase().IsTemplateType()) {
				j2gort.CheckStateFmt(false, "expected:%s", type_)
			}
			return type_
		}
	}
	type_ = this_.GetType_2(Upcast_StaticScope(scope), jsTypeName)
	if type_ == nil || type_.IsJavaNil_() {
		return ToIface_JSType(this_.CreateNamedType(scope, jsTypeName, sourceName, lineno, charno))
	}
	return type_
}

// Java method: JSTypeRegistry.getNativeType(com.google.javascript.rhino.jstype.JSTypeNative typeId)
func (this_ *JSTypeRegistry) GetNativeType(typeId JSTypeNative) JSType {
	_ = typeId
	return this_.nativeTypes[(typeId).Ordinal()]
}

// Java method: JSTypeRegistry.getNativeObjectType(com.google.javascript.rhino.jstype.JSTypeNative typeId)
func (this_ *JSTypeRegistry) GetNativeObjectType(typeId JSTypeNative) ObjectType {
	_ = typeId
	return j2gort.AsIface[ObjectType](this_.GetNativeType(typeId))
}

// Java method: JSTypeRegistry.getNativeFunctionType(com.google.javascript.rhino.jstype.JSTypeNative typeId)
func (this_ *JSTypeRegistry) GetNativeFunctionType(typeId JSTypeNative) *FunctionType {
	_ = typeId
	return j2gort.CastPtr[FunctionType](this_.GetNativeType(typeId))
}

// Java method: JSTypeRegistry.getResolver()
func (this_ *JSTypeRegistry) GetResolver() *JSTypeResolver {
	return this_.resolver
}

// Java method: JSTypeRegistry.isResolverClosed()
func (this_ *JSTypeRegistry) IsResolverClosed() bool {
	return ((this_.resolver != nil) && this_.resolver.IsClosed())
}

// Java method: JSTypeRegistry.evaluateTypeExpressionInGlobalScope(com.google.javascript.rhino.JSTypeExpression expr)
func (this_ *JSTypeRegistry) EvaluateTypeExpressionInGlobalScope(expr JSTypeExpression) JSType {
	_ = expr
	return expr.Evaluate(nil, this_)
}

// Java method: JSTypeRegistry.createOptionalType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeRegistry) CreateOptionalType(type_ JSType) JSType {
	_ = type_
	if j2gort.InstanceOf[*UnknownType](type_) || type_.GetJSTypeBase().IsAllType() {
		return type_
	} else {
		return this_.CreateUnionType_2(type_, this_.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE()))
	}
}

// Java method: JSTypeRegistry.createNullableType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeRegistry) CreateNullableType(type_ JSType) JSType {
	_ = type_
	return this_.CreateUnionType_2(type_, this_.GetNativeType(StaticGet_JSTypeNative_NULL_TYPE()))
}

// Java method: JSTypeRegistry.createOptionalNullableType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeRegistry) CreateOptionalNullableType(type_ JSType) JSType {
	_ = type_
	var builder *UnionType_Builder = this_.AcquireUnionTypeBuilder()
	_ = builder
	defer func() {
		this_.ReleaseUnionTypeBuilder(builder)
	}()
	{
		builder.AddAlternate(j2gort.CheckNotNullIface(type_))
		builder.AddAlternate(this_.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE()))
		builder.AddAlternate(this_.GetNativeType(StaticGet_JSTypeNative_NULL_TYPE()))
		return builder.Build()
	}
	panic("unreachable")
}

// Java method: JSTypeRegistry.createUnionType(com.google.javascript.rhino.jstype.JSType a, com.google.javascript.rhino.jstype.JSType b)
func (this_ *JSTypeRegistry) CreateUnionType_2(a JSType, b JSType) JSType {
	_ = a
	_ = b
	j2gort.CheckNotNullIface(a)
	j2gort.CheckNotNullIface(b)
	if (a == b || ((a == nil || a.IsJavaNil_()) && (b == nil || b.IsJavaNil_()))) && ((a.GetJSTypeBase().typeKindMask & (int32(1) | int32(64))) == int32(0)) {
		return a
	}
	if this_.resolver.IsClosed() && (((int32(a.GetJSTypeBase().jstypeFlags) & int32(b.GetJSTypeBase().jstypeFlags)) & int32(int8(8))) != int32(0)) {
		var h int32 = ((a.GetJSTypeBase().typeId * int32(-1640531527)) ^ (b.GetJSTypeBase().typeId * int32(-2048144789)))
		_ = h
		var slot int32 = ((h ^ (javalang.Ushr32(h, int(int32(16))))) & int32(16383))
		_ = slot
		var cacheA []JSType = this_.union2A
		_ = cacheA
		if ((cacheA != nil) && j2gort.IfaceEquals[JSType](cacheA[slot], a)) && j2gort.IfaceEquals[JSType](this_.union2B[slot], b) {
			return this_.union2Res[slot]
		}
		var builder *UnionType_Builder = this_.AcquireUnionTypeBuilder()
		_ = builder
		var res JSType = nil
		_ = res
		func() {
			defer func() {
				this_.ReleaseUnionTypeBuilder(builder)
			}()
			builder.AddAlternate(a)
			builder.AddAlternate(b)
			res = builder.Build()
		}()
		if cacheA == nil {
			cacheA = make([]JSType, (int32(16383) + int32(1)))
			this_.union2A = cacheA
			this_.union2B = make([]JSType, (int32(16383) + int32(1)))
			this_.union2Res = make([]JSType, (int32(16383) + int32(1)))
		}
		{
			a_27 := cacheA
			i_28 := slot
			a_27[i_28] = a
		}
		{
			a_29 := this_.union2B
			i_30 := slot
			a_29[i_30] = b
		}
		{
			a_31 := this_.union2Res
			i_32 := slot
			a_31[i_32] = res
		}
		return res
	}
	var builder_2 *UnionType_Builder = this_.AcquireUnionTypeBuilder()
	_ = builder_2
	defer func() {
		this_.ReleaseUnionTypeBuilder(builder_2)
	}()
	{
		builder_2.AddAlternate(a)
		builder_2.AddAlternate(b)
		return builder_2.Build()
	}
	panic("unreachable")
}

// Java method: JSTypeRegistry.createUnionType(com.google.javascript.rhino.jstype.JSType[] variants)
func (this_ *JSTypeRegistry) CreateUnionType_1_JSType(variants []JSType) JSType {
	_ = variants
	var builder *UnionType_Builder = this_.AcquireUnionTypeBuilder()
	_ = builder
	defer func() {
		this_.ReleaseUnionTypeBuilder(builder)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			for ; i < int32(len(variants)); i = i + 1 {
				builder.AddAlternate(j2gort.CheckNotNullIface(variants[i]))
			}
		}
		return builder.Build()
	}
	panic("unreachable")
}

// Java method: JSTypeRegistry.createUnionType(java.util.List<? extends com.google.javascript.rhino.jstype.JSType> variants)
func (this_ *JSTypeRegistry) CreateUnionType_1_Ext_java_util_List(variants j2gort.Ext_java_util_List) JSType {
	_ = variants
	var builder *UnionType_Builder = this_.AcquireUnionTypeBuilder()
	_ = builder
	defer func() {
		this_.ReleaseUnionTypeBuilder(builder)
	}()
	{
		builder.AddAlternates_1_Ext_java_util_List(variants)
		return builder.Build()
	}
	panic("unreachable")
}

// Java method: JSTypeRegistry.createUnionType(com.google.javascript.rhino.jstype.JSTypeNative[] variants)
func (this_ *JSTypeRegistry) CreateUnionType_1_JSTypeNative(variants []JSTypeNative) JSType {
	_ = variants
	var builder *UnionType_Builder = this_.AcquireUnionTypeBuilder()
	_ = builder
	defer func() {
		this_.ReleaseUnionTypeBuilder(builder)
	}()
	{
		for _, el_33 := range j2gort.CheckNotNullSlice(variants) {
			type_ := el_33
			_ = type_
			builder.AddAlternate(this_.GetNativeType(type_))
		}
		return builder.Build()
	}
	panic("unreachable")
}

// Java method: JSTypeRegistry.createFunctionType(com.google.javascript.rhino.jstype.JSType returnType, com.google.javascript.rhino.jstype.JSType[] parameterTypes)
func (this_ *JSTypeRegistry) CreateFunctionType_2_JSType_JSType(returnType JSType, parameterTypes []JSType) *FunctionType {
	_ = returnType
	_ = parameterTypes
	return this_.CreateFunctionType_2_JSType_Ext_java_util_List(returnType, j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateParameters_1(parameterTypes)))
}

// Java method: JSTypeRegistry.createFunctionType(com.google.javascript.rhino.jstype.JSType returnType, java.util.List<com.google.javascript.rhino.jstype.FunctionType.Parameter> parameters)
func (this_ *JSTypeRegistry) CreateFunctionType_2_JSType_Ext_java_util_List(returnType JSType, parameters j2gort.Ext_java_util_List) *FunctionType {
	_ = returnType
	_ = parameters
	var fb_ FunctionType_Builder
	InitFunctionType_Builder(&fb_, this_)
	(&fb_).WithParameters_1(parameters)
	(&fb_).WithReturnType_1(returnType)
	return j2gort.CastPtr[FunctionType]((&fb_).Build())
}

// Java method: JSTypeRegistry.createFunctionTypeWithVarArgs(com.google.javascript.rhino.jstype.JSType returnType, com.google.javascript.rhino.jstype.JSType[] parameterTypes)
func (this_ *JSTypeRegistry) CreateFunctionTypeWithVarArgs(returnType JSType, parameterTypes []JSType) *FunctionType {
	_ = returnType
	_ = parameterTypes
	return this_.CreateFunctionType_2_JSType_Ext_java_util_List(returnType, j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateParametersWithVarArgs(parameterTypes)))
}

// Java method: JSTypeRegistry.createNativeFunctionTypeWithVarArgs(com.google.javascript.rhino.jstype.JSType returnType, com.google.javascript.rhino.jstype.JSType[] parameterTypes)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_createNativeFunctionTypeWithVarArgs_9296(returnType JSType, parameterTypes []JSType) *FunctionType {
	_ = returnType
	_ = parameterTypes
	return this_.priv_JSTypeRegistry_createNativeFunctionType_9306(returnType, j2gort.AsIface[j2gort.Ext_java_util_List](this_.CreateParametersWithVarArgs(parameterTypes)))
}

// Java method: JSTypeRegistry.createFunctionTypeWithInstanceType(com.google.javascript.rhino.jstype.ObjectType instanceType, com.google.javascript.rhino.jstype.JSType returnType, java.util.List<com.google.javascript.rhino.jstype.JSType> parameterTypes)
func (this_ *JSTypeRegistry) CreateFunctionTypeWithInstanceType(instanceType ObjectType, returnType JSType, parameterTypes j2gort.Ext_java_util_List) JSType {
	_ = instanceType
	_ = returnType
	_ = parameterTypes
	var iifeArg_2 []JSType
iifeSwitch_2:
	switch {
	default:

		arr_37 := parameterTypes.ToArray_1_slice_any(func() []any {
			arr_34 := make([]JSType, int32(0))
			if arr_34 == nil {
				return nil
			}
			out_35 := make([]any, len(arr_34))
			for i_, el_36 := range arr_34 {
				out_35[i_] = el_36
			}
			return out_35
		}())
		if arr_37 == nil {
			iifeArg_2 = nil
			break iifeSwitch_2
		}
		out_38 := make([]JSType, len(arr_37))
		for i_, el_39 := range arr_37 {
			out_38[i_] = j2gort.AsIface[JSType](el_39)
		}
		iifeArg_2 = out_38
		break iifeSwitch_2
	}
	var paramsNode j2gort.Ext_com_google_common_collect_ImmutableList = this_.CreateParameters_1(iifeArg_2)
	_ = paramsNode
	var fb_ FunctionType_Builder
	InitFunctionType_Builder(&fb_, this_)
	(&fb_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](paramsNode))
	(&fb_).WithReturnType_1(returnType)
	(&fb_).WithTypeOfThis(Upcast_JSType(instanceType))
	return ToIface_JSType(j2gort.CastPtr[FunctionType]((&fb_).Build()))
}

// Java method: JSTypeRegistry.createParameters(com.google.javascript.rhino.jstype.JSType[] parameterTypes)
func (this_ *JSTypeRegistry) CreateParameters_1(parameterTypes []JSType) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = parameterTypes
	return this_.priv_JSTypeRegistry_createParameters_9299(false, parameterTypes)
}

// Java method: JSTypeRegistry.createParameters(boolean lastVarArgs, com.google.javascript.rhino.jstype.JSType[] parameterTypes)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_createParameters_9299(lastVarArgs bool, parameterTypes []JSType) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = lastVarArgs
	_ = parameterTypes
	var val_40 FunctionParamBuilder
	InitFunctionParamBuilder_1(&val_40, this_)
	var builder *FunctionParamBuilder = &val_40
	_ = builder
	var max_ int32 = (int32(len(parameterTypes)) - int32(1))
	_ = max_
	{
		var i int32 = int32(0)
		_ = i
		for ; i <= max_; i = i + 1 {
			if lastVarArgs && (i == max_) {
				builder.AddVarArgs(parameterTypes[i])
			} else {
				builder.AddRequiredParams([]JSType{parameterTypes[i]})
			}
		}
	}
	return builder.Build()
}

// Java method: JSTypeRegistry.getDefaultVarArgsParameterArray()
func (this_ *JSTypeRegistry) GetDefaultVarArgsParameterArray() []FunctionType_Parameter {
	var cached []FunctionType_Parameter = this_.defaultVarArgsParameterArray
	_ = cached
	if cached == nil {
		cached = []FunctionType_Parameter{FunctionType_Parameter_Create(this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE()), false, true)}
		this_.defaultVarArgsParameterArray = cached
	}
	return cached
}

// Java method: JSTypeRegistry.getUnknownRequiredParameter()
func (this_ *JSTypeRegistry) GetUnknownRequiredParameter() FunctionType_Parameter {
	var p FunctionType_Parameter = this_.unknownRequiredParameter
	_ = p
	if p == (FunctionType_Parameter{}) {
		p = FunctionType_Parameter_Create(this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE()), false, false)
		this_.unknownRequiredParameter = p
	}
	return p
}

// Java method: JSTypeRegistry.getUnknownRequiredParameterArray(int count)
func (this_ *JSTypeRegistry) GetUnknownRequiredParameterArray(count int32) []FunctionType_Parameter {
	_ = count
	if count == int32(0) {
		return StaticGet_FunctionType_EMPTY_PARAMETERS()
	}
	var p FunctionType_Parameter = this_.GetUnknownRequiredParameter()
	_ = p
	if count <= int32(4) {
		var table [][]FunctionType_Parameter = this_.unknownRequiredParameterArrays
		_ = table
		if table == nil {
			table = [][]FunctionType_Parameter{StaticGet_FunctionType_EMPTY_PARAMETERS(), []FunctionType_Parameter{p}, []FunctionType_Parameter{p, p}, []FunctionType_Parameter{p, p, p}, []FunctionType_Parameter{p, p, p, p}}
			this_.unknownRequiredParameterArrays = table
		}
		return table[count]
	}
	var arr []FunctionType_Parameter = make([]FunctionType_Parameter, count)
	_ = arr
	{
		var i int32 = int32(0)
		_ = i
		for ; i < count; i = i + 1 {
			{
				a_41 := arr
				i_42 := i
				a_41[i_42] = p
			}
		}
	}
	return arr
}

// Java method: JSTypeRegistry.createParametersWithVarArgs(com.google.javascript.rhino.jstype.JSType[] parameterTypes)
func (this_ *JSTypeRegistry) CreateParametersWithVarArgs(parameterTypes []JSType) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = parameterTypes
	return this_.priv_JSTypeRegistry_createParameters_9299(true, parameterTypes)
}

// Java method: JSTypeRegistry.createOptionalParameters(com.google.javascript.rhino.jstype.JSType[] parameterTypes)
func (this_ *JSTypeRegistry) CreateOptionalParameters(parameterTypes []JSType) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = parameterTypes
	var val_43 FunctionParamBuilder
	InitFunctionParamBuilder_1(&val_43, this_)
	var builder *FunctionParamBuilder = &val_43
	_ = builder
	builder.AddOptionalParams(parameterTypes)
	return builder.Build()
}

// Java method: JSTypeRegistry.createFunctionTypeWithNewReturnType(com.google.javascript.rhino.jstype.FunctionType existingFunctionType, com.google.javascript.rhino.jstype.JSType returnType)
func (this_ *JSTypeRegistry) CreateFunctionTypeWithNewReturnType(existingFunctionType *FunctionType, returnType JSType) *FunctionType {
	_ = existingFunctionType
	_ = returnType
	var fb_ FunctionType_Builder
	existingFunctionType.InitBuilder_(&fb_)
	(&fb_).WithReturnType_1(returnType)
	return j2gort.CastPtr[FunctionType]((&fb_).Build())
}

// Java method: JSTypeRegistry.createNativeFunctionType(com.google.javascript.rhino.jstype.JSType returnType, java.util.List<com.google.javascript.rhino.jstype.FunctionType.Parameter> parameters)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_createNativeFunctionType_9306(returnType JSType, parameters j2gort.Ext_java_util_List) *FunctionType {
	_ = returnType
	_ = parameters
	var fb_ FunctionType_Builder
	InitFunctionType_Builder(&fb_, this_)
	(&fb_).WithParameters_1(parameters)
	(&fb_).WithReturnType_1(returnType)
	(&fb_).ForNativeType()
	return j2gort.CastPtr[FunctionType]((&fb_).Build())
}

// Java method: JSTypeRegistry.buildRecordTypeFromObject(com.google.javascript.rhino.jstype.ObjectType objType)
func (this_ *JSTypeRegistry) BuildRecordTypeFromObject(objType ObjectType) JSType {
	_ = objType
	var recType *RecordType = objType.ToMaybeRecordType()
	_ = recType
	if recType != nil {
		return ToIface_JSType(recType)
	}
	var propNames j2gort.Ext_java_util_Set = objType.GetOwnPropertyNames()
	_ = propNames
	if propNames.IsEmpty() {
		return this_.GetNativeType(StaticGet_JSTypeNative_OBJECT_TYPE())
	}
	var props j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.NewExt_com_google_common_collect_ImmutableMap_Builder()
	_ = props
	{
		cur_44 := j2gort.NewIterCursor(j2gort.CheckNotNull(propNames))
		for cur_44.HasNext() {
			propName := cur_44.NextStr()
			_ = propName
			props.Put_2(j2gort.BoxString(propName), objType.GetObjectTypeBase().GetPropertyType_1_string(propName))
		}
	}
	return this_.CreateRecordType(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](props.BuildOrThrow()))
}

// Java method: JSTypeRegistry.createRecordType(java.util.Map<java.lang.String,? extends com.google.javascript.rhino.jstype.JSType> props)
func (this_ *JSTypeRegistry) CreateRecordType(props j2gort.Ext_java_util_Map) JSType {
	_ = props
	var propMap j2gort.Ext_java_util_Map = props
	_ = propMap
	var val_45 RecordTypeBuilder
	InitRecordTypeBuilder(&val_45, this_)
	var builder *RecordTypeBuilder = &val_45
	_ = builder
	{
		cur_46 := j2gort.NewMapEntryCursor(propMap)
		for cur_46.NextEntry() {
			builder.AddProperty(cur_46.KeyStr(), j2gort.AsIface[JSType](cur_46.GetValue()), Node(0))
		}
	}
	return builder.Build()
}

// Java method: JSTypeRegistry.createObjectType(java.lang.String name, com.google.javascript.rhino.jstype.ObjectType implicitPrototype)
func (this_ *JSTypeRegistry) CreateObjectType(name string, implicitPrototype ObjectType) ObjectType {
	_ = name
	_ = implicitPrototype
	return Upcast_ObjectType(j2gort.AsIface[PrototypeObjectType](j2gort.AsIface[PrototypeObjectType_Builder](PrototypeObjectType_S_Builder(this_).SetName(name)).GetPrototypeObjectType_BuilderBase().SetImplicitPrototype(implicitPrototype).Build()))
}

// Java method: JSTypeRegistry.createAnonymousObjectType(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info)
func (this_ *JSTypeRegistry) CreateAnonymousObjectType(info *JSDocInfo) ObjectType {
	_ = info
	var type_ PrototypeObjectType = j2gort.AsIface[PrototypeObjectType](PrototypeObjectType_S_Builder(this_).GetPrototypeObjectType_BuilderBase().SetAnonymous(true).Build())
	_ = type_
	base_type__ := type_.GetPrototypeObjectTypeBase()
	base_type__.SetPrettyPrint(true)
	base_type__.SetJSDocInfo(info)
	return Upcast_ObjectType(type_)
}

// Java method: JSTypeRegistry.resetImplicitPrototype(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.jstype.ObjectType newImplicitProto)
func (this_ *JSTypeRegistry) ResetImplicitPrototype(type_ JSType, newImplicitProto ObjectType) {
	_ = type_
	_ = newImplicitProto
	var poType PrototypeObjectType
	_ = poType
	var ok_47 bool
	poType, ok_47 = j2gort.TryCast[PrototypeObjectType](type_)
	if ok_47 {
		poType.ClearCachedValues()
		poType.GetPrototypeObjectTypeBase().SetImplicitPrototype(newImplicitProto)
		this_.union2A = nil
		this_.union2B = nil
		this_.union2Res = nil
	}
}

// Java method: JSTypeRegistry.createConstructorType(java.lang.String name, com.google.javascript.rhino.Node source, java.util.List<com.google.javascript.rhino.jstype.FunctionType.Parameter> parameters, com.google.javascript.rhino.jstype.JSType returnType, com.google.common.collect.@org.jspecify.annotations.Nullable ImmutableList<com.google.javascript.rhino.jstype.TemplateType> templateKeys, boolean isAbstract)
func (this_ *JSTypeRegistry) CreateConstructorType(name string, source Node, parameters j2gort.Ext_java_util_List, returnType JSType, templateKeys j2gort.Ext_com_google_common_collect_ImmutableList, isAbstract bool) *FunctionType {
	_ = name
	_ = source
	_ = parameters
	_ = returnType
	_ = templateKeys
	_ = isAbstract
	if !(((source == 0) || source.IsFunction()) || source.IsClass()) {
		j2gort.CheckArgument(false)
	}
	var fb_ FunctionType_Builder
	InitFunctionType_Builder(&fb_, this_)
	(&fb_).ForConstructor()
	(&fb_).WithName(name)
	(&fb_).WithSourceNode(source)
	(&fb_).WithParameters_1(parameters)
	(&fb_).WithReturnType_1(returnType)
	(&fb_).WithTemplateKeys_1_Ext_com_google_common_collect_ImmutableList(templateKeys)
	(&fb_).WithIsAbstract(isAbstract)
	return j2gort.CastPtr[FunctionType]((&fb_).Build())
}

// Java method: JSTypeRegistry.createTemplateType(java.lang.String name)
func (this_ *JSTypeRegistry) CreateTemplateType_1(name string) *TemplateType {
	_ = name
	return NewTemplateType_2(this_, name)
}

// Java method: JSTypeRegistry.createTemplateType(java.lang.String name, com.google.javascript.rhino.jstype.JSType bound)
func (this_ *JSTypeRegistry) CreateTemplateType_2(name string, bound JSType) *TemplateType {
	_ = name
	_ = bound
	return NewTemplateType_3_JSTypeRegistry_string_JSType(this_, name, bound)
}

// Java method: JSTypeRegistry.createTemplateTypeWithTransformation(java.lang.String name, com.google.javascript.rhino.Node expr)
func (this_ *JSTypeRegistry) CreateTemplateTypeWithTransformation(name string, expr Node) *TemplateType {
	_ = name
	_ = expr
	return NewTemplateType_3_JSTypeRegistry_string_Node(this_, name, expr)
}

// Java method: JSTypeRegistry.getEmptyTemplateTypeMap()
func (this_ *JSTypeRegistry) GetEmptyTemplateTypeMap() *TemplateTypeMap {
	return this_.emptyTemplateTypeMap
}

// Java method: JSTypeRegistry.createTemplatizedType(com.google.javascript.rhino.jstype.ObjectType baseType, com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.JSType> templatizedTypes)
func (this_ *JSTypeRegistry) CreateTemplatizedType_2_ObjectType_Ext_com_google_common_collect_ImmutableList(baseType ObjectType, templatizedTypes j2gort.Ext_com_google_common_collect_ImmutableList) *TemplatizedType {
	_ = baseType
	_ = templatizedTypes
	j2gort.CheckNotNullIface(baseType)
	return allocPooled_TemplatizedType(&this_.slab_TemplatizedType, this_, baseType, templatizedTypes)
}

// Java method: JSTypeRegistry.createTemplatizedType(com.google.javascript.rhino.jstype.ObjectType baseType, java.util.Map<com.google.javascript.rhino.jstype.TemplateType,com.google.javascript.rhino.jstype.JSType> templatizedTypes)
func (this_ *JSTypeRegistry) CreateTemplatizedType_2_ObjectType_Ext_java_util_Map(baseType ObjectType, templatizedTypes j2gort.Ext_java_util_Map) *TemplatizedType {
	_ = baseType
	_ = templatizedTypes
	var unknownType JSType = this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
	_ = unknownType
	var typeParams []*TemplateType = j2gort.ImmListToSliceNullable[*TemplateType](baseType.GetObjectTypeBase().GetTypeParameters())
	_ = typeParams
	var n int32 = int32(len(typeParams))
	_ = n
	if n == int32(0) {
		return this_.CreateTemplatizedType_2_ObjectType_Ext_com_google_common_collect_ImmutableList(baseType, j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, n)
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			builder.AddVoid_(j2gort.AsIface[JSType](templatizedTypes.GetOrDefault(j2gort.PtrToAny(typeParams[i]), unknownType)))
		}
	}
	return this_.CreateTemplatizedType_2_ObjectType_Ext_com_google_common_collect_ImmutableList(baseType, builder.BuildImmList_())
}

// Java method: JSTypeRegistry.createTemplatizedType(com.google.javascript.rhino.jstype.ObjectType baseType, com.google.javascript.rhino.jstype.JSType[] templatizedTypes)
func (this_ *JSTypeRegistry) CreateTemplatizedType_2_ObjectType_JSType(baseType ObjectType, templatizedTypes []JSType) *TemplatizedType {
	_ = baseType
	_ = templatizedTypes
	var iifeArg_3 []any
iifeSwitch_3:
	switch {
	default:

		arr_48 := templatizedTypes
		if arr_48 == nil {
			iifeArg_3 = nil
			break iifeSwitch_3
		}
		out_49 := make([]any, len(arr_48))
		for i_, el_50 := range arr_48 {
			out_49[i_] = el_50
		}
		iifeArg_3 = out_49
		break iifeSwitch_3
	}
	return this_.CreateTemplatizedType_2_ObjectType_Ext_com_google_common_collect_ImmutableList(baseType, j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_slice_any(iifeArg_3))
}

// Java method: JSTypeRegistry.createNamedType(com.google.javascript.rhino.jstype.StaticTypedScope scope, java.lang.String reference, java.lang.String sourceName, int lineno, int charno)
func (this_ *JSTypeRegistry) CreateNamedType(scope StaticTypedScope, reference string, sourceName string, lineno int32, charno int32) *NamedType {
	_ = scope
	_ = reference
	_ = sourceName
	_ = lineno
	_ = charno
	return NamedType_S_Builder(this_, reference).SetScope(scope).SetResolutionKind(StaticGet_NamedType_ResolutionKind_TYPE_NAME()).SetErrorReportingLocation(sourceName, lineno, charno).Build()
}

// Java method: JSTypeRegistry.identifyNonNullableName(com.google.javascript.rhino.@org.jspecify.annotations.Nullable StaticScope scope, java.lang.String name)
func (this_ *JSTypeRegistry) IdentifyNonNullableName(scope StaticScope, name string) {
	_ = scope
	_ = name
	j2gort.CheckNotNullString(name)
	var lookupScope StaticScope = JSTypeRegistry_GetLookupScope(scope, name)
	_ = lookupScope
	this_.nonNullableTypeNames.Put(this_.priv_JSTypeRegistry_getRootNodeForScope_9248(lookupScope), name)
}

// Java method: JSTypeRegistry.isNonNullableName(com.google.javascript.rhino.StaticScope scope, java.lang.String name)
func (this_ *JSTypeRegistry) IsNonNullableName(scope StaticScope, name string) bool {
	_ = scope
	_ = name
	j2gort.CheckNotNullString(name)
	scope = JSTypeRegistry_GetLookupScope(scope, name)
	return this_.nonNullableTypeNames.ContainsEntry(this_.priv_JSTypeRegistry_getRootNodeForScope_9248(scope), name)
}

// Java method: JSTypeRegistry.evaluateTypeExpression(com.google.javascript.rhino.JSTypeExpression expr, com.google.javascript.rhino.jstype.StaticTypedScope scope)
func (this_ *JSTypeRegistry) EvaluateTypeExpression(expr JSTypeExpression, scope StaticTypedScope) JSType {
	_ = expr
	_ = scope
	return this_.CreateTypeFromCommentNode_3(expr.GetRoot(), expr.GetSourceName(), scope)
}

// Java method: JSTypeRegistry.createTypeFromCommentNode(com.google.javascript.rhino.Node n)
func (this_ *JSTypeRegistry) CreateTypeFromCommentNode_1(n Node) JSType {
	_ = n
	return this_.CreateTypeFromCommentNode_3(n, "[internal]", nil)
}

// Java method: JSTypeRegistry.createTypeFromCommentNode(com.google.javascript.rhino.Node n, java.lang.String sourceName, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable StaticTypedScope scope)
func (this_ *JSTypeRegistry) CreateTypeFromCommentNode_3(n Node, sourceName string, scope StaticTypedScope) JSType {
	_ = n
	_ = sourceName
	_ = scope
	{

		switch n.TokenSwitchOrd_() {
		case int32(162) /* Token.LC */ :
			return this_.priv_JSTypeRegistry_createRecordTypeFromNodes_9328(n.GetFirstChild(), sourceName, scope)
		case int32(159) /* Token.BANG */ :
			var child JSType = this_.CreateTypeFromCommentNode_3(n.GetFirstChild(), sourceName, scope)
			_ = child
			var namedType *NamedType
			_ = namedType
			var ok_51 bool
			namedType, ok_51 = j2gort.TryCast[*NamedType](child)
			if ok_51 {
				return namedType.GetBangType()
			}
			return child.RestrictByNotNullOrUndefined()
		case int32(158) /* Token.QMARK */ :
			var firstChild Node = n.GetFirstChild()
			_ = firstChild
			if firstChild == 0 {
				return this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
			}
			return this_.CreateNullableType(this_.CreateTypeFromCommentNode_3(firstChild, sourceName, scope))
		case int32(160) /* Token.EQUALS */ :
			return this_.CreateOptionalType(this_.CreateTypeFromCommentNode_3(n.GetFirstChild(), sourceName, scope))
		case int32(125) /* Token.ITER_REST */ :
			return this_.CreateTypeFromCommentNode_3(n.GetFirstChild(), sourceName, scope)
		case int32(156) /* Token.STAR */ :
			return this_.GetNativeType(StaticGet_JSTypeNative_ALL_TYPE())
		case int32(155) /* Token.PIPE */ :
			var builder *UnionType_Builder = this_.AcquireUnionTypeBuilder()
			_ = builder
			var didRet_52 bool
			var retVal_53 JSType
			func() {
				defer func() {
					this_.ReleaseUnionTypeBuilder(builder)
				}()
				{
					var child_2 Node = n.GetFirstChild()
					_ = child_2
					for ; child_2 != 0; child_2 = child_2.GetNext() {
						builder.AddAlternate(j2gort.CheckNotNullIface(this_.CreateTypeFromCommentNode_3(child_2, sourceName, scope)))
					}
				}
				retVal_53 = builder.Build()
				didRet_52 = true
				return
			}()
			if didRet_52 {
				return retVal_53
			}
		case int32(88) /* Token.EMPTY */ :
			return this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
		case int32(87) /* Token.VOID */ :
			return this_.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE())
		case int32(25) /* Token.TYPEOF */ :
			var name string = n.GetFirstChild().GetString()
			_ = name
			var firstDot int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
			_ = firstDot
			var root string
			_ = root
			if firstDot < int32(0) {
				root = name
			} else {
				root = j2gort.JStrSub(name, int32(0), firstDot)
			}
			var declarationScope StaticScope = scope.GetTopmostScopeOfEventualDeclaration(root)
			_ = declarationScope
			var rootSlot StaticSlot = Upcast_StaticSlot(j2gort.AsIface[StaticTypedSlot](scope.GetSlot(root)))
			_ = rootSlot
			var type_ JSType = scope.LookupQualifiedName_1_string(name)
			_ = type_
			if ((type_ == nil || type_.IsJavaNil_()) || type_.GetJSTypeBase().IsUnknownType()) || !(j2gort.IfaceEquals[StaticScope](j2gort.AsIface[StaticScope](rootSlot.GetScope()), declarationScope)) {
				return ToIface_JSType(NamedType_S_Builder(this_, ("typeof " + j2gort.StrValString(name))).SetScope(scope).SetResolutionKind(StaticGet_NamedType_ResolutionKind_TYPEOF()).SetErrorReportingLocationFrom(n).Build())
			}
			if type_.GetJSTypeBase().IsLiteralObject() {
				var scopeType JSType = type_
				_ = scopeType
				type_ = ToIface_JSType(NamedType_S_Builder(this_, ("typeof " + j2gort.StrValString(name))).SetResolutionKind(StaticGet_NamedType_ResolutionKind_NONE()).SetReferencedType(scopeType).Build())
			}
			return type_
		case int32(35) /* Token.STRINGLIT */ :
			var nominalType JSType = this_.GetType_5(scope, n.GetString(), sourceName, n.GetLineno(), n.GetCharno())
			_ = nominalType
			var templateArgs j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_JSTypeRegistry_parseTemplateArgs_9327(nominalType, n, sourceName, scope)
			_ = templateArgs
			if nominalType.GetJSTypeBase().IsNamedType() && !(nominalType.GetJSTypeBase().IsResolved()) {
				if templateArgs != nil {
					nominalType = ToIface_JSType(nominalType.ToMaybeNamedType().ToBuilder().SetTemplateTypes(templateArgs).Build())
				}
				return this_.priv_JSTypeRegistry_addNullabilityBasedOnParseContext_9326(n, nominalType, Upcast_StaticScope(scope))
			}
			var objectType ObjectType
			_ = objectType
			var ok_54 bool
			objectType, ok_54 = j2gort.TryCast[ObjectType](nominalType)
			if !ok_54 {
				return nominalType
			}
			if (templateArgs == nil) || !(nominalType.GetJSTypeBase().IsRawTypeOfTemplatizedType()) {
				return this_.priv_JSTypeRegistry_addNullabilityBasedOnParseContext_9326(n, nominalType, Upcast_StaticScope(scope))
			}
			return this_.priv_JSTypeRegistry_addNullabilityBasedOnParseContext_9326(n, ToIface_JSType(this_.CreateTemplatizedType_2_ObjectType_Ext_com_google_common_collect_ImmutableList(objectType, templateArgs)), Upcast_StaticScope(scope))
		case int32(73) /* Token.FUNCTION */ :
			var thisType JSType = nil
			_ = thisType
			var isConstructor bool = false
			_ = isConstructor
			var current Node = n.GetFirstChild()
			_ = current
			if current.IsThis() || current.IsNew() {
				var contextNode Node = current.GetFirstChild()
				_ = contextNode
				var candidateThisType JSType = this_.CreateTypeFromCommentNode_3(contextNode, sourceName, scope)
				_ = candidateThisType
				if candidateThisType.GetJSTypeBase().IsNullType() || candidateThisType.GetJSTypeBase().IsVoidType() {
					thisType = candidateThisType
				} else if current.IsThis() {
					thisType = candidateThisType.RestrictByNotNullOrUndefined()
				} else if current.IsNew() {
					thisType = Upcast_JSType(ObjectType_Cast(candidateThisType.RestrictByNotNullOrUndefined()))
					if thisType == nil || thisType.IsJavaNil_() {
						this_.reporter.Warning(StaticGet_Msg_JSDOC_FUNCTION_NEWNOTOBJECT().Format_0(), sourceName, contextNode.GetLineno(), contextNode.GetCharno())
					}
				}
				isConstructor = (current.RawToken_() == uint8(23))
				current = current.GetNext()
			}
			var val_55 FunctionParamBuilder
			InitFunctionParamBuilder_1(&val_55, this_)
			var paramBuilder *FunctionParamBuilder = &val_55
			_ = paramBuilder
			if current.RawToken_() == uint8(49) {
				{
					var arg Node = current.GetFirstChild()
					_ = arg
					for ; arg != 0; arg = arg.GetNext() {
						if arg.RawToken_() == uint8(125) {
							if !(arg.HasChildren()) {
								paramBuilder.AddVarArgs(this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE()))
							} else {
								paramBuilder.AddVarArgs(this_.CreateTypeFromCommentNode_3(arg.GetFirstChild(), sourceName, scope))
							}
						} else {
							var type__2 JSType = this_.CreateTypeFromCommentNode_3(arg, sourceName, scope)
							_ = type__2
							if arg.RawToken_() == uint8(160) {
								var addSuccess bool = paramBuilder.AddOptionalParams([]JSType{type__2})
								_ = addSuccess
								if !(addSuccess) {
									this_.reporter.Warning(StaticGet_Msg_JSDOC_FUNCTION_VARARGS().Format_0(), sourceName, arg.GetLineno(), arg.GetCharno())
								}
							} else {
								paramBuilder.AddRequiredParams([]JSType{type__2})
							}
						}
					}
				}
				current = current.GetNext()
			}
			var returnType JSType = this_.CreateTypeFromCommentNode_3(current, sourceName, scope)
			_ = returnType
			var fb_ FunctionType_Builder
			InitFunctionType_Builder(&fb_, this_)
			(&fb_).WithParameters_1(j2gort.AsIface[j2gort.Ext_java_util_List](paramBuilder.Build()))
			(&fb_).WithReturnType_1(returnType)
			(&fb_).WithTypeOfThis(thisType)
			(&fb_).WithKind(func() FunctionType_Kind {
				if isConstructor {
					return StaticGet_FunctionType_Kind_CONSTRUCTOR()
				}
				return StaticGet_FunctionType_Kind_ORDINARY()
			}())
			return ToIface_JSType(j2gort.CastPtr[FunctionType]((&fb_).Build()))
		default:
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected node in type expression: "+j2gort.StrVal(n))))
		}
		panic("unreachable switch expression")
	}

}

// Java method: JSTypeRegistry.addNullabilityBasedOnParseContext(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.StaticScope scope)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_addNullabilityBasedOnParseContext_9326(n Node, type_ JSType, scope StaticScope) JSType {
	_ = n
	_ = type_
	_ = scope
	if !(n.IsName() || n.IsStringLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	j2gort.CheckNotNullIface(type_)
	if this_.IsNonNullableName(scope, n.GetString()) {
		return type_
	} else if type_.GetJSTypeBase().IsTemplateType() {
		return type_
	} else if n.HasParent() && (n.GetParent().RawToken_() == uint8(159)) {
		return type_
	} else {
		return this_.CreateNullableType(type_)
	}
}

// Java method: JSTypeRegistry.parseTemplateArgs(com.google.javascript.rhino.jstype.JSType nominalType, com.google.javascript.rhino.Node typeNode, java.lang.String sourceName, com.google.javascript.rhino.jstype.StaticTypedScope scope)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_parseTemplateArgs_9327(nominalType JSType, typeNode Node, sourceName string, scope StaticTypedScope) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = nominalType
	_ = typeNode
	_ = sourceName
	_ = scope
	var typeList Node = typeNode.GetFirstChild()
	_ = typeList
	if typeList == 0 {
		return nil
	}
	var templateArgs j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = templateArgs
	{
		var templateNode Node = typeList.GetFirstChild()
		_ = templateNode
		for ; templateNode != 0; templateNode = templateNode.GetNext() {
			templateArgs.Add_1(this_.CreateTypeFromCommentNode_3(templateNode, sourceName, scope))
		}
	}
	var isObject bool = ((j2gort.CheckNotNullString(typeNode.GetString()) == "Object") || (j2gort.CheckNotNullString(typeNode.GetString()) == "window.Object"))
	_ = isObject
	if isObject && (templateArgs.Size() == int32(1)) {
		templateArgs.Add_2(int32(0), this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE()))
	}
	base_nominalType_ := nominalType.GetJSTypeBase()
	if base_nominalType_.IsNamedType() && !(base_nominalType_.IsResolved()) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(templateArgs)
	}
	var requiredTemplateArgCount int32 = nominalType.GetTemplateParamCount()
	_ = requiredTemplateArgCount
	if templateArgs.Size() <= requiredTemplateArgCount {
		return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(templateArgs)
	}
	if !(base_nominalType_.IsUnknownType()) && !(this_.IsNonNullableName(Upcast_StaticScope(scope), typeNode.GetString())) {
		var firstExtraTemplateParam Node = typeNode.GetFirstChild().GetChildAtIndex(requiredTemplateArgCount)
		_ = firstExtraTemplateParam
		var message string = ("Too many template parameters\nFound " + j2gort.StrValInt32(templateArgs.Size()) + ", required at most " + j2gort.StrValInt32(requiredTemplateArgCount))
		_ = message
		this_.reporter.Warning(message, sourceName, firstExtraTemplateParam.GetLineno(), firstExtraTemplateParam.GetCharno())
	}
	return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_List](templateArgs.SubList(int32(0), requiredTemplateArgCount)))
}

// Java method: JSTypeRegistry.createRecordTypeFromNodes(com.google.javascript.rhino.Node n, java.lang.String sourceName, com.google.javascript.rhino.jstype.StaticTypedScope scope)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_createRecordTypeFromNodes_9328(n Node, sourceName string, scope StaticTypedScope) JSType {
	_ = n
	_ = sourceName
	_ = scope
	var val_56 RecordTypeBuilder
	InitRecordTypeBuilder(&val_56, this_)
	var builder *RecordTypeBuilder = &val_56
	_ = builder
	{
		var fieldTypeNode Node = n.GetFirstChild()
		_ = fieldTypeNode
		for ; fieldTypeNode != 0; fieldTypeNode = fieldTypeNode.GetNext() {
			var fieldNameNode Node = fieldTypeNode
			_ = fieldNameNode
			var hasType bool = false
			_ = hasType
			if fieldTypeNode.RawToken_() == uint8(163) {
				fieldNameNode = fieldTypeNode.GetFirstChild()
				hasType = true
			}
			var fieldName string = fieldNameNode.GetString()
			_ = fieldName
			if javalang.StringStartsWith(j2gort.CheckNotNullString(fieldName), "'") || javalang.StringStartsWith(j2gort.CheckNotNullString(fieldName), "\"") {
				fieldName = j2gort.JStrSub(fieldName, int32(1), (j2gort.JStrLen(fieldName) - int32(1)))
			}
			var fieldType JSType = nil
			_ = fieldType
			if hasType {
				fieldType = this_.CreateTypeFromCommentNode_3(fieldTypeNode.GetLastChild(), sourceName, scope)
			} else {
				fieldType = this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
			}
			builder.AddProperty(fieldName, fieldType, fieldNameNode)
		}
	}
	return builder.Build()
}

// Java method: JSTypeRegistry.registerTemplateTypeNamesInScope(java.lang.Iterable<com.google.javascript.rhino.jstype.TemplateType> keys, com.google.javascript.rhino.Node scopeRoot)
func (this_ *JSTypeRegistry) RegisterTemplateTypeNamesInScope_2_Ext_java_lang_Iterable_Node(keys j2gort.Ext_java_lang_Iterable, scopeRoot Node) {
	_ = keys
	_ = scopeRoot
	{
		cur_57 := j2gort.NewIterCursor(j2gort.CheckNotNull(keys))
		for cur_57.HasNext() {
			rawElem_58 := cur_57.Next()
			key := j2gort.CastPtr[TemplateType](rawElem_58)
			_ = key
			this_.RegisterSingleTemplateTypeNameInScope(key, scopeRoot)
		}
	}
}

// Java method: JSTypeRegistry.registerTemplateTypeNamesInScope(com.google.javascript.rhino.jstype.TemplateTypeMap map, com.google.javascript.rhino.Node scopeRoot)
func (this_ *JSTypeRegistry) RegisterTemplateTypeNamesInScope_2_TemplateTypeMap_Node(map_ *TemplateTypeMap, scopeRoot Node) {
	_ = map_
	_ = scopeRoot
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = map_.Size()
		_ = n
		for ; i < n; i = i + 1 {
			this_.RegisterSingleTemplateTypeNameInScope(map_.GetTemplateKeyAt(i), scopeRoot)
		}
	}
}

// Java method: JSTypeRegistry.registerSingleTemplateTypeNameInScope(com.google.javascript.rhino.jstype.TemplateType key, com.google.javascript.rhino.Node scopeRoot)
func (this_ *JSTypeRegistry) RegisterSingleTemplateTypeNameInScope(key *TemplateType, scopeRoot Node) {
	_ = key
	_ = scopeRoot
	var refName string = key.GetReferenceName()
	_ = refName
	this_.scopedNameTable.Put(scopeRoot, refName, ToIface_JSType(key))
	if j2gort.JStrIndexOfCp(refName, int32(uint16(46))) >= int32(0) {
		this_.priv_JSTypeRegistry_addQualifiedScopedNameHash_9260(Node_HashQualifiedNameString(refName))
	}
}

// Java method: JSTypeRegistry.createScopeWithTemplates(com.google.javascript.rhino.jstype.StaticTypedScope scope, java.lang.Iterable<com.google.javascript.rhino.jstype.TemplateType> templates)
func (this_ *JSTypeRegistry) CreateScopeWithTemplates_2_StaticTypedScope_Ext_java_lang_Iterable(scope StaticTypedScope, templates j2gort.Ext_java_lang_Iterable) StaticTypedScope {
	_ = scope
	_ = templates
	return NewJSTypeRegistry_SyntheticTemplateScope_2_StaticTypedScope_Ext_java_lang_Iterable(scope, templates)
}

// Java method: JSTypeRegistry.createScopeWithTemplates(com.google.javascript.rhino.jstype.StaticTypedScope scope, com.google.javascript.rhino.jstype.TemplateTypeMap templateMap)
func (this_ *JSTypeRegistry) CreateScopeWithTemplates_2_StaticTypedScope_TemplateTypeMap(scope StaticTypedScope, templateMap *TemplateTypeMap) StaticTypedScope {
	_ = scope
	_ = templateMap
	if templateMap.IsEmpty() {
		return scope
	}
	return NewJSTypeRegistry_SyntheticTemplateScope_2_StaticTypedScope_TemplateTypeMap(scope, templateMap)
}

// Java method: JSTypeRegistry.createScopeWithTemplates(com.google.javascript.rhino.jstype.StaticTypedScope scope, com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.TemplateType> primaryTemplates, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable TemplateTypeMap secondaryTemplateMap)
func (this_ *JSTypeRegistry) CreateScopeWithTemplates_3(scope StaticTypedScope, primaryTemplates j2gort.Ext_com_google_common_collect_ImmutableList, secondaryTemplateMap *TemplateTypeMap) StaticTypedScope {
	_ = scope
	_ = primaryTemplates
	_ = secondaryTemplateMap
	if primaryTemplates.IsEmpty() {
		if (secondaryTemplateMap != nil) && !(secondaryTemplateMap.IsEmpty()) {
			return NewJSTypeRegistry_SyntheticTemplateScope_2_StaticTypedScope_TemplateTypeMap(scope, secondaryTemplateMap)
		} else {
			return scope
		}
	}
	return NewJSTypeRegistry_SyntheticTemplateScope_3(scope, primaryTemplates, secondaryTemplateMap)
}

// Java method: JSTypeRegistry.nativeConstructorBuilder(java.lang.String name)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_nativeConstructorBuilder_9335(fb_ *FunctionType_Builder, name string) {
	_ = name
	InitFunctionType_Builder(fb_, this_)
	fb_.ForNativeType()
	fb_.ForConstructor()
	fb_.WithName(name)
}

// Java method: JSTypeRegistry.nativeInterface(java.lang.String name, com.google.javascript.rhino.jstype.TemplateType[] templateKeys)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_nativeInterface_9336(name string, templateKeys []*TemplateType) *FunctionType {
	_ = name
	_ = templateKeys
	var fb_ FunctionType_Builder
	InitFunctionType_Builder(&fb_, this_)
	(&fb_).ForNativeType()
	(&fb_).ForInterface()
	(&fb_).WithName(name)
	var builder *FunctionType_Builder = &fb_
	_ = builder
	if int32(len(templateKeys)) > int32(0) {
		builder.WithTemplateKeys_1_TemplateType(templateKeys)
	}
	return j2gort.CastPtr[FunctionType](builder.Build())
}

// Java method: JSTypeRegistry.nativeRecord(java.lang.String name, com.google.javascript.rhino.jstype.TemplateType[] templateKeys)
func (this_ *JSTypeRegistry) priv_JSTypeRegistry_nativeRecord_9337(name string, templateKeys []*TemplateType) *FunctionType {
	_ = name
	_ = templateKeys
	var type_ *FunctionType = this_.priv_JSTypeRegistry_nativeInterface_9336(name, templateKeys)
	_ = type_
	type_.SetImplicitMatch(true)
	return type_
}

// Java method: JSTypeRegistry.registerNonLegacyClosureNamespace(java.lang.String moduleName, com.google.javascript.rhino.Node definitionNode, com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeRegistry) RegisterNonLegacyClosureNamespace(moduleName string, definitionNode Node, type_ JSType) {
	_ = moduleName
	_ = definitionNode
	_ = type_
	this_.closureNamespaces.Put(moduleName, JSTypeRegistry_ClosureNamespace_Create(false, definitionNode, type_))
	this_.nonLegacyClosureNamespaceRoots.Add(JSTypeRegistry_GetRootElementOfName(moduleName))
}

// Java method: JSTypeRegistry.registerLegacyClosureNamespace(java.lang.String moduleName)
func (this_ *JSTypeRegistry) RegisterLegacyClosureNamespace(moduleName string) {
	_ = moduleName
	this_.closureNamespaces.Put(moduleName, JSTypeRegistry_ClosureNamespace_Create(true, Node(0), nil))
}

// Java class: com.google.javascript.rhino.jstype.JSTypeRegistry.ScopedNameTable
type JSTypeRegistry_ScopedNameTable struct {
	keyNodes []Node   `gson:"keyNodes"`
	keyNames []string `gson:"keyNames"`
	values   []JSType `gson:"values"`
	size     int32    `gson:"size"`
	idHash_  atomic.Int32
}

func (this_ *JSTypeRegistry_ScopedNameTable) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeRegistry_ScopedNameTable) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeRegistry_ScopedNameTable) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSTypeRegistry$ScopedNameTable", "ScopedNameTable"
}
func (this_ *JSTypeRegistry_ScopedNameTable) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSTypeRegistry_ScopedNameTable) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeRegistry.ScopedNameTable", "com.google.javascript.rhino.jstype.JSTypeRegistry$ScopedNameTable", func(v_ any) bool {
		return j2gort.InstanceOf[*JSTypeRegistry_ScopedNameTable](v_)
	})
}

func (this_ *JSTypeRegistry_ScopedNameTable) initFields_JSTypeRegistry_ScopedNameTable() {
	_ = this_
}

// Java constructor: ScopedNameTable()
func NewJSTypeRegistry_ScopedNameTable() *JSTypeRegistry_ScopedNameTable {
	this_ := &JSTypeRegistry_ScopedNameTable{}
	InitJSTypeRegistry_ScopedNameTable(this_)
	return this_
}

func InitJSTypeRegistry_ScopedNameTable(this_ *JSTypeRegistry_ScopedNameTable) {
	this_.initFields_JSTypeRegistry_ScopedNameTable()
}

// Java static method: ScopedNameTable.mix(int h)
func JSTypeRegistry_ScopedNameTable_Mix(h int32) int32 {
	_ = h
	return (h ^ (javalang.Ushr32(h, int(int32(16)))))
}

// Java method: ScopedNameTable.get(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rootNode, java.lang.@org.jspecify.annotations.Nullable String name)
func (this_ *JSTypeRegistry_ScopedNameTable) Get(rootNode Node, name string) JSType {
	_ = rootNode
	_ = name
	var names []string = this_.keyNames
	_ = names
	if ((names == nil) || (rootNode == 0)) || (name == j2gort.NullString) {
		return nil
	}
	var mask int32 = (int32(len(names)) - int32(1))
	_ = mask
	var slot int32 = (JSTypeRegistry_ScopedNameTable_Mix(((int32(j2gort.CheckNotNullNode(rootNode)) * int32(31)) + javalang.StringHashCode(j2gort.CheckNotNullString(name)))) & mask)
	_ = slot
	for {
		var existingName string = names[slot]
		_ = existingName
		if existingName == j2gort.NullString {
			return nil
		}
		if (this_.keyNodes[slot] == rootNode) && (j2gort.CheckNotNullString(existingName) == name) {
			return this_.values[slot]
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: ScopedNameTable.put(com.google.javascript.rhino.Node rootNode, java.lang.String name, com.google.javascript.rhino.jstype.JSType value)
func (this_ *JSTypeRegistry_ScopedNameTable) Put(rootNode Node, name string, value JSType) {
	_ = rootNode
	_ = name
	_ = value
	j2gort.CheckNotNullNode(rootNode)
	j2gort.CheckNotNullString(name)
	j2gort.CheckNotNullIface(value)
	var names []string = this_.keyNames
	_ = names
	if names == nil {
		names = j2gort.NewStringArray(int32(256))
		this_.keyNames = names
		this_.keyNodes = make([]Node, int32(256))
		this_.values = make([]JSType, int32(256))
	}
	var mask int32 = (int32(len(names)) - int32(1))
	_ = mask
	var slot int32 = (JSTypeRegistry_ScopedNameTable_Mix(((int32(j2gort.CheckNotNullNode(rootNode)) * int32(31)) + javalang.StringHashCode(j2gort.CheckNotNullString(name)))) & mask)
	_ = slot
	for {
		var existingName string = names[slot]
		_ = existingName
		if existingName == j2gort.NullString {
			{
				a_59 := this_.keyNodes
				i_60 := slot
				a_59[i_60] = rootNode
			}
			{
				a_61 := names
				i_62 := slot
				a_61[i_62] = name
			}
			{
				a_63 := this_.values
				i_64 := slot
				a_63[i_64] = value
			}
			this_.size = this_.size + 1
			if (this_.size * int32(4)) >= (int32(len(names)) * int32(3)) {
				this_.priv_ScopedNameTable_rehash_9342((int32(len(names)) * int32(2)))
			}
			return
		}
		if (this_.keyNodes[slot] == rootNode) && (j2gort.CheckNotNullString(existingName) == name) {
			{
				a_65 := this_.values
				i_66 := slot
				a_65[i_66] = value
			}
			return
		}
		slot = ((slot + int32(1)) & mask)
	}
}

// Java method: ScopedNameTable.rehash(int newCap)
func (this_ *JSTypeRegistry_ScopedNameTable) priv_ScopedNameTable_rehash_9342(newCap int32) {
	_ = newCap
	var oldNodes []Node = this_.keyNodes
	_ = oldNodes
	var oldNames []string = this_.keyNames
	_ = oldNames
	var oldValues []JSType = this_.values
	_ = oldValues
	var newNodes []Node = make([]Node, newCap)
	_ = newNodes
	var newNames []string = j2gort.NewStringArray(newCap)
	_ = newNames
	var newValues []JSType = make([]JSType, newCap)
	_ = newValues
	var mask int32 = (newCap - int32(1))
	_ = mask
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(oldNames))
		_ = len_
		for ; i < len_; i = i + 1 {
			var name string = oldNames[i]
			_ = name
			if name != j2gort.NullString {
				var node Node = oldNodes[i]
				_ = node
				var slot int32 = (JSTypeRegistry_ScopedNameTable_Mix(((int32(j2gort.CheckNotNullNode(node)) * int32(31)) + javalang.StringHashCode(j2gort.CheckNotNullString(name)))) & mask)
				_ = slot
				for newNames[slot] != j2gort.NullString {
					slot = ((slot + int32(1)) & mask)
				}
				{
					a_67 := newNodes
					i_68 := slot
					a_67[i_68] = node
				}
				{
					a_69 := newNames
					i_70 := slot
					a_69[i_70] = name
				}
				{
					a_71 := newValues
					i_72 := slot
					a_71[i_72] = oldValues[i]
				}
			}
		}
	}
	this_.keyNodes = newNodes
	this_.keyNames = newNames
	this_.values = newValues
}

// Java class: com.google.javascript.rhino.jstype.JSTypeRegistry.ScopedNameSet
type JSTypeRegistry_ScopedNameSet struct {
	keyNodes []Node   `gson:"keyNodes"`
	keyNames []string `gson:"keyNames"`
	size     int32    `gson:"size"`
	idHash_  atomic.Int32
}

func (this_ *JSTypeRegistry_ScopedNameSet) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeRegistry_ScopedNameSet) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeRegistry_ScopedNameSet) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSTypeRegistry$ScopedNameSet", "ScopedNameSet"
}
func (this_ *JSTypeRegistry_ScopedNameSet) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSTypeRegistry_ScopedNameSet) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeRegistry.ScopedNameSet", "com.google.javascript.rhino.jstype.JSTypeRegistry$ScopedNameSet", func(v_ any) bool {
		return j2gort.InstanceOf[*JSTypeRegistry_ScopedNameSet](v_)
	})
}

func (this_ *JSTypeRegistry_ScopedNameSet) initFields_JSTypeRegistry_ScopedNameSet() {
	_ = this_
}

// Java constructor: ScopedNameSet()
func NewJSTypeRegistry_ScopedNameSet() *JSTypeRegistry_ScopedNameSet {
	this_ := &JSTypeRegistry_ScopedNameSet{}
	InitJSTypeRegistry_ScopedNameSet(this_)
	return this_
}

func InitJSTypeRegistry_ScopedNameSet(this_ *JSTypeRegistry_ScopedNameSet) {
	this_.initFields_JSTypeRegistry_ScopedNameSet()
}

// Java static method: ScopedNameSet.mix(int h)
func JSTypeRegistry_ScopedNameSet_Mix(h int32) int32 {
	_ = h
	return (h ^ (javalang.Ushr32(h, int(int32(16)))))
}

// Java method: ScopedNameSet.containsEntry(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rootNode, java.lang.@org.jspecify.annotations.Nullable String name)
func (this_ *JSTypeRegistry_ScopedNameSet) ContainsEntry(rootNode Node, name string) bool {
	_ = rootNode
	_ = name
	var names []string = this_.keyNames
	_ = names
	if ((names == nil) || (rootNode == 0)) || (name == j2gort.NullString) {
		return false
	}
	var mask int32 = (int32(len(names)) - int32(1))
	_ = mask
	var slot int32 = (JSTypeRegistry_ScopedNameSet_Mix(((int32(j2gort.CheckNotNullNode(rootNode)) * int32(31)) + javalang.StringHashCode(j2gort.CheckNotNullString(name)))) & mask)
	_ = slot
	for {
		var existingName string = names[slot]
		_ = existingName
		if existingName == j2gort.NullString {
			return false
		}
		if (this_.keyNodes[slot] == rootNode) && (j2gort.CheckNotNullString(existingName) == name) {
			return true
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: ScopedNameSet.put(com.google.javascript.rhino.Node rootNode, java.lang.String name)
func (this_ *JSTypeRegistry_ScopedNameSet) Put(rootNode Node, name string) {
	_ = rootNode
	_ = name
	j2gort.CheckNotNullNode(rootNode)
	j2gort.CheckNotNullString(name)
	var names []string = this_.keyNames
	_ = names
	if names == nil {
		names = j2gort.NewStringArray(int32(64))
		this_.keyNames = names
		this_.keyNodes = make([]Node, int32(64))
	}
	var mask int32 = (int32(len(names)) - int32(1))
	_ = mask
	var slot int32 = (JSTypeRegistry_ScopedNameSet_Mix(((int32(j2gort.CheckNotNullNode(rootNode)) * int32(31)) + javalang.StringHashCode(j2gort.CheckNotNullString(name)))) & mask)
	_ = slot
	for {
		var existingName string = names[slot]
		_ = existingName
		if existingName == j2gort.NullString {
			{
				a_73 := this_.keyNodes
				i_74 := slot
				a_73[i_74] = rootNode
			}
			{
				a_75 := names
				i_76 := slot
				a_75[i_76] = name
			}
			this_.size = this_.size + 1
			if (this_.size * int32(4)) >= (int32(len(names)) * int32(3)) {
				this_.priv_ScopedNameSet_rehash_9345((int32(len(names)) * int32(2)))
			}
			return
		}
		if (this_.keyNodes[slot] == rootNode) && (j2gort.CheckNotNullString(existingName) == name) {
			return
		}
		slot = ((slot + int32(1)) & mask)
	}
}

// Java method: ScopedNameSet.rehash(int newCap)
func (this_ *JSTypeRegistry_ScopedNameSet) priv_ScopedNameSet_rehash_9345(newCap int32) {
	_ = newCap
	var oldNodes []Node = this_.keyNodes
	_ = oldNodes
	var oldNames []string = this_.keyNames
	_ = oldNames
	var newNodes []Node = make([]Node, newCap)
	_ = newNodes
	var newNames []string = j2gort.NewStringArray(newCap)
	_ = newNames
	var mask int32 = (newCap - int32(1))
	_ = mask
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(oldNames))
		_ = len_
		for ; i < len_; i = i + 1 {
			var name string = oldNames[i]
			_ = name
			if name != j2gort.NullString {
				var node Node = oldNodes[i]
				_ = node
				var slot int32 = (JSTypeRegistry_ScopedNameSet_Mix(((int32(j2gort.CheckNotNullNode(node)) * int32(31)) + javalang.StringHashCode(j2gort.CheckNotNullString(name)))) & mask)
				_ = slot
				for newNames[slot] != j2gort.NullString {
					slot = ((slot + int32(1)) & mask)
				}
				{
					a_77 := newNodes
					i_78 := slot
					a_77[i_78] = node
				}
				{
					a_79 := newNames
					i_80 := slot
					a_79[i_80] = name
				}
			}
		}
	}
	this_.keyNodes = newNodes
	this_.keyNames = newNames
}

// Java enum: com.google.javascript.rhino.jstype.JSTypeRegistry.PropDefinitionKind
type JSTypeRegistry_PropDefinitionKind uint8

const (
	enum_JSTypeRegistry_PropDefinitionKind_UNKNOWN     JSTypeRegistry_PropDefinitionKind = 1
	enum_JSTypeRegistry_PropDefinitionKind_KNOWN       JSTypeRegistry_PropDefinitionKind = 2
	enum_JSTypeRegistry_PropDefinitionKind_LOOSE       JSTypeRegistry_PropDefinitionKind = 3
	enum_JSTypeRegistry_PropDefinitionKind_LOOSE_UNION JSTypeRegistry_PropDefinitionKind = 4
)

var enumBaseTable_JSTypeRegistry_PropDefinitionKind = [...]j2gort.EnumBase{
	{Name_: "UNKNOWN", Ordinal_: 0},
	{Name_: "KNOWN", Ordinal_: 1},
	{Name_: "LOOSE", Ordinal_: 2},
	{Name_: "LOOSE_UNION", Ordinal_: 3},
}
var enumValues_JSTypeRegistry_PropDefinitionKind = []JSTypeRegistry_PropDefinitionKind{1, 2, 3, 4}

func StaticGet_JSTypeRegistry_PropDefinitionKind_UNKNOWN() JSTypeRegistry_PropDefinitionKind {
	return enum_JSTypeRegistry_PropDefinitionKind_UNKNOWN
}
func StaticGet_JSTypeRegistry_PropDefinitionKind_KNOWN() JSTypeRegistry_PropDefinitionKind {
	return enum_JSTypeRegistry_PropDefinitionKind_KNOWN
}
func StaticGet_JSTypeRegistry_PropDefinitionKind_LOOSE() JSTypeRegistry_PropDefinitionKind {
	return enum_JSTypeRegistry_PropDefinitionKind_LOOSE
}
func StaticGet_JSTypeRegistry_PropDefinitionKind_LOOSE_UNION() JSTypeRegistry_PropDefinitionKind {
	return enum_JSTypeRegistry_PropDefinitionKind_LOOSE_UNION
}

func JSTypeRegistry_PropDefinitionKind_Values() []JSTypeRegistry_PropDefinitionKind {
	return j2gort.CloneSlice(enumValues_JSTypeRegistry_PropDefinitionKind)
}

func JSTypeRegistry_PropDefinitionKind_ValueOf(name string) JSTypeRegistry_PropDefinitionKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSTypeRegistry_PropDefinitionKind {
		if enumBaseTable_JSTypeRegistry_PropDefinitionKind[i_].Name_ == name {
			return JSTypeRegistry_PropDefinitionKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.JSTypeRegistry.PropDefinitionKind." + name))
}

func (this_ JSTypeRegistry_PropDefinitionKind) IsJavaNil_() bool                              { return this_ == 0 }
func (JSTypeRegistry_PropDefinitionKind) HasIdentityEquals_()                                 {}
func (JSTypeRegistry_PropDefinitionKind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JSTypeRegistry_PropDefinitionKind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSTypeRegistry_PropDefinitionKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JSTypeRegistry_PropDefinitionKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSTypeRegistry$PropDefinitionKind", "PropDefinitionKind"
}
func (this_ JSTypeRegistry_PropDefinitionKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSTypeRegistry_PropDefinitionKind[this_-1]
}
func (this_ JSTypeRegistry_PropDefinitionKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSTypeRegistry_PropDefinitionKind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSTypeRegistry_PropDefinitionKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSTypeRegistry_PropDefinitionKind[this_-1].Name_
}
func (this_ JSTypeRegistry_PropDefinitionKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSTypeRegistry_PropDefinitionKind[this_-1].Name_
}
func (this_ JSTypeRegistry_PropDefinitionKind) Equals(other_ any) bool {
	o_, ok_ := other_.(JSTypeRegistry_PropDefinitionKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSTypeRegistry_PropDefinitionKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSTypeRegistry_PropDefinitionKind[this_-1].HashCode()
}
func (this_ JSTypeRegistry_PropDefinitionKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSTypeRegistry_PropDefinitionKind[this_-1].CompareTo(other_)
}
func (this_ JSTypeRegistry_PropDefinitionKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSTypeRegistry_PropDefinitionKind[this_-1].Name_)
}
func init() {
	vals_ := []any{JSTypeRegistry_PropDefinitionKind(1), JSTypeRegistry_PropDefinitionKind(2), JSTypeRegistry_PropDefinitionKind(3), JSTypeRegistry_PropDefinitionKind(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSTypeRegistry.PropDefinitionKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSTypeRegistry$PropDefinitionKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeRegistry.PropDefinitionKind", "com.google.javascript.rhino.jstype.JSTypeRegistry$PropDefinitionKind", func(v_ any) bool {
		return j2gort.InstanceOf[JSTypeRegistry_PropDefinitionKind](v_)
	})
}

// Java class: com.google.javascript.rhino.jstype.JSTypeRegistry.SyntheticTemplateScope
type JSTypeRegistry_SyntheticTemplateScope struct {
	delegate   StaticTypedScope                 `gson:"delegate"`
	singleKey  string                           `gson:"singleKey"`
	singleType *TemplateType                    `gson:"singleType"`
	types      *HamtPMap[string, *TemplateType] `gson:"types"`
	idHash_    atomic.Int32
}

func (this_ *JSTypeRegistry_SyntheticTemplateScope) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeRegistry_SyntheticTemplateScope) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeRegistry_SyntheticTemplateScope) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSTypeRegistry$SyntheticTemplateScope", "SyntheticTemplateScope"
}
func (this_ *JSTypeRegistry_SyntheticTemplateScope) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSTypeRegistry_SyntheticTemplateScope) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeRegistry.SyntheticTemplateScope", "com.google.javascript.rhino.jstype.JSTypeRegistry$SyntheticTemplateScope", func(v_ any) bool {
		return j2gort.InstanceOf[*JSTypeRegistry_SyntheticTemplateScope](v_)
	})
}
func (*JSTypeRegistry_SyntheticTemplateScope) Implements_rhino_StaticTypedScope() {}
func (*JSTypeRegistry_SyntheticTemplateScope) Implements_rhino_StaticScope()      {}

var static_JSTypeRegistry_SyntheticTemplateScope_EMPTY_TEMPLATE_TYPES *HamtPMap[string, *TemplateType]

func StaticGet_JSTypeRegistry_SyntheticTemplateScope_EMPTY_TEMPLATE_TYPES() *HamtPMap[string, *TemplateType] {
	clinit_JSTypeRegistry_SyntheticTemplateScope()
	return static_JSTypeRegistry_SyntheticTemplateScope_EMPTY_TEMPLATE_TYPES
}
func StaticSet_JSTypeRegistry_SyntheticTemplateScope_EMPTY_TEMPLATE_TYPES(v *HamtPMap[string, *TemplateType]) {
	clinit_JSTypeRegistry_SyntheticTemplateScope()
	static_JSTypeRegistry_SyntheticTemplateScope_EMPTY_TEMPLATE_TYPES = v
}

var clinitDone_JSTypeRegistry_SyntheticTemplateScope uint32
var clinitState_JSTypeRegistry_SyntheticTemplateScope j2gort.ClinitState

func clinit_JSTypeRegistry_SyntheticTemplateScope() {
	if atomic.LoadUint32(&clinitDone_JSTypeRegistry_SyntheticTemplateScope) != 0 {
		return
	}
	clinitSlow_JSTypeRegistry_SyntheticTemplateScope()
}

func clinitSlow_JSTypeRegistry_SyntheticTemplateScope() {
	if !j2gort.EnterClinit(&clinitState_JSTypeRegistry_SyntheticTemplateScope, &clinitDone_JSTypeRegistry_SyntheticTemplateScope) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSTypeRegistry_SyntheticTemplateScope, &clinitDone_JSTypeRegistry_SyntheticTemplateScope)
	static_JSTypeRegistry_SyntheticTemplateScope_EMPTY_TEMPLATE_TYPES = HamtPMap_Empty[string, *TemplateType]()
}

func (this_ *JSTypeRegistry_SyntheticTemplateScope) initFields_JSTypeRegistry_SyntheticTemplateScope() {
	_ = this_
	this_.singleKey = j2gort.NullString
}

// Java constructor: SyntheticTemplateScope(com.google.javascript.rhino.jstype.StaticTypedScope delegate, java.lang.Iterable<com.google.javascript.rhino.jstype.TemplateType> templates)
func NewJSTypeRegistry_SyntheticTemplateScope_2_StaticTypedScope_Ext_java_lang_Iterable(delegate StaticTypedScope, templates j2gort.Ext_java_lang_Iterable) *JSTypeRegistry_SyntheticTemplateScope {
	clinit_JSTypeRegistry_SyntheticTemplateScope()
	this_ := &JSTypeRegistry_SyntheticTemplateScope{}
	InitJSTypeRegistry_SyntheticTemplateScope_2_StaticTypedScope_Ext_java_lang_Iterable(this_, delegate, templates)
	return this_
}

func InitJSTypeRegistry_SyntheticTemplateScope_2_StaticTypedScope_Ext_java_lang_Iterable(this_ *JSTypeRegistry_SyntheticTemplateScope, delegate StaticTypedScope, templates j2gort.Ext_java_lang_Iterable) {
	_ = delegate
	_ = templates
	this_.initFields_JSTypeRegistry_SyntheticTemplateScope()
	this_.delegate = delegate
	var syntheticTemplateScope *JSTypeRegistry_SyntheticTemplateScope
	_ = syntheticTemplateScope
	var list j2gort.Ext_com_google_common_collect_ImmutableList
	_ = list
	var ok_81 bool
	syntheticTemplateScope, ok_81 = j2gort.TryCast[*JSTypeRegistry_SyntheticTemplateScope](delegate)
	if ok_81 {
		var map_ *HamtPMap[string, *TemplateType] = syntheticTemplateScope.priv_SyntheticTemplateScope_asMap_9346()
		_ = map_
		{
			cur_82 := j2gort.NewIterCursor(j2gort.CheckNotNull(templates))
			for cur_82.HasNext() {
				rawElem_83 := cur_82.Next()
				key := j2gort.CastPtr[TemplateType](rawElem_83)
				_ = key
				map_ = j2gort.CastPtr[HamtPMap[string, *TemplateType]](map_.Plus_2(key.GetReferenceName(), key))
			}
		}
		this_.singleKey = j2gort.NullString
		this_.singleType = nil
		this_.types = map_
	} else {
		var ok_84 bool
		list, ok_84 = j2gort.TryCast[j2gort.Ext_com_google_common_collect_ImmutableList](templates)
		if ok_84 && (list.Size() == int32(1)) {
			var t0 *TemplateType = j2gort.CastPtr[TemplateType](list.Get(int32(0)))
			_ = t0
			this_.singleKey = t0.GetReferenceName()
			this_.singleType = t0
			this_.types = nil
		} else {
			var map__2 *HamtPMap[string, *TemplateType] = StaticGet_JSTypeRegistry_SyntheticTemplateScope_EMPTY_TEMPLATE_TYPES()
			_ = map__2
			{
				cur_85 := j2gort.NewIterCursor(j2gort.CheckNotNull(templates))
				for cur_85.HasNext() {
					rawElem_86 := cur_85.Next()
					key_2 := j2gort.CastPtr[TemplateType](rawElem_86)
					_ = key_2
					map__2 = j2gort.CastPtr[HamtPMap[string, *TemplateType]](map__2.Plus_2(key_2.GetReferenceName(), key_2))
				}
			}
			this_.singleKey = j2gort.NullString
			this_.singleType = nil
			this_.types = map__2
		}
	}
}

// Java constructor: SyntheticTemplateScope(com.google.javascript.rhino.jstype.StaticTypedScope delegate, com.google.javascript.rhino.jstype.TemplateTypeMap templateMap)
func NewJSTypeRegistry_SyntheticTemplateScope_2_StaticTypedScope_TemplateTypeMap(delegate StaticTypedScope, templateMap *TemplateTypeMap) *JSTypeRegistry_SyntheticTemplateScope {
	clinit_JSTypeRegistry_SyntheticTemplateScope()
	this_ := &JSTypeRegistry_SyntheticTemplateScope{}
	InitJSTypeRegistry_SyntheticTemplateScope_2_StaticTypedScope_TemplateTypeMap(this_, delegate, templateMap)
	return this_
}

func InitJSTypeRegistry_SyntheticTemplateScope_2_StaticTypedScope_TemplateTypeMap(this_ *JSTypeRegistry_SyntheticTemplateScope, delegate StaticTypedScope, templateMap *TemplateTypeMap) {
	_ = delegate
	_ = templateMap
	this_.initFields_JSTypeRegistry_SyntheticTemplateScope()
	this_.delegate = delegate
	var n int32 = templateMap.Size()
	_ = n
	if !(j2gort.InstanceOf[*JSTypeRegistry_SyntheticTemplateScope](delegate)) && (n == int32(1)) {
		var t0 *TemplateType = templateMap.GetTemplateKeyAt(int32(0))
		_ = t0
		this_.singleKey = t0.GetReferenceName()
		this_.singleType = t0
		this_.types = nil
	} else {
		var syntheticTemplateScope *JSTypeRegistry_SyntheticTemplateScope
		_ = syntheticTemplateScope
		var map_ *HamtPMap[string, *TemplateType]
		_ = map_
		var ok_87 bool
		syntheticTemplateScope, ok_87 = j2gort.TryCast[*JSTypeRegistry_SyntheticTemplateScope](delegate)
		if ok_87 {
			map_ = syntheticTemplateScope.priv_SyntheticTemplateScope_asMap_9346()
		} else {
			map_ = StaticGet_JSTypeRegistry_SyntheticTemplateScope_EMPTY_TEMPLATE_TYPES()
		}
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				var key *TemplateType = templateMap.GetTemplateKeyAt(i)
				_ = key
				map_ = j2gort.CastPtr[HamtPMap[string, *TemplateType]](map_.Plus_2(key.GetReferenceName(), key))
			}
		}
		this_.singleKey = j2gort.NullString
		this_.singleType = nil
		this_.types = map_
	}
}

// Java constructor: SyntheticTemplateScope(com.google.javascript.rhino.jstype.StaticTypedScope delegate, com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.TemplateType> primaryTemplates, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable TemplateTypeMap secondaryTemplateMap)
func NewJSTypeRegistry_SyntheticTemplateScope_3(delegate StaticTypedScope, primaryTemplates j2gort.Ext_com_google_common_collect_ImmutableList, secondaryTemplateMap *TemplateTypeMap) *JSTypeRegistry_SyntheticTemplateScope {
	clinit_JSTypeRegistry_SyntheticTemplateScope()
	this_ := &JSTypeRegistry_SyntheticTemplateScope{}
	InitJSTypeRegistry_SyntheticTemplateScope_3(this_, delegate, primaryTemplates, secondaryTemplateMap)
	return this_
}

func InitJSTypeRegistry_SyntheticTemplateScope_3(this_ *JSTypeRegistry_SyntheticTemplateScope, delegate StaticTypedScope, primaryTemplates j2gort.Ext_com_google_common_collect_ImmutableList, secondaryTemplateMap *TemplateTypeMap) {
	_ = delegate
	_ = primaryTemplates
	_ = secondaryTemplateMap
	this_.initFields_JSTypeRegistry_SyntheticTemplateScope()
	this_.delegate = delegate
	var n1 int32 = primaryTemplates.Size()
	_ = n1
	var n2 int32
	_ = n2
	if secondaryTemplateMap != nil {
		n2 = secondaryTemplateMap.Size()
	} else {
		n2 = int32(0)
	}
	if (!(j2gort.InstanceOf[*JSTypeRegistry_SyntheticTemplateScope](delegate)) && (n1 == int32(1))) && (n2 == int32(0)) {
		var t0 *TemplateType = j2gort.CastPtr[TemplateType](primaryTemplates.Get(int32(0)))
		_ = t0
		this_.singleKey = t0.GetReferenceName()
		this_.singleType = t0
		this_.types = nil
	} else {
		var syntheticTemplateScope *JSTypeRegistry_SyntheticTemplateScope
		_ = syntheticTemplateScope
		var map_ *HamtPMap[string, *TemplateType]
		_ = map_
		var ok_88 bool
		syntheticTemplateScope, ok_88 = j2gort.TryCast[*JSTypeRegistry_SyntheticTemplateScope](delegate)
		if ok_88 {
			map_ = syntheticTemplateScope.priv_SyntheticTemplateScope_asMap_9346()
		} else {
			map_ = StaticGet_JSTypeRegistry_SyntheticTemplateScope_EMPTY_TEMPLATE_TYPES()
		}
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n1; i = i + 1 {
				var key *TemplateType = j2gort.CastPtr[TemplateType](primaryTemplates.Get(i))
				_ = key
				map_ = j2gort.CastPtr[HamtPMap[string, *TemplateType]](map_.Plus_2(key.GetReferenceName(), key))
			}
		}
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < n2; i_2 = i_2 + 1 {
				var key_2 *TemplateType = secondaryTemplateMap.GetTemplateKeyAt(i_2)
				_ = key_2
				map_ = j2gort.CastPtr[HamtPMap[string, *TemplateType]](map_.Plus_2(key_2.GetReferenceName(), key_2))
			}
		}
		this_.singleKey = j2gort.NullString
		this_.singleType = nil
		this_.types = map_
	}
}

// Java method: SyntheticTemplateScope.asMap()
func (this_ *JSTypeRegistry_SyntheticTemplateScope) priv_SyntheticTemplateScope_asMap_9346() *HamtPMap[string, *TemplateType] {
	if this_.types != nil {
		return this_.types
	}
	if this_.singleKey != j2gort.NullString {
		return HamtPMap_Singleton[string, *TemplateType](this_.singleKey, this_.singleType)
	} else {
		return StaticGet_JSTypeRegistry_SyntheticTemplateScope_EMPTY_TEMPLATE_TYPES()
	}
}

// Java method: SyntheticTemplateScope.getRootNode()
func (this_ *JSTypeRegistry_SyntheticTemplateScope) GetRootNode_0() Node {
	return this_.delegate.GetRootNode_0()
}

// Java method: SyntheticTemplateScope.getParentScope()
func (this_ *JSTypeRegistry_SyntheticTemplateScope) GetParentScope() any {
	return j2gort.AsIface[StaticTypedScope](this_.delegate.GetParentScope())
}

// Java method: SyntheticTemplateScope.getSlot(java.lang.String name)
func (this_ *JSTypeRegistry_SyntheticTemplateScope) GetSlot(name string) any {
	_ = name
	return j2gort.AsIface[StaticTypedSlot](this_.delegate.GetSlot(name))
}

// Java method: SyntheticTemplateScope.getOwnSlot(java.lang.String name)
func (this_ *JSTypeRegistry_SyntheticTemplateScope) GetOwnSlot(name string) any {
	_ = name
	return j2gort.AsIface[StaticTypedSlot](this_.delegate.GetOwnSlot(name))
}

// Java method: SyntheticTemplateScope.getTypeOfThis()
func (this_ *JSTypeRegistry_SyntheticTemplateScope) GetTypeOfThis() JSType {
	return this_.delegate.GetTypeOfThis()
}

// Java method: SyntheticTemplateScope.getTemplateType(java.lang.String name)
func (this_ *JSTypeRegistry_SyntheticTemplateScope) GetTemplateType(name string) *TemplateType {
	_ = name
	if this_.types != nil {
		return this_.types.Get_1(name)
	}
	if (this_.singleKey != j2gort.NullString) && (j2gort.CheckNotNullString(this_.singleKey) == name) {
		return this_.singleType
	} else {
		return nil
	}
}

// Java method: SyntheticTemplateScope.getTopmostScopeOfEventualDeclaration(java.lang.String name)
func (this_ *JSTypeRegistry_SyntheticTemplateScope) GetTopmostScopeOfEventualDeclaration(name string) StaticScope {
	_ = name
	if this_.GetTemplateType(name) != nil {
		return this_
	}
	return this_.delegate.GetTopmostScopeOfEventualDeclaration(name)
}

func (this_ *JSTypeRegistry_SyntheticTemplateScope) LookupQualifiedName_1_QualifiedName(p0_ QualifiedName) JSType {
	return IfaceDefault_StaticTypedScope_LookupQualifiedName_1_QualifiedName(this_, p0_)
}

func (this_ *JSTypeRegistry_SyntheticTemplateScope) LookupQualifiedName_1_string(p0_ string) JSType {
	return IfaceDefault_StaticTypedScope_LookupQualifiedName_1_string(this_, p0_)
}

func (this_ *JSTypeRegistry_SyntheticTemplateScope) LookupQualifiedNameNode(p0_ Node) JSType {
	return IfaceDefault_StaticTypedScope_LookupQualifiedNameNode(this_, p0_)
}

// Java polymorphic class interface: com.google.javascript.rhino.jstype.JSTypeRegistry.ClosureNamespace
type JSTypeRegistry_ClosureNamespace interface {
	Implements_rhino_JSTypeRegistry_ClosureNamespace()
	IsJavaNil_() bool
	GetJSTypeRegistry_ClosureNamespaceBase() *JSTypeRegistry_ClosureNamespaceBase
	IsLegacy() bool
	DefinitionNode() Node
	Type() JSType
}

func ToIface_JSTypeRegistry_ClosureNamespace[T interface {
	JSTypeRegistry_ClosureNamespace
	comparable
}](v T) JSTypeRegistry_ClosureNamespace {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JSTypeRegistry_ClosureNamespace[T JSTypeRegistry_ClosureNamespace](v T) JSTypeRegistry_ClosureNamespace {
	var vi JSTypeRegistry_ClosureNamespace = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.rhino.jstype.JSTypeRegistry.ClosureNamespace
type JSTypeRegistry_ClosureNamespaceBase struct {
	idHash_ atomic.Int32
}

func (this_ *JSTypeRegistry_ClosureNamespaceBase) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeRegistry_ClosureNamespaceBase) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeRegistry_ClosureNamespaceBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSTypeRegistry$ClosureNamespace", "ClosureNamespace"
}
func (this_ *JSTypeRegistry_ClosureNamespaceBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeRegistry.ClosureNamespace", "com.google.javascript.rhino.jstype.JSTypeRegistry$ClosureNamespace", func(v_ any) bool {
		return j2gort.InstanceOf[JSTypeRegistry_ClosureNamespace](v_)
	})
}
func (b *JSTypeRegistry_ClosureNamespaceBase) GetJSTypeRegistry_ClosureNamespaceBase() *JSTypeRegistry_ClosureNamespaceBase {
	return b
}
func (*JSTypeRegistry_ClosureNamespaceBase) Implements_rhino_JSTypeRegistry_ClosureNamespace() {}

func (this_ *JSTypeRegistry_ClosureNamespaceBase) initFields_JSTypeRegistry_ClosureNamespace() {
	_ = this_
}

func InitJSTypeRegistry_ClosureNamespace(this_ *JSTypeRegistry_ClosureNamespaceBase) {
	this_.initFields_JSTypeRegistry_ClosureNamespace()
}

// Java method: ClosureNamespace.isLegacy()
func (this_ *JSTypeRegistry_ClosureNamespaceBase) IsLegacy() bool {
	panic("abstract or native method invoked")
}

// Java method: ClosureNamespace.definitionNode()
func (this_ *JSTypeRegistry_ClosureNamespaceBase) DefinitionNode() Node {
	panic("abstract or native method invoked")
}

// Java method: ClosureNamespace.type()
func (this_ *JSTypeRegistry_ClosureNamespaceBase) Type() JSType {
	panic("abstract or native method invoked")
}

// Java static method: ClosureNamespace.create(boolean isLegacy, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node definitionNode, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type)
func JSTypeRegistry_ClosureNamespace_Create(isLegacy bool, definitionNode Node, type_ JSType) JSTypeRegistry_ClosureNamespace {
	_ = isLegacy
	_ = definitionNode
	_ = type_
	return NewAutoValue_JSTypeRegistry_ClosureNamespace(isLegacy, definitionNode, type_)
}
func (this_ *JSTypeRegistry) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_UnionType.ReleaseWith(func(v_ *UnionType) { *v_ = UnionType{} })
	this_.slab_FunctionType.ReleaseWith(func(v_ *FunctionType) { *v_ = FunctionType{} })
	this_.slab_TemplatizedType.ReleaseWith(func(v_ *TemplatizedType) { *v_ = TemplatizedType{} })
}
