/*
 * Copyright 2025 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/js/RuntimeJsLibManager.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_resources "github.com/paullepoulpe/jsgo/transpiled/gen/resources"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_runtime_js_lib_manager_1 j2gort.Ext_java_util_function_Predicate

func init() {
	lambdaSingleton_runtime_js_lib_manager_1 = j2gort.NewLambda_Ext_java_util_function_Predicate(func(mrp0_ any) bool {
		return j2gort.CastPtr[RuntimeJsLibManager_InternalField](mrp0_).priv_InternalField_isInjected_6480()
	})
}

// Java class: com.google.javascript.jscomp.js.RuntimeJsLibManager
type RuntimeJsLibManager struct {
	changeTracker        *ChangeTracker                                                `gson:"changeTracker"`
	resourceProvider     RuntimeJsLibManager_ResourceProvider                          `gson:"resourceProvider"`
	nodeForCodeInsertion j2gort.Ext_java_util_function_Supplier                        `gson:"nodeForCodeInsertion"`
	injectedLibs         j2gort.Ext_java_util_Set                                      `gson:"injectedLibs"`
	internedFields       j2gort.OrderedMap[string, *RuntimeJsLibManager_InternalField] `gson:"internedFields"`
	lastInjectedLibrary  pkg_rhino.Node                                                `gson:"lastInjectedLibrary"`
	idHash_              atomic.Int32
	mode                 RuntimeJsLibManager_RuntimeLibraryMode `gson:"mode"`
}

func (this_ *RuntimeJsLibManager) BindSelf_(s any) {
	_ = s
}
func (this_ *RuntimeJsLibManager) IsJavaNil_() bool { return this_ == nil }
func (*RuntimeJsLibManager) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.js.RuntimeJsLibManager", "RuntimeJsLibManager"
}
func (this_ *RuntimeJsLibManager) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RuntimeJsLibManager) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.js.RuntimeJsLibManager", "com.google.javascript.jscomp.js.RuntimeJsLibManager", func(v_ any) bool {
		return j2gort.InstanceOf[*RuntimeJsLibManager](v_)
	})
}

func StaticGet_RuntimeJsLibManager_RUNTIME_LIB_DIR() string {
	return "src/com/google/javascript/jscomp/js/"
}
func (this_ *RuntimeJsLibManager) initFields_RuntimeJsLibManager() {
	_ = this_
	this_.injectedLibs = j2gort.NewExt_java_util_LinkedHashSet()
	this_.internedFields = j2gort.NewOrderedMap[string, *RuntimeJsLibManager_InternalField](

	// Java constructor: RuntimeJsLibManager(com.google.javascript.jscomp.ChangeTracker changeTracker, com.google.javascript.jscomp.js.RuntimeJsLibManager.RuntimeLibraryMode mode, com.google.javascript.jscomp.js.RuntimeJsLibManager.ResourceProvider resourceProvider, java.util.function.Supplier<com.google.javascript.rhino.Node> nodeForCodeInsertion)
	)

}

func NewRuntimeJsLibManager(changeTracker *ChangeTracker, mode RuntimeJsLibManager_RuntimeLibraryMode, resourceProvider RuntimeJsLibManager_ResourceProvider, nodeForCodeInsertion j2gort.Ext_java_util_function_Supplier) *RuntimeJsLibManager {
	this_ := &RuntimeJsLibManager{}
	InitRuntimeJsLibManager(this_, changeTracker, mode, resourceProvider, nodeForCodeInsertion)
	return this_
}

func InitRuntimeJsLibManager(this_ *RuntimeJsLibManager, changeTracker *ChangeTracker, mode RuntimeJsLibManager_RuntimeLibraryMode, resourceProvider RuntimeJsLibManager_ResourceProvider, nodeForCodeInsertion j2gort.Ext_java_util_function_Supplier) {
	_ = changeTracker
	_ = mode
	_ = resourceProvider
	_ = nodeForCodeInsertion
	this_.initFields_RuntimeJsLibManager()
	this_.mode = mode
	this_.changeTracker = changeTracker
	this_.resourceProvider = resourceProvider
	this_.nodeForCodeInsertion = nodeForCodeInsertion
}

// Java static method: RuntimeJsLibManager.create(com.google.javascript.jscomp.js.RuntimeJsLibManager.RuntimeLibraryMode mode, com.google.javascript.jscomp.js.RuntimeJsLibManager.ResourceProvider resourceProvider, com.google.javascript.jscomp.ChangeTracker changeTracker, java.util.function.Supplier<com.google.javascript.rhino.Node> nodeForCodeInsertion)
func RuntimeJsLibManager_Create(mode RuntimeJsLibManager_RuntimeLibraryMode, resourceProvider RuntimeJsLibManager_ResourceProvider, changeTracker *ChangeTracker, nodeForCodeInsertion j2gort.Ext_java_util_function_Supplier) *RuntimeJsLibManager {
	_ = mode
	_ = resourceProvider
	_ = changeTracker
	_ = nodeForCodeInsertion
	return NewRuntimeJsLibManager(changeTracker, mode, resourceProvider, nodeForCodeInsertion)
}

// Java method: RuntimeJsLibManager.recordLibraryInjected(java.lang.String resourceName)
func (this_ *RuntimeJsLibManager) RecordLibraryInjected(resourceName string) {
	_ = resourceName
	j2gort.SetAddStr(this_.injectedLibs, resourceName)
}

// Java method: RuntimeJsLibManager.hasInjectedLibrary(java.lang.String resourceName)
func (this_ *RuntimeJsLibManager) HasInjectedLibrary(resourceName string) bool {
	_ = resourceName
	return j2gort.SetContainsStr(this_.injectedLibs, resourceName)
}

// Java method: RuntimeJsLibManager.getInjectedLibraries()
func (this_ *RuntimeJsLibManager) GetInjectedLibraries() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.injectedLibs))
}

// Java method: RuntimeJsLibManager.getExternedFields()
func (this_ *RuntimeJsLibManager) GetExternedFields() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.mode != StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES() {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.internedFields.LiveValues()).Stream().Filter(lambdaSingleton_runtime_js_lib_manager_1).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))
}

// Java method: RuntimeJsLibManager.createField(java.lang.String fieldName)
func (this_ *RuntimeJsLibManager) priv_RuntimeJsLibManager_createField_6470(fieldName string) *RuntimeJsLibManager_InternalField {
	_ = fieldName
	if !(!(j2gort.CheckNotNullString(fieldName) == "")) {
		j2gort.CheckArgumentStr(false, fieldName)
	}
	var parts j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(46)), fieldName)
	_ = parts
	if !(parts.Size() >= int32(2)) {
		j2gort.CheckArgumentFmt(false, "Field name must start with $jscomp., found %s", j2gort.BoxString(fieldName))
	}
	if !(j2gort.CheckNotNullString(j2gort.ListGetStr(parts, int32(0))) == "$jscomp") {
		j2gort.CheckArgumentFmt(false, "Field name must start with $jscomp, found %s", j2gort.BoxString(fieldName))
	}
	var propName string = j2gort.ListGetStr(parts, int32(1))
	_ = propName
	var resourceName string = j2gort.CheckNotNullString(j2gort.CastString(j2gort.MapGetStr(StaticGet_RuntimeJsLibManager_FieldsTable_INSTANCE(), propName)))
	_ = resourceName
	var compiledName string = j2gort.NullString
	_ = compiledName
	switch (this_.mode).Ordinal() {
	case int32(4) /* RuntimeLibraryMode.EXTERN_FIELD_NAMES */ :
		compiledName = j2gort.JStrReplaceChar(fieldName, uint16(46), uint16(95))
	case int32(1) /* RuntimeLibraryMode.NO_OP */, int32(2) /* RuntimeLibraryMode.RECORD_ONLY */, int32(3) /* RuntimeLibraryMode.RECORD_AND_VALIDATE_FIELDS */, int32(0) /* RuntimeLibraryMode.INJECT */ :
		compiledName = fieldName
	}
	return NewRuntimeJsLibManager_InternalField(this_, resourceName, compiledName, fieldName)
}

