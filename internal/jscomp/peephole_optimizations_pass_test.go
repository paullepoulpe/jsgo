// Port of com/google/javascript/jscomp/PeepholeOptimizationsPassTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/PeepholeOptimizationsPassTest.java)

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PeepholeOptimizationsPassTest is the unit tests for
// PeepholeOptimizationsPass.
type PeepholeOptimizationsPassTest struct {
	CompilerTestCaseBase

	currentPeepholePasses []AbstractPeepholeOptimization
}

// Java: PeepholeOptimizationsPassTest()
func NewPeepholeOptimizationsPassTest() *PeepholeOptimizationsPassTest {
	c := &PeepholeOptimizationsPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestPeepholeOptimizationsPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPeepholeOptimizationsPassTest() })
}

// Java: PeepholeOptimizationsPassTest.getProcessor(Compiler)
func (c *PeepholeOptimizationsPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPeepholeOptimizationsPass(
		compiler, c.GetName(), append([]AbstractPeepholeOptimization(nil), c.currentPeepholePasses...)...)
}

// peepholeOptimizationsPassTest_anonymous is an anonymous subclass of
// AbstractPeepholeOptimization (`new AbstractPeepholeOptimization() {...}`);
// its optimizeSubtree override is the closure.
type peepholeOptimizationsPassTest_anonymous struct {
	AbstractPeepholeOptimizationBase
	optimizeSubtree func(self *peepholeOptimizationsPassTest_anonymous, node rhino.Node) rhino.Node
}

func newPeepholeOptimizationsPassTest_anonymous(
	optimizeSubtree func(self *peepholeOptimizationsPassTest_anonymous, node rhino.Node) rhino.Node) *peepholeOptimizationsPassTest_anonymous {
	o := &peepholeOptimizationsPassTest_anonymous{optimizeSubtree: optimizeSubtree}
	return o
}

// OptimizeSubtree is the anonymous class's optimizeSubtree(Node).
func (o *peepholeOptimizationsPassTest_anonymous) OptimizeSubtree(node rhino.Node) rhino.Node {
	return o.optimizeSubtree(o, node)
}

// PeepholeOptimizationsPass should handle the case when no peephole
// optimizations are turned on.
//
// Java: PeepholeOptimizationsPassTest.testEmptyPass()
func (c *PeepholeOptimizationsPassTest) TestEmptyPass() {
	c.currentPeepholePasses = []AbstractPeepholeOptimization{}

	c.TestSame("var x; var y;")
}

// Java: PeepholeOptimizationsPassTest.testOptimizationOrder()
func (c *PeepholeOptimizationsPassTest) TestOptimizationOrder() {
	// We need to make sure that: 1) We are only traversing the AST once 2) For
	// each node, we visit the optimizations in the client-supplied order
	//
	// To test this, we create two fake optimizations that each make an entry in
	// the visitationLog when they are passed a name node to optimize.
	//
	// Each entry is of the form nameX where 'name' is the name of the name node
	// visited and X is the identity of the optimization (1 or 2 in this case).
	// After the pass is run, we verify the correct ordering by querying the
	// log.
	//
	// Using a log, rather than, say, transforming nodes, allows us to ensure
	// not only that we are visiting each node but that our visits occur in the
	// right order (i.e. we need to make sure we're not traversing the entire
	// AST for the first optimization and then a second time for the second).

	visitationLog := jutil.NewArrayList[string]()

	note1Applied :=
		newPeepholeOptimizationsPassTest_anonymous(func(self *peepholeOptimizationsPassTest_anonymous, node rhino.Node) rhino.Node {
			if node.IsName() {
				visitationLog.Add(node.GetString() + "1")
			}

			return node
		})

	note2Applied :=
		newPeepholeOptimizationsPassTest_anonymous(func(self *peepholeOptimizationsPassTest_anonymous, node rhino.Node) rhino.Node {
			if node.IsName() {
				visitationLog.Add(node.GetString() + "2")
			}

			return node
		})

	c.currentPeepholePasses = []AbstractPeepholeOptimization{note1Applied, note2Applied}

	c.TestSame("var x; var y")

	// We expect the optimization order to be: "x" visited by optimization1 "x"
	// visited by optimization2 "y" visited by optimization1 "y" visited by
	// optimization2
	AssertThat(visitationLog).ContainsExactly("x1", "x2", "y1", "y2").InOrder()
}

