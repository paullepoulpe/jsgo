// Port of com/google/javascript/jscomp/modules/UnresolvedModule.java @ closure-compiler v20261021

/*
 * Copyright 2019 The Closure Compiler Authors.
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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// UnresolvedModule is a module which has had some of its imports and
// exports statements scanned but has yet to resolve anything transitively.
//
// Reference equality is expected in ExportTrace, so implementations must be
// pointer types (Go interface equality is then identity, as the final
// Java equals/hashCode require).
type UnresolvedModule interface {
	// Reset clears any caching so that Resolve(ModuleRequestResolver) will
	// create a new Module.
	//
	// Java: UnresolvedModule.reset()
	Reset()

	// Java: UnresolvedModule.resolve(ModuleRequestResolver)
	Resolve(moduleRequestResolver ModuleRequestResolver) *Module

	// ResolveWithModuleSpecifier resolves all imports and exports and
	// returns a resolved module. moduleSpecifier is the module specifier
	// that was used to import this module, if resolving an import (may be
	// nil).
	//
	// Java: UnresolvedModule.resolve(ModuleRequestResolver, String)
	ResolveWithModuleSpecifier(moduleRequestResolver ModuleRequestResolver, moduleSpecifier *string) *Module

	// Metadata returns the metadata corresponding to this module.
	//
	// Java: UnresolvedModule.metadata()
	Metadata() *ModuleMetadataMap_ModuleMetadata

	// GetExportedNames returns all names in this module's namespace. Names
	// are sorted per Java's string ordering, which should be the same as
	// JavaScript's Array.protype.sort, which is how the spec says these keys
	// should be ordered in the ES module object.
	//
	// Java: UnresolvedModule.getExportedNames(ModuleRequestResolver)
	GetExportedNames(moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashSet[string]

	// GetExportedNamesWithVisited returns all names in this module's
	// namespace. Names are sorted per Java's string ordering, which should be
	// the same as JavaScript's Array.prototype.sort, which is how the spec
	// says these keys should be ordered in the ES module object. visited is
	// the set used to detect export * cycles.
	//
	// Java: UnresolvedModule.getExportedNames(ModuleRequestResolver, Set<UnresolvedModule>)
	GetExportedNamesWithVisited(moduleRequestResolver ModuleRequestResolver, visited *jutil.LinkedHashSet[UnresolvedModule]) *jutil.LinkedHashSet[string]

	// Java: UnresolvedModule.resolveExport(ModuleRequestResolver, String)
	ResolveExport(moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult

	// ResolveExportWithTraces resolves exportName (the name of the export to
	// resolve). moduleSpecifier is the specifier used to reference this
	// module, if this trace is from an import (may be nil); resolveSet is the
	// set used to detect invalid cycles (it is invalid to reach the same
	// exact export (same module with the same export name) in a given
	// cycle); exportStarSet is the set used for cycle checking with export *
	// statements. It returns the result of resolving the export, which can be
	// one of several states:
	//   - The resolved export with the binding, if found.
	//   - A result indicating that the export is ambiguous.
	//   - A result indicating that the module has no such export.
	//   - A result indicating that there was some other error resolving, like
	//     a cycle, or a module transitively returned that there was no such
	//     export.
	//
	// Java: UnresolvedModule.resolveExport(ModuleRequestResolver, String, String, Set<ExportTrace>, Set<UnresolvedModule>)
	ResolveExportWithTraces(
		moduleRequestResolver ModuleRequestResolver,
		moduleSpecifier *string,
		exportName string,
		resolveSet *jutil.LinkedHashSet[ExportTrace],
		exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult

	// Java: UnresolvedModule.equals(Object)
	Equals(other any) bool

	// Java: UnresolvedModule.hashCode()
	HashCode() int32
}

// UnresolvedModuleBase holds the base implementations of UnresolvedModule;
// subclasses embed it and call InitUnresolvedModuleBase.
type UnresolvedModuleBase struct {
	self UnresolvedModule
}

// InitUnresolvedModuleBase is the Go form of super(); self is the object
// being constructed.
func (m *UnresolvedModuleBase) InitUnresolvedModuleBase(self UnresolvedModule) {
	m.self = self
}

// Java: UnresolvedModule.resolve(ModuleRequestResolver)
func (m *UnresolvedModuleBase) Resolve(moduleRequestResolver ModuleRequestResolver) *Module {
	return m.self.ResolveWithModuleSpecifier(moduleRequestResolver /* moduleSpecifier= */, nil)
}

// ResolveExport resolves exportName; see ResolveExportWithTraces.
//
// Java: UnresolvedModule.resolveExport(ModuleRequestResolver, String)
func (m *UnresolvedModuleBase) ResolveExport(
	moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult {
	return m.self.ResolveExportWithTraces(
		moduleRequestResolver,
		/* moduleSpecifier= */ nil,
		exportName,
		jutil.NewLinkedHashSet[ExportTrace](),
		jutil.NewLinkedHashSet[UnresolvedModule]())
}

// Equals is reference equality (final in Java).
//
// Java: UnresolvedModule.equals(Object)
func (m *UnresolvedModuleBase) Equals(other any) bool {
	o, ok := other.(UnresolvedModule)
	return ok && o == m.self
}

// HashCode is the identity hash code (final in Java).
//
// Java: UnresolvedModule.hashCode()
func (m *UnresolvedModuleBase) HashCode() int32 {
	return jutil.IdentityHashCode(m.self)
}
