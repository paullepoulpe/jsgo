// Port of com/google/javascript/jscomp/instrumentation/CoverageInstrumentationCallback.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CoverageInstrumentationCallback implements a traversal to instrument an AST
// for code coverage.
//
// Java: class CoverageInstrumentationCallback implements NodeTraversal.Callback
type CoverageInstrumentationCallback struct {
	instrumentationData *jutil.LinkedHashMap[string, *FileInstrumentationData]
	reach               CoverageInstrumentationPass_CoverageReach
}

// Java: CoverageInstrumentationCallback.ARRAY_NAME_PREFIX
const CoverageInstrumentationCallback_ARRAY_NAME_PREFIX = "JSCompiler_lcov_data_"

// Java: CoverageInstrumentationCallback(Map<String, FileInstrumentationData>, CoverageInstrumentationPass.CoverageReach)
func NewCoverageInstrumentationCallback(
	instrumentationData *jutil.LinkedHashMap[string, *FileInstrumentationData],
	reach CoverageInstrumentationPass_CoverageReach) *CoverageInstrumentationCallback {
	c := &CoverageInstrumentationCallback{}
	c.instrumentationData = instrumentationData
	c.reach = reach
	return c
}

// CoverageInstrumentationCallback_GetFileName returns the name of the source
// file from which the given node originates.
//
// Go-only: language-adapter — NodeTraversal.getSourceName() may return null,
// which a Go string key cannot hold; a script always has a source name here
// (a null one would make CoverageUtil.createIdentifierFromText throw in Java),
// so the nil case panics with a NullPointerException.
//
// Java: CoverageInstrumentationCallback.getFileName(NodeTraversal)
func CoverageInstrumentationCallback_GetFileName(traversal *NodeTraversal) string {
	name := traversal.GetSourceName()
	if name == nil {
		panic(jutil.NewNullPointerException())
	}
	return *name
}

// CreateArrayName returns a string that can be used as array name. The name
// is based on the source filename of the AST node.
//
// Java: CoverageInstrumentationCallback.createArrayName(NodeTraversal)
func (c *CoverageInstrumentationCallback) CreateArrayName(traversal *NodeTraversal) string {
	return CoverageInstrumentationCallback_ARRAY_NAME_PREFIX +
		CoverageUtil_CreateIdentifierFromText(CoverageInstrumentationCallback_GetFileName(traversal))
}

// NewInstrumentationNode creates and return a new instrumentation node. The
// instrumentation Node is of the form: "arrayName[lineNumber] = true;"
// Note 1: Node returns a 1-based line number.
// Note 2: Line numbers in instrumentation are made 0-based. This is done to
// map to their bit representation in BitField. Thus, there's a one-to-one
// correspondence of the line number seen on instrumented files and their bit
// encodings.
//
// Java: CoverageInstrumentationCallback.newInstrumentationNode(NodeTraversal, Node)
func (c *CoverageInstrumentationCallback) NewInstrumentationNode(traversal *NodeTraversal, node rhino.Node) rhino.Node {
	lineNumber := node.GetLineno()
	arrayName := c.CreateArrayName(traversal)

	// Create instrumentation Node
	getElemNode := rhino.IR_Getelem(
		rhino.IR_Name(arrayName),
		rhino.IR_Number(float64(lineNumber-1))) // Make line number 0-based
	exprNode := rhino.IR_ExprResult(rhino.IR_Assign(getElemNode, rhino.IR_TrueNode()))

	// Note line as instrumented
	fileName := CoverageInstrumentationCallback_GetFileName(traversal)
	c.instrumentationData.ComputeIfAbsent(
		fileName, func(k string) *FileInstrumentationData { return NewFileInstrumentationData(k, arrayName) })
	c.instrumentationData.Get(fileName).SetLineAsInstrumented(lineNumber)

	return exprNode.SrcrefTreeIfMissing(node)
}

// NewArrayDeclarationNode creates and returns a new array declaration node.
// The array name is generated based on the source filename, and declaration
// is of the form: "var arrayNameUsedInFile = [];"
//
// Java: CoverageInstrumentationCallback.newArrayDeclarationNode(NodeTraversal)
func (c *CoverageInstrumentationCallback) NewArrayDeclarationNode(traversal *NodeTraversal) rhino.Node {
	return rhino.IR_VarWithValue(
		rhino.IR_Name(c.CreateArrayName(traversal)),
		rhino.IR_Arraylit())
}

