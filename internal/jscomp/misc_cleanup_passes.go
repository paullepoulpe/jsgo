// Port of com/google/javascript/jscomp/CheckDebuggerStatement.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/GatherRawExports.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/RemoveCastNodes.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/RemoveUnnecessarySyntheticExterns.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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

/*
 * Copyright 2021 The Closure Compiler Authors.
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

// CheckDebuggerStatement checks for the presence of the "debugger" statement
// in JavaScript code. It is appropriate to use this statement while developing
// JavaScript; however, it is generally undesirable to include it in production
// code.
type CheckDebuggerStatement struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler
}

// Java: CheckDebuggerStatement.DEBUGGER_STATEMENT_PRESENT
var CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT = DiagnosticType_Disabled(
	"JSC_DEBUGGER_STATEMENT_PRESENT",
	"Using the debugger statement can halt your application if the user has a JavaScript"+
		" debugger running.\n"+
		"To disable this check when you want to do debugging, you can suppress this message"+
		" like this:\n"+
		// NOTE: Single quotes around curly braces are necessary to prevent them being
		// interpreted as a placeholder for replacement. See MessageFormat javadoc.
		"  /** @suppress '{'checkDebuggerStatement'}' */\n"+
		"  debugger;")

// Java: CheckDebuggerStatement(AbstractCompiler)
func NewCheckDebuggerStatement(compiler AbstractCompiler) *CheckDebuggerStatement {
	c := &CheckDebuggerStatement{}
	c.compiler = compiler
	return c
}

// Java: CheckDebuggerStatement.process(Node, Node)
func (c *CheckDebuggerStatement) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckDebuggerStatement.visit(NodeTraversal, Node, Node)
func (c *CheckDebuggerStatement) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsDebugger() {
		t.Report(n, CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT)
	}
}

// GatherRawExports collects known properties on the global windows object to
// avoid creating global variable names that conflict with these names. This
// pass runs after property renaming but before variable renaming.
//
// This is a best effort pass and does not guarantee that there are no
// conflicts.
//
// This is not required if the global variables are isolated from global
// scope (isolation_mode=IIFE) or equivalent.
type GatherRawExports struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler

	exportedVariables *jutil.LinkedHashSet[string]
}

// TODO(johnlenz): "goog$global" should be part of a coding convention.
// Note: GatherRawExports runs after property renaming and
// collapse properties, so the two entries here protect goog.global in the
// two common cases "collapse properties and renaming on" or both off
// but not the case where only property renaming is on.
var gatherRawExports_GLOBAL_THIS_NAMES = []string{
	"window",
	"globalThis",
	"top",
	"self",
	"goog$global",
	"goog.global",
	"$jscomp.global",
	"$jscomp$global",
}

// Java: GatherRawExports(AbstractCompiler)
func NewGatherRawExports(compiler AbstractCompiler) *GatherRawExports {
	g := &GatherRawExports{}
	g.exportedVariables = jutil.NewLinkedHashSet[string]()
	g.compiler = compiler
	return g
}

// Java: GatherRawExports.process(Node, Node)
func (g *GatherRawExports) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(g.compiler.GetLifeCycleStage().IsNormalized())
	NodeTraversal_Traverse(g.compiler, root, g)
}

// Java: GatherRawExports.visit(NodeTraversal, Node, Node)
func (g *GatherRawExports) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if NodeUtil_IsNormalOrOptChainGet(n) && GatherRawExports_IsGlobalThisObject(t, n.GetFirstChild()) {
		if NodeUtil_IsNormalOrOptChainGetProp(n) {
			g.exportedVariables.Add(n.GetString())
		} else if NodeUtil_IsNormalOrOptChainGet(n) && n.GetSecondChild().IsStringLit() {
			g.exportedVariables.Add(n.GetSecondChild().GetString())
		}
	}
}

// Java: GatherRawExports.isGlobalThisObject(NodeTraversal, Node)
func GatherRawExports_IsGlobalThisObject(t *NodeTraversal, n rhino.Node) bool {
	// We should consider using scopes or types to check for references to the "global this".
	if n.IsThis() {
		return t.InGlobalHoistScope()
	} else if n.IsQualifiedName() {
		items := len(gatherRawExports_GLOBAL_THIS_NAMES)
		for i := 0; i < items; i++ {
			if n.MatchesQualifiedName(gatherRawExports_GLOBAL_THIS_NAMES[i]) {
				return true
			}
		}
	}
	return false
}

// Java: GatherRawExports.getExportedVariableNames()
func (g *GatherRawExports) GetExportedVariableNames() *jutil.LinkedHashSet[string] {
	return g.exportedVariables
}

// RemoveCastNodes normalizes all CAST nodes out of the AST.
type RemoveCastNodes struct {
	compiler AbstractCompiler
}

// Java: RemoveCastNodes(AbstractCompiler)
func NewRemoveCastNodes(compiler AbstractCompiler) *RemoveCastNodes {
	return &RemoveCastNodes{compiler: compiler}
}

// Java: RemoveCastNodes.process(Node, Node)
func (r *RemoveCastNodes) Process(externs rhino.Node, root rhino.Node) {
	if r.processParallel(externs, root) {
		return
	}
	NodeTraversal_TraverseRoots(r.compiler, &RemoveCastNodes_RemoveCastsCallback{outer: r}, externs, root)
}

