// Port of com/google/javascript/jscomp/deps/DependencyFile.java @ closure-compiler v20260415

/*
 * Copyright 2008 The Closure Compiler Authors.
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

// DependencyFile is a SourceFile containing dependency information.
// Delegates file handling to another SourceFile such that a VirtualFile,
// LocalFile or RemoteFile can be used.
type DependencyFile struct {
	// Map of name spaces to their dependency info.
	//
	// Go-only: language-adapter — a Java HashMap that is only looked up
	// (get, put, clear, isEmpty), never iterated, so insertion order is
	// unobservable.
	dependencies *jutil.LinkedHashMap[string, DependencyInfo]

	// A source file to delegate functionality too.
	delegate Deps_SourceFile
}

var _ Deps_SourceFile = (*DependencyFile)(nil)

// Logger for DependencyResolver.
var dependencyFile_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.deps.DependencyFile")

// NewDependencyFile creates a new dependency file.
//
// Java: DependencyFile(SourceFile)
func NewDependencyFile(delegate Deps_SourceFile) *DependencyFile {
	return &DependencyFile{
		dependencies: jutil.NewLinkedHashMap[string, DependencyInfo](),
		delegate:     delegate,
	}
}

// Java: DependencyFile.getName()
func (f *DependencyFile) GetName() string {
	return f.delegate.GetName()
}

// Java: DependencyFile.getContent()
func (f *DependencyFile) GetContent() string {
	return f.delegate.GetContent()
}

// Java: DependencyFile.wasModified()
func (f *DependencyFile) WasModified() bool {
	return f.delegate.WasModified()
}

// EnsureUpToDate ensures that the dependency graph is up to date and reloads
// the graph if necessary.
//
// Java: DependencyFile.ensureUpToDate()
func (f *DependencyFile) EnsureUpToDate() {
	if f.dependencies.IsEmpty() || f.WasModified() {
		f.LoadGraph()
	}
}

// GetDependencyInfo gets the dependency info for the provided symbol, if
// contained in this dependency file.
//
// Java: DependencyFile.getDependencyInfo(String)
func (f *DependencyFile) GetDependencyInfo(symbol string) DependencyInfo {
	return f.dependencies.Get(symbol)
}

// LoadGraph loads the dependency graph.
//
// Java: DependencyFile.loadGraph()
func (f *DependencyFile) LoadGraph() {
	f.dependencies.Clear()

	dependencyFile_logger.Info("Loading dependency graph")

	// Parse the deps.js file.
	var errorManager ErrorManager = NewLoggerErrorManager(dependencyFile_logger)
	parser := NewDepsFileRegexParser(errorManager)
	depInfos := parser.ParseFileWithContents(f.GetName(), f.GetContent())

	// Ensure the parse succeeded.
	if !parser.DidParseSucceed() {
		panic(NewServiceException("Problem parsing " + f.GetName() + ". See logs for details."))
	}
	// Incorporate the dependencies into our maps.
	for _, depInfo := range depInfos {
		for _, provide := range depInfo.GetProvides() {
			existing := f.dependencies.Get(provide)
			if existing != nil && !jutil.ObjectsEquals(existing, depInfo) {
				panic(NewServiceException("Duplicate provide of " + provide +
					". Was provided by " + existing.GetPathRelativeToClosureBase() +
					" and " + depInfo.GetPathRelativeToClosureBase()))
			}
			f.dependencies.Put(provide, depInfo)
		}
	}

	provides := []string{}
	provides = append(provides, DefaultDependencyResolver_CLOSURE_BASE_PROVIDE)

	// Add implicit base.js entry.
	f.dependencies.Put(
		DefaultDependencyResolver_CLOSURE_BASE_PROVIDE,
		SimpleDependencyInfo_Builder_(DefaultDependencyResolver_CLOSURE_BASE, DefaultDependencyResolver_CLOSURE_BASE).
			SetProvides(provides).Build())
	errorManager.GenerateReport()

	dependencyFile_logger.Info("Dependencies loaded")
}
