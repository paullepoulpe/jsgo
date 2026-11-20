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
// Port of com/google/javascript/jscomp/SourceFile.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.SourceFile
type SourceFile struct {
	fileName                                string                `gson:"fileName"`
	loader                                  SourceFile_CodeLoader `gson:"loader"`
	lineOffsets                             []int32               `gson:"lineOffsets"`
	code                                    string                `gson:"code"`
	_sfHandle                               atomic.Uint32
	numLines                                int32 `gson:"numLines"`
	numBytes                                int32 `gson:"numBytes"`
	idHash_                                 atomic.Int32
	kind                                    pkg_rhino.StaticSourceFile_SourceKind `gson:"kind"`
	isStubSourceFileForAlreadyProvidedInput bool                                  `gson:"isStubSourceFileForAlreadyProvidedInput"`
	isClosureUnawareCode                    bool                                  `gson:"isClosureUnawareCode"`
}

func (this_ *SourceFile) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceFile) IsJavaNil_() bool { return this_ == nil }
func (*SourceFile) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceFile", "SourceFile"
}
func (this_ *SourceFile) SourceFileHandlePtr_() *atomic.Uint32 { return &this_._sfHandle }
func (this_ *SourceFile) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceFile) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceFile", "com.google.javascript.jscomp.SourceFile", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceFile](v_)
	})
}
func (*SourceFile) Implements_rhino_StaticSourceFile() {}

func StaticGet_SourceFile_UTF8_BOM() string { return "\ufeff" }

var static_SourceFile_BANG_SLASH string

func StaticGet_SourceFile_BANG_SLASH() string {
	clinit_SourceFile()
	return static_SourceFile_BANG_SLASH
}
func StaticSet_SourceFile_BANG_SLASH(v string) {
	clinit_SourceFile()
	static_SourceFile_BANG_SLASH = v
}
func StaticGet_SourceFile_SOURCE_EXCERPT_REGION_LENGTH() int32 { return int32(5) }

var clinitDone_SourceFile uint32
var clinitState_SourceFile j2gort.ClinitState

func clinit_SourceFile() {
	if atomic.LoadUint32(&clinitDone_SourceFile) != 0 {
		return
	}
	clinitSlow_SourceFile()
}

func clinitSlow_SourceFile() {
	if !j2gort.EnterClinit(&clinitState_SourceFile, &clinitDone_SourceFile) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SourceFile, &clinitDone_SourceFile)
	static_SourceFile_BANG_SLASH = ("!" + j2gort.StrValString(Platform_GetFileSeperator()))
}

func (this_ *SourceFile) initFields_SourceFile() {
	_ = this_
	this_.fileName = j2gort.NullString

	this_.code = j2gort.NullString
	this_.numLines = int32(-1)
	this_.numBytes = int32(-1)

}

// Java constructor: SourceFile(com.google.javascript.jscomp.SourceFile.CodeLoader loader, java.lang.String fileName, com.google.javascript.rhino.StaticSourceFile.SourceKind kind)
func NewSourceFile(loader SourceFile_CodeLoader, fileName string, kind pkg_rhino.StaticSourceFile_SourceKind) *SourceFile {
	clinit_SourceFile()
	this_ := &SourceFile{}
	InitSourceFile(this_, loader, fileName, kind)
	return this_
}

func InitSourceFile(this_ *SourceFile, loader SourceFile_CodeLoader, fileName string, kind pkg_rhino.StaticSourceFile_SourceKind) {
	_ = loader
	_ = fileName
	_ = kind
	this_.initFields_SourceFile()
	if j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(fileName) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "a source must have a name")))
	}
	{
		var iifeCond_1 bool
		{

			sr_1 := "/"
			se_2 := Platform_GetFileSeperator()
			iifeCond_1 = j2gort.CheckNotNullString(sr_1) == se_2
		}
		if !(iifeCond_1) {
			fileName = javalang.StringReplace(fileName, j2gort.CheckNotNullString(Platform_GetFileSeperator()), "/")
		}
	}
	this_.loader = loader
	this_.fileName = fileName
	this_.kind = kind
}

// Java method: SourceFile.isStubSourceFileForAlreadyProvidedInput()
func (this_ *SourceFile) IsStubSourceFileForAlreadyProvidedInput() bool {
	return this_.isStubSourceFileForAlreadyProvidedInput
}

// Java method: SourceFile.setIsStubSourceFileForAlreadyProvidedInput(boolean isStubSourceFileForAlreadyProvidedInput)
func (this_ *SourceFile) priv_SourceFile_setIsStubSourceFileForAlreadyProvidedInput_5124(isStubSourceFileForAlreadyProvidedInput bool) {
	_ = isStubSourceFileForAlreadyProvidedInput
	this_.isStubSourceFileForAlreadyProvidedInput = isStubSourceFileForAlreadyProvidedInput
}

// Java method: SourceFile.getLineOffset(int lineno)
func (this_ *SourceFile) GetLineOffset(lineno int32) int32 {
	_ = lineno
	this_.priv_SourceFile_findLineOffsets_5128()
	if (lineno < int32(1)) || (lineno > int32(len(this_.lineOffsets))) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString(("Expected line number between 1 and " + j2gort.StrValInt32(int32(len(this_.lineOffsets))) + "\nActual: " + j2gort.StrValInt32(lineno))))))
	}
	return this_.lineOffsets[(lineno - int32(1))]
}

// Java method: SourceFile.getNumLines()
func (this_ *SourceFile) GetNumLines() int32 {
	if this_.numLines < int32(0) {
		var didRet_3 bool
		var retVal_4 int32
		func() {
			defer func() {
				if rec_5 := recover(); rec_5 != nil {
					jex_6 := javalang.RecoverException(rec_5)
					_ = jex_6
					if j2gort.InstanceOf[javaio.IIOException](jex_6) {
						e := j2gort.AsIface[javaio.IIOException](jex_6)
						_ = e
						retVal_4 = int32(0)
						didRet_3 = true
						return
					} else {
						panic(rec_5)
					}
				}
			}()
			var unused string = this_.GetCode()
			_ = unused
		}()
		if didRet_3 {
			return retVal_4
		}
	}
	return this_.numLines
}

// Java method: SourceFile.getNumBytes()
func (this_ *SourceFile) GetNumBytes() int32 {
	if this_.numBytes < int32(0) {
		var didRet_7 bool
		var retVal_8 int32
		func() {
			defer func() {
				if rec_9 := recover(); rec_9 != nil {
					jex_10 := javalang.RecoverException(rec_9)
					_ = jex_10
					if j2gort.InstanceOf[javaio.IIOException](jex_10) {
						e := j2gort.AsIface[javaio.IIOException](jex_10)
						_ = e
						retVal_8 = int32(0)
						didRet_7 = true
						return
					} else {
						panic(rec_9)
					}
				}
			}()
			var unused string = this_.GetCode()
			_ = unused
		}()
		if didRet_7 {
			return retVal_8
		}
	}
	return this_.numBytes
}

