/*
 * Copyright 2004 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/TypedScopeCreator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.TypedScopeCreator
type TypedScopeCreator struct {
	compiler                     AbstractCompiler                               `gson:"compiler"`
	typeParsingErrorReporter     pkg_rhino.ErrorReporter                        `gson:"typeParsingErrorReporter"`
	validator                    *TypeValidator                                 `gson:"validator"`
	codingConvention             CodingConvention                               `gson:"codingConvention"`
	typeRegistry                 *pkg_rhino.JSTypeRegistry                      `gson:"typeRegistry"`
	moduleMap                    *ModuleMap                                     `gson:"moduleMap"`
	metadataMap                  *ModuleMetadataMap                             `gson:"metadataMap"`
	moduleImportResolver         *ModuleImportResolver                          `gson:"moduleImportResolver"`
	memoized                     j2gort.OrderedMap[pkg_rhino.Node, *TypedScope] `gson:"memoized"`
	reservedNamesForScope        j2gort.PlainMap[pkg_rhino.Node, []string]      `gson:"reservedNamesForScope"`
	functionsWithNonEmptyReturns j2gort.PlainSet[pkg_rhino.Node]                `gson:"functionsWithNonEmptyReturns"`
	fnTypeBuilderPool            []*FunctionTypeBuilder                         `gson:"fnTypeBuilderPool"`
	overflowPool                 []*TypedScope_Overflow                         `gson:"overflowPool"`
	scopedVarStates              *TypedScopeCreator_ScopedVarStateTable         `gson:"scopedVarStates"`
	unknownType                  pkg_rhino.ObjectType                           `gson:"unknownType"`
	weakImports                  []*TypedScopeCreator_WeakModuleImport          `gson:"weakImports"`
	unresolvedNodes              []pkg_rhino.Node                               `gson:"unresolvedNodes"`
	undeclaredNamesForClosure    j2gort.OrderedSet[pkg_rhino.Node]              `gson:"undeclaredNamesForClosure"`
	providedNamesFromCall        j2gort.Ext_com_google_common_collect_Multimap  `gson:"providedNamesFromCall"`
	fnTypeBuilderPoolTop         int32                                          `gson:"fnTypeBuilderPoolTop"`
	overflowPoolTop              int32                                          `gson:"overflowPoolTop"`
	idHash_                      atomic.Int32
	processClosurePrimitives     bool                    `gson:"processClosurePrimitives"`
	stage                        TypedScopeCreator_Stage `gson:"stage"`
}

func (this_ *TypedScopeCreator) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedScopeCreator) IsJavaNil_() bool { return this_ == nil }
func (*TypedScopeCreator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator", "TypedScopeCreator"
}
func (this_ *TypedScopeCreator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedScopeCreator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator", "com.google.javascript.jscomp.TypedScopeCreator", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScopeCreator](v_)
	})
}
func (*TypedScopeCreator) Implements_jscomp_ScopeCreator()     {}
func (*TypedScopeCreator) Implements_rhino_StaticSymbolTable() {}

var static_TypedScopeCreator_MALFORMED_TYPEDEF *DiagnosticType

func StaticGet_TypedScopeCreator_MALFORMED_TYPEDEF() *DiagnosticType {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_MALFORMED_TYPEDEF
}
func StaticSet_TypedScopeCreator_MALFORMED_TYPEDEF(v *DiagnosticType) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_MALFORMED_TYPEDEF = v
}

var static_TypedScopeCreator_ENUM_INITIALIZER *DiagnosticType

func StaticGet_TypedScopeCreator_ENUM_INITIALIZER() *DiagnosticType {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_ENUM_INITIALIZER
}
func StaticSet_TypedScopeCreator_ENUM_INITIALIZER(v *DiagnosticType) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_ENUM_INITIALIZER = v
}

var static_TypedScopeCreator_INVALID_ENUM_KEY *DiagnosticType

func StaticGet_TypedScopeCreator_INVALID_ENUM_KEY() *DiagnosticType {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_INVALID_ENUM_KEY
}
func StaticSet_TypedScopeCreator_INVALID_ENUM_KEY(v *DiagnosticType) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_INVALID_ENUM_KEY = v
}

var static_TypedScopeCreator_CTOR_INITIALIZER *DiagnosticType

func StaticGet_TypedScopeCreator_CTOR_INITIALIZER() *DiagnosticType {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_CTOR_INITIALIZER
}
func StaticSet_TypedScopeCreator_CTOR_INITIALIZER(v *DiagnosticType) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_CTOR_INITIALIZER = v
}

var static_TypedScopeCreator_IFACE_INITIALIZER *DiagnosticType

func StaticGet_TypedScopeCreator_IFACE_INITIALIZER() *DiagnosticType {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_IFACE_INITIALIZER
}
func StaticSet_TypedScopeCreator_IFACE_INITIALIZER(v *DiagnosticType) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_IFACE_INITIALIZER = v
}

var static_TypedScopeCreator_CONSTRUCTOR_EXPECTED *DiagnosticType

func StaticGet_TypedScopeCreator_CONSTRUCTOR_EXPECTED() *DiagnosticType {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_CONSTRUCTOR_EXPECTED
}
func StaticSet_TypedScopeCreator_CONSTRUCTOR_EXPECTED(v *DiagnosticType) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_CONSTRUCTOR_EXPECTED = v
}

var static_TypedScopeCreator_UNKNOWN_LENDS *DiagnosticType

func StaticGet_TypedScopeCreator_UNKNOWN_LENDS() *DiagnosticType {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_UNKNOWN_LENDS
}
func StaticSet_TypedScopeCreator_UNKNOWN_LENDS(v *DiagnosticType) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_UNKNOWN_LENDS = v
}

var static_TypedScopeCreator_LENDS_ON_NON_OBJECT *DiagnosticType

func StaticGet_TypedScopeCreator_LENDS_ON_NON_OBJECT() *DiagnosticType {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_LENDS_ON_NON_OBJECT
}
func StaticSet_TypedScopeCreator_LENDS_ON_NON_OBJECT(v *DiagnosticType) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_LENDS_ON_NON_OBJECT = v
}

var static_TypedScopeCreator_INCOMPATIBLE_ALIAS_ANNOTATION *DiagnosticType

func StaticGet_TypedScopeCreator_INCOMPATIBLE_ALIAS_ANNOTATION() *DiagnosticType {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_INCOMPATIBLE_ALIAS_ANNOTATION
}
func StaticSet_TypedScopeCreator_INCOMPATIBLE_ALIAS_ANNOTATION(v *DiagnosticType) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_INCOMPATIBLE_ALIAS_ANNOTATION = v
}

var static_TypedScopeCreator_DYNAMIC_EXTENDS_WITHOUT_JSDOC *DiagnosticType

func StaticGet_TypedScopeCreator_DYNAMIC_EXTENDS_WITHOUT_JSDOC() *DiagnosticType {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_DYNAMIC_EXTENDS_WITHOUT_JSDOC
}
func StaticSet_TypedScopeCreator_DYNAMIC_EXTENDS_WITHOUT_JSDOC(v *DiagnosticType) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_DYNAMIC_EXTENDS_WITHOUT_JSDOC = v
}

var static_TypedScopeCreator_EXPORTS_ONLY []string

func StaticGet_TypedScopeCreator_EXPORTS_ONLY() []string {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_EXPORTS_ONLY
}
func StaticSet_TypedScopeCreator_EXPORTS_ONLY(v []string) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_EXPORTS_ONLY = v
}

var static_TypedScopeCreator_WELL_KNOWN_SYMBOLS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_TypedScopeCreator_WELL_KNOWN_SYMBOLS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_TypedScopeCreator()
	return static_TypedScopeCreator_WELL_KNOWN_SYMBOLS
}
func StaticSet_TypedScopeCreator_WELL_KNOWN_SYMBOLS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_TypedScopeCreator()
	static_TypedScopeCreator_WELL_KNOWN_SYMBOLS = v
}

var clinitDone_TypedScopeCreator uint32
var clinitState_TypedScopeCreator j2gort.ClinitState

func clinit_TypedScopeCreator() {
	if atomic.LoadUint32(&clinitDone_TypedScopeCreator) != 0 {
		return
	}
	clinitSlow_TypedScopeCreator()
}

func clinitSlow_TypedScopeCreator() {
	if !j2gort.EnterClinit(&clinitState_TypedScopeCreator, &clinitDone_TypedScopeCreator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TypedScopeCreator, &clinitDone_TypedScopeCreator)
	static_TypedScopeCreator_MALFORMED_TYPEDEF = DiagnosticType_Warning("JSC_MALFORMED_TYPEDEF", "Typedef for {0} does not have any type information")
	static_TypedScopeCreator_ENUM_INITIALIZER = DiagnosticType_Warning("JSC_ENUM_INITIALIZER_NOT_ENUM", "enum initializer must be an object literal or an enum")
	static_TypedScopeCreator_INVALID_ENUM_KEY = DiagnosticType_Warning("JSC_INVALID_ENUM_KEY", "enum key must be a string or numeric literal")
	static_TypedScopeCreator_CTOR_INITIALIZER = DiagnosticType_Warning("JSC_CTOR_INITIALIZER_NOT_CTOR", "Constructor {0} must be initialized at declaration")
	static_TypedScopeCreator_IFACE_INITIALIZER = DiagnosticType_Warning("JSC_IFACE_INITIALIZER_NOT_IFACE", "Interface {0} must be initialized at declaration")
	static_TypedScopeCreator_CONSTRUCTOR_EXPECTED = DiagnosticType_Warning("JSC_REFLECT_CONSTRUCTOR_EXPECTED", "Constructor expected as first argument")
	static_TypedScopeCreator_UNKNOWN_LENDS = DiagnosticType_Warning("JSC_UNKNOWN_LENDS", "Variable {0} not declared before @lends annotation.")
	static_TypedScopeCreator_LENDS_ON_NON_OBJECT = DiagnosticType_Warning("JSC_LENDS_ON_NON_OBJECT", "May only lend properties to object types. {0} has type {1}.")
	static_TypedScopeCreator_INCOMPATIBLE_ALIAS_ANNOTATION = DiagnosticType_Warning("JSC_INCOMPATIBLE_ALIAS_ANNOTATION", "Annotation {0} on {1} incompatible with aliased type.")
	static_TypedScopeCreator_DYNAMIC_EXTENDS_WITHOUT_JSDOC = DiagnosticType_Warning("JSC_DYNAMIC_EXTENDS_WITHOUT_JSDOC", "The right-hand side of an extends clause must be a qualified name, or else @extends must be specified in JSDoc")
	static_TypedScopeCreator_EXPORTS_ONLY = []string{"exports"}
	static_TypedScopeCreator_WELL_KNOWN_SYMBOLS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("Symbol.asyncDispose", "Symbol.asyncIterator", "Symbol.dispose", "Symbol.hasInstance", "Symbol.isConcatSpreadable", "Symbol.iterator", []any{"Symbol.match", "Symbol.replace", "Symbol.species", "Symbol.toPrimitive", "Symbol.toStringTag", "Symbol.unscopables"})
}

func (this_ *TypedScopeCreator) initFields_TypedScopeCreator() {
	_ = this_
	this_.memoized = j2gort.NewOrderedMapCap[pkg_rhino.Node, *TypedScope](int32(1024))
	this_.functionsWithNonEmptyReturns = j2gort.NewPlainSetVal[pkg_rhino.Node]()
	this_.fnTypeBuilderPool = make([]*FunctionTypeBuilder, int32(8))

	this_.overflowPool = make([]*TypedScope_Overflow, int32(16))

	this_.scopedVarStates = NewTypedScopeCreator_ScopedVarStateTable()
	this_.weakImports = []*TypedScopeCreator_WeakModuleImport{}
	this_.unresolvedNodes = []pkg_rhino.Node{}
	this_.undeclaredNamesForClosure = j2gort.NewOrderedSet[pkg_rhino.Node]()
	this_.stage = StaticGet_TypedScopeCreator_Stage_BUILDING()
}

// Java constructor: TypedScopeCreator(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewTypedScopeCreator_1(compiler AbstractCompiler) *TypedScopeCreator {
	clinit_TypedScopeCreator()
	this_ := &TypedScopeCreator{}
	InitTypedScopeCreator_1(this_, compiler)
	return this_
}

func InitTypedScopeCreator_1(this_ *TypedScopeCreator, compiler AbstractCompiler) {
	_ = compiler
	InitTypedScopeCreator_2(this_, compiler, compiler.GetCodingConvention())
}

// Java constructor: TypedScopeCreator(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.CodingConvention codingConvention)
func NewTypedScopeCreator_2(compiler AbstractCompiler, codingConvention CodingConvention) *TypedScopeCreator {
	clinit_TypedScopeCreator()
	this_ := &TypedScopeCreator{}
	InitTypedScopeCreator_2(this_, compiler, codingConvention)
	return this_
}

func InitTypedScopeCreator_2(this_ *TypedScopeCreator, compiler AbstractCompiler, codingConvention CodingConvention) {
	_ = compiler
	_ = codingConvention
	this_.initFields_TypedScopeCreator()
	this_.compiler = compiler
	this_.validator = compiler.GetTypeValidator()
	this_.codingConvention = codingConvention
	this_.typeRegistry = compiler.GetTypeRegistry()
	this_.typeParsingErrorReporter = this_.typeRegistry.GetErrorReporter()
	this_.unknownType = this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	if compiler.GetModuleMetadataMap() != nil {
		this_.metadataMap = compiler.GetModuleMetadataMap()
	} else {
		this_.metadataMap = NewModuleMetadataMap(j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()), j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()))
	}
	this_.moduleMap = compiler.GetModuleMap()
	this_.moduleImportResolver = NewModuleImportResolver(this_.moduleMap, this_.GetNodeToScopeMapper(), this_.typeRegistry)
	this_.processClosurePrimitives = !(this_.metadataMap.GetModulesByGoogNamespace().IsEmpty())
	this_.priv_TypedScopeCreator_clearCommonState_5738()
}

// Java method: TypedScopeCreator.acquireFunctionTypeBuilder(java.lang.String fnName, com.google.javascript.rhino.Node errorRoot, com.google.javascript.jscomp.TypedScope scope)
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_acquireFunctionTypeBuilder_5720(fnName string, errorRoot pkg_rhino.Node, scope *TypedScope) *FunctionTypeBuilder {
	_ = fnName
	_ = errorRoot
	_ = scope
	var top int32 = this_.fnTypeBuilderPoolTop
	_ = top
	if top > int32(0) {
		var nextTop int32 = (top - int32(1))
		_ = nextTop
		this_.fnTypeBuilderPoolTop = nextTop
		var b *FunctionTypeBuilder = this_.fnTypeBuilderPool[nextTop]
		_ = b
		{
			a_1 := this_.fnTypeBuilderPool
			i_2 := nextTop
			a_1[i_2] = nil
		}
		b.Reset(fnName, errorRoot, scope)
		return b
	}
	return NewFunctionTypeBuilder(fnName, this_.compiler, errorRoot, scope)
}

// Java method: TypedScopeCreator.releaseFunctionTypeBuilder(com.google.javascript.jscomp.FunctionTypeBuilder builder)
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_releaseFunctionTypeBuilder_5721(builder *FunctionTypeBuilder) {
	_ = builder
	builder.ClearForPool()
	var top int32 = this_.fnTypeBuilderPoolTop
	_ = top
	if top < int32(len(this_.fnTypeBuilderPool)) {
		{
			a_3 := this_.fnTypeBuilderPool
			i_4 := top
			a_3[i_4] = builder
		}
		this_.fnTypeBuilderPoolTop = (top + int32(1))
	}
}

// Java method: TypedScopeCreator.acquireOverflow()
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_acquireOverflow_5722() *TypedScope_Overflow {
	var top int32 = this_.overflowPoolTop
	_ = top
	if top > int32(0) {
		var nextTop int32 = (top - int32(1))
		_ = nextTop
		this_.overflowPoolTop = nextTop
		var o *TypedScope_Overflow = this_.overflowPool[nextTop]
		_ = o
		{
			a_5 := this_.overflowPool
			i_6 := nextTop
			a_5[i_6] = nil
		}
		return o
	}
	return NewTypedScope_Overflow()
}

// Java method: TypedScopeCreator.releaseOverflow(com.google.javascript.jscomp.TypedScope.Overflow overflow)
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_releaseOverflow_5723(overflow *TypedScope_Overflow) {
	_ = overflow
	overflow.ResetForPool()
	var top int32 = this_.overflowPoolTop
	_ = top
	if top < int32(len(this_.overflowPool)) {
		{
			a_7 := this_.overflowPool
			i_8 := top
			a_7[i_8] = overflow
		}
		this_.overflowPoolTop = (top + int32(1))
	}
}

// Java method: TypedScopeCreator.report(com.google.javascript.jscomp.JSError error)
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_report_5724(error_ *JSError) {
	_ = error_
	this_.compiler.Report_1(error_)
}

// Java method: TypedScopeCreator.getReferences(com.google.javascript.jscomp.TypedVar var)
func (this_ *TypedScopeCreator) GetReferences(arg_9 any) any {
	var_ := j2gort.CastPtr[TypedVar](arg_9)
	_ = var_
	return j2gort.ImmListOfPtr(var_)
}

// Java method: TypedScopeCreator.getScope(com.google.javascript.jscomp.TypedVar var)
func (this_ *TypedScopeCreator) GetScope(arg_10 any) any {
	var_ := j2gort.CastPtr[TypedVar](arg_10)
	_ = var_
	return j2gort.PtrToAny(j2gort.CastPtr[TypedScope](var_.GetScope()))
}

// Java method: TypedScopeCreator.getAllSymbols()
func (this_ *TypedScopeCreator) GetAllSymbols() any {
	var vars j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList_1_int32((this_.memoized.Size() * int32(2)))
	_ = vars
	{
		cur_11 := this_.memoized.IterValues()
		for cur_11.HasNext() {

			s := cur_11.
				Next()
			_ = s
			s.AddVarsTo(vars)
		}
	}
	return vars
}

// Java method: TypedScopeCreator.getNodeToScopeMapper()
func (this_ *TypedScopeCreator) GetNodeToScopeMapper() j2gort.Ext_java_util_function_Function {
	return func() j2gort.Ext_java_util_function_Function {
		bound_13 := &this_.memoized
		return j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.OrderedMapGetNodeAny(bound_13, mrp0_) })
	}()
}

// Java method: TypedScopeCreator.getAllMemoizedScopes()
func (this_ *TypedScopeCreator) GetAllMemoizedScopes() j2gort.Ext_java_util_Collection {
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_Lists_Reverse(j2gort.AsIface[j2gort.Ext_java_util_List](this_.memoized.ValuesToImmutableList())))
}

// Java method: TypedScopeCreator.createScope(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator) CreateScope_1(n pkg_rhino.Node) *TypedScope {
	_ = n
	var s *TypedScope = this_.memoized.Get(n)
	_ = s
	if s != nil {
		return s
	} else {
		return j2gort.CastPtr[TypedScope](this_.CreateScope_2_Node_AbstractScope(n, ToIface_AbstractScope(this_.CreateScope_1(NodeUtil_GetEnclosingScopeRoot(n.GetParent())))))
	}
}

// Java method: TypedScopeCreator.createScope(com.google.javascript.rhino.Node root, com.google.javascript.jscomp.AbstractScope<?,?> parent)
func (this_ *TypedScopeCreator) CreateScope_2_Node_AbstractScope(root pkg_rhino.Node, parent AbstractScope) any {
	_ = root
	_ = parent
	if !((parent == nil || parent.IsJavaNil_()) || j2gort.InstanceOf[*TypedScope](parent)) {
		j2gort.CheckArgument(false)
	}
	var typedParent *TypedScope = j2gort.CastPtr[TypedScope](parent)
	_ = typedParent
	var scope *TypedScope = this_.memoized.Get(root)
	_ = scope
	if scope != nil {
		if !(typedParent == j2gort.CastPtr[TypedScope](scope.GetParent())) {
			j2gort.CheckState(false)
		}
	} else {
		scope = this_.priv_TypedScopeCreator_createScopeInternal_5728(root, typedParent)
		this_.memoized.Put(root, scope)
	}
	return j2gort.PtrToAny(scope)
}

// Java method: TypedScopeCreator.createScopeInternal(com.google.javascript.rhino.Node root, com.google.javascript.jscomp.TypedScope typedParent)
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_createScopeInternal_5728(root pkg_rhino.Node, typedParent *TypedScope) *TypedScope {
	_ = root
	_ = typedParent
	{
		var iifeCond_1 bool
		{

			re_14 := this_.stage
			ae_15 := StaticGet_TypedScopeCreator_Stage_BUILDING()
			if re_14 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_14 == ae_15
		}
		if !(iifeCond_1) {
			j2gort.CheckStateFmt(false, "Cannot create scope for %s after TypedScopeCreator is frozen", j2gort.BoxNode(root))
		}
	}
	var newScope *TypedScope = nil
	_ = newScope
	var scopeBuilder TypedScopeCreator_AbstractScopeBuilder = nil
	_ = scopeBuilder
	var module *Module = ModuleImportResolver_GetModuleFromScopeRoot(this_.moduleMap, Upcast_CompilerInputProvider(this_.compiler), root)
	_ = module
	var pooledOverflow *TypedScope_Overflow = nil
	_ = pooledOverflow
	if typedParent == nil {
		if !(root.IsRoot()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(root))
		}
		var externsRoot pkg_rhino.Node = root.GetFirstChild()
		_ = externsRoot
		var jsRoot pkg_rhino.Node = root.GetSecondChild()
		_ = jsRoot
		if !(externsRoot.IsRoot()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(externsRoot))
		}
		if !(jsRoot.IsRoot()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(jsRoot))
		}
		var globalThis pkg_rhino.JSType = pkg_rhino.Upcast_JSType(this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_GLOBAL_THIS()))
		_ = globalThis
		root.SetJSType(globalThis)
		externsRoot.SetJSType(globalThis)
		jsRoot.SetJSType(globalThis)
		newScope = this_.CreateInitialScope(root)
	} else {
		var reservedNames []string = this_.reservedNamesForScope.Remove(root)
		_ = reservedNames
		if (module != nil) && module.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsGoogModule() {
			reservedNames = TypedScopeCreator_WithReservedExports(reservedNames)
		} else if ((reservedNames != nil) && root.IsFunction()) && NodeUtil_IsBundledGoogModuleCall(root.GetParent()) {
			reservedNames = TypedScopeCreator_WithoutReservedExports(reservedNames)
		}
		if ((reservedNames != nil) && (int32(len(reservedNames)) != int32(0))) || (module != nil) {
			pooledOverflow = this_.priv_TypedScopeCreator_acquireOverflow_5722()
		}
		newScope = NewTypedScope_5(typedParent, root, reservedNames, module, pooledOverflow)
	}
	if root.IsFunction() {
		scopeBuilder = NewTypedScopeCreator_FunctionScopeBuilder(this_, newScope)
	} else if (root.IsClass() || root.IsMemberFieldDef()) || root.IsComputedFieldDef() {
		scopeBuilder = NewTypedScopeCreator_ClassScopeBuilder(this_, newScope)
	} else {
		scopeBuilder = NewTypedScopeCreator_NormalScopeBuilder(this_, newScope)
		if root.IsModuleBody() {
			var inputId *pkg_rhino.InputId = j2gort.CheckNotNullPtr(NodeUtil_GetInputId(root))
			_ = inputId
			var compilerInput *CompilerInput = j2gort.CheckNotNullPtr(this_.compiler.GetInput(inputId))
			_ = compilerInput
			compilerInput.SetTypedScope(newScope)
		}
	}
	func() {
		defer func() {
			if (pooledOverflow != nil) && !(newScope.IsUsingOverflow(pooledOverflow)) {
				this_.priv_TypedScopeCreator_releaseOverflow_5723(pooledOverflow)
			}
		}()
		scopeBuilder.GetTypedScopeCreator_AbstractScopeBuilderBase().Build()
	}()
	if (module != nil) && module.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsEs6Module() {
		var namespaceType pkg_rhino.ObjectType = this_.typeRegistry.CreateAnonymousObjectType(nil)
		_ = namespaceType
		newScope.Declare("*exports*", root, pkg_rhino.Upcast_JSType(namespaceType), this_.compiler.GetInput(NodeUtil_GetInputId(root)), false)
		root.SetJSType(pkg_rhino.Upcast_JSType(namespaceType))
		this_.moduleImportResolver.UpdateEsModuleNamespaceType(namespaceType, module, newScope)
	}
	return newScope
}

// Java static method: TypedScopeCreator.withReservedExports(java.lang.String @org.jspecify.annotations.Nullable [] reservedNames)
func TypedScopeCreator_WithReservedExports(reservedNames []string) []string {
	clinit_TypedScopeCreator()
	_ = reservedNames
	if (reservedNames == nil) || (int32(len(reservedNames)) == int32(0)) {
		return StaticGet_TypedScopeCreator_EXPORTS_ONLY()
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(reservedNames)); i = i + 1 {
			{
				var iifeCond_1 bool
				{

					sr_16 := "exports"
					se_17 := reservedNames[i]
					iifeCond_1 = j2gort.CheckNotNullString(sr_16) == se_17
				}
				if iifeCond_1 {
					return reservedNames
				}
			}
		}
	}
	var out []string = j2gort.NewStringArray((int32(len(reservedNames)) + int32(1)))
	_ = out
	func() {
		acSrc_18 := reservedNames
		acSP_19 := int32(0)
		acDst_20 := out
		acDP_21 := int32(0)
		acLen_22 := int32(len(reservedNames))
		j2gort.CheckNotNullSlice(acSrc_18)
		j2gort.CheckNotNullSlice(acDst_20)
		if acSP_19 < 0 || acDP_21 < 0 || acLen_22 < 0 || int64(acSP_19)+int64(acLen_22) > int64(len(acSrc_18)) || int64(acDP_21)+int64(acLen_22) > int64(len(acDst_20)) {
			panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
		}
		copy(acDst_20[int(acDP_21):int(acDP_21)+int(acLen_22)], acSrc_18[int(acSP_19):int(acSP_19)+int(acLen_22)])
	}()
	{
		a_23 := out
		i_24 := int32(len(reservedNames))
		a_23[i_24] = "exports"
	}
	return out
}

// Java static method: TypedScopeCreator.withoutReservedExports(java.lang.String[] reservedNames)
func TypedScopeCreator_WithoutReservedExports(reservedNames []string) []string {
	clinit_TypedScopeCreator()
	_ = reservedNames
	var foundIdx int32 = int32(-1)
	_ = foundIdx
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(reservedNames)); i = i + 1 {
			{
				var iifeCond_1 bool
				{

					sr_25 := "exports"
					se_26 := reservedNames[i]
					iifeCond_1 = j2gort.CheckNotNullString(sr_25) == se_26
				}
				if iifeCond_1 {
					foundIdx = i
					break
				}
			}
		}
	}
	if foundIdx < int32(0) {
		return reservedNames
	}
	if int32(len(reservedNames)) == int32(1) {
		return nil
	}
	var out []string = j2gort.NewStringArray((int32(len(reservedNames)) - int32(1)))
	_ = out
	if foundIdx > int32(0) {
		func() {
			acSrc_27 := reservedNames
			acSP_28 := int32(0)
			acDst_29 := out
			acDP_30 := int32(0)
			acLen_31 := foundIdx
			j2gort.CheckNotNullSlice(acSrc_27)
			j2gort.CheckNotNullSlice(acDst_29)
			if acSP_28 < 0 || acDP_30 < 0 || acLen_31 < 0 || int64(acSP_28)+int64(acLen_31) > int64(len(acSrc_27)) || int64(acDP_30)+int64(acLen_31) > int64(len(acDst_29)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_29[int(acDP_30):int(acDP_30)+int(acLen_31)], acSrc_27[int(acSP_28):int(acSP_28)+int(acLen_31)])
		}()
	}
	var tail int32 = ((int32(len(reservedNames)) - foundIdx) - int32(1))
	_ = tail
	if tail > int32(0) {
		func() {
			acSrc_32 := reservedNames
			acSP_33 := (foundIdx + int32(1))
			acDst_34 := out
			acDP_35 := foundIdx
			acLen_36 := tail
			j2gort.CheckNotNullSlice(acSrc_32)
			j2gort.CheckNotNullSlice(acDst_34)
			if acSP_33 < 0 || acDP_35 < 0 || acLen_36 < 0 || int64(acSP_33)+int64(acLen_36) > int64(len(acSrc_32)) || int64(acDP_35)+int64(acLen_36) > int64(len(acDst_34)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_34[int(acDP_35):int(acDP_35)+int(acLen_36)], acSrc_32[int(acSP_33):int(acSP_33)+int(acLen_36)])
		}()
	}
	return out
}

// Java method: TypedScopeCreator.markGoogModuleExportsAsConst(com.google.javascript.rhino.Node moduleBody)
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_markGoogModuleExportsAsConst_5729(moduleBody pkg_rhino.Node) {
	_ = moduleBody
	{
		var statement pkg_rhino.Node = moduleBody.GetFirstChild()
		_ = statement
		for ; statement != 0; statement = statement.GetNext() {
			if !(NodeUtil_IsExprAssign(statement)) {
				continue
			}
			var lhs pkg_rhino.Node = statement.GetFirstFirstChild()
			_ = lhs
			if lhs.MatchesName_1_string("exports") {
				this_.undeclaredNamesForClosure.Add(lhs)
				if NodeUtil_IsNamedExportsLiteral(lhs.GetNext()) {
					{
						var key pkg_rhino.Node = lhs.GetNext().GetFirstChild()
						_ = key
						for ; key != 0; key = key.GetNext() {
							this_.undeclaredNamesForClosure.Add(key)
						}
					}
				}
			} else if lhs.IsGetProp() && lhs.GetFirstChild().MatchesName_1_string("exports") {
				this_.undeclaredNamesForClosure.Add(lhs)
			}
		}
	}
}

// Java method: TypedScopeCreator.gatherAllProvides(com.google.javascript.rhino.Node root)
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_gatherAllProvides_5730(root pkg_rhino.Node) {
	_ = root
	if !(this_.processClosurePrimitives) {
		return
	}
	var externs pkg_rhino.Node = root.GetFirstChild()
	_ = externs
	var js pkg_rhino.Node = root.GetSecondChild()
	_ = js
	var providedNames j2gort.Ext_java_util_Map = NewProcessClosureProvidesAndRequires(this_.compiler, true).CollectProvidedNames(externs, js)
	_ = providedNames
	{
		cur_37 := j2gort.NewMapValCursor(providedNames)
		for cur_37.HasNext() {
			rawElem_38 := cur_37.Next()
			name := j2gort.CastPtr[ProcessClosureProvidesAndRequires_ProvidedName](rawElem_38)
			_ = name
			if name.HasImplicitInitialization() {
				continue
			}
			if name.GetCandidateDefinition() != 0 {
				var firstDefinitionNode pkg_rhino.Node = name.GetCandidateDefinition()
				_ = firstDefinitionNode
				if NodeUtil_IsExprAssign(firstDefinitionNode) && firstDefinitionNode.GetFirstFirstChild().IsName() {
					this_.undeclaredNamesForClosure.Add(firstDefinitionNode.GetFirstFirstChild())
				}
			} else if ((name.GetFirstProvideCall() != 0) && NodeUtil_IsExprCall(name.GetFirstProvideCall())) && !(name.IsFromLegacyModule()) {
				j2gort.MultimapPutNode(this_.providedNamesFromCall, name.GetFirstProvideCall(), j2gort.PtrToAny(name))
			}
			if name.IsExplicitlyProvided() && !(name.IsFromLegacyModule()) {
				this_.typeRegistry.RegisterLegacyClosureNamespace(name.GetNamespace())
			}
		}
	}
}

// Java method: TypedScopeCreator.createInitialScope(com.google.javascript.rhino.Node root)
func (this_ *TypedScopeCreator) CreateInitialScope(root pkg_rhino.Node) *TypedScope {
	_ = root
	if !(root.IsRoot()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(root))
	}
	var val_39 MemoizedScopeCreator
	InitMemoizedScopeCreator_2(&val_39, NewSyntacticScopeCreator_1(this_.compiler), int32(2048))
	var scopeCreator *MemoizedScopeCreator = &val_39
	_ = scopeCreator
	defer scopeCreator.ReleaseSlabs_()

	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(NewTypedScopeCreator_FirstOrderFunctionAnalyzer(this_))
		(&ntb_).SetScopeCreator(scopeCreator)
		(&ntb_).TraverseRoots(root.GetFirstChild(), root.GetLastChild())
	}
	{
		clinit_NodeTraversal()
		var ntb_1_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_1_)
		(&ntb_1_).SetCompiler(this_.compiler)
		(&ntb_1_).SetCallback_1_NodeTraversal_Callback(NewTypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable(this_.typeRegistry))
		(&ntb_1_).SetScopeCreator(scopeCreator)
		(&ntb_1_).Traverse(root)
	}
	var s *TypedScope = TypedScope_CreateGlobalScope(root)
	_ = s
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_ARRAY_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_BIGINT_OBJECT_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_DATE_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_FUNCTION_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_GENERATOR_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_ITERABLE_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_ITERATOR_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_NUMBER_OBJECT_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_OBJECT_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_REGEXP_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_STRING_OBJECT_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeFunctionType_5733(s, pkg_rhino.StaticGet_JSTypeNative_SYMBOL_OBJECT_FUNCTION_TYPE())
	this_.priv_TypedScopeCreator_declareNativeValueType_5734(s, "undefined", pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
	this_.priv_TypedScopeCreator_addWellKnownSymbols_5732(s)
	this_.priv_TypedScopeCreator_gatherAllProvides_5730(root)
	this_.memoized.Put(root, s)
	return s
}

// Java method: TypedScopeCreator.addWellKnownSymbols(com.google.javascript.jscomp.TypedScope scope)
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_addWellKnownSymbols_5732(scope *TypedScope) {
	_ = scope
	{
		cur_40 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_TypedScopeCreator_WELL_KNOWN_SYMBOLS()))
		for cur_40.HasNext() {
			symbol := cur_40.NextStr()
			_ = symbol
			var type_ pkg_rhino.JSType = pkg_rhino.NewKnownSymbolType(this_.typeRegistry, symbol)
			_ = type_
			TypedScopeCreator_DeclareNativeType(scope, symbol, type_)
		}
	}
}

// Java method: TypedScopeCreator.declareNativeFunctionType(com.google.javascript.jscomp.TypedScope scope, com.google.javascript.rhino.jstype.JSTypeNative tId)
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_declareNativeFunctionType_5733(scope *TypedScope, tId pkg_rhino.JSTypeNative) {
	_ = scope
	_ = tId
	var t *pkg_rhino.FunctionType = this_.typeRegistry.GetNativeFunctionType(tId)
	_ = t
	TypedScopeCreator_DeclareNativeType(scope, t.GetInstanceType().GetReferenceName(), pkg_rhino.ToIface_JSType(t))
	TypedScopeCreator_DeclareNativeType(scope, t.GetPrototype().GetReferenceName(), pkg_rhino.Upcast_JSType(t.GetPrototype()))
}

// Java method: TypedScopeCreator.declareNativeValueType(com.google.javascript.jscomp.TypedScope scope, java.lang.String name, com.google.javascript.rhino.jstype.JSTypeNative tId)
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_declareNativeValueType_5734(scope *TypedScope, name string, tId pkg_rhino.JSTypeNative) {
	_ = scope
	_ = name
	_ = tId
	TypedScopeCreator_DeclareNativeType(scope, name, this_.typeRegistry.GetNativeType(tId))
}

// Java static method: TypedScopeCreator.declareNativeType(com.google.javascript.jscomp.TypedScope scope, java.lang.String name, com.google.javascript.rhino.jstype.JSType t)
func TypedScopeCreator_DeclareNativeType(scope *TypedScope, name string, t pkg_rhino.JSType) {
	clinit_TypedScopeCreator()
	_ = scope
	_ = name
	_ = t
	scope.Declare(name, pkg_rhino.Node(0), t, nil, false)
}

// Java method: TypedScopeCreator.setDeferredType(com.google.javascript.rhino.Node node, com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypedScopeCreator) SetDeferredType(node pkg_rhino.Node, type_ pkg_rhino.JSType) {
	_ = node
	_ = type_
	node.SetJSType(type_)
	if !(type_.GetJSTypeBase().IsResolved()) {
		this_.unresolvedNodes = append(this_.unresolvedNodes, node)
	}
}

// Java method: TypedScopeCreator.resolveWeakImportsPreResolution()
func (this_ *TypedScopeCreator) ResolveWeakImportsPreResolution() {
	for _, weakImport := range this_.weakImports {
		_ = weakImport
		weakImport.Resolve()
	}

}

// Java method: TypedScopeCreator.finishAndFreeze()
func (this_ *TypedScopeCreator) FinishAndFreeze() {
	for _, node := range this_.unresolvedNodes {
		_ = node
		var resolved pkg_rhino.JSType = node.GetJSType().GetJSTypeBase().Resolve(this_.typeParsingErrorReporter)
		_ = resolved
		node.SetJSType(resolved)
	}

	{
		cur_44 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllMemoizedScopes()))
		for cur_44.HasNext() {
			rawElem_45 := cur_44.Next()
			scope := j2gort.CastPtr[TypedScope](rawElem_45)
			_ = scope
			scope.ResolveOwnVarTypes(this_.typeParsingErrorReporter)
			scope.ValidateCompletelyBuilt()
		}
	}
	this_.priv_TypedScopeCreator_clearCommonState_5738()
	this_.stage = StaticGet_TypedScopeCreator_Stage_FROZEN()
}

// Java method: TypedScopeCreator.clearCommonState()
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_clearCommonState_5738() {
	this_.reservedNamesForScope = j2gort.NewPlainMapVal[pkg_rhino.Node, []string]()
	this_.functionsWithNonEmptyReturns = j2gort.NewPlainSetVal[pkg_rhino.Node]()
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.fnTypeBuilderPoolTop; i = i + 1 {
			{
				a_46 := this_.fnTypeBuilderPool
				i_47 := i
				a_46[i_47] = nil
			}
		}
	}
	this_.fnTypeBuilderPoolTop = int32(0)
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < this_.overflowPoolTop; i_2 = i_2 + 1 {
			{
				a_48 := this_.overflowPool
				i_49 := i_2
				a_48[i_49] = nil
			}
		}
	}
	this_.overflowPoolTop = int32(0)
	this_.scopedVarStates.Clear()
	this_.weakImports = []*TypedScopeCreator_WeakModuleImport{}
	this_.unresolvedNodes = []pkg_rhino.Node{}
	this_.undeclaredNamesForClosure = j2gort.NewOrderedSet[pkg_rhino.Node]()
	this_.providedNamesFromCall = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](j2gort.Ext_com_google_common_collect_LinkedHashMultimap_Create_0())
}

// Java method: TypedScopeCreator.getNativeType(com.google.javascript.rhino.jstype.JSTypeNative nativeType)
func (this_ *TypedScopeCreator) priv_TypedScopeCreator_getNativeType_5739(nativeType pkg_rhino.JSTypeNative) pkg_rhino.JSType {
	_ = nativeType
	return this_.typeRegistry.GetNativeType(nativeType)
}

// Java static method: TypedScopeCreator.containingGoogModuleIdOf(com.google.javascript.jscomp.TypedScope scope)
func TypedScopeCreator_ContainingGoogModuleIdOf(scope *TypedScope) string {
	clinit_TypedScopeCreator()
	_ = scope
	var module *Module = scope.GetModule()
	_ = module
	if module == nil {
		var parent *TypedScope = j2gort.CastPtr[TypedScope](scope.GetParent())
		_ = parent
		if parent != nil {
			return TypedScopeCreator_ContainingGoogModuleIdOf(parent)
		} else {
			return j2gort.NullString
		}
	}
	var metadata ModuleMetadataMap_ModuleMetadata = module.Metadata()
	_ = metadata
	if !(metadata.GetModuleMetadataMap_ModuleMetadataBase().IsModule()) {
		j2gort.CheckStateMsg(false, metadata)
	}
	return j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetFirst(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](metadata.GoogNamespaces()), nil))
}

// Java static method: TypedScopeCreator.extractKnownSymbolKey(com.google.javascript.jscomp.TypedScope currentScope, com.google.javascript.rhino.Node key)
func TypedScopeCreator_ExtractKnownSymbolKey(currentScope *TypedScope, key pkg_rhino.Node) *pkg_rhino.KnownSymbolType {
	clinit_TypedScopeCreator()
	_ = currentScope
	_ = key
	if !(key.IsQualifiedName()) {
		return nil
	}
	var type_ pkg_rhino.JSType = currentScope.LookupQualifiedNameNode(key)
	_ = type_
	if type_ != nil && !type_.IsJavaNil_() {
		return type_.ToMaybeKnownSymbolType()
	} else {
		return nil
	}
}

// Java class: com.google.javascript.jscomp.TypedScopeCreator.ScopedVarStateTable
type TypedScopeCreator_ScopedVarStateTable struct {
	keyNodes []pkg_rhino.Node `gson:"keyNodes"`
	keyNames []string         `gson:"keyNames"`
	states   []int32          `gson:"states"`
	size     int32            `gson:"size"`
	idHash_  atomic.Int32
}

func (this_ *TypedScopeCreator_ScopedVarStateTable) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedScopeCreator_ScopedVarStateTable) IsJavaNil_() bool { return this_ == nil }
func (*TypedScopeCreator_ScopedVarStateTable) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator$ScopedVarStateTable", "ScopedVarStateTable"
}
func (this_ *TypedScopeCreator_ScopedVarStateTable) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedScopeCreator_ScopedVarStateTable) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator.ScopedVarStateTable", "com.google.javascript.jscomp.TypedScopeCreator$ScopedVarStateTable", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScopeCreator_ScopedVarStateTable](v_)
	})
}

func (this_ *TypedScopeCreator_ScopedVarStateTable) initFields_TypedScopeCreator_ScopedVarStateTable() {
	_ = this_
}

// Java constructor: ScopedVarStateTable()
func NewTypedScopeCreator_ScopedVarStateTable() *TypedScopeCreator_ScopedVarStateTable {
	this_ := &TypedScopeCreator_ScopedVarStateTable{}
	InitTypedScopeCreator_ScopedVarStateTable(this_)
	return this_
}

func InitTypedScopeCreator_ScopedVarStateTable(this_ *TypedScopeCreator_ScopedVarStateTable) {
	this_.initFields_TypedScopeCreator_ScopedVarStateTable()
}

// Java static method: ScopedVarStateTable.hash(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rootNode, java.lang.String name)
func TypedScopeCreator_ScopedVarStateTable_Hash(rootNode pkg_rhino.Node, name string) int32 {
	_ = rootNode
	_ = name
	var h int32 = (((func() int32 {
		if rootNode != 0 {
			return int32(j2gort.CheckNotNullNode(rootNode))
		}
		return int32(0)
	}()) * int32(31)) + javalang.StringHashCode(j2gort.CheckNotNullString(name)))
	_ = h
	return (h ^ (javalang.Ushr32(h, int(int32(16)))))
}

// Java method: ScopedVarStateTable.clear()
func (this_ *TypedScopeCreator_ScopedVarStateTable) Clear() {
	this_.keyNodes = nil
	this_.keyNames = nil
	this_.states = nil
	this_.size = int32(0)
}

// Java method: ScopedVarStateTable.recordAssigned(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rootNode, java.lang.String name, boolean escaped)
func (this_ *TypedScopeCreator_ScopedVarStateTable) RecordAssigned(rootNode pkg_rhino.Node, name string, escaped bool) {
	_ = rootNode
	_ = name
	_ = escaped
	j2gort.CheckNotNullString(name)
	var names []string = this_.keyNames
	_ = names
	if names == nil {
		names = j2gort.NewStringArray(int32(256))
		this_.keyNames = names
		this_.keyNodes = make([]pkg_rhino.Node, int32(256))
		this_.states = make([]int32, int32(256))
	}
	var mask int32 = (int32(len(names)) - int32(1))
	_ = mask
	var slot int32 = (TypedScopeCreator_ScopedVarStateTable_Hash(rootNode, name) & mask)
	_ = slot
	for {
		var existingName string = names[slot]
		_ = existingName
		if existingName == j2gort.NullString {
			{
				a_50 := this_.keyNodes
				i_51 := slot
				a_50[i_51] = rootNode
			}
			{
				a_52 := names
				i_53 := slot
				a_52[i_53] = name
			}
			{
				a_54 := this_.states
				i_55 := slot
				a_54[i_55] = (int32(2) | (func() int32 {
					if escaped {
						return int32(1)
					}
					return int32(0)
				}()))
			}
			this_.size = this_.size + 1
			if (this_.size * int32(4)) >= (int32(len(names)) * int32(3)) {
				this_.priv_ScopedVarStateTable_rehash_5745((int32(len(names)) * int32(2)))
			}
			return
		}
		if (this_.keyNodes[slot] == rootNode) && (j2gort.CheckNotNullString(existingName) == name) {
			{
				a_56 := this_.states
				i_57 := slot
				a_56[i_57] = ((this_.states[slot] + int32(2)) | (func() int32 {
					if escaped {
						return int32(1)
					}
					return int32(0)
				}()))
			}
			return
		}
		slot = ((slot + int32(1)) & mask)
	}
}

// Java method: ScopedVarStateTable.recordEscaped(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rootNode, java.lang.String name)
func (this_ *TypedScopeCreator_ScopedVarStateTable) RecordEscaped(rootNode pkg_rhino.Node, name string) {
	_ = rootNode
	_ = name
	j2gort.CheckNotNullString(name)
	var names []string = this_.keyNames
	_ = names
	if names == nil {
		names = j2gort.NewStringArray(int32(256))
		this_.keyNames = names
		this_.keyNodes = make([]pkg_rhino.Node, int32(256))
		this_.states = make([]int32, int32(256))
	}
	var mask int32 = (int32(len(names)) - int32(1))
	_ = mask
	var slot int32 = (TypedScopeCreator_ScopedVarStateTable_Hash(rootNode, name) & mask)
	_ = slot
	for {
		var existingName string = names[slot]
		_ = existingName
		if existingName == j2gort.NullString {
			{
				a_58 := this_.keyNodes
				i_59 := slot
				a_58[i_59] = rootNode
			}
			{
				a_60 := names
				i_61 := slot
				a_60[i_61] = name
			}
			{
				a_62 := this_.states
				i_63 := slot
				a_62[i_63] = int32(1)
			}
			this_.size = this_.size + 1
			if (this_.size * int32(4)) >= (int32(len(names)) * int32(3)) {
				this_.priv_ScopedVarStateTable_rehash_5745((int32(len(names)) * int32(2)))
			}
			return
		}
		if (this_.keyNodes[slot] == rootNode) && (j2gort.CheckNotNullString(existingName) == name) {
			a_64 := this_.states
			i_65 := slot
			a_64[i_65] = (a_64[i_65] | int32(1))
			return
		}
		slot = ((slot + int32(1)) & mask)
	}
}

// Java method: ScopedVarStateTable.isEscaped(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rootNode, java.lang.@org.jspecify.annotations.Nullable String name)
func (this_ *TypedScopeCreator_ScopedVarStateTable) IsEscaped(rootNode pkg_rhino.Node, name string) bool {
	_ = rootNode
	_ = name
	var names []string = this_.keyNames
	_ = names
	if (names == nil) || (name == j2gort.NullString) {
		return false
	}
	var mask int32 = (int32(len(names)) - int32(1))
	_ = mask
	var slot int32 = (TypedScopeCreator_ScopedVarStateTable_Hash(rootNode, name) & mask)
	_ = slot
	for {
		var existingName string = names[slot]
		_ = existingName
		if existingName == j2gort.NullString {
			return false
		}
		if (this_.keyNodes[slot] == rootNode) && (j2gort.CheckNotNullString(existingName) == name) {
			return ((this_.states[slot] & int32(1)) != int32(0))
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: ScopedVarStateTable.lookupDeclareAndDecrementAssigned(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rootNode, java.lang.String name)
func (this_ *TypedScopeCreator_ScopedVarStateTable) LookupDeclareAndDecrementAssigned(rootNode pkg_rhino.Node, name string) int32 {
	_ = rootNode
	_ = name
	var names []string = this_.keyNames
	_ = names
	if (names == nil) || (name == j2gort.NullString) {
		return int32(0)
	}
	var mask int32 = (int32(len(names)) - int32(1))
	_ = mask
	var slot int32 = (TypedScopeCreator_ScopedVarStateTable_Hash(rootNode, name) & mask)
	_ = slot
	for {
		var existingName string = names[slot]
		_ = existingName
		if existingName == j2gort.NullString {
			return int32(0)
		}
		if (this_.keyNodes[slot] == rootNode) && (j2gort.CheckNotNullString(existingName) == name) {
			var st int32 = this_.states[slot]
			_ = st
			var escapedBit int32 = (st & int32(1))
			_ = escapedBit
			var count int32 = javalang.Ushr32(st, int(int32(1)))
			_ = count
			if count > int32(0) {
				{
					a_66 := this_.states
					i_67 := slot
					a_66[i_67] = ((javalang.Shl32((count - int32(1)), int(int32(1)))) | escapedBit)
				}
			}
			return (escapedBit | (func() int32 {
				if count == int32(1) {
					return int32(2)
				}
				return int32(0)
			}()))
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: ScopedVarStateTable.rehash(int newCap)
func (this_ *TypedScopeCreator_ScopedVarStateTable) priv_ScopedVarStateTable_rehash_5745(newCap int32) {
	_ = newCap
	var oldNodes []pkg_rhino.Node = this_.keyNodes
	_ = oldNodes
	var oldNames []string = this_.keyNames
	_ = oldNames
	var oldStates []int32 = this_.states
	_ = oldStates
	var newNodes []pkg_rhino.Node = make([]pkg_rhino.Node, newCap)
	_ = newNodes
	var newNames []string = j2gort.NewStringArray(newCap)
	_ = newNames
	var newStates []int32 = make([]int32, newCap)
	_ = newStates
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
				var node pkg_rhino.Node = oldNodes[i]
				_ = node
				var slot int32 = (TypedScopeCreator_ScopedVarStateTable_Hash(node, name) & mask)
				_ = slot
				for newNames[slot] != j2gort.NullString {
					slot = ((slot + int32(1)) & mask)
				}
				{
					a_68 := newNodes
					i_69 := slot
					a_68[i_69] = node
				}
				{
					a_70 := newNames
					i_71 := slot
					a_70[i_71] = name
				}
				{
					a_72 := newStates
					i_73 := slot
					a_72[i_73] = oldStates[i]
				}
			}
		}
	}
	this_.keyNodes = newNodes
	this_.keyNames = newNames
	this_.states = newStates
}

// Java enum: com.google.javascript.jscomp.TypedScopeCreator.Stage
type TypedScopeCreator_Stage uint8

const (
	enum_TypedScopeCreator_Stage_BUILDING TypedScopeCreator_Stage = 1
	enum_TypedScopeCreator_Stage_FROZEN   TypedScopeCreator_Stage = 2
)

var enumBaseTable_TypedScopeCreator_Stage = [...]j2gort.EnumBase{
	{Name_: "BUILDING", Ordinal_: 0},
	{Name_: "FROZEN", Ordinal_: 1},
}
var enumValues_TypedScopeCreator_Stage = []TypedScopeCreator_Stage{1, 2}

func StaticGet_TypedScopeCreator_Stage_BUILDING() TypedScopeCreator_Stage {
	return enum_TypedScopeCreator_Stage_BUILDING
}
func StaticGet_TypedScopeCreator_Stage_FROZEN() TypedScopeCreator_Stage {
	return enum_TypedScopeCreator_Stage_FROZEN
}

func TypedScopeCreator_Stage_Values() []TypedScopeCreator_Stage {
	return j2gort.CloneSlice(enumValues_TypedScopeCreator_Stage)
}

func TypedScopeCreator_Stage_ValueOf(name string) TypedScopeCreator_Stage {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_TypedScopeCreator_Stage {
		if enumBaseTable_TypedScopeCreator_Stage[i_].Name_ == name {
			return TypedScopeCreator_Stage(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.TypedScopeCreator.Stage." + name))
}

func (this_ TypedScopeCreator_Stage) IsJavaNil_() bool                              { return this_ == 0 }
func (TypedScopeCreator_Stage) HasIdentityEquals_()                                 {}
func (TypedScopeCreator_Stage) Implements_j2gort_Ext_java_lang_Enum()               {}
func (TypedScopeCreator_Stage) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ TypedScopeCreator_Stage) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (TypedScopeCreator_Stage) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator$Stage", "Stage"
}
func (this_ TypedScopeCreator_Stage) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_TypedScopeCreator_Stage[this_-1]
}
func (this_ TypedScopeCreator_Stage) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ TypedScopeCreator_Stage) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ TypedScopeCreator_Stage) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypedScopeCreator_Stage[this_-1].Name_
}
func (this_ TypedScopeCreator_Stage) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_TypedScopeCreator_Stage[this_-1].Name_
}
func (this_ TypedScopeCreator_Stage) Equals(other_ any) bool {
	o_, ok_ := other_.(TypedScopeCreator_Stage)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ TypedScopeCreator_Stage) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_TypedScopeCreator_Stage[this_-1].HashCode()
}
func (this_ TypedScopeCreator_Stage) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypedScopeCreator_Stage[this_-1].CompareTo(other_)
}
func (this_ TypedScopeCreator_Stage) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_TypedScopeCreator_Stage[this_-1].Name_)
}
func init() {
	vals_ := []any{TypedScopeCreator_Stage(1), TypedScopeCreator_Stage(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.TypedScopeCreator.Stage", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.TypedScopeCreator$Stage", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator.Stage", "com.google.javascript.jscomp.TypedScopeCreator$Stage", func(v_ any) bool {
		return j2gort.InstanceOf[TypedScopeCreator_Stage](v_)
	})
}

// Java class: com.google.javascript.jscomp.TypedScopeCreator.WeakModuleImport
type TypedScopeCreator_WeakModuleImport struct {
	outer_           *TypedScopeCreator
	scopedImport     ScopedName     `gson:"scopedImport"`
	localModuleScope *TypedScope    `gson:"localModuleScope"`
	moduleLocalNode  pkg_rhino.Node `gson:"moduleLocalNode"`
	idHash_          atomic.Int32
}

func (this_ *TypedScopeCreator_WeakModuleImport) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedScopeCreator_WeakModuleImport) IsJavaNil_() bool { return this_ == nil }
func (*TypedScopeCreator_WeakModuleImport) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator$WeakModuleImport", "WeakModuleImport"
}
func (this_ *TypedScopeCreator_WeakModuleImport) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedScopeCreator_WeakModuleImport) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator.WeakModuleImport", "com.google.javascript.jscomp.TypedScopeCreator$WeakModuleImport", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScopeCreator_WeakModuleImport](v_)
	})
}
func (this_ *TypedScopeCreator_WeakModuleImport) SetOuter_(o *TypedScopeCreator) { this_.outer_ = o }
func (this_ *TypedScopeCreator_WeakModuleImport) GetOuter_() *TypedScopeCreator  { return this_.outer_ }

func (this_ *TypedScopeCreator_WeakModuleImport) initFields_TypedScopeCreator_WeakModuleImport() {
	_ = this_
}

// Java constructor: WeakModuleImport(com.google.javascript.rhino.Node moduleLocalNode, com.google.javascript.jscomp.ScopedName scopedImport, com.google.javascript.jscomp.TypedScope localModuleScope)
func NewTypedScopeCreator_WeakModuleImport(outer_ *TypedScopeCreator, moduleLocalNode pkg_rhino.Node, scopedImport ScopedName, localModuleScope *TypedScope) *TypedScopeCreator_WeakModuleImport {
	this_ := &TypedScopeCreator_WeakModuleImport{}
	this_.outer_ = outer_
	InitTypedScopeCreator_WeakModuleImport(this_, moduleLocalNode, scopedImport, localModuleScope)
	return this_
}

func InitTypedScopeCreator_WeakModuleImport(this_ *TypedScopeCreator_WeakModuleImport, moduleLocalNode pkg_rhino.Node, scopedImport ScopedName, localModuleScope *TypedScope) {
	_ = moduleLocalNode
	_ = scopedImport
	_ = localModuleScope
	this_.initFields_TypedScopeCreator_WeakModuleImport()
	if !(moduleLocalNode.IsName() || moduleLocalNode.IsImportStar()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(moduleLocalNode))
	}
	this_.moduleLocalNode = moduleLocalNode
	this_.scopedImport = scopedImport
	this_.localModuleScope = localModuleScope
}

// Java method: WeakModuleImport.resolve()
func (this_ *TypedScopeCreator_WeakModuleImport) Resolve() {
	var resolvedScope *TypedScope = this_.outer_.memoized.Get(this_.scopedImport.GetScopedNameBase().GetScopeRoot())
	_ = resolvedScope
	var requiredVar *TypedVar
	_ = requiredVar
	if resolvedScope != nil {
		requiredVar = j2gort.CastPtr[TypedVar](resolvedScope.GetSlot(this_.scopedImport.GetScopedNameBase().GetName()))
	} else {
		requiredVar = nil
	}
	var type_ pkg_rhino.JSType = nil
	_ = type_
	var isInferred bool = false
	_ = isInferred
	if requiredVar != nil {
		type_ = requiredVar.GetType()
		isInferred = requiredVar.IsTypeInferred()
	} else if resolvedScope != nil {
		var nsType pkg_rhino.JSType = resolvedScope.GetTypeThroughNamespace(this_.scopedImport.GetScopedNameBase().GetName())
		_ = nsType
		if nsType != nil && !nsType.IsJavaNil_() {
			type_ = nsType
		} else {
			type_ = pkg_rhino.Upcast_JSType(this_.outer_.unknownType)
		}
		isInferred = (nsType == nil || nsType.IsJavaNil_())
	} else {
		type_ = pkg_rhino.Upcast_JSType(this_.outer_.unknownType)
		isInferred = true
	}
	this_.localModuleScope.Declare(this_.moduleLocalNode.GetString(), this_.moduleLocalNode, type_, this_.outer_.compiler.GetInput(NodeUtil_GetInputId(this_.moduleLocalNode)), isInferred)
	if (requiredVar != nil) && !(j2gort.IfaceEquals[pkg_rhino.JSType](requiredVar.GetNameNode().GetTypedefTypeProp(), nil)) {
		var typedefType pkg_rhino.JSType = requiredVar.GetNameNode().GetTypedefTypeProp()
		_ = typedefType
		this_.moduleLocalNode.SetTypedefTypeProp(typedefType)
		this_.outer_.typeRegistry.DeclareType(pkg_rhino.ToIface_StaticScope(this_.localModuleScope), this_.moduleLocalNode.GetString(), typedefType)
	}
}

// Java class: com.google.javascript.jscomp.TypedScopeCreator.RValueInfo
type TypedScopeCreator_RValueInfo struct {
	type_         pkg_rhino.JSType        `gson:"type"`
	qualifiedName pkg_rhino.QualifiedName `gson:"qualifiedName"`
	idHash_       atomic.Int32
}

func (this_ *TypedScopeCreator_RValueInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedScopeCreator_RValueInfo) IsJavaNil_() bool { return this_ == nil }
func (*TypedScopeCreator_RValueInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator$RValueInfo", "RValueInfo"
}
func (this_ *TypedScopeCreator_RValueInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedScopeCreator_RValueInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator.RValueInfo", "com.google.javascript.jscomp.TypedScopeCreator$RValueInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScopeCreator_RValueInfo](v_)
	})
}

var static_TypedScopeCreator_RValueInfo_EMPTY *TypedScopeCreator_RValueInfo

func StaticGet_TypedScopeCreator_RValueInfo_EMPTY() *TypedScopeCreator_RValueInfo {
	clinit_TypedScopeCreator_RValueInfo()
	return static_TypedScopeCreator_RValueInfo_EMPTY
}
func StaticSet_TypedScopeCreator_RValueInfo_EMPTY(v *TypedScopeCreator_RValueInfo) {
	clinit_TypedScopeCreator_RValueInfo()
	static_TypedScopeCreator_RValueInfo_EMPTY = v
}

var clinitDone_TypedScopeCreator_RValueInfo uint32
var clinitState_TypedScopeCreator_RValueInfo j2gort.ClinitState

func clinit_TypedScopeCreator_RValueInfo() {
	if atomic.LoadUint32(&clinitDone_TypedScopeCreator_RValueInfo) != 0 {
		return
	}
	clinitSlow_TypedScopeCreator_RValueInfo()
}

func clinitSlow_TypedScopeCreator_RValueInfo() {
	if !j2gort.EnterClinit(&clinitState_TypedScopeCreator_RValueInfo, &clinitDone_TypedScopeCreator_RValueInfo) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TypedScopeCreator_RValueInfo, &clinitDone_TypedScopeCreator_RValueInfo)
	static_TypedScopeCreator_RValueInfo_EMPTY = NewTypedScopeCreator_RValueInfo(nil, nil)
}

func (this_ *TypedScopeCreator_RValueInfo) initFields_TypedScopeCreator_RValueInfo() {
	_ = this_
}

// Java constructor: RValueInfo(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type, com.google.javascript.rhino.@org.jspecify.annotations.Nullable QualifiedName qualifiedName)
func NewTypedScopeCreator_RValueInfo(type_ pkg_rhino.JSType, qualifiedName pkg_rhino.QualifiedName) *TypedScopeCreator_RValueInfo {
	clinit_TypedScopeCreator_RValueInfo()
	this_ := &TypedScopeCreator_RValueInfo{}
	InitTypedScopeCreator_RValueInfo(this_, type_, qualifiedName)
	return this_
}

func InitTypedScopeCreator_RValueInfo(this_ *TypedScopeCreator_RValueInfo, type_ pkg_rhino.JSType, qualifiedName pkg_rhino.QualifiedName) {
	_ = type_
	_ = qualifiedName
	this_.initFields_TypedScopeCreator_RValueInfo()
	this_.type_ = type_
	this_.qualifiedName = qualifiedName
}

// Java static method: RValueInfo.empty()
func TypedScopeCreator_RValueInfo_Empty() *TypedScopeCreator_RValueInfo {
	clinit_TypedScopeCreator_RValueInfo()
	return StaticGet_TypedScopeCreator_RValueInfo_EMPTY()
}

// Java class: com.google.javascript.jscomp.TypedScopeCreator.IdentifyEnumsAndTypedefsAsNonNullable
type TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	registry *pkg_rhino.JSTypeRegistry `gson:"registry"`
}

func (this_ *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) IsJavaNil_() bool {
	return this_ == nil
}
func (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator$IdentifyEnumsAndTypedefsAsNonNullable", "IdentifyEnumsAndTypedefsAsNonNullable"
}
func (this_ *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator.IdentifyEnumsAndTypedefsAsNonNullable", "com.google.javascript.jscomp.TypedScopeCreator$IdentifyEnumsAndTypedefsAsNonNullable", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable](v_)
	})
}
func (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) initFields_TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable() {
	_ = this_
}

// Java constructor: IdentifyEnumsAndTypedefsAsNonNullable(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewTypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable(registry *pkg_rhino.JSTypeRegistry) *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable {
	this_ := &TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable{}
	InitTypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable(this_, registry)
	return this_
}

func InitTypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable(this_ *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable, registry *pkg_rhino.JSTypeRegistry) {
	_ = registry
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable()
	this_.registry = registry
}

// Java method: IdentifyEnumsAndTypedefsAsNonNullable.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) Visit_3(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = node
	_ = parent
	switch node.TokenSwitchOrd_() {
	case int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(84) /* Token.VAR */ :
		{
			var child pkg_rhino.Node = node.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				this_.priv_IdentifyEnumsAndTypedefsAsNonNullable_identifyEnumOrTypedefDeclaration_5747(t, child, child.GetFirstChild(), NodeUtil_GetBestJSDocInfo(child))
			}
		}
	case int32(93) /* Token.EXPR_RESULT */ :
		var firstChild pkg_rhino.Node = node.GetFirstChild()
		_ = firstChild
		if firstChild.IsAssign() {
			var assign pkg_rhino.Node = firstChild
			_ = assign
			this_.priv_IdentifyEnumsAndTypedefsAsNonNullable_identifyEnumOrTypedefDeclaration_5747(t, assign.GetFirstChild(), assign.GetSecondChild(), assign.GetJSDocInfo())
		} else if firstChild.IsGetProp() {
			this_.priv_IdentifyEnumsAndTypedefsAsNonNullable_identifyEnumOrTypedefDeclaration_5747(t, firstChild, pkg_rhino.Node(0), firstChild.GetJSDocInfo())
		}
	default:
	}
}

