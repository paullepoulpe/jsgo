package sourcemap

import (
	"reflect"
	"runtime"
	"testing"
	"unsafe"
)

// Candidate-only representation check; omitted from stock API positive controls.
func TestSourceMapMappingSourceCellLayout(t *testing.T) {
	field, ok := reflect.TypeOf(SourceMapGeneratorV3_Mapping{}).FieldByName("sourceFile")
	if !ok || field.Type != reflect.TypeOf((*string)(nil)) {
		t.Fatal("Mapping source must be a GC-tracked string pointer")
	}
	if runtime.GOARCH == "amd64" && unsafe.Sizeof(SourceMapGeneratorV3_Mapping{}) != 40 {
		t.Fatalf("Mapping size=%d, expected40", unsafe.Sizeof(SourceMapGeneratorV3_Mapping{}))
	}
	t.Logf("COMPILED_SOURCE_MAPPING_LAYOUT bytes=%d align=%d", unsafe.Sizeof(SourceMapGeneratorV3_Mapping{}), unsafe.Alignof(SourceMapGeneratorV3_Mapping{}))
}
