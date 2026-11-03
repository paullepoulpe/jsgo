package sourcemap

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

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

func TestMappingOriginalMapping(t *testing.T) {
	dm := Mapping_OriginalMapping_GetDefaultInstance()
	if !dm.Equals(&Mapping_OriginalMapping{}) {
		t.Fatalf("default instance should equal zero-value Mapping_OriginalMapping")
	}
	if dm.Equals(nil) || dm.Equals((*Mapping_OriginalMapping)(nil)) || dm.Equals("other") {
		t.Fatalf("Equals should reject nil and non-Mapping_OriginalMapping values")
	}

	m := &Mapping_OriginalMapping{
		OriginalFile:   "a.js",
		LineNumber:     3,
		ColumnPosition: 15,
		Identifier:     "foo",
		Precision:      Mapping_OriginalMapping_Precision_EXACT,
	}
	if m.GetOriginalFile() != "a.js" || m.GetLineNumber() != 3 || m.GetColumnPosition() != 15 ||
		m.GetIdentifier() != "foo" || m.GetPrecision() != Mapping_OriginalMapping_Precision_EXACT ||
		m.GetPrecision().String() != "EXACT" {
		t.Fatalf("unexpected getters on %+v", m)
	}

	m2 := *m
	if !m.Equals(&m2) {
		t.Fatalf("expected copies to compare equal")
	}
	m2.SetOriginalFileAndColumnInPlace("dir/a.js", 14)
	if m2.GetOriginalFile() != "dir/a.js" || m2.GetColumnPosition() != 14 {
		t.Fatalf("SetOriginalFileAndColumnInPlace failed: %+v", m2)
	}
	if m.Equals(&m2) {
		t.Fatalf("expected mutated copy to differ")
	}
}
