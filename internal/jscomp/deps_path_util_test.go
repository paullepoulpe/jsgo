// Port of com/google/javascript/jscomp/deps/PathUtilTest.java @ closure-compiler v20260415

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"os"
	"testing"
)

func TestPathUtil_CollapseDots(t *testing.T) {
	assertEqual(t, "1", PathUtil_CollapseDots("/foo"), "/foo")
	assertEqual(t, "2", PathUtil_CollapseDots("/foo/bar//"), "/foo/bar")
	assertEqual(t, "3", PathUtil_CollapseDots("..//foo/bar"), "foo/bar")
	assertEqual(t, "4", PathUtil_CollapseDots("foo/.././bar"), "bar")
	assertEqual(t, "5", PathUtil_CollapseDots("foo/../../bar/."), "bar")
	assertEqual(t, "6", PathUtil_CollapseDots("./foo/bar/../../baz"), "baz")
	assertEqual(t, "7", PathUtil_CollapseDots("//foo/../bar/../baz/"), "/baz")
	assertEqual(t, "8", PathUtil_CollapseDots("/foo/.."), "/")
	assertEqual(t, "9", PathUtil_CollapseDots("foo/.."), "")
	assertEqual(t, "10", PathUtil_CollapseDots("foo/./././bar//.."), "foo")
}

func TestPathUtil_MakeAbsolute(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, "1", PathUtil_MakeAbsolute("/foo/"), "/foo")
	assertEqual(t, "2", PathUtil_MakeAbsolute("/foo//bar"), "/foo/bar")
	assertEqual(t, "3", PathUtil_MakeAbsolute("foo/"), cwd+"/foo")
	assertEqual(t, "4", PathUtil_MakeAbsolute("foo/b/.."), cwd+"/foo")

	cwd = "/some/root"
	assertEqual(t, "5", PathUtil_MakeAbsoluteWithRoot("/foo/", cwd), "/foo")
	assertEqual(t, "6", PathUtil_MakeAbsoluteWithRoot("/foo//bar", cwd), "/foo/bar")
	assertEqual(t, "7", PathUtil_MakeAbsoluteWithRoot("foo/", cwd), "/some/root/foo")
	assertEqual(t, "8", PathUtil_MakeAbsoluteWithRoot("../foo//./", cwd), "/some/foo")
}

func TestPathUtil_MakeRelative(t *testing.T) {
	cases := [][3]string{
		{"/gfs/bp/pso", "/gfs/bp/pso/simba", "simba"},
		{"/gfs/bp/pso", "/gfs/bp/pso/simba/", "simba"},
		{"gfs/bp/pso/", "gfs/bp/pso/simba", "simba"},
		{"gfs/bp/pso/", "gfs/bp/pso/simba/", "simba"},
		{"/gfs/bp/pso/", "/gfs/bp/pso/simba/ore", "simba/ore"},
		{"/gfs/bp/pso/", "/gfs/bp/pso/simba/ore/", "simba/ore"},
		{"gfs/bp/pso/", "gfs/bp/pso", "."},
		{"gfs/bp/pso", "gfs/bp/pso/", "."},
		{"/gfs/bp/pso", "/gfs/bp/pso///", "."},
		{"/gfs/bp/pso", "//gfs//bp//pso///", "."},

		{"/gfs/bp/pso", "/gfs/bp", ".."},
		{"/gfs/bp/pso", "/gfs", "../.."},
		{"gfs/bp/pso", "gfs/bp", ".."},
		{"gfs/bp/pso", "gfs", "../.."},

		{"gfs/bp/pso", "file", "../../../file"},
		{"/gfs/bp/pso", "/abs/path", "../../../abs/path"},

		{"gfs/bp/pso", "gfs/c/file", "../../c/file"},
		{"/gfs/bp/pso", "/gfs/c/file", "../../c/file"},
	}
	for _, c := range cases {
		assertEqual(t, c[0]+" -> "+c[1], PathUtil_MakeRelative(c[0], c[1]), c[2])
	}
}
