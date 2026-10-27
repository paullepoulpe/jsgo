// Port of com/google/javascript/jscomp/Result.java @ closure-compiler v20261021

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
)

// Result holds compilation results.
type Result struct {
	// Java's public final fields; the ones other Go packages read are
	// exported (PORTING.md §3).
	Success  bool
	Errors   []*JSError
	Warnings []*JSError

	variableMap             *VariableMap
	propertyMap             *VariableMap
	namedAnonFunctionMap    *VariableMap
	stringMap               *VariableMap                 // nil is null
	instrumentationMappings *VariableMap                 // nil is null
	SourceMap               *SourceMap                   // nil is null
	cssNames                *jutil.LinkedHashSet[string] // nil is null
	ExternExport            *string                      // nil is null
	idGeneratorMap          *string                      // nil is null
	TranspiledFiles         bool
}

// Java: Result(ImmutableList<JSError>, ImmutableList<JSError>, VariableMap, VariableMap, VariableMap, VariableMap, VariableMap, SourceMap, String, Set<String>, String, boolean)
func NewResultWithAllMaps(
	errors []*JSError,
	warnings []*JSError,
	variableMap *VariableMap,
	propertyMap *VariableMap,
	namedAnonFunctionMap *VariableMap,
	stringMap *VariableMap,
	instrumentationMappings *VariableMap,
	sourceMap *SourceMap,
	externExport *string,
	cssNames *jutil.LinkedHashSet[string],
	idGeneratorMap *string,
	transpiledFiles bool) *Result {
	r := &Result{}
	r.Success = len(errors) == 0
	r.Errors = errors
	r.Warnings = warnings
	r.variableMap = variableMap
	r.propertyMap = propertyMap
	r.namedAnonFunctionMap = namedAnonFunctionMap
	r.stringMap = stringMap
	r.instrumentationMappings = instrumentationMappings
	r.SourceMap = sourceMap
	r.ExternExport = externExport
	r.cssNames = cssNames
	r.idGeneratorMap = idGeneratorMap
	r.TranspiledFiles = transpiledFiles
	return r
}

// NewResult is @VisibleForTesting and deprecated in Java.
//
// Java: Result(ImmutableList<JSError>, ImmutableList<JSError>, VariableMap, VariableMap, VariableMap, SourceMap, String)
func NewResult(
	errors []*JSError,
	warnings []*JSError,
	variableMap *VariableMap,
	propertyMap *VariableMap,
	namedAnonFunctionMap *VariableMap,
	sourceMap *SourceMap,
	externExport *string) *Result {
	return NewResultWithAllMaps(
		errors,
		warnings,
		variableMap,
		propertyMap,
		namedAnonFunctionMap,
		/* stringMap= */ nil,
		/* instrumentationMappings= */ nil,
		sourceMap,
		externExport,
		/* cssNames= */ nil,
		/* idGeneratorMap= */ nil,
		/* transpiledFiles= */ false)
}

// Result_PruneResultForPartialCompilation returns an almost empty result
// that is more appropriate for a partial compilation.
//
// For a partial compilation we only care about errors and warnings. It is
// unnecessary to examine all of the other results.
//
// result is the full `Result` object provided by the compiler.
//
// Java: Result.pruneResultForPartialCompilation(Result)
func Result_PruneResultForPartialCompilation(result *Result) *Result {
	emptyVariableMap := NewVariableMap(jutil.NewLinkedHashMap[string, string]())
	return NewResultWithAllMaps(
		/* errors= */ result.Errors,
		/* warnings= */ result.Warnings,
		/* variableMap= */ emptyVariableMap,
		/* propertyMap= */ emptyVariableMap,
		/* namedAnonFunctionMap= */ emptyVariableMap,
		/* stringMap= */ nil,
		/* instrumentationMappings= */ emptyVariableMap,
		/* sourceMap= */ nil,
		/* externExport= */ jutil.StringPtr(""),
		/* cssNames= */ nil,
		/* idGeneratorMap= */ nil,
		/* transpiledFiles= */ false)
}
