/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ErrorFormat.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.ErrorFormat
type ErrorFormat struct {
	j2gort.EnumBase
	fn_ToFormatter func(*ErrorFormat, SourceExcerptProvider, bool) MessageFormatter
}

func (this_ *ErrorFormat) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *ErrorFormat) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*ErrorFormat) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*ErrorFormat) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *ErrorFormat) initFields_ErrorFormat() {
	_ = this_
}

func initEnumErrorFormat(this_ *ErrorFormat, enumName_ string, enumOrd_ int32) *ErrorFormat {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_ErrorFormat()
	return this_
}

var enum_ErrorFormat_SINGLELINE *ErrorFormat

func StaticGet_ErrorFormat_SINGLELINE() *ErrorFormat {
	clinit_ErrorFormat()
	return enum_ErrorFormat_SINGLELINE
}

var enum_ErrorFormat_FULL *ErrorFormat

func StaticGet_ErrorFormat_FULL() *ErrorFormat {
	clinit_ErrorFormat()
	return enum_ErrorFormat_FULL
}

var enum_ErrorFormat_MULTILINE *ErrorFormat

func StaticGet_ErrorFormat_MULTILINE() *ErrorFormat {
	clinit_ErrorFormat()
	return enum_ErrorFormat_MULTILINE
}

var enum_ErrorFormat_SOURCELESS *ErrorFormat

func StaticGet_ErrorFormat_SOURCELESS() *ErrorFormat {
	clinit_ErrorFormat()
	return enum_ErrorFormat_SOURCELESS
}

var clinitDone_ErrorFormat uint32
var clinitState_ErrorFormat j2gort.ClinitState

func clinit_ErrorFormat() {
	if atomic.LoadUint32(&clinitDone_ErrorFormat) != 0 {
		return
	}
	clinitSlow_ErrorFormat()
}

func clinitSlow_ErrorFormat() {
	if !j2gort.EnterClinit(&clinitState_ErrorFormat, &clinitDone_ErrorFormat) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ErrorFormat, &clinitDone_ErrorFormat)
	enum_ErrorFormat_SINGLELINE = &ErrorFormat{}
	enum_ErrorFormat_SINGLELINE.fn_ToFormatter = func(this_ *ErrorFormat, source SourceExcerptProvider, colorize bool) MessageFormatter {
		_ = source
		_ = colorize
		var formatter *LightweightMessageFormatter = NewLightweightMessageFormatter_1(source)
		_ = formatter
		formatter.SetColorize(colorize)
		return formatter
	}
	initEnumErrorFormat(enum_ErrorFormat_SINGLELINE, "SINGLELINE", int32(0))
	enum_ErrorFormat_FULL = &ErrorFormat{}
	enum_ErrorFormat_FULL.fn_ToFormatter = func(this_ *ErrorFormat, source SourceExcerptProvider, colorize bool) MessageFormatter {
		_ = source
		_ = colorize
		var formatter *LightweightMessageFormatter = NewLightweightMessageFormatter_2(source, StaticGet_SourceExcerptProvider_SourceExcerpt_FULL())
		_ = formatter
		formatter.SetColorize(colorize)
		return formatter
	}
	initEnumErrorFormat(enum_ErrorFormat_FULL, "FULL", int32(1))
	enum_ErrorFormat_MULTILINE = &ErrorFormat{}
	enum_ErrorFormat_MULTILINE.fn_ToFormatter = func(this_ *ErrorFormat, source SourceExcerptProvider, colorize bool) MessageFormatter {
		_ = source
		_ = colorize
		var formatter *LightweightMessageFormatter = NewLightweightMessageFormatter_2(source, StaticGet_SourceExcerptProvider_SourceExcerpt_REGION())
		_ = formatter
		formatter.SetColorize(colorize)
		return formatter
	}
	initEnumErrorFormat(enum_ErrorFormat_MULTILINE, "MULTILINE", int32(2))
	enum_ErrorFormat_SOURCELESS = &ErrorFormat{}
	enum_ErrorFormat_SOURCELESS.fn_ToFormatter = func(this_ *ErrorFormat, source SourceExcerptProvider, colorize bool) MessageFormatter {
		_ = source
		_ = colorize
		var formatter *LightweightMessageFormatter = LightweightMessageFormatter_WithoutSource()
		_ = formatter
		formatter.SetColorize(colorize)
		return ToIface_MessageFormatter(formatter)
	}
	initEnumErrorFormat(enum_ErrorFormat_SOURCELESS, "SOURCELESS", int32(3))
	enumValues_ErrorFormat = []*ErrorFormat{enum_ErrorFormat_SINGLELINE, enum_ErrorFormat_FULL, enum_ErrorFormat_MULTILINE, enum_ErrorFormat_SOURCELESS}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ErrorFormat", []any{enum_ErrorFormat_SINGLELINE, enum_ErrorFormat_FULL, enum_ErrorFormat_MULTILINE, enum_ErrorFormat_SOURCELESS})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.ErrorFormat", "com.google.javascript.jscomp.ErrorFormat", func() {
		clinit_ErrorFormat()
	})
}

var enumValues_ErrorFormat []*ErrorFormat

func ErrorFormat_Values() []*ErrorFormat {
	clinit_ErrorFormat()
	return j2gort.CloneSlice(enumValues_ErrorFormat)
}

func ErrorFormat_ValueOf(name string) *ErrorFormat {
	clinit_ErrorFormat()
	for _, c := range enumValues_ErrorFormat {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ErrorFormat." + name))
}

// Java method: ErrorFormat.toFormatter(com.google.javascript.jscomp.SourceExcerptProvider source, boolean colorize)
func (this_ *ErrorFormat) ToFormatter(source SourceExcerptProvider, colorize bool) MessageFormatter {
	if this_.fn_ToFormatter != nil {
		return this_.fn_ToFormatter(this_, source, colorize)
	}
	return this_.baseEnum_ToFormatter(source, colorize)
}

func (this_ *ErrorFormat) baseEnum_ToFormatter(source SourceExcerptProvider, colorize bool) MessageFormatter {
	_ = source
	_ = colorize
	panic("abstract or native method invoked")
}

func (this_ *ErrorFormat) IsJavaNil_() bool { return this_ == nil }
func (*ErrorFormat) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ErrorFormat", "ErrorFormat"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ErrorFormat", "com.google.javascript.jscomp.ErrorFormat", func(v_ any) bool {
		return j2gort.InstanceOf[*ErrorFormat](v_)
	})
}
