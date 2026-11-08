// Port of com/google/javascript/jscomp/RescopeGlobalSymbolsRenameHelper.java @ closure-compiler v20261104

/*
 * Copyright 2026 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	rescopeDisambiguationSuffix     = "$"
	rescopeUnwrappedLocalNamePrefix = "_"
)

// rescopeRenameHelper disambiguates variable names in RescopeGlobalSymbols and
// rescopeRewriteCallback.
//
// Java: RescopeGlobalSymbolsRenameHelper
type rescopeRenameHelper struct {
	compiler                           AbstractCompiler
	globalSymbolNamespace              string
	disambigPrefix                     string
	unwrappedLocalNamesInDefiningChunk map[string]string
}

// rescopeRenameHelperBuilder records identifiers encountered during the
// pre-traversal and collects disambiguated unwrapped local variable names in
// the defining chunk when build is called.
//
// Java: RescopeGlobalSymbolsRenameHelper.HelperBuilder
type rescopeRenameHelperBuilder struct {
	compiler                    AbstractCompiler
	globalSymbolNamespace       string
	disambigPrefix              string
	optimizeLocalAccess         CompilerOptions_OptimizeLocalAccess
	namesStartingWithUnderscore map[string]struct{}
}

func newRescopeRenameHelperBuilder(
	compiler AbstractCompiler,
	globalSymbolNamespace string,
	optimizeLocalAccess CompilerOptions_OptimizeLocalAccess,
) rescopeRenameHelperBuilder {
	return rescopeRenameHelperBuilder{
		compiler:              compiler,
		globalSymbolNamespace: globalSymbolNamespace,
		disambigPrefix:        globalSymbolNamespace + rescopeDisambiguationSuffix,
		optimizeLocalAccess:   optimizeLocalAccess,
	}
}

func (b *rescopeRenameHelperBuilder) recordExistingName(name string) {
	if !b.optimizeLocalAccess.ShouldKeepLocalAliasForWrappedReassignables() {
		return
	}
	if strings.HasPrefix(name, rescopeUnwrappedLocalNamePrefix) {
		if b.namesStartingWithUnderscore == nil {
			b.namesStartingWithUnderscore = make(map[string]struct{})
		}
		b.namesStartingWithUnderscore[name] = struct{}{}
	}
	if rescopeConflictsWithGlobalSymbolNamespace(name, b.globalSymbolNamespace, b.disambigPrefix) {
		disambiguatedName := name + rescopeDisambiguationSuffix
		if strings.HasPrefix(disambiguatedName, rescopeUnwrappedLocalNamePrefix) {
			if b.namesStartingWithUnderscore == nil {
				b.namesStartingWithUnderscore = make(map[string]struct{})
			}
			b.namesStartingWithUnderscore[disambiguatedName] = struct{}{}
		}
	}
}

func (b *rescopeRenameHelperBuilder) build(symbolInfo *rescopeSymbolInfo) rescopeRenameHelper {
	var unwrapped map[string]string
	if b.optimizeLocalAccess.ShouldKeepLocalAliasForWrappedReassignables() && !symbolInfo.wrappedCrossChunkNames.IsEmpty() {
		// Reuse b.namesStartingWithUnderscore directly as usedNames (builder is consumed once).
		usedNames := b.namesStartingWithUnderscore
		if usedNames == nil {
			usedNames = make(map[string]struct{})
		}
		for name := range symbolInfo.wrappedCrossChunkNames.All() {
			_, hasUnsafeWrite := symbolInfo.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk[name]
			_, hasRead := symbolInfo.globalNamesWithReadInDefiningChunk[name]
			if !symbolInfo.crossChunkNamesWithWriteFromOtherChunk.Has(name) &&
				!hasUnsafeWrite &&
				hasRead {
				unwrappedName := rescopeUnwrappedLocalNamePrefix + name
				for {
					_, used := usedNames[unwrappedName]
					if !used &&
						!symbolInfo.hasExternName(unwrappedName) &&
						unwrappedName != symbolInfo.globalSymbolNamespace {
						break
					}
					unwrappedName = rescopeUnwrappedLocalNamePrefix + unwrappedName
				}
				usedNames[unwrappedName] = struct{}{}
				if unwrapped == nil {
					unwrapped = make(map[string]string)
				}
				unwrapped[name] = unwrappedName
			}
		}
	}
	return rescopeRenameHelper{
		compiler:                           b.compiler,
		globalSymbolNamespace:              b.globalSymbolNamespace,
		disambigPrefix:                     b.disambigPrefix,
		unwrappedLocalNamesInDefiningChunk: unwrapped,
	}
}

func rescopeConflictsWithGlobalSymbolNamespace(name, globalSymbolNamespace, disambigPrefix string) bool {
	return name == globalSymbolNamespace || strings.HasPrefix(name, disambigPrefix)
}

func (h *rescopeRenameHelper) hasUnwrappedLocalNameInDefiningChunk(name string) bool {
	_, ok := h.unwrappedLocalNamesInDefiningChunk[name]
	return ok
}

func (h *rescopeRenameHelper) getLocalNameInDefiningChunk(name string, isGlobalDeclaration bool) string {
	if isGlobalDeclaration {
		if localName, ok := h.unwrappedLocalNamesInDefiningChunk[name]; ok {
			return localName
		}
	}
	return name
}

func (h *rescopeRenameHelper) maybeDisambiguateLocalVar(n rhino.Node, v *Var) {
	name := n.GetString()
	if !v.IsGlobal() && rescopeConflictsWithGlobalSymbolNamespace(name, h.globalSymbolNamespace, h.disambigPrefix) {
		n.SetString(name + rescopeDisambiguationSuffix)
		h.compiler.ReportChangeToEnclosingScope(n)
	}
}

func (h *rescopeRenameHelper) maybeRenameGlobalInDefiningChunk(n rhino.Node) string {
	name := n.GetString()
	if localName, ok := h.unwrappedLocalNamesInDefiningChunk[name]; ok && localName != name {
		n.SetString(localName)
		h.compiler.ReportChangeToEnclosingScope(n)
		return localName
	}
	return name
}
