// Port of com/google/javascript/jscomp/modules/Module.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/AutoBuilder_Module_Builder.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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

// Module holds information for modules, particularly ES modules, that is
// useful for rewriting. The primary pieces of information are what
// variables are exported (transitive or local), and what names are
// imported.
//
// Go-only: language-adapter — closureNamespace is stored unboxed on Module
// (returning an interior pointer from ClosureNamespace when present), and
// AutoBuilder_Module_Builder is replaced by NewModule / WithBoundNames.
type Module struct {
	metadata               *ModuleMetadataMap_ModuleMetadata
	path                   *ModuleLoader_ModulePath
	namespace              *jutil.LinkedHashMap[string, *Binding]
	boundNames             *jutil.LinkedHashMap[string, *Binding]
	localNameToLocalExport *jutil.LinkedHashMap[string, *Export]
	closureNamespace       string
	hasClosureNamespace    bool
}

// Java: Module(ModuleMetadata, ModulePath, ImmutableMap<String, Binding>, ImmutableMap<String, Binding>, ImmutableMap<String, Export>, String)
func NewModule(
	metadata *ModuleMetadataMap_ModuleMetadata,
	path *ModuleLoader_ModulePath,
	namespace *jutil.LinkedHashMap[string, *Binding],
	boundNames *jutil.LinkedHashMap[string, *Binding],
	localNameToLocalExport *jutil.LinkedHashMap[string, *Export],
	closureNamespace string,
	hasClosureNamespace bool) *Module {
	jutil.CheckNotNullMsg(metadata, "metadata")
	jutil.CheckNotNullMsg(namespace, "namespace")
	jutil.CheckNotNullMsg(boundNames, "boundNames")
	jutil.CheckNotNullMsg(localNameToLocalExport, "localNameToLocalExport")
	return &Module{
		metadata:               metadata,
		path:                   path,
		namespace:              namespace,
		boundNames:             boundNames,
		localNameToLocalExport: localNameToLocalExport,
		closureNamespace:       closureNamespace,
		hasClosureNamespace:    hasClosureNamespace,
	}
}

// Java: Module.metadata()
func (m *Module) Metadata() *ModuleMetadataMap_ModuleMetadata { return m.metadata }

// Path may return nil.
//
// Java: Module.path()
func (m *Module) Path() *ModuleLoader_ModulePath { return m.path }

// Java: Module.namespace()
func (m *Module) Namespace() *jutil.LinkedHashMap[string, *Binding] { return m.namespace }

// Java: Module.boundNames()
func (m *Module) BoundNames() *jutil.LinkedHashMap[string, *Binding] { return m.boundNames }

// Java: Module.localNameToLocalExport()
func (m *Module) LocalNameToLocalExport() *jutil.LinkedHashMap[string, *Export] {
	return m.localNameToLocalExport
}

// ClosureNamespace returns an interior pointer to closureNamespace when
// present, or nil when absent (zero extra heap allocation).
//
// Java: Module.closureNamespace()
func (m *Module) ClosureNamespace() *string {
	if !m.hasClosureNamespace {
		return nil
	}
	return &m.closureNamespace
}

// ClosureNamespaceString returns the unboxed closureNamespace string ("" when absent).
func (m *Module) ClosureNamespaceString() string { return m.closureNamespace }

// HasClosureNamespace reports whether closureNamespace was set.
func (m *Module) HasClosureNamespace() bool { return m.hasClosureNamespace }

// WithBoundNames returns a shallow copy of m with boundNames replaced.
func (m *Module) WithBoundNames(boundNames *jutil.LinkedHashMap[string, *Binding]) *Module {
	if boundNames == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null boundNames"))
	}
	cp := *m
	cp.boundNames = boundNames
	return &cp
}
