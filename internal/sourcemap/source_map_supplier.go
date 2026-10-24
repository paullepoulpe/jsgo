// Port of com/google/debugging/sourcemap/SourceMapSupplier.java @ closure-compiler v20261021

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

// SourceMapSupplier is a class for mapping source map names to the actual
// contents. Used when parsing index maps.
type SourceMapSupplier interface {
	// GetSourceMap returns the contents of the map associated with the URL
	// (nil for Java null). It may panic with a jutil IOException.
	//
	// Java: SourceMapSupplier.getSourceMap(String)
	GetSourceMap(url string) *string
}