// Java method: IdentifyEnumsAndTypedefsAsNonNullable.identifyEnumOrTypedefDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nameNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rvalue, com.google.javascript.rhino.JSDocInfo info)
func (this_ *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) priv_IdentifyEnumsAndTypedefsAsNonNullable_identifyEnumOrTypedefDeclaration_5747(t *NodeTraversal, nameNode pkg_rhino.Node, rvalue pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = t
	_ = nameNode
	_ = rvalue
	_ = info
	if !(nameNode.IsQualifiedName()) {
		return
	}
	if (info != nil) && info.HasEnumParameterType() {
		this_.registry.IdentifyNonNullableName(pkg_rhino.ToIface_StaticScope(t.GetScope()), nameNode.GetQualifiedName())
	} else if (info != nil) && info.HasTypedefType() {
		this_.registry.IdentifyNonNullableName(pkg_rhino.ToIface_StaticScope(t.GetScope()), nameNode.GetQualifiedName())
	} else if (((rvalue != 0) && rvalue.IsQualifiedName()) && this_.registry.IsNonNullableName(pkg_rhino.ToIface_StaticScope(t.GetScope()), rvalue.GetQualifiedName())) && NodeUtil_IsConstantDeclaration(info, nameNode) {
		this_.registry.IdentifyNonNullableName(pkg_rhino.ToIface_StaticScope(t.GetScope()), nameNode.GetQualifiedName())
	}
}