// Java static method: SourceFile.isLineTerminator(char ch)
func SourceFile_IsLineTerminator(ch uint16) bool {
	clinit_SourceFile()
	_ = ch
	switch ch {
	case uint16(10), uint16(13), uint16(8232), uint16(8233):
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: SourceFile.findLineOffsets()
func (this_ *SourceFile) priv_SourceFile_findLineOffsets_5128() {
	if this_.lineOffsets != nil {
		return
	}
	var localCode string = this_.code
	_ = localCode
	if localCode == j2gort.NullString {
		var didRet_11 bool
		func() {
			defer func() {
				if rec_12 := recover(); rec_12 != nil {
					jex_13 := javalang.RecoverException(rec_12)
					_ = jex_13
					if j2gort.InstanceOf[javaio.IIOException](jex_13) {
						e := j2gort.AsIface[javaio.IIOException](jex_13)
						_ = e
						localCode = ""
						this_.lineOffsets = make([]int32, int32(1))
						didRet_11 = true
						return
					} else {
						panic(rec_12)
					}
				}
			}()
			localCode = this_.GetCode()
		}()
		if didRet_11 {
			return
		}
	}
	var offsets []int32 = make([]int32, this_.numLines)
	_ = offsets
	var index int32 = int32(1)
	_ = index
	var len_ int32 = j2gort.JStrLen(localCode)
	_ = len_
	{
		var i int32 = int32(0)
		_ = i
		for ; i < len_; i = i + 1 {
			var ch uint16 = j2gort.JStrCharAt(localCode, i)
			_ = ch
			if SourceFile_IsLineTerminator(ch) {
				if ((int32(ch) == int32(uint16(13))) && ((i + int32(1)) < len_)) && (int32(j2gort.JStrCharAt(localCode, (i+int32(1)))) == int32(uint16(10))) {
					i = i + 1
				}
				{
					a_14 := offsets
					var i_15 int32
					{

						old_16 := index
						nw_17 := old_16 + 1
						index = nw_17
						i_15 = old_16
					}
					a_14[i_15] = (i + int32(1))
				}
			}
		}
	}
	if !(index == int32(len(offsets))) {
		j2gort.CheckState(false)
	}
	this_.lineOffsets = offsets
}

// Java method: SourceFile.getCode()
func (this_ *SourceFile) GetCode() string {
	if this_.code == j2gort.NullString {
		{
			unlock_18 := j2gort.LockMonitor(this_)
			defer unlock_18()
			if this_.code == j2gort.NullString {
				this_.priv_SourceFile_setCodeAndDoBookkeeping_5132(this_.loader.LoadUncachedCode())
			}
		}
	}
	return this_.code
}

// Java method: SourceFile.setCodeDeprecated(java.lang.String code)
func (this_ *SourceFile) SetCodeDeprecated(code string) {
	_ = code
	this_.priv_SourceFile_setCodeAndDoBookkeeping_5132(code)
}

// Java method: SourceFile.getCodeReader()
func (this_ *SourceFile) GetCodeReader() j2gort.Ext_java_io_Reader {
	if this_.code == j2gort.NullString {
		{
			unlock_19 := j2gort.LockMonitor(this_)
			defer unlock_19()
			if this_.code == j2gort.NullString {
				var uncachedReader j2gort.Ext_java_io_Reader = this_.loader.OpenUncachedReader()
				_ = uncachedReader
				if uncachedReader != nil {
					return uncachedReader
				}
			}
		}
	}
	return j2gort.NewExt_java_io_StringReader(this_.GetCode())
}

// Java method: SourceFile.setCodeAndDoBookkeeping(java.lang.@org.jspecify.annotations.Nullable String sourceCode)
func (this_ *SourceFile) priv_SourceFile_setCodeAndDoBookkeeping_5132(sourceCode string) {
	_ = sourceCode
	this_.code = j2gort.NullString
	this_.lineOffsets = nil
	if sourceCode != j2gort.NullString {
		if javalang.StringStartsWith(j2gort.CheckNotNullString(sourceCode), "\ufeff") {
			sourceCode = j2gort.JStrSubFrom(sourceCode, j2gort.JStrLen("\ufeff"))
		}
		this_.code = sourceCode
		var len_ int32 = j2gort.JStrLen(sourceCode)
		_ = len_
		this_.numBytes = len_
		var numLines int32 = int32(1)
		_ = numLines
		{
			var i int32 = int32(0)
			_ = i
			for ; i < len_; i = i + 1 {
				var ch uint16 = j2gort.JStrCharAt(sourceCode, i)
				_ = ch
				if SourceFile_IsLineTerminator(ch) {
					if ((int32(ch) == int32(uint16(13))) && ((i + int32(1)) < len_)) && (int32(j2gort.JStrCharAt(sourceCode, (i+int32(1)))) == int32(uint16(10))) {
						i = i + 1
					}
					numLines = numLines + 1
				}
			}
		}
		this_.numLines = numLines
	}
}

// Java method: SourceFile.getOriginalPath()
func (this_ *SourceFile) GetOriginalPath() string {
	return this_.GetName()
}

// Java method: SourceFile.clearCachedSource()
func (this_ *SourceFile) ClearCachedSource() {
	this_.priv_SourceFile_setCodeAndDoBookkeeping_5132(j2gort.NullString)
}

// Java method: SourceFile.hasSourceInMemory()
func (this_ *SourceFile) HasSourceInMemory() bool {
	return (this_.code != j2gort.NullString)
}

// Java method: SourceFile.getName()
func (this_ *SourceFile) GetName() string {
	return this_.fileName
}

// Java method: SourceFile.getKind()
func (this_ *SourceFile) GetKind() pkg_rhino.StaticSourceFile_SourceKind {
	return this_.kind
}

// Java method: SourceFile.setKind(com.google.javascript.rhino.StaticSourceFile.SourceKind kind)
func (this_ *SourceFile) SetKind(kind pkg_rhino.StaticSourceFile_SourceKind) {
	_ = kind
	this_.kind = kind
}

// Java method: SourceFile.isClosureUnawareCode()
func (this_ *SourceFile) IsClosureUnawareCode() bool {
	return this_.isClosureUnawareCode
}

// Java method: SourceFile.markAsClosureUnawareCode()
func (this_ *SourceFile) MarkAsClosureUnawareCode() {
	this_.isClosureUnawareCode = true
}

// Java method: SourceFile.getLineOfOffset(int offset)
func (this_ *SourceFile) GetLineOfOffset(offset int32) int32 {
	_ = offset
	if !(offset >= int32(0)) {
		j2gort.CheckArgumentFmt(false, "offset must be non-negative: %s", j2gort.ValueOfInt(offset))
	}
	this_.priv_SourceFile_findLineOffsets_5128()
	var search int32 = j2gort.Ext_java_util_Arrays_BinarySearch_2_slice_int32_int32(this_.lineOffsets, offset)
	_ = search
	if search >= int32(0) {
		return (search + int32(1))
	} else {
		var insertionPoint int32 = (int32(-1) * (search + int32(1)))
		_ = insertionPoint
		return (min((insertionPoint-int32(1)), (int32(len(this_.lineOffsets))-int32(1))) + int32(1))
	}
}

// Java method: SourceFile.getColumnOfOffset(int offset)
func (this_ *SourceFile) GetColumnOfOffset(offset int32) int32 {
	_ = offset
	var line int32 = this_.GetLineOfOffset(offset)
	_ = line
	return (offset - this_.lineOffsets[(line-int32(1))])
}

// Java method: SourceFile.getLine(int lineNumber)
func (this_ *SourceFile) GetLine(lineNumber int32) string {
	_ = lineNumber
	var js string = j2gort.NullString
	_ = js
	var didRet_20 bool
	var retVal_21 string
	func() {
		defer func() {
			if rec_22 := recover(); rec_22 != nil {
				jex_23 := javalang.RecoverException(rec_22)
				_ = jex_23
				if j2gort.InstanceOf[javaio.IIOException](jex_23) {
					e := j2gort.AsIface[javaio.IIOException](jex_23)
					_ = e
					retVal_21 = j2gort.NullString
					didRet_20 = true
					return
				} else {
					panic(rec_22)
				}
			}
		}()
		js = this_.GetCode()
	}()
	if didRet_20 {
		return retVal_21
	}
	this_.priv_SourceFile_findLineOffsets_5128()
	if lineNumber > int32(len(this_.lineOffsets)) {
		return j2gort.NullString
	}
	if lineNumber < int32(1) {
		lineNumber = int32(1)
	}
	var pos int32 = this_.lineOffsets[(lineNumber - int32(1))]
	_ = pos
	if lineNumber < int32(len(this_.lineOffsets)) {
		var nextPos int32 = this_.lineOffsets[lineNumber]
		_ = nextPos
		var endPos int32 = (nextPos - int32(1))
		_ = endPos
		if ((int32(j2gort.JStrCharAt(js, endPos)) == int32(uint16(10))) && (endPos > pos)) && (int32(j2gort.JStrCharAt(js, (endPos-int32(1)))) == int32(uint16(13))) {
			endPos = endPos - 1
		}
		return j2gort.JStrSub(js, pos, endPos)
	} else {
		if pos >= j2gort.JStrLen(js) {
			return j2gort.NullString
		} else {
			return j2gort.JStrSubFrom(js, pos)
		}
	}
}

// Java method: SourceFile.getLines(int lineNumber, int length)
func (this_ *SourceFile) GetLines(lineNumber int32, length int32) Region {
	_ = lineNumber
	_ = length
	var js string = j2gort.NullString
	_ = js
	var didRet_24 bool
	var retVal_25 Region
	func() {
		defer func() {
			if rec_26 := recover(); rec_26 != nil {
				jex_27 := javalang.RecoverException(rec_26)
				_ = jex_27
				if j2gort.InstanceOf[javaio.IIOException](jex_27) {
					e := j2gort.AsIface[javaio.IIOException](jex_27)
					_ = e
					retVal_25 = nil
					didRet_24 = true
					return
				} else {
					panic(rec_26)
				}
			}
		}()
		js = this_.GetCode()
	}()
	if didRet_24 {
		return retVal_25
	}
	this_.priv_SourceFile_findLineOffsets_5128()
	if lineNumber > int32(len(this_.lineOffsets)) {
		return nil
	}
	if lineNumber < int32(1) {
		lineNumber = int32(1)
	}
	if length <= int32(0) {
		length = int32(1)
	}
	var pos int32 = this_.lineOffsets[(lineNumber - int32(1))]
	_ = pos
	if pos == j2gort.JStrLen(js) {
		return NewSimpleRegion(lineNumber, lineNumber, "")
	}
	var endChar int32 = pos
	_ = endChar
	var endLine int32 = lineNumber
	_ = endLine
	for ; (endChar < (pos + length)) && (endLine <= int32(len(this_.lineOffsets))); endLine = endLine + 1 {
		if endLine < int32(len(this_.lineOffsets)) {
			endChar = this_.lineOffsets[endLine]
		} else {
			endChar = j2gort.JStrLen(js)
		}
	}
	if endChar > pos {
		if int32(j2gort.JStrCharAt(js, (endChar-int32(1)))) == int32(uint16(10)) {
			if ((endChar - int32(1)) > pos) && (int32(j2gort.JStrCharAt(js, (endChar-int32(2)))) == int32(uint16(13))) {
				return NewSimpleRegion(lineNumber, endLine, j2gort.JStrSub(js, pos, (endChar-int32(2))))
			}
			return NewSimpleRegion(lineNumber, endLine, j2gort.JStrSub(js, pos, (endChar-int32(1))))
		} else if SourceFile_IsLineTerminator(j2gort.JStrCharAt(js, (endChar - int32(1)))) {
			return NewSimpleRegion(lineNumber, endLine, j2gort.JStrSub(js, pos, (endChar-int32(1))))
		}
	}
	return NewSimpleRegion(lineNumber, endLine, j2gort.JStrSub(js, pos, endChar))
}

// Java method: SourceFile.getRegion(int lineNumber)
func (this_ *SourceFile) GetRegion(lineNumber int32) Region {
	_ = lineNumber
	var js string = ""
	_ = js
	var didRet_28 bool
	var retVal_29 Region
	func() {
		defer func() {
			if rec_30 := recover(); rec_30 != nil {
				jex_31 := javalang.RecoverException(rec_30)
				_ = jex_31
				if j2gort.InstanceOf[javaio.IIOException](jex_31) {
					e := j2gort.AsIface[javaio.IIOException](jex_31)
					_ = e
					retVal_29 = nil
					didRet_28 = true
					return
				} else {
					panic(rec_30)
				}
			}
		}()
		js = this_.GetCode()
	}()
	if didRet_28 {
		return retVal_29
	}
	var pos int32 = int32(0)
	_ = pos
	var startLine int32 = max(int32(1), ((lineNumber - ((int32(5) + int32(1)) / int32(2))) + int32(1)))
	_ = startLine
	{
		var n int32 = int32(1)
		_ = n
		for ; n < startLine; n = n + 1 {
			var nextpos int32 = j2gort.JStrIndexOfCpFrom(js, int32(uint16(10)), pos)
			_ = nextpos
			if nextpos == int32(-1) {
				break
			}
			pos = (nextpos + int32(1))
		}
	}
	var end int32 = pos
	_ = end
	var endLine int32 = startLine
	_ = endLine
	{
		var n_2 int32 = int32(0)
		_ = n_2
		for ; n_2 < int32(5); func() { n_2 = n_2 + 1; endLine = endLine + 1 }() {
			end = j2gort.JStrIndexOfCpFrom(js, int32(uint16(10)), end)
			if end == int32(-1) {
				break
			}
			end = end + 1
		}
	}
	if lineNumber >= endLine {
		return nil
	}
	if end == int32(-1) {
		var last int32 = (j2gort.JStrLen(js) - int32(1))
		_ = last
		if int32(j2gort.JStrCharAt(js, last)) == int32(uint16(10)) {
			return NewSimpleRegion(startLine, endLine, j2gort.JStrSub(js, pos, last))
		} else {
			return NewSimpleRegion(startLine, endLine, j2gort.JStrSubFrom(js, pos))
		}
	} else {
		return NewSimpleRegion(startLine, endLine, j2gort.JStrSub(js, pos, end))
	}
}

// Java method: SourceFile.toString()
func (this_ *SourceFile) String() string {
	return this_.fileName
}

// Java static method: SourceFile.fromZipFile(java.lang.String zipName, java.nio.charset.Charset inputCharset)
func SourceFile_FromZipFile(zipName string, inputCharset j2gort.Ext_java_nio_charset_Charset) j2gort.Ext_java_util_List {
	clinit_SourceFile()
	_ = zipName
	_ = inputCharset
	var didRet_32 bool
	var retVal_33 j2gort.Ext_java_util_List
	func() {
		var input j2gort.Ext_java_io_InputStream = j2gort.NewExt_java_io_FileInputStream_1_string(zipName)
		_ = input
		defer j2gort.CloseResource(input)
		retVal_33 = SourceFile_FromZipInput(zipName, input, inputCharset)
		didRet_32 = true
		return
	}()
	if didRet_32 {
		return retVal_33
	}
	panic("unreachable")
}

// Java static method: SourceFile.fromZipInput(java.lang.String zipName, java.io.InputStream input, java.nio.charset.Charset inputCharset)
func SourceFile_FromZipInput(zipName string, input j2gort.Ext_java_io_InputStream, inputCharset j2gort.Ext_java_nio_charset_Charset) j2gort.Ext_java_util_List {
	clinit_SourceFile()
	_ = zipName
	_ = input
	_ = inputCharset
	var absoluteZipPath string = j2gort.NewExt_java_io_File_1_string(zipName).GetAbsolutePath()
	_ = absoluteZipPath
	var sourceFiles []*SourceFile = []*SourceFile{}
	_ = sourceFiles
	func() {
		var in *j2gort.Ext_java_util_zip_ZipInputStream = j2gort.NewExt_java_util_zip_ZipInputStream_2(input, inputCharset)
		_ = in
		defer j2gort.CloseResource(in)
		var zipEntry *j2gort.Ext_java_util_zip_ZipEntry = nil
		_ = zipEntry
		for (func() *j2gort.Ext_java_util_zip_ZipEntry {
			var asgn_34 *j2gort.Ext_java_util_zip_ZipEntry = in.GetNextEntry()
			zipEntry = asgn_34
			return asgn_34
		}()) != nil {
			var entryName string = zipEntry.GetName()
			_ = entryName
			if !(javalang.StringEndsWith(j2gort.CheckNotNullString(entryName), ".js")) {
				continue
			}
			sourceFiles = append(sourceFiles, SourceFile_S_Builder().WithCharset(inputCharset).WithOriginalPath((j2gort.StrValString(zipName)+j2gort.StrValString(StaticGet_SourceFile_BANG_SLASH())+j2gort.StrValString(entryName))).WithZipEntryPath(absoluteZipPath, entryName).Build())
		}
	}()
	return j2gort.WrapTypedArrayList(sourceFiles)
}

// Java static method: SourceFile.fromFile(java.lang.String fileName, java.nio.charset.Charset charset)
func SourceFile_FromFile_2(fileName string, charset j2gort.Ext_java_nio_charset_Charset) *SourceFile {
	clinit_SourceFile()
	_ = fileName
	_ = charset
	return SourceFile_S_Builder().WithPath_1_string(fileName).WithCharset(charset).Build()
}

// Java static method: SourceFile.fromFile(java.lang.String fileName)
func SourceFile_FromFile_1(fileName string) *SourceFile {
	clinit_SourceFile()
	_ = fileName
	return SourceFile_S_Builder().WithPath_1_string(fileName).Build()
}

// Java static method: SourceFile.fromPath(java.nio.file.Path path, java.nio.charset.Charset charset)
func SourceFile_FromPath(path j2gort.Ext_java_nio_file_Path, charset j2gort.Ext_java_nio_charset_Charset) *SourceFile {
	clinit_SourceFile()
	_ = path
	_ = charset
	return SourceFile_S_Builder().WithPath_1_Ext_java_nio_file_Path(path).WithCharset(charset).Build()
}

// Java static method: SourceFile.fromCode(java.lang.String fileName, java.lang.String code, com.google.javascript.rhino.StaticSourceFile.SourceKind kind)
func SourceFile_FromCode_3(fileName string, code string, kind pkg_rhino.StaticSourceFile_SourceKind) *SourceFile {
	clinit_SourceFile()
	_ = fileName
	_ = code
	_ = kind
	return SourceFile_S_Builder().WithPath_1_string(fileName).WithKind(kind).WithContent_1_string(code).Build()
}

// Java static method: SourceFile.fromCode(java.lang.String fileName, java.lang.String code)
func SourceFile_FromCode_2(fileName string, code string) *SourceFile {
	clinit_SourceFile()
	_ = fileName
	_ = code
	return SourceFile_S_Builder().WithPath_1_string(fileName).WithContent_1_string(code).Build()
}

// Java static method: SourceFile.stubSourceFile(java.lang.String fileName, com.google.javascript.rhino.StaticSourceFile.SourceKind kind)
func SourceFile_StubSourceFile(fileName string, kind pkg_rhino.StaticSourceFile_SourceKind) *SourceFile {
	clinit_SourceFile()
	_ = fileName
	_ = kind
	return SourceFile_S_Builder().WithPath_1_string(fileName).WithKind(kind).SetIsStubSourceFileForAlreadyProvidedInput().Build()
}

// Java static method: SourceFile.stripConfigSegment(java.lang.String path)
func SourceFile_StripConfigSegment(path string) string {
	clinit_SourceFile()
	_ = path
	if javalang.StringStartsWith(j2gort.CheckNotNullString(path), "bin/") {
		return j2gort.JStrSubFrom(path, int32(4))
	}
	if javalang.StringStartsWith(j2gort.CheckNotNullString(path), "genfiles/") {
		return j2gort.JStrSubFrom(path, int32(9))
	}
	if javalang.StringStartsWith(j2gort.CheckNotNullString(path), "testlogs/") {
		return j2gort.JStrSubFrom(path, int32(9))
	}
	var blazeOutIdx int32 = j2gort.JStrIndexOf(path, "blaze-out/")
	_ = blazeOutIdx
	if blazeOutIdx == int32(-1) {
		blazeOutIdx = j2gort.JStrIndexOf(path, "bazel-out/")
	}
	if blazeOutIdx != int32(-1) {
		var nextSlash int32 = j2gort.JStrIndexOfCpFrom(path, int32(uint16(47)), (blazeOutIdx + int32(10)))
		_ = nextSlash
		if nextSlash != int32(-1) {
			if j2gort.JStrStartsWithOffset(path, "bin/", (nextSlash + int32(1))) {
				return j2gort.JStrSubFrom(path, (nextSlash + int32(5)))
			}
			if j2gort.JStrStartsWithOffset(path, "genfiles/", (nextSlash + int32(1))) {
				return j2gort.JStrSubFrom(path, (nextSlash + int32(10)))
			}
			if j2gort.JStrStartsWithOffset(path, "testlogs/", (nextSlash + int32(1))) {
				return j2gort.JStrSubFrom(path, (nextSlash + int32(10)))
			}
		}
	}
	return path
}

// Java static method: SourceFile.canonicalizePathForMatching(java.lang.String path)
func SourceFile_CanonicalizePathForMatching(path string) string {
	clinit_SourceFile()
	_ = path
	return SourceFile_StripConfigSegment(path)
}

// Java method: SourceFile.restoreCachedStateFrom(com.google.javascript.jscomp.serialization.SourceFileProto protoSourceFile)
func (this_ *SourceFile) RestoreCachedStateFrom(protoSourceFile *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto) {
	_ = protoSourceFile
	{
		var iifeCond_1 bool
		{

			sr_35 := protoSourceFile.GetFilename()
			se_36 := this_.GetName()
			iifeCond_1 = j2gort.CheckNotNullString(sr_35) == se_36
		}
		if !(iifeCond_1 || func() bool {
			sr_37 := SourceFile_CanonicalizePathForMatching(protoSourceFile.GetFilename())
			se_38 := SourceFile_CanonicalizePathForMatching(this_.GetName())
			return j2gort.CheckNotNullString(sr_37) == se_38
		}()) {
			j2gort.CheckStateFmt(false, "Cannot restore state for %s from %s", j2gort.BoxString(this_.GetName()), j2gort.BoxString(protoSourceFile.GetFilename()))
		}
	}
	if protoSourceFile.GetSourceKind() == j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_SourceKind_F_EXTERN {
		if !(this_.GetKind() == pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN()) {
			j2gort.CheckStateFmt(false, "TypedAST compilations must pass all extern files as externs, not js, but found %s", j2gort.BoxString(this_.GetName()))
		}
	}
	var protoNumLines int32 = (protoSourceFile.GetNumLinesPlusOne() - int32(1))
	_ = protoNumLines
	var protoNumBytes int32 = (protoSourceFile.GetNumBytesPlusOne() - int32(1))
	_ = protoNumBytes
	if !(protoNumLines >= int32(-1)) {
		j2gort.CheckStateFmt(false, "protoNumLines must be >= -1: %s", j2gort.ValueOfInt(protoNumLines))
	}
	if !(protoNumBytes >= int32(-1)) {
		j2gort.CheckStateFmt(false, "protoNumBytes must be >= -1: %s", j2gort.ValueOfInt(protoNumBytes))
	}
	if protoNumLines != int32(-1) {
		this_.numLines = protoNumLines
	} else {
		this_.numLines = this_.numLines
	}
	if protoNumBytes != int32(-1) {
		this_.numBytes = protoNumBytes
	} else {
		this_.numBytes = this_.numBytes
	}
	this_.isClosureUnawareCode = protoSourceFile.GetIsClosureUnawareCode()
}

// Java static method: SourceFile.fromProto(com.google.javascript.jscomp.serialization.SourceFileProto protoSourceFile)
func SourceFile_FromProto_1(protoSourceFile *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto) *SourceFile {
	clinit_SourceFile()
	_ = protoSourceFile
	var sourceKind pkg_rhino.StaticSourceFile_SourceKind = SourceFile_GetSourceKindFromProto(protoSourceFile)
	_ = sourceKind
	var sourceFile *SourceFile = SourceFile_FromProto_2(protoSourceFile, sourceKind)
	_ = sourceFile
	var protoNumLines int32 = (protoSourceFile.GetNumLinesPlusOne() - int32(1))
	_ = protoNumLines
	var protoNumBytes int32 = (protoSourceFile.GetNumBytesPlusOne() - int32(1))
	_ = protoNumBytes
	if !(protoNumLines >= int32(-1)) {
		j2gort.CheckStateFmt(false, "protoNumLines must be >= -1: %s", j2gort.ValueOfInt(protoNumLines))
	}
	if !(protoNumBytes >= int32(-1)) {
		j2gort.CheckStateFmt(false, "protoNumBytes must be >= -1: %s", j2gort.ValueOfInt(protoNumBytes))
	}
	{
		r_39 := sourceFile
		r_39.numLines = protoNumLines
	}
	{
		r_40 := sourceFile
		r_40.numBytes = protoNumBytes
	}
	if protoSourceFile.GetIsClosureUnawareCode() {
		sourceFile.MarkAsClosureUnawareCode()
	}
	return sourceFile
}

// Java static method: SourceFile.fromProto(com.google.javascript.jscomp.serialization.SourceFileProto protoSourceFile, com.google.javascript.rhino.StaticSourceFile.SourceKind sourceKind)
func SourceFile_FromProto_2(protoSourceFile *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto, sourceKind pkg_rhino.StaticSourceFile_SourceKind) *SourceFile {
	clinit_SourceFile()
	_ = protoSourceFile
	_ = sourceKind
	switch j2gort.CheckNotNullPtr(protoSourceFile.GetLoaderCase()).Ordinal() {
	case int32(0) /* LoaderCase.PRELOADED_CONTENTS */ :
		return SourceFile_FromCode_3(protoSourceFile.GetFilename(), protoSourceFile.GetPreloadedContents(), sourceKind)
	case int32(1) /* LoaderCase.FILE_ON_DISK */ :
		var pathOnDisk string
		_ = pathOnDisk
		if j2gort.CheckNotNullString(protoSourceFile.GetFileOnDisk().GetActualPath()) == "" {
			pathOnDisk = protoSourceFile.GetFilename()
		} else {
			pathOnDisk = protoSourceFile.GetFileOnDisk().GetActualPath()
		}
		return SourceFile_S_Builder().WithCharset(SourceFile_ToCharset(protoSourceFile.GetFileOnDisk().GetCharset())).WithOriginalPath(protoSourceFile.GetFilename()).WithKind(sourceKind).WithPath_1_string(pathOnDisk).Build()
	case int32(2) /* LoaderCase.ZIP_ENTRY */ :
		var zipEntry *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk = protoSourceFile.GetZipEntry()
		_ = zipEntry
		return SourceFile_S_Builder().WithKind(sourceKind).WithOriginalPath(protoSourceFile.GetFilename()).WithCharset(SourceFile_ToCharset(zipEntry.GetCharset())).WithZipEntryPath(zipEntry.GetZipPath(), zipEntry.GetEntryName()).Build()
	case int32(3) /* LoaderCase.STUB_FILE */ :
		return SourceFile_S_Builder().WithKind(sourceKind).WithOriginalPath(protoSourceFile.GetFilename()).SetIsStubSourceFileForAlreadyProvidedInput().Build()
	case int32(4) /* LoaderCase.LOADER_NOT_SET */ :
	}
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java static method: SourceFile.getSourceKindFromProto(com.google.javascript.jscomp.serialization.SourceFileProto protoSourceFile)
func SourceFile_GetSourceKindFromProto(protoSourceFile *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto) pkg_rhino.StaticSourceFile_SourceKind {
	clinit_SourceFile()
	_ = protoSourceFile
	switch j2gort.CheckNotNullPtr(protoSourceFile.GetSourceKind()).Ordinal() {
	case int32(1) /* SourceKind.EXTERN */ :
		return pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN()
	case int32(2) /* SourceKind.CODE */ :
		return pkg_rhino.StaticGet_StaticSourceFile_SourceKind_STRONG()
	case int32(0) /* SourceKind.NOT_SPECIFIED */, int32(3) /* SourceKind.UNRECOGNIZED */ :
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.BoxString(protoSourceFile.GetSourceKind().String()))))
	}
	panic("unreachable switch expression")
}

