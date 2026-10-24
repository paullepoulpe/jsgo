// Port of com/google/javascript/jscomp/instrumentation/FileInstrumentationData.java @ closure-compiler v20261021

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

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// FileInstrumentationData holds instrumentation details related to a file,
// namely, the filename, the array name used in instrumentation, and the
// lines which were instrumented (in encoded form).
type FileInstrumentationData struct {
	instrumentedBits *jutil.BitSet // Instrumented lines, a bit per line
	arrayName        string
	fileName         string

	// branchPresent denotes the lines containing at least one branch
	// bit[i] set to 1 denotes the (i+1)-th line containing at least one
	// branch bit[i] bit set to 0 denotes the (i+1)-th line containing no
	// branches
	branchPresent *jutil.BitSet

	// Number of branches in line, the index is zero-based.
	branchesInLine *jutil.LinkedHashMap[int, int]

	// Map of (line-no, branch-idx) to nodes of blocks of conditional
	// branches. The two indices are all zero-based.
	//
	// For example, if the source code has an 'if' statement on line 3 with
	// an if and else branch, the two branches will be indexed as (2, 0) and
	// (2, 1).
	branchNodes *jutil.LinkedHashMap[FileInstrumentationData_BranchIndexPair, rhino.Node]
}

// FileInstrumentationData_BranchIndexPair is a (line, branch) pair.
//
// Go-only: language-adapter — its equals/hashCode compare the two ints, so
// it is a Go value type (usable as a map key with the same equality).
type FileInstrumentationData_BranchIndexPair struct {
	line   int
	branch int
}

// Java: FileInstrumentationData.BranchIndexPair.getLine()
func (p FileInstrumentationData_BranchIndexPair) GetLine() int {
	return p.line
}

// Java: FileInstrumentationData.BranchIndexPair.getBranch()
func (p FileInstrumentationData_BranchIndexPair) GetBranch() int {
	return p.branch
}

// Java: FileInstrumentationData.BranchIndexPair(int, int)
func NewFileInstrumentationData_BranchIndexPair(line int, branch int) FileInstrumentationData_BranchIndexPair {
	return FileInstrumentationData_BranchIndexPair{line: line, branch: branch}
}

// Java: FileInstrumentationData.BranchIndexPair.of(int, int)
func FileInstrumentationData_BranchIndexPair_Of(line int, branch int) FileInstrumentationData_BranchIndexPair {
	return NewFileInstrumentationData_BranchIndexPair(line, branch)
}

// Java: FileInstrumentationData.BranchIndexPair.hashCode()
func (p FileInstrumentationData_BranchIndexPair) HashCode() int32 {
	return 31*int32(p.line) + int32(p.branch)
}

// Java: FileInstrumentationData.BranchIndexPair.equals(Object)
func (p FileInstrumentationData_BranchIndexPair) Equals(object any) bool {
	if that, ok := object.(FileInstrumentationData_BranchIndexPair); ok {
		return p.GetLine() == that.GetLine() && p.GetBranch() == that.GetBranch()
	}
	return false
}

// Java: FileInstrumentationData(String, String)
func NewFileInstrumentationData(fileName string, arrayName string) *FileInstrumentationData {
	return &FileInstrumentationData{
		fileName:         fileName,
		arrayName:        arrayName,
		instrumentedBits: jutil.NewBitSet(),
		branchPresent:    jutil.NewBitSet(),
		branchesInLine:   jutil.NewLinkedHashMap[int, int](),
		branchNodes:      jutil.NewLinkedHashMap[FileInstrumentationData_BranchIndexPair, rhino.Node](),
	}
}

// Java: FileInstrumentationData.getArrayName()
func (d *FileInstrumentationData) GetArrayName() string {
	return d.arrayName
}

// Java: FileInstrumentationData.getFileName()
func (d *FileInstrumentationData) GetFileName() string {
	return d.fileName
}

// Java: FileInstrumentationData.maxInstrumentedLine()
func (d *FileInstrumentationData) MaxInstrumentedLine() int {
	return d.instrumentedBits.Length()
}

// Java: FileInstrumentationData.maxBranchPresentLine()
func (d *FileInstrumentationData) MaxBranchPresentLine() int {
	return d.branchPresent.Length()
}