// Java polymorphic class interface: com.google.javascript.jscomp.TypedScopeCreator.AbstractScopeBuilder
type TypedScopeCreator_AbstractScopeBuilder interface {
	Implements_jscomp_TypedScopeCreator_AbstractScopeBuilder()
	IsJavaNil_() bool
	GetTypedScopeCreator_AbstractScopeBuilderBase() *TypedScopeCreator_AbstractScopeBuilderBase
	NodeTraversal_Callback
	DeferAction(pkg_rhino.Node, j2gort.Ext_java_lang_Runnable)
	GetCompilerInput() *CompilerInput
	GetModule() *Module
	Build()
	VisitPreorder(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
	VisitPostorder(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
	AttachLiteralTypes(pkg_rhino.Node)
	ProcessObjectLitProperties(pkg_rhino.Node, pkg_rhino.ObjectType, bool)
	FindOwnerTemplateMap(pkg_rhino.Node) *pkg_rhino.TemplateTypeMap
	AssertDefinitionNode(pkg_rhino.Node, pkg_rhino.Token)
	DefineCatch(pkg_rhino.Node)
	DefineAssignAsIfVarDeclaration(pkg_rhino.Node)
	DefineVars(pkg_rhino.Node)
	DefineVarChild(*pkg_rhino.JSDocInfo, pkg_rhino.Node, *TypedScope)
	DefineDestructuringPatternInVarDeclaration(pkg_rhino.Node, *TypedScope, j2gort.Ext_com_google_common_base_Supplier)
	DefineClassLiteral(pkg_rhino.Node)
	DefineFunctionLiteral(pkg_rhino.Node)
	GetDeclaredType(*pkg_rhino.JSDocInfo, pkg_rhino.Node, pkg_rhino.Node, j2gort.Ext_com_google_common_base_Supplier) pkg_rhino.JSType
	IsGoogModuleExports(pkg_rhino.Node) bool
	CheckForClassDefiningCalls(pkg_rhino.Node)
	MaybeDeclareGetElem(*pkg_rhino.JSDocInfo, pkg_rhino.Node, pkg_rhino.Node)
	MaybeDeclareQualifiedName(*NodeTraversal, *pkg_rhino.JSDocInfo, pkg_rhino.Node, pkg_rhino.Node, pkg_rhino.Node)
	DeclareProvidedNs(pkg_rhino.Node, *ProcessClosureProvidesAndRequires_ProvidedName)
	ResolveStubDeclaration(pkg_rhino.Node, bool, string)
	DeclareTypedefType(pkg_rhino.Node, *pkg_rhino.JSDocInfo)
	DeclarePropertyIfNamespaceType_6_ObjectType_Node_string_JSType_Node_JSDocInfo(pkg_rhino.ObjectType, pkg_rhino.Node, string, pkg_rhino.JSType, pkg_rhino.Node, *pkg_rhino.JSDocInfo)
	DeclarePropertyIfNamespaceType_6_ObjectType_Node_Property_Key_JSType_Node_JSDocInfo(pkg_rhino.ObjectType, pkg_rhino.Node, pkg_rhino.Property_Key, pkg_rhino.JSType, pkg_rhino.Node, *pkg_rhino.JSDocInfo)
}

func ToIface_TypedScopeCreator_AbstractScopeBuilder[T interface {
	TypedScopeCreator_AbstractScopeBuilder
	comparable
}](v T) TypedScopeCreator_AbstractScopeBuilder {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_TypedScopeCreator_AbstractScopeBuilder[T TypedScopeCreator_AbstractScopeBuilder](v T) TypedScopeCreator_AbstractScopeBuilder {
	var vi TypedScopeCreator_AbstractScopeBuilder = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.TypedScopeCreator.AbstractScopeBuilder
type TypedScopeCreator_AbstractScopeBuilderBase struct {
	self_             TypedScopeCreator_AbstractScopeBuilder
	outer_            *TypedScopeCreator
	currentScope      *TypedScope                                      `gson:"currentScope"`
	currentHoistScope *TypedScope                                      `gson:"currentHoistScope"`
	sourceName        string                                           `gson:"sourceName"`
	inputId           *pkg_rhino.InputId                               `gson:"inputId"`
	deferredActions   j2gort.Ext_com_google_common_collect_SetMultimap `gson:"deferredActions"`
	idHash_           atomic.Int32
}

func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(TypedScopeCreator_AbstractScopeBuilder)
}
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*TypedScopeCreator_AbstractScopeBuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator$AbstractScopeBuilder", "AbstractScopeBuilder"
}
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator.AbstractScopeBuilder", "com.google.javascript.jscomp.TypedScopeCreator$AbstractScopeBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[TypedScopeCreator_AbstractScopeBuilder](v_)
	})
}
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetTypedScopeCreator_AbstractScopeBuilderBase() *TypedScopeCreator_AbstractScopeBuilderBase {
	return b
}
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) SetOuter_(o *TypedScopeCreator) {
	this_.outer_ = o
}
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) GetOuter_() *TypedScopeCreator {
	return this_.outer_
}
func (*TypedScopeCreator_AbstractScopeBuilderBase) Implements_jscomp_TypedScopeCreator_AbstractScopeBuilder() {
}
func (*TypedScopeCreator_AbstractScopeBuilderBase) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) initFields_TypedScopeCreator_AbstractScopeBuilder() {
	_ = this_
	this_.sourceName = j2gort.NullString

}

func InitTypedScopeCreator_AbstractScopeBuilder(this_ *TypedScopeCreator_AbstractScopeBuilderBase, scope *TypedScope) {
	_ = scope
	this_.initFields_TypedScopeCreator_AbstractScopeBuilder()
	this_.currentScope = scope
	this_.currentHoistScope = j2gort.CastPtr[TypedScope](scope.GetClosestHoistScope())
}

// Java method: AbstractScopeBuilder.deferAction(com.google.javascript.rhino.Node key, java.lang.Runnable action)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DeferAction(key pkg_rhino.Node, action j2gort.Ext_java_lang_Runnable) {
	_ = key
	_ = action
	var actions j2gort.Ext_com_google_common_collect_SetMultimap = this_.deferredActions
	_ = actions
	if actions == nil {
		actions = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_LinkedHashMultimap_Create_0())
		this_.deferredActions = actions
	}
	j2gort.MultimapPutNode(actions, key, action)
}

// Java method: AbstractScopeBuilder.getCompilerInput()
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) GetCompilerInput() *CompilerInput {
	return this_.outer_.compiler.GetInput(this_.inputId)
}

// Java method: AbstractScopeBuilder.getModule()
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) GetModule() *Module {
	return this_.currentScope.GetModule()
}

// Java method: AbstractScopeBuilder.build()
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) Build() {
	this_.priv_AbstractScopeBuilder_initializeModuleScope_5752(this_.currentScope.GetRootNode_0())
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.outer_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(Upcast_NodeTraversal_Callback(this_.self_))
		(&ntb_).SetScopeCreator(StaticGet_ScopeCreator_ASSERT_NO_SCOPES_CREATED())
		(&ntb_).TraverseAtScope(ToIface_AbstractScope(this_.currentScope))
	}
	this_.priv_AbstractScopeBuilder_finishDeclaringGoogModule_5754()
}

// Java method: AbstractScopeBuilder.initializeModuleScope(com.google.javascript.rhino.Node moduleBody)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_initializeModuleScope_5752(moduleBody pkg_rhino.Node) {
	_ = moduleBody
	if this_.GetModule() == nil {
		return
	}
	var moduleType ModuleMetadataMap_ModuleType = this_.GetModule().Metadata().ModuleType()
	_ = moduleType
	switch (moduleType).Ordinal() {
	case int32(2) /* ModuleType.LEGACY_GOOG_MODULE */, int32(1) /* ModuleType.GOOG_MODULE */ :
		this_.priv_AbstractScopeBuilder_declareExportsInGoogModuleScope_5753(this_.GetModule(), moduleBody)
		this_.outer_.priv_TypedScopeCreator_markGoogModuleExportsAsConst_5729(moduleBody)
	case int32(0) /* ModuleType.ES6_MODULE */ :
		var unresolvedImports j2gort.Ext_com_google_common_collect_ImmutableMap = this_.outer_.moduleImportResolver.DeclareEsModuleImports(this_.GetModule(), this_.currentScope, this_.outer_.compiler.GetInput(NodeUtil_GetInputId(moduleBody)))
		_ = unresolvedImports
		{
			cur_74 := j2gort.NewMapEntryCursor(unresolvedImports)
			for cur_74.NextEntry() {
				this_.outer_.weakImports = append(this_.outer_.weakImports, NewTypedScopeCreator_WeakModuleImport(this_.outer_, pkg_rhino.Node(cur_74.KeyNode()), j2gort.AsIface[ScopedName](cur_74.GetValue()), this_.currentScope))
			}
		}
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("Unexpected module type %s in module %s", []any{j2gort.BoxNullableVal(moduleType), j2gort.PtrToAny(this_.GetModule())})))))
	}
}

// Java method: AbstractScopeBuilder.declareExportsInGoogModuleScope(com.google.javascript.jscomp.modules.Module googModule, com.google.javascript.rhino.Node moduleBody)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_declareExportsInGoogModuleScope_5753(googModule *Module, moduleBody pkg_rhino.Node) {
	_ = googModule
	_ = moduleBody
	{
		var statement pkg_rhino.Node = moduleBody.GetFirstChild()
		_ = statement
		for ; statement != 0; statement = statement.GetNext() {
			if !(NodeUtil_IsExprAssign(statement)) {
				continue
			}
			var lhs pkg_rhino.Node = statement.GetFirstFirstChild()
			_ = lhs
			if lhs.MatchesName_1_string("exports") {
				return
			}
		}
	}
	var root pkg_rhino.Node = googModule.Metadata().RootNode()
	_ = root
	NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).ForDeclarationNode(root.GetFirstChild()).ForVariableName("exports").WithType(pkg_rhino.Upcast_JSType(this_.outer_.typeRegistry.CreateAnonymousObjectType(nil))).AllowLaterTypeInference(false).InScope(this_.currentScope).DefineSlot()
}

// Java method: AbstractScopeBuilder.finishDeclaringGoogModule()
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_finishDeclaringGoogModule_5754() {
	if (this_.GetModule() == nil) || !(this_.GetModule().Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsGoogModule()) {
		return
	}
	var exportsVar *TypedVar = j2gort.CheckNotNullPtr(j2gort.CastPtr[TypedVar](this_.currentScope.GetSlot("exports")))
	_ = exportsVar
	if this_.GetModule().Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsLegacyGoogModule() {
		var moduleNamespace string = this_.GetModule().ClosureNamespace()
		_ = moduleNamespace
		this_.outer_.typeRegistry.RegisterLegacyClosureNamespace(moduleNamespace)
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).InScope(j2gort.CastPtr[TypedScope](this_.currentScope.GetGlobalScope())).ForDeclarationNode(exportsVar.GetNameNode()).ForVariableName(moduleNamespace).WithType(exportsVar.GetType()).AllowLaterTypeInference(exportsVar.IsTypeInferred()).ForGoogProvidedName().DefineSlot()
		var lastDot int32 = j2gort.JStrLastIndexOfCp(moduleNamespace, int32(uint16(46)))
		_ = lastDot
		if (lastDot >= int32(0)) && !(exportsVar.IsTypeInferred()) {
			var parentType pkg_rhino.JSType = j2gort.CastPtr[TypedScope](this_.currentScope.GetGlobalScope()).LookupQualifiedName_1_string(j2gort.JStrSub(moduleNamespace, int32(0), lastDot))
			_ = parentType
			if (parentType != nil && !parentType.IsJavaNil_()) && !(j2gort.IfaceEquals[pkg_rhino.ObjectType](parentType.GetJSTypeBase().ToMaybeObjectType(), nil)) {
				var typeDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(exportsVar.GetNameNode())
				_ = typeDoc
				this_.DeclarePropertyIfNamespaceType_6_ObjectType_Node_string_JSType_Node_JSDocInfo(parentType.GetJSTypeBase().ToMaybeObjectType(), exportsVar.GetNameNode(), j2gort.JStrSubFrom(moduleNamespace, (lastDot+int32(1))), exportsVar.GetType(), exportsVar.GetNameNode(), typeDoc)
			}
		}
		this_.priv_AbstractScopeBuilder_declareAliasTypeIfRvalueIsAliasable_5791(moduleNamespace, exportsVar.GetNameNode(), "exports", exportsVar.GetType(), this_.currentScope, j2gort.CastPtr[TypedScope](this_.currentScope.GetGlobalScope()))
	} else {
		this_.outer_.typeRegistry.RegisterNonLegacyClosureNamespace(this_.GetModule().ClosureNamespace(), exportsVar.GetNameNode(), exportsVar.GetType())
	}
	var rootNode pkg_rhino.Node = this_.currentScope.GetRootNode_0()
	_ = rootNode
	if rootNode.IsModuleBody() {
		rootNode.SetJSType(exportsVar.GetType())
	} else {
		if !(rootNode.IsBlock()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(rootNode))
		}
		var fn pkg_rhino.Node = rootNode.GetParent()
		_ = fn
		var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(fn)
		_ = paramList
		paramList.GetOnlyChild().SetJSType(exportsVar.GetType())
	}
}

// Java method: AbstractScopeBuilder.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	this_.inputId = t.GetInputId()
	if (n.IsFunction() || n.IsScript()) || ((parent == 0) && (this_.inputId != nil)) {
		j2gort.CheckNotNullPtr(this_.inputId)
		this_.sourceName = NodeUtil_GetSourceName(n)
	}
	if (parent == 0) || this_.priv_AbstractScopeBuilder_inCurrentScope_5755(t) {
		this_.self_.VisitPreorder(t, n, parent)
		return true
	} else if n.IsModuleBody() {
		this_.outer_.CreateScope_2_Node_AbstractScope(n, ToIface_AbstractScope(this_.currentScope))
	} else if NodeUtil_IsBundledGoogModuleScopeRoot(n) {
		var functionScope *TypedScope = j2gort.CastPtr[TypedScope](this_.outer_.CreateScope_2_Node_AbstractScope(parent, ToIface_AbstractScope(this_.currentScope)))
		_ = functionScope
		this_.outer_.CreateScope_2_Node_AbstractScope(n, ToIface_AbstractScope(functionScope))
	}
	return false
}

// Java method: AbstractScopeBuilder.inCurrentScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_inCurrentScope_5755(t *NodeTraversal) bool {
	_ = t
	return (t.GetScopeRoot() == this_.currentScope.GetRootNode_0())
}

// Java method: AbstractScopeBuilder.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var isScopeRoot bool = (parent == 0)
	_ = isScopeRoot
	this_.inputId = t.GetInputId()
	if !(isScopeRoot) || !(n.IsClass()) {
		this_.AttachLiteralTypes(n)
		this_.self_.VisitPostorder(t, n, parent)
	}
	var actions j2gort.Ext_com_google_common_collect_SetMultimap = this_.deferredActions
	_ = actions
	if actions != nil {
		if !(isScopeRoot) && j2gort.MultimapContainsKeyNode(actions, n) {
			{
				cur_75 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.MultimapRemoveAllNode(actions, n))))
				for cur_75.HasNext() {
					rawElem_76 := cur_75.Next()
					action := j2gort.AsIface[j2gort.Ext_java_lang_Runnable](rawElem_76)
					_ = action
					action.Run()
				}
			}
		} else if isScopeRoot && !(actions.IsEmpty()) {
			{
				cur_77 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Collection](actions.Values())))
				for cur_77.HasNext() {
					rawElem_78 := cur_77.Next()
					action_2 := j2gort.AsIface[j2gort.Ext_java_lang_Runnable](rawElem_78)
					_ = action_2
					action_2.Run()
				}
			}
		}
	}
}

// Java method: AbstractScopeBuilder.visitPreorder(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) VisitPreorder(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}

// Java method: AbstractScopeBuilder.visitPostorder(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) VisitPostorder(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}

// Java method: AbstractScopeBuilder.attachLiteralTypes(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) AttachLiteralTypes(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(36) /* Token.NULL */ :
		n.SetJSType(this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_NULL_TYPE()))
	case int32(87) /* Token.VOID */ :
		n.SetJSType(this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE()))
	case int32(35) /* Token.STRINGLIT */, int32(133) /* Token.TEMPLATELIT_STRING */ :
		n.SetJSType(this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE()))
	case int32(33) /* Token.NUMBER */ :
		n.SetJSType(this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE()))
	case int32(34) /* Token.BIGINT */ :
		n.SetJSType(this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_BIGINT_TYPE()))
	case int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */ :
		n.SetJSType(this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE()))
	case int32(42) /* Token.REGEXP */ :
		n.SetJSType(this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_REGEXP_TYPE()))
	case int32(47) /* Token.OBJECTLIT */ :
		var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
		_ = info
		if (info != nil) && info.HasLendsName() {
			this_.DeferAction(NodeUtil_GetEnclosingStatement(n), j2gort.NewLambda_Ext_java_lang_Runnable(func() {
				this_.priv_AbstractScopeBuilder_defineObjectLiteral_5759(n)
			}))
		} else {
			this_.priv_AbstractScopeBuilder_defineObjectLiteral_5759(n)
		}
	case int32(105) /* Token.CLASS */ :
		this_.DefineClassLiteral(n)
	case int32(46) /* Token.ARRAYLIT */ :
		n.SetJSType(this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE()))
	default:
	}
}

