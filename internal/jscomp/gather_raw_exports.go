// Port of com/google/javascript/jscomp/GatherRawExports.java @ closure-compiler v20261015

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
