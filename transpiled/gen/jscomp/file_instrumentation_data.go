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
// Port of com/google/javascript/jscomp/instrumentation/FileInstrumentationData.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.instrumentation.FileInstrumentationData
type FileInstrumentationData struct {
	instrumentedBits *j2gort.Ext_java_util_BitSet                                                      `gson:"instrumentedBits"`
	arrayName        string                                                                            `gson:"arrayName"`
	fileName         string                                                                            `gson:"fileName"`
	branchPresent    *j2gort.Ext_java_util_BitSet                                                      `gson:"branchPresent"`
	branchesInLine   j2gort.OrderedMap[int32, int32]                                                   `gson:"branchesInLine"`
	branchNodes      j2gort.CustomOrderedMap[*FileInstrumentationData_BranchIndexPair, pkg_rhino.Node] `gson:"branchNodes"`
	idHash_          atomic.Int32
}

func (this_ *FileInstrumentationData) BindSelf_(s any) {
	_ = s
}
func (this_ *FileInstrumentationData) IsJavaNil_() bool { return this_ == nil }
func (*FileInstrumentationData) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.FileInstrumentationData", "FileInstrumentationData"
}
func (this_ *FileInstrumentationData) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FileInstrumentationData) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.FileInstrumentationData", "com.google.javascript.jscomp.instrumentation.FileInstrumentationData", func(v_ any) bool {
		return j2gort.InstanceOf[*FileInstrumentationData](v_)
	})
}

func (this_ *FileInstrumentationData) initFields_FileInstrumentationData() {
	_ = this_
	this_.arrayName = j2gort.NullString
	this_.fileName = j2gort.NullString
}

// Java constructor: FileInstrumentationData(java.lang.String fileName, java.lang.String arrayName)
func NewFileInstrumentationData(fileName string, arrayName string) *FileInstrumentationData {
	this_ := &FileInstrumentationData{}
	InitFileInstrumentationData(this_, fileName, arrayName)
	return this_
}

func InitFileInstrumentationData(this_ *FileInstrumentationData, fileName string, arrayName string) {
	_ = fileName
	_ = arrayName
	this_.initFields_FileInstrumentationData()
	this_.fileName = fileName
	this_.arrayName = arrayName
	this_.instrumentedBits = j2gort.NewExt_java_util_BitSet()
	this_.branchPresent = j2gort.NewExt_java_util_BitSet()
	this_.branchesInLine = j2gort.NewOrderedMap[int32, int32]()
	this_.branchNodes = j2gort.NewCustomOrderedMap[*FileInstrumentationData_BranchIndexPair,

		// Java method: FileInstrumentationData.getArrayName()
		pkg_rhino.Node]()
}

func (this_ *FileInstrumentationData) GetArrayName() string {
	return this_.arrayName
}

// Java method: FileInstrumentationData.getFileName()
func (this_ *FileInstrumentationData) GetFileName() string {
	return this_.fileName
}

// Java method: FileInstrumentationData.maxInstrumentedLine()
func (this_ *FileInstrumentationData) MaxInstrumentedLine() int32 {
	return this_.instrumentedBits.Length()
}

// Java method: FileInstrumentationData.maxBranchPresentLine()
func (this_ *FileInstrumentationData) MaxBranchPresentLine() int32 {
	return this_.branchPresent.Length()
}

// Java method: FileInstrumentationData.putBranchNode(int lineNumber, int branchNumber, com.google.javascript.rhino.Node block)
func (this_ *FileInstrumentationData) PutBranchNode(lineNumber int32, branchNumber int32, block pkg_rhino.Node) {
	_ = lineNumber
	_ = branchNumber
	_ = block
	if !(lineNumber > int32(0)) {
		j2gort.CheckArgumentFmt(false, "Expected non-zero positive integer as line number: %s", j2gort.ValueOfInt(lineNumber))
	}
	if !(branchNumber > int32(0)) {
		j2gort.CheckArgumentFmt(false, "Expected non-zero positive integer as branch number: %s", j2gort.ValueOfInt(branchNumber))
	}
	this_.branchNodes.Put(FileInstrumentationData_BranchIndexPair_Of((lineNumber-int32(1)), (branchNumber-int32(1))), block)
}

