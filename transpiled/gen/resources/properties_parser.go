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
// Port of com/google/javascript/jscomp/resources/PropertiesParser.java (generated mechanically by j2go).

package resources

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.resources.PropertiesParser
type PropertiesParser struct {
	idHash_ atomic.Int32
}

func (this_ *PropertiesParser) BindSelf_(s any) {
	_ = s
}
func (this_ *PropertiesParser) IsJavaNil_() bool { return this_ == nil }
func (*PropertiesParser) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.resources.PropertiesParser", "PropertiesParser"
}
func (this_ *PropertiesParser) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PropertiesParser) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.resources.PropertiesParser", "com.google.javascript.jscomp.resources.PropertiesParser", func(v_ any) bool {
		return j2gort.InstanceOf[*PropertiesParser](v_)
	})
}

func (this_ *PropertiesParser) initFields_PropertiesParser() {
	_ = this_
}

// Java constructor: PropertiesParser()
func NewPropertiesParser() *PropertiesParser {
	this_ := &PropertiesParser{}
	InitPropertiesParser(this_)
	return this_
}

func InitPropertiesParser(this_ *PropertiesParser) {
	this_.initFields_PropertiesParser()
}

// Java static method: PropertiesParser.parse(java.lang.String source)
func PropertiesParser_Parse(source string) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = source
	var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = builder
	var lines []string = j2gort.JStrSplit(source, "\r?\n")
	_ = lines
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(lines)); i = i + 1 {
			var line string = lines[i]
			_ = line
			if ((j2gort.CheckNotNullString(line) == "") || javalang.StringStartsWith(j2gort.CheckNotNullString(line), "#")) || javalang.StringStartsWith(j2gort.CheckNotNullString(line), "!") {
				continue
			}
			var delimeterIndex int32 = PropertiesParser_FindDelimiter(line)
			_ = delimeterIndex
			if delimeterIndex == int32(-1) {
				continue
			}
			var data string = ""
			_ = data
			var key string = javalang.StringTrim(j2gort.CheckNotNullString(j2gort.JStrSub(line, int32(0), delimeterIndex)))
			_ = key
			line = PropertiesParser_TrimLeft(j2gort.JStrSubFrom(line, (delimeterIndex + int32(1))))
			for {
				if javalang.StringEndsWith(j2gort.CheckNotNullString(line), "\\") {
					data = (j2gort.StrValString(data) + j2gort.StrValString(j2gort.JStrSub(line, int32(0), (j2gort.JStrLen(line)-int32(1)))))
					if (i + int32(1)) == int32(len(lines)) {
						break
					}
					line = PropertiesParser_TrimLeft(lines[func() int32 { old_1 := i; nw_2 := old_1 + 1; i = nw_2; return nw_2 }()])
				} else {
					data = (j2gort.StrValString(data) + j2gort.StrValString(line))
					break
				}
			}
			builder.Put_2(j2gort.BoxString(key), j2gort.BoxString(data))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow())
}

// Java static method: PropertiesParser.trimLeft(java.lang.String str)
func PropertiesParser_TrimLeft(str string) string {
	_ = str
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.JStrLen(str); i = i + 1 {
			if int32(j2gort.JStrCharAt(str, i)) != int32(uint16(32)) {
				return j2gort.JStrSubFrom(str, i)
			}
		}
	}
	return str
}

// Java static method: PropertiesParser.findDelimiter(java.lang.String line)
func PropertiesParser_FindDelimiter(line string) int32 {
	_ = line
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.JStrLen(line); i = i + 1 {
			switch j2gort.JStrCharAt(line, i) {
			case uint16(58), uint16(61):
				return i
			default:
			}
		}
	}
	return j2gort.JStrIndexOfCp(line, int32(uint16(32)))
}
