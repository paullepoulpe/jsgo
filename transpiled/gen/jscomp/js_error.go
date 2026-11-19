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
// Port of com/google/javascript/jscomp/JSError.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.JSError
type JSError struct {
	type_        *DiagnosticType                                      `gson:"type"`
	description  string                                               `gson:"description"`
	sourceName   string                                               `gson:"sourceName"`
	defaultLevel *CheckLevel                                          `gson:"defaultLevel"`
	requirement  *j2gort.Ext_com_google_javascript_jscomp_Requirement `gson:"requirement"`
	lineno       int32                                                `gson:"lineno"`
	charno       int32                                                `gson:"charno"`
	length       int32                                                `gson:"length"`
	node         pkg_rhino.Node                                       `gson:"node"`
}

func (this_ *JSError) BindSelf_(s any) {
	_ = s
}
func (this_ *JSError) IsJavaNil_() bool { return this_ == nil }
func (*JSError) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSError", "JSError"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JSError", "com.google.javascript.jscomp.JSError", func(v_ any) bool {
		return j2gort.InstanceOf[*JSError](v_)
	})
}
func (*JSError) Implements_j2gort_Ext_java_lang_Record() {}

func StaticGet_JSError_DEFAULT_LINENO() int32 { return int32(-1) }
func StaticGet_JSError_DEFAULT_CHARNO() int32 { return int32(-1) }
func StaticGet_JSError_DEFAULT_LENGTH() int32 { return int32(0) }

var static_JSError_DEFAULT_SOURCENAME string

func StaticGet_JSError_DEFAULT_SOURCENAME() string {
	clinit_JSError()
	return static_JSError_DEFAULT_SOURCENAME
}
func StaticSet_JSError_DEFAULT_SOURCENAME(v string) {
	clinit_JSError()
	static_JSError_DEFAULT_SOURCENAME = v
}

var static_JSError_DEFAULT_NODE pkg_rhino.Node

func StaticGet_JSError_DEFAULT_NODE() pkg_rhino.Node {
	clinit_JSError()
	return static_JSError_DEFAULT_NODE
}
func StaticSet_JSError_DEFAULT_NODE(v pkg_rhino.Node) {
	clinit_JSError()
	static_JSError_DEFAULT_NODE = v
}

var static_JSError_DEFAULT_REQUIREMENT *j2gort.Ext_com_google_javascript_jscomp_Requirement

func StaticGet_JSError_DEFAULT_REQUIREMENT() *j2gort.Ext_com_google_javascript_jscomp_Requirement {
	clinit_JSError()
	return static_JSError_DEFAULT_REQUIREMENT
}
func StaticSet_JSError_DEFAULT_REQUIREMENT(v *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	clinit_JSError()
	static_JSError_DEFAULT_REQUIREMENT = v
}

var clinitDone_JSError uint32
var clinitState_JSError j2gort.ClinitState

func clinit_JSError() {
	if atomic.LoadUint32(&clinitDone_JSError) != 0 {
		return
	}
	clinitSlow_JSError()
}

func clinitSlow_JSError() {
	if !j2gort.EnterClinit(&clinitState_JSError, &clinitDone_JSError) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSError, &clinitDone_JSError)
	static_JSError_DEFAULT_SOURCENAME = j2gort.NullString
	static_JSError_DEFAULT_NODE = pkg_rhino.Node(0)
	static_JSError_DEFAULT_REQUIREMENT = nil
}

func (this_ *JSError) initFields_JSError() {
	_ = this_
}

