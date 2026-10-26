// Port of com/google/javascript/jscomp/deps/DefaultDependencyResolver.java @ closure-compiler v20261021

/*
 * Copyright 2006 The Closure Compiler Authors.
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

// DefaultDependencyResolver is the class for resolving Closure
// dependencies.
//
// Given a valid deps.js file the dependency tree is parsed and stored in
// memory. The DependencyResolver can then be used to calculate the full list
// of transitive dependencies from: a block of code (GetDependencies), or a
// list of symbols (GetDependenciesOfSymbols).
type DefaultDependencyResolver struct {
	// Source files used to look up the dependencies.
	depsFiles []*DependencyFile

	// Flag that determines if the resolver will strictly require all
	// goog.requires.
	strictRequires bool
}

var _ DependencyResolver = (*DefaultDependencyResolver)(nil)

// DefaultDependencyResolver_CLOSURE_BASE is the filename for Closure's
// base.js file which is always added.
const DefaultDependencyResolver_CLOSURE_BASE = "base.js"

// DefaultDependencyResolver_CLOSURE_BASE_PROVIDE is the provide for
// Closure's base.js.
const DefaultDependencyResolver_CLOSURE_BASE_PROVIDE = "goog"

// Logger for DependencyResolver.
var defaultDependencyResolver_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.deps.DefaultDependencyResolver")

// NewDefaultDependencyResolver creates a new dependency resolver. depsFiles
// is the list of deps files; strictRequires determines if the resolver will
// throw an exception on a missing dependency.
//
// Java: DefaultDependencyResolver(List<DependencyFile>, boolean)
func NewDefaultDependencyResolver(depsFiles []*DependencyFile, strictRequires bool) *DefaultDependencyResolver {
	return &DefaultDependencyResolver{depsFiles: depsFiles, strictRequires: strictRequires}
}

// GetDependencies gets a list of dependencies for the provided code.
//
// Java: DefaultDependencyResolver.getDependencies(String)
func (r *DefaultDependencyResolver) GetDependencies(code string) []string {
	return r.GetDependenciesOfSymbols(defaultDependencyResolver_parseRequires(code, true))
}

// GetDependenciesOfSymbols gets a list of dependencies for the provided list
// of symbols.
//
// Java: DefaultDependencyResolver.getDependencies(Collection<String>)
func (r *DefaultDependencyResolver) GetDependenciesOfSymbols(symbols []string) []string {
	// Go-only: language-adapter — the HashSet is only queried (containsAll,
	// addAll), never iterated, so insertion order is unobservable.
	return r.GetDependenciesOfSymbolsWithSeen(symbols, jutil.NewLinkedHashSet[string]())
}

// GetDependenciesWithSeen: code is the raw code to be parsed for requires;
// seen is the set of already seen symbols; addClosureBaseFile indicates
// whether the closure base file should be added to the dependency list.
// Returns a list of filenames for each of the dependencies for the provided
// code.
//
// Java: DefaultDependencyResolver.getDependencies(String, Set<String>, boolean)
func (r *DefaultDependencyResolver) GetDependenciesWithSeen(
	code string, seen *jutil.LinkedHashSet[string], addClosureBaseFile bool) []string {
	return r.GetDependenciesOfSymbolsWithSeen(defaultDependencyResolver_parseRequires(code, addClosureBaseFile), seen)
}

// GetDependenciesOfSymbolsWithSeen: symbols is a list of required symbols;
// seen is the set of already seen symbols. Returns a list of filenames for
// each of the required symbols.
//
// Java: DefaultDependencyResolver.getDependencies(Collection<String>, Set<String>)
func (r *DefaultDependencyResolver) GetDependenciesOfSymbolsWithSeen(
	symbols []string, seen *jutil.LinkedHashSet[string]) []string {
	var list []string
	for _, depsFile := range r.depsFiles {
		depsFile.EnsureUpToDate()
	}

	for _, symbol := range symbols {
		r.addDependency(symbol, seen, &list)
	}
	return list
}

// addDependency adds all the transitive dependencies for a symbol to the
// provided list. The set is used to avoid adding dupes while keeping the
// correct order. NOTE: Use of a LinkedHashSet would require reversing the
// results to get correct dependency ordering.
//
// Java: DefaultDependencyResolver.addDependency(String, Set<String>, List<String>)
func (r *DefaultDependencyResolver) addDependency(symbol string, seen *jutil.LinkedHashSet[string], list *[]string) {
	dependency := r.GetDependencyInfo(symbol)
	if dependency == nil {
		if r.strictRequires {
			panic(NewServiceException("Unknown require of " + symbol))
		}
	} else if !defaultDependencyResolver_containsAll(seen, dependency.GetProvides()) {
		seen.AddAll(dependency.GetProvides()...)
		for _, require := range dependency.GetRequiredSymbols() {
			r.addDependency(require, seen, list)
		}
		*list = append(*list, dependency.GetPathRelativeToClosureBase())
	}
}

// Go-only: language-adapter — Set.containsAll(Collection).
func defaultDependencyResolver_containsAll(seen *jutil.LinkedHashSet[string], elems []string) bool {
	for _, e := range elems {
		if !seen.Contains(e) {
			return false
		}
	}
	return true
}

// Parses a block of code for goog.require statements and extracts the
// required symbols.
//
// Java: DefaultDependencyResolver.parseRequires(String, boolean)
func defaultDependencyResolver_parseRequires(code string, addClosureBase bool) []string {
	var errorManager ErrorManager = NewLoggerErrorManager(defaultDependencyResolver_logger)
	parser := NewJsFileRegexParser(errorManager)
	deps := parser.ParseFile("<unknown path>", "<unknown path>", code)
	requires := []string{}
	if addClosureBase {
		requires = append(requires, DefaultDependencyResolver_CLOSURE_BASE_PROVIDE)
	}
	requires = append(requires, deps.GetRequiredSymbols()...)
	errorManager.GenerateReport()
	return requires
}

// GetDependencyInfo looks at each of the dependency files for dependency
// information.
//
// Java: DefaultDependencyResolver.getDependencyInfo(String)
func (r *DefaultDependencyResolver) GetDependencyInfo(symbol string) DependencyInfo {
	for _, depsFile := range r.depsFiles {
		di := depsFile.GetDependencyInfo(symbol)
		if di != nil {
			return di
		}
	}
	return nil
}