// Java method: RuntimeJsLibManager.getJsLibField(java.lang.String fieldName)
func (this_ *RuntimeJsLibManager) GetJsLibField(fieldName string) RuntimeJsLibManager_JsLibField {
	_ = fieldName
	return ToIface_RuntimeJsLibManager_JsLibField(this_.priv_RuntimeJsLibManager_getJsLibFieldInternal_6472(fieldName))
}

// Java method: RuntimeJsLibManager.getJsLibFieldInternal(java.lang.String fieldName)
func (this_ *RuntimeJsLibManager) priv_RuntimeJsLibManager_getJsLibFieldInternal_6472(fieldName string) *RuntimeJsLibManager_InternalField {
	_ = fieldName
	return func() *RuntimeJsLibManager_InternalField {
		bound_1 := this_
		return j2gort.OrderedMapComputeIfAbsent(&this_.internedFields, fieldName, nil, func(mrp0_ string) *RuntimeJsLibManager_InternalField {
			return bound_1.priv_RuntimeJsLibManager_createField_6470(mrp0_)
		})
	}()

}

// Java method: RuntimeJsLibManager.injectLibForField(java.lang.String fieldName)
func (this_ *RuntimeJsLibManager) InjectLibForField(fieldName string) {
	_ = fieldName
	var field *RuntimeJsLibManager_InternalField = this_.priv_RuntimeJsLibManager_getJsLibFieldInternal_6472(fieldName)
	_ = field
	if field.priv_InternalField_isInjected_6480() {
		return
	}
	field.priv_InternalField_markInjected_6481()
	if this_.mode == StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES() {
		var externRoot pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](this_.nodeForCodeInsertion.Get())
		_ = externRoot
		if !(externRoot.IsFromExterns()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(externRoot))
		}
		var declaration pkg_rhino.Node = pkg_rhino.IR_Var_1(pkg_rhino.IR_Name(field.QualifiedName()))
		_ = declaration
		declaration.GetFirstChild().PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
		externRoot.AddChildToBack(declaration)
	}
	this_.EnsureLibraryInjected(field.resourceName, false)
}