// Java method: AbstractScopeBuilder.defineObjectLiteral(com.google.javascript.rhino.Node objectLit)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_defineObjectLiteral_5759(objectLit pkg_rhino.Node) {
	_ = objectLit
	var type_ pkg_rhino.JSType = nil
	_ = type_
	var info *pkg_rhino.JSDocInfo = objectLit.GetJSDocInfo()
	_ = info
	if (info != nil) && info.HasLendsName() {
		var lendsName string = info.GetLendsName().GetRoot().GetString()
		_ = lendsName
		var lendsVar *TypedVar = this_.currentScope.LookupTypedVar(lendsName)
		_ = lendsVar
		if lendsVar == nil {
			this_.outer_.priv_TypedScopeCreator_report_5724(JSError_Make_3(objectLit, StaticGet_TypedScopeCreator_UNKNOWN_LENDS(), []string{lendsName}))
		} else {
			type_ = lendsVar.GetType()
			if type_ == nil || type_.IsJavaNil_() {
				type_ = pkg_rhino.Upcast_JSType(this_.outer_.unknownType)
			}
			if !(type_.GetJSTypeBase().IsSubtypeOf_1(this_.outer_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE()))) {
				this_.outer_.priv_TypedScopeCreator_report_5724(JSError_Make_3(objectLit, StaticGet_TypedScopeCreator_LENDS_ON_NON_OBJECT(), []string{lendsName, type_.GetJSTypeBase().String()}))
				type_ = nil
			} else {
				objectLit.SetJSType(type_)
			}
		}
	}
	info = NodeUtil_GetBestJSDocInfo(objectLit)
	var createEnumType bool = ((info != nil) && info.HasEnumParameterType())
	_ = createEnumType
	if createEnumType {
		var lValue pkg_rhino.Node = NodeUtil_GetBestLValue(objectLit)
		_ = lValue
		var lValueName string = NodeUtil_GetBestLValueName(lValue)
		_ = lValueName
		type_ = pkg_rhino.ToIface_JSType(this_.priv_AbstractScopeBuilder_createEnumTypeFromNodes_5784(objectLit, lValueName, lValue, info))
	}
	if type_ == nil || type_.IsJavaNil_() {
		type_ = pkg_rhino.Upcast_JSType(this_.outer_.typeRegistry.CreateAnonymousObjectType(info))
	}
	this_.outer_.SetDeferredType(objectLit, type_)
	this_.ProcessObjectLitProperties(objectLit, pkg_rhino.ObjectType_Cast(objectLit.GetJSType()), !(createEnumType))
}

// Java method: AbstractScopeBuilder.processObjectLitProperties(com.google.javascript.rhino.Node objLit, com.google.javascript.rhino.jstype.ObjectType objLitType, boolean declareOnOwner)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) ProcessObjectLitProperties(objLit pkg_rhino.Node, objLitType pkg_rhino.ObjectType, declareOnOwner bool) {
	_ = objLit
	_ = objLitType
	_ = declareOnOwner
	var ownerName string = j2gort.NullString
	_ = ownerName
	var ownerNameComputed bool = false
	_ = ownerNameComputed
	{
		var keyNode pkg_rhino.Node = objLit.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			if keyNode.IsSpread() {
				continue
			}
			var memberSymbol *pkg_rhino.KnownSymbolType = nil
			_ = memberSymbol
			var rawKeyName string = j2gort.NullString
			_ = rawKeyName
			var value pkg_rhino.Node = pkg_rhino.Node(0)
			_ = value
			var qualifiedName string = j2gort.NullString
			_ = qualifiedName
			if keyNode.IsComputedProp() {
				memberSymbol = TypedScopeCreator_ExtractKnownSymbolKey(this_.currentScope, keyNode.GetFirstChild())
				if memberSymbol == nil {
					continue
				}
				rawKeyName = j2gort.NullString
				qualifiedName = j2gort.NullString
				value = keyNode.GetSecondChild()
			} else {
				memberSymbol = nil
				value = keyNode.GetFirstChild()
				rawKeyName = j2gort.CheckNotNullString(NodeUtil_GetObjectOrClassLitKeyName(keyNode))
				if NodeUtil_MayBeObjectLitKey(keyNode) {
					if !(ownerNameComputed) {
						ownerName = NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(objLit))
						ownerNameComputed = true
					}
					if (ownerName != j2gort.NullString) && pkg_rhino.TokenStream_IsJSIdentifier(rawKeyName) {
						qualifiedName = (j2gort.StrValString(ownerName) + "." + j2gort.StrValString(rawKeyName))
					} else {
						qualifiedName = j2gort.NullString
					}
				} else {
					qualifiedName = NodeUtil_GetBestLValueName(keyNode)
				}
			}
			var info *pkg_rhino.JSDocInfo = keyNode.GetJSDocInfo()
			_ = info
			var valueType pkg_rhino.JSType = this_.GetDeclaredType(info, keyNode, value, nil)
			_ = valueType
			var keyType pkg_rhino.JSType
			_ = keyType
			if objLitType.GetObjectTypeBase().IsEnumType() {
				keyType = pkg_rhino.ToIface_JSType(objLitType.ToMaybeEnumType().GetElementsType())
			} else {
				keyType = TypeCheck_GetObjectLitKeyTypeFromValueType(keyNode, valueType)
			}
			if qualifiedName != j2gort.NullString {
				NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).ForDeclarationNode(keyNode).ForVariableName(qualifiedName).InScope(this_.priv_AbstractScopeBuilder_getLValueRootScope_5788(keyNode)).WithType(keyType).AllowLaterTypeInference((keyType == nil || keyType.IsJavaNil_())).DefineSlot()
			} else if keyType != nil && !keyType.IsJavaNil_() {
				this_.outer_.SetDeferredType(keyNode, keyType)
			}
			if ((keyType != nil && !keyType.IsJavaNil_()) && (objLitType != nil && !objLitType.IsJavaNil_())) && declareOnOwner {
				if memberSymbol != nil {
					objLitType.GetObjectTypeBase().DefineDeclaredProperty_3_Property_Key_JSType_Node(pkg_rhino.NewProperty_SymbolKey(memberSymbol), keyType, keyNode)
				} else {
					objLitType.GetObjectTypeBase().DefineDeclaredProperty_3_string_JSType_Node(rawKeyName, keyType, keyNode)
				}
			}
		}
	}
}

// Java method: AbstractScopeBuilder.getDeclaredTypeInAnnotation(com.google.javascript.rhino.Node node, com.google.javascript.rhino.JSDocInfo info)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_getDeclaredTypeInAnnotation_5761(node pkg_rhino.Node, info *pkg_rhino.JSDocInfo) pkg_rhino.JSType {
	_ = node
	_ = info
	if !(info.HasType()) {
		j2gort.CheckArgumentMsg(false, j2gort.PtrToAny(info))
	}
	var ownerTemplateMap *pkg_rhino.TemplateTypeMap = this_.FindOwnerTemplateMap(node)
	_ = ownerTemplateMap
	var templateScope pkg_rhino.StaticTypedScope
	_ = templateScope
	if (ownerTemplateMap != nil) && !(ownerTemplateMap.IsEmpty()) {
		templateScope = this_.outer_.typeRegistry.CreateScopeWithTemplates_2_StaticTypedScope_TemplateTypeMap(pkg_rhino.ToIface_StaticTypedScope(this_.currentScope), ownerTemplateMap)
	} else {
		templateScope = pkg_rhino.ToIface_StaticTypedScope(this_.currentScope)
	}
	return info.GetType().Evaluate(templateScope, this_.outer_.typeRegistry)
}

// Java method: AbstractScopeBuilder.findOwnerTemplateMap(com.google.javascript.rhino.Node node)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) FindOwnerTemplateMap(node pkg_rhino.Node) *pkg_rhino.TemplateTypeMap {
	_ = node
	if ((node != 0) && !(node.IsStaticMember())) && node.GetParent().IsClassMembers() {
		var classType pkg_rhino.JSType = node.GetGrandparent().GetJSType()
		_ = classType
		if classType != nil && !classType.IsJavaNil_() {
			return classType.ToMaybeFunctionType().GetTemplateTypeMap()
		}
		return nil
	}
	var ownerNode pkg_rhino.Node = NodeUtil_GetBestLValueOwner(node)
	_ = ownerNode
	var ownerName string = NodeUtil_GetBestLValueName(ownerNode)
	_ = ownerName
	if ownerName != j2gort.NullString {
		var ownerVar *TypedVar = this_.currentScope.LookupTypedVar(ownerName)
		_ = ownerVar
		if ownerVar != nil {
			var ownerType pkg_rhino.JSType = pkg_rhino.Upcast_JSType(this_.priv_AbstractScopeBuilder_getPrototypeOwnerType_5782(pkg_rhino.ObjectType_Cast(ownerVar.GetType())))
			_ = ownerType
			if ownerType != nil && !ownerType.IsJavaNil_() {
				return ownerType.GetTemplateTypeMap()
			}
		}
	}
	return nil
}

// Java method: AbstractScopeBuilder.assertDefinitionNode(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token type)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) AssertDefinitionNode(n pkg_rhino.Node, type_ pkg_rhino.Token) {
	_ = n
	_ = type_
	if !(this_.sourceName != j2gort.NullString) {
		j2gort.CheckState(false)
	}
	if !(n.GetToken() == type_) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
}

// Java method: AbstractScopeBuilder.defineCatch(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DefineCatch(n pkg_rhino.Node) {
	_ = n
	this_.AssertDefinitionNode(n, pkg_rhino.StaticGet_Token_CATCH())
	NodeUtil_VisitLhsNodesInNode(n, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_79 any) {
		catchName := j2gort.CastNode[pkg_rhino.Node](lp_79)
		_ = catchName
		var type_ pkg_rhino.JSType = this_.GetDeclaredType(catchName.GetJSDocInfo(), catchName, pkg_rhino.Node(0), nil)
		_ = type_
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).ForDeclarationNode(catchName).ForVariableName(catchName.GetString()).InScope(this_.currentScope).WithType(type_).AllowLaterTypeInference((type_ == nil || type_.IsJavaNil_())).DefineSlot()
	}))
}

// Java method: AbstractScopeBuilder.defineAssignAsIfVarDeclaration(com.google.javascript.rhino.Node assignment)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DefineAssignAsIfVarDeclaration(assignment pkg_rhino.Node) {
	_ = assignment
	var info *pkg_rhino.JSDocInfo = assignment.GetJSDocInfo()
	_ = info
	var name pkg_rhino.Node = assignment.GetFirstChild()
	_ = name
	if !(name.IsName()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(name))
	}
	var rvalue pkg_rhino.Node = assignment.GetSecondChild()
	_ = rvalue
	this_.priv_AbstractScopeBuilder_defineName_5773(name, rvalue, j2gort.CastPtr[TypedScope](this_.currentScope.GetClosestHoistScope()), info)
}

// Java method: AbstractScopeBuilder.defineVars(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DefineVars(n pkg_rhino.Node) {
	_ = n
	if !(this_.sourceName != j2gort.NullString) {
		j2gort.CheckState(false)
	}
	if !(NodeUtil_IsNameDeclaration(n)) {
		j2gort.CheckState(false)
	}
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	var scope *TypedScope
	_ = scope
	if n.IsVar() {
		scope = this_.currentHoistScope
	} else {
		scope = this_.currentScope
	}
	if n.HasMoreThanOneChild() && (info != nil) {
		this_.outer_.priv_TypedScopeCreator_report_5724(JSError_Make_3(n, StaticGet_TypeCheck_MULTIPLE_VAR_DEF(), []string{}))
	}
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.DefineVarChild(info, child, scope)
		}
	}
	if n.HasOneChild() && this_.priv_AbstractScopeBuilder_isValidTypedefDeclaration_5809(n.GetOnlyChild(), n.GetJSDocInfo()) {
		this_.DeclareTypedefType(n.GetOnlyChild(), n.GetJSDocInfo())
	}
}

// Java method: AbstractScopeBuilder.defineVarChild(com.google.javascript.rhino.JSDocInfo declarationInfo, com.google.javascript.rhino.Node child, com.google.javascript.jscomp.TypedScope scope)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DefineVarChild(declarationInfo *pkg_rhino.JSDocInfo, child pkg_rhino.Node, scope *TypedScope) {
	_ = declarationInfo
	_ = child
	_ = scope
	if child.IsName() {
		if declarationInfo == nil {
			declarationInfo = child.GetJSDocInfo()
		}
		this_.priv_AbstractScopeBuilder_defineName_5773(child, child.GetFirstChild(), scope, declarationInfo)
	} else {
		if !(child.IsDestructuringLhs()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(child))
		}
		var pattern pkg_rhino.Node = child.GetFirstChild()
		_ = pattern
		var value pkg_rhino.Node = child.GetSecondChild()
		_ = value
		if ModuleImportResolver_IsGoogModuleDependencyCall(value) {
			var defaultImport ScopedName = this_.outer_.moduleImportResolver.GetClosureNamespaceTypeFromCall(value)
			_ = defaultImport
			{
				var key pkg_rhino.Node = pattern.GetFirstChild()
				_ = key
				for ; key != 0; key = key.GetNext() {
					this_.priv_AbstractScopeBuilder_defineModuleImport_5774(key.GetFirstChild(), defaultImport, key.GetString(), scope)
				}
			}
			return
		}
		this_.DefineDestructuringPatternInVarDeclaration(pattern, scope, j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
			return func() *TypedScopeCreator_RValueInfo {
				if value != 0 {
					return NewTypedScopeCreator_RValueInfo(this_.priv_AbstractScopeBuilder_getDeclaredRValueType_5797(pkg_rhino.Node(0), value), value.GetQualifiedNameObject())
				}
				return NewTypedScopeCreator_RValueInfo(pkg_rhino.Upcast_JSType(this_.outer_.unknownType), nil)
			}()
		}))
	}
}

// Java method: AbstractScopeBuilder.inferTypeForDestructuredTarget(com.google.javascript.jscomp.DestructuredTarget target, com.google.common.base.Supplier<com.google.javascript.jscomp.TypedScopeCreator.RValueInfo> patternTypeSupplier)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_inferTypeForDestructuredTarget_5768(target *DestructuredTarget, patternTypeSupplier j2gort.Ext_com_google_common_base_Supplier) *TypedScopeCreator_RValueInfo {
	_ = target
	_ = patternTypeSupplier
	if !(target.HasStringKey()) || target.HasDefaultValue() {
		return TypedScopeCreator_RValueInfo_Empty()
	}
	var rvalue *TypedScopeCreator_RValueInfo = j2gort.CastPtr[TypedScopeCreator_RValueInfo](patternTypeSupplier.Get())
	_ = rvalue
	var patternType pkg_rhino.JSType = rvalue.type_
	_ = patternType
	var propertyName string = target.GetStringKey().GetString()
	_ = propertyName
	var qname pkg_rhino.QualifiedName
	_ = qname
	if rvalue.qualifiedName != nil && !rvalue.qualifiedName.IsJavaNil_() {
		qname = rvalue.qualifiedName.GetQualifiedNameBase().Getprop(propertyName)
	} else {
		qname = nil
	}
	if (patternType == nil || patternType.IsJavaNil_()) || patternType.GetJSTypeBase().IsUnknownType() {
		return NewTypedScopeCreator_RValueInfo(nil, qname)
	}
	base_patternType_ := patternType.GetJSTypeBase()
	if base_patternType_.HasProperty_1_string(propertyName) {
		var type_ pkg_rhino.JSType = base_patternType_.FindPropertyType_1_string(propertyName)
		_ = type_
		return NewTypedScopeCreator_RValueInfo(type_, qname)
	}
	return NewTypedScopeCreator_RValueInfo(nil, qname)
}

// Java method: AbstractScopeBuilder.defineDestructuringPatternInVarDeclaration(com.google.javascript.rhino.Node pattern, com.google.javascript.jscomp.TypedScope scope, com.google.common.base.Supplier<com.google.javascript.jscomp.TypedScopeCreator.RValueInfo> patternTypeSupplier)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DefineDestructuringPatternInVarDeclaration(pattern pkg_rhino.Node, scope *TypedScope, patternTypeSupplier j2gort.Ext_com_google_common_base_Supplier) {
	_ = pattern
	_ = scope
	_ = patternTypeSupplier
	{
		cur_80 := j2gort.NewIterCursor(j2gort.CheckNotNull(DestructuredTarget_CreateAllNonEmptyTargetsInPattern_3_JSTypeRegistry_Ext_com_google_common_base_Supplier_Node(this_.outer_.typeRegistry, j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
			return j2gort.CastPtr[TypedScopeCreator_RValueInfo](patternTypeSupplier.Get()).type_
		}), pattern)))
		for cur_80.HasNext() {
			rawElem_81 := cur_80.Next()
			target := j2gort.CastPtr[DestructuredTarget](rawElem_81)
			_ = target
			var typeSupplier j2gort.Ext_com_google_common_base_Supplier = j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
				return j2gort.PtrToAny(this_.priv_AbstractScopeBuilder_inferTypeForDestructuredTarget_5768(target, patternTypeSupplier))
			})
			_ = typeSupplier
			if target.GetNode().IsDestructuringPattern() {
				this_.DefineDestructuringPatternInVarDeclaration(target.GetNode(), scope, typeSupplier)
			} else {
				var name pkg_rhino.Node = target.GetNode()
				_ = name
				if !(name.IsName()) {
					j2gort.CheckStateFmt(false, "This method is only for declaring variables: %s", j2gort.BoxNode(name))
				}
				var type_ pkg_rhino.JSType = this_.GetDeclaredType(name.GetJSDocInfo(), name, pkg_rhino.Node(0), typeSupplier)
				_ = type_
				if type_ == nil || type_.IsJavaNil_() {
					if name.IsFromExterns() {
						type_ = pkg_rhino.Upcast_JSType(this_.outer_.unknownType)
					} else {
						type_ = nil
					}
				}
				NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).ForDeclarationNode(name).ForVariableName(name.GetString()).InScope(scope).WithType(type_).AllowLaterTypeInference((type_ == nil || type_.IsJavaNil_())).DefineSlot()
			}
		}
	}
}

// Java method: AbstractScopeBuilder.defineClassLiteral(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DefineClassLiteral(n pkg_rhino.Node) {
	_ = n
	this_.AssertDefinitionNode(n, pkg_rhino.StaticGet_Token_CLASS())
	var lValue pkg_rhino.Node = NodeUtil_GetBestLValue(n)
	_ = lValue
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
	_ = info
	var className string = NodeUtil_GetBestLValueName(lValue)
	_ = className
	var classTypeIdentifier string = this_.priv_AbstractScopeBuilder_getBestTypeName_5771(lValue, className)
	_ = classTypeIdentifier
	var classType *pkg_rhino.FunctionType = this_.priv_AbstractScopeBuilder_createClassTypeFromNodes_5776(n, classTypeIdentifier, className, info, lValue)
	_ = classType
	this_.outer_.SetDeferredType(n, pkg_rhino.ToIface_JSType(classType))
	if NodeUtil_IsClassDeclaration(n) {
		j2gort.CheckNotNullString(className)
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).ForDeclarationNode(n.GetFirstChild()).ForVariableName(className).InScope(this_.currentScope).WithType(pkg_rhino.ToIface_JSType(classType)).AllowLaterTypeInference((classType == nil)).DefineSlot()
	}
}

// Java method: AbstractScopeBuilder.getBestTypeName(com.google.javascript.rhino.Node lvalue, java.lang.@org.jspecify.annotations.Nullable String syntacticLvalueName)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_getBestTypeName_5771(lvalue pkg_rhino.Node, syntacticLvalueName string) string {
	_ = lvalue
	_ = syntacticLvalueName
	if syntacticLvalueName == j2gort.NullString {
		return j2gort.NullString
	}
	if this_.IsGoogModuleExports(lvalue) {
		return javalang.StringReplace(syntacticLvalueName, "exports", j2gort.CheckNotNullString(this_.GetModule().ClosureNamespace()))
	}
	return syntacticLvalueName
}

// Java method: AbstractScopeBuilder.defineFunctionLiteral(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DefineFunctionLiteral(n pkg_rhino.Node) {
	_ = n
	this_.AssertDefinitionNode(n, pkg_rhino.StaticGet_Token_FUNCTION())
	var lValue pkg_rhino.Node = NodeUtil_GetBestLValue(n)
	_ = lValue
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
	_ = info
	var functionName string = NodeUtil_GetBestLValueName(lValue)
	_ = functionName
	var functionType *pkg_rhino.FunctionType = this_.priv_AbstractScopeBuilder_createFunctionTypeFromNodes_5780(n, functionName, info, lValue)
	_ = functionType
	this_.outer_.SetDeferredType(n, pkg_rhino.ToIface_JSType(functionType))
	if NodeUtil_IsFunctionDeclaration(n) {
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).ForDeclarationNode(n.GetFirstChild()).ForVariableName(functionName).InScope(this_.currentScope).WithType(pkg_rhino.ToIface_JSType(functionType)).AllowLaterTypeInference((functionType == nil)).DefineSlot()
	}
}

// Java method: AbstractScopeBuilder.defineName(com.google.javascript.rhino.Node name, com.google.javascript.rhino.Node value, com.google.javascript.jscomp.TypedScope scope, com.google.javascript.rhino.JSDocInfo info)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_defineName_5773(name pkg_rhino.Node, value pkg_rhino.Node, scope *TypedScope, info *pkg_rhino.JSDocInfo) {
	_ = name
	_ = value
	_ = scope
	_ = info
	if ModuleImportResolver_IsGoogModuleDependencyCall(value) {
		this_.priv_AbstractScopeBuilder_defineModuleImport_5774(name, this_.outer_.moduleImportResolver.GetClosureNamespaceTypeFromCall(value), j2gort.NullString, scope)
		return
	}
	var type_ pkg_rhino.JSType = this_.GetDeclaredType(info, name, value, nil)
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		if name.IsFromExterns() {
			type_ = pkg_rhino.Upcast_JSType(this_.outer_.unknownType)
		} else {
			type_ = nil
		}
	}
	NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).ForDeclarationNode(name).ForVariableName(name.GetString()).InScope(scope).WithType(type_).AllowLaterTypeInference((type_ == nil || type_.IsJavaNil_())).DefineSlot()
}

// Java method: AbstractScopeBuilder.defineModuleImport(com.google.javascript.rhino.Node localNameNode, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable ScopedName importedModuleObject, java.lang.@org.jspecify.annotations.Nullable String optionalProperty, com.google.javascript.jscomp.TypedScope scopeToDeclareIn)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_defineModuleImport_5774(localNameNode pkg_rhino.Node, importedModuleObject ScopedName, optionalProperty string, scopeToDeclareIn *TypedScope) {
	_ = localNameNode
	_ = importedModuleObject
	_ = optionalProperty
	_ = scopeToDeclareIn
	if importedModuleObject == nil || importedModuleObject.IsJavaNil_() {
		var fallbackType pkg_rhino.JSType
		_ = fallbackType
		if this_.outer_.compiler.GetOptions().AssumeForwardDeclaredForMissingTypes() {
			fallbackType = pkg_rhino.ToIface_JSType(this_.outer_.typeRegistry.CreateNamedType(pkg_rhino.ToIface_StaticTypedScope(this_.currentScope), localNameNode.GetString(), localNameNode.GetSourceFileName(), localNameNode.GetLineno(), localNameNode.GetCharno()))
		} else {
			fallbackType = pkg_rhino.Upcast_JSType(this_.outer_.unknownType)
		}
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).ForDeclarationNode(localNameNode).ForVariableName(localNameNode.GetString()).InScope(scopeToDeclareIn).WithType(fallbackType).AllowLaterTypeInference((fallbackType.GetJSTypeBase().IsUnknownType() && !(this_.outer_.compiler.GetOptions().AssumeForwardDeclaredForMissingTypes()))).DefineSlot()
		return
	}
	base_importedModuleObject_ := importedModuleObject.GetScopedNameBase()
	var exportedScopeRoot pkg_rhino.Node = base_importedModuleObject_.GetScopeRoot()
	_ = exportedScopeRoot
	var exportedQname string
	_ = exportedQname
	if optionalProperty == j2gort.NullString {
		exportedQname = base_importedModuleObject_.GetName()
	} else {
		exportedQname = (j2gort.StrValString(base_importedModuleObject_.GetName()) + "." + j2gort.StrValString(optionalProperty))
	}
	var exportScope *TypedScope
	_ = exportScope
	if exportedScopeRoot != 0 {
		exportScope = this_.outer_.memoized.Get(exportedScopeRoot)
	} else {
		exportScope = nil
	}
	if exportScope != nil {
		var type_ pkg_rhino.JSType = exportScope.LookupQualifiedName_1_string(exportedQname)
		_ = type_
		if type_ != nil && !type_.IsJavaNil_() {
			this_.priv_AbstractScopeBuilder_declareAliasTypeIfRvalueIsAliasable_5791(localNameNode.GetQualifiedName(), localNameNode, exportedQname, type_, exportScope, this_.currentScope)
			NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).ForDeclarationNode(localNameNode).ForVariableName(localNameNode.GetString()).InScope(scopeToDeclareIn).WithType(type_).AllowLaterTypeInference((type_ == nil || type_.IsJavaNil_())).DefineSlot()
			return
		}
	}
	var exportedName ScopedName
	_ = exportedName
	if optionalProperty == j2gort.NullString {
		exportedName = importedModuleObject
	} else {
		exportedName = ScopedName_Of(exportedQname, exportedScopeRoot)
	}
	this_.outer_.weakImports = append(this_.outer_.weakImports, NewTypedScopeCreator_WeakModuleImport(this_.outer_, localNameNode, exportedName, scopeToDeclareIn))
}

// Java method: AbstractScopeBuilder.shouldUseFunctionLiteralType(com.google.javascript.rhino.jstype.FunctionType type, com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.Node lValue)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_shouldUseFunctionLiteralType_5775(type_ *pkg_rhino.FunctionType, info *pkg_rhino.JSDocInfo, lValue pkg_rhino.Node) bool {
	_ = type_
	_ = info
	_ = lValue
	if info != nil {
		return true
	}
	if (lValue != 0) && NodeUtil_MayBeObjectLitKey(lValue) {
		return false
	}
	return (this_.priv_AbstractScopeBuilder_isLValueRootedInGlobalScope_5787(lValue) || !(type_.IsReturnTypeInferred()))
}

// Java method: AbstractScopeBuilder.createClassTypeFromNodes(com.google.javascript.rhino.Node clazz, java.lang.@org.jspecify.annotations.Nullable String name, java.lang.@org.jspecify.annotations.Nullable String syntacticName, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node lvalueNode)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_createClassTypeFromNodes_5776(clazz pkg_rhino.Node, name string, syntacticName string, info *pkg_rhino.JSDocInfo, lvalueNode pkg_rhino.Node) *pkg_rhino.FunctionType {
	_ = clazz
	_ = name
	_ = syntacticName
	_ = info
	_ = lvalueNode
	if !(clazz.IsClass()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(clazz))
	}
	var classType *pkg_rhino.FunctionType = nil
	_ = classType
	var constructor pkg_rhino.Node = pkg_rhino.Node(0)
	_ = constructor
	var builder *FunctionTypeBuilder = this_.outer_.priv_TypedScopeCreator_acquireFunctionTypeBuilder_5720(name, clazz, this_.currentScope)
	_ = builder
	func() {
		defer func() {
			this_.outer_.priv_TypedScopeCreator_releaseFunctionTypeBuilder_5721(builder)
		}()
		builder.UsingClassSyntax().SetSyntacticFunctionName(syntacticName).SetAstContents(clazz, false).SetDeclarationScope(func() *TypedScope {
			if lvalueNode != 0 {
				return this_.priv_AbstractScopeBuilder_getLValueRootScope_5788(lvalueNode)
			}
			return this_.currentScope
		}()).InferKind(info).InferTemplateTypeName(info, nil)
		var extendsClause pkg_rhino.Node = clazz.GetSecondChild()
		_ = extendsClause
		var baseType pkg_rhino.ObjectType = this_.priv_AbstractScopeBuilder_findSuperClassFromNodes_5777(extendsClause, info)
		_ = baseType
		builder.InferInheritance(info, baseType)
		constructor = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(clazz)
		if constructor != 0 {
			constructor = constructor.GetOnlyChild()
		}
		if constructor != 0 {
			var ctorInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(constructor)
			_ = ctorInfo
			builder.InferConstructorParameters(constructor.GetSecondChild(), ctorInfo)
		} else if extendsClause.IsEmpty() {
			builder.InferImplicitConstructorParameters(pkg_rhino.StaticGet_FunctionType_EMPTY_PARAMETERS())
		} else {
			var extendsCtor *pkg_rhino.FunctionType
			_ = extendsCtor
			if baseType != nil && !baseType.IsJavaNil_() {
				extendsCtor = baseType.GetConstructor()
			} else {
				extendsCtor = nil
			}
			if extendsCtor != nil {
				builder.InferImplicitConstructorParameters(extendsCtor.GetParameterArray())
			} else {
				builder.InferImplicitConstructorParameters(this_.outer_.typeRegistry.GetDefaultVarArgsParameterArray())
			}
		}
		classType = builder.BuildAndRegister()
	}()
	if classType.IsConstructor() {
		var qmarkCtor *pkg_rhino.FunctionType = classType.ForgetParameterAndReturnTypes()
		_ = qmarkCtor
		var classPrototypeType pkg_rhino.ObjectType = classType.GetPrototypeProperty()
		_ = classPrototypeType
		classPrototypeType.GetObjectTypeBase().DefineDeclaredProperty_3_string_JSType_Node("constructor", pkg_rhino.ToIface_JSType(qmarkCtor), constructor)
	}
	if classType.HasInstanceType() {
		var classPrototype *pkg_rhino.Property = classType.GetSlot_1_string("prototype")
		_ = classPrototype
		classPrototype.SetNode(func() pkg_rhino.Node {
			if lvalueNode != 0 {
				return lvalueNode
			}
			return classPrototype.GetNode()
		}())
	}
	return classType
}