// PutBranchNode stores a node to be instrumented later for branch
// coverage.
//
// lineNumber is the 1-based line number, branchNumber the 1-based branch
// number, block the node of the conditional block.
//
// Java: FileInstrumentationData.putBranchNode(int, int, Node)
func (d *FileInstrumentationData) PutBranchNode(lineNumber int, branchNumber int, block rhino.Node) {
	jutil.CheckArgumentMsg1(
		lineNumber > 0, "Expected non-zero positive integer as line number: %s", lineNumber)
	jutil.CheckArgumentMsg1(
		branchNumber > 0, "Expected non-zero positive integer as branch number: %s", branchNumber)

	d.branchNodes.Put(FileInstrumentationData_BranchIndexPair_Of(lineNumber-1, branchNumber-1), block)
}

// GetBranchNode gets the block node to be instrumented for branch
// coverage.
//
// lineNumber is the 1-based line number, branchNumber the 1-based branch
// number. Returns the node of the conditional block.
//
// Java: FileInstrumentationData.getBranchNode(int, int)
func (d *FileInstrumentationData) GetBranchNode(lineNumber int, branchNumber int) rhino.Node {
	jutil.CheckArgumentMsg1(
		lineNumber > 0, "Expected non-zero positive integer as line number: %s", lineNumber)
	jutil.CheckArgumentMsg1(
		branchNumber > 0, "Expected non-zero positive integer as branch number: %s", branchNumber)

	return d.branchNodes.Get(FileInstrumentationData_BranchIndexPair_Of(lineNumber-1, branchNumber-1))
}

// fileInstrumentationData_getHexString returns a byte-wise hex string
// representation of the BitField from MSB (Most Significant Byte) to LSB
// (Least Significant Byte). Eg. Single byte: a setting of "0001 1111",
// returns "1f" Eg. Multiple bytes: a setting of "0000 0010 0001 1111",
// returns "1f02"
//
// Java: FileInstrumentationData.getHexString(BitSet)
func fileInstrumentationData_getHexString(bitSet *jutil.BitSet) string {
	builder := jutil.NewStringBuilder()

	// Build the hex string.
	for _, byteEntry := range bitSet.ToByteArray() {
		// Java bytes are signed, but we want the value as if it were
		// unsigned.
		value := int(byteEntry)
		// Integer.toHexString, padded to be two characters (if it isn't
		// already).
		builder.Append(fmt.Sprintf("%02x", value))
	}
	return builder.String()
}

// GetInstrumentedLinesAsHexString gets a hex string representation of the
// instrumentedBits bit vector.
//
// Java: FileInstrumentationData.getInstrumentedLinesAsHexString()
func (d *FileInstrumentationData) GetInstrumentedLinesAsHexString() string {
	return fileInstrumentationData_getHexString(d.instrumentedBits)
}

// GetBranchPresentAsHexString gets a hex string representation of the
// branchPresent bit vector.
//
// Java: FileInstrumentationData.getBranchPresentAsHexString()
func (d *FileInstrumentationData) GetBranchPresentAsHexString() string {
	return fileInstrumentationData_getHexString(d.branchPresent)
}

// SetLineAsInstrumented marks given 1-based line number as instrumented.
// Zero, Negative numbers are not allowed.
//
// Java: FileInstrumentationData.setLineAsInstrumented(int)
func (d *FileInstrumentationData) SetLineAsInstrumented(lineNumber int) {
	jutil.CheckArgumentMsg1(
		lineNumber > 0, "Expected non-zero positive integer as line number: %s", lineNumber)

	// Map the 1-based line number to 0-based bit position
	d.instrumentedBits.Set(lineNumber - 1)
}

// SetBranchPresent marks a given 1-based line number has branch presented.
//
// Java: FileInstrumentationData.setBranchPresent(int)
func (d *FileInstrumentationData) SetBranchPresent(lineNumber int) {
	jutil.CheckArgumentMsg1(
		lineNumber > 0, "Expected non-zero positive integer as line number: %s", lineNumber)

	// Map the 1-based line number to 0-based bit position
	d.branchPresent.Set(lineNumber - 1)
}

// AddBranches adds a number of branches to a line.
//
// Java: FileInstrumentationData.addBranches(int, int)
func (d *FileInstrumentationData) AddBranches(lineNumber int, numberOfBranches int) {
	lineIdx := lineNumber - 1
	currentValue, ok := d.branchesInLine.Lookup(lineIdx)
	if !ok {
		d.branchesInLine.Put(lineIdx, numberOfBranches)
	} else {
		d.branchesInLine.Put(lineIdx, currentValue+numberOfBranches)
	}
}

// GetNumBranches gets the number of branches on a line (1-based line
// number).
//
// Java: FileInstrumentationData.getNumBranches(int)
func (d *FileInstrumentationData) GetNumBranches(lineNumber int) int {
	numBranches, ok := d.branchesInLine.Lookup(lineNumber - 1)
	if !ok {
		return 0
	} else {
		return numBranches
	}
}