// Java method: RuntimeJsLibManager.ensureLibraryInjected(java.lang.String resourceName, boolean force)
func (this_ *RuntimeJsLibManager) EnsureLibraryInjected(resourceName string, force bool) pkg_rhino.Node {
	_ = resourceName
	_ = force
	if !(force) {
		switch (this_.mode).Ordinal() {
		case int32(1) /* RuntimeLibraryMode.NO_OP */ :
			return this_.lastInjectedLibrary
		case int32(2) /* RuntimeLibraryMode.RECORD_ONLY */, int32(3) /* RuntimeLibraryMode.RECORD_AND_VALIDATE_FIELDS */, int32(4) /* RuntimeLibraryMode.EXTERN_FIELD_NAMES */ :
			this_.RecordLibraryInjected(resourceName)
			return this_.lastInjectedLibrary
		case int32(0) /* RuntimeLibraryMode.INJECT */ :
		}
	}
	if j2gort.SetContainsStr(this_.injectedLibs, resourceName) {
		return this_.lastInjectedLibrary
	}
	this_.RecordLibraryInjected(resourceName)
	var path string = j2gort.Ext_java_lang_String_Join_2_any_slice_any("", []any{"src/com/google/javascript/jscomp/js/", j2gort.BoxString(resourceName), ".js"})
	_ = path
	var ast pkg_rhino.Node = this_.resourceProvider.Parse(resourceName, path)
	_ = ast
	this_.priv_RuntimeJsLibManager_injectLibraryDependencies_6476(ast, force)
	return this_.priv_RuntimeJsLibManager_injectLibrary_6475(ast)
}

// Java method: RuntimeJsLibManager.injectLibrary(com.google.javascript.rhino.Node ast)
func (this_ *RuntimeJsLibManager) priv_RuntimeJsLibManager_injectLibrary_6475(ast pkg_rhino.Node) pkg_rhino.Node {
	_ = ast
	if !(ast.HasChildren()) {
		return this_.lastInjectedLibrary
	}
	{
		var child pkg_rhino.Node = ast.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.changeTracker.MarkNewScopesChanged(child)
		}
	}
	var endOfLib pkg_rhino.Node = ast.GetLastChild()
	_ = endOfLib
	var firstChild pkg_rhino.Node = ast.RemoveChildren()
	_ = firstChild
	var parent pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](this_.nodeForCodeInsertion.Get())
	_ = parent
	if this_.lastInjectedLibrary == 0 {
		parent.AddChildrenToFront(firstChild)
	} else {
		parent.AddChildrenAfter(firstChild, this_.lastInjectedLibrary)
	}
	this_.lastInjectedLibrary = endOfLib
	this_.changeTracker.ReportChangeToEnclosingScope(parent)
	return endOfLib
}

