// Port of com/google/javascript/jscomp/J2clUtilGetDefineRewriterPass.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// J2clUtilGetDefineRewriterPass is a normalization pass to re-write
// Util.$getDefine calls to make them work in compiled mode.
//
// Java: public class J2clUtilGetDefineRewriterPass extends AbstractPostOrderCallback implements CompilerPass
type J2clUtilGetDefineRewriterPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
	defines  *jutil.LinkedHashSet[string]
}

// Java: J2clUtilGetDefineRewriterPass(AbstractCompiler)
func NewJ2clUtilGetDefineRewriterPass(compiler AbstractCompiler) *J2clUtilGetDefineRewriterPass {
	return &J2clUtilGetDefineRewriterPass{compiler: compiler}
}

// Java: J2clUtilGetDefineRewriterPass.process(Node, Node)
func (p *J2clUtilGetDefineRewriterPass) Process(externs *rhino.Node, root *rhino.Node) {
	if !J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
		return
	}
	p.defines = p.compiler.GetDefineNames()

	checkGetDefineCalls := p.newCheckGetDefineCalls()
	NodeTraversal_Traverse(p.compiler, root, checkGetDefineCalls)
	checkGetDefineCalls.validateRegisteredDefines()

	NodeTraversal_Traverse(p.compiler, root, p)
}

// Java: J2clUtilGetDefineRewriterPass.visit(NodeTraversal, Node, Node)
func (p *J2clUtilGetDefineRewriterPass) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if j2clUtilGetDefineRewriterPass_isUtilGetDefineCall(n) {
		p.rewriteUtilGetDefine(t, n)
	}
}

// Java: J2clUtilGetDefineRewriterPass.rewriteUtilGetDefine(NodeTraversal, Node)
func (p *J2clUtilGetDefineRewriterPass) rewriteUtilGetDefine(t *NodeTraversal, callNode *rhino.Node) {
	firstExpr := callNode.GetSecondChild()
	secondExpr := callNode.GetLastChild()

	if secondExpr != firstExpr {
		secondExpr.Detach()
	} else {
		// There is no secondExpr; default to null.
		secondExpr = rhino.IR_NullNode()
	}

	replacement := p.getDefineReplacement(firstExpr, secondExpr)
	replacement.SrcrefTreeIfMissing(callNode)
	callNode.ReplaceWith(replacement)
	t.ReportCodeChange()
}

// Java: J2clUtilGetDefineRewriterPass.getDefineReplacement(Node, Node)
func (p *J2clUtilGetDefineRewriterPass) getDefineReplacement(firstExpr *rhino.Node, secondExpr *rhino.Node) *rhino.Node {
	defineName := firstExpr.GetString()
	if p.defines.Contains(defineName) {
		define := NodeUtil_NewQName(p.compiler, ProcessDefines_GetGlobalDefineAlias(defineName))
		defineStringValue := NodeUtil_NewCallNode(rhino.IR_Name("String"), define)
		return rhino.IR_Comma(secondExpr, defineStringValue)
	} else {
		return secondExpr
	}
}

// Java: J2clUtilGetDefineRewriterPass.isUtilGetDefineCall(Node)
func j2clUtilGetDefineRewriterPass_isUtilGetDefineCall(n *rhino.Node) bool {
	return n.IsCall() && j2clUtilGetDefineRewriterPass_isUtilGetDefineMethodName(n.GetFirstChild().GetQualifiedName())
}

// Java: J2clUtilGetDefineRewriterPass.isUtilGetDefineMethodName(String)
func j2clUtilGetDefineRewriterPass_isUtilGetDefineMethodName(fnName *string) bool {
	// TODO: Switch this to the filename + property name heuristic which is less brittle.
	return fnName != nil && strings.HasSuffix(*fnName, ".$getDefine") && strings.Contains(*fnName, "Util")
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
type J2clUtilGetDefineRewriterPass_CheckGetDefineCalls struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer                             *J2clUtilGetDefineRewriterPass
	defineNameToAddSystemPropertyCall *jutil.LinkedHashMap[string, *rhino.Node]
	defineNameToGetDefineCall         *jutil.LinkedHashMap[string, *rhino.Node]
}

// Java: J2clUtilGetDefineRewriterPass.CheckGetDefineCalls()
func (p *J2clUtilGetDefineRewriterPass) newCheckGetDefineCalls() *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls {
	return &J2clUtilGetDefineRewriterPass_CheckGetDefineCalls{
		outer:                             p,
		defineNameToAddSystemPropertyCall: jutil.NewLinkedHashMap[string, *rhino.Node](),
		defineNameToGetDefineCall:         jutil.NewLinkedHashMap[string, *rhino.Node](),
	}
}

// Java: J2clUtilGetDefineRewriterPass.CheckGetDefineCalls.visit(NodeTraversal, Node, Node)
func (c *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	compiler := c.outer.compiler
	if j2clUtilGetDefineRewriterPass_CheckGetDefineCalls_isAddSystemPropertyFromGoogDefine(n) {
		defineName := j2clUtilGetDefineRewriterPass_CheckGetDefineCalls_getDefineNameParameter(n)
		if defineName != nil {
			c.defineNameToAddSystemPropertyCall.Put(*defineName, n)
		} else {
			compiler.Report(JSError_MakeWithNode(n, J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME))
		}
	} else if j2clUtilGetDefineRewriterPass_isUtilGetDefineCall(n) {
		defineName := j2clUtilGetDefineRewriterPass_CheckGetDefineCalls_getDefineNameParameter(n)
		if defineName != nil {
			c.defineNameToGetDefineCall.Put(*defineName, n)
		} else {
			compiler.Report(JSError_MakeWithNode(n, J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME))
		}
	}
}

// Java: J2clUtilGetDefineRewriterPass.CheckGetDefineCalls.validateRegisteredDefines()
func (c *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) validateRegisteredDefines() {
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
func j2clUtilGetDefineRewriterPass_CheckGetDefineCalls_getDefineNameParameter(callNode *rhino.Node) *string {
	nameNode := callNode.GetSecondChild()
	if nameNode == nil || !nameNode.IsString() {
		return nil
	}
	name := nameNode.GetString()
	return &name
}

// Java: J2clUtilGetDefineRewriterPass.CheckGetDefineCalls.isAddSystemPropertyFromGoogDefine(Node)
func j2clUtilGetDefineRewriterPass_CheckGetDefineCalls_isAddSystemPropertyFromGoogDefine(n *rhino.Node) bool {
	return n.IsCall() &&
		n.GetFirstChild().GetQualifiedName() != nil &&
		strings.HasSuffix(*n.GetFirstChild().GetQualifiedName(), "addSystemPropertyFromGoogDefine")
}
