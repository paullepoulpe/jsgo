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
// Port of com/google/javascript/jscomp/VerboseMessageFormatter.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.VerboseMessageFormatter
type VerboseMessageFormatter struct {
	AbstractMessageFormatterBase
}

func (this_ *VerboseMessageFormatter) BindSelf_(s any) {
	_ = s
}
func (this_ *VerboseMessageFormatter) IsJavaNil_() bool { return this_ == nil }
func (*VerboseMessageFormatter) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.VerboseMessageFormatter", "VerboseMessageFormatter" }
func (this_ *VerboseMessageFormatter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractMessageFormatterBase.IdentityHash_()
}
func (*VerboseMessageFormatter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.VerboseMessageFormatter", "com.google.javascript.jscomp.VerboseMessageFormatter", func(v_ any) bool {
		return j2gort.InstanceOf[*VerboseMessageFormatter](v_)
	})
}
func (*VerboseMessageFormatter) Implements_jscomp_AbstractMessageFormatter() {}
func (*VerboseMessageFormatter) Implements_jscomp_MessageFormatter() {}

func (this_ *VerboseMessageFormatter) initFields_VerboseMessageFormatter() {
	_ = this_
}

// Java constructor: VerboseMessageFormatter(com.google.javascript.jscomp.SourceExcerptProvider source)
func NewVerboseMessageFormatter(source SourceExcerptProvider) *VerboseMessageFormatter {
	this_ := &VerboseMessageFormatter{}
	InitVerboseMessageFormatter(this_, source)
	return this_
}

func InitVerboseMessageFormatter(this_ *VerboseMessageFormatter, source SourceExcerptProvider) {
	_ = source
	InitAbstractMessageFormatter(&this_.AbstractMessageFormatterBase, source)
	this_.initFields_VerboseMessageFormatter()
}

// Java method: VerboseMessageFormatter.formatError(com.google.javascript.jscomp.JSError error)
func (this_ *VerboseMessageFormatter) FormatError(error_ *JSError) string {
	_ = error_
	return (j2gort.StrValString(this_.GetLevelName(StaticGet_CheckLevel_ERROR())) + ": " + j2gort.StrValString(this_.priv_VerboseMessageFormatter_format_5864(error_)))
}

// Java method: VerboseMessageFormatter.formatWarning(com.google.javascript.jscomp.JSError warning)
func (this_ *VerboseMessageFormatter) FormatWarning(warning *JSError) string {
	_ = warning
	return (j2gort.StrValString(this_.GetLevelName(StaticGet_CheckLevel_WARNING())) + ": " + j2gort.StrValString(this_.priv_VerboseMessageFormatter_format_5864(warning)))
}

// Java method: VerboseMessageFormatter.format(com.google.javascript.jscomp.JSError message)
func (this_ *VerboseMessageFormatter) priv_VerboseMessageFormatter_format_5864(message *JSError) string {
	_ = message
	var description string = message.Description()
	_ = description
	var sourceName string = message.SourceName()
	_ = sourceName
	var lineNumber int32 = message.GetLineNumber()
	_ = lineNumber
	var sourceRegion Region = this_.GetSource().GetSourceRegion(sourceName, lineNumber)
	_ = sourceRegion
	var lineSource string = j2gort.NullString
	_ = lineSource
	if ((sourceRegion != nil && !sourceRegion.IsJavaNil_())) {
		lineSource = sourceRegion.GetSourceExcerpt()
	}
	return (j2gort.StrValString(description) + " at " + j2gort.StrValString(func() string { if j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(sourceName) { return "(unknown source)" }; return sourceName }()) + " line " + j2gort.StrValString(func() string { if ((lineNumber < int32(0))) { return "(unknown line)" }; return strconv.FormatInt(int64(lineNumber), 10) }()) + " " + j2gort.StrValString(func() string { if ((lineSource != j2gort.NullString)) { return (":\n\n" + j2gort.StrValString(lineSource)) }; return "." }()))
}