// Java method: RuntimeJsLibManager.injectLibraryDependencies(com.google.javascript.rhino.Node ast, boolean force)
func (this_ *RuntimeJsLibManager) priv_RuntimeJsLibManager_injectLibraryDependencies_6476(ast pkg_rhino.Node, force bool) {
	_ = ast
	_ = force
	{
		var node pkg_rhino.Node = ast.GetFirstChild()
		_ = node
		for ; ((node != 0) && node.IsExprResult()) && node.GetFirstChild().IsStringLit(); node = ast.GetFirstChild() {
			var directive string = node.GetFirstChild().GetString()
			_ = directive
			var words j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(32)).Limit(int32(2)), directive)
			_ = words
			switch j2gort.CheckNotNullString(j2gort.ListGetStr(words, int32(0))) {
			case "use":
			case "require":
				this_.EnsureLibraryInjected(j2gort.ListGetStr(words, int32(1)), force)
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Bad directive: "+j2gort.StrValString(directive))))
			}
			node.Detach()
		}
	}
}

// Java method: RuntimeJsLibManager.setLastInjectedLibrary(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node lastInjectedLibrary)
func (this_ *RuntimeJsLibManager) SetLastInjectedLibrary(lastInjectedLibrary pkg_rhino.Node) {
	_ = lastInjectedLibrary
	this_.lastInjectedLibrary = lastInjectedLibrary
}

// Java method: RuntimeJsLibManager.getLastInjectedLibrary()
func (this_ *RuntimeJsLibManager) GetLastInjectedLibrary() pkg_rhino.Node {
	return this_.lastInjectedLibrary
}

// Java enum: com.google.javascript.jscomp.js.RuntimeJsLibManager.RuntimeLibraryMode
type RuntimeJsLibManager_RuntimeLibraryMode uint8

const (
	enum_RuntimeJsLibManager_RuntimeLibraryMode_INJECT                     RuntimeJsLibManager_RuntimeLibraryMode = 1
	enum_RuntimeJsLibManager_RuntimeLibraryMode_NO_OP                      RuntimeJsLibManager_RuntimeLibraryMode = 2
	enum_RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY                RuntimeJsLibManager_RuntimeLibraryMode = 3
	enum_RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS RuntimeJsLibManager_RuntimeLibraryMode = 4
	enum_RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES         RuntimeJsLibManager_RuntimeLibraryMode = 5
)

var enumBaseTable_RuntimeJsLibManager_RuntimeLibraryMode = [...]j2gort.EnumBase{
	{Name_: "INJECT", Ordinal_: 0},
	{Name_: "NO_OP", Ordinal_: 1},
	{Name_: "RECORD_ONLY", Ordinal_: 2},
	{Name_: "RECORD_AND_VALIDATE_FIELDS", Ordinal_: 3},
	{Name_: "EXTERN_FIELD_NAMES", Ordinal_: 4},
}
var enumValues_RuntimeJsLibManager_RuntimeLibraryMode = []RuntimeJsLibManager_RuntimeLibraryMode{1, 2, 3, 4, 5}

func StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_INJECT() RuntimeJsLibManager_RuntimeLibraryMode {
	return enum_RuntimeJsLibManager_RuntimeLibraryMode_INJECT
}
func StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_NO_OP() RuntimeJsLibManager_RuntimeLibraryMode {
	return enum_RuntimeJsLibManager_RuntimeLibraryMode_NO_OP
}
func StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY() RuntimeJsLibManager_RuntimeLibraryMode {
	return enum_RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY
}
func StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS() RuntimeJsLibManager_RuntimeLibraryMode {
	return enum_RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS
}
func StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES() RuntimeJsLibManager_RuntimeLibraryMode {
	return enum_RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES
}

func RuntimeJsLibManager_RuntimeLibraryMode_Values() []RuntimeJsLibManager_RuntimeLibraryMode {
	return j2gort.CloneSlice(enumValues_RuntimeJsLibManager_RuntimeLibraryMode)
}