// Java static method: SourceFile.toCharset(java.lang.String protoCharset)
func SourceFile_ToCharset(protoCharset string) j2gort.Ext_java_nio_charset_Charset {
	clinit_SourceFile()
	_ = protoCharset
	if j2gort.CheckNotNullString(protoCharset) == "" {
		return j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8
	}
	return j2gort.Ext_java_nio_charset_Charset_ForName_1(protoCharset)
}

// Java static method: SourceFile.builder()
func SourceFile_S_Builder() *SourceFile_Builder {
	clinit_SourceFile()
	return NewSourceFile_Builder()
}

// Java method: SourceFile.getProto()
func (this_ *SourceFile) GetProto() *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto {
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto](this_.loader.ToProtoLocationBuilder(this_.GetName()).SetFilename(this_.GetName()).SetSourceKind(SourceFile_SourceKindToProto(this_.GetKind())).SetIsClosureUnawareCode(this_.isClosureUnawareCode).SetNumLinesPlusOne((this_.numLines + int32(1))).SetNumBytesPlusOne((this_.numBytes + int32(1))).Build())
}

// Java static method: SourceFile.sourceKindToProto(com.google.javascript.rhino.StaticSourceFile.SourceKind sourceKind)
func SourceFile_SourceKindToProto(sourceKind pkg_rhino.StaticSourceFile_SourceKind) *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_SourceKind {
	clinit_SourceFile()
	_ = sourceKind
	switch (sourceKind).Ordinal() {
	case int32(2) /* SourceKind.EXTERN */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_SourceKind_F_EXTERN
	case int32(0) /* SourceKind.STRONG */, int32(1) /* SourceKind.WEAK */ :
		return j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_SourceKind_F_CODE
	case int32(3) /* SourceKind.NON_CODE */ :
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.BoxString(j2gort.CheckNotNullEnum(sourceKind).String()))))
	}
	panic("unreachable switch expression")
}

