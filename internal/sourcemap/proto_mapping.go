// Reconstruction of com.google.debugging.sourcemap.proto.Mapping (generated
// by protoc from src/com/google/debugging/sourcemap/proto/mapping.proto) @
// closure-compiler v20261104.
//
// Go-only: language-adapter — the Java class is protobuf-generated and not
// part of the Java sources; it is represented here from mapping.proto as a
// plain Go struct with value equality.

package sourcemap

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Mapping_OriginalMapping represents sourcemap.OriginalMapping.
type Mapping_OriginalMapping struct {
	OriginalFile   string
	LineNumber     int
	ColumnPosition int
	Identifier     string
	Precision      Mapping_OriginalMapping_Precision
}

// Mapping_OriginalMapping_Precision is OriginalMapping.Precision.
type Mapping_OriginalMapping_Precision int32

const (
	Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION Mapping_OriginalMapping_Precision = 0
	Mapping_OriginalMapping_Precision_EXACT             Mapping_OriginalMapping_Precision = 1
	Mapping_OriginalMapping_Precision_APPROXIMATE_LINE  Mapping_OriginalMapping_Precision = 2
)

var mapping_OriginalMapping_Precision_names = [...]string{"UNKNOWN_PRECISION", "EXACT", "APPROXIMATE_LINE"}

// String is Enum.name().
func (p Mapping_OriginalMapping_Precision) String() string {
	return mapping_OriginalMapping_Precision_names[p]
}

var mapping_OriginalMapping_DEFAULT_INSTANCE = &Mapping_OriginalMapping{}

// Mapping_OriginalMapping_GetDefaultInstance returns the default empty mapping.
//
// Java: OriginalMapping.getDefaultInstance()
func Mapping_OriginalMapping_GetDefaultInstance() *Mapping_OriginalMapping {
	return mapping_OriginalMapping_DEFAULT_INSTANCE
}

// Java: OriginalMapping.getOriginalFile()
func (m *Mapping_OriginalMapping) GetOriginalFile() string {
	return m.OriginalFile
}

// Java: OriginalMapping.getLineNumber()
func (m *Mapping_OriginalMapping) GetLineNumber() int {
	return m.LineNumber
}

// Java: OriginalMapping.getColumnPosition()
func (m *Mapping_OriginalMapping) GetColumnPosition() int {
	return m.ColumnPosition
}

// SetOriginalFileAndColumnInPlace mutates OriginalFile and ColumnPosition in place.
func (m *Mapping_OriginalMapping) SetOriginalFileAndColumnInPlace(originalFile string, columnPosition int) {
	m.OriginalFile = originalFile
	m.ColumnPosition = columnPosition
}

// Java: OriginalMapping.getIdentifier()
func (m *Mapping_OriginalMapping) GetIdentifier() string {
	return m.Identifier
}

// Java: OriginalMapping.getPrecision()
func (m *Mapping_OriginalMapping) GetPrecision() Mapping_OriginalMapping_Precision {
	if uint(m.Precision) <= 2 {
		return m.Precision
	}
	return Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION
}

// Equals compares two OriginalMapping values for field-wise equality.
//
// Java: OriginalMapping.equals(Object)
func (m *Mapping_OriginalMapping) Equals(obj any) bool {
	other, ok := obj.(*Mapping_OriginalMapping)
	return ok && other != nil && *m == *other
}

// LookupMappingForDiff resolves (lineNumber, column) via GetMappingForLine and
// additionally reports whether the underlying VLQ entry carried an explicit
// name index (distinguishing an explicit empty-string name from UNMAPPED in the
// differential oracle).
func LookupMappingForDiff(c *SourceMapConsumerV3, lineNumber, column int) (*Mapping_OriginalMapping, bool) {
	ln := int32(lineNumber) - 1
	col := int32(column) - 1
	if ln < 0 || int(ln) >= c.mappings.GetParsedLineCount() {
		return nil, false
	}
	entries := c.mappings.line(int(ln))
	if len(entries) == 0 || entries[0].genCol > col {
		prevLine := int(ln)
		for {
			if prevLine == 0 {
				return nil, false
			}
			prevLine--
			if !c.mappings.lineIsEmpty(prevLine) {
				break
			}
		}
		prevEntries := c.mappings.line(prevLine)
		e := &prevEntries[len(prevEntries)-1]
		m := c.GetOriginalMappingForEntry(e, Mapping_OriginalMapping_Precision_APPROXIMATE_LINE)
		return m, m != nil && e.nameId != SourceMapConsumerV3_UNMAPPED
	}
	index := sourceMapConsumerV3_search(entries, col, 0, len(entries)-1)
	jutil.CheckStateMsg1(index >= 0, "unexpected:%s", index*c.mappings.entrySize)
	e := &entries[index]
	m := c.GetOriginalMappingForEntry(e, Mapping_OriginalMapping_Precision_EXACT)
	return m, m != nil && e.nameId != SourceMapConsumerV3_UNMAPPED
}