// Java method: FileInstrumentationData.getBranchNode(int lineNumber, int branchNumber)
func (this_ *FileInstrumentationData) GetBranchNode(lineNumber int32, branchNumber int32) pkg_rhino.Node {
	_ = lineNumber
	_ = branchNumber
	if !(lineNumber > int32(0)) {
		j2gort.CheckArgumentFmt(false, "Expected non-zero positive integer as line number: %s", j2gort.ValueOfInt(lineNumber))
	}
	if !(branchNumber > int32(0)) {
		j2gort.CheckArgumentFmt(false, "Expected non-zero positive integer as branch number: %s", j2gort.ValueOfInt(branchNumber))
	}
	return this_.branchNodes.Get(FileInstrumentationData_BranchIndexPair_Of((lineNumber - int32(1)), (branchNumber - int32(1))))
}

// Java static method: FileInstrumentationData.getHexString(java.util.BitSet bitSet)
func FileInstrumentationData_GetHexString(bitSet *j2gort.Ext_java_util_BitSet) string {
	_ = bitSet
	var builder j2gort.LocalStringBuilder

	for _, el_1 := range j2gort.CheckNotNullSlice(bitSet.ToByteArray()) {
		byteEntry := el_1
		_ = byteEntry
		var value int32 = j2gort.BoxByte_ToUnsignedInt(byteEntry)
		_ = value
		var hexString string = javalang.IntegerToHexString(value)
		_ = hexString
		hexString = j2gort.Ext_com_google_common_base_Strings_PadStart(hexString, int32(2), uint16(48))
		builder.AppendStr(hexString)
	}
	return builder.String()
}

// Java method: FileInstrumentationData.getInstrumentedLinesAsHexString()
func (this_ *FileInstrumentationData) GetInstrumentedLinesAsHexString() string {
	return FileInstrumentationData_GetHexString(this_.instrumentedBits)
}

// Java method: FileInstrumentationData.getBranchPresentAsHexString()
func (this_ *FileInstrumentationData) GetBranchPresentAsHexString() string {
	return FileInstrumentationData_GetHexString(this_.branchPresent)
}

// Java method: FileInstrumentationData.setLineAsInstrumented(int lineNumber)
func (this_ *FileInstrumentationData) SetLineAsInstrumented(lineNumber int32) {
	_ = lineNumber
	if !(lineNumber > int32(0)) {
		j2gort.CheckArgumentFmt(false, "Expected non-zero positive integer as line number: %s", j2gort.ValueOfInt(lineNumber))
	}
	this_.instrumentedBits.Set_1((lineNumber - int32(1)))
}

// Java method: FileInstrumentationData.setBranchPresent(int lineNumber)
func (this_ *FileInstrumentationData) SetBranchPresent(lineNumber int32) {
	_ = lineNumber
	if !(lineNumber > int32(0)) {
		j2gort.CheckArgumentFmt(false, "Expected non-zero positive integer as line number: %s", j2gort.ValueOfInt(lineNumber))
	}
	this_.branchPresent.Set_1((lineNumber - int32(1)))
}

// Java method: FileInstrumentationData.addBranches(int lineNumber, int numberOfBranches)
func (this_ *FileInstrumentationData) AddBranches(lineNumber int32, numberOfBranches int32) {
	_ = lineNumber
	_ = numberOfBranches
	var lineIdx int32 = (lineNumber - int32(1))
	_ = lineIdx
	currentValue, currentValueOk_ := this_.branchesInLine.GetOk(lineIdx)
	_ = currentValueOk_
	_ = currentValue
	if !currentValueOk_ {
		this_.branchesInLine.Put(lineIdx, numberOfBranches)
	} else {
		this_.branchesInLine.Put(lineIdx, (func() int32 {
			if !currentValueOk_ {
				j2gort.ThrowNPE()
			}
			return currentValue

			// Java method: FileInstrumentationData.getNumBranches(int lineNumber)
		}() + numberOfBranches))
	}
}