func RuntimeJsLibManager_RuntimeLibraryMode_ValueOf(name string) RuntimeJsLibManager_RuntimeLibraryMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_RuntimeJsLibManager_RuntimeLibraryMode {
		if enumBaseTable_RuntimeJsLibManager_RuntimeLibraryMode[i_].Name_ == name {
			return RuntimeJsLibManager_RuntimeLibraryMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.js.RuntimeJsLibManager.RuntimeLibraryMode." + name))
}

func (this_ RuntimeJsLibManager_RuntimeLibraryMode) IsJavaNil_() bool                              { return this_ == 0 }
func (RuntimeJsLibManager_RuntimeLibraryMode) HasIdentityEquals_()                                 {}
func (RuntimeJsLibManager_RuntimeLibraryMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (RuntimeJsLibManager_RuntimeLibraryMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ RuntimeJsLibManager_RuntimeLibraryMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (RuntimeJsLibManager_RuntimeLibraryMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.js.RuntimeJsLibManager$RuntimeLibraryMode", "RuntimeLibraryMode"
}
func (this_ RuntimeJsLibManager_RuntimeLibraryMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_RuntimeJsLibManager_RuntimeLibraryMode[this_-1]
}
func (this_ RuntimeJsLibManager_RuntimeLibraryMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ RuntimeJsLibManager_RuntimeLibraryMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ RuntimeJsLibManager_RuntimeLibraryMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_RuntimeJsLibManager_RuntimeLibraryMode[this_-1].Name_
}
func (this_ RuntimeJsLibManager_RuntimeLibraryMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_RuntimeJsLibManager_RuntimeLibraryMode[this_-1].Name_
}
func (this_ RuntimeJsLibManager_RuntimeLibraryMode) Equals(other_ any) bool {
	o_, ok_ := other_.(RuntimeJsLibManager_RuntimeLibraryMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ RuntimeJsLibManager_RuntimeLibraryMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_RuntimeJsLibManager_RuntimeLibraryMode[this_-1].HashCode()
}
func (this_ RuntimeJsLibManager_RuntimeLibraryMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_RuntimeJsLibManager_RuntimeLibraryMode[this_-1].CompareTo(other_)
}
func (this_ RuntimeJsLibManager_RuntimeLibraryMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_RuntimeJsLibManager_RuntimeLibraryMode[this_-1].Name_)
}
func init() {
	vals_ := []any{RuntimeJsLibManager_RuntimeLibraryMode(1), RuntimeJsLibManager_RuntimeLibraryMode(2), RuntimeJsLibManager_RuntimeLibraryMode(3), RuntimeJsLibManager_RuntimeLibraryMode(4), RuntimeJsLibManager_RuntimeLibraryMode(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.js.RuntimeJsLibManager.RuntimeLibraryMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.js.RuntimeJsLibManager$RuntimeLibraryMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.js.RuntimeJsLibManager.RuntimeLibraryMode", "com.google.javascript.jscomp.js.RuntimeJsLibManager$RuntimeLibraryMode", func(v_ any) bool {
		return j2gort.InstanceOf[RuntimeJsLibManager_RuntimeLibraryMode](v_)
	})
}

// Java interface: com.google.javascript.jscomp.js.RuntimeJsLibManager.ResourceProvider
type RuntimeJsLibManager_ResourceProvider interface {
	Implements_jscomp_RuntimeJsLibManager_ResourceProvider()
	IsJavaNil_() bool
	Parse(string, string) pkg_rhino.Node
}

func ToIface_RuntimeJsLibManager_ResourceProvider[T interface {
	RuntimeJsLibManager_ResourceProvider
	comparable
}](v T) RuntimeJsLibManager_ResourceProvider {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.js.RuntimeJsLibManager.ResourceProvider", "com.google.javascript.jscomp.js.RuntimeJsLibManager$ResourceProvider", func(v_ any) bool {
		return j2gort.InstanceOf[RuntimeJsLibManager_ResourceProvider](v_)
	})
}

type lambda_RuntimeJsLibManager_ResourceProvider struct {
	fn func(string, string) pkg_rhino.Node
}

func (l *lambda_RuntimeJsLibManager_ResourceProvider) IsJavaNil_() bool { return l == nil }
func NewLambda_RuntimeJsLibManager_ResourceProvider(fn func(string, string) pkg_rhino.Node) RuntimeJsLibManager_ResourceProvider {
	if fn == nil {
		return nil
	}
	return &lambda_RuntimeJsLibManager_ResourceProvider{fn: fn}
}

func (*lambda_RuntimeJsLibManager_ResourceProvider) Implements_jscomp_RuntimeJsLibManager_ResourceProvider() {
}
func (l *lambda_RuntimeJsLibManager_ResourceProvider) Parse(p0_ string, p1_ string) pkg_rhino.Node {
	return l.fn(p0_, p1_)
}

// Java class: com.google.javascript.jscomp.js.RuntimeJsLibManager.FieldsTable
type RuntimeJsLibManager_FieldsTable struct {
	idHash_ atomic.Int32
}

func (this_ *RuntimeJsLibManager_FieldsTable) BindSelf_(s any) {
	_ = s
}
func (this_ *RuntimeJsLibManager_FieldsTable) IsJavaNil_() bool { return this_ == nil }
func (*RuntimeJsLibManager_FieldsTable) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.js.RuntimeJsLibManager$FieldsTable", "FieldsTable"
}
func (this_ *RuntimeJsLibManager_FieldsTable) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RuntimeJsLibManager_FieldsTable) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.js.RuntimeJsLibManager.FieldsTable", "com.google.javascript.jscomp.js.RuntimeJsLibManager$FieldsTable", func(v_ any) bool {
		return j2gort.InstanceOf[*RuntimeJsLibManager_FieldsTable](v_)
	})
}

