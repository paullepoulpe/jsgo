/*
 * Copyright 2026 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/SimpleTable.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.SimpleTable
type SimpleTable[R comparable, C comparable, V any] struct {
	rowMap j2gort.Ext_java_util_Map `gson:"rowMap"`
}

func (this_ *SimpleTable[R, C, V]) BindSelf_(s any) {
	_ = s
}
func (this_ *SimpleTable[R, C, V]) IsJavaNil_() bool { return this_ == nil }
func (*SimpleTable[R, C, V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SimpleTable", "SimpleTable"
}
func (*SimpleTable[R, C, V]) IsGenericInst_jscomp_SimpleTable_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SimpleTable", "com.google.javascript.jscomp.SimpleTable", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_jscomp_SimpleTable_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

func StaticGet_SimpleTable_SerialVersionUID() int64 { return int64(1) }
func (this_ *SimpleTable[R, C, V]) initFields_SimpleTable() {
	_ = this_
}

// Java constructor: SimpleTable()
func NewSimpleTable[R comparable, C comparable, V any]() *SimpleTable[R, C, V] {
	this_ := &SimpleTable[R, C, V]{}
	InitSimpleTable[R, C, V](this_)
	return this_
}

func InitSimpleTable[R comparable, C comparable, V any](this_ *SimpleTable[R, C, V]) {
	this_.initFields_SimpleTable()
	this_.rowMap = j2gort.NewExt_java_util_LinkedHashMap()
}

// Java static method: SimpleTable.create()
func SimpleTable_Create[R comparable, C comparable, V any]() *SimpleTable[R, C, V] {
	return NewSimpleTable[R, C, V]()
}

// Java method: SimpleTable.contains(java.lang.Object rowKey, java.lang.Object columnKey)
func (this_ *SimpleTable[R, C, V]) Contains(rowKey any, columnKey any) bool {
	_ = rowKey
	_ = columnKey
	var row j2gort.Ext_java_util_Map = j2gort.AsIface[j2gort.Ext_java_util_Map](this_.rowMap.Get(rowKey))
	_ = row
	return ((row != nil) && row.ContainsKey(columnKey))
}

// Java method: SimpleTable.containsRow(java.lang.Object rowKey)
func (this_ *SimpleTable[R, C, V]) ContainsRow(rowKey any) bool {
	_ = rowKey
	return this_.rowMap.ContainsKey(rowKey)
}

// Java method: SimpleTable.containsColumn(java.lang.Object columnKey)
func (this_ *SimpleTable[R, C, V]) ContainsColumn(columnKey any) bool {
	_ = columnKey
	{
		cur_1 := j2gort.NewMapValCursor(this_.rowMap)
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			row := j2gort.AsIface[j2gort.Ext_java_util_Map](rawElem_2)
			_ = row
			if row.ContainsKey(columnKey) {
				return true
			}
		}
	}
	return false
}

// Java method: SimpleTable.containsValue(java.lang.Object value)
func (this_ *SimpleTable[R, C, V]) ContainsValue(value any) bool {
	_ = value
	{
		cur_3 := j2gort.NewMapValCursor(this_.rowMap)
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			row := j2gort.AsIface[j2gort.Ext_java_util_Map](rawElem_4)
			_ = row
			if row.ContainsValue(value) {
				return true
			}
		}
	}
	return false
}

// Java method: SimpleTable.get(java.lang.Object rowKey, java.lang.Object columnKey)
func (this_ *SimpleTable[R, C, V]) Get(rowKey any, columnKey any) V {
	_ = rowKey
	_ = columnKey
	var row j2gort.Ext_java_util_Map = j2gort.AsIface[j2gort.Ext_java_util_Map](this_.rowMap.Get(rowKey))
	_ = row
	if row == nil {
		{
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				return any(j2gort.NullString).(V)
			}
			var z_ V
			return z_
		}
	} else {
		{
			tv_5 := row.Get(columnKey)
			if tv_5 == nil {
				{
					if _, ok_ := any((*V)(nil)).(*string); ok_ {
						return any(j2gort.NullString).(V)
					}
					var z_ V
					return z_
				}
			}
			return j2gort.Cast[V](tv_5)
		}
	}
}

// Java method: SimpleTable.isEmpty()
func (this_ *SimpleTable[R, C, V]) IsEmpty() bool {
	return this_.rowMap.IsEmpty()
}

// Java method: SimpleTable.size()
func (this_ *SimpleTable[R, C, V]) Size() int32 {
	var total int32 = int32(0)
	_ = total
	{
		cur_6 := j2gort.NewMapValCursor(this_.rowMap)
		for cur_6.HasNext() {
			rawElem_7 := cur_6.Next()
			row := j2gort.AsIface[j2gort.Ext_java_util_Map](rawElem_7)
			_ = row
			total = (total + row.Size())
		}
	}
	return total
}

// Java method: SimpleTable.clear()
func (this_ *SimpleTable[R, C, V]) Clear() {
	this_.rowMap.Clear()
}

// Java method: SimpleTable.put(R rowKey, C columnKey, V value)
func (this_ *SimpleTable[R, C, V]) Put(rowKey R, columnKey C, value V) V {
	_ = rowKey
	_ = columnKey
	_ = value
	func() R {
		tv_8 := rowKey
		if _, ok_ := any((*R)(nil)).(*string); ok_ {
			if tv_8 == any(j2gort.NullString).(R) {
				j2gort.ThrowNPE()
			}
			return tv_8
		}
		var z_ R
		if tv_8 == z_ {
			j2gort.ThrowNPE()
		}
		return tv_8
	}()
	func() C {
		tv_9 := columnKey
		if _, ok_ := any((*C)(nil)).(*string); ok_ {
			if tv_9 == any(j2gort.NullString).(C) {
				j2gort.ThrowNPE()
			}
			return tv_9
		}
		var z_ C
		if tv_9 == z_ {
			j2gort.ThrowNPE()
		}
		return tv_9
	}()
	j2gort.CheckNotNullGeneric(value)
	var row j2gort.Ext_java_util_Map
iifeSwitch_2:
	switch {
	default:
		m_10 := this_.rowMap
		k_11 := j2gort.ToAny(rowKey)
		if v_12, ok_13 := j2gort.MapTryGetComputeAny(m_10, k_11); ok_13 {
			row = j2gort.AsIface[j2gort.Ext_java_util_Map](v_12)
			break iifeSwitch_2
		}
		row = j2gort.AsIface[j2gort.Ext_java_util_Map](m_10.ComputeIfAbsent(k_11, j2gort.NewLambda_Ext_java_util_function_Function(func(lp_14 any) any {
			var k R
		iifeSwitch_1:
			switch {
			default:

				tv_15 := lp_14
				if tv_15 == nil {
					{
						if _, ok_ := any((*R)(nil)).(*string); ok_ {
							k = any(j2gort.NullString).(R)
							break iifeSwitch_1
						}
						var z_ R
						k = z_
						break iifeSwitch_1
					}
				}
				k = j2gort.Cast[R](tv_15)
				break iifeSwitch_1
			}
			_ = k
			return j2gort.NewExt_java_util_LinkedHashMap()
		})))
		break iifeSwitch_2
	}

	_ = row
	{
		tv_16 := row.Put_2(j2gort.ToAny(columnKey), j2gort.ToAny(value))
		if tv_16 == nil {
			{
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					return any(j2gort.NullString).(V)
				}
				var z_ V
				return z_
			}
		}
		return j2gort.Cast[V](tv_16)
	}
}

// Java method: SimpleTable.remove(java.lang.Object rowKey, java.lang.Object columnKey)
func (this_ *SimpleTable[R, C, V]) Remove(rowKey any, columnKey any) V {
	_ = rowKey
	_ = columnKey
	var row j2gort.Ext_java_util_Map = j2gort.AsIface[j2gort.Ext_java_util_Map](this_.rowMap.Get(rowKey))
	_ = row
	if row == nil {
		{
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				return any(j2gort.NullString).(V)
			}
			var z_ V
			return z_
		}
	}
	var removed V
iifeSwitch_3:
	switch {
	default:
		tv_17 := row.Remove_1(columnKey)
		if tv_17 == nil {
			{
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					removed = any(j2gort.NullString).(V)
					break iifeSwitch_3
				}
				var z_ V
				removed = z_
				break iifeSwitch_3
			}
		}
		removed = j2gort.Cast[V](tv_17)
		break iifeSwitch_3
	}
	_ = removed
	if row.IsEmpty() {
		this_.rowMap.Remove_1(rowKey)
	}
	return removed
}

// Java method: SimpleTable.row(R rowKey)
func (this_ *SimpleTable[R, C, V]) Row(rowKey R) j2gort.Ext_java_util_Map {
	_ = rowKey
	var row j2gort.Ext_java_util_Map = j2gort.AsIface[j2gort.Ext_java_util_Map](this_.rowMap.Get(j2gort.ToAny(rowKey)))
	_ = row
	if row == nil {
		return j2gort.Ext_java_util_Collections_EmptyMap()
	}
	return row
}

// Java method: SimpleTable.column(C columnKey)
func (this_ *SimpleTable[R, C, V]) Column(columnKey C) j2gort.Ext_java_util_Map {
	_ = columnKey
	var col j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = col
	{
		cur_18 := j2gort.NewMapEntryCursor(this_.rowMap)
		for cur_18.NextEntry() {
			var val V
		iifeSwitch_4:
			switch {
			default:
				tv_19 := j2gort.AsIface[j2gort.Ext_java_util_Map](cur_18.GetValue()).Get(j2gort.ToAny(columnKey))
				if tv_19 == nil {
					{
						if _, ok_ := any((*V)(nil)).(*string); ok_ {
							val = any(j2gort.NullString).(V)
							break iifeSwitch_4
						}
						var z_ V
						val = z_
						break iifeSwitch_4
					}
				}
				val = j2gort.Cast[V](tv_19)
				break iifeSwitch_4
			}
			_ = val
			if !j2gort.IsNullGeneric(val) {
				func() V {
					tv_21 := col.Put_2(j2gort.ToAny(func() R {
						tv_20 := cur_18.GetKey()
						if tv_20 == nil {
							{
								if _, ok_ := any((*R)(nil)).(*string); ok_ {
									return any(j2gort.NullString).(R)
								}
								var z_ R
								return z_
							}
						}
						return j2gort.Cast[R](tv_20)
					}()), j2gort.ToAny(val))
					if tv_21 == nil {
						{
							if _, ok_ := any((*V)(nil)).(*string); ok_ {
								return any(j2gort.NullString).(V)
							}
							var z_ V
							return z_
						}
					}
					return j2gort.Cast[V](tv_21)
				}()
			}
		}
	}
	return col
}

// Java method: SimpleTable.cellSet()
func (this_ *SimpleTable[R, C, V]) CellSet() j2gort.Ext_java_util_Set {
	return NewSimpleTable__Anon_82[R, C, V](this_)
}

// Java method: SimpleTable.rowKeySet()
func (this_ *SimpleTable[R, C, V]) RowKeySet() j2gort.Ext_java_util_Set {
	return j2gort.AsIface[j2gort.Ext_java_util_Set](this_.rowMap.KeySet_0())
}

// Java method: SimpleTable.columnKeySet()
func (this_ *SimpleTable[R, C, V]) ColumnKeySet() j2gort.Ext_java_util_Set {
	var cols j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = cols
	{
		cur_22 := j2gort.NewMapValCursor(this_.rowMap)
		for cur_22.HasNext() {
			rawElem_23 := cur_22.Next()
			row := j2gort.AsIface[j2gort.Ext_java_util_Map](rawElem_23)
			_ = row
			cols.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Set](row.KeySet_0()))
		}
	}
	return cols
}

// Java method: SimpleTable.values()
func (this_ *SimpleTable[R, C, V]) Values() j2gort.Ext_java_util_Collection {
	var vals j2gort.Ext_java_util_Collection = j2gort.NewExt_java_util_ArrayList_1_int32(this_.Size())
	_ = vals
	{
		cur_24 := j2gort.NewMapValCursor(this_.rowMap)
		for cur_24.HasNext() {
			rawElem_25 := cur_24.Next()
			row := j2gort.AsIface[j2gort.Ext_java_util_Map](rawElem_25)
			_ = row
			vals.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](row.Values()))
		}
	}
	return vals
}

// Java method: SimpleTable.rowMap()
func (this_ *SimpleTable[R, C, V]) RowMap() j2gort.Ext_java_util_Map {
	return j2gort.Ext_java_util_Collections_UnmodifiableMap(this_.rowMap)
}

// Java method: SimpleTable.equals(java.lang.Object o)
func (this_ *SimpleTable[R, C, V]) Equals(arg_26 any) bool {
	o := arg_26
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	{
		var iifeCond_1 bool
		{

			v_, ok_ := o.(interface {
				IsGenericInst_jscomp_SimpleTable_()
				IsJavaNil_() bool
			})
			iifeCond_1 = ok_ && !v_.IsJavaNil_()
		}
		if !(iifeCond_1) {
			return false
		}
	}
	var other *SimpleTable[R, C, V] = j2gort.CastPtr[SimpleTable[R, C, V]](j2gort.CastPtr[SimpleTable[any, any, any]](o))
	_ = other
	return this_.CellSet().Equals(other.CellSet())
}

// Java method: SimpleTable.hashCode()
func (this_ *SimpleTable[R, C, V]) HashCode() int32 {
	return this_.CellSet().HashCode()
}

// Java method: SimpleTable.toString()
func (this_ *SimpleTable[R, C, V]) String() string {
	return javalang.ToJavaString(j2gort.CheckNotNull(this_.rowMap))
}

// Java interface: com.google.javascript.jscomp.SimpleTable.Cell
type SimpleTable_Cell interface {
	Implements_jscomp_SimpleTable_Cell()
	IsJavaNil_() bool
	GetRowKey() any
	GetColumnKey() any
	GetValue() any
}

func ToIface_SimpleTable_Cell[T interface {
	SimpleTable_Cell
	comparable
}](v T) SimpleTable_Cell {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SimpleTable_Cell[T SimpleTable_Cell](v T) SimpleTable_Cell {
	var vi SimpleTable_Cell = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SimpleTable.Cell", "com.google.javascript.jscomp.SimpleTable$Cell", func(v_ any) bool {
		return j2gort.InstanceOf[SimpleTable_Cell](v_)
	})
}

// Java class: com.google.javascript.jscomp.SimpleTable.SimpleCell
type SimpleTable_SimpleCell[R comparable, C comparable, V any] struct {
	rowKey    R `gson:"rowKey"`
	columnKey C `gson:"columnKey"`
	value     V `gson:"value"`
}

func (this_ *SimpleTable_SimpleCell[R, C, V]) BindSelf_(s any) {
	_ = s
}
func (this_ *SimpleTable_SimpleCell[R, C, V]) IsJavaNil_() bool { return this_ == nil }
func (*SimpleTable_SimpleCell[R, C, V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SimpleTable$SimpleCell", "SimpleCell"
}
func (*SimpleTable_SimpleCell[R, C, V]) IsGenericInst_jscomp_SimpleTable_SimpleCell_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SimpleTable.SimpleCell", "com.google.javascript.jscomp.SimpleTable$SimpleCell", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_jscomp_SimpleTable_SimpleCell_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*SimpleTable_SimpleCell[R, C, V]) Implements_jscomp_SimpleTable_Cell() {}

func StaticGet_SimpleTable_SimpleCell_SerialVersionUID() int64 { return int64(1) }
func (this_ *SimpleTable_SimpleCell[R, C, V]) initFields_SimpleTable_SimpleCell() {
	_ = this_
}

// Java constructor: SimpleCell(R rowKey, C columnKey, V value)
func NewSimpleTable_SimpleCell[R comparable, C comparable, V any](rowKey R, columnKey C, value V) *SimpleTable_SimpleCell[R, C, V] {
	this_ := &SimpleTable_SimpleCell[R, C, V]{}
	InitSimpleTable_SimpleCell[R, C, V](this_, rowKey, columnKey, value)
	return this_
}

func InitSimpleTable_SimpleCell[R comparable, C comparable, V any](this_ *SimpleTable_SimpleCell[R, C, V], rowKey R, columnKey C, value V) {
	_ = rowKey
	_ = columnKey
	_ = value
	this_.initFields_SimpleTable_SimpleCell()
	this_.rowKey = rowKey
	this_.columnKey = columnKey
	this_.value = value
}

// Java method: SimpleCell.getRowKey()
func (this_ *SimpleTable_SimpleCell[R, C, V]) GetRowKey() any {
	return j2gort.ToAny(this_.rowKey)
}

// Java method: SimpleCell.getColumnKey()
func (this_ *SimpleTable_SimpleCell[R, C, V]) GetColumnKey() any {
	return j2gort.ToAny(this_.columnKey)
}

// Java method: SimpleCell.getValue()
func (this_ *SimpleTable_SimpleCell[R, C, V]) GetValue() any {
	return j2gort.ToAny(this_.value)
}

// Java method: SimpleCell.equals(java.lang.Object o)
func (this_ *SimpleTable_SimpleCell[R, C, V]) Equals(arg_27 any) bool {
	o := arg_27
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	if !(j2gort.InstanceOf[SimpleTable_Cell](o)) {
		return false
	}
	var other SimpleTable_Cell = j2gort.AsIface[SimpleTable_Cell](o)
	_ = other
	return ((javalang.ObjectsEquals(j2gort.ToAny(this_.rowKey), other.GetRowKey()) && javalang.ObjectsEquals(j2gort.ToAny(this_.columnKey), other.GetColumnKey())) && javalang.ObjectsEquals(j2gort.ToAny(this_.value), other.GetValue()))
}

// Java method: SimpleCell.hashCode()
func (this_ *SimpleTable_SimpleCell[R, C, V]) HashCode() int32 {
	{
		var oh_28 int32 = 1
		oh_28 = 31*oh_28 + j2gort.ObjectHash(this_.rowKey)
		oh_28 = 31*oh_28 + j2gort.ObjectHash(this_.columnKey)
		oh_28 = 31*oh_28 + j2gort.ObjectHash(this_.value)
		return oh_28
	}
}

// Java method: SimpleCell.toString()
func (this_ *SimpleTable_SimpleCell[R, C, V]) String() string {
	return ("(" + j2gort.StrVal(j2gort.ToAny(this_.rowKey)) + "," + j2gort.StrVal(j2gort.ToAny(this_.columnKey)) + ")=" + j2gort.StrVal(j2gort.ToAny(this_.value)))
}

// Java class:
type SimpleTable__Anon_82[R comparable, C comparable, V any] struct {
	j2gort.Ext_java_util_AbstractSetBase
	outer_ *SimpleTable[R, C, V]
}

func (this_ *SimpleTable__Anon_82[R, C, V]) BindSelf_(s any) {
	_ = s
	this_.Ext_java_util_AbstractSetBase.BindExtSelf_(s)
}
func (this_ *SimpleTable__Anon_82[R, C, V]) IsJavaNil_() bool { return this_ == nil }
func (*SimpleTable__Anon_82[R, C, V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SimpleTable$1", ""
}
func (*SimpleTable__Anon_82[R, C, V]) IsGenericInst_jscomp_SimpleTable__Anon_82_()         {}
func (this_ *SimpleTable__Anon_82[R, C, V]) SetOuter_(o *SimpleTable[R, C, V])             { this_.outer_ = o }
func (this_ *SimpleTable__Anon_82[R, C, V]) GetOuter_() *SimpleTable[R, C, V]              { return this_.outer_ }
func (*SimpleTable__Anon_82[R, C, V]) Implements_j2gort_Ext_java_util_AbstractSet()        {}
func (*SimpleTable__Anon_82[R, C, V]) Implements_j2gort_Ext_java_util_AbstractCollection() {}
func (*SimpleTable__Anon_82[R, C, V]) Implements_j2gort_Ext_java_util_Set()                {}
func (*SimpleTable__Anon_82[R, C, V]) Implements_j2gort_Ext_java_util_Collection()         {}
func (*SimpleTable__Anon_82[R, C, V]) Implements_j2gort_Ext_java_lang_Iterable()           {}

func (this_ *SimpleTable__Anon_82[R, C, V]) initFields_SimpleTable__Anon_82() {
	_ = this_
}

func NewSimpleTable__Anon_82[R comparable, C comparable, V any](outer_ *SimpleTable[R, C, V]) *SimpleTable__Anon_82[R, C, V] {
	this_ := &SimpleTable__Anon_82[R, C, V]{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	j2gort.InitExt_java_util_AbstractSet(&this_.Ext_java_util_AbstractSetBase)
	this_.initFields_SimpleTable__Anon_82()
	return this_
}

// Java method: .iterator()
func (this_ *SimpleTable__Anon_82[R, C, V]) Iterator() any {
	var cells []*SimpleTable_SimpleCell[R, C, V] = j2gort.MakeNonNilSlice[*SimpleTable_SimpleCell[R, C, V]](this_.Size())
	_ = cells
	{
		cur_29 := j2gort.NewMapEntryCursor(this_.outer_.rowMap)
		for cur_29.NextEntry() {
			var r R
		iifeSwitch_5:
			switch {
			default:
				tv_30 := cur_29.GetKey()
				if tv_30 == nil {
					{
						if _, ok_ := any((*R)(nil)).(*string); ok_ {
							r = any(j2gort.NullString).(R)
							break iifeSwitch_5
						}
						var z_ R
						r = z_
						break iifeSwitch_5
					}
				}
				r = j2gort.Cast[R](tv_30)
				break iifeSwitch_5
			}
			_ = r
			{
				cur_31 := j2gort.NewMapEntryCursor(j2gort.AsIface[j2gort.Ext_java_util_Map](cur_29.GetValue()))
				for cur_31.NextEntry() {
					var iifeArg_1 C
				iifeSwitch_6:
					switch {
					default:

						tv_32 := cur_31.GetKey()
						if tv_32 == nil {
							{
								if _, ok_ := any((*C)(nil)).(*string); ok_ {
									iifeArg_1 = any(j2gort.NullString).(C)
									break iifeSwitch_6
								}
								var z_ C
								iifeArg_1 = z_
								break iifeSwitch_6
							}
						}
						iifeArg_1 = j2gort.Cast[C](tv_32)
						break iifeSwitch_6
					}
					var iifeArg_2 V
				iifeSwitch_7:
					switch {
					default:

						tv_33 := cur_31.GetValue()
						if tv_33 == nil {
							{
								if _, ok_ := any((*V)(nil)).(*string); ok_ {
									iifeArg_2 = any(j2gort.NullString).(V)
									break iifeSwitch_7
								}
								var z_ V
								iifeArg_2 = z_
								break iifeSwitch_7
							}
						}
						iifeArg_2 = j2gort.Cast[V](tv_33)
						break iifeSwitch_7
					}
					cells = append(cells, NewSimpleTable_SimpleCell[R, C, V](r, iifeArg_1,

						iifeArg_2,
					))
				}
			}
		}
	}
	return j2gort.NewTypedSliceIter(cells)
}

// Java method: .size()
func (this_ *SimpleTable__Anon_82[R, C, V]) Size() int32 {
	return this_.outer_.Size()
}
