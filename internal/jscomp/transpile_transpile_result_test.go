// Port of test/com/google/javascript/jscomp/transpile/TranspileResultTest.java @ closure-compiler v20261015

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestTranspileResult_Equals(t *testing.T) {
	// EqualsTester: members of a group are equal (with equal hash codes),
	// members of different groups are not.
	groups := [][]*TranspileResult{
		{NewTranspileResult(jutil.NewURI("a"), "b", "c", "d"), NewTranspileResult(jutil.NewURI("a"), "b", "c", "d")},
		{NewTranspileResult(jutil.NewURI("A"), "b", "c", "d")},
		{NewTranspileResult(jutil.NewURI("a"), "B", "c", "d")},
		{NewTranspileResult(jutil.NewURI("a"), "b", "C", "d")},
		{NewTranspileResult(jutil.NewURI("a"), "b", "c", "D")},
	}
	for i, g := range groups {
		for j, h := range groups {
			for _, x := range g {
				jsaFalse(t, x.Equals(nil), "equals(null)")
				jsaFalse(t, x.Equals("a string"), "equals(other type)")
				for _, y := range h {
					jsaEq(t, x.Equals(y), i == j, i, " vs ", j)
					if i == j {
						jsaEq(t, x.HashCode(), y.HashCode())
					}
				}
			}
		}
	}
}

func TestTranspileResult_EmbedSourceMap_noSourceMap(t *testing.T) {
	result := NewTranspileResult(jutil.NewURI("a"), "b", "c", "")
	jsaTrue(t, result.EmbedSourcemap() == result)
	jsaTrue(t, result.EmbedSourcemapUrl("foo") == result)
}

func TestTranspileResult_EmbedSourceMap(t *testing.T) {
	result := NewTranspileResult(jutil.NewURI("a"), "b", "c", "{\"version\": 3}")
	jsaJavaEq(t, result.EmbedSourcemap(),
		NewTranspileResult(jutil.NewURI("a"), "b", "c\n//# sourceMappingURL=data:,%7B%22version%22%3A%203%7D\n", ""))
	jsaJavaEq(t, result.EmbedSourcemapUrl("foo.js.map"),
		NewTranspileResult(jutil.NewURI("a"), "b", "c\n//# sourceMappingURL=foo.js.map\n", "{\"version\": 3}"))
}
