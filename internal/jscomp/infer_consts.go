// Port of com/google/javascript/jscomp/InferConsts.java @ closure-compiler v20260415

/*
 * Copyright 2014 The Closure Compiler Authors.
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

// InferConsts attaches the CONST_VAR annotation to any variable that's
// 1) Provably well-defined and assigned once in its lifetime.
// 2) Annotated 'const'
// 3) Declared with the 'const' keyword.
//
// These 3 are considered semantically equivalent. Notice that a variable in
// a loop is never considered const.
//
// Note that criteria (1) is only used for normal code, not externs.
type InferConsts struct {
	compiler AbstractCompiler
}

// Java: InferConsts(AbstractCompiler)
func NewInferConsts(compiler AbstractCompiler) *InferConsts {
	return &InferConsts{compiler: compiler}
}

// Java: InferConsts.process(Node, Node)
func (c *InferConsts) Process(externs *rhino.Node, js *rhino.Node) {
	collector := NewReferenceCollector(
		c.compiler, ReferenceCollector_DO_NOTHING_BEHAVIOR, NewSyntacticScopeCreator(c.compiler))
	collector.ProcessRoot(js)

	for v := range collector.GetAllSymbols() {
		c.ConsiderVar(v, collector.GetReferences(v))
	}

	globalExternsScope := NewSyntacticScopeCreator(c.compiler).CreateScopeOfScope(externs, nil)
	for v := range globalExternsScope.GetAllSymbols() {
		c.ConsiderVar(v, nil)
	}
}

// Java: InferConsts.considerVar(Var, ReferenceCollection)
func (c *InferConsts) ConsiderVar(v *Var, refCollection *ReferenceCollection) {
	if v.IsImplicitGoogNamespace() {
		return // no name node for provided variables.
	}
	nameNode := v.GetNameNode()
	docInfo := v.GetJSDocInfo()
	if docInfo != nil && docInfo.IsConstant() {
		nameNode.SetDeclaredConstantVar(true)
	} else if nameNode != nil && nameNode.GetParent().IsConst() {
		nameNode.SetDeclaredConstantVar(true)
	}
	if InferConsts_IsInferredConst(v, refCollection) {
		nameNode.SetInferredConstantVar(true)
	}
}

// Java: InferConsts.isInferredConst(Var, ReferenceCollection)
func InferConsts_IsInferredConst(v *Var, refCollection *ReferenceCollection) bool {
	nameNode := v.GetNameNode()
	if nameNode == nil || refCollection == nil || !refCollection.IsAssignedOnceInLifetime() {
		return false
	}
	if v.IsImplicitGoogNamespace() {
		return false
	}
	declarationType := v.DeclarationType()

	// A switch on a null Token throws NullPointerException.
	switch *declarationType {
	case rhino.Token_LET:
		// Check that non-destructuring let names are actually assigned at declaration.
		return !nameNode.GetParent().IsLet() || nameNode.HasChildren()
	case rhino.Token_CONST,
		rhino.Token_CATCH,
		rhino.Token_CLASS,
		rhino.Token_PARAM_LIST, // Parameters cannot be referenced before declaration.
		rhino.Token_FUNCTION:   // Function hoisting means no references before declaration.
		return true
	case rhino.Token_IMPORT:
		// ES module exports are mutable.
		// TODO(lharker): make this check smarter if we start optimizing unrewritten modules.
		return false
	case rhino.Token_VAR:
		// var hoisting requires this extra work to make sure the 'declaration' is also the first
		// reference.
		return refCollection.FirstReferenceIsAssigningDeclaration() && refCollection.IsWellDefined()
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unrecognized declaration type " + declarationType.String()))
	}
}
