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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Import is an imported name in a module.
//
// See the ImportEntry in the ECMAScript spec
// (https://www.ecma-international.org/ecma-262/9.0/index.html#table-39).
//
// Go-only: language-adapter — callers construct Import directly via NewImport
// rather than heap-allocating AutoBuilder_Import_Builder.
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
