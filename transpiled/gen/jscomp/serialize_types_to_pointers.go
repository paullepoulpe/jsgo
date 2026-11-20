/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/serialization/SerializeTypesToPointers.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.serialization.SerializeTypesToPointers
type SerializeTypesToPointers struct {
	compiler                  AbstractCompiler                                                  `gson:"compiler"`
	jstypeReconserializer     *JSTypeReconserializer                                            `gson:"jstypeReconserializer"`
	propertiesReferencedInAst *SerializeTypesToPointers_AstPropertySet                          `gson:"propertiesReferencedInAst"`
	typePointersByJstype      *pkg_base.LinkedIdentityHashMap[pkg_rhino.JSType, *j2gort.BoxInt] `gson:"typePointersByJstype"`
	typePool                  *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool   `gson:"typePool"`
	idHash_                   atomic.Int32
}

func (this_ *SerializeTypesToPointers) BindSelf_(s any) {
	_ = s
}
func (this_ *SerializeTypesToPointers) IsJavaNil_() bool { return this_ == nil }
func (*SerializeTypesToPointers) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.SerializeTypesToPointers", "SerializeTypesToPointers"
}
func (this_ *SerializeTypesToPointers) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SerializeTypesToPointers) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.SerializeTypesToPointers", "com.google.javascript.jscomp.serialization.SerializeTypesToPointers", func(v_ any) bool {
		return j2gort.InstanceOf[*SerializeTypesToPointers](v_)
	})
}

var static_SerializeTypesToPointers_GSON *j2gort.Ext_com_google_gson_Gson

func StaticGet_SerializeTypesToPointers_GSON() *j2gort.Ext_com_google_gson_Gson {
	clinit_SerializeTypesToPointers()
	return static_SerializeTypesToPointers_GSON
}
func StaticSet_SerializeTypesToPointers_GSON(v *j2gort.Ext_com_google_gson_Gson) {
	clinit_SerializeTypesToPointers()
	static_SerializeTypesToPointers_GSON = v
}

var clinitDone_SerializeTypesToPointers uint32
var clinitState_SerializeTypesToPointers j2gort.ClinitState

func clinit_SerializeTypesToPointers() {
	if atomic.LoadUint32(&clinitDone_SerializeTypesToPointers) != 0 {
		return
	}
	clinitSlow_SerializeTypesToPointers()
}

func clinitSlow_SerializeTypesToPointers() {
	if !j2gort.EnterClinit(&clinitState_SerializeTypesToPointers, &clinitDone_SerializeTypesToPointers) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SerializeTypesToPointers, &clinitDone_SerializeTypesToPointers)
	static_SerializeTypesToPointers_GSON = j2gort.NewExt_com_google_gson_Gson()
}

func (this_ *SerializeTypesToPointers) initFields_SerializeTypesToPointers() {
	_ = this_
	this_.typePointersByJstype = pkg_base.NewLinkedIdentityHashMap[pkg_rhino.JSType, *j2gort.BoxInt]()

}

// Java constructor: SerializeTypesToPointers(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.serialization.JSTypeReconserializer jstypeReconserializer, com.google.javascript.jscomp.serialization.SerializeTypesToPointers.AstPropertySet propertiesReferencedInAst)
func NewSerializeTypesToPointers(compiler AbstractCompiler, jstypeReconserializer *JSTypeReconserializer, propertiesReferencedInAst *SerializeTypesToPointers_AstPropertySet) *SerializeTypesToPointers {
	clinit_SerializeTypesToPointers()
	this_ := &SerializeTypesToPointers{}
	InitSerializeTypesToPointers(this_, compiler, jstypeReconserializer, propertiesReferencedInAst)
	return this_
}

func InitSerializeTypesToPointers(this_ *SerializeTypesToPointers, compiler AbstractCompiler, jstypeReconserializer *JSTypeReconserializer, propertiesReferencedInAst *SerializeTypesToPointers_AstPropertySet) {
	_ = compiler
	_ = jstypeReconserializer
	_ = propertiesReferencedInAst
	this_.initFields_SerializeTypesToPointers()
	this_.compiler = compiler
	this_.jstypeReconserializer = jstypeReconserializer
	this_.propertiesReferencedInAst = propertiesReferencedInAst
}

