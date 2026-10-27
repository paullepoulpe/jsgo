// Port of test/com/google/javascript/jscomp/bundle/SourceTest.java @ closure-compiler v20261021

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
	foo := Source_Builder_().SetCode("foo").Build()
	flags := jutil.NewLinkedHashMap[string, string]()
	flags.Put("x", "y")
	jsaEqualsTester(t,
		jsaGroup(Source_Builder_().Build(), Source_Builder_().Build()),
		jsaGroup(
			Source_Builder_().SetCodeFromSupplier(func() string { return "foo" }).Build(),
			Source_Builder_().SetCodeFromSupplier(func() string { return "foo" }).Build(),
			Source_Builder_().SetCode("foo").Build(),
			Source_Builder_().SetCode("foo").Build()),
		jsaGroup(Source_Builder_().SetCode("bar").Build()),
		jsaGroup(
			foo.ToBuilder().SetPath(jutil.Path_Of("/x")).Build(),
			foo.ToBuilder().SetPath(jutil.Path_Of("/x")).Build()),
		jsaGroup(
			foo.ToBuilder().SetOriginalCode("bar").Build(),
			foo.ToBuilder().SetOriginalCode("bar").Build()),
		jsaGroup(foo.ToBuilder().SetSourceUrl("foo").Build()),
		jsaGroup(foo.ToBuilder().SetSourceMap("foo").Build()),
		jsaGroup(foo.ToBuilder().SetSourceMappingUrl("foo").Build()),
		jsaGroup(foo.ToBuilder().SetEstimatedSize(12).Build()),
		jsaGroup(foo.ToBuilder().SetLoadFlags(flags).Build()),
		jsaGroup(foo.ToBuilder().SetRuntimes(jutil.NewLinkedHashSet("x")).Build()))
}

func TestSource_Code(t *testing.T) {
	jsaEq(t, Source_Builder_().SetCode("foo").Build().Code(), "foo")
	jsaEq(t, Source_Builder_().SetCodeFromSupplier(func() string { return "foo" }).Build().Code(), "foo")
}

func TestSource_Original(t *testing.T) {
	jsaEq(t, Source_Builder_().SetCode("foo").Build().OriginalCode(), "foo")
	jsaEq(t, Source_Builder_().SetCodeFromSupplier(func() string { return "foo" }).Build().OriginalCode(), "foo")

	jsaEq(t, Source_Builder_().SetCode("foo").SetOriginalCode("bar").Build().OriginalCode(), "bar")
	jsaEq(t, Source_Builder_().
		SetCodeFromSupplier(func() string { return "foo" }).
		SetOriginalCode("bar").
		Build().
		OriginalCode(), "bar")
}

func TestSource_Runtimes(t *testing.T) {
	foo := Source_Builder_().AddRuntime("foo").Build()
	jsaDeepEq(t, foo.Runtimes().ToSlice(), []string{"foo"})
	jsaDeepEq(t, foo.ToBuilder().AddRuntime("bar").Build().Runtimes().ToSlice(), []string{"foo", "bar"})
	jsaDeepEq(t, foo.ToBuilder().AddRuntime("bar").AddRuntime("baz").Build().Runtimes().ToSlice(), []string{"foo", "bar", "baz"})
	jsaDeepEq(t, foo.ToBuilder().AddRuntime("bar").AddRuntime("foo").Build().Runtimes().ToSlice(), []string{"foo", "bar"})
}
