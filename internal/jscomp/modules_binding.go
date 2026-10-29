// Port of com/google/javascript/jscomp/modules/Binding.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Binding represents a variable bound by an import or export statement, or
// goog.require. This can either be a single variable or a entire module
// namespace created by an import * statement.
//
// See Module#namespace() and Module#boundNames() for how Bindings are used.
//
// metadata is the metadata of the module this is bound to: if this was made
// from an import * then it is the module that this imported, otherwise it
// is the module this binding was created in (e.g. the module of the
// originating export). sourceNode is the AST node to use for source location
// when rewriting: this is generally a NAME or IMPORT_STAR node inside an
// import or export statement that represents where the name was bound;
// however as export * from has no NAME nodes the source node in that
// instance should be the entire export node; nil for missing ES modules and
// non-ES modules as they are currently not scanned. originatingExport is the
// original export if this binding was created by an export (for transitive
// exports this will still be the *original* export, not the transitive
// link). isModuleNamespace is true if this represents a module namespace,
// e.g. created by import *.
type Binding struct {
	metadata            *ModuleMetadataMap_ModuleMetadata
	sourceNode          rhino.Node
	originatingExport   *Export
	closureNamespace    string
	createdBy           Binding_CreatedBy
	isModuleNamespace   bool
	hasClosureNamespace bool
}

// Java: Binding(ModuleMetadata, Node, Export, boolean, String, CreatedBy)
func NewBinding(
	metadata *ModuleMetadataMap_ModuleMetadata,
	sourceNode rhino.Node,
	originatingExport *Export,
	isModuleNamespace bool,
	closureNamespace string,
	hasClosureNamespace bool,
	createdBy Binding_CreatedBy) *Binding {
	jutil.CheckNotNullMsg(metadata, "metadata")
	return &Binding{
		metadata:            metadata,
		sourceNode:          sourceNode,
		originatingExport:   originatingExport,
		closureNamespace:    closureNamespace,
		createdBy:           createdBy,
		isModuleNamespace:   isModuleNamespace,
		hasClosureNamespace: hasClosureNamespace,
	}
}

// Java: Binding.metadata()
func (b *Binding) Metadata() *ModuleMetadataMap_ModuleMetadata { return b.metadata }

// SourceNode may return nil.
//
// Java: Binding.sourceNode()
func (b *Binding) SourceNode() rhino.Node { return b.sourceNode }

// OriginatingExport may return nil.
//
// Java: Binding.originatingExport()
func (b *Binding) OriginatingExport() *Export { return b.originatingExport }

// Java: Binding.isModuleNamespace()
func (b *Binding) IsModuleNamespace() bool { return b.isModuleNamespace }

// Java: Binding.closureNamespace()
func (b *Binding) ClosureNamespace() string { return b.closureNamespace }

// HasClosureNamespace reports whether closureNamespace was set.
func (b *Binding) HasClosureNamespace() bool { return b.hasClosureNamespace }

// Java: Binding.createdBy()
func (b *Binding) CreatedBy() Binding_CreatedBy { return b.createdBy }

// Binding_CreatedBy lists the different ways that Bindings can be created.
//
// Java: Binding.CreatedBy
type Binding_CreatedBy int32

const (
	// Binding_CreatedBy_EXPORT is a binding created by an export in an ES
	// module.
	//
	//	export const x = 0;
	//	export function x() {}
	//	export {x};
	//	export * from ''; // creates bindings in the namespace
	Binding_CreatedBy_EXPORT Binding_CreatedBy = iota
	// Binding_CreatedBy_IMPORT is a binding created by an ES import.
	//
	//	import {x} from '';
	//	import * as x from '';
	//	import x from '';
	Binding_CreatedBy_IMPORT
	// Binding_CreatedBy_GOOG_REQUIRE is a binding created by a goog.require
	// statement.
	//
	//	const x = goog.require();
	//	const {x} = goog.require();
	Binding_CreatedBy_GOOG_REQUIRE
	// Binding_CreatedBy_GOOG_REQUIRE_TYPE is a binding created by a
	// goog.requireType statement.
	//
	//	const x = goog.requireType();
	//	const {x} = goog.requireType();
	Binding_CreatedBy_GOOG_REQUIRE_TYPE
	// Binding_CreatedBy_GOOG_FORWARD_DECLARE is a binding created by a
	// goog.forwardDeclare statement.
	//
	//	const x = goog.forwardDeclare();
	//	const {x} = goog.forwardDeclare();
	Binding_CreatedBy_GOOG_FORWARD_DECLARE
)