// Java static method: SerializeTypesToPointers.create(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.serialization.StringPool.Builder stringPoolBuilder, com.google.javascript.jscomp.serialization.SerializationOptions serializationOptions)
func SerializeTypesToPointers_Create(compiler AbstractCompiler, stringPoolBuilder *StringPool_Builder, serializationOptions *SerializationOptions) *SerializeTypesToPointers {
	clinit_SerializeTypesToPointers()
	_ = compiler
	_ = stringPoolBuilder
	_ = serializationOptions
	var invalidatingTypes *InvalidatingTypes = NewInvalidatingTypes_Builder(compiler.GetTypeRegistry()).AddAllTypeMismatches(compiler.GetTypeMismatches()).Build()
	_ = invalidatingTypes
	var val_1 SerializeTypesToPointers_AstPropertySet
	InitSerializeTypesToPointers_AstPropertySet(&val_1)
	var propertiesReferencedInAst *SerializeTypesToPointers_AstPropertySet = &val_1
	_ = propertiesReferencedInAst
	var jsTypeReconserializer *JSTypeReconserializer = JSTypeReconserializer_CreateWithPropertySet(compiler.GetTypeRegistry(), invalidatingTypes, stringPoolBuilder, propertiesReferencedInAst, serializationOptions)
	_ = jsTypeReconserializer
	return NewSerializeTypesToPointers(compiler, jsTypeReconserializer, propertiesReferencedInAst)
}

// Java method: SerializeTypesToPointers.gatherTypesOnAst(com.google.javascript.rhino.Node root)
func (this_ *SerializeTypesToPointers) GatherTypesOnAst(root pkg_rhino.Node) {
	_ = root
	if !(this_.typePool == nil) {
		j2gort.CheckStateStr(false, "Cannot call process() twice")
	}
	NodeTraversal_Traverse(this_.compiler, root, NewSerializeTypesToPointers_PropertySearchCallback(this_))
	NodeTraversal_Traverse(this_.compiler, root, NewSerializeTypesToPointers_TypeSearchCallback(this_))
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.compiler.GetTypeMismatches()))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			mismatch := j2gort.CastPtr[TypeMismatch](rawElem_3)
			_ = mismatch
			this_.jstypeReconserializer.SerializeType(mismatch.Found())
			this_.jstypeReconserializer.SerializeType(mismatch.Required())
		}
	}
	this_.typePool = this_.jstypeReconserializer.GenerateTypePool()
	this_.priv_SerializeTypesToPointers_logSerializationDebugInfo_7799(this_.jstypeReconserializer, this_.typePool)
}

// Java method: SerializeTypesToPointers.getTypePointersByJstype()
func (this_ *SerializeTypesToPointers) GetTypePointersByJstype() *pkg_base.LinkedIdentityHashMap[pkg_rhino.JSType, *j2gort.BoxInt] {
	return this_.typePointersByJstype
}

// Java method: SerializeTypesToPointers.getTypePool()
func (this_ *SerializeTypesToPointers) GetTypePool() *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool {
	return this_.typePool
}

// Java method: SerializeTypesToPointers.logSerializationDebugInfo(com.google.javascript.jscomp.serialization.JSTypeReconserializer serializer, com.google.javascript.jscomp.serialization.TypePool typePool)
func (this_ *SerializeTypesToPointers) priv_SerializeTypesToPointers_logSerializationDebugInfo_7799(serializer *JSTypeReconserializer, typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool) {
	_ = serializer
	_ = typePool
	func() {
		var log pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenLog(j2gort.GetObjectClass(this_), "object_uuids.log", []string{})
		_ = log
		defer j2gort.CloseResource(log)
		if log.IsLogging() {
			var allSerializedTypes j2gort.Ext_com_google_common_collect_ImmutableMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](serializer.GetColorIdToJSTypeMapForDebugging().AsMap())
			_ = allSerializedTypes
			log.LogJson_1_LogFile_StreamedJsonProducer(NewSerializeTypesToPointers_StreamObjectUuidsJson(allSerializedTypes))
		}
	}()
	func() {
		var log_2 pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenLog(j2gort.GetObjectClass(this_), "mismatches.log", []string{})
		_ = log_2
		defer j2gort.CloseResource(log_2)
		if log_2.IsLogging() {
			log_2.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
				return j2gort.BoxString(StaticGet_SerializeTypesToPointers_GSON().ToJson_1_any(this_.priv_SerializeTypesToPointers_logTypeMismatches_7800(this_.compiler.GetTypeMismatches(), serializer, typePool)))
			}))
		}
	}()
}

