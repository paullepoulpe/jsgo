// Port of test/com/google/javascript/jscomp/ControlFlowAnalysisTest.java @ closure-compiler v20261021

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
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// cfaTestCfg: given an input in JavaScript, test if the control flow
// analysis creates the proper control flow graph by comparing the expected
// Dot file output.
//
// shouldTraverseFunctions is whether to traverse functions when constructing
// the CFG. Passed in to the builder for ControlFlowAnalysis. Defaults to true
// for the two-arg testCfg(String, String).
//
// Java: ControlFlowAnalysisTest.testCfg(String, String, boolean)
func cfaTestCfg(t *testing.T, input string, expected string, shouldTraverseFunctions bool) {
	t.Helper()
	compiler := NewCompiler()
	root := compiler.ParseSyntheticCode("cfgtest", input)

	cfg := ControlFlowAnalysis_Builder_().
		SetCompiler(compiler).
		SetCfgRoot(root).
		SetTraverseFunctions(shouldTraverseFunctions).
		SetIncludeEdgeAnnotations(true).
		ComputeCfg()

	if got := DotFormatter_ToDotWithCfg(root, cfg); got != expected {
		t.Errorf("DotFormatter.toDot(root, cfg):\n%s\nwant:\n%s", got, expected)
	}
}

// cfaGetAllEdges gets all the edges of the graph.
//
// Java: ControlFlowAnalysisTest.getAllEdges(ControlFlowGraph<Node>)
func cfaGetAllEdges(cfg ControlFlowGraph[rhino.Node]) []DiGraph_DiGraphEdge[rhino.Node, ControlFlowGraph_Branch] {
	var edges []DiGraph_DiGraphEdge[rhino.Node, ControlFlowGraph_Branch]
	for _, n := range cfg.GetDiGraphNodes() {
		for _, e := range cfg.GetOutEdges(n.GetValue()) {
			edges = append(edges, e)
		}
	}
	return edges
}

// cfaGetAllEdgesBetween gets all the control flow edges from some node with
// the first token to some node with the second token.
//
// Java: ControlFlowAnalysisTest.getAllEdges(ControlFlowGraph<Node>, Token, Token)
func cfaGetAllEdgesBetween(
	cfg ControlFlowGraph[rhino.Node], startToken rhino.Token, endToken rhino.Token) []DiGraph_DiGraphEdge[rhino.Node, ControlFlowGraph_Branch] {
	var edges []DiGraph_DiGraphEdge[rhino.Node, ControlFlowGraph_Branch]
	for _, edge := range cfaGetAllEdges(cfg) {
		startNode := edge.GetSource().GetValue()
		endNode := edge.GetDestination().GetValue()
		if startNode.IsNil() ||
			endNode.IsNil() ||
			startNode.GetToken() != startToken ||
			endNode.GetToken() != endToken {
			continue
		}
		edges = append(edges, edge)
	}
	return edges
}

// cfaGetAllEdgesOfType gets all the control flow edges of the given type
// from some node with the first token to some node with the second token.
//
// Java: ControlFlowAnalysisTest.getAllEdges(ControlFlowGraph<Node>, Token, Token, Branch)
func cfaGetAllEdgesOfType(
	cfg ControlFlowGraph[rhino.Node], startToken rhino.Token, endToken rhino.Token, type_ ControlFlowGraph_Branch) []DiGraph_DiGraphEdge[rhino.Node, ControlFlowGraph_Branch] {
	var edges []DiGraph_DiGraphEdge[rhino.Node, ControlFlowGraph_Branch]
	for _, elem := range cfaGetAllEdgesBetween(cfg, startToken, endToken) {
		if type_ == elem.GetValue() {
			edges = append(edges, elem)
		}
	}
	return edges
}

// Java: ControlFlowAnalysisTest.isAncestor(Node, Node)
func cfaIsAncestor(n rhino.Node, maybeDescendant rhino.Node) bool {
	for current := n.GetFirstChild(); !current.IsNil(); current = current.GetNext() {
		if current == maybeDescendant || cfaIsAncestor(current, maybeDescendant) {
			return true
		}
	}

	return false
}

// cfaGetAllDownEdges gets all the control flow edges of the given type from
// some node with the first token to some node with the second token. This
// edge must flow from a parent to one of its descendants.
//
// Java: ControlFlowAnalysisTest.getAllDownEdges(ControlFlowGraph<Node>, Token, Token, Branch)
func cfaGetAllDownEdges(
	cfg ControlFlowGraph[rhino.Node], startToken rhino.Token, endToken rhino.Token, type_ ControlFlowGraph_Branch) []DiGraph_DiGraphEdge[rhino.Node, ControlFlowGraph_Branch] {
	var edges []DiGraph_DiGraphEdge[rhino.Node, ControlFlowGraph_Branch]
	for _, edge := range cfaGetAllEdgesOfType(cfg, startToken, endToken, type_) {
		source := edge.GetSource().GetValue()
		dest := edge.GetDestination().GetValue()
		if cfaIsAncestor(source, dest) {
			edges = append(edges, edge)
		}
	}

	return edges
}

// cfaAssertNoEdge asserts that there exists no control flow edge of the
// given type from some node with the first token to some node with the
// second token.
//
// Java: ControlFlowAnalysisTest.assertNoEdge(ControlFlowGraph<Node>, Token, Token)
func cfaAssertNoEdge(t *testing.T, cfg ControlFlowGraph[rhino.Node], startToken rhino.Token, endToken rhino.Token) {
	t.Helper()
	if edges := cfaGetAllEdgesBetween(cfg, startToken, endToken); len(edges) != 0 {
		t.Errorf("expected no edge %v -> %v, found %d", startToken, endToken, len(edges))
	}
}

// Java: ControlFlowAnalysisTest.assertNoEdge(ControlFlowGraph<Node>, Token, Token, Branch)
func cfaAssertNoEdgeOfType(
	t *testing.T, cfg ControlFlowGraph[rhino.Node], startToken rhino.Token, endToken rhino.Token, type_ ControlFlowGraph_Branch) {
	t.Helper()
	if edges := cfaGetAllEdgesOfType(cfg, startToken, endToken, type_); len(edges) != 0 {
		t.Errorf("expected no edge %v -%v-> %v, found %d", startToken, type_, endToken, len(edges))
	}
}

// cfaAssertDownEdge asserts that there exists a control flow edge of the
// given type from some node with the first token to some node with the
// second token. This edge must flow from a parent to one of its descendants.
//
// Java: ControlFlowAnalysisTest.assertDownEdge(ControlFlowGraph<Node>, Token, Token, Branch)
func cfaAssertDownEdge(
	t *testing.T, cfg ControlFlowGraph[rhino.Node], startToken rhino.Token, endToken rhino.Token, type_ ControlFlowGraph_Branch) {
	t.Helper()
	if len(cfaGetAllDownEdges(cfg, startToken, endToken, type_)) == 0 {
		t.Errorf("No down edge found (%v -%v-> %v)", startToken, type_, endToken)
	}
}

// cfaAssertEdge asserts that there exists a control flow edge of the given
// type from some node with the first token to some node with the second
// token.
//
// Java: ControlFlowAnalysisTest.assertEdge(ControlFlowGraph<Node>, Token, Token, Branch)
func cfaAssertEdge(
	t *testing.T, cfg ControlFlowGraph[rhino.Node], startToken rhino.Token, endToken rhino.Token, type_ ControlFlowGraph_Branch) {
	t.Helper()
	if len(cfaGetAllEdgesOfType(cfg /* startToken= */, startToken /* endToken= */, endToken, type_)) == 0 {
		t.Errorf("No up edge found. (%v -%v-> %v)", startToken, type_, endToken)
	}
}

// cfaAssertUpEdge asserts that there exists a control flow edge of the given
// type from some node with the first token to some node with the second
// token. This edge must flow from a node to one of its ancestors.
//
// Java: ControlFlowAnalysisTest.assertUpEdge(ControlFlowGraph<Node>, Token, Token, Branch)
func cfaAssertUpEdge(
	t *testing.T, cfg ControlFlowGraph[rhino.Node], startToken rhino.Token, endToken rhino.Token, type_ ControlFlowGraph_Branch) {
	t.Helper()
	if len(cfaGetAllDownEdges(cfg /* startToken= */, endToken /* endToken= */, startToken, type_)) == 0 {
		t.Errorf("No up edge found. (%v -%v-> %v)", startToken, type_, endToken)
	}
}

