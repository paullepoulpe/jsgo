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
// Port of com/google/javascript/jscomp/RenamingToken.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.RenamingToken
type RenamingToken struct {
	j2gort.EnumBase
}

func (this_ *RenamingToken) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *RenamingToken) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*RenamingToken) Implements_j2gort_Ext_java_lang_Enum() {}
func (*RenamingToken) Implements_jscomp_RenamingMap() {}
func (*RenamingToken) Implements_j2gort_Ext_java_lang_constant_Constable() {}
var enumBacking_RenamingToken = [3]RenamingToken{
	{EnumBase: j2gort.EnumBase{Name_: "INCONSISTENT", Ordinal_: 0}},
	{EnumBase: j2gort.EnumBase{Name_: "STABLE", Ordinal_: 1}},
	{EnumBase: j2gort.EnumBase{Name_: "DISABLE", Ordinal_: 2}},
}
var enum_RenamingToken_INCONSISTENT *RenamingToken = &enumBacking_RenamingToken[0]
func StaticGet_RenamingToken_INCONSISTENT() *RenamingToken {
	return &enumBacking_RenamingToken[0]
}
var enum_RenamingToken_STABLE *RenamingToken = &enumBacking_RenamingToken[1]
func StaticGet_RenamingToken_STABLE() *RenamingToken {
	return &enumBacking_RenamingToken[1]
}
var enum_RenamingToken_DISABLE *RenamingToken = &enumBacking_RenamingToken[2]
func StaticGet_RenamingToken_DISABLE() *RenamingToken {
	return &enumBacking_RenamingToken[2]
}
var enumValues_RenamingToken = []*RenamingToken{&enumBacking_RenamingToken[0], &enumBacking_RenamingToken[1], &enumBacking_RenamingToken[2]}
func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.RenamingToken", "com.google.javascript.jscomp.RenamingToken", func() {
		vals_ := make([]any, len(enumValues_RenamingToken))
		for i_, v_ := range enumValues_RenamingToken {
			vals_[i_] = v_
		}
		j2gort.RegisterEnumConstants("com.google.javascript.jscomp.RenamingToken", vals_)
	})
}

func RenamingToken_Values() []*RenamingToken {
	return j2gort.CloneSlice(enumValues_RenamingToken)
}

func RenamingToken_ValueOf(name string) *RenamingToken {
	for _, c := range enumValues_RenamingToken {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.RenamingToken." + name))
}

// Java method: RenamingToken.get(java.lang.String value)
func (this_ *RenamingToken) Get(value string) string {
	_ = value
	return j2gort.NullString
}

func (this_ *RenamingToken) IsJavaNil_() bool { return this_ == nil }
func (*RenamingToken) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.RenamingToken", "RenamingToken" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenamingToken", "com.google.javascript.jscomp.RenamingToken", func(v_ any) bool {
		return j2gort.InstanceOf[*RenamingToken](v_)
	})
}