// peepholeOptimizationsPassTest_RemoveNodesNamedXUnderVarOptimization is a
// peephole optimization that, given a subtree consisting of a VAR node,
// removes children of that node named "x".
type peepholeOptimizationsPassTest_RemoveNodesNamedXUnderVarOptimization struct {
	AbstractPeepholeOptimizationBase
}

func newPeepholeOptimizationsPassTest_RemoveNodesNamedXUnderVarOptimization() *peepholeOptimizationsPassTest_RemoveNodesNamedXUnderVarOptimization {
	o := &peepholeOptimizationsPassTest_RemoveNodesNamedXUnderVarOptimization{}
	return o
}

// Java: PeepholeOptimizationsPassTest.RemoveNodesNamedXUnderVarOptimization.optimizeSubtree(Node)
func (o *peepholeOptimizationsPassTest_RemoveNodesNamedXUnderVarOptimization) OptimizeSubtree(node rhino.Node) rhino.Node {
	if node.IsVar() {
		nodesToRemove := jutil.NewLinkedHashSet[rhino.Node]()

		for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if "x" == child.GetString() {
				nodesToRemove.Add(child)
			}
		}

		for childToRemove := range nodesToRemove.All() {
			o.ReportChangeToEnclosingScope(node)
			childToRemove.Detach()
		}
	}

	return node
}

// peepholeOptimizationsPassTest_RemoveNodesNamedXOptimization is a peephole
// optimization that, given a subtree consisting of a name node named "x"
// removes that node.
type peepholeOptimizationsPassTest_RemoveNodesNamedXOptimization struct {
	AbstractPeepholeOptimizationBase
}

func newPeepholeOptimizationsPassTest_RemoveNodesNamedXOptimization() *peepholeOptimizationsPassTest_RemoveNodesNamedXOptimization {
	o := &peepholeOptimizationsPassTest_RemoveNodesNamedXOptimization{}
	return o
}

// Java: PeepholeOptimizationsPassTest.RemoveNodesNamedXOptimization.optimizeSubtree(Node)
func (o *peepholeOptimizationsPassTest_RemoveNodesNamedXOptimization) OptimizeSubtree(node rhino.Node) rhino.Node {
	if node.IsName() && "x" == node.GetString() {
		o.ReportChangeToEnclosingScope(node)
		node.Detach()

		return rhino.NilNode
	}

	return node
}

// peepholeOptimizationsPassTest_RemoveParentVarsForNodesNamedX is a peephole
// optimization that, given a subtree consisting of a name node named "x"
// whose parent is a VAR node, removes the parent VAR node.
type peepholeOptimizationsPassTest_RemoveParentVarsForNodesNamedX struct {
	AbstractPeepholeOptimizationBase
}

func newPeepholeOptimizationsPassTest_RemoveParentVarsForNodesNamedX() *peepholeOptimizationsPassTest_RemoveParentVarsForNodesNamedX {
	o := &peepholeOptimizationsPassTest_RemoveParentVarsForNodesNamedX{}
	return o
}

// Java: PeepholeOptimizationsPassTest.RemoveParentVarsForNodesNamedX.optimizeSubtree(Node)
func (o *peepholeOptimizationsPassTest_RemoveParentVarsForNodesNamedX) OptimizeSubtree(node rhino.Node) rhino.Node {
	if node.IsName() && "x" == node.GetString() {
		parent := node.GetParent()
		if parent.IsVar() {
			o.ReportChangeToEnclosingScope(parent)
			parent.Detach()
			return rhino.NilNode
		}
	}
	return node
}

// peepholeOptimizationsPassTest_RenameYToX is a peephole optimization that,
// given a subtree consisting of a name node named "y", replaces it with a
// name node named "x";
type peepholeOptimizationsPassTest_RenameYToX struct {
	AbstractPeepholeOptimizationBase
}

func newPeepholeOptimizationsPassTest_RenameYToX() *peepholeOptimizationsPassTest_RenameYToX {
	o := &peepholeOptimizationsPassTest_RenameYToX{}
	return o
}

// Java: PeepholeOptimizationsPassTest.RenameYToX.optimizeSubtree(Node)
func (o *peepholeOptimizationsPassTest_RenameYToX) OptimizeSubtree(node rhino.Node) rhino.Node {
	if node.IsName() && "y" == node.GetString() {
		replacement := rhino.Node_NewStringWithToken(rhino.Token_NAME, "x")

		node.ReplaceWith(replacement)
		o.ReportChangeToEnclosingScope(replacement)

		return replacement
	}
	return node
}