func (this_ *SourceFile) IsStrong() bool {
	return pkg_rhino.IfaceDefault_StaticSourceFile_IsStrong(this_)
}

func (this_ *SourceFile) IsWeak() bool {
	return pkg_rhino.IfaceDefault_StaticSourceFile_IsWeak(this_)
}

func (this_ *SourceFile) IsExtern() bool {
	return pkg_rhino.IfaceDefault_StaticSourceFile_IsExtern(this_)
}

func (this_ *SourceFile) IsNonCode() bool {
	return pkg_rhino.IfaceDefault_StaticSourceFile_IsNonCode(this_)
}

func (this_ *SourceFile) IsTypeScriptSource() bool {
	return pkg_rhino.IfaceDefault_StaticSourceFile_IsTypeScriptSource(this_)
}

// Java class: com.google.javascript.jscomp.SourceFile.Builder
type SourceFile_Builder struct {
	charset                                 j2gort.Ext_java_nio_charset_Charset    `gson:"charset"`
	originalPath                            string                                 `gson:"originalPath"`
	path                                    string                                 `gson:"path"`
	pathWithFilesystem                      j2gort.Ext_java_nio_file_Path          `gson:"pathWithFilesystem"`
	zipEntryPath                            string                                 `gson:"zipEntryPath"`
	lazyContent                             j2gort.Ext_java_util_function_Supplier `gson:"lazyContent"`
	idHash_                                 atomic.Int32
	kind                                    pkg_rhino.StaticSourceFile_SourceKind `gson:"kind"`
	isStubSourceFileForAlreadyProvidedInput bool                                  `gson:"isStubSourceFileForAlreadyProvidedInput"`
}

