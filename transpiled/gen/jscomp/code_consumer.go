/*
 * Copyright 2004 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CodeConsumer.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.CodeConsumer
type CodeConsumer interface {
	Implements_jscomp_CodeConsumer()
	IsJavaNil_() bool
	GetCodeConsumerBase() *CodeConsumerBase
	StartSourceMapping(pkg_rhino.Node)
	EndSourceMapping(pkg_rhino.Node)
	TrackLicenses(pkg_rhino.Node)
	ContinueProcessing() bool
	GetLastChar() uint16
	Add_1_string(string)
	Append(string)
	AddIdentifier(string)
	AppendBlockStart()
	AppendBlockEnd()
	StartNewLine()
	MaybeLineBreak()
	MaybeCutLine()
	EndLine()
	NotePreferredLineBreak()
	BeginBlock()
	EndBlock_0()
	EndBlock_1(bool)
	ListSeparator()
	OptionalListSeparator()
	EndStatement_1(bool)
	EndStatement_2(bool, bool)
	MaybeEndStatement()
	EndFunction(bool)
	EndClass(bool)
	BeginCaseBody()
	EndCaseBody()
	BeginTemplateLit()
	BeginTemplateLitSub()
	EndTemplateLitSub()
	EndTemplateLit()
	IsInTemplateLiteral() bool
	AppendOp(string, bool)
	AddOp(string, bool)
	AddNumber(float64, pkg_rhino.Node)
	AddBigInt(*j2gort.Ext_java_math_BigInteger)
	AddConstant(string)
	ShouldPreserveExtras(pkg_rhino.Node) bool
	MaybeInsertSpace()
	BreakAfterBlockFor(pkg_rhino.Node, bool) bool
	EndFile()
}

func ToIface_CodeConsumer[T interface {
	CodeConsumer
	comparable
}](v T) CodeConsumer {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_CodeConsumer[T CodeConsumer](v T) CodeConsumer {
	var vi CodeConsumer = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.CodeConsumer
type CodeConsumerBase struct {
	self_               CodeConsumer
	templateLitDepth    int32 `gson:"templateLitDepth"`
	templateLitSubDepth int32 `gson:"templateLitSubDepth"`
	idHash_             atomic.Int32
	statementNeedsEnded bool `gson:"statementNeedsEnded"`
	statementStarted    bool `gson:"statementStarted"`
	sawFunction         bool `gson:"sawFunction"`
}

func (this_ *CodeConsumerBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(CodeConsumer)
}
func (this_ *CodeConsumerBase) IsJavaNil_() bool { return this_ == nil }
func (*CodeConsumerBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodeConsumer", "CodeConsumer"
}
func (this_ *CodeConsumerBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodeConsumer", "com.google.javascript.jscomp.CodeConsumer", func(v_ any) bool {
		return j2gort.InstanceOf[CodeConsumer](v_)
	})
}
func (b *CodeConsumerBase) GetCodeConsumerBase() *CodeConsumerBase { return b }
func (*CodeConsumerBase) Implements_jscomp_CodeConsumer()          {}

func (this_ *CodeConsumerBase) initFields_CodeConsumer() {
	_ = this_

}

func InitCodeConsumer(this_ *CodeConsumerBase) {
	this_.initFields_CodeConsumer()
}

// Java method: CodeConsumer.startSourceMapping(com.google.javascript.rhino.Node node)
func (this_ *CodeConsumerBase) StartSourceMapping(node pkg_rhino.Node) {
	_ = node
}

// Java method: CodeConsumer.endSourceMapping(com.google.javascript.rhino.Node node)
func (this_ *CodeConsumerBase) EndSourceMapping(node pkg_rhino.Node) {
	_ = node
}

// Java method: CodeConsumer.trackLicenses(com.google.javascript.rhino.Node node)
func (this_ *CodeConsumerBase) TrackLicenses(node pkg_rhino.Node) {
	_ = node
}

// Java method: CodeConsumer.continueProcessing()
func (this_ *CodeConsumerBase) ContinueProcessing() bool {
	return true
}

// Java method: CodeConsumer.getLastChar()
func (this_ *CodeConsumerBase) GetLastChar() uint16 {
	panic("abstract or native method invoked")
}

// Java method: CodeConsumer.add(java.lang.String newcode)
func (this_ *CodeConsumerBase) Add_1_string(newcode string) {
	_ = newcode
	this_.MaybeEndStatement()
	if j2gort.CheckNotNullString(newcode) == "" {
		return
	}
	var c uint16 = j2gort.JStrCharAt(newcode, int32(0))
	_ = c
	if (CodeConsumer_IsWordChar(c) || (int32(c) == int32(uint16(92)))) && CodeConsumer_IsWordChar(this_.self_.GetLastChar()) {
		this_.self_.Append(" ")
	} else if (int32(c) == int32(uint16(47))) && ((int32(this_.self_.GetLastChar()) == int32(uint16(47))) || (int32(this_.self_.GetLastChar()) == int32(uint16(60)))) {
		this_.self_.Append(" ")
	} else if ((int32(c) == int32(uint16(34))) || (int32(c) == int32(uint16(39)))) && CodeConsumer_IsWordChar(this_.self_.GetLastChar()) {
		this_.self_.MaybeInsertSpace()
	}
	var startOfLine int32 = int32(0)
	_ = startOfLine
	var endOfLine int32 = j2gort.JStrIndexOfCp(newcode, int32(uint16(10)))
	_ = endOfLine
	for endOfLine >= int32(0) {
		if endOfLine > startOfLine {
			this_.self_.Append(j2gort.JStrSub(newcode, startOfLine, endOfLine))
		}
		this_.self_.StartNewLine()
		startOfLine = (endOfLine + int32(1))
		endOfLine = j2gort.JStrIndexOfCpFrom(newcode, int32(uint16(10)), startOfLine)
	}
	if j2gort.JStrLen(newcode) > startOfLine {
		this_.self_.Append(j2gort.JStrSubFrom(newcode, startOfLine))
	}
}

// Java method: CodeConsumer.append(java.lang.String str)
func (this_ *CodeConsumerBase) Append(str string) {
	_ = str
	panic("abstract or native method invoked")
}

// Java method: CodeConsumer.addIdentifier(java.lang.String identifier)
func (this_ *CodeConsumerBase) AddIdentifier(identifier string) {
	_ = identifier
	this_.Add_1_string(identifier)
}

// Java method: CodeConsumer.appendBlockStart()
func (this_ *CodeConsumerBase) AppendBlockStart() {
	this_.self_.Append("{")
}

// Java method: CodeConsumer.appendBlockEnd()
func (this_ *CodeConsumerBase) AppendBlockEnd() {
	this_.self_.Append("}")
}

// Java method: CodeConsumer.startNewLine()
func (this_ *CodeConsumerBase) StartNewLine() {
}

// Java method: CodeConsumer.maybeLineBreak()
func (this_ *CodeConsumerBase) MaybeLineBreak() {
	this_.self_.MaybeCutLine()
}

// Java method: CodeConsumer.maybeCutLine()
func (this_ *CodeConsumerBase) MaybeCutLine() {
}

// Java method: CodeConsumer.endLine()
func (this_ *CodeConsumerBase) EndLine() {
}

// Java method: CodeConsumer.notePreferredLineBreak()
func (this_ *CodeConsumerBase) NotePreferredLineBreak() {
}

// Java method: CodeConsumer.beginBlock()
func (this_ *CodeConsumerBase) BeginBlock() {
	if this_.statementNeedsEnded {
		this_.self_.Append(";")
		this_.self_.MaybeLineBreak()
	}
	this_.self_.AppendBlockStart()
	this_.self_.EndLine()
	this_.statementNeedsEnded = false
}

// Java method: CodeConsumer.endBlock()
func (this_ *CodeConsumerBase) EndBlock_0() {
	this_.EndBlock_1(false)
}

// Java method: CodeConsumer.endBlock(boolean shouldEndLine)
func (this_ *CodeConsumerBase) EndBlock_1(shouldEndLine bool) {
	_ = shouldEndLine
	this_.self_.AppendBlockEnd()
	if shouldEndLine {
		this_.self_.EndLine()
	}
	this_.statementNeedsEnded = false
}

// Java method: CodeConsumer.listSeparator()
func (this_ *CodeConsumerBase) ListSeparator() {
	this_.Add_1_string(",")
	this_.self_.MaybeLineBreak()
}

// Java method: CodeConsumer.optionalListSeparator()
func (this_ *CodeConsumerBase) OptionalListSeparator() {
}

// Java method: CodeConsumer.endStatement(boolean hasTrailingCommentOnSameLine)
func (this_ *CodeConsumerBase) EndStatement_1(hasTrailingCommentOnSameLine bool) {
	_ = hasTrailingCommentOnSameLine
	this_.self_.EndStatement_2(false, hasTrailingCommentOnSameLine)
}

// Java method: CodeConsumer.endStatement(boolean needSemiColon, boolean hasTrailingCommentOnSameLine)
func (this_ *CodeConsumerBase) EndStatement_2(needSemiColon bool, hasTrailingCommentOnSameLine bool) {
	_ = needSemiColon
	_ = hasTrailingCommentOnSameLine
	if needSemiColon {
		this_.self_.Append(";")
		if !(hasTrailingCommentOnSameLine) {
			this_.self_.MaybeLineBreak()
		}
		this_.statementNeedsEnded = false
	} else if this_.statementStarted {
		this_.statementNeedsEnded = true
	}
}

// Java method: CodeConsumer.maybeEndStatement()
func (this_ *CodeConsumerBase) MaybeEndStatement() {
	if this_.statementNeedsEnded {
		this_.self_.Append(";")
		this_.self_.MaybeLineBreak()
		this_.self_.EndLine()
		this_.statementNeedsEnded = false
	}
	this_.statementStarted = true
}

// Java method: CodeConsumer.endFunction(boolean statementContext)
func (this_ *CodeConsumerBase) EndFunction(statementContext bool) {
	_ = statementContext
	this_.sawFunction = true
	if statementContext {
		this_.self_.EndLine()
	}
}

// Java method: CodeConsumer.endClass(boolean statementContext)
func (this_ *CodeConsumerBase) EndClass(statementContext bool) {
	_ = statementContext
	if statementContext {
		this_.self_.EndLine()
	}
}

// Java method: CodeConsumer.beginCaseBody()
func (this_ *CodeConsumerBase) BeginCaseBody() {
	this_.self_.Append(":")
}

// Java method: CodeConsumer.endCaseBody()
func (this_ *CodeConsumerBase) EndCaseBody() {
}

// Java method: CodeConsumer.beginTemplateLit()
func (this_ *CodeConsumerBase) BeginTemplateLit() {
	if !(this_.templateLitDepth == this_.templateLitSubDepth) {
		j2gort.CheckState(false)
	}
	this_.MaybeEndStatement()
	this_.self_.Append("`")
	this_.templateLitDepth = this_.templateLitDepth + 1
}

// Java method: CodeConsumer.beginTemplateLitSub()
func (this_ *CodeConsumerBase) BeginTemplateLitSub() {
	if !(this_.IsInTemplateLiteral()) {
		j2gort.CheckState(false)
	}
	this_.self_.Append("${")
	this_.templateLitSubDepth = this_.templateLitSubDepth + 1
}

// Java method: CodeConsumer.endTemplateLitSub()
func (this_ *CodeConsumerBase) EndTemplateLitSub() {
	if !(this_.templateLitSubDepth > int32(0)) {
		j2gort.CheckState(false)
	}
	if !(this_.templateLitDepth == this_.templateLitSubDepth) {
		j2gort.CheckState(false)
	}
	this_.self_.Append("}")
	this_.templateLitSubDepth = this_.templateLitSubDepth - 1
}

// Java method: CodeConsumer.endTemplateLit()
func (this_ *CodeConsumerBase) EndTemplateLit() {
	if !(this_.templateLitDepth > int32(0)) {
		j2gort.CheckState(false)
	}
	if !(this_.IsInTemplateLiteral()) {
		j2gort.CheckState(false)
	}
	this_.self_.Append("`")
	this_.templateLitDepth = this_.templateLitDepth - 1
}

// Java method: CodeConsumer.isInTemplateLiteral()
func (this_ *CodeConsumerBase) IsInTemplateLiteral() bool {
	return (this_.templateLitDepth == (this_.templateLitSubDepth + int32(1)))
}

// Java method: CodeConsumer.appendOp(java.lang.String op, boolean binOp)
func (this_ *CodeConsumerBase) AppendOp(op string, binOp bool) {
	_ = op
	_ = binOp
	this_.self_.Append(op)
}

// Java method: CodeConsumer.addOp(java.lang.String op, boolean binOp)
func (this_ *CodeConsumerBase) AddOp(op string, binOp bool) {
	_ = op
	_ = binOp
	this_.MaybeEndStatement()
	var first uint16 = j2gort.JStrCharAt(op, int32(0))
	_ = first
	var prev uint16 = this_.self_.GetLastChar()
	_ = prev
	if ((int32(first) == int32(uint16(43))) || (int32(first) == int32(uint16(45)))) && (int32(prev) == int32(first)) {
		this_.self_.Append(" ")
	} else if j2gort.BoxChar_IsLetter_1_uint16(first) && CodeConsumer_IsWordChar(prev) {
		this_.self_.Append(" ")
	} else if ((int32(prev) == int32(uint16(45))) && (int32(first) == int32(uint16(62)))) || ((int32(prev) == int32(uint16(60))) && (int32(first) == int32(uint16(33)))) {
		this_.self_.Append(" ")
	}
	this_.self_.AppendOp(op, binOp)
	if binOp {
		this_.self_.MaybeCutLine()
	}
}

// Java method: CodeConsumer.addNumber(double x, com.google.javascript.rhino.Node n)
func (this_ *CodeConsumerBase) AddNumber(x float64, n pkg_rhino.Node) {
	_ = x
	_ = n
	if !(pkg_base.JSCompDoubles_IsPositive(x)) {
		j2gort.CheckStateMsg(false, j2gort.ValueOfDouble(x))
	}
	this_.self_.AddConstant(CodeConsumer_FormatNumber(x))
}

// Java static method: CodeConsumer.formatNumber(double x)
func CodeConsumer_FormatNumber(x float64) string {
	_ = x
	if float64(javalang.D2L(x)) != x {
		return j2gort.JStrReplaceFirst(javalang.StringReplace(javalang.DoubleToString(x), ".0E", "E"), "^0\\.", ".")
	}
	var value int64 = javalang.D2L(x)
	_ = value
	var mantissa int64 = value
	_ = mantissa
	var exp int32 = int32(0)
	_ = exp
	if x >= float64(int32(100)) {
		for (mantissa % int64(int32(10))) == int64(int32(0)) {
			mantissa = (mantissa / int64(int32(10)))
			exp = exp + 1
		}
	}
	if exp > int32(2) {
		return (j2gort.StrValInt64(mantissa) + "E" + j2gort.StrValInt32(exp))
	}
	var decValueString string = strconv.FormatInt(value, 10)
	_ = decValueString
	if value <= int64(1000000000000) {
		return decValueString
	}
	var hexValueString string = javalang.LongToHexString(value)
	_ = hexValueString
	if (j2gort.JStrLen(hexValueString) + int32(2)) < j2gort.JStrLen(decValueString) {
		return ("0x" + j2gort.StrValString(hexValueString))
	} else {
		return decValueString
	}
}

// Java method: CodeConsumer.addBigInt(java.math.BigInteger bi)
func (this_ *CodeConsumerBase) AddBigInt(bi *j2gort.Ext_java_math_BigInteger) {
	_ = bi
	var hexEncoded string = ("0x" + j2gort.StrValString(bi.ToString_1(int32(16))) + "n")
	_ = hexEncoded
	var decimalEncoded string = (j2gort.StrVal(bi) + "n")
	_ = decimalEncoded
	this_.self_.AddConstant(func() string {
		if j2gort.JStrLen(hexEncoded) < j2gort.JStrLen(decimalEncoded) {
			return hexEncoded
		}
		return decimalEncoded
	}())
}

// Java method: CodeConsumer.addConstant(java.lang.String newcode)
func (this_ *CodeConsumerBase) AddConstant(newcode string) {
	_ = newcode
	this_.Add_1_string(newcode)
}

// Java static method: CodeConsumer.isWordChar(char ch)
func CodeConsumer_IsWordChar(ch uint16) bool {
	_ = ch
	return (((int32(ch) == int32(uint16(95))) || (int32(ch) == int32(uint16(36)))) || j2gort.BoxChar_IsLetterOrDigit_1_uint16(ch))
}

// Java method: CodeConsumer.shouldPreserveExtras(com.google.javascript.rhino.Node n)
func (this_ *CodeConsumerBase) ShouldPreserveExtras(n pkg_rhino.Node) bool {
	_ = n
	return false
}

// Java method: CodeConsumer.maybeInsertSpace()
func (this_ *CodeConsumerBase) MaybeInsertSpace() {
}

// Java method: CodeConsumer.breakAfterBlockFor(com.google.javascript.rhino.Node n, boolean statementContext)
func (this_ *CodeConsumerBase) BreakAfterBlockFor(n pkg_rhino.Node, statementContext bool) bool {
	_ = n
	_ = statementContext
	return statementContext
}

// Java method: CodeConsumer.endFile()
func (this_ *CodeConsumerBase) EndFile() {
}
