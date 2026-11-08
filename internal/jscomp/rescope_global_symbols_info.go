// Port of com/google/javascript/jscomp/RescopeGlobalSymbolsInfo.java @ closure-compiler v20261104

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

import "github.com/paullepoulpe/jsgo/internal/strpool"

// rescopeSymbolInfo holds cross-chunk and defining-chunk symbol analysis
// information collected during the pre-traversal of RescopeGlobalSymbols.
//
// Java: RescopeGlobalSymbolsInfo
type rescopeSymbolInfo struct {
	globalSymbolNamespace                                  string
	externNameIDs                                          map[strpool.ID]struct{}
	crossChunkNames                                        OrderedSet[string]
	crossChunkNamesWithWriteFromOtherChunk                 OrderedSet[string]
	crossChunkNamesUsedInNonDependentChunk                 OrderedSet[string]
	globalNamesWithReadInDefiningChunk                     map[string]struct{}
	globalNamesWithInnerScopeWriteInDefiningChunk          OrderedSet[string]
	globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk map[string]struct{}
	maybeReferencesThis                                    map[string]struct{}
	wrappedCrossChunkNames                                 OrderedSet[string]
}

func (info *rescopeSymbolInfo) hasExternName(name string) bool {
	if id := strpool.Lookup(name); id > strpool.EmptyID {
		_, ok := info.externNameIDs[id]
		return ok
	}
	return false
}

func (info *rescopeSymbolInfo) populateWrappedCrossChunkNames() {
	for name := range info.crossChunkNamesWithWriteFromOtherChunk.All() {
		if !info.hasExternName(name) {
			info.wrappedCrossChunkNames.Add(name)
		}
	}
	for name := range info.crossChunkNamesUsedInNonDependentChunk.All() {
		if !info.hasExternName(name) {
			info.wrappedCrossChunkNames.Add(name)
		}
	}
	for name := range info.globalNamesWithInnerScopeWriteInDefiningChunk.All() {
		if info.crossChunkNames.Has(name) && !info.hasExternName(name) {
			info.wrappedCrossChunkNames.Add(name)
		}
	}
}