func (this_ *SourceFile_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceFile_Builder) IsJavaNil_() bool { return this_ == nil }
func (*SourceFile_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceFile$Builder", "Builder"
}
func (this_ *SourceFile_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceFile_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceFile.Builder", "com.google.javascript.jscomp.SourceFile$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceFile_Builder](v_)
	})
}

func (this_ *SourceFile_Builder) initFields_SourceFile_Builder() {
	_ = this_
	this_.kind = pkg_rhino.StaticGet_StaticSourceFile_SourceKind_STRONG()
	this_.charset = j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8
	this_.originalPath = j2gort.NullString

	this_.path = j2gort.NullString

	this_.zipEntryPath = j2gort.NullString

}

// Java constructor: Builder()
func NewSourceFile_Builder() *SourceFile_Builder {
	this_ := &SourceFile_Builder{}
	InitSourceFile_Builder(this_)
	return this_
}

func InitSourceFile_Builder(this_ *SourceFile_Builder) {
	this_.initFields_SourceFile_Builder()
}

// Java method: Builder.withKind(com.google.javascript.rhino.StaticSourceFile.SourceKind kind)
func (this_ *SourceFile_Builder) WithKind(kind pkg_rhino.StaticSourceFile_SourceKind) *SourceFile_Builder {
	_ = kind
	this_.kind = kind
	return this_
}

