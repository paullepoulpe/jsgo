// Port of com/google/javascript/jscomp/J2clUtilGetDefineRewriterPass.java @ closure-compiler v20261028

/*
 * Copyright 2018 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var j2clUtilGetDefineVisitTokens = makeVisitTokenSet(rhino.Token_CALL)

// J2clUtilGetDefineRewriterPass is a normalization pass to re-write
// Util.$getDefine calls to make them work in compiled mode.
//
// Java: public class J2clUtilGetDefineRewriterPass extends AbstractPostOrderCallback implements CompilerPass
type J2clUtilGetDefineRewriterPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
	defines  *OrderedSet[string]
}

func (*J2clUtilGetDefineRewriterPass) noVisitLeaves() {}
func (*J2clUtilGetDefineRewriterPass) visitTokens() *ast.TokenSet {
	return &j2clUtilGetDefineVisitTokens
}

// Java: J2clUtilGetDefineRewriterPass(AbstractCompiler)
func NewJ2clUtilGetDefineRewriterPass(compiler AbstractCompiler) *J2clUtilGetDefineRewriterPass {
	return &J2clUtilGetDefineRewriterPass{compiler: compiler}
}

// Java: J2clUtilGetDefineRewriterPass.process(Node, Node)
func (p *J2clUtilGetDefineRewriterPass) Process(externs rhino.Node, root rhino.Node) {
	if !J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
		return
	}
	p.defines = p.compiler.GetDefineNames()

	checker := &j2clCheckGetDefineCalls{outer: p}
	NodeTraversal_Traverse(p.compiler, root, checker)
	checker.validateRegisteredDefines()

	NodeTraversal_Traverse(p.compiler, root, p)
}

// Java: J2clUtilGetDefineRewriterPass.visit(NodeTraversal, Node, Node)
func (p *J2clUtilGetDefineRewriterPass) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if call, ok := ast.AsCall(n); ok && isUtilGetDefineCall(call) {
		p.rewriteUtilGetDefine(t, call)
	}
}

// Java: J2clUtilGetDefineRewriterPass.rewriteUtilGetDefine(NodeTraversal, Node)
func (p *J2clUtilGetDefineRewriterPass) rewriteUtilGetDefine(t *NodeTraversal, call ast.Call) {
	firstArg, ok := call.FirstArg()
	if !ok {
		return
	}
	firstExpr := firstArg.Node()
	lastArg, _ := call.LastArg()
	secondExpr := lastArg.Node()

	if secondExpr != firstExpr {
		secondExpr.Detach()
	} else {
		// There is no secondExpr; default to null.
		secondExpr = rhino.IR_NullNode()
	}

	callNode := call.Node()
	replacement := p.getDefineReplacement(firstExpr, secondExpr)
	replacement.SrcrefTreeIfMissing(callNode)
	callNode.ReplaceWith(replacement)
	t.ReportCodeChange()
}

// Java: J2clUtilGetDefineRewriterPass.getDefineReplacement(Node, Node)
func (p *J2clUtilGetDefineRewriterPass) getDefineReplacement(firstExpr rhino.Node, secondExpr rhino.Node) rhino.Node {
	if firstExpr.IsNil() || !firstExpr.IsStringLit() {
		return secondExpr
	}
	defineName := firstExpr.GetString()
	if p.defines.Contains(defineName) {
		define := NodeUtil_NewQName(p.compiler, ProcessDefines_GetGlobalDefineAlias(defineName))
		defineStringValue := NodeUtil_NewCallNode(rhino.IR_Name("String"), define)
		return rhino.IR_Comma(secondExpr, defineStringValue)
	}
	return secondExpr
}

// Java: J2clUtilGetDefineRewriterPass.isUtilGetDefineCall(Node)
func isUtilGetDefineCall(call ast.Call) bool {
	fnName, ok := call.Callee().Node().GetQualifiedNameString()
	return ok && isUtilGetDefineMethodName(fnName)
}

// Java: J2clUtilGetDefineRewriterPass.isUtilGetDefineMethodName(String)
func isUtilGetDefineMethodName(fnName string) bool {
	// TODO: Switch this to the filename + property name heuristic which is less brittle.
	return strings.HasSuffix(fnName, ".$getDefine") && strings.Contains(fnName, "Util")
}

var (
	J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME = DiagnosticType_Error(
		"JSC_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME",
		"Calls to System.getProperty must use a String literal for the property name")

	J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY = DiagnosticType_Error(
		"JSC_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY",
		"Unknown system property \"{0}\", are you missing a call to"+
			" jre.addSystemPropertyFromGoogDefine?")

	J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME = DiagnosticType_Error(
		"JSC_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME",
		"Calls to addSystemPropertyFromGoogDefine must use a String literal for the define"+
			" name")

	J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE = DiagnosticType_Error(
		"JSC_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE", "Unknown goog.define with name: {0}")
)

// Java: private final class J2clUtilGetDefineRewriterPass.CheckGetDefineCalls extends AbstractPostOrderCallback
type j2clCheckGetDefineCalls struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer                             *J2clUtilGetDefineRewriterPass
	defineNameToAddSystemPropertyCall OrderedStringMap[rhino.Node]
	defineNameToGetDefineCall         OrderedStringMap[rhino.Node]
}

func (*j2clCheckGetDefineCalls) noVisitLeaves() {}
func (*j2clCheckGetDefineCalls) visitTokens() *ast.TokenSet {
	return &j2clUtilGetDefineVisitTokens
}

// Java: J2clUtilGetDefineRewriterPass.CheckGetDefineCalls.visit(NodeTraversal, Node, Node)
func (c *j2clCheckGetDefineCalls) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	call, ok := ast.AsCall(n)
	if !ok {
		return
	}
	compiler := c.outer.compiler
	if isAddSystemPropertyFromGoogDefine(call) {
		if defineName, ok := defineNameParameter(call); ok {
			c.defineNameToAddSystemPropertyCall.Put(defineName, n)
		} else {
			compiler.Report(JSError_MakeWithNode(n, J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME))
		}
	} else if isUtilGetDefineCall(call) {
		if defineName, ok := defineNameParameter(call); ok {
			c.defineNameToGetDefineCall.Put(defineName, n)
		} else {
			compiler.Report(JSError_MakeWithNode(n, J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME))
		}
	}
}

// Java: J2clUtilGetDefineRewriterPass.CheckGetDefineCalls.validateRegisteredDefines()
func (c *j2clCheckGetDefineCalls) validateRegisteredDefines() {
	compiler := c.outer.compiler
	defines := c.outer.defines
	// Validate that all addSystemPropertyFromGoogDefine calls reference a known define.
	for key, value := range c.defineNameToAddSystemPropertyCall.All() {
		if !defines.Contains(key) {
			compiler.Report(
				JSError_MakeWithNode(
					value, J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE, key))
		}
	}

	// Validate that Util.$getDefine() calls reference a define registered with
	// addSystemPropertyFromGoogDefine.
	for key, value := range c.defineNameToGetDefineCall.All() {
		// Tolerate System.getProperty() calls that reference something that is never goog.define'd.
		// This provides compatibility with code that never intended reference a define (ex. third
		// party code).
		if !defines.Contains(key) {
			continue
		}

		if !c.defineNameToAddSystemPropertyCall.ContainsKey(key) {
			compiler.Report(
				JSError_MakeWithNode(
					value, J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY, key))
		}
	}
}

// Java: J2clUtilGetDefineRewriterPass.CheckGetDefineCalls.getDefineNameParameter(Node)
func defineNameParameter(call ast.Call) (string, bool) {
	firstArg, ok := call.FirstArg()
	if !ok {
		return "", false
	}
	str, ok := ast.AsStringLit(firstArg.Node())
	if !ok {
		return "", false
	}
	return str.Value(), true
}

// Java: J2clUtilGetDefineRewriterPass.CheckGetDefineCalls.isAddSystemPropertyFromGoogDefine(Node)
func isAddSystemPropertyFromGoogDefine(call ast.Call) bool {
	qname, ok := call.Callee().Node().GetQualifiedNameString()
	return ok && strings.HasSuffix(qname, "addSystemPropertyFromGoogDefine")
}
