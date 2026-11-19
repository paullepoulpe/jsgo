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
// Port of com/google/javascript/jscomp/AbstractMessageFormatter.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.AbstractMessageFormatter
type AbstractMessageFormatter interface {
	Implements_jscomp_AbstractMessageFormatter()
	IsJavaNil_() bool
	GetAbstractMessageFormatterBase() *AbstractMessageFormatterBase
	MessageFormatter
	SetColorize(bool)
	GetSource() SourceExcerptProvider
	GetLevelName(*CheckLevel) string
	MaybeEmbolden(string) string
}

func ToIface_AbstractMessageFormatter[T interface {
	AbstractMessageFormatter
	comparable
}](v T) AbstractMessageFormatter {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.AbstractMessageFormatter
type AbstractMessageFormatterBase struct {
	source   SourceExcerptProvider `gson:"source"`
	idHash_  atomic.Int32
	colorize bool `gson:"colorize"`
}

func (this_ *AbstractMessageFormatterBase) BindSelf_(s any) {
	_ = s
}
func (this_ *AbstractMessageFormatterBase) IsJavaNil_() bool { return this_ == nil }
func (*AbstractMessageFormatterBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractMessageFormatter", "AbstractMessageFormatter"
}
func (this_ *AbstractMessageFormatterBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractMessageFormatter", "com.google.javascript.jscomp.AbstractMessageFormatter", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractMessageFormatter](v_)
	})
}
func (b *AbstractMessageFormatterBase) GetAbstractMessageFormatterBase() *AbstractMessageFormatterBase {
	return b
}
func (*AbstractMessageFormatterBase) Implements_jscomp_AbstractMessageFormatter() {}
func (*AbstractMessageFormatterBase) Implements_jscomp_MessageFormatter()         {}

var static_AbstractMessageFormatter_SUPPORTED_COLOR_TERMINALS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_AbstractMessageFormatter_SUPPORTED_COLOR_TERMINALS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_AbstractMessageFormatter()
	return static_AbstractMessageFormatter_SUPPORTED_COLOR_TERMINALS
}
func StaticSet_AbstractMessageFormatter_SUPPORTED_COLOR_TERMINALS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_AbstractMessageFormatter()
	static_AbstractMessageFormatter_SUPPORTED_COLOR_TERMINALS = v
}

var clinitDone_AbstractMessageFormatter uint32
var clinitState_AbstractMessageFormatter j2gort.ClinitState

func clinit_AbstractMessageFormatter() {
	if atomic.LoadUint32(&clinitDone_AbstractMessageFormatter) != 0 {
		return
	}
	clinitSlow_AbstractMessageFormatter()
}

func clinitSlow_AbstractMessageFormatter() {
	if !j2gort.EnterClinit(&clinitState_AbstractMessageFormatter, &clinitDone_AbstractMessageFormatter) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AbstractMessageFormatter, &clinitDone_AbstractMessageFormatter)
	static_AbstractMessageFormatter_SUPPORTED_COLOR_TERMINALS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("xterm", "xterm-color", "xterm-256color", "screen-bce")
}

func (this_ *AbstractMessageFormatterBase) initFields_AbstractMessageFormatter() {
	_ = this_
}

func InitAbstractMessageFormatter(this_ *AbstractMessageFormatterBase, source SourceExcerptProvider) {
	_ = source
	this_.initFields_AbstractMessageFormatter()
	this_.source = source
}

// Java method: AbstractMessageFormatter.setColorize(boolean colorize)
func (this_ *AbstractMessageFormatterBase) SetColorize(colorize bool) {
	_ = colorize
	this_.colorize = colorize
}

// Java method: AbstractMessageFormatter.getSource()
func (this_ *AbstractMessageFormatterBase) GetSource() SourceExcerptProvider {
	return this_.source
}

// Java static method: AbstractMessageFormatter.termSupportsColor(java.lang.String term)
func AbstractMessageFormatter_TermSupportsColor(term string) bool {
	_ = term
	return j2gort.SetContainsStr(StaticGet_AbstractMessageFormatter_SUPPORTED_COLOR_TERMINALS(), term)
}

// Java method: AbstractMessageFormatter.getLevelName(com.google.javascript.jscomp.CheckLevel level)
func (this_ *AbstractMessageFormatterBase) GetLevelName(level *CheckLevel) string {
	_ = level
	switch j2gort.CheckNotNullPtr(level).Ordinal() {
	case int32(0) /* CheckLevel.ERROR */ :
		return this_.priv_AbstractMessageFormatter_maybeColorize_414("ERROR", StaticGet_AbstractMessageFormatter_Color_ERROR())
	case int32(1) /* CheckLevel.WARNING */ :
		return this_.priv_AbstractMessageFormatter_maybeColorize_414("WARNING", StaticGet_AbstractMessageFormatter_Color_WARNING())
	default:
		return level.String()
	}
	panic("unreachable switch expression")
}

// Java method: AbstractMessageFormatter.maybeEmbolden(java.lang.String text)
func (this_ *AbstractMessageFormatterBase) MaybeEmbolden(text string) string {
	_ = text
	if !(this_.colorize) {
		return text
	}
	return (j2gort.StrValString(StaticGet_AbstractMessageFormatter_Color_BOLD().GetControlCharacter()) + j2gort.StrValString(text) + j2gort.StrValString(StaticGet_AbstractMessageFormatter_Color_UNBOLD().GetControlCharacter()))
}