// Java method: Builder.withCharset(java.nio.charset.Charset charset)
func (this_ *SourceFile_Builder) WithCharset(charset j2gort.Ext_java_nio_charset_Charset) *SourceFile_Builder {
	_ = charset
	this_.charset = charset
	return this_
}

// Java method: Builder.withPath(java.lang.String path)
func (this_ *SourceFile_Builder) WithPath_1_string(path string) *SourceFile_Builder {
	_ = path
	return this_.priv_Builder_withPathInternal_5159(path, nil)
}

// Java method: Builder.withPath(java.nio.file.Path path)
func (this_ *SourceFile_Builder) WithPath_1_Ext_java_nio_file_Path(path j2gort.Ext_java_nio_file_Path) *SourceFile_Builder {
	_ = path
	return this_.priv_Builder_withPathInternal_5159(path.String(), path)
}

// Java method: Builder.withContent(java.lang.String x)
func (this_ *SourceFile_Builder) WithContent_1_string(x string) *SourceFile_Builder {
	_ = x
	this_.lazyContent = func() j2gort.Ext_java_util_function_Supplier {
		bound_41 := x
		return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.CheckNotNullString(bound_41) })
	}()
	return this_
}

// Java method: Builder.withContent(java.io.InputStream x)
func (this_ *SourceFile_Builder) WithContent_1_Ext_java_io_InputStream(x j2gort.Ext_java_io_InputStream) *SourceFile_Builder {
	_ = x
	this_.lazyContent = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		if !(this_.charset != nil) {
			j2gort.CheckState(false)
		}
		var didRet_42 bool
		var retVal_43 any
		func() {
			defer func() {
				if rec_44 := recover(); rec_44 != nil {
					jex_45 := javalang.RecoverException(rec_44)
					_ = jex_45
					if j2gort.InstanceOf[javaio.IIOException](jex_45) {
						e := j2gort.AsIface[javaio.IIOException](jex_45)
						_ = e
						panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
					} else {
						panic(rec_44)
					}
				}
			}()
			retVal_43 = j2gort.BoxString(j2gort.Ext_com_google_common_io_CharStreams_ToString(j2gort.NewExt_java_io_InputStreamReader_2_Ext_java_io_InputStream_Ext_java_nio_charset_Charset(x, this_.charset)))
			didRet_42 = true
			return
		}()
		if didRet_42 {
			return retVal_43
		}
		panic("unreachable")
	})
	return this_
}

// Java method: Builder.withZipEntryPath(java.lang.String zipPath, java.lang.String entryPath)
func (this_ *SourceFile_Builder) WithZipEntryPath(zipPath string, entryPath string) *SourceFile_Builder {
	_ = zipPath
	_ = entryPath
	this_.path = zipPath
	this_.zipEntryPath = entryPath
	return this_
}

// Java method: Builder.withOriginalPath(java.lang.String originalPath)
func (this_ *SourceFile_Builder) WithOriginalPath(originalPath string) *SourceFile_Builder {
	_ = originalPath
	this_.originalPath = originalPath
	return this_
}

// Java method: Builder.setIsStubSourceFileForAlreadyProvidedInput()
func (this_ *SourceFile_Builder) SetIsStubSourceFileForAlreadyProvidedInput() *SourceFile_Builder {
	this_.isStubSourceFileForAlreadyProvidedInput = true
	return this_
}

// Java method: Builder.build()
func (this_ *SourceFile_Builder) Build() *SourceFile {
	var displayPath string
	_ = displayPath
	if this_.originalPath != j2gort.NullString {
		displayPath = this_.originalPath
	} else if this_.zipEntryPath == j2gort.NullString {
		displayPath = this_.path
	} else {
		displayPath = (j2gort.StrValString(this_.path) + j2gort.StrValString(StaticGet_SourceFile_BANG_SLASH()) + j2gort.StrValString(this_.zipEntryPath))
	}
	if this_.isStubSourceFileForAlreadyProvidedInput {
		clinit_SourceFile()
		var val_46 SourceFile
		InitSourceFile(&val_46, NewSourceFile_CodeLoader_StubSourceFileCodeLoader(), displayPath, this_.kind)
		var file *SourceFile = &val_46
		_ = file
		file.priv_SourceFile_setIsStubSourceFileForAlreadyProvidedInput_5124(true)
		return file
	}
	if this_.lazyContent != nil {
		return NewSourceFile(NewSourceFile_CodeLoader_Preloaded(j2gort.CastString(this_.lazyContent.Get())), displayPath, this_.kind)
	}
	if this_.zipEntryPath != j2gort.NullString {
		return NewSourceFile(NewSourceFile_CodeLoader_AtZip(this_.path, this_.zipEntryPath, this_.charset), displayPath, this_.kind)
	}
	var iifeArg_1 j2gort.
		Ext_java_nio_file_Path
iifeSwitch_1:
	switch {
	default:

		if this_.pathWithFilesystem != nil {
			iifeArg_1 = this_.pathWithFilesystem
			break iifeSwitch_1
		}
		iifeArg_1 = j2gort.Ext_java_nio_file_Path_Of_2(this_.path, []string{})
		break iifeSwitch_1
	}
	return NewSourceFile(NewSourceFile_CodeLoader_OnDisk(iifeArg_1, this_.charset), displayPath, this_.kind)
}

// Java method: Builder.withPathInternal(java.lang.String path, java.nio.file.@org.jspecify.annotations.Nullable Path pathWithFilesystem)
func (this_ *SourceFile_Builder) priv_Builder_withPathInternal_5159(path string, pathWithFilesystem j2gort.Ext_java_nio_file_Path) *SourceFile_Builder {
	_ = path
	_ = pathWithFilesystem
	var bangSlashIndex int32 = j2gort.JStrIndexOf(path, StaticGet_SourceFile_BANG_SLASH())
	_ = bangSlashIndex
	if bangSlashIndex >= int32(0) {
		var zipPath string = j2gort.JStrSub(path, int32(0), bangSlashIndex)
		_ = zipPath
		var entryPath string = j2gort.JStrSubFrom(path, (bangSlashIndex + j2gort.JStrLen(StaticGet_SourceFile_BANG_SLASH())))
		_ = entryPath
		if javalang.StringEndsWith(j2gort.CheckNotNullString(zipPath), ".zip") && (javalang.StringEndsWith(j2gort.CheckNotNullString(entryPath), ".js") || javalang.StringEndsWith(j2gort.CheckNotNullString(entryPath), ".js.map")) {
			return this_.WithZipEntryPath(zipPath, entryPath)
		}
	}
	this_.path = path
	this_.pathWithFilesystem = pathWithFilesystem
	return this_
}

// Java polymorphic class interface: com.google.javascript.jscomp.SourceFile.CodeLoader
type SourceFile_CodeLoader interface {
	Implements_jscomp_SourceFile_CodeLoader()
	IsJavaNil_() bool
	GetSourceFile_CodeLoaderBase() *SourceFile_CodeLoaderBase
	LoadUncachedCode() string
	OpenUncachedReader() j2gort.Ext_java_io_Reader
	ToProtoLocationBuilder(string) *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder
}

func ToIface_SourceFile_CodeLoader[T interface {
	SourceFile_CodeLoader
	comparable
}](v T) SourceFile_CodeLoader {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SourceFile_CodeLoader[T SourceFile_CodeLoader](v T) SourceFile_CodeLoader {
	var vi SourceFile_CodeLoader = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.SourceFile.CodeLoader
type SourceFile_CodeLoaderBase struct {
	idHash_ atomic.Int32
}

func (this_ *SourceFile_CodeLoaderBase) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceFile_CodeLoaderBase) IsJavaNil_() bool { return this_ == nil }
func (*SourceFile_CodeLoaderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceFile$CodeLoader", "CodeLoader"
}
func (this_ *SourceFile_CodeLoaderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceFile.CodeLoader", "com.google.javascript.jscomp.SourceFile$CodeLoader", func(v_ any) bool {
		return j2gort.InstanceOf[SourceFile_CodeLoader](v_)
	})
}
func (b *SourceFile_CodeLoaderBase) GetSourceFile_CodeLoaderBase() *SourceFile_CodeLoaderBase {
	return b
}
func (*SourceFile_CodeLoaderBase) Implements_jscomp_SourceFile_CodeLoader() {}

