// Port of com/google/javascript/jscomp/Es6RelativizeImportPaths.java @ closure-compiler v20261015

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RelativizeImportPaths rewrites ES6 import paths to be relative after
// resolving according to the compiler's module resolver.
//
// Useful for servers that wish to preserve ES6 modules, meaning their paths
// need to be valid in the browser.
type Es6RelativizeImportPaths struct {
	compiler AbstractCompiler
}

// Java: Es6RelativizeImportPaths(AbstractCompiler)
func NewEs6RelativizeImportPaths(compiler AbstractCompiler) *Es6RelativizeImportPaths {
	e := &Es6RelativizeImportPaths{}
	e.compiler = compiler
	return e
}

// Java: Es6RelativizeImportPaths.process(Node, Node)
func (e *Es6RelativizeImportPaths) Process(externs *rhino.Node, root *rhino.Node) {
	for script := root.GetFirstChild(); script != nil; script = script.GetNext() {
		if Es6RewriteModules_IsEs6ModuleRoot(script) {
			NodeTraversal_Traverse(e.compiler, script, NewEs6RelativizeImportPaths_Rewriter())
		}
	}
}

// Es6RelativizeImportPaths_Rewriter is the private static nested class
// Es6RelativizeImportPaths.Rewriter.
type Es6RelativizeImportPaths_Rewriter struct {
	NodeTraversal_AbstractPreOrderCallbackBase
}

// Java: Es6RelativizeImportPaths.Rewriter()
func NewEs6RelativizeImportPaths_Rewriter() *Es6RelativizeImportPaths_Rewriter {
	return &Es6RelativizeImportPaths_Rewriter{}
}

// Java: Es6RelativizeImportPaths.Rewriter.shouldTraverse(NodeTraversal, Node, Node)
func (r *Es6RelativizeImportPaths_Rewriter) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_MODULE_BODY, rhino.Token_SCRIPT:
		return true
	case rhino.Token_IMPORT:
		r.VisitImport(nodeTraversal, n)
		return false
	case rhino.Token_EXPORT:
		r.VisitExport(nodeTraversal, n)
		return false
	default:
		return false
	}
}

// Java: Es6RelativizeImportPaths.Rewriter.visitImport(NodeTraversal, Node)
func (r *Es6RelativizeImportPaths_Rewriter) VisitImport(t *NodeTraversal, importDecl *rhino.Node) {
	specifierNode := importDecl.GetLastChild()
	specifier := specifierNode.GetString()

	// Leave relative and truly absolute paths (those with a scheme) alone. Only transform
	// absolute paths without a scheme (starting with "/") and ambiguous paths.
	if ModuleLoader_IsRelativeIdentifier(specifier) {
		return
	} else {
		returned := false
		jutil.Try(func() {
			specifierURI := jutil.NewURI(specifier)
			if specifierURI.IsAbsolute() {
				returned = true
			}
		}, func(e *jutil.URISyntaxException) {
			returned = true
		})
		if returned {
			return
		}
	}

	scriptPath := t.GetInput().GetPath().String()

	returned := false
	jutil.Try(func() {
		// If the script path has a scheme / host / port then just use the path part.
		scriptPath = *jutil.NewURI(scriptPath).GetPath()
	}, func(e *jutil.URISyntaxException) {
		returned = true
	})
	if returned {
		return
	}

	newSpecifier := t.GetInput().GetPath().ResolveModuleAsPath(specifier).String()

	// If a module root is stripped then this won't start with "/" when it probably should.
	if !strings.HasPrefix(newSpecifier, "/") {
		newSpecifier = "/" + newSpecifier
	}

	newSpecifier = jutil.Path_Of(scriptPath).GetParent().Relativize(jutil.Path_Of(newSpecifier)).String()

	// Relativizing two paths with the same directory yields an ambiguous path rather than one
	// starting with "./".
	if ModuleLoader_IsAmbiguousIdentifier(newSpecifier) {
		newSpecifier = "./" + newSpecifier
	}

	if newSpecifier != specifier {
		specifierNode.SetString(newSpecifier)
		t.ReportCodeChangeWithNode(specifierNode)
	}
}

// Java: Es6RelativizeImportPaths.Rewriter.visitExport(NodeTraversal, Node)
func (r *Es6RelativizeImportPaths_Rewriter) VisitExport(t *NodeTraversal, export *rhino.Node) {
	if export.HasTwoChildren() {
		// export from
		r.VisitImport(t, export)
	}
}
