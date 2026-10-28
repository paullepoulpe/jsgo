// Port of com/google/javascript/jscomp/modules/Export.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/AutoBuilder_Export_Builder.java @ closure-compiler v20261021

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
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Export is an exported name in a module.
//
// See the ExportEntry in the ECMAScript spec
// (https://www.ecma-international.org/ecma-262/9.0/index.html#table-41).
//
// Go-only: language-adapter — the AutoBuilder-generated
// AutoBuilder_Export_Builder is folded into Export_Builder.
type Export struct {
	exportName       *string
	moduleRequest    *string
	importName       *string
	localName        *string
	modulePath       *ModuleLoader_ModulePath
	exportNode       rhino.Node
	nameNode         rhino.Node
	moduleMetadata   *ModuleMetadataMap_ModuleMetadata
	closureNamespace *string
	mutated          bool
	hash             int32
}

// Export_DEFAULT_EXPORT_NAME is the Export.localName() of anonymous ES
// module default exports, e.g. export default 0.
const Export_DEFAULT_EXPORT_NAME = "*default*"

// Export_DEFAULT is the Export.exportName() of anonymous ES module default
// exports, e.g. export default 0.
const Export_DEFAULT = "default"

// Export_NAMESPACE is the Export.exportName() of goog.module default
// exports, e.g. exports = class {};, and the 'namespace' of an ES module
// consisting of all exported names.
const Export_NAMESPACE = "*exports*"

// Java: Export(String, String, String, String, ModulePath, Node, Node, ModuleMetadata, String, boolean)
func NewExport(
	exportName *string,
	moduleRequest *string,
	importName *string,
	localName *string,
	modulePath *ModuleLoader_ModulePath,
	exportNode rhino.Node,
	nameNode rhino.Node,
	moduleMetadata *ModuleMetadataMap_ModuleMetadata,
	closureNamespace *string,
	mutated bool) *Export {
	jutil.CheckNotNullMsg(moduleMetadata, "moduleMetadata")
	return &Export{
		exportName:       exportName,
		moduleRequest:    moduleRequest,
		importName:       importName,
		localName:        localName,
		modulePath:       modulePath,
		exportNode:       exportNode,
		nameNode:         nameNode,
		moduleMetadata:   moduleMetadata,
		closureNamespace: closureNamespace,
		mutated:          mutated,
	}
}

// ExportName may return nil.
//
// Java: Export.exportName()
func (e *Export) ExportName() *string { return e.exportName }

// ModuleRequest may return nil.
//
// Java: Export.moduleRequest()
func (e *Export) ModuleRequest() *string { return e.moduleRequest }

// ImportName may return nil.
//
// Java: Export.importName()
func (e *Export) ImportName() *string { return e.importName }

// LocalName may return nil.
//
// Java: Export.localName()
func (e *Export) LocalName() *string { return e.localName }

// ModulePath may return nil.
//
// Java: Export.modulePath()
func (e *Export) ModulePath() *ModuleLoader_ModulePath { return e.modulePath }

// ExportNode may return nil.
//
// Java: Export.exportNode()
func (e *Export) ExportNode() rhino.Node { return e.exportNode }

// NameNode may return nil.
//
// Java: Export.nameNode()
func (e *Export) NameNode() rhino.Node { return e.nameNode }

// Java: Export.moduleMetadata()
func (e *Export) ModuleMetadata() *ModuleMetadataMap_ModuleMetadata { return e.moduleMetadata }

// ClosureNamespace may return nil.
//
// Java: Export.closureNamespace()
func (e *Export) ClosureNamespace() *string { return e.closureNamespace }

// Java: Export.mutated()
func (e *Export) Mutated() bool { return e.mutated }

func (e *Export) equalsExport(other *Export) bool {
	if e == other {
		return true
	}
	if e == nil || other == nil {
		return false
	}
	return modules_nullableStringEquals(e.exportName, other.exportName) &&
		modules_nullableStringEquals(e.moduleRequest, other.moduleRequest) &&
		modules_nullableStringEquals(e.importName, other.importName) &&
		modules_nullableStringEquals(e.localName, other.localName) &&
		e.modulePath == other.modulePath &&
		e.exportNode == other.exportNode &&
		e.nameNode == other.nameNode &&
		e.moduleMetadata == other.moduleMetadata &&
		modules_nullableStringEquals(e.closureNamespace, other.closureNamespace) &&
		e.mutated == other.mutated
}

// Equals is the record's componentwise equals (ModulePath, Node and
// ModuleMetadata compare by identity).
func (e *Export) Equals(o any) bool {
	other, ok := o.(*Export)
	return ok && other != nil && e.equalsExport(other)
}

