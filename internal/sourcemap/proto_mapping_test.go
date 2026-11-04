package sourcemap

import (
	"testing"
)

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