// Java constructor: JSError(com.google.javascript.jscomp.DiagnosticType type, java.lang.String description, java.lang.@org.jspecify.annotations.Nullable String sourceName, int lineno, int charno, int length, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node node, com.google.javascript.jscomp.CheckLevel defaultLevel, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Requirement requirement)
func NewJSError(type_ *DiagnosticType, description string, sourceName string, lineno int32, charno int32, length int32, node pkg_rhino.Node, defaultLevel *CheckLevel, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *JSError {
	clinit_JSError()
	this_ := &JSError{}
	InitJSError(this_, type_, description, sourceName, lineno, charno, length, node, defaultLevel, requirement)
	return this_
}

func InitJSError(this_ *JSError, type_ *DiagnosticType, description string, sourceName string, lineno int32, charno int32, length int32, node pkg_rhino.Node, defaultLevel *CheckLevel, requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) {
	_ = type_
	_ = description
	_ = sourceName
	_ = lineno
	_ = charno
	_ = length
	_ = node
	_ = defaultLevel
	_ = requirement
	this_.initFields_JSError()
	j2gort.CheckNotNullPtr(type_)
	j2gort.CheckNotNullString(description)
	j2gort.CheckNotNullPtr(defaultLevel)
	this_.type_ = type_
	this_.description = description
	this_.sourceName = sourceName
	this_.lineno = lineno
	this_.charno = charno
	this_.length = length
	this_.node = node
	this_.defaultLevel = defaultLevel
	this_.requirement = requirement
}

// Java method: JSError.getType()
func (this_ *JSError) GetType() *DiagnosticType {
	return this_.Type()
}

// Java method: JSError.getDescription()
func (this_ *JSError) GetDescription() string {
	return this_.Description()
}

// Java method: JSError.getSourceName()
func (this_ *JSError) GetSourceName() string {
	return this_.SourceName()
}

// Java method: JSError.getLineno()
func (this_ *JSError) GetLineno() int32 {
	return this_.Lineno()
}

// Java method: JSError.getCharno()
func (this_ *JSError) GetCharno() int32 {
	return this_.Charno()
}

// Java method: JSError.getLength()
func (this_ *JSError) GetLength() int32 {
	return this_.Length()
}

// Java method: JSError.getNode()
func (this_ *JSError) GetNode() pkg_rhino.Node {
	return this_.Node()
}

// Java method: JSError.getDefaultLevel()
func (this_ *JSError) GetDefaultLevel() *CheckLevel {
	return this_.DefaultLevel()
}

// Java method: JSError.getRequirement()
func (this_ *JSError) GetRequirement() *j2gort.Ext_com_google_javascript_jscomp_Requirement {
	return this_.Requirement()
}

// Java static method: JSError.make(com.google.javascript.jscomp.DiagnosticType type, java.lang.String[] arguments)
func JSError_Make_2(type_ *DiagnosticType, arguments []string) *JSError {
	clinit_JSError()
	_ = type_
	_ = arguments
	return JSError_S_Builder(type_, arguments).Build()
}

// Java static method: JSError.make(java.lang.String sourceName, int lineno, int charno, com.google.javascript.jscomp.DiagnosticType type, java.lang.String[] arguments)
func JSError_Make_5(sourceName string, lineno int32, charno int32, type_ *DiagnosticType, arguments []string) *JSError {
	clinit_JSError()
	_ = sourceName
	_ = lineno
	_ = charno
	_ = type_
	_ = arguments
	return JSError_S_Builder(type_, arguments).SetSourceLocation(sourceName, lineno, charno).Build()
}

// Java static method: JSError.make(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DiagnosticType type, java.lang.String[] arguments)
func JSError_Make_3(n pkg_rhino.Node, type_ *DiagnosticType, arguments []string) *JSError {
	clinit_JSError()
	_ = n
	_ = type_
	_ = arguments
	return JSError_S_Builder(type_, arguments).SetNode(n).Build()
}

// Java static method: JSError.make(com.google.javascript.rhino.Node start, com.google.javascript.rhino.Node end, com.google.javascript.jscomp.DiagnosticType type, java.lang.String[] arguments)
func JSError_Make_4_Node_Node_DiagnosticType_slice_string(start pkg_rhino.Node, end pkg_rhino.Node, type_ *DiagnosticType, arguments []string) *JSError {
	clinit_JSError()
	_ = start
	_ = end
	_ = type_
	_ = arguments
	return JSError_S_Builder(type_, arguments).SetNodeRange(start, end).Build()
}

// Java static method: JSError.make(com.google.javascript.jscomp.Requirement requirement, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DiagnosticType type, java.lang.String[] arguments)
func JSError_Make_4_Ext_com_google_javascript_jscomp_Requirement_Node_DiagnosticType_slice_string(requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement, n pkg_rhino.Node, type_ *DiagnosticType, arguments []string) *JSError {
	clinit_JSError()
	_ = requirement
	_ = n
	_ = type_
	_ = arguments
	return JSError_S_Builder(type_, arguments).SetNode(n).SetRequirement(requirement).Build()
}

// Java static method: JSError.builder(com.google.javascript.jscomp.DiagnosticType type, java.lang.String[] arguments)
func JSError_S_Builder(type_ *DiagnosticType, arguments []string) *JSError_Builder {
	clinit_JSError()
	_ = type_
	_ = arguments
	return NewJSError_Builder(type_, arguments)
}

// Java method: JSError.toString()
func (this_ *JSError) String() string {
	var sourceName string
	_ = sourceName
	if j2gort.Ext_com_google_common_base_Strings_EmptyToNull(this_.SourceName()) != j2gort.NullString {
		sourceName = this_.SourceName()
	} else {
		sourceName = "(unknown source)"
	}
	var lineno string
	_ = lineno
	if this_.Lineno() != int32(-1) {
		lineno = strconv.FormatInt(int64(this_.Lineno()), 10)
	} else {
		lineno = "(unknown line)"
	}
	var charno string
	_ = charno
	if this_.Charno() != int32(-1) {
		charno = strconv.FormatInt(int64(this_.Charno()), 10)
	} else {
		charno = "(unknown column)"
	}
	return (j2gort.StrValString(this_.Type().key) + ". " + j2gort.StrValString(this_.Description()) + " at " + j2gort.StrValString(sourceName) + " line " + j2gort.StrValString(lineno) + " : " + j2gort.StrValString(charno))
}

// Java method: JSError.format(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.MessageFormatter formatter)
func (this_ *JSError) Format(level *CheckLevel, formatter MessageFormatter) string {
	_ = level
	_ = formatter
	switch j2gort.CheckNotNullPtr(level).Ordinal() {
	case int32(0) /* CheckLevel.ERROR */ :
		return formatter.FormatError(this_)
	case int32(1) /* CheckLevel.WARNING */ :
		return formatter.FormatWarning(this_)
	default:
		return j2gort.NullString
	}
	panic("unreachable switch expression")
}

// Java method: JSError.getNodeSourceOffset()
func (this_ *JSError) GetNodeSourceOffset() int32 {
	if this_.Node() != 0 {
		return this_.Node().GetSourceOffset()
	} else {
		return int32(-1)
	}
}

// Java method: JSError.getLineNumber()
func (this_ *JSError) GetLineNumber() int32 {
	return this_.Lineno()
}

func (this_ *JSError) Type() *DiagnosticType     { return this_.type_ }
func (this_ *JSError) Description() string       { return this_.description }
func (this_ *JSError) SourceName() string        { return this_.sourceName }
func (this_ *JSError) Lineno() int32             { return this_.lineno }
func (this_ *JSError) Charno() int32             { return this_.charno }
func (this_ *JSError) Length() int32             { return this_.length }
func (this_ *JSError) Node() pkg_rhino.Node      { return this_.node }
func (this_ *JSError) DefaultLevel() *CheckLevel { return this_.defaultLevel }
func (this_ *JSError) Requirement() *j2gort.Ext_com_google_javascript_jscomp_Requirement {
	return this_.requirement
}
func (this_ *JSError) Equals(other_ any) bool {
	o_, ok_ := other_.(*JSError)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.type_), j2gort.PtrToAny(o_.type_)) {
		return false
	}
	if this_.description != o_.description {
		return false
	}
	if this_.sourceName != o_.sourceName {
		return false
	}
	if this_.lineno != o_.lineno {
		return false
	}
	if this_.charno != o_.charno {
		return false
	}
	if this_.length != o_.length {
		return false
	}
	if this_.node != o_.node {
		return false
	}
	if this_.defaultLevel != o_.defaultLevel {
		return false
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.requirement), j2gort.PtrToAny(o_.requirement)) {
		return false
	}
	return true
}
func (this_ *JSError) HashCode() int32 {
	var h_ int32 = 0
	{
		hp_1 := this_.type_
		h_ = 31 * h_
		if hp_1 != nil {
			h_ += hp_1.HashCode()
		}
	}
	{
		hs_2 := this_.description
		h_ = 31 * h_
		if hs_2 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_2)
		}
	}
	{
		hs_3 := this_.sourceName
		h_ = 31 * h_
		if hs_3 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_3)
		}
	}
	h_ = 31*h_ + (this_.lineno)
	h_ = 31*h_ + (this_.charno)
	h_ = 31*h_ + (this_.length)
	h_ = 31*h_ + int32(this_.node)
	{
		hp_4 := this_.defaultLevel
		h_ = 31 * h_
		if hp_4 != nil {
			h_ += hp_4.HashCode()
		}
	}
	h_ = 31*h_ + j2gort.ObjectHash(j2gort.PtrToAny(this_.requirement))
	return h_
}