// HashCode is the record's hashCode.
//
// Go-only: language-adapter — identity hashes are jutil.IdentityHashCode.
func (e *Export) HashCode() int32 {
	if h := atomic.LoadInt32(&e.hash); h != 0 {
		return h
	}
	return e.computeHashCode()
}

//go:noinline
func (e *Export) computeHashCode() int32 {
	h := int32(0)
	h = 31*h + modules_nullableStringHash(e.exportName)
	h = 31*h + modules_nullableStringHash(e.moduleRequest)
	h = 31*h + modules_nullableStringHash(e.importName)
	h = 31*h + modules_nullableStringHash(e.localName)
	h = 31*h + jutil.IdentityHashCodeOfPointer(e.modulePath)
	h = 31*h + e.exportNode.IdentityHashCode()
	h = 31*h + e.nameNode.IdentityHashCode()
	h = 31*h + jutil.IdentityHashCodeOfPointer(e.moduleMetadata)
	h = 31*h + modules_nullableStringHash(e.closureNamespace)
	if e.mutated {
		h = 31*h + 1231
	} else {
		h = 31*h + 1237
	}
	atomic.StoreInt32(&e.hash, h)
	return h
}

// String is the record's toString.
func (e *Export) String() string {
	modulePath := "null"
	if e.modulePath != nil {
		modulePath = e.modulePath.String()
	}
	exportNode := "null"
	if !e.exportNode.IsNil() {
		exportNode = e.exportNode.String()
	}
	nameNode := "null"
	if !e.nameNode.IsNil() {
		nameNode = e.nameNode.String()
	}
	return "Export[exportName=" + modules_nullableString(e.exportName) +
		", moduleRequest=" + modules_nullableString(e.moduleRequest) +
		", importName=" + modules_nullableString(e.importName) +
		", localName=" + modules_nullableString(e.localName) +
		", modulePath=" + modulePath +
		", exportNode=" + exportNode +
		", nameNode=" + nameNode +
		", moduleMetadata=" + e.moduleMetadata.String() +
		", closureNamespace=" + modules_nullableString(e.closureNamespace) +
		", mutated=" + strconv.FormatBool(e.mutated) + "]"
}

// modules_nullableStringEquals is Objects.equals for nullable Strings.
func modules_nullableStringEquals(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// modules_nullableStringHash is Objects.hashCode for a nullable String.
func modules_nullableStringHash(s *string) int32 {
	if s == nil {
		return 0
	}
	return jutil.StringHashCode(*s)
}

// modules_nullableString is String.valueOf for a nullable String.
func modules_nullableString(s *string) string {
	if s == nil {
		return "null"
	}
	return *s
}

// Export_Builder is Export.Builder with its AutoBuilder implementation.
type Export_Builder struct {
	e     Export
	built bool
}

// ExportName takes a possibly-nil value.
//
// Java: AutoBuilder_Export_Builder.exportName(String)
func (b *Export_Builder) ExportName(exportName *string) *Export_Builder {
	b.e.exportName = exportName
	return b
}

// ModuleRequest takes a possibly-nil value.
//
// Java: AutoBuilder_Export_Builder.moduleRequest(String)
func (b *Export_Builder) ModuleRequest(moduleRequest *string) *Export_Builder {
	b.e.moduleRequest = moduleRequest
	return b
}

// ImportName takes a possibly-nil value.
//
// Java: AutoBuilder_Export_Builder.importName(String)
func (b *Export_Builder) ImportName(importName *string) *Export_Builder {
	b.e.importName = importName
	return b
}

// LocalName takes a possibly-nil value.
//
// Java: AutoBuilder_Export_Builder.localName(String)
func (b *Export_Builder) LocalName(localName *string) *Export_Builder {
	b.e.localName = localName
	return b
}

// ModulePath takes a possibly-nil value.
//
// Java: AutoBuilder_Export_Builder.modulePath(ModuleLoader.ModulePath)
func (b *Export_Builder) ModulePath(modulePath *ModuleLoader_ModulePath) *Export_Builder {
	b.e.modulePath = modulePath
	return b
}

// ExportNode takes a possibly-nil value.
//
// Java: AutoBuilder_Export_Builder.exportNode(Node)
func (b *Export_Builder) ExportNode(exportNode rhino.Node) *Export_Builder {
	b.e.exportNode = exportNode
	return b
}

// NameNode takes a possibly-nil value.
//
// Java: AutoBuilder_Export_Builder.nameNode(Node)
func (b *Export_Builder) NameNode(nameNode rhino.Node) *Export_Builder {
	b.e.nameNode = nameNode
	return b
}

// Java: AutoBuilder_Export_Builder.moduleMetadata(ModuleMetadataMap.ModuleMetadata)
func (b *Export_Builder) ModuleMetadata(moduleMetadata *ModuleMetadataMap_ModuleMetadata) *Export_Builder {
	if moduleMetadata == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null moduleMetadata"))
	}
	b.e.moduleMetadata = moduleMetadata
	return b
}

