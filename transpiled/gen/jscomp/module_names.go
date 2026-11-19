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
// Port of com/google/javascript/jscomp/deps/ModuleNames.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.deps.ModuleNames
type ModuleNames struct {
	idHash_ atomic.Int32
}

func (this_ *ModuleNames) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleNames) IsJavaNil_() bool { return this_ == nil }
func (*ModuleNames) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.ModuleNames", "ModuleNames"
}
func (this_ *ModuleNames) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ModuleNames) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ModuleNames", "com.google.javascript.jscomp.deps.ModuleNames", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleNames](v_)
	})
}

func StaticGet_ModuleNames_MODULE_SLASH() string { return "/" }
func (this_ *ModuleNames) initFields_ModuleNames() {
	_ = this_
}

// Java constructor: ModuleNames()
func NewModuleNames() *ModuleNames {
	this_ := &ModuleNames{}
	InitModuleNames(this_)
	return this_
}

func InitModuleNames(this_ *ModuleNames) {
	this_.initFields_ModuleNames()
}

// Java static method: ModuleNames.fileToModuleName(java.lang.String path)
func ModuleNames_FileToModuleName(path string) string {
	_ = path
	return ModuleNames_ToModuleName(ModuleNames_EscapePath(path))
}

// Java static method: ModuleNames.fileToJsIdentifier(java.lang.String path)
func ModuleNames_FileToJsIdentifier(path string) string {
	_ = path
	return ModuleNames_ToJSIdentifier(ModuleNames_EscapePath(path))
}

// Java static method: ModuleNames.escapePath(java.lang.String input)
func ModuleNames_EscapePath(input string) string {
	_ = input
	var encodedInput string = javalang.StringReplace(javalang.StringReplace(javalang.StringReplace(javalang.StringReplace(javalang.StringReplace(j2gort.JStrReplaceChar(j2gort.JStrReplaceChar(input, uint16(58), uint16(45)), uint16(92), uint16(47)), " ", "%20"), "[", "%5B"), "]", "%5D"), "<", "%3C"), ">", "%3E")
	_ = encodedInput
	return ModuleNames_CanonicalizePath(encodedInput)
}

// Java static method: ModuleNames.toJSIdentifier(java.lang.String path)
func ModuleNames_ToJSIdentifier(path string) string {
	_ = path
	var stripped string = ModuleNames_StripJsExtension(path)
	_ = stripped
	if javalang.StringStartsWith(j2gort.CheckNotNullString(stripped), "./") {
		stripped = j2gort.JStrSubFrom(stripped, int32(2))
	}
	return javalang.StringReplace(j2gort.JStrReplaceChar(j2gort.JStrReplaceChar(j2gort.JStrReplaceChar(j2gort.JStrReplaceChar(j2gort.JStrReplaceChar(j2gort.JStrReplaceChar(j2gort.JStrReplaceChar(javalang.StringReplace(stripped, "/", "$"), uint16(92), uint16(36)), uint16(64), uint16(36)), uint16(43), uint16(36)), uint16(33), uint16(36)), uint16(45), uint16(95)), uint16(58), uint16(95)), uint16(46), uint16(95)), "%20", "_")
}

// Java static method: ModuleNames.toModuleName(java.lang.String path)
func ModuleNames_ToModuleName(path string) string {
	_ = path
	if javalang.StringStartsWith(j2gort.CheckNotNullString(path), "/") {
		path = j2gort.JStrSubFrom(path, int32(1))
	}
	return ("module$" + j2gort.StrValString(ModuleNames_ToJSIdentifier(path)))
}

// Java static method: ModuleNames.stripJsExtension(java.lang.String fileName)
func ModuleNames_StripJsExtension(fileName string) string {
	_ = fileName
	if javalang.StringEndsWith(j2gort.CheckNotNullString(fileName), ".js") {
		fileName = j2gort.JStrSub(fileName, int32(0), (j2gort.JStrLen(fileName) - j2gort.JStrLen(".js")))
		if javalang.StringEndsWith(j2gort.CheckNotNullString(fileName), ".js.i") {
			fileName = j2gort.JStrSub(fileName, int32(0), (j2gort.JStrLen(fileName) - j2gort.JStrLen(".js.i")))
		}
	}
	return fileName
}

