// Port of com/google/javascript/jscomp/SourceFileMapping.java @ closure-compiler v20261021

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// SourceFileMapping maps a source file, line, and column into an
// OriginalMapping.
//
// See com.google.debugging.sourcemap.SourceMapping.
type SourceFileMapping interface {
	// GetSourceMapping returns the original mapping for the file name, line
	// number and column position found in the source map. Returns null if
	// the source file does not have a input map, and returns
	// OriginalMapping.getDefaultInstance() if the input map does not map
	// the location.
	//
	// lineNo is the line number, 1-based; columnNo the column index,
	// 1-based.
	//
	// Java: SourceFileMapping.getSourceMapping(String, int, int)
	GetSourceMapping(fileName string, lineNo int, columnNo int) *sourcemap.Mapping_OriginalMapping
}
