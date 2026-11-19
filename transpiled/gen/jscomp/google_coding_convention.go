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
// Port of com/google/javascript/jscomp/GoogleCodingConvention.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.GoogleCodingConvention
type GoogleCodingConvention struct {
	CodingConventions_ProxyBase
}

func (this_ *GoogleCodingConvention) BindSelf_(s any) {
	_ = s
	this_.CodingConventions_ProxyBase.BindSelf_(s)
}
func (this_ *GoogleCodingConvention) IsJavaNil_() bool { return this_ == nil }
func (*GoogleCodingConvention) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GoogleCodingConvention", "GoogleCodingConvention"
}
func (this_ *GoogleCodingConvention) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodingConventions_ProxyBase.IdentityHash_()
}
func (*GoogleCodingConvention) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GoogleCodingConvention", "com.google.javascript.jscomp.GoogleCodingConvention", func(v_ any) bool {
		return j2gort.InstanceOf[*GoogleCodingConvention](v_)
	})
}
func (*GoogleCodingConvention) Implements_jscomp_CodingConventions_Proxy() {}
func (*GoogleCodingConvention) Implements_jscomp_CodingConvention()        {}

func StaticGet_GoogleCodingConvention_SerialVersionUID() int64     { return int64(1) }
func StaticGet_GoogleCodingConvention_OPTIONAL_ARG_PREFIX() string { return "opt_" }
func StaticGet_GoogleCodingConvention_VAR_ARGS_NAME() string       { return "var_args" }

var static_GoogleCodingConvention_ENUM_KEY_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_GoogleCodingConvention_ENUM_KEY_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_GoogleCodingConvention()
	return static_GoogleCodingConvention_ENUM_KEY_PATTERN
}
func StaticSet_GoogleCodingConvention_ENUM_KEY_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_GoogleCodingConvention()
	static_GoogleCodingConvention_ENUM_KEY_PATTERN = v
}

var static_GoogleCodingConvention_PACKAGE_WITH_TEST_DIR *j2gort.Ext_java_util_regex_Pattern

func StaticGet_GoogleCodingConvention_PACKAGE_WITH_TEST_DIR() *j2gort.Ext_java_util_regex_Pattern {
	clinit_GoogleCodingConvention()
	return static_GoogleCodingConvention_PACKAGE_WITH_TEST_DIR
}
func StaticSet_GoogleCodingConvention_PACKAGE_WITH_TEST_DIR(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_GoogleCodingConvention()
	static_GoogleCodingConvention_PACKAGE_WITH_TEST_DIR = v
}

var static_GoogleCodingConvention_GENFILES_DIR *j2gort.Ext_java_util_regex_Pattern

func StaticGet_GoogleCodingConvention_GENFILES_DIR() *j2gort.Ext_java_util_regex_Pattern {
	clinit_GoogleCodingConvention()
	return static_GoogleCodingConvention_GENFILES_DIR
}
func StaticSet_GoogleCodingConvention_GENFILES_DIR(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_GoogleCodingConvention()
	static_GoogleCodingConvention_GENFILES_DIR = v
}

var clinitDone_GoogleCodingConvention uint32
var clinitState_GoogleCodingConvention j2gort.ClinitState

func clinit_GoogleCodingConvention() {
	if atomic.LoadUint32(&clinitDone_GoogleCodingConvention) != 0 {
		return
	}
	clinitSlow_GoogleCodingConvention()
}

func clinitSlow_GoogleCodingConvention() {
	if !j2gort.EnterClinit(&clinitState_GoogleCodingConvention, &clinitDone_GoogleCodingConvention) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_GoogleCodingConvention, &clinitDone_GoogleCodingConvention)
	static_GoogleCodingConvention_ENUM_KEY_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("[A-Z0-9][A-Z0-9_]*")
	static_GoogleCodingConvention_PACKAGE_WITH_TEST_DIR = j2gort.Ext_java_util_regex_Pattern_Compile_1("^(.*)/(?:test|tests|testing)/(?:[^/]+)$")
	static_GoogleCodingConvention_GENFILES_DIR = j2gort.Ext_java_util_regex_Pattern_Compile_1("-out/.*/(bin|genfiles)/(.*)$")
}

func (this_ *GoogleCodingConvention) initFields_GoogleCodingConvention() {
	_ = this_
}

// Java constructor: GoogleCodingConvention()
func NewGoogleCodingConvention() *GoogleCodingConvention {
	clinit_GoogleCodingConvention()
	this_ := &GoogleCodingConvention{}
	this_.BindSelf_(this_)
	InitGoogleCodingConvention(this_)
	return this_
}

func InitGoogleCodingConvention(this_ *GoogleCodingConvention) {
	InitGoogleCodingConvention_1(this_, NewClosureCodingConvention())
}

// Java constructor: GoogleCodingConvention(com.google.javascript.jscomp.CodingConvention convention)
func NewGoogleCodingConvention_1(convention CodingConvention) *GoogleCodingConvention {
	clinit_GoogleCodingConvention()
	this_ := &GoogleCodingConvention{}
	this_.BindSelf_(this_)
	InitGoogleCodingConvention_1(this_, convention)
	return this_
}

