// Port of com/google/javascript/jscomp/deps/DependencyResolver.java @ closure-compiler v20261021

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

// DependencyResolver is the interface for DependencyResolver to keep in
// line with com.google.javascript.jscomp.deps.DefaultDependencyResolver,
// which doesn't provide an interface. This is so that later on, we can
// merge back with it.
//
// Every method throws ServiceException (as a panic).
type DependencyResolver interface {
	// GetDependencies gets a list of dependencies for the provided code.
	//
	// Java: DependencyResolver.getDependencies(String)
	GetDependencies(code string) []string

	// GetDependenciesOfSymbols gets a list of dependencies for the provided
	// list of symbols.
	//
	// Java: DependencyResolver.getDependencies(Collection<String>)
	GetDependenciesOfSymbols(symbols []string) []string

	// GetDependenciesWithSeen: code is the raw code to be parsed for
	// requires; seen is the set of already seen symbols;
	// addClosureBaseFile indicates whether the closure base file should be
	// added to the dependency list. Returns a list of filenames for each of
	// the dependencies for the provided code.
	//
	// Java: DependencyResolver.getDependencies(String, Set<String>, boolean)
	GetDependenciesWithSeen(code string, seen *jutil.LinkedHashSet[string], addClosureBaseFile bool) []string

	// GetDependenciesOfSymbolsWithSeen: symbols is a list of required
	// symbols; seen is the set of already seen symbols. Returns a list of
	// filenames for each of the required symbols.
	//
	// Java: DependencyResolver.getDependencies(Collection<String>, Set<String>)
	GetDependenciesOfSymbolsWithSeen(symbols []string, seen *jutil.LinkedHashSet[string]) []string
}