func (r *RemoveCastNodes) processParallel(externs rhino.Node, root rhino.Node) bool {
	scopeRoot, items, ok := collectExternAndRootScripts(externs, root)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(r.compiler, len(items))
	if numWorkers <= 1 {
		return false
	}
	prewarmCompilerForParallelReads(r.compiler)
	runParallelWork(r.compiler, numWorkers, len(items), func(_ int) (func(int, *recordingCompiler), func()) {
		var workerRC *recordingCompiler
		var workerCb *RemoveCastNodes_RemoveCastsCallback
		runItem := func(idx int, rc *recordingCompiler) {
			if workerCb == nil || workerRC != rc {
				workerRC = rc
				workerCb = &RemoveCastNodes_RemoveCastsCallback{
					outer: &RemoveCastNodes{compiler: rc},
				}
			}
			item := items[idx]
			b := NodeTraversal_Builder{
				compiler: rc,
				callback: workerCb,
			}
			b.traverseScriptUnderRoot(scopeRoot, item.parent, item.script, nil)
		}
		return runItem, nil
	}, false)
	return true
}

// RemoveCastNodes_RemoveCastsCallback is the inner class RemoveCastsCallback.
type RemoveCastNodes_RemoveCastsCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *RemoveCastNodes
}

// Java: RemoveCastNodes.RemoveCastsCallback.visit(NodeTraversal, Node, Node)
func (c *RemoveCastNodes_RemoveCastsCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsCast() {
		c.outer.compiler.ReportChangeToEnclosingScope(n)
		n.ReplaceWith(n.RemoveFirstChild())
	}
}

// RemoveUnnecessarySyntheticExterns deletes name declarations from synthetic
// externs if they are in fact declared in source
//
// The VarCheck pass will create synthetic extern declarations (`var
// someName;`) for every variable name referenced without a corresponding
// externs or code declaration. Such names are "unfulfilled". VarCheck only
// does this because later compiler passes expect the following invariant to
// hold: all variable names referenced in externs or code are also declared
// in externs or code.
//
// When running optimizations based on precompiled TypedAST library shards,
// it's possible for a name to be unfulfilled in shard A, but declared in
// another shard B. Shard B "fulfills" the reference from shard A. This pass
// removes the corresponding synthetic extern that came from shard A now that
// the name is fulfilled. The purpose is to improve optimizations, e.g. to
// allow mangling the name if it is not defined in externs.
type RemoveUnnecessarySyntheticExterns struct {
	compiler      AbstractCompiler
	nodesToDetach *jutil.LinkedHashSet[rhino.Node]
}

// Java: RemoveUnnecessarySyntheticExterns(AbstractCompiler)
func NewRemoveUnnecessarySyntheticExterns(compiler AbstractCompiler) *RemoveUnnecessarySyntheticExterns {
	return &RemoveUnnecessarySyntheticExterns{
		compiler:      compiler,
		nodesToDetach: jutil.NewLinkedHashSet[rhino.Node](),
	}
}

// Java: RemoveUnnecessarySyntheticExterns.process(Node, Node)
func (r *RemoveUnnecessarySyntheticExterns) Process(externs rhino.Node, root rhino.Node) {
	globalRoot := externs.GetParent()
	scopeCreator := NewSyntacticScopeCreatorWithTreatProvidesAsRedeclarations(
		r.compiler, &RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler{outer: r}, true)

	// This call will invoke the RedeclarationCheckHandler on duplicate synthetic names
	scopeCreator.CreateScopeOfScope(globalRoot, nil /* parent= */)

	if r.nodesToDetach.IsEmpty() {
		return
	}

	for n := range r.nodesToDetach.All() {
		n.Detach()
	}

	syntheticExternsScript := r.compiler.GetSynthesizedExternsInput().GetAstRoot(r.compiler)
	r.compiler.ReportChangeToEnclosingScope(syntheticExternsScript)
}

// RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler is the handler
// for duplicate declarations.
type RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler struct {
	outer *RemoveUnnecessarySyntheticExterns
}

// Java: RemoveUnnecessarySyntheticExterns.RedeclarationCheckHandler.onRedeclaration(Scope, String, Node, CompilerInput)
func (h *RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) OnRedeclaration(s *Scope, name string, n rhino.Node, input *CompilerInput) {
	jutil.CheckArgumentMsg2(s.IsGlobal(), "Unexpected non-global redeclaration %s in scope %s", n, s)

	isSynthetic := input == h.outer.compiler.GetSynthesizedExternsInput()
	origVar := s.GetVar(name)
	// origNode will be null for `arguments`, since there's no node that declares it.
	origNode := origVar.GetNode()

	if isSynthetic {
		// Delete duplicate unfulfilled synthetic declarations.
		// The original declaration will be deleted later iff it's also synthetic and we visit a
		// non-synthetic duplicate declaration.
		if n.GetParent().IsSynthesizedUnfulfilledNameDeclaration() {
			h.outer.nodesToDetach.Add(n.GetParent())
		}
	} else if !origNode.IsNil() &&
		origNode.GetParent().IsSynthesizedUnfulfilledNameDeclaration() {
		// Found a synthetic extern declaration that's fulfilled by some non-synthetic code.
		h.outer.nodesToDetach.Add(origNode.GetParent())
	}
}
