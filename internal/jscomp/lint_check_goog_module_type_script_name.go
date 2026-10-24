// Port of com/google/javascript/jscomp/lint/CheckGoogModuleTypeScriptName.java @ closure-compiler v20261021

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckGoogModuleTypeScriptName checks that goog.module statement matches the
// generated TypeScript module namespace, which is based on the file path.
type CheckGoogModuleTypeScriptName struct {
	compiler AbstractCompiler

	finished bool
}

var CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE = DiagnosticType_Disabled(
	"JSC_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE",
	"goog.module namespace does not match the future TypeScript namespace, which is generated"+
		" from the file path."+
		" The correct namespace is: \"{0}\"")

var checkGoogModuleTypeScriptName_allowedDirectories = jutil.NewLinkedHashSet[string]()

// MOE::begin_strip
func init() {
	checkGoogModuleTypeScriptName_allowedDirectories.Add("google3/gws/")
	checkGoogModuleTypeScriptName_allowedDirectories.Add("google3/java/com/google/gws/")
	checkGoogModuleTypeScriptName_allowedDirectories.Add("google3/javascript/search/")
}

// MOE::end_strip

// NewCheckGoogModuleTypeScriptName is new CheckGoogModuleTypeScriptName(compiler).
//
// Java: CheckGoogModuleTypeScriptName(AbstractCompiler)
func NewCheckGoogModuleTypeScriptName(compiler AbstractCompiler) *CheckGoogModuleTypeScriptName {
	return &CheckGoogModuleTypeScriptName{compiler: compiler, finished: false}
}

// Java: CheckGoogModuleTypeScriptName.process(Node, Node)
func (c *CheckGoogModuleTypeScriptName) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckGoogModuleTypeScriptName.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckGoogModuleTypeScriptName) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Traverse top-level statements until a `goog.module` statement is found.
	return !c.finished &&
		(parent.IsNil() || parent.IsRoot() || parent.IsScript() || parent.IsModuleBody())
}

// Java: CheckGoogModuleTypeScriptName.visit(NodeTraversal, Node, Node)
func (c *CheckGoogModuleTypeScriptName) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if NodeUtil_IsGoogModuleCall(n) {
		c.CheckGoogModuleNamespace(t, n)
		c.finished = true
	}
}

// Java: CheckGoogModuleTypeScriptName.checkGoogModuleNamespace(NodeTraversal, Node)
func (c *CheckGoogModuleTypeScriptName) CheckGoogModuleNamespace(t *NodeTraversal, n rhino.Node) {
	originalNamespace := n.GetFirstChild().GetSecondChild().GetString()
	sourceNamePtr := t.GetSourceName()
	// Go-only: language-adapter — Node.getString() is a Go string, which
	// cannot be null, so `originalNamespace == null` is always false.
	if sourceNamePtr == nil {
		return
	}
	sourceName := *sourceNamePtr
	// MOE::begin_strip
	google3Index := jutil.StringIndexOf(sourceName, "google3")
	if google3Index == -1 {
		sourceName = "google3/" + sourceName
	} else if google3Index != 0 {
		sourceName = jutil.StringSubstringFrom(sourceName, google3Index)
	}
	// MOE::end_strip
	replacementNamespace := jutil.StringSubstring(
		strings.ReplaceAll(sourceName, "/", "."), 0, jutil.StringLength(sourceName)-jutil.StringLength(".js"))
	if originalNamespace != replacementNamespace {
		for allowedDirectory := range checkGoogModuleTypeScriptName_allowedDirectories.All() {
			if strings.HasPrefix(sourceName, allowedDirectory) {
				t.Report(n, CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE, replacementNamespace)
				return
			}
		}
	}
}
