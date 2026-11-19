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
// Port of com/google/javascript/jscomp/CheckLevel.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.CheckLevel
type CheckLevel struct {
	j2gort.EnumBase
}

func (this_ *CheckLevel) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *CheckLevel) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*CheckLevel) Implements_j2gort_Ext_java_lang_Enum() {}
func (*CheckLevel) Implements_j2gort_Ext_java_lang_constant_Constable() {}
var enumBacking_CheckLevel = [3]CheckLevel{
	{EnumBase: j2gort.EnumBase{Name_: "ERROR", Ordinal_: 0}},
	{EnumBase: j2gort.EnumBase{Name_: "WARNING", Ordinal_: 1}},
	{EnumBase: j2gort.EnumBase{Name_: "OFF", Ordinal_: 2}},
}
var enum_CheckLevel_ERROR *CheckLevel = &enumBacking_CheckLevel[0]
func StaticGet_CheckLevel_ERROR() *CheckLevel {
	return &enumBacking_CheckLevel[0]
}
var enum_CheckLevel_WARNING *CheckLevel = &enumBacking_CheckLevel[1]
func StaticGet_CheckLevel_WARNING() *CheckLevel {
	return &enumBacking_CheckLevel[1]
}
var enum_CheckLevel_OFF *CheckLevel = &enumBacking_CheckLevel[2]
func StaticGet_CheckLevel_OFF() *CheckLevel {
	return &enumBacking_CheckLevel[2]
}
var enumValues_CheckLevel = []*CheckLevel{&enumBacking_CheckLevel[0], &enumBacking_CheckLevel[1], &enumBacking_CheckLevel[2]}
func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.CheckLevel", "com.google.javascript.jscomp.CheckLevel", func() {
		vals_ := make([]any, len(enumValues_CheckLevel))
		for i_, v_ := range enumValues_CheckLevel {
			vals_[i_] = v_
		}
		j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CheckLevel", vals_)
	})
}

func CheckLevel_Values() []*CheckLevel {
	return j2gort.CloneSlice(enumValues_CheckLevel)
}

func CheckLevel_ValueOf(name string) *CheckLevel {
	for _, c := range enumValues_CheckLevel {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CheckLevel." + name))
}

// Java method: CheckLevel.isOn()
func (this_ *CheckLevel) IsOn() bool {
	return (this_ != StaticGet_CheckLevel_OFF())
}

func (this_ *CheckLevel) IsJavaNil_() bool { return this_ == nil }
func (*CheckLevel) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.CheckLevel", "CheckLevel" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckLevel", "com.google.javascript.jscomp.CheckLevel", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckLevel](v_)
	})
}