var static_RuntimeJsLibManager_FieldsTable_INSTANCE j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_RuntimeJsLibManager_FieldsTable_INSTANCE() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_RuntimeJsLibManager_FieldsTable()
	return static_RuntimeJsLibManager_FieldsTable_INSTANCE
}
func StaticSet_RuntimeJsLibManager_FieldsTable_INSTANCE(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_RuntimeJsLibManager_FieldsTable()
	static_RuntimeJsLibManager_FieldsTable_INSTANCE = v
}

var clinitDone_RuntimeJsLibManager_FieldsTable uint32
var clinitState_RuntimeJsLibManager_FieldsTable j2gort.ClinitState

func clinit_RuntimeJsLibManager_FieldsTable() {
	if atomic.LoadUint32(&clinitDone_RuntimeJsLibManager_FieldsTable) != 0 {
		return
	}
	clinitSlow_RuntimeJsLibManager_FieldsTable()
}

func clinitSlow_RuntimeJsLibManager_FieldsTable() {
	if !j2gort.EnterClinit(&clinitState_RuntimeJsLibManager_FieldsTable, &clinitDone_RuntimeJsLibManager_FieldsTable) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RuntimeJsLibManager_FieldsTable, &clinitDone_RuntimeJsLibManager_FieldsTable)
	static_RuntimeJsLibManager_FieldsTable_INSTANCE = RuntimeJsLibManager_FieldsTable_LoadFields()
}

func (this_ *RuntimeJsLibManager_FieldsTable) initFields_RuntimeJsLibManager_FieldsTable() {
	_ = this_
}

// Java constructor: FieldsTable()
func NewRuntimeJsLibManager_FieldsTable() *RuntimeJsLibManager_FieldsTable {
	clinit_RuntimeJsLibManager_FieldsTable()
	this_ := &RuntimeJsLibManager_FieldsTable{}
	InitRuntimeJsLibManager_FieldsTable(this_)
	return this_
}

func InitRuntimeJsLibManager_FieldsTable(this_ *RuntimeJsLibManager_FieldsTable) {
	this_.initFields_RuntimeJsLibManager_FieldsTable()
}

