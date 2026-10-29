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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type exportPresence uint8

const (
	exportHasExportName exportPresence = 1 << iota
	exportHasModuleRequest
	exportHasImportName
	exportHasLocalName
	exportHasClosureNamespace
)

// Export is an exported name in a module.
//
// See the ExportEntry in the ECMAScript spec
// (https://www.ecma-international.org/ecma-262/9.0/index.html#table-41).
//
// Go-only: language-adapter — the AutoBuilder-generated
// AutoBuilder_Export_Builder is folded into the stack-allocated Export_Builder,
// and nullable strings are unboxed with a compact presence bitmask.
type Export struct {
	exportName       string
	moduleRequest    string
	importName       string
	localName        string
	closureNamespace string
	modulePath       *ModuleLoader_ModulePath
	exportNode       rhino.Node
	nameNode         rhino.Node
	moduleMetadata   *ModuleMetadataMap_ModuleMetadata
	present          exportPresence
	mutated          bool
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

// Java: Export.exportName()
func (e *Export) ExportName() string { return e.exportName }

// HasExportName reports whether exportName was set (distinguishing "" from absent).
func (e *Export) HasExportName() bool { return e.present&exportHasExportName != 0 }

// Java: Export.moduleRequest()
func (e *Export) ModuleRequest() string { return e.moduleRequest }

// HasModuleRequest reports whether moduleRequest was set (distinguishing "" from absent).
func (e *Export) HasModuleRequest() bool { return e.present&exportHasModuleRequest != 0 }

// Java: Export.importName()
func (e *Export) ImportName() string { return e.importName }

// HasImportName reports whether importName was set (distinguishing "" from absent).
func (e *Export) HasImportName() bool { return e.present&exportHasImportName != 0 }

// Java: Export.localName()
func (e *Export) LocalName() string { return e.localName }

// HasLocalName reports whether localName was set (distinguishing "" from absent).
func (e *Export) HasLocalName() bool { return e.present&exportHasLocalName != 0 }

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

// Java: Export.closureNamespace()
func (e *Export) ClosureNamespace() string { return e.closureNamespace }

// HasClosureNamespace reports whether closureNamespace was set (distinguishing "" from absent).
func (e *Export) HasClosureNamespace() bool { return e.present&exportHasClosureNamespace != 0 }

// Java: Export.mutated()
func (e *Export) Mutated() bool { return e.mutated }

// modules_nullableString is String.valueOf for a nullable String.
func modules_nullableString(s *string) string {
	if s == nil {
		return "null"
	}
	return *s
}

// modules_optString formats an unboxed optional string as "null" when absent.
func modules_optString(s string, ok bool) string {
	if !ok {
		return "null"
	}
	return s
}

// Export_Builder is a stack-allocated builder for Export.
type Export_Builder struct {
	e Export
}

// Java: Export.builder()
func Export_Builder_() Export_Builder {
	return Export_Builder{}
}

// Java: AutoBuilder_Export_Builder.exportName(String)
func (b Export_Builder) ExportName(exportName string) Export_Builder {
	b.e.exportName = exportName
	b.e.present |= exportHasExportName
	return b
}

// Java: AutoBuilder_Export_Builder.moduleRequest(String)
func (b Export_Builder) ModuleRequest(moduleRequest string) Export_Builder {
	b.e.moduleRequest = moduleRequest
	b.e.present |= exportHasModuleRequest
	return b
}

// Java: AutoBuilder_Export_Builder.importName(String)
func (b Export_Builder) ImportName(importName string) Export_Builder {
	b.e.importName = importName
	b.e.present |= exportHasImportName
	return b
}

// Java: AutoBuilder_Export_Builder.localName(String)
func (b Export_Builder) LocalName(localName string) Export_Builder {
	b.e.localName = localName
	b.e.present |= exportHasLocalName
	return b
}

// ModulePath takes a possibly-nil value.
//
// Java: AutoBuilder_Export_Builder.modulePath(ModuleLoader.ModulePath)
func (b Export_Builder) ModulePath(modulePath *ModuleLoader_ModulePath) Export_Builder {
	b.e.modulePath = modulePath
	return b
}

// ExportNode takes a possibly-nil value.
//
// Java: AutoBuilder_Export_Builder.exportNode(Node)
func (b Export_Builder) ExportNode(exportNode rhino.Node) Export_Builder {
	b.e.exportNode = exportNode
	return b
}

// NameNode takes a possibly-nil value.
//
// Java: AutoBuilder_Export_Builder.nameNode(Node)
func (b Export_Builder) NameNode(nameNode rhino.Node) Export_Builder {
	b.e.nameNode = nameNode
	return b
}

// Java: AutoBuilder_Export_Builder.moduleMetadata(ModuleMetadataMap.ModuleMetadata)
func (b Export_Builder) ModuleMetadata(moduleMetadata *ModuleMetadataMap_ModuleMetadata) Export_Builder {
	if moduleMetadata == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null moduleMetadata"))
	}
	b.e.moduleMetadata = moduleMetadata
	return b
}

