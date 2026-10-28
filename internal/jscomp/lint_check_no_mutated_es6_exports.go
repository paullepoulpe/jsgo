// Port of com/google/javascript/jscomp/lint/CheckNoMutatedEs6Exports.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// CheckNoMutatedEs6Exports checks that exports of ES6 modules are not mutated
// outside of module initialization.
type CheckNoMutatedEs6Exports struct {
	compiler AbstractCompiler
	// Java's MultimapBuilder.hashKeys().hashSetValues().build() sizes a
	// java.util.HashMap for 8 expected keys (initial capacity 11 -> 16),
	// whose bucket order determines MUTATED_EXPORT report order.
	mutatedNames       *jutil.JavaHashMap[string, []rhino.Node]
	exportedLocalNames map[uint32]struct{}
}

var CheckNoMutatedEs6Exports_MUTATED_EXPORT = DiagnosticType_Warning(
	"JSC_MUTATED_EXPORT",
	"The name \"{0}\" is exported and should not be mutated outside of module "+
		"initialization. Mutable exports are generally difficult to reason about. You "+
		"can work around this by exporting getter/setter functions, or an object with "+
		"mutable properties instead.")

func NewCheckNoMutatedEs6Exports(compiler AbstractCompiler) *CheckNoMutatedEs6Exports {
	return &CheckNoMutatedEs6Exports{
		compiler:           compiler,
		mutatedNames:       jutil.NewJavaHashMapWithCapacity[string, []rhino.Node](11),
		exportedLocalNames: make(map[uint32]struct{}),
	}
}

func (c *CheckNoMutatedEs6Exports) CheckNoMutations() {
	for mutatedExport, mutations := range c.mutatedNames.All() {
		if _, ok := c.exportedLocalNames[uint32(strpool.Lookup(mutatedExport))]; !ok {
			continue
		}
		for _, mutation := range mutations {
			c.compiler.Report(JSError_MakeWithNode(mutation, CheckNoMutatedEs6Exports_MUTATED_EXPORT, mutatedExport))
		}
	}

	c.mutatedNames.Clear()
	clear(c.exportedLocalNames)
}

func (c *CheckNoMutatedEs6Exports) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsModuleBody() {
		c.CheckNoMutations()
	}
}

func (c *CheckNoMutatedEs6Exports) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		return n.GetBooleanProp(rhino.Node_ES6_MODULE)
	case rhino.Token_EXPORT:
		c.VisitExport(n)
	case rhino.Token_NAME:
		c.VisitName(t, n)
	}
	return true
}

func (c *CheckNoMutatedEs6Exports) VisitExport(export rhino.Node) {
	if !export.HasOneChild() {
		return
	}
	first := export.GetFirstChild()
	if first.IsExportSpecs() {
		for exportSpec := first.GetFirstChild(); !exportSpec.IsNil(); exportSpec = exportSpec.GetNext() {
			jutil.CheckState(exportSpec.HasTwoChildren())
			c.exportedLocalNames[exportSpec.GetFirstChild().GetStringID()] = struct{}{}
		}
		return
	}
	if export.GetBooleanProp(rhino.Node_EXPORT_ALL_FROM) {
		return
	}
	if NodeUtil_IsNameDeclaration(first) {
		NodeUtil_VisitLhsNodesInNode(first, func(lhs rhino.Node) {
			jutil.CheckState(lhs.IsName())
			c.exportedLocalNames[lhs.GetStringID()] = struct{}{}
		})
	} else if export.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
		if !first.IsClass() && !first.IsFunction() {
			return
		}
		nameNode := first.GetFirstChild()
		if !nameNode.IsEmpty() && nameNode.GetString() != "" {
			c.exportedLocalNames[nameNode.GetStringID()] = struct{}{}
		}
	} else {
		jutil.CheckState(first.IsClass() || first.IsFunction())
		c.exportedLocalNames[first.GetFirstChild().GetStringID()] = struct{}{}
	}
}

func (c *CheckNoMutatedEs6Exports) VisitName(t *NodeTraversal, name rhino.Node) {
	if !NodeUtil_IsLValue(name) {
		return
	}
	scope := t.GetScope()
	if scope.GetClosestHoistScope().IsModuleScope() {
		return
	}
	id := name.GetStringID()
	var_ := scope.GetVarByID(id)
	if var_ == nil || !var_.GetScope().(*Scope).IsModuleScope() {
		return
	}
	key := name.GetString()
	nodes, _ := c.mutatedNames.Lookup(key)
	if !slices.Contains(nodes, name) {
		c.mutatedNames.Put(key, append(nodes, name))
	}
}

func (c *CheckNoMutatedEs6Exports) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}
