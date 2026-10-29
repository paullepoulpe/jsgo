// Port of com/google/javascript/jscomp/IdMappingUtilTest.java @ closure-compiler v20261021

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
	"strings"
	"testing"
)

func TestIdMappingUtil_ParseIdMapping(t *testing.T) {
	s := "[gen1]\nid1:data1\nid2:data2:data22\n\n[gen2]"
	result := IdMappingUtil_ParseSerializedIdMappings(&s)

	assertEqual(t, "size", len(result), 2)

	gen1, ok := result["gen1"]
	if !ok {
		t.Fatal("missing gen1")
	}
	assertEqual(t, "gen1 size", gen1.Size(), 2)
	assertEqual(t, "id1", gen1.Get("id1"), "data1")
	assertEqual(t, "id2", gen1.Get("id2"), "data2:data22")

	gen2, ok := result["gen2"]
	if !ok {
		t.Fatal("missing gen2")
	}
	assertEqual(t, "gen2 empty", gen2.IsEmpty(), true)
}

func TestIdMappingUtil_ParseSectionAsStream(t *testing.T) {
	s := "[gen1]\nid1:data1\n[gen2]\nid2:data2"
	result := IdMappingUtil_ParseSectionAsStream(strings.NewReader(s), "gen1")

	assertEqual(t, "result", result.String(), "{id1=data1}")
}
