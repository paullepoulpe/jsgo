/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/UniqueIdSupplier.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.UniqueIdSupplier
type UniqueIdSupplier struct {
	counter j2gort.Ext_com_google_common_collect_Multiset `gson:"counter"`
	idHash_ atomic.Int32
}

func (this_ *UniqueIdSupplier) BindSelf_(s any) {
	_ = s
}
func (this_ *UniqueIdSupplier) IsJavaNil_() bool { return this_ == nil }
func (*UniqueIdSupplier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.UniqueIdSupplier", "UniqueIdSupplier"
}
func (this_ *UniqueIdSupplier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*UniqueIdSupplier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.UniqueIdSupplier", "com.google.javascript.jscomp.UniqueIdSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[*UniqueIdSupplier](v_)
	})
}

func (this_ *UniqueIdSupplier) initFields_UniqueIdSupplier() {
	_ = this_
}

// Java constructor: UniqueIdSupplier()
func NewUniqueIdSupplier() *UniqueIdSupplier {
	this_ := &UniqueIdSupplier{}
	InitUniqueIdSupplier(this_)
	return this_
}

func InitUniqueIdSupplier(this_ *UniqueIdSupplier) {
	this_.initFields_UniqueIdSupplier()
	this_.counter = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multiset](j2gort.Ext_com_google_common_collect_LinkedHashMultiset_Create_0())
}

// Java method: UniqueIdSupplier.getUniqueId(com.google.javascript.jscomp.CompilerInput input)
func (this_ *UniqueIdSupplier) GetUniqueId(input *CompilerInput) string {
	_ = input
	var filePath string = input.GetSourceFile().GetName()
	_ = filePath
	var fileHashCode int32 = javalang.StringHashCode(j2gort.CheckNotNullString(filePath))
	_ = fileHashCode
	var id int32 = this_.counter.Add_2(j2gort.ValueOfInt(fileHashCode), int32(1))
	_ = id
	var fileHashString string
	_ = fileHashString
	if fileHashCode < int32(0) {
		fileHashString = ("m" + j2gort.StrValInt32(-(fileHashCode)))
	} else {
		fileHashString = ("" + j2gort.StrValInt32(fileHashCode))
	}
	return (j2gort.StrValString(fileHashString) + "$" + j2gort.StrValInt32(id))
}

// Java method: UniqueIdSupplier.toProto()
func (this_ *UniqueIdSupplier) ToProto() j2gort.Ext_java_util_List {
	var result []*j2gort.Ext_com_google_javascript_jscomp_UniqueIdProto = []*j2gort.Ext_com_google_javascript_jscomp_UniqueIdProto{}
	_ = result
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.counter.EntrySet())))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			entry := j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multiset_Entry](rawElem_2)
			_ = entry
			result = append(result, j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_UniqueIdProto](j2gort.Ext_com_google_javascript_jscomp_UniqueIdProto_NewBuilder_0().SetHash(j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](entry.GetElement()))).SetCounter(entry.GetCount()).Build()))
		}
	}
	return j2gort.WrapTypedArrayList(result)
}

// Java static method: UniqueIdSupplier.fromProto(java.util.List<com.google.javascript.jscomp.UniqueIdProto> protos)
func UniqueIdSupplier_FromProto(protos j2gort.Ext_java_util_List) *UniqueIdSupplier {
	_ = protos
	var val_3 UniqueIdSupplier
	InitUniqueIdSupplier(&val_3)
	var supplier *UniqueIdSupplier = &val_3
	_ = supplier
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(protos))
		for cur_4.HasNext() {
			rawElem_5 := cur_4.Next()
			p := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_UniqueIdProto](rawElem_5)
			_ = p
			supplier.counter.Add_2(j2gort.ValueOfInt(p.GetHash()), p.GetCounter())
		}
	}
	return supplier
}
