// Port of com/google/javascript/jscomp/graph/GraphReachabilityTest.java @ closure-compiler v20261015

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
	"testing"
)

type graphReachabilityTest struct {
	t            *testing.T
	reachability *GraphReachability[string, string]
	graph        DiGraph[string, string]
}

func TestGraphReachability_Simple(t *testing.T) {
	r := &graphReachabilityTest{t: t}
	r.graph = LinkedDirectedGraph_Create[string, string]()
	r.graph.CreateNode("A")
	r.reachability = NewGraphReachability(r.graph)
	r.reachability.Compute("A")
	r.assertReachable("A")

	r.graph.CreateNode("B")
	r.reachability = NewGraphReachability(r.graph)
	r.reachability.Compute("A")
	r.assertReachable("A")
	r.assertNotReachable("B")

	r.graph.Connect("A", "--->", "B")
	r.reachability = NewGraphReachability(r.graph)
	r.reachability.Compute("B")
	r.assertNotReachable("A")
	r.assertReachable("B")

	r.graph.Connect("B", "--->", "A")
	r.reachability = NewGraphReachability(r.graph)
	r.reachability.Compute("B")
	r.assertReachable("A")
	r.assertReachable("B")

	r.graph.CreateNode("C")
	r.reachability = NewGraphReachability(r.graph)
	r.reachability.Compute("A")
	r.assertReachable("A")
	r.assertReachable("B")
	r.assertNotReachable("C")

	r.graph.CreateNode("D")
	r.graph.Connect("C", "--->", "D")
	r.reachability = NewGraphReachability(r.graph)
	r.reachability.Compute("A")
	r.assertReachable("A")
	r.assertReachable("B")
	r.assertNotReachable("C")
	r.assertNotReachable("D")
	r.reachability.Recompute("C")
	r.assertReachable("C")
	r.assertReachable("D")
}

func (r *graphReachabilityTest) assertReachable(s string) {
	r.t.Helper()
	if GraphReachability_REACHABLE != r.graph.GetNode(s).GetAnnotation() {
		r.t.Errorf("%s should be reachable", s)
	}
}

func (r *graphReachabilityTest) assertNotReachable(s string) {
	r.t.Helper()
	if GraphReachability_REACHABLE == r.graph.GetNode(s).GetAnnotation() {
		r.t.Errorf("%s should not be reachable", s)
	}
}