// cfaAssertCrossEdge asserts that there exists a control flow edge of the
// given type from some node with the first token to some node with the
// second token. This edge must flow between two nodes that are not in the
// same subtree.
//
// Java: ControlFlowAnalysisTest.assertCrossEdge(ControlFlowGraph<Node>, Token, Token, Branch)
func cfaAssertCrossEdge(
	t *testing.T, cfg ControlFlowGraph[rhino.Node], startToken rhino.Token, endToken rhino.Token, type_ ControlFlowGraph_Branch) {
	t.Helper()
	numDownEdges := len(cfaGetAllDownEdges(cfg, startToken, endToken, type_))
	numUpEdges := len(cfaGetAllDownEdges(cfg, endToken, startToken, type_))
	numEdges := len(cfaGetAllEdgesOfType(cfg, startToken, endToken, type_))
	if !(numDownEdges+numUpEdges < numEdges) {
		t.Errorf("No cross edges found (%v -%v-> %v)", startToken, type_, endToken)
	}
}

// cfaAssertReturnEdge asserts that there exists a control flow edge of the
// given type from some node with the first token to the return node.
//
// Java: ControlFlowAnalysisTest.assertReturnEdge(ControlFlowGraph<Node>, Token)
func cfaAssertReturnEdge(t *testing.T, cfg ControlFlowGraph[rhino.Node], startToken rhino.Token) {
	t.Helper()
	edges := cfaGetAllEdges(cfg)
	for _, edge := range edges {
		source := edge.GetSource().GetValue()
		dest := edge.GetDestination()
		if source.GetToken() == startToken && cfg.IsImplicitReturn(dest) {
			return
		}
	}

	t.Errorf("No return edge found (%v)", startToken)
}

// cfaAssertNoReturnEdge asserts that there exists no control flow edge of
// the given type from some node with the first token to the return node.
//
// Java: ControlFlowAnalysisTest.assertNoReturnEdge(ControlFlowGraph<Node>, Token)
func cfaAssertNoReturnEdge(t *testing.T, cfg ControlFlowGraph[rhino.Node], startToken rhino.Token) {
	t.Helper()
	edges := cfaGetAllEdges(cfg)
	for _, edge := range edges {
		source := edge.GetSource().GetValue()
		dest := edge.GetDestination()
		if source.GetToken() == startToken {
			if cfg.IsImplicitReturn(dest) {
				t.Errorf("Token %v should not have an out going edge to the implicit return", startToken)
			}
			return
		}
	}
}

// cfaTestCreateCfg: given an input in JavaScript, get a control flow graph
// for it.
//
// Java: ControlFlowAnalysisTest.createCfg(String, boolean), createCfg(String)
func cfaTestCreateCfg(t *testing.T, input string, runSynBlockPass bool) ControlFlowGraph[rhino.Node] {
	t.Helper()
	compiler := NewCompiler()
	root := compiler.ParseSyntheticCode("cfgtest", input)
	if runSynBlockPass {
		pass := NewCreateSyntheticBlocks(compiler, "START", "END")
		pass.Process(rhino.NilNode, root)
	}
	return ControlFlowAnalysis_Builder_().
		SetCompiler(compiler).
		SetCfgRoot(root).
		SetTraverseFunctions(true).
		SetIncludeEdgeAnnotations(true).
		ComputeCfg()
}

// cfaFindNodes is CodeSubTree.findNodesNonEmpty (jscomp.testing): all the
// nodes matching predicate, in preorder (NodeUtil.visitPreOrder).
//
// Go-only: test-helper — CodeSubTree lives in com.google.javascript.jscomp.testing.
func cfaFindNodes(t *testing.T, rootNode rhino.Node, predicate func(rhino.Node) bool) []rhino.Node {
	t.Helper()
	var results []rhino.Node
	var visit func(n rhino.Node)
	visit = func(n rhino.Node) {
		if predicate(n) {
			results = append(results, n)
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			visit(c)
		}
	}
	visit(rootNode)
	if len(results) == 0 {
		t.Fatal("no nodes found")
	}
	return results
}

// cfaAssertNodeOrder asserts the priority order of CFG nodes.
//
// Checks that the node type of the highest-priority node matches the first
// element of the list, the type of the second node matches the second
// element of the list, and so on.
//
// Java: ControlFlowAnalysisTest.assertNodeOrder(ControlFlowGraph<Node>, List<Token>)
func cfaAssertNodeOrder(t *testing.T, cfg ControlFlowGraph[rhino.Node], nodeTypes []rhino.Token) {
	t.Helper()
	cfgNodes := cfg.GetDiGraphNodes()
	slices.SortStableFunc(cfgNodes, cfg.GetOptionalNodeComparator(true))

	// IMPLICIT RETURN must always be last.
	implicitReturn := cfgNodes[len(cfgNodes)-1].GetValue()
	cfgNodes = cfgNodes[:len(cfgNodes)-1]
	if !implicitReturn.IsNil() {
		t.Errorf("%s: expected the implicit return last", implicitReturn.ToStringTree())
	}

	var got []rhino.Token
	for _, n := range cfgNodes {
		got = append(got, n.GetValue().GetToken())
	}
	if !slices.Equal(got, nodeTypes) {
		t.Errorf("node order %v, want %v", got, nodeTypes)
	}
}

func TestControlFlowAnalysis_SimpleStatementsInScript(t *testing.T) {
	src := "var a; a = a; a = a"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_SCRIPT, rhino.Token_VAR, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_VAR, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_SimpleStatementsInGoogModule(t *testing.T) {
	src := "goog.module('myMod'); var a; a = a; a = a"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, ControlFlowGraph_Branch_UNCOND)
	// the EXPR_RESULT is the goog.module(...) call
	cfaAssertDownEdge(t, cfg, rhino.Token_MODULE_BODY, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_VAR, ControlFlowGraph_Branch_UNCOND) // goog.module() -> var a;
}

func TestControlFlowAnalysis_SimpleStatementsInEsModule(t *testing.T) {
	src := "var a; a = a; export default a;"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, ControlFlowGraph_Branch_UNCOND)
	cfaAssertDownEdge(t, cfg, rhino.Token_MODULE_BODY, rhino.Token_VAR, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_VAR, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_EXPORT, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_SimpleIf(t *testing.T) {
	src := "var x; if (x) { x() } else { x() };"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_SCRIPT, rhino.Token_VAR, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_VAR, rhino.Token_IF, ControlFlowGraph_Branch_UNCOND)
	cfaAssertDownEdge(t, cfg, rhino.Token_IF, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	cfaAssertNoEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_CALL)
	cfaAssertDownEdge(t, cfg, rhino.Token_IF, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_FALSE)
	cfaAssertReturnEdge(t, cfg, rhino.Token_EMPTY)
}

