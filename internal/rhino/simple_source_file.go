// Port of com/google/javascript/rhino/SimpleSourceFile.java
// @ closure-compiler v20260415

package rhino

import "math"

// SimpleSourceFile is a simple implementation of StaticSourceFile used for
// testing. Mirrors com.google.javascript.rhino.SimpleSourceFile.
type SimpleSourceFile struct {
	name string
	kind SourceKind
}

// NewSimpleSourceFile constructs a SimpleSourceFile with the given name and kind.
func NewSimpleSourceFile(name string, kind SourceKind) *SimpleSourceFile {
	return &SimpleSourceFile{name: name, kind: kind}
}

func (s *SimpleSourceFile) GetName() string     { return s.name }
func (s *SimpleSourceFile) GetKind() SourceKind { return s.kind }

func (s *SimpleSourceFile) IsClosureUnawareCode() bool { return false }

func (s *SimpleSourceFile) GetColumnOfOffset(offset int) int { return 0 }
func (s *SimpleSourceFile) GetLineOfOffset(offset int) int   { return 1 }

// GetLineOffset mirrors SimpleSourceFile.getLineOffset: returns
// math.MinInt32 when the source file has no line information (which is
// always the case for SimpleSourceFile). Panics if line < 1.
func (s *SimpleSourceFile) GetLineOffset(line int) int {
	if line < 1 {
		panic("Should not call GetLineOffset with line number " + itoa(line))
	}
	return math.MinInt32
}

func (s *SimpleSourceFile) String() string { return s.name }

// itoa is a minimal int-to-string helper to avoid importing strconv in a
// panic message that is only reached on programming errors.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := make([]byte, 0, 20)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	if neg {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}
