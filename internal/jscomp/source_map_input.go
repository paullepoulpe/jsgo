// Port of com/google/javascript/jscomp/SourceMapInput.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// SourceMapInput is a lazy-loaded SourceMapConsumerV3 instance.
type SourceMapInput struct {
	// Go-only: language-adapter — the object's monitor, which the
	// synchronized getSourceMap holds.
	monitor         jutil.Monitor
	sourceFile      *SourceFile
	parsedSourceMap *sourcemap.SourceMapConsumerV3
	cached          bool
	ready           atomic.Bool
}

var SourceMapInput_SOURCEMAP_RESOLVE_FAILED = DiagnosticType_Warning("SOURCEMAP_RESOLVE_FAILED", "Failed to resolve sourcemap at {0}: {1}")

var SourceMapInput_SOURCEMAP_PARSE_FAILED = DiagnosticType_Warning(
	"SOURCEMAP_PARSE_FAILED", "Failed to parse malformed sourcemap in {0}: {1}")

// Java: SourceMapInput(SourceFile)
func NewSourceMapInput(sourceFile *SourceFile) *SourceMapInput {
	return &SourceMapInput{sourceFile: sourceFile}
}

// GetSourceMap gets the source map, reading from disk and parsing if
// necessary. Returns nil if the sourcemap cannot be resolved or is
// malformed.
//
// Java: SourceMapInput.getSourceMap(ErrorManager)
func (i *SourceMapInput) GetSourceMap(errorManager ErrorManager) *sourcemap.SourceMapConsumerV3 {
	if i.ready.Load() {
		return i.parsedSourceMap
	}
	i.monitor.Lock()
	defer i.monitor.Unlock()
	if !i.cached {
		// Avoid re-reading or reparsing files.
		i.cached = true
		sourceMapPath := i.sourceFile.GetName()
		jutil.Try(func() {
			sourceMapContents := i.sourceFile.GetCode()
			consumer := sourcemap.NewSourceMapConsumerV3()
			consumer.Parse(sourceMapContents)

			// Detect an empty "default" source map provided by Gulp. We do not want
			// to treat this file as having a real source map (implying this is an
			// intermediate source file). Rather, this is the original source file.
			// https://github.com/gulp-sourcemaps/gulp-sourcemaps/blob/91b94954/src/init/index.js#L116-L124
			if consumer.GetFile() != nil &&
				*consumer.GetFile() != "" &&
				len(consumer.GetOriginalSources()) == 1 &&
				sourceMapInput_contains(consumer.GetOriginalSources(), consumer.GetFile()) &&
				len(consumer.GetOriginalSourcesContent()) == 1 &&
				len(consumer.GetOriginalNames()) == 0 &&
				consumer.GetLineCount() == -1 {
				consumer = nil
			}

			i.parsedSourceMap = consumer
			// Go-only: performance — the parsed map holds all that is
			// needed; do not also keep the map's text (read again on demand).
			i.sourceFile.ClearCachedSource()
		}, func(e jutil.IException) {
			switch e.(type) {
			case jutil.IIOException:
				error :=
					JSError_Make(SourceMapInput_SOURCEMAP_RESOLVE_FAILED, sourceMapPath, e.GetMessage())
				errorManager.Report(error.DefaultLevel().Ptr(), error)
			case sourcemap.ISourceMapParseException:
				error :=
					JSError_Make(SourceMapInput_SOURCEMAP_PARSE_FAILED, sourceMapPath, e.GetMessage())
				errorManager.Report(error.DefaultLevel().Ptr(), error)
			default:
				panic(e)
			}
		})
		i.ready.Store(true)
	}
	return i.parsedSourceMap
}

// sourceMapInput_contains is Collection.contains(Object) on a list of
// nullable Strings.
func sourceMapInput_contains(list []*string, o *string) bool {
	for _, s := range list {
		if (s == nil && o == nil) || (s != nil && o != nil && *s == *o) {
			return true
		}
	}
	return false
}

// GetOriginalPath gets the original location of this sourcemap file on
// disk.
//
// Java: SourceMapInput.getOriginalPath()
func (i *SourceMapInput) GetOriginalPath() string {
	return i.sourceFile.GetName()
}

// GetRawSourceMapContents may return nil.
//
// Java: SourceMapInput.getRawSourceMapContents()
func (i *SourceMapInput) GetRawSourceMapContents() (result *string) {
	jutil.Try(func() {
		code := i.sourceFile.GetCode()
		result = &code
	}, func(e jutil.IIOException) {
		result = nil
	})
	return result
}