// Java method: SerializeTypesToPointers.logTypeMismatches(java.lang.Iterable<com.google.javascript.jscomp.TypeMismatch> typeMismatches, com.google.javascript.jscomp.serialization.JSTypeReconserializer serializer, com.google.javascript.jscomp.serialization.TypePool typePool)
func (this_ *SerializeTypesToPointers) priv_SerializeTypesToPointers_logTypeMismatches_7800(typeMismatches j2gort.Ext_java_lang_Iterable, serializer *JSTypeReconserializer, typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool) j2gort.Ext_com_google_common_collect_ImmutableSortedSet {
	_ = typeMismatches
	_ = serializer
	_ = typePool
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSortedSet](j2gort.Ext_com_google_common_collect_Streams_Stream_1_Ext_java_lang_Iterable(typeMismatches).Map(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_4 any) any {
		mismatch := j2gort.CastPtr[TypeMismatch](lp_4)
		_ = mismatch
		return j2gort.PtrToAny(SerializeTypesToPointers_TypeMismatchJson_Create(mismatch, serializer, typePool))
	})).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSortedSet_ToImmutableSortedSet(j2gort.Ext_java_util_Comparator_NaturalOrder())))
}

// Java class: com.google.javascript.jscomp.serialization.SerializeTypesToPointers.AstPropertySet
type SerializeTypesToPointers_AstPropertySet struct {
	keys    []string `gson:"keys"`
	size    int32    `gson:"size"`
	idHash_ atomic.Int32
}

func (this_ *SerializeTypesToPointers_AstPropertySet) BindSelf_(s any) {
	_ = s
}
func (this_ *SerializeTypesToPointers_AstPropertySet) IsJavaNil_() bool { return this_ == nil }
func (*SerializeTypesToPointers_AstPropertySet) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.SerializeTypesToPointers$AstPropertySet", "AstPropertySet"
}
func (this_ *SerializeTypesToPointers_AstPropertySet) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SerializeTypesToPointers_AstPropertySet) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.SerializeTypesToPointers.AstPropertySet", "com.google.javascript.jscomp.serialization.SerializeTypesToPointers$AstPropertySet", func(v_ any) bool {
		return j2gort.InstanceOf[*SerializeTypesToPointers_AstPropertySet](v_)
	})
}

func (this_ *SerializeTypesToPointers_AstPropertySet) initFields_SerializeTypesToPointers_AstPropertySet() {
	_ = this_
	this_.keys = j2gort.NewStringArray(int32(4096))

}

// Java constructor: AstPropertySet()
func NewSerializeTypesToPointers_AstPropertySet() *SerializeTypesToPointers_AstPropertySet {
	this_ := &SerializeTypesToPointers_AstPropertySet{}
	InitSerializeTypesToPointers_AstPropertySet(this_)
	return this_
}

func InitSerializeTypesToPointers_AstPropertySet(this_ *SerializeTypesToPointers_AstPropertySet) {
	this_.initFields_SerializeTypesToPointers_AstPropertySet()
}

// Java method: AstPropertySet.add(java.lang.String name)
func (this_ *SerializeTypesToPointers_AstPropertySet) Add(name string) {
	_ = name
	if name == j2gort.NullString {
		return
	}
	var table []string = this_.keys
	_ = table
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var h int32 = javalang.StringHashCode(j2gort.CheckNotNullString(name))
	_ = h
	var slot int32 = ((h ^ (javalang.Ushr32(h, int(int32(16))))) & mask)
	_ = slot
	for {
		var existing string = table[slot]
		_ = existing
		if existing == j2gort.NullString {
			{
				a_5 := table
				i_6 := slot
				a_5[i_6] = name
			}
			this_.size = this_.size + 1
			if (this_.size * int32(4)) >= (int32(len(table)) * int32(3)) {
				this_.priv_AstPropertySet_rehash_7803((int32(len(table)) * int32(2)))
			}
			return
		}
		if j2gort.CheckNotNullString(existing) == name {
			return
		}
		slot = ((slot + int32(1)) & mask)
	}
}