// Java static method: ModuleNames.isAlreadyCanonical(java.lang.String path)
func ModuleNames_IsAlreadyCanonical(path string) bool {
	_ = path
	var len_ int32 = j2gort.JStrLen(path)
	_ = len_
	if len_ == int32(0) {
		return true
	}
	if int32(j2gort.JStrCharAt(path, (len_-int32(1)))) == int32(uint16(47)) {
		return false
	}
	var segStart int32 = int32(0)
	_ = segStart
	for segStart <= len_ {
		var nextSlash int32 = j2gort.JStrIndexOfCpFrom(path, int32(uint16(47)), segStart)
		_ = nextSlash
		var segEnd int32
		_ = segEnd
		if nextSlash < int32(0) {
			segEnd = len_
		} else {
			segEnd = nextSlash
		}
		var segLen int32 = (segEnd - segStart)
		_ = segLen
		if (segLen == int32(1)) && (int32(j2gort.JStrCharAt(path, segStart)) == int32(uint16(46))) {
			return false
		}
		if ((segLen == int32(2)) && (int32(j2gort.JStrCharAt(path, segStart)) == int32(uint16(46)))) && (int32(j2gort.JStrCharAt(path, (segStart+int32(1)))) == int32(uint16(46))) {
			return false
		}
		if nextSlash < int32(0) {
			break
		}
		segStart = (nextSlash + int32(1))
	}
	return true
}

// Java static method: ModuleNames.canonicalizePath(java.lang.String path)
func ModuleNames_CanonicalizePath(path string) string {
	_ = path
	if ModuleNames_IsAlreadyCanonical(path) {
		return path
	}
	var len_ int32 = j2gort.JStrLen(path)
	_ = len_
	var effectiveEnd int32 = len_
	_ = effectiveEnd
	for (effectiveEnd > int32(0)) && (int32(j2gort.JStrCharAt(path, (effectiveEnd-int32(1)))) == int32(uint16(47))) {
		effectiveEnd = effectiveEnd - 1
	}
	if effectiveEnd == int32(0) {
		return ""
	}
	var slashCount int32 = int32(0)
	_ = slashCount
	{
		var i int32 = int32(0)
		_ = i
		for ; i < effectiveEnd; i = i + 1 {
			if int32(j2gort.JStrCharAt(path, i)) == int32(uint16(47)) {
				slashCount = slashCount + 1
			}
		}
	}
	var numParts int32 = (slashCount + int32(1))
	_ = numParts
	var buffer []string = j2gort.NewStringArray(numParts)
	_ = buffer
	var position int32 = int32(0)
	_ = position
	var available int32 = int32(0)
	_ = available
	var absolutePath bool = ((numParts > int32(1)) && (int32(j2gort.JStrCharAt(path, int32(0))) == int32(uint16(47))))
	_ = absolutePath
	if absolutePath {
		available = available - 1
	}
	var segStart int32 = int32(0)
	_ = segStart
	{
		var partIdx int32 = int32(0)
		_ = partIdx
		for ; partIdx < numParts; partIdx = partIdx + 1 {
			var nextSlash int32 = j2gort.JStrIndexOfCpFrom(path, int32(uint16(47)), segStart)
			_ = nextSlash
			var segEnd int32
			_ = segEnd
			if (nextSlash < int32(0)) || (nextSlash > effectiveEnd) {
				segEnd = effectiveEnd
			} else {
				segEnd = nextSlash
			}
			var segLen int32 = (segEnd - segStart)
			_ = segLen
			if (segLen == int32(1)) && (int32(j2gort.JStrCharAt(path, segStart)) == int32(uint16(46))) {
				segStart = (segEnd + int32(1))
				continue
			}
			if ((segLen == int32(2)) && (int32(j2gort.JStrCharAt(path, segStart)) == int32(uint16(46)))) && (int32(j2gort.JStrCharAt(path, (segStart+int32(1)))) == int32(uint16(46))) {
				if available > int32(0) {
					position = position - 1
					available = available - 1
					{
						a_1 := buffer
						i_2 := position
						a_1[i_2] = j2gort.NullString
					}
				} else if !(absolutePath) {
					{
						a_3 := buffer
						i_4 := position
						a_3[i_4] = ".."
					}
					position = position + 1
				}
				segStart = (segEnd + int32(1))
				continue
			}
			{
				a_5 := buffer
				i_6 := position
				a_5[i_6] = func() string {
					if segLen == int32(0) {
						return ""
					}
					return j2gort.JStrSub(path, segStart, segEnd)
				}()
			}
			position = position + 1
			available = available + 1
			segStart = (segEnd + int32(1))
		}
	}
	if absolutePath && (position == int32(1)) {
		return "/"
	}
	if position == int32(0) {
		return ""
	}
	if position == int32(1) {
		return buffer[int32(0)]
	}
	var totalLen int32 = (position - int32(1))
	_ = totalLen
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < position; i_2 = i_2 + 1 {
			totalLen = (totalLen + j2gort.JStrLen(buffer[i_2]))
		}
	}
	var sb j2gort.LocalStringBuilder
	sb.InitCap(totalLen)

	sb.AppendStr(buffer[int32(0)])
	{
		var i_3 int32 = int32(1)
		_ = i_3
		for ; i_3 < position; i_3 = i_3 + 1 {
			sb.AppendChar(uint16(47)).AppendStr(buffer[i_3])
		}
	}
	return sb.String()
}