func (this_ *SourceFile_CodeLoaderBase) initFields_SourceFile_CodeLoader() {
	_ = this_
}

func InitSourceFile_CodeLoader(this_ *SourceFile_CodeLoaderBase) {
	this_.initFields_SourceFile_CodeLoader()
}

// Java method: CodeLoader.loadUncachedCode()
func (this_ *SourceFile_CodeLoaderBase) LoadUncachedCode() string {
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java method: CodeLoader.openUncachedReader()
func (this_ *SourceFile_CodeLoaderBase) OpenUncachedReader() j2gort.Ext_java_io_Reader {
	return nil
}

// Java method: CodeLoader.toProtoLocationBuilder(java.lang.String fileName)
func (this_ *SourceFile_CodeLoaderBase) ToProtoLocationBuilder(fileName string) *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder {
	_ = fileName
	panic("abstract or native method invoked")
}

// Java class: com.google.javascript.jscomp.SourceFile.CodeLoader.Preloaded
type SourceFile_CodeLoader_Preloaded struct {
	SourceFile_CodeLoaderBase
	preloadedCode string `gson:"preloadedCode"`
}

func (this_ *SourceFile_CodeLoader_Preloaded) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceFile_CodeLoader_Preloaded) IsJavaNil_() bool { return this_ == nil }
func (*SourceFile_CodeLoader_Preloaded) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceFile$CodeLoader$Preloaded", "Preloaded"
}
func (this_ *SourceFile_CodeLoader_Preloaded) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SourceFile_CodeLoaderBase.IdentityHash_()
}
func (*SourceFile_CodeLoader_Preloaded) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceFile.CodeLoader.Preloaded", "com.google.javascript.jscomp.SourceFile$CodeLoader$Preloaded", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceFile_CodeLoader_Preloaded](v_)
	})
}
func (*SourceFile_CodeLoader_Preloaded) Implements_jscomp_SourceFile_CodeLoader() {}

func StaticGet_SourceFile_CodeLoader_Preloaded_SerialVersionUID() int64 { return int64(2) }
func (this_ *SourceFile_CodeLoader_Preloaded) initFields_SourceFile_CodeLoader_Preloaded() {
	_ = this_
	this_.preloadedCode = j2gort.NullString
}

// Java constructor: Preloaded(java.lang.String preloadedCode)
func NewSourceFile_CodeLoader_Preloaded(preloadedCode string) *SourceFile_CodeLoader_Preloaded {
	this_ := &SourceFile_CodeLoader_Preloaded{}
	InitSourceFile_CodeLoader_Preloaded(this_, preloadedCode)
	return this_
}

func InitSourceFile_CodeLoader_Preloaded(this_ *SourceFile_CodeLoader_Preloaded, preloadedCode string) {
	_ = preloadedCode
	InitSourceFile_CodeLoader(&this_.SourceFile_CodeLoaderBase)
	this_.initFields_SourceFile_CodeLoader_Preloaded()
	this_.preloadedCode = j2gort.CheckNotNullString(preloadedCode)
}

// Java method: Preloaded.loadUncachedCode()
func (this_ *SourceFile_CodeLoader_Preloaded) LoadUncachedCode() string {
	return this_.preloadedCode
}

// Java method: Preloaded.toProtoLocationBuilder(java.lang.String fileName)
func (this_ *SourceFile_CodeLoader_Preloaded) ToProtoLocationBuilder(fileName string) *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder {
	_ = fileName
	return j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_NewBuilder_0().SetPreloadedContents(this_.preloadedCode)
}

// Java class: com.google.javascript.jscomp.SourceFile.CodeLoader.StubSourceFileCodeLoader
type SourceFile_CodeLoader_StubSourceFileCodeLoader struct {
	SourceFile_CodeLoaderBase
}

func (this_ *SourceFile_CodeLoader_StubSourceFileCodeLoader) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceFile_CodeLoader_StubSourceFileCodeLoader) IsJavaNil_() bool { return this_ == nil }
func (*SourceFile_CodeLoader_StubSourceFileCodeLoader) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceFile$CodeLoader$StubSourceFileCodeLoader", "StubSourceFileCodeLoader"
}
func (this_ *SourceFile_CodeLoader_StubSourceFileCodeLoader) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SourceFile_CodeLoaderBase.IdentityHash_()
}
func (*SourceFile_CodeLoader_StubSourceFileCodeLoader) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceFile.CodeLoader.StubSourceFileCodeLoader", "com.google.javascript.jscomp.SourceFile$CodeLoader$StubSourceFileCodeLoader", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceFile_CodeLoader_StubSourceFileCodeLoader](v_)
	})
}
func (*SourceFile_CodeLoader_StubSourceFileCodeLoader) Implements_jscomp_SourceFile_CodeLoader() {}

func StaticGet_SourceFile_CodeLoader_StubSourceFileCodeLoader_SerialVersionUID() int64 {
	return int64(2)
}
func (this_ *SourceFile_CodeLoader_StubSourceFileCodeLoader) initFields_SourceFile_CodeLoader_StubSourceFileCodeLoader() {
	_ = this_
}

// Java constructor: StubSourceFileCodeLoader()
func NewSourceFile_CodeLoader_StubSourceFileCodeLoader() *SourceFile_CodeLoader_StubSourceFileCodeLoader {
	this_ := &SourceFile_CodeLoader_StubSourceFileCodeLoader{}
	InitSourceFile_CodeLoader_StubSourceFileCodeLoader(this_)
	return this_
}

func InitSourceFile_CodeLoader_StubSourceFileCodeLoader(this_ *SourceFile_CodeLoader_StubSourceFileCodeLoader) {
	InitSourceFile_CodeLoader(&this_.SourceFile_CodeLoaderBase)
	this_.initFields_SourceFile_CodeLoader_StubSourceFileCodeLoader()
}

// Java method: StubSourceFileCodeLoader.loadUncachedCode()
func (this_ *SourceFile_CodeLoader_StubSourceFileCodeLoader) LoadUncachedCode() string {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Attempting to load code from a stub SourceFile.")))
}

// Java method: StubSourceFileCodeLoader.toProtoLocationBuilder(java.lang.String fileName)
func (this_ *SourceFile_CodeLoader_StubSourceFileCodeLoader) ToProtoLocationBuilder(fileName string) *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder {
	_ = fileName
	return j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_NewBuilder_0().SetFilename(fileName).SetStubFile(true)
}

// Java class: com.google.javascript.jscomp.SourceFile.CodeLoader.OnDisk
type SourceFile_CodeLoader_OnDisk struct {
	SourceFile_CodeLoaderBase
	serializableCharset string                        `gson:"serializableCharset"`
	relativePath        j2gort.Ext_java_nio_file_Path `gson:"relativePath"`
}

func (this_ *SourceFile_CodeLoader_OnDisk) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceFile_CodeLoader_OnDisk) IsJavaNil_() bool { return this_ == nil }
func (*SourceFile_CodeLoader_OnDisk) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceFile$CodeLoader$OnDisk", "OnDisk"
}
func (this_ *SourceFile_CodeLoader_OnDisk) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SourceFile_CodeLoaderBase.IdentityHash_()
}
func (*SourceFile_CodeLoader_OnDisk) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceFile.CodeLoader.OnDisk", "com.google.javascript.jscomp.SourceFile$CodeLoader$OnDisk", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceFile_CodeLoader_OnDisk](v_)
	})
}
func (*SourceFile_CodeLoader_OnDisk) Implements_jscomp_SourceFile_CodeLoader() {}