// Java method: AbstractScopeBuilder.findSuperClassFromNodes(com.google.javascript.rhino.Node extendsNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_findSuperClassFromNodes_5777(extendsNode pkg_rhino.Node, info *pkg_rhino.JSDocInfo) pkg_rhino.ObjectType {
	_ = extendsNode
	_ = info
	if extendsNode.IsEmpty() {
		return nil
	}
	var ctorType pkg_rhino.JSType = extendsNode.GetJSType()
	_ = ctorType
	if ctorType == nil || ctorType.IsJavaNil_() {
		if extendsNode.IsQualifiedName() {
			var superclass string = extendsNode.GetQualifiedName()
			_ = superclass
			ctorType = this_.currentScope.LookupQualifiedName_1_string(superclass)
			if ctorType == nil || ctorType.IsJavaNil_() {
				var var_ *TypedVar = this_.currentScope.LookupTypedVar(superclass)
				_ = var_
				if var_ != nil {
					ctorType = var_.GetType()
				} else {
					ctorType = nil
				}
			}
			if ctorType == nil || ctorType.IsJavaNil_() {
				return pkg_rhino.ObjectType_Cast(this_.outer_.typeRegistry.GetType_5(pkg_rhino.ToIface_StaticTypedScope(this_.currentScope), superclass, extendsNode.GetSourceFileName(), extendsNode.GetLineno(), extendsNode.GetCharno()))
			}
		} else if this_.priv_AbstractScopeBuilder_isGoogModuleGetProperty_5778(extendsNode) {
			ctorType = this_.priv_AbstractScopeBuilder_getCtorForGoogModuleGet_5779(extendsNode)
		} else {
			if (info == nil) || !(info.HasBaseType()) {
				this_.outer_.priv_TypedScopeCreator_report_5724(JSError_Make_3(extendsNode, StaticGet_TypedScopeCreator_DYNAMIC_EXTENDS_WITHOUT_JSDOC(), []string{}))
			}
		}
	}
	if ctorType != nil && !ctorType.IsJavaNil_() {
		if ctorType.GetJSTypeBase().IsConstructor() || ctorType.GetJSTypeBase().IsInterface() {
			return ctorType.ToMaybeFunctionType().GetInstanceType()
		} else if ctorType.GetJSTypeBase().IsUnknownType() || ctorType.GetJSTypeBase().IsNoResolvedType() {
			return ctorType.GetJSTypeBase().ToMaybeObjectType()
		}
	}
	return nil
}

// Java method: AbstractScopeBuilder.isGoogModuleGetProperty(com.google.javascript.rhino.Node extendsClause)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_isGoogModuleGetProperty_5778(extendsClause pkg_rhino.Node) bool {
	_ = extendsClause
	if extendsClause.IsCall() {
		return ModuleImportResolver_IsGoogModuleDependencyCall(extendsClause)
	} else if extendsClause.IsGetProp() {
		return this_.priv_AbstractScopeBuilder_isGoogModuleGetProperty_5778(extendsClause.GetFirstChild())
	} else {
		return false
	}
}

// Java method: AbstractScopeBuilder.getCtorForGoogModuleGet(com.google.javascript.rhino.Node extendsNode)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_getCtorForGoogModuleGet_5779(extendsNode pkg_rhino.Node) pkg_rhino.JSType {
	_ = extendsNode
	var call pkg_rhino.Node = extendsNode
	_ = call
	var properties []string = []string{}
	_ = properties
	for !(call.IsCall()) {
		properties = j2gort.SliceInsertAt(properties, int32(0), call.GetString())
		call = call.GetFirstChild()
	}
	var extendsCall ScopedName = this_.outer_.moduleImportResolver.GetClosureNamespaceTypeFromCall(call)
	_ = extendsCall
	if (extendsCall == nil || extendsCall.IsJavaNil_()) || !(this_.outer_.memoized.ContainsKey(extendsCall.GetScopedNameBase().GetScopeRoot())) {
		return nil
	}
	base_extendsCall_ := extendsCall.GetScopedNameBase()
	var moduleScope *TypedScope = this_.outer_.memoized.Get(base_extendsCall_.GetScopeRoot())
	_ = moduleScope
	properties = j2gort.SliceInsertAt(properties, int32(0), base_extendsCall_.GetName())
	return moduleScope.LookupQualifiedName_1_string(j2gort.Ext_com_google_common_base_Joiner_On_1_uint16(uint16(46)).Join_1_Ext_java_lang_Iterable(j2gort.SliceReadOnlyIterable(properties)))
}

// Java method: AbstractScopeBuilder.createFunctionTypeFromNodes(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rValue, java.lang.@org.jspecify.annotations.Nullable String name, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node lvalueNode)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_createFunctionTypeFromNodes_5780(rValue pkg_rhino.Node, name string, info *pkg_rhino.JSDocInfo, lvalueNode pkg_rhino.Node) *pkg_rhino.FunctionType {
	_ = rValue
	_ = name
	_ = info
	_ = lvalueNode
	if ((rValue != 0) && rValue.IsQualifiedName()) && (lvalueNode != 0) {
		var var_ *TypedVar = this_.currentScope.LookupTypedVar(rValue.GetQualifiedName())
		_ = var_
		if ((var_ != nil) && !(j2gort.IfaceEquals[pkg_rhino.JSType](var_.GetType(), nil))) && var_.GetType().GetJSTypeBase().IsFunctionType() {
			var aliasedType *pkg_rhino.FunctionType = var_.GetType().ToMaybeFunctionType()
			_ = aliasedType
			if aliasedType.IsConstructor() || aliasedType.IsInterface() {
				if name != j2gort.NullString {
					this_.outer_.typeRegistry.DeclareType(pkg_rhino.ToIface_StaticScope(this_.currentScope), name, pkg_rhino.Upcast_JSType(aliasedType.GetInstanceType()))
				}
				this_.priv_AbstractScopeBuilder_checkFunctionAliasAnnotations_5781(lvalueNode, aliasedType, info)
				return aliasedType
			}
		}
	}
	if (info != nil) && info.HasType() {
		var type_ pkg_rhino.JSType = info.GetType().Evaluate(pkg_rhino.ToIface_StaticTypedScope(this_.currentScope), this_.outer_.typeRegistry)
		_ = type_
		type_ = type_.RestrictByNotNullOrUndefined()
		if type_.GetJSTypeBase().IsFunctionType() {
			var functionType *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
			_ = functionType
			functionType.SetJSDocInfo(info)
			return functionType
		}
	}
	var errorRoot pkg_rhino.Node
	_ = errorRoot
	if rValue == 0 {
		errorRoot = lvalueNode
	} else {
		errorRoot = rValue
	}
	var isFnLiteral bool = ((rValue != 0) && rValue.IsFunction())
	_ = isFnLiteral
	var fnRoot pkg_rhino.Node
	_ = fnRoot
	if isFnLiteral {
		fnRoot = rValue
	} else {
		fnRoot = pkg_rhino.Node(0)
	}
	var parametersNode pkg_rhino.Node
	_ = parametersNode
	if isFnLiteral {
		parametersNode = rValue.GetSecondChild()
	} else {
		parametersNode = pkg_rhino.Node(0)
	}
	var classRoot pkg_rhino.Node
	_ = classRoot
	if (lvalueNode != 0) && lvalueNode.GetParent().IsClassMembers() {
		classRoot = lvalueNode.GetGrandparent()
	} else {
		classRoot = pkg_rhino.Node(0)
	}
	var ownerNode pkg_rhino.Node = NodeUtil_GetBestLValueOwner(lvalueNode)
	_ = ownerNode
	var ownerType pkg_rhino.ObjectType = nil
	_ = ownerType
	var propName string = j2gort.NullString
	_ = propName
	if classRoot != 0 {
		ownerType = pkg_rhino.ToIface_ObjectType(pkg_rhino.JSType_ToMaybeFunctionType(classRoot.GetJSType()))
		if !(lvalueNode.IsStaticMember()) && (ownerType != nil && !ownerType.IsJavaNil_()) {
			ownerType = (j2gort.CastPtr[pkg_rhino.FunctionType](ownerType)).GetPrototype()
		}
		if lvalueNode.IsComputedProp() {
			propName = j2gort.NullString
		} else {
			propName = lvalueNode.GetString()
		}
	} else {
		var ownerName string = j2gort.NullString
		_ = ownerName
		if ((lvalueNode != 0) && lvalueNode.IsGetProp()) && (name != j2gort.NullString) {
			var lastDot int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(46)))
			_ = lastDot
			if lastDot >= int32(0) {
				ownerName = j2gort.JStrSub(name, int32(0), lastDot)
			} else {
				ownerName = NodeUtil_GetBestLValueName(ownerNode)
			}
			propName = lvalueNode.GetString()
		} else {
			ownerName = NodeUtil_GetBestLValueName(ownerNode)
			if (ownerName != j2gort.NullString) && (name != j2gort.NullString) {
				if !(javalang.StringStartsWith(j2gort.CheckNotNullString(name), j2gort.CheckNotNullString(ownerName))) {
					j2gort.CheckStateFmt(false, "Expected \"%s\" to start with \"%s\"", j2gort.BoxString(name), j2gort.BoxString(ownerName))
				}
				propName = j2gort.JStrSubFrom(name, (j2gort.JStrLen(ownerName) + int32(1)))
			}
		}
		if ownerName != j2gort.NullString {
			ownerType = pkg_rhino.ObjectType_Cast(this_.currentScope.LookupQualifiedName_1_string(ownerName))
		}
	}
	var prototypeOwner pkg_rhino.ObjectType = this_.priv_AbstractScopeBuilder_getPrototypeOwnerType_5782(ownerType)
	_ = prototypeOwner
	var prototypeOwnerTypeMap *pkg_rhino.TemplateTypeMap = nil
	_ = prototypeOwnerTypeMap
	if (prototypeOwner != nil && !prototypeOwner.IsJavaNil_()) && !(j2gort.IfaceEquals[pkg_rhino.JSType](prototypeOwner.GetTypeOfThis(), nil)) {
		prototypeOwnerTypeMap = prototypeOwner.GetTypeOfThis().GetTemplateTypeMap()
	}
	var overriddenType *pkg_rhino.FunctionType = nil
	_ = overriddenType
	if (ownerType != nil && !ownerType.IsJavaNil_()) && (propName != j2gort.NullString) {
		var overriddenPropType pkg_rhino.JSType = this_.priv_AbstractScopeBuilder_findOverriddenProperty_5783(ownerType, propName, prototypeOwnerTypeMap)
		_ = overriddenPropType
		if overriddenPropType != nil && !overriddenPropType.IsJavaNil_() {
			if lvalueNode.IsGetterDef() {
				overriddenType = this_.outer_.typeRegistry.CreateFunctionType_2_JSType_JSType(overriddenPropType, []pkg_rhino.JSType{})
			} else if lvalueNode.IsSetterDef() {
				overriddenType = this_.outer_.typeRegistry.CreateFunctionType_2_JSType_JSType(this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE()), []pkg_rhino.JSType{overriddenPropType})
			} else if overriddenPropType.GetJSTypeBase().IsFunctionType() {
				overriddenType = overriddenPropType.ToMaybeFunctionType()
			}
		}
	}
	var hasNonEmptyReturns bool = ((fnRoot != 0) && this_.outer_.functionsWithNonEmptyReturns.Contains(fnRoot))
	_ = hasNonEmptyReturns
	var bestTypeName string = this_.priv_AbstractScopeBuilder_getBestTypeName_5771(lvalueNode, name)
	_ = bestTypeName
	var fnType *pkg_rhino.FunctionType = nil
	_ = fnType
	var builder *FunctionTypeBuilder = this_.outer_.priv_TypedScopeCreator_acquireFunctionTypeBuilder_5720(bestTypeName, errorRoot, this_.currentScope)
	_ = builder
	func() {
		defer func() {
			this_.outer_.priv_TypedScopeCreator_releaseFunctionTypeBuilder_5721(builder)
		}()
		builder.SetSyntacticFunctionName(name).SetAstContents(fnRoot, hasNonEmptyReturns).SetDeclarationScope(func() *TypedScope {
			if lvalueNode != 0 {
				return this_.priv_AbstractScopeBuilder_getLValueRootScope_5788(lvalueNode)
			}
			return this_.currentScope
		}()).InferFromOverriddenFunction(overriddenType, parametersNode).InferKind(info).InferClosurePrimitive(info).InferTemplateTypeName(info, pkg_rhino.Upcast_JSType(prototypeOwner)).InferInheritance(info, nil)
		if (info == nil) || !(info.HasReturnType()) {
			if ((rValue != 0) && rValue.IsFunction()) && rValue.HasChildren() {
				var nameDocInfo *pkg_rhino.JSDocInfo = rValue.GetFirstChild().GetJSDocInfo()
				_ = nameDocInfo
				builder.InferReturnType(nameDocInfo, true)
			}
		} else {
			builder.InferReturnType(info, false)
		}
		var fallbackReceiverType pkg_rhino.JSType = nil
		_ = fallbackReceiverType
		if ((ownerType != nil && !ownerType.IsJavaNil_()) && ownerType.GetObjectTypeBase().IsFunctionPrototypeType()) && ownerType.GetOwnerFunction().HasInstanceType() {
			fallbackReceiverType = pkg_rhino.Upcast_JSType(ownerType.GetOwnerFunction().GetInstanceType())
		} else if ((((ownerType != nil && !ownerType.IsJavaNil_()) && ownerType.GetObjectTypeBase().IsFunctionType()) && ownerType.ToMaybeFunctionType().HasInstanceType()) && (lvalueNode != 0)) && lvalueNode.IsStaticMember() {
			fallbackReceiverType = pkg_rhino.Upcast_JSType(ownerType)
		} else if (ownerNode != 0) && ownerNode.IsThis() {
			fallbackReceiverType = this_.currentScope.GetTypeOfThis()
		}
		fnType = builder.InferThisType_2(info, fallbackReceiverType).InferParameterTypes_2(parametersNode, info).BuildAndRegister()
	}()
	if fnType.HasInstanceType() && (lvalueNode != 0) {
		var prototypeSlot *pkg_rhino.Property = fnType.GetSlot_1_string("prototype")
		_ = prototypeSlot
		prototypeSlot.SetNode(lvalueNode)
	}
	return fnType
}

// Java method: AbstractScopeBuilder.checkFunctionAliasAnnotations(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType type, com.google.javascript.rhino.JSDocInfo info)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_checkFunctionAliasAnnotations_5781(n pkg_rhino.Node, type_ *pkg_rhino.FunctionType, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = type_
	_ = info
	if info == nil {
		return
	}
	var annotation string = j2gort.NullString
	_ = annotation
	if info.UsesImplicitMatch() {
		if !(type_.IsStructuralInterface()) {
			annotation = "@record"
		}
	} else if info.IsInterface() {
		if !(type_.IsInterface()) {
			annotation = "@interface"
		}
	} else if info.IsConstructor() && !(type_.IsConstructor()) {
		annotation = "@constructor"
	}
	if (annotation != j2gort.NullString) && (!(n.IsFromExterns()) || (j2gort.CheckNotNullString(annotation) == "@record")) {
		this_.outer_.priv_TypedScopeCreator_report_5724(JSError_Make_3(n, StaticGet_TypedScopeCreator_INCOMPATIBLE_ALIAS_ANNOTATION(), []string{annotation, n.GetQualifiedName()}))
	}
}

// Java method: AbstractScopeBuilder.getPrototypeOwnerType(com.google.javascript.rhino.jstype.ObjectType ownerType)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_getPrototypeOwnerType_5782(ownerType pkg_rhino.ObjectType) pkg_rhino.ObjectType {
	_ = ownerType
	if (ownerType != nil && !ownerType.IsJavaNil_()) && ownerType.GetObjectTypeBase().IsFunctionPrototypeType() {
		return pkg_rhino.ToIface_ObjectType(ownerType.GetOwnerFunction())
	}
	return nil
}

// Java method: AbstractScopeBuilder.findOverriddenProperty(com.google.javascript.rhino.jstype.ObjectType ownerType, java.lang.String propName, com.google.javascript.rhino.jstype.TemplateTypeMap typeMap)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_findOverriddenProperty_5783(ownerType pkg_rhino.ObjectType, propName string, typeMap *pkg_rhino.TemplateTypeMap) pkg_rhino.JSType {
	_ = ownerType
	_ = propName
	_ = typeMap
	var result pkg_rhino.JSType = nil
	_ = result
	var propType pkg_rhino.JSType = ownerType.GetObjectTypeBase().GetPropertyType_1_string(propName)
	_ = propType
	if (propType != nil && !propType.IsJavaNil_()) && !(propType.GetJSTypeBase().IsUnknownType()) {
		result = propType
	} else {
		{
			cur_82 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](ownerType.GetCtorImplementedInterfaces())))
			for cur_82.HasNext() {
				rawElem_83 := cur_82.Next()
				iface := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_83)
				_ = iface
				propType = iface.GetObjectTypeBase().GetPropertyType_1_string(propName)
				if (propType != nil && !propType.IsJavaNil_()) && !(propType.GetJSTypeBase().IsUnknownType()) {
					result = propType
					break
				}
			}
		}
	}
	if ((result != nil && !result.IsJavaNil_()) && (typeMap != nil)) && !(typeMap.IsEmpty()) {
		result = this_.outer_.typeRegistry.ReplaceTemplateTypesPartial(result, typeMap)
	}
	return result
}

// Java method: AbstractScopeBuilder.createEnumTypeFromNodes(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rValue, java.lang.@org.jspecify.annotations.Nullable String name, com.google.javascript.rhino.Node lValue, com.google.javascript.rhino.JSDocInfo info)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_createEnumTypeFromNodes_5784(rValue pkg_rhino.Node, name string, lValue pkg_rhino.Node, info *pkg_rhino.JSDocInfo) *pkg_rhino.EnumType {
	_ = rValue
	_ = name
	_ = lValue
	_ = info
	j2gort.CheckNotNullPtr(info)
	if !(info.HasEnumParameterType()) {
		j2gort.CheckState(false)
	}
	if !((lValue != 0) || (rValue != 0)) {
		j2gort.CheckStateStr(false, "An enum initializer should come from either an lvalue or rvalue")
	}
	var enumType *pkg_rhino.EnumType = nil
	_ = enumType
	if (rValue != 0) && rValue.IsQualifiedName() {
		var var_ *TypedVar = this_.currentScope.LookupTypedVar(rValue.GetQualifiedName())
		_ = var_
		if ((var_ != nil) && !(j2gort.IfaceEquals[pkg_rhino.JSType](var_.GetType(), nil))) && var_.GetType().GetJSTypeBase().IsEnumType() {
			enumType = var_.GetType().ToMaybeEnumType()
		}
	}
	if enumType == nil {
		enumType = j2gort.CastPtr[pkg_rhino.EnumType](j2gort.CastPtr[pkg_rhino.EnumType_Builder](pkg_rhino.EnumType_S_Builder(this_.outer_.typeRegistry).SetName(this_.priv_AbstractScopeBuilder_getBestTypeName_5771(lValue, name))).SetGoogModuleId(TypedScopeCreator_ContainingGoogModuleIdOf(this_.currentScope)).SetSource(rValue).SetElementType(info.GetEnumParameterType().Evaluate(pkg_rhino.ToIface_StaticTypedScope(this_.currentScope), this_.outer_.typeRegistry)).Build())
		if (rValue != 0) && rValue.IsObjectLit() {
			var key pkg_rhino.Node = rValue.GetFirstChild()
			_ = key
			for key != 0 {
				if key.IsComputedProp() {
					this_.outer_.priv_TypedScopeCreator_report_5724(JSError_Make_3(key, StaticGet_TypedScopeCreator_INVALID_ENUM_KEY(), []string{}))
					key = key.GetNext()
					continue
				}
				var keyName string = key.GetString()
				_ = keyName
				j2gort.CheckNotNullString(keyName)
				enumType.DefineElement(keyName, key)
				key = key.GetNext()
			}
		}
	}
	if name != j2gort.NullString {
		this_.outer_.typeRegistry.DeclareType(pkg_rhino.ToIface_StaticScope(this_.currentScope), name, pkg_rhino.ToIface_JSType(enumType.GetElementsType()))
	}
	if (rValue == 0) || !(rValue.IsObjectLit() || rValue.IsQualifiedName()) {
		this_.outer_.priv_TypedScopeCreator_report_5724(JSError_Make_3(func() pkg_rhino.Node {
			if lValue != 0 {
				return lValue
			}
			return rValue
		}(), StaticGet_TypedScopeCreator_ENUM_INITIALIZER(), []string{}))
	}
	return enumType
}

// Java method: AbstractScopeBuilder.declare(com.google.javascript.jscomp.TypedScope scope, java.lang.String name, com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.jscomp.CompilerInput input, boolean inferred)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_declare_5785(scope *TypedScope, name string, n pkg_rhino.Node, type_ pkg_rhino.JSType, input *CompilerInput, inferred bool) *TypedVar {
	_ = scope
	_ = name
	_ = n
	_ = type_
	_ = input
	_ = inferred
	var var_ *TypedVar = scope.Declare(name, n, type_, input, inferred)
	_ = var_
	var varState int32 = this_.outer_.scopedVarStates.LookupDeclareAndDecrementAssigned(scope.GetRootNode_0(), name)
	_ = varState
	if (varState & int32(1)) != int32(0) {
		var_.MarkEscaped()
	}
	if (varState & int32(2)) != int32(0) {
		var_.MarkAssignedExactlyOnce()
	}
	return var_
}

// Java method: AbstractScopeBuilder.finishConstructorDefinition(com.google.javascript.rhino.Node declarationNode, java.lang.String variableName, com.google.javascript.rhino.jstype.FunctionType fnType, com.google.javascript.jscomp.TypedScope scopeToDeclareIn, com.google.javascript.jscomp.CompilerInput input)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_finishConstructorDefinition_5786(declarationNode pkg_rhino.Node, variableName string, fnType *pkg_rhino.FunctionType, scopeToDeclareIn *TypedScope, input *CompilerInput) {
	_ = declarationNode
	_ = variableName
	_ = fnType
	_ = scopeToDeclareIn
	_ = input
	var superClassCtor *pkg_rhino.FunctionType = fnType.GetSuperClassConstructor()
	_ = superClassCtor
	var prototypeSlot *pkg_rhino.Property = fnType.GetSlot_1_string("prototype")
	_ = prototypeSlot
	var prototypeName string = (j2gort.StrValString(variableName) + ".prototype")
	_ = prototypeName
	var prototypeVar *TypedVar = scopeToDeclareIn.LookupTypedVar(prototypeName)
	_ = prototypeVar
	if (prototypeVar != nil) && (j2gort.CastPtr[TypedScope](prototypeVar.GetScope()) == scopeToDeclareIn) {
		scopeToDeclareIn.Undeclare(ToIface_AbstractVar(prototypeVar))
	}
	scopeToDeclareIn.Declare(prototypeName, declarationNode, prototypeSlot.GetType(), input, ((superClassCtor == nil) || superClassCtor.GetInstanceType().GetObjectTypeBase().GetJSTypeBase().EqualsJSType(this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE()))))
}

// Java method: AbstractScopeBuilder.isLValueRootedInGlobalScope(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_isLValueRootedInGlobalScope_5787(n pkg_rhino.Node) bool {
	_ = n
	return this_.priv_AbstractScopeBuilder_getLValueRootScope_5788(n).IsGlobal()
}

// Java method: AbstractScopeBuilder.getLValueRootScope(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_getLValueRootScope_5788(n pkg_rhino.Node) *TypedScope {
	_ = n
	var root pkg_rhino.Node = NodeUtil_GetBestLValueRoot(n)
	_ = root
	if root != 0 {
		if root.IsName() {
			var nameParent pkg_rhino.Node = root.GetParent()
			_ = nameParent
			switch nameParent.TokenSwitchOrd_() {
			case int32(84) /* Token.VAR */ :
				return this_.currentHoistScope
			case int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(105) /* Token.CLASS */, int32(73) /* Token.FUNCTION */, int32(49) /* Token.PARAM_LIST */, int32(86) /* Token.CATCH */ :
				return this_.currentScope
			case int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */ :
				if !(nameParent.GetParent().IsParamList()) {
					j2gort.CheckStateMsg(false, j2gort.BoxNode(nameParent))
				}
				return this_.currentScope
			default:
				if this_.IsGoogModuleExports(root) {
					return this_.currentScope
				}
				var var_ *TypedVar = this_.currentScope.LookupTypedVar(root.GetString())
				_ = var_
				if var_ != nil {
					return j2gort.CastPtr[TypedScope](var_.GetScope())
				}
			}
		} else if root.IsThis() || root.IsSuper() {
			return j2gort.CastPtr[TypedScope](this_.currentHoistScope.GetScopeOfThis())
		}
	}
	return j2gort.CastPtr[TypedScope](this_.currentHoistScope.GetGlobalScope())
}