// Java method: AstPropertySet.contains(java.lang.String name)
func (this_ *SerializeTypesToPointers_AstPropertySet) Contains(name string) bool {
	_ = name
	if (this_.size == int32(0)) || (name == j2gort.NullString) {
		return false
	}
	var table []string = this_.keys
	_ = table
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var h int32 = javalang.StringHashCode(j2gort.CheckNotNullString(name))
	_ = h
	var slot int32 = ((h ^ (javalang.Ushr32(h, int(int32(16))))) & mask)
	_ = slot
	for {
		var existing string = table[slot]
		_ = existing
		if existing == j2gort.NullString {
			return false
		}
		if j2gort.CheckNotNullString(existing) == name {
			return true
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: AstPropertySet.rehash(int newCap)
func (this_ *SerializeTypesToPointers_AstPropertySet) priv_AstPropertySet_rehash_7803(newCap int32) {
	_ = newCap
	var oldTable []string = this_.keys
	_ = oldTable
	var newTable []string = j2gort.NewStringArray(newCap)
	_ = newTable
	var mask int32 = (newCap - int32(1))
	_ = mask
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(oldTable))
		_ = len_
		for ; i < len_; i = i + 1 {
			var k string = oldTable[i]
			_ = k
			if k != j2gort.NullString {
				var h int32 = javalang.StringHashCode(j2gort.CheckNotNullString(k))
				_ = h
				var slot int32 = ((h ^ (javalang.Ushr32(h, int(int32(16))))) & mask)
				_ = slot
				for newTable[slot] != j2gort.NullString {
					slot = ((slot + int32(1)) & mask)
				}
				{
					a_7 := newTable
					i_8 := slot
					a_7[i_8] = k
				}
			}
		}
	}
	this_.keys = newTable
}

// Java class: com.google.javascript.jscomp.serialization.SerializeTypesToPointers.PropertySearchCallback
type SerializeTypesToPointers_PropertySearchCallback struct {
	outer_  *SerializeTypesToPointers
	idHash_ atomic.Int32
}

func (this_ *SerializeTypesToPointers_PropertySearchCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *SerializeTypesToPointers_PropertySearchCallback) IsJavaNil_() bool { return this_ == nil }
func (*SerializeTypesToPointers_PropertySearchCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.SerializeTypesToPointers$PropertySearchCallback", "PropertySearchCallback"
}
func (this_ *SerializeTypesToPointers_PropertySearchCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SerializeTypesToPointers_PropertySearchCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.SerializeTypesToPointers.PropertySearchCallback", "com.google.javascript.jscomp.serialization.SerializeTypesToPointers$PropertySearchCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*SerializeTypesToPointers_PropertySearchCallback](v_)
	})
}
func (this_ *SerializeTypesToPointers_PropertySearchCallback) SetOuter_(o *SerializeTypesToPointers) {
	this_.outer_ = o
}
func (this_ *SerializeTypesToPointers_PropertySearchCallback) GetOuter_() *SerializeTypesToPointers {
	return this_.outer_
}
func (*SerializeTypesToPointers_PropertySearchCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *SerializeTypesToPointers_PropertySearchCallback) initFields_SerializeTypesToPointers_PropertySearchCallback() {
	_ = this_
}

// Java constructor: PropertySearchCallback()
func NewSerializeTypesToPointers_PropertySearchCallback(outer_ *SerializeTypesToPointers) *SerializeTypesToPointers_PropertySearchCallback {
	this_ := &SerializeTypesToPointers_PropertySearchCallback{}
	this_.outer_ = outer_
	InitSerializeTypesToPointers_PropertySearchCallback(this_)
	return this_
}

func InitSerializeTypesToPointers_PropertySearchCallback(this_ *SerializeTypesToPointers_PropertySearchCallback) {
	this_.initFields_SerializeTypesToPointers_PropertySearchCallback()
}

// Java method: PropertySearchCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SerializeTypesToPointers_PropertySearchCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(n.IsScript()) || !(NodeUtil_IsFromTypeSummary(n)))
}

// Java method: PropertySearchCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SerializeTypesToPointers_PropertySearchCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		this_.outer_.propertiesReferencedInAst.Add(n.GetString())
	case int32(100) /* Token.STRING_KEY */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
		if !(n.IsQuotedStringKey()) {
			this_.outer_.propertiesReferencedInAst.Add(n.GetString())
		}
	default:
	}
}

