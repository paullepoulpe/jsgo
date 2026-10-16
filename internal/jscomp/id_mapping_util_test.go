// Port of com/google/javascript/jscomp/IdMappingUtilTest.java @ closure-compiler v20260415

package jscomp

import (
	"bytes"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestIdMappingUtil_ParseIdMapping(t *testing.T) {
	mapping := jutil.NewStringBuilder()
	mapping.
		Append("[gen1]").
		AppendChar(IdMappingUtil_NEW_LINE).
		Append("id1:data1").
		AppendChar(IdMappingUtil_NEW_LINE).
		Append("id2:data2:data22").
		AppendChar(IdMappingUtil_NEW_LINE).
		AppendChar(IdMappingUtil_NEW_LINE).
		Append("[gen2]")

	s := mapping.String()
	result := IdMappingUtil_ParseSerializedIdMappings(&s)

	assertEqual(t, "size", result.Size(), 2)

	if !result.ContainsKey("gen1") {
		t.Fatal("missing gen1")
	}
	gen1 := result.Get("gen1")
	assertEqual(t, "gen1 size", gen1.Size(), 2)
	assertEqual(t, "id1", gen1.Get("id1"), "data1")
	assertEqual(t, "id2", gen1.Get("id2"), "data2:data22")

	if !result.ContainsKey("gen2") {
		t.Fatal("missing gen2")
	}
	assertEqual(t, "gen2 empty", result.Get("gen2").IsEmpty(), true)
}

func TestIdMappingUtil_ParseSectionAsStream(t *testing.T) {
	mapping := jutil.NewStringBuilder()
	mapping.
		Append("[gen1]").
		AppendChar(IdMappingUtil_NEW_LINE).
		Append("id1:data1").
		AppendChar(IdMappingUtil_NEW_LINE).
		Append("[gen2]").
		AppendChar(IdMappingUtil_NEW_LINE).
		Append("id2:data2")

	result := IdMappingUtil_ParseSectionAsStream(
		bytes.NewReader(jutil.StringGetBytesUTF8(mapping.String())), "gen1")

	assertEqual(t, "result", result.String(), "{id1=data1}")
}
