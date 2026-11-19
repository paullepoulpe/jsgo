/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/TypeStringBuilder.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.TypeStringBuilder
type TypeStringBuilder struct {
	builder          *j2gort.Ext_java_lang_StringBuilder `gson:"builder"`
	indentation      string                              `gson:"indentation"`
	idHash_          atomic.Int32
	isForAnnotations bool `gson:"isForAnnotations"`
}

func (this_ *TypeStringBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeStringBuilder) IsJavaNil_() bool { return this_ == nil }
func (*TypeStringBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.TypeStringBuilder", "TypeStringBuilder"
}
func (this_ *TypeStringBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeStringBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.TypeStringBuilder", "com.google.javascript.rhino.jstype.TypeStringBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeStringBuilder](v_)
	})
}

func (this_ *TypeStringBuilder) initFields_TypeStringBuilder() {
	_ = this_
	this_.builder = j2gort.NewExt_java_lang_StringBuilder()

}

// Java constructor: TypeStringBuilder(boolean isForAnnotations)
func NewTypeStringBuilder(isForAnnotations bool) *TypeStringBuilder {
	this_ := &TypeStringBuilder{}
	InitTypeStringBuilder(this_, isForAnnotations)
	return this_
}

func InitTypeStringBuilder(this_ *TypeStringBuilder, isForAnnotations bool) {
	_ = isForAnnotations
	this_.initFields_TypeStringBuilder()
	this_.isForAnnotations = isForAnnotations
}

// Java method: TypeStringBuilder.isForAnnotations()
func (this_ *TypeStringBuilder) IsForAnnotations() bool {
	return this_.isForAnnotations
}

// Java method: TypeStringBuilder.append(java.lang.String x)
func (this_ *TypeStringBuilder) Append_1_string(x string) *TypeStringBuilder {
	_ = x
	this_.builder.Append_1_string(x)
	return this_
}

// Java method: TypeStringBuilder.append(com.google.javascript.rhino.jstype.JSType x)
func (this_ *TypeStringBuilder) Append_1_JSType(x JSType) *TypeStringBuilder {
	_ = x
	x.AppendTo(this_)
	return this_
}

// Java method: TypeStringBuilder.appendAll(java.lang.Iterable<?> elements, java.lang.String separator)
func (this_ *TypeStringBuilder) AppendAll(elements j2gort.Ext_java_lang_Iterable, separator string) *TypeStringBuilder {
	_ = elements
	_ = separator
	var separate bool = false
	_ = separate
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(elements))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			e := rawElem_2
			_ = e
			if separate {
				this_.Append_1_string(separator)
			} else {
				separate = true
			}
			var jSType JSType
			_ = jSType
			var ok_3 bool
			jSType, ok_3 = j2gort.TryCast[JSType](e)
			if ok_3 {
				this_.Append_1_JSType(jSType)
			} else {
				this_.Append_1_string(j2gort.CastString(e))
			}
		}
	}
	return this_
}

// Java method: TypeStringBuilder.appendNonNull(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypeStringBuilder) AppendNonNull(type_ JSType) *TypeStringBuilder {
	_ = type_
	if ((((((this_.isForAnnotations && type_.IsObject()) && !(type_.IsUnknownType())) && !(type_.IsTemplateType())) && !(type_.IsRecordType())) && !(type_.IsFunctionType())) && !(type_.IsUnionType())) && !(type_.IsLiteralObject()) {
		this_.Append_1_string("!")
	}
	return this_.Append_1_JSType(type_)
}

// Java method: TypeStringBuilder.breakLineAndIndent()
func (this_ *TypeStringBuilder) BreakLineAndIndent() *TypeStringBuilder {
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](this_.builder.Append_1_string("\n")).Append_1_string(this_.indentation)
	return this_
}

// Java method: TypeStringBuilder.indent(java.lang.Runnable cb)
func (this_ *TypeStringBuilder) Indent(cb j2gort.Ext_java_lang_Runnable) *TypeStringBuilder {
	_ = cb
	var lastIndent string = this_.indentation
	_ = lastIndent
	this_.indentation = (j2gort.StrValString(lastIndent) + "  ")
	cb.Run()
	this_.indentation = lastIndent
	return this_
}

// Java method: TypeStringBuilder.build()
func (this_ *TypeStringBuilder) Build() string {
	return this_.builder.String()
}

// Java method: TypeStringBuilder.cloneWithConfig()
func (this_ *TypeStringBuilder) CloneWithConfig() *TypeStringBuilder {
	var val_4 TypeStringBuilder
	InitTypeStringBuilder(&val_4, this_.isForAnnotations)
	var clone *TypeStringBuilder = &val_4
	_ = clone
	{
		r_5 := clone
		r_5.indentation = this_.indentation
	}
	return clone
}
