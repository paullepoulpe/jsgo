// Port of com/google/javascript/jscomp/ForbidDynamicImportUsage.java @ closure-compiler v20261021

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ForbidDynamicImportUsage warns at any usage of Dynamic Import expressions
// that they are unable to be transpiled.
type ForbidDynamicImportUsage struct {
	compiler AbstractCompiler
}

// Java: ForbidDynamicImportUsage.DYNAMIC_IMPORT_USAGE
var ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE = DiagnosticType_Error(
	"JSC_DYNAMIC_IMPORT_USAGE", "Dynamic import expressions cannot be transpiled.")

// Java: ForbidDynamicImportUsage(AbstractCompiler)
func NewForbidDynamicImportUsage(compiler AbstractCompiler) *ForbidDynamicImportUsage {
	f := &ForbidDynamicImportUsage{}
	f.compiler = compiler
	return f
}

// Java: ForbidDynamicImportUsage.process(Node, Node)
func (f *ForbidDynamicImportUsage) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(f.compiler, root, f)
	f.compiler.MarkFeatureNotAllowed(parser.FeatureSet_Feature_DYNAMIC_IMPORT)
}

// Java: ForbidDynamicImportUsage.shouldTraverse(NodeTraversal, Node, Node)
func (f *ForbidDynamicImportUsage) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		// We only want to traverse scripts that contain dynamic imports, so we can generate error
		// messages that point to them.
		return f.GetFeatureSetOfScript(n).ContainsFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT)
	} else {
		return true
	}
}

// Java: ForbidDynamicImportUsage.getFeatureSetOfScript(Node)
func (f *ForbidDynamicImportUsage) GetFeatureSetOfScript(n rhino.Node) *parser.FeatureSet {
	featureSet := NodeUtil_GetFeatureSetOfScript(n)
	if featureSet == nil {
		return parser.FeatureSet_BARE_MINIMUM
	}
	return featureSet
}

// Java: ForbidDynamicImportUsage.visit(NodeTraversal, Node, Node)
func (f *ForbidDynamicImportUsage) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_DYNAMIC_IMPORT:
		t.Report(n, ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE)
	default:
	}
}
