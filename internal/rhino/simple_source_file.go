// Port of com/google/javascript/rhino/SimpleSourceFile.java @ closure-compiler v20260415

package rhino

import (
	"math"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SimpleSourceFile is a simple implementation of StaticSourceFile for
// testing.
type SimpleSourceFile struct {
	name string
	kind StaticSourceFile_SourceKind
}

// Java: SimpleSourceFile(String, SourceKind)
func NewSimpleSourceFile(name string, kind StaticSourceFile_SourceKind) *SimpleSourceFile {
	return &SimpleSourceFile{name: name, kind: kind}
}

// Java: SimpleSourceFile.getName()
func (f *SimpleSourceFile) GetName() string {
	return f.name
}

// Java: SimpleSourceFile.getKind()
func (f *SimpleSourceFile) GetKind() StaticSourceFile_SourceKind {
	return f.kind
}

// Java: SimpleSourceFile.isClosureUnawareCode()
func (f *SimpleSourceFile) IsClosureUnawareCode() bool {
	return false
}

// Java: SimpleSourceFile.getColumnOfOffset(int)
func (f *SimpleSourceFile) GetColumnOfOffset(offset int) int {
	return 0
}

// Java: SimpleSourceFile.getLineOfOffset(int)
func (f *SimpleSourceFile) GetLineOfOffset(offset int) int {
	return 1
}

// Java: SimpleSourceFile.getLineOffset(int)
func (f *SimpleSourceFile) GetLineOffset(line int) int {
	if line < 1 {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Should not call getLineOffset with line number " + strconv.Itoa(line)))
	}
	return math.MinInt32
}

// Java: SimpleSourceFile.toString()
func (f *SimpleSourceFile) String() string {
	return f.name
}

// Inherited StaticSourceFile default methods.

// Java: StaticSourceFile.isStrong()
func (f *SimpleSourceFile) IsStrong() bool { return StaticSourceFile_IsStrong(f) }

// Java: StaticSourceFile.isWeak()
func (f *SimpleSourceFile) IsWeak() bool { return StaticSourceFile_IsWeak(f) }

// Java: StaticSourceFile.isExtern()
func (f *SimpleSourceFile) IsExtern() bool { return StaticSourceFile_IsExtern(f) }

// Java: StaticSourceFile.isNonCode()
func (f *SimpleSourceFile) IsNonCode() bool { return StaticSourceFile_IsNonCode(f) }

// Java: StaticSourceFile.isTypeScriptSource()
func (f *SimpleSourceFile) IsTypeScriptSource() bool { return StaticSourceFile_IsTypeScriptSource(f) }