// Java method: AbstractScopeBuilder.getDeclaredType(com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.Node lValue, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rValue, com.google.common.base.@org.jspecify.annotations.Nullable Supplier<com.google.javascript.jscomp.TypedScopeCreator.RValueInfo> declaredRValueTypeSupplier)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) GetDeclaredType(info *pkg_rhino.JSDocInfo, lValue pkg_rhino.Node, rValue pkg_rhino.Node, declaredRValueTypeSupplier j2gort.Ext_com_google_common_base_Supplier) pkg_rhino.JSType {
	_ = info
	_ = lValue
	_ = rValue
	_ = declaredRValueTypeSupplier
	if (info != nil) && info.HasType() {
		var type_ pkg_rhino.JSType = this_.priv_AbstractScopeBuilder_getDeclaredTypeInAnnotation_5761(lValue, info)
		_ = type_
		if (((type_.GetJSTypeBase().IsSymbol() && lValue.IsFromExterns()) && lValue.IsGetProp()) && lValue.GetFirstChild().MatchesName_1_string("Symbol")) && this_.currentScope.IsGlobal() {
			return pkg_rhino.NewKnownSymbolType(this_.outer_.typeRegistry, lValue.GetQualifiedName())
		}
		return type_
	} else if ((rValue != 0) && rValue.IsFunction()) && this_.priv_AbstractScopeBuilder_shouldUseFunctionLiteralType_5775(pkg_rhino.JSType_ToMaybeFunctionType(rValue.GetJSType()), info, lValue) {
		return rValue.GetJSType()
	} else if (rValue != 0) && rValue.IsClass() {
		return rValue.GetJSType()
	} else if info != nil {
		if info.HasEnumParameterType() {
			if (rValue != 0) && rValue.IsObjectLit() {
				return rValue.GetJSType()
			} else {
				return pkg_rhino.ToIface_JSType(this_.priv_AbstractScopeBuilder_createEnumTypeFromNodes_5784(rValue, lValue.GetQualifiedName(), lValue, info))
			}
		} else if info.IsConstructorOrInterface() {
			var fnType *pkg_rhino.FunctionType = this_.priv_AbstractScopeBuilder_createFunctionTypeFromNodes_5780(rValue, lValue.GetQualifiedName(), info, lValue)
			_ = fnType
			if (rValue == 0) && !(lValue.IsFromExterns()) {
				this_.outer_.priv_TypedScopeCreator_report_5724(JSError_Make_3(lValue, func() *DiagnosticType {
					if fnType.IsConstructor() {
						return StaticGet_TypedScopeCreator_CTOR_INITIALIZER()
					}
					return StaticGet_TypedScopeCreator_IFACE_INITIALIZER()
				}(), []string{lValue.GetQualifiedName()}))
			}
			return pkg_rhino.ToIface_JSType(fnType)
		}
	}
	if ((lValue.IsExport() && lValue.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT())) || NodeUtil_IsConstantDeclaration(info, lValue)) || this_.IsGoogModuleExports(lValue) {
		if rValue != 0 {
			var rValueType pkg_rhino.JSType = this_.priv_AbstractScopeBuilder_getDeclaredRValueType_5797(lValue, rValue)
			_ = rValueType
			if rValue.IsQualifiedName() {
				this_.priv_AbstractScopeBuilder_declareAliasTypeIfRvalueIsAliasable_5791(lValue.GetQualifiedName(), lValue, rValue.GetQualifiedName(), rValueType, this_.currentScope, this_.currentScope)
			}
			if rValueType != nil && !rValueType.IsJavaNil_() {
				return rValueType
			}
		} else if declaredRValueTypeSupplier != nil {
			var rvalueInfo *TypedScopeCreator_RValueInfo = j2gort.CastPtr[TypedScopeCreator_RValueInfo](declaredRValueTypeSupplier.Get())
			_ = rvalueInfo
			if rvalueInfo != nil {
				if rvalueInfo.qualifiedName != nil && !rvalueInfo.qualifiedName.IsJavaNil_() {
					this_.priv_AbstractScopeBuilder_declareAliasTypeIfRvalueIsAliasable_5790(lValue, rvalueInfo.qualifiedName, rvalueInfo.type_, this_.currentScope)
				}
				if rvalueInfo.type_ != nil && !rvalueInfo.type_.IsJavaNil_() {
					return rvalueInfo.type_
				}
			}
		}
	}
	if (rValue != 0) && rValue.IsAssign() {
		return this_.GetDeclaredType(info, lValue, rValue.GetSecondChild(), nil)
	}
	if (info != nil) && FunctionTypeBuilder_IsFunctionTypeDeclaration(info) {
		var fnName string = lValue.GetQualifiedName()
		_ = fnName
		return pkg_rhino.ToIface_JSType(this_.priv_AbstractScopeBuilder_createFunctionTypeFromNodes_5780(pkg_rhino.Node(0), fnName, info, lValue))
	}
	if this_.priv_AbstractScopeBuilder_isValidTypedefDeclaration_5809(lValue, info) {
		return this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_NO_TYPE())
	}
	return nil
}

// Java method: AbstractScopeBuilder.declareAliasTypeIfRvalueIsAliasable(com.google.javascript.rhino.Node lValue, com.google.javascript.rhino.@org.jspecify.annotations.Nullable QualifiedName rValue, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType rValueType, com.google.javascript.jscomp.TypedScope rValueLookupScope)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_declareAliasTypeIfRvalueIsAliasable_5790(lValue pkg_rhino.Node, rValue pkg_rhino.QualifiedName, rValueType pkg_rhino.JSType, rValueLookupScope *TypedScope) {
	_ = lValue
	_ = rValue
	_ = rValueType
	_ = rValueLookupScope
	this_.priv_AbstractScopeBuilder_declareAliasTypeIfRvalueIsAliasable_5792(lValue.GetQualifiedName(), lValue, rValue, rValueType, rValueLookupScope, this_.currentScope)
}

// Java method: AbstractScopeBuilder.declareAliasTypeIfRvalueIsAliasable(java.lang.@org.jspecify.annotations.Nullable String lValueName, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node actualLvalueNode, java.lang.@org.jspecify.annotations.Nullable String rValueName, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType rValueType, com.google.javascript.jscomp.TypedScope rValueLookupScope, com.google.javascript.jscomp.TypedScope aliasDeclarationScope)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_declareAliasTypeIfRvalueIsAliasable_5791(lValueName string, actualLvalueNode pkg_rhino.Node, rValueName string, rValueType pkg_rhino.JSType, rValueLookupScope *TypedScope, aliasDeclarationScope *TypedScope) {
	_ = lValueName
	_ = actualLvalueNode
	_ = rValueName
	_ = rValueType
	_ = rValueLookupScope
	_ = aliasDeclarationScope
	if rValueName == j2gort.NullString {
		return
	}
	var definitionNode pkg_rhino.Node = this_.priv_AbstractScopeBuilder_getDefinitionNode_5795(rValueName, rValueLookupScope)
	_ = definitionNode
	this_.priv_AbstractScopeBuilder_finishDeclareAliasType_5793(lValueName, actualLvalueNode, definitionNode, rValueType, aliasDeclarationScope)
}

// Java method: AbstractScopeBuilder.declareAliasTypeIfRvalueIsAliasable(java.lang.@org.jspecify.annotations.Nullable String lValueName, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node actualLvalueNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable QualifiedName rValue, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType rValueType, com.google.javascript.jscomp.TypedScope rValueLookupScope, com.google.javascript.jscomp.TypedScope aliasDeclarationScope)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_declareAliasTypeIfRvalueIsAliasable_5792(lValueName string, actualLvalueNode pkg_rhino.Node, rValue pkg_rhino.QualifiedName, rValueType pkg_rhino.JSType, rValueLookupScope *TypedScope, aliasDeclarationScope *TypedScope) {
	_ = lValueName
	_ = actualLvalueNode
	_ = rValue
	_ = rValueType
	_ = rValueLookupScope
	_ = aliasDeclarationScope
	if rValue == nil || rValue.IsJavaNil_() {
		return
	}
	var definitionNode pkg_rhino.Node = this_.priv_AbstractScopeBuilder_getDefinitionNode_5796(rValue, rValueLookupScope)
	_ = definitionNode
	this_.priv_AbstractScopeBuilder_finishDeclareAliasType_5793(lValueName, actualLvalueNode, definitionNode, rValueType, aliasDeclarationScope)
}

// Java method: AbstractScopeBuilder.finishDeclareAliasType(java.lang.@org.jspecify.annotations.Nullable String lValueName, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node actualLvalueNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node definitionNode, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType rValueType, com.google.javascript.jscomp.TypedScope aliasDeclarationScope)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_finishDeclareAliasType_5793(lValueName string, actualLvalueNode pkg_rhino.Node, definitionNode pkg_rhino.Node, rValueType pkg_rhino.JSType, aliasDeclarationScope *TypedScope) {
	_ = lValueName
	_ = actualLvalueNode
	_ = definitionNode
	_ = rValueType
	_ = aliasDeclarationScope
	if definitionNode != 0 {
		var typedefType pkg_rhino.JSType = definitionNode.GetTypedefTypeProp()
		_ = typedefType
		if typedefType != nil && !typedefType.IsJavaNil_() {
			actualLvalueNode.SetTypedefTypeProp(typedefType)
			if lValueName != j2gort.NullString {
				this_.outer_.typeRegistry.IdentifyNonNullableName(pkg_rhino.ToIface_StaticScope(aliasDeclarationScope), lValueName)
				this_.outer_.typeRegistry.DeclareType(pkg_rhino.ToIface_StaticScope(aliasDeclarationScope), lValueName, typedefType)
			}
			return
		}
	}
	if lValueName == j2gort.NullString {
		return
	}
	if ((rValueType != nil && !rValueType.IsJavaNil_()) && rValueType.GetJSTypeBase().IsFunctionType()) && rValueType.ToMaybeFunctionType().HasInstanceType() {
		var functionType *pkg_rhino.FunctionType = rValueType.ToMaybeFunctionType()
		_ = functionType
		this_.outer_.typeRegistry.DeclareType(pkg_rhino.ToIface_StaticScope(aliasDeclarationScope), lValueName, pkg_rhino.Upcast_JSType(functionType.GetInstanceType()))
		return
	}
	if (rValueType != nil && !rValueType.IsJavaNil_()) && rValueType.GetJSTypeBase().IsEnumType() {
		this_.outer_.typeRegistry.DeclareType(pkg_rhino.ToIface_StaticScope(aliasDeclarationScope), lValueName, pkg_rhino.ToIface_JSType(rValueType.ToMaybeEnumType().GetElementsType()))
		this_.outer_.typeRegistry.IdentifyNonNullableName(pkg_rhino.ToIface_StaticScope(aliasDeclarationScope), lValueName)
	}
}

// Java method: AbstractScopeBuilder.isGoogModuleExports(com.google.javascript.rhino.Node lValue)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) IsGoogModuleExports(lValue pkg_rhino.Node) bool {
	_ = lValue
	if (this_.GetModule() == nil) || (lValue == 0) {
		return false
	}
	if this_.outer_.undeclaredNamesForClosure.Contains(lValue) {
		return true
	}
	return ((((lValue.IsStringKey() && lValue.GetParent().IsObjectLit()) && lValue.GetGrandparent().IsAssign()) && lValue.GetParent().GetPrevious().MatchesName_1_string("exports")) && this_.outer_.undeclaredNamesForClosure.Contains(lValue.GetParent().GetPrevious()))
}

// Java method: AbstractScopeBuilder.getDefinitionNode(java.lang.String qname, com.google.javascript.jscomp.TypedScope scope)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_getDefinitionNode_5795(qname string, scope *TypedScope) pkg_rhino.Node {
	_ = qname
	_ = scope
	var lastDot int32 = j2gort.JStrLastIndexOfCp(qname, int32(uint16(46)))
	_ = lastDot
	if lastDot < int32(0) {
		var var_ *TypedVar = scope.LookupTypedVar(qname)
		_ = var_
		if var_ != nil {
			return var_.GetNameNode()
		} else {
			return pkg_rhino.Node(0)
		}
	}
	var parent pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(scope.LookupQualifiedName_1_string(j2gort.JStrSub(qname, int32(0), lastDot)))
	_ = parent
	if parent != nil && !parent.IsJavaNil_() {
		return parent.GetObjectTypeBase().GetPropertyDefSite_1_string(j2gort.JStrSubFrom(qname, (lastDot + int32(1))))
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: AbstractScopeBuilder.getDefinitionNode(com.google.javascript.rhino.QualifiedName qname, com.google.javascript.jscomp.TypedScope scope)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_getDefinitionNode_5796(qname pkg_rhino.QualifiedName, scope *TypedScope) pkg_rhino.Node {
	_ = qname
	_ = scope
	if qname.IsSimple() {
		var var_ *TypedVar = scope.LookupTypedVar(qname.GetComponent())
		_ = var_
		if var_ != nil {
			return var_.GetNameNode()
		} else {
			return pkg_rhino.Node(0)
		}
	}
	var parent pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(scope.LookupQualifiedName_1_QualifiedName(qname.GetOwner()))
	_ = parent
	if parent != nil && !parent.IsJavaNil_() {
		return parent.GetObjectTypeBase().GetPropertyDefSite_1_string(qname.GetComponent())
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: AbstractScopeBuilder.getDeclaredRValueType(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node lValue, com.google.javascript.rhino.Node rValue)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_getDeclaredRValueType_5797(lValue pkg_rhino.Node, rValue pkg_rhino.Node) pkg_rhino.JSType {
	_ = lValue
	_ = rValue
	var rValueInfo *pkg_rhino.JSDocInfo = rValue.GetJSDocInfo()
	_ = rValueInfo
	if (rValue.IsCast() && (rValueInfo != nil)) && rValueInfo.HasType() {
		return rValueInfo.GetType().Evaluate(pkg_rhino.ToIface_StaticTypedScope(this_.currentScope), this_.outer_.typeRegistry)
	}
	var type_ pkg_rhino.JSType = rValue.GetJSType()
	_ = type_
	if (type_ != nil && !type_.IsJavaNil_()) && !(type_.GetJSTypeBase().IsUnknownType()) {
		return type_
	}
	if rValue.IsQualifiedName() {
		return this_.currentScope.LookupQualifiedNameNode(rValue)
	}
	if NodeUtil_IsBooleanResult(rValue) {
		return this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE())
	}
	if NodeUtil_IsNumericResult(rValue) {
		return this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
	}
	if NodeUtil_IsBigIntResult(rValue) {
		return this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_BIGINT_TYPE())
	}
	if NodeUtil_IsStringResult(rValue) {
		return this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
	}
	if rValue.IsNew() && rValue.GetFirstChild().IsQualifiedName() {
		var targetType pkg_rhino.JSType = this_.currentScope.LookupQualifiedNameNode(rValue.GetFirstChild())
		_ = targetType
		if targetType != nil && !targetType.IsJavaNil_() {
			var fnType *pkg_rhino.FunctionType = targetType.RestrictByNotNullOrUndefined().ToMaybeFunctionType()
			_ = fnType
			if (fnType != nil) && fnType.HasInstanceType() {
				return pkg_rhino.Upcast_JSType(fnType.GetInstanceType())
			}
		}
	}
	if rValue.IsCall() && rValue.GetFirstChild().MatchesName_1_string("Symbol") {
		var symbolVar *TypedVar = j2gort.CastPtr[TypedVar](this_.currentScope.GetVar("Symbol"))
		_ = symbolVar
		if (symbolVar != nil) && j2gort.CastPtr[TypedScope](symbolVar.GetScope()).IsGlobal() {
			return this_.outer_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_SYMBOL_TYPE())
		}
	}
	if rValue.IsOr() {
		var firstClause pkg_rhino.Node = rValue.GetFirstChild()
		_ = firstClause
		var secondClause pkg_rhino.Node = firstClause.GetNext()
		_ = secondClause
		var namesMatch bool = (((firstClause.IsName() && (lValue != 0)) && lValue.IsName()) && func() bool {
			sr_84 := firstClause.GetString()
			se_85 := lValue.GetString()
			return j2gort.CheckNotNullString(sr_84) == se_85
		}())
		_ = namesMatch
		if namesMatch {
			type_ = secondClause.GetJSType()
			if (type_ != nil && !type_.IsJavaNil_()) && !(type_.GetJSTypeBase().IsUnknownType()) {
				return type_
			}
		}
	}
	return nil
}

// Java method: AbstractScopeBuilder.checkForClassDefiningCalls(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) CheckForClassDefiningCalls(n pkg_rhino.Node) {
	_ = n
	var relationship *CodingConvention_SubclassRelationship = this_.outer_.codingConvention.GetClassesDefinedByCall(n)
	_ = relationship
	if relationship != nil {
		var superClass pkg_rhino.ObjectType = TypeValidator_GetInstanceOfCtor(this_.currentScope.LookupQualifiedName_1_string(relationship.superclassName))
		_ = superClass
		var subClass pkg_rhino.ObjectType = TypeValidator_GetInstanceOfCtor(this_.currentScope.LookupQualifiedName_1_string(relationship.subclassName))
		_ = subClass
		if (superClass != nil && !superClass.IsJavaNil_()) && (subClass != nil && !subClass.IsJavaNil_()) {
			var superCtor *pkg_rhino.FunctionType = superClass.GetConstructor()
			_ = superCtor
			var subCtor *pkg_rhino.FunctionType = subClass.GetConstructor()
			_ = subCtor
			if (superCtor != nil) && (subCtor != nil) {
				this_.outer_.codingConvention.ApplySubclassRelationship(pkg_rhino.NewNominalTypeBuilder(superCtor, superClass), pkg_rhino.NewNominalTypeBuilder(subCtor, subClass), relationship.type_)
			}
		}
	}
	var singletonGetterClassName string = this_.outer_.codingConvention.GetSingletonGetterClassName(n)
	_ = singletonGetterClassName
	if singletonGetterClassName != j2gort.NullString {
		var objectType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(this_.outer_.typeRegistry.GetType_2(pkg_rhino.ToIface_StaticScope(this_.currentScope), singletonGetterClassName))
		_ = objectType
		if objectType != nil && !objectType.IsJavaNil_() {
			var functionType *pkg_rhino.FunctionType = objectType.GetConstructor()
			_ = functionType
			if functionType != nil {
				var getterType *pkg_rhino.FunctionType = this_.outer_.typeRegistry.CreateFunctionType_2_JSType_JSType(pkg_rhino.Upcast_JSType(objectType), []pkg_rhino.JSType{})
				_ = getterType
				this_.outer_.codingConvention.ApplySingletonGetter(pkg_rhino.NewNominalTypeBuilder(functionType, objectType), getterType)
			}
		}
	}
	var objectLiteralCast *CodingConvention_ObjectLiteralCast = this_.outer_.codingConvention.GetObjectLiteralCast(n)
	_ = objectLiteralCast
	if objectLiteralCast != nil {
		if objectLiteralCast.diagnosticType == nil {
			var type_ pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(this_.outer_.typeRegistry.GetType_2(pkg_rhino.ToIface_StaticScope(this_.currentScope), objectLiteralCast.typeName))
			_ = type_
			if (type_ != nil && !type_.IsJavaNil_()) && (type_.GetConstructor() != nil) {
				this_.outer_.SetDeferredType(objectLiteralCast.objectNode, pkg_rhino.Upcast_JSType(type_))
				objectLiteralCast.objectNode.PutBooleanProp(pkg_rhino.StaticGet_Node_REFLECTED_OBJECT(), true)
			} else {
				this_.outer_.priv_TypedScopeCreator_report_5724(JSError_Make_3(n, StaticGet_TypedScopeCreator_CONSTRUCTOR_EXPECTED(), []string{}))
			}
		} else {
			this_.outer_.priv_TypedScopeCreator_report_5724(JSError_Make_3(n, objectLiteralCast.diagnosticType, []string{}))
		}
	}
}

// Java method: AbstractScopeBuilder.maybeDeclareGetElem(com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rhsValue)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) MaybeDeclareGetElem(info *pkg_rhino.JSDocInfo, n pkg_rhino.Node, rhsValue pkg_rhino.Node) {
	_ = info
	_ = n
	_ = rhsValue
	if !(n.IsGetElem()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	var ownerNode pkg_rhino.Node = n.GetFirstChild()
	_ = ownerNode
	var key pkg_rhino.Node = n.GetSecondChild()
	_ = key
	var keyType *pkg_rhino.KnownSymbolType = TypedScopeCreator_ExtractKnownSymbolKey(this_.currentScope, key)
	_ = keyType
	if keyType == nil {
		return
	}
	var valueType pkg_rhino.JSType = this_.GetDeclaredType(info, n, rhsValue, nil)
	_ = valueType
	if valueType == nil || valueType.IsJavaNil_() {
		return
	}
	var ownerName string = ownerNode.GetQualifiedName()
	_ = ownerName
	var ownerType pkg_rhino.ObjectType = this_.priv_AbstractScopeBuilder_getObjectSlot_5806(ownerName)
	_ = ownerType
	if ownerType == nil || ownerType.IsJavaNil_() {
		return
	}
	this_.DeclarePropertyIfNamespaceType_6_ObjectType_Node_Property_Key_JSType_Node_JSDocInfo(ownerType, ownerNode, pkg_rhino.NewProperty_SymbolKey(keyType), valueType, n, nil)
}

// Java method: AbstractScopeBuilder.maybeDeclareQualifiedName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rhsValue)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) MaybeDeclareQualifiedName(t *NodeTraversal, info *pkg_rhino.JSDocInfo, n pkg_rhino.Node, parent pkg_rhino.Node, rhsValue pkg_rhino.Node) {
	_ = t
	_ = info
	_ = n
	_ = parent
	_ = rhsValue
	var isTypedef bool = this_.priv_AbstractScopeBuilder_isValidTypedefDeclaration_5809(n, info)
	_ = isTypedef
	if isTypedef {
		this_.DeclareTypedefType(n, info)
	}
	var ownerNode pkg_rhino.Node = n.GetFirstChild()
	_ = ownerNode
	var propName string = n.GetString()
	_ = propName
	if !(n.IsGetProp() && ownerNode.IsQualifiedName()) {
		j2gort.CheckArgument(false)
	}
	var ownerName string = j2gort.NullString
	_ = ownerName
	var qName string = j2gort.NullString
	_ = qName
	var valueType pkg_rhino.JSType = this_.GetDeclaredType(info, n, rhsValue, nil)
	_ = valueType
	if (valueType == nil || valueType.IsJavaNil_()) && (rhsValue != 0) {
		valueType = rhsValue.GetJSType()
	}
	if "prototype" == propName {
		ownerName = ownerNode.GetQualifiedName()
		qName = (j2gort.StrValString(ownerName) + "." + j2gort.StrValString(propName))
		var qVar *TypedVar = this_.currentScope.LookupTypedVar(qName)
		_ = qVar
		if qVar != nil {
			var qVarType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(qVar.GetType())
			_ = qVarType
			if ((qVarType != nil && !qVarType.IsJavaNil_()) && (rhsValue != 0)) && rhsValue.IsObjectLit() {
				this_.outer_.typeRegistry.ResetImplicitPrototype(rhsValue.GetJSType(), qVarType.GetImplicitPrototype())
			} else if !(qVar.IsTypeInferred()) {
				return
			}
			j2gort.CastPtr[TypedScope](qVar.GetScope()).Undeclare(ToIface_AbstractVar(qVar))
		}
	}
	if valueType == nil || valueType.IsJavaNil_() {
		if parent.IsExprResult() {
			if ownerName == j2gort.NullString {
				ownerName = ownerNode.GetQualifiedName()
			}
			var stubOwnerName string = ownerName
			_ = stubOwnerName
			var isExtern bool = ((t.GetInput() != nil) && t.GetInput().IsExtern())
			_ = isExtern
			this_.DeferAction(this_.currentScope.GetRootNode_0(), j2gort.NewLambda_Ext_java_lang_Runnable(func() {
				this_.ResolveStubDeclaration(n, isExtern, stubOwnerName)
			}))
		}
		return
	}
	var inferred bool = this_.priv_AbstractScopeBuilder_isQualifiedNameInferred_5801(ownerNode, ownerName, qName, propName, n, info, rhsValue, valueType)
	_ = inferred
	if !(inferred) {
		if ownerName == j2gort.NullString {
			ownerName = ownerNode.GetQualifiedName()
		}
		if qName == j2gort.NullString {
			qName = (j2gort.StrValString(ownerName) + "." + j2gort.StrValString(propName))
		}
		var ownerType pkg_rhino.ObjectType = this_.priv_AbstractScopeBuilder_getObjectSlot_5806(ownerName)
		_ = ownerType
		if ownerType != nil && !ownerType.IsJavaNil_() {
			this_.DeclarePropertyIfNamespaceType_6_ObjectType_Node_string_JSType_Node_JSDocInfo(ownerType, ownerNode, propName, valueType, n, nil)
		}
		var declareInScope bool = ((ownerType == nil || ownerType.IsJavaNil_()) || !(ownerType.GetObjectTypeBase().IsFunctionPrototypeType()))
		_ = declareInScope
		if declareInScope {
			NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).ForDeclarationNode(n).ForVariableName(qName).InScope(this_.priv_AbstractScopeBuilder_getLValueRootScope_5788(n)).WithType(valueType).AllowLaterTypeInference(inferred).DefineSlot()
		}
	}
}

// Java method: AbstractScopeBuilder.isQualifiedNameInferred(com.google.javascript.rhino.Node ownerNode, java.lang.@org.jspecify.annotations.Nullable String ownerName, java.lang.@org.jspecify.annotations.Nullable String qName, java.lang.String propName, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rhsValue, com.google.javascript.rhino.jstype.JSType valueType)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_isQualifiedNameInferred_5801(ownerNode pkg_rhino.Node, ownerName string, qName string, propName string, n pkg_rhino.Node, info *pkg_rhino.JSDocInfo, rhsValue pkg_rhino.Node, valueType pkg_rhino.JSType) bool {
	_ = ownerNode
	_ = ownerName
	_ = qName
	_ = propName
	_ = n
	_ = info
	_ = rhsValue
	_ = valueType
	if ("prototype" == propName) || javalang.StringEndsWith(j2gort.CheckNotNullString(propName), ".prototype") {
		if ownerName == j2gort.NullString {
			ownerName = ownerNode.GetQualifiedName()
		}
		var className string = j2gort.NullString
		_ = className
		if ("prototype" == propName) && (j2gort.JStrIndexOfCp(propName, int32(uint16(46))) < int32(0)) {
			className = ownerName
		} else {
			if qName == j2gort.NullString {
				qName = (j2gort.StrValString(ownerName) + "." + j2gort.StrValString(propName))
			}
			className = j2gort.JStrSub(qName, int32(0), j2gort.JStrLastIndexOf(qName, ".prototype"))
		}
		var slot *TypedVar = this_.currentScope.LookupTypedVar(className)
		_ = slot
		var classType pkg_rhino.JSType
		_ = classType
		if slot == nil {
			classType = nil
		} else {
			classType = slot.GetType()
		}
		if (classType != nil && !classType.IsJavaNil_()) && (classType.GetJSTypeBase().IsConstructor() || classType.GetJSTypeBase().IsInterface()) {
			return false
		}
	}
	rhsValue = this_.priv_AbstractScopeBuilder_unwrapIfAssign_5802(rhsValue)
	if ((info != nil) && ((((info.HasType() || info.HasEnumParameterType()) || this_.priv_AbstractScopeBuilder_isValidTypedefDeclaration_5809(n, info)) || FunctionTypeBuilder_IsFunctionTypeDeclaration(info)) || ((rhsValue != 0) && rhsValue.IsFunction()))) || this_.priv_AbstractScopeBuilder_isTypedConstantDeclaration_5804(info, n, valueType) {
		return false
	}
	if (rhsValue == 0) || (!(rhsValue.IsFunction()) && !(rhsValue.IsClass())) {
		return true
	}
	if !(n.IsUnscopedQualifiedName()) {
		return true
	}
	var ownerScope *TypedScope = this_.priv_AbstractScopeBuilder_getLValueRootScope_5788(n)
	_ = ownerScope
	if qName == j2gort.NullString {
		if ownerName == j2gort.NullString {
			ownerName = ownerNode.GetQualifiedName()
		}
		qName = (j2gort.StrValString(ownerName) + "." + j2gort.StrValString(propName))
	}
	if (ownerScope != nil) && ownerScope.HasOwnSlot(qName) {
		return true
	}
	if this_.priv_AbstractScopeBuilder_hasControlStructureAncestor_5805(n.GetParent()) {
		return true
	}
	if (ownerScope != nil) && this_.outer_.scopedVarStates.IsEscaped(ownerScope.GetRootNode_0(), qName) {
		return true
	}
	return false
}

// Java method: AbstractScopeBuilder.unwrapIfAssign(com.google.javascript.rhino.Node maybeAssign)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_unwrapIfAssign_5802(maybeAssign pkg_rhino.Node) pkg_rhino.Node {
	_ = maybeAssign
	if (maybeAssign == 0) || !(maybeAssign.IsAssign()) {
		return maybeAssign
	}
	return this_.priv_AbstractScopeBuilder_unwrapIfAssign_5802(maybeAssign.GetSecondChild())
}

// Java method: AbstractScopeBuilder.declareProvidedNs(com.google.javascript.rhino.Node provideCall, com.google.javascript.jscomp.ProcessClosureProvidesAndRequires.ProvidedName providedName)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DeclareProvidedNs(provideCall pkg_rhino.Node, providedName *ProcessClosureProvidesAndRequires_ProvidedName) {
	_ = provideCall
	_ = providedName
	var anonymousObjectType pkg_rhino.ObjectType = this_.outer_.typeRegistry.CreateAnonymousObjectType(nil)
	_ = anonymousObjectType
	NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).InScope(j2gort.CastPtr[TypedScope](this_.currentScope.GetGlobalScope())).AllowLaterTypeInference(false).ForVariableName(providedName.GetNamespace()).ForDeclarationNode(provideCall).WithType(pkg_rhino.Upcast_JSType(anonymousObjectType)).ForGoogProvidedName().DefineSlot()
	var ns string = providedName.GetNamespace()
	_ = ns
	var lastDot int32 = j2gort.JStrLastIndexOfCp(ns, int32(uint16(46)))
	_ = lastDot
	if lastDot >= int32(0) {
		var ownerType pkg_rhino.JSType = this_.currentScope.LookupQualifiedName_1_string(j2gort.JStrSub(ns, int32(0), lastDot))
		_ = ownerType
		if (ownerType != nil && !ownerType.IsJavaNil_()) && ownerType.GetJSTypeBase().IsObjectType() {
			ownerType.GetJSTypeBase().ToMaybeObjectType().GetObjectTypeBase().DefineDeclaredProperty_3_string_JSType_Node(j2gort.JStrSubFrom(ns, (lastDot+int32(1))), pkg_rhino.Upcast_JSType(anonymousObjectType), provideCall)
		}
	}
}

