// Port of com/google/javascript/jscomp/deps/DependencyResolverTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/deps/DependencyResolverTest.java)

/*
 * Copyright 2009 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Tests for DependencyResolver.
//
// Go-only: test-helper — the JUnit class is a plain Go test per @Test
// method; the fields and setUp are dependencyResolverTest_setUp.

type dependencyResolverTest struct {
	fakeDeps1 *DependencyFile
	fakeDeps2 *DependencyFile
	resolver  *DefaultDependencyResolver
}

// Java: DependencyResolverTest.setUp()
func dependencyResolverTest_setUp() *dependencyResolverTest {
	d := &dependencyResolverTest{}
	d.fakeDeps1 = NewDependencyFile(NewVirtualFile("deps1",
		"goog.addDependency('a.js', ['a'], []);\n"+
			"goog.addDependency('b.js', ['b'], []);\n"+
			"goog.addDependency('c.js', ['c', 'c2'], ['a']);\n"+
			"goog.addDependency('d.js', ['d'], ['b', 'c']);\n"))
	d.fakeDeps2 = NewDependencyFile(NewVirtualFile("deps2",
		"goog.addDependency('e.js', ['e'], ['c2']);\n"+
			"goog.addDependency('f.js', ['f'], ['b', 'c']);\n"+
			"goog.addDependency('g.js', ['g'], ['a', 'b', 'c']);\n"+
			"goog.addDependency('h.js', ['h', 'i'], ['g', 'd', 'c']);\n"))
	d.resolver = NewDefaultDependencyResolver([]*DependencyFile{d.fakeDeps1, d.fakeDeps2}, false)
	return d
}

func dependencyResolverTest_assertJoined(t *testing.T, deps []string, want string) {
	t.Helper()
	if got := strings.Join(deps, ","); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestDependencyResolver_BasicCase(t *testing.T) {
	d := dependencyResolverTest_setUp()
	deps := d.resolver.GetDependencies("goog.require('a');")
	dependencyResolverTest_assertJoined(t, deps, "base.js,a.js")
}

func TestDependencyResolver_SimpleDependencies(t *testing.T) {
	d := dependencyResolverTest_setUp()
	deps := d.resolver.GetDependencies("goog.require('c');")
	dependencyResolverTest_assertJoined(t, deps, "base.js,a.js,c.js")
}

func TestDependencyResolver_TransitiveDependencies(t *testing.T) {
	d := dependencyResolverTest_setUp()
	deps := d.resolver.GetDependencies("goog.require('e');")
	dependencyResolverTest_assertJoined(t, deps, "base.js,a.js,c.js,e.js")
}

func TestDependencyResolver_MultipleRequires(t *testing.T) {
	d := dependencyResolverTest_setUp()
	deps := d.resolver.GetDependencies("goog.require('e');goog.require('a');goog.require('b');")
	dependencyResolverTest_assertJoined(t, deps, "base.js,a.js,c.js,e.js,b.js")
}

func TestDependencyResolver_OneMoreForGoodMeasure(t *testing.T) {
	d := dependencyResolverTest_setUp()
	deps := d.resolver.GetDependencies("goog.require('g');goog.require('f');goog.require('c');")
	dependencyResolverTest_assertJoined(t, deps, "base.js,a.js,b.js,c.js,g.js,f.js")
}

func TestDependencyResolver_SharedSeenSetNoBaseFile(t *testing.T) {
	d := dependencyResolverTest_setUp()
	seen := jutil.NewLinkedHashSet[string]()

	deps := d.resolver.GetDependenciesWithSeen(
		"goog.require('g');goog.require('f');goog.require('c');", seen, false)

	depsLater := d.resolver.GetDependenciesWithSeen("goog.require('f');goog.require('c');", seen, false)

	dependencyResolverTest_assertJoined(t, deps, "a.js,b.js,c.js,g.js,f.js")
	if len(depsLater) != 0 {
		t.Fatalf("expected empty, got %v", depsLater)
	}
}

func TestDependencyResolver_SharedSeenSetNoBaseFileNewRequires(t *testing.T) {
	d := dependencyResolverTest_setUp()
	seen := jutil.NewLinkedHashSet[string]()

	deps := d.resolver.GetDependenciesWithSeen("goog.require('f');goog.require('c');", seen, false)

	depsLater := d.resolver.GetDependenciesWithSeen("goog.require('g');goog.require('c');", seen, false)

	dependencyResolverTest_assertJoined(t, deps, "b.js,a.js,c.js,f.js")
	dependencyResolverTest_assertJoined(t, depsLater, "g.js")
}

func TestDependencyResolver_SharedSeenSetNoBaseFileMultipleProvides(t *testing.T) {
	d := dependencyResolverTest_setUp()
	seen := jutil.NewLinkedHashSet[string]()

	deps := d.resolver.GetDependenciesWithSeen("goog.require('h');goog.require('i');", seen, false)

	dependencyResolverTest_assertJoined(t, deps, "a.js,b.js,c.js,g.js,d.js,h.js")
}

func TestDependencyResolver_NonExistentProvideLoose(t *testing.T) {
	d := dependencyResolverTest_setUp()
	seen := jutil.NewLinkedHashSet[string]()
	d.resolver = NewDefaultDependencyResolver([]*DependencyFile{d.fakeDeps1}, false)
	deps := d.resolver.GetDependenciesWithSeen("goog.require('foo');goog.require('d');", seen, false)

	dependencyResolverTest_assertJoined(t, deps, "b.js,a.js,c.js,d.js")
}

func TestDependencyResolver_NonExistentProvideStrict(t *testing.T) {
	d := dependencyResolverTest_setUp()
	seen := jutil.NewLinkedHashSet[string]()
	d.resolver = NewDefaultDependencyResolver([]*DependencyFile{d.fakeDeps1}, true)
	thrown := false
	jutil.Try(func() {
		d.resolver.GetDependenciesWithSeen("goog.require('foo');goog.require('a');", seen, false)
		t.Fatal("Service exception should be thrown")
	}, func(expected IServiceException) {
		thrown = true
	})
	if !thrown {
		t.Fatal("Service exception should be thrown")
	}
}
