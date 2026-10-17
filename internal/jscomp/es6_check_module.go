// Port of com/google/javascript/jscomp/Es6CheckModule.java @ closure-compiler v20260415

/*
 * Copyright 2017 The Closure Compiler Authors.
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

// Es6CheckModule checks that ES6 Modules are used correctly, and do not
// reference undefined keywords or features.
type Es6CheckModule struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler
}

var Es6CheckModule_ES6_MODULE_REFERENCES_THIS = DiagnosticType_Warning(
	"ES6_MODULE_REFERENCES_THIS", "The body of an ES6 module cannot reference 'this'.")

var Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED = DiagnosticType_Error(
	"JSC_IMPORT_CANNOT_BE_REASSIGNED", "Assignment to constant variable \"{0}\".")

// Java: Es6CheckModule(AbstractCompiler)
func NewEs6CheckModule(compiler AbstractCompiler) *Es6CheckModule {
	e := &Es6CheckModule{}
	e.compiler = compiler
	return e
}

// Java: Es6CheckModule.process(Node, Node)
func (e *Es6CheckModule) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(e.compiler, root, e)
}

// Java: Es6CheckModule.visit(NodeTraversal, Node, Node)
func (e *Es6CheckModule) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_THIS:
		if t.InModuleHoistScope() {
			t.Report(n, Es6CheckModule_ES6_MODULE_REFERENCES_THIS)
		}
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		if NodeUtil_IsLValue(n) &&
			!NodeUtil_IsDeclarationLValue(n) &&
			n.GetFirstChild().IsName() {
			var_ := t.GetScope().GetVar(n.GetFirstChild().GetString())
			if var_ != nil {
				nameNode := var_.GetNameNode()
				if nameNode != nil && nameNode.IsImportStar() {
					// import * as M from '';
					// M.x = 2;
					e.compiler.Report(JSError_MakeWithNode(n, Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED, nameNode.GetString()))
				}
			}
		}
	case rhino.Token_NAME:
		if NodeUtil_IsLValue(n) && !NodeUtil_IsDeclarationLValue(n) {
			var_ := t.GetScope().GetVar(n.GetString())
			if var_ != nil {
				nameNode := var_.GetNameNode()
				if nameNode != nil && nameNode != n {
					if NodeUtil_IsImportedName(nameNode) {
						// import { x } from '';
						// x = 2;
						e.compiler.Report(JSError_MakeWithNode(n, Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED, nameNode.GetString()))
					}
				}
			}
		}
	default:
	}
}