// Java class: com.google.javascript.jscomp.JSError.Builder
type JSError_Builder struct {
	type_       *DiagnosticType                                      `gson:"type"`
	args        []string                                             `gson:"args"`
	level       *CheckLevel                                          `gson:"level"`
	sourceName  string                                               `gson:"sourceName"`
	requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement `gson:"requirement"`
	n           pkg_rhino.Node                                       `gson:"n"`
	lineno      int32                                                `gson:"lineno"`
	charno      int32                                                `gson:"charno"`
	length      int32                                                `gson:"length"`
	idHash_     atomic.Int32
}

func (this_ *JSError_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *JSError_Builder) IsJavaNil_() bool { return this_ == nil }
func (*JSError_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSError$Builder", "Builder"
}
func (this_ *JSError_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSError_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JSError.Builder", "com.google.javascript.jscomp.JSError$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*JSError_Builder](v_)
	})
}

func (this_ *JSError_Builder) initFields_JSError_Builder() {
	_ = this_
	this_.n = StaticGet_JSError_DEFAULT_NODE()
	this_.sourceName = StaticGet_JSError_DEFAULT_SOURCENAME()
	this_.lineno = int32(-1)
	this_.charno = int32(-1)

	this_.requirement = StaticGet_JSError_DEFAULT_REQUIREMENT()
}

