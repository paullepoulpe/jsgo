/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/SourceExcerptProvider.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java interface: com.google.javascript.jscomp.SourceExcerptProvider
type SourceExcerptProvider interface {
	Implements_jscomp_SourceExcerptProvider()
	IsJavaNil_() bool
	GetSourceLine(string, int32) string
	GetSourceLines(string, int32, int32) Region
	GetSourceRegion(string, int32) Region
	GetSourceMapping(string, int32, int32) *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping
}
func ToIface_SourceExcerptProvider[T interface {
	SourceExcerptProvider
	comparable
}](v T) SourceExcerptProvider {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SourceExcerptProvider[T SourceExcerptProvider](v T) SourceExcerptProvider {
	var vi SourceExcerptProvider = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceExcerptProvider", "com.google.javascript.jscomp.SourceExcerptProvider", func(v_ any) bool {
		return j2gort.InstanceOf[SourceExcerptProvider](v_)
	})
}
// Java enum: com.google.javascript.jscomp.SourceExcerptProvider.SourceExcerpt
type SourceExcerptProvider_SourceExcerpt struct {
	j2gort.EnumBase
	fn_Get_5 func(*SourceExcerptProvider_SourceExcerpt, SourceExcerptProvider, string, int32, int32, SourceExcerptProvider_ExcerptFormatter) string
}

func (this_ *SourceExcerptProvider_SourceExcerpt) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *SourceExcerptProvider_SourceExcerpt) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*SourceExcerptProvider_SourceExcerpt) Implements_j2gort_Ext_java_lang_Enum() {}
func (*SourceExcerptProvider_SourceExcerpt) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *SourceExcerptProvider_SourceExcerpt) initFields_SourceExcerptProvider_SourceExcerpt() {
	_ = this_
}

func initEnumSourceExcerptProvider_SourceExcerpt(this_ *SourceExcerptProvider_SourceExcerpt, enumName_ string, enumOrd_ int32) *SourceExcerptProvider_SourceExcerpt {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_SourceExcerptProvider_SourceExcerpt()
	return this_
}

var enum_SourceExcerptProvider_SourceExcerpt_LINE *SourceExcerptProvider_SourceExcerpt
func StaticGet_SourceExcerptProvider_SourceExcerpt_LINE() *SourceExcerptProvider_SourceExcerpt {
	clinit_SourceExcerptProvider_SourceExcerpt()
	return enum_SourceExcerptProvider_SourceExcerpt_LINE
}
var enum_SourceExcerptProvider_SourceExcerpt_FULL *SourceExcerptProvider_SourceExcerpt
func StaticGet_SourceExcerptProvider_SourceExcerpt_FULL() *SourceExcerptProvider_SourceExcerpt {
	clinit_SourceExcerptProvider_SourceExcerpt()
	return enum_SourceExcerptProvider_SourceExcerpt_FULL
}
var enum_SourceExcerptProvider_SourceExcerpt_REGION *SourceExcerptProvider_SourceExcerpt
func StaticGet_SourceExcerptProvider_SourceExcerpt_REGION() *SourceExcerptProvider_SourceExcerpt {
	clinit_SourceExcerptProvider_SourceExcerpt()
	return enum_SourceExcerptProvider_SourceExcerpt_REGION
}
var clinitDone_SourceExcerptProvider_SourceExcerpt uint32
var clinitState_SourceExcerptProvider_SourceExcerpt j2gort.ClinitState
func clinit_SourceExcerptProvider_SourceExcerpt() {
	if atomic.LoadUint32(&clinitDone_SourceExcerptProvider_SourceExcerpt) != 0 {
		return
	}
	clinitSlow_SourceExcerptProvider_SourceExcerpt()
}

func clinitSlow_SourceExcerptProvider_SourceExcerpt() {
	if !j2gort.EnterClinit(&clinitState_SourceExcerptProvider_SourceExcerpt, &clinitDone_SourceExcerptProvider_SourceExcerpt) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SourceExcerptProvider_SourceExcerpt, &clinitDone_SourceExcerptProvider_SourceExcerpt)
	enum_SourceExcerptProvider_SourceExcerpt_LINE = &SourceExcerptProvider_SourceExcerpt{}
	enum_SourceExcerptProvider_SourceExcerpt_LINE.fn_Get_5 = func(this_ *SourceExcerptProvider_SourceExcerpt, source SourceExcerptProvider, sourceName string, lineNumber int32, unused int32, formatter SourceExcerptProvider_ExcerptFormatter) string {
	_ = source
	_ = sourceName
	_ = lineNumber
	_ = unused
	_ = formatter
	return formatter.FormatLine(source.GetSourceLine(sourceName, lineNumber), lineNumber)
}
	initEnumSourceExcerptProvider_SourceExcerpt(enum_SourceExcerptProvider_SourceExcerpt_LINE, "LINE", int32(0))
	enum_SourceExcerptProvider_SourceExcerpt_FULL = &SourceExcerptProvider_SourceExcerpt{}
	enum_SourceExcerptProvider_SourceExcerpt_FULL.fn_Get_5 = func(this_ *SourceExcerptProvider_SourceExcerpt, source SourceExcerptProvider, sourceName string, startLineNumber int32, length int32, formatter SourceExcerptProvider_ExcerptFormatter) string {
	_ = source
	_ = sourceName
	_ = startLineNumber
	_ = length
	_ = formatter
	return formatter.FormatRegion(source.GetSourceLines(sourceName, startLineNumber, length))
}
	initEnumSourceExcerptProvider_SourceExcerpt(enum_SourceExcerptProvider_SourceExcerpt_FULL, "FULL", int32(1))
	enum_SourceExcerptProvider_SourceExcerpt_REGION = &SourceExcerptProvider_SourceExcerpt{}
	enum_SourceExcerptProvider_SourceExcerpt_REGION.fn_Get_5 = func(this_ *SourceExcerptProvider_SourceExcerpt, source SourceExcerptProvider, sourceName string, lineNumber int32, length int32, formatter SourceExcerptProvider_ExcerptFormatter) string {
	_ = source
	_ = sourceName
	_ = lineNumber
	_ = length
	_ = formatter
	return formatter.FormatRegion(source.GetSourceRegion(sourceName, lineNumber))
}
	initEnumSourceExcerptProvider_SourceExcerpt(enum_SourceExcerptProvider_SourceExcerpt_REGION, "REGION", int32(2))
	enumValues_SourceExcerptProvider_SourceExcerpt = []*SourceExcerptProvider_SourceExcerpt{enum_SourceExcerptProvider_SourceExcerpt_LINE, enum_SourceExcerptProvider_SourceExcerpt_FULL, enum_SourceExcerptProvider_SourceExcerpt_REGION}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.SourceExcerptProvider.SourceExcerpt", []any{enum_SourceExcerptProvider_SourceExcerpt_LINE, enum_SourceExcerptProvider_SourceExcerpt_FULL, enum_SourceExcerptProvider_SourceExcerpt_REGION})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.SourceExcerptProvider$SourceExcerpt", []any{enum_SourceExcerptProvider_SourceExcerpt_LINE, enum_SourceExcerptProvider_SourceExcerpt_FULL, enum_SourceExcerptProvider_SourceExcerpt_REGION})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.SourceExcerptProvider.SourceExcerpt", "com.google.javascript.jscomp.SourceExcerptProvider$SourceExcerpt", func() {
		clinit_SourceExcerptProvider_SourceExcerpt()
	})
}