// Java static method: FieldsTable.loadFields()
func RuntimeJsLibManager_FieldsTable_LoadFields() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_RuntimeJsLibManager_FieldsTable()
	var transpilationLibsTxt string = pkg_resources.ResourceLoader_LoadTextResource(j2gort.ClassLiteral("com.google.javascript.jscomp.js.RuntimeJsLibManager", "RuntimeJsLibManager"), "transpilation_libs.txt")
	_ = transpilationLibsTxt
	var allFields j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = allFields
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.SplitterSplitStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(10)).OmitEmptyStrings(), transpilationLibsTxt)))
		for cur_2.HasNext() {
			line := cur_2.NextStr()
			_ = line
			var tokens j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(44)).Limit(int32(2)), line)
			_ = tokens
			if !(tokens.Size() == int32(2)) {
				j2gort.CheckStateMsg(false, tokens)
			}
			var fieldName string = j2gort.ListGetStr(tokens, int32(0))
			_ = fieldName
			var resourcePath string = j2gort.ListGetStr(tokens, int32(1))
			_ = resourcePath
			allFields.Put_2(j2gort.BoxString(fieldName), j2gort.BoxString(resourcePath))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](allFields.BuildOrThrow())
}

// Java class: com.google.javascript.jscomp.js.RuntimeJsLibManager.InternalField
type RuntimeJsLibManager_InternalField struct {
	outer_         *RuntimeJsLibManager
	resourceName   string `gson:"resourceName"`
	qualifiedName  string `gson:"qualifiedName"`
	uncompiledName string `gson:"uncompiledName"`
	idHash_        atomic.Int32
	injected       bool `gson:"injected"`
}

func (this_ *RuntimeJsLibManager_InternalField) BindSelf_(s any) {
	_ = s
}
func (this_ *RuntimeJsLibManager_InternalField) IsJavaNil_() bool { return this_ == nil }
func (*RuntimeJsLibManager_InternalField) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.js.RuntimeJsLibManager$InternalField", "InternalField"
}
func (this_ *RuntimeJsLibManager_InternalField) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RuntimeJsLibManager_InternalField) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.js.RuntimeJsLibManager.InternalField", "com.google.javascript.jscomp.js.RuntimeJsLibManager$InternalField", func(v_ any) bool {
		return j2gort.InstanceOf[*RuntimeJsLibManager_InternalField](v_)
	})
}
func (this_ *RuntimeJsLibManager_InternalField) SetOuter_(o *RuntimeJsLibManager)               { this_.outer_ = o }
func (this_ *RuntimeJsLibManager_InternalField) GetOuter_() *RuntimeJsLibManager                { return this_.outer_ }
func (*RuntimeJsLibManager_InternalField) Implements_jscomp_RuntimeJsLibManager_ExternedField() {}
func (*RuntimeJsLibManager_InternalField) Implements_jscomp_RuntimeJsLibManager_InjectedJsLibField() {
}
func (*RuntimeJsLibManager_InternalField) Implements_jscomp_RuntimeJsLibManager_JsLibField() {}

func (this_ *RuntimeJsLibManager_InternalField) initFields_RuntimeJsLibManager_InternalField() {
	_ = this_
	this_.resourceName = j2gort.NullString
	this_.qualifiedName = j2gort.NullString
	this_.uncompiledName = j2gort.NullString

}

// Java constructor: InternalField(java.lang.String resourceName, java.lang.String qualifiedName, java.lang.String uncompiledName)
func NewRuntimeJsLibManager_InternalField(outer_ *RuntimeJsLibManager, resourceName string, qualifiedName string, uncompiledName string) *RuntimeJsLibManager_InternalField {
	this_ := &RuntimeJsLibManager_InternalField{}
	this_.outer_ = outer_
	InitRuntimeJsLibManager_InternalField(this_, resourceName, qualifiedName, uncompiledName)
	return this_
}

func InitRuntimeJsLibManager_InternalField(this_ *RuntimeJsLibManager_InternalField, resourceName string, qualifiedName string, uncompiledName string) {
	_ = resourceName
	_ = qualifiedName
	_ = uncompiledName
	this_.initFields_RuntimeJsLibManager_InternalField()
	this_.resourceName = resourceName
	this_.qualifiedName = qualifiedName
	this_.uncompiledName = uncompiledName
}

// Java method: InternalField.isInjected()
func (this_ *RuntimeJsLibManager_InternalField) priv_InternalField_isInjected_6480() bool {
	return this_.injected
}

// Java method: InternalField.markInjected()
func (this_ *RuntimeJsLibManager_InternalField) priv_InternalField_markInjected_6481() {
	this_.injected = true
}