// Java class: com.google.javascript.jscomp.serialization.SerializeTypesToPointers.TypeSearchCallback
type SerializeTypesToPointers_TypeSearchCallback struct {
	outer_  *SerializeTypesToPointers
	idHash_ atomic.Int32
}

func (this_ *SerializeTypesToPointers_TypeSearchCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *SerializeTypesToPointers_TypeSearchCallback) IsJavaNil_() bool { return this_ == nil }
func (*SerializeTypesToPointers_TypeSearchCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.SerializeTypesToPointers$TypeSearchCallback", "TypeSearchCallback"
}
func (this_ *SerializeTypesToPointers_TypeSearchCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SerializeTypesToPointers_TypeSearchCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.SerializeTypesToPointers.TypeSearchCallback", "com.google.javascript.jscomp.serialization.SerializeTypesToPointers$TypeSearchCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*SerializeTypesToPointers_TypeSearchCallback](v_)
	})
}
func (this_ *SerializeTypesToPointers_TypeSearchCallback) SetOuter_(o *SerializeTypesToPointers) {
	this_.outer_ = o
}
func (this_ *SerializeTypesToPointers_TypeSearchCallback) GetOuter_() *SerializeTypesToPointers {
	return this_.outer_
}
func (*SerializeTypesToPointers_TypeSearchCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *SerializeTypesToPointers_TypeSearchCallback) initFields_SerializeTypesToPointers_TypeSearchCallback() {
	_ = this_
}

// Java constructor: TypeSearchCallback()
func NewSerializeTypesToPointers_TypeSearchCallback(outer_ *SerializeTypesToPointers) *SerializeTypesToPointers_TypeSearchCallback {
	this_ := &SerializeTypesToPointers_TypeSearchCallback{}
	this_.outer_ = outer_
	InitSerializeTypesToPointers_TypeSearchCallback(this_)
	return this_
}

func InitSerializeTypesToPointers_TypeSearchCallback(this_ *SerializeTypesToPointers_TypeSearchCallback) {
	this_.initFields_SerializeTypesToPointers_TypeSearchCallback()
}

// Java method: TypeSearchCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SerializeTypesToPointers_TypeSearchCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(n.IsScript()) || !(NodeUtil_IsFromTypeSummary(n)))
}

// Java method: TypeSearchCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SerializeTypesToPointers_TypeSearchCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsRoot() {
		return
	}
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	if type_ != nil && !type_.IsJavaNil_() {
		var existing *j2gort.BoxInt = this_.outer_.typePointersByJstype.Get(type_)
		_ = existing
		if existing == nil {
			this_.outer_.typePointersByJstype.Put(type_, j2gort.ValueOfInt(this_.outer_.jstypeReconserializer.SerializeType(type_)))
		}
	}
	var shadow pkg_rhino.Node = n.GetClosureUnawareShadow()
	_ = shadow
	if shadow != 0 {
		NodeTraversal_Traverse(this_.outer_.compiler, shadow.GetFirstFirstChild().GetFirstChild(), this_)
	}
}