func (this_ *FileInstrumentationData) GetNumBranches(lineNumber int32) int32 {
	_ = lineNumber
	numBranches, numBranchesOk_ := this_.branchesInLine.GetOk(lineNumber - int32(1))
	_ = numBranchesOk_
	_ = numBranches
	if !numBranchesOk_ {
		return int32(0)
	} else {
		{

			if !numBranchesOk_ {
				j2gort.ThrowNPE(

				// Java class: com.google.javascript.jscomp.instrumentation.FileInstrumentationData.BranchIndexPair
				)
			}
			return numBranches
		}

	}
}

type FileInstrumentationData_BranchIndexPair struct {
	line   int32 `gson:"line"`
	branch int32 `gson:"branch"`
}

func (this_ *FileInstrumentationData_BranchIndexPair) BindSelf_(s any) {
	_ = s
}
func (this_ *FileInstrumentationData_BranchIndexPair) IsJavaNil_() bool { return this_ == nil }
func (*FileInstrumentationData_BranchIndexPair) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.FileInstrumentationData$BranchIndexPair", "BranchIndexPair"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.FileInstrumentationData.BranchIndexPair", "com.google.javascript.jscomp.instrumentation.FileInstrumentationData$BranchIndexPair", func(v_ any) bool {
		return j2gort.InstanceOf[*FileInstrumentationData_BranchIndexPair](v_)
	})
}

func (this_ *FileInstrumentationData_BranchIndexPair) initFields_FileInstrumentationData_BranchIndexPair() {
	_ = this_
}

// Java constructor: BranchIndexPair(int line, int branch)
func NewFileInstrumentationData_BranchIndexPair(line int32, branch int32) *FileInstrumentationData_BranchIndexPair {
	this_ := &FileInstrumentationData_BranchIndexPair{}
	InitFileInstrumentationData_BranchIndexPair(this_, line, branch)
	return this_
}

func InitFileInstrumentationData_BranchIndexPair(this_ *FileInstrumentationData_BranchIndexPair, line int32, branch int32) {
	_ = line
	_ = branch
	this_.initFields_FileInstrumentationData_BranchIndexPair()
	this_.line = line
	this_.branch = branch
}

// Java method: BranchIndexPair.getLine()
func (this_ *FileInstrumentationData_BranchIndexPair) GetLine() int32 {
	return this_.line
}

// Java method: BranchIndexPair.getBranch()
func (this_ *FileInstrumentationData_BranchIndexPair) GetBranch() int32 {
	return this_.branch
}

// Java static method: BranchIndexPair.of(int line, int branch)
func FileInstrumentationData_BranchIndexPair_Of(line int32, branch int32) *FileInstrumentationData_BranchIndexPair {
	_ = line
	_ = branch
	return NewFileInstrumentationData_BranchIndexPair(line, branch)
}

// Java method: BranchIndexPair.hashCode()
func (this_ *FileInstrumentationData_BranchIndexPair) HashCode() int32 {
	return ((int32(31) * this_.line) + this_.branch)
}

// Java method: BranchIndexPair.equals(java.lang.Object object)
func (this_ *FileInstrumentationData_BranchIndexPair) Equals(arg_2 any) bool {
	object := arg_2
	_ = object
	var that *FileInstrumentationData_BranchIndexPair
	_ = that
	var ok_3 bool
	that, ok_3 = j2gort.TryCast[*FileInstrumentationData_BranchIndexPair](object)
	if ok_3 {
		return ((this_.GetLine() == that.GetLine()) && (this_.GetBranch() == that.GetBranch()))
	}
	return false
}