func StaticGet_SourceFile_CodeLoader_OnDisk_SerialVersionUID() int64 { return int64(1) }
func (this_ *SourceFile_CodeLoader_OnDisk) initFields_SourceFile_CodeLoader_OnDisk() {
	_ = this_
	this_.serializableCharset = j2gort.NullString
}

// Java constructor: OnDisk(java.nio.file.Path relativePath, java.nio.charset.Charset c)
func NewSourceFile_CodeLoader_OnDisk(relativePath j2gort.Ext_java_nio_file_Path, c j2gort.Ext_java_nio_charset_Charset) *SourceFile_CodeLoader_OnDisk {
	this_ := &SourceFile_CodeLoader_OnDisk{}
	InitSourceFile_CodeLoader_OnDisk(this_, relativePath, c)
	return this_
}

func InitSourceFile_CodeLoader_OnDisk(this_ *SourceFile_CodeLoader_OnDisk, relativePath j2gort.Ext_java_nio_file_Path, c j2gort.Ext_java_nio_charset_Charset) {
	_ = relativePath
	_ = c
	InitSourceFile_CodeLoader(&this_.SourceFile_CodeLoaderBase)
	this_.initFields_SourceFile_CodeLoader_OnDisk()
	this_.serializableCharset = c.Name()
	this_.relativePath = relativePath
}

// Java method: OnDisk.loadUncachedCode()
func (this_ *SourceFile_CodeLoader_OnDisk) LoadUncachedCode() string {
	var didRet_47 bool
	var retVal_48 string
	func() {
		defer func() {
			if rec_49 := recover(); rec_49 != nil {
				jex_50 := javalang.RecoverException(rec_49)
				_ = jex_50
				if j2gort.InstanceOf[j2gort.Ext_java_nio_charset_CharacterCodingException](jex_50) {
					e := j2gort.AsIface[j2gort.Ext_java_nio_charset_CharacterCodingException](jex_50)
					_ = e
					panic(j2gort.AsIface[javaio.IIOException](j2gort.ConstructShimThrowable("java.io.IOException", "Failed to read: "+j2gort.StrVal(this_.relativePath)+", is this input UTF-8 encoded?", e)))
				} else {
					panic(rec_49)
				}
			}
		}()
		retVal_48 = j2gort.Ext_java_nio_file_Files_ReadString_2(this_.relativePath, this_.priv_OnDisk_getCharset_5163())
		didRet_47 = true
		return
	}()
	if didRet_47 {
		return retVal_48
	}
	panic("unreachable")
}

// Java method: OnDisk.openUncachedReader()
func (this_ *SourceFile_CodeLoader_OnDisk) OpenUncachedReader() j2gort.Ext_java_io_Reader {
	return j2gort.AsIface[j2gort.Ext_java_io_Reader](j2gort.Ext_java_nio_file_Files_NewBufferedReader_2(this_.relativePath, this_.priv_OnDisk_getCharset_5163()))
}

// Java method: OnDisk.getCharset()
func (this_ *SourceFile_CodeLoader_OnDisk) priv_OnDisk_getCharset_5163() j2gort.Ext_java_nio_charset_Charset {
	return j2gort.Ext_java_nio_charset_Charset_ForName_1(this_.serializableCharset)
}

// Java method: OnDisk.toProtoLocationBuilder(java.lang.String fileName)
func (this_ *SourceFile_CodeLoader_OnDisk) ToProtoLocationBuilder(fileName string) *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder {
	_ = fileName
	var actualPath string = this_.relativePath.String()
	_ = actualPath
	return j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_NewBuilder_0().SetFileOnDisk_1_Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk_Builder(j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk_NewBuilder_0().SetActualPath(func() string {
		if j2gort.CheckNotNullString(fileName) == actualPath {
			return ""
		}
		return actualPath
	}()).SetCharset(func() string {
		if this_.priv_OnDisk_getCharset_5163().Equals(j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8) {
			return ""
		}
		return this_.serializableCharset
	}()))
}

// Java class: com.google.javascript.jscomp.SourceFile.CodeLoader.AtZip
type SourceFile_CodeLoader_AtZip struct {
	SourceFile_CodeLoaderBase
	zipName             string `gson:"zipName"`
	entryName           string `gson:"entryName"`
	serializableCharset string `gson:"serializableCharset"`
}

func (this_ *SourceFile_CodeLoader_AtZip) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceFile_CodeLoader_AtZip) IsJavaNil_() bool { return this_ == nil }
func (*SourceFile_CodeLoader_AtZip) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceFile$CodeLoader$AtZip", "AtZip"
}
func (this_ *SourceFile_CodeLoader_AtZip) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SourceFile_CodeLoaderBase.IdentityHash_()
}
func (*SourceFile_CodeLoader_AtZip) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceFile.CodeLoader.AtZip", "com.google.javascript.jscomp.SourceFile$CodeLoader$AtZip", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceFile_CodeLoader_AtZip](v_)
	})
}
func (*SourceFile_CodeLoader_AtZip) Implements_jscomp_SourceFile_CodeLoader() {}

func StaticGet_SourceFile_CodeLoader_AtZip_SerialVersionUID() int64 { return int64(1) }
func (this_ *SourceFile_CodeLoader_AtZip) initFields_SourceFile_CodeLoader_AtZip() {
	_ = this_
	this_.zipName = j2gort.NullString
	this_.entryName = j2gort.NullString
	this_.serializableCharset = j2gort.NullString
}

// Java constructor: AtZip(java.lang.String zipName, java.lang.String entryName, java.nio.charset.Charset c)
func NewSourceFile_CodeLoader_AtZip(zipName string, entryName string, c j2gort.Ext_java_nio_charset_Charset) *SourceFile_CodeLoader_AtZip {
	this_ := &SourceFile_CodeLoader_AtZip{}
	InitSourceFile_CodeLoader_AtZip(this_, zipName, entryName, c)
	return this_
}

func InitSourceFile_CodeLoader_AtZip(this_ *SourceFile_CodeLoader_AtZip, zipName string, entryName string, c j2gort.Ext_java_nio_charset_Charset) {
	_ = zipName
	_ = entryName
	_ = c
	InitSourceFile_CodeLoader(&this_.SourceFile_CodeLoaderBase)
	this_.initFields_SourceFile_CodeLoader_AtZip()
	this_.zipName = zipName
	this_.entryName = entryName
	this_.serializableCharset = c.Name()
}

// Java method: AtZip.loadUncachedCode()
func (this_ *SourceFile_CodeLoader_AtZip) LoadUncachedCode() string {
	return j2gort.Ext_com_google_common_io_CharStreams_ToString(j2gort.AsIface[j2gort.Ext_java_lang_Readable](this_.OpenUncachedReader()))
}

// Java method: AtZip.openUncachedReader()
func (this_ *SourceFile_CodeLoader_AtZip) OpenUncachedReader() j2gort.Ext_java_io_Reader {
	return j2gort.NewExt_java_io_InputStreamReader_2_Ext_java_io_InputStream_Ext_java_nio_charset_Charset(JSCompZipFileCache_GetEntryStream(this_.zipName, this_.entryName), this_.priv_AtZip_getCharset_5164())
}

// Java method: AtZip.getCharset()
func (this_ *SourceFile_CodeLoader_AtZip) priv_AtZip_getCharset_5164() j2gort.Ext_java_nio_charset_Charset {
	return j2gort.Ext_java_nio_charset_Charset_ForName_1(this_.serializableCharset)
}

// Java method: AtZip.toProtoLocationBuilder(java.lang.String fileName)
func (this_ *SourceFile_CodeLoader_AtZip) ToProtoLocationBuilder(fileName string) *j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder {
	_ = fileName
	return j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_NewBuilder_0().SetFilename(fileName).SetZipEntry_1_Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk](j2gort.Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk_NewBuilder_0().SetEntryName(this_.entryName).SetZipPath(this_.zipName).SetCharset(func() string {
		if this_.priv_AtZip_getCharset_5164().Equals(j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8) {
			return ""
		}
		return this_.serializableCharset
	}()).Build()))
}