// Java method: InternalField.assertInjected()
func (this_ *RuntimeJsLibManager_InternalField) AssertInjected() any {
	switch (this_.outer_.mode).Ordinal() {
	case int32(1) /* RuntimeLibraryMode.NO_OP */, int32(2) /* RuntimeLibraryMode.RECORD_ONLY */ :
	case int32(3) /* RuntimeLibraryMode.RECORD_AND_VALIDATE_FIELDS */, int32(0) /* RuntimeLibraryMode.INJECT */, int32(4) /* RuntimeLibraryMode.EXTERN_FIELD_NAMES */ :
		if !(this_.injected) {
			j2gort.CheckStateFmt(false, "Field %s is not injected", j2gort.BoxString(this_.uncompiledName))
		}
	}
	return this_
}

// Java method: InternalField.matches(com.google.javascript.rhino.Node node)
func (this_ *RuntimeJsLibManager_InternalField) Matches(node pkg_rhino.Node) bool {
	_ = node
	return node.MatchesQualifiedName_1_string(this_.qualifiedName)
}

// Java method: InternalField.toString()
func (this_ *RuntimeJsLibManager_InternalField) String() string {
	return j2gort.Ext_com_google_common_base_Strings_LenientFormat("Field<%s, %s>", []any{j2gort.BoxString(this_.uncompiledName), j2gort.BoxString(this_.resourceName)})
}

// Java method: InternalField.qualifiedName()
func (this_ *RuntimeJsLibManager_InternalField) QualifiedName() string {
	return this_.qualifiedName
}

// Java method: InternalField.uncompiledName()
func (this_ *RuntimeJsLibManager_InternalField) UncompiledName() string {
	return this_.uncompiledName
}

// Java method: InternalField.resourceName()
func (this_ *RuntimeJsLibManager_InternalField) ResourceName() string {
	return this_.resourceName
}

// Java interface: com.google.javascript.jscomp.js.RuntimeJsLibManager.JsLibField
type RuntimeJsLibManager_JsLibField interface {
	Implements_jscomp_RuntimeJsLibManager_JsLibField()
	IsJavaNil_() bool
	Matches(pkg_rhino.Node) bool
	ResourceName() string
	AssertInjected() any
}

func ToIface_RuntimeJsLibManager_JsLibField[T interface {
	RuntimeJsLibManager_JsLibField
	comparable
}](v T) RuntimeJsLibManager_JsLibField {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.js.RuntimeJsLibManager.JsLibField", "com.google.javascript.jscomp.js.RuntimeJsLibManager$JsLibField", func(v_ any) bool {
		return j2gort.InstanceOf[RuntimeJsLibManager_JsLibField](v_)
	})
}

// Java interface: com.google.javascript.jscomp.js.RuntimeJsLibManager.InjectedJsLibField
type RuntimeJsLibManager_InjectedJsLibField interface {
	Implements_jscomp_RuntimeJsLibManager_InjectedJsLibField()
	IsJavaNil_() bool
	RuntimeJsLibManager_JsLibField
	QualifiedName() string
}

func ToIface_RuntimeJsLibManager_InjectedJsLibField[T interface {
	RuntimeJsLibManager_InjectedJsLibField
	comparable
}](v T) RuntimeJsLibManager_InjectedJsLibField {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.js.RuntimeJsLibManager.InjectedJsLibField", "com.google.javascript.jscomp.js.RuntimeJsLibManager$InjectedJsLibField", func(v_ any) bool {
		return j2gort.InstanceOf[RuntimeJsLibManager_InjectedJsLibField](v_)
	})
}

// Java interface: com.google.javascript.jscomp.js.RuntimeJsLibManager.ExternedField
type RuntimeJsLibManager_ExternedField interface {
	Implements_jscomp_RuntimeJsLibManager_ExternedField()
	IsJavaNil_() bool
	RuntimeJsLibManager_InjectedJsLibField
	UncompiledName() string
}

func ToIface_RuntimeJsLibManager_ExternedField[T interface {
	RuntimeJsLibManager_ExternedField
	comparable
}](v T) RuntimeJsLibManager_ExternedField {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.js.RuntimeJsLibManager.ExternedField", "com.google.javascript.jscomp.js.RuntimeJsLibManager$ExternedField", func(v_ any) bool {
		return j2gort.InstanceOf[RuntimeJsLibManager_ExternedField](v_)
	})
}