// Java: PeepholeOptimizationsPassTest.testOptimizationRemovingSubtreeChild()
func (c *PeepholeOptimizationsPassTest) TestOptimizationRemovingSubtreeChild() {
	c.currentPeepholePasses =
		[]AbstractPeepholeOptimization{newPeepholeOptimizationsPassTest_RemoveNodesNamedXUnderVarOptimization()}

	c.Test("var x,y;", "var y;")
	c.Test("var y,x;", "var y;")
	c.Test("var x,y,x;", "var y;")
}

// Java: PeepholeOptimizationsPassTest.testOptimizationRemovingSubtree()
func (c *PeepholeOptimizationsPassTest) TestOptimizationRemovingSubtree() {
	c.currentPeepholePasses =
		[]AbstractPeepholeOptimization{newPeepholeOptimizationsPassTest_RemoveNodesNamedXOptimization()}

	c.Test("var x,y;", "var y;")
	c.Test("var y,x;", "var y;")
	c.Test("var x,y,x;", "var y;")
}

// Java: PeepholeOptimizationsPassTest.testOptimizationRemovingSubtreeParent()
func (c *PeepholeOptimizationsPassTest) TestOptimizationRemovingSubtreeParent() {
	c.currentPeepholePasses =
		[]AbstractPeepholeOptimization{newPeepholeOptimizationsPassTest_RemoveParentVarsForNodesNamedX()}

	c.Test("var x; var y", "var y")
}

// Test the case where the first peephole optimization removes a node and the
// second wants to remove (the now nonexistent) parent of that node.
//
// Java: PeepholeOptimizationsPassTest.testOptimizationsRemoveParentAfterRemoveChild()
func (c *PeepholeOptimizationsPassTest) TestOptimizationsRemoveParentAfterRemoveChild() {
	c.currentPeepholePasses =
		[]AbstractPeepholeOptimization{newPeepholeOptimizationsPassTest_RemoveNodesNamedXOptimization(), newPeepholeOptimizationsPassTest_RemoveParentVarsForNodesNamedX()}

	c.Test("var x,y; var z;", "var y; var z;")
}

// Java: PeepholeOptimizationsPassTest.testOptimizationReplacingNode()
func (c *PeepholeOptimizationsPassTest) TestOptimizationReplacingNode() {
	c.currentPeepholePasses =
		[]AbstractPeepholeOptimization{newPeepholeOptimizationsPassTest_RenameYToX(), newPeepholeOptimizationsPassTest_RemoveParentVarsForNodesNamedX()}

	c.Test("var y; var z;", "var z;")
}

// Java: PeepholeOptimizationsPassTest.testAddFeatureToEnclosingScript()
func (c *PeepholeOptimizationsPassTest) TestAddFeatureToEnclosingScript() {
	c.currentPeepholePasses =
		[]AbstractPeepholeOptimization{
			newPeepholeOptimizationsPassTest_anonymous(func(self *peepholeOptimizationsPassTest_anonymous, node rhino.Node) rhino.Node {
				if node.IsAdd() {
					self.AddFeatureToEnclosingScript(parser.FeatureSet_Feature_LET_DECLARATIONS)
					self.AddFeatureToEnclosingScript(parser.FeatureSet_Feature_LET_DECLARATIONS)
					self.AddFeatureToEnclosingScript(parser.FeatureSet_Feature_CLASSES)
				}
				return node
			}),
			newPeepholeOptimizationsPassTest_anonymous(func(self *peepholeOptimizationsPassTest_anonymous, node rhino.Node) rhino.Node {
				if node.IsSub() {
					self.AddFeatureToEnclosingScript(parser.FeatureSet_Feature_CONST_DECLARATIONS)
				}
				return node
			}),
		}
	c.TestSame("(3 + 4); function sub() { return 3 - 4; }")
	compiler := c.GetLastCompiler()
	script := rhino.CheckNotNull(compiler.GetScriptNode("testcode"))
	AssertThat(script.GetProp(rhino.Node_FEATURE_SET)).
		IsEqualTo(
			parser.FeatureSet_BARE_MINIMUM.WithFeatures(
				parser.FeatureSet_Feature_LET_DECLARATIONS, parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_CONST_DECLARATIONS))
}