// ClosureNamespace takes a possibly-nil value.
//
// Java: AutoBuilder_Export_Builder.closureNamespace(String)
func (b *Export_Builder) ClosureNamespace(closureNamespace *string) *Export_Builder {
	b.e.closureNamespace = closureNamespace
	return b
}

// Java: AutoBuilder_Export_Builder.mutated(boolean)
func (b *Export_Builder) Mutated(mutated bool) *Export_Builder {
	b.e.mutated = mutated
	return b
}

// Java: AutoBuilder_Export_Builder.autoBuild()
func (b *Export_Builder) AutoBuild() *Export {
	if b.e.moduleMetadata == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties: moduleMetadata"))
	}
	if !b.built {
		b.built = true
		return &b.e
	}
	cp := b.e
	cp.hash = 0
	return &cp
}

// Java: Export.Builder.build()
func (b *Export_Builder) Build() *Export {
	e := b.AutoBuild()
	if e.ModuleMetadata().IsEs6Module() {
		b.ValidateEsModule(e)
	} else if e.ModuleMetadata().IsGoogModule() {
		Export_Builder_ValidateGoogModule(e)
	} else {
		Export_Builder_ValidateOtherModule(e)
	}
	return e
}

// ValidateEsModule validates an export from an ES module.
//
// Java: Export.Builder.validateEsModule(Export)
func (b *Export_Builder) ValidateEsModule(e *Export) {
	jutil.CheckState(e.ClosureNamespace() == nil)

	jutil.CheckStateMsg(
		!(e.ImportName() != nil && *e.ImportName() == "*") ||
			(e.ModuleRequest() != nil && e.ExportName() == nil && e.LocalName() == nil),
		"Star exports should not have exported / local names.")

	jutil.CheckStateMsg(
		e.LocalName() == nil || e.ModuleRequest() == nil,
		"Local exports should not have module requests.")
	jutil.CheckStateMsg(
		e.ModuleRequest() == nil || e.LocalName() == nil,
		"Reexports should not have local names.")

	jutil.CheckStateMsg(
		e.ModuleRequest() == nil || e.ImportName() != nil,
		"Reexports should have import names.")
	jutil.CheckStateMsg(
		e.ImportName() == nil || e.ModuleRequest() != nil,
		"Exports with an import name should be a reexport.")
}

// Export_Builder_ValidateGoogModule validates some export from a goog
// module.
//
// Java: Export.Builder.validateGoogModule(Export)
func Export_Builder_ValidateGoogModule(e *Export) {
	jutil.CheckStateMsg(e.ClosureNamespace() != nil, "Exports should be associated with a namespace")
	jutil.CheckStateMsg(e.ExportName() != nil, "Exports should be named")
	jutil.CheckStateMsg(!e.ExportNode().IsNil(), "Exports should have a node")
	jutil.CheckStateMsg(e.LocalName() == nil, "goog.module Exports don't set a localName")
	jutil.CheckStateMsg(e.ModuleRequest() == nil, "goog modules cannot export from other modules")
}

// Export_Builder_ValidateOtherModule validates some faux export from a
// non-ES module.
//
// Java: Export.Builder.validateOtherModule(Export)
func Export_Builder_ValidateOtherModule(e *Export) {
	jutil.CheckNotNull(e.ExportName())

	// Fields ignored for these fake exports. Should not set these.
	jutil.CheckState(e.ExportNode().IsNil())
	jutil.CheckState(e.LocalName() == nil)
	jutil.CheckState(e.ModuleRequest() == nil)
	jutil.CheckState(e.ImportName() == nil)
	jutil.CheckState(e.NameNode().IsNil())
}

// Java: Export.builder()
func Export_Builder_() *Export_Builder {
	return &Export_Builder{}
}

// Java: Export.toBuilder()
func (e *Export) ToBuilder() *Export_Builder {
	b := &Export_Builder{e: *e}
	b.e.hash = 0
	return b
}

// MutatedCopy returns a copy of this export that has the mutated() bit set.
//
// Java: Export.mutatedCopy()
func (e *Export) MutatedCopy() *Export {
	if e.moduleMetadata == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties: moduleMetadata"))
	}
	cp := *e
	cp.mutated = true
	cp.hash = 0
	return &cp
}
