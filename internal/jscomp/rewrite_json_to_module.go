// Port of com/google/javascript/jscomp/RewriteJsonToModule.java @ closure-compiler v20261021

/*
 * Copyright 2016 The Closure Compiler Authors.
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

// RewriteJsonToModule rewrites a JSON file to be a module export. So that the
// JSON file parses correctly, it is wrapped in an EXPR_RESULT. The pass makes
// only basic checks that the file provided is valid JSON. It is not a full
// JSON validator.
//
// Looks for JSON files named "package.json" so that the "main" property can
// be used as an alias in module name resolution. See
// https://docs.npmjs.com/files/package.json#main
type RewriteJsonToModule struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	packageJsonMainEntries *jutil.LinkedHashMap[string, string]
	compiler               AbstractCompiler
}

var RewriteJsonToModule_JSON_UNEXPECTED_TOKEN = DiagnosticType_Error("JSC_JSON_UNEXPECTED_TOKEN", "Unexpected JSON token")

// NewRewriteJsonToModule creates a new RewriteJsonToModule instance which
// can be used to rewrite JSON files to modules.
//
// Java: RewriteJsonToModule(AbstractCompiler)
func NewRewriteJsonToModule(compiler AbstractCompiler) *RewriteJsonToModule {
	r := &RewriteJsonToModule{}
	r.compiler = compiler
	r.packageJsonMainEntries = jutil.NewLinkedHashMap[string, string]()
	return r
}

// Java: RewriteJsonToModule.getPackageJsonMainEntries()
func (r *RewriteJsonToModule) GetPackageJsonMainEntries() *jutil.LinkedHashMap[string, string] {
	return jutil.NewLinkedHashMapFrom(r.packageJsonMainEntries)
}

// Process: module rewriting is done a on per-file basis prior to main
// compilation. The root node for each file is a SCRIPT - not the typical
// jsRoot of other passes.
//
// Java: RewriteJsonToModule.process(Node, Node)
func (r *RewriteJsonToModule) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(root.IsScript())
	NodeTraversal_Traverse(r.compiler, root, r)
}

// Java: RewriteJsonToModule.visit(NodeTraversal, Node, Node)
func (r *RewriteJsonToModule) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if !n.HasOneChild() {
			r.compiler.Report(JSError_MakeWithNode(n, RewriteJsonToModule_JSON_UNEXPECTED_TOKEN))
		} else {
			r.VisitScript(t, n)
		}
		return
	case rhino.Token_OBJECTLIT, rhino.Token_ARRAYLIT, rhino.Token_NUMBER, rhino.Token_TRUE, rhino.Token_FALSE, rhino.Token_NULL, rhino.Token_STRINGLIT:
	case rhino.Token_STRING_KEY:
		if !n.IsQuotedStringKey() || !n.HasOneChild() {
			r.compiler.Report(JSError_MakeWithNode(n, RewriteJsonToModule_JSON_UNEXPECTED_TOKEN))
		}
	case rhino.Token_EXPR_RESULT:
		if !parent.IsScript() {
			r.compiler.Report(JSError_MakeWithNode(n, RewriteJsonToModule_JSON_UNEXPECTED_TOKEN))
		}
	default:
		r.compiler.Report(JSError_MakeWithNode(n, RewriteJsonToModule_JSON_UNEXPECTED_TOKEN))
	}

	if n.GetLineno() == 1 {
		// We wrapped the expression in parens so our first-line columns are off by one.
		// We need to correct for this.
		n.SetLinenoCharno(1, n.GetCharno()-1)
		t.ReportCodeChange()
	}
}

// VisitScript: for script nodes of JSON objects, add a module variable
// assignment so the result is exported.
//
// If the file path ends with "/package.json", look for main entries in their
// specified order in the object literal and track them as module aliases.
// Main entries default to "main" and can be overridden with the
// `--package_json_entry_names` option.
//
// Java: RewriteJsonToModule.visitScript(NodeTraversal, Node)
func (r *RewriteJsonToModule) VisitScript(t *NodeTraversal, n rhino.Node) {
	if !n.HasOneChild() || !n.GetFirstChild().IsExprResult() {
		r.compiler.Report(JSError_MakeWithNode(n, RewriteJsonToModule_JSON_UNEXPECTED_TOKEN))
		return
	}

	jsdoc := rhino.JSDocInfo_Builder_()
	jsdoc.RecordFileOverview("Suppresses undefined var goog error")
	jsdoc.RecordSuppression("undefinedVars")
	n.SetJSDocInfo(jsdoc.Build())

	jsonObject := n.GetFirstFirstChild().Detach()
	n.RemoveFirstChild()

	moduleName := t.GetInput().GetPath().ToModuleName()

	n.AddChildToFront(
		rhino.IR_VarWithValue(rhino.IR_Name(moduleName).Srcref(jsonObject), jsonObject).Srcref(jsonObject))

	n.AddChildToFront(
		rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name("goog"), "provide"), rhino.IR_String(moduleName))).
			SrcrefTreeIfMissing(n))

	inputPath := t.GetInput().GetSourceFile().GetName()
	if jutil.StringEndsWith(inputPath, "/package.json") && jsonObject.IsObjectLit() {
		possibleMainEntries := r.compiler.GetOptions().GetPackageJsonEntryNames()

		for _, entryName := range possibleMainEntries {
			entry := NodeUtil_GetFirstPropMatchingKey(jsonObject, entryName)

			if !entry.IsNil() && (entry.IsStringLit() || entry.IsObjectLit()) {
				dirName := jutil.StringSubstring(inputPath, 0, jutil.StringLength(inputPath)-jutil.StringLength("package.json"))

				if entry.IsStringLit() {
					r.packageJsonMainEntries.Put(inputPath, dirName+entry.GetString())
					break
				} else if entry.IsObjectLit() {
					jutil.CheckStateMsg1(entryName == "browser", "%s", entryName)

					// don't break if we're processing a browser field that is an object literal
					// because one of its entries may override the package.json main, which
					// we will get in the next iteration.
					r.ProcessBrowserFieldAdvancedUsage(dirName, entry)
				}
			}
		}
	}

	t.ReportCodeChange()
}

// ProcessBrowserFieldAdvancedUsage: for browser field entries in
// package.json files that are used in an advanced manner
// (https://github.com/defunctzombie/package-browser-field-spec/#replace-specific-files---advanced),
// track the entries in that object literal as module file replacements.
//
// Java: RewriteJsonToModule.processBrowserFieldAdvancedUsage(String, Node)
func (r *RewriteJsonToModule) ProcessBrowserFieldAdvancedUsage(dirName string, entry rhino.Node) {
	for child := entry.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		value := child.GetFirstChild()

		jutil.CheckState(child.IsStringKey() && (value.IsStringLit() || value.IsFalse()))

		path := child.GetString()

		if jutil.StringStartsWith(path, ModuleLoader_DEFAULT_FILENAME_PREFIX) {
			path = jutil.StringSubstringFrom(path, jutil.StringLength(ModuleLoader_DEFAULT_FILENAME_PREFIX))
		}

		var replacement string
		if value.IsStringLit() {
			replacement = dirName + value.GetString()
		} else {
			replacement = ModuleLoader_JSC_BROWSER_SKIPLISTED_MARKER
		}

		r.packageJsonMainEntries.Put(dirName+path, replacement)
	}
}