// Java class: com.google.javascript.jscomp.serialization.SerializeTypesToPointers.StreamObjectUuidsJson
type SerializeTypesToPointers_StreamObjectUuidsJson struct {
	allSerializedTypes j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"allSerializedTypes"`
	idHash_            atomic.Int32
}

func (this_ *SerializeTypesToPointers_StreamObjectUuidsJson) BindSelf_(s any) {
	_ = s
}
func (this_ *SerializeTypesToPointers_StreamObjectUuidsJson) IsJavaNil_() bool { return this_ == nil }
func (*SerializeTypesToPointers_StreamObjectUuidsJson) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.SerializeTypesToPointers$StreamObjectUuidsJson", "StreamObjectUuidsJson"
}
func (this_ *SerializeTypesToPointers_StreamObjectUuidsJson) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SerializeTypesToPointers_StreamObjectUuidsJson) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.SerializeTypesToPointers.StreamObjectUuidsJson", "com.google.javascript.jscomp.serialization.SerializeTypesToPointers$StreamObjectUuidsJson", func(v_ any) bool {
		return j2gort.InstanceOf[*SerializeTypesToPointers_StreamObjectUuidsJson](v_)
	})
}
func (*SerializeTypesToPointers_StreamObjectUuidsJson) Implements_diagnostic_LogFile_StreamedJsonProducer() {
}

func (this_ *SerializeTypesToPointers_StreamObjectUuidsJson) initFields_SerializeTypesToPointers_StreamObjectUuidsJson() {
	_ = this_
}

// Java constructor: StreamObjectUuidsJson(com.google.common.collect.ImmutableMap<java.lang.String,java.util.Collection<com.google.javascript.rhino.jstype.JSType>> allSerializedTypes)
func NewSerializeTypesToPointers_StreamObjectUuidsJson(allSerializedTypes j2gort.Ext_com_google_common_collect_ImmutableMap) *SerializeTypesToPointers_StreamObjectUuidsJson {
	this_ := &SerializeTypesToPointers_StreamObjectUuidsJson{}
	InitSerializeTypesToPointers_StreamObjectUuidsJson(this_, allSerializedTypes)
	return this_
}

func InitSerializeTypesToPointers_StreamObjectUuidsJson(this_ *SerializeTypesToPointers_StreamObjectUuidsJson, allSerializedTypes j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = allSerializedTypes
	this_.initFields_SerializeTypesToPointers_StreamObjectUuidsJson()
	this_.allSerializedTypes = allSerializedTypes
}

// Java method: StreamObjectUuidsJson.writeJson(com.google.gson.stream.JsonWriter jsonWriter)
func (this_ *SerializeTypesToPointers_StreamObjectUuidsJson) WriteJson(jsonWriter *j2gort.Ext_com_google_gson_stream_JsonWriter) {
	_ = jsonWriter
	jsonWriter.BeginObject()
	{
		cur_9 := j2gort.NewMapEntryCursor(this_.allSerializedTypes)
		for cur_9.NextEntry() {
			jsonWriter.Name(cur_9.KeyStr())
			jsonWriter.BeginArray()
			var jstypes []string = []string{}
			_ = jstypes
			{
				cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Collection](cur_9.GetValue())))
				for cur_10.HasNext() {
					rawElem_11 := cur_10.Next()
					jstype := j2gort.AsIface[pkg_rhino.JSType](rawElem_11)
					_ = jstype
					jstypes = append(jstypes, jstype.GetJSTypeBase().String())
				}
			}
			j2gort.SliceSortNatural(jstypes)
			for _, typeName := range jstypes {
				_ = typeName
				jsonWriter.Value_1_string(typeName)
			}

			jsonWriter.EndArray()
		}
	}
	jsonWriter.EndObject()
}

// Java class: com.google.javascript.jscomp.serialization.SerializeTypesToPointers.TypeMismatchJson
type SerializeTypesToPointers_TypeMismatchJson struct {
	location        string `gson:"location"`
	foundColorId    string `gson:"foundColorId"`
	requiredColorId string `gson:"requiredColorId"`
	idHash_         atomic.Int32
}

func (this_ *SerializeTypesToPointers_TypeMismatchJson) BindSelf_(s any) {
	_ = s
}
func (this_ *SerializeTypesToPointers_TypeMismatchJson) IsJavaNil_() bool { return this_ == nil }
func (*SerializeTypesToPointers_TypeMismatchJson) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.SerializeTypesToPointers$TypeMismatchJson", "TypeMismatchJson"
}
func (this_ *SerializeTypesToPointers_TypeMismatchJson) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SerializeTypesToPointers_TypeMismatchJson) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.SerializeTypesToPointers.TypeMismatchJson", "com.google.javascript.jscomp.serialization.SerializeTypesToPointers$TypeMismatchJson", func(v_ any) bool {
		return j2gort.InstanceOf[*SerializeTypesToPointers_TypeMismatchJson](v_)
	})
}

func (this_ *SerializeTypesToPointers_TypeMismatchJson) initFields_SerializeTypesToPointers_TypeMismatchJson() {
	_ = this_
	this_.location = j2gort.NullString
	this_.foundColorId = j2gort.NullString
	this_.requiredColorId = j2gort.NullString
}

// Java constructor: TypeMismatchJson(com.google.javascript.jscomp.TypeMismatch x, com.google.javascript.jscomp.colors.ColorId found, com.google.javascript.jscomp.colors.ColorId required)
func NewSerializeTypesToPointers_TypeMismatchJson(x *TypeMismatch, found pkg_colors.ColorId, required pkg_colors.ColorId) *SerializeTypesToPointers_TypeMismatchJson {
	this_ := &SerializeTypesToPointers_TypeMismatchJson{}
	InitSerializeTypesToPointers_TypeMismatchJson(this_, x, found, required)
	return this_
}

func InitSerializeTypesToPointers_TypeMismatchJson(this_ *SerializeTypesToPointers_TypeMismatchJson, x *TypeMismatch, found pkg_colors.ColorId, required pkg_colors.ColorId) {
	_ = x
	_ = found
	_ = required
	this_.initFields_SerializeTypesToPointers_TypeMismatchJson()
	this_.location = x.Location().GetLocation()
	this_.foundColorId = found.String()
	this_.requiredColorId = required.String()
}

// Java static method: TypeMismatchJson.create(com.google.javascript.jscomp.TypeMismatch x, com.google.javascript.jscomp.serialization.JSTypeReconserializer serializer, com.google.javascript.jscomp.serialization.TypePool typePool)
func SerializeTypesToPointers_TypeMismatchJson_Create(x *TypeMismatch, serializer *JSTypeReconserializer, typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool) *SerializeTypesToPointers_TypeMismatchJson {
	_ = x
	_ = serializer
	_ = typePool
	var foundPointer int32 = serializer.SerializeType(x.Found())
	_ = foundPointer
	var requiredPointer int32 = serializer.SerializeType(x.Required())
	_ = requiredPointer
	return NewSerializeTypesToPointers_TypeMismatchJson(x, SerializeTypesToPointers_TypeMismatchJson_TypePointerToId(foundPointer, typePool), SerializeTypesToPointers_TypeMismatchJson_TypePointerToId(requiredPointer, typePool))
}

// Java static method: TypeMismatchJson.typePointerToId(int poolOffset, com.google.javascript.jscomp.serialization.TypePool typePool)
func SerializeTypesToPointers_TypeMismatchJson_TypePointerToId(poolOffset int32, typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool) pkg_colors.ColorId {
	_ = poolOffset
	_ = typePool
	if TypePointers_IsAxiomatic(poolOffset) {
		return j2gort.CastPtr[pkg_colors.Color](StaticGet_TypePointers_OFFSET_TO_AXIOMATIC_COLOR().Get(poolOffset)).GetId()
	}
	var typeProto *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](typePool.GetTypeList().Get(TypePointers_TrimOffset(poolOffset)))
	_ = typeProto
	switch j2gort.CheckNotNullPtr(typeProto.GetKindCase()).Ordinal() {
	case int32(1) /* KindCase.UNION */ :
		return pkg_colors.ColorId_Union(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](typeProto.GetUnion().GetUnionMemberList().Stream().Map(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_13 any) any {
			pointer := j2gort.CastPtr[j2gort.BoxInt](lp_13)
			_ = pointer
			return j2gort.BoxNullableVal(SerializeTypesToPointers_TypeMismatchJson_TypePointerToId(j2gort.UnboxInt(pointer), typePool))
		})).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet())))
	case int32(0) /* KindCase.OBJECT */ :
		return pkg_colors.ColorId_FromBytes_1_Ext_com_google_protobuf_ByteString(typeProto.GetObject().GetUuid())
	case int32(2) /* KindCase.KIND_NOT_SET */ :
	}
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Unrecognized TypeProto "+j2gort.StrVal(typeProto))))
}

// Java method: TypeMismatchJson.compareTo(com.google.javascript.jscomp.serialization.SerializeTypesToPointers.TypeMismatchJson x)
func (this_ *SerializeTypesToPointers_TypeMismatchJson) CompareTo(arg_14 any) int32 {
	x := j2gort.CastPtr[SerializeTypesToPointers_TypeMismatchJson](arg_14)
	_ = x
	return j2gort.Ext_com_google_common_collect_ComparisonChain_Start().Compare_2_any_any(j2gort.BoxString(this_.foundColorId), j2gort.BoxString(x.foundColorId)).Compare_2_any_any(j2gort.BoxString(this_.requiredColorId), j2gort.BoxString(x.requiredColorId)).Compare_2_any_any(j2gort.BoxString(this_.location), j2gort.BoxString(x.location)).Result()
}
func (*SerializeTypesToPointers_PropertySearchCallback) CanReleaseSyntacticScopes_() {}
func (*SerializeTypesToPointers_TypeSearchCallback) CanReleaseSyntacticScopes_()     {}
