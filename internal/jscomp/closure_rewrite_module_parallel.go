// Copyright 2026 The Closure Compiler Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

func (m *ClosureRewriteModule) processPhase2Parallel(
	externs rhino.Node,
	root rhino.Node,
	scriptDescriptions []*closureScriptDescription,
) bool {
	if m.preprocessorSymbolTable != nil ||
		m.globalTypedScope != nil ||
		m.compiler == nil ||
		m.compiler.HasTypeCheckingRun() ||
		m.compiler.IsDebugLoggingEnabled() {
		return false
	}
	scopeRoot, items, ok := collectExternAndRootScripts(externs, root)
	if !ok || len(items) != len(scriptDescriptions) {
		return false
	}
	if len(items) == 0 {
		return true
	}

	hasNonLegacyModules := len(m.rewriteState.nonLegacyModuleRoots) > 0
	isActive := func(idx int) bool {
		s := items[idx].script
		if s.IsFromExterns() && !NodeUtil_IsFromTypeSummary(s) {
			return false
		}
		return scriptDescriptions[idx].hasGoogOrModuleBody || hasNonLegacyModules
	}

	activeCount := 0
	hasNonExtern := false
	for i := range items {
		s := items[i].script
		if !s.IsFromExterns() || NodeUtil_IsFromTypeSummary(s) {
			hasNonExtern = true
		}
		if isActive(i) {
			activeCount++
		}
	}
	if activeCount == 0 {
		if len(m.syntheticExternNodes) > 0 && hasNonExtern && m.globalScope == nil {
			sc := NewSyntacticScopeCreator(m.compiler)
			m.globalScope = sc.CreateScopeOfScope(scopeRoot, nil)
		}
		return true
	}
	numWorkers := parallelThreadCount(m.compiler, activeCount)
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(m.compiler)
	// Capture the global scope before any goog.module MODULE_BODY is rewritten to BLOCK.
	lazyGlobal := newLazyGlobalScope(m.compiler, scopeRoot)
	m.globalScope = lazyGlobal.get()

	runParallelWork(m.compiler, numWorkers, len(items), func(_ int) (func(int, *recordingCompiler), func()) {
		var wm ClosureRewriteModule
		var wsc *SyntacticScopeCreator
		initialized := false

		runItem := func(idx int, rc *recordingCompiler) {
			if !isActive(idx) {
				return
			}
			if !initialized {
				wm = *m
				wm.scriptStack = nil
				wm.currentScript = nil
				wsc = NewSyntacticScopeCreator(rc)
				initialized = true
			}
			wm.compiler = rc
			wm.currentScript = nil
			wm.scriptStack = wm.scriptStack[:0]
			updater := closureScriptUpdater{
				outer:             &wm,
				singleDescription: scriptDescriptions[idx],
			}
			t := acquireNodeTraversalDirect(rc, &updater, wsc, false)
			traverseScriptUnderRoot(t, scopeRoot, items[idx].parent, items[idx].script, lazyGlobal)
			releaseNodeTraversal(t)
		}
		return runItem, nil
	}, false)
	return true
}
