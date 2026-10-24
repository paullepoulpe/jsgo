// Port of com/google/debugging/sourcemap/SourceMapGeneratorFactory.java @ closure-compiler v20261021

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

package sourcemap

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SourceMapGeneratorFactory creates source map generators.
type SourceMapGeneratorFactory struct{}

// SourceMapGeneratorFactory_GetInstance returns the appropriate source map
// object for the given source map format.
//
// Java: SourceMapGeneratorFactory.getInstance(SourceMapFormat)
func SourceMapGeneratorFactory_GetInstance(format SourceMapFormat) SourceMapGenerator {
	switch format {
	case SourceMapFormat_DEFAULT, SourceMapFormat_V3:
		return NewSourceMapGeneratorV3()
	}
	// Go-only: language-adapter — unreachable for a valid enum value: the
	// Java switch expression is exhaustive.
	panic(jutil.NewIllegalStateException())
}