// Java: AutoBuilder_Export_Builder.closureNamespace(String)
func (b Export_Builder) ClosureNamespace(closureNamespace string) Export_Builder {
	b.e.closureNamespace = closureNamespace
	b.e.present |= exportHasClosureNamespace
	return b
}

// Java: Export.Builder.build()
func (b Export_Builder) Build() *Export {
	if b.e.moduleMetadata == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties: moduleMetadata"))
	}
	e := b.e
	if e.moduleMetadata.IsEs6Module() {
		export_validateEsModule(&e)
	} else if e.moduleMetadata.IsGoogModule() {
		export_validateGoogModule(&e)
	} else {
		export_validateOtherModule(&e)
	}
	return &e
}

// export_validateEsModule validates an export from an ES module.
//
// Java: Export.Builder.validateEsModule(Export)
func export_validateEsModule(e *Export) {
	jutil.CheckState(!e.HasClosureNamespace())

	jutil.CheckStateMsg(
		!(e.HasImportName() && e.ImportName() == "*") ||
			(e.HasModuleRequest() && !e.HasExportName() && !e.HasLocalName()),
		"Star exports should not have exported / local names.")

	jutil.CheckStateMsg(
		!e.HasLocalName() || !e.HasModuleRequest(),
		"Local exports should not have module requests.")

	jutil.CheckStateMsg(
		!e.HasModuleRequest() || e.HasImportName(),
		"Reexports should have import names.")
	jutil.CheckStateMsg(
		!e.HasImportName() || e.HasModuleRequest(),
		"Exports with an import name should be a reexport.")
}

// export_validateGoogModule validates some export from a goog
// module.
//
// Java: Export.Builder.validateGoogModule(Export)
func export_validateGoogModule(e *Export) {
	jutil.CheckStateMsg(e.HasClosureNamespace(), "Exports should be associated with a namespace")
	jutil.CheckStateMsg(e.HasExportName(), "Exports should be named")
	jutil.CheckStateMsg(!e.ExportNode().IsNil(), "Exports should have a node")
	jutil.CheckStateMsg(!e.HasLocalName(), "goog.module Exports don't set a localName")
	jutil.CheckStateMsg(!e.HasModuleRequest(), "goog modules cannot export from other modules")
}

// export_validateOtherModule validates some faux export from a
// non-ES module.
//
// Java: Export.Builder.validateOtherModule(Export)
func export_validateOtherModule(e *Export) {
	if !e.HasExportName() {
		panic(jutil.NewNullPointerException())
	}

	// Fields ignored for these fake exports. Should not set these.
	jutil.CheckState(e.ExportNode().IsNil())
	jutil.CheckState(!e.HasLocalName())
	jutil.CheckState(!e.HasModuleRequest())
	jutil.CheckState(!e.HasImportName())
	jutil.CheckState(e.NameNode().IsNil())
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
	return &cp
}
