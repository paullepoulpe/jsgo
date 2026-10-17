// Port of com/google/javascript/jscomp/IdMappingUtilTest.java @ closure-compiler v20261015

/*
 * Copyright 2015 The Closure Compiler Authors.
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