// Java method: AbstractMessageFormatter.maybeColorize(java.lang.String text, com.google.javascript.jscomp.AbstractMessageFormatter.Color color)
func (this_ *AbstractMessageFormatterBase) priv_AbstractMessageFormatter_maybeColorize_414(text string, color AbstractMessageFormatter_Color) string {
	_ = text
	_ = color
	if !(this_.colorize) {
		return text
	}
	return (j2gort.StrValString(color.GetControlCharacter()) + j2gort.StrValString(text) + j2gort.StrValString(StaticGet_AbstractMessageFormatter_Color_NO_COLOR().GetControlCharacter()))
}

// Java enum: com.google.javascript.jscomp.AbstractMessageFormatter.Color
type AbstractMessageFormatter_Color uint8

const (
	enum_AbstractMessageFormatter_Color_ERROR    AbstractMessageFormatter_Color = 1
	enum_AbstractMessageFormatter_Color_WARNING  AbstractMessageFormatter_Color = 2
	enum_AbstractMessageFormatter_Color_NO_COLOR AbstractMessageFormatter_Color = 3
	enum_AbstractMessageFormatter_Color_BOLD     AbstractMessageFormatter_Color = 4
	enum_AbstractMessageFormatter_Color_UNBOLD   AbstractMessageFormatter_Color = 5
)

var enumBaseTable_AbstractMessageFormatter_Color = [...]j2gort.EnumBase{
	{Name_: "ERROR", Ordinal_: 0},
	{Name_: "WARNING", Ordinal_: 1},
	{Name_: "NO_COLOR", Ordinal_: 2},
	{Name_: "BOLD", Ordinal_: 3},
	{Name_: "UNBOLD", Ordinal_: 4},
}
var enumValues_AbstractMessageFormatter_Color = []AbstractMessageFormatter_Color{1, 2, 3, 4, 5}

var enumField_AbstractMessageFormatter_Color_controlCharacter = [...]string{
	"\x1b[31m",
	"\x1b[35m",
	"\x1b[39m",
	"\x1b[1m",
	"\x1b[0m",
}

func (this_ AbstractMessageFormatter_Color) FieldGet_controlCharacter() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_AbstractMessageFormatter_Color_controlCharacter[this_-1]
}

func StaticGet_AbstractMessageFormatter_Color_ERROR() AbstractMessageFormatter_Color {
	return enum_AbstractMessageFormatter_Color_ERROR
}
func StaticGet_AbstractMessageFormatter_Color_WARNING() AbstractMessageFormatter_Color {
	return enum_AbstractMessageFormatter_Color_WARNING
}
func StaticGet_AbstractMessageFormatter_Color_NO_COLOR() AbstractMessageFormatter_Color {
	return enum_AbstractMessageFormatter_Color_NO_COLOR
}
func StaticGet_AbstractMessageFormatter_Color_BOLD() AbstractMessageFormatter_Color {
	return enum_AbstractMessageFormatter_Color_BOLD
}
func StaticGet_AbstractMessageFormatter_Color_UNBOLD() AbstractMessageFormatter_Color {
	return enum_AbstractMessageFormatter_Color_UNBOLD
}

func AbstractMessageFormatter_Color_Values() []AbstractMessageFormatter_Color {
	return j2gort.CloneSlice(enumValues_AbstractMessageFormatter_Color)
}

func AbstractMessageFormatter_Color_ValueOf(name string) AbstractMessageFormatter_Color {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_AbstractMessageFormatter_Color {
		if enumBaseTable_AbstractMessageFormatter_Color[i_].Name_ == name {
			return AbstractMessageFormatter_Color(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.AbstractMessageFormatter.Color." + name))
}

// Java method: Color.getControlCharacter()
func (this_ AbstractMessageFormatter_Color) GetControlCharacter() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.FieldGet_controlCharacter()
}

func (this_ AbstractMessageFormatter_Color) IsJavaNil_() bool                              { return this_ == 0 }
func (AbstractMessageFormatter_Color) HasIdentityEquals_()                                 {}
func (AbstractMessageFormatter_Color) Implements_j2gort_Ext_java_lang_Enum()               {}
func (AbstractMessageFormatter_Color) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ AbstractMessageFormatter_Color) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (AbstractMessageFormatter_Color) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractMessageFormatter$Color", "Color"
}
func (this_ AbstractMessageFormatter_Color) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_AbstractMessageFormatter_Color[this_-1]
}
func (this_ AbstractMessageFormatter_Color) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ AbstractMessageFormatter_Color) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ AbstractMessageFormatter_Color) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AbstractMessageFormatter_Color[this_-1].Name_
}
func (this_ AbstractMessageFormatter_Color) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_AbstractMessageFormatter_Color[this_-1].Name_
}
func (this_ AbstractMessageFormatter_Color) Equals(other_ any) bool {
	o_, ok_ := other_.(AbstractMessageFormatter_Color)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ AbstractMessageFormatter_Color) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_AbstractMessageFormatter_Color[this_-1].HashCode()
}
func (this_ AbstractMessageFormatter_Color) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AbstractMessageFormatter_Color[this_-1].CompareTo(other_)
}
func (this_ AbstractMessageFormatter_Color) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_AbstractMessageFormatter_Color[this_-1].Name_)
}
func init() {
	vals_ := []any{AbstractMessageFormatter_Color(1), AbstractMessageFormatter_Color(2), AbstractMessageFormatter_Color(3), AbstractMessageFormatter_Color(4), AbstractMessageFormatter_Color(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractMessageFormatter.Color", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractMessageFormatter$Color", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractMessageFormatter.Color", "com.google.javascript.jscomp.AbstractMessageFormatter$Color", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractMessageFormatter_Color](v_)
	})
}