// Java method: AbstractScopeBuilder.isTypedConstantDeclaration(com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType valueType)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_isTypedConstantDeclaration_5804(info *pkg_rhino.JSDocInfo, n pkg_rhino.Node, valueType pkg_rhino.JSType) bool {
	_ = info
	_ = n
	_ = valueType
	return ((NodeUtil_IsConstantDeclaration(info, n) || this_.IsGoogModuleExports(n)) && (valueType != nil && !valueType.IsJavaNil_()))
}

// Java method: AbstractScopeBuilder.hasControlStructureAncestor(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_hasControlStructureAncestor_5805(n pkg_rhino.Node) bool {
	_ = n
	for !(n.IsScript() || n.IsFunction()) {
		if NodeUtil_IsControlStructure(n) {
			return true
		}
		n = n.GetParent()
	}
	return false
}

// Java method: AbstractScopeBuilder.getObjectSlot(java.lang.String slotName)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_getObjectSlot_5806(slotName string) pkg_rhino.ObjectType {
	_ = slotName
	var ownerVar *TypedVar = this_.currentScope.LookupTypedVar(slotName)
	_ = ownerVar
	if ownerVar != nil {
		var ownerVarType pkg_rhino.JSType = ownerVar.GetType()
		_ = ownerVarType
		return pkg_rhino.ObjectType_Cast(func() pkg_rhino.JSType {
			if ownerVarType == nil || ownerVarType.IsJavaNil_() {
				return nil
			}
			return ownerVarType.RestrictByNotNullOrUndefined()
		}())
	}
	return nil
}

// Java method: AbstractScopeBuilder.getInheritedInterfacePropertyType(com.google.javascript.rhino.jstype.ObjectType obj, java.lang.String propName)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_getInheritedInterfacePropertyType_5807(obj pkg_rhino.ObjectType, propName string) pkg_rhino.JSType {
	_ = obj
	_ = propName
	if (obj != nil && !obj.IsJavaNil_()) && obj.GetObjectTypeBase().IsFunctionPrototypeType() {
		var f *pkg_rhino.FunctionType = obj.GetOwnerFunction()
		_ = f
		{
			cur_86 := j2gort.NewIterCursor(j2gort.CheckNotNull(f.GetImplementedInterfaces()))
			for cur_86.HasNext() {
				rawElem_87 := cur_86.Next()
				i := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_87)
				_ = i
				if i.GetObjectTypeBase().HasProperty_1_string(propName) {
					return i.GetObjectTypeBase().GetPropertyType_1_string(propName)
				}
			}
		}
	}
	return nil
}

// Java method: AbstractScopeBuilder.resolveStubDeclaration(com.google.javascript.rhino.Node n, boolean isExtern, java.lang.String ownerName)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) ResolveStubDeclaration(n pkg_rhino.Node, isExtern bool, ownerName string) {
	_ = n
	_ = isExtern
	_ = ownerName
	var qName string = n.GetQualifiedName()
	_ = qName
	var propName string = n.GetString()
	_ = propName
	if this_.currentScope.HasOwnSlot(qName) {
		return
	}
	var ownerType pkg_rhino.ObjectType = this_.priv_AbstractScopeBuilder_getObjectSlot_5806(ownerName)
	_ = ownerType
	var inheritedType pkg_rhino.JSType = this_.priv_AbstractScopeBuilder_getInheritedInterfacePropertyType_5807(ownerType, propName)
	_ = inheritedType
	var stubType pkg_rhino.JSType
	_ = stubType
	if inheritedType == nil || inheritedType.IsJavaNil_() {
		stubType = pkg_rhino.Upcast_JSType(this_.outer_.unknownType)
	} else {
		stubType = inheritedType
	}
	NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_.self_).ForDeclarationNode(n).ReadVariableNameFromDeclarationNode().InScope(this_.priv_AbstractScopeBuilder_getLValueRootScope_5788(n)).WithType(stubType).AllowLaterTypeInference(true).DefineSlot()
	if (ownerType != nil && !ownerType.IsJavaNil_()) && (isExtern || ownerType.GetObjectTypeBase().IsFunctionPrototypeType()) {
		ownerType.GetObjectTypeBase().DefineInferredProperty_3_string_JSType_Node(propName, stubType, n)
	} else {
		this_.outer_.typeRegistry.RegisterPropertyOnType(propName, func() pkg_rhino.JSType {
			if ownerType == nil || ownerType.IsJavaNil_() {
				return stubType
			}
			return pkg_rhino.Upcast_JSType(ownerType)
		}())
	}
}

// Java method: AbstractScopeBuilder.isValidTypedefDeclaration(com.google.javascript.rhino.Node candidate, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) priv_AbstractScopeBuilder_isValidTypedefDeclaration_5809(candidate pkg_rhino.Node, info *pkg_rhino.JSDocInfo) bool {
	_ = candidate
	_ = info
	if (info == nil) || !(info.HasTypedefType()) {
		return false
	}
	return (candidate.IsUnscopedQualifiedName() && !(NodeUtil_IsPrototypeProperty(candidate)))
}

// Java method: AbstractScopeBuilder.declareTypedefType(com.google.javascript.rhino.Node candidate, com.google.javascript.rhino.JSDocInfo info)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DeclareTypedefType(candidate pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = candidate
	_ = info
	var typedef string = candidate.GetQualifiedName()
	_ = typedef
	this_.outer_.typeRegistry.DeclareType(pkg_rhino.ToIface_StaticScope(this_.currentScope), typedef, pkg_rhino.Upcast_JSType(this_.outer_.unknownType))
	var realType pkg_rhino.JSType = info.GetTypedefType().Evaluate(pkg_rhino.ToIface_StaticTypedScope(this_.currentScope), this_.outer_.typeRegistry)
	_ = realType
	if realType == nil || realType.IsJavaNil_() {
		this_.outer_.priv_TypedScopeCreator_report_5724(JSError_Make_3(candidate, StaticGet_TypedScopeCreator_MALFORMED_TYPEDEF(), []string{typedef}))
	} else {
		candidate.SetTypedefTypeProp(realType)
	}
	this_.outer_.typeRegistry.OverwriteDeclaredType(pkg_rhino.ToIface_StaticScope(this_.currentScope), typedef, realType)
}

// Java method: AbstractScopeBuilder.declarePropertyIfNamespaceType(com.google.javascript.rhino.jstype.ObjectType ownerType, com.google.javascript.rhino.Node ownerNode, java.lang.String propName, com.google.javascript.rhino.jstype.JSType valueType, com.google.javascript.rhino.Node declarationNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo jsdocInfo)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DeclarePropertyIfNamespaceType_6_ObjectType_Node_string_JSType_Node_JSDocInfo(ownerType pkg_rhino.ObjectType, ownerNode pkg_rhino.Node, propName string, valueType pkg_rhino.JSType, declarationNode pkg_rhino.Node, jsdocInfo *pkg_rhino.JSDocInfo) {
	_ = ownerType
	_ = ownerNode
	_ = propName
	_ = valueType
	_ = declarationNode
	_ = jsdocInfo
	j2gort.CheckNotNullString(propName)
	base_ownerType_ := ownerType.GetObjectTypeBase()
	if base_ownerType_.HasOwnProperty_1_string(propName) && !(base_ownerType_.IsPropertyTypeInferred_1_string(propName)) {
		return
	}
	var isNonNativeExtern bool = (((this_.GetCompilerInput() != nil) && this_.GetCompilerInput().IsExtern()) && !(base_ownerType_.IsNativeObjectType()))
	_ = isNonNativeExtern
	if (isNonNativeExtern || !(ownerType.IsInstanceType())) || ownerNode.IsThis() {
		base_ownerType_.DefineDeclaredProperty_3_string_JSType_Node(propName, valueType, declarationNode)
		ownerType.SetPropertyJSDocInfo_2_string_JSDocInfo(propName, jsdocInfo)
	}
}

// Java method: AbstractScopeBuilder.declarePropertyIfNamespaceType(com.google.javascript.rhino.jstype.ObjectType ownerType, com.google.javascript.rhino.Node ownerNode, com.google.javascript.rhino.jstype.Property.Key propName, com.google.javascript.rhino.jstype.JSType valueType, com.google.javascript.rhino.Node declarationNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo jsdocInfo)
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) DeclarePropertyIfNamespaceType_6_ObjectType_Node_Property_Key_JSType_Node_JSDocInfo(ownerType pkg_rhino.ObjectType, ownerNode pkg_rhino.Node, propName pkg_rhino.Property_Key, valueType pkg_rhino.JSType, declarationNode pkg_rhino.Node, jsdocInfo *pkg_rhino.JSDocInfo) {
	_ = ownerType
	_ = ownerNode
	_ = propName
	_ = valueType
	_ = declarationNode
	_ = jsdocInfo
	base_ownerType_ := ownerType.GetObjectTypeBase()
	if base_ownerType_.HasOwnProperty_1_Property_Key(propName) && !(base_ownerType_.IsPropertyTypeInferred_1_Property_Key(propName)) {
		return
	}
	var isNonNativeExtern bool = (((this_.GetCompilerInput() != nil) && this_.GetCompilerInput().IsExtern()) && !(base_ownerType_.IsNativeObjectType()))
	_ = isNonNativeExtern
	if (isNonNativeExtern || !(ownerType.IsInstanceType())) || ownerNode.IsThis() {
		base_ownerType_.DefineDeclaredProperty_3_Property_Key_JSType_Node(propName, valueType, declarationNode)
		ownerType.SetPropertyJSDocInfo_2_Property_Key_JSDocInfo(propName, jsdocInfo)
	}
}

// Java class: com.google.javascript.jscomp.TypedScopeCreator.AbstractScopeBuilder.SlotDefiner
type TypedScopeCreator_AbstractScopeBuilder_SlotDefiner struct {
	outer_                  TypedScopeCreator_AbstractScopeBuilder
	variableName            string           `gson:"variableName"`
	scope                   *TypedScope      `gson:"scope"`
	type_                   pkg_rhino.JSType `gson:"type"`
	declarationNode         pkg_rhino.Node   `gson:"declarationNode"`
	idHash_                 atomic.Int32
	allowLaterTypeInference bool `gson:"allowLaterTypeInference"`
	forGoogProvidedName     bool `gson:"forGoogProvidedName"`
}

func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) IsJavaNil_() bool {
	return this_ == nil
}
func (*TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator$AbstractScopeBuilder$SlotDefiner", "SlotDefiner"
}
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator.AbstractScopeBuilder.SlotDefiner", "com.google.javascript.jscomp.TypedScopeCreator$AbstractScopeBuilder$SlotDefiner", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScopeCreator_AbstractScopeBuilder_SlotDefiner](v_)
	})
}
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) SetOuter_(o TypedScopeCreator_AbstractScopeBuilder) {
	this_.outer_ = o
}
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) GetOuter_() TypedScopeCreator_AbstractScopeBuilder {
	return this_.outer_
}

func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) initFields_TypedScopeCreator_AbstractScopeBuilder_SlotDefiner() {
	_ = this_
	this_.variableName = j2gort.NullString

	this_.allowLaterTypeInference = true

}

// Java constructor: SlotDefiner()
func NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(outer_ TypedScopeCreator_AbstractScopeBuilder) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	this_ := &TypedScopeCreator_AbstractScopeBuilder_SlotDefiner{}
	this_.outer_ = outer_
	InitTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_)
	return this_
}

func InitTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) {
	this_.initFields_TypedScopeCreator_AbstractScopeBuilder_SlotDefiner()
}

// Java method: SlotDefiner.forDeclarationNode(com.google.javascript.rhino.Node declarationNode)
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) ForDeclarationNode(declarationNode pkg_rhino.Node) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	_ = declarationNode
	this_.declarationNode = declarationNode
	return this_
}

// Java method: SlotDefiner.readVariableNameFromDeclarationNode()
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) ReadVariableNameFromDeclarationNode() *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	var parent pkg_rhino.Node = this_.declarationNode.GetParent()
	_ = parent
	if this_.declarationNode.IsName() {
		if !(((((parent.IsFunction() || parent.IsClass()) || NodeUtil_IsNameDeclaration(parent)) || parent.IsParamList()) || (parent.IsRest() && parent.GetParent().IsParamList())) || parent.IsCatch()) {
			j2gort.CheckArgument(false)
		}
	} else {
		if !(this_.declarationNode.IsGetProp() && (parent.IsAssign() || parent.IsExprResult())) {
			j2gort.CheckArgument(false)
		}
	}
	this_.variableName = this_.declarationNode.GetQualifiedName()
	return this_
}

// Java method: SlotDefiner.forVariableName(java.lang.String variableName)
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) ForVariableName(variableName string) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	_ = variableName
	this_.variableName = variableName
	return this_
}

// Java method: SlotDefiner.inScope(com.google.javascript.jscomp.TypedScope scope)
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) InScope(scope *TypedScope) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	_ = scope
	this_.scope = j2gort.CheckNotNullPtr(scope)
	return this_
}

// Java method: SlotDefiner.withType(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type)
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) WithType(type_ pkg_rhino.JSType) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	_ = type_
	this_.type_ = type_
	return this_
}

// Java method: SlotDefiner.allowLaterTypeInference(boolean allowLaterTypeInference)
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) AllowLaterTypeInference(allowLaterTypeInference bool) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	_ = allowLaterTypeInference
	this_.allowLaterTypeInference = allowLaterTypeInference
	return this_
}

// Java method: SlotDefiner.forGoogProvidedName()
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) ForGoogProvidedName() *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	this_.forGoogProvidedName = true
	return this_
}

// Java method: SlotDefiner.defineSlot()
func (this_ *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) DefineSlot() {
	j2gort.CheckNotNullNode(this_.declarationNode)
	if !(((((this_.declarationNode.IsName() || this_.declarationNode.IsGetProp()) || NodeUtil_MayBeObjectLitKey(this_.declarationNode)) || this_.declarationNode.IsModuleBody()) || this_.declarationNode.IsExport()) || (this_.forGoogProvidedName && NodeUtil_IsExprCall(this_.declarationNode))) {
		j2gort.CheckStateFmt(false, "declaration node must be an lvalue or goog.provide call, found %s", j2gort.BoxNode(this_.declarationNode))
	}
	j2gort.CheckNotNullString(this_.variableName)
	if !(this_.allowLaterTypeInference || this_.type_ != nil && !this_.type_.IsJavaNil_()) {
		j2gort.CheckStateStr(false, "null type but inference not allowed")
	}
	if !(!(j2gort.CheckNotNullString(this_.variableName) == "")) {
		j2gort.CheckState(false)
	}
	j2gort.CheckNotNullPtr(this_.scope)
	var parent pkg_rhino.Node = this_.declarationNode.GetParent()
	_ = parent
	var scopeToDeclareIn *TypedScope = this_.scope
	_ = scopeToDeclareIn
	var isGlobalVar bool = (this_.declarationNode.IsName() && scopeToDeclareIn.IsGlobal())
	_ = isGlobalVar
	var shouldDeclareOnGlobalThis bool = ((isGlobalVar && (parent.IsVar() || parent.IsFunction())) || (this_.forGoogProvidedName && !(j2gort.JStrIndexOf(this_.variableName, ".") >= 0)))
	_ = shouldDeclareOnGlobalThis
	if (((scopeToDeclareIn != this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentHoistScope) && scopeToDeclareIn.IsGlobal()) && scopeToDeclareIn.HasOwnSlot(this_.variableName)) && !(this_.forGoogProvidedName) {
		scopeToDeclareIn = this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentHoistScope
	}
	var input *CompilerInput = this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().outer_.compiler.GetInput(this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().inputId)
	_ = input
	if !(scopeToDeclareIn.CanDeclare(this_.variableName)) {
		var oldVar *TypedVar = scopeToDeclareIn.LookupTypedVar(this_.variableName)
		_ = oldVar
		this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().outer_.validator.ExpectUndeclaredVariable(this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().sourceName, input, this_.declarationNode, parent, oldVar, this_.variableName, this_.type_)
	} else {
		if this_.type_ != nil && !this_.type_.IsJavaNil_() {
			this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().outer_.SetDeferredType(this_.declarationNode, this_.type_)
		}
		this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().priv_AbstractScopeBuilder_declare_5785(scopeToDeclareIn, this_.variableName, this_.declarationNode, this_.type_, input, this_.allowLaterTypeInference)
	}
	var fnType *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(this_.type_)
	_ = fnType
	if (fnType != nil) && !(this_.type_.GetJSTypeBase().IsEmptyType()) {
		if fnType.IsConstructor() || fnType.IsInterface() {
			this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().priv_AbstractScopeBuilder_finishConstructorDefinition_5786(this_.declarationNode, this_.variableName, fnType, scopeToDeclareIn, input)
		}
	}
	if shouldDeclareOnGlobalThis {
		var globalThis pkg_rhino.ObjectType = this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().outer_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_GLOBAL_THIS())
		_ = globalThis
		if this_.allowLaterTypeInference {
			globalThis.GetObjectTypeBase().DefineInferredProperty_3_string_JSType_Node(this_.variableName, func() pkg_rhino.JSType {
				if this_.type_ == nil || this_.type_.IsJavaNil_() {
					return this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_NO_TYPE())
				}
				return this_.type_
			}(), this_.declarationNode)
		} else {
			globalThis.GetObjectTypeBase().DefineDeclaredProperty_3_string_JSType_Node(this_.variableName, this_.type_, this_.declarationNode)
		}
	}
	if (((isGlobalVar && func() bool {
		sr_88 := "Window"
		se_89 := this_.variableName
		return j2gort.CheckNotNullString(sr_88) == se_89
	}()) && (this_.type_ != nil && !this_.type_.IsJavaNil_())) && this_.type_.GetJSTypeBase().IsFunctionType()) && this_.type_.GetJSTypeBase().IsConstructor() {
		var globalThisCtor *pkg_rhino.FunctionType = this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().outer_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_GLOBAL_THIS()).GetConstructor()
		_ = globalThisCtor
		globalThisCtor.GetInstanceType().ClearCachedValues()
		this_.outer_.GetTypedScopeCreator_AbstractScopeBuilderBase().outer_.typeRegistry.ResetImplicitPrototype(pkg_rhino.Upcast_JSType(globalThisCtor.GetPrototype()), this_.type_.ToMaybeFunctionType().GetInstanceType())
	}
}

// Java class: com.google.javascript.jscomp.TypedScopeCreator.NormalScopeBuilder
type TypedScopeCreator_NormalScopeBuilder struct {
	TypedScopeCreator_AbstractScopeBuilderBase
	outer_ *TypedScopeCreator
}

func (this_ *TypedScopeCreator_NormalScopeBuilder) BindSelf_(s any) {
	_ = s
	this_.TypedScopeCreator_AbstractScopeBuilderBase.BindSelf_(s)
}
func (this_ *TypedScopeCreator_NormalScopeBuilder) IsJavaNil_() bool { return this_ == nil }
func (*TypedScopeCreator_NormalScopeBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator$NormalScopeBuilder", "NormalScopeBuilder"
}
func (this_ *TypedScopeCreator_NormalScopeBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.TypedScopeCreator_AbstractScopeBuilderBase.IdentityHash_()
}
func (*TypedScopeCreator_NormalScopeBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator.NormalScopeBuilder", "com.google.javascript.jscomp.TypedScopeCreator$NormalScopeBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScopeCreator_NormalScopeBuilder](v_)
	})
}
func (this_ *TypedScopeCreator_NormalScopeBuilder) SetOuter_(o *TypedScopeCreator) { this_.outer_ = o }
func (this_ *TypedScopeCreator_NormalScopeBuilder) GetOuter_() *TypedScopeCreator {
	return this_.outer_
}
func (*TypedScopeCreator_NormalScopeBuilder) Implements_jscomp_TypedScopeCreator_AbstractScopeBuilder() {
}
func (*TypedScopeCreator_NormalScopeBuilder) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *TypedScopeCreator_NormalScopeBuilder) initFields_TypedScopeCreator_NormalScopeBuilder() {
	_ = this_
}

// Java constructor: NormalScopeBuilder(com.google.javascript.jscomp.TypedScope scope)
func NewTypedScopeCreator_NormalScopeBuilder(outer_ *TypedScopeCreator, scope *TypedScope) *TypedScopeCreator_NormalScopeBuilder {
	this_ := &TypedScopeCreator_NormalScopeBuilder{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitTypedScopeCreator_NormalScopeBuilder(this_, scope)
	return this_
}

func InitTypedScopeCreator_NormalScopeBuilder(this_ *TypedScopeCreator_NormalScopeBuilder, scope *TypedScope) {
	_ = scope
	this_.TypedScopeCreator_AbstractScopeBuilderBase.outer_ = this_.outer_
	InitTypedScopeCreator_AbstractScopeBuilder(&this_.TypedScopeCreator_AbstractScopeBuilderBase, scope)
	this_.initFields_TypedScopeCreator_NormalScopeBuilder()
}

// Java method: NormalScopeBuilder.visitPreorder(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedScopeCreator_NormalScopeBuilder) VisitPreorder(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if ((parent != 0) && NodeUtil_CreatesBlockScope(n)) && !(n.IsClass()) {
		this_.outer_.CreateScope_2_Node_AbstractScope(n, ToIface_AbstractScope(this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope))
	}
	if n.IsFunction() {
		this_.DefineFunctionLiteral(n)
	}
}

// Java method: NormalScopeBuilder.visitPostorder(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedScopeCreator_NormalScopeBuilder) VisitPostorder(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */ :
		this_.CheckForClassDefiningCalls(n)
	case int32(51) /* Token.ASSIGN */ :
		var firstChild pkg_rhino.Node = n.GetFirstChild()
		_ = firstChild
		if firstChild.IsGetProp() && firstChild.IsQualifiedName() {
			this_.MaybeDeclareQualifiedName(t, n.GetJSDocInfo(), firstChild, n, firstChild.GetNext())
		} else if firstChild.IsGetElem() && firstChild.GetFirstChild().IsQualifiedName() {
			this_.MaybeDeclareGetElem(n.GetJSDocInfo(), firstChild, firstChild.GetNext())
		} else if this_.outer_.undeclaredNamesForClosure.Contains(firstChild) {
			this_.DefineAssignAsIfVarDeclaration(n)
		}
	case int32(86) /* Token.CATCH */ :
		this_.DefineCatch(n)
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		this_.DefineVars(n)
	case int32(26) /* Token.GETPROP */ :
		if parent.IsExprResult() && n.IsQualifiedName() {
			this_.MaybeDeclareQualifiedName(t, n.GetJSDocInfo(), n, parent, pkg_rhino.Node(0))
		}
	case int32(27) /* Token.GETELEM */ :
		if parent.IsExprResult() && n.GetFirstChild().IsQualifiedName() {
			this_.MaybeDeclareGetElem(n.GetJSDocInfo(), n, pkg_rhino.Node(0))
		}
	case int32(105) /* Token.CLASS */ :
		this_.outer_.CreateScope_2_Node_AbstractScope(n, ToIface_AbstractScope(this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope))
	case int32(93) /* Token.EXPR_RESULT */ :
		if !(this_.outer_.providedNamesFromCall.IsEmpty()) && j2gort.MultimapContainsKeyNode(this_.outer_.providedNamesFromCall, n) {
			{
				cur_90 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.MultimapGetNode(this_.outer_.providedNamesFromCall, n))))
				for cur_90.HasNext() {
					rawElem_91 := cur_90.Next()
					name := j2gort.CastPtr[ProcessClosureProvidesAndRequires_ProvidedName](rawElem_91)
					_ = name
					this_.DeclareProvidedNs(n, name)
				}
			}
		}
	case int32(120) /* Token.EXPORT */ :
		if n.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
			var declaredType pkg_rhino.JSType = this_.GetDeclaredType(nil, n, n.GetOnlyChild(), nil)
			_ = declaredType
			NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_).InScope(this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope).ForDeclarationNode(n).WithType(declaredType).ForVariableName("*default*").AllowLaterTypeInference((declaredType == nil || declaredType.IsJavaNil_())).DefineSlot()
		}
	default:
	}
}

// Java class: com.google.javascript.jscomp.TypedScopeCreator.FunctionScopeBuilder
type TypedScopeCreator_FunctionScopeBuilder struct {
	TypedScopeCreator_AbstractScopeBuilderBase
	outer_ *TypedScopeCreator
}

func (this_ *TypedScopeCreator_FunctionScopeBuilder) BindSelf_(s any) {
	_ = s
	this_.TypedScopeCreator_AbstractScopeBuilderBase.BindSelf_(s)
}
func (this_ *TypedScopeCreator_FunctionScopeBuilder) IsJavaNil_() bool { return this_ == nil }
func (*TypedScopeCreator_FunctionScopeBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator$FunctionScopeBuilder", "FunctionScopeBuilder"
}
func (this_ *TypedScopeCreator_FunctionScopeBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.TypedScopeCreator_AbstractScopeBuilderBase.IdentityHash_()
}
func (*TypedScopeCreator_FunctionScopeBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator.FunctionScopeBuilder", "com.google.javascript.jscomp.TypedScopeCreator$FunctionScopeBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScopeCreator_FunctionScopeBuilder](v_)
	})
}
func (this_ *TypedScopeCreator_FunctionScopeBuilder) SetOuter_(o *TypedScopeCreator) {
	this_.outer_ = o
}
func (this_ *TypedScopeCreator_FunctionScopeBuilder) GetOuter_() *TypedScopeCreator {
	return this_.outer_
}
func (*TypedScopeCreator_FunctionScopeBuilder) Implements_jscomp_TypedScopeCreator_AbstractScopeBuilder() {
}
func (*TypedScopeCreator_FunctionScopeBuilder) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *TypedScopeCreator_FunctionScopeBuilder) initFields_TypedScopeCreator_FunctionScopeBuilder() {
	_ = this_
}

// Java constructor: FunctionScopeBuilder(com.google.javascript.jscomp.TypedScope scope)
func NewTypedScopeCreator_FunctionScopeBuilder(outer_ *TypedScopeCreator, scope *TypedScope) *TypedScopeCreator_FunctionScopeBuilder {
	this_ := &TypedScopeCreator_FunctionScopeBuilder{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitTypedScopeCreator_FunctionScopeBuilder(this_, scope)
	return this_
}

func InitTypedScopeCreator_FunctionScopeBuilder(this_ *TypedScopeCreator_FunctionScopeBuilder, scope *TypedScope) {
	_ = scope
	this_.TypedScopeCreator_AbstractScopeBuilderBase.outer_ = this_.outer_
	InitTypedScopeCreator_AbstractScopeBuilder(&this_.TypedScopeCreator_AbstractScopeBuilderBase, scope)
	this_.initFields_TypedScopeCreator_FunctionScopeBuilder()
}

// Java method: FunctionScopeBuilder.visitPreorder(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedScopeCreator_FunctionScopeBuilder) VisitPreorder(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if parent == 0 {
		this_.HandleFunctionInputs()
	} else if n.IsFunction() {
		this_.DefineFunctionLiteral(n)
	}
}

// Java method: FunctionScopeBuilder.handleFunctionInputs()
func (this_ *TypedScopeCreator_FunctionScopeBuilder) HandleFunctionInputs() {
	var fnNode pkg_rhino.Node = this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope.GetRootNode_0()
	_ = fnNode
	var fnNameNode pkg_rhino.Node = fnNode.GetFirstChild()
	_ = fnNameNode
	var fnName string = fnNameNode.GetString()
	_ = fnName
	if !(j2gort.CheckNotNullString(fnName) == "") && NodeUtil_IsFunctionExpression(fnNode) {
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_).ForDeclarationNode(fnNameNode).ForVariableName(fnName).InScope(this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope).WithType(fnNode.GetJSType()).AllowLaterTypeInference(false).DefineSlot()
	}
	this_.DeclareParameters(fnNode)
}

