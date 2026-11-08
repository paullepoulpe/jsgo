// Port of test/com/google/javascript/jscomp/bundle/SourceTest.java @ closure-compiler v20261104

/*
 * Copyright 2017 The Closure Compiler Authors.
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

package bundle

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// jsaEqualsTester is Guava's EqualsTester: members of a group are equal to
// each other with equal hash codes, and unequal to members of other groups,
// to null and to an unrelated object.
func jsaEqualsTester(t *testing.T, groups ...[]jsaHashable) {
	t.Helper()
	for i, g := range groups {
		for _, x := range g {
			jsaFalse(t, x.Equals(nil), "group ", i, " equals(null)")
			jsaFalse(t, x.Equals(jsaNotAnything{}), "group ", i, " equals(other)")
			for j, h := range groups {
				for _, y := range h {
					if x.Equals(y) != (i == j) {
						t.Errorf("group %d item %v equals group %d item %v: %v", i, x, j, y, x.Equals(y))
					}
					if i == j && x.HashCode() != y.HashCode() {
						t.Errorf("group %d: hash codes differ", i)
					}
				}
			}
		}
	}
}

type jsaNotAnything struct{}

type jsaHashable = interface {
	Equals(any) bool
	HashCode() int32
}

func jsaGroup[T jsaHashable](xs ...T) []jsaHashable {
	var out []jsaHashable
	for _, x := range xs {
		out = append(out, x)
	}
	return out
}

func TestSource_Equals(t *testing.T) {
	flags := map[string]string{"x": "y"}
	jsaEqualsTester(t,
		jsaGroup(NewSource(SourceConfig{}), NewSource(SourceConfig{})),
		jsaGroup(
			NewSource(SourceConfig{CodeSupplier: func() string { return "foo" }}),
			NewSource(SourceConfig{CodeSupplier: func() string { return "foo" }}),
			NewSource(SourceConfig{Code: "foo"}),
			NewSource(SourceConfig{Code: "foo"})),
		jsaGroup(NewSource(SourceConfig{Code: "bar"})),
		jsaGroup(
			NewSource(SourceConfig{Code: "foo", Path: jutil.Path_Of("/x")}),
			NewSource(SourceConfig{Code: "foo", Path: jutil.Path_Of("/x")})),
		jsaGroup(
			NewSource(SourceConfig{Code: "foo", OriginalCode: "bar"}),
			NewSource(SourceConfig{Code: "foo", OriginalCode: "bar"})),
		jsaGroup(NewSource(SourceConfig{Code: "foo", SourceUrl: "foo"})),
		jsaGroup(NewSource(SourceConfig{Code: "foo", SourceMap: "foo"})),
		jsaGroup(NewSource(SourceConfig{Code: "foo", SourceMappingUrl: "foo"})),
		jsaGroup(NewSource(SourceConfig{Code: "foo", EstimatedSize: 12})),
		jsaGroup(NewSource(SourceConfig{Code: "foo", LoadFlags: flags})),
		jsaGroup(NewSource(SourceConfig{Code: "foo", Runtimes: []string{"x"}})))
}

func TestSource_Code(t *testing.T) {
	jsaEq(t, NewSource(SourceConfig{Code: "foo"}).Code(), "foo")
	jsaEq(t, NewSource(SourceConfig{CodeSupplier: func() string { return "foo" }}).Code(), "foo")
}

func TestSource_Original(t *testing.T) {
	jsaEq(t, NewSource(SourceConfig{Code: "foo"}).OriginalCode(), "foo")
	jsaEq(t, NewSource(SourceConfig{CodeSupplier: func() string { return "foo" }}).OriginalCode(), "foo")

	jsaEq(t, NewSource(SourceConfig{Code: "foo", OriginalCode: "bar"}).OriginalCode(), "bar")
	jsaEq(t, NewSource(SourceConfig{
		CodeSupplier: func() string { return "foo" },
		OriginalCode: "bar",
	}).OriginalCode(), "bar")
}

func TestSource_Runtimes(t *testing.T) {
	foo := NewSource(SourceConfig{Runtimes: []string{"foo"}})
	jsaDeepEq(t, foo.Runtimes(), []string{"foo"})
	jsaDeepEq(t, foo.WithTransformed(foo.Code(), foo.SourceMap(), "bar").Runtimes(), []string{"foo", "bar"})
	jsaDeepEq(t, foo.WithTransformed(foo.Code(), foo.SourceMap(), "bar", "baz").Runtimes(), []string{"foo", "bar", "baz"})
	jsaDeepEq(t, foo.WithTransformed(foo.Code(), foo.SourceMap(), "bar", "foo").Runtimes(), []string{"foo", "bar"})
	jsaDeepEq(t, NewSource(SourceConfig{Runtimes: []string{"foo", "bar", "foo"}}).Runtimes(), []string{"foo", "bar"})
}
