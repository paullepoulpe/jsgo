// Port of com/google/javascript/jscomp/RewriteGlobalDeclarationsForTryCatchWrapping.java @ closure-compiler v20261015

/*
 * Copyright 2008 The Closure Compiler Authors.
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

// RewriteGlobalDeclarationsForTryCatchWrapping moves top-level function
// declarations to the top of the enclosing JSChunk and rewrites class
// declarations.
//
// Enable this pass if a try catch block wraps the output after compilation,
// and the output runs on Firefox because function declarations are only
// defined when reached inside a try catch block on older versions of
// Firefox.
//
// On Firefox versions <= 45 (March 2016), this code works:
//
// var g = f; function f() {}
//
// but this code does not work:
//
// try { var g = f; function f() {} } catch(e) {}
//
// NOTE: As of Firefox version 46 the above code works. Projects not
// supporting older versions of Firefox shouldn't need this pass.
//
// RescopeGlobalSymbols still depends on this pass running first to preserve
// function hoisting semantics.
//
// NOTE: This pass is safe to turn on by default and delete the associated
// compiler option. However, we don't do that because the pass is only
// useful for code wrapped in a try/catch, and otherwise it makes debugging
// harder because it moves code around.
//
// NOTE: now that this pass also moves class declarations it should be
// renamed along with the associated options.
type RewriteGlobalDeclarationsForTryCatchWrapping struct {
	compiler AbstractCompiler
	// Go-only: language-adapter — Java's ArrayListMultimap is backed by a
	// HashMap keyed by JSChunk, which has identity hashing, so Java's
	// iteration order over the chunks is not deterministic; each chunk has
	// its own insertion root, so the order is not observable. Insertion order
	// is used here.
	functions *jutil.LinkedHashMap[*JSChunk, *jutil.ArrayList[rhino.Node]]
	classes   *jutil.ArrayList[rhino.Node]
}

// Java: RewriteGlobalDeclarationsForTryCatchWrapping(AbstractCompiler)
func NewRewriteGlobalDeclarationsForTryCatchWrapping(compiler AbstractCompiler) *RewriteGlobalDeclarationsForTryCatchWrapping {
	r := &RewriteGlobalDeclarationsForTryCatchWrapping{
		functions: jutil.NewLinkedHashMap[*JSChunk, *jutil.ArrayList[rhino.Node]](),
		classes:   jutil.NewArrayList[rhino.Node](),
	}
	r.compiler = compiler
	return r
}

// Java: RewriteGlobalDeclarationsForTryCatchWrapping.process(Node, Node)
func (r *RewriteGlobalDeclarationsForTryCatchWrapping) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, r)
	for chunk, values := range r.functions.All() {
		addingRoot := r.compiler.GetNodeForCodeInsertion(chunk)
		fnNodes := values.Slice()
		if len(fnNodes) != 0 {
			// Lists.reverse(entry.getValue())
			for i := len(fnNodes) - 1; i >= 0; i-- {
				n := fnNodes[i]
				parent := n.GetParent()
				n.Detach()
				r.compiler.ReportChangeToEnclosingScope(parent)

				nameNode := n.GetFirstChild()
				name := nameNode.GetString()
				nameNode.SetString("")
				addingRoot.AddChildToFront(rhino.IR_VarWithValue(rhino.IR_Name(name), n).SrcrefTreeIfMissing(n))
				r.compiler.ReportChangeToEnclosingScope(nameNode)
			}
			r.compiler.ReportChangeToEnclosingScope(addingRoot)
		}
	}

	for _, n := range r.classes.Slice() {
		// rewrite CLASS > NAME ... => VAR > NAME > CLASS > EMPTY ...
		originalNameNode := n.GetFirstChild()
		originalNameNode.ReplaceWith(rhino.IR_Empty())
		var_ := rhino.IR_Var(originalNameNode)
		n.ReplaceWith(var_)
		originalNameNode.AddChildToFront(n)
		r.compiler.ReportChangeToEnclosingScope(n)
	}
}

// Java: RewriteGlobalDeclarationsForTryCatchWrapping.shouldTraverse(NodeTraversal, Node, Node)
func (r *RewriteGlobalDeclarationsForTryCatchWrapping) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	grandparent := n.GetAncestor(2)
	return grandparent.IsNil() || !grandparent.IsScript()
}

// Java: RewriteGlobalDeclarationsForTryCatchWrapping.visit(NodeTraversal, Node, Node)
func (r *RewriteGlobalDeclarationsForTryCatchWrapping) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if parent.IsNil() || !parent.IsScript() {
		return
	}

	if NodeUtil_IsFunctionDeclaration(n) {
		chunk := t.GetChunk()
		list := r.functions.Get(chunk)
		if list == nil {
			list = jutil.NewArrayList[rhino.Node]()
			r.functions.Put(chunk, list)
		}
		list.Add(n)
	}

	if NodeUtil_IsClassDeclaration(n) {
		r.classes.Add(n)
	}

	if n.IsConst() || n.IsLet() {
		n.SetToken(rhino.Token_VAR)
		r.compiler.ReportChangeToEnclosingScope(n)
	}
}
