// Port of com/google/javascript/jscomp/DotFormatterTest.java @ closure-compiler v20261015

/*
 * Copyright 2007 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Tests that keys are assigned sequentially.
func TestDotFormatter_KeyAssignementSequential(t *testing.T) {
	dot := DotFormatter_NewInstanceForTesting()
	for want := 0; want < 5; want++ {
		if got := dot.Key(rhino.NewNode(rhino.Token_BLOCK)); got != want {
			t.Errorf("key: got %d, want %d", got, want)
		}
	}
}

// Tests that keys are assigned once per node.
func TestDotFormatter_KeyAssignementOncePerNode(t *testing.T) {
	dot := DotFormatter_NewInstanceForTesting()
	node0 := rhino.NewNode(rhino.Token_BLOCK)
	node1 := rhino.NewNode(rhino.Token_BLOCK)
	node2 := rhino.NewNode(rhino.Token_BLOCK)

	for _, tc := range []struct {
		n    *rhino.Node
		want int
	}{{node0, 0}, {node1, 1}, {node2, 2}, {node0, 0}, {node1, 1}, {node2, 2}} {
		if got := dot.Key(tc.n); got != tc.want {
			t.Errorf("key: got %d, want %d", got, tc.want)
		}
	}
}

// Tests the formatting (simple tree).
func TestDotFormatter_ToDotSimple(t *testing.T) {
	ast := rhino.NewNode(rhino.Token_BITOR)

	expected := "digraph AST {\n" +
		"  node [color=lightblue2, style=filled];\n" +
		"  node0 [label=\"BITOR\"];\n" +
		"}\n"
	dotFormatterTestTest(t, expected, ast)
}

// Tests the formatting (simple tree).
func TestDotFormatter_ToDotSimpleName(t *testing.T) {
	ast := rhino.Node_NewStringWithToken(rhino.Token_NAME, "dummy")

	expected := "digraph AST {\n" +
		"  node [color=lightblue2, style=filled];\n" +
		"  node0 [label=\"NAME(dummy)\"];\n" +
		"}\n"
	dotFormatterTestTest(t, expected, ast)
}

// Tests the formatting (simple tree).
func TestDotFormatter_ToDotSimpleStringKey(t *testing.T) {
	ast := rhino.Node_NewStringWithToken(rhino.Token_STRING_KEY, "key")

	expected := "digraph AST {\n" +
		"  node [color=lightblue2, style=filled];\n" +
		"  node0 [label=\"STRING_KEY(key)\"];\n" +
		"}\n"
	dotFormatterTestTest(t, expected, ast)
}

// Tests the formatting (3 element tree).
func TestDotFormatter_ToDot3Elements_nodesWithoutNames(t *testing.T) {
	ast := rhino.NewNode(rhino.Token_BLOCK)
	ast.AddChildToBack(rhino.Node_NewStringWithToken(rhino.Token_NAME, ""))
	ast.AddChildToBack(rhino.Node_NewStringWithToken(rhino.Token_STRINGLIT, ""))

	expected := "digraph AST {\n" +
		"  node [color=lightblue2, style=filled];\n" +
		"  node0 [label=\"BLOCK\"];\n" +
		"  node1 [label=\"NAME\"];\n" +
		"  node0 -> node1 [weight=1];\n" +
		"  node2 [label=\"STRINGLIT\"];\n" +
		"  node0 -> node2 [weight=1];\n" +
		"}\n"
	dotFormatterTestTest(t, expected, ast)
}

// Tests the formatting (3 element tree).
func TestDotFormatter_ToDot3Elements_NodesCreatedWithNewString(t *testing.T) {
	ast := rhino.NewNode(rhino.Token_BLOCK)
	ast.AddChildToBack(rhino.Node_NewStringWithToken(rhino.Token_NAME, "a"))
	ast.AddChildToBack(rhino.Node_NewStringWithToken(rhino.Token_STRINGLIT, "b"))

	expected := "digraph AST {\n" +
		"  node [color=lightblue2, style=filled];\n" +
		"  node0 [label=\"BLOCK\"];\n" +
		"  node1 [label=\"NAME(a)\"];\n" +
		"  node0 -> node1 [weight=1];\n" +
		"  node2 [label=\"STRINGLIT(b)\"];\n" +
		"  node0 -> node2 [weight=1];\n" +
		"}\n"
	dotFormatterTestTest(t, expected, ast)
}

// Tests the labels for `import * as name from "module-name";`
func TestDotFormatter_ImportStarLabel(t *testing.T) {
	ast := rhino.NewNode(rhino.Token_IMPORT)
	ast.AddChildToBack(rhino.NewNode(rhino.Token_EMPTY))
	ast.AddChildToBack(rhino.Node_NewStringWithToken(rhino.Token_IMPORT_STAR, "name"))
	ast.AddChildToBack(rhino.Node_NewStringWithToken(rhino.Token_STRINGLIT, "module-name"))

	expected := "digraph AST {\n" +
		"  node [color=lightblue2, style=filled];\n" +
		"  node0 [label=\"IMPORT\"];\n" +
		"  node1 [label=\"EMPTY\"];\n" +
		"  node0 -> node1 [weight=1];\n" +
		"  node2 [label=\"IMPORT_STAR(name)\"];\n" +
		"  node0 -> node2 [weight=1];\n" +
		"  node3 [label=\"STRINGLIT(module-nam)\"];\n" +
		"  node0 -> node3 [weight=1];\n" +
		"}\n"
	dotFormatterTestTest(t, expected, ast)
}

// Java: DotFormatterTest.test(String, Node)
func dotFormatterTestTest(t *testing.T, expected string, ast *rhino.Node) {
	t.Helper()
	if got := DotFormatter_ToDot(ast); got != expected {
		t.Errorf("toDot:\n got: %q\nwant: %q", got, expected)
	}
}