// Java constructor: Builder(com.google.javascript.jscomp.DiagnosticType type, java.lang.String[] args)
func NewJSError_Builder(type_ *DiagnosticType, args []string) *JSError_Builder {
	this_ := &JSError_Builder{}
	InitJSError_Builder(this_, type_, args)
	return this_
}

func InitJSError_Builder(this_ *JSError_Builder, type_ *DiagnosticType, args []string) {
	_ = type_
	_ = args
	this_.initFields_JSError_Builder()
	this_.type_ = type_
	this_.args = args
	this_.level = type_.level
}

// Java method: Builder.setNode(com.google.javascript.rhino.Node n)
func (this_ *JSError_Builder) SetNode(n pkg_rhino.Node) *JSError_Builder {
	_ = n
	if !(StaticGet_JSError_DEFAULT_SOURCENAME() == this_.sourceName) {
		j2gort.CheckStateStr(false, "Cannot provide a Node when there's already a source name")
	}
	this_.n = n
	this_.sourceName = n.GetSourceFileName()
	this_.lineno = n.GetLineno()
	this_.charno = n.GetCharno()
	this_.length = n.GetLength()
	return this_
}

// Java method: Builder.setNodeRange(com.google.javascript.rhino.Node start, com.google.javascript.rhino.Node end)
func (this_ *JSError_Builder) SetNodeRange(start pkg_rhino.Node, end pkg_rhino.Node) *JSError_Builder {
	_ = start
	_ = end
	if !(StaticGet_JSError_DEFAULT_SOURCENAME() == this_.sourceName) {
		j2gort.CheckStateStr(false, "Cannot provide a Node when there's already a source name")
	}
	this_.n = start
	this_.sourceName = this_.n.GetSourceFileName()
	this_.lineno = this_.n.GetLineno()
	this_.charno = this_.n.GetCharno()
	var endOffset int32 = end.GetSourceOffset()
	_ = endOffset
	var startOffset int32 = start.GetSourceOffset()
	_ = startOffset
	if (endOffset != int32(-1)) && (startOffset != int32(-1)) {
		this_.length = ((end.GetSourceOffset() + end.GetLength()) - start.GetSourceOffset())
	}
	return this_
}

// Java method: Builder.setLevel(com.google.javascript.jscomp.CheckLevel level)
func (this_ *JSError_Builder) SetLevel(level *CheckLevel) *JSError_Builder {
	_ = level
	this_.level = j2gort.CheckNotNullPtr(level)
	return this_
}

// Java method: Builder.setRequirement(com.google.javascript.jscomp.Requirement requirement)
func (this_ *JSError_Builder) SetRequirement(requirement *j2gort.Ext_com_google_javascript_jscomp_Requirement) *JSError_Builder {
	_ = requirement
	this_.requirement = j2gort.CheckNotNullPtr(requirement)
	return this_
}

// Java method: Builder.setSourceLocation(java.lang.String sourceName, int lineno, int charno)
func (this_ *JSError_Builder) SetSourceLocation(sourceName string, lineno int32, charno int32) *JSError_Builder {
	_ = sourceName
	_ = lineno
	_ = charno
	if !(this_.n == StaticGet_JSError_DEFAULT_NODE()) {
		j2gort.CheckStateStr(false, "Cannot provide a source location when there is already a Node")
	}
	this_.sourceName = sourceName
	this_.lineno = lineno
	this_.charno = charno
	return this_
}

// Java method: Builder.build()
func (this_ *JSError_Builder) Build() *JSError {
	return NewJSError(this_.type_, this_.type_.Format(this_.args), this_.sourceName, this_.lineno, this_.charno, this_.length, this_.n, this_.level, this_.requirement)
}