var enumValues_SourceExcerptProvider_SourceExcerpt []*SourceExcerptProvider_SourceExcerpt
func SourceExcerptProvider_SourceExcerpt_Values() []*SourceExcerptProvider_SourceExcerpt {
	clinit_SourceExcerptProvider_SourceExcerpt()
	return j2gort.CloneSlice(enumValues_SourceExcerptProvider_SourceExcerpt)
}

func SourceExcerptProvider_SourceExcerpt_ValueOf(name string) *SourceExcerptProvider_SourceExcerpt {
	clinit_SourceExcerptProvider_SourceExcerpt()
	for _, c := range enumValues_SourceExcerptProvider_SourceExcerpt {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.SourceExcerptProvider.SourceExcerpt." + name))
}

// Java method: SourceExcerpt.get(com.google.javascript.jscomp.SourceExcerptProvider source, java.lang.String sourceName, int lineNumber, com.google.javascript.jscomp.SourceExcerptProvider.ExcerptFormatter formatter)
func (this_ *SourceExcerptProvider_SourceExcerpt) Get_4(source SourceExcerptProvider, sourceName string, lineNumber int32, formatter SourceExcerptProvider_ExcerptFormatter) string {
	_ = source
	_ = sourceName
	_ = lineNumber
	_ = formatter
	return this_.Get_5(source, sourceName, lineNumber, int32(-1), formatter)
}

// Java method: SourceExcerpt.get(com.google.javascript.jscomp.SourceExcerptProvider source, java.lang.String sourceName, int lineNumber, int length, com.google.javascript.jscomp.SourceExcerptProvider.ExcerptFormatter formatter)
func (this_ *SourceExcerptProvider_SourceExcerpt) Get_5(source SourceExcerptProvider, sourceName string, lineNumber int32, length int32, formatter SourceExcerptProvider_ExcerptFormatter) string {
	if this_.fn_Get_5 != nil {
		return this_.fn_Get_5(this_, source, sourceName, lineNumber, length, formatter)
	}
	return this_.baseEnum_Get_5(source, sourceName, lineNumber, length, formatter)
}

func (this_ *SourceExcerptProvider_SourceExcerpt) baseEnum_Get_5(source SourceExcerptProvider, sourceName string, lineNumber int32, length int32, formatter SourceExcerptProvider_ExcerptFormatter) string {
	_ = source
	_ = sourceName
	_ = lineNumber
	_ = length
	_ = formatter
	panic("abstract or native method invoked")
}

func (this_ *SourceExcerptProvider_SourceExcerpt) IsJavaNil_() bool { return this_ == nil }
func (*SourceExcerptProvider_SourceExcerpt) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.SourceExcerptProvider$SourceExcerpt", "SourceExcerpt" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceExcerptProvider.SourceExcerpt", "com.google.javascript.jscomp.SourceExcerptProvider$SourceExcerpt", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceExcerptProvider_SourceExcerpt](v_)
	})
}
// Java interface: com.google.javascript.jscomp.SourceExcerptProvider.ExcerptFormatter
type SourceExcerptProvider_ExcerptFormatter interface {
	Implements_jscomp_SourceExcerptProvider_ExcerptFormatter()
	IsJavaNil_() bool
	FormatLine(string, int32) string
	FormatRegion(Region) string
}
func ToIface_SourceExcerptProvider_ExcerptFormatter[T interface {
	SourceExcerptProvider_ExcerptFormatter
	comparable
}](v T) SourceExcerptProvider_ExcerptFormatter {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SourceExcerptProvider_ExcerptFormatter[T SourceExcerptProvider_ExcerptFormatter](v T) SourceExcerptProvider_ExcerptFormatter {
	var vi SourceExcerptProvider_ExcerptFormatter = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceExcerptProvider.ExcerptFormatter", "com.google.javascript.jscomp.SourceExcerptProvider$ExcerptFormatter", func(v_ any) bool {
		return j2gort.InstanceOf[SourceExcerptProvider_ExcerptFormatter](v_)
	})
}
