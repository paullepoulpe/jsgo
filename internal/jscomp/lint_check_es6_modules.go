// Port of com/google/javascript/jscomp/lint/CheckEs6Modules.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckEs6Modules: miscellaneous checks for style in ES6 modules.
type CheckEs6Modules struct {
	compiler         AbstractCompiler
	importSpecifiers *jutil.LinkedHashMap[string, rhino.Node]
}

var CheckEs6Modules_DUPLICATE_IMPORT = DiagnosticType_Warning(
	"JSC_DUPLICATE_IMPORT", "The module \"{0}\" has already been imported at {1}, {2}.")

var CheckEs6Modules_NO_DEFAULT_EXPORT = DiagnosticType_Warning(
	"JSC_DEFAULT_EXPORT",
	"Do not use the default export. There is no way to force consistent naming when "+
		"imported.")

// NewCheckEs6Modules is new CheckEs6Modules(compiler).
//
// Java: CheckEs6Modules(AbstractCompiler)
func NewCheckEs6Modules(compiler AbstractCompiler) *CheckEs6Modules {
	return &CheckEs6Modules{
		compiler:         compiler,
		importSpecifiers: jutil.NewLinkedHashMap[string, rhino.Node](),
	}
}

// Java: CheckEs6Modules.visit(NodeTraversal, Node, Node)
func (c *CheckEs6Modules) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsModuleBody() {
		c.importSpecifiers.Clear()
	}
}

// Java: CheckEs6Modules.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckEs6Modules) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_MODULE_BODY:
		return true
	case rhino.Token_SCRIPT:
		return n.GetBooleanProp(rhino.Node_ES6_MODULE)
	case rhino.Token_IMPORT:
		c.VisitImport(t, n)
		return false
	case rhino.Token_EXPORT:
		c.VisitExport(t, n)
		return false
	default:
		return false
	}
}

// Java: CheckEs6Modules.visitImport(NodeTraversal, Node)
func (c *CheckEs6Modules) VisitImport(t *NodeTraversal, importNode rhino.Node) {
	specifier := importNode.GetLastChild().GetString()
	duplicateImport := c.importSpecifiers.PutIfAbsent(specifier, importNode)
	if !duplicateImport.IsNil() {
		t.Report(
			importNode,
			CheckEs6Modules_DUPLICATE_IMPORT,
			specifier,
			strconv.Itoa(duplicateImport.GetLineno()),
			strconv.Itoa(duplicateImport.GetCharno()))
	}
}

// Java: CheckEs6Modules.visitExport(NodeTraversal, Node)
func (c *CheckEs6Modules) VisitExport(t *NodeTraversal, export rhino.Node) {
	if export.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
		t.Report(export, CheckEs6Modules_NO_DEFAULT_EXPORT)
	}
}

// Java: CheckEs6Modules.process(Node, Node)
func (c *CheckEs6Modules) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}
