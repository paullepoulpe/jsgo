// Port of com/google/javascript/jscomp/RemoveUnnecessarySyntheticExterns.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RemoveUnnecessarySyntheticExterns deletes name declarations from synthetic
// externs if they are in fact declared in source
//
// The VarCheck pass will create synthetic extern declarations (`var
// someName;`) for every variable name referenced without a corresponding
// externs or code declaration. Such names are "unfulfilled". VarCheck only
// does this because later compiler passes expect the following invariant to
// hold: all variable names referenced in externs or code are also declared
// in externs or code.
//
// When running optimizations based on precompiled TypedAST library shards,
// it's possible for a name to be unfulfilled in shard A, but declared in
// another shard B. Shard B "fulfills" the reference from shard A. This pass
// removes the corresponding synthetic extern that came from shard A now that
// the name is fulfilled. The purpose is to improve optimizations, e.g. to
// allow mangling the name if it is not defined in externs.
type RemoveUnnecessarySyntheticExterns struct {
	compiler      AbstractCompiler
	nodesToDetach *jutil.LinkedHashSet[*rhino.Node]
}

// Java: RemoveUnnecessarySyntheticExterns(AbstractCompiler)
func NewRemoveUnnecessarySyntheticExterns(compiler AbstractCompiler) *RemoveUnnecessarySyntheticExterns {
	return &RemoveUnnecessarySyntheticExterns{
		compiler:      compiler,
		nodesToDetach: jutil.NewLinkedHashSet[*rhino.Node](),
	}
}

// Java: RemoveUnnecessarySyntheticExterns.process(Node, Node)
func (r *RemoveUnnecessarySyntheticExterns) Process(externs *rhino.Node, root *rhino.Node) {
	globalRoot := externs.GetParent()
	scopeCreator := NewSyntacticScopeCreatorWithTreatProvidesAsRedeclarations(
		r.compiler, &RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler{outer: r}, true)

	// This call will invoke the RedeclarationCheckHandler on duplicate synthetic names
	scopeCreator.CreateScopeOfScope(globalRoot, nil /* parent= */)

	if r.nodesToDetach.IsEmpty() {
		return
	}

	for n := range r.nodesToDetach.All() {
		n.Detach()
	}

	syntheticExternsScript := r.compiler.GetSynthesizedExternsInput().GetAstRoot(r.compiler)
	r.compiler.ReportChangeToEnclosingScope(syntheticExternsScript)
}

// RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler is the handler
// for duplicate declarations.
type RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler struct {
	outer *RemoveUnnecessarySyntheticExterns
}

// Java: RemoveUnnecessarySyntheticExterns.RedeclarationCheckHandler.onRedeclaration(Scope, String, Node, CompilerInput)
func (h *RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) OnRedeclaration(s *Scope, name string, n *rhino.Node, input *CompilerInput) {
	jutil.CheckArgumentMsg2(s.IsGlobal(), "Unexpected non-global redeclaration %s in scope %s", n, s)

	isSynthetic := input == h.outer.compiler.GetSynthesizedExternsInput()
	origVar := s.GetVar(name)
	// origNode will be null for `arguments`, since there's no node that declares it.
	origNode := origVar.GetNode()

	if isSynthetic {
		// Delete duplicate unfulfilled synthetic declarations.
		// The original declaration will be deleted later iff it's also synthetic and we visit a
		// non-synthetic duplicate declaration.
		if n.GetParent().IsSynthesizedUnfulfilledNameDeclaration() {
			h.outer.nodesToDetach.Add(n.GetParent())
		}
	} else if origNode != nil &&
		origNode.GetParent().IsSynthesizedUnfulfilledNameDeclaration() {
		// Found a synthetic extern declaration that's fulfilled by some non-synthetic code.
		h.outer.nodesToDetach.Add(origNode.GetParent())
	}
}
