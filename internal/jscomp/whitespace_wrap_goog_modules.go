// Port of com/google/javascript/jscomp/WhitespaceWrapGoogModules.java @ closure-compiler v20260415

/*
 * Copyright 2015 The Closure Compiler Authors.
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

// WhitespaceWrapGoogModules replicates the effect of ClosureBundler in
// whitespace-only mode and wraps goog.modules in goog.loadModule calls. See
// comment block below.
type WhitespaceWrapGoogModules struct {
	compiler AbstractCompiler
}

// Java: WhitespaceWrapGoogModules(AbstractCompiler)
func NewWhitespaceWrapGoogModules(compiler AbstractCompiler) *WhitespaceWrapGoogModules {
	w := &WhitespaceWrapGoogModules{}
	w.compiler = compiler
	return w
}

// Java: WhitespaceWrapGoogModules.process(Node, Node)
func (w *WhitespaceWrapGoogModules) Process(externs *rhino.Node, root *rhino.Node) {
	for c := root.GetFirstChild(); c != nil; c = c.GetNext() {
		jutil.CheckState(c.IsScript())
		w.VisitScript(c)
	}
}

// Java: WhitespaceWrapGoogModules.visitScript(Node)
func (w *WhitespaceWrapGoogModules) VisitScript(scriptRoot *rhino.Node) {
	if !NodeUtil_IsGoogModuleFile(scriptRoot) {
		return
	}
	moduleBody := scriptRoot.GetFirstChild()
	moduleBody.SetToken(rhino.Token_BLOCK)
	NodeUtil_TryMergeBlock(moduleBody, true)
	w.compiler.ReportChangeToEnclosingScope(scriptRoot)

	// As per ClosureBundler:
	/*
	   // add the prefix on the first line so the line numbers aren't affected.
	   out.append(
	       "goog.loadModule(function(exports) {"
	       + "'use strict';");
	   append(out, Mode.NORMAL, contents);
	   out.append(
	       "\n" // terminate any trailing single line comment.
	       + ";" // terminate any trailing expression.
	       + "return exports;});\n");
	   appendSourceUrl(out, Mode.NORMAL);
	*/

	block := rhino.IR_Block()
	block.AddChildToBack(
		rhino.IR_ExprResult(rhino.IR_String("use strict"))) // needs to be explicit, to match ClosureBundler

	function := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(rhino.IR_Name("exports")), block)
	w.compiler.ReportChangeToChangeScope(function)
	loadMod :=
		rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name("goog"), "loadModule"), function)).
			SrcrefTree(scriptRoot)

	if scriptRoot.HasChildren() {
		block.AddChildrenToBack(scriptRoot.RemoveChildren())
	}
	block.AddChildToBack(rhino.IR_ReturnNodeWithExpr(rhino.IR_Name("exports")).SrcrefTree(scriptRoot))

	scriptRoot.AddChildToBack(loadMod)
	w.compiler.ReportChangeToChangeScope(scriptRoot)
}
