// Port of com/google/javascript/jscomp/lint/CheckNoMutatedEs6Exports.java @ closure-compiler v20260415

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

// CheckNoMutatedEs6Exports checks that exports of ES6 modules are not mutated
// outside of module initialization.
type CheckNoMutatedEs6Exports struct {
	compiler AbstractCompiler
	// Go-only: language-adapter — Java HashMap order: the Java field is
	// MultimapBuilder.hashKeys().hashSetValues().build(), whose keys iterate
	// in java.util.HashMap order (String hash buckets) and whose Node values
	// iterate in identity-hash order (not deterministic in Java itself).
	// Here keys and values iterate in insertion order. The order is
	// observable only in the order of MUTATED_EXPORT reports within one
	// module body.
	mutatedNames       *jutil.SetMultimap[string, *rhino.Node]
	exportedLocalNames *jutil.LinkedHashSet[string]
}

var CheckNoMutatedEs6Exports_MUTATED_EXPORT = DiagnosticType_Warning(
	"JSC_MUTATED_EXPORT",
	"The name \"{0}\" is exported and should not be mutated outside of module "+
		"initialization. Mutable exports are generally difficult to reason about. You "+
		"can work around this by exporting getter/setter functions, or an object with "+
		"mutable properties instead.")

// NewCheckNoMutatedEs6Exports is new CheckNoMutatedEs6Exports(compiler).
//
// Java: CheckNoMutatedEs6Exports(AbstractCompiler)
func NewCheckNoMutatedEs6Exports(compiler AbstractCompiler) *CheckNoMutatedEs6Exports {
	return &CheckNoMutatedEs6Exports{
		compiler:           compiler,
		mutatedNames:       jutil.NewSetMultimap[string, *rhino.Node](),
		exportedLocalNames: jutil.NewLinkedHashSet[string](),
	}
}

// Java: CheckNoMutatedEs6Exports.checkNoMutations()
func (c *CheckNoMutatedEs6Exports) CheckNoMutations() {
	// Sets.intersection(mutatedNames.keySet(), exportedLocalNames): the keys of
	// mutatedNames, in its order, that are also in exportedLocalNames.
	for mutatedExport := range c.mutatedNames.KeySet() {
		if !c.exportedLocalNames.Contains(mutatedExport) {
			continue
		}
		for mutation := range c.mutatedNames.Get(mutatedExport).All() {
			c.compiler.Report(JSError_MakeWithNode(mutation, CheckNoMutatedEs6Exports_MUTATED_EXPORT, mutatedExport))
		}
	}

	c.mutatedNames.Clear()
	c.exportedLocalNames.Clear()
}

// Java: CheckNoMutatedEs6Exports.visit(NodeTraversal, Node, Node)
func (c *CheckNoMutatedEs6Exports) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsModuleBody() {
		c.CheckNoMutations()
	}
}

// Java: CheckNoMutatedEs6Exports.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckNoMutatedEs6Exports) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		return n.GetBooleanProp(rhino.Node_ES6_MODULE)
	case rhino.Token_EXPORT:
		c.VisitExport(n)
		return true
	case rhino.Token_NAME:
		c.VisitName(t, n)
		return true
	default:
		return true
	}
}

// Java: CheckNoMutatedEs6Exports.visitExport(Node)
func (c *CheckNoMutatedEs6Exports) VisitExport(export *rhino.Node) {
	if export.HasOneChild() && export.GetFirstChild().IsExportSpecs() {
		// export {a, b as c};
		for exportSpec := export.GetFirstFirstChild(); exportSpec != nil; exportSpec = exportSpec.GetNext() {
			jutil.CheckState(exportSpec.HasTwoChildren())
			c.exportedLocalNames.Add(exportSpec.GetFirstChild().GetString())
		}
	} else if export.HasOneChild() && !export.GetBooleanProp(rhino.Node_EXPORT_ALL_FROM) {
		declaration := export.GetFirstChild()

		if NodeUtil_IsNameDeclaration(declaration) {
			// export const x = 0;
			// export let a, b, c;
			NodeUtil_VisitLhsNodesInNode(
				declaration,
				func(lhs *rhino.Node) {
					jutil.CheckState(lhs.IsName())
					c.exportedLocalNames.Add(lhs.GetString())
				})
		} else if export.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
			// export default function() {}
			// export default function foo() {}
			// export default 0;
			// Note: default exports may or may not be declarations.
			if !declaration.IsClass() && !declaration.IsFunction() {
				return
			}
			nameNode := declaration.GetFirstChild()
			if !nameNode.IsEmpty() && nameNode.GetString() != "" {
				c.exportedLocalNames.Add(nameNode.GetString())
			}
		} else {
			// export function foo() {}
			// export class Bar {}
			jutil.CheckState(declaration.IsClass() || declaration.IsFunction())
			nameNode := declaration.GetFirstChild()
			c.exportedLocalNames.Add(nameNode.GetString())
		}
	}
}

// Java: CheckNoMutatedEs6Exports.visitName(NodeTraversal, Node)
func (c *CheckNoMutatedEs6Exports) VisitName(t *NodeTraversal, name *rhino.Node) {
	scope := t.GetScope()
	if NodeUtil_IsLValue(name) && !scope.GetClosestHoistScope().IsModuleScope() {
		var_ := scope.GetVar(name.GetString())
		if var_ != nil && var_.GetScope().(*Scope).IsModuleScope() {
			c.mutatedNames.Put(name.GetString(), name)
		}
	}
}

// Java: CheckNoMutatedEs6Exports.process(Node, Node)
func (c *CheckNoMutatedEs6Exports) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}
