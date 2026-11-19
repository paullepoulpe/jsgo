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
// Port of com/google/javascript/jscomp/DefaultExterns.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.DefaultExterns
type DefaultExterns struct {
	idHash_ atomic.Int32
}

func (this_ *DefaultExterns) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultExterns) IsJavaNil_() bool { return this_ == nil }
func (*DefaultExterns) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultExterns", "DefaultExterns"
}
func (this_ *DefaultExterns) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultExterns) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DefaultExterns", "com.google.javascript.jscomp.DefaultExterns", func(v_ any) bool {
		return j2gort.InstanceOf[*DefaultExterns](v_)
	})
}

var static_DefaultExterns_BUILTIN_LANG_EXTERNS j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_DefaultExterns_BUILTIN_LANG_EXTERNS() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_DefaultExterns()
	return static_DefaultExterns_BUILTIN_LANG_EXTERNS
}
func StaticSet_DefaultExterns_BUILTIN_LANG_EXTERNS(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_DefaultExterns()
	static_DefaultExterns_BUILTIN_LANG_EXTERNS = v
}

var static_DefaultExterns_BROWSER_EXTERN_DEP_ORDER j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_DefaultExterns_BROWSER_EXTERN_DEP_ORDER() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_DefaultExterns()
	return static_DefaultExterns_BROWSER_EXTERN_DEP_ORDER
}
func StaticSet_DefaultExterns_BROWSER_EXTERN_DEP_ORDER(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_DefaultExterns()
	static_DefaultExterns_BROWSER_EXTERN_DEP_ORDER = v
}

var clinitDone_DefaultExterns uint32
var clinitState_DefaultExterns j2gort.ClinitState

func clinit_DefaultExterns() {
	if atomic.LoadUint32(&clinitDone_DefaultExterns) != 0 {
		return
	}
	clinitSlow_DefaultExterns()
}

func clinitSlow_DefaultExterns() {
	if !j2gort.EnterClinit(&clinitState_DefaultExterns, &clinitDone_DefaultExterns) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DefaultExterns, &clinitDone_DefaultExterns)
	static_DefaultExterns_BUILTIN_LANG_EXTERNS = j2gort.Ext_com_google_common_collect_ImmutableList_Of_5("es3.js", "es5.js", "es6.js", "es6_typed_array.js", "es6_collections.js")
	static_DefaultExterns_BROWSER_EXTERN_DEP_ORDER = j2gort.Ext_com_google_common_collect_ImmutableList_Of_13("intl.js", "w3c_event.js", "w3c_event3.js", "gecko_event.js", "ie_event.js", "webkit_event.js", "w3c_device_sensor_event.js", "w3c_dom1.js", "w3c_dom2.js", "w3c_dom3.js", "w3c_dom4.js", "gecko_dom.js", []any{"ie_dom.js", "webkit_dom.js", "w3c_css.js", "gecko_css.js", "ie_css.js", "webkit_css.js", "w3c_touch_event.js", "webstorage.js", "whatwg_console.js", "nonstandard_console.js", "w3c_rtc.js"})
}

func (this_ *DefaultExterns) initFields_DefaultExterns() {
	_ = this_
}

// Java constructor: DefaultExterns()
func NewDefaultExterns() *DefaultExterns {
	this_ := &DefaultExterns{}
	InitDefaultExterns(this_)
	return this_
}

func InitDefaultExterns(this_ *DefaultExterns) {
	this_.initFields_DefaultExterns()
}

// Java static method: DefaultExterns.prepareExterns(com.google.javascript.jscomp.CompilerOptions.Environment env, java.util.Map<java.lang.String,com.google.javascript.jscomp.SourceFile> externs)
func DefaultExterns_PrepareExterns(env CompilerOptions_Environment, externs j2gort.Ext_java_util_Map) j2gort.Ext_java_util_List {
	_ = env
	_ = externs
	var out []*SourceFile = []*SourceFile{}
	_ = out
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_DefaultExterns_BUILTIN_LANG_EXTERNS()))
		for cur_1.HasNext() {
			key := cur_1.NextStr()
			_ = key
			if !(j2gort.MapContainsKeyStr(externs, key)) {
				j2gort.CheckStateFmt(false, "Externs must contain builtin: %s", j2gort.BoxString(key))
			}
			out = append(out, j2gort.CastPtr[SourceFile](j2gort.MapRemoveStr(externs, key)))
		}
	}
	if env == StaticGet_CompilerOptions_Environment_BROWSER() {
		{
			cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_DefaultExterns_BROWSER_EXTERN_DEP_ORDER()))
			for cur_2.HasNext() {
				key_2 := cur_2.NextStr()
				_ = key_2
				if !(j2gort.MapContainsKeyStr(externs, key_2)) {
					j2gort.CheckStateFmt(false, "Externs must contain builtin for env %s: %s", j2gort.BoxNullableVal(env), j2gort.BoxString(key_2))
				}
				out = append(out, j2gort.CastPtr[SourceFile](j2gort.MapRemoveStr(externs, key_2)))
			}
		}
		out = j2gort.AppendCollection[*SourceFile](out, j2gort.AsIface[j2gort.Ext_java_util_Collection](externs.Values()))
	}
	return j2gort.WrapTypedArrayList(out)
}