var binding_CreatedBy_names = [...]string{"EXPORT", "IMPORT", "GOOG_REQUIRE", "GOOG_REQUIRE_TYPE", "GOOG_FORWARD_DECLARE"}

// String is Enum.name().
func (c Binding_CreatedBy) String() string { return binding_CreatedBy_names[c] }

// Ordinal is Enum.ordinal().
func (c Binding_CreatedBy) Ordinal() int { return int(c) }

// IsClosureImport is whether this is some goog.* dependency import.
//
// Java: Binding.CreatedBy.isClosureImport()
func (c Binding_CreatedBy) IsClosureImport() bool {
	return c == Binding_CreatedBy_GOOG_REQUIRE ||
		c == Binding_CreatedBy_GOOG_REQUIRE_TYPE ||
		c == Binding_CreatedBy_GOOG_FORWARD_DECLARE
}

// Binding_From is the binding for an exported value that is not a module
// namespace object.
//
// Java: Binding.from(Export, Node)
func Binding_From(boundExport *Export, sourceNode rhino.Node) *Binding {
	return NewBinding(
		boundExport.ModuleMetadata(),
		sourceNode,
		boundExport,
		/* isModuleNamespace= */ false,
		boundExport.ClosureNamespace(),
		boundExport.HasClosureNamespace(),
		Binding_CreatedBy_EXPORT)
}

// Binding_FromClosureImport is the binding for an entire module namespace
// created by const x = goog.require(Type)('...').
//
// Java: Binding.from(ModuleMetadata, Node, String, CreatedBy)
func Binding_FromClosureImport(
	metadata *ModuleMetadataMap_ModuleMetadata, sourceNode rhino.Node, closureNamespace string, createdBy Binding_CreatedBy) *Binding {
	jutil.CheckArgumentMsg1(
		createdBy.IsClosureImport(),
		"Expected goog.require(Type) or goog.forwardDeclare, got %s",
		createdBy)
	return NewBinding(
		metadata,
		sourceNode,
		/* originatingExport= */ nil,
		/* isModuleNamespace= */ true,
		closureNamespace,
		true,
		createdBy)
}

// Binding_FromImportStar is the binding for an entire module namespace
// created by an import *.
//
// Java: Binding.from(ModuleMetadata, String, Node)
func Binding_FromImportStar(
	metadataOfBoundModule *ModuleMetadataMap_ModuleMetadata,
	closureNamespace string,
	hasClosureNamespace bool,
	sourceNode rhino.Node) *Binding {
	return NewBinding(
		metadataOfBoundModule,
		sourceNode,
		/* originatingExport= */ nil,
		/* isModuleNamespace= */ true,
		closureNamespace,
		hasClosureNamespace,
		Binding_CreatedBy_IMPORT)
}

// Copy copies the binding with a new source node and CreatedBy binding.
//
// Java: Binding.copy(Node, CreatedBy)
func (b *Binding) Copy(sourceNode rhino.Node, createdBy Binding_CreatedBy) *Binding {
	rhino.CheckNotNull(sourceNode)
	cp := *b
	cp.sourceNode = sourceNode
	cp.createdBy = createdBy
	return &cp
}

// BoundName returns the name of the variable this export is bound to, assuming
// it is not a binding of a module namespace, and whether it is present.
//
// Java: Binding.boundName()
func (b *Binding) BoundName() (string, bool) {
	jutil.CheckState(!b.IsModuleNamespace())
	return b.OriginatingExport().LocalName(), b.OriginatingExport().HasLocalName()
}

// IsMutated returns whether or not this export is potentially mutated after
// module execution (i.e. in a function scope).
//
// Java: Binding.isMutated()
func (b *Binding) IsMutated() bool {
	// Module namespaces can never be mutated. They are always imported, and import bound names
	// are const.
	return !b.IsModuleNamespace() && b.OriginatingExport().Mutated()
}

// IsCreatedByEsImport returns whether this Binding originated from an ES
// import, as opposed to an export or goog.require.
//
// Java: Binding.isCreatedByEsImport()
func (b *Binding) IsCreatedByEsImport() bool {
	return b.CreatedBy() == Binding_CreatedBy_IMPORT
}

// IsCreatedByEsExport returns whether this Binding originated from an ES
// import, as opposed to an export or goog.require.
//
// Java: Binding.isCreatedByEsExport()
func (b *Binding) IsCreatedByEsExport() bool {
	return b.CreatedBy() == Binding_CreatedBy_EXPORT
}

// IsSomeImport returns whether this Binding originated from an ES import or
// goog.require.
//
// Java: Binding.isSomeImport()
func (b *Binding) IsSomeImport() bool {
	return b.CreatedBy() != Binding_CreatedBy_EXPORT
}