// Java method: FunctionScopeBuilder.declareParameters(com.google.javascript.rhino.Node functionNode)
func (this_ *TypedScopeCreator_FunctionScopeBuilder) DeclareParameters(functionNode pkg_rhino.Node) {
	_ = functionNode
	if NodeUtil_IsBundledGoogModuleCall(functionNode.GetParent()) {
		return
	}
	var astParameters pkg_rhino.Node = functionNode.GetSecondChild()
	_ = astParameters
	var iifeArgumentNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = iifeArgumentNode
	if NodeUtil_IsInvocationTarget(functionNode) {
		iifeArgumentNode = functionNode.GetNext()
	}
	var functionType *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(functionNode.GetJSType())
	_ = functionType
	if functionType != nil {
		var jsdocParameters []pkg_rhino.FunctionType_Parameter = functionType.GetParameterArray()
		_ = jsdocParameters
		var jsdocParamLen int32 = int32(len(jsdocParameters))
		_ = jsdocParamLen
		var jsdocParamIdx int32 = int32(0)
		_ = jsdocParamIdx
		{
			var astParameter pkg_rhino.Node = astParameters.GetFirstChild()
			_ = astParameter
			for ; astParameter != 0; astParameter = astParameter.GetNext() {
				if (iifeArgumentNode != 0) && iifeArgumentNode.IsSpread() {
					iifeArgumentNode = pkg_rhino.Node(0)
				}
				var declaredType pkg_rhino.JSType
				_ = declaredType
				if jsdocParamIdx < jsdocParamLen {
					declaredType = jsdocParameters[func() int32 { old_92 := jsdocParamIdx; nw_93 := old_92 + 1; jsdocParamIdx = nw_93; return old_92 }()].GetJSType()
				} else {
					declaredType = pkg_rhino.Upcast_JSType(this_.outer_.unknownType)
				}
				this_.priv_FunctionScopeBuilder_declareNamesInPositionalParameter_5824(astParameter, declaredType, iifeArgumentNode)
				if iifeArgumentNode != 0 {
					iifeArgumentNode = iifeArgumentNode.GetNext()
				}
			}
		}
		this_.priv_FunctionScopeBuilder_declareTemplatesInFunctionBody_5823(functionNode, functionType)
	}
}

// Java method: FunctionScopeBuilder.declareTemplatesInFunctionBody(com.google.javascript.rhino.Node functionNode, com.google.javascript.rhino.jstype.FunctionType functionType)
func (this_ *TypedScopeCreator_FunctionScopeBuilder) priv_FunctionScopeBuilder_declareTemplatesInFunctionBody_5823(functionNode pkg_rhino.Node, functionType *pkg_rhino.FunctionType) {
	_ = functionNode
	_ = functionType
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(functionNode)
	_ = info
	var ownerTemplateMap *pkg_rhino.TemplateTypeMap = nil
	_ = ownerTemplateMap
	var isPossiblePrototypeMethod bool = (!(functionType.HasInstanceType()) && ((info == nil) || !(info.HasThisType())))
	_ = isPossiblePrototypeMethod
	if isPossiblePrototypeMethod {
		var ownerNode pkg_rhino.Node = NodeUtil_GetBestLValue(functionNode)
		_ = ownerNode
		ownerTemplateMap = this_.FindOwnerTemplateMap(ownerNode)
	}
	var infoTemplateNames j2gort.Ext_com_google_common_collect_ImmutableList
	_ = infoTemplateNames
	if info != nil {
		infoTemplateNames = info.GetTemplateTypeNames()
	} else {
		infoTemplateNames = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var infoTypeTransformations j2gort.Ext_com_google_common_collect_ImmutableMap
	_ = infoTypeTransformations
	if info != nil {
		infoTypeTransformations = info.GetTypeTransformations()
	} else {
		infoTypeTransformations = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	}
	var ownerKeyCount int32
	_ = ownerKeyCount
	if ownerTemplateMap != nil {
		ownerKeyCount = ownerTemplateMap.Size()
	} else {
		ownerKeyCount = int32(0)
	}
	if ((ownerKeyCount == int32(0)) && infoTemplateNames.IsEmpty()) && infoTypeTransformations.IsEmpty() {
		return
	}
	var templateNames j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = templateNames
	{
		var i int32 = int32(0)
		_ = i
		for ; i < ownerKeyCount; i = i + 1 {
			templateNames.Add(ownerTemplateMap.GetTemplateKeyAt(i).GetReferenceName())
		}
	}
	j2gort.LocalOrderedSetAddAllStr(&templateNames, j2gort.AsIface[j2gort.Ext_java_util_Collection](infoTemplateNames))
	if !(infoTypeTransformations.IsEmpty()) {
		j2gort.LocalOrderedSetAddAllStr(&templateNames, j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](infoTypeTransformations.KeySet_0()))
	}
	var input *CompilerInput = this_.GetCompilerInput()
	_ = input
	var voidType pkg_rhino.JSType = this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
	_ = voidType
	{
		cur_94 := templateNames.Iter()
		for cur_94.HasNext() {
			name := cur_94.
				Next()
			_ = name
			if !(this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope.CanDeclare(name)) {
				this_.outer_.validator.ExpectUndeclaredVariable(NodeUtil_GetSourceName(functionNode), input, functionNode, functionNode.GetParent(), j2gort.CastPtr[TypedVar](this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope.GetVar(name)), name, voidType)
			}
			this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope.Declare(name, functionNode, voidType, input, false)
		}
	}
}

// Java method: FunctionScopeBuilder.declareNamesInPositionalParameter(com.google.javascript.rhino.Node astParameter, com.google.javascript.rhino.jstype.JSType declaredParameterType, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node iifeArgumentNode)
func (this_ *TypedScopeCreator_FunctionScopeBuilder) priv_FunctionScopeBuilder_declareNamesInPositionalParameter_5824(astParameter pkg_rhino.Node, declaredParameterType pkg_rhino.JSType, iifeArgumentNode pkg_rhino.Node) {
	_ = astParameter
	_ = declaredParameterType
	_ = iifeArgumentNode
	var paramType pkg_rhino.JSType = declaredParameterType
	_ = paramType
	var isInferred bool = paramType.GetJSTypeBase().EqualsJSType(this_.outer_.unknownType)
	_ = isInferred
	if (iifeArgumentNode != 0) && isInferred {
		var argumentName string = iifeArgumentNode.GetQualifiedName()
		_ = argumentName
		var argumentVar *TypedVar
		_ = argumentVar
		if (argumentName == j2gort.NullString) || (j2gort.CastPtr[TypedScope](this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope.GetParent()) == nil) {
			argumentVar = nil
		} else {
			argumentVar = j2gort.CastPtr[TypedVar](j2gort.CastPtr[TypedScope](this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope.GetParent()).GetVar(argumentName))
		}
		if (argumentVar != nil) && !(argumentVar.IsTypeInferred()) {
			paramType = argumentVar.GetType()
		}
	}
	if paramType == nil || paramType.IsJavaNil_() {
		paramType = pkg_rhino.Upcast_JSType(this_.outer_.unknownType)
	}
	switch astParameter.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		this_.priv_FunctionScopeBuilder_declareSingleParameterName_5826(isInferred, astParameter, paramType)
	case int32(125) /* Token.ITER_REST */ :
		var param pkg_rhino.Node = astParameter.GetFirstChild()
		_ = param
		var arrayType pkg_rhino.ObjectType = this_.outer_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE())
		_ = arrayType
		var restParamType pkg_rhino.JSType = pkg_rhino.ToIface_JSType(this_.outer_.typeRegistry.CreateTemplatizedType_2_ObjectType_JSType(arrayType, []pkg_rhino.JSType{paramType}))
		_ = restParamType
		if param.IsName() {
			this_.priv_FunctionScopeBuilder_declareSingleParameterName_5826(isInferred, astParameter.GetFirstChild(), restParamType)
		} else {
			this_.priv_FunctionScopeBuilder_declareDestructuringParameter_5825(isInferred, param, restParamType)
		}
	case int32(134) /* Token.DEFAULT_VALUE */ :
		var actualParam pkg_rhino.Node = astParameter.GetFirstChild()
		_ = actualParam
		if actualParam.IsName() {
			this_.priv_FunctionScopeBuilder_declareSingleParameterName_5826(isInferred, actualParam, paramType)
		} else {
			this_.priv_FunctionScopeBuilder_declareDestructuringParameter_5825(isInferred, actualParam, paramType)
		}
	case int32(102) /* Token.ARRAY_PATTERN */, int32(103) /* Token.OBJECT_PATTERN */ :
		this_.priv_FunctionScopeBuilder_declareDestructuringParameter_5825(isInferred, astParameter, paramType)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected function parameter node "+j2gort.StrVal(astParameter))))
	}
}

// Java method: FunctionScopeBuilder.declareDestructuringParameter(boolean isInferred, com.google.javascript.rhino.Node pattern, com.google.javascript.rhino.jstype.JSType patternType)
func (this_ *TypedScopeCreator_FunctionScopeBuilder) priv_FunctionScopeBuilder_declareDestructuringParameter_5825(isInferred bool, pattern pkg_rhino.Node, patternType pkg_rhino.JSType) {
	_ = isInferred
	_ = pattern
	_ = patternType
	{
		cur_95 := j2gort.NewIterCursor(j2gort.CheckNotNull(DestructuredTarget_CreateAllNonEmptyTargetsInPattern_3_JSTypeRegistry_JSType_Node(this_.outer_.typeRegistry, patternType, pattern)))
		for cur_95.HasNext() {
			rawElem_96 := cur_95.Next()
			target := j2gort.CastPtr[DestructuredTarget](rawElem_96)
			_ = target
			var parameterType pkg_rhino.JSType = target.InferTypeWithoutUsingDefaultValue()
			_ = parameterType
			if target.GetNode().IsDestructuringPattern() {
				this_.priv_FunctionScopeBuilder_declareDestructuringParameter_5825(isInferred, target.GetNode(), parameterType)
			} else {
				var paramName pkg_rhino.Node = target.GetNode()
				_ = paramName
				if !(paramName.IsName()) {
					j2gort.CheckStateFmt(false, "Expected all parameters to be names, got %s", j2gort.BoxNode(paramName))
				}
				if (parameterType == nil || parameterType.IsJavaNil_()) || parameterType.GetJSTypeBase().IsUnknownType() {
					var paramJSDoc *pkg_rhino.JSDocInfo = paramName.GetJSDocInfo()
					_ = paramJSDoc
					if (paramJSDoc != nil) && paramJSDoc.HasType() {
						parameterType = this_.outer_.typeRegistry.EvaluateTypeExpression(paramJSDoc.GetType(), pkg_rhino.ToIface_StaticTypedScope(this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope))
						isInferred = false
					} else {
						isInferred = true
					}
				}
				this_.priv_FunctionScopeBuilder_declareSingleParameterName_5826(isInferred, paramName, parameterType)
			}
		}
	}
}

// Java method: FunctionScopeBuilder.declareSingleParameterName(boolean isInferred, com.google.javascript.rhino.Node name, com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypedScopeCreator_FunctionScopeBuilder) priv_FunctionScopeBuilder_declareSingleParameterName_5826(isInferred bool, name pkg_rhino.Node, type_ pkg_rhino.JSType) {
	_ = isInferred
	_ = name
	_ = type_
	NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_).ForDeclarationNode(name).ForVariableName(name.GetString()).InScope(this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope).WithType(type_).AllowLaterTypeInference(isInferred).DefineSlot()
}

// Java class: com.google.javascript.jscomp.TypedScopeCreator.ClassScopeBuilder
type TypedScopeCreator_ClassScopeBuilder struct {
	TypedScopeCreator_AbstractScopeBuilderBase
	outer_ *TypedScopeCreator
}

func (this_ *TypedScopeCreator_ClassScopeBuilder) BindSelf_(s any) {
	_ = s
	this_.TypedScopeCreator_AbstractScopeBuilderBase.BindSelf_(s)
}
func (this_ *TypedScopeCreator_ClassScopeBuilder) IsJavaNil_() bool { return this_ == nil }
func (*TypedScopeCreator_ClassScopeBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator$ClassScopeBuilder", "ClassScopeBuilder"
}
func (this_ *TypedScopeCreator_ClassScopeBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.TypedScopeCreator_AbstractScopeBuilderBase.IdentityHash_()
}
func (*TypedScopeCreator_ClassScopeBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator.ClassScopeBuilder", "com.google.javascript.jscomp.TypedScopeCreator$ClassScopeBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScopeCreator_ClassScopeBuilder](v_)
	})
}
func (this_ *TypedScopeCreator_ClassScopeBuilder) SetOuter_(o *TypedScopeCreator) { this_.outer_ = o }
func (this_ *TypedScopeCreator_ClassScopeBuilder) GetOuter_() *TypedScopeCreator  { return this_.outer_ }
func (*TypedScopeCreator_ClassScopeBuilder) Implements_jscomp_TypedScopeCreator_AbstractScopeBuilder() {
}
func (*TypedScopeCreator_ClassScopeBuilder) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *TypedScopeCreator_ClassScopeBuilder) initFields_TypedScopeCreator_ClassScopeBuilder() {
	_ = this_
}

// Java constructor: ClassScopeBuilder(com.google.javascript.jscomp.TypedScope scope)
func NewTypedScopeCreator_ClassScopeBuilder(outer_ *TypedScopeCreator, scope *TypedScope) *TypedScopeCreator_ClassScopeBuilder {
	this_ := &TypedScopeCreator_ClassScopeBuilder{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitTypedScopeCreator_ClassScopeBuilder(this_, scope)
	return this_
}

func InitTypedScopeCreator_ClassScopeBuilder(this_ *TypedScopeCreator_ClassScopeBuilder, scope *TypedScope) {
	_ = scope
	this_.TypedScopeCreator_AbstractScopeBuilderBase.outer_ = this_.outer_
	InitTypedScopeCreator_AbstractScopeBuilder(&this_.TypedScopeCreator_AbstractScopeBuilderBase, scope)
	this_.initFields_TypedScopeCreator_ClassScopeBuilder()
}

// Java method: ClassScopeBuilder.visitPreorder(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedScopeCreator_ClassScopeBuilder) VisitPreorder(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsFunction()) {
		return
	}
	if NodeUtil_IsEs6Constructor(n) {
		this_.outer_.SetDeferredType(n, this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope.GetRootNode_0().GetJSType())
	} else {
		this_.DefineFunctionLiteral(n)
	}
}

// Java method: ClassScopeBuilder.visitPostorder(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedScopeCreator_ClassScopeBuilder) VisitPostorder(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if (n.IsName() && (parent == this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope.GetRootNode_0())) && NodeUtil_IsClassExpression(parent) {
		if !(!(j2gort.CheckNotNullString(n.GetString()) == "")) {
			j2gort.CheckState(false)
		}
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(this_).ForDeclarationNode(n).ReadVariableNameFromDeclarationNode().InScope(this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope).WithType(parent.GetJSType()).AllowLaterTypeInference(false).DefineSlot()
	} else if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) {
	} else if n.IsMemberFunctionDef() && parent.IsClassMembers() {
		this_.DefineMemberFunction(n)
	} else if n.IsMemberFieldDef() {
		this_.DefineMemberField(n)
	} else if (n.IsGetterDef() || n.IsSetterDef()) && parent.IsClassMembers() {
		this_.DefineGetterSetter(n)
	} else if n.IsComputedFieldDef() || ((n.IsComputedProp() && n.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_METHOD())) && parent.IsClassMembers()) {
		this_.DefineComputedMemberField(n)
	}
}

// Java method: ClassScopeBuilder.defineMemberFunction(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_ClassScopeBuilder) DefineMemberFunction(n pkg_rhino.Node) {
	_ = n
	var ownerType pkg_rhino.ObjectType = this_.priv_ClassScopeBuilder_determineOwnerTypeForClassMember_5831(n)
	_ = ownerType
	ownerType.GetObjectTypeBase().DefineDeclaredProperty_3_string_JSType_Node(n.GetString(), n.GetLastChild().GetJSType(), n)
}

// Java method: ClassScopeBuilder.defineMemberField(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_ClassScopeBuilder) DefineMemberField(n pkg_rhino.Node) {
	_ = n
	var rhs pkg_rhino.Node = n.GetLastChild()
	_ = rhs
	var declaredType pkg_rhino.JSType = this_.GetDeclaredType(n.GetJSDocInfo(), n, rhs, nil)
	_ = declaredType
	if declaredType == nil || declaredType.IsJavaNil_() {
		var rhsType pkg_rhino.JSType
		_ = rhsType
		if rhs == 0 {
			rhsType = nil
		} else {
			rhsType = this_.TypedScopeCreator_AbstractScopeBuilderBase.GetTypedScopeCreator_AbstractScopeBuilderBase().priv_AbstractScopeBuilder_getDeclaredRValueType_5797(pkg_rhino.Node(0), rhs)
		}
		if (rhsType != nil && !rhsType.IsJavaNil_()) && !(rhsType.GetJSTypeBase().IsUnknownType()) {
			declaredType = rhsType
		} else {
			declaredType = this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_ALL_TYPE())
		}
	}
	var ownerType pkg_rhino.ObjectType = this_.priv_ClassScopeBuilder_determineOwnerTypeForClassMember_5831(n)
	_ = ownerType
	var propName string = n.GetString()
	_ = propName
	ownerType.GetObjectTypeBase().DefineDeclaredProperty_3_string_JSType_Node(propName, declaredType, n)
	n.SetJSType(declaredType)
}

// Java method: ClassScopeBuilder.defineComputedMemberField(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_ClassScopeBuilder) DefineComputedMemberField(n pkg_rhino.Node) {
	_ = n
	var key pkg_rhino.Node = n.GetFirstChild()
	_ = key
	if !(key.IsQualifiedName()) {
		return
	}
	var keyType pkg_rhino.JSType = this_.GetTypedScopeCreator_AbstractScopeBuilderBase().currentScope.LookupQualifiedNameNode(key)
	_ = keyType
	if (keyType == nil || keyType.IsJavaNil_()) || !(keyType.GetJSTypeBase().IsKnownSymbolValueType()) {
		return
	}
	var rhs pkg_rhino.Node = n.GetLastChild()
	_ = rhs
	var declaredType pkg_rhino.JSType = this_.GetDeclaredType(n.GetJSDocInfo(), n, rhs, nil)
	_ = declaredType
	if declaredType == nil || declaredType.IsJavaNil_() {
		var rhsType pkg_rhino.JSType
		_ = rhsType
		if rhs == 0 {
			rhsType = nil
		} else {
			rhsType = this_.TypedScopeCreator_AbstractScopeBuilderBase.GetTypedScopeCreator_AbstractScopeBuilderBase().priv_AbstractScopeBuilder_getDeclaredRValueType_5797(pkg_rhino.Node(0), rhs)
		}
		if (rhsType != nil && !rhsType.IsJavaNil_()) && !(rhsType.GetJSTypeBase().IsUnknownType()) {
			declaredType = rhsType
		} else {
			declaredType = this_.outer_.priv_TypedScopeCreator_getNativeType_5739(pkg_rhino.StaticGet_JSTypeNative_ALL_TYPE())
		}
	}
	var ownerType pkg_rhino.ObjectType = this_.priv_ClassScopeBuilder_determineOwnerTypeForClassMember_5831(n)
	_ = ownerType
	ownerType.GetObjectTypeBase().DefineDeclaredProperty_3_Property_Key_JSType_Node(pkg_rhino.NewProperty_SymbolKey(keyType.ToMaybeKnownSymbolType()), declaredType, n)
	n.SetJSType(declaredType)
}

// Java method: ClassScopeBuilder.defineGetterSetter(com.google.javascript.rhino.Node n)
func (this_ *TypedScopeCreator_ClassScopeBuilder) DefineGetterSetter(n pkg_rhino.Node) {
	_ = n
	var name string = n.GetString()
	_ = name
	var methodType *pkg_rhino.FunctionType = n.GetLastChild().GetJSType().ToMaybeFunctionType()
	_ = methodType
	var propertyType pkg_rhino.JSType = nil
	_ = propertyType
	switch n.TokenSwitchOrd_() {
	case int32(95) /* Token.GETTER_DEF */ :
		propertyType = func() pkg_rhino.JSType {
			if methodType.IsReturnTypeInferred() {
				return pkg_rhino.Upcast_JSType(this_.outer_.unknownType)
			}
			return methodType.GetReturnType()
		}()
	case int32(96) /* Token.SETTER_DEF */ :
		propertyType = func() pkg_rhino.JSType {
			if methodType.GetParameterCount() == int32(0) {
				return pkg_rhino.Upcast_JSType(this_.outer_.unknownType)
			}
			return methodType.GetParameter(int32(0)).GetJSType()
		}()
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.BoxString(n.ToStringTree()))))
	}
	var ownerType pkg_rhino.ObjectType = this_.priv_ClassScopeBuilder_determineOwnerTypeForClassMember_5831(n)
	_ = ownerType
	ownerType.GetObjectTypeBase().DefineDeclaredProperty_3_string_JSType_Node(name, propertyType, n)
}

// Java method: ClassScopeBuilder.determineOwnerTypeForClassMember(com.google.javascript.rhino.Node member)
func (this_ *TypedScopeCreator_ClassScopeBuilder) priv_ClassScopeBuilder_determineOwnerTypeForClassMember_5831(member pkg_rhino.Node) pkg_rhino.ObjectType {
	_ = member
	var ownerNode pkg_rhino.Node = member.GetGrandparent()
	_ = ownerNode
	if !(ownerNode.IsClass()) {
		j2gort.CheckState(false)
	}
	var ownerType *pkg_rhino.FunctionType = ownerNode.GetJSType().ToMaybeFunctionType()
	_ = ownerType
	if member.IsStaticMember() {
		return pkg_rhino.ToIface_ObjectType(ownerType)
	} else if member.IsMemberFieldDef() {
		return ownerType.GetInstanceType()
	} else {
		if !((((member.IsMemberFunctionDef() || member.IsGetterDef()) || member.IsSetterDef()) || member.IsComputedProp()) || member.IsComputedFieldDef()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(member))
		}
		return ownerType.GetPrototype()
	}
}

// Java class: com.google.javascript.jscomp.TypedScopeCreator.FirstOrderFunctionAnalyzer
type TypedScopeCreator_FirstOrderFunctionAnalyzer struct {
	NodeTraversal_AbstractScopedCallbackBase
	outer_ *TypedScopeCreator
}

func (this_ *TypedScopeCreator_FirstOrderFunctionAnalyzer) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedScopeCreator_FirstOrderFunctionAnalyzer) IsJavaNil_() bool { return this_ == nil }
func (*TypedScopeCreator_FirstOrderFunctionAnalyzer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScopeCreator$FirstOrderFunctionAnalyzer", "FirstOrderFunctionAnalyzer"
}
func (this_ *TypedScopeCreator_FirstOrderFunctionAnalyzer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractScopedCallbackBase.IdentityHash_()
}
func (*TypedScopeCreator_FirstOrderFunctionAnalyzer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScopeCreator.FirstOrderFunctionAnalyzer", "com.google.javascript.jscomp.TypedScopeCreator$FirstOrderFunctionAnalyzer", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScopeCreator_FirstOrderFunctionAnalyzer](v_)
	})
}
func (this_ *TypedScopeCreator_FirstOrderFunctionAnalyzer) SetOuter_(o *TypedScopeCreator) {
	this_.outer_ = o
}
func (this_ *TypedScopeCreator_FirstOrderFunctionAnalyzer) GetOuter_() *TypedScopeCreator {
	return this_.outer_
}
func (*TypedScopeCreator_FirstOrderFunctionAnalyzer) Implements_jscomp_NodeTraversal_AbstractScopedCallback() {
}
func (*TypedScopeCreator_FirstOrderFunctionAnalyzer) Implements_jscomp_NodeTraversal_ScopedCallback() {
}
func (*TypedScopeCreator_FirstOrderFunctionAnalyzer) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *TypedScopeCreator_FirstOrderFunctionAnalyzer) initFields_TypedScopeCreator_FirstOrderFunctionAnalyzer() {
	_ = this_
}

// Java constructor: FirstOrderFunctionAnalyzer()
func NewTypedScopeCreator_FirstOrderFunctionAnalyzer(outer_ *TypedScopeCreator) *TypedScopeCreator_FirstOrderFunctionAnalyzer {
	this_ := &TypedScopeCreator_FirstOrderFunctionAnalyzer{}
	this_.outer_ = outer_
	InitTypedScopeCreator_FirstOrderFunctionAnalyzer(this_)
	return this_
}

func InitTypedScopeCreator_FirstOrderFunctionAnalyzer(this_ *TypedScopeCreator_FirstOrderFunctionAnalyzer) {
	InitNodeTraversal_AbstractScopedCallback(&this_.NodeTraversal_AbstractScopedCallbackBase)
	this_.initFields_TypedScopeCreator_FirstOrderFunctionAnalyzer()
}

// Java method: FirstOrderFunctionAnalyzer.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *TypedScopeCreator_FirstOrderFunctionAnalyzer) EnterScope(t *NodeTraversal) {
	_ = t
	var scope *Scope = t.GetScope()
	_ = scope
	if !(scope.IsGlobal()) && (scope.GetVarCount() > int32(0)) {
		var root pkg_rhino.Node = scope.GetRootNode_0()
		_ = root
		var existing []string = this_.outer_.reservedNamesForScope.Get(root)
		_ = existing
		if existing == nil {
			this_.outer_.reservedNamesForScope.PutVoid(root, scope.GetOwnVarNamesArray())
		} else {
			this_.outer_.reservedNamesForScope.PutVoid(root, scope.MergeVarNamesArray(existing))
		}
	}
}

// Java method: FirstOrderFunctionAnalyzer.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypedScopeCreator_FirstOrderFunctionAnalyzer) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if t.InGlobalScope() {
		return
	}
	var containerScope *Scope = j2gort.CastPtr[Scope](t.GetClosestContainerScope())
	_ = containerScope
	if (n.IsReturn() && n.HasChildren()) || (NodeUtil_IsBlocklessArrowFunctionResult(n)) {
		this_.outer_.functionsWithNonEmptyReturns.AddVoid(containerScope.GetRootNode_0())
	}
	if (n.IsName() && NodeUtil_IsLValue(n)) && !(NodeUtil_IsBleedingFunctionName(n)) {
		var name string = n.GetString()
		_ = name
		var scope *Scope = t.GetScope()
		_ = scope
		var var_ *Var = scope.LookupVar(name)
		_ = var_
		if var_ != nil {
			var ownerScope *Scope = j2gort.CastPtr[Scope](var_.GetScope())
			_ = ownerScope
			if ownerScope.IsLocal() {
				this_.outer_.scopedVarStates.RecordAssigned(ownerScope.GetRootNode_0(), name, !(containerScope.HasSameContainerScope(ToIface_AbstractScope(ownerScope))))
			}
		}
	} else if (n.IsGetProp() && n.IsUnscopedQualifiedName()) && NodeUtil_IsLValue(n) {
		var name_2 string = NodeUtil_GetRootOfQualifiedName_1_Node(n).GetString()
		_ = name_2
		var scope_2 *Scope = t.GetScope()
		_ = scope_2
		var var__2 *Var = scope_2.LookupVar(name_2)
		_ = var__2
		if var__2 != nil {
			var ownerScope_2 *Scope = j2gort.CastPtr[Scope](var__2.GetScope())
			_ = ownerScope_2
			if ownerScope_2.IsLocal() && !(containerScope.HasSameContainerScope(ToIface_AbstractScope(ownerScope_2))) {
				this_.outer_.scopedVarStates.RecordEscaped(ownerScope_2.GetRootNode_0(), n.GetQualifiedName())
			}
		}
	}
}
func (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) AlwaysShouldTraverse_()      {}
func (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) CanReleaseSyntacticScopes_() {}
func (*TypedScopeCreator_NormalScopeBuilder) CanReleaseSyntacticScopes_()                    {}
func (*TypedScopeCreator_FunctionScopeBuilder) CanReleaseSyntacticScopes_()                  {}
func (*TypedScopeCreator_ClassScopeBuilder) CanReleaseSyntacticScopes_()                     {}
func (*TypedScopeCreator_FirstOrderFunctionAnalyzer) AlwaysShouldTraverse_()                 {}