func TestControlFlowAnalysis_BreakingBlock(t *testing.T) {
	src := "X: { while(1) { break } }"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertUpEdge(t, cfg, rhino.Token_BREAK, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_BreakToStatement(t *testing.T) {
	src := "X: break X;"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertReturnEdge(t, cfg, rhino.Token_BREAK)
}

func TestControlFlowAnalysis_BreakWithMultipleLabels(t *testing.T) {
	src := "La: Lb: break La;"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertReturnEdge(t, cfg, rhino.Token_BREAK)
}

func TestControlFlowAnalysis_BreakingWhile(t *testing.T) {
	src := "var x; while(true) { break; } x();"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_WHILE, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_BREAK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_BREAK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_InifiteLoopWhile(t *testing.T) {
	src := "var x; while(true) { } x();"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_WHILE, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertNoEdge(t, cfg, rhino.Token_WHILE, rhino.Token_EXPR_RESULT)
}

func TestControlFlowAnalysis_InifiteLoopDoWhile(t *testing.T) {
	src := "var x; do { } while (true); x();"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_DO, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertNoEdge(t, cfg, rhino.Token_DO, rhino.Token_EXPR_RESULT)
}

func TestControlFlowAnalysis_InifiteLoopFor_emptyCond(t *testing.T) {
	src := "var x; for(;;) { } x();"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_FOR, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertNoEdge(t, cfg, rhino.Token_FOR, rhino.Token_EXPR_RESULT)
}

func TestControlFlowAnalysis_BreakingFor_emptyCond(t *testing.T) {
	src := "var x; for(;;) { break; } x();"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_FOR, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_BREAK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_BREAK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_InifiteLoopFor_trueCond(t *testing.T) {
	src := "var x; for(;true;) { } x();"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_FOR, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertNoEdge(t, cfg, rhino.Token_FOR, rhino.Token_EXPR_RESULT)
}

func TestControlFlowAnalysis_BreakingFor_trueCond(t *testing.T) {
	src := "var x; for(;true;) { break; } x();"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_FOR, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_BREAK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_BREAK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_ThrowInCatchBlock(t *testing.T) {
	src := "try { throw ''; } catch (e) { throw e;} finally {}"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="TRY"];
  node0 -> node1 [weight=1];
  node2 [label="BLOCK"];
  node1 -> node2 [weight=1];
  node3 [label="THROW"];
  node2 -> node3 [weight=1];
  node4 [label="STRINGLIT"];
  node3 -> node4 [weight=1];
  node5 [label="BLOCK"];
  node3 -> node5 [label="ON_EX", fontcolor="red", weight=0.01, color="red"];
  node2 -> node3 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node5 [weight=1];
  node6 [label="CATCH"];
  node5 -> node6 [weight=1];
  node7 [label="NAME(e)"];
  node6 -> node7 [weight=1];
  node8 [label="BLOCK"];
  node6 -> node8 [weight=1];
  node9 [label="THROW"];
  node8 -> node9 [weight=1];
  node10 [label="NAME(e)"];
  node9 -> node10 [weight=1];
  node11 [label="BLOCK"];
  node9 -> node11 [label="ON_EX", fontcolor="red", weight=0.01, color="red"];
  node8 -> node9 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node6 -> node8 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node5 -> node6 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node11 [weight=1];
  node11 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node2 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_BreakingTryBlock(t *testing.T) {
	src := "a: try { break a; } finally {} if(x) {}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_BREAK, rhino.Token_IF, ControlFlowGraph_Branch_UNCOND)

	src = "a: try {} finally {break a;} if(x) {}"
	cfg = cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_BREAK, rhino.Token_IF, ControlFlowGraph_Branch_UNCOND)

	src = "a: try {} catch(e) {break a;} if(x) {}"
	cfg = cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_BREAK, rhino.Token_IF, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_WithStatement(t *testing.T) {
	src := "var x, y; with(x) { y() }"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_WITH, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertNoEdge(t, cfg, rhino.Token_WITH, rhino.Token_NAME)
	cfaAssertNoEdge(t, cfg, rhino.Token_NAME, rhino.Token_BLOCK)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	cfaAssertReturnEdge(t, cfg, rhino.Token_EXPR_RESULT)
}

func TestControlFlowAnalysis_SimpleWhile(t *testing.T) {
	src := "var x; while (x) { x(); if (x) { break; } x() }"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_WHILE, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	cfaAssertDownEdge(t, cfg, rhino.Token_IF, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertReturnEdge(t, cfg, rhino.Token_BREAK)
}

func TestControlFlowAnalysis_SimpleSwitch(t *testing.T) {
	src := "var x; switch(x){ case(1): x(); case('x'): x(); break; default: x();}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_VAR, rhino.Token_SWITCH, ControlFlowGraph_Branch_UNCOND)
	cfaAssertNoEdge(t, cfg, rhino.Token_SWITCH, rhino.Token_NAME)
	cfaAssertDownEdge(t, cfg, rhino.Token_SWITCH, rhino.Token_SWITCH_BODY, ControlFlowGraph_Branch_UNCOND)
	// Transfer between cases and default.
	cfaAssertDownEdge(t, cfg, rhino.Token_SWITCH_BODY, rhino.Token_CASE, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_CASE, rhino.Token_CASE, ControlFlowGraph_Branch_ON_FALSE)
	cfaAssertCrossEdge(t, cfg, rhino.Token_CASE, rhino.Token_DEFAULT_CASE, ControlFlowGraph_Branch_ON_FALSE)
	// Within each case.
	cfaAssertDownEdge(t, cfg, rhino.Token_CASE, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	cfaAssertNoEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_CALL)
	cfaAssertNoEdge(t, cfg, rhino.Token_CALL, rhino.Token_NAME)
}

func TestControlFlowAnalysis_SimpleNoDefault(t *testing.T) {
	src := "var x; switch(x){ case(1): break; } x();"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_CASE, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_ON_FALSE)
}

func TestControlFlowAnalysis_SwitchDefaultFirst(t *testing.T) {
	// DEFAULT appears first. But it is should evaluated last.
	src := "var x; switch(x){ default: break; case 1: break; }"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_SWITCH_BODY, rhino.Token_CASE, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_CASE, rhino.Token_DEFAULT_CASE, ControlFlowGraph_Branch_ON_FALSE)
}

func TestControlFlowAnalysis_SwitchDefaultInMiddle(t *testing.T) {
	// DEFAULT appears in the middle. But it is should evaluated last.
	src := "var x; switch(x){ case 1: break; default: break; case 2: break; }"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_SWITCH_BODY, rhino.Token_CASE, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_CASE, rhino.Token_CASE, ControlFlowGraph_Branch_ON_FALSE)
	cfaAssertCrossEdge(t, cfg, rhino.Token_CASE, rhino.Token_DEFAULT_CASE, ControlFlowGraph_Branch_ON_FALSE)
}

func TestControlFlowAnalysis_SwitchEmpty(t *testing.T) {
	// DEFAULT appears first. But it is should evaluated last.
	src := "var x; switch(x){}; x()"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_SWITCH_BODY, rhino.Token_EMPTY, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_EMPTY, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_ReturnThrowingException(t *testing.T) {
	src := "function f() {try { return a(); } catch (e) {e()}}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_RETURN, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_CATCH, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_CallThrowingException(t *testing.T) {
	src := "function f() {try { obj.a(); } catch (e) {e()}}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_CATCH, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_OptChainCallThrowingException(t *testing.T) {
	src := "function f() {try { obj?.a(); } catch (e) {e()}}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_CATCH, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_OptChainGetPropThrowingException(t *testing.T) {
	src := "function f() {try { obj?.a; } catch (e) {e()}}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_CATCH, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_OptChainGetElemThrowingException(t *testing.T) {
	src := "function f() {try { obj?.[a]; } catch (e) {e()}}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_CATCH, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_ForOfThrowingException(t *testing.T) {
	src := "function f() {try { for (var x of y) {} } catch (e) {e()}}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_FOR_OF, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_CATCH, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_ForAwaitOfThrowingException(t *testing.T) {
	src := "async function f() {try { for await (var x of y) {} } catch (e) {e()}}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_FOR_AWAIT_OF, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_CATCH, ControlFlowGraph_Branch_UNCOND)
}

func TestControlFlowAnalysis_ForInNoThrowingException(t *testing.T) {
	src := "function f() {try { for (var x in y) {} } catch (e) {e()}}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertNoEdgeOfType(t, cfg, rhino.Token_FOR_IN, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
}

// Test a simple FOR loop.
func TestControlFlowAnalysis_SimpleFor(t *testing.T) {
	src := "var a; for (var x = 0; x < 100; x++) { a(); }"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="VAR"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(a)"];
  node1 -> node2 [weight=1];
  node3 [label="VAR"];
  node1 -> node3 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node4 [label="FOR"];
  node0 -> node4 [weight=1];
  node4 -> node3 [weight=1];
  node5 [label="NAME(x)"];
  node3 -> node5 [weight=1];
  node6 [label="NUMBER"];
  node5 -> node6 [weight=1];
  node3 -> node4 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node7 [label="LT"];
  node4 -> node7 [weight=1];
  node8 [label="NAME(x)"];
  node7 -> node8 [weight=1];
  node9 [label="NUMBER"];
  node7 -> node9 [weight=1];
  node10 [label="INC"];
  node4 -> node10 [weight=1];
  node11 [label="NAME(x)"];
  node10 -> node11 [weight=1];
  node10 -> node4 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node12 [label="BLOCK"];
  node4 -> node12 [weight=1];
  node13 [label="EXPR_RESULT"];
  node12 -> node13 [weight=1];
  node14 [label="CALL"];
  node13 -> node14 [weight=1];
  node15 [label="NAME(a)"];
  node14 -> node15 [weight=1];
  node13 -> node10 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node12 -> node13 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node4 -> RETURN [label="ON_FALSE", fontcolor="red", weight=0.01, color="red"];
  node4 -> node12 [label="ON_TRUE", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_SimpleForWithContinue(t *testing.T) {
	src := "var a; for (var x = 0; x < 100; x++) {a();continue;a()}"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="VAR"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(a)"];
  node1 -> node2 [weight=1];
  node3 [label="VAR"];
  node1 -> node3 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node4 [label="FOR"];
  node0 -> node4 [weight=1];
  node4 -> node3 [weight=1];
  node5 [label="NAME(x)"];
  node3 -> node5 [weight=1];
  node6 [label="NUMBER"];
  node5 -> node6 [weight=1];
  node3 -> node4 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node7 [label="LT"];
  node4 -> node7 [weight=1];
  node8 [label="NAME(x)"];
  node7 -> node8 [weight=1];
  node9 [label="NUMBER"];
  node7 -> node9 [weight=1];
  node10 [label="INC"];
  node4 -> node10 [weight=1];
  node11 [label="NAME(x)"];
  node10 -> node11 [weight=1];
  node10 -> node4 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node12 [label="BLOCK"];
  node4 -> node12 [weight=1];
  node13 [label="EXPR_RESULT"];
  node12 -> node13 [weight=1];
  node14 [label="CALL"];
  node13 -> node14 [weight=1];
  node15 [label="NAME(a)"];
  node14 -> node15 [weight=1];
  node16 [label="CONTINUE"];
  node13 -> node16 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node12 -> node16 [weight=1];
  node16 -> node10 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node17 [label="EXPR_RESULT"];
  node12 -> node17 [weight=1];
  node18 [label="CALL"];
  node17 -> node18 [weight=1];
  node19 [label="NAME(a)"];
  node18 -> node19 [weight=1];
  node17 -> node10 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node12 -> node13 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node4 -> RETURN [label="ON_FALSE", fontcolor="red", weight=0.01, color="red"];
  node4 -> node12 [label="ON_TRUE", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_NestedFor(t *testing.T) {
	// This is tricky as the inner FOR branches to "x++" ON_FALSE.
	src := "var a,b;a();for(var x=0;x<100;x++){for(var y=0;y<100;y++){continue;b();}}"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="VAR"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(a)"];
  node1 -> node2 [weight=1];
  node3 [label="NAME(b)"];
  node1 -> node3 [weight=1];
  node4 [label="EXPR_RESULT"];
  node1 -> node4 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node4 [weight=1];
  node5 [label="CALL"];
  node4 -> node5 [weight=1];
  node6 [label="NAME(a)"];
  node5 -> node6 [weight=1];
  node7 [label="VAR"];
  node4 -> node7 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node8 [label="FOR"];
  node0 -> node8 [weight=1];
  node8 -> node7 [weight=1];
  node9 [label="NAME(x)"];
  node7 -> node9 [weight=1];
  node10 [label="NUMBER"];
  node9 -> node10 [weight=1];
  node7 -> node8 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node11 [label="LT"];
  node8 -> node11 [weight=1];
  node12 [label="NAME(x)"];
  node11 -> node12 [weight=1];
  node13 [label="NUMBER"];
  node11 -> node13 [weight=1];
  node14 [label="INC"];
  node8 -> node14 [weight=1];
  node15 [label="NAME(x)"];
  node14 -> node15 [weight=1];
  node14 -> node8 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node16 [label="BLOCK"];
  node8 -> node16 [weight=1];
  node17 [label="FOR"];
  node16 -> node17 [weight=1];
  node18 [label="VAR"];
  node17 -> node18 [weight=1];
  node19 [label="NAME(y)"];
  node18 -> node19 [weight=1];
  node20 [label="NUMBER"];
  node19 -> node20 [weight=1];
  node18 -> node17 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node21 [label="LT"];
  node17 -> node21 [weight=1];
  node22 [label="NAME(y)"];
  node21 -> node22 [weight=1];
  node23 [label="NUMBER"];
  node21 -> node23 [weight=1];
  node24 [label="INC"];
  node17 -> node24 [weight=1];
  node25 [label="NAME(y)"];
  node24 -> node25 [weight=1];
  node24 -> node17 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node26 [label="BLOCK"];
  node17 -> node26 [weight=1];
  node27 [label="CONTINUE"];
  node26 -> node27 [weight=1];
  node27 -> node24 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node28 [label="EXPR_RESULT"];
  node26 -> node28 [weight=1];
  node29 [label="CALL"];
  node28 -> node29 [weight=1];
  node30 [label="NAME(b)"];
  node29 -> node30 [weight=1];
  node28 -> node24 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node26 -> node27 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node17 -> node14 [label="ON_FALSE", fontcolor="red", weight=0.01, color="red"];
  node17 -> node26 [label="ON_TRUE", fontcolor="red", weight=0.01, color="red"];
  node16 -> node18 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node8 -> RETURN [label="ON_FALSE", fontcolor="red", weight=0.01, color="red"];
  node8 -> node16 [label="ON_TRUE", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_NestedDoWithBreak(t *testing.T) {
	// The BREAK branches to a() with UNCOND.
	src := "var a;do{do{break}while(a);do{a()}while(a)}while(a);"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="VAR"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(a)"];
  node1 -> node2 [weight=1];
  node3 [label="BLOCK"];
  node1 -> node3 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node4 [label="DO"];
  node0 -> node4 [weight=1];
  node4 -> node3 [weight=1];
  node5 [label="DO"];
  node3 -> node5 [weight=1];
  node6 [label="BLOCK"];
  node5 -> node6 [weight=1];
  node7 [label="BREAK"];
  node6 -> node7 [weight=1];
  node8 [label="BLOCK"];
  node7 -> node8 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node6 -> node7 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node9 [label="NAME(a)"];
  node5 -> node9 [weight=1];
  node5 -> node6 [label="ON_TRUE", fontcolor="red", weight=0.01, color="red"];
  node5 -> node8 [label="ON_FALSE", fontcolor="red", weight=0.01, color="red"];
  node10 [label="DO"];
  node3 -> node10 [weight=1];
  node10 -> node8 [weight=1];
  node11 [label="EXPR_RESULT"];
  node8 -> node11 [weight=1];
  node12 [label="CALL"];
  node11 -> node12 [weight=1];
  node13 [label="NAME(a)"];
  node12 -> node13 [weight=1];
  node11 -> node10 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node8 -> node11 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node14 [label="NAME(a)"];
  node10 -> node14 [weight=1];
  node10 -> node4 [label="ON_FALSE", fontcolor="red", weight=0.01, color="red"];
  node10 -> node8 [label="ON_TRUE", fontcolor="red", weight=0.01, color="red"];
  node3 -> node6 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node15 [label="NAME(a)"];
  node4 -> node15 [weight=1];
  node4 -> RETURN [label="ON_FALSE", fontcolor="red", weight=0.01, color="red"];
  node4 -> node3 [label="ON_TRUE", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_ForIn(t *testing.T) {
	src := "var a,b;for(a in b){a()};"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="VAR"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(a)"];
  node1 -> node2 [weight=1];
  node3 [label="NAME(b)"];
  node1 -> node3 [weight=1];
  node4 [label="NAME(b)"];
  node1 -> node4 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node5 [label="FOR_IN"];
  node0 -> node5 [weight=1];
  node6 [label="NAME(a)"];
  node5 -> node6 [weight=1];
  node5 -> node4 [weight=1];
  node4 -> node5 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node7 [label="BLOCK"];
  node5 -> node7 [weight=1];
  node8 [label="EXPR_RESULT"];
  node7 -> node8 [weight=1];
  node9 [label="CALL"];
  node8 -> node9 [weight=1];
  node10 [label="NAME(a)"];
  node9 -> node10 [weight=1];
  node8 -> node5 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node7 -> node8 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node11 [label="EMPTY"];
  node5 -> node11 [label="ON_FALSE", fontcolor="red", weight=0.01, color="red"];
  node5 -> node7 [label="ON_TRUE", fontcolor="red", weight=0.01, color="red"];
  node0 -> node11 [weight=1];
  node11 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_Throw(t *testing.T) {
	src := "function f() { throw 1; f() }"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="FUNCTION"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(f)"];
  node1 -> node2 [weight=1];
  node3 [label="PARAM_LIST"];
  node1 -> node3 [weight=1];
  node4 [label="BLOCK"];
  node1 -> node4 [weight=1];
  node5 [label="THROW"];
  node4 -> node5 [weight=1];
  node6 [label="NUMBER"];
  node5 -> node6 [weight=1];
  node7 [label="EXPR_RESULT"];
  node4 -> node7 [weight=1];
  node8 [label="CALL"];
  node7 -> node8 [weight=1];
  node9 [label="NAME(f)"];
  node8 -> node9 [weight=1];
  node7 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node4 -> node5 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node4 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_SimpleFunction(t *testing.T) {
	src := "function f() { f() } f()"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="FUNCTION"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(f)"];
  node1 -> node2 [weight=1];
  node3 [label="PARAM_LIST"];
  node1 -> node3 [weight=1];
  node4 [label="BLOCK"];
  node1 -> node4 [weight=1];
  node5 [label="EXPR_RESULT"];
  node4 -> node5 [weight=1];
  node6 [label="CALL"];
  node5 -> node6 [weight=1];
  node7 [label="NAME(f)"];
  node6 -> node7 [weight=1];
  node5 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node4 -> node5 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node4 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node8 [label="EXPR_RESULT"];
  node0 -> node8 [weight=1];
  node9 [label="CALL"];
  node8 -> node9 [weight=1];
  node10 [label="NAME(f)"];
  node9 -> node10 [weight=1];
  node8 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node8 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_SimpleClass(t *testing.T) {
	src := "class C{} f();"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="CLASS"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(C)"];
  node1 -> node2 [weight=1];
  node3 [label="EMPTY"];
  node1 -> node3 [weight=1];
  node4 [label="CLASS_MEMBERS"];
  node1 -> node4 [weight=1];
  node5 [label="EXPR_RESULT"];
  node1 -> node5 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node5 [weight=1];
  node6 [label="CALL"];
  node5 -> node6 [weight=1];
  node7 [label="NAME(f)"];
  node6 -> node7 [weight=1];
  node5 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_Class_withPublicFieldsAndMethod(t *testing.T) {
	src := "class C{ x; y; foo() {}}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertEdge(t, cfg, rhino.Token_FUNCTION, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertReturnEdge(t, cfg, rhino.Token_BLOCK)
}

func TestControlFlowAnalysis_ClassWithMemberFunctions(t *testing.T) {
	src := "class C{ f(){} g(){} }"
	expectedWithoutShouldTraverseFunctions := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="CLASS"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(C)"];
  node1 -> node2 [weight=1];
  node3 [label="EMPTY"];
  node1 -> node3 [weight=1];
  node4 [label="CLASS_MEMBERS"];
  node1 -> node4 [weight=1];
  node5 [label="MEMBER_FUNCTION_DEF"];
  node4 -> node5 [weight=1];
  node6 [label="FUNCTION"];
  node5 -> node6 [weight=1];
  node7 [label="NAME"];
  node6 -> node7 [weight=1];
  node8 [label="PARAM_LIST"];
  node6 -> node8 [weight=1];
  node9 [label="BLOCK"];
  node6 -> node9 [weight=1];
  node10 [label="MEMBER_FUNCTION_DEF"];
  node4 -> node10 [weight=1];
  node11 [label="FUNCTION"];
  node10 -> node11 [weight=1];
  node12 [label="NAME"];
  node11 -> node12 [weight=1];
  node13 [label="PARAM_LIST"];
  node11 -> node13 [weight=1];
  node14 [label="BLOCK"];
  node11 -> node14 [weight=1];
  node1 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	expectedWithShouldTraverseFunctions := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="CLASS"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(C)"];
  node1 -> node2 [weight=1];
  node3 [label="EMPTY"];
  node1 -> node3 [weight=1];
  node4 [label="CLASS_MEMBERS"];
  node1 -> node4 [weight=1];
  node5 [label="MEMBER_FUNCTION_DEF"];
  node4 -> node5 [weight=1];
  node6 [label="FUNCTION"];
  node5 -> node6 [weight=1];
  node7 [label="NAME"];
  node6 -> node7 [weight=1];
  node8 [label="PARAM_LIST"];
  node6 -> node8 [weight=1];
  node9 [label="BLOCK"];
  node6 -> node9 [weight=1];
  node9 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node6 -> node9 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node10 [label="MEMBER_FUNCTION_DEF"];
  node4 -> node10 [weight=1];
  node11 [label="FUNCTION"];
  node10 -> node11 [weight=1];
  node12 [label="NAME"];
  node11 -> node12 [weight=1];
  node13 [label="PARAM_LIST"];
  node11 -> node13 [weight=1];
  node14 [label="BLOCK"];
  node11 -> node14 [weight=1];
  node14 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node11 -> node14 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expectedWithShouldTraverseFunctions, true)
	cfaTestCfg(t, src, expectedWithoutShouldTraverseFunctions /* shouldTraverseFunctions= */, false)
}

func TestControlFlowAnalysis_SimpleCatch(t *testing.T) {
	src := "try{ throw x; x(); x['stuff']; x.x; x} catch (e) { e() }"

	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="TRY"];
  node0 -> node1 [weight=1];
  node2 [label="BLOCK"];
  node1 -> node2 [weight=1];
  node3 [label="THROW"];
  node2 -> node3 [weight=1];
  node4 [label="NAME(x)"];
  node3 -> node4 [weight=1];
  node5 [label="BLOCK"];
  node3 -> node5 [label="ON_EX", fontcolor="red", weight=0.01, color="red"];
  node6 [label="EXPR_RESULT"];
  node2 -> node6 [weight=1];
  node7 [label="CALL"];
  node6 -> node7 [weight=1];
  node8 [label="NAME(x)"];
  node7 -> node8 [weight=1];
  node9 [label="EXPR_RESULT"];
  node6 -> node5 [label="ON_EX", fontcolor="red", weight=0.01, color="red"];
  node6 -> node9 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node2 -> node9 [weight=1];
  node10 [label="GETELEM"];
  node9 -> node10 [weight=1];
  node11 [label="NAME(x)"];
  node10 -> node11 [weight=1];
  node12 [label="STRINGLIT(stuff)"];
  node10 -> node12 [weight=1];
  node13 [label="EXPR_RESULT"];
  node9 -> node13 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node9 -> node5 [label="ON_EX", fontcolor="red", weight=0.01, color="red"];
  node2 -> node13 [weight=1];
  node14 [label="GETPROP(x)"];
  node13 -> node14 [weight=1];
  node15 [label="NAME(x)"];
  node14 -> node15 [weight=1];
  node16 [label="EXPR_RESULT"];
  node13 -> node16 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node13 -> node5 [label="ON_EX", fontcolor="red", weight=0.01, color="red"];
  node2 -> node16 [weight=1];
  node17 [label="NAME(x)"];
  node16 -> node17 [weight=1];
  node16 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node2 -> node3 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node5 [weight=1];
  node18 [label="CATCH"];
  node5 -> node18 [weight=1];
  node19 [label="NAME(e)"];
  node18 -> node19 [weight=1];
  node20 [label="BLOCK"];
  node18 -> node20 [weight=1];
  node21 [label="EXPR_RESULT"];
  node20 -> node21 [weight=1];
  node22 [label="CALL"];
  node21 -> node22 [weight=1];
  node23 [label="NAME(e)"];
  node22 -> node23 [weight=1];
  node21 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node20 -> node21 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node18 -> node20 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node5 -> node18 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node2 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`

	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_FunctionWithinTry(t *testing.T) {
	// Make sure we don't search for the handler outside of the function.
	src := "try { var f = function() {throw 1;} } catch (e) { }"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="TRY"];
  node0 -> node1 [weight=1];
  node2 [label="BLOCK"];
  node1 -> node2 [weight=1];
  node3 [label="VAR"];
  node2 -> node3 [weight=1];
  node4 [label="NAME(f)"];
  node3 -> node4 [weight=1];
  node5 [label="FUNCTION"];
  node4 -> node5 [weight=1];
  node6 [label="NAME"];
  node5 -> node6 [weight=1];
  node7 [label="PARAM_LIST"];
  node5 -> node7 [weight=1];
  node8 [label="BLOCK"];
  node5 -> node8 [weight=1];
  node9 [label="THROW"];
  node8 -> node9 [weight=1];
  node10 [label="NUMBER"];
  node9 -> node10 [weight=1];
  node3 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node2 -> node3 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node11 [label="BLOCK"];
  node1 -> node11 [weight=1];
  node12 [label="CATCH"];
  node11 -> node12 [weight=1];
  node13 [label="NAME(e)"];
  node12 -> node13 [weight=1];
  node14 [label="BLOCK"];
  node12 -> node14 [weight=1];
  node14 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node12 -> node14 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node11 -> node12 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node2 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_NestedCatch(t *testing.T) {
	// Make sure we are going to the right handler.
	src := "try{try{throw 1;}catch(e){throw 2}}catch(f){}"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="TRY"];
  node0 -> node1 [weight=1];
  node2 [label="BLOCK"];
  node1 -> node2 [weight=1];
  node3 [label="TRY"];
  node2 -> node3 [weight=1];
  node4 [label="BLOCK"];
  node3 -> node4 [weight=1];
  node5 [label="THROW"];
  node4 -> node5 [weight=1];
  node6 [label="NUMBER"];
  node5 -> node6 [weight=1];
  node7 [label="BLOCK"];
  node5 -> node7 [label="ON_EX", fontcolor="red", weight=0.01, color="red"];
  node4 -> node5 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node3 -> node7 [weight=1];
  node8 [label="CATCH"];
  node7 -> node8 [weight=1];
  node9 [label="NAME(e)"];
  node8 -> node9 [weight=1];
  node10 [label="BLOCK"];
  node8 -> node10 [weight=1];
  node11 [label="THROW"];
  node10 -> node11 [weight=1];
  node12 [label="NUMBER"];
  node11 -> node12 [weight=1];
  node13 [label="BLOCK"];
  node11 -> node13 [label="ON_EX", fontcolor="red", weight=0.01, color="red"];
  node10 -> node11 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node8 -> node10 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node7 -> node8 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node3 -> node4 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node2 -> node3 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node13 [weight=1];
  node14 [label="CATCH"];
  node13 -> node14 [weight=1];
  node15 [label="NAME(f)"];
  node14 -> node15 [weight=1];
  node16 [label="BLOCK"];
  node14 -> node16 [weight=1];
  node16 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node14 -> node16 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node13 -> node14 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node2 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_SimpleFinally(t *testing.T) {
	src := "try{var x; foo()}finally{}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_TRY, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_VAR, ControlFlowGraph_Branch_UNCOND)
	// VAR to FINALLY.
	cfaAssertCrossEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	// No CATCH to FINALLY.
	cfaAssertNoEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_BLOCK)
}

func TestControlFlowAnalysis_SimpleCatchFinally(t *testing.T) {
	// Make sure we are going to the right handler.
	src := "try{ if(a){throw 1}else{a} } catch(e){a}finally{a}"
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="TRY"];
  node0 -> node1 [weight=1];
  node2 [label="BLOCK"];
  node1 -> node2 [weight=1];
  node3 [label="IF"];
  node2 -> node3 [weight=1];
  node4 [label="NAME(a)"];
  node3 -> node4 [weight=1];
  node5 [label="BLOCK"];
  node3 -> node5 [weight=1];
  node6 [label="THROW"];
  node5 -> node6 [weight=1];
  node7 [label="NUMBER"];
  node6 -> node7 [weight=1];
  node8 [label="BLOCK"];
  node6 -> node8 [label="ON_EX", fontcolor="red", weight=0.01, color="red"];
  node5 -> node6 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node9 [label="BLOCK"];
  node3 -> node9 [weight=1];
  node10 [label="EXPR_RESULT"];
  node9 -> node10 [weight=1];
  node11 [label="NAME(a)"];
  node10 -> node11 [weight=1];
  node12 [label="BLOCK"];
  node10 -> node12 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node9 -> node10 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node3 -> node5 [label="ON_TRUE", fontcolor="red", weight=0.01, color="red"];
  node3 -> node9 [label="ON_FALSE", fontcolor="red", weight=0.01, color="red"];
  node2 -> node3 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node8 [weight=1];
  node13 [label="CATCH"];
  node8 -> node13 [weight=1];
  node14 [label="NAME(e)"];
  node13 -> node14 [weight=1];
  node15 [label="BLOCK"];
  node13 -> node15 [weight=1];
  node16 [label="EXPR_RESULT"];
  node15 -> node16 [weight=1];
  node17 [label="NAME(a)"];
  node16 -> node17 [weight=1];
  node16 -> node12 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node15 -> node16 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node13 -> node15 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node8 -> node13 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node12 [weight=1];
  node18 [label="EXPR_RESULT"];
  node12 -> node18 [weight=1];
  node19 [label="NAME(a)"];
  node18 -> node19 [weight=1];
  node18 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node12 -> node18 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node1 -> node2 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_ComplicatedFinally2(t *testing.T) {
	// Now the most nasty case.....
	src := `while(1){try{if(a){a;continue;}else if(b){b;break;} else if(c) throw 1; else a}catch(e){}finally{c()}bar}foo`

	cfg := cfaTestCreateCfg(t, src, false)
	// Focus only on the ON_EX edges.
	cfaAssertCrossEdge(t, cfg, rhino.Token_CONTINUE, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_BREAK, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_THROW, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
}

func TestControlFlowAnalysis_DeepNestedBreakwithFinally(t *testing.T) {
	src := "X:while(1){try{while(2){try{var a;break X;}finally{}}}finally{}}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertDownEdge(t, cfg, rhino.Token_WHILE, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_TRUE)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_TRY, ControlFlowGraph_Branch_UNCOND)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_VAR, ControlFlowGraph_Branch_UNCOND)
	// BREAK to FINALLY.
	cfaAssertCrossEdge(t, cfg, rhino.Token_BREAK, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	// FINALLY to FINALLY.
	cfaAssertCrossEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
	cfaAssertCrossEdge(t, cfg, rhino.Token_WHILE, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_FALSE)
	cfaAssertReturnEdge(t, cfg, rhino.Token_BLOCK)
}

func TestControlFlowAnalysis_DeepNestedFinally(t *testing.T) {
	src := "try{try{try{throw 1}finally{1;var a}}finally{2;if(a);}}finally{3;a()}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_THROW, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
	cfaAssertCrossEdge(t, cfg, rhino.Token_VAR, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_IF, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
}

func TestControlFlowAnalysis_Return(t *testing.T) {
	src := "function f() { return; }"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertReturnEdge(t, cfg, rhino.Token_RETURN)
}

func TestControlFlowAnalysis_ReturnInFinally(t *testing.T) {
	src := "function f(x){ try{} finally {return x;} }"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertReturnEdge(t, cfg, rhino.Token_RETURN)
}

func TestControlFlowAnalysis_ReturnInFinally2(t *testing.T) {
	src := "function f(x){ try{ try{}finally{var dummy; return x;} } finally {} }"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_VAR, rhino.Token_RETURN, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_RETURN, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertReturnEdge(t, cfg, rhino.Token_BLOCK)
	cfaAssertNoReturnEdge(t, cfg, rhino.Token_RETURN)
}

func TestControlFlowAnalysis_ReturnInTry(t *testing.T) {
	src := "function f(x){ try{x; return x()} finally {} var y;}"
	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_RETURN, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_RETURN, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertCrossEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_VAR, ControlFlowGraph_Branch_UNCOND)
	cfaAssertReturnEdge(t, cfg, rhino.Token_VAR)
	cfaAssertReturnEdge(t, cfg, rhino.Token_BLOCK)
	cfaAssertNoReturnEdge(t, cfg, rhino.Token_RETURN)
}

func TestControlFlowAnalysis_OptionNotToTraverseFunctions(t *testing.T) {
	src := "var x = 1; function f() { x = null; }"
	expectedWhenNotTraversingFunctions := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="VAR"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(x)"];
  node1 -> node2 [weight=1];
  node3 [label="NUMBER"];
  node2 -> node3 [weight=1];
  node1 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node4 [label="FUNCTION"];
  node0 -> node4 [weight=1];
  node5 [label="NAME(f)"];
  node4 -> node5 [weight=1];
  node6 [label="PARAM_LIST"];
  node4 -> node6 [weight=1];
  node7 [label="BLOCK"];
  node4 -> node7 [weight=1];
  node8 [label="EXPR_RESULT"];
  node7 -> node8 [weight=1];
  node9 [label="ASSIGN"];
  node8 -> node9 [weight=1];
  node10 [label="NAME(x)"];
  node9 -> node10 [weight=1];
  node11 [label="NULL"];
  node9 -> node11 [weight=1];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="VAR"];
  node0 -> node1 [weight=1];
  node2 [label="NAME(x)"];
  node1 -> node2 [weight=1];
  node3 [label="NUMBER"];
  node2 -> node3 [weight=1];
  node1 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node4 [label="FUNCTION"];
  node0 -> node4 [weight=1];
  node5 [label="NAME(f)"];
  node4 -> node5 [weight=1];
  node6 [label="PARAM_LIST"];
  node4 -> node6 [weight=1];
  node7 [label="BLOCK"];
  node4 -> node7 [weight=1];
  node8 [label="EXPR_RESULT"];
  node7 -> node8 [weight=1];
  node9 [label="ASSIGN"];
  node8 -> node9 [weight=1];
  node10 [label="NAME(x)"];
  node9 -> node10 [weight=1];
  node11 [label="NULL"];
  node9 -> node11 [weight=1];
  node8 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node7 -> node8 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node4 -> node7 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
	cfaTestCfg(t, src, expectedWhenNotTraversingFunctions, false)
}

func TestControlFlowAnalysis_InstanceOf(t *testing.T) {
	src := "try { x instanceof 'x' } catch (e) { }"
	cfg := cfaTestCreateCfg(t, src, true)
	cfaAssertCrossEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_BLOCK, ControlFlowGraph_Branch_ON_EX)
}

func TestControlFlowAnalysis_SynBlock(t *testing.T) {
	src := "START(); var x; END(); var y;"
	cfg := cfaTestCreateCfg(t, src, true)
	cfaAssertCrossEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_SYN_BLOCK)
}

func TestControlFlowAnalysis_PartialTraversalOfScope(t *testing.T) {
	compiler := NewCompiler()

	script1 := compiler.ParseSyntheticCode("cfgtest", "var foo;")
	script2 := compiler.ParseSyntheticCode("cfgtest2", "var bar;")
	// Create a parent node for the scripts
	rhino.NewNodeWith2Children(rhino.Token_BLOCK, script1, script2)

	cfg := ControlFlowAnalysis_Builder_().
		SetCompiler(compiler).
		SetCfgRoot(script1).
		SetTraverseFunctions(true).
		SetIncludeEdgeAnnotations(true).
		ComputeCfg()

	if cfg.GetDiGraphNode(script1) == nil {
		t.Error("script1 should be in the CFG")
	}
	if cfg.GetDiGraphNode(script2) != nil {
		t.Error("script2 should not be in the CFG")
	}
}

func TestControlFlowAnalysis_ForLoopOrder(t *testing.T) {
	cfaAssertNodeOrder(t,
		cfaTestCreateCfg(t, "for (var i = 0; i < 5; i++) { var x = 3; } if (true) {}", false),
		[]rhino.Token{
			rhino.Token_SCRIPT,
			rhino.Token_VAR,
			rhino.Token_FOR,
			rhino.Token_BLOCK,
			rhino.Token_VAR,
			rhino.Token_INC, /* i++ */
			rhino.Token_IF,
			rhino.Token_BLOCK,
		})
}

func TestControlFlowAnalysis_LabelledForInLoopOrder(t *testing.T) {
	cfaAssertNodeOrder(t,
		cfaTestCreateCfg(t,
			"var i = 0; var y = {}; "+
				"label: for (var x in y) { "+
				"    if (x) { break label; } else { i++ } x(); }", false),
		[]rhino.Token{
			rhino.Token_SCRIPT,
			rhino.Token_VAR,
			rhino.Token_VAR,
			rhino.Token_NAME,
			rhino.Token_FOR_IN,
			rhino.Token_BLOCK,
			rhino.Token_IF,
			rhino.Token_BLOCK,
			rhino.Token_BREAK,
			rhino.Token_BLOCK,
			rhino.Token_EXPR_RESULT,
			rhino.Token_EXPR_RESULT,
		})
}

func TestControlFlowAnalysis_SimpleExportDeclarationsInEsModule(t *testing.T) {
	cfaAssertNodeOrder(t,
		cfaTestCreateCfg(t, "export let a = 0; export default a; export {b} from './mod';", false),
		[]rhino.Token{
			rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_EXPORT, rhino.Token_EXPORT, rhino.Token_EXPORT,
		})
}

func TestControlFlowAnalysis_SimpleImportDeclarationsInEsModule(t *testing.T) {
	cfaAssertNodeOrder(t,
		cfaTestCreateCfg(t, "import x from './mod'; import {y} from './mod'; import * as z from './mod';", false),
		[]rhino.Token{
			rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_IMPORT, rhino.Token_IMPORT, rhino.Token_IMPORT,
		})
}

func TestControlFlowAnalysis_LocalFunctionOrder(t *testing.T) {
	cfg := cfaTestCreateCfg(t, "function f() { while (x) { x++; } } var x = 3;", false)
	cfaAssertNodeOrder(t,
		cfg,
		[]rhino.Token{
			rhino.Token_SCRIPT,
			rhino.Token_VAR,
			rhino.Token_FUNCTION,
			rhino.Token_BLOCK,
			rhino.Token_WHILE,
			rhino.Token_BLOCK,
			rhino.Token_EXPR_RESULT,
		})
}

func TestControlFlowAnalysis_DoWhileOrder(t *testing.T) {
	cfaAssertNodeOrder(t,
		cfaTestCreateCfg(t, "do { var x = 3; } while (true); void x;", false),
		[]rhino.Token{rhino.Token_SCRIPT, rhino.Token_BLOCK, rhino.Token_VAR, rhino.Token_DO, rhino.Token_EXPR_RESULT})
}

func TestControlFlowAnalysis_ForOfOrder(t *testing.T) {
	cfaAssertNodeOrder(t,
		cfaTestCreateCfg(t, "async function f() { for (x of y) { z; } return 0; }", false),
		[]rhino.Token{
			rhino.Token_SCRIPT,
			rhino.Token_FUNCTION,
			rhino.Token_BLOCK,
			rhino.Token_NAME,
			rhino.Token_FOR_OF,
			rhino.Token_BLOCK,
			rhino.Token_EXPR_RESULT,
			rhino.Token_RETURN,
		})
}

func TestControlFlowAnalysis_ForAwaitOfOrder(t *testing.T) {
	cfaAssertNodeOrder(t,
		cfaTestCreateCfg(t, "async function f() { for await (x of y) { z; } return 0; }", false),
		[]rhino.Token{
			rhino.Token_SCRIPT,
			rhino.Token_FUNCTION,
			rhino.Token_BLOCK,
			rhino.Token_NAME,
			rhino.Token_FOR_AWAIT_OF,
			rhino.Token_BLOCK,
			rhino.Token_EXPR_RESULT,
			rhino.Token_RETURN,
		})
}

func TestControlFlowAnalysis_ForAwaitOfOrderBreakAndContinue(t *testing.T) {
	cfaAssertNodeOrder(t,
		cfaTestCreateCfg(t, `async function f() {
  outer: for await (let x of y) {
    inner: for await (let z of x) {
      if (z) break inner;
      else continue outer;
    }
  }
  return 0;
}
`, false),
		[]rhino.Token{
			rhino.Token_SCRIPT,
			rhino.Token_FUNCTION,
			rhino.Token_BLOCK,
			rhino.Token_NAME,
			rhino.Token_FOR_AWAIT_OF,
			rhino.Token_BLOCK,
			rhino.Token_NAME,
			rhino.Token_FOR_AWAIT_OF,
			rhino.Token_BLOCK,
			rhino.Token_IF,
			rhino.Token_BLOCK,
			rhino.Token_BREAK,
			rhino.Token_BLOCK,
			rhino.Token_CONTINUE,
			rhino.Token_RETURN,
		})
}

func TestControlFlowAnalysis_ForAwaitOfOrderBreakAndContinueAndYield(t *testing.T) {
	cfaAssertNodeOrder(t,
		cfaTestCreateCfg(t, `async function* f() {
  outer: for await (let x of y) {
    inner: for await (let z of x) {
      if (z > 0) break inner;
      else if (z < 0) continue outer;
      yield z;
    }
  }
}
`, false),
		[]rhino.Token{
			rhino.Token_SCRIPT,
			rhino.Token_FUNCTION,
			rhino.Token_BLOCK,
			rhino.Token_NAME,
			rhino.Token_FOR_AWAIT_OF,
			rhino.Token_BLOCK,
			rhino.Token_NAME,
			rhino.Token_FOR_AWAIT_OF,
			rhino.Token_BLOCK,
			rhino.Token_IF,
			rhino.Token_BLOCK,
			rhino.Token_BREAK,
			rhino.Token_BLOCK,
			rhino.Token_IF,
			rhino.Token_BLOCK,
			rhino.Token_CONTINUE,
			rhino.Token_EXPR_RESULT,
		})
}

func TestControlFlowAnalysis_BreakInFinally1(t *testing.T) {
	src := `f = function() {
  var action;
  a: {
    var proto = null;
    try {
      proto = new Proto
    } finally {
      action = proto;
      break a // Remove this...
    }
  }
  alert(action) // but not this.
};
`
	expected := `digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="SCRIPT"];
  node1 [label="EXPR_RESULT"];
  node0 -> node1 [weight=1];
  node2 [label="ASSIGN"];
  node1 -> node2 [weight=1];
  node3 [label="NAME(f)"];
  node2 -> node3 [weight=1];
  node4 [label="FUNCTION"];
  node2 -> node4 [weight=1];
  node5 [label="NAME"];
  node4 -> node5 [weight=1];
  node6 [label="PARAM_LIST"];
  node4 -> node6 [weight=1];
  node7 [label="BLOCK"];
  node4 -> node7 [weight=1];
  node8 [label="VAR"];
  node7 -> node8 [weight=1];
  node9 [label="NAME(action)"];
  node8 -> node9 [weight=1];
  node10 [label="LABEL"];
  node7 -> node10 [weight=1];
  node11 [label="LABEL_NAME"];
  node10 -> node11 [weight=1];
  node12 [label="BLOCK"];
  node10 -> node12 [weight=1];
  node13 [label="VAR"];
  node12 -> node13 [weight=1];
  node14 [label="NAME(proto)"];
  node13 -> node14 [weight=1];
  node15 [label="NULL"];
  node14 -> node15 [weight=1];
  node16 [label="TRY"];
  node12 -> node16 [weight=1];
  node17 [label="BLOCK"];
  node16 -> node17 [weight=1];
  node18 [label="EXPR_RESULT"];
  node17 -> node18 [weight=1];
  node19 [label="ASSIGN"];
  node18 -> node19 [weight=1];
  node20 [label="NAME(proto)"];
  node19 -> node20 [weight=1];
  node21 [label="NEW"];
  node19 -> node21 [weight=1];
  node22 [label="NAME(Proto)"];
  node21 -> node22 [weight=1];
  node23 [label="BLOCK"];
  node16 -> node23 [weight=1];
  node24 [label="BLOCK"];
  node16 -> node24 [weight=1];
  node25 [label="EXPR_RESULT"];
  node24 -> node25 [weight=1];
  node26 [label="ASSIGN"];
  node25 -> node26 [weight=1];
  node27 [label="NAME(action)"];
  node26 -> node27 [weight=1];
  node28 [label="NAME(proto)"];
  node26 -> node28 [weight=1];
  node29 [label="BREAK"];
  node24 -> node29 [weight=1];
  node30 [label="LABEL_NAME"];
  node29 -> node30 [weight=1];
  node31 [label="EXPR_RESULT"];
  node7 -> node31 [weight=1];
  node32 [label="CALL"];
  node31 -> node32 [weight=1];
  node33 [label="NAME(alert)"];
  node32 -> node33 [weight=1];
  node34 [label="NAME(action)"];
  node32 -> node34 [weight=1];
  node1 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}
`
	cfaTestCfg(t, src, expected, true)
}

func TestControlFlowAnalysis_BreakInFinally2(t *testing.T) {
	src := `var action;
a: {
  var proto = null;
  try {
    proto = new Proto
  } finally {
    action = proto;
    break a
  }
}
alert(action)
`

	cfg := cfaTestCreateCfg(t, src, false)
	cfaAssertCrossEdge(t, cfg, rhino.Token_BREAK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	cfaAssertNoEdge(t, cfg, rhino.Token_BREAK, rhino.Token_BLOCK)
}

func TestControlFlowAnalysis_TryFinallyDeterminism(t *testing.T) {
	src := `a: {
  b: {
    try {
      if (x) { break a; } else { break b; }
    } finally {
      fin();
    }
  }
  afterB();
}
afterA();
`
	cfg := cfaTestCreateCfg(t, src, false)

	var finExpr *LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	for _, n := range cfg.GetLinkedDiGraphNodes() {
		if !n.GetValue().IsNil() &&
			n.GetValue().IsExprResult() &&
			n.GetValue().GetFirstChild().GetFirstChild().MatchesName("fin") {
			finExpr = n
			break
		}
	}
	if finExpr == nil {
		t.Fatal("no fin() node")
	}

	var onExTargets []string
	for _, e := range cfg.GetOutEdges(finExpr.GetValue()) {
		if e.GetValue() == ControlFlowGraph_Branch_ON_EX && !e.GetDestination().GetValue().IsNil() {
			onExTargets = append(onExTargets, e.GetDestination().GetValue().GetFirstChild().GetFirstChild().GetString())
		}
	}

	if want := []string{"afterA", "afterB"}; !slices.Equal(onExTargets, want) {
		t.Errorf("ON_EX targets %v, want %v", onExTargets, want)
	}
}

func TestControlFlowAnalysis_CfgRootedAtEmptyFunctionDeclaration(t *testing.T) {
	src := "function fn() {}"
	compiler := NewCompiler()
	globalRoot := compiler.ParseSyntheticCode("cfgtest", src)
	fnRoot := cfaFindNodes(t, globalRoot, (rhino.Node).IsFunction)[0]

	cfg := ControlFlowAnalysis_Builder_().SetCompiler(compiler).SetCfgRoot(fnRoot).ComputeCfg()

	cfaAssertDownEdge(t, cfg, rhino.Token_FUNCTION, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertReturnEdge(t, cfg, rhino.Token_BLOCK)
}

func TestControlFlowAnalysis_CfgRootedAtFunctionDeclarationWithSingleStatement(t *testing.T) {
	src := "function fn() {  alert(3); }"
	compiler := NewCompiler()
	globalRoot := compiler.ParseSyntheticCode("cfgtest", src)
	fnRoot := cfaFindNodes(t, globalRoot, (rhino.Node).IsFunction)[0]

	cfg := ControlFlowAnalysis_Builder_().SetCompiler(compiler).SetCfgRoot(fnRoot).ComputeCfg()

	cfaAssertDownEdge(t, cfg, rhino.Token_FUNCTION, rhino.Token_BLOCK, ControlFlowGraph_Branch_UNCOND)
	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	cfaAssertReturnEdge(t, cfg, rhino.Token_EXPR_RESULT)
}

func TestControlFlowAnalysis_CfgRootedAtEmptyClassStaticBlock(t *testing.T) {
	src := "class C { static {} }"
	compiler := NewCompiler()
	globalRoot := compiler.ParseSyntheticCode("cfgtest", src)
	staticBlock := cfaFindNodes(t, globalRoot, (rhino.Node).IsBlock)[0]

	cfg := ControlFlowAnalysis_Builder_().SetCompiler(compiler).SetCfgRoot(staticBlock).ComputeCfg()

	cfaAssertReturnEdge(t, cfg, rhino.Token_BLOCK)
}

func TestControlFlowAnalysis_CfgRootedAtClassStaticBlockSingleStatementAndSubsequentMembers(t *testing.T) {
	src := "class C { static { alert(0); } x = 0; fn() {} }"
	compiler := NewCompiler()
	globalRoot := compiler.ParseSyntheticCode("cfgtest", src)
	staticBlock := cfaFindNodes(t, globalRoot, (rhino.Node).IsBlock)[0]

	cfg := ControlFlowAnalysis_Builder_().SetCompiler(compiler).SetCfgRoot(staticBlock).ComputeCfg()

	cfaAssertDownEdge(t, cfg, rhino.Token_BLOCK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	// Double check edge removal from EXPR_RESULT to next child of CLASS_MEMBER
	cfaAssertNoEdge(t, cfg, rhino.Token_EXPR_RESULT, rhino.Token_MEMBER_FIELD_DEF)
	cfaAssertReturnEdge(t, cfg, rhino.Token_EXPR_RESULT)
}

func TestControlFlowAnalysis_CfgRootedAtClassStaticBlockMultipleStatementAndSubsequentMembers(t *testing.T) {
	src := `class C {
  static { alert(0); }
  x = 0;
  [0+1]() {}
  static { alert(1); }
  fn(){}
}
`
	compiler := NewCompiler()
	globalRoot := compiler.ParseSyntheticCode("cfgtest", src)
	staticBlocks := cfaFindNodes(t, globalRoot, NodeUtil_IsClassStaticBlock)
	if len(staticBlocks) != 2 {
		t.Fatalf("found %d static blocks, want 2", len(staticBlocks))
	}
	staticBlock1 := staticBlocks[0]
	staticBlock2 := staticBlocks[1]

	// First Static Block
	cfg1 := ControlFlowAnalysis_Builder_().SetCompiler(compiler).SetCfgRoot(staticBlock1).ComputeCfg()

	cfaAssertDownEdge(t, cfg1, rhino.Token_BLOCK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	cfaAssertNoEdge(t, cfg1, rhino.Token_EXPR_RESULT, rhino.Token_MEMBER_FIELD_DEF)
	cfaAssertReturnEdge(t, cfg1, rhino.Token_EXPR_RESULT)

	cfaAssertNoEdge(t, cfg1, rhino.Token_COMPUTED_PROP, rhino.Token_BLOCK)

	// Second Static Block
	cfg2 := ControlFlowAnalysis_Builder_().SetCompiler(compiler).SetCfgRoot(staticBlock2).ComputeCfg()

	cfaAssertDownEdge(t, cfg2, rhino.Token_BLOCK, rhino.Token_EXPR_RESULT, ControlFlowGraph_Branch_UNCOND)
	cfaAssertNoEdge(t, cfg2, rhino.Token_EXPR_RESULT, rhino.Token_MEMBER_FIELD_DEF)
	cfaAssertReturnEdge(t, cfg2, rhino.Token_EXPR_RESULT)
}

func TestControlFlowAnalysis_IsEnteringNewCfgNode_enhancedForLoops(t *testing.T) {
	compiler := NewCompiler()

	check := func(root rhino.Node, predicate func(rhino.Node) bool) {
		t.Helper()
		loop := cfaFindNodes(t, root, predicate)[0]
		item := loop.GetFirstChild()
		collection := item.GetNext()
		body := collection.GetNext()

		if ControlFlowGraph_IsEnteringNewCfgNode(item) {
			t.Errorf("%s: item is entering a new CFG node", loop.GetToken())
		}
		if !ControlFlowGraph_IsEnteringNewCfgNode(collection) {
			t.Errorf("%s: collection is not entering a new CFG node", loop.GetToken())
		}
		if !ControlFlowGraph_IsEnteringNewCfgNode(body) {
			t.Errorf("%s: body is not entering a new CFG node", loop.GetToken())
		}
	}

	// FOR_IN
	check(compiler.ParseSyntheticCode("for_in", "for (var item in collection) { body(); }"), (rhino.Node).IsForIn)

	// FOR_OF
	check(compiler.ParseSyntheticCode("for_of", "for (var item of collection) { body(); }"), (rhino.Node).IsForOf)

	// FOR_AWAIT_OF
	check(
		compiler.ParseSyntheticCode(
			"for_await_of",
			"async function f() { for await (var item of collection) { body(); } }"),
		(rhino.Node).IsForAwaitOf)
}