// NewHeaderNode returns a Node containing file specific setup logic.
//
// Java: CoverageInstrumentationCallback.newHeaderNode(NodeTraversal, Node)
func (c *CoverageInstrumentationCallback) NewHeaderNode(traversal *NodeTraversal, srcref rhino.Node) rhino.Node {
	fileName := CoverageInstrumentationCallback_GetFileName(traversal)
	arrayName := c.CreateArrayName(traversal)
	data := c.instrumentationData.Get(fileName)
	jutil.CheckNotNull(data)

	objName := CoverageInstrumentationPass_JS_INSTRUMENTATION_OBJECT_NAME

	// var JSCompiler_lcov_data_xx = [];
	// __jscov['executedLines'].push(JSCompiler_lcov_data_xx);
	// __jscov['instrumentedLines'].push(hex-data);
	// __jscov['fileNames'].push(filename);
	return rhino.IR_BlockOfNodes(
		c.NewArrayDeclarationNode(traversal),
		rhino.IR_ExprResult(
			rhino.IR_Call(
				rhino.IR_Getprop(rhino.IR_Getelem(rhino.IR_Name(objName), rhino.IR_String("executedLines")), "push"),
				rhino.IR_Name(arrayName))),
		rhino.IR_ExprResult(
			rhino.IR_Call(
				rhino.IR_Getprop(
					rhino.IR_Getelem(rhino.IR_Name(objName), rhino.IR_String("instrumentedLines")), "push"),
				rhino.IR_String(data.GetInstrumentedLinesAsHexString()))),
		rhino.IR_ExprResult(
			rhino.IR_Call(
				rhino.IR_Getprop(rhino.IR_Getelem(rhino.IR_Name(objName), rhino.IR_String("fileNames")), "push"),
				rhino.IR_String(fileName)))).
		SrcrefTreeIfMissing(srcref)
}

// ShouldTraverse is final in Java.
//
// Java: CoverageInstrumentationCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *CoverageInstrumentationCallback) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Skip scripts that are marked with @nocoverage.
	return !n.IsScript() || n.GetJSDocInfo() == nil || !n.GetJSDocInfo().IsNoCoverage()
}

// Visit instruments the JS code by inserting appropriate nodes into the AST.
// The instrumentation logic is tuned to collect "line coverage" data only.
//
// Java: CoverageInstrumentationCallback.visit(NodeTraversal, Node, Node)
func (c *CoverageInstrumentationCallback) Visit(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) {
	// SCRIPT node is special - it is the root of the AST for code from a file.
	// Append code to declare and initialize structures used in instrumentation.
	if node.IsScript() {
		// Go-only: language-adapter — the null source name (see GetFileName)
		// is only looked up here, where Java's map.get(null) is null.
		if name := traversal.GetSourceName(); name != nil && c.instrumentationData.Get(*name) != nil {
			if node.HasChildren() && node.GetFirstChild().IsModuleBody() {
				node = node.GetFirstChild()
			}
			node.AddChildrenToFront(c.NewHeaderNode(traversal, node).RemoveChildren())
		}
		traversal.ReportCodeChange()
		return
	}

	// Don't instrument global statements
	if c.reach == CoverageInstrumentationPass_CoverageReach_CONDITIONAL &&
		!parent.IsNil() && parent.IsScript() {
		return
	}

	// For arrow functions whose body is an expression instead of a block,
	// convert it to a block so that it can be instrumented.
	if node.IsFunction() && !NodeUtil_GetFunctionBody(node).IsBlock() {
		returnValue := NodeUtil_GetFunctionBody(node)
		body := rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(returnValue.Detach()))
		body.SrcrefTreeIfMissing(returnValue)
		node.AddChildToBack(body)
	}

	// Add instrumentation code just before a function block.
	// Similarly before other constructs: 'with', 'case', 'default', 'catch'
	if node.IsFunction() ||
		node.IsWith() ||
		node.IsCase() ||
		node.IsDefaultCase() ||
		node.IsCatch() {
		codeBlock := node.GetLastChild()
		codeBlock.AddChildToFront(
			c.NewInstrumentationNode(traversal, node))
		traversal.ReportCodeChange()
		return
	}

	// Add instrumentation code as the first child of a 'try' block.
	if node.IsTry() {
		firstChild := node.GetFirstChild()
		firstChild.AddChildToFront(
			c.NewInstrumentationNode(traversal, node))
		traversal.ReportCodeChange()
		return
	}

	// For any other statement, add instrumentation code just before it.
	if !parent.IsNil() && NodeUtil_IsStatementBlock(parent) && !node.IsModuleBody() {
		c.NewInstrumentationNode(traversal, node).InsertBefore(node)
		traversal.ReportCodeChange()
		return
	}
}
