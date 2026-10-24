// Port of com/google/javascript/jscomp/modules/Import.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/AutoBuilder_Import_Builder.java @ closure-compiler v20261021

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

// Import is an imported name in a module.
//
// See the ImportEntry in the ECMAScript spec
// (https://www.ecma-international.org/ecma-262/9.0/index.html#table-39).
//
// moduleRequest is the module identifier of this import; importName the
// name that was imported from the requested module (for import * this will
// return "*"); localName the local name the imported value is bound to;
// modulePath the path of the containing module, if from an ES module (may
// be nil); importNode the import node for source information; nameNode the
// name node for source information.
//
// Go-only: language-adapter — the AutoBuilder-generated
// AutoBuilder_Import_Builder is folded into Import_Builder.
type Import struct {
	moduleRequest string
	importName    string
	localName     string
	modulePath    *ModuleLoader_ModulePath
	importNode    rhino.Node
	nameNode      rhino.Node
}

// Java: Import(String, String, String, ModulePath, Node, Node)
func NewImport(
	moduleRequest string,
	importName string,
	localName string,
	modulePath *ModuleLoader_ModulePath,
	importNode rhino.Node,
	nameNode rhino.Node) *Import {
	rhino.CheckNotNullMsg(importNode, "importNode")
	rhino.CheckNotNullMsg(nameNode, "nameNode")
	return &Import{
		moduleRequest: moduleRequest,
		importName:    importName,
		localName:     localName,
		modulePath:    modulePath,
		importNode:    importNode,
		nameNode:      nameNode,
	}
}

// Java: Import.moduleRequest()
func (i *Import) ModuleRequest() string { return i.moduleRequest }

// Java: Import.importName()
func (i *Import) ImportName() string { return i.importName }

// Java: Import.localName()
func (i *Import) LocalName() string { return i.localName }

// ModulePath may return nil.
//
// Java: Import.modulePath()
func (i *Import) ModulePath() *ModuleLoader_ModulePath { return i.modulePath }

// Java: Import.importNode()
func (i *Import) ImportNode() rhino.Node { return i.importNode }

// Java: Import.nameNode()
func (i *Import) NameNode() rhino.Node { return i.nameNode }

// Equals is the record's componentwise equals (ModulePath and Node compare
// by identity).
func (i *Import) Equals(o any) bool {
	other, ok := o.(*Import)
	return ok && other != nil &&
		i.moduleRequest == other.moduleRequest &&
		i.importName == other.importName &&
		i.localName == other.localName &&
		i.modulePath == other.modulePath &&
		i.importNode == other.importNode &&
		i.nameNode == other.nameNode
}

// HashCode is the record's hashCode.
//
// Go-only: language-adapter — the identity hashes of ModulePath and Node
// are jutil.IdentityHashCode.
func (i *Import) HashCode() int32 {
	h := int32(0)
	h = 31*h + jutil.StringHashCode(i.moduleRequest)
	h = 31*h + jutil.StringHashCode(i.importName)
	h = 31*h + jutil.StringHashCode(i.localName)
	if i.modulePath != nil {
		h = 31*h + jutil.IdentityHashCode(i.modulePath)
	} else {
		h = 31 * h
	}
	h = 31*h + jutil.IdentityHashCode(i.importNode)
	h = 31*h + jutil.IdentityHashCode(i.nameNode)
	return h
}

// String is the record's toString.
func (i *Import) String() string {
	modulePath := "null"
	if i.modulePath != nil {
		modulePath = i.modulePath.String()
	}
	return "Import[moduleRequest=" + i.moduleRequest +
		", importName=" + i.importName +
		", localName=" + i.localName +
		", modulePath=" + modulePath +
		", importNode=" + i.importNode.String() +
		", nameNode=" + i.nameNode.String() + "]"
}

// Import_Builder is Import.Builder with its AutoBuilder implementation.
type Import_Builder struct {
	moduleRequest *string
	importName    *string
	localName     *string
	modulePath    *ModuleLoader_ModulePath
	importNode    rhino.Node
	nameNode      rhino.Node
}

// Java: Import.builder()
func Import_Builder_() *Import_Builder {
	return newAutoBuilder_Import_Builder()
}

// Java: AutoBuilder_Import_Builder()
func newAutoBuilder_Import_Builder() *Import_Builder {
	return &Import_Builder{}
}

// Java: AutoBuilder_Import_Builder(Import)
func newAutoBuilder_Import_BuilderFrom(source *Import) *Import_Builder {
	b := &Import_Builder{}
	b.moduleRequest = &source.moduleRequest
	b.importName = &source.importName
	b.localName = &source.localName
	b.modulePath = source.modulePath
	b.importNode = source.importNode
	b.nameNode = source.nameNode
	return b
}

// Java: AutoBuilder_Import_Builder.moduleRequest(String)
func (b *Import_Builder) ModuleRequest(moduleRequest string) *Import_Builder {
	b.moduleRequest = &moduleRequest
	return b
}

// Java: AutoBuilder_Import_Builder.importName(String)
func (b *Import_Builder) ImportName(importName string) *Import_Builder {
	b.importName = &importName
	return b
}

// Java: AutoBuilder_Import_Builder.localName(String)
func (b *Import_Builder) LocalName(localName string) *Import_Builder {
	b.localName = &localName
	return b
}

// Java: AutoBuilder_Import_Builder.modulePath(ModuleLoader.ModulePath)
func (b *Import_Builder) ModulePath(modulePath *ModuleLoader_ModulePath) *Import_Builder {
	b.modulePath = modulePath
	return b
}

// Java: AutoBuilder_Import_Builder.importNode(Node)
func (b *Import_Builder) ImportNode(importNode rhino.Node) *Import_Builder {
	if importNode.IsNil() {
		panic(jutil.NewNullPointerExceptionMsg("Null importNode"))
	}
	b.importNode = importNode
	return b
}

// Java: AutoBuilder_Import_Builder.nameNode(Node)
func (b *Import_Builder) NameNode(nameNode rhino.Node) *Import_Builder {
	if nameNode.IsNil() {
		panic(jutil.NewNullPointerExceptionMsg("Null nameNode"))
	}
	b.nameNode = nameNode
	return b
}

// Java: AutoBuilder_Import_Builder.build()
func (b *Import_Builder) Build() *Import {
	if b.moduleRequest == nil ||
		b.importName == nil ||
		b.localName == nil ||
		b.importNode.IsNil() ||
		b.nameNode.IsNil() {
		missing := ""
		if b.moduleRequest == nil {
			missing += " moduleRequest"
		}
		if b.importName == nil {
			missing += " importName"
		}
		if b.localName == nil {
			missing += " localName"
		}
		if b.importNode.IsNil() {
			missing += " importNode"
		}
		if b.nameNode.IsNil() {
			missing += " nameNode"
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing))
	}
	return NewImport(
		*b.moduleRequest,
		*b.importName,
		*b.localName,
		b.modulePath,
		b.importNode,
		b.nameNode)
}
