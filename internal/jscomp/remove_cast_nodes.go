// Port of com/google/javascript/jscomp/RemoveCastNodes.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
			t := acquireNodeTraversal(&b)
			traverseScriptUnderRoot(t, scopeRoot, item.parent, item.script, nil)
			releaseNodeTraversal(t)
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
