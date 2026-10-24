// Port of com/google/debugging/sourcemap/SourceMappingReversable.java @ closure-compiler v20261021

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

package sourcemap

// SourceMappingReversable is a SourceMapping that can provide the reverse
// (source → target) source mapping.
type SourceMappingReversable interface {
	SourceMapping

	// GetOriginalSources returns the collection of original sources in this
	// source mapping.
	//
	// Java: SourceMappingReversable.getOriginalSources()
	GetOriginalSources() []*string

	// GetReverseMapping returns, given a source file, line, and column, the
	// reverse mapping (source → target). A collection is returned as in some
	// cases (like a function being inlined), one source line may map to more
	// then one target location. An empty collection is returned if there
	// were no matches.
	//
	// Java: SourceMappingReversable.getReverseMapping(String, int, int)
	GetReverseMapping(originalFile string, line, column int) []*Mapping_OriginalMapping
}