func InitGoogleCodingConvention_1(this_ *GoogleCodingConvention, convention CodingConvention) {
	_ = convention
	InitCodingConventions_Proxy(&this_.CodingConventions_ProxyBase, convention)
	this_.initFields_GoogleCodingConvention()
}

// Java method: GoogleCodingConvention.isConstant(java.lang.String name)
func (this_ *GoogleCodingConvention) IsConstant(name string) bool {
	_ = name
	if j2gort.JStrLen(name) <= int32(1) {
		return false
	}
	var pos int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(36)))
	_ = pos
	var start int32
	_ = start
	if pos >= int32(0) {
		start = (pos + int32(1))
	} else {
		start = int32(0)
	}
	return GoogleCodingConvention_IsConstantKey(name, start, j2gort.JStrLen(name))
}

// Java method: GoogleCodingConvention.isConstantKey(java.lang.String name)
func (this_ *GoogleCodingConvention) IsConstantKey(name string) bool {
	_ = name
	return GoogleCodingConvention_IsConstantKey(name, int32(0), j2gort.JStrLen(name))
}

// Java static method: GoogleCodingConvention.isConstantKey(java.lang.String name, int start, int end)
func GoogleCodingConvention_IsConstantKey(name string, start int32, end int32) bool {
	clinit_GoogleCodingConvention()
	_ = name
	_ = start
	_ = end
	if (start >= end) || !(j2gort.BoxChar_IsUpperCase_1_uint16(j2gort.JStrCharAt(name, start))) {
		return false
	}
	{
		var i int32 = (start + int32(1))
		_ = i
		for ; i < end; i = i + 1 {
			if int32(j2gort.BoxChar_ToUpperCase_1_uint16(j2gort.JStrCharAt(name, i))) != int32(j2gort.JStrCharAt(name, i)) {
				return false
			}
		}
	}
	return true
}

// Java method: GoogleCodingConvention.isValidEnumKey(java.lang.String key)
func (this_ *GoogleCodingConvention) IsValidEnumKey(key string) bool {
	_ = key
	return StaticGet_GoogleCodingConvention_ENUM_KEY_PATTERN().Matcher(j2gort.BoxString(key)).Matches()
}

// Java method: GoogleCodingConvention.isOptionalParameter(com.google.javascript.rhino.Node parameter)
func (this_ *GoogleCodingConvention) IsOptionalParameter(parameter pkg_rhino.Node) bool {
	_ = parameter
	return (this_.CodingConventions_ProxyBase.IsOptionalParameter(parameter) || (parameter.IsName() && javalang.StringStartsWith(j2gort.CheckNotNullString(parameter.GetString()), "opt_")))
}

// Java method: GoogleCodingConvention.isVarArgsParameter(com.google.javascript.rhino.Node parameter)
func (this_ *GoogleCodingConvention) IsVarArgsParameter(parameter pkg_rhino.Node) bool {
	_ = parameter
	return (this_.CodingConventions_ProxyBase.IsVarArgsParameter(parameter) || (parameter.IsName() && func() bool {
		sr_1 := "var_args"
		se_2 := parameter.GetString()
		return j2gort.CheckNotNullString(sr_1) == se_2
	}()))
}

// Java method: GoogleCodingConvention.isExported(java.lang.String name, boolean local)
func (this_ *GoogleCodingConvention) IsExported_2(name string, local bool) bool {
	_ = name
	_ = local
	return (this_.CodingConventions_ProxyBase.IsExported_2(name, local) || (!(local) && javalang.StringStartsWith(j2gort.CheckNotNullString(name), "_")))
}

// Java method: GoogleCodingConvention.isClassFactoryCall(com.google.javascript.rhino.Node callNode)
func (this_ *GoogleCodingConvention) IsClassFactoryCall(callNode pkg_rhino.Node) bool {
	_ = callNode
	var callTarget pkg_rhino.Node = callNode.GetFirstChild()
	_ = callTarget
	return (this_.CodingConventions_ProxyBase.IsClassFactoryCall(callNode) || (callTarget.IsName() && (j2gort.CheckNotNullString(callTarget.String()) == "Polymer")))
}

// Java method: GoogleCodingConvention.getPackageName(com.google.javascript.rhino.StaticSourceFile source)
func (this_ *GoogleCodingConvention) GetPackageName(source pkg_rhino.StaticSourceFile) string {
	_ = source
	var name string = source.GetName()
	_ = name
	var genfilesMatcher *j2gort.Ext_java_util_regex_Matcher = StaticGet_GoogleCodingConvention_GENFILES_DIR().Matcher(j2gort.BoxString(name))
	_ = genfilesMatcher
	if genfilesMatcher.Find_0() {
		name = genfilesMatcher.Group_1_int32(int32(2))
	}
	var m *j2gort.Ext_java_util_regex_Matcher = StaticGet_GoogleCodingConvention_PACKAGE_WITH_TEST_DIR().Matcher(j2gort.BoxString(name))
	_ = m
	if m.Find_0() {
		return m.Group_1_int32(int32(1))
	} else {
		var lastSlash int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(47)))
		_ = lastSlash
		if lastSlash == int32(-1) {
			return ""
		} else {
			return j2gort.JStrSub(name, int32(0), lastSlash)
		}
	}
}
